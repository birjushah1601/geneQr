package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeys            map[string]APIKeyInfo `json:"api_keys"`
	RequireAuth        bool                  `json:"require_auth"`
	AuthHeaderName     string               `json:"auth_header_name"`
	BypassPaths        []string             `json:"bypass_paths"`
	EnableRateLimiting bool                 `json:"enable_rate_limiting"`
}

// APIKeyInfo contains metadata for an API key
type APIKeyInfo struct {
	Name         string            `json:"name"`
	Permissions  []string          `json:"permissions"`
	RateLimit    int              `json:"rate_limit"`    // requests per hour
	CreatedAt    time.Time        `json:"created_at"`
	LastUsedAt   *time.Time       `json:"last_used_at"`
	IsActive     bool             `json:"is_active"`
	Owner        string           `json:"owner"`
	Description  string           `json:"description"`
	Metadata     map[string]string `json:"metadata"`
}

// AuthContext holds authentication information for the request
type AuthContext struct {
	APIKey      string      `json:"api_key"`
	KeyInfo     APIKeyInfo  `json:"key_info"`
	Permissions []string    `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
}

// AuthMiddleware handles API key authentication
type AuthMiddleware struct {
	config *AuthConfig
	logger *slog.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(config *AuthConfig, logger *slog.Logger) *AuthMiddleware {
	if config.AuthHeaderName == "" {
		config.AuthHeaderName = "X-API-Key"
	}
	
	return &AuthMiddleware{
		config: config,
		logger: logger.With(slog.String("component", "auth_middleware")),
	}
}

// DefaultAuthConfig returns a default authentication configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		APIKeys: map[string]APIKeyInfo{
			"dev-key-001": {
				Name:        "Development Key",
				Permissions: []string{"read", "write", "upload", "analyze"},
				RateLimit:   1000, // 1000 requests per hour
				CreatedAt:   time.Now(),
				IsActive:    true,
				Owner:       "development",
				Description: "Development and testing key",
				Metadata: map[string]string{
					"environment": "development",
					"team":       "engineering",
				},
			},
			"admin-key-001": {
				Name:        "Admin Key",
				Permissions: []string{"read", "write", "upload", "analyze", "admin", "delete"},
				RateLimit:   5000, // 5000 requests per hour
				CreatedAt:   time.Now(),
				IsActive:    true,
				Owner:       "admin",
				Description: "Administrative access key",
				Metadata: map[string]string{
					"environment": "all",
					"team":       "operations",
				},
			},
		},
		RequireAuth:        true,
		AuthHeaderName:     "X-API-Key",
		BypassPaths:        []string{"/health", "/metrics", "/ping", "/docs"},
		EnableRateLimiting: true,
	}
}

// Middleware returns the HTTP middleware function
func (m *AuthMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if authentication is required
			if !m.config.RequireAuth {
				next.ServeHTTP(w, r)
				return
			}

			// Check if path should bypass authentication
			if m.shouldBypassAuth(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract API key from request
			apiKey := m.extractAPIKey(r)
			if apiKey == "" {
				m.logger.Warn("Missing API key", 
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("remote_addr", r.RemoteAddr))
				m.writeUnauthorized(w, "API key is required")
				return
			}

			// Validate API key
			keyInfo, valid := m.validateAPIKey(apiKey)
			if !valid {
				m.logger.Warn("Invalid API key", 
					slog.String("api_key", maskAPIKey(apiKey)),
					slog.String("path", r.URL.Path),
					slog.String("remote_addr", r.RemoteAddr))
				m.writeUnauthorized(w, "Invalid API key")
				return
			}

			// Update last used timestamp
			keyInfo.LastUsedAt = &time.Time{}
			*keyInfo.LastUsedAt = time.Now()

			// Create auth context
			authCtx := &AuthContext{
				APIKey:      apiKey,
				KeyInfo:     keyInfo,
				Permissions: keyInfo.Permissions,
				Metadata:    keyInfo.Metadata,
			}

			// Add auth context to request
			ctx := context.WithValue(r.Context(), "auth", authCtx)
			r = r.WithContext(ctx)

			// Log successful authentication
			m.logger.Debug("API request authenticated", 
				slog.String("api_key_name", keyInfo.Name),
				slog.String("owner", keyInfo.Owner),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method))

			next.ServeHTTP(w, r)
		})
	}
}

// extractAPIKey extracts the API key from the request
func (m *AuthMiddleware) extractAPIKey(r *http.Request) string {
	// Try header first
	if key := r.Header.Get(m.config.AuthHeaderName); key != "" {
		return key
	}

	// Try Authorization header with Bearer token
	if auth := r.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
	}

	// Try query parameter as fallback
	if key := r.URL.Query().Get("api_key"); key != "" {
		return key
	}

	return ""
}

// validateAPIKey validates an API key and returns key info
func (m *AuthMiddleware) validateAPIKey(apiKey string) (APIKeyInfo, bool) {
	keyInfo, exists := m.config.APIKeys[apiKey]
	if !exists {
		return APIKeyInfo{}, false
	}

	// Check if key is active
	if !keyInfo.IsActive {
		return APIKeyInfo{}, false
	}

	return keyInfo, true
}

// shouldBypassAuth checks if a path should bypass authentication
func (m *AuthMiddleware) shouldBypassAuth(path string) bool {
	for _, bypassPath := range m.config.BypassPaths {
		if strings.HasPrefix(path, bypassPath) {
			return true
		}
	}
	return false
}

// writeUnauthorized writes an unauthorized response
func (m *AuthMiddleware) writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"success": false, "error": "` + message + `"}`))
}

// maskAPIKey masks an API key for logging
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
}

// GetAuthContext extracts authentication context from request context
func GetAuthContext(r *http.Request) (*AuthContext, bool) {
	if auth := r.Context().Value("auth"); auth != nil {
		if authCtx, ok := auth.(*AuthContext); ok {
			return authCtx, true
		}
	}
	return nil, false
}

// HasPermission checks if the authenticated user has a specific permission
func HasPermission(r *http.Request, permission string) bool {
	authCtx, exists := GetAuthContext(r)
	if !exists {
		return false
	}

	for _, perm := range authCtx.Permissions {
		if perm == permission || perm == "admin" {
			return true
		}
	}
	return false
}

// RequirePermission creates a middleware that requires a specific permission
func RequirePermission(permission string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasPermission(r, permission) {
				logger.Warn("Insufficient permissions", 
					slog.String("required_permission", permission),
					slog.String("path", r.URL.Path))
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"success": false, "error": "Insufficient permissions"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}