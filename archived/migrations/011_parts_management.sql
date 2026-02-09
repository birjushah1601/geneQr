-- Migration: Parts Management System
-- Description: Complete parts catalog, variants, compatibility, accessories, inventory, and suppliers
-- Dependencies: equipment_types table
-- Phase: 2C - Foundation for AI Parts Recommender

-- ============================================================================
-- 1. EQUIPMENT VARIANTS (Installation Types)
-- ============================================================================

CREATE TABLE IF NOT EXISTS equipment_variants (
    variant_id BIGSERIAL PRIMARY KEY,
    equipment_type_id BIGINT NOT NULL REFERENCES equipment_types(equipment_type_id) ON DELETE CASCADE,
    
    -- Variant details
    variant_name VARCHAR(100) NOT NULL, -- "ICU Installation", "Portable", "Mobile", "Wall-Mounted"
    variant_code VARCHAR(50) NOT NULL,
    description TEXT,
    
    -- Characteristics
    is_stationary BOOLEAN DEFAULT true,
    requires_installation BOOLEAN DEFAULT true,
    environment_type VARCHAR(50), -- ICU, Ward, Emergency, Outpatient, etc.
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(equipment_type_id, variant_code)
);

CREATE INDEX idx_equipment_variants_type ON equipment_variants(equipment_type_id);
CREATE INDEX idx_equipment_variants_active ON equipment_variants(is_active) WHERE is_active = true;
CREATE INDEX idx_equipment_variants_environment ON equipment_variants(environment_type);

COMMENT ON TABLE equipment_variants IS 'Equipment installation types and variants';
COMMENT ON COLUMN equipment_variants.variant_name IS 'Display name (e.g., "ICU Installation", "Portable Unit")';
COMMENT ON COLUMN equipment_variants.environment_type IS 'Where this variant is typically used';

-- ============================================================================
-- 2. PARTS CATALOG (Master List)
-- ============================================================================

CREATE TABLE IF NOT EXISTS parts_catalog (
    part_id BIGSERIAL PRIMARY KEY,
    
    -- Part identification
    part_number VARCHAR(100) UNIQUE NOT NULL,
    part_name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Classification
    category VARCHAR(100) NOT NULL, -- Filter, Valve, Sensor, Circuit Board, Battery, Cable, etc.
    subcategory VARCHAR(100),
    
    -- Part type
    part_type VARCHAR(50) NOT NULL, -- Component, Accessory, Consumable, Tool
    is_oem_part BOOLEAN DEFAULT true, -- Original Equipment Manufacturer
    is_universal BOOLEAN DEFAULT false, -- Fits multiple equipment types
    
    -- Specifications
    specifications JSONB, -- Technical specs (voltage, size, material, etc.)
    
    -- Manufacturer info
    manufacturer_name VARCHAR(200),
    manufacturer_part_number VARCHAR(100),
    
    -- Lifecycle
    is_discontinued BOOLEAN DEFAULT false,
    replacement_part_id BIGINT REFERENCES parts_catalog(part_id),
    
    -- Pricing (reference only)
    unit_price DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'USD',
    
    -- Images
    image_url TEXT,
    datasheet_url TEXT,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_parts_catalog_number ON parts_catalog(part_number);
CREATE INDEX idx_parts_catalog_name ON parts_catalog(part_name);
CREATE INDEX idx_parts_catalog_category ON parts_catalog(category, subcategory);
CREATE INDEX idx_parts_catalog_type ON parts_catalog(part_type);
CREATE INDEX idx_parts_catalog_active ON parts_catalog(is_active) WHERE is_active = true;
CREATE INDEX idx_parts_catalog_manufacturer ON parts_catalog(manufacturer_name);

-- GIN index for specifications JSONB
CREATE INDEX idx_parts_catalog_specs ON parts_catalog USING GIN (specifications);

COMMENT ON TABLE parts_catalog IS 'Master catalog of all parts and accessories';
COMMENT ON COLUMN parts_catalog.part_type IS 'Component (replacement part), Accessory (add-on), Consumable (regular replacement), Tool';
COMMENT ON COLUMN parts_catalog.is_universal IS 'Part that fits multiple equipment types';
COMMENT ON COLUMN parts_catalog.specifications IS 'Technical specifications as JSONB';

-- ============================================================================
-- 3. EQUIPMENT PARTS COMPATIBILITY
-- ============================================================================

CREATE TABLE IF NOT EXISTS equipment_parts (
    equipment_part_id BIGSERIAL PRIMARY KEY,
    
    -- Relationships
    equipment_type_id BIGINT NOT NULL REFERENCES equipment_types(equipment_type_id) ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts_catalog(part_id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES equipment_variants(variant_id) ON DELETE SET NULL, -- NULL = all variants
    
    -- Compatibility details
    is_standard_part BOOLEAN DEFAULT true, -- Standard vs optional
    is_critical_part BOOLEAN DEFAULT false, -- Critical for operation
    quantity_per_unit INTEGER DEFAULT 1, -- How many per equipment
    
    -- Replacement info
    recommended_replacement_interval INTERVAL, -- e.g., '6 months', '1 year'
    replacement_interval_hours INTEGER, -- Operating hours before replacement
    replacement_interval_cycles INTEGER, -- Cycles before replacement
    
    -- Usage notes
    installation_notes TEXT,
    compatibility_notes TEXT,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(equipment_type_id, part_id, variant_id)
);

CREATE INDEX idx_equipment_parts_equipment ON equipment_parts(equipment_type_id);
CREATE INDEX idx_equipment_parts_part ON equipment_parts(part_id);
CREATE INDEX idx_equipment_parts_variant ON equipment_parts(variant_id);
CREATE INDEX idx_equipment_parts_critical ON equipment_parts(is_critical_part) WHERE is_critical_part = true;
CREATE INDEX idx_equipment_parts_active ON equipment_parts(is_active) WHERE is_active = true;

COMMENT ON TABLE equipment_parts IS 'Maps parts to equipment types with compatibility rules';
COMMENT ON COLUMN equipment_parts.variant_id IS 'NULL means compatible with all variants';
COMMENT ON COLUMN equipment_parts.is_standard_part IS 'Standard (included) vs optional part';
COMMENT ON COLUMN equipment_parts.is_critical_part IS 'Critical for equipment operation';

-- ============================================================================
-- 4. ACCESSORIES BY VARIANT (for Sales)
-- ============================================================================

CREATE TABLE IF NOT EXISTS equipment_accessories (
    accessory_id BIGSERIAL PRIMARY KEY,
    
    -- Relationships
    equipment_type_id BIGINT NOT NULL REFERENCES equipment_types(equipment_type_id) ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts_catalog(part_id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES equipment_variants(variant_id) ON DELETE CASCADE, -- Specific variant
    
    -- Accessory details
    is_recommended BOOLEAN DEFAULT true,
    is_required_for_variant BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    
    -- Sales info
    upsell_priority INTEGER DEFAULT 0, -- Higher = show first
    bundle_discount_percent DECIMAL(5,2),
    
    -- Marketing
    marketing_description TEXT,
    benefits TEXT[], -- Array of benefit points
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(equipment_type_id, part_id, variant_id)
);

CREATE INDEX idx_equipment_accessories_equipment ON equipment_accessories(equipment_type_id);
CREATE INDEX idx_equipment_accessories_variant ON equipment_accessories(variant_id);
CREATE INDEX idx_equipment_accessories_part ON equipment_accessories(part_id);
CREATE INDEX idx_equipment_accessories_priority ON equipment_accessories(upsell_priority DESC);
CREATE INDEX idx_equipment_accessories_active ON equipment_accessories(is_active) WHERE is_active = true;

COMMENT ON TABLE equipment_accessories IS 'Accessories for sale based on equipment variant';
COMMENT ON COLUMN equipment_accessories.upsell_priority IS 'Higher number = show first in recommendations';
COMMENT ON COLUMN equipment_accessories.is_required_for_variant IS 'Required for this installation type';

-- ============================================================================
-- 5. PARTS SUPPLIERS
-- ============================================================================

CREATE TABLE IF NOT EXISTS parts_suppliers (
    supplier_id BIGSERIAL PRIMARY KEY,
    
    -- Supplier details
    supplier_name VARCHAR(200) NOT NULL,
    supplier_code VARCHAR(50) UNIQUE,
    
    -- Contact info
    contact_person VARCHAR(200),
    email VARCHAR(255),
    phone VARCHAR(50),
    
    -- Address
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    
    -- Business details
    is_oem_supplier BOOLEAN DEFAULT false,
    is_authorized BOOLEAN DEFAULT true,
    quality_rating INTEGER CHECK (quality_rating >= 1 AND quality_rating <= 5),
    
    -- Terms
    payment_terms TEXT,
    lead_time_days INTEGER,
    minimum_order_value DECIMAL(10,2),
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_parts_suppliers_name ON parts_suppliers(supplier_name);
CREATE INDEX idx_parts_suppliers_active ON parts_suppliers(is_active) WHERE is_active = true;

COMMENT ON TABLE parts_suppliers IS 'Suppliers for parts and accessories';
COMMENT ON COLUMN parts_suppliers.is_oem_supplier IS 'Original equipment manufacturer supplier';
COMMENT ON COLUMN parts_suppliers.quality_rating IS 'Quality rating 1-5 stars';

-- ============================================================================
-- 6. SUPPLIER PARTS (Supplier-specific pricing and availability)
-- ============================================================================

CREATE TABLE IF NOT EXISTS supplier_parts (
    supplier_part_id BIGSERIAL PRIMARY KEY,
    
    -- Relationships
    supplier_id BIGINT NOT NULL REFERENCES parts_suppliers(supplier_id) ON DELETE CASCADE,
    part_id BIGINT NOT NULL REFERENCES parts_catalog(part_id) ON DELETE CASCADE,
    
    -- Supplier-specific details
    supplier_part_number VARCHAR(100),
    supplier_description TEXT,
    
    -- Pricing
    unit_price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    minimum_order_quantity INTEGER DEFAULT 1,
    bulk_price DECIMAL(10,2), -- Price for bulk orders
    bulk_quantity_threshold INTEGER,
    
    -- Availability
    lead_time_days INTEGER,
    is_in_stock BOOLEAN DEFAULT true,
    last_stock_check TIMESTAMP WITH TIME ZONE,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_preferred BOOLEAN DEFAULT false, -- Preferred supplier for this part
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(supplier_id, part_id)
);

CREATE INDEX idx_supplier_parts_supplier ON supplier_parts(supplier_id);
CREATE INDEX idx_supplier_parts_part ON supplier_parts(part_id);
CREATE INDEX idx_supplier_parts_active ON supplier_parts(is_active) WHERE is_active = true;
CREATE INDEX idx_supplier_parts_preferred ON supplier_parts(is_preferred) WHERE is_preferred = true;

COMMENT ON TABLE supplier_parts IS 'Supplier-specific pricing and availability for parts';
COMMENT ON COLUMN supplier_parts.is_preferred IS 'Preferred supplier for this part';

-- ============================================================================
-- 7. PARTS INVENTORY
-- ============================================================================

CREATE TABLE IF NOT EXISTS parts_inventory (
    inventory_id BIGSERIAL PRIMARY KEY,
    
    -- Part reference
    part_id BIGINT NOT NULL REFERENCES parts_catalog(part_id) ON DELETE CASCADE,
    
    -- Location
    warehouse_location VARCHAR(100),
    bin_location VARCHAR(50),
    
    -- Quantity
    quantity_on_hand INTEGER DEFAULT 0,
    quantity_reserved INTEGER DEFAULT 0, -- Reserved for orders
    quantity_available INTEGER GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
    
    -- Thresholds
    reorder_point INTEGER DEFAULT 5,
    reorder_quantity INTEGER DEFAULT 10,
    minimum_stock_level INTEGER DEFAULT 0,
    maximum_stock_level INTEGER DEFAULT 100,
    
    -- Cost tracking
    average_cost DECIMAL(10,2),
    last_purchase_price DECIMAL(10,2),
    last_purchase_date TIMESTAMP WITH TIME ZONE,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Timestamps
    last_stock_count TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(part_id, warehouse_location)
);

CREATE INDEX idx_parts_inventory_part ON parts_inventory(part_id);
CREATE INDEX idx_parts_inventory_location ON parts_inventory(warehouse_location);
CREATE INDEX idx_parts_inventory_low_stock ON parts_inventory(quantity_available) WHERE quantity_available <= reorder_point;

COMMENT ON TABLE parts_inventory IS 'Current inventory levels for parts';
COMMENT ON COLUMN parts_inventory.quantity_available IS 'Auto-calculated: on_hand - reserved';
COMMENT ON COLUMN parts_inventory.reorder_point IS 'Reorder when quantity falls below this';

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Update timestamp triggers
CREATE OR REPLACE FUNCTION update_parts_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_equipment_variants_timestamp
    BEFORE UPDATE ON equipment_variants
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_parts_catalog_timestamp
    BEFORE UPDATE ON parts_catalog
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_equipment_parts_timestamp
    BEFORE UPDATE ON equipment_parts
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_equipment_accessories_timestamp
    BEFORE UPDATE ON equipment_accessories
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_parts_suppliers_timestamp
    BEFORE UPDATE ON parts_suppliers
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_supplier_parts_timestamp
    BEFORE UPDATE ON supplier_parts
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

CREATE TRIGGER trg_update_parts_inventory_timestamp
    BEFORE UPDATE ON parts_inventory
    FOR EACH ROW EXECUTE FUNCTION update_parts_timestamp();

-- ============================================================================
-- VIEWS
-- ============================================================================

-- Complete parts view with compatibility
CREATE OR REPLACE VIEW v_parts_compatibility AS
SELECT 
    pc.part_id,
    pc.part_number,
    pc.part_name,
    pc.category,
    pc.part_type,
    et.equipment_type_id,
    et.equipment_name,
    ev.variant_id,
    ev.variant_name,
    ep.is_standard_part,
    ep.is_critical_part,
    ep.recommended_replacement_interval,
    pi.quantity_available,
    pi.reorder_point,
    CASE 
        WHEN pi.quantity_available <= pi.reorder_point THEN 'Low Stock'
        WHEN pi.quantity_available = 0 THEN 'Out of Stock'
        ELSE 'In Stock'
    END as stock_status
FROM parts_catalog pc
JOIN equipment_parts ep ON pc.part_id = ep.part_id
JOIN equipment_types et ON ep.equipment_type_id = et.equipment_type_id
LEFT JOIN equipment_variants ev ON ep.variant_id = ev.variant_id
LEFT JOIN parts_inventory pi ON pc.part_id = pi.part_id
WHERE pc.is_active = true AND ep.is_active = true;

COMMENT ON VIEW v_parts_compatibility IS 'Complete view of parts with equipment compatibility and stock status';

-- Accessories for sale
CREATE OR REPLACE VIEW v_accessories_catalog AS
SELECT 
    ea.accessory_id,
    et.equipment_type_id,
    et.equipment_name,
    ev.variant_id,
    ev.variant_name,
    pc.part_id,
    pc.part_number,
    pc.part_name,
    pc.description,
    pc.category,
    pc.unit_price,
    ea.is_recommended,
    ea.is_required_for_variant,
    ea.upsell_priority,
    ea.bundle_discount_percent,
    ea.marketing_description,
    ea.benefits,
    pi.quantity_available as stock_available
FROM equipment_accessories ea
JOIN equipment_types et ON ea.equipment_type_id = et.equipment_type_id
LEFT JOIN equipment_variants ev ON ea.variant_id = ev.variant_id
JOIN parts_catalog pc ON ea.part_id = pc.part_id
LEFT JOIN parts_inventory pi ON pc.part_id = pi.part_id
WHERE ea.is_active = true AND pc.is_active = true
ORDER BY ea.upsell_priority DESC;

COMMENT ON VIEW v_accessories_catalog IS 'Accessories available for sale by equipment variant';

-- Supplier pricing comparison
CREATE OR REPLACE VIEW v_supplier_pricing AS
SELECT 
    pc.part_id,
    pc.part_number,
    pc.part_name,
    ps.supplier_id,
    ps.supplier_name,
    ps.is_oem_supplier,
    sp.unit_price,
    sp.currency,
    sp.lead_time_days,
    sp.is_in_stock,
    sp.is_preferred,
    RANK() OVER (PARTITION BY pc.part_id ORDER BY sp.unit_price ASC) as price_rank
FROM parts_catalog pc
JOIN supplier_parts sp ON pc.part_id = sp.part_id
JOIN parts_suppliers ps ON sp.supplier_id = ps.supplier_id
WHERE pc.is_active = true 
  AND ps.is_active = true 
  AND sp.is_active = true;

COMMENT ON VIEW v_supplier_pricing IS 'Supplier pricing comparison with ranking';

-- Low stock alerts
CREATE OR REPLACE VIEW v_low_stock_parts AS
SELECT 
    pi.inventory_id,
    pc.part_id,
    pc.part_number,
    pc.part_name,
    pc.category,
    pi.quantity_available,
    pi.reorder_point,
    pi.reorder_quantity,
    pi.warehouse_location,
    CASE 
        WHEN pi.quantity_available = 0 THEN 'Critical - Out of Stock'
        WHEN pi.quantity_available <= pi.minimum_stock_level THEN 'Urgent - Below Minimum'
        WHEN pi.quantity_available <= pi.reorder_point THEN 'Low - Reorder Needed'
    END as alert_level
FROM parts_inventory pi
JOIN parts_catalog pc ON pi.part_id = pc.part_id
WHERE pi.quantity_available <= pi.reorder_point
  AND pc.is_active = true
  AND pi.is_active = true
ORDER BY pi.quantity_available ASC;

COMMENT ON VIEW v_low_stock_parts IS 'Parts with low or out of stock alerts';

-- ============================================================================
-- SAMPLE DATA COMMENTS (for reference)
-- ============================================================================

-- Example: Ventilator variants
-- INSERT INTO equipment_variants (equipment_type_id, variant_name, variant_code, environment_type, is_stationary)
-- VALUES 
--   (1, 'ICU Installation', 'ICU-VENT', 'ICU', true),
--   (1, 'Portable Unit', 'PORT-VENT', 'General', false);

-- Example: Parts
-- INSERT INTO parts_catalog (part_number, part_name, category, part_type, manufacturer_name, unit_price)
-- VALUES 
--   ('VENT-FILT-001', 'HEPA Filter Assembly', 'Filter', 'Component', 'MedTech Inc', 45.00),
--   ('VENT-VALVE-001', 'Pressure Relief Valve', 'Valve', 'Component', 'MedTech Inc', 85.00),
--   ('VENT-ACC-HUM', 'Heated Humidifier (ICU)', 'Humidifier', 'Accessory', 'MedTech Inc', 450.00),
--   ('VENT-ACC-BAT', 'Portable Battery Pack', 'Battery', 'Accessory', 'MedTech Inc', 350.00);

-- Example: Equipment parts compatibility
-- INSERT INTO equipment_parts (equipment_type_id, part_id, variant_id, is_standard_part, recommended_replacement_interval)
-- VALUES 
--   (1, 1, NULL, true, '6 months'), -- Filter for all variants
--   (1, 2, NULL, true, '1 year');   -- Valve for all variants

-- Example: Accessories by variant
-- INSERT INTO equipment_accessories (equipment_type_id, part_id, variant_id, is_recommended, upsell_priority)
-- VALUES 
--   (1, 3, 1, true, 10), -- Humidifier for ICU variant
--   (1, 4, 2, true, 10); -- Battery for Portable variant
