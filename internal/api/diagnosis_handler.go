package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/aby-med/medical-platform/internal/diagnosis"
	"github.com/gorilla/mux"
)

// DiagnosisHandler handles diagnosis-related HTTP requests
type DiagnosisHandler struct {
	engine *diagnosis.Engine
	db     *sql.DB
}

// NewDiagnosisHandler creates a new diagnosis handler
func NewDiagnosisHandler(aiManager *ai.Manager, db *sql.DB) *DiagnosisHandler {
	return &DiagnosisHandler{
		engine: diagnosis.NewEngine(aiManager, db),
		db:     db,
	}
}

// RegisterRoutes registers diagnosis routes
func (h *DiagnosisHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/diagnose", h.DiagnoseTicket).Methods("POST")
	r.HandleFunc("/api/diagnoses/{diagnosisId}", h.GetDiagnosis).Methods("GET")
	r.HandleFunc("/api/tickets/{ticketId}/diagnoses", h.GetTicketDiagnoses).Methods("GET")
	r.HandleFunc("/api/diagnoses/{diagnosisId}/feedback", h.ProvideFeedback).Methods("POST")
	r.HandleFunc("/api/diagnoses/analytics", h.GetAnalytics).Methods("GET")
}

// DiagnoseTicket handles POST /api/diagnose
func (h *DiagnosisHandler) DiagnoseTicket(w http.ResponseWriter, r *http.Request) {
	var req diagnosis.DiagnosisRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Set default options if not provided
	if req.Options == (diagnosis.DiagnosisOptions{}) {
		req.Options = diagnosis.DefaultDiagnosisOptions()
	}

	// Perform diagnosis
	result, err := h.engine.Diagnose(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Diagnosis failed",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetDiagnosis handles GET /api/diagnoses/{diagnosisId}
func (h *DiagnosisHandler) GetDiagnosis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	diagnosisID := vars["diagnosisId"]

	result, err := h.engine.GetDiagnosis(r.Context(), diagnosisID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, ErrorResponse{
			Error:   "Diagnosis not found",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetTicketDiagnoses handles GET /api/tickets/{ticketId}/diagnoses
func (h *DiagnosisHandler) GetTicketDiagnoses(w http.ResponseWriter, r *http.Request) {
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

	diagnoses, err := h.engine.GetTicketDiagnoses(r.Context(), ticketID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve diagnoses",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id":  ticketID,
		"count":      len(diagnoses),
		"diagnoses":  diagnoses,
	})
}

// FeedbackRequest represents feedback on a diagnosis
type FeedbackRequest struct {
	WasAccurate        bool   `json:"was_accurate"`
	AccuracyScore      int    `json:"accuracy_score"`       // 0-100
	FeedbackNotes      string `json:"feedback_notes"`
	ActualResolution   string `json:"actual_resolution"`
}

// ProvideFeedback handles POST /api/diagnoses/{diagnosisId}/feedback
func (h *DiagnosisHandler) ProvideFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	diagnosisID := vars["diagnosisId"]

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Validate accuracy score
	if req.AccuracyScore < 0 || req.AccuracyScore > 100 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid accuracy score",
			Message: "Accuracy score must be between 0 and 100",
		})
		return
	}

	err := h.engine.ProvideFeedback(
		r.Context(),
		diagnosisID,
		req.WasAccurate,
		req.AccuracyScore,
		req.FeedbackNotes,
		req.ActualResolution,
	)

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save feedback",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Feedback recorded successfully",
	})
}

// GetAnalytics handles GET /api/diagnoses/analytics
func (h *DiagnosisHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
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
			diagnosis_date,
			provider,
			model,
			total_diagnoses,
			ROUND(avg_confidence, 2) as avg_confidence,
			ROUND(accuracy_rate, 2) as accuracy_rate,
			ROUND(avg_accuracy_score, 2) as avg_accuracy_score,
			total_tokens,
			ROUND(total_cost, 4) as total_cost,
			accurate_count,
			inaccurate_count,
			pending_feedback
		FROM v_diagnosis_analytics
		WHERE diagnosis_date > NOW() - INTERVAL '$1 days'
		ORDER BY diagnosis_date DESC
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
			diagnosisDate    string
			provider         string
			model            string
			totalDiagnoses   int
			avgConfidence    *float64
			accuracyRate     *float64
			avgAccuracyScore *float64
			totalTokens      *int
			totalCost        *float64
			accurateCount    int
			inaccurateCount  int
			pendingFeedback  int
		)

		err := rows.Scan(
			&diagnosisDate,
			&provider,
			&model,
			&totalDiagnoses,
			&avgConfidence,
			&accuracyRate,
			&avgAccuracyScore,
			&totalTokens,
			&totalCost,
			&accurateCount,
			&inaccurateCount,
			&pendingFeedback,
		)
		if err != nil {
			continue
		}

		analytics = append(analytics, map[string]interface{}{
			"date":               diagnosisDate,
			"provider":           provider,
			"model":              model,
			"total_diagnoses":    totalDiagnoses,
			"avg_confidence":     avgConfidence,
			"accuracy_rate":      accuracyRate,
			"avg_accuracy_score": avgAccuracyScore,
			"total_tokens":       totalTokens,
			"total_cost":         totalCost,
			"accurate_count":     accurateCount,
			"inaccurate_count":   inaccurateCount,
			"pending_feedback":   pendingFeedback,
		})
	}

	// Get summary stats
	summaryQuery := `
		SELECT 
			COUNT(*) as total,
			SUM(cost_usd) as total_cost,
			AVG(confidence_score) as avg_confidence,
			COUNT(CASE WHEN was_accurate = true THEN 1 END) as accurate,
			COUNT(CASE WHEN was_accurate = false THEN 1 END) as inaccurate
		FROM ai_diagnoses
		WHERE created_at > NOW() - INTERVAL '$1 days'
	`

	var (
		total          int
		totalCost      *float64
		avgConfidence  *float64
		accurateCount  int
		inaccurateCount int
	)

	err = h.db.QueryRowContext(r.Context(), summaryQuery, days).Scan(
		&total,
		&totalCost,
		&avgConfidence,
		&accurateCount,
		&inaccurateCount,
	)

	if err != nil && err != sql.ErrNoRows {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve summary",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"period_days": days,
		"summary": map[string]interface{}{
			"total_diagnoses":   total,
			"total_cost":        totalCost,
			"avg_confidence":    avgConfidence,
			"accurate_count":    accurateCount,
			"inaccurate_count":  inaccurateCount,
			"accuracy_rate":     float64(accurateCount) / float64(accurateCount+inaccurateCount) * 100,
		},
		"daily_analytics": analytics,
	})
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteStatus(status)
	json.NewEncoder(w).Encode(data)
}

