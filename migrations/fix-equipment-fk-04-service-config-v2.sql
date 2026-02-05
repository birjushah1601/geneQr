-- Migration: Fix equipment_service_config foreign key (Version 2)
-- Date: 2026-02-05
-- Purpose: Point to equipment_registry instead of equipment
-- Impact: Orphaned records deleted, now applying FK constraint

BEGIN;

-- Step 1: Verify no records reference non-existent equipment_registry
DO $$
DECLARE
    orphaned_count INT;
    total_count INT;
BEGIN
    SELECT COUNT(*) INTO total_count FROM equipment_service_config;
    
    SELECT COUNT(*) INTO orphaned_count
    FROM equipment_service_config sc
    LEFT JOIN equipment_registry er ON sc.equipment_id = er.id
    WHERE sc.equipment_id IS NOT NULL AND er.id IS NULL;
    
    IF orphaned_count > 0 THEN
        RAISE EXCEPTION 'Found % orphaned records out of % total in equipment_service_config. Must fix data first.', 
            orphaned_count, total_count;
    END IF;
    
    RAISE NOTICE 'Pre-check passed: % total records, 0 orphaned records.', total_count;
END $$;

-- Step 2: Drop old FK constraint
ALTER TABLE equipment_service_config 
DROP CONSTRAINT IF EXISTS equipment_service_config_equipment_id_fkey;

-- Step 3: Add new FK constraint to equipment_registry
ALTER TABLE equipment_service_config
ADD CONSTRAINT equipment_service_config_equipment_id_fkey
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id)
ON DELETE CASCADE;

-- Step 4: Final verification
DO $$
DECLARE
    total_count INT;
BEGIN
    SELECT COUNT(*) INTO total_count FROM equipment_service_config;
    
    RAISE NOTICE '✓ equipment_service_config FK migration successful. % records, all valid.', total_count;
END $$;

COMMIT;
