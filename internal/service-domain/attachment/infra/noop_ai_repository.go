package infra

import (
    "context"
    "errors"
    "github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
    "github.com/google/uuid"
)

// ErrNotImplemented is returned by the noop AI repository for all operations
var ErrNotImplemented = errors.New("ai analysis repository not configured")

// NoopAIAnalysisRepository is a placeholder implementation used when AI storage is not configured
type NoopAIAnalysisRepository struct{}

func NewNoopAIAnalysisRepository() *NoopAIAnalysisRepository { return &NoopAIAnalysisRepository{} }

func (r *NoopAIAnalysisRepository) Create(ctx context.Context, analysis *domain.AIVisionAnalysis) error {
    return ErrNotImplemented
}

func (r *NoopAIAnalysisRepository) GetByAttachmentID(ctx context.Context, attachmentID uuid.UUID) (*domain.AIVisionAnalysis, error) {
    return nil, ErrNotImplemented
}

func (r *NoopAIAnalysisRepository) GetByTicketID(ctx context.Context, ticketID string) ([]*domain.AIVisionAnalysis, error) {
    return nil, ErrNotImplemented
}

func (r *NoopAIAnalysisRepository) Update(ctx context.Context, analysis *domain.AIVisionAnalysis) error {
    return ErrNotImplemented
}

func (r *NoopAIAnalysisRepository) List(ctx context.Context, req *domain.ListAIAnalysisRequest) (*domain.AIAnalysisListResult, error) {
    return nil, ErrNotImplemented
}
