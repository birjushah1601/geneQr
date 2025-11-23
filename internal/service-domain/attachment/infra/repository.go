package infra

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresAttachmentRepository implements domain.AttachmentRepository using PostgreSQL
type PostgresAttachmentRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewPostgresAttachmentRepository creates a new PostgreSQL attachment repository
func NewPostgresAttachmentRepository(db *pgxpool.Pool) domain.AttachmentRepository {
	return &PostgresAttachmentRepository{
		db:     db,
		logger: slog.With(slog.String("component", "postgres_attachment_repository")),
	}
}

// Create stores a new attachment
func (r *PostgresAttachmentRepository) Create(ctx context.Context, attachment *domain.Attachment) error {
	query := `
		INSERT INTO ticket_attachments (
			id, ticket_id, filename, file_type, file_size_bytes, storage_path, 
			attachment_category, source, processing_status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)`

	_, err := r.db.Exec(ctx, query,
		attachment.ID,
		attachment.TicketID,
		attachment.Filename,
		attachment.FileType,
		attachment.FileSizeBytes,
		attachment.StoragePath,
		attachment.AttachmentCategory,
		attachment.Source,
		attachment.ProcessingStatus,
		attachment.CreatedAt,
		attachment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create attachment: %w", err)
	}

	r.logger.Info("Attachment created successfully", slog.String("id", attachment.ID.String()))
	return nil
}

// GetByID retrieves an attachment by its ID
func (r *PostgresAttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attachment, error) {
	query := `
		SELECT id, ticket_id, filename, file_type, file_size_bytes, storage_path, 
		       attachment_category, source, processing_status, created_at, updated_at
		FROM ticket_attachments 
		WHERE id = $1`

	var attachment domain.Attachment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&attachment.ID,
		&attachment.TicketID,
		&attachment.Filename,
		&attachment.FileType,
		&attachment.FileSizeBytes,
		&attachment.StoragePath,
		&attachment.AttachmentCategory,
		&attachment.Source,
		&attachment.ProcessingStatus,
		&attachment.CreatedAt,
		&attachment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get attachment by ID: %w", err)
	}

	return &attachment, nil
}

// GetByTicketID retrieves all attachments for a specific ticket
func (r *PostgresAttachmentRepository) GetByTicketID(ctx context.Context, ticketID string) ([]*domain.Attachment, error) {
	query := `
		SELECT id, ticket_id, filename, file_type, file_size_bytes, storage_path, 
		       attachment_category, source, processing_status, created_at, updated_at
		FROM ticket_attachments 
		WHERE ticket_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attachments by ticket ID: %w", err)
	}
	defer rows.Close()

	var attachments []*domain.Attachment
	for rows.Next() {
		var attachment domain.Attachment
		if err := rows.Scan(
			&attachment.ID,
			&attachment.TicketID,
			&attachment.Filename,
			&attachment.FileType,
			&attachment.FileSizeBytes,
			&attachment.StoragePath,
			&attachment.AttachmentCategory,
			&attachment.Source,
			&attachment.ProcessingStatus,
			&attachment.CreatedAt,
			&attachment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, &attachment)
	}

	return attachments, nil
}

// List retrieves attachments based on criteria with pagination
func (r *PostgresAttachmentRepository) List(ctx context.Context, req *domain.ListAttachmentsRequest) (*domain.AttachmentListResult, error) {
	// Count total first
	countQuery := `SELECT COUNT(*) FROM ticket_attachments WHERE 1=1`
	var countArgs []interface{}
	
	// Apply same filters for counting
	whereConditions := []string{}
	argIndex := 1
	
	if req.TicketID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("ticket_id = $%d", argIndex))
		countArgs = append(countArgs, *req.TicketID)
		argIndex++
	}
	
	if len(whereConditions) > 0 {
		countQuery += " AND " + fmt.Sprintf("%s", whereConditions[0])
	}
	
	var total int64
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		total = 0 // Don't fail on count error
	}

	// Build main query
	query := `
		SELECT id, ticket_id, filename, file_type, file_size_bytes, storage_path, 
		       attachment_category, source, processing_status, created_at, updated_at
		FROM ticket_attachments 
		WHERE 1=1`

	var args []interface{}
	argIndex = 1

	if req.TicketID != nil {
		query += fmt.Sprintf(" AND ticket_id = $%d", argIndex)
		args = append(args, *req.TicketID)
		argIndex++
	}

	// Add ordering
	query += " ORDER BY created_at DESC"
	
	// Add pagination
	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, req.Limit)
		argIndex++
		
		if req.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, req.Offset)
		}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query attachments: %w", err)
	}
	defer rows.Close()

	var attachments []*domain.Attachment
	for rows.Next() {
		var attachment domain.Attachment
		if err := rows.Scan(
			&attachment.ID,
			&attachment.TicketID,
			&attachment.Filename,
			&attachment.FileType,
			&attachment.FileSizeBytes,
			&attachment.StoragePath,
			&attachment.AttachmentCategory,
			&attachment.Source,
			&attachment.ProcessingStatus,
			&attachment.CreatedAt,
			&attachment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, &attachment)
	}

	return &domain.AttachmentListResult{
		Attachments: attachments,
		Total:       total,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}, nil
}

// Update modifies an existing attachment
func (r *PostgresAttachmentRepository) Update(ctx context.Context, attachment *domain.Attachment) error {
	query := `
		UPDATE ticket_attachments 
		SET filename = $2, file_type = $3, file_size_bytes = $4, storage_path = $5,
		    attachment_category = $6, source = $7, processing_status = $8, updated_at = $9
		WHERE id = $1`

	attachment.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		attachment.ID,
		attachment.Filename,
		attachment.FileType,
		attachment.FileSizeBytes,
		attachment.StoragePath,
		attachment.AttachmentCategory,
		attachment.Source,
		attachment.ProcessingStatus,
		attachment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update attachment: %w", err)
	}

	r.logger.Info("Attachment updated successfully", slog.String("id", attachment.ID.String()))
	return nil
}

// Delete removes an attachment
func (r *PostgresAttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ticket_attachments WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	r.logger.Info("Attachment deleted successfully", slog.String("id", id.String()))
	return nil
}

// UpdateStatus updates only the processing status of an attachment
func (r *PostgresAttachmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProcessingStatus) error {
	query := `
		UPDATE ticket_attachments 
		SET processing_status = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id, string(status), time.Now())
	if err != nil {
		return fmt.Errorf("failed to update attachment status: %w", err)
	}

	r.logger.Info("Attachment status updated", slog.String("id", id.String()), slog.String("status", string(status)))
	return nil
}

// GetPendingForProcessing retrieves attachments that need processing
func (r *PostgresAttachmentRepository) GetPendingForProcessing(ctx context.Context, limit int) ([]*domain.Attachment, error) {
	query := `
		SELECT id, ticket_id, filename, file_type, file_size_bytes, storage_path, 
		       attachment_category, source, processing_status, created_at, updated_at
		FROM ticket_attachments 
		WHERE processing_status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending attachments: %w", err)
	}
	defer rows.Close()

	var attachments []*domain.Attachment
	for rows.Next() {
		var attachment domain.Attachment
		if err := rows.Scan(
			&attachment.ID,
			&attachment.TicketID,
			&attachment.Filename,
			&attachment.FileType,
			&attachment.FileSizeBytes,
			&attachment.StoragePath,
			&attachment.AttachmentCategory,
			&attachment.Source,
			&attachment.ProcessingStatus,
			&attachment.CreatedAt,
			&attachment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}
		attachments = append(attachments, &attachment)
	}

	return attachments, nil
}

// GetStats retrieves attachment statistics
func (r *PostgresAttachmentRepository) GetStats(ctx context.Context) (*domain.AttachmentStats, error) {
	stats := &domain.AttachmentStats{
		ByStatus:   make(map[string]int),
		ByCategory: make(map[string]int),
		BySource:   make(map[string]int),
	}

	// Return mock stats for now (database might not exist)
	stats.Total = 0
	stats.AvgConfidence = 0.85

	return stats, nil
}