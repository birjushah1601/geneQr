// Package parts provides AI-powered parts recommendations
package parts

import (
	"time"
)

// RecommendationRequest represents a request for parts recommendations
type RecommendationRequest struct {
	// TicketID for context
	TicketID int64

	// Equipment information
	EquipmentID     *int64 // Specific equipment instance
	EquipmentTypeID int64
	EquipmentType   string
	VariantID       *int64 // Installation variant
	VariantName     *string
	Manufacturer    *string
	ModelNumber     *string

	// Problem context
	ProblemType        string
	ProblemDescription string
	Severity           string

	// Diagnosis context (if available)
	DiagnosisID         *string
	DiagnosisConfidence *float64
	IdentifiedIssues    []string
	RecommendedActions  []string

	// Equipment history
	LastMaintenanceDate *time.Time
	OperatingHours      *int
	TotalCycles         *int

	// Request options
	Options RecommendationOptions
}

// RecommendationOptions controls recommendation behavior
type RecommendationOptions struct {
	// IncludeReplacementParts whether to recommend replacement parts
	IncludeReplacementParts bool

	// IncludeAccessories whether to show accessories for upselling
	IncludeAccessories bool

	// IncludePreventiveParts parts due for replacement based on intervals
	IncludePreventiveParts bool

	// CheckInventory whether to check stock availability
	CheckInventory bool

	// IncludePricing whether to include pricing information
	IncludePricing bool

	// UseAI whether to use AI for recommendations
	UseAI bool

	// MaxRecommendations maximum parts to return
	MaxRecommendations int

	// MinConfidence minimum confidence threshold (0-100)
	MinConfidence float64
}

// RecommendationResponse represents parts recommendations
type RecommendationResponse struct {
	// RequestID unique identifier
	RequestID string

	// TicketID reference
	TicketID int64

	// ReplacementParts parts likely needed for repair
	ReplacementParts []PartRecommendation

	// Accessories upsell opportunities
	Accessories []AccessoryRecommendation

	// PreventiveParts parts due for scheduled replacement
	PreventiveParts []PartRecommendation

	// Metadata about recommendations
	Metadata RecommendationMetadata

	// CreatedAt timestamp
	CreatedAt time.Time
}

// PartRecommendation represents a single part recommendation
type PartRecommendation struct {
	// Rank in recommendation list
	Rank int

	// Part details
	PartID          int64
	PartNumber      string
	PartName        string
	Description     string
	Category        string
	Subcategory     *string

	// Recommendation details
	Confidence       float64 // 0-100
	ReasonCode       string  // DiagnosisMatch, HistoricalPattern, PreventiveMaintenance, etc.
	ReasonText       string
	Evidence         []string
	IsOEMPart        bool
	IsCriticalPart   bool

	// Quantity
	RecommendedQuantity int
	QuantityReasoning   string

	// Availability
	StockStatus      string // InStock, LowStock, OutOfStock, OrderRequired
	QuantityAvailable int
	LeadTimeDays      *int

	// Pricing (if requested)
	UnitPrice    *float64
	TotalPrice   *float64
	Currency     string
	SupplierInfo *SupplierInfo

	// Installation info
	RequiresSpecialist      bool
	EstimatedInstallTime    *time.Duration
	InstallationNotes       *string
	CompatibilityNotes      *string

	// Alternative parts
	AlternativeParts []AlternativePart
}

// AccessoryRecommendation represents an accessory upsell
type AccessoryRecommendation struct {
	// Rank in accessory list
	Rank int

	// Part details
	PartID      int64
	PartNumber  string
	PartName    string
	Description string
	Category    string

	// Upsell details
	UpsellPriority      int
	IsRecommended       bool
	IsRequiredForVariant bool
	ReasonText          string
	Benefits            []string

	// Pricing
	UnitPrice            float64
	BundleDiscountPercent *float64
	DiscountedPrice      *float64
	Currency             string

	// Availability
	StockStatus       string
	QuantityAvailable int

	// Marketing
	MarketingDescription *string
	ImageURL             *string
}

// AlternativePart represents an alternative/compatible part
type AlternativePart struct {
	PartID          int64
	PartNumber      string
	PartName        string
	IsUniversalPart bool
	PriceDifference float64 // Positive = more expensive, Negative = cheaper
	LeadTimeDays    *int
	StockStatus     string
}

// SupplierInfo contains supplier information
type SupplierInfo struct {
	SupplierID   int64
	SupplierName string
	IsOEMSupplier bool
	IsPreferred  bool
	LeadTimeDays int
	InStock      bool
}

// RecommendationMetadata contains metadata about the recommendation
type RecommendationMetadata struct {
	// Totals
	TotalReplacementParts int
	TotalAccessories      int
	TotalPreventiveParts  int

	// Estimated costs
	EstimatedPartsCost       *float64
	EstimatedAccessoriesCost *float64
	EstimatedTotalCost       *float64
	Currency                 string

	// AI usage
	UsedAI     bool
	AIProvider string
	AIModel    string
	CostUSD    float64

	// Processing
	ProcessingTime time.Duration
	Version        string
}

// RecommendationHistory represents saved recommendations
type RecommendationHistory struct {
	ID                  int64
	RequestID           string
	TicketID            int64
	ReplacementParts    []PartRecommendation
	Accessories         []AccessoryRecommendation
	PreventiveParts     []PartRecommendation
	Metadata            RecommendationMetadata
	PartsUsed           []int64 // Parts actually used
	AccessoriesSold     []int64 // Accessories sold
	WasAccurate         *bool
	AccuracyFeedback    *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// DefaultRecommendationOptions returns default options
func DefaultRecommendationOptions() RecommendationOptions {
	return RecommendationOptions{
		IncludeReplacementParts: true,
		IncludeAccessories:      true,
		IncludePreventiveParts:  true,
		CheckInventory:          true,
		IncludePricing:          true,
		UseAI:                   true,
		MaxRecommendations:      10,
		MinConfidence:           40.0,
	}
}

// RepairOnlyOptions returns options for repair parts only
func RepairOnlyOptions() RecommendationOptions {
	return RecommendationOptions{
		IncludeReplacementParts: true,
		IncludeAccessories:      false,
		IncludePreventiveParts:  false,
		CheckInventory:          true,
		IncludePricing:          true,
		UseAI:                   true,
		MaxRecommendations:      5,
		MinConfidence:           50.0,
	}
}

// UpsellFocusedOptions returns options focused on accessories
func UpsellFocusedOptions() RecommendationOptions {
	return RecommendationOptions{
		IncludeReplacementParts: true,
		IncludeAccessories:      true,
		IncludePreventiveParts:  false,
		CheckInventory:          true,
		IncludePricing:          true,
		UseAI:                   true,
		MaxRecommendations:      15,
		MinConfidence:           30.0, // Lower threshold for accessories
	}
}

// Constants
const (
	// Reason codes
	ReasonDiagnosisMatch       = "DiagnosisMatch"
	ReasonHistoricalPattern    = "HistoricalPattern"
	ReasonPreventiveMaintenance = "PreventiveMaintenance"
	ReasonCommonFailure        = "CommonFailure"
	ReasonCriticalPart         = "CriticalPart"
	ReasonManufacturerRecommended = "ManufacturerRecommended"
	ReasonVariantRequired      = "VariantRequired"
	ReasonUpsellOpportunity    = "UpsellOpportunity"

	// Stock statuses
	StockInStock      = "InStock"
	StockLowStock     = "LowStock"
	StockOutOfStock   = "OutOfStock"
	StockOrderRequired = "OrderRequired"
	StockUnknown      = "Unknown"
)

