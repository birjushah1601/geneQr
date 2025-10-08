package infra

import (
    "context"
    "log/slog"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Organization struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    OrgType  string `json:"org_type"`
    Status   string `json:"status"`
    Metadata []byte `json:"metadata"`
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

func (r *Repository) ListOrgs(ctx context.Context, limit, offset int) ([]Organization, error) {
    if limit <= 0 || limit > 500 { limit = 100 }
    if offset < 0 { offset = 0 }
    rows, err := r.db.Query(ctx, `SELECT id, name, org_type, status, COALESCE(metadata, '{}'::jsonb) FROM organizations ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
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

func (r *Repository) ListRelationships(ctx context.Context, orgID string) ([]Relationship, error) {
    const q = `SELECT id, parent_org_id, child_org_id, rel_type FROM org_relationships WHERE parent_org_id=$1 OR child_org_id=$1 ORDER BY created_at DESC`
    rows, err := r.db.Query(ctx, q, orgID)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Relationship
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
