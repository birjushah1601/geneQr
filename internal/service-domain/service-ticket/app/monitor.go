package app

import (
    "context"
    "log/slog"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type SLAMonitor struct {
    pool   *pgxpool.Pool
    logger *slog.Logger
}

func NewSLAMonitor(pool *pgxpool.Pool, logger *slog.Logger) *SLAMonitor {
    return &SLAMonitor{pool: pool, logger: logger.With(slog.String("component", "sla_monitor"))}
}

func (m *SLAMonitor) Run(ctx context.Context) {
    if !enabled(os.Getenv("ENABLE_SLA_MONITOR")) { return }
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            _ = m.checkBreaches(ctx)
        }
    }
}

func (m *SLAMonitor) checkBreaches(ctx context.Context) error {
    // Response breach: not acknowledged and response_due passed
    const q1 = `UPDATE service_tickets
                SET sla_breached = true
                WHERE COALESCE(sla_breached,false) = false
                  AND sla_response_due IS NOT NULL
                  AND acknowledged_at IS NULL
                  AND NOW() > sla_response_due`
    // Resolution breach: not resolved and resolution_due passed
    const q2 = `UPDATE service_tickets
                SET sla_breached = true
                WHERE COALESCE(sla_breached,false) = false
                  AND sla_resolution_due IS NOT NULL
                  AND resolved_at IS NULL
                  AND NOW() > sla_resolution_due`
    if _, err := m.pool.Exec(ctx, q1); err != nil { return err }
    if _, err := m.pool.Exec(ctx, q2); err != nil { return err }
    return nil
}
