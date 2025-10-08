package infra

import (
    "context"
    "log/slog"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

// SeedOrgDemoData inserts minimal demo records if tables are empty.
func SeedOrgDemoData(ctx context.Context, db PgxIface, logger *slog.Logger) error {
    // organizations
    var cnt int
    if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM organizations`).Scan(&cnt); err != nil { return err }
    if cnt == 0 {
        _, err := db.Exec(ctx, `INSERT INTO organizations (name, org_type) VALUES
            ('Global Manufacturer A','manufacturer'),
            ('Regional Distributor X','distributor'),
            ('Local Dealer Z','dealer'),
            ('Supplier S','supplier')`)
        if err != nil { return err }
    }

    // channels
    if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM channels`).Scan(&cnt); err != nil { return err }
    if cnt == 0 {
        _, err := db.Exec(ctx, `INSERT INTO channels (code, name, channel_type) VALUES
            ('DIRECT','Direct','direct'),
            ('PARTNER','Partner','partner'),
            ('ONLINE','Online','online')`)
        if err != nil { return err }
    }

    // products
    if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM products`).Scan(&cnt); err != nil { return err }
    if cnt == 0 {
        _, err := db.Exec(ctx, `INSERT INTO products (name) VALUES
            ('ECG Machine Pro'),('Ventilator Max'),('Infusion Pump Lite')`)
        if err != nil { return err }
    }

    // skus
    if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM skus`).Scan(&cnt); err != nil { return err }
    if cnt == 0 {
        _, err := db.Exec(ctx, `INSERT INTO skus (product_id, sku_code, status)
            SELECT p.id, 'SKU-' || LEFT(md5(p.name),8), 'active' FROM products p LIMIT 3`)
        if err != nil { return err }
    }
    logger.Info("Seeded org/channel/product/sku demo data if empty")
    return nil
}

// PgxIface captures the minimal pgx methods used (for easier testing)
type PgxIface interface {
    Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
    QueryRow(context.Context, string, ...interface{}) pgx.Row
}
