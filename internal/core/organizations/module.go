package organizations

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "github.com/aby-med/medical-platform/internal/core/organizations/api"
    "github.com/aby-med/medical-platform/internal/core/organizations/infra"
    "github.com/aby-med/medical-platform/internal/shared/config"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

const ModuleName = "organizations"

type Module struct {
    cfg     *config.Config
    logger  *slog.Logger
    pool    *pgxpool.Pool
    handler *api.Handler
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
    if !isEnabled(os.Getenv("ENABLE_ORG")) {
        m.logger.Info("Organizations module disabled by flag ENABLE_ORG=false")
        return nil
    }

    // DB pool
    pool, err := pgxpool.New(ctx, m.cfg.GetDSN())
    if err != nil { return fmt.Errorf("pgxpool new: %w", err) }
    if err := pool.Ping(ctx); err != nil { return fmt.Errorf("db ping: %w", err) }
    m.pool = pool

    // Ensure schema (create tables if not exists)
    if err := infra.EnsureOrgSchema(ctx, pool, m.logger); err != nil {
        return fmt.Errorf("ensure org schema: %w", err)
    }

    repo := infra.NewRepository(pool, m.logger)
    m.handler = api.NewHandler(repo, m.logger)
    m.logger.Info("Organizations module initialized")
    return nil
}

func (m *Module) MountRoutes(r chi.Router) {
    if !isEnabled(os.Getenv("ENABLE_ORG")) {
        return
    }
    if m.handler == nil {
        m.logger.Warn("Organizations handler not initialized; skipping routes")
        return
    }
    m.logger.Info("Mounting organizations routes")
    r.Route("/orgs", func(r chi.Router) {
        r.Get("/", m.handler.ListOrgs)
        r.Get("/{id}/relationships", m.handler.ListRelationships)
    })
}

func (m *Module) Start(ctx context.Context) error {
    // no background processes yet
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
