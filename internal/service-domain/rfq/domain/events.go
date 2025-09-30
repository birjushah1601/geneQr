package domain

import (
	"time"
)

// EventType represents the type of domain event
type EventType string

const (
	EventTypeRFQCreated    EventType = "rfq.created"
	EventTypeRFQUpdated    EventType = "rfq.updated"
	EventTypeRFQPublished  EventType = "rfq.published"
	EventTypeRFQClosed     EventType = "rfq.closed"
	EventTypeRFQCancelled  EventType = "rfq.cancelled"
	EventTypeRFQAwarded    EventType = "rfq.awarded"
	EventTypeSupplierInvited EventType = "rfq.supplier_invited"
)

// DomainEvent is the base structure for all domain events
type DomainEvent struct {
	EventID   string    `json:"event_id"`
	EventType EventType `json:"event_type"`
	TenantID  string    `json:"tenant_id"`
	Timestamp time.Time `json:"timestamp"`
}

// RFQCreatedEvent is published when a new RFQ is created
type RFQCreatedEvent struct {
	DomainEvent
	RFQID      string      `json:"rfq_id"`
	RFQNumber  string      `json:"rfq_number"`
	Title      string      `json:"title"`
	Priority   RFQPriority `json:"priority"`
	CreatedBy  string      `json:\"created_by\"`
}

// RFQPublishedEvent is published when an RFQ is published
type RFQPublishedEvent struct {
	DomainEvent
	RFQID            string    `json:"rfq_id"`
	RFQNumber        string    `json:"rfq_number"`
	Title            string    `json:"title"`
	ResponseDeadline time.Time `json:"response_deadline"`
	ItemCount        int       `json:"item_count"`
}

// RFQClosedEvent is published when an RFQ is closed
type RFQClosedEvent struct {
	DomainEvent
	RFQID     string `json:"rfq_id"`
	RFQNumber string `json:"rfq_number"`
	ClosedBy  string `json:"closed_by"`
}

// RFQCancelledEvent is published when an RFQ is cancelled
type RFQCancelledEvent struct {
	DomainEvent
	RFQID       string `json:"rfq_id"`
	RFQNumber   string `json:"rfq_number"`
	CancelledBy string `json:"cancelled_by"`
	Reason      string `json:"reason,omitempty"`
}

// RFQAwardedEvent is published when an RFQ is awarded
type RFQAwardedEvent struct {
	DomainEvent
	RFQID      string `json:"rfq_id"`
	RFQNumber  string `json:"rfq_number"`
	SupplierID string `json:"supplier_id"`
	AwardedBy  string `json:"awarded_by"`
}

// SupplierInvitedEvent is published when a supplier is invited to quote
type SupplierInvitedEvent struct {
	DomainEvent
	RFQID        string    `json:"rfq_id"`
	RFQNumber    string    `json:"rfq_number"`
	SupplierID   string    `json:"supplier_id"`
	InvitationID string    `json:"invitation_id"`
	Deadline     time.Time `json:"deadline"`
}

// NewRFQCreatedEvent creates a new RFQ created event
func NewRFQCreatedEvent(rfq *RFQ) *RFQCreatedEvent {
	return &RFQCreatedEvent{
		DomainEvent: DomainEvent{
			EventID:   generateEventID(),
			EventType: EventTypeRFQCreated,
			TenantID:  rfq.TenantID,
			Timestamp: time.Now(),
		},
		RFQID:     rfq.ID,
		RFQNumber: rfq.RFQNumber,
		Title:     rfq.Title,
		Priority:  rfq.Priority,
		CreatedBy: rfq.CreatedBy,
	}
}

// NewRFQPublishedEvent creates a new RFQ published event
func NewRFQPublishedEvent(rfq *RFQ) *RFQPublishedEvent {
	return &RFQPublishedEvent{
		DomainEvent: DomainEvent{
			EventID:   generateEventID(),
			EventType: EventTypeRFQPublished,
			TenantID:  rfq.TenantID,
			Timestamp: time.Now(),
		},
		RFQID:            rfq.ID,
		RFQNumber:        rfq.RFQNumber,
		Title:            rfq.Title,
		ResponseDeadline: rfq.ResponseDeadline,
		ItemCount:        len(rfq.Items),
	}
}

// NewRFQClosedEvent creates a new RFQ closed event
func NewRFQClosedEvent(rfq *RFQ, closedBy string) *RFQClosedEvent {
	return &RFQClosedEvent{
		DomainEvent: DomainEvent{
			EventID:   generateEventID(),
			EventType: EventTypeRFQClosed,
			TenantID:  rfq.TenantID,
			Timestamp: time.Now(),
		},
		RFQID:     rfq.ID,
		RFQNumber: rfq.RFQNumber,
		ClosedBy:  closedBy,
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	// In production, use a proper ID generation library (e.g., ULID, UUID)
	return time.Now().Format("20060102150405") + "-" + "event"
}
