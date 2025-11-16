-- Migration: 016-create-equipment-catalog.sql
-- Description: Equipment Catalog & Parts Management
-- Ticket: T2B.1
-- Date: 2025-11-16
-- 
-- This migration creates:
-- 1. equipment_catalog - Master list of equipment types
-- 2. equipment_parts - Parts and accessories catalog
-- 3. equipment_parts_context - Context-specific parts (ICU, General Ward, etc.)
-- 4. equipment_compatibility - Parts compatibility matrix
--
-- Purpose:
-- - Separate equipment types (catalog) from installed instances (equipment_registry)
-- - Support context-aware parts recommendations
-- - Track parts specifications, pricing, and availability

-- =====================================================================
-- 1. EQUIPMENT CATALOG (Master List)
-- =====================================================================

CREATE TABLE IF NOT EXISTS equipment_catalog (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic Information
    manufacturer_id UUID NOT NULL REFERENCES organizations(id),
    equipment_type TEXT NOT NULL,               -- 'Ventilator', 'MRI', 'CT Scanner', etc.
    model_number TEXT NOT NULL,
    model_name TEXT NOT NULL,
    
    -- Classification
    category TEXT NOT NULL,                     -- 'Diagnostic', 'Life Support', 'Surgical', 'Laboratory', etc.
    sub_category TEXT,                          -- More specific classification
    
    -- Specifications
    specifications JSONB DEFAULT '{}',          -- Technical specs
    dimensions JSONB,                           -- {"length": "120cm", "width": "60cm", "height": "80cm"}
    weight_kg NUMERIC(10,2),
    power_requirements JSONB,                   -- {"voltage": "220V", "frequency": "50Hz", "power": "1500W"}
    
    -- Documentation
    description TEXT,
    features TEXT[],                            -- Array of key features
    service_manual_url TEXT,
    user_manual_url TEXT,
    brochure_url TEXT,
    image_urls TEXT[],                          -- Product images
    
    -- Service Information
    typical_lifespan_years INT,
    maintenance_interval_months INT,
    requires_certification BOOLEAN DEFAULT false,
    
    -- Regulatory
    regulatory_approvals JSONB,                 -- {"FDA": "K123456", "CE": "CE-2024-001"}
    compliance_standards TEXT[],                -- ['ISO 13485', 'ISO 80601', ...]
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    discontinued_date DATE,
    replacement_model_id UUID REFERENCES equipment_catalog(id),
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_equipment_type CHECK (equipment_type != ''),
    CONSTRAINT chk_model_number CHECK (model_number != ''),
    CONSTRAINT chk_category CHECK (category IN (
        'Diagnostic', 'Life Support', 'Surgical', 'Laboratory', 
        'Monitoring', 'Therapeutic', 'Imaging', 'Sterilization', 'Other'
    )),
    CONSTRAINT chk_weight CHECK (weight_kg IS NULL OR weight_kg > 0),
    CONSTRAINT chk_lifespan CHECK (typical_lifespan_years IS NULL OR typical_lifespan_years > 0),
    CONSTRAINT chk_maintenance CHECK (maintenance_interval_months IS NULL OR maintenance_interval_months > 0)
);

-- Unique constraint: One model per manufacturer
CREATE UNIQUE INDEX idx_equipment_catalog_unique_model 
ON equipment_catalog(manufacturer_id, model_number) 
WHERE is_active = true;

-- Indexes for common queries
CREATE INDEX idx_equipment_catalog_manufacturer ON equipment_catalog(manufacturer_id);
CREATE INDEX idx_equipment_catalog_type ON equipment_catalog(equipment_type);
CREATE INDEX idx_equipment_catalog_category ON equipment_catalog(category);
CREATE INDEX idx_equipment_catalog_active ON equipment_catalog(is_active) WHERE is_active = true;

-- GIN index for JSONB queries
CREATE INDEX idx_equipment_catalog_specs ON equipment_catalog USING GIN (specifications);

-- Full-text search
CREATE INDEX idx_equipment_catalog_search ON equipment_catalog USING GIN (
    to_tsvector('english', 
        COALESCE(model_name, '') || ' ' || 
        COALESCE(equipment_type, '') || ' ' || 
        COALESCE(description, '')
    )
);

COMMENT ON TABLE equipment_catalog IS 'Master catalog of equipment types and models (NOT installed instances)';
COMMENT ON COLUMN equipment_catalog.equipment_type IS 'Generic type: Ventilator, MRI, CT Scanner, etc.';
COMMENT ON COLUMN equipment_catalog.model_number IS 'Manufacturer model number';
COMMENT ON COLUMN equipment_catalog.specifications IS 'Technical specifications in JSONB format';

-- =====================================================================
-- 2. EQUIPMENT PARTS (Parts Catalog)
-- =====================================================================

CREATE TABLE IF NOT EXISTS equipment_parts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic Information
    equipment_catalog_id UUID NOT NULL REFERENCES equipment_catalog(id),
    part_number TEXT NOT NULL,
    part_name TEXT NOT NULL,
    
    -- Classification
    part_category TEXT NOT NULL,                -- 'consumable', 'replaceable', 'optional', 'tool'
    part_type TEXT NOT NULL,                    -- 'accessory', 'component', 'attachment', 'tool', 'supply'
    
    -- Specifications
    description TEXT,
    specifications JSONB DEFAULT '{}',
    dimensions JSONB,
    weight_kg NUMERIC(10,2),
    material TEXT,
    
    -- Compatibility
    compatible_models TEXT[],                   -- Model numbers this part works with
    replaces_part_number TEXT,                  -- Old part number being replaced
    
    -- Lifecycle
    is_oem BOOLEAN DEFAULT true,                -- Original Equipment Manufacturer part
    is_universal BOOLEAN DEFAULT false,         -- Works with multiple equipment types
    is_critical BOOLEAN DEFAULT false,          -- Critical for equipment operation
    lifespan_hours INT,                         -- Expected lifespan in hours
    replacement_frequency_months INT,           -- Recommended replacement interval
    
    -- Inventory & Pricing
    unit_of_measure TEXT DEFAULT 'piece',       -- 'piece', 'set', 'box', 'meter', etc.
    min_order_quantity INT DEFAULT 1,
    standard_price NUMERIC(10,2),
    currency TEXT DEFAULT 'INR',
    lead_time_days INT,                         -- Procurement lead time
    
    -- Storage
    storage_conditions TEXT,                    -- 'Room temperature', 'Refrigerated', etc.
    shelf_life_months INT,                      -- For consumables
    
    -- Documentation
    image_urls TEXT[],
    manual_url TEXT,
    datasheet_url TEXT,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    discontinued_date DATE,
    replacement_part_id UUID REFERENCES equipment_parts(id),
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_part_number CHECK (part_number != ''),
    CONSTRAINT chk_part_category CHECK (part_category IN (
        'consumable', 'replaceable', 'optional', 'tool', 'calibration', 'maintenance'
    )),
    CONSTRAINT chk_part_type CHECK (part_type IN (
        'accessory', 'component', 'attachment', 'tool', 'supply', 'sensor', 
        'filter', 'battery', 'cable', 'adapter', 'probe', 'other'
    )),
    CONSTRAINT chk_price CHECK (standard_price IS NULL OR standard_price >= 0),
    CONSTRAINT chk_min_order CHECK (min_order_quantity > 0),
    CONSTRAINT chk_lead_time CHECK (lead_time_days IS NULL OR lead_time_days >= 0),
    CONSTRAINT chk_lifespan_hours CHECK (lifespan_hours IS NULL OR lifespan_hours > 0)
);

-- Unique constraint
CREATE UNIQUE INDEX idx_equipment_parts_unique_number 
ON equipment_parts(equipment_catalog_id, part_number) 
WHERE is_active = true;

-- Indexes
CREATE INDEX idx_equipment_parts_catalog ON equipment_parts(equipment_catalog_id);
CREATE INDEX idx_equipment_parts_category ON equipment_parts(part_category);
CREATE INDEX idx_equipment_parts_type ON equipment_parts(part_type);
CREATE INDEX idx_equipment_parts_critical ON equipment_parts(is_critical) WHERE is_critical = true;
CREATE INDEX idx_equipment_parts_active ON equipment_parts(is_active) WHERE is_active = true;

-- GIN index for array searches
CREATE INDEX idx_equipment_parts_compatible ON equipment_parts USING GIN (compatible_models);

-- Full-text search
CREATE INDEX idx_equipment_parts_search ON equipment_parts USING GIN (
    to_tsvector('english', 
        COALESCE(part_name, '') || ' ' || 
        COALESCE(part_number, '') || ' ' || 
        COALESCE(description, '')
    )
);

COMMENT ON TABLE equipment_parts IS 'Catalog of parts, accessories, and consumables for equipment';
COMMENT ON COLUMN equipment_parts.part_category IS 'High-level category: consumable, replaceable, optional, etc.';
COMMENT ON COLUMN equipment_parts.is_critical IS 'Critical for equipment operation (cannot run without it)';

-- =====================================================================
-- 3. EQUIPMENT PARTS CONTEXT (Context-Specific Parts)
-- =====================================================================

CREATE TABLE IF NOT EXISTS equipment_parts_context (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    equipment_catalog_id UUID NOT NULL REFERENCES equipment_catalog(id),
    part_id UUID NOT NULL REFERENCES equipment_parts(id),
    
    -- Context
    installation_context TEXT NOT NULL,         -- 'ICU', 'General Ward', 'OT', 'ER', 'Laboratory', etc.
    use_case TEXT,                              -- Specific use case within context
    
    -- Recommendation
    is_required BOOLEAN DEFAULT false,          -- Must-have for this context
    is_recommended BOOLEAN DEFAULT true,        -- Recommended but optional
    recommended_quantity INT DEFAULT 1,
    priority INT DEFAULT 5,                     -- Display order (1=highest)
    
    -- Rationale
    reason TEXT,                                -- Why this part for this context
    benefits TEXT[],                            -- Benefits in this context
    alternatives JSONB,                         -- Alternative parts for this context
    
    -- Usage
    typical_usage_frequency TEXT,              -- 'Daily', 'Weekly', 'Per Procedure', etc.
    estimated_monthly_consumption NUMERIC(10,2),
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_installation_context CHECK (installation_context IN (
        'ICU', 'General Ward', 'OT', 'ER', 'Laboratory', 
        'Dialysis', 'Radiology', 'Cardiology', 'Neonatal', 
        'Outpatient', 'Home Care', 'Ambulance', 'Other'
    )),
    CONSTRAINT chk_recommended_quantity CHECK (recommended_quantity > 0),
    CONSTRAINT chk_priority CHECK (priority >= 1 AND priority <= 10),
    CONSTRAINT chk_consumption CHECK (estimated_monthly_consumption IS NULL OR estimated_monthly_consumption >= 0)
);

-- Unique constraint: One part-context combination per equipment
CREATE UNIQUE INDEX idx_equipment_parts_context_unique 
ON equipment_parts_context(equipment_catalog_id, part_id, installation_context);

-- Indexes
CREATE INDEX idx_equipment_parts_context_equipment ON equipment_parts_context(equipment_catalog_id);
CREATE INDEX idx_equipment_parts_context_part ON equipment_parts_context(part_id);
CREATE INDEX idx_equipment_parts_context_context ON equipment_parts_context(installation_context);
CREATE INDEX idx_equipment_parts_context_required ON equipment_parts_context(is_required) WHERE is_required = true;
CREATE INDEX idx_equipment_parts_context_priority ON equipment_parts_context(equipment_catalog_id, installation_context, priority);

COMMENT ON TABLE equipment_parts_context IS 'Context-specific parts recommendations (ICU vs General Ward accessories)';
COMMENT ON COLUMN equipment_parts_context.installation_context IS 'Where equipment is installed: ICU, General Ward, OT, etc.';
COMMENT ON COLUMN equipment_parts_context.priority IS 'Display priority (1=show first, 10=show last)';

-- =====================================================================
-- 4. EQUIPMENT COMPATIBILITY (Cross-Equipment Parts)
-- =====================================================================

CREATE TABLE IF NOT EXISTS equipment_compatibility (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    part_id UUID NOT NULL REFERENCES equipment_parts(id),
    compatible_equipment_id UUID NOT NULL REFERENCES equipment_catalog(id),
    
    -- Compatibility Details
    compatibility_type TEXT NOT NULL,           -- 'direct', 'with_adapter', 'replacement', 'upgrade'
    compatibility_notes TEXT,
    requires_adapter BOOLEAN DEFAULT false,
    adapter_part_id UUID REFERENCES equipment_parts(id),
    
    -- Validation
    tested BOOLEAN DEFAULT false,
    test_date DATE,
    test_notes TEXT,
    certification_ref TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_compatibility_type CHECK (compatibility_type IN (
        'direct', 'with_adapter', 'replacement', 'upgrade', 'alternative'
    ))
);

-- Unique constraint
CREATE UNIQUE INDEX idx_equipment_compatibility_unique 
ON equipment_compatibility(part_id, compatible_equipment_id);

-- Indexes
CREATE INDEX idx_equipment_compatibility_part ON equipment_compatibility(part_id);
CREATE INDEX idx_equipment_compatibility_equipment ON equipment_compatibility(compatible_equipment_id);
CREATE INDEX idx_equipment_compatibility_tested ON equipment_compatibility(tested) WHERE tested = true;

COMMENT ON TABLE equipment_compatibility IS 'Tracks which parts are compatible with which equipment models';
COMMENT ON COLUMN equipment_compatibility.compatibility_type IS 'Type of compatibility: direct, needs adapter, etc.';

-- =====================================================================
-- 5. HELPER FUNCTIONS
-- =====================================================================

-- Function: Get all parts for an equipment type
CREATE OR REPLACE FUNCTION get_equipment_parts(
    p_equipment_catalog_id UUID,
    p_include_optional BOOLEAN DEFAULT true
) RETURNS TABLE (
    part_id UUID,
    part_number TEXT,
    part_name TEXT,
    part_category TEXT,
    part_type TEXT,
    is_critical BOOLEAN,
    standard_price NUMERIC,
    lead_time_days INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ep.id,
        ep.part_number,
        ep.part_name,
        ep.part_category,
        ep.part_type,
        ep.is_critical,
        ep.standard_price,
        ep.lead_time_days
    FROM equipment_parts ep
    WHERE ep.equipment_catalog_id = p_equipment_catalog_id
      AND ep.is_active = true
      AND (p_include_optional = true OR ep.is_critical = true)
    ORDER BY ep.is_critical DESC, ep.part_category, ep.part_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_equipment_parts IS 'Get all parts for an equipment type, optionally filtering critical only';

-- Function: Get context-specific parts for equipment
CREATE OR REPLACE FUNCTION get_context_specific_parts(
    p_equipment_catalog_id UUID,
    p_installation_context TEXT
) RETURNS TABLE (
    part_id UUID,
    part_number TEXT,
    part_name TEXT,
    part_category TEXT,
    is_required BOOLEAN,
    is_recommended BOOLEAN,
    recommended_quantity INT,
    priority INT,
    reason TEXT,
    standard_price NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ep.id,
        ep.part_number,
        ep.part_name,
        ep.part_category,
        epc.is_required,
        epc.is_recommended,
        epc.recommended_quantity,
        epc.priority,
        epc.reason,
        ep.standard_price
    FROM equipment_parts_context epc
    JOIN equipment_parts ep ON epc.part_id = ep.id
    WHERE epc.equipment_catalog_id = p_equipment_catalog_id
      AND epc.installation_context = p_installation_context
      AND ep.is_active = true
    ORDER BY 
        epc.is_required DESC,
        epc.priority ASC,
        ep.part_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_context_specific_parts IS 'Get parts recommended for specific installation context (ICU, General Ward, etc.)';

-- Function: Find compatible parts across equipment
CREATE OR REPLACE FUNCTION find_compatible_parts(
    p_equipment_catalog_id UUID
) RETURNS TABLE (
    part_id UUID,
    part_number TEXT,
    part_name TEXT,
    original_equipment_id UUID,
    original_equipment_name TEXT,
    compatibility_type TEXT,
    requires_adapter BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ep.id,
        ep.part_number,
        ep.part_name,
        ec_original.id,
        ec_original.model_name,
        compat.compatibility_type,
        compat.requires_adapter
    FROM equipment_compatibility compat
    JOIN equipment_parts ep ON compat.part_id = ep.id
    JOIN equipment_catalog ec_original ON ep.equipment_catalog_id = ec_original.id
    WHERE compat.compatible_equipment_id = p_equipment_catalog_id
      AND ep.is_active = true
    ORDER BY compat.compatibility_type, ep.part_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION find_compatible_parts IS 'Find parts from other equipment that are compatible with this equipment';

-- Function: Search equipment catalog
CREATE OR REPLACE FUNCTION search_equipment_catalog(
    p_search_term TEXT
) RETURNS TABLE (
    equipment_id UUID,
    equipment_type TEXT,
    model_number TEXT,
    model_name TEXT,
    manufacturer_name TEXT,
    category TEXT,
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ec.id,
        ec.equipment_type,
        ec.model_number,
        ec.model_name,
        o.name,
        ec.category,
        ts_rank(
            to_tsvector('english', 
                COALESCE(ec.model_name, '') || ' ' || 
                COALESCE(ec.equipment_type, '') || ' ' || 
                COALESCE(ec.description, '')
            ),
            plainto_tsquery('english', p_search_term)
        ) as rank
    FROM equipment_catalog ec
    JOIN organizations o ON ec.manufacturer_id = o.id
    WHERE ec.is_active = true
      AND (
          to_tsvector('english', 
              COALESCE(ec.model_name, '') || ' ' || 
              COALESCE(ec.equipment_type, '') || ' ' || 
              COALESCE(ec.description, '')
          ) @@ plainto_tsquery('english', p_search_term)
      )
    ORDER BY rank DESC, ec.model_name
    LIMIT 50;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION search_equipment_catalog IS 'Full-text search across equipment catalog';

-- =====================================================================
-- 6. VIEWS
-- =====================================================================

-- View: Equipment catalog with manufacturer details
CREATE OR REPLACE VIEW equipment_catalog_with_manufacturer AS
SELECT 
    ec.*,
    o.name as manufacturer_name,
    o.org_type as manufacturer_type,
    o.country as manufacturer_country
FROM equipment_catalog ec
JOIN organizations o ON ec.manufacturer_id = o.id
WHERE ec.is_active = true;

COMMENT ON VIEW equipment_catalog_with_manufacturer IS 'Equipment catalog with manufacturer information';

-- View: Parts with equipment details
CREATE OR REPLACE VIEW parts_with_equipment AS
SELECT 
    ep.*,
    ec.equipment_type,
    ec.model_number as equipment_model,
    ec.model_name as equipment_name,
    o.name as manufacturer_name
FROM equipment_parts ep
JOIN equipment_catalog ec ON ep.equipment_catalog_id = ec.id
JOIN organizations o ON ec.manufacturer_id = o.id
WHERE ep.is_active = true;

COMMENT ON VIEW parts_with_equipment IS 'Parts catalog with equipment and manufacturer details';

-- View: Context-specific parts summary
CREATE OR REPLACE VIEW context_parts_summary AS
SELECT 
    ec.id as equipment_id,
    ec.model_name as equipment_name,
    epc.installation_context,
    COUNT(*) as total_parts,
    COUNT(*) FILTER (WHERE epc.is_required) as required_parts,
    COUNT(*) FILTER (WHERE epc.is_recommended AND NOT epc.is_required) as recommended_parts,
    SUM(ep.standard_price * epc.recommended_quantity) as estimated_total_cost
FROM equipment_catalog ec
JOIN equipment_parts_context epc ON ec.id = epc.equipment_catalog_id
JOIN equipment_parts ep ON epc.part_id = ep.id
WHERE ec.is_active = true AND ep.is_active = true
GROUP BY ec.id, ec.model_name, epc.installation_context;

COMMENT ON VIEW context_parts_summary IS 'Summary of parts by equipment and installation context';

-- =====================================================================
-- 7. TRIGGERS FOR UPDATED_AT
-- =====================================================================

CREATE OR REPLACE FUNCTION update_equipment_catalog_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_equipment_catalog_updated
    BEFORE UPDATE ON equipment_catalog
    FOR EACH ROW
    EXECUTE FUNCTION update_equipment_catalog_timestamp();

CREATE TRIGGER trigger_equipment_parts_updated
    BEFORE UPDATE ON equipment_parts
    FOR EACH ROW
    EXECUTE FUNCTION update_equipment_catalog_timestamp();

CREATE TRIGGER trigger_equipment_parts_context_updated
    BEFORE UPDATE ON equipment_parts_context
    FOR EACH ROW
    EXECUTE FUNCTION update_equipment_catalog_timestamp();

-- =====================================================================
-- 8. SAMPLE DATA (for testing)
-- =====================================================================

-- Note: Sample data will be added in migration 022-data-migration-seeding.sql

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

-- Verification queries
DO $$
BEGIN
    RAISE NOTICE 'Migration 016 complete!';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - equipment_catalog';
    RAISE NOTICE '  - equipment_parts';
    RAISE NOTICE '  - equipment_parts_context';
    RAISE NOTICE '  - equipment_compatibility';
    RAISE NOTICE 'Created 4 helper functions';
    RAISE NOTICE '  Created 3 views';
    RAISE NOTICE 'Ready for use!';
END $$;
