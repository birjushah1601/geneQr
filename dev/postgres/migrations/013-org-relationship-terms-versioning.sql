-- Migration: Organization Relationship Terms Versioning
-- Ticket: T2.1
-- Purpose: Track business terms changes over time with full audit trail

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS btree_gist;  -- For EXCLUDE constraint

-- =============================================================================
-- STEP 1: Create org_relationship_terms table
-- =============================================================================

CREATE TABLE IF NOT EXISTS org_relationship_terms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Link to relationship
    relationship_id UUID NOT NULL REFERENCES org_relationships(id) ON DELETE CASCADE,
    version INT NOT NULL,
    
    -- Temporal validity
    effective_from TIMESTAMPTZ NOT NULL,
    effective_to TIMESTAMPTZ,
    
    -- Business terms (what changes over time)
    commission_percentage NUMERIC(5,2),
    payment_terms JSONB,
    credit_limit NUMERIC(18,2),
    credit_currency VARCHAR(3) DEFAULT 'INR',
    annual_target NUMERIC(18,2),
    annual_target_currency VARCHAR(3) DEFAULT 'INR',
    performance_tier TEXT,
    priority_level INT,
    
    -- Pricing rules
    discount_percentage NUMERIC(5,2),
    special_pricing_applicable BOOLEAN DEFAULT false,
    volume_incentives JSONB,
    
    -- Service commitments
    sla_terms JSONB,
    support_level TEXT,
    training_included BOOLEAN DEFAULT false,
    
    -- Audit trail
    changed_by TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    change_reason TEXT,
    approved_by TEXT,
    approval_date TIMESTAMPTZ,
    
    -- Metadata
    notes TEXT,
    metadata JSONB,
    
    -- Constraints
    CONSTRAINT effective_dates_valid CHECK (
        effective_to IS NULL OR effective_to >= effective_from
    ),
    CONSTRAINT version_positive CHECK (version > 0),
    CONSTRAINT chk_performance_tier CHECK (
        performance_tier IS NULL OR 
        performance_tier IN ('bronze', 'silver', 'gold', 'platinum', 'diamond')
    ),
    CONSTRAINT chk_support_level CHECK (
        support_level IS NULL OR
        support_level IN ('basic', 'standard', 'premium', 'enterprise')
    ),
    CONSTRAINT chk_priority_level CHECK (
        priority_level IS NULL OR
        (priority_level >= 1 AND priority_level <= 5)
    ),
    UNIQUE(relationship_id, version)
);

-- Add EXCLUDE constraint to prevent overlapping terms
ALTER TABLE org_relationship_terms
    ADD CONSTRAINT org_terms_no_overlap
    EXCLUDE USING gist (
        relationship_id WITH =,
        tstzrange(effective_from, COALESCE(effective_to, '9999-12-31'::timestamptz), '[]') WITH &&
    );

-- =============================================================================
-- STEP 2: Create indexes for performance
-- =============================================================================

-- Primary lookups
CREATE INDEX idx_org_terms_relationship ON org_relationship_terms(relationship_id);
CREATE INDEX idx_org_terms_version ON org_relationship_terms(relationship_id, version DESC);

-- Temporal queries
CREATE INDEX idx_org_terms_dates ON org_relationship_terms(effective_from, effective_to);
CREATE INDEX idx_org_terms_effective_from ON org_relationship_terms(effective_from);

-- Current terms (most common query)
CREATE INDEX idx_org_terms_current ON org_relationship_terms(relationship_id)
    WHERE effective_to IS NULL;

-- Business queries
CREATE INDEX idx_org_terms_tier ON org_relationship_terms(performance_tier)
    WHERE effective_to IS NULL;
CREATE INDEX idx_org_terms_changed_by ON org_relationship_terms(changed_by);

-- GIN indexes for JSONB
CREATE INDEX idx_org_terms_payment ON org_relationship_terms USING gin(payment_terms);
CREATE INDEX idx_org_terms_incentives ON org_relationship_terms USING gin(volume_incentives);
CREATE INDEX idx_org_terms_sla ON org_relationship_terms USING gin(sla_terms);

-- =============================================================================
-- STEP 3: Migrate existing data from org_relationships
-- =============================================================================

DO $$
DECLARE
    v_migrated_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Starting org relationship terms migration...';
    
    -- Migrate current terms as version 1
    INSERT INTO org_relationship_terms (
        relationship_id,
        version,
        effective_from,
        commission_percentage,
        payment_terms,
        credit_limit,
        annual_target,
        performance_tier,
        priority_level,
        changed_by,
        changed_at,
        change_reason
    )
    SELECT
        id as relationship_id,
        1 as version,
        COALESCE(start_date, created_at)::timestamptz as effective_from,
        commission_percentage,
        payment_terms,
        credit_limit,
        annual_target,
        performance_tier,
        priority_level,
        COALESCE(created_by, 'system') as changed_by,
        COALESCE(created_at, NOW()) as changed_at,
        'Migrated from org_relationships table - initial terms' as change_reason
    FROM org_relationships
    WHERE relationship_status = 'active'
        OR relationship_status IS NULL
    ON CONFLICT DO NOTHING;
    
    GET DIAGNOSTICS v_migrated_count = ROW_COUNT;
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration Summary:';
    RAISE NOTICE '  - Terms records migrated: %', v_migrated_count;
    RAISE NOTICE '============================================';
END $$;

-- =============================================================================
-- STEP 4: Create helper functions
-- =============================================================================

-- Function to get current terms for a relationship
CREATE OR REPLACE FUNCTION get_current_terms(p_relationship_id UUID)
RETURNS TABLE (
    version INT,
    commission_percentage NUMERIC,
    credit_limit NUMERIC,
    annual_target NUMERIC,
    performance_tier TEXT,
    effective_from TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        t.version,
        t.commission_percentage,
        t.credit_limit,
        t.annual_target,
        t.performance_tier,
        t.effective_from
    FROM org_relationship_terms t
    WHERE t.relationship_id = p_relationship_id
        AND t.effective_to IS NULL
    ORDER BY t.version DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function to get terms at specific timestamp (temporal query)
CREATE OR REPLACE FUNCTION get_terms_at_timestamp(
    p_relationship_id UUID,
    p_timestamp TIMESTAMPTZ
)
RETURNS TABLE (
    version INT,
    commission_percentage NUMERIC,
    credit_limit NUMERIC,
    annual_target NUMERIC,
    performance_tier TEXT,
    effective_from TIMESTAMPTZ,
    effective_to TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        t.version,
        t.commission_percentage,
        t.credit_limit,
        t.annual_target,
        t.performance_tier,
        t.effective_from,
        t.effective_to
    FROM org_relationship_terms t
    WHERE t.relationship_id = p_relationship_id
        AND p_timestamp >= t.effective_from
        AND (t.effective_to IS NULL OR p_timestamp <= t.effective_to)
    ORDER BY t.version DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function to update relationship terms (creates new version)
CREATE OR REPLACE FUNCTION update_relationship_terms(
    p_relationship_id UUID,
    p_commission_percentage NUMERIC DEFAULT NULL,
    p_payment_terms JSONB DEFAULT NULL,
    p_credit_limit NUMERIC DEFAULT NULL,
    p_annual_target NUMERIC DEFAULT NULL,
    p_performance_tier TEXT DEFAULT NULL,
    p_priority_level INT DEFAULT NULL,
    p_discount_percentage NUMERIC DEFAULT NULL,
    p_effective_from TIMESTAMPTZ DEFAULT NOW(),
    p_changed_by TEXT DEFAULT 'system',
    p_change_reason TEXT DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    v_new_version INT;
    v_new_id UUID;
    v_old_terms RECORD;
BEGIN
    -- Get current terms to copy unchanged values
    SELECT * INTO v_old_terms
    FROM org_relationship_terms
    WHERE relationship_id = p_relationship_id
        AND effective_to IS NULL
    LIMIT 1;
    
    -- Close current terms
    UPDATE org_relationship_terms
    SET effective_to = p_effective_from,
        notes = COALESCE(notes || E'\n', '') || 'Superseded by version ' || (v_old_terms.version + 1)::text
    WHERE relationship_id = p_relationship_id
        AND effective_to IS NULL;
    
    -- Get next version number
    SELECT COALESCE(MAX(version), 0) + 1 INTO v_new_version
    FROM org_relationship_terms
    WHERE relationship_id = p_relationship_id;
    
    -- Create new version (use new values if provided, else copy from old)
    INSERT INTO org_relationship_terms (
        relationship_id,
        version,
        effective_from,
        commission_percentage,
        payment_terms,
        credit_limit,
        credit_currency,
        annual_target,
        annual_target_currency,
        performance_tier,
        priority_level,
        discount_percentage,
        special_pricing_applicable,
        volume_incentives,
        sla_terms,
        support_level,
        training_included,
        changed_by,
        changed_at,
        change_reason,
        metadata
    ) VALUES (
        p_relationship_id,
        v_new_version,
        p_effective_from,
        COALESCE(p_commission_percentage, v_old_terms.commission_percentage),
        COALESCE(p_payment_terms, v_old_terms.payment_terms),
        COALESCE(p_credit_limit, v_old_terms.credit_limit),
        v_old_terms.credit_currency,
        COALESCE(p_annual_target, v_old_terms.annual_target),
        v_old_terms.annual_target_currency,
        COALESCE(p_performance_tier, v_old_terms.performance_tier),
        COALESCE(p_priority_level, v_old_terms.priority_level),
        COALESCE(p_discount_percentage, v_old_terms.discount_percentage),
        v_old_terms.special_pricing_applicable,
        v_old_terms.volume_incentives,
        v_old_terms.sla_terms,
        v_old_terms.support_level,
        v_old_terms.training_included,
        p_changed_by,
        NOW(),
        p_change_reason,
        v_old_terms.metadata
    )
    RETURNING id INTO v_new_id;
    
    RAISE NOTICE 'Created version % for relationship %', v_new_version, p_relationship_id;
    
    RETURN v_new_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get complete terms history
CREATE OR REPLACE FUNCTION get_terms_history(p_relationship_id UUID)
RETURNS TABLE (
    version INT,
    effective_from TIMESTAMPTZ,
    effective_to TIMESTAMPTZ,
    commission_percentage NUMERIC,
    credit_limit NUMERIC,
    annual_target NUMERIC,
    performance_tier TEXT,
    changed_by TEXT,
    changed_at TIMESTAMPTZ,
    change_reason TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        t.version,
        t.effective_from,
        t.effective_to,
        t.commission_percentage,
        t.credit_limit,
        t.annual_target,
        t.performance_tier,
        t.changed_by,
        t.changed_at,
        t.change_reason
    FROM org_relationship_terms t
    WHERE t.relationship_id = p_relationship_id
    ORDER BY t.version DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================================================
-- STEP 5: Create views for backward compatibility
-- =============================================================================

-- View: Relationships with current terms
CREATE OR REPLACE VIEW org_relationships_with_current_terms AS
SELECT
    r.*,
    t.version as terms_version,
    t.commission_percentage as current_commission_percentage,
    t.payment_terms as current_payment_terms,
    t.credit_limit as current_credit_limit,
    t.annual_target as current_annual_target,
    t.performance_tier as current_performance_tier,
    t.priority_level as current_priority_level,
    t.discount_percentage as current_discount_percentage,
    t.effective_from as terms_effective_from,
    t.changed_by as terms_changed_by,
    t.changed_at as terms_changed_at
FROM org_relationships r
LEFT JOIN org_relationship_terms t ON (
    r.id = t.relationship_id
    AND t.effective_to IS NULL
);

COMMENT ON VIEW org_relationships_with_current_terms IS 
    'Organization relationships joined with their current terms (for backward compatibility)';

-- =============================================================================
-- STEP 6: Add comments for documentation
-- =============================================================================

COMMENT ON TABLE org_relationship_terms IS 
    'Version-controlled business terms for organization relationships. Tracks commission, credit, targets over time.';

COMMENT ON COLUMN org_relationship_terms.relationship_id IS 
    'Foreign key to org_relationships table';

COMMENT ON COLUMN org_relationship_terms.version IS 
    'Sequential version number (1, 2, 3, ...) for each relationship';

COMMENT ON COLUMN org_relationship_terms.effective_from IS 
    'Start timestamp when these terms become effective';

COMMENT ON COLUMN org_relationship_terms.effective_to IS 
    'End timestamp when these terms are superseded. NULL means currently active.';

COMMENT ON COLUMN org_relationship_terms.commission_percentage IS 
    'Commission percentage (e.g., 10.00 for 10%)';

COMMENT ON COLUMN org_relationship_terms.payment_terms IS 
    'JSONB containing payment terms: {type: "net_days", days: 30, discount_for_early: 2}';

COMMENT ON COLUMN org_relationship_terms.volume_incentives IS 
    'JSONB containing volume-based incentive structure';

COMMENT ON COLUMN org_relationship_terms.changed_by IS 
    'Email or ID of person who made this change';

COMMENT ON COLUMN org_relationship_terms.change_reason IS 
    'Human-readable reason for term change (e.g., "Annual review", "Performance upgrade")';

-- =============================================================================
-- STEP 7: Validation queries
-- =============================================================================

DO $$
DECLARE
    v_relationships_count INTEGER;
    v_terms_count INTEGER;
    v_missing_terms_count INTEGER;
BEGIN
    -- Count relationships
    SELECT COUNT(*) INTO v_relationships_count 
    FROM org_relationships 
    WHERE relationship_status = 'active';
    
    -- Count terms records
    SELECT COUNT(*) INTO v_terms_count 
    FROM org_relationship_terms;
    
    -- Check for relationships without terms
    SELECT COUNT(*) INTO v_missing_terms_count
    FROM org_relationships r
    WHERE r.relationship_status = 'active'
        AND NOT EXISTS (
            SELECT 1 FROM org_relationship_terms t
            WHERE t.relationship_id = r.id
        );
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Validation Results:';
    RAISE NOTICE '  - Active relationships: %', v_relationships_count;
    RAISE NOTICE '  - Terms records: %', v_terms_count;
    RAISE NOTICE '  - Relationships missing terms: %', v_missing_terms_count;
    
    IF v_missing_terms_count > 0 THEN
        RAISE WARNING '% active relationships have no terms', v_missing_terms_count;
    ELSE
        RAISE NOTICE '✓ All active relationships have terms';
    END IF;
    
    RAISE NOTICE '============================================';
END $$;

-- =============================================================================
-- Migration Complete
-- =============================================================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'T2.1 Migration Complete!';
    RAISE NOTICE 'Organization relationship terms versioning enabled';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Features enabled:';
    RAISE NOTICE '  ✓ Version-controlled business terms';
    RAISE NOTICE '  ✓ Temporal queries (terms at specific date)';
    RAISE NOTICE '  ✓ Complete change history';
    RAISE NOTICE '  ✓ Audit trail with reasons';
    RAISE NOTICE '  ✓ Automatic version management';
    RAISE NOTICE '  ✓ No-overlap constraint';
    RAISE NOTICE '  ✓ Backward compatible view';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Example usage:';
    RAISE NOTICE '  -- Get current terms:';
    RAISE NOTICE '  SELECT * FROM get_current_terms(''relationship-uuid'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Get terms on specific date:';
    RAISE NOTICE '  SELECT * FROM get_terms_at_timestamp(''relationship-uuid'', ''2024-01-15 00:00:00+00'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Update terms (commission upgrade):';
    RAISE NOTICE '  SELECT update_relationship_terms(';
    RAISE NOTICE '      ''relationship-uuid'',';
    RAISE NOTICE '      12.0,  -- new commission';
    RAISE NOTICE '      NULL,  -- payment terms unchanged';
    RAISE NOTICE '      NULL,  -- credit unchanged';
    RAISE NOTICE '      NULL,  -- target unchanged';
    RAISE NOTICE '      NULL,  -- tier unchanged';
    RAISE NOTICE '      NULL,  -- priority unchanged';
    RAISE NOTICE '      NULL,  -- discount unchanged';
    RAISE NOTICE '      NOW(), -- effective now';
    RAISE NOTICE '      ''manager@company.com'',';
    RAISE NOTICE '      ''Performance-based upgrade''';
    RAISE NOTICE '  );';
    RAISE NOTICE '============================================';
END $$;
