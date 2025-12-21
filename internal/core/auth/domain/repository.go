package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// User operations
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	GetByEmailOrPhone(ctx context.Context, identifier string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error
	UnlockAccount(ctx context.Context, userID uuid.UUID) error

	// User organization operations
	GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]UserOrganization, error)
	AddUserToOrganization(ctx context.Context, userOrg *UserOrganization) error
	RemoveUserFromOrganization(ctx context.Context, userID, orgID uuid.UUID) error
	UpdateUserRole(ctx context.Context, userID, orgID uuid.UUID, role string, permissions []string) error
}

// OTPRepository defines the interface for OTP data access
type OTPRepository interface {
	Create(ctx context.Context, otp *OTPCode) error
	GetByCode(ctx context.Context, identifier, code string) (*OTPCode, error)
	GetLatest(ctx context.Context, identifier, purpose string) (*OTPCode, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) (int, error)
	CountRecentOTPs(ctx context.Context, identifier string, since time.Time) (int, error)
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]RefreshToken, error)
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID, reason string) error
	DeleteExpired(ctx context.Context) (int, error)
}

// AuditRepository defines the interface for audit log data access
type AuditRepository interface {
	Log(ctx context.Context, log *AuthAuditLog) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]AuthAuditLog, error)
	GetFailedLoginsByIP(ctx context.Context, ipAddress string, since time.Time) (int, error)
	GetRecentActivity(ctx context.Context, userID uuid.UUID, limit int) ([]AuthAuditLog, error)
}

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	GetByOrgType(ctx context.Context, orgType string) ([]Role, error)
	GetAll(ctx context.Context) ([]Role, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// NotificationPreferencesRepository defines the interface for notification preferences data access
type NotificationPreferencesRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error)
	Create(ctx context.Context, prefs *NotificationPreferences) error
	Update(ctx context.Context, prefs *NotificationPreferences) error
}
