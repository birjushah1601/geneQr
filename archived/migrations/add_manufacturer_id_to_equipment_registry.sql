-- ============================================================================
-- Add manufacturer_id to equipment_registry and populate it
-- ============================================================================
-- This migration adds manufacturer_id to equipment_registry and links it
-- through equipment_catalog for proper manufacturer dashboard support

-- ============================================================================
-- 1. ADD manufacturer_id COLUMN
-- ============================================================================

-- Add the column (nullable initially)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS manufacturer_id uuid;

-- Add index for performance
CREATE INDEX IF NOT EXISTS idx_equipment_registry_manufacturer 
ON equipment_registry(manufacturer_id);

-- Add foreign key constraint
ALTER TABLE equipment_registry
ADD CONSTRAINT fk_equipment_registry_manufacturer
FOREIGN KEY (manufacturer_id) 
REFERENCES organizations(id)
ON DELETE SET NULL;

-- ============================================================================
-- 2. POPULATE manufacturer_id FROM equipment_catalog
-- ============================================================================

-- Update equipment_registry.manufacturer_id from equipment_catalog.manufacturer_id
-- This links registry items to manufacturers through the catalog relationship
UPDATE equipment_registry er
SET manufacturer_id = ec.manufacturer_id
FROM equipment_catalog ec
WHERE er.equipment_catalog_id = ec.id
  AND ec.manufacturer_id IS NOT NULL;

-- ============================================================================
-- 3. VERIFY THE LINKING
-- ============================================================================

-- Show counts before and after
SELECT 
    'Equipment Registry Items' as metric,
    COUNT(*) as total,
    COUNT(manufacturer_id) as with_manufacturer_id,
    COUNT(*) - COUNT(manufacturer_id) as missing_manufacturer_id
FROM equipment_registry;

-- Show distribution by manufacturer
SELECT 
    o.name as manufacturer,
    COUNT(er.id) as installed_units
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY installed_units DESC, o.name;

-- ============================================================================
-- 4. UPDATE REMAINING ITEMS BY NAME MATCHING
-- ============================================================================

-- For items that didn't get linked via catalog, try name matching
-- This handles legacy data where equipment_catalog_id might be null

UPDATE equipment_registry er
SET manufacturer_id = o.id
FROM organizations o
WHERE er.manufacturer_id IS NULL
  AND o.org_type = 'manufacturer'
  AND (
    -- Direct name match
    er.manufacturer_name = o.name
    OR
    -- Partial match (first word)
    er.manufacturer_name ILIKE '%' || SPLIT_PART(o.name, ' ', 1) || '%'
    OR
    -- Handle "Siemens Healthineers" vs "Siemens Healthineers India"
    SPLIT_PART(er.manufacturer_name, ' ', 1) = SPLIT_PART(o.name, ' ', 1)
  );

-- ============================================================================
-- 5. FINAL VERIFICATION
-- ============================================================================

-- Show final counts
SELECT 
    'FINAL STATUS' as stage,
    COUNT(*) as total_items,
    COUNT(manufacturer_id) as linked_items,
    COUNT(*) - COUNT(manufacturer_id) as unlinked_items,
    ROUND(100.0 * COUNT(manufacturer_id) / COUNT(*), 1) as percent_linked
FROM equipment_registry;

-- Show detailed manufacturer mapping
SELECT 
    o.name as manufacturer_name,
    o.id as manufacturer_id,
    COUNT(er.id) as installed_equipment,
    COUNT(DISTINCT er.customer_id) as unique_customers,
    array_agg(DISTINCT er.category) FILTER (WHERE er.category IS NOT NULL) as equipment_types
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY installed_equipment DESC, o.name;

-- Show unlinked items (if any)
SELECT 
    id,
    equipment_name,
    manufacturer_name,
    'No matching manufacturer org' as reason
FROM equipment_registry
WHERE manufacturer_id IS NULL
ORDER BY manufacturer_name;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

SELECT '✅ manufacturer_id added and linked to equipment_registry!' as result;
