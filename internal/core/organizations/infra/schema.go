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

        // Channels (read-only Phase 1)
        `CREATE TABLE IF NOT EXISTS channels (
            id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            code        TEXT NOT NULL UNIQUE,
            name        TEXT NOT NULL,
            channel_type TEXT NULL, -- online|offline|partner|direct|marketplace
            metadata    JSONB NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_channels_type ON channels(channel_type);`,

        // Products (manufacturer optional)
        `CREATE TABLE IF NOT EXISTS products (
            id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name                TEXT NOT NULL,
            manufacturer_org_id UUID NULL REFERENCES organizations(id) ON DELETE SET NULL,
            external_ref        TEXT NULL,
            metadata            JSONB NULL,
            created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_products_mfr ON products(manufacturer_org_id);`,

        // SKUs (per product)
        `CREATE TABLE IF NOT EXISTS skus (
            id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
            sku_code    TEXT NOT NULL UNIQUE,
            status      TEXT NOT NULL DEFAULT 'active',
            attributes  JSONB NULL,
            metadata    JSONB NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_skus_product ON skus(product_id);`,

        // Phase 2: offerings + channel_catalog (publish flow)
        `CREATE TABLE IF NOT EXISTS offerings (
            id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            sku_id       UUID NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
            owner_org_id UUID NULL REFERENCES organizations(id) ON DELETE SET NULL,
            status       TEXT NOT NULL DEFAULT 'draft', -- draft|published
            version      INT  NOT NULL DEFAULT 1,
            data         JSONB NULL,
            created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_offerings_sku ON offerings(sku_id);`,
        `CREATE INDEX IF NOT EXISTS idx_offerings_owner ON offerings(owner_org_id);`,

        `CREATE TABLE IF NOT EXISTS channel_catalog (
            id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            channel_id    UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
            offering_id   UUID NOT NULL REFERENCES offerings(id) ON DELETE CASCADE,
            listed        BOOLEAN NOT NULL DEFAULT true,
            published_version INT NOT NULL DEFAULT 0,
            created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
            UNIQUE(channel_id, offering_id)
        );`,
        `CREATE INDEX IF NOT EXISTS idx_channel_catalog_channel ON channel_catalog(channel_id);`,

        // Phase 3: price books + rules
        `CREATE TABLE IF NOT EXISTS price_books (
            id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name         TEXT NOT NULL,
            org_id       UUID NULL REFERENCES organizations(id) ON DELETE SET NULL,
            channel_id   UUID NULL REFERENCES channels(id) ON DELETE SET NULL,
            currency     TEXT NOT NULL DEFAULT 'INR',
            created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE INDEX IF NOT EXISTS idx_price_books_scope ON price_books(org_id, channel_id);`,

        `CREATE TABLE IF NOT EXISTS price_rules (
            id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            book_id      UUID NOT NULL REFERENCES price_books(id) ON DELETE CASCADE,
            sku_id       UUID NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
            price        NUMERIC(18,2) NOT NULL,
            valid_from   TIMESTAMPTZ NULL,
            valid_to     TIMESTAMPTZ NULL,
            created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
            UNIQUE(book_id, sku_id)
        );`,
        `CREATE INDEX IF NOT EXISTS idx_price_rules_sku ON price_rules(sku_id);`,

        // Phase 5: engineers + memberships + coverage (optional)
        `CREATE TABLE IF NOT EXISTS engineers (
            id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name        TEXT NOT NULL,
            phone       TEXT NULL,
            email       TEXT NULL,
            skills      TEXT[] NULL,
            home_region TEXT NULL,
            metadata    JSONB NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
        `CREATE TABLE IF NOT EXISTS engineer_org_memberships (
            engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
            org_id      UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
            role        TEXT NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
            PRIMARY KEY(engineer_id, org_id)
        );`,
        `CREATE TABLE IF NOT EXISTS engineer_coverage (
            id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
            region      TEXT NULL,
            skills      TEXT[] NULL,
            metadata    JSONB NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,

        // Phase 6: agreements/contracts (backend-only skeleton)
        `CREATE TABLE IF NOT EXISTS agreements (
            id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            org_id      UUID NULL REFERENCES organizations(id) ON DELETE SET NULL,
            agreement_type TEXT NOT NULL, -- warranty|amc|subscription|other
            start_date  DATE NULL,
            end_date    DATE NULL,
            metadata    JSONB NULL,
            created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
        );`,
    }

    for _, s := range stmts {
        if _, err := db.Exec(ctx, s); err != nil {
            logger.Error("schema exec failed", slog.String("stmt", s), slog.String("err", err.Error()))
            return err
        }
    }
    return nil
}
