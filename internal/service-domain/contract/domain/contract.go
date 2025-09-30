package domain

import (
	"errors"
	"time"
)

var (
	ErrContractNotFound       = errors.New("contract not found")
	ErrInvalidContractStatus  = errors.New("invalid contract status")
	ErrContractAlreadyActive  = errors.New("contract already active")
	ErrContractAlreadySigned  = errors.New("contract already signed")
	ErrCannotAmendContract    = errors.New("cannot amend contract in current state")
	ErrInvalidPaymentSchedule = errors.New("invalid payment schedule")
)

// ContractStatus represents the status of a contract
type ContractStatus string

const (
	ContractStatusDraft     ContractStatus = "draft"
	ContractStatusActive    ContractStatus = "active"
	ContractStatusCompleted ContractStatus = "completed"
	ContractStatusCancelled ContractStatus = "cancelled"
	ContractStatusExpired   ContractStatus = "expired"
	ContractStatusSuspended ContractStatus = "suspended"
)

// PaymentTerm represents payment terms
type PaymentTerm struct {
	DueDate     time.Time `json:"due_date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Paid        bool      `json:"paid"`
	PaidDate    *time.Time `json:"paid_date,omitempty"`
}

// DeliverySchedule represents delivery milestones
type DeliverySchedule struct {
	MilestoneDate time.Time `json:"milestone_date"`
	Description   string    `json:"description"`
	Completed     bool      `json:"completed"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`
}

// ContractItem represents an item in the contract
type ContractItem struct {
	ID                string  `json:"id"`
	EquipmentID       string  `json:"equipment_id"`
	EquipmentName     string  `json:"equipment_name"`
	Quantity          int     `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TotalPrice        float64 `json:"total_price"`
	ManufacturerName  string  `json:"manufacturer_name"`
	ModelNumber       string  `json:"model_number"`
	Specifications    string  `json:"specifications"`
	WarrantyPeriod    string  `json:"warranty_period"`
}

// Amendment represents a contract amendment
type Amendment struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Changes     string    `json:"changes"` // JSON string of what changed
	AmendedBy   string    `json:"amended_by"`
}

// Contract is the aggregate root for contract management
type Contract struct {
	ID                 string             `json:"id"`
	TenantID           string             `json:"tenant_id"`
	ContractNumber     string             `json:"contract_number"`
	RFQID              string             `json:"rfq_id"`
	QuoteID            string             `json:"quote_id"`
	SupplierID         string             `json:"supplier_id"`
	SupplierName       string             `json:"supplier_name"`
	Status             ContractStatus     `json:"status"`
	
	// Financial
	TotalAmount        float64            `json:"total_amount"`
	Currency           string             `json:"currency"`
	TaxAmount          float64            `json:"tax_amount"`
	
	// Dates
	StartDate          time.Time          `json:"start_date"`
	EndDate            time.Time          `json:"end_date"`
	SignedDate         *time.Time         `json:"signed_date,omitempty"`
	
	// Terms
	PaymentTerms       string             `json:"payment_terms"`
	DeliveryTerms      string             `json:"delivery_terms"`
	WarrantyTerms      string             `json:"warranty_terms"`
	TermsAndConditions string             `json:"terms_and_conditions"`
	
	// Schedules
	PaymentSchedule    []PaymentTerm      `json:"payment_schedule"`
	DeliverySchedule   []DeliverySchedule `json:"delivery_schedule"`
	
	// Items
	Items              []ContractItem     `json:"items"`
	
	// Amendments
	Amendments         []Amendment        `json:"amendments"`
	
	// Metadata
	Notes              string             `json:"notes"`
	CreatedBy          string             `json:"created_by"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	SignedBy           string             `json:"signed_by,omitempty"`
}

// NewContract creates a new contract
func NewContract(tenantID, rfqID, quoteID, supplierID, supplierName, createdBy string) *Contract {
	now := time.Now()
	return &Contract{
		TenantID:         tenantID,
		RFQID:            rfqID,
		QuoteID:          quoteID,
		SupplierID:       supplierID,
		SupplierName:     supplierName,
		Status:           ContractStatusDraft,
		Currency:         "USD",
		PaymentSchedule:  []PaymentTerm{},
		DeliverySchedule: []DeliverySchedule{},
		Items:            []ContractItem{},
		Amendments:       []Amendment{},
		CreatedBy:        createdBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// Activate activates the contract
func (c *Contract) Activate() error {
	if c.Status != ContractStatusDraft {
		return errors.New("can only activate draft contracts")
	}
	if c.ContractNumber == "" {
		return errors.New("contract number required before activation")
	}
	if len(c.Items) == 0 {
		return errors.New("contract must have at least one item")
	}
	c.Status = ContractStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

// Sign marks the contract as signed
func (c *Contract) Sign(signedBy string) error {
	if c.Status != ContractStatusActive {
		return errors.New("can only sign active contracts")
	}
	if c.SignedDate != nil {
		return ErrContractAlreadySigned
	}
	now := time.Now()
	c.SignedDate = &now
	c.SignedBy = signedBy
	c.UpdatedAt = now
	return nil
}

// Complete marks the contract as completed
func (c *Contract) Complete() error {
	if c.Status != ContractStatusActive {
		return errors.New("can only complete active contracts")
	}
	c.Status = ContractStatusCompleted
	c.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the contract
func (c *Contract) Cancel(reason string) error {
	if c.Status == ContractStatusCompleted || c.Status == ContractStatusCancelled {
		return errors.New("cannot cancel completed or already cancelled contracts")
	}
	c.Status = ContractStatusCancelled
	c.Notes = c.Notes + "\nCancellation reason: " + reason
	c.UpdatedAt = time.Now()
	return nil
}

// Suspend suspends the contract
func (c *Contract) Suspend(reason string) error {
	if c.Status != ContractStatusActive {
		return errors.New("can only suspend active contracts")
	}
	c.Status = ContractStatusSuspended
	c.Notes = c.Notes + "\nSuspension reason: " + reason
	c.UpdatedAt = time.Now()
	return nil
}

// Resume resumes a suspended contract
func (c *Contract) Resume() error {
	if c.Status != ContractStatusSuspended {
		return errors.New("can only resume suspended contracts")
	}
	c.Status = ContractStatusActive
	c.UpdatedAt = time.Now()
	return nil
}

// AddAmendment adds an amendment to the contract
func (c *Contract) AddAmendment(amendment Amendment) error {
	if c.Status != ContractStatusActive {
		return ErrCannotAmendContract
	}
	c.Amendments = append(c.Amendments, amendment)
	c.UpdatedAt = time.Now()
	return nil
}

// AddPaymentTerm adds a payment term
func (c *Contract) AddPaymentTerm(term PaymentTerm) error {
	if c.Status == ContractStatusCompleted || c.Status == ContractStatusCancelled {
		return errors.New("cannot add payment terms to completed or cancelled contracts")
	}
	c.PaymentSchedule = append(c.PaymentSchedule, term)
	c.UpdatedAt = time.Now()
	return nil
}

// MarkPaymentPaid marks a payment as paid
func (c *Contract) MarkPaymentPaid(index int) error {
	if index < 0 || index >= len(c.PaymentSchedule) {
		return errors.New("invalid payment index")
	}
	if c.PaymentSchedule[index].Paid {
		return errors.New("payment already marked as paid")
	}
	now := time.Now()
	c.PaymentSchedule[index].Paid = true
	c.PaymentSchedule[index].PaidDate = &now
	c.UpdatedAt = now
	return nil
}

// AddDeliveryMilestone adds a delivery milestone
func (c *Contract) AddDeliveryMilestone(milestone DeliverySchedule) error {
	if c.Status == ContractStatusCompleted || c.Status == ContractStatusCancelled {
		return errors.New("cannot add delivery milestones to completed or cancelled contracts")
	}
	c.DeliverySchedule = append(c.DeliverySchedule, milestone)
	c.UpdatedAt = time.Now()
	return nil
}

// MarkDeliveryCompleted marks a delivery milestone as completed
func (c *Contract) MarkDeliveryCompleted(index int) error {
	if index < 0 || index >= len(c.DeliverySchedule) {
		return errors.New("invalid delivery milestone index")
	}
	if c.DeliverySchedule[index].Completed {
		return errors.New("delivery already marked as completed")
	}
	now := time.Now()
	c.DeliverySchedule[index].Completed = true
	c.DeliverySchedule[index].CompletedDate = &now
	c.UpdatedAt = now
	return nil
}

// AddItem adds an item to the contract
func (c *Contract) AddItem(item ContractItem) error {
	if c.Status != ContractStatusDraft {
		return errors.New("can only add items to draft contracts")
	}
	c.Items = append(c.Items, item)
	c.UpdatedAt = time.Now()
	return nil
}

// CalculateTotals recalculates the contract totals
func (c *Contract) CalculateTotals() {
	total := 0.0
	for _, item := range c.Items {
		total += item.TotalPrice
	}
	c.TotalAmount = total + c.TaxAmount
	c.UpdatedAt = time.Now()
}

// IsExpired checks if the contract has expired
func (c *Contract) IsExpired() bool {
	return time.Now().After(c.EndDate) && c.Status == ContractStatusActive
}

// GetPaymentProgress returns the payment completion percentage
func (c *Contract) GetPaymentProgress() float64 {
	if len(c.PaymentSchedule) == 0 {
		return 0.0
	}
	
	paidCount := 0
	for _, payment := range c.PaymentSchedule {
		if payment.Paid {
			paidCount++
		}
	}
	
	return (float64(paidCount) / float64(len(c.PaymentSchedule))) * 100.0
}

// GetDeliveryProgress returns the delivery completion percentage
func (c *Contract) GetDeliveryProgress() float64 {
	if len(c.DeliverySchedule) == 0 {
		return 0.0
	}
	
	completedCount := 0
	for _, delivery := range c.DeliverySchedule {
		if delivery.Completed {
			completedCount++
		}
	}
	
	return (float64(completedCount) / float64(len(c.DeliverySchedule))) * 100.0
}
