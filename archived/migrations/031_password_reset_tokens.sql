-- Migration: 031 - Create password_reset_tokens table
-- Purpose: Store secure tokens for password setup/reset
-- Used when creating new manufacturers to allow secure password setup

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    used_at TIMESTAMP,
    CONSTRAINT valid_expiry CHECK (expires_at > created_at)
);

-- Index for fast token lookup
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);

-- Index for cleanup of expired tokens
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at ON password_reset_tokens(expires_at);

-- Index for user lookup
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);

COMMENT ON TABLE password_reset_tokens IS 'Secure tokens for password setup and reset operations';
COMMENT ON COLUMN password_reset_tokens.token IS 'Cryptographically secure token sent to user';
COMMENT ON COLUMN password_reset_tokens.expires_at IS 'Token expiration time (typically 24-48 hours)';
COMMENT ON COLUMN password_reset_tokens.used_at IS 'When the token was used (null if not yet used)';
