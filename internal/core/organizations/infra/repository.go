package infra

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
    ID                string          `json:"id"`
    Name              string          `json:"name"`
    OrgType           string          `json:"org_type"`
    Status            string          `json:"status"`
    Metadata          json.RawMessage `json:"metadata"`
    EquipmentCount    int             `json:"equipment_count,omitempty"`
    EngineersCount    int             `json:"engineers_count,omitempty"`
    ActiveTickets     int             `json:"active_tickets,omitempty"`
}

type Relationship struct {
    ID          string `json:"id"`
    ParentOrgID string `json:"parent_org_id"`
    ChildOrgID  string `json:"child_org_id"`
    RelType     string `json:"rel_type"`
}

type Repository struct {
    db     *pgxpool.Pool
    logger *slog.Logger
}

func NewRepository(db *pgxpool.Pool, logger *slog.Logger) *Repository {
    return &Repository{db: db, logger: logger}
}

func (r *Repository) DB() *pgxpool.Pool {
    return r.db
}

func (r *Repository) ListOrgs(ctx context.Context, limit, offset int, orgType, status string) ([]Organization, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    
    query := `SELECT id, name, org_type, status, COALESCE(metadata, '{}'::jsonb) FROM organizations WHERE 1=1`
    args := []interface{}{}
    argPos := 1
    
    // Add filters if provided
    if orgType != "" {
        query += ` AND org_type = $` + fmt.Sprintf("%d", argPos)
        args = append(args, orgType)
        argPos++
    }
    
    if status != "" {
        query += ` AND status = $` + fmt.Sprintf("%d", argPos)
        args = append(args, status)
        argPos++
    }
    
    query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)
    args = append(args, limit, offset)
    
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Organization
    for rows.Next() {
        var o Organization
        if err := rows.Scan(&o.ID, &o.Name, &o.OrgType, &o.Status, &o.Metadata); err != nil { return nil, err }
        out = append(out, o)
    }
    return out, rows.Err()
}

func (r *Repository) GetOrgByID(ctx context.Context, id string) (*Organization, error) {
    row := r.db.QueryRow(ctx, `SELECT id, name, org_type, status, COALESCE(metadata, '{}'::jsonb) FROM organizations WHERE id=$1`, id)
    var o Organization
    if err := row.Scan(&o.ID, &o.Name, &o.OrgType, &o.Status, &o.Metadata); err != nil {
        return nil, err
    }
    return &o, nil
}

func (r *Repository) GetEquipmentCount(ctx context.Context, manufacturerID string) (int, error) {
    var count int
    err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM equipment_registry WHERE manufacturer_id = $1`, manufacturerID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (r *Repository) GetEngineersCount(ctx context.Context, organizationID string) (int, error) {
    var count int
    err := r.db.QueryRow(ctx, `SELECT COUNT(DISTINCT engineer_id) FROM engineer_org_memberships WHERE org_id = $1`, organizationID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (r *Repository) GetActiveTicketsCount(ctx context.Context, manufacturerID string) (int, error) {
    var count int
    query := `
        SELECT COUNT(DISTINCT st.id) 
        FROM service_tickets st
        JOIN equipment_registry er ON st.equipment_id = er.id
        WHERE er.manufacturer_id = $1 
        AND st.status != 'closed'
    `
    err := r.db.QueryRow(ctx, query, manufacturerID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

type Facility struct {
    ID             string `json:"id"`
    OrgID          string `json:"org_id"`
    FacilityName   string `json:"facility_name"`
    FacilityCode   string `json:"facility_code"`
    FacilityType   string `json:"facility_type"`
    Address        []byte `json:"address"`
    Status         string `json:"status"`
}

func (r *Repository) ListFacilities(ctx context.Context, orgID string) ([]Facility, error) {
    rows, err := r.db.Query(ctx, `SELECT id, org_id, facility_name, facility_code, facility_type, COALESCE(address, '{}'::jsonb), status FROM organization_facilities WHERE org_id=$1 ORDER BY created_at DESC`, orgID)
    if err != nil { return nil, err }
    defer rows.Close()
    // ensure empty slice instead of null when no rows
    out := make([]Facility, 0)
    for rows.Next() {
        var f Facility
        if err := rows.Scan(&f.ID, &f.OrgID, &f.FacilityName, &f.FacilityCode, &f.FacilityType, &f.Address, &f.Status); err != nil { return nil, err }
        out = append(out, f)
    }
    return out, rows.Err()
}

func (r *Repository) ListRelationships(ctx context.Context, orgID string) ([]Relationship, error) {
    const q = `SELECT id, parent_org_id, child_org_id, rel_type FROM org_relationships WHERE parent_org_id=$1 OR child_org_id=$1 ORDER BY created_at DESC`
    rows, err := r.db.Query(ctx, q, orgID)
    if err != nil { return nil, err }
    defer rows.Close()
    // ensure empty slice instead of null when no rows
    out := make([]Relationship, 0)
    for rows.Next() {
        var rel Relationship
        if err := rows.Scan(&rel.ID, &rel.ParentOrgID, &rel.ChildOrgID, &rel.RelType); err != nil { return nil, err }
        out = append(out, rel)
    }
    return out, rows.Err()
}

func (r *Repository) HealthCheck(ctx context.Context) error {
    return r.db.Ping(ctx)
}

func scanOne[T any](rows pgx.Row, dest *T) error { return nil }

// Channels
type Channel struct {
    ID     string `json:"id"`
    Code   string `json:"code"`
    Name   string `json:"name"`
    Type   string `json:"channel_type"`
}

func (r *Repository) ListChannels(ctx context.Context, limit, offset int) ([]Channel, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, code, name, COALESCE(channel_type,'') FROM channels ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Channel
    for rows.Next() {
        var c Channel
        if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Type); err != nil { return nil, err }
        out = append(out, c)
    }
    return out, rows.Err()
}

// Products
type Product struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
}

func (r *Repository) ListProducts(ctx context.Context, limit, offset int) ([]Product, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, name FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name); err != nil { return nil, err }
        out = append(out, p)
    }
    return out, rows.Err()
}

// SKUs
type SKU struct {
    ID      string `json:"id"`
    Product string `json:"product_id"`
    Code    string `json:"sku_code"`
    Status  string `json:"status"`
}

func (r *Repository) ListSkus(ctx context.Context, limit, offset int) ([]SKU, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, product_id, sku_code, status FROM skus ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []SKU
    for rows.Next() {
        var s SKU
        if err := rows.Scan(&s.ID, &s.Product, &s.Code, &s.Status); err != nil { return nil, err }
        out = append(out, s)
    }
    return out, rows.Err()
}

// Offerings and Channel Catalog (Phase 2)
type Offering struct {
    ID      string `json:"id"`
    SkuID   string `json:"sku_id"`
    Status  string `json:"status"`
    Version int    `json:"version"`
}

func (r *Repository) CreateOffering(ctx context.Context, skuID string, ownerOrgID *string, data []byte) (Offering, error) {
    var o Offering
    err := r.db.QueryRow(ctx,
        `INSERT INTO offerings (sku_id, owner_org_id, data) VALUES ($1,$2,$3)
         RETURNING id, sku_id, status, version`, skuID, ownerOrgID, data).
        Scan(&o.ID, &o.SkuID, &o.Status, &o.Version)
    return o, err
}

func (r *Repository) ListOfferings(ctx context.Context, limit, offset int) ([]Offering, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, sku_id, status, version FROM offerings ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Offering
    for rows.Next() {
        var o Offering
        if err := rows.Scan(&o.ID, &o.SkuID, &o.Status, &o.Version); err != nil { return nil, err }
        out = append(out, o)
    }
    return out, rows.Err()
}

func (r *Repository) PublishToChannel(ctx context.Context, channelID, offeringID string) error {
    _, err := r.db.Exec(ctx, `INSERT INTO channel_catalog(channel_id, offering_id, listed, published_version)
        VALUES ($1,$2,true,(SELECT version FROM offerings WHERE id=$2))
        ON CONFLICT (channel_id, offering_id)
        DO UPDATE SET listed=true, published_version=EXCLUDED.published_version, updated_at=now()`, channelID, offeringID)
    return err
}

func (r *Repository) UnlistFromChannel(ctx context.Context, channelID, offeringID string) error {
    _, err := r.db.Exec(ctx, `UPDATE channel_catalog SET listed=false, updated_at=now() WHERE channel_id=$1 AND offering_id=$2`, channelID, offeringID)
    return err
}

// Pricing (Phase 3)
type PriceBook struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    OrgID     *string `json:"org_id"`
    ChannelID *string `json:"channel_id"`
    Currency  string `json:"currency"`
}

func (r *Repository) CreatePriceBook(ctx context.Context, name string, orgID, channelID *string, currency string) (PriceBook, error) {
    var b PriceBook
    err := r.db.QueryRow(ctx, `INSERT INTO price_books(name, org_id, channel_id, currency) VALUES($1,$2,$3,$4)
        RETURNING id, name, org_id, channel_id, currency`, name, orgID, channelID, currency).
        Scan(&b.ID, &b.Name, &b.OrgID, &b.ChannelID, &b.Currency)
    return b, err
}

func (r *Repository) AddPriceRule(ctx context.Context, bookID, skuID string, price float64) error {
    _, err := r.db.Exec(ctx, `INSERT INTO price_rules(book_id, sku_id, price) VALUES($1,$2,$3)
        ON CONFLICT (book_id, sku_id) DO UPDATE SET price=EXCLUDED.price, updated_at=now()`, bookID, skuID, price)
    return err
}

type PriceResult struct { Price float64 `json:"price"`; Currency string `json:"currency"` }

// Resolve price with precedence: org+channel > channel > org > global (null,null)
func (r *Repository) ResolvePrice(ctx context.Context, skuID string, orgID, channelID *string) (*PriceResult, error) {
    const q = `
WITH candidates AS (
  SELECT pb.currency, pr.price,
         CASE
           WHEN pb.org_id IS NOT NULL AND pb.channel_id IS NOT NULL THEN 4
           WHEN pb.channel_id IS NOT NULL AND pb.org_id IS NULL THEN 3
           WHEN pb.org_id IS NOT NULL AND pb.channel_id IS NULL THEN 2
           ELSE 1
         END AS precedence
  FROM price_rules pr
  JOIN price_books pb ON pb.id = pr.book_id
  WHERE pr.sku_id = $1
    AND (pb.org_id IS NULL OR pb.org_id = $2)
    AND (pb.channel_id IS NULL OR pb.channel_id = $3)
)
SELECT price, currency FROM candidates
ORDER BY precedence DESC
LIMIT 1`
    var res PriceResult
    var o, c *string = orgID, channelID
    err := r.db.QueryRow(ctx, q, skuID, o, c).Scan(&res.Price, &res.Currency)
    if err != nil { return nil, err }
    return &res, nil
}

// Engineers (Phase 5)
type Engineer struct {
    ID    string   `json:"id"`
    Name  string   `json:"name"`
    Skills []string `json:"skills"`
}

func (r *Repository) ListEngineers(ctx context.Context, limit, offset int) ([]Engineer, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, full_name, COALESCE(skills, ARRAY[]::text[]) FROM engineers ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Engineer
    for rows.Next() {
        var e Engineer
        if err := rows.Scan(&e.ID, &e.Name, &e.Skills); err != nil { return nil, err }
        out = append(out, e)
    }
    return out, rows.Err()
}

func (r *Repository) EligibleEngineers(ctx context.Context, skills []string, region string, limit int) ([]Engineer, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    // Build query with optional filters
    // Skills: e.skills @> $1 OR EXISTS coverage with skills @> $1
    // Region: e.home_region=$2 OR EXISTS coverage with region=$2
    const base = `
SELECT DISTINCT e.id, e.full_name, COALESCE(e.skills, ARRAY[]::text[])
FROM engineers e
LEFT JOIN engineer_coverage c ON c.engineer_id = e.id
WHERE ( $1::text[] IS NULL OR (COALESCE(e.skills, ARRAY[]::text[]) @> $1::text[] OR COALESCE(c.skills, ARRAY[]::text[]) @> $1::text[]) )
  AND ( $2::text IS NULL OR e.home_region = $2 OR c.region = $2 )
ORDER BY e.created_at DESC
LIMIT $3`
    var skillsArray []string
    if len(skills) > 0 { skillsArray = skills } else { skillsArray = nil }
    var regionPtr *string
    if region != "" { regionPtr = &region }
    rows, err := r.db.Query(ctx, base, skillsArray, regionPtr, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Engineer
    for rows.Next() {
        var e Engineer
        if err := rows.Scan(&e.ID, &e.Name, &e.Skills); err != nil { return nil, err }
        out = append(out, e)
    }
    return out, rows.Err()
}
