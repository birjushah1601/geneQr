# T2B.1: Equipment Catalog & Parts Management

**Status:** ✅ Complete  
**Started:** November 16, 2025  
**Completed:** November 16, 2025  
**Effort:** 3-4 days  
**Priority:** Critical

---

## 🎯 Objective

Create a comprehensive equipment catalog system that separates **equipment types** (master catalog) from **installed instances** (equipment registry), and supports **context-aware parts recommendations** (ICU vs General Ward accessories).

---

## 🔍 Problem Statement

### **Current State:**
- `equipment_registry` table contains only installed equipment instances
- No master list of equipment types/models
- No parts/accessories catalog
- No way to recommend context-specific parts (ICU vs General Ward)
- Cannot track parts specifications, pricing, compatibility

### **Impact:**
- Cannot recommend parts during service tickets
- No AI-powered parts suggestions possible
- Engineers don't know which accessories are needed for which context
- No pricing information for cost estimation
- Cannot track parts compatibility across equipment

---

## ✅ Solution

Created 4 new tables + helper functions + views:

### **1. equipment_catalog** - Equipment Master List
**Purpose:** Master catalog of equipment types and models (NOT installed instances)

**Key Fields:**
- `manufacturer_id` - Which manufacturer makes this
- `equipment_type` - Generic type (Ventilator, MRI, CT Scanner)
- `model_number` + `model_name` - Specific model
- `category` - Diagnostic, Life Support, Surgical, etc.
- `specifications` - JSONB technical specs
- `service_manual_url`, `user_manual_url` - Documentation
- `typical_lifespan_years`, `maintenance_interval_months` - Service info

**Example:**
```sql
-- SIEMENS Ventilator Model ABC-123
INSERT INTO equipment_catalog (manufacturer_id, equipment_type, model_number, model_name, category)
VALUES ('siemens-uuid', 'Ventilator', 'ABC-123', 'SIEMENS PrismaVent', 'Life Support');
```

---

### **2. equipment_parts** - Parts Catalog
**Purpose:** Catalog of parts, accessories, and consumables for equipment

**Key Fields:**
- `equipment_catalog_id` - Which equipment this part belongs to
- `part_number`, `part_name` - Part identification
- `part_category` - consumable, replaceable, optional, tool
- `part_type` - accessory, component, filter, cable, etc.
- `is_critical` - Critical for equipment operation
- `is_oem` - Original Equipment Manufacturer part
- `standard_price`, `currency`, `lead_time_days` - Pricing/availability
- `compatible_models` - Array of compatible model numbers
- `lifespan_hours`, `replacement_frequency_months` - Lifecycle

**Example:**
```sql
-- Ventilator tube (consumable)
INSERT INTO equipment_parts (
    equipment_catalog_id, part_number, part_name, 
    part_category, part_type, is_critical, standard_price, lead_time_days
)
VALUES (
    'ventilator-catalog-uuid', 'TUBE-001', 'High-Flow Ventilator Tube',
    'consumable', 'accessory', true, 500.00, 2
);
```

---

### **3. equipment_parts_context** - Context-Specific Parts ⭐
**Purpose:** Recommend different parts based on installation context (ICU vs General Ward)

**Key Fields:**
- `equipment_catalog_id` - Which equipment
- `part_id` - Which part
- `installation_context` - ICU, General Ward, OT, ER, etc.
- `is_required` - Must-have for this context
- `is_recommended` - Recommended but optional
- `recommended_quantity` - How many needed
- `priority` - Display order (1=highest)
- `reason` - Why this part for this context
- `typical_usage_frequency` - Daily, Weekly, Per Procedure

**This is the KEY table for context-aware AI recommendations!**

**Example:**
```sql
-- ICU ventilator needs high-flow tubes
INSERT INTO equipment_parts_context (
    equipment_catalog_id, part_id, installation_context,
    is_required, recommended_quantity, priority, reason
)
VALUES (
    'ventilator-uuid', 'high-flow-tube-uuid', 'ICU',
    true, 5, 1, 'Critical for ICU high-flow oxygen delivery'
);

-- General Ward ventilator needs standard tubes
INSERT INTO equipment_parts_context (
    equipment_catalog_id, part_id, installation_context,
    is_required, recommended_quantity, priority, reason
)
VALUES (
    'ventilator-uuid', 'standard-tube-uuid', 'General Ward',
    true, 3, 1, 'Standard tubes sufficient for general ward use'
);
```

---

### **4. equipment_compatibility** - Cross-Equipment Parts
**Purpose:** Track which parts from one equipment can work with another

**Key Fields:**
- `part_id` - The part
- `compatible_equipment_id` - Equipment it's compatible with
- `compatibility_type` - direct, with_adapter, replacement, upgrade
- `requires_adapter` - Needs adapter part
- `adapter_part_id` - Which adapter needed
- `tested`, `test_date` - Compatibility validation

**Example:**
```sql
-- GE ventilator tube works with SIEMENS with adapter
INSERT INTO equipment_compatibility (
    part_id, compatible_equipment_id, compatibility_type, requires_adapter, adapter_part_id
)
VALUES (
    'ge-tube-uuid', 'siemens-ventilator-uuid', 'with_adapter', true, 'adapter-xyz-uuid'
);
```

---

## 🔧 Helper Functions

### **1. get_equipment_parts(equipment_catalog_id, include_optional)**
Get all parts for an equipment type.

```sql
-- Get all parts for a ventilator
SELECT * FROM get_equipment_parts('ventilator-uuid', true);

-- Get only critical parts
SELECT * FROM get_equipment_parts('ventilator-uuid', false);
```

**Returns:** part_id, part_number, part_name, part_category, is_critical, standard_price, lead_time_days

---

### **2. get_context_specific_parts(equipment_catalog_id, installation_context)** ⭐
**THIS IS THE KEY FUNCTION FOR AI!**

Get parts recommended for specific context (ICU, General Ward, etc.)

```sql
-- Get ICU-specific parts for ventilator
SELECT * FROM get_context_specific_parts('ventilator-uuid', 'ICU');

-- Get General Ward parts
SELECT * FROM get_context_specific_parts('ventilator-uuid', 'General Ward');
```

**Returns:** part_id, part_number, part_name, is_required, is_recommended, recommended_quantity, priority, reason, standard_price

**Ordered by:** is_required DESC, priority ASC (required first, then by priority)

---

###  **3. find_compatible_parts(equipment_catalog_id)**
Find parts from other equipment that work with this equipment.

```sql
SELECT * FROM find_compatible_parts('siemens-ventilator-uuid');
```

**Returns:** Compatible parts with compatibility_type and adapter requirements

---

### **4. search_equipment_catalog(search_term)**
Full-text search across equipment catalog.

```sql
-- Find all ventilators
SELECT * FROM search_equipment_catalog('ventilator');

-- Find SIEMENS MRI
SELECT * FROM search_equipment_catalog('SIEMENS MRI');
```

**Returns:** Top 50 results ranked by relevance

---

## 📊 Views

### **1. equipment_catalog_with_manufacturer**
Equipment catalog with manufacturer details (name, type, country).

```sql
SELECT * FROM equipment_catalog_with_manufacturer WHERE equipment_type = 'Ventilator';
```

---

### **2. parts_with_equipment**
Parts catalog with equipment and manufacturer details.

```sql
SELECT * FROM parts_with_equipment WHERE equipment_type = 'Ventilator';
```

---

### **3. context_parts_summary**
Summary of parts by equipment and installation context.

```sql
-- See parts breakdown per context
SELECT * FROM context_parts_summary WHERE equipment_name = 'SIEMENS PrismaVent';
```

**Shows:** total_parts, required_parts, recommended_parts, estimated_total_cost per context

---

## 🎯 Use Cases

### **Use Case 1: AI Diagnosis Recommends Parts**

When AI analyzes a service ticket:

```sql
-- 1. Get equipment catalog ID from ticket's equipment_registry
SELECT equipment_catalog_id FROM equipment_registry WHERE id = :ticket_equipment_id;

-- 2. Get installation context
SELECT installation_context FROM equipment_registry WHERE id = :ticket_equipment_id;

-- 3. Get context-specific parts
SELECT * FROM get_context_specific_parts(:catalog_id, :context);

-- Result: AI presents engineer with context-aware parts list!
```

---

### **Use Case 2: Remote Engineer Selects Parts**

During Stage 1 (Remote Diagnosis):

```sql
-- Present ALL parts (general + context-specific)
SELECT 
    p.*,
    epc.is_required,
    epc.recommended_quantity,
    epc.reason
FROM equipment_parts p
LEFT JOIN equipment_parts_context epc ON (
    p.id = epc.part_id 
    AND epc.installation_context = :context
)
WHERE p.equipment_catalog_id = :catalog_id
ORDER BY epc.is_required DESC NULLS LAST, p.part_name;
```

---

### **Use Case 3: Parts Procurement Cost Estimation**

```sql
-- Estimate total cost for ICU ventilator setup
SELECT 
    SUM(ep.standard_price * epc.recommended_quantity) as total_cost
FROM equipment_parts_context epc
JOIN equipment_parts ep ON epc.part_id = ep.id
WHERE epc.equipment_catalog_id = :ventilator_id
  AND epc.installation_context = 'ICU'
  AND epc.is_required = true;
```

---

### **Use Case 4: Find Alternative Compatible Parts**

If OEM part not available, find alternatives:

```sql
SELECT * FROM find_compatible_parts(:equipment_catalog_id)
WHERE compatibility_type IN ('direct', 'with_adapter');
```

---

## 🔐 Data Integrity

### **Constraints:**
- ✅ Unique model per manufacturer (can't have duplicate model numbers)
- ✅ Unique part numbers per equipment
- ✅ One part-context combination per equipment
- ✅ Valid categories and types (CHECK constraints)
- ✅ Positive prices, quantities, weights
- ✅ Foreign keys to organizations, equipment_catalog

### **Indexes:**
- ✅ Fast manufacturer lookups
- ✅ Fast equipment type searches
- ✅ Full-text search on names/descriptions
- ✅ GIN indexes on JSONB and arrays
- ✅ Partial indexes on is_active=true
- ✅ Context-specific priority ordering

---

## 📈 Performance

### **Query Performance Targets:**
- Equipment catalog search: <50ms
- Get equipment parts: <100ms
- Get context-specific parts: <100ms
- Find compatible parts: <200ms

### **Optimization:**
- Partial indexes on is_active=true (only active records)
- GIN indexes for JSONB and array searches
- Full-text search indexes
- Proper foreign key indexes

---

## 🔄 Relationship to Other Tables

```
equipment_catalog (Master List)
    ↓
equipment_registry (Installed Instances)
    - references equipment_catalog_id
    - adds installation_context
    - tracks specific serial numbers, locations

equipment_catalog
    ↓
equipment_parts
    ↓
equipment_parts_context (Context mapping)
    ↓
Used by: AI Parts Recommender (T2C.4)
```

---

## ✅ Success Criteria

**All Met:**
- ✅ All 4 tables created successfully
- ✅ Constraints and indexes in place
- ✅ 4 helper functions working
- ✅ 3 views created
- ✅ Backward compatible (equipment_registry still works)
- ✅ Ready for AI integration (Phase 2C)
- ✅ Documentation complete

---

## 📝 Migration Files

**Migration:** `dev/postgres/migrations/016-create-equipment-catalog.sql`  
**Documentation:** This file

---

## 🚀 Next Steps

**After this ticket:**
- ✅ Ready for T2B.2 (Engineer Expertise)
- ✅ Ready for data seeding (T2B.7)
- ✅ Ready for AI Parts Recommender (T2C.4)

**Immediate needs:**
1. Populate equipment_catalog from existing equipment_registry
2. Seed sample parts data
3. Define context mappings for common equipment

---

## 💡 Key Insights

### **Why Separate Catalog from Registry?**
- **Catalog** = Equipment *types* (1 SIEMENS Ventilator model)
- **Registry** = Installed *instances* (50 ventilators installed across hospitals)

### **Why Context-Specific Parts?**
- ICU ventilator needs different accessories than General Ward
- OT equipment needs different setup than Outpatient
- AI can recommend correct parts based on WHERE equipment is installed

### **Why Compatibility Table?**
- Parts shortage? Use compatible alternatives
- Upgrade paths for older equipment
- Cross-manufacturer parts (with adapters)

---

**Status:** ✅ **COMPLETE**  
**Ready for:** T2B.2 - Engineer Expertise & Service Configuration

---

**Excellent progress!** This foundation enables AI-powered context-aware parts recommendations! 🎉
