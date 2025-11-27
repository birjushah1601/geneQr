-- Migration: 020-backfill-equipment-registry-catalog-id.sql
-- Description: Phase 2 – Backfill equipment_registry.equipment_catalog_id by matching manufacturer + model_number
-- Date: 2025-11-27

BEGIN;

-- 1) Build a temporary matching table from registry -> catalog
CREATE TEMP TABLE tmp_registry_catalog_match AS
SELECT 
  er.id AS registry_id,
  ec.id AS equipment_catalog_id
FROM equipment_registry er
JOIN organizations o ON lower(o.name) = lower(er.manufacturer_name) AND o.org_type = 'manufacturer'
JOIN equipment_catalog ec ON ec.manufacturer_id = o.id
WHERE er.model_number IS NOT NULL AND ec.model_number IS NOT NULL
  AND trim(lower(ec.model_number)) = trim(lower(er.model_number));

CREATE INDEX ON tmp_registry_catalog_match(registry_id);

-- 2) Update equipment_registry with matched catalog IDs
UPDATE equipment_registry er
SET equipment_catalog_id = m.equipment_catalog_id
FROM tmp_registry_catalog_match m
WHERE er.id = m.registry_id
  AND (er.equipment_catalog_id IS NULL OR er.equipment_catalog_id <> m.equipment_catalog_id);

COMMIT;

DO $$
DECLARE
  updated_count integer;
BEGIN
  SELECT COUNT(*) INTO updated_count FROM equipment_registry WHERE equipment_catalog_id IS NOT NULL;
  RAISE NOTICE '020 complete: equipment_registry.equipment_catalog_id backfilled. Total linked: %', updated_count;
END $$;
