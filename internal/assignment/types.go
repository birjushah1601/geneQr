// Package assignment provides intelligent engineer assignment recommendations
package assignment

import (
	"time"
)

// AssignmentRequest represents a request for engineer recommendations
type AssignmentRequest struct {
	// TicketID to assign
	TicketID int64

	// Equipment information
	EquipmentTypeID int64
	EquipmentType   string
	Manufacturer    *string
	ModelNumber     *string

	// Ticket details
	Priority     string
	ProblemType  string
	Description  string
	LocationID   int64
	LocationName string

	// Requirements
	RequiredSkills      []string
	RequiresSpecialist  bool
	SpecialistType      *string
	EstimatedComplexity string // Low, Medium, High

	// Diagnosis insights (if available)
	DiagnosisConfidence *float64
	DiagnosisSeverity   *string
	RequiredParts       []string

	// Constraints
	MaxRecommendations int
	MinScoreThreshold  float64

	// Options
	Options AssignmentOptions
}

// AssignmentOptions controls assignment behavior
type AssignmentOptions struct {
	// ConsiderWorkload whether to factor in current workload
	ConsiderWorkload bool

	// ConsiderLocation whether to prioritize nearby engineers
	ConsiderLocation bool

	// ConsiderPerformance whether to factor in past performance
	ConsiderPerformance bool

	// ConsiderSpecialization whether to match skills
	ConsiderSpecialization bool

	// ConsiderAvailability whether to check availability
	ConsiderAvailability bool

	// UseAI whether to use AI for final ranking
	UseAI bool

	// Weights for scoring factors (0-1, should sum to ~1.0)
	Weights ScoringWeights
}

// ScoringWeights defines weights for each scoring factor
type ScoringWeights struct {
	Expertise    float64 // 0-1
	Location     float64 // 0-1
	Performance  float64 // 0-1
	Workload     float64 // 0-1
	Availability float64 // 0-1
}

// AssignmentResponse represents engineer recommendations
type AssignmentResponse struct {
	// RequestID unique identifier
	RequestID string

	// TicketID being assigned
	TicketID int64

	// Recommendations ranked list of engineers
	Recommendations []EngineerRecommendation

	// Metadata about the recommendation process
	Metadata AssignmentMetadata

	// CreatedAt timestamp
	CreatedAt time.Time
}

// EngineerRecommendation represents a single engineer recommendation
type EngineerRecommendation struct {
	// Rank position in recommendation list (1 = best match)
	Rank int

	// Engineer information
	EngineerID   int64
	EngineerName string
	Email        string
	Phone        *string

	// Overall score (0-100)
	OverallScore float64

	// Score breakdown
	Scores ScoreBreakdown

	// Match reasons
	MatchReasons []MatchReason

	// Warnings if any
	Warnings []string

	// Current status
	CurrentWorkload  int     // Number of open tickets
	AvailabilityInfo *string // e.g., "Available now", "Available in 2 hours"

	// Location info
	CurrentLocation *string
	DistanceKM      *float64

	// Skills matching
	MatchingSkills []string
	MissingSkills  []string

	// Performance history
	AverageResolutionTime *time.Duration
	SuccessRate           *float64
	RecentTicketsCount    int
}

// ScoreBreakdown shows individual scoring components
type ScoreBreakdown struct {
	// Individual scores (0-100)
	ExpertiseScore    float64
	LocationScore     float64
	PerformanceScore  float64
	WorkloadScore     float64
	AvailabilityScore float64

	// Weights applied
	Weights ScoringWeights

	// Raw weighted scores
	WeightedExpertise    float64
	WeightedLocation     float64
	WeightedPerformance  float64
	WeightedWorkload     float64
	WeightedAvailability float64
}

// MatchReason explains why an engineer was recommended
type MatchReason struct {
	// Category of the reason
	Category string // Expertise, Location, Performance, etc.

	// Reason description
	Reason string

	// Impact on score (High, Medium, Low)
	Impact string

	// Supporting evidence
	Evidence string
}

// AssignmentMetadata contains metadata about recommendations
type AssignmentMetadata struct {
	// EngineersEvaluated total engineers considered
	EngineersEvaluated int

	// EngineersFiltered engineers after filtering
	EngineersFiltered int

	// UsedAI whether AI was used for ranking
	UsedAI bool

	// AIProvider if AI was used
	AIProvider string

	// AIModel if AI was used
	AIModel string

	// CostUSD if AI was used
	CostUSD float64

	// ProcessingTime time taken
	ProcessingTime time.Duration

	// Version assignment engine version
	Version string
}

// EngineerProfile represents detailed engineer information for scoring
type EngineerProfile struct {
	EngineerID   int64
	FullName     string
	Email        string
	Phone        *string
	IsActive     bool
	HireDate     time.Time
	Department   string
	Specialization *string

	// Skills and expertise
	Skills              []EngineerSkill
	Certifications      []string
	EquipmentExpertise  []EquipmentExpertise

	// Performance metrics
	TotalTicketsResolved int
	SuccessRate          float64
	AverageResolutionTime time.Duration
	AverageRating        *float64

	// Current status
	CurrentWorkload   int
	OpenTicketsCount  int
	AvailabilityStatus string // Available, Busy, OnLeave, etc.
	CurrentLocation    *string

	// Recent activity
	LastTicketDate     *time.Time
	RecentTickets      []RecentTicketSummary
}

// EngineerSkill represents a skill with proficiency
type EngineerSkill struct {
	SkillName    string
	SkillType    string // Technical, Soft, Equipment-Specific
	ProficiencyLevel string // Beginner, Intermediate, Advanced, Expert
	YearsExperience  *int
}

// EquipmentExpertise represents expertise with equipment types
type EquipmentExpertise struct {
	EquipmentType       string
	ManufacturerCode    *string
	TicketsHandled      int
	SuccessRate         float64
	AverageResolutionTime time.Duration
}

// RecentTicketSummary summarizes recent ticket activity
type RecentTicketSummary struct {
	TicketID         int64
	EquipmentType    string
	ProblemType      string
	ResolutionTime   time.Duration
	WasSuccessful    bool
	CustomerRating   *int
}

// AssignmentHistory represents a saved assignment recommendation
type AssignmentHistory struct {
	ID              int64
	RequestID       string
	TicketID        int64
	Recommendations []EngineerRecommendation
	SelectedEngineerID *int64
	SelectionReason    *string
	WasSuccessful      *bool
	ActualResolutionTime *time.Duration
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// DefaultAssignmentOptions returns default assignment options
func DefaultAssignmentOptions() AssignmentOptions {
	return AssignmentOptions{
		ConsiderWorkload:       true,
		ConsiderLocation:       true,
		ConsiderPerformance:    true,
		ConsiderSpecialization: true,
		ConsiderAvailability:   true,
		UseAI:                  true,
		Weights: DefaultScoringWeights(),
	}
}

// DefaultScoringWeights returns balanced scoring weights
func DefaultScoringWeights() ScoringWeights {
	return ScoringWeights{
		Expertise:    0.35, // 35% - Most important
		Location:     0.15, // 15% - Important for response time
		Performance:  0.25, // 25% - Historical success
		Workload:     0.15, // 15% - Current capacity
		Availability: 0.10, // 10% - Immediate availability
	}
}

// UrgentScoringWeights returns weights optimized for urgent tickets
func UrgentScoringWeights() ScoringWeights {
	return ScoringWeights{
		Expertise:    0.30,
		Location:     0.25, // Higher for urgent
		Performance:  0.20,
		Workload:     0.10,
		Availability: 0.15, // Higher for urgent
	}
}

// QualityScoringWeights returns weights optimized for quality
func QualityScoringWeights() ScoringWeights {
	return ScoringWeights{
		Expertise:    0.40, // Highest
		Location:     0.10,
		Performance:  0.35, // High
		Workload:     0.10,
		Availability: 0.05,
	}
}

// Constants
const (
	// Availability statuses
	AvailabilityAvailable   = "Available"
	AvailabilityBusy        = "Busy"
	AvailabilityOnLeave     = "OnLeave"
	AvailabilityOffline     = "Offline"

	// Proficiency levels
	ProficiencyBeginner     = "Beginner"
	ProficiencyIntermediate = "Intermediate"
	ProficiencyAdvanced     = "Advanced"
	ProficiencyExpert       = "Expert"

	// Impact levels
	ImpactHigh   = "High"
	ImpactMedium = "Medium"
	ImpactLow    = "Low"

	// Match categories
	CategoryExpertise    = "Expertise"
	CategoryLocation     = "Location"
	CategoryPerformance  = "Performance"
	CategoryWorkload     = "Workload"
	CategoryAvailability = "Availability"
)

