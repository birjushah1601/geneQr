-- ============================================================================
-- PHASE 1: FULL ORGANIZATIONS ARCHITECTURE - DATABASE SCHEMA
-- ============================================================================
-- Date: October 11, 2025
-- Description: Complete organizations architecture with engineer management
-- Version: 1.0.0
-- ============================================================================

-- Begin Transaction
BEGIN;

-- ============================================================================
-- 1. CORE ORGANIZATIONS TABLES (Create from scratch)
-- ============================================================================

-- Create organizations table (consolidated manufacturers + suppliers)
CREATE TABLE IF NOT EXISTS organizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  display_name TEXT,
  org_type TEXT NOT NULL,
  sub_type TEXT,
  status TEXT NOT NULL DEFAULT 'active',
  verified BOOLEAN DEFAULT false,
  verification_date TIMESTAMPTZ,
  onboarded_date TIMESTAMPTZ DEFAULT now(),
  
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
  external_ref TEXT,
  metadata JSONB,
  tenant_id TEXT DEFAULT 'default',
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_by TEXT,
  
  CONSTRAINT chk_org_type CHECK (org_type IN (
    'manufacturer', 'channel_partner', 'sub_dealer', 'supplier',
    'hospital', 'laboratory', 'diagnostic_center', 'clinic',
    'service_provider', 'logistics_partner', 'insurance_provider',
    'government_body', 'other'
  ))
);

-- Enhance existing organizations table (if it already exists, these will be ignored)
ALTER TABLE organizations 
ADD COLUMN IF NOT EXISTS display_name TEXT,
ADD COLUMN IF NOT EXISTS sub_type TEXT,
ADD COLUMN IF NOT EXISTS verified BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS verification_date TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS onboarded_date TIMESTAMPTZ DEFAULT now(),
ADD COLUMN IF NOT EXISTS legal_entity_name TEXT,
ADD COLUMN IF NOT EXISTS registration_number TEXT,
ADD COLUMN IF NOT EXISTS tax_id TEXT,
ADD COLUMN IF NOT EXISTS incorporation_date DATE,
ADD COLUMN IF NOT EXISTS year_established INT,
ADD COLUMN IF NOT EXISTS annual_turnover NUMERIC(18,2),
ADD COLUMN IF NOT EXISTS employee_count INT,
ADD COLUMN IF NOT EXISTS industry_segments TEXT[],
ADD COLUMN IF NOT EXISTS website TEXT,
ADD COLUMN IF NOT EXISTS logo_url TEXT,
ADD COLUMN IF NOT EXISTS tenant_id TEXT DEFAULT 'default',
ADD COLUMN IF NOT EXISTS created_by TEXT;

-- Update org_type constraint to include all types
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS chk_org_type;
ALTER TABLE organizations ADD CONSTRAINT chk_org_type CHECK (org_type IN (
  'manufacturer', 'channel_partner', 'sub_dealer', 'supplier',
  'hospital', 'laboratory', 'diagnostic_center', 'clinic',
  'service_provider', 'logistics_partner', 'insurance_provider',
  'government_body', 'other'
));

-- Additional indexes
CREATE INDEX IF NOT EXISTS idx_org_verified ON organizations(verified);
CREATE INDEX IF NOT EXISTS idx_org_tenant ON organizations(tenant_id);

-- ============================================================================
-- 2. ORGANIZATION FACILITIES (Multi-Location Support)
-- ============================================================================

CREATE TABLE IF NOT EXISTS organization_facilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  
  -- Facility Identity
  facility_name TEXT NOT NULL,
  facility_code TEXT UNIQUE,
  facility_type TEXT NOT NULL,
  
  -- Address (JSONB for flexibility)
  address JSONB NOT NULL,
  geo_location POINT,
  
  -- Operational Details
  capacity TEXT,
  operational_hours JSONB,
  services_offered TEXT[],
  equipment_types TEXT[],
  
  -- Coverage
  service_radius_km INT,
  coverage_pincodes TEXT[],
  coverage_states TEXT[],
  
  -- Status
  status TEXT DEFAULT 'active',
  operational_since DATE,
  
  -- System
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_facility_type CHECK (facility_type IN (
    'manufacturing_plant', 'assembly_unit', 'rnd_center',
    'warehouse', 'distribution_center', 'service_center',
    'training_center', 'sales_office', 'showroom',
    'hospital_unit', 'laboratory_unit', 'diagnostic_center', 'clinic'
  )),
  CONSTRAINT chk_facility_status CHECK (status IN ('active', 'inactive', 'under_construction'))
);

CREATE INDEX IF NOT EXISTS idx_facility_org ON organization_facilities(org_id);
CREATE INDEX IF NOT EXISTS idx_facility_type ON organization_facilities(facility_type);
CREATE INDEX IF NOT EXISTS idx_facility_status ON organization_facilities(status);
CREATE INDEX IF NOT EXISTS idx_facility_location ON organization_facilities USING GIST(geo_location);

-- ============================================================================
-- 3. ORGANIZATION RELATIONSHIPS (Enhanced)
-- ============================================================================

-- Create org_relationships table first
CREATE TABLE IF NOT EXISTS org_relationships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  parent_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  child_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  rel_type TEXT NOT NULL,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rel_parent ON org_relationships(parent_org_id);
CREATE INDEX IF NOT EXISTS idx_rel_child ON org_relationships(child_org_id);
CREATE INDEX IF NOT EXISTS idx_rel_type ON org_relationships(rel_type);

-- Enhance existing org_relationships table
ALTER TABLE org_relationships
ADD COLUMN IF NOT EXISTS relationship_status TEXT DEFAULT 'active',
ADD COLUMN IF NOT EXISTS start_date DATE DEFAULT CURRENT_DATE,
ADD COLUMN IF NOT EXISTS end_date DATE,
ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS exclusive BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS territory_id UUID,
ADD COLUMN IF NOT EXISTS commission_percentage NUMERIC(5,2),
ADD COLUMN IF NOT EXISTS payment_terms JSONB,
ADD COLUMN IF NOT EXISTS credit_limit NUMERIC(18,2),
ADD COLUMN IF NOT EXISTS annual_target NUMERIC(18,2),
ADD COLUMN IF NOT EXISTS performance_tier TEXT,
ADD COLUMN IF NOT EXISTS priority_level INT,
ADD COLUMN IF NOT EXISTS contract_reference TEXT,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS created_by TEXT,
ADD COLUMN IF NOT EXISTS notes TEXT;

-- Update rel_type constraint
ALTER TABLE org_relationships DROP CONSTRAINT IF EXISTS chk_rel_type;
ALTER TABLE org_relationships ADD CONSTRAINT chk_rel_type CHECK (rel_type IN (
  'authorized_channel_partner', 'exclusive_channel_partner', 'regional_channel_partner',
  'authorized_sub_dealer', 'service_partner', 'sub_dealer_network', 'sub_channel_partner',
  'amc_provider', 'spare_parts_supplier', 'strategic_partner', 'oem_partner',
  'direct_buyer', 'institutional_buyer', 'logistics_partner', 'financing_partner',
  'manufacturer_of', 'channel_partner_of', 'sub_dealer_of', 'supplier_of', 'partner_of'
));

-- Add constraints
ALTER TABLE org_relationships ADD CONSTRAINT chk_rel_status 
  CHECK (relationship_status IN ('active', 'inactive', 'pending', 'expired'));
ALTER TABLE org_relationships ADD CONSTRAINT chk_no_self_rel 
  CHECK (parent_org_id != child_org_id);

-- Additional indexes
CREATE INDEX IF NOT EXISTS idx_rel_status ON org_relationships(relationship_status);
CREATE INDEX IF NOT EXISTS idx_rel_territory ON org_relationships(territory_id);

-- ============================================================================
-- 4. TERRITORIES
-- ============================================================================

CREATE TABLE IF NOT EXISTS territories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  
  -- Coverage
  coverage_type TEXT NOT NULL,
  states TEXT[],
  cities TEXT[],
  districts TEXT[],
  pincodes TEXT[],
  custom_boundaries JSONB,
  
  -- Hierarchy
  parent_territory_id UUID REFERENCES territories(id),
  assigned_to_org_id UUID REFERENCES organizations(id),
  assigned_to_facility_id UUID REFERENCES organization_facilities(id),
  
  -- Market Data
  estimated_market_size NUMERIC(18,2),
  potential_customers INT,
  metadata JSONB,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_coverage_type CHECK (coverage_type IN (
    'pincode', 'city', 'district', 'state', 'region', 'custom'
  ))
);

CREATE INDEX IF NOT EXISTS idx_territory_org ON territories(assigned_to_org_id);
CREATE INDEX IF NOT EXISTS idx_territory_facility ON territories(assigned_to_facility_id);
CREATE INDEX IF NOT EXISTS idx_territory_parent ON territories(parent_territory_id);

-- Add foreign key from org_relationships
ALTER TABLE org_relationships 
  ADD CONSTRAINT fk_rel_territory 
  FOREIGN KEY (territory_id) REFERENCES territories(id);

-- ============================================================================
-- 5. CONTACT PERSONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS contact_persons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  
  -- Identity
  name TEXT NOT NULL,
  designation TEXT,
  department TEXT,
  
  -- Contact
  email TEXT NOT NULL,
  primary_phone TEXT NOT NULL,
  alternate_phone TEXT,
  whatsapp_number TEXT,
  
  -- Permissions
  is_primary BOOLEAN DEFAULT false,
  can_approve_orders BOOLEAN DEFAULT false,
  can_raise_tickets BOOLEAN DEFAULT false,
  
  -- Preferences
  preferred_contact_method TEXT DEFAULT 'email',
  language_preferences TEXT[],
  
  -- Status
  active BOOLEAN DEFAULT true,
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_contact_method CHECK (preferred_contact_method IN ('email', 'phone', 'whatsapp'))
);

CREATE INDEX IF NOT EXISTS idx_contact_org ON contact_persons(org_id);
CREATE INDEX IF NOT EXISTS idx_contact_primary ON contact_persons(org_id, is_primary) WHERE is_primary = true;
CREATE INDEX IF NOT EXISTS idx_contact_active ON contact_persons(active);

-- ============================================================================
-- 6. ORGANIZATION CERTIFICATIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS organization_certifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  facility_id UUID REFERENCES organization_facilities(id),
  
  -- Certification Details
  certification_type TEXT NOT NULL,
  certification_number TEXT,
  issued_by TEXT,
  issue_date DATE,
  expiry_date DATE,
  status TEXT DEFAULT 'active',
  
  -- Documents
  certificate_document_url TEXT,
  verification_url TEXT,
  
  -- Scope
  scope TEXT,
  applicable_products TEXT[],
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_cert_status CHECK (status IN ('active', 'expired', 'suspended'))
);

CREATE INDEX IF NOT EXISTS idx_cert_org ON organization_certifications(org_id);
CREATE INDEX IF NOT EXISTS idx_cert_facility ON organization_certifications(facility_id);
CREATE INDEX IF NOT EXISTS idx_cert_status ON organization_certifications(status);
CREATE INDEX IF NOT EXISTS idx_cert_expiry ON organization_certifications(expiry_date) 
  WHERE status = 'active';

-- ============================================================================
-- 7. ENGINEERS TABLE (Multi-Entity Support)
-- ============================================================================

-- Create engineers table from scratch
CREATE TABLE IF NOT EXISTS engineers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- Identity
  employee_id TEXT,
  full_name TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  email TEXT,
  phone TEXT,
  whatsapp_number TEXT,
  
  -- Employment
  org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  org_type TEXT,
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
  home_region TEXT,
  
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
  
  -- Legacy
  skills TEXT[],
  
  -- System
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  CONSTRAINT chk_eng_status CHECK (status IN ('available', 'on_job', 'on_leave', 'inactive')),
  CONSTRAINT chk_eng_org_type CHECK (org_type IN (
    'manufacturer', 'channel_partner', 'sub_dealer', 'hospital', 'clinic',
    'service_provider', 'laboratory', 'diagnostic_center'
  )),
  CONSTRAINT chk_eng_employment CHECK (employment_type IN ('full_time', 'part_time', 'contract', 'freelance'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_eng_org ON engineers(org_id);
CREATE INDEX IF NOT EXISTS idx_eng_org_type ON engineers(org_type);
CREATE INDEX IF NOT EXISTS idx_eng_status ON engineers(status);
CREATE INDEX IF NOT EXISTS idx_eng_location ON engineers USING GIST(current_location);
CREATE INDEX IF NOT EXISTS idx_eng_facility ON engineers(primary_facility_id);
CREATE INDEX IF NOT EXISTS idx_eng_mobile ON engineers(mobile_engineer);

-- ============================================================================
-- 7B. SERVICE TICKETS TABLE (Must be created before engineer_assignments)
-- ============================================================================

-- Create service_tickets table first (using VARCHAR to match Go code expectations)
CREATE TABLE IF NOT EXISTS service_tickets (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    tenant_id VARCHAR(255) NOT NULL,
    equipment_id VARCHAR(255),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority VARCHAR(50) DEFAULT 'medium',
    status VARCHAR(50) DEFAULT 'open',
    organization_id UUID REFERENCES organizations(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tickets_status ON service_tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON service_tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_organization ON service_tickets(organization_id);
CREATE INDEX IF NOT EXISTS idx_tickets_equipment ON service_tickets(equipment_id);

-- ============================================================================
-- 7C. EQUIPMENT TABLE (Must be created before engineer_assignments)
-- ============================================================================

-- Create equipment table first (using VARCHAR to match Go code expectations)
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    serial_number VARCHAR(200),
    model VARCHAR(200),
    manufacturer VARCHAR(200),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_equipment_status ON equipment(status);
CREATE INDEX IF NOT EXISTS idx_equipment_serial ON equipment(serial_number);

-- ============================================================================
-- 8. ENGINEER SKILLS & CERTIFICATIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS engineer_skills (
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

CREATE INDEX IF NOT EXISTS idx_skill_engineer ON engineer_skills(engineer_id);
CREATE INDEX IF NOT EXISTS idx_skill_manufacturer ON engineer_skills(manufacturer_id);
CREATE INDEX IF NOT EXISTS idx_skill_equipment_type ON engineer_skills(equipment_type);
CREATE INDEX IF NOT EXISTS idx_skill_equipment_category ON engineer_skills(equipment_category);
CREATE INDEX IF NOT EXISTS idx_skill_verified ON engineer_skills(verified);
CREATE INDEX IF NOT EXISTS idx_skill_expiry ON engineer_skills(expiry_date) 
  WHERE expiry_date IS NOT NULL;

-- ============================================================================
-- 9. ENGINEER AVAILABILITY
-- ============================================================================

CREATE TABLE IF NOT EXISTS engineer_availability (
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
  
  UNIQUE(engineer_id, date),
  
  CONSTRAINT chk_avail_reason CHECK (reason IN (
    'on_leave', 'on_training', 'sick', 'on_site', 'other'
  ))
);

CREATE INDEX IF NOT EXISTS idx_avail_engineer ON engineer_availability(engineer_id);
CREATE INDEX IF NOT EXISTS idx_avail_date ON engineer_availability(date);
CREATE INDEX IF NOT EXISTS idx_avail_available ON engineer_availability(available);

-- ============================================================================
-- 10. ENGINEER ASSIGNMENTS
-- ============================================================================

CREATE TABLE IF NOT EXISTS engineer_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  
  -- Assignment
  engineer_id UUID NOT NULL REFERENCES engineers(id),
  ticket_id VARCHAR(255) NOT NULL REFERENCES service_tickets(id),
  equipment_id VARCHAR(255) NOT NULL REFERENCES equipment(id),
  
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
  CONSTRAINT chk_assignment_type CHECK (assignment_type IN ('auto', 'manual')),
  CONSTRAINT chk_customer_rating CHECK (customer_rating BETWEEN 1 AND 5)
);

CREATE INDEX IF NOT EXISTS idx_assignment_engineer ON engineer_assignments(engineer_id);
CREATE INDEX IF NOT EXISTS idx_assignment_ticket ON engineer_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_assignment_equipment ON engineer_assignments(equipment_id);
CREATE INDEX IF NOT EXISTS idx_assignment_status ON engineer_assignments(status);
CREATE INDEX IF NOT EXISTS idx_assignment_date ON engineer_assignments(assigned_at);

-- ============================================================================
-- 11. ENHANCE SERVICE_TICKETS TABLE (table created in section 7B)
-- ============================================================================

-- Add engineer assignment fields
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS assigned_engineer_id UUID REFERENCES engineers(id),
ADD COLUMN IF NOT EXISTS assignment_tier INT,
ADD COLUMN IF NOT EXISTS assignment_tier_name TEXT;

CREATE INDEX IF NOT EXISTS idx_ticket_engineer ON service_tickets(assigned_engineer_id);

-- ============================================================================
-- 12. ENHANCE EQUIPMENT TABLE (table created in section 7C)
-- ============================================================================

-- Link equipment to organizations
ALTER TABLE equipment
ADD COLUMN IF NOT EXISTS manufacturer_org_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS sold_by_sub_dealer_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS owned_by_org_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS installed_facility_id UUID REFERENCES organization_facilities(id);

CREATE INDEX IF NOT EXISTS idx_equipment_manufacturer ON equipment(manufacturer_org_id);
CREATE INDEX IF NOT EXISTS idx_equipment_sub_dealer ON equipment(sold_by_sub_dealer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_owner ON equipment(owned_by_org_id);
CREATE INDEX IF NOT EXISTS idx_equipment_facility ON equipment(installed_facility_id);

-- Add foreign key from service_tickets to equipment (now that equipment table exists)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'service_tickets_equipment_id_fkey'
    ) THEN
        ALTER TABLE service_tickets 
        ADD CONSTRAINT service_tickets_equipment_id_fkey 
        FOREIGN KEY (equipment_id) REFERENCES equipment(id);
    END IF;
END $$;

-- ============================================================================
-- 13. CREATE VIEWS FOR COMMON QUERIES
-- ============================================================================

-- View: Organization with primary contact
CREATE OR REPLACE VIEW vw_organizations_with_contact AS
SELECT 
  o.*,
  cp.name AS primary_contact_name,
  cp.email AS primary_contact_email,
  cp.primary_phone AS primary_contact_phone
FROM organizations o
LEFT JOIN contact_persons cp ON o.id = cp.org_id AND cp.is_primary = true;

-- View: Engineers with skills summary
CREATE OR REPLACE VIEW vw_engineers_with_skills AS
SELECT 
  e.*,
  o.name AS organization_name,
  o.org_type AS organization_type,
  f.facility_name,
  COUNT(DISTINCT es.id) AS total_skills,
  COUNT(DISTINCT CASE WHEN es.manufacturer_authorized THEN es.id END) AS authorized_skills
FROM engineers e
LEFT JOIN organizations o ON e.org_id = o.id
LEFT JOIN organization_facilities f ON e.primary_facility_id = f.id
LEFT JOIN engineer_skills es ON e.id = es.engineer_id
GROUP BY e.id, o.name, o.org_type, f.facility_name;

-- View: Active assignments with details
CREATE OR REPLACE VIEW vw_active_assignments AS
SELECT 
  ea.*,
  e.full_name AS engineer_name,
  e.phone AS engineer_phone,
  st.title AS ticket_title,
  st.priority AS ticket_priority,
  eq.name AS equipment_name,
  eq.model AS equipment_model,
  o.name AS customer_org_name
FROM engineer_assignments ea
JOIN engineers e ON ea.engineer_id = e.id
JOIN service_tickets st ON ea.ticket_id = st.id
JOIN equipment eq ON ea.equipment_id = eq.id
JOIN organizations o ON st.organization_id = o.id
WHERE ea.status IN ('assigned', 'accepted', 'en_route', 'on_site');

-- ============================================================================
-- 14. TRIGGERS FOR UPDATED_AT
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_facilities_updated_at BEFORE UPDATE ON organization_facilities 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_relationships_updated_at BEFORE UPDATE ON org_relationships 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_territories_updated_at BEFORE UPDATE ON territories 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_contact_persons_updated_at BEFORE UPDATE ON contact_persons 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certifications_updated_at BEFORE UPDATE ON organization_certifications 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_engineers_updated_at BEFORE UPDATE ON engineers 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_engineer_skills_updated_at BEFORE UPDATE ON engineer_skills 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assignments_updated_at BEFORE UPDATE ON engineer_assignments 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Commit Transaction
COMMIT;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
-- Summary:
-- - Enhanced organizations table with full profile support
-- - Added organization_facilities for multi-location support
-- - Enhanced org_relationships with business terms
-- - Added territories for geographic management
-- - Added contact_persons for organization contacts
-- - Added organization_certifications for compliance tracking
-- - Enhanced engineers table with full profile
-- - Added engineer_skills for skill-based routing
-- - Added engineer_availability for scheduling
-- - Added engineer_assignments for work tracking
-- - Enhanced existing equipment and service_tickets tables
-- - Created useful views for common queries
-- - Added triggers for automatic timestamp updates
-- ============================================================================
