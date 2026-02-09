-- Migration: Ticket Notification System
-- Description: Add tables for tracking tokens and notification logs
-- Date: 2026-02-07

BEGIN;

-- Create ticket_tracking_tokens table
CREATE TABLE IF NOT EXISTS ticket_tracking_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id VARCHAR(50) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Indexes
    CONSTRAINT check_expires_future CHECK (expires_at > created_at)
);

CREATE INDEX idx_tracking_token ON ticket_tracking_tokens(token);
CREATE INDEX idx_tracking_ticket_id ON ticket_tracking_tokens(ticket_id);
CREATE INDEX idx_tracking_expires ON ticket_tracking_tokens(expires_at);

COMMENT ON TABLE ticket_tracking_tokens IS 'Tracking tokens for public ticket access';
COMMENT ON COLUMN ticket_tracking_tokens.token IS 'Secure random token for public access';
COMMENT ON COLUMN ticket_tracking_tokens.expires_at IS 'Token expiration timestamp (default 30 days)';

-- Create notification_log table
CREATE TABLE IF NOT EXISTS notification_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id VARCHAR(50) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'sent',
    error_message TEXT,
    
    -- Constraints
    CONSTRAINT check_notification_type CHECK (
        notification_type IN ('manual', 'ticket_created', 'daily_digest')
    ),
    CONSTRAINT check_notification_status CHECK (
        status IN ('sent', 'failed')
    )
);

CREATE INDEX idx_notification_ticket_id ON notification_log(ticket_id);
CREATE INDEX idx_notification_type ON notification_log(notification_type);
CREATE INDEX idx_notification_status ON notification_log(status);
CREATE INDEX idx_notification_sent_at ON notification_log(sent_at DESC);
CREATE INDEX idx_notification_failed ON notification_log(status, sent_at) 
    WHERE status = 'failed';

COMMENT ON TABLE notification_log IS 'Log of all ticket notification attempts';
COMMENT ON COLUMN notification_log.notification_type IS 'Type: manual, ticket_created, daily_digest';
COMMENT ON COLUMN notification_log.status IS 'Status: sent or failed';
COMMENT ON COLUMN notification_log.error_message IS 'Error details if status is failed';

-- Add helpful function to clean up expired tokens
CREATE OR REPLACE FUNCTION cleanup_expired_tracking_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM ticket_tracking_tokens
    WHERE expires_at <= NOW();
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_tracking_tokens IS 'Delete expired tracking tokens, returns count deleted';

COMMIT;
