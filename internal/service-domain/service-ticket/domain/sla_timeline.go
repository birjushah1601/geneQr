package domain

import "time"

// TicketMilestone represents a stage in the ticket lifecycle with its own ETA
type TicketMilestone struct {
	Stage             MilestoneStage `json:"stage"`
	Status            MilestoneStatus `json:"status"` // pending, in_progress, completed, delayed, blocked
	EstimatedStart    *time.Time     `json:"estimated_start,omitempty"`
	EstimatedComplete *time.Time     `json:"estimated_complete,omitempty"`
	ActualStart       *time.Time     `json:"actual_start,omitempty"`
	ActualComplete    *time.Time     `json:"actual_complete,omitempty"`
	Description       string         `json:"description"`
	BlockerReason     string         `json:"blocker_reason,omitempty"`
	AssignedTo        string         `json:"assigned_to,omitempty"` // Engineer name
}

type MilestoneStage string

const (
	// Standard workflow stages
	MilestoneAcknowledgment MilestoneStage = "acknowledgment"  // Engineer acknowledges ticket
	MilestoneDiagnosis      MilestoneStage = "diagnosis"       // Engineer diagnoses issue
	MilestonePartsOrdered   MilestoneStage = "parts_ordered"   // Parts ordered from supplier
	MilestonePartsDelivery  MilestoneStage = "parts_delivery"  // Waiting for parts to arrive
	MilestonePartsReceived  MilestoneStage = "parts_received"  // Parts arrived
	MilestoneRepairSchedule MilestoneStage = "repair_schedule" // Scheduling engineer visit
	MilestoneRepairStart    MilestoneStage = "repair_start"    // Engineer on-site
	MilestoneRepairComplete MilestoneStage = "repair_complete" // Fix applied
	MilestoneVerification   MilestoneStage = "verification"    // Testing/verification
	MilestoneResolution     MilestoneStage = "resolution"      // Final closure
)

type MilestoneStatus string

const (
	MilestoneStatusPending    MilestoneStatus = "pending"     // Not started yet
	MilestoneStatusInProgress MilestoneStatus = "in_progress" // Currently working on this
	MilestoneStatusCompleted  MilestoneStatus = "completed"   // Done
	MilestoneStatusDelayed    MilestoneStatus = "delayed"     // Behind schedule
	MilestoneStatusBlocked    MilestoneStatus = "blocked"     // Cannot proceed (waiting for parts, etc.)
	MilestoneStatusSkipped    MilestoneStatus = "skipped"     // Not needed for this ticket
)

// TicketTimeline represents the complete journey of a ticket
type TicketTimeline struct {
	TicketID            string             `json:"ticket_id"`
	OverallStatus       string             `json:"overall_status"` // on_track, at_risk, delayed, blocked
	CurrentStage        MilestoneStage     `json:"current_stage"`
	Milestones          []TicketMilestone  `json:"milestones"`
	EstimatedResolution *time.Time         `json:"estimated_resolution"`
	RequiresParts       bool               `json:"requires_parts"`
	PartsETA            *time.Time         `json:"parts_eta,omitempty"`
	LastUpdated         time.Time          `json:"last_updated"`
}

// PublicTimeline is the customer-facing view of the timeline
type PublicTimeline struct {
	OverallStatus       string                `json:"overall_status"`
	StatusMessage       string                `json:"status_message"`       // Human-friendly message
	CurrentStage        string                `json:"current_stage"`
	CurrentStageDesc    string                `json:"current_stage_desc"`   // What we're doing now
	NextStage           string                `json:"next_stage,omitempty"`
	NextStageDesc       string                `json:"next_stage_desc,omitempty"`
	EstimatedResolution *time.Time            `json:"estimated_resolution"`
	TimeRemaining       string                `json:"time_remaining"`        // Human readable: "2 days, 3 hours"
	RequiresParts       bool                  `json:"requires_parts"`
	PartsStatus         string                `json:"parts_status,omitempty"` // "ordering", "in_transit", "received"
	PartsETA            *time.Time            `json:"parts_eta,omitempty"`
	AssignedEngineer    string                `json:"assigned_engineer,omitempty"`
	Priority            string                `json:"priority"`
	IsUrgent            bool                  `json:"is_urgent"`
	Milestones          []PublicMilestone     `json:"milestones"`
	ProgressPercentage  int                   `json:"progress_percentage"` // 0-100
}

// PublicMilestone is what customers see for each stage
type PublicMilestone struct {
	Stage       string     `json:"stage"`
	Title       string     `json:"title"`        // Customer-friendly title
	Description string     `json:"description"`  // What happens in this stage
	Status      string     `json:"status"`       // pending, in_progress, completed, delayed
	ETA         *time.Time `json:"eta,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	IsActive    bool       `json:"is_active"`    // Currently on this stage
}

// SLAConfig holds SLA configuration based on priority and complexity
type SLAConfig struct {
	Priority                string
	ResponseHours           int // Time to acknowledge
	DiagnosisHours          int // Time to diagnose issue
	SimpleRepairHours       int // Time to fix if no parts needed
	PartsOrderHours         int // Time to order parts
	StandardPartsDelivery   int // Days for standard parts delivery
	UrgentPartsDelivery     int // Days for urgent/express parts
	RepairAfterPartsHours   int // Time to complete repair after parts arrive
	VerificationHours       int // Time for testing/verification
}

// GetSLAConfig returns SLA configuration based on ticket priority
func GetSLAConfig(priority TicketPriority) SLAConfig {
	switch priority {
	case PriorityCritical:
		return SLAConfig{
			Priority:              "critical",
			ResponseHours:         1,  // 1 hour to acknowledge
			DiagnosisHours:        2,  // 2 hours to diagnose
			SimpleRepairHours:     4,  // 4 hours total for simple fix
			PartsOrderHours:       2,  // 2 hours to order parts
			StandardPartsDelivery: 1,  // Next-day for critical
			UrgentPartsDelivery:   0,  // Same-day if available
			RepairAfterPartsHours: 4,  // 4 hours after parts arrive
			VerificationHours:     1,  // 1 hour verification
		}
	case PriorityHigh:
		return SLAConfig{
			Priority:              "high",
			ResponseHours:         2,
			DiagnosisHours:        4,
			SimpleRepairHours:     8,
			PartsOrderHours:       4,
			StandardPartsDelivery: 2,  // 2 business days
			UrgentPartsDelivery:   1,  // 1 day express
			RepairAfterPartsHours: 8,
			VerificationHours:     2,
		}
	case PriorityMedium:
		return SLAConfig{
			Priority:              "medium",
			ResponseHours:         4,
			DiagnosisHours:        8,
			SimpleRepairHours:     24,
			PartsOrderHours:       8,
			StandardPartsDelivery: 5,  // 5 business days
			UrgentPartsDelivery:   2,  // 2 days express
			RepairAfterPartsHours: 16,
			VerificationHours:     4,
		}
	case PriorityLow:
		return SLAConfig{
			Priority:              "low",
			ResponseHours:         8,
			DiagnosisHours:        16,
			SimpleRepairHours:     48,
			PartsOrderHours:       16,
			StandardPartsDelivery: 7,  // 1 week
			UrgentPartsDelivery:   3,  // 3 days
			RepairAfterPartsHours: 24,
			VerificationHours:     8,
		}
	default:
		return GetSLAConfig(PriorityMedium)
	}
}
