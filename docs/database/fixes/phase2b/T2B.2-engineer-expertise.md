# T2B.2: Engineer Expertise & Service Configuration

**Status:** ✅ Complete  
**Started:** November 16, 2025  
**Completed:** November 16, 2025  
**Effort:** 2-3 days  
**Priority:** Critical  
**Dependencies:** T2B.1

---

## 🎯 Objective

Track engineer expertise levels (L1/L2/L3) for each equipment type and configure service ownership (who handles what equipment - manufacturer vs client vs dealer). Enable AI to intelligently match engineers to tickets based on skills, certifications, and service configuration.

---

## 🔍 Problem Statement

### **Current State:**
- No tracking of engineer equipment expertise
- No L1/L2/L3 support level differentiation
- No way to know which engineer can service which equipment
- No manufacturer service configuration (who handles service)
- Cannot filter engineers by certification or experience
- AI cannot match engineers to tickets intelligently

### **Impact:**
- Manual engineer assignment (no automation possible)
- Wrong engineers assigned to complex issues
- No consideration of certifications
- Cannot enforce manufacturer service contracts
- No visibility into engineer capabilities

---

## ✅ Solution

Created 3 new tables + 5 helper functions + 3 views:

### **1. engineer_equipment_expertise** - Engineer Skills Database ⭐

**Purpose:** Track which equipment types each engineer can service and at what support level

**Support Levels:**
- **L1:** Basic/Remote support - Can diagnose remotely, handle simple issues
- **L2:** Advanced/Field support - Can handle complex issues, onsite visits
- **L3:** Expert support - Can handle critical/complex issues, installations, calibrations

**Key Fields:**
- `engineer_id` + `equipment_catalog_id` + `manufacturer_id` - What they can service
- `support_level` - L1, L2, or L3
- `certified` - Manufacturer certified?
- `years_experience` - How long working with this equipment
- `total_repairs_completed` + `successful_repairs` - Track record
- `first_time_fix_rate` - % of issues fixed on first attempt
- `customer_satisfaction_avg` - Average rating (0-5)
- `escalation_rate` - % of tickets escalated to higher level
- `can_do_remote`, `can_do_onsite`, `can_do_installation`, `can_do_calibration` - Capabilities
- `max_concurrent_tickets` - Workload limit

**Example:**
```sql
-- Ramesh can handle SIEMENS Ventilators at L2 level
INSERT INTO engineer_equipment_expertise (
    engineer_id, equipment_catalog_id, manufacturer_id,
    support_level, certified, years_experience,
    can_do_remote, can_do_onsite
)
VALUES (
    'ramesh-uuid', 'ventilator-catalog-uuid', 'siemens-uuid',
    'L2', true, 5,
    true, true
);
```

---

### **2. manufacturer_service_config** - Service Ownership ⭐

**Purpose:** Configure who handles service for each equipment (manufacturer vs client vs dealer)

**Key Concept - Hierarchy:**
```
Equipment-Specific Config (Priority 10)
    ↓ (if not found)
Manufacturer-Level Config (Priority 5)
    ↓ (if not found)
Default/Fallback
```

**Example Scenario:**
```
SIEMENS (Manufacturer Level):
  └── "We handle ALL our equipment" (service_provider_type = 'manufacturer')

SIEMENS Refrigerator (Equipment Override):
  └── "Client handles refrigerators" (service_provider_type = 'client')

Result: 
- SIEMENS MRI → SIEMENS handles (manufacturer config)
- SIEMENS Refrigerator → Client handles (equipment override)
```

**Key Fields:**
- `manufacturer_id` + `equipment_catalog_id` (NULL = all equipment)
- `service_provider_type` - manufacturer, client, dealer, distributor, third_party
- `service_provider_org_id` - Specific organization providing service
- `service_scope` - ['installation', 'repair', 'maintenance', 'calibration']
- `requires_oem_parts` - Must use OEM parts?
- `requires_certified_engineer` - Must be certified?
- `requires_manufacturer_approval` - Need approval before service?
- `warranty_void_if_third_party` - Using non-authorized voids warranty?
- `sla_response_hours`, `sla_resolution_hours` - Service level agreements
- `priority` - For hierarchy resolution

**Example:**
```sql
-- SIEMENS handles all their MRI service
INSERT INTO manufacturer_service_config (
    manufacturer_id, equipment_catalog_id,
    service_provider_type, service_provider_org_id,
    requires_oem_parts, requires_certified_engineer,
    sla_response_hours, sla_resolution_hours,
    priority
)
VALUES (
    'siemens-uuid', 'mri-catalog-uuid',
    'manufacturer', 'siemens-uuid',
    true, true,  -- Must use OEM parts and certified engineers
    4, 24,       -- 4 hours response, 24 hours resolution
    10           -- Equipment-specific (highest priority)
);
```

---

### **3. engineer_certifications** - Formal Certifications

**Purpose:** Track formal certifications from manufacturers and training organizations

**Key Fields:**
- `engineer_id` + `manufacturer_id` + `equipment_catalog_id`
- `certification_name`, `certification_number`
- `certification_level` - Basic, Advanced, Expert, Trainer
- `issued_by` - Organization that issued
- `issue_date`, `expiry_date` - Validity period
- `status` - active, expired, suspended, revoked
- `verified` - Verification status
- `renewable` - Needs periodic renewal?
- `skills_covered` - Array of specific skills

**Example:**
```sql
-- Ramesh has SIEMENS Ventilator L2 certification
INSERT INTO engineer_certifications (
    engineer_id, manufacturer_id, equipment_catalog_id,
    certification_name, certification_number, certification_level,
    issued_by, issue_date, expiry_date,
    status, skills_covered
)
VALUES (
    'ramesh-uuid', 'siemens-uuid', 'ventilator-uuid',
    'SIEMENS Ventilator Service Technician Level 2',
    'SIEM-VENT-L2-12345', 'Advanced',
    'SIEMENS Training Center', '2023-01-15', '2026-01-15',
    'active', ARRAY['Installation', 'Calibration', 'Advanced Repair']
);
```

---

## 🔧 Helper Functions

### **1. find_eligible_engineers(equipment_catalog_id, support_level, must_be_certified)** ⭐

**THE KEY FUNCTION FOR AI ASSIGNMENT!**

Find engineers qualified for equipment at specified support level.

```sql
-- Find L1 engineers for SIEMENS Ventilator
SELECT * FROM find_eligible_engineers('ventilator-uuid', 'L1', false);

-- Find only certified L2 engineers
SELECT * FROM find_eligible_engineers('ventilator-uuid', 'L2', true);
```

**Returns:** engineer_id, name, email, support_level, certified, years_experience, total_repairs, first_time_fix_rate, customer_satisfaction, can_do_remote, can_do_onsite

**Ordering:** certified DESC, first_time_fix_rate DESC, customer_satisfaction DESC, years_experience DESC

---

### **2. get_service_configuration(manufacturer_id, equipment_catalog_id)** ⭐

**Get service configuration with hierarchy resolution**

```sql
-- Who handles SIEMENS MRI service?
SELECT * FROM get_service_configuration('siemens-uuid', 'mri-uuid');

-- Who handles SIEMENS equipment in general?
SELECT * FROM get_service_configuration('siemens-uuid', NULL);
```

**Returns:** service_provider_type, service_provider_org_id, service_scope, requires_oem_parts, requires_certified_engineer, SLA details, config_level (equipment-specific or manufacturer-level)

**Logic:**
1. Try equipment-specific config first (priority 10)
2. If not found, try manufacturer-level (priority 5)
3. Returns which level was used (equipment-specific or manufacturer-level)

---

### **3. is_engineer_qualified(engineer_id, equipment_catalog_id, required_support_level)**

Quick boolean check if engineer is qualified.

```sql
-- Can Ramesh handle L2 support for this ventilator?
SELECT is_engineer_qualified('ramesh-uuid', 'ventilator-uuid', 'L2');
-- Returns: true/false
```

**Note:** L3 engineers can handle L1 and L2 tickets (higher levels cover lower)

---

### **4. get_engineer_expertise_summary(engineer_id)**

Get complete expertise profile for an engineer.

```sql
SELECT * FROM get_engineer_expertise_summary('ramesh-uuid');
```

**Returns:** All equipment types, support levels, certifications, experience, success rates

---

### **5. get_expiring_certifications(days_ahead)**

Find certifications expiring soon (for renewal reminders).

```sql
-- Find certifications expiring in next 30 days
SELECT * FROM get_expiring_certifications(30);

-- Find certifications expiring in next 90 days
SELECT * FROM get_expiring_certifications(90);
```

**Returns:** engineer_id, engineer_name, certification_name, expiry_date, days_until_expiry, equipment_type

---

## 📊 Views

### **1. engineer_expertise_with_details**

Complete engineer expertise with all related details.

```sql
SELECT * FROM engineer_expertise_with_details 
WHERE equipment_type = 'Ventilator' AND support_level = 'L2';
```

Shows: engineer details, equipment details, manufacturer, performance metrics, success rate

---

### **2. active_certifications**

All active certifications with expiry status.

```sql
SELECT * FROM active_certifications WHERE manufacturer_name = 'SIEMENS';
```

**Expiry Status:**
- "Never expires" - No expiry date
- "Valid" - >90 days remaining
- "Expiring Soon" - 30-90 days remaining
- "Expiring Critical" - <30 days remaining

---

### **3. service_configuration_summary**

Active service configurations with all details.

```sql
SELECT * FROM service_configuration_summary 
WHERE manufacturer_name = 'SIEMENS';
```

Shows: manufacturer, equipment scope, service provider, config scope (equipment vs manufacturer level)

---

## 🎯 Use Cases

### **Use Case 1: AI Assignment - Find Best Engineer**

```sql
-- Ticket created for SIEMENS Ventilator needing L1 support
-- 1. Get service configuration
SELECT * FROM get_service_configuration('siemens-uuid', 'ventilator-uuid');
-- Result: service_provider_type = 'client' (we handle it)

-- 2. Find eligible engineers
SELECT * FROM find_eligible_engineers('ventilator-uuid', 'L1', false);
-- Returns engineers ranked by:
--   - Certified first
--   - Best first-time fix rate
--   - Best customer satisfaction
--   - Most experience

-- 3. AI scores each engineer and assigns best match
```

---

### **Use Case 2: Enforce Manufacturer Service Contract**

```sql
-- Ticket for SIEMENS MRI
SELECT * FROM get_service_configuration('siemens-uuid', 'mri-uuid');
-- Result: service_provider_type = 'manufacturer'
--         requires_certified_engineer = true
--         requires_oem_parts = true

-- System enforces:
-- ✅ Only assigns to manufacturer's engineers
-- ✅ Only shows OEM parts in parts recommendation
-- ✅ Requires certified engineer
```

---

### **Use Case 3: Equipment-Specific Override**

```sql
-- SIEMENS wants to handle most equipment, but client handles refrigerators

-- Manufacturer-level (priority 5):
INSERT INTO manufacturer_service_config (manufacturer_id, service_provider_type, priority)
VALUES ('siemens-uuid', 'manufacturer', 5);

-- Equipment override (priority 10):
INSERT INTO manufacturer_service_config (
    manufacturer_id, equipment_catalog_id, service_provider_type, priority
)
VALUES ('siemens-uuid', 'refrigerator-uuid', 'client', 10);

-- Query for MRI:
SELECT * FROM get_service_configuration('siemens-uuid', 'mri-uuid');
-- Returns: 'manufacturer' (uses manufacturer-level config)

-- Query for Refrigerator:
SELECT * FROM get_service_configuration('siemens-uuid', 'refrigerator-uuid');
-- Returns: 'client' (uses equipment-specific override)
```

---

### **Use Case 4: Track Engineer Performance**

```sql
-- Update after ticket completion
UPDATE engineer_equipment_expertise
SET 
    total_repairs_completed = total_repairs_completed + 1,
    successful_repairs = successful_repairs + 1,  -- if successful
    first_time_fix_rate = (successful_repairs::NUMERIC / total_repairs_completed) * 100
WHERE engineer_id = 'ramesh-uuid' 
  AND equipment_catalog_id = 'ventilator-uuid';

-- AI uses these metrics for future assignments!
```

---

## 🔐 Data Integrity

### **Constraints:**
- ✅ Unique engineer-equipment-manufacturer combination
- ✅ Support level must be L1, L2, or L3
- ✅ Performance metrics within valid ranges (0-100%, 0-5 stars)
- ✅ Successful repairs ≤ total repairs
- ✅ Certification expiry > issue date
- ✅ Effective dates logical (to ≥ from)
- ✅ SLA hours positive
- ✅ Priority 1-10

### **Indexes:**
- ✅ Fast engineer lookups
- ✅ Fast equipment lookups
- ✅ Composite index for assignment queries (most common)
- ✅ GIN indexes for arrays (specializations, service_scope, skills)
- ✅ Partial indexes on is_active=true
- ✅ Expiry date indexes for certification tracking

---

## 📈 Performance

### **Query Performance Targets:**
- Find eligible engineers: <100ms
- Get service configuration: <50ms (with hierarchy)
- Check engineer qualified: <10ms (boolean check)
- Get expertise summary: <100ms

### **Optimization:**
- Composite index on (equipment_catalog_id, support_level, is_active, available_for_assignment)
- Partial indexes for active/available records only
- GIN indexes for array searches
- Proper foreign key indexes

---

## 🔄 Integration with Other Tables

```
engineer_equipment_expertise
    ↓
Used by: AI Assignment Optimizer (T2C.3)
    - Filters eligible engineers
    - Scores based on performance metrics
    - Considers certifications

manufacturer_service_config
    ↓
Used by: Workflow Orchestrator (T2D.1)
    - Determines which engineers to consider
    - Enforces OEM parts requirement
    - Applies SLA targets

engineer_certifications
    ↓
Used by: Engineer Assignment
    - Validates engineer qualifications
    - Tracks expiry for renewal reminders
```

---

## ✅ Success Criteria

**All Met:**
- ✅ All 3 tables created successfully
- ✅ Constraints and indexes in place
- ✅ 5 helper functions working
- ✅ 3 views created
- ✅ Service configuration hierarchy works correctly
- ✅ Ready for AI Assignment Optimizer (T2C.3)
- ✅ Documentation complete

---

## 📝 Migration Files

**Migration:** `dev/postgres/migrations/017-create-engineer-expertise.sql`  
**Documentation:** This file

---

## 🚀 Next Steps

**After this ticket:**
- ✅ Ready for T2B.3 (Configurable Workflows)
- ✅ Ready for T2C.3 (AI Assignment Optimizer) - **Can now intelligently match engineers!**
- ✅ Ready for data seeding (T2B.7)

**Immediate needs:**
1. Seed engineer expertise data for existing engineers
2. Create manufacturer service configurations
3. Import existing certifications

---

## 💡 Key Insights

### **Why L1/L2/L3 Support Levels?**
- **L1 (Remote):** Can diagnose remotely, handle simple issues - First line of support
- **L2 (Field):** Can handle complex issues, onsite visits - Most common
- **L3 (Expert):** Critical issues, installations, calibrations - Rare/complex

**Benefit:** Right engineer for right job. Don't send L3 expert for simple issue!

### **Why Service Configuration Hierarchy?**
- **Equipment-Specific:** SIEMENS Refrigerator → Client handles
- **Manufacturer-Level:** All other SIEMENS equipment → SIEMENS handles
- **Default:** Fallback if nothing configured

**Benefit:** Flexible contracts. One manufacturer can have different arrangements for different equipment.

### **Why Track Performance Metrics?**
- `first_time_fix_rate` - How often fixed on first attempt?
- `customer_satisfaction_avg` - Customer feedback
- `escalation_rate` - How often needs help?

**Benefit:** AI learns from history. Assigns best-performing engineers to critical tickets.

---

## 🎨 Real-World Example

**Scenario:** Hospital in Mumbai has SIEMENS MRI malfunctioning

**System Flow:**
```sql
-- 1. Check service configuration
SELECT * FROM get_service_configuration('siemens-uuid', 'mri-uuid');
-- Result: service_provider_type = 'manufacturer' (SIEMENS handles)
--         requires_certified_engineer = true

-- 2. Find eligible SIEMENS engineers (only from SIEMENS org)
SELECT * FROM find_eligible_engineers('mri-uuid', 'L2', true);
-- Returns certified L2 engineers only

-- 3. AI scores each engineer:
--    - Location (Mumbai proximity)
--    - Performance (first-time fix rate, customer satisfaction)
--    - Availability (current workload)
--    - Experience (years, total repairs)

-- 4. Assigns best match!
```

---

**Status:** ✅ **COMPLETE**  
**Ready for:** T2B.3 - Configurable Workflows

---

**Excellent progress!** This enables intelligent AI-powered engineer assignment! 🎉
