-- Migration: 017-create-engineer-expertise.sql
-- Description: Engineer Expertise & Service Configuration
-- Ticket: T2B.2
-- Date: 2025-11-16
-- 
-- This migration creates:
-- 1. engineer_equipment_expertise - Engineer skills per equipment (L1/L2/L3)
-- 2. manufacturer_service_config - Who handles service (manufacturer, client, dealer)
-- 3. engineer_certifications - Formal certifications tracking
--
-- Purpose:
-- - Track engineer expertise levels for each equipment type
-- - Configure who handles service for each equipment (manufacturer vs client)
-- - Support L1/L2/L3 support level filtering for assignment
-- - Enable AI to match engineers to tickets based on skills

-- =====================================================================
-- 1. ENGINEER EQUIPMENT EXPERTISE
-- =====================================================================

CREATE TABLE IF NOT EXISTS engineer_equipment_expertise (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    engineer_id UUID NOT NULL REFERENCES engineers(id),
    equipment_catalog_id UUID NOT NULL REFERENCES equipment_catalog(id),
    manufacturer_id UUID NOT NULL REFERENCES organizations(id),
    
    -- Support Level
    support_level TEXT NOT NULL,                -- 'L1', 'L2', 'L3'
    
    -- Certification & Training
    certified BOOLEAN DEFAULT false,
    certification_number TEXT,
    certification_date DATE,
    certification_expiry DATE,
    training_completed BOOLEAN DEFAULT false,
    training_date DATE,
    training_provider TEXT,
    
    -- Experience
    years_experience INT DEFAULT 0,
    total_repairs_completed INT DEFAULT 0,
    successful_repairs INT DEFAULT 0,
    average_resolution_hours NUMERIC(5,2),
    
    -- Performance Metrics
    first_time_fix_rate NUMERIC(5,2),          -- Percentage (0-100)
    customer_satisfaction_avg NUMERIC(3,2),    -- Average rating (0-5)
    escalation_rate NUMERIC(5,2),              -- Percentage of tickets escalated
    
    -- Specializations
    specializations TEXT[],                     -- ['Remote Diagnosis', 'Field Repair', 'Installation']
    can_do_remote BOOLEAN DEFAULT true,
    can_do_onsite BOOLEAN DEFAULT true,
    can_do_installation BOOLEAN DEFAULT false,
    can_do_calibration BOOLEAN DEFAULT false,
    
    -- Availability
    is_active BOOLEAN DEFAULT true,
    available_for_assignment BOOLEAN DEFAULT true,
    max_concurrent_tickets INT DEFAULT 5,
    
    -- Validity Period (for temporary assignments)
    effective_from DATE DEFAULT CURRENT_DATE,
    effective_to DATE,                          -- NULL = permanent
    
    -- Notes
    notes TEXT,
    internal_notes TEXT,                        -- Not visible to engineer
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_support_level CHECK (support_level IN ('L1', 'L2', 'L3')),
    CONSTRAINT chk_years_experience CHECK (years_experience >= 0),
    CONSTRAINT chk_total_repairs CHECK (total_repairs_completed >= 0),
    CONSTRAINT chk_successful_repairs CHECK (successful_repairs >= 0 AND successful_repairs <= total_repairs_completed),
    CONSTRAINT chk_first_time_fix CHECK (first_time_fix_rate IS NULL OR (first_time_fix_rate >= 0 AND first_time_fix_rate <= 100)),
    CONSTRAINT chk_customer_satisfaction CHECK (customer_satisfaction_avg IS NULL OR (customer_satisfaction_avg >= 0 AND customer_satisfaction_avg <= 5)),
    CONSTRAINT chk_escalation_rate CHECK (escalation_rate IS NULL OR (escalation_rate >= 0 AND escalation_rate <= 100)),
    CONSTRAINT chk_max_tickets CHECK (max_concurrent_tickets > 0),
    CONSTRAINT chk_effective_dates CHECK (effective_to IS NULL OR effective_to >= effective_from)
);

-- Unique constraint: One engineer-equipment-manufacturer combination
CREATE UNIQUE INDEX idx_engineer_expertise_unique 
ON engineer_equipment_expertise(engineer_id, equipment_catalog_id, manufacturer_id);

-- Indexes for common queries
CREATE INDEX idx_engineer_expertise_engineer ON engineer_equipment_expertise(engineer_id);
CREATE INDEX idx_engineer_expertise_equipment ON engineer_equipment_expertise(equipment_catalog_id);
CREATE INDEX idx_engineer_expertise_manufacturer ON engineer_equipment_expertise(manufacturer_id);
CREATE INDEX idx_engineer_expertise_level ON engineer_equipment_expertise(support_level);
CREATE INDEX idx_engineer_expertise_active ON engineer_equipment_expertise(is_active) WHERE is_active = true;
CREATE INDEX idx_engineer_expertise_available ON engineer_equipment_expertise(available_for_assignment) WHERE available_for_assignment = true;
CREATE INDEX idx_engineer_expertise_certified ON engineer_equipment_expertise(certified) WHERE certified = true;

-- Composite index for assignment queries (most common)
CREATE INDEX idx_engineer_expertise_assignment ON engineer_equipment_expertise(
    equipment_catalog_id, support_level, is_active, available_for_assignment
) WHERE is_active = true AND available_for_assignment = true;

-- GIN index for specializations array
CREATE INDEX idx_engineer_expertise_specializations ON engineer_equipment_expertise USING GIN (specializations);

COMMENT ON TABLE engineer_equipment_expertise IS 'Engineer skills and expertise per equipment type with L1/L2/L3 support levels';
COMMENT ON COLUMN engineer_equipment_expertise.support_level IS 'L1=Basic/Remote, L2=Advanced/Field, L3=Expert/Complex';
COMMENT ON COLUMN engineer_equipment_expertise.first_time_fix_rate IS 'Percentage of issues fixed on first attempt';
COMMENT ON COLUMN engineer_equipment_expertise.escalation_rate IS 'Percentage of tickets escalated to higher level';

-- =====================================================================
-- 2. MANUFACTURER SERVICE CONFIGURATION
-- =====================================================================

CREATE TABLE IF NOT EXISTS manufacturer_service_config (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Scope
    manufacturer_id UUID NOT NULL REFERENCES organizations(id),
    equipment_catalog_id UUID REFERENCES equipment_catalog(id),  -- NULL = applies to all equipment from manufacturer
    
    -- Service Provider
    service_provider_type TEXT NOT NULL,        -- 'manufacturer', 'client', 'dealer', 'distributor'
    service_provider_org_id UUID REFERENCES organizations(id),
    
    -- Service Details
    service_scope TEXT[],                       -- ['installation', 'repair', 'maintenance', 'calibration', 'training']
    warranty_service_only BOOLEAN DEFAULT false,
    post_warranty_service BOOLEAN DEFAULT true,
    
    -- Requirements
    requires_oem_parts BOOLEAN DEFAULT false,   -- Must use OEM parts
    requires_certified_engineer BOOLEAN DEFAULT false,
    requires_manufacturer_approval BOOLEAN DEFAULT false,
    warranty_void_if_third_party BOOLEAN DEFAULT false,
    
    -- SLA
    sla_response_hours INT,                     -- Response time SLA
    sla_resolution_hours INT,                   -- Resolution time SLA
    sla_priority TEXT,                          -- 'critical', 'high', 'medium', 'low'
    
    -- Geographic Coverage
    coverage_regions TEXT[],                    -- ['North India', 'South India', 'Pan India']
    coverage_cities TEXT[],                     -- Specific cities if not regional
    
    -- Priority (for hierarchy)
    priority INT DEFAULT 5,                     -- Equipment-specific (10) > Manufacturer (5) > Default (1)
    
    -- Contract Details
    contract_ref TEXT,
    contract_start_date DATE,
    contract_end_date DATE,
    
    -- Validity
    effective_from DATE DEFAULT CURRENT_DATE,
    effective_to DATE,                          -- NULL = currently active
    
    -- Additional Terms
    terms JSONB DEFAULT '{}',                   -- Additional configuration
    notes TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_service_provider_type CHECK (service_provider_type IN (
        'manufacturer', 'client', 'dealer', 'distributor', 'third_party'
    )),
    CONSTRAINT chk_sla_priority CHECK (sla_priority IS NULL OR sla_priority IN (
        'critical', 'high', 'medium', 'low'
    )),
    CONSTRAINT chk_priority CHECK (priority >= 1 AND priority <= 10),
    CONSTRAINT chk_sla_response CHECK (sla_response_hours IS NULL OR sla_response_hours > 0),
    CONSTRAINT chk_sla_resolution CHECK (sla_resolution_hours IS NULL OR sla_resolution_hours > 0),
    CONSTRAINT chk_effective_dates CHECK (effective_to IS NULL OR effective_to >= effective_from),
    CONSTRAINT chk_contract_dates CHECK (
        (contract_start_date IS NULL AND contract_end_date IS NULL) OR
        (contract_start_date IS NOT NULL AND (contract_end_date IS NULL OR contract_end_date >= contract_start_date))
    )
);

-- Indexes
CREATE INDEX idx_manufacturer_service_manufacturer ON manufacturer_service_config(manufacturer_id);
CREATE INDEX idx_manufacturer_service_equipment ON manufacturer_service_config(equipment_catalog_id);
CREATE INDEX idx_manufacturer_service_provider ON manufacturer_service_config(service_provider_org_id);
CREATE INDEX idx_manufacturer_service_priority ON manufacturer_service_config(manufacturer_id, priority DESC);
CREATE INDEX idx_manufacturer_service_active ON manufacturer_service_config(effective_to) WHERE effective_to IS NULL;

-- Composite index for service provider lookup (most common query)
CREATE INDEX idx_manufacturer_service_lookup ON manufacturer_service_config(
    manufacturer_id, equipment_catalog_id, priority DESC
) WHERE effective_to IS NULL;

-- GIN indexes for arrays
CREATE INDEX idx_manufacturer_service_scope ON manufacturer_service_config USING GIN (service_scope);
CREATE INDEX idx_manufacturer_service_regions ON manufacturer_service_config USING GIN (coverage_regions);

COMMENT ON TABLE manufacturer_service_config IS 'Configure who handles service for equipment (manufacturer vs client vs dealer)';
COMMENT ON COLUMN manufacturer_service_config.equipment_catalog_id IS 'NULL = applies to all equipment from manufacturer';
COMMENT ON COLUMN manufacturer_service_config.priority IS 'Equipment-specific (10) > Manufacturer-level (5) > Default (1)';
COMMENT ON COLUMN manufacturer_service_config.service_provider_type IS 'Who provides the service: manufacturer, client, dealer, distributor';

-- =====================================================================
-- 3. ENGINEER CERTIFICATIONS (Formal tracking)
-- =====================================================================

CREATE TABLE IF NOT EXISTS engineer_certifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    engineer_id UUID NOT NULL REFERENCES engineers(id),
    manufacturer_id UUID REFERENCES organizations(id),
    equipment_catalog_id UUID REFERENCES equipment_catalog(id),
    
    -- Certification Details
    certification_name TEXT NOT NULL,
    certification_number TEXT NOT NULL,
    certification_level TEXT,                   -- 'Basic', 'Advanced', 'Expert', 'Trainer'
    
    -- Issuer
    issued_by TEXT NOT NULL,                    -- Organization/authority that issued
    issued_by_org_id UUID REFERENCES organizations(id),
    
    -- Dates
    issue_date DATE NOT NULL,
    expiry_date DATE,                           -- NULL = does not expire
    
    -- Status
    status TEXT DEFAULT 'active',               -- 'active', 'expired', 'suspended', 'revoked'
    
    -- Verification
    verification_url TEXT,
    verification_code TEXT,
    verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMPTZ,
    verified_by TEXT,
    
    -- Scope
    scope_description TEXT,                     -- What this certification covers
    skills_covered TEXT[],                      -- Specific skills/areas
    
    -- Documents
    certificate_url TEXT,
    certificate_document_id UUID,
    
    -- Renewal
    renewable BOOLEAN DEFAULT true,
    renewal_required_months_before INT DEFAULT 1,
    last_renewed_date DATE,
    
    -- Notes
    notes TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_certification_status CHECK (status IN (
        'active', 'expired', 'suspended', 'revoked', 'pending'
    )),
    CONSTRAINT chk_expiry_after_issue CHECK (expiry_date IS NULL OR expiry_date > issue_date)
);

-- Indexes
CREATE INDEX idx_engineer_certifications_engineer ON engineer_certifications(engineer_id);
CREATE INDEX idx_engineer_certifications_manufacturer ON engineer_certifications(manufacturer_id);
CREATE INDEX idx_engineer_certifications_equipment ON engineer_certifications(equipment_catalog_id);
CREATE INDEX idx_engineer_certifications_status ON engineer_certifications(status);
CREATE INDEX idx_engineer_certifications_expiry ON engineer_certifications(expiry_date) WHERE expiry_date IS NOT NULL AND status = 'active';

-- Composite index for valid certifications
CREATE INDEX idx_engineer_certifications_valid ON engineer_certifications(
    engineer_id, manufacturer_id, status
) WHERE status = 'active' AND (expiry_date IS NULL OR expiry_date > CURRENT_DATE);

-- GIN index for skills
CREATE INDEX idx_engineer_certifications_skills ON engineer_certifications USING GIN (skills_covered);

COMMENT ON TABLE engineer_certifications IS 'Formal certifications from manufacturers and training organizations';
COMMENT ON COLUMN engineer_certifications.certification_level IS 'Basic, Advanced, Expert, or Trainer level';
COMMENT ON COLUMN engineer_certifications.renewable IS 'Whether certification needs periodic renewal';

-- =====================================================================
-- 4. HELPER FUNCTIONS
-- =====================================================================

-- Function: Find eligible engineers for equipment and support level
CREATE OR REPLACE FUNCTION find_eligible_engineers(
    p_equipment_catalog_id UUID,
    p_support_level TEXT DEFAULT 'L1',
    p_must_be_certified BOOLEAN DEFAULT false
) RETURNS TABLE (
    engineer_id UUID,
    engineer_name TEXT,
    engineer_email TEXT,
    support_level TEXT,
    certified BOOLEAN,
    years_experience INT,
    total_repairs INT,
    first_time_fix_rate NUMERIC,
    customer_satisfaction NUMERIC,
    can_do_remote BOOLEAN,
    can_do_onsite BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.name,
        e.email,
        eee.support_level,
        eee.certified,
        eee.years_experience,
        eee.total_repairs_completed,
        eee.first_time_fix_rate,
        eee.customer_satisfaction_avg,
        eee.can_do_remote,
        eee.can_do_onsite
    FROM engineers e
    JOIN engineer_equipment_expertise eee ON e.id = eee.engineer_id
    WHERE eee.equipment_catalog_id = p_equipment_catalog_id
      AND eee.support_level = p_support_level
      AND eee.is_active = true
      AND eee.available_for_assignment = true
      AND e.status = 'available'
      AND (p_must_be_certified = false OR eee.certified = true)
      AND (eee.effective_to IS NULL OR eee.effective_to >= CURRENT_DATE)
    ORDER BY 
        eee.certified DESC,
        eee.first_time_fix_rate DESC NULLS LAST,
        eee.customer_satisfaction_avg DESC NULLS LAST,
        eee.years_experience DESC;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION find_eligible_engineers IS 'Find engineers qualified for equipment at specified support level';

-- Function: Get service configuration for equipment
CREATE OR REPLACE FUNCTION get_service_configuration(
    p_manufacturer_id UUID,
    p_equipment_catalog_id UUID DEFAULT NULL
) RETURNS TABLE (
    config_id UUID,
    service_provider_type TEXT,
    service_provider_org_id UUID,
    service_scope TEXT[],
    requires_oem_parts BOOLEAN,
    requires_certified_engineer BOOLEAN,
    sla_response_hours INT,
    sla_resolution_hours INT,
    priority INT,
    config_level TEXT
) AS $$
BEGIN
    RETURN QUERY
    -- First try equipment-specific config
    SELECT 
        msc.id,
        msc.service_provider_type,
        msc.service_provider_org_id,
        msc.service_scope,
        msc.requires_oem_parts,
        msc.requires_certified_engineer,
        msc.sla_response_hours,
        msc.sla_resolution_hours,
        msc.priority,
        'equipment-specific'::TEXT as config_level
    FROM manufacturer_service_config msc
    WHERE msc.manufacturer_id = p_manufacturer_id
      AND msc.equipment_catalog_id = p_equipment_catalog_id
      AND msc.effective_to IS NULL
    ORDER BY msc.priority DESC
    LIMIT 1;
    
    -- If no equipment-specific config, try manufacturer-level
    IF NOT FOUND THEN
        RETURN QUERY
        SELECT 
            msc.id,
            msc.service_provider_type,
            msc.service_provider_org_id,
            msc.service_scope,
            msc.requires_oem_parts,
            msc.requires_certified_engineer,
            msc.sla_response_hours,
            msc.sla_resolution_hours,
            msc.priority,
            'manufacturer-level'::TEXT as config_level
        FROM manufacturer_service_config msc
        WHERE msc.manufacturer_id = p_manufacturer_id
          AND msc.equipment_catalog_id IS NULL
          AND msc.effective_to IS NULL
        ORDER BY msc.priority DESC
        LIMIT 1;
    END IF;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_service_configuration IS 'Get service configuration with hierarchy: Equipment-specific > Manufacturer-level';

-- Function: Check if engineer is qualified for equipment
CREATE OR REPLACE FUNCTION is_engineer_qualified(
    p_engineer_id UUID,
    p_equipment_catalog_id UUID,
    p_required_support_level TEXT
) RETURNS BOOLEAN AS $$
DECLARE
    v_qualified BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1
        FROM engineer_equipment_expertise eee
        WHERE eee.engineer_id = p_engineer_id
          AND eee.equipment_catalog_id = p_equipment_catalog_id
          AND eee.support_level >= p_required_support_level  -- L3 can do L1, L2 can do L1, etc.
          AND eee.is_active = true
          AND eee.available_for_assignment = true
          AND (eee.effective_to IS NULL OR eee.effective_to >= CURRENT_DATE)
    ) INTO v_qualified;
    
    RETURN v_qualified;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION is_engineer_qualified IS 'Check if engineer is qualified for equipment at required support level';

-- Function: Get engineer expertise summary
CREATE OR REPLACE FUNCTION get_engineer_expertise_summary(
    p_engineer_id UUID
) RETURNS TABLE (
    equipment_type TEXT,
    equipment_model TEXT,
    manufacturer_name TEXT,
    support_level TEXT,
    certified BOOLEAN,
    years_experience INT,
    total_repairs INT,
    success_rate NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ec.equipment_type,
        ec.model_name,
        o.name,
        eee.support_level,
        eee.certified,
        eee.years_experience,
        eee.total_repairs_completed,
        CASE 
            WHEN eee.total_repairs_completed > 0 
            THEN ROUND((eee.successful_repairs::NUMERIC / eee.total_repairs_completed::NUMERIC) * 100, 2)
            ELSE NULL
        END as success_rate
    FROM engineer_equipment_expertise eee
    JOIN equipment_catalog ec ON eee.equipment_catalog_id = ec.id
    JOIN organizations o ON ec.manufacturer_id = o.id
    WHERE eee.engineer_id = p_engineer_id
      AND eee.is_active = true
    ORDER BY eee.support_level DESC, ec.equipment_type;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_engineer_expertise_summary IS 'Get summary of engineer expertise across all equipment';

-- Function: Get expiring certifications
CREATE OR REPLACE FUNCTION get_expiring_certifications(
    p_days_ahead INT DEFAULT 30
) RETURNS TABLE (
    engineer_id UUID,
    engineer_name TEXT,
    certification_name TEXT,
    expiry_date DATE,
    days_until_expiry INT,
    equipment_type TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.name,
        ec_cert.certification_name,
        ec_cert.expiry_date,
        (ec_cert.expiry_date - CURRENT_DATE)::INT as days_until_expiry,
        COALESCE(ec.equipment_type, 'General') as equipment_type
    FROM engineer_certifications ec_cert
    JOIN engineers e ON ec_cert.engineer_id = e.id
    LEFT JOIN equipment_catalog ec ON ec_cert.equipment_catalog_id = ec.id
    WHERE ec_cert.status = 'active'
      AND ec_cert.expiry_date IS NOT NULL
      AND ec_cert.expiry_date <= CURRENT_DATE + p_days_ahead
      AND ec_cert.expiry_date >= CURRENT_DATE
    ORDER BY ec_cert.expiry_date ASC, e.name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_expiring_certifications IS 'Get certifications expiring within specified days';

-- =====================================================================
-- 5. VIEWS
-- =====================================================================

-- View: Engineer expertise with equipment details
CREATE OR REPLACE VIEW engineer_expertise_with_details AS
SELECT 
    eee.*,
    e.name as engineer_name,
    e.email as engineer_email,
    e.phone as engineer_phone,
    e.status as engineer_status,
    ec.equipment_type,
    ec.model_name as equipment_model,
    o.name as manufacturer_name,
    CASE 
        WHEN eee.total_repairs_completed > 0 
        THEN ROUND((eee.successful_repairs::NUMERIC / eee.total_repairs_completed::NUMERIC) * 100, 2)
        ELSE NULL
    END as success_rate_percentage
FROM engineer_equipment_expertise eee
JOIN engineers e ON eee.engineer_id = e.id
JOIN equipment_catalog ec ON eee.equipment_catalog_id = ec.id
JOIN organizations o ON ec.manufacturer_id = o.id
WHERE eee.is_active = true;

COMMENT ON VIEW engineer_expertise_with_details IS 'Engineer expertise with engineer, equipment, and manufacturer details';

-- View: Active certifications
CREATE OR REPLACE VIEW active_certifications AS
SELECT 
    ec.*,
    e.name as engineer_name,
    e.email as engineer_email,
    o_mfr.name as manufacturer_name,
    eq.equipment_type,
    eq.model_name as equipment_model,
    CASE 
        WHEN ec.expiry_date IS NULL THEN 'Never expires'
        WHEN ec.expiry_date > CURRENT_DATE + INTERVAL '90 days' THEN 'Valid'
        WHEN ec.expiry_date > CURRENT_DATE + INTERVAL '30 days' THEN 'Expiring Soon'
        ELSE 'Expiring Critical'
    END as expiry_status
FROM engineer_certifications ec
JOIN engineers e ON ec.engineer_id = e.id
LEFT JOIN organizations o_mfr ON ec.manufacturer_id = o_mfr.id
LEFT JOIN equipment_catalog eq ON ec.equipment_catalog_id = eq.id
WHERE ec.status = 'active'
  AND (ec.expiry_date IS NULL OR ec.expiry_date >= CURRENT_DATE);

COMMENT ON VIEW active_certifications IS 'All active certifications with expiry status';

-- View: Service configuration summary
CREATE OR REPLACE VIEW service_configuration_summary AS
SELECT 
    msc.*,
    o_mfr.name as manufacturer_name,
    COALESCE(ec.model_name, 'All Equipment') as equipment_scope,
    o_provider.name as service_provider_name,
    CASE 
        WHEN msc.equipment_catalog_id IS NOT NULL THEN 'Equipment-Specific'
        ELSE 'Manufacturer-Level'
    END as config_scope
FROM manufacturer_service_config msc
JOIN organizations o_mfr ON msc.manufacturer_id = o_mfr.id
LEFT JOIN equipment_catalog ec ON msc.equipment_catalog_id = ec.id
LEFT JOIN organizations o_provider ON msc.service_provider_org_id = o_provider.id
WHERE msc.effective_to IS NULL;

COMMENT ON VIEW service_configuration_summary IS 'Active service configurations with details';

-- =====================================================================
-- 6. TRIGGERS FOR UPDATED_AT
-- =====================================================================

CREATE OR REPLACE FUNCTION update_engineer_expertise_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_engineer_expertise_updated
    BEFORE UPDATE ON engineer_equipment_expertise
    FOR EACH ROW
    EXECUTE FUNCTION update_engineer_expertise_timestamp();

CREATE TRIGGER trigger_manufacturer_service_updated
    BEFORE UPDATE ON manufacturer_service_config
    FOR EACH ROW
    EXECUTE FUNCTION update_engineer_expertise_timestamp();

CREATE TRIGGER trigger_engineer_certifications_updated
    BEFORE UPDATE ON engineer_certifications
    FOR EACH ROW
    EXECUTE FUNCTION update_engineer_expertise_timestamp();

-- =====================================================================
-- 7. SAMPLE DATA NOTES
-- =====================================================================

-- Note: Sample data will be added in migration 022-data-migration-seeding.sql

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

-- Verification queries
DO $$
BEGIN
    RAISE NOTICE 'Migration 017 complete!';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - engineer_equipment_expertise';
    RAISE NOTICE '  - manufacturer_service_config';
    RAISE NOTICE '  - engineer_certifications';
    RAISE NOTICE 'Created 5 helper functions';
    RAISE NOTICE 'Created 3 views';
    RAISE NOTICE 'Ready for AI Assignment Optimizer (T2C.3)!';
END $$;
