-- Migration: Fix spare_parts_catalog FK to use equipment_registry
-- Date: 2026-02-05
-- Issue: spare_parts_catalog references equipment (catalog) instead of equipment_registry (installations)
-- Impact: Parts can't be properly linked to installed equipment

-- Step 1: Check current constraint
SELECT 
    conname as constraint_name,
    conrelid::regclass as table_name,
    confrelid::regclass as referenced_table,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint 
WHERE conrelid = 'spare_parts_catalog'::regclass 
  AND contype = 'f'
  AND confrelid = 'equipment'::regclass;

-- Step 2: Check for orphaned records
SELECT 
    COUNT(*) as orphaned_records,
    'Records in spare_parts_catalog referencing non-existent equipment_registry' as description
FROM spare_parts_catalog spc
WHERE spc.equipment_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM equipment_registry er 
    WHERE er.id = spc.equipment_id
  );

-- Step 3: List orphaned records (if any)
SELECT 
    spc.id,
    spc.part_number,
    spc.part_name,
    spc.equipment_id,
    'Orphaned - no matching equipment_registry' as status
FROM spare_parts_catalog spc
WHERE spc.equipment_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM equipment_registry er 
    WHERE er.id = spc.equipment_id
  );

-- Step 4: Backup orphaned records (if any exist)
-- Run this manually if orphaned records found:
/*
CREATE TEMP TABLE backup_orphaned_spare_parts AS
SELECT * FROM spare_parts_catalog
WHERE equipment_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM equipment_registry er 
    WHERE er.id = equipment_id
  );
*/

-- Step 5: Handle orphaned records
-- Option A: Set equipment_id to NULL for orphaned records
UPDATE spare_parts_catalog
SET equipment_id = NULL
WHERE equipment_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM equipment_registry er 
    WHERE er.id = equipment_id
  );

-- Step 6: Drop old constraint
ALTER TABLE spare_parts_catalog
DROP CONSTRAINT IF EXISTS spare_parts_catalog_equipment_id_fkey;

ALTER TABLE spare_parts_catalog
DROP CONSTRAINT IF EXISTS fk_spare_parts_equipment;

ALTER TABLE spare_parts_catalog
DROP CONSTRAINT IF EXISTS fk_equipment;

-- Step 7: Add new constraint to equipment_registry
ALTER TABLE spare_parts_catalog
ADD CONSTRAINT spare_parts_catalog_equipment_id_fkey 
FOREIGN KEY (equipment_id) 
REFERENCES equipment_registry(id) 
ON DELETE SET NULL;

-- Step 8: Verify the new constraint
SELECT 
    conname as constraint_name,
    conrelid::regclass as table_name,
    confrelid::regclass as referenced_table,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint 
WHERE conrelid = 'spare_parts_catalog'::regclass 
  AND contype = 'f'
  AND conname = 'spare_parts_catalog_equipment_id_fkey';

-- Step 9: Test the constraint
-- Try to insert invalid equipment_id (should fail)
-- DO $$
-- BEGIN
--     INSERT INTO spare_parts_catalog (
--         part_number, part_name, equipment_id
--     ) VALUES (
--         'TEST-PART', 'Test Part', '00000000-0000-0000-0000-000000000000'
--     );
-- EXCEPTION WHEN foreign_key_violation THEN
--     RAISE NOTICE 'Foreign key constraint working correctly';
-- END $$;

-- Step 10: Summary
SELECT 
    'spare_parts_catalog FK constraint updated successfully' as status,
    COUNT(*) as total_records,
    COUNT(equipment_id) as records_with_equipment_link
FROM spare_parts_catalog;

-- Notes:
-- 1. This migration updates spare_parts_catalog to reference equipment_registry
-- 2. Orphaned records (if any) have equipment_id set to NULL
-- 3. The constraint now uses ON DELETE SET NULL for safety
-- 4. Parts catalog should ideally reference equipment (catalog) for compatibility
-- 5. But if equipment_id is meant to track which installation uses the part, 
--    then equipment_registry is correct

-- Rollback (if needed):
/*
ALTER TABLE spare_parts_catalog
DROP CONSTRAINT spare_parts_catalog_equipment_id_fkey;

ALTER TABLE spare_parts_catalog
ADD CONSTRAINT spare_parts_catalog_equipment_id_fkey 
FOREIGN KEY (equipment_id) 
REFERENCES equipment(id) 
ON DELETE CASCADE;
*/
