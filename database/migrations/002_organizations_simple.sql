-- ============================================================================
-- SIMPLE ORGANIZATIONS SCHEMA - No Transaction (for Incremental Execution)
-- ============================================================================

-- 1. Create org_relationships table first
CREATE TABLE IF NOT EXISTS org_relationships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  parent_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  child_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  rel_type TEXT NOT NULL,
  relationship_status TEXT DEFAULT 'active',
  start_date DATE DEFAULT CURRENT_DATE,
  end_date DATE,
  auto_renew BOOLEAN DEFAULT false,
  exclusive BOOLEAN DEFAULT false,
  territory_id UUID,
  commission_percentage NUMERIC(5,2),
  payment_terms JSONB,
  credit_limit NUMERIC(18,2),
  annual_target NUMERIC(18,2),
  performance_tier TEXT,
  priority_level INT,
  contract_reference TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  created_by TEXT,
  notes TEXT,
  
  CONSTRAINT chk_no_self_rel CHECK (parent_org_id != child_org_id)
);

CREATE INDEX IF NOT EXISTS idx_rel_parent ON org_relationships(parent_org_id);
CREATE INDEX IF NOT EXISTS idx_rel_child ON org_relationships(child_org_id);
CREATE INDEX IF NOT EXISTS idx_rel_type ON org_relationships(rel_type);
CREATE INDEX IF NOT EXISTS idx_rel_status ON org_relationships(relationship_status);

-- 2. Create organization_facilities
CREATE TABLE IF NOT EXISTS organization_facilities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  facility_name TEXT NOT NULL,
  facility_code TEXT UNIQUE,
  facility_type TEXT NOT NULL,
  address JSONB NOT NULL,
  geo_location POINT,
  capacity TEXT,
  operational_hours JSONB,
  services_offered TEXT[],
  equipment_types TEXT[],
  service_radius_km INT,
  coverage_pincodes TEXT[],
  coverage_states TEXT[],
  status TEXT DEFAULT 'active',
  operational_since DATE,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_facility_org ON organization_facilities(org_id);
CREATE INDEX IF NOT EXISTS idx_facility_type ON organization_facilities(facility_type);
CREATE INDEX IF NOT EXISTS idx_facility_status ON organization_facilities(status);

-- 3. Create territories
CREATE TABLE IF NOT EXISTS territories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  coverage_type TEXT NOT NULL,
  states TEXT[],
  cities TEXT[],
  districts TEXT[],
  pincodes TEXT[],
  custom_boundaries JSONB,
  parent_territory_id UUID REFERENCES territories(id),
  assigned_to_org_id UUID REFERENCES organizations(id),
  assigned_to_facility_id UUID REFERENCES organization_facilities(id),
  estimated_market_size NUMERIC(18,2),
  potential_customers INT,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_territory_org ON territories(assigned_to_org_id);
CREATE INDEX IF NOT EXISTS idx_territory_facility ON territories(assigned_to_facility_id);
CREATE INDEX IF NOT EXISTS idx_territory_parent ON territories(parent_territory_id);

-- 4. Add foreign key to org_relationships
DO $$ 
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'fk_rel_territory'
  ) THEN
    ALTER TABLE org_relationships 
      ADD CONSTRAINT fk_rel_territory 
      FOREIGN KEY (territory_id) REFERENCES territories(id);
  END IF;
END $$;

-- 5. Create contact_persons
CREATE TABLE IF NOT EXISTS contact_persons (
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
  preferred_contact_method TEXT DEFAULT 'email',
  language_preferences TEXT[],
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contact_org ON contact_persons(org_id);
CREATE INDEX IF NOT EXISTS idx_contact_primary ON contact_persons(org_id, is_primary) WHERE is_primary = true;
CREATE INDEX IF NOT EXISTS idx_contact_active ON contact_persons(active);

-- 6. Create organization_certifications
CREATE TABLE IF NOT EXISTS organization_certifications (
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

CREATE INDEX IF NOT EXISTS idx_cert_org ON organization_certifications(org_id);
CREATE INDEX IF NOT EXISTS idx_cert_facility ON organization_certifications(facility_id);
CREATE INDEX IF NOT EXISTS idx_cert_status ON organization_certifications(status);
CREATE INDEX IF NOT EXISTS idx_cert_expiry ON organization_certifications(expiry_date) WHERE status = 'active';

-- 7. Create engineers table
CREATE TABLE IF NOT EXISTS engineers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id TEXT,
  full_name TEXT NOT NULL,
  first_name TEXT,
  last_name TEXT,
  email TEXT,
  phone TEXT,
  whatsapp_number TEXT,
  org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
  org_type TEXT,
  employment_type TEXT DEFAULT 'full_time',
  joining_date DATE,
  primary_facility_id UUID REFERENCES organization_facilities(id),
  mobile_engineer BOOLEAN DEFAULT true,
  current_location POINT,
  coverage_radius_km INT,
  coverage_pincodes TEXT[],
  coverage_cities TEXT[],
  coverage_states TEXT[],
  home_region TEXT,
  status TEXT DEFAULT 'available',
  active_tickets INT DEFAULT 0,
  max_daily_tickets INT DEFAULT 5,
  working_hours JSONB,
  on_call_24x7 BOOLEAN DEFAULT false,
  total_tickets_resolved INT DEFAULT 0,
  avg_resolution_time_hours NUMERIC(6,2),
  customer_rating NUMERIC(3,2) DEFAULT 0,
  first_time_fix_rate NUMERIC(5,2) DEFAULT 0,
  preferred_contact_method TEXT DEFAULT 'phone',
  language_preferences TEXT[],
  skills TEXT[],
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_eng_org ON engineers(org_id);
CREATE INDEX IF NOT EXISTS idx_eng_org_type ON engineers(org_type);
CREATE INDEX IF NOT EXISTS idx_eng_status ON engineers(status);
CREATE INDEX IF NOT EXISTS idx_eng_facility ON engineers(primary_facility_id);
CREATE INDEX IF NOT EXISTS idx_eng_mobile ON engineers(mobile_engineer);

-- 8. Create engineer_skills
CREATE TABLE IF NOT EXISTS engineer_skills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
  skill_type TEXT NOT NULL,
  equipment_category TEXT,
  equipment_type TEXT,
  equipment_models TEXT[],
  manufacturer_id UUID REFERENCES organizations(id),
  manufacturer_name TEXT,
  manufacturer_authorized BOOLEAN DEFAULT false,
  proficiency_level TEXT DEFAULT 'intermediate',
  certification_name TEXT,
  certification_number TEXT,
  certification_authority TEXT,
  certified_date DATE,
  expiry_date DATE,
  certificate_document_url TEXT,
  can_install BOOLEAN DEFAULT false,
  can_calibrate BOOLEAN DEFAULT false,
  can_repair BOOLEAN DEFAULT true,
  can_train_users BOOLEAN DEFAULT false,
  years_of_experience INT DEFAULT 0,
  tickets_resolved_for_this_skill INT DEFAULT 0,
  verified BOOLEAN DEFAULT false,
  verified_by TEXT,
  verified_date DATE,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_skill_engineer ON engineer_skills(engineer_id);
CREATE INDEX IF NOT EXISTS idx_skill_manufacturer ON engineer_skills(manufacturer_id);
CREATE INDEX IF NOT EXISTS idx_skill_equipment_type ON engineer_skills(equipment_type);
CREATE INDEX IF NOT EXISTS idx_skill_equipment_category ON engineer_skills(equipment_category);
CREATE INDEX IF NOT EXISTS idx_skill_verified ON engineer_skills(verified);
CREATE INDEX IF NOT EXISTS idx_skill_expiry ON engineer_skills(expiry_date) WHERE expiry_date IS NOT NULL;

-- 9. Create engineer_availability
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
  UNIQUE(engineer_id, date)
);

CREATE INDEX IF NOT EXISTS idx_avail_engineer ON engineer_availability(engineer_id);
CREATE INDEX IF NOT EXISTS idx_avail_date ON engineer_availability(date);
CREATE INDEX IF NOT EXISTS idx_avail_available ON engineer_availability(available);

-- 10. Create engineer_assignments
CREATE TABLE IF NOT EXISTS engineer_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  engineer_id UUID NOT NULL REFERENCES engineers(id),
  ticket_id UUID NOT NULL REFERENCES service_tickets(id),
  equipment_id UUID NOT NULL REFERENCES equipment(id),
  assigned_by UUID,
  assigned_at TIMESTAMPTZ DEFAULT NOW(),
  assignment_type TEXT DEFAULT 'auto',
  status TEXT DEFAULT 'assigned',
  accepted_at TIMESTAMPTZ,
  en_route_at TIMESTAMPTZ,
  reached_site_at TIMESTAMPTZ,
  work_started_at TIMESTAMPTZ,
  work_completed_at TIMESTAMPTZ,
  engineer_start_location POINT,
  customer_location POINT,
  travel_distance_km NUMERIC(8,2),
  estimated_arrival TIMESTAMPTZ,
  actual_arrival TIMESTAMPTZ,
  issue_description TEXT,
  diagnosis TEXT,
  actions_taken TEXT,
  parts_used JSONB,
  customer_signature TEXT,
  customer_rating INT,
  customer_feedback TEXT,
  before_photos TEXT[],
  after_photos TEXT[],
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_assignment_engineer ON engineer_assignments(engineer_id);
CREATE INDEX IF NOT EXISTS idx_assignment_ticket ON engineer_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_assignment_equipment ON engineer_assignments(equipment_id);
CREATE INDEX IF NOT EXISTS idx_assignment_status ON engineer_assignments(status);
CREATE INDEX IF NOT EXISTS idx_assignment_date ON engineer_assignments(assigned_at);

-- 11. Enhance service_tickets table
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS assigned_engineer_id UUID REFERENCES engineers(id),
ADD COLUMN IF NOT EXISTS assignment_tier INT,
ADD COLUMN IF NOT EXISTS assignment_tier_name TEXT;

CREATE INDEX IF NOT EXISTS idx_ticket_engineer ON service_tickets(assigned_engineer_id);

-- 12. Enhance equipment table
ALTER TABLE equipment
ADD COLUMN IF NOT EXISTS manufacturer_org_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS sold_by_dealer_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS owned_by_org_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS installed_facility_id UUID REFERENCES organization_facilities(id);

CREATE INDEX IF NOT EXISTS idx_equipment_manufacturer ON equipment(manufacturer_org_id);
CREATE INDEX IF NOT EXISTS idx_equipment_dealer ON equipment(sold_by_dealer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_owner ON equipment(owned_by_org_id);
CREATE INDEX IF NOT EXISTS idx_equipment_facility ON equipment(installed_facility_id);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
