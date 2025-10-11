# 🏗️ Detailed Organizations Architecture - Complete Design

**Document Type:** Brainstorming & Detailed Design  
**Date:** October 11, 2025, 10:30 PM IST  
**Status:** Draft for Review

---

## 📋 Table of Contents

1. [Entity Model - Detailed Design](#entity-model)
2. [Real-World Scenarios](#real-world-scenarios)
3. [Relationship Patterns](#relationship-patterns)
4. [Data Structures](#data-structures)
5. [Dashboard Designs by User Type](#dashboard-designs)
6. [Use Cases & Workflows](#use-cases)
7. [Technical Implementation](#technical-implementation)

---

## 1. Entity Model - Detailed Design

### 1.1 Organization Entity (Core)

```typescript
interface Organization {
  // Identity
  id: UUID;
  name: string;
  display_name?: string;
  org_type: OrgType; // See below
  sub_type?: string; // For granular classification
  
  // Status & Lifecycle
  status: 'active' | 'inactive' | 'suspended' | 'pending_approval';
  verified: boolean;
  verification_date?: Date;
  onboarded_date: Date;
  
  // Legal & Registration
  legal_entity_name?: string;
  registration_number?: string; // CIN/GSTIN/PAN
  tax_id?: string;
  incorporation_date?: Date;
  
  // Contact Information
  primary_contact: ContactPerson;
  secondary_contacts?: ContactPerson[];
  
  // Address Hierarchy
  headquarters: Address;
  registered_office?: Address;
  
  // Business Information
  industry_segments?: string[]; // Imaging, Diagnostics, Cardiology, etc.
  certifications?: Certification[]; // ISO, FDA, CE Mark, etc.
  annual_turnover?: Money;
  employee_count?: number;
  year_established?: number;
  
  // Digital Presence
  website?: string;
  social_media?: SocialMedia;
  logo_url?: string;
  
  // Platform Integration
  external_refs: ExternalReference[]; // ERP IDs, SAP IDs, etc.
  metadata: JSONB; // Flexible additional data
  
  // Multi-tenant
  tenant_id: string;
  
  // Audit
  created_at: DateTime;
  updated_at: DateTime;
  created_by: string;
  updated_by?: string;
}

enum OrgType {
  MANUFACTURER = 'manufacturer',
  DISTRIBUTOR = 'distributor',
  DEALER = 'dealer',
  SUPPLIER = 'supplier',
  HOSPITAL = 'hospital',
  LABORATORY = 'laboratory',
  DIAGNOSTIC_CENTER = 'diagnostic_center',
  CLINIC = 'clinic',
  SERVICE_PROVIDER = 'service_provider',
  LOGISTICS_PARTNER = 'logistics_partner',
  INSURANCE_PROVIDER = 'insurance_provider',
  GOVERNMENT_BODY = 'government_body',
  OTHER = 'other'
}
```

---

### 1.2 Facilities/Locations (Multi-Location Support)

**Real-World Scenario:** 
- Siemens has manufacturing in Mumbai, Bangalore, Chennai
- Apollo Hospitals has 70+ hospitals across India
- Distributor has warehouses in 10 cities

```typescript
interface OrganizationFacility {
  id: UUID;
  org_id: UUID; // Parent organization
  
  // Facility Identity
  facility_name: string; // "Mumbai Manufacturing Plant"
  facility_code: string; // "SIE-MUM-01"
  facility_type: FacilityType;
  
  // Location
  address: Address;
  geo_location: Point; // Lat/Long for mapping
  
  // Operational Details
  capacity?: string; // "500 units/month"
  operational_hours?: OperatingHours;
  contact_person: ContactPerson;
  
  // Capabilities
  services_offered?: string[]; // Manufacturing, Servicing, Training
  certifications?: Certification[];
  equipment_types?: string[]; // What they manufacture/service
  
  // Coverage
  service_radius_km?: number;
  coverage_pincodes?: string[];
  coverage_states?: string[];
  
  // Status
  status: 'active' | 'inactive' | 'under_construction';
  operational_since?: Date;
  
  // Integration
  external_refs: ExternalReference[];
  metadata: JSONB;
  
  created_at: DateTime;
  updated_at: DateTime;
}

enum FacilityType {
  MANUFACTURING_PLANT = 'manufacturing_plant',
  ASSEMBLY_UNIT = 'assembly_unit',
  R&D_CENTER = 'rnd_center',
  WAREHOUSE = 'warehouse',
  DISTRIBUTION_CENTER = 'distribution_center',
  SERVICE_CENTER = 'service_center',
  TRAINING_CENTER = 'training_center',
  SALES_OFFICE = 'sales_office',
  SHOWROOM = 'showroom',
  HOSPITAL_UNIT = 'hospital_unit',
  LABORATORY_UNIT = 'laboratory_unit',
  DIAGNOSTIC_CENTER = 'diagnostic_center',
  CLINIC = 'clinic'
}
```

---

### 1.3 Organization Relationships (Complex Network)

**Real-World Scenarios:**
- Dealer works with 5 manufacturers
- Manufacturer has 20 distributors across India
- Distributor covers 100 dealers in their region
- Service provider partners with multiple manufacturers

```typescript
interface OrganizationRelationship {
  id: UUID;
  
  // Relationship Parties
  parent_org_id: UUID; // "From" organization
  child_org_id: UUID; // "To" organization
  
  // Relationship Type & Details
  relationship_type: RelationshipType;
  relationship_status: 'active' | 'inactive' | 'pending' | 'expired';
  
  // Business Terms
  start_date: Date;
  end_date?: Date;
  auto_renew: boolean;
  
  // Territory/Coverage
  exclusive: boolean; // Exclusive distributor/dealer?
  territory?: Territory;
  
  // Commercial Terms
  commission_percentage?: number;
  volume_discounts?: VolumeDiscount[];
  payment_terms?: PaymentTerms;
  credit_limit?: Money;
  
  // Product Scope
  product_categories?: string[]; // Which products covered
  excluded_products?: string[];
  
  // Performance Metrics
  annual_target?: Money;
  quarterly_targets?: QuarterlyTarget[];
  performance_tier?: 'platinum' | 'gold' | 'silver' | 'bronze';
  
  // Operational
  priority_level?: number; // For multi-manufacturer dealers
  default_supplier?: boolean;
  
  // Legal
  contract_reference?: string;
  agreement_documents?: Document[];
  
  // Audit
  created_at: DateTime;
  updated_at: DateTime;
  created_by: string;
  notes?: string;
}

enum RelationshipType {
  // Manufacturer relationships
  AUTHORIZED_DISTRIBUTOR = 'authorized_distributor',
  EXCLUSIVE_DISTRIBUTOR = 'exclusive_distributor',
  REGIONAL_DISTRIBUTOR = 'regional_distributor',
  AUTHORIZED_DEALER = 'authorized_dealer',
  SERVICE_PARTNER = 'service_partner',
  
  // Distributor relationships
  DEALER_NETWORK = 'dealer_network',
  SUB_DISTRIBUTOR = 'sub_distributor',
  
  // Service relationships
  AMC_PROVIDER = 'amc_provider',
  SPARE_PARTS_SUPPLIER = 'spare_parts_supplier',
  
  // Business relationships
  STRATEGIC_PARTNER = 'strategic_partner',
  OEM_PARTNER = 'oem_partner',
  WHITE_LABEL_PARTNER = 'white_label_partner',
  
  // Sales relationships
  DIRECT_BUYER = 'direct_buyer',
  INSTITUTIONAL_BUYER = 'institutional_buyer',
  
  // Support relationships
  LOGISTICS_PARTNER = 'logistics_partner',
  INSURANCE_PARTNER = 'insurance_partner',
  FINANCING_PARTNER = 'financing_partner'
}
```

---

### 1.4 Territory Management

```typescript
interface Territory {
  id: UUID;
  name: string; // "North India", "Mumbai Metro"
  code: string; // "TERR-NORTH-01"
  
  // Geographic Coverage
  coverage_type: 'pincode' | 'city' | 'district' | 'state' | 'region' | 'custom';
  
  // Specific Areas
  states?: string[];
  cities?: string[];
  districts?: string[];
  pincodes?: string[];
  custom_boundaries?: Polygon; // GeoJSON
  
  // Hierarchy
  parent_territory_id?: UUID;
  
  // Assignment
  assigned_to_org_id?: UUID;
  assigned_to_facility_id?: UUID;
  
  // Market Data
  estimated_market_size?: Money;
  potential_customers?: number;
  competitor_presence?: CompetitorInfo[];
  
  metadata: JSONB;
  created_at: DateTime;
  updated_at: DateTime;
}
```

---

## 2. Real-World Scenarios

### 2.1 Manufacturer Scenario: Siemens Healthineers

**Organization Structure:**

```
Siemens Healthineers India Ltd. (MANUFACTURER)
├── Facilities
│   ├── Mumbai Manufacturing Plant (MANUFACTURING_PLANT)
│   │   - Products: CT Scanners, X-Ray
│   │   - Capacity: 100 units/month
│   │   - Certifications: ISO 13485, CE Mark, FDA
│   ├── Bangalore R&D Center (R&D_CENTER)
│   ├── Chennai Service Hub (SERVICE_CENTER)
│   │   - Coverage: Tamil Nadu, Kerala, Karnataka
│   └── Delhi Sales Office (SALES_OFFICE)
│
├── Distribution Network
│   ├── North Region
│   │   └── MedEquip Distributors Pvt Ltd (EXCLUSIVE_DISTRIBUTOR)
│   │       - Territory: Delhi, UP, Punjab, Haryana
│   │       - Products: All Siemens products
│   │       - Dealer Network: 25 dealers
│   ├── South Region
│   │   └── HealthTech Solutions (REGIONAL_DISTRIBUTOR)
│   │       - Territory: TN, Kerala, Karnataka, AP
│   │       - Dealer Network: 30 dealers
│   └── West Region
│       └── Western Medical Supplies (AUTHORIZED_DISTRIBUTOR)
│           - Territory: Maharashtra, Gujarat, MP
│           - Dealer Network: 20 dealers
│
├── Service Partners
│   ├── QuickFix Medical Services (SERVICE_PARTNER)
│   │   - Coverage: Pan-India
│   │   - SLA: 24-hour response
│   └── TechCare India (AMC_PROVIDER)
│       - Coverage: Metro cities
│
└── Direct Institutional Customers
    ├── Apollo Hospitals (DIRECT_BUYER)
    ├── Fortis Healthcare (DIRECT_BUYER)
    └── AIIMS Delhi (INSTITUTIONAL_BUYER)
```

---

### 2.2 Distributor Scenario: MedEquip Distributors

**Organization Structure:**

```
MedEquip Distributors Pvt Ltd (DISTRIBUTOR)
├── Facilities
│   ├── Delhi Warehouse (DISTRIBUTION_CENTER)
│   │   - Inventory: ₹50 Cr
│   │   - Coverage: Delhi NCR
│   ├── Chandigarh Service Center (SERVICE_CENTER)
│   ├── Lucknow Branch (SALES_OFFICE)
│   └── Jaipur Showroom (SHOWROOM)
│
├── Manufacturer Partnerships (Multi-Brand)
│   ├── Siemens Healthineers (EXCLUSIVE_DISTRIBUTOR)
│   │   - Products: All imaging equipment
│   │   - Territory: North India
│   │   - Commission: 8-12%
│   │   - Annual Target: ₹100 Cr
│   ├── GE Healthcare (AUTHORIZED_DISTRIBUTOR)
│   │   - Products: Patient monitoring
│   │   - Territory: Delhi, UP
│   │   - Commission: 6-10%
│   ├── Philips Healthcare (REGIONAL_DISTRIBUTOR)
│   │   - Products: Ultrasound systems
│   │   - Territory: North India
│   └── Local OEMs (AUTHORIZED_DISTRIBUTOR)
│       - Products: Basic equipment
│
├── Dealer Network (75 dealers)
│   ├── Premium Dealers (10)
│   │   - Annual Purchase: >₹2 Cr
│   │   - Credit Limit: ₹50 L
│   │   - Priority: Platinum
│   ├── Standard Dealers (40)
│   │   - Annual Purchase: ₹50L - ₹2Cr
│   │   - Credit Limit: ₹20 L
│   └── Small Dealers (25)
│       - Annual Purchase: <₹50L
│       - Credit Limit: ₹5 L
│
└── Customer Segments
    ├── Corporate Hospitals (Direct Sales)
    ├── Private Clinics (Through Dealers)
    └── Government Hospitals (Tender Sales)
```

---

### 2.3 Dealer Scenario: City Medical Equipment Co.

**Organization Structure:**

```
City Medical Equipment Co. (DEALER)
├── Facilities
│   ├── Main Showroom (SHOWROOM)
│   │   - Location: Connaught Place, Delhi
│   │   - Display: 50+ equipment
│   ├── Service Center (SERVICE_CENTER)
│   │   - 5 service engineers
│   │   - Coverage: Delhi NCR
│   └── Warehouse (WAREHOUSE)
│       - Inventory: ₹2 Cr
│
├── Supplier Relationships (Multi-Manufacturer)
│   ├── MedEquip Distributors (PRIMARY)
│   │   - Brands: Siemens, GE, Philips
│   │   - Credit Terms: 45 days
│   │   - Discount: 15%
│   ├── HealthTech Solutions (SECONDARY)
│   │   - Brands: Medtronic, Abbott
│   │   - Credit Terms: 30 days
│   └── Direct from Manufacturers
│       - Small orders, special items
│
├── Customer Base
│   ├── Private Hospitals (60%)
│   │   - 25 regular customers
│   │   - AMC contracts: 100+ equipment
│   ├── Polyclinics (25%)
│   │   - 40 customers
│   ├── Diagnostic Centers (10%)
│   │   - 15 customers
│   └── Individual Doctors (5%)
│       - 50+ customers
│
└── Service Operations
    ├── Installation Services
    ├── AMC Contracts (150 active)
    ├── Spare Parts Supply
    └── Training Services
```

---

### 2.4 Hospital Scenario: Apollo Hospitals

**Organization Structure:**

```
Apollo Hospitals Enterprise Ltd (HOSPITAL)
├── Hospital Network (70+ locations)
│   ├── Apollo Hospital Delhi (HOSPITAL_UNIT)
│   │   - Equipment: 500+ units
│   │   - AMCs: 450 active
│   │   - Annual Procurement: ₹20 Cr
│   ├── Apollo Spectra Bangalore (HOSPITAL_UNIT)
│   ├── Apollo Diagnostics Mumbai (DIAGNOSTIC_CENTER)
│   └── [68 more locations...]
│
├── Procurement Structure
│   ├── Centralized Procurement (70%)
│   │   - High-value equipment
│   │   - Direct from manufacturers
│   │   - Annual contracts
│   ├── Regional Procurement (20%)
│   │   - Mid-value equipment
│   │   - Through distributors
│   └── Local Procurement (10%)
│       - Consumables, small equipment
│       - Through local dealers
│
├── Vendor Relationships
│   ├── Direct from Manufacturers (Tier 1)
│   │   - Siemens: CT, MRI (10-year contract)
│   │   - GE: Patient monitors (5-year contract)
│   │   - Philips: Ultrasound systems
│   ├── Through Distributors (Tier 2)
│   │   - MedEquip: General equipment
│   │   - HealthTech: Lab equipment
│   └── Local Dealers (Tier 3)
│       - Quick replacements
│       - Emergency supplies
│
└── Service Management
    ├── In-House BME Team (50 engineers)
    ├── OEM Service Contracts
    ├── Third-Party AMCs
    └── Parts Inventory (₹5 Cr)
```

---

## 3. Relationship Patterns

### 3.1 Hierarchical Patterns

#### Pattern 1: Manufacturer → Distributor → Dealer → Hospital

```
Siemens (MANUFACTURER)
  ↓ EXCLUSIVE_DISTRIBUTOR
MedEquip (DISTRIBUTOR)
  ↓ DEALER_NETWORK
City Medical (DEALER)
  ↓ DIRECT_BUYER
Apollo Hospital Delhi (HOSPITAL)
```

**Key Attributes:**
- **Siemens ↔ MedEquip:**
  - Exclusive territory: North India
  - All Siemens products
  - 8-12% commission
  - ₹100 Cr annual target
  
- **MedEquip ↔ City Medical:**
  - Non-exclusive
  - Credit limit: ₹50 L
  - 45-day payment terms
  - 15% dealer discount
  
- **City Medical ↔ Apollo:**
  - Equipment sales
  - AMC contracts
  - Installation services
  - Training support

---

#### Pattern 2: Multi-Brand Distributor

```
MedEquip Distributors (DISTRIBUTOR)
  ├── Siemens (EXCLUSIVE_DISTRIBUTOR)
  ├── GE (AUTHORIZED_DISTRIBUTOR)
  ├── Philips (REGIONAL_DISTRIBUTOR)
  └── Medtronic (AUTHORIZED_DISTRIBUTOR)
```

**Decision Logic:**
- Exclusive: Only MedEquip can sell Siemens in North India
- Non-Exclusive: Multiple distributors for GE
- Priority: Siemens gets first priority for resources

---

#### Pattern 3: Multi-Manufacturer Dealer

```
City Medical (DEALER)
  ├── MedEquip → Siemens, GE, Philips (PRIMARY, 70%)
  ├── HealthTech → Medtronic, Abbott (SECONDARY, 20%)
  └── Direct → Small OEMs (10%)
```

**Selection Logic:**
- Check MedEquip first (better terms)
- HealthTech for brands not with MedEquip
- Direct for special/urgent orders

---

### 3.2 Service Network Patterns

```
Equipment Installation at Apollo Delhi
  ↓ Needs Service
Service Request Created
  ↓ AI Routing
Check Service Hierarchy:
  1. Manufacturer's Own Service Center (Siemens Chennai)
  2. Authorized Service Partner (QuickFix Delhi)
  3. Dealer Service Team (City Medical)
  4. Third-Party Provider (TechCare)
```

---

### 3.3 Territory Management Patterns

#### Scenario: New Dealer Registration

```
Request: New dealer in Jaipur wants Siemens dealership

Check:
1. Jaipur in which territory?
   → North India (covered by MedEquip exclusive)
   
2. Can dealer register?
   → No: MedEquip is exclusive distributor
   → Dealer must buy through MedEquip
   → Dealer becomes MedEquip's sub-dealer

3. If non-exclusive:
   → Yes: Can become direct dealer
   → Territory check: No overlap with exclusive dealers
```

---

## 4. Data Structures

### 4.1 Contact Person

```typescript
interface ContactPerson {
  id: UUID;
  name: string;
  designation: string;
  department?: string;
  
  // Contact Methods
  email: string;
  primary_phone: string;
  alternate_phone?: string;
  whatsapp_number?: string;
  
  // Address
  office_address?: Address;
  
  // Preferences
  preferred_contact_method: 'email' | 'phone' | 'whatsapp';
  language_preference?: string[];
  
  // Role
  is_primary: boolean;
  can_approve_orders?: boolean;
  can_raise_tickets?: boolean;
  
  // Status
  active: boolean;
}
```

---

### 4.2 Address (Indian Context)

```typescript
interface Address {
  id: UUID;
  address_type: 'headquarters' | 'registered_office' | 'branch' | 'warehouse' | 'service_center';
  
  // Address Components
  building_name?: string;
  street_address: string;
  locality: string;
  landmark?: string;
  city: string;
  district?: string;
  state: string;
  pincode: string;
  country: string;
  
  // Geo
  latitude?: number;
  longitude?: number;
  
  // Contact
  phone?: string;
  email?: string;
  
  // Logistics
  delivery_instructions?: string;
  access_hours?: OperatingHours;
  
  // Verification
  verified: boolean;
  verified_date?: Date;
}
```

---

### 4.3 Certification

```typescript
interface Certification {
  id: UUID;
  certification_type: string; // 'ISO 13485', 'CE Mark', 'FDA', etc.
  certification_number: string;
  issued_by: string;
  issue_date: Date;
  expiry_date?: Date;
  status: 'active' | 'expired' | 'suspended';
  
  // Documents
  certificate_document_url?: string;
  verification_url?: string;
  
  // Scope
  scope?: string; // What it covers
  applicable_products?: string[];
  applicable_facilities?: UUID[];
}
```

---

### 4.4 Operating Hours

```typescript
interface OperatingHours {
  monday?: TimeRange;
  tuesday?: TimeRange;
  wednesday?: TimeRange;
  thursday?: TimeRange;
  friday?: TimeRange;
  saturday?: TimeRange;
  sunday?: TimeRange;
  
  // Special
  public_holidays_open: boolean;
  24x7: boolean;
  emergency_contact?: string;
  
  // Regional
  timezone: string;
}

interface TimeRange {
  open: string; // "09:00"
  close: string; // "18:00"
  breaks?: {start: string, end: string}[];
}
```

---

## 5. Dashboard Designs by User Type

### 5.1 Manufacturer Dashboard

**Primary User:** Siemens Sales Manager

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Active          │ Monthly Sales   │ Pending Orders  │ Service         │
│ Distributors    │ ₹45.2 Cr       │ 156 orders      │ Tickets         │
│ 23 ↑2          │ ↑12% MoM       │ ₹15.8 Cr       │ 234 ↓12%       │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Distribution Network (Map + List)
- **Map View:** India map with distributor locations
- **Filters:** Region, Performance Tier, Product Category
- **Actions:** Add Distributor, View Details, Manage Territory

**Distributor Table:**
| Distributor | Territory | Products | Monthly Sales | Dealers | Performance | Actions |
|-------------|-----------|----------|---------------|---------|-------------|---------|
| MedEquip Distributors | North India | All Products | ₹8.2 Cr | 75 | 🥇 Platinum | View Details |
| HealthTech Solutions | South India | Imaging Only | ₹6.5 Cr | 60 | 🥈 Gold | View Details |

#### Section 2: Sales Analytics
- **Chart:** Regional sales breakdown (bar chart)
- **Chart:** Product category performance (pie chart)
- **Trend:** Monthly sales trend (line chart)

#### Section 3: Territory Management
- **Map:** Territory coverage visualization
- **Gaps:** Uncovered areas highlighted
- **Opportunities:** Potential new distributor locations

#### Section 4: Dealer Performance (Top 20)
| Dealer | Location | Monthly Sales | Growth | AMC Contracts | Service Rating |
|--------|----------|---------------|--------|---------------|----------------|

#### Section 5: Service Network Status
- Active service tickets by region
- Average resolution time
- Customer satisfaction scores
- Parts inventory levels

#### Section 6: Quick Actions
- [ ] Add New Distributor
- [ ] Create Territory
- [ ] Bulk Product Upload
- [ ] Generate Sales Report
- [ ] Schedule Training

---

### 5.2 Distributor Dashboard

**Primary User:** MedEquip Operations Manager

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Active Dealers  │ Monthly Sales   │ Pending Orders  │ Inventory       │
│ 75 ↑5          │ ₹8.2 Cr        │ 89 orders       │ ₹12.5 Cr       │
│                 │ ↑15% MoM       │ ₹2.1 Cr        │ 850 SKUs       │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Multi-Brand Overview
**Manufacturer Performance:**
| Manufacturer | Products | Monthly Sales | Commission | Target Progress | Actions |
|--------------|----------|---------------|------------|-----------------|---------|
| Siemens | All Imaging | ₹4.2 Cr | ₹0.42 Cr | 84% | View Contract |
| GE Healthcare | Monitoring | ₹2.1 Cr | ₹0.15 Cr | 70% | View Products |
| Philips | Ultrasound | ₹1.9 Cr | ₹0.14 Cr | 95% | View Territory |

#### Section 2: Dealer Network
**Map View:** Dealers plotted on map with color coding
- 🟢 Green: High performers (>₹2 Cr/year)
- 🟡 Yellow: Standard (₹50L-₹2Cr/year)
- 🔴 Red: Needs attention (<₹50L/year)

**Dealer Table:**
| Dealer | Location | Brands | Monthly Sales | Outstanding | Credit Limit | Performance |
|--------|----------|--------|---------------|-------------|--------------|-------------|
| City Medical | Delhi | Multi | ₹45 L | ₹12 L | ₹50 L | 🥇 Platinum |
| Metro Healthcare | Gurgaon | Multi | ₹32 L | ₹8 L | ₹30 L | 🥈 Gold |

#### Section 3: Inventory Management
- **Stock Levels:** By product category
- **Fast-Moving Items:** Top 20 SKUs
- **Slow-Moving:** Items to push
- **Alerts:** Low stock, expiring items

#### Section 4: Orders & Logistics
- **Pending Orders:** From dealers
- **In-Transit:** Shipments tracking
- **Delivery Schedule:** Next 7 days
- **Backorders:** Items awaiting stock

#### Section 5: Financial Dashboard
- **Receivables:** Outstanding from dealers
- **Payables:** Due to manufacturers
- **Credit Utilization:** By dealer
- **Commission Earned:** By manufacturer

#### Section 6: Territory Insights
- **Coverage Map:** Areas served
- **Gaps:** Potential new dealer locations
- **Competition:** Competitor presence
- **Market Size:** Potential by region

---

### 5.3 Dealer Dashboard

**Primary User:** City Medical Owner

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Monthly Sales   │ Active AMCs     │ Pending Orders  │ Service Jobs    │
│ ₹45 L          │ 150 contracts   │ 12 orders       │ 23 active       │
│ ↑8% MoM        │ ₹18 L MRR      │ ₹8.5 L         │ 5 pending       │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Supplier Hub (Multi-Manufacturer View)
**Active Suppliers:**
| Supplier | Brands Available | Credit Used | Credit Limit | Payment Due | Next Delivery |
|----------|------------------|-------------|--------------|-------------|---------------|
| MedEquip | Siemens, GE, Philips | ₹12 L | ₹50 L | ₹8 L (3 days) | Tomorrow |
| HealthTech | Medtronic, Abbott | ₹5 L | ₹20 L | ₹3 L (5 days) | 3 days |

**Quick Actions:**
- Place Order with MedEquip
- Check Product Availability
- Request Credit Extension
- View Price Lists

#### Section 2: Product Catalog (Aggregated)
**Multi-Manufacturer Search:**
- Search: "CT Scanner"
- Results show products from all suppliers
- Compare prices, delivery times, credit terms
- "Add to Cart" → Auto-selects best supplier

| Product | Manufacturer | Supplier | Price | Delivery | Credit Terms | Stock |
|---------|--------------|----------|-------|----------|--------------|-------|
| CT Scanner Pro | Siemens | MedEquip | ₹1.2 Cr | 7 days | 45 days | ✓ Available |
| CT Elite 500 | GE | MedEquip | ₹1.1 Cr | 10 days | 45 days | On Order |

#### Section 3: Customer Management
**Hospital Accounts:**
| Customer | Type | Monthly Avg | Outstanding | AMC Value | Next Service | Status |
|----------|------|-------------|-------------|-----------|--------------|--------|
| Apollo Delhi | Hospital | ₹12 L | ₹5 L | ₹8 L/yr | 5 days | 🟢 Active |
| Max Hospital | Hospital | ₹8 L | ₹2 L | ₹6 L/yr | 15 days | 🟢 Active |
| City Clinic | Clinic | ₹2 L | ₹50 K | ₹1 L/yr | 2 days | 🟡 Follow-up |

#### Section 4: AMC Management
**Active Contracts:**
- **Total AMCs:** 150
- **Monthly Recurring Revenue:** ₹18 L
- **Upcoming Renewals:** 15 (next 30 days)
- **Expiring Soon:** 8 (action needed)

**AMC Calendar:**
| Equipment | Customer | Next Service | Status | Engineer | Actions |
|-----------|----------|--------------|--------|----------|---------|
| X-Ray Machine | Apollo Delhi | Tomorrow | Scheduled | Ramesh K | View Details |
| CT Scanner | Max Hospital | 3 days | Pending Parts | Suresh M | Order Parts |

#### Section 5: Service Operations
**Today's Schedule:**
- 5 AM service visits
- 3 installations pending
- 2 training sessions
- 1 demo scheduled

**Engineer Availability:**
| Engineer | Skills | Today's Jobs | Location | Status |
|----------|--------|--------------|----------|--------|
| Ramesh K | CT, MRI | 2 jobs | Apollo Delhi | 🟢 Available |
| Suresh M | X-Ray, Ultrasound | 3 jobs | Max Hospital | 🟡 Busy |

#### Section 6: Financial Summary
- **Sales This Month:** ₹45 L
- **Collections:** ₹38 L
- **Outstanding:** ₹28 L
- **Expenses:** ₹12 L (salaries, rent, etc.)
- **Net Profit:** ₹14 L

**Payment Alerts:**
- ⚠️ MedEquip: ₹8 L due in 3 days
- ⚠️ Customer: Apollo Hospital ₹5 L overdue

---

### 5.4 Hospital Dashboard

**Primary User:** Apollo Delhi Biomedical Engineer

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Total Equipment │ Active AMCs     │ Open Tickets    │ Monthly Spend   │
│ 500 units       │ 450 contracts   │ 15 tickets      │ ₹12 L          │
│ 485 operational │ ₹22 L/month    │ Avg: 4.2 hrs   │ Budget: 92%    │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Equipment Inventory
**By Department:**
| Department | Equipment Count | Operational | Under Service | Down | Utilization |
|------------|-----------------|-------------|---------------|------|-------------|
| Radiology | 45 | 43 | 1 | 1 | 87% |
| ICU | 120 | 118 | 2 | 0 | 94% |
| OT | 85 | 82 | 3 | 0 | 89% |
| Lab | 150 | 148 | 2 | 0 | 92% |

**Equipment Map:** Floor plan with equipment locations

#### Section 2: Service Request Portal
**Quick Service:**
- Scan QR code on equipment
- Auto-fills equipment details
- Select issue type
- Attach photos
- Submit → Auto-routes to correct vendor

**Open Tickets:**
| Equipment | Issue | Vendor | Priority | Created | Status | SLA |
|-----------|-------|--------|----------|---------|--------|-----|
| CT Scanner | Calibration | Siemens | High | 2 hours | Engineer En Route | 🟢 On Time |
| X-Ray | No Power | City Medical | Critical | 30 mins | Parts Ordered | 🟡 At Risk |

#### Section 3: AMC Management
**Active Contracts:**
| Vendor | Equipment Count | Monthly Cost | Next Service | Contract Expiry | Actions |
|--------|-----------------|--------------|--------------|-----------------|---------|
| Siemens Direct | 120 | ₹8 L | Tomorrow | 2 years | View Details |
| City Medical | 150 | ₹10 L | Next Week | 6 months | Renew Soon |
| QuickFix | 180 | ₹4 L | 15 days | 1 year | View Contract |

#### Section 4: Vendor Performance
**Scorecard:**
| Vendor | Active AMCs | Avg Response Time | Resolution Time | Uptime % | Rating | Trend |
|--------|-------------|-------------------|-----------------|----------|--------|-------|
| Siemens Direct | 120 | 2.1 hrs | 4.5 hrs | 99.2% | ⭐⭐⭐⭐⭐ | ↑ |
| City Medical | 150 | 3.5 hrs | 6.2 hrs | 98.5% | ⭐⭐⭐⭐ | → |
| QuickFix | 180 | 4.2 hrs | 8.1 hrs | 97.8% | ⭐⭐⭐ | ↓ |

#### Section 5: Procurement Planning
**Budget Tracker:**
- Annual Budget: ₹144 L
- Spent: ₹132 L (92%)
- Remaining: ₹12 L
- Forecast: On track

**Upcoming Purchases:**
| Item | Department | Quantity | Estimated Cost | Approval Status | Procurement Route |
|------|------------|----------|----------------|-----------------|-------------------|
| New CT Scanner | Radiology | 1 | ₹1.5 Cr | Approved | Direct (Siemens) |
| Patient Monitors | ICU | 20 | ₹40 L | Pending | RFQ (3 vendors) |

#### Section 6: Preventive Maintenance Calendar
**This Month:**
- 45 PM schedules
- 42 completed
- 3 pending
- 0 overdue

**Upcoming (Next 7 Days):**
| Date | Equipment | Type | Vendor | Status |
|------|-----------|------|--------|--------|
| Tomorrow | MRI Scanner | PM | Siemens | Scheduled |
| 2 days | Ventilators (10) | PM | GE | Confirmed |

---

### 5.5 Service Provider Dashboard

**Primary User:** QuickFix Service Manager

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Active Tickets  │ Engineers       │ Monthly Revenue │ Customer        │
│ 45 tickets      │ 25 engineers    │ ₹18 L          │ Satisfaction    │
│ 12 high priority│ 22 available    │ ↑10% MoM       │ 4.5/5.0        │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Ticket Management
**Active Tickets Board (Kanban):**
| New (8) | Assigned (15) | In Progress (12) | Pending Parts (6) | Resolved (4) |
|---------|---------------|------------------|-------------------|--------------|

**High Priority:**
| Ticket | Equipment | Customer | Issue | SLA | Engineer | Status |
|--------|-----------|----------|-------|-----|----------|--------|
| #1234 | CT Scanner | Apollo Delhi | No Power | 2 hrs left | Ramesh K | En Route |
| #1235 | Ventilator | Max Hospital | Alarm Issue | 4 hrs left | Suresh M | Diagnosing |

#### Section 2: Engineer Management
**Engineer Dashboard:**
| Engineer | Location | Skills | Today's Jobs | Completed | Rating | Status |
|----------|----------|--------|--------------|-----------|--------|--------|
| Ramesh K | Delhi | CT, MRI, X-Ray | 3 | 1 | ⭐4.8 | 🚗 En Route |
| Suresh M | Gurgaon | All Equipment | 4 | 2 | ⭐4.6 | 🔧 On Job |

**Coverage Map:**
- Engineers plotted on map
- Open tickets shown
- Auto-suggest nearest engineer
- Traffic-aware routing

#### Section 3: Parts Inventory
**Stock Levels:**
| Part Category | In Stock | Low Stock | Out of Stock | Ordered |
|---------------|----------|-----------|--------------|---------|
| CT Components | 45 | 5 | 2 | 8 |
| X-Ray Parts | 120 | 12 | 0 | 5 |
| General | 850 | 45 | 8 | 23 |

**Alerts:**
- ⚠️ CT Detector: Only 2 left (reorder now)
- ⚠️ X-Ray Tube: Lead time 15 days

#### Section 4: Customer Accounts
**Top Customers:**
| Customer | AMC Value | Monthly Calls | Avg Resolution | Outstanding | Actions |
|----------|-----------|---------------|----------------|-------------|---------|
| Apollo Network | ₹8 L/month | 45 tickets | 4.2 hrs | ₹0 | View Details |
| Max Healthcare | ₹6 L/month | 38 tickets | 5.1 hrs | ₹1.2 L | Follow Up |

#### Section 5: Performance Analytics
**This Month:**
- Tickets Resolved: 145
- Avg Response Time: 3.2 hrs (Target: 4 hrs) ✓
- Avg Resolution Time: 6.8 hrs (Target: 8 hrs) ✓
- First-Time Fix Rate: 78%
- Customer Satisfaction: 4.5/5.0

**Trends:**
- Response time improving
- Parts availability improving
- Engineer productivity up 12%

#### Section 6: Financial Dashboard
**Revenue Breakdown:**
- AMC Contracts: ₹12 L (67%)
- Break-Fix: ₹4 L (22%)
- Parts Sales: ₹2 L (11%)

**Profitability:**
- Revenue: ₹18 L
- Engineer Costs: ₹8 L
- Parts Cost: ₹3 L
- Overheads: ₹2 L
- Net Profit: ₹5 L (28%)

---

### 5.6 Platform Admin Dashboard

**Primary User:** Platform Operations Team

#### Top Stats Row
```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│ Total Orgs      │ Active Users    │ Monthly GMV     │ Platform        │
│ 1,245 orgs      │ 5,420 users     │ ₹245 Cr        │ Commission      │
│ ↑45 this month  │ ↑234 this month │ ↑18% MoM       │ ₹9.8 Cr        │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

#### Section 1: Organization Overview
**By Type:**
| Org Type | Count | Active | Verified | Pending Approval | Actions |
|----------|-------|--------|----------|------------------|---------|
| Manufacturers | 45 | 42 | 40 | 3 | View All |
| Distributors | 180 | 165 | 150 | 12 | View All |
| Dealers | 650 | 580 | 520 | 45 | View All |
| Hospitals | 280 | 270 | 260 | 8 | View All |
| Service Providers | 90 | 85 | 80 | 5 | View All |

**Recent Registrations:**
| Organization | Type | Registration Date | Status | Actions |
|--------------|------|-------------------|--------|---------|
| NewMed Distributors | Distributor | Today | Pending Verification | Review |
| TechCare Services | Service Provider | Yesterday | Verified | View Profile |

#### Section 2: Network Visualization
**Relationship Graph:**
- Interactive network diagram
- Nodes: Organizations
- Edges: Relationships
- Color-coded by type
- Click to explore

**Network Stats:**
- Total Relationships: 4,520
- Avg Relationships per Org: 3.6
- Highly Connected: Top 20 hubs
- Isolated: 45 orgs (need attention)

#### Section 3: Transaction Monitoring
**Real-Time Activity:**
- Orders Created: Live feed
- Payments Processed: ₹X Cr today
- Service Tickets: Active count
- RFQs Submitted: Today's count

**Top Transactions:**
| Transaction | Type | Buyer | Seller | Value | Status | Time |
|-------------|------|-------|--------|-------|--------|------|
| #ORD-1234 | Purchase Order | Apollo Delhi | Siemens | ₹1.2 Cr | Confirmed | 5 mins ago |
| #RFQ-5678 | RFQ | Max Hospital | Multiple | ₹80 L | 5 Quotes Received | 1 hr ago |

#### Section 4: Compliance & Verification
**Pending Verifications:**
- Documents to Review: 23
- Certifications Expiring: 12 (next 30 days)
- Address Verifications: 8
- Bank Details Verification: 5

**Compliance Dashboard:**
| Organization | Missing Documents | Expired Certs | Action Needed | Priority |
|--------------|-------------------|---------------|---------------|----------|
| NewMed Dist | GST Certificate | None | Upload GST | High |
| ABC Hospital | None | ISO Cert | Renew Cert | Medium |

#### Section 5: Platform Health
**System Metrics:**
- API Response Time: 245 ms (p95)
- Database Queries: 1,245/sec
- Active Sessions: 2,340
- Error Rate: 0.02%

**Service Status:**
| Service | Status | Uptime (30d) | Last Incident | Actions |
|---------|--------|--------------|---------------|---------|
| Equipment Registry | ✓ Healthy | 99.98% | None | Monitor |
| Organizations | ✓ Healthy | 99.95% | None | Monitor |
| RFQ Module | ⚠️ Slow | 99.80% | Yesterday | Investigate |

#### Section 6: Analytics & Insights
**Growth Metrics:**
- New Organizations: +45 this month
- GMV Growth: +18% MoM
- User Engagement: +24% MoM
- Platform Commission: ₹9.8 Cr

**Popular Categories:**
- Imaging Equipment: 35% of GMV
- Patient Monitoring: 25% of GMV
- Lab Equipment: 20% of GMV
- Others: 20%

---

## 6. Use Cases & Workflows

### 6.1 Use Case: New Equipment Purchase

**Scenario:** Apollo Delhi needs a new CT Scanner

```
Step 1: Create RFQ
  - Apollo creates RFQ on platform
  - Specifies: CT Scanner, specifications, budget
  - System identifies potential vendors:
    * Siemens (Direct manufacturer)
    * GE Healthcare (Direct)
    * MedEquip (Distributor - multi-brand)
    * City Medical (Dealer - local)

Step 2: AI Routing
  - Check Apollo's existing relationships
  - Check vendor capabilities
  - Check territory coverage
  - Route to: Siemens, GE, Philips (through MedEquip)

Step 3: Quote Submission
  - Siemens: ₹1.5 Cr, 7 days delivery
  - GE: ₹1.4 Cr, 10 days delivery
  - Philips (via MedEquip): ₹1.35 Cr, 14 days delivery

Step 4: AI Advisory
  - Compare quotes
  - Show market benchmarks
  - Highlight pros/cons
  - Negotiation tips

Step 5: Apollo Selects
  - Chooses Siemens (best brand, good price, fast delivery)
  - Creates PO on platform
  - Auto-syncs to Apollo's SAP

Step 6: Execution
  - Siemens confirms order
  - Payment terms: 30% advance, 70% on installation
  - Delivery tracking: Live updates
  - Installation scheduled
  - Equipment auto-registered with QR code
  - AMC contract created

Step 7: Ongoing Service
  - QR code on equipment
  - Scan for service requests
  - Auto-routes to Siemens service team
  - Service history tracked
```

---

### 6.2 Use Case: Multi-Manufacturer Dealer Operations

**Scenario:** City Medical (Dealer) daily operations

```
Morning: Stock Check
  - System shows inventory from multiple suppliers:
    * MedEquip: 50 SKUs from Siemens, GE, Philips
    * HealthTech: 30 SKUs from Medtronic, Abbott
  - Low stock alerts:
    * Siemens X-Ray tube: Only 2 left
    * GE Monitor parts: 5 left

Action: Reorder
  - Auto-suggests reorder quantities
  - Checks credit limits:
    * MedEquip: ₹38 L used / ₹50 L limit → OK
    * HealthTech: ₹18 L used / ₹20 L limit → Near limit
  - Creates orders:
    * MedEquip: ₹5 L order (Siemens + GE parts)
    * HealthTech: ₹1 L order (small items only)

Mid-Day: Customer Inquiry
  - Hospital calls: "Need patient monitors, 20 units"
  - Dealer checks:
    * GE Monitor: Available via MedEquip, ₹2.2 L each
    * Philips Monitor: Available via MedEquip, ₹2.4 L each
    * Siemens Monitor: Out of stock
  - Dealer quotes both options
  - Customer selects GE
  - Dealer creates order to MedEquip

Afternoon: Service Call
  - Apollo Hospital calls: CT Scanner issue
  - Dealer checks:
    * Equipment: Siemens CT (sold 2 years ago)
    * AMC: Active with dealer
    * Issue: Calibration error
  - Dealer dispatches engineer
  - Engineer fixes issue using parts from inventory
  - Updates service record
  - Customer charged ₹8,000 (covered under AMC)

Evening: Financial Review
  - Today's sales: ₹15 L
  - Orders placed to suppliers: ₹6 L
  - Collections: ₹10 L
  - Outstanding: ₹28 L
  - Payment due to MedEquip: ₹8 L (3 days)
  - Action: Schedule payment

Multi-Brand Intelligence:
  - GE products selling faster this month
  - Siemens has better margins
  - Philips has longest delivery times
  - System suggests: Push Siemens products for better profit
```

---

### 6.3 Use Case: Distributor Territory Expansion

**Scenario:** MedEquip wants to expand to East India

```
Step 1: Analysis
  - Current: North India (exclusive for Siemens)
  - Opportunity: East India (no exclusive distributor)
  - Market size: ₹200 Cr/year potential
  - Competition: 3 distributors (non-exclusive)

Step 2: Proposal to Siemens
  - MedEquip proposes exclusive distributorship for East India
  - Shows track record in North India
  - Commits: ₹50 Cr annual target

Step 3: Siemens Reviews
  - Dashboard shows MedEquip performance:
    * North India: 84% of target achieved
    * 75 active dealers
    * 4.5/5.0 rating
    * Growing steadily
  - Decision: Approve with conditions

Step 4: Territory Setup
  - Platform creates new territory: "East India"
  - States: West Bengal, Odisha, Bihar, Jharkhand
  - Assignment: MedEquip (exclusive for Siemens)
  - Relationship created:
    * Type: EXCLUSIVE_DISTRIBUTOR
    * Territory: East India
    * Products: All Siemens
    * Commission: 10-14% (tiered)
    * Target: ₹50 Cr/year
    * Duration: 3 years

Step 5: Execution
  - MedEquip opens warehouse in Kolkata
  - Registers facility on platform
  - Starts recruiting dealers
  - Platform shows available dealers in region
  - MedEquip onboards 30 dealers in 6 months

Step 6: Monitoring
  - Siemens dashboard shows East India performance
  - MedEquip dashboard shows dealer network
  - AI suggests: Focus on Odisha (untapped market)
  - Quarterly reviews automated
```

---

## 7. Technical Implementation

### 7.1 Database Schema

```sql
-- Core Organizations
CREATE TABLE organizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  display_name TEXT,
  org_type TEXT NOT NULL,
  sub_type TEXT,
  status TEXT NOT NULL DEFAULT 'active',
  verified BOOLEAN DEFAULT false,
  
  -- Legal
  legal_entity_name TEXT,
  registration_number TEXT,
  tax_id TEXT,
  incorporation_date DATE,
  
  -- Business
  year_established INT,
  annual_turnover NUMERIC(18,2),
  employee_count INT,
  industry_segments TEXT[],
  
  -- Digital
  website TEXT,
  logo_url TEXT,
  
  -- System
  external_refs JSONB,
  metadata JSONB,
  tenant_id TEXT NOT NULL,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  created_by TEXT,
  
  CONSTRAINT chk_org_type CHECK (org_type IN (
    'manufacturer', 'distributor', 'dealer', 'supplier',
    'hospital', 'laboratory', 'diagnostic_center', 'clinic',
    'service_provider', 'logistics_partner', 'insurance_provider',
    'government_body', 'other'
  ))
);

CREATE INDEX idx_org_type ON organizations(org_type);
CREATE INDEX idx_org_status ON organizations(status);
CREATE INDEX idx_org_verified ON organizations(verified);
CREATE INDEX idx_org_tenant ON organizations(tenant_id);

-- Facilities
CREATE TABLE organization_facilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  
  facility_name TEXT NOT NULL,
  facility_code TEXT UNIQUE,
  facility_type TEXT NOT NULL,
  
  -- Address
  address JSONB NOT NULL,
  geo_location POINT,
  
  -- Operations
  capacity TEXT,
  operational_hours JSONB,
  services_offered TEXT[],
  equipment_types TEXT[],
  
  -- Coverage
  service_radius_km INT,
  coverage_pincodes TEXT[],
  coverage_states TEXT[],
  
  status TEXT DEFAULT 'active',
  operational_since DATE,
  
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_facility_org ON organization_facilities(org_id);
CREATE INDEX idx_facility_type ON organization_facilities(facility_type);
CREATE INDEX idx_facility_location ON organization_facilities USING GIST(geo_location);

-- Relationships
CREATE TABLE organization_relationships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  parent_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  child_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  
  relationship_type TEXT NOT NULL,
  relationship_status TEXT DEFAULT 'active',
  
  -- Terms
  start_date DATE NOT NULL,
  end_date DATE,
  auto_renew BOOLEAN DEFAULT false,
  exclusive BOOLEAN DEFAULT false,
  
  -- Territory
  territory_id UUID REFERENCES territories(id),
  
  -- Commercial
  commission_percentage NUMERIC(5,2),
  payment_terms JSONB,
  credit_limit NUMERIC(18,2),
  
  -- Performance
  annual_target NUMERIC(18,2),
  performance_tier TEXT,
  priority_level INT,
  
  -- Legal
  contract_reference TEXT,
  
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_rel_type CHECK (relationship_type IN (
    'authorized_distributor', 'exclusive_distributor', 'regional_distributor',
    'authorized_dealer', 'service_partner', 'dealer_network', 'sub_distributor',
    'amc_provider', 'spare_parts_supplier', 'strategic_partner', 'oem_partner',
    'direct_buyer', 'institutional_buyer', 'logistics_partner', 'financing_partner'
  )),
  
  CONSTRAINT chk_no_self_rel CHECK (parent_org_id != child_org_id)
);

CREATE INDEX idx_rel_parent ON organization_relationships(parent_org_id);
CREATE INDEX idx_rel_child ON organization_relationships(child_org_id);
CREATE INDEX idx_rel_type ON organization_relationships(relationship_type);
CREATE INDEX idx_rel_status ON organization_relationships(relationship_status);

-- Territories
CREATE TABLE territories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  
  coverage_type TEXT NOT NULL,
  states TEXT[],
  cities TEXT[],
  districts TEXT[],
  pincodes TEXT[],
  custom_boundaries JSONB, -- GeoJSON
  
  parent_territory_id UUID REFERENCES territories(id),
  assigned_to_org_id UUID REFERENCES organizations(id),
  
  estimated_market_size NUMERIC(18,2),
  metadata JSONB,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_territory_org ON territories(assigned_to_org_id);
CREATE INDEX idx_territory_parent ON territories(parent_territory_id);

-- Contact Persons
CREATE TABLE contact_persons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  
  name TEXT NOT NULL,
  designation TEXT,
  department TEXT,
  
  email TEXT NOT NULL,
  primary_phone TEXT NOT NULL,
  alternate_phone TEXT,
  whatsapp_number TEXT,
  
  is_primary BOOLEAN DEFAULT false,
  can_approve_orders BOOLEAN DEFAULT false,
  can_raise_tickets BOOLEAN DEFAULT false,
  
  preferred_contact_method TEXT,
  active BOOLEAN DEFAULT true,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_contact_org ON contact_persons(org_id);
CREATE INDEX idx_contact_primary ON contact_persons(org_id, is_primary) WHERE is_primary = true;

-- Certifications
CREATE TABLE organization_certifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  facility_id UUID REFERENCES organization_facilities(id),
  
  certification_type TEXT NOT NULL,
  certification_number TEXT,
  issued_by TEXT,
  issue_date DATE,
  expiry_date DATE,
  status TEXT DEFAULT 'active',
  
  certificate_document_url TEXT,
  verification_url TEXT,
  
  scope TEXT,
  applicable_products TEXT[],
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_cert_org ON organization_certifications(org_id);
CREATE INDEX idx_cert_expiry ON organization_certifications(expiry_date) WHERE status = 'active';
```

---

### 7.2 API Endpoints

```typescript
// Organizations API
GET    /api/v1/organizations
GET    /api/v1/organizations/:id
POST   /api/v1/organizations
PATCH  /api/v1/organizations/:id
DELETE /api/v1/organizations/:id

// Filters
GET    /api/v1/organizations?org_type=manufacturer
GET    /api/v1/organizations?status=active&verified=true
GET    /api/v1/organizations?search=siemens

// Facilities
GET    /api/v1/organizations/:id/facilities
POST   /api/v1/organizations/:id/facilities
PATCH  /api/v1/facilities/:id
DELETE /api/v1/facilities/:id

// Relationships
GET    /api/v1/organizations/:id/relationships
POST   /api/v1/relationships
PATCH  /api/v1/relationships/:id
DELETE /api/v1/relationships/:id

// Specific relationship queries
GET    /api/v1/organizations/:id/distributors
GET    /api/v1/organizations/:id/dealers
GET    /api/v1/organizations/:id/manufacturers
GET    /api/v1/organizations/:id/service-providers

// Network visualization
GET    /api/v1/organizations/:id/network?depth=2

// Territories
GET    /api/v1/territories
POST   /api/v1/territories
PATCH  /api/v1/territories/:id
GET    /api/v1/territories/:id/organizations

// Contact Persons
GET    /api/v1/organizations/:id/contacts
POST   /api/v1/organizations/:id/contacts
PATCH  /api/v1/contacts/:id

// Certifications
GET    /api/v1/organizations/:id/certifications
POST   /api/v1/organizations/:id/certifications
PATCH  /api/v1/certifications/:id

// Dashboard Data
GET    /api/v1/dashboard/manufacturer
GET    /api/v1/dashboard/distributor
GET    /api/v1/dashboard/dealer
GET    /api/v1/dashboard/hospital
GET    /api/v1/dashboard/service-provider
GET    /api/v1/dashboard/admin

// Analytics
GET    /api/v1/analytics/network-stats
GET    /api/v1/analytics/relationship-graph
GET    /api/v1/analytics/territory-coverage
GET    /api/v1/analytics/performance-metrics
```

---

## 8. Next Steps

### Phase 1: Database & Backend (Week 1)
1. Create all database tables
2. Add seed data for testing
3. Implement Organizations API
4. Implement Facilities API
5. Implement Relationships API
6. Test with real scenarios

### Phase 2: Frontend Core (Week 2)
1. Organizations list page
2. Organization detail page
3. Facilities management
4. Relationships management
5. Multi-select filters

### Phase 3: Dashboards (Week 3)
1. Manufacturer dashboard
2. Distributor dashboard
3. Dealer dashboard
4. Hospital dashboard
5. Service provider dashboard
6. Admin dashboard

### Phase 4: Advanced Features (Week 4)
1. Network visualization
2. Territory management UI
3. Performance analytics
4. AI recommendations
5. Reporting system

---

**Status:** 📝 DETAILED DESIGN COMPLETE  
**Ready for:** Technical Review & Implementation Planning  
**Next:** Get your feedback and approval to proceed

