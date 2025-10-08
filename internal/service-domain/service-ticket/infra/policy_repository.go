package infra

import (
    "context"

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

var _ domain.PolicyRepository = (*PolicyRepository)(nil)
