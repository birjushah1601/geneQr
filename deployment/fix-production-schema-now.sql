-- Quick fix to make current production database work with current code
-- Run this on production to bridge the gap until new repo is ready

BEGIN;

-- 1. Engineers: Code uses e.name in queries, but some DBs might have full_name
-- Add name as alias if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='engineers' AND column_name='name') THEN
        -- If name column doesn't exist, add it
        -- Check if full_name exists instead
        IF EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='engineers' AND column_name='full_name') THEN
            ALTER TABLE engineers ADD COLUMN name TEXT GENERATED ALWAYS AS (full_name) STORED;
            RAISE NOTICE 'Added engineers.name as alias to full_name';
        ELSE
            ALTER TABLE engineers ADD COLUMN name TEXT;
            RAISE NOTICE 'Added engineers.name column';
        END IF;
    END IF;
END $$;

-- 2. Service Tickets: Code expects requester_org_id but DB has organization_id
-- Add requester_org_id as copy of organization_id
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='service_tickets' AND column_name='requester_org_id') THEN
        ALTER TABLE service_tickets ADD COLUMN requester_org_id VARCHAR(50);
        UPDATE service_tickets SET requester_org_id = organization_id WHERE requester_org_id IS NULL;
        RAISE NOTICE 'Added service_tickets.requester_org_id';
    END IF;
END $$;

-- 3. Keep organization_id and requester_org_id in sync
CREATE OR REPLACE FUNCTION sync_ticket_org_ids()
RETURNS TRIGGER AS $$
BEGIN
    -- Sync both ways
    IF NEW.organization_id IS DISTINCT FROM OLD.organization_id THEN
        NEW.requester_org_id = NEW.organization_id;
    END IF;
    IF NEW.requester_org_id IS DISTINCT FROM OLD.requester_org_id THEN
        NEW.organization_id = NEW.requester_org_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sync_ticket_orgs ON service_tickets;
CREATE TRIGGER trg_sync_ticket_orgs
    BEFORE INSERT OR UPDATE ON service_tickets
    FOR EACH ROW
    EXECUTE FUNCTION sync_ticket_org_ids();

-- 4. Verify equipment_registry exists (it should from your restore)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables 
                   WHERE table_name='equipment_registry') THEN
        RAISE EXCEPTION 'equipment_registry table missing! Re-run database restore.';
    END IF;
END $$;

COMMIT;

-- Verification
SELECT 
    'Schema Fix Complete' as status,
    (SELECT COUNT(*) FROM engineers WHERE name IS NOT NULL) as engineers_with_name,
    (SELECT COUNT(*) FROM service_tickets WHERE requester_org_id IS NOT NULL) as tickets_with_requester_org,
    (SELECT COUNT(*) FROM equipment_registry) as equipment_in_registry;
