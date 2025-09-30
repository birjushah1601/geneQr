package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/rfq/domain"
	"github.com/segmentio/ksuid"
)

// RFQService orchestrates domain operations for the RFQ module
type RFQService struct {
	repository domain.RFQRepository
	eventBus   domain.EventPublisher
	logger     *slog.Logger
}

// NewRFQService creates a new RFQ service
func NewRFQService(
	repository domain.RFQRepository,
	eventBus domain.EventPublisher,
	logger *slog.Logger,
) *RFQService {
	return &RFQService{
		repository: repository,
		eventBus:   eventBus,
		logger:     logger.With(slog.String("component", "rfq_service")),
	}
}

// CreateRFQ creates a new RFQ in draft status
func (s *RFQService) CreateRFQ(ctx context.Context, req CreateRFQRequest) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	userID := domain.GetUserID(ctx)
	if userID == "" {
		userID = "system" // Fallback for development
	}

	// Generate IDs
	rfqID := ksuid.New().String()
	
	// Generate RFQ number using database function would be better,
	// but for now we'll generate it here
	rfqNumber := fmt.Sprintf("RFQ-%s-%04d", time.Now().Format("2006"), time.Now().Unix()%10000)

	// Create domain entity
	rfq, err := domain.NewRFQ(
		rfqID,
		rfqNumber,
		tenantID,
		req.Title,
		req.Description,
		req.Priority,
		req.ResponseDeadline,
		req.DeliveryTerms,
		req.PaymentTerms,
		userID,
	)

	if err != nil {
		s.logger.Error("Failed to create RFQ domain entity",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create RFQ: %w", err)
	}

	// Add items if provided
	for _, itemReq := range req.Items {
		item := domain.RFQItem{
			ID:             ksuid.New().String(),
			EquipmentID:    itemReq.EquipmentID,
			Name:           itemReq.Name,
			Description:    itemReq.Description,
			Specifications: itemReq.Specifications,
			Quantity:       itemReq.Quantity,
			Unit:           itemReq.Unit,
			EstimatedPrice: itemReq.EstimatedPrice,
			Notes:          itemReq.Notes,
			CategoryID:     itemReq.CategoryID,
		}

		if err := rfq.AddItem(item); err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}
	}

	// Set internal notes
	rfq.InternalNotes = req.InternalNotes

	// Persist to database
	if err := s.repository.Create(ctx, rfq); err != nil {
		s.logger.Error("Failed to persist RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID))
		return nil, fmt.Errorf("failed to persist RFQ: %w", err)
	}

	// Add items to database
	for _, item := range rfq.Items {
		if err := s.repository.AddItem(ctx, rfq.ID, &item); err != nil {
			s.logger.Error("Failed to add RFQ item",
				slog.String("error", err.Error()),
				slog.String("item_id", item.ID))
			// Continue adding other items
		}
	}

	// Publish domain event
	event := domain.NewRFQCreatedEvent(rfq)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish RFQ created event",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID))
		// Continue despite event publishing failure
	}

	s.logger.Info("RFQ created successfully",
		slog.String("rfq_id", rfqID),
		slog.String("rfq_number", rfqNumber),
		slog.String("tenant_id", tenantID))

	return s.mapToDTO(rfq), nil
}

// GetRFQ retrieves an RFQ by ID
func (s *RFQService) GetRFQ(ctx context.Context, id string) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			s.logger.Info("RFQ not found",
				slog.String("rfq_id", id),
				slog.String("tenant_id", tenantID))
			return nil, err
		}
		s.logger.Error("Failed to get RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to get RFQ: %w", err)
	}

	return s.mapToDTO(rfq), nil
}

// UpdateRFQ updates an existing RFQ (only in draft status)
func (s *RFQService) UpdateRFQ(ctx context.Context, id string, req UpdateRFQRequest) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch existing RFQ
	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	// Check if RFQ can be edited
	if !rfq.CanBeEdited() {
		return nil, errors.New("RFQ cannot be edited in current status")
	}

	// Update fields
	rfq.Title = req.Title
	rfq.Description = req.Description
	rfq.Priority = req.Priority
	rfq.ResponseDeadline = req.ResponseDeadline
	rfq.DeliveryTerms = req.DeliveryTerms
	rfq.PaymentTerms = req.PaymentTerms
	rfq.InternalNotes = req.InternalNotes
	rfq.UpdatedAt = time.Now()

	// Persist changes
	if err := s.repository.Update(ctx, rfq); err != nil {
		s.logger.Error("Failed to update RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to update RFQ: %w", err)
	}

	s.logger.Info("RFQ updated successfully",
		slog.String("rfq_id", id),
		slog.String("tenant_id", tenantID))

	return s.mapToDTO(rfq), nil
}

// DeleteRFQ deletes an RFQ (only in draft status)
func (s *RFQService) DeleteRFQ(ctx context.Context, id string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	// Fetch RFQ to check status
	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		return err
	}

	// Only allow deletion of draft RFQs
	if rfq.Status != domain.RFQStatusDraft {
		return errors.New("only draft RFQs can be deleted")
	}

	if err := s.repository.Delete(ctx, id, tenantID); err != nil {
		s.logger.Error("Failed to delete RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return fmt.Errorf("failed to delete RFQ: %w", err)
	}

	s.logger.Info("RFQ deleted successfully",
		slog.String("rfq_id", id),
		slog.String("tenant_id", tenantID))

	return nil
}

// ListRFQs lists RFQs with filtering and pagination
func (s *RFQService) ListRFQs(ctx context.Context, req ListRFQsRequest) (*PaginatedResponse, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// Build criteria
	criteria := domain.ListCriteria{
		TenantID:      tenantID,
		Status:        req.Status,
		Priority:      req.Priority,
		CreatedBy:     req.CreatedBy,
		SearchQuery:   req.SearchQuery,
		Page:          req.Page,
		PageSize:      req.PageSize,
		SortBy:        req.SortBy,
		SortDirection: req.SortDirection,
	}

	// Fetch from repository
	rfqs, total, err := s.repository.List(ctx, criteria)
	if err != nil {
		s.logger.Error("Failed to list RFQs",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list RFQs: %w", err)
	}

	// Convert to DTOs
	dtos := make([]*RFQDTO, len(rfqs))
	for i, rfq := range rfqs {
		dtos[i] = s.mapToDTO(rfq)
	}

	// Calculate total pages
	totalPages := total / req.PageSize
	if total%req.PageSize > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Items:      dtos,
		TotalItems: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// PublishRFQ transitions an RFQ from draft to published
func (s *RFQService) PublishRFQ(ctx context.Context, id string) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch RFQ
	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	// Publish (domain logic validates)
	if err := rfq.Publish(); err != nil {
		return nil, err
	}

	// Persist changes
	if err := s.repository.Update(ctx, rfq); err != nil {
		s.logger.Error("Failed to update RFQ status",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to publish RFQ: %w", err)
	}

	// Publish domain event
	event := domain.NewRFQPublishedEvent(rfq)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish RFQ published event",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
	}

	s.logger.Info("RFQ published successfully",
		slog.String("rfq_id", id),
		slog.String("rfq_number", rfq.RFQNumber))

	return s.mapToDTO(rfq), nil
}

// CloseRFQ closes an RFQ for new quotes
func (s *RFQService) CloseRFQ(ctx context.Context, id string) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	userID := domain.GetUserID(ctx)

	// Fetch RFQ
	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	// Close
	if err := rfq.Close(); err != nil {
		return nil, err
	}

	// Persist changes
	if err := s.repository.Update(ctx, rfq); err != nil {
		s.logger.Error("Failed to close RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to close RFQ: %w", err)
	}

	// Publish domain event
	event := domain.NewRFQClosedEvent(rfq, userID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish RFQ closed event",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
	}

	s.logger.Info("RFQ closed successfully",
		slog.String("rfq_id", id),
		slog.String("rfq_number", rfq.RFQNumber))

	return s.mapToDTO(rfq), nil
}

// CancelRFQ cancels an RFQ
func (s *RFQService) CancelRFQ(ctx context.Context, id string) (*RFQDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch RFQ
	rfq, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	// Cancel
	if err := rfq.Cancel(); err != nil {
		return nil, err
	}

	// Persist changes
	if err := s.repository.Update(ctx, rfq); err != nil {
		s.logger.Error("Failed to cancel RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to cancel RFQ: %w", err)
	}

	s.logger.Info("RFQ cancelled successfully",
		slog.String("rfq_id", id),
		slog.String("rfq_number", rfq.RFQNumber))

	return s.mapToDTO(rfq), nil
}

// AddItem adds an item to an RFQ
func (s *RFQService) AddItem(ctx context.Context, rfqID string, req AddItemRequest) (*RFQItemDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch RFQ
	rfq, err := s.repository.GetByID(ctx, rfqID, tenantID)
	if err != nil {
		return nil, err
	}

	// Create item
	item := domain.RFQItem{
		ID:             ksuid.New().String(),
		RFQID:          rfqID,
		EquipmentID:    req.EquipmentID,
		Name:           req.Name,
		Description:    req.Description,
		Specifications: req.Specifications,
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		EstimatedPrice: req.EstimatedPrice,
		Notes:          req.Notes,
		CategoryID:     req.CategoryID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Add to domain entity (validates)
	if err := rfq.AddItem(item); err != nil {
		return nil, err
	}

	// Persist
	if err := s.repository.AddItem(ctx, rfqID, &item); err != nil {
		s.logger.Error("Failed to add RFQ item",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID))
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	s.logger.Info("Item added to RFQ",
		slog.String("rfq_id", rfqID),
		slog.String("item_id", item.ID))

	return s.mapItemToDTO(&item), nil
}

// RemoveItem removes an item from an RFQ
func (s *RFQService) RemoveItem(ctx context.Context, rfqID, itemID string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	// Fetch RFQ to check status
	rfq, err := s.repository.GetByID(ctx, rfqID, tenantID)
	if err != nil {
		return err
	}

	// Check if can be edited
	if !rfq.CanBeEdited() {
		return errors.New("cannot remove items from non-draft RFQ")
	}

	// Remove from repository
	if err := s.repository.RemoveItem(ctx, rfqID, itemID); err != nil {
		s.logger.Error("Failed to remove RFQ item",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID),
			slog.String("item_id", itemID))
		return fmt.Errorf("failed to remove item: %w", err)
	}

	s.logger.Info("Item removed from RFQ",
		slog.String("rfq_id", rfqID),
		slog.String("item_id", itemID))

	return nil
}

// Helper methods

func (s *RFQService) mapToDTO(rfq *domain.RFQ) *RFQDTO {
	dto := &RFQDTO{
		ID:               rfq.ID,
		RFQNumber:        rfq.RFQNumber,
		TenantID:         rfq.TenantID,
		Title:            rfq.Title,
		Description:      rfq.Description,
		Priority:         string(rfq.Priority),
		Status:           string(rfq.Status),
		DeliveryTerms:    s.deliveryTermsToMap(rfq.DeliveryTerms),
		PaymentTerms:     s.paymentTermsToMap(rfq.PaymentTerms),
		PublishedAt:      rfq.PublishedAt,
		ResponseDeadline: rfq.ResponseDeadline,
		ClosedAt:         rfq.ClosedAt,
		CreatedBy:        rfq.CreatedBy,
		CreatedAt:        rfq.CreatedAt,
		UpdatedAt:        rfq.UpdatedAt,
		InternalNotes:    rfq.InternalNotes,
		Items:            []RFQItemDTO{},
		Invitations:      []RFQInvitationDTO{},
	}

	// Map items
	for _, item := range rfq.Items {
		dto.Items = append(dto.Items, *s.mapItemToDTO(&item))
	}

	// Map invitations
	for _, inv := range rfq.Invitations {
		dto.Invitations = append(dto.Invitations, RFQInvitationDTO{
			ID:          inv.ID,
			RFQID:       inv.RFQID,
			SupplierID:  inv.SupplierID,
			Status:      inv.Status,
			InvitedAt:   inv.InvitedAt,
			ViewedAt:    inv.ViewedAt,
			RespondedAt: inv.RespondedAt,
			Message:     inv.Message,
		})
	}

	return dto
}

func (s *RFQService) mapItemToDTO(item *domain.RFQItem) *RFQItemDTO {
	return &RFQItemDTO{
		ID:             item.ID,
		RFQID:          item.RFQID,
		EquipmentID:    item.EquipmentID,
		Name:           item.Name,
		Description:    item.Description,
		Specifications: item.Specifications,
		Quantity:       item.Quantity,
		Unit:           item.Unit,
		EstimatedPrice: item.EstimatedPrice,
		Notes:          item.Notes,
		CategoryID:     item.CategoryID,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

// deliveryTermsToMap converts DeliveryTerms struct to map for JSON serialization
func (s *RFQService) deliveryTermsToMap(terms domain.DeliveryTerms) map[string]interface{} {
	return map[string]interface{}{
		"address":              terms.Address,
		"city":                 terms.City,
		"state":                terms.State,
		"postal_code":          terms.PostalCode,
		"country":              terms.Country,
		"required_by":          terms.RequiredBy,
		"special_notes":        terms.SpecialNotes,
		"installation_required": terms.InstallationReq,
	}
}

// paymentTermsToMap converts PaymentTerms struct to map for JSON serialization
func (s *RFQService) paymentTermsToMap(terms domain.PaymentTerms) map[string]interface{} {
	return map[string]interface{}{
		"payment_method":         terms.PaymentMethod,
		"payment_days":           terms.PaymentDays,
		"advance_payment_percent": terms.AdvancePayment,
		"special_terms":          terms.SpecialTerms,
	}
}
