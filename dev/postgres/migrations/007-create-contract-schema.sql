-- Migration: Create contract schema for contract management and tracking
-- This creates tables for managing contracts with payment/delivery schedules

-- Contracts table
CREATE TABLE IF NOT EXISTS contracts (
    id VARCHAR(32) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    contract_number VARCHAR(50) UNIQUE NOT NULL,
    rfq_id VARCHAR(32) NOT NULL,
    quote_id VARCHAR(32) NOT NULL,
    supplier_id VARCHAR(32) NOT NULL,
    supplier_name VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    
    -- Financial
    total_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    tax_amount DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    
    -- Dates
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    signed_date TIMESTAMP WITH TIME ZONE,
    
    -- Terms (text fields)
    payment_terms TEXT,
    delivery_terms TEXT,
    warranty_terms TEXT,
    terms_and_conditions TEXT,
    
    -- Schedules and items (stored as JSONB for flexibility)
    payment_schedule JSONB NOT NULL DEFAULT '[]'::jsonb,
    delivery_schedule JSONB NOT NULL DEFAULT '[]'::jsonb,
    items JSONB NOT NULL DEFAULT '[]'::jsonb,
    amendments JSONB NOT NULL DEFAULT '[]'::jsonb,
    
    -- Metadata
    notes TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    signed_by VARCHAR(255),
    
    CONSTRAINT contracts_status_check CHECK (status IN ('draft', 'active', 'completed', 'cancelled', 'expired', 'suspended'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_contracts_tenant ON contracts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contracts_contract_number ON contracts(contract_number);
CREATE INDEX IF NOT EXISTS idx_contracts_rfq ON contracts(rfq_id);
CREATE INDEX IF NOT EXISTS idx_contracts_quote ON contracts(quote_id);
CREATE INDEX IF NOT EXISTS idx_contracts_supplier ON contracts(supplier_id);
CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status);
CREATE INDEX IF NOT EXISTS idx_contracts_start_date ON contracts(start_date);
CREATE INDEX IF NOT EXISTS idx_contracts_end_date ON contracts(end_date);
CREATE INDEX IF NOT EXISTS idx_contracts_created_at ON contracts(created_at DESC);

-- GIN indexes for JSONB queries
CREATE INDEX IF NOT EXISTS idx_contracts_items ON contracts USING GIN (items);
CREATE INDEX IF NOT EXISTS idx_contracts_payment_schedule ON contracts USING GIN (payment_schedule);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_contracts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER contracts_updated_at_trigger
    BEFORE UPDATE ON contracts
    FOR EACH ROW
    EXECUTE FUNCTION update_contracts_updated_at();

-- Row Level Security (RLS) for multi-tenancy
ALTER TABLE contracts ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see contracts from their tenant
CREATE POLICY contracts_tenant_isolation ON contracts
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant', TRUE));

-- Helper function to generate contract numbers
CREATE OR REPLACE FUNCTION generate_contract_number(p_tenant_id VARCHAR)
RETURNS VARCHAR AS $$
DECLARE
    v_count INT;
    v_number VARCHAR;
BEGIN
    -- Count existing contracts for this tenant
    SELECT COUNT(*) INTO v_count
    FROM contracts
    WHERE tenant_id = p_tenant_id;
    
    -- Generate number: CT-YYYYMMDD-XXXX
    v_number := 'CT-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD((v_count + 1)::TEXT, 4, '0');
    
    RETURN v_number;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get contract statistics
CREATE OR REPLACE FUNCTION get_contract_statistics(p_tenant_id VARCHAR)
RETURNS TABLE (
    total_contracts BIGINT,
    active_contracts BIGINT,
    completed_contracts BIGINT,
    total_contract_value NUMERIC,
    avg_contract_value NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_contracts,
        COUNT(*) FILTER (WHERE status = 'active')::BIGINT as active_contracts,
        COUNT(*) FILTER (WHERE status = 'completed')::BIGINT as completed_contracts,
        COALESCE(SUM(total_amount), 0)::NUMERIC as total_contract_value,
        COALESCE(AVG(total_amount), 0)::NUMERIC as avg_contract_value
    FROM contracts
    WHERE tenant_id = p_tenant_id;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get contracts expiring soon
CREATE OR REPLACE FUNCTION get_expiring_contracts(p_tenant_id VARCHAR, p_days INT DEFAULT 30)
RETURNS TABLE (
    contract_id VARCHAR,
    contract_number VARCHAR,
    supplier_name VARCHAR,
    end_date TIMESTAMP WITH TIME ZONE,
    days_until_expiry INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.contract_number,
        c.supplier_name,
        c.end_date,
        EXTRACT(DAY FROM (c.end_date - NOW()))::INT as days_until_expiry
    FROM contracts c
    WHERE c.tenant_id = p_tenant_id
      AND c.status = 'active'
      AND c.end_date BETWEEN NOW() AND (NOW() + INTERVAL '1 day' * p_days)
    ORDER BY c.end_date ASC;
END;
$$ LANGUAGE plpgsql;

-- Helper function to check payment progress
CREATE OR REPLACE FUNCTION get_payment_progress(p_contract_id VARCHAR)
RETURNS TABLE (
    total_payments INT,
    paid_payments INT,
    payment_progress NUMERIC
) AS $$
DECLARE
    v_payment_schedule JSONB;
    v_total INT;
    v_paid INT;
BEGIN
    -- Get payment schedule
    SELECT payment_schedule INTO v_payment_schedule
    FROM contracts
    WHERE id = p_contract_id;
    
    -- Count total and paid payments
    v_total := jsonb_array_length(v_payment_schedule);
    
    SELECT COUNT(*) INTO v_paid
    FROM jsonb_array_elements(v_payment_schedule) AS payment
    WHERE (payment->>'paid')::BOOLEAN = TRUE;
    
    -- Calculate progress percentage
    RETURN QUERY
    SELECT
        v_total,
        v_paid,
        CASE
            WHEN v_total > 0 THEN (v_paid::NUMERIC / v_total::NUMERIC * 100)
            ELSE 0
        END;
END;
$$ LANGUAGE plpgsql;

-- Seed data for testing (optional)
-- Commented out for production, uncomment for development testing
/*
INSERT INTO contracts (
    id, tenant_id, contract_number, rfq_id, quote_id,
    supplier_id, supplier_name, status, total_amount,
    start_date, end_date, created_by
) VALUES (
    '33QContractTestID1234567',
    'city-hospital',
    'CT-20250930-0001',
    'test-rfq-001',
    '33QPQkOOjzSIZdRjArSvjwGtFWr',
    '33Q6Ziuic9GoTjlBFjyUQT03N79',
    'Test Medical Supplier Inc.',
    'draft',
    330000.00,
    NOW(),
    NOW() + INTERVAL '1 year',
    'test-user'
);
*/

-- Grant permissions (adjust based on your user roles)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON contracts TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;
