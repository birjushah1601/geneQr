package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// MockAttachmentHandler provides simple mock responses for testing
type MockAttachmentHandler struct {
	logger *slog.Logger
}

// NewMockAttachmentHandler creates a new mock attachment HTTP handler
func NewMockAttachmentHandler(logger *slog.Logger) *MockAttachmentHandler {
	return &MockAttachmentHandler{
		logger: logger.With(slog.String("component", "mock_attachment_handler")),
	}
}

// ListAttachments handles GET /api/v1/attachments
func (h *MockAttachmentHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	mockResponse := AttachmentResponse{
		Items: []AttachmentInfo{
			{
				ID:         "mock-1",
				FileName:   "test-mri-scan.jpg",
				FileSize:   2456789,
				FileType:   "image/jpeg",
				UploadDate: "2024-11-19T00:30:00Z",
				TicketID:   "TK-2025-001",
				Category:   "equipment_photo",
				Status:     "completed",
				Source:     "whatsapp",
			},
			{
				ID:         "mock-2", 
				FileName:   "ct-scanner-error.jpg",
				FileSize:   1234567,
				FileType:   "image/jpeg",
				UploadDate: "2024-11-19T00:35:00Z",
				TicketID:   "TK-2025-002",
				Category:   "issue_photo",
				Status:     "processing",
				Source:     "whatsapp",
			},
		},
		Total:    2,
		Page:     1,
		PageSize: 20,
		HasNext:  false,
		HasPrev:  false,
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    mockResponse,
	})
}

// GetAttachment handles GET /api/v1/attachments/{id}
func (h *MockAttachmentHandler) GetAttachment(w http.ResponseWriter, r *http.Request) {
	mockAttachment := AttachmentInfo{
		ID:         "mock-1",
		FileName:   "test-mri-scan.jpg",
		FileSize:   2456789,
		FileType:   "image/jpeg",
		UploadDate: "2024-11-19T00:30:00Z",
		TicketID:   "TK-2025-001",
		Category:   "equipment_photo",
		Status:     "completed",
		Source:     "whatsapp",
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    mockAttachment,
	})
}

// GetAIAnalysis handles GET /api/v1/attachments/{id}/ai-analysis
func (h *MockAttachmentHandler) GetAIAnalysis(w http.ResponseWriter, r *http.Request) {
	mockAnalysis := AIAnalysisResult{
		ID:           "ai-mock-1",
		AttachmentID: "mock-1",
		TicketID:     "TK-2025-001",
		AIProvider:   "openai",
		AIModel:      "gpt-4-vision-preview",
		Confidence:   0.87,
		ImageQualityScore: 0.92,
		AnalysisQuality:   "good",
		ProcessingDurationMs: 2340,
		TokensUsed:   1245,
		CostUsd:     0.0089,
		Status:      "completed",
		AnalyzedAt:  "2024-11-19T00:32:00Z",
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    mockAnalysis,
	})
}

// GetStats handles GET /api/v1/attachments/stats
func (h *MockAttachmentHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	mockStats := AttachmentStats{
		Total:         2,
		ByStatus:      map[string]int{"completed": 1, "processing": 1},
		ByCategory:    map[string]int{"equipment_photo": 1, "issue_photo": 1},
		BySource:      map[string]int{"whatsapp": 2},
		AvgConfidence: 0.87,
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    mockStats,
	})
}

// CreateAttachment handles POST /api/v1/attachments (mock implementation)
func (h *MockAttachmentHandler) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	// For mock, just return a successful response
	mockCreatedAttachment := AttachmentInfo{
		ID:         "mock-new-1",
		FileName:   "uploaded-file.jpg",
		FileSize:   1024000,
		FileType:   "image/jpeg",
		UploadDate: "2024-11-19T12:00:00Z",
		TicketID:   "TK-2025-003",
		Category:   "equipment_photo",
		Status:     "pending",
		Source:     "web_upload",
	}

	h.logger.Info("Mock file upload completed", 
		slog.String("attachment_id", mockCreatedAttachment.ID),
		slog.String("filename", mockCreatedAttachment.FileName))

	h.respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    mockCreatedAttachment,
	})
}

// Helper methods
func (h *MockAttachmentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *MockAttachmentHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}