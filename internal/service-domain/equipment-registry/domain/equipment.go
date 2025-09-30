package domain

import (
	"errors"
	"time"
)

var (
	ErrEquipmentNotFound     = errors.New("equipment not found")
	ErrInvalidSerialNumber   = errors.New("invalid serial number")
	ErrDuplicateSerialNumber = errors.New("duplicate serial number")
	ErrInvalidQRCode         = errors.New("invalid QR code")
)

// EquipmentStatus represents the operational status of equipment
type EquipmentStatus string

const (
	StatusOperational      EquipmentStatus = "operational"
	StatusDown             EquipmentStatus = "down"
	StatusUnderMaintenance EquipmentStatus = "under_maintenance"
	StatusDecommissioned   EquipmentStatus = "decommissioned"
)

// Equipment represents a registered medical equipment in the field
type Equipment struct {
	ID           string          `json:"id"`
	QRCode       string          `json:"qr_code"`        // Unique QR identifier
	SerialNumber string          `json:"serial_number"`  // Manufacturer serial number
	
	// Equipment details
	EquipmentID      string `json:"equipment_id"`       // Link to catalog
	EquipmentName    string `json:"equipment_name"`
	ManufacturerName string `json:"manufacturer_name"`
	ModelNumber      string `json:"model_number"`
	Category         string `json:"category"`
	
	// Installation details
	CustomerID          string                 `json:"customer_id"`
	CustomerName        string                 `json:"customer_name"`
	InstallationLocation string                `json:"installation_location"`
	InstallationAddress  map[string]interface{} `json:"installation_address"`
	InstallationDate     *time.Time            `json:"installation_date"`
	
	// Contract details
	ContractID    string     `json:"contract_id,omitempty"`    // Link to procurement contract
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice float64    `json:"purchase_price"`
	WarrantyExpiry *time.Time `json:"warranty_expiry,omitempty"`
	AMCContractID string     `json:"amc_contract_id,omitempty"`
	
	// Status and service
	Status           EquipmentStatus `json:"status"`
	LastServiceDate  *time.Time      `json:"last_service_date,omitempty"`
	NextServiceDate  *time.Time      `json:"next_service_date,omitempty"`
	ServiceCount     int             `json:"service_count"`
	
	// Technical details
	Specifications map[string]interface{} `json:"specifications"`
	Photos         []string               `json:"photos"`  // Array of photo URLs
	Documents      []string               `json:"documents"` // Manuals, certificates
	
	// QR Code URL
	QRCodeURL string `json:"qr_code_url"` // URL encoded in QR code
	
	// Metadata
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
}

// NewEquipment creates a new equipment registration
func NewEquipment(serialNumber, equipmentName, manufacturerName, modelNumber, customerName, createdBy string) *Equipment {
	now := time.Now()
	return &Equipment{
		SerialNumber:     serialNumber,
		EquipmentName:    equipmentName,
		ManufacturerName: manufacturerName,
		ModelNumber:      modelNumber,
		CustomerName:     customerName,
		Status:           StatusOperational,
		ServiceCount:     0,
		Specifications:   make(map[string]interface{}),
		Photos:           []string{},
		Documents:        []string{},
		CreatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// MarkAsDown marks equipment as non-operational
func (e *Equipment) MarkAsDown() {
	e.Status = StatusDown
	e.UpdatedAt = time.Now()
}

// MarkUnderMaintenance marks equipment as under maintenance
func (e *Equipment) MarkUnderMaintenance() {
	e.Status = StatusUnderMaintenance
	e.UpdatedAt = time.Now()
}

// MarkAsOperational marks equipment as operational
func (e *Equipment) MarkAsOperational() {
	e.Status = StatusOperational
	e.UpdatedAt = time.Now()
}

// RecordService records a service completion
func (e *Equipment) RecordService(serviceDate time.Time) {
	e.LastServiceDate = &serviceDate
	e.ServiceCount++
	e.Status = StatusOperational
	e.UpdatedAt = time.Now()
}

// ScheduleNextService sets the next service date
func (e *Equipment) ScheduleNextService(nextDate time.Time) {
	e.NextServiceDate = &nextDate
	e.UpdatedAt = time.Now()
}

// Decommission marks equipment as decommissioned
func (e *Equipment) Decommission() {
	e.Status = StatusDecommissioned
	e.UpdatedAt = time.Now()
}

// IsUnderWarranty checks if equipment is still under warranty
func (e *Equipment) IsUnderWarranty() bool {
	if e.WarrantyExpiry == nil {
		return false
	}
	return time.Now().Before(*e.WarrantyExpiry)
}

// HasAMC checks if equipment has an active AMC
func (e *Equipment) HasAMC() bool {
	return e.AMCContractID != ""
}

// CSVImportRow represents a row from CSV import
type CSVImportRow struct {
	SerialNumber     string  `csv:"serial_number"`
	EquipmentName    string  `csv:"equipment_name"`
	ManufacturerName string  `csv:"manufacturer_name"`
	ModelNumber      string  `csv:"model_number"`
	Category         string  `csv:"category"`
	CustomerName     string  `csv:"customer_name"`
	CustomerID       string  `csv:"customer_id"`
	InstallationLocation string `csv:"installation_location"`
	InstallationDate string  `csv:"installation_date"` // Format: YYYY-MM-DD
	PurchaseDate     string  `csv:"purchase_date"`     // Format: YYYY-MM-DD
	PurchasePrice    float64 `csv:"purchase_price"`
	WarrantyMonths   int     `csv:"warranty_months"`
	Notes            string  `csv:"notes"`
}

// CSVImportResult contains results of CSV import
type CSVImportResult struct {
	TotalRows      int      `json:"total_rows"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	Errors         []string `json:"errors"`
	ImportedIDs    []string `json:"imported_ids"`
}
