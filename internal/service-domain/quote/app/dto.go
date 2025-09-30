package app

import (
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/quote/domain"
)

// CreateQuoteRequest represents the request to create a new quote
type CreateQuoteRequest struct {
	RFQID         string             `json:"rfq_id" validate:"required"`
	SupplierID    string             `json:"supplier_id" validate:"required"`
	ValidUntil    time.Time          `json:"valid_until" validate:"required"`
	DeliveryTerms string             `json:"delivery_terms"`
	PaymentTerms  string             `json:"payment_terms"`
	WarrantyTerms string             `json:"warranty_terms"`
	Notes         string             `json:"notes"`
	Items         []QuoteItemRequest `json:"items" validate:"required,min=1"`
}

// QuoteItemRequest represents a line item in a quote request
type QuoteItemRequest struct {
	RFQItemID         string  `json:"rfq_item_id" validate:"required"`
	EquipmentID       string  `json:"equipment_id" validate:"required"`
	EquipmentName     string  `json:"equipment_name" validate:"required"`
	Quantity          int     `json:"quantity" validate:"required,min=1"`
	UnitPrice         float64 `json:"unit_price" validate:"required,min=0"`
	TaxRate           float64 `json:"tax_rate"`
	DeliveryTimeframe string  `json:"delivery_timeframe"`
	ManufacturerName  string  `json:"manufacturer_name"`
	ModelNumber       string  `json:"model_number"`
	Specifications    string  `json:"specifications"`
	ComplianceCerts   string  `json:"compliance_certs"`
	Notes             string  `json:"notes"`
}

// UpdateQuoteRequest represents the request to update quote details
type UpdateQuoteRequest struct {
	DeliveryTerms string  `json:"delivery_terms"`
	PaymentTerms  string  `json:"payment_terms"`
	WarrantyTerms string  `json:"warranty_terms"`
	Notes         string  `json:"notes"`
	ValidUntil    *time.Time `json:"valid_until"`
}

// AcceptQuoteRequest represents the request to accept a quote
type AcceptQuoteRequest struct {
	ReviewedBy string `json:"reviewed_by" validate:"required"`
	Notes      string `json:"notes"`
}

// RejectQuoteRequest represents the request to reject a quote
type RejectQuoteRequest struct {
	ReviewedBy string `json:"reviewed_by" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
}

// ReviseQuoteRequest represents the request to revise a quote
type ReviseQuoteRequest struct {
	Changes   string `json:"changes" validate:"required"`
	RevisedBy string `json:"revised_by" validate:"required"`
}

// QuoteResponse represents the response DTO for a quote
type QuoteResponse struct {
	ID              string                   `json:"id"`
	TenantID        string                   `json:"tenant_id"`
	RFQID           string                   `json:"rfq_id"`
	SupplierID      string                   `json:"supplier_id"`
	QuoteNumber     string                   `json:"quote_number"`
	Status          string                   `json:"status"`
	TotalAmount     float64                  `json:"total_amount"`
	Currency        string                   `json:"currency"`
	ValidUntil      time.Time                `json:"valid_until"`
	DeliveryTerms   string                   `json:"delivery_terms,omitempty"`
	PaymentTerms    string                   `json:"payment_terms,omitempty"`
	WarrantyTerms   string                   `json:"warranty_terms,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	Items           []QuoteItemResponse      `json:"items"`
	RevisionNumber  int                      `json:"revision_number"`
	Revisions       []QuoteRevisionResponse  `json:"revisions,omitempty"`
	ReviewedAt      *time.Time               `json:"reviewed_at,omitempty"`
	ReviewedBy      *string                  `json:"reviewed_by,omitempty"`
	ReviewNotes     *string                  `json:"review_notes,omitempty"`
	RejectionReason *string                  `json:"rejection_reason,omitempty"`
	Metadata        map[string]interface{}   `json:"metadata,omitempty"`
	CreatedBy       string                   `json:"created_by"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// QuoteItemResponse represents a line item in a quote response
type QuoteItemResponse struct {
	ID                string  `json:"id"`
	RFQItemID         string  `json:"rfq_item_id"`
	EquipmentID       string  `json:"equipment_id"`
	EquipmentName     string  `json:"equipment_name"`
	Quantity          int     `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TotalPrice        float64 `json:"total_price"`
	TaxRate           float64 `json:"tax_rate"`
	TaxAmount         float64 `json:"tax_amount"`
	DeliveryTimeframe string  `json:"delivery_timeframe,omitempty"`
	ManufacturerName  string  `json:"manufacturer_name,omitempty"`
	ModelNumber       string  `json:"model_number,omitempty"`
	Specifications    string  `json:"specifications,omitempty"`
	ComplianceCerts   string  `json:"compliance_certs,omitempty"`
	Notes             string  `json:"notes,omitempty"`
}

// QuoteRevisionResponse represents a quote revision in the response
type QuoteRevisionResponse struct {
	RevisionNumber int                    `json:"revision_number"`
	RevisedAt      time.Time              `json:"revised_at"`
	RevisedBy      string                 `json:"revised_by"`
	Changes        string                 `json:"changes"`
	PreviousTotal  float64                `json:"previous_total"`
	NewTotal       float64                `json:"new_total"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ListQuotesResponse represents the response for listing quotes
type ListQuotesResponse struct {
	Quotes []QuoteResponse `json:"quotes"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
}

// ToQuoteResponse converts a domain Quote to QuoteResponse
func ToQuoteResponse(quote *domain.Quote) QuoteResponse {
	items := make([]QuoteItemResponse, len(quote.Items))
	for i, item := range quote.Items {
		items[i] = QuoteItemResponse{
			ID:                item.ID,
			RFQItemID:         item.RFQItemID,
			EquipmentID:       item.EquipmentID,
			EquipmentName:     item.EquipmentName,
			Quantity:          item.Quantity,
			UnitPrice:         item.UnitPrice,
			TotalPrice:        item.TotalPrice,
			TaxRate:           item.TaxRate,
			TaxAmount:         item.TaxAmount,
			DeliveryTimeframe: item.DeliveryTimeframe,
			ManufacturerName:  item.ManufacturerName,
			ModelNumber:       item.ModelNumber,
			Specifications:    item.Specifications,
			ComplianceCerts:   item.ComplianceCerts,
			Notes:             item.Notes,
		}
	}

	revisions := make([]QuoteRevisionResponse, len(quote.Revisions))
	for i, rev := range quote.Revisions {
		revisions[i] = QuoteRevisionResponse{
			RevisionNumber: rev.RevisionNumber,
			RevisedAt:      rev.RevisedAt,
			RevisedBy:      rev.RevisedBy,
			Changes:        rev.Changes,
			PreviousTotal:  rev.PreviousTotal,
			NewTotal:       rev.NewTotal,
			Metadata:       rev.Metadata,
		}
	}

	return QuoteResponse{
		ID:              quote.ID,
		TenantID:        quote.TenantID,
		RFQID:           quote.RFQID,
		SupplierID:      quote.SupplierID,
		QuoteNumber:     quote.QuoteNumber,
		Status:          string(quote.Status),
		TotalAmount:     quote.TotalAmount,
		Currency:        quote.Currency,
		ValidUntil:      quote.ValidUntil,
		DeliveryTerms:   quote.DeliveryTerms,
		PaymentTerms:    quote.PaymentTerms,
		WarrantyTerms:   quote.WarrantyTerms,
		Notes:           quote.Notes,
		Items:           items,
		RevisionNumber:  quote.RevisionNumber,
		Revisions:       revisions,
		ReviewedAt:      quote.ReviewedAt,
		ReviewedBy:      quote.ReviewedBy,
		ReviewNotes:     quote.ReviewNotes,
		RejectionReason: quote.RejectionReason,
		Metadata:        quote.Metadata,
		CreatedBy:       quote.CreatedBy,
		CreatedAt:       quote.CreatedAt,
		UpdatedAt:       quote.UpdatedAt,
	}
}

// ToQuoteResponses converts a list of domain Quotes to QuoteResponse list
func ToQuoteResponses(quotes []*domain.Quote) []QuoteResponse {
	responses := make([]QuoteResponse, len(quotes))
	for i, quote := range quotes {
		responses[i] = ToQuoteResponse(quote)
	}
	return responses
}
