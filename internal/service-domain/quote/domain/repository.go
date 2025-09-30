package domain

import "context"

// QuoteRepository defines the interface for quote persistence
type QuoteRepository interface {
	// Create persists a new quote
	Create(ctx context.Context, quote *Quote) error

	// GetByID retrieves a quote by ID
	GetByID(ctx context.Context, id string, tenantID string) (*Quote, error)

	// GetByRFQID retrieves all quotes for an RFQ
	GetByRFQID(ctx context.Context, rfqID string, tenantID string) ([]*Quote, error)

	// GetBySupplierID retrieves all quotes from a supplier
	GetBySupplierID(ctx context.Context, supplierID string, tenantID string) ([]*Quote, error)

	// List retrieves quotes with filtering and pagination
	List(ctx context.Context, criteria ListCriteria) ([]*Quote, int, error)

	// Update updates an existing quote
	Update(ctx context.Context, quote *Quote) error

	// Delete removes a quote
	Delete(ctx context.Context, id string, tenantID string) error
}

// ListCriteria defines filtering and pagination options for listing quotes
type ListCriteria struct {
	TenantID       string        `json:"tenant_id"`
	RFQID          string        `json:"rfq_id,omitempty"`
	SupplierID     string        `json:"supplier_id,omitempty"`
	Status         []QuoteStatus `json:"status,omitempty"`
	MinAmount      float64       `json:"min_amount,omitempty"`
	MaxAmount      float64       `json:"max_amount,omitempty"`
	SortBy         string        `json:"sort_by,omitempty"`
	SortDirection  string        `json:"sort_direction,omitempty"` // asc or desc
	Page           int           `json:"page,omitempty"`
	PageSize       int           `json:"page_size,omitempty"`
}
