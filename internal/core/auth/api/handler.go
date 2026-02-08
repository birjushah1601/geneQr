package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aby-med/medical-platform/internal/core/auth/app"
	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *app.AuthService
	jwtService  *app.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *app.AuthService, jwtService *app.JWTService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier          string                 `json:"identifier"`
		Password            string                 `json:"password"`
		FullName            string                 `json:"full_name"`
		PreferredAuthMethod string                 `json:"preferred_auth_method"`
		DeviceInfo          map[string]interface{} `json:"device_info"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Identifier == "" || req.FullName == "" {
		respondError(w, http.StatusBadRequest, "identifier and full_name are required")
		return
	}

	// Default preferred auth method
	if req.PreferredAuthMethod == "" {
		req.PreferredAuthMethod = domain.AuthMethodOTP
	}

	resp, err := h.authService.Register(r.Context(), &app.RegisterRequest{
		Identifier:          req.Identifier,
		Password:            req.Password,
		FullName:            req.FullName,
		PreferredAuthMethod: req.PreferredAuthMethod,
		DeviceInfo:          req.DeviceInfo,
		IPAddress:           getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// SendOTP sends OTP for login
// POST /api/v1/auth/send-otp
func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier     string                 `json:"identifier"`
		DeliveryMethod string                 `json:"delivery_method"`
		DeviceInfo     map[string]interface{} `json:"device_info"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Identifier == "" {
		respondError(w, http.StatusBadRequest, "identifier is required")
		return
	}

	resp, err := h.authService.LoginWithOTP(r.Context(), &app.LoginOTPRequest{
		Identifier:     req.Identifier,
		DeliveryMethod: req.DeliveryMethod,
		DeviceInfo:     req.DeviceInfo,
		IPAddress:      getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// VerifyOTP verifies OTP and logs in
// POST /api/v1/auth/verify-otp
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string                 `json:"identifier"`
		Code       string                 `json:"code"`
		DeviceInfo map[string]interface{} `json:"device_info"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Identifier == "" || req.Code == "" {
		respondError(w, http.StatusBadRequest, "identifier and code are required")
		return
	}

	tokens, err := h.authService.VerifyOTPAndLogin(r.Context(), &app.VerifyOTPLoginRequest{
		Identifier: req.Identifier,
		Code:       req.Code,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tokens)
}

// LoginWithPassword handles password login
// POST /api/v1/auth/login-password
func (h *AuthHandler) LoginWithPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string                 `json:"identifier"`
		Password   string                 `json:"password"`
		DeviceInfo map[string]interface{} `json:"device_info"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("[DEBUG] Failed to decode login request: %v\n", err)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	fmt.Printf("[DEBUG] Login attempt: identifier=%s, has_password=%v\n", req.Identifier, len(req.Password) > 0)

	if req.Identifier == "" || req.Password == "" {
		fmt.Printf("[DEBUG] Missing identifier or password\n")
		respondError(w, http.StatusBadRequest, "identifier and password are required")
		return
	}

	tokens, err := h.authService.LoginWithPassword(r.Context(), &app.LoginPasswordRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  getIPAddress(r),
	})
	if err != nil {
		fmt.Printf("[DEBUG] Login failed: identifier=%s, error=%v\n", req.Identifier, err)
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	fmt.Printf("[DEBUG] Login successful: identifier=%s\n", req.Identifier)
	respondJSON(w, http.StatusOK, tokens)
}

// RefreshToken refreshes access token
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tokens)
}

// Logout logs out user
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id")
	if userID == nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.authService.Logout(r.Context(), &app.LogoutRequest{
		UserID:       userUUID,
		RefreshToken: req.RefreshToken,
		IPAddress:    getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GetCurrentUser returns current authenticated user
// GET /api/v1/auth/me
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := r.Context().Value("claims")
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Claims is now a map[string]interface{} (set by AuthMiddleware)
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		respondError(w, http.StatusInternalServerError, "Invalid claims format")
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":           claimsMap["user_id"],
		"email":             claimsMap["email"],
		"name":              claimsMap["name"],
		"organization_id":   claimsMap["organization_id"],
		"organization_name": claimsMap["organization_name"],
		"organization_type": claimsMap["organization_type"],
		"role":              claimsMap["role"],
		"permissions":       claimsMap["permissions"],
	})
}

// ForgotPassword initiates password reset
// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string                 `json:"identifier"`
		DeviceInfo map[string]interface{} `json:"device_info"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Identifier == "" {
		respondError(w, http.StatusBadRequest, "identifier is required")
		return
	}

	resp, err := h.authService.ForgotPassword(r.Context(), &app.ForgotPasswordRequest{
		Identifier: req.Identifier,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ResetPassword resets password with OTP
// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier  string `json:"identifier"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Identifier == "" || req.Code == "" || req.NewPassword == "" {
		respondError(w, http.StatusBadRequest, "identifier, code, and new_password are required")
		return
	}

	err := h.authService.ResetPassword(r.Context(), &app.ResetPasswordRequest{
		Identifier:  req.Identifier,
		Code:        req.Code,
		NewPassword: req.NewPassword,
		IPAddress:   getIPAddress(r),
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// ValidateToken validates a JWT token
// POST /api/v1/auth/validate
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondError(w, http.StatusUnauthorized, "Authorization header required")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		respondError(w, http.StatusUnauthorized, "Bearer token required")
		return
	}

	claims, err := h.jwtService.ValidateToken(tokenString)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
		"email":   claims.Email,
		"exp":     claims.ExpiresAt.Time,
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}

func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take first IP if multiple
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to remote address
	return r.RemoteAddr
}

// RegisterRoutes registers all auth routes
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Public routes (no auth required)
		r.Post("/register", h.Register)
		r.Post("/send-otp", h.SendOTP)
		r.Post("/verify-otp", h.VerifyOTP)
		r.Post("/login-password", h.LoginWithPassword)
		r.Post("/refresh", h.RefreshToken)
		r.Post("/forgot-password", h.ForgotPassword)
		r.Post("/reset-password", h.ResetPassword)
		r.Post("/validate", h.ValidateToken)
		
		// Password setup via secure token
		r.Get("/validate-reset-token", h.ValidateResetToken)
		r.Post("/set-password", h.SetPassword)
		r.Post("/resend-setup-link", h.ResendSetupLink)

		// Protected routes (auth middleware required)
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Get("/me", h.GetCurrentUser)
			r.Post("/logout", h.Logout)
		})
	})
}

// ValidateResetToken validates a password reset token
// GET /api/v1/auth/validate-reset-token?token=xyz
func (h *AuthHandler) ValidateResetToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "token parameter is required")
		return
	}

	// Query database for token
	var userID uuid.UUID
	var email string
	var fullName string
	var expiresAt string
	var usedAt *string

	query := `
		SELECT u.id, u.email, u.full_name, t.expires_at, t.used_at
		FROM password_reset_tokens t
		JOIN users u ON t.user_id = u.id
		WHERE t.token = $1
	`

	err := h.authService.GetDB().QueryRow(r.Context(), query, token).Scan(&userID, &email, &fullName, &expiresAt, &usedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "Invalid or expired token")
		return
	}

	// Check if already used
	if usedAt != nil {
		respondError(w, http.StatusBadRequest, "Token has already been used")
		return
	}

	// Check if expired (compare with NOW())
	var isExpired bool
	expireCheck := `SELECT NOW() > $1::timestamp`
	err = h.authService.GetDB().QueryRow(r.Context(), expireCheck, expiresAt).Scan(&isExpired)
	if err != nil || isExpired {
		respondError(w, http.StatusBadRequest, "Token has expired")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":     true,
		"user_id":   userID,
		"email":     email,
		"full_name": fullName,
	})
}

// SetPassword sets password using a valid reset token
// POST /api/v1/auth/set-password
func (h *AuthHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.NewPassword == "" {
		respondError(w, http.StatusBadRequest, "new_password is required")
		return
	}

	// Validate token and get user
	var userID uuid.UUID
	var usedAt *string
	var expiresAt string

	query := `
		SELECT user_id, used_at, expires_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	err := h.authService.GetDB().QueryRow(r.Context(), query, req.Token).Scan(&userID, &usedAt, &expiresAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "Invalid or expired token")
		return
	}

	if usedAt != nil {
		respondError(w, http.StatusBadRequest, "Token has already been used")
		return
	}

	// Check if expired
	var isExpired bool
	expireCheck := `SELECT NOW() > $1::timestamp`
	err = h.authService.GetDB().QueryRow(r.Context(), expireCheck, expiresAt).Scan(&isExpired)
	if err != nil || isExpired {
		respondError(w, http.StatusBadRequest, "Token has expired")
		return
	}

	// Update user password and status
	err = h.authService.SetPasswordByUserID(r.Context(), userID, req.NewPassword)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Mark token as used
	updateToken := `UPDATE password_reset_tokens SET used_at = NOW() WHERE token = $1`
	_, err = h.authService.GetDB().Exec(r.Context(), updateToken, req.Token)
	if err != nil {
		// Log but don't fail - password is already set
		fmt.Printf("Warning: Failed to mark token as used: %v\n", err)
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Password set successfully. You can now login.",
	})
}

// ResendSetupLink resends a password setup link for expired tokens
// POST /api/v1/auth/resend-setup-link
func (h *AuthHandler) ResendSetupLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	// Get user by email
	var userID uuid.UUID
	var status string
	var hasPassword bool

	userQuery := `SELECT id, status, (password_hash IS NOT NULL) as has_password FROM users WHERE email = $1`
	err := h.authService.GetDB().QueryRow(r.Context(), userQuery, req.Email).Scan(&userID, &status, &hasPassword)
	if err != nil {
		// Don't reveal if user exists (security)
		respondJSON(w, http.StatusOK, map[string]string{
			"message": "If the account exists and is pending setup, a new link will be sent.",
		})
		return
	}

	// Only send link if user hasn't set password yet
	if hasPassword {
		respondError(w, http.StatusBadRequest, "Account is already active. Use forgot password instead.")
		return
	}

	// Generate new token
	token := generateSecureToken()

	// Insert new token (old ones remain but will be ignored)
	tokenQuery := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at) 
		VALUES (gen_random_uuid(), $1, $2, NOW() + INTERVAL '48 hours', NOW())
	`
	_, err = h.authService.GetDB().Exec(r.Context(), tokenQuery, userID, token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate new token")
		return
	}

	// Generate setup link
	setupLink := fmt.Sprintf("http://localhost:3000/set-password?token=%s", token)

	// TODO: Send email here
	fmt.Printf("ðŸ“§ New setup link for %s: %s\n", req.Email, setupLink)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "A new setup link has been sent to your email.",
		"setup_link": setupLink, // Remove in production
	})
}

// AuthMiddleware validates JWT token and adds claims to context
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public endpoints that don't require authentication
		// Note: /api/v1/auth/me is NOT public - it requires auth to get user info
		publicPaths := map[string]bool{
			"/health":                             true,
			"/metrics":                            true,
			"/api/v1/auth/login-password":         true,
			"/api/v1/auth/register":               true,
			"/api/v1/auth/send-otp":               true,
			"/api/v1/auth/verify-otp":             true,
			"/api/v1/auth/refresh":                true,
			"/api/v1/auth/forgot-password":        true,
			"/api/v1/auth/reset-password":         true,
			"/api/v1/auth/validate":               true,
			"/api/v1/auth/validate-reset-token":   true,
			"/api/v1/auth/set-password":           true,
			"/api/v1/auth/resend-setup-link":      true,
		}
		
		// Skip auth for public endpoints
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow QR code lookups (for public service request creation)
		if strings.HasPrefix(r.URL.Path, "/api/v1/equipment/qr/") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow public tracking pages (customer ticket status)
		if strings.HasPrefix(r.URL.Path, "/api/v1/track/") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow public timeline for tracking pages (customer ticket status)
		if strings.HasPrefix(r.URL.Path, "/api/v1/tickets/") && strings.Contains(r.URL.Path, "/timeline") {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow public ticket creation (for QR-based service requests)
		if r.URL.Path == "/api/v1/tickets" && r.Method == "POST" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow attachments operations for public ticket tracking
		// POST: upload files, GET: list/view attachments, GET with /download: download files
		if strings.HasPrefix(r.URL.Path, "/api/v1/attachments") {
			// Allow POST (upload), GET (list/view), and download
			if r.Method == "POST" || r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}
		}
		
		// Allow invitation routes (they use tokens, not JWT) - check prefix
		if strings.HasPrefix(r.URL.Path, "/api/v1/invitations/validate/") || 
		   (strings.HasPrefix(r.URL.Path, "/api/v1/invitations/") && strings.HasSuffix(r.URL.Path, "/accept")) {
			next.ServeHTTP(w, r)
			return
		}
		
		// Allow public ticket tracking (no auth required - uses secure tokens)
		if strings.HasPrefix(r.URL.Path, "/api/v1/track/") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			respondError(w, http.StatusUnauthorized, "Bearer token required")
			return
		}

		claims, err := h.jwtService.ValidateToken(tokenString)
		if err != nil {
			fmt.Printf("âŒ Token validation failed for path %s: %v\n", r.URL.Path, err)
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		fmt.Printf("âœ… Token validated for path %s - user_id=%s, org_id=%s, org_type=%s\n", 
			r.URL.Path, claims.UserID, claims.OrganizationID, claims.OrganizationType)

		// Add claims to context for downstream middleware
		// Convert Claims struct to map for OrganizationContextMiddleware
		ctx := r.Context()
		claimsMap := map[string]interface{}{
			"user_id":           claims.UserID,
			"email":             claims.Email,
			"name":              claims.Name,
			"organization_id":   claims.OrganizationID,
			"organization_name": claims.OrganizationName,
			"organization_type": claims.OrganizationType,
			"role":              claims.Role,
			"permissions":       claims.Permissions,
		}
		
		ctx = context.WithValue(ctx, "claims", claimsMap)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() string {
	b := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(b); err != nil {
		// Fallback (should never happen)
		return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 1234567890)))
	}
	return base64.URLEncoding.EncodeToString(b)
}
