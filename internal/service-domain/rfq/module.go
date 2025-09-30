package rfq

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aby-med/medical-platform/internal/service-domain/rfq/api"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq/app"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/rfq/infra"
	"github.com/aby-med/medical-platform/internal/shared/service"
	"github.com/go-chi/chi/v5"
)

// Module represents the RFQ service module
type Module struct {
	config      *Config
	logger      *slog.Logger
	db          *infra.PostgresDB
	repository  domain.RFQRepository
	eventBus    domain.EventPublisher
	appService  *app.RFQService
	httpHandler *api.RFQHandler
}

// Config holds configuration for the RFQ module
type Config struct {
	DatabaseDSN  string
	KafkaBrokers []string
}

// NewModule creates a new RFQ module instance
func NewModule(config *Config, logger *slog.Logger) *Module {
	return &Module{
		config: config,
		logger: logger.With(slog.String("module", "rfq")),
	}
}

// Name returns the module name
func (m *Module) Name() string {
	return "rfq"
}

// Initialize sets up the module dependencies
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing RFQ module")

	// Initialize database connection
	db, err := infra.NewPostgresDB(ctx, m.config.DatabaseDSN, m.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Initialize repository
	m.repository = infra.NewRFQRepository(db, m.logger)

	// Initialize event bus
	m.eventBus = infra.NewKafkaEventPublisher(m.config.KafkaBrokers, m.logger)

	// Initialize application service
	m.appService = app.NewRFQService(m.repository, m.eventBus, m.logger)

	// Initialize HTTP handler
	m.httpHandler = api.NewRFQHandler(m.appService, m.logger)

	m.logger.Info("RFQ module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes with the router
func (m *Module) MountRoutes(router chi.Router) {
	m.logger.Info("Mounting RFQ routes")

	// Create subrouter with tenant context middleware
	router.Route("/rfq", func(r chi.Router) {
		r.Use(m.tenantContextMiddleware)

		// RFQ endpoints
		r.Get("/", m.httpHandler.ListRFQs)
		r.Post("/", m.httpHandler.CreateRFQ)
		r.Get("/{id}", m.httpHandler.GetRFQ)
		r.Put("/{id}", m.httpHandler.UpdateRFQ)
		r.Delete("/{id}", m.httpHandler.DeleteRFQ)

		// RFQ action endpoints
		r.Post("/{id}/publish", m.httpHandler.PublishRFQ)
		r.Post("/{id}/close", m.httpHandler.CloseRFQ)
		r.Post("/{id}/cancel", m.httpHandler.CancelRFQ)

		// RFQ items endpoints
		r.Post("/{id}/items", m.httpHandler.AddItem)
		r.Delete("/{id}/items/{item_id}", m.httpHandler.RemoveItem)
	})

	m.logger.Info("RFQ routes mounted successfully")
}

// Start starts the module background services
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting RFQ module")
	// No background services to start for now
	return nil
}

// Stop gracefully stops the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Stopping RFQ module")

	// Close database connection
	if m.db != nil {
		m.db.Close()
	}

	// Close event bus
	if m.eventBus != nil {
		if publisher, ok := m.eventBus.(*infra.KafkaEventPublisher); ok {
			if err := publisher.Close(); err != nil {
				m.logger.Error("Failed to close event publisher",
					slog.String("error", err.Error()))
			}
		}
	}

	m.logger.Info("RFQ module stopped")
	return nil
}

// tenantContextMiddleware extracts tenant ID from request header and adds it to context
func (m *Module) tenantContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			// For development, use a default tenant
			tenantID = "city-hospital"
			m.logger.Warn("No tenant ID provided, using default",
				slog.String("tenant_id", tenantID))
		}

		// Add tenant ID to context
		ctx := domain.WithTenantID(r.Context(), tenantID)

		// Try to extract user ID from header (if available)
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			ctx = domain.WithUserID(ctx, userID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Ensure Module implements service.Module interface
var _ service.Module = (*Module)(nil)
