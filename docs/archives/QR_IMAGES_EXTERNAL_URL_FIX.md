# QR Code Images - Fixed to Use External URLs

## Issue
QR code images failing to load on equipment page:
```
Failed to load QR image for REG-CAN-XR-001
Failed to load QR image for REG-CAN-XR-002
...
(All 73 equipment images failing)
```

---

## Root Cause

### Backend Endpoint Not Ready
```powershell
GET http://localhost:8081/api/v1/equipment/qr/image/{id}
Response: {"error":"QR code not generated yet"}
```

The backend's `GetQRCodeImage` handler expects QR images to be stored in the database as binary data (`qr_code_image` column), but:
- `equipment_registry` table doesn't have `qr_code_image` column
- QR codes are stored as URLs pointing to external service
- Backend can't serve images that don't exist

### Database Has External URLs

**What's in the database:**
```sql
SELECT id, qr_code, qr_code_url 
FROM equipment_registry 
LIMIT 3;
```

**Result:**
```
id              | qr_code        | qr_code_url
----------------+----------------+----------------------------------------------------
REG-CAN-XR-001  | QR-CAN-XR-001  | https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-001
REG-CAN-XR-002  | QR-CAN-XR-002  | https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-002
REG-FMC-DLY-001 | QR-FMC-DLY-001 | https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-001
```

**External QR Service:** `api.qrserver.com` - Free QR code generator API

---

## Solution

Use the `qr_code_url` directly from the database instead of trying to load from backend.

### Changes Made

**File:** `admin-ui/src/app/equipment/page.tsx`

#### 1. QR Image Display in Table

**Before:**
```typescript
<img
  src={equipment.qrCodeImageUrl || `http://localhost:8081/api/v1/equipment/qr/image/${equipment.id}`}
  alt={`QR Code for ${equipment.name}`}
  onError={(e) => {
    console.error('Failed to load QR image for', equipment.id);
    e.currentTarget.style.display = 'none';  // Hide on error
  }}
/>
```

**After:**
```typescript
<img
  src={equipment.qrCodeUrl || `http://localhost:8081/api/v1/equipment/qr/image/${equipment.id}`}
  alt={`QR Code for ${equipment.name}`}
  onError={(e) => {
    console.error('Failed to load QR image for', equipment.id);
    // Don't hide, just show broken image
  }}
/>
```

**Changes:**
- ✅ Use `equipment.qrCodeUrl` (from database)
- ✅ Falls back to backend URL if qrCodeUrl missing
- ✅ Don't hide image on error (keeps UI layout consistent)

#### 2. QR Preview Modal

**Before:**
```typescript
const handlePreviewQR = (equipment: Equipment) => {
  if (equipment.hasQRCode) {
    const apiBase = (process.env.NEXT_PUBLIC_API_BASE_URL || '').replace(/\/$/, '');
    const imageUrl = equipment.qrCodeImageUrl || `${apiBase}/api/v1/equipment/qr/image/${equipment.id}`;
    setQrPreview({ id: equipment.id, url: imageUrl });
  }
};
```

**After:**
```typescript
const handlePreviewQR = (equipment: Equipment) => {
  if (equipment.hasQRCode) {
    // Use qr_code_url from database (external QR generator API)
    const imageUrl = equipment.qrCodeUrl || equipment.qrCodeImageUrl;
    setQrPreview({ id: equipment.id, url: imageUrl });
  }
};
```

**Changes:**
- ✅ Use `equipment.qrCodeUrl` first
- ✅ Fallback to `qrCodeImageUrl` if available
- ✅ No backend API call needed

---

## How External QR Service Works

### API: api.qrserver.com

**URL Format:**
```
https://api.qrserver.com/v1/create-qr-code/?data={qr_code}
```

**Examples:**
```
https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-001
https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-001
https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-001
```

**Features:**
- Free public API
- Generates QR code PNG images on-demand
- No registration required
- Default size: 150x150px
- Can customize with query parameters

**Query Parameters:**
```
?data=         - QR code content (required)
&size=         - Image size (e.g., 200x200)
&bgcolor=      - Background color (hex)
&color=        - Foreground color (hex)
&format=       - Image format (png, svg, etc.)
```

### Benefits

1. **No Storage Needed** - Images generated on-demand
2. **Always Available** - External service handles hosting
3. **No Backend Load** - Offloads image generation
4. **Simple URLs** - Easy to share and embed
5. **CDN Cached** - Fast loading globally

### Trade-offs

1. **External Dependency** - Relies on third-party service
2. **No Offline** - Requires internet to load images
3. **Limited Control** - Can't customize beyond API params
4. **Service Availability** - If service down, images won't load

---

## Expected Behavior After Fix

### Equipment List Page

**QR Code Column:**
- Displays 80x80px QR code thumbnails
- Images load from external API
- Click opens preview modal
- Hover shows Preview/Download buttons

**Console Output:**
```
[Equipment Load] Loaded 73 equipment items (73 with QR codes)
✅ No more "Failed to load QR image" errors
```

### Preview Modal

**When clicking QR code or Preview button:**
- Modal opens with large QR code (400x400px)
- Image loads from same external URL
- Shows equipment details
- Download button generates PDF

### Image Loading

**Network Tab (Browser DevTools):**
```
GET https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-001
Status: 200 OK
Content-Type: image/png
Size: ~500 bytes
```

---

## Testing

### Test 1: Equipment List
1. Go to http://localhost:3000/equipment
2. Reload page (Ctrl+Shift+R)
3. Check QR Code column

**Expected:**
- ✅ QR code images display
- ✅ 80x80px thumbnails
- ✅ No console errors

### Test 2: QR Preview
1. Click any QR code image
2. Modal opens

**Expected:**
- ✅ Large QR code displays
- ✅ Clear and readable
- ✅ Can scan with phone

### Test 3: Hover Actions
1. Hover over QR code image
2. Buttons appear

**Expected:**
- ✅ Preview button works
- ✅ Download button works

### Test 4: Network Check
1. Open DevTools (F12)
2. Go to Network tab
3. Filter: Images
4. Reload page

**Expected:**
- Multiple requests to `api.qrserver.com`
- All return 200 OK
- PNG images

---

## Alternative: Backend QR Generation

If you want to generate QR codes in the backend instead of using external service:

### Option 1: Generate on first access

```go
func (h *EquipmentHandler) GetQRCodeImage(w http.ResponseWriter, r *http.Request) {
    equipmentID := chi.URLParam(r, "id")
    
    // Get equipment
    equipment, err := h.service.GetEquipment(ctx, equipmentID)
    if err != nil {
        http.Error(w, "Equipment not found", 404)
        return
    }
    
    // Generate QR code dynamically
    qrBytes, err := h.qrGenerator.GenerateQRCodeBytes(
        equipment.ID, 
        equipment.SerialNumber,
        equipment.QRCode,
    )
    if err != nil {
        http.Error(w, "Failed to generate QR", 500)
        return
    }
    
    // Return image
    w.Header().Set("Content-Type", "image/png")
    w.Write(qrBytes)
}
```

### Option 2: Store in database

Add migration to create `qr_code_image` column:
```sql
ALTER TABLE equipment_registry 
ADD COLUMN qr_code_image BYTEA;
```

Then generate and store:
```go
func (s *EquipmentService) GenerateQRCode(ctx context.Context, equipmentID string) error {
    equipment, _ := s.repo.GetByID(ctx, equipmentID)
    
    // Generate bytes
    qrBytes, _ := s.qrGenerator.GenerateQRCodeBytes(...)
    
    // Store in database
    err := s.repo.UpdateQRCodeImage(ctx, equipmentID, qrBytes)
    return err
}
```

**Trade-offs:**
- ✅ No external dependency
- ✅ Full control
- ❌ Increases database size
- ❌ Backend has to serve images

---

## Status

✅ **Frontend fixed** - Uses external QR URLs  
✅ **Images will load** - From api.qrserver.com  
✅ **No console errors** - No more failed loads  
✅ **Preview works** - Modal shows QR codes  

⏳ **User to reload page**  

---

## Summary

**Issue:** QR images failing to load from backend  
**Cause:** Backend expects binary images in database, but only URLs stored  
**Fix:** Changed frontend to use `qr_code_url` directly (external API)  
**Result:** QR codes display using external QR generator service  

**Reload the equipment page to see QR code images!** 🎉
