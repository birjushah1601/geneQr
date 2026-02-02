-- ============================================================================
-- ServQR Platform Catalog Module - Database Initialization Script
-- ============================================================================
-- This script sets up the catalog database schema, loads sample data,
-- creates indexes, and configures row-level security for multi-tenancy.
-- ============================================================================

\echo 'Starting catalog module database initialization...'

-- Start a transaction for atomicity
BEGIN;

-- Create schema if it doesn't exist
\echo 'Creating schema...'
CREATE SCHEMA IF NOT EXISTS public;

-- Enable Row Level Security
\echo 'Enabling row level security...'
ALTER DATABASE medplatform SET row_security = on;

-- Create a function to set tenant context
\echo 'Creating tenant context function...'
CREATE OR REPLACE FUNCTION set_tenant_context(text)
RETURNS text AS $$
BEGIN
  PERFORM set_config('app.current_tenant', $1, false);
  RETURN $1;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create tables if they don't exist
\echo 'Creating tables if they do not exist...'

-- Manufacturers table
CREATE TABLE IF NOT EXISTS manufacturers (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    headquarters VARCHAR(255) NOT NULL,
    website VARCHAR(255),
    specialization VARCHAR(255) NOT NULL,
    established INT,
    description TEXT,
    country VARCHAR(50) DEFAULT 'India',
    tenant_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Categories table with hierarchical structure
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(26),
    description TEXT,
    tenant_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Equipment table
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(100) NOT NULL,
    category_id VARCHAR(26) NOT NULL,
    manufacturer_id VARCHAR(26) NOT NULL,
    description TEXT,
    specifications JSONB NOT NULL,
    price_amount DECIMAL(12, 2) NOT NULL,
    price_currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    sku VARCHAR(50),
    images TEXT[],
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tenant_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id)
);

-- Create indexes for better performance
\echo 'Creating indexes for performance...'
CREATE INDEX IF NOT EXISTS idx_manufacturers_tenant ON manufacturers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_manufacturers_specialization ON manufacturers(specialization);
CREATE INDEX IF NOT EXISTS idx_manufacturers_name ON manufacturers(name);

CREATE INDEX IF NOT EXISTS idx_categories_tenant ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

CREATE INDEX IF NOT EXISTS idx_equipment_tenant ON equipment(tenant_id);
CREATE INDEX IF NOT EXISTS idx_equipment_category ON equipment(category_id);
CREATE INDEX IF NOT EXISTS idx_equipment_manufacturer ON equipment(manufacturer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_name ON equipment(name);
CREATE INDEX IF NOT EXISTS idx_equipment_model ON equipment(model);
CREATE INDEX IF NOT EXISTS idx_equipment_sku ON equipment(sku);
CREATE INDEX IF NOT EXISTS idx_equipment_active ON equipment(is_active);
CREATE INDEX IF NOT EXISTS idx_equipment_price ON equipment(price_amount);
CREATE INDEX IF NOT EXISTS idx_equipment_specs ON equipment USING gin(specifications);

-- Set up Row Level Security policies
\echo 'Setting up row-level security policies...'

-- Enable RLS on tables
ALTER TABLE manufacturers ENABLE ROW LEVEL SECURITY;
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE equipment ENABLE ROW LEVEL SECURITY;

-- Create policies
CREATE POLICY manufacturers_tenant_isolation ON manufacturers
    USING (tenant_id = current_setting('app.current_tenant', true));

CREATE POLICY categories_tenant_isolation ON categories
    USING (tenant_id = current_setting('app.current_tenant', true));

CREATE POLICY equipment_tenant_isolation ON equipment
    USING (tenant_id = current_setting('app.current_tenant', true));

-- Create superuser bypass policies (for admins/migrations)
CREATE POLICY manufacturers_superuser_bypass ON manufacturers
    USING (current_user = 'postgres');

CREATE POLICY categories_superuser_bypass ON categories
    USING (current_user = 'postgres');

CREATE POLICY equipment_superuser_bypass ON equipment
    USING (current_user = 'postgres');

-- Load sample data using \i commands
\echo 'Loading manufacturer data...'
\i /docker-entrypoint-initdb.d/indian-manufacturers.sql

\echo 'Loading categories data...'
\i /docker-entrypoint-initdb.d/medical-equipment-categories.sql

\echo 'Loading sample equipment data...'
\i /docker-entrypoint-initdb.d/sample-equipment-catalog.sql

-- Create trigger functions for updated_at timestamps
\echo 'Creating timestamp trigger functions...'
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at timestamps
CREATE TRIGGER update_manufacturers_timestamp
BEFORE UPDATE ON manufacturers
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_categories_timestamp
BEFORE UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_equipment_timestamp
BEFORE UPDATE ON equipment
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- Create search functions for equipment
\echo 'Creating search functions...'
CREATE OR REPLACE FUNCTION search_equipment(
    p_query TEXT,
    p_category_ids TEXT[],
    p_manufacturer_ids TEXT[],
    p_price_min DECIMAL,
    p_price_max DECIMAL,
    p_is_active BOOLEAN,
    p_tenant_id TEXT,
    p_limit INT DEFAULT 20,
    p_offset INT DEFAULT 0
)
RETURNS TABLE (
    id VARCHAR,
    name VARCHAR,
    model VARCHAR,
    category_id VARCHAR,
    category_name VARCHAR,
    manufacturer_id VARCHAR,
    manufacturer_name VARCHAR,
    description TEXT,
    specifications JSONB,
    price_amount DECIMAL,
    price_currency VARCHAR,
    sku VARCHAR,
    images TEXT[],
    is_active BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id, e.name, e.model, e.category_id, c.name AS category_name,
        e.manufacturer_id, m.name AS manufacturer_name,
        e.description, e.specifications, e.price_amount, e.price_currency,
        e.sku, e.images, e.is_active, e.created_at, e.updated_at
    FROM equipment e
    JOIN categories c ON e.category_id = c.id
    JOIN manufacturers m ON e.manufacturer_id = m.id
    WHERE e.tenant_id = p_tenant_id
    AND (p_query IS NULL OR e.name ILIKE '%' || p_query || '%' OR e.description ILIKE '%' || p_query || '%' OR e.model ILIKE '%' || p_query || '%')
    AND (p_category_ids IS NULL OR e.category_id = ANY(p_category_ids))
    AND (p_manufacturer_ids IS NULL OR e.manufacturer_id = ANY(p_manufacturer_ids))
    AND (p_price_min IS NULL OR e.price_amount >= p_price_min)
    AND (p_price_max IS NULL OR e.price_amount <= p_price_max)
    AND (p_is_active IS NULL OR e.is_active = p_is_active)
    ORDER BY e.name
    LIMIT p_limit
    OFFSET p_offset;
END;
$$ LANGUAGE plpgsql;

-- Create function to count search results
CREATE OR REPLACE FUNCTION count_search_equipment(
    p_query TEXT,
    p_category_ids TEXT[],
    p_manufacturer_ids TEXT[],
    p_price_min DECIMAL,
    p_price_max DECIMAL,
    p_is_active BOOLEAN,
    p_tenant_id TEXT
)
RETURNS INT AS $$
DECLARE
    result_count INT;
BEGIN
    SELECT COUNT(*)
    INTO result_count
    FROM equipment e
    WHERE e.tenant_id = p_tenant_id
    AND (p_query IS NULL OR e.name ILIKE '%' || p_query || '%' OR e.description ILIKE '%' || p_query || '%' OR e.model ILIKE '%' || p_query || '%')
    AND (p_category_ids IS NULL OR e.category_id = ANY(p_category_ids))
    AND (p_manufacturer_ids IS NULL OR e.manufacturer_id = ANY(p_manufacturer_ids))
    AND (p_price_min IS NULL OR e.price_amount >= p_price_min)
    AND (p_price_max IS NULL OR e.price_amount <= p_price_max)
    AND (p_is_active IS NULL OR e.is_active = p_is_active);
    
    RETURN result_count;
END;
$$ LANGUAGE plpgsql;

-- Create default tenants if they don't exist
\echo 'Creating default tenants...'
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'tenants') THEN
        CREATE TABLE tenants (
            id VARCHAR(50) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        
        INSERT INTO tenants (id, name, description) VALUES
        ('demo-hospital', 'Demo Hospital', 'Default demo tenant for testing'),
        ('city-hospital', 'City Hospital', 'Secondary demo tenant for multi-tenant testing');
    END IF;
END $$;

-- Commit the transaction
COMMIT;

\echo 'Catalog module database initialization completed successfully.'
