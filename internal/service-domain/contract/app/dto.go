package app

import (
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/contract/domain"
)

// CreateContractRequest represents a request to create a new contract
type CreateContractRequest struct {
	RFQID              string                     `json:"rfq_id"`
	QuoteID            string                     `json:"quote_id"`
	SupplierID         string                     `json:"supplier_id"`
	SupplierName       string                     `json:"supplier_name"`
	StartDate          time.Time                  `json:"start_date"`
	EndDate            time.Time                  `json:"end_date"`
	PaymentTerms       string                     `json:"payment_terms"`
	DeliveryTerms      string                     `json:"delivery_terms"`
	WarrantyTerms      string                     `json:"warranty_terms"`
	TermsAndConditions string                     `json:"terms_and_conditions"`
	Items              []ContractItemRequest      `json:"items"`
	PaymentSchedule    []PaymentTermRequest       `json:"payment_schedule,omitempty"`
	DeliverySchedule   []DeliveryScheduleRequest  `json:"delivery_schedule,omitempty"`
	TaxAmount          float64                    `json:"tax_amount"`
	Notes              string                     `json:"notes,omitempty"`
}

// ContractItemRequest represents a contract item in requests
type ContractItemRequest struct {
	EquipmentID      string  `json:"equipment_id"`
	EquipmentName    string  `json:"equipment_name"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
	ManufacturerName string  `json:"manufacturer_name"`
	ModelNumber      string  `json:"model_number"`
	Specifications   string  `json:"specifications"`
	WarrantyPeriod   string  `json:"warranty_period"`
}

// PaymentTermRequest represents a payment term in requests
type PaymentTermRequest struct {
	DueDate     time.Time `json:"due_date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

// DeliveryScheduleRequest represents a delivery milestone in requests
type DeliveryScheduleRequest struct {
	MilestoneDate time.Time `json:"milestone_date"`
	Description   string    `json:"description"`
}

// UpdateContractRequest represents a request to update a contract
type UpdateContractRequest struct {
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	PaymentTerms       *string    `json:"payment_terms,omitempty"`
	DeliveryTerms      *string    `json:"delivery_terms,omitempty"`
	WarrantyTerms      *string    `json:"warranty_terms,omitempty"`
	TermsAndConditions *string    `json:"terms_and_conditions,omitempty"`
	TaxAmount          *float64   `json:"tax_amount,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
}

// SignContractRequest represents a request to sign a contract
type SignContractRequest struct {
	SignedBy string `json:"signed_by"`
}

// CancelContractRequest represents a request to cancel a contract
type CancelContractRequest struct {
	Reason string `json:"reason"`
}

// SuspendContractRequest represents a request to suspend a contract
type SuspendContractRequest struct {
	Reason string `json:"reason"`
}

// AddAmendmentRequest represents a request to add an amendment
type AddAmendmentRequest struct {
	Description string `json:"description"`
	Changes     string `json:"changes"`
	AmendedBy   string `json:"amended_by"`
}

// MarkPaymentPaidRequest represents a request to mark a payment as paid
type MarkPaymentPaidRequest struct {
	PaymentIndex int `json:"payment_index"`
}

// MarkDeliveryCompletedRequest represents a request to mark a delivery as completed
type MarkDeliveryCompletedRequest struct {
	DeliveryIndex int `json:"delivery_index"`
}

// ContractResponse represents a contract in responses
type ContractResponse struct {
	ID                 string                       `json:"id"`
	TenantID           string                       `json:"tenant_id"`
	ContractNumber     string                       `json:"contract_number"`
	RFQID              string                       `json:"rfq_id"`
	QuoteID            string                       `json:"quote_id"`
	SupplierID         string                       `json:"supplier_id"`
	SupplierName       string                       `json:"supplier_name"`
	Status             string                       `json:"status"`
	TotalAmount        float64                      `json:"total_amount"`
	Currency           string                       `json:"currency"`
	TaxAmount          float64                      `json:"tax_amount"`
	StartDate          time.Time                    `json:"start_date"`
	EndDate            time.Time                    `json:"end_date"`
	SignedDate         *time.Time                   `json:"signed_date,omitempty"`
	PaymentTerms       string                       `json:"payment_terms"`
	DeliveryTerms      string                       `json:"delivery_terms"`
	WarrantyTerms      string                       `json:"warranty_terms"`
	TermsAndConditions string                       `json:"terms_and_conditions"`
	PaymentSchedule    []domain.PaymentTerm         `json:"payment_schedule"`
	DeliverySchedule   []domain.DeliverySchedule    `json:"delivery_schedule"`
	Items              []domain.ContractItem        `json:"items"`
	Amendments         []domain.Amendment           `json:"amendments"`
	Notes              string                       `json:"notes"`
	CreatedBy          string                       `json:"created_by"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	SignedBy           string                       `json:"signed_by,omitempty"`
	PaymentProgress    float64                      `json:"payment_progress"`
	DeliveryProgress   float64                      `json:"delivery_progress"`
}

// ListContractsRequest represents filtering criteria for listing contracts
type ListContractsRequest struct {
	RFQID        string   `json:"rfq_id,omitempty"`
	SupplierID   string   `json:"supplier_id,omitempty"`
	Status       []string `json:"status,omitempty"`
	CreatedBy    string   `json:"created_by,omitempty"`
	SortBy       string   `json:"sort_by,omitempty"`
	SortDirection string  `json:"sort_direction,omitempty"`
	Page         int      `json:"page,omitempty"`
	PageSize     int      `json:"page_size,omitempty"`
}

// ListContractsResponse represents a paginated list of contracts
type ListContractsResponse struct {
	Contracts  []ContractResponse `json:"contracts"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// ToContractResponse converts a domain contract to a response DTO
func ToContractResponse(contract *domain.Contract) ContractResponse {
	return ContractResponse{
		ID:                 contract.ID,
		TenantID:           contract.TenantID,
		ContractNumber:     contract.ContractNumber,
		RFQID:              contract.RFQID,
		QuoteID:            contract.QuoteID,
		SupplierID:         contract.SupplierID,
		SupplierName:       contract.SupplierName,
		Status:             string(contract.Status),
		TotalAmount:        contract.TotalAmount,
		Currency:           contract.Currency,
		TaxAmount:          contract.TaxAmount,
		StartDate:          contract.StartDate,
		EndDate:            contract.EndDate,
		SignedDate:         contract.SignedDate,
		PaymentTerms:       contract.PaymentTerms,
		DeliveryTerms:      contract.DeliveryTerms,
		WarrantyTerms:      contract.WarrantyTerms,
		TermsAndConditions: contract.TermsAndConditions,
		PaymentSchedule:    contract.PaymentSchedule,
		DeliverySchedule:   contract.DeliverySchedule,
		Items:              contract.Items,
		Amendments:         contract.Amendments,
		Notes:              contract.Notes,
		CreatedBy:          contract.CreatedBy,
		CreatedAt:          contract.CreatedAt,
		UpdatedAt:          contract.UpdatedAt,
		SignedBy:           contract.SignedBy,
		PaymentProgress:    contract.GetPaymentProgress(),
		DeliveryProgress:   contract.GetDeliveryProgress(),
	}
}
