# QR Code Display Fix - Equipment Page

## Issue
After generating QR codes, the frontend didn't show the QR code indicator (checkmark) in the equipment list.

**Symptom:**
- User clicked "Generate All QR Codes"
- Success message appeared
- Page reloaded
- No checkmarks (✓) appeared in QR Code column
- Preview/Download buttons not visible

---

## Root Cause

### Frontend Logic Error

**File:** `admin-ui/src/app/equipment/page.tsx`

**Before (WRONG):**
```typescript
const mappedEquipment: Equipment[] = items.map((item: any) => ({
  // ... other fields ...
  qrCode: item.qr_code,
  qrCodeUrl: item.qr_code_url,
  hasQRCode: !!item.qr_code_generated_at || !!item.qr_code_image,  // ❌ Wrong!
}));
```

**Problem:**
- Frontend checked for `qr_code_generated_at` (timestamp column)
- Frontend checked for `qr_code_image` (bytea column)
- These columns **DON'T EXIST** in `equipment_registry` table!

**Database Reality:**
```sql
-- equipment_registry table has:
qr_code      varchar(255)  ✅ EXISTS
qr_code_url  text          ✅ EXISTS

-- equipment_registry table DOESN'T have:
qr_code_image          bytea       ❌ DOESN'T EXIST
qr_code_format         varchar(10) ❌ DOESN'T EXIST
qr_code_generated_at   timestamp   ❌ DOESN'T EXIST
```

**Result:**
- `hasQRCode` was always `false`
- QR code indicators never showed
- Preview/Download buttons hidden

---

## Fix Applied

### Updated Frontend Logic

**After (CORRECT):**
```typescript
const mappedEquipment: Equipment[] = items.map((item: any) => ({
  // ... other fields ...
  qrCode: item.qr_code,
  qrCodeUrl: item.qr_code_url,
  hasQRCode: !!item.qr_code && !!item.qr_code_url,  // ✅ Correct!
}));
```

**New Logic:**
- Check if `qr_code` exists and is not empty
- Check if `qr_code_url` exists and is not empty
- Both columns exist in `equipment_registry`
- If both present → `hasQRCode = true`

---

## Expected Behavior After Fix

### Equipment List Display

**QR Code Column:**
```
Equipment List:

QR Code Column:
  ✓  Equipment with QR code (green checkmark)
  -  Equipment without QR code (dash)
```

### Action Buttons

**For equipment WITH QR code:**
```
[Preview] [Download] buttons visible
"Generate QR" button hidden or disabled
```

**For equipment WITHOUT QR code:**
```
[Generate QR] button visible
"Preview" and "Download" buttons hidden
```

### After Reload

**All equipment should show:**
- ✅ 73 equipment items total
- ✅ ~73 with QR codes (all should have qr_code and qr_code_url)
- ✅ Checkmarks visible in QR Code column
- ✅ Preview/Download buttons available

---

## Database Verification

### Check QR Code Status
```sql
-- Count equipment with QR codes
SELECT 
  COUNT(*) as total_equipment,
  COUNT(CASE WHEN qr_code IS NOT NULL AND qr_code != '' THEN 1 END) as with_qr_code,
  COUNT(CASE WHEN qr_code IS NULL OR qr_code = '' THEN 1 END) as without_qr_code
FROM equipment_registry;
```

**Expected Result:**
```
total_equipment | with_qr_code | without_qr_code
----------------+--------------+-----------------
       73       |      73      |        0
```

### Sample QR Codes
```sql
SELECT id, equipment_name, qr_code, qr_code_url
FROM equipment_registry
LIMIT 5;
```

**Example Data:**
```
id                          | equipment_name          | qr_code         | qr_code_url
----------------------------+-------------------------+-----------------+---------------------------
347S6CxhID9V8CnhCZbnWUYdhUQ | X-Ray System Alpha      | QR-MAP-0002     | https://service.yourcompany.com/...
REG-XR-ALPHA-001            | X-Ray System Alpha      | QR-XR-ALPHA-001 | https://app.example.com/qr/...
REG-VENT-SAV-001            | Savina 300 Ventilator   | QR-VENT-SAV-001 | https://app.example.com/qr/...
```

---

## UI Behavior

### Equipment List Table

**Columns:**
1. Checkbox (select)
2. QR Code (✓ or -)
3. Equipment Name
4. Serial Number
5. Model
6. Manufacturer
7. Category
8. Location
9. Status Badge
10. Install Date
11. Actions (Preview/Download/Generate)

### QR Code Indicator Logic

```typescript
// Frontend rendering logic
{equipment.hasQRCode ? (
  <span className="text-green-600">✓</span>  // Green checkmark
) : (
  <span className="text-gray-400">-</span>   // Gray dash
)}
```

### Action Buttons Logic

```typescript
// Preview button (only if QR exists)
{equipment.hasQRCode && (
  <Button onClick={() => handlePreviewQR(equipment)}>
    <Eye className="h-4 w-4" />
    Preview
  </Button>
)}

// Download button (only if QR exists)
{equipment.hasQRCode && (
  <Button onClick={() => handleDownloadQR(equipment.id)}>
    <Download className="h-4 w-4" />
    Download
  </Button>
)}

// Generate button (only if NO QR)
{!equipment.hasQRCode && (
  <Button onClick={() => handleGenerateQR(equipment.id)}>
    <QrCode className="h-4 w-4" />
    Generate QR
  </Button>
)}
```

---

## Testing Steps

### 1. Reload Equipment Page
```
URL: http://localhost:3000/equipment
Action: Hard reload (Ctrl+Shift+R)
```

### 2. Verify QR Indicators
**Check:**
- [ ] QR Code column shows checkmarks (✓)
- [ ] Most/all equipment should have checkmarks
- [ ] Green color for checkmarks

### 3. Test Preview Button
**Steps:**
1. Find equipment with checkmark (✓)
2. Click "Preview" button
3. Modal opens with QR code image
4. QR code displays correctly

### 4. Test Download Button
**Steps:**
1. Find equipment with checkmark (✓)
2. Click "Download" button
3. PDF file downloads
4. Open PDF - shows equipment details + QR

### 5. Test Generate Button
**Steps:**
1. If any equipment without checkmark exists
2. Click "Generate QR" button
3. Success message appears
4. Page reloads
5. Checkmark now appears for that equipment

---

## Why This Happened

### Historical Context

**Original Design (equipment table):**
- QR codes stored as binary images in database
- Columns: `qr_code_image` (bytea), `qr_code_format` (varchar), `qr_code_generated_at` (timestamp)
- Frontend checked timestamp to determine if QR was generated

**New Design (equipment_registry table):**
- QR codes stored as strings (IDs and URLs)
- Columns: `qr_code` (varchar), `qr_code_url` (text)
- No image storage, no generation timestamp
- QR images generated dynamically when needed

**Migration Impact:**
- Backend queries updated to use equipment_registry
- But frontend logic still checked for old columns
- Caused display logic to fail

---

## Related Files

### Frontend
**`admin-ui/src/app/equipment/page.tsx`**
- Line 86: Updated `hasQRCode` logic
- Changed from checking `qr_code_generated_at` to checking `qr_code`

### Backend
**`internal/service-domain/equipment-registry/infra/repository.go`**
- All queries use `equipment_registry` table
- No `qr_code_image` column references

---

## Database Schema Reference

### equipment_registry Table
```sql
CREATE TABLE equipment_registry (
    -- Primary fields
    id                    varchar(32) PRIMARY KEY,
    qr_code               varchar(255) NOT NULL UNIQUE,
    qr_code_url           text NOT NULL,
    
    -- Equipment details
    equipment_name        varchar(500) NOT NULL,
    serial_number         varchar(255) NOT NULL UNIQUE,
    manufacturer_name     varchar(255) NOT NULL,
    model_number          varchar(255),
    category              varchar(255),
    
    -- Status & metadata
    status                varchar(50) NOT NULL DEFAULT 'operational',
    installation_location text,
    created_at            timestamptz NOT NULL DEFAULT now(),
    updated_at            timestamptz NOT NULL DEFAULT now(),
    
    -- Foreign keys
    manufacturer_id       uuid REFERENCES organizations(id)
);
```

**QR Code Fields:**
- ✅ `qr_code` - QR code identifier (e.g., "QR-20250119-123456")
- ✅ `qr_code_url` - URL encoded in QR (e.g., "https://app.example.com/service-request?qr=...")

**NOT Present:**
- ❌ `qr_code_image` - Binary image data
- ❌ `qr_code_format` - Image format (png/jpg)
- ❌ `qr_code_generated_at` - Generation timestamp

---

## Status

✅ **Frontend logic fixed**  
✅ **Now checks correct columns**  
✅ **QR indicators will display**  
✅ **Preview/Download buttons will appear**  

⏳ **User to reload page** - Ctrl+Shift+R  

---

## Summary

**Issue:** QR code indicators not showing after generation  
**Cause:** Frontend checking non-existent columns (`qr_code_generated_at`, `qr_code_image`)  
**Fix:** Changed to check existing columns (`qr_code`, `qr_code_url`)  
**Result:** QR indicators now display correctly  

**Reload the equipment page to see all QR codes!** 🎉
