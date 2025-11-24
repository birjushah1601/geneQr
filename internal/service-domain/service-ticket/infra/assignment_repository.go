package infra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
)

// AssignmentRepository implements domain.EngineerSuggestionRepository
type AssignmentRepository struct {
	pool *pgxpool.Pool
}

// NewAssignmentRepository creates a new assignment repository
func NewAssignmentRepository(pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{pool: pool}
}

// ListEngineers retrieves engineers, optionally filtered by organization
func (r *AssignmentRepository) ListEngineers(ctx context.Context, organizationID *string, limit, offset int) ([]*domain.Engineer, error) {
	if limit <= 0 {
		limit = 100
	}
	
	query := `
		SELECT DISTINCT
			e.id, 
			COALESCE(eom.org_id::TEXT, '') as organization_id,
			COALESCE(o.name, '') as organization_name,
			e.name, 
			COALESCE(e.email, '') as email, 
			COALESCE(e.phone, '') as phone, 
			COALESCE(e.engineer_level, 1) as engineer_level,
			true as is_active, 
			e.created_at, 
			e.updated_at
		FROM engineers e
		LEFT JOIN engineer_org_memberships eom ON e.id = eom.engineer_id
		LEFT JOIN organizations o ON eom.org_id = o.id
		WHERE 1=1
	`
	
	args := []interface{}{}
	argPos := 1
	
	if organizationID != nil && *organizationID != "" {
		query += fmt.Sprintf(" AND eom.org_id = $%d", argPos)
		args = append(args, *organizationID)
		argPos++
	}
	
	query += fmt.Sprintf(" ORDER BY e.name LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)
	
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var engineers []*domain.Engineer
	for rows.Next() {
		var eng domain.Engineer
		var level int
		err := rows.Scan(
			&eng.ID, &eng.OrganizationID, &eng.OrganizationName,
			&eng.Name, &eng.Email, &eng.Phone,
			&level, &eng.IsActive, &eng.CreatedAt, &eng.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Convert int level to string format (L1, L2, L3)
		eng.EngineerLevel = domain.EngineerLevel(fmt.Sprintf("L%d", level))
		engineers = append(engineers, &eng)
	}
	
	return engineers, nil
}

// GetEngineerByID retrieves a single engineer by ID
func (r *AssignmentRepository) GetEngineerByID(ctx context.Context, engineerID string) (*domain.Engineer, error) {
	query := `
		SELECT 
			e.id,
			COALESCE(eom.org_id::TEXT, '') as organization_id,
			COALESCE(o.name, '') as organization_name,
			e.name,
			COALESCE(e.email, '') as email,
			COALESCE(e.phone, '') as phone,
			COALESCE(e.engineer_level, 1) as engineer_level,
			true as is_active,
			e.created_at,
			e.updated_at
		FROM engineers e
		LEFT JOIN engineer_org_memberships eom ON e.id = eom.engineer_id
		LEFT JOIN organizations o ON eom.org_id = o.id
		WHERE e.id = $1
		LIMIT 1
	`
	
	var eng domain.Engineer
	var level int
	err := r.pool.QueryRow(ctx, query, engineerID).Scan(
		&eng.ID, &eng.OrganizationID, &eng.OrganizationName,
		&eng.Name, &eng.Email, &eng.Phone,
		&level, &eng.IsActive, &eng.CreatedAt, &eng.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("engineer not found: %w", err)
	}
	
	// Convert int level to string format (L1, L2, L3)
	eng.EngineerLevel = domain.EngineerLevel(fmt.Sprintf("L%d", level))
	return &eng, nil
}

// UpdateEngineerLevel updates an engineer's skill level
func (r *AssignmentRepository) UpdateEngineerLevel(ctx context.Context, engineerID string, level domain.EngineerLevel) error {
	query := `UPDATE engineers SET engineer_level = $2, updated_at = $3 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, engineerID, string(level), time.Now())
	return err
}

// ListEngineerEquipmentTypes retrieves all equipment types an engineer can service
func (r *AssignmentRepository) ListEngineerEquipmentTypes(ctx context.Context, engineerID string) ([]*domain.EngineerEquipmentType, error) {
	query := `
		SELECT id, engineer_id, manufacturer_name as manufacturer, equipment_category as category, created_at
		FROM engineer_equipment_types
		WHERE engineer_id = $1
		ORDER BY manufacturer_name, equipment_category
	`
	
	rows, err := r.pool.Query(ctx, query, engineerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var types []*domain.EngineerEquipmentType
	for rows.Next() {
		var t domain.EngineerEquipmentType
		err := rows.Scan(&t.ID, &t.EngineerID, &t.Manufacturer, &t.Category, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		types = append(types, &t)
	}
	
	return types, nil
}

// AddEngineerEquipmentType adds an equipment type capability to an engineer
func (r *AssignmentRepository) AddEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error {
	// Check if already exists
	checkQuery := `SELECT COUNT(*) FROM engineer_equipment_types WHERE engineer_id = $1 AND manufacturer_name = $2 AND equipment_category = $3`
	var count int
	err := r.pool.QueryRow(ctx, checkQuery, engineerID, manufacturer, category).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already exists, no error
	}
	
	query := `
		INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.pool.Exec(ctx, query, engineerID, manufacturer, category, time.Now())
	return err
}

// RemoveEngineerEquipmentType removes an equipment type capability from an engineer
func (r *AssignmentRepository) RemoveEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error {
	query := `DELETE FROM engineer_equipment_types WHERE engineer_id = $1 AND manufacturer_name = $2 AND equipment_category = $3`
	_, err := r.pool.Exec(ctx, query, engineerID, manufacturer, category)
	return err
}

// GetEquipmentServiceConfig retrieves service configuration for equipment
func (r *AssignmentRepository) GetEquipmentServiceConfig(ctx context.Context, equipmentID string) (*domain.EquipmentServiceConfig, error) {
	query := `
		SELECT 
			id, equipment_id, under_warranty, under_amc,
			primary_service_org_id, secondary_service_org_id,
			tertiary_service_org_id, fallback_service_org_id,
			created_at, updated_at
		FROM equipment_service_config
		WHERE equipment_id = $1
	`
	
	var config domain.EquipmentServiceConfig
	err := r.pool.QueryRow(ctx, query, equipmentID).Scan(
		&config.ID, &config.EquipmentID, &config.UnderWarranty, &config.UnderAMC,
		&config.PrimaryServiceOrgID, &config.SecondaryServiceOrgID,
		&config.TertiaryServiceOrgID, &config.FallbackServiceOrgID,
		&config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("service config not found: %w", err)
	}
	
	return &config, nil
}

// CreateEquipmentServiceConfig creates a new service configuration
func (r *AssignmentRepository) CreateEquipmentServiceConfig(ctx context.Context, config *domain.EquipmentServiceConfig) error {
	if config.ID == "" {
		config.ID = ksuid.New().String()
	}
	
	query := `
		INSERT INTO equipment_service_config (
			id, equipment_id, under_warranty, under_amc,
			primary_service_org_id, secondary_service_org_id,
			tertiary_service_org_id, fallback_service_org_id,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	now := time.Now()
	_, err := r.pool.Exec(ctx, query,
		config.ID, config.EquipmentID, config.UnderWarranty, config.UnderAMC,
		config.PrimaryServiceOrgID, config.SecondaryServiceOrgID,
		config.TertiaryServiceOrgID, config.FallbackServiceOrgID,
		now, now,
	)
	
	return err
}

// UpdateEquipmentServiceConfig updates an existing service configuration
func (r *AssignmentRepository) UpdateEquipmentServiceConfig(ctx context.Context, config *domain.EquipmentServiceConfig) error {
	query := `
		UPDATE equipment_service_config SET
			under_warranty = $2, under_amc = $3,
			primary_service_org_id = $4, secondary_service_org_id = $5,
			tertiary_service_org_id = $6, fallback_service_org_id = $7,
			updated_at = $8
		WHERE id = $1
	`
	
	_, err := r.pool.Exec(ctx, query,
		config.ID, config.UnderWarranty, config.UnderAMC,
		config.PrimaryServiceOrgID, config.SecondaryServiceOrgID,
		config.TertiaryServiceOrgID, config.FallbackServiceOrgID,
		time.Now(),
	)
	
	return err
}

// GetSuggestedEngineers retrieves suggested engineers for a service ticket
// This implements the core assignment algorithm
func (r *AssignmentRepository) GetSuggestedEngineers(ctx context.Context, equipmentID string, manufacturer, category string, minLevel domain.EngineerLevel) ([]*domain.SuggestedEngineer, error) {
	// Step 1: Get eligible service organizations from equipment config
	eligibleOrgsQuery := `SELECT get_eligible_service_orgs($1)`
	var orgIDsJSON []byte
	err := r.pool.QueryRow(ctx, eligibleOrgsQuery, equipmentID).Scan(&orgIDsJSON)
	if err != nil {
		// If no config found, return empty list (no suggestions)
		return []*domain.SuggestedEngineer{}, nil
	}
	
	// Parse org IDs from JSON array
	orgIDsStr := strings.Trim(string(orgIDsJSON), "[]\"")
	if orgIDsStr == "" {
		return []*domain.SuggestedEngineer{}, nil
	}
	orgIDs := strings.Split(orgIDsStr, "\",\"")
	
	// Step 2: Find engineers in eligible orgs who can service this equipment type
	query := `
		SELECT DISTINCT
			e.id, e.name, e.organization_id, o.name as org_name,
			COALESCE(e.engineer_level, 'L1') as engineer_level
		FROM engineers e
		JOIN organizations o ON e.organization_id = o.id
		JOIN engineer_equipment_types eet ON e.id = eet.engineer_id
		WHERE e.is_active = true
			AND e.organization_id = ANY($1)
			AND eet.manufacturer = $2
			AND eet.category = $3
			AND COALESCE(e.engineer_level, 'L1') >= $4
		ORDER BY e.engineer_level DESC, e.name
	`
	
	rows, err := r.pool.Query(ctx, query, orgIDs, manufacturer, category, string(minLevel))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var suggestions []*domain.SuggestedEngineer
	priority := 1
	
	for rows.Next() {
		var suggestion domain.SuggestedEngineer
		var levelStr string
		
		err := rows.Scan(
			&suggestion.EngineerID, &suggestion.EngineerName,
			&suggestion.OrganizationID, &suggestion.OrganizationName,
			&levelStr,
		)
		if err != nil {
			return nil, err
		}
		
		suggestion.EngineerLevel = domain.EngineerLevel(levelStr)
		suggestion.AssignmentTier = r.determineAssignmentTier(equipmentID, suggestion.OrganizationID)
		suggestion.AssignmentTierName = r.formatTierName(suggestion.AssignmentTier)
		suggestion.MatchReason = fmt.Sprintf("%s %s engineer, Level %s", 
			manufacturer, category, levelStr)
		suggestion.Priority = priority
		priority++
		
		suggestions = append(suggestions, &suggestion)
	}
	
	return suggestions, nil
}

// determineAssignmentTier determines which tier an organization represents for the equipment
func (r *AssignmentRepository) determineAssignmentTier(equipmentID, orgID string) string {
	query := `
		SELECT 
			CASE 
				WHEN primary_service_org_id = $2 THEN 
					CASE WHEN under_warranty THEN 'warranty_primary'
						 WHEN under_amc THEN 'amc_primary'
						 ELSE 'primary' END
				WHEN secondary_service_org_id = $2 THEN 'secondary'
				WHEN tertiary_service_org_id = $2 THEN 'tertiary'
				WHEN fallback_service_org_id = $2 THEN 'fallback'
				ELSE 'unmatched'
			END as tier
		FROM equipment_service_config
		WHERE equipment_id = $1
	`
	
	var tier string
	err := r.pool.QueryRow(context.Background(), query, equipmentID, orgID).Scan(&tier)
	if err != nil {
		return "unmatched"
	}
	
	return tier
}

// formatTierName converts tier code to human-readable name
func (r *AssignmentRepository) formatTierName(tier string) string {
	switch tier {
	case "warranty_primary":
		return "Warranty Coverage"
	case "amc_primary":
		return "AMC Coverage"
	case "primary":
		return "Primary Service"
	case "secondary":
		return "Secondary Service"
	case "tertiary":
		return "Tertiary Service"
	case "fallback":
		return "Fallback Service"
	default:
		return "Other"
	}
}

// AssignEngineerToTicket assigns an engineer to a service ticket with full assignment data
func (r *AssignmentRepository) AssignEngineerToTicket(ctx context.Context, req domain.AssignmentRequest) error {
	now := time.Now()
	
	query := `
		UPDATE service_tickets SET
			assigned_engineer_id = $2,
			assigned_engineer_name = $3,
			assigned_org_id = $4,
			assignment_tier = $5,
			assignment_tier_name = $6,
			assigned_at = $7,
			status = 'assigned',
			updated_at = $7
		WHERE id = $1
	`
	
	_, err := r.pool.Exec(ctx, query,
		req.TicketID, req.EngineerID, req.EngineerName,
		req.OrganizationID, req.AssignmentTier, req.AssignmentTierName,
		now,
	)
	
	return err
}
