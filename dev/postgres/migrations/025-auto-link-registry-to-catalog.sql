-- 025: Auto-link equipment_registry to equipment_catalog on insert/update by manufacturer+model

-- Optional performance index for matching
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE schemaname = 'public' AND indexname = 'idx_equipment_catalog_mfg_model'
    ) THEN
        CREATE INDEX idx_equipment_catalog_mfg_model 
            ON equipment_catalog (lower(manufacturer_name), lower(model_number));
    END IF;
END $$;

-- Function to set equipment_catalog_id based on manufacturer_name + model_number
CREATE OR REPLACE FUNCTION set_registry_catalog_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.equipment_catalog_id IS NULL 
       AND NEW.manufacturer_name IS NOT NULL 
       AND NEW.model_number IS NOT NULL THEN
        SELECT ec.id
        INTO NEW.equipment_catalog_id
        FROM equipment_catalog ec
        WHERE trim(lower(ec.manufacturer_name)) = trim(lower(NEW.manufacturer_name))
          AND trim(lower(ec.model_number)) = trim(lower(NEW.model_number))
        ORDER BY ec.created_at DESC
        LIMIT 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers: before insert and before update of identifying fields
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_registry_set_catalog_on_insert'
    ) THEN
        CREATE TRIGGER trg_registry_set_catalog_on_insert
        BEFORE INSERT ON equipment_registry
        FOR EACH ROW
        EXECUTE FUNCTION set_registry_catalog_id();
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_registry_set_catalog_on_update'
    ) THEN
        CREATE TRIGGER trg_registry_set_catalog_on_update
        BEFORE UPDATE OF manufacturer_name, model_number, equipment_catalog_id ON equipment_registry
        FOR EACH ROW
        WHEN (NEW.equipment_catalog_id IS DISTINCT FROM OLD.equipment_catalog_id)
        EXECUTE FUNCTION set_registry_catalog_id();
    END IF;
END $$;

COMMENT ON FUNCTION set_registry_catalog_id IS 'Auto-populates equipment_registry.equipment_catalog_id by matching manufacturer_name+model_number';
