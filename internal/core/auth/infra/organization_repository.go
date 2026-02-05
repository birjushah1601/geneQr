package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aby-med/medical-platform/internal/core/auth/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type organizationRepository struct {
	db *sqlx.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *sqlx.DB) domain.OrganizationRepository {
	return &organizationRepository{
		db: db,
	}
}

// GetByID retrieves an organization by its ID
func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization
	query := `
		SELECT id, name, org_type, status, external_ref, metadata, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &org, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	
	return &org, nil
}

// GetByUserID retrieves all organizations for a specific user
func (r *organizationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Organization, error) {
	var orgs []*domain.Organization
	query := `
		SELECT o.id, o.name, o.org_type, o.status, o.external_ref, o.metadata, o.created_at, o.updated_at
		FROM organizations o
		JOIN user_organizations uo ON o.id = uo.organization_id
		WHERE uo.user_id = $1
		AND uo.status = 'active'
		ORDER BY uo.is_primary DESC, o.name
	`
	
	err := r.db.SelectContext(ctx, &orgs, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}
	
	return orgs, nil
}
