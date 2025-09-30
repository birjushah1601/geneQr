package domain

import (
	"errors"
	"time"
)

var (
	ErrComparisonNotFound        = errors.New("comparison not found")
	ErrNoQuotesToCompare         = errors.New("at least two quotes required for comparison")
	ErrInvalidWeights            = errors.New("scoring weights must sum to 100")
	ErrInvalidCriteria           = errors.New("invalid scoring criteria")
	ErrQuoteNotInComparison      = errors.New("quote not included in this comparison")
)

// ComparisonStatus represents the status of a comparison
type ComparisonStatus string

const (
	ComparisonStatusDraft      ComparisonStatus = "draft"
	ComparisonStatusActive     ComparisonStatus = "active"
	ComparisonStatusCompleted  ComparisonStatus = "completed"
	ComparisonStatusArchived   ComparisonStatus = "archived"
)

// ScoringCriteria defines the criteria used for scoring quotes
type ScoringCriteria struct {
	PriceWeight       float64 `json:"price_weight"`       // 0-100
	QualityWeight     float64 `json:"quality_weight"`     // 0-100
	DeliveryWeight    float64 `json:"delivery_weight"`    // 0-100
	ComplianceWeight  float64 `json:"compliance_weight"`  // 0-100
}

// ValidateWeights ensures weights sum to 100
func (sc *ScoringCriteria) ValidateWeights() error {
	total := sc.PriceWeight + sc.QualityWeight + sc.DeliveryWeight + sc.ComplianceWeight
	if total < 99.9 || total > 100.1 { // Allow small floating point errors
		return ErrInvalidWeights
	}
	return nil
}

// QuoteScore represents the calculated score for a quote
type QuoteScore struct {
	QuoteID           string    `json:"quote_id"`
	QuoteNumber       string    `json:"quote_number"`
	SupplierID        string    `json:"supplier_id"`
	SupplierName      string    `json:"supplier_name"`
	TotalAmount       float64   `json:"total_amount"`
	PriceScore        float64   `json:"price_score"`          // 0-100
	QualityScore      float64   `json:"quality_score"`        // 0-100
	DeliveryScore     float64   `json:"delivery_score"`       // 0-100
	ComplianceScore   float64   `json:"compliance_score"`     // 0-100
	OverallScore      float64   `json:"overall_score"`        // Weighted average
	Rank              int       `json:"rank"`                 // 1 = best
	Strengths         []string  `json:"strengths"`
	Weaknesses        []string  `json:"weaknesses"`
	Recommendation    string    `json:"recommendation"`
	CalculatedAt      time.Time `json:"calculated_at"`
}

// PriceDifference represents price comparison between quotes
type PriceDifference struct {
	QuoteID         string  `json:"quote_id"`
	QuoteNumber     string  `json:"quote_number"`
	TotalAmount     float64 `json:"total_amount"`
	DifferenceFromLowest  float64 `json:"difference_from_lowest"`   // Absolute difference
	PercentageFromLowest  float64 `json:"percentage_from_lowest"`   // Percentage difference
}

// ItemComparison represents comparison of specific items across quotes
type ItemComparison struct {
	EquipmentID     string                 `json:"equipment_id"`
	EquipmentName   string                 `json:"equipment_name"`
	Quotes          map[string]ItemDetails `json:"quotes"` // quote_id -> details
}

// ItemDetails contains item-specific details for comparison
type ItemDetails struct {
	Quantity          int     `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TotalPrice        float64 `json:"total_price"`
	DeliveryTimeframe string  `json:"delivery_timeframe"`
	ManufacturerName  string  `json:"manufacturer_name"`
	ModelNumber       string  `json:"model_number"`
	Specifications    string  `json:"specifications"`
	ComplianceCerts   string  `json:"compliance_certs"`
}

// Comparison is the aggregate root for quote comparisons
type Comparison struct {
	ID                string              `json:"id"`
	TenantID          string              `json:"tenant_id"`
	RFQID             string              `json:"rfq_id"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	QuoteIDs          []string            `json:"quote_ids"`
	Status            ComparisonStatus    `json:"status"`
	ScoringCriteria   ScoringCriteria     `json:"scoring_criteria"`
	QuoteScores       []QuoteScore        `json:"quote_scores"`
	PriceDifferences  []PriceDifference   `json:"price_differences"`
	ItemComparisons   []ItemComparison    `json:"item_comparisons"`
	BestOverallQuote  string              `json:"best_overall_quote"`
	BestPriceQuote    string              `json:"best_price_quote"`
	Recommendation    string              `json:"recommendation"`
	Notes             string              `json:"notes"`
	CreatedBy         string              `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	CompletedAt       *time.Time          `json:"completed_at,omitempty"`
}

// NewComparison creates a new comparison
func NewComparison(tenantID, rfqID, title, createdBy string, quoteIDs []string) (*Comparison, error) {
	if len(quoteIDs) < 2 {
		return nil, ErrNoQuotesToCompare
	}

	now := time.Now()
	return &Comparison{
		TenantID:    tenantID,
		RFQID:       rfqID,
		Title:       title,
		QuoteIDs:    quoteIDs,
		Status:      ComparisonStatusDraft,
		ScoringCriteria: ScoringCriteria{
			PriceWeight:      40.0,
			QualityWeight:    30.0,
			DeliveryWeight:   20.0,
			ComplianceWeight: 10.0,
		},
		QuoteScores:      []QuoteScore{},
		PriceDifferences: []PriceDifference{},
		ItemComparisons:  []ItemComparison{},
		CreatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// UpdateScoringCriteria updates the weights for scoring
func (c *Comparison) UpdateScoringCriteria(criteria ScoringCriteria) error {
	if err := criteria.ValidateWeights(); err != nil {
		return err
	}
	c.ScoringCriteria = criteria
	c.UpdatedAt = time.Now()
	return nil
}

// AddQuote adds a quote to the comparison
func (c *Comparison) AddQuote(quoteID string) error {
	if c.Status == ComparisonStatusCompleted || c.Status == ComparisonStatusArchived {
		return errors.New("cannot add quotes to completed or archived comparison")
	}
	
	// Check if already added
	for _, id := range c.QuoteIDs {
		if id == quoteID {
			return errors.New("quote already in comparison")
		}
	}
	
	c.QuoteIDs = append(c.QuoteIDs, quoteID)
	c.UpdatedAt = time.Now()
	return nil
}

// RemoveQuote removes a quote from the comparison
func (c *Comparison) RemoveQuote(quoteID string) error {
	if c.Status == ComparisonStatusCompleted || c.Status == ComparisonStatusArchived {
		return errors.New("cannot remove quotes from completed or archived comparison")
	}
	
	newQuoteIDs := []string{}
	found := false
	for _, id := range c.QuoteIDs {
		if id != quoteID {
			newQuoteIDs = append(newQuoteIDs, id)
		} else {
			found = true
		}
	}
	
	if !found {
		return ErrQuoteNotInComparison
	}
	
	if len(newQuoteIDs) < 2 {
		return ErrNoQuotesToCompare
	}
	
	c.QuoteIDs = newQuoteIDs
	c.UpdatedAt = time.Now()
	return nil
}

// Activate marks the comparison as active
func (c *Comparison) Activate() error {
	if c.Status != ComparisonStatusDraft {
		return errors.New("can only activate draft comparisons")
	}
	if len(c.QuoteIDs) < 2 {
		return ErrNoQuotesToCompare
	}
	c.Status = ComparisonStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

// Complete marks the comparison as completed
func (c *Comparison) Complete() error {
	if c.Status != ComparisonStatusActive {
		return errors.New("can only complete active comparisons")
	}
	now := time.Now()
	c.Status = ComparisonStatusCompleted
	c.CompletedAt = &now
	c.UpdatedAt = now
	return nil
}

// Archive archives the comparison
func (c *Comparison) Archive() error {
	if c.Status == ComparisonStatusArchived {
		return errors.New("comparison already archived")
	}
	c.Status = ComparisonStatusArchived
	c.UpdatedAt = time.Now()
	return nil
}

// SetScores sets the calculated scores for quotes
func (c *Comparison) SetScores(scores []QuoteScore) {
	c.QuoteScores = scores
	c.UpdatedAt = time.Now()
	
	// Find best overall and best price
	if len(scores) > 0 {
		bestOverall := scores[0]
		bestPrice := scores[0]
		
		for _, score := range scores {
			if score.OverallScore > bestOverall.OverallScore {
				bestOverall = score
			}
			if score.TotalAmount < bestPrice.TotalAmount {
				bestPrice = score
			}
		}
		
		c.BestOverallQuote = bestOverall.QuoteID
		c.BestPriceQuote = bestPrice.QuoteID
	}
}

// SetPriceDifferences sets the price comparison data
func (c *Comparison) SetPriceDifferences(differences []PriceDifference) {
	c.PriceDifferences = differences
	c.UpdatedAt = time.Now()
}

// SetItemComparisons sets the item-level comparison data
func (c *Comparison) SetItemComparisons(comparisons []ItemComparison) {
	c.ItemComparisons = comparisons
	c.UpdatedAt = time.Now()
}

// SetRecommendation sets the overall recommendation
func (c *Comparison) SetRecommendation(recommendation string) {
	c.Recommendation = recommendation
	c.UpdatedAt = time.Now()
}

// GetQuoteScore retrieves the score for a specific quote
func (c *Comparison) GetQuoteScore(quoteID string) (*QuoteScore, error) {
	for _, score := range c.QuoteScores {
		if score.QuoteID == quoteID {
			return &score, nil
		}
	}
	return nil, errors.New("score not found for quote")
}

// IsQuoteIncluded checks if a quote is included in the comparison
func (c *Comparison) IsQuoteIncluded(quoteID string) bool {
	for _, id := range c.QuoteIDs {
		if id == quoteID {
			return true
		}
	}
	return false
}
