# QR Code Table Design Analysis

**Date:** December 23, 2025  
**Status:** 🧠 **Brainstorming - Design Decision Required**

---

## 🎯 Question

**Do we need a separate table for QR code generation, or can we use the existing `equipment_registry` table?**

---

## 📊 Current Schema Analysis

### **Existing Table: `equipment_registry`**

```sql
CREATE TABLE equipment_registry (
    id VARCHAR(32) PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,           -- ✅ Already has QR code
    serial_number VARCHAR(255) UNIQUE NOT NULL,
    
    -- Equipment details
    equipment_id VARCHAR(32),                        -- Link to catalog
    equipment_name VARCHAR(500) NOT NULL,            -- ❌ Required but might not be known yet
    manufacturer_name VARCHAR(255) NOT NULL,         -- ❌ Required
    model_number VARCHAR(255),
    category VARCHAR(255),
    
    -- Installation details
    customer_id VARCHAR(32),                         -- ✅ Can be NULL (unassigned)
    customer_name VARCHAR(500) NOT NULL,             -- ❌ Required but might not exist yet
    installation_location TEXT,                      -- ✅ Can be NULL
    installation_address JSONB,                      -- ✅ Can be NULL
    installation_date DATE,                          -- ✅ Can be NULL
    
    -- Contract details
    purchase_date DATE,
    warranty_expiry DATE,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'operational',
    
    -- QR Code URL
    qr_code_url TEXT NOT NULL,                       -- ✅ Already has this
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL
);
```

### **Issues with Current Schema for Pre-Generation:**

❌ **Problem 1:** `equipment_name` is **NOT NULL** but we might not know it yet  
❌ **Problem 2:** `customer_name` is **NOT NULL** but unassigned equipment has no customer  
❌ **Problem 3:** `manufacturer_name` is **NOT NULL** but might want generic QR codes  
❌ **Problem 4:** No batch tracking for bulk generation  
❌ **Problem 5:** No way to track "unused" vs "assigned" QR codes  
❌ **Problem 6:** No link to equipment_catalog_id for pre-generation  

---

## 💡 Design Options

### **Option 1: Modify Existing `equipment_registry` Table**

**Approach:** Make fields nullable and add status tracking

```sql
ALTER TABLE equipment_registry 
    ALTER COLUMN equipment_name DROP NOT NULL,
    ALTER COLUMN manufacturer_name DROP NOT NULL,
    ALTER COLUMN customer_name DROP NOT NULL,
    ADD COLUMN qr_status VARCHAR(50) DEFAULT 'unassigned',
    ADD COLUMN batch_id UUID NULL,
    ADD COLUMN reserved_for_customer_id UUID NULL;

-- Add check constraint
ALTER TABLE equipment_registry 
    ADD CONSTRAINT chk_equipment_qr_status 
    CHECK (qr_status IN ('generated', 'reserved', 'assigned', 'decommissioned'));
```

**Pros:**
- ✅ Simple - uses existing table
- ✅ No additional tables needed
- ✅ QR code always linked to equipment record
- ✅ Single source of truth

**Cons:**
- ❌ Breaks existing NOT NULL constraints
- ❌ Mixing two concepts: QR generation vs Equipment installation
- ❌ Can't track batch metadata well
- ❌ No audit trail for QR lifecycle
- ❌ Complex validation (some fields required based on status)
- ❌ Potential for data inconsistency

---

### **Option 2: Create Separate `qr_codes` Table** ⭐ **RECOMMENDED**

**Approach:** Dedicated table for QR code lifecycle management

```sql
-- New dedicated QR codes table
CREATE TABLE qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    qr_code VARCHAR(255) UNIQUE NOT NULL,           -- The actual QR code
    qr_code_url TEXT NOT NULL,                      -- Full URL to equipment page
    qr_image_url TEXT,                              -- Generated QR image URL
    
    -- Optional linking (NULL until assigned)
    equipment_catalog_id UUID REFERENCES equipment_catalog(id),  -- Model type
    manufacturer_id UUID REFERENCES organizations(id),           -- Manufacturer
    batch_id UUID REFERENCES qr_batches(id),                     -- Generation batch
    
    -- Assignment tracking
    equipment_registry_id UUID REFERENCES equipment_registry(id) NULL,  -- NULL = unassigned
    assigned_at TIMESTAMPTZ,
    assigned_by VARCHAR(255),
    
    -- Status lifecycle
    status VARCHAR(50) NOT NULL DEFAULT 'generated',
    -- Lifecycle: generated → reserved → assigned → decommissioned
    
    -- Physical tracking
    serial_number VARCHAR(255),                     -- Can be pre-assigned or added later
    printed BOOLEAN DEFAULT false,
    printed_at TIMESTAMPTZ,
    
    -- Metadata
    metadata JSONB,                                 -- Flexible data
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255)
);

-- Batch tracking table
CREATE TABLE qr_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_number VARCHAR(100) UNIQUE NOT NULL,
    
    -- Generation details
    manufacturer_id UUID REFERENCES organizations(id),
    equipment_catalog_id UUID REFERENCES equipment_catalog(id),  -- All QRs for this model
    
    -- Batch info
    quantity_requested INT NOT NULL,
    quantity_generated INT NOT NULL,
    start_serial_number VARCHAR(255),               -- First serial in batch
    end_serial_number VARCHAR(255),                 -- Last serial in batch
    
    -- Files
    pdf_url TEXT,                                   -- Downloadable PDF
    csv_url TEXT,                                   -- CSV export
    
    -- Status
    status VARCHAR(50) DEFAULT 'pending',           -- pending|generating|completed|failed
    
    -- Metadata
    generated_at TIMESTAMPTZ,
    generated_by VARCHAR(255),
    metadata JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_qr_codes_status ON qr_codes(status);
CREATE INDEX idx_qr_codes_manufacturer ON qr_codes(manufacturer_id);
CREATE INDEX idx_qr_codes_batch ON qr_codes(batch_id);
CREATE INDEX idx_qr_codes_equipment_registry ON qr_codes(equipment_registry_id);
CREATE INDEX idx_qr_codes_catalog ON qr_codes(equipment_catalog_id);
CREATE UNIQUE INDEX idx_qr_codes_qr_code ON qr_codes(qr_code);

CREATE INDEX idx_qr_batches_manufacturer ON qr_batches(manufacturer_id);
CREATE INDEX idx_qr_batches_status ON qr_batches(status);
```

**Pros:**
- ✅ Clean separation of concerns
- ✅ QR codes can exist independently
- ✅ Perfect for pre-generation workflow
- ✅ Batch tracking built-in
- ✅ Full lifecycle management (generated → reserved → assigned)
- ✅ Audit trail for QR code usage
- ✅ Can track printing status
- ✅ Flexible - doesn't break existing constraints
- ✅ Easy to query unassigned QR codes
- ✅ Supports QR code reservation
- ✅ Can link to equipment later via foreign key

**Cons:**
- ⚠️ Additional table(s) to maintain
- ⚠️ Join required to get full equipment details
- ⚠️ Need to sync data between tables on assignment

---

### **Option 3: Hybrid Approach**

**Approach:** Use both tables with clear responsibilities

```sql
-- qr_codes: Pre-generation and inventory
-- equipment_registry: Installed equipment only

-- Rule 1: QR generated → entry in qr_codes table (status=generated)
-- Rule 2: Equipment installed → entry in equipment_registry + update qr_codes (status=assigned)
-- Rule 3: qr_codes.equipment_registry_id links the two
```

**Workflow:**
```
1. Manufacturer generates QR codes
   → Creates records in qr_codes (status=generated)
   → Creates batch record in qr_batches
   → Downloads PDF

2. QR codes printed and applied to equipment
   → Update qr_codes.printed = true

3. Equipment sold/installed
   → Create record in equipment_registry
   → Update qr_codes.status = 'assigned'
   → Link: qr_codes.equipment_registry_id = equipment_registry.id

4. Customer scans QR
   → Lookup qr_code in qr_codes table
   → If assigned → fetch from equipment_registry
   → If unassigned → show "Equipment not yet registered"
```

**Pros:**
- ✅ Best of both worlds
- ✅ Clear separation: generation vs installation
- ✅ Equipment registry stays clean (only installed equipment)
- ✅ QR codes table handles lifecycle
- ✅ No breaking changes to existing schema

**Cons:**
- ⚠️ Two tables to maintain
- ⚠️ Synchronization logic required
- ⚠️ Developers need to understand relationship

---

## 🔍 Use Case Analysis

### **Use Case 1: Manufacturer Pre-Generates 1000 QR Codes**

**With Option 1 (Modify equipment_registry):**
```sql
-- Create 1000 equipment_registry records with minimal data
INSERT INTO equipment_registry (id, qr_code, qr_code_url, qr_status, batch_id, created_by)
VALUES 
  ('EQ-001', 'QR-20251223-000001', 'https://...', 'generated', 'BATCH-123', 'manufacturer-admin'),
  ('EQ-002', 'QR-20251223-000002', 'https://...', 'generated', 'BATCH-123', 'manufacturer-admin'),
  ...
  -- ❌ But equipment_name is NOT NULL - what to put?
  -- ❌ manufacturer_name is NOT NULL - what to put?
  -- ❌ customer_name is NOT NULL - what to put?
```

**With Option 2 (Separate qr_codes table):**
```sql
-- Create 1000 QR code records
INSERT INTO qr_codes (qr_code, qr_code_url, manufacturer_id, batch_id, status)
VALUES 
  ('QR-20251223-000001', 'https://...', 'uuid-mfr', 'uuid-batch', 'generated'),
  ('QR-20251223-000002', 'https://...', 'uuid-mfr', 'uuid-batch', 'generated'),
  ...
  -- ✅ Clean! Only required fields
```

### **Use Case 2: Assign QR to Equipment During Installation**

**With Option 1:**
```sql
-- Update existing record
UPDATE equipment_registry 
SET 
  equipment_name = 'Ventilator Pro 3000',
  manufacturer_name = 'MedTech Inc',
  customer_id = 'uuid-customer',
  customer_name = 'City Hospital',
  installation_location = 'ICU Ward 3',
  installation_date = '2025-12-23',
  qr_status = 'assigned',
  updated_at = NOW()
WHERE qr_code = 'QR-20251223-000001';
-- ✅ Simple update
```

**With Option 2:**
```sql
-- Create equipment record
INSERT INTO equipment_registry (id, qr_code, equipment_name, ...)
VALUES ('EQ-001', 'QR-20251223-000001', 'Ventilator Pro 3000', ...);

-- Update QR code record
UPDATE qr_codes 
SET 
  equipment_registry_id = 'EQ-001',
  status = 'assigned',
  assigned_at = NOW(),
  assigned_by = 'installer-123'
WHERE qr_code = 'QR-20251223-000001';
-- ⚠️ Two operations, but cleaner separation
```

### **Use Case 3: Query Unassigned QR Codes**

**With Option 1:**
```sql
SELECT qr_code, batch_id 
FROM equipment_registry 
WHERE qr_status = 'generated';
-- ✅ Simple query
```

**With Option 2:**
```sql
SELECT qr_code, batch_id, manufacturer_id
FROM qr_codes 
WHERE status = 'generated';
-- ✅ Even simpler - dedicated table
```

### **Use Case 4: Track QR Code Lifecycle**

**With Option 1:**
```sql
-- No audit trail
-- Can only see current status
SELECT qr_code, qr_status, created_at, updated_at
FROM equipment_registry;
-- ❌ Limited history
```

**With Option 2:**
```sql
-- Can add audit table or use temporal tables
CREATE TABLE qr_code_history (
  id UUID PRIMARY KEY,
  qr_code_id UUID REFERENCES qr_codes(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50),
  changed_at TIMESTAMPTZ,
  changed_by VARCHAR(255),
  notes TEXT
);
-- ✅ Full audit trail possible
```

---

## 📊 Comparison Matrix

| Feature | Option 1: Modify Existing | Option 2: New Table | Option 3: Hybrid |
|---------|---------------------------|---------------------|------------------|
| **Simplicity** | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **Pre-generation Support** | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Batch Tracking** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Lifecycle Management** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Audit Trail** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Data Integrity** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Backward Compatibility** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Query Performance** | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **Maintenance Complexity** | ⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **Scalability** | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Flexibility** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐ |

---

## 🎯 Recommendation

### **⭐ RECOMMENDED: Option 2 + Option 3 (New Tables with Hybrid Approach)**

**Why:**

1. **Clean Separation of Concerns**
   - `qr_codes` = QR code lifecycle management
   - `qr_batches` = Batch tracking and metadata
   - `equipment_registry` = Installed equipment only

2. **Supports All Workflows**
   - ✅ Pre-generation (bulk QR codes)
   - ✅ Reservation (hold for customer)
   - ✅ Assignment (link to installed equipment)
   - ✅ Batch tracking
   - ✅ Audit trail

3. **Future-Proof**
   - Can add more QR-related features
   - Can track QR code reuse
   - Can support QR code transfers
   - Can track physical printing

4. **No Breaking Changes**
   - Existing `equipment_registry` stays intact
   - New tables are additive only
   - Backward compatible

5. **Better Data Integrity**
   - No nullable fields in equipment_registry
   - Clear status transitions
   - Foreign key constraints

---

## 📋 Implementation Plan

### **Phase 1: Create New Tables**

```sql
-- 1. Create qr_batches table
CREATE TABLE qr_batches (...);

-- 2. Create qr_codes table
CREATE TABLE qr_codes (...);

-- 3. Add indexes
CREATE INDEX ...;

-- 4. Add foreign key to equipment_registry (optional)
ALTER TABLE equipment_registry 
ADD COLUMN qr_code_id UUID REFERENCES qr_codes(id);
```

### **Phase 2: Migration Strategy**

```sql
-- Migrate existing QR codes from equipment_registry to qr_codes
INSERT INTO qr_codes (
  qr_code, 
  qr_code_url, 
  equipment_registry_id,
  status,
  assigned_at,
  created_at
)
SELECT 
  qr_code,
  qr_code_url,
  id,
  'assigned', -- All existing are assigned
  installation_date,
  created_at
FROM equipment_registry;

-- Link back
UPDATE equipment_registry er
SET qr_code_id = qc.id
FROM qr_codes qc
WHERE qc.qr_code = er.qr_code;
```

### **Phase 3: Update Application Logic**

```go
// New workflow for QR generation
func GenerateBulkQRCodes(req BulkQRRequest) error {
  // 1. Create batch record
  batch := CreateBatch(req)
  
  // 2. Generate QR codes
  for i := 0; i < req.Quantity; i++ {
    qr := GenerateQRCode(batch.ID)
    qr.Status = "generated"
    db.Create(&qr)
  }
  
  // 3. Generate PDF
  pdf := GeneratePDF(batch.QRCodes)
  batch.PDFURL = pdf.URL
  batch.Status = "completed"
  db.Save(&batch)
  
  return nil
}

// Updated workflow for equipment installation
func RegisterEquipment(req EquipmentRequest) error {
  // 1. Create equipment record
  equipment := CreateEquipment(req)
  db.Create(&equipment)
  
  // 2. Update QR code
  qr := db.FindByQRCode(req.QRCode)
  qr.EquipmentRegistryID = equipment.ID
  qr.Status = "assigned"
  qr.AssignedAt = time.Now()
  db.Save(&qr)
  
  return nil
}
```

---

## ✅ Decision Summary

**CREATE NEW TABLES:**
1. ✅ `qr_batches` - For batch tracking
2. ✅ `qr_codes` - For QR lifecycle management
3. ✅ Keep `equipment_registry` as-is for installed equipment

**RELATIONSHIPS:**
- `qr_codes.batch_id` → `qr_batches.id`
- `qr_codes.equipment_registry_id` → `equipment_registry.id` (NULL = unassigned)
- `qr_codes.manufacturer_id` → `organizations.id`
- `qr_codes.equipment_catalog_id` → `equipment_catalog.id`

**BENEFITS:**
- ✅ Clean architecture
- ✅ Full lifecycle management
- ✅ No breaking changes
- ✅ Batch tracking built-in
- ✅ Audit trail support
- ✅ Future-proof design

---

## 🚀 Next Steps

1. **Approve design** ✅
2. Create migration scripts
3. Implement new tables
4. Build QR generation API
5. Update equipment registration API
6. Test workflows
7. Update documentation

---

**Status:** ✅ **RECOMMENDATION: Create separate `qr_codes` and `qr_batches` tables**  
**Confidence:** High  
**Breaking Changes:** None  
**Migration Required:** Yes (one-time, non-breaking)
