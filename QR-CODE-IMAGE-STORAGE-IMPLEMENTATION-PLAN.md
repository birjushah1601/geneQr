# QR Code Image Storage - Comprehensive Implementation Plan

**Date:** 2026-02-05 | **Author:** Droid AI | **Estimated Time:** 2.5-3 hours

---

## ?? Executive Summary

### Problem
- **Download QR PDF:** ? Broken (returns 404)
- **Equipment Detail QR Display:** ? Broken (image fails to load)
- **Root Cause:** `equipment_registry` table missing `qr_code_image` column

### Solution
Add 3 columns to `equipment_registry` + Fix repository code to store/load QR images

### What Works Already
- ? QR thumbnails in equipment list (uses external qrserver.com API)
- ? QR preview modal (uses external API)
- ? QR code scan ? ticket creation (uses qr_code text, not image)

---

## ?? Implementation Steps

### PHASE 1: Database Migration (30 min)

#### Step 1.1: Create Migration SQL
**File:** `add-qr-image-columns.sql`
```sql
BEGIN;

ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS qr_code_image BYTEA,
ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png',
ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_equipment_registry_has_qr 
ON equipment_registry(id) WHERE qr_code_image IS NOT NULL;

COMMENT ON COLUMN equipment_registry.qr_code_image IS 
    'Binary PNG for /api/v1/equipment/qr/image/{id}';

COMMIT;
```

#### Step 1.2: Migrate Existing QR Images  
**File:** `migrate-existing-qr-images.sql`
```sql
BEGIN;

-- Copy QR images from equipment table (10 records exist in both tables)
UPDATE equipment_registry er
SET qr_code_image = e.qr_code_image,
    qr_code_format = e.qr_code_format,
    qr_code_generated_at = e.qr_code_generated_at
FROM equipment e
WHERE er.id = e.id AND e.qr_code_image IS NOT NULL;

COMMIT;
```

#### Step 1.3: Apply Migrations
```bash
# Get container ID
docker ps | grep postgres

# Apply migrations
docker exec -i <container_id> psql -U postgres -d med_platform < add-qr-image-columns.sql
docker exec -i <container_id> psql -U postgres -d med_platform < migrate-existing-qr-images.sql
```

#### Step 1.4: Verify
```sql
-- Check columns added
SELECT column_name, data_type 
FROM information_schema.columns
WHERE table_name = 'equipment_registry'
AND column_name IN ('qr_code_image', 'qr_code_format', 'qr_code_generated_at');

-- Check migrated data
SELECT 
    COUNT(*) as total,
    COUNT(qr_code_image) FILTER (WHERE qr_code_image IS NOT NULL) as with_image
FROM equipment_registry;
-- Expected: total=73, with_image=10
```

---

### PHASE 2: Backend Code Fixes (45 min)

#### Fix 2.1: Update SELECT Query
**File:** `internal/service-domain/equipment-registry/infra/repository.go`  
**Location:** Line ~20-48

**Add these 3 lines to `equipmentSelectColumns`:**
```go
const equipmentSelectColumns = `
    id,
    COALESCE(qr_code,'') AS qr_code,
    // ... existing columns ...
    COALESCE(qr_code_url,'') AS qr_code_url,
    qr_code_image,                    // ? ADD
    qr_code_format,                   // ? ADD
    qr_code_generated_at,             // ? ADD
    COALESCE(notes,'') AS notes,
    // ... rest ...`
```

#### Fix 2.2: Update scanEquipment()
**File:** `internal/service-domain/equipment-registry/infra/repository.go`  
**Location:** ~Line 650-700

**Add 3 fields to row.Scan():**
```go
func (r *EquipmentRepository) scanEquipment(row pgx.Row) (*domain.Equipment, error) {
    err := row.Scan(
        &equipment.ID,
        &equipment.QRCode,
        // ... existing fields ...
        &equipment.QRCodeURL,
        &equipment.QRCodeImage,          // ? ADD
        &equipment.QRCodeFormat,         // ? ADD
        &equipment.QRCodeGeneratedAt,   // ? ADD
        &equipment.Notes,
        // ... rest ...
    )
    return &equipment, err
}
```

#### Fix 2.3: Implement UpdateQRCode()
**File:** `internal/service-domain/equipment-registry/infra/repository.go`  
**Location:** Line 706-719

**Replace NO-OP function with actual storage:**
```go
func (r *EquipmentRepository) UpdateQRCode(ctx context.Context, equipmentID string, qrImage []byte, format string) error {
    query := `
        UPDATE equipment_registry 
        SET qr_code_image = $2,
            qr_code_format = $3,
            qr_code_generated_at = NOW(),
            updated_at = NOW()
        WHERE id = $1`
    
    tag, err := r.pool.Exec(ctx, query, equipmentID, qrImage, format)
    if err != nil {
        return fmt.Errorf("failed to update QR code: %w", err)
    }
    if tag.RowsAffected() == 0 {
        return fmt.Errorf("equipment not found: %s", equipmentID)
    }
    return nil
}
```

#### Fix 2.4: Rebuild Backend
```bash
cd C:\Users\birju\aby-med
go build -o aby-med-platform.exe ./cmd/server

# Restart backend
Stop-Process -Name "aby-med-platform" -ErrorAction SilentlyContinue
Start-Process -FilePath ".\aby-med-platform.exe" -NoNewWindow
```

---

### PHASE 3: Backend Testing (30 min)

#### Test 3.1: Generate QR Code
```bash
# Generate QR for test equipment
curl -X POST http://localhost:8081/api/v1/equipment/347S6CxhID9V8CnhCZbnWUYdhUQ/qr

# Verify in database
docker exec <container_id> psql -U postgres -d med_platform -c "
SELECT id, LENGTH(qr_code_image) as size, qr_code_format, qr_code_generated_at 
FROM equipment_registry 
WHERE id = '347S6CxhID9V8CnhCZbnWUYdhUQ';"
-- Expected: size=1041, format=png, timestamp=recent
```

#### Test 3.2: Retrieve QR Image
```bash
# Get QR image
curl http://localhost:8081/api/v1/equipment/qr/image/347S6CxhID9V8CnhCZbnWUYdhUQ --output test-qr.png

# Verify it's a valid PNG
file test-qr.png
# Expected: PNG image data, 200 x 200

# View in browser
start test-qr.png
```

#### Test 3.3: Download QR PDF
```bash
# Download PDF
curl http://localhost:8081/api/v1/equipment/347S6CxhID9V8CnhCZbnWUYdhUQ/qr/pdf --output test-label.pdf

# Verify PDF
file test-label.pdf
# Expected: PDF document

# Open PDF - should contain QR code + equipment info
start test-label.pdf
```

#### Test 3.4: Bulk Generate
```bash
# Generate for all 73 equipment
curl -X POST http://localhost:8081/api/v1/equipment/qr/bulk-generate

# Expected response: "Generated 63 QR codes" (10 already exist)

# Verify all have images
docker exec <container_id> psql -U postgres -d med_platform -c "
SELECT COUNT(qr_code_image) as count FROM equipment_registry;"
# Expected: count=73
```

---

### PHASE 4: Frontend Verification (20 min)

#### Test 4.1: Equipment Detail Page
**URL:** http://localhost:3000/equipment/347S6CxhID9V8CnhCZbnWUYdhUQ

**Expected:**
- ? QR code image displays (no 404 error)
- ? "Download PDF" button works
- ? PDF contains QR code + equipment info

**Note:** Frontend code already correct - uses backend API endpoint

#### Test 4.2: Equipment List
**URL:** http://localhost:3000/equipment

**Expected:**
- ? QR thumbnails display in first column (uses external API - no changes)
- ? Click thumbnail opens preview modal (uses external API - no changes)
- ? "Download" button in hover menu works (now uses database)

---

### PHASE 5: Integration Testing (15 min)

#### Test 5.1: End-to-End New Equipment
```bash
# 1. Create equipment
curl -X POST http://localhost:8081/api/v1/equipment \
  -H "Content-Type: application/json" \
  -d '{"serial_number":"E2E-001","equipment_name":"E2E Test","qr_code":"QR-E2E-001","qr_code_url":"http://localhost:3000/service-request?qr=QR-E2E-001"}'

# 2. Generate QR
curl -X POST http://localhost:8081/api/v1/equipment/<new_id>/qr

# 3. View in frontend
# Open: http://localhost:3000/equipment/<new_id>
# Expected: QR displays, PDF downloads

# 4. Test scan flow
# Open: http://localhost:3000/service-request?qr=QR-E2E-001
# Expected: Equipment loads, ticket creation works
```

#### Test 5.2: Performance Check
```bash
# Time bulk generation for all 73 equipment
time curl -X POST http://localhost:8081/api/v1/equipment/qr/bulk-generate

# Expected: < 30 seconds (~100-200ms per QR)
```

---

## ?? Rollback Plan

### If Migration Fails
```sql
BEGIN;
ALTER TABLE equipment_registry 
DROP COLUMN IF EXISTS qr_code_image,
DROP COLUMN IF EXISTS qr_code_format,
DROP COLUMN IF EXISTS qr_code_generated_at;
DROP INDEX IF EXISTS idx_equipment_registry_has_qr;
COMMIT;
```

### If Code Changes Break
```bash
# Revert code
git checkout HEAD -- internal/service-domain/equipment-registry/infra/repository.go

# Rebuild
go build -o aby-med-platform.exe ./cmd/server
```

---

## ? Success Criteria

### Must Pass All Tests

**Database:**
- [ ] 3 columns added to equipment_registry
- [ ] 10 existing QR images migrated
- [ ] No errors or constraint violations

**Backend:**
- [ ] POST /equipment/{id}/qr stores QR in database
- [ ] GET /equipment/qr/image/{id} returns PNG (200 status)
- [ ] GET /equipment/{id}/qr/pdf returns PDF with QR
- [ ] POST /equipment/qr/bulk-generate works for all
- [ ] Repository loads qr_code_image correctly

**Frontend:**
- [ ] Equipment detail page displays QR (no 404)
- [ ] PDF download button works
- [ ] List thumbnails work (external API)
- [ ] Preview modal works (external API)
- [ ] No console errors

**Performance:**
- [ ] QR generation < 200ms per equipment
- [ ] Bulk generation < 30 seconds for 73 equipment

---

## ?? Estimated Timeline

| Phase | Task | Time |
|-------|------|------|
| 1 | Database migration + verification | 30 min |
| 2 | Backend code fixes + rebuild | 45 min |
| 3 | Backend testing | 30 min |
| 4 | Frontend verification | 20 min |
| 5 | Integration testing | 15 min |
| | **TOTAL** | **2.5 hrs** |

---

## ?? Pre-Implementation Checklist

Before proceeding:

- [ ] User has reviewed this plan
- [ ] Timeline is acceptable
- [ ] Database backup will be taken
- [ ] Rollback plan is understood
- [ ] Success criteria agreed upon

---

## ?? Next Steps

Once approved, I will execute:

1. ? Create migration files
2. ? Apply migrations to database
3. ? Fix repository code (3 locations)
4. ? Rebuild backend
5. ? Run comprehensive tests
6. ? Verify frontend
7. ? Provide summary report

**Ready to proceed?**
