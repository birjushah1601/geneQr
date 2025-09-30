package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aby-med/medical-platform/internal/service-domain/supplier/app"
	"github.com/aby-med/medical-platform/internal/service-domain/supplier/domain"
	"github.com/go-chi/chi/v5"
)

// SupplierHandler handles HTTP requests for supplier operations
type SupplierHandler struct {
	service *app.SupplierService
	logger  *slog.Logger
}

// NewSupplierHandler creates a new supplier handler
func NewSupplierHandler(service *app.SupplierService, logger *slog.Logger) *SupplierHandler {
	return &SupplierHandler{
		service: service,
		logger:  logger.With(slog.String("component", "supplier_handler")),
	}
}

// RegisterRoutes registers all supplier routes
func (h *SupplierHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateSupplier)
	r.Get("/", h.ListSuppliers)
	r.Get("/{id}", h.GetSupplier)
	r.Put("/{id}", h.UpdateSupplier)
	r.Delete("/{id}", h.DeleteSupplier)
	
	// Lifecycle operations
	r.Post("/{id}/verify", h.VerifySupplier)
	r.Post("/{id}/reject", h.RejectSupplier)
	r.Post("/{id}/suspend", h.SuspendSupplier)
	r.Post("/{id}/activate", h.ActivateSupplier)
	
	// Additional operations
	r.Post("/{id}/certifications", h.AddCertification)
	r.Get("/category/{categoryId}", h.GetSuppliersByCategory)
}

// CreateSupplier handles supplier creation
func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	var req app.CreateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Get user from auth context
	createdBy := "system"

	response, err := h.service.CreateSupplier(r.Context(), tenantID, req, createdBy)
	if err != nil {
		if err == domain.ErrSupplierAlreadyExists {
			h.respondError(w, http.StatusConflict, "Supplier with this tax ID already exists")
			return
		}
		h.logger.Error("Failed to create supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to create supplier")
		return
	}

	h.respondJSON(w, http.StatusCreated, response)
}

// GetSupplier handles retrieving a supplier by ID
func (h *SupplierHandler) GetSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	response, err := h.service.GetSupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to get supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// UpdateSupplier handles updating a supplier
func (h *SupplierHandler) UpdateSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	var req app.UpdateSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.UpdateSupplier(r.Context(), tenantID, supplierID, req)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		if err == domain.ErrCannotModifySupplier {
			h.respondError(w, http.StatusForbidden, "Supplier cannot be modified in current status")
			return
		}
		h.logger.Error("Failed to update supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// DeleteSupplier handles deleting a supplier
func (h *SupplierHandler) DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	err := h.service.DeleteSupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to delete supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to delete supplier")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSuppliers handles listing suppliers with filters and pagination
func (h *SupplierHandler) ListSuppliers(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	// Parse query parameters
	req := app.ListSuppliersRequest{
		Status:             r.URL.Query()["status"],
		VerificationStatus: r.URL.Query()["verification_status"],
		CategoryID:         r.URL.Query().Get("category_id"),
		SearchQuery:        r.URL.Query().Get("search"),
		SortBy:             r.URL.Query().Get("sort_by"),
		SortDirection:      r.URL.Query().Get("sort_direction"),
	}

	// Parse page and page_size
	// TODO: Add proper parsing with error handling
	
	response, err := h.service.ListSuppliers(r.Context(), tenantID, req)
	if err != nil {
		h.logger.Error("Failed to list suppliers", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list suppliers")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// VerifySupplier handles verifying a supplier
func (h *SupplierHandler) VerifySupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	// TODO: Get user from auth context
	verifiedBy := "system"

	response, err := h.service.VerifySupplier(r.Context(), tenantID, supplierID, verifiedBy)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to verify supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to verify supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// RejectSupplier handles rejecting a supplier verification
func (h *SupplierHandler) RejectSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	// TODO: Get user from auth context
	rejectedBy := "system"

	response, err := h.service.RejectSupplier(r.Context(), tenantID, supplierID, rejectedBy)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to reject supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to reject supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// SuspendSupplier handles suspending a supplier
func (h *SupplierHandler) SuspendSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	response, err := h.service.SuspendSupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to suspend supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to suspend supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// ActivateSupplier handles activating a supplier
func (h *SupplierHandler) ActivateSupplier(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	response, err := h.service.ActivateSupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to activate supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to activate supplier")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// AddCertification handles adding a certification to a supplier
func (h *SupplierHandler) AddCertification(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	supplierID := chi.URLParam(r, "id")
	if supplierID == "" {
		h.respondError(w, http.StatusBadRequest, "Supplier ID is required")
		return
	}

	var cert app.CertificationDTO
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.AddCertification(r.Context(), tenantID, supplierID, cert)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			h.respondError(w, http.StatusNotFound, "Supplier not found")
			return
		}
		h.logger.Error("Failed to add certification", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to add certification")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// GetSuppliersByCategory handles retrieving suppliers by category
func (h *SupplierHandler) GetSuppliersByCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header is required")
		return
	}

	categoryID := chi.URLParam(r, "categoryId")
	if categoryID == "" {
		h.respondError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	response, err := h.service.GetSuppliersByCategory(r.Context(), tenantID, categoryID)
	if err != nil {
		h.logger.Error("Failed to get suppliers by category", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get suppliers by category")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// Helper methods

func (h *SupplierHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *SupplierHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
