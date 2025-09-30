-- Migration: Create comparison schema for quote comparison and scoring
-- This creates tables for managing quote comparisons with scoring analysis

-- Comparisons table
CREATE TABLE IF NOT EXISTS comparisons (
    id VARCHAR(32) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    rfq_id VARCHAR(32) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    quote_ids TEXT[] NOT NULL, -- Array of quote IDs
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    
    -- Scoring criteria (stored as JSONB for flexibility)
    scoring_criteria JSONB NOT NULL DEFAULT '{"price_weight": 40, "quality_weight": 30, "delivery_weight": 20, "compliance_weight": 10}'::jsonb,
    
    -- Analysis results (stored as JSONB)
    quote_scores JSONB NOT NULL DEFAULT '[]'::jsonb,
    price_differences JSONB NOT NULL DEFAULT '[]'::jsonb,
    item_comparisons JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    -- Best recommendations
    best_overall_quote VARCHAR(32),
    best_price_quote VARCHAR(32),
    recommendation TEXT,
    
    notes TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT comparisons_status_check CHECK (status IN ('draft', 'active', 'completed', 'archived'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_comparisons_tenant ON comparisons(tenant_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_rfq ON comparisons(rfq_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_status ON comparisons(status);
CREATE INDEX IF NOT EXISTS idx_comparisons_created_by ON comparisons(created_by);
CREATE INDEX IF NOT EXISTS idx_comparisons_created_at ON comparisons(created_at DESC);

-- GIN index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_comparisons_quote_scores ON comparisons USING GIN (quote_scores);
CREATE INDEX IF NOT EXISTS idx_comparisons_quote_ids ON comparisons USING GIN (quote_ids);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_comparisons_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER comparisons_updated_at_trigger
    BEFORE UPDATE ON comparisons
    FOR EACH ROW
    EXECUTE FUNCTION update_comparisons_updated_at();

-- Row Level Security (RLS) for multi-tenancy
ALTER TABLE comparisons ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see comparisons from their tenant
CREATE POLICY comparisons_tenant_isolation ON comparisons
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant', TRUE));

-- Helper function to calculate comparison statistics
CREATE OR REPLACE FUNCTION get_comparison_statistics(p_tenant_id VARCHAR)
RETURNS TABLE (
    total_comparisons BIGINT,
    active_comparisons BIGINT,
    completed_comparisons BIGINT,
    avg_quotes_per_comparison NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_comparisons,
        COUNT(*) FILTER (WHERE status = 'active')::BIGINT as active_comparisons,
        COUNT(*) FILTER (WHERE status = 'completed')::BIGINT as completed_comparisons,
        AVG(array_length(quote_ids, 1))::NUMERIC as avg_quotes_per_comparison
    FROM comparisons
    WHERE tenant_id = p_tenant_id;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get comparison summary by RFQ
CREATE OR REPLACE FUNCTION get_rfq_comparison_summary(p_tenant_id VARCHAR, p_rfq_id VARCHAR)
RETURNS TABLE (
    comparison_id VARCHAR,
    title VARCHAR,
    status VARCHAR,
    quote_count INT,
    best_overall_quote VARCHAR,
    best_price_quote VARCHAR,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.title,
        c.status,
        array_length(c.quote_ids, 1) as quote_count,
        c.best_overall_quote,
        c.best_price_quote,
        c.created_at
    FROM comparisons c
    WHERE c.tenant_id = p_tenant_id
      AND c.rfq_id = p_rfq_id
    ORDER BY c.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Seed data for testing (optional)
-- This will create a sample comparison once quotes exist
-- Commented out for production, uncomment for development testing
/*
INSERT INTO comparisons (id, tenant_id, rfq_id, title, quote_ids, status, created_by)
VALUES (
    '33QComparisonTestID1234567',
    'city-hospital',
    'test-rfq-001',
    'X-Ray Machine Comparison',
    ARRAY['quote-id-1', 'quote-id-2'],
    'draft',
    'test-user'
);
*/

-- Grant permissions (adjust based on your user roles)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON comparisons TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
