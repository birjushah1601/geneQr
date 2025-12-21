# Equipment QR Code - Troubleshooting Guide

## Issue Reported
After generating QR codes, they don't show up on the frontend.

---

## Fixes Applied

### 1. Fixed `hasQRCode` Logic ✅
**Problem:** Frontend was checking non-existent columns  
**Fix:** Now checks `qr_code` and `qr_code_url` which exist in equipment_registry

### 2. Enhanced Console Logging ✅
**Added:** Detailed logging throughout QR generation flow  
**Purpose:** Help identify where issues occur

### 3. Improved UX ✅
**Changed:** Removed alert popups (smoother experience)  
**Added:** Confirmation dialog for bulk generation  
**Added:** Success/failure counters

---

## How to Debug

### Step 1: Open Browser Console
1. Go to http://localhost:3000/equipment
2. Press **F12** to open DevTools
3. Click **Console** tab

### Step 2: Check Data Loading
**Look for:**
```
[Equipment Load] Loaded 73 equipment items (73 with QR codes)
```

**If you see:**
```
[Equipment Load] Loaded 73 equipment items (0 with QR codes)
```
→ **Problem:** Equipment doesn't have QR codes in database

**If you see:**
```
Failed to fetch equipment from API
```
→ **Problem:** Backend not responding or wrong API URL

### Step 3: Test QR Generation
1. Select one equipment item (checkbox)
2. Click "Generate All QR Codes" button
3. Confirm dialog
4. Watch console for:

```
[Bulk QR] Generating for 1 equipment items...
[Bulk QR] Generating for: {equipment-id}
[Bulk QR] Complete! Success: 1, Failed: 0
```

### Step 4: Check After Reload
After page reloads, look for:
```
[Equipment Load] Item: {
  id: "...",
  name: "...",
  qr_code: "QR-...",      ← Should have value
  qr_code_url: "https://", ← Should have value
  hasQRCode: true          ← Should be true
}
```

---

## Common Issues & Solutions

### Issue 1: QR Codes Not Showing

**Symptoms:**
- Equipment list loads
- All QR code cells empty or show "Generate" button
- Console shows: `(0 with QR codes)`

**Check Database:**
```sql
SELECT COUNT(*) as total,
       COUNT(CASE WHEN qr_code IS NOT NULL AND qr_code != '' THEN 1 END) as with_qr
FROM equipment_registry;
```

**Expected:**
```
total | with_qr
------+---------
  73  |   73
```

**If with_qr = 0:**
→ Run bulk generate to create QR codes

**Solution:**
```powershell
# Test bulk generation
Invoke-WebRequest -Method POST `
  -Uri "http://localhost:8081/api/v1/equipment/qr/bulk-generate" `
  -Headers @{"X-Tenant-ID"="default"}
```

### Issue 2: QR Images Not Displaying

**Symptoms:**
- Console shows equipment has QR codes
- But QR code images broken/not loading
- Browser console shows 404 errors

**Check Image URL:**
Open browser console → Network tab → Look for:
```
GET /api/v1/equipment/qr/image/{id}  → 404
```

**Solution:**
The QR image is generated dynamically. Check backend is running:
```powershell
# Test QR image endpoint
Invoke-WebRequest -Uri "http://localhost:8081/api/v1/equipment/qr/image/REG-CAN-XR-001" `
  -Headers @{"X-Tenant-ID"="default"} `
  -OutFile "test-qr.png"
```

Should create `test-qr.png` with QR code image.

### Issue 3: "Generate QR" Button Doesn't Work

**Symptoms:**
- Click "Generate QR" button
- Nothing happens
- No console logs

**Check:**
1. Open browser console (F12)
2. Click Generate QR
3. Look for error messages

**Common Errors:**

**A) Network Error:**
```
Failed to generate QR code: NetworkError
```
→ Backend not running or wrong URL

**B) CORS Error:**
```
CORS policy: No 'Access-Control-Allow-Origin' header
```
→ Backend CORS not configured (shouldn't happen with our setup)

**C) 500 Error:**
```
Failed to generate QR code: Internal Server Error
```
→ Check backend logs

### Issue 4: Page Doesn't Reload After Generation

**Symptoms:**
- QR generation succeeds (console shows success)
- Page doesn't reload
- QR codes still not visible

**Check Console:**
Look for:
```
[QR Generation] Success! Result: {...}
[QR Generation] Reloading page to show updated QR code...
```

**If reload doesn't happen:**
- Browser popup blocker might prevent reload
- JavaScript error preventing reload

**Manual Fix:**
Press **Ctrl+Shift+R** to hard reload

---

## Testing QR Code Functionality

### Test 1: View Equipment with QR Codes

**Steps:**
1. Load http://localhost:3000/equipment
2. Check QR Code column

**Expected:**
- QR code images displayed (80x80px thumbnails)
- Hover shows Preview/Download buttons
- Click image opens preview modal

**Console:**
```
[Equipment Load] Loaded 73 equipment items (73 with QR codes)
```

### Test 2: Generate Single QR

**Steps:**
1. Find equipment WITHOUT QR (if any)
2. Click "Generate" button in QR column
3. Wait for reload

**Expected Console:**
```
[QR Generation] Starting for equipment: {id}
[QR Generation] Success! Result: {qr_code: "QR-..."}
[QR Generation] Reloading page...
```

**After Reload:**
- QR code image appears
- Preview/Download buttons available

### Test 3: Generate Selected QR Codes

**Steps:**
1. Select 2-3 equipment items (checkboxes)
2. Click "Generate All QR Codes" button
3. Confirm dialog
4. Wait for completion

**Expected Console:**
```
[Bulk QR] Generating for 3 equipment items...
[Bulk QR] Generating for: {id-1}
[Bulk QR] Generating for: {id-2}
[Bulk QR] Generating for: {id-3}
[Bulk QR] Complete! Success: 3, Failed: 0
```

**After Reload:**
- All selected equipment show QR codes

### Test 4: Preview QR Code

**Steps:**
1. Click QR code image (or Preview button on hover)
2. Modal opens

**Expected:**
- Large QR code image (400x400px)
- Equipment details shown
- Download button available
- Close button works

### Test 5: Download QR Label

**Steps:**
1. Click Download button (or hover menu)
2. Check downloads folder

**Expected:**
- PDF file downloads
- Filename: `equipment-qr-{id}.pdf`
- PDF contains:
  - Equipment name
  - Serial number
  - Manufacturer
  - QR code image

---

## API Endpoints Reference

### List Equipment
```http
GET /api/v1/equipment?limit=1000
Headers: X-Tenant-ID: default
```

**Response:**
```json
{
  "equipment": [
    {
      "id": "REG-CAN-XR-001",
      "equipment_name": "Digital X-Ray System",
      "qr_code": "QR-CAN-XR-001",
      "qr_code_url": "https://api.qrserver.com/...",
      "serial_number": "CAN-CXDI-509001",
      "manufacturer_name": "Canon Medical Systems"
    }
  ]
}
```

### Generate QR Code
```http
POST /api/v1/equipment/{id}/qr
Headers: X-Tenant-ID: default
```

**Response:**
```json
{
  "qr_code": "QR-20250119-123456",
  "message": "QR code generated successfully"
}
```

### Get QR Code Image
```http
GET /api/v1/equipment/qr/image/{id}
Headers: X-Tenant-ID: default
```

**Response:**
- Content-Type: image/png
- Binary PNG image data

### Download QR Label (PDF)
```http
GET /api/v1/equipment/{id}/qr/pdf
Headers: X-Tenant-ID: default
```

**Response:**
- Content-Type: application/pdf
- Binary PDF file

### Bulk Generate QR Codes
```http
POST /api/v1/equipment/qr/bulk-generate
Headers: X-Tenant-ID: default
```

**Response:**
```json
{
  "total_processed": 73,
  "successful": 73,
  "failed": 0,
  "message": "Successfully generated 73 QR codes"
}
```

---

## Database Verification

### Check Equipment with QR Codes
```sql
SELECT 
  id,
  equipment_name,
  qr_code,
  CASE 
    WHEN qr_code IS NOT NULL AND qr_code != '' THEN 'HAS QR'
    ELSE 'NO QR'
  END as qr_status
FROM equipment_registry
LIMIT 10;
```

### Count QR Code Status
```sql
SELECT 
  COUNT(*) as total_equipment,
  COUNT(CASE WHEN qr_code IS NOT NULL AND qr_code != '' THEN 1 END) as with_qr,
  COUNT(CASE WHEN qr_code IS NULL OR qr_code = '' THEN 1 END) as without_qr
FROM equipment_registry;
```

### Sample QR Code Data
```sql
SELECT id, equipment_name, qr_code, qr_code_url
FROM equipment_registry
WHERE qr_code IS NOT NULL
LIMIT 5;
```

**Expected Result:**
```
id              | equipment_name          | qr_code        | qr_code_url
----------------+-------------------------+----------------+-------------
REG-CAN-XR-001  | Digital X-Ray System    | QR-CAN-XR-001  | https://...
REG-FMC-DLY-001 | Fresenius Dialysis      | QR-FMC-DLY-001 | https://...
```

---

## Quick Fixes

### Fix 1: All Equipment Missing QR Codes
```powershell
# Generate QR codes for all equipment via API
Invoke-WebRequest -Method POST `
  -Uri "http://localhost:8081/api/v1/equipment/qr/bulk-generate" `
  -Headers @{"X-Tenant-ID"="default"}
```

### Fix 2: Frontend Not Showing QR Codes
```powershell
# Hard reload browser
# Press: Ctrl+Shift+R
```

### Fix 3: Backend Not Responding
```powershell
# Check if backend is running
Get-Process | Where-Object {$_.Name -eq "backend"}

# If not running, start it
cd C:\Users\birju\aby-med
.\backend.exe
```

### Fix 4: Database Connection Issue
```powershell
# Check PostgreSQL container
docker ps | Select-String "med_platform_pg"

# If not running, start it
docker start med_platform_pg
```

---

## Expected Console Output (Complete Flow)

### On Page Load:
```
[Equipment Load] API Response: {...}
[Equipment Load] Item: {id: "...", hasQRCode: true}
[Equipment Load] Item: {id: "...", hasQRCode: true}
[Equipment Load] Item: {id: "...", hasQRCode: true}
[Equipment Load] Loaded 73 equipment items (73 with QR codes)
```

### On QR Generation:
```
[QR Generation] Starting for equipment: REG-CAN-XR-001
[QR Generation] Success! Result: {qr_code: "QR-20250119-123456"}
[QR Generation] Reloading page to show updated QR code...
```

### On Bulk Generation:
```
[Bulk QR] Generating for 5 equipment items...
[Bulk QR] Generating for: REG-CAN-XR-001
[Bulk QR] Generating for: REG-CAN-XR-002
[Bulk QR] Generating for: REG-CAN-XR-003
[Bulk QR] Generating for: REG-CAN-XR-004
[Bulk QR] Generating for: REG-CAN-XR-005
[Bulk QR] Complete! Success: 5, Failed: 0
```

---

## Contact Points for Support

### Frontend Issues
- File: `admin-ui/src/app/equipment/page.tsx`
- Check: Console logs, Network tab in DevTools
- Debug: hasQRCode logic, API responses

### Backend Issues
- File: `internal/service-domain/equipment-registry/`
- Check: Backend logs, API responses
- Debug: SQL queries, QR generation

### Database Issues
- Check: equipment_registry table
- Verify: qr_code and qr_code_url columns have data
- Debug: SQL queries in psql

---

## Status Checklist

Before reporting issues, verify:

- [ ] Backend is running (`.\backend.exe`)
- [ ] Database is running (`docker ps`)
- [ ] Frontend loaded (`http://localhost:3000/equipment`)
- [ ] Browser console open (F12)
- [ ] Network tab shows API calls
- [ ] Console logs show equipment count
- [ ] Equipment data loaded successfully
- [ ] QR codes exist in database

---

**For additional help, check console logs and share the output!** 🔍
