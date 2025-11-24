package assignment

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	
)

// Scorer calculates engineer matching scores
type Scorer struct {
	db *sql.DB
}

// NewScorer creates a new scorer
func NewScorer(db *sql.DB) *Scorer {
	return &Scorer{
		db: db,
	}
}

// ScoreEngineer calculates overall score for an engineer
func (s *Scorer) ScoreEngineer(ctx context.Context, engineer *EngineerProfile, req *AssignmentRequest) (*ScoreBreakdown, []MatchReason, []string, error) {
	breakdown := &ScoreBreakdown{
		Weights: req.Options.Weights,
	}

	var reasons []MatchReason
	var warnings []string

	// Calculate individual scores
	if req.Options.ConsiderSpecialization {
		score, matchReasons := s.calculateExpertiseScore(engineer, req)
		breakdown.ExpertiseScore = score
		breakdown.WeightedExpertise = score * req.Options.Weights.Expertise
		reasons = append(reasons, matchReasons...)
	}

	if req.Options.ConsiderLocation {
		score, matchReasons, warns := s.calculateLocationScore(engineer, req)
		breakdown.LocationScore = score
		breakdown.WeightedLocation = score * req.Options.Weights.Location
		reasons = append(reasons, matchReasons...)
		warnings = append(warnings, warns...)
	}

	if req.Options.ConsiderPerformance {
		score, matchReasons := s.calculatePerformanceScore(engineer, req)
		breakdown.PerformanceScore = score
		breakdown.WeightedPerformance = score * req.Options.Weights.Performance
		reasons = append(reasons, matchReasons...)
	}

	if req.Options.ConsiderWorkload {
		score, matchReasons, warns := s.calculateWorkloadScore(engineer, req)
		breakdown.WorkloadScore = score
		breakdown.WeightedWorkload = score * req.Options.Weights.Workload
		reasons = append(reasons, matchReasons...)
		warnings = append(warnings, warns...)
	}

	if req.Options.ConsiderAvailability {
		score, matchReasons, warns := s.calculateAvailabilityScore(engineer, req)
		breakdown.AvailabilityScore = score
		breakdown.WeightedAvailability = score * req.Options.Weights.Availability
		reasons = append(reasons, matchReasons...)
		warnings = append(warnings, warns...)
	}

	return breakdown, reasons, warnings, nil
}

// calculateExpertiseScore scores based on skills and equipment expertise
func (s *Scorer) calculateExpertiseScore(engineer *EngineerProfile, req *AssignmentRequest) (float64, []MatchReason) {
	var reasons []MatchReason
	score := 0.0

	// Check equipment expertise
	_ = false
	var bestMatch *EquipmentExpertise

	for i := range engineer.EquipmentExpertise {
		exp := &engineer.EquipmentExpertise[i]
		if exp.EquipmentType == req.EquipmentType {
			_ = true // hasEquipmentExpertise
			if bestMatch == nil || exp.SuccessRate > bestMatch.SuccessRate {
				bestMatch = exp
			}
		}
	}

	if bestMatch != nil {
		// Score based on tickets handled and success rate
		ticketsScore := math.Min(float64(bestMatch.TicketsHandled)/20.0, 1.0) * 50 // Max 50 points
		successScore := bestMatch.SuccessRate * 50                                   // Max 50 points
		score = ticketsScore + successScore

		reasons = append(reasons, MatchReason{
			Category: CategoryExpertise,
			Reason:   fmt.Sprintf("Handled %d %s tickets with %.0f%% success rate", bestMatch.TicketsHandled, req.EquipmentType, bestMatch.SuccessRate),
			Impact:   ImpactHigh,
			Evidence: fmt.Sprintf("%d tickets, %.1f%% success", bestMatch.TicketsHandled, bestMatch.SuccessRate),
		})
	} else {
		// No direct equipment experience - check for related skills
		score = 30.0 // Base score
		reasons = append(reasons, MatchReason{
			Category: CategoryExpertise,
			Reason:   "No direct experience with this equipment type",
			Impact:   ImpactMedium,
			Evidence: "General technical skills",
		})
	}

	// Check for required skills
	if len(req.RequiredSkills) > 0 {
		matchingSkills := 0
		totalRequired := len(req.RequiredSkills)

		for _, reqSkill := range req.RequiredSkills {
			for _, engSkill := range engineer.Skills {
				if strings.EqualFold(engSkill.SkillName, reqSkill) {
					matchingSkills++
					// Bonus for advanced proficiency
					if engSkill.ProficiencyLevel == ProficiencyExpert {
						score += 5
					} else if engSkill.ProficiencyLevel == ProficiencyAdvanced {
						score += 3
					}
					break
				}
			}
		}

		skillMatchRate := float64(matchingSkills) / float64(totalRequired)
		if skillMatchRate > 0.8 {
			reasons = append(reasons, MatchReason{
				Category: CategoryExpertise,
				Reason:   fmt.Sprintf("Matches %d/%d required skills", matchingSkills, totalRequired),
				Impact:   ImpactHigh,
				Evidence: fmt.Sprintf("%.0f%% skill match", skillMatchRate*100),
			})
		}
	}

	// Check specialization
	if req.RequiresSpecialist && engineer.Specialization != nil {
		if req.SpecialistType != nil && strings.EqualFold(*engineer.Specialization, *req.SpecialistType) {
			score += 15
			reasons = append(reasons, MatchReason{
				Category: CategoryExpertise,
				Reason:   fmt.Sprintf("Specialist: %s", *engineer.Specialization),
				Impact:   ImpactHigh,
				Evidence: "Matching specialization",
			})
		}
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score, reasons
}

// calculateLocationScore scores based on proximity to ticket location
func (s *Scorer) calculateLocationScore(engineer *EngineerProfile, req *AssignmentRequest) (float64, []MatchReason, []string) {
	var reasons []MatchReason
	var warnings []string

	// If no location info, return neutral score
	if engineer.CurrentLocation == nil {
		return 50.0, reasons, []string{"Engineer location unknown"}
	}

	// Simple location matching (in production, use actual distance calculation)
	// For now, check if locations match or are nearby

	if strings.Contains(strings.ToLower(*engineer.CurrentLocation), strings.ToLower(req.LocationName)) {
		reasons = append(reasons, MatchReason{
			Category: CategoryLocation,
			Reason:   "Currently at same location",
			Impact:   ImpactHigh,
			Evidence: fmt.Sprintf("Location: %s", *engineer.CurrentLocation),
		})
		return 100.0, reasons, warnings
	}

	// Check if in same building/area (simplified)
	if strings.Contains(strings.ToLower(req.LocationName), "icu") &&
		strings.Contains(strings.ToLower(*engineer.CurrentLocation), "icu") {
		reasons = append(reasons, MatchReason{
			Category: CategoryLocation,
			Reason:   "In same department area",
			Impact:   ImpactMedium,
			Evidence: "ICU area",
		})
		return 80.0, reasons, warnings
	}

	// Different location - lower score but not critical
	reasons = append(reasons, MatchReason{
		Category: CategoryLocation,
		Reason:   "At different location, travel time required",
		Impact:   ImpactMedium,
		Evidence: fmt.Sprintf("Current: %s, Ticket: %s", *engineer.CurrentLocation, req.LocationName),
	})

	return 40.0, reasons, warnings
}

// calculatePerformanceScore scores based on historical performance
func (s *Scorer) calculatePerformanceScore(engineer *EngineerProfile, req *AssignmentRequest) (float64, []MatchReason) {
	var reasons []MatchReason
	score := 50.0 // Base score

	// Success rate component (50 points max)
	successScore := engineer.SuccessRate * 50
	score = successScore

	if engineer.SuccessRate >= 0.9 {
		reasons = append(reasons, MatchReason{
			Category: CategoryPerformance,
			Reason:   fmt.Sprintf("Excellent success rate: %.0f%%", engineer.SuccessRate*100),
			Impact:   ImpactHigh,
			Evidence: fmt.Sprintf("%d tickets resolved", engineer.TotalTicketsResolved),
		})
	} else if engineer.SuccessRate >= 0.75 {
		reasons = append(reasons, MatchReason{
			Category: CategoryPerformance,
			Reason:   fmt.Sprintf("Good success rate: %.0f%%", engineer.SuccessRate*100),
			Impact:   ImpactMedium,
			Evidence: fmt.Sprintf("%d tickets resolved", engineer.TotalTicketsResolved),
		})
	}

	// Resolution time component (25 points max)
	avgHours := engineer.AverageResolutionTime.Hours()
	if avgHours > 0 {
		// Score inversely proportional to time (faster = better)
		// Assume target is 4 hours, scale accordingly
		timeScore := math.Max(0, 25*(1-avgHours/24.0))
		score += timeScore

		if avgHours <= 2 {
			reasons = append(reasons, MatchReason{
				Category: CategoryPerformance,
				Reason:   fmt.Sprintf("Fast average resolution: %.1f hours", avgHours),
				Impact:   ImpactMedium,
				Evidence: "Quick turnaround",
			})
		}
	}

	// Customer rating component (25 points max)
	if engineer.AverageRating != nil {
		ratingScore := (*engineer.AverageRating / 5.0) * 25
		score += ratingScore

		if *engineer.AverageRating >= 4.5 {
			reasons = append(reasons, MatchReason{
				Category: CategoryPerformance,
				Reason:   fmt.Sprintf("Excellent customer rating: %.1f/5", *engineer.AverageRating),
				Impact:   ImpactMedium,
				Evidence: "High satisfaction",
			})
		}
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score, reasons
}

// calculateWorkloadScore scores based on current workload
func (s *Scorer) calculateWorkloadScore(engineer *EngineerProfile, req *AssignmentRequest) (float64, []MatchReason, []string) {
	var reasons []MatchReason
	var warnings []string

	// Ideal workload is 3-5 open tickets
	// Score decreases with too many or too few tickets

	openTickets := engineer.OpenTicketsCount
	score := 100.0

	if openTickets == 0 {
		// Available but might be rusty or on leave
		score = 90.0
		reasons = append(reasons, MatchReason{
			Category: CategoryWorkload,
			Reason:   "No current assignments - fully available",
			Impact:   ImpactHigh,
			Evidence: "0 open tickets",
		})
	} else if openTickets <= 3 {
		// Optimal workload
		score = 100.0
		reasons = append(reasons, MatchReason{
			Category: CategoryWorkload,
			Reason:   "Light workload - good availability",
			Impact:   ImpactHigh,
			Evidence: fmt.Sprintf("%d open tickets", openTickets),
		})
	} else if openTickets <= 5 {
		// Moderate workload
		score = 80.0
		reasons = append(reasons, MatchReason{
			Category: CategoryWorkload,
			Reason:   "Moderate workload",
			Impact:   ImpactMedium,
			Evidence: fmt.Sprintf("%d open tickets", openTickets),
		})
	} else if openTickets <= 8 {
		// Heavy workload
		score = 50.0
		warnings = append(warnings, fmt.Sprintf("Heavy workload: %d open tickets", openTickets))
		reasons = append(reasons, MatchReason{
			Category: CategoryWorkload,
			Reason:   "Heavy workload - may affect response time",
			Impact:   ImpactMedium,
			Evidence: fmt.Sprintf("%d open tickets", openTickets),
		})
	} else {
		// Overloaded
		score = 20.0
		warnings = append(warnings, fmt.Sprintf("Overloaded: %d open tickets", openTickets))
		reasons = append(reasons, MatchReason{
			Category: CategoryWorkload,
			Reason:   "Overloaded - not recommended for new assignments",
			Impact:   ImpactHigh,
			Evidence: fmt.Sprintf("%d open tickets", openTickets),
		})
	}

	// Adjust for priority
	if req.Priority == "Critical" && openTickets > 5 {
		score *= 0.5 // Penalize heavily for critical tickets
		warnings = append(warnings, "High workload not ideal for critical ticket")
	}

	return score, reasons, warnings
}

// calculateAvailabilityScore scores based on current availability
func (s *Scorer) calculateAvailabilityScore(engineer *EngineerProfile, req *AssignmentRequest) (float64, []MatchReason, []string) {
	var reasons []MatchReason
	var warnings []string

	switch engineer.AvailabilityStatus {
	case AvailabilityAvailable:
		reasons = append(reasons, MatchReason{
			Category: CategoryAvailability,
			Reason:   "Currently available",
			Impact:   ImpactHigh,
			Evidence: "Active status",
		})
		return 100.0, reasons, warnings

	case AvailabilityBusy:
		reasons = append(reasons, MatchReason{
			Category: CategoryAvailability,
			Reason:   "Currently busy with another ticket",
			Impact:   ImpactMedium,
			Evidence: "Busy status",
		})
		return 60.0, reasons, warnings

	case AvailabilityOnLeave:
		warnings = append(warnings, "Engineer is on leave")
		return 0.0, reasons, warnings

	case AvailabilityOffline:
		warnings = append(warnings, "Engineer is offline")
		return 10.0, reasons, warnings

	default:
		return 50.0, reasons, warnings
	}
}

// CalculateOverallScore combines weighted scores
func (s *Scorer) CalculateOverallScore(breakdown *ScoreBreakdown) float64 {
	overall := breakdown.WeightedExpertise +
		breakdown.WeightedLocation +
		breakdown.WeightedPerformance +
		breakdown.WeightedWorkload +
		breakdown.WeightedAvailability

	// Normalize to 0-100 range
	// Since weights should sum to 1.0, overall should be 0-100
	return math.Min(100.0, math.Max(0.0, overall))
}




