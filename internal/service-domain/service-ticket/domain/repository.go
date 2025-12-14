package domain

import (
    "context"
    "encoding/json"
)

// TicketRepository defines data access operations for service tickets
type TicketRepository interface {
	// Create creates a new service ticket
	Create(ctx context.Context, ticket *ServiceTicket) error
	
	// GetByID retrieves a ticket by ID
	GetByID(ctx context.Context, id string) (*ServiceTicket, error)
	
	// GetByTicketNumber retrieves a ticket by ticket number
	GetByTicketNumber(ctx context.Context, ticketNumber string) (*ServiceTicket, error)
	
	// Update updates an existing ticket
	Update(ctx context.Context, ticket *ServiceTicket) error
	
	// List retrieves tickets based on criteria
	List(ctx context.Context, criteria ListCriteria) (*TicketListResult, error)
	
	// GetByEquipment retrieves all tickets for an equipment
	GetByEquipment(ctx context.Context, equipmentID string) ([]*ServiceTicket, error)
	
	// GetByCustomer retrieves all tickets for a customer
	GetByCustomer(ctx context.Context, customerID string) ([]*ServiceTicket, error)
	
	// GetByEngineer retrieves all tickets assigned to an engineer
	GetByEngineer(ctx context.Context, engineerID string) ([]*ServiceTicket, error)
	
	// GetBySource retrieves tickets by source
	GetBySource(ctx context.Context, source TicketSource) ([]*ServiceTicket, error)
	
	// AddComment adds a comment to a ticket
	AddComment(ctx context.Context, comment *TicketComment) error
	
	// GetComments retrieves all comments for a ticket
	GetComments(ctx context.Context, ticketID string) ([]*TicketComment, error)
	DeleteComment(ctx context.Context, commentID string, ticketID string) error
	
	// AddStatusHistory records a status change
	AddStatusHistory(ctx context.Context, history *StatusHistory) error
	
	// GetStatusHistory retrieves status history for a ticket
	GetStatusHistory(ctx context.Context, ticketID string) ([]*StatusHistory, error)

    // UpdateResponsibility sets responsible_org_id and policy_provenance (Phase 4 optional)
    UpdateResponsibility(ctx context.Context, ticketID string, responsibleOrgID *string, provenance json.RawMessage) error
}

// ListCriteria defines filtering criteria for listing tickets
type ListCriteria struct {
	Status           []TicketStatus
	Priority         []TicketPriority
	Source           []TicketSource
	EquipmentID      string
	CustomerID       string
	EngineerID       string
	SLABreached      *bool
	CoveredUnderAMC  *bool
	CreatedAfter     *string
	CreatedBefore    *string
	SortBy           string
	SortDirection    string
	Page             int
	PageSize         int
}

// TicketListResult contains paginated ticket results
type TicketListResult struct {
	Tickets    []*ServiceTicket `json:"tickets"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// TicketComment represents a comment on a ticket
type TicketComment struct {
	ID          string    `json:"id"`
	TicketID    string    `json:"ticket_id"`
	CommentType string    `json:"comment_type"` // customer, engineer, internal, system
	AuthorID    string    `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	Comment     string    `json:"comment"`
	Attachments []string  `json:"attachments"`
	CreatedAt   string    `json:"created_at"`
}

// StatusHistory tracks status changes
type StatusHistory struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  string    `json:"changed_at"`
	Reason     string    `json:"reason,omitempty"`
}
