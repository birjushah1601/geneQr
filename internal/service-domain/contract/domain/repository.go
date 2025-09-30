package domain

import "context"

// Repository defines the interface for contract persistence
type Repository interface {
	// Create creates a new contract
	Create(ctx context.Context, contract *Contract) error

	// GetByID retrieves a contract by ID
	GetByID(ctx context.Context, tenantID, id string) (*Contract, error)

	// GetByContractNumber retrieves a contract by contract number
	GetByContractNumber(ctx context.Context, tenantID, contractNumber string) (*Contract, error)

	// GetByRFQ retrieves all contracts for an RFQ
	GetByRFQ(ctx context.Context, tenantID, rfqID string) ([]*Contract, error)

	// GetBySupplier retrieves all contracts for a supplier
	GetBySupplier(ctx context.Context, tenantID, supplierID string) ([]*Contract, error)

	// List retrieves contracts with filtering
	List(ctx context.Context, criteria ListCriteria) (*ListResult, error)

	// Update updates a contract
	Update(ctx context.Context, contract *Contract) error

	// Delete deletes a contract
	Delete(ctx context.Context, tenantID, id string) error
}

// ListCriteria defines filtering criteria for listing contracts
type ListCriteria struct {
	TenantID     string
	RFQID        string
	SupplierID   string
	Status       []ContractStatus
	CreatedBy    string
	StartDateFrom *string
	StartDateTo   *string
	SortBy       string // created_at, updated_at, start_date, total_amount
	SortDirection string // asc, desc
	Page         int
	PageSize     int
}

// ListResult contains the paginated list of contracts
type ListResult struct {
	Contracts  []*Contract `json:"contracts"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
