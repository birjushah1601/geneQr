-- Migration: Authentication System
-- Version: 1.0
-- Date: 2025-12-20
-- Description: Core authentication tables for OTP-first login system

-- ============================================================================
-- 1. USERS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Primary Login Identifiers (at least one required)
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    
    -- Password (OPTIONAL - for fallback only)
    password_hash VARCHAR(255),
    
    -- Auth Preferences
    preferred_auth_method VARCHAR(20) DEFAULT 'otp' CHECK (preferred_auth_method IN ('otp', 'password')),
    
    -- Verification Status
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    
    -- Profile Information
    full_name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    
    -- Account Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted', 'pending')),
    
    -- Security
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMP,
    
    -- Timestamps
    last_login TIMESTAMP,
    last_otp_sent TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Metadata
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Constraints
    CONSTRAINT email_or_phone_required CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- Indexes
CREATE INDEX idx_users_email ON users(email) WHERE email IS NOT NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Comments
COMMENT ON TABLE users IS 'User accounts with OTP-first authentication';
COMMENT ON COLUMN users.preferred_auth_method IS 'Preferred authentication method: otp or password';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt password hash (optional - for fallback only)';

-- ============================================================================
-- 2. OTP_CODES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User Reference (optional for first-time users)
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    -- Delivery Target
    email VARCHAR(255),
    phone VARCHAR(20),
    
    -- OTP Details
    code VARCHAR(10) NOT NULL,
    code_hash VARCHAR(255) NOT NULL, -- Hashed OTP for security
    
    -- Delivery Method
    delivery_method VARCHAR(20) NOT NULL CHECK (delivery_method IN ('email', 'sms', 'whatsapp')),
    
    -- Purpose
    purpose VARCHAR(50) NOT NULL CHECK (purpose IN ('login', 'verify', 'reset', 'register')),
    
    -- Status
    used BOOLEAN DEFAULT FALSE,
    attempts INT DEFAULT 0,
    
    -- Expiry
    expires_at TIMESTAMP NOT NULL,
    
    -- Metadata
    device_info JSONB,
    ip_address INET,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    used_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT otp_email_or_phone CHECK (email IS NOT NULL OR phone IS NOT NULL)
);

-- Indexes
CREATE INDEX idx_otp_user_id ON otp_codes(user_id);
CREATE INDEX idx_otp_email ON otp_codes(email) WHERE email IS NOT NULL;
CREATE INDEX idx_otp_phone ON otp_codes(phone) WHERE phone IS NOT NULL;
CREATE INDEX idx_otp_expires_at ON otp_codes(expires_at);
CREATE INDEX idx_otp_used ON otp_codes(used);

-- Auto-delete expired OTPs (cleanup job)
CREATE INDEX idx_otp_expired ON otp_codes(expires_at) WHERE used = FALSE;

-- Comments
COMMENT ON TABLE otp_codes IS 'One-Time Password codes for authentication';
COMMENT ON COLUMN otp_codes.code_hash IS 'Hashed OTP code for security (compare on verify)';
COMMENT ON COLUMN otp_codes.attempts IS 'Number of verification attempts (max 3)';

-- ============================================================================
-- 3. REFRESH_TOKENS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User Reference
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Token
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    
    -- Device Information
    device_info JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address INET,
    user_agent TEXT,
    
    -- Status
    revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    revoke_reason VARCHAR(100),
    
    -- Expiry
    expires_at TIMESTAMP NOT NULL,
    
    -- Usage Tracking
    last_used_at TIMESTAMP,
    usage_count INT DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);

-- Auto-delete expired tokens (cleanup job)
CREATE INDEX idx_refresh_tokens_expired ON refresh_tokens(expires_at) WHERE revoked = FALSE;

-- Comments
COMMENT ON TABLE refresh_tokens IS 'Refresh tokens for JWT token rotation';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'Hashed refresh token (SHA-256)';
COMMENT ON COLUMN refresh_tokens.usage_count IS 'Number of times token was used (for rotation tracking)';

-- ============================================================================
-- 4. AUTH_AUDIT_LOG TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS auth_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- User Reference
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Action Details
    action VARCHAR(100) NOT NULL,
    -- Actions: login_success, login_failed, otp_sent, otp_verified, otp_failed,
    --          logout, password_changed, password_reset, account_locked, etc.
    
    -- Result
    success BOOLEAN NOT NULL,
    
    -- Context
    ip_address INET,
    user_agent TEXT,
    device_info JSONB,
    
    -- Additional Data
    metadata JSONB DEFAULT '{}'::jsonb,
    error_message TEXT,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_auth_audit_user_id ON auth_audit_log(user_id);
CREATE INDEX idx_auth_audit_action ON auth_audit_log(action);
CREATE INDEX idx_auth_audit_created_at ON auth_audit_log(created_at);
CREATE INDEX idx_auth_audit_success ON auth_audit_log(success);
CREATE INDEX idx_auth_audit_ip_address ON auth_audit_log(ip_address);

-- Partitioning by month (optional for large scale)
-- CREATE INDEX idx_auth_audit_created_at_month ON auth_audit_log((date_trunc('month', created_at)));

-- Comments
COMMENT ON TABLE auth_audit_log IS 'Comprehensive audit log for all authentication events';
COMMENT ON COLUMN auth_audit_log.action IS 'Action performed (login_success, login_failed, etc.)';
COMMENT ON COLUMN auth_audit_log.metadata IS 'Additional context-specific data';

-- ============================================================================
-- 5. ORGANIZATIONS TABLE (Enhanced)
-- ============================================================================

-- Note: This may already exist, so we use CREATE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Basic Information
    name VARCHAR(255) NOT NULL,
    org_type VARCHAR(50) NOT NULL CHECK (org_type IN ('manufacturer', 'hospital', 'laboratory', 'distributor', 'dealer')),
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'pending', 'suspended', 'deleted')),
    onboarding_status VARCHAR(50) DEFAULT 'created' CHECK (onboarding_status IN ('created', 'documents_pending', 'under_review', 'verified', 'active')),
    
    -- Contact Information
    address JSONB,
    contact JSONB,
    
    -- Configuration
    settings JSONB DEFAULT '{}'::jsonb,
    
    -- Metadata (org-specific fields)
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    verified_at TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_organizations_type ON organizations(org_type);
CREATE INDEX IF NOT EXISTS idx_organizations_status ON organizations(status);
CREATE INDEX IF NOT EXISTS idx_organizations_created_at ON organizations(created_at);

-- Comments
COMMENT ON COLUMN organizations.org_type IS 'Type: manufacturer, hospital, laboratory, distributor, dealer';
COMMENT ON COLUMN organizations.metadata IS 'Org-specific fields (license_number, bed_count, etc.)';

-- ============================================================================
-- 6. USER_ORGANIZATIONS TABLE (Many-to-Many)
-- ============================================================================

CREATE TABLE IF NOT EXISTS user_organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- References
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    
    -- Role & Permissions
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'manager', 'engineer', 'viewer', 'technician')),
    permissions JSONB DEFAULT '[]'::jsonb, -- Array of permission strings
    
    -- Status
    is_primary BOOLEAN DEFAULT FALSE, -- Primary organization for user
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'pending')),
    
    -- Timestamps
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP,
    
    -- Unique constraint
    UNIQUE(user_id, organization_id)
);

-- Indexes
CREATE INDEX idx_user_orgs_user_id ON user_organizations(user_id);
CREATE INDEX idx_user_orgs_org_id ON user_organizations(organization_id);
CREATE INDEX idx_user_orgs_role ON user_organizations(role);
CREATE INDEX idx_user_orgs_status ON user_organizations(status);
CREATE INDEX idx_user_orgs_primary ON user_organizations(is_primary) WHERE is_primary = TRUE;

-- Comments
COMMENT ON TABLE user_organizations IS 'Many-to-many mapping of users to organizations';
COMMENT ON COLUMN user_organizations.permissions IS 'Array of permission strings specific to this user-org relationship';
COMMENT ON COLUMN user_organizations.is_primary IS 'User\'s primary/default organization';

-- ============================================================================
-- 7. ROLES TABLE (RBAC)
-- ============================================================================

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Role Details
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Organization Type (NULL = global role)
    org_type VARCHAR(50) CHECK (org_type IN ('manufacturer', 'hospital', 'laboratory', 'distributor', 'dealer')),
    
    -- Permissions
    permissions JSONB NOT NULL DEFAULT '[]'::jsonb, -- Array of permission strings
    
    -- System Role (cannot be deleted)
    is_system BOOLEAN DEFAULT FALSE,
    
    -- Status
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_org_type ON roles(org_type);
CREATE INDEX idx_roles_status ON roles(status);

-- Comments
COMMENT ON TABLE roles IS 'Role definitions for RBAC system';
COMMENT ON COLUMN roles.org_type IS 'Organization type this role applies to (NULL = all types)';
COMMENT ON COLUMN roles.is_system IS 'System roles cannot be deleted or modified';

-- ============================================================================
-- 8. NOTIFICATION_PREFERENCES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Owner
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Notification Channels
    email_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    whatsapp_notifications BOOLEAN DEFAULT TRUE,
    in_app_notifications BOOLEAN DEFAULT TRUE,
    
    -- Event Subscriptions
    events JSONB DEFAULT '{
        "ticket_created": true,
        "ticket_assigned": true,
        "ticket_updated": true,
        "ticket_completed": true,
        "engineer_assigned": true,
        "parts_ordered": true,
        "payment_received": true
    }'::jsonb,
    
    -- Quiet Hours
    quiet_hours JSONB, -- {"start": "22:00", "end": "08:00", "timezone": "America/New_York"}
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Unique per user
    UNIQUE(user_id)
);

-- Indexes
CREATE INDEX idx_notif_prefs_user_id ON notification_preferences(user_id);

-- Comments
COMMENT ON TABLE notification_preferences IS 'User notification preferences';
COMMENT ON COLUMN notification_preferences.events IS 'JSON object of event subscriptions';

-- ============================================================================
-- 9. HELPER FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notif_prefs_updated_at BEFORE UPDATE ON notification_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- 10. SEED DATA - Default Roles
-- ============================================================================

-- Global Admin Role
INSERT INTO roles (id, name, display_name, description, org_type, permissions, is_system)
VALUES (
    gen_random_uuid(),
    'system_admin',
    'System Administrator',
    'Full system access - can manage all organizations and users',
    NULL,
    '["*"]'::jsonb,
    TRUE
) ON CONFLICT (name) DO NOTHING;

-- Manufacturer Roles
INSERT INTO roles (id, name, display_name, description, org_type, permissions, is_system)
VALUES 
(gen_random_uuid(), 'manufacturer_admin', 'Manufacturer Admin', 'Full access to manufacturer organization', 'manufacturer', 
 '["manage_organization", "manage_users", "manage_engineers", "view_all_tickets", "manage_equipment", "view_reports", "manage_contracts"]'::jsonb, TRUE),
(gen_random_uuid(), 'manufacturer_manager', 'Service Manager', 'Manage service operations', 'manufacturer',
 '["view_organization", "assign_engineers", "view_all_tickets", "update_tickets", "view_reports"]'::jsonb, TRUE),
(gen_random_uuid(), 'manufacturer_engineer', 'Field Engineer', 'Field service engineer', 'manufacturer',
 '["view_assigned_tickets", "update_assigned_tickets", "add_parts", "upload_photos"]'::jsonb, TRUE)
ON CONFLICT (name) DO NOTHING;

-- Hospital/Lab Roles
INSERT INTO roles (id, name, display_name, description, org_type, permissions, is_system)
VALUES
(gen_random_uuid(), 'hospital_admin', 'Hospital Administrator', 'Full access to hospital operations', 'hospital',
 '["manage_organization", "manage_users", "create_tickets", "view_all_tickets", "manage_equipment", "view_reports", "approve_invoices"]'::jsonb, TRUE),
(gen_random_uuid(), 'hospital_manager', 'Biomedical Manager', 'Manage biomedical equipment', 'hospital',
 '["create_tickets", "view_all_tickets", "manage_equipment", "view_reports"]'::jsonb, TRUE),
(gen_random_uuid(), 'hospital_technician', 'Biomedical Technician', 'Create and track service tickets', 'hospital',
 '["create_tickets", "view_own_tickets", "track_tickets"]'::jsonb, TRUE),
(gen_random_uuid(), 'hospital_viewer', 'Viewer', 'Read-only access', 'hospital',
 '["view_tickets", "view_equipment", "view_reports"]'::jsonb, TRUE)
ON CONFLICT (name) DO NOTHING;

-- Laboratory Roles (similar to hospital)
INSERT INTO roles (id, name, display_name, description, org_type, permissions, is_system)
SELECT 
    gen_random_uuid(),
    REPLACE(name, 'hospital_', 'laboratory_'),
    REPLACE(display_name, 'Hospital', 'Laboratory'),
    REPLACE(description, 'hospital', 'laboratory'),
    'laboratory',
    permissions,
    TRUE
FROM roles 
WHERE org_type = 'hospital'
ON CONFLICT (name) DO NOTHING;

-- Distributor Roles
INSERT INTO roles (id, name, display_name, description, org_type, permissions, is_system)
VALUES
(gen_random_uuid(), 'distributor_admin', 'Distributor Admin', 'Full access to distributor operations', 'distributor',
 '["manage_organization", "manage_users", "manage_inventory", "view_sales", "create_tickets", "view_reports"]'::jsonb, TRUE),
(gen_random_uuid(), 'distributor_sales', 'Sales Representative', 'Manage sales and customers', 'distributor',
 '["view_inventory", "create_orders", "view_customers", "create_tickets"]'::jsonb, TRUE)
ON CONFLICT (name) DO NOTHING;

-- ============================================================================
-- 11. VIEWS FOR CONVENIENCE
-- ============================================================================

-- View: Active users with their organizations
CREATE OR REPLACE VIEW v_users_with_organizations AS
SELECT 
    u.id AS user_id,
    u.email,
    u.phone,
    u.full_name,
    u.status AS user_status,
    u.last_login,
    uo.organization_id,
    o.name AS organization_name,
    o.org_type,
    uo.role,
    uo.is_primary,
    uo.joined_at
FROM users u
JOIN user_organizations uo ON u.id = uo.user_id
JOIN organizations o ON uo.organization_id = o.id
WHERE u.status = 'active' AND uo.status = 'active';

-- View: Failed login attempts summary
CREATE OR REPLACE VIEW v_failed_login_summary AS
SELECT 
    DATE(created_at) AS date,
    COUNT(*) AS total_attempts,
    COUNT(DISTINCT ip_address) AS unique_ips,
    COUNT(DISTINCT user_id) AS unique_users
FROM auth_audit_log
WHERE action = 'login_failed'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- ============================================================================
-- 12. CLEANUP JOBS (To be run by scheduler)
-- ============================================================================

-- Function to delete expired OTPs
CREATE OR REPLACE FUNCTION cleanup_expired_otps()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM otp_codes
    WHERE expires_at < NOW() AND used = FALSE;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Function to delete expired refresh tokens
CREATE OR REPLACE FUNCTION cleanup_expired_refresh_tokens()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM refresh_tokens
    WHERE expires_at < NOW() AND revoked = FALSE;
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Function to archive old audit logs (optional)
-- Run monthly to archive logs older than 1 year
CREATE OR REPLACE FUNCTION archive_old_audit_logs()
RETURNS INTEGER AS $$
DECLARE
    archived_count INTEGER;
BEGIN
    -- This could move to an archive table or delete
    -- For now, just count (implement archival strategy as needed)
    SELECT COUNT(*) INTO archived_count
    FROM auth_audit_log
    WHERE created_at < NOW() - INTERVAL '1 year';
    
    RETURN archived_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

-- Add migration tracking
INSERT INTO schema_migrations (version, name, applied_at)
VALUES (
    20,
    'authentication_system',
    NOW()
) ON CONFLICT (version) DO NOTHING;

COMMENT ON SCHEMA public IS 'Authentication system migration applied - Version 1.0';
