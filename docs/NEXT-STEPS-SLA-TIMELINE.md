# SLA/ETA Timeline Implementation - Next Steps

## ✅ What's Complete

### Phase 1: Backend Integration ✅
- ✅ Timeline service created and initialized
- ✅ GET `/tickets/{id}/timeline` endpoint working
- ✅ Backend compiles and runs successfully
- ✅ Logs show: "Timeline service initialized" and "Timeline service wired to ticket handler"

### Phase 2: Frontend Component ✅  
- ✅ `TicketTimeline.tsx` React component created
- ✅ TypeScript types (`PublicTimeline`, `PublicMilestone`) added
- ✅ Beautiful UI with animations, progress bars, milestones

### Phase 3: Integration (Partial) ✅
- ✅ Timeline query added to ticket details page
- ✅ Component imported
- ⏸️ **Need to add component to page layout**

---

## 🎯 TO-DO: Complete Integration

### 1. Ticket Details Page (`admin-ui/src/app/tickets/[id]/page.tsx`)

**Find the main content area** (around line 580-900) and add timeline after Parts section:

```tsx
{/* SLA/ETA Timeline - Add this section */}
{timeline && !timelineLoading && (
  <div className="mt-6">
    <h2 className="text-xl font-semibold mb-4">Service Timeline & ETA</h2>
    <TicketTimeline timeline={timeline} />
  </div>
)}

{timelineLoading && (
  <div className="mt-6 flex items-center gap-2 text-gray-500">
    <Loader2 className="h-5 w-5 animate-spin" />
    Loading timeline...
  </div>
)}
```

---

### 2. Public Tracking Page (`admin-ui/src/app/track/[token]/page.tsx`)

This page needs similar integration but for public (non-authenticated) access.

**Current file structure:**
- Fetches public ticket data
- Shows status, comments, history
- Located at `/track/[token]`

**Add timeline:**

1. **Update the API endpoint in handler.go:**
```go
// In GetPublicTicket handler, add timeline
if h.timelineService != nil {
    timeline, err := h.timelineService.GenerateTimeline(ctx, ticket)
    if err == nil {
        publicView.Timeline = h.timelineService.ConvertToPublicTimeline(timeline, ticket)
    }
}
```

2. **Or create new public endpoint:**
```go
// GET /track/{token}/timeline
func (h *TicketHandler) GetPublicTimeline(w http.ResponseWriter, r *http.Request) {
    token := chi.URLParam(r, "token")
    
    // Validate token and get ticket
    trackingToken, err := h.notificationService.GetTrackingToken(token)
    if err != nil {
        h.respondError(w, http.StatusNotFound, "Invalid tracking token")
        return
    }
    
    ticket, err := h.service.GetTicket(ctx, trackingToken.TicketID)
    if err != nil {
        h.respondError(w, http.StatusNotFound, "Ticket not found")
        return
    }
    
    // Generate timeline
    timeline, err := h.timelineService.GenerateTimeline(ctx, ticket)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "Failed to generate timeline")
        return
    }
    
    publicTimeline := h.timelineService.ConvertToPublicTimeline(timeline, ticket)
    h.respondJSON(w, http.StatusOK, publicTimeline)
}
```

3. **Register route:**
```go
// In module.go
r.Get("/track/{token}/timeline", m.ticketHandler.GetPublicTimeline)
```

4. **Frontend:**
```tsx
// In app/track/[token]/page.tsx
const { data: timeline } = useQuery({
  queryKey: ['public-timeline', token],
  queryFn: async () => {
    const response = await fetch(`/api/v1/track/${token}/timeline`);
    return response.json();
  },
});

// In render:
{timeline && <TicketTimeline timeline={timeline} />}
```

---

### 3. Email Notifications (`SendNotificationModal.tsx`)

Add timeline breakdown to email body:

```tsx
// In generateSummary function
if (timeline) {
  summary += `\n📅 Service Timeline:\n`;
  summary += `${'='.repeat(50)}\n\n`;
  
  summary += `Expected Resolution: ${formatDate(timeline.estimated_resolution)}\n`;
  summary += `Time Remaining: ${timeline.time_remaining}\n`;
  summary += `Progress: ${timeline.progress_percentage}%\n\n`;
  
  if (timeline.requires_parts) {
    summary += `📦 Parts Required: ${timeline.parts_status}\n`;
    if (timeline.parts_eta) {
      summary += `Parts Expected: ${formatDate(timeline.parts_eta)}\n`;
    }
    summary += `\n`;
  }
  
  summary += `\nService Journey:\n`;
  timeline.milestones.forEach((m) => {
    const icon = m.status === 'completed' ? '✓' : 
                 m.is_active ? '→' : '○';
    summary += `${icon} ${m.title}\n`;
    if (m.completed_at) {
      summary += `   Completed: ${formatDate(m.completed_at)}\n`;
    } else if (m.eta) {
      summary += `   ETA: ${formatDate(m.eta)}\n`;
    }
    summary += `\n`;
  });
}
```

---

## 🧪 Testing Checklist

### Backend Testing:
```bash
# Test timeline endpoint
curl http://localhost:8081/api/v1/tickets/{ticket_id}/timeline \
  -H "Authorization: Bearer {token}"

# Should return PublicTimeline JSON
```

### Frontend Testing:
1. Open ticket details page
2. Verify timeline loads
3. Check progress bar animation
4. Verify milestone icons and status
5. Test parts status display (if applicable)
6. Check engineer display

### Public Tracking:
1. Get tracking URL from ticket
2. Open in incognito window (no auth)
3. Verify timeline shows
4. Check all milestones visible

### Email Testing:
1. Send notification
2. Check email content includes timeline
3. Verify milestone breakdown
4. Check formatting

---

## 🎨 Optional Enhancements

### 1. Real-Time Updates
Add WebSocket or polling for live progress updates:
```tsx
useEffect(() => {
  const interval = setInterval(() => {
    refetch(); // Refresh timeline every 5 minutes
  }, 5 * 60 * 1000);
  
  return () => clearInterval(interval);
}, []);
```

### 2. Mobile Optimization
Add responsive breakpoints:
```tsx
<div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
  {/* Timeline cards */}
</div>
```

### 3. Status Notifications
Add toast notifications when milestones complete:
```tsx
useEffect(() => {
  if (timeline?.milestones) {
    const justCompleted = timeline.milestones.find(
      m => m.status === 'completed' && 
      new Date(m.completed_at) > Date.now() - 60000 // Last minute
    );
    
    if (justCompleted) {
      toast.success(`${justCompleted.title} completed!`);
    }
  }
}, [timeline]);
```

### 4. ETA Accuracy Tracking
Log actual vs. estimated times for ML improvements:
```go
// When milestone completes
actualDuration := time.Since(milestone.EstimatedStart)
estimatedDuration := milestone.EstimatedComplete.Sub(milestone.EstimatedStart)
accuracy := float64(actualDuration) / float64(estimatedDuration)

// Log for analysis
logger.Info("Milestone accuracy", 
  "stage", milestone.Stage,
  "accuracy", accuracy)
```

---

## 📊 Success Metrics

Track these to measure impact:

1. **Customer Satisfaction**
   - Fewer "where's my ticket?" calls
   - Higher satisfaction scores

2. **Transparency**
   - % of customers who view tracking page
   - Average time on tracking page

3. **SLA Compliance**
   - % of tickets meeting ETA
   - Average delay time

4. **Communication**
   - Reduction in status inquiry emails
   - Proactive notification effectiveness

---

## 🚀 Quick Start

**To complete now:**

1. Add `<TicketTimeline timeline={timeline} />` to ticket details page
2. Test with existing ticket
3. Verify timeline displays correctly
4. Repeat for public tracking page
5. Update email templates

**Estimated time:** 30-60 minutes

---

## 📝 Notes

- Timeline generates automatically based on ticket state
- Parts workflow detection is currently based on ticket status
- Future: Detect parts from `ticket_parts` table
- SLA times configurable per priority
- Milestones skip if not applicable

---

**Current Status:** Backend ✅ Component ✅ Integration ⏸️  
**Blocking:** Add component to page layout (5 min task)  
**Ready for:** Full testing and deployment
