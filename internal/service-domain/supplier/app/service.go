package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aby-med/medical-platform/internal/service-domain/supplier/domain"
	"github.com/segmentio/ksuid"
)

// SupplierService provides application-level supplier operations
type SupplierService struct {
	repo   domain.SupplierRepository
	logger *slog.Logger
}

// NewSupplierService creates a new supplier service
func NewSupplierService(repo domain.SupplierRepository, logger *slog.Logger) *SupplierService {
	return &SupplierService{
		repo:   repo,
		logger: logger.With(slog.String("component", "supplier_service")),
	}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(ctx context.Context, tenantID string, req CreateSupplierRequest, createdBy string) (*SupplierResponse, error) {
	// Generate ID
	id := ksuid.New().String()

	// Convert DTOs to domain models
	contactInfo := ToContactInfo(req.ContactInfo)
	address := ToAddress(req.Address)

	// Create the supplier
	supplier, err := domain.NewSupplier(
		id,
		tenantID,
		req.CompanyName,
		req.BusinessRegistrationNum,
		req.TaxID,
		contactInfo,
		address,
		createdBy,
	)
	if err != nil {
		s.logger.Error("Failed to create supplier domain model",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("invalid supplier data: %w", err)
	}

	// Set additional fields
	supplier.YearEstablished = req.YearEstablished
	supplier.Description = req.Description
	supplier.Specializations = req.Specializations

	// Check if tax ID already exists
	if req.TaxID != "" {
		existing, err := s.repo.GetByTaxID(ctx, req.TaxID, tenantID)
		if err == nil && existing != nil {
			return nil, domain.ErrSupplierAlreadyExists
		}
	}

	// Persist the supplier
	if err := s.repo.Create(ctx, supplier); err != nil {
		s.logger.Error("Failed to persist supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", id))
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	s.logger.Info("Supplier created successfully",
		slog.String("supplier_id", id),
		slog.String("company_name", req.CompanyName))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// GetSupplier retrieves a supplier by ID
func (s *SupplierService) GetSupplier(ctx context.Context, tenantID string, supplierID string) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			return nil, err
		}
		s.logger.Error("Failed to get supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(ctx context.Context, tenantID string, supplierID string, req UpdateSupplierRequest) (*SupplierResponse, error) {
	// Get existing supplier
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	// Check if supplier can be modified
	if !supplier.CanBeModified() {
		return nil, domain.ErrCannotModifySupplier
	}

	// Update fields
	if req.CompanyName != "" {
		supplier.CompanyName = req.CompanyName
	}
	if req.BusinessRegistrationNum != "" {
		supplier.BusinessRegistrationNum = req.BusinessRegistrationNum
	}
	if req.TaxID != "" {
		supplier.TaxID = req.TaxID
	}
	if req.YearEstablished > 0 {
		supplier.YearEstablished = req.YearEstablished
	}
	if req.Description != "" {
		supplier.Description = req.Description
	}
	if req.ContactInfo != nil {
		supplier.ContactInfo = ToContactInfo(*req.ContactInfo)
	}
	if req.Address != nil {
		supplier.Address = ToAddress(*req.Address)
	}
	if req.Specializations != nil {
		supplier.Specializations = req.Specializations
	}

	// Save updated supplier
	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to update supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// DeleteSupplier deletes a supplier
func (s *SupplierService) DeleteSupplier(ctx context.Context, tenantID string, supplierID string) error {
	err := s.repo.Delete(ctx, supplierID, tenantID)
	if err != nil {
		if err == domain.ErrSupplierNotFound {
			return err
		}
		s.logger.Error("Failed to delete supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	s.logger.Info("Supplier deleted successfully",
		slog.String("supplier_id", supplierID))

	return nil
}

// ListSuppliers retrieves a paginated list of suppliers
func (s *SupplierService) ListSuppliers(ctx context.Context, tenantID string, req ListSuppliersRequest) (*ListSuppliersResponse, error) {
	// Build criteria
	criteria := domain.ListCriteria{
		TenantID:      tenantID,
		CategoryID:    req.CategoryID,
		SearchQuery:   req.SearchQuery,
		MinRating:     req.MinRating,
		Page:          req.Page,
		PageSize:      req.PageSize,
		SortBy:        req.SortBy,
		SortDirection: req.SortDirection,
	}

	// Convert status strings to domain types
	if len(req.Status) > 0 {
		criteria.Status = make([]domain.SupplierStatus, len(req.Status))
		for i, status := range req.Status {
			criteria.Status[i] = domain.SupplierStatus(status)
		}
	}

	// Convert verification status strings to domain types
	if len(req.VerificationStatus) > 0 {
		criteria.VerificationStatus = make([]domain.VerificationStatus, len(req.VerificationStatus))
		for i, status := range req.VerificationStatus {
			criteria.VerificationStatus[i] = domain.VerificationStatus(status)
		}
	}

	// Query suppliers
	suppliers, total, err := s.repo.List(ctx, criteria)
	if err != nil {
		s.logger.Error("Failed to list suppliers",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}

	// Convert to response DTOs
	supplierResponses := make([]SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		supplierResponses[i] = ToSupplierResponse(supplier)
	}

	// Calculate pagination info
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (total + pageSize - 1) / pageSize

	return &ListSuppliersResponse{
		Suppliers:  supplierResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// VerifySupplier verifies a supplier for business
func (s *SupplierService) VerifySupplier(ctx context.Context, tenantID string, supplierID string, verifiedBy string) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := supplier.Verify(verifiedBy); err != nil {
		return nil, fmt.Errorf("failed to verify supplier: %w", err)
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to save verified supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to save verified supplier: %w", err)
	}

	s.logger.Info("Supplier verified successfully",
		slog.String("supplier_id", supplierID),
		slog.String("verified_by", verifiedBy))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// RejectSupplier rejects a supplier verification
func (s *SupplierService) RejectSupplier(ctx context.Context, tenantID string, supplierID string, rejectedBy string) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := supplier.Reject(rejectedBy); err != nil {
		return nil, fmt.Errorf("failed to reject supplier: %w", err)
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to save rejected supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to save rejected supplier: %w", err)
	}

	s.logger.Info("Supplier rejected",
		slog.String("supplier_id", supplierID),
		slog.String("rejected_by", rejectedBy))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// SuspendSupplier suspends a supplier
func (s *SupplierService) SuspendSupplier(ctx context.Context, tenantID string, supplierID string) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := supplier.Suspend(); err != nil {
		return nil, fmt.Errorf("failed to suspend supplier: %w", err)
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to save suspended supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to save suspended supplier: %w", err)
	}

	s.logger.Info("Supplier suspended",
		slog.String("supplier_id", supplierID))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// ActivateSupplier activates a suspended supplier
func (s *SupplierService) ActivateSupplier(ctx context.Context, tenantID string, supplierID string) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	if err := supplier.Activate(); err != nil {
		return nil, fmt.Errorf("failed to activate supplier: %w", err)
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to save activated supplier",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to save activated supplier: %w", err)
	}

	s.logger.Info("Supplier activated",
		slog.String("supplier_id", supplierID))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// AddCertification adds a certification to a supplier
func (s *SupplierService) AddCertification(ctx context.Context, tenantID string, supplierID string, cert CertificationDTO) (*SupplierResponse, error) {
	supplier, err := s.repo.GetByID(ctx, supplierID, tenantID)
	if err != nil {
		return nil, err
	}

	// Generate certification ID if not provided
	if cert.ID == "" {
		cert.ID = ksuid.New().String()
	}

	domainCert := ToCertification(cert)
	if err := supplier.AddCertification(domainCert); err != nil {
		return nil, fmt.Errorf("failed to add certification: %w", err)
	}

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("Failed to save supplier with new certification",
			slog.String("error", err.Error()),
			slog.String("supplier_id", supplierID))
		return nil, fmt.Errorf("failed to save certification: %w", err)
	}

	s.logger.Info("Certification added to supplier",
		slog.String("supplier_id", supplierID),
		slog.String("certification_name", cert.Name))

	response := ToSupplierResponse(supplier)
	return &response, nil
}

// GetSuppliersByCategory retrieves suppliers specialized in a category
func (s *SupplierService) GetSuppliersByCategory(ctx context.Context, tenantID string, categoryID string) ([]SupplierResponse, error) {
	suppliers, err := s.repo.GetByCategory(ctx, categoryID, tenantID)
	if err != nil {
		s.logger.Error("Failed to get suppliers by category",
			slog.String("error", err.Error()),
			slog.String("category_id", categoryID))
		return nil, fmt.Errorf("failed to get suppliers by category: %w", err)
	}

	// Convert to response DTOs
	responses := make([]SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		responses[i] = ToSupplierResponse(supplier)
	}

	return responses, nil
}
