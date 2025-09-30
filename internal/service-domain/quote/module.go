package quote

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/quote/api"
	"github.com/aby-med/medical-platform/internal/service-domain/quote/app"
	"github.com/aby-med/medical-platform/internal/service-domain/quote/infra"
	"github.com/go-chi/chi/v5"
)

// Config holds configuration for the quote module
type Config struct {
	DatabaseURL string
}

// Module represents the quote service module
type Module struct {
	config  Config
	logger  *slog.Logger
	db      *infra.PostgresDB
	handler *api.QuoteHandler
}

// NewModule creates a new quote module instance
func NewModule(config Config, logger *slog.Logger) *Module {
	return &Module{
		config: config,
		logger: logger.With(slog.String("module", "quote")),
	}
}

// Initialize sets up the module dependencies
func (m *Module) Initialize(ctx context.Context) error {
	m.logger.Info("Initializing quote module")

	// Initialize database
	db, err := infra.NewPostgresDB(ctx, m.config.DatabaseURL, m.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	m.db = db

	// Initialize layers
	repo := infra.NewQuoteRepository(db, m.logger)
	service := app.NewQuoteService(repo, m.logger)
	m.handler = api.NewQuoteHandler(service, m.logger)

	m.logger.Info("Quote module initialized successfully")
	return nil
}

// MountRoutes registers HTTP routes for the quote module
func (m *Module) MountRoutes(r chi.Router) {
	m.logger.Info("Mounting quote routes at /quotes")

	r.Route("/quotes", func(r chi.Router) {
		// Core quote operations
		r.Post("/", m.handler.CreateQuote)
		r.Get("/", m.handler.ListQuotes)
		r.Get("/{id}", m.handler.GetQuote)
		r.Patch("/{id}", m.handler.UpdateQuote)
		r.Delete("/{id}", m.handler.DeleteQuote)

		// Quote item operations
		r.Post("/{id}/items", m.handler.AddQuoteItem)
		r.Patch("/{id}/items/{item_id}", m.handler.UpdateQuoteItem)
		r.Delete("/{id}/items/{item_id}", m.handler.RemoveQuoteItem)

		// Quote lifecycle operations
		r.Post("/{id}/submit", m.handler.SubmitQuote)
		r.Post("/{id}/revise", m.handler.ReviseQuote)
		r.Post("/{id}/accept", m.handler.AcceptQuote)
		r.Post("/{id}/reject", m.handler.RejectQuote)
		r.Post("/{id}/withdraw", m.handler.WithdrawQuote)
		r.Post("/{id}/under-review", m.handler.MarkUnderReview)
	})

	// Additional routes for querying quotes by RFQ or Supplier
	r.Get("/rfqs/{rfq_id}/quotes", m.handler.GetQuotesByRFQ)
	r.Get("/suppliers/{supplier_id}/quotes", m.handler.GetQuotesBySupplier)

	m.logger.Info("Quote routes mounted successfully")
}

// Start begins any background processes
func (m *Module) Start(ctx context.Context) error {
	m.logger.Info("Starting quote module")
	// No background processes for now
	return nil
}

// Stop gracefully shuts down the module
func (m *Module) Stop(ctx context.Context) error {
	m.logger.Info("Stopping quote module")

	if m.db != nil {
		m.db.Close()
	}

	return nil
}

// Name returns the module name
func (m *Module) Name() string {
	return "quote"
}
