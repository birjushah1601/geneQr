-- Migration: Create equipment relationships history table
-- Ticket: T1.4
-- Purpose: Track equipment ownership, location, and relationship changes over time

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS btree_gist;  -- For EXCLUDE constraint with date ranges

-- =============================================================================
-- STEP 1: Create equipment_relationships table
-- =============================================================================

CREATE TABLE IF NOT EXISTS equipment_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Equipment reference
    equipment_id VARCHAR(32) NOT NULL REFERENCES equipment_registry(id) ON DELETE CASCADE,
    
    -- Relationship type
    relationship_type VARCHAR(50) NOT NULL,
    
    -- Related entity (organization or facility)
    related_org_id UUID,
    related_facility_id UUID,
    
    -- Entity details (denormalized for historical accuracy)
    entity_name VARCHAR(500) NOT NULL,
    entity_type VARCHAR(50),
    
    -- Location details (for facility relationships)
    location_details JSONB,
    address JSONB,
    
    -- Temporal validity
    effective_from DATE NOT NULL,
    effective_to DATE,
    
    -- Relationship details
    relationship_status VARCHAR(50) NOT NULL DEFAULT 'active',
    
    -- Transfer/change information
    reason VARCHAR(255),
    transfer_type VARCHAR(50),
    transfer_documents JSONB,
    
    -- Financial details (if applicable)
    transaction_amount DECIMAL(15,2),
    transaction_currency VARCHAR(3) DEFAULT 'INR',
    contract_reference VARCHAR(100),
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    terminated_at TIMESTAMPTZ,
    terminated_by VARCHAR(255),
    termination_reason TEXT,
    
    -- Constraints
    CONSTRAINT relationship_type_check CHECK (
        relationship_type IN (
            'manufacturer', 'distributor', 'dealer', 'owner', 
            'facility', 'service_provider', 'leasing_company', 'other'
        )
    ),
    CONSTRAINT relationship_status_check CHECK (
        relationship_status IN ('active', 'terminated', 'transferred', 'expired')
    ),
    CONSTRAINT effective_dates_valid CHECK (
        effective_to IS NULL OR effective_to >= effective_from
    )
);

-- Add EXCLUDE constraint to prevent overlapping active relationships of same type
-- Note: This ensures only one active owner/facility relationship exists at any point in time
ALTER TABLE equipment_relationships
    ADD CONSTRAINT equipment_rel_no_overlap 
    EXCLUDE USING gist (
        equipment_id WITH =,
        relationship_type WITH =,
        daterange(effective_from, COALESCE(effective_to, '9999-12-31'::date), '[]') WITH &&
    );

-- =============================================================================
-- STEP 2: Create indexes for performance
-- =============================================================================

-- Primary lookups
CREATE INDEX idx_equipment_rel_equipment ON equipment_relationships(equipment_id);
CREATE INDEX idx_equipment_rel_type ON equipment_relationships(relationship_type);
CREATE INDEX idx_equipment_rel_org ON equipment_relationships(related_org_id) WHERE related_org_id IS NOT NULL;
CREATE INDEX idx_equipment_rel_facility ON equipment_relationships(related_facility_id) WHERE related_facility_id IS NOT NULL;
CREATE INDEX idx_equipment_rel_status ON equipment_relationships(relationship_status);

-- Temporal queries
CREATE INDEX idx_equipment_rel_dates ON equipment_relationships(effective_from, effective_to);
CREATE INDEX idx_equipment_rel_effective_from ON equipment_relationships(effective_from);

-- Current relationships (most common query)
CREATE INDEX idx_equipment_rel_current ON equipment_relationships(equipment_id, relationship_type)
    WHERE relationship_status = 'active' AND effective_to IS NULL;

-- GIN indexes for JSONB searches
CREATE INDEX idx_equipment_rel_location ON equipment_relationships USING gin(location_details);
CREATE INDEX idx_equipment_rel_address ON equipment_relationships USING gin(address);
CREATE INDEX idx_equipment_rel_documents ON equipment_relationships USING gin(transfer_documents);

-- Composite index for common query pattern
CREATE INDEX idx_equipment_rel_lookup ON equipment_relationships(
    equipment_id, relationship_type, relationship_status, effective_from DESC
);

-- =============================================================================
-- STEP 3: Migrate existing data from equipment_registry
-- =============================================================================

DO $$
DECLARE
    v_owner_count INTEGER := 0;
    v_facility_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Starting equipment relationships migration...';
    
    -- 3a. Migrate ownership relationships
    INSERT INTO equipment_relationships (
        equipment_id,
        relationship_type,
        related_org_id,
        entity_name,
        entity_type,
        location_details,
        address,
        effective_from,
        relationship_status,
        reason,
        created_by
    )
    SELECT
        id as equipment_id,
        'owner' as relationship_type,
        CASE 
            WHEN customer_id ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
            THEN customer_id::uuid
            ELSE NULL
        END as related_org_id,
        customer_name as entity_name,
        'hospital' as entity_type,  -- Default type
        jsonb_build_object('location', installation_location) as location_details,
        installation_address as address,
        COALESCE(installation_date, created_at::date) as effective_from,
        'active' as relationship_status,
        'Migrated from equipment_registry' as reason,
        'migration_T1.4' as created_by
    FROM equipment_registry
    WHERE customer_id IS NOT NULL
        AND customer_name IS NOT NULL
        AND customer_name != ''
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_owner_count = ROW_COUNT;
    RAISE NOTICE 'Migrated % owner relationships', v_owner_count;
    
    -- 3b. Migrate facility/location relationships
    INSERT INTO equipment_relationships (
        equipment_id,
        relationship_type,
        entity_name,
        entity_type,
        location_details,
        address,
        effective_from,
        relationship_status,
        reason,
        created_by
    )
    SELECT
        id as equipment_id,
        'facility' as relationship_type,
        COALESCE(customer_name || ' - ' || installation_location, customer_name) as entity_name,
        'hospital' as entity_type,
        jsonb_build_object(
            'location', installation_location,
            'facility_name', customer_name,
            'department', COALESCE(category, 'General')
        ) as location_details,
        installation_address as address,
        COALESCE(installation_date, created_at::date) as effective_from,
        'active' as relationship_status,
        'Migrated from equipment_registry - facility location' as reason,
        'migration_T1.4' as created_by
    FROM equipment_registry
    WHERE installation_location IS NOT NULL
        AND installation_location != ''
        AND customer_name IS NOT NULL
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_facility_count = ROW_COUNT;
    RAISE NOTICE 'Migrated % facility relationships', v_facility_count;
    
    -- Summary
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Migration Summary:';
    RAISE NOTICE '  - Owner relationships: %', v_owner_count;
    RAISE NOTICE '  - Facility relationships: %', v_facility_count;
    RAISE NOTICE '  - Total relationships: %', v_owner_count + v_facility_count;
    RAISE NOTICE '==============================================';
END $$;

-- =============================================================================
-- STEP 4: Create helper functions
-- =============================================================================

-- Function to get current owner of equipment
CREATE OR REPLACE FUNCTION get_current_owner(p_equipment_id VARCHAR)
RETURNS TABLE (
    owner_id UUID,
    owner_name VARCHAR,
    since_date DATE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        er.related_org_id as owner_id,
        er.entity_name as owner_name,
        er.effective_from as since_date
    FROM equipment_relationships er
    WHERE er.equipment_id = p_equipment_id
        AND er.relationship_type = 'owner'
        AND er.relationship_status = 'active'
        AND er.effective_to IS NULL
    ORDER BY er.effective_from DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function to get owner at specific date (temporal query)
CREATE OR REPLACE FUNCTION get_owner_at_date(
    p_equipment_id VARCHAR,
    p_date DATE
)
RETURNS TABLE (
    owner_id UUID,
    owner_name VARCHAR,
    from_date DATE,
    to_date DATE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        er.related_org_id as owner_id,
        er.entity_name as owner_name,
        er.effective_from as from_date,
        er.effective_to as to_date
    FROM equipment_relationships er
    WHERE er.equipment_id = p_equipment_id
        AND er.relationship_type = 'owner'
        AND p_date >= er.effective_from
        AND (er.effective_to IS NULL OR p_date <= er.effective_to)
    ORDER BY er.effective_from DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function to transfer equipment ownership
CREATE OR REPLACE FUNCTION transfer_equipment(
    p_equipment_id VARCHAR,
    p_new_owner_id UUID,
    p_new_owner_name VARCHAR,
    p_transfer_date DATE,
    p_reason VARCHAR,
    p_transfer_type VARCHAR DEFAULT 'transfer',
    p_amount DECIMAL DEFAULT NULL,
    p_created_by VARCHAR DEFAULT 'system'
)
RETURNS UUID AS $$
DECLARE
    v_new_rel_id UUID;
    v_old_owner_name VARCHAR;
BEGIN
    -- Get current owner name for logging
    SELECT entity_name INTO v_old_owner_name
    FROM equipment_relationships
    WHERE equipment_id = p_equipment_id
        AND relationship_type = 'owner'
        AND relationship_status = 'active'
        AND effective_to IS NULL
    LIMIT 1;
    
    -- Close current ownership
    UPDATE equipment_relationships
    SET effective_to = p_transfer_date,
        relationship_status = 'transferred',
        terminated_at = NOW(),
        terminated_by = p_created_by,
        termination_reason = p_reason
    WHERE equipment_id = p_equipment_id
        AND relationship_type = 'owner'
        AND relationship_status = 'active'
        AND effective_to IS NULL;
    
    -- Create new ownership
    INSERT INTO equipment_relationships (
        equipment_id,
        relationship_type,
        related_org_id,
        entity_name,
        effective_from,
        relationship_status,
        reason,
        transfer_type,
        transaction_amount,
        created_by,
        notes
    ) VALUES (
        p_equipment_id,
        'owner',
        p_new_owner_id,
        p_new_owner_name,
        p_transfer_date,
        'active',
        p_reason,
        p_transfer_type,
        p_amount,
        p_created_by,
        'Transferred from: ' || COALESCE(v_old_owner_name, 'Unknown')
    )
    RETURNING id INTO v_new_rel_id;
    
    RAISE NOTICE 'Equipment % transferred from % to % on %', 
        p_equipment_id, v_old_owner_name, p_new_owner_name, p_transfer_date;
    
    RETURN v_new_rel_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get relationship history
CREATE OR REPLACE FUNCTION get_relationship_history(
    p_equipment_id VARCHAR,
    p_relationship_type VARCHAR DEFAULT NULL
)
RETURNS TABLE (
    relationship_id UUID,
    relationship_type VARCHAR,
    entity_name VARCHAR,
    effective_from DATE,
    effective_to DATE,
    status VARCHAR,
    reason VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        er.id as relationship_id,
        er.relationship_type,
        er.entity_name,
        er.effective_from,
        er.effective_to,
        er.relationship_status as status,
        er.reason
    FROM equipment_relationships er
    WHERE er.equipment_id = p_equipment_id
        AND (p_relationship_type IS NULL OR er.relationship_type = p_relationship_type)
    ORDER BY er.effective_from DESC, er.created_at DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================================================
-- STEP 5: Create views for easy access
-- =============================================================================

-- View: Current equipment relationships (active only)
CREATE OR REPLACE VIEW equipment_current_relationships AS
SELECT
    er.equipment_id,
    er.relationship_type,
    er.related_org_id,
    er.entity_name,
    er.entity_type,
    er.location_details,
    er.effective_from,
    er.created_at
FROM equipment_relationships er
WHERE er.relationship_status = 'active'
    AND er.effective_to IS NULL;

COMMENT ON VIEW equipment_current_relationships IS 
    'All currently active equipment relationships (owner, facility, etc.)';

-- View: Equipment with current owner and location
CREATE OR REPLACE VIEW equipment_with_current_info AS
SELECT
    e.*,
    owner.entity_name as current_owner_name,
    owner.related_org_id as current_owner_id,
    owner.effective_from as owner_since,
    facility.entity_name as current_facility_name,
    facility.location_details as current_location,
    facility.address as current_address
FROM equipment_registry e
LEFT JOIN equipment_relationships owner ON (
    e.id = owner.equipment_id
    AND owner.relationship_type = 'owner'
    AND owner.relationship_status = 'active'
    AND owner.effective_to IS NULL
)
LEFT JOIN equipment_relationships facility ON (
    e.id = facility.equipment_id
    AND facility.relationship_type = 'facility'
    AND facility.relationship_status = 'active'
    AND facility.effective_to IS NULL
);

COMMENT ON VIEW equipment_with_current_info IS 
    'Equipment registry joined with current owner and facility information';

-- =============================================================================
-- STEP 6: Create trigger for timestamp updates
-- =============================================================================

CREATE OR REPLACE FUNCTION update_equipment_rel_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Note: We don't have updated_at column, so this is for future if we add it
-- CREATE TRIGGER trigger_update_equipment_rel_timestamp
--     BEFORE UPDATE ON equipment_relationships
--     FOR EACH ROW
--     EXECUTE FUNCTION update_equipment_rel_timestamp();

-- =============================================================================
-- STEP 7: Add comments for documentation
-- =============================================================================

COMMENT ON TABLE equipment_relationships IS 
    'Tracks equipment ownership, location, and relationship history over time. Supports temporal queries and full audit trail.';

COMMENT ON COLUMN equipment_relationships.relationship_type IS 
    'Type of relationship: manufacturer, distributor, dealer, owner, facility, service_provider, leasing_company, other';

COMMENT ON COLUMN equipment_relationships.effective_from IS 
    'Start date of this relationship. Required.';

COMMENT ON COLUMN equipment_relationships.effective_to IS 
    'End date of this relationship. NULL means currently active.';

COMMENT ON COLUMN equipment_relationships.relationship_status IS 
    'Current status: active, terminated, transferred, expired';

COMMENT ON COLUMN equipment_relationships.location_details IS 
    'JSONB containing detailed location info: building, floor, room, ward, department';

COMMENT ON COLUMN equipment_relationships.transfer_documents IS 
    'JSONB array of document URLs: invoices, contracts, agreements';

-- =============================================================================
-- STEP 8: Validation queries
-- =============================================================================

DO $$
DECLARE
    v_equipment_count INTEGER;
    v_relationships_count INTEGER;
    v_orphaned_count INTEGER;
BEGIN
    -- Count equipment
    SELECT COUNT(*) INTO v_equipment_count FROM equipment_registry;
    
    -- Count relationships
    SELECT COUNT(*) INTO v_relationships_count FROM equipment_relationships;
    
    -- Check for equipment without relationships
    SELECT COUNT(*) INTO v_orphaned_count
    FROM equipment_registry e
    LEFT JOIN equipment_relationships er ON e.id = er.equipment_id
    WHERE er.id IS NULL;
    
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Validation Results:';
    RAISE NOTICE '  - Total equipment: %', v_equipment_count;
    RAISE NOTICE '  - Total relationships: %', v_relationships_count;
    RAISE NOTICE '  - Equipment without relationships: %', v_orphaned_count;
    
    IF v_orphaned_count > 0 THEN
        RAISE WARNING '% equipment items have no relationships', v_orphaned_count;
    ELSE
        RAISE NOTICE '✓ All equipment have at least one relationship';
    END IF;
    
    RAISE NOTICE '==============================================';
END $$;

-- =============================================================================
-- Migration Complete
-- =============================================================================

DO $$
BEGIN
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'T1.4 Migration Complete!';
    RAISE NOTICE 'Equipment relationships history tracking enabled';
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Features enabled:';
    RAISE NOTICE '  ✓ Ownership history tracking';
    RAISE NOTICE '  ✓ Location change history';
    RAISE NOTICE '  ✓ Temporal queries (owner at specific date)';
    RAISE NOTICE '  ✓ Transfer equipment function';
    RAISE NOTICE '  ✓ Relationship audit trail';
    RAISE NOTICE '  ✓ Views for easy access';
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Example queries:';
    RAISE NOTICE '  SELECT * FROM get_current_owner(''EQ-12345'');';
    RAISE NOTICE '  SELECT * FROM get_owner_at_date(''EQ-12345'', ''2024-01-15'');';
    RAISE NOTICE '  SELECT * FROM equipment_with_current_info WHERE id = ''EQ-12345'';';
    RAISE NOTICE '==============================================';
END $$;
