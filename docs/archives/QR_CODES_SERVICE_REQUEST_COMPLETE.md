# QR Codes - Service Request Integration Complete

## Overview

QR codes now encode URLs that directly open the service request page with equipment pre-selected, allowing field technicians to scan equipment and immediately create service tickets.

---

## How It Works

### Workflow

```
1. Field Technician scans QR code on equipment
   ↓
2. Phone camera/QR app reads the QR code
   ↓
3. QR code contains: http://localhost:3000/service-request?qr=QR-XR-ALPHA-001
   ↓
4. Browser opens service-request page with ?qr= parameter
   ↓
5. Page automatically fetches equipment details using QR code
   ↓
6. Equipment information pre-filled (name, serial, manufacturer, etc.)
   ↓
7. Technician fills in issue description and priority
   ↓
8. Ticket created and assigned
```

---

## Database Changes

### QR Code URL Update

**Before:**
```sql
SELECT id, qr_code, qr_code_url 
FROM equipment_registry 
LIMIT 3;
```

```
id               | qr_code         | qr_code_url
-----------------+-----------------+------------------------------------------------
REG-XR-ALPHA-001 | QR-XR-ALPHA-001 | https://app.example.com/qr/REG-XR-ALPHA-001
REG-FMC-DLY-001  | QR-FMC-DLY-001  | https://app.example.com/qr/REG-FMC-DLY-001
```

**After:**
```
id               | qr_code         | qr_code_url
-----------------+-----------------+----------------------------------------------------------
REG-XR-ALPHA-001 | QR-XR-ALPHA-001 | http://localhost:3000/service-request?qr=QR-XR-ALPHA-001
REG-FMC-DLY-001  | QR-FMC-DLY-001  | http://localhost:3000/service-request?qr=QR-FMC-DLY-001
```

**Update Query:**
```sql
UPDATE equipment_registry
SET qr_code_url = 'http://localhost:3000/service-request?qr=' || qr_code,
    updated_at = NOW()
WHERE qr_code IS NOT NULL AND qr_code != '';
```

**Result:** 73 equipment items updated

---

## Frontend Changes

### Equipment Page QR Code Display

**File:** `admin-ui/src/app/equipment/page.tsx`

#### QR Code Image in Table

**Before:**
```typescript
<img
  src={equipment.qrCodeUrl}  // Just the URL string
  alt={`QR Code for ${equipment.name}`}
/>
```

**After:**
```typescript
<img
  src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(equipment.qrCodeUrl || '')}`}
  alt={`QR Code for ${equipment.name}`}
/>
```

**What Changed:**
- QR code image now **encodes the full service-request URL**
- Uses external QR generator API to create image
- URL in database: `http://localhost:3000/service-request?qr=QR-XR-ALPHA-001`
- Gets encoded into QR code image
- Scanning reveals the URL and opens it

#### QR Code Preview Modal

**Before:**
```typescript
const handlePreviewQR = (equipment: Equipment) => {
  const imageUrl = equipment.qrCodeUrl;  // Just show URL
  setQrPreview({ id: equipment.id, url: imageUrl });
};
```

**After:**
```typescript
const handlePreviewQR = (equipment: Equipment) => {
  // Generate QR code image that encodes the service-request URL
  const qrImageUrl = `https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=${encodeURIComponent(equipment.qrCodeUrl || '')}`;
  setQrPreview({ id: equipment.id, url: qrImageUrl });
};
```

**What Changed:**
- Preview shows larger QR code (400x400px)
- Still encodes the service-request URL
- Can be scanned from screen or downloaded

---

## QR Code Generation

### External API: api.qrserver.com

**URL Format:**
```
https://api.qrserver.com/v1/create-qr-code/?size={size}&data={url}
```

**Parameters:**
- `size` - Image dimensions (e.g., 200x200, 400x400)
- `data` - The URL to encode (must be URL-encoded)

**Example:**
```
https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=http%3A%2F%2Flocalhost%3A3000%2Fservice-request%3Fqr%3DQR-XR-ALPHA-001
```

**What Gets Encoded:**
```
http://localhost:3000/service-request?qr=QR-XR-ALPHA-001
```

**When Scanned:**
- Phone opens this URL in browser
- Service request page loads
- Equipment details fetched automatically

---

## Service Request Page Integration

### How Page Detects QR Code

**File:** `admin-ui/src/app/service-request/page.tsx`

```typescript
function ServiceRequestPageInner() {
  const searchParams = useSearchParams();
  const qrCode = searchParams?.get('qr');  // Get ?qr= parameter
  
  useEffect(() => {
    if (!qrCode) {
      setError('No QR code provided...');
      return;
    }

    // Fetch equipment by QR code
    const fetchEquipment = async () => {
      try {
        const response = await equipmentApi.getByQRCode(qrCode);
        setEquipment(response);
      } catch (err) {
        setError('Equipment not found');
      }
    };

    fetchEquipment();
  }, [qrCode]);
  
  // Form is pre-filled with equipment details
}
```

### API Call

**Backend Endpoint:**
```
GET /api/v1/equipment/qr/{qr_code}
```

**Example:**
```
GET /api/v1/equipment/qr/QR-XR-ALPHA-001
```

**Response:**
```json
{
  "id": "REG-XR-ALPHA-001",
  "equipment_name": "X-Ray System Alpha",
  "serial_number": "SN-XR-ALPHA-001",
  "manufacturer_name": "Canon Medical Systems",
  "category": "X-Ray",
  "qr_code": "QR-XR-ALPHA-001",
  "installation_location": "X-Ray Room 1",
  "customer_name": "Apollo Hospital Mumbai"
}
```

### Form Pre-Fill

When equipment is loaded:
1. ✅ Equipment name displayed
2. ✅ Serial number shown
3. ✅ Manufacturer info shown
4. ✅ Location displayed
5. ✅ Equipment ID stored for ticket creation

User fills in:
- Issue description
- Priority level
- Their name

Click "Submit" → Ticket created with equipment linked

---

## Testing

### Test 1: View QR Codes on Equipment Page

**Steps:**
1. Go to http://localhost:3000/equipment
2. Reload page (Ctrl+Shift+R)
3. Look at QR Code column

**Expected:**
- QR code images displayed (200x200px)
- Each QR encodes service-request URL
- Hover shows Preview/Download buttons

### Test 2: Preview QR Code

**Steps:**
1. Click any QR code image
2. Modal opens with larger QR

**Expected:**
- 400x400px QR code displayed
- Same URL encoded
- Can be scanned from screen

### Test 3: Download QR Label

**Steps:**
1. Click "Download" button on equipment
2. PDF downloads

**Expected:**
- PDF contains equipment details
- QR code embedded in PDF
- Can print and affix to equipment

### Test 4: Scan QR Code with Phone

**Steps:**
1. Use phone camera or QR scanner app
2. Point at QR code on screen
3. Phone shows URL preview
4. Tap to open

**Expected:**
- URL: `http://localhost:3000/service-request?qr=QR-XR-ALPHA-001`
- Browser opens service-request page
- Equipment info automatically loaded
- Form ready to fill

### Test 5: Create Ticket from Scan

**Steps:**
1. Scan QR code with phone
2. Service request page opens
3. Equipment info displayed
4. Fill in description: "Screen flickering"
5. Set priority: High
6. Enter name: "John Doe"
7. Click Submit

**Expected:**
- Ticket created successfully
- Ticket linked to equipment
- Can view in tickets list
- Shows equipment details

---

## Production Deployment

### Update URLs for Production

**Current (Development):**
```
http://localhost:3000/service-request?qr=QR-XR-ALPHA-001
```

**Production:**
```
https://yourapp.com/service-request?qr=QR-XR-ALPHA-001
```

**Update Query:**
```sql
UPDATE equipment_registry
SET qr_code_url = 'https://yourapp.com/service-request?qr=' || qr_code,
    updated_at = NOW()
WHERE qr_code IS NOT NULL;
```

**Or use environment variable:**
```typescript
const APP_BASE_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000';
const qrCodeUrl = `${APP_BASE_URL}/service-request?qr=${qrCode}`;
```

---

## QR Code Printing

### Print QR Labels

**Option 1: PDF Download**
- Click "Download" on equipment
- PDF includes QR code + equipment details
- Print on label paper
- Affix to equipment

**Option 2: Bulk Print**
- Click "Export CSV" on equipment page
- Import to label printing software
- Generate QR codes in bulk
- Print sheet of labels

**Option 3: API Integration**
```
GET /api/v1/equipment/{id}/qr/pdf
```

Returns PDF with:
- Equipment name
- Serial number
- Manufacturer
- QR code (high resolution)
- Installation date

---

## QR Code Best Practices

### Size Guidelines

**Minimum Size:**
- At least 1" x 1" (2.5cm x 2.5cm)
- Ensures scannability from 6-12 inches away

**Recommended Sizes:**
- Small equipment: 1.5" x 1.5"
- Medium equipment: 2" x 2"
- Large equipment: 3" x 3"

### Placement

**Best Locations:**
- Front panel (easily accessible)
- Near manufacturer label
- Away from vents/heat sources
- Protected from wear

**Avoid:**
- Curved surfaces
- Areas with scratches
- Behind access panels
- In direct sunlight (fading)

### Material

**Label Material:**
- Laminated labels (water-resistant)
- Polyester labels (durable)
- Avoid paper (tears easily)

**Adhesive:**
- Permanent adhesive
- Temperature-resistant
- Medical-grade (if applicable)

---

## Benefits

### For Field Technicians

1. **Fast Ticket Creation**
   - Scan QR → Equipment auto-selected
   - No manual searching/typing
   - Reduces errors

2. **Equipment Verification**
   - Correct equipment every time
   - No confusion about asset IDs
   - Serial numbers verified

3. **Mobile-Friendly**
   - Works on any smartphone
   - No special app required
   - Offline QR scanning (online for form)

### For Management

1. **Accurate Data**
   - Equipment always correctly linked
   - Audit trail of scans
   - Location verification

2. **Faster Response**
   - Tickets created immediately
   - Assignment happens faster
   - Better SLA compliance

3. **Analytics**
   - Track scans per equipment
   - Identify problem equipment
   - Service frequency patterns

---

## Troubleshooting

### QR Code Won't Scan

**Possible Causes:**
1. Too small - print larger
2. Damaged/scratched - reprint
3. Poor lighting - use flashlight
4. Camera focus - move closer/farther
5. Wrong scanner app - use native camera

**Solution:**
- Reprint QR at larger size
- Clean surface before applying
- Test scan before affixing

### Wrong Equipment Loaded

**Possible Causes:**
1. QR code mismatch
2. Database outdated
3. Incorrect QR code ID

**Solution:**
- Verify QR code matches equipment
- Check database: `SELECT * FROM equipment_registry WHERE qr_code = 'QR-...'`
- Regenerate QR if needed

### Service Request Page Not Loading

**Possible Causes:**
1. Network issue
2. Backend down
3. Incorrect URL in QR

**Solution:**
- Check internet connection
- Verify backend running
- Check QR URL: should be `http://localhost:3000/service-request?qr=...`

---

## Future Enhancements

### Potential Features

1. **Offline Support**
   - Cache equipment data
   - Create tickets offline
   - Sync when online

2. **Location Tracking**
   - GPS coordinates on scan
   - Verify equipment location
   - Track movement

3. **Scan Analytics**
   - Count scans per equipment
   - Track technician activity
   - Identify frequently serviced equipment

4. **Smart Routing**
   - Auto-assign based on location
   - Nearest available engineer
   - Skill matching

5. **Photo Capture**
   - Take photo on scan
   - Attach to ticket automatically
   - Visual documentation

---

## Status

✅ **Database updated** - All 73 QR URLs point to service-request  
✅ **Frontend updated** - QR images encode service-request URLs  
✅ **Preview working** - Modal shows scannable QR codes  
✅ **Service integration** - Page detects ?qr= and loads equipment  

⏳ **Ready to test** - Scan and create tickets!  

---

## Quick Reference

### QR Code URL Pattern
```
http://localhost:3000/service-request?qr={QR_CODE_ID}
```

### QR Image Generator
```
https://api.qrserver.com/v1/create-qr-code/?size=200x200&data={URL}
```

### Service Request API
```
GET /api/v1/equipment/qr/{qr_code}
POST /api/v1/tickets (with equipment_id)
```

---

**QR codes are now fully integrated with service request workflow!** 🎉
