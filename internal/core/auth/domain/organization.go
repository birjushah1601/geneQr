package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization in the system
type Organization struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Type       string    `json:"org_type" db:"org_type"` // manufacturer, hospital, distributor, dealer, supplier, imaging_center
	Status     string    `json:"status" db:"status"`
	ExternalRef *string   `json:"external_ref" db:"external_ref"`
	Metadata   JSONBMap  `json:"metadata" db:"metadata"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// OrganizationRepository provides data access for organizations
type OrganizationRepository interface {
	// GetByID retrieves an organization by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	
	// GetByUserID retrieves organizations for a specific user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Organization, error)
}

// Organization types
const (
	OrgTypeManufacturer  = "manufacturer"
	OrgTypeHospital      = "hospital"
	OrgTypeDistributor   = "distributor"
	OrgTypeDealer        = "dealer"
	OrgTypeSupplier      = "supplier"
	OrgTypeImagingCenter = "imaging_center"
)

// Organization status
const (
	OrgStatusActive    = "active"
	OrgStatusInactive  = "inactive"
	OrgStatusSuspended = "suspended"
)
