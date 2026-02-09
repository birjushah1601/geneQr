-- ============================================================================
-- Link Equipment Records to Equipment Catalog
-- This allows parts to be displayed for tickets
-- ============================================================================

-- Update equipment records to link them to their catalog entries
UPDATE equipment e
SET catalog_id = ec.id
FROM equipment_catalog ec
WHERE LOWER(TRIM(e.equipment_name)) = LOWER(TRIM(ec.product_name))
  AND e.catalog_id IS NULL;

-- Show results
SELECT 
    'Equipment linked to catalog' as message,
    COUNT(*) as linked_count
FROM equipment
WHERE catalog_id IS NOT NULL;

-- Show which equipment now has parts
SELECT 
    e.id as equipment_id,
    e.equipment_name,
    ec.product_name as catalog_name,
    COUNT(epa.id) as parts_count
FROM equipment e
JOIN equipment_catalog ec ON e.catalog_id = ec.id
LEFT JOIN equipment_part_assignments epa ON epa.equipment_catalog_id = ec.id
GROUP BY e.id, e.equipment_name, ec.product_name
ORDER BY parts_count DESC;
