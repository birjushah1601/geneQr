package api

import (
    "database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/aby-med/medical-platform/internal/shared/audit"
	emailInfra "github.com/aby-med/medical-platform/internal/infrastructure/email"
	"github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// TicketHandler handles HTTP requests for service tickets
type TicketHandler struct {
	service             *app.TicketService
	notificationService *app.NotificationService
	logger              *slog.Logger
    pool                *pgxpool.Pool
	auditLogger         *audit.AuditLogger
}

// NewTicketHandler creates a new ticket HTTP handler
func NewTicketHandler(service *app.TicketService, logger *slog.Logger, pool *pgxpool.Pool, auditLogger *audit.AuditLogger) *TicketHandler {
	return &TicketHandler{
		service:     service,
		logger:      logger.With(slog.String("component", "ticket_handler")),
        pool:        pool,
		auditLogger: auditLogger,
	}
}

// SetNotificationService sets the notification service (called after initialization)
func (h *TicketHandler) SetNotificationService(notificationService *app.NotificationService) {
	h.notificationService = notificationService
}

// CreateTicket handles POST /tickets
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()

	h.logger.Info("🎫 CreateTicket handler - Starting",
		slog.Int64("content_length", r.ContentLength),
		slog.String("content_type", r.Header.Get("Content-Type")))

	var req app.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("❌ Failed to decode request body",
			slog.String("error", err.Error()),
			slog.Int64("content_length", r.ContentLength))
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	h.logger.Info("✅ Request decoded successfully",
		slog.String("equipment_id", req.EquipmentID),
		slog.String("qr_code", req.QRCode),
		slog.String("customer_phone", req.CustomerPhone),
		slog.String("customer_email", req.CustomerEmail),
		slog.String("customer_name", req.CustomerName))

	ticket, err := h.service.CreateTicket(ctx, req)
	if err != nil {
		h.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to create ticket: "+err.Error())
		
		// Log failure to audit
		if h.auditLogger != nil {
			errorMsg := err.Error()
			durationMs := int(time.Since(startTime).Milliseconds())
			h.auditLogger.LogAsync(ctx, &audit.AuditEvent{
				EventType:     "ticket_create_failed",
				EventCategory: audit.CategoryTicket,
				EventAction:   audit.ActionCreate,
				EventStatus:   audit.StatusFailure,
				IPAddress:     stringPtr(audit.ExtractIPAddress(r)),
				UserAgent:     stringPtr(audit.ExtractUserAgent(r)),
				RequestMethod: stringPtr(r.Method),
				RequestPath:   stringPtr(r.URL.Path),
				ErrorMessage:  &errorMsg,
				DurationMs:    &durationMs,
				Metadata: map[string]interface{}{
					"qr_code":      req.QRCode,
					"equipment_id": req.EquipmentID,
				},
			})
		}
		return
	}

	// If parts_requested are provided, create ticket_parts entries
	if len(req.PartsRequested) > 0 && h.pool != nil {
		h.logger.Info("Creating ticket_parts entries",
			slog.String("ticket_id", ticket.ID),
			slog.Int("parts_count", len(req.PartsRequested)))

		for _, part := range req.PartsRequested {
			// Query spare_parts_catalog to get spare_part_id by part_number
			var sparePartID string
			err := h.pool.QueryRow(ctx,
				`SELECT id FROM spare_parts_catalog WHERE part_number = $1 LIMIT 1`,
				part.PartNumber,
			).Scan(&sparePartID)

			if err != nil {
				h.logger.Warn("Part not found in catalog, skipping",
					slog.String("part_number", part.PartNumber),
					slog.String("error", err.Error()))
				continue
			}

			// Insert into ticket_parts
			_, err = h.pool.Exec(ctx, `
				INSERT INTO ticket_parts (
					ticket_id, spare_part_id, quantity_required,
					unit_price, total_price, status, notes, assigned_at
				) VALUES ($1, $2, $3, $4, $5, 'pending', $6, NOW())
			`,
				ticket.ID,
				sparePartID,
				part.Quantity,
				part.UnitPrice,
				part.TotalPrice,
				part.Description,
			)

			if err != nil {
				h.logger.Warn("Failed to create ticket_part entry",
					slog.String("part_number", part.PartNumber),
					slog.String("error", err.Error()))
			} else {
				h.logger.Info("Created ticket_part entry",
					slog.String("ticket_id", ticket.ID),
					slog.String("spare_part_id", sparePartID))
			}
		}
	}

	// Generate tracking token for the ticket
	var trackingToken string
	var trackingURL string
	if h.notificationService != nil {
		token, err := h.notificationService.GetOrCreateTrackingToken(ticket.ID)
		if err != nil {
			h.logger.Warn("Failed to create tracking token",
				slog.String("ticket_id", ticket.ID),
				slog.String("error", err.Error()))
		} else {
			trackingToken = token
			baseURL := "http://localhost:3000" // TODO: Get from config
			if envURL := os.Getenv("FRONTEND_BASE_URL"); envURL != "" {
				baseURL = envURL
			}
			trackingURL = fmt.Sprintf("%s/track/%s", baseURL, trackingToken)
			
			h.logger.Info("Tracking token generated",
				slog.String("ticket_id", ticket.ID),
				slog.String("token", trackingToken[:8]+"..."))
		}
	}

	// Log successful ticket creation to audit
	if h.auditLogger != nil {
		durationMs := int(time.Since(startTime).Milliseconds())
		ipAddress := audit.ExtractIPAddress(r)
		
		metadata := map[string]interface{}{
			"qr_code":        req.QRCode,
			"equipment_id":   req.EquipmentID,
			"priority":       req.Priority,
			"parts_count":    len(req.PartsRequested),
			"source":         req.Source,
		}
		
		if req.CustomerName != "" {
			metadata["customer_name"] = req.CustomerName
		}
		if req.CustomerPhone != "" {
			metadata["customer_phone"] = req.CustomerPhone
		}
		if trackingToken != "" {
			metadata["tracking_token_generated"] = true
		}

		h.auditLogger.LogAsync(ctx, &audit.AuditEvent{
			EventType:     "ticket_created",
			EventCategory: audit.CategoryTicket,
			EventAction:   audit.ActionCreate,
			EventStatus:   audit.StatusSuccess,
			ResourceType:  stringPtr("ticket"),
			ResourceID:    &ticket.ID,
			ResourceName:  stringPtr(ticket.TicketNumber),
			IPAddress:     &ipAddress,
			UserAgent:     stringPtr(audit.ExtractUserAgent(r)),
			RequestMethod: stringPtr(r.Method),
			RequestPath:   stringPtr(r.URL.Path),
			DurationMs:    &durationMs,
			Metadata:      metadata,
		})
	}

	// Build response with ticket and tracking information
	response := map[string]interface{}{
		"ticket":        ticket,
		"tracking_url":  trackingURL,
		"tracking_token": trackingToken,
	}

	h.respondJSON(w, http.StatusCreated, response)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

// GetTicket handles GET /tickets/{id}
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	ticket, err := h.service.GetTicket(ctx, id)
	if err != nil {
		if err == domain.ErrTicketNotFound {
			h.respondError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		h.logger.Error("Failed to get ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get ticket")
		return
	}

	h.respondJSON(w, http.StatusOK, ticket)
}

// GetTicketByNumber handles GET /tickets/number/{number}
func (h *TicketHandler) GetTicketByNumber(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketNumber := chi.URLParam(r, "number")

	if ticketNumber == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket number is required")
		return
	}

	ticket, err := h.service.GetTicketByNumber(ctx, ticketNumber)
	if err != nil {
		if err == domain.ErrTicketNotFound {
			h.respondError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		h.logger.Error("Failed to get ticket by number", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get ticket")
		return
	}

	h.respondJSON(w, http.StatusOK, ticket)
}

// ListTickets handles GET /tickets
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	criteria := domain.ListCriteria{
		EquipmentID: r.URL.Query().Get("equipment_id"),
		CustomerID:  r.URL.Query().Get("customer_id"),
		EngineerID:  r.URL.Query().Get("engineer_id"),
		SortBy:      r.URL.Query().Get("sort_by"),
		SortDirection: r.URL.Query().Get("sort_dir"),
	}

	// Parse status filter (multiple values)
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		criteria.Status = []domain.TicketStatus{domain.TicketStatus(statusStr)}
	}

	// Parse priority filter
	if priorityStr := r.URL.Query().Get("priority"); priorityStr != "" {
		criteria.Priority = []domain.TicketPriority{domain.TicketPriority(priorityStr)}
	}

	// Parse source filter
	if sourceStr := r.URL.Query().Get("source"); sourceStr != "" {
		criteria.Source = []domain.TicketSource{domain.TicketSource(sourceStr)}
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
	if slaBreached := r.URL.Query().Get("sla_breached"); slaBreached != "" {
		val := slaBreached == "true"
		criteria.SLABreached = &val
	}

	if coveredUnderAMC := r.URL.Query().Get("covered_under_amc"); coveredUnderAMC != "" {
		val := coveredUnderAMC == "true"
		criteria.CoveredUnderAMC = &val
	}

	result, err := h.service.ListTickets(ctx, criteria)
	if err != nil {
		h.logger.Error("Failed to list tickets", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to list tickets")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// AssignTicket handles POST /tickets/{id}/assign
func (h *TicketHandler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		EngineerID   string `json:"engineer_id"`
		EngineerName string `json:"engineer_name"`
		AssignedBy   string `json:"assigned_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.AssignTicket(ctx, id, req.EngineerID, req.EngineerName, req.AssignedBy); err != nil {
		h.logger.Error("Failed to assign ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to assign ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket assigned successfully"})
}

// AcknowledgeTicket handles POST /tickets/{id}/acknowledge
func (h *TicketHandler) AcknowledgeTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.AcknowledgeTicket(ctx, id, req.AcknowledgedBy); err != nil {
		h.logger.Error("Failed to acknowledge ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to acknowledge ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket acknowledged successfully"})
}

// StartWork handles POST /tickets/{id}/start
func (h *TicketHandler) StartWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		StartedBy string `json:"started_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.StartWork(ctx, id, req.StartedBy); err != nil {
		h.logger.Error("Failed to start work on ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to start work: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Work started successfully"})
}

// PutOnHold handles POST /tickets/{id}/hold
func (h *TicketHandler) PutOnHold(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		Reason    string `json:"reason"`
		ChangedBy string `json:"changed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.PutOnHold(ctx, id, req.Reason, req.ChangedBy); err != nil {
		h.logger.Error("Failed to put ticket on hold", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to put on hold: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket put on hold successfully"})
}

// ResumeWork handles POST /tickets/{id}/resume
func (h *TicketHandler) ResumeWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		ResumedBy string `json:"resumed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.ResumeWork(ctx, id, req.ResumedBy); err != nil {
		h.logger.Error("Failed to resume work on ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to resume work: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Work resumed successfully"})
}

// ResolveTicket handles POST /tickets/{id}/resolve
func (h *TicketHandler) ResolveTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.ResolveTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.ResolveTicket(ctx, id, req); err != nil {
		h.logger.Error("Failed to resolve ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to resolve ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket resolved successfully"})
}

// CloseTicket handles POST /tickets/{id}/close
func (h *TicketHandler) CloseTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		ClosedBy string `json:"closed_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CloseTicket(ctx, id, req.ClosedBy); err != nil {
		h.logger.Error("Failed to close ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to close ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket closed successfully"})
}

// CancelTicket handles POST /tickets/{id}/cancel
func (h *TicketHandler) CancelTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		Reason      string `json:"reason"`
		CancelledBy string `json:"cancelled_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CancelTicket(ctx, id, req.Reason, req.CancelledBy); err != nil {
		h.logger.Error("Failed to cancel ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to cancel ticket: "+err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Ticket cancelled successfully"})
}

// AddComment handles POST /tickets/{id}/comments
func (h *TicketHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req app.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	req.TicketID = id

	if err := h.service.AddComment(ctx, req); err != nil {
		h.logger.Error("Failed to add comment", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to add comment")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Comment added successfully"})
}

// GetComments handles GET /tickets/{id}/comments
func (h *TicketHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	comments, err := h.service.GetComments(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get comments", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get comments")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"comments": comments})
}

// GetStatusHistory handles GET /tickets/{id}/history
func (h *TicketHandler) GetStatusHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	history, err := h.service.GetStatusHistory(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get status history", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get status history")
		return
	}

	h.respondJSON(w, http.StatusOK, history)
}

// GetTicketParts handles GET /tickets/{id}/parts
// Returns parts assigned to this specific ticket from ticket_parts table
func (h *TicketHandler) GetTicketParts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    if id == "" {
        h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
        return
    }

    if h.pool == nil {
        h.respondError(w, http.StatusInternalServerError, "DB pool not initialized")
        return
    }

    // Query ticket_parts table with join to spare_parts_catalog for details
    const q = `
        SELECT 
            tp.id,
            tp.spare_part_id,
            sp.part_number,
            sp.part_name,
            tp.quantity_required,
            tp.quantity_used,
            tp.is_critical,
            tp.status,
            tp.unit_price,
            tp.total_price,
            tp.currency,
            sp.category,
            sp.stock_status,
            sp.lead_time_days,
            tp.assigned_by,
            tp.assigned_at,
            tp.installed_at,
            tp.notes
        FROM ticket_parts tp
        JOIN spare_parts_catalog sp ON tp.spare_part_id = sp.id
        WHERE tp.ticket_id = $1
        ORDER BY tp.assigned_at DESC
    `

    rows, err := h.pool.Query(ctx, q, id)
    if err != nil {
        h.logger.Error("Failed to fetch ticket parts", slog.String("error", err.Error()))
        h.respondError(w, http.StatusInternalServerError, "Failed to fetch ticket parts")
        return
    }
    defer rows.Close()

    type Part struct {
        ID               string    `json:"id"` // ticket_parts.id for deletion
        AssignmentID     string    `json:"assignment_id"` // alias for backwards compatibility
        SparePartID      string    `json:"spare_part_id"`
        PartNumber       string    `json:"part_number"`
        PartName         string    `json:"part_name"`
        QuantityRequired int       `json:"quantity_required"`
        QuantityUsed     *int      `json:"quantity_used,omitempty"`
        IsCritical       bool      `json:"is_critical"`
        Status           string    `json:"status"`
        UnitPrice        *float64  `json:"unit_price,omitempty"`
        TotalPrice       *float64  `json:"total_price,omitempty"`
        Currency         string    `json:"currency"`
        Category         *string   `json:"category,omitempty"`
        StockStatus      *string   `json:"stock_status,omitempty"`
        LeadTimeDays     *int      `json:"lead_time_days,omitempty"`
        AssignedBy       *string   `json:"assigned_by,omitempty"`
        AssignedAt       time.Time `json:"assigned_at"`
        InstalledAt      *time.Time `json:"installed_at,omitempty"`
        Notes            *string   `json:"notes,omitempty"`
    }

    parts := make([]Part, 0, 8)
    for rows.Next() {
        var p Part
        var (
            qtyUsed     sql.NullInt64
            unitPrice   sql.NullFloat64
            totalPrice  sql.NullFloat64
            category    sql.NullString
            stockStatus sql.NullString
            leadTime    sql.NullInt64
            assignedBy  sql.NullString
            installedAt sql.NullTime
            notes       sql.NullString
        )
        
        if err := rows.Scan(
            &p.ID, &p.SparePartID, &p.PartNumber, &p.PartName,
            &p.QuantityRequired, &qtyUsed, &p.IsCritical, &p.Status,
            &unitPrice, &totalPrice, &p.Currency,
            &category, &stockStatus, &leadTime,
            &assignedBy, &p.AssignedAt, &installedAt, &notes,
        ); err != nil {
            h.logger.Warn("Failed to scan ticket part", slog.String("error", err.Error()))
            continue
        }
        
        // Set AssignmentID for backwards compatibility
        p.AssignmentID = p.ID
        
        if qtyUsed.Valid { v := int(qtyUsed.Int64); p.QuantityUsed = &v }
        if unitPrice.Valid { v := unitPrice.Float64; p.UnitPrice = &v }
        if totalPrice.Valid { v := totalPrice.Float64; p.TotalPrice = &v }
        if category.Valid { v := category.String; p.Category = &v }
        if stockStatus.Valid { v := stockStatus.String; p.StockStatus = &v }
        if leadTime.Valid { v := int(leadTime.Int64); p.LeadTimeDays = &v }
        if assignedBy.Valid { v := assignedBy.String; p.AssignedBy = &v }
        if installedAt.Valid { p.InstalledAt = &installedAt.Time }
        if notes.Valid { v := notes.String; p.Notes = &v }
        
        parts = append(parts, p)
    }

    h.respondJSON(w, http.StatusOK, map[string]any{
        "ticket_id": id,
        "count": len(parts),
        "parts": parts,
    })
}

// AddTicketPart handles POST /tickets/{id}/parts
// Adds a single part to the ticket_parts table
func (h *TicketHandler) AddTicketPart(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ticketID := chi.URLParam(r, "id")
    
    if ticketID == "" {
        h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
        return
    }
    
    if h.pool == nil {
        h.respondError(w, http.StatusInternalServerError, "DB pool not initialized")
        return
    }
    
    // Parse request body
    var req struct {
        SparePartID      string   `json:"spare_part_id"`
        QuantityRequired int      `json:"quantity_required"`
        UnitPrice        *float64 `json:"unit_price"`
        TotalPrice       *float64 `json:"total_price"`
        IsCritical       bool     `json:"is_critical"`
        Status           string   `json:"status"`
        Notes            string   `json:"notes"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }
    
    // Validate required fields
    if req.SparePartID == "" {
        h.respondError(w, http.StatusBadRequest, "spare_part_id is required")
        return
    }
    if req.QuantityRequired <= 0 {
        req.QuantityRequired = 1
    }
    if req.Status == "" {
        req.Status = "pending"
    }
    
    // Insert into ticket_parts table
    const insertQuery = `
        INSERT INTO ticket_parts (
            ticket_id, spare_part_id, quantity_required, 
            unit_price, total_price, is_critical, status, notes,
            assigned_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
        RETURNING id, assigned_at
    `
    
    var partID string
    var assignedAt time.Time
    err := h.pool.QueryRow(ctx, insertQuery,
        ticketID, req.SparePartID, req.QuantityRequired,
        req.UnitPrice, req.TotalPrice, req.IsCritical,
        req.Status, req.Notes,
    ).Scan(&partID, &assignedAt)
    
    if err != nil {
        h.logger.Error("Failed to add part to ticket",
            slog.String("ticket_id", ticketID),
            slog.String("spare_part_id", req.SparePartID),
            slog.String("error", err.Error()))
        h.respondError(w, http.StatusInternalServerError, "Failed to add part: "+err.Error())
        return
    }
    
    h.logger.Info("Part added to ticket",
        slog.String("ticket_id", ticketID),
        slog.String("part_id", partID))
    
    h.respondJSON(w, http.StatusCreated, map[string]interface{}{
        "id":                partID,
        "ticket_id":         ticketID,
        "spare_part_id":     req.SparePartID,
        "quantity_required": req.QuantityRequired,
        "unit_price":        req.UnitPrice,
        "total_price":       req.TotalPrice,
        "is_critical":       req.IsCritical,
        "status":            req.Status,
        "notes":             req.Notes,
        "assigned_at":       assignedAt,
    })
}

// DeleteTicketPart handles DELETE /tickets/{ticket_id}/parts/{part_id}
// Removes a specific part from the ticket
func (h *TicketHandler) DeleteTicketPart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	partID := chi.URLParam(r, "partId")
	
	if ticketID == "" || partID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID and Part ID are required")
		return
	}
	
	if h.pool == nil {
		h.respondError(w, http.StatusInternalServerError, "DB pool not initialized")
		return
	}
	
	// Delete the part from ticket_parts table
	const deleteQuery = `DELETE FROM ticket_parts WHERE id = $1 AND ticket_id = $2`
	
	result, err := h.pool.Exec(ctx, deleteQuery, partID, ticketID)
	if err != nil {
		h.logger.Error("Failed to delete ticket part", 
			slog.String("ticket_id", ticketID),
			slog.String("part_id", partID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to delete part")
		return
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Part not found")
		return
	}
	
	h.logger.Info("Ticket part deleted successfully",
		slog.String("ticket_id", ticketID),
		slog.String("part_id", partID))
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Part deleted successfully",
		"part_id": partID,
	})
}

// UpdateParts handles PATCH /tickets/{id}/parts
func (h *TicketHandler) UpdateParts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}

	var req struct {
		Parts []map[string]interface{} `json:"parts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.UpdateParts(ctx, id, req.Parts); err != nil {
		h.logger.Error("Failed to update parts", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update parts")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Parts updated successfully"})
}

// no-op: database/sql Null* used for scanning

// respondJSON writes JSON response
func (h *TicketHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes error response
func (h *TicketHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}


// DeleteComment handles DELETE /api/v1/tickets/{id}/comments/{commentId}
func (h *TicketHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")
	commentID := chi.URLParam(r, "commentId")

	if ticketID == "" || commentID == "" {
		http.Error(w, "ticket ID and comment ID are required", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteComment(ctx, ticketID, commentID); err != nil {
		http.Error(w, "failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Comment deleted"})
}

// UpdateTicketPriority handles PATCH /api/v1/tickets/{id}/priority (admin-only)
func (h *TicketHandler) UpdateTicketPriority(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()
	ticketID := chi.URLParam(r, "id")
	
	if ticketID == "" {
		h.respondError(w, http.StatusBadRequest, "Ticket ID is required")
		return
	}
	
	// Extract user info from JWT claims in context
	claims, ok := ctx.Value("claims").(map[string]interface{})
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "No authentication claims found")
		return
	}
	
	userRole, _ := claims["role"].(string)
	userID, _ := claims["user_id"].(string)
	orgType, _ := claims["organization_type"].(string)
	
	// Allow admins and system admins to update priority
	isAdmin := userRole == "admin" || userRole == "super_admin" || userRole == "system_admin"
	isSystemOrg := orgType == "system"
	
	if !isAdmin && !isSystemOrg {
		h.logger.Warn("Unauthorized priority update attempt",
			slog.String("ticket_id", ticketID),
			slog.String("user_role", userRole),
			slog.String("org_type", orgType))
		h.respondError(w, http.StatusForbidden, "Only admins can update ticket priority")
		return
	}
	
	var req struct {
		Priority string `json:"priority"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Validate priority
	validPriorities := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
	}
	
	if !validPriorities[req.Priority] {
		h.respondError(w, http.StatusBadRequest, "Invalid priority value. Must be: critical, high, medium, or low")
		return
	}
	
	// Get current ticket to log old priority
	ticket, err := h.service.GetTicket(ctx, ticketID)
	if err != nil {
		if err == domain.ErrTicketNotFound {
			h.respondError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		h.logger.Error("Failed to get ticket", slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to get ticket")
		return
	}
	
	oldPriority := string(ticket.Priority)
	
	// Update ticket priority
	if err := h.service.UpdatePriority(ctx, ticketID, domain.TicketPriority(req.Priority)); err != nil {
		h.logger.Error("Failed to update ticket priority",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to update priority")
		
		// Log failure to audit
		if h.auditLogger != nil {
			errorMsg := err.Error()
			durationMs := int(time.Since(startTime).Milliseconds())
			h.auditLogger.LogAsync(ctx, &audit.AuditEvent{
				EventType:     "ticket_priority_update_failed",
				EventCategory: audit.CategoryTicket,
				EventAction:   audit.ActionUpdate,
				EventStatus:   audit.StatusFailure,
				ResourceType:  stringPtr("ticket"),
				ResourceID:    &ticketID,
				IPAddress:     stringPtr(audit.ExtractIPAddress(r)),
				UserAgent:     stringPtr(audit.ExtractUserAgent(r)),
				RequestMethod: stringPtr(r.Method),
				RequestPath:   stringPtr(r.URL.Path),
				ErrorMessage:  &errorMsg,
				DurationMs:    &durationMs,
				Metadata: map[string]interface{}{
					"old_priority": oldPriority,
					"new_priority": req.Priority,
					"user_id":      userID,
					"user_role":    userRole,
				},
			})
		}
		return
	}
	
	// Log successful priority update to audit
	if h.auditLogger != nil {
		durationMs := int(time.Since(startTime).Milliseconds())
		ipAddress := audit.ExtractIPAddress(r)
		
		h.auditLogger.LogAsync(ctx, &audit.AuditEvent{
			EventType:     "ticket_priority_updated",
			EventCategory: audit.CategoryTicket,
			EventAction:   audit.ActionUpdate,
			EventStatus:   audit.StatusSuccess,
			ResourceType:  stringPtr("ticket"),
			ResourceID:    &ticketID,
			ResourceName:  stringPtr(ticket.TicketNumber),
			IPAddress:     &ipAddress,
			UserAgent:     stringPtr(audit.ExtractUserAgent(r)),
			RequestMethod: stringPtr(r.Method),
			RequestPath:   stringPtr(r.URL.Path),
			DurationMs:    &durationMs,
			Metadata: map[string]interface{}{
				"old_priority": oldPriority,
				"new_priority": req.Priority,
				"user_id":      userID,
				"user_role":    userRole,
			},
		})
	}
	
	h.logger.Info("Ticket priority updated",
		slog.String("ticket_id", ticketID),
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("old_priority", oldPriority),
		slog.String("new_priority", req.Priority),
		slog.String("updated_by", userID))
	
	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "success",
		"ticket_id":    ticketID,
		"old_priority": oldPriority,
		"new_priority": req.Priority,
		"message":      "Priority updated successfully",
	})
}

// SendEmailNotification handles POST /tickets/:id/send-notification
func (h *TicketHandler) SendEmailNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ticketID := chi.URLParam(r, "id")

	if h.notificationService == nil {
		h.respondError(w, http.StatusServiceUnavailable, "Notification service not available")
		return
	}

	var req struct {
		IncludeComments bool `json:"include_comments"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.IncludeComments = true // Default to including comments
	}

	err := h.notificationService.SendManualEmail(ctx, ticketID, req.IncludeComments)
	if err != nil {
		h.logger.Error("Failed to send email notification",
			slog.String("ticket_id", ticketID),
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusInternalServerError, "Failed to send email: "+err.Error())
		return
	}

	// Get ticket to return recipient info (optional for now)
	// ticket, _ := h.service.GetByID(ctx, ticketID)
	recipientEmail := "not-implemented"

	h.logger.Info("Email notification sent",
		slog.String("ticket_id", ticketID),
		slog.String("recipient", recipientEmail))

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"sent":      true,
		"recipient": recipientEmail,
		"message":   "Email notification sent successfully",
	})
}

// NotifyCustomer handles POST /tickets/:id/notify - Manual notification by admin
func (h *TicketHandler) NotifyCustomer(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")

	// Check feature flag
	if os.Getenv("FEATURE_MANUAL_NOTIFICATIONS") == "false" {
		h.respondError(w, http.StatusForbidden, "Manual notifications are disabled")
		return
	}

	var req struct {
		Email   string `json:"email"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Comment == "" {
		h.respondError(w, http.StatusBadRequest, "Email and comment are required")
		return
	}

	h.logger.Info("Manual notification requested",
		slog.String("ticket_id", ticketID),
		slog.String("email", req.Email),
		slog.String("comment_preview", req.Comment[:min(50, len(req.Comment))]))

	// Get ticket details for the email
	ticket, err := h.service.GetTicket(r.Context(), ticketID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Ticket not found")
		return
	}

	// Send email using SendGrid
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridAPIKey == "" {
		h.logger.Warn("SENDGRID_API_KEY not configured, email not sent")
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Notification logged (email service not configured)",
			"email":   req.Email,
		})
		return
	}

	// Send email via SendGrid
	emailService := emailInfra.NewNotificationService(
		sendgridAPIKey,
		os.Getenv("SENDGRID_FROM_EMAIL"),
		os.Getenv("SENDGRID_FROM_NAME"),
	)

	err = emailService.SendManualNotification(r.Context(), req.Email, ticket.CustomerName, ticket.TicketNumber, req.Comment)
	if err != nil {
		h.logger.Error("Failed to send notification email",
			slog.String("error", err.Error()),
			slog.String("email", req.Email))
		h.respondError(w, http.StatusInternalServerError, "Failed to send email: "+err.Error())
		return
	}

	h.logger.Info("Manual notification email sent successfully",
		slog.String("ticket_id", ticketID),
		slog.String("email", req.Email))

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Notification sent successfully",
		"email":   req.Email,
	})
}

// GetPublicTicket handles GET /track/:token
func (h *TicketHandler) GetPublicTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	if h.notificationService == nil {
		h.respondError(w, http.StatusServiceUnavailable, "Tracking service not available")
		return
	}

	publicView, err := h.notificationService.GetPublicTicketView(ctx, token)
	if err != nil {
		h.logger.Warn("Invalid tracking token access attempt",
			slog.String("token", token[:min(8, len(token))]), // Log first 8 chars only
			slog.String("error", err.Error()))
		h.respondError(w, http.StatusNotFound, "Invalid or expired tracking link")
		return
	}

	h.logger.Info("Public ticket viewed",
		slog.String("ticket_number", publicView.TicketNumber),
		slog.String("token", token[:min(8, len(token))]))

	h.respondJSON(w, http.StatusOK, publicView)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
