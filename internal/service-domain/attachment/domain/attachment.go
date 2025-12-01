package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Attachment represents a file attachment to a ticket
type Attachment struct {
	ID               uuid.UUID `json:"id" db:"id"`
    TicketID         *string   `json:"ticket_id" db:"ticket_id"`
	Filename         string    `json:"filename" db:"filename"`
	OriginalFilename string    `json:"original_filename" db:"original_filename"`
	FileType         string    `json:"file_type" db:"file_type"`
	FileSizeBytes    int64     `json:"file_size_bytes" db:"file_size_bytes"`
	StoragePath      string    `json:"storage_path" db:"storage_path"`
	UploadedByID     *string   `json:"uploaded_by_id" db:"uploaded_by_id"`
	Source           string    `json:"source" db:"source"` // whatsapp, web_upload, email
	SourceMessageID  *string   `json:"source_message_id" db:"source_message_id"`
	AttachmentCategory string  `json:"attachment_category" db:"attachment_category"` // equipment_photo, repair_photo, issue_photo, document, video, audio, other
	ProcessingStatus string    `json:"processing_status" db:"processing_status"` // pending, processing, completed, failed
	UploadedAt       time.Time `json:"uploaded_at" db:"uploaded_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// AttachmentSource represents the source of an attachment
type AttachmentSource string

const (
	AttachmentSourceWhatsApp   AttachmentSource = "whatsapp"
	AttachmentSourceWebUpload  AttachmentSource = "web_upload"
	AttachmentSourceEmail      AttachmentSource = "email"
	AttachmentSourceAPI        AttachmentSource = "api"
)

// AttachmentCategory represents the type/category of attachment
type AttachmentCategory string

const (
	AttachmentCategoryEquipmentPhoto AttachmentCategory = "equipment_photo"
	AttachmentCategoryRepairPhoto    AttachmentCategory = "repair_photo"
	AttachmentCategoryIssuePhoto     AttachmentCategory = "issue_photo"
	AttachmentCategoryDocument       AttachmentCategory = "document"
	AttachmentCategoryVideo          AttachmentCategory = "video"
	AttachmentCategoryAudio          AttachmentCategory = "audio"
	AttachmentCategoryOther          AttachmentCategory = "other"
)

// ProcessingStatus represents the processing status of an attachment
type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "pending"
	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusProcessed  ProcessingStatus = "processed"  // Successfully processed by AI
	ProcessingStatusCompleted  ProcessingStatus = "completed"
	ProcessingStatusFailed     ProcessingStatus = "failed"
)

// AttachmentProcessingQueueItem represents an item in the processing queue
type AttachmentProcessingQueueItem struct {
	ID           uuid.UUID `json:"id" db:"id"`
	AttachmentID uuid.UUID `json:"attachment_id" db:"attachment_id"`
	QueuedAt     time.Time `json:"queued_at" db:"queued_at"`
	ProcessedAt  *time.Time `json:"processed_at" db:"processed_at"`
	Status       string    `json:"status" db:"status"` // pending, processing, completed, failed
	ErrorMessage *string   `json:"error_message" db:"error_message"`
	RetryCount   int       `json:"retry_count" db:"retry_count"`
	Priority     string    `json:"priority" db:"priority"` // low, medium, high, urgent
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// QueuePriority represents the priority of a queue item
type QueuePriority string

const (
	QueuePriorityLow    QueuePriority = "low"
	QueuePriorityMedium QueuePriority = "medium"
	QueuePriorityHigh   QueuePriority = "high"
	QueuePriorityUrgent QueuePriority = "urgent"
)

// CreateAttachmentRequest represents a request to create an attachment
type CreateAttachmentRequest struct {
    TicketID         *string            `json:"ticket_id,omitempty"`
	Filename         string             `json:"filename" validate:"required"`
	OriginalFilename string             `json:"original_filename"`
	FileType         string             `json:"file_type" validate:"required"`
	FileSizeBytes    int64              `json:"file_size_bytes" validate:"min=1"`
	StoragePath      string             `json:"storage_path" validate:"required"`
	UploadedByID     *string            `json:"uploaded_by_id"`
	Source           AttachmentSource   `json:"source" validate:"required"`
	SourceMessageID  *string            `json:"source_message_id"`
	Category         AttachmentCategory `json:"category" validate:"required"`
}

// UpdateAttachmentRequest represents a request to update an attachment
type UpdateAttachmentRequest struct {
	ID               uuid.UUID          `json:"id" validate:"required"`
	ProcessingStatus *ProcessingStatus  `json:"processing_status,omitempty"`
	Category         *AttachmentCategory `json:"category,omitempty"`
}

// ListAttachmentsRequest represents a request to list attachments
type ListAttachmentsRequest struct {
	TicketID   *string             `json:"ticket_id"`
	Source     *AttachmentSource   `json:"source"`
	Category   *AttachmentCategory `json:"category"`
	Status     *ProcessingStatus   `json:"status"`
	UploadedBy *string             `json:"uploaded_by"`
	Limit      int                 `json:"limit"`
	Offset     int                 `json:"offset"`
	SortBy     string              `json:"sort_by"` // created_at, file_size, filename
	SortOrder  string              `json:"sort_order"` // asc, desc
}

// AttachmentListResult represents the result of listing attachments
type AttachmentListResult struct {
	Attachments []*Attachment `json:"attachments"`
	Total       int64         `json:"total"`
	Limit       int           `json:"limit"`
	Offset      int           `json:"offset"`
}

// NewAttachment creates a new attachment from a request
func NewAttachment(req *CreateAttachmentRequest) *Attachment {
	now := time.Now()
	return &Attachment{
		ID:               uuid.New(),
        TicketID:         req.TicketID,
		Filename:         req.Filename,
		OriginalFilename: req.OriginalFilename,
		FileType:         req.FileType,
		FileSizeBytes:    req.FileSizeBytes,
		StoragePath:      req.StoragePath,
		UploadedByID:     req.UploadedByID,
		Source:           string(req.Source),
		SourceMessageID:  req.SourceMessageID,
		AttachmentCategory: string(req.Category),
		ProcessingStatus: string(ProcessingStatusPending),
		UploadedAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// MarkProcessing marks the attachment as being processed
func (a *Attachment) MarkProcessing() {
	a.ProcessingStatus = string(ProcessingStatusProcessing)
	a.UpdatedAt = time.Now()
}

// MarkCompleted marks the attachment as processing completed
func (a *Attachment) MarkCompleted() {
	a.ProcessingStatus = string(ProcessingStatusCompleted)
	a.UpdatedAt = time.Now()
}

// MarkFailed marks the attachment as processing failed
func (a *Attachment) MarkFailed() {
	a.ProcessingStatus = string(ProcessingStatusFailed)
	a.UpdatedAt = time.Now()
}

// IsImage returns true if the attachment is an image file
func (a *Attachment) IsImage() bool {
	switch a.FileType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp":
		return true
	default:
		return false
	}
}

// IsVideo returns true if the attachment is a video file
func (a *Attachment) IsVideo() bool {
	switch a.FileType {
	case "video/mp4", "video/avi", "video/mov", "video/wmv", "video/flv", "video/webm":
		return true
	default:
		return false
	}
}

// IsDocument returns true if the attachment is a document file
func (a *Attachment) IsDocument() bool {
	switch a.FileType {
	case "application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return true
	case "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return true
	case "text/plain", "text/csv":
		return true
	default:
		return false
	}
}

// GetDisplaySize returns a human-readable file size
func (a *Attachment) GetDisplaySize() string {
	size := float64(a.FileSizeBytes)
	units := []string{"B", "KB", "MB", "GB", "TB"}
	
	for _, unit := range units {
		if size < 1024.0 {
			if unit == "B" {
				return fmt.Sprintf("%.0f %s", size, unit)
			}
			return fmt.Sprintf("%.2f %s", size, unit)
		}
		size /= 1024.0
	}
	
	return fmt.Sprintf("%.2f TB", size)
}

// AttachmentStats represents statistics about attachments
type AttachmentStats struct {
	Total         int                `json:"total"`
	ByStatus      map[string]int     `json:"by_status"`
	ByCategory    map[string]int     `json:"by_category"`
	BySource      map[string]int     `json:"by_source"`
	AvgConfidence float64            `json:"avg_confidence"`
}