# Equipment & Parts System - Complete Documentation

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Equipment Catalog System](#equipment-catalog-system)
- [Equipment Registry & QR Codes](#equipment-registry--qr-codes)
- [Service Tickets](#service-tickets)
- [Parts Management](#parts-management)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)
- [Data Flow](#data-flow)

---

## Architecture Overview

### Core Concept

The system follows a **3-tier equipment architecture**:

```
equipment_catalog (Master Products)
    ↓ Manufacturers define products & assign compatible parts
equipment_part_assignments (Parts per Product Type)
    ↓ Links to parts catalog
spare_parts_catalog (Parts Master Inventory)
    ↓ Products are installed as
equipment_registry (Specific Units at Customer Locations)
    ↓ Customers create support requests
service_tickets (Support Tickets)
    ↓ Admins assign parts from catalog
ticket_parts (Parts Assigned to Tickets)
```

### Key Principles

1. **equipment_catalog** = Master product definitions by manufacturers (e.g., "MAGNETOM Vida 3T MRI Scanner")
2. **equipment_registry** = Specific installed units at customer sites with serial numbers
3. **service_tickets** = Always reference equipment_registry (installed units), NOT catalog
4. **Parts are assigned at catalog level** but used at ticket level
5. **QR codes** are unique per installed equipment unit

---

## Equipment Catalog System

### Purpose
Master product catalog maintained by manufacturers defining equipment types and their specifications.

### Table: `equipment_catalog`

```sql
CREATE TABLE equipment_catalog (
    id UUID PRIMARY KEY,
    product_name VARCHAR(255),           -- e.g., "MAGNETOM Vida 3T MRI Scanner"
    manufacturer_name VARCHAR(255),       -- e.g., "Siemens Healthineers"
    model_number VARCHAR(100),
    category VARCHAR(100),                -- e.g., "MRI", "CT", "X-Ray"
    specifications JSONB,                 -- Technical specs
    certifications JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Current Data

**14 Equipment Catalog Items:**

| Category | Product | Manufacturer | Parts |
|----------|---------|--------------|-------|
| MRI | MAGNETOM Vida 3T MRI Scanner | Siemens Healthineers | 2 |
| MRI | SIGNA Explorer 1.5T | GE Healthcare | 2 |
| CT | CT Scanner Nova | Wipro GE Healthcare | 2 |
| CT | Ingenuity CT 128-slice | Philips Healthcare | 2 |
| X-Ray | X-Ray System Alpha | SouthCare Distributors | 2 |
| X-Ray | Digital X-Ray System CXDI-410C | Canon Medical | 2 |
| Ultrasound | LOGIQ E10 Ultrasound | GE Healthcare | 2 |
| Patient Monitor | Patient Monitor Visionary | Medtronic India | 2 |
| Ventilator | Savina 300 Ventilator | Dräger Medical | 2 |
| Dialysis | Fresenius 5008 Dialysis Machine | Fresenius Medical Care | 2 |
| Infusion Pump | Infusion Pump Lite | Philips Healthcare India | 2 |

**Total:** 14 catalog items, 11 with parts assigned

---

## Equipment Registry & QR Codes

### Purpose
Tracks specific equipment units installed at customer locations. Each unit has a unique serial number and QR code.

### Table: `equipment_registry`

```sql
CREATE TABLE equipment_registry (
    id VARCHAR(32) PRIMARY KEY,           -- e.g., "REG-MRI-VIDA-001"
    equipment_catalog_id UUID,            -- FK to equipment_catalog
    qr_code VARCHAR(255) UNIQUE,          -- Unique QR code
    serial_number VARCHAR(255) UNIQUE,    -- e.g., "SN-MRI-VIDA-001"
    equipment_name VARCHAR(500),
    manufacturer_name VARCHAR(255),
    customer_id VARCHAR(32),
    customer_name VARCHAR(500),
    installation_location TEXT,           -- e.g., "Apollo Hospital - Imaging Wing"
    installation_address JSONB,
    installation_date DATE,
    warranty_expiry DATE,
    status VARCHAR(50),                   -- operational, down, under_maintenance, decommissioned
    specifications JSONB,
    qr_code_url TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    
    CONSTRAINT fk_equipment_registry_catalog 
        FOREIGN KEY (equipment_catalog_id) 
        REFERENCES equipment_catalog(id)
);
```

### QR Code System

**QR Code Format:**
- Format: `QR-{TYPE}-{MODEL}-{SEQUENCE}`
- Example: `QR-MRI-VIDA-001`
- Each QR code is **globally unique**

**QR Code URL:**
- Pattern: `https://app.example.com/qr/{equipment_registry_id}`
- Scanning redirects to equipment details or ticket creation

**QR Code Workflow:**
1. Manufacturer creates equipment in catalog
2. Equipment installed at customer location → Registry entry created
3. Unique QR code generated and printed
4. QR code attached to physical equipment
5. Customer/technician scans QR → Opens ticket creation form
6. Ticket automatically linked to correct equipment unit

### Current Registry Data

**23 Equipment Registry Entries:**

| Registry ID | Serial | Equipment | Location | Status |
|-------------|--------|-----------|----------|--------|
| REG-MRI-VIDA-001 | SN-MRI-VIDA-001 | MAGNETOM Vida 3T | Apollo Hospital | operational |
| REG-MRI-VIDA-002 | SN-MRI-VIDA-002 | MAGNETOM Vida 3T | Fortis Hospital | operational |
| REG-CT-NOVA-001 | SN-CT-NOVA-001 | CT Scanner Nova | AIIMS Delhi | operational |
| REG-XR-ALPHA-001 | SN-XR-ALPHA-001 | X-Ray System | Local Clinic | operational |
| ... | ... | ... | ... | ... |

**Total:** 23 installed equipment units across 11 equipment types

### Equipment Registration Flow

```
1. Manufacturer ships equipment
   ↓
2. Customer receives equipment
   ↓
3. Installation team registers in system:
   - Scans/enters serial number
   - Links to equipment_catalog entry
   - Records installation location
   - Generates QR code
   ↓
4. QR code printed and affixed to equipment
   ↓
5. Equipment ready for service requests
```

---

## Service Tickets

### Purpose
Customer support requests for installed equipment. Always references `equipment_registry` entries.

### Table: `service_tickets`

```sql
CREATE TABLE service_tickets (
    id VARCHAR(50) PRIMARY KEY,
    ticket_number VARCHAR(50) UNIQUE,     -- e.g., "TKT-20251219-123456"
    equipment_id VARCHAR(32),             -- FK to equipment_registry.id
    equipment_name VARCHAR(500),
    customer_id VARCHAR(32),
    customer_name VARCHAR(500),
    issue_category VARCHAR(100),
    issue_description TEXT,
    priority VARCHAR(20),                 -- low, medium, high, critical
    status VARCHAR(50),                   -- open, assigned, in_progress, resolved, closed
    assigned_engineer_id UUID,
    assigned_engineer_name VARCHAR(255),
    parts_used JSONB DEFAULT '[]',        -- Legacy field (not used)
    created_at TIMESTAMP,
    resolved_at TIMESTAMP,
    
    CONSTRAINT fk_service_tickets_equipment 
        FOREIGN KEY (equipment_id) 
        REFERENCES equipment_registry(id)
);
```

### Ticket Creation Flow

**Method 1: Via QR Code (Recommended)**
```
1. Customer scans QR code on equipment
   ↓
2. System reads equipment_registry_id from QR
   ↓
3. Ticket creation form pre-filled with:
   - Equipment details
   - Serial number
   - Location
   - Customer info
   ↓
4. Customer describes issue
   ↓
5. Ticket created and assigned to engineer
```

**Method 2: Manual Selection**
```
1. Customer opens ticket creation form
   ↓
2. Selects equipment from their registered equipment list
   ↓
3. System fetches equipment_registry entry
   ↓
4. Ticket created with equipment reference
```

### Current Tickets

**8 Service Tickets:**

| Ticket Number | Equipment | Status | Parts Assigned |
|---------------|-----------|--------|----------------|
| SR-DEM-0001 | MRI Scanner | open | 0 (ready) |
| TKT-20251015-175800 | Patient Monitor | open | 0 (ready) |
| TKT-20251015-175801 | Ventilator | open | 0 (ready) |
| TKT-20251212-202732 | CT Scanner | open | 0 (ready) |

All tickets cleared and ready for fresh parts assignment.

---

## Parts Management

### Architecture

```
equipment_catalog
    ↓ has compatible parts defined in
equipment_part_assignments
    ↓ references
spare_parts_catalog
    ↓ assigned to tickets via
ticket_parts (junction table)
```

### Table: `equipment_part_assignments`

Links equipment types to their compatible spare parts.

```sql
CREATE TABLE equipment_part_assignments (
    id UUID PRIMARY KEY,
    equipment_catalog_id UUID,            -- FK to equipment_catalog
    spare_part_id UUID,                   -- FK to spare_parts_catalog
    quantity_required INTEGER DEFAULT 1,
    is_critical BOOLEAN DEFAULT false,
    installation_complexity VARCHAR(50),
    
    CONSTRAINT fk_equipment_catalog 
        FOREIGN KEY (equipment_catalog_id) 
        REFERENCES equipment_catalog(id),
    CONSTRAINT fk_spare_part 
        FOREIGN KEY (spare_part_id) 
        REFERENCES spare_parts_catalog(id)
);
```

**Current Data:** 22 equipment-parts assignments

### Table: `spare_parts_catalog`

Master inventory of all spare parts.

```sql
CREATE TABLE spare_parts_catalog (
    id UUID PRIMARY KEY,
    part_number VARCHAR(100) UNIQUE,      -- e.g., "MRI-COIL-HEAD-8CH"
    part_name VARCHAR(255),               -- e.g., "Head Coil 8-Channel"
    category VARCHAR(100),                -- component, consumable, accessory
    unit_price NUMERIC(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    stock_status VARCHAR(50),             -- in_stock, low_stock, out_of_stock
    lead_time_days INTEGER,
    is_available BOOLEAN DEFAULT true,
    supplier_info JSONB,
    created_at TIMESTAMP
);
```

**Current Data:** 16 spare parts available

**Example Parts:**

| Part Number | Part Name | Category | Price | Critical |
|-------------|-----------|----------|-------|----------|
| MRI-COIL-HEAD-8CH | Head Coil 8-Channel | accessory | $12,500 | Yes |
| MRI-CRYOGEN-HELIUM | Liquid Helium Refill | consumable | $8,500 | No |
| VENT-FLOW-SENSOR | Flow Sensor | component | $650 | Yes |
| XR-PANEL-FLAT | Flat Panel Detector | component | $45,000 | Yes |

### Table: `ticket_parts`

**Junction table** storing which parts are assigned to which tickets.

```sql
CREATE TABLE ticket_parts (
    id UUID PRIMARY KEY,
    ticket_id VARCHAR(50),                -- FK to service_tickets
    spare_part_id UUID,                   -- FK to spare_parts_catalog
    quantity_required INTEGER DEFAULT 1,
    quantity_used INTEGER,
    is_critical BOOLEAN DEFAULT false,
    status VARCHAR(50),                   -- pending, ordered, received, installed, returned
    unit_price NUMERIC(12,2),
    total_price NUMERIC(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    notes TEXT,
    assigned_by VARCHAR(100),
    assigned_at TIMESTAMP,
    installed_at TIMESTAMP,
    
    CONSTRAINT fk_ticket_parts_ticket 
        FOREIGN KEY (ticket_id) 
        REFERENCES service_tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_parts_part 
        FOREIGN KEY (spare_part_id) 
        REFERENCES spare_parts_catalog(id)
);
```

### Parts Assignment Flow

```
1. Ticket created for equipment (e.g., MRI Scanner)
   ↓
2. System looks up equipment_registry → equipment_catalog_id
   ↓
3. Query equipment_part_assignments for this catalog_id
   ↓
4. Display available parts to admin/engineer
   ↓
5. Admin selects parts to assign (with quantity)
   ↓
6. Parts saved to ticket_parts table
   ↓
7. Status tracked: pending → ordered → received → installed
```

### Getting Parts for a Ticket

**Database Function:**
```sql
-- Returns parts available for equipment type
get_parts_for_ticket(ticket_id) → parts list

-- Implementation:
1. Get ticket's equipment_id
2. Look up equipment_registry entry
3. Get equipment_catalog_id
4. Find parts in equipment_part_assignments
5. Join with spare_parts_catalog for details
6. Return parts with pricing, criticality, stock status
```

**Example Query:**
```sql
SELECT * FROM get_parts_for_ticket('TKT-20251219-123456');

-- Returns:
spare_part_id | part_number | part_name | unit_price | is_critical | stock_status
--------------|-------------|-----------|------------|-------------|-------------
uuid-1        | MRI-COIL... | Head Coil | 12500.00   | true        | in_stock
uuid-2        | MRI-CRYO... | Helium    | 8500.00    | false       | in_stock
```

---

## API Endpoints

### Equipment Registry

**GET /api/v1/equipment/registry**
- List all registered equipment
- Filters: customer_id, status, category

**GET /api/v1/equipment/registry/{id}**
- Get specific equipment details
- Includes: serial, QR code, location, warranty

**POST /api/v1/equipment/registry**
- Register new equipment installation
- Generates QR code automatically

**GET /api/v1/equipment/registry/qr/{qr_code}**
- Look up equipment by QR code
- Used when QR is scanned

### Service Tickets

**GET /api/v1/tickets**
- List all tickets
- Filters: status, priority, engineer, equipment

**GET /api/v1/tickets/{id}**
- Get ticket details
- Includes: equipment info, parts, comments

**POST /api/v1/tickets**
- Create new ticket
- Body: equipment_id (from registry), issue_description, priority

**PATCH /api/v1/tickets/{id}**
- Update ticket status, assignment, etc.

### Parts Management

**GET /api/v1/tickets/{id}/parts**
- Get parts assigned to this ticket
- Returns from `ticket_parts` table

**PATCH /api/v1/tickets/{id}/parts**
- Assign/update parts for a ticket
- Body:
```json
{
  "parts": [
    {
      "part_id": "uuid",
      "quantity": 2,
      "is_critical": true,
      "unit_price": 650.00
    }
  ]
}
```

**GET /api/v1/equipment/{equipment_catalog_id}/parts**
- Get available parts for equipment type
- Used for showing options during assignment

**GET /api/v1/spare-parts**
- List all spare parts in catalog
- Filters: category, availability, search

---

## Database Schema

### Complete ER Diagram

```
organizations
    ↓
engineers (via engineer_org_memberships)
    ↓ assigned to
service_tickets
    ↓ references
equipment_registry
    ↓ links to
equipment_catalog
    ↓ has parts via
equipment_part_assignments
    ↓ references
spare_parts_catalog
    ↓ assigned via
ticket_parts
```

### Key Relationships

| Parent | Child | Relationship | Description |
|--------|-------|--------------|-------------|
| equipment_catalog | equipment_registry | 1:N | One product → many installations |
| equipment_catalog | equipment_part_assignments | 1:N | One product → many compatible parts |
| spare_parts_catalog | equipment_part_assignments | 1:N | One part → used in many equipment types |
| equipment_registry | service_tickets | 1:N | One unit → many support tickets |
| service_tickets | ticket_parts | 1:N | One ticket → many assigned parts |
| spare_parts_catalog | ticket_parts | 1:N | One part → assigned to many tickets |

---

## Data Flow

### Complete Lifecycle Example

**1. Manufacturer Onboards Product**
```sql
INSERT INTO equipment_catalog (
    id, product_name, manufacturer_name, category
) VALUES (
    gen_random_uuid(),
    'MAGNETOM Vida 3T MRI Scanner',
    'Siemens Healthineers',
    'MRI'
);
```

**2. Manufacturer Defines Compatible Parts**
```sql
INSERT INTO equipment_part_assignments (
    equipment_catalog_id, spare_part_id, quantity_required, is_critical
) VALUES (
    'catalog-id', 'part-id-1', 2, true
);
```

**3. Equipment Installed at Customer Site**
```sql
INSERT INTO equipment_registry (
    id, equipment_catalog_id, serial_number, qr_code,
    customer_name, installation_location, status
) VALUES (
    'REG-MRI-001', 'catalog-id', 'SN-MRI-001', 'QR-MRI-001',
    'Apollo Hospital', 'Imaging Wing', 'operational'
);
```

**4. Customer Creates Ticket (via QR scan)**
```sql
INSERT INTO service_tickets (
    id, ticket_number, equipment_id, 
    customer_name, issue_description, priority
) VALUES (
    ksuid(), 'TKT-20251219-001', 'REG-MRI-001',
    'Apollo Hospital', 'MRI not producing images', 'high'
);
```

**5. Engineer Diagnoses and Assigns Parts**
```sql
INSERT INTO ticket_parts (
    ticket_id, spare_part_id, quantity_required, 
    is_critical, status, unit_price
) VALUES (
    'TKT-20251219-001', 'part-id-1', 1, 
    true, 'pending', 12500.00
);
```

**6. Parts Ordered → Received → Installed**
```sql
UPDATE ticket_parts 
SET status = 'installed', installed_at = NOW()
WHERE ticket_id = 'TKT-20251219-001';

UPDATE service_tickets
SET status = 'resolved', resolved_at = NOW()
WHERE id = 'TKT-20251219-001';
```

---

## Current System State

### Summary Statistics

| Metric | Count |
|--------|-------|
| Equipment Catalog Items | 14 |
| Equipment Registry Entries | 23 |
| Equipment Types with Parts | 11 |
| Total Part Assignments | 22 |
| Spare Parts in Catalog | 16 |
| Service Tickets | 8 |
| Parts Assigned to Tickets | 0 (cleared for testing) |

### Data Integrity

✅ All equipment_registry entries link to valid equipment_catalog entries
✅ All tickets reference valid equipment_registry entries  
✅ All equipment_part_assignments reference valid catalog and parts
✅ ticket_parts table exists and ready for use
✅ QR codes are unique per equipment unit
✅ Serial numbers are unique

### Ready for Testing

✅ Equipment registration working
✅ QR code system in place
✅ Ticket creation functional
✅ Parts catalog populated
✅ Parts assignment API ready
✅ Backend updated to use ticket_parts table

---

## Best Practices

### Equipment Registration
1. Always link to equipment_catalog (don't create standalone entries)
2. Generate unique QR codes for each unit
3. Record installation location and date
4. Set warranty expiry dates

### Ticket Creation
1. Always reference equipment_registry.id (not catalog)
2. Pre-fill from QR code when available
3. Set appropriate priority based on equipment criticality
4. Link to customer organization

### Parts Assignment
1. Show only compatible parts (from equipment_part_assignments)
2. Mark critical parts clearly
3. Check stock availability before assigning
4. Calculate total price (unit_price × quantity)
5. Track status through lifecycle

### QR Code Usage
1. Print QR codes on equipment labels
2. Include equipment_registry_id in QR data
3. Test QR scanning before deployment
4. Update QR code URL if equipment relocated

---

## Troubleshooting

### Issue: Parts not showing for ticket
**Check:**
1. Ticket references equipment_registry? (not old equipment table)
2. Registry entry has equipment_catalog_id set?
3. Catalog has parts in equipment_part_assignments?
4. Parts are marked is_available = true?

**Solution:**
```sql
-- Verify chain
SELECT 
    t.id as ticket_id,
    er.id as registry_id,
    er.equipment_catalog_id,
    COUNT(epa.id) as parts_available
FROM service_tickets t
JOIN equipment_registry er ON t.equipment_id = er.id
LEFT JOIN equipment_part_assignments epa ON epa.equipment_catalog_id = er.equipment_catalog_id
WHERE t.id = 'your-ticket-id'
GROUP BY t.id, er.id, er.equipment_catalog_id;
```

### Issue: QR code not working
**Check:**
1. QR code is unique in equipment_registry?
2. QR code URL format correct?
3. Frontend QR scanner configured?

**Solution:**
```sql
-- Verify QR code
SELECT id, qr_code, serial_number, equipment_name
FROM equipment_registry
WHERE qr_code = 'scanned-qr-code';
```

### Issue: Ticket creation fails
**Check:**
1. equipment_id exists in equipment_registry?
2. Required fields provided (customer_name, issue_description)?
3. Foreign key constraints satisfied?

---

## Migration History

### Executed
1. ✅ create_ticket_parts_table.sql
2. ✅ correct_equipment_data.sql  
3. ✅ fix_parts_function.sql
4. ✅ link_equipment_to_catalog.sql

### Pending (Optional - AI Features)
1. ⏳ 012_parts_recommendations.sql (skip ticket_parts creation)
2. ⏳ 013_feedback_system.sql

---

## References

- [ER Diagram](./ER_DIAGRAM.md)
- [API Documentation](./api/)
- [Database Migrations](../database/migrations/)
- [Testing Guide](./TESTING-GUIDE.md)

---

**Last Updated:** 2025-12-19
**Version:** 1.0
**Status:** Production Ready ✅
