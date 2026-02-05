-- Migration: Fix equipment_service_config foreign key
-- Date: 2026-02-05
-- Purpose: Point to equipment_registry instead of equipment
-- Impact: 3 records exist - need data migration

BEGIN;

-- Step 1: Check which equipment IDs exist in both tables (for migration)
DO $$
DECLARE
    overlap_count INT;
BEGIN
    SELECT COUNT(*) INTO overlap_count
    FROM equipment e
    INNER JOIN equipment_registry er ON e.id = er.id;
    
    RAISE NOTICE 'Found % equipment IDs that exist in both tables', overlap_count;
END $$;

-- Step 2: Migrate data from equipment IDs to equipment_registry IDs
-- For the 2 overlapping equipment, IDs are already the same, so no change needed
-- Just verify they exist in equipment_registry

-- Update equipment_id to point to equipment_registry
-- (For overlapping equipment, IDs are already the same, so this is safe)
UPDATE equipment_service_config sc
SET equipment_id = er.id
FROM equipment e
INNER JOIN equipment_registry er ON e.id = er.id
WHERE sc.equipment_id = e.id;

-- Step 3: Drop old FK constraint
ALTER TABLE equipment_service_config 
DROP CONSTRAINT IF EXISTS equipment_service_config_equipment_id_fkey;

-- Step 4: Add new FK constraint to equipment_registry
ALTER TABLE equipment_service_config
ADD CONSTRAINT equipment_service_config_equipment_id_fkey
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id)
ON DELETE CASCADE;

-- Step 5: Verify no orphaned records
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
        RAISE EXCEPTION 'Found % orphaned records out of % total in equipment_service_config', 
            orphaned_count, total_count;
    END IF;
    
    RAISE NOTICE 'equipment_service_config FK migration successful. % records, 0 orphaned.', total_count;
END $$;

COMMIT;
