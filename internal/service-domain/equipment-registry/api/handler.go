package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	"github.com/go-chi/chi/v5"
)

// EquipmentHandler handles HTTP requests for equipment operations
type EquipmentHandler struct {
	service *app.EquipmentService
	logger  *slog.Logger
}

// NewEquipmentHandler creates a new equipment HTTP handler
func NewEquipmentHandler(service *app.EquipmentService, logger *slog.Logger) *EquipmentHandler {
	return &EquipmentHandler{
		service: service,
		logger:  logger.With(slog.String("component", "equipment_handler")),
	}
}

// RegisterEquipment handles POST /equipment
func (h *EquipmentHandler) RegisterEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req app.RegisterEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	equipment, err := h.service.RegisterEquipment(ctx, req)
	if err != nil {
		h.logger.Error("Failed to register equipment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to register equipment: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, equipment)
}

// GetEquipment handles GET /equipment/{id}
func (h *EquipmentHandler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	equipment, err := h.service.GetEquipmentByID(ctx, id)
	if err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to get equipment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get equipment")
		return
	}

	h.respondJSON(w, http.StatusOK, equipment)
}

// GetEquipmentByQR handles GET /equipment/qr/{qr_code}
func (h *EquipmentHandler) GetEquipmentByQR(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	qrCode := chi.URLParam(r, "qr_code")

	if qrCode == "" {
		h.respondError(w, http.StatusBadRequest, "QR code is required")
		return
	}

	equipment, err := h.service.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to get equipment by QR", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get equipment")
		return
	}

	h.respondJSON(w, http.StatusOK, equipment)
}

// GetEquipmentBySerial handles GET /equipment/serial/{serial}
func (h *EquipmentHandler) GetEquipmentBySerial(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serial := chi.URLParam(r, "serial")

	if serial == "" {
		h.respondError(w, http.StatusBadRequest, "Serial number is required")
		return
	}

	equipment, err := h.service.GetEquipmentBySerial(ctx, serial)
	if err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to get equipment by serial", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get equipment")
		return
	}

	h.respondJSON(w, http.StatusOK, equipment)
}

// ListEquipment handles GET /equipment
func (h *EquipmentHandler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	criteria := domain.ListCriteria{
		CustomerID:       r.URL.Query().Get("customer_id"),
		ManufacturerName: r.URL.Query().Get("manufacturer"),
		Category:         r.URL.Query().Get("category"),
		SortBy:           r.URL.Query().Get("sort_by"),
		SortDirection:    r.URL.Query().Get("sort_dir"),
	}

	// Parse status filter
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		criteria.Status = []domain.EquipmentStatus{domain.EquipmentStatus(statusStr)}
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
	if hasAMC := r.URL.Query().Get("has_amc"); hasAMC != "" {
		val := hasAMC == "true"
		criteria.HasAMC = &val
	}

	if underWarranty := r.URL.Query().Get("under_warranty"); underWarranty != "" {
		val := underWarranty == "true"
		criteria.UnderWarranty = &val
	}

	result, err := h.service.ListEquipment(ctx, criteria)
	if err != nil {
		h.logger.Error("Failed to list equipment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list equipment")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// UpdateEquipment handles PATCH /equipment/{id}
func (h *EquipmentHandler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	var equipment domain.Equipment
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	equipment.ID = id

	if err := h.service.UpdateEquipment(ctx, &equipment); err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to update equipment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update equipment")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Equipment updated successfully"})
}

// GenerateQRCode handles POST /equipment/{id}/qr
func (h *EquipmentHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	qrPath, err := h.service.GenerateQRCode(ctx, id)
	if err != nil {
		h.logger.Error("Failed to generate QR code", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to generate QR code: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"message": "QR code generated successfully",
		"path":    qrPath,
	})
}

// GetQRCodeImage handles GET /equipment/{id}/qr/image
func (h *EquipmentHandler) GetQRCodeImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	// Get equipment with QR image
	equipment, err := h.service.GetEquipmentByID(ctx, id)
	if err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to get equipment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get equipment")
		return
	}

	if len(equipment.QRCodeImage) == 0 {
		h.respondError(w, http.StatusNotFound, "QR code not generated yet")
		return
	}

	// Serve QR image from database
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(equipment.QRCodeImage)))
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 1 day
	w.WriteHeader(http.StatusOK)
	w.Write(equipment.QRCodeImage)
}
// DownloadQRLabel handles GET /equipment/{id}/qr/pdf
func (h *EquipmentHandler) DownloadQRLabel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	pdfBytes, err := h.service.GenerateQRLabel(ctx, id)
	if err != nil {
		h.logger.Error("Failed to generate QR label", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to generate QR label: "+err.Error())
		return
	}

	// Serve the PDF from database bytes
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=qr_label_"+id+".pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

// ImportCSV handles POST /equipment/import
func (h *EquipmentHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("csv_file")
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "CSV file is required")
		return
	}
	defer file.Close()

	createdBy := r.FormValue("created_by")
	if createdBy == "" {
		createdBy = "system"
	}

	// Save uploaded file temporarily
	tempFilePath := "/tmp/" + header.Filename

	// Create temp file
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFilePath)

	// Copy uploaded file to temp location
	_, err = io.Copy(tempFile, file)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Import from CSV
	result, err := h.service.BulkImportFromCSV(ctx, tempFilePath, createdBy)
	if err != nil {
		h.logger.Error("Failed to import CSV", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to import CSV: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// RecordService handles POST /equipment/{id}/service
func (h *EquipmentHandler) RecordService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}

	var req struct {
		ServiceDate time.Time `json:"service_date"`
		Notes       string    `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.RecordService(ctx, id, req.ServiceDate, req.Notes); err != nil {
		if err == domain.ErrEquipmentNotFound {
			h.respondError(w, http.StatusNotFound, "Equipment not found")
			return
		}
		h.logger.Error("Failed to record service", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to record service")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Service recorded successfully"})
}

// BulkGenerateQRCodes handles POST /equipment/qr/bulk-generate
// Generates QR codes for all equipment that doesn't have one
func (h *EquipmentHandler) BulkGenerateQRCodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger.Info("Starting bulk QR code generation")

	result, err := h.service.BulkGenerateQRCodes(ctx)
	if err != nil {
		h.logger.Error("Failed to bulk generate QR codes", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to bulk generate QR codes: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// respondJSON writes JSON response
func (h *EquipmentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func (h *EquipmentHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

