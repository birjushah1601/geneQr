package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// GetDB returns the underlying database connection for direct queries
func (r *userRepository) GetDB() *sqlx.DB {
	return r.db
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, email, phone, password_hash, preferred_auth_method,
			email_verified, phone_verified, full_name, avatar_url,
			status, metadata
		) VALUES (
			:id, :email, :phone, :password_hash, :preferred_auth_method,
			:email_verified, :phone_verified, :full_name, :avatar_url,
			:status, :metadata
		)
	`

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, phone, password_hash, preferred_auth_method,
			   email_verified, phone_verified, full_name, avatar_url,
			   status, failed_login_attempts, locked_until, last_login,
			   last_otp_sent, created_at, updated_at, metadata
		FROM users
		WHERE id = $1 AND status != 'deleted'
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, phone, password_hash, preferred_auth_method,
			   email_verified, phone_verified, full_name, avatar_url,
			   status, failed_login_attempts, locked_until, last_login,
			   last_otp_sent, created_at, updated_at, metadata
		FROM users
		WHERE email = $1 AND status != 'deleted'
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByPhone retrieves a user by phone
func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, phone, password_hash, preferred_auth_method,
			   email_verified, phone_verified, full_name, avatar_url,
			   status, failed_login_attempts, locked_until, last_login,
			   last_otp_sent, created_at, updated_at, metadata
		FROM users
		WHERE phone = $1 AND status != 'deleted'
	`

	err := r.db.GetContext(ctx, &user, query, phone)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

// GetByEmailOrPhone retrieves a user by email or phone
func (r *userRepository) GetByEmailOrPhone(ctx context.Context, identifier string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, email, phone, password_hash, preferred_auth_method,
			   email_verified, phone_verified, full_name, avatar_url,
			   status, failed_login_attempts, locked_until, last_login,
			   last_otp_sent, created_at, updated_at, metadata
		FROM users
		WHERE (email = $1 OR phone = $1) AND status != 'deleted'
	`

	err := r.db.GetContext(ctx, &user, query, identifier)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by identifier: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = :email,
			phone = :phone,
			password_hash = :password_hash,
			preferred_auth_method = :preferred_auth_method,
			email_verified = :email_verified,
			phone_verified = :phone_verified,
			full_name = :full_name,
			avatar_url = :avatar_url,
			status = :status,
			metadata = :metadata,
			updated_at = NOW()
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdatePassword updates user's password
func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateLastLogin updates user's last login timestamp
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET last_login = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// IncrementFailedAttempts increments failed login attempts
func (r *userRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to increment failed attempts: %w", err)
	}

	return nil
}

// ResetFailedAttempts resets failed login attempts
func (r *userRepository) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to reset failed attempts: %w", err)
	}

	return nil
}

// LockAccount locks a user account until specified time
func (r *userRepository) LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error {
	query := `
		UPDATE users
		SET locked_until = $1,
			updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, until, userID)
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	return nil
}

// UnlockAccount unlocks a user account
func (r *userRepository) UnlockAccount(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE users
		SET locked_until = NULL,
			failed_login_attempts = 0,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	return nil
}

// GetUserOrganizations retrieves all organizations for a user
func (r *userRepository) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]domain.UserOrganization, error) {
	query := `
		SELECT id, user_id, organization_id, role, permissions,
			   is_primary, status, joined_at, left_at
		FROM user_organizations
		WHERE user_id = $1 AND status = 'active'
		ORDER BY is_primary DESC, joined_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user organizations: %w", err)
	}
	defer rows.Close()

	var userOrgs []domain.UserOrganization
	for rows.Next() {
		var userOrg domain.UserOrganization
		var permissions pq.StringArray
		
		err := rows.Scan(
			&userOrg.ID,
			&userOrg.UserID,
			&userOrg.OrganizationID,
			&userOrg.Role,
			&permissions,
			&userOrg.IsPrimary,
			&userOrg.Status,
			&userOrg.JoinedAt,
			&userOrg.LeftAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user organization: %w", err)
		}
		
		userOrg.Permissions = []string(permissions)
		userOrgs = append(userOrgs, userOrg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user organizations: %w", err)
	}

	return userOrgs, nil
}

// AddUserToOrganization adds a user to an organization
func (r *userRepository) AddUserToOrganization(ctx context.Context, userOrg *domain.UserOrganization) error {
	query := `
		INSERT INTO user_organizations (
			id, user_id, organization_id, role, permissions,
			is_primary, status
		) VALUES (
			:id, :user_id, :organization_id, :role, :permissions,
			:is_primary, :status
		)
	`

	if userOrg.ID == uuid.Nil {
		userOrg.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, userOrg)
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}

// RemoveUserFromOrganization removes a user from an organization
func (r *userRepository) RemoveUserFromOrganization(ctx context.Context, userID, orgID uuid.UUID) error {
	query := `
		UPDATE user_organizations
		SET status = 'inactive',
			left_at = NOW()
		WHERE user_id = $1 AND organization_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, orgID)
	if err != nil {
		return fmt.Errorf("failed to remove user from organization: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user organization not found")
	}

	return nil
}

// UpdateUserRole updates a user's role and permissions in an organization
func (r *userRepository) UpdateUserRole(ctx context.Context, userID, orgID uuid.UUID, role string, permissions []string) error {
	query := `
		UPDATE user_organizations
		SET role = $1,
			permissions = $2
		WHERE user_id = $3 AND organization_id = $4
	`

	result, err := r.db.ExecContext(ctx, query, role, permissions, userID, orgID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user organization not found")
	}

	return nil
}
