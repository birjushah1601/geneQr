package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/contract/app"
	"github.com/go-chi/chi/v5"
)

// ContractHandler handles HTTP requests for contracts
type ContractHandler struct {
	service *app.ContractService
	logger  *slog.Logger
}

// NewContractHandler creates a new contract handler
func NewContractHandler(service *app.ContractService, logger *slog.Logger) *ContractHandler {
	return &ContractHandler{
		service: service,
		logger:  logger.With(slog.String("handler", "contract")),
	}
}

// CreateContract handles contract creation
func (h *ContractHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var req app.CreateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	createdBy := r.Header.Get("X-User-ID")
	if createdBy == "" {
		createdBy = "system"
	}

	response, err := h.service.CreateContract(r.Context(), tenantID, createdBy, req)
	if err != nil {
		h.logger.Error("Failed to create contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, response)
}

// GetContract retrieves a contract by ID
func (h *ContractHandler) GetContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	response, err := h.service.GetContract(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to get contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Contract not found")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// GetContractByNumber retrieves a contract by contract number
func (h *ContractHandler) GetContractByNumber(w http.ResponseWriter, r *http.Request) {
	contractNumber := chi.URLParam(r, "number")
	tenantID := r.Header.Get("X-Tenant-ID")

	response, err := h.service.GetContractByNumber(r.Context(), tenantID, contractNumber)
	if err != nil {
		h.logger.Error("Failed to get contract by number", slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Contract not found")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// ListContracts retrieves contracts with filtering
func (h *ContractHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")

	// Parse query parameters
	req := app.ListContractsRequest{
		RFQID:         r.URL.Query().Get("rfq_id"),
		SupplierID:    r.URL.Query().Get("supplier_id"),
		CreatedBy:     r.URL.Query().Get("created_by"),
		SortBy:        r.URL.Query().Get("sort_by"),
		SortDirection: r.URL.Query().Get("sort_direction"),
	}

	// Parse status filter
	if statuses := r.URL.Query()["status"]; len(statuses) > 0 {
		req.Status = statuses
	}

	// Parse pagination
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	response, err := h.service.ListContracts(r.Context(), tenantID, req)
	if err != nil {
		h.logger.Error("Failed to list contracts", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// UpdateContract updates a contract
func (h *ContractHandler) UpdateContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.UpdateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.UpdateContract(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to update contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// ActivateContract activates a contract
func (h *ContractHandler) ActivateContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	response, err := h.service.ActivateContract(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to activate contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// SignContract signs a contract
func (h *ContractHandler) SignContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.SignContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.SignContract(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to sign contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// CompleteContract completes a contract
func (h *ContractHandler) CompleteContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	response, err := h.service.CompleteContract(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to complete contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// CancelContract cancels a contract
func (h *ContractHandler) CancelContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.CancelContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.CancelContract(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to cancel contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// SuspendContract suspends a contract
func (h *ContractHandler) SuspendContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.SuspendContractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.SuspendContract(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to suspend contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// ResumeContract resumes a suspended contract
func (h *ContractHandler) ResumeContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	response, err := h.service.ResumeContract(r.Context(), tenantID, id)
	if err != nil {
		h.logger.Error("Failed to resume contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// AddAmendment adds an amendment to a contract
func (h *ContractHandler) AddAmendment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.AddAmendmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.AddAmendment(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to add amendment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// MarkPaymentPaid marks a payment as paid
func (h *ContractHandler) MarkPaymentPaid(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.MarkPaymentPaidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.MarkPaymentPaid(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to mark payment paid", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// MarkDeliveryCompleted marks a delivery milestone as completed
func (h *ContractHandler) MarkDeliveryCompleted(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	var req app.MarkDeliveryCompletedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.service.MarkDeliveryCompleted(r.Context(), tenantID, id, req)
	if err != nil {
		h.logger.Error("Failed to mark delivery completed", slog.String("error", err.Error()))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

// DeleteContract deletes a contract
func (h *ContractHandler) DeleteContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenantID := r.Header.Get("X-Tenant-ID")

	if err := h.service.DeleteContract(r.Context(), tenantID, id); err != nil {
		h.logger.Error("Failed to delete contract", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetContractsByRFQ retrieves all contracts for an RFQ
func (h *ContractHandler) GetContractsByRFQ(w http.ResponseWriter, r *http.Request) {
	rfqID := chi.URLParam(r, "rfq_id")
	tenantID := r.Header.Get("X-Tenant-ID")

	responses, err := h.service.GetContractsByRFQ(r.Context(), tenantID, rfqID)
	if err != nil {
		h.logger.Error("Failed to get contracts by RFQ", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"contracts": responses,
		"total":     len(responses),
	})
}

// GetContractsBySupplier retrieves all contracts for a supplier
func (h *ContractHandler) GetContractsBySupplier(w http.ResponseWriter, r *http.Request) {
	supplierID := chi.URLParam(r, "supplier_id")
	tenantID := r.Header.Get("X-Tenant-ID")

	responses, err := h.service.GetContractsBySupplier(r.Context(), tenantID, supplierID)
	if err != nil {
		h.logger.Error("Failed to get contracts by supplier", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"contracts": responses,
		"total":     len(responses),
	})
}

// respondJSON sends a JSON response
func (h *ContractHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", slog.String("error", err.Error()))
	}
}

// respondError sends an error response
func (h *ContractHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
