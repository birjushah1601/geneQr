package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	equipmentApp "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	ticketDomain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	ticketApp "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
)

// WhatsAppMessage represents an incoming WhatsApp message
type WhatsAppMessage struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`       // Customer WhatsApp number
	To        string    `json:"to"`         // Business WhatsApp number
	Text      string    `json:"text"`       // Message content
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // text, image, document, location
	MediaURL  string    `json:"media_url,omitempty"`
}

// WhatsAppWebhookPayload is the incoming webhook payload
type WhatsAppWebhookPayload struct {
	Event   string           `json:"event"` // "message" or "status"
	Message *WhatsAppMessage `json:"message,omitempty"`
	Status  *MessageStatus   `json:"status,omitempty"`
}

// MessageStatus represents message delivery status
type MessageStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // sent, delivered, read, failed
	Timestamp time.Time `json:"timestamp"`
}

// TicketCreatedResponse is sent back to customer
type TicketCreatedResponse struct {
	Success      bool   `json:"success"`
	TicketNumber string `json:"ticket_number"`
	Message      string `json:"message"`
}

// WhatsAppHandler handles WhatsApp webhook events
type WhatsAppHandler struct {
	equipmentService *equipmentApp.EquipmentService
	ticketService    *ticketApp.ServiceTicketService
	whatsappService  *WhatsAppService // For sending messages back
	logger           *slog.Logger
}

// NewWhatsAppHandler creates a new WhatsApp webhook handler
func NewWhatsAppHandler(
	equipmentService *equipmentApp.EquipmentService,
	ticketService *ticketApp.ServiceTicketService,
	whatsappService *WhatsAppService,
	logger *slog.Logger,
) *WhatsAppHandler {
	return &WhatsAppHandler{
		equipmentService: equipmentService,
		ticketService:    ticketService,
		whatsappService:  whatsappService,
		logger:           logger.With(slog.String("component", "whatsapp_handler")),
	}
}

// HandleWebhook processes incoming WhatsApp webhooks
func (h *WhatsAppHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse webhook payload
	var payload WhatsAppWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Error("Failed to parse webhook payload", slog.String("error", err.Error()))
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Handle based on event type
	switch payload.Event {
	case "message":
		if payload.Message != nil {
			h.handleIncomingMessage(ctx, payload.Message)
		}
	case "status":
		if payload.Status != nil {
			h.handleMessageStatus(ctx, payload.Status)
		}
	default:
		h.logger.Warn("Unknown event type", slog.String("event", payload.Event))
	}

	// Acknowledge receipt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleIncomingMessage processes incoming WhatsApp messages
func (h *WhatsAppHandler) handleIncomingMessage(ctx context.Context, msg *WhatsAppMessage) {
	h.logger.Info("Received WhatsApp message",
		slog.String("from", msg.From),
		slog.String("text", msg.Text),
	)

	// Extract QR code from message
	qrCode := h.extractQRCode(msg.Text)
	
	if qrCode == "" {
		// No QR code found - send help message
		h.sendHelpMessage(ctx, msg.From)
		return
	}

	// Lookup equipment by QR code
	equipment, err := h.equipmentService.GetEquipmentByQR(ctx, qrCode)
	if err != nil {
		h.logger.Error("Equipment not found",
			slog.String("qr_code", qrCode),
			slog.String("error", err.Error()),
		)
		h.sendErrorMessage(ctx, msg.From, "Equipment not found. Please check the QR code and try again.")
		return
	}

	// Extract issue description
	issueDescription := h.extractIssueDescription(msg.Text, qrCode)
	
	// Determine priority based on keywords
	priority := h.determinePriority(issueDescription)

	// Create service ticket
	ticket, err := h.createTicketFromWhatsApp(ctx, equipment, msg, issueDescription, priority)
	if err != nil {
		h.logger.Error("Failed to create ticket",
			slog.String("error", err.Error()),
		)
		h.sendErrorMessage(ctx, msg.From, "Failed to create service ticket. Please try again or call support.")
		return
	}

	h.logger.Info("Ticket created from WhatsApp",
		slog.String("ticket_number", ticket.TicketNumber),
		slog.String("equipment_id", equipment.ID),
		slog.String("customer_phone", msg.From),
	)

	// Send confirmation to customer
	h.sendTicketConfirmation(ctx, msg.From, ticket)
	
	// TODO: Notify admin dashboard via WebSocket
	// TODO: Send notification to available engineers in the area
}

// extractQRCode extracts QR code from message text
func (h *WhatsAppHandler) extractQRCode(text string) string {
	// Try to find QR code pattern: QR-YYYYMMDD-XXXXXX
	qrPattern := regexp.MustCompile(`QR-\d{8}-\d{6}`)
	matches := qrPattern.FindString(text)
	if matches != "" {
		return matches
	}

	// Also check for standalone alphanumeric codes
	// In case user sends just "QR20251001832300"
	qrPatternAlt := regexp.MustCompile(`QR\d{14}`)
	matchesAlt := qrPatternAlt.FindString(text)
	if matchesAlt != "" {
		// Convert QR20251001832300 to QR-20251001-832300
		return fmt.Sprintf("QR-%s-%s", matchesAlt[2:10], matchesAlt[10:])
	}

	return ""
}

// extractIssueDescription extracts the issue description from message
func (h *WhatsAppHandler) extractIssueDescription(text, qrCode string) string {
	// Remove QR code from text
	description := strings.ReplaceAll(text, qrCode, "")
	
	// Clean up
	description = strings.TrimSpace(description)
	
	if description == "" {
		return "Equipment issue reported via WhatsApp"
	}
	
	return description
}

// determinePriority determines ticket priority based on keywords
func (h *WhatsAppHandler) determinePriority(description string) ticketDomain.TicketPriority {
	descLower := strings.ToLower(description)
	
	// Critical keywords
	criticalKeywords := []string{"urgent", "emergency", "critical", "down", "not working", "stopped", "patient"}
	for _, keyword := range criticalKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityCritical
		}
	}
	
	// High priority keywords
	highKeywords := []string{"error", "alarm", "warning", "issue", "problem", "broken"}
	for _, keyword := range highKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityHigh
		}
	}
	
	// Medium priority keywords
	mediumKeywords := []string{"maintenance", "service", "check", "noise", "slow"}
	for _, keyword := range mediumKeywords {
		if strings.Contains(descLower, keyword) {
			return ticketDomain.PriorityMedium
		}
	}
	
	return ticketDomain.PriorityMedium
}

// createTicketFromWhatsApp creates a service ticket from WhatsApp message
func (h *WhatsAppHandler) createTicketFromWhatsApp(
	ctx context.Context,
	equipment *domain.Equipment,
	msg *WhatsAppMessage,
	issueDescription string,
	priority ticketDomain.TicketPriority,
) (*ticketDomain.ServiceTicket, error) {
	
	// Create ticket request
	// NOTE: Engineer assignment should filter by equipment.ManufacturerID
	// Only engineers belonging to the equipment's manufacturer can be assigned
	// This ensures manufacturer-specific service teams handle their own equipment
	req := ticketApp.CreateTicketRequest{
		EquipmentID:      equipment.ID,
		QRCode:           equipment.QRCode,
		SerialNumber:     equipment.SerialNumber,
		CustomerPhone:    msg.From,
		CustomerWhatsApp: msg.From,
		IssueCategory:    "breakdown", // Could be enhanced with keyword detection
		IssueDescription: issueDescription,
		Priority:         priority,
		Source:           ticketDomain.SourceWhatsApp,
		SourceMessageID:  msg.ID,
		CreatedBy:        "whatsapp-bot",
	}

	// Create ticket
	ticket, err := h.ticketService.CreateTicket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return ticket, nil
}

// sendTicketConfirmation sends ticket confirmation to customer
func (h *WhatsAppHandler) sendTicketConfirmation(ctx context.Context, to string, ticket *ticketDomain.ServiceTicket) {
	message := fmt.Sprintf(
		"✅ *Service Request Confirmed*\n\n"+
			"Ticket Number: *%s*\n"+
			"Equipment: %s\n"+
			"Serial: %s\n"+
			"Priority: %s\n\n"+
			"Our engineer will contact you soon.\n"+
			"Thank you!",
		ticket.TicketNumber,
		ticket.EquipmentName,
		ticket.SerialNumber,
		strings.ToUpper(string(ticket.Priority)),
	)

	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send confirmation",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// sendErrorMessage sends error message to customer
func (h *WhatsAppHandler) sendErrorMessage(ctx context.Context, to, errorMsg string) {
	message := fmt.Sprintf("❌ %s", errorMsg)
	
	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send error message",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// sendHelpMessage sends help message when no QR code is found
func (h *WhatsAppHandler) sendHelpMessage(ctx context.Context, to string) {
	message := "🤖 *ABY-MED Service Bot*\n\n" +
		"To report an equipment issue:\n\n" +
		"1️⃣ Scan the QR code on your equipment\n" +
		"2️⃣ Send the QR code along with issue description\n\n" +
		"Example:\n" +
		"QR-20251001-832300\n" +
		"MRI machine not starting, showing error code E-503\n\n" +
		"We'll create a ticket and assign an engineer immediately!"

	if err := h.whatsappService.SendMessage(ctx, to, message); err != nil {
		h.logger.Error("Failed to send help message",
			slog.String("to", to),
			slog.String("error", err.Error()),
		)
	}
}

// handleMessageStatus handles message delivery status updates
func (h *WhatsAppHandler) handleMessageStatus(ctx context.Context, status *MessageStatus) {
	h.logger.Info("Message status update",
		slog.String("message_id", status.ID),
		slog.String("status", status.Status),
	)
	
	// TODO: Update ticket with delivery status if needed
}
