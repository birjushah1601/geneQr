-- Migration: 023-ticket-parts-helpers-and-backfill.sql
-- Purpose:
-- 1) Backfill equipment_registry.equipment_catalog_id by matching on manufacturer_name + model_number
-- 2) Provide helper functions/views to fetch parts for an installed equipment (registry) or a ticket
--
-- Notes:
-- - This migration aligns with existing schema where equipment_catalog has columns:
--   manufacturer_name, model_number, product_name/product_code, etc.
-- - Uses existing parts tables: spare_parts_catalog and equipment_spare_parts

DO $$
BEGIN
  -- Backfill link from equipment_registry -> equipment_catalog using name + model match
  -- Only set when NULL; do not override existing links
  UPDATE equipment_registry er
  SET equipment_catalog_id = ec.id
  FROM equipment_catalog ec
  WHERE er.equipment_catalog_id IS NULL
    AND er.model_number IS NOT NULL AND ec.model_number IS NOT NULL
    AND trim(lower(er.model_number)) = trim(lower(ec.model_number))
    AND er.manufacturer_name IS NOT NULL AND ec.manufacturer_name IS NOT NULL
    AND trim(lower(er.manufacturer_name)) = trim(lower(ec.manufacturer_name));
END $$;

-- Helper: parts for a registry equipment
CREATE OR REPLACE FUNCTION get_parts_for_registry(
  p_registry_id VARCHAR
) RETURNS TABLE (
  spare_part_id UUID,
  part_number TEXT,
  part_name TEXT,
  unit_price NUMERIC,
  currency TEXT,
  is_critical BOOLEAN,
  quantity_required INT,
  part_category TEXT,
  stock_status TEXT,
  lead_time_days INT
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    sp.id AS spare_part_id,
    sp.part_number,
    sp.part_name,
    sp.unit_price,
    sp.currency,
    esp.is_critical,
    COALESCE(esp.quantity_required, 1) AS quantity_required,
    sp.category AS part_category,
    sp.stock_status,
    sp.lead_time_days
  FROM equipment_registry er
  JOIN equipment_catalog ec ON ec.id = er.equipment_catalog_id
  JOIN equipment_spare_parts esp ON esp.equipment_catalog_id = ec.id
  JOIN spare_parts_catalog sp ON sp.id = esp.spare_part_id AND sp.is_available = true
  WHERE er.id = p_registry_id
  ORDER BY esp.is_critical DESC, COALESCE(esp.quantity_required, 1) DESC, sp.part_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_parts_for_registry IS 'List compatible/linked parts for an installed equipment (equipment_registry row) via its linked equipment_catalog_id';

-- Helper: parts for a ticket (resolves equipment via service_tickets)
CREATE OR REPLACE FUNCTION get_parts_for_ticket(
  p_ticket_id VARCHAR
) RETURNS TABLE (
  spare_part_id UUID,
  part_number TEXT,
  part_name TEXT,
  unit_price NUMERIC,
  currency TEXT,
  is_critical BOOLEAN,
  quantity_required INT,
  part_category TEXT,
  stock_status TEXT,
  lead_time_days INT
) AS $$
BEGIN
  RETURN QUERY
  SELECT *
  FROM get_parts_for_registry(
    (
      SELECT t.equipment_id
      FROM service_tickets t
      WHERE t.id = p_ticket_id
      LIMIT 1
    )
  );
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_parts_for_ticket IS 'List parts for a service ticket by resolving its equipment to equipment_registry → equipment_catalog → parts';

-- Convenience view to query parts directly for tickets
CREATE OR REPLACE VIEW v_ticket_parts AS
SELECT
  t.id AS ticket_id,
  t.ticket_number,
  er.id AS registry_id,
  ec.id AS equipment_catalog_id,
  sp.id AS spare_part_id,
  sp.part_number,
  sp.part_name,
  sp.unit_price,
  sp.currency,
  esp.is_critical,
  COALESCE(esp.quantity_required, 1) AS quantity_required,
  sp.category AS part_category,
  sp.stock_status,
  sp.lead_time_days
FROM service_tickets t
JOIN equipment_registry er ON er.id = t.equipment_id
JOIN equipment_catalog ec ON ec.id = er.equipment_catalog_id
JOIN equipment_spare_parts esp ON esp.equipment_catalog_id = ec.id
JOIN spare_parts_catalog sp ON sp.id = esp.spare_part_id AND sp.is_available = true;

COMMENT ON VIEW v_ticket_parts IS 'Parts available for assignment/purchase for each ticket, via registry→catalog linkage';

DO $$
BEGIN
  RAISE NOTICE '023 complete: backfilled registry→catalog where possible; added get_parts_for_registry(), get_parts_for_ticket(), and v_ticket_parts';
END $$;
