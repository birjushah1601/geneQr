# Engineers Count API - Complete Implementation

## Issue
Manufacturer dashboard showing **Engineers: 0** even though database has engineers assigned.

## Solution Implemented

### Backend API Enhancement

#### 1. Added GetEngineersCount Method

**File:** `internal/core/organizations/infra/repository.go`

```go
func (r *Repository) GetEngineersCount(ctx context.Context, organizationID string) (int, error) {
    var count int
    err := r.db.QueryRow(ctx, 
        `SELECT COUNT(DISTINCT engineer_id) 
         FROM engineer_org_memberships 
         WHERE org_id = $1`, 
        organizationID).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}
```

#### 2. Added EngineersCount Field

**File:** `internal/core/organizations/infra/repository.go`

```go
type Organization struct {
    ID             string          `json:"id"`
    Name           string          `json:"name"`
    OrgType        string          `json:"org_type"`
    Status         string          `json:"status"`
    Metadata       json.RawMessage `json:"metadata"`
    EquipmentCount int             `json:"equipment_count,omitempty"`
    EngineersCount int             `json:"engineers_count,omitempty"` // NEW
}
```

#### 3. Updated ListOrgs Handler

**File:** `internal/core/organizations/api/handler.go`

```go
func (h *Handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // If include_counts is requested, fetch equipment and engineer counts
    if includeCounts && orgType == "manufacturer" {
        for i := range items {
            equipmentCount, _ := h.repo.GetEquipmentCount(ctx, items[i].ID)
            items[i].EquipmentCount = equipmentCount
            
            engineersCount, _ := h.repo.GetEngineersCount(ctx, items[i].ID)
            items[i].EngineersCount = engineersCount
        }
    }
}
```

#### 4. Updated GetOrg Handler

**File:** `internal/core/organizations/api/handler.go`

```go
func (h *Handler) GetOrg(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    org, err := h.repo.GetOrgByID(ctx, id)
    
    // If include_counts parameter is set, fetch counts
    includeCounts := r.URL.Query().Get("include_counts") == "true"
    if includeCounts && org.OrgType == "manufacturer" {
        equipmentCount, _ := h.repo.GetEquipmentCount(ctx, org.ID)
        org.EquipmentCount = equipmentCount
        
        engineersCount, _ := h.repo.GetEngineersCount(ctx, org.ID)
        org.EngineersCount = engineersCount
    }
    
    h.respondJSON(w, http.StatusOK, org)
}
```

---

## Frontend Updates

### 1. Manufacturers List Page

**File:** `admin-ui/src/app/manufacturers/page.tsx`

**Before:**
```typescript
engineersCount: 0, // TODO: Get from engineers count API
```

**After:**
```typescript
engineersCount: org.engineers_count || 0,
```

### 2. Manufacturer Dashboard

**File:** `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`

**Before:**
```typescript
// Fetch engineers count (TODO: create engineers endpoint)
let engineersCount = 0;
```

**After:**
```typescript
// Fetch engineers count
let engineersCount = 0;
try {
  const engineersResponse = await fetch(
    `${apiBaseUrl}/v1/organizations/${manufacturerId}?include_counts=true`,
    { headers: { 'X-Tenant-ID': 'default' } }
  );
  if (engineersResponse.ok) {
    const engineersData = await engineersResponse.json();
    engineersCount = engineersData.engineers_count || 0;
  }
} catch (e) {
  console.error('Failed to fetch engineer count:', e);
}
```

---

## API Usage

### Get All Manufacturers with Counts

```http
GET /api/v1/organizations?type=manufacturer&include_counts=true
X-Tenant-ID: default
```

**Response:**
```json
{
  "items": [
    {
      "id": "f1c1ebfb-57fd-4307-93db-2f72e9d004ad",
      "name": "Philips Healthcare India",
      "org_type": "manufacturer",
      "status": "active",
      "equipment_count": 10,
      "engineers_count": 5,
      "metadata": { ... }
    }
  ]
}
```

### Get Single Manufacturer with Counts

```http
GET /api/v1/organizations/{id}?include_counts=true
X-Tenant-ID: default
```

**Response:**
```json
{
  "id": "f1c1ebfb-57fd-4307-93db-2f72e9d004ad",
  "name": "Philips Healthcare India",
  "org_type": "manufacturer",
  "status": "active",
  "equipment_count": 10,
  "engineers_count": 5,
  "metadata": {
    "contact_person": "Mr. Ankit Desai",
    "email": "ankit.desai@philips.com",
    ...
  }
}
```

---

## Test Results

### API Test - All Manufacturers

```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/organizations?type=manufacturer&include_counts=true" `
  -Headers @{"X-Tenant-ID"="default"} | 
  ConvertFrom-Json | 
  Select-Object -ExpandProperty items | 
  Select-Object name, equipment_count, engineers_count
```

**Result:**
```
name                         equipment_count engineers_count
----                         --------------- ---------------
Siemens Healthineers India                10               6
Wipro GE Healthcare                       10               6
Canon Medical Systems India               10               5
Philips Healthcare India                  10               5
Medtronic India                           10               4
Dräger Medical India                      10               3
Fresenius Medical Care India              10               2
Global Manufacturer A                      0               2
```

### API Test - Single Manufacturer (Philips)

```powershell
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/organizations/f1c1ebfb-57fd-4307-93db-2f72e9d004ad?include_counts=true" `
  -Headers @{"X-Tenant-ID"="default"} |
  ConvertFrom-Json |
  Select-Object name, equipment_count, engineers_count
```

**Result:**
```
name                     : Philips Healthcare India
equipment_count          : 10
engineers_count          : 5
```

---

## Expected Frontend Display

### Manufacturers List Page

After reload:
```
┌────────────────────────┬───────────┬───────────┐
│ Manufacturer           │ Equipment │ Engineers │
├────────────────────────┼───────────┼───────────┤
│ Siemens Healthineers   │    10     │     6     │
│ Wipro GE Healthcare    │    10     │     6     │
│ Canon Medical Systems  │    10     │     5     │
│ Philips Healthcare     │    10     │     5     │
│ Medtronic India        │    10     │     4     │
│ Dräger Medical         │    10     │     3     │
│ Fresenius Medical      │    10     │     2     │
│ Global Manufacturer A  │     0     │     2     │
└────────────────────────┴───────────┴───────────┘
```

### Manufacturer Dashboard (Philips Healthcare)

After reload:
```
Manufacturer Dashboard
Manage equipment, engineers, and service operations for Philips Healthcare India

┌─────────────┬─────────────┬───────────────┬──────────────┐
│ Equipment   │ Engineers   │ Active Tickets│ Member Since │
│     10      │      5      │       0       │  Dec 2025    │
│ Registered  │ Service     │ Open requests │ Partner      │
│  devices    │   team      │               │  status      │
└─────────────┴─────────────┴───────────────┴──────────────┘

Equipment Registry
10 equipment items registered. View, import, or manage equipment.

Service Engineers
5 engineers in the service team. View, add, or manage engineers.
```

---

## Engineer Distribution by Manufacturer

| Manufacturer | Engineers | Engineers List |
|--------------|-----------|----------------|
| **Siemens Healthineers** | **6** | Amit Patel, Arun Menon, Kavita Nair, Manish Joshi, Rajesh Kumar Singh, Vikram Reddy |
| **Wipro GE Healthcare** | **6** | Amit Patel, Kavita Nair, Manish Joshi, Rajesh Kumar Singh, Suresh Gupta, Vikram Reddy |
| **Canon Medical** | **5** | Amit Patel, Manish Joshi, Priya Sharma, Ravi Iyer, Suresh Gupta |
| **Philips Healthcare** | **5** | Amit Patel, Arun Menon, Kavita Nair, Rajesh Kumar Singh, Suresh Gupta |
| **Medtronic India** | **4** | Arjun Malhotra, Priya Sharma, Rajesh Kumar Singh, Shreya Patel |
| **Dräger Medical** | **3** | Deepak Verma, Karthik Raghavan, Manish Joshi |
| **Fresenius Medical** | **2** | Neha Kulkarni, Sanjay Mehta |
| **Global Manufacturer A** | **2** | Divya Krishnan, Karthik Raghavan |

---

## Database Queries Used

### Get Engineer Count for Manufacturer
```sql
SELECT COUNT(DISTINCT engineer_id) 
FROM engineer_org_memberships 
WHERE org_id = 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad';
-- Result: 5
```

### Get All Counts for Manufacturer
```sql
SELECT 
    o.name,
    COUNT(DISTINCT m.engineer_id) as engineers,
    COUNT(DISTINCT er.id) as equipment
FROM organizations o
LEFT JOIN engineer_org_memberships m ON m.org_id = o.id
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.id = 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad'
GROUP BY o.id, o.name;
```

---

## Files Modified

### Backend
1. ✅ `internal/core/organizations/infra/repository.go`
   - Added `EngineersCount` field to Organization struct
   - Added `GetEngineersCount()` method

2. ✅ `internal/core/organizations/api/handler.go`
   - Updated `ListOrgs()` to fetch engineer counts
   - Updated `GetOrg()` to support include_counts parameter

### Frontend
1. ✅ `admin-ui/src/app/manufacturers/page.tsx`
   - Use `engineers_count` from API response
   - Display on manufacturers list

2. ✅ `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`
   - Fetch organization with include_counts=true
   - Display engineer count on dashboard

---

## Performance Considerations

### Single Query per Organization
The `GetEngineersCount()` method uses:
```sql
SELECT COUNT(DISTINCT engineer_id) 
FROM engineer_org_memberships 
WHERE org_id = $1
```

- **Index exists:** `engineer_org_memberships_pkey` on (engineer_id, org_id)
- **Performance:** Fast lookup, uses index
- **Result:** Cached by database for repeated calls

### Batch Query (Alternative for Many Manufacturers)
For listing many manufacturers, could optimize with:
```sql
SELECT 
    org_id, 
    COUNT(DISTINCT engineer_id) as count
FROM engineer_org_memberships
GROUP BY org_id;
```

Current implementation is fine for 8-10 manufacturers.

---

## Status

✅ **Backend API complete** - Returns engineer counts  
✅ **Frontend updated** - Uses engineer counts from API  
✅ **Database populated** - 16 engineers assigned to 8 manufacturers  
✅ **API tested** - Returns correct counts (6, 6, 5, 5, 4, 3, 2, 2)  
✅ **Backend rebuilt** - New code compiled and running  

⏳ **Browser reload needed** - To see engineer counts on UI  

---

## Next Step

**Hard reload browser:** `Ctrl + Shift + R` (or `Cmd + Shift + R` on Mac)

Then visit:
- http://localhost:3000/manufacturers (see engineer counts in list)
- http://localhost:3000/manufacturers/f1c1ebfb-57fd-4307-93db-2f72e9d004ad/dashboard (see 5 engineers)

**Engineer counts will now display correctly!** 🎉
