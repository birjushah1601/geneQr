-- Migration: 021-compatibility-views.sql
-- Description: Phase 4 – Create compatibility views to aid transition from legacy marketplace equipment to v2 catalog
-- Date: 2025-11-27

-- View synthesizing marketplace equipment rows from listings + equipment_catalog + organizations
-- This view aligns closely with fields used by the marketplace repository read paths.
CREATE OR REPLACE VIEW v_market_equipment AS
SELECT
  -- Primary identity: listing-specific
  ml.id::text                         AS id,
  ec.model_name                       AS name,
  ec.model_number                     AS model,
  ec.description                      AS description,
  ec.specifications                   AS specifications,
  ml.price_amount                     AS price_amount,
  ml.price_currency                   AS price_currency,
  ml.sku                              AS sku,
  COALESCE(ec.image_urls, ARRAY[]::text[]) AS images,
  ml.is_active                        AS is_active,
  ml.created_at                       AS created_at,
  ml.updated_at                       AS updated_at,
  COALESCE(ml.tenant_id, 'default')   AS tenant_id,

  -- Derived category mapping (text → slug id)
  lower(replace(coalesce(ec.category,'Other'), ' ', '_')) AS category_id,
  ec.category                           AS category_name,
  NULL::text                            AS category_parent_id,

  -- Manufacturer (organizations)
  o.id                                  AS manufacturer_id,
  o.name                                AS manufacturer_name,
  coalesce(o.country, 'India')          AS manufacturer_country,
  o.website                             AS manufacturer_website
FROM marketplace_listings ml
JOIN equipment_catalog ec ON ml.equipment_catalog_id = ec.id
JOIN organizations o ON ec.manufacturer_id = o.id
WHERE ml.is_active = true AND ec.is_active = true;

COMMENT ON VIEW v_market_equipment IS 'Compatibility view mapping v2 catalog + listings to marketplace-equipment shape';

-- Helpful view for category rollups per tenant
CREATE OR REPLACE VIEW v_market_categories AS
SELECT
  lower(replace(coalesce(ec.category,'Other'), ' ', '_')) AS id,
  ec.category AS name,
  NULL::text AS parent_id,
  ml.tenant_id
FROM equipment_catalog ec
JOIN marketplace_listings ml ON ml.equipment_catalog_id = ec.id
WHERE ml.is_active = true AND ec.is_active = true
GROUP BY 1,2,3,4;

COMMENT ON VIEW v_market_categories IS 'Derived categories per tenant from v2 catalog + listings';

DO $$
BEGIN
  RAISE NOTICE '021 complete: compatibility views v_market_equipment and v_market_categories created';
END $$;
