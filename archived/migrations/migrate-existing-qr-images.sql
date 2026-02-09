-- Migration: Copy existing QR images from equipment table to equipment_registry
-- Date: 2026-02-05
-- Author: Droid AI
-- Purpose: Preserve existing QR images for equipment records that exist in both tables
-- Context: 10 equipment records exist in both equipment and equipment_registry tables

BEGIN;

-- Migrate QR images from equipment to equipment_registry
UPDATE equipment_registry er
SET 
    qr_code_image = e.qr_code_image,
    qr_code_format = e.qr_code_format,
    qr_code_generated_at = e.qr_code_generated_at
FROM equipment e
WHERE er.id = e.id
AND e.qr_code_image IS NOT NULL;

-- Log migration results
DO $$
DECLARE
    migrated_count INT;
    total_count INT;
BEGIN
    -- Count migrated records
    SELECT COUNT(*) INTO migrated_count
    FROM equipment_registry
    WHERE qr_code_image IS NOT NULL;
    
    -- Count total records
    SELECT COUNT(*) INTO total_count
    FROM equipment_registry;
    
    RAISE NOTICE 'Migration complete: % of % equipment records now have QR images', 
                 migrated_count, total_count;
    RAISE NOTICE 'Remaining % equipment records need QR generation', 
                 (total_count - migrated_count);
END $$;

COMMIT;

-- Verification query (run separately)
-- SELECT 
--     COUNT(*) as total_equipment,
--     COUNT(qr_code_image) FILTER (WHERE qr_code_image IS NOT NULL) as with_qr_image,
--     COUNT(qr_code_image) FILTER (WHERE qr_code_image IS NULL) as without_qr_image,
--     ROUND(100.0 * COUNT(qr_code_image) FILTER (WHERE qr_code_image IS NOT NULL) / COUNT(*), 2) as percentage_with_qr
-- FROM equipment_registry;
