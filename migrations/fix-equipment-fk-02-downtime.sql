-- Migration: Fix equipment_downtime foreign key
-- Date: 2026-02-05
-- Purpose: Point to equipment_registry instead of equipment
-- Impact: 0 records exist, so no data migration needed

BEGIN;

-- Step 1: Drop old FK constraint
ALTER TABLE equipment_downtime 
DROP CONSTRAINT IF EXISTS equipment_downtime_equipment_id_fkey;

-- Step 2: Add new FK constraint to equipment_registry
ALTER TABLE equipment_downtime
ADD CONSTRAINT equipment_downtime_equipment_id_fkey
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id)
ON DELETE CASCADE;

-- Step 3: Verify no orphaned records
DO $$
DECLARE
    orphaned_count INT;
BEGIN
    SELECT COUNT(*) INTO orphaned_count
    FROM equipment_downtime ed
    LEFT JOIN equipment_registry er ON ed.equipment_id = er.id
    WHERE ed.equipment_id IS NOT NULL AND er.id IS NULL;
    
    IF orphaned_count > 0 THEN
        RAISE EXCEPTION 'Found % orphaned records in equipment_downtime', orphaned_count;
    END IF;
    
    RAISE NOTICE 'equipment_downtime FK migration successful. 0 orphaned records.';
END $$;

COMMIT;
