package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/supplier/domain"
	"github.com/jackc/pgx/v5"
)

// SupplierRepository implements the domain.SupplierRepository interface
type SupplierRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewSupplierRepository creates a new supplier repository
func NewSupplierRepository(db *PostgresDB, logger *slog.Logger) *SupplierRepository {
	return &SupplierRepository{
		db:     db,
		logger: logger.With(slog.String("component", "supplier_repository")),
	}
}

// Create persists a new supplier
func (r *SupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		INSERT INTO suppliers (
			id, tenant_id, company_name, business_registration_number, tax_id,
			year_established, description, contact_info, address, specializations,
			certifications, performance_rating, total_orders, completed_orders,
			status, verification_status, verified_at, verified_by, metadata,
			created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	// Convert structs to JSON
	contactInfoJSON, err := json.Marshal(supplier.ContactInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal contact info: %w", err)
	}

	addressJSON, err := json.Marshal(supplier.Address)
	if err != nil {
		return fmt.Errorf("failed to marshal address: %w", err)
	}

	certificationsJSON, err := json.Marshal(supplier.Certifications)
	if err != nil {
		return fmt.Errorf("failed to marshal certifications: %w", err)
	}

	metadataJSON, err := json.Marshal(supplier.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Ensure specializations array is not nil
	specializations := supplier.Specializations
	if specializations == nil {
		specializations = []string{}
	}

	_, err = r.db.pool.Exec(
		ctx,
		query,
		supplier.ID,
		supplier.TenantID,
		supplier.CompanyName,
		supplier.BusinessRegistrationNum,
		supplier.TaxID,
		supplier.YearEstablished,
		supplier.Description,
		contactInfoJSON,
		addressJSON,
		specializations, // pgx v5 handles slices natively
		certificationsJSON,
		supplier.PerformanceRating,
		supplier.TotalOrders,
		supplier.CompletedOrders,
		string(supplier.Status),
		string(supplier.VerificationStatus),
		supplier.VerifiedAt,
		supplier.VerifiedBy,
		metadataJSON,
		supplier.CreatedBy,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplier.ID))
		return fmt.Errorf("failed to create supplier: %w", err)
	}

	r.logger.Info("Supplier created successfully",
		slog.String("supplier_id", supplier.ID),
		slog.String("company_name", supplier.CompanyName))

	return nil
}

// GetByID retrieves a supplier by ID
func (r *SupplierRepository) GetByID(ctx context.Context, id string, tenantID string) (*domain.Supplier, error) {
	query := `
		SELECT 
			id, tenant_id, company_name, business_registration_number, tax_id,
			year_established, description, contact_info, address, specializations,
			certifications, performance_rating, total_orders, completed_orders,
			status, verification_status, verified_at, verified_by, metadata,
			created_by, created_at, updated_at
		FROM suppliers
		WHERE id = $1 AND tenant_id = $2
	`

	row := r.db.pool.QueryRow(ctx, query, id, tenantID)
	return r.scanSupplier(row)
}

// GetByTaxID retrieves a supplier by Tax ID
func (r *SupplierRepository) GetByTaxID(ctx context.Context, taxID string, tenantID string) (*domain.Supplier, error) {
	query := `
		SELECT 
			id, tenant_id, company_name, business_registration_number, tax_id,
			year_established, description, contact_info, address, specializations,
			certifications, performance_rating, total_orders, completed_orders,
			status, verification_status, verified_at, verified_by, metadata,
			created_by, created_at, updated_at
		FROM suppliers
		WHERE tax_id = $1 AND tenant_id = $2
	`

	row := r.db.pool.QueryRow(ctx, query, taxID, tenantID)

	return r.scanSupplier(row)
}

// Update updates an existing supplier
func (r *SupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	query := `
		UPDATE suppliers SET
			company_name = $1,
			business_registration_number = $2,
			tax_id = $3,
			year_established = $4,
			description = $5,
			contact_info = $6,
			address = $7,
			specializations = $8,
			certifications = $9,
			performance_rating = $10,
			total_orders = $11,
			completed_orders = $12,
			status = $13,
			verification_status = $14,
			verified_at = $15,
			verified_by = $16,
			metadata = $17,
			updated_at = $18
		WHERE id = $19 AND tenant_id = $20
	`

	contactInfoJSON, _ := json.Marshal(supplier.ContactInfo)
	addressJSON, _ := json.Marshal(supplier.Address)
	certificationsJSON, _ := json.Marshal(supplier.Certifications)
	metadataJSON, _ := json.Marshal(supplier.Metadata)

	result, err := r.db.pool.Exec(
		ctx,
		query,
		supplier.CompanyName,
		supplier.BusinessRegistrationNum,
		supplier.TaxID,
		supplier.YearEstablished,
		supplier.Description,
		contactInfoJSON,
		addressJSON,
		supplier.Specializations, // pgx v5 handles slices natively
		certificationsJSON,
		supplier.PerformanceRating,
		supplier.TotalOrders,
		supplier.CompletedOrders,
		string(supplier.Status),
		string(supplier.VerificationStatus),
		supplier.VerifiedAt,
		supplier.VerifiedBy,
		metadataJSON,
		supplier.UpdatedAt,
		supplier.ID,
		supplier.TenantID,
	)

	if err != nil {
		r.logger.Error("Failed to update supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplier.ID))
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSupplierNotFound
	}

	return nil
}

// Delete removes a supplier
func (r *SupplierRepository) Delete(ctx context.Context, id string, tenantID string) error {
	query := `DELETE FROM suppliers WHERE id = $1 AND tenant_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		r.logger.Error("Failed to delete supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", id))
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSupplierNotFound
	}

	return nil
}

// List retrieves suppliers with filtering and pagination
func (r *SupplierRepository) List(ctx context.Context, criteria domain.ListCriteria) ([]*domain.Supplier, int, error) {
	// Build the query dynamically
	queryParts := []string{`
		SELECT 
			id, tenant_id, company_name, business_registration_number, tax_id,
			year_established, description, contact_info, address, specializations,
			certifications, performance_rating, total_orders, completed_orders,
			status, verification_status, verified_at, verified_by, metadata,
			created_by, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1
	`}

	params := []interface{}{criteria.TenantID}
	paramIndex := 2

	// Add status filter
	if len(criteria.Status) > 0 {
		placeholders := make([]string, len(criteria.Status))
		for i, status := range criteria.Status {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, string(status))
			paramIndex++
		}
		queryParts = append(queryParts, fmt.Sprintf("AND status IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add verification status filter
	if len(criteria.VerificationStatus) > 0 {
		placeholders := make([]string, len(criteria.VerificationStatus))
		for i, status := range criteria.VerificationStatus {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, string(status))
			paramIndex++
		}
		queryParts = append(queryParts, fmt.Sprintf("AND verification_status IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add category filter
	if criteria.CategoryID != "" {
		queryParts = append(queryParts, fmt.Sprintf("AND $%d = ANY(specializations)", paramIndex))
		params = append(params, criteria.CategoryID)
		paramIndex++
	}

	// Add search query
	if criteria.SearchQuery != "" {
		queryParts = append(queryParts, fmt.Sprintf("AND company_name ILIKE $%d", paramIndex))
		params = append(params, "%"+criteria.SearchQuery+"%")
		paramIndex++
	}

	// Add minimum rating filter
	if criteria.MinRating > 0 {
		queryParts = append(queryParts, fmt.Sprintf("AND performance_rating >= $%d", paramIndex))
		params = append(params, criteria.MinRating)
		paramIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM suppliers WHERE tenant_id = $1"
	var total int
	err := r.db.pool.QueryRow(ctx, countQuery, criteria.TenantID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	// Add sorting
	sortBy := "created_at"
	if criteria.SortBy != "" {
		sortBy = criteria.SortBy
	}
	sortDirection := "DESC"
	if criteria.SortDirection == "asc" {
		sortDirection = "ASC"
	}
	queryParts = append(queryParts, fmt.Sprintf("ORDER BY %s %s", sortBy, sortDirection))

	// Add pagination
	pageSize := 20
	if criteria.PageSize > 0 {
		pageSize = criteria.PageSize
	}
	page := 1
	if criteria.Page > 0 {
		page = criteria.Page
	}
	offset := (page - 1) * pageSize

	queryParts = append(queryParts, fmt.Sprintf("LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1))
	params = append(params, pageSize, offset)

	// Execute query
	query := strings.Join(queryParts, " ")
	rows, err := r.db.pool.Query(ctx, query, params...)
	if err != nil {
		r.logger.Error("Failed to list suppliers",
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	suppliers := []*domain.Supplier{}
	for rows.Next() {
		supplier, err := r.scanSupplierFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, total, nil
}

// GetByCategory retrieves suppliers specialized in a category
func (r *SupplierRepository) GetByCategory(ctx context.Context, categoryID string, tenantID string) ([]*domain.Supplier, error) {
	query := `
		SELECT 
			id, tenant_id, company_name, business_registration_number, tax_id,
			year_established, description, contact_info, address, specializations,
			certifications, performance_rating, total_orders, completed_orders,
			status, verification_status, verified_at, verified_by, metadata,
			created_by, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1
		  AND status = 'active'
		  AND verification_status = 'approved'
		  AND $2 = ANY(specializations)
		ORDER BY performance_rating DESC
	`

	rows, err := r.db.pool.Query(ctx, query, tenantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get suppliers by category: %w", err)
	}
	defer rows.Close()

	suppliers := []*domain.Supplier{}
	for rows.Next() {
		supplier, err := r.scanSupplierFromRows(rows)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil
}

// Helper methods

func (r *SupplierRepository) scanSupplier(row pgx.Row) (*domain.Supplier, error) {
	var supplier domain.Supplier
	var contactInfoJSON, addressJSON, certificationsJSON, metadataJSON []byte
	var specializations []string
	var status, verificationStatus string

	err := row.Scan(
		&supplier.ID,
		&supplier.TenantID,
		&supplier.CompanyName,
		&supplier.BusinessRegistrationNum,
		&supplier.TaxID,
		&supplier.YearEstablished,
		&supplier.Description,
		&contactInfoJSON,
		&addressJSON,
		&specializations, // pgx v5 scans arrays directly
		&certificationsJSON,
		&supplier.PerformanceRating,
		&supplier.TotalOrders,
		&supplier.CompletedOrders,
		&status,
		&verificationStatus,
		&supplier.VerifiedAt,
		&supplier.VerifiedBy,
		&metadataJSON,
		&supplier.CreatedBy,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSupplierNotFound
		}
		r.logger.Error("Failed to scan supplier",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to scan supplier: %w", err)
	}

	supplier.Status = domain.SupplierStatus(status)
	supplier.VerificationStatus = domain.VerificationStatus(verificationStatus)
	supplier.Specializations = specializations

	// Unmarshal JSON fields
	if err := json.Unmarshal(contactInfoJSON, &supplier.ContactInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contact info: %w", err)
	}

	if err := json.Unmarshal(addressJSON, &supplier.Address); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address: %w", err)
	}

	if err := json.Unmarshal(certificationsJSON, &supplier.Certifications); err != nil {
		return nil, fmt.Errorf("failed to unmarshal certifications: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &supplier.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &supplier, nil
}

func (r *SupplierRepository) scanSupplierFromRows(rows pgx.Rows) (*domain.Supplier, error) {
	var supplier domain.Supplier
	var contactInfoJSON, addressJSON, certificationsJSON, metadataJSON []byte
	var specializations []string
	var status, verificationStatus string

	err := rows.Scan(
		&supplier.ID,
		&supplier.TenantID,
		&supplier.CompanyName,
		&supplier.BusinessRegistrationNum,
		&supplier.TaxID,
		&supplier.YearEstablished,
		&supplier.Description,
		&contactInfoJSON,
		&addressJSON,
		&specializations, // pgx v5 scans arrays directly
		&certificationsJSON,
		&supplier.PerformanceRating,
		&supplier.TotalOrders,
		&supplier.CompletedOrders,
		&status,
		&verificationStatus,
		&supplier.VerifiedAt,
		&supplier.VerifiedBy,
		&metadataJSON,
		&supplier.CreatedBy,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan supplier: %w", err)
	}

	supplier.Status = domain.SupplierStatus(status)
	supplier.VerificationStatus = domain.VerificationStatus(verificationStatus)
	supplier.Specializations = specializations

	// Unmarshal JSON fields
	json.Unmarshal(contactInfoJSON, &supplier.ContactInfo)
	json.Unmarshal(addressJSON, &supplier.Address)
	json.Unmarshal(certificationsJSON, &supplier.Certifications)
	json.Unmarshal(metadataJSON, &supplier.Metadata)

	return &supplier, nil
}
