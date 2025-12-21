package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
)

// TicketRepository implements the domain.TicketRepository interface
type TicketRepository struct {
	pool *pgxpool.Pool
}

// NewTicketRepository creates a new ticket repository
func NewTicketRepository(pool *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{pool: pool}
}

// UpdateTicketParts updates the parts assigned to a ticket in ticket_parts table
func (r *TicketRepository) UpdateTicketParts(ctx context.Context, ticketID string, parts []map[string]interface{}) error {
	// Clear existing parts for this ticket
	deleteQuery := `DELETE FROM ticket_parts WHERE ticket_id = $1`
	if _, err := r.pool.Exec(ctx, deleteQuery, ticketID); err != nil {
		return fmt.Errorf("failed to clear existing parts: %w", err)
	}

	// Insert new parts into ticket_parts table
	insertQuery := `
		INSERT INTO ticket_parts (
			ticket_id, 
			spare_part_id, 
			quantity_required, 
			is_critical,
			unit_price,
			total_price,
			currency,
			assigned_by,
			status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, part := range parts {
		partID, ok := part["part_id"].(string)
		if !ok || partID == "" {
			continue // Skip invalid parts
		}

		quantity := 1
		if q, ok := part["quantity"].(float64); ok {
			quantity = int(q)
		}

		isCritical := false
		if c, ok := part["is_critical"].(bool); ok {
			isCritical = c
		}

		var unitPrice, totalPrice float64
		if p, ok := part["unit_price"].(float64); ok {
			unitPrice = p
			totalPrice = unitPrice * float64(quantity)
		}

		currency := "USD"
		if c, ok := part["currency"].(string); ok && c != "" {
			currency = c
		}

		assignedBy := "system"
		if a, ok := part["assigned_by"].(string); ok && a != "" {
			assignedBy = a
		}

		status := "pending"
		if st, ok := part["status"].(string); ok && st != "" {
			status = st
		}

		_, err := r.pool.Exec(ctx, insertQuery,
			ticketID, partID, quantity, isCritical,
			unitPrice, totalPrice, currency, assignedBy, status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert part %s: %w", partID, err)
		}
	}

	return nil
}

// Create creates a new service ticket
func (r *TicketRepository) Create(ctx context.Context, ticket *domain.ServiceTicket) error {
	if ticket.ID == "" {
		ticket.ID = ksuid.New().String()
	}

	// Marshal JSONB fields
	partsUsed, _ := json.Marshal(ticket.PartsUsed)
	photos, _ := json.Marshal(ticket.Photos)
	videos, _ := json.Marshal(ticket.Videos)
	documents, _ := json.Marshal(ticket.Documents)

	query := `
		INSERT INTO service_tickets (
			id, ticket_number, equipment_id, qr_code, serial_number, equipment_name,
			customer_id, customer_name, customer_phone, customer_whatsapp,
			issue_category, issue_description, priority, severity,
			source, source_message_id,
			assigned_engineer_id, assigned_engineer_name, assigned_at,
			status, created_at, acknowledged_at, started_at, resolved_at, closed_at,
			sla_response_due, sla_resolution_due, sla_breached,
			resolution_notes, parts_used, labor_hours, cost,
			photos, videos, documents,
			amc_contract_id, covered_under_amc,
			updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
		)
	`

	_, err := r.pool.Exec(ctx, query,
		ticket.ID, ticket.TicketNumber, ticket.EquipmentID, ticket.QRCode, ticket.SerialNumber, ticket.EquipmentName,
		ticket.CustomerID, ticket.CustomerName, ticket.CustomerPhone, ticket.CustomerWhatsApp,
		ticket.IssueCategory, ticket.IssueDescription, ticket.Priority, ticket.Severity,
		ticket.Source, ticket.SourceMessageID,
		ticket.AssignedEngineerID, ticket.AssignedEngineerName, ticket.AssignedAt,
		ticket.Status, ticket.CreatedAt, ticket.AcknowledgedAt, ticket.StartedAt, ticket.ResolvedAt, ticket.ClosedAt,
		ticket.SLAResponseDue, ticket.SLAResolutionDue, ticket.SLABreached,
		ticket.ResolutionNotes, partsUsed, ticket.LaborHours, ticket.Cost,
		photos, videos, documents,
		ticket.AMCContractID, ticket.CoveredUnderAMC,
		ticket.UpdatedAt, ticket.CreatedBy,
	)

	return err
}

// UpdateResponsibility sets responsible_org_id and policy_provenance for a ticket (Phase 4 optional)
func (r *TicketRepository) UpdateResponsibility(ctx context.Context, ticketID string, responsibleOrgID *string, provenance json.RawMessage) error {
    query := `UPDATE service_tickets SET responsible_org_id = $2, policy_provenance = COALESCE($3, '{}'::jsonb) WHERE id = $1`
    _, err := r.pool.Exec(ctx, query, ticketID, responsibleOrgID, provenance)
    return err
}

// GetByID retrieves a ticket by ID
func (r *TicketRepository) GetByID(ctx context.Context, id string) (*domain.ServiceTicket, error) {
	query := `
		SELECT 
			id, ticket_number, equipment_id, qr_code, serial_number, equipment_name,
			customer_id, customer_name, customer_phone, customer_whatsapp,
			issue_category, issue_description, priority, severity,
			source, source_message_id,
			assigned_engineer_id, assigned_engineer_name, assigned_at,
			status, created_at, acknowledged_at, started_at, resolved_at, closed_at,
			sla_response_due, sla_resolution_due, sla_breached,
			resolution_notes, parts_used, labor_hours, cost,
			photos, videos, documents,
			amc_contract_id, covered_under_amc,
			updated_at, created_by
		FROM service_tickets
		WHERE id = $1
	`

	var ticket domain.ServiceTicket
	var partsUsed, photos, videos, documents []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.EquipmentID, &ticket.QRCode, &ticket.SerialNumber, &ticket.EquipmentName,
		&ticket.CustomerID, &ticket.CustomerName, &ticket.CustomerPhone, &ticket.CustomerWhatsApp,
		&ticket.IssueCategory, &ticket.IssueDescription, &ticket.Priority, &ticket.Severity,
		&ticket.Source, &ticket.SourceMessageID,
		&ticket.AssignedEngineerID, &ticket.AssignedEngineerName, &ticket.AssignedAt,
		&ticket.Status, &ticket.CreatedAt, &ticket.AcknowledgedAt, &ticket.StartedAt, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.SLAResponseDue, &ticket.SLAResolutionDue, &ticket.SLABreached,
		&ticket.ResolutionNotes, &partsUsed, &ticket.LaborHours, &ticket.Cost,
		&photos, &videos, &documents,
		&ticket.AMCContractID, &ticket.CoveredUnderAMC,
		&ticket.UpdatedAt, &ticket.CreatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}

	// Unmarshal JSONB fields
	json.Unmarshal(partsUsed, &ticket.PartsUsed)
	json.Unmarshal(photos, &ticket.Photos)
	json.Unmarshal(videos, &ticket.Videos)
	json.Unmarshal(documents, &ticket.Documents)

	return &ticket, nil
}

// GetByTicketNumber retrieves a ticket by ticket number
func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*domain.ServiceTicket, error) {
	query := `
		SELECT 
			id, ticket_number, equipment_id, qr_code, serial_number, equipment_name,
			customer_id, customer_name, customer_phone, customer_whatsapp,
			issue_category, issue_description, priority, severity,
			source, source_message_id,
			assigned_engineer_id, assigned_engineer_name, assigned_at,
			status, created_at, acknowledged_at, started_at, resolved_at, closed_at,
			sla_response_due, sla_resolution_due, sla_breached,
			resolution_notes, parts_used, labor_hours, cost,
			photos, videos, documents,
			amc_contract_id, covered_under_amc,
			updated_at, created_by
		FROM service_tickets
		WHERE ticket_number = $1
	`

	var ticket domain.ServiceTicket
	var partsUsed, photos, videos, documents []byte

	err := r.pool.QueryRow(ctx, query, ticketNumber).Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.EquipmentID, &ticket.QRCode, &ticket.SerialNumber, &ticket.EquipmentName,
		&ticket.CustomerID, &ticket.CustomerName, &ticket.CustomerPhone, &ticket.CustomerWhatsApp,
		&ticket.IssueCategory, &ticket.IssueDescription, &ticket.Priority, &ticket.Severity,
		&ticket.Source, &ticket.SourceMessageID,
		&ticket.AssignedEngineerID, &ticket.AssignedEngineerName, &ticket.AssignedAt,
		&ticket.Status, &ticket.CreatedAt, &ticket.AcknowledgedAt, &ticket.StartedAt, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.SLAResponseDue, &ticket.SLAResolutionDue, &ticket.SLABreached,
		&ticket.ResolutionNotes, &partsUsed, &ticket.LaborHours, &ticket.Cost,
		&photos, &videos, &documents,
		&ticket.AMCContractID, &ticket.CoveredUnderAMC,
		&ticket.UpdatedAt, &ticket.CreatedBy,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrTicketNotFound
		}
		return nil, err
	}

	// Unmarshal JSONB fields
	json.Unmarshal(partsUsed, &ticket.PartsUsed)
	json.Unmarshal(photos, &ticket.Photos)
	json.Unmarshal(videos, &ticket.Videos)
	json.Unmarshal(documents, &ticket.Documents)

	return &ticket, nil
}

// Update updates an existing ticket
func (r *TicketRepository) Update(ctx context.Context, ticket *domain.ServiceTicket) error {
	// Marshal JSONB fields
	partsUsed, _ := json.Marshal(ticket.PartsUsed)
	photos, _ := json.Marshal(ticket.Photos)
	videos, _ := json.Marshal(ticket.Videos)
	documents, _ := json.Marshal(ticket.Documents)

	query := `
		UPDATE service_tickets SET
			ticket_number = $2, equipment_id = $3, qr_code = $4, serial_number = $5, equipment_name = $6,
			customer_id = $7, customer_name = $8, customer_phone = $9, customer_whatsapp = $10,
			issue_category = $11, issue_description = $12, priority = $13, severity = $14,
			source = $15, source_message_id = $16,
			assigned_engineer_id = $17, assigned_engineer_name = $18, assigned_at = $19,
			status = $20, acknowledged_at = $21, started_at = $22, resolved_at = $23, closed_at = $24,
			sla_response_due = $25, sla_resolution_due = $26, sla_breached = $27,
			resolution_notes = $28, parts_used = $29, labor_hours = $30, cost = $31,
			photos = $32, videos = $33, documents = $34,
			amc_contract_id = $35, covered_under_amc = $36
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		ticket.ID, ticket.TicketNumber, ticket.EquipmentID, ticket.QRCode, ticket.SerialNumber, ticket.EquipmentName,
		ticket.CustomerID, ticket.CustomerName, ticket.CustomerPhone, ticket.CustomerWhatsApp,
		ticket.IssueCategory, ticket.IssueDescription, ticket.Priority, ticket.Severity,
		ticket.Source, ticket.SourceMessageID,
		ticket.AssignedEngineerID, ticket.AssignedEngineerName, ticket.AssignedAt,
		ticket.Status, ticket.AcknowledgedAt, ticket.StartedAt, ticket.ResolvedAt, ticket.ClosedAt,
		ticket.SLAResponseDue, ticket.SLAResolutionDue, ticket.SLABreached,
		ticket.ResolutionNotes, partsUsed, ticket.LaborHours, ticket.Cost,
		photos, videos, documents,
		ticket.AMCContractID, ticket.CoveredUnderAMC,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTicketNotFound
	}

	return nil
}

// List retrieves tickets based on criteria
func (r *TicketRepository) List(ctx context.Context, criteria domain.ListCriteria) (*domain.TicketListResult, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argPos := 1

	if len(criteria.Status) > 0 {
		placeholders := make([]string, len(criteria.Status))
		for i, status := range criteria.Status {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, status)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(criteria.Priority) > 0 {
		placeholders := make([]string, len(criteria.Priority))
		for i, priority := range criteria.Priority {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, priority)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("priority IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(criteria.Source) > 0 {
		placeholders := make([]string, len(criteria.Source))
		for i, source := range criteria.Source {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, source)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("source IN (%s)", strings.Join(placeholders, ",")))
	}

	if criteria.EquipmentID != "" {
		conditions = append(conditions, fmt.Sprintf("equipment_id = $%d", argPos))
		args = append(args, criteria.EquipmentID)
		argPos++
	}

	if criteria.CustomerID != "" {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", argPos))
		args = append(args, criteria.CustomerID)
		argPos++
	}

	if criteria.EngineerID != "" {
		conditions = append(conditions, fmt.Sprintf("assigned_engineer_id = $%d", argPos))
		args = append(args, criteria.EngineerID)
		argPos++
	}

	if criteria.SLABreached != nil {
		conditions = append(conditions, fmt.Sprintf("sla_breached = $%d", argPos))
		args = append(args, *criteria.SLABreached)
		argPos++
	}

	if criteria.CoveredUnderAMC != nil {
		conditions = append(conditions, fmt.Sprintf("covered_under_amc = $%d", argPos))
		args = append(args, *criteria.CoveredUnderAMC)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM service_tickets %s", whereClause)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Sorting
	sortBy := "created_at"
	if criteria.SortBy != "" {
		sortBy = criteria.SortBy
	}
	sortDir := "DESC"
	if criteria.SortDirection != "" {
		sortDir = strings.ToUpper(criteria.SortDirection)
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
			id, ticket_number, equipment_id, qr_code, serial_number, equipment_name,
			customer_id, customer_name, customer_phone, customer_whatsapp,
			issue_category, issue_description, priority, severity,
			source, source_message_id,
			assigned_engineer_id, assigned_engineer_name, assigned_at,
			status, created_at, acknowledged_at, started_at, resolved_at, closed_at,
			sla_response_due, sla_resolution_due, sla_breached,
			resolution_notes, parts_used, labor_hours, cost,
			photos, videos, documents,
			amc_contract_id, covered_under_amc,
			updated_at, created_by
		FROM service_tickets
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortBy, sortDir, argPos, argPos+1)

	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*domain.ServiceTicket
	for rows.Next() {
		var ticket domain.ServiceTicket
		var partsUsed, photos, videos, documents []byte

		err := rows.Scan(
			&ticket.ID, &ticket.TicketNumber, &ticket.EquipmentID, &ticket.QRCode, &ticket.SerialNumber, &ticket.EquipmentName,
			&ticket.CustomerID, &ticket.CustomerName, &ticket.CustomerPhone, &ticket.CustomerWhatsApp,
			&ticket.IssueCategory, &ticket.IssueDescription, &ticket.Priority, &ticket.Severity,
			&ticket.Source, &ticket.SourceMessageID,
			&ticket.AssignedEngineerID, &ticket.AssignedEngineerName, &ticket.AssignedAt,
			&ticket.Status, &ticket.CreatedAt, &ticket.AcknowledgedAt, &ticket.StartedAt, &ticket.ResolvedAt, &ticket.ClosedAt,
			&ticket.SLAResponseDue, &ticket.SLAResolutionDue, &ticket.SLABreached,
			&ticket.ResolutionNotes, &partsUsed, &ticket.LaborHours, &ticket.Cost,
			&photos, &videos, &documents,
			&ticket.AMCContractID, &ticket.CoveredUnderAMC,
			&ticket.UpdatedAt, &ticket.CreatedBy,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSONB fields
		json.Unmarshal(partsUsed, &ticket.PartsUsed)
		json.Unmarshal(photos, &ticket.Photos)
		json.Unmarshal(videos, &ticket.Videos)
		json.Unmarshal(documents, &ticket.Documents)

		tickets = append(tickets, &ticket)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.TicketListResult{
		Tickets:    tickets,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByEquipment retrieves all tickets for an equipment
func (r *TicketRepository) GetByEquipment(ctx context.Context, equipmentID string) ([]*domain.ServiceTicket, error) {
	criteria := domain.ListCriteria{
		EquipmentID: equipmentID,
		PageSize:    1000,
	}
	result, err := r.List(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

// GetByCustomer retrieves all tickets for a customer
func (r *TicketRepository) GetByCustomer(ctx context.Context, customerID string) ([]*domain.ServiceTicket, error) {
	criteria := domain.ListCriteria{
		CustomerID: customerID,
		PageSize:   1000,
	}
	result, err := r.List(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

// GetByEngineer retrieves all tickets assigned to an engineer
func (r *TicketRepository) GetByEngineer(ctx context.Context, engineerID string) ([]*domain.ServiceTicket, error) {
	criteria := domain.ListCriteria{
		EngineerID: engineerID,
		PageSize:   1000,
	}
	result, err := r.List(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

// GetBySource retrieves tickets by source
func (r *TicketRepository) GetBySource(ctx context.Context, source domain.TicketSource) ([]*domain.ServiceTicket, error) {
	criteria := domain.ListCriteria{
		Source:   []domain.TicketSource{source},
		PageSize: 1000,
	}
	result, err := r.List(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return result.Tickets, nil
}

// AddComment adds a comment to a ticket
func (r *TicketRepository) AddComment(ctx context.Context, comment *domain.TicketComment) error {
	if comment.ID == "" {
		comment.ID = ksuid.New().String()
	}

	attachments, _ := json.Marshal(comment.Attachments)

	query := `
		INSERT INTO ticket_comments (id, ticket_id, comment_type, author_id, author_name, comment, attachments)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.TicketID, comment.CommentType, comment.AuthorID, comment.AuthorName, comment.Comment, attachments,
	)

	return err
}

// GetComments retrieves all comments for a ticket
func (r *TicketRepository) GetComments(ctx context.Context, ticketID string) ([]*domain.TicketComment, error) {
	query := `
		SELECT id, ticket_id, comment_type, author_id, author_name, comment, attachments, created_at::text
		FROM ticket_comments
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.TicketComment
	for rows.Next() {
		var comment domain.TicketComment
		var attachments []byte

		err := rows.Scan(
			&comment.ID, &comment.TicketID, &comment.CommentType, &comment.AuthorID, &comment.AuthorName,
			&comment.Comment, &attachments, &comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(attachments, &comment.Attachments)
		comments = append(comments, &comment)
	}

	return comments, nil
}

// AddStatusHistory records a status change
func (r *TicketRepository) AddStatusHistory(ctx context.Context, history *domain.StatusHistory) error {
	if history.ID == "" {
		history.ID = ksuid.New().String()
	}

	query := `
		INSERT INTO ticket_status_history (id, ticket_id, from_status, to_status, changed_by, reason)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		history.ID, history.TicketID, history.FromStatus, history.ToStatus, history.ChangedBy, history.Reason,
	)

	return err
}

// GetStatusHistory retrieves status history for a ticket
func (r *TicketRepository) GetStatusHistory(ctx context.Context, ticketID string) ([]*domain.StatusHistory, error) {
	query := `
		SELECT id, ticket_id, from_status, to_status, changed_by, changed_at, reason
		FROM ticket_status_history
		WHERE ticket_id = $1
		ORDER BY changed_at ASC
	`

	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*domain.StatusHistory
	for rows.Next() {
		var h domain.StatusHistory
		err := rows.Scan(&h.ID, &h.TicketID, &h.FromStatus, &h.ToStatus, &h.ChangedBy, &h.ChangedAt, &h.Reason)
		if err != nil {
			return nil, err
		}
		history = append(history, &h)
	}

	return history, nil
}

// DeleteComment removes a comment from a ticket
func (r *TicketRepository) DeleteComment(ctx context.Context, commentID string, ticketID string) error {
	query := `
		DELETE FROM ticket_comments 
		WHERE id = $1 AND ticket_id = $2
	`
	
	result, err := r.pool.Exec(ctx, query, commentID, ticketID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or does not belong to this ticket")
	}
	
	return nil
}
