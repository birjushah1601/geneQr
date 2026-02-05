# ðŸ”§ Engineer Management Architecture - Comprehensive Design

**Document Type:** Technical Design  
**Date:** October 11, 2025  
**Status:** Final Design before Implementation

---

## ðŸ“‹ Overview

This document covers the **multi-entity engineer management system** where:
- **Manufacturers** have engineers (per facility)
- **Sub-Sub-sub_sub_SUB_DEALERs** have engineers (per location)
- **Channel Partners** have service engineers (per service center)
- **Hospitals/Clients** have in-house BME teams
- **Independent Service Providers** have engineer teams

**Critical Routing Logic:** If manufacturer/Sub-sub_SUB_DEALER/Channel Partner don't have available engineers, the system should route service requests to the client's in-house engineers as a fallback.

---

## 1. Entity Model

### 1.1 Engineer Entity

```typescript
interface Engineer {
  id: UUID;
  
  // Identity
  employee_id: string; // Organization's internal ID
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  whatsapp_number?: string;
  
  // Employment
  org_id: UUID; // Which organization employs this engineer
  org_type: OrgType; // manufacturer, Sub-sub_SUB_DEALER, hospital, service_provider
  employment_type: 'full_time' | 'part_time' | 'contract' | 'freelance';
  joining_date: Date;
  
  // Location Assignment
  primary_facility_id?: UUID; // Home facility
  assigned_facilities?: UUID[]; // Can work at these facilities
  mobile_engineer: boolean; // Can travel to customer sites
  
  // Coverage
  coverage_radius_km?: number; // For mobile engineers
  coverage_pincodes?: string[];
  coverage_cities?: string[];
  coverage_states?: string[];
  
  // Availability
  status: 'available' | 'on_job' | 'on_leave' | 'inactive';
  current_location?: Point; // Real-time GPS (if mobile)
  
  // Workload
  active_tickets: number;
  max_daily_tickets: number;
  
  // Schedule
  working_hours: OperatingHours;
  on_call_24x7: boolean;
  
  // Performance
  total_tickets_resolved: number;
  avg_resolution_time_hours: number;
  customer_rating: number; // 0-5
  first_time_fix_rate: number; // Percentage
  
  // Contact Preferences
  preferred_contact_method: 'phone' | 'whatsapp' | 'email';
  language_preferences: string[];
  
  // System
  metadata: JSONB;
  created_at: DateTime;
  updated_at: DateTime;
}
```

---

### 1.2 Engineer Skills & Certifications

```typescript
interface EngineerSkill {
  id: UUID;
  engineer_id: UUID;
  
  // Skill Definition
  skill_type: SkillType;
  
  // Equipment Type Skills
  equipment_category?: string; // "Imaging", "Patient Monitoring", "Lab"
  equipment_type?: string; // "CT Scanner", "X-Ray", "MRI"
  equipment_models?: string[]; // ["Siemens SOMATOM", "GE Revolution"]
  
  // Manufacturer-Specific Skills
  manufacturer_id?: UUID; // Certified for which manufacturer
  manufacturer_name?: string; // "Siemens", "GE", "Philips"
  manufacturer_authorized: boolean; // Official OEM training?
  
  // Skill Level
  proficiency_level: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  
  // Certification
  certification_name?: string; // "Siemens CT Advanced Certification"
  certification_number?: string;
  certification_authority?: string; // "Siemens Training Center"
  certified_date?: Date;
  expiry_date?: Date;
  certificate_document_url?: string;
  
  // Capabilities
  can_install: boolean;
  can_calibrate: boolean;
  can_repair: boolean;
  can_train_users: boolean;
  
  // Experience
  years_of_experience: number;
  tickets_resolved_for_this_skill: number;
  
  // Verification
  verified: boolean;
  verified_by?: string;
  verified_date?: Date;
  
  metadata: JSONB;
  created_at: DateTime;
  updated_at: DateTime;
}

enum SkillType {
  EQUIPMENT_CATEGORY = 'equipment_category', // Broad: All imaging equipment
  EQUIPMENT_TYPE = 'equipment_type', // Specific: CT Scanners
  EQUIPMENT_MODEL = 'equipment_model', // Very specific: Siemens SOMATOM
  MANUFACTURER_GENERAL = 'manufacturer_general', // All Siemens equipment
  SERVICE_TYPE = 'service_type' // Installation, Calibration, etc.
}
```

---

### 1.3 Engineer Availability & Schedule

```typescript
interface EngineerAvailability {
  id: UUID;
  engineer_id: UUID;
  
  // Date Range
  date: Date;
  available: boolean;
  
  // Reason for Unavailability
  reason?: 'on_leave' | 'on_training' | 'sick' | 'on_site' | 'other';
  notes?: string;
  
  // Time Slots (for partial day availability)
  available_slots?: TimeSlot[];
  blocked_slots?: TimeSlot[];
  
  created_at: DateTime;
  updated_at: DateTime;
}

interface TimeSlot {
  start_time: string; // "09:00"
  end_time: string; // "12:00"
}
```

---

### 1.4 Engineer Work Assignment

```typescript
interface EngineerAssignment {
  id: UUID;
  
  // Assignment
  engineer_id: UUID;
  ticket_id: UUID;
  equipment_id: UUID;
  
  // Assignment Details
  assigned_by: UUID; // User who assigned
  assigned_at: DateTime;
  assignment_type: 'auto' | 'manual'; // Auto-routed or manually assigned
  
  // Status
  status: 'assigned' | 'accepted' | 'en_route' | 'on_site' | 'completed' | 'cancelled';
  
  // Timeline
  accepted_at?: DateTime;
  en_route_at?: DateTime;
  reached_site_at?: DateTime;
  work_started_at?: DateTime;
  work_completed_at?: DateTime;
  
  // Location Tracking
  engineer_start_location?: Point;
  customer_location: Point;
  travel_distance_km?: number;
  estimated_arrival?: DateTime;
  actual_arrival?: DateTime;
  
  // Work Details
  issue_description: string;
  diagnosis?: string;
  actions_taken?: string;
  parts_used?: PartUsage[];
  
  // Customer Feedback
  customer_signature?: string; // Base64 image
  customer_rating?: number; // 1-5
  customer_feedback?: string;
  
  // Photos
  before_photos?: string[];
  after_photos?: string[];
  
  metadata: JSONB;
  created_at: DateTime;
  updated_at: DateTime;
}

interface PartUsage {
  part_id: UUID;
  part_name: string;
  quantity: number;
  serial_numbers?: string[];
}
```

---

## 2. Service Request Routing Logic

### 2.1 Multi-Tier Routing Strategy

```typescript
interface ServiceRoutingConfig {
  ticket_id: UUID;
  equipment_id: UUID;
  
  // Equipment Details
  manufacturer_id: UUID;
  equipment_type: string;
  equipment_model: string;
  
  // Location
  customer_location: Point;
  customer_pincode: string;
  customer_city: string;
  customer_state: string;
  
  // Priority
  priority: 'critical' | 'high' | 'medium' | 'low';
  sla_hours: number;
  
  // Routing Preferences
  routing_tiers: RoutingTier[];
}

interface RoutingTier {
  tier_number: number;
  tier_name: string; // "OEM Service", "Sub-sub_SUB_DEALER Service", "Client In-House"
  org_ids: UUID[]; // Organizations to check in this tier
  auto_route: boolean; // Auto-assign or wait for manual?
  timeout_minutes: number; // Wait time before moving to next tier
}
```

---

### 2.2 Routing Algorithm

```typescript
/**
 * TIER-BASED SERVICE REQUEST ROUTING
 * 
 * Priority Order:
 * 1. OEM/Manufacturer's Own Engineers (if available)
 * 2. Authorized Sub-sub_SUB_DEALER's Engineers (who sold the equipment)
 * 3. Channel Partner's Service Centers (in the region)
 * 4. Independent Service Providers (with manufacturer authorization)
 * 5. Client's In-House Engineers (BME Team) - FALLBACK
 * 6. Any Available Engineer (with required skills)
 */

async function routeServiceRequest(ticket: ServiceTicket): Promise<RoutingResult> {
  const equipment = await getEquipment(ticket.equipment_id);
  const customer = await getOrganization(ticket.customer_org_id);
  
  // Build routing configuration
  const config: ServiceRoutingConfig = {
    ticket_id: ticket.id,
    equipment_id: equipment.id,
    manufacturer_id: equipment.manufacturer_id,
    equipment_type: equipment.equipment_type,
    equipment_model: equipment.model,
    customer_location: customer.location,
    priority: ticket.priority,
    sla_hours: ticket.sla_hours,
    routing_tiers: []
  };
  
  // TIER 1: Manufacturer's Engineers (if manufacturer has service team)
  const manufacturerEngineers = await findEligibleEngineers({
    org_id: equipment.manufacturer_id,
    org_type: 'manufacturer',
    equipment_id: equipment.id,
    location: customer.location,
    max_distance_km: 100,
    available: true,
    skills_match: true
  });
  
  if (manufacturerEngineers.length > 0) {
    config.routing_tiers.push({
      tier_number: 1,
      tier_name: 'OEM Service (Manufacturer)',
      org_ids: [equipment.manufacturer_id],
      auto_route: true,
      timeout_minutes: 30
    });
  }
  
  // TIER 2: Sub-sub_SUB_DEALER's Engineers (who sold this equipment)
  if (equipment.sold_by_sub_sub_Sub-sub_SUB_DEALER_id) {
    const Sub-sub_SUB_DEALEREngineers = await findEligibleEngineers({
      org_id: equipment.sold_by_sub_sub_Sub-sub_SUB_DEALER_id,
      org_type: 'Sub-sub_SUB_DEALER',
      equipment_id: equipment.id,
      location: customer.location,
      max_distance_km: 50,
      available: true,
      skills_match: true
    });
    
    if (Sub-sub_SUB_DEALEREngineers.length > 0) {
      config.routing_tiers.push({
        tier_number: 2,
        tier_name: 'Sub-sub_SUB_DEALER Service',
        org_ids: [equipment.sold_by_sub_sub_Sub-sub_SUB_DEALER_id],
        auto_route: true,
        timeout_minutes: 20
      });
    }
  }
  
  // TIER 3: Channel Partner's Service Centers
  const Channel Partners = await getChannel PartnersForManufacturer({
    manufacturer_id: equipment.manufacturer_id,
    covers_location: customer.location
  });
  
  const Channel PartnerEngineers = await findEligibleEngineers({
    org_ids: Channel Partners.map(d => d.id),
    org_type: 'Channel Partner',
    equipment_id: equipment.id,
    location: customer.location,
    max_distance_km: 100,
    available: true,
    skills_match: true
  });
  
  if (Channel PartnerEngineers.length > 0) {
    config.routing_tiers.push({
      tier_number: 3,
      tier_name: 'Channel Partner Service Centers',
      org_ids: Channel Partners.map(d => d.id),
      auto_route: true,
      timeout_minutes: 30
    });
  }
  
  // TIER 4: Authorized Service Providers
  const serviceProviders = await getAuthorizedServiceProviders({
    manufacturer_id: equipment.manufacturer_id,
    covers_location: customer.location
  });
  
  const serviceProviderEngineers = await findEligibleEngineers({
    org_ids: serviceProviders.map(sp => sp.id),
    org_type: 'service_provider',
    equipment_id: equipment.id,
    location: customer.location,
    max_distance_km: 150,
    available: true,
    skills_match: true
  });
  
  if (serviceProviderEngineers.length > 0) {
    config.routing_tiers.push({
      tier_number: 4,
      tier_name: 'Authorized Service Providers',
      org_ids: serviceProviders.map(sp => sp.id),
      auto_route: true,
      timeout_minutes: 45
    });
  }
  
  // TIER 5: CLIENT'S IN-HOUSE ENGINEERS (CRITICAL FALLBACK)
  const clientEngineers = await findEligibleEngineers({
    org_id: customer.id,
    org_type: 'hospital', // or clinic, lab, etc.
    equipment_id: equipment.id,
    location: customer.location,
    available: true,
    skills_match: true // Only if they have the skills
  });
  
  if (clientEngineers.length > 0) {
    config.routing_tiers.push({
      tier_number: 5,
      tier_name: 'In-House BME Team (Client)',
      org_ids: [customer.id],
      auto_route: true,
      timeout_minutes: 15
    });
  }
  
  // TIER 6: Any Available Engineer with Skills (Emergency)
  config.routing_tiers.push({
    tier_number: 6,
    tier_name: 'Any Available Engineer',
    org_ids: [],
    auto_route: false, // Manual approval required
    timeout_minutes: 0
  });
  
  // Execute routing through tiers
  return await executeRoutingTiers(config);
}

/**
 * Find eligible engineers based on criteria
 */
async function findEligibleEngineers(criteria: {
  org_id?: UUID;
  org_ids?: UUID[];
  org_type?: OrgType;
  equipment_id: UUID;
  location: Point;
  max_distance_km?: number;
  available: boolean;
  skills_match: boolean;
}): Promise<Engineer[]> {
  
  const equipment = await getEquipment(criteria.equipment_id);
  
  // Build SQL query
  let query = `
    SELECT DISTINCT e.*
    FROM engineers e
    WHERE e.status = 'available'
  `;
  
  // Filter by organization
  if (criteria.org_id) {
    query += ` AND e.org_id = '${criteria.org_id}'`;
  }
  if (criteria.org_ids && criteria.org_ids.length > 0) {
    query += ` AND e.org_id IN ('${criteria.org_ids.join("','")}')`;
  }
  if (criteria.org_type) {
    query += ` AND e.org_type = '${criteria.org_type}'`;
  }
  
  // Filter by location (if mobile engineer)
  if (criteria.max_distance_km) {
    query += `
      AND (
        e.mobile_engineer = true
        AND ST_DWithin(
          e.current_location::geography,
          '${criteria.location}'::geography,
          ${criteria.max_distance_km * 1000}
        )
      )
    `;
  }
  
  // Filter by skills
  if (criteria.skills_match) {
    query += `
      AND EXISTS (
        SELECT 1 FROM engineer_skills es
        WHERE es.engineer_id = e.id
        AND (
          -- Match by equipment type
          (es.skill_type = 'equipment_type' AND es.equipment_type = '${equipment.equipment_type}')
          OR
          -- Match by equipment model
          (es.skill_type = 'equipment_model' AND '${equipment.model}' = ANY(es.equipment_models))
          OR
          -- Match by manufacturer
          (es.skill_type = 'manufacturer_general' AND es.manufacturer_id = '${equipment.manufacturer_id}')
          OR
          -- Match by category
          (es.skill_type = 'equipment_category' AND es.equipment_category = '${equipment.category}')
        )
        AND (es.expiry_date IS NULL OR es.expiry_date > NOW())
      )
    `;
  }
  
  // Check workload
  query += ` AND e.active_tickets < e.max_daily_tickets`;
  
  // Order by best match
  query += `
    ORDER BY 
      e.customer_rating DESC,
      e.first_time_fix_rate DESC,
      ST_Distance(e.current_location::geography, '${criteria.location}'::geography) ASC
    LIMIT 10
  `;
  
  return await db.query(query);
}
```

---

## 3. Real-World Scenarios

### 3.1 Scenario: Manufacturer with Multi-Location Engineers

**Siemens Healthineers:**

```
Organization: Siemens Healthineers India Ltd
â”œâ”€â”€ Chennai Service Hub (Facility)
â”‚   â”œâ”€â”€ Ramesh Kumar (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: CT Scanner, MRI (Siemens only)
â”‚   â”‚   â”œâ”€â”€ Certifications: Siemens CT Advanced, MRI Expert
â”‚   â”‚   â”œâ”€â”€ Coverage: Tamil Nadu, Kerala, Karnataka
â”‚   â”‚   â”œâ”€â”€ Mobile: Yes (200 km radius)
â”‚   â”‚   â””â”€â”€ Active Tickets: 2/5
â”‚   â”œâ”€â”€ Suresh M (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: X-Ray, Ultrasound (Siemens only)
â”‚   â”‚   â”œâ”€â”€ Coverage: Tamil Nadu
â”‚   â”‚   â””â”€â”€ Mobile: Yes (150 km radius)
â”‚   â””â”€â”€ Priya S (Engineer)
â”‚       â”œâ”€â”€ Skills: All Imaging Equipment
â”‚       â”œâ”€â”€ Certifications: Siemens Master Technician
â”‚       â””â”€â”€ Mobile: Yes (300 km radius)
â”‚
â”œâ”€â”€ Mumbai Service Hub (Facility)
â”‚   â”œâ”€â”€ 5 Engineers covering West India
â”‚   â””â”€â”€ Skills: All Siemens equipment
â”‚
â””â”€â”€ Delhi Service Hub (Facility)
    â”œâ”€â”€ 8 Engineers covering North India
    â””â”€â”€ Skills: All Siemens equipment
```

**When Apollo Chennai raises a ticket for Siemens CT Scanner:**
1. System checks: Equipment manufacturer = Siemens
2. Finds Siemens Chennai Service Hub (same city)
3. Filters engineers: CT Scanner skills + available + Chennai location
4. **Result:** Ramesh Kumar (CT expert, 5-star rating, only 2 active tickets)
5. Auto-assigns to Ramesh Kumar

---

### 3.2 Scenario: Sub-sub_SUB_DEALER with Service Team

**City Medical Equipment Co. (Delhi):**

```
Organization: City Medical Equipment Co.
â”œâ”€â”€ Main Showroom (Facility)
â”œâ”€â”€ Service Center (Facility)
â”‚   â”œâ”€â”€ Amit Sharma (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: Multi-brand (Siemens, GE, Philips CT/MRI)
â”‚   â”‚   â”œâ”€â”€ Certifications: Siemens CT, GE CT, Philips Ultrasound
â”‚   â”‚   â”œâ”€â”€ Coverage: Delhi NCR
â”‚   â”‚   â””â”€â”€ Mobile: Yes (50 km)
â”‚   â”œâ”€â”€ Rajesh Verma (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: X-Ray, Ultrasound (all brands)
â”‚   â”‚   â”œâ”€â”€ Coverage: Delhi NCR
â”‚   â”‚   â””â”€â”€ Mobile: Yes (30 km)
â”‚   â”œâ”€â”€ Vikram Singh (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: Patient Monitors, Lab Equipment
â”‚   â”‚   â””â”€â”€ Mobile: Yes (40 km)
â”‚   â”œâ”€â”€ Anil Kumar (Engineer)
â”‚   â”‚   â”œâ”€â”€ Skills: Installation Specialist (all equipment)
â”‚   â”‚   â””â”€â”€ Mobile: Yes (60 km)
â”‚   â””â”€â”€ Deepak R (Engineer)
â”‚       â”œâ”€â”€ Skills: General Service (all equipment)
â”‚       â””â”€â”€ Mobile: Yes (50 km)
â””â”€â”€ Warehouse (Facility)
```

**When Max Hospital Gurgaon raises a ticket for GE CT Scanner (purchased from City Medical):**
1. System checks: Sold by = City Medical
2. Routing Tier 2: Sub-sub_SUB_DEALER Service (City Medical has priority)
3. Filters: GE CT skills + available + Gurgaon coverage
4. **Result:** Amit Sharma (multi-brand expert, has GE CT certification)
5. Auto-assigns to Amit

---

### 3.3 Scenario: Hospital with In-House BME Team

**Apollo Hospital Delhi:**

```
Organization: Apollo Hospital Delhi
â”œâ”€â”€ Biomedical Engineering Department
â”‚   â”œâ”€â”€ Head BME: Dr. Rajiv Mehta
â”‚   â”œâ”€â”€ Senior Engineers (5)
â”‚   â”‚   â”œâ”€â”€ Mohan L (Senior BME)
â”‚   â”‚   â”‚   â”œâ”€â”€ Skills: All imaging equipment (Siemens, GE, Philips)
â”‚   â”‚   â”‚   â”œâ”€â”€ 15 years experience
â”‚   â”‚   â”‚   â”œâ”€â”€ Certifications: Siemens, GE, Philips authorized
â”‚   â”‚   â”‚   â””â”€â”€ Can handle: 80% of issues
â”‚   â”‚   â”œâ”€â”€ Prakash K (Senior BME)
â”‚   â”‚   â”‚   â”œâ”€â”€ Skills: Patient monitoring, ICU equipment
â”‚   â”‚   â”‚   â””â”€â”€ 12 years experience
â”‚   â”‚   â””â”€â”€ [3 more senior engineers]
â”‚   â””â”€â”€ Junior Engineers (10)
â”‚       â”œâ”€â”€ Skills: Basic maintenance, preventive maintenance
â”‚       â””â”€â”€ Supervised by senior engineers
```

**Scenario A: Siemens CT Scanner issue (manufacturer has engineers nearby):**
1. Routing Tier 1: Siemens Service â†’ Auto-assigned to Siemens engineer
2. Apollo's BME team NOT involved (vendor handling)

**Scenario B: GE Patient Monitor issue (GE has no engineers in Delhi, Sub-sub_SUB_DEALER unavailable):**
1. Routing Tier 1: GE Service â†’ No engineers available in Delhi
2. Routing Tier 2: Sub-sub_SUB_DEALER Service â†’ Sub-sub_SUB_DEALER engineers busy (3/3 active)
3. Routing Tier 3: Channel Partner Service â†’ None in Delhi
4. **Routing Tier 5: Apollo's In-House BME Team** (FALLBACK)
5. Filters: GE Patient Monitor skills + available
6. **Result:** Prakash K (patient monitoring expert) assigned
7. Prakash K fixes the issue using Apollo's parts inventory
8. Cost: Internal (no vendor charges)

**Scenario C: Philips Ultrasound issue (no external engineer available within SLA):**
1. All external tiers timeout or unavailable
2. **Routing Tier 5: Apollo's In-House BME** (CRITICAL FALLBACK)
3. **Result:** Mohan L (multi-brand expert) assigned
4. Mohan L handles it, orders parts if needed

---

### 3.4 Scenario: Service Provider with Engineer Network

**QuickFix Medical Services (Pan-India):**

```
Organization: QuickFix Medical Services
â”œâ”€â”€ Delhi Service Center
â”‚   â”œâ”€â”€ 8 Engineers
â”‚   â””â”€â”€ Coverage: Delhi NCR
â”œâ”€â”€ Mumbai Service Center
â”‚   â”œâ”€â”€ 12 Engineers
â”‚   â””â”€â”€ Coverage: Mumbai, Pune, Ahmedabad
â”œâ”€â”€ Bangalore Service Center
â”‚   â”œâ”€â”€ 10 Engineers
â”‚   â””â”€â”€ Coverage: Bangalore, Hyderabad
â””â”€â”€ Regional Engineers (Freelance/Contract)
    â”œâ”€â”€ 50+ engineers across India
    â””â”€â”€ On-demand availability
```

**When a hospital in Jaipur raises a ticket for Siemens X-Ray:**
1. Tier 1: Siemens Service â†’ No engineer in Jaipur
2. Tier 2: Sub-sub_SUB_DEALER Service â†’ No Sub-sub_SUB_DEALER engineers
3. Tier 4: **Service Provider (QuickFix) â†’ Authorized for Siemens**
4. Finds nearest QuickFix engineer in Jaipur (freelance contract)
5. Auto-assigns

---

## 4. Database Schema Additions

```sql
-- Engineers Table
CREATE TABLE engineers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- Identity
  employee_id TEXT,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  whatsapp_number TEXT,
  
  -- Employment
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  org_type TEXT NOT NULL,
  employment_type TEXT DEFAULT 'full_time',
  joining_date DATE,
  
  -- Location
  primary_facility_id UUID REFERENCES organization_facilities(id),
  mobile_engineer BOOLEAN DEFAULT true,
  current_location POINT,
  
  -- Coverage
  coverage_radius_km INT,
  coverage_pincodes TEXT[],
  coverage_cities TEXT[],
  coverage_states TEXT[],
  
  -- Availability
  status TEXT DEFAULT 'available',
  active_tickets INT DEFAULT 0,
  max_daily_tickets INT DEFAULT 5,
  
  -- Schedule
  working_hours JSONB,
  on_call_24x7 BOOLEAN DEFAULT false,
  
  -- Performance
  total_tickets_resolved INT DEFAULT 0,
  avg_resolution_time_hours NUMERIC(6,2),
  customer_rating NUMERIC(3,2) DEFAULT 0,
  first_time_fix_rate NUMERIC(5,2) DEFAULT 0,
  
  -- Contact
  preferred_contact_method TEXT DEFAULT 'phone',
  language_preferences TEXT[],
  
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_eng_status CHECK (status IN ('available', 'on_job', 'on_leave', 'inactive')),
  CONSTRAINT chk_eng_org_type CHECK (org_type IN (
    'manufacturer', 'Channel Partner', 'Sub-sub_SUB_DEALER', 'hospital', 'clinic',
    'service_provider', 'laboratory', 'diagnostic_center'
  ))
);

CREATE INDEX idx_eng_org ON engineers(org_id);
CREATE INDEX idx_eng_status ON engineers(status);
CREATE INDEX idx_eng_location ON engineers USING GIST(current_location);
CREATE INDEX idx_eng_facility ON engineers(primary_facility_id);

-- Engineer Skills Table
CREATE TABLE engineer_skills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
  
  -- Skill Definition
  skill_type TEXT NOT NULL,
  
  -- Equipment Skills
  equipment_category TEXT,
  equipment_type TEXT,
  equipment_models TEXT[],
  
  -- Manufacturer Skills
  manufacturer_id UUID REFERENCES organizations(id),
  manufacturer_name TEXT,
  manufacturer_authorized BOOLEAN DEFAULT false,
  
  -- Level
  proficiency_level TEXT DEFAULT 'intermediate',
  
  -- Certification
  certification_name TEXT,
  certification_number TEXT,
  certification_authority TEXT,
  certified_date DATE,
  expiry_date DATE,
  certificate_document_url TEXT,
  
  -- Capabilities
  can_install BOOLEAN DEFAULT false,
  can_calibrate BOOLEAN DEFAULT false,
  can_repair BOOLEAN DEFAULT true,
  can_train_users BOOLEAN DEFAULT false,
  
  -- Experience
  years_of_experience INT DEFAULT 0,
  tickets_resolved_for_this_skill INT DEFAULT 0,
  
  -- Verification
  verified BOOLEAN DEFAULT false,
  verified_by TEXT,
  verified_date DATE,
  
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_skill_type CHECK (skill_type IN (
    'equipment_category', 'equipment_type', 'equipment_model',
    'manufacturer_general', 'service_type'
  )),
  CONSTRAINT chk_proficiency CHECK (proficiency_level IN (
    'beginner', 'intermediate', 'advanced', 'expert'
  ))
);

CREATE INDEX idx_skill_engineer ON engineer_skills(engineer_id);
CREATE INDEX idx_skill_manufacturer ON engineer_skills(manufacturer_id);
CREATE INDEX idx_skill_equipment_type ON engineer_skills(equipment_type);
CREATE INDEX idx_skill_equipment_category ON engineer_skills(equipment_category);
CREATE INDEX idx_skill_expiry ON engineer_skills(expiry_date) WHERE expiry_date IS NOT NULL;

-- Engineer Availability Table
CREATE TABLE engineer_availability (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
  
  date DATE NOT NULL,
  available BOOLEAN DEFAULT true,
  
  reason TEXT,
  notes TEXT,
  
  available_slots JSONB,
  blocked_slots JSONB,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  UNIQUE(engineer_id, date)
);

CREATE INDEX idx_avail_engineer ON engineer_availability(engineer_id);
CREATE INDEX idx_avail_date ON engineer_availability(date);

-- Engineer Assignments Table
CREATE TABLE engineer_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- Assignment
  engineer_id UUID NOT NULL REFERENCES engineers(id),
  ticket_id UUID NOT NULL REFERENCES service_tickets(id),
  equipment_id UUID NOT NULL REFERENCES equipment(id),
  
  -- Assignment Details
  assigned_by UUID,
  assigned_at TIMESTAMPTZ DEFAULT NOW(),
  assignment_type TEXT DEFAULT 'auto',
  
  -- Status
  status TEXT DEFAULT 'assigned',
  
  -- Timeline
  accepted_at TIMESTAMPTZ,
  en_route_at TIMESTAMPTZ,
  reached_site_at TIMESTAMPTZ,
  work_started_at TIMESTAMPTZ,
  work_completed_at TIMESTAMPTZ,
  
  -- Location
  engineer_start_location POINT,
  customer_location POINT,
  travel_distance_km NUMERIC(8,2),
  estimated_arrival TIMESTAMPTZ,
  actual_arrival TIMESTAMPTZ,
  
  -- Work Details
  issue_description TEXT,
  diagnosis TEXT,
  actions_taken TEXT,
  parts_used JSONB,
  
  -- Customer Feedback
  customer_signature TEXT,
  customer_rating INT,
  customer_feedback TEXT,
  
  -- Photos
  before_photos TEXT[],
  after_photos TEXT[],
  
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_assignment_status CHECK (status IN (
    'assigned', 'accepted', 'en_route', 'on_site', 'completed', 'cancelled'
  )),
  CONSTRAINT chk_assignment_type CHECK (assignment_type IN ('auto', 'manual'))
);

CREATE INDEX idx_assignment_engineer ON engineer_assignments(engineer_id);
CREATE INDEX idx_assignment_ticket ON engineer_assignments(ticket_id);
CREATE INDEX idx_assignment_equipment ON engineer_assignments(equipment_id);
CREATE INDEX idx_assignment_status ON engineer_assignments(status);

-- Add foreign key to service_tickets for assigned engineer
ALTER TABLE service_tickets 
ADD COLUMN assigned_engineer_id UUID REFERENCES engineers(id),
ADD COLUMN assignment_tier INT, -- Which routing tier was used
ADD COLUMN assignment_tier_name TEXT; -- "OEM Service", "Sub-sub_SUB_DEALER Service", etc.
```

---

## 5. API Endpoints

```typescript
// Engineers API
GET    /api/v1/engineers
GET    /api/v1/engineers/:id
POST   /api/v1/engineers
PATCH  /api/v1/engineers/:id
DELETE /api/v1/engineers/:id

// Organization's engineers
GET    /api/v1/organizations/:id/engineers
POST   /api/v1/organizations/:id/engineers

// Facility's engineers
GET    /api/v1/facilities/:id/engineers

// Engineer Skills
GET    /api/v1/engineers/:id/skills
POST   /api/v1/engineers/:id/skills
PATCH  /api/v1/skills/:id
DELETE /api/v1/skills/:id

// Engineer Availability
GET    /api/v1/engineers/:id/availability
POST   /api/v1/engineers/:id/availability
PATCH  /api/v1/availability/:id

// Engineer Assignments
GET    /api/v1/engineers/:id/assignments
GET    /api/v1/assignments/:id
PATCH  /api/v1/assignments/:id

// Assignment Status Updates (Real-time)
POST   /api/v1/assignments/:id/accept
POST   /api/v1/assignments/:id/en-route
POST   /api/v1/assignments/:id/reach-site
POST   /api/v1/assignments/:id/start-work
POST   /api/v1/assignments/:id/complete

// Service Request Routing
POST   /api/v1/service-tickets/:id/route        // Auto-route to best engineer
GET    /api/v1/service-tickets/:id/routing      // Get routing options
POST   /api/v1/service-tickets/:id/assign       // Manual assignment

// Engineer Search & Availability
GET    /api/v1/engineers/search?skills=CT&location=Delhi&available=true
GET    /api/v1/engineers/available?equipment_id=xxx&date=2025-10-12
GET    /api/v1/engineers/nearby?lat=28.7041&lng=77.1025&radius=50

// Performance Analytics
GET    /api/v1/engineers/:id/performance
GET    /api/v1/organizations/:id/engineers/performance
```

---

## 6. Implementation Checklist

### Phase 1: Database & Core Models âœ…
- [ ] Create engineers table
- [ ] Create engineer_skills table
- [ ] Create engineer_availability table
- [ ] Create engineer_assignments table
- [ ] Add seed data (30+ engineers across orgs)

### Phase 2: Routing Logic âœ…
- [ ] Implement tier-based routing algorithm
- [ ] Implement skill matching logic
- [ ] Implement location-based filtering
- [ ] Implement availability checking
- [ ] Add fallback to client's in-house engineers

### Phase 3: APIs âœ…
- [ ] Engineers CRUD APIs
- [ ] Skills management APIs
- [ ] Availability management APIs
- [ ] Assignment tracking APIs
- [ ] Routing APIs

### Phase 4: Frontend âœ…
- [ ] Engineer management page (per organization)
- [ ] Engineer profile page
- [ ] Skills & certifications management
- [ ] Assignment tracking UI
- [ ] Real-time engineer location tracking

### Phase 5: Mobile App (Future) ðŸ”®
- [ ] Engineer mobile app for assignment acceptance
- [ ] Real-time status updates
- [ ] Navigation to customer site
- [ ] Digital signature capture
- [ ] Photo upload

---

## 7. Key Design Decisions

### âœ… Multi-Entity Engineer Support
- Engineers belong to manufacturers, Sub-Sub-sub_sub_SUB_DEALERs, Channel Partners, hospitals, service providers
- Each organization manages their own engineer team
- Flexible assignment based on availability and skills

### âœ… Skill-Based Routing
- Engineers have specific skills for equipment types, models, manufacturers
- Certification tracking with expiry dates
- Proficiency levels (beginner â†’ expert)

### âœ… Tier-Based Routing with Fallback
- **Priority:** OEM â†’ Sub-sub_SUB_DEALER â†’ Channel Partner â†’ Service Provider â†’ **Client In-House**
- Timeout-based escalation
- Client's engineers as critical fallback option

### âœ… Location-Based Assignment
- Mobile engineers with coverage radius
- Real-time location tracking
- Distance-based optimization

### âœ… Workload Management
- Max daily tickets per engineer
- Active ticket tracking
- Availability scheduling

### âœ… Performance Tracking
- Customer ratings
- Resolution times
- First-time fix rates
- Tickets resolved

---

**Status:** ðŸ“ ENGINEER MANAGEMENT DESIGN COMPLETE  
**Ready for:** Implementation  
**Integration:** Will be integrated with Organizations Architecture

