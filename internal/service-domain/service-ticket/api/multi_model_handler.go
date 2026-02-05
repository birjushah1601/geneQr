package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/go-chi/chi/v5"
)

// MultiModelAssignmentHandler handles HTTP requests for multi-model engineer assignment
type MultiModelAssignmentHandler struct {
	service *app.MultiModelAssignmentService
	logger  *slog.Logger
}

// NewMultiModelAssignmentHandler creates a new multi-model assignment HTTP handler
func NewMultiModelAssignmentHandler(service *app.MultiModelAssignmentService, logger *slog.Logger) *MultiModelAssignmentHandler {
	return &MultiModelAssignmentHandler{
		service: service,
		logger:  logger.With(slog.String("component", "multi_model_assignment_handler")),
	}
}

// GetAssignmentSuggestions handles GET /tickets/{id}/assignment-suggestions
func (h *MultiModelAssignmentHandler) GetAssignmentSuggestions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	
	if ticketID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}
	
	h.logger.Info("Fetching assignment suggestions", slog.String("ticket_id", ticketID))
	
	suggestions, err := h.service.GetMultiModelSuggestions(ctx, ticketID)
	if err != nil {
		h.logger.Error("Failed to get assignment suggestions", 
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get suggestions: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, suggestions)
}

// respondJSON writes JSON response
func (h *MultiModelAssignmentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func (h *MultiModelAssignmentHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
