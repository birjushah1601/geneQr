package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

// AssignmentService provides business logic for engineer assignment
type AssignmentService struct {
	assignRepo domain.EngineerSuggestionRepository
	ticketRepo domain.TicketRepository
	logger     *slog.Logger
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(
	assignRepo domain.EngineerSuggestionRepository,
	ticketRepo domain.TicketRepository,
	logger *slog.Logger,
) *AssignmentService {
	return &AssignmentService{
		assignRepo: assignRepo,
		ticketRepo: ticketRepo,
		logger:     logger.With(slog.String("component", "assignment_service")),
	}
}

// ListEngineers retrieves engineers, optionally filtered by organization
func (s *AssignmentService) ListEngineers(ctx context.Context, organizationID *string, limit, offset int) ([]*domain.Engineer, error) {
	s.logger.Info("Listing engineers", slog.Any("org_id", organizationID))
	return s.assignRepo.ListEngineers(ctx, organizationID, limit, offset)
}

// GetEngineer retrieves a single engineer by ID
func (s *AssignmentService) GetEngineer(ctx context.Context, engineerID string) (*domain.Engineer, error) {
	return s.assignRepo.GetEngineerByID(ctx, engineerID)
}

// UpdateEngineerLevel updates an engineer's skill level
func (s *AssignmentService) UpdateEngineerLevel(ctx context.Context, engineerID string, level domain.EngineerLevel) error {
	s.logger.Info("Updating engineer level", 
		slog.String("engineer_id", engineerID),
		slog.String("level", string(level)))
	return s.assignRepo.UpdateEngineerLevel(ctx, engineerID, level)
}

// ListEngineerEquipmentTypes retrieves all equipment types an engineer can service
func (s *AssignmentService) ListEngineerEquipmentTypes(ctx context.Context, engineerID string) ([]*domain.EngineerEquipmentType, error) {
	return s.assignRepo.ListEngineerEquipmentTypes(ctx, engineerID)
}

// AddEngineerEquipmentType adds an equipment type capability to an engineer
func (s *AssignmentService) AddEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error {
	s.logger.Info("Adding engineer equipment type",
		slog.String("engineer_id", engineerID),
		slog.String("manufacturer", manufacturer),
		slog.String("category", category))
	
	// Validate engineer exists
	_, err := s.assignRepo.GetEngineerByID(ctx, engineerID)
	if err != nil {
		return fmt.Errorf("engineer not found: %w", err)
	}
	
	return s.assignRepo.AddEngineerEquipmentType(ctx, engineerID, manufacturer, category)
}

// RemoveEngineerEquipmentType removes an equipment type capability from an engineer
func (s *AssignmentService) RemoveEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error {
	s.logger.Info("Removing engineer equipment type",
		slog.String("engineer_id", engineerID),
		slog.String("manufacturer", manufacturer),
		slog.String("category", category))
	return s.assignRepo.RemoveEngineerEquipmentType(ctx, engineerID, manufacturer, category)
}

// GetEquipmentServiceConfig retrieves service configuration for equipment
func (s *AssignmentService) GetEquipmentServiceConfig(ctx context.Context, equipmentID string) (*domain.EquipmentServiceConfig, error) {
	return s.assignRepo.GetEquipmentServiceConfig(ctx, equipmentID)
}

// CreateEquipmentServiceConfig creates a new service configuration
func (s *AssignmentService) CreateEquipmentServiceConfig(ctx context.Context, config *domain.EquipmentServiceConfig) error {
	s.logger.Info("Creating equipment service config", slog.String("equipment_id", config.EquipmentID))
	return s.assignRepo.CreateEquipmentServiceConfig(ctx, config)
}

// UpdateEquipmentServiceConfig updates an existing service configuration
func (s *AssignmentService) UpdateEquipmentServiceConfig(ctx context.Context, config *domain.EquipmentServiceConfig) error {
	s.logger.Info("Updating equipment service config", slog.String("config_id", config.ID))
	return s.assignRepo.UpdateEquipmentServiceConfig(ctx, config)
}

// GetSuggestedEngineers retrieves suggested engineers for a service ticket
func (s *AssignmentService) GetSuggestedEngineers(ctx context.Context, ticketID string) ([]*domain.SuggestedEngineer, error) {
	s.logger.Info("Getting suggested engineers", slog.String("ticket_id", ticketID))
	
	// Get ticket details
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}
	
	// Default to L1 minimum level (can be enhanced with priority mapping)
	minLevel := domain.EngineerLevelL1
	switch ticket.Priority {
	case domain.PriorityCritical:
		minLevel = domain.EngineerLevelL3
	case domain.PriorityHigh:
		minLevel = domain.EngineerLevelL2
	default:
		minLevel = domain.EngineerLevelL1
	}
	
	// Extract manufacturer and category from equipment
	_, manufacturerName, category, err := s.assignRepo.GetEquipmentDetails(ctx, ticket.EquipmentID)
	if err != nil {
		s.logger.Warn("Failed to get equipment details, continuing with empty manufacturer/category",
			slog.String("equipment_id", ticket.EquipmentID),
			slog.String("error", err.Error()))
		manufacturerName = ""
		category = ""
	}
	
	s.logger.Info("Retrieved equipment details for engineer assignment",
		slog.String("equipment_id", ticket.EquipmentID),
		slog.String("manufacturer", manufacturerName),
		slog.String("category", category))
	
	// Get suggestions from repository
	suggestions, err := s.assignRepo.GetSuggestedEngineers(
		ctx,
		ticket.EquipmentID,
		manufacturerName,
		category,
		minLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}
	
	s.logger.Info("Found suggested engineers",
		slog.String("ticket_id", ticketID),
		slog.Int("count", len(suggestions)))
	
	return suggestions, nil
}

// AssignEngineer assigns an engineer to a service ticket
func (s *AssignmentService) AssignEngineer(ctx context.Context, req AssignEngineerRequest) error {
	s.logger.Info("Assigning engineer to ticket",
		slog.String("ticket_id", req.TicketID),
		slog.String("engineer_id", req.EngineerID))
	
	// Validate ticket exists
	ticket, err := s.ticketRepo.GetByID(ctx, req.TicketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}
	
	// Validate engineer exists
	engineer, err := s.assignRepo.GetEngineerByID(ctx, req.EngineerID)
	if err != nil {
		return fmt.Errorf("engineer not found: %w", err)
	}
	
	// Build assignment request
	assignReq := domain.AssignmentRequest{
		TicketID:           req.TicketID,
		EngineerID:         engineer.ID,
		EngineerName:       engineer.Name,
		OrganizationID:     engineer.OrganizationID,
		AssignmentTier:     req.AssignmentTier,
		AssignmentTierName: req.AssignmentTierName,
		AssignedBy:         req.AssignedBy,
	}
	
	// Perform assignment
	err = s.assignRepo.AssignEngineerToTicket(ctx, assignReq)
	if err != nil {
		return fmt.Errorf("failed to assign engineer: %w", err)
	}
	
	// Add status history
	history := &domain.StatusHistory{
		TicketID:   req.TicketID,
		FromStatus: string(ticket.Status),
		ToStatus:   "assigned",
		ChangedBy:  req.AssignedBy,
		Reason:     fmt.Sprintf("Assigned to %s (%s) - %s", engineer.Name, engineer.OrganizationName, req.AssignmentTierName),
	}
	s.ticketRepo.AddStatusHistory(ctx, history)
	
	// Add comment
	comment := &domain.TicketComment{
		TicketID:    req.TicketID,
		CommentType: "system",
		AuthorName:  "System",
		Comment:     fmt.Sprintf("Ticket assigned to %s from %s (%s)", engineer.Name, engineer.OrganizationName, req.AssignmentTierName),
	}
	s.ticketRepo.AddComment(ctx, comment)
	
	s.logger.Info("Engineer assigned successfully",
		slog.String("ticket_id", req.TicketID),
		slog.String("engineer_id", engineer.ID))
	
	return nil
}

// Request DTOs

type AssignEngineerRequest struct {
	TicketID           string                `json:"ticket_id"`
	EngineerID         string                `json:"engineer_id"`
	AssignmentTier     string                `json:"assignment_tier"`
	AssignmentTierName string                `json:"assignment_tier_name"`
	AssignedBy         string                `json:"assigned_by"`
}

type AddEquipmentTypeRequest struct {
	EngineerID   string `json:"engineer_id"`
	Manufacturer string `json:"manufacturer"`
	Category     string `json:"category"`
}

type RemoveEquipmentTypeRequest struct {
	EngineerID   string `json:"engineer_id"`
	Manufacturer string `json:"manufacturer"`
	Category     string `json:"category"`
}

// GetAssignmentHistory retrieves assignment history for a ticket
func (s *AssignmentService) GetAssignmentHistory(ctx context.Context, ticketID string) ([]*domain.EngineerAssignment, error) {
	s.logger.Info("Getting assignment history",
		slog.String("ticket_id", ticketID))
	
	history, err := s.assignRepo.GetAssignmentHistoryByTicketID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment history: %w", err)
	}
	
	s.logger.Info("Retrieved assignment history",
		slog.String("ticket_id", ticketID),
		slog.Int("count", len(history)))
	
	return history, nil
}
