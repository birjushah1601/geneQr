-- Migration: Create customers table and normalize customer data
-- Ticket: T1.2
-- Purpose: Create proper customer entity and eliminate data duplication

-- Enable pg_trgm extension for full-text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- =============================================================================
-- STEP 1: Create customers table
-- =============================================================================

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic Information
    customer_type VARCHAR(50) NOT NULL,
    name VARCHAR(500) NOT NULL,
    display_name VARCHAR(255),
    
    -- Contact Information
    primary_phone VARCHAR(20) NOT NULL,
    secondary_phone VARCHAR(20),
    primary_email VARCHAR(255),
    whatsapp_number VARCHAR(20),
    
    -- Address
    address_line1 VARCHAR(500),
    address_line2 VARCHAR(500),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'India',
    
    -- Organization Link
    organization_id UUID,
    
    -- Business Details
    registration_number VARCHAR(100),
    tax_id VARCHAR(50),
    website VARCHAR(255),
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    tier VARCHAR(50),
    
    -- Preferences
    preferred_language VARCHAR(20) DEFAULT 'en',
    preferred_contact_method VARCHAR(50),
    
    -- Metadata
    notes TEXT,
    tags JSONB DEFAULT '[]'::jsonb,
    
    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    
    CONSTRAINT customer_type_check CHECK (customer_type IN ('individual', 'hospital', 'clinic', 'lab', 'manufacturer', 'Sub-sub_SUB_DEALER', 'other')),
    CONSTRAINT customer_status_check CHECK (status IN ('active', 'inactive', 'blocked'))
);

-- =============================================================================
-- STEP 2: Create indexes for performance
-- =============================================================================

-- Primary lookups
CREATE INDEX idx_customers_type ON customers(customer_type);
CREATE INDEX idx_customers_phone ON customers(primary_phone);
CREATE INDEX idx_customers_email ON customers(primary_email) WHERE primary_email IS NOT NULL;
CREATE INDEX idx_customers_status ON customers(status) WHERE status = 'active';

-- Geographic queries
CREATE INDEX idx_customers_city_state ON customers(city, state) WHERE city IS NOT NULL;

-- Organization relationship
CREATE INDEX idx_customers_org ON customers(organization_id) WHERE organization_id IS NOT NULL;

-- Full-text search on name
CREATE INDEX idx_customers_name_trgm ON customers USING gin(name gin_trgm_ops);

-- Unique constraint: same phone can't belong to multiple active customers
CREATE UNIQUE INDEX idx_customers_unique_phone ON customers(primary_phone) WHERE status = 'active';

-- =============================================================================
-- STEP 3: Migrate existing customer data from service_tickets
-- =============================================================================

DO $$
DECLARE
    v_migrated_count INTEGER := 0;
    v_new_customers_count INTEGER := 0;
    v_updated_tickets_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'Starting customer data migration...';
    
    -- 3a. Extract customers with existing customer_id
    INSERT INTO customers (
        id,
        customer_type,
        name,
        primary_phone,
        whatsapp_number,
        created_at,
        created_by
    )
    SELECT DISTINCT
        customer_id::uuid,
        'hospital' as customer_type,  -- Default type, can be updated
        customer_name,
        COALESCE(customer_phone, 'UNKNOWN'),
        customer_whatsapp,
        MIN(created_at) as created_at,
        'migration_T1.2' as created_by
    FROM service_tickets
    WHERE customer_id IS NOT NULL 
        AND customer_id != ''
        AND customer_name IS NOT NULL
        AND customer_name != ''
    GROUP BY customer_id, customer_name, customer_phone, customer_whatsapp
    ON CONFLICT (id) DO NOTHING;
    
    GET DIAGNOSTICS v_migrated_count = ROW_COUNT;
    RAISE NOTICE 'Migrated % existing customers', v_migrated_count;
    
    -- 3b. Create new customers for tickets without customer_id
    WITH new_customers AS (
        SELECT DISTINCT
            customer_name,
            customer_phone,
            customer_whatsapp,
            MIN(created_at) as first_ticket_date
        FROM service_tickets
        WHERE (customer_id IS NULL OR customer_id = '')
            AND customer_name IS NOT NULL
            AND customer_name != ''
        GROUP BY customer_name, customer_phone, customer_whatsapp
    )
    INSERT INTO customers (
        customer_type,
        name,
        primary_phone,
        whatsapp_number,
        created_at,
        created_by
    )
    SELECT
        'hospital' as customer_type,
        customer_name,
        COALESCE(customer_phone, 'UNKNOWN'),
        customer_whatsapp,
        first_ticket_date,
        'migration_T1.2' as created_by
    FROM new_customers;
    
    GET DIAGNOSTICS v_new_customers_count = ROW_COUNT;
    RAISE NOTICE 'Created % new customers', v_new_customers_count;
    
    -- 3c. Update tickets to reference new customer IDs
    UPDATE service_tickets st
    SET customer_id = c.id::text
    FROM customers c
    WHERE st.customer_name = c.name
        AND COALESCE(st.customer_phone, 'UNKNOWN') = c.primary_phone
        AND (st.customer_id IS NULL OR st.customer_id = '');
    
    GET DIAGNOSTICS v_updated_tickets_count = ROW_COUNT;
    RAISE NOTICE 'Updated % tickets with customer references', v_updated_tickets_count;
    
    -- Summary
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Migration Summary:';
    RAISE NOTICE '  - Existing customers migrated: %', v_migrated_count;
    RAISE NOTICE '  - New customers created: %', v_new_customers_count;
    RAISE NOTICE '  - Tickets updated: %', v_updated_tickets_count;
    RAISE NOTICE '==============================================';
END $$;

-- =============================================================================
-- STEP 4: Validation queries
-- =============================================================================

-- Check for orphaned customer_ids in service_tickets
DO $$
DECLARE
    v_orphaned_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_orphaned_count
    FROM service_tickets st
    LEFT JOIN customers c ON st.customer_id = c.id::text
    WHERE st.customer_id IS NOT NULL 
        AND st.customer_id != ''
        AND c.id IS NULL;
    
    IF v_orphaned_count > 0 THEN
        RAISE WARNING 'Found % orphaned customer references in service_tickets', v_orphaned_count;
    ELSE
        RAISE NOTICE 'âœ“ No orphaned customer references found';
    END IF;
END $$;

-- Check for duplicate customers
DO $$
DECLARE
    v_duplicate_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_duplicate_count
    FROM (
        SELECT primary_phone, COUNT(*) as cnt
        FROM customers
        WHERE status = 'active'
        GROUP BY primary_phone
        HAVING COUNT(*) > 1
    ) duplicates;
    
    IF v_duplicate_count > 0 THEN
        RAISE WARNING 'Found % duplicate customer phone numbers', v_duplicate_count;
    ELSE
        RAISE NOTICE 'âœ“ No duplicate customers found';
    END IF;
END $$;

-- =============================================================================
-- STEP 5: Deprecate old denormalized columns (gradual migration)
-- =============================================================================

-- Rename columns instead of dropping for backward compatibility
DO $$
BEGIN
    -- Only rename if columns exist and haven't been renamed already
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'service_tickets' AND column_name = 'customer_name') THEN
        ALTER TABLE service_tickets RENAME COLUMN customer_name TO customer_name_deprecated;
        RAISE NOTICE 'Renamed customer_name to customer_name_deprecated';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'service_tickets' AND column_name = 'customer_phone') THEN
        ALTER TABLE service_tickets RENAME COLUMN customer_phone TO customer_phone_deprecated;
        RAISE NOTICE 'Renamed customer_phone to customer_phone_deprecated';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'service_tickets' AND column_name = 'customer_whatsapp') THEN
        ALTER TABLE service_tickets RENAME COLUMN customer_whatsapp TO customer_whatsapp_deprecated;
        RAISE NOTICE 'Renamed customer_whatsapp to customer_whatsapp_deprecated';
    END IF;
END $$;

-- Add deprecation comments
COMMENT ON COLUMN service_tickets.customer_name_deprecated IS 
    'DEPRECATED: Use customers.name via customer_id foreign key. Will be dropped after 30 days.';
COMMENT ON COLUMN service_tickets.customer_phone_deprecated IS 
    'DEPRECATED: Use customers.primary_phone via customer_id foreign key. Will be dropped after 30 days.';
COMMENT ON COLUMN service_tickets.customer_whatsapp_deprecated IS 
    'DEPRECATED: Use customers.whatsapp_number via customer_id foreign key. Will be dropped after 30 days.';

-- =============================================================================
-- STEP 6: Create view for easy migration of existing queries
-- =============================================================================

CREATE OR REPLACE VIEW service_tickets_with_customer AS
SELECT 
    st.*,
    c.name as customer_name,
    c.primary_phone as customer_phone,
    c.whatsapp_number as customer_whatsapp,
    c.primary_email as customer_email,
    c.customer_type,
    c.city as customer_city,
    c.state as customer_state
FROM service_tickets st
LEFT JOIN customers c ON st.customer_id = c.id::text;

COMMENT ON VIEW service_tickets_with_customer IS 
    'View that joins service_tickets with customers for backward compatibility during migration';

-- =============================================================================
-- STEP 7: Add foreign key constraint (commented out initially for safety)
-- =============================================================================

-- Uncomment after validation in production:
-- ALTER TABLE service_tickets
--     ADD CONSTRAINT fk_ticket_customer 
--     FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT;

-- COMMENT ON CONSTRAINT fk_ticket_customer ON service_tickets IS
--     'Foreign key to customers table ensuring referential integrity';

-- =============================================================================
-- STEP 8: Create trigger to update updated_at timestamp
-- =============================================================================

CREATE OR REPLACE FUNCTION update_customer_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_customer_timestamp
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_customer_timestamp();

-- =============================================================================
-- Migration Complete
-- =============================================================================

DO $$
BEGIN
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'T1.2 Migration Complete!';
    RAISE NOTICE 'Customers table created and data migrated';
    RAISE NOTICE 'Old columns deprecated (not dropped)';
    RAISE NOTICE 'View created for backward compatibility';
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Next steps:';
    RAISE NOTICE '1. Validate customer data';
    RAISE NOTICE '2. Update backend code to use customers table';
    RAISE NOTICE '3. Test thoroughly';
    RAISE NOTICE '4. Enable foreign key constraint';
    RAISE NOTICE '5. After 30 days, drop deprecated columns';
    RAISE NOTICE '==============================================';
END $$;
