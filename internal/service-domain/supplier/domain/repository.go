package domain

import "context"

// SupplierRepository defines the interface for supplier data access
type SupplierRepository interface {
	// Create adds a new supplier
	Create(ctx context.Context, supplier *Supplier) error
	
	// GetByID retrieves a supplier by ID
	GetByID(ctx context.Context, id string, tenantID string) (*Supplier, error)
	
	// GetByTaxID retrieves a supplier by Tax ID
	GetByTaxID(ctx context.Context, taxID string, tenantID string) (*Supplier, error)
	
	// Update updates an existing supplier
	Update(ctx context.Context, supplier *Supplier) error
	
	// Delete removes a supplier
	Delete(ctx context.Context, id string, tenantID string) error
	
	// List retrieves suppliers with filtering and pagination
	List(ctx context.Context, criteria ListCriteria) ([]*Supplier, int, error)
	
	// GetByCategory retrieves suppliers specialized in a category
	GetByCategory(ctx context.Context, categoryID string, tenantID string) ([]*Supplier, error)
}

// ListCriteria defines filtering criteria for listing suppliers
type ListCriteria struct {
	TenantID           string
	Status             []SupplierStatus
	VerificationStatus []VerificationStatus
	CategoryID         string // Filter by specialization
	SearchQuery        string // Search in company name
	MinRating          float64
	Page               int
	PageSize           int
	SortBy             string // e.g., "company_name", "performance_rating", "created_at"
	SortDirection      string // "asc" or "desc"
}
