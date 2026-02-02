package catalog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/marketplace/catalog/app"
	"github.com/aby-med/medical-platform/internal/marketplace/catalog/domain"
	"github.com/aby-med/medical-platform/internal/marketplace/catalog/infra"
	cataloghttp "github.com/aby-med/medical-platform/internal/marketplace/catalog/http"
	"github.com/aby-med/medical-platform/internal/shared/config"
	"github.com/go-chi/chi/v5"
)

const (
	// ModuleName is the unique identifier for this module
	ModuleName = "catalog"
)

// Module represents the catalog module that manages medical equipment catalog
type Module struct {
	config     *config.Config
	logger     *slog.Logger
	repository domain.CatalogRepository
	service    *app.CatalogService
	handler    *cataloghttp.Handler
	eventBus   domain.EventPublisher
}

// New creates a new catalog module
func New(cfg *config.Config, logger *slog.Logger) *Module {
	moduleLogger := logger.With(slog.String("module", ModuleName))
	
	return &Module{
		config: cfg,
		logger: moduleLogger,
	}
}

// Name returns the module's unique identifier
func (m *Module) Name() string {
	return ModuleName
}

// MountRoutes registers the module's HTTP routes on the provided router
func (m *Module) MountRoutes(r chi.Router) {
	if m.handler == nil {
		m.logger.Error("Cannot mount routes: handler not initialized")
		return
	}

	m.logger.Info("Mounting catalog routes")
	
	// Create a subrouter for catalog endpoints
	catalogRouter := chi.NewRouter()
	
	// Mount equipment catalog endpoints
	catalogRouter.Get("/", m.handler.ListEquipment)
	catalogRouter.Get("/{id}", m.handler.GetEquipment)
	catalogRouter.Post("/", m.handler.CreateEquipment)
	catalogRouter.Put("/{id}", m.handler.UpdateEquipment)
	catalogRouter.Delete("/{id}", m.handler.DeleteEquipment)
	catalogRouter.Get("/search", m.handler.SearchEquipment)
	catalogRouter.Get("/categories", m.handler.ListCategories)
	catalogRouter.Get("/manufacturers", m.handler.ListManufacturers)
	
	// Mount the catalog router under /catalog path
	r.Mount("/catalog", catalogRouter)
}

// Initialize sets up the module's dependencies (called before MountRoutes)
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing catalog module")
	
	// Initialize repositories
	db, err := infra.NewPostgresDB(ctx, m.config.GetDSN(), m.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize catalog database: %w", err)
	}
	
	// Initialize repositories
	m.repository = infra.NewCatalogRepository(db, m.logger)
	
	// Initialize event publisher
	m.eventBus = infra.NewKafkaEventPublisher(m.config.Kafka.Brokers, m.logger)
	
	// Initialize application services
	m.service = app.NewCatalogService(m.repository, m.eventBus, m.logger)
	
	// Initialize HTTP handlers (must be done before MountRoutes)
	m.handler = cataloghttp.NewHandler(m.service, m.logger)
	
	// Run database migrations if needed
	if err := m.runMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run catalog migrations: %w", err)
	}
	
	m.logger.Info("Catalog module initialized successfully")
	return nil
}

// Start initializes and starts the module's background processes
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting catalog module background processes")
	
	// Start background processes if any
	if err := m.startBackgroundProcesses(ctx); err != nil {
		return fmt.Errorf("failed to start background processes: %w", err)
	}
	
	m.logger.Info("Catalog module started successfully")
	
	// Block until context is canceled
	<-ctx.Done()
	m.logger.Info("Shutting down catalog module")
	
	return nil
}

// runMigrations runs database migrations for the catalog module
func (m *Module) runMigrations(ctx context.Context) error {
	m.logger.Info("Running catalog database migrations")
	
	// In a real implementation, this would use a migration library
	// like golang-migrate to run SQL migrations
	
	return nil
}

// startBackgroundProcesses starts any background processes for the catalog module
func (m *Module) startBackgroundProcesses(ctx context.Context) error {
	// Start any background processes like cache warming, data synchronization, etc.
	go func() {
		// Use context to handle cancellation
		<-ctx.Done()
		m.logger.Info("Background processes shutting down")
	}()
	
	return nil
}

// Health returns the health status of the module
func (m *Module) Health() map[string]interface{} {
	status := "ok"
	info := map[string]interface{}{
		"status": status,
	}
	
	// Add more health information if available
	if m.repository != nil {
		if healthCheck, ok := m.repository.(interface{ HealthCheck() error }); ok {
			if err := healthCheck.HealthCheck(); err != nil {
				status = "degraded"
				info["repository_error"] = err.Error()
			}
		}
	}
	
	return info
}

// Metrics returns the module's metrics
func (m *Module) Metrics() map[string]interface{} {
	// Return any module-specific metrics
	return map[string]interface{}{
		"equipment_count": 0, // This would be a real count in production
	}
}
