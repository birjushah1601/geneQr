package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/contract/domain"
	"github.com/segmentio/ksuid"
)

// ContractService provides application-level contract management operations
type ContractService struct {
	repo   domain.Repository
	logger *slog.Logger
}

// NewContractService creates a new contract service
func NewContractService(repo domain.Repository, logger *slog.Logger) *ContractService {
	return &ContractService{
		repo:   repo,
		logger: logger.With(slog.String("service", "contract")),
	}
}

// CreateContract creates a new contract
func (s *ContractService) CreateContract(ctx context.Context, tenantID, createdBy string, req CreateContractRequest) (*ContractResponse, error) {
	s.logger.Info("Creating contract", slog.String("tenant_id", tenantID), slog.String("rfq_id", req.RFQID))

	// Create contract entity
	contract := domain.NewContract(tenantID, req.RFQID, req.QuoteID, req.SupplierID, req.SupplierName, createdBy)
	contract.ID = ksuid.New().String()
	
	// Generate contract number using database function (we'll call it via direct query in repository)
	// For now, generate a simple contract number
	contract.ContractNumber = fmt.Sprintf("CT-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	
	// Set dates and terms
	contract.StartDate = req.StartDate
	contract.EndDate = req.EndDate
	contract.PaymentTerms = req.PaymentTerms
	contract.DeliveryTerms = req.DeliveryTerms
	contract.WarrantyTerms = req.WarrantyTerms
	contract.TermsAndConditions = req.TermsAndConditions
	contract.TaxAmount = req.TaxAmount
	contract.Notes = req.Notes

	// Add items
	for _, itemReq := range req.Items {
		item := domain.ContractItem{
			ID:               ksuid.New().String(),
			EquipmentID:      itemReq.EquipmentID,
			EquipmentName:    itemReq.EquipmentName,
			Quantity:         itemReq.Quantity,
			UnitPrice:        itemReq.UnitPrice,
			TotalPrice:       float64(itemReq.Quantity) * itemReq.UnitPrice,
			ManufacturerName: itemReq.ManufacturerName,
			ModelNumber:      itemReq.ModelNumber,
			Specifications:   itemReq.Specifications,
			WarrantyPeriod:   itemReq.WarrantyPeriod,
		}
		if err := contract.AddItem(item); err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}
	}

	// Add payment schedule
	for _, paymentReq := range req.PaymentSchedule {
		payment := domain.PaymentTerm{
			DueDate:     paymentReq.DueDate,
			Amount:      paymentReq.Amount,
			Description: paymentReq.Description,
			Paid:        false,
		}
		if err := contract.AddPaymentTerm(payment); err != nil {
			return nil, fmt.Errorf("failed to add payment term: %w", err)
		}
	}

	// Add delivery schedule
	for _, deliveryReq := range req.DeliverySchedule {
		delivery := domain.DeliverySchedule{
			MilestoneDate: deliveryReq.MilestoneDate,
			Description:   deliveryReq.Description,
			Completed:     false,
		}
		if err := contract.AddDeliveryMilestone(delivery); err != nil {
			return nil, fmt.Errorf("failed to add delivery milestone: %w", err)
		}
	}

	// Calculate totals
	contract.CalculateTotals()

	// Save to repository
	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	s.logger.Info("Contract created successfully", slog.String("id", contract.ID))
	response := ToContractResponse(contract)
	return &response, nil
}

// GetContract retrieves a contract by ID
func (s *ContractService) GetContract(ctx context.Context, tenantID, id string) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// GetContractByNumber retrieves a contract by contract number
func (s *ContractService) GetContractByNumber(ctx context.Context, tenantID, contractNumber string) (*ContractResponse, error) {
	contract, err := s.repo.GetByContractNumber(ctx, tenantID, contractNumber)
	if err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// GetContractsByRFQ retrieves all contracts for an RFQ
func (s *ContractService) GetContractsByRFQ(ctx context.Context, tenantID, rfqID string) ([]ContractResponse, error) {
	contracts, err := s.repo.GetByRFQ(ctx, tenantID, rfqID)
	if err != nil {
		return nil, err
	}

	responses := make([]ContractResponse, len(contracts))
	for i, contract := range contracts {
		responses[i] = ToContractResponse(contract)
	}

	return responses, nil
}

// GetContractsBySupplier retrieves all contracts for a supplier
func (s *ContractService) GetContractsBySupplier(ctx context.Context, tenantID, supplierID string) ([]ContractResponse, error) {
	contracts, err := s.repo.GetBySupplier(ctx, tenantID, supplierID)
	if err != nil {
		return nil, err
	}

	responses := make([]ContractResponse, len(contracts))
	for i, contract := range contracts {
		responses[i] = ToContractResponse(contract)
	}

	return responses, nil
}

// ListContracts retrieves contracts with filtering and pagination
func (s *ContractService) ListContracts(ctx context.Context, tenantID string, req ListContractsRequest) (*ListContractsResponse, error) {
	// Convert to domain criteria
	statuses := make([]domain.ContractStatus, len(req.Status))
	for i, status := range req.Status {
		statuses[i] = domain.ContractStatus(status)
	}

	criteria := domain.ListCriteria{
		TenantID:      tenantID,
		RFQID:         req.RFQID,
		SupplierID:    req.SupplierID,
		Status:        statuses,
		CreatedBy:     req.CreatedBy,
		SortBy:        req.SortBy,
		SortDirection: req.SortDirection,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}

	result, err := s.repo.List(ctx, criteria)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]ContractResponse, len(result.Contracts))
	for i, contract := range result.Contracts {
		responses[i] = ToContractResponse(contract)
	}

	return &ListContractsResponse{
		Contracts:  responses,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UpdateContract updates a contract
func (s *ContractService) UpdateContract(ctx context.Context, tenantID, id string, req UpdateContractRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.StartDate != nil {
		contract.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		contract.EndDate = *req.EndDate
	}
	if req.PaymentTerms != nil {
		contract.PaymentTerms = *req.PaymentTerms
	}
	if req.DeliveryTerms != nil {
		contract.DeliveryTerms = *req.DeliveryTerms
	}
	if req.WarrantyTerms != nil {
		contract.WarrantyTerms = *req.WarrantyTerms
	}
	if req.TermsAndConditions != nil {
		contract.TermsAndConditions = *req.TermsAndConditions
	}
	if req.TaxAmount != nil {
		contract.TaxAmount = *req.TaxAmount
		contract.CalculateTotals()
	}
	if req.Notes != nil {
		contract.Notes = *req.Notes
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// ActivateContract activates a contract
func (s *ContractService) ActivateContract(ctx context.Context, tenantID, id string) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Activate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// SignContract signs a contract
func (s *ContractService) SignContract(ctx context.Context, tenantID, id string, req SignContractRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Sign(req.SignedBy); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// CompleteContract completes a contract
func (s *ContractService) CompleteContract(ctx context.Context, tenantID, id string) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Complete(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// CancelContract cancels a contract
func (s *ContractService) CancelContract(ctx context.Context, tenantID, id string, req CancelContractRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Cancel(req.Reason); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// SuspendContract suspends a contract
func (s *ContractService) SuspendContract(ctx context.Context, tenantID, id string, req SuspendContractRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Suspend(req.Reason); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// ResumeContract resumes a suspended contract
func (s *ContractService) ResumeContract(ctx context.Context, tenantID, id string) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.Resume(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// AddAmendment adds an amendment to a contract
func (s *ContractService) AddAmendment(ctx context.Context, tenantID, id string, req AddAmendmentRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	amendment := domain.Amendment{
		ID:          ksuid.New().String(),
		Date:        time.Now(),
		Description: req.Description,
		Changes:     req.Changes,
		AmendedBy:   req.AmendedBy,
	}

	if err := contract.AddAmendment(amendment); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// MarkPaymentPaid marks a payment as paid
func (s *ContractService) MarkPaymentPaid(ctx context.Context, tenantID, id string, req MarkPaymentPaidRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.MarkPaymentPaid(req.PaymentIndex); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// MarkDeliveryCompleted marks a delivery milestone as completed
func (s *ContractService) MarkDeliveryCompleted(ctx context.Context, tenantID, id string, req MarkDeliveryCompletedRequest) (*ContractResponse, error) {
	contract, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if err := contract.MarkDeliveryCompleted(req.DeliveryIndex); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	response := ToContractResponse(contract)
	return &response, nil
}

// DeleteContract deletes a contract
func (s *ContractService) DeleteContract(ctx context.Context, tenantID, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
