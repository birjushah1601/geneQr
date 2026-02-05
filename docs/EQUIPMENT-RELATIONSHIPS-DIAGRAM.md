# Equipment Table Relationships

## Complete Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                     equipment (CATALOG)                      │
│  ─────────────────────────────────────────────────────────  │
│  • Generic equipment MODELS                                  │
│  • Example: "Siemens MAGNETOM Vida (MRI)"                   │
│  • Fields: model_number, manufacturer, category, specs      │
└────────────────────┬─────────────────────────────────────────┘
                     │
                     │ Referenced by:
                     │
      ┌──────────────┼──────────────┐
      │              │              │
      ▼              ▼              ▼
┌─────────┐  ┌──────────────┐  ┌──────────────┐
│ spare_  │  │  equipment_  │  │   Other      │
│ parts_  │  │  registry    │  │  catalog     │
│ catalog │  │              │  │  tables      │
└─────────┘  └──────┬───────┘  └──────────────┘
                    │
                    │ Referenced by:
                    │ (6 operational tables)
                    │
      ┌─────────────┼─────────────┐
      │             │             │
      ▼             ▼             ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│ mainten- │  │ equip_   │  │ equip_   │
│ ance_    │  │ downtime │  │ usage_   │
│ schedules│  │          │  │ logs     │
└──────────┘  └──────────┘  └──────────┘
      │             │             │
      ▼             ▼             ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│ equip_   │  │ equip_   │  │ service_ │
│ service_ │  │ document │  │ tickets  │
│ config   │  │ s        │  │          │
└──────────┘  └──────────┘  └──────────┘
      │
      ▼
┌──────────┐
│ equip_   │
│ attach-  │
│ ments    │
└──────────┘
```

---

## Critical Relationship

### equipment_registry → equipment (MUST EXIST)

**This is the KEY relationship that connects everything:**

```sql
CREATE TABLE equipment_registry (
    id UUID PRIMARY KEY,
    serial_number VARCHAR(255) UNIQUE,
    equipment_id UUID NOT NULL,  -- ← CRITICAL FK
    qr_code VARCHAR(255) UNIQUE,
    installation_date DATE,
    installation_location TEXT,
    warranty_expiry DATE,
    -- ... other fields
    
    -- CRITICAL FOREIGN KEY
    CONSTRAINT fk_equipment_registry_equipment 
    FOREIGN KEY (equipment_id) 
    REFERENCES equipment(id)
    ON DELETE RESTRICT
);
```

**Why this is critical:**
1. **Links installation to model** - Each installation must know its model type
2. **Enables part lookup** - Installation → Model → Compatible Parts
3. **Provides specifications** - Installation inherits model specs
4. **Supports reporting** - Aggregate data by model type

---

## Data Flow Examples

### Example 1: Finding Parts for an Installation

**Scenario:** Field engineer scans QR code on equipment, needs compatible parts

```sql
-- Step 1: Get installation details
SELECT * FROM equipment_registry 
WHERE qr_code = 'QR-12345';
-- Returns: serial_number='SN-12345', equipment_id='model-abc'

-- Step 2: Get model information
SELECT * FROM equipment 
WHERE id = 'model-abc';
-- Returns: model_number='MAGNETOM Vida', manufacturer='Siemens'

-- Step 3: Find compatible parts
SELECT * FROM spare_parts_catalog 
WHERE equipment_id = 'model-abc';
-- Returns: All parts compatible with this model

-- Combined query:
SELECT 
    er.serial_number,
    e.model_number,
    e.manufacturer,
    p.part_number,
    p.part_name,
    p.unit_price
FROM equipment_registry er
JOIN equipment e ON er.equipment_id = e.id
JOIN spare_parts_catalog p ON p.equipment_id = e.id
WHERE er.qr_code = 'QR-12345';
```

### Example 2: Service History by Model

**Scenario:** Manufacturer wants to see all service tickets for a specific model

```sql
-- Find all installations of a model and their tickets
SELECT 
    e.model_number,
    COUNT(DISTINCT er.id) as installation_count,
    COUNT(st.id) as total_tickets,
    AVG(EXTRACT(EPOCH FROM (st.resolved_at - st.created_at))/3600) as avg_resolution_hours
FROM equipment e
JOIN equipment_registry er ON er.equipment_id = e.id
LEFT JOIN service_tickets st ON st.equipment_registry_id = er.id
WHERE e.model_number = 'MAGNETOM Vida'
GROUP BY e.id, e.model_number;
```

### Example 3: Installation with Maintenance Schedule

**Scenario:** Create maintenance schedule for newly installed equipment

```sql
-- Step 1: Register installation
INSERT INTO equipment_registry (
    serial_number,
    equipment_id,  -- ← Links to equipment model
    installation_location,
    installation_date
) VALUES (
    'SN-NEW-001',
    (SELECT id FROM equipment WHERE model_number = 'MAGNETOM Vida'),
    'City Hospital, Room 204',
    CURRENT_DATE
);

-- Step 2: Create maintenance schedule based on model
INSERT INTO maintenance_schedules (
    equipment_id,  -- ← Links to equipment_registry
    maintenance_type,
    frequency_days,
    next_due_date
)
SELECT 
    er.id,  -- The new installation
    'Preventive Maintenance',
    90,  -- Every 90 days (from model specs)
    CURRENT_DATE + INTERVAL '90 days'
FROM equipment_registry er
WHERE er.serial_number = 'SN-NEW-001';
```

---

## Relationship Rules

### ✅ MUST HAVE: equipment_registry → equipment

**Without this relationship:**
- ❌ Can't determine which model an installation is
- ❌ Can't find compatible parts
- ❌ Can't apply model specifications
- ❌ Can't aggregate data by model
- ❌ System breaks completely

**Schema requirement:**
```sql
ALTER TABLE equipment_registry
ADD CONSTRAINT fk_equipment_registry_equipment 
FOREIGN KEY (equipment_id) 
REFERENCES equipment(id)
ON DELETE RESTRICT;
```

**Why RESTRICT:**
- Can't delete a model if installations exist
- Prevents orphaned installation records
- Maintains data integrity

---

## Complete FK Chain

```
equipment (models)
  ↑
  │ FK: equipment_id
  │
equipment_registry (installations)
  ↑
  │ FK: equipment_id (or equipment_registry_id)
  │
  ├─ maintenance_schedules
  ├─ equipment_downtime
  ├─ equipment_usage_logs
  ├─ equipment_service_config
  ├─ equipment_documents
  ├─ equipment_attachments
  └─ service_tickets

spare_parts_catalog
  ↑
  │ FK: equipment_id
  │
equipment (models)
```

---

## Verification Queries

### Check if FK exists
```sql
SELECT 
    conname as constraint_name,
    conrelid::regclass as table_name,
    confrelid::regclass as referenced_table,
    pg_get_constraintdef(oid) as definition
FROM pg_constraint 
WHERE conrelid = 'equipment_registry'::regclass 
  AND confrelid = 'equipment'::regclass
  AND contype = 'f';
```

**Expected result:**
- constraint_name: `fk_equipment_registry_equipment` (or similar)
- table_name: `equipment_registry`
- referenced_table: `equipment`

### Check for orphaned registrations
```sql
-- Find installations without valid equipment model
SELECT 
    er.serial_number,
    er.equipment_id,
    'Orphaned - no matching equipment model' as issue
FROM equipment_registry er
WHERE NOT EXISTS (
    SELECT 1 FROM equipment e 
    WHERE e.id = er.equipment_id
);
```

**Should return 0 rows** if FK constraint exists and is enforced.

### Verify data integrity
```sql
-- Count installations per model
SELECT 
    e.model_number,
    e.manufacturer,
    COUNT(er.id) as installation_count
FROM equipment e
LEFT JOIN equipment_registry er ON er.equipment_id = e.id
GROUP BY e.id, e.model_number, e.manufacturer
ORDER BY installation_count DESC;
```

---

## Migration to Add FK (if missing)

```sql
-- Step 1: Check for orphaned records
SELECT COUNT(*) FROM equipment_registry er
WHERE NOT EXISTS (
    SELECT 1 FROM equipment e WHERE e.id = er.equipment_id
);

-- Step 2: Fix orphaned records (if any)
-- Option A: Delete orphaned records
DELETE FROM equipment_registry
WHERE NOT EXISTS (
    SELECT 1 FROM equipment e WHERE e.id = equipment_id
);

-- Option B: Set to NULL and handle manually
UPDATE equipment_registry
SET equipment_id = NULL
WHERE NOT EXISTS (
    SELECT 1 FROM equipment e WHERE e.id = equipment_id
);

-- Step 3: Add FK constraint
ALTER TABLE equipment_registry
ADD CONSTRAINT fk_equipment_registry_equipment 
FOREIGN KEY (equipment_id) 
REFERENCES equipment(id)
ON DELETE RESTRICT;

-- Step 4: Make NOT NULL if appropriate
ALTER TABLE equipment_registry
ALTER COLUMN equipment_id SET NOT NULL;
```

---

## Common Queries

### Get full equipment details for an installation
```sql
SELECT 
    er.serial_number,
    er.qr_code,
    er.installation_location,
    er.installation_date,
    e.model_number,
    e.manufacturer,
    e.category,
    e.specifications
FROM equipment_registry er
JOIN equipment e ON er.equipment_id = e.id
WHERE er.serial_number = 'SN-12345';
```

### Get all installations at a location
```sql
SELECT 
    er.serial_number,
    e.model_number,
    e.category,
    er.installation_date,
    er.warranty_expiry
FROM equipment_registry er
JOIN equipment e ON er.equipment_id = e.id
WHERE er.installation_location LIKE '%City Hospital%'
ORDER BY e.category, er.installation_date DESC;
```

### Get installations needing service (by model)
```sql
SELECT 
    e.model_number,
    er.serial_number,
    er.installation_location,
    ms.next_due_date,
    ms.maintenance_type
FROM maintenance_schedules ms
JOIN equipment_registry er ON ms.equipment_id = er.id
JOIN equipment e ON er.equipment_id = e.id
WHERE ms.next_due_date <= CURRENT_DATE + INTERVAL '7 days'
  AND ms.status = 'pending'
ORDER BY ms.next_due_date;
```

---

## Summary

**The relationship `equipment_registry.equipment_id → equipment.id` is CRITICAL and MUST exist.**

Without it:
- Installation records have no model information
- Parts cannot be matched to installations
- System loses fundamental functionality

**This is different from operational FKs:**
- Operational tables → equipment_registry (track operations per installation)
- Parts catalog → equipment (parts fit models)
- **Equipment registry → equipment (installations have models)** ← CRITICAL LINK

---

**Status:** ⏳ **Needs Verification**  
**Action:** Run verification queries to confirm FK exists

---

**Related Documentation:**
- [Equipment Architecture Final](./EQUIPMENT-ARCHITECTURE-FINAL.md)
- [Equipment Table Architecture Fix Plan](../EQUIPMENT-TABLE-ARCHITECTURE-FIX-PLAN.md)
