-- Migration: 019-populate-catalog-and-listings.sql
-- Description: Phase 1 – Populate equipment_catalog from legacy equipment tables and seed marketplace_listings
-- Date: 2025-11-27

BEGIN;

-- 0) Ensure required extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto; -- for gen_random_uuid()

-- 1) Create manufacturer organizations from legacy manufacturers (idempotent)
--    Strategy: upsert by (name, org_type='manufacturer'); keep manufacturers.id in external_ref for traceability
INSERT INTO organizations (id, name, org_type, website, year_established, external_ref, tenant_id, created_by)
SELECT 
    gen_random_uuid() AS id,
    m.name,
    'manufacturer'::text AS org_type,
    m.website,
    NULLIF(m.established, 0) AS year_established,
    m.id AS external_ref,
    COALESCE(m.tenant_id, 'default') AS tenant_id,
    'migration-019'
FROM manufacturers m
LEFT JOIN organizations o
  ON o.external_ref = m.id OR (lower(o.name) = lower(m.name) AND o.org_type = 'manufacturer')
WHERE o.id IS NULL
ON CONFLICT DO NOTHING;

-- 2) Helper temp mapping from legacy manufacturers to organizations
CREATE TEMP TABLE tmp_manufacturer_org_map AS
SELECT m.id AS legacy_manufacturer_id, o.id AS org_id
FROM manufacturers m
JOIN organizations o
  ON o.external_ref = m.id OR (lower(o.name) = lower(m.name) AND o.org_type = 'manufacturer');

CREATE INDEX ON tmp_manufacturer_org_map(legacy_manufacturer_id);

-- 3) Insert equipment models into equipment_catalog (idempotent via unique idx on (manufacturer_id, model_number) where is_active)
INSERT INTO equipment_catalog (
    id,
    manufacturer_id,
    equipment_type,
    model_number,
    model_name,
    category,
    sub_category,
    specifications,
    description,
    image_urls,
    is_active,
    created_by
)
SELECT
    gen_random_uuid() AS id,
    map.org_id AS manufacturer_id,
    -- Use category name as equipment_type fallback
    COALESCE(NULLIF(c.name, ''), 'Other') AS equipment_type,
    NULLIF(e.model, '') AS model_number,
    e.name AS model_name,
    -- Normalize category into allowed set
    CASE 
        WHEN c.name ILIKE '%Diagnostic%' THEN 'Diagnostic'
        WHEN c.name ILIKE '%Surgical%' THEN 'Surgical'
        WHEN c.name ILIKE '%Laborator%' THEN 'Laboratory'
        WHEN c.name ILIKE '%Imaging%' THEN 'Imaging'
        WHEN c.name ILIKE '%Monitor%' THEN 'Monitoring'
        WHEN c.name ILIKE '%Therap%' THEN 'Therapeutic'
        WHEN c.name ILIKE '%Steriliz%' THEN 'Sterilization'
        ELSE 'Other'
    END AS category,
    NULL::text AS sub_category,
    COALESCE(e.specifications, '{}'::jsonb) AS specifications,
    e.description,
    e.images AS image_urls,
    COALESCE(e.is_active, true) AS is_active,
    'migration-019' AS created_by
FROM equipment e
JOIN categories c ON e.category_id = c.id
JOIN tmp_manufacturer_org_map map ON map.legacy_manufacturer_id = e.manufacturer_id
LEFT JOIN equipment_catalog ec
  ON ec.manufacturer_id = map.org_id AND ec.model_number = NULLIF(e.model, '')
WHERE ec.id IS NULL
  AND NULLIF(e.model, '') IS NOT NULL; -- require a model number to ensure uniqueness

-- 4) Seed marketplace_listings (OEM listing per model) if missing
--    Use manufacturer organization as seller, carry over price fields
INSERT INTO marketplace_listings (
    id,
    equipment_catalog_id,
    seller_org_id,
    tenant_id,
    title,
    sku,
    price_amount,
    price_currency,
    availability_status,
    stock_quantity,
    images,
    is_active,
    created_by
)
SELECT
    gen_random_uuid() AS id,
    ec.id AS equipment_catalog_id,
    ec.manufacturer_id AS seller_org_id,
    COALESCE(e.tenant_id, 'default') AS tenant_id,
    e.name AS title,
    e.sku,
    e.price_amount,
    COALESCE(NULLIF(e.price_currency, ''), 'INR') AS price_currency,
    CASE WHEN COALESCE(e.is_active, true) THEN 'in_stock' ELSE 'discontinued' END AS availability_status,
    NULL::int AS stock_quantity,
    e.images,
    COALESCE(e.is_active, true) AS is_active,
    'migration-019' AS created_by
FROM equipment e
JOIN categories c ON e.category_id = c.id
JOIN tmp_manufacturer_org_map map ON map.legacy_manufacturer_id = e.manufacturer_id
JOIN equipment_catalog ec ON ec.manufacturer_id = map.org_id AND ec.model_number = NULLIF(e.model, '')
LEFT JOIN marketplace_listings ml
  ON ml.equipment_catalog_id = ec.id AND ml.seller_org_id = ec.manufacturer_id
WHERE ml.id IS NULL;

COMMIT;

DO $$
BEGIN
  RAISE NOTICE '019 complete: Populated organizations (manufacturers), equipment_catalog, and marketplace_listings';
END $$;
