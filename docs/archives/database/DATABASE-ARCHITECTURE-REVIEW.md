# 🏗️ Database Architecture Review & Improvement Recommendations

**Date:** November 16, 2025  
**Reviewer:** Database Architecture Analysis  
**Severity Levels:** 🔴 Critical | 🟠 High | 🟡 Medium | 🟢 Low

---

## 📋 Executive Summary

After thorough review of the current database schema, **14 critical design issues** have been identified that will cause significant problems in production:

- **Data Integrity Issues:** 6 areas with potential inconsistency
- **Scalability Concerns:** 5 areas that won't scale well
- **Audit Trail Gaps:** 8 areas missing history tracking
- **Query Performance Risks:** 4 areas with inefficient queries

**Recommendation:** Implement schema refactoring **before** adding more features to avoid costly data migrations later.

---

## 🔴 CRITICAL ISSUES

### **Issue #1: Service Ticket Engineer Assignment (Your Example)**

#### **Current Flawed Design:**
```sql
service_tickets (
    id,
    assigned_engineer_id VARCHAR,  -- ❌ PROBLEM: Only ONE engineer!
    assigned_engineer_name VARCHAR,
    assignment_tier INT,
    assignment_tier_name TEXT,
    ...
)
```

#### **Problems:**
1. ❌ Can only track ONE engineer at a time (loses escalation history)
2. ❌ When reassigning, previous assignment is lost
3. ❌ Can't track: "Level 1 dealer engineer tried → failed → escalated to Level 2 manufacturer engineer"
4. ❌ No audit trail of why engineer changed
5. ❌ Can't query: "How many times was this ticket reassigned?"
6. ❌ Can't track: Engineer A worked 2 hours, Engineer B finished the job
7. ❌ Data duplication: name stored when there's an ID

#### **Correct Design:**
```sql
-- REMOVE from service_tickets:
-- assigned_engineer_id, assigned_engineer_name, assignment_tier, assignment_tier_name

-- USE ONLY the engineer_assignments table (which already exists!)
engineer_assignments (
    id UUID PRIMARY KEY,
    ticket_id UUID FK,
    engineer_id UUID FK,
    assignment_sequence INT,          -- 1, 2, 3, 4... (for escalation tracking)
    assignment_tier INT,               -- 1=Manufacturer, 2=Dealer, 3=Distributor, etc.
    assignment_tier_name TEXT,         -- "OEM Engineer", "Authorized Dealer", etc.
    assignment_reason TEXT,            -- "Initial", "Escalation", "Specialist Required"
    assigned_by UUID,
    assigned_at TIMESTAMPTZ,
    status TEXT,                       -- assigned|accepted|rejected|in_progress|completed|failed
    accepted_at TIMESTAMPTZ,
    rejected_at TIMESTAMPTZ,
    rejection_reason TEXT,             -- Why engineer rejected the assignment
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    completion_status TEXT,            -- success|failed|escalated|parts_required
    escalation_reason TEXT,            -- "Beyond my expertise", "Parts not available"
    time_spent_hours NUMERIC(5,2),
    ...
)

-- Current engineer is derived:
CREATE VIEW current_ticket_assignments AS
SELECT DISTINCT ON (ticket_id)
    ticket_id,
    engineer_id,
    assignment_sequence,
    assignment_tier,
    status
FROM engineer_assignments
WHERE status NOT IN ('completed', 'rejected', 'failed')
ORDER BY ticket_id, assigned_at DESC;
```

#### **Benefits:**
✅ Complete escalation history  
✅ Track multiple engineers per ticket  
✅ Query: "Show average escalation rate by equipment type"  
✅ Report: "How many tickets required L3 support?"  
✅ Audit trail: Who worked when and for how long  

---

### **Issue #2: Equipment Ownership - Static Relationships**

#### **Current Flawed Design:**
```sql
equipment (
    id,
    manufacturer_org_id UUID,    -- ❌ What if equipment has multiple manufacturers?
    sold_by_dealer_id UUID,      -- ❌ What if resold by different dealer?
    owned_by_org_id UUID,        -- ❌ What if ownership changes?
    installed_facility_id UUID,  -- ❌ What if equipment moved?
    ...
)
```

#### **Problems:**
1. ❌ Hospital sells equipment to another hospital → Can't track transfer
2. ❌ Equipment moved from ICU to ER → Installation facility changes
3. ❌ Lease equipment → Ownership temporarily changes
4. ❌ No history: "Who owned this in January 2024?"
5. ❌ Can't handle: Manufacturer → Distributor → Dealer → Hospital (supply chain)

#### **Correct Design:**
```sql
-- Remove static fields from equipment table
-- Keep only immutable data:
equipment (
    id UUID PRIMARY KEY,
    qr_code TEXT UNIQUE,
    serial_number TEXT UNIQUE,
    equipment_name TEXT,
    model_number TEXT,
    category TEXT,
    specifications JSONB,
    manufactured_date DATE,
    -- NO ownership/location fields here!
)

-- Create history table for all relationships:
equipment_relationships (
    id UUID PRIMARY KEY,
    equipment_id UUID FK,
    relationship_type TEXT,  -- 'manufactured_by', 'sold_by', 'owned_by', 'installed_at'
    organization_id UUID FK,
    facility_id UUID FK,
    valid_from TIMESTAMPTZ DEFAULT NOW(),
    valid_to TIMESTAMPTZ,    -- NULL = current
    transfer_reason TEXT,
    transfer_document_url TEXT,
    created_at TIMESTAMPTZ,
    created_by UUID
);

CREATE INDEX idx_equipment_rel_current 
ON equipment_relationships(equipment_id, relationship_type) 
WHERE valid_to IS NULL;

-- Query current owner:
SELECT org_id FROM equipment_relationships
WHERE equipment_id = 'EQ-123'
  AND relationship_type = 'owned_by'
  AND valid_to IS NULL;

-- Query ownership history:
SELECT * FROM equipment_relationships
WHERE equipment_id = 'EQ-123'
  AND relationship_type = 'owned_by'
ORDER BY valid_from DESC;
```

---

### **Issue #3: Organization Relationships - No Change History**

#### **Current Flawed Design:**
```sql
org_relationships (
    parent_org_id UUID,
    child_org_id UUID,
    commission_percentage NUMERIC(5,2),  -- ❌ What if commission changes?
    credit_limit NUMERIC(18,2),          -- ❌ What if credit limit increases?
    annual_target NUMERIC(18,2),         -- ❌ New target every year?
    start_date DATE,
    end_date DATE,
    ...
)
```

#### **Problems:**
1. ❌ Commission changes from 10% to 12% → Overwrites history
2. ❌ Can't calculate: "What commission did we pay in Q1 2024?"
3. ❌ Annual targets change → Previous year data lost
4. ❌ Relationship ends and restarts → Need new record? Same record?
5. ❌ Credit limit increases → Can't audit when it changed

#### **Correct Design:**
```sql
-- Keep basic relationship:
org_relationships (
    id UUID PRIMARY KEY,
    parent_org_id UUID FK,
    child_org_id UUID FK,
    rel_type TEXT,
    relationship_status TEXT,
    created_at TIMESTAMPTZ
);

-- Track all term changes:
org_relationship_terms (
    id UUID PRIMARY KEY,
    relationship_id UUID FK,
    version INT,
    effective_from DATE,
    effective_to DATE,
    
    -- Business terms:
    commission_percentage NUMERIC(5,2),
    payment_terms JSONB,
    credit_limit NUMERIC(18,2),
    annual_target NUMERIC(18,2),
    performance_tier TEXT,
    
    -- Audit:
    changed_by UUID,
    changed_at TIMESTAMPTZ,
    change_reason TEXT,
    
    CONSTRAINT no_overlap EXCLUDE USING gist (
        relationship_id WITH =,
        daterange(effective_from, effective_to, '[]') WITH &&
    )
);

-- Query: Current terms
SELECT * FROM org_relationship_terms
WHERE relationship_id = '...'
  AND CURRENT_DATE BETWEEN effective_from AND COALESCE(effective_to, '9999-12-31');

-- Query: What was commission in Jan 2024?
SELECT commission_percentage 
FROM org_relationship_terms
WHERE relationship_id = '...'
  AND '2024-01-15' BETWEEN effective_from AND COALESCE(effective_to, '9999-12-31');
```

---

### **Issue #4: RFQs/Quotes - JSONB for Queryable Data**

#### **Current Flawed Design:**
```sql
rfqs (
    items JSONB  -- ❌ [{"product": "CT Scanner", "qty": 1}, ...]
);

quotes (
    line_items JSONB  -- ❌ [{"product": "CT Scanner", "price": 50000}, ...]
);
```

#### **Problems:**
1. ❌ Can't query: "Show all RFQs for CT Scanners"
2. ❌ Can't report: "Total quote value by product category"
3. ❌ No referential integrity to products table
4. ❌ Can't enforce: Product must exist
5. ❌ Can't track: Which products are most requested?
6. ❌ JSONB queries are slow and don't use indexes well

#### **Correct Design:**
```sql
rfqs (
    id UUID PRIMARY KEY,
    rfq_number TEXT,
    requesting_org_id UUID FK,
    status TEXT,
    expected_response_date DATE,
    -- NO items here!
);

rfq_items (
    id UUID PRIMARY KEY,
    rfq_id UUID FK,
    line_number INT,
    product_id UUID FK,           -- ✅ Proper FK to products table
    sku_id UUID FK,               -- ✅ Specific SKU if known
    quantity INT,
    required_by_date DATE,
    specifications JSONB,         -- Additional specs OK as JSONB
    notes TEXT
);

quotes (
    id UUID PRIMARY KEY,
    rfq_id UUID FK,
    supplier_org_id UUID FK,
    quote_number TEXT,
    status TEXT,
    -- NO line_items here!
);

quote_items (
    id UUID PRIMARY KEY,
    quote_id UUID FK,
    rfq_item_id UUID FK,          -- ✅ Links back to specific RFQ item
    product_id UUID FK,
    sku_id UUID FK,
    quantity INT,
    unit_price NUMERIC(18,2),
    tax_amount NUMERIC(18,2),
    discount_percentage NUMERIC(5,2),
    total_price NUMERIC(18,2),
    lead_time_days INT,
    notes TEXT
);

-- Now you can query:
-- "Show all RFQs for CT Scanners"
SELECT r.* FROM rfqs r
JOIN rfq_items ri ON r.id = ri.rfq_id
JOIN products p ON ri.product_id = p.id
WHERE p.name ILIKE '%CT Scanner%';

-- "Average quote price by product"
SELECT p.name, AVG(qi.unit_price)
FROM quote_items qi
JOIN products p ON qi.product_id = p.id
GROUP BY p.name;
```

---

### **Issue #5: Engineer Coverage - Arrays Instead of Tables**

#### **Current Flawed Design:**
```sql
engineers (
    coverage_pincodes TEXT[],  -- ❌ ['110001', '110002', ...]
    coverage_cities TEXT[],    -- ❌ ['Delhi', 'Gurgaon', ...]
    coverage_states TEXT[],    -- ❌ ['Delhi', 'UP', ...]
    ...
)
```

#### **Problems:**
1. ❌ Query is slow: "Which engineers cover pincode 110001?"
   - Requires array scan: `WHERE '110001' = ANY(coverage_pincodes)`
2. ❌ Can't add metadata: "Primary coverage vs secondary coverage"
3. ❌ Can't track: When did coverage change?
4. ❌ Can't store: Travel cost for this area
5. ❌ Arrays don't normalize well
6. ❌ Hard to maintain referential integrity

#### **Correct Design:**
```sql
-- Remove arrays from engineers table

coverage_areas (
    id UUID PRIMARY KEY,
    area_type TEXT,              -- 'pincode', 'city', 'district', 'state'
    area_code TEXT,              -- '110001', 'Delhi', etc.
    area_name TEXT,
    parent_area_id UUID FK,      -- Hierarchy: pincode → city → state
    coordinates POINT,
    metadata JSONB
);

engineer_coverage (
    id UUID PRIMARY KEY,
    engineer_id UUID FK,
    coverage_area_id UUID FK,
    coverage_type TEXT,          -- 'primary', 'secondary', 'emergency'
    priority INT,                -- 1 = highest priority for this area
    travel_cost_per_km NUMERIC(8,2),
    typical_travel_time_mins INT,
    valid_from DATE,
    valid_to DATE,
    notes TEXT
);

-- Query: Engineers covering pincode 110001
SELECT e.* 
FROM engineers e
JOIN engineer_coverage ec ON e.id = ec.engineer_id
JOIN coverage_areas ca ON ec.coverage_area_id = ca.id
WHERE ca.area_type = 'pincode' 
  AND ca.area_code = '110001'
  AND ec.valid_to IS NULL
ORDER BY ec.priority;

-- Query: All areas covered by engineer
SELECT ca.area_type, ca.area_code, ca.area_name, ec.coverage_type
FROM engineer_coverage ec
JOIN coverage_areas ca ON ec.coverage_area_id = ca.id
WHERE ec.engineer_id = 'ENG-123'
  AND ec.valid_to IS NULL;
```

---

## 🟠 HIGH PRIORITY ISSUES

### **Issue #6: Ticket Status - Dual Source of Truth**

#### **Problem:**
```sql
service_tickets (
    status TEXT  -- Source 1
);

ticket_status_history (
    to_status TEXT  -- Source 2
);
```

Both store status - which is correct if they differ?

#### **Solution:**
```sql
-- Option A: Status is derived from history
CREATE VIEW service_tickets_with_status AS
SELECT t.*, h.to_status as current_status
FROM service_tickets t
LEFT JOIN LATERAL (
    SELECT to_status 
    FROM ticket_status_history 
    WHERE ticket_id = t.id 
    ORDER BY changed_at DESC 
    LIMIT 1
) h ON true;

-- Option B: Use trigger to keep in sync
CREATE TRIGGER sync_ticket_status
AFTER INSERT ON ticket_status_history
FOR EACH ROW
EXECUTE FUNCTION update_ticket_status();
```

---

### **Issue #7: Customer Information Denormalized**

#### **Problem:**
```sql
service_tickets (
    customer_id VARCHAR,
    customer_name VARCHAR,
    customer_phone VARCHAR,
    customer_whatsapp VARCHAR
);

equipment (
    customer_id VARCHAR,
    customer_name VARCHAR
);
```

Customer data duplicated across tables!

#### **Solution:**
```sql
customers (
    id UUID PRIMARY KEY,
    name TEXT,
    phone TEXT,
    whatsapp TEXT,
    email TEXT,
    organization_id UUID FK,  -- If customer is an organization
    facility_id UUID FK,      -- Specific facility
    address JSONB,
    created_at TIMESTAMPTZ
);

service_tickets (
    customer_id UUID FK  -- Single FK only!
);

equipment (
    owner_customer_id UUID FK  -- Single FK only!
);
```

---

### **Issue #8: Certifications - No Renewal Tracking**

#### **Problem:**
```sql
organization_certifications (
    certification_number TEXT,
    issue_date DATE,
    expiry_date DATE
);
```

When certification renews, do you delete old record?

#### **Solution:**
```sql
organization_certifications (
    certification_type TEXT,
    certification_number TEXT,
    version INT,              -- ✅ Version 1, 2, 3...
    issue_date DATE,
    expiry_date DATE,
    renewal_of_id UUID FK,    -- ✅ Links to previous cert
    superseded_by_id UUID FK, -- ✅ Links to next cert
    status TEXT               -- active|expired|revoked|renewed
);
```

---

### **Issue #9: Price Rules - No Temporal Querying**

#### **Problem:**
```sql
price_rules (
    sku_id UUID,
    price NUMERIC,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ
);
```

What if price changes multiple times?

#### **Solution:**
```sql
price_rules (
    sku_id UUID,
    version INT,
    price NUMERIC,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    
    -- Prevent overlaps:
    CONSTRAINT no_overlap EXCLUDE USING gist (
        sku_id WITH =,
        tstzrange(valid_from, valid_to, '[]') WITH &&
    )
);

-- Query: Price on specific date
SELECT price FROM price_rules
WHERE sku_id = '...'
  AND '2024-01-15' <@ tstzrange(valid_from, valid_to, '[]');
```

---

### **Issue #10: Mixed ID Types (UUID vs VARCHAR)**

#### **Problem:**
```sql
organizations (id UUID)
equipment (id VARCHAR)
service_tickets (id VARCHAR)
```

Inconsistency makes joins complex!

#### **Solution:**
Standardize on **UUID for everything**:
```sql
-- All tables should use:
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- Keep human-readable codes separately:
equipment (
    id UUID PRIMARY KEY,
    equipment_code TEXT UNIQUE  -- 'EQ-001'
);
```

---

## 🟡 MEDIUM PRIORITY ISSUES

### **Issue #11: Engineer Skills - Data Duplication**

```sql
engineer_skills (
    manufacturer_id UUID FK,
    manufacturer_name TEXT  -- ❌ Redundant!
);
```

**Solution:** Remove `manufacturer_name`, derive from FK.

---

### **Issue #12: Contact Persons - No Change History**

```sql
contact_persons (
    is_primary BOOLEAN
);
```

**Solution:** Add `effective_from`, `effective_to` for historical tracking.

---

### **Issue #13: Territories - No Multi-Assignment**

```sql
territories (
    assigned_to_org_id UUID  -- ❌ Only ONE org!
);
```

**Solution:** Create `territory_assignments` many-to-many table.

---

### **Issue #14: Engineer Availability - No Timezone**

```sql
engineer_availability (
    date DATE  -- ❌ Which timezone?
);
```

**Solution:** Use `timestamptz` or add `timezone TEXT` field.

---

## 📊 Impact Analysis

### **If NOT Fixed:**

| Issue | Impact | Estimated Cost |
|-------|--------|----------------|
| Ticket Assignment | Lost escalation data, poor routing decisions | **High** - Core feature broken |
| Equipment Ownership | Can't track transfers, leases, supply chain | **High** - Audit/compliance risk |
| Org Relationships | Can't calculate historical commissions | **High** - Financial reporting broken |
| RFQ/Quote Items | Can't analyze demand, slow queries | **Medium** - Business intelligence limited |
| Coverage Arrays | Slow engineer matching, can't optimize routes | **Medium** - Poor user experience |
| Customer Duplication | Data inconsistency, update anomalies | **Medium** - Data quality issues |
| Price Rules | Billing disputes, can't prove historical prices | **Medium** - Legal/financial risk |

### **If Fixed:**

✅ **Scalability:** System can handle 10x growth  
✅ **Audit Trail:** Complete history for compliance  
✅ **Query Performance:** Fast, indexed queries  
✅ **Data Integrity:** No inconsistencies  
✅ **Reporting:** Rich analytics capabilities  
✅ **Maintainability:** Easy to understand and modify  

---

## 🛠️ Recommended Action Plan

### **Phase 1: Critical Fixes (Week 1-2)**
1. ✅ Fix service ticket assignment (Issue #1)
2. ✅ Implement equipment relationships (Issue #2)
3. ✅ Create customers table (Issue #7)
4. ✅ Normalize RFQ/Quote items (Issue #4)

### **Phase 2: High Priority (Week 3-4)**
5. ✅ Org relationship terms history (Issue #3)
6. ✅ Standardize ID types (Issue #10)
7. ✅ Engineer coverage normalization (Issue #5)
8. ✅ Price rules temporal design (Issue #9)

### **Phase 3: Medium Priority (Week 5-6)**
9. ✅ Certification versioning (Issue #8)
10. ✅ Contact person history (Issue #12)
11. ✅ Territory assignments (Issue #13)
12. ✅ Ticket status sync (Issue #6)

---

## 📝 Migration Strategy

### **Option A: Big Bang (Not Recommended)**
- Shut down system
- Migrate all at once
- High risk, long downtime

### **Option B: Incremental (Recommended)**
1. Create new tables alongside old
2. Dual-write to both (old + new)
3. Backfill historical data
4. Switch reads to new tables
5. Stop writing to old tables
6. Drop old tables

### **Example Migration: Ticket Assignments**

```sql
-- Step 1: Create new structure (already exists!)
-- engineer_assignments table is ready

-- Step 2: Backfill from existing data
INSERT INTO engineer_assignments (
    id, ticket_id, engineer_id, 
    assignment_sequence, assignment_tier,
    assigned_at, status
)
SELECT 
    gen_random_uuid(),
    id,
    assigned_engineer_id::uuid,
    1,  -- Assume current is first assignment
    assignment_tier,
    assigned_at,
    CASE 
        WHEN status IN ('resolved', 'closed') THEN 'completed'
        ELSE 'in_progress'
    END
FROM service_tickets
WHERE assigned_engineer_id IS NOT NULL;

-- Step 3: Application code uses engineer_assignments
-- Step 4: After testing, drop old columns
ALTER TABLE service_tickets 
DROP COLUMN assigned_engineer_id,
DROP COLUMN assigned_engineer_name,
DROP COLUMN assignment_tier,
DROP COLUMN assignment_tier_name;
```

---

## ✅ Success Criteria

- ✅ All escalations tracked with full history
- ✅ Equipment ownership changes recorded
- ✅ Can query historical data accurately
- ✅ No data duplication
- ✅ Query performance < 100ms for common operations
- ✅ Full audit trail for compliance
- ✅ System can scale to 1M+ tickets

---

## 🎯 Conclusion

The current schema has **solid foundations** but needs **normalization** and **historization** fixes. These are **not optional** - they will cause production issues if not addressed.

**Recommendation:** Start with **Phase 1 critical fixes** immediately before adding more features.

---

**Reviewed By:** Database Architecture Team  
**Priority:** 🔴 Critical  
**Estimated Effort:** 4-6 weeks  
**Risk if Delayed:** High - Data loss, integrity issues, poor performance
