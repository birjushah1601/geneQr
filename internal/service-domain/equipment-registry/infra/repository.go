package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EquipmentRepository implements the domain.Repository interface
type EquipmentRepository struct {
	pool *pgxpool.Pool
}

// NewEquipmentRepository creates a new equipment repository
func NewEquipmentRepository(pool *pgxpool.Pool) *EquipmentRepository {
	return &EquipmentRepository{pool: pool}
}

// Create creates a new equipment registration
func (r *EquipmentRepository) Create(ctx context.Context, equipment *domain.Equipment) error {
	query := `
		INSERT INTO equipment_registry (
			id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)
	`

	// Marshal JSONB fields
	specs, err := json.Marshal(equipment.Specifications)
	if err != nil {
		return fmt.Errorf("failed to marshal specifications: %w", err)
	}

	photos, err := json.Marshal(equipment.Photos)
	if err != nil {
		return fmt.Errorf("failed to marshal photos: %w", err)
	}

	docs, err := json.Marshal(equipment.Documents)
	if err != nil {
		return fmt.Errorf("failed to marshal documents: %w", err)
	}

	address, err := json.Marshal(equipment.InstallationAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal installation address: %w", err)
	}

	_, err = r.pool.Exec(ctx, query,
		equipment.ID,
		equipment.QRCode,
		equipment.SerialNumber,
		equipment.EquipmentID,
		equipment.EquipmentName,
		equipment.ManufacturerName,
		equipment.ModelNumber,
		equipment.Category,
		equipment.CustomerID,
		equipment.CustomerName,
		equipment.InstallationLocation,
		address,
		equipment.InstallationDate,
		equipment.ContractID,
		equipment.PurchaseDate,
		equipment.PurchasePrice,
		equipment.WarrantyExpiry,
		equipment.AMCContractID,
		equipment.Status,
		equipment.LastServiceDate,
		equipment.NextServiceDate,
		equipment.ServiceCount,
		specs,
		photos,
		docs,
		equipment.QRCodeURL,
		equipment.Notes,
		equipment.CreatedAt,
		equipment.UpdatedAt,
		equipment.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}

	return nil
}

// GetByID retrieves equipment by ID
func (r *EquipmentRepository) GetByID(ctx context.Context, id string) (*domain.Equipment, error) {
	query := `
		SELECT id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		FROM equipment_registry
		WHERE id = $1
	`

	equipment, err := r.scanEquipment(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("failed to get equipment by ID: %w", err)
	}

	return equipment, nil
}

// GetByQRCode retrieves equipment by QR code
func (r *EquipmentRepository) GetByQRCode(ctx context.Context, qrCode string) (*domain.Equipment, error) {
	query := `
		SELECT id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		FROM equipment_registry
		WHERE qr_code = $1
	`

	equipment, err := r.scanEquipment(r.pool.QueryRow(ctx, query, qrCode))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("failed to get equipment by QR code: %w", err)
	}

	return equipment, nil
}

// GetBySerialNumber retrieves equipment by serial number
func (r *EquipmentRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Equipment, error) {
	query := `
		SELECT id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		FROM equipment_registry
		WHERE serial_number = $1
	`

	equipment, err := r.scanEquipment(r.pool.QueryRow(ctx, query, serialNumber))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrEquipmentNotFound
		}
		return nil, fmt.Errorf("failed to get equipment by serial number: %w", err)
	}

	return equipment, nil
}

// List retrieves equipment with filtering
func (r *EquipmentRepository) List(ctx context.Context, criteria domain.ListCriteria) (*domain.ListResult, error) {
	// Build query with filters
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		FROM equipment_registry
		WHERE 1=1
	`)

	args := []interface{}{}
	argCount := 1

	// Apply filters
	if criteria.CustomerID != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND customer_id = $%d", argCount))
		args = append(args, criteria.CustomerID)
		argCount++
	}

	if criteria.ManufacturerName != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND manufacturer_name ILIKE $%d", argCount))
		args = append(args, "%"+criteria.ManufacturerName+"%")
		argCount++
	}

	if len(criteria.Status) > 0 {
		statuses := make([]string, len(criteria.Status))
		for i, status := range criteria.Status {
			statuses[i] = string(status)
		}
		queryBuilder.WriteString(fmt.Sprintf(" AND status = ANY($%d)", argCount))
		args = append(args, statuses)
		argCount++
	}

	if criteria.Category != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND category ILIKE $%d", argCount))
		args = append(args, "%"+criteria.Category+"%")
		argCount++
	}

	if criteria.HasAMC != nil {
		if *criteria.HasAMC {
			queryBuilder.WriteString(" AND amc_contract_id IS NOT NULL AND amc_contract_id != ''")
		} else {
			queryBuilder.WriteString(" AND (amc_contract_id IS NULL OR amc_contract_id = '')")
		}
	}

	if criteria.UnderWarranty != nil {
		if *criteria.UnderWarranty {
			queryBuilder.WriteString(" AND warranty_expiry IS NOT NULL AND warranty_expiry > NOW()")
		} else {
			queryBuilder.WriteString(" AND (warranty_expiry IS NULL OR warranty_expiry <= NOW())")
		}
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM equipment_registry WHERE 1=1"
	// Add same filters for count
	countBuilder := strings.Builder{}
	countBuilder.WriteString(countQuery)
	
	// Re-apply filters for count (simplified)
	if criteria.CustomerID != "" {
		countBuilder.WriteString(" AND customer_id = $1")
	}

	var total int
	err := r.pool.QueryRow(ctx, countBuilder.String(), args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count equipment: %w", err)
	}

	// Add sorting
	sortBy := "created_at"
	if criteria.SortBy != "" {
		sortBy = criteria.SortBy
	}

	sortDirection := "DESC"
	if criteria.SortDirection != "" {
		sortDirection = strings.ToUpper(criteria.SortDirection)
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDirection))

	// Add pagination
	page := criteria.Page
	if page < 1 {
		page = 1
	}

	pageSize := criteria.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1))
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list equipment: %w", err)
	}
	defer rows.Close()

	equipment := []*domain.Equipment{}
	for rows.Next() {
		eq, err := r.scanEquipmentFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan equipment: %w", err)
		}
		equipment = append(equipment, eq)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &domain.ListResult{
		Equipment:  equipment,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates equipment
func (r *EquipmentRepository) Update(ctx context.Context, equipment *domain.Equipment) error {
	query := `
		UPDATE equipment_registry SET
			qr_code = $2, serial_number = $3, equipment_id = $4, equipment_name = $5,
			manufacturer_name = $6, model_number = $7, category = $8, customer_id = $9,
			customer_name = $10, installation_location = $11, installation_address = $12,
			installation_date = $13, contract_id = $14, purchase_date = $15, purchase_price = $16,
			warranty_expiry = $17, amc_contract_id = $18, status = $19, last_service_date = $20,
			next_service_date = $21, service_count = $22, specifications = $23, photos = $24,
			documents = $25, qr_code_url = $26, notes = $27, updated_at = $28
		WHERE id = $1
	`

	// Marshal JSONB fields
	specs, _ := json.Marshal(equipment.Specifications)
	photos, _ := json.Marshal(equipment.Photos)
	docs, _ := json.Marshal(equipment.Documents)
	address, _ := json.Marshal(equipment.InstallationAddress)

	equipment.UpdatedAt = time.Now()

	result, err := r.pool.Exec(ctx, query,
		equipment.ID,
		equipment.QRCode,
		equipment.SerialNumber,
		equipment.EquipmentID,
		equipment.EquipmentName,
		equipment.ManufacturerName,
		equipment.ModelNumber,
		equipment.Category,
		equipment.CustomerID,
		equipment.CustomerName,
		equipment.InstallationLocation,
		address,
		equipment.InstallationDate,
		equipment.ContractID,
		equipment.PurchaseDate,
		equipment.PurchasePrice,
		equipment.WarrantyExpiry,
		equipment.AMCContractID,
		equipment.Status,
		equipment.LastServiceDate,
		equipment.NextServiceDate,
		equipment.ServiceCount,
		specs,
		photos,
		docs,
		equipment.QRCodeURL,
		equipment.Notes,
		equipment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update equipment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEquipmentNotFound
	}

	return nil
}

// Delete deletes equipment
func (r *EquipmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM equipment_registry WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrEquipmentNotFound
	}

	return nil
}

// BulkCreate creates multiple equipment registrations
func (r *EquipmentRepository) BulkCreate(ctx context.Context, equipmentList []*domain.Equipment) error {
	// Use transaction for bulk insert
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO equipment_registry (
			id, qr_code, serial_number, equipment_id, equipment_name, manufacturer_name,
			model_number, category, customer_id, customer_name, installation_location,
			installation_address, installation_date, contract_id, purchase_date, purchase_price,
			warranty_expiry, amc_contract_id, status, last_service_date, next_service_date,
			service_count, specifications, photos, documents, qr_code_url, notes,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		)
	`

	for _, equipment := range equipmentList {
		specs, _ := json.Marshal(equipment.Specifications)
		photos, _ := json.Marshal(equipment.Photos)
		docs, _ := json.Marshal(equipment.Documents)
		address, _ := json.Marshal(equipment.InstallationAddress)

		_, err = tx.Exec(ctx, query,
			equipment.ID,
			equipment.QRCode,
			equipment.SerialNumber,
			equipment.EquipmentID,
			equipment.EquipmentName,
			equipment.ManufacturerName,
			equipment.ModelNumber,
			equipment.Category,
			equipment.CustomerID,
			equipment.CustomerName,
			equipment.InstallationLocation,
			address,
			equipment.InstallationDate,
			equipment.ContractID,
			equipment.PurchaseDate,
			equipment.PurchasePrice,
			equipment.WarrantyExpiry,
			equipment.AMCContractID,
			equipment.Status,
			equipment.LastServiceDate,
			equipment.NextServiceDate,
			equipment.ServiceCount,
			specs,
			photos,
			docs,
			equipment.QRCodeURL,
			equipment.Notes,
			equipment.CreatedAt,
			equipment.UpdatedAt,
			equipment.CreatedBy,
		)

		if err != nil {
			return fmt.Errorf("failed to bulk create equipment: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// scanEquipment scans a single equipment row
func (r *EquipmentRepository) scanEquipment(row pgx.Row) (*domain.Equipment, error) {
	var equipment domain.Equipment
	var specs, photos, docs, address []byte

	err := row.Scan(
		&equipment.ID,
		&equipment.QRCode,
		&equipment.SerialNumber,
		&equipment.EquipmentID,
		&equipment.EquipmentName,
		&equipment.ManufacturerName,
		&equipment.ModelNumber,
		&equipment.Category,
		&equipment.CustomerID,
		&equipment.CustomerName,
		&equipment.InstallationLocation,
		&address,
		&equipment.InstallationDate,
		&equipment.ContractID,
		&equipment.PurchaseDate,
		&equipment.PurchasePrice,
		&equipment.WarrantyExpiry,
		&equipment.AMCContractID,
		&equipment.Status,
		&equipment.LastServiceDate,
		&equipment.NextServiceDate,
		&equipment.ServiceCount,
		&specs,
		&photos,
		&docs,
		&equipment.QRCodeURL,
		&equipment.Notes,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
		&equipment.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	if len(specs) > 0 {
		json.Unmarshal(specs, &equipment.Specifications)
	}
	if len(photos) > 0 {
		json.Unmarshal(photos, &equipment.Photos)
	}
	if len(docs) > 0 {
		json.Unmarshal(docs, &equipment.Documents)
	}
	if len(address) > 0 {
		json.Unmarshal(address, &equipment.InstallationAddress)
	}

	return &equipment, nil
}

// scanEquipmentFromRows scans equipment from rows
func (r *EquipmentRepository) scanEquipmentFromRows(rows pgx.Rows) (*domain.Equipment, error) {
	var equipment domain.Equipment
	var specs, photos, docs, address []byte

	err := rows.Scan(
		&equipment.ID,
		&equipment.QRCode,
		&equipment.SerialNumber,
		&equipment.EquipmentID,
		&equipment.EquipmentName,
		&equipment.ManufacturerName,
		&equipment.ModelNumber,
		&equipment.Category,
		&equipment.CustomerID,
		&equipment.CustomerName,
		&equipment.InstallationLocation,
		&address,
		&equipment.InstallationDate,
		&equipment.ContractID,
		&equipment.PurchaseDate,
		&equipment.PurchasePrice,
		&equipment.WarrantyExpiry,
		&equipment.AMCContractID,
		&equipment.Status,
		&equipment.LastServiceDate,
		&equipment.NextServiceDate,
		&equipment.ServiceCount,
		&specs,
		&photos,
		&docs,
		&equipment.QRCodeURL,
		&equipment.Notes,
		&equipment.CreatedAt,
		&equipment.UpdatedAt,
		&equipment.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	if len(specs) > 0 {
		json.Unmarshal(specs, &equipment.Specifications)
	}
	if len(photos) > 0 {
		json.Unmarshal(photos, &equipment.Photos)
	}
	if len(docs) > 0 {
		json.Unmarshal(docs, &equipment.Documents)
	}
	if len(address) > 0 {
		json.Unmarshal(address, &equipment.InstallationAddress)
	}

	return &equipment, nil
}
