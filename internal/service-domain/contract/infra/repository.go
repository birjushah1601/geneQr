package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/contract/domain"
	"github.com/jackc/pgx/v5"
)

// ContractRepository implements the domain.Repository interface
type ContractRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewContractRepository creates a new contract repository
func NewContractRepository(db *PostgresDB, logger *slog.Logger) *ContractRepository {
	return &ContractRepository{
		db:     db,
		logger: logger.With(slog.String("component", "contract_repository")),
	}
}

// Create creates a new contract
func (r *ContractRepository) Create(ctx context.Context, contract *domain.Contract) error {
	// Marshal JSONB fields
	paymentScheduleJSON, err := json.Marshal(contract.PaymentSchedule)
	if err != nil {
		return fmt.Errorf("failed to marshal payment schedule: %w", err)
	}

	deliveryScheduleJSON, err := json.Marshal(contract.DeliverySchedule)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery schedule: %w", err)
	}

	itemsJSON, err := json.Marshal(contract.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	amendmentsJSON, err := json.Marshal(contract.Amendments)
	if err != nil {
		return fmt.Errorf("failed to marshal amendments: %w", err)
	}

	query := `
		INSERT INTO contracts (
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err = r.db.Pool().Exec(ctx, query,
		contract.ID, contract.TenantID, contract.ContractNumber, contract.RFQID, contract.QuoteID,
		contract.SupplierID, contract.SupplierName, contract.Status, contract.TotalAmount,
		contract.Currency, contract.TaxAmount, contract.StartDate, contract.EndDate, contract.SignedDate,
		contract.PaymentTerms, contract.DeliveryTerms, contract.WarrantyTerms, contract.TermsAndConditions,
		paymentScheduleJSON, deliveryScheduleJSON, itemsJSON, amendmentsJSON,
		contract.Notes, contract.CreatedBy, contract.CreatedAt, contract.UpdatedAt, contract.SignedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	r.logger.Info("Contract created", slog.String("id", contract.ID), slog.String("contract_number", contract.ContractNumber))
	return nil
}

// GetByID retrieves a contract by ID
func (r *ContractRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.Contract, error) {
	query := `
		SELECT
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		FROM contracts
		WHERE id = $1 AND tenant_id = $2
	`

	var contract domain.Contract
	var paymentScheduleJSON, deliveryScheduleJSON, itemsJSON, amendmentsJSON []byte

	err := r.db.Pool().QueryRow(ctx, query, id, tenantID).Scan(
		&contract.ID, &contract.TenantID, &contract.ContractNumber, &contract.RFQID, &contract.QuoteID,
		&contract.SupplierID, &contract.SupplierName, &contract.Status, &contract.TotalAmount,
		&contract.Currency, &contract.TaxAmount, &contract.StartDate, &contract.EndDate, &contract.SignedDate,
		&contract.PaymentTerms, &contract.DeliveryTerms, &contract.WarrantyTerms, &contract.TermsAndConditions,
		&paymentScheduleJSON, &deliveryScheduleJSON, &itemsJSON, &amendmentsJSON,
		&contract.Notes, &contract.CreatedBy, &contract.CreatedAt, &contract.UpdatedAt, &contract.SignedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Unmarshal JSONB fields
	if err := json.Unmarshal(paymentScheduleJSON, &contract.PaymentSchedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment schedule: %w", err)
	}
	if err := json.Unmarshal(deliveryScheduleJSON, &contract.DeliverySchedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delivery schedule: %w", err)
	}
	if err := json.Unmarshal(itemsJSON, &contract.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}
	if err := json.Unmarshal(amendmentsJSON, &contract.Amendments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal amendments: %w", err)
	}

	return &contract, nil
}

// GetByContractNumber retrieves a contract by contract number
func (r *ContractRepository) GetByContractNumber(ctx context.Context, tenantID, contractNumber string) (*domain.Contract, error) {
	query := `
		SELECT
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		FROM contracts
		WHERE contract_number = $1 AND tenant_id = $2
	`

	var contract domain.Contract
	var paymentScheduleJSON, deliveryScheduleJSON, itemsJSON, amendmentsJSON []byte

	err := r.db.Pool().QueryRow(ctx, query, contractNumber, tenantID).Scan(
		&contract.ID, &contract.TenantID, &contract.ContractNumber, &contract.RFQID, &contract.QuoteID,
		&contract.SupplierID, &contract.SupplierName, &contract.Status, &contract.TotalAmount,
		&contract.Currency, &contract.TaxAmount, &contract.StartDate, &contract.EndDate, &contract.SignedDate,
		&contract.PaymentTerms, &contract.DeliveryTerms, &contract.WarrantyTerms, &contract.TermsAndConditions,
		&paymentScheduleJSON, &deliveryScheduleJSON, &itemsJSON, &amendmentsJSON,
		&contract.Notes, &contract.CreatedBy, &contract.CreatedAt, &contract.UpdatedAt, &contract.SignedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrContractNotFound
		}
		return nil, fmt.Errorf("failed to get contract by number: %w", err)
	}

	// Unmarshal JSONB fields
	if err := json.Unmarshal(paymentScheduleJSON, &contract.PaymentSchedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment schedule: %w", err)
	}
	if err := json.Unmarshal(deliveryScheduleJSON, &contract.DeliverySchedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delivery schedule: %w", err)
	}
	if err := json.Unmarshal(itemsJSON, &contract.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}
	if err := json.Unmarshal(amendmentsJSON, &contract.Amendments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal amendments: %w", err)
	}

	return &contract, nil
}

// GetByRFQ retrieves all contracts for an RFQ
func (r *ContractRepository) GetByRFQ(ctx context.Context, tenantID, rfqID string) ([]*domain.Contract, error) {
	query := `
		SELECT
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		FROM contracts
		WHERE rfq_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool().Query(ctx, query, rfqID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts by RFQ: %w", err)
	}
	defer rows.Close()

	return r.scanContracts(rows)
}

// GetBySupplier retrieves all contracts for a supplier
func (r *ContractRepository) GetBySupplier(ctx context.Context, tenantID, supplierID string) ([]*domain.Contract, error) {
	query := `
		SELECT
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		FROM contracts
		WHERE supplier_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool().Query(ctx, query, supplierID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts by supplier: %w", err)
	}
	defer rows.Close()

	return r.scanContracts(rows)
}

// List retrieves contracts with filtering
func (r *ContractRepository) List(ctx context.Context, criteria domain.ListCriteria) (*domain.ListResult, error) {
	// Build WHERE clause
	whereClauses := []string{"tenant_id = $1"}
	args := []interface{}{criteria.TenantID}
	argPos := 2

	if criteria.RFQID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("rfq_id = $%d", argPos))
		args = append(args, criteria.RFQID)
		argPos++
	}

	if criteria.SupplierID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("supplier_id = $%d", argPos))
		args = append(args, criteria.SupplierID)
		argPos++
	}

	if len(criteria.Status) > 0 {
		placeholders := make([]string, len(criteria.Status))
		for i, status := range criteria.Status {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, status)
			argPos++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if criteria.CreatedBy != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("created_by = $%d", argPos))
		args = append(args, criteria.CreatedBy)
		argPos++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contracts WHERE %s", whereClause)
	var total int
	err := r.db.Pool().QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if criteria.SortBy != "" {
		direction := "DESC"
		if criteria.SortDirection == "asc" {
			direction = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", criteria.SortBy, direction)
	}

	// Pagination
	page := criteria.Page
	if page < 1 {
		page = 1
	}
	pageSize := criteria.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Query with pagination
	query := fmt.Sprintf(`
		SELECT
			id, tenant_id, contract_number, rfq_id, quote_id,
			supplier_id, supplier_name, status, total_amount,
			currency, tax_amount, start_date, end_date, signed_date,
			payment_terms, delivery_terms, warranty_terms, terms_and_conditions,
			payment_schedule, delivery_schedule, items, amendments,
			notes, created_by, created_at, updated_at, signed_by
		FROM contracts
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argPos, argPos+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Pool().Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}
	defer rows.Close()

	contracts, err := r.scanContracts(rows)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &domain.ListResult{
		Contracts:  contracts,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a contract
func (r *ContractRepository) Update(ctx context.Context, contract *domain.Contract) error {
	// Marshal JSONB fields
	paymentScheduleJSON, err := json.Marshal(contract.PaymentSchedule)
	if err != nil {
		return fmt.Errorf("failed to marshal payment schedule: %w", err)
	}

	deliveryScheduleJSON, err := json.Marshal(contract.DeliverySchedule)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery schedule: %w", err)
	}

	itemsJSON, err := json.Marshal(contract.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	amendmentsJSON, err := json.Marshal(contract.Amendments)
	if err != nil {
		return fmt.Errorf("failed to marshal amendments: %w", err)
	}

	query := `
		UPDATE contracts SET
			contract_number = $1, rfq_id = $2, quote_id = $3,
			supplier_id = $4, supplier_name = $5, status = $6,
			total_amount = $7, currency = $8, tax_amount = $9,
			start_date = $10, end_date = $11, signed_date = $12,
			payment_terms = $13, delivery_terms = $14, warranty_terms = $15,
			terms_and_conditions = $16, payment_schedule = $17,
			delivery_schedule = $18, items = $19, amendments = $20,
			notes = $21, updated_at = $22, signed_by = $23
		WHERE id = $24 AND tenant_id = $25
	`

	result, err := r.db.Pool().Exec(ctx, query,
		contract.ContractNumber, contract.RFQID, contract.QuoteID,
		contract.SupplierID, contract.SupplierName, contract.Status,
		contract.TotalAmount, contract.Currency, contract.TaxAmount,
		contract.StartDate, contract.EndDate, contract.SignedDate,
		contract.PaymentTerms, contract.DeliveryTerms, contract.WarrantyTerms,
		contract.TermsAndConditions, paymentScheduleJSON,
		deliveryScheduleJSON, itemsJSON, amendmentsJSON,
		contract.Notes, contract.UpdatedAt, contract.SignedBy,
		contract.ID, contract.TenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrContractNotFound
	}

	r.logger.Info("Contract updated", slog.String("id", contract.ID))
	return nil
}

// Delete deletes a contract
func (r *ContractRepository) Delete(ctx context.Context, tenantID, id string) error {
	query := "DELETE FROM contracts WHERE id = $1 AND tenant_id = $2"

	result, err := r.db.Pool().Exec(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrContractNotFound
	}

	r.logger.Info("Contract deleted", slog.String("id", id))
	return nil
}

// scanContracts is a helper to scan multiple contracts from rows
func (r *ContractRepository) scanContracts(rows pgx.Rows) ([]*domain.Contract, error) {
	var contracts []*domain.Contract

	for rows.Next() {
		var contract domain.Contract
		var paymentScheduleJSON, deliveryScheduleJSON, itemsJSON, amendmentsJSON []byte

		err := rows.Scan(
			&contract.ID, &contract.TenantID, &contract.ContractNumber, &contract.RFQID, &contract.QuoteID,
			&contract.SupplierID, &contract.SupplierName, &contract.Status, &contract.TotalAmount,
			&contract.Currency, &contract.TaxAmount, &contract.StartDate, &contract.EndDate, &contract.SignedDate,
			&contract.PaymentTerms, &contract.DeliveryTerms, &contract.WarrantyTerms, &contract.TermsAndConditions,
			&paymentScheduleJSON, &deliveryScheduleJSON, &itemsJSON, &amendmentsJSON,
			&contract.Notes, &contract.CreatedBy, &contract.CreatedAt, &contract.UpdatedAt, &contract.SignedBy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}

		// Unmarshal JSONB fields
		if err := json.Unmarshal(paymentScheduleJSON, &contract.PaymentSchedule); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payment schedule: %w", err)
		}
		if err := json.Unmarshal(deliveryScheduleJSON, &contract.DeliverySchedule); err != nil {
			return nil, fmt.Errorf("failed to unmarshal delivery schedule: %w", err)
		}
		if err := json.Unmarshal(itemsJSON, &contract.Items); err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}
		if err := json.Unmarshal(amendmentsJSON, &contract.Amendments); err != nil {
			return nil, fmt.Errorf("failed to unmarshal amendments: %w", err)
		}

		contracts = append(contracts, &contract)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contracts: %w", err)
	}

	return contracts, nil
}
