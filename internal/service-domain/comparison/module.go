package comparison

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/comparison/api"
	"github.com/aby-med/medical-platform/internal/service-domain/comparison/app"
	"github.com/aby-med/medical-platform/internal/service-domain/comparison/infra"
	"github.com/go-chi/chi/v5"
)

// Config holds configuration for the comparison module
type Config struct {
	DatabaseURL string
}

// Module represents the comparison service module
type Module struct {
	config  Config
	logger  *slog.Logger
	db      *infra.PostgresDB
	handler *api.ComparisonHandler
}

// NewModule creates a new comparison module instance
func NewModule(config Config, logger *slog.Logger) *Module {
	return &Module{
		config: config,
		logger: logger.With(slog.String("module", "comparison")),
	}
}

// Initialize sets up the module dependencies
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing comparison module")

	// Initialize database
	db, err := infra.NewPostgresDB(ctx, m.config.DatabaseURL, m.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Initialize layers
	repo := infra.NewComparisonRepository(db, m.logger)
	service := app.NewComparisonService(repo, m.logger)
	m.handler = api.NewComparisonHandler(service, m.logger)

	m.logger.Info("Comparison module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes for the comparison module
func (m *Module) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting comparison routes at /comparisons")

	r.Route("/comparisons", func(r chi.Router) {
		// Core comparison operations
		r.Post("/", m.handler.CreateComparison)
		r.Get("/", m.handler.ListComparisons)
		r.Get("/{id}", m.handler.GetComparison)
		r.Patch("/{id}", m.handler.UpdateComparison)
		r.Delete("/{id}", m.handler.DeleteComparison)

		// Comparison management
		r.Post("/{id}/quotes", m.handler.AddQuote)
		r.Delete("/{id}/quotes/{quote_id}", m.handler.RemoveQuote)
		r.Post("/{id}/scoring", m.handler.UpdateScoringCriteria)
		r.Post("/{id}/calculate", m.handler.CalculateScores)

		// Lifecycle operations
		r.Post("/{id}/activate", m.handler.ActivateComparison)
		r.Post("/{id}/complete", m.handler.CompleteComparison)
		r.Post("/{id}/archive", m.handler.ArchiveComparison)
	})

	// Additional route for getting comparisons by RFQ
	r.Get("/rfqs/{rfq_id}/comparisons", m.handler.GetComparisonsByRFQ)

	m.logger.Info("Comparison routes mounted successfully")
}

// Start begins any background processes
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting comparison module")
	// No background processes for now
	return nil
}

// Stop gracefully shuts down the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Stopping comparison module")

	if m.db != nil {
		m.db.Close()
	}

	return nil
}

// Name returns the module name
func (m *Module) Name() string {
	return "comparison"
}
