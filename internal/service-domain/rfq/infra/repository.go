package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/rfq/domain"
	"github.com/jackc/pgx/v5"
)

// RFQRepository implements the domain.RFQRepository interface
type RFQRepository struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewRFQRepository creates a new RFQ repository
func NewRFQRepository(db *PostgresDB, logger *slog.Logger) *RFQRepository {
	return &RFQRepository{
		db:     db,
		logger: logger.With(slog.String("component", "rfq_repository")),
	}
}

// Create persists a new RFQ
func (r *RFQRepository) Create(ctx context.Context, rfq *domain.RFQ) error {
	query := `
		INSERT INTO rfqs (
			id, rfq_number, tenant_id, title, description, priority, status,
			delivery_terms, payment_terms, response_deadline,
			created_by, created_at, updated_at, internal_notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	// Convert terms to JSON
	deliveryTermsJSON, err := json.Marshal(rfq.DeliveryTerms)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery terms: %w", err)
	}

	paymentTermsJSON, err := json.Marshal(rfq.PaymentTerms)
	if err != nil {
		return fmt.Errorf("failed to marshal payment terms: %w", err)
	}

	_, err = r.db.pool.Exec(
		ctx,
		query,
		rfq.ID,
		rfq.RFQNumber,
		rfq.TenantID,
		rfq.Title,
		rfq.Description,
		string(rfq.Priority),
		string(rfq.Status),
		deliveryTermsJSON,
		paymentTermsJSON,
		rfq.ResponseDeadline,
		rfq.CreatedBy,
		rfq.CreatedAt,
		rfq.UpdatedAt,
		rfq.InternalNotes,
	)

	if err != nil {
		r.logger.Error("Failed to create RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfq.ID))
		return fmt.Errorf("failed to create RFQ: %w", err)
	}

	r.logger.Info("RFQ created successfully",
		slog.String("rfq_id", rfq.ID),
		slog.String("rfq_number", rfq.RFQNumber))

	return nil
}

// GetByID retrieves an RFQ by ID
func (r *RFQRepository) GetByID(ctx context.Context, id string, tenantID string) (*domain.RFQ, error) {
	query := `
		SELECT 
			id, rfq_number, tenant_id, title, description, priority, status,
			delivery_terms, payment_terms, published_at, response_deadline, closed_at,
			created_by, created_at, updated_at, internal_notes
		FROM rfqs
		WHERE id = $1 AND tenant_id = $2
	`

	row := r.db.pool.QueryRow(ctx, query, id, tenantID)

	var rfq domain.RFQ
	var deliveryTermsJSON, paymentTermsJSON []byte
	var priority, status string

	err := row.Scan(
		&rfq.ID,
		&rfq.RFQNumber,
		&rfq.TenantID,
		&rfq.Title,
		&rfq.Description,
		&priority,
		&status,
		&deliveryTermsJSON,
		&paymentTermsJSON,
		&rfq.PublishedAt,
		&rfq.ResponseDeadline,
		&rfq.ClosedAt,
		&rfq.CreatedBy,
		&rfq.CreatedAt,
		&rfq.UpdatedAt,
		&rfq.InternalNotes,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRFQNotFound
		}
		r.logger.Error("Failed to get RFQ by ID",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return nil, fmt.Errorf("failed to get RFQ: %w", err)
	}

	rfq.Priority = domain.RFQPriority(priority)
	rfq.Status = domain.RFQStatus(status)

	// Unmarshal terms
	if err := json.Unmarshal(deliveryTermsJSON, &rfq.DeliveryTerms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delivery terms: %w", err)
	}

	if err := json.Unmarshal(paymentTermsJSON, &rfq.PaymentTerms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment terms: %w", err)
	}

	// Load items
	items, err := r.GetItems(ctx, rfq.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load RFQ items: %w", err)
	}
	rfq.Items = items

	// Load invitations
	invitations, err := r.GetInvitations(ctx, rfq.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load RFQ invitations: %w", err)
	}
	rfq.Invitations = invitations

	return &rfq, nil
}

// GetByRFQNumber retrieves an RFQ by its RFQ number
func (r *RFQRepository) GetByRFQNumber(ctx context.Context, rfqNumber string, tenantID string) (*domain.RFQ, error) {
	query := `
		SELECT id FROM rfqs
		WHERE rfq_number = $1 AND tenant_id = $2
	`

	var id string
	err := r.db.pool.QueryRow(ctx, query, rfqNumber, tenantID).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRFQNotFound
		}
		return nil, fmt.Errorf("failed to get RFQ by number: %w", err)
	}

	return r.GetByID(ctx, id, tenantID)
}

// Update updates an existing RFQ
func (r *RFQRepository) Update(ctx context.Context, rfq *domain.RFQ) error {
	query := `
		UPDATE rfqs SET
			title = $1,
			description = $2,
			priority = $3,
			status = $4,
			delivery_terms = $5,
			payment_terms = $6,
			published_at = $7,
			response_deadline = $8,
			closed_at = $9,
			updated_at = $10,
			internal_notes = $11
		WHERE id = $12 AND tenant_id = $13
	`

	deliveryTermsJSON, err := json.Marshal(rfq.DeliveryTerms)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery terms: %w", err)
	}

	paymentTermsJSON, err := json.Marshal(rfq.PaymentTerms)
	if err != nil {
		return fmt.Errorf("failed to marshal payment terms: %w", err)
	}

	result, err := r.db.pool.Exec(
		ctx,
		query,
		rfq.Title,
		rfq.Description,
		string(rfq.Priority),
		string(rfq.Status),
		deliveryTermsJSON,
		paymentTermsJSON,
		rfq.PublishedAt,
		rfq.ResponseDeadline,
		rfq.ClosedAt,
		rfq.UpdatedAt,
		rfq.InternalNotes,
		rfq.ID,
		rfq.TenantID,
	)

	if err != nil {
		r.logger.Error("Failed to update RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfq.ID))
		return fmt.Errorf("failed to update RFQ: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRFQNotFound
	}

	return nil
}

// Delete removes an RFQ
func (r *RFQRepository) Delete(ctx context.Context, id string, tenantID string) error {
	query := `DELETE FROM rfqs WHERE id = $1 AND tenant_id = $2`

	result, err := r.db.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		r.logger.Error("Failed to delete RFQ",
			slog.String("error", err.Error()),
			slog.String("rfq_id", id))
		return fmt.Errorf("failed to delete RFQ: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRFQNotFound
	}

	return nil
}

// List retrieves RFQs with pagination and filtering
func (r *RFQRepository) List(ctx context.Context, criteria domain.ListCriteria) ([]*domain.RFQ, int, error) {
	// Build the query dynamically
	queryParts := []string{`
		SELECT 
			id, rfq_number, tenant_id, title, description, priority, status,
			delivery_terms, payment_terms, published_at, response_deadline, closed_at,
			created_by, created_at, updated_at, internal_notes
		FROM rfqs
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

	// Add priority filter
	if len(criteria.Priority) > 0 {
		placeholders := make([]string, len(criteria.Priority))
		for i, priority := range criteria.Priority {
			placeholders[i] = fmt.Sprintf("$%d", paramIndex)
			params = append(params, string(priority))
			paramIndex++
		}
		queryParts = append(queryParts, fmt.Sprintf("AND priority IN (%s)", strings.Join(placeholders, ", ")))
	}

	// Add search query
	if criteria.SearchQuery != "" {
		queryParts = append(queryParts, fmt.Sprintf(`AND (
			title ILIKE $%d OR description ILIKE $%d
		)`, paramIndex, paramIndex))
		params = append(params, "%"+criteria.SearchQuery+"%")
		paramIndex++
	}

	// Add created_by filter
	if criteria.CreatedBy != "" {
		queryParts = append(queryParts, fmt.Sprintf("AND created_by = $%d", paramIndex))
		params = append(params, criteria.CreatedBy)
		paramIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM rfqs WHERE tenant_id = $1"
	var total int
	err := r.db.pool.QueryRow(ctx, countQuery, criteria.TenantID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count RFQs: %w", err)
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
		r.logger.Error("Failed to list RFQs",
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to list RFQs: %w", err)
	}
	defer rows.Close()

	rfqs := []*domain.RFQ{}
	for rows.Next() {
		var rfq domain.RFQ
		var deliveryTermsJSON, paymentTermsJSON []byte
		var priority, status string

		err := rows.Scan(
			&rfq.ID,
			&rfq.RFQNumber,
			&rfq.TenantID,
			&rfq.Title,
			&rfq.Description,
			&priority,
			&status,
			&deliveryTermsJSON,
			&paymentTermsJSON,
			&rfq.PublishedAt,
			&rfq.ResponseDeadline,
			&rfq.ClosedAt,
			&rfq.CreatedBy,
			&rfq.CreatedAt,
			&rfq.UpdatedAt,
			&rfq.InternalNotes,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan RFQ: %w", err)
		}

		rfq.Priority = domain.RFQPriority(priority)
		rfq.Status = domain.RFQStatus(status)

		// Unmarshal terms
		if err := json.Unmarshal(deliveryTermsJSON, &rfq.DeliveryTerms); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal delivery terms: %w", err)
		}

		if err := json.Unmarshal(paymentTermsJSON, &rfq.PaymentTerms); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal payment terms: %w", err)
		}

		// Load items count (for list view, we don't need all details)
		items, _ := r.GetItems(ctx, rfq.ID)
		rfq.Items = items

		rfqs = append(rfqs, &rfq)
	}

	return rfqs, total, nil
}

// AddItem adds an item to an RFQ
func (r *RFQRepository) AddItem(ctx context.Context, rfqID string, item *domain.RFQItem) error {
	query := `
		INSERT INTO rfq_items (
			id, rfq_id, equipment_id, category_id, name, description,
			specifications, quantity, unit, estimated_price, notes,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	specsJSON, err := json.Marshal(item.Specifications)
	if err != nil {
		return fmt.Errorf("failed to marshal specifications: %w", err)
	}

	_, err = r.db.pool.Exec(
		ctx,
		query,
		item.ID,
		rfqID,
		item.EquipmentID,
		item.CategoryID,
		item.Name,
		item.Description,
		specsJSON,
		item.Quantity,
		item.Unit,
		item.EstimatedPrice,
		item.Notes,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to add RFQ item",
			slog.String("error", err.Error()),
			slog.String("rfq_id", rfqID))
		return fmt.Errorf("failed to add RFQ item: %w", err)
	}

	return nil
}

// UpdateItem updates an RFQ item
func (r *RFQRepository) UpdateItem(ctx context.Context, item *domain.RFQItem) error {
	query := `
		UPDATE rfq_items SET
			name = $1,
			description = $2,
			specifications = $3,
			quantity = $4,
			unit = $5,
			estimated_price = $6,
			notes = $7,
			updated_at = $8
		WHERE id = $9
	`

	specsJSON, err := json.Marshal(item.Specifications)
	if err != nil {
		return fmt.Errorf("failed to marshal specifications: %w", err)
	}

	result, err := r.db.pool.Exec(
		ctx,
		query,
		item.Name,
		item.Description,
		specsJSON,
		item.Quantity,
		item.Unit,
		item.EstimatedPrice,
		item.Notes,
		item.UpdatedAt,
		item.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update RFQ item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("RFQ item not found")
	}

	return nil
}

// RemoveItem removes an item from an RFQ
func (r *RFQRepository) RemoveItem(ctx context.Context, rfqID, itemID string) error {
	query := `DELETE FROM rfq_items WHERE id = $1 AND rfq_id = $2`

	result, err := r.db.pool.Exec(ctx, query, itemID, rfqID)
	if err != nil {
		return fmt.Errorf("failed to remove RFQ item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("RFQ item not found")
	}

	return nil
}

// GetItems retrieves all items for an RFQ
func (r *RFQRepository) GetItems(ctx context.Context, rfqID string) ([]domain.RFQItem, error) {
	query := `
		SELECT 
			id, rfq_id, equipment_id, category_id, name, description,
			specifications, quantity, unit, estimated_price, notes,
			created_at, updated_at
		FROM rfq_items
		WHERE rfq_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.pool.Query(ctx, query, rfqID)
	if err != nil {
		return nil, fmt.Errorf("failed to get RFQ items: %w", err)
	}
	defer rows.Close()

	items := []domain.RFQItem{}
	for rows.Next() {
		var item domain.RFQItem
		var specsJSON []byte

		err := rows.Scan(
			&item.ID,
			&item.RFQID,
			&item.EquipmentID,
			&item.CategoryID,
			&item.Name,
			&item.Description,
			&specsJSON,
			&item.Quantity,
			&item.Unit,
			&item.EstimatedPrice,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan RFQ item: %w", err)
		}

		if len(specsJSON) > 0 {
			if err := json.Unmarshal(specsJSON, &item.Specifications); err != nil {
				return nil, fmt.Errorf("failed to unmarshal specifications: %w", err)
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// AddInvitation adds a supplier invitation
func (r *RFQRepository) AddInvitation(ctx context.Context, invitation *domain.RFQInvitation) error {
	query := `
		INSERT INTO rfq_invitations (
			id, rfq_id, supplier_id, status, invited_at, message
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`

	_, err := r.db.pool.Exec(
		ctx,
		query,
		invitation.ID,
		invitation.RFQID,
		invitation.SupplierID,
		invitation.Status,
		invitation.InvitedAt,
		invitation.Message,
	)

	if err != nil {
		r.logger.Error("Failed to add RFQ invitation",
			slog.String("error", err.Error()),
			slog.String("rfq_id", invitation.RFQID))
		return fmt.Errorf("failed to add RFQ invitation: %w", err)
	}

	return nil
}

// GetInvitations retrieves all invitations for an RFQ
func (r *RFQRepository) GetInvitations(ctx context.Context, rfqID string) ([]domain.RFQInvitation, error) {
	query := `
		SELECT 
			id, rfq_id, supplier_id, status, invited_at, viewed_at, responded_at, message
		FROM rfq_invitations
		WHERE rfq_id = $1
		ORDER BY invited_at DESC
	`

	rows, err := r.db.pool.Query(ctx, query, rfqID)
	if err != nil {
		return nil, fmt.Errorf("failed to get RFQ invitations: %w", err)
	}
	defer rows.Close()

	invitations := []domain.RFQInvitation{}
	for rows.Next() {
		var inv domain.RFQInvitation

		err := rows.Scan(
			&inv.ID,
			&inv.RFQID,
			&inv.SupplierID,
			&inv.Status,
			&inv.InvitedAt,
			&inv.ViewedAt,
			&inv.RespondedAt,
			&inv.Message,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan RFQ invitation: %w", err)
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

// UpdateInvitation updates an invitation status
func (r *RFQRepository) UpdateInvitation(ctx context.Context, invitation *domain.RFQInvitation) error {
	query := `
		UPDATE rfq_invitations SET
			status = $1,
			viewed_at = $2,
			responded_at = $3
		WHERE id = $4
	`

	result, err := r.db.pool.Exec(
		ctx,
		query,
		invitation.Status,
		invitation.ViewedAt,
		invitation.RespondedAt,
		invitation.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update RFQ invitation: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("RFQ invitation not found")
	}

	return nil
}
