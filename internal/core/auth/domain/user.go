package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system
type User struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Email                 *string    `json:"email" db:"email"`
	Phone                 *string    `json:"phone" db:"phone"`
	PasswordHash          *string    `json:"-" db:"password_hash"` // Never expose in JSON
	PreferredAuthMethod   string     `json:"preferred_auth_method" db:"preferred_auth_method"`
	EmailVerified         bool       `json:"email_verified" db:"email_verified"`
	PhoneVerified         bool       `json:"phone_verified" db:"phone_verified"`
	FullName              string     `json:"full_name" db:"full_name"`
	AvatarURL             *string    `json:"avatar_url" db:"avatar_url"`
	Status                string     `json:"status" db:"status"`
	FailedLoginAttempts   int        `json:"-" db:"failed_login_attempts"`
	LockedUntil           *time.Time `json:"locked_until,omitempty" db:"locked_until"`
	LastLogin             *time.Time `json:"last_login" db:"last_login"`
	LastOTPSent           *time.Time `json:"-" db:"last_otp_sent"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	Metadata              JSONBMap   `json:"metadata" db:"metadata"`
}

// UserOrganization represents a user's membership in an organization
type UserOrganization struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	Role           string     `json:"role" db:"role"`
	Permissions    []string   `json:"permissions" db:"permissions"`
	IsPrimary      bool       `json:"is_primary" db:"is_primary"`
	Status         string     `json:"status" db:"status"`
	JoinedAt       time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty" db:"left_at"`
}

// UserWithOrganizations represents a user with their organization memberships
type UserWithOrganizations struct {
	User          User                      `json:"user"`
	Organizations []OrganizationMembership  `json:"organizations"`
}

// OrganizationMembership represents an organization the user belongs to
type OrganizationMembership struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Type           string     `json:"type"`
	Role           string     `json:"role"`
	Permissions    []string   `json:"permissions"`
	IsPrimary      bool       `json:"is_primary"`
}

// OTPCode represents a one-time password
type OTPCode struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Email          *string    `json:"email,omitempty" db:"email"`
	Phone          *string    `json:"phone,omitempty" db:"phone"`
	Code           string     `json:"-" db:"code"` // Never expose in JSON
	CodeHash       string     `json:"-" db:"code_hash"`
	DeliveryMethod string     `json:"delivery_method" db:"delivery_method"`
	Purpose        string     `json:"purpose" db:"purpose"`
	Used           bool       `json:"used" db:"used"`
	Attempts       int        `json:"attempts" db:"attempts"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	DeviceInfo     JSONBMap   `json:"device_info,omitempty" db:"device_info"`
	IPAddress      *string    `json:"ip_address,omitempty" db:"ip_address"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UsedAt         *time.Time `json:"used_at,omitempty" db:"used_at"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash    string     `json:"-" db:"token_hash"`
	DeviceInfo   JSONBMap   `json:"device_info" db:"device_info"`
	IPAddress    *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string    `json:"user_agent,omitempty" db:"user_agent"`
	Revoked      bool       `json:"revoked" db:"revoked"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	RevokeReason *string    `json:"revoke_reason,omitempty" db:"revoke_reason"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	UsageCount   int        `json:"usage_count" db:"usage_count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// AuthAuditLog represents an audit log entry for authentication events
type AuthAuditLog struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Action       string     `json:"action" db:"action"`
	Success      bool       `json:"success" db:"success"`
	IPAddress    *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string    `json:"user_agent,omitempty" db:"user_agent"`
	DeviceInfo   JSONBMap   `json:"device_info,omitempty" db:"device_info"`
	Metadata     JSONBMap   `json:"metadata,omitempty" db:"metadata"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Role represents a user role with permissions
type Role struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	DisplayName string     `json:"display_name" db:"display_name"`
	Description *string    `json:"description,omitempty" db:"description"`
	OrgType     *string    `json:"org_type,omitempty" db:"org_type"`
	Permissions JSONBArray `json:"permissions" db:"permissions"`
	IsSystem    bool       `json:"is_system" db:"is_system"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// NotificationPreferences represents user's notification settings
type NotificationPreferences struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	UserID                uuid.UUID  `json:"user_id" db:"user_id"`
	EmailNotifications    bool       `json:"email_notifications" db:"email_notifications"`
	SMSNotifications      bool       `json:"sms_notifications" db:"sms_notifications"`
	WhatsAppNotifications bool       `json:"whatsapp_notifications" db:"whatsapp_notifications"`
	InAppNotifications    bool       `json:"in_app_notifications" db:"in_app_notifications"`
	Events                JSONBMap   `json:"events" db:"events"`
	QuietHours            JSONBMap   `json:"quiet_hours,omitempty" db:"quiet_hours"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// JSONBMap represents a JSONB map column
type JSONBMap map[string]interface{}

// Scan implements sql.Scanner interface for JSONBMap
func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBMap: %T", value)
	}
	
	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface for JSONBMap
func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONBArray represents a JSONB array column
type JSONBArray []string

// Scan implements sql.Scanner interface for JSONBArray
func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBArray: %T", value)
	}
	
	return json.Unmarshal(bytes, j)
}

// Value implements driver.Valuer interface for JSONBArray
func (j JSONBArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// UserStatus constants
const (
	UserStatusActive    = "active"
	UserStatusSuspended = "suspended"
	UserStatusDeleted   = "deleted"
	UserStatusPending   = "pending"
)

// AuthMethod constants
const (
	AuthMethodOTP      = "otp"
	AuthMethodPassword = "password"
)

// DeliveryMethod constants
const (
	DeliveryMethodEmail    = "email"
	DeliveryMethodSMS      = "sms"
	DeliveryMethodWhatsApp = "whatsapp"
)

// OTPPurpose constants
const (
	OTPPurposeLogin    = "login"
	OTPPurposeVerify   = "verify"
	OTPPurposeReset    = "reset"
	OTPPurposeRegister = "register"
)

// AuditAction constants
const (
	AuditActionLoginSuccess   = "login_success"
	AuditActionLoginFailed    = "login_failed"
	AuditActionOTPSent        = "otp_sent"
	AuditActionOTPVerified    = "otp_verified"
	AuditActionOTPFailed      = "otp_failed"
	AuditActionLogout         = "logout"
	AuditActionPasswordChange = "password_changed"
	AuditActionPasswordReset  = "password_reset"
	AuditActionAccountLocked  = "account_locked"
	AuditActionAccountUnlock  = "account_unlocked"
	AuditActionRegister       = "register"
	AuditActionTokenRefresh   = "token_refresh"
)

// Helper methods

// IsLocked checks if the user account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CanLogin checks if the user can login
func (u *User) CanLogin() bool {
	return u.IsActive() && !u.IsLocked()
}

// IsExpired checks if the OTP is expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// CanRetry checks if the OTP can be retried (max 3 attempts)
func (o *OTPCode) CanRetry() bool {
	return o.Attempts < 3
}

// IsExpired checks if the refresh token is expired
func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (r *RefreshToken) IsValid() bool {
	return !r.IsExpired() && !r.Revoked
}
