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

// AssignmentHandler handles HTTP requests for engineer assignments
type AssignmentHandler struct {
	service *app.AssignmentService
	logger  *slog.Logger
}

// NewAssignmentHandler creates a new assignment HTTP handler
func NewAssignmentHandler(service *app.AssignmentService, logger *slog.Logger) *AssignmentHandler {
	return &AssignmentHandler{
		service: service,
		logger:  logger.With(slog.String("component", "assignment_handler")),
	}
}

// AssignTicket handles POST /tickets/{ticketId}/assign
func (h *AssignmentHandler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "ticketId")

	if ticketID == "" {
		respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.AssignTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set ticket ID from URL
	req.TicketID = ticketID

	assignment, err := h.service.AssignTicket(ctx, req)
	if err != nil {
		h.logger.Error("Failed to assign ticket", 
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to assign ticket: "+err.Error())
		return
	}

	h.logger.Info("Ticket assigned successfully", 
		slog.String("ticket_id", ticketID),
		slog.String("assignment_id", assignment.ID))

	respondJSON(w, http.StatusCreated, assignment)
}

// EscalateTicket handles POST /tickets/{ticketId}/escalate
func (h *AssignmentHandler) EscalateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "ticketId")

	if ticketID == "" {
		respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.EscalateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set ticket ID from URL
	req.TicketID = ticketID

	assignment, err := h.service.EscalateTicket(ctx, req)
	if err != nil {
		h.logger.Error("Failed to escalate ticket",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to escalate ticket: "+err.Error())
		return
	}

	h.logger.Info("Ticket escalated successfully",
		slog.String("ticket_id", ticketID),
		slog.String("new_assignment_id", assignment.ID))

	respondJSON(w, http.StatusCreated, assignment)
}

// GetCurrentAssignment handles GET /tickets/{ticketId}/current-assignment
func (h *AssignmentHandler) GetCurrentAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "ticketId")

	if ticketID == "" {
		respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	assignment, err := h.service.GetCurrentAssignment(ctx, ticketID)
	if err != nil {
		if err == domain.ErrAssignmentNotFound {
			respondError(w, http.StatusNotFound, "No active assignment found for ticket")
			return
		}
		h.logger.Error("Failed to get current assignment",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to get current assignment")
		return
	}

	respondJSON(w, http.StatusOK, assignment)
}

// GetAssignmentHistory handles GET /tickets/{ticketId}/assignments
func (h *AssignmentHandler) GetAssignmentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "ticketId")

	if ticketID == "" {
		respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	assignments, err := h.service.GetAssignmentHistory(ctx, ticketID)
	if err != nil {
		h.logger.Error("Failed to get assignment history",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to get assignment history")
		return
	}

	respondJSON(w, http.StatusOK, assignments)
}

// AcceptAssignment handles POST /assignments/{assignmentId}/accept
func (h *AssignmentHandler) AcceptAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignmentID := chi.URLParam(r, "assignmentId")

	if assignmentID == "" {
		respondError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	var req app.AcceptAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set assignment ID from URL
	req.AssignmentID = assignmentID

	if err := h.service.AcceptAssignment(ctx, req); err != nil {
		h.logger.Error("Failed to accept assignment",
			slog.String("assignment_id", assignmentID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to accept assignment: "+err.Error())
		return
	}

	h.logger.Info("Assignment accepted",
		slog.String("assignment_id", assignmentID))

	respondJSON(w, http.StatusOK, map[string]string{"message": "Assignment accepted successfully"})
}

// RejectAssignment handles POST /assignments/{assignmentId}/reject
func (h *AssignmentHandler) RejectAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignmentID := chi.URLParam(r, "assignmentId")

	if assignmentID == "" {
		respondError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	var req app.RejectAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set assignment ID from URL
	req.AssignmentID = assignmentID

	if err := h.service.RejectAssignment(ctx, req); err != nil {
		h.logger.Error("Failed to reject assignment",
			slog.String("assignment_id", assignmentID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to reject assignment: "+err.Error())
		return
	}

	h.logger.Info("Assignment rejected",
		slog.String("assignment_id", assignmentID))

	respondJSON(w, http.StatusOK, map[string]string{"message": "Assignment rejected successfully"})
}

// StartAssignment handles POST /assignments/{assignmentId}/start
func (h *AssignmentHandler) StartAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignmentID := chi.URLParam(r, "assignmentId")

	if assignmentID == "" {
		respondError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	var req app.StartAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set assignment ID from URL
	req.AssignmentID = assignmentID

	if err := h.service.StartAssignment(ctx, req); err != nil {
		h.logger.Error("Failed to start assignment",
			slog.String("assignment_id", assignmentID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to start assignment: "+err.Error())
		return
	}

	h.logger.Info("Assignment started",
		slog.String("assignment_id", assignmentID))

	respondJSON(w, http.StatusOK, map[string]string{"message": "Assignment started successfully"})
}

// CompleteAssignment handles POST /assignments/{assignmentId}/complete
func (h *AssignmentHandler) CompleteAssignment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignmentID := chi.URLParam(r, "assignmentId")

	if assignmentID == "" {
		respondError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	var req app.CompleteAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set assignment ID from URL
	req.AssignmentID = assignmentID

	if err := h.service.CompleteAssignment(ctx, req); err != nil {
		h.logger.Error("Failed to complete assignment",
			slog.String("assignment_id", assignmentID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to complete assignment: "+err.Error())
		return
	}

	h.logger.Info("Assignment completed",
		slog.String("assignment_id", assignmentID))

	respondJSON(w, http.StatusOK, map[string]string{"message": "Assignment completed successfully"})
}

// AddCustomerFeedback handles POST /assignments/{assignmentId}/feedback
func (h *AssignmentHandler) AddCustomerFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignmentID := chi.URLParam(r, "assignmentId")

	if assignmentID == "" {
		respondError(w, http.StatusBadRequest, "Assignment ID is required")
		return
	}

	var req app.AddCustomerFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Set assignment ID from URL
	req.AssignmentID = assignmentID

	if err := h.service.AddCustomerFeedback(ctx, req); err != nil {
		h.logger.Error("Failed to add customer feedback",
			slog.String("assignment_id", assignmentID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to add feedback: "+err.Error())
		return
	}

	h.logger.Info("Customer feedback added",
		slog.String("assignment_id", assignmentID))

	respondJSON(w, http.StatusOK, map[string]string{"message": "Feedback added successfully"})
}

// GetEngineerAssignments handles GET /engineers/{engineerId}/assignments
func (h *AssignmentHandler) GetEngineerAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	engineerID := chi.URLParam(r, "engineerId")

	if engineerID == "" {
		respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}

	// Parse limit from query params
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	assignments, err := h.service.GetEngineerAssignments(ctx, engineerID, limit)
	if err != nil {
		h.logger.Error("Failed to get engineer assignments",
			slog.String("engineer_id", engineerID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to get assignments")
		return
	}

	respondJSON(w, http.StatusOK, assignments)
}

// GetActiveEngineerAssignments handles GET /engineers/{engineerId}/assignments/active
func (h *AssignmentHandler) GetActiveEngineerAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	engineerID := chi.URLParam(r, "engineerId")

	if engineerID == "" {
		respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}

	assignments, err := h.service.GetActiveEngineerAssignments(ctx, engineerID)
	if err != nil {
		h.logger.Error("Failed to get active engineer assignments",
			slog.String("engineer_id", engineerID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to get active assignments")
		return
	}

	respondJSON(w, http.StatusOK, assignments)
}

// GetEngineerWorkload handles GET /engineers/{engineerId}/workload
func (h *AssignmentHandler) GetEngineerWorkload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	engineerID := chi.URLParam(r, "engineerId")

	if engineerID == "" {
		respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}

	activeCount, completedCount, avgHours, err := h.service.GetEngineerWorkload(ctx, engineerID)
	if err != nil {
		h.logger.Error("Failed to get engineer workload",
			slog.String("engineer_id", engineerID),
			slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "Failed to get workload")
		return
	}

	workload := map[string]interface{}{
		"engineer_id":      engineerID,
		"active_count":     activeCount,
		"completed_count":  completedCount,
		"avg_hours":        avgHours,
	}

	respondJSON(w, http.StatusOK, workload)
}

// Helper functions for JSON responses
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, map[string]string{"error": message})
}
