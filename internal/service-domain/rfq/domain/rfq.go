package domain

import (
	"errors"
	"time"
)

// RFQ Status represents the lifecycle state of an RFQ
type RFQStatus string

const (
	RFQStatusDraft     RFQStatus = "draft"
	RFQStatusPublished RFQStatus = "published"
	RFQStatusClosed    RFQStatus = "closed"
	RFQStatusAwarded   RFQStatus = "awarded"
	RFQStatusCancelled RFQStatus = "cancelled"
)

// RFQPriority represents the urgency level
type RFQPriority string

const (
	RFQPriorityLow      RFQPriority = "low"
	RFQPriorityMedium   RFQPriority = "medium"
	RFQPriorityHigh     RFQPriority = "high"
	RFQPriorityCritical RFQPriority = "critical"
)

// DeliveryTerms represents delivery requirements
type DeliveryTerms struct {
	Address         string    `json:"address"`
	City            string    `json:"city"`
	State           string    `json:"state"`
	PostalCode      string    `json:"postal_code"`
	Country         string    `json:"country"`
	RequiredBy      time.Time `json:"required_by"`
	SpecialNotes    string    `json:"special_notes,omitempty"`
	InstallationReq bool      `json:"installation_required"`
}

// PaymentTerms represents payment conditions
type PaymentTerms struct {
	PaymentMethod   string `json:"payment_method"`
	PaymentDays     int    `json:"payment_days"`
	AdvancePayment  int    `json:"advance_payment_percent"`
	SpecialTerms    string `json:"special_terms,omitempty"`
}

// RFQ represents a Request for Quote aggregate root
type RFQ struct {
	// Identity
	ID       string `json:"id"`
	RFQNumber string `json:"rfq_number"`
	TenantID string `json:"tenant_id"`
	
	// Basic Info
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Priority    RFQPriority `json:"priority"`
	Status      RFQStatus   `json:"status"`
	
	// Items
	Items []RFQItem `json:"items"`
	
	// Delivery & Payment
	DeliveryTerms DeliveryTerms `json:"delivery_terms"`
	PaymentTerms  PaymentTerms  `json:"payment_terms"`
	
	// Timeline
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	ResponseDeadline time.Time `json:"response_deadline"`
	ClosedAt       *time.Time `json:"closed_at,omitempty"`
	
	// Metadata
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Invitations
	Invitations []RFQInvitation `json:"invitations,omitempty"`
	
	// Notes
	InternalNotes string `json:"internal_notes,omitempty"`
}

// RFQItem represents an equipment item in the RFQ
type RFQItem struct {
	ID             string                 `json:"id"`
	RFQID          string                 `json:"rfq_id"`
	EquipmentID    *string                `json:"equipment_id,omitempty"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications,omitempty"`
	Quantity       int                    `json:"quantity"`
	Unit           string                 `json:"unit"`
	EstimatedPrice *float64               `json:"estimated_price,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
	CategoryID     *string                `json:"category_id,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// RFQInvitation represents an invitation sent to a supplier
type RFQInvitation struct {
	ID         string    `json:"id"`
	RFQID      string    `json:"rfq_id"`
	SupplierID string    `json:"supplier_id"`
	Status     string    `json:"status"` // invited, viewed, quoted, declined
	InvitedAt  time.Time `json:"invited_at"`
	ViewedAt   *time.Time `json:"viewed_at,omitempty"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
	Message    string    `json:"message,omitempty"`
}

// Domain errors
var (
	ErrRFQNotFound            = errors.New("rfq not found")
	ErrInvalidRFQStatus       = errors.New("invalid rfq status")
	ErrInvalidTransition      = errors.New("invalid status transition")
	ErrRFQAlreadyPublished    = errors.New("rfq already published")
	ErrRFQNotPublished        = errors.New("rfq not published")
	ErrRFQAlreadyClosed       = errors.New("rfq already closed")
	ErrDeadlinePassed         = errors.New("response deadline has passed")
	ErrNoItems                = errors.New("rfq must have at least one item")
	ErrInvalidDeliveryTerms   = errors.New("invalid delivery terms")
	ErrInvalidPaymentTerms    = errors.New("invalid payment terms")
	ErrInvalidQuantity        = errors.New("quantity must be greater than zero")
)

// NewRFQ creates a new RFQ in draft status
func NewRFQ(
	id, rfqNumber, tenantID, title, description string,
	priority RFQPriority,
	responseDeadline time.Time,
	deliveryTerms DeliveryTerms,
	paymentTerms PaymentTerms,
	createdBy string,
) (*RFQ, error) {
	now := time.Now()
	
	// Validate deadline is in the future
	if responseDeadline.Before(now) {
		return nil, errors.New("response deadline must be in the future")
	}
	
	rfq := &RFQ{
		ID:               id,
		RFQNumber:        rfqNumber,
		TenantID:         tenantID,
		Title:            title,
		Description:      description,
		Priority:         priority,
		Status:           RFQStatusDraft,
		Items:            []RFQItem{},
		DeliveryTerms:    deliveryTerms,
		PaymentTerms:     paymentTerms,
		ResponseDeadline: responseDeadline,
		CreatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
		Invitations:      []RFQInvitation{},
	}
	
	return rfq, nil
}

// AddItem adds an equipment item to the RFQ
func (r *RFQ) AddItem(item RFQItem) error {
	if r.Status != RFQStatusDraft {
		return errors.New("can only add items to draft RFQ")
	}
	
	if item.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	item.RFQID = r.ID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	
	r.Items = append(r.Items, item)
	r.UpdatedAt = time.Now()
	
	return nil
}

// RemoveItem removes an item from the RFQ
func (r *RFQ) RemoveItem(itemID string) error {
	if r.Status != RFQStatusDraft {
		return errors.New("can only remove items from draft RFQ")
	}
	
	for i, item := range r.Items {
		if item.ID == itemID {
			r.Items = append(r.Items[:i], r.Items[i+1:]...)
			r.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return errors.New("item not found")
}

// Publish transitions the RFQ from draft to published
func (r *RFQ) Publish() error {
	if r.Status != RFQStatusDraft {
		return ErrRFQAlreadyPublished
	}
	
	if len(r.Items) == 0 {
		return ErrNoItems
	}
	
	// Validate delivery terms
	if r.DeliveryTerms.Address == "" || r.DeliveryTerms.City == "" {
		return ErrInvalidDeliveryTerms
	}
	
	now := time.Now()
	r.Status = RFQStatusPublished
	r.PublishedAt = &now
	r.UpdatedAt = now
	
	return nil
}

// Close transitions the RFQ to closed status
func (r *RFQ) Close() error {
	if r.Status != RFQStatusPublished {
		return ErrRFQNotPublished
	}
	
	now := time.Now()
	r.Status = RFQStatusClosed
	r.ClosedAt = &now
	r.UpdatedAt = now
	
	return nil
}

// Cancel cancels the RFQ
func (r *RFQ) Cancel() error {
	if r.Status == RFQStatusAwarded || r.Status == RFQStatusCancelled {
		return errors.New("cannot cancel RFQ in current status")
	}
	
	r.Status = RFQStatusCancelled
	r.UpdatedAt = time.Now()
	
	return nil
}

// Award marks the RFQ as awarded
func (r *RFQ) Award() error {
	if r.Status != RFQStatusClosed {
		return errors.New("can only award closed RFQ")
	}
	
	r.Status = RFQStatusAwarded
	r.UpdatedAt = time.Now()
	
	return nil
}

// InviteSupplier adds a supplier invitation
func (r *RFQ) InviteSupplier(invitation RFQInvitation) error {
	if r.Status != RFQStatusPublished {
		return errors.New("can only invite suppliers to published RFQ")
	}
	
	// Check if supplier already invited
	for _, inv := range r.Invitations {
		if inv.SupplierID == invitation.SupplierID {
			return errors.New("supplier already invited")
		}
	}
	
	invitation.RFQID = r.ID
	invitation.InvitedAt = time.Now()
	invitation.Status = "invited"
	
	r.Invitations = append(r.Invitations, invitation)
	r.UpdatedAt = time.Now()
	
	return nil
}

// IsExpired checks if the response deadline has passed
func (r *RFQ) IsExpired() bool {
	return time.Now().After(r.ResponseDeadline)
}

// CanBeEdited returns true if the RFQ can be edited
func (r *RFQ) CanBeEdited() bool {
	return r.Status == RFQStatusDraft
}

// CanBePublished returns true if the RFQ can be published
func (r *RFQ) CanBePublished() bool {
	return r.Status == RFQStatusDraft && len(r.Items) > 0
}
