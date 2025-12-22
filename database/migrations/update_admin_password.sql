-- Update admin password to a simple one for testing
-- Password: admin123

UPDATE users 
SET password_hash = '$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW',
    updated_at = CURRENT_TIMESTAMP
WHERE email = 'admin@geneqr.com';

-- Verify update
SELECT 
    '✅ Admin password updated!' as status,
    email,
    full_name,
    status,
    preferred_auth_method
FROM users 
WHERE email = 'admin@geneqr.com';

SELECT 
    '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════",
    'Email: admin@geneqr.com' as "New Credentials",
    'Password: admin123' as " ",
    '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "═════════════════════════════════════════";
