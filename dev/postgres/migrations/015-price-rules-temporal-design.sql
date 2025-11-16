-- Migration: Price Rules Temporal Design
-- Ticket: T2.4
-- Purpose: Add version control and temporal tracking to pricing

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS btree_gist;  -- For EXCLUDE constraint

-- =============================================================================
-- STEP 1: Remove blocking constraint
-- =============================================================================

-- Drop UNIQUE(book_id, sku_id) constraint that prevents versioning
ALTER TABLE price_rules 
DROP CONSTRAINT IF EXISTS price_rules_book_id_sku_id_key;

RAISE NOTICE 'Removed UNIQUE constraint blocking price versioning';

-- =============================================================================
-- STEP 2: Add version control columns
-- =============================================================================

-- Add versioning and audit trail columns
ALTER TABLE price_rules
ADD COLUMN IF NOT EXISTS version INT DEFAULT 1,
ADD COLUMN IF NOT EXISTS pricing_type TEXT DEFAULT 'regular',
ADD COLUMN IF NOT EXISTS min_quantity INT,
ADD COLUMN IF NOT EXISTS max_quantity INT,
ADD COLUMN IF NOT EXISTS discount_percentage NUMERIC(5,2),
ADD COLUMN IF NOT EXISTS changed_by TEXT,
ADD COLUMN IF NOT EXISTS changed_at TIMESTAMPTZ DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS change_reason TEXT,
ADD COLUMN IF NOT EXISTS approved_by TEXT,
ADD COLUMN IF NOT EXISTS approval_date TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS notes TEXT,
ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Update existing records
UPDATE price_rules
SET 
    version = 1,
    pricing_type = 'regular',
    changed_by = 'migration_T2.4',
    changed_at = NOW(),
    change_reason = 'Initial price version from migration'
WHERE version IS NULL OR version = 0;

-- Make version NOT NULL
ALTER TABLE price_rules ALTER COLUMN version SET NOT NULL;
ALTER TABLE price_rules ALTER COLUMN version SET DEFAULT 1;

-- Add constraints on new columns
ALTER TABLE price_rules
ADD CONSTRAINT chk_version_positive CHECK (version > 0),
ADD CONSTRAINT chk_price_positive CHECK (price > 0),
ADD CONSTRAINT chk_pricing_type CHECK (
    pricing_type IS NULL OR 
    pricing_type IN ('regular', 'promotional', 'seasonal', 'contract', 'volume', 'clearance')
),
ADD CONSTRAINT chk_valid_dates CHECK (
    valid_to IS NULL OR valid_to >= valid_from
),
ADD CONSTRAINT chk_quantity_range CHECK (
    (min_quantity IS NULL AND max_quantity IS NULL) OR
    (min_quantity IS NOT NULL AND (max_quantity IS NULL OR max_quantity >= min_quantity))
);

-- =============================================================================
-- STEP 3: Add new constraints
-- =============================================================================

-- Add version-based unique constraint
ALTER TABLE price_rules
ADD CONSTRAINT price_rules_version_unique 
UNIQUE(book_id, sku_id, version);

-- Add no-overlap constraint to enforce temporal integrity
ALTER TABLE price_rules
ADD CONSTRAINT price_rules_no_overlap
EXCLUDE USING gist (
    book_id WITH =,
    sku_id WITH =,
    tstzrange(valid_from, COALESCE(valid_to, '9999-12-31'::timestamptz), '[]') WITH &&
);

RAISE NOTICE 'Added version control and temporal integrity constraints';

-- =============================================================================
-- STEP 4: Create indexes for performance
-- =============================================================================

-- Primary lookups
CREATE INDEX IF NOT EXISTS idx_price_rules_book_sku ON price_rules(book_id, sku_id);
CREATE INDEX IF NOT EXISTS idx_price_rules_version ON price_rules(book_id, sku_id, version DESC);

-- Current prices (most common query)
CREATE INDEX IF NOT EXISTS idx_price_rules_current ON price_rules(book_id, sku_id)
    WHERE valid_to IS NULL;

-- Temporal queries
CREATE INDEX IF NOT EXISTS idx_price_rules_dates ON price_rules(valid_from, valid_to);
CREATE INDEX IF NOT EXISTS idx_price_rules_valid_from ON price_rules(valid_from);

-- Pricing type queries
CREATE INDEX IF NOT EXISTS idx_price_rules_type ON price_rules(pricing_type)
    WHERE valid_to IS NULL;

-- Promotional pricing
CREATE INDEX IF NOT EXISTS idx_price_rules_promo ON price_rules(book_id, pricing_type)
    WHERE pricing_type = 'promotional' AND valid_to IS NULL;

-- =============================================================================
-- STEP 5: Create helper functions
-- =============================================================================

-- Function: Get current price for SKU
CREATE OR REPLACE FUNCTION get_current_price(
    p_book_id UUID,
    p_sku_id UUID
)
RETURNS TABLE (
    version INT,
    price NUMERIC,
    currency TEXT,
    pricing_type TEXT,
    valid_from TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        pr.version,
        pr.price,
        pr.currency,
        pr.pricing_type,
        pr.valid_from
    FROM price_rules pr
    WHERE pr.book_id = p_book_id
        AND pr.sku_id = p_sku_id
        AND pr.valid_to IS NULL
    ORDER BY pr.version DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function: Get price at specific date (temporal query)
CREATE OR REPLACE FUNCTION get_price_at_date(
    p_book_id UUID,
    p_sku_id UUID,
    p_date TIMESTAMPTZ
)
RETURNS TABLE (
    version INT,
    price NUMERIC,
    currency TEXT,
    pricing_type TEXT,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    change_reason TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        pr.version,
        pr.price,
        pr.currency,
        pr.pricing_type,
        pr.valid_from,
        pr.valid_to,
        pr.change_reason
    FROM price_rules pr
    WHERE pr.book_id = p_book_id
        AND pr.sku_id = p_sku_id
        AND p_date >= pr.valid_from
        AND (pr.valid_to IS NULL OR p_date <= pr.valid_to)
    ORDER BY pr.version DESC
    LIMIT 1;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function: Get complete price history
CREATE OR REPLACE FUNCTION get_price_history(
    p_book_id UUID,
    p_sku_id UUID
)
RETURNS TABLE (
    version INT,
    price NUMERIC,
    currency TEXT,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    pricing_type TEXT,
    discount_percentage NUMERIC,
    changed_by TEXT,
    changed_at TIMESTAMPTZ,
    change_reason TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        pr.version,
        pr.price,
        pr.currency,
        pr.valid_from,
        pr.valid_to,
        pr.pricing_type,
        pr.discount_percentage,
        pr.changed_by,
        pr.changed_at,
        pr.change_reason
    FROM price_rules pr
    WHERE pr.book_id = p_book_id
        AND pr.sku_id = p_sku_id
    ORDER BY pr.version DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- Function: Update price (creates new version)
CREATE OR REPLACE FUNCTION update_price(
    p_book_id UUID,
    p_sku_id UUID,
    p_new_price NUMERIC,
    p_valid_from TIMESTAMPTZ DEFAULT NOW(),
    p_pricing_type TEXT DEFAULT 'regular',
    p_changed_by TEXT DEFAULT 'system',
    p_change_reason TEXT DEFAULT NULL,
    p_currency TEXT DEFAULT 'INR'
)
RETURNS UUID AS $$
DECLARE
    v_new_version INT;
    v_new_id UUID;
BEGIN
    -- Close current price
    UPDATE price_rules
    SET valid_to = p_valid_from,
        notes = COALESCE(notes || E'\n', '') || 'Superseded by version ' || (
            SELECT COALESCE(MAX(version), 0) + 1 
            FROM price_rules 
            WHERE book_id = p_book_id AND sku_id = p_sku_id
        )::text
    WHERE book_id = p_book_id
        AND sku_id = p_sku_id
        AND valid_to IS NULL
        AND valid_from < p_valid_from;
    
    -- Get next version number
    SELECT COALESCE(MAX(version), 0) + 1 INTO v_new_version
    FROM price_rules
    WHERE book_id = p_book_id AND sku_id = p_sku_id;
    
    -- Insert new version
    INSERT INTO price_rules (
        book_id,
        sku_id,
        version,
        price,
        currency,
        valid_from,
        pricing_type,
        changed_by,
        changed_at,
        change_reason
    ) VALUES (
        p_book_id,
        p_sku_id,
        v_new_version,
        p_new_price,
        p_currency,
        p_valid_from,
        p_pricing_type,
        p_changed_by,
        NOW(),
        p_change_reason
    )
    RETURNING id INTO v_new_id;
    
    RAISE NOTICE 'Created price version % for SKU % in book %', v_new_version, p_sku_id, p_book_id;
    
    RETURN v_new_id;
END;
$$ LANGUAGE plpgsql;

-- Function: Schedule promotional price (with end date)
CREATE OR REPLACE FUNCTION schedule_promotion(
    p_book_id UUID,
    p_sku_id UUID,
    p_promo_price NUMERIC,
    p_discount_pct NUMERIC,
    p_start_date TIMESTAMPTZ,
    p_end_date TIMESTAMPTZ,
    p_changed_by TEXT,
    p_reason TEXT
)
RETURNS UUID AS $$
DECLARE
    v_new_version INT;
    v_new_id UUID;
    v_regular_price NUMERIC;
BEGIN
    -- Get current regular price
    SELECT price INTO v_regular_price
    FROM price_rules
    WHERE book_id = p_book_id
        AND sku_id = p_sku_id
        AND valid_to IS NULL
        AND pricing_type = 'regular';
    
    IF v_regular_price IS NULL THEN
        RAISE EXCEPTION 'No regular price found for SKU %', p_sku_id;
    END IF;
    
    -- Get next version
    SELECT COALESCE(MAX(version), 0) + 1 INTO v_new_version
    FROM price_rules
    WHERE book_id = p_book_id AND sku_id = p_sku_id;
    
    -- Insert promotional price
    INSERT INTO price_rules (
        book_id,
        sku_id,
        version,
        price,
        valid_from,
        valid_to,
        pricing_type,
        discount_percentage,
        changed_by,
        changed_at,
        change_reason
    ) VALUES (
        p_book_id,
        p_sku_id,
        v_new_version,
        p_promo_price,
        p_start_date,
        p_end_date,
        'promotional',
        p_discount_pct,
        p_changed_by,
        NOW(),
        p_reason
    )
    RETURNING id INTO v_new_id;
    
    RAISE NOTICE 'Scheduled promotion v% for SKU % from % to %', 
        v_new_version, p_sku_id, p_start_date, p_end_date;
    
    RETURN v_new_id;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- STEP 6: Create views for backward compatibility
-- =============================================================================

-- View: Current prices only
CREATE OR REPLACE VIEW current_prices AS
SELECT
    pr.id,
    pr.book_id,
    pb.name as book_name,
    pb.org_id,
    pb.channel_id,
    pr.sku_id,
    pr.version,
    pr.price,
    pr.currency,
    pr.pricing_type,
    pr.discount_percentage,
    pr.valid_from,
    pr.changed_by,
    pr.changed_at
FROM price_rules pr
JOIN price_books pb ON pr.book_id = pb.id
WHERE pr.valid_to IS NULL;

COMMENT ON VIEW current_prices IS 
    'Current active prices across all price books';

-- View: Promotional prices
CREATE OR REPLACE VIEW active_promotions AS
SELECT
    pr.id,
    pr.book_id,
    pb.name as book_name,
    pr.sku_id,
    pr.price as promo_price,
    pr.discount_percentage,
    pr.valid_from as promo_start,
    pr.valid_to as promo_end,
    pr.change_reason as promo_name,
    pr.changed_by
FROM price_rules pr
JOIN price_books pb ON pr.book_id = pb.id
WHERE pr.pricing_type = 'promotional'
    AND pr.valid_to IS NULL;

COMMENT ON VIEW active_promotions IS 
    'Currently active promotional prices';

-- View: Upcoming price changes
CREATE OR REPLACE VIEW upcoming_price_changes AS
SELECT
    pr.id,
    pr.book_id,
    pb.name as book_name,
    pr.sku_id,
    pr.price as new_price,
    pr.valid_from as effective_date,
    pr.pricing_type,
    pr.change_reason,
    pr.changed_by,
    (pr.valid_from - NOW()) as time_until_change
FROM price_rules pr
JOIN price_books pb ON pr.book_id = pb.id
WHERE pr.valid_from > NOW()
    AND pr.valid_to IS NULL
ORDER BY pr.valid_from ASC;

COMMENT ON VIEW upcoming_price_changes IS 
    'Future scheduled price changes';

-- =============================================================================
-- STEP 7: Add comments for documentation
-- =============================================================================

COMMENT ON TABLE price_rules IS 
    'Version-controlled price rules with temporal tracking. Supports promotional pricing, future scheduling, and complete audit trail.';

COMMENT ON COLUMN price_rules.version IS 
    'Sequential version number (1, 2, 3, ...) for price changes';

COMMENT ON COLUMN price_rules.pricing_type IS 
    'Type of pricing: regular, promotional, seasonal, contract, volume, clearance';

COMMENT ON COLUMN price_rules.valid_from IS 
    'Start timestamp when this price becomes effective';

COMMENT ON COLUMN price_rules.valid_to IS 
    'End timestamp when this price is superseded. NULL means currently active.';

COMMENT ON COLUMN price_rules.discount_percentage IS 
    'Discount percentage for promotional pricing';

COMMENT ON COLUMN price_rules.min_quantity IS 
    'Minimum quantity for volume-based pricing';

COMMENT ON COLUMN price_rules.max_quantity IS 
    'Maximum quantity for volume-based pricing (NULL = no limit)';

-- =============================================================================
-- STEP 8: Validation
-- =============================================================================

DO $$
DECLARE
    v_price_books_count INTEGER;
    v_price_rules_count INTEGER;
    v_current_prices_count INTEGER;
BEGIN
    -- Count price books
    SELECT COUNT(*) INTO v_price_books_count FROM price_books;
    
    -- Count price rules
    SELECT COUNT(*) INTO v_price_rules_count FROM price_rules;
    
    -- Count current prices
    SELECT COUNT(*) INTO v_current_prices_count 
    FROM price_rules 
    WHERE valid_to IS NULL;
    
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Validation Results:';
    RAISE NOTICE '  - Price books: %', v_price_books_count;
    RAISE NOTICE '  - Total price rules: %', v_price_rules_count;
    RAISE NOTICE '  - Current active prices: %', v_current_prices_count;
    RAISE NOTICE '============================================';
END $$;

-- =============================================================================
-- Migration Complete
-- =============================================================================

DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'T2.4 Migration Complete!';
    RAISE NOTICE 'Price rules temporal design enabled';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Features enabled:';
    RAISE NOTICE '  ✓ Version-controlled pricing';
    RAISE NOTICE '  ✓ Temporal price queries';
    RAISE NOTICE '  ✓ Future price scheduling';
    RAISE NOTICE '  ✓ Promotional pricing support';
    RAISE NOTICE '  ✓ Complete price history';
    RAISE NOTICE '  ✓ No-overlap constraint';
    RAISE NOTICE '  ✓ Audit trail with reasons';
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Example usage:';
    RAISE NOTICE '  -- Get current price:';
    RAISE NOTICE '  SELECT * FROM get_current_price(''book-uuid'', ''sku-uuid'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Get price on specific date:';
    RAISE NOTICE '  SELECT * FROM get_price_at_date(''book-uuid'', ''sku-uuid'', ''2024-06-15 00:00:00+00'');';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Update price:';
    RAISE NOTICE '  SELECT update_price(';
    RAISE NOTICE '      ''book-uuid'',';
    RAISE NOTICE '      ''sku-uuid'',';
    RAISE NOTICE '      15000.00,  -- new price';
    RAISE NOTICE '      NOW(),     -- effective from';
    RAISE NOTICE '      ''regular'',';
    RAISE NOTICE '      ''pricing_manager@company.com'',';
    RAISE NOTICE '      ''Q4 2024 price update''';
    RAISE NOTICE '  );';
    RAISE NOTICE '';
    RAISE NOTICE '  -- Schedule promotion:';
    RAISE NOTICE '  SELECT schedule_promotion(';
    RAISE NOTICE '      ''book-uuid'', ''sku-uuid'',';
    RAISE NOTICE '      12000.00,  -- promo price';
    RAISE NOTICE '      20.00,     -- 20% discount';
    RAISE NOTICE '      ''2024-10-15''::timestamptz,';
    RAISE NOTICE '      ''2024-10-31''::timestamptz,';
    RAISE NOTICE '      ''marketing@company.com'',';
    RAISE NOTICE '      ''Diwali Festival Sale''';
    RAISE NOTICE '  );';
    RAISE NOTICE '============================================';
END $$;
