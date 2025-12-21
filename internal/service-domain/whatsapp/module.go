package whatsapp

import (
	"log/slog"
	"net/http"
	
	equipmentApp "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/app"
	ticketApp "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WhatsAppModule encapsulates WhatsApp integration
type WhatsAppModule struct {
	handler *WhatsAppHandler
	logger  *slog.Logger
}

// NewWhatsAppModule creates a new WhatsApp module with all dependencies
func NewWhatsAppModule(
	db *pgxpool.Pool,
	equipmentService *equipmentApp.EquipmentService,
	ticketService *ticketApp.TicketService,
	twilioAccountSID string,
	twilioAuthToken string,
	twilioWhatsAppNumber string,
	logger *slog.Logger,
) *WhatsAppModule {
	logger.Info("Initializing WhatsApp module",
		slog.String("whatsapp_number", maskPhoneNumber(twilioWhatsAppNumber)))
	
	// Create WhatsApp service for sending messages
	whatsappService := NewWhatsAppService(
		twilioAccountSID,
		twilioAuthToken,
		twilioWhatsAppNumber,
		logger,
	)
	
	// Create handler for incoming webhooks
	handler := NewWhatsAppHandler(
		equipmentService,
		ticketService,
		whatsappService,
		logger,
	)
	
	logger.Info("WhatsApp module initialized successfully")
	
	return &WhatsAppModule{
		handler: handler,
		logger:  logger,
	}
}

// MountRoutes registers WhatsApp webhook routes
func (m *WhatsAppModule) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting WhatsApp routes")
	
	// Webhook for incoming messages
	r.Post("/whatsapp/webhook", m.handler.HandleWebhook)
	
	// Verification endpoint (required by Twilio)
	r.Get("/whatsapp/webhook", m.handleVerification)
	
	m.logger.Info("WhatsApp routes mounted",
		slog.String("webhook_path", "/api/v1/whatsapp/webhook"))
}

// handleVerification handles Twilio webhook verification
func (m *WhatsAppModule) handleVerification(w http.ResponseWriter, r *http.Request) {
	// Twilio sends a GET request to verify the webhook URL
	// We just need to return 200 OK
	m.logger.Info("WhatsApp webhook verification request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WhatsApp webhook verified"))
}
