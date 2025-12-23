# Equipment Registry Schema Fix - Column Mismatch

## Issue
After fixing the table name from `equipment` to `equipment_registry`, got new error:
```
Failed to get suggestions: equipment not found: failed to get equipment by ID: 
ERROR: column "qr_code_image" does not exist (SQLSTATE 42703)
```

---

## Root Cause

The equipment repository's SELECT query was written for the `equipment` table schema, but we're now querying the `equipment_registry` table which has a different schema.

### Schema Differences

**`equipment` table has:**
```sql
qr_code                     varchar(100)
qr_code_url                 text
qr_code_image               bytea          ← EXISTS
qr_code_format              varchar(10)    ← EXISTS
qr_code_generated_at        timestamp      ← EXISTS
```

**`equipment_registry` table has:**
```sql
qr_code                     varchar(255)
qr_code_url                 text
-- NO qr_code_image
-- NO qr_code_format
-- NO qr_code_generated_at
```

### The Problem

**Repository SELECT query included:**
```go
const equipmentSelectColumns = `
    ...
    qr_code_image,              // ❌ Doesn't exist in equipment_registry
    qr_code_format,             // ❌ Doesn't exist in equipment_registry
    qr_code_generated_at        // ❌ Doesn't exist in equipment_registry
`
```

**And scan functions tried to scan them:**
```go
err := row.Scan(
    ...
    &equipment.QRCodeImage,        // ❌ Column doesn't exist
    &equipment.QRCodeFormat,       // ❌ Column doesn't exist
    &equipment.QRCodeGeneratedAt,  // ❌ Column doesn't exist
)
```

---

## Solution

Removed the three non-existent columns from SELECT query and scan functions.

### Changes Made

**File:** `internal/service-domain/equipment-registry/infra/repository.go`

#### 1. Updated equipmentSelectColumns constant

**Before:**
```go
const equipmentSelectColumns = `
    ...
    COALESCE(created_by,'') AS created_by,
    qr_code_image,
    COALESCE(qr_code_format,'png') AS qr_code_format,
    qr_code_generated_at`
```

**After:**
```go
// Note: equipment_registry table doesn't have qr_code_image, qr_code_format, qr_code_generated_at
const equipmentSelectColumns = `
    ...
    COALESCE(created_by,'') AS created_by`
```

**Removed:**
- `qr_code_image`
- `qr_code_format`
- `qr_code_generated_at`

#### 2. Updated scanEquipment function

**Before (33 fields):**
```go
err := row.Scan(
    &equipment.ID,
    &equipment.QRCode,
    ...
    &equipment.CreatedBy,
    &equipment.QRCodeImage,        // ❌ Removed
    &equipment.QRCodeFormat,       // ❌ Removed
    &equipment.QRCodeGeneratedAt,  // ❌ Removed
)
```

**After (30 fields):**
```go
err := row.Scan(
    &equipment.ID,
    &equipment.QRCode,
    ...
    &equipment.CreatedBy,
)
```

#### 3. Updated scanEquipmentFromRows function

**Before (33 fields):**
```go
err := rows.Scan(
    ...
    &equipment.QRCodeImage,        // ❌ Removed
    &equipment.QRCodeFormat,       // ❌ Removed
    &equipment.QRCodeGeneratedAt,  // ❌ Removed
)
log.Printf("[ERROR] Number of columns: %d, Number of scan destinations: 33", ...)
```

**After (30 fields):**
```go
err := rows.Scan(
    ...
    &equipment.CreatedBy,
)
log.Printf("[ERROR] Number of columns: %d, Number of scan destinations: 30", ...)
```

---

## Equipment Registry vs Equipment Table

### Purpose Difference

**`equipment_registry` (Main Table - 73 rows)**
- **Purpose:** Tracks installed/registered equipment at customer sites
- **Used by:** Service tickets, QR codes, maintenance tracking
- **IDs:** REG-* format (e.g., REG-VENT-SAV-001)
- **Key fields:** installation_location, manufacturer_id, customer details
- **References:**
  - Referenced by `service_tickets.equipment_id`
  - References `organizations.id` for manufacturer
  - References `equipment_catalog.id` for catalog entry

**`equipment` (Marketplace Table - 10 rows)**
- **Purpose:** Equipment catalog for marketplace/inventory
- **Used by:** Equipment marketplace features, QR generation tools
- **IDs:** Various formats
- **Key fields:** QR code image storage, catalog information
- **Note:** Has QR code image fields for storing generated QR codes

### Why Two Tables?

1. **Different domains:** Registry (installed equipment) vs Catalog (available equipment)
2. **Different schemas:** Registry focuses on installation/service, Catalog on product details
3. **Different scale:** Registry grows with installations, Catalog is relatively static

---

## All Fixes Applied

### Summary of Engineer Reassignment Fixes

1. **First Issue: Wrong Table**
   - Changed: `FROM equipment` → `FROM equipment_registry`
   - File: `repository.go` line 143

2. **Second Issue: Wrong Columns**
   - Removed: `qr_code_image`, `qr_code_format`, `qr_code_generated_at`
   - Updated: `equipmentSelectColumns` constant
   - Updated: `scanEquipment()` function
   - Updated: `scanEquipmentFromRows()` function

---

## Equipment Registry Schema

### All Columns in equipment_registry

```sql
CREATE TABLE equipment_registry (
    -- Identity
    id                    varchar(32) PRIMARY KEY,
    qr_code               varchar(255) NOT NULL UNIQUE,
    serial_number         varchar(255) NOT NULL UNIQUE,
    equipment_id          varchar(32),
    
    -- Equipment Details
    equipment_name        varchar(500) NOT NULL,
    manufacturer_name     varchar(255) NOT NULL,
    model_number          varchar(255),
    category              varchar(255),
    
    -- Customer Details
    customer_id           varchar(32),
    customer_name         varchar(500) NOT NULL,
    
    -- Installation Details
    installation_location text,
    installation_address  jsonb,
    installation_date     date,
    
    -- Contract Details
    contract_id           varchar(32),
    purchase_date         date,
    purchase_price        numeric(15,2),
    warranty_expiry       date,
    amc_contract_id       varchar(32),
    
    -- Status & Service
    status                varchar(50) NOT NULL DEFAULT 'operational',
    last_service_date     date,
    next_service_date     date,
    service_count         integer NOT NULL DEFAULT 0,
    
    -- Additional Info
    specifications        jsonb DEFAULT '{}'::jsonb,
    photos                jsonb DEFAULT '[]'::jsonb,
    documents             jsonb DEFAULT '[]'::jsonb,
    qr_code_url           text NOT NULL,
    notes                 text,
    
    -- Metadata
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    created_by            varchar(255) NOT NULL,
    
    -- Foreign Keys
    equipment_catalog_id  uuid REFERENCES equipment_catalog(id),
    manufacturer_id       uuid REFERENCES organizations(id)
);
```

**Total:** 32 columns (not 35 with the removed QR fields)

---

## Testing

### Backend Compilation
```bash
go build -o backend.exe ./cmd/platform
```
✅ **Result:** Compiled successfully

### Backend Restart
```bash
Stop-Process -Name "backend" -Force
Start-Process backend.exe
```
✅ **Result:** Restarted successfully

### Expected API Behavior

**Test Engineer Suggestions:**
```http
POST /api/v1/tickets/346Pvjnwx3HArBM69q0ggFkBL6B/suggest-engineers
X-Tenant-ID: default
```

**Should now:**
1. ✅ Get ticket from service_tickets
2. ✅ Get equipment from equipment_registry (correct table)
3. ✅ Scan all 30 columns (correct count)
4. ✅ Return equipment details
5. ✅ Load engineer suggestions
6. ✅ Return ranked engineers

**Should NOT:**
- ❌ Error: "column qr_code_image does not exist"
- ❌ Error: "equipment not found"

---

## Files Modified

1. ✅ `internal/service-domain/equipment-registry/infra/repository.go`
   - Updated `equipmentSelectColumns` constant (removed 3 columns)
   - Updated `scanEquipment()` function (removed 3 scan fields)
   - Updated `scanEquipmentFromRows()` function (removed 3 scan fields)

---

## Status

✅ **Table name fixed** - equipment → equipment_registry  
✅ **Columns fixed** - Removed qr_code_image, qr_code_format, qr_code_generated_at  
✅ **Scan functions updated** - Match new column count (30)  
✅ **Backend rebuilt** - New code compiled  
✅ **Backend restarted** - Running with all fixes  

⏳ **Testing needed** - User to verify engineer suggestions now load  

---

## Test Checklist

### Test Engineer Reassignment
1. [ ] Go to: http://localhost:3000/tickets/346Pvjnwx3HArBM69q0ggFkBL6B
2. [ ] Click "Reassign" button
3. [ ] Modal should open with "Loading..."
4. [ ] Should load engineer list (no column error)
5. [ ] Should show engineers with skills/details
6. [ ] Can select and reassign engineer

### Verify Equipment Data
1. [ ] Equipment details should display in modal
2. [ ] Manufacturer name should be correct
3. [ ] Equipment category should match (Patient Monitor, etc.)

---

## Next Steps

**Engineer suggestions should now load correctly!** 🎉

Test by:
1. Opening any ticket details page
2. Clicking "Reassign" 
3. Verifying engineer list loads without errors
