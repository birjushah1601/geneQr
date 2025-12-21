-- ============================================================================
-- Clear All Parts Assigned to Tickets
-- ============================================================================
-- This removes all parts assignments from tickets so they can be tested fresh
-- Parts remain in the catalog and can be assigned again

BEGIN;

-- Show current state before clearing
SELECT 
    '=== BEFORE: Tickets with Parts Assigned ===' as status;

SELECT 
    ticket_number,
    equipment_name,
    jsonb_array_length(COALESCE(parts_used, '[]'::jsonb)) as parts_count,
    parts_used
FROM service_tickets
WHERE parts_used IS NOT NULL 
  AND parts_used != '[]'::jsonb
ORDER BY ticket_number;

-- Clear all parts assignments
UPDATE service_tickets
SET parts_used = '[]'::jsonb
WHERE parts_used IS NOT NULL 
  AND parts_used != '[]'::jsonb
RETURNING ticket_number, equipment_name;

-- Show final state
SELECT 
    '=== AFTER: Verification ===' as status;

SELECT 
    COUNT(*) as total_tickets,
    COUNT(CASE WHEN parts_used IS NULL OR parts_used = '[]'::jsonb THEN 1 END) as tickets_with_no_parts,
    COUNT(CASE WHEN parts_used IS NOT NULL AND parts_used != '[]'::jsonb THEN 1 END) as tickets_with_parts
FROM service_tickets;

-- Verify parts are still available in catalog
SELECT 
    '=== Parts Still Available in Catalog ===' as status;

SELECT 
    COUNT(*) as total_parts,
    COUNT(CASE WHEN is_available = true THEN 1 END) as available_parts
FROM spare_parts_catalog;

SELECT 
    '=== Equipment Part Assignments Still Intact ===' as status;

SELECT 
    COUNT(*) as total_assignments,
    COUNT(DISTINCT equipment_catalog_id) as equipment_types_with_parts,
    COUNT(DISTINCT spare_part_id) as unique_parts_assigned
FROM equipment_part_assignments;

COMMIT;

SELECT '✅ All parts assignments cleared from tickets. Parts remain available for fresh assignment.' as final_status;
