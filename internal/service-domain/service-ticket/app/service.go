package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	equipmentDomain "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	ticketDomain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/whatsapp"
)

// TicketService provides business logic for service tickets
type TicketService struct {
	repo           ticketDomain.TicketRepository
	equipmentRepo  equipmentDomain.Repository
    policyRepo     ticketDomain.PolicyRepository
    eventRepo      ticketDomain.EventRepository
	logger         *slog.Logger
	defaultSLA     SLAConfig
}

// SLAConfig holds default SLA settings
type SLAConfig struct {
	CriticalResponseHours   int
	CriticalResolutionHours int
	HighResponseHours       int
	HighResolutionHours     int
	MediumResponseHours     int
	MediumResolutionHours   int
	LowResponseHours        int
	LowResolutionHours      int
}

// DefaultSLAConfig returns default SLA configuration
func DefaultSLAConfig() SLAConfig {
	return SLAConfig{
		CriticalResponseHours:   1,
		CriticalResolutionHours: 4,
		HighResponseHours:       2,
		HighResolutionHours:     8,
		MediumResponseHours:     4,
		MediumResolutionHours:   24,
		LowResponseHours:        8,
		LowResolutionHours:      48,
	}
}

// NewTicketService creates a new ticket service
func NewTicketService(
	repo ticketDomain.TicketRepository,
	equipmentRepo equipmentDomain.Repository,
    policyRepo ticketDomain.PolicyRepository,
    eventRepo ticketDomain.EventRepository,
	logger *slog.Logger,
) *TicketService {
	return &TicketService{
		repo:          repo,
		equipmentRepo: equipmentRepo,
        policyRepo:    policyRepo,
        eventRepo:     eventRepo,
		logger:        logger.With(slog.String("component", "ticket_service")),
		defaultSLA:    DefaultSLAConfig(),
	}
}

// CreateTicket creates a new service ticket
func (s *TicketService) CreateTicket(ctx context.Context, req CreateTicketRequest) (*ticketDomain.ServiceTicket, error) {
	s.logger.Info("Creating service ticket",
		slog.String("equipment_id", req.EquipmentID),
		slog.String("customer_name", req.CustomerName))

	// Create ticket
	ticket := ticketDomain.NewServiceTicket(
		req.EquipmentID,
		req.SerialNumber,
		req.EquipmentName,
		req.CustomerName,
		req.IssueDescription,
		req.Source,
		req.CreatedBy,
	)

	ticket.TicketNumber = ticketDomain.GenerateTicketNumber()
	ticket.CustomerID = req.CustomerID
	ticket.CustomerPhone = req.CustomerPhone
	ticket.CustomerWhatsApp = req.CustomerWhatsApp
	ticket.IssueCategory = req.IssueCategory
	ticket.Priority = req.Priority
	ticket.QRCode = req.QRCode
	ticket.SourceMessageID = req.SourceMessageID

	// Add media
	if len(req.Photos) > 0 {
		ticket.Photos = req.Photos
	}
	if len(req.Videos) > 0 {
		ticket.Videos = req.Videos
	}

	// Add requested parts
	if len(req.PartsRequested) > 0 {
		ticket.PartsUsed = req.PartsRequested
	}

    // Set SLA based on policy if available, else defaults
    if s.policyRepo != nil {
        if rules, _ := s.policyRepo.GetSLARules(ctx, nil); rules != nil {
            var resp, res int
            switch ticket.Priority {
            case ticketDomain.PriorityCritical:
                resp, res = rules.Critical.Response, rules.Critical.Resolution
            case ticketDomain.PriorityHigh:
                resp, res = rules.High.Response, rules.High.Resolution
            case ticketDomain.PriorityMedium:
                resp, res = rules.Medium.Response, rules.Medium.Resolution
            case ticketDomain.PriorityLow:
                resp, res = rules.Low.Response, rules.Low.Resolution
            }
            if resp > 0 && res > 0 {
                ticket.SetSLA(resp, res)
            } else {
                s.setSLA(ticket)
            }
        } else {
            s.setSLA(ticket)
        }
    } else {
        s.setSLA(ticket)
    }

	// Save ticket
	if err := s.repo.Create(ctx, ticket); err != nil {
		s.logger.Error("Failed to create ticket", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

    // Emit event: ticket.created
    s.emitEvent(ctx, ticketDomain.EventTicketCreated, "ticket", ticket.ID, map[string]any{
        "ticket_id": ticket.ID,
        "ticket_number": ticket.TicketNumber,
        "priority": ticket.Priority,
    })

	// Optional: minimal responsibility resolver (Phase 4)
	if enabled(os.Getenv("ENABLE_RESP_ORG_ASSIGNMENT")) {
        var resolvedOrg *string
        if s.policyRepo != nil {
            id, _ := s.policyRepo.GetDefaultResponsibleOrg(ctx)
            resolvedOrg = id
        }
        decision := "none"
        reason := "no policy configured"
        if resolvedOrg != nil {
            decision = "default_org"
            reason = "assigned to default_org_id"
        }
        prov := map[string]any{
            "resolver": "policy",
            "decision": decision,
            "reason":   reason,
            "ts":       time.Now().UTC().Format(time.RFC3339),
        }
        b, _ := json.Marshal(prov)
        _ = s.repo.UpdateResponsibility(ctx, ticket.ID, resolvedOrg, b)
	}

	// Add initial comment
	if req.InitialComment != "" {
		comment := &ticketDomain.TicketComment{
			TicketID:    ticket.ID,
			CommentType: "system",
			AuthorName:  "System",
			Comment:     req.InitialComment,
		}
		s.repo.AddComment(ctx, comment)
	}

	s.logger.Info("Ticket created successfully",
		slog.String("ticket_id", ticket.ID),
		slog.String("ticket_number", ticket.TicketNumber))

	return ticket, nil
}

func enabled(v string) bool {
    switch v {
    case "1", "true", "TRUE", "True", "yes", "on":
        return true
    default:
        return false
    }
}

// CreateFromWhatsApp creates a ticket from a WhatsApp message
func (s *TicketService) CreateFromWhatsApp(ctx context.Context, req whatsapp.WhatsAppTicketRequest) (string, error) {
	s.logger.Info("Creating ticket from WhatsApp",
		slog.String("equipment_id", req.EquipmentID),
		slog.String("qr_code", req.QRCode))

	// Look up equipment details
	var equipment *equipmentDomain.Equipment
	var err error

	if req.EquipmentID != "" {
		equipment, err = s.equipmentRepo.GetByID(ctx, req.EquipmentID)
	} else if req.QRCode != "" {
		equipment, err = s.equipmentRepo.GetByQRCode(ctx, req.QRCode)
	} else if req.SerialNumber != "" {
		equipment, err = s.equipmentRepo.GetBySerialNumber(ctx, req.SerialNumber)
	}

	if err != nil {
		s.logger.Error("Failed to find equipment", slog.String("error", err.Error()))
		return "", fmt.Errorf("equipment not found: %w", err)
	}

	// Determine priority based on equipment status
	priority := ticketDomain.PriorityMedium
	if equipment.Status == equipmentDomain.StatusDown {
		priority = ticketDomain.PriorityHigh
	}

	// Check if covered under AMC
	coveredUnderAMC := equipment.HasAMC()

	// Create ticket
	createReq := CreateTicketRequest{
		EquipmentID:      equipment.ID,
		QRCode:           equipment.QRCode,
		SerialNumber:     equipment.SerialNumber,
		EquipmentName:    equipment.EquipmentName,
		CustomerID:       equipment.CustomerID,
		CustomerName:     equipment.CustomerName,
		CustomerPhone:    req.CustomerPhone,
		CustomerWhatsApp: req.CustomerWhatsApp,
		IssueCategory:    "breakdown",
		IssueDescription: req.IssueDescription,
		Priority:         priority,
		Source:           ticketDomain.SourceWhatsApp,
		SourceMessageID:  req.SourceMessageID,
		Photos:           req.Photos,
		Videos:           req.Videos,
		CreatedBy:        "whatsapp-bot",
		InitialComment:   fmt.Sprintf("Ticket created via WhatsApp from customer %s", req.CustomerName),
	}

	ticket, err := s.CreateTicket(ctx, createReq)
	if err != nil {
		return "", err
	}

	// Set AMC coverage
	if coveredUnderAMC {
		ticket.CoveredUnderAMC = true
		ticket.AMCContractID = equipment.AMCContractID
		s.repo.Update(ctx, ticket)
	}

	return ticket.TicketNumber, nil
}

// GetTicket retrieves a ticket by ID
func (s *TicketService) GetTicket(ctx context.Context, id string) (*ticketDomain.ServiceTicket, error) {
	return s.repo.GetByID(ctx, id)
}

// GetTicketByNumber retrieves a ticket by ticket number
func (s *TicketService) GetTicketByNumber(ctx context.Context, ticketNumber string) (*ticketDomain.ServiceTicket, error) {
	return s.repo.GetByTicketNumber(ctx, ticketNumber)
}

// ListTickets lists tickets with filtering and pagination
func (s *TicketService) ListTickets(ctx context.Context, criteria ticketDomain.ListCriteria) (*ticketDomain.TicketListResult, error) {
	return s.repo.List(ctx, criteria)
}

// AssignTicket assigns a ticket to an engineer
func (s *TicketService) AssignTicket(ctx context.Context, ticketID, engineerID, engineerName, assignedBy string) error {
	s.logger.Info("Assigning ticket",
		slog.String("ticket_id", ticketID),
		slog.String("engineer_id", engineerID))

	// Get ticket
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	// Record old status
	oldStatus := string(ticket.Status)

	// Assign engineer
	if err := ticket.AssignEngineer(engineerID, engineerName); err != nil {
		return err
	}

	// Update ticket
	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	// Add status history
	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  assignedBy,
		Reason:     fmt.Sprintf("Assigned to engineer %s", engineerName),
	}
	s.repo.AddStatusHistory(ctx, history)

	// Add comment
	comment := &ticketDomain.TicketComment{
		TicketID:    ticketID,
		CommentType: "system",
		AuthorName:  "System",
		Comment:     fmt.Sprintf("Ticket assigned to engineer %s", engineerName),
	}
	s.repo.AddComment(ctx, comment)

    // Emit event: ticket.assigned
    s.emitEvent(ctx, ticketDomain.EventTicketAssigned, "ticket", ticketID, map[string]any{
        "engineer_id": engineerID,
        "engineer_name": engineerName,
    })
	s.logger.Info("Ticket assigned successfully", slog.String("ticket_id", ticketID))
	return nil
}

// AcknowledgeTicket marks a ticket as acknowledged by the engineer
func (s *TicketService) AcknowledgeTicket(ctx context.Context, ticketID, acknowledgedBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if err := ticket.Acknowledge(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	comment := &ticketDomain.TicketComment{
		TicketID:    ticketID,
		CommentType: "engineer",
		AuthorID:    ticket.AssignedEngineerID,
		AuthorName:  ticket.AssignedEngineerName,
		Comment:     "Ticket acknowledged. Will start work soon.",
	}
	s.repo.AddComment(ctx, comment)

    // Emit event: ticket.acknowledged
    s.emitEvent(ctx, ticketDomain.EventTicketAck, "ticket", ticketID, map[string]any{})
	return nil
}

// StartWork marks a ticket as in progress
func (s *TicketService) StartWork(ctx context.Context, ticketID, startedBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.Start(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  startedBy,
		Reason:     "Work started on ticket",
	}
	s.repo.AddStatusHistory(ctx, history)

	comment := &ticketDomain.TicketComment{
		TicketID:    ticketID,
		CommentType: "engineer",
		AuthorID:    ticket.AssignedEngineerID,
		AuthorName:  ticket.AssignedEngineerName,
		Comment:     "Started working on the issue.",
	}
	s.repo.AddComment(ctx, comment)

    // Emit event: ticket.started
    s.emitEvent(ctx, ticketDomain.EventTicketStarted, "ticket", ticketID, map[string]any{})
	return nil
}

// PutOnHold puts a ticket on hold
func (s *TicketService) PutOnHold(ctx context.Context, ticketID, reason, changedBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.PutOnHold(reason); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  changedBy,
		Reason:     reason,
	}
	s.repo.AddStatusHistory(ctx, history)

	comment := &ticketDomain.TicketComment{
		TicketID:    ticketID,
		CommentType: "engineer",
		AuthorID:    ticket.AssignedEngineerID,
		AuthorName:  ticket.AssignedEngineerName,
		Comment:     fmt.Sprintf("Ticket put on hold: %s", reason),
	}
	s.repo.AddComment(ctx, comment)

    // Emit event: ticket.on_hold
    s.emitEvent(ctx, ticketDomain.EventTicketOnHold, "ticket", ticketID, map[string]any{"reason": reason})
	return nil
}

// ResumeWork resumes work on a held ticket
func (s *TicketService) ResumeWork(ctx context.Context, ticketID, resumedBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.Resume(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  resumedBy,
		Reason:     "Work resumed",
	}
	s.repo.AddStatusHistory(ctx, history)

    // Emit event: ticket.resumed
    s.emitEvent(ctx, ticketDomain.EventTicketResumed, "ticket", ticketID, map[string]any{})
	return nil
}

// ResolveTicket marks a ticket as resolved
func (s *TicketService) ResolveTicket(ctx context.Context, ticketID string, req ResolveTicketRequest) error {
	s.logger.Info("Resolving ticket", slog.String("ticket_id", ticketID))

	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.Resolve(req.ResolutionNotes, req.PartsUsed, req.LaborHours, req.Cost); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  req.ResolvedBy,
		Reason:     "Ticket resolved",
	}
	s.repo.AddStatusHistory(ctx, history)

	comment := &ticketDomain.TicketComment{
		TicketID:    ticketID,
		CommentType: "engineer",
		AuthorID:    ticket.AssignedEngineerID,
		AuthorName:  ticket.AssignedEngineerName,
		Comment:     fmt.Sprintf("Issue resolved. %s", req.ResolutionNotes),
	}
	s.repo.AddComment(ctx, comment)

	// Update equipment service history
	if s.equipmentRepo != nil && ticket.EquipmentID != "" {
		// This would record service in equipment registry
		// s.equipmentRepo.RecordService(...)
	}

	s.logger.Info("Ticket resolved successfully", slog.String("ticket_id", ticketID))

    // Emit event: ticket.resolved
    s.emitEvent(ctx, ticketDomain.EventTicketResolved, "ticket", ticketID, map[string]any{"notes": req.ResolutionNotes})
	return nil
}

// CloseTicket closes a resolved ticket
func (s *TicketService) CloseTicket(ctx context.Context, ticketID, closedBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.Close(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  closedBy,
		Reason:     "Ticket closed",
	}
	s.repo.AddStatusHistory(ctx, history)

    // Emit event: ticket.closed
    s.emitEvent(ctx, ticketDomain.EventTicketClosed, "ticket", ticketID, map[string]any{})
	return nil
}

// CancelTicket cancels a ticket
func (s *TicketService) CancelTicket(ctx context.Context, ticketID, reason, cancelledBy string) error {
	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldStatus := string(ticket.Status)

	if err := ticket.Cancel(reason); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, ticket); err != nil {
		return err
	}

	history := &ticketDomain.StatusHistory{
		TicketID:   ticketID,
		FromStatus: oldStatus,
		ToStatus:   string(ticket.Status),
		ChangedBy:  cancelledBy,
		Reason:     reason,
	}
	s.repo.AddStatusHistory(ctx, history)

    // Emit event: ticket.cancelled
    s.emitEvent(ctx, ticketDomain.EventTicketCancelled, "ticket", ticketID, map[string]any{"reason": reason})
	return nil
}

// AddComment adds a comment to a ticket
func (s *TicketService) AddComment(ctx context.Context, req AddCommentRequest) error {
	// Set defaults for empty fields
	authorID := req.AuthorID
	if authorID == "" {
		authorID = "system"
	}
	
	authorName := req.AuthorName
	if authorName == "" {
		authorName = "System User"
	}
	
	attachments := req.Attachments
	if attachments == nil {
		attachments = []string{}
	}
	
	comment := &ticketDomain.TicketComment{
		TicketID:    req.TicketID,
		CommentType: req.CommentType,
		AuthorID:    authorID,
		AuthorName:  authorName,
		Comment:     req.Comment,
		Attachments: attachments,
	}

	return s.repo.AddComment(ctx, comment)
}

// GetComments retrieves all comments for a ticket
func (s *TicketService) GetComments(ctx context.Context, ticketID string) ([]*ticketDomain.TicketComment, error) {
	return s.repo.GetComments(ctx, ticketID)
}

// GetStatusHistory retrieves status history for a ticket
func (s *TicketService) GetStatusHistory(ctx context.Context, ticketID string) ([]*ticketDomain.StatusHistory, error) {
	return s.repo.GetStatusHistory(ctx, ticketID)
}

// UpdateParts updates the parts assigned to a ticket
func (s *TicketService) UpdateParts(ctx context.Context, ticketID string, parts []map[string]interface{}) error {
	s.logger.Info("Updating ticket parts", slog.String("ticket_id", ticketID))

	ticket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	// Update parts used - convert to []interface{} which is JSON compatible
	partsInterface := make([]interface{}, len(parts))
	for i, p := range parts {
		partsInterface[i] = p
	}
	
	// Store as interface slice (will be marshaled to JSONB)
	ticket.PartsUsed = partsInterface

	if err := s.repo.Update(ctx, ticket); err != nil {
		s.logger.Error("Failed to update ticket parts", slog.String("error", err.Error()))
		return fmt.Errorf("failed to update ticket parts: %w", err)
	}

	// Emit event: ticket.parts_updated
	s.emitEvent(ctx, "ticket.parts_updated", "ticket", ticketID, map[string]any{
		"parts_count": len(parts),
	})

	s.logger.Info("Ticket parts updated successfully", slog.String("ticket_id", ticketID))
	return nil
}

// emitEvent is a best-effort outbox writer (no-op if repo is nil)
func (s *TicketService) emitEvent(ctx context.Context, eventType, aggregateType, aggregateID string, payload map[string]any) {
    if s.eventRepo == nil { return }
    b, _ := json.Marshal(payload)
    if id, err := s.eventRepo.CreateEvent(ctx, eventType, aggregateType, aggregateID, b); err == nil {
        _ = s.eventRepo.EnqueueDeliveriesForEvent(ctx, id, eventType)
    }
}

// setSLA sets SLA deadlines based on priority
func (s *TicketService) setSLA(ticket *ticketDomain.ServiceTicket) {
	var responseHours, resolutionHours int

	switch ticket.Priority {
	case ticketDomain.PriorityCritical:
		responseHours = s.defaultSLA.CriticalResponseHours
		resolutionHours = s.defaultSLA.CriticalResolutionHours
	case ticketDomain.PriorityHigh:
		responseHours = s.defaultSLA.HighResponseHours
		resolutionHours = s.defaultSLA.HighResolutionHours
	case ticketDomain.PriorityMedium:
		responseHours = s.defaultSLA.MediumResponseHours
		resolutionHours = s.defaultSLA.MediumResolutionHours
	case ticketDomain.PriorityLow:
		responseHours = s.defaultSLA.LowResponseHours
		resolutionHours = s.defaultSLA.LowResolutionHours
	default:
		responseHours = s.defaultSLA.MediumResponseHours
		resolutionHours = s.defaultSLA.MediumResolutionHours
	}

	ticket.SetSLA(responseHours, resolutionHours)
}

// Request/Response DTOs

type CreateTicketRequest struct {
	EquipmentID      string
	QRCode           string
	SerialNumber     string
	EquipmentName    string
	CustomerID       string
	CustomerName     string
	CustomerPhone    string
	CustomerWhatsApp string
	IssueCategory    string
	IssueDescription string
	Priority         ticketDomain.TicketPriority
	Source           ticketDomain.TicketSource
	SourceMessageID  string
	Photos           []string
	Videos           []string
	CreatedBy        string
	InitialComment   string
	PartsRequested   []ticketDomain.Part // Parts requested for this service
}

type ResolveTicketRequest struct {
	ResolutionNotes string
	PartsUsed       []ticketDomain.Part
	LaborHours      float64
	Cost            float64
	ResolvedBy      string
}

type AddCommentRequest struct {
	TicketID    string   `json:"ticket_id"`
	CommentType string   `json:"comment_type"`
	AuthorID    string   `json:"author_id"`
	AuthorName  string   `json:"author_name"`
	Comment     string   `json:"comment"`
	Attachments []string   `json:"attachments"`
}


// DeleteComment deletes a comment from a ticket
func (s *TicketService) DeleteComment(ctx context.Context, ticketID string, commentID string) error {
	// Validate ticket exists
	_, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}
	
	// Delete the comment
	if err := s.repo.DeleteComment(ctx, commentID, ticketID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	
	return nil
}
