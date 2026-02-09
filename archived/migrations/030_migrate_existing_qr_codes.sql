-- Migration: 030-migrate-existing-qr-codes.sql
-- Description: Migrate existing QR codes from equipment_registry to qr_codes table
-- Purpose: Populate new qr_codes table with existing data (one-time migration)
-- Date: 2025-12-23
-- IMPORTANT: This migration is safe to run multiple times (idempotent)

-- =====================================================================
-- 1. MIGRATE EXISTING QR CODES TO QR_CODES TABLE
-- =====================================================================

-- Insert existing QR codes from equipment_registry into qr_codes
-- Mark all as 'assigned' since they're already linked to equipment
INSERT INTO qr_codes (
    qr_code,
    qr_code_url,
    manufacturer_id,
    equipment_registry_id,
    serial_number,
    status,
    assigned_at,
    created_at,
    created_by
)
SELECT 
    er.qr_code,
    er.qr_code_url,
    er.manufacturer_id,
    er.id as equipment_registry_id,
    er.serial_number,
    'assigned' as status,  -- Already assigned to equipment
    er.created_at as assigned_at,
    er.created_at,
    er.created_by
FROM equipment_registry er
WHERE er.qr_code IS NOT NULL
AND er.qr_code != ''
AND NOT EXISTS (
    -- Don't insert duplicates if migration already ran
    SELECT 1 FROM qr_codes qc 
    WHERE qc.qr_code = er.qr_code
);

-- =====================================================================
-- 2. UPDATE EQUIPMENT_REGISTRY WITH QR_CODE_ID
-- =====================================================================

-- Link equipment_registry back to qr_codes table
UPDATE equipment_registry er
SET qr_code_id = qc.id
FROM qr_codes qc
WHERE er.qr_code = qc.qr_code
AND er.qr_code_id IS NULL;

-- =====================================================================
-- 3. CREATE A MIGRATION BATCH FOR EXISTING QR CODES
-- =====================================================================

-- Create a special batch to represent "migrated" QR codes
INSERT INTO qr_batches (
    batch_number,
    manufacturer_id,
    quantity_requested,
    quantity_generated,
    status,
    generated_at,
    generated_by,
    notes
)
SELECT 
    'BATCH-MIGRATED-' || TO_CHAR(NOW(), 'YYYYMMDD'),
    NULL,  -- No specific manufacturer (mixed batch)
    COUNT(*),
    COUNT(*),
    'completed',
    NOW(),
    'migration-script',
    'Migrated from existing equipment_registry QR codes'
FROM qr_codes
WHERE batch_id IS NULL
AND status = 'assigned'
ON CONFLICT DO NOTHING;

-- Link migrated QR codes to the migration batch
UPDATE qr_codes qc
SET batch_id = qb.id
FROM qr_batches qb
WHERE qc.batch_id IS NULL
AND qc.status = 'assigned'
AND qb.batch_number LIKE 'BATCH-MIGRATED-%'
AND qb.generated_by = 'migration-script';

-- =====================================================================
-- 4. VALIDATION QUERIES (INFORMATIONAL)
-- =====================================================================

-- Count records in each table
DO $$
DECLARE
    v_equipment_count INT;
    v_qr_count INT;
    v_linked_count INT;
    v_unlinked_count INT;
BEGIN
    -- Total equipment with QR codes
    SELECT COUNT(*) INTO v_equipment_count
    FROM equipment_registry
    WHERE qr_code IS NOT NULL AND qr_code != '';
    
    -- Total QR codes in new table
    SELECT COUNT(*) INTO v_qr_count
    FROM qr_codes;
    
    -- Equipment with qr_code_id linked
    SELECT COUNT(*) INTO v_linked_count
    FROM equipment_registry
    WHERE qr_code_id IS NOT NULL;
    
    -- Equipment with QR but not linked
    SELECT COUNT(*) INTO v_unlinked_count
    FROM equipment_registry
    WHERE qr_code IS NOT NULL 
    AND qr_code != ''
    AND qr_code_id IS NULL;
    
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Migration 030: QR Codes Migration Complete';
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Equipment with QR codes: %', v_equipment_count;
    RAISE NOTICE 'QR codes migrated: %', v_qr_count;
    RAISE NOTICE 'Equipment linked to qr_codes: %', v_linked_count;
    RAISE NOTICE 'Equipment unlinked: %', v_unlinked_count;
    
    IF v_unlinked_count > 0 THEN
        RAISE WARNING 'Found % equipment with QR codes but not linked!', v_unlinked_count;
        RAISE NOTICE 'Run this query to investigate:';
        RAISE NOTICE 'SELECT id, qr_code, equipment_name FROM equipment_registry WHERE qr_code IS NOT NULL AND qr_code_id IS NULL;';
    ELSE
        RAISE NOTICE '✅ All equipment QR codes successfully migrated and linked!';
    END IF;
    
    RAISE NOTICE '========================================';
END $$;

-- =====================================================================
-- 5. CREATE HELPER VIEW FOR MIGRATION STATUS
-- =====================================================================

CREATE OR REPLACE VIEW migration_qr_status AS
SELECT 
    'Total Equipment' as metric,
    COUNT(*) as count
FROM equipment_registry
WHERE qr_code IS NOT NULL AND qr_code != ''

UNION ALL

SELECT 
    'QR Codes in qr_codes table',
    COUNT(*)
FROM qr_codes

UNION ALL

SELECT 
    'Equipment Linked to qr_codes',
    COUNT(*)
FROM equipment_registry
WHERE qr_code_id IS NOT NULL

UNION ALL

SELECT 
    'Assigned QR Codes',
    COUNT(*)
FROM qr_codes
WHERE status = 'assigned'

UNION ALL

SELECT 
    'Unassigned QR Codes',
    COUNT(*)
FROM qr_codes
WHERE status = 'generated' AND equipment_registry_id IS NULL;

COMMENT ON VIEW migration_qr_status IS 'Migration status for QR codes';

-- =====================================================================
-- 6. POST-MIGRATION VALIDATION
-- =====================================================================

-- Ensure referential integrity
DO $$
DECLARE
    v_orphaned_qr INT;
    v_orphaned_equipment INT;
BEGIN
    -- Check for QR codes pointing to non-existent equipment
    SELECT COUNT(*) INTO v_orphaned_qr
    FROM qr_codes qc
    WHERE qc.equipment_registry_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM equipment_registry er 
        WHERE er.id = qc.equipment_registry_id
    );
    
    -- Check for equipment pointing to non-existent QR codes
    SELECT COUNT(*) INTO v_orphaned_equipment
    FROM equipment_registry er
    WHERE er.qr_code_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM qr_codes qc 
        WHERE qc.id = er.qr_code_id
    );
    
    IF v_orphaned_qr > 0 OR v_orphaned_equipment > 0 THEN
        RAISE WARNING 'Data integrity issues found:';
        IF v_orphaned_qr > 0 THEN
            RAISE WARNING '  - % QR codes point to non-existent equipment', v_orphaned_qr;
        END IF;
        IF v_orphaned_equipment > 0 THEN
            RAISE WARNING '  - % equipment records point to non-existent QR codes', v_orphaned_equipment;
        END IF;
    ELSE
        RAISE NOTICE '✅ Data integrity check passed!';
    END IF;
END $$;

-- =====================================================================
-- MIGRATION NOTES
-- =====================================================================

-- This migration is safe to run multiple times (idempotent)
-- It will not create duplicate records
-- 
-- If you need to re-run this migration:
-- 1. It will skip existing QR codes (INSERT ... WHERE NOT EXISTS)
-- 2. It will only link unlinked equipment (UPDATE ... WHERE qr_code_id IS NULL)
-- 
-- To view migration status:
-- SELECT * FROM migration_qr_status;
-- 
-- To view migrated QR codes:
-- SELECT * FROM qr_codes WHERE status = 'assigned';
-- 
-- To view unassigned QR codes:
-- SELECT * FROM qr_codes_unassigned;
