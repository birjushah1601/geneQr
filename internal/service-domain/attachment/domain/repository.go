package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AttachmentRepository defines data access operations for attachments
type AttachmentRepository interface {
	// Create creates a new attachment record
	Create(ctx context.Context, attachment *Attachment) error
	
	// GetByID retrieves an attachment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Attachment, error)
	
	// GetByTicketID retrieves all attachments for a ticket
	GetByTicketID(ctx context.Context, ticketID string) ([]*Attachment, error)
	
	// List retrieves attachments based on criteria with pagination
	List(ctx context.Context, req *ListAttachmentsRequest) (*AttachmentListResult, error)
	
	// Update updates an existing attachment
	Update(ctx context.Context, attachment *Attachment) error
	
    // LinkToTicket sets/changes the ticket_id and optionally updates storage_path
    LinkToTicket(ctx context.Context, id uuid.UUID, ticketID string, newStoragePath *string) error

	// Delete removes an attachment record
	Delete(ctx context.Context, id uuid.UUID) error
	
	// UpdateStatus updates only the processing status of an attachment
	UpdateStatus(ctx context.Context, id uuid.UUID, status ProcessingStatus) error
	
	// GetPendingForProcessing retrieves attachments that need AI processing
	GetPendingForProcessing(ctx context.Context, limit int) ([]*Attachment, error)
	
	// GetStats retrieves attachment statistics
	GetStats(ctx context.Context) (*AttachmentStats, error)
}

// ProcessingQueueRepository defines operations for the attachment processing queue
type ProcessingQueueRepository interface {
	// Enqueue adds an attachment to the processing queue
	Enqueue(ctx context.Context, attachmentID uuid.UUID, priority QueuePriority) error
	
	// Dequeue gets the next item to process from the queue
	Dequeue(ctx context.Context) (*AttachmentProcessingQueueItem, error)
	
	// MarkProcessing marks a queue item as being processed
	MarkProcessing(ctx context.Context, id uuid.UUID) error
	
	// MarkCompleted marks a queue item as completed
	MarkCompleted(ctx context.Context, id uuid.UUID) error
	
	// MarkFailed marks a queue item as failed and increments retry count
	MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error
	
	// GetQueueStats returns statistics about the processing queue
	GetQueueStats(ctx context.Context) (*QueueStats, error)
	
	// RetryFailed requeues failed items that haven't exceeded max retries
	RetryFailed(ctx context.Context, maxRetries int) error
	
	// CleanupCompleted removes completed queue items older than the specified duration
	CleanupCompleted(ctx context.Context, olderThan time.Duration) error
	
	// GetStaleProcessingItems returns items that have been processing for too long
	GetStaleProcessingItems(ctx context.Context, staleAfter time.Duration) ([]*AttachmentProcessingQueueItem, error)
}

// QueueStats represents statistics about the processing queue
type QueueStats struct {
	PendingCount    int64 `json:"pending_count"`
	ProcessingCount int64 `json:"processing_count"`
	CompletedCount  int64 `json:"completed_count"`
	FailedCount     int64 `json:"failed_count"`
	AvgProcessingTime float64 `json:"avg_processing_time_seconds"`
}

// AIAnalysisRepository defines operations for AI vision analysis results
type AIAnalysisRepository interface {
	// Create stores a new AI vision analysis result
	Create(ctx context.Context, analysis *AIVisionAnalysis) error
	
	// GetByAttachmentID retrieves AI analysis for an attachment
	GetByAttachmentID(ctx context.Context, attachmentID uuid.UUID) (*AIVisionAnalysis, error)
	
	// GetByTicketID retrieves all AI analyses for a ticket
	GetByTicketID(ctx context.Context, ticketID string) ([]*AIVisionAnalysis, error)
	
	// Update updates an existing analysis record
	Update(ctx context.Context, analysis *AIVisionAnalysis) error
	
	// List retrieves AI analyses with filtering and pagination
	List(ctx context.Context, req *ListAIAnalysisRequest) (*AIAnalysisListResult, error)
}

// AIVisionAnalysis represents AI vision analysis results stored in database
type AIVisionAnalysis struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	AttachmentID uuid.UUID  `json:"attachment_id" db:"attachment_id"`
	TicketID     string     `json:"ticket_id" db:"ticket_id"`
	AIProvider   string     `json:"ai_provider" db:"ai_provider"`
	AIModel      string     `json:"ai_model" db:"ai_model"`
	
	// Analysis Results (stored as JSON)
	AnalysisResultsJSON string `json:"-" db:"analysis_results_json"`
	
	// Key metrics extracted for quick querying
	AnalysisConfidence float64 `json:"analysis_confidence" db:"analysis_confidence"`
	ImageQualityScore  float64 `json:"image_quality_score" db:"image_quality_score"`
	AnalysisQuality    string  `json:"analysis_quality" db:"analysis_quality"`
	
	// Processing info
	ProcessingDurationMS int     `json:"processing_duration_ms" db:"processing_duration_ms"`
	TokensUsed          int     `json:"tokens_used" db:"tokens_used"`
	CostUSD             float64 `json:"cost_usd" db:"cost_usd"`
	
	// Status
	Status       string     `json:"status" db:"status"` // completed, failed
	ErrorMessage *string    `json:"error_message" db:"error_message"`
	
	// Timestamps
	AnalyzedAt   time.Time  `json:"analyzed_at" db:"analyzed_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ListAIAnalysisRequest represents a request to list AI analyses
type ListAIAnalysisRequest struct {
	TicketID           *string  `json:"ticket_id"`
	AttachmentID       *uuid.UUID `json:"attachment_id"`
	AIProvider         *string  `json:"ai_provider"`
	MinConfidence      *float64 `json:"min_confidence"`
	AnalysisQuality    *string  `json:"analysis_quality"`
	Status             *string  `json:"status"`
	Limit              int      `json:"limit"`
	Offset             int      `json:"offset"`
	SortBy             string   `json:"sort_by"` // analyzed_at, confidence, cost
	SortOrder          string   `json:"sort_order"` // asc, desc
}

// AIAnalysisListResult represents the result of listing AI analyses
type AIAnalysisListResult struct {
	Analyses []*AIVisionAnalysis `json:"analyses"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

// AIProcessor defines operations for processing attachments with AI
type AIProcessor interface {
	// ProcessAttachment processes an attachment with AI and stores the results
	ProcessAttachment(ctx context.Context, attachmentID uuid.UUID, analysisResult interface{}) error
}
