package diagnosis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// FeedbackHandler handles AI diagnosis feedback endpoints
type FeedbackHandler struct {
	engine *Engine
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(engine *Engine) *FeedbackHandler {
	return &FeedbackHandler{
		engine: engine,
	}
}

// RegisterRoutes registers feedback routes
func (h *FeedbackHandler) RegisterRoutes(r chi.Router) {
	r.Route("/diagnosis", func(r chi.Router) {
		r.Post("/analyze", h.AnalyzeDiagnosis)         // Enhanced with confidence
		r.Post("/{diagnosisId}/feedback", h.SubmitFeedback) // NEW: Submit decision feedback
		r.Get("/{diagnosisId}", h.GetDiagnosis)          // Get diagnosis with decision status
	})
}

// AnalyzeDiagnosis runs AI diagnosis with confidence scoring
func (h *FeedbackHandler) AnalyzeDiagnosis(w http.ResponseWriter, r *http.Request) {
	var req DiagnosisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
		return
	}

	// Run AI diagnosis
	response, err := h.engine.Diagnose(r.Context(), &req)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to diagnose: %v", err),
		})
		return
	}

	// Enhance with confidence scoring
	enhanced, err := h.engine.EnhanceResponseWithConfidence(response)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to calculate confidence: %v", err),
		})
		return
	}

	// TODO: Store enhanced response in database
	// For now, just return it
	respondJSON(w, http.StatusOK, enhanced)
}

// SubmitFeedback handles user decision feedback (accept/reject)
func (h *FeedbackHandler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	diagnosisID := chi.URLParam(r, "diagnosisId")
	if diagnosisID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Missing diagnosis ID",
		})
		return
	}

	var feedback AIDecisionFeedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid feedback format",
		})
		return
	}

	// Validate decision
	if feedback.Decision != "accepted" && feedback.Decision != "rejected" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Decision must be 'accepted' or 'rejected'",
		})
		return
	}

	// Set diagnosis ID from URL
	feedback.DiagnosisID = diagnosisID

	// Store feedback in database
	err := h.storeFeedback(&feedback)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to store feedback: %v", err),
		})
		return
	}

	var message string
	if feedback.Decision == "accepted" {
		message = "Thank you! Your acceptance helps our AI learn what works well."
	} else {
		message = "Thank you for the feedback! This helps our AI improve for next time."
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":     true,
		"message":     message,
		"feedback_id": fmt.Sprintf("fb_%s_%d", diagnosisID, time.Now().Unix()),
	})
}

// GetDiagnosis retrieves a diagnosis with current decision status
func (h *FeedbackHandler) GetDiagnosis(w http.ResponseWriter, r *http.Request) {
	diagnosisID := chi.URLParam(r, "diagnosisId")
	if diagnosisID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Missing diagnosis ID",
		})
		return
	}

	// TODO: Retrieve from database
	// For now, return mock data
	respondJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "Get diagnosis not yet implemented - database integration needed",
	})
}

// storeFeedback stores user feedback in database
func (h *FeedbackHandler) storeFeedback(feedback *AIDecisionFeedback) error {
	// TODO: Implement database storage
	// This would typically:
	// 1. Update diagnosis record with decision status
	// 2. Store feedback in ai_feedback table
	// 3. Trigger learning pipeline if needed
	
	// For now, just log it
	fmt.Printf("[FEEDBACK] DiagnosisID: %s, Decision: %s, UserID: %d, Role: %s\n",
		feedback.DiagnosisID, feedback.Decision, feedback.UserID, feedback.UserRole)
	
	if feedback.FeedbackText != "" {
		fmt.Printf("[FEEDBACK] Text: %s\n", feedback.FeedbackText)
	}
	
	if len(feedback.Corrections) > 0 {
		fmt.Printf("[FEEDBACK] Corrections: %+v\n", feedback.Corrections)
	}

	return nil
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}