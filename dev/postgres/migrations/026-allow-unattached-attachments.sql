-- Migration 026: Allow unattached attachments (nullable ticket_id) and linking later
-- - Make ticket_attachments.ticket_id nullable to support pre-creation uploads (e.g., WhatsApp)
-- - Add partial indexes for performance

-- 1) Relax NOT NULL on ticket_id
ALTER TABLE ticket_attachments
    ALTER COLUMN ticket_id DROP NOT NULL;

-- 2) Ensure foreign key still enforced when present
-- (FK already exists; no change needed)

-- 3) Index optimizations
-- Existing idx_attachments_ticket remains valid; add partial index for non-null values for planner
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE schemaname = 'public' AND indexname = 'idx_attachments_ticket_not_null'
    ) THEN
        CREATE INDEX idx_attachments_ticket_not_null ON ticket_attachments(ticket_id) WHERE ticket_id IS NOT NULL;
    END IF;
END $$;

-- 4) Convenience unique constraint on (source, source_message_id) when provided to avoid duplicates from WhatsApp
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uniq_attachment_source_message'
    ) THEN
        ALTER TABLE ticket_attachments
        ADD CONSTRAINT uniq_attachment_source_message UNIQUE (source, source_message_id);
    END IF;
EXCEPTION WHEN undefined_column THEN
    -- source_message_id column may be nullable or absent in some envs; ignore if not present
    NULL;
END $$;
