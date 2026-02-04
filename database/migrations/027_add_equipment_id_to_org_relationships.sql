-- ============================================================================
-- Migration: Add equipment_id to org_relationships for Partner Association
-- ============================================================================
-- Date: February 4, 2026
-- Description: Enables hybrid partner association (org-level + equipment-level)
--              NULL equipment_id = partner services ALL equipment (general)
--              Specific equipment_id = partner services only that equipment (override)
-- ============================================================================

BEGIN;

-- Step 1: Add equipment_id column (nullable)
-- Note: equipment.id is VARCHAR, not UUID
ALTER TABLE org_relationships
ADD COLUMN equipment_id VARCHAR NULL;

-- Step 1b: Add foreign key constraint
ALTER TABLE org_relationships
ADD CONSTRAINT org_relationships_equipment_id_fkey
FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE;

-- Step 2: Create index for equipment-based queries (partial index for performance)
CREATE INDEX idx_org_rel_equipment 
ON org_relationships(equipment_id)
WHERE equipment_id IS NOT NULL;

-- Step 3: Create index for parent org queries (if not exists)
CREATE INDEX IF NOT EXISTS idx_org_rel_parent 
ON org_relationships(parent_org_id);

-- Step 4: Create composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_org_rel_parent_child 
ON org_relationships(parent_org_id, child_org_id);

-- Step 5: Update unique constraint to include equipment_id
-- Drop existing unique constraint if it exists
ALTER TABLE org_relationships 
DROP CONSTRAINT IF EXISTS org_relationships_parent_org_id_child_org_id_rel_type_key;

ALTER TABLE org_relationships
DROP CONSTRAINT IF EXISTS org_relationships_unique;

-- Add new unique constraint with equipment_id
-- Using unique index with COALESCE to treat NULL as a specific value
-- This allows:
--   - One general association (equipment_id = NULL) per manufacturer-partner pair
--   - Multiple equipment-specific associations for same manufacturer-partner pair
CREATE UNIQUE INDEX org_relationships_unique_with_equipment 
ON org_relationships (parent_org_id, child_org_id, rel_type, COALESCE(equipment_id, ''));

-- Step 6: Add column comment for documentation
COMMENT ON COLUMN org_relationships.equipment_id IS 
'Equipment-specific association. NULL = partner services ALL equipment (general). UUID = partner services only this equipment (override).';

-- Step 7: Add helpful view for querying partner associations
CREATE OR REPLACE VIEW partner_associations AS
SELECT 
    r.id,
    r.parent_org_id,
    p_org.name as parent_org_name,
    p_org.org_type as parent_org_type,
    r.child_org_id,
    c_org.name as child_org_name,
    c_org.org_type as child_org_type,
    r.rel_type,
    r.equipment_id,
    CASE 
        WHEN r.equipment_id IS NULL THEN 'general'
        ELSE 'equipment-specific'
    END as association_type,
    e.equipment_name,
    e.serial_number as equipment_serial,
    r.created_at
FROM org_relationships r
JOIN organizations p_org ON p_org.id = r.parent_org_id
JOIN organizations c_org ON c_org.id = r.child_org_id
LEFT JOIN equipment e ON e.id = r.equipment_id
WHERE r.rel_type IN ('services_for', 'partner_of');

COMMENT ON VIEW partner_associations IS 
'Convenient view for querying partner associations with org names and equipment details';

COMMIT;

-- ============================================================================
-- Verification Queries
-- ============================================================================

-- Check if column was added
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'org_relationships' AND column_name = 'equipment_id';

-- Check indexes
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'org_relationships'
ORDER BY indexname;

-- Check constraints
SELECT 
    conname as constraint_name,
    contype as constraint_type,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conrelid = 'org_relationships'::regclass
ORDER BY conname;

-- Show sample data structure
SELECT 
    parent_org_name,
    child_org_name,
    child_org_type,
    association_type,
    equipment_name
FROM partner_associations
LIMIT 5;

-- ============================================================================
-- Example Usage
-- ============================================================================

-- Example 1: Create general partner association (services all equipment)
-- INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, equipment_id)
-- VALUES ('manufacturer-uuid', 'channel-partner-uuid', 'services_for', NULL);

-- Example 2: Create equipment-specific association (override)
-- INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, equipment_id)
-- VALUES ('manufacturer-uuid', 'subdeal-uuid', 'services_for', 'equipment-uuid');

-- ============================================================================
