-- ============================================================================
-- Populate Manufacturer Sample Data
-- ============================================================================
-- This migration adds comprehensive metadata to manufacturers and ensures
-- proper relationships with equipment_catalog and equipment_registry

-- ============================================================================
-- 1. UPDATE MANUFACTURER METADATA
-- ============================================================================

-- Siemens Healthineers India
UPDATE organizations 
SET metadata = jsonb_build_object(
    'contact_person', 'Dr. Rajesh Kumar',
    'email', 'rajesh.kumar@siemens-healthineers.com',
    'phone', '+91-80-4141-4141',
    'website', 'https://www.siemens-healthineers.com/en-in',
    'address', jsonb_build_object(
        'street', 'Olympia Technology Park, 1-A, SIDCO Industrial Estate',
        'city', 'Guindy, Chennai',
        'state', 'Tamil Nadu',
        'postal_code', '600032',
        'country', 'India'
    ),
    'business_info', jsonb_build_object(
        'gst_number', '33AACCS1119F1Z5',
        'pan_number', 'AACCS1119F',
        'established_year', 1992,
        'employee_count', 5000,
        'headquarters', 'Mumbai, Maharashtra'
    ),
    'support_info', jsonb_build_object(
        'support_email', 'service.india@siemens-healthineers.com',
        'support_phone', '+91-80-4141-4200',
        'support_hours', '24/7 Available',
        'response_time_sla', '4 hours'
    )
)
WHERE id = '11afdeec-5dee-44d4-aa5b-952703536f10';

-- Wipro GE Healthcare
UPDATE organizations 
SET metadata = jsonb_build_object(
    'contact_person', 'Priya Sharma',
    'email', 'priya.sharma@wipro-ge.com',
    'phone', '+91-80-6623-3000',
    'website', 'https://www.wipro.com/healthcare',
    'address', jsonb_build_object(
        'street', 'Wipro GE Healthcare, 72 Whitefield Main Road',
        'city', 'Bengaluru',
        'state', 'Karnataka',
        'postal_code', '560066',
        'country', 'India'
    ),
    'business_info', jsonb_build_object(
        'gst_number', '29AAACW3775F1ZH',
        'pan_number', 'AAACW3775F',
        'established_year', 1989,
        'employee_count', 3500,
        'headquarters', 'Bengaluru, Karnataka'
    ),
    'support_info', jsonb_build_object(
        'support_email', 'support@wipro-ge.com',
        'support_phone', '+91-80-6623-3100',
        'support_hours', 'Mon-Sat: 8AM-8PM',
        'response_time_sla', '6 hours'
    )
)
WHERE id = 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad';

-- Philips Healthcare India
UPDATE organizations 
SET metadata = jsonb_build_object(
    'contact_person', 'Mr. Ankit Desai',
    'email', 'ankit.desai@philips.com',
    'phone', '+91-20-6602-6000',
    'website', 'https://www.philips.co.in/healthcare',
    'address', jsonb_build_object(
        'street', 'Philips Innovation Campus, Manyata Embassy Business Park',
        'city', 'Pune',
        'state', 'Maharashtra',
        'postal_code', '411045',
        'country', 'India'
    ),
    'business_info', jsonb_build_object(
        'gst_number', '27AABCP2635L1ZN',
        'pan_number', 'AABCP2635L',
        'established_year', 1996,
        'employee_count', 4200,
        'headquarters', 'Pune, Maharashtra'
    ),
    'support_info', jsonb_build_object(
        'support_email', 'india.support@philips.com',
        'support_phone', '+91-20-6602-6100',
        'support_hours', '24/7 Available',
        'response_time_sla', '3 hours'
    )
)
WHERE id = 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad';

-- Global Manufacturer A
UPDATE organizations 
SET metadata = jsonb_build_object(
    'contact_person', 'Mr. Suresh Menon',
    'email', 'suresh.menon@globalmed.com',
    'phone', '+91-22-2875-4000',
    'website', 'https://www.globalmedmanufacturing.com',
    'address', jsonb_build_object(
        'street', 'Global Medical Tower, Bandra Kurla Complex',
        'city', 'Mumbai',
        'state', 'Maharashtra',
        'postal_code', '400051',
        'country', 'India'
    ),
    'business_info', jsonb_build_object(
        'gst_number', '27AACGM4567M1Z8',
        'pan_number', 'AACGM4567M',
        'established_year', 2005,
        'employee_count', 1500,
        'headquarters', 'Mumbai, Maharashtra'
    ),
    'support_info', jsonb_build_object(
        'support_email', 'support@globalmed.com',
        'support_phone', '+91-22-2875-4100',
        'support_hours', 'Mon-Fri: 9AM-6PM',
        'response_time_sla', '8 hours'
    )
)
WHERE id = '31370ba0-b49f-4bb6-9a6f-5d06d31b61c9';

-- ============================================================================
-- 2. UPDATE EQUIPMENT_CATALOG WITH MANUFACTURER_ID
-- ============================================================================

-- Link Siemens products to Siemens organization
UPDATE equipment_catalog 
SET manufacturer_id = '11afdeec-5dee-44d4-aa5b-952703536f10'
WHERE manufacturer_name ILIKE '%Siemens%';

-- Link Wipro GE products to Wipro GE organization
UPDATE equipment_catalog 
SET manufacturer_id = 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad'
WHERE manufacturer_name = 'Wipro GE Healthcare';

-- Link Philips products to Philips organization
UPDATE equipment_catalog 
SET manufacturer_id = 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad'
WHERE manufacturer_name ILIKE '%Philips%';

-- Link GE Healthcare products to Wipro GE (as GE Healthcare operates through Wipro in India)
UPDATE equipment_catalog 
SET manufacturer_id = 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad'
WHERE manufacturer_name = 'GE Healthcare';

-- ============================================================================
-- 3. ADD MORE EQUIPMENT TO CATALOG FOR EACH MANUFACTURER
-- ============================================================================

-- Siemens Healthineers - Add more products
INSERT INTO equipment_catalog (
    id, product_code, product_name, manufacturer_id, manufacturer_name,
    model_number, category, subcategory, base_price, currency,
    recommended_service_interval_days, estimated_lifespan_years,
    maintenance_complexity, is_active
) VALUES
(
    gen_random_uuid(), 'SIE-MRI-002', 'MAGNETOM Skyra 3T',
    '11afdeec-5dee-44d4-aa5b-952703536f10', 'Siemens Healthineers',
    'SKYRA-3T', 'MRI Scanner', '3 Tesla', 45000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'SIE-CT-001', 'SOMATOM Definition AS',
    '11afdeec-5dee-44d4-aa5b-952703536f10', 'Siemens Healthineers',
    'AS-64', 'CT Scanner', '64-slice', 35000000.00, 'INR',
    180, 12, 'high', true
),
(
    gen_random_uuid(), 'SIE-XRAY-001', 'Multix Fusion',
    '11afdeec-5dee-44d4-aa5b-952703536f10', 'Siemens Healthineers',
    'MF-MAX', 'X-Ray System', 'Digital Radiography', 8500000.00, 'INR',
    90, 10, 'medium', true
);

-- Wipro GE Healthcare - Add more products
INSERT INTO equipment_catalog (
    id, product_code, product_name, manufacturer_id, manufacturer_name,
    model_number, category, subcategory, base_price, currency,
    recommended_service_interval_days, estimated_lifespan_years,
    maintenance_complexity, is_active
) VALUES
(
    gen_random_uuid(), 'WGE-MRI-001', 'Optima MR450w',
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad', 'Wipro GE Healthcare',
    'MR450W-1.5T', 'MRI Scanner', '1.5 Tesla', 38000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'WGE-US-001', 'Voluson E10',
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad', 'Wipro GE Healthcare',
    'E10-BT20', 'Ultrasound System', '4D Imaging', 12500000.00, 'INR',
    180, 10, 'medium', true
),
(
    gen_random_uuid(), 'WGE-XRAY-001', 'Brivo XR575',
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad', 'Wipro GE Healthcare',
    'XR575-PLUS', 'X-Ray System', 'Digital', 7500000.00, 'INR',
    90, 10, 'medium', true
);

-- Philips Healthcare - Add more products
INSERT INTO equipment_catalog (
    id, product_code, product_name, manufacturer_id, manufacturer_name,
    model_number, category, subcategory, base_price, currency,
    recommended_service_interval_days, estimated_lifespan_years,
    maintenance_complexity, is_active
) VALUES
(
    gen_random_uuid(), 'PHI-MRI-001', 'Ingenia 1.5T',
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad', 'Philips Healthcare India',
    'ING-1.5T-S', 'MRI Scanner', '1.5 Tesla', 42000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'PHI-US-001', 'EPIQ Elite',
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad', 'Philips Healthcare India',
    'EPIQ-7C', 'Ultrasound System', 'Premium Cardiac', 18500000.00, 'INR',
    180, 12, 'medium', true
),
(
    gen_random_uuid(), 'PHI-PM-001', 'IntelliVue MX850',
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad', 'Philips Healthcare India',
    'MX850-ICU', 'Patient Monitor', 'ICU Monitor', 850000.00, 'INR',
    90, 8, 'low', true
);

-- ============================================================================
-- 4. VERIFY DATA
-- ============================================================================

-- Show manufacturer counts
SELECT 
    o.name as manufacturer,
    COUNT(DISTINCT ec.id) as equipment_catalog_count,
    COUNT(DISTINCT er.id) as equipment_registry_count
FROM organizations o
LEFT JOIN equipment_catalog ec ON ec.manufacturer_id = o.id
LEFT JOIN equipment_registry er ON er.manufacturer_name = o.name
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY o.name;

-- Show manufacturer details with contact info
SELECT 
    id,
    name,
    metadata->>'contact_person' as contact,
    metadata->>'email' as email,
    metadata->>'phone' as phone,
    metadata->'address'->>'city' as city
FROM organizations 
WHERE org_type = 'manufacturer'
ORDER BY name;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
