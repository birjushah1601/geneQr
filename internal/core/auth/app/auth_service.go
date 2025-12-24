package app

import (
	"context"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
)

// AuthService orchestrates authentication business logic
type AuthService struct {
	userRepo     domain.UserRepository
	otpService   *OTPService
	jwtService   *JWTService
	pwdService   *PasswordService
	auditRepo    domain.AuditRepository
	orgRepo      domain.OrganizationRepository
	config       *AuthConfig
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	MaxFailedAttempts int           // Default: 5
	LockoutDuration   time.Duration // Default: 30 minutes
	AllowRegistration bool          // Default: true
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo domain.UserRepository,
	otpService *OTPService,
	jwtService *JWTService,
	pwdService *PasswordService,
	auditRepo domain.AuditRepository,
	orgRepo domain.OrganizationRepository,
	config *AuthConfig,
) *AuthService {
	if config == nil {
		config = &AuthConfig{
			MaxFailedAttempts: 5,
			LockoutDuration:   30 * time.Minute,
			AllowRegistration: true,
		}
	}

	return &AuthService{
		userRepo:   userRepo,
		otpService: otpService,
		jwtService: jwtService,
		pwdService: pwdService,
		auditRepo:  auditRepo,
		orgRepo:    orgRepo,
		config:     config,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if !s.config.AllowRegistration {
		return nil, fmt.Errorf("registration is currently disabled")
	}

	// Validate password strength
	if req.Password != "" {
		if err := s.pwdService.ValidatePasswordStrength(req.Password); err != nil {
			return nil, fmt.Errorf("weak password: %w", err)
		}
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmailOrPhone(ctx, req.Identifier)
	if existingUser != nil {
		return nil, fmt.Errorf("user already exists with this email/phone")
	}

	// Create user
	user := &domain.User{
		ID:                 uuid.New(),
		FullName:           req.FullName,
		PreferredAuthMethod: req.PreferredAuthMethod,
		Status:             domain.UserStatusPending,
	}

	// Set email or phone
	if isEmail(req.Identifier) {
		user.Email = &req.Identifier
	} else {
		user.Phone = &req.Identifier
	}

	// Hash password if provided
	if req.Password != "" {
		hash, err := s.pwdService.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = &hash
	}

	// Create user record
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send OTP for verification
	deliveryMethod := domain.DeliveryMethodEmail
	if user.Phone != nil {
		deliveryMethod = domain.DeliveryMethodSMS
	}

	otpResp, err := s.otpService.SendOTP(ctx, &SendOTPRequest{
		Identifier:     req.Identifier,
		UserID:         &user.ID,
		DeliveryMethod: deliveryMethod,
		Purpose:        domain.OTPPurposeVerify,
		DeviceInfo:     req.DeviceInfo,
		IPAddress:      req.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send verification OTP: %w", err)
	}

	// Log registration
	s.logAudit(ctx, &user.ID, domain.AuditActionRegister, true, req.IPAddress, "")

	return &RegisterResponse{
		UserID:      user.ID,
		RequiresOTP: true,
		OTPSentTo:   otpResp.DeliveredTo,
		ExpiresIn:   otpResp.ExpiresIn,
	}, nil
}

// LoginWithOTP handles OTP-based login (send OTP step)
func (s *AuthService) LoginWithOTP(ctx context.Context, req *LoginOTPRequest) (*LoginOTPResponse, error) {
	// Get user
	user, err := s.userRepo.GetByEmailOrPhone(ctx, req.Identifier)
	if err != nil {
		s.logAudit(ctx, nil, domain.AuditActionLoginFailed, false, req.IPAddress, "User not found")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked
	if user.IsLocked() {
		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "Account locked")
		return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	// Check if account can login
	if !user.CanLogin() {
		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "Account cannot login")
		return nil, fmt.Errorf("account is not active")
	}

	// Determine delivery method
	deliveryMethod := req.DeliveryMethod
	if deliveryMethod == "" {
		if user.Email != nil {
			deliveryMethod = domain.DeliveryMethodEmail
		} else {
			deliveryMethod = domain.DeliveryMethodSMS
		}
	}

	// Send OTP
	otpResp, err := s.otpService.SendOTP(ctx, &SendOTPRequest{
		Identifier:     req.Identifier,
		UserID:         &user.ID,
		DeliveryMethod: deliveryMethod,
		Purpose:        domain.OTPPurposeLogin,
		DeviceInfo:     req.DeviceInfo,
		IPAddress:      req.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send OTP: %w", err)
	}

	return &LoginOTPResponse{
		SentTo:    otpResp.DeliveredTo,
		ExpiresIn: otpResp.ExpiresIn,
	}, nil
}

// VerifyOTPAndLogin verifies OTP and completes login
func (s *AuthService) VerifyOTPAndLogin(ctx context.Context, req *VerifyOTPLoginRequest) (*TokenResponse, error) {
	// Verify OTP
	otpResp, err := s.otpService.VerifyOTP(ctx, &VerifyOTPRequest{
		Identifier: req.Identifier,
		Code:       req.Code,
		IPAddress:  req.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("OTP verification failed: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, *otpResp.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// If user was pending, activate them
	if user.Status == domain.UserStatusPending {
		user.Status = domain.UserStatusActive
		if otpResp.Email != nil {
			user.EmailVerified = true
		}
		if otpResp.Phone != nil {
			user.PhoneVerified = true
		}
		s.userRepo.Update(ctx, user)
	}

	// Reset failed attempts and update last login
	s.userRepo.ResetFailedAttempts(ctx, user.ID)
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Get user organizations
	userOrgs, err := s.userRepo.GetUserOrganizations(ctx, user.ID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get user organizations for user %s: %v\n", user.ID, err)
	}
	var primaryOrg *domain.UserOrganization
	if len(userOrgs) > 0 {
		primaryOrg = &userOrgs[0]
		fmt.Printf("[DEBUG] Found primary org for user %s: org_id=%s, role=%s\n", user.ID, primaryOrg.OrganizationID, primaryOrg.Role)
	} else {
		fmt.Printf("[WARN] No organizations found for user %s\n", user.ID)
	}

	// Generate tokens
	tokenReq := &TokenRequest{
		UserID:     user.ID,
		Email:      *user.Email,
		Name:       user.FullName,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  req.IPAddress,
	}
	if primaryOrg != nil {
		tokenReq.OrganizationID = primaryOrg.OrganizationID.String()
		tokenReq.Role = primaryOrg.Role
		tokenReq.Permissions = primaryOrg.Permissions
		fmt.Printf("[DEBUG] Token request includes: org_id=%s, role=%s\n", tokenReq.OrganizationID, tokenReq.Role)
	}

	tokens, err := s.jwtService.GenerateTokenPair(ctx, tokenReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log successful login
	s.logAudit(ctx, &user.ID, domain.AuditActionLoginSuccess, true, req.IPAddress, "")

	return tokens, nil
}

// LoginWithPassword handles password-based login
func (s *AuthService) LoginWithPassword(ctx context.Context, req *LoginPasswordRequest) (*TokenResponse, error) {
	fmt.Printf("[DEBUG] LoginWithPassword called, identifier=%s\n", req.Identifier)
	
	// Get user
	user, err := s.userRepo.GetByEmailOrPhone(ctx, req.Identifier)
	if err != nil {
		fmt.Printf("[DEBUG] User not found: %v\n", err)
		s.logAudit(ctx, nil, domain.AuditActionLoginFailed, false, req.IPAddress, "User not found")
		return nil, fmt.Errorf("invalid credentials")
	}

	fmt.Printf("[DEBUG] User found: id=%s, email=%v, status=%s\n", user.ID, user.Email, user.Status)

	// Check if account is locked
	if user.IsLocked() {
		fmt.Printf("[DEBUG] Account is locked\n")
		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "Account locked")
		return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	// Check if account can login
	if !user.CanLogin() {
		fmt.Printf("[DEBUG] Account cannot login, status=%s\n", user.Status)
		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "Account inactive")
		return nil, fmt.Errorf("account is not active")
	}

	// Check if user has password set
	if user.PasswordHash == nil {
		fmt.Printf("[DEBUG] No password set\n")
		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "No password set")
		return nil, fmt.Errorf("password login not available, please use OTP")
	}

	fmt.Printf("[DEBUG] Verifying password, hash_length=%d\n", len(*user.PasswordHash))

	// Verify password
	if err := s.pwdService.VerifyPassword(req.Password, *user.PasswordHash); err != nil {
		fmt.Printf("[DEBUG] Password verification failed: %v\n", err)
		// Increment failed attempts
		s.userRepo.IncrementFailedAttempts(ctx, user.ID)

		// Check if should lock account
		if user.FailedLoginAttempts+1 >= s.config.MaxFailedAttempts {
			lockUntil := time.Now().Add(s.config.LockoutDuration)
			s.userRepo.LockAccount(ctx, user.ID, lockUntil)
			s.logAudit(ctx, &user.ID, domain.AuditActionAccountLocked, true, req.IPAddress, "Too many failed attempts")
			return nil, fmt.Errorf("account locked due to too many failed attempts")
		}

		s.logAudit(ctx, &user.ID, domain.AuditActionLoginFailed, false, req.IPAddress, "Invalid password")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Reset failed attempts and update last login
	s.userRepo.ResetFailedAttempts(ctx, user.ID)
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Get user organizations
	userOrgs, err := s.userRepo.GetUserOrganizations(ctx, user.ID)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get user organizations for user %s: %v\n", user.ID, err)
	}
	var primaryOrg *domain.UserOrganization
	if len(userOrgs) > 0 {
		primaryOrg = &userOrgs[0]
		fmt.Printf("[DEBUG] Found primary org for user %s: org_id=%s, role=%s\n", user.ID, primaryOrg.OrganizationID, primaryOrg.Role)
	} else {
		fmt.Printf("[WARN] No organizations found for user %s\n", user.ID)
	}

	// Fetch organization details to get org_type
	var orgType string
	if primaryOrg != nil {
		org, err := s.orgRepo.GetByID(ctx, primaryOrg.OrganizationID)
		if err == nil {
			orgType = org.Type
			fmt.Printf("[DEBUG] Fetched organization: name=%s, type=%s\n", org.Name, org.Type)
		} else {
			fmt.Printf("[ERROR] Failed to fetch organization: %v\n", err)
		}
	}

	// Generate tokens
	tokenReq := &TokenRequest{
		UserID:     user.ID,
		Email:      *user.Email,
		Name:       user.FullName,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  req.IPAddress,
	}
	if primaryOrg != nil {
		tokenReq.OrganizationID = primaryOrg.OrganizationID.String()
		tokenReq.OrganizationType = orgType
		tokenReq.Role = primaryOrg.Role
		tokenReq.Permissions = primaryOrg.Permissions
		fmt.Printf("[DEBUG] Token request includes: org_id=%s, org_type=%s, role=%s\n", tokenReq.OrganizationID, tokenReq.OrganizationType, tokenReq.Role)
	} else {
		fmt.Printf("[WARN] No primary org, token will not include organization data\n")
	}

	tokens, err := s.jwtService.GenerateTokenPair(ctx, tokenReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log successful login
	s.logAudit(ctx, &user.ID, domain.AuditActionLoginSuccess, true, req.IPAddress, "")

	return tokens, nil
}

// Logout revokes refresh token
func (s *AuthService) Logout(ctx context.Context, req *LogoutRequest) error {
	if err := s.jwtService.RevokeToken(ctx, req.RefreshToken); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	s.logAudit(ctx, &req.UserID, domain.AuditActionLogout, true, req.IPAddress, "")
	return nil
}

// RefreshToken refreshes access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	tokens, err := s.jwtService.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	return tokens, nil
}

// ForgotPassword initiates password reset flow
func (s *AuthService) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	// Get user
	user, err := s.userRepo.GetByEmailOrPhone(ctx, req.Identifier)
	if err != nil {
		// Don't reveal if user exists or not (security)
		return &ForgotPasswordResponse{Message: "If the account exists, a reset code will be sent"}, nil
	}

	// Determine delivery method
	deliveryMethod := domain.DeliveryMethodEmail
	if user.Phone != nil {
		deliveryMethod = domain.DeliveryMethodSMS
	}

	// Send OTP
	otpResp, err := s.otpService.SendOTP(ctx, &SendOTPRequest{
		Identifier:     req.Identifier,
		UserID:         &user.ID,
		DeliveryMethod: deliveryMethod,
		Purpose:        domain.OTPPurposeReset,
		DeviceInfo:     req.DeviceInfo,
		IPAddress:      req.IPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send reset OTP: %w", err)
	}

	s.logAudit(ctx, &user.ID, domain.AuditActionPasswordReset, true, req.IPAddress, "")

	return &ForgotPasswordResponse{
		Message:   "Reset code sent",
		SentTo:    otpResp.DeliveredTo,
		ExpiresIn: otpResp.ExpiresIn,
	}, nil
}

// ResetPassword resets password with OTP verification
func (s *AuthService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	// Verify OTP
	otpResp, err := s.otpService.VerifyOTP(ctx, &VerifyOTPRequest{
		Identifier: req.Identifier,
		Code:       req.Code,
		IPAddress:  req.IPAddress,
	})
	if err != nil {
		return fmt.Errorf("invalid or expired reset code: %w", err)
	}

	// Validate new password
	if err := s.pwdService.ValidatePasswordStrength(req.NewPassword); err != nil {
		return fmt.Errorf("weak password: %w", err)
	}

	// Hash new password
	hash, err := s.pwdService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, *otpResp.UserID, hash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Reset failed attempts and unlock account
	s.userRepo.ResetFailedAttempts(ctx, *otpResp.UserID)
	s.userRepo.UnlockAccount(ctx, *otpResp.UserID)

	// Revoke all existing tokens (security)
	s.jwtService.RevokeAllUserTokens(ctx, *otpResp.UserID)

	s.logAudit(ctx, otpResp.UserID, domain.AuditActionPasswordChange, true, req.IPAddress, "")

	return nil
}

// logAudit logs an audit entry
func (s *AuthService) logAudit(ctx context.Context, userID *uuid.UUID, action string, success bool, ipAddress, errorMsg string) {
	log := &domain.AuthAuditLog{
		UserID:    userID,
		Action:    action,
		Success:   success,
		IPAddress: &ipAddress,
	}
	if errorMsg != "" {
		log.ErrorMessage = &errorMsg
	}
	s.auditRepo.Log(ctx, log)
}

// Helper function
func isEmail(identifier string) bool {
	return len(identifier) > 0 && identifier[0] != '+'
}

// Request/Response types

type RegisterRequest struct {
	Identifier          string
	Password            string
	FullName            string
	PreferredAuthMethod string
	DeviceInfo          domain.JSONBMap
	IPAddress           string
}

type RegisterResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	RequiresOTP bool      `json:"requires_otp"`
	OTPSentTo   string    `json:"otp_sent_to"`
	ExpiresIn   int       `json:"expires_in"`
}

type LoginOTPRequest struct {
	Identifier     string
	DeliveryMethod string
	DeviceInfo     domain.JSONBMap
	IPAddress      string
}

type LoginOTPResponse struct {
	SentTo    string `json:"sent_to"`
	ExpiresIn int    `json:"expires_in"`
}

type VerifyOTPLoginRequest struct {
	Identifier string
	Code       string
	DeviceInfo domain.JSONBMap
	IPAddress  string
}

type LoginPasswordRequest struct {
	Identifier string
	Password   string
	DeviceInfo domain.JSONBMap
	IPAddress  string
}

type LogoutRequest struct {
	UserID       uuid.UUID
	RefreshToken string
	IPAddress    string
}

type ForgotPasswordRequest struct {
	Identifier string
	DeviceInfo domain.JSONBMap
	IPAddress  string
}

type ForgotPasswordResponse struct {
	Message   string `json:"message"`
	SentTo    string `json:"sent_to,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

type ResetPasswordRequest struct {
	Identifier  string
	Code        string
	NewPassword string
	IPAddress   string
}
