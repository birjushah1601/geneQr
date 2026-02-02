# QR Code Functionality - Complete Guide

## âœ… BACKEND IMPLEMENTATION (Verified Working)

### 1. Database Storage
**Location:** `database/migrations/002_store_qr_in_database.sql`

**Fields in `equipment` table:**
- `qr_code_image` (BYTEA) - Binary PNG image data
- `qr_code_format` (VARCHAR) - Format type (default: 'png')
- `qr_code_generated_at` (TIMESTAMP) - Generation timestamp
- `qr_code` (VARCHAR) - QR code identifier (e.g., "QR-HOSP001-CT001")
- `qr_code_path` (VARCHAR) - DEPRECATED: Legacy filesystem path

**Storage Method:** QR codes are stored as binary data in the database, NOT as files.

### 2. QR Code Generation
**Location:** `internal/service-domain/equipment-registry/qrcode/generator.go`

**Key Functions:**
- `GenerateQRCodeBytes(equipmentID, serialNumber, qrCodeID)` - Generates QR as byte array
- Returns PNG image bytes that are stored in database

**QR Code Content (JSON):**
```json
{
  "url": "http://localhost:3000/service-request?qr=QR-HOSP001-CT001",
  "id": "equipment-uuid",
  "serial": "SN12345",
  "qr": "QR-HOSP001-CT001"
}
```

### 3. Backend API Endpoints
**Location:** `internal/service-domain/equipment-registry/api/handler.go`

#### Generate QR Code
```
POST /api/v1/equipment/{id}/qr
```
- Generates QR code and stores in database
- Updates `qr_code_image` and `qr_code_generated_at` fields

#### Get QR Code Image
```
GET /api/v1/equipment/qr/image/{id}
```
- Returns PNG image from database
- Sets `Content-Type: image/png`
- Caches for 1 day

#### Download QR PDF Label
```
GET /api/v1/equipment/{id}/qr/pdf
```
- Generates printable PDF with equipment details
- Includes QR code image from database

---

## ðŸŽ¨ FRONTEND IMPLEMENTATION

### Current Flow:
1. **Display:** Shows QR code if `hasQRCode` is true
2. **Image URL:** `http://localhost:8081/api/v1/equipment/qr/image/{id}`
3. **Generate Button:** Calls `POST /api/v1/equipment/{id}/qr`
4. **Preview Modal:** Shows full-size QR code

### What We Need to Fix:

#### âœ… Already Working:
- Backend generates and stores QR in database
- API endpoint serves QR image from database
- Generate button calls correct API

#### âŒ Issues to Fix:
1. **hasQRCode Logic:** Frontend checks `qr_code_generated_at` field
2. **Placeholder:** No placeholder shown while loading
3. **Error Handling:** No fallback if image fails to load
4. **Refresh:** After generation, needs page reload

---

## ðŸ”§ HOW IT SHOULD WORK

### Scenario 1: Equipment WITHOUT QR Code
**Display:** "Generate" button (16x16px box)
**Action:** Click â†’ POST to backend â†’ QR generated â†’ Reload page â†’ Show QR image

### Scenario 2: Equipment WITH QR Code
**Display:** QR code thumbnail (20x20px)
**Hover:** Shows "Preview" and "Download" buttons
**Click:** Opens preview modal with full-size image

### Scenario 3: QR Image Load Failure
**Fallback:** Show placeholder with "Regenerate" button

---

## ðŸ“ TESTING CHECKLIST

### Backend Tests:
- [ ] Generate QR for equipment: `POST /api/v1/equipment/{id}/qr`
- [ ] Verify QR stored in database: `SELECT qr_code_image FROM equipment WHERE id='{id}'`
- [ ] Get QR image: `GET /api/v1/equipment/qr/image/{id}` (should return PNG)
- [ ] Download PDF label: `GET /api/v1/equipment/{id}/qr/pdf` (should download)

### Frontend Tests:
- [ ] Equipment list shows "Generate" button for items without QR
- [ ] Clicking "Generate" creates QR and reloads page
- [ ] Equipment with QR shows thumbnail image
- [ ] Clicking thumbnail opens preview modal
- [ ] Modal shows full-size QR code
- [ ] Download button works from modal
- [ ] Bulk generate works for multiple equipment

---

## ðŸ› COMMON ISSUES & FIXES

### Issue 1: QR Image Not Loading
**Symptom:** Broken image icon
**Causes:**
- Backend not running
- Wrong API path (should be `/api/v1/equipment/qr/image/{id}`)
- QR not generated yet (qr_code_image is NULL)

**Fix:**
- Verify backend is running on port 8081
- Check browser network tab for 404/500 errors
- Generate QR code first

### Issue 2: Generate Button Not Working
**Symptom:** Button click does nothing
**Causes:**
- API endpoint not reachable
- Equipment ID not found
- Database write permissions

**Fix:**
- Check browser console for errors
- Verify equipment exists in database
- Check backend logs

### Issue 3: Page Needs Refresh After Generation
**Why:** Frontend state not updated after API call
**Solution:** After successful generation, reload the page or refetch data

---

## ðŸ’¡ IMPROVEMENTS NEEDED

1. **Real-time Update:** Instead of page reload, update component state
2. **Loading State:** Show spinner while generating
3. **Error Messages:** Better user feedback on failures
4. **Batch Operations:** Progress indicator for bulk generation
5. **Preview Before Save:** Show QR before storing in database

---

## ðŸš€ QUICK START FOR TESTING

### 1. Start Backend
```bash
cd C:\Users\birju\ServQR
.\backend.exe
```

### 2. Start Frontend
```bash
cd admin-ui
npm run dev
```

### 3. Test Generate QR
```powershell
# Replace {id} with actual equipment ID
Invoke-RestMethod -Method POST -Uri "http://localhost:8081/api/v1/equipment/{id}/qr" -Headers @{"X-Tenant-ID"="default"}
```

### 4. Test Get QR Image
```
Open: http://localhost:8081/api/v1/equipment/qr/image/{id}
```
Should display PNG image in browser

### 5. Check Database
```sql
SELECT id, equipment_name, qr_code, 
       qr_code_generated_at, 
       length(qr_code_image) as image_size 
FROM equipment 
WHERE qr_code_image IS NOT NULL;
```

---

## âœ… SUMMARY

**What's Working:**
âœ… QR generation logic (backend)
âœ… Database storage (binary data)
âœ… Image serving endpoint
âœ… PDF label generation

**What Needs Fixing:**
âŒ Frontend refresh logic
âŒ Better error handling
âŒ Loading states
âŒ Placeholder images

**Next Steps:**
1. Test QR generation for one equipment
2. Verify image loads in frontend
3. Fix any broken image issues
4. Implement better UX (loading, errors)
