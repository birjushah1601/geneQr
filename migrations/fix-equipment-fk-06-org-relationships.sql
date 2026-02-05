-- Migration: Fix org_relationships foreign key
-- Date: 2026-02-05
-- Purpose: Point to equipment_registry instead of equipment
-- Impact: 0 records with equipment_id, so no data migration needed

BEGIN;

-- Step 1: Drop old FK constraint
ALTER TABLE org_relationships 
DROP CONSTRAINT IF EXISTS org_relationships_equipment_id_fkey;

-- Step 2: Add new FK constraint to equipment_registry
ALTER TABLE org_relationships
ADD CONSTRAINT org_relationships_equipment_id_fkey
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id)
ON DELETE CASCADE;

-- Step 3: Verify no orphaned records
DO $$
DECLARE
    orphaned_count INT;
    total_with_equipment INT;
BEGIN
    SELECT COUNT(*) INTO total_with_equipment 
    FROM org_relationships 
    WHERE equipment_id IS NOT NULL;
    
    SELECT COUNT(*) INTO orphaned_count
    FROM org_relationships r
    LEFT JOIN equipment_registry er ON r.equipment_id = er.id
    WHERE r.equipment_id IS NOT NULL AND er.id IS NULL;
    
    IF orphaned_count > 0 THEN
        RAISE EXCEPTION 'Found % orphaned records out of % total in org_relationships', 
            orphaned_count, total_with_equipment;
    END IF;
    
    RAISE NOTICE 'org_relationships FK migration successful. % records with equipment_id, 0 orphaned.', 
        total_with_equipment;
END $$;

COMMIT;
