-- ============================================================================
-- Create Admin User
-- ============================================================================
-- Creates a system admin user with full permissions
-- Email: admin@geneqr.com
-- Password: Admin@123456
-- ============================================================================

DO $$
DECLARE
    admin_role_id UUID;
    admin_user_id UUID;
BEGIN
    -- Get or create system_admin role
    SELECT id INTO admin_role_id FROM roles WHERE name = 'system_admin' LIMIT 1;
    
    IF admin_role_id IS NULL THEN
        -- Create system_admin role if it doesn't exist
        INSERT INTO roles (name, display_name, description, org_type, permissions, is_system)
        VALUES (
            'system_admin',
            'System Administrator',
            'Full system access - can manage all organizations and users',
            NULL,
            '["*"]'::jsonb,
            TRUE
        )
        RETURNING id INTO admin_role_id;
        
        RAISE NOTICE 'Created system_admin role: %', admin_role_id;
    ELSE
        RAISE NOTICE 'Using existing system_admin role: %', admin_role_id;
    END IF;
    
    -- Check if admin user already exists
    SELECT id INTO admin_user_id FROM users WHERE email = 'admin@geneqr.com';
    
    IF admin_user_id IS NOT NULL THEN
        RAISE NOTICE 'Admin user already exists: %', admin_user_id;
        -- Update password and settings
        UPDATE users SET
            password_hash = '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIeWeCrZZG', -- Admin@123456
            preferred_auth_method = 'password',
            email_verified = TRUE,
            is_active = TRUE,
            role_id = admin_role_id,
            full_name = 'System Administrator',
            updated_at = CURRENT_TIMESTAMP
        WHERE email = 'admin@geneqr.com';
        
        RAISE NOTICE 'Updated admin user password and settings';
    ELSE
        -- Insert new admin user
        INSERT INTO users (
            email,
            password_hash,
            preferred_auth_method,
            email_verified,
            phone_verified,
            is_active,
            role_id,
            full_name,
            created_at,
            updated_at
        ) VALUES (
            'admin@geneqr.com',
            '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIeWeCrZZG', -- Admin@123456
            'password',
            TRUE,
            FALSE,
            TRUE,
            admin_role_id,
            'System Administrator',
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        )
        RETURNING id INTO admin_user_id;
        
        RAISE NOTICE 'Created new admin user: %', admin_user_id;
    END IF;
    
    -- Display success message
    RAISE NOTICE '';
    RAISE NOTICE '✅ Admin user ready!';
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    RAISE NOTICE 'Email: admin@geneqr.com';
    RAISE NOTICE 'Password: Admin@123456';
    RAISE NOTICE 'Role: System Administrator (Full Access)';
    RAISE NOTICE '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━';
    
END $$;

-- Verify admin user was created
SELECT 
    u.id,
    u.email,
    u.full_name,
    r.name as role_name,
    r.display_name as role_display_name,
    u.email_verified,
    u.is_active,
    u.preferred_auth_method
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.email = 'admin@geneqr.com';
