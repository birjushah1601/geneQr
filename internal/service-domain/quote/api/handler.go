package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/quote/app"
	"github.com/aby-med/medical-platform/internal/service-domain/quote/domain"
	"github.com/go-chi/chi/v5"
)

// QuoteHandler handles HTTP requests for quotes
type QuoteHandler struct {
	service *app.QuoteService
	logger  *slog.Logger
}

// NewQuoteHandler creates a new quote handler
func NewQuoteHandler(service *app.QuoteService, logger *slog.Logger) *QuoteHandler {
	return &QuoteHandler{
		service: service,
		logger:  logger.With(slog.String("component", "quote_handler")),
	}
}

// CreateQuote handles POST /quotes
func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req app.CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdBy := r.Header.Get("X-User-ID")
	if createdBy == "" {
		createdBy = "system"
	}

	quote, err := h.service.CreateQuote(r.Context(), tenantID, createdBy, req)
	if err != nil {
		h.logger.Error("Failed to create quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quote)
}

// GetQuote handles GET /quotes/{id}
func (h *QuoteHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	quote, err := h.service.GetQuote(r.Context(), tenantID, quoteID)
	if err != nil {
		if err == domain.ErrQuoteNotFound {
			http.Error(w, "Quote not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// ListQuotes handles GET /quotes
func (h *QuoteHandler) ListQuotes(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	criteria := domain.ListCriteria{
		TenantID: tenantID,
	}

	if rfqID := r.URL.Query().Get("rfq_id"); rfqID != "" {
		criteria.RFQID = rfqID
	}
	if supplierID := r.URL.Query().Get("supplier_id"); supplierID != "" {
		criteria.SupplierID = supplierID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		criteria.Status = []domain.QuoteStatus{domain.QuoteStatus(status)}
	}
	if minAmount := r.URL.Query().Get("min_amount"); minAmount != "" {
		if val, err := strconv.ParseFloat(minAmount, 64); err == nil {
			criteria.MinAmount = val
		}
	}
	if maxAmount := r.URL.Query().Get("max_amount"); maxAmount != "" {
		if val, err := strconv.ParseFloat(maxAmount, 64); err == nil {
			criteria.MaxAmount = val
		}
	}
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		criteria.SortBy = sortBy
	}
	if sortDir := r.URL.Query().Get("sort_direction"); sortDir != "" {
		criteria.SortDirection = sortDir
	}
	if page := r.URL.Query().Get("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil {
			criteria.Page = val
		}
	}
	if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
		if val, err := strconv.Atoi(pageSize); err == nil {
			criteria.PageSize = val
		}
	}

	result, err := h.service.ListQuotes(r.Context(), criteria)
	if err != nil {
		h.logger.Error("Failed to list quotes", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetQuotesByRFQ handles GET /rfqs/{rfq_id}/quotes
func (h *QuoteHandler) GetQuotesByRFQ(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	rfqID := chi.URLParam(r, "rfq_id")
	if rfqID == "" {
		http.Error(w, "RFQ ID required", http.StatusBadRequest)
		return
	}

	quotes, err := h.service.GetQuotesByRFQ(r.Context(), tenantID, rfqID)
	if err != nil {
		h.logger.Error("Failed to get quotes by RFQ", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

// GetQuotesBySupplier handles GET /suppliers/{supplier_id}/quotes
func (h *QuoteHandler) GetQuotesBySupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	supplierID := chi.URLParam(r, "supplier_id")
	if supplierID == "" {
		http.Error(w, "Supplier ID required", http.StatusBadRequest)
		return
	}

	quotes, err := h.service.GetQuotesBySupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		h.logger.Error("Failed to get quotes by supplier", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

// UpdateQuote handles PATCH /quotes/{id}
func (h *QuoteHandler) UpdateQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	var req app.UpdateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.UpdateQuote(r.Context(), tenantID, quoteID, req)
	if err != nil {
		if err == domain.ErrQuoteNotFound {
			http.Error(w, "Quote not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// AddQuoteItem handles POST /quotes/{id}/items
func (h *QuoteHandler) AddQuoteItem(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	var req app.QuoteItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.AddQuoteItem(r.Context(), tenantID, quoteID, req)
	if err != nil {
		h.logger.Error("Failed to add item", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// UpdateQuoteItem handles PATCH /quotes/{id}/items/{item_id}
func (h *QuoteHandler) UpdateQuoteItem(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "item_id")
	if quoteID == "" || itemID == "" {
		http.Error(w, "quote ID and item ID required", http.StatusBadRequest)
		return
	}

	var req app.QuoteItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.UpdateQuoteItem(r.Context(), tenantID, quoteID, itemID, req)
	if err != nil {
		h.logger.Error("Failed to update item", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// RemoveQuoteItem handles DELETE /quotes/{id}/items/{item_id}
func (h *QuoteHandler) RemoveQuoteItem(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "item_id")
	if quoteID == "" || itemID == "" {
		http.Error(w, "quote ID and item ID required", http.StatusBadRequest)
		return
	}

	quote, err := h.service.RemoveQuoteItem(r.Context(), tenantID, quoteID, itemID)
	if err != nil {
		h.logger.Error("Failed to remove item", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// SubmitQuote handles POST /quotes/{id}/submit
func (h *QuoteHandler) SubmitQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	quote, err := h.service.SubmitQuote(r.Context(), tenantID, quoteID)
	if err != nil {
		h.logger.Error("Failed to submit quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// ReviseQuote handles POST /quotes/{id}/revise
func (h *QuoteHandler) ReviseQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	var req app.ReviseQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.ReviseQuote(r.Context(), tenantID, quoteID, req)
	if err != nil {
		h.logger.Error("Failed to revise quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// AcceptQuote handles POST /quotes/{id}/accept
func (h *QuoteHandler) AcceptQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	var req app.AcceptQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.AcceptQuote(r.Context(), tenantID, quoteID, req)
	if err != nil {
		h.logger.Error("Failed to accept quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// RejectQuote handles POST /quotes/{id}/reject
func (h *QuoteHandler) RejectQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	var req app.RejectQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	quote, err := h.service.RejectQuote(r.Context(), tenantID, quoteID, req)
	if err != nil {
		h.logger.Error("Failed to reject quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// WithdrawQuote handles POST /quotes/{id}/withdraw
func (h *QuoteHandler) WithdrawQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	quote, err := h.service.WithdrawQuote(r.Context(), tenantID, quoteID)
	if err != nil {
		h.logger.Error("Failed to withdraw quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// MarkUnderReview handles POST /quotes/{id}/under-review
func (h *QuoteHandler) MarkUnderReview(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	quote, err := h.service.MarkQuoteUnderReview(r.Context(), tenantID, quoteID)
	if err != nil {
		h.logger.Error("Failed to mark quote under review", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// DeleteQuote handles DELETE /quotes/{id}
func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "quote ID required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteQuote(r.Context(), tenantID, quoteID); err != nil {
		if err == domain.ErrQuoteNotFound {
			http.Error(w, "Quote not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
