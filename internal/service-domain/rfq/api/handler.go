package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/rfq/app"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq/domain"
	"github.com/go-chi/chi/v5"
)

// RFQHandler handles HTTP requests for RFQ operations
type RFQHandler struct {
	service *app.RFQService
	logger  *slog.Logger
}

// NewRFQHandler creates a new RFQ HTTP handler
func NewRFQHandler(service *app.RFQService, logger *slog.Logger) *RFQHandler {
	return &RFQHandler{
		service: service,
		logger:  logger.With(slog.String("component", "rfq_handler")),
	}
}

// CreateRFQ handles POST /api/v1/rfq
func (h *RFQHandler) CreateRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req app.CreateRFQRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Basic validation
	if req.Title == "" || req.Description == "" {
		h.respondError(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	rfq, err := h.service.CreateRFQ(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to create RFQ: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, app.APIResponse{
		Success: true,
		Message: "RFQ created successfully",
		Data:    rfq,
	})
}

// GetRFQ handles GET /api/v1/rfq/{id}
func (h *RFQHandler) GetRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	rfq, err := h.service.GetRFQ(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to get RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get RFQ: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Data:    rfq,
	})
}

// UpdateRFQ handles PUT /api/v1/rfq/{id}
func (h *RFQHandler) UpdateRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	var req app.UpdateRFQRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	rfq, err := h.service.UpdateRFQ(ctx, id, req)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to update RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "RFQ updated successfully",
		Data:    rfq,
	})
}

// DeleteRFQ handles DELETE /api/v1/rfq/{id}
func (h *RFQHandler) DeleteRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	err := h.service.DeleteRFQ(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to delete RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "RFQ deleted successfully",
	})
}

// ListRFQs handles GET /api/v1/rfq
func (h *RFQHandler) ListRFQs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	query := r.URL.Query()

	// Pagination
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))

	// Build request
	req := app.ListRFQsRequest{
		Page:          page,
		PageSize:      pageSize,
		SearchQuery:   query.Get("search"),
		CreatedBy:     query.Get("created_by"),
		SortBy:        query.Get("sort_by"),
		SortDirection: query.Get("sort_direction"),
	}

	// Parse status filter (can be multiple)
	if statusStr := query.Get("status"); statusStr != "" {
		req.Status = []domain.RFQStatus{domain.RFQStatus(statusStr)}
	}

	// Parse priority filter (can be multiple)
	if priorityStr := query.Get("priority"); priorityStr != "" {
		req.Priority = []domain.RFQPriority{domain.RFQPriority(priorityStr)}
	}

	result, err := h.service.ListRFQs(ctx, req)
	if err != nil {
		h.logger.Error("Failed to list RFQs", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list RFQs: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Data:    result,
	})
}

// PublishRFQ handles POST /api/v1/rfq/{id}/publish
func (h *RFQHandler) PublishRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	rfq, err := h.service.PublishRFQ(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to publish RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "RFQ published successfully",
		Data:    rfq,
	})
}

// CloseRFQ handles POST /api/v1/rfq/{id}/close
func (h *RFQHandler) CloseRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	rfq, err := h.service.CloseRFQ(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to close RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "RFQ closed successfully",
		Data:    rfq,
	})
}

// CancelRFQ handles POST /api/v1/rfq/{id}/cancel
func (h *RFQHandler) CancelRFQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	rfq, err := h.service.CancelRFQ(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to cancel RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "RFQ cancelled successfully",
		Data:    rfq,
	})
}

// AddItem handles POST /api/v1/rfq/{id}/items
func (h *RFQHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rfqID := chi.URLParam(r, "id")

	if rfqID == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID is required")
		return
	}

	var req app.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	item, err := h.service.AddItem(ctx, rfqID, req)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to add item", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, app.APIResponse{
		Success: true,
		Message: "Item added successfully",
		Data:    item,
	})
}

// RemoveItem handles DELETE /api/v1/rfq/{id}/items/{item_id}
func (h *RFQHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rfqID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "item_id")

	if rfqID == "" || itemID == "" {
		h.respondError(w, http.StatusBadRequest, "RFQ ID and Item ID are required")
		return
	}

	err := h.service.RemoveItem(ctx, rfqID, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrRFQNotFound) {
			h.respondError(w, http.StatusNotFound, "RFQ not found")
			return
		}
		h.logger.Error("Failed to remove item", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, app.APIResponse{
		Success: true,
		Message: "Item removed successfully",
	})
}

// Helper methods

func (h *RFQHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *RFQHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, app.APIResponse{
		Success: false,
		Error:   message,
	})
}
