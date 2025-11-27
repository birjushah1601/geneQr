-- Migration: 018-create-marketplace-listings.sql
-- Description: Create marketplace_listings to decouple marketplace pricing/listing from master equipment catalog (Phase 0)
-- Date: 2025-11-27

-- 1) Table: marketplace_listings
CREATE TABLE IF NOT EXISTS marketplace_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Master model linkage
    equipment_catalog_id UUID NOT NULL REFERENCES equipment_catalog(id),

    -- Seller/tenant context
    seller_org_id UUID NOT NULL REFERENCES organizations(id),
    tenant_id TEXT, -- optional multi-tenant scoping

    -- Listing details
    title TEXT,
    sku TEXT,
    price_amount NUMERIC(12,2) NOT NULL,
    price_currency TEXT NOT NULL DEFAULT 'INR',
    availability_status TEXT NOT NULL DEFAULT 'in_stock',
    stock_quantity INT,
    images TEXT[],
    is_active BOOLEAN NOT NULL DEFAULT true,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,

    CONSTRAINT chk_availability_status CHECK (
        availability_status IN ('in_stock','out_of_stock','preorder','discontinued')
    ),
    CONSTRAINT chk_price_amount CHECK (price_amount >= 0)
);

-- 2) Indexes
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_equipment ON marketplace_listings(equipment_catalog_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_seller ON marketplace_listings(seller_org_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_active ON marketplace_listings(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_price ON marketplace_listings(price_amount);

-- 3) updated_at trigger
CREATE OR REPLACE FUNCTION update_marketplace_listing_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_marketplace_listings_updated
    BEFORE UPDATE ON marketplace_listings
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_listing_timestamp();

-- 4) Comments
COMMENT ON TABLE marketplace_listings IS 'Per-seller/per-tenant marketplace listings referencing master equipment_catalog';
COMMENT ON COLUMN marketplace_listings.seller_org_id IS 'Organization offering the listing (seller)';
COMMENT ON COLUMN marketplace_listings.tenant_id IS 'Optional tenant scope for multitenancy';

-- Notices
DO $$
BEGIN
    RAISE NOTICE '018 complete: marketplace_listings created';
END $$;
