package domain

import (
	"errors"
	"time"
)

// Supplier Status
type SupplierStatus string

const (
	SupplierStatusPending  SupplierStatus = "pending"   // Registration submitted, awaiting verification
	SupplierStatusActive   SupplierStatus = "active"    // Verified and active
	SupplierStatusSuspended SupplierStatus = "suspended" // Temporarily suspended
	SupplierStatusInactive SupplierStatus = "inactive"  // Deactivated
)

// Verification Status
type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationApproved VerificationStatus = "approved"
	VerificationRejected VerificationStatus = "rejected"
)

// Errors
var (
	ErrSupplierNotFound      = errors.New("supplier not found")
	ErrInvalidSupplierData   = errors.New("invalid supplier data")
	ErrSupplierAlreadyExists = errors.New("supplier already exists")
	ErrCannotModifySupplier  = errors.New("cannot modify supplier in current status")
)

// ContactInfo represents supplier contact details
type ContactInfo struct {
	PrimaryContactName  string `json:"primary_contact_name"`
	PrimaryContactEmail string `json:"primary_contact_email"`
	PrimaryContactPhone string `json:"primary_contact_phone"`
	SecondaryContactName  string `json:"secondary_contact_name,omitempty"`
	SecondaryContactEmail string `json:"secondary_contact_email,omitempty"`
	SecondaryContactPhone string `json:"secondary_contact_phone,omitempty"`
	Website             string `json:"website,omitempty"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Certification represents a supplier certification
type Certification struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`          // e.g., "ISO 9001", "FDA Registration"
	IssuingBody   string    `json:"issuing_body"`  // e.g., "ISO", "FDA"
	CertNumber    string    `json:"cert_number"`
	IssueDate     time.Time `json:"issue_date"`
	ExpiryDate    time.Time `json:"expiry_date"`
	DocumentURL   string    `json:"document_url,omitempty"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	VerifiedBy    string    `json:"verified_by,omitempty"`
}

// Supplier represents a medical equipment supplier
type Supplier struct {
	// Identity
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	
	// Company Information
	CompanyName             string `json:"company_name"`
	BusinessRegistrationNum string `json:"business_registration_number"`
	TaxID                   string `json:"tax_id"`
	YearEstablished         int    `json:"year_established,omitempty"`
	
	// Contact & Location
	ContactInfo ContactInfo `json:"contact_info"`
	Address     Address     `json:"address"`
	
	// Categorization
	Specializations []string `json:"specializations"` // Category IDs they specialize in
	
	// Certifications & Compliance
	Certifications []Certification `json:"certifications"`
	
	// Performance & Rating
	PerformanceRating float64 `json:"performance_rating"` // 0.0 to 5.0
	TotalOrders       int     `json:"total_orders"`
	CompletedOrders   int     `json:"completed_orders"`
	
	// Status & Verification
	Status             SupplierStatus     `json:"status"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	VerifiedAt         *time.Time         `json:"verified_at,omitempty"`
	VerifiedBy         string             `json:"verified_by,omitempty"`
	
	// Additional Information
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	
	// Audit
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewSupplier creates a new supplier in pending status
func NewSupplier(
	id, tenantID, companyName, businessRegNum, taxID string,
	contactInfo ContactInfo,
	address Address,
	createdBy string,
) (*Supplier, error) {
	// Validation
	if companyName == "" {
		return nil, errors.New("company name is required")
	}
	if contactInfo.PrimaryContactEmail == "" {
		return nil, errors.New("primary contact email is required")
	}
	if address.City == "" || address.Country == "" {
		return nil, errors.New("address city and country are required")
	}

	now := time.Now()
	
	return &Supplier{
		ID:                      id,
		TenantID:                tenantID,
		CompanyName:             companyName,
		BusinessRegistrationNum: businessRegNum,
		TaxID:                   taxID,
		ContactInfo:             contactInfo,
		Address:                 address,
		Specializations:         []string{},
		Certifications:          []Certification{},
		PerformanceRating:       0.0,
		TotalOrders:             0,
		CompletedOrders:         0,
		Status:                  SupplierStatusPending,
		VerificationStatus:      VerificationPending,
		Metadata:                make(map[string]interface{}),
		CreatedBy:               createdBy,
		CreatedAt:               now,
		UpdatedAt:               now,
	}, nil
}

// Verify approves the supplier for business
func (s *Supplier) Verify(verifiedBy string) error {
	if s.VerificationStatus == VerificationApproved {
		return errors.New("supplier is already verified")
	}
	
	now := time.Now()
	s.VerificationStatus = VerificationApproved
	s.VerifiedAt = &now
	s.VerifiedBy = verifiedBy
	s.Status = SupplierStatusActive
	s.UpdatedAt = now
	
	return nil
}

// Reject rejects the supplier verification
func (s *Supplier) Reject(rejectedBy string) error {
	if s.VerificationStatus == VerificationRejected {
		return errors.New("supplier is already rejected")
	}
	
	now := time.Now()
	s.VerificationStatus = VerificationRejected
	s.VerifiedBy = rejectedBy
	s.VerifiedAt = &now
	s.Status = SupplierStatusInactive
	s.UpdatedAt = now
	
	return nil
}

// Suspend temporarily suspends the supplier
func (s *Supplier) Suspend() error {
	if s.Status == SupplierStatusInactive {
		return errors.New("cannot suspend inactive supplier")
	}
	
	s.Status = SupplierStatusSuspended
	s.UpdatedAt = time.Now()
	
	return nil
}

// Activate reactivates a suspended supplier
func (s *Supplier) Activate() error {
	if s.VerificationStatus != VerificationApproved {
		return errors.New("supplier must be verified to be activated")
	}
	
	s.Status = SupplierStatusActive
	s.UpdatedAt = time.Now()
	
	return nil
}

// Deactivate permanently deactivates the supplier
func (s *Supplier) Deactivate() error {
	s.Status = SupplierStatusInactive
	s.UpdatedAt = time.Now()
	return nil
}

// AddCertification adds a certification to the supplier
func (s *Supplier) AddCertification(cert Certification) error {
	if cert.Name == "" || cert.CertNumber == "" {
		return errors.New("certification name and number are required")
	}
	
	s.Certifications = append(s.Certifications, cert)
	s.UpdatedAt = time.Now()
	
	return nil
}

// AddSpecialization adds a category specialization
func (s *Supplier) AddSpecialization(categoryID string) error {
	// Check if already exists
	for _, spec := range s.Specializations {
		if spec == categoryID {
			return nil // Already exists
		}
	}
	
	s.Specializations = append(s.Specializations, categoryID)
	s.UpdatedAt = time.Now()
	
	return nil
}

// UpdatePerformanceRating updates the supplier's performance rating
func (s *Supplier) UpdatePerformanceRating(newRating float64) error {
	if newRating < 0.0 || newRating > 5.0 {
		return errors.New("rating must be between 0.0 and 5.0")
	}
	
	s.PerformanceRating = newRating
	s.UpdatedAt = time.Now()
	
	return nil
}

// RecordOrder records a new order for statistics
func (s *Supplier) RecordOrder(completed bool) {
	s.TotalOrders++
	if completed {
		s.CompletedOrders++
	}
	s.UpdatedAt = time.Now()
}

// CanBeModified checks if supplier details can be modified
func (s *Supplier) CanBeModified() bool {
	return s.Status == SupplierStatusPending || s.Status == SupplierStatusActive
}

// IsActive checks if supplier is active and can receive RFQs
func (s *Supplier) IsActive() bool {
	return s.Status == SupplierStatusActive && s.VerificationStatus == VerificationApproved
}
