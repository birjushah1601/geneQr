package domain

import (
	"errors"
	"time"
)

var (
	ErrTicketNotFound      = errors.New("ticket not found")
	ErrInvalidStatus       = errors.New("invalid status transition")
	ErrEngineerNotAssigned = errors.New("engineer not assigned")
)

// TicketStatus represents the lifecycle status of a service ticket
type TicketStatus string

const (
	StatusNew        TicketStatus = "new"
	StatusAssigned   TicketStatus = "assigned"
	StatusInProgress TicketStatus = "in_progress"
	StatusOnHold     TicketStatus = "on_hold"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"
	StatusCancelled  TicketStatus = "cancelled"
)

// TicketPriority represents the urgency level
type TicketPriority string

const (
	PriorityCritical TicketPriority = "critical"
	PriorityHigh     TicketPriority = "high"
	PriorityMedium   TicketPriority = "medium"
	PriorityLow      TicketPriority = "low"
)

// TicketSource represents where the ticket originated from
type TicketSource string

const (
	SourceWhatsApp  TicketSource = "whatsapp"
	SourceWeb       TicketSource = "web"
	SourcePhone     TicketSource = "phone"
	SourceEmail     TicketSource = "email"
	SourceScheduled TicketSource = "scheduled"
)

// ServiceTicket represents a customer service request
type ServiceTicket struct {
	ID           string         `json:"id"`
	TicketNumber string         `json:"ticket_number"` // TKT-YYYYMMDD-XXXX
	
	// Equipment & Customer
	EquipmentID   string `json:"equipment_id"`
	QRCode        string `json:"qr_code,omitempty"`
	SerialNumber  string `json:"serial_number"`
	EquipmentName string `json:"equipment_name"`
	CustomerID       string  `json:"customer_id"`
	CustomerName     string  `json:"customer_name"`
	CustomerPhone    string  `json:"customer_phone"`
	CustomerEmail    *string `json:"customer_email,omitempty"`
	CustomerWhatsApp string  `json:"customer_whatsapp,omitempty"`
	
	// Issue details
	IssueCategory    string         `json:"issue_category"` // breakdown, maintenance, installation, inspection
	IssueDescription string         `json:"issue_description"`
	Priority         TicketPriority `json:"priority"`
	Severity         string         `json:"severity,omitempty"`
	
	// Source
	Source          TicketSource `json:"source"`
	SourceMessageID string       `json:"source_message_id,omitempty"` // WhatsApp message ID
	
	// Assignment
	AssignedEngineerID   string `json:"assigned_engineer_id,omitempty"`
	AssignedEngineerName string `json:"assigned_engineer_name,omitempty"`
	AssignedAt           *time.Time `json:"assigned_at,omitempty"`
	
	// Status & Timeline
	Status       TicketStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	ResolvedAt   *time.Time   `json:"resolved_at,omitempty"`
	ClosedAt     *time.Time   `json:"closed_at,omitempty"`
	
	// SLA tracking
	SLAResponseDue   *time.Time `json:"sla_response_due,omitempty"`
	SLAResolutionDue *time.Time `json:"sla_resolution_due,omitempty"`
	SLABreached      bool       `json:"sla_breached"`
	
	// Resolution
	ResolutionNotes string                   `json:"resolution_notes,omitempty"`
	PartsUsed       interface{}              `json:"parts_used,omitempty"` // Can be []Part or []interface{} for flexibility
	LaborHours      float64                  `json:"labor_hours"`
	Cost            float64                  `json:"cost"`
	
	// Media
	Photos    []string `json:"photos"`
	Videos    []string `json:"videos"`
	Documents []string `json:"documents"`
	
	// AMC linkage
	AMCContractID   string `json:"amc_contract_id,omitempty"`
	CoveredUnderAMC bool   `json:"covered_under_amc"`
	
	// Metadata
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
}

// Part represents a spare part used in service
type Part struct {
	PartNumber  string  `json:"part_number"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

// NewServiceTicket creates a new service ticket
func NewServiceTicket(
	equipmentID, serialNumber, equipmentName, customerName, issueDescription string,
	source TicketSource,
	createdBy string,
) *ServiceTicket {
	now := time.Now()
	return &ServiceTicket{
		EquipmentID:      equipmentID,
		SerialNumber:     serialNumber,
		EquipmentName:    equipmentName,
		CustomerName:     customerName,
		IssueDescription: issueDescription,
		Source:           source,
		Priority:         PriorityMedium,
		Status:           StatusNew,
		Photos:           []string{},
		Videos:           []string{},
		Documents:        []string{},
		PartsUsed:        []Part{},
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedBy:        createdBy,
	}
}

// AssignEngineer assigns an engineer to the ticket
func (t *ServiceTicket) AssignEngineer(engineerID, engineerName string) error {
	if t.Status != StatusNew && t.Status != StatusAssigned {
		return ErrInvalidStatus
	}
	
	now := time.Now()
	t.AssignedEngineerID = engineerID
	t.AssignedEngineerName = engineerName
	t.AssignedAt = &now
	t.Status = StatusAssigned
	t.UpdatedAt = now
	
	return nil
}

// Acknowledge acknowledges the ticket (admin/system has seen it)
// This just sets acknowledged_at timestamp, doesn't change status
// Status changes to 'assigned' only when engineer is actually assigned
func (t *ServiceTicket) Acknowledge() error {
	// Allow acknowledging from 'new' status
	if t.Status != StatusNew {
		return ErrInvalidStatus
	}
	
	now := time.Now()
	t.AcknowledgedAt = &now
	// Note: Status stays as 'new' until engineer is assigned
	t.UpdatedAt = now
	
	return nil
}

// Start marks the ticket as in progress
func (t *ServiceTicket) Start() error {
	if t.Status != StatusAssigned {
		return ErrInvalidStatus
	}
	
	if t.AssignedEngineerID == "" {
		return ErrEngineerNotAssigned
	}
	
	now := time.Now()
	t.StartedAt = &now
	t.Status = StatusInProgress
	t.UpdatedAt = now
	
	return nil
}

// PutOnHold temporarily halts work on the ticket
func (t *ServiceTicket) PutOnHold(reason string) error {
	if t.Status != StatusInProgress {
		return ErrInvalidStatus
	}
	
	t.Status = StatusOnHold
	t.UpdatedAt = time.Now()
	
	return nil
}

// Resume resumes work on a held ticket
func (t *ServiceTicket) Resume() error {
	if t.Status != StatusOnHold {
		return ErrInvalidStatus
	}
	
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
	
	return nil
}

// Resolve marks the ticket as resolved
func (t *ServiceTicket) Resolve(resolutionNotes string, partsUsed []Part, laborHours, cost float64) error {
	if t.Status != StatusInProgress {
		return ErrInvalidStatus
	}
	
	now := time.Now()
	t.ResolutionNotes = resolutionNotes
	t.PartsUsed = partsUsed
	t.LaborHours = laborHours
	t.Cost = cost
	t.ResolvedAt = &now
	t.Status = StatusResolved
	t.UpdatedAt = now
	
	return nil
}

// Close closes the ticket
func (t *ServiceTicket) Close() error {
	if t.Status != StatusResolved {
		return ErrInvalidStatus
	}
	
	now := time.Now()
	t.ClosedAt = &now
	t.Status = StatusClosed
	t.UpdatedAt = now
	
	return nil
}

// Cancel cancels the ticket
func (t *ServiceTicket) Cancel(reason string) error {
	if t.Status == StatusClosed || t.Status == StatusCancelled {
		return ErrInvalidStatus
	}
	
	t.ResolutionNotes = "Cancelled: " + reason
	t.Status = StatusCancelled
	t.UpdatedAt = time.Now()
	
	return nil
}

// SetPriority updates the ticket priority
func (t *ServiceTicket) SetPriority(priority TicketPriority) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

// AddPhoto adds a photo URL to the ticket
func (t *ServiceTicket) AddPhoto(photoURL string) {
	t.Photos = append(t.Photos, photoURL)
	t.UpdatedAt = time.Now()
}

// AddVideo adds a video URL to the ticket
func (t *ServiceTicket) AddVideo(videoURL string) {
	t.Videos = append(t.Videos, videoURL)
	t.UpdatedAt = time.Now()
}

// SetSLA sets SLA deadlines for the ticket
func (t *ServiceTicket) SetSLA(responseHours, resolutionHours int) {
	now := time.Now()
	responseDue := now.Add(time.Duration(responseHours) * time.Hour)
	resolutionDue := now.Add(time.Duration(resolutionHours) * time.Hour)
	
	t.SLAResponseDue = &responseDue
	t.SLAResolutionDue = &resolutionDue
	t.UpdatedAt = now
}

// CheckSLABreach checks if SLA has been breached
func (t *ServiceTicket) CheckSLABreach() bool {
	now := time.Now()
	
	// Check response SLA
	if t.SLAResponseDue != nil && t.AcknowledgedAt == nil && now.After(*t.SLAResponseDue) {
		t.SLABreached = true
		return true
	}
	
	// Check resolution SLA
	if t.SLAResolutionDue != nil && t.ResolvedAt == nil && now.After(*t.SLAResolutionDue) {
		t.SLABreached = true
		return true
	}
	
	return false
}

// GenerateTicketNumber generates a unique ticket number
func GenerateTicketNumber() string {
	now := time.Now()
	return "TKT-" + now.Format("20060102") + "-" + now.Format("150405")
}
