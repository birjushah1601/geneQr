# Spare Parts with Images - Complete Implementation

## Overview

Added visual support to the spare parts catalog, allowing parts to be displayed with images in the "Add Parts" modal for better identification and user experience.

---

## What Was Added

### **1. Database Schema**

The `spare_parts_catalog` table already had image fields:

```sql
-- Existing columns
image_url         TEXT        -- Single primary image URL
photos            TEXT[]      -- Array of additional photo URLs
```

### **2. Sample Images Added**

Updated 44 spare parts with placeholder images from Unsplash:

```sql
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=800',
        'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800'
    ]
WHERE part_number = 'XR-TUBE-001';
```

**Image Sources:**
- Unsplash (free, high-quality medical/tech images)
- Images sized at 400px (thumbnails) and 800px (full size)
- All images are placeholder URLs (can be replaced with real product photos)

---

## Implementation Details

### **Backend API Changes**

**File:** `cmd/platform/main.go`

#### Updated Query

```go
query := `
    SELECT 
        id, part_number, part_name, category, subcategory, 
        description, unit_price, currency, is_available, 
        stock_status, requires_engineer, engineer_level_required,
        installation_time_minutes, lead_time_days, minimum_order_quantity,
        image_url, photos    -- Added these fields
    FROM spare_parts_catalog
    WHERE is_available = true AND is_obsolete = false`
```

#### Updated Scan Variables

```go
var (
    id, partNumber, partName, cat, description, currency, stockStatus string
    subcategory, engineerLevel, imageURL *string  // Added imageURL
    unitPrice float64
    isAvailable, requiresEngineer bool
    installTime, leadTime, minOrderQty *int
    photos []string  // Added photos array
)

rows.Scan(&id, &partNumber, &partName, &cat, &subcategory,
    &description, &unitPrice, &currency, &isAvailable, &stockStatus,
    &requiresEngineer, &engineerLevel, &installTime, &leadTime, &minOrderQty,
    &imageURL, &photos)  // Added image fields
```

#### Updated Response

```go
if imageURL != nil {
    part["image_url"] = *imageURL
}
if photos != nil && len(photos) > 0 {
    part["photos"] = photos
}
```

**API Response Example:**

```json
{
  "parts": [
    {
      "id": "uuid",
      "part_number": "XR-TUBE-001",
      "part_name": "X-Ray Tube Assembly",
      "category": "X-Ray",
      "description": "High-voltage X-ray tube assembly, 150kV capacity",
      "unit_price": 12500.00,
      "currency": "USD",
      "is_available": true,
      "stock_status": "in_stock",
      "requires_engineer": true,
      "engineer_level_required": "L3",
      "installation_time_minutes": 180,
      "lead_time_days": 7,
      "minimum_order_quantity": 1,
      "image_url": "https://images.unsplash.com/photo-1516549655169-df83a0774514?w=400",
      "photos": [
        "https://images.unsplash.com/photo-1516549655169-df83a0774514?w=800",
        "https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800"
      ]
    }
  ],
  "count": 50
}
```

---

### **Frontend Changes**

**File:** `admin-ui/src/components/PartsAssignmentModal.tsx`

#### Updated TypeScript Interface

```typescript
interface SparePart {
  id: string;
  part_number: string;
  part_name: string;
  category: string;
  subcategory?: string;
  description?: string;
  unit_price: number;
  currency: string;
  is_available: boolean;
  stock_status: string;
  requires_engineer: boolean;
  engineer_level_required?: string;
  installation_time_minutes?: number;
  lead_time_days?: number;
  minimum_order_quantity: number;
  image_url?: string;     // Added
  photos?: string[];      // Added
}
```

#### Updated Part Card Display

```tsx
<CardContent className="p-4">
  <div className="flex items-start gap-4">
    <Checkbox
      checked={isSelected}
      onCheckedChange={() => handleSelectPart(part)}
      onClick={(e) => e.stopPropagation()}
    />
    
    {/* Part Image - NEW */}
    {part.image_url && (
      <div className="w-20 h-20 flex-shrink-0 rounded-md overflow-hidden bg-gray-100">
        <img 
          src={part.image_url} 
          alt={part.part_name}
          className="w-full h-full object-cover"
          onError={(e) => {
            // Hide image if it fails to load
            (e.target as HTMLImageElement).style.display = 'none';
          }}
        />
      </div>
    )}
    
    <div className="flex-1 space-y-2">
      {/* Part details... */}
    </div>
  </div>
</CardContent>
```

**Features:**
- ✅ 80x80px thumbnail displayed next to part details
- ✅ Rounded corners with gray background
- ✅ Graceful fallback if image fails to load
- ✅ Maintains aspect ratio with `object-cover`

---

## Visual Layout

### Before (Text Only)
```
┌────────────────────────────────────────┐
│ [ ] X-Ray Tube Assembly                │
│     Part: XR-TUBE-001                  │
│     $12,500                            │
│     Description: High-voltage...       │
│     [Critical] [L3 Required]           │
└────────────────────────────────────────┘
```

### After (With Image)
```
┌────────────────────────────────────────┐
│ [ ] [IMAGE]  X-Ray Tube Assembly       │
│     80x80    Part: XR-TUBE-001         │
│              $12,500                   │
│              Description: High-voltage │
│              [Critical] [L3 Required]  │
└────────────────────────────────────────┘
```

---

## Parts with Images

### X-Ray Parts (5 parts)
- **XR-TUBE-001** - X-Ray Tube Assembly
- **XR-DET-001** - Flat Panel Detector
- **XR-COL-001** - Collimator Assembly
- **XR-FILT-001** - X-Ray Filter Set
- **XR-GRID-001** - Anti-Scatter Grid

### CT Scanner Parts (4 parts)
- **CT-TUBE-001** - CT X-Ray Tube
- **CT-DET-001** - CT Detector Module
- **CT-SLIP-001** - Slip Ring Assembly
- **CT-COL-001** - CT Collimator

### MRI Parts (5 parts)
- **MRI-COIL-HEAD** - MRI Head Coil
- **MRI-COIL-BODY** - MRI Body Coil
- **MRI-GRAD-001** - Gradient Coil
- **MRI-CRYO-001** - Cryogen System
- **MRI-RF-AMP** - RF Power Amplifier

### Ultrasound Parts (4 parts)
- **US-PROBE-C60** - Convex Probe
- **US-PROBE-L38** - Linear Probe
- **US-GEL-001** - Ultrasound Gel
- **US-BATT-001** - Battery Pack

### Ventilator Parts (7 parts)
- **VENT-VALVE-001** - Expiratory Valve
- **VENT-VALVE-002** - Inspiratory Valve
- **VENT-SENS-O2** - Oxygen Sensor
- **VENT-SENS-CO2** - CO2 Sensor
- **VENT-FILT-001** - HEPA Filter
- **VENT-TUBE-001** - Breathing Circuit
- **VENT-BATT-001** - Ventilator Battery

### Patient Monitor Parts (7 parts)
- **PM-ECG-CABLE** - ECG Cable 5-Lead
- **PM-SPO2-SENSOR** - SpO2 Sensor
- **PM-NIBP-CUFF** - NIBP Cuff
- **PM-TEMP-PROBE** - Temperature Probe
- **PM-IBP-CABLE** - IBP Cable
- **PM-BATT-001** - Monitor Battery
- **PM-DISPLAY-001** - LCD Display Module

### Dialysis Parts (7 parts)
- **DIAL-PUMP-001** - Blood Pump Head
- **DIAL-FILT-001** - Dialyzer Filter
- **DIAL-LINE-001** - Bloodline Set
- **DIAL-CONC-BIC** - Bicarbonate Concentrate
- **DIAL-CONC-ACID** - Acid Concentrate
- **DIAL-PRES-001** - Pressure Transducer
- **DIAL-VALVE-001** - Solenoid Valve

### Anesthesia Parts (5 parts)
- **ANES-VAPOR-ISO** - Isoflurane Vaporizer
- **ANES-VAPOR-SEV** - Sevoflurane Vaporizer
- **ANES-CO2-ABS** - CO2 Absorbent
- **ANES-O2-SENS** - Oxygen Analyzer
- **ANES-BELLOW** - Ventilator Bellows

---

## Testing

### Test 1: API Endpoint

```bash
curl http://localhost:8081/api/v1/catalog/parts
```

**Expected Response:**
```json
{
  "parts": [
    {
      "part_name": "X-Ray Tube Assembly",
      "image_url": "https://images.unsplash.com/photo-1516549655169-df83a0774514?w=400",
      "photos": ["..."]
    }
  ],
  "count": 50
}
```

### Test 2: Frontend Modal

1. Open: `http://localhost:3000/service-request?qr=QR-CAN-XR-005`
2. Click "Add Parts" button
3. Modal opens showing parts list
4. **Expected:** Each part card shows thumbnail image (80x80px)
5. **Expected:** Images load from Unsplash
6. **Expected:** Parts without images show text only

### Test 3: Image Fallback

1. Check browser console for any image loading errors
2. Images that fail to load should be hidden gracefully
3. Part information should still be fully readable

---

## Future Enhancements

### 1. Image Gallery Modal
```tsx
// Click on part image to open full gallery
<Dialog>
  <DialogContent>
    <Carousel>
      {part.photos.map(photo => (
        <img src={photo} alt="Part photo" />
      ))}
    </Carousel>
  </DialogContent>
</Dialog>
```

### 2. Video Support
```typescript
interface SparePart {
  // ... existing fields
  video_url?: string;        // Installation video
  video_thumbnail?: string;  // Video preview image
}
```

### 3. 3D Model Viewer
```tsx
{part.model_3d_url && (
  <ModelViewer 
    src={part.model_3d_url} 
    alt={part.part_name}
    ar
    auto-rotate
  />
)}
```

### 4. Real Product Photos

Replace placeholder Unsplash images with:
- Manufacturer product photos
- Internal photography
- CAD renderings
- Installation diagrams

**Storage Options:**
- AWS S3
- Cloudinary
- Local storage: `/data/parts-images/`

**Update Query:**
```sql
UPDATE spare_parts_catalog 
SET image_url = 'https://yourdomain.com/parts/XR-TUBE-001.jpg',
    photos = ARRAY[
        'https://yourdomain.com/parts/XR-TUBE-001-1.jpg',
        'https://yourdomain.com/parts/XR-TUBE-001-2.jpg'
    ]
WHERE part_number = 'XR-TUBE-001';
```

### 5. Image Optimization

```typescript
// Lazy loading
<img 
  src={part.image_url} 
  loading="lazy"
  srcSet={`
    ${part.image_url}?w=400 400w,
    ${part.image_url}?w=800 800w
  `}
  sizes="(max-width: 400px) 400px, 800px"
/>
```

---

## Database Summary

```sql
-- Check parts with images
SELECT 
    COUNT(*) as total_parts,
    COUNT(image_url) as with_images,
    COUNT(photos) as with_photo_arrays,
    ROUND(100.0 * COUNT(image_url) / COUNT(*), 1) as percentage_with_images
FROM spare_parts_catalog;
```

**Result:**
```
total_parts | with_images | with_photo_arrays | percentage_with_images
------------+-------------+-------------------+------------------------
         50 |          44 |                44 |                   88.0
```

---

## Summary

✅ **Database:** 44/50 parts have images (88%)  
✅ **Backend:** API returns image_url and photos arrays  
✅ **Frontend:** Modal displays 80x80px thumbnails  
✅ **Fallback:** Graceful handling of missing images  
✅ **Ready:** System ready for real product images  

**Next Steps:**
1. Test modal with images
2. Replace placeholder images with real photos
3. Add image gallery feature
4. Consider video/3D model support

---

**Visual support for spare parts is complete!** 🎉
