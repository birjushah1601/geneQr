-- Migration: 029-extend-equipment-registry.sql
-- Description: Extend equipment_registry for better integration
-- Purpose: Add foreign keys for manufacturer, catalog, and QR code linking
-- Date: 2025-12-23

-- =====================================================================
-- 1. ADD NEW COLUMNS TO EQUIPMENT_REGISTRY
-- =====================================================================

-- Add manufacturer_id (link to organizations table)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS manufacturer_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

-- Add equipment_catalog_id (link to equipment catalog)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS equipment_catalog_id UUID REFERENCES equipment_catalog(id) ON DELETE SET NULL;

-- Add customer_org_id (link customer to organizations table)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS customer_org_id UUID REFERENCES organizations(id) ON DELETE SET NULL;

-- Add qr_code_id (link to qr_codes table)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS qr_code_id UUID REFERENCES qr_codes(id) ON DELETE SET NULL;

-- Make some fields nullable for unassigned equipment
-- (Already nullable in most schemas, but ensuring consistency)
ALTER TABLE equipment_registry 
ALTER COLUMN customer_name DROP NOT NULL,
ALTER COLUMN equipment_name DROP NOT NULL,
ALTER COLUMN manufacturer_name DROP NOT NULL;

-- =====================================================================
-- 2. CREATE INDEXES FOR NEW FOREIGN KEYS
-- =====================================================================

CREATE INDEX IF NOT EXISTS idx_equipment_registry_manufacturer_id 
ON equipment_registry(manufacturer_id);

CREATE INDEX IF NOT EXISTS idx_equipment_registry_equipment_catalog_id 
ON equipment_registry(equipment_catalog_id);

CREATE INDEX IF NOT EXISTS idx_equipment_registry_customer_org_id 
ON equipment_registry(customer_org_id);

CREATE INDEX IF NOT EXISTS idx_equipment_registry_qr_code_id 
ON equipment_registry(qr_code_id);

-- =====================================================================
-- 3. UPDATE EXISTING DATA (OPTIONAL - RUN CAREFULLY)
-- =====================================================================

-- Try to link existing equipment to manufacturers by name matching
UPDATE equipment_registry er
SET manufacturer_id = o.id
FROM organizations o
WHERE er.manufacturer_id IS NULL
AND er.manufacturer_name IS NOT NULL
AND LOWER(er.manufacturer_name) = LOWER(o.name)
AND o.org_type = 'manufacturer';

-- Try to link existing equipment to customers by name matching
UPDATE equipment_registry er
SET customer_org_id = o.id
FROM organizations o
WHERE er.customer_org_id IS NULL
AND er.customer_name IS NOT NULL
AND LOWER(er.customer_name) = LOWER(o.name)
AND o.org_type IN ('hospital', 'clinic');

-- =====================================================================
-- 4. CREATE HELPER VIEWS
-- =====================================================================

-- View for equipment with full organization details
CREATE OR REPLACE VIEW equipment_registry_full AS
SELECT 
    er.*,
    mfr.name as mfr_org_name,
    mfr.org_type as mfr_org_type,
    cust.name as cust_org_name,
    cust.org_type as cust_org_type,
    ec.model_name as catalog_model_name,
    ec.category as catalog_category,
    qc.qr_code as qr_code_value,
    qc.status as qr_code_status,
    qc.batch_id as qr_batch_id
FROM equipment_registry er
LEFT JOIN organizations mfr ON er.manufacturer_id = mfr.id
LEFT JOIN organizations cust ON er.customer_org_id = cust.id
LEFT JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
LEFT JOIN qr_codes qc ON er.qr_code_id = qc.id;

COMMENT ON VIEW equipment_registry_full IS 'Equipment registry with full organization and catalog details';

-- View for unassigned equipment (equipment without customer)
CREATE OR REPLACE VIEW equipment_inventory AS
SELECT 
    er.*,
    mfr.name as manufacturer_org_name,
    ec.model_name as equipment_model,
    qc.qr_code as qr_code_value
FROM equipment_registry er
LEFT JOIN organizations mfr ON er.manufacturer_id = mfr.id
LEFT JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
LEFT JOIN qr_codes qc ON er.qr_code_id = qc.id
WHERE er.customer_org_id IS NULL
AND er.status IN ('operational', 'under_maintenance');

COMMENT ON VIEW equipment_inventory IS 'Unassigned equipment (manufacturer inventory)';

-- =====================================================================
-- 5. ADD COMMENTS FOR DOCUMENTATION
-- =====================================================================

COMMENT ON COLUMN equipment_registry.manufacturer_id IS 'Link to manufacturer organization (from organizations table)';
COMMENT ON COLUMN equipment_registry.equipment_catalog_id IS 'Link to equipment catalog (model/type reference)';
COMMENT ON COLUMN equipment_registry.customer_org_id IS 'Link to customer organization (hospital/clinic)';
COMMENT ON COLUMN equipment_registry.qr_code_id IS 'Link to QR code record (from qr_codes table)';

-- =====================================================================
-- 6. CREATE TRIGGER TO SYNC QR CODE ON INSERT
-- =====================================================================

-- When equipment is registered with a QR code, update qr_codes table
CREATE OR REPLACE FUNCTION sync_qr_code_on_equipment_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- If QR code is provided, try to find it and link
    IF NEW.qr_code IS NOT NULL THEN
        UPDATE qr_codes
        SET 
            equipment_registry_id = NEW.id,
            status = 'assigned',
            assigned_at = NOW()
        WHERE qr_code = NEW.qr_code
        AND status = 'generated';
        
        -- Also try to set qr_code_id in equipment_registry
        NEW.qr_code_id := (SELECT id FROM qr_codes WHERE qr_code = NEW.qr_code LIMIT 1);
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER equipment_registry_sync_qr_trigger
    BEFORE INSERT ON equipment_registry
    FOR EACH ROW
    EXECUTE FUNCTION sync_qr_code_on_equipment_insert();

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

DO $$
DECLARE
    v_updated_mfr INT;
    v_updated_cust INT;
BEGIN
    -- Count updated records
    SELECT COUNT(*) INTO v_updated_mfr 
    FROM equipment_registry 
    WHERE manufacturer_id IS NOT NULL;
    
    SELECT COUNT(*) INTO v_updated_cust 
    FROM equipment_registry 
    WHERE customer_org_id IS NOT NULL;
    
    RAISE NOTICE 'Migration 029: Equipment registry extended successfully';
    RAISE NOTICE '  - Added manufacturer_id column';
    RAISE NOTICE '  - Added equipment_catalog_id column';
    RAISE NOTICE '  - Added customer_org_id column';
    RAISE NOTICE '  - Added qr_code_id column';
    RAISE NOTICE '  - Linked % equipment to manufacturers', v_updated_mfr;
    RAISE NOTICE '  - Linked % equipment to customers', v_updated_cust;
    RAISE NOTICE '  - Created 2 views for queries';
    RAISE NOTICE '  - Created 1 trigger for QR sync';
END $$;
