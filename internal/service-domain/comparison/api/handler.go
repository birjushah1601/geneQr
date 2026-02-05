package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/comparison/app"
	"github.com/aby-med/medical-platform/internal/service-domain/comparison/domain"
	"github.com/go-chi/chi/v5"
)

// ComparisonHandler handles HTTP requests for comparisons
type ComparisonHandler struct {
	service *app.ComparisonService
	logger  *slog.Logger
}

// NewComparisonHandler creates a new comparison handler
func NewComparisonHandler(service *app.ComparisonService, logger *slog.Logger) *ComparisonHandler {
	return &ComparisonHandler{
		service: service,
		logger:  logger.With(slog.String("component", "comparison_handler")),
	}
}

// CreateComparison handles POST /comparisons
func (h *ComparisonHandler) CreateComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req app.CreateComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdBy := r.Header.Get("X-User-ID")
	if createdBy == "" {
		createdBy = "system"
	}

	comparison, err := h.service.CreateComparison(r.Context(), tenantID, createdBy, req)
	if err != nil {
		h.logger.Error("Failed to create comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// GetComparison handles GET /comparisons/{id}
func (h *ComparisonHandler) GetComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.GetComparison(r.Context(), tenantID, id)
	if err != nil {
		if err == domain.ErrComparisonNotFound {
			http.Error(w, "Comparison not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// ListComparisons handles GET /comparisons
func (h *ComparisonHandler) ListComparisons(w http.ResponseWriter, r *http.Request) {
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
	if status := r.URL.Query().Get("status"); status != "" {
		criteria.Status = []domain.ComparisonStatus{domain.ComparisonStatus(status)}
	}
	if createdBy := r.URL.Query().Get("created_by"); createdBy != "" {
		criteria.CreatedBy = createdBy
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

	result, err := h.service.ListComparisons(r.Context(), criteria)
	if err != nil {
		h.logger.Error("Failed to list comparisons", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetComparisonsByRFQ handles GET /rfqs/{rfq_id}/comparisons
func (h *ComparisonHandler) GetComparisonsByRFQ(w http.ResponseWriter, r *http.Request) {
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

	comparisons, err := h.service.GetComparisonsByRFQ(r.Context(), tenantID, rfqID)
	if err != nil {
		h.logger.Error("Failed to get comparisons by RFQ", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToListResponse(comparisons))
}

// UpdateComparison handles PATCH /comparisons/{id}
func (h *ComparisonHandler) UpdateComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	var req app.UpdateComparisonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.UpdateComparison(r.Context(), tenantID, id, req)
	if err != nil {
		if err == domain.ErrComparisonNotFound {
			http.Error(w, "Comparison not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to update comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// UpdateScoringCriteria handles POST /comparisons/{id}/scoring
func (h *ComparisonHandler) UpdateScoringCriteria(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	var req app.UpdateScoringCriteriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.UpdateScoringCriteria(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to update scoring criteria", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// AddQuote handles POST /comparisons/{id}/quotes
func (h *ComparisonHandler) AddQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	var req app.AddQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.AddQuote(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to add quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// RemoveQuote handles DELETE /comparisons/{id}/quotes/{quote_id}
func (h *ComparisonHandler) RemoveQuote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	quoteID := chi.URLParam(r, "quote_id")
	if id == "" || quoteID == "" {
		http.Error(w, "comparison ID and quote ID required", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.RemoveQuote(r.Context(), tenantID, id, quoteID)
	if err != nil {
		h.logger.Error("Failed to remove quote", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// CalculateScores handles POST /comparisons/{id}/calculate
func (h *ComparisonHandler) CalculateScores(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	// Expect quote data in request body
	var quotes []app.Quote
	if err := json.NewDecoder(r.Body).Decode(&quotes); err != nil {
		http.Error(w, "Invalid request body - expected array of quotes", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.CalculateScores(r.Context(), tenantID, id, quotes)
	if err != nil {
		h.logger.Error("Failed to calculate scores", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// ActivateComparison handles POST /comparisons/{id}/activate
func (h *ComparisonHandler) ActivateComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.ActivateComparison(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to activate comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// CompleteComparison handles POST /comparisons/{id}/complete
func (h *ComparisonHandler) CompleteComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.CompleteComparison(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to complete comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// ArchiveComparison handles POST /comparisons/{id}/archive
func (h *ComparisonHandler) ArchiveComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	comparison, err := h.service.ArchiveComparison(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to archive comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.ToResponse(comparison))
}

// DeleteComparison handles DELETE /comparisons/{id}
func (h *ComparisonHandler) DeleteComparison(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "comparison ID required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteComparison(r.Context(), tenantID, id); err != nil {
		if err == domain.ErrComparisonNotFound {
			http.Error(w, "Comparison not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to delete comparison", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
