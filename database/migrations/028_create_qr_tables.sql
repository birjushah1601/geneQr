-- Migration: 028-create-qr-tables.sql
-- Description: Create QR code lifecycle management tables
-- Purpose: Support QR code generation without equipment assignment
-- Date: 2025-12-23

-- =====================================================================
-- 1. QR BATCHES TABLE
-- =====================================================================
-- Tracks bulk QR code generation batches

CREATE TABLE IF NOT EXISTS qr_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_number VARCHAR(100) UNIQUE NOT NULL,
    
    -- Generation details
    manufacturer_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    equipment_catalog_id UUID REFERENCES equipment_catalog(id) ON DELETE SET NULL,
    
    -- Batch info
    quantity_requested INT NOT NULL CHECK (quantity_requested > 0),
    quantity_generated INT NOT NULL DEFAULT 0 CHECK (quantity_generated >= 0),
    start_serial_number VARCHAR(255),
    end_serial_number VARCHAR(255),
    
    -- Files
    pdf_url TEXT,
    csv_url TEXT,
    
    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Status lifecycle: pending → generating → completed → failed
    
    -- Metadata
    generated_at TIMESTAMPTZ,
    generated_by VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_qr_batch_status CHECK (status IN ('pending', 'generating', 'completed', 'failed')),
    CONSTRAINT chk_qr_batch_quantities CHECK (quantity_generated <= quantity_requested)
);

-- Indexes for qr_batches
CREATE INDEX idx_qr_batches_manufacturer ON qr_batches(manufacturer_id);
CREATE INDEX idx_qr_batches_equipment_catalog ON qr_batches(equipment_catalog_id);
CREATE INDEX idx_qr_batches_status ON qr_batches(status);
CREATE INDEX idx_qr_batches_created_at ON qr_batches(created_at DESC);

-- Comments
COMMENT ON TABLE qr_batches IS 'Tracks bulk QR code generation batches';
COMMENT ON COLUMN qr_batches.status IS 'Batch generation status: pending|generating|completed|failed';
COMMENT ON COLUMN qr_batches.quantity_requested IS 'Number of QR codes requested in this batch';
COMMENT ON COLUMN qr_batches.quantity_generated IS 'Number of QR codes successfully generated';

-- =====================================================================
-- 2. QR CODES TABLE
-- =====================================================================
-- Individual QR code records with lifecycle management

CREATE TABLE IF NOT EXISTS qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    qr_code_url TEXT NOT NULL,
    qr_image_url TEXT,
    
    -- Optional linking (NULL until assigned)
    equipment_catalog_id UUID REFERENCES equipment_catalog(id) ON DELETE SET NULL,
    manufacturer_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    batch_id UUID REFERENCES qr_batches(id) ON DELETE CASCADE,
    
    -- Assignment tracking (NULL = unassigned)
    equipment_registry_id UUID REFERENCES equipment_registry(id) ON DELETE SET NULL,
    assigned_at TIMESTAMPTZ,
    assigned_by VARCHAR(255),
    
    -- Status lifecycle
    status VARCHAR(50) NOT NULL DEFAULT 'generated',
    -- Lifecycle: generated → reserved → assigned → decommissioned
    
    -- Physical tracking
    serial_number VARCHAR(255),
    printed BOOLEAN DEFAULT false,
    printed_at TIMESTAMPTZ,
    
    -- Metadata
    metadata JSONB DEFAULT '{}',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    
    -- Constraints
    CONSTRAINT chk_qr_code_status CHECK (status IN ('generated', 'reserved', 'assigned', 'decommissioned'))
);

-- Indexes for qr_codes
CREATE UNIQUE INDEX idx_qr_codes_qr_code ON qr_codes(qr_code);
CREATE INDEX idx_qr_codes_status ON qr_codes(status);
CREATE INDEX idx_qr_codes_manufacturer ON qr_codes(manufacturer_id);
CREATE INDEX idx_qr_codes_batch ON qr_codes(batch_id);
CREATE INDEX idx_qr_codes_equipment_registry ON qr_codes(equipment_registry_id);
CREATE INDEX idx_qr_codes_equipment_catalog ON qr_codes(equipment_catalog_id);
CREATE INDEX idx_qr_codes_serial_number ON qr_codes(serial_number);
CREATE INDEX idx_qr_codes_created_at ON qr_codes(created_at DESC);

-- Composite index for common queries
CREATE INDEX idx_qr_codes_status_manufacturer ON qr_codes(status, manufacturer_id);
CREATE INDEX idx_qr_codes_status_batch ON qr_codes(status, batch_id);

-- Comments
COMMENT ON TABLE qr_codes IS 'Individual QR codes with lifecycle management';
COMMENT ON COLUMN qr_codes.status IS 'QR code lifecycle status: generated|reserved|assigned|decommissioned';
COMMENT ON COLUMN qr_codes.equipment_registry_id IS 'NULL = unassigned QR code, UUID = assigned to equipment';
COMMENT ON COLUMN qr_codes.qr_code IS 'Unique QR code identifier (e.g., QR-20251223-000001)';
COMMENT ON COLUMN qr_codes.qr_code_url IS 'Full URL to equipment page (e.g., https://app.com/equipment/qr/<code>)';

-- =====================================================================
-- 3. TRIGGERS
-- =====================================================================

-- Auto-update updated_at timestamp for qr_batches
CREATE OR REPLACE FUNCTION update_qr_batches_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER qr_batches_updated_at_trigger
    BEFORE UPDATE ON qr_batches
    FOR EACH ROW
    EXECUTE FUNCTION update_qr_batches_updated_at();

-- Auto-update updated_at timestamp for qr_codes
CREATE OR REPLACE FUNCTION update_qr_codes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER qr_codes_updated_at_trigger
    BEFORE UPDATE ON qr_codes
    FOR EACH ROW
    EXECUTE FUNCTION update_qr_codes_updated_at();

-- Auto-update batch quantity when QR codes are created
CREATE OR REPLACE FUNCTION increment_batch_quantity()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE qr_batches 
    SET quantity_generated = quantity_generated + 1
    WHERE id = NEW.batch_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER qr_codes_increment_batch_trigger
    AFTER INSERT ON qr_codes
    FOR EACH ROW
    WHEN (NEW.batch_id IS NOT NULL)
    EXECUTE FUNCTION increment_batch_quantity();

-- =====================================================================
-- 4. HELPER FUNCTIONS
-- =====================================================================

-- Function to generate unique QR code
CREATE OR REPLACE FUNCTION generate_unique_qr_code()
RETURNS VARCHAR AS $$
DECLARE
    v_qr_code VARCHAR;
    v_exists BOOLEAN;
BEGIN
    LOOP
        -- Generate format: QR-YYYYMMDD-XXXXXX (random 6 digit)
        v_qr_code := 'QR-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');
        
        -- Check if exists
        SELECT EXISTS(SELECT 1 FROM qr_codes WHERE qr_code = v_qr_code) INTO v_exists;
        
        -- If doesn't exist, use it
        IF NOT v_exists THEN
            EXIT;
        END IF;
    END LOOP;
    
    RETURN v_qr_code;
END;
$$ LANGUAGE plpgsql;

-- Function to generate batch number
CREATE OR REPLACE FUNCTION generate_batch_number()
RETURNS VARCHAR AS $$
DECLARE
    v_batch_number VARCHAR;
    v_exists BOOLEAN;
BEGIN
    LOOP
        -- Generate format: BATCH-YYYYMMDD-XXX (random 3 digit)
        v_batch_number := 'BATCH-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(FLOOR(RANDOM() * 1000)::TEXT, 3, '0');
        
        -- Check if exists
        SELECT EXISTS(SELECT 1 FROM qr_batches WHERE batch_number = v_batch_number) INTO v_exists;
        
        -- If doesn't exist, use it
        IF NOT v_exists THEN
            EXIT;
        END IF;
    END LOOP;
    
    RETURN v_batch_number;
END;
$$ LANGUAGE plpgsql;

-- Function to get unassigned QR codes count for a batch
CREATE OR REPLACE FUNCTION get_unassigned_qr_count(batch_uuid UUID)
RETURNS INT AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM qr_codes
    WHERE batch_id = batch_uuid 
    AND status = 'generated'
    AND equipment_registry_id IS NULL;
    
    RETURN v_count;
END;
$$ LANGUAGE plpgsql;

-- Function to get batch statistics
CREATE OR REPLACE FUNCTION get_batch_stats(batch_uuid UUID)
RETURNS TABLE (
    total_generated INT,
    assigned INT,
    unassigned INT,
    reserved INT,
    decommissioned INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*)::INT as total_generated,
        COUNT(*) FILTER (WHERE status = 'assigned')::INT as assigned,
        COUNT(*) FILTER (WHERE status = 'generated' AND equipment_registry_id IS NULL)::INT as unassigned,
        COUNT(*) FILTER (WHERE status = 'reserved')::INT as reserved,
        COUNT(*) FILTER (WHERE status = 'decommissioned')::INT as decommissioned
    FROM qr_codes
    WHERE batch_id = batch_uuid;
END;
$$ LANGUAGE plpgsql;

-- =====================================================================
-- 5. VIEWS FOR COMMON QUERIES
-- =====================================================================

-- View for unassigned QR codes (available for assignment)
CREATE OR REPLACE VIEW qr_codes_unassigned AS
SELECT 
    qc.id,
    qc.qr_code,
    qc.qr_code_url,
    qc.serial_number,
    qc.manufacturer_id,
    qc.equipment_catalog_id,
    qc.batch_id,
    qc.printed,
    qc.created_at,
    o.name as manufacturer_name,
    ec.model_name as equipment_model,
    qb.batch_number
FROM qr_codes qc
LEFT JOIN organizations o ON qc.manufacturer_id = o.id
LEFT JOIN equipment_catalog ec ON qc.equipment_catalog_id = ec.id
LEFT JOIN qr_batches qb ON qc.batch_id = qb.id
WHERE qc.status = 'generated' 
AND qc.equipment_registry_id IS NULL;

COMMENT ON VIEW qr_codes_unassigned IS 'Unassigned QR codes available for assignment';

-- View for batch summary with statistics
CREATE OR REPLACE VIEW qr_batches_summary AS
SELECT 
    qb.id,
    qb.batch_number,
    qb.manufacturer_id,
    qb.equipment_catalog_id,
    qb.quantity_requested,
    qb.quantity_generated,
    qb.status,
    qb.pdf_url,
    qb.csv_url,
    qb.created_at,
    qb.generated_at,
    o.name as manufacturer_name,
    ec.model_name as equipment_model,
    COUNT(qc.id) FILTER (WHERE qc.status = 'assigned') as assigned_count,
    COUNT(qc.id) FILTER (WHERE qc.status = 'generated' AND qc.equipment_registry_id IS NULL) as unassigned_count,
    COUNT(qc.id) FILTER (WHERE qc.printed = true) as printed_count
FROM qr_batches qb
LEFT JOIN organizations o ON qb.manufacturer_id = o.id
LEFT JOIN equipment_catalog ec ON qb.equipment_catalog_id = ec.id
LEFT JOIN qr_codes qc ON qb.id = qc.batch_id
GROUP BY qb.id, qb.batch_number, qb.manufacturer_id, qb.equipment_catalog_id,
         qb.quantity_requested, qb.quantity_generated, qb.status, qb.pdf_url,
         qb.csv_url, qb.created_at, qb.generated_at, o.name, ec.model_name;

COMMENT ON VIEW qr_batches_summary IS 'QR batches with statistics and manufacturer info';

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

-- Log migration
DO $$
BEGIN
    RAISE NOTICE 'Migration 028: QR code tables created successfully';
    RAISE NOTICE '  - qr_batches table: Batch tracking';
    RAISE NOTICE '  - qr_codes table: Individual QR lifecycle';
    RAISE NOTICE '  - 4 helper functions created';
    RAISE NOTICE '  - 2 views created for common queries';
    RAISE NOTICE '  - 3 triggers for auto-updates';
END $$;
