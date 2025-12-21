package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type otpRepository struct {
	db *sqlx.DB
}

// NewOTPRepository creates a new instance of OTPRepository
func NewOTPRepository(db *sqlx.DB) domain.OTPRepository {
	return &otpRepository{db: db}
}

// Create creates a new OTP code
func (r *otpRepository) Create(ctx context.Context, otp *domain.OTPCode) error {
	query := `
		INSERT INTO otp_codes (
			id, user_id, email, phone, code, code_hash,
			delivery_method, purpose, expires_at,
			device_info, ip_address
		) VALUES (
			:id, :user_id, :email, :phone, :code, :code_hash,
			:delivery_method, :purpose, :expires_at,
			:device_info, :ip_address
		)
	`

	if otp.ID == uuid.Nil {
		otp.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, otp)
	if err != nil {
		return fmt.Errorf("failed to create OTP: %w", err)
	}

	return nil
}

// GetByCode retrieves an OTP by identifier and code
func (r *otpRepository) GetByCode(ctx context.Context, identifier, code string) (*domain.OTPCode, error) {
	var otp domain.OTPCode
	query := `
		SELECT id, user_id, email, phone, code, code_hash,
			   delivery_method, purpose, used, attempts,
			   expires_at, device_info, ip_address,
			   created_at, used_at
		FROM otp_codes
		WHERE (email = $1 OR phone = $1)
		  AND code_hash = $2
		  AND used = false
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &otp, query, identifier, code)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("OTP not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	return &otp, nil
}

// GetLatest retrieves the latest OTP for an identifier and purpose
func (r *otpRepository) GetLatest(ctx context.Context, identifier, purpose string) (*domain.OTPCode, error) {
	var otp domain.OTPCode
	query := `
		SELECT id, user_id, email, phone, code, code_hash,
			   delivery_method, purpose, used, attempts,
			   expires_at, device_info, ip_address,
			   created_at, used_at
		FROM otp_codes
		WHERE (email = $1 OR phone = $1)
		  AND purpose = $2
		  AND used = false
		  AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &otp, query, identifier, purpose)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active OTP found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest OTP: %w", err)
	}

	return &otp, nil
}

// MarkAsUsed marks an OTP as used
func (r *otpRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_codes
		SET used = true,
			used_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("OTP not found")
	}

	return nil
}

// IncrementAttempts increments the verification attempts for an OTP
func (r *otpRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE otp_codes
		SET attempts = attempts + 1
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment OTP attempts: %w", err)
	}

	return nil
}

// DeleteExpired deletes expired OTP codes
func (r *otpRepository) DeleteExpired(ctx context.Context) (int, error) {
	query := `
		DELETE FROM otp_codes
		WHERE expires_at < NOW()
		  AND used = false
	`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired OTPs: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rows), nil
}

// CountRecentOTPs counts OTPs sent to an identifier since a certain time
func (r *otpRepository) CountRecentOTPs(ctx context.Context, identifier string, since time.Time) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM otp_codes
		WHERE (email = $1 OR phone = $1)
		  AND created_at > $2
	`

	err := r.db.GetContext(ctx, &count, query, identifier, since)
	if err != nil {
		return 0, fmt.Errorf("failed to count recent OTPs: %w", err)
	}

	return count, nil
}
