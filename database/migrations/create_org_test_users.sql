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
        RAISE NOTICE '❌ No manufacturer organization found!';
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
    
    RAISE NOTICE '✅ Manufacturer user created: % for org: %', manufacturer_user_id, manufacturer_org_name;
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
        RAISE NOTICE '❌ No hospital organization found!';
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
    
    RAISE NOTICE '✅ Hospital user created: % for org: %', hospital_user_id, hospital_org_name;
END $$;

-- 3. DISTRIBUTOR USER
DO $$
DECLARE
    distributor_org_id UUID;
    distributor_user_id UUID;
    distributor_org_name TEXT;
BEGIN
    SELECT id, name INTO distributor_org_id, distributor_org_name
    FROM organizations 
    WHERE org_type = 'distributor' 
    AND status = 'active'
    LIMIT 1;
    
    IF distributor_org_id IS NULL THEN
        RAISE NOTICE '❌ No distributor organization found!';
        RETURN;
    END IF;
    
    DELETE FROM users WHERE email = 'distributor@geneqr.com';
    
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'distributor@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG',
        'password',
        TRUE,
        'Distributor Admin',
        'active'
    )
    RETURNING id INTO distributor_user_id;
    
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        distributor_user_id,
        distributor_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE '✅ Distributor user created: % for org: %', distributor_user_id, distributor_org_name;
END $$;

-- 4. DEALER USER
DO $$
DECLARE
    dealer_org_id UUID;
    dealer_user_id UUID;
    dealer_org_name TEXT;
BEGIN
    SELECT id, name INTO dealer_org_id, dealer_org_name
    FROM organizations 
    WHERE org_type = 'dealer' 
    AND status = 'active'
    LIMIT 1;
    
    IF dealer_org_id IS NULL THEN
        RAISE NOTICE '❌ No dealer organization found!';
        RETURN;
    END IF;
    
    DELETE FROM users WHERE email = 'dealer@geneqr.com';
    
    INSERT INTO users (email, password_hash, preferred_auth_method, email_verified, full_name, status)
    VALUES (
        'dealer@geneqr.com',
        '$2a$12$7LrnRmZTJ.4qXs1Qdmr4w.zpAzY0OoFN2elf1H8UApwZASfKr4yNG',
        'password',
        TRUE,
        'Dealer Admin',
        'active'
    )
    RETURNING id INTO dealer_user_id;
    
    INSERT INTO user_organizations (user_id, organization_id, role, is_primary, status)
    VALUES (
        dealer_user_id,
        dealer_org_id,
        'admin',
        TRUE,
        'active'
    );
    
    RAISE NOTICE '✅ Dealer user created: % for org: %', dealer_user_id, dealer_org_name;
END $$;

-- Display summary
SELECT 
    '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════════════════════";
SELECT '✅ TEST USERS CREATED SUCCESSFULLY!' as "Status";
SELECT '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════════════════════";

-- Display all created test users
SELECT 
    u.email as "Email",
    u.full_name as "Full Name",
    o.name as "Organization",
    o.org_type as "Org Type",
    uo.role as "Role",
    CASE WHEN uo.is_primary THEN '✅' ELSE '❌' END as "Primary"
FROM users u
JOIN user_organizations uo ON u.id = uo.user_id
JOIN organizations o ON uo.organization_id = o.id
WHERE u.email LIKE '%@geneqr.com'
ORDER BY o.org_type;

SELECT '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════════════════════";
SELECT 'Password for all users: password' as "Credentials";
SELECT 'Login at: http://localhost:3000/login' as "URL";
SELECT '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━' as "════════════════════════════════════════════════════════";
