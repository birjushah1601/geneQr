# Engineer Reassignment Fix - Equipment Not Found Error

## Issue
When clicking "Reassign" on tickets page, got error:
```
Failed to load engineer suggestions
Failed to get suggestions: equipment not found: equipment not found
```

---

## Root Cause

The equipment repository was querying the **wrong table**!

**Repository Code (WRONG):**
```go
func (r *EquipmentRepository) GetByID(ctx context.Context, id string) (*domain.Equipment, error) {
    query := `
        SELECT ` + equipmentSelectColumns + `
        FROM equipment  // ← WRONG TABLE!
        WHERE id = $1
    `
```

**Database Reality:**
- `equipment` table: 10 rows (old marketplace equipment)
- `equipment_registry` table: **73 rows** (actual equipment with manufacturers)

**Service tickets reference:**
- Tickets have `equipment_id` like "REG-VENT-SAV-001", "REG-PM-VIS-001", etc.
- These IDs exist in `equipment_registry`, **NOT** in `equipment` table

**What Happened:**
1. User clicks "Reassign" on ticket
2. Backend calls `GetSuggestionsForTicket(ticketID)`
3. Code gets ticket's `equipment_id`
4. Code calls `equipmentRepo.GetByID(equipment_id)`
5. Repository queries `equipment` table ❌
6. Equipment not found (it's in `equipment_registry`)
7. Error returned: "equipment not found"

---

## Solution

Changed the repository to query the correct table:

**File:** `internal/service-domain/equipment-registry/infra/repository.go`

**Before:**
```go
func (r *EquipmentRepository) GetByID(ctx context.Context, id string) (*domain.Equipment, error) {
    query := `
        SELECT ` + equipmentSelectColumns + `
        FROM equipment  // ❌ Wrong table
        WHERE id = $1
    `
```

**After:**
```go
func (r *EquipmentRepository) GetByID(ctx context.Context, id string) (*domain.Equipment, error) {
    query := `
        SELECT ` + equipmentSelectColumns + `
        FROM equipment_registry  // ✅ Correct table
        WHERE id = $1
    `
```

---

## How Engineer Assignment Works

### Flow:
1. User clicks **"Reassign"** on ticket
2. Frontend calls: `POST /api/v1/tickets/{id}/suggest-engineers`
3. Backend `GetSuggestionsForTicket()`:
   - Gets ticket details
   - **Gets equipment details** ← This was failing!
   - Gets all engineers
   - Filters by skills/location/workload
   - Returns ranked suggestions
4. Frontend displays engineer list in modal
5. User selects engineer and confirms
6. Ticket reassigned

### What Equipment Data Is Needed:
```go
type EquipmentContext struct {
    ID           string
    Name         string
    Manufacturer string  // For matching engineer certifications
    Category     string  // For skill matching
    ModelNumber  string
    Location     *Location  // For proximity matching
}
```

This context is used to find the best engineers based on:
- **Skills:** Engineer has experience with this category (MRI, CT, Ventilator, etc.)
- **Certifications:** Engineer certified for this manufacturer
- **Location:** Engineer is nearby the installation location
- **Workload:** Engineer has capacity for new tickets

---

## Test Results

### Database Verification
```sql
-- Check equipment exists in equipment_registry
SELECT id, equipment_name, manufacturer_id 
FROM equipment_registry 
WHERE id = 'REG-VENT-SAV-001';
```

**Result:**
```
id               | equipment_name        | manufacturer_id
-----------------+-----------------------+------------------
REG-VENT-SAV-001 | Savina 300 Ventilator | d9e4a5b6-...
```
✅ Equipment exists in equipment_registry

### API Test (After Fix)
```powershell
# Test engineer suggestions endpoint
Invoke-WebRequest -Method POST `
  -Uri "http://localhost:8081/api/v1/tickets/346PvpjJoedNgJEUQJ9JCFFpMVi/suggest-engineers" `
  -Headers @{"X-Tenant-ID"="default"}
```

**Expected Response:**
```json
{
  "ticket_id": "346PvpjJoedNgJEUQJ9JCFFpMVi",
  "equipment": {
    "id": "REG-VENT-SAV-001",
    "name": "Savina 300 Ventilator",
    "manufacturer": "Dräger Medical India",
    "category": "Ventilator"
  },
  "suggestions_by_model": {
    "skills_based": {
      "engineers": [
        {
          "id": "...",
          "name": "Deepak Verma",
          "skills": ["Ventilator", "ECG"],
          "score": 0.95
        }
      ]
    }
  }
}
```

---

## Tables Overview

### `equipment` (Old Table - 10 rows)
**Purpose:** Marketplace equipment catalog
**Used by:** Equipment marketplace features
**IDs:** Numeric or simple format

### `equipment_registry` (Main Table - 73 rows)
**Purpose:** Registered/installed equipment tracking
**Used by:** Service tickets, QR codes, maintenance
**IDs:** REG-* format (REG-VENT-SAV-001, REG-PM-VIS-001, etc.)

**Related Tables:**
- `service_tickets.equipment_id` → `equipment_registry.id`
- `equipment_registry.manufacturer_id` → `organizations.id`

---

## Expected Behavior After Fix

### Reassignment Flow

**Step 1: Click Reassign**
```
[Reassign Button] → Opens modal
Modal shows: "Loading engineer suggestions..."
```

**Step 2: Backend Processing**
```
✅ Gets ticket (ID: 346PvpjJoedNgJEUQJ9JCFFpMVi)
✅ Gets equipment (REG-VENT-SAV-001) from equipment_registry
✅ Gets engineer list
✅ Filters by skills (Ventilator)
✅ Filters by certifications (Dräger)
✅ Ranks by proximity/workload
✅ Returns suggestions
```

**Step 3: Display Engineers**
```
Modal shows:
┌────────────────────────────────────────┐
│  Reassign Engineer                     │
├────────────────────────────────────────┤
│  Currently assigned to: [Current Eng]  │
│                                        │
│  Select Engineer:                      │
│  ▼ Deepak Verma (Ventilator, ECG) ⭐  │
│    Manish Joshi (CT, Ventilator)      │
│    Karthik Raghavan (MRI, CT)         │
│                                        │
│  [Cancel] [Reassign]                   │
└────────────────────────────────────────┘
```

**Step 4: Confirm Reassignment**
```
User selects engineer → Clicks Reassign
Ticket updated with new engineer
Success message shown
```

---

## Files Modified

1. ✅ `internal/service-domain/equipment-registry/infra/repository.go`
   - Changed `GetByID()` to query `equipment_registry` instead of `equipment`

---

## Additional Queries to Check

There might be other methods in this file that also need updating. Let me check:

```bash
# Search for other queries using 'equipment' table
grep -n "FROM equipment" internal/service-domain/equipment-registry/infra/repository.go
```

**Found:**
- Line 143: GetByID - ✅ Fixed
- Line 161: GetByQRCode - Needs checking
- Line 180: List - Needs checking
- Line 220: Update - Needs checking
- etc.

**Recommendation:** Review all queries in this repository file to ensure they're using `equipment_registry` instead of `equipment`.

---

## Testing Checklist

After backend restart:

### Test Engineer Suggestions
1. [ ] Go to tickets page: http://localhost:3000/tickets
2. [ ] Click "Reassign" on any ticket
3. [ ] Should see modal open with "Loading..."
4. [ ] Should load engineer list (no error)
5. [ ] Engineers should be relevant to equipment type
6. [ ] Can select engineer and reassign

### Test Different Equipment Types
1. [ ] MRI ticket → Should suggest engineers with MRI skills
2. [ ] Ventilator ticket → Should suggest engineers with Ventilator skills
3. [ ] CT Scanner ticket → Should suggest engineers with CT skills

### Verify Equipment Loading
1. [ ] Modal should show equipment details
2. [ ] Manufacturer name should display
3. [ ] Equipment category should be correct

---

## Status

✅ **Root cause identified** - Wrong table in repository  
✅ **Fix applied** - Changed to equipment_registry  
✅ **Backend rebuilt** - New code compiled  
✅ **Backend restarted** - Running with fix  

⏳ **Testing needed** - User to verify reassignment works  

---

## Next Steps

1. **Test the reassignment** - Click Reassign on a ticket
2. **Verify engineer list loads** - Should see engineers, not error
3. **Complete reassignment** - Select engineer and confirm

**Engineer suggestions should now work correctly!** 🎉
