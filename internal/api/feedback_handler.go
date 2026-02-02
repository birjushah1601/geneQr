package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/feedback"
	"github.com/gorilla/mux"
)

// FeedbackHandler handles feedback HTTP requests
type FeedbackHandler struct {
	collector *feedback.Collector
	analyzer  *feedback.Analyzer
	learner   *feedback.Learner
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(db *sql.DB) *FeedbackHandler {
	return &FeedbackHandler{
		collector: feedback.NewCollector(db),
		analyzer:  feedback.NewAnalyzer(db),
		learner:   feedback.NewLearner(db),
	}
}

// RegisterRoutes registers feedback routes
func (h *FeedbackHandler) RegisterRoutes(r *mux.Router) {
	// Human feedback endpoints
	r.HandleFunc("/api/feedback/human", h.SubmitHumanFeedback).Methods("POST")
	
	// Machine feedback endpoints (internal use)
	r.HandleFunc("/api/feedback/machine", h.SubmitMachineFeedback).Methods("POST")
	r.HandleFunc("/api/tickets/{ticketId}/auto-feedback", h.AutoCollectFeedback).Methods("POST")
	
	// Retrieval endpoints
	r.HandleFunc("/api/feedback/{serviceType}/{requestId}", h.GetFeedbackByRequest).Methods("GET")
	r.HandleFunc("/api/tickets/{ticketId}/feedback", h.GetFeedbackByTicket).Methods("GET")
	
	// Analytics endpoints
	r.HandleFunc("/api/feedback/analytics", h.GetFeedbackAnalytics).Methods("GET")
	r.HandleFunc("/api/feedback/summary", h.GetFeedbackSummary).Methods("GET")
	
	// Learning endpoints
	r.HandleFunc("/api/feedback/improvements", h.GetImprovements).Methods("GET")
	r.HandleFunc("/api/feedback/improvements/{opportunityId}/apply", h.ApplyImprovement).Methods("POST")
	r.HandleFunc("/api/feedback/actions/{actionId}/evaluate", h.EvaluateAction).Methods("POST")
	r.HandleFunc("/api/feedback/learning-progress", h.GetLearningProgress).Methods("GET")
}

// SubmitHumanFeedback handles POST /api/feedback/human
func (h *FeedbackHandler) SubmitHumanFeedback(w http.ResponseWriter, r *http.Request) {
	var req feedback.HumanFeedbackRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Validate required fields
	if req.ServiceType == "" || req.RequestID == "" || req.UserID == 0 {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Missing required fields",
			Message: "service_type, request_id, and user_id are required",
		})
		return
	}
	
	// Collect feedback
	entry, err := h.collector.CollectHumanFeedback(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to submit feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":     true,
		"feedback_id": entry.FeedbackID,
		"message":     "Thank you for your feedback! It helps us improve our AI systems.",
	})
}

// SubmitMachineFeedback handles POST /api/feedback/machine
func (h *FeedbackHandler) SubmitMachineFeedback(w http.ResponseWriter, r *http.Request) {
	var req feedback.MachineFeedbackRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Collect feedback
	entry, err := h.collector.CollectMachineFeedback(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to submit machine feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":     true,
		"feedback_id": entry.FeedbackID,
	})
}

// AutoCollectFeedback handles POST /api/tickets/{ticketId}/auto-feedback
func (h *FeedbackHandler) AutoCollectFeedback(w http.ResponseWriter, r *http.Request) {
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
	
	// Auto-collect feedback for all AI services used in the ticket
	err = h.collector.CollectTicketCompletionFeedback(r.Context(), ticketID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to auto-collect feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Automatically collected feedback from ticket outcomes",
	})
}

// GetFeedbackByRequest handles GET /api/feedback/{serviceType}/{requestId}
func (h *FeedbackHandler) GetFeedbackByRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceType := vars["serviceType"]
	requestID := vars["requestId"]
	
	entries, err := h.collector.GetFeedbackByRequest(r.Context(), serviceType, requestID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"service_type": serviceType,
		"request_id":   requestID,
		"count":        len(entries),
		"feedback":     entries,
	})
}

// GetFeedbackByTicket handles GET /api/tickets/{ticketId}/feedback
func (h *FeedbackHandler) GetFeedbackByTicket(w http.ResponseWriter, r *http.Request) {
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
	
	entries, err := h.collector.GetFeedbackByTicket(r.Context(), ticketID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id": ticketID,
		"count":     len(entries),
		"feedback":  entries,
	})
}

// GetFeedbackAnalytics handles GET /api/feedback/analytics
func (h *FeedbackHandler) GetFeedbackAnalytics(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("service_type")
	daysStr := r.URL.Query().Get("days")
	
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}
	
	if serviceType == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Missing service_type parameter",
			Message: "Provide service_type: diagnosis, assignment, or parts",
		})
		return
	}
	
	analysis, err := h.analyzer.AnalyzeFeedback(r.Context(), serviceType, days)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to analyze feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, analysis)
}

// GetFeedbackSummary handles GET /api/feedback/summary
func (h *FeedbackHandler) GetFeedbackSummary(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}
	
	summary, err := h.analyzer.GetSummary(r.Context(), days)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get feedback summary",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, summary)
}

// GetImprovements handles GET /api/feedback/improvements
func (h *FeedbackHandler) GetImprovements(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("service_type")
	status := r.URL.Query().Get("status")
	
	if serviceType == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Missing service_type parameter",
			Message: "Provide service_type: diagnosis, assignment, or parts",
		})
		return
	}
	
	// Analyze feedback to generate improvements
	analysis, err := h.analyzer.AnalyzeFeedback(r.Context(), serviceType, 30)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get improvements",
			Message: err.Error(),
		})
		return
	}
	
	// Filter by status if provided
	var filteredImprovements []feedback.ImprovementOpportunity
	for _, imp := range analysis.Improvements {
		if status == "" || imp.Status == status {
			filteredImprovements = append(filteredImprovements, imp)
		}
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"service_type":  serviceType,
		"count":         len(filteredImprovements),
		"improvements":  filteredImprovements,
	})
}

// ApplyImprovementRequest represents a request to apply an improvement
type ApplyImprovementRequest struct {
	AppliedBy string `json:"applied_by"`
}

// ApplyImprovement handles POST /api/feedback/improvements/{opportunityId}/apply
func (h *FeedbackHandler) ApplyImprovement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	opportunityID := vars["opportunityId"]
	
	var req ApplyImprovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	if req.AppliedBy == "" {
		req.AppliedBy = "system"
	}
	
	// Apply the improvement
	action, err := h.learner.ApplyImprovement(r.Context(), opportunityID, req.AppliedBy)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to apply improvement",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Improvement applied successfully and is now in testing mode",
		"action":  action,
	})
}

// EvaluateAction handles POST /api/feedback/actions/{actionId}/evaluate
func (h *FeedbackHandler) EvaluateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionID := vars["actionId"]
	
	// Evaluate the action
	err := h.learner.EvaluateAction(r.Context(), actionID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to evaluate action",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Action evaluated successfully",
	})
}

// GetLearningProgress handles GET /api/feedback/learning-progress
func (h *FeedbackHandler) GetLearningProgress(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("service_type")
	
	if serviceType == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Missing service_type parameter",
			Message: "Provide service_type: diagnosis, assignment, or parts",
		})
		return
	}
	
	progress, err := h.learner.GetLearningProgress(r.Context(), serviceType)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get learning progress",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusOK, progress)
}

