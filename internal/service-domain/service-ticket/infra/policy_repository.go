package infra

import (
    "context"
    "encoding/json"

    domain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PolicyRepository struct {
    pool *pgxpool.Pool
}

func NewPolicyRepository(pool *pgxpool.Pool) *PolicyRepository { return &PolicyRepository{pool: pool} }

func (r *PolicyRepository) GetDefaultResponsibleOrg(ctx context.Context) (*string, error) {
    const q = `SELECT rules->>'default_org_id' AS default_org_id
               FROM service_policies
               WHERE enabled = true
               ORDER BY created_at DESC
               LIMIT 1`
    var id *string
    if err := r.pool.QueryRow(ctx, q).Scan(&id); err != nil {
        // If no rows, return nil without error
        // pgx returns ErrNoRows; treat as no policy configured
        return nil, nil
    }
    if id != nil && *id == "" { return nil, nil }
    return id, nil
}

func (r *PolicyRepository) GetSLARules(ctx context.Context, orgID *string) (*domain.SLARules, error) {
    const q = `SELECT rules FROM sla_policies
               WHERE active = true AND ((org_id IS NULL) OR (org_id = COALESCE($1::uuid, org_id)))
               ORDER BY CASE WHEN org_id IS NULL THEN 1 ELSE 0 END, updated_at DESC
               LIMIT 1`
    var raw []byte
    if err := r.pool.QueryRow(ctx, q, orgID).Scan(&raw); err != nil {
        return nil, nil
    }
    var rules domain.SLARules
    if err := json.Unmarshal(raw, &rules); err == nil {
        return &rules, nil
    }
    var env struct{ Priority domain.SLARules `json:"priority"` }
    if err := json.Unmarshal(raw, &env); err == nil {
        return &env.Priority, nil
    }
    return nil, nil
}

var _ domain.PolicyRepository = (*PolicyRepository)(nil)
