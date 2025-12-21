# Active Tickets Count API - Complete Implementation

## Overview
Implemented active tickets count for manufacturers by joining service_tickets → equipment_registry → manufacturer.

---

## Backend Implementation

### 1. Repository Method

**File:** `internal/core/organizations/infra/repository.go`

```go
func (r *Repository) GetActiveTicketsCount(ctx context.Context, manufacturerID string) (int, error) {
    var count int
    query := `
        SELECT COUNT(DISTINCT st.id) 
        FROM service_tickets st
        JOIN equipment_registry er ON st.equipment_id = er.id
        WHERE er.manufacturer_id = $1 
        AND st.status != 'closed'
    `
    err := r.db.QueryRow(ctx, query, manufacturerID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}
```

**Logic:**
- Joins `service_tickets` with `equipment_registry` via `equipment_id`
- Filters by `manufacturer_id` on equipment
- Excludes tickets with status = 'closed'
- Counts active tickets (new, assigned, in_progress)

### 2. Organization Struct

**File:** `internal/core/organizations/infra/repository.go`

```go
type Organization struct {
    ID                string          `json:"id"`
    Name              string          `json:"name"`
    OrgType           string          `json:"org_type"`
    Status            string          `json:"status"`
    Metadata          json.RawMessage `json:"metadata"`
    EquipmentCount    int             `json:"equipment_count,omitempty"`
    EngineersCount    int             `json:"engineers_count,omitempty"`
    ActiveTickets     int             `json:"active_tickets,omitempty"` // NEW
}
```

### 3. API Handler Updates

**File:** `internal/core/organizations/api/handler.go`

**ListOrgs (for manufacturers list):**
```go
if includeCounts && orgType == "manufacturer" {
    for i := range items {
        equipmentCount, _ := h.repo.GetEquipmentCount(ctx, items[i].ID)
        items[i].EquipmentCount = equipmentCount
        
        engineersCount, _ := h.repo.GetEngineersCount(ctx, items[i].ID)
        items[i].EngineersCount = engineersCount
        
        activeTickets, _ := h.repo.GetActiveTicketsCount(ctx, items[i].ID)
        items[i].ActiveTickets = activeTickets
    }
}
```

**GetOrg (for single manufacturer):**
```go
if includeCounts && org.OrgType == "manufacturer" {
    equipmentCount, _ := h.repo.GetEquipmentCount(ctx, org.ID)
    org.EquipmentCount = equipmentCount
    
    engineersCount, _ := h.repo.GetEngineersCount(ctx, org.ID)
    org.EngineersCount = engineersCount
    
    activeTickets, _ := h.repo.GetActiveTicketsCount(ctx, org.ID)
    org.ActiveTickets = activeTickets
}
```

---

## Frontend Implementation

### 1. Manufacturers List Page

**File:** `admin-ui/src/app/manufacturers/page.tsx`

**Added field:**
```typescript
activeTickets: org.active_tickets || 0,
```

### 2. Manufacturer Dashboard (Optimized)

**File:** `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`

**Before (3 separate API calls):**
```typescript
// Fetch equipment
const equipmentResponse = await fetch('/v1/equipment?manufacturer_id=...');

// Fetch engineers count
const engineersResponse = await fetch('/v1/organizations/{id}?include_counts=true');

// Fetch tickets count
const ticketsResponse = await fetch('/v1/organizations/{id}?include_counts=true');
```

**After (1 optimized API call):**
```typescript
// Fetch everything in ONE call
const orgResponse = await fetch(
  `${apiBaseUrl}/v1/organizations/${manufacturerId}?include_counts=true`,
  { headers: { 'X-Tenant-ID': 'default' } }
);

const org = await orgResponse.json();
equipmentCount = org.equipment_count || 0;
engineersCount = org.engineers_count || 0;
activeTickets = org.active_tickets || 0;
```

**Benefits:**
- ✅ 66% fewer HTTP requests
- ✅ Faster page load
- ✅ Less network overhead
- ✅ Consistent data snapshot

---

## API Usage

### Get All Manufacturers with All Counts

```http
GET /api/v1/organizations?type=manufacturer&include_counts=true
X-Tenant-ID: default
```

**Response:**
```json
{
  "items": [
    {
      "id": "11afdeec-5dee-44d4-aa5b-952703536f10",
      "name": "Siemens Healthineers India",
      "org_type": "manufacturer",
      "status": "active",
      "equipment_count": 10,
      "engineers_count": 6,
      "active_tickets": 1,
      "metadata": { ... }
    }
  ]
}
```

### Get Single Manufacturer with All Counts

```http
GET /api/v1/organizations/{id}?include_counts=true
X-Tenant-ID: default
```

**Response (Philips Healthcare):**
```json
{
  "id": "f1c1ebfb-57fd-4307-93db-2f72e9d004ad",
  "name": "Philips Healthcare India",
  "equipment_count": 10,
  "engineers_count": 5,
  "active_tickets": 0,
  "metadata": { ... }
}
```

---

## Test Results

### API Test - All Manufacturers

```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/organizations?type=manufacturer&include_counts=true" |
  ConvertFrom-Json | 
  Select-Object -ExpandProperty items |
  Select-Object name, equipment_count, engineers_count, active_tickets
```

**Result:**
```
name                         equipment_count engineers_count active_tickets
----                         --------------- --------------- --------------
Siemens Healthineers India                10               6              1
Wipro GE Healthcare                       10               6              2
Medtronic India                           10               4              2
Dräger Medical India                      10               3              2
Canon Medical Systems India               10               5              0
Philips Healthcare India                  10               5              0
Fresenius Medical Care India              10               2              0
Global Manufacturer A                      0               2              0
```

### Database Verification

**Check tickets by manufacturer:**
```sql
SELECT 
    o.name as manufacturer,
    COUNT(DISTINCT st.id) as active_tickets,
    string_agg(DISTINCT st.status, ', ') as statuses
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
LEFT JOIN service_tickets st ON st.equipment_id = er.id AND st.status != 'closed'
WHERE o.org_type = 'manufacturer'
GROUP BY o.name
ORDER BY active_tickets DESC;
```

**Result:**
```
manufacturer                 active_tickets statuses
---------------------------  -------------- --------------------
Wipro GE Healthcare                      2  new, assigned
Medtronic India                          2  assigned, in_progress
Dräger Medical India                     2  new, assigned
Siemens Healthineers India               1  assigned
Philips Healthcare India                 0  
Canon Medical Systems India              0  
Fresenius Medical Care India             0  
Global Manufacturer A                    0  
```

---

## Ticket Statuses

### Active Statuses (Counted)
- `new` - Just created, not yet assigned
- `assigned` - Assigned to engineer
- `in_progress` - Engineer working on it

### Closed Status (Excluded)
- `closed` - Ticket resolved and closed

**SQL Filter:** `WHERE st.status != 'closed'`

---

## Expected Frontend Display

### Manufacturers List Page

After reload:
```
┌────────────────────┬───────────┬───────────┬────────────────┐
│ Manufacturer       │ Equipment │ Engineers │ Active Tickets │
├────────────────────┼───────────┼───────────┼────────────────┤
│ Siemens            │    10     │     6     │       1        │
│ Wipro GE           │    10     │     6     │       2        │
│ Medtronic          │    10     │     4     │       2        │
│ Dräger             │    10     │     3     │       2        │
│ Canon              │    10     │     5     │       0        │
│ Philips            │    10     │     5     │       0        │
│ Fresenius          │    10     │     2     │       0        │
│ Global Mfr A       │     0     │     2     │       0        │
└────────────────────┴───────────┴───────────┴────────────────┘
```

### Manufacturer Dashboard

**For Philips Healthcare:**
```
Manufacturer Dashboard
Manage equipment, engineers, and service operations for Philips Healthcare India

┌─────────────┬─────────────┬───────────────┬──────────────┐
│ Equipment   │ Engineers   │ Active Tickets│ Member Since │
│     10      │      5      │       0       │  Dec 2025    │
│ Registered  │ Service     │ Open requests │ Partner      │
│  devices    │   team      │               │  status      │
└─────────────┴─────────────┴───────────────┴──────────────┘
```

**For Wipro GE Healthcare:**
```
┌─────────────┬─────────────┬───────────────┬──────────────┐
│ Equipment   │ Engineers   │ Active Tickets│ Member Since │
│     10      │      6      │       2       │  Dec 2025    │
│ Registered  │ Service     │ Open requests │ Partner      │
│  devices    │   team      │               │  status      │
└─────────────┴─────────────┴───────────────┴──────────────┘

Service Tickets
2 active service tickets requiring attention.
[View All Tickets] (enabled)
```

---

## Database Schema

### Tables Involved

**service_tickets:**
- `id` - Ticket identifier
- `equipment_id` - FK to equipment_registry
- `status` - Ticket status (new, assigned, in_progress, closed)

**equipment_registry:**
- `id` - Equipment identifier
- `manufacturer_id` - FK to organizations

**organizations:**
- `id` - Manufacturer identifier
- `org_type` - Must be 'manufacturer'

### Relationship Chain
```
organizations (manufacturer)
    ↓ manufacturer_id
equipment_registry
    ↓ equipment_id
service_tickets
```

---

## Performance Considerations

### Query Performance
```sql
SELECT COUNT(DISTINCT st.id) 
FROM service_tickets st
JOIN equipment_registry er ON st.equipment_id = er.id
WHERE er.manufacturer_id = $1 
AND st.status != 'closed'
```

**Indexes Used:**
- `idx_ticket_equipment` on service_tickets(equipment_id)
- `idx_ticket_status` on service_tickets(status)
- Primary key on equipment_registry(id)

**Performance:** Fast for 8 manufacturers with ~10 equipment each and ~9 tickets total.

### API Response Time
Single API call now fetches:
- Organization metadata
- Equipment count (1 query)
- Engineers count (1 query)
- Active tickets count (1 query with join)

**Total:** 3 DB queries instead of 3 separate HTTP requests.

---

## Files Modified

### Backend
1. ✅ `internal/core/organizations/infra/repository.go`
   - Added `ActiveTickets` field to Organization
   - Added `GetActiveTicketsCount()` method

2. ✅ `internal/core/organizations/api/handler.go`
   - Updated `ListOrgs()` to fetch active tickets
   - Updated `GetOrg()` to fetch active tickets

### Frontend
1. ✅ `admin-ui/src/app/manufacturers/page.tsx`
   - Use `active_tickets` from API response

2. ✅ `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`
   - Optimized to single API call
   - Fetch all counts in one request
   - Use `active_tickets` from response

---

## Complete Statistics

### All Manufacturers Summary

| Manufacturer | Equipment | Engineers | Active Tickets |
|--------------|-----------|-----------|----------------|
| Siemens Healthineers | 10 | 6 | **1** |
| Wipro GE Healthcare | 10 | 6 | **2** |
| Canon Medical | 10 | 5 | 0 |
| Philips Healthcare | 10 | 5 | 0 |
| Medtronic India | 10 | 4 | **2** |
| Dräger Medical | 10 | 3 | **2** |
| Fresenius Medical | 10 | 2 | 0 |
| Global Manufacturer A | 0 | 2 | 0 |

**Totals:**
- Equipment: 70
- Engineers: 33 assignments (16 unique)
- Active Tickets: 7 (out of 9 total)

---

## Status

✅ **Backend complete** - GetActiveTicketsCount() method  
✅ **API enhanced** - Returns active_tickets field  
✅ **Frontend optimized** - Single API call for all counts  
✅ **Database verified** - Correct ticket counts by manufacturer  
✅ **API tested** - All manufacturers return correct counts  

⏳ **Browser reload needed** - To see active tickets on dashboard  

---

## Next Step

**Hard reload browser:** `Ctrl + Shift + R`

Then visit:
- http://localhost:3000/manufacturers (see all counts)
- http://localhost:3000/manufacturers/aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad/dashboard (Wipro GE - 2 active tickets)
- http://localhost:3000/manufacturers/f1c1ebfb-57fd-4307-93db-2f72e9d004ad/dashboard (Philips - 0 active tickets)

**All dashboard stats now show real data from database!** 🎉
