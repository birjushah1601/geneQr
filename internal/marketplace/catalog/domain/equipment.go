package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

// Common errors
var (
	ErrEquipmentNotFound      = errors.New("equipment not found")
	ErrInvalidEquipmentID     = errors.New("invalid equipment ID")
	ErrInvalidEquipmentName   = errors.New("equipment name cannot be empty")
	ErrInvalidCategory        = errors.New("invalid equipment category")
	ErrInvalidManufacturer    = errors.New("invalid equipment manufacturer")
	ErrInvalidPrice           = errors.New("equipment price must be positive")
	ErrInvalidSpecifications  = errors.New("invalid equipment specifications")
)

// Category represents a medical equipment category
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// ParentID allows for hierarchical categories (optional)
	ParentID *string `json:"parent_id,omitempty"`
}

// Validate ensures the category is valid
func (c Category) Validate() error {
	if c.ID == "" {
		return errors.New("category ID cannot be empty")
	}
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}
	return nil
}

// Manufacturer represents a medical equipment manufacturer
type Manufacturer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
	Website string `json:"website,omitempty"`
}

// Validate ensures the manufacturer is valid
func (m Manufacturer) Validate() error {
	if m.ID == "" {
		return errors.New("manufacturer ID cannot be empty")
	}
	if m.Name == "" {
		return errors.New("manufacturer name cannot be empty")
	}
	return nil
}

// Price represents a monetary value with currency
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Validate ensures the price is valid
func (p Price) Validate() error {
	if p.Amount < 0 {
		return ErrInvalidPrice
	}
	if p.Currency == "" {
		return errors.New("currency cannot be empty")
	}
	return nil
}

// Specifications represents flexible equipment specifications as JSON
type Specifications map[string]interface{}

// Validate ensures the specifications are valid
func (s Specifications) Validate() error {
	// Ensure specifications can be marshaled to JSON
	_, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSpecifications, err)
	}
	return nil
}

// Equipment represents a medical equipment in the catalog
type Equipment struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Category      Category       `json:"category"`
	Manufacturer  Manufacturer   `json:"manufacturer"`
	Model         string         `json:"model"`
	Description   string         `json:"description"`
	Specifications Specifications `json:"specifications"`
	Price         Price          `json:"price"`
	SKU           string         `json:"sku,omitempty"`
	Images        []string       `json:"images,omitempty"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	tenantID      string         `json:"tenant_id"`
}

// TenantID provides read-only access to the private tenant identifier.
func (e *Equipment) TenantID() string {
	return e.tenantID
}

// SetTenantID sets the tenant identifier (used by repository layer)
func (e *Equipment) SetTenantID(tenantID string) {
	e.tenantID = tenantID
}

// NewEquipment creates a new equipment with a generated ID
func NewEquipment(
	name string,
	category Category,
	manufacturer Manufacturer,
	model string,
	description string,
	specs Specifications,
	price Price,
	tenantID string,
) (*Equipment, error) {
	// Generate a new ULID for the equipment ID
	id := ulid.Make().String()
	
	now := time.Now().UTC()
	
	equipment := &Equipment{
		ID:            id,
		Name:          name,
		Category:      category,
		Manufacturer:  manufacturer,
		Model:         model,
		Description:   description,
		Specifications: specs,
		Price:         price,
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
		tenantID:      tenantID,
	}
	
	// Validate the new equipment
	if err := equipment.Validate(); err != nil {
		return nil, err
	}
	
	return equipment, nil
}

// Validate ensures the equipment is valid
func (e *Equipment) Validate() error {
	if e.ID == "" {
		return ErrInvalidEquipmentID
	}
	
	if e.Name == "" {
		return ErrInvalidEquipmentName
	}
	
	if err := e.Category.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidCategory, err)
	}
	
	if err := e.Manufacturer.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidManufacturer, err)
	}
	
	if err := e.Price.Validate(); err != nil {
		return err
	}
	
	if e.Specifications != nil {
		if err := e.Specifications.Validate(); err != nil {
			return err
		}
	}
	
	if e.tenantID == "" {
		return errors.New("tenant ID cannot be empty")
	}
	
	return nil
}

// Update updates the equipment with new values
func (e *Equipment) Update(
	name string,
	category Category,
	manufacturer Manufacturer,
	model string,
	description string,
	specs Specifications,
	price Price,
	isActive bool,
) error {
	e.Name = name
	e.Category = category
	e.Manufacturer = manufacturer
	e.Model = model
	e.Description = description
	e.Specifications = specs
	e.Price = price
	e.IsActive = isActive
	e.UpdatedAt = time.Now().UTC()
	
	return e.Validate()
}

// SearchCriteria defines parameters for searching equipment
type SearchCriteria struct {
	Query         string   `json:"query"`
	CategoryIDs   []string `json:"category_ids,omitempty"`
	ManufacturerIDs []string `json:"manufacturer_ids,omitempty"`
	PriceMin      *float64 `json:"price_min,omitempty"`
	PriceMax      *float64 `json:"price_max,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
	Page          int      `json:"page"`
	PageSize      int      `json:"page_size"`
	SortBy        string   `json:"sort_by,omitempty"`
	SortDirection string   `json:"sort_direction,omitempty"`
	TenantID      string   `json:"tenant_id"`
}

// CatalogRepository defines the interface for equipment persistence
type CatalogRepository interface {
	// CRUD operations
	Create(ctx context.Context, equipment *Equipment) error
	GetByID(ctx context.Context, id string, tenantID string) (*Equipment, error)
	Update(ctx context.Context, equipment *Equipment) error
	Delete(ctx context.Context, id string, tenantID string) error
	
	// Search operations
	Search(ctx context.Context, criteria SearchCriteria) ([]*Equipment, int, error)
	ListByCategory(ctx context.Context, categoryID string, tenantID string, page, pageSize int) ([]*Equipment, int, error)
	ListByManufacturer(ctx context.Context, manufacturerID string, tenantID string, page, pageSize int) ([]*Equipment, int, error)
	
	// Category and manufacturer operations
	ListCategories(ctx context.Context, tenantID string) ([]Category, error)
	ListManufacturers(ctx context.Context, tenantID string) ([]Manufacturer, error)
}

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	TenantID() string
}

// BaseDomainEvent provides common fields for all domain events
type BaseDomainEvent struct {
	Type       string    `json:"type"`
	ID         string    `json:"id"`
	Aggregate  string    `json:"aggregate_id"`
	Timestamp  time.Time `json:"occurred_at"`
	tenantID   string    `json:"tenant_id"`
}

func (e BaseDomainEvent) EventType() string {
	return e.Type
}

func (e BaseDomainEvent) AggregateID() string {
	return e.Aggregate
}

func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.Timestamp
}

func (e BaseDomainEvent) TenantID() string {
	return e.tenantID
}

// EquipmentCreatedEvent is emitted when a new equipment is created
type EquipmentCreatedEvent struct {
	BaseDomainEvent
	Equipment *Equipment `json:"equipment"`
}

// NewEquipmentCreatedEvent creates a new equipment created event
func NewEquipmentCreatedEvent(equipment *Equipment) *EquipmentCreatedEvent {
	return &EquipmentCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:      "equipment.created",
			ID:        ulid.Make().String(),
			Aggregate: equipment.ID,
			Timestamp: time.Now().UTC(),
			tenantID:  equipment.TenantID(),
		},
		Equipment: equipment,
	}
}

// EquipmentUpdatedEvent is emitted when an equipment is updated
type EquipmentUpdatedEvent struct {
	BaseDomainEvent
	Equipment *Equipment `json:"equipment"`
}

// NewEquipmentUpdatedEvent creates a new equipment updated event
func NewEquipmentUpdatedEvent(equipment *Equipment) *EquipmentUpdatedEvent {
	return &EquipmentUpdatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:      "equipment.updated",
			ID:        ulid.Make().String(),
			Aggregate: equipment.ID,
			Timestamp: time.Now().UTC(),
			tenantID:  equipment.TenantID(),
		},
		Equipment: equipment,
	}
}

// EquipmentDeletedEvent is emitted when an equipment is deleted
type EquipmentDeletedEvent struct {
	BaseDomainEvent
	EquipmentID string `json:"equipment_id"`
}

// NewEquipmentDeletedEvent creates a new equipment deleted event
func NewEquipmentDeletedEvent(equipmentID string, tenantID string) *EquipmentDeletedEvent {
	return &EquipmentDeletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			Type:      "equipment.deleted",
			ID:        ulid.Make().String(),
			Aggregate: equipmentID,
			Timestamp: time.Now().UTC(),
			tenantID:  tenantID,
		},
		EquipmentID: equipmentID,
	}
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}
