package infra

import (
    "context"
    "log/slog"

    "github.com/jackc/pgx/v5/pgxpool"
)

// EnsureOrgSchema creates core tables with nullable/optional fields.
func EnsureOrgSchema(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) error {
    stmts := []string{
        `CREATE TABLE IF NOT EXISTS organizations (
            id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name             TEXT NOT NULL,
            org_type         TEXT NOT NULL, -- manufacturer|supplier|distributor|dealer|hospital|service_provider|other
            status           TEXT NOT NULL DEFAULT 'active',
            external_ref     TEXT NULL,
            metadata         JSONB NULL,
            created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_organizations_org_type ON organizations(org_type);`,
        `CREATE INDEX IF NOT EXISTS idx_organizations_status ON organizations(status);`,

        `CREATE TABLE IF NOT EXISTS org_relationships (
            id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            parent_org_id    UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
            child_org_id     UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
            rel_type         TEXT NOT NULL, -- manufacturer_of|distributor_of|dealer_of|supplier_of|partner_of
            metadata         JSONB NULL,
            created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_org_rel_parent ON org_relationships(parent_org_id);`,
        `CREATE INDEX IF NOT EXISTS idx_org_rel_child ON org_relationships(child_org_id);`,
        `CREATE INDEX IF NOT EXISTS idx_org_rel_type ON org_relationships(rel_type);`,
    }

    for _, s := range stmts {
        if _, err := db.Exec(ctx, s); err != nil {
            logger.Error("schema exec failed", slog.String("stmt", s), slog.String("err", err.Error()))
            return err
        }
    }
    return nil
}
