-- Migration: Fix equipment_usage_logs foreign key
-- Date: 2026-02-05
-- Purpose: Point to equipment_registry instead of equipment
-- Impact: 0 records exist, so no data migration needed

BEGIN;

-- Step 1: Drop old FK constraint
ALTER TABLE equipment_usage_logs 
DROP CONSTRAINT IF EXISTS equipment_usage_logs_equipment_id_fkey;

-- Step 2: Add new FK constraint to equipment_registry
ALTER TABLE equipment_usage_logs
ADD CONSTRAINT equipment_usage_logs_equipment_id_fkey
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id)
ON DELETE CASCADE;

-- Step 3: Verify no orphaned records
DO $$
DECLARE
    orphaned_count INT;
BEGIN
    SELECT COUNT(*) INTO orphaned_count
    FROM equipment_usage_logs ul
    LEFT JOIN equipment_registry er ON ul.equipment_id = er.id
    WHERE ul.equipment_id IS NOT NULL AND er.id IS NULL;
    
    IF orphaned_count > 0 THEN
        RAISE EXCEPTION 'Found % orphaned records in equipment_usage_logs', orphaned_count;
    END IF;
    
    RAISE NOTICE 'equipment_usage_logs FK migration successful. 0 orphaned records.';
END $$;

COMMIT;
