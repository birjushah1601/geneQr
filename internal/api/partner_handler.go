package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/aby-med/medical-platform/internal/services"
	"github.com/gorilla/mux"
)

// PartnerHandler handles partner association HTTP requests
type PartnerHandler struct {
	partnerService *services.PartnerService
	networkService *services.NetworkEngineersService
}

// NewPartnerHandler creates a new partner handler
func NewPartnerHandler(db *sql.DB) *PartnerHandler {
	return &PartnerHandler{
		partnerService: services.NewPartnerService(db),
		networkService: services.NewNetworkEngineersService(db),
	}
}

// RegisterRoutes registers partner routes
func (h *PartnerHandler) RegisterRoutes(r *mux.Router) {
	// Partner management routes
	r.HandleFunc("/api/v1/organizations/{manufacturerId}/partners", h.ListPartners).Methods("GET")
	r.HandleFunc("/api/v1/organizations/{manufacturerId}/available-partners", h.GetAvailablePartners).Methods("GET")
	r.HandleFunc("/api/v1/organizations/{manufacturerId}/partners", h.AssociatePartner).Methods("POST")
	r.HandleFunc("/api/v1/organizations/{manufacturerId}/partners/{partnerId}", h.RemovePartner).Methods("DELETE")
	
	// Network engineers route
	r.HandleFunc("/api/v1/engineers/network/{manufacturerId}", h.GetNetworkEngineers).Methods("GET")
}

// ListPartners handles GET /api/v1/organizations/:manufacturerId/partners
func (h *PartnerHandler) ListPartners(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	manufacturerID := vars["manufacturerId"]

	// Get query parameters for filtering
	filters := make(map[string]string)
	if orgType := r.URL.Query().Get("type"); orgType != "" {
		filters["type"] = orgType
	}
	if assocType := r.URL.Query().Get("association_type"); assocType != "" {
		filters["association_type"] = assocType
	}

	partners, err := h.partnerService.GetPartners(r.Context(), manufacturerID, filters)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get partners",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"partners": partners,
		"total":    len(partners),
	})
}

// GetAvailablePartners handles GET /api/v1/organizations/:manufacturerId/available-partners
func (h *PartnerHandler) GetAvailablePartners(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	manufacturerID := vars["manufacturerId"]
	search := r.URL.Query().Get("search")

	orgs, err := h.partnerService.GetAvailablePartners(r.Context(), manufacturerID, search)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get available partners",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"organizations": orgs,
		"total":         len(orgs),
	})
}

// AssociatePartner handles POST /api/v1/organizations/:manufacturerId/partners
func (h *PartnerHandler) AssociatePartner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	manufacturerID := vars["manufacturerId"]

	var req struct {
		PartnerOrgID string  `json:"partner_org_id"`
		EquipmentID  *string `json:"equipment_id,omitempty"`
		RelType      string  `json:"rel_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Validate required fields
	if req.PartnerOrgID == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Validation error",
			Message: "partner_org_id is required",
		})
		return
	}

	// Set default rel_type if not provided
	if req.RelType == "" {
		req.RelType = "services_for"
	}

	// Create association
	assoc, err := h.partnerService.CreateAssociation(r.Context(), services.CreateAssociationRequest{
		ManufacturerID: manufacturerID,
		PartnerOrgID:   req.PartnerOrgID,
		EquipmentID:    req.EquipmentID,
		RelType:        req.RelType,
	})

	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create association",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusCreated, assoc)
}

// RemovePartner handles DELETE /api/v1/organizations/:manufacturerId/partners/:partnerId
func (h *PartnerHandler) RemovePartner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	manufacturerID := vars["manufacturerId"]
	partnerID := vars["partnerId"]

	// Check for equipment_id query parameter
	var equipmentID *string
	if eqID := r.URL.Query().Get("equipment_id"); eqID != "" {
		equipmentID = &eqID
	}

	err := h.partnerService.RemoveAssociation(r.Context(), manufacturerID, partnerID, equipmentID)
	if err != nil {
		if err.Error() == "association not found" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "Association not found",
				Message: "The specified partner association does not exist",
			})
			return
		}
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to remove association",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Partner association removed successfully",
		"removed": map[string]interface{}{
			"manufacturer_id": manufacturerID,
			"partner_id":      partnerID,
			"equipment_id":    equipmentID,
		},
	})
}

// GetNetworkEngineers handles GET /api/v1/engineers/network/:manufacturerId
func (h *PartnerHandler) GetNetworkEngineers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	manufacturerID := vars["manufacturerId"]

	// Check for equipment_id query parameter
	var equipmentID *string
	if eqID := r.URL.Query().Get("equipment_id"); eqID != "" {
		equipmentID = &eqID
	}

	result, err := h.networkService.GetNetworkEngineers(r.Context(), manufacturerID, equipmentID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get network engineers",
			Message: err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, result)
}
