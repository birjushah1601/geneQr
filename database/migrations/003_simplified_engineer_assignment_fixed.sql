-- ============================================================================
-- SIMPLIFIED ENGINEER ASSIGNMENT SYSTEM - MVP (FIXED VERSION)
-- ============================================================================
-- Purpose: Simplified, configuration-driven engineer assignment
-- Approach: Organization-centric, loosely coupled, extensible
-- Date: 2025-11-21
-- Fix: Changed equipment_id from UUID to VARCHAR to match equipment table
-- ============================================================================

-- ============================================================================
-- 1. SIMPLIFY ENGINEERS TABLE (Add engineer_level)
-- ============================================================================
ALTER TABLE engineers 
ADD COLUMN IF NOT EXISTS engineer_level INT DEFAULT 1;

COMMENT ON COLUMN engineers.engineer_level IS 'Engineer level: 1=Junior (basic repairs), 2=Senior (most repairs), 3=Expert (complex/specialized)';

CREATE INDEX IF NOT EXISTS idx_engineer_level ON engineers(engineer_level);

-- ============================================================================
-- 2. ENGINEER EQUIPMENT TYPES (What Engineers Can Repair)
-- ============================================================================
CREATE TABLE IF NOT EXISTS engineer_equipment_types (\n  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
  
  -- Equipment they can service (simple manufacturer + category approach)
  manufacturer_name TEXT NOT NULL,  -- e.g., "Siemens Healthineers", "GE Healthcare"
  equipment_category TEXT NOT NULL, -- e.g., "MRI", "CT Scanner", "X-Ray", "Ultrasound"
  model_pattern TEXT,               -- Optional: e.g., "Magnetom%" for specific models, NULL = all models
  
  -- Certification (optional for now)
  is_certified BOOLEAN DEFAULT false,
  certification_number TEXT,
  certification_expiry DATE,
  
  -- Metadata
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Unique constraint to prevent duplicates
  UNIQUE(engineer_id, manufacturer_name, equipment_category)
);

CREATE INDEX IF NOT EXISTS idx_engineer_equip_engineer ON engineer_equipment_types(engineer_id);
CREATE INDEX IF NOT EXISTS idx_engineer_equip_mfr ON engineer_equipment_types(manufacturer_name);
CREATE INDEX IF NOT EXISTS idx_engineer_equip_cat ON engineer_equipment_types(equipment_category);
CREATE INDEX IF NOT EXISTS idx_engineer_equip_certified ON engineer_equipment_types(is_certified) WHERE is_certified = true;

COMMENT ON TABLE engineer_equipment_types IS 'Defines which equipment types an engineer can service (manufacturer + category based)';

-- ============================================================================
-- 3. EQUIPMENT SERVICE CONFIGURATION (Who Services What)
-- ============================================================================
-- Drop table if it exists from previous failed migration
DROP TABLE IF EXISTS equipment_service_config CASCADE;

CREATE TABLE equipment_service_config (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  equipment_id VARCHAR(255) NOT NULL REFERENCES equipment(id) ON DELETE CASCADE UNIQUE,
  
  -- Service hierarchy (ordered by priority)
  -- These are the organizations that should be tried in order
  primary_service_org_id UUID REFERENCES organizations(id),    -- Usually manufacturer (Tier 1)
  secondary_service_org_id UUID REFERENCES organizations(id),  -- Usually dealer who sold it (Tier 2)
  tertiary_service_org_id UUID REFERENCES organizations(id),   -- Usually distributor (Tier 3)
  fallback_service_org_id UUID REFERENCES organizations(id),   -- Hospital in-house BME team (Tier 4/5)
  
  -- Warranty/AMC status (affects who gets priority)
  warranty_provider_org_id UUID REFERENCES organizations(id),
  warranty_active BOOLEAN DEFAULT false,
  warranty_start_date DATE,
  warranty_end_date DATE,
  
  amc_provider_org_id UUID REFERENCES organizations(id),
  amc_active BOOLEAN DEFAULT false,
  amc_start_date DATE,
  amc_end_date DATE,
  amc_contract_number TEXT,
  
  -- Minimum engineer level required for this equipment
  min_engineer_level INT DEFAULT 1,  -- 1=any engineer, 2=senior+, 3=expert only
  
  -- Notes
  service_notes TEXT,
  special_instructions TEXT,
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  created_by TEXT,
  updated_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_equip_config_equipment ON equipment_service_config(equipment_id);
CREATE INDEX IF NOT EXISTS idx_equip_config_primary ON equipment_service_config(primary_service_org_id);
CREATE INDEX IF NOT EXISTS idx_equip_config_secondary ON equipment_service_config(secondary_service_org_id);
CREATE INDEX IF NOT EXISTS idx_equip_config_warranty_active ON equipment_service_config(warranty_active) WHERE warranty_active = true;
CREATE INDEX IF NOT EXISTS idx_equip_config_amc_active ON equipment_service_config(amc_active) WHERE amc_active = true;

COMMENT ON TABLE equipment_service_config IS 'Configuration for service routing per equipment - who should service it and in what order';
COMMENT ON COLUMN equipment_service_config.primary_service_org_id IS 'First choice for service (usually manufacturer)';
COMMENT ON COLUMN equipment_service_config.secondary_service_org_id IS 'Second choice (usually dealer)';
COMMENT ON COLUMN equipment_service_config.tertiary_service_org_id IS 'Third choice (usually distributor)';
COMMENT ON COLUMN equipment_service_config.fallback_service_org_id IS 'Last resort (usually hospital in-house team)';
COMMENT ON COLUMN equipment_service_config.min_engineer_level IS 'Minimum engineer level required: 1=Junior, 2=Senior, 3=Expert';

-- ============================================================================
-- 4. ENHANCE SERVICE_TICKETS TABLE
-- ============================================================================
-- Add columns for assignment tracking (if not already exist)
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS assigned_org_id UUID REFERENCES organizations(id),
ADD COLUMN IF NOT EXISTS assignment_tier INT,
ADD COLUMN IF NOT EXISTS assignment_tier_name TEXT,
ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_ticket_assigned_org ON service_tickets(assigned_org_id);
CREATE INDEX IF NOT EXISTS idx_ticket_assignment_tier ON service_tickets(assignment_tier);

COMMENT ON COLUMN service_tickets.assigned_org_id IS 'Organization whose engineer was assigned';
COMMENT ON COLUMN service_tickets.assignment_tier IS 'Which tier was used: 1=Primary, 2=Secondary, 3=Tertiary, 4=Fallback';
COMMENT ON COLUMN service_tickets.assignment_tier_name IS 'Human-readable tier name: warranty, amc, manufacturer, dealer, distributor, in_house';

-- ============================================================================
-- 5. HELPER FUNCTIONS
-- ============================================================================

-- Function to get eligible service organizations for an equipment
CREATE OR REPLACE FUNCTION get_eligible_service_orgs(p_equipment_id VARCHAR)
RETURNS TABLE(
  org_id UUID,
  org_name TEXT,
  org_type TEXT,
  tier INT,
  tier_name TEXT,
  reason TEXT
) AS $$
BEGIN
  RETURN QUERY
  WITH config AS (
    SELECT 
      CASE 
        WHEN warranty_active THEN warranty_provider_org_id
        WHEN amc_active THEN amc_provider_org_id
        ELSE primary_service_org_id
      END as priority_org_id,
      CASE 
        WHEN warranty_active THEN 'warranty'
        WHEN amc_active THEN 'amc'
        ELSE 'standard'
      END as priority_reason,
      secondary_service_org_id,
      tertiary_service_org_id,
      fallback_service_org_id,
      warranty_active,
      amc_active
    FROM equipment_service_config
    WHERE equipment_id = p_equipment_id
  )
  SELECT 
    o.id,
    o.name,
    o.org_type,
    1 as tier,
    CASE 
      WHEN c.warranty_active THEN 'warranty'
      WHEN c.amc_active THEN 'amc'
      ELSE 'primary'
    END as tier_name,
    c.priority_reason
  FROM config c
  JOIN organizations o ON o.id = c.priority_org_id
  WHERE c.priority_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    2 as tier,
    'secondary' as tier_name,
    'authorized_dealer' as reason
  FROM config c
  JOIN organizations o ON o.id = c.secondary_service_org_id
  WHERE c.secondary_service_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    3 as tier,
    'tertiary' as tier_name,
    'distributor' as reason
  FROM config c
  JOIN organizations o ON o.id = c.tertiary_service_org_id
  WHERE c.tertiary_service_org_id IS NOT NULL
  
  UNION ALL
  
  SELECT 
    o.id,
    o.name,
    o.org_type,
    4 as tier,
    'fallback' as tier_name,
    'in_house' as reason
  FROM config c
  JOIN organizations o ON o.id = c.fallback_service_org_id
  WHERE c.fallback_service_org_id IS NOT NULL
  
  ORDER BY tier;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_eligible_service_orgs IS 'Returns list of organizations eligible to service an equipment, in priority order';

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
