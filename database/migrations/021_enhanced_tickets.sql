-- Migration: Enhanced Service Tickets for Anonymous Creation
-- Version: 1.0
-- Date: 2025-12-20
-- Description: Add support for anonymous ticket creation (QR/WhatsApp) with flexible contact info

-- ============================================================================
-- 1. ENHANCE SERVICE_TICKETS TABLE
-- ============================================================================

-- Add new columns to service_tickets table
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS contact_info JSONB,
ADD COLUMN IF NOT EXISTS source VARCHAR(20) CHECK (source IN ('web', 'qr', 'whatsapp', 'mobile', 'api')),
ADD COLUMN IF NOT EXISTS recaptcha_score FLOAT,
ADD COLUMN IF NOT EXISTS tracking_token VARCHAR(100) UNIQUE;

-- Update existing tickets
UPDATE service_tickets 
SET source = 'web' 
WHERE source IS NULL;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_tickets_created_by ON service_tickets(created_by);
CREATE INDEX IF NOT EXISTS idx_tickets_source ON service_tickets(source);
CREATE INDEX IF NOT EXISTS idx_tickets_tracking_token ON service_tickets(tracking_token) WHERE tracking_token IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tickets_contact_email ON service_tickets((contact_info->>'email')) WHERE contact_info IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tickets_contact_phone ON service_tickets((contact_info->>'phone')) WHERE contact_info IS NOT NULL;

-- Comments
COMMENT ON COLUMN service_tickets.created_by IS 'User who created ticket (NULL for anonymous)';
COMMENT ON COLUMN service_tickets.contact_info IS 'Contact information for anonymous tickets: {email, phone, whatsapp, name, preferred_channel}';
COMMENT ON COLUMN service_tickets.source IS 'How ticket was created: web, qr, whatsapp, mobile';
COMMENT ON COLUMN service_tickets.recaptcha_score IS 'reCAPTCHA v3 score (0.0-1.0) for spam detection';
COMMENT ON COLUMN service_tickets.tracking_token IS 'Public tracking token for anonymous users';

-- ============================================================================
-- 2. TICKET_NOTIFICATIONS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS ticket_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Ticket Reference
    ticket_id UUID NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    
    -- Notification Details
    notification_type VARCHAR(50) NOT NULL CHECK (notification_type IN (
        'created', 'assigned', 'status_updated', 'engineer_arrived', 
        'parts_ordered', 'completed', 'feedback_requested'
    )),
    
    -- Recipient (for anonymous tickets)
    recipient_email VARCHAR(255),
    recipient_phone VARCHAR(20),
    recipient_whatsapp VARCHAR(20),
    
    -- Message
    subject VARCHAR(255),
    message TEXT NOT NULL,
    
    -- Delivery Status
    sent_via JSONB DEFAULT '[]'::jsonb, -- Array: ["email", "sms", "whatsapp"]
    delivery_status JSONB DEFAULT '{}'::jsonb, -- {"email": "delivered", "sms": "failed"}
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_ticket_notif_ticket_id ON ticket_notifications(ticket_id);
CREATE INDEX idx_ticket_notif_type ON ticket_notifications(notification_type);
CREATE INDEX idx_ticket_notif_created_at ON ticket_notifications(created_at);

-- Comments
COMMENT ON TABLE ticket_notifications IS 'Notification history for service tickets';
COMMENT ON COLUMN ticket_notifications.sent_via IS 'Channels used: ["email", "sms", "whatsapp"]';

-- ============================================================================
-- 3. WHATSAPP_CONVERSATIONS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS whatsapp_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- WhatsApp Details
    phone_number VARCHAR(20) NOT NULL,
    whatsapp_id VARCHAR(100), -- Twilio WhatsApp ID
    
    -- Conversation State
    state VARCHAR(50) DEFAULT 'idle' CHECK (state IN (
        'idle', 'awaiting_equipment_id', 'awaiting_issue_description',
        'creating_ticket', 'active_ticket', 'completed'
    )),
    
    -- Current Context
    ticket_id UUID REFERENCES service_tickets(id) ON DELETE SET NULL,
    equipment_id UUID REFERENCES equipment_registry(id) ON DELETE SET NULL,
    
    -- Conversation Data
    context JSONB DEFAULT '{}'::jsonb,
    
    -- Message Count
    message_count INT DEFAULT 0,
    last_message_at TIMESTAMP,
    
    -- Timestamps
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Session timeout (30 minutes of inactivity)
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '30 minutes')
);

-- Indexes
CREATE INDEX idx_whatsapp_conv_phone ON whatsapp_conversations(phone_number);
CREATE INDEX idx_whatsapp_conv_ticket_id ON whatsapp_conversations(ticket_id);
CREATE INDEX idx_whatsapp_conv_state ON whatsapp_conversations(state);
CREATE INDEX idx_whatsapp_conv_expires_at ON whatsapp_conversations(expires_at);

-- Comments
COMMENT ON TABLE whatsapp_conversations IS 'WhatsApp conversation state management';
COMMENT ON COLUMN whatsapp_conversations.state IS 'Current conversation state for wizard flow';
COMMENT ON COLUMN whatsapp_conversations.context IS 'Temporary data collected during conversation';

-- ============================================================================
-- 4. WHATSAPP_MESSAGES TABLE (Audit Trail)
-- ============================================================================

CREATE TABLE IF NOT EXISTS whatsapp_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Conversation Reference
    conversation_id UUID REFERENCES whatsapp_conversations(id) ON DELETE CASCADE,
    
    -- Message Details
    message_sid VARCHAR(100) UNIQUE, -- Twilio Message SID
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('inbound', 'outbound')),
    
    -- Content
    from_number VARCHAR(20) NOT NULL,
    to_number VARCHAR(20) NOT NULL,
    body TEXT,
    
    -- Media
    media_url TEXT,
    media_type VARCHAR(50),
    
    -- Status (for outbound)
    status VARCHAR(20), -- queued, sent, delivered, failed
    error_message TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    delivered_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_whatsapp_msg_conv_id ON whatsapp_messages(conversation_id);
CREATE INDEX idx_whatsapp_msg_sid ON whatsapp_messages(message_sid);
CREATE INDEX idx_whatsapp_msg_direction ON whatsapp_messages(direction);
CREATE INDEX idx_whatsapp_msg_created_at ON whatsapp_messages(created_at);

-- Comments
COMMENT ON TABLE whatsapp_messages IS 'Complete audit trail of WhatsApp messages';

-- ============================================================================
-- 5. RECAPTCHA_SCORES TABLE (Anti-Spam)
-- ============================================================================

CREATE TABLE IF NOT EXISTS recaptcha_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Request Details
    ip_address INET NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'create_ticket', 'register', etc.
    
    -- reCAPTCHA Response
    score FLOAT NOT NULL CHECK (score >= 0.0 AND score <= 1.0),
    success BOOLEAN NOT NULL,
    
    -- Additional Data
    hostname VARCHAR(255),
    challenge_ts TIMESTAMP,
    
    -- Action Taken
    allowed BOOLEAN NOT NULL,
    reason VARCHAR(255),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_recaptcha_ip ON recaptcha_scores(ip_address);
CREATE INDEX idx_recaptcha_action ON recaptcha_scores(action);
CREATE INDEX idx_recaptcha_score ON recaptcha_scores(score);
CREATE INDEX idx_recaptcha_created_at ON recaptcha_scores(created_at);

-- Comments
COMMENT ON TABLE recaptcha_scores IS 'reCAPTCHA v3 scores for spam detection';
COMMENT ON COLUMN recaptcha_scores.score IS 'Score from 0.0 (bot) to 1.0 (human)';

-- ============================================================================
-- 6. HELPER FUNCTIONS
-- ============================================================================

-- Function to generate tracking token
CREATE OR REPLACE FUNCTION generate_tracking_token()
RETURNS VARCHAR AS $$
BEGIN
    RETURN 'TRK-' || UPPER(substring(md5(random()::text) from 1 for 12));
END;
$$ LANGUAGE plpgsql;

-- Function to get or create WhatsApp conversation
CREATE OR REPLACE FUNCTION get_or_create_whatsapp_conversation(p_phone VARCHAR)
RETURNS UUID AS $$
DECLARE
    v_conversation_id UUID;
BEGIN
    -- Try to get active conversation
    SELECT id INTO v_conversation_id
    FROM whatsapp_conversations
    WHERE phone_number = p_phone
      AND expires_at > NOW()
      AND state NOT IN ('completed')
    ORDER BY started_at DESC
    LIMIT 1;
    
    -- Create new if not found
    IF v_conversation_id IS NULL THEN
        INSERT INTO whatsapp_conversations (phone_number)
        VALUES (p_phone)
        RETURNING id INTO v_conversation_id;
    ELSE
        -- Update expiry
        UPDATE whatsapp_conversations
        SET expires_at = NOW() + INTERVAL '30 minutes',
            last_message_at = NOW()
        WHERE id = v_conversation_id;
    END IF;
    
    RETURN v_conversation_id;
END;
$$ LANGUAGE plpgsql;

-- Function to cleanup expired WhatsApp conversations
CREATE OR REPLACE FUNCTION cleanup_expired_whatsapp_conversations()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    UPDATE whatsapp_conversations
    SET state = 'completed',
        completed_at = NOW()
    WHERE expires_at < NOW()
      AND state NOT IN ('completed');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- 7. TRIGGERS
-- ============================================================================

-- Trigger to auto-generate tracking token for anonymous tickets
CREATE OR REPLACE FUNCTION auto_generate_tracking_token()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.created_by IS NULL AND NEW.tracking_token IS NULL THEN
        NEW.tracking_token := generate_tracking_token();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_tracking_token
BEFORE INSERT ON service_tickets
FOR EACH ROW
EXECUTE FUNCTION auto_generate_tracking_token();

-- Trigger to send notification on ticket creation
CREATE OR REPLACE FUNCTION notify_on_ticket_create()
RETURNS TRIGGER AS $$
BEGIN
    -- Insert notification record
    INSERT INTO ticket_notifications (
        ticket_id,
        notification_type,
        recipient_email,
        recipient_phone,
        recipient_whatsapp,
        subject,
        message
    )
    SELECT
        NEW.id,
        'created',
        NEW.contact_info->>'email',
        NEW.contact_info->>'phone',
        NEW.contact_info->>'whatsapp',
        'Service Ticket Created: ' || NEW.ticket_number,
        format('Your service ticket %s has been created successfully. Track it at: %s',
               NEW.ticket_number,
               'https://app.ServQR.com/track/' || COALESCE(NEW.tracking_token, NEW.ticket_number));
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_notify_ticket_created
AFTER INSERT ON service_tickets
FOR EACH ROW
EXECUTE FUNCTION notify_on_ticket_create();

-- ============================================================================
-- 8. VIEWS
-- ============================================================================

-- View: Anonymous tickets pending assignment
CREATE OR REPLACE VIEW v_anonymous_tickets_pending AS
SELECT 
    t.id,
    t.ticket_number,
    t.issue_description,
    t.priority,
    t.source,
    t.contact_info,
    t.recaptcha_score,
    t.created_at,
    e.equipment_model,
    e.serial_number,
    e.location AS equipment_location
FROM service_tickets t
LEFT JOIN equipment_registry e ON t.equipment_id = e.id
WHERE t.created_by IS NULL
  AND t.status IN ('pending', 'new')
ORDER BY t.created_at DESC;

-- View: WhatsApp ticket statistics
CREATE OR REPLACE VIEW v_whatsapp_ticket_stats AS
SELECT 
    DATE(t.created_at) AS date,
    COUNT(*) AS total_tickets,
    AVG(EXTRACT(EPOCH FROM (t.updated_at - t.created_at))/60) AS avg_response_time_minutes,
    COUNT(CASE WHEN t.status = 'completed' THEN 1 END) AS completed_tickets
FROM service_tickets t
WHERE t.source = 'whatsapp'
GROUP BY DATE(t.created_at)
ORDER BY date DESC;

-- View: reCAPTCHA score analysis
CREATE OR REPLACE VIEW v_recaptcha_analysis AS
SELECT 
    DATE(created_at) AS date,
    action,
    COUNT(*) AS total_requests,
    AVG(score) AS avg_score,
    COUNT(CASE WHEN score < 0.5 THEN 1 END) AS suspicious_count,
    COUNT(CASE WHEN allowed = FALSE THEN 1 END) AS blocked_count
FROM recaptcha_scores
GROUP BY DATE(created_at), action
ORDER BY date DESC, action;

-- ============================================================================
-- 9. SAMPLE DATA (for testing)
-- ============================================================================

-- Sample anonymous ticket (commented out for production)
/*
INSERT INTO service_tickets (
    ticket_number,
    equipment_id,
    issue_description,
    priority,
    status,
    contact_info,
    source,
    recaptcha_score
)
SELECT
    'TKT-ANON-001',
    id,
    'Equipment not powering on - reported via QR code',
    'high',
    'pending',
    '{"email": "user@hospital.com", "phone": "+1234567890", "name": "John Maintenance", "preferred_channel": "whatsapp"}'::jsonb,
    'qr',
    0.9
FROM equipment_registry
WHERE qr_code = 'QR-CAN-XR-002'
LIMIT 1;
*/

-- ============================================================================
-- 10. GRANT PERMISSIONS (adjust as needed)
-- ============================================================================

-- Grant read access to application role
-- GRANT SELECT ON ticket_notifications TO app_role;
-- GRANT SELECT ON whatsapp_conversations TO app_role;
-- GRANT SELECT ON whatsapp_messages TO app_role;
-- GRANT SELECT ON recaptcha_scores TO app_role;

-- Grant write access for ticket operations
-- GRANT INSERT, UPDATE ON service_tickets TO app_role;
-- GRANT INSERT ON ticket_notifications TO app_role;
-- GRANT ALL ON whatsapp_conversations TO app_role;
-- GRANT ALL ON whatsapp_messages TO app_role;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

-- Add migration tracking
INSERT INTO schema_migrations (version, name, applied_at)
VALUES (
    21,
    'enhanced_tickets',
    NOW()
) ON CONFLICT (version) DO NOTHING;

COMMENT ON SCHEMA public IS 'Enhanced tickets migration applied - Support for anonymous creation via QR/WhatsApp';
