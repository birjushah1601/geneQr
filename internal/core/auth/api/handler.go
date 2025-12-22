package api

import (
	"context"
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

	userClaims := claims.(*app.Claims)
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":         userClaims.UserID,
		"email":           userClaims.Email,
		"name":            userClaims.Name,
		"organization_id": userClaims.OrganizationID,
		"role":            userClaims.Role,
		"permissions":     userClaims.Permissions,
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

		// Protected routes (auth middleware required)
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Get("/me", h.GetCurrentUser)
			r.Post("/logout", h.Logout)
		})
	})
}

// AuthMiddleware validates JWT token and adds claims to context
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// Add claims to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
