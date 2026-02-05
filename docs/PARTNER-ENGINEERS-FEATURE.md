# Partner Engineers Feature

## Overview
The Partner Engineers feature allows manufacturers to view and assign engineers from their partner organizations (Channel Partners and Sub-Dealers) to service tickets.

---

## Problem Statement

### Before
- Engineers API only returned engineers from the manufacturer's own organization
- Query filtered: `WHERE eom.org_id = '<manufacturer_org_id>'`
- Partner engineers were invisible to manufacturers
- Couldn't assign partner engineers to tickets

### After
- Engineers API includes partner engineers when requested
- Manufacturers can see engineers from authorized partner organizations
- Smart filtering based on `org_relationships` table
- New category in multi-model assignment

---

## Architecture

### Database Schema

```sql
-- org_relationships table defines partner relationships
CREATE TABLE org_relationships (
    parent_org_id UUID REFERENCES organizations(id),
    child_org_id UUID REFERENCES organizations(id),
    relationship_type VARCHAR(50),
    -- Types: 'channel_partner', 'sub_dealer'
);
```

### API Enhancement

#### Endpoint
```
GET /api/tickets/{ticket_id}/engineers?include_partners=true
```

#### Parameters
- `include_partners` (boolean, optional, default: false)
  - `false`: Only own engineers
  - `true`: Own + partner engineers

#### Response Structure
```json
{
  "ticket_id": "...",
  "equipment": {...},
  "ticket": {...},
  "suggestions_by_model": {
    "best_match": {
      "model_name": "Best Overall Match",
      "engineers": [
        {
          "id": "...",
          "name": "John Doe",
          "organization_id": "...",
          "organization_name": "...",
          "engineer_level": 3,
          "match_score": 90,
          ...
        }
      ],
      "count": 5
    },
    ...
  }
}
```

---

## Backend Implementation

### Files Modified

#### 1. `domain/assignment_repository.go`
```go
// Added parameter to interface
GetAvailableEngineers(
    ticketID string, 
    includePartners bool
) ([]*model.EngineerSuggestion, error)
```

#### 2. `infra/assignment_repository.go`
Enhanced SQL query with partner organization support:

```sql
-- Base query for own engineers
SELECT * FROM engineer_org_mapping eom
WHERE eom.org_id = $1

UNION

-- Additional query when include_partners=true
SELECT * FROM engineer_org_mapping eom
WHERE eom.org_id IN (
    SELECT child_org_id 
    FROM org_relationships 
    WHERE parent_org_id = $1
)
```

#### 3. `app/assignment_service.go`
```go
// Updated method signature
func (s *AssignmentService) GetEngineersForTicket(
    ticketID string,
    includePartners bool,
) ([]*model.EngineerSuggestion, error) {
    return s.repo.GetAvailableEngineers(ticketID, includePartners)
}
```

#### 4. `api/assignment_handler.go`
```go
// Parse query parameter
includePartners := r.URL.Query().Get("include_partners") == "true"

// Call service
suggestions, err := h.service.GetEngineersForTicket(ticketID, includePartners)
```

#### 5. `app/multi_model_assignment.go`
```go
// Include partners in AI model suggestions
if includePartners {
    engineers = append(engineers, partnerEngineers...)
}
```

---

## Frontend Implementation

### Files Modified

#### 1. `admin-ui/src/app/tickets/[id]/page.tsx`
```typescript
// Fetch with include_partners parameter
const engineers = await engineersApi.getForTicket(
  ticketId, 
  { include_partners: true }
);
```

#### 2. `admin-ui/src/components/EngineerSelectionModal.tsx`
```typescript
// Updated API endpoint
const response = await fetch(
  `/api/tickets/${ticketId}/engineers?include_partners=true`
);
```

#### 3. `admin-ui/src/components/MultiModelAssignment.tsx`

**Dynamic Category Creation:**
```typescript
// Get main org ID from first engineer
const mainOrgId = bestMatch?.engineers?.[0]?.organization_id;

// Filter unique partner engineers
const seenIds = new Set<string>();
const partnerEngineers: any[] = [];

Object.values(suggestions_by_model).forEach((model: any) => {
  model.engineers?.forEach((eng: any) => {
    if (eng.organization_id !== mainOrgId && !seenIds.has(eng.id)) {
      seenIds.add(eng.id);
      partnerEngineers.push(eng);
    }
  });
});

// Create partner engineers category
if (partnerEngineers.length > 0) {
  suggestions.suggestions_by_model.partner_engineers = {
    model_name: "Partner Engineers",
    description: "Engineers from partner organizations",
    engineers: partnerEngineers,
    count: partnerEngineers.length
  };
}
```

**Dynamic Sorting:**
```typescript
// Sort categories by engineer count
const sortedModelEntries = Object.entries(models).sort((a, b) => {
  const countA = a[1].count || 0;
  const countB = b[1].count || 0;
  return countB - countA; // Descending order
});
```

---

## UI Features

### 1. Partner Engineers Category

**6 Categories Total:**
1. Best Overall Match
2. Senior Engineers Only
3. Low Workload
4. Manufacturer Certified
5. Skills Match
6. **Partner Engineers** ← NEW

**Visual Design:**
- Purple active tab color (distinguishes from other categories)
- Badge shows engineer count
- Automatically sorted by count
- Empty categories disabled

### 2. Badge Color System

**High Contrast Design:**

| Organization Type | Background | Text Color |
|-------------------|------------|------------|
| Manufacturer | Blue (#3b82f6) | Blue-800 |
| Channel Partner | **Dark Orange (#ea580c)** | **White** |
| Sub-Dealer | **Dark Purple (#7e22ce)** | **White** |

**Before:** Light backgrounds with colored text (poor contrast)
**After:** Dark backgrounds with white text (high contrast, readable)

### 3. Dynamic Category Sorting

**Algorithm:**
1. Count engineers in each category
2. Sort by count (descending)
3. Move empty categories to bottom
4. Disable empty categories
5. Auto-select first non-empty category

**Benefits:**
- Most relevant categories appear first
- Users see best options immediately
- Empty categories don't clutter UI

### 4. Deduplication Logic

**Problem:** Same engineer can appear in multiple categories
**Solution:** Use Set to track seen engineer IDs

```typescript
const seenIds = new Set<string>();

categories.forEach(category => {
  category.engineers.forEach(engineer => {
    if (!seenIds.has(engineer.id)) {
      seenIds.add(engineer.id);
      uniqueEngineers.push(engineer);
    }
  });
});
```

---

## Usage Examples

### Example 1: Manufacturer Assigning Partner Engineer

**Scenario:**
- Manufacturer: Siemens Healthineers India
- Ticket: TKT-20260202-091553
- Partner: Local Dealer Z
- Engineer: Suresh Gupta (from Local Dealer Z)

**Steps:**
1. Open ticket detail page
2. Click "Multi-Model Assignment"
3. See "Partner Engineers (1)" category
4. Select Suresh Gupta
5. Assign to ticket

**Result:**
- Ticket shows: "Assigned to: Suresh Gupta"
- Organization: Local Dealer Z
- Badge: Purple "Sub-Dealer"

### Example 2: No Partner Engineers Available

**Scenario:**
- Manufacturer has no partner organizations
- OR partners have no available engineers

**Behavior:**
- Partner Engineers category not created
- Only 5 categories shown
- No UI clutter from empty category

---

## Testing

### Backend Tests

```bash
# Test 1: Without include_partners (default)
curl http://localhost:8080/api/tickets/TKT-XXX/engineers

# Expected: Only manufacturer's engineers

# Test 2: With include_partners=true
curl http://localhost:8080/api/tickets/TKT-XXX/engineers?include_partners=true

# Expected: Manufacturer + partner engineers
```

### Database Verification

```sql
-- Check org relationships
SELECT 
    parent.name as manufacturer,
    child.name as partner,
    r.relationship_type
FROM org_relationships r
JOIN organizations parent ON r.parent_org_id = parent.id
JOIN organizations child ON r.child_org_id = child.id;

-- Check engineer assignments
SELECT 
    e.name as engineer,
    o.name as organization,
    o.org_type
FROM users e
JOIN engineer_org_mapping eom ON e.id = eom.engineer_id
JOIN organizations o ON eom.org_id = o.id
WHERE e.role = 'engineer';
```

### Frontend Tests

1. **Partner Category Creation**
   - [ ] Category appears when partners exist
   - [ ] Correct count in badge
   - [ ] Purple color when active

2. **Dynamic Sorting**
   - [ ] Categories sorted by count
   - [ ] Empty categories at bottom
   - [ ] Empty categories disabled

3. **Deduplication**
   - [ ] No duplicate engineers
   - [ ] Each engineer appears once
   - [ ] Set tracking works correctly

4. **Badge Colors**
   - [ ] Channel Partner: Orange background + white text
   - [ ] Sub-Dealer: Purple background + white text
   - [ ] High contrast and readable

---

## Configuration

### Database Setup

```sql
-- Example: Create partner relationship
INSERT INTO org_relationships (
    parent_org_id,
    child_org_id,
    relationship_type
) VALUES (
    '<manufacturer-org-id>',
    '<partner-org-id>',
    'channel_partner'  -- or 'sub_dealer'
);
```

### Frontend Configuration

No configuration needed - feature is automatic when:
1. Backend returns partner engineers
2. Engineers have different `organization_id` than main org

---

## Performance Considerations

### Database
- Uses indexed columns (`org_id`, `parent_org_id`, `child_org_id`)
- UNION query is efficient with proper indexes
- Consider adding composite index on `org_relationships`

### Frontend
- Set-based deduplication is O(n)
- Sorting is O(n log n) - acceptable for typical counts
- Category creation only happens once per ticket load

---

## Security

### Access Control
- Only authenticated users can access engineers API
- JWT token verifies user's organization
- Query filters ensure only authorized partners are included

### Data Privacy
- Engineers from non-partner organizations are never exposed
- Relationship must exist in `org_relationships` table
- No transitive relationships (parent-child-grandchild)

---

## Future Enhancements

### Potential Improvements
1. **Nested Relationships**
   - Support multi-level partner hierarchies
   - Manufacturer → Channel Partner → Sub-Dealer

2. **Partner Preferences**
   - Preferred partners for specific equipment
   - Partner rating system

3. **Capacity Management**
   - Track partner engineer availability
   - Load balancing across partners

4. **Analytics**
   - Partner performance metrics
   - Assignment success rates
   - Response time tracking

---

## Troubleshooting

### Issue: No Partner Engineers Showing

**Check:**
1. Org relationships exist in database
2. Frontend uses `include_partners=true`
3. Partner engineers are active
4. Partner organizations have correct `org_type`

**SQL Debug:**
```sql
-- Verify relationships
SELECT * FROM org_relationships 
WHERE parent_org_id = '<manufacturer-id>';

-- Verify partner engineers
SELECT e.*, o.name as org_name
FROM users e
JOIN engineer_org_mapping eom ON e.id = eom.engineer_id  
JOIN organizations o ON eom.org_id = o.id
WHERE eom.org_id IN (
    SELECT child_org_id FROM org_relationships 
    WHERE parent_org_id = '<manufacturer-id>'
);
```

### Issue: Duplicate Engineers

**Solution:**
- Check `cleanMatchReason()` deduplication logic
- Verify Set is being used correctly
- Ensure engineer IDs are unique

### Issue: Category Not Sorted

**Solution:**
- Check `sortedModelEntries` implementation
- Verify `count` field exists on all categories
- Ensure sorting happens before render

---

## Related Documentation

- [Engineer Assignment System](./ENGINEER-ASSIGNMENT-README.md)
- [Organization Management](./ENGINEER-MANAGEMENT-FOR-ALL-ORGS.md)
- [Multi-Model Assignment](./PHASE4-ENGINEER-SELECTION-ENHANCEMENT.md)

---

## Changelog

**2026-02-05:** Initial implementation
- Backend API parameter added
- Frontend category created
- Dynamic sorting implemented
- Badge colors improved
- Documentation created
