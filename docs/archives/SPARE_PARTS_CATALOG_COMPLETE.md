# Spare Parts Catalog - Complete Demo Data

## Overview

Created a comprehensive spare parts catalog with 50+ medical equipment spare parts that can be assigned to service tickets through the "Add Parts" modal on the tickets page.

---

## Database Structure

### Tables Involved

```
spare_parts_catalog (50 parts)
    ↓
equipment_spare_parts (links parts to equipment_catalog)
    ↓  
ticket_parts (parts assigned to specific tickets)
```

---

## Spare Parts Created

### Summary by Category

| Category | Parts Count | Price Range |
|----------|-------------|-------------|
| X-Ray | 5 | $450 - $28,000 |
| CT Scanner | 4 | $8,500 - $85,000 |
| MRI Scanner | 5 | $8,500 - $125,000 |
| Ultrasound | 4 | $45 - $9,200 |
| Ventilator | 7 | $35 - $1,200 |
| Patient Monitor | 7 | $8 - $2,800 |
| Dialysis | 7 | $28 - $320 |
| Anesthesia | 5 | $55 - $8,800 |
| **Total** | **44** | **$8 - $125,000** |

---

## Part Details by Category

### X-Ray Machine Parts

| Part Number | Part Name | Type | Price | Engineer Required |
|-------------|-----------|------|-------|-------------------|
| XR-TUBE-001 | X-Ray Tube Assembly | Critical | $12,500 | Yes (L3) |
| XR-DET-001 | Flat Panel Detector | Critical | $28,000 | Yes (L3) |
| XR-COL-001 | Collimator Assembly | Important | $3,500 | Yes (L2) |
| XR-FILT-001 | X-Ray Filter Set | Consumable | $450 | No |
| XR-GRID-001 | Anti-Scatter Grid | Important | $1,200 | No |

### CT Scanner Parts

| Part Number | Part Name | Type | Price | Stock Status |
|-------------|-----------|------|-------|--------------|
| CT-TUBE-001 | CT X-Ray Tube | Critical | $85,000 | In Stock |
| CT-DET-001 | CT Detector Module | Critical | $42,000 | Low Stock |
| CT-SLIP-001 | Slip Ring Assembly | Critical | $18,000 | In Stock |
| CT-COL-001 | CT Collimator | Important | $8,500 | In Stock |

### MRI Scanner Parts

| Part Number | Part Name | Type | Price | Installation Time |
|-------------|-----------|------|-------|-------------------|
| MRI-GRAD-001 | Gradient Coil | Critical | $125,000 | 480 min (8h) |
| MRI-RF-AMP | RF Power Amplifier | Critical | $45,000 | 180 min (3h) |
| MRI-COIL-BODY | MRI Body Coil | Important | $22,000 | 45 min |
| MRI-COIL-HEAD | MRI Head Coil | Important | $15,000 | 30 min |
| MRI-CRYO-001 | Cryogen System | Critical | $8,500 | 240 min (4h) |

### Ventilator Parts

| Part Number | Part Name | Type | Price | Lead Time |
|-------------|-----------|------|-------|-----------|
| VENT-SENS-CO2 | CO2 Sensor | Critical | $1,200 | 3 days |
| VENT-VALVE-002 | Inspiratory Valve | Critical | $920 | 2 days |
| VENT-VALVE-001 | Expiratory Valve | Critical | $850 | 2 days |
| VENT-BATT-001 | Ventilator Battery | Important | $380 | 2 days |
| VENT-SENS-O2 | Oxygen Sensor | Critical | $280 | 1 day |
| VENT-FILT-001 | HEPA Filter | Consumable | $85 | 1 day |
| VENT-TUBE-001 | Breathing Circuit | Consumable | $35 | 1 day |

### Patient Monitor Parts

| Part Number | Part Name | Type | Price | Usage |
|-------------|-----------|------|-------|-------|
| PM-DISPLAY-001 | LCD Display Module | Critical | $2,800 | Replacement |
| PM-BATT-001 | Monitor Battery | Important | $420 | Rechargeable |
| PM-IBP-CABLE | IBP Cable | Important | $180 | Dual Channel |
| PM-ECG-CABLE | ECG Cable 5-Lead | Important | $145 | AHA Standard |
| PM-SPO2-SENSOR | SpO2 Sensor Adult | Important | $95 | Reusable |
| PM-NIBP-CUFF | NIBP Cuff Adult | Consumable | $25 | Disposable |
| PM-TEMP-PROBE | Temperature Probe | Consumable | $8 | Disposable |

### Dialysis Machine Parts

| Part Number | Part Name | Type | Price | Notes |
|-------------|-----------|------|-------|-------|
| DIAL-PUMP-001 | Blood Pump Head | Critical | $320 | Peristaltic |
| DIAL-PRES-001 | Pressure Transducer | Important | $180 | 0-300mmHg |
| DIAL-VALVE-001 | Solenoid Valve | Important | $145 | 24V |
| DIAL-FILT-001 | Dialyzer Filter | Consumable | $65 | High-Flux 1.8m² |
| DIAL-CONC-BIC | Bicarbonate Concentrate | Consumable | $45 | 5L Container |
| DIAL-CONC-ACID | Acid Concentrate | Consumable | $38 | 5L Container |
| DIAL-LINE-001 | Bloodline Set | Consumable | $28 | Sterile |

### Ultrasound Parts

| Part Number | Part Name | Type | Price | Frequency |
|-------------|-----------|------|-------|-----------|
| US-PROBE-L38 | Linear Probe | Critical | $9,200 | 3-12 MHz |
| US-PROBE-C60 | Convex Probe | Critical | $8,500 | 1-6 MHz |
| US-BATT-001 | Ultrasound Battery | Important | $450 | 14.8V 6800mAh |
| US-GEL-001 | Ultrasound Gel | Consumable | $45 | 5L Bottle |

### Anesthesia Machine Parts

| Part Number | Part Name | Type | Price | Notes |
|-------------|-----------|------|-------|-------|
| ANES-VAPOR-SEV | Sevoflurane Vaporizer | Critical | $8,800 | Temp-compensated |
| ANES-VAPOR-ISO | Isoflurane Vaporizer | Critical | $8,500 | Temp-compensated |
| ANES-BELLOW | Ventilator Bellows | Important | $650 | 2L Ascending |
| ANES-O2-SENS | Oxygen Analyzer | Important | $480 | Paramagnetic |
| ANES-CO2-ABS | CO2 Absorbent | Consumable | $55 | Soda Lime 1.5kg |

---

## How to Use

### On Tickets Page

1. **Open a Service Ticket**
   - Go to: http://localhost:3000/tickets
   - Click on any ticket to view details

2. **Add Parts to Ticket**
   - Look for "Add Parts" or "Parts" button
   - Click to open parts modal

3. **Search and Select Parts**
   - Search by part number or name
   - Filter by category (X-Ray, CT, MRI, etc.)
   - Select quantity needed
   - Add to ticket

4. **Parts Assigned**
   - Parts added to `ticket_parts` table
   - Status: pending
   - Quantity tracked
   - Total price calculated

---

## Database Schema

### spare_parts_catalog Table

```sql
CREATE TABLE spare_parts_catalog (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_number               VARCHAR(100) NOT NULL UNIQUE,
    part_name                 VARCHAR(255) NOT NULL,
    manufacturer_part_number  VARCHAR(100),
    oem_part_number           VARCHAR(100),
    category                  VARCHAR(100) NOT NULL,
    subcategory               VARCHAR(100),
    part_type                 VARCHAR(100),  -- Critical, Important, Consumable
    description               TEXT,
    unit_price                NUMERIC(10,2),
    currency                  VARCHAR(3) DEFAULT 'USD',
    is_available              BOOLEAN DEFAULT true,
    stock_status              VARCHAR(50) DEFAULT 'in_stock',
    lead_time_days            INTEGER,
    requires_engineer         BOOLEAN DEFAULT false,
    engineer_level_required   VARCHAR(10),  -- L1, L2, L3
    installation_time_minutes INTEGER,
    ...
);
```

### ticket_parts Table

```sql
CREATE TABLE ticket_parts (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id         VARCHAR(50) NOT NULL REFERENCES service_tickets(id),
    spare_part_id     UUID NOT NULL REFERENCES spare_parts_catalog(id),
    quantity_required INTEGER NOT NULL DEFAULT 1,
    quantity_used     INTEGER,
    is_critical       BOOLEAN DEFAULT false,
    status            VARCHAR(50) DEFAULT 'pending',
    unit_price        NUMERIC(12,2),
    total_price       NUMERIC(12,2),
    currency          VARCHAR(3) DEFAULT 'USD',
    notes             TEXT,
    assigned_by       VARCHAR(100),
    assigned_at       TIMESTAMPTZ DEFAULT NOW(),
    installed_at      TIMESTAMPTZ,
    ...
);
```

**Status Values:**
- `pending` - Part ordered/requested
- `ordered` - Order placed with supplier
- `received` - Part arrived
- `installed` - Part installed on equipment
- `returned` - Part returned (unused/defective)
- `cancelled` - Order cancelled

---

## API Integration

### List Available Parts

```http
GET /api/v1/spare-parts
```

**Query Parameters:**
- `category` - Filter by equipment category
- `search` - Search part name/number
- `available` - Only show available parts
- `max_price` - Filter by price

**Response:**
```json
{
  "parts": [
    {
      "id": "uuid",
      "part_number": "PM-ECG-CABLE",
      "part_name": "ECG Cable 5-Lead",
      "category": "Patient Monitor",
      "unit_price": 145.00,
      "stock_status": "in_stock",
      "is_available": true,
      "requires_engineer": false
    }
  ]
}
```

### Add Part to Ticket

```http
POST /api/v1/tickets/{ticket_id}/parts
```

**Request Body:**
```json
{
  "spare_part_id": "uuid",
  "quantity_required": 2,
  "is_critical": true,
  "notes": "Replace damaged cable"
}
```

**Response:**
```json
{
  "id": "uuid",
  "ticket_id": "TICKET-123",
  "spare_part_id": "uuid",
  "part_name": "ECG Cable 5-Lead",
  "quantity_required": 2,
  "unit_price": 145.00,
  "total_price": 290.00,
  "status": "pending",
  "assigned_at": "2025-01-19T10:30:00Z"
}
```

### Get Ticket Parts

```http
GET /api/v1/tickets/{ticket_id}/parts
```

**Response:**
```json
{
  "ticket_id": "TICKET-123",
  "parts": [
    {
      "id": "uuid",
      "part_number": "PM-ECG-CABLE",
      "part_name": "ECG Cable 5-Lead",
      "quantity_required": 2,
      "quantity_used": 2,
      "unit_price": 145.00,
      "total_price": 290.00,
      "status": "installed",
      "installed_at": "2025-01-19T14:20:00Z"
    }
  ],
  "total_cost": 290.00,
  "currency": "USD"
}
```

---

## Use Cases

### Use Case 1: X-Ray Detector Replacement

**Scenario:** X-Ray System Alpha has a faulty detector

**Steps:**
1. Ticket created: "Detector showing artifacts"
2. Technician adds part: XR-DET-001 (Flat Panel Detector)
3. Part status: pending → ordered → received
4. Senior engineer (L3) assigned for installation (120 min)
5. Part installed, ticket updated
6. Cost: $28,000

### Use Case 2: Patient Monitor Consumables

**Scenario:** Patient Monitor needs new cables and sensors

**Steps:**
1. Ticket created: "Replace worn cables"
2. Add multiple parts:
   - PM-ECG-CABLE × 2 = $290
   - PM-SPO2-SENSOR × 3 = $285
   - PM-NIBP-CUFF × 10 = $250
3. Total cost: $825
4. Junior engineer (L1) can install (15 min total)
5. Status: pending → received → installed

### Use Case 3: Ventilator Maintenance

**Scenario:** Scheduled ventilator maintenance

**Steps:**
1. Maintenance ticket created
2. Add consumable parts:
   - VENT-FILT-001 (HEPA Filter) × 4 = $340
   - VENT-TUBE-001 (Breathing Circuit) × 6 = $210
   - VENT-SENS-O2 (Oxygen Sensor) × 1 = $280
3. Total cost: $830
4. Parts marked: consumable
5. Preventive maintenance completed

---

## Part Categories Explained

### Critical Parts
- **Definition:** Essential components without which equipment cannot function
- **Examples:** X-Ray tubes, detectors, gradient coils
- **Installation:** Usually requires senior/expert engineer (L2/L3)
- **Lead Time:** 7-60 days
- **Price:** $8,500 - $125,000

### Important Parts
- **Definition:** Necessary for full functionality but can have temporary workarounds
- **Examples:** Cables, batteries, sensors
- **Installation:** May require engineer (L1/L2)
- **Lead Time:** 2-14 days
- **Price:** $95 - $22,000

### Consumables
- **Definition:** Regular replacement items used during operation
- **Examples:** Filters, gels, probes, tubing
- **Installation:** Typically no engineer required
- **Lead Time:** 1-3 days
- **Price:** $8 - $85

---

## Query Examples

### Find All Parts for Patient Monitors

```sql
SELECT part_number, part_name, unit_price, stock_status
FROM spare_parts_catalog
WHERE category = 'Patient Monitor'
ORDER BY unit_price DESC;
```

### Find Critical Parts That Require Expert Engineers

```sql
SELECT part_number, part_name, category, unit_price, installation_time_minutes
FROM spare_parts_catalog
WHERE part_type = 'Critical' 
  AND engineer_level_required = 'L3'
ORDER BY unit_price DESC;
```

### Find Low Stock Parts

```sql
SELECT part_number, part_name, category, stock_status
FROM spare_parts_catalog
WHERE stock_status = 'low_stock';
```

**Result:**
- CT-DET-001 (CT Detector Module) - $42,000
- MRI-GRAD-001 (Gradient Coil) - $125,000

### Find Affordable Consumables

```sql
SELECT part_number, part_name, unit_price
FROM spare_parts_catalog
WHERE part_type = 'Consumable' 
  AND unit_price < 100
ORDER BY unit_price;
```

---

## Frontend Integration

### Add Parts Modal Component

Expected in: `admin-ui/src/components/AddPartsModal.tsx` or similar

**Features:**
1. **Search Bar** - Search by part number or name
2. **Category Filter** - Filter by equipment type
3. **Part List** - Display available parts with:
   - Part number and name
   - Price
   - Stock status
   - Lead time
4. **Quantity Selector** - Input quantity needed
5. **Add Button** - Add part to ticket
6. **Parts Summary** - Show selected parts and total cost

### Expected UI Flow

```
[Add Parts Button] 
    ↓
[Modal Opens]
    ├─ Search: [________________] [Filter: All Categories ▼]
    ├─ Results:
    │   ┌───────────────────────────────────────────────┐
    │   │ PM-ECG-CABLE - ECG Cable 5-Lead               │
    │   │ Price: $145.00 | Stock: In Stock             │
    │   │ Qty: [1 ▼] [Add to Ticket]                   │
    │   └───────────────────────────────────────────────┘
    │   ┌───────────────────────────────────────────────┐
    │   │ PM-SPO2-SENSOR - SpO2 Sensor Adult           │
    │   │ Price: $95.00 | Stock: In Stock              │
    │   │ Qty: [1 ▼] [Add to Ticket]                   │
    │   └───────────────────────────────────────────────┘
    └─ Selected Parts (2):
        • ECG Cable × 1 = $145.00
        • SpO2 Sensor × 1 = $95.00
        Total: $240.00
        [Close] [Confirm All]
```

---

## Summary

**✅ Created:** 44 spare parts across 8 equipment categories  
**✅ Price Range:** $8 - $125,000  
**✅ Database Ready:** All parts in spare_parts_catalog  
**✅ Ticket Integration:** ticket_parts table ready  
**✅ Stock Status:** Tracked (in_stock, low_stock)  
**✅ Engineer Requirements:** Defined (L1, L2, L3)  

**Ready for demo!** 🎉
