-- Quote Service Schema
-- Migration 005: Create tables for quote management

-- Enable necessary extensions (if not already enabled)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search

-- ============================================================================
-- 1. QUOTES TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS quotes (
    -- Identity
    id VARCHAR(32) PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    rfq_id VARCHAR(32) NOT NULL,
    supplier_id VARCHAR(32) NOT NULL,

    -- Quote details
    quote_number VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    valid_until TIMESTAMP NOT NULL,
    delivery_terms TEXT,
    payment_terms TEXT,
    warranty_terms TEXT,
    notes TEXT,

    -- Revision tracking
    revision_number INT NOT NULL DEFAULT 1,

    -- Review information
    reviewed_at TIMESTAMP,
    reviewed_by VARCHAR(100),
    review_notes TEXT,
    rejection_reason TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for quotes
CREATE INDEX idx_quotes_tenant ON quotes(tenant_id);
CREATE INDEX idx_quotes_rfq ON quotes(rfq_id, tenant_id);
CREATE INDEX idx_quotes_supplier ON quotes(supplier_id, tenant_id);
CREATE INDEX idx_quotes_status ON quotes(status, tenant_id);
CREATE INDEX idx_quotes_created ON quotes(created_at DESC);
CREATE INDEX idx_quotes_amount ON quotes(total_amount);
CREATE INDEX idx_quotes_valid_until ON quotes(valid_until);
CREATE INDEX idx_quotes_quote_number ON quotes(quote_number) WHERE quote_number IS NOT NULL;

-- GIN index for metadata JSONB
CREATE INDEX idx_quotes_metadata ON quotes USING gin(metadata);

-- ============================================================================
-- 2. QUOTE_ITEMS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS quote_items (
    -- Identity
    id VARCHAR(32) PRIMARY KEY,
    quote_id VARCHAR(32) NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    rfq_item_id VARCHAR(32) NOT NULL,
    
    -- Item details
    equipment_id VARCHAR(32) NOT NULL,
    equipment_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(15, 2) NOT NULL,
    total_price DECIMAL(15, 2) NOT NULL,
    tax_rate DECIMAL(5, 4) DEFAULT 0,
    tax_amount DECIMAL(15, 2) DEFAULT 0,
    delivery_timeframe VARCHAR(100),
    
    -- Product information
    manufacturer_name VARCHAR(255),
    model_number VARCHAR(100),
    specifications TEXT,
    compliance_certs TEXT,
    notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for quote_items
CREATE INDEX idx_quote_items_quote ON quote_items(quote_id);
CREATE INDEX idx_quote_items_rfq_item ON quote_items(rfq_item_id);
CREATE INDEX idx_quote_items_equipment ON quote_items(equipment_id);

-- ============================================================================
-- 3. QUOTE_REVISIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS quote_revisions (
    -- Identity
    id SERIAL PRIMARY KEY,
    quote_id VARCHAR(32) NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    
    -- Revision details
    revision_number INT NOT NULL,
    revised_at TIMESTAMP NOT NULL,
    revised_by VARCHAR(100) NOT NULL,
    changes TEXT,
    previous_total DECIMAL(15, 2) NOT NULL,
    new_total DECIMAL(15, 2) NOT NULL,
    
    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,
    
    UNIQUE(quote_id, revision_number)
);

-- Indexes for quote_revisions
CREATE INDEX idx_quote_revisions_quote ON quote_revisions(quote_id, revision_number DESC);
CREATE INDEX idx_quote_revisions_revised_at ON quote_revisions(revised_at DESC);

-- ============================================================================
-- 4. HELPER FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_quote_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for auto-updating updated_at
CREATE TRIGGER trigger_update_quotes_updated_at
    BEFORE UPDATE ON quotes
    FOR EACH ROW
    EXECUTE FUNCTION update_quote_updated_at();

CREATE TRIGGER trigger_update_quote_items_updated_at
    BEFORE UPDATE ON quote_items
    FOR EACH ROW
    EXECUTE FUNCTION update_quote_updated_at();

-- Function to validate quote status transitions
CREATE OR REPLACE FUNCTION validate_quote_status_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Allow any transition for new records
    IF TG_OP = 'INSERT' THEN
        RETURN NEW;
    END IF;

    -- Validate status transitions
    IF OLD.status = 'accepted' OR OLD.status = 'rejected' THEN
        IF NEW.status != OLD.status AND NEW.status != 'expired' THEN
            RAISE EXCEPTION 'Cannot change status from % to %', OLD.status, NEW.status;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for quote status validation
CREATE TRIGGER trigger_validate_quote_status
    BEFORE UPDATE ON quotes
    FOR EACH ROW
    EXECUTE FUNCTION validate_quote_status_transition();

-- ============================================================================
-- 5. ROW LEVEL SECURITY (RLS)
-- ============================================================================

-- Enable RLS on quotes table
ALTER TABLE quotes ENABLE ROW LEVEL SECURITY;

-- Policy for tenant isolation
CREATE POLICY quote_tenant_isolation_policy ON quotes
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant', true))
    WITH CHECK (tenant_id = current_setting('app.current_tenant', true));

-- Policy for superuser bypass
CREATE POLICY quote_superuser_bypass ON quotes
    FOR ALL
    TO PUBLIC
    USING (current_user = 'postgres');

-- Note: For production, create specific application users and restrict postgres access

-- ============================================================================
-- 6. COMMENTS
-- ============================================================================

COMMENT ON TABLE quotes IS 'Stores vendor quotes in response to RFQs';
COMMENT ON TABLE quote_items IS 'Line items for each quote';
COMMENT ON TABLE quote_revisions IS 'Revision history for quotes';

COMMENT ON COLUMN quotes.status IS 'Quote lifecycle status: draft, submitted, under_review, revised, accepted, rejected, expired, withdrawn';
COMMENT ON COLUMN quotes.revision_number IS 'Current revision number, incremented with each revision';
COMMENT ON COLUMN quotes.valid_until IS 'Expiration date/time for this quote';
COMMENT ON COLUMN quote_items.tax_rate IS 'Tax rate as decimal (e.g., 0.10 for 10%)';
COMMENT ON COLUMN quote_revisions.changes IS 'Human-readable description of what changed in this revision';

-- ============================================================================
-- 7. INITIAL DATA / SEED
-- ============================================================================

-- No initial seed data for quotes - created by suppliers dynamically

-- Migration complete
SELECT 'Quote schema migration completed successfully' AS status;
