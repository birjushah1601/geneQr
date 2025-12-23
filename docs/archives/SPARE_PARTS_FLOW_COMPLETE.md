# Spare Parts Flow - Complete Chronology

## Overview

This document explains the complete flow of spare parts from catalog to service ticket fulfillment, including marketplace integration.

---

## Architecture Overview

```
spare_parts_catalog (50 parts)
    ↓
equipment_spare_parts (113 associations) ← Links parts to equipment types
    ↓
equipment_catalog (23 models) ← Product catalog for marketplace
    ↓
equipment_registry (73 installed units) ← Actual deployed equipment at customer sites
    ↓
service_tickets ← Service requests for equipment issues
    ↓
ticket_parts ← Parts assigned to tickets for fulfillment
    ↓
[SHIPMENT TO CLIENT]
```

---

## The Complete Flow

### **Phase 1: Catalog Setup (Marketplace)**

#### 1.1 Spare Parts Catalog
**Table:** `spare_parts_catalog`

**Purpose:** Master catalog of all available spare parts

**Example:**
```sql
SELECT part_number, part_name, category, unit_price 
FROM spare_parts_catalog 
WHERE part_number = 'XR-TUBE-001';
```

```
part_number  | part_name             | category | unit_price
-------------+-----------------------+----------+------------
XR-TUBE-001  | X-Ray Tube Assembly   | X-Ray    | 12500.00
```

**Use Case:** 
- Parts available for sale in marketplace
- Parts that can be ordered for maintenance
- Replacement parts for equipment

#### 1.2 Equipment Catalog (Product Models)
**Table:** `equipment_catalog`

**Purpose:** Product catalog of equipment models available in marketplace

**Example:**
```sql
SELECT product_code, product_name, manufacturer_name, category 
FROM equipment_catalog 
WHERE category = 'X-Ray' LIMIT 3;
```

```
product_code     | product_name                    | manufacturer_name       | category
-----------------+---------------------------------+-------------------------+----------
XR-SOU-ALPHA     | X-Ray System Alpha              | SouthCare Distributors  | X-Ray
550e8400-...-005 | Digital X-Ray System CXDI-410C  | Canon Medical Systems   | X-Ray
d2e959d0-...-ef9 | Multix Fusion                   | Siemens Healthineers    | X-Ray
```

**Use Case:**
- Equipment products listed in marketplace
- Specifications for buyers
- Pricing information

#### 1.3 Equipment-Parts Association
**Table:** `equipment_spare_parts`

**Purpose:** Links spare parts to equipment models

**Example:**
```sql
SELECT 
    ec.product_name,
    sp.part_name,
    esp.is_critical
FROM equipment_spare_parts esp
JOIN equipment_catalog ec ON esp.equipment_catalog_id = ec.id
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
WHERE ec.category = 'X-Ray'
LIMIT 5;
```

```
product_name                   | part_name              | is_critical
-------------------------------+------------------------+-------------
X-Ray System Alpha             | X-Ray Tube Assembly    | true
X-Ray System Alpha             | Flat Panel Detector    | true
X-Ray System Alpha             | Collimator Assembly    | false
Digital X-Ray System CXDI-410C | X-Ray Tube Assembly    | true
Digital X-Ray System CXDI-410C | X-Ray Filter Set       | false
```

**Use Case:**
- Show compatible parts for each equipment model in marketplace
- Suggest parts when equipment is purchased
- Parts bundling options

**Current Status:**
- ✅ 113 associations created
- ✅ X-Ray equipment linked to 5 parts
- ✅ MRI equipment linked to 5 parts
- ✅ CT equipment linked to 4 parts
- ✅ Patient Monitor linked to 7 parts
- ✅ Ventilator linked to 7 parts

---

### **Phase 2: Equipment Deployment (Field Service)**

#### 2.1 Equipment Registry (Installed Units)
**Table:** `equipment_registry`

**Purpose:** Actual equipment deployed at customer sites

**Example:**
```sql
SELECT 
    id,
    equipment_name,
    serial_number,
    customer_name,
    installation_location,
    equipment_catalog_id
FROM equipment_registry 
WHERE id = 'REG-CAN-XR-005';
```

```
id              | equipment_name                  | serial_number    | customer_name      | installation_location | equipment_catalog_id
----------------+---------------------------------+------------------+--------------------+-----------------------+----------------------
REG-CAN-XR-005  | Digital X-Ray System CXDI-410C  | SN-CAN-XR-005    | Fortis Hospital    | X-Ray Room 2          | 550e8400-...-005
```

**Link to Catalog:**
```sql
-- equipment_registry has foreign key to equipment_catalog
equipment_catalog_id → equipment_catalog(id)
```

**This means:**
1. ✅ Each installed equipment links to a product model
2. ✅ Through equipment_catalog, we know which parts are compatible
3. ✅ Can lookup spare parts via: registry → catalog → equipment_spare_parts → spare_parts_catalog

**Example Query:**
```sql
-- Get all compatible spare parts for an installed equipment
SELECT 
    er.equipment_name,
    er.serial_number,
    sp.part_name,
    sp.unit_price
FROM equipment_registry er
JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
JOIN equipment_spare_parts esp ON ec.id = esp.equipment_catalog_id
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
WHERE er.id = 'REG-CAN-XR-005';
```

**Current Status:**
- ✅ 73 equipment units deployed
- ✅ All linked to equipment_catalog
- ✅ Can lookup compatible parts for each unit

---

### **Phase 3: Service Request (Fault Occurs)**

#### 3.1 Service Ticket Creation
**Table:** `service_tickets`

**When:** Equipment fails or needs maintenance

**Example:**
```sql
-- User scans QR code: QR-CAN-XR-005
-- System creates ticket
INSERT INTO service_tickets (
    id,
    equipment_id,
    customer_name,
    issue_description,
    priority,
    status
) VALUES (
    'TICKET-123',
    'REG-CAN-XR-005',
    'Fortis Hospital',
    'X-Ray tube not firing properly',
    'high',
    'open'
);
```

**Workflow:**
1. ✅ Technician scans QR code on equipment
2. ✅ QR code opens service-request page: `?qr=QR-CAN-XR-005`
3. ✅ Page fetches equipment details from equipment_registry
4. ✅ Equipment info pre-filled in form
5. ✅ Technician describes issue
6. ✅ Ticket created and linked to equipment

---

### **Phase 4: Parts Assignment (Diagnosis)**

#### 4.1 Parts Assignment Modal
**Table:** `ticket_parts`

**When:** Engineer diagnoses issue and identifies required parts

**Modal Features:**
- ✅ Opens on service-request page
- ✅ Calls API: `GET /api/v1/catalog/parts`
- ✅ Shows all available spare parts (50 parts)
- ✅ Filter by category (X-Ray, MRI, CT, etc.)
- ✅ Search by part name/number
- ✅ Select quantity needed
- ✅ Shows price, stock status, lead time

**Example Scenario:**
```
Ticket: TICKET-123
Issue: "X-Ray tube not firing properly"
Equipment: Digital X-Ray System CXDI-410C

Engineer opens "Add Parts" modal:
  → Shows parts compatible with X-Ray equipment
  → Searches for "tube"
  → Finds: XR-TUBE-001 - X-Ray Tube Assembly - $12,500
  → Selects quantity: 1
  → Adds to ticket
```

**Database Insert:**
```sql
INSERT INTO ticket_parts (
    ticket_id,
    spare_part_id,
    quantity_required,
    unit_price,
    total_price,
    status,
    is_critical
) VALUES (
    'TICKET-123',
    '<uuid-of-XR-TUBE-001>',
    1,
    12500.00,
    12500.00,
    'pending',
    true
);
```

**Status Flow:**
```
pending → ordered → received → installed
```

**Current Status:**
- ✅ ticket_parts table ready
- ✅ Modal configured and working
- ✅ API endpoint created: /api/v1/catalog/parts
- ✅ Frontend integrated

---

### **Phase 5: Parts Fulfillment (Shipping)**

#### 5.1 Order Processing

**Status: `pending`**
```sql
-- Parts assigned to ticket, awaiting approval
SELECT * FROM ticket_parts 
WHERE ticket_id = 'TICKET-123' AND status = 'pending';
```

**Manager reviews:**
- Part needed: X-Ray Tube Assembly
- Cost: $12,500
- Customer: Fortis Hospital
- Priority: High
- Approves order

**Status: `ordered`**
```sql
-- Order placed with supplier/warehouse
UPDATE ticket_parts 
SET status = 'ordered', 
    assigned_at = NOW(),
    assigned_by = 'manager@company.com'
WHERE ticket_id = 'TICKET-123';
```

**Warehouse prepares shipment:**
- Picks part from inventory
- Packages for shipping
- Creates shipping label
- Dispatches to customer location

**Status: `received`**
```sql
-- Part delivered to customer site
UPDATE ticket_parts 
SET status = 'received',
    notes = 'Delivered to Fortis Hospital - X-Ray Room 2'
WHERE ticket_id = 'TICKET-123';
```

**Status: `installed`**
```sql
-- Engineer installs part and completes repair
UPDATE ticket_parts 
SET status = 'installed',
    quantity_used = 1,
    installed_at = NOW()
WHERE ticket_id = 'TICKET-123';

-- Close ticket
UPDATE service_tickets 
SET status = 'resolved',
    resolution_notes = 'X-Ray tube replaced successfully. Equipment operational.'
WHERE id = 'TICKET-123';
```

---

## Marketplace Use Cases

### **Use Case 1: Equipment Purchase with Parts Bundle**

**Scenario:** Hospital buying X-Ray System

```sql
-- Show equipment in marketplace
SELECT * FROM equipment_catalog WHERE id = '550e8400-...-005';

-- Show recommended spare parts bundle
SELECT 
    sp.part_name,
    sp.unit_price,
    esp.is_critical
FROM equipment_spare_parts esp
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
WHERE esp.equipment_catalog_id = '550e8400-...-005'
  AND esp.is_critical = true;
```

**Output:**
```
part_name              | unit_price | is_critical
-----------------------+------------+-------------
X-Ray Tube Assembly    | 12500.00   | true
Flat Panel Detector    | 28000.00   | true
```

**Marketplace Offer:**
```
X-Ray System: $150,000
+ Spare Parts Bundle: $40,500
  - X-Ray Tube (backup)
  - Flat Panel Detector (backup)
------------------------
Total: $190,500 (Save 10%)
```

### **Use Case 2: Preventive Maintenance Parts Order**

**Scenario:** Hospital orders parts for scheduled maintenance

```sql
-- Find consumable parts for their X-Ray
SELECT 
    sp.part_name,
    sp.unit_price,
    sp.part_type
FROM equipment_registry er
JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
JOIN equipment_spare_parts esp ON ec.id = esp.equipment_catalog_id
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
WHERE er.customer_id = '<fortis-hospital-id>'
  AND er.category = 'X-Ray'
  AND sp.part_type = 'consumable';
```

**Output:**
```
part_name         | unit_price | part_type
------------------+------------+------------
X-Ray Filter Set  | 450.00     | consumable
```

**Order in Marketplace:**
```
Preventive Maintenance Kit
- X-Ray Filter Set × 4 = $1,800
- Anti-Scatter Grid × 2 = $2,400
------------------------------
Total: $4,200
```

### **Use Case 3: Emergency Part Rush Order**

**Scenario:** Equipment breaks down, need urgent part

**Flow:**
1. ✅ Engineer creates service ticket from mobile
2. ✅ Adds critical part: X-Ray Tube Assembly
3. ✅ Marks as urgent/high priority
4. ✅ System alerts warehouse
5. ✅ Warehouse ships via express delivery
6. ✅ Track shipment in real-time
7. ✅ Engineer gets notification when delivered
8. ✅ Install and close ticket

**Database Tracking:**
```sql
SELECT 
    st.id as ticket_id,
    st.priority,
    tp.status as part_status,
    sp.part_name,
    tp.assigned_at,
    tp.installed_at,
    st.status as ticket_status
FROM service_tickets st
JOIN ticket_parts tp ON st.id = tp.ticket_id
JOIN spare_parts_catalog sp ON tp.spare_part_id = sp.id
WHERE st.priority = 'high'
  AND tp.is_critical = true;
```

---

## Database Schema Relationships

### **Visual Diagram**

```
┌─────────────────────────┐
│  spare_parts_catalog    │ ← Master parts catalog (50 parts)
│  - id (PK)              │
│  - part_number          │
│  - part_name            │
│  - category             │
│  - unit_price           │
└──────────┬──────────────┘
           │
           │ Referenced by
           │
┌──────────▼──────────────┐
│ equipment_spare_parts   │ ← Links parts to equipment models (113 links)
│  - equipment_catalog_id │──┐
│  - spare_part_id        │  │
│  - is_critical          │  │
└─────────────────────────┘  │
                             │
                             │ References
                             │
           ┌─────────────────▼──────────────┐
           │   equipment_catalog            │ ← Product catalog for marketplace (23 models)
           │   - id (PK)                    │
           │   - product_code               │
           │   - product_name               │
           │   - category                   │
           └──────────┬─────────────────────┘
                      │
                      │ Referenced by
                      │
           ┌──────────▼─────────────────────┐
           │   equipment_registry           │ ← Deployed equipment (73 units)
           │   - id (PK)                    │
           │   - equipment_catalog_id (FK)  │
           │   - serial_number              │
           │   - customer_name              │
           │   - qr_code                    │
           └──────────┬─────────────────────┘
                      │
                      │ Referenced by
                      │
           ┌──────────▼─────────────────────┐
           │   service_tickets              │ ← Service requests
           │   - id (PK)                    │
           │   - equipment_id (FK)          │
           │   - issue_description          │
           │   - status                     │
           └──────────┬─────────────────────┘
                      │
                      │ Referenced by
                      │
           ┌──────────▼─────────────────────┐
           │   ticket_parts                 │ ← Parts assigned to tickets
           │   - ticket_id (FK)             │
           │   - spare_part_id (FK)         │
           │   - quantity_required          │
           │   - status                     │
           └────────────────────────────────┘
```

---

## Verification Queries

### **Check Complete Chain**

```sql
-- Verify complete chain: Registry → Catalog → Parts
SELECT 
    er.id as registry_id,
    er.equipment_name,
    ec.product_name as catalog_model,
    sp.part_name,
    sp.unit_price,
    esp.is_critical
FROM equipment_registry er
JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
JOIN equipment_spare_parts esp ON ec.id = esp.equipment_catalog_id
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
WHERE er.id = 'REG-CAN-XR-005'
ORDER BY sp.unit_price DESC;
```

### **Count Parts per Equipment**

```sql
-- How many parts available for each deployed equipment?
SELECT 
    er.equipment_name,
    er.customer_name,
    COUNT(DISTINCT sp.id) as available_parts,
    COUNT(DISTINCT CASE WHEN esp.is_critical THEN sp.id END) as critical_parts
FROM equipment_registry er
JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
JOIN equipment_spare_parts esp ON ec.id = esp.equipment_catalog_id
JOIN spare_parts_catalog sp ON esp.spare_part_id = sp.id
GROUP BY er.id, er.equipment_name, er.customer_name
ORDER BY available_parts DESC
LIMIT 10;
```

### **Find Most Used Parts**

```sql
-- Which parts are assigned to tickets most often?
SELECT 
    sp.part_name,
    sp.category,
    COUNT(tp.id) as times_assigned,
    SUM(tp.quantity_required) as total_quantity,
    SUM(tp.total_price) as total_revenue
FROM ticket_parts tp
JOIN spare_parts_catalog sp ON tp.spare_part_id = sp.id
GROUP BY sp.id, sp.part_name, sp.category
ORDER BY times_assigned DESC
LIMIT 10;
```

---

## Summary

### **Complete Flow Verified ✅**

1. **Marketplace (Catalog Phase)**
   - ✅ spare_parts_catalog: 50 parts available
   - ✅ equipment_catalog: 23 equipment models
   - ✅ equipment_spare_parts: 113 part-to-model associations

2. **Field Deployment**
   - ✅ equipment_registry: 73 deployed units
   - ✅ Each unit links to equipment_catalog
   - ✅ Through catalog, parts are accessible

3. **Service Tickets**
   - ✅ Tickets created when equipment fails
   - ✅ Linked to equipment_registry
   - ✅ Parts assigned via ticket_parts table

4. **Fulfillment**
   - ✅ Parts status tracking: pending → ordered → received → installed
   - ✅ Shipment to client location
   - ✅ Installation and ticket closure

### **The Chronology is Correct! ✅**

**YES - These are the SAME spare parts used for:**
1. ✅ Marketplace catalog (selling with equipment)
2. ✅ Service tickets (repair/replacement)
3. ✅ Fulfillment (shipping to clients)

**The link is:**
```
spare_parts_catalog 
  → equipment_spare_parts (compatibility)
  → equipment_catalog (product models)
  → equipment_registry (deployed units)
  → service_tickets (issues)
  → ticket_parts (fulfillment)
```

**Everything is connected and ready to use! 🎉**
