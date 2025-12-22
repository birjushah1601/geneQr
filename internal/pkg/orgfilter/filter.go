package orgfilter

import (
	"context"
	"fmt"

	"github.com/aby-med/medical-platform/internal/middleware"
	"github.com/google/uuid"
)

// OrgContext holds organization context extracted from request
type OrgContext struct {
	OrgID   uuid.UUID
	OrgType string
	Role    string
}

// GetOrgContext extracts organization context from request context
func GetOrgContext(ctx context.Context) (*OrgContext, error) {
	orgID, ok := middleware.GetOrganizationID(ctx)
	if !ok {
		return nil, fmt.Errorf("organization ID not found in context")
	}

	orgType, _ := middleware.GetOrganizationType(ctx)
	role, _ := middleware.GetUserRole(ctx)

	return &OrgContext{
		OrgID:   orgID,
		OrgType: orgType,
		Role:    role,
	}, nil
}

// EquipmentFilter builds WHERE clause for equipment queries based on organization type
// Returns: SQL condition string and parameter value
func EquipmentFilter(orgCtx *OrgContext) (string, uuid.UUID) {
	switch orgCtx.OrgType {
	case "manufacturer":
		// Manufacturers see ALL equipment they manufactured (across all organizations)
		return "manufacturer_id = $%d", orgCtx.OrgID

	case "hospital", "imaging_center":
		// Hospitals see ONLY equipment they own
		return "(organization_id = $%d OR owner_org_id = $%d)", orgCtx.OrgID

	case "distributor", "dealer":
		// Distributors see equipment they sold/service
		return "(distributor_org_id = $%d OR service_provider_org_id = $%d)", orgCtx.OrgID

	case "supplier":
		// Suppliers see equipment where they supply parts
		return "supplier_org_id = $%d", orgCtx.OrgID

	default:
		// Default: only owned equipment
		return "organization_id = $%d", orgCtx.OrgID
	}
}

// TicketFilter builds WHERE clause for ticket queries based on organization type
func TicketFilter(orgCtx *OrgContext) (string, uuid.UUID) {
	switch orgCtx.OrgType {
	case "manufacturer":
		// Manufacturers see tickets for their equipment
		return `EXISTS (
			SELECT 1 FROM equipment_registry e 
			WHERE e.id = service_tickets.equipment_id 
			AND e.manufacturer_id = $%d
		)`, orgCtx.OrgID

	case "hospital", "imaging_center":
		// Hospitals see tickets they created
		return "requester_org_id = $%d", orgCtx.OrgID

	case "distributor", "dealer", "supplier":
		// Service providers see tickets assigned to them
		return "(assigned_org_id = $%d OR service_provider_org_id = $%d)", orgCtx.OrgID

	default:
		return "requester_org_id = $%d", orgCtx.OrgID
	}
}

// EngineerFilter builds WHERE clause for engineer queries
func EngineerFilter(orgCtx *OrgContext) (string, uuid.UUID) {
	// Engineers belong to organizations through engineer_org_memberships
	return `EXISTS (
		SELECT 1 FROM engineer_org_memberships eom 
		WHERE eom.engineer_id = engineers.id 
		AND eom.org_id = $%d 
		AND eom.status = 'active'
	)`, orgCtx.OrgID
}

// BuildWhereClause constructs complete WHERE clause with org filter
// filterTemplate example: "manufacturer_id = $%d"
// paramIndex: current parameter index (for PostgreSQL $1, $2, etc.)
func BuildWhereClause(filterTemplate string, paramIndex int, additionalConditions ...string) string {
	// Replace %d with actual parameter index
	orgFilter := fmt.Sprintf(filterTemplate, paramIndex)

	where := "WHERE " + orgFilter

	// Add additional conditions
	for _, condition := range additionalConditions {
		if condition != "" {
			where += " AND " + condition
		}
	}

	return where
}

// IsSystemAdmin checks if user has system admin role (can see all orgs)
func IsSystemAdmin(ctx context.Context) bool {
	role, ok := middleware.GetUserRole(ctx)
	if !ok {
		return false
	}
	return role == "system_admin" || role == "super_admin"
}
