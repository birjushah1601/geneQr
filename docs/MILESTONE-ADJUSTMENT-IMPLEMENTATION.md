# Milestone Adjustment Logic Implementation Plan

## Problem

When admin changes an intermediate milestone to a later time, subsequent milestones don't automatically adjust. This creates illogical situations where:
- Diagnosis ETA: Feb 10, 3:00 PM (moved later)
- Repair ETA: Feb 9, 5:00 PM (earlier than diagnosis!)
- Target completion: Feb 9, 8:00 PM (earlier than diagnosis!)

## Solution

Auto-adjust all subsequent milestones when an intermediate milestone is moved, maintaining 3-5 hour gaps.

## Implementation

### 1. Frontend: Real-time Adjustment (1 hour)

**File:** `admin-ui/src/components/TimelineEditModal.tsx`

Add adjustment logic:

```typescript
// Add after imports
import { addBusinessHours, isBusinessHours } from '@/lib/businessHours';

// Helper: Add hours considering business hours if needed
const addHours = (date: Date, hours: number, useBusinessHours: boolean): Date => {
  if (useBusinessHours) {
    return addBusinessHours(date, hours);
  }
  return new Date(date.getTime() + hours * 60 * 60 * 1000);
};

// Helper: Adjust subsequent milestones
const adjustSubsequentMilestones = (
  milestones: PublicMilestone[],
  changedIndex: number,
  newETA: string,
  useBusinessHours: boolean = true
): PublicMilestone[] => {
  const adjusted = [...milestones];
  let previousDate = new Date(newETA);
  
  console.log(`Adjusting milestones after index ${changedIndex}`, {
    changedMilestone: adjusted[changedIndex].stage,
    newTime: newETA
  });
  
  // Adjust all milestones after the changed one
  for (let i = changedIndex + 1; i < adjusted.length; i++) {
    const milestone = adjusted[i];
    
    // Determine gap hours based on milestone type
    let gapHours = 4; // Default gap
    
    if (milestone.stage === 'verification' || i === adjusted.length - 1) {
      gapHours = 5; // Last milestone or verification gets 5 hours
    } else if (milestone.stage === 'parts_delivery') {
      gapHours = 24; // Parts delivery typically 1 day later
    } else if (milestone.stage === 'repair_complete') {
      gapHours = 8; // Repair takes longer
    }
    
    // Calculate new ETA
    const newMilestoneETA = addHours(previousDate, gapHours, useBusinessHours);
    
    adjusted[i] = {
      ...milestone,
      eta: newMilestoneETA.toISOString()
    };
    
    previousDate = newMilestoneETA;
    
    console.log(`Adjusted ${milestone.stage} to ${newMilestoneETA.toISOString()}`);
  }
  
  // Update target completion (last milestone)
  const finalMilestone = adjusted[adjusted.length - 1];
  if (finalMilestone.eta) {
    const targetDate = new Date(finalMilestone.eta);
    // Add 1 hour buffer after last milestone
    const targetCompletion = addHours(targetDate, 1, useBusinessHours);
    
    console.log('Updated target completion:', targetCompletion.toISOString());
    
    // Store this in editedTimeline state
    // (will need to update parent state)
  }
  
  return adjusted;
};

// Update handleMilestoneEtaChange function
const handleMilestoneEtaChange = (index: number, newDate: string) => {
  const updatedMilestones = [...editedTimeline.milestones];
  updatedMilestones[index] = {
    ...updatedMilestones[index],
    eta: newDate
  };
  
  // Auto-adjust subsequent milestones
  const adjustedMilestones = adjustSubsequentMilestones(
    updatedMilestones,
    index,
    newDate,
    true // use business hours
  );
  
  // Also update target completion time
  const finalMilestone = adjustedMilestones[adjustedMilestones.length - 1];
  if (finalMilestone.eta) {
    const finalDate = new Date(finalMilestone.eta);
    const targetCompletion = addHours(finalDate, 1, true);
    
    setEditedTimeline({
      ...editedTimeline,
      milestones: adjustedMilestones,
      estimated_resolution: targetCompletion.toISOString()
    });
  } else {
    setEditedTimeline({
      ...editedTimeline,
      milestones: adjustedMilestones
    });
  }
};
```

### 2. Frontend: Business Hours Utility (30 min)

**File:** `admin-ui/src/lib/businessHours.ts`

```typescript
export const BUSINESS_START_HOUR = 9;
export const BUSINESS_END_HOUR = 18;
export const BUSINESS_HOURS_PER_DAY = 9;

/**
 * Check if a date/time is within business hours (9 AM - 6 PM IST, Mon-Fri)
 */
export function isBusinessHours(date: Date): boolean {
  // Convert to IST
  const istDate = new Date(date.toLocaleString('en-US', { timeZone: 'Asia/Kolkata' }));
  
  // Check weekend
  const day = istDate.getDay();
  if (day === 0 || day === 6) return false; // Sunday or Saturday
  
  // Check hour
  const hour = istDate.getHours();
  return hour >= BUSINESS_START_HOUR && hour < BUSINESS_END_HOUR;
}

/**
 * Get next business hour start time
 */
export function nextBusinessHourStart(date: Date): Date {
  let current = new Date(date);
  
  // If already in business hours, return as is
  if (isBusinessHours(current)) {
    return current;
  }
  
  // Convert to IST for calculations
  const istDate = new Date(current.toLocaleString('en-US', { timeZone: 'Asia/Kolkata' }));
  
  // If after business hours, move to next day
  if (istDate.getHours() >= BUSINESS_END_HOUR) {
    current.setDate(current.getDate() + 1);
  }
  
  // Skip weekends
  while (current.getDay() === 0 || current.getDay() === 6) {
    current.setDate(current.getDate() + 1);
  }
  
  // Set to 9 AM IST
  current.setHours(BUSINESS_START_HOUR, 0, 0, 0);
  
  return current;
}

/**
 * Add business hours to a date
 */
export function addBusinessHours(start: Date, hours: number): Date {
  let current = nextBusinessHourStart(new Date(start));
  let remainingHours = hours;
  
  while (remainingHours > 0) {
    // How many hours left in current business day?
    const currentHour = current.getHours();
    const hoursLeftToday = BUSINESS_END_HOUR - currentHour;
    
    if (hoursLeftToday >= remainingHours) {
      // Can finish today
      current.setHours(current.getHours() + remainingHours);
      remainingHours = 0;
    } else {
      // Use rest of today, move to next business day
      remainingHours -= hoursLeftToday;
      current.setDate(current.getDate() + 1);
      current = nextBusinessHourStart(current);
    }
  }
  
  return current;
}
```

### 3. Backend Validation (Optional, 30 min)

**File:** `internal/service-domain/service-ticket/api/handler.go`

Add validation in UpdateTimeline:

```go
// Validate milestone sequence
for i := 1; i < len(req.Milestones); i++ {
	prev := req.Milestones[i-1]
	curr := req.Milestones[i]
	
	if prev.ETA != nil && curr.ETA != nil {
		if curr.ETA.Before(*prev.ETA) {
			h.respondError(w, http.StatusBadRequest, 
				fmt.Sprintf("Milestone %s cannot be before %s", 
					curr.Stage, prev.Stage))
			return
		}
	}
}

// Validate target completion is after last milestone
if len(req.Milestones) > 0 && req.EstimatedResolution != nil {
	lastMilestone := req.Milestones[len(req.Milestones)-1]
	if lastMilestone.ETA != nil && req.EstimatedResolution.Before(*lastMilestone.ETA) {
		h.respondError(w, http.StatusBadRequest, 
			"Target completion cannot be before last milestone")
		return
	}
}
```

### 4. UI Feedback (15 min)

Add visual indicator when milestones are auto-adjusted:

```typescript
// In TimelineEditModal.tsx
const [autoAdjusted, setAutoAdjusted] = useState<number[]>([]);

// After adjusting milestones
setAutoAdjusted(Array.from({length: adjustedMilestones.length - changedIndex - 1}, 
  (_, i) => changedIndex + i + 1));

// Clear after 3 seconds
setTimeout(() => setAutoAdjusted([]), 3000);

// In render
{autoAdjusted.includes(index) && (
  <span className="text-xs text-blue-600 ml-2">
    ✓ Auto-adjusted
  </span>
)}
```

## Testing Scenarios

### Test 1: Move Diagnosis Later
1. Create ticket with default timeline
2. Move "Diagnosis" from 2 PM to 5 PM
3. **Expected:**
   - Repair starts at 9 PM (4 hours after diagnosis)
   - Verification at 5 AM next day (8 hours after repair)
   - Target completion at 6 AM (1 hour after verification)

### Test 2: Move with Business Hours
1. Move "Diagnosis" to 5 PM (near end of business day)
2. **Expected:**
   - Repair starts next day 9 AM (business hours logic)
   - All subsequent milestones follow business hours

### Test 3: Move Last Milestone
1. Move "Verification" 2 days later
2. **Expected:**
   - Target completion adjusts to 1 hour after verification
   - Previous milestones unchanged

### Test 4: Validation
1. Try to set milestone before previous one
2. **Expected:**
   - Error message
   - Changes not saved

## Gap Hour Configuration

Different milestone types should have different gaps:

| From Milestone | To Milestone | Gap Hours |
|----------------|--------------|-----------|
| Acknowledgment | Diagnosis | 3-4 hours |
| Diagnosis | Repair Start | 4-5 hours |
| Repair Start | Parts Delivery | 24-48 hours |
| Parts Delivery | Repair Complete | 6-8 hours |
| Repair Complete | Verification | 4-5 hours |
| Verification | Target Completion | 1 hour |

## Configuration

Add to config:
```typescript
export const MILESTONE_GAPS = {
  acknowledgment_to_diagnosis: 4,
  diagnosis_to_repair: 5,
  repair_to_parts: 24,
  parts_to_repair_complete: 8,
  repair_to_verification: 5,
  verification_to_completion: 1
};
```

## Estimated Time

- Frontend adjustment logic: 1 hour
- Business hours utility: 30 min
- Backend validation: 30 min
- UI feedback: 15 min
- Testing: 45 min
- **Total: 3 hours**

## Implementation Order

1. Create business hours utility
2. Add adjustment logic to TimelineEditModal
3. Test with various scenarios
4. Add backend validation
5. Add UI feedback

## Notes

- Keep existing behavior as default
- Make auto-adjustment optional (checkbox?)
- Log all adjustments for debugging
- Consider making gap hours configurable
