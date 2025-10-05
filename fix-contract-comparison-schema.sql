-- ============================================================================
-- Fix Contract and Comparison Service Schemas
-- ============================================================================

-- DROP AND RECREATE CONTRACTS TABLE WITH ALL REQUIRED COLUMNS
DROP TABLE IF EXISTS contract_items CASCADE;
DROP TABLE IF EXISTS contracts CASCADE;

CREATE TABLE contracts (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    contract_number VARCHAR(100) NOT NULL,
    rfq_id VARCHAR(255),
    quote_id VARCHAR(255),
    supplier_id VARCHAR(255) NOT NULL,
    supplier_name VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    tax_amount DECIMAL(15,2),
    start_date DATE,
    end_date DATE,
    signed_date DATE,
    payment_terms TEXT,
    delivery_terms TEXT,
    warranty_terms TEXT,
    terms_and_conditions TEXT,
    payment_schedule JSONB,
    delivery_schedule JSONB,
    items JSONB,
    amendments JSONB,
    notes TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signed_by VARCHAR(255),
    
    CONSTRAINT contracts_tenant_contract_number_idx UNIQUE (tenant_id, contract_number)
);

CREATE INDEX idx_contracts_tenant_id ON contracts(tenant_id);
CREATE INDEX idx_contracts_rfq_id ON contracts(rfq_id);
CREATE INDEX idx_contracts_supplier_id ON contracts(supplier_id);
CREATE INDEX idx_contracts_status ON contracts(status);

-- DROP AND RECREATE COMPARISONS TABLE WITH ALL REQUIRED COLUMNS
DROP TABLE IF EXISTS comparison_items CASCADE;
DROP TABLE IF EXISTS comparisons CASCADE;

CREATE TABLE comparisons (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    rfq_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    quote_ids TEXT[] NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL,
    scoring_criteria JSONB,
    quote_scores JSONB,
    price_differences JSONB,
    item_comparisons JSONB,
    best_overall_quote VARCHAR(255),
    best_price_quote VARCHAR(255),
    recommendation TEXT,
    notes TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_comparisons_tenant_id ON comparisons(tenant_id);
CREATE INDEX idx_comparisons_rfq_id ON comparisons(rfq_id);
CREATE INDEX idx_comparisons_status ON comparisons(status);
CREATE INDEX idx_comparisons_quote_ids ON comparisons USING GIN(quote_ids);

-- Verify tables exist
SELECT 
    'Contracts Table' as table_name,
    COUNT(*) as column_count
FROM information_schema.columns
WHERE table_name = 'contracts'
UNION ALL
SELECT 
    'Comparisons Table',
    COUNT(*)
FROM information_schema.columns
WHERE table_name = 'comparisons';

-- Show contract columns
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'contracts' 
ORDER BY ordinal_position;
