-- Migration: Create service tickets schema for Field Service Management

-- Service tickets table
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

-- Indexes
CREATE INDEX IF NOT EXISTS idx_ticket_number ON service_tickets(ticket_number);
CREATE INDEX IF NOT EXISTS idx_ticket_equipment ON service_tickets(equipment_id);
CREATE INDEX IF NOT EXISTS idx_ticket_customer ON service_tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_ticket_status ON service_tickets(status);
CREATE INDEX IF NOT EXISTS idx_ticket_priority ON service_tickets(priority);
CREATE INDEX IF NOT EXISTS idx_ticket_engineer ON service_tickets(assigned_engineer_id);
CREATE INDEX IF NOT EXISTS idx_ticket_source ON service_tickets(source);
CREATE INDEX IF NOT EXISTS idx_ticket_created_at ON service_tickets(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ticket_sla_breach ON service_tickets(sla_breached) WHERE sla_breached = true;

-- GIN indexes for JSONB
CREATE INDEX IF NOT EXISTS idx_ticket_parts_used ON service_tickets USING GIN (parts_used);
CREATE INDEX IF NOT EXISTS idx_ticket_photos ON service_tickets USING GIN (photos);

-- Auto-update trigger
CREATE OR REPLACE FUNCTION update_ticket_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ticket_updated_at_trigger
    BEFORE UPDATE ON service_tickets
    FOR EACH ROW
    EXECUTE FUNCTION update_ticket_updated_at();

-- Ticket comments table
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

CREATE INDEX IF NOT EXISTS idx_comment_ticket ON ticket_comments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_comment_created_at ON ticket_comments(created_at DESC);

-- Ticket status history table
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

CREATE INDEX IF NOT EXISTS idx_history_ticket ON ticket_status_history(ticket_id);
CREATE INDEX IF NOT EXISTS idx_history_changed_at ON ticket_status_history(changed_at DESC);

-- Helper function: Get ticket statistics
CREATE OR REPLACE FUNCTION get_ticket_statistics()
RETURNS TABLE (
    total_tickets BIGINT,
    new_count BIGINT,
    assigned_count BIGINT,
    in_progress_count BIGINT,
    resolved_count BIGINT,
    closed_count BIGINT,
    sla_breached_count BIGINT,
    avg_resolution_hours NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_tickets,
        COUNT(*) FILTER (WHERE status = 'new')::BIGINT as new_count,
        COUNT(*) FILTER (WHERE status = 'assigned')::BIGINT as assigned_count,
        COUNT(*) FILTER (WHERE status = 'in_progress')::BIGINT as in_progress_count,
        COUNT(*) FILTER (WHERE status = 'resolved')::BIGINT as resolved_count,
        COUNT(*) FILTER (WHERE status = 'closed')::BIGINT as closed_count,
        COUNT(*) FILTER (WHERE sla_breached = true)::BIGINT as sla_breached_count,
        COALESCE(AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)) / 3600), 0)::NUMERIC as avg_resolution_hours
    FROM service_tickets;
END;
$$ LANGUAGE plpgsql;

-- Helper function: Get overdue tickets
CREATE OR REPLACE FUNCTION get_overdue_tickets()
RETURNS TABLE (
    ticket_id VARCHAR,
    ticket_number VARCHAR,
    equipment_name VARCHAR,
    customer_name VARCHAR,
    priority VARCHAR,
    sla_resolution_due TIMESTAMP WITH TIME ZONE,
    hours_overdue NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        t.id,
        t.ticket_number,
        t.equipment_name,
        t.customer_name,
        t.priority,
        t.sla_resolution_due,
        EXTRACT(EPOCH FROM (NOW() - t.sla_resolution_due)) / 3600 as hours_overdue
    FROM service_tickets t
    WHERE t.sla_resolution_due IS NOT NULL
      AND t.resolved_at IS NULL
      AND NOW() > t.sla_resolution_due
      AND t.status NOT IN ('closed', 'cancelled')
    ORDER BY t.sla_resolution_due ASC;
END;
$$ LANGUAGE plpgsql;

COMMENT ON TABLE service_tickets IS 'Service tickets for field service management - tracks equipment service requests from creation to resolution';
