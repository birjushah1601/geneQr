package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
)

// AIAttachmentProcessor processes attachments using AI vision analysis
type AIAttachmentProcessor struct {
	visionEngine *ai.VisionAnalysisEngine
	aiRepo       domain.AIAnalysisRepository
	logger       *slog.Logger
}

// NewAIAttachmentProcessor creates a new AI attachment processor
func NewAIAttachmentProcessor(
	visionEngine *ai.VisionAnalysisEngine,
	aiRepo domain.AIAnalysisRepository,
	logger *slog.Logger,
) *AIAttachmentProcessor {
	return &AIAttachmentProcessor{
		visionEngine: visionEngine,
		aiRepo:       aiRepo,
		logger:       logger.With(slog.String("component", "ai_processor")),
	}
}

// Process processes an attachment using AI vision analysis
func (p *AIAttachmentProcessor) Process(ctx context.Context, attachment *domain.Attachment) error {
	if !attachment.IsImage() {
		p.logger.Info("Skipping non-image attachment",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("file_type", attachment.FileType),
		)
		return nil // Not an error, just skip processing
	}

	p.logger.Info("Starting AI analysis of attachment",
		slog.String("attachment_id", attachment.ID.String()),
		slog.String("filename", attachment.Filename),
		slog.String("ticket_id", attachment.TicketID),
	)

	// Create equipment context (TODO: get from ticket/equipment service)
	equipment := p.createEquipmentContext(attachment)

	// Create AI analysis request
	request := &ai.VisionAnalysisRequest{
		AttachmentID: attachment.ID,
		TicketID:     attachment.TicketID,
		ImagePath:    attachment.StoragePath,
		FileType:     attachment.FileType,
		Equipment:    equipment,
		Purpose:      p.determinePurpose(attachment),
	}

	// Perform AI analysis
	result, err := p.visionEngine.AnalyzeImage(ctx, request)
	if err != nil {
		p.logger.Error("AI analysis failed",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("AI analysis failed: %w", err)
	}

	// Store analysis results in database
	if err := p.storeAnalysisResult(ctx, result); err != nil {
		p.logger.Error("Failed to store analysis result",
			slog.String("attachment_id", attachment.ID.String()),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to store analysis result: %w", err)
	}

	p.logger.Info("AI analysis completed successfully",
		slog.String("attachment_id", attachment.ID.String()),
		slog.Float64("confidence", result.AnalysisConfidence),
		slog.String("quality", result.AnalysisQuality),
		slog.Int("detected_objects", len(result.DetectedObjects)),
		slog.Int("detected_issues", len(result.DetectedIssues)),
	)

	return nil
}

// createEquipmentContext creates equipment context for AI analysis
// TODO: This should fetch actual equipment data from the ticket/equipment service
func (p *AIAttachmentProcessor) createEquipmentContext(attachment *domain.Attachment) *ai.EquipmentContext {
	// For now, create a generic context based on attachment category
	// In real implementation, this would query the equipment service
	return &ai.EquipmentContext{
		ID:           "unknown",
		Name:         "Medical Equipment",
		Manufacturer: "Unknown",
		Model:        "Generic Model",
		SerialNumber: "UNKNOWN",
		Category:     p.guessEquipmentCategory(attachment),
		Age:          5, // Default age
	}
}

// guessEquipmentCategory guesses equipment category from attachment
func (p *AIAttachmentProcessor) guessEquipmentCategory(attachment *domain.Attachment) string {
	switch attachment.AttachmentCategory {
	case string(domain.AttachmentCategoryEquipmentPhoto):
		return "Medical Equipment"
	case string(domain.AttachmentCategoryIssuePhoto):
		return "Diagnostic Equipment"
	case string(domain.AttachmentCategoryRepairPhoto):
		return "Maintenance Equipment"
	default:
		return "Medical Equipment"
	}
}

// determinePurpose determines the analysis purpose based on attachment category
func (p *AIAttachmentProcessor) determinePurpose(attachment *domain.Attachment) string {
	switch attachment.AttachmentCategory {
	case string(domain.AttachmentCategoryIssuePhoto):
		return "issue_evidence"
	case string(domain.AttachmentCategoryRepairPhoto):
		return "after_repair"
	case string(domain.AttachmentCategoryEquipmentPhoto):
		return "before_repair"
	default:
		return "issue_evidence"
	}
}

// storeAnalysisResult stores the AI analysis result in the database
func (p *AIAttachmentProcessor) storeAnalysisResult(ctx context.Context, result *ai.VisionAnalysisResult) error {
	// Convert the complete result to JSON for storage
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal analysis result: %w", err)
	}

	// Create database record
	analysisRecord := &domain.AIVisionAnalysis{
		ID:           uuid.New(),
		AttachmentID: result.AttachmentID,
		TicketID:     result.TicketID,
		AIProvider:   result.AIProvider,
		AIModel:      result.AIModel,

		// Store complete analysis as JSON
		AnalysisResultsJSON: string(resultJSON),

		// Extract key metrics for querying
		AnalysisConfidence: result.AnalysisConfidence,
		ImageQualityScore:  result.ImageQualityScore,
		AnalysisQuality:    result.AnalysisQuality,

		// Processing info
		ProcessingDurationMS: int(result.ProcessingDuration.Milliseconds()),
		TokensUsed:          result.TokensUsed,
		CostUSD:             result.CostUSD,

		// Status
		Status: "completed",

		// Timestamps
		AnalyzedAt: result.AnalyzedAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Store in database
	if err := p.aiRepo.Create(ctx, analysisRecord); err != nil {
		return fmt.Errorf("failed to create AI analysis record: %w", err)
	}

	return nil
}

// GetAnalysisResult retrieves and deserializes an AI analysis result
func (p *AIAttachmentProcessor) GetAnalysisResult(ctx context.Context, attachmentID uuid.UUID) (*ai.VisionAnalysisResult, error) {
	// Get analysis record from database
	record, err := p.aiRepo.GetByAttachmentID(ctx, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis record: %w", err)
	}

	// Deserialize the JSON result
	var result ai.VisionAnalysisResult
	if err := json.Unmarshal([]byte(record.AnalysisResultsJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal analysis result: %w", err)
	}

	return &result, nil
}