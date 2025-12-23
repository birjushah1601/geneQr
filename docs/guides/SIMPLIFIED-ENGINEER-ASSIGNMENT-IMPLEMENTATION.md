# 🎯 Simplified Engineer Assignment System - Implementation Guide

**Date:** November 21, 2025  
**Status:** Phase 1 Complete (Database) → Phase 2 Starting (APIs)  
**Approach:** Loosely coupled, configuration-driven, organization-centric

---

## 📋 Overview

A **simple, extensible engineer assignment system** that routes service tickets based on:
- Equipment manufacturer & category
- Warranty/AMC status
- Engineer levels (L1=Junior, L2=Senior, L3=Expert)
- Organization hierarchy (OEM → Dealer → Distributor → Hospital)

**No AI initially** - pure configuration-based routing that can be enhanced later.

---

## 🗄️ Database Schema

### **1. Engineers Table (Enhanced)**
```sql
engineers (
  id UUID PRIMARY KEY,
  org_id UUID,                  -- Which organization they belong to
  full_name TEXT,
  first_name TEXT,
  last_name TEXT,
  email TEXT,
  phone TEXT,
  engineer_level INT DEFAULT 1, -- NEW: 1=Junior, 2=Senior, 3=Expert
  status TEXT DEFAULT 'available',
  base_location TEXT,           -- City-based for now
  ...
)
```

**Engineer Levels:**
- **L1 (Junior):** Basic repairs, routine maintenance
- **L2 (Senior):** Most repairs, complex diagnostics
- **L3 (Expert):** Specialized equipment, critical systems

### **2. Engineer Equipment Types (NEW)**
```sql
engineer_equipment_types (
  id UUID PRIMARY KEY,
  engineer_id UUID,
  manufacturer_name TEXT,       -- e.g., "Siemens Healthineers"
  equipment_category TEXT,      -- e.g., "MRI", "X-Ray"
  model_pattern TEXT,           -- Optional: "Magnetom%" or NULL for all
  is_certified BOOLEAN,
  certification_number TEXT,
  certification_expiry DATE,
  UNIQUE(engineer_id, manufacturer_name, equipment_category)
)
```

**Purpose:** Defines what each engineer can repair.

### **3. Equipment Service Config (NEW)**
```sql
equipment_service_config (
  id UUID PRIMARY KEY,
  equipment_id UUID UNIQUE,
  
  -- Service hierarchy (ordered by priority)
  primary_service_org_id UUID,      -- Usually manufacturer
  secondary_service_org_id UUID,    -- Usually dealer
  tertiary_service_org_id UUID,     -- Usually distributor
  fallback_service_org_id UUID,     -- Hospital in-house
  
  -- Warranty/AMC (affects priority)
  warranty_provider_org_id UUID,
  warranty_active BOOLEAN,
  warranty_start_date DATE,
  warranty_end_date DATE,
  
  amc_provider_org_id UUID,
  amc_active BOOLEAN,
  amc_start_date DATE,
  amc_end_date DATE,
  amc_contract_number TEXT,
  
  -- Minimum engineer level required
  min_engineer_level INT DEFAULT 1,
  
  service_notes TEXT
)
```

**Purpose:** Per-equipment routing configuration.

---

## 🔧 Assignment Algorithm (Simple)

### **Step 1: Determine Eligible Organizations**

```
IF warranty_active THEN
  Priority Org = warranty_provider_org_id (Tier 1)
ELSE IF amc_active THEN
  Priority Org = amc_provider_org_id (Tier 1)
ELSE
  Priority Org = primary_service_org_id (Tier 1)
END

Fallback tiers:
  Tier 2: secondary_service_org_id
  Tier 3: tertiary_service_org_id
  Tier 4: fallback_service_org_id
```

### **Step 2: Find Matching Engineers**

```sql
SELECT e.*
FROM engineers e
JOIN engineer_equipment_types eet ON e.id = eet.engineer_id
WHERE 
  e.org_id IN (eligible_org_ids)          -- From step 1
  AND e.status = 'available'
  AND e.engineer_level >= min_level       -- From config
  AND eet.manufacturer_name = equipment.manufacturer
  AND eet.equipment_category = equipment.category
ORDER BY 
  ARRAY_POSITION(eligible_org_ids, e.org_id),  -- Tier priority
  e.engineer_level DESC,                         -- Higher level first
  e.id
```

### **Step 3: Return Suggestions**

Frontend shows list of suggested engineers with:
- Priority number (1, 2, 3...)
- Engineer name & level
- Organization
- Location
- Certification status

User manually selects one to assign.

---

## 📊 Seed Data Created

### **Engineers:**

| Organization | Name | Level | Location | Can Service |
|-------------|------|-------|----------|-------------|
| **Siemens** | Rajesh Singh | L3 | Mumbai | Siemens MRI, CT (certified) |
| **Siemens** | Priya Sharma | L2 | Delhi | Siemens X-Ray, Ultrasound |
| **Siemens** | Amit Patel | L2 | Bangalore | Siemens MRI, CT, X-Ray |
| **GE Healthcare** | Vikram Reddy | L3 | Hyderabad | GE CT, MRI, PET-CT (certified) |
| **GE Healthcare** | Sneha Desai | L2 | Mumbai | GE X-Ray, Ultrasound |
| **Philips** | Arun Menon | L3 | Chennai | Philips MRI, CT (certified) |
| **City Medical (Dealer)** | Suresh Gupta | L2 | Delhi | Multi-brand X-Ray, Ultrasound |
| **City Medical (Dealer)** | Rahul Verma | L1 | Delhi | Siemens X-Ray, GE Ultrasound |
| **Apollo Hospital** | Manish Joshi | L2 | Delhi | Multi-brand X-Ray, Ultrasound (in-house) |
| **Apollo Hospital** | Deepak Yadav | L1 | Delhi | Siemens X-Ray, GE Ultrasound |

### **Equipment Configs:**

1. **Siemens MRI** (Warranty active)
   - Priority: Siemens → City Medical (dealer) → Apollo (in-house)
   - Min Level: L2
   
2. **GE X-Ray** (AMC with dealer)
   - Priority: City Medical (dealer) → GE → Apollo
   - Min Level: L1

3. **Philips Ultrasound** (No warranty/AMC)
   - Priority: Philips → City Medical → Apollo
   - Min Level: L1

---

## 🚀 Implementation Plan

### ✅ Phase 1: Database (COMPLETE)
- [x] Create migration file: `003_simplified_engineer_assignment.sql`
- [x] Create seed data: `005_engineer_assignment_data.sql`
- [x] Add helper function: `get_eligible_service_orgs()`

### 📝 Phase 2: Backend APIs (NEXT)

#### **2.1 Engineers Management APIs**
- `GET /api/v1/organizations/{orgId}/engineers` - List org's engineers
- `POST /api/v1/organizations/{orgId}/engineers` - Add engineer
- `PUT /api/v1/engineers/{id}` - Update engineer
- `DELETE /api/v1/engineers/{id}` - Remove engineer

#### **2.2 Engineer Equipment Types APIs**
- `GET /api/v1/engineers/{id}/equipment-types` - List capabilities
- `POST /api/v1/engineers/{id}/equipment-types` - Add capability
- `DELETE /api/v1/engineers/{id}/equipment-types/{typeId}` - Remove

#### **2.3 Assignment APIs**
- `GET /api/v1/service-tickets/{id}/suggested-engineers` - Get suggestions
- `POST /api/v1/service-tickets/{id}/assign` - Manual assignment
- `GET /api/v1/engineers/{id}/workload` - Current assignments

### 🎨 Phase 3: Frontend (AFTER APIs)

#### **3.1 Engineers Management Page**
- List view with filters (by level, location, status)
- Add/Edit engineer form
- Equipment capabilities management

#### **3.2 Service Ticket Assignment UI**
- Show suggested engineers (sorted by priority)
- Display badges (priority, level, certified)
- One-click assignment
- Workload indicator

---

## 🔌 API Examples

### **Get Suggested Engineers**

**Request:**
```http
GET /api/v1/service-tickets/ticket-123/suggested-engineers
```

**Response:**
```json
{
  "ticket_id": "ticket-123",
  "equipment": {
    "id": "eq-001",
    "name": "Magnetom Vida 3T",
    "manufacturer": "Siemens Healthineers",
    "category": "MRI"
  },
  "suggested_engineers": [
    {
      "id": "eng-001",
      "name": "Rajesh Singh",
      "level": 3,
      "organization": "Siemens Healthineers",
      "organization_id": "org-siemens",
      "base_location": "Mumbai",
      "is_certified": true,
      "priority": 1,
      "reason": "warranty_coverage"
    },
    {
      "id": "eng-003",
      "name": "Amit Patel",
      "level": 2,
      "organization": "Siemens Healthineers",
      "organization_id": "org-siemens",
      "base_location": "Bangalore",
      "is_certified": true,
      "priority": 2,
      "reason": "warranty_coverage"
    },
    {
      "id": "eng-dealer-001",
      "name": "Suresh Gupta",
      "level": 2,
      "organization": "City Medical Equipment",
      "organization_id": "org-dealer",
      "base_location": "Delhi",
      "is_certified": false,
      "priority": 3,
      "reason": "authorized_dealer"
    }
  ]
}
```

### **Assign Engineer**

**Request:**
```http
POST /api/v1/service-tickets/ticket-123/assign
Content-Type: application/json

{
  "engineer_id": "eng-001",
  "notes": "Assigned due to warranty coverage and proximity"
}
```

**Response:**
```json
{
  "message": "Engineer assigned successfully",
  "ticket_id": "ticket-123",
  "assigned_engineer": {
    "id": "eng-001",
    "name": "Rajesh Singh"
  },
  "assigned_at": "2025-11-21T10:30:00Z"
}
```

---

## 🧪 Testing Scenarios

### **Scenario 1: Warranty Equipment**
- Equipment: Siemens MRI (warranty active)
- Expected: Siemens engineer (L2+) suggested first
- Test: Create ticket → Check suggestions → Assign Rajesh (L3)

### **Scenario 2: AMC Equipment**
- Equipment: GE X-Ray (AMC with dealer)
- Expected: Dealer engineer suggested first
- Test: Create ticket → Check suggestions → Assign Suresh (L2 dealer)

### **Scenario 3: No Coverage**
- Equipment: Philips Ultrasound (no warranty/AMC)
- Expected: Manufacturer → Dealer → Hospital order
- Test: Create ticket → Check all tiers available

### **Scenario 4: Level Requirements**
- Equipment: Siemens MRI (requires L2+)
- Expected: L1 engineers filtered out
- Test: Verify only L2 and L3 engineers appear

---

## 🔮 Future Enhancements (Not Now)

Once the basic system works, we can add:

1. **Availability tracking** (calendar, working hours)
2. **Real-time location** (GPS from mobile app)
3. **Workload balancing** (current ticket count)
4. **Travel time calculation** (distance-based ETA)
5. **Automatic assignment** (skip manual selection)
6. **AI suggestions** (historical performance, success rate)
7. **Skills matrix** (more granular than manufacturer + category)
8. **Multi-step workflows** (diagnostic → repair → testing)

---

## 📁 Files Created

```
database/
├── migrations/
│   └── 003_simplified_engineer_assignment.sql  ✅ Created
└── seed/
    └── 005_engineer_assignment_data.sql        ✅ Created

docs/
└── SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md  ✅ This file
```

---

## ✅ Next Steps

1. **Run migration:**
   ```bash
   psql -U postgres -d medplatform -f database/migrations/003_simplified_engineer_assignment.sql
   ```

2. **Load seed data:**
   ```bash
   psql -U postgres -d medplatform -f database/seed/005_engineer_assignment_data.sql
   ```

3. **Verify data:**
   ```sql
   SELECT COUNT(*) FROM engineers;
   SELECT COUNT(*) FROM engineer_equipment_types;
   SELECT COUNT(*) FROM equipment_service_config;
   ```

4. **Start building APIs** (Phase 2)

---

**Simple, clean, extensible!** 🎯
