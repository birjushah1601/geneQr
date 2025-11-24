package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/assignment"
	"github.com/gorilla/mux"
)

// AssignmentHandler handles assignment-related HTTP requests
type AssignmentHandler struct {
	engine *assignment.Engine
	db     *sql.DB
}

// NewAssignmentHandler creates a new assignment handler
func NewAssignmentHandler(aiManager *ai.Manager, db *sql.DB) *AssignmentHandler {
	return &AssignmentHandler{
		engine: assignment.NewEngine(aiManager, db),
		db:     db,
	}
}

// RegisterRoutes registers assignment routes
func (h *AssignmentHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/assign", h.RecommendEngineers).Methods("POST")
	r.HandleFunc("/api/tickets/{ticketId}/assignments", h.GetTicketAssignments).Methods("GET")
	r.HandleFunc("/api/assignments/{requestId}", h.GetAssignment).Methods("GET")
	r.HandleFunc("/api/assignments/{requestId}/select", h.SelectEngineer).Methods("POST")
	r.HandleFunc("/api/assignments/{requestId}/feedback", h.ProvideFeedback).Methods("POST")
	r.HandleFunc("/api/assignments/analytics", h.GetAnalytics).Methods("GET")
}

// RecommendEngineers handles POST /api/assign
func (h *AssignmentHandler) RecommendEngineers(w http.ResponseWriter, r *http.Request) {
	var req assignment.AssignmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Set default options if not provided
	if req.Options == (assignment.AssignmentOptions{}) {
		req.Options = assignment.DefaultAssignmentOptions()
	}

	// Adjust weights based on priority
	if req.Priority == "Critical" || req.Priority == "High" {
		req.Options.Weights = assignment.UrgentScoringWeights()
	}

	// Get recommendations
	result, err := h.engine.RecommendEngineers(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Assignment recommendation failed",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetAssignment handles GET /api/assignments/{requestId}
func (h *AssignmentHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var result assignment.AssignmentResponse

	query := `
		SELECT 
			request_id,
			ticket_id,
			recommendations,
			metadata,
			created_at
		FROM assignment_history
		WHERE request_id = $1
	`

	var recsJSON, metadataJSON []byte
	err := h.db.QueryRowContext(r.Context(), query, requestID).Scan(
		&result.RequestID,
		&result.TicketID,
		&recsJSON,
		&metadataJSON,
		&result.CreatedAt,
	)

	if err == sql.ErrNoRows {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Assignment not found",
			Message: "No assignment with this request ID",
		})
		return
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve assignment",
			Message: err.Error(),
		})
		return
	}

	json.Unmarshal(recsJSON, &result.Recommendations)
	json.Unmarshal(metadataJSON, &result.Metadata)

	respondJSON(w, http.StatusOK, result)
}

// GetTicketAssignments handles GET /api/tickets/{ticketId}/assignments
func (h *AssignmentHandler) GetTicketAssignments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketIDStr := vars["ticketId"]

	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ticket ID",
			Message: err.Error(),
		})
		return
	}

	history, err := h.engine.GetAssignmentHistory(r.Context(), ticketID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve assignments",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id":    ticketID,
		"count":        len(history),
		"assignments":  history,
	})
}

// SelectEngineerRequest represents engineer selection
type SelectEngineerRequest struct {
	EngineerID     int64  `json:"engineer_id"`
	SelectionReason string `json:"selection_reason"`
}

// SelectEngineer handles POST /api/assignments/{requestId}/select
func (h *AssignmentHandler) SelectEngineer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var req SelectEngineerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	query := `
		UPDATE assignment_history
		SET 
			selected_engineer_id = $1,
			selection_reason = $2,
			selection_time = NOW(),
			updated_at = NOW()
		WHERE request_id = $3
	`

	result, err := h.db.ExecContext(
		r.Context(),
		query,
		req.EngineerID,
		req.SelectionReason,
		requestID,
	)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to record selection",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Assignment not found",
			Message: "No assignment with this request ID",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Engineer selection recorded",
	})
}

// AssignmentFeedbackRequest represents feedback on assignment
type AssignmentFeedbackRequest struct {
	WasSuccessful          bool   `json:"was_successful"`
	ActualResolutionHours  float64 `json:"actual_resolution_hours"`
	FeedbackNotes          string `json:"feedback_notes"`
}

// ProvideFeedback handles POST /api/assignments/{requestId}/feedback
func (h *AssignmentHandler) ProvideFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["requestId"]

	var req AssignmentFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	query := `
		UPDATE assignment_history
		SET 
			was_successful = $1,
			actual_resolution_time = make_interval(secs => $2 * 3600),
			feedback_notes = $3,
			updated_at = NOW()
		WHERE request_id = $4
	`

	result, err := h.db.ExecContext(
		r.Context(),
		query,
		req.WasSuccessful,
		req.ActualResolutionHours,
		req.FeedbackNotes,
		requestID,
	)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save feedback",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Assignment not found",
			Message: "No assignment with this request ID",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Feedback recorded successfully",
	})
}

// GetAnalytics handles GET /api/assignments/analytics
func (h *AssignmentHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	daysStr := r.URL.Query().Get("days")
	days := 30 // Default to last 30 days
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// Query analytics view
	query := `
		SELECT 
			assignment_date,
			total_assignments,
			assignments_made,
			successful_assignments,
			ROUND(success_rate, 2) as success_rate,
			ROUND(avg_resolution_hours, 2) as avg_resolution_hours,
			ROUND(total_ai_cost, 4) as total_ai_cost,
			ai_assisted_count
		FROM v_assignment_analytics
		WHERE assignment_date > NOW() - INTERVAL '$1 days'
		ORDER BY assignment_date DESC
	`

	rows, err := h.db.QueryContext(r.Context(), query, days)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve analytics",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var analytics []map[string]interface{}
	for rows.Next() {
		var (
			assignmentDate        string
			totalAssignments      int
			assignmentsMade       int
			successfulAssignments int
			successRate           *float64
			avgResolutionHours    *float64
			totalAICost           *float64
			aiAssistedCount       int
		)

		err := rows.Scan(
			&assignmentDate,
			&totalAssignments,
			&assignmentsMade,
			&successfulAssignments,
			&successRate,
			&avgResolutionHours,
			&totalAICost,
			&aiAssistedCount,
		)
		if err != nil {
			continue
		}

		analytics = append(analytics, map[string]interface{}{
			"date":                  assignmentDate,
			"total_assignments":     totalAssignments,
			"assignments_made":      assignmentsMade,
			"successful_assignments": successfulAssignments,
			"success_rate":          successRate,
			"avg_resolution_hours":  avgResolutionHours,
			"total_ai_cost":         totalAICost,
			"ai_assisted_count":     aiAssistedCount,
		})
	}

	// Get summary stats
	summaryQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(selected_engineer_id) as assigned,
			AVG(CASE WHEN was_successful = true THEN 100.0 ELSE 0.0 END) as success_rate,
			SUM((metadata->>'cost_usd')::DECIMAL) as total_cost
		FROM assignment_history
		WHERE created_at > NOW() - INTERVAL '$1 days'
	`

	var (
		total       int
		assigned    int
		successRate *float64
		totalCost   *float64
	)

	err = h.db.QueryRowContext(r.Context(), summaryQuery, days).Scan(
		&total,
		&assigned,
		&successRate,
		&totalCost,
	)

	if err != nil && err != sql.ErrNoRows {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve summary",
			Message: err.Error(),
		})
		return
	}

	// Get top engineers
	engineersQuery := `
		SELECT 
			user_id,
			full_name,
			total_assignments,
			successful_count,
			success_percentage,
			avg_resolution_hours
		FROM v_engineer_assignment_success
		ORDER BY success_percentage DESC, total_assignments DESC
		LIMIT 10
	`

	engineerRows, err := h.db.QueryContext(r.Context(), engineersQuery)
	if err == nil {
		defer engineerRows.Close()
		var topEngineers []map[string]interface{}

		for engineerRows.Next() {
			var (
				userID               int64
				fullName             string
				totalAssignments     int
				successfulCount      int
				successPercentage    float64
				avgResolutionHours   *float64
			)

			engineerRows.Scan(
				&userID,
				&fullName,
				&totalAssignments,
				&successfulCount,
				&successPercentage,
				&avgResolutionHours,
			)

			topEngineers = append(topEngineers, map[string]interface{}{
				"engineer_id":          userID,
				"name":                 fullName,
				"total_assignments":    totalAssignments,
				"successful_count":     successfulCount,
				"success_percentage":   successPercentage,
				"avg_resolution_hours": avgResolutionHours,
			})
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"period_days": days,
			"summary": map[string]interface{}{
				"total_recommendations": total,
				"engineers_assigned":    assigned,
				"success_rate":          successRate,
				"total_ai_cost":         totalCost,
			},
			"daily_analytics":  analytics,
			"top_engineers":    topEngineers,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"period_days": days,
		"summary": map[string]interface{}{
			"total_recommendations": total,
			"engineers_assigned":    assigned,
			"success_rate":          successRate,
			"total_ai_cost":         totalCost,
		},
		"daily_analytics": analytics,
	})
}

