package domain

import "context"

// AssignmentRepository defines data access operations for engineer assignment
type AssignmentRepository interface {
	// Engineer CRUD operations
	ListEngineers(ctx context.Context, organizationID *string, limit, offset int) ([]*Engineer, error)
	GetEngineerByID(ctx context.Context, engineerID string) (*Engineer, error)
	UpdateEngineerLevel(ctx context.Context, engineerID string, level EngineerLevel) error
	
	// Engineer equipment types (capabilities)
	ListEngineerEquipmentTypes(ctx context.Context, engineerID string) ([]*EngineerEquipmentType, error)
	AddEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error
	RemoveEngineerEquipmentType(ctx context.Context, engineerID, manufacturer, category string) error
	
	// Equipment service configuration
	GetEquipmentServiceConfig(ctx context.Context, equipmentID string) (*EquipmentServiceConfig, error)
	CreateEquipmentServiceConfig(ctx context.Context, config *EquipmentServiceConfig) error
	UpdateEquipmentServiceConfig(ctx context.Context, config *EquipmentServiceConfig) error
	
	// Assignment suggestion algorithm
	GetSuggestedEngineers(ctx context.Context, equipmentID string, manufacturer, category string, minLevel EngineerLevel) ([]*SuggestedEngineer, error)
	
	// Manual assignment
	AssignEngineerToTicket(ctx context.Context, req AssignmentRequest) error
}
