-- ============================================================================
-- Additional Service Tables for ServQR Platform
-- ============================================================================

-- QUOTES SERVICE TABLE
CREATE TABLE IF NOT EXISTS quotes (
    id VARCHAR(255) PRIMARY KEY,
    rfq_id VARCHAR(255) NOT NULL,
    supplier_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    quote_number VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    valid_until TIMESTAMP,
    notes TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP,
    
    CONSTRAINT quotes_tenant_quote_number_idx UNIQUE (tenant_id, quote_number)
);

CREATE INDEX IF NOT EXISTS idx_quotes_tenant_id ON quotes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_quotes_rfq_id ON quotes(rfq_id);
CREATE INDEX IF NOT EXISTS idx_quotes_supplier_id ON quotes(supplier_id);
CREATE INDEX IF NOT EXISTS idx_quotes_status ON quotes(status);

-- QUOTE ITEMS TABLE
CREATE TABLE IF NOT EXISTS quote_items (
    id VARCHAR(255) PRIMARY KEY,
    quote_id VARCHAR(255) NOT NULL,
    rfq_item_id VARCHAR(255),
    description TEXT,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    delivery_time VARCHAR(100),
    specifications JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_quote_items_quote_id ON quote_items(quote_id);

-- COMPARISONS SERVICE TABLE
CREATE TABLE IF NOT EXISTS comparisons (
    id VARCHAR(255) PRIMARY KEY,
    rfq_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_comparisons_tenant_id ON comparisons(tenant_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_rfq_id ON comparisons(rfq_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_status ON comparisons(status);

-- COMPARISON ITEMS TABLE (links quotes being compared)
CREATE TABLE IF NOT EXISTS comparison_items (
    id VARCHAR(255) PRIMARY KEY,
    comparison_id VARCHAR(255) NOT NULL,
    quote_id VARCHAR(255) NOT NULL,
    score DECIMAL(5,2),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comparison_items_comparison_id ON comparison_items(comparison_id);
CREATE INDEX IF NOT EXISTS idx_comparison_items_quote_id ON comparison_items(quote_id);

-- CONTRACTS SERVICE TABLE
CREATE TABLE IF NOT EXISTS contracts (
    id VARCHAR(255) PRIMARY KEY,
    contract_number VARCHAR(100) NOT NULL,
    rfq_id VARCHAR(255),
    quote_id VARCHAR(255),
    supplier_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    contract_value DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    start_date DATE,
    end_date DATE,
    payment_terms JSONB,
    delivery_terms JSONB,
    terms_and_conditions TEXT,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    signed_at TIMESTAMP,
    
    CONSTRAINT contracts_tenant_contract_number_idx UNIQUE (tenant_id, contract_number)
);

CREATE INDEX IF NOT EXISTS idx_contracts_tenant_id ON contracts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contracts_rfq_id ON contracts(rfq_id);
CREATE INDEX IF NOT EXISTS idx_contracts_supplier_id ON contracts(supplier_id);
CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status);

-- CONTRACT ITEMS TABLE
CREATE TABLE IF NOT EXISTS contract_items (
    id VARCHAR(255) PRIMARY KEY,
    contract_id VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50) NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    specifications JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_contract_items_contract_id ON contract_items(contract_id);

-- Display summary
SELECT 
    'Quotes' as table_name, 
    COUNT(*) as record_count 
FROM quotes
UNION ALL
SELECT 'Quote Items', COUNT(*) FROM quote_items
UNION ALL
SELECT 'Comparisons', COUNT(*) FROM comparisons
UNION ALL
SELECT 'Comparison Items', COUNT(*) FROM comparison_items
UNION ALL
SELECT 'Contracts', COUNT(*) FROM contracts
UNION ALL
SELECT 'Contract Items', COUNT(*) FROM contract_items;
