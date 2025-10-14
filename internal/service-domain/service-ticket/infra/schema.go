package infra

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

// EnsureServiceTicketSchema creates the required tables and indices if they don't exist.
func EnsureServiceTicketSchema(ctx context.Context, pool *pgxpool.Pool) error {
    // Service tickets table and related objects
    schema := `
CREATE TABLE IF NOT EXISTS service_tickets (
    id VARCHAR(32) PRIMARY KEY,
    ticket_number VARCHAR(50) UNIQUE NOT NULL,

    -- Equipment & Customer
    equipment_id VARCHAR(32) NOT NULL,
    qr_code VARCHAR(255),
    serial_number VARCHAR(255) NOT NULL,
    equipment_name VARCHAR(500) NOT NULL,
    customer_id VARCHAR(32),
    customer_name VARCHAR(500) NOT NULL,
    customer_phone VARCHAR(20),
    customer_whatsapp VARCHAR(20),

    -- Issue details
    issue_category VARCHAR(100),
    issue_description TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    severity VARCHAR(50),

    -- Source
    source VARCHAR(50) NOT NULL,
    source_message_id VARCHAR(255),

    -- Assignment
    assigned_engineer_id VARCHAR(32),
    assigned_engineer_name VARCHAR(255),
    assigned_at TIMESTAMP WITH TIME ZONE,

    -- Status & Timeline
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    closed_at TIMESTAMP WITH TIME ZONE,

    -- SLA tracking
    sla_response_due TIMESTAMP WITH TIME ZONE,
    sla_resolution_due TIMESTAMP WITH TIME ZONE,
    sla_breached BOOLEAN NOT NULL DEFAULT false,

    -- Resolution
    resolution_notes TEXT,
    parts_used JSONB DEFAULT '[]'::jsonb,
    labor_hours DECIMAL(5,2) DEFAULT 0,
    cost DECIMAL(15,2) DEFAULT 0,

    -- Media
    photos JSONB DEFAULT '[]'::jsonb,
    videos JSONB DEFAULT '[]'::jsonb,
    documents JSONB DEFAULT '[]'::jsonb,

    -- AMC linkage
    amc_contract_id VARCHAR(32),
    covered_under_amc BOOLEAN DEFAULT false,

    -- Metadata
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,

    CONSTRAINT ticket_status_check CHECK (status IN ('new', 'assigned', 'in_progress', 'on_hold', 'resolved', 'closed', 'cancelled')),
    CONSTRAINT ticket_priority_check CHECK (priority IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT ticket_source_check CHECK (source IN ('whatsapp', 'web', 'phone', 'email', 'scheduled'))
);

-- Indexes will be created after ensuring legacy columns (see below)

CREATE OR REPLACE FUNCTION update_ticket_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'ticket_updated_at_trigger'
    ) THEN
        CREATE TRIGGER ticket_updated_at_trigger
            BEFORE UPDATE ON service_tickets
            FOR EACH ROW
            EXECUTE FUNCTION update_ticket_updated_at();
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS ticket_comments (
    id VARCHAR(32) PRIMARY KEY,
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    comment_type VARCHAR(50) NOT NULL,
    author_id VARCHAR(32),
    author_name VARCHAR(255) NOT NULL,
    comment TEXT NOT NULL,
    attachments JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT comment_type_check CHECK (comment_type IN ('customer', 'engineer', 'internal', 'system'))
);

-- Indexes will be created after ensuring legacy columns (see below)

CREATE TABLE IF NOT EXISTS ticket_status_history (
    id VARCHAR(32) PRIMARY KEY,
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    from_status VARCHAR(50),
    to_status VARCHAR(50) NOT NULL,
    changed_by VARCHAR(32),
    changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    reason TEXT,
    CONSTRAINT history_from_status_check CHECK (from_status IN ('new', 'assigned', 'in_progress', 'on_hold', 'resolved', 'closed', 'cancelled')),
    CONSTRAINT history_to_status_check CHECK (to_status IN ('new', 'assigned', 'in_progress', 'on_hold', 'resolved', 'closed', 'cancelled'))
);

-- Indexes will be created after ensuring legacy columns (see below)
-- Backfill/compatibility: ensure columns exist on legacy installations
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS ticket_number VARCHAR(50);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS qr_code VARCHAR(255);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS serial_number VARCHAR(255);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS equipment_name VARCHAR(500);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS customer_id VARCHAR(32);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS customer_name VARCHAR(500);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(20);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS customer_whatsapp VARCHAR(20);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS issue_category VARCHAR(100);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS issue_description TEXT;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS severity VARCHAR(50);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS source VARCHAR(50);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS source_message_id VARCHAR(255);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS assigned_engineer_id VARCHAR(32);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS assigned_engineer_name VARCHAR(255);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS acknowledged_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS started_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS closed_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS sla_response_due TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS sla_resolution_due TIMESTAMP WITH TIME ZONE;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS sla_breached BOOLEAN DEFAULT false;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS resolution_notes TEXT;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS labor_hours DECIMAL(5,2) DEFAULT 0;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS cost DECIMAL(15,2) DEFAULT 0;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS photos JSONB DEFAULT '[]'::jsonb;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS videos JSONB DEFAULT '[]'::jsonb;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS documents JSONB DEFAULT '[]'::jsonb;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS amc_contract_id VARCHAR(32);
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS covered_under_amc BOOLEAN DEFAULT false;
-- Optional Phase 4 columns
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS responsible_org_id UUID NULL;
ALTER TABLE service_tickets ADD COLUMN IF NOT EXISTS policy_provenance JSONB DEFAULT '{}'::jsonb;
-- Create indexes (after columns are ensured)
CREATE UNIQUE INDEX IF NOT EXISTS uq_ticket_number ON service_tickets(ticket_number);
CREATE INDEX IF NOT EXISTS idx_ticket_number ON service_tickets(ticket_number);
CREATE INDEX IF NOT EXISTS idx_ticket_equipment ON service_tickets(equipment_id);
CREATE INDEX IF NOT EXISTS idx_ticket_customer ON service_tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_ticket_status ON service_tickets(status);
CREATE INDEX IF NOT EXISTS idx_ticket_priority ON service_tickets(priority);
CREATE INDEX IF NOT EXISTS idx_ticket_engineer ON service_tickets(assigned_engineer_id);
CREATE INDEX IF NOT EXISTS idx_ticket_source ON service_tickets(source);
CREATE INDEX IF NOT EXISTS idx_ticket_created_at ON service_tickets(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ticket_sla_breach ON service_tickets(sla_breached) WHERE sla_breached = true;
CREATE INDEX IF NOT EXISTS idx_ticket_parts_used ON service_tickets USING GIN (parts_used);
CREATE INDEX IF NOT EXISTS idx_ticket_photos ON service_tickets USING GIN (photos);

CREATE INDEX IF NOT EXISTS idx_comment_ticket ON ticket_comments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_comment_created_at ON ticket_comments(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_history_ticket ON ticket_status_history(ticket_id);
CREATE INDEX IF NOT EXISTS idx_history_changed_at ON ticket_status_history(changed_at DESC);

-- Service policies (Phase 4)
CREATE TABLE IF NOT EXISTS service_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    rules JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- SLA policies (Phase 6)
CREATE TABLE IF NOT EXISTS sla_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NULL,
    name TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    rules JSONB NOT NULL DEFAULT '{}'::jsonb, -- {priority:{critical:{resp:1,res:4},...}}
    effective_from TIMESTAMP WITH TIME ZONE NULL,
    effective_to   TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sla_policies_org_active ON sla_policies(org_id, active);

-- Events + Webhooks (Phase 6)
CREATE TABLE IF NOT EXISTS service_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status TEXT NOT NULL DEFAULT 'queued', -- queued|delivered|failed
    attempt_count INT NOT NULL DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMP WITH TIME ZONE NULL
);
CREATE INDEX IF NOT EXISTS idx_events_status_created ON service_events(status, created_at);

CREATE TABLE IF NOT EXISTS webhook_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    endpoint_url TEXT NOT NULL,
    event_types TEXT[] NOT NULL, -- ['ticket.created', ...] or ['*']
    secret TEXT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_webhooks_active ON webhook_subscriptions(active);

CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES service_events(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES webhook_subscriptions(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'queued', -- queued|delivered|failed
    attempt_count INT NOT NULL DEFAULT 0,
    last_error TEXT NULL,
    last_attempt_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMP WITH TIME ZONE NULL
);
CREATE INDEX IF NOT EXISTS idx_deliveries_status ON webhook_deliveries(status);
`

    _, err := pool.Exec(ctx, schema)
    return err
}
