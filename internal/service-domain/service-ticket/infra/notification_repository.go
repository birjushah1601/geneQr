package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepository implements domain.NotificationRepository
type NotificationRepository struct {
	pool *pgxpool.Pool
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		pool: pool,
	}
}

// CreateTrackingToken creates a new tracking token for a ticket
func (r *NotificationRepository) CreateTrackingToken(token *domain.TrackingToken) error {
	ctx := context.Background()
	
	query := `
		INSERT INTO ticket_tracking_tokens (id, ticket_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	if token.ID == "" {
		token.ID = uuid.New().String()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}
	
	_, err := r.pool.Exec(ctx, query,
		token.ID,
		token.TicketID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create tracking token: %w", err)
	}
	
	return nil
}

// GetTrackingToken retrieves a tracking token by token string
func (r *NotificationRepository) GetTrackingToken(token string) (*domain.TrackingToken, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, ticket_id, token, expires_at, created_at
		FROM ticket_tracking_tokens
		WHERE token = $1 AND expires_at > NOW()
	`
	
	var trackingToken domain.TrackingToken
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&trackingToken.ID,
		&trackingToken.TicketID,
		&trackingToken.Token,
		&trackingToken.ExpiresAt,
		&trackingToken.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking token: %w", err)
	}
	
	return &trackingToken, nil
}

// DeleteExpiredTokens removes expired tracking tokens
func (r *NotificationRepository) DeleteExpiredTokens() error {
	ctx := context.Background()
	
	query := `DELETE FROM ticket_tracking_tokens WHERE expires_at <= NOW()`
	
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}
	
	return nil
}

// LogNotification logs a notification attempt
func (r *NotificationRepository) LogNotification(log *domain.NotificationLog) error {
	ctx := context.Background()
	
	query := `
		INSERT INTO notification_log (
			id, ticket_id, notification_type, recipient_email, 
			sent_at, status, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	if log.ID == "" {
		log.ID = uuid.New().String()
	}
	if log.SentAt.IsZero() {
		log.SentAt = time.Now()
	}
	
	_, err := r.pool.Exec(ctx, query,
		log.ID,
		log.TicketID,
		log.NotificationType,
		log.RecipientEmail,
		log.SentAt,
		log.Status,
		log.ErrorMessage,
	)
	
	if err != nil {
		return fmt.Errorf("failed to log notification: %w", err)
	}
	
	return nil
}

// GetNotificationHistory retrieves notification history for a ticket
func (r *NotificationRepository) GetNotificationHistory(ticketID string) ([]*domain.NotificationLog, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, ticket_id, notification_type, recipient_email, 
		       sent_at, status, error_message
		FROM notification_log
		WHERE ticket_id = $1
		ORDER BY sent_at DESC
	`
	
	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification history: %w", err)
	}
	defer rows.Close()
	
	var logs []*domain.NotificationLog
	for rows.Next() {
		var log domain.NotificationLog
		err := rows.Scan(
			&log.ID,
			&log.TicketID,
			&log.NotificationType,
			&log.RecipientEmail,
			&log.SentAt,
			&log.Status,
			&log.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification log: %w", err)
		}
		logs = append(logs, &log)
	}
	
	return logs, nil
}

// GetFailedNotifications retrieves failed notifications since a given time
func (r *NotificationRepository) GetFailedNotifications(since time.Time) ([]*domain.NotificationLog, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, ticket_id, notification_type, recipient_email, 
		       sent_at, status, error_message
		FROM notification_log
		WHERE status = $1 AND sent_at >= $2
		ORDER BY sent_at DESC
	`
	
	rows, err := r.pool.Query(ctx, query, domain.NotificationStatusFailed, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed notifications: %w", err)
	}
	defer rows.Close()
	
	var logs []*domain.NotificationLog
	for rows.Next() {
		var log domain.NotificationLog
		err := rows.Scan(
			&log.ID,
			&log.TicketID,
			&log.NotificationType,
			&log.RecipientEmail,
			&log.SentAt,
			&log.Status,
			&log.ErrorMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification log: %w", err)
		}
		logs = append(logs, &log)
	}
	
	return logs, nil
}
