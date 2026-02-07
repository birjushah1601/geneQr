package serviceticket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	equipmentInfra "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/infra"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/qrcode"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/api"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/app"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/infra"
	"github.com/aby-med/medical-platform/internal/service-domain/whatsapp"
	"github.com/aby-med/medical-platform/internal/shared/audit"
	sharedMiddleware "github.com/aby-med/medical-platform/internal/shared/middleware"
	attachmentDomain "github.com/aby-med/medical-platform/internal/service-domain/attachment/domain"
	attachmentInfra "github.com/aby-med/medical-platform/internal/service-domain/attachment/infra"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module represents the service ticket module
type Module struct {
	config                     ModuleConfig
	ticketHandler              *api.TicketHandler
	assignmentHandler          *api.AssignmentHandler
	multiModelAssignmentHandler *api.MultiModelAssignmentHandler
	whatsappHandler            *whatsapp.WebhookHandler
	logger                     *slog.Logger
	dispatcher                 *app.WebhookDispatcher
	slaMonitor                 *app.SLAMonitor
	qrRateLimiter              *sharedMiddleware.QRRateLimiter
	ipRateLimiter              *sharedMiddleware.IPRateLimiter
	inputSanitizer             *sharedMiddleware.InputSanitizer
	auditLogger                *audit.AuditLogger
}

// ModuleConfig holds module configuration
type ModuleConfig struct {
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	BaseURL            string
	QROutputDir        string
	WhatsAppVerifyToken string
	WhatsAppAccessToken string
	WhatsAppPhoneID     string
	WhatsAppMediaDir    string
}

// NewModule creates a new service ticket module
func NewModule(cfg ModuleConfig, logger *slog.Logger) (*Module, error) {
	return &Module{
		config: cfg,
		logger: logger.With(slog.String("module", "service-ticket")),
	}, nil
}

// Initialize initializes the module (database connections, etc.)
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing Service Ticket module")

	// Create database connection pool
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		m.config.DBUser, m.config.DBPassword, m.config.DBHost, m.config.DBPort, m.config.DBName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

    // Ensure schema exists
    if err := infra.EnsureServiceTicketSchema(ctx, pool); err != nil {
        return fmt.Errorf("failed to ensure service ticket schema: %w", err)
    }

	// Create repositories
	ticketRepo := infra.NewTicketRepository(pool)
	assignmentRepo := infra.NewAssignmentRepository(pool)
	notificationRepo := infra.NewNotificationRepository(pool)
	
	// Create equipment repository for WhatsApp integration
	equipmentRepo := equipmentInfra.NewEquipmentRepository(pool)

    // Create ticket service
    policyRepo := infra.NewPolicyRepository(pool)
    eventRepo := infra.NewEventRepository(pool)
    ticketService := app.NewTicketService(ticketRepo, equipmentRepo, policyRepo, eventRepo, m.logger)
	
	// Create notification service
	// TODO: Replace nil with actual email service when configured
	notificationService := app.NewNotificationService(
		ticketRepo,
		notificationRepo,
		nil, // emailService - to be added when SendGrid is configured
		m.logger,
		&app.NotificationConfig{
			BaseURL:         m.config.BaseURL,
			TokenExpiryDays: 30,
		},
	)

	// Create assignment service
	assignmentService := app.NewAssignmentService(assignmentRepo, ticketRepo, pool, m.logger)
	
	// Create multi-model assignment service
	multiModelService := app.NewMultiModelAssignmentService(assignmentRepo, ticketRepo, equipmentRepo, pool, m.logger)

    // Create dispatcher (started conditionally)
    m.dispatcher = app.NewWebhookDispatcher(pool, m.logger)
    // Create SLA monitor (started conditionally)
    m.slaMonitor = app.NewSLAMonitor(pool, m.logger)

	// Initialize audit logger first (needed by handlers)
	m.auditLogger = audit.NewAuditLogger(pool, m.logger)
	m.logger.Info("Audit logger initialized")
	
    // Create HTTP handlers (with audit logger)
    m.ticketHandler = api.NewTicketHandler(ticketService, m.logger, pool, m.auditLogger)
	m.ticketHandler.SetNotificationService(notificationService)
	m.logger.Info("Notification service wired to ticket handler")
	
	m.assignmentHandler = api.NewAssignmentHandler(assignmentService, m.logger)
	m.multiModelAssignmentHandler = api.NewMultiModelAssignmentHandler(multiModelService, m.logger)

	// Create QR generator for WhatsApp
	qrGenerator := qrcode.NewGenerator(m.config.BaseURL, m.config.QROutputDir)

	// Initialize minimal AttachmentService for WhatsApp intake (using same DB pool)
	attRepo := attachmentInfra.NewPostgresAttachmentRepository(pool)
	queueRepo := attachmentInfra.NewPostgresProcessingQueueRepository(pool)
	aiRepo := attachmentInfra.NewNoopAIAnalysisRepository()
	attService := attachmentDomain.NewAttachmentService(attRepo, queueRepo, aiRepo, m.logger)

	// Create WhatsApp webhook handler
	whatsappConfig := whatsapp.WebhookConfig{
		VerifyToken:   m.config.WhatsAppVerifyToken,
		AccessToken:   m.config.WhatsAppAccessToken,
		PhoneNumberID: m.config.WhatsAppPhoneID,
		MediaDir:      m.config.WhatsAppMediaDir,
	}
	m.whatsappHandler = whatsapp.NewWebhookHandler(whatsappConfig, qrGenerator, ticketService, m.logger, attService)

	// Initialize QR-based rate limiter
	// 5 tickets per QR code per hour (prevents spam while allowing legitimate use)
	m.qrRateLimiter = sharedMiddleware.NewQRRateLimiter(5, 1*time.Hour, m.logger)
	m.logger.Info("QR rate limiter initialized", slog.Int("limit", 5), slog.String("window", "1 hour"))

	// Initialize IP-based rate limiter
	// 20 tickets per IP per hour (prevents multi-QR spam attacks)
	m.ipRateLimiter = sharedMiddleware.NewIPRateLimiter(20, 1*time.Hour, m.logger)
	m.logger.Info("IP rate limiter initialized", slog.Int("limit", 20), slog.String("window", "1 hour"))

	// Initialize input sanitizer
	m.inputSanitizer = sharedMiddleware.NewInputSanitizer(m.logger)
	m.logger.Info("Input sanitizer initialized")

	m.logger.Info("Service Ticket module initialized successfully")
	return nil
}

// MountRoutes mounts module routes to the router
func (m *Module) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting Service Ticket routes")

	// QR code extractor for rate limiting
	qrCodeExtractor := func(r *http.Request) string {
		var body map[string]interface{}
		var bodyBytes []byte
		
		// Read body for QR code
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			json.Unmarshal(bodyBytes, &body)
			// Restore body with the actual content
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		
		if qrCode, ok := body["QRCode"].(string); ok {
			return qrCode
		}
		if qrCode, ok := body["qr_code"].(string); ok {
			return qrCode
		}
		return ""
	}

	// Service ticket management routes
	r.Route("/tickets", func(r chi.Router) {
		// Apply multiple security middlewares to ticket creation
		// 1. Input sanitization (strip HTML/scripts, limit sizes)
		// 2. IP rate limiting (20 tickets per IP per hour)
		// 3. QR rate limiting (5 tickets per QR code per hour)
		r.With(
			m.inputSanitizer.Middleware(sharedMiddleware.DefaultSanitizeConfig()),
			m.ipRateLimiter.Middleware,
			m.qrRateLimiter.Middleware(qrCodeExtractor),
		).Post("/", m.ticketHandler.CreateTicket)
		r.Get("/", m.ticketHandler.ListTickets)                // List tickets
		r.Get("/number/{number}", m.ticketHandler.GetTicketByNumber) // Get by ticket number
		r.Get("/{id}", m.ticketHandler.GetTicket)              // Get by ID
        r.Get("/{id}/parts", m.ticketHandler.GetTicketParts)       // Get parts for ticket
		r.Post("/{id}/parts", m.ticketHandler.AddTicketPart)       // Add single part to ticket
		r.Delete("/{id}/parts/{partId}", m.ticketHandler.DeleteTicketPart) // Delete specific part
		r.Patch("/{id}/parts", m.ticketHandler.UpdateParts)        // Update parts for ticket
		
		// Admin-only: Update ticket priority
		// NOTE: In production, add proper JWT auth middleware to verify admin role
		// For now, checks X-User-Role header in handler
		r.Patch("/{id}/priority", m.ticketHandler.UpdateTicketPriority)
		
		// Ticket lifecycle operations
		r.Post("/{id}/assign", m.ticketHandler.AssignTicket)       // Assign engineer (legacy)
		r.Post("/{id}/acknowledge", m.ticketHandler.AcknowledgeTicket) // Acknowledge
		r.Post("/{id}/start", m.ticketHandler.StartWork)           // Start work
		r.Post("/{id}/hold", m.ticketHandler.PutOnHold)            // Put on hold
		r.Post("/{id}/resume", m.ticketHandler.ResumeWork)         // Resume work
		r.Post("/{id}/resolve", m.ticketHandler.ResolveTicket)     // Resolve
		r.Post("/{id}/close", m.ticketHandler.CloseTicket)         // Close
		r.Post("/{id}/cancel", m.ticketHandler.CancelTicket)       // Cancel
		
		// Assignment operations (new)
		r.Get("/{id}/suggested-engineers", m.assignmentHandler.GetSuggestedEngineers) // Get engineer suggestions (legacy)
		r.Get("/{id}/assignment-suggestions", m.multiModelAssignmentHandler.GetAssignmentSuggestions) // Multi-model suggestions
		r.Post("/{id}/assign-engineer", m.assignmentHandler.AssignEngineer)           // Manual assignment with tier
		r.Get("/{id}/assignments/history", m.assignmentHandler.GetAssignmentHistory)  // Get assignment history
		
		// Comments and history
		r.Post("/{id}/comments", m.ticketHandler.AddComment)       // Add comment
		r.Get("/{id}/comments", m.ticketHandler.GetComments)       // Get comments
		r.Delete("/{id}/comments/{commentId}", m.ticketHandler.DeleteComment) // Delete comment
		r.Get("/{id}/history", m.ticketHandler.GetStatusHistory)   // Get status history
		
		// Notification routes
		r.Post("/{id}/send-notification", m.ticketHandler.SendEmailNotification) // Send manual email (auto)
		r.Post("/{id}/notify", m.ticketHandler.NotifyCustomer)                   // Manual notification by admin
	})
	
	// Public tracking route (no auth required)
	r.Get("/track/{token}", m.ticketHandler.GetPublicTicket)

	// Engineer management routes
	r.Route("/engineers", func(r chi.Router) {
		r.Get("/", m.assignmentHandler.ListEngineers)           // List all engineers
		r.Post("/import", m.assignmentHandler.ImportEngineersCSV) // CSV import
		r.Get("/{id}", m.assignmentHandler.GetEngineer)         // Get engineer details
		r.Put("/{id}/level", m.assignmentHandler.UpdateEngineerLevel) // Update engineer level
		
		// Engineer equipment type capabilities
		r.Get("/{id}/equipment-types", m.assignmentHandler.ListEngineerEquipmentTypes)    // List capabilities
		r.Post("/{id}/equipment-types", m.assignmentHandler.AddEngineerEquipmentType)     // Add capability
		r.Delete("/{id}/equipment-types", m.assignmentHandler.RemoveEngineerEquipmentType) // Remove capability
	})

	// Equipment service configuration routes (under service-tickets to avoid conflict)
	r.Route("/equipment-service-config", func(r chi.Router) {
		r.Get("/{id}", m.assignmentHandler.GetEquipmentServiceConfig)    // Get config
		r.Post("/{id}", m.assignmentHandler.CreateEquipmentServiceConfig) // Create config
		r.Put("/{id}", m.assignmentHandler.UpdateEquipmentServiceConfig)  // Update config
	})

	// Note: Organization-specific engineer routes removed to avoid conflict with organizations module
	// Use /engineers?orgId={orgId} instead to filter engineers by organization

	// WhatsApp webhook routes
	r.Route("/whatsapp", func(r chi.Router) {
		r.Get("/webhook", m.whatsappHandler.VerifyWebhook)  // Webhook verification
		r.Post("/webhook", m.whatsappHandler.HandleWebhook) // Webhook handler
	})

	m.logger.Info("Service Ticket routes mounted successfully")
}

// Start starts background tasks (if any)
func (m *Module) Start(ctx context.Context) error {
    m.logger.Info("Service Ticket module started")
    // Start dispatcher if enabled
    if m.dispatcher != nil {
        go m.dispatcher.Run(ctx)
    }
    // Start SLA monitor if enabled
    if m.slaMonitor != nil {
        go m.slaMonitor.Run(ctx)
    }
	return nil
}

// Stop gracefully stops the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Service Ticket module stopped")
	return nil
}

// Health returns the health status
func (m *Module) Health(ctx context.Context) error {
	// TODO: Check database connection
	return nil
}

// Name returns the module name
func (m *Module) Name() string {
	return "service-ticket"
}
