package app

import (
	"time"

	"github.com/aby-med/medical-platform/internal/service-domain/supplier/domain"
)

// CreateSupplierRequest represents a request to create a new supplier
type CreateSupplierRequest struct {
	CompanyName             string                `json:"company_name"`
	BusinessRegistrationNum string                `json:"business_registration_number,omitempty"`
	TaxID                   string                `json:"tax_id,omitempty"`
	YearEstablished         int                   `json:"year_established,omitempty"`
	Description             string                `json:"description,omitempty"`
	ContactInfo             ContactInfoDTO        `json:"contact_info"`
	Address                 AddressDTO            `json:"address"`
	Specializations         []string              `json:"specializations,omitempty"`
}

// UpdateSupplierRequest represents a request to update an existing supplier
type UpdateSupplierRequest struct {
	CompanyName             string                `json:"company_name,omitempty"`
	BusinessRegistrationNum string                `json:"business_registration_number,omitempty"`
	TaxID                   string                `json:"tax_id,omitempty"`
	YearEstablished         int                   `json:"year_established,omitempty"`
	Description             string                `json:"description,omitempty"`
	ContactInfo             *ContactInfoDTO       `json:"contact_info,omitempty"`
	Address                 *AddressDTO           `json:"address,omitempty"`
	Specializations         []string              `json:"specializations,omitempty"`
}

// ContactInfoDTO represents contact information
type ContactInfoDTO struct {
	PrimaryContactName    string `json:"primary_contact_name"`
	PrimaryContactEmail   string `json:"primary_contact_email"`
	PrimaryContactPhone   string `json:"primary_contact_phone"`
	SecondaryContactName  string `json:"secondary_contact_name,omitempty"`
	SecondaryContactEmail string `json:"secondary_contact_email,omitempty"`
	SecondaryContactPhone string `json:"secondary_contact_phone,omitempty"`
	Website               string `json:"website,omitempty"`
}

// AddressDTO represents an address
type AddressDTO struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country"`
}

// CertificationDTO represents a certification
type CertificationDTO struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	IssuingBody string    `json:"issuing_body"`
	CertNumber  string    `json:"cert_number"`
	IssueDate   time.Time `json:"issue_date"`
	ExpiryDate  time.Time `json:"expiry_date"`
	DocumentURL string    `json:"document_url,omitempty"`
}

// SupplierResponse represents a supplier response
type SupplierResponse struct {
	ID                      string                 `json:"id"`
	TenantID                string                 `json:"tenant_id"`
	CompanyName             string                 `json:"company_name"`
	BusinessRegistrationNum string                 `json:"business_registration_number,omitempty"`
	TaxID                   string                 `json:"tax_id,omitempty"`
	YearEstablished         int                    `json:"year_established,omitempty"`
	Description             string                 `json:"description,omitempty"`
	ContactInfo             ContactInfoDTO         `json:"contact_info"`
	Address                 AddressDTO             `json:"address"`
	Specializations         []string               `json:"specializations"`
	Certifications          []CertificationDTO     `json:"certifications"`
	PerformanceRating       float64                `json:"performance_rating"`
	TotalOrders             int                    `json:"total_orders"`
	CompletedOrders         int                    `json:"completed_orders"`
	Status                  string                 `json:"status"`
	VerificationStatus      string                 `json:"verification_status"`
	VerifiedAt              *time.Time             `json:"verified_at,omitempty"`
	VerifiedBy              string                 `json:"verified_by,omitempty"`
	Metadata                map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy               string                 `json:"created_by"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
}

// ListSuppliersRequest represents filters for listing suppliers
type ListSuppliersRequest struct {
	Status             []string `json:"status,omitempty"`
	VerificationStatus []string `json:"verification_status,omitempty"`
	CategoryID         string   `json:"category_id,omitempty"`
	SearchQuery        string   `json:"search_query,omitempty"`
	MinRating          float64  `json:"min_rating,omitempty"`
	Page               int      `json:"page,omitempty"`
	PageSize           int      `json:"page_size,omitempty"`
	SortBy             string   `json:"sort_by,omitempty"`
	SortDirection      string   `json:"sort_direction,omitempty"`
}

// ListSuppliersResponse represents a paginated list of suppliers
type ListSuppliersResponse struct {
	Suppliers  []SupplierResponse `json:"suppliers"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// Conversion functions

func ToContactInfo(dto ContactInfoDTO) domain.ContactInfo {
	return domain.ContactInfo{
		PrimaryContactName:    dto.PrimaryContactName,
		PrimaryContactEmail:   dto.PrimaryContactEmail,
		PrimaryContactPhone:   dto.PrimaryContactPhone,
		SecondaryContactName:  dto.SecondaryContactName,
		SecondaryContactEmail: dto.SecondaryContactEmail,
		SecondaryContactPhone: dto.SecondaryContactPhone,
		Website:               dto.Website,
	}
}

func ToAddress(dto AddressDTO) domain.Address {
	return domain.Address{
		Street:     dto.Street,
		City:       dto.City,
		State:      dto.State,
		PostalCode: dto.PostalCode,
		Country:    dto.Country,
	}
}

func ToCertification(dto CertificationDTO) domain.Certification {
	return domain.Certification{
		ID:          dto.ID,
		Name:        dto.Name,
		IssuingBody: dto.IssuingBody,
		CertNumber:  dto.CertNumber,
		IssueDate:   dto.IssueDate,
		ExpiryDate:  dto.ExpiryDate,
		DocumentURL: dto.DocumentURL,
	}
}

func ToSupplierResponse(s *domain.Supplier) SupplierResponse {
	contactInfo := ContactInfoDTO{
		PrimaryContactName:    s.ContactInfo.PrimaryContactName,
		PrimaryContactEmail:   s.ContactInfo.PrimaryContactEmail,
		PrimaryContactPhone:   s.ContactInfo.PrimaryContactPhone,
		SecondaryContactName:  s.ContactInfo.SecondaryContactName,
		SecondaryContactEmail: s.ContactInfo.SecondaryContactEmail,
		SecondaryContactPhone: s.ContactInfo.SecondaryContactPhone,
		Website:               s.ContactInfo.Website,
	}

	address := AddressDTO{
		Street:     s.Address.Street,
		City:       s.Address.City,
		State:      s.Address.State,
		PostalCode: s.Address.PostalCode,
		Country:    s.Address.Country,
	}

	certifications := make([]CertificationDTO, len(s.Certifications))
	for i, cert := range s.Certifications {
		certifications[i] = CertificationDTO{
			ID:          cert.ID,
			Name:        cert.Name,
			IssuingBody: cert.IssuingBody,
			CertNumber:  cert.CertNumber,
			IssueDate:   cert.IssueDate,
			ExpiryDate:  cert.ExpiryDate,
			DocumentURL: cert.DocumentURL,
		}
	}

	return SupplierResponse{
		ID:                      s.ID,
		TenantID:                s.TenantID,
		CompanyName:             s.CompanyName,
		BusinessRegistrationNum: s.BusinessRegistrationNum,
		TaxID:                   s.TaxID,
		YearEstablished:         s.YearEstablished,
		Description:             s.Description,
		ContactInfo:             contactInfo,
		Address:                 address,
		Specializations:         s.Specializations,
		Certifications:          certifications,
		PerformanceRating:       s.PerformanceRating,
		TotalOrders:             s.TotalOrders,
		CompletedOrders:         s.CompletedOrders,
		Status:                  string(s.Status),
		VerificationStatus:      string(s.VerificationStatus),
		VerifiedAt:              s.VerifiedAt,
		VerifiedBy:              s.VerifiedBy,
		Metadata:                s.Metadata,
		CreatedBy:               s.CreatedBy,
		CreatedAt:               s.CreatedAt,
		UpdatedAt:               s.UpdatedAt,
	}
}
