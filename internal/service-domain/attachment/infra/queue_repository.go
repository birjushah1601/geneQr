package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
)

// PostgresProcessingQueueRepository implements the ProcessingQueueRepository interface
type PostgresProcessingQueueRepository struct {
	db *pgxpool.Pool
}

// NewPostgresProcessingQueueRepository creates a new queue repository
func NewPostgresProcessingQueueRepository(db *pgxpool.Pool) *PostgresProcessingQueueRepository {
	return &PostgresProcessingQueueRepository{db: db}
}

// Enqueue adds an attachment to the processing queue
func (r *PostgresProcessingQueueRepository) Enqueue(ctx context.Context, attachmentID uuid.UUID, priority domain.QueuePriority) error {
	query := `
		INSERT INTO attachment_processing_queue (
			id, attachment_id, queued_at, status, priority, retry_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, 'pending', $4, 0, $5, $6
		)`

	now := time.Now()
	id := uuid.New()

	_, err := r.db.Exec(ctx, query, id, attachmentID, now, string(priority), now, now)
	if err != nil {
		return fmt.Errorf("failed to enqueue attachment: %w", err)
	}

	return nil
}

// Dequeue gets the next item to process from the queue
func (r *PostgresProcessingQueueRepository) Dequeue(ctx context.Context) (*domain.AttachmentProcessingQueueItem, error) {
	// Use a transaction to atomically get and lock the next item
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get the next item based on priority and creation time
	query := `
		SELECT id, attachment_id, queued_at, processed_at, status, error_message, 
			   retry_count, priority, created_at, updated_at
		FROM attachment_processing_queue
		WHERE status = 'pending'
		ORDER BY 
			CASE priority 
				WHEN 'urgent' THEN 1 
				WHEN 'high' THEN 2 
				WHEN 'medium' THEN 3 
				WHEN 'low' THEN 4 
				ELSE 5 
			END,
			created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED`

	var item domain.AttachmentProcessingQueueItem
	err = tx.QueryRow(ctx, query).Scan(
		&item.ID,
		&item.AttachmentID,
		&item.QueuedAt,
		&item.ProcessedAt,
		&item.Status,
		&item.ErrorMessage,
		&item.RetryCount,
		&item.Priority,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No items available
		}
		return nil, fmt.Errorf("failed to dequeue item: %w", err)
	}

	// Mark as processing
	updateQuery := `
		UPDATE attachment_processing_queue 
		SET status = 'processing', updated_at = $2
		WHERE id = $1`

	_, err = tx.Exec(ctx, updateQuery, item.ID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to mark item as processing: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &item, nil
}

// MarkProcessing marks a queue item as being processed
func (r *PostgresProcessingQueueRepository) MarkProcessing(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE attachment_processing_queue 
		SET status = 'processing', updated_at = $2
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark as processing: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("queue item not found")
	}

	return nil
}

// MarkCompleted marks a queue item as completed
func (r *PostgresProcessingQueueRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE attachment_processing_queue 
		SET status = 'completed', processed_at = $2, updated_at = $3
		WHERE id = $1`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, id, now, now)
	if err != nil {
		return fmt.Errorf("failed to mark as completed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("queue item not found")
	}

	return nil
}

// MarkFailed marks a queue item as failed and increments retry count
func (r *PostgresProcessingQueueRepository) MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error {
	query := `
		UPDATE attachment_processing_queue 
		SET status = 'failed', error_message = $2, retry_count = retry_count + 1, 
			processed_at = $3, updated_at = $4
		WHERE id = $1`

	now := time.Now()
	result, err := r.db.Exec(ctx, query, id, errorMessage, now, now)
	if err != nil {
		return fmt.Errorf("failed to mark as failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("queue item not found")
	}

	return nil
}

// GetQueueStats returns statistics about the processing queue
func (r *PostgresProcessingQueueRepository) GetQueueStats(ctx context.Context) (*domain.QueueStats, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing_count,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_count,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
			COALESCE(AVG(EXTRACT(EPOCH FROM (processed_at - queued_at))), 0) as avg_processing_time
		FROM attachment_processing_queue
		WHERE processed_at IS NOT NULL`

	var stats domain.QueueStats
	err := r.db.QueryRow(ctx, query).Scan(
		&stats.PendingCount,
		&stats.ProcessingCount,
		&stats.CompletedCount,
		&stats.FailedCount,
		&stats.AvgProcessingTime,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}

	return &stats, nil
}

// RetryFailed requeues failed items that haven't exceeded max retries
func (r *PostgresProcessingQueueRepository) RetryFailed(ctx context.Context, maxRetries int) error {
	query := `
		UPDATE attachment_processing_queue 
		SET status = 'pending', error_message = NULL, updated_at = $2
		WHERE status = 'failed' AND retry_count < $1`

	result, err := r.db.Exec(ctx, query, maxRetries, time.Now())
	if err != nil {
		return fmt.Errorf("failed to retry failed items: %w", err)
	}

	_ = result.RowsAffected() // We don't need to check this for retry operation

	return nil
}

// CleanupCompleted removes completed queue items older than the specified duration
func (r *PostgresProcessingQueueRepository) CleanupCompleted(ctx context.Context, olderThan time.Duration) error {
	query := `
		DELETE FROM attachment_processing_queue 
		WHERE status = 'completed' AND processed_at < $1`

	cutoff := time.Now().Add(-olderThan)
	result, err := r.db.Exec(ctx, query, cutoff)
	if err != nil {
		return fmt.Errorf("failed to cleanup completed items: %w", err)
	}

	_ = result.RowsAffected()

	return nil
}

// GetPendingCount returns the number of pending items in the queue
func (r *PostgresProcessingQueueRepository) GetPendingCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM attachment_processing_queue WHERE status = 'pending'`
	
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending count: %w", err)
	}

	return count, nil
}

// GetProcessingCount returns the number of items currently being processed
func (r *PostgresProcessingQueueRepository) GetProcessingCount(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM attachment_processing_queue WHERE status = 'processing'`
	
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get processing count: %w", err)
	}

	return count, nil
}

// GetStaleProcessingItems returns items that have been processing for too long
func (r *PostgresProcessingQueueRepository) GetStaleProcessingItems(ctx context.Context, staleAfter time.Duration) ([]*domain.AttachmentProcessingQueueItem, error) {
	query := `
		SELECT id, attachment_id, queued_at, processed_at, status, error_message, 
			   retry_count, priority, created_at, updated_at
		FROM attachment_processing_queue
		WHERE status = 'processing' AND updated_at < $1
		ORDER BY updated_at ASC`

	cutoff := time.Now().Add(-staleAfter)
	rows, err := r.db.Query(ctx, query, cutoff)
	if err != nil {
		return nil, fmt.Errorf("failed to query stale items: %w", err)
	}
	defer rows.Close()

	var items []*domain.AttachmentProcessingQueueItem
	for rows.Next() {
		item := &domain.AttachmentProcessingQueueItem{}
		err := rows.Scan(
			&item.ID,
			&item.AttachmentID,
			&item.QueuedAt,
			&item.ProcessedAt,
			&item.Status,
			&item.ErrorMessage,
			&item.RetryCount,
			&item.Priority,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stale item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}