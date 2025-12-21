-- ============================================================================
-- Fix get_parts_for_registry function to use equipment_part_assignments table
-- ============================================================================

-- Drop and recreate the function to use the correct table
CREATE OR REPLACE FUNCTION public.get_parts_for_registry(p_registry_id character varying)
RETURNS TABLE(
  spare_part_id uuid,
  part_number text,
  part_name text,
  unit_price numeric,
  currency text,
  is_critical boolean,
  quantity_required integer,
  part_category text,
  stock_status text,
  lead_time_days integer
)
LANGUAGE plpgsql
STABLE
AS $function$
DECLARE
  v_catalog_id uuid;
BEGIN
  -- First, try to find the catalog_id from equipment_registry
  SELECT equipment_catalog_id INTO v_catalog_id
  FROM equipment_registry
  WHERE id = p_registry_id
  LIMIT 1;
  
  -- If not found, try equipment table
  IF v_catalog_id IS NULL THEN
    SELECT catalog_id INTO v_catalog_id
    FROM equipment
    WHERE id = p_registry_id
    LIMIT 1;
  END IF;
  
  -- If we found a catalog_id, get parts via equipment_part_assignments
  IF v_catalog_id IS NOT NULL THEN
    RETURN QUERY
    SELECT
      sp.id AS spare_part_id,
      sp.part_number::text,
      sp.part_name::text,
      sp.unit_price,
      sp.currency::text,
      epa.is_critical,
      COALESCE(epa.quantity_required, 1) AS quantity_required,
      sp.category::text AS part_category,
      sp.stock_status::text,
      sp.lead_time_days
    FROM equipment_part_assignments epa
    JOIN spare_parts_catalog sp ON sp.id = epa.spare_part_id AND sp.is_available = true
    WHERE epa.equipment_catalog_id = v_catalog_id
    ORDER BY epa.is_critical DESC, COALESCE(epa.quantity_required, 1) DESC, sp.part_name;
    
    RETURN;
  END IF;
  
  -- If no results, fall back to equipment_spare_parts (legacy)
  IF NOT FOUND THEN
    RETURN QUERY
    SELECT
      sp.id AS spare_part_id,
      sp.part_number::text,
      sp.part_name::text,
      sp.unit_price,
      sp.currency::text,
      esp.is_critical,
      COALESCE(esp.quantity_required, 1) AS quantity_required,
      sp.category::text AS part_category,
      sp.stock_status::text,
      sp.lead_time_days
    FROM equipment_registry er
    JOIN equipment_catalog ec ON ec.id = er.equipment_catalog_id
    JOIN equipment_spare_parts esp ON esp.equipment_catalog_id = ec.id
    JOIN spare_parts_catalog sp ON sp.id = esp.spare_part_id AND sp.is_available = true
    WHERE er.id = p_registry_id
    ORDER BY esp.is_critical DESC, COALESCE(esp.quantity_required, 1) DESC, sp.part_name;
  END IF;
END;
$function$;

-- Add comment
COMMENT ON FUNCTION get_parts_for_registry(character varying) IS 
'Returns spare parts for a given equipment registry ID. Tries equipment_part_assignments first, falls back to equipment_spare_parts.';

-- Test the function
SELECT 'Testing get_parts_for_registry function...' as message;
SELECT * FROM get_parts_for_registry('test-id') LIMIT 1;
