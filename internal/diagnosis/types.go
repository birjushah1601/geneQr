// Package diagnosis provides AI-powered ticket diagnosis capabilities
package diagnosis

import (
	"time"
)

// DiagnosisRequest represents a request for ticket diagnosis
type DiagnosisRequest struct {
	// TicketID to diagnose
	TicketID int64

	// Description of the issue
	Description string

	// EquipmentType (e.g., "CT Scanner", "MRI Machine")
	EquipmentType string

	// EquipmentID if known
	EquipmentID *int64

	// Manufacturer if known
	Manufacturer *string

	// ModelNumber if known
	ModelNumber *string

	// Location context
	Location string

	// LocationType (ICU, Ward, Emergency, etc.)
	LocationType *string

	// Attachments (images, videos)
	Attachments []Attachment

	// ReportedBy user information
	ReportedBy UserContext

	// Priority level
	Priority string

	// Additional context
	AdditionalContext map[string]interface{}

	// Options for diagnosis
	Options DiagnosisOptions
}

// Attachment represents a file attachment (image, video, etc.)
type Attachment struct {
	ID          int64
	FileName    string
	ContentType string
	Size        int64
	URL         string
	Base64Data  string // For direct base64 encoding
}

// UserContext provides information about the reporting user
type UserContext struct {
	UserID   int64
	Username string
	Role     string
}

// DiagnosisOptions controls diagnosis behavior
type DiagnosisOptions struct {
	// IncludeVisionAnalysis whether to analyze images
	IncludeVisionAnalysis bool

	// IncludeHistoricalContext whether to include equipment history
	IncludeHistoricalContext bool

	// IncludeSimilarTickets whether to find similar past tickets
	IncludeSimilarTickets bool

	// MaxSimilarTickets maximum number of similar tickets to retrieve
	MaxSimilarTickets int

	// MinConfidenceThreshold minimum confidence to return (0-100)
	MinConfidenceThreshold float64

	// Model to use (defaults to configured model)
	Model string

	// Temperature for AI (0-1)
	Temperature *float32

	// MaxTokens for response
	MaxTokens *int
}

// DiagnosisResponse represents the AI diagnosis result
type DiagnosisResponse struct {
	// DiagnosisID unique identifier
	DiagnosisID string

	// TicketID being diagnosed
	TicketID int64

	// PrimaryDiagnosis the main identified problem
	PrimaryDiagnosis DiagnosisResult

	// AlternateDiagnoses other possible diagnoses
	AlternateDiagnoses []DiagnosisResult

	// VisionAnalysis if images were analyzed
	VisionAnalysis *VisionAnalysisResult

	// ContextUsed what context was used
	ContextUsed DiagnosisContext

	// RecommendedActions suggested next steps
	RecommendedActions []RecommendedAction

	// RequiredParts parts likely needed
	RequiredParts []RequiredPart

	// EstimatedResolutionTime estimated time to fix
	EstimatedResolutionTime *time.Duration

	// Metadata additional information
	Metadata DiagnosisMetadata

	// CreatedAt timestamp
	CreatedAt time.Time
}

// DiagnosisResult represents a single diagnosis
type DiagnosisResult struct {
	// ProblemCategory (Hardware, Software, User Error, etc.)
	ProblemCategory string

	// ProblemType specific problem classification
	ProblemType string

	// Description detailed description
	Description string

	// Confidence score (0-100)
	Confidence float64

	// Severity (Low, Medium, High, Critical)
	Severity string

	// RootCause identified root cause
	RootCause string

	// Symptoms observed symptoms
	Symptoms []string

	// PossibleCauses list of possible causes
	PossibleCauses []string

	// ReasoningExplanation why AI thinks this is the issue
	ReasoningExplanation string
}

// VisionAnalysisResult represents image/video analysis
type VisionAnalysisResult struct {
	// AttachmentsAnalyzed number of attachments processed
	AttachmentsAnalyzed int

	// Findings visual findings
	Findings []VisualFinding

	// OverallAssessment summary of visual analysis
	OverallAssessment string

	// DetectedComponents equipment components visible
	DetectedComponents []string

	// VisibleDamage any visible damage
	VisibleDamage []DamageDescription

	// Confidence overall confidence in vision analysis
	Confidence float64
}

// VisualFinding represents a finding from image analysis
type VisualFinding struct {
	// AttachmentID which attachment
	AttachmentID int64

	// Finding description
	Finding string

	// Confidence in this finding
	Confidence float64

	// Category (Damage, Normal, Warning, Error Display, etc.)
	Category string

	// Location in image if applicable
	Location string
}

// DamageDescription describes visible damage
type DamageDescription struct {
	// Type of damage (Physical, Display Error, LED indicators, etc.)
	Type string

	// Description
	Description string

	// Severity (Minor, Moderate, Severe)
	Severity string

	// Location where damage is visible
	Location string
}

// DiagnosisContext represents context used in diagnosis
type DiagnosisContext struct {
	// EquipmentHistoryUsed whether equipment history was used
	EquipmentHistoryUsed bool

	// EquipmentHistoryCount number of historical records
	EquipmentHistoryCount int

	// SimilarTicketsUsed whether similar tickets were used
	SimilarTicketsUsed bool

	// SimilarTicketsCount number of similar tickets found
	SimilarTicketsCount int

	// SimilarTickets details of similar tickets
	SimilarTickets []SimilarTicketInfo

	// ManufacturerGuidelinesUsed whether manufacturer docs were referenced
	ManufacturerGuidelinesUsed bool

	// KnownIssuesUsed whether known issues DB was checked
	KnownIssuesUsed bool
}

// SimilarTicketInfo represents a similar past ticket
type SimilarTicketInfo struct {
	TicketID    int64
	Description string
	Resolution  string
	TimeTaken   time.Duration
	Similarity  float64 // 0-1 similarity score
}

// RecommendedAction represents a suggested action
type RecommendedAction struct {
	// Order priority order (1 = first)
	Order int

	// Action description
	Action string

	// ActionType (Diagnostic, Repair, Replace, Escalate, etc.)
	ActionType string

	// EstimatedTime to complete
	EstimatedTime *time.Duration

	// RequiresSpecialist whether specialist is needed
	RequiresSpecialist bool

	// SpecialistType if specialist needed
	SpecialistType *string

	// RequiredTools tools needed
	RequiredTools []string

	// RequiredParts parts needed
	RequiredParts []string

	// SafetyPrecautions safety notes
	SafetyPrecautions []string
}

// RequiredPart represents a part likely needed
type RequiredPart struct {
	// PartCode internal part code
	PartCode string

	// PartName human-readable name
	PartName string

	// PartCategory category
	PartCategory string

	// Probability likelihood of needing this part (0-100)
	Probability float64

	// Quantity estimated quantity needed
	Quantity int

	// IsOEMRequired whether OEM part is required
	IsOEMRequired bool

	// Manufacturer part manufacturer
	Manufacturer string

	// EstimatedCost estimated cost
	EstimatedCost *float64

	// TypicalLeadTime typical delivery time
	TypicalLeadTime *time.Duration

	// AlternativeParts alternative part codes
	AlternativeParts []string
}

// DiagnosisMetadata contains metadata about the diagnosis process
type DiagnosisMetadata struct {
	// Provider AI provider used (openai, anthropic)
	Provider string

	// Model AI model used
	Model string

	// TokensUsed total tokens
	TokensUsed int

	// CostUSD cost in USD
	CostUSD float64

	// Latency time taken
	Latency time.Duration

	// VisionAnalysisPerformed whether vision was used
	VisionAnalysisPerformed bool

	// ContextEnrichmentPerformed whether context was enriched
	ContextEnrichmentPerformed bool

	// Version diagnosis engine version
	Version string
}

// DiagnosisHistory represents a saved diagnosis
type DiagnosisHistory struct {
	ID              int64
	DiagnosisID     string
	TicketID        int64
	DiagnosisData   DiagnosisResponse
	WasAccurate     *bool   // Feedback: was this accurate?
	AccuracyScore   *int    // 0-100 if provided
	FeedbackNotes   *string // Human feedback
	ActualResolution *string // What actually fixed it
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProblemCategory constants
const (
	CategoryHardware      = "Hardware"
	CategorySoftware      = "Software"
	CategoryConfiguration = "Configuration"
	CategoryUserError     = "User Error"
	CategoryNetwork       = "Network"
	CategoryPower         = "Power"
	CategoryEnvironmental = "Environmental"
	CategoryUnknown       = "Unknown"
)

// Severity constants
const (
	SeverityLow      = "Low"
	SeverityMedium   = "Medium"
	SeverityHigh     = "High"
	SeverityCritical = "Critical"
)

// ActionType constants
const (
	ActionTypeDiagnostic = "Diagnostic"
	ActionTypeRepair     = "Repair"
	ActionTypeReplace    = "Replace"
	ActionTypeCalibrate  = "Calibrate"
	ActionTypeClean      = "Clean"
	ActionTypeUpdate     = "Update"
	ActionTypeRestart    = "Restart"
	ActionTypeEscalate   = "Escalate"
	ActionTypeMonitor    = "Monitor"
)

// DefaultDiagnosisOptions returns default diagnosis options
func DefaultDiagnosisOptions() DiagnosisOptions {
	return DiagnosisOptions{
		IncludeVisionAnalysis:    true,
		IncludeHistoricalContext: true,
		IncludeSimilarTickets:    true,
		MaxSimilarTickets:        5,
		MinConfidenceThreshold:   20.0, // Show diagnoses with 20%+ confidence
	}
}

