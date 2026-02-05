# Equipment Architecture - Final Design

## Overview
The equipment system uses two tables with distinct purposes that must NOT be confused.

---

## Table Purposes

### 1. `equipment` (Catalog Table)
**Purpose:** Generic catalog of equipment MODELS

**What it stores:**
- Equipment model information
- Generic specifications
- Manufacturer details
- Model numbers and categories

**Example Records:**
- "Siemens MAGNETOM Vida" (MRI Model)
- "GE Discovery CT750 HD" (CT Scanner Model)
- "Philips Azurion 7" (Cath Lab Model)

**Used by:**
- Parts catalog (which parts fit which models)
- Equipment browsing/selection
- Model specifications
- Compatibility matching

---

### 2. `equipment_registry` (Installation Table)
**Purpose:** Specific equipment INSTALLATIONS at customer sites

**What it stores:**
- Specific installed equipment
- Serial numbers
- Installation locations (hospitals)
- QR codes for field service
- Warranty information
- Installation dates

**Example Records:**
- "MRI #SN-12345 at City Hospital, Room 204"
- "CT Scanner #SN-67890 at County Medical, 3rd Floor"
- "Cath Lab #SN-54321 at Regional Heart Center"

**Used by:**
- Operational data (maintenance, downtime, usage)
- Service tickets and history
- Field service (QR codes)
- Warranty tracking
- Customer installations

---

## Foreign Key Architecture

### ✅ CORRECT: Tables Referencing `equipment_registry`

These tables track OPERATIONAL data for specific installations:

1. **`maintenance_schedules`** → `equipment_registry`
   - Schedules for specific installed equipment
   - Example: "MRI #12345 needs service on March 15"

2. **`equipment_downtime`** → `equipment_registry`
   - Downtime of specific installations
   - Example: "CT #67890 was down 3 hours on Feb 5"

3. **`equipment_usage_logs`** → `equipment_registry`
   - Usage tracking of specific equipment
   - Example: "MRI #12345 ran 24 scans today"

4. **`equipment_service_config`** → `equipment_registry`
   - Service SLAs for specific installations
   - Example: "MRI #12345 has 4-hour SLA"

5. **`equipment_documents`** → `equipment_registry`
   - Documents for specific installations
   - Example: "Installation cert for MRI #12345"

6. **`equipment_attachments`** → `equipment_registry`
   - Photos/files of specific installations
   - Example: "Photo of CT #67890 setup"

### ✅ CORRECT: Tables Referencing `equipment` (Catalog)

These tables relate to equipment MODELS, not installations:

1. **`spare_parts_catalog`** → `equipment`
   - Parts compatible with equipment MODELS
   - Example: "Part ABC-123 fits MAGNETOM Vida model"
   - **Why:** Parts are defined by model compatibility, not specific installations
   - **Logic:** A part fits ALL installations of that model type

---

## Common Mistakes to Avoid

### ❌ WRONG: spare_parts_catalog → equipment_registry
**Why wrong:**
- Parts aren't tied to specific installations
- A part fits a MODEL, not a serial number
- Would require duplicate part records for each installation

**Example of the problem:**
```sql
-- WRONG: Tying part to specific installation
INSERT INTO spare_parts_catalog (
    part_number, 
    part_name,
    equipment_id  -- pointing to equipment_registry
) VALUES (
    'ABC-123',
    'Cooling Fan',
    'installation-#12345'  -- Specific MRI installation
);

-- This means the part ONLY works for MRI #12345?
-- What about MRI #67890 of the same model?
-- Need to create the same part record again? ❌
```

**Correct approach:**
```sql
-- CORRECT: Tying part to model
INSERT INTO spare_parts_catalog (
    part_number,
    part_name,
    equipment_id  -- pointing to equipment (catalog)
) VALUES (
    'ABC-123',
    'Cooling Fan',
    'model-magnetom-vida'  -- The MODEL
);

-- This part fits ALL MAGNETOM Vida installations ✅
```

---

## Relationship Diagram

```
equipment (MODELS)
├── spare_parts_catalog (which parts fit which models)
└── equipment_registry (specific installations)
    ├── maintenance_schedules
    ├── equipment_downtime
    ├── equipment_usage_logs
    ├── equipment_service_config
    ├── equipment_documents
    ├── equipment_attachments
    └── service_tickets
```

---

## Real-World Analogy

### Car Parts Store (spare_parts_catalog → equipment)
**Question:** "Do you have brake pads for a 2023 Honda Civic?"
**Answer:** "Yes, part #BP-2023-HC fits the 2023 Honda Civic MODEL"

**Not:** "Which specific Honda Civic? VIN number please?"
- That would be ridiculous!
- The part fits the MODEL, not a specific car

### Service History (operational tables → equipment_registry)
**Question:** "When was the last oil change for my car?"
**Answer:** "Your car (VIN: 12345) had an oil change on Feb 1, 2026"

**Not:** "All 2023 Honda Civics?"
- That would make no sense!
- Maintenance is tracked per SPECIFIC vehicle

---

## Database Query Examples

### Finding Parts for Equipment Model
```sql
-- Find all parts for a specific equipment model
SELECT 
    p.part_number,
    p.part_name,
    e.model_number as equipment_model
FROM spare_parts_catalog p
JOIN equipment e ON p.equipment_id = e.id
WHERE e.model_number = 'MAGNETOM Vida';
-- Returns parts that fit this MODEL
```

### Finding Maintenance for Specific Installation
```sql
-- Find maintenance history for specific installed equipment
SELECT 
    m.scheduled_date,
    m.maintenance_type,
    er.serial_number,
    er.location
FROM maintenance_schedules m
JOIN equipment_registry er ON m.equipment_id = er.id
WHERE er.serial_number = 'SN-12345';
-- Returns maintenance for THIS specific installation
```

### Finding All Installations of a Model (with their service history)
```sql
-- Find all installations of a model and their operational data
SELECT 
    er.serial_number,
    er.installation_location,
    COUNT(m.id) as maintenance_count
FROM equipment_registry er
JOIN equipment e ON er.equipment_id = e.id
LEFT JOIN maintenance_schedules m ON er.id = m.equipment_id
WHERE e.model_number = 'MAGNETOM Vida'
GROUP BY er.id;
-- Returns all MRI installations of this model with their history
```

---

## Migration Summary

### Completed (Feb 5, 2026)
✅ 6 FK constraints updated to reference `equipment_registry`:
1. `maintenance_schedules`
2. `equipment_downtime`
3. `equipment_usage_logs`
4. `equipment_service_config`
5. `equipment_documents`
6. `equipment_attachments`

### No Changes Needed
✅ `spare_parts_catalog` correctly references `equipment` (catalog)

---

## Design Principles

### When to use `equipment` (catalog):
- ✅ Parts compatibility
- ✅ Model specifications
- ✅ Generic equipment info
- ✅ Equipment selection/browsing
- ✅ Model-level configuration

### When to use `equipment_registry` (installations):
- ✅ Operational data (maintenance, downtime, usage)
- ✅ Service tickets and history
- ✅ QR codes for field service
- ✅ Warranty tracking
- ✅ Customer-specific data
- ✅ Installation-specific documents

---

## Testing Checklist

### Verify FK Constraints
```sql
-- Check all FKs pointing to equipment_registry
SELECT 
    conname as constraint_name,
    conrelid::regclass as table_name,
    confrelid::regclass as referenced_table
FROM pg_constraint 
WHERE confrelid = 'equipment_registry'::regclass 
  AND contype = 'f'
ORDER BY conrelid::regclass::text;

-- Should return 6 operational tables
```

### Verify spare_parts_catalog
```sql
-- Confirm spare_parts_catalog points to equipment
SELECT 
    conname as constraint_name,
    confrelid::regclass as referenced_table
FROM pg_constraint 
WHERE conrelid = 'spare_parts_catalog'::regclass 
  AND contype = 'f';

-- Should reference 'equipment', not 'equipment_registry'
```

---

## Future Considerations

### Potential Confusion Points
1. **Field name ambiguity:** Both relationships use `equipment_id`
   - Consider renaming in `spare_parts_catalog` to `equipment_model_id`
   - Would make the intent clearer

2. **Join complexity:** Queries might need both tables
   ```sql
   -- Example: Find parts for a specific installation
   SELECT p.*
   FROM equipment_registry er
   JOIN equipment e ON er.equipment_id = e.id
   JOIN spare_parts_catalog p ON p.equipment_id = e.id
   WHERE er.serial_number = 'SN-12345';
   ```

3. **Table naming:** Consider renaming for clarity
   - `equipment` → `equipment_catalog` or `equipment_models`
   - `equipment_registry` → `equipment_installations`
   - Would eliminate confusion

---

## Related Documentation

- [Equipment Table Architecture Fix Plan](../EQUIPMENT-TABLE-ARCHITECTURE-FIX-PLAN.md)
- [Parts Architecture Explained](../PARTS_ARCHITECTURE_EXPLAINED.md)
- [Parts Catalog Template Guide](../PARTS-CATALOG-TEMPLATE-GUIDE.md)

---

## Conclusion

The current architecture is **CORRECT**:

- **Operational tables** reference `equipment_registry` (specific installations)
- **Parts catalog** references `equipment` (generic models)

This design follows database normalization principles and real-world logic where:
- Parts are compatible with MODELS
- Operations are tracked per INSTALLATION

**No further changes needed.**

---

**Last Updated:** 2026-02-05  
**Status:** ✅ Architecture Validated and Documented
