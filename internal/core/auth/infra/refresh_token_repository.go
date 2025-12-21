package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository
func NewRefreshTokenRepository(db *sqlx.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash, device_info,
			ip_address, user_agent, expires_at
		) VALUES (
			:id, :user_id, :token_hash, :device_info,
			:ip_address, :user_agent, :expires_at
		)
	`

	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetByTokenHash retrieves a refresh token by token hash
func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	query := `
		SELECT id, user_id, token_hash, device_info,
			   ip_address, user_agent, revoked, revoked_at,
			   revoke_reason, expires_at, last_used_at,
			   usage_count, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	err := r.db.GetContext(ctx, &token, query, tokenHash)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *refreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.RefreshToken, error) {
	var tokens []domain.RefreshToken
	query := `
		SELECT id, user_id, token_hash, device_info,
			   ip_address, user_agent, revoked, revoked_at,
			   revoke_reason, expires_at, last_used_at,
			   usage_count, created_at
		FROM refresh_tokens
		WHERE user_id = $1
		  AND revoked = false
		  AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &tokens, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user refresh tokens: %w", err)
	}

	return tokens, nil
}

// UpdateLastUsed updates the last used timestamp and increments usage count
func (r *refreshTokenRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET last_used_at = NOW(),
			usage_count = usage_count + 1
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last used: %w", err)
	}

	return nil
}

// Revoke revokes a refresh token
func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true,
			revoked_at = NOW(),
			revoke_reason = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, reason, id)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID, reason string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true,
			revoked_at = NOW(),
			revoke_reason = $1
		WHERE user_id = $2 AND revoked = false
	`

	_, err := r.db.ExecContext(ctx, query, reason, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens: %w", err)
	}

	return nil
}

// DeleteExpired deletes expired refresh tokens
func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) (int, error) {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
		  AND revoked = false
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rows), nil
}
