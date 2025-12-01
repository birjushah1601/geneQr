package domain

import (
	"context"
	"fmt"
	"log/slog"
    "strings"
    "os"
    "path/filepath"

	"github.com/google/uuid"
)

// AttachmentService handles business logic for attachments
type AttachmentService struct {
	attachmentRepo AttachmentRepository
	queueRepo      ProcessingQueueRepository
	aiRepo         AIAnalysisRepository
	logger         *slog.Logger
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(
	attachmentRepo AttachmentRepository,
	queueRepo ProcessingQueueRepository,
	aiRepo AIAnalysisRepository,
	logger *slog.Logger,
) *AttachmentService {
	return &AttachmentService{
		attachmentRepo: attachmentRepo,
		queueRepo:      queueRepo,
		aiRepo:         aiRepo,
		logger:         logger.With(slog.String("service", "attachment")),
	}
}

// CreateAttachment creates a new attachment and optionally queues it for processing
func (s *AttachmentService) CreateAttachment(ctx context.Context, req *CreateAttachmentRequest) (*Attachment, error) {
	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create attachment
	attachment := NewAttachment(req)
	
    if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
        s.logger.Error("Failed to create attachment",
            slog.String("error", err.Error()),
            slog.String("ticket_id", func() string { if req.TicketID==nil {return ""}; return *req.TicketID }()),
            slog.String("filename", req.Filename),
        )
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}

    s.logger.Info("Attachment created successfully",
        slog.String("attachment_id", attachment.ID.String()),
        slog.String("ticket_id", func() string { if attachment.TicketID==nil {return ""}; return *attachment.TicketID }()),
		slog.String("filename", attachment.Filename),
		slog.String("source", attachment.Source),
	)

	// Queue for processing if it's an image
	if attachment.IsImage() {
		priority := s.determinePriority(attachment)
		if err := s.queueRepo.Enqueue(ctx, attachment.ID, priority); err != nil {
			s.logger.Error("Failed to queue attachment for processing",
				slog.String("attachment_id", attachment.ID.String()),
				slog.String("error", err.Error()),
			)
			// Don't return error - attachment was created successfully
		} else {
			s.logger.Info("Attachment queued for AI analysis",
				slog.String("attachment_id", attachment.ID.String()),
				slog.String("priority", string(priority)),
			)
		}
	}

	return attachment, nil
}

// GetAttachment retrieves an attachment by ID
func (s *AttachmentService) GetAttachment(ctx context.Context, id uuid.UUID) (*Attachment, error) {
	attachment, err := s.attachmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	return attachment, nil
}

// GetAttachmentsForTicket retrieves all attachments for a ticket
func (s *AttachmentService) GetAttachmentsForTicket(ctx context.Context, ticketID string) ([]*Attachment, error) {
	attachments, err := s.attachmentRepo.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments for ticket: %w", err)
	}
	return attachments, nil
}

// ListAttachments retrieves attachments with filtering and pagination
func (s *AttachmentService) ListAttachments(ctx context.Context, req *ListAttachmentsRequest) (*AttachmentListResult, error) {
	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	result, err := s.attachmentRepo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	return result, nil
}

// UpdateAttachment updates an attachment
func (s *AttachmentService) UpdateAttachment(ctx context.Context, req *UpdateAttachmentRequest) error {
	attachment, err := s.attachmentRepo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	// Update fields if provided
	if req.ProcessingStatus != nil {
		attachment.ProcessingStatus = string(*req.ProcessingStatus)
	}
	if req.Category != nil {
		attachment.AttachmentCategory = string(*req.Category)
	}

	if err := s.attachmentRepo.Update(ctx, attachment); err != nil {
		return fmt.Errorf("failed to update attachment: %w", err)
	}

	s.logger.Info("Attachment updated successfully",
		slog.String("attachment_id", attachment.ID.String()),
	)

	return nil
}

// DeleteAttachment removes an attachment
func (s *AttachmentService) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	// TODO: Also remove physical file from storage
	if err := s.attachmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	s.logger.Info("Attachment deleted successfully",
		slog.String("attachment_id", id.String()),
	)

	return nil
}

// ProcessNextInQueue processes the next attachment in the queue
func (s *AttachmentService) ProcessNextInQueue(ctx context.Context, processor AttachmentProcessor) error {
	// Get next item from queue
	queueItem, err := s.queueRepo.Dequeue(ctx)
	if err != nil {
		return fmt.Errorf("failed to dequeue item: %w", err)
	}
	if queueItem == nil {
		return nil // No items in queue
	}

	// Get attachment details
	attachment, err := s.attachmentRepo.GetByID(ctx, queueItem.AttachmentID)
	if err != nil {
		// Mark queue item as failed
		s.queueRepo.MarkFailed(ctx, queueItem.ID, fmt.Sprintf("attachment not found: %v", err))
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	// Mark as processing
	if err := s.queueRepo.MarkProcessing(ctx, queueItem.ID); err != nil {
		return fmt.Errorf("failed to mark as processing: %w", err)
	}

	attachment.MarkProcessing()
	s.attachmentRepo.Update(ctx, attachment)

	s.logger.Info("Processing attachment",
		slog.String("attachment_id", attachment.ID.String()),
		slog.String("filename", attachment.Filename),
	)

	// Process the attachment
	if err := processor.Process(ctx, attachment); err != nil {
		// Mark as failed
		s.queueRepo.MarkFailed(ctx, queueItem.ID, err.Error())
		attachment.MarkFailed()
		s.attachmentRepo.Update(ctx, attachment)
		
		s.logger.Error("Failed to process attachment",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("processing failed: %w", err)
	}

	// Mark as completed
	s.queueRepo.MarkCompleted(ctx, queueItem.ID)
	attachment.MarkCompleted()
	s.attachmentRepo.Update(ctx, attachment)

	s.logger.Info("Attachment processed successfully",
		slog.String("attachment_id", attachment.ID.String()),
	)

	return nil
}

// GetQueueStats returns processing queue statistics
func (s *AttachmentService) GetQueueStats(ctx context.Context) (*QueueStats, error) {
	stats, err := s.queueRepo.GetQueueStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}
	return stats, nil
}

// RetryFailedItems requeues failed items for retry
func (s *AttachmentService) RetryFailedItems(ctx context.Context, maxRetries int) error {
	if err := s.queueRepo.RetryFailed(ctx, maxRetries); err != nil {
		return fmt.Errorf("failed to retry failed items: %w", err)
	}
	
	s.logger.Info("Retried failed queue items",
		slog.Int("max_retries", maxRetries),
	)
	
	return nil
}

// GetAIAnalysis retrieves AI analysis results for an attachment
func (s *AttachmentService) GetAIAnalysis(ctx context.Context, attachmentID uuid.UUID) (*AIVisionAnalysis, error) {
	analysis, err := s.aiRepo.GetByAttachmentID(ctx, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI analysis: %w", err)
	}
	return analysis, nil
}

// GetAIAnalysesForTicket retrieves all AI analyses for a ticket
func (s *AttachmentService) GetAIAnalysesForTicket(ctx context.Context, ticketID string) ([]*AIVisionAnalysis, error) {
	analyses, err := s.aiRepo.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI analyses for ticket: %w", err)
	}
	return analyses, nil
}

// GetAttachmentStats retrieves attachment statistics
func (s *AttachmentService) GetAttachmentStats(ctx context.Context) (*AttachmentStats, error) {
	stats, err := s.attachmentRepo.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment stats: %w", err)
	}
	return stats, nil
}

// validateCreateRequest validates a create attachment request
func (s *AttachmentService) validateCreateRequest(req *CreateAttachmentRequest) error {
	if req.Filename == "" {
		return fmt.Errorf("filename is required")
	}
	if req.FileType == "" {
		return fmt.Errorf("file_type is required")
	}
	if req.FileSizeBytes <= 0 {
		return fmt.Errorf("file_size_bytes must be positive")
	}
	if req.StoragePath == "" {
		return fmt.Errorf("storage_path is required")
	}
	if req.Source == "" {
		return fmt.Errorf("source is required")
	}
	if req.Category == "" {
		return fmt.Errorf("category is required")
	}

	// Validate file size limits (50MB default)
	maxSize := int64(50 * 1024 * 1024)
	if req.FileSizeBytes > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", req.FileSizeBytes, maxSize)
	}

	return nil
}

// LinkAttachmentToTicket associates an existing attachment with a ticket and updates storage path if moved
func (s *AttachmentService) LinkAttachmentToTicket(ctx context.Context, id uuid.UUID, ticketID string) error {
    // Fetch attachment
    att, err := s.attachmentRepo.GetByID(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to get attachment: %w", err)
    }

    // Compute new storage path if current path is in unassigned bucket and move file on disk
    var newPath *string
    if att.StoragePath != "" && strings.Contains(att.StoragePath, "/attachments/unassigned/") {
        targetPath := strings.Replace(att.StoragePath, "/attachments/unassigned/", "/attachments/"+ticketID+"/", 1)
        // Ensure target directory exists
        if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
            return fmt.Errorf("failed to create target directory: %w", err)
        }
        // Move file
        if err := os.Rename(att.StoragePath, targetPath); err != nil {
            return fmt.Errorf("failed to move attachment file: %w", err)
        }
        newPath = &targetPath
    }

    if err := s.attachmentRepo.LinkToTicket(ctx, id, ticketID, newPath); err != nil {
        return fmt.Errorf("failed to link attachment to ticket: %w", err)
    }
    return nil
}

// determinePriority determines processing priority based on attachment properties
func (s *AttachmentService) determinePriority(attachment *Attachment) QueuePriority {
	switch attachment.AttachmentCategory {
	case string(AttachmentCategoryIssuePhoto):
		return QueuePriorityHigh
	case string(AttachmentCategoryEquipmentPhoto), string(AttachmentCategoryRepairPhoto):
		return QueuePriorityMedium
	default:
		return QueuePriorityLow
	}
}

// AttachmentProcessor defines the interface for processing attachments
type AttachmentProcessor interface {
	Process(ctx context.Context, attachment *Attachment) error
}