package domain

import "context"

// Repository defines the interface for comparison persistence
type Repository interface {
	// Create creates a new comparison
	Create(ctx context.Context, comparison *Comparison) error

	// GetByID retrieves a comparison by ID
	GetByID(ctx context.Context, tenantID, id string) (*Comparison, error)

	// GetByRFQ retrieves all comparisons for an RFQ
	GetByRFQ(ctx context.Context, tenantID, rfqID string) ([]*Comparison, error)

	// List retrieves comparisons with filtering
	List(ctx context.Context, criteria ListCriteria) (*ListResult, error)

	// Update updates a comparison
	Update(ctx context.Context, comparison *Comparison) error

	// Delete deletes a comparison
	Delete(ctx context.Context, tenantID, id string) error
}

// ListCriteria defines filtering criteria for listing comparisons
type ListCriteria struct {
	TenantID    string
	RFQID       string
	Status      []ComparisonStatus
	CreatedBy   string
	SortBy      string // created_at, updated_at, title
	SortDirection string // asc, desc
	Page        int
	PageSize    int
}

// ListResult contains the paginated list of comparisons
type ListResult struct {
	Comparisons []*Comparison `json:"comparisons"`
	Total       int           `json:"total"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
	TotalPages  int           `json:"total_pages"`
}
