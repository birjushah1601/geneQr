package contract

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/contract/api"
	"github.com/aby-med/medical-platform/internal/service-domain/contract/app"
	"github.com/aby-med/medical-platform/internal/service-domain/contract/infra"
	"github.com/go-chi/chi/v5"
)

// Config holds configuration for the contract module
type Config struct {
	DatabaseURL string
}

// Module represents the contract service module
type Module struct {
	config  Config
	logger  *slog.Logger
	db      *infra.PostgresDB
	handler *api.ContractHandler
}

// NewModule creates a new contract module instance
func NewModule(config Config, logger *slog.Logger) *Module {
	return &Module{
		config: config,
		logger: logger.With(slog.String("module", "contract")),
	}
}

// Initialize sets up the module dependencies
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing contract module")

	// Initialize database
	db, err := infra.NewPostgresDB(ctx, m.config.DatabaseURL, m.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Initialize layers
	repo := infra.NewContractRepository(db, m.logger)
	service := app.NewContractService(repo, m.logger)
	m.handler = api.NewContractHandler(service, m.logger)

	m.logger.Info("Contract module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes for the contract module
func (m *Module) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting contract routes at /contracts")

	r.Route("/contracts", func(r chi.Router) {
		// Core operations
		r.Post("/", m.handler.CreateContract)
		r.Get("/", m.handler.ListContracts)
		r.Get("/{id}", m.handler.GetContract)
		r.Patch("/{id}", m.handler.UpdateContract)
		r.Delete("/{id}", m.handler.DeleteContract)

		// Lifecycle operations
		r.Post("/{id}/activate", m.handler.ActivateContract)
		r.Post("/{id}/sign", m.handler.SignContract)
		r.Post("/{id}/complete", m.handler.CompleteContract)
		r.Post("/{id}/cancel", m.handler.CancelContract)
		r.Post("/{id}/suspend", m.handler.SuspendContract)
		r.Post("/{id}/resume", m.handler.ResumeContract)

		// Amendment
		r.Post("/{id}/amendments", m.handler.AddAmendment)

		// Payment and delivery tracking
		r.Post("/{id}/payments/paid", m.handler.MarkPaymentPaid)
		r.Post("/{id}/deliveries/completed", m.handler.MarkDeliveryCompleted)
	})

	// Query routes
	r.Get("/contracts/by-number/{number}", m.handler.GetContractByNumber)
	r.Get("/rfqs/{rfq_id}/contracts", m.handler.GetContractsByRFQ)
	r.Get("/suppliers/{supplier_id}/contracts", m.handler.GetContractsBySupplier)

	m.logger.Info("Contract routes mounted successfully")
}

// Start begins any background processes
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting contract module")
	// No background processes for now
	return nil
}

// Stop gracefully shuts down the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Stopping contract module")

	if m.db != nil {
		m.db.Close()
	}

	return nil
}

// Name returns the module name
func (m *Module) Name() string {
	return "contract"
}
