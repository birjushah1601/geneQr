package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
)

// ProcessorConfig holds configuration for the queue processor
type ProcessorConfig struct {
	// Number of concurrent workers
	WorkerCount int
	
	// Time between queue polling attempts
	PollInterval time.Duration
	
	// Maximum processing time before considering an item stale
	StaleTimeout time.Duration
	
	// Maximum number of retry attempts
	MaxRetries int
	
	// Cleanup interval for completed items
	CleanupInterval time.Duration
	
	// Keep completed items for this duration before cleanup
	KeepCompletedFor time.Duration
}

// DefaultProcessorConfig returns sensible defaults
func DefaultProcessorConfig() *ProcessorConfig {
	return &ProcessorConfig{
		WorkerCount:      3,                // 3 concurrent workers
		PollInterval:     5 * time.Second,  // Check every 5 seconds
		StaleTimeout:     10 * time.Minute, // 10 minutes timeout
		MaxRetries:       3,                // Try 3 times
		CleanupInterval:  1 * time.Hour,    // Cleanup every hour
		KeepCompletedFor: 24 * time.Hour,   // Keep for 24 hours
	}
}

// QueueProcessor handles automated processing of attachment queue items
type QueueProcessor struct {
	config           *ProcessorConfig
	queueRepo        domain.ProcessingQueueRepository
	attachmentRepo   domain.AttachmentRepository
	visionEngine     *ai.VisionAnalysisEngine
	aiProcessor      domain.AIProcessor
	logger           *slog.Logger
	
	// Control channels
	stopCh           chan struct{}
	doneCh           chan struct{}
	
	// Worker management
	workerWg         sync.WaitGroup
	cleanupWg        sync.WaitGroup
	
	// Status tracking
	mu               sync.RWMutex
	isRunning        bool
	processedCount   int64
	errorCount       int64
	lastProcessedAt  time.Time
}

// NewQueueProcessor creates a new queue processor
func NewQueueProcessor(
	config *ProcessorConfig,
	queueRepo domain.ProcessingQueueRepository,
	attachmentRepo domain.AttachmentRepository,
	visionEngine *ai.VisionAnalysisEngine,
	aiProcessor domain.AIProcessor,
	logger *slog.Logger,
) *QueueProcessor {
	if config == nil {
		config = DefaultProcessorConfig()
	}

	return &QueueProcessor{
		config:         config,
		queueRepo:      queueRepo,
		attachmentRepo: attachmentRepo,
		visionEngine:   visionEngine,
		aiProcessor:    aiProcessor,
		logger:         logger,
		stopCh:         make(chan struct{}),
		doneCh:         make(chan struct{}),
	}
}

// Start begins processing queue items in the background
func (p *QueueProcessor) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isRunning {
		return fmt.Errorf("queue processor is already running")
	}

	p.logger.Info("Starting attachment queue processor",
		slog.Int("worker_count", p.config.WorkerCount),
		slog.Duration("poll_interval", p.config.PollInterval),
		slog.Duration("stale_timeout", p.config.StaleTimeout),
	)

	p.isRunning = true

	// Start worker goroutines
	for i := 0; i < p.config.WorkerCount; i++ {
		p.workerWg.Add(1)
		go p.worker(ctx, i)
	}

	// Start cleanup goroutine
	p.cleanupWg.Add(1)
	go p.cleanupWorker(ctx)

	// Start stale item monitor
	p.cleanupWg.Add(1)
	go p.staleItemMonitor(ctx)

	p.logger.Info("Attachment queue processor started successfully")
	return nil
}

// Stop gracefully shuts down the queue processor
func (p *QueueProcessor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isRunning {
		return nil
	}

	p.logger.Info("Stopping attachment queue processor...")

	// Signal all workers to stop
	close(p.stopCh)

	// Wait for all workers to finish
	p.workerWg.Wait()
	p.cleanupWg.Wait()

	p.isRunning = false
	close(p.doneCh)

	p.logger.Info("Attachment queue processor stopped successfully",
		slog.Int64("total_processed", p.processedCount),
		slog.Int64("total_errors", p.errorCount),
	)

	return nil
}

// GetStatus returns the current status of the processor
func (p *QueueProcessor) GetStatus() *ProcessorStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return &ProcessorStatus{
		IsRunning:       p.isRunning,
		ProcessedCount:  p.processedCount,
		ErrorCount:      p.errorCount,
		LastProcessedAt: p.lastProcessedAt,
		WorkerCount:     p.config.WorkerCount,
	}
}

// ProcessorStatus represents the current status of the processor
type ProcessorStatus struct {
	IsRunning       bool      `json:"is_running"`
	ProcessedCount  int64     `json:"processed_count"`
	ErrorCount      int64     `json:"error_count"`
	LastProcessedAt time.Time `json:"last_processed_at"`
	WorkerCount     int       `json:"worker_count"`
}

// worker is the main processing loop for each worker goroutine
func (p *QueueProcessor) worker(ctx context.Context, workerID int) {
	defer p.workerWg.Done()

	workerLogger := p.logger.With(slog.Int("worker_id", workerID))
	workerLogger.Info("Queue worker started")

	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			workerLogger.Info("Queue worker stopped due to context cancellation")
			return
		case <-p.stopCh:
			workerLogger.Info("Queue worker stopped")
			return
		case <-ticker.C:
			// Process next queue item
			p.processNextItem(ctx, workerLogger)
		}
	}
}

// processNextItem processes a single item from the queue
func (p *QueueProcessor) processNextItem(ctx context.Context, logger *slog.Logger) {
	// Get next item from queue
	queueItem, err := p.queueRepo.Dequeue(ctx)
	if err != nil {
		logger.Error("Failed to dequeue item", slog.String("error", err.Error()))
		p.incrementErrorCount()
		return
	}

	if queueItem == nil {
		// No items in queue, continue polling
		return
	}

	logger.Info("Processing queue item",
		slog.String("queue_item_id", queueItem.ID.String()),
		slog.String("attachment_id", queueItem.AttachmentID.String()),
		slog.String("priority", string(queueItem.Priority)),
	)

	// Get attachment details
	attachment, err := p.attachmentRepo.GetByID(ctx, queueItem.AttachmentID)
	if err != nil {
		logger.Error("Failed to get attachment",
			slog.String("attachment_id", queueItem.AttachmentID.String()),
			slog.String("error", err.Error()),
		)
		p.markItemFailed(ctx, queueItem.ID, fmt.Sprintf("Failed to get attachment: %v", err), logger)
		return
	}

	// Process the attachment
	err = p.processAttachment(ctx, attachment, logger)
	if err != nil {
		logger.Error("Failed to process attachment",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("error", err.Error()),
		)
		p.markItemFailed(ctx, queueItem.ID, fmt.Sprintf("Processing failed: %v", err), logger)
		return
	}

	// Mark as completed
	err = p.queueRepo.MarkCompleted(ctx, queueItem.ID)
	if err != nil {
		logger.Error("Failed to mark item as completed",
			slog.String("queue_item_id", queueItem.ID.String()),
			slog.String("error", err.Error()),
		)
		p.incrementErrorCount()
		return
	}

	// Update attachment status
	err = p.attachmentRepo.UpdateStatus(ctx, attachment.ID, domain.ProcessingStatusProcessed)
	if err != nil {
		logger.Error("Failed to update attachment status",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("error", err.Error()),
		)
		// Don't fail the queue item for this, just log
	}

	p.incrementProcessedCount()
	
	logger.Info("Successfully processed queue item",
		slog.String("queue_item_id", queueItem.ID.String()),
		slog.String("attachment_id", attachment.ID.String()),
	)
}

// processAttachment handles the actual AI processing of an attachment
func (p *QueueProcessor) processAttachment(ctx context.Context, attachment *domain.Attachment, logger *slog.Logger) error {
	// Only process images for now
	if !attachment.IsImage() {
		return fmt.Errorf("attachment is not an image: %s", attachment.FileType)
	}

	// Create AI processing request
	request := &ai.VisionAnalysisRequest{
		AttachmentID: attachment.ID,
		TicketID:     attachment.TicketID,
		ImagePath:    attachment.StoragePath,
		FileType:     attachment.FileType,
		Purpose:      "automated_analysis",
		// Equipment context would be loaded from ticket/QR code
		Equipment: p.getEquipmentContextFromTicket(attachment.TicketID),
	}

	// Process with AI
	result, err := p.visionEngine.AnalyzeImage(ctx, request)
	if err != nil {
		return fmt.Errorf("AI vision analysis failed: %w", err)
	}

	// Use AI processor to store results
	err = p.aiProcessor.ProcessAttachment(ctx, attachment.ID, result)
	if err != nil {
		return fmt.Errorf("failed to store AI analysis results: %w", err)
	}

	logger.Info("AI analysis completed",
		slog.String("attachment_id", attachment.ID.String()),
		slog.Float64("confidence", result.AnalysisConfidence),
		slog.String("quality", result.AnalysisQuality),
		slog.Int("objects_detected", len(result.DetectedObjects)),
		slog.Int("issues_detected", len(result.DetectedIssues)),
	)

	return nil
}

// getEquipmentContextFromTicket loads equipment information based on ticket
func (p *QueueProcessor) getEquipmentContextFromTicket(ticketID string) *ai.EquipmentContext {
	// This would integrate with your existing equipment/QR system
	// For now, return a default context
	return &ai.EquipmentContext{
		ID:           "unknown",
		Name:         "Medical Equipment",
		Manufacturer: "Unknown",
		Model:        "Unknown",
		Category:     "General",
		Age:          0,
	}
}

// markItemFailed marks a queue item as failed
func (p *QueueProcessor) markItemFailed(ctx context.Context, itemID uuid.UUID, errorMessage string, logger *slog.Logger) {
	err := p.queueRepo.MarkFailed(ctx, itemID, errorMessage)
	if err != nil {
		logger.Error("Failed to mark item as failed",
			slog.String("item_id", itemID.String()),
			slog.String("error", err.Error()),
		)
	}
	p.incrementErrorCount()
}

// cleanupWorker handles periodic cleanup of completed items
func (p *QueueProcessor) cleanupWorker(ctx context.Context) {
	defer p.cleanupWg.Done()

	logger := p.logger.With(slog.String("component", "cleanup_worker"))
	logger.Info("Cleanup worker started")

	ticker := time.NewTicker(p.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Cleanup worker stopped due to context cancellation")
			return
		case <-p.stopCh:
			logger.Info("Cleanup worker stopped")
			return
		case <-ticker.C:
			p.performCleanup(ctx, logger)
		}
	}
}

// performCleanup removes old completed items and retries failed ones
func (p *QueueProcessor) performCleanup(ctx context.Context, logger *slog.Logger) {
	// Cleanup completed items
	err := p.queueRepo.CleanupCompleted(ctx, p.config.KeepCompletedFor)
	if err != nil {
		logger.Error("Failed to cleanup completed items", slog.String("error", err.Error()))
	}

	// Retry failed items that haven't exceeded max retries
	err = p.queueRepo.RetryFailed(ctx, p.config.MaxRetries)
	if err != nil {
		logger.Error("Failed to retry failed items", slog.String("error", err.Error()))
	} else {
		logger.Debug("Performed cleanup and retry cycle")
	}
}

// staleItemMonitor checks for items that have been processing too long
func (p *QueueProcessor) staleItemMonitor(ctx context.Context) {
	defer p.cleanupWg.Done()

	logger := p.logger.With(slog.String("component", "stale_monitor"))
	logger.Info("Stale item monitor started")

	ticker := time.NewTicker(p.config.StaleTimeout / 2) // Check twice as often as timeout
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stale item monitor stopped due to context cancellation")
			return
		case <-p.stopCh:
			logger.Info("Stale item monitor stopped")
			return
		case <-ticker.C:
			p.checkStaleItems(ctx, logger)
		}
	}
}

// checkStaleItems finds and requeues stale processing items
func (p *QueueProcessor) checkStaleItems(ctx context.Context, logger *slog.Logger) {
	staleItems, err := p.queueRepo.GetStaleProcessingItems(ctx, p.config.StaleTimeout)
	if err != nil {
		logger.Error("Failed to get stale items", slog.String("error", err.Error()))
		return
	}

	if len(staleItems) == 0 {
		return
	}

	logger.Warn("Found stale processing items", slog.Int("count", len(staleItems)))

	for _, item := range staleItems {
		// Mark as failed with stale message
		err := p.queueRepo.MarkFailed(ctx, item.ID, "Processing timed out (stale)")
		if err != nil {
			logger.Error("Failed to mark stale item as failed",
				slog.String("item_id", item.ID.String()),
				slog.String("error", err.Error()),
			)
		} else {
			logger.Info("Marked stale item as failed",
				slog.String("item_id", item.ID.String()),
				slog.String("attachment_id", item.AttachmentID.String()),
			)
		}
	}
}

// Helper methods for thread-safe counter updates
func (p *QueueProcessor) incrementProcessedCount() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processedCount++
	p.lastProcessedAt = time.Now()
}

func (p *QueueProcessor) incrementErrorCount() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.errorCount++
}