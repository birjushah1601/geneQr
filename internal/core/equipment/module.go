package equipment

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aby-med/medical-platform/internal/core/equipment/api"
	"github.com/aby-med/medical-platform/internal/shared/config"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ModuleName = "equipment"

type Module struct {
	cfg    *config.Config
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func New(cfg *config.Config, logger *slog.Logger) *Module {
	return &Module{
		cfg:    cfg,
		logger: logger.With(slog.String("module", ModuleName)),
	}
}

func (m *Module) Name() string { return ModuleName }

func (m *Module) Initialize(ctx context.Context) error {
	// Feature flag guard
	if !isEnabled(os.Getenv("ENABLE_EQUIPMENT")) {
		m.logger.Info("Equipment module disabled by flag ENABLE_EQUIPMENT=false")
		return nil
	}

	// DB pool
	pool, err := pgxpool.New(ctx, m.cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("pgxpool new: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}
	m.pool = pool

	m.logger.Info("Equipment module initialized")
	return nil
}

// GetDB returns the database pool for use by other services
func (m *Module) GetDB() *pgxpool.Pool {
	return m.pool
}

func (m *Module) MountRoutes(r chi.Router) {
	if !isEnabled(os.Getenv("ENABLE_EQUIPMENT")) {
		return
	}
	if m.pool == nil {
		m.logger.Warn("Equipment module not initialized; skipping routes")
		return
	}

	m.logger.Info("Mounting equipment routes")

	// Catalog bulk import handler
	catalogImportHandler := api.NewCatalogBulkImportHandler(m.pool, m.logger)

	r.Route("/equipment", func(r chi.Router) {
		r.Post("/catalog/import", catalogImportHandler.HandleCatalogBulkImport)
	})
}

func (m *Module) Start(ctx context.Context) error {
	// No background processes yet
	<-ctx.Done()
	return nil
}

func isEnabled(v string) bool {
	switch v {
	case "1", "true", "TRUE", "True", "yes", "on":
		return true
	default:
		return false
	}
}
