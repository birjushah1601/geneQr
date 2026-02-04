package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// PartnerService handles partner association operations
type PartnerService struct {
	db *sql.DB
}

// NewPartnerService creates a new partner service
func NewPartnerService(db *sql.DB) *PartnerService {
	return &PartnerService{db: db}
}

// Partner represents a partner organization association
type Partner struct {
	ID              string     `json:"id"`
	PartnerOrgID    string     `json:"partner_org_id"`
	PartnerName     string     `json:"partner_name"`
	OrgType         string     `json:"org_type"`
	EquipmentID     *string    `json:"equipment_id,omitempty"`
	EquipmentName   *string    `json:"equipment_name,omitempty"`
	AssociationType string     `json:"association_type"`
	EngineersCount  int        `json:"engineers_count"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Organization represents an available organization
type Organization struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrgType        string `json:"org_type"`
	Location       string `json:"location,omitempty"`
	EngineersCount int    `json:"engineers_count"`
	ContactEmail   string `json:"contact_email,omitempty"`
}

// CreateAssociationRequest represents a request to create a partner association
type CreateAssociationRequest struct {
	ManufacturerID string  `json:"manufacturer_id"`
	PartnerOrgID   string  `json:"partner_org_id"`
	EquipmentID    *string `json:"equipment_id,omitempty"`
	RelType        string  `json:"rel_type"`
}

// PartnerAssociation represents a created association
type PartnerAssociation struct {
	ID            string     `json:"id"`
	ParentOrgID   string     `json:"parent_org_id"`
	ChildOrgID    string     `json:"child_org_id"`
	PartnerName   string     `json:"partner_name"`
	OrgType       string     `json:"org_type"`
	EquipmentID   *string    `json:"equipment_id,omitempty"`
	EquipmentName *string    `json:"equipment_name,omitempty"`
	RelType       string     `json:"rel_type"`
	CreatedAt     time.Time  `json:"created_at"`
}

// GetPartners returns partners associated with a manufacturer
func (s *PartnerService) GetPartners(ctx context.Context, manufacturerID string, filters map[string]string) ([]Partner, error) {
	query := `
		SELECT 
			r.id,
			r.child_org_id as partner_id,
			o.name as partner_name,
			o.org_type,
			r.equipment_id,
			e.equipment_name,
			CASE 
				WHEN r.equipment_id IS NULL THEN 'general'
				ELSE 'equipment-specific'
			END as association_type,
			COUNT(DISTINCT eom.engineer_id) as engineers_count,
			r.created_at
		FROM org_relationships r
		JOIN organizations o ON o.id = r.child_org_id
		LEFT JOIN equipment e ON e.id = r.equipment_id
		LEFT JOIN engineer_org_memberships eom ON eom.org_id = o.id
		WHERE r.parent_org_id = $1 
		AND r.rel_type = 'services_for'
		AND o.org_type IN ('channel_partner', 'sub_dealer')
	`

	args := []interface{}{manufacturerID}
	argCount := 1

	// Add filters
	if orgType, ok := filters["type"]; ok && orgType != "" {
		argCount++
		query += fmt.Sprintf(" AND o.org_type = $%d", argCount)
		args = append(args, orgType)
	}

	if assocType, ok := filters["association_type"]; ok && assocType != "" {
		if assocType == "general" {
			query += " AND r.equipment_id IS NULL"
		} else if assocType == "equipment-specific" {
			query += " AND r.equipment_id IS NOT NULL"
		}
	}

	query += " GROUP BY r.id, r.child_org_id, o.name, o.org_type, r.equipment_id, e.equipment_name, r.created_at"
	query += " ORDER BY o.org_type, o.name"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query partners: %w", err)
	}
	defer rows.Close()

	var partners []Partner
	for rows.Next() {
		var p Partner
		var equipmentID, equipmentName sql.NullString
		
		err := rows.Scan(
			&p.ID,
			&p.PartnerOrgID,
			&p.PartnerName,
			&p.OrgType,
			&equipmentID,
			&equipmentName,
			&p.AssociationType,
			&p.EngineersCount,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan partner: %w", err)
		}

		if equipmentID.Valid {
			p.EquipmentID = &equipmentID.String
			if equipmentName.Valid {
				p.EquipmentName = &equipmentName.String
			}
		}

		partners = append(partners, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate partners: %w", err)
	}

	return partners, nil
}

// GetAvailablePartners returns partners not yet associated with manufacturer
func (s *PartnerService) GetAvailablePartners(ctx context.Context, manufacturerID, search string) ([]Organization, error) {
	query := `
		SELECT 
			o.id,
			o.name,
			o.org_type,
			COALESCE(o.metadata->>'city', '') as location,
			COUNT(DISTINCT eom.engineer_id) as engineers_count,
			COALESCE(o.metadata->>'contact_email', '') as contact_email
		FROM organizations o
		LEFT JOIN engineer_org_memberships eom ON eom.org_id = o.id
		WHERE o.org_type IN ('channel_partner', 'sub_dealer')
		AND o.id NOT IN (
			SELECT child_org_id 
			FROM org_relationships 
			WHERE parent_org_id = $1 
			AND rel_type = 'services_for'
			AND equipment_id IS NULL
		)
	`

	args := []interface{}{manufacturerID}

	if search != "" {
		query += " AND (o.name ILIKE $2 OR o.metadata->>'city' ILIKE $2)"
		args = append(args, "%"+search+"%")
	}

	query += " GROUP BY o.id, o.name, o.org_type, o.metadata"
	query += " ORDER BY o.name LIMIT 50"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query available partners: %w", err)
	}
	defer rows.Close()

	var orgs []Organization
	for rows.Next() {
		var org Organization
		
		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.OrgType,
			&org.Location,
			&org.EngineersCount,
			&org.ContactEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("scan organization: %w", err)
		}

		orgs = append(orgs, org)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate organizations: %w", err)
	}

	return orgs, nil
}

// CreateAssociation creates a new partner association
func (s *PartnerService) CreateAssociation(ctx context.Context, req CreateAssociationRequest) (*PartnerAssociation, error) {
	// Validate partner org_type
	var orgType string
	err := s.db.QueryRowContext(ctx,
		"SELECT org_type FROM organizations WHERE id = $1",
		req.PartnerOrgID,
	).Scan(&orgType)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("partner organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query partner org_type: %w", err)
	}

	if orgType != "channel_partner" && orgType != "sub_dealer" {
		return nil, fmt.Errorf("invalid partner org_type: %s (must be channel_partner or sub_dealer)", orgType)
	}

	// If equipment_id provided, verify it belongs to manufacturer
	if req.EquipmentID != nil && *req.EquipmentID != "" {
		var mfgID sql.NullString
		err := s.db.QueryRowContext(ctx,
			"SELECT manufacturer_id FROM equipment WHERE id = $1",
			*req.EquipmentID,
		).Scan(&mfgID)
		
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("equipment not found")
		}
		if err != nil {
			return nil, fmt.Errorf("query equipment: %w", err)
		}
		
		if !mfgID.Valid || mfgID.String != req.ManufacturerID {
			return nil, fmt.Errorf("equipment does not belong to manufacturer")
		}
	}

	// Set default rel_type if not provided
	if req.RelType == "" {
		req.RelType = "services_for"
	}

	// Create association
	var id string
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO org_relationships 
		(parent_org_id, child_org_id, rel_type, equipment_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, req.ManufacturerID, req.PartnerOrgID, req.RelType, req.EquipmentID).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("create association: %w", err)
	}

	// Fetch created association with details
	var assoc PartnerAssociation
	var equipmentID, equipmentName sql.NullString
	
	err = s.db.QueryRowContext(ctx, `
		SELECT 
			r.id,
			r.parent_org_id,
			r.child_org_id,
			o.name as partner_name,
			o.org_type,
			r.equipment_id,
			e.equipment_name,
			r.rel_type,
			r.created_at
		FROM org_relationships r
		JOIN organizations o ON o.id = r.child_org_id
		LEFT JOIN equipment e ON e.id = r.equipment_id
		WHERE r.id = $1
	`, id).Scan(
		&assoc.ID,
		&assoc.ParentOrgID,
		&assoc.ChildOrgID,
		&assoc.PartnerName,
		&assoc.OrgType,
		&equipmentID,
		&equipmentName,
		&assoc.RelType,
		&assoc.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("fetch association: %w", err)
	}

	if equipmentID.Valid {
		assoc.EquipmentID = &equipmentID.String
		if equipmentName.Valid {
			assoc.EquipmentName = &equipmentName.String
		}
	}

	return &assoc, nil
}

// RemoveAssociation removes a partner association
func (s *PartnerService) RemoveAssociation(ctx context.Context, manufacturerID, partnerID string, equipmentID *string) error {
	query := `
		DELETE FROM org_relationships
		WHERE parent_org_id = $1 
		AND child_org_id = $2 
		AND rel_type = 'services_for'
	`

	args := []interface{}{manufacturerID, partnerID}

	if equipmentID != nil && *equipmentID != "" {
		query += " AND equipment_id = $3"
		args = append(args, *equipmentID)
	} else {
		query += " AND equipment_id IS NULL"
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("remove association: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("association not found")
	}

	return nil
}
