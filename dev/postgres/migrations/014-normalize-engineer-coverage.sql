-- Migration: Normalize Engineer Coverage
-- Ticket: T2.3
-- Purpose: Move engineer coverage from arrays to proper normalized table

-- =============================================================================
-- STEP 1: Create engineer_coverage_areas table
-- =============================================================================

CREATE TABLE IF NOT EXISTS engineer_coverage_areas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Link to engineer
    engineer_id UUID NOT NULL REFERENCES engineers(id) ON DELETE CASCADE,
    
    -- Coverage definition
    coverage_type TEXT NOT NULL,
    coverage_value TEXT NOT NULL,
    
    -- Territory link (optional)
    territory_id UUID REFERENCES territories(id),
    
    -- Priority & assignment
    priority INT DEFAULT 1,
    coverage_radius_km INT,
    
    -- Temporal validity
    effective_from DATE NOT NULL DEFAULT CURRENT_DATE,
    effective_to DATE,
    
    -- Assignment rules
    can_emergency BOOLEAN DEFAULT true,
    can_scheduled BOOLEAN DEFAULT true,
    max_concurrent_tickets INT DEFAULT 3,
    
    -- Metadata
    assigned_by TEXT,
    assignment_reason TEXT,
    notes TEXT,
    metadata JSONB,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_coverage_type CHECK (coverage_type IN (
        'pincode', 'city', 'district', 'state', 'region', 'country'
    )),
    CONSTRAINT chk_priority CHECK (priority >= 1 AND priority <= 10),
    CONSTRAINT chk_effective_dates CHECK (
        effective_to IS NULL OR effective_to >= effective_from
    ),
    CONSTRAINT unique_coverage_per_engineer UNIQUE (
        engineer_id, coverage_type, coverage_value, effective_from
    )
);

COMMENT ON TABLE engineer_coverage_areas IS 
    'Normalized engineer coverage areas with temporal tracking and priority-based assignment';

COMMENT ON COLUMN engineer_coverage_areas.coverage_type IS 
    'Type of coverage: pincode, city, district, state, region, country';

COMMENT ON COLUMN engineer_coverage_areas.coverage_value IS 
    'Actual value (e.g., ''400001'' for pincode, ''Mumbai'' for city)';

COMMENT ON COLUMN engineer_coverage_areas.priority IS 
    '1=primary (first choice), 2=secondary (backup), 3+=tertiary';

COMMENT ON COLUMN engineer_coverage_areas.effective_from IS 
    'Start date when this coverage begins. Use for temporal queries.';

COMMENT ON COLUMN engineer_coverage_areas.effective_to IS 
    'End date when this coverage ends. NULL means currently active.';

-- =============================================================================
-- STEP 2: Create indexes for performance
-- =============================================================================

-- Primary lookups
CREATE INDEX idx_eng_coverage_engineer ON engineer_coverage_areas(engineer_id);
CREATE INDEX idx_eng_coverage_territory ON engineer_coverage_areas(territory_id);

-- "Who covers X?" queries
CREATE INDEX idx_eng_coverage_type_value ON engineer_coverage_areas(coverage_type, coverage_value);
CREATE INDEX idx_eng_coverage_value ON engineer_coverage_areas(coverage_value);

-- Current coverage (most common query)
CREATE INDEX idx_eng_coverage_current ON engineer_coverage_areas(coverage_type, coverage_value)
    WHERE effective_to IS NULL;

-- Priority-based lookups
CREATE INDEX idx_eng_coverage_priority ON engineer_coverage_areas(coverage_type, coverage_value, priority)
    WHERE effective_to IS NULL;

-- Temporal queries
CREATE INDEX idx_eng_coverage_dates ON engineer_coverage_areas(effective_from, effective_to);

-- Assignment rules
CREATE INDEX idx_eng_coverage_emergency ON engineer_coverage_areas(coverage_type, coverage_value)
    WHERE can_emergency = true AND effective_to IS NULL;

-- =============================================================================
-- STEP 3: Migrate existing data from arrays
-- =============================================================================

DO $$
DECLARE
    v_pincodes_migrated INTEGER := 0;
    v_cities_migrated INTEGER := 0;
    v_states_migrated INTEGER := 0;
BEGIN
    RAISE NOTICE 'Starting engineer coverage migration...';
    
    -- Migrate pincodes
    INSERT INTO engineer_coverage_areas (
        engineer_id,
        coverage_type,
        coverage_value,
        priority,
        effective_from,
        coverage_radius_km,
        assigned_by,
        assignment_reason
    )
    SELECT 
        e.id as engineer_id,
        'pincode' as coverage_type,
        unnest(e.coverage_pincodes) as coverage_value,
        1 as priority,
        COALESCE(e.joining_date, e.created_at::date, CURRENT_DATE) as effective_from,
        e.coverage_radius_km,
        'migration_T2.3' as assigned_by,
        'Migrated from engineers.coverage_pincodes array' as assignment_reason
    FROM engineers e
    WHERE e.coverage_pincodes IS NOT NULL 
        AND array_length(e.coverage_pincodes, 1) > 0
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_pincodes_migrated = ROW_COUNT;
    
    -- Migrate cities
    INSERT INTO engineer_coverage_areas (
        engineer_id,
        coverage_type,
        coverage_value,
        priority,
        effective_from,
        coverage_radius_km,
        assigned_by,
        assignment_reason
    )
    SELECT 
        e.id as engineer_id,
        'city' as coverage_type,
        unnest(e.coverage_cities) as coverage_value,
        1 as priority,
        COALESCE(e.joining_date, e.created_at::date, CURRENT_DATE) as effective_from,
        e.coverage_radius_km,
        'migration_T2.3' as assigned_by,
        'Migrated from engineers.coverage_cities array' as assignment_reason
    FROM engineers e
    WHERE e.coverage_cities IS NOT NULL 
        AND array_length(e.coverage_cities, 1) > 0
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_cities_migrated = ROW_COUNT;
    
    -- Migrate states
    INSERT INTO engineer_coverage_areas (
        engineer_id,
        coverage_type,
        coverage_value,
        priority,
        effective_from,
        coverage_radius_km,
        assigned_by,
        assignment_reason
    )
    SELECT 
        e.id as engineer_id,
        'state' as coverage_type,
        unnest(e.coverage_states) as coverage_value,
        1 as priority,
        COALESCE(e.joining_date, e.created_at::date, CURRENT_DATE) as effective_from,
        e.coverage_radius_km,
        'migration_T2.3' as assigned_by,
        'Migrated from engineers.coverage_states array' as assignment_reason
    FROM engineers e
    WHERE e.coverage_states IS NOT NULL 
        AND array_length(e.coverage_states, 1) > 0
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_states_migrated = ROW_COUNT;
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration Summary:';
    RAISE NOTICE '  - Pincode coverage records: %', v_pincodes_migrated;
    RAISE NOTICE '  - City coverage records: %', v_cities_migrated;
    RAISE NOTICE '  - State coverage records: %', v_states_migrated;
    RAISE NOTICE '  - Total coverage records: %', (v_pincodes_migrated + v_cities_migrated + v_states_migrated);
    RAISE NOTICE '============================================';
END $$;

-- =============================================================================
-- STEP 4: Create helper functions
-- =============================================================================

-- Function: Find available engineer for location
CREATE OR REPLACE FUNCTION find_engineer_for_location(
    p_coverage_type TEXT,
    p_coverage_value TEXT,
    p_priority INT DEFAULT 1
)
RETURNS TABLE (
    engineer_id UUID,
    engineer_name TEXT,
    engineer_status TEXT,
    active_tickets INT,
    priority INT,
    can_emergency BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        e.id as engineer_id,
        e.full_name as engineer_name,
        e.status as engineer_status,
        e.active_tickets,
        eca.priority,
        eca.can_emergency
    FROM engineers e
    JOIN engineer_coverage_areas eca ON e.id = eca.engineer_id
    WHERE eca.coverage_type = p_coverage_type
        AND eca.coverage_value = p_coverage_value
        AND eca.priority <= p_priority
        AND eca.effective_to IS NULL
        AND e.status = 'available'
    ORDER BY eca.priority ASC, e.active_tickets ASC;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function: Get engineer's current coverage
CREATE OR REPLACE FUNCTION get_engineer_current_coverage(p_engineer_id UUID)
RETURNS TABLE (
    coverage_type TEXT,
    coverage_value TEXT,
    priority INT,
    since_date DATE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        eca.coverage_type,
        eca.coverage_value,
        eca.priority,
        eca.effective_from as since_date
    FROM engineer_coverage_areas eca
    WHERE eca.engineer_id = p_engineer_id
        AND eca.effective_to IS NULL
    ORDER BY eca.coverage_type, eca.priority, eca.coverage_value;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function: Update engineer coverage (closes old, adds new)
CREATE OR REPLACE FUNCTION update_engineer_coverage(
    p_engineer_id UUID,
    p_coverage_type TEXT,
    p_new_coverage_values TEXT[],
    p_priority INT DEFAULT 1,
    p_effective_from DATE DEFAULT CURRENT_DATE,
    p_assigned_by TEXT DEFAULT 'system',
    p_reason TEXT DEFAULT NULL
)
RETURNS INTEGER AS $$
DECLARE
    v_inserted_count INTEGER := 0;
BEGIN
    -- Close existing coverage of this type
    UPDATE engineer_coverage_areas
    SET effective_to = p_effective_from,
        notes = COALESCE(notes || E'\n', '') || 'Coverage updated: ' || COALESCE(p_reason, 'Coverage change')
    WHERE engineer_id = p_engineer_id
        AND coverage_type = p_coverage_type
        AND effective_to IS NULL
        AND effective_from < p_effective_from;
    
    -- Add new coverage
    INSERT INTO engineer_coverage_areas (
        engineer_id,
        coverage_type,
        coverage_value,
        priority,
        effective_from,
        assigned_by,
        assignment_reason
    )
    SELECT
        p_engineer_id,
        p_coverage_type,
        unnest(p_new_coverage_values),
        p_priority,
        p_effective_from,
        p_assigned_by,
        p_reason;
    
    GET DIAGNOSTICS v_inserted_count = ROW_COUNT;
    
    RETURN v_inserted_count;
END;
$$ LANGUAGE plpgsql;

-- Function: Find coverage gaps
CREATE OR REPLACE FUNCTION find_coverage_gaps(
    p_coverage_type TEXT,
    p_required_values TEXT[]
)
RETURNS TABLE (
    coverage_value TEXT,
    has_primary BOOLEAN,
    has_secondary BOOLEAN,
    total_engineers INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        required_value as coverage_value,
        COUNT(*) FILTER (WHERE eca.priority = 1) > 0 as has_primary,
        COUNT(*) FILTER (WHERE eca.priority = 2) > 0 as has_secondary,
        COUNT(*)::INT as total_engineers
    FROM unnest(p_required_values) AS required_value
    LEFT JOIN engineer_coverage_areas eca ON (
        eca.coverage_type = p_coverage_type
        AND eca.coverage_value = required_value
        AND eca.effective_to IS NULL
    )
    GROUP BY required_value
    ORDER BY total_engineers ASC, required_value;
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================================================
-- STEP 5: Create views for backward compatibility
-- =============================================================================

-- View: Engineers with current coverage (backward compatible)
CREATE OR REPLACE VIEW engineers_with_coverage AS
SELECT
    e.*,
    -- Current pincodes
    array_agg(DISTINCT eca.coverage_value) FILTER (
        WHERE eca.coverage_type = 'pincode' AND eca.effective_to IS NULL
    ) as current_pincodes,
    -- Current cities
    array_agg(DISTINCT eca.coverage_value) FILTER (
        WHERE eca.coverage_type = 'city' AND eca.effective_to IS NULL
    ) as current_cities,
    -- Current states
    array_agg(DISTINCT eca.coverage_value) FILTER (
        WHERE eca.coverage_type = 'state' AND eca.effective_to IS NULL
    ) as current_states,
    -- Total coverage count
    COUNT(eca.id) FILTER (WHERE eca.effective_to IS NULL) as total_coverage_areas
FROM engineers e
LEFT JOIN engineer_coverage_areas eca ON e.id = eca.engineer_id
GROUP BY e.id;

COMMENT ON VIEW engineers_with_coverage IS 
    'Engineers with aggregated coverage arrays for backward compatibility';

-- View: Active coverage summary
CREATE OR REPLACE VIEW engineer_coverage_summary AS
SELECT
    e.id as engineer_id,
    e.full_name,
    e.org_id,
    e.status,
    COUNT(DISTINCT eca.id) FILTER (WHERE eca.coverage_type = 'pincode') as pincodes_covered,
    COUNT(DISTINCT eca.id) FILTER (WHERE eca.coverage_type = 'city') as cities_covered,
    COUNT(DISTINCT eca.id) FILTER (WHERE eca.coverage_type = 'state') as states_covered,
    COUNT(DISTINCT eca.id) FILTER (WHERE eca.priority = 1) as primary_areas,
    COUNT(DISTINCT eca.id) FILTER (WHERE eca.priority > 1) as backup_areas
FROM engineers e
LEFT JOIN engineer_coverage_areas eca ON (
    e.id = eca.engineer_id 
    AND eca.effective_to IS NULL
)
GROUP BY e.id, e.full_name, e.org_id, e.status;

COMMENT ON VIEW engineer_coverage_summary IS 
    'Quick summary of engineer coverage counts and priorities';

-- =============================================================================
-- STEP 6: Validation
-- =============================================================================

DO $$
DECLARE
    v_engineers_count INTEGER;
    v_coverage_count INTEGER;
    v_engineers_no_coverage INTEGER;
BEGIN
    -- Count engineers
    SELECT COUNT(*) INTO v_engineers_count FROM engineers;
    
    -- Count coverage records
    SELECT COUNT(*) INTO v_coverage_count 
    FROM engineer_coverage_areas 
    WHERE effective_to IS NULL;
    
    -- Count engineers without coverage
    SELECT COUNT(*) INTO v_engineers_no_coverage
    FROM engineers e
    WHERE NOT EXISTS (
        SELECT 1 FROM engineer_coverage_areas eca
        WHERE eca.engineer_id = e.id
            AND eca.effective_to IS NULL
    );
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Validation Results:';
    RAISE NOTICE '  - Total engineers: %', v_engineers_count;
    RAISE NOTICE '  - Active coverage records: %', v_coverage_count;
    RAISE NOTICE '  - Engineers without coverage: %', v_engineers_no_coverage;
    
    IF v_engineers_no_coverage > 0 THEN
        RAISE WARNING '% engineers have no coverage defined', v_engineers_no_coverage;
    ELSE
        RAISE NOTICE '✓ All engineers have coverage defined';
    END IF;
    
    RAISE NOTICE '============================================';
END $$;

-- =============================================================================
-- Migration Complete
-- =============================================================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'T2.3 Migration Complete!';
    RAISE NOTICE 'Engineer coverage normalized';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Features enabled:';
    RAISE NOTICE '  ✓ Normalized coverage table';
    RAISE NOTICE '  ✓ Fast "who covers X?" queries';
    RAISE NOTICE '  ✓ Temporal coverage tracking';
    RAISE NOTICE '  ✓ Priority-based assignment';
    RAISE NOTICE '  ✓ Territory integration ready';
    RAISE NOTICE '  ✓ Coverage history audit trail';
    RAISE NOTICE '  ✓ Backward compatible view';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Example usage:';
    RAISE NOTICE '  -- Find engineer for pincode:';
    RAISE NOTICE '  SELECT * FROM find_engineer_for_location(''pincode'', ''400001'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Get engineer coverage:';
    RAISE NOTICE '  SELECT * FROM get_engineer_current_coverage(''engineer-uuid'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Update coverage:';
    RAISE NOTICE '  SELECT update_engineer_coverage(';
    RAISE NOTICE '      ''engineer-uuid'',';
    RAISE NOTICE '      ''pincode'',';
    RAISE NOTICE '      ARRAY[''400001'', ''400002'', ''400003''],';
    RAISE NOTICE '      1,  -- priority';
    RAISE NOTICE '      CURRENT_DATE,';
    RAISE NOTICE '      ''manager@company.com'',';
    RAISE NOTICE '      ''Coverage area expansion''';
    RAISE NOTICE '  );';
    RAISE NOTICE '============================================';
END $$;
