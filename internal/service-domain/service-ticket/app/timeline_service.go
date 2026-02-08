package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ticketDomain "github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

type TimelineService struct {
	ticketRepo ticketDomain.TicketRepository
	logger     *slog.Logger
}

func NewTimelineService(ticketRepo ticketDomain.TicketRepository, logger *slog.Logger) *TimelineService {
	return &TimelineService{
		ticketRepo: ticketRepo,
		logger:     logger.With(slog.String("service", "timeline")),
	}
}

// GenerateTimeline creates a timeline based on current ticket state
func (s *TimelineService) GenerateTimeline(ctx context.Context, ticket *ticketDomain.ServiceTicket) (*ticketDomain.TicketTimeline, error) {
	slaConfig := ticketDomain.GetSLAConfig(ticket.Priority)
	
	timeline := &ticketDomain.TicketTimeline{
		TicketID:      ticket.ID,
		CurrentStage:  s.determineCurrentStage(ticket),
		RequiresParts: s.requiresParts(ticket),
		LastUpdated:   time.Now(),
	}

	// Build milestones based on workflow
	milestones := s.buildMilestones(ticket, slaConfig, timeline.RequiresParts)
	timeline.Milestones = milestones

	// Calculate overall status and ETA
	timeline.OverallStatus = s.calculateOverallStatus(milestones)
	timeline.EstimatedResolution = s.calculateFinalETA(milestones)

	return timeline, nil
}

// ConvertToPublicTimeline converts internal timeline to customer-facing view
func (s *TimelineService) ConvertToPublicTimeline(timeline *ticketDomain.TicketTimeline, ticket *ticketDomain.ServiceTicket) *ticketDomain.PublicTimeline {
	publicTimeline := &ticketDomain.PublicTimeline{
		OverallStatus:       timeline.OverallStatus,
		StatusMessage:       s.generateStatusMessage(timeline, ticket),
		CurrentStage:        string(timeline.CurrentStage),
		CurrentStageDesc:    s.getStageDescription(timeline.CurrentStage),
		EstimatedResolution: timeline.EstimatedResolution,
		TimeRemaining:       s.formatTimeRemaining(timeline.EstimatedResolution),
		RequiresParts:       timeline.RequiresParts,
		AssignedEngineer:    ticket.AssignedEngineerName,
		Priority:            string(ticket.Priority),
		IsUrgent:            ticket.Priority == ticketDomain.PriorityCritical || ticket.Priority == ticketDomain.PriorityHigh,
		ProgressPercentage:  s.calculateProgress(timeline.Milestones),
	}

	// Add parts info if applicable
	if timeline.RequiresParts {
		publicTimeline.PartsStatus = s.getPartsStatus(timeline.Milestones)
		publicTimeline.PartsETA = timeline.PartsETA
	}

	// Convert milestones to public view
	publicTimeline.Milestones = s.convertMilestonesToPublic(timeline.Milestones)

	// Determine next stage
	if idx := s.findCurrentMilestoneIndex(timeline.Milestones); idx < len(timeline.Milestones)-1 {
		nextMilestone := timeline.Milestones[idx+1]
		publicTimeline.NextStage = string(nextMilestone.Stage)
		publicTimeline.NextStageDesc = s.getStageDescription(nextMilestone.Stage)
	}

	return publicTimeline
}

// determineCurrentStage figures out what stage the ticket is in
func (s *TimelineService) determineCurrentStage(ticket *ticketDomain.ServiceTicket) ticketDomain.MilestoneStage {
	switch ticket.Status {
	case ticketDomain.StatusNew:
		return ticketDomain.MilestoneAcknowledgment
	case ticketDomain.StatusAssigned:
		if ticket.AcknowledgedAt != nil {
			return ticketDomain.MilestoneDiagnosis
		}
		return ticketDomain.MilestoneAcknowledgment
	case ticketDomain.StatusInProgress:
		// Check if parts are needed
		// TODO: Check if parts are assigned to ticket
		return ticketDomain.MilestoneRepairStart
	case ticketDomain.StatusOnHold:
		// Usually on hold means waiting for parts
		return ticketDomain.MilestonePartsDelivery
	case ticketDomain.StatusResolved:
		return ticketDomain.MilestoneVerification
	case ticketDomain.StatusClosed:
		return ticketDomain.MilestoneResolution
	default:
		return ticketDomain.MilestoneAcknowledgment
	}
}

// requiresParts checks if ticket needs parts (stub for now)
func (s *TimelineService) requiresParts(ticket *ticketDomain.ServiceTicket) bool {
	// TODO: Check if parts are assigned to ticket from database
	// For now, check status
	return ticket.Status == ticketDomain.StatusOnHold
}

// buildMilestones creates milestone list based on workflow
func (s *TimelineService) buildMilestones(ticket *ticketDomain.ServiceTicket, config ticketDomain.SLAConfig, needsParts bool) []ticketDomain.TicketMilestone {
	milestones := []ticketDomain.TicketMilestone{}
	baseTime := ticket.CreatedAt

	// 1. Acknowledgment
	ackETA := baseTime.Add(time.Duration(config.ResponseHours) * time.Hour)
	milestones = append(milestones, ticketDomain.TicketMilestone{
		Stage:             ticketDomain.MilestoneAcknowledgment,
		Status:            s.getMilestoneStatus(ticket, ticketDomain.MilestoneAcknowledgment),
		EstimatedComplete: &ackETA,
		ActualComplete:    ticket.AcknowledgedAt,
		Description:       "Engineer acknowledges ticket and reviews details",
		AssignedTo:        ticket.AssignedEngineerName,
	})

	// 2. Diagnosis
	diagnosisStart := ackETA
	if ticket.AcknowledgedAt != nil {
		diagnosisStart = *ticket.AcknowledgedAt
	}
	diagnosisETA := diagnosisStart.Add(time.Duration(config.DiagnosisHours) * time.Hour)
	milestones = append(milestones, ticketDomain.TicketMilestone{
		Stage:             ticketDomain.MilestoneDiagnosis,
		Status:            s.getMilestoneStatus(ticket, ticketDomain.MilestoneDiagnosis),
		EstimatedStart:    &diagnosisStart,
		EstimatedComplete: &diagnosisETA,
		Description:       "Engineer diagnoses the issue and determines solution",
		AssignedTo:        ticket.AssignedEngineerName,
	})

	// 3-6. Parts workflow (if needed)
	if needsParts {
		// Parts ordering
		partsOrderStart := diagnosisETA
		partsOrderETA := partsOrderStart.Add(time.Duration(config.PartsOrderHours) * time.Hour)
		milestones = append(milestones, ticketDomain.TicketMilestone{
			Stage:             ticketDomain.MilestonePartsOrdered,
			Status:            ticketDomain.MilestoneStatusPending,
			EstimatedStart:    &partsOrderStart,
			EstimatedComplete: &partsOrderETA,
			Description:       "Required parts ordered from supplier",
		})

		// Parts delivery
		deliveryDays := config.StandardPartsDelivery
		if ticket.Priority == ticketDomain.PriorityCritical {
			deliveryDays = config.UrgentPartsDelivery
		}
		partsDeliveryETA := partsOrderETA.Add(time.Duration(deliveryDays) * 24 * time.Hour)
		milestones = append(milestones, ticketDomain.TicketMilestone{
			Stage:             ticketDomain.MilestonePartsDelivery,
			Status:            ticketDomain.MilestoneStatusBlocked,
			EstimatedComplete: &partsDeliveryETA,
			Description:       fmt.Sprintf("Waiting for parts delivery (%d business days)", deliveryDays),
			BlockerReason:     "Waiting for parts from supplier",
		})

		// Parts received
		milestones = append(milestones, ticketDomain.TicketMilestone{
			Stage:             ticketDomain.MilestonePartsReceived,
			Status:            ticketDomain.MilestoneStatusPending,
			EstimatedComplete: &partsDeliveryETA,
			Description:       "Parts received and ready for installation",
		})

		// Repair after parts
		repairStart := partsDeliveryETA
		repairETA := repairStart.Add(time.Duration(config.RepairAfterPartsHours) * time.Hour)
		milestones = append(milestones, ticketDomain.TicketMilestone{
			Stage:             ticketDomain.MilestoneRepairStart,
			Status:            ticketDomain.MilestoneStatusPending,
			EstimatedStart:    &repairStart,
			EstimatedComplete: &repairETA,
			Description:       "Engineer installs parts and completes repair",
			AssignedTo:        ticket.AssignedEngineerName,
		})
	} else {
		// Simple repair without parts
		repairStart := diagnosisETA
		repairETA := repairStart.Add(time.Duration(config.SimpleRepairHours) * time.Hour)
		milestones = append(milestones, ticketDomain.TicketMilestone{
			Stage:             ticketDomain.MilestoneRepairStart,
			Status:            s.getMilestoneStatus(ticket, ticketDomain.MilestoneRepairStart),
			EstimatedStart:    &repairStart,
			EstimatedComplete: &repairETA,
			Description:       "Engineer performs repair",
			AssignedTo:        ticket.AssignedEngineerName,
		})
	}

	// 7. Verification
	lastMilestone := milestones[len(milestones)-1]
	verifyStart := *lastMilestone.EstimatedComplete
	verifyETA := verifyStart.Add(time.Duration(config.VerificationHours) * time.Hour)
	milestones = append(milestones, ticketDomain.TicketMilestone{
		Stage:             ticketDomain.MilestoneVerification,
		Status:            s.getMilestoneStatus(ticket, ticketDomain.MilestoneVerification),
		EstimatedStart:    &verifyStart,
		EstimatedComplete: &verifyETA,
		Description:       "Testing and verification of the fix",
		AssignedTo:        ticket.AssignedEngineerName,
	})

	// 8. Final resolution
	milestones = append(milestones, ticketDomain.TicketMilestone{
		Stage:             ticketDomain.MilestoneResolution,
		Status:            s.getMilestoneStatus(ticket, ticketDomain.MilestoneResolution),
		EstimatedComplete: &verifyETA,
		Description:       "Ticket closed and resolved",
	})

	return milestones
}

// getMilestoneStatus determines status based on ticket state
func (s *TimelineService) getMilestoneStatus(ticket *ticketDomain.ServiceTicket, stage ticketDomain.MilestoneStage) ticketDomain.MilestoneStatus {
	currentStage := s.determineCurrentStage(ticket)
	
	// Compare stages to determine status
	stageOrder := map[ticketDomain.MilestoneStage]int{
		ticketDomain.MilestoneAcknowledgment: 0,
		ticketDomain.MilestoneDiagnosis:      1,
		ticketDomain.MilestonePartsOrdered:   2,
		ticketDomain.MilestonePartsDelivery:  3,
		ticketDomain.MilestonePartsReceived:  4,
		ticketDomain.MilestoneRepairStart:    5,
		ticketDomain.MilestoneVerification:   6,
		ticketDomain.MilestoneResolution:     7,
	}

	currentOrder := stageOrder[currentStage]
	stageOrderNum := stageOrder[stage]

	if stageOrderNum < currentOrder {
		return ticketDomain.MilestoneStatusCompleted
	} else if stageOrderNum == currentOrder {
		return ticketDomain.MilestoneStatusInProgress
	}
	return ticketDomain.MilestoneStatusPending
}

// Helper functions
func (s *TimelineService) calculateOverallStatus(milestones []ticketDomain.TicketMilestone) string {
	// Check if any milestone is delayed or blocked
	for _, m := range milestones {
		if m.Status == ticketDomain.MilestoneStatusDelayed {
			return "delayed"
		}
		if m.Status == ticketDomain.MilestoneStatusBlocked {
			return "blocked"
		}
	}
	return "on_track"
}

func (s *TimelineService) calculateFinalETA(milestones []ticketDomain.TicketMilestone) *time.Time {
	if len(milestones) == 0 {
		return nil
	}
	return milestones[len(milestones)-1].EstimatedComplete
}

func (s *TimelineService) generateStatusMessage(timeline *ticketDomain.TicketTimeline, ticket *ticketDomain.ServiceTicket) string {
	switch timeline.OverallStatus {
	case "on_track":
		if timeline.RequiresParts {
			return fmt.Sprintf("Your service request is progressing well. We've identified the parts needed and are working to complete the repair by %s.", 
				s.formatDateTime(timeline.EstimatedResolution))
		}
		return fmt.Sprintf("Your service request is on track for completion by %s.", s.formatDateTime(timeline.EstimatedResolution))
	case "blocked":
		if timeline.RequiresParts {
			return fmt.Sprintf("We're waiting for parts to arrive (expected: %s). Once received, our engineer will complete the repair promptly.", 
				s.formatDateTime(timeline.PartsETA))
		}
		return "Your service request is temporarily on hold. We're working to resolve the blocker."
	case "delayed":
		return fmt.Sprintf("Your service request is taking longer than initially expected. Revised completion estimate: %s. We apologize for the delay.", 
			s.formatDateTime(timeline.EstimatedResolution))
	default:
		return "Your service request is being processed."
	}
}

func (s *TimelineService) getStageDescription(stage ticketDomain.MilestoneStage) string {
	descriptions := map[ticketDomain.MilestoneStage]string{
		ticketDomain.MilestoneAcknowledgment: "Engineer reviewing your request",
		ticketDomain.MilestoneDiagnosis:      "Diagnosing the issue",
		ticketDomain.MilestonePartsOrdered:   "Ordering required parts",
		ticketDomain.MilestonePartsDelivery:  "Waiting for parts delivery",
		ticketDomain.MilestonePartsReceived:  "Parts received",
		ticketDomain.MilestoneRepairSchedule: "Scheduling repair visit",
		ticketDomain.MilestoneRepairStart:    "Performing repair",
		ticketDomain.MilestoneVerification:   "Testing and verification",
		ticketDomain.MilestoneResolution:     "Service completed",
	}
	return descriptions[stage]
}

func (s *TimelineService) formatTimeRemaining(eta *time.Time) string {
	if eta == nil {
		return ""
	}
	
	remaining := time.Until(*eta)
	if remaining < 0 {
		return "Overdue"
	}
	
	days := int(remaining.Hours() / 24)
	hours := int(remaining.Hours()) % 24
	
	if days > 0 {
		return fmt.Sprintf("%d days, %d hours", days, hours)
	}
	return fmt.Sprintf("%d hours", hours)
}

func (s *TimelineService) formatDateTime(t *time.Time) string {
	if t == nil {
		return "TBD"
	}
	return t.Format("Jan 2, 2006 at 3:04 PM")
}

func (s *TimelineService) getPartsStatus(milestones []ticketDomain.TicketMilestone) string {
	for _, m := range milestones {
		switch m.Stage {
		case ticketDomain.MilestonePartsOrdered:
			if m.Status == ticketDomain.MilestoneStatusCompleted {
				continue
			}
			return "ordering"
		case ticketDomain.MilestonePartsDelivery:
			if m.Status == ticketDomain.MilestoneStatusCompleted {
				continue
			}
			return "in_transit"
		case ticketDomain.MilestonePartsReceived:
			if m.Status == ticketDomain.MilestoneStatusCompleted {
				return "received"
			}
		}
	}
	return "unknown"
}

func (s *TimelineService) calculateProgress(milestones []ticketDomain.TicketMilestone) int {
	if len(milestones) == 0 {
		return 0
	}
	
	completed := 0
	for _, m := range milestones {
		if m.Status == ticketDomain.MilestoneStatusCompleted {
			completed++
		}
	}
	
	return (completed * 100) / len(milestones)
}

func (s *TimelineService) convertMilestonesToPublic(milestones []ticketDomain.TicketMilestone) []ticketDomain.PublicMilestone {
	public := make([]ticketDomain.PublicMilestone, len(milestones))
	
	for i, m := range milestones {
		public[i] = ticketDomain.PublicMilestone{
			Stage:       string(m.Stage),
			Title:       s.getMilestoneTitle(m.Stage),
			Description: m.Description,
			Status:      string(m.Status),
			ETA:         m.EstimatedComplete,
			CompletedAt: m.ActualComplete,
			IsActive:    m.Status == ticketDomain.MilestoneStatusInProgress,
		}
	}
	
	return public
}

func (s *TimelineService) getMilestoneTitle(stage ticketDomain.MilestoneStage) string {
	titles := map[ticketDomain.MilestoneStage]string{
		ticketDomain.MilestoneAcknowledgment: "Acknowledgment",
		ticketDomain.MilestoneDiagnosis:      "Diagnosis",
		ticketDomain.MilestonePartsOrdered:   "Parts Ordered",
		ticketDomain.MilestonePartsDelivery:  "Parts Delivery",
		ticketDomain.MilestonePartsReceived:  "Parts Received",
		ticketDomain.MilestoneRepairSchedule: "Repair Scheduled",
		ticketDomain.MilestoneRepairStart:    "Repair",
		ticketDomain.MilestoneVerification:   "Verification",
		ticketDomain.MilestoneResolution:     "Completed",
	}
	return titles[stage]
}

func (s *TimelineService) findCurrentMilestoneIndex(milestones []ticketDomain.TicketMilestone) int {
	for i, m := range milestones {
		if m.Status == ticketDomain.MilestoneStatusInProgress {
			return i
		}
	}
	return 0
}
