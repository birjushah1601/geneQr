package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// Organization context keys
type contextKey string

const (
	OrganizationIDKey      contextKey = "organization_id"
	OrganizationTypeKey    contextKey = "organization_type"
	UserRoleKey            contextKey = "user_role"
	UserPermissionsKey     contextKey = "user_permissions"
	UserIDKey              contextKey = "user_id"
	UserEmailKey           contextKey = "user_email"
)

// OrganizationContextMiddleware extracts organization info from JWT claims
// and injects it into the request context for downstream handlers
func OrganizationContextMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT claims from context (set by auth middleware)
			claims, ok := r.Context().Value("claims").(map[string]interface{})
			if !ok {
				// No claims - might be public endpoint or auth not required
				logger.Debug("No JWT claims found in request context", "path", r.URL.Path)
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			// Extract user_id
			if userIDStr, ok := claims["user_id"].(string); ok && userIDStr != "" {
				if userID, err := uuid.Parse(userIDStr); err == nil {
					ctx = context.WithValue(ctx, UserIDKey, userID)
				}
			}

			// Extract email
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, UserEmailKey, email)
			}

			// Extract organization_id
			if orgIDStr, ok := claims["organization_id"].(string); ok && orgIDStr != "" {
				if orgID, err := uuid.Parse(orgIDStr); err == nil {
					ctx = context.WithValue(ctx, OrganizationIDKey, orgID)
					logger.Debug("Organization context set",
						"organization_id", orgID,
						"path", r.URL.Path,
						"method", r.Method)
				} else {
					logger.Warn("Invalid organization_id in JWT",
						"organization_id", orgIDStr,
						"error", err)
				}
			} else {
				logger.Debug("No organization_id in JWT claims", "path", r.URL.Path)
			}

			// Extract organization_type
			if orgType, ok := claims["organization_type"].(string); ok {
				ctx = context.WithValue(ctx, OrganizationTypeKey, orgType)
				logger.Debug("Organization type set", "organization_type", orgType)
			}

			// Extract user role
			if role, ok := claims["role"].(string); ok {
				ctx = context.WithValue(ctx, UserRoleKey, role)
			}

			// Extract permissions (can be []interface{} or []string)
			if permsInterface, ok := claims["permissions"].([]interface{}); ok {
				permissions := make([]string, 0, len(permsInterface))
				for _, p := range permsInterface {
					if pStr, ok := p.(string); ok {
						permissions = append(permissions, pStr)
					}
				}
				ctx = context.WithValue(ctx, UserPermissionsKey, permissions)
			} else if permsStrings, ok := claims["permissions"].([]string); ok {
				ctx = context.WithValue(ctx, UserPermissionsKey, permsStrings)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetOrganizationID extracts organization ID from context
func GetOrganizationID(ctx context.Context) (uuid.UUID, bool) {
	orgID, ok := ctx.Value(OrganizationIDKey).(uuid.UUID)
	return orgID, ok
}

// GetOrganizationType extracts organization type from context
func GetOrganizationType(ctx context.Context) (string, bool) {
	orgType, ok := ctx.Value(OrganizationTypeKey).(string)
	return orgType, ok
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

// GetUserPermissions extracts user permissions from context
func GetUserPermissions(ctx context.Context) ([]string, bool) {
	perms, ok := ctx.Value(UserPermissionsKey).([]string)
	return perms, ok
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// RequireOrganizationContext middleware ensures organization context exists
// Use this for endpoints that MUST have an organization context
func RequireOrganizationContext(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgID, ok := GetOrganizationID(r.Context())
			if !ok || orgID == uuid.Nil {
				logger.Warn("Organization context required but not found",
					"path", r.URL.Path,
					"method", r.Method)
				
				http.Error(w, `{"error":{"message":"Organization context required"}}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
