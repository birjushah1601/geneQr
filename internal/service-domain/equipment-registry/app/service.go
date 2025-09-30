package app

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/qrcode"
	"github.com/segmentio/ksuid"
)

// EquipmentService handles equipment business logic
type EquipmentService struct {
	repo        domain.Repository
	qrGenerator *qrcode.Generator
	logger      *slog.Logger
	baseURL     string
}

// NewEquipmentService creates a new equipment service
func NewEquipmentService(repo domain.Repository, qrGenerator *qrcode.Generator, logger *slog.Logger, baseURL string) *EquipmentService {
	return &EquipmentService{
		repo:        repo,
		qrGenerator: qrGenerator,
		logger:      logger.With(slog.String("component", "equipment_service")),
		baseURL:     baseURL,
	}
}

// RegisterEquipmentRequest represents equipment registration request
type RegisterEquipmentRequest struct {
	SerialNumber         string                 `json:"serial_number"`
	EquipmentID          string                 `json:"equipment_id,omitempty"`
	EquipmentName        string                 `json:"equipment_name"`
	ManufacturerName     string                 `json:"manufacturer_name"`
	ModelNumber          string                 `json:"model_number,omitempty"`
	Category             string                 `json:"category,omitempty"`
	CustomerID           string                 `json:"customer_id,omitempty"`
	CustomerName         string                 `json:"customer_name"`
	InstallationLocation string                 `json:"installation_location,omitempty"`
	InstallationAddress  map[string]interface{} `json:"installation_address,omitempty"`
	InstallationDate     *time.Time             `json:"installation_date,omitempty"`
	ContractID           string                 `json:"contract_id,omitempty"`
	PurchaseDate         *time.Time             `json:"purchase_date,omitempty"`
	PurchasePrice        float64                `json:"purchase_price,omitempty"`
	WarrantyMonths       int                    `json:"warranty_months,omitempty"`
	AMCContractID        string                 `json:"amc_contract_id,omitempty"`
	Specifications       map[string]interface{} `json:"specifications,omitempty"`
	Notes                string                 `json:"notes,omitempty"`
	CreatedBy            string                 `json:"created_by"`
}

// RegisterEquipment registers a new equipment
func (s *EquipmentService) RegisterEquipment(ctx context.Context, req RegisterEquipmentRequest) (*domain.Equipment, error) {
	// Generate IDs
	equipmentID := ksuid.New().String()
	qrCodeID := s.generateQRCodeID()

	// Create equipment entity
	equipment := domain.NewEquipment(
		req.SerialNumber,
		req.EquipmentName,
		req.ManufacturerName,
		req.ModelNumber,
		req.CustomerName,
		req.CreatedBy,
	)

	equipment.ID = equipmentID
	equipment.QRCode = qrCodeID
	equipment.EquipmentID = req.EquipmentID
	equipment.Category = req.Category
	equipment.CustomerID = req.CustomerID
	equipment.InstallationLocation = req.InstallationLocation
	equipment.InstallationAddress = req.InstallationAddress
	equipment.InstallationDate = req.InstallationDate
	equipment.ContractID = req.ContractID
	equipment.PurchaseDate = req.PurchaseDate
	equipment.PurchasePrice = req.PurchasePrice
	equipment.AMCContractID = req.AMCContractID
	equipment.Notes = req.Notes

	if req.Specifications != nil {
		equipment.Specifications = req.Specifications
	}

	// Calculate warranty expiry
	if req.WarrantyMonths > 0 {
		if req.PurchaseDate != nil {
			warrantyExpiry := req.PurchaseDate.AddDate(0, req.WarrantyMonths, 0)
			equipment.WarrantyExpiry = &warrantyExpiry
		} else {
			warrantyExpiry := time.Now().AddDate(0, req.WarrantyMonths, 0)
			equipment.WarrantyExpiry = &warrantyExpiry
		}
	}

	// Generate QR code URL
	equipment.QRCodeURL = fmt.Sprintf("%s/equipment/%s", s.baseURL, equipmentID)

	// Save to database
	if err := s.repo.Create(ctx, equipment); err != nil {
		s.logger.Error("Failed to register equipment", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to register equipment: %w", err)
	}

	s.logger.Info("Equipment registered successfully",
		slog.String("equipment_id", equipmentID),
		slog.String("serial_number", req.SerialNumber),
	)

	return equipment, nil
}

// GenerateQRCode generates QR code for equipment
func (s *EquipmentService) GenerateQRCode(ctx context.Context, equipmentID string) (string, error) {
	// Get equipment
	equipment, err := s.repo.GetByID(ctx, equipmentID)
	if err != nil {
		return "", fmt.Errorf("equipment not found: %w", err)
	}

	// Generate QR code image
	qrImagePath, err := s.qrGenerator.GenerateQRCode(equipment.ID, equipment.SerialNumber, equipment.QRCode)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	s.logger.Info("QR code generated",
		slog.String("equipment_id", equipmentID),
		slog.String("qr_path", qrImagePath),
	)

	return qrImagePath, nil
}

// GenerateQRLabel generates printable PDF label with QR code
func (s *EquipmentService) GenerateQRLabel(ctx context.Context, equipmentID string) (string, error) {
	// Get equipment
	equipment, err := s.repo.GetByID(ctx, equipmentID)
	if err != nil {
		return "", fmt.Errorf("equipment not found: %w", err)
	}

	// Generate QR code image first
	qrImagePath, err := s.qrGenerator.GenerateQRCode(equipment.ID, equipment.SerialNumber, equipment.QRCode)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Generate PDF label
	pdfPath, err := s.qrGenerator.GenerateQRLabel(
		equipment.ID,
		equipment.EquipmentName,
		equipment.SerialNumber,
		equipment.ManufacturerName,
		equipment.QRCode,
		qrImagePath,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF label: %w", err)
	}

	s.logger.Info("PDF label generated",
		slog.String("equipment_id", equipmentID),
		slog.String("pdf_path", pdfPath),
	)

	return pdfPath, nil
}

// GetEquipmentByID retrieves equipment by ID
func (s *EquipmentService) GetEquipmentByID(ctx context.Context, id string) (*domain.Equipment, error) {
	return s.repo.GetByID(ctx, id)
}

// GetEquipmentByQR retrieves equipment by QR code
func (s *EquipmentService) GetEquipmentByQR(ctx context.Context, qrCode string) (*domain.Equipment, error) {
	return s.repo.GetByQRCode(ctx, qrCode)
}

// GetEquipmentBySerial retrieves equipment by serial number
func (s *EquipmentService) GetEquipmentBySerial(ctx context.Context, serialNumber string) (*domain.Equipment, error) {
	return s.repo.GetBySerialNumber(ctx, serialNumber)
}

// ListEquipment retrieves equipment with filtering
func (s *EquipmentService) ListEquipment(ctx context.Context, criteria domain.ListCriteria) (*domain.ListResult, error) {
	return s.repo.List(ctx, criteria)
}

// UpdateEquipment updates equipment details
func (s *EquipmentService) UpdateEquipment(ctx context.Context, equipment *domain.Equipment) error {
	return s.repo.Update(ctx, equipment)
}

// RecordService records a service completion
func (s *EquipmentService) RecordService(ctx context.Context, equipmentID string, serviceDate time.Time, notes string) error {
	equipment, err := s.repo.GetByID(ctx, equipmentID)
	if err != nil {
		return fmt.Errorf("equipment not found: %w", err)
	}

	equipment.RecordService(serviceDate)
	equipment.Notes = notes

	if err := s.repo.Update(ctx, equipment); err != nil {
		return fmt.Errorf("failed to record service: %w", err)
	}

	s.logger.Info("Service recorded",
		slog.String("equipment_id", equipmentID),
		slog.Time("service_date", serviceDate),
	)

	return nil
}

// BulkImportFromCSV imports equipment from CSV file
func (s *EquipmentService) BulkImportFromCSV(ctx context.Context, csvFilePath, createdBy string) (*domain.CSVImportResult, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	s.logger.Info("CSV header", slog.Any("columns", header))

	result := &domain.CSVImportResult{
		Errors:      []string{},
		ImportedIDs: []string{},
	}

	equipmentList := []*domain.Equipment{}
	rowNum := 1

	// Read rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to read: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		// Parse row (expected columns match CSVImportRow structure)
		if len(row) < 13 {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Insufficient columns", rowNum))
			result.FailureCount++
			rowNum++
			continue
		}

		// Create equipment from CSV row
		equipment, err := s.parseCSVRow(row, createdBy)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %v", rowNum, err))
			result.FailureCount++
			rowNum++
			continue
		}

		equipmentList = append(equipmentList, equipment)
		result.SuccessCount++
		rowNum++
	}

	result.TotalRows = rowNum - 1

	// Bulk insert
	if len(equipmentList) > 0 {
		if err := s.repo.BulkCreate(ctx, equipmentList); err != nil {
			return nil, fmt.Errorf("failed to bulk insert: %w", err)
		}

		// Collect IDs
		for _, eq := range equipmentList {
			result.ImportedIDs = append(result.ImportedIDs, eq.ID)
		}
	}

	s.logger.Info("CSV import completed",
		slog.Int("total", result.TotalRows),
		slog.Int("success", result.SuccessCount),
		slog.Int("failure", result.FailureCount),
	)

	return result, nil
}

// parseCSVRow parses a CSV row into equipment
func (s *EquipmentService) parseCSVRow(row []string, createdBy string) (*domain.Equipment, error) {
	// CSV columns: serial_number, equipment_name, manufacturer_name, model_number, category,
	// customer_name, customer_id, installation_location, installation_date, purchase_date,
	// purchase_price, warranty_months, notes

	serialNumber := row[0]
	if serialNumber == "" {
		return nil, fmt.Errorf("serial_number is required")
	}

	equipmentName := row[1]
	if equipmentName == "" {
		return nil, fmt.Errorf("equipment_name is required")
	}

	manufacturerName := row[2]
	if manufacturerName == "" {
		return nil, fmt.Errorf("manufacturer_name is required")
	}

	modelNumber := row[3]
	category := row[4]
	customerName := row[5]
	if customerName == "" {
		return nil, fmt.Errorf("customer_name is required")
	}

	customerID := row[6]
	installationLocation := row[7]

	// Parse dates
	var installationDate, purchaseDate *time.Time
	if row[8] != "" {
		date, err := time.Parse("2006-01-02", row[8])
		if err == nil {
			installationDate = &date
		}
	}

	if row[9] != "" {
		date, err := time.Parse("2006-01-02", row[9])
		if err == nil {
			purchaseDate = &date
		}
	}

	// Parse price
	purchasePrice := 0.0
	if row[10] != "" {
		price, err := strconv.ParseFloat(row[10], 64)
		if err == nil {
			purchasePrice = price
		}
	}

	// Parse warranty months
	warrantyMonths := 0
	if row[11] != "" {
		months, err := strconv.Atoi(row[11])
		if err == nil {
			warrantyMonths = months
		}
	}

	notes := row[12]

	// Generate IDs
	equipmentID := ksuid.New().String()
	qrCodeID := s.generateQRCodeID()

	// Create equipment
	equipment := domain.NewEquipment(
		serialNumber,
		equipmentName,
		manufacturerName,
		modelNumber,
		customerName,
		createdBy,
	)

	equipment.ID = equipmentID
	equipment.QRCode = qrCodeID
	equipment.Category = category
	equipment.CustomerID = customerID
	equipment.InstallationLocation = installationLocation
	equipment.InstallationDate = installationDate
	equipment.PurchaseDate = purchaseDate
	equipment.PurchasePrice = purchasePrice
	equipment.Notes = notes

	// Calculate warranty expiry
	if warrantyMonths > 0 {
		if purchaseDate != nil {
			warrantyExpiry := purchaseDate.AddDate(0, warrantyMonths, 0)
			equipment.WarrantyExpiry = &warrantyExpiry
		} else {
			warrantyExpiry := time.Now().AddDate(0, warrantyMonths, 0)
			equipment.WarrantyExpiry = &warrantyExpiry
		}
	}

	// Generate QR code URL
	equipment.QRCodeURL = fmt.Sprintf("%s/equipment/%s", s.baseURL, equipmentID)

	return equipment, nil
}

// generateQRCodeID generates a unique QR code identifier
func (s *EquipmentService) generateQRCodeID() string {
	now := time.Now()
	// Format: QR-YYYYMMDD-XXXXXX (random 6 digits)
	return fmt.Sprintf("QR-%s-%06d", now.Format("20060102"), now.UnixNano()%1000000)
}
