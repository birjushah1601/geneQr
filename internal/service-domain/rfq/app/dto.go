package app

import (
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/rfq/domain"
)

// CreateRFQRequest represents the request to create a new RFQ
type CreateRFQRequest struct {
	Title            string                `json:"title" validate:"required,min=5,max=255"`
	Description      string                `json:"description" validate:"required,min=10"`
	Priority         domain.RFQPriority    `json:"priority" validate:"required,oneof=low medium high critical"`
	ResponseDeadline time.Time             `json:"response_deadline" validate:"required"`
	DeliveryTerms    domain.DeliveryTerms  `json:"delivery_terms"`
	PaymentTerms     domain.PaymentTerms   `json:"payment_terms"`
	InternalNotes    string                `json:"internal_notes"`
	Items            []AddItemRequest      `json:"items"`
}

// UpdateRFQRequest represents the request to update an existing RFQ
type UpdateRFQRequest struct {
	Title            string               `json:"title" validate:"required,min=5,max=255"`
	Description      string               `json:"description" validate:"required,min=10"`
	Priority         domain.RFQPriority   `json:"priority" validate:"required,oneof=low medium high critical"`
	ResponseDeadline time.Time            `json:"response_deadline" validate:"required"`
	DeliveryTerms    domain.DeliveryTerms `json:"delivery_terms"`
	PaymentTerms     domain.PaymentTerms  `json:"payment_terms"`
	InternalNotes    string               `json:"internal_notes"`
}

// AddItemRequest represents the request to add an item to an RFQ
type AddItemRequest struct {
	EquipmentID    *string                `json:"equipment_id"`
	CategoryID     *string                `json:"category_id"`
	Name           string                 `json:"name" validate:"required,min=3,max=255"`
	Description    string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Quantity       int                    `json:"quantity" validate:"required,min=1"`
	Unit           string                 `json:"unit" validate:"required"`
	EstimatedPrice *float64               `json:"estimated_price"`
	Notes          string                 `json:"notes"`
}

// ListRFQsRequest represents the request to list RFQs with filtering
type ListRFQsRequest struct {
	Status        []domain.RFQStatus   `json:"status"`
	Priority      []domain.RFQPriority `json:"priority"`
	CreatedBy     string               `json:"created_by"`
	SearchQuery   string               `json:"search_query"`
	Page          int                  `json:"page"`
	PageSize      int                  `json:"page_size"`
	SortBy        string               `json:"sort_by"`
	SortDirection string               `json:"sort_direction"`
}

// RFQDTO represents an RFQ in the API response
type RFQDTO struct {
	ID               string                 `json:"id"`
	RFQNumber        string                 `json:"rfq_number"`
	TenantID         string                 `json:"tenant_id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Priority         string                 `json:"priority"`
	Status           string                 `json:"status"`
	DeliveryTerms    map[string]interface{} `json:"delivery_terms"`
	PaymentTerms     map[string]interface{} `json:"payment_terms"`
	PublishedAt      *time.Time             `json:"published_at,omitempty"`
	ResponseDeadline time.Time              `json:"response_deadline"`
	ClosedAt         *time.Time             `json:"closed_at,omitempty"`
	CreatedBy        string                 `json:"created_by"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	InternalNotes    string                 `json:"internal_notes,omitempty"`
	Items            []RFQItemDTO           `json:"items"`
	Invitations      []RFQInvitationDTO     `json:"invitations"`
}

// RFQItemDTO represents an RFQ item in the API response
type RFQItemDTO struct {
	ID             string                 `json:"id"`
	RFQID          string                 `json:"rfq_id"`
	EquipmentID    *string                `json:"equipment_id,omitempty"`
	CategoryID     *string                `json:"category_id,omitempty"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Quantity       int                    `json:"quantity"`
	Unit           string                 `json:"unit"`
	EstimatedPrice *float64               `json:"estimated_price,omitempty"`
	Notes          string                 `json:"notes"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// RFQInvitationDTO represents an RFQ invitation in the API response
type RFQInvitationDTO struct {
	ID          string     `json:"id"`
	RFQID       string     `json:"rfq_id"`
	SupplierID  string     `json:"supplier_id"`
	Status      string     `json:"status"`
	InvitedAt   time.Time  `json:"invited_at"`
	ViewedAt    *time.Time `json:"viewed_at,omitempty"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
	Message     string     `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated list response
type PaginatedResponse struct {
	Items      []*RFQDTO `json:"items"`
	TotalItems int       `json:"total_items"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
