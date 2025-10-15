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
        _, err := db.Exec(ctx, `INSERT INTO organizations (name, org_type, status, metadata) VALUES
            -- Manufacturers (India)
            ('Wipro GE Healthcare','manufacturer','active', '{"country":"IN","city":"Bengaluru"}'),
            ('Siemens Healthineers India','manufacturer','active', '{"country":"IN","city":"Gurugram"}'),
            ('Philips Healthcare India','manufacturer','active', '{"country":"IN","city":"Pune"}'),

            -- Suppliers / Distributors (India)
            ('SouthCare Distributors','distributor','active', '{"country":"IN","city":"Chennai"}'),
            ('MedSupply Mumbai','distributor','active', '{"country":"IN","city":"Mumbai"}'),
            ('Aknamed Medical Supplies','supplier','active', '{"country":"IN","city":"Bengaluru"}'),

            -- Hospitals (India)
            ('AIIMS New Delhi','hospital','active', '{"country":"IN","city":"New Delhi"}'),
            ('Apollo Hospitals Chennai','hospital','active', '{"country":"IN","city":"Chennai"}'),
            ('Fortis Hospital Mumbai','hospital','active', '{"country":"IN","city":"Mumbai"}'),
            ('Manipal Hospitals Bengaluru','hospital','active', '{"country":"IN","city":"Bengaluru"}'),
            ('Yashoda Hospitals Hyderabad','hospital','active', '{"country":"IN","city":"Hyderabad"}'),

            -- Imaging centers (India)
            ('Aarthi Scans & Labs - Chennai','imaging_center','active', '{"country":"IN","city":"Chennai"}'),
            ('Vijaya Diagnostic Centre - Hyderabad','imaging_center','active', '{"country":"IN","city":"Hyderabad"}'),
            ('SRL Diagnostics Imaging - Mumbai','imaging_center','active', '{"country":"IN","city":"Mumbai"}')`)
        if err != nil { return err }

        // Sample supply chain relationships (manufacturer -> distributor) and (distributor -> hospital)
        _, err = db.Exec(ctx, `
          WITH m1 AS (SELECT id FROM organizations WHERE name='Wipro GE Healthcare'),
               m2 AS (SELECT id FROM organizations WHERE name='Siemens Healthineers India'),
               d1 AS (SELECT id FROM organizations WHERE name='SouthCare Distributors'),
               d2 AS (SELECT id FROM organizations WHERE name='MedSupply Mumbai'),
               h1 AS (SELECT id FROM organizations WHERE name='Apollo Hospitals Chennai'),
               h2 AS (SELECT id FROM organizations WHERE name='Fortis Hospital Mumbai')
          INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type)
          SELECT m1.id, d1.id, 'distributor_of' FROM m1, d1
          UNION ALL
          SELECT m2.id, d2.id, 'distributor_of' FROM m2, d2
          UNION ALL
          SELECT d1.id, h1.id, 'supplier_of' FROM d1, h1
          UNION ALL
          SELECT d2.id, h2.id, 'supplier_of' FROM d2, h2;
        `)
        if err != nil { return err }
    }

    // Ensure key India entities exist (idempotent; insert if missing)
    ensureStmt := func(name, otype, city string) string {
        return "INSERT INTO organizations(name, org_type, status, metadata) " +
            "SELECT '" + name + "','" + otype + "','active', '{\"country\":\"IN\",\"city\":\"" + city + "\"}' " +
            "WHERE NOT EXISTS (SELECT 1 FROM organizations WHERE name='" + name + "');"
    }
    stmts := []string{
        ensureStmt("Wipro GE Healthcare", "manufacturer", "Bengaluru"),
        ensureStmt("Siemens Healthineers India", "manufacturer", "Gurugram"),
        ensureStmt("Philips Healthcare India", "manufacturer", "Pune"),
        ensureStmt("SouthCare Distributors", "distributor", "Chennai"),
        ensureStmt("MedSupply Mumbai", "distributor", "Mumbai"),
        ensureStmt("Aknamed Medical Supplies", "supplier", "Bengaluru"),
        ensureStmt("AIIMS New Delhi", "hospital", "New Delhi"),
        ensureStmt("Apollo Hospitals Chennai", "hospital", "Chennai"),
        ensureStmt("Fortis Hospital Mumbai", "hospital", "Mumbai"),
        ensureStmt("Manipal Hospitals Bengaluru", "hospital", "Bengaluru"),
        ensureStmt("Yashoda Hospitals Hyderabad", "hospital", "Hyderabad"),
        ensureStmt("Aarthi Scans & Labs - Chennai", "imaging_center", "Chennai"),
        ensureStmt("Vijaya Diagnostic Centre - Hyderabad", "imaging_center", "Hyderabad"),
        ensureStmt("SRL Diagnostics Imaging - Mumbai", "imaging_center", "Mumbai"),
    }
    for _, s := range stmts {
        if _, err := db.Exec(ctx, s); err != nil { return err }
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
