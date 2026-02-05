-- ============================================================================
-- ServQR Platform Database Schema Fix
-- This script creates all missing tables and fixes schema mismatches
-- ============================================================================

-- Drop existing tables if they exist (for clean setup)
DROP TABLE IF EXISTS rfq_invitations CASCADE;
DROP TABLE IF EXISTS rfq_items CASCADE;
DROP TABLE IF EXISTS rfqs CASCADE;

-- ============================================================================
-- RFQ SERVICE TABLES
-- ============================================================================

-- Main RFQ table
CREATE TABLE IF NOT EXISTS rfqs (
    id VARCHAR(255) PRIMARY KEY,
    rfq_number VARCHAR(100) NOT NULL UNIQUE,
    tenant_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority VARCHAR(50) NOT NULL,  -- low, medium, high, critical
    status VARCHAR(50) NOT NULL,     -- draft, published, closed, awarded, cancelled
    delivery_terms JSONB NOT NULL,
    payment_terms JSONB NOT NULL,
    published_at TIMESTAMP,
    response_deadline TIMESTAMP NOT NULL,
    closed_at TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    internal_notes TEXT,
    
    CONSTRAINT rfqs_tenant_id_idx UNIQUE (tenant_id, rfq_number)
);

-- RFQ Items table
CREATE TABLE IF NOT EXISTS rfq_items (
    id VARCHAR(255) PRIMARY KEY,
    rfq_id VARCHAR(255) NOT NULL REFERENCES rfqs(id) ON DELETE CASCADE,
    equipment_id VARCHAR(255),
    category_id VARCHAR(255),
    name VARCHAR(500) NOT NULL,
    description TEXT,
    specifications JSONB,
    quantity INTEGER NOT NULL,
    unit VARCHAR(50) NOT NULL,
    estimated_price DECIMAL(15,2),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- RFQ Invitations table
CREATE TABLE IF NOT EXISTS rfq_invitations (
    id VARCHAR(255) PRIMARY KEY,
    rfq_id VARCHAR(255) NOT NULL REFERENCES rfqs(id) ON DELETE CASCADE,
    supplier_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,  -- invited, viewed, quoted, declined
    invited_at TIMESTAMP NOT NULL,
    viewed_at TIMESTAMP,
    responded_at TIMESTAMP,
    message TEXT
);

-- Create indexes for RFQ tables
CREATE INDEX IF NOT EXISTS idx_rfqs_tenant_id ON rfqs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_status ON rfqs(status);
CREATE INDEX IF NOT EXISTS idx_rfqs_created_at ON rfqs(created_at);
CREATE INDEX IF NOT EXISTS idx_rfq_items_rfq_id ON rfq_items(rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_invitations_rfq_id ON rfq_invitations(rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_invitations_supplier_id ON rfq_invitations(supplier_id);

-- ============================================================================
-- SUPPLIER SERVICE TABLE FIX
-- ============================================================================

-- Drop and recreate suppliers table with correct schema
DROP TABLE IF EXISTS suppliers CASCADE;

CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    company_name VARCHAR(500) NOT NULL,
    business_registration_number VARCHAR(100),
    tax_id VARCHAR(100),
    year_established INTEGER,
    description TEXT,
    contact_info JSONB NOT NULL,  -- Contains contact details
    address JSONB NOT NULL,        -- Contains address details
    specializations TEXT[] NOT NULL DEFAULT '{}',  -- Array of category IDs
    certifications JSONB NOT NULL DEFAULT '[]',    -- Array of certification objects
    performance_rating DECIMAL(3,2) DEFAULT 0.0,
    total_orders INTEGER DEFAULT 0,
    completed_orders INTEGER DEFAULT 0,
    status VARCHAR(50) NOT NULL,   -- pending, active, suspended, inactive
    verification_status VARCHAR(50) NOT NULL,  -- pending, approved, rejected
    verified_at TIMESTAMP,
    verified_by VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT suppliers_tenant_tax_id_idx UNIQUE (tenant_id, tax_id)
);

-- Create indexes for suppliers
CREATE INDEX IF NOT EXISTS idx_suppliers_tenant_id ON suppliers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_status ON suppliers(status);
CREATE INDEX IF NOT EXISTS idx_suppliers_verification_status ON suppliers(verification_status);
CREATE INDEX IF NOT EXISTS idx_suppliers_specializations ON suppliers USING GIN(specializations);

-- ============================================================================
-- CATALOG SERVICE TABLES (ensure they exist with correct schema)
-- ============================================================================

-- Equipment table (used by catalog service)
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    category_id VARCHAR(255) NOT NULL,
    manufacturer_id VARCHAR(255) NOT NULL,
    model VARCHAR(255),
    description TEXT,
    specifications JSONB,
    price_amount DECIMAL(15,2) NOT NULL,
    price_currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    sku VARCHAR(100),
    images TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tenant_id VARCHAR(255) NOT NULL
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(255) REFERENCES categories(id),
    tenant_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Manufacturers table
CREATE TABLE IF NOT EXISTS manufacturers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100),
    website VARCHAR(500),
    tenant_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for equipment/catalog
CREATE INDEX IF NOT EXISTS idx_equipment_tenant_id ON equipment(tenant_id);
CREATE INDEX IF NOT EXISTS idx_equipment_category_id ON equipment(category_id);
CREATE INDEX IF NOT EXISTS idx_equipment_manufacturer_id ON equipment(manufacturer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_is_active ON equipment(is_active);
CREATE INDEX IF NOT EXISTS idx_categories_tenant_id ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_manufacturers_tenant_id ON manufacturers(tenant_id);

-- ============================================================================
-- SAMPLE DATA
-- ============================================================================

-- Insert sample categories
INSERT INTO categories (id, name, parent_id, tenant_id) VALUES
    ('cat-001', 'Diagnostic Equipment', NULL, 'city-hospital'),
    ('cat-002', 'Surgical Equipment', NULL, 'city-hospital'),
    ('cat-003', 'Laboratory Equipment', NULL, 'city-hospital'),
    ('cat-004', 'Imaging Equipment', 'cat-001', 'city-hospital')
ON CONFLICT (id) DO NOTHING;

-- Insert sample manufacturers
INSERT INTO manufacturers (id, name, country, website, tenant_id) VALUES
    ('mfr-001', 'MedTech Solutions', 'USA', 'https://medtechsolutions.com', 'city-hospital'),
    ('mfr-002', 'Global Medical Devices', 'Germany', 'https://globalmedical.de', 'city-hospital'),
    ('mfr-003', 'HealthCare Innovations', 'Japan', 'https://healthcareinnovations.jp', 'city-hospital')
ON CONFLICT (id) DO NOTHING;

-- Insert sample equipment (for catalog service)
INSERT INTO equipment (id, name, category_id, manufacturer_id, model, description, specifications, price_amount, price_currency, sku, is_active, tenant_id) VALUES
    ('eq-001', 'Digital X-Ray Machine', 'cat-004', 'mfr-001', 'DXR-5000', 'High-resolution digital X-ray system', 
     '{"resolution": "5MP", "power": "220V", "weight": "450kg"}', 125000.00, 'USD', 'SKU-XR-001', true, 'city-hospital'),
    ('eq-002', 'Ultrasound Scanner', 'cat-004', 'mfr-002', 'US-PRO-300', 'Portable ultrasound scanner with 3D imaging',
     '{"modes": ["2D", "3D", "Doppler"], "display": "15-inch LCD", "weight": "8kg"}', 35000.00, 'USD', 'SKU-US-002', true, 'city-hospital'),
    ('eq-003', 'Surgical Microscope', 'cat-002', 'mfr-003', 'SM-8000', 'Advanced surgical microscope with LED illumination',
     '{"magnification": "5x-40x", "illumination": "LED", "camera": "4K"}', 89000.00, 'USD', 'SKU-SM-003', true, 'city-hospital')
ON CONFLICT (id) DO NOTHING;

-- Insert sample suppliers
INSERT INTO suppliers (
    id, tenant_id, company_name, business_registration_number, tax_id,
    year_established, description, contact_info, address, specializations,
    certifications, performance_rating, status, verification_status, created_by
) VALUES
    (
        'sup-001', 'city-hospital', 'Premier Medical Supplies Inc', 'BRN-2018-001', 'TAX-PMS-001',
        2018, 'Leading supplier of medical diagnostic equipment',
        '{"primary_contact_name": "John Smith", "primary_contact_email": "john@premiermed.com", "primary_contact_phone": "+1-555-0101", "website": "https://premiermed.com"}',
        '{"street": "123 Medical Plaza", "city": "New York", "state": "NY", "postal_code": "10001", "country": "USA"}',
        ARRAY['cat-001', 'cat-004'],
        '[{"id": "cert-001", "name": "ISO 9001", "issuing_body": "ISO", "cert_number": "ISO-2023-001", "issue_date": "2023-01-15", "expiry_date": "2026-01-15"}]',
        4.5, 'active', 'approved', 'admin'
    ),
    (
        'sup-002', 'city-hospital', 'Surgical Instruments Corp', 'BRN-2015-002', 'TAX-SIC-002',
        2015, 'Specialized in surgical and laboratory equipment',
        '{"primary_contact_name": "Sarah Johnson", "primary_contact_email": "sarah@surgicalinstruments.com", "primary_contact_phone": "+1-555-0102", "website": "https://surgicalinstruments.com"}',
        '{"street": "456 Healthcare Ave", "city": "Boston", "state": "MA", "postal_code": "02101", "country": "USA"}',
        ARRAY['cat-002', 'cat-003'],
        '[{"id": "cert-002", "name": "FDA Registration", "issuing_body": "FDA", "cert_number": "FDA-2022-002", "issue_date": "2022-06-01", "expiry_date": "2025-06-01"}]',
        4.8, 'active', 'approved', 'admin'
    ),
    (
        'sup-003', 'city-hospital', 'Global HealthTech Channel Partners', 'BRN-2020-003', 'TAX-GHT-003',
        2020, 'International Channel Partner of cutting-edge medical technology',
        '{"primary_contact_name": "Michael Chen", "primary_contact_email": "michael@globalhealthtech.com", "primary_contact_phone": "+1-555-0103", "website": "https://globalhealthtech.com"}',
        '{"street": "789 Tech Park Drive", "city": "San Francisco", "state": "CA", "postal_code": "94102", "country": "USA"}',
        ARRAY['cat-001', 'cat-002', 'cat-004'],
        '[{"id": "cert-003", "name": "CE Marking", "issuing_body": "EU", "cert_number": "CE-2023-003", "issue_date": "2023-03-20", "expiry_date": "2028-03-20"}]',
        4.3, 'active', 'approved', 'admin'
    )
ON CONFLICT (id) DO NOTHING;

-- Insert sample RFQ
INSERT INTO rfqs (
    id, rfq_number, tenant_id, title, description, priority, status,
    delivery_terms, payment_terms, response_deadline, created_by
) VALUES
    (
        'rfq-001', 'RFQ-2025-001', 'city-hospital',
        'Purchase of Diagnostic Equipment for Radiology Department',
        'Request for quotes on diagnostic imaging equipment including X-ray and ultrasound systems',
        'high', 'published',
        '{"address": "100 Hospital Road", "city": "New York", "state": "NY", "postal_code": "10001", "country": "USA", "required_by": "2025-12-31T00:00:00Z", "installation_required": true}',
        '{"payment_method": "Net 30", "payment_days": 30, "advance_payment_percent": 20}',
        '2025-11-30T23:59:59Z', 'admin'
    )
ON CONFLICT (id) DO NOTHING;

-- Insert sample RFQ items
INSERT INTO rfq_items (
    id, rfq_id, equipment_id, category_id, name, description, specifications,
    quantity, unit, estimated_price
) VALUES
    (
        'item-001', 'rfq-001', 'eq-001', 'cat-004',
        'Digital X-Ray Machine', 'High-resolution digital X-ray system for radiology department',
        '{"resolution": "5MP minimum", "power": "220V", "warranty": "3 years"}',
        2, 'unit', 120000.00
    ),
    (
        'item-002', 'rfq-001', 'eq-002', 'cat-004',
        'Portable Ultrasound Scanner', '3D imaging capable ultrasound scanner',
        '{"modes": ["2D", "3D", "Doppler"], "portability": "required"}',
        3, 'unit', 33000.00
    )
ON CONFLICT (id) DO NOTHING;

-- Insert sample RFQ invitations
INSERT INTO rfq_invitations (
    id, rfq_id, supplier_id, status, invited_at, message
) VALUES
    ('inv-001', 'rfq-001', 'sup-001', 'invited', CURRENT_TIMESTAMP, 'Please provide your best quote for the diagnostic equipment'),
    ('inv-002', 'rfq-001', 'sup-003', 'invited', CURRENT_TIMESTAMP, 'We value your expertise in imaging equipment')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- Display summary of created tables
SELECT 
    'RFQs' as table_name, 
    COUNT(*) as record_count 
FROM rfqs
UNION ALL
SELECT 'RFQ Items', COUNT(*) FROM rfq_items
UNION ALL
SELECT 'RFQ Invitations', COUNT(*) FROM rfq_invitations
UNION ALL
SELECT 'Suppliers', COUNT(*) FROM suppliers
UNION ALL
SELECT 'Equipment', COUNT(*) FROM equipment
UNION ALL
SELECT 'Categories', COUNT(*) FROM categories
UNION ALL
SELECT 'Manufacturers', COUNT(*) FROM manufacturers;
