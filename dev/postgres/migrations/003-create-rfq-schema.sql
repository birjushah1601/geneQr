-- =============================================
-- RFQ Service - Database Schema
-- =============================================
-- This migration creates the schema for the RFQ (Request for Quote) service

-- ======================
-- 1. RFQs TABLE
-- ======================

CREATE TABLE IF NOT EXISTS rfqs (
    -- Identity
    id VARCHAR(26) PRIMARY KEY,
    rfq_number VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(50) NOT NULL,
    
    -- Basic Information
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'published', 'closed', 'awarded', 'cancelled')),
    
    -- Delivery Terms (stored as JSONB for flexibility)
    delivery_terms JSONB NOT NULL,
    
    -- Payment Terms (stored as JSONB for flexibility)
    payment_terms JSONB NOT NULL,
    
    -- Timeline
    published_at TIMESTAMPTZ,
    response_deadline TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ,
    
    -- Metadata
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Notes
    internal_notes TEXT,
    
    -- Constraints
    UNIQUE(rfq_number, tenant_id)
);

-- Indexes for RFQs table
CREATE INDEX IF NOT EXISTS idx_rfqs_tenant_id ON rfqs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_status ON rfqs(status);
CREATE INDEX IF NOT EXISTS idx_rfqs_priority ON rfqs(priority);
CREATE INDEX IF NOT EXISTS idx_rfqs_created_by ON rfqs(created_by);
CREATE INDEX IF NOT EXISTS idx_rfqs_created_at ON rfqs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rfqs_response_deadline ON rfqs(response_deadline);
CREATE INDEX IF NOT EXISTS idx_rfqs_tenant_status ON rfqs(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_rfqs_rfq_number ON rfqs(rfq_number);

-- Full-text search index on title and description
CREATE INDEX IF NOT EXISTS idx_rfqs_search ON rfqs USING gin(
    to_tsvector('english', title || ' ' || description)
);

-- ======================
-- 2. RFQ ITEMS TABLE
-- ======================

CREATE TABLE IF NOT EXISTS rfq_items (
    -- Identity
    id VARCHAR(26) PRIMARY KEY,
    rfq_id VARCHAR(26) NOT NULL REFERENCES rfqs(id) ON DELETE CASCADE,
    
    -- Equipment Reference (optional - can be custom item)
    equipment_id VARCHAR(26),
    category_id VARCHAR(26),
    
    -- Item Details
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    specifications JSONB,
    
    -- Quantity
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL DEFAULT 'piece',
    
    -- Pricing
    estimated_price DECIMAL(12,2),
    
    -- Notes
    notes TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for RFQ Items table
CREATE INDEX IF NOT EXISTS idx_rfq_items_rfq_id ON rfq_items(rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_items_equipment_id ON rfq_items(equipment_id);
CREATE INDEX IF NOT EXISTS idx_rfq_items_category_id ON rfq_items(category_id);

-- ======================
-- 3. RFQ INVITATIONS TABLE
-- ======================

CREATE TABLE IF NOT EXISTS rfq_invitations (
    -- Identity
    id VARCHAR(26) PRIMARY KEY,
    rfq_id VARCHAR(26) NOT NULL REFERENCES rfqs(id) ON DELETE CASCADE,
    supplier_id VARCHAR(255) NOT NULL,
    
    -- Status
    status VARCHAR(20) NOT NULL CHECK (status IN ('invited', 'viewed', 'quoted', 'declined')),
    
    -- Timeline
    invited_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    viewed_at TIMESTAMPTZ,
    responded_at TIMESTAMPTZ,
    
    -- Message
    message TEXT,
    
    -- Constraints
    UNIQUE(rfq_id, supplier_id)
);

-- Indexes for RFQ Invitations table
CREATE INDEX IF NOT EXISTS idx_rfq_invitations_rfq_id ON rfq_invitations(rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_invitations_supplier_id ON rfq_invitations(supplier_id);
CREATE INDEX IF NOT EXISTS idx_rfq_invitations_status ON rfq_invitations(status);

-- ======================
-- 4. TRIGGERS
-- ======================

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_rfq_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for rfqs table
DROP TRIGGER IF EXISTS trigger_update_rfqs_timestamp ON rfqs;
CREATE TRIGGER trigger_update_rfqs_timestamp
    BEFORE UPDATE ON rfqs
    FOR EACH ROW
    EXECUTE FUNCTION update_rfq_timestamp();

-- Trigger for rfq_items table
DROP TRIGGER IF EXISTS trigger_update_rfq_items_timestamp ON rfq_items;
CREATE TRIGGER trigger_update_rfq_items_timestamp
    BEFORE UPDATE ON rfq_items
    FOR EACH ROW
    EXECUTE FUNCTION update_rfq_timestamp();

-- ======================
-- 5. ROW LEVEL SECURITY (Optional - for multi-tenancy)
-- ======================

-- Enable RLS on rfqs table
ALTER TABLE rfqs ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only access RFQs from their tenant
CREATE POLICY rfq_tenant_isolation_policy ON rfqs
    USING (tenant_id = current_setting('app.current_tenant', true))
    WITH CHECK (tenant_id = current_setting('app.current_tenant', true));

-- Policy: Superuser bypass
CREATE POLICY rfq_superuser_bypass ON rfqs
    USING (current_user = 'postgres');

-- ======================
-- 6. HELPER FUNCTIONS
-- ======================

-- Function to generate RFQ number
CREATE OR REPLACE FUNCTION generate_rfq_number(p_tenant_id VARCHAR)
RETURNS VARCHAR AS $$
DECLARE
    v_count INTEGER;
    v_year VARCHAR(4);
    v_rfq_number VARCHAR(50);
BEGIN
    -- Get current year
    v_year := TO_CHAR(CURRENT_DATE, 'YYYY');
    
    -- Get count of RFQs for this tenant this year
    SELECT COUNT(*) INTO v_count
    FROM rfqs
    WHERE tenant_id = p_tenant_id
    AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE);
    
    -- Generate RFQ number: RFQ-YYYY-NNNN
    v_rfq_number := 'RFQ-' || v_year || '-' || LPAD((v_count + 1)::TEXT, 4, '0');
    
    RETURN v_rfq_number;
END;
$$ LANGUAGE plpgsql;

-- Function to get RFQ statistics
CREATE OR REPLACE FUNCTION get_rfq_stats(p_tenant_id VARCHAR)
RETURNS TABLE (
    total_rfqs BIGINT,
    draft_rfqs BIGINT,
    published_rfqs BIGINT,
    closed_rfqs BIGINT,
    awarded_rfqs BIGINT,
    cancelled_rfqs BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT AS total_rfqs,
        COUNT(*) FILTER (WHERE status = 'draft')::BIGINT AS draft_rfqs,
        COUNT(*) FILTER (WHERE status = 'published')::BIGINT AS published_rfqs,
        COUNT(*) FILTER (WHERE status = 'closed')::BIGINT AS closed_rfqs,
        COUNT(*) FILTER (WHERE status = 'awarded')::BIGINT AS awarded_rfqs,
        COUNT(*) FILTER (WHERE status = 'cancelled')::BIGINT AS cancelled_rfqs
    FROM rfqs
    WHERE tenant_id = p_tenant_id;
END;
$$ LANGUAGE plpgsql;

-- ======================
-- 7. COMMENTS
-- ======================

COMMENT ON TABLE rfqs IS 'Request for Quote (RFQ) records for procurement workflow';
COMMENT ON TABLE rfq_items IS 'Items/equipment requested in each RFQ';
COMMENT ON TABLE rfq_invitations IS 'Supplier invitations for RFQs';

COMMENT ON COLUMN rfqs.rfq_number IS 'Human-readable unique RFQ number (e.g., RFQ-2025-0001)';
COMMENT ON COLUMN rfqs.priority IS 'Urgency level: low, medium, high, critical';
COMMENT ON COLUMN rfqs.status IS 'Lifecycle status: draft, published, closed, awarded, cancelled';
COMMENT ON COLUMN rfqs.delivery_terms IS 'Delivery requirements stored as JSONB';
COMMENT ON COLUMN rfqs.payment_terms IS 'Payment terms stored as JSONB';
COMMENT ON COLUMN rfqs.response_deadline IS 'Deadline for suppliers to submit quotes';

-- ======================
-- 8. SAMPLE DATA (for development)
-- ======================

-- This will be populated by seed data scripts

-- Migration complete
SELECT 'RFQ schema migration completed successfully' AS status;
