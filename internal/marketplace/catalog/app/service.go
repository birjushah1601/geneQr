package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aby-med/medical-platform/internal/marketplace/catalog/domain"
)

// DTOs for input/output operations

// EquipmentDTO represents equipment data for API operations
type EquipmentDTO struct {
	ID            string                 `json:"id,omitempty"`
	Name          string                 `json:"name"`
	CategoryID    string                 `json:"category_id"`
	CategoryName  string                 `json:"category_name"`
	ManufacturerID string                `json:"manufacturer_id"`
	ManufacturerName string              `json:"manufacturer_name"`
	Model         string                 `json:"model"`
	Description   string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Price         float64                `json:"price"`
	Currency      string                 `json:"currency"`
	SKU           string                 `json:"sku,omitempty"`
	Images        []string               `json:"images,omitempty"`
	IsActive      bool                   `json:"is_active"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
	UpdatedAt     time.Time              `json:"updated_at,omitempty"`
}

// CreateEquipmentRequest represents the data needed to create new equipment
type CreateEquipmentRequest struct {
	Name          string                 `json:"name" validate:"required"`
	CategoryID    string                 `json:"category_id" validate:"required"`
	ManufacturerID string                `json:"manufacturer_id" validate:"required"`
	Model         string                 `json:"model" validate:"required"`
	Description   string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Price         float64                `json:"price" validate:"required,gt=0"`
	Currency      string                 `json:"currency" validate:"required"`
	SKU           string                 `json:"sku,omitempty"`
	Images        []string               `json:"images,omitempty"`
}

// UpdateEquipmentRequest represents the data needed to update equipment
type UpdateEquipmentRequest struct {
	Name          string                 `json:"name" validate:"required"`
	CategoryID    string                 `json:"category_id" validate:"required"`
	ManufacturerID string                `json:"manufacturer_id" validate:"required"`
	Model         string                 `json:"model" validate:"required"`
	Description   string                 `json:"description"`
	Specifications map[string]interface{} `json:"specifications"`
	Price         float64                `json:"price" validate:"required,gt=0"`
	Currency      string                 `json:"currency" validate:"required"`
	SKU           string                 `json:"sku,omitempty"`
	Images        []string               `json:"images,omitempty"`
	IsActive      bool                   `json:"is_active"`
}

// SearchEquipmentRequest represents search criteria for equipment
type SearchEquipmentRequest struct {
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
}

// PaginatedResponse represents a paginated list of items
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	TotalItems int         `json:"total_items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// CategoryDTO represents category data for API operations
type CategoryDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id,omitempty"`
}

// ManufacturerDTO represents manufacturer data for API operations
type ManufacturerDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
	Website string `json:"website,omitempty"`
}

// CatalogService orchestrates domain operations for the catalog module
type CatalogService struct {
	repository domain.CatalogRepository
	eventBus   domain.EventPublisher
	logger     *slog.Logger
}

// NewCatalogService creates a new catalog service
func NewCatalogService(
	repository domain.CatalogRepository,
	eventBus domain.EventPublisher,
	logger *slog.Logger,
) *CatalogService {
	return &CatalogService{
		repository: repository,
		eventBus:   eventBus,
		logger:     logger.With(slog.String("component", "catalog_service")),
	}
}

// CreateEquipment creates new equipment in the catalog
func (s *CatalogService) CreateEquipment(ctx context.Context, req CreateEquipmentRequest) (*EquipmentDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch category and manufacturer to ensure they exist
	categories, err := s.repository.ListCategories(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to fetch categories", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var category domain.Category
	categoryFound := false
	for _, c := range categories {
		if c.ID == req.CategoryID {
			category = c
			categoryFound = true
			break
		}
	}

	if !categoryFound {
		return nil, domain.ErrInvalidCategory
	}

	manufacturers, err := s.repository.ListManufacturers(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to fetch manufacturers", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to fetch manufacturers: %w", err)
	}

	var manufacturer domain.Manufacturer
	manufacturerFound := false
	for _, m := range manufacturers {
		if m.ID == req.ManufacturerID {
			manufacturer = m
			manufacturerFound = true
			break
		}
	}

	if !manufacturerFound {
		return nil, domain.ErrInvalidManufacturer
	}

	// Create domain entity
	price := domain.Price{
		Amount:   req.Price,
		Currency: req.Currency,
	}

	equipment, err := domain.NewEquipment(
		req.Name,
		category,
		manufacturer,
		req.Model,
		req.Description,
		domain.Specifications(req.Specifications),
		price,
		tenantID,
	)
	if err != nil {
		s.logger.Error("Failed to create equipment entity", slog.String("error", err.Error()))
		return nil, err
	}

	// Set optional fields
	equipment.SKU = req.SKU
	equipment.Images = req.Images

	// Persist to repository
	if err := s.repository.Create(ctx, equipment); err != nil {
		s.logger.Error("Failed to persist equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		return nil, fmt.Errorf("failed to create equipment: %w", err)
	}

	// Publish domain event
	event := domain.NewEquipmentCreatedEvent(equipment)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish equipment created event",
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		// Continue despite event publishing failure
	}

	s.logger.Info("Equipment created successfully", 
		slog.String("equipment_id", equipment.ID),
		slog.String("tenant_id", tenantID))

	// Convert to DTO and return
	return s.mapEquipmentToDTO(equipment), nil
}

// GetEquipment retrieves equipment by ID
func (s *CatalogService) GetEquipment(ctx context.Context, id string) (*EquipmentDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	equipment, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			s.logger.Info("Equipment not found", 
				slog.String("equipment_id", id),
				slog.String("tenant_id", tenantID))
			return nil, err
		}
		s.logger.Error("Failed to get equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return nil, fmt.Errorf("failed to get equipment: %w", err)
	}

	return s.mapEquipmentToDTO(equipment), nil
}

// UpdateEquipment updates existing equipment
func (s *CatalogService) UpdateEquipment(ctx context.Context, id string, req UpdateEquipmentRequest) (*EquipmentDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Fetch existing equipment
	equipment, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			s.logger.Info("Equipment not found for update", 
				slog.String("equipment_id", id),
				slog.String("tenant_id", tenantID))
			return nil, err
		}
		s.logger.Error("Failed to get equipment for update", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return nil, fmt.Errorf("failed to get equipment for update: %w", err)
	}

	// Fetch category and manufacturer to ensure they exist
	categories, err := s.repository.ListCategories(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to fetch categories", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var category domain.Category
	categoryFound := false
	for _, c := range categories {
		if c.ID == req.CategoryID {
			category = c
			categoryFound = true
			break
		}
	}

	if !categoryFound {
		return nil, domain.ErrInvalidCategory
	}

	manufacturers, err := s.repository.ListManufacturers(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to fetch manufacturers", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to fetch manufacturers: %w", err)
	}

	var manufacturer domain.Manufacturer
	manufacturerFound := false
	for _, m := range manufacturers {
		if m.ID == req.ManufacturerID {
			manufacturer = m
			manufacturerFound = true
			break
		}
	}

	if !manufacturerFound {
		return nil, domain.ErrInvalidManufacturer
	}

	// Update domain entity
	price := domain.Price{
		Amount:   req.Price,
		Currency: req.Currency,
	}

	if err := equipment.Update(
		req.Name,
		category,
		manufacturer,
		req.Model,
		req.Description,
		domain.Specifications(req.Specifications),
		price,
		req.IsActive,
	); err != nil {
		s.logger.Error("Failed to update equipment entity", slog.String("error", err.Error()))
		return nil, err
	}

	// Set optional fields
	equipment.SKU = req.SKU
	equipment.Images = req.Images

	// Persist to repository
	if err := s.repository.Update(ctx, equipment); err != nil {
		s.logger.Error("Failed to persist updated equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		return nil, fmt.Errorf("failed to update equipment: %w", err)
	}

	// Publish domain event
	event := domain.NewEquipmentUpdatedEvent(equipment)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish equipment updated event",
			slog.String("error", err.Error()),
			slog.String("equipment_id", equipment.ID))
		// Continue despite event publishing failure
	}

	s.logger.Info("Equipment updated successfully", 
		slog.String("equipment_id", equipment.ID),
		slog.String("tenant_id", tenantID))

	// Convert to DTO and return
	return s.mapEquipmentToDTO(equipment), nil
}

// DeleteEquipment removes equipment from the catalog
func (s *CatalogService) DeleteEquipment(ctx context.Context, id string) error {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return errors.New("tenant ID is required")
	}

	// Check if equipment exists
	_, err := s.repository.GetByID(ctx, id, tenantID)
	if err != nil {
		if errors.Is(err, domain.ErrEquipmentNotFound) {
			s.logger.Info("Equipment not found for deletion", 
				slog.String("equipment_id", id),
				slog.String("tenant_id", tenantID))
			return err
		}
		s.logger.Error("Failed to get equipment for deletion", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return fmt.Errorf("failed to get equipment for deletion: %w", err)
	}

	// Delete from repository
	if err := s.repository.Delete(ctx, id, tenantID); err != nil {
		s.logger.Error("Failed to delete equipment", 
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		return fmt.Errorf("failed to delete equipment: %w", err)
	}

	// Publish domain event
	event := domain.NewEquipmentDeletedEvent(id, tenantID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		s.logger.Error("Failed to publish equipment deleted event",
			slog.String("error", err.Error()),
			slog.String("equipment_id", id))
		// Continue despite event publishing failure
	}

	s.logger.Info("Equipment deleted successfully", 
		slog.String("equipment_id", id),
		slog.String("tenant_id", tenantID))

	return nil
}

// SearchEquipment searches for equipment based on criteria
func (s *CatalogService) SearchEquipment(ctx context.Context, req SearchEquipmentRequest) (*PaginatedResponse, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Set default pagination if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// Convert request to domain search criteria
	criteria := domain.SearchCriteria{
		Query:           req.Query,
		CategoryIDs:     req.CategoryIDs,
		ManufacturerIDs: req.ManufacturerIDs,
		PriceMin:        req.PriceMin,
		PriceMax:        req.PriceMax,
		IsActive:        req.IsActive,
		Page:            req.Page,
		PageSize:        req.PageSize,
		SortBy:          req.SortBy,
		SortDirection:   req.SortDirection,
		TenantID:        tenantID,
	}

	// Perform search
	equipment, total, err := s.repository.Search(ctx, criteria)
	if err != nil {
		s.logger.Error("Failed to search equipment", 
			slog.String("error", err.Error()),
			slog.String("query", req.Query))
		return nil, fmt.Errorf("failed to search equipment: %w", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]*EquipmentDTO, len(equipment))
	for i, e := range equipment {
		dtos[i] = s.mapEquipmentToDTO(e)
	}

	// Calculate total pages
	totalPages := total / req.PageSize
	if total%req.PageSize > 0 {
		totalPages++
	}

	// Create paginated response
	response := &PaginatedResponse{
		Items:      dtos,
		TotalItems: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// ListCategories retrieves all categories
func (s *CatalogService) ListCategories(ctx context.Context) ([]CategoryDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	categories, err := s.repository.ListCategories(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to list categories", 
			slog.String("error", err.Error()),
			slog.String("tenant_id", tenantID))
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]CategoryDTO, len(categories))
	for i, c := range categories {
		dtos[i] = CategoryDTO{
			ID:       c.ID,
			Name:     c.Name,
			ParentID: c.ParentID,
		}
	}

	return dtos, nil
}

// ListManufacturers retrieves all manufacturers
func (s *CatalogService) ListManufacturers(ctx context.Context) ([]ManufacturerDTO, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	manufacturers, err := s.repository.ListManufacturers(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to list manufacturers", 
			slog.String("error", err.Error()),
			slog.String("tenant_id", tenantID))
		return nil, fmt.Errorf("failed to list manufacturers: %w", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]ManufacturerDTO, len(manufacturers))
	for i, m := range manufacturers {
		dtos[i] = ManufacturerDTO{
			ID:      m.ID,
			Name:    m.Name,
			Country: m.Country,
			Website: m.Website,
		}
	}

	return dtos, nil
}

// ListEquipmentByCategory retrieves equipment by category
func (s *CatalogService) ListEquipmentByCategory(ctx context.Context, categoryID string, page, pageSize int) (*PaginatedResponse, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Set default pagination if not provided
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Fetch equipment by category
	equipment, total, err := s.repository.ListByCategory(ctx, categoryID, tenantID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list equipment by category", 
			slog.String("error", err.Error()),
			slog.String("category_id", categoryID))
		return nil, fmt.Errorf("failed to list equipment by category: %w", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]*EquipmentDTO, len(equipment))
	for i, e := range equipment {
		dtos[i] = s.mapEquipmentToDTO(e)
	}

	// Calculate total pages
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	// Create paginated response
	response := &PaginatedResponse{
		Items:      dtos,
		TotalItems: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// ListEquipmentByManufacturer retrieves equipment by manufacturer
func (s *CatalogService) ListEquipmentByManufacturer(ctx context.Context, manufacturerID string, page, pageSize int) (*PaginatedResponse, error) {
	tenantID := domain.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.New("tenant ID is required")
	}

	// Set default pagination if not provided
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Fetch equipment by manufacturer
	equipment, total, err := s.repository.ListByManufacturer(ctx, manufacturerID, tenantID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list equipment by manufacturer", 
			slog.String("error", err.Error()),
			slog.String("manufacturer_id", manufacturerID))
		return nil, fmt.Errorf("failed to list equipment by manufacturer: %w", err)
	}

	// Convert domain entities to DTOs
	dtos := make([]*EquipmentDTO, len(equipment))
	for i, e := range equipment {
		dtos[i] = s.mapEquipmentToDTO(e)
	}

	// Calculate total pages
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	// Create paginated response
	response := &PaginatedResponse{
		Items:      dtos,
		TotalItems: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return response, nil
}

// Helper method to map domain entity to DTO
func (s *CatalogService) mapEquipmentToDTO(equipment *domain.Equipment) *EquipmentDTO {
	return &EquipmentDTO{
		ID:              equipment.ID,
		Name:            equipment.Name,
		CategoryID:      equipment.Category.ID,
		CategoryName:    equipment.Category.Name,
		ManufacturerID:  equipment.Manufacturer.ID,
		ManufacturerName: equipment.Manufacturer.Name,
		Model:           equipment.Model,
		Description:     equipment.Description,
		Specifications:  equipment.Specifications,
		Price:           equipment.Price.Amount,
		Currency:        equipment.Price.Currency,
		SKU:             equipment.SKU,
		Images:          equipment.Images,
		IsActive:        equipment.IsActive,
		CreatedAt:       equipment.CreatedAt,
		UpdatedAt:       equipment.UpdatedAt,
	}
}
