package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
)

// OTPService handles OTP generation, validation, and delivery
type OTPService struct {
	otpRepo     domain.OTPRepository
	auditRepo   domain.AuditRepository
	emailSender EmailSender
	smsSender   SMSSender
	config      *OTPConfig
}

// OTPConfig holds OTP service configuration
type OTPConfig struct {
	Length           int
	ExpiryMinutes    int
	MaxAttempts      int
	RateLimitPerHour int
	CooldownSeconds  int
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendOTP(ctx context.Context, to, otp string) error
}

// SMSSender defines the interface for sending SMS
type SMSSender interface {
	SendOTP(ctx context.Context, to, otp string) error
	SendWhatsAppOTP(ctx context.Context, to, otp string) error
}

// NewOTPService creates a new OTP service
func NewOTPService(
	otpRepo domain.OTPRepository,
	auditRepo domain.AuditRepository,
	emailSender EmailSender,
	smsSender SMSSender,
	config *OTPConfig,
) *OTPService {
	if config == nil {
		config = &OTPConfig{
			Length:           6,
			ExpiryMinutes:    5,
			MaxAttempts:      3,
			RateLimitPerHour: 3,
			CooldownSeconds:  60,
		}
	}

	return &OTPService{
		otpRepo:     otpRepo,
		auditRepo:   auditRepo,
		emailSender: emailSender,
		smsSender:   smsSender,
		config:      config,
	}
}

// SendOTP generates and sends an OTP code
func (s *OTPService) SendOTP(ctx context.Context, req *SendOTPRequest) (*SendOTPResponse, error) {
	// Check rate limiting
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	recentCount, err := s.otpRepo.CountRecentOTPs(ctx, req.Identifier, oneHourAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	if recentCount >= s.config.RateLimitPerHour {
		s.logAudit(ctx, req.UserID, domain.AuditActionOTPSent, false, req.IPAddress, "Rate limit exceeded")
		return nil, fmt.Errorf("rate limit exceeded: too many OTP requests")
	}

	// Check cooldown period
	latest, err := s.otpRepo.GetLatest(ctx, req.Identifier, req.Purpose)
	if err == nil && latest != nil {
		cooldownUntil := latest.CreatedAt.Add(time.Duration(s.config.CooldownSeconds) * time.Second)
		if time.Now().Before(cooldownUntil) {
			remaining := int(time.Until(cooldownUntil).Seconds())
			return nil, fmt.Errorf("please wait %d seconds before requesting another OTP", remaining)
		}
	}

	// Generate OTP code
	code, err := s.generateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Hash the code for storage
	codeHash := s.hashOTPCode(code)

	// Create OTP record
	otp := &domain.OTPCode{
		ID:             uuid.New(),
		UserID:         req.UserID,
		Code:           code,
		CodeHash:       codeHash,
		DeliveryMethod: req.DeliveryMethod,
		Purpose:        req.Purpose,
		ExpiresAt:      time.Now().Add(time.Duration(s.config.ExpiryMinutes) * time.Minute),
		DeviceInfo:     req.DeviceInfo,
		IPAddress:      &req.IPAddress,
	}

	// Set email or phone based on identifier
	if req.DeliveryMethod == domain.DeliveryMethodEmail {
		otp.Email = &req.Identifier
	} else {
		otp.Phone = &req.Identifier
	}

	// Save OTP to database
	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP via appropriate channel
	if err := s.deliverOTP(ctx, req.Identifier, code, req.DeliveryMethod); err != nil {
		s.logAudit(ctx, req.UserID, domain.AuditActionOTPSent, false, req.IPAddress, err.Error())
		return nil, fmt.Errorf("failed to deliver OTP: %w", err)
	}

	// Log successful OTP send
	s.logAudit(ctx, req.UserID, domain.AuditActionOTPSent, true, req.IPAddress, "")

	return &SendOTPResponse{
		ExpiresIn:   s.config.ExpiryMinutes * 60,
		RetryAfter:  s.config.CooldownSeconds,
		DeliveredTo: s.maskIdentifier(req.Identifier, req.DeliveryMethod),
	}, nil
}

// VerifyOTP verifies an OTP code
func (s *OTPService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*VerifyOTPResponse, error) {
	// Hash the provided code
	codeHash := s.hashOTPCode(req.Code)

	// Get OTP from database
	otp, err := s.otpRepo.GetByCode(ctx, req.Identifier, codeHash)
	if err != nil {
		s.logAudit(ctx, nil, domain.AuditActionOTPFailed, false, req.IPAddress, "OTP not found")
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	// Check if already used
	if otp.Used {
		s.logAudit(ctx, otp.UserID, domain.AuditActionOTPFailed, false, req.IPAddress, "OTP already used")
		return nil, fmt.Errorf("OTP has already been used")
	}

	// Check attempts
	if otp.Attempts >= s.config.MaxAttempts {
		s.logAudit(ctx, otp.UserID, domain.AuditActionOTPFailed, false, req.IPAddress, "Max attempts exceeded")
		return nil, fmt.Errorf("maximum verification attempts exceeded")
	}

	// Check expiry
	if otp.IsExpired() {
		s.logAudit(ctx, otp.UserID, domain.AuditActionOTPFailed, false, req.IPAddress, "OTP expired")
		return nil, fmt.Errorf("OTP has expired")
	}

	// Verify the code matches
	if otp.CodeHash != codeHash {
		// Increment attempts
		s.otpRepo.IncrementAttempts(ctx, otp.ID)
		s.logAudit(ctx, otp.UserID, domain.AuditActionOTPFailed, false, req.IPAddress, "Invalid code")
		return nil, fmt.Errorf("invalid OTP code")
	}

	// Mark as used
	if err := s.otpRepo.MarkAsUsed(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	// Log successful verification
	s.logAudit(ctx, otp.UserID, domain.AuditActionOTPVerified, true, req.IPAddress, "")

	return &VerifyOTPResponse{
		UserID: otp.UserID,
		Email:  otp.Email,
		Phone:  otp.Phone,
	}, nil
}

// generateOTPCode generates a cryptographically secure OTP code
func (s *OTPService) generateOTPCode() (string, error) {
	// Generate random number with specified length
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(int64(s.config.Length)), nil).Sub(max, big.NewInt(1))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", s.config.Length)
	return fmt.Sprintf(format, n), nil
}

// hashOTPCode hashes an OTP code using SHA-256
func (s *OTPService) hashOTPCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

// deliverOTP sends the OTP via the specified delivery method
func (s *OTPService) deliverOTP(ctx context.Context, identifier, code, method string) error {
	switch method {
	case domain.DeliveryMethodEmail:
		return s.emailSender.SendOTP(ctx, identifier, code)
	case domain.DeliveryMethodSMS:
		return s.smsSender.SendOTP(ctx, identifier, code)
	case domain.DeliveryMethodWhatsApp:
		return s.smsSender.SendWhatsAppOTP(ctx, identifier, code)
	default:
		return fmt.Errorf("unsupported delivery method: %s", method)
	}
}

// maskIdentifier masks the identifier for privacy
func (s *OTPService) maskIdentifier(identifier, method string) string {
	if method == domain.DeliveryMethodEmail {
		// Mask email: a***@example.com
		at := 0
		for i, c := range identifier {
			if c == '@' {
				at = i
				break
			}
		}
		if at > 2 {
			return identifier[0:1] + "***" + identifier[at:]
		}
		return identifier
	}

	// Mask phone: +1***7890
	if len(identifier) > 6 {
		return identifier[0:2] + "***" + identifier[len(identifier)-4:]
	}
	return identifier
}

// logAudit logs an audit entry
func (s *OTPService) logAudit(ctx context.Context, userID *uuid.UUID, action string, success bool, ipAddress, errorMsg string) {
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

// Request/Response types

type SendOTPRequest struct {
	Identifier     string
	UserID         *uuid.UUID
	DeliveryMethod string
	Purpose        string
	DeviceInfo     domain.JSONBMap
	IPAddress      string
}

type SendOTPResponse struct {
	ExpiresIn   int    // seconds
	RetryAfter  int    // seconds
	DeliveredTo string // masked identifier
}

type VerifyOTPRequest struct {
	Identifier string
	Code       string
	IPAddress  string
}

type VerifyOTPResponse struct {
	UserID *uuid.UUID
	Email  *string
	Phone  *string
}
