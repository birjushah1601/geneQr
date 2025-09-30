package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/quote/domain"
)

// QuoteService provides application-level quote operations
type QuoteService struct {
	repo   domain.QuoteRepository
	logger *slog.Logger
}

// NewQuoteService creates a new quote service
func NewQuoteService(repo domain.QuoteRepository, logger *slog.Logger) *QuoteService {
	return &QuoteService{
		repo:   repo,
		logger: logger.With(slog.String("component", "quote_service")),
	}
}

// CreateQuote creates a new quote in draft status
func (s *QuoteService) CreateQuote(ctx context.Context, tenantID string, createdBy string, req CreateQuoteRequest) (*QuoteResponse, error) {
	// Create quote
	quote, err := domain.NewQuote(tenantID, req.RFQID, req.SupplierID, createdBy, req.ValidUntil)
	if err != nil {
		s.logger.Error("Failed to create quote", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	// Set terms
	quote.DeliveryTerms = req.DeliveryTerms
	quote.PaymentTerms = req.PaymentTerms
	quote.WarrantyTerms = req.WarrantyTerms
	quote.Notes = req.Notes

	// Add items
	for _, itemReq := range req.Items {
		item := domain.QuoteItem{
			RFQItemID:         itemReq.RFQItemID,
			EquipmentID:       itemReq.EquipmentID,
			EquipmentName:     itemReq.EquipmentName,
			Quantity:          itemReq.Quantity,
			UnitPrice:         itemReq.UnitPrice,
			TaxRate:           itemReq.TaxRate,
			DeliveryTimeframe: itemReq.DeliveryTimeframe,
			ManufacturerName:  itemReq.ManufacturerName,
			ModelNumber:       itemReq.ModelNumber,
			Specifications:    itemReq.Specifications,
			ComplianceCerts:   itemReq.ComplianceCerts,
			Notes:             itemReq.Notes,
		}

		if err := quote.AddItem(item); err != nil {
			s.logger.Error("Failed to add item to quote", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to add item: %w", err)
		}
	}

	// Persist
	if err := s.repo.Create(ctx, quote); err != nil {
		s.logger.Error("Failed to persist quote", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to persist quote: %w", err)
	}

	s.logger.Info("Quote created successfully",
		slog.String("quote_id", quote.ID),
		slog.String("rfq_id", quote.RFQID),
		slog.String("supplier_id", quote.SupplierID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// GetQuote retrieves a quote by ID
func (s *QuoteService) GetQuote(ctx context.Context, tenantID string, quoteID string) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		s.logger.Error("Failed to get quote",
			slog.String("error", err.Error()),
			slog.String("quote_id", quoteID))
		return nil, err
	}

	response := ToQuoteResponse(quote)
	return &response, nil
}

// GetQuotesByRFQ retrieves all quotes for a specific RFQ
func (s *QuoteService) GetQuotesByRFQ(ctx context.Context, tenantID string, rfqID string) ([]QuoteResponse, error) {
	quotes, err := s.repo.GetByRFQID(ctx, rfqID, tenantID)
	if err != nil {
		s.logger.Error("Failed to get quotes by RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID))
		return nil, err
	}

	return ToQuoteResponses(quotes), nil
}

// GetQuotesBySupplier retrieves all quotes from a specific supplier
func (s *QuoteService) GetQuotesBySupplier(ctx context.Context, tenantID string, supplierID string) ([]QuoteResponse, error) {
	quotes, err := s.repo.GetBySupplierID(ctx, supplierID, tenantID)
	if err != nil {
		s.logger.Error("Failed to get quotes by supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, err
	}

	return ToQuoteResponses(quotes), nil
}

// ListQuotes retrieves quotes with filtering and pagination
func (s *QuoteService) ListQuotes(ctx context.Context, criteria domain.ListCriteria) (*ListQuotesResponse, error) {
	quotes, total, err := s.repo.List(ctx, criteria)
	if err != nil {
		s.logger.Error("Failed to list quotes", slog.String("error", err.Error()))
		return nil, err
	}

	page := criteria.Page
	if page == 0 {
		page = 1
	}
	pageSize := criteria.PageSize
	if pageSize == 0 {
		pageSize = 20
	}

	return &ListQuotesResponse{
		Quotes: ToQuoteResponses(quotes),
		Total:  total,
		Page:   page,
		Size:   pageSize,
	}, nil
}

// UpdateQuote updates quote details (only for drafts/revised quotes)
func (s *QuoteService) UpdateQuote(ctx context.Context, tenantID string, quoteID string, req UpdateQuoteRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if !quote.IsEditable() {
		return nil, fmt.Errorf("quote is not editable in status: %s", quote.Status)
	}

	// Update fields
	quote.DeliveryTerms = req.DeliveryTerms
	quote.PaymentTerms = req.PaymentTerms
	quote.WarrantyTerms = req.WarrantyTerms
	quote.Notes = req.Notes
	if req.ValidUntil != nil {
		quote.ValidUntil = *req.ValidUntil
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote updated successfully", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// AddQuoteItem adds a new item to a quote
func (s *QuoteService) AddQuoteItem(ctx context.Context, tenantID string, quoteID string, req QuoteItemRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	item := domain.QuoteItem{
		RFQItemID:         req.RFQItemID,
		EquipmentID:       req.EquipmentID,
		EquipmentName:     req.EquipmentName,
		Quantity:          req.Quantity,
		UnitPrice:         req.UnitPrice,
		TaxRate:           req.TaxRate,
		DeliveryTimeframe: req.DeliveryTimeframe,
		ManufacturerName:  req.ManufacturerName,
		ModelNumber:       req.ModelNumber,
		Specifications:    req.Specifications,
		ComplianceCerts:   req.ComplianceCerts,
		Notes:             req.Notes,
	}

	if err := quote.AddItem(item); err != nil {
		s.logger.Error("Failed to add item to quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Item added to quote", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// RemoveQuoteItem removes an item from a quote
func (s *QuoteService) RemoveQuoteItem(ctx context.Context, tenantID string, quoteID string, itemID string) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.RemoveItem(itemID); err != nil {
		s.logger.Error("Failed to remove item from quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Item removed from quote", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// UpdateQuoteItem updates an existing item in a quote
func (s *QuoteService) UpdateQuoteItem(ctx context.Context, tenantID string, quoteID string, itemID string, req QuoteItemRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	item := domain.QuoteItem{
		RFQItemID:         req.RFQItemID,
		EquipmentID:       req.EquipmentID,
		EquipmentName:     req.EquipmentName,
		Quantity:          req.Quantity,
		UnitPrice:         req.UnitPrice,
		TaxRate:           req.TaxRate,
		DeliveryTimeframe: req.DeliveryTimeframe,
		ManufacturerName:  req.ManufacturerName,
		ModelNumber:       req.ModelNumber,
		Specifications:    req.Specifications,
		ComplianceCerts:   req.ComplianceCerts,
		Notes:             req.Notes,
	}

	if err := quote.UpdateItem(itemID, item); err != nil {
		s.logger.Error("Failed to update item in quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Item updated in quote", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// SubmitQuote submits a quote for review
func (s *QuoteService) SubmitQuote(ctx context.Context, tenantID string, quoteID string) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.Submit(); err != nil {
		s.logger.Error("Failed to submit quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote submitted successfully",
		slog.String("quote_id", quote.ID),
		slog.String("quote_number", quote.QuoteNumber))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// ReviseQuote creates a revision of the quote
func (s *QuoteService) ReviseQuote(ctx context.Context, tenantID string, quoteID string, req ReviseQuoteRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.Revise(req.Changes, req.RevisedBy); err != nil {
		s.logger.Error("Failed to revise quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote revised successfully",
		slog.String("quote_id", quote.ID),
		slog.Int("revision_number", quote.RevisionNumber))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// AcceptQuote marks a quote as accepted
func (s *QuoteService) AcceptQuote(ctx context.Context, tenantID string, quoteID string, req AcceptQuoteRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.Accept(req.ReviewedBy, req.Notes); err != nil {
		s.logger.Error("Failed to accept quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote accepted successfully",
		slog.String("quote_id", quote.ID),
		slog.String("reviewed_by", req.ReviewedBy))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// RejectQuote marks a quote as rejected
func (s *QuoteService) RejectQuote(ctx context.Context, tenantID string, quoteID string, req RejectQuoteRequest) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.Reject(req.ReviewedBy, req.Reason); err != nil {
		s.logger.Error("Failed to reject quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote rejected successfully",
		slog.String("quote_id", quote.ID),
		slog.String("reviewed_by", req.ReviewedBy))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// WithdrawQuote allows supplier to withdraw a quote
func (s *QuoteService) WithdrawQuote(ctx context.Context, tenantID string, quoteID string) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.Withdraw(); err != nil {
		s.logger.Error("Failed to withdraw quote", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote withdrawn successfully", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// MarkQuoteUnderReview marks a quote as under review
func (s *QuoteService) MarkQuoteUnderReview(ctx context.Context, tenantID string, quoteID string) (*QuoteResponse, error) {
	quote, err := s.repo.GetByID(ctx, quoteID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := quote.MarkUnderReview(); err != nil {
		s.logger.Error("Failed to mark quote under review", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.Update(ctx, quote); err != nil {
		s.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return nil, err
	}

	s.logger.Info("Quote marked under review", slog.String("quote_id", quote.ID))

	response := ToQuoteResponse(quote)
	return &response, nil
}

// DeleteQuote deletes a quote
func (s *QuoteService) DeleteQuote(ctx context.Context, tenantID string, quoteID string) error {
	if err := s.repo.Delete(ctx, quoteID, tenantID); err != nil {
		s.logger.Error("Failed to delete quote", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("Quote deleted successfully", slog.String("quote_id", quoteID))
	return nil
}
