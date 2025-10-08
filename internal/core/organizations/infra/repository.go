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
