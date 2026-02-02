-- ============================================================================
-- Create Test Users for Different Organization Types
-- ============================================================================
-- Password for all users: password
-- ============================================================================

-- 1. MANUFACTURER USER
DO $$
DECLARE
    manufacturer_org_id UUID;
    manufacturer_user_id UUID;
    manufacturer_org_name TEXT;
BEGIN
    -- Get a manufacturer organization
    SELECT id, name INTO manufacturer_org_id, manufacturer_org_name
    FROM organizations 
    WHERE org_type = 'manufacturer' 
    AND status = 'active'
    LIMIT 1;
    
    IF manufacturer_org_id IS NULL THEN
        RAISE NOTICE 'âŒ No manufacturer organization found!';
        RETURN;
    END IF;
    
    -- Delete existing manufacturer test user
    DELETE FROM users WHERE email = 'manufacturer@geneqr.com';
    
    -- Create user
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'manufacturer@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG', -- password: password
        'password',
        TRUE,
        'Manufacturer Admin',
        'active'
    )
    RETURNING id INTO manufacturer_user_id;
    
    -- Link user to organization
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        manufacturer_user_id,
        manufacturer_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE 'âœ… Manufacturer user created: % for org: %', manufacturer_user_id, manufacturer_org_name;
END $$;

-- 2. HOSPITAL USER
DO $$
DECLARE
    hospital_org_id UUID;
    hospital_user_id UUID;
    hospital_org_name TEXT;
BEGIN
    SELECT id, name INTO hospital_org_id, hospital_org_name
    FROM organizations 
    WHERE org_type = 'hospital' 
    AND status = 'active'
    LIMIT 1;
    
    IF hospital_org_id IS NULL THEN
        RAISE NOTICE 'âŒ No hospital organization found!';
        RETURN;
    END IF;
    
    DELETE FROM users WHERE email = 'hospital@geneqr.com';
    
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'hospital@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG',
        'password',
        TRUE,
        'Hospital Admin',
        'active'
    )
    RETURNING id INTO hospital_user_id;
    
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        hospital_user_id,
        hospital_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE 'âœ… Hospital user created: % for org: %', hospital_user_id, hospital_org_name;
END $$;

-- 3. Channel Partner USER
DO $$
DECLARE
    channel_partner_org_id UUID;
    channel_partner_user_id UUID;
    channel_partner_org_name TEXT;
BEGIN
    SELECT id, name INTO channel_partner_org_id, channel_partner_org_name
    FROM organizations 
    WHERE org_type = 'Channel Partner' 
    AND status = 'active'
    LIMIT 1;
    
    IF channel_partner_org_id IS NULL THEN
        RAISE NOTICE 'âŒ No Channel Partner organization found!';
        RETURN;
    END IF;
    
    DELETE FROM users WHERE email = 'Channel Partner@geneqr.com';
    
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'Channel Partner@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG',
        'password',
        TRUE,
        'Channel Partner Admin',
        'active'
    )
    RETURNING id INTO channel_partner_user_id;
    
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        channel_partner_user_id,
        channel_partner_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE 'âœ… Channel Partner user created: % for org: %', channel_partner_user_id, channel_partner_org_name;
END $$;

-- 4. Sub-sub_SUB_DEALER USER
DO $$
DECLARE
    sub_sub_Sub-sub_SUB_DEALER_org_id UUID;
    sub_sub_Sub-sub_SUB_DEALER_user_id UUID;
    sub_sub_Sub-sub_SUB_DEALER_org_name TEXT;
BEGIN
    SELECT id, name INTO sub_sub_Sub-sub_SUB_DEALER_org_id, sub_sub_Sub-sub_SUB_DEALER_org_name
    FROM organizations 
    WHERE org_type = 'Sub-sub_SUB_DEALER' 
    AND status = 'active'
    LIMIT 1;
    
    IF sub_sub_Sub-sub_SUB_DEALER_org_id IS NULL THEN
        RAISE NOTICE 'âŒ No Sub-sub_SUB_DEALER organization found!';
        RETURN;
    END IF;
    
    DELETE FROM users WHERE email = 'Sub-sub_SUB_DEALER@geneqr.com';
    
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'Sub-sub_SUB_DEALER@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG',
        'password',
        TRUE,
        'Sub-sub_SUB_DEALER Admin',
        'active'
    )
    RETURNING id INTO sub_sub_Sub-sub_SUB_DEALER_user_id;
    
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        sub_sub_Sub-sub_SUB_DEALER_user_id,
        sub_sub_Sub-sub_SUB_DEALER_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE 'âœ… Sub-sub_SUB_DEALER user created: % for org: %', sub_sub_Sub-sub_SUB_DEALER_user_id, sub_sub_Sub-sub_SUB_DEALER_org_name;
END $$;

-- Display summary
SELECT 
    'â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”' as "â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•";
SELECT 'âœ… TEST USERS CREATED SUCCESSFULLY!' as "Status";
SELECT 'â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”' as "â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•";

-- Display all created test users
SELECT 
    u.email as "Email",
    u.full_name as "Full Name",
    o.name as "Organization",
    o.org_type as "Org Type",
    uo.role as "Role",
    CASE WHEN uo.is_primary THEN 'âœ…' ELSE 'âŒ' END as "Primary"
FROM users u
JOIN user_organizations uo ON u.id = uo.user_id
JOIN organizations o ON uo.organization_id = o.id
WHERE u.email LIKE '%@geneqr.com'
ORDER BY o.org_type;

SELECT 'â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”' as "â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•";
SELECT 'Password for all users: password' as "Credentials";
SELECT 'Login at: http://localhost:3000/login' as "URL";
SELECT 'â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”' as "â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•â•";
