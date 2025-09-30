package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/go-chi/chi/v5"
)

// TicketHandler handles HTTP requests for service tickets
type TicketHandler struct {
	service *app.TicketService
	logger  *slog.Logger
}

// NewTicketHandler creates a new ticket HTTP handler
func NewTicketHandler(service *app.TicketService, logger *slog.Logger) *TicketHandler {
	return &TicketHandler{
		service: service,
		logger:  logger.With(slog.String("component", "ticket_handler")),
	}
}

// CreateTicket handles POST /tickets
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req app.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	ticket, err := h.service.CreateTicket(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to create ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, ticket)
}

// GetTicket handles GET /tickets/{id}
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	ticket, err := h.service.GetTicket(ctx, id)
	if err != nil {
		if err == domain.ErrTicketNotFound {
			h.respondError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		h.logger.Error("Failed to get ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get ticket")
		return
	}

	h.respondJSON(w, http.StatusOK, ticket)
}

// GetTicketByNumber handles GET /tickets/number/{number}
func (h *TicketHandler) GetTicketByNumber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketNumber := chi.URLParam(r, "number")

	if ticketNumber == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket number is required")
		return
	}

	ticket, err := h.service.GetTicketByNumber(ctx, ticketNumber)
	if err != nil {
		if err == domain.ErrTicketNotFound {
			h.respondError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		h.logger.Error("Failed to get ticket by number", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get ticket")
		return
	}

	h.respondJSON(w, http.StatusOK, ticket)
}

// ListTickets handles GET /tickets
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	criteria := domain.ListCriteria{
		EquipmentID: r.URL.Query().Get("equipment_id"),
		CustomerID:  r.URL.Query().Get("customer_id"),
		EngineerID:  r.URL.Query().Get("engineer_id"),
		SortBy:      r.URL.Query().Get("sort_by"),
		SortDirection: r.URL.Query().Get("sort_dir"),
	}

	// Parse status filter (multiple values)
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		criteria.Status = []domain.TicketStatus{domain.TicketStatus(statusStr)}
	}

	// Parse priority filter
	if priorityStr := r.URL.Query().Get("priority"); priorityStr != "" {
		criteria.Priority = []domain.TicketPriority{domain.TicketPriority(priorityStr)}
	}

	// Parse source filter
	if sourceStr := r.URL.Query().Get("source"); sourceStr != "" {
		criteria.Source = []domain.TicketSource{domain.TicketSource(sourceStr)}
	}

	// Parse pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	criteria.Page = page

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}
	criteria.PageSize = pageSize

	// Parse boolean filters
	if slaBreached := r.URL.Query().Get("sla_breached"); slaBreached != "" {
		val := slaBreached == "true"
		criteria.SLABreached = &val
	}

	if coveredUnderAMC := r.URL.Query().Get("covered_under_amc"); coveredUnderAMC != "" {
		val := coveredUnderAMC == "true"
		criteria.CoveredUnderAMC = &val
	}

	result, err := h.service.ListTickets(ctx, criteria)
	if err != nil {
		h.logger.Error("Failed to list tickets", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list tickets")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// AssignTicket handles POST /tickets/{id}/assign
func (h *TicketHandler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		EngineerID   string `json:"engineer_id"`
		EngineerName string `json:"engineer_name"`
		AssignedBy   string `json:"assigned_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.AssignTicket(ctx, id, req.EngineerID, req.EngineerName, req.AssignedBy); err != nil {
		h.logger.Error("Failed to assign ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to assign ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket assigned successfully"})
}

// AcknowledgeTicket handles POST /tickets/{id}/acknowledge
func (h *TicketHandler) AcknowledgeTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.AcknowledgeTicket(ctx, id, req.AcknowledgedBy); err != nil {
		h.logger.Error("Failed to acknowledge ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to acknowledge ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket acknowledged successfully"})
}

// StartWork handles POST /tickets/{id}/start
func (h *TicketHandler) StartWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		StartedBy string `json:"started_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.StartWork(ctx, id, req.StartedBy); err != nil {
		h.logger.Error("Failed to start work on ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to start work: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Work started successfully"})
}

// PutOnHold handles POST /tickets/{id}/hold
func (h *TicketHandler) PutOnHold(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		Reason    string `json:"reason"`
		ChangedBy string `json:"changed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.PutOnHold(ctx, id, req.Reason, req.ChangedBy); err != nil {
		h.logger.Error("Failed to put ticket on hold", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to put on hold: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket put on hold successfully"})
}

// ResumeWork handles POST /tickets/{id}/resume
func (h *TicketHandler) ResumeWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		ResumedBy string `json:"resumed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.ResumeWork(ctx, id, req.ResumedBy); err != nil {
		h.logger.Error("Failed to resume work on ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to resume work: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Work resumed successfully"})
}

// ResolveTicket handles POST /tickets/{id}/resolve
func (h *TicketHandler) ResolveTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.ResolveTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.ResolveTicket(ctx, id, req); err != nil {
		h.logger.Error("Failed to resolve ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to resolve ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket resolved successfully"})
}

// CloseTicket handles POST /tickets/{id}/close
func (h *TicketHandler) CloseTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		ClosedBy string `json:"closed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CloseTicket(ctx, id, req.ClosedBy); err != nil {
		h.logger.Error("Failed to close ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to close ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket closed successfully"})
}

// CancelTicket handles POST /tickets/{id}/cancel
func (h *TicketHandler) CancelTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		Reason      string `json:"reason"`
		CancelledBy string `json:"cancelled_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CancelTicket(ctx, id, req.Reason, req.CancelledBy); err != nil {
		h.logger.Error("Failed to cancel ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to cancel ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket cancelled successfully"})
}

// AddComment handles POST /tickets/{id}/comments
func (h *TicketHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	req.TicketID = id

	if err := h.service.AddComment(ctx, req); err != nil {
		h.logger.Error("Failed to add comment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to add comment")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Comment added successfully"})
}

// GetComments handles GET /tickets/{id}/comments
func (h *TicketHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	comments, err := h.service.GetComments(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get comments", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	h.respondJSON(w, http.StatusOK, comments)
}

// GetStatusHistory handles GET /tickets/{id}/history
func (h *TicketHandler) GetStatusHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	history, err := h.service.GetStatusHistory(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get status history", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get status history")
		return
	}

	h.respondJSON(w, http.StatusOK, history)
}

// respondJSON writes JSON response
func (h *TicketHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func (h *TicketHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
