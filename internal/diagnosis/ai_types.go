package diagnosis

import "time"

// AISuggestionMetadata provides metadata about AI suggestion confidence and decision
type AISuggestionMetadata struct {
	Provider           string   `json:"provider"`            // openai/anthropic
	Model              string   `json:"model"`
	Confidence         float64  `json:"confidence"`
	ConfidenceFactors  []string `json:"confidence_factors"`  // What influenced confidence
	AlternativesCount  int      `json:"alternatives_count"`
	RequiresFeedback   bool     `json:"requires_feedback"`   // If confidence < threshold
	SuggestionOnly     bool     `json:"suggestion_only"`     // Always true for now
}

// EnhancedDiagnosisResponse extends DiagnosisResponse with AI-assisted features
type EnhancedDiagnosisResponse struct {
	// Original diagnosis fields
	DiagnosisID        string    `json:"diagnosis_id"`
	TicketID           int64     `json:"ticket_id"`
	PrimaryDiagnosis   DiagnosisResult `json:"primary_diagnosis"`
	AlternateDiagnoses []DiagnosisResult `json:"alternate_diagnoses"`
	
	// AI Confidence and Decision Tracking (NEW)
	Confidence      float64  `json:"confidence"`        // 0.0-1.0
	ConfidenceLevel string   `json:"confidence_level"`  // HIGH/MEDIUM/LOW
	DecisionStatus  string   `json:"decision_status"`   // pending/accepted/rejected
	DecidedBy       *int64   `json:"decided_by,omitempty"`       // user_id who made decision
	DecidedAt       *time.Time `json:"decided_at,omitempty"`    // when decision was made
	FeedbackText    string   `json:"feedback_text,omitempty"`    // user feedback on decision

	// AI Metadata for suggestion context (NEW)
	AIMetadata AISuggestionMetadata `json:"ai_metadata"`

	// Vision and context analysis
	VisionAnalysis     *VisionAnalysisResult `json:"vision_analysis,omitempty"`
	ContextUsed        DiagnosisContext `json:"context_used"`

	// Recommended actions and parts
	RecommendedActions []RecommendedAction `json:"recommended_actions"`
	RequiredParts      []RequiredPart `json:"required_parts"`

	// Timing and metadata
	EstimatedResolutionTime *time.Duration `json:"estimated_resolution_time,omitempty"`
	Metadata               DiagnosisMetadata `json:"metadata"`
	CreatedAt              time.Time `json:"created_at"`
}

// AIDecisionFeedback represents user feedback on AI diagnosis
type AIDecisionFeedback struct {
	DiagnosisID  string `json:"diagnosis_id"`
	Decision     string `json:"decision"`      // "accepted" or "rejected"
	UserID       int64  `json:"user_id"`
	UserRole     string `json:"user_role"`     // engineer, dispatcher, admin
	FeedbackText string `json:"feedback_text,omitempty"`
	Corrections  map[string]interface{} `json:"corrections,omitempty"`
}

// ConfidenceCalculation holds factors used to calculate confidence
type ConfidenceCalculation struct {
	VisionConfidence      *float64 `json:"vision_confidence,omitempty"`
	HistoricalMatch       *float64 `json:"historical_match,omitempty"`
	SymptomClarity        *float64 `json:"symptom_clarity,omitempty"`
	ModelConfidence       *float64 `json:"model_confidence,omitempty"`
	FinalConfidence       float64  `json:"final_confidence"`
	ConfidenceFactorCount int      `json:"confidence_factor_count"`
}