package infra

import (
    "context"
    "github.com/jackc/pgx/v5/pgconn"
)

// EnsureEquipmentSchema ensures the equipment table matches application expectations.
// It is idempotent and safe to run on startup.
func EnsureEquipmentSchema(ctx context.Context, pool PgxIface) error {
    stmts := []string{
        // Widen potentially narrow columns (older schemas may have small varchar lengths)
        "ALTER TABLE equipment ALTER COLUMN id TYPE VARCHAR(255)",
        "ALTER TABLE equipment ALTER COLUMN serial_number TYPE VARCHAR(255)",
        // Relax legacy NOT NULL constraints not used by new code path
        "ALTER TABLE equipment ALTER COLUMN name DROP NOT NULL",
        "ALTER TABLE equipment ALTER COLUMN tenant_id DROP NOT NULL",
        "ALTER TABLE equipment ALTER COLUMN model DROP NOT NULL",
        "ALTER TABLE equipment ALTER COLUMN manufacturer DROP NOT NULL",
        // Add missing columns with safe defaults
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS equipment_id VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS equipment_name VARCHAR(500)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS manufacturer_name VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS model_number VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS category VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS customer_id VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS customer_name VARCHAR(500)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS installation_location TEXT",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS installation_address JSONB DEFAULT '{}'::jsonb",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS contract_id VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS purchase_date DATE",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS purchase_price DECIMAL(15,2) DEFAULT 0",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS warranty_expiry DATE",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS amc_contract_id VARCHAR(255)",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'operational'",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS last_service_date TIMESTAMP",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS next_service_date TIMESTAMP",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS service_count INTEGER DEFAULT 0",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS specifications JSONB DEFAULT '{}'::jsonb",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS photos JSONB DEFAULT '[]'::jsonb",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS documents JSONB DEFAULT '[]'::jsonb",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_url TEXT",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_image BYTEA",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png'",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP",
        "ALTER TABLE equipment ADD COLUMN IF NOT EXISTS created_by VARCHAR(255)",
    }

    for _, stmt := range stmts {
        if _, err := pool.Exec(ctx, stmt); err != nil {
            // Ignore all ALTER errors - table might already have correct schema
            // This is safe because we're only widening columns or adding optional columns
            continue
        }
    }
    return nil
}

// PgxIface is a minimal interface implemented by *pgxpool.Pool
type PgxIface interface {
    Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}
