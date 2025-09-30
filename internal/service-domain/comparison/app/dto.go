package app

import "github.com/aby-med/medical-platform/internal/service-domain/comparison/domain"

// CreateComparisonRequest represents the request to create a comparison
type CreateComparisonRequest struct {
	RFQID       string   `json:"rfq_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	QuoteIDs    []string `json:"quote_ids"`
}

// UpdateComparisonRequest represents the request to update a comparison
type UpdateComparisonRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}

// UpdateScoringCriteriaRequest represents the request to update scoring weights
type UpdateScoringCriteriaRequest struct {
	PriceWeight      float64 `json:"price_weight"`
	QualityWeight    float64 `json:"quality_weight"`
	DeliveryWeight   float64 `json:"delivery_weight"`
	ComplianceWeight float64 `json:"compliance_weight"`
}

// AddQuoteRequest represents the request to add a quote to a comparison
type AddQuoteRequest struct {
	QuoteID string `json:"quote_id"`
}

// ComparisonResponse represents a comparison in API responses
type ComparisonResponse struct {
	*domain.Comparison
}

// ToResponse converts a domain comparison to a response DTO
func ToResponse(comparison *domain.Comparison) *ComparisonResponse {
	return &ComparisonResponse{
		Comparison: comparison,
	}
}

// ToListResponse converts a list of comparisons to response DTOs
func ToListResponse(comparisons []*domain.Comparison) []*ComparisonResponse {
	responses := make([]*ComparisonResponse, len(comparisons))
	for i, c := range comparisons {
		responses[i] = ToResponse(c)
	}
	return responses
}

// Quote represents quote data needed for scoring (simplified)
type Quote struct {
	ID                string      `json:"id"`
	QuoteNumber       string      `json:"quote_number"`
	SupplierID        string      `json:"supplier_id"`
	SupplierName      string      `json:"supplier_name"`
	RFQID             string      `json:"rfq_id"`
	TotalAmount       float64     `json:"total_amount"`
	ValidUntil        string      `json:"valid_until"`
	DeliveryTerms     string      `json:"delivery_terms"`
	PaymentTerms      string      `json:"payment_terms"`
	WarrantyTerms     string      `json:"warranty_terms"`
	Items             []QuoteItem `json:"items"`
}

// QuoteItem represents an item in a quote
type QuoteItem struct {
	ID                string  `json:"id"`
	RFQItemID         string  `json:"rfq_item_id"`
	EquipmentID       string  `json:"equipment_id"`
	EquipmentName     string  `json:"equipment_name"`
	Quantity          int     `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TotalPrice        float64 `json:"total_price"`
	TaxAmount         float64 `json:"tax_amount"`
	DeliveryTimeframe string  `json:"delivery_timeframe"`
	ManufacturerName  string  `json:"manufacturer_name"`
	ModelNumber       string  `json:"model_number"`
	Specifications    string  `json:"specifications"`
	ComplianceCerts   string  `json:"compliance_certs"`
}
