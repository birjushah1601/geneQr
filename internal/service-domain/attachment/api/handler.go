package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
    "github.com/aby-med/medical-platform/internal/service-domain/attachment/infra"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// AttachmentHandler handles HTTP requests for attachment operations
type AttachmentHandler struct {
	service *domain.AttachmentService
	logger  *slog.Logger
}

// NewAttachmentHandler creates a new attachment HTTP handler
func NewAttachmentHandler(service *domain.AttachmentService, logger *slog.Logger) *AttachmentHandler {
	return &AttachmentHandler{
		service: service,
		logger:  logger.With(slog.String("component", "attachment_handler")),
	}
}

// ListAttachments handles GET /api/v1/attachments
func (h *AttachmentHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	// Parse query parameters
	page, _ := strconv.Atoi(query.Get("page"))
	pageSize, _ := strconv.Atoi(query.Get("page_size"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

    // Build filter request
	listReq := &domain.ListAttachmentsRequest{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
	
	if ticketID := query.Get("ticket_id"); ticketID != "" {
		listReq.TicketID = &ticketID
	}
	if statusStr := query.Get("status"); statusStr != "" {
		status := domain.ProcessingStatus(statusStr)
		listReq.Status = &status
	}
	if categoryStr := query.Get("category"); categoryStr != "" {
		category := domain.AttachmentCategory(categoryStr)
		listReq.Category = &category
	}
	if sourceStr := query.Get("source"); sourceStr != "" {
		source := domain.AttachmentSource(sourceStr)
		listReq.Source = &source
	}
	
    listReq.SortBy = query.Get("sort_by")
	listReq.SortOrder = query.Get("sort_direction")
    if ua := query.Get("unassigned"); ua == "true" || ua == "1" || ua == "yes" {
        listReq.UnassignedOnly = true
    }

	result, err := h.service.ListAttachments(ctx, listReq)
	if err != nil {
		h.logger.Error("Failed to list attachments", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list attachments: "+err.Error())
		return
	}

	// Convert to response format
	response := AttachmentResponse{
		Items:    make([]AttachmentInfo, len(result.Attachments)),
		Total:    int(result.Total),
		Page:     page,
		PageSize: pageSize,
		HasNext:  int64(page*pageSize) < result.Total,
		HasPrev:  page > 1,
	}

	for i, item := range result.Attachments {
        response.Items[i] = AttachmentInfo{
			ID:         item.ID.String(),
			FileName:   item.Filename,
			FileSize:   item.FileSizeBytes,
			FileType:   item.FileType,
			UploadDate: item.CreatedAt.Format("2006-01-02T15:04:05Z"),
            TicketID:   func() string { if item.TicketID==nil {return ""}; return *item.TicketID }(),
			Category:   item.AttachmentCategory,
			Status:     item.ProcessingStatus,
			Source:     item.Source,
		}
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetAttachment handles GET /api/v1/attachments/{id}
func (h *AttachmentHandler) GetAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.respondError(w, http.StatusBadRequest, "Attachment ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid attachment ID format")
		return
	}

	attachment, err := h.service.GetAttachment(ctx, id)
	if err != nil {
		if errors.Is(err, fmt.Errorf("attachment not found")) {
			h.respondError(w, http.StatusNotFound, "Attachment not found")
			return
		}
		h.logger.Error("Failed to get attachment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get attachment: "+err.Error())
		return
	}

    response := AttachmentInfo{
		ID:         attachment.ID.String(),
		FileName:   attachment.Filename,
		FileSize:   attachment.FileSizeBytes,
		FileType:   attachment.FileType,
		UploadDate: attachment.CreatedAt.Format("2006-01-02T15:04:05Z"),
        TicketID:   func() string { if attachment.TicketID==nil {return ""}; return *attachment.TicketID }(),
		Category:   attachment.AttachmentCategory,
		Status:     attachment.ProcessingStatus,
		Source:     attachment.Source,
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// DownloadAttachment handles GET /api/v1/attachments/{id}/download
func (h *AttachmentHandler) DownloadAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.respondError(w, http.StatusBadRequest, "Attachment ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid attachment ID format")
		return
	}

	// Get attachment metadata
	attachment, err := h.service.GetAttachment(ctx, id)
	if err != nil {
		if errors.Is(err, fmt.Errorf("attachment not found")) {
			h.respondError(w, http.StatusNotFound, "Attachment not found")
			return
		}
		h.logger.Error("Failed to get attachment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get attachment: "+err.Error())
		return
	}

	// Open the file
	file, err := os.Open(attachment.StoragePath)
	if err != nil {
		h.logger.Error("Failed to open attachment file", 
			slog.String("attachment_id", id.String()),
			slog.String("storage_path", attachment.StoragePath),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Attachment file not found on disk")
		return
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		h.logger.Error("Failed to stat attachment file", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to read file info")
		return
	}

	// Set headers for download
	w.Header().Set("Content-Type", attachment.FileType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.OriginalFilename))
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Stream the file to response
	_, err = io.Copy(w, file)
	if err != nil {
		h.logger.Error("Failed to stream attachment file", slog.String("error", err.Error()))
		return
	}

	h.logger.Info("Attachment downloaded successfully",
		slog.String("attachment_id", id.String()),
		slog.String("filename", attachment.OriginalFilename),
		slog.Int64("size", fileInfo.Size()))
}

// GetAIAnalysis handles GET /api/v1/attachments/{id}/ai-analysis
func (h *AttachmentHandler) GetAIAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.respondError(w, http.StatusBadRequest, "Attachment ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid attachment ID format")
		return
	}

    analysis, err := h.service.GetAIAnalysis(ctx, id)
	if err != nil {
        if errors.Is(err, infra.ErrNotImplemented) {
            h.respondError(w, http.StatusNotImplemented, "AI analysis not configured")
            return
        }
		if errors.Is(err, fmt.Errorf("analysis not found")) {
			h.respondError(w, http.StatusNotFound, "AI analysis not found")
			return
		}
		h.logger.Error("Failed to get AI analysis", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get AI analysis: "+err.Error())
		return
	}

	// Convert to response format
	response := AIAnalysisResult{
		ID:           analysis.ID.String(),
		AttachmentID: analysis.AttachmentID.String(),
		TicketID:     analysis.TicketID,
		AIProvider:   analysis.AIProvider,
		AIModel:      analysis.AIModel,
		Confidence:   analysis.AnalysisConfidence,
		ImageQualityScore: analysis.ImageQualityScore,
		AnalysisQuality:   analysis.AnalysisQuality,
		ProcessingDurationMs: int64(analysis.ProcessingDurationMS),
		TokensUsed:    analysis.TokensUsed,
		CostUsd:      analysis.CostUSD,
		Status:       analysis.Status,
		AnalyzedAt:   analysis.AnalyzedAt.Format("2006-01-02T15:04:05Z"),
		// Convert detected objects, issues, etc.
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// GetStats handles GET /api/v1/attachments/stats
func (h *AttachmentHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.service.GetAttachmentStats(ctx)
	if err != nil {
		h.logger.Error("Failed to get attachment stats", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get stats: "+err.Error())
		return
	}

	response := AttachmentStats{
		Total:         stats.Total,
		ByStatus:      stats.ByStatus,
		ByCategory:    stats.ByCategory,
		BySource:      stats.BySource,
		AvgConfidence: stats.AvgConfidence,
	}

	h.respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// Helper types for API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type AttachmentResponse struct {
	Items    []AttachmentInfo `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	HasNext  bool             `json:"has_next"`
	HasPrev  bool             `json:"has_prev"`
}

type AttachmentInfo struct {
	ID         string `json:"id"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	FileType   string `json:"fileType"`
	UploadDate string `json:"uploadDate"`
	TicketID   string `json:"ticketId"`
	Category   string `json:"category"`
	Status     string `json:"status"`
	Source     string `json:"source"`
}

type AIAnalysisResult struct {
	ID           string `json:"id"`
	AttachmentID string `json:"attachmentId"`
	TicketID     string `json:"ticketId"`
	AIProvider   string `json:"aiProvider"`
	AIModel      string `json:"aiModel"`
	Confidence   float64 `json:"confidence"`
	ImageQualityScore  float64 `json:"imageQualityScore"`
	AnalysisQuality    string  `json:"analysisQuality"`
	ProcessingDurationMs int64 `json:"processingDurationMs"`
	TokensUsed    int    `json:"tokensUsed"`
	CostUsd      float64 `json:"costUsd"`
	Status       string  `json:"status"`
	AnalyzedAt   string  `json:"analyzedAt"`
	// Additional fields would include detected objects, issues, etc.
}

type AttachmentStats struct {
	Total         int                `json:"total"`
	ByStatus      map[string]int     `json:"by_status"`
	ByCategory    map[string]int     `json:"by_category"`
	BySource      map[string]int     `json:"by_source"`
	AvgConfidence float64            `json:"avg_confidence"`
}

// CreateAttachment handles POST /api/v1/attachments
func (h *AttachmentHandler) CreateAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Parse multipart form with max memory 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error("Failed to parse multipart form", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, "Invalid multipart form")
		return
	}
	
	// Get the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get uploaded file", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()
	
	// Validate file type and size
	if err := h.validateFileUpload(header); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	// Get required form fields
	ticketID := r.FormValue("ticket_id")
	category := r.FormValue("category")
	source := r.FormValue("source")
	
    if category == "" || source == "" {
        h.respondError(w, http.StatusBadRequest, "category and source are required")
		return
	}
	
	// Create unique storage path
    // If ticketID is empty, store under unassigned bucket
    baseTicket := ticketID
    if baseTicket == "" {
        baseTicket = "unassigned"
    }
    storagePath, err := h.createStoragePath(header.Filename, baseTicket)
	if err != nil {
		h.logger.Error("Failed to create storage path", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to process file")
		return
	}
	
	// Save file to storage
	if err := h.saveFile(file, storagePath); err != nil {
		h.logger.Error("Failed to save file", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	
	// Create attachment request
    // Build request; TicketID optional
    var ticketPtr *string
    if ticketID != "" { ticketPtr = &ticketID }
    req := &domain.CreateAttachmentRequest{
		TicketID:         ticketPtr,
		Filename:         header.Filename,
		OriginalFilename: header.Filename,
		FileType:         header.Header.Get("Content-Type"),
		FileSizeBytes:    header.Size,
		StoragePath:      storagePath,
		Source:           domain.AttachmentSource(source),
		Category:         domain.AttachmentCategory(category),
	}
	
	// Optional fields
	if uploadedBy := r.FormValue("uploaded_by"); uploadedBy != "" {
		req.UploadedByID = &uploadedBy
	}
	if sourceMessageID := r.FormValue("source_message_id"); sourceMessageID != "" {
		req.SourceMessageID = &sourceMessageID
	}
	
	// Create attachment
	attachment, err := h.service.CreateAttachment(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create attachment", 
			slog.String("error", err.Error()),
			slog.String("ticket_id", ticketID),
			slog.String("filename", header.Filename))
		h.respondError(w, http.StatusInternalServerError, "Failed to create attachment")
		return
	}
	
	// Convert to API response format
    apiAttachment := AttachmentInfo{
		ID:         attachment.ID.String(),
		FileName:   attachment.Filename,
		FileSize:   attachment.FileSizeBytes,
		FileType:   attachment.FileType,
		UploadDate: attachment.UploadedAt.Format(time.RFC3339),
		TicketID:   func() string { if attachment.TicketID==nil {return ""}; return *attachment.TicketID }(),
		Category:   attachment.AttachmentCategory,
		Status:     attachment.ProcessingStatus,
		Source:     attachment.Source,
	}
	
	h.logger.Info("Attachment created successfully",
		slog.String("attachment_id", attachment.ID.String()),
		slog.String("ticket_id", func() string { if attachment.TicketID==nil {return ""}; return *attachment.TicketID }()),
		slog.String("filename", header.Filename))
	
	h.respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    apiAttachment,
	})
}

// Helper methods
func (h *AttachmentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *AttachmentHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

// validateFileUpload validates uploaded file type and size
func (h *AttachmentHandler) validateFileUpload(header *multipart.FileHeader) error {
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	
	if header.Size > maxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", header.Size, maxFileSize)
	}
	
	// Get file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	
	// Define allowed file types
	allowedTypes := map[string]bool{
		".jpg":  true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, // Images
		".pdf":  true, ".doc": true, ".docx": true, ".txt": true, ".rtf": true, // Documents
		".mp4":  true, ".avi": true, ".mov": true, ".wmv": true, ".mkv": true, // Videos
		".mp3":  true, ".wav": true, ".m4a": true, ".aac": true, // Audio
		".zip":  true, ".rar": true, ".7z": true, // Archives
		".xml":  true, ".json": true, ".csv": true, // Data files
	}
	
	if !allowedTypes[ext] {
		return fmt.Errorf("file type %s is not allowed", ext)
	}
	
	return nil
}

// createStoragePath generates a unique storage path for the file
func (h *AttachmentHandler) createStoragePath(filename, ticketID string) (string, error) {
	// Create base storage directory
	baseDir := filepath.Join("storage", "attachments", ticketID)
	
	// Generate unique filename to avoid conflicts
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	uniqueID := uuid.New().String()
	uniqueFilename := fmt.Sprintf("%s_%s%s", baseName, uniqueID, ext)
	
	// Full storage path
	storagePath := filepath.Join(baseDir, uniqueFilename)
	
	return storagePath, nil
}

// LinkAttachment handles POST /api/v1/attachments/{id}/link
// Body: { "ticket_id": "..." }
func (h *AttachmentHandler) LinkAttachment(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    idStr := chi.URLParam(r, "id")
    if idStr == "" { h.respondError(w, http.StatusBadRequest, "Attachment ID is required"); return }
    id, err := uuid.Parse(idStr)
    if err != nil { h.respondError(w, http.StatusBadRequest, "Invalid attachment ID format"); return }

    var body struct{ TicketID string `json:"ticket_id"` }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid JSON body")
        return
    }
    if strings.TrimSpace(body.TicketID) == "" {
        h.respondError(w, http.StatusBadRequest, "ticket_id is required")
        return
    }

    if err := h.service.LinkAttachmentToTicket(ctx, id, body.TicketID); err != nil {
        h.logger.Error("Failed to link attachment", slog.String("error", err.Error()))
        h.respondError(w, http.StatusInternalServerError, "Failed to link attachment")
        return
    }

    h.respondJSON(w, http.StatusOK, APIResponse{ Success: true })
}

// saveFile saves the uploaded file to the specified path
func (h *AttachmentHandler) saveFile(file multipart.File, storagePath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	// Create destination file
	dst, err := os.Create(storagePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", storagePath, err)
	}
	defer dst.Close()
	
	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		// Clean up partial file on error
		os.Remove(storagePath)
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	
	return nil
}

// DeleteAttachment handles DELETE /api/v1/attachments/{id}
func (h *AttachmentHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	// Parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid attachment ID")
		return
	}

	// Delete attachment
	if err := h.service.DeleteAttachment(ctx, id); err != nil {
		h.logger.Error("Failed to delete attachment",
			slog.String("id", idStr),
			slog.String("error", err.Error()),
		)
		
		if strings.Contains(err.Error(), "not found") {
			h.respondError(w, http.StatusNotFound, "Attachment not found")
			return
		}
		
		h.respondError(w, http.StatusInternalServerError, "Failed to delete attachment")
		return
	}

	// Return success response
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Attachment deleted successfully",
	})
}
