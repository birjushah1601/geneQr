# Multi-Model Engineer Assignment System

**Status:** ✅ Complete and Production Ready  
**Date:** December 13, 2025  
**Version:** 1.0.0

---

## Overview

The Multi-Model Engineer Assignment System is an intelligent ticket assignment feature that provides **5 different assignment algorithms** to help administrators select the most suitable engineer for a service ticket. Each model uses different criteria and scoring mechanisms to rank engineers based on various factors like certifications, workload, experience, and skills.

---

## Features

### 5 Assignment Models

1. **Best Overall Match** 🎯
   - **Algorithm:** Combines multiple weighted factors
   - **Factors:**
     - Manufacturer certifications (30%)
     - Experience level (25%)
     - Current workload (25%)
     - Organization tier (20%)
   - **Scoring:** 100-point scale
   - **Use Case:** General-purpose assignment, balanced approach

2. **Manufacturer Certified** 🏭
   - **Algorithm:** Filters engineers certified for specific manufacturers
   - **Factors:**
     - Active manufacturer certifications
     - Certification validity (not expired)
     - Equipment manufacturer matching
   - **Scoring:** Based on certification count and specificity
   - **Use Case:** Equipment requiring manufacturer-certified service

3. **Skills Matched** 🔧
   - **Algorithm:** Matches engineer skills with ticket requirements
   - **Factors:**
     - Skill set alignment
     - Equipment category expertise
     - Relevant experience
   - **Scoring:** Skill overlap percentage
   - **Use Case:** Specialized equipment or specific issue categories

4. **Low Workload** ⚡
   - **Algorithm:** Prioritizes engineers with fewer active tickets
   - **Factors:**
     - Active ticket count
     - In-progress ticket count
     - Average resolution time
   - **Scoring:** Inverse workload ranking
   - **Use Case:** Quick assignment, load balancing

5. **High Seniority** 👔
   - **Algorithm:** Prioritizes most experienced engineers
   - **Factors:**
     - Engineer level (1-5)
     - Years of experience
     - Resolution success rate
   - **Scoring:** Experience-weighted ranking
   - **Use Case:** Critical issues, complex equipment

---

## Architecture

### Backend Components

#### 1. Multi-Model Assignment Service
**Location:** `internal/service-domain/service-ticket/app/multi_model_assignment.go`

```go
type MultiModelAssignmentService struct {
    ticketRepo      domain.TicketRepository
    engineerRepo    domain.EngineerRepository
    equipmentRepo   equipmentDomain.Repository
    assignmentRepo  domain.AssignmentRepository
    logger          *slog.Logger
}
```

**Key Methods:**
- `GetAssignmentSuggestions(ctx, ticketID)` - Returns all 5 models with ranked engineers
- `extractEquipmentContext(ctx, equipmentID)` - Gets manufacturer, category from equipment
- `calculateWorkload(ctx, engineerID)` - Calculates active/in-progress ticket counts
- `scoreEngineer(engineer, equipment, ticket, workload)` - 100-point scoring algorithm

#### 2. Assignment Helpers
**Location:** `internal/service-domain/service-ticket/app/multi_model_assignment_helpers.go`

**Functions:**
- `buildBestOverallMatch()` - Weighted multi-factor scoring
- `buildManufacturerCertified()` - Certification filtering and ranking
- `buildSkillsMatch()` - Skill set alignment scoring
- `buildLowWorkload()` - Workload-based ranking
- `buildHighSeniority()` - Experience-based ranking

#### 3. API Handler
**Location:** `internal/service-domain/service-ticket/api/multi_model_handler.go`

**Endpoint:**
```
GET /api/v1/tickets/{id}/assignment-suggestions
```

**Response Structure:**
```json
{
  "ticket_id": "string",
  "equipment": {
    "id": "string",
    "name": "string",
    "manufacturer": "string",
    "category": "string",
    "model_number": "string"
  },
  "ticket": {
    "priority": "string",
    "min_level_required": 3,
    "requires_certification": true
  },
  "suggestions_by_model": {
    "best_match": {
      "model_name": "Best Overall Match",
      "description": "Combines certification, experience, workload, and organization tier",
      "count": 10,
      "engineers": [
        {
          "id": "string",
          "name": "string",
          "email": "string",
          "engineer_level": 3,
          "match_score": 85,
          "match_reasons": [
            "Level 3 engineer (meets requirement)",
            "2 certifications",
            "0 active tickets"
          ],
          "workload": {
            "active_tickets": 0,
            "in_progress_tickets": 0
          },
          "certifications": [
            {
              "manufacturer": "Wipro GE Healthcare",
              "category": "CT Scanner",
              "is_certified": true
            }
          ]
        }
      ]
    }
  }
}
```

---

### Frontend Components

#### 1. EngineerCard Component
**Location:** `admin-ui/src/components/EngineerCard.tsx`

**Features:**
- Clean, simplified UI
- Prominent match score display
- Level badge with color coding
- Workload status indicator
- Certification count badge
- Active ticket count
- Single "Assign Engineer" action button

**Design:**
```
┌────────────────────────────────┐
│  RK  Rajesh Kumar Singh   60% ││
│      Level 3 - Senior          │
│                                │
│  ✓ 2 certifications            │
│  📊 0 active tickets           │
│                                │
│  [Assign Engineer]             │
└────────────────────────────────┘
```

#### 2. MultiModelAssignment Component
**Location:** `admin-ui/src/components/MultiModelAssignment.tsx`

**Layout:** Side-by-side master-detail interface

**Structure:**
```
┌─────────────────────────────────────────────┐
│ Equipment Context                            │
│ CT Scanner Nova | Wipro GE | Priority: High │
├───────────┬─────────────────────────────────┤
│ FILTERS   │ ENGINEER CARDS                  │
│ (25%)     │ (75%)                           │
│           │                                 │
│ ┌───────┐ │ ┌──────┐ ┌──────┐              │
│ │ Best  │ │ │ Card │ │ Card │              │
│ │ Match │ │ └──────┘ └──────┘              │
│ └───────┘ │                                 │
│           │ ┌──────┐ ┌──────┐              │
│ ┌───────┐ │ │ Card │ │ Card │              │
│ │ Cert. │ │ └──────┘ └──────┘              │
│ └───────┘ │                                 │
│           │                                 │
│ ┌───────┐ │                                 │
│ │ Skills│ │                                 │
│ └───────┘ │                                 │
└───────────┴─────────────────────────────────┘
```

**Key Features:**
- Vertical filter tabs on left (full-width buttons)
- 2-column engineer card grid on right
- Instant updates when switching models
- Horizontal layout option for narrow spaces
- Loading states and error handling
- Confirmation modal before assignment

#### 3. Ticket Detail Page Integration
**Location:** `admin-ui/src/app/tickets/[id]/page.tsx`

**Placement:**
- Shows below Comments section
- Only visible when ticket is unassigned
- Hides after engineer is assigned
- Replaces simple dropdown with intelligent interface

---

## Scoring Algorithm Details

### Best Overall Match Scoring

```typescript
Total Score (100 points) = 
  Certification Score (30 points) +
  Level Score (25 points) +
  Workload Score (25 points) +
  Organization Tier Score (20 points)
```

**Certification Score (30 points):**
- Has manufacturer certification: +15 points
- Has category certification: +10 points
- Additional certifications: +5 points

**Level Score (25 points):**
- Meets min level requirement: +15 points
- Exceeds by 1 level: +20 points
- Exceeds by 2+ levels: +25 points

**Workload Score (25 points):**
- 0 active tickets: +25 points
- 1-2 active tickets: +20 points
- 3-5 active tickets: +15 points
- 6-10 active tickets: +10 points
- 11+ active tickets: +5 points

**Organization Tier Score (20 points):**
- Tier 1 org: +20 points
- Tier 2 org: +15 points
- Tier 3 org: +10 points
- Direct employee: +20 points

---

## Database Schema Changes

### Engineers Table
```sql
-- No changes needed, existing structure supports assignment
CREATE TABLE engineers (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    engineer_level INT DEFAULT 1,
    skills TEXT[],
    home_region VARCHAR(100),
    organization_id VARCHAR(32),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Service Tickets Table
```sql
-- Extended engineer ID field to support longer IDs
ALTER TABLE service_tickets 
ALTER COLUMN assigned_engineer_id TYPE VARCHAR(255);

-- Already has all necessary assignment fields:
-- assigned_engineer_id VARCHAR(255)
-- assigned_engineer_name VARCHAR(255)
-- assigned_at TIMESTAMP
```

### Engineer Assignment Tracking Table
```sql
CREATE TABLE ticket_engineer_assignments (
    id VARCHAR(32) PRIMARY KEY,
    ticket_id VARCHAR(32) NOT NULL,
    engineer_id VARCHAR(32) NOT NULL,
    assignment_tier VARCHAR(10),
    assignment_tier_name VARCHAR(100),
    assigned_by VARCHAR(255),
    assigned_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (ticket_id) REFERENCES service_tickets(id),
    FOREIGN KEY (engineer_id) REFERENCES engineers(id)
);
```

---

## API Integration

### Get Assignment Suggestions

**Request:**
```http
GET /api/v1/tickets/36jsCfunFXbVrGdzZDiOLo0VvQI/assignment-suggestions
```

**Response:** (See Response Structure section above)

### Assign Engineer to Ticket

**Request:**
```http
POST /api/v1/tickets/36jsCfunFXbVrGdzZDiOLo0VvQI/assign-engineer
Content-Type: application/json

{
  "ticket_id": "36jsCfunFXbVrGdzZDiOLo0VvQI",
  "engineer_id": "347S6Gzmxeb1lQzGKeOi7nMYu0N",
  "assignment_tier": "1",
  "assignment_tier_name": "Direct Assignment",
  "assigned_by": "admin"
}
```

**Response:**
```json
{
  "message": "Engineer assigned successfully"
}
```

---

## Usage Guide

### For Administrators

1. **Open Ticket Detail Page**
   - Navigate to unassigned ticket
   - Scroll to "Assign Engineer" section (below comments)

2. **Select Assignment Model**
   - Click vertical tabs on left to switch between models
   - Each model shows different engineer rankings

3. **Review Engineer Details**
   - Check match score (higher is better)
   - Review match reasons
   - See workload and certifications
   - Compare multiple engineers

4. **Assign Engineer**
   - Click "Assign Engineer" button on chosen card
   - Confirm assignment in modal
   - System updates ticket and notifies engineer

### For System Integration

```typescript
// Fetch assignment suggestions
const suggestions = await ticketsApi.getAssignmentSuggestions(ticketId);

// Access specific model
const bestMatch = suggestions.suggestions_by_model.best_match;
const topEngineer = bestMatch.engineers[0];

// Assign engineer
await ticketsApi.assignEngineerToTicket(ticketId, {
  ticket_id: ticketId,
  engineer_id: topEngineer.id,
  assignment_tier: "1",
  assignment_tier_name: "Direct Assignment",
  assigned_by: "admin"
});
```

---

## Performance Considerations

### Backend Optimization
- Equipment context cached per request
- Workload calculations batched for all engineers
- Certification lookups use indexed queries
- Scoring done in-memory after data fetch

### Frontend Optimization
- Assignment data cached for 30 seconds
- Engineer cards rendered in batches
- Horizontal scrolling with virtual rendering
- Lazy loading for large engineer lists

---

## Testing

### Test Coverage

**Backend Tests:**
- ✅ Equipment context extraction
- ✅ Workload calculations
- ✅ Certification matching logic
- ✅ Scoring algorithms for each model
- ✅ API endpoint response structure

**Frontend Tests:**
- ✅ Component rendering
- ✅ Model switching
- ✅ Engineer selection
- ✅ Assignment submission
- ✅ Error handling

### Manual Testing Checklist

- [x] Load assignment suggestions for ticket
- [x] Switch between all 5 models
- [x] Verify engineer rankings differ per model
- [x] Check match scores are calculated correctly
- [x] Test assignment flow end-to-end
- [x] Verify ticket updates after assignment
- [x] Test with no available engineers
- [x] Test with engineers at different levels
- [x] Test with certified vs non-certified engineers
- [x] Test workload-based ranking

---

## Troubleshooting

### Common Issues

**Issue: No engineers returned**
- **Cause:** No engineers in system or all filtered out
- **Solution:** Check engineer database, ensure at least basic level engineers exist

**Issue: Match scores all 0**
- **Cause:** Missing equipment context or engineer data
- **Solution:** Verify equipment has manufacturer/category, engineers have levels

**Issue: Assignment fails with 500 error**
- **Cause:** Database constraint violation (assigned_engineer_id too short)
- **Solution:** Run migration to extend varchar(32) to varchar(255)

**Issue: UI shows "Loading..." indefinitely**
- **Cause:** API endpoint not accessible or returning error
- **Solution:** Check backend logs, verify `/api/v1/tickets/{id}/assignment-suggestions` endpoint

---

## Future Enhancements

### Planned Features
- [ ] Machine learning-based scoring refinement
- [ ] Historical assignment success rate tracking
- [ ] Geographic proximity scoring (distance-based)
- [ ] Availability calendar integration
- [ ] Customer preference tracking
- [ ] Real-time workload monitoring
- [ ] Automated assignment based on rules
- [ ] Multi-engineer assignment for complex tickets

### Potential Improvements
- Add engineer shift schedule consideration
- Include customer satisfaction ratings
- Factor in parts availability at engineer location
- Consider traffic and travel time estimates
- Add peak hour load balancing
- Include language/communication preferences

---

## Related Documentation

- [Engineer Management](./ENGINEER-ASSIGNMENT-COMPLETE-WITH-POSTMAN.md)
- [Service Ticket System](./field-service-management-implementation.md)
- [API Reference](../api/README.md)
- [Database Schema](../database/schema.md)

---

## Changelog

### Version 1.0.0 (December 13, 2025)
- ✅ Initial implementation with 5 assignment models
- ✅ Side-by-side UI with vertical filters
- ✅ Simplified engineer cards
- ✅100-point scoring system
- ✅ Equipment context extraction
- ✅ Workload calculations
- ✅ Certification matching
- ✅ API endpoint and frontend integration
- ✅ Production-ready and tested

---

**Status:** Feature is complete, tested, and deployed to production.  
**Next Steps:** Monitor usage patterns and gather feedback for v2.0 enhancements.
