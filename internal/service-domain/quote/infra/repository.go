package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/quote/domain"
	"github.com/jackc/pgx/v5"
)

// QuoteRepository implements the domain.QuoteRepository interface
type QuoteRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewQuoteRepository creates a new quote repository
func NewQuoteRepository(db *PostgresDB, logger *slog.Logger) *QuoteRepository {
	return &QuoteRepository{
		db:     db,
		logger: logger.With(slog.String("component", "quote_repository")),
	}
}

// Create persists a new quote with its items
func (r *QuoteRepository) Create(ctx context.Context, quote *domain.Quote) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert quote
	quoteQuery := `
		INSERT INTO quotes (
			id, tenant_id, rfq_id, supplier_id, quote_number, status,
			total_amount, currency, valid_until, delivery_terms, payment_terms,
			warranty_terms, notes, revision_number, reviewed_at, reviewed_by,
			review_notes, rejection_reason, metadata, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	metadataJSON, _ := json.Marshal(quote.Metadata)

	_, err = tx.Exec(ctx, quoteQuery,
		quote.ID, quote.TenantID, quote.RFQID, quote.SupplierID, quote.QuoteNumber,
		string(quote.Status), quote.TotalAmount, quote.Currency, quote.ValidUntil,
		quote.DeliveryTerms, quote.PaymentTerms, quote.WarrantyTerms, quote.Notes,
		quote.RevisionNumber, quote.ReviewedAt, quote.ReviewedBy, quote.ReviewNotes,
		quote.RejectionReason, metadataJSON, quote.CreatedBy, quote.CreatedAt, quote.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to create quote", slog.String("error", err.Error()))
		return fmt.Errorf("failed to create quote: %w", err)
	}

	// Insert quote items
	for _, item := range quote.Items {
		itemQuery := `
			INSERT INTO quote_items (
				id, quote_id, rfq_item_id, equipment_id, equipment_name,
				quantity, unit_price, total_price, tax_rate, tax_amount,
				delivery_timeframe, manufacturer_name, model_number,
				specifications, compliance_certs, notes, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			)
		`

		_, err = tx.Exec(ctx, itemQuery,
			item.ID, quote.ID, item.RFQItemID, item.EquipmentID, item.EquipmentName,
			item.Quantity, item.UnitPrice, item.TotalPrice, item.TaxRate, item.TaxAmount,
			item.DeliveryTimeframe, item.ManufacturerName, item.ModelNumber,
			item.Specifications, item.ComplianceCerts, item.Notes,
		)
		if err != nil {
			r.logger.Error("Failed to create quote item", slog.String("error", err.Error()))
			return fmt.Errorf("failed to create quote item: %w", err)
		}
	}

	// Insert revisions if any
	for _, revision := range quote.Revisions {
		revisionQuery := `
			INSERT INTO quote_revisions (
				quote_id, revision_number, revised_at, revised_by, changes,
				previous_total, new_total, metadata
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		revisionMetadata, _ := json.Marshal(revision.Metadata)

		_, err = tx.Exec(ctx, revisionQuery,
			quote.ID, revision.RevisionNumber, revision.RevisedAt, revision.RevisedBy,
			revision.Changes, revision.PreviousTotal, revision.NewTotal, revisionMetadata,
		)
		if err != nil {
			r.logger.Error("Failed to create quote revision", slog.String("error", err.Error()))
			return fmt.Errorf("failed to create quote revision: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Quote created successfully", slog.String("quote_id", quote.ID))
	return nil
}

// GetByID retrieves a quote by ID with its items and revisions
func (r *QuoteRepository) GetByID(ctx context.Context, id string, tenantID string) (*domain.Quote, error) {
	// Get quote
	quoteQuery := `
		SELECT 
			id, tenant_id, rfq_id, supplier_id, quote_number, status,
			total_amount, currency, valid_until, delivery_terms, payment_terms,
			warranty_terms, notes, revision_number, reviewed_at, reviewed_by,
			review_notes, rejection_reason, metadata, created_by, created_at, updated_at
		FROM quotes
		WHERE id = $1 AND tenant_id = $2
	`

	quote := &domain.Quote{}
	var metadataJSON []byte
	var status string

	err := r.db.pool.QueryRow(ctx, quoteQuery, id, tenantID).Scan(
		&quote.ID, &quote.TenantID, &quote.RFQID, &quote.SupplierID, &quote.QuoteNumber,
		&status, &quote.TotalAmount, &quote.Currency, &quote.ValidUntil,
		&quote.DeliveryTerms, &quote.PaymentTerms, &quote.WarrantyTerms, &quote.Notes,
		&quote.RevisionNumber, &quote.ReviewedAt, &quote.ReviewedBy, &quote.ReviewNotes,
		&quote.RejectionReason, &metadataJSON, &quote.CreatedBy, &quote.CreatedAt, &quote.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrQuoteNotFound
		}
		r.logger.Error("Failed to get quote", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	quote.Status = domain.QuoteStatus(status)
	if err := json.Unmarshal(metadataJSON, &quote.Metadata); err != nil {
		quote.Metadata = make(map[string]interface{})
	}

	// Get quote items
	itemsQuery := `
		SELECT 
			id, rfq_item_id, equipment_id, equipment_name, quantity,
			unit_price, total_price, tax_rate, tax_amount, delivery_timeframe,
			manufacturer_name, model_number, specifications, compliance_certs, notes
		FROM quote_items
		WHERE quote_id = $1
		ORDER BY id
	`

	rows, err := r.db.pool.Query(ctx, itemsQuery, quote.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote items: %w", err)
	}
	defer rows.Close()

	quote.Items = []domain.QuoteItem{}
	for rows.Next() {
		var item domain.QuoteItem
		err := rows.Scan(
			&item.ID, &item.RFQItemID, &item.EquipmentID, &item.EquipmentName, &item.Quantity,
			&item.UnitPrice, &item.TotalPrice, &item.TaxRate, &item.TaxAmount, &item.DeliveryTimeframe,
			&item.ManufacturerName, &item.ModelNumber, &item.Specifications, &item.ComplianceCerts, &item.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote item: %w", err)
		}
		quote.Items = append(quote.Items, item)
	}

	// Get revisions
	revisionsQuery := `
		SELECT 
			revision_number, revised_at, revised_by, changes,
			previous_total, new_total, metadata
		FROM quote_revisions
		WHERE quote_id = $1
		ORDER BY revision_number
	`

	revRows, err := r.db.pool.Query(ctx, revisionsQuery, quote.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote revisions: %w", err)
	}
	defer revRows.Close()

	quote.Revisions = []domain.QuoteRevision{}
	for revRows.Next() {
		var revision domain.QuoteRevision
		var revMetadata []byte
		err := revRows.Scan(
			&revision.RevisionNumber, &revision.RevisedAt, &revision.RevisedBy, &revision.Changes,
			&revision.PreviousTotal, &revision.NewTotal, &revMetadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote revision: %w", err)
		}
		if err := json.Unmarshal(revMetadata, &revision.Metadata); err != nil {
			revision.Metadata = make(map[string]interface{})
		}
		quote.Revisions = append(quote.Revisions, revision)
	}

	return quote, nil
}

// GetByRFQID retrieves all quotes for an RFQ
func (r *QuoteRepository) GetByRFQID(ctx context.Context, rfqID string, tenantID string) ([]*domain.Quote, error) {
	query := `
		SELECT id
		FROM quotes
		WHERE rfq_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, rfqID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes by RFQ: %w", err)
	}
	defer rows.Close()

	quotes := []*domain.Quote{}
	for rows.Next() {
		var quoteID string
		if err := rows.Scan(&quoteID); err != nil {
			return nil, fmt.Errorf("failed to scan quote ID: %w", err)
		}

		quote, err := r.GetByID(ctx, quoteID, tenantID)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// GetBySupplierID retrieves all quotes from a supplier
func (r *QuoteRepository) GetBySupplierID(ctx context.Context, supplierID string, tenantID string) ([]*domain.Quote, error) {
	query := `
		SELECT id
		FROM quotes
		WHERE supplier_id = $1 AND tenant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, supplierID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes by supplier: %w", err)
	}
	defer rows.Close()

	quotes := []*domain.Quote{}
	for rows.Next() {
		var quoteID string
		if err := rows.Scan(&quoteID); err != nil {
			return nil, fmt.Errorf("failed to scan quote ID: %w", err)
		}

		quote, err := r.GetByID(ctx, quoteID, tenantID)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// List retrieves quotes with filtering and pagination
func (r *QuoteRepository) List(ctx context.Context, criteria domain.ListCriteria) ([]*domain.Quote, int, error) {
	// Build query dynamically
	queryParts := []string{`
		SELECT id
		FROM quotes
		WHERE tenant_id = $1
	`}

	params := []interface{}{criteria.TenantID}
	paramIndex := 2

	// Add RFQ filter
	if criteria.RFQID != "" {
		queryParts = append(queryParts, fmt.Sprintf("AND rfq_id = $%d", paramIndex))
		params = append(params, criteria.RFQID)
		paramIndex++
	}

	// Add supplier filter
	if criteria.SupplierID != "" {
		queryParts = append(queryParts, fmt.Sprintf("AND supplier_id = $%d", paramIndex))
		params = append(params, criteria.SupplierID)
		paramIndex++
	}

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

	// Add amount filters
	if criteria.MinAmount > 0 {
		queryParts = append(queryParts, fmt.Sprintf("AND total_amount >= $%d", paramIndex))
		params = append(params, criteria.MinAmount)
		paramIndex++
	}
	if criteria.MaxAmount > 0 {
		queryParts = append(queryParts, fmt.Sprintf("AND total_amount <= $%d", paramIndex))
		params = append(params, criteria.MaxAmount)
		paramIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM quotes WHERE tenant_id = $1"
	var total int
	err := r.db.pool.QueryRow(ctx, countQuery, criteria.TenantID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count quotes: %w", err)
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
		r.logger.Error("Failed to list quotes", slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	quotes := []*domain.Quote{}
	for rows.Next() {
		var quoteID string
		if err := rows.Scan(&quoteID); err != nil {
			return nil, 0, fmt.Errorf("failed to scan quote ID: %w", err)
		}

		quote, err := r.GetByID(ctx, quoteID, criteria.TenantID)
		if err != nil {
			return nil, 0, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, total, nil
}

// Update updates an existing quote
func (r *QuoteRepository) Update(ctx context.Context, quote *domain.Quote) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update quote
	quoteQuery := `
		UPDATE quotes SET
			quote_number = $1, status = $2, total_amount = $3, currency = $4,
			valid_until = $5, delivery_terms = $6, payment_terms = $7,
			warranty_terms = $8, notes = $9, revision_number = $10,
			reviewed_at = $11, reviewed_by = $12, review_notes = $13,
			rejection_reason = $14, metadata = $15, updated_at = $16
		WHERE id = $17 AND tenant_id = $18
	`

	metadataJSON, _ := json.Marshal(quote.Metadata)

	result, err := tx.Exec(ctx, quoteQuery,
		quote.QuoteNumber, string(quote.Status), quote.TotalAmount, quote.Currency,
		quote.ValidUntil, quote.DeliveryTerms, quote.PaymentTerms, quote.WarrantyTerms,
		quote.Notes, quote.RevisionNumber, quote.ReviewedAt, quote.ReviewedBy,
		quote.ReviewNotes, quote.RejectionReason, metadataJSON, quote.UpdatedAt,
		quote.ID, quote.TenantID,
	)
	if err != nil {
		r.logger.Error("Failed to update quote", slog.String("error", err.Error()))
		return fmt.Errorf("failed to update quote: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrQuoteNotFound
	}

	// Delete existing items
	_, err = tx.Exec(ctx, "DELETE FROM quote_items WHERE quote_id = $1", quote.ID)
	if err != nil {
		return fmt.Errorf("failed to delete quote items: %w", err)
	}

	// Re-insert items
	for _, item := range quote.Items {
		itemQuery := `
			INSERT INTO quote_items (
				id, quote_id, rfq_item_id, equipment_id, equipment_name,
				quantity, unit_price, total_price, tax_rate, tax_amount,
				delivery_timeframe, manufacturer_name, model_number,
				specifications, compliance_certs, notes, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			)
		`

		_, err = tx.Exec(ctx, itemQuery,
			item.ID, quote.ID, item.RFQItemID, item.EquipmentID, item.EquipmentName,
			item.Quantity, item.UnitPrice, item.TotalPrice, item.TaxRate, item.TaxAmount,
			item.DeliveryTimeframe, item.ManufacturerName, item.ModelNumber,
			item.Specifications, item.ComplianceCerts, item.Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to create quote item: %w", err)
		}
	}

	// Update revisions (append new ones)
	for _, revision := range quote.Revisions {
		// Check if revision exists
		var exists bool
		err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM quote_revisions WHERE quote_id = $1 AND revision_number = $2)",
			quote.ID, revision.RevisionNumber).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check revision existence: %w", err)
		}

		if !exists {
			revisionQuery := `
				INSERT INTO quote_revisions (
					quote_id, revision_number, revised_at, revised_by, changes,
					previous_total, new_total, metadata
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			`

			revisionMetadata, _ := json.Marshal(revision.Metadata)

			_, err = tx.Exec(ctx, revisionQuery,
				quote.ID, revision.RevisionNumber, revision.RevisedAt, revision.RevisedBy,
				revision.Changes, revision.PreviousTotal, revision.NewTotal, revisionMetadata,
			)
			if err != nil {
				return fmt.Errorf("failed to create quote revision: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.logger.Info("Quote updated successfully", slog.String("quote_id", quote.ID))
	return nil
}

// Delete removes a quote
func (r *QuoteRepository) Delete(ctx context.Context, id string, tenantID string) error {
	query := `DELETE FROM quotes WHERE id = $1 AND tenant_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		r.logger.Error("Failed to delete quote", slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrQuoteNotFound
	}

	r.logger.Info("Quote deleted successfully", slog.String("quote_id", id))
	return nil
}
