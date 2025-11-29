package domain

import "context"

// Repository defines the interface for equipment persistence
type Repository interface {
	// Create creates a new equipment registration
	Create(ctx context.Context, equipment *Equipment) error

	// GetByID retrieves equipment by ID
	GetByID(ctx context.Context, id string) (*Equipment, error)

	// GetByQRCode retrieves equipment by QR code
	GetByQRCode(ctx context.Context, qrCode string) (*Equipment, error)

	// GetBySerialNumber retrieves equipment by serial number
	GetBySerialNumber(ctx context.Context, serialNumber string) (*Equipment, error)

	// List retrieves equipment with filtering
	List(ctx context.Context, criteria ListCriteria) (*ListResult, error)

	// Update updates equipment
	Update(ctx context.Context, equipment *Equipment) error

	// Delete deletes equipment
	Delete(ctx context.Context, id string) error

	// BulkCreate creates multiple equipment registrations
	BulkCreate(ctx context.Context, equipment []*Equipment) error
	
	// UpdateQRCode updates the QR code image in database
	UpdateQRCode(ctx context.Context, equipmentID string, qrImage []byte, format string) error

    // SetQRCodeByID maps a QR code (and URL) to an equipment by ID
    SetQRCodeByID(ctx context.Context, id, qrCode, qrURL string) error

    // SetQRCodeBySerial maps a QR code (and URL) to an equipment by serial number
    SetQRCodeBySerial(ctx context.Context, serial, qrCode, qrURL string) error
}

// ListCriteria defines filtering criteria for listing equipment
type ListCriteria struct {
	CustomerID       string
	ManufacturerName string
	Status           []EquipmentStatus
	Category         string
	HasAMC           *bool
	UnderWarranty    *bool
	SortBy           string // created_at, serial_number, installation_date
	SortDirection    string // asc, desc
	Page             int
	PageSize         int
}

// ListResult contains the paginated list of equipment
type ListResult struct {
	Equipment  []*Equipment `json:"equipment"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
