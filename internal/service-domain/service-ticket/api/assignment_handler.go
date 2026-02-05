package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/go-chi/chi/v5"
)

// AssignmentHandler handles HTTP requests for engineer assignment
type AssignmentHandler struct {
	service *app.AssignmentService
	logger  *slog.Logger
}

// NewAssignmentHandler creates a new assignment HTTP handler
func NewAssignmentHandler(service *app.AssignmentService, logger *slog.Logger) *AssignmentHandler {
	return &AssignmentHandler{
		service: service,
		logger:  logger.With(slog.String("component", "assignment_handler")),
	}
}

// ListEngineers handles GET /engineers or GET /organizations/{orgId}/engineers
// Query parameters:
//   - limit: max number of engineers to return (default: 100)
//   - offset: pagination offset (default: 0)
//   - include_partners: include engineers from partner organizations (default: false)
func (h *AssignmentHandler) ListEngineers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Check if organization ID is in URL path
	var orgID *string
	if orgIDParam := chi.URLParam(r, "orgId"); orgIDParam != "" {
		orgID = &orgIDParam
	} else if orgIDQuery := r.URL.Query().Get("organization_id"); orgIDQuery != "" {
		orgID = &orgIDQuery
	}
	
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	
	// Parse include_partners parameter (default: false)
	includePartners := r.URL.Query().Get("include_partners") == "true"
	
	engineers, err := h.service.ListEngineers(ctx, orgID, includePartners, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list engineers", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list engineers: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"engineers":        engineers,
		"total":            len(engineers),
		"include_partners": includePartners,
	})
}

// GetEngineer handles GET /engineers/{id}
func (h *AssignmentHandler) GetEngineer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	engineer, err := h.service.GetEngineer(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get engineer", slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Engineer not found")
		return
	}
	
	h.respondJSON(w, http.StatusOK, engineer)
}

// UpdateEngineerLevel handles PUT /engineers/{id}/level
func (h *AssignmentHandler) UpdateEngineerLevel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	var req struct {
		Level string `json:"level"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	level := domain.EngineerLevel(req.Level)
	if level != domain.EngineerLevelL1 && level != domain.EngineerLevelL2 && level != domain.EngineerLevelL3 {
		h.respondError(w, http.StatusBadRequest, "Invalid level. Must be L1, L2, or L3")
		return
	}
	
	if err := h.service.UpdateEngineerLevel(ctx, id, level); err != nil {
		h.logger.Error("Failed to update engineer level", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update engineer level")
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Engineer level updated successfully"})
}

// ListEngineerEquipmentTypes handles GET /engineers/{id}/equipment-types
func (h *AssignmentHandler) ListEngineerEquipmentTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	types, err := h.service.ListEngineerEquipmentTypes(ctx, id)
	if err != nil {
		h.logger.Error("Failed to list engineer equipment types", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list equipment types")
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"equipment_types": types,
		"total":           len(types),
	})
}

// AddEngineerEquipmentType handles POST /engineers/{id}/equipment-types
func (h *AssignmentHandler) AddEngineerEquipmentType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	var req struct {
		Manufacturer string `json:"manufacturer"`
		Category     string `json:"category"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	if req.Manufacturer == "" || req.Category == "" {
		h.respondError(w, http.StatusBadRequest, "Manufacturer and category are required")
		return
	}
	
	if err := h.service.AddEngineerEquipmentType(ctx, id, req.Manufacturer, req.Category); err != nil {
		h.logger.Error("Failed to add engineer equipment type", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to add equipment type: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Equipment type added successfully"})
}

// RemoveEngineerEquipmentType handles DELETE /engineers/{id}/equipment-types
func (h *AssignmentHandler) RemoveEngineerEquipmentType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	manufacturer := r.URL.Query().Get("manufacturer")
	category := r.URL.Query().Get("category")
	
	if manufacturer == "" || category == "" {
		h.respondError(w, http.StatusBadRequest, "Manufacturer and category query parameters are required")
		return
	}
	
	if err := h.service.RemoveEngineerEquipmentType(ctx, id, manufacturer, category); err != nil {
		h.logger.Error("Failed to remove engineer equipment type", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to remove equipment type")
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Equipment type removed successfully"})
}

// GetSuggestedEngineers handles GET /tickets/{id}/suggested-engineers
func (h *AssignmentHandler) GetSuggestedEngineers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	
	if ticketID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}
	
	suggestions, err := h.service.GetSuggestedEngineers(ctx, ticketID)
	if err != nil {
		h.logger.Error("Failed to get suggested engineers", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get suggestions: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"suggested_engineers": suggestions,
		"total":               len(suggestions),
	})
}

// AssignEngineer handles POST /tickets/{id}/assign-engineer
func (h *AssignmentHandler) AssignEngineer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	
	if ticketID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}
	
	var req app.AssignEngineerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	req.TicketID = ticketID // Override with URL parameter
	
	if req.EngineerID == "" {
		h.respondError(w, http.StatusBadRequest, "Engineer ID is required")
		return
	}
	
	if err := h.service.AssignEngineer(ctx, req); err != nil {
		h.logger.Error("Failed to assign engineer", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to assign engineer: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Engineer assigned successfully"})
}

// GetEquipmentServiceConfig handles GET /equipment/{id}/service-config
func (h *AssignmentHandler) GetEquipmentServiceConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	equipmentID := chi.URLParam(r, "id")
	
	if equipmentID == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}
	
	config, err := h.service.GetEquipmentServiceConfig(ctx, equipmentID)
	if err != nil {
		h.logger.Error("Failed to get equipment service config", slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Service config not found")
		return
	}
	
	h.respondJSON(w, http.StatusOK, config)
}

// CreateEquipmentServiceConfig handles POST /equipment/{id}/service-config
func (h *AssignmentHandler) CreateEquipmentServiceConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	equipmentID := chi.URLParam(r, "id")
	
	if equipmentID == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}
	
	var config domain.EquipmentServiceConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	config.EquipmentID = equipmentID // Override with URL parameter
	
	if err := h.service.CreateEquipmentServiceConfig(ctx, &config); err != nil {
		h.logger.Error("Failed to create equipment service config", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to create service config: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Service config created successfully"})
}

// UpdateEquipmentServiceConfig handles PUT /equipment/{id}/service-config
func (h *AssignmentHandler) UpdateEquipmentServiceConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	equipmentID := chi.URLParam(r, "id")
	
	if equipmentID == "" {
		h.respondError(w, http.StatusBadRequest, "Equipment ID is required")
		return
	}
	
	var config domain.EquipmentServiceConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	config.EquipmentID = equipmentID // Override with URL parameter
	
	if err := h.service.UpdateEquipmentServiceConfig(ctx, &config); err != nil {
		h.logger.Error("Failed to update equipment service config", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update service config: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Service config updated successfully"})
}

// respondJSON writes JSON response
func (h *AssignmentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func (h *AssignmentHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

// GetAssignmentHistory handles GET /tickets/{id}/assignments/history
func (h *AssignmentHandler) GetAssignmentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	
	if ticketID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}
	
	history, err := h.service.GetAssignmentHistory(ctx, ticketID)
	if err != nil {
		h.logger.Error("Failed to get assignment history",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get assignment history: "+err.Error())
		return
	}
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"ticket_id":   ticketID,
		"count":       len(history),
		"assignments": history,
	})
}
