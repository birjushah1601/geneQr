package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/ksuid"
)

// AssignmentRepository implements domain.AssignmentRepository
type AssignmentRepository struct {
	pool *pgxpool.Pool
}

// NewAssignmentRepository creates a new assignment repository
func NewAssignmentRepository(pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{pool: pool}
}

// Create creates a new engineer assignment
func (r *AssignmentRepository) Create(ctx context.Context, assignment *domain.EngineerAssignment) error {
	if assignment.ID == "" {
		assignment.ID = ksuid.New().String()
	}

	// Marshal JSONB fields
	partsUsed, _ := json.Marshal(assignment.PartsUsed)

	query := `
		INSERT INTO engineer_assignments (
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27, $28
		)
	`

	_, err := r.pool.Exec(ctx, query,
		assignment.ID, assignment.TicketID, assignment.EngineerID, assignment.EquipmentID,
		assignment.AssignmentSequence, assignment.AssignmentTier, assignment.AssignmentTierName, assignment.AssignmentReason,
		assignment.AssignmentType, assignment.Status, assignment.AssignedBy, assignment.AssignedAt,
		assignment.AcceptedAt, assignment.RejectedAt, assignment.RejectionReason,
		assignment.StartedAt, assignment.CompletedAt, assignment.CompletionStatus, assignment.EscalationReason,
		assignment.TimeSpentHours, assignment.Diagnosis, assignment.ActionsTaken, partsUsed,
		assignment.CustomerRating, assignment.CustomerFeedback, assignment.Notes,
		assignment.CreatedAt, assignment.UpdatedAt,
	)

	return err
}

// GetByID retrieves an assignment by ID
func (r *AssignmentRepository) GetByID(ctx context.Context, id string) (*domain.EngineerAssignment, error) {
	query := `
		SELECT 
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM engineer_assignments
		WHERE id = $1
	`

	assignment := &domain.EngineerAssignment{}
	var partsUsedJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&assignment.ID, &assignment.TicketID, &assignment.EngineerID, &assignment.EquipmentID,
		&assignment.AssignmentSequence, &assignment.AssignmentTier, &assignment.AssignmentTierName, &assignment.AssignmentReason,
		&assignment.AssignmentType, &assignment.Status, &assignment.AssignedBy, &assignment.AssignedAt,
		&assignment.AcceptedAt, &assignment.RejectedAt, &assignment.RejectionReason,
		&assignment.StartedAt, &assignment.CompletedAt, &assignment.CompletionStatus, &assignment.EscalationReason,
		&assignment.TimeSpentHours, &assignment.Diagnosis, &assignment.ActionsTaken, &partsUsedJSON,
		&assignment.CustomerRating, &assignment.CustomerFeedback, &assignment.Notes,
		&assignment.CreatedAt, &assignment.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrAssignmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment: %w", err)
	}

	// Unmarshal JSONB fields
	if len(partsUsedJSON) > 0 {
		json.Unmarshal(partsUsedJSON, &assignment.PartsUsed)
	}

	return assignment, nil
}

// Update updates an existing assignment
func (r *AssignmentRepository) Update(ctx context.Context, assignment *domain.EngineerAssignment) error {
	// Marshal JSONB fields
	partsUsed, _ := json.Marshal(assignment.PartsUsed)

	query := `
		UPDATE engineer_assignments
		SET 
			status = $2,
			accepted_at = $3,
			rejected_at = $4,
			rejection_reason = $5,
			started_at = $6,
			completed_at = $7,
			completion_status = $8,
			escalation_reason = $9,
			time_spent_hours = $10,
			diagnosis = $11,
			actions_taken = $12,
			parts_used = $13,
			customer_rating = $14,
			customer_feedback = $15,
			notes = $16,
			updated_at = $17
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		assignment.ID, assignment.Status,
		assignment.AcceptedAt, assignment.RejectedAt, assignment.RejectionReason,
		assignment.StartedAt, assignment.CompletedAt, assignment.CompletionStatus, assignment.EscalationReason,
		assignment.TimeSpentHours, assignment.Diagnosis, assignment.ActionsTaken, partsUsed,
		assignment.CustomerRating, assignment.CustomerFeedback, assignment.Notes,
		assignment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrAssignmentNotFound
	}

	return nil
}

// Delete deletes an assignment (soft delete in practice)
func (r *AssignmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM engineer_assignments WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrAssignmentNotFound
	}

	return nil
}

// GetCurrentAssignmentByTicketID retrieves the current active assignment for a ticket
func (r *AssignmentRepository) GetCurrentAssignmentByTicketID(ctx context.Context, ticketID string) (*domain.EngineerAssignment, error) {
	query := `
		SELECT 
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM engineer_assignments
		WHERE ticket_id = $1
		  AND status NOT IN ('completed', 'rejected', 'failed', 'escalated')
		ORDER BY assigned_at DESC
		LIMIT 1
	`

	assignment := &domain.EngineerAssignment{}
	var partsUsedJSON []byte

	err := r.pool.QueryRow(ctx, query, ticketID).Scan(
		&assignment.ID, &assignment.TicketID, &assignment.EngineerID, &assignment.EquipmentID,
		&assignment.AssignmentSequence, &assignment.AssignmentTier, &assignment.AssignmentTierName, &assignment.AssignmentReason,
		&assignment.AssignmentType, &assignment.Status, &assignment.AssignedBy, &assignment.AssignedAt,
		&assignment.AcceptedAt, &assignment.RejectedAt, &assignment.RejectionReason,
		&assignment.StartedAt, &assignment.CompletedAt, &assignment.CompletionStatus, &assignment.EscalationReason,
		&assignment.TimeSpentHours, &assignment.Diagnosis, &assignment.ActionsTaken, &partsUsedJSON,
		&assignment.CustomerRating, &assignment.CustomerFeedback, &assignment.Notes,
		&assignment.CreatedAt, &assignment.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no active assignment found for ticket %s", ticketID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get current assignment: %w", err)
	}

	// Unmarshal JSONB fields
	if len(partsUsedJSON) > 0 {
		json.Unmarshal(partsUsedJSON, &assignment.PartsUsed)
	}

	return assignment, nil
}

// GetAssignmentHistoryByTicketID retrieves all assignments for a ticket
func (r *AssignmentRepository) GetAssignmentHistoryByTicketID(ctx context.Context, ticketID string) ([]*domain.EngineerAssignment, error) {
	query := `
		SELECT 
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM engineer_assignments
		WHERE ticket_id = $1
		ORDER BY assignment_sequence ASC, assigned_at ASC
	`

	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment history: %w", err)
	}
	defer rows.Close()

	var assignments []*domain.EngineerAssignment
	for rows.Next() {
		assignment := &domain.EngineerAssignment{}
		var partsUsedJSON []byte

		err := rows.Scan(
			&assignment.ID, &assignment.TicketID, &assignment.EngineerID, &assignment.EquipmentID,
			&assignment.AssignmentSequence, &assignment.AssignmentTier, &assignment.AssignmentTierName, &assignment.AssignmentReason,
			&assignment.AssignmentType, &assignment.Status, &assignment.AssignedBy, &assignment.AssignedAt,
			&assignment.AcceptedAt, &assignment.RejectedAt, &assignment.RejectionReason,
			&assignment.StartedAt, &assignment.CompletedAt, &assignment.CompletionStatus, &assignment.EscalationReason,
			&assignment.TimeSpentHours, &assignment.Diagnosis, &assignment.ActionsTaken, &partsUsedJSON,
			&assignment.CustomerRating, &assignment.CustomerFeedback, &assignment.Notes,
			&assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}

		// Unmarshal JSONB fields
		if len(partsUsedJSON) > 0 {
			json.Unmarshal(partsUsedJSON, &assignment.PartsUsed)
		}

		assignments = append(assignments, assignment)
	}

	return assignments, rows.Err()
}

// GetAssignmentsByEngineerID retrieves assignments for an engineer
func (r *AssignmentRepository) GetAssignmentsByEngineerID(ctx context.Context, engineerID string, limit int) ([]*domain.EngineerAssignment, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT 
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM engineer_assignments
		WHERE engineer_id = $1
		ORDER BY assigned_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, engineerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get engineer assignments: %w", err)
	}
	defer rows.Close()

	var assignments []*domain.EngineerAssignment
	for rows.Next() {
		assignment := &domain.EngineerAssignment{}
		var partsUsedJSON []byte

		err := rows.Scan(
			&assignment.ID, &assignment.TicketID, &assignment.EngineerID, &assignment.EquipmentID,
			&assignment.AssignmentSequence, &assignment.AssignmentTier, &assignment.AssignmentTierName, &assignment.AssignmentReason,
			&assignment.AssignmentType, &assignment.Status, &assignment.AssignedBy, &assignment.AssignedAt,
			&assignment.AcceptedAt, &assignment.RejectedAt, &assignment.RejectionReason,
			&assignment.StartedAt, &assignment.CompletedAt, &assignment.CompletionStatus, &assignment.EscalationReason,
			&assignment.TimeSpentHours, &assignment.Diagnosis, &assignment.ActionsTaken, &partsUsedJSON,
			&assignment.CustomerRating, &assignment.CustomerFeedback, &assignment.Notes,
			&assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}

		// Unmarshal JSONB fields
		if len(partsUsedJSON) > 0 {
			json.Unmarshal(partsUsedJSON, &assignment.PartsUsed)
		}

		assignments = append(assignments, assignment)
	}

	return assignments, rows.Err()
}

// GetActiveAssignmentsByEngineerID retrieves active assignments for an engineer
func (r *AssignmentRepository) GetActiveAssignmentsByEngineerID(ctx context.Context, engineerID string) ([]*domain.EngineerAssignment, error) {
	query := `
		SELECT 
			id, ticket_id, engineer_id, equipment_id,
			assignment_sequence, assignment_tier, assignment_tier_name, assignment_reason,
			assignment_type, status, assigned_by, assigned_at,
			accepted_at, rejected_at, rejection_reason,
			started_at, completed_at, completion_status, escalation_reason,
			time_spent_hours, diagnosis, actions_taken, parts_used,
			customer_rating, customer_feedback, notes,
			created_at, updated_at
		FROM engineer_assignments
		WHERE engineer_id = $1
		  AND status NOT IN ('completed', 'rejected', 'failed', 'escalated')
		ORDER BY assigned_at DESC
	`

	rows, err := r.pool.Query(ctx, query, engineerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active assignments: %w", err)
	}
	defer rows.Close()

	var assignments []*domain.EngineerAssignment
	for rows.Next() {
		assignment := &domain.EngineerAssignment{}
		var partsUsedJSON []byte

		err := rows.Scan(
			&assignment.ID, &assignment.TicketID, &assignment.EngineerID, &assignment.EquipmentID,
			&assignment.AssignmentSequence, &assignment.AssignmentTier, &assignment.AssignmentTierName, &assignment.AssignmentReason,
			&assignment.AssignmentType, &assignment.Status, &assignment.AssignedBy, &assignment.AssignedAt,
			&assignment.AcceptedAt, &assignment.RejectedAt, &assignment.RejectionReason,
			&assignment.StartedAt, &assignment.CompletedAt, &assignment.CompletionStatus, &assignment.EscalationReason,
			&assignment.TimeSpentHours, &assignment.Diagnosis, &assignment.ActionsTaken, &partsUsedJSON,
			&assignment.CustomerRating, &assignment.CustomerFeedback, &assignment.Notes,
			&assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}

		// Unmarshal JSONB fields
		if len(partsUsedJSON) > 0 {
			json.Unmarshal(partsUsedJSON, &assignment.PartsUsed)
		}

		assignments = append(assignments, assignment)
	}

	return assignments, rows.Err()
}

// CountActiveAssignmentsByEngineerID counts active assignments for an engineer
func (r *AssignmentRepository) CountActiveAssignmentsByEngineerID(ctx context.Context, engineerID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM engineer_assignments
		WHERE engineer_id = $1
		  AND status NOT IN ('completed', 'rejected', 'failed', 'escalated')
	`

	var count int
	err := r.pool.QueryRow(ctx, query, engineerID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active assignments: %w", err)
	}

	return count, nil
}

// GetEngineerWorkload returns workload statistics for an engineer
func (r *AssignmentRepository) GetEngineerWorkload(ctx context.Context, engineerID string) (int, float64, error) {
	query := `
		SELECT 
			COUNT(*) as total_assignments,
			COALESCE(AVG(time_spent_hours), 0) as avg_hours
		FROM engineer_assignments
		WHERE engineer_id = $1
		  AND status = 'completed'
	`

	var count int
	var avgHours sql.NullFloat64

	err := r.pool.QueryRow(ctx, query, engineerID).Scan(&count, &avgHours)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get engineer workload: %w", err)
	}

	return count, avgHours.Float64, nil
}
