package domain

import (
	"errors"
	"time"

	"github.com/segmentio/ksuid"
)

// QuoteStatus represents the lifecycle status of a quote
type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "draft"     // Being prepared
	QuoteStatusSubmitted QuoteStatus = "submitted" // Submitted to hospital
	QuoteStatusUnderReview QuoteStatus = "under_review" // Hospital reviewing
	QuoteStatusRevised   QuoteStatus = "revised"   // Supplier made revisions
	QuoteStatusAccepted  QuoteStatus = "accepted"  // Hospital accepted
	QuoteStatusRejected  QuoteStatus = "rejected"  // Hospital rejected
	QuoteStatusExpired   QuoteStatus = "expired"   // Validity period expired
	QuoteStatusWithdrawn QuoteStatus = "withdrawn" // Supplier withdrew
)

// Quote aggregate root - represents a supplier's quote for an RFQ
type Quote struct {
	// Identity
	ID         string `json:"id"`
	TenantID   string `json:"tenant_id"`
	RFQID      string `json:"rfq_id"`
	SupplierID string `json:"supplier_id"`

	// Quote details
	QuoteNumber    string      `json:"quote_number"`
	Status         QuoteStatus `json:"status"`
	TotalAmount    float64     `json:"total_amount"`
	Currency       string      `json:"currency"`
	ValidUntil     time.Time   `json:"valid_until"`
	DeliveryTerms  string      `json:"delivery_terms"`
	PaymentTerms   string      `json:"payment_terms"`
	WarrantyTerms  string      `json:"warranty_terms"`
	Notes          string      `json:"notes,omitempty"`

	// Items
	Items []QuoteItem `json:"items"`

	// Revision tracking
	RevisionNumber int             `json:"revision_number"`
	Revisions      []QuoteRevision `json:"revisions,omitempty"`

	// Review information
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy      *string    `json:"reviewed_by,omitempty"`
	ReviewNotes     *string    `json:"review_notes,omitempty"`
	RejectionReason *string    `json:"rejection_reason,omitempty"`

	// Metadata
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy string                 `json:"created_by"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// QuoteItem represents a line item in a quote
type QuoteItem struct {
	ID                 string  `json:"id"`
	RFQItemID          string  `json:"rfq_item_id"`
	EquipmentID        string  `json:"equipment_id"`
	EquipmentName      string  `json:"equipment_name"`
	Quantity           int     `json:"quantity"`
	UnitPrice          float64 `json:"unit_price"`
	TotalPrice         float64 `json:"total_price"`
	TaxRate            float64 `json:"tax_rate,omitempty"`
	TaxAmount          float64 `json:"tax_amount,omitempty"`
	DeliveryTimeframe  string  `json:"delivery_timeframe"`
	ManufacturerName   string  `json:"manufacturer_name,omitempty"`
	ModelNumber        string  `json:"model_number,omitempty"`
	Specifications     string  `json:"specifications,omitempty"`
	ComplianceCerts    string  `json:"compliance_certs,omitempty"`
	Notes              string  `json:"notes,omitempty"`
}

// QuoteRevision tracks the history of quote revisions
type QuoteRevision struct {
	RevisionNumber int                    `json:"revision_number"`
	RevisedAt      time.Time              `json:"revised_at"`
	RevisedBy      string                 `json:"revised_by"`
	Changes        string                 `json:"changes"` // Description of changes
	PreviousTotal  float64                `json:"previous_total"`
	NewTotal       float64                `json:"new_total"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Domain errors
var (
	ErrQuoteNotFound         = errors.New("quote not found")
	ErrQuoteAlreadySubmitted = errors.New("quote already submitted")
	ErrQuoteExpired          = errors.New("quote has expired")
	ErrQuoteNotInDraft       = errors.New("quote is not in draft status")
	ErrQuoteNotSubmitted     = errors.New("quote is not submitted")
	ErrInvalidStatus         = errors.New("invalid quote status")
	ErrNoItems               = errors.New("quote must have at least one item")
)

// NewQuote creates a new quote
func NewQuote(tenantID, rfqID, supplierID, createdBy string, validUntil time.Time) (*Quote, error) {
	now := time.Now()

	// Validate validity period
	if validUntil.Before(now) {
		return nil, errors.New("valid_until must be in the future")
	}

	quote := &Quote{
		ID:             ksuid.New().String(),
		TenantID:       tenantID,
		RFQID:          rfqID,
		SupplierID:     supplierID,
		Status:         QuoteStatusDraft,
		Currency:       "USD", // Default
		ValidUntil:     validUntil,
		RevisionNumber: 1,
		Items:          []QuoteItem{},
		Revisions:      []QuoteRevision{},
		Metadata:       make(map[string]interface{}),
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return quote, nil
}

// AddItem adds a line item to the quote
func (q *Quote) AddItem(item QuoteItem) error {
	if q.Status != QuoteStatusDraft && q.Status != QuoteStatusRevised {
		return ErrQuoteAlreadySubmitted
	}

	// Generate item ID if not set
	if item.ID == "" {
		item.ID = ksuid.New().String()
	}

	// Calculate total price
	item.TotalPrice = float64(item.Quantity) * item.UnitPrice
	item.TaxAmount = item.TotalPrice * item.TaxRate

	q.Items = append(q.Items, item)
	q.recalculateTotal()
	q.UpdatedAt = time.Now()

	return nil
}

// RemoveItem removes a line item from the quote
func (q *Quote) RemoveItem(itemID string) error {
	if q.Status != QuoteStatusDraft && q.Status != QuoteStatusRevised {
		return ErrQuoteAlreadySubmitted
	}

	for i, item := range q.Items {
		if item.ID == itemID {
			q.Items = append(q.Items[:i], q.Items[i+1:]...)
			q.recalculateTotal()
			q.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found")
}

// UpdateItem updates a line item in the quote
func (q *Quote) UpdateItem(itemID string, updatedItem QuoteItem) error {
	if q.Status != QuoteStatusDraft && q.Status != QuoteStatusRevised {
		return ErrQuoteAlreadySubmitted
	}

	for i, item := range q.Items {
		if item.ID == itemID {
			// Preserve ID
			updatedItem.ID = item.ID
			updatedItem.TotalPrice = float64(updatedItem.Quantity) * updatedItem.UnitPrice
			updatedItem.TaxAmount = updatedItem.TotalPrice * updatedItem.TaxRate
			
			q.Items[i] = updatedItem
			q.recalculateTotal()
			q.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("item not found")
}

// Submit submits the quote for review
func (q *Quote) Submit() error {
	if q.Status != QuoteStatusDraft && q.Status != QuoteStatusRevised {
		return errors.New("only draft or revised quotes can be submitted")
	}

	if len(q.Items) == 0 {
		return ErrNoItems
	}

	if time.Now().After(q.ValidUntil) {
		return ErrQuoteExpired
	}

	// Generate quote number if not set
	if q.QuoteNumber == "" {
		q.QuoteNumber = "QT-" + ksuid.New().String()[:10]
	}

	q.Status = QuoteStatusSubmitted
	q.UpdatedAt = time.Now()

	return nil
}

// Revise creates a new revision of the quote
func (q *Quote) Revise(changes string, revisedBy string) error {
	if q.Status != QuoteStatusUnderReview && q.Status != QuoteStatusSubmitted {
		return errors.New("only submitted or under-review quotes can be revised")
	}

	// Save revision history
	revision := QuoteRevision{
		RevisionNumber: q.RevisionNumber,
		RevisedAt:      time.Now(),
		RevisedBy:      revisedBy,
		Changes:        changes,
		PreviousTotal:  q.TotalAmount,
		NewTotal:       q.TotalAmount, // Will be updated as items change
		Metadata:       make(map[string]interface{}),
	}

	q.Revisions = append(q.Revisions, revision)
	q.RevisionNumber++
	q.Status = QuoteStatusRevised
	q.UpdatedAt = time.Now()

	return nil
}

// Accept marks the quote as accepted by the hospital
func (q *Quote) Accept(reviewedBy string, notes string) error {
	if q.Status != QuoteStatusSubmitted && q.Status != QuoteStatusUnderReview {
		return errors.New("only submitted or under-review quotes can be accepted")
	}

	if time.Now().After(q.ValidUntil) {
		return ErrQuoteExpired
	}

	now := time.Now()
	q.Status = QuoteStatusAccepted
	q.ReviewedAt = &now
	q.ReviewedBy = &reviewedBy
	q.ReviewNotes = &notes
	q.UpdatedAt = now

	return nil
}

// Reject marks the quote as rejected by the hospital
func (q *Quote) Reject(reviewedBy string, reason string) error {
	if q.Status != QuoteStatusSubmitted && q.Status != QuoteStatusUnderReview {
		return errors.New("only submitted or under-review quotes can be rejected")
	}

	now := time.Now()
	q.Status = QuoteStatusRejected
	q.ReviewedAt = &now
	q.ReviewedBy = &reviewedBy
	q.RejectionReason = &reason
	q.UpdatedAt = now

	return nil
}

// Withdraw allows supplier to withdraw the quote
func (q *Quote) Withdraw() error {
	if q.Status == QuoteStatusAccepted || q.Status == QuoteStatusRejected {
		return errors.New("cannot withdraw accepted or rejected quotes")
	}

	q.Status = QuoteStatusWithdrawn
	q.UpdatedAt = time.Now()

	return nil
}

// MarkUnderReview marks the quote as under review by the hospital
func (q *Quote) MarkUnderReview() error {
	if q.Status != QuoteStatusSubmitted {
		return errors.New("only submitted quotes can be marked as under review")
	}

	q.Status = QuoteStatusUnderReview
	q.UpdatedAt = time.Now()

	return nil
}

// MarkExpired marks the quote as expired
func (q *Quote) MarkExpired() error {
	if time.Now().Before(q.ValidUntil) {
		return errors.New("quote is not yet expired")
	}

	q.Status = QuoteStatusExpired
	q.UpdatedAt = time.Now()

	return nil
}

// Private helper methods

func (q *Quote) recalculateTotal() {
	total := 0.0
	for _, item := range q.Items {
		total += item.TotalPrice + item.TaxAmount
	}
	q.TotalAmount = total
}

// IsEditable checks if the quote can be edited
func (q *Quote) IsEditable() bool {
	return q.Status == QuoteStatusDraft || q.Status == QuoteStatusRevised
}

// IsExpired checks if the quote has expired
func (q *Quote) IsExpired() bool {
	return time.Now().After(q.ValidUntil)
}
