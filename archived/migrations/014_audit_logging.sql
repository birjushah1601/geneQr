-- Migration: 014_audit_logging.sql
-- Description: Comprehensive audit logging system for tracking all important operations
-- Created: December 22, 2025

-- ============================================================================
-- AUDIT LOG TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Event Information
    event_type VARCHAR(100) NOT NULL,           -- e.g., 'ticket_created', 'equipment_updated'
    event_category VARCHAR(50) NOT NULL,        -- e.g., 'equipment', 'ticket', 'auth', 'parts'
    event_action VARCHAR(50) NOT NULL,          -- e.g., 'create', 'update', 'delete', 'view'
    event_status VARCHAR(20) NOT NULL,          -- 'success', 'failure', 'denied'
    
    -- User/Actor Information
    user_id UUID,                               -- User who performed the action (NULL for public)
    user_email VARCHAR(255),                    -- User email
    user_role VARCHAR(50),                      -- User role at time of action
    organization_id UUID,                       -- Organization context
    organization_type VARCHAR(50),              -- Organization type
    
    -- Resource Information
    resource_type VARCHAR(50),                  -- e.g., 'equipment', 'ticket', 'engineer'
    resource_id VARCHAR(100),                   -- ID of the resource affected
    resource_name VARCHAR(255),                 -- Name/title of the resource
    
    -- Request Information
    ip_address INET,                            -- Client IP address
    user_agent TEXT,                            -- Browser/client info
    request_method VARCHAR(10),                 -- HTTP method (GET, POST, etc.)
    request_path VARCHAR(500),                  -- API endpoint path
    request_query TEXT,                         -- Query parameters (JSON)
    
    -- Change Tracking
    old_values JSONB,                           -- Previous state (for updates)
    new_values JSONB,                           -- New state (for creates/updates)
    changed_fields TEXT[],                      -- Array of changed field names
    
    -- Additional Context
    metadata JSONB,                             -- Additional context-specific data
    error_message TEXT,                         -- Error details if event_status = 'failure'
    duration_ms INTEGER,                        -- Operation duration in milliseconds
    
    -- Rate Limiting Context
    is_rate_limited BOOLEAN DEFAULT FALSE,      -- Was this a rate limit violation?
    rate_limit_key VARCHAR(255),                -- Rate limit key (QR code, IP, etc.)
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes for common queries
    CONSTRAINT audit_logs_event_status_check CHECK (event_status IN ('success', 'failure', 'denied'))
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Index for searching by event type and category
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX idx_audit_logs_event_category ON audit_logs(event_category);
CREATE INDEX idx_audit_logs_event_action ON audit_logs(event_action);
CREATE INDEX idx_audit_logs_event_status ON audit_logs(event_status);

-- Index for user tracking
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_user_email ON audit_logs(user_email);

-- Index for organization filtering
CREATE INDEX idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_organization_type ON audit_logs(organization_type);

-- Index for resource tracking
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);

-- Index for IP-based queries (security monitoring)
CREATE INDEX idx_audit_logs_ip_address ON audit_logs(ip_address);

-- Index for rate limiting queries
CREATE INDEX idx_audit_logs_rate_limited ON audit_logs(is_rate_limited) WHERE is_rate_limited = TRUE;

-- Index for time-based queries (most common)
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- Composite index for common query patterns
CREATE INDEX idx_audit_logs_user_time ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_org_time ON audit_logs(organization_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource_time ON audit_logs(resource_type, resource_id, created_at DESC);

-- ============================================================================
-- PARTITION STRATEGY (Optional - for high-volume systems)
-- ============================================================================

-- Note: Partitioning can be added later if audit log volume grows
-- Example: Partition by month for easier archiving and performance

-- ============================================================================
-- RETENTION POLICY (COMMENT FOR REFERENCE)
-- ============================================================================

-- Consider implementing data retention:
-- - Keep 90 days of detailed logs
-- - Archive older logs to cold storage
-- - Delete logs older than 1 year (based on compliance requirements)

-- Example cleanup function (to be scheduled):
/*
CREATE OR REPLACE FUNCTION cleanup_old_audit_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM audit_logs 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '1 year';
END;
$$ LANGUAGE plpgsql;
*/

-- ============================================================================
-- SAMPLE EVENT TYPES (REFERENCE)
-- ============================================================================

-- Authentication Events:
-- - auth_login_success, auth_login_failure
-- - auth_logout, auth_token_refresh
-- - auth_password_reset

-- Equipment Events:
-- - equipment_created, equipment_updated, equipment_deleted
-- - equipment_qr_generated, equipment_qr_scanned
-- - equipment_imported

-- Ticket Events:
-- - ticket_created, ticket_updated, ticket_deleted
-- - ticket_assigned, ticket_status_changed
-- - ticket_parts_added, ticket_diagnosed

-- Engineer Events:
-- - engineer_created, engineer_updated, engineer_deleted
-- - engineer_assigned, engineer_unassigned

-- Parts Events:
-- - parts_viewed, parts_added_to_cart
-- - parts_assigned_to_ticket

-- Security Events:
-- - rate_limit_exceeded
-- - unauthorized_access_attempt
-- - suspicious_activity_detected

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE audit_logs IS 'Comprehensive audit log for tracking all system operations, security events, and user actions';
COMMENT ON COLUMN audit_logs.event_type IS 'Specific event that occurred (e.g., ticket_created, equipment_qr_scanned)';
COMMENT ON COLUMN audit_logs.event_category IS 'High-level category of event (equipment, ticket, auth, parts, etc.)';
COMMENT ON COLUMN audit_logs.event_action IS 'Action type: create, update, delete, view, assign, etc.';
COMMENT ON COLUMN audit_logs.event_status IS 'Outcome: success, failure, or denied';
COMMENT ON COLUMN audit_logs.ip_address IS 'Client IP address for security tracking';
COMMENT ON COLUMN audit_logs.old_values IS 'Previous state before update (NULL for creates)';
COMMENT ON COLUMN audit_logs.new_values IS 'New state after update (NULL for deletes)';
COMMENT ON COLUMN audit_logs.metadata IS 'Additional context-specific information as JSON';
COMMENT ON COLUMN audit_logs.is_rate_limited IS 'Flag indicating if this event triggered rate limiting';

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

-- Verify table creation
SELECT 'Audit logging table created successfully' AS status;
SELECT 'Created ' || COUNT(*) || ' indexes for audit_logs' AS index_status 
FROM pg_indexes 
WHERE tablename = 'audit_logs';
