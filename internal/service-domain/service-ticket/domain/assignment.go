package domain

import (
	"context"
	"errors"
	"time"
)

// ==================== ERROR DEFINITIONS ====================

var (
	ErrAssignmentNotFound      = errors.New("assignment not found")
	ErrInvalidAssignmentStatus = errors.New("invalid assignment status transition")
	ErrTicketAlreadyAssigned   = errors.New("ticket already has an active assignment")
)

// ==================== ENGINEER LEVEL ====================

// EngineerLevel represents engineer skill level
type EngineerLevel string

const (
	EngineerLevelL1 EngineerLevel = "L1"
	EngineerLevelL2 EngineerLevel = "L2"
	EngineerLevelL3 EngineerLevel = "L3"
)

// ==================== ASSIGNMENT WORKFLOW SYSTEM (from main) ====================

// AssignmentStatus represents the lifecycle status of an engineer assignment
type AssignmentStatus string

const (
	AssignmentStatusAssigned   AssignmentStatus = "assigned"
	AssignmentStatusAccepted   AssignmentStatus = "accepted"
	AssignmentStatusRejected   AssignmentStatus = "rejected"
	AssignmentStatusInProgress AssignmentStatus = "in_progress"
	AssignmentStatusCompleted  AssignmentStatus = "completed"
	AssignmentStatusFailed     AssignmentStatus = "failed"
	AssignmentStatusEscalated  AssignmentStatus = "escalated"
)

// AssignmentType represents how the assignment was created
type AssignmentType string

const (
	AssignmentTypeAuto       AssignmentType = "auto"
	AssignmentTypeManual     AssignmentType = "manual"
	AssignmentTypeEscalation AssignmentType = "escalation"
)

// CompletionStatus represents the outcome of an assignment
type CompletionStatus string

const (
	CompletionStatusSuccess             CompletionStatus = "success"
	CompletionStatusFailed              CompletionStatus = "failed"
	CompletionStatusEscalated           CompletionStatus = "escalated"
	CompletionStatusPartsRequired       CompletionStatus = "parts_required"
	CompletionStatusCustomerUnavailable CompletionStatus = "customer_unavailable"
)

// Note: Part type is defined in ticket.go to avoid duplication

// EngineerAssignment represents an engineer assigned to work on a service ticket
type EngineerAssignment struct {
	ID           string `json:"id"`
	TicketID     string `json:"ticket_id"`
	EngineerID   string `json:"engineer_id"`
	EngineerName string `json:"engineer_name"`
	EquipmentID  string `json:"equipment_id"`

	// Sequence tracking
	AssignmentSequence int    `json:"assignment_sequence"` // 1, 2, 3... for escalations
	AssignmentTier     int    `json:"assignment_tier"`     // 1=OEM, 2=Sub-sub_SUB_DEALER, etc.
	AssignmentTierName string `json:"assignment_tier_name"`
	AssignmentReason   string `json:"assignment_reason"` // "Initial", "Escalation", etc.

	// Workflow
	AssignmentType  AssignmentType   `json:"assignment_type"`
	Status          AssignmentStatus `json:"status"`
	AssignedBy      string           `json:"assigned_by"`
	AssignedAt      time.Time        `json:"assigned_at"`
	AcceptedAt      *time.Time       `json:"accepted_at,omitempty"`
	RejectedAt      *time.Time       `json:"rejected_at,omitempty"`
	RejectionReason string           `json:"rejection_reason,omitempty"`

	// Execution
	StartedAt        *time.Time       `json:"started_at,omitempty"`
	CompletedAt      *time.Time       `json:"completed_at,omitempty"`
	CompletionStatus CompletionStatus `json:"completion_status,omitempty"`
	EscalationReason string           `json:"escalation_reason,omitempty"`
	TimeSpentHours   float64          `json:"time_spent_hours"`

	// Details
	Diagnosis    string `json:"diagnosis,omitempty"`
	ActionsTaken string `json:"actions_taken,omitempty"`
	PartsUsed    []Part `json:"parts_used,omitempty"`

	// Customer feedback
	CustomerRating   int    `json:"customer_rating,omitempty"` // 1-5
	CustomerFeedback string `json:"customer_feedback,omitempty"`

	// Metadata
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewEngineerAssignment creates a new engineer assignment
func NewEngineerAssignment(
	ticketID, engineerID, equipmentID, assignedBy string,
	sequence, tier int,
	tierName, reason string,
	assignmentType AssignmentType,
) *EngineerAssignment {
	now := time.Now()
	return &EngineerAssignment{
		TicketID:           ticketID,
		EngineerID:         engineerID,
		EquipmentID:        equipmentID,
		AssignmentSequence: sequence,
		AssignmentTier:     tier,
		AssignmentTierName: tierName,
		AssignmentReason:   reason,
		AssignmentType:     assignmentType,
		Status:             AssignmentStatusAssigned,
		AssignedBy:         assignedBy,
		AssignedAt:         now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// Accept marks the assignment as accepted by the engineer
func (a *EngineerAssignment) Accept() error {
	if a.Status != AssignmentStatusAssigned {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.AcceptedAt = &now
	a.Status = AssignmentStatusAccepted
	a.UpdatedAt = now

	return nil
}

// Reject marks the assignment as rejected by the engineer
func (a *EngineerAssignment) Reject(reason string) error {
	if a.Status != AssignmentStatusAssigned && a.Status != AssignmentStatusAccepted {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.RejectedAt = &now
	a.RejectionReason = reason
	a.Status = AssignmentStatusRejected
	a.UpdatedAt = now

	return nil
}

// Start marks the assignment as in progress
func (a *EngineerAssignment) Start() error {
	if a.Status != AssignmentStatusAssigned && a.Status != AssignmentStatusAccepted {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.StartedAt = &now
	a.Status = AssignmentStatusInProgress
	a.UpdatedAt = now

	return nil
}

// Complete marks the assignment as completed
func (a *EngineerAssignment) Complete(
	completionStatus CompletionStatus,
	diagnosis, actionsTaken string,
	partsUsed []Part,
	timeSpentHours float64,
) error {
	if a.Status != AssignmentStatusInProgress {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.CompletedAt = &now
	a.CompletionStatus = completionStatus
	a.Diagnosis = diagnosis
	a.ActionsTaken = actionsTaken
	a.PartsUsed = partsUsed
	a.TimeSpentHours = timeSpentHours
	a.Status = AssignmentStatusCompleted
	a.UpdatedAt = now

	return nil
}

// Escalate marks the assignment as escalated to next tier
func (a *EngineerAssignment) Escalate(reason string) error {
	if a.Status != AssignmentStatusInProgress && a.Status != AssignmentStatusAccepted {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.CompletedAt = &now
	a.EscalationReason = reason
	a.CompletionStatus = CompletionStatusEscalated
	a.Status = AssignmentStatusEscalated
	a.UpdatedAt = now

	return nil
}

// Fail marks the assignment as failed
func (a *EngineerAssignment) Fail(reason string) error {
	if a.Status != AssignmentStatusInProgress {
		return ErrInvalidAssignmentStatus
	}

	now := time.Now()
	a.CompletedAt = &now
	a.EscalationReason = reason
	a.CompletionStatus = CompletionStatusFailed
	a.Status = AssignmentStatusFailed
	a.UpdatedAt = now

	return nil
}

// AddCustomerFeedback adds customer rating and feedback
func (a *EngineerAssignment) AddCustomerFeedback(rating int, feedback string) error {
	if a.Status != AssignmentStatusCompleted {
		return ErrInvalidAssignmentStatus
	}

	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	a.CustomerRating = rating
	a.CustomerFeedback = feedback
	a.UpdatedAt = time.Now()

	return nil
}

// IsActive returns true if the assignment is currently active (not completed/rejected/failed)
func (a *EngineerAssignment) IsActive() bool {
	return a.Status != AssignmentStatusCompleted &&
		a.Status != AssignmentStatusRejected &&
		a.Status != AssignmentStatusFailed &&
		a.Status != AssignmentStatusEscalated
}

// CanEscalate returns true if the assignment can be escalated
func (a *EngineerAssignment) CanEscalate() bool {
	return a.Status == AssignmentStatusInProgress || a.Status == AssignmentStatusAccepted
}

// AssignmentRepository defines the interface for assignment workflow persistence
type AssignmentRepository interface {
	Create(ctx context.Context, assignment *EngineerAssignment) error
	GetByID(ctx context.Context, id string) (*EngineerAssignment, error)
	Update(ctx context.Context, assignment *EngineerAssignment) error
	Delete(ctx context.Context, id string) error

	// Query methods
	GetCurrentAssignmentByTicketID(ctx context.Context, ticketID string) (*EngineerAssignment, error)
	GetAssignmentHistoryByTicketID(ctx context.Context, ticketID string) ([]*EngineerAssignment, error)
	GetAssignmentsByEngineerID(ctx context.Context, engineerID string, limit int) ([]*EngineerAssignment, error)
	GetActiveAssignmentsByEngineerID(ctx context.Context, engineerID string) ([]*EngineerAssignment, error)

	// Statistics
	CountActiveAssignmentsByEngineerID(ctx context.Context, engineerID string) (int, error)
	GetEngineerWorkload(ctx context.Context, engineerID string) (int, float64, error) // count, avg hours
}

// ==================== ENGINEER SUGGESTION SYSTEM (from feat/database-refactor-phase1) ====================

// Engineer represents an engineer from the organizations table with assignment-specific data
type Engineer struct {
	ID                   string        `json:"id"`
	OrganizationID       string        `json:"organization_id"`
	OrganizationName     string        `json:"organization_name,omitempty"`
	Name                 string        `json:"name"`
	Email                string        `json:"email"`
	Phone                string        `json:"phone"`
	EngineerLevel        EngineerLevel `json:"engineer_level"`
	IsActive             bool          `json:"is_active"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`

	// For eligible engineers list
	EligibleEquipmentTypes []string `json:"eligible_equipment_types,omitempty"`
}

// EngineerEquipmentType maps engineers to equipment types they can service
type EngineerEquipmentType struct {
	ID           string    `json:"id"`
	EngineerID   string    `json:"engineer_id"`
	Manufacturer string    `json:"manufacturer"` // e.g., "Siemens"
	Category     string    `json:"category"`     // e.g., "MRI", "CT", "X-Ray"
	CreatedAt    time.Time `json:"created_at"`
}

// EquipmentServiceConfig defines service hierarchy for specific equipment
type EquipmentServiceConfig struct {
	ID                    string    `json:"id"`
	EquipmentID           string    `json:"equipment_id"`
	UnderWarranty         bool      `json:"under_warranty"`
	UnderAMC              bool      `json:"under_amc"`
	PrimaryServiceOrgID   *string   `json:"primary_service_org_id"`
	SecondaryServiceOrgID *string   `json:"secondary_service_org_id"`
	TertiaryServiceOrgID  *string   `json:"tertiary_service_org_id"`
	FallbackServiceOrgID  *string   `json:"fallback_service_org_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// SuggestedEngineer represents an engineer suggestion for assignment
type SuggestedEngineer struct {
	EngineerID         string        `json:"engineer_id"`
	EngineerName       string        `json:"engineer_name"`
	OrganizationID     string        `json:"organization_id"`
	OrganizationName   string        `json:"organization_name"`
	EngineerLevel      EngineerLevel `json:"engineer_level"`
	AssignmentTier     string        `json:"assignment_tier"`      // "warranty_primary", "amc_primary", "secondary", "tertiary", "fallback"
	AssignmentTierName string        `json:"assignment_tier_name"` // Human-readable tier name
	MatchReason        string        `json:"match_reason"`         // Why this engineer was suggested
	Priority           int           `json:"priority"`             // Lower is higher priority
}

// AssignmentRequest represents a manual assignment request
type AssignmentRequest struct {
	TicketID           string `json:"ticket_id"`
	EngineerID         string `json:"engineer_id"`
	EngineerName       string `json:"engineer_name"`
	OrganizationID     string `json:"organization_id"`
	AssignmentTier     string `json:"assignment_tier"`
	AssignmentTierName string `json:"assignment_tier_name"`
	AssignedBy         string `json:"assigned_by"`
}
