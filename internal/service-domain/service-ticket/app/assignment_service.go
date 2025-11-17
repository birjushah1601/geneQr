package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

// AssignmentService provides business logic for engineer assignments
type AssignmentService struct {
	assignmentRepo domain.AssignmentRepository
	ticketRepo     domain.TicketRepository
	logger         *slog.Logger
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(
	assignmentRepo domain.AssignmentRepository,
	ticketRepo domain.TicketRepository,
	logger *slog.Logger,
) *AssignmentService {
	return &AssignmentService{
		assignmentRepo: assignmentRepo,
		ticketRepo:     ticketRepo,
		logger:         logger.With(slog.String("component", "assignment_service")),
	}
}

// AssignTicketRequest represents a request to assign a ticket to an engineer
type AssignTicketRequest struct {
	TicketID    string                 `json:"ticket_id"`
	EngineerID  string                 `json:"engineer_id"`
	EquipmentID string                 `json:"equipment_id"`
	Tier        int                    `json:"tier"`
	TierName    string                 `json:"tier_name"`
	Reason      string                 `json:"reason"`
	Type        domain.AssignmentType  `json:"type"` // "auto", "manual"
	AssignedBy  string                 `json:"assigned_by"`
}

// AssignTicket assigns a ticket to an engineer
func (s *AssignmentService) AssignTicket(ctx context.Context, req AssignTicketRequest) (*domain.EngineerAssignment, error) {
	s.logger.Info("Assigning ticket to engineer",
		slog.String("ticket_id", req.TicketID),
		slog.String("engineer_id", req.EngineerID),
		slog.Int("tier", req.Tier))

	// Get current assignment history to determine sequence
	history, err := s.assignmentRepo.GetAssignmentHistoryByTicketID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment history: %w", err)
	}

	// Create new assignment
	assignment := domain.NewEngineerAssignment(
		req.TicketID,
		req.EngineerID,
		req.EquipmentID,
		req.AssignedBy,
		len(history)+1, // Sequence
		req.Tier,
		req.TierName,
		req.Reason,
		req.Type,
	)

	// Save assignment
	if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	// Update ticket status to "assigned"
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	if err := ticket.AssignEngineer(req.EngineerID, ""); err != nil {
		s.logger.Warn("Failed to update ticket status", slog.String("error", err.Error()))
		// Continue - assignment is created, ticket update is best effort for backward compatibility
	} else {
		if err := s.ticketRepo.Update(ctx, ticket); err != nil {
			s.logger.Warn("Failed to persist ticket update", slog.String("error", err.Error()))
		}
	}

	s.logger.Info("Ticket assigned successfully",
		slog.String("assignment_id", assignment.ID),
		slog.Int("sequence", assignment.AssignmentSequence))

	return assignment, nil
}

// EscalateTicketRequest represents a request to escalate a ticket
type EscalateTicketRequest struct {
	TicketID        string `json:"ticket_id"`
	Reason          string `json:"reason"`
	NextEngineerID  string `json:"next_engineer_id"`
	NextTierName    string `json:"next_tier_name"`
	EscalatedBy     string `json:"escalated_by"`
}

// EscalateTicket escalates a ticket to the next tier
func (s *AssignmentService) EscalateTicket(ctx context.Context, req EscalateTicketRequest) (*domain.EngineerAssignment, error) {
	s.logger.Info("Escalating ticket",
		slog.String("ticket_id", req.TicketID),
		slog.String("reason", req.Reason))

	// Get current assignment
	current, err := s.assignmentRepo.GetCurrentAssignmentByTicketID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current assignment: %w", err)
	}

	if !current.CanEscalate() {
		return nil, fmt.Errorf("assignment cannot be escalated in current status: %s", current.Status)
	}

	// Mark current assignment as escalated
	if err := current.Escalate(req.Reason); err != nil {
		return nil, fmt.Errorf("failed to escalate current assignment: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, current); err != nil {
		return nil, fmt.Errorf("failed to update current assignment: %w", err)
	}

	// Get ticket to find equipment_id
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Create new assignment at higher tier
	newAssignment, err := s.AssignTicket(ctx, AssignTicketRequest{
		TicketID:    req.TicketID,
		EngineerID:  req.NextEngineerID,
		EquipmentID: ticket.EquipmentID,
		Tier:        current.AssignmentTier + 1,
		TierName:    req.NextTierName,
		Reason:      "Escalation from tier " + fmt.Sprintf("%d", current.AssignmentTier),
		Type:        domain.AssignmentTypeEscalation,
		AssignedBy:  req.EscalatedBy,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create escalated assignment: %w", err)
	}

	s.logger.Info("Ticket escalated successfully",
		slog.String("new_assignment_id", newAssignment.ID),
		slog.Int("new_tier", newAssignment.AssignmentTier))

	return newAssignment, nil
}

// AcceptAssignmentRequest represents a request to accept an assignment
type AcceptAssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
	EngineerID   string `json:"engineer_id"`
}

// AcceptAssignment marks an assignment as accepted by the engineer
func (s *AssignmentService) AcceptAssignment(ctx context.Context, req AcceptAssignmentRequest) error {
	s.logger.Info("Engineer accepting assignment",
		slog.String("assignment_id", req.AssignmentID),
		slog.String("engineer_id", req.EngineerID))

	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Verify engineer
	if assignment.EngineerID != req.EngineerID {
		return fmt.Errorf("assignment not assigned to this engineer")
	}

	if err := assignment.Accept(); err != nil {
		return fmt.Errorf("failed to accept assignment: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	s.logger.Info("Assignment accepted", slog.String("assignment_id", req.AssignmentID))
	return nil
}

// RejectAssignmentRequest represents a request to reject an assignment
type RejectAssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
	EngineerID   string `json:"engineer_id"`
	Reason       string `json:"reason"`
}

// RejectAssignment marks an assignment as rejected by the engineer
func (s *AssignmentService) RejectAssignment(ctx context.Context, req RejectAssignmentRequest) error {
	s.logger.Info("Engineer rejecting assignment",
		slog.String("assignment_id", req.AssignmentID),
		slog.String("reason", req.Reason))

	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Verify engineer
	if assignment.EngineerID != req.EngineerID {
		return fmt.Errorf("assignment not assigned to this engineer")
	}

	if err := assignment.Reject(req.Reason); err != nil {
		return fmt.Errorf("failed to reject assignment: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	s.logger.Info("Assignment rejected", slog.String("assignment_id", req.AssignmentID))
	return nil
}

// StartAssignmentRequest represents a request to start work on an assignment
type StartAssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
	EngineerID   string `json:"engineer_id"`
}

// StartAssignment marks an assignment as in progress
func (s *AssignmentService) StartAssignment(ctx context.Context, req StartAssignmentRequest) error {
	s.logger.Info("Engineer starting assignment",
		slog.String("assignment_id", req.AssignmentID))

	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Verify engineer
	if assignment.EngineerID != req.EngineerID {
		return fmt.Errorf("assignment not assigned to this engineer")
	}

	if err := assignment.Start(); err != nil {
		return fmt.Errorf("failed to start assignment: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	// Update ticket status to in_progress
	ticket, err := s.ticketRepo.GetByID(ctx, assignment.TicketID)
	if err == nil {
		if err := ticket.Start(); err == nil {
			s.ticketRepo.Update(ctx, ticket)
		}
	}

	s.logger.Info("Assignment started", slog.String("assignment_id", req.AssignmentID))
	return nil
}

// CompleteAssignmentRequest represents a request to complete an assignment
type CompleteAssignmentRequest struct {
	AssignmentID     string                   `json:"assignment_id"`
	EngineerID       string                   `json:"engineer_id"`
	CompletionStatus domain.CompletionStatus  `json:"completion_status"`
	Diagnosis        string                   `json:"diagnosis"`
	ActionsTaken     string                   `json:"actions_taken"`
	PartsUsed        []domain.Part            `json:"parts_used"`
	TimeSpentHours   float64                  `json:"time_spent_hours"`
}

// CompleteAssignment marks an assignment as completed
func (s *AssignmentService) CompleteAssignment(ctx context.Context, req CompleteAssignmentRequest) error {
	s.logger.Info("Engineer completing assignment",
		slog.String("assignment_id", req.AssignmentID),
		slog.String("status", string(req.CompletionStatus)))

	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// Verify engineer
	if assignment.EngineerID != req.EngineerID {
		return fmt.Errorf("assignment not assigned to this engineer")
	}

	if err := assignment.Complete(
		req.CompletionStatus,
		req.Diagnosis,
		req.ActionsTaken,
		req.PartsUsed,
		req.TimeSpentHours,
	); err != nil {
		return fmt.Errorf("failed to complete assignment: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	// Update ticket if successfully resolved
	if req.CompletionStatus == domain.CompletionStatusSuccess {
		ticket, err := s.ticketRepo.GetByID(ctx, assignment.TicketID)
		if err == nil {
			if err := ticket.Resolve(req.ActionsTaken, req.PartsUsed, req.TimeSpentHours, 0); err == nil {
				s.ticketRepo.Update(ctx, ticket)
			}
		}
	}

	s.logger.Info("Assignment completed", slog.String("assignment_id", req.AssignmentID))
	return nil
}

// AddCustomerFeedbackRequest represents customer feedback for an assignment
type AddCustomerFeedbackRequest struct {
	AssignmentID string `json:"assignment_id"`
	Rating       int    `json:"rating"`
	Feedback     string `json:"feedback"`
}

// AddCustomerFeedback adds customer rating and feedback to a completed assignment
func (s *AssignmentService) AddCustomerFeedback(ctx context.Context, req AddCustomerFeedbackRequest) error {
	s.logger.Info("Adding customer feedback",
		slog.String("assignment_id", req.AssignmentID),
		slog.Int("rating", req.Rating))

	assignment, err := s.assignmentRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	if err := assignment.AddCustomerFeedback(req.Rating, req.Feedback); err != nil {
		return fmt.Errorf("failed to add feedback: %w", err)
	}

	if err := s.assignmentRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	s.logger.Info("Customer feedback added", slog.String("assignment_id", req.AssignmentID))
	return nil
}

// GetCurrentAssignment retrieves the current active assignment for a ticket
func (s *AssignmentService) GetCurrentAssignment(ctx context.Context, ticketID string) (*domain.EngineerAssignment, error) {
	return s.assignmentRepo.GetCurrentAssignmentByTicketID(ctx, ticketID)
}

// GetAssignmentHistory retrieves all assignments for a ticket
func (s *AssignmentService) GetAssignmentHistory(ctx context.Context, ticketID string) ([]*domain.EngineerAssignment, error) {
	return s.assignmentRepo.GetAssignmentHistoryByTicketID(ctx, ticketID)
}

// GetEngineerAssignments retrieves assignments for an engineer
func (s *AssignmentService) GetEngineerAssignments(ctx context.Context, engineerID string, limit int) ([]*domain.EngineerAssignment, error) {
	return s.assignmentRepo.GetAssignmentsByEngineerID(ctx, engineerID, limit)
}

// GetActiveEngineerAssignments retrieves active assignments for an engineer
func (s *AssignmentService) GetActiveEngineerAssignments(ctx context.Context, engineerID string) ([]*domain.EngineerAssignment, error) {
	return s.assignmentRepo.GetActiveAssignmentsByEngineerID(ctx, engineerID)
}

// GetEngineerWorkload returns workload statistics for an engineer
func (s *AssignmentService) GetEngineerWorkload(ctx context.Context, engineerID string) (activeCount int, completedCount int, avgHours float64, error error) {
	activeCount, err := s.assignmentRepo.CountActiveAssignmentsByEngineerID(ctx, engineerID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count active assignments: %w", err)
	}

	completedCount, avgHours, err = s.assignmentRepo.GetEngineerWorkload(ctx, engineerID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get workload: %w", err)
	}

	return activeCount, completedCount, avgHours, nil
}
