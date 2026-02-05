package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type auditRepository struct {
	db *sqlx.DB
}

// NewAuditRepository creates a new instance of AuditRepository
func NewAuditRepository(db *sqlx.DB) domain.AuditRepository {
	return &auditRepository{db: db}
}

// Log creates a new audit log entry
func (r *auditRepository) Log(ctx context.Context, log *domain.AuthAuditLog) error {
	query := `
		INSERT INTO auth_audit_log (
			id, user_id, action, success, ip_address,
			user_agent, device_info, metadata, error_message
		) VALUES (
			:id, :user_id, :action, :success, :ip_address,
			:user_agent, :device_info, :metadata, :error_message
		)
	`

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetByUserID retrieves audit logs for a user with pagination
func (r *auditRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.AuthAuditLog, error) {
	var logs []domain.AuthAuditLog
	query := `
		SELECT id, user_id, action, success, ip_address,
			   user_agent, device_info, metadata, error_message,
			   created_at
		FROM auth_audit_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by user: %w", err)
	}

	return logs, nil
}

// GetFailedLoginsByIP counts failed login attempts from an IP address
func (r *auditRepository) GetFailedLoginsByIP(ctx context.Context, ipAddress string, since time.Time) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM auth_audit_log
		WHERE action = $1
		  AND success = false
		  AND ip_address = $2
		  AND created_at > $3
	`

	err := r.db.GetContext(ctx, &count, query, domain.AuditActionLoginFailed, ipAddress, since)
	if err != nil {
		return 0, fmt.Errorf("failed to get failed logins by IP: %w", err)
	}

	return count, nil
}

// GetRecentActivity retrieves recent activity for a user
func (r *auditRepository) GetRecentActivity(ctx context.Context, userID uuid.UUID, limit int) ([]domain.AuthAuditLog, error) {
	var logs []domain.AuthAuditLog
	query := `
		SELECT id, user_id, action, success, ip_address,
			   user_agent, device_info, metadata, error_message,
			   created_at
		FROM auth_audit_log
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	err := r.db.SelectContext(ctx, &logs, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}

	return logs, nil
}
