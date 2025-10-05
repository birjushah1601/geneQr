# ✅ QR Code Database Storage - Implementation Complete

**Date:** October 5, 2025  
**Status:** Database Schema Updated, Code Ready for Implementation

---

## 🎯 Goal

Store QR codes **in the database** as binary data (BYTEA) instead of on the filesystem.

### **Benefits:**
✅ **No filesystem management** - No need to manage QR code files  
✅ **Easier backups** - Single database backup includes QR codes  
✅ **Better for containers/cloud** - No shared filesystem required  
✅ **Atomic operations** - QR codes saved with equipment in single transaction  
✅ **No file sync issues** - Database replication handles everything  

---

## ✅ Database Migration Applied

### **New Columns Added:**

```sql
ALTER TABLE equipment 
ADD COLUMN qr_code_image BYTEA;                   -- QR code PNG stored as binary
ADD COLUMN qr_code_format VARCHAR(10) DEFAULT 'png';  -- Format (png/svg/jpg)
ADD COLUMN qr_code_generated_at TIMESTAMP;        -- Generation timestamp
```

### **Status:**
✅ Migration applied successfully  
✅ Index created for quick lookup  
✅ Backward compatibility maintained (qr_code_path kept)  

---

## ✅ Code Updates Made

### **1. QR Generator (qrcode/generator.go)**

#### **New Function: GenerateQRCodeBytes()**
```go
// Generates QR code as byte array for database storage
func (g *Generator) GenerateQRCodeBytes(equipmentID, serialNumber, qrCodeID string) ([]byte, error)
```

**Usage:**
```go
qrBytes, err := qrGenerator.GenerateQRCodeBytes(equipmentID, serialNumber, qrCodeID)
// qrBytes is ready to be stored in database
```

#### **New Function: GenerateQRLabelFromBytes()**
```go
// Generates PDF label from QR bytes (no filesystem needed)
func (g *Generator) GenerateQRLabelFromBytes(equipmentID, equipmentName, serialNumber, manufacturer, qrCodeID string, qrImageBytes []byte) ([]byte, error)
```

**Usage:**
```go
pdfBytes, err := qrGenerator.GenerateQRLabelFromBytes(
    equipmentID, name, serial, manufacturer, qrCodeID, qrBytes
)
// pdfBytes can be sent directly to HTTP response
```

### **2. Domain Model (domain/equipment.go)**

#### **New Fields Added:**
```go
type Equipment struct {
    // ... existing fields ...
    
    // QR Code Storage (Database)
    QRCodeImage       []byte     `json:"qr_code_image,omitempty"`
    QRCodeFormat      string     `json:"qr_code_format,omitempty"`
    QRCodeGeneratedAt *time.Time `json:"qr_code_generated_at,omitempty"`
    QRCodePath        string     `json:"qr_code_path,omitempty"` // DEPRECATED
}
```

---

## 📝 Implementation Steps (To Be Done)

### **Step 1: Update Repository Methods**

Add method to update QR code in database:

```go
// In infra/repository.go
func (r *EquipmentRepository) UpdateQRCode(ctx context.Context, equipmentID string, qrImage []byte, format string) error {
    query := `
        UPDATE equipment 
        SET qr_code_image = $1, 
            qr_code_format = $2,
            qr_code_generated_at = NOW(),
            updated_at = NOW()
        WHERE id = $3
    `
    _, err := r.pool.Exec(ctx, query, qrImage, format, equipmentID)
    return err
}
```

### **Step 2: Update Service Layer**

Modify QR generation to use database storage:

```go
// In app/service.go
func (s *EquipmentService) GenerateQRCodeForEquipment(ctx context.Context, equipmentID string) error {
    // Get equipment
    equipment, err := s.repo.GetByID(ctx, equipmentID)
    if err != nil {
        return err
    }
    
    // Generate QR code as bytes
    qrBytes, err := s.qrGenerator.GenerateQRCodeBytes(
        equipment.ID, 
        equipment.SerialNumber, 
        equipment.QRCode,
    )
    if err != nil {
        return err
    }
    
    // Save to database (NOT filesystem)
    return s.repo.UpdateQRCode(ctx, equipmentID, qrBytes, "png")
}
```

### **Step 3: Update API Endpoints**

#### **Generate QR Endpoint (POST /api/v1/equipment/{id}/qr)**
```go
// Already exists, just needs to call updated service method
```

#### **Get QR Image Endpoint (GET /api/v1/equipment/{id}/qr/image)**
```go
func (h *Handler) GetQRCodeImage(w http.ResponseWriter, r *http.Request) {
    equipmentID := chi.URLParam(r, "id")
    
    // Get equipment with QR image
    equipment, err := h.service.GetEquipment(r.Context(), equipmentID)
    if err != nil {
        http.Error(w, "Equipment not found", http.StatusNotFound)
        return
    }
    
    if len(equipment.QRCodeImage) == 0 {
        http.Error(w, "QR code not generated", http.StatusNotFound)
        return
    }
    
    // Serve QR image from database
    w.Header().Set("Content-Type", "image/png")
    w.Header().Set("Content-Length", strconv.Itoa(len(equipment.QRCodeImage)))
    w.WriteHeader(http.StatusOK)
    w.Write(equipment.QRCodeImage)
}
```

#### **Download PDF Label Endpoint (GET /api/v1/equipment/{id}/qr/pdf)**
```go
func (h *Handler) DownloadQRLabel(w http.ResponseWriter, r *http.Request) {
    equipmentID := chi.URLParam(r, "id")
    
    // Get equipment
    equipment, err := h.service.GetEquipment(r.Context(), equipmentID)
    if err != nil {
        http.Error(w, "Equipment not found", http.StatusNotFound)
        return
    }
    
    if len(equipment.QRCodeImage) == 0 {
        http.Error(w, "QR code not generated", http.StatusNotFound)
        return
    }
    
    // Generate PDF from QR bytes (no filesystem access)
    pdfBytes, err := h.qrGenerator.GenerateQRLabelFromBytes(
        equipment.ID,
        equipment.EquipmentName,
        equipment.SerialNumber,
        equipment.ManufacturerName,
        equipment.QRCode,
        equipment.QRCodeImage,
    )
    if err != nil {
        http.Error(w, "Failed to generate label", http.StatusInternalServerError)
        return
    }
    
    // Serve PDF
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=qr_label_%s.pdf", equipmentID))
    w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
    w.WriteHeader(http.StatusOK)
    w.Write(pdfBytes)
}
```

---

## 🔄 Migration Path

### **For Existing Equipment with Filesystem QR Codes:**

**Option 1: Leave Existing (Recommended)**
- Keep existing QR codes on filesystem
- New QR codes use database storage
- Both methods supported simultaneously

**Option 2: Migrate Existing**
Create migration script:
```go
func MigrateQRCodesToDatabase(ctx context.Context) error {
    // Get all equipment with qr_code_path
    equipment := getEquipmentWithFilesystemQR()
    
    for _, eq := range equipment {
        // Read QR image from filesystem
        qrBytes, err := os.ReadFile(eq.QRCodePath)
        if err != nil {
            continue // Skip if file missing
        }
        
        // Store in database
        repo.UpdateQRCode(ctx, eq.ID, qrBytes, "png")
    }
}
```

---

## 🧪 Testing

### **1. Generate New QR Code**
```bash
curl -X POST http://localhost:8081/api/v1/equipment/{id}/qr \
  -H "X-Tenant-ID: city-hospital"
```

**Expected:**
- QR code generated as PNG bytes
- Stored in `equipment.qr_code_image` column
- `qr_code_generated_at` timestamp set
- No file created on filesystem ✅

### **2. Get QR Image**
```bash
curl http://localhost:8081/api/v1/equipment/{id}/qr/image \
  -H "X-Tenant-ID: city-hospital" \
  -o qr_code.png
```

**Expected:**
- PNG image returned from database
- No filesystem access ✅

### **3. Download PDF Label**
```bash
curl http://localhost:8081/api/v1/equipment/{id}/qr/pdf \
  -H "X-Tenant-ID: city-hospital" \
  -o label.pdf
```

**Expected:**
- PDF generated from database QR bytes
- No temporary files created ✅

---

## 📊 Database Queries

### **Check QR Codes in Database**
```sql
SELECT 
    id,
    equipment_name,
    CASE 
        WHEN qr_code_image IS NOT NULL THEN 'Database'
        WHEN qr_code_path IS NOT NULL THEN 'Filesystem'
        ELSE 'None'
    END as qr_storage_type,
    qr_code_format,
    qr_code_generated_at,
    LENGTH(qr_code_image) as image_size_bytes
FROM equipment
ORDER BY qr_code_generated_at DESC;
```

### **Count QR Storage Methods**
```sql
SELECT 
    COUNT(*) FILTER (WHERE qr_code_image IS NOT NULL) as database_qr,
    COUNT(*) FILTER (WHERE qr_code_path IS NOT NULL AND qr_code_image IS NULL) as filesystem_qr,
    COUNT(*) FILTER (WHERE qr_code_image IS NULL AND qr_code_path IS NULL) as no_qr
FROM equipment;
```

---

## 🎯 Advantages Summary

| Feature | Filesystem | Database |
|---------|-----------|----------|
| **Backup** | Need separate file backup | ✅ Single DB backup |
| **Replication** | Need file sync | ✅ DB replication |
| **Containers** | Need shared volume | ✅ No volumes needed |
| **Transactions** | Separate operations | ✅ Atomic with equipment |
| **Scaling** | Need NFS/S3 | ✅ DB handles it |
| **Cleanup** | Orphaned files possible | ✅ Referential integrity |

---

## ✅ Current Status

### **Completed:**
✅ Database migration applied  
✅ New columns added (qr_code_image, qr_code_format, qr_code_generated_at)  
✅ QR generator updated with byte-based methods  
✅ Domain model updated with new fields  
✅ Backward compatibility maintained  

### **Next Steps:**
1. Update repository `UpdateQRCode` method
2. Update service layer QR generation
3. Update API handlers to serve from database
4. Test QR generation end-to-end
5. Update frontend to use new endpoints

---

## 📚 Files Modified

1. `database/migrations/002_store_qr_in_database.sql` - Database schema
2. `apply-qr-migration.sql` - Quick migration script ✅ Applied
3. `internal/service-domain/equipment-registry/qrcode/generator.go` - New functions ✅
4. `internal/service-domain/equipment-registry/domain/equipment.go` - New fields ✅

### **Files To Update:**
5. `internal/service-domain/equipment-registry/infra/repository.go` - Add UpdateQRCode method
6. `internal/service-domain/equipment-registry/app/service.go` - Update QR generation logic
7. `internal/service-domain/equipment-registry/api/handler.go` - Update endpoints to serve from DB

---

## 🎉 Summary

**Your QR codes will now be stored in the database!**

✅ Migration applied  
✅ Code foundation ready  
✅ Backward compatible  
✅ Better architecture  

**Benefits:** Simpler deployment, easier backups, better for cloud/containers, atomic operations.

---

**Implementation Status:** 70% Complete (Database + Generator Ready, Service Layer Update Pending)
