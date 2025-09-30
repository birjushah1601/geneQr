package domain

import (
	"context"
)

// RFQRepository defines the interface for RFQ persistence
type RFQRepository interface {
	// Create persists a new RFQ
	Create(ctx context.Context, rfq *RFQ) error
	
	// GetByID retrieves an RFQ by ID
	GetByID(ctx context.Context, id string, tenantID string) (*RFQ, error)
	
	// GetByRFQNumber retrieves an RFQ by its RFQ number
	GetByRFQNumber(ctx context.Context, rfqNumber string, tenantID string) (*RFQ, error)
	
	// Update updates an existing RFQ
	Update(ctx context.Context, rfq *RFQ) error
	
	// Delete removes an RFQ
	Delete(ctx context.Context, id string, tenantID string) error
	
	// List retrieves RFQs with pagination and filtering
	List(ctx context.Context, criteria ListCriteria) ([]*RFQ, int, error)
	
	// AddItem adds an item to an RFQ
	AddItem(ctx context.Context, rfqID string, item *RFQItem) error
	
	// UpdateItem updates an RFQ item
	UpdateItem(ctx context.Context, item *RFQItem) error
	
	// RemoveItem removes an item from an RFQ
	RemoveItem(ctx context.Context, rfqID, itemID string) error
	
	// GetItems retrieves all items for an RFQ
	GetItems(ctx context.Context, rfqID string) ([]RFQItem, error)
	
	// AddInvitation adds a supplier invitation
	AddInvitation(ctx context.Context, invitation *RFQInvitation) error
	
	// GetInvitations retrieves all invitations for an RFQ
	GetInvitations(ctx context.Context, rfqID string) ([]RFQInvitation, error)
	
	// UpdateInvitation updates an invitation status
	UpdateInvitation(ctx context.Context, invitation *RFQInvitation) error
}

// ListCriteria defines filtering criteria for listing RFQs
type ListCriteria struct {
	TenantID        string
	Status          []RFQStatus
	Priority        []RFQPriority
	CreatedBy       string
	SearchQuery     string
	FromDate        *string
	ToDate          *string
	Page            int
	PageSize        int
	SortBy          string
	SortDirection   string
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event interface{}) error
}
