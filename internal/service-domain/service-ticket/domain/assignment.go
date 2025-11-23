package domain

import (
	"time"
)

// EngineerLevel represents engineer skill level
type EngineerLevel string

const (
	EngineerLevelL1 EngineerLevel = "L1"
	EngineerLevelL2 EngineerLevel = "L2"
	EngineerLevelL3 EngineerLevel = "L3"
)

// Engineer represents an engineer from the organizations table with assignment-specific data
type Engineer struct {
	ID                   string        `json:"id"`
	OrganizationID       string        `json:"organization_id"`
	OrganizationName     string        `json:"organization_name,omitempty"`
	Name                 string        `json:"name"`
	Email                string        `json:"email"`
	Phone                string        `json:"phone"`
	EngineerLevel        EngineerLevel `json:"engineer_level"`
	IsActive             bool          `json:"is_active"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	
	// For eligible engineers list
	EligibleEquipmentTypes []string `json:"eligible_equipment_types,omitempty"`
}

// EngineerEquipmentType maps engineers to equipment types they can service
type EngineerEquipmentType struct {
	ID            string    `json:"id"`
	EngineerID    string    `json:"engineer_id"`
	Manufacturer  string    `json:"manufacturer"`  // e.g., "Siemens"
	Category      string    `json:"category"`      // e.g., "MRI", "CT", "X-Ray"
	CreatedAt     time.Time `json:"created_at"`
}

// EquipmentServiceConfig defines service hierarchy for specific equipment
type EquipmentServiceConfig struct {
	ID                       string    `json:"id"`
	EquipmentID              string    `json:"equipment_id"`
	UnderWarranty            bool      `json:"under_warranty"`
	UnderAMC                 bool      `json:"under_amc"`
	PrimaryServiceOrgID      *string   `json:"primary_service_org_id"`
	SecondaryServiceOrgID    *string   `json:"secondary_service_org_id"`
	TertiaryServiceOrgID     *string   `json:"tertiary_service_org_id"`
	FallbackServiceOrgID     *string   `json:"fallback_service_org_id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

// SuggestedEngineer represents an engineer suggestion for assignment
type SuggestedEngineer struct {
	EngineerID         string        `json:"engineer_id"`
	EngineerName       string        `json:"engineer_name"`
	OrganizationID     string        `json:"organization_id"`
	OrganizationName   string        `json:"organization_name"`
	EngineerLevel      EngineerLevel `json:"engineer_level"`
	AssignmentTier     string        `json:"assignment_tier"`      // "warranty_primary", "amc_primary", "secondary", "tertiary", "fallback"
	AssignmentTierName string        `json:"assignment_tier_name"` // Human-readable tier name
	MatchReason        string        `json:"match_reason"`         // Why this engineer was suggested
	Priority           int           `json:"priority"`             // Lower is higher priority
}

// AssignmentRequest represents a manual assignment request
type AssignmentRequest struct {
	TicketID           string  `json:"ticket_id"`
	EngineerID         string  `json:"engineer_id"`
	EngineerName       string  `json:"engineer_name"`
	OrganizationID     string  `json:"organization_id"`
	AssignmentTier     string  `json:"assignment_tier"`
	AssignmentTierName string  `json:"assignment_tier_name"`
	AssignedBy         string  `json:"assigned_by"`
}
