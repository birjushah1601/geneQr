-- ============================================================================
-- Create Admin User (Simple)
-- ============================================================================
-- Creates a system admin user
-- Email: admin@geneqr.com
-- Password: Admin@123456
-- ============================================================================

-- Delete existing admin user if exists
DELETE FROM users WHERE email = 'admin@geneqr.com';

-- Insert admin user
-- Password hash for: Admin@123456 (bcrypt cost 12)
INSERT INTO users (
    email,
    password_hash,
    preferred_auth_method,
    email_verified,
    phone_verified,
    full_name,
    status,
    created_at,
    updated_at
) VALUES (
    'admin@geneqr.com',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIeWeCrZZG',
    'password',
    TRUE,
    FALSE,
    'System Administrator',
    'active',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Display result
SELECT 
    '✅ Admin user created successfully!' as status,
    email,
    full_name,
    status,
    email_verified,
    preferred_auth_method
FROM users 
WHERE email = 'admin@geneqr.com';

-- Display credentials
SELECT 
    '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════",
    'Email: admin@geneqr.com' as "Credentials",
    'Password: Admin@123456' as " ",
    'Login at: http://localhost:3000/login' as "  ",
    '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "═════════════════════════════════════════";
