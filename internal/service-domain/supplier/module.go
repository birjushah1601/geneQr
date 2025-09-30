package supplier

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/supplier/api"
	"github.com/aby-med/medical-platform/internal/service-domain/supplier/app"
	"github.com/aby-med/medical-platform/internal/service-domain/supplier/infra"
	"github.com/go-chi/chi/v5"
)

// Config holds configuration for the Supplier module
type Config struct {
	DatabaseDSN string
}

// Module implements the service.Module interface for the supplier service
type Module struct {
	config  *Config
	db      *infra.PostgresDB
	handler *api.SupplierHandler
	logger  *slog.Logger
}

// NewModule creates a new supplier module
func NewModule(config *Config, logger *slog.Logger) *Module {
	return &Module{
		config: config,
		logger: logger.With(slog.String("module", "supplier")),
	}
}

// Name returns the module name
func (m *Module) Name() string {
	return "supplier"
}

// Initialize initializes the supplier module
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing supplier module")

	// Create database connection
	db, err := infra.NewPostgresDB(ctx, m.config.DatabaseDSN, m.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Create repository
	repo := infra.NewSupplierRepository(db, m.logger)

	// Create application service
	supplierService := app.NewSupplierService(repo, m.logger)

	// Create HTTP handler
	m.handler = api.NewSupplierHandler(supplierService, m.logger)

	m.logger.Info("Supplier module initialized successfully")
	return nil
}

// MountRoutes mounts the module's HTTP routes
func (m *Module) MountRoutes(router chi.Router) {
	m.logger.Info("Mounting supplier routes at /api/v1/suppliers")
	router.Route("/suppliers", func(r chi.Router) {
		m.handler.RegisterRoutes(r)
	})
}

// Start starts the supplier module background processes
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting supplier module")
	// No background processes to start for now
	return nil
}

// Stop stops the supplier module (graceful shutdown)
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Stopping supplier module")
	
	if m.db != nil {
		m.db.Close()
	}
	
	return nil
}
