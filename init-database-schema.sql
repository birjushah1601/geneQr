-- ABY-MED Platform - Database Initialization Script
-- This script creates the basic schema for testing all services

-- Drop existing tables if they exist (for clean reinstall)
DROP TABLE IF EXISTS service_tickets CASCADE;
DROP TABLE IF EXISTS equipment CASCADE;
DROP TABLE IF EXISTS contracts CASCADE;
DROP TABLE IF EXISTS quote_comparisons CASCADE;
DROP TABLE IF EXISTS quote_items CASCADE;
DROP TABLE IF EXISTS quotes CASCADE;
DROP TABLE IF EXISTS rfq_items CASCADE;
DROP TABLE IF EXISTS rfqs CASCADE;
DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS catalog_items CASCADE;

-- ============================================================================
-- CATALOG SERVICE TABLES
-- ============================================================================

CREATE TABLE catalog_items (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    category VARCHAR(200),
    manufacturer VARCHAR(200),
    description TEXT,
    specifications JSONB,
    base_price DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'INR',
    compliance_certifications JSONB,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    UNIQUE(tenant_id, sku)
);

CREATE INDEX idx_catalog_tenant ON catalog_items(tenant_id);
CREATE INDEX idx_catalog_category ON catalog_items(category);
CREATE INDEX idx_catalog_manufacturer ON catalog_items(manufacturer);
CREATE INDEX idx_catalog_status ON catalog_items(status);

-- ============================================================================
-- SUPPLIER SERVICE TABLES
-- ============================================================================

CREATE TABLE suppliers (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    company_name VARCHAR(500) NOT NULL,
    gstin VARCHAR(50),
    pan VARCHAR(20),
    contact_person VARCHAR(200),
    email VARCHAR(200),
    phone VARCHAR(50),
    address JSONB,
    categories JSONB,
    certifications JSONB,
    status VARCHAR(50) DEFAULT 'active',
    rating DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);

CREATE INDEX idx_suppliers_tenant ON suppliers(tenant_id);
CREATE INDEX idx_suppliers_status ON suppliers(status);
CREATE INDEX idx_suppliers_email ON suppliers(email);

-- ============================================================================
-- RFQ SERVICE TABLES
-- ============================================================================

CREATE TABLE rfqs (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    deadline TIMESTAMP,
    delivery_location JSONB,
    published_at TIMESTAMP,
    closed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);

CREATE TABLE rfq_items (
    id VARCHAR(255) PRIMARY KEY,
    rfq_id VARCHAR(255) NOT NULL REFERENCES rfqs(id) ON DELETE CASCADE,
    catalog_id VARCHAR(255),
    quantity INTEGER NOT NULL,
    specifications TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (catalog_id) REFERENCES catalog_items(id)
);

CREATE INDEX idx_rfqs_tenant ON rfqs(tenant_id);
CREATE INDEX idx_rfqs_status ON rfqs(status);
CREATE INDEX idx_rfqs_deadline ON rfqs(deadline);
CREATE INDEX idx_rfq_items_rfq ON rfq_items(rfq_id);

-- ============================================================================
-- QUOTE SERVICE TABLES
-- ============================================================================

CREATE TABLE quotes (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    rfq_id VARCHAR(255) NOT NULL REFERENCES rfqs(id),
    supplier_id VARCHAR(255) NOT NULL REFERENCES suppliers(id),
    validity_days INTEGER,
    payment_terms TEXT,
    delivery_timeline VARCHAR(200),
    warranty TEXT,
    notes TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    submitted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255)
);

CREATE TABLE quote_items (
    id VARCHAR(255) PRIMARY KEY,
    quote_id VARCHAR(255) NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    rfq_item_id VARCHAR(255) NOT NULL REFERENCES rfq_items(id),
    unit_price DECIMAL(15,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    tax_percent DECIMAL(5,2) DEFAULT 18,
    delivery_charges DECIMAL(15,2) DEFAULT 0,
    installation_charges DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_quotes_tenant ON quotes(tenant_id);
CREATE INDEX idx_quotes_rfq ON quotes(rfq_id);
CREATE INDEX idx_quotes_supplier ON quotes(supplier_id);
CREATE INDEX idx_quotes_status ON quotes(status);
CREATE INDEX idx_quote_items_quote ON quote_items(quote_id);

-- ============================================================================
-- COMPARISON SERVICE TABLES
-- ============================================================================

CREATE TABLE quote_comparisons (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    rfq_id VARCHAR(255) NOT NULL REFERENCES rfqs(id),
    quote_ids JSONB NOT NULL,
    comparison_criteria JSONB,
    comparison_results JSONB,
    recommended_quote_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255)
);

CREATE INDEX idx_comparisons_tenant ON quote_comparisons(tenant_id);
CREATE INDEX idx_comparisons_rfq ON quote_comparisons(rfq_id);

-- ============================================================================
-- CONTRACT SERVICE TABLES
-- ============================================================================

CREATE TABLE contracts (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    quote_id VARCHAR(255) NOT NULL REFERENCES quotes(id),
    contract_number VARCHAR(100),
    buyer_signatory VARCHAR(200),
    supplier_signatory VARCHAR(200),
    special_terms JSONB,
    payment_schedule JSONB,
    status VARCHAR(50) DEFAULT 'draft',
    signed_at TIMESTAMP,
    effective_date DATE,
    expiry_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    UNIQUE(tenant_id, contract_number)
);

CREATE INDEX idx_contracts_tenant ON contracts(tenant_id);
CREATE INDEX idx_contracts_quote ON contracts(quote_id);
CREATE INDEX idx_contracts_status ON contracts(status);

-- ============================================================================
-- EQUIPMENT REGISTRY TABLES
-- ============================================================================

CREATE TABLE equipment (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    serial_number VARCHAR(200),
    model VARCHAR(200),
    manufacturer VARCHAR(200),
    purchase_date DATE,
    installation_date DATE,
    location JSONB,
    warranty_expiry DATE,
    maintenance_schedule VARCHAR(100),
    status VARCHAR(50) DEFAULT 'active',
    qr_code_path VARCHAR(500),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    UNIQUE(tenant_id, serial_number)
);

CREATE INDEX idx_equipment_tenant ON equipment(tenant_id);
CREATE INDEX idx_equipment_status ON equipment(status);
CREATE INDEX idx_equipment_serial ON equipment(serial_number);
CREATE INDEX idx_equipment_manufacturer ON equipment(manufacturer);

-- ============================================================================
-- SERVICE TICKET TABLES
-- ============================================================================

CREATE TABLE service_tickets (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    equipment_id VARCHAR(255) REFERENCES equipment(id),
    ticket_number VARCHAR(100),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority VARCHAR(50) DEFAULT 'medium',
    status VARCHAR(50) DEFAULT 'open',
    reported_by VARCHAR(200),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(200),
    assigned_to VARCHAR(255),
    resolution TEXT,
    parts_used JSONB,
    labor_hours DECIMAL(5,2),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    UNIQUE(tenant_id, ticket_number)
);

CREATE TABLE ticket_comments (
    id VARCHAR(255) PRIMARY KEY,
    ticket_id VARCHAR(255) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(255)
);

CREATE INDEX idx_tickets_tenant ON service_tickets(tenant_id);
CREATE INDEX idx_tickets_equipment ON service_tickets(equipment_id);
CREATE INDEX idx_tickets_status ON service_tickets(status);
CREATE INDEX idx_tickets_priority ON service_tickets(priority);
CREATE INDEX idx_ticket_comments_ticket ON ticket_comments(ticket_id);

-- ============================================================================
-- INSERT SAMPLE DATA FOR TESTING
-- ============================================================================

-- Sample Catalog Items
INSERT INTO catalog_items (id, tenant_id, name, sku, category, manufacturer, description, base_price, currency, status) VALUES
('cat-001', 'city-hospital', 'MRI Scanner - Siemens Magnetom', 'MRI-001-SIEMENS', 'Diagnostic Imaging', 'Siemens Healthineers', '1.5T MRI Scanner with advanced imaging capabilities', 1500000.00, 'INR', 'active'),
('cat-002', 'city-hospital', 'CT Scanner - GE Revolution', 'CT-001-GE', 'Diagnostic Imaging', 'GE Healthcare', '128-slice CT scanner with low dose technology', 2500000.00, 'INR', 'active'),
('cat-003', 'city-hospital', 'Ultrasound - Philips EPIQ', 'US-001-PHILIPS', 'Diagnostic Imaging', 'Philips Healthcare', 'Premium ultrasound system', 750000.00, 'INR', 'active');

-- Sample Suppliers
INSERT INTO suppliers (id, tenant_id, company_name, email, phone, status, rating) VALUES
('sup-001', 'city-hospital', 'MedTech Supplies Pvt Ltd', 'info@medtechsupplies.com', '+91-9876543210', 'active', 4.5),
('sup-002', 'city-hospital', 'Healthcare Solutions India', 'contact@healthcaresolutions.in', '+91-9876543211', 'active', 4.2),
('sup-003', 'city-hospital', 'Advanced Medical Equipment Co', 'sales@advmedequip.com', '+91-9876543212', 'active', 4.8);

-- Sample Equipment
INSERT INTO equipment (id, tenant_id, name, serial_number, model, manufacturer, status) VALUES
('eq-001', 'city-hospital', 'MRI Scanner Unit 1', 'MRI-UNIT-001', 'Magnetom Skyra 1.5T', 'Siemens Healthineers', 'operational'),
('eq-002', 'city-hospital', 'CT Scanner Unit 1', 'CT-UNIT-001', 'Revolution 128', 'GE Healthcare', 'operational');

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

-- Grant permissions to postgres user (adjust as needed)
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- Show all tables
SELECT 
    table_name, 
    (SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
FROM information_schema.tables t
WHERE table_schema = 'public' 
    AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Show sample data
SELECT 'Catalog Items' as table_name, COUNT(*) as record_count FROM catalog_items
UNION ALL
SELECT 'Suppliers', COUNT(*) FROM suppliers
UNION ALL
SELECT 'Equipment', COUNT(*) FROM equipment
UNION ALL
SELECT 'RFQs', COUNT(*) FROM rfqs
UNION ALL
SELECT 'Quotes', COUNT(*) FROM quotes
UNION ALL
SELECT 'Service Tickets', COUNT(*) FROM service_tickets;

-- Success message
SELECT 'Database schema initialized successfully!' as status;
