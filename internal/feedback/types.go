package feedback

import "time"

// FeedbackSource indicates where feedback originated
type FeedbackSource string

const (
	SourceHuman  FeedbackSource = "human"  // Explicit feedback from users
	SourceMachine FeedbackSource = "machine" // Implicit feedback from system outcomes
)

// FeedbackType categorizes the feedback
type FeedbackType string

const (
	// Diagnosis feedback
	FeedbackDiagnosisAccuracy  FeedbackType = "diagnosis_accuracy"
	FeedbackDiagnosisCorrection FeedbackType = "diagnosis_correction"
	
	// Assignment feedback
	FeedbackAssignmentAcceptance FeedbackType = "assignment_acceptance"
	FeedbackAssignmentPerformance FeedbackType = "assignment_performance"
	FeedbackAssignmentCorrection  FeedbackType = "assignment_correction"
	
	// Parts feedback
	FeedbackPartsAccuracy   FeedbackType = "parts_accuracy"
	FeedbackPartsUsed       FeedbackType = "parts_used"
	FeedbackPartsCorrection FeedbackType = "parts_correction"
	
	// General feedback
	FeedbackGeneral FeedbackType = "general"
)

// FeedbackSentiment represents the overall sentiment
type FeedbackSentiment string

const (
	SentimentPositive FeedbackSentiment = "positive"
	SentimentNeutral  FeedbackSentiment = "neutral"
	SentimentNegative FeedbackSentiment = "negative"
)

// FeedbackEntry represents a single piece of feedback
type FeedbackEntry struct {
	FeedbackID   int64             `json:"feedback_id"`
	Source       FeedbackSource    `json:"source"`        // human or machine
	Type         FeedbackType      `json:"type"`          // what aspect
	
	// Context
	TicketID     *int64            `json:"ticket_id,omitempty"`
	RequestID    *string           `json:"request_id,omitempty"` // diagnosis/assignment/parts request ID
	ServiceType  string            `json:"service_type"`          // diagnosis, assignment, parts
	
	// Human feedback
	UserID       *int64            `json:"user_id,omitempty"`
	Rating       *int              `json:"rating,omitempty"`      // 1-5 scale
	Sentiment    FeedbackSentiment `json:"sentiment"`
	Comments     string            `json:"comments,omitempty"`
	
	// Machine feedback (outcomes)
	Outcomes     map[string]interface{} `json:"outcomes,omitempty"`
	
	// Corrections (what should have been)
	Corrections  map[string]interface{} `json:"corrections,omitempty"`
	
	// Metadata
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	
	// Timestamps
	CreatedAt    time.Time         `json:"created_at"`
	ProcessedAt  *time.Time        `json:"processed_at,omitempty"`
}

// HumanFeedbackRequest represents explicit human feedback
type HumanFeedbackRequest struct {
	// Required
	ServiceType string            `json:"service_type"`  // diagnosis, assignment, parts
	RequestID   string            `json:"request_id"`    // ID of the AI request
	
	// Optional context
	TicketID    *int64            `json:"ticket_id,omitempty"`
	
	// Feedback content
	Rating      *int              `json:"rating,omitempty"`     // 1-5 scale
	WasAccurate bool              `json:"was_accurate"`         // Binary accuracy
	Comments    string            `json:"comments,omitempty"`
	
	// Corrections (what should have been recommended)
	Corrections map[string]interface{} `json:"corrections,omitempty"`
	
	// Metadata
	UserID      int64             `json:"user_id"`
	UserRole    string            `json:"user_role,omitempty"` // engineer, dispatcher, etc.
}

// MachineFeedbackRequest represents implicit system-generated feedback
type MachineFeedbackRequest struct {
	// Required
	ServiceType string            `json:"service_type"`
	RequestID   string            `json:"request_id"`
	TicketID    int64             `json:"ticket_id"`
	
	// Outcomes (what actually happened)
	Outcomes    MachineFeedbackOutcomes `json:"outcomes"`
	
	// Metadata
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MachineFeedbackOutcomes contains actual outcomes
type MachineFeedbackOutcomes struct {
	// Diagnosis outcomes
	ActualProblem           *string   `json:"actual_problem,omitempty"`
	ActualRootCause         *string   `json:"actual_root_cause,omitempty"`
	DiagnosisMatchedAI      *bool     `json:"diagnosis_matched_ai,omitempty"`
	
	// Assignment outcomes
	AssignmentAccepted      *bool     `json:"assignment_accepted,omitempty"`
	AssignedEngineerID      *int64    `json:"assigned_engineer_id,omitempty"`
	WasTopRecommendation    *bool     `json:"was_top_recommendation,omitempty"`
	AssignmentRank          *int      `json:"assignment_rank,omitempty"`
	
	// Parts outcomes
	PartsUsedIDs            []int64   `json:"parts_used_ids,omitempty"`
	RecommendedPartsUsed    []int64   `json:"recommended_parts_used,omitempty"`
	UnrecommendedPartsUsed  []int64   `json:"unrecommended_parts_used,omitempty"`
	AccessoriesSold         []int64   `json:"accessories_sold,omitempty"`
	
	// Resolution outcomes
	ResolutionTime          *int      `json:"resolution_time_minutes,omitempty"`
	FirstTimeFixRate        *bool     `json:"first_time_fix,omitempty"`
	CustomerSatisfaction    *int      `json:"customer_satisfaction,omitempty"` // 1-5
	
	// Cost outcomes
	ActualCost              *float64  `json:"actual_cost,omitempty"`
	EstimatedCost           *float64  `json:"estimated_cost,omitempty"`
	CostVariancePercent     *float64  `json:"cost_variance_percent,omitempty"`
}

// FeedbackAnalysis represents analyzed feedback patterns
type FeedbackAnalysis struct {
	ServiceType string            `json:"service_type"`
	Period      string            `json:"period"` // daily, weekly, monthly
	
	// Volume metrics
	TotalFeedback       int     `json:"total_feedback"`
	HumanFeedbackCount  int     `json:"human_feedback_count"`
	MachineFeedbackCount int    `json:"machine_feedback_count"`
	
	// Sentiment metrics
	PositiveFeedback    int     `json:"positive_feedback"`
	NeutralFeedback     int     `json:"neutral_feedback"`
	NegativeFeedback    int     `json:"negative_feedback"`
	AvgRating           float64 `json:"avg_rating,omitempty"`
	
	// Accuracy metrics
	AccuracyRate        float64 `json:"accuracy_rate"`
	
	// Common issues
	CommonIssues        []FeedbackIssue `json:"common_issues,omitempty"`
	
	// Improvement opportunities
	Improvements        []ImprovementOpportunity `json:"improvements,omitempty"`
	
	// Generated
	GeneratedAt         time.Time `json:"generated_at"`
}

// FeedbackIssue represents a common problem identified from feedback
type FeedbackIssue struct {
	IssueType   string  `json:"issue_type"`
	Description string  `json:"description"`
	Frequency   int     `json:"frequency"`
	Severity    string  `json:"severity"` // high, medium, low
	Examples    []int64 `json:"examples,omitempty"` // feedback IDs
}

// ImprovementOpportunity represents an actionable improvement
type ImprovementOpportunity struct {
	OpportunityID   string                 `json:"opportunity_id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	ServiceType     string                 `json:"service_type"` // diagnosis, assignment, parts
	ImpactLevel     string                 `json:"impact_level"` // high, medium, low
	ImplementationType string              `json:"implementation_type"` // prompt_tuning, weight_adjustment, training_data
	SuggestedChanges map[string]interface{} `json:"suggested_changes"`
	SupportingData   []int64                `json:"supporting_data"` // feedback IDs
	CreatedAt        time.Time              `json:"created_at"`
	Status           string                 `json:"status"` // pending, applied, rejected
}

// LearningAction represents an action taken based on feedback
type LearningAction struct {
	ActionID         string                 `json:"action_id"`
	OpportunityID    string                 `json:"opportunity_id"`
	ActionType       string                 `json:"action_type"` // prompt_update, weight_adjustment, config_change
	ServiceType      string                 `json:"service_type"`
	
	// Changes made
	Changes          map[string]interface{} `json:"changes"`
	
	// Before/after comparison
	BeforeMetrics    map[string]float64     `json:"before_metrics,omitempty"`
	AfterMetrics     map[string]float64     `json:"after_metrics,omitempty"`
	
	// Status
	Status           string                 `json:"status"` // testing, deployed, rolled_back
	AppliedAt        time.Time              `json:"applied_at"`
	AppliedBy        string                 `json:"applied_by"` // system or user ID
	
	// Results
	ResultNotes      string                 `json:"result_notes,omitempty"`
	RolledBackAt     *time.Time             `json:"rolled_back_at,omitempty"`
	RollbackReason   string                 `json:"rollback_reason,omitempty"`
}

// FeedbackMetrics provides high-level metrics
type FeedbackMetrics struct {
	ServiceType string    `json:"service_type"`
	DateRange   DateRange `json:"date_range"`
	
	// Volume
	TotalRequests       int `json:"total_requests"`
	FeedbackReceived    int `json:"feedback_received"`
	FeedbackRate        float64 `json:"feedback_rate"` // percentage
	
	// Quality metrics
	AvgAccuracyRate     float64 `json:"avg_accuracy_rate"`
	AvgRating           float64 `json:"avg_rating"`
	PositiveSentiment   float64 `json:"positive_sentiment_percent"`
	
	// Learning metrics
	ImprovementsFound   int `json:"improvements_found"`
	ActionsApplied      int `json:"actions_applied"`
	MeasuredImpact      float64 `json:"measured_impact_percent,omitempty"`
	
	// Trends
	AccuracyTrend       string `json:"accuracy_trend"` // improving, stable, declining
	SentimentTrend      string `json:"sentiment_trend"`
}

// DateRange represents a time period
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// FeedbackSummary provides a summary for dashboard
type FeedbackSummary struct {
	OverallMetrics      FeedbackMetrics              `json:"overall_metrics"`
	ByServiceType       map[string]FeedbackMetrics   `json:"by_service_type"`
	RecentImprovements  []ImprovementOpportunity     `json:"recent_improvements"`
	ActiveActions       []LearningAction             `json:"active_actions"`
	TopIssues           []FeedbackIssue              `json:"top_issues"`
	GeneratedAt         time.Time                    `json:"generated_at"`
}


