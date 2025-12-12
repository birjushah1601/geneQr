package app

import (
	"context"
	"fmt"
	"log/slog"

	equipmentDomain "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MultiModelAssignmentService provides multi-model engineer suggestions
type MultiModelAssignmentService struct {
	assignRepo    domain.EngineerSuggestionRepository
	ticketRepo    domain.TicketRepository
	equipmentRepo equipmentDomain.Repository
	pool          *pgxpool.Pool
	logger        *slog.Logger
}

// NewMultiModelAssignmentService creates a new multi-model assignment service
func NewMultiModelAssignmentService(
	assignRepo domain.EngineerSuggestionRepository,
	ticketRepo domain.TicketRepository,
	equipmentRepo equipmentDomain.Repository,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) *MultiModelAssignmentService {
	return &MultiModelAssignmentService{
		assignRepo:    assignRepo,
		ticketRepo:    ticketRepo,
		equipmentRepo: equipmentRepo,
		pool:          pool,
		logger:        logger.With(slog.String("component", "multi_model_assignment")),
	}
}

// AssignmentSuggestionsResponse holds all assignment models
type AssignmentSuggestionsResponse struct {
	TicketID     string                        `json:"ticket_id"`
	Equipment    *EquipmentContext             `json:"equipment"`
	Ticket       *TicketContext                `json:"ticket"`
	ModelResults map[string]*AssignmentModel   `json:"suggestions_by_model"`
	TierInfo     []*TierInfo                   `json:"assignment_tiers"`
}

// Equipment Context
type EquipmentContext struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Manufacturer string  `json:"manufacturer"`
	Category     string  `json:"category"`
	ModelNumber  string  `json:"model_number"`
	Location     *Location `json:"location,omitempty"`
}

// Location info
type Location struct {
	Region  string   `json:"region"`
	Address string   `json:"address"`
	Lat     *float64 `json:"lat,omitempty"`
	Lng     *float64 `json:"lng,omitempty"`
}

// Ticket Context
type TicketContext struct {
	Priority         string `json:"priority"`
	MinLevelRequired int    `json:"min_level_required"`
	RequiresCertification bool `json:"requires_certification"`
}

// AssignmentModel represents one assignment algorithm's results
type AssignmentModel struct {
	ModelName   string                    `json:"model_name"`
	Description string                    `json:"description"`
	Engineers   []*EngineerSuggestion     `json:"engineers"`
	Count       int                       `json:"count"`
}

// EngineerSuggestion is an enhanced engineer with match details
type EngineerSuggestion struct {
	// Engineer basic info
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	EngineerLevel    int    `json:"engineer_level"`
	Skills           []string `json:"skills,omitempty"`
	HomeRegion       string   `json:"home_region,omitempty"`
	
	// Organization info
	OrganizationID   string `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	OrganizationTier int    `json:"organization_tier,omitempty"`
	
	// Match scoring
	MatchScore   int      `json:"match_score"`
	MatchReasons []string `json:"match_reasons"`
	
	// Workload
	Workload *WorkloadInfo `json:"workload,omitempty"`
	
	// Certifications
	Certifications []*CertificationInfo `json:"certifications,omitempty"`
	
	// Distance (if available)
	DistanceKm *float64 `json:"distance_km,omitempty"`
	TravelTimeMins *int `json:"estimated_travel_time_mins,omitempty"`
}

// WorkloadInfo tracks engineer's current assignments
type WorkloadInfo struct {
	ActiveTickets     int     `json:"active_tickets"`
	InProgressTickets int     `json:"in_progress_tickets"`
	AvgResolutionHours *float64 `json:"avg_resolution_hours,omitempty"`
}

// CertificationInfo tracks certifications
type CertificationInfo struct {
	Manufacturer        string `json:"manufacturer"`
	Category            string `json:"category"`
	IsCertified         bool   `json:"is_certified"`
	CertificationNumber string `json:"certification_number,omitempty"`
	Expiry              string `json:"expiry,omitempty"`
}

// TierInfo describes assignment tiers
type TierInfo struct {
	Tier            int      `json:"tier"`
	Name            string   `json:"name"`
	OrganizationIDs []string `json:"organization_ids"`
	AvailableCount  int      `json:"available_count"`
}

// GetMultiModelSuggestions returns engineer suggestions across all models
func (s *MultiModelAssignmentService) GetMultiModelSuggestions(ctx context.Context, ticketID string) (*AssignmentSuggestionsResponse, error) {
	s.logger.Info("Getting multi-model assignment suggestions", slog.String("ticket_id", ticketID))
	
	// 1. Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}
	
	// 2. Get equipment details
	equipment, err := s.equipmentRepo.GetByID(ctx, ticket.EquipmentID)
	if err != nil {
		return nil, fmt.Errorf("equipment not found: %w", err)
	}
	
	// 3. Build contexts
	equipmentCtx := &EquipmentContext{
		ID:           equipment.ID,
		Name:         equipment.EquipmentName,
		Manufacturer: equipment.ManufacturerName,
		Category:     equipment.Category,
		ModelNumber:  equipment.ModelNumber,
		Location: &Location{
			Region:  equipment.InstallationLocation,
			Address: "", // TODO: extract from installation_address map
		},
	}
	
	minLevel := s.getMinLevelFromPriority(ticket.Priority)
	ticketCtx := &TicketContext{
		Priority:              string(ticket.Priority),
		MinLevelRequired:      minLevel,
		RequiresCertification: ticket.Priority == domain.PriorityCritical,
	}
	
	// 4. Get all engineers
	allEngineers, err := s.assignRepo.ListEngineers(ctx, nil, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list engineers: %w", err)
	}
	
	// 5. Get workload for all engineers
	workloadMap, err := s.getEngineerWorkloads(ctx, allEngineers)
	if err != nil {
		s.logger.Warn("Failed to get workloads", slog.String("error", err.Error()))
		workloadMap = make(map[string]*WorkloadInfo)
	}
	
	// 6. Get certifications for relevant engineers
	certMap, err := s.getEngineerCertifications(ctx, allEngineers, equipment.ManufacturerName, equipment.Category)
	if err != nil {
		s.logger.Warn("Failed to get certifications", slog.String("error", err.Error()))
		certMap = make(map[string][]*CertificationInfo)
	}
	
	// 7. Run all assignment models
	modelResults := make(map[string]*AssignmentModel)
	
	// Model 1: Best Match (combines all factors)
	modelResults["best_match"] = s.getBestMatchModel(allEngineers, equipment, ticket, workloadMap, certMap, minLevel)
	
	// Model 2: Manufacturer Certified
	modelResults["manufacturer_certified"] = s.getManufacturerCertifiedModel(allEngineers, equipment, certMap, minLevel)
	
	// Model 3: Skills Match
	modelResults["skills_match"] = s.getSkillsMatchModel(allEngineers, equipment, minLevel)
	
	// Model 4: Low Workload
	modelResults["low_workload"] = s.getLowWorkloadModel(allEngineers, workloadMap, minLevel)
	
	// Model 5: High Seniority
	modelResults["high_seniority"] = s.getHighSeniorityModel(allEngineers, equipment)
	
	// Model 6: Geographic Proximity (if location data available)
	// modelResults["geographic_proximity"] = s.getGeographicModel(allEngineers, equipment)
	
	// 8. Get tier information
	tierInfo := s.getTierInformation(ctx, allEngineers)
	
	response := &AssignmentSuggestionsResponse{
		TicketID:     ticketID,
		Equipment:    equipmentCtx,
		Ticket:       ticketCtx,
		ModelResults: modelResults,
		TierInfo:     tierInfo,
	}
	
	s.logger.Info("Generated multi-model suggestions",
		slog.String("ticket_id", ticketID),
		slog.Int("total_engineers", len(allEngineers)),
		slog.Int("models_count", len(modelResults)))
	
	return response, nil
}

// getMinLevelFromPriority maps ticket priority to minimum engineer level
func (s *MultiModelAssignmentService) getMinLevelFromPriority(priority domain.TicketPriority) int {
	switch priority {
	case domain.PriorityCritical:
		return 3
	case domain.PriorityHigh:
		return 2
	default:
		return 1
	}
}

// Helper methods for each model
// (Continued in next message due to length...)
