package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditLogger handles audit logging operations
type AuditLogger struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(db *pgxpool.Pool, logger *slog.Logger) *AuditLogger {
	return &AuditLogger{
		db:     db,
		logger: logger.With(slog.String("component", "audit_logger")),
	}
}

// EventCategory represents the high-level category of an event
type EventCategory string

const (
	CategoryAuth      EventCategory = "auth"
	CategoryEquipment EventCategory = "equipment"
	CategoryTicket    EventCategory = "ticket"
	CategoryEngineer  EventCategory = "engineer"
	CategoryParts     EventCategory = "parts"
	CategorySecurity  EventCategory = "security"
	CategoryOrg       EventCategory = "organization"
	CategorySystem    EventCategory = "system"
)

// EventAction represents the action type
type EventAction string

const (
	ActionCreate EventAction = "create"
	ActionRead   EventAction = "read"
	ActionUpdate EventAction = "update"
	ActionDelete EventAction = "delete"
	ActionView   EventAction = "view"
	ActionAssign EventAction = "assign"
	ActionScan   EventAction = "scan"
	ActionImport EventAction = "import"
	ActionExport EventAction = "export"
	ActionLogin  EventAction = "login"
	ActionLogout EventAction = "logout"
)

// EventStatus represents the outcome of an event
type EventStatus string

const (
	StatusSuccess EventStatus = "success"
	StatusFailure EventStatus = "failure"
	StatusDenied  EventStatus = "denied"
)

// AuditEvent represents a single audit log entry
type AuditEvent struct {
	// Event Information
	EventType     string        `json:"event_type"`
	EventCategory EventCategory `json:"event_category"`
	EventAction   EventAction   `json:"event_action"`
	EventStatus   EventStatus   `json:"event_status"`

	// User/Actor Information
	UserID           *string `json:"user_id,omitempty"`
	UserEmail        *string `json:"user_email,omitempty"`
	UserRole         *string `json:"user_role,omitempty"`
	OrganizationID   *string `json:"organization_id,omitempty"`
	OrganizationType *string `json:"organization_type,omitempty"`

	// Resource Information
	ResourceType *string `json:"resource_type,omitempty"`
	ResourceID   *string `json:"resource_id,omitempty"`
	ResourceName *string `json:"resource_name,omitempty"`

	// Request Information
	IPAddress     *string `json:"ip_address,omitempty"`
	UserAgent     *string `json:"user_agent,omitempty"`
	RequestMethod *string `json:"request_method,omitempty"`
	RequestPath   *string `json:"request_path,omitempty"`
	RequestQuery  *string `json:"request_query,omitempty"`

	// Change Tracking
	OldValues     map[string]interface{} `json:"old_values,omitempty"`
	NewValues     map[string]interface{} `json:"new_values,omitempty"`
	ChangedFields []string               `json:"changed_fields,omitempty"`

	// Additional Context
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ErrorMessage *string                `json:"error_message,omitempty"`
	DurationMs   *int                   `json:"duration_ms,omitempty"`

	// Rate Limiting Context
	IsRateLimited *bool   `json:"is_rate_limited,omitempty"`
	RateLimitKey  *string `json:"rate_limit_key,omitempty"`
}

// Log writes an audit event to the database
func (al *AuditLogger) Log(ctx context.Context, event *AuditEvent) error {
	query := `
		INSERT INTO audit_logs (
			event_type, event_category, event_action, event_status,
			user_id, user_email, user_role, organization_id, organization_type,
			resource_type, resource_id, resource_name,
			ip_address, user_agent, request_method, request_path, request_query,
			old_values, new_values, changed_fields,
			metadata, error_message, duration_ms,
			is_rate_limited, rate_limit_key
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20,
			$21, $22, $23,
			$24, $25
		)
	`

	// Convert maps to JSON
	var oldValuesJSON, newValuesJSON, metadataJSON []byte
	var err error

	if event.OldValues != nil {
		oldValuesJSON, err = json.Marshal(event.OldValues)
		if err != nil {
			al.logger.Error("Failed to marshal old_values", slog.Any("error", err))
			oldValuesJSON = nil
		}
	}

	if event.NewValues != nil {
		newValuesJSON, err = json.Marshal(event.NewValues)
		if err != nil {
			al.logger.Error("Failed to marshal new_values", slog.Any("error", err))
			newValuesJSON = nil
		}
	}

	if event.Metadata != nil {
		metadataJSON, err = json.Marshal(event.Metadata)
		if err != nil {
			al.logger.Error("Failed to marshal metadata", slog.Any("error", err))
			metadataJSON = nil
		}
	}

	_, err = al.db.Exec(ctx, query,
		event.EventType, event.EventCategory, event.EventAction, event.EventStatus,
		event.UserID, event.UserEmail, event.UserRole, event.OrganizationID, event.OrganizationType,
		event.ResourceType, event.ResourceID, event.ResourceName,
		event.IPAddress, event.UserAgent, event.RequestMethod, event.RequestPath, event.RequestQuery,
		oldValuesJSON, newValuesJSON, event.ChangedFields,
		metadataJSON, event.ErrorMessage, event.DurationMs,
		event.IsRateLimited, event.RateLimitKey,
	)

	if err != nil {
		al.logger.Error("Failed to write audit log",
			slog.String("event_type", event.EventType),
			slog.Any("error", err))
		return err
	}

	al.logger.Debug("Audit log written",
		slog.String("event_type", event.EventType),
		slog.String("status", string(event.EventStatus)))

	return nil
}

// LogAsync writes an audit event asynchronously (non-blocking)
func (al *AuditLogger) LogAsync(ctx context.Context, event *AuditEvent) {
	go func() {
		// Use a background context with timeout
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := al.Log(bgCtx, event); err != nil {
			al.logger.Error("Failed to write async audit log",
				slog.String("event_type", event.EventType),
				slog.Any("error", err))
		}
	}()
}

// Helper functions to extract information from HTTP requests

// ExtractIPAddress gets the client IP from the request
func ExtractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take first IP if multiple
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

// ExtractUserAgent gets the user agent from the request
func ExtractUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// Helper function to create event from request
func EventFromRequest(r *http.Request, eventType string, category EventCategory, action EventAction) *AuditEvent {
	ip := ExtractIPAddress(r)
	ua := ExtractUserAgent(r)
	method := r.Method
	path := r.URL.Path
	query := r.URL.RawQuery

	return &AuditEvent{
		EventType:     eventType,
		EventCategory: category,
		EventAction:   action,
		EventStatus:   StatusSuccess, // Default, can be changed
		IPAddress:     &ip,
		UserAgent:     &ua,
		RequestMethod: &method,
		RequestPath:   &path,
		RequestQuery:  &query,
	}
}

// Predefined event creators for common operations

// LogTicketCreated logs a ticket creation event
func (al *AuditLogger) LogTicketCreated(ctx context.Context, ticketID, qrCode, ipAddress string, metadata map[string]interface{}) {
	event := &AuditEvent{
		EventType:     "ticket_created",
		EventCategory: CategoryTicket,
		EventAction:   ActionCreate,
		EventStatus:   StatusSuccess,
		ResourceType:  stringPtr("ticket"),
		ResourceID:    &ticketID,
		IPAddress:     &ipAddress,
		Metadata:      metadata,
	}

	if qrCode != "" {
		event.Metadata["qr_code"] = qrCode
	}

	al.LogAsync(ctx, event)
}

// LogRateLimitExceeded logs a rate limit violation
func (al *AuditLogger) LogRateLimitExceeded(ctx context.Context, rateLimitKey, ipAddress string, metadata map[string]interface{}) {
	rateLimited := true
	event := &AuditEvent{
		EventType:     "rate_limit_exceeded",
		EventCategory: CategorySecurity,
		EventAction:   ActionCreate,
		EventStatus:   StatusDenied,
		IPAddress:     &ipAddress,
		IsRateLimited: &rateLimited,
		RateLimitKey:  &rateLimitKey,
		Metadata:      metadata,
	}

	al.LogAsync(ctx, event)
}

// LogQRScanned logs a QR code scan event
func (al *AuditLogger) LogQRScanned(ctx context.Context, qrCode, equipmentID, ipAddress string) {
	event := &AuditEvent{
		EventType:     "equipment_qr_scanned",
		EventCategory: CategoryEquipment,
		EventAction:   ActionScan,
		EventStatus:   StatusSuccess,
		ResourceType:  stringPtr("equipment"),
		ResourceID:    &equipmentID,
		IPAddress:     &ipAddress,
		Metadata: map[string]interface{}{
			"qr_code": qrCode,
		},
	}

	al.LogAsync(ctx, event)
}

// LogAuthLogin logs a login event
func (al *AuditLogger) LogAuthLogin(ctx context.Context, userID, userEmail, ipAddress string, success bool) {
	status := StatusSuccess
	if !success {
		status = StatusFailure
	}

	event := &AuditEvent{
		EventType:     "auth_login",
		EventCategory: CategoryAuth,
		EventAction:   ActionLogin,
		EventStatus:   status,
		UserID:        &userID,
		UserEmail:     &userEmail,
		IPAddress:     &ipAddress,
	}

	al.LogAsync(ctx, event)
}

// LogAuthLogout logs a logout event
func (al *AuditLogger) LogAuthLogout(ctx context.Context, userID, userEmail, ipAddress string) {
	event := &AuditEvent{
		EventType:     "auth_logout",
		EventCategory: CategoryAuth,
		EventAction:   ActionLogout,
		EventStatus:   StatusSuccess,
		UserID:        &userID,
		UserEmail:     &userEmail,
		IPAddress:     &ipAddress,
	}

	al.LogAsync(ctx, event)
}

// LogEquipmentCreated logs equipment creation
func (al *AuditLogger) LogEquipmentCreated(ctx context.Context, equipmentID, equipmentName, userID string, orgID *string) {
	event := &AuditEvent{
		EventType:      "equipment_created",
		EventCategory:  CategoryEquipment,
		EventAction:    ActionCreate,
		EventStatus:    StatusSuccess,
		UserID:         &userID,
		OrganizationID: orgID,
		ResourceType:   stringPtr("equipment"),
		ResourceID:     &equipmentID,
		ResourceName:   &equipmentName,
	}

	al.LogAsync(ctx, event)
}

// LogEngineerAssigned logs engineer assignment to ticket
func (al *AuditLogger) LogEngineerAssigned(ctx context.Context, ticketID, engineerID, assignedByUserID string, orgID *string) {
	event := &AuditEvent{
		EventType:      "ticket_engineer_assigned",
		EventCategory:  CategoryTicket,
		EventAction:    ActionAssign,
		EventStatus:    StatusSuccess,
		UserID:         &assignedByUserID,
		OrganizationID: orgID,
		ResourceType:   stringPtr("ticket"),
		ResourceID:     &ticketID,
		Metadata: map[string]interface{}{
			"engineer_id": engineerID,
		},
	}

	al.LogAsync(ctx, event)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
