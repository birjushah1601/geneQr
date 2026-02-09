# SLA Business Hours Implementation Plan

## Requirements

### Current State
- SLA calculated from ticket creation time
- No business hours consideration
- Different hours for different priorities

### Target State  
- **Overall SLA:** 2 days (48 business hours)
- **Business Hours:** 9 AM - 6 PM IST (Monday-Friday)
- **Start Time Logic:**
  - If ticket created during business hours: SLA starts immediately
  - If created after hours/weekend: SLA starts next business day 9 AM
  - Acknowledgment can happen anytime (24/7)
  - SLA timer starts after acknowledgment, respecting business hours

## Implementation Steps

### 1. Create Business Hours Helper (30 min)

**File:** `internal/service-domain/service-ticket/domain/business_hours.go`

```go
package domain

import (
	"time"
)

const (
	BusinessDayStartHour = 9  // 9 AM IST
	BusinessDayEndHour   = 18 // 6 PM IST
	BusinessHoursPerDay  = 9  // 9 hours per day
)

// IsBusinessHours checks if given time is within business hours
func IsBusinessHours(t time.Time) bool {
	// Convert to IST
	ist, _ := time.LoadLocation("Asia/Kolkata")
	istTime := t.In(ist)
	
	// Check if weekend
	if istTime.Weekday() == time.Saturday || istTime.Weekday() == time.Sunday {
		return false
	}
	
	// Check hour
	hour := istTime.Hour()
	return hour >= BusinessDayStartHour && hour < BusinessDayEndHour
}

// NextBusinessHourStart returns next business hour start time
func NextBusinessHourStart(t time.Time) time.Time {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	istTime := t.In(ist)
	
	// If already in business hours, return as is
	if IsBusinessHours(istTime) {
		return istTime
	}
	
	// If same day but after hours, move to next business day
	hour := istTime.Hour()
	if hour >= BusinessDayEndHour {
		istTime = istTime.AddDate(0, 0, 1)
	}
	
	// Skip weekends
	for istTime.Weekday() == time.Saturday || istTime.Weekday() == time.Sunday {
		istTime = istTime.AddDate(0, 0, 1)
	}
	
	// Set to 9 AM
	istTime = time.Date(
		istTime.Year(), istTime.Month(), istTime.Day(),
		BusinessDayStartHour, 0, 0, 0, ist,
	)
	
	return istTime
}

// AddBusinessHours adds business hours to a start time
func AddBusinessHours(start time.Time, hours int) time.Time {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	current := start.In(ist)
	
	// Ensure we start from business hours
	current = NextBusinessHourStart(current)
	
	remainingHours := hours
	
	for remainingHours > 0 {
		// How many hours left in current business day?
		currentHour := current.Hour()
		hoursLeftToday := BusinessDayEndHour - currentHour
		
		if hoursLeftToday >= remainingHours {
			// Can finish today
			current = current.Add(time.Duration(remainingHours) * time.Hour)
			remainingHours = 0
		} else {
			// Use rest of today, move to next business day
			remainingHours -= hoursLeftToday
			current = current.AddDate(0, 0, 1)
			current = NextBusinessHourStart(current)
		}
	}
	
	return current
}

// CalculateBusinessHoursBetween calculates business hours between two times
func CalculateBusinessHoursBetween(start, end time.Time) int {
	ist, _ := time.LoadLocation("Asia/Kolkata")
	current := start.In(ist)
	endTime := end.In(ist)
	
	totalHours := 0
	
	for current.Before(endTime) {
		if IsBusinessHours(current) {
			totalHours++
		}
		current = current.Add(time.Hour)
	}
	
	return totalHours
}
```

### 2. Update SLA Config (15 min)

**File:** `internal/service-domain/service-ticket/domain/sla_timeline.go`

Add new field to SLAConfig:
```go
type SLAConfig struct {
	// ... existing fields ...
	TotalBusinessHours    int  // Total business hours for completion (48 for 2 days)
	UseBusinessHoursOnly  bool // Whether to use business hours calculation
}
```

Update GetSLAConfig for all priorities:
```go
func GetSLAConfig(priority TicketPriority) SLAConfig {
	// All priorities now have same overall SLA: 2 business days
	return SLAConfig{
		Priority:              string(priority),
		TotalBusinessHours:    48, // 2 days * 9 hours * 2 (approx, adjust based on actual calc)
		UseBusinessHoursOnly:  true,
		// Individual stage hours can stay for granularity
		ResponseHours:         varying by priority,
		DiagnosisHours:        varying by priority,
		// ... rest
	}
}
```

### 3. Update Timeline Calculation (1 hour)

**File:** `internal/service-domain/service-ticket/app/timeline_service.go`

Update `buildMilestones` function:

```go
func (s *TimelineService) buildMilestones(ticket *ticketDomain.ServiceTicket, config ticketDomain.SLAConfig, needsParts bool) []ticketDomain.TicketMilestone {
	milestones := []ticketDomain.TicketMilestone{}
	
	// Determine SLA start time
	slaStartTime := ticket.CreatedAt
	
	// If acknowledgment exists, SLA starts from there
	if ticket.AcknowledgedAt != nil {
		slaStartTime = *ticket.AcknowledgedAt
	}
	
	// Adjust to next business hour if needed
	if config.UseBusinessHoursOnly {
		slaStartTime = ticketDomain.NextBusinessHourStart(slaStartTime)
	}
	
	// Build milestones using business hours
	currentTime := slaStartTime
	
	// 1. Acknowledgment (can happen 24/7)
	ackMilestone := ticketDomain.TicketMilestone{
		Stage:       ticketDomain.MilestoneAcknowledgment,
		StartTime:   ticket.CreatedAt,
		ETA:         ticket.CreatedAt.Add(time.Duration(config.ResponseHours) * time.Hour),
		Status:      s.getMilestoneStatus(ticket, ticketDomain.MilestoneAcknowledgment),
		CompletedAt: ticket.AcknowledgedAt,
	}
	milestones = append(milestones, ackMilestone)
	
	// 2. Diagnosis (business hours)
	diagnosisETA := ticketDomain.AddBusinessHours(currentTime, config.DiagnosisHours)
	diagMilestone := ticketDomain.TicketMilestone{
		Stage:     ticketDomain.MilestoneDiagnosis,
		StartTime: currentTime,
		ETA:       &diagnosisETA,
		Status:    s.getMilestoneStatus(ticket, ticketDomain.MilestoneDiagnosis),
	}
	milestones = append(milestones, diagMilestone)
	currentTime = diagnosisETA
	
	// Continue for other milestones...
	// Ensure final ETA is within 2 business days from SLA start
	
	return milestones
}
```

### 4. Update Manual Timeline Edit (30 min)

**File:** `admin-ui/src/components/TimelineEditModal.tsx`

When admin manually changes milestone dates, need to validate and adjust:
- If moving intermediate milestone later, shift all subsequent milestones
- Maintain 3-5 hour gaps between milestones
- Keep everything within business hours

Add helper function:
```typescript
const adjustSubsequentMilestones = (
  milestones: Milestone[],
  changedIndex: number,
  newTime: Date
): Milestone[] => {
  const adjusted = [...milestones];
  let previousTime = newTime;
  
  // Adjust all milestones after the changed one
  for (let i = changedIndex + 1; i < adjusted.length; i++) {
    const gapHours = i === adjusted.length - 1 ? 5 : 4; // Last milestone gets 5 hours
    const newETA = addBusinessHours(previousTime, gapHours);
    adjusted[i] = {
      ...adjusted[i],
      eta: newETA.toISOString()
    };
    previousTime = newETA;
  }
  
  return adjusted;
};
```

### 5. Testing Checklist

**Business Hours Tests:**
- [ ] Ticket created at 10 AM IST → SLA starts at 10 AM
- [ ] Ticket created at 8 PM IST → SLA starts next day 9 AM
- [ ] Ticket created Saturday → SLA starts Monday 9 AM
- [ ] Ticket created Friday 7 PM → SLA starts Monday 9 AM

**SLA Calculation Tests:**
- [ ] 2 day SLA = 48 business hours (approx)
- [ ] Acknowledgment at 11 PM → SLA starts next day 9 AM
- [ ] Milestone ETA calculation respects business hours
- [ ] Weekend hours not counted

**Manual Adjustment Tests:**
- [ ] Move diagnosis 5 hours later → repair starts 4 hours after diagnosis
- [ ] Final completion time adjusts automatically
- [ ] Can't set milestone outside business hours

## Migration Path

1. **Phase 1 (Immediate):** Deploy business hours helper functions
2. **Phase 2 (Next):** Update SLA config, keep existing logic as fallback
3. **Phase 3 (Testing):** Enable business hours for new tickets only
4. **Phase 4 (Full):** Apply to all tickets, remove old logic

## Configuration

Add to `.env`:
```
BUSINESS_HOURS_START=9
BUSINESS_HOURS_END=18
BUSINESS_HOURS_TIMEZONE=Asia/Kolkata
USE_BUSINESS_HOURS_SLA=true
```

## Estimated Time

- Business hours helper: 30 min
- SLA config update: 15 min  
- Timeline calculation: 1 hour
- Manual edit adjustment: 30 min
- Testing: 1 hour
- **Total: 3-4 hours**

## Notes

- Keep existing SLA config as fallback
- Make business hours configurable via environment
- Log SLA calculations for debugging
- Consider public holidays in future iteration
