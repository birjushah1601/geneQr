package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"context"

	"github.com/aby-med/medical-platform/internal/marketplace/catalog/app"
	"github.com/aby-med/medical-platform/internal/marketplace/catalog/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

const (
	// Default pagination values
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100

	// Header keys
	headerTenantID = "X-Tenant-ID"
)

// Handler handles HTTP requests for the catalog module
type Handler struct {
	service   *app.CatalogService
	logger    *slog.Logger
	validator *validator.Validate
}

// NewHandler creates a new catalog HTTP handler
func NewHandler(service *app.CatalogService, logger *slog.Logger) *Handler {
	return &Handler{
		service:   service,
		logger:    logger.With(slog.String("component", "catalog_http_handler")),
		validator: validator.New(),
	}
}

// apiError represents an API error response
type apiError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ListEquipment handles GET /catalog requests to list equipment
func (h *Handler) ListEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Parse pagination parameters
	page, pageSize := h.getPaginationParams(r)

	// Create search request with minimal filters
	searchReq := app.SearchEquipmentRequest{
		Page:     page,
		PageSize: pageSize,
		IsActive: boolPtr(true), // Default to active equipment only
	}

	// Execute search
	result, err := h.service.SearchEquipment(ctx, searchReq)
	if err != nil {
		h.logger.Error("Failed to list equipment", 
			slog.String("error", err.Error()),
			slog.Int("page", page),
			slog.Int("page_size", pageSize))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "list_failed",
			Message: "Failed to retrieve equipment list",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, result)
}

// GetEquipment handles GET /catalog/{id} requests to get a specific equipment
func (h *Handler) GetEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Extract equipment ID from URL
	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_id",
			Message: "Equipment ID is required",
		})
		return
	}

	// Get equipment from service
	equipment, err := h.service.GetEquipment(ctx, equipmentID)
	if err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			h.respondWithError(w, apiError{
				Status:  http.StatusNotFound,
				Code:    "equipment_not_found",
				Message: "Equipment not found",
			})
			return
		}
		h.logger.Error("Failed to get equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipmentID))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "retrieval_failed",
			Message: "Failed to retrieve equipment",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, equipment)
}

// CreateEquipment handles POST /catalog requests to create new equipment
func (h *Handler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Parse request body
	var req app.CreateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "Validation failed",
			Details: validationErrors,
		})
		return
	}

	// Create equipment
	equipment, err := h.service.CreateEquipment(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCategory) || 
		   errors.Is(err, domain.ErrInvalidManufacturer) ||
		   errors.Is(err, domain.ErrInvalidPrice) {
			h.respondWithError(w, apiError{
				Status:  http.StatusBadRequest,
				Code:    "invalid_input",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to create equipment", slog.String("error", err.Error()))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "creation_failed",
			Message: "Failed to create equipment",
		})
		return
	}

	h.respondWithJSON(w, http.StatusCreated, equipment)
}

// UpdateEquipment handles PUT /catalog/{id} requests to update equipment
func (h *Handler) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Extract equipment ID from URL
	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_id",
			Message: "Equipment ID is required",
		})
		return
	}

	// Parse request body
	var req app.UpdateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "invalid_request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := h.formatValidationErrors(err)
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "Validation failed",
			Details: validationErrors,
		})
		return
	}

	// Update equipment
	equipment, err := h.service.UpdateEquipment(ctx, equipmentID, req)
	if err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			h.respondWithError(w, apiError{
				Status:  http.StatusNotFound,
				Code:    "equipment_not_found",
				Message: "Equipment not found",
			})
			return
		}
		if errors.Is(err, domain.ErrInvalidCategory) || 
		   errors.Is(err, domain.ErrInvalidManufacturer) ||
		   errors.Is(err, domain.ErrInvalidPrice) {
			h.respondWithError(w, apiError{
				Status:  http.StatusBadRequest,
				Code:    "invalid_input",
				Message: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to update equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipmentID))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "update_failed",
			Message: "Failed to update equipment",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, equipment)
}

// DeleteEquipment handles DELETE /catalog/{id} requests to delete equipment
func (h *Handler) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Extract equipment ID from URL
	equipmentID := chi.URLParam(r, "id")
	if equipmentID == "" {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_id",
			Message: "Equipment ID is required",
		})
		return
	}

	// Delete equipment
	if err := h.service.DeleteEquipment(ctx, equipmentID); err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			h.respondWithError(w, apiError{
				Status:  http.StatusNotFound,
				Code:    "equipment_not_found",
				Message: "Equipment not found",
			})
			return
		}
		h.logger.Error("Failed to delete equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipmentID))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "deletion_failed",
			Message: "Failed to delete equipment",
		})
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}

// SearchEquipment handles GET /catalog/search requests to search equipment
func (h *Handler) SearchEquipment(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Parse pagination parameters
	page, pageSize := h.getPaginationParams(r)

	// Parse search parameters
	query := r.URL.Query().Get("q")
	
	// Parse category IDs (can be multiple)
	var categoryIDs []string
	if categoryIDsParam := r.URL.Query().Get("category_ids"); categoryIDsParam != "" {
		categoryIDs = strings.Split(categoryIDsParam, ",")
	}
	
	// Parse manufacturer IDs (can be multiple)
	var manufacturerIDs []string
	if manufacturerIDsParam := r.URL.Query().Get("manufacturer_ids"); manufacturerIDsParam != "" {
		manufacturerIDs = strings.Split(manufacturerIDsParam, ",")
	}
	
	// Parse price range
	var priceMin, priceMax *float64
	if priceMinParam := r.URL.Query().Get("price_min"); priceMinParam != "" {
		if val, err := strconv.ParseFloat(priceMinParam, 64); err == nil {
			priceMin = &val
		}
	}
	if priceMaxParam := r.URL.Query().Get("price_max"); priceMaxParam != "" {
		if val, err := strconv.ParseFloat(priceMaxParam, 64); err == nil {
			priceMax = &val
		}
	}
	
	// Parse active status
	var isActive *bool
	if activeParam := r.URL.Query().Get("active"); activeParam != "" {
		if val, err := strconv.ParseBool(activeParam); err == nil {
			isActive = &val
		}
	}
	
	// Parse sort parameters
	sortBy := r.URL.Query().Get("sort_by")
	sortDirection := r.URL.Query().Get("sort_direction")

	// Create search request
	searchReq := app.SearchEquipmentRequest{
		Query:           query,
		CategoryIDs:     categoryIDs,
		ManufacturerIDs: manufacturerIDs,
		PriceMin:        priceMin,
		PriceMax:        priceMax,
		IsActive:        isActive,
		Page:            page,
		PageSize:        pageSize,
		SortBy:          sortBy,
		SortDirection:   sortDirection,
	}

	// Execute search
	result, err := h.service.SearchEquipment(ctx, searchReq)
	if err != nil {
		h.logger.Error("Failed to search equipment", 
			slog.String("error", err.Error()),
			slog.String("query", query))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "search_failed",
			Message: "Failed to search equipment",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, result)
}

// ListCategories handles GET /catalog/categories requests to list categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Get categories from service
	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		h.logger.Error("Failed to list categories", slog.String("error", err.Error()))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "list_failed",
			Message: "Failed to retrieve categories",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, categories)
}

// ListManufacturers handles GET /catalog/manufacturers requests to list manufacturers
func (h *Handler) ListManufacturers(w http.ResponseWriter, r *http.Request) {
	// Extract tenant ID from header
	ctx, err := h.extractTenantContext(r)
	if err != nil {
		h.respondWithError(w, apiError{
			Status:  http.StatusBadRequest,
			Code:    "missing_tenant",
			Message: "Tenant ID is required",
		})
		return
	}

	// Get manufacturers from service
	manufacturers, err := h.service.ListManufacturers(ctx)
	if err != nil {
		h.logger.Error("Failed to list manufacturers", slog.String("error", err.Error()))
		h.respondWithError(w, apiError{
			Status:  http.StatusInternalServerError,
			Code:    "list_failed",
			Message: "Failed to retrieve manufacturers",
		})
		return
	}

	h.respondWithJSON(w, http.StatusOK, manufacturers)
}

// Helper methods

// extractTenantContext extracts tenant ID from request header and adds it to context
func (h *Handler) extractTenantContext(r *http.Request) (context.Context, error) {
	tenantID := r.Header.Get(headerTenantID)
	if tenantID == "" {
		return nil, errors.New("tenant ID not found in request header")
	}
	return domain.WithTenantID(r.Context(), tenantID), nil
}

// getPaginationParams extracts and validates pagination parameters
func (h *Handler) getPaginationParams(r *http.Request) (page, pageSize int) {
	// Parse page parameter
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = defaultPage
	} else {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		} else {
			page = defaultPage
		}
	}

	// Parse page_size parameter
	pageSizeStr := r.URL.Query().Get("page_size")
	if pageSizeStr == "" {
		pageSize = defaultPageSize
	} else {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
			// Limit maximum page size
			if pageSize > maxPageSize {
				pageSize = maxPageSize
			}
		} else {
			pageSize = defaultPageSize
		}
	}

	return page, pageSize
}

// respondWithJSON sends a JSON response with the given status code
func (h *Handler) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
			// If we can't encode the JSON, send a plain text error
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error: failed to encode response"))
		}
	}
}

// respondWithError sends an error response with the given status code
func (h *Handler) respondWithError(w http.ResponseWriter, err apiError) {
	h.respondWithJSON(w, err.Status, err)
}

// formatValidationErrors formats validation errors into a user-friendly format
func (h *Handler) formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = "This field is required"
			case "min":
				errors[field] = "Value is below minimum"
			case "max":
				errors[field] = "Value is above maximum"
			case "gt":
				errors[field] = "Value must be greater than " + e.Param()
			default:
				errors[field] = "Invalid value"
			}
		}
	}
	
	return errors
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}
