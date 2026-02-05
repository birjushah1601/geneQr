package services

import (
	"context"
	"database/sql"
	"fmt"
)

// NetworkEngineersService handles network engineer queries with smart filtering
type NetworkEngineersService struct {
	db *sql.DB
}

// NewNetworkEngineersService creates a new network engineers service
func NewNetworkEngineersService(db *sql.DB) *NetworkEngineersService {
	return &NetworkEngineersService{db: db}
}

// Engineer represents an engineer with organization details
type Engineer struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Organization OrgInfo `json:"organization"`
	Category     string  `json:"category"`
}

// OrgInfo represents organization information
type OrgInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OrgType string `json:"org_type"`
}

// EngineerInfo represents simplified engineer information for grouping
type EngineerInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// NetworkEngineersResponse represents the response for network engineers query
type NetworkEngineersResponse struct {
	EquipmentID     *string                   `json:"equipment_id,omitempty"`
	Engineers       []Engineer                `json:"engineers"`
	Grouped         map[string][]EngineerInfo `json:"grouped"`
	TotalEngineers  int                       `json:"total_engineers"`
	AssociationType string                    `json:"association_type"`
}

// GetNetworkEngineers returns engineers from manufacturer's network with smart filtering
// 
// Smart Filtering Logic:
// 1. If equipment_id provided:
//    - Check for equipment-specific partner associations
//    - If found: Return ONLY engineers from those partners + manufacturer
//    - If not found: Return engineers from general partners + manufacturer
// 2. If no equipment_id:
//    - Return all general partners + manufacturer engineers
func (s *NetworkEngineersService) GetNetworkEngineers(ctx context.Context, manufacturerID string, equipmentID *string) (*NetworkEngineersResponse, error) {
	query := `
		WITH equipment_partners AS (
			SELECT child_org_id 
			FROM org_relationships
			WHERE parent_org_id = $1 
			AND equipment_id = $2 
			AND rel_type = 'services_for'
		),
		general_partners AS (
			SELECT child_org_id 
			FROM org_relationships
			WHERE parent_org_id = $1 
			AND equipment_id IS NULL 
			AND rel_type = 'services_for'
		)
		SELECT 
			e.id,
			e.name,
			COALESCE(e.phone, '') as phone,
			COALESCE(e.email, '') as email,
			o.id as org_id,
			o.name as org_name,
			o.org_type,
			CASE o.org_type
				WHEN 'manufacturer' THEN 'Manufacturer'
				WHEN 'channel_partner' THEN 'Channel Partner'
				WHEN 'sub_dealer' THEN 'Sub-Dealer'
				ELSE 'Other'
			END as category
		FROM engineers e
		JOIN engineer_org_memberships eom ON eom.engineer_id = e.id
		JOIN organizations o ON o.id = eom.org_id
		WHERE 
			-- If equipment-specific partners exist, use only those
			(EXISTS (SELECT 1 FROM equipment_partners) AND o.id IN (SELECT child_org_id FROM equipment_partners))
			-- Otherwise, use general partners + manufacturer
			OR (NOT EXISTS (SELECT 1 FROM equipment_partners) AND (
				o.id = $1 OR o.id IN (SELECT child_org_id FROM general_partners)
			))
		ORDER BY o.org_type, o.name, e.name
	`

	args := []interface{}{manufacturerID}
	if equipmentID != nil && *equipmentID != "" {
		args = append(args, *equipmentID)
	} else {
		args = append(args, nil)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query network engineers: %w", err)
	}
	defer rows.Close()

	engineers := []Engineer{}
	grouped := make(map[string][]EngineerInfo)

	for rows.Next() {
		var eng Engineer
		err := rows.Scan(
			&eng.ID,
			&eng.Name,
			&eng.Phone,
			&eng.Email,
			&eng.Organization.ID,
			&eng.Organization.Name,
			&eng.Organization.OrgType,
			&eng.Category,
		)
		if err != nil {
			return nil, fmt.Errorf("scan engineer: %w", err)
		}

		engineers = append(engineers, eng)

		// Group engineers by category
		groupKey := eng.Category
		if eng.Category != "Manufacturer" {
			// For partners, include organization name in group key
			groupKey = fmt.Sprintf("%s - %s", eng.Category, eng.Organization.Name)
		}

		grouped[groupKey] = append(grouped[groupKey], EngineerInfo{
			ID:    eng.ID,
			Name:  eng.Name,
			Phone: eng.Phone,
			Email: eng.Email,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate engineers: %w", err)
	}

	// Determine association type
	associationType := "general"
	if equipmentID != nil && *equipmentID != "" && len(engineers) > 0 {
		// Check if we used equipment-specific associations
		var count int
		err = s.db.QueryRowContext(ctx, `
			SELECT COUNT(*) 
			FROM org_relationships
			WHERE parent_org_id = $1 
			AND equipment_id = $2 
			AND rel_type = 'services_for'
		`, manufacturerID, *equipmentID).Scan(&count)
		
		if err == nil && count > 0 {
			associationType = "equipment-specific"
		}
	}

	return &NetworkEngineersResponse{
		EquipmentID:     equipmentID,
		Engineers:       engineers,
		Grouped:         grouped,
		TotalEngineers:  len(engineers),
		AssociationType: associationType,
	}, nil
}

// GetNetworkEngineersByManufacturer returns all engineers in manufacturer's network (no equipment filter)
func (s *NetworkEngineersService) GetNetworkEngineersByManufacturer(ctx context.Context, manufacturerID string) (*NetworkEngineersResponse, error) {
	return s.GetNetworkEngineers(ctx, manufacturerID, nil)
}

// GetNetworkEngineersForEquipment returns engineers for specific equipment (with override logic)
func (s *NetworkEngineersService) GetNetworkEngineersForEquipment(ctx context.Context, manufacturerID, equipmentID string) (*NetworkEngineersResponse, error) {
	return s.GetNetworkEngineers(ctx, manufacturerID, &equipmentID)
}
