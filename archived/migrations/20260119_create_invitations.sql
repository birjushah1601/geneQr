-- Migration: Create invitations table for user invitations
-- Created: 2026-01-19

-- Invitations table: Store user invitation tokens
CREATE TABLE IF NOT EXISTS invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Invitation details
    email VARCHAR(255) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL, -- admin, manager, viewer, engineer
    
    -- Inviter information
    invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
    invited_at TIMESTAMP DEFAULT NOW(),
    
    -- Token and security
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    
    -- Status tracking
    status VARCHAR(20) DEFAULT 'pending', -- pending, accepted, expired, revoked
    
    -- Acceptance tracking
    accepted_at TIMESTAMP,
    accepted_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Additional metadata (name, phone, etc.)
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_invitations_token ON invitations(token);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_invitations_status ON invitations(status);
CREATE INDEX idx_invitations_organization ON invitations(organization_id);
CREATE INDEX idx_invitations_expires ON invitations(expires_at);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_invitations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at
CREATE TRIGGER trg_invitations_updated_at
    BEFORE UPDATE ON invitations
    FOR EACH ROW
    EXECUTE FUNCTION update_invitations_updated_at();

-- Comments for documentation
COMMENT ON TABLE invitations IS 'Stores user invitation tokens for organization membership';
COMMENT ON COLUMN invitations.token IS 'Unique invitation token sent via email';
COMMENT ON COLUMN invitations.status IS 'pending: awaiting acceptance, accepted: user created, expired: token expired, revoked: manually cancelled';
COMMENT ON COLUMN invitations.metadata IS 'Additional data like name, phone from onboarding';
