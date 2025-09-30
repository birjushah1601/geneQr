-- Supplier Management Schema
-- Version: 1.0
-- Description: Creates tables for supplier profiles, certifications, and performance tracking

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Suppliers table
CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(32) PRIMARY KEY,
    tenant_id VARCHAR(50) NOT NULL,
    
    -- Company Information
    company_name VARCHAR(255) NOT NULL,
    business_registration_number VARCHAR(100),
    tax_id VARCHAR(100),
    year_established INTEGER,
    description TEXT,
    
    -- Contact Information (JSONB for flexibility)
    contact_info JSONB NOT NULL,
    
    -- Address (JSONB)
    address JSONB NOT NULL,
    
    -- Specializations (array of category IDs)
    specializations TEXT[] DEFAULT '{}',
    
    -- Certifications (JSONB array)
    certifications JSONB DEFAULT '[]',
    
    -- Performance Metrics
    performance_rating DECIMAL(3,2) DEFAULT 0.00 CHECK (performance_rating >= 0.00 AND performance_rating <= 5.00),
    total_orders INTEGER DEFAULT 0,
    completed_orders INTEGER DEFAULT 0,
    
    -- Status
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'active', 'suspended', 'inactive')),
    verification_status VARCHAR(20) NOT NULL CHECK (verification_status IN ('pending', 'approved', 'rejected')),
    verified_at TIMESTAMP WITH TIME ZONE,
    verified_by VARCHAR(255),
    
    -- Metadata
    metadata JSONB DEFAULT '{}',
    
    -- Audit fields
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    UNIQUE(tax_id, tenant_id)
);

-- Indexes for suppliers
CREATE INDEX idx_suppliers_tenant_id ON suppliers(tenant_id);
CREATE INDEX idx_suppliers_company_name ON suppliers(company_name);
CREATE INDEX idx_suppliers_status ON suppliers(status);
CREATE INDEX idx_suppliers_verification_status ON suppliers(verification_status);
CREATE INDEX idx_suppliers_performance_rating ON suppliers(performance_rating DESC);
CREATE INDEX idx_suppliers_specializations ON suppliers USING GIN(specializations);
CREATE INDEX idx_suppliers_created_at ON suppliers(created_at DESC);
CREATE INDEX idx_suppliers_search ON suppliers USING GIN(to_tsvector('english', company_name || ' ' || COALESCE(description, '')));
CREATE INDEX idx_suppliers_tenant_status ON suppliers(tenant_id, status);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_supplier_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for updated_at
DROP TRIGGER IF EXISTS trigger_update_suppliers_timestamp ON suppliers;
CREATE TRIGGER trigger_update_suppliers_timestamp
    BEFORE UPDATE ON suppliers
    FOR EACH ROW
    EXECUTE FUNCTION update_supplier_timestamp();

-- Enable Row Level Security
ALTER TABLE suppliers ENABLE ROW LEVEL SECURITY;

-- RLS Policy for tenant isolation
CREATE POLICY supplier_tenant_isolation_policy ON suppliers
    USING (tenant_id = current_setting('app.current_tenant', true))
    WITH CHECK (tenant_id = current_setting('app.current_tenant', true));

-- Policy for superuser bypass
CREATE POLICY supplier_superuser_bypass ON suppliers
    USING (current_user = 'postgres');

-- Helper function to get supplier statistics
CREATE OR REPLACE FUNCTION get_supplier_stats(p_supplier_id VARCHAR(32))
RETURNS JSON AS $$
DECLARE
    result JSON;
BEGIN
    SELECT json_build_object(
        'total_orders', total_orders,
        'completed_orders', completed_orders,
        'completion_rate', 
            CASE 
                WHEN total_orders > 0 THEN ROUND((completed_orders::DECIMAL / total_orders) * 100, 2)
                ELSE 0
            END,
        'performance_rating', performance_rating,
        'certifications_count', jsonb_array_length(COALESCE(certifications, '[]'::jsonb)),
        'specializations_count', array_length(specializations, 1)
    ) INTO result
    FROM suppliers
    WHERE id = p_supplier_id;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Helper function to find suppliers by category
CREATE OR REPLACE FUNCTION find_suppliers_by_category(
    p_category_id VARCHAR(50),
    p_tenant_id VARCHAR(50)
)
RETURNS SETOF suppliers AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM suppliers
    WHERE tenant_id = p_tenant_id
      AND status = 'active'
      AND verification_status = 'approved'
      AND p_category_id = ANY(specializations)
    ORDER BY performance_rating DESC;
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE suppliers IS 'Stores supplier/vendor profiles and information';
COMMENT ON COLUMN suppliers.contact_info IS 'JSON object containing primary and secondary contact details';
COMMENT ON COLUMN suppliers.address IS 'JSON object containing full address information';
COMMENT ON COLUMN suppliers.specializations IS 'Array of category IDs the supplier specializes in';
COMMENT ON COLUMN suppliers.certifications IS 'JSON array of certifications (ISO, FDA, etc.)';
COMMENT ON COLUMN suppliers.performance_rating IS 'Supplier performance rating from 0.00 to 5.00';
COMMENT ON COLUMN suppliers.metadata IS 'Additional flexible data storage';

-- Verification completed successfully
SELECT 'Supplier schema migration completed successfully' AS status;
