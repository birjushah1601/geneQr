package domain

import (
    "context"
    "encoding/json"
)

// Event types for the service-ticket domain
const (
    EventTicketCreated   = "ticket.created"
    EventTicketAssigned  = "ticket.assigned"
    EventTicketAck       = "ticket.acknowledged"
    EventTicketStarted   = "ticket.started"
    EventTicketOnHold    = "ticket.on_hold"
    EventTicketResumed   = "ticket.resumed"
    EventTicketResolved  = "ticket.resolved"
    EventTicketClosed    = "ticket.closed"
    EventTicketCancelled = "ticket.cancelled"
    EventTicketCommented = "ticket.commented"
)

// EventRepository abstraction to persist events and enqueue deliveries
type EventRepository interface {
    CreateEvent(ctx context.Context, eventType, aggregateType, aggregateID string, payload json.RawMessage) (string, error)
    EnqueueDeliveriesForEvent(ctx context.Context, eventID string, eventType string) error
}
