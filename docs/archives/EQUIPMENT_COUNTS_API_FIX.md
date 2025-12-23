# Equipment Counts Now Showing on Manufacturers Page

## Issue
Manufacturers page was showing `0 equipment` for all manufacturers even though database has 10 equipment per manufacturer (70 total).

## Root Cause
Frontend was hardcoded to `equipmentCount: 0` with a TODO comment. No backend endpoint existed to fetch equipment counts.

## Solution

### 1. Backend API Enhancement

**Added equipment count support to organizations API:**

#### File: `internal/core/organizations/infra/repository.go`

Added `EquipmentCount` field to Organization struct:
```go
type Organization struct {
    ID             string          `json:"id"`
    Name           string          `json:"name"`
    OrgType        string          `json:"org_type"`
    Status         string          `json:"status"`
    Metadata       json.RawMessage `json:"metadata"`
    EquipmentCount int             `json:"equipment_count,omitempty"` // NEW
}
```

Added repository method:
```go
func (r *Repository) GetEquipmentCount(ctx context.Context, manufacturerID string) (int, error) {
    var count int
    err := r.db.QueryRow(ctx, 
        `SELECT COUNT(*) FROM equipment_registry WHERE manufacturer_id = $1`, 
        manufacturerID).Scan(&count)
    return count, err
}
```

#### File: `internal/core/organizations/api/handler.go`

Enhanced ListOrgs handler to support `include_counts` parameter:
```go
func (h *Handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    includeCounts := r.URL.Query().Get("include_counts") == "true"
    
    items, err := h.repo.ListOrgs(ctx, limit, offset, orgType, status)
    
    // If include_counts is requested, fetch equipment counts
    if includeCounts && orgType == "manufacturer" {
        for i := range items {
            count, _ := h.repo.GetEquipmentCount(ctx, items[i].ID)
            items[i].EquipmentCount = count
        }
    }
    
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}
```

### 2. Frontend Update

#### File: `admin-ui/src/app/manufacturers/page.tsx`

**Before:**
```typescript
queryFn: () => organizationsApi.list({ type: 'manufacturer', limit: 1000 }),

// Transform
equipmentCount: 0, // TODO: Get from equipment count API
```

**After:**
```typescript
queryFn: async () => {
  const response = await fetch('/api/v1/organizations?type=manufacturer&include_counts=true&limit=1000', {
    headers: { 'X-Tenant-ID': 'default' }
  });
  const data = await response.json();
  return data.items || [];
},

// Transform
equipmentCount: org.equipment_count || 0,
```

## API Usage

### Get Manufacturers WITHOUT Counts (Fast)
```http
GET /api/v1/organizations?type=manufacturer
X-Tenant-ID: default

Response:
{
  "items": [
    {
      "id": "uuid",
      "name": "Siemens Healthineers India",
      "org_type": "manufacturer",
      "status": "active"
      // No equipment_count
    }
  ]
}
```

### Get Manufacturers WITH Counts (Includes Count)
```http
GET /api/v1/organizations?type=manufacturer&include_counts=true
X-Tenant-ID: default

Response:
{
  "items": [
    {
      "id": "uuid",
      "name": "Siemens Healthineers India",
      "org_type": "manufacturer",
      "status": "active",
      "equipment_count": 10  // ← NEW!
    }
  ]
}
```

## Expected Results After Frontend Restart

### Manufacturers Page (`/manufacturers`)

**Equipment counts will show:**

| Manufacturer | Equipment Count |
|--------------|-----------------|
| Siemens Healthineers India | **10** |
| Wipro GE Healthcare | **10** |
| Philips Healthcare India | **10** |
| Medtronic India | **10** |
| Dräger Medical India | **10** |
| Fresenius Medical Care India | **10** |
| Canon Medical Systems India | **10** |
| Global Manufacturer A | **0** |

**Total Equipment Card:** 70 items

## Database Verification

```sql
-- Verify counts in database
SELECT 
    o.name as manufacturer,
    COUNT(er.id) as equipment_count
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.name
ORDER BY equipment_count DESC;
```

**Database has:**
- 7 manufacturers with 10 equipment each
- 1 manufacturer with 0 equipment (placeholder)
- **Total: 70 equipment items**

## Testing

### Test Backend API
```powershell
# With counts
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/organizations?type=manufacturer&include_counts=true" `
  -Headers @{"X-Tenant-ID"="default"} | ConvertFrom-Json | 
  Select-Object -ExpandProperty items | 
  Select-Object name, equipment_count
```

### Test Frontend
After restart:
```
Visit: http://localhost:3000/manufacturers
Expected: Each manufacturer shows equipment count (10 for most)
```

## Benefits

### 1. Performance
- `include_counts=false` (default) - Fast list without joins
- `include_counts=true` - Includes counts when needed
- Frontend can choose based on needs

### 2. Accurate Data
- Counts come directly from database
- No hardcoded values
- Real-time accuracy

### 3. Scalable
- Works for any number of manufacturers
- Works for any number of equipment
- No manual updates needed

## Files Modified

1. **Backend:**
   - `internal/core/organizations/infra/repository.go`
     - Added EquipmentCount field to struct
     - Added GetEquipmentCount method
   
   - `internal/core/organizations/api/handler.go`
     - Added include_counts parameter support
     - Fetch counts when requested

2. **Frontend:**
   - `admin-ui/src/app/manufacturers/page.tsx`
     - Updated query to use include_counts=true
     - Use equipment_count from API response

## Status

✅ **Backend API enhanced** - Supports equipment counts
✅ **Backend rebuilt and restarted** - Running with new code
✅ **API tested** - Returns correct counts (10 per manufacturer)
✅ **Frontend updated** - Fetches with include_counts=true
✅ **Frontend restart needed** - To see equipment counts

## Result

After frontend restart, the manufacturers page will display:
- ✅ Real equipment counts (10 per manufacturer)
- ✅ Total equipment card shows 70
- ✅ Each manufacturer card shows correct count
- ✅ Data updates automatically from database

**No more 0 values - real counts from database!** 🎉
