-- Migration: 017-link-servicing-to-catalog.sql
-- Description: Link installed equipment to master catalog and enforce FK from tickets to registry (Phase 0)
-- Date: 2025-11-27

-- 1) Add equipment_catalog_id to equipment_registry and FK to equipment_catalog
ALTER TABLE IF EXISTS equipment_registry
    ADD COLUMN IF NOT EXISTS equipment_catalog_id UUID NULL;

DO $$
BEGIN
    -- Add FK if not present
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints tc
        WHERE tc.table_name = 'equipment_registry'
          AND tc.constraint_type = 'FOREIGN KEY'
          AND tc.constraint_name = 'fk_equipment_registry_catalog'
    ) THEN
        ALTER TABLE equipment_registry
            ADD CONSTRAINT fk_equipment_registry_catalog
            FOREIGN KEY (equipment_catalog_id)
            REFERENCES equipment_catalog(id)
            ON UPDATE CASCADE
            ON DELETE SET NULL;
    END IF;
END $$;

-- Helpful index for lookups by model
CREATE INDEX IF NOT EXISTS idx_equipment_registry_catalog_id
    ON equipment_registry(equipment_catalog_id);

-- 2) Add FK for service_tickets(equipment_id) → equipment_registry(id)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints tc
        WHERE tc.table_name = 'service_tickets'
          AND tc.constraint_type = 'FOREIGN KEY'
          AND tc.constraint_name = 'fk_service_tickets_equipment'
    ) THEN
        -- Ensure column types are compatible (both are VARCHAR(32) per current schema)
        ALTER TABLE service_tickets
            ADD CONSTRAINT fk_service_tickets_equipment
            FOREIGN KEY (equipment_id)
            REFERENCES equipment_registry(id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;
    END IF;
END $$;

-- Notices
DO $$
BEGIN
    RAISE NOTICE '017 complete: equipment_registry.equipment_catalog_id added with FK to equipment_catalog';
    RAISE NOTICE '017 complete: service_tickets.equipment_id now FK to equipment_registry(id)';
END $$;
