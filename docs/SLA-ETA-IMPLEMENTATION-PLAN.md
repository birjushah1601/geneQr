# SLA/ETA with Parts Workflow - Implementation Summary

## ✅ What We've Created

### 1. Domain Models (`internal/service-domain/service-ticket/domain/sla_timeline.go`)

**New Types:**
- `TicketMilestone` - Individual stages with ETAs
- `TicketTimeline` - Complete ticket journey
- `PublicTimeline` - Customer-facing view
- `SLAConfig` - Priority-based SLA rules

**Milestone Stages:**
1. Acknowledgment - Engineer accepts ticket
2. Diagnosis - Issue identification  
3. Parts Ordered - If parts needed
4. Parts Delivery - Waiting for parts
5. Parts Received - Parts arrived
6. Repair Start - Fix being applied
7. Verification - Testing
8. Resolution - Completed

**SLA Times by Priority:**
```
Critical: 1h response, 4h simple repair, next-day parts
High:     2h response, 8h simple repair, 2-day parts
Medium:   4h response, 24h simple repair, 5-day parts
Low:      8h response, 48h simple repair, 7-day parts
```

### 2. Timeline Service (`internal/service-domain/service-ticket/app/timeline_service.go`)

**Key Methods:**
- `GenerateTimeline()` - Creates milestone timeline
- `ConvertToPublicTimeline()` - Customer-friendly view
- `determineCurrentStage()` - Figures out where ticket is
- `buildMilestones()` - Creates stage-by-stage breakdown

**Smart Features:**
- Detects if parts are required
- Adjusts timeline for parts delivery
- Calculates progress percentage
- Generates empathetic status messages

## 🎯 What Customers Will See

### Public Tracking Page View

```json
{
  "overall_status": "on_track",
  "status_message": "Your service request is progressing well...",
  "current_stage": "parts_delivery",
  "current_stage_desc": "Waiting for parts delivery",
  "next_stage": "repair_start",
  "next_stage_desc": "Performing repair",
  "estimated_resolution": "2026-02-15T15:00:00Z",
  "time_remaining": "3 days, 4 hours",
  "requires_parts": true,
  "parts_status": "in_transit",
  "parts_eta": "2026-02-12T10:00:00Z",
  "assigned_engineer": "Amit Patel",
  "priority": "high",
  "is_urgent": true,
  "progress_percentage": 45,
  "milestones": [
    {
      "stage": "acknowledgment",
      "title": "Acknowledgment",
      "description": "Engineer acknowledges ticket",
      "status": "completed",
      "completed_at": "2026-02-08T10:30:00Z",
      "is_active": false
    },
    {
      "stage": "diagnosis",
      "title": "Diagnosis",
      "description": "Diagnosing the issue",
      "status": "completed",
      "completed_at": "2026-02-08T14:00:00Z",
      "is_active": false
    },
    {
      "stage": "parts_delivery",
      "title": "Parts Delivery",
      "description": "Waiting for parts (2 business days)",
      "status": "in_progress",
      "eta": "2026-02-12T10:00:00Z",
      "is_active": true
    },
    {
      "stage": "repair_start",
      "title": "Repair",
      "description": "Engineer installs parts",
      "status": "pending",
      "eta": "2026-02-12T18:00:00Z",
      "is_active": false
    }
  ]
}
```

## 📋 Next Steps to Complete

### Backend Integration

1. **Add Timeline Service to Module** (`module.go`)
```go
timelineService := app.NewTimelineService(ticketRepo, logger)
ticketHandler.SetTimelineService(timelineService)
```

2. **Add API Endpoint** (`handler.go`)
```go
// GET /tickets/{id}/timeline
func (h *TicketHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
    ticket := h.service.GetTicket(ctx, id)
    timeline := h.timelineService.GenerateTimeline(ctx, ticket)
    publicTimeline := h.timelineService.ConvertToPublicTimeline(timeline, ticket)
    h.respondJSON(w, http.StatusOK, publicTimeline)
}
```

3. **Register Route** (`module.go`)
```go
r.Get("/{id}/timeline", m.ticketHandler.GetTimeline)
```

4. **Add to Public Tracking** (`notification_service.go`)
```go
// GetPublicTicketView should include timeline
view.Timeline = timelineService.ConvertToPublicTimeline(...)
```

### Frontend Implementation

#### 1. Add Timeline Component (`admin-ui/src/components/TicketTimeline.tsx`)

```tsx
interface TimelineProps {
  timeline: PublicTimeline;
}

export function TicketTimeline({ timeline }: TimelineProps) {
  return (
    <div className="space-y-6">
      {/* Overall Status Card */}
      <Card className="bg-gradient-to-r from-blue-50 to-indigo-50">
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              Expected Timeline
            </CardTitle>
            <StatusBadge status={timeline.overall_status} />
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-gray-700 mb-4">
            {timeline.status_message}
          </p>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <p className="text-sm text-gray-600">Target Resolution</p>
              <p className="text-xl font-semibold">
                {formatDate(timeline.estimated_resolution)}
              </p>
              <p className="text-sm text-gray-500 mt-1">
                {timeline.time_remaining}
              </p>
            </div>
            
            {timeline.is_urgent && (
              <div className="flex items-center gap-2 text-orange-600">
                <AlertCircle className="h-5 w-5" />
                <span className="font-semibold">High Priority</span>
              </div>
            )}
          </div>

          {/* Progress Bar */}
          <div className="mt-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-600">Progress</span>
              <span className="text-sm font-semibold">{timeline.progress_percentage}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className="bg-blue-600 h-2 rounded-full transition-all"
                style={{ width: `${timeline.progress_percentage}%` }}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Parts Status (if applicable) */}
      {timeline.requires_parts && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Parts Status
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <div className="flex-1">
                <p className="text-sm text-gray-600">Status</p>
                <p className="text-lg font-semibold capitalize">
                  {timeline.parts_status.replace('_', ' ')}
                </p>
              </div>
              {timeline.parts_eta && (
                <div className="flex-1">
                  <p className="text-sm text-gray-600">Expected Arrival</p>
                  <p className="text-lg font-semibold">
                    {formatDate(timeline.parts_eta)}
                  </p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Milestone Timeline */}
      <Card>
        <CardHeader>
          <CardTitle>Service Journey</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="relative">
            {timeline.milestones.map((milestone, index) => (
              <div 
                key={milestone.stage}
                className={`flex gap-4 pb-8 ${
                  index === timeline.milestones.length - 1 ? 'pb-0' : ''
                }`}
              >
                {/* Timeline line */}
                {index < timeline.milestones.length - 1 && (
                  <div className="absolute left-4 top-10 w-0.5 h-full bg-gray-200" />
                )}
                
                {/* Milestone icon */}
                <div className={`
                  relative z-10 flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center
                  ${milestone.status === 'completed' ? 'bg-green-500 text-white' : ''}
                  ${milestone.is_active ? 'bg-blue-500 text-white animate-pulse' : ''}
                  ${milestone.status === 'pending' ? 'bg-gray-200 text-gray-500' : ''}
                  ${milestone.status === 'blocked' ? 'bg-yellow-500 text-white' : ''}
                `}>
                  {milestone.status === 'completed' && <CheckCircle className="h-5 w-5" />}
                  {milestone.is_active && <Loader2 className="h-5 w-5 animate-spin" />}
                  {milestone.status === 'pending' && <Clock className="h-4 w-4" />}
                  {milestone.status === 'blocked' && <AlertTriangle className="h-5 w-5" />}
                </div>

                {/* Milestone content */}
                <div className="flex-1">
                  <div className="flex items-center justify-between">
                    <h3 className={`font-semibold ${
                      milestone.is_active ? 'text-blue-600' : 'text-gray-900'
                    }`}>
                      {milestone.title}
                    </h3>
                    {milestone.eta && !milestone.completed_at && (
                      <span className="text-sm text-gray-500">
                        ETA: {formatTime(milestone.eta)}
                      </span>
                    )}
                    {milestone.completed_at && (
                      <span className="text-sm text-green-600">
                        ✓ {formatTime(milestone.completed_at)}
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-gray-600 mt-1">
                    {milestone.description}
                  </p>
                  {milestone.is_active && (
                    <p className="text-sm text-blue-600 mt-2 font-medium">
                      → Currently in progress
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Assigned Engineer */}
      {timeline.assigned_engineer && (
        <Card>
          <CardContent className="py-4">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white font-semibold">
                {timeline.assigned_engineer.split(' ').map(n => n[0]).join('')}
              </div>
              <div>
                <p className="text-sm text-gray-600">Assigned Engineer</p>
                <p className="font-semibold">{timeline.assigned_engineer}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
```

#### 2. Add to Public Tracking Page (`admin-ui/src/app/track/[token]/page.tsx`)

```tsx
// Fetch timeline
const { data: timeline } = useQuery({
  queryKey: ['ticket-timeline', token],
  queryFn: () => apiClient.get(`/v1/track/${token}/timeline`),
});

// In render:
{timeline && <TicketTimeline timeline={timeline} />}
```

#### 3. Add to Email Template

Update `SendNotificationModal.tsx` to include timeline info:

```typescript
if (timeline) {
  summary += `\n📅 Timeline Breakdown:\n`;
  summary += `${'='.repeat(50)}\n\n`;
  
  timeline.milestones.forEach((m) => {
    const status = m.status === 'completed' ? '✓' : 
                   m.is_active ? '→' : '○';
    summary += `${status} ${m.title}\n`;
    summary += `   ${m.description}\n`;
    if (m.eta) {
      summary += `   ETA: ${formatDate(m.eta)}\n`;
    }
    summary += `\n`;
  });
}
```

## 🎨 Visual Design Mockup

```
┌─────────────────────────────────────────────────┐
│ ⏰ Expected Timeline              🟢 On Track   │
├─────────────────────────────────────────────────┤
│ Your service request is progressing well.       │
│ We've identified the parts needed and are       │
│ working to complete the repair by Feb 15.       │
│                                                  │
│ Target Resolution: Feb 15, 2026 at 3:00 PM     │
│ Time Remaining: 3 days, 4 hours                 │
│                                                  │
│ Progress: ████████████░░░░░░░░░ 60%             │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ 📦 Parts Status                                  │
├─────────────────────────────────────────────────┤
│ Status: In Transit                               │
│ Expected Arrival: Feb 12, 2026                   │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ Service Journey                                  │
├─────────────────────────────────────────────────┤
│                                                  │
│ ✓ Acknowledgment          ✓ Feb 8, 10:30 AM    │
│   Engineer reviewed request                      │
│ │                                                │
│ ✓ Diagnosis               ✓ Feb 8, 2:00 PM     │
│   Issue identified                               │
│ │                                                │
│ ✓ Parts Ordered           ✓ Feb 8, 4:00 PM     │
│   Required parts ordered                         │
│ │                                                │
│ ⟳ Parts Delivery          ETA: Feb 12, 10:00 AM │
│   Waiting for parts (2 business days)           │
│   → Currently in progress                        │
│ │                                                │
│ ○ Repair                  ETA: Feb 12, 6:00 PM  │
│   Engineer will install parts                    │
│ │                                                │
│ ○ Verification            ETA: Feb 13, 10:00 AM │
│   Testing and verification                       │
│ │                                                │
│ ○ Completed               ETA: Feb 15, 3:00 PM  │
│   Service resolved                               │
│                                                  │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ 👤 Assigned Engineer: Amit Patel                │
└─────────────────────────────────────────────────┘
```

## 🚀 Benefits

1. **Transparency** - Customers know exactly where their ticket stands
2. **Expectation Management** - Clear ETAs reduce anxiety
3. **Parts Visibility** - Customers understand delays caused by parts
4. **Progress Tracking** - Visual progress bar shows advancement
5. **Engineer Accountability** - Name attached to timeline builds trust
6. **Empathetic Messaging** - Human-friendly communication
7. **Priority Awareness** - Urgent tickets show special handling

## 🔄 Real-World Scenarios

### Scenario 1: Simple Repair (No Parts)
```
Created → Acknowledged (1h) → Diagnosed (2h) → Repaired (4h) → Verified (1h)
Total: 8 hours for high-priority ticket
```

### Scenario 2: Parts Required
```
Created → Acknowledged (1h) → Diagnosed (2h) → Parts Ordered (2h) → 
Parts Delivery (2 days) → Repair (4h) → Verified (1h)
Total: ~2.5 days for high-priority ticket
```

### Scenario 3: Critical with Express Parts
```
Created → Acknowledged (30min) → Diagnosed (1h) → Parts Ordered (1h) → 
Parts Delivery (same-day) → Repair (2h) → Verified (30min)
Total: Same day resolution
```

## 📊 Future Enhancements

1. **SMS Notifications** - Alert on milestone completion
2. **Real-Time Updates** - WebSocket for live progress
3. **Photo Updates** - Engineer uploads progress photos
4. **Customer Feedback** - Rate each milestone
5. **Predictive ETA** - ML-based timeline adjustments
6. **Multi-Engineer** - Show handoffs between engineers
7. **Location Tracking** - Engineer en-route visibility

---

**Status:** Domain models and service created ✅  
**Next:** Wire into handlers and create frontend UI  
**Ready to implement?** Yes! Just need integration work.
