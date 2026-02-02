-- ============================================================================
-- Link Remaining Equipment to Manufacturers
-- ============================================================================
-- Associate equipment catalog items that don't have manufacturer_id yet

-- ============================================================================
-- 1. CHECK CURRENT STATUS
-- ============================================================================

SELECT 
    'BEFORE UPDATE' as status,
    COUNT(*) as total_equipment,
    COUNT(manufacturer_id) as with_manufacturer_id,
    COUNT(*) - COUNT(manufacturer_id) as missing_manufacturer_id
FROM equipment_catalog;

-- ============================================================================
-- 2. CREATE NEW MANUFACTURER ORGANIZATIONS FOR MISSING ONES
-- ============================================================================

-- Canon Medical (Japanese manufacturer)
INSERT INTO organizations (id, name, org_type, status, metadata)
VALUES (
    'c8d3f4e5-6789-4abc-def0-123456789abc',
    'Canon Medical Systems India',
    'manufacturer',
    'active',
    jsonb_build_object(
        'contact_person', 'Mr. Takeshi Yamamoto',
        'email', 'takeshi.yamamoto@canon-medical.co.in',
        'phone', '+91-124-4819-700',
        'website', 'https://in.medical.canon',
        'address', jsonb_build_object(
            'street', 'Plot No. 11, Udyog Vihar, Phase IV',
            'city', 'Gurugram',
            'state', 'Haryana',
            'postal_code', '122015',
            'country', 'India'
        ),
        'business_info', jsonb_build_object(
            'gst_number', '06AABCC5678K1Z9',
            'pan_number', 'AABCC5678K',
            'established_year', 2016,
            'employee_count', 800,
            'headquarters', 'Gurugram, Haryana'
        ),
        'support_info', jsonb_build_object(
            'support_email', 'support.india@canon-medical.co.in',
            'support_phone', '+91-124-4819-750',
            'support_hours', 'Mon-Sat: 9AM-6PM',
            'response_time_sla', '6 hours'
        )
    )
)
ON CONFLICT (id) DO NOTHING;

-- DrÃ¤ger Medical (German manufacturer - ventilators, anesthesia)
INSERT INTO organizations (id, name, org_type, status, metadata)
VALUES (
    'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
    'DrÃ¤ger Medical India',
    'manufacturer',
    'active',
    jsonb_build_object(
        'contact_person', 'Dr. Klaus Weber',
        'email', 'klaus.weber@draeger.com',
        'phone', '+91-22-6112-2000',
        'website', 'https://www.draeger.com/en-in',
        'address', jsonb_build_object(
            'street', 'DrÃ¤ger India, 5th Floor, Tower A, Unitech Business Park',
            'city', 'Navi Mumbai',
            'state', 'Maharashtra',
            'postal_code', '400614',
            'country', 'India'
        ),
        'business_info', jsonb_build_object(
            'gst_number', '27AABCD6789L1Z1',
            'pan_number', 'AABCD6789L',
            'established_year', 2005,
            'employee_count', 600,
            'headquarters', 'Navi Mumbai, Maharashtra'
        ),
        'support_info', jsonb_build_object(
            'support_email', 'service.india@draeger.com',
            'support_phone', '+91-22-6112-2100',
            'support_hours', '24/7 Available',
            'response_time_sla', '4 hours'
        )
    )
)
ON CONFLICT (id) DO NOTHING;

-- Fresenius Medical Care (German manufacturer - dialysis)
INSERT INTO organizations (id, name, org_type, status, metadata)
VALUES (
    'e0f5b6c7-8901-4cde-f012-345678901cde',
    'Fresenius Medical Care India',
    'manufacturer',
    'active',
    jsonb_build_object(
        'contact_person', 'Mr. Ravi Chandran',
        'email', 'ravi.chandran@fmc-ag.com',
        'phone', '+91-44-4203-5000',
        'website', 'https://www.freseniusmedicalcare.asia/en-in',
        'address', jsonb_build_object(
            'street', 'Fresenius House, 3rd Floor, Tower 3, Express Towers',
            'city', 'Chennai',
            'state', 'Tamil Nadu',
            'postal_code', '600002',
            'country', 'India'
        ),
        'business_info', jsonb_build_object(
            'gst_number', '33AABCF7890M1Z2',
            'pan_number', 'AABCF7890M',
            'established_year', 2010,
            'employee_count', 500,
            'headquarters', 'Chennai, Tamil Nadu'
        ),
        'support_info', jsonb_build_object(
            'support_email', 'india.support@fmc-ag.com',
            'support_phone', '+91-44-4203-5100',
            'support_hours', 'Mon-Sat: 8AM-8PM',
            'response_time_sla', '8 hours'
        )
    )
)
ON CONFLICT (id) DO NOTHING;

-- Medtronic India (US manufacturer - patient monitors, surgical equipment)
INSERT INTO organizations (id, name, org_type, status, metadata)
VALUES (
    'f1a6b7c8-9012-4def-0123-456789012def',
    'Medtronic India',
    'manufacturer',
    'active',
    jsonb_build_object(
        'contact_person', 'Ms. Anjali Verma',
        'email', 'anjali.verma@medtronic.com',
        'phone', '+91-40-4888-3000',
        'website', 'https://www.medtronic.com/in-en',
        'address', jsonb_build_object(
            'street', 'Medtronic India, Raheja Mindspace, Building 9, 5th Floor',
            'city', 'Hyderabad',
            'state', 'Telangana',
            'postal_code', '500081',
            'country', 'India'
        ),
        'business_info', jsonb_build_object(
            'gst_number', '36AABCM8901N1Z3',
            'pan_number', 'AABCM8901N',
            'established_year', 1998,
            'employee_count', 2000,
            'headquarters', 'Hyderabad, Telangana'
        ),
        'support_info', jsonb_build_object(
            'support_email', 'india.service@medtronic.com',
            'support_phone', '+91-40-4888-3100',
            'support_hours', '24/7 Available',
            'response_time_sla', '4 hours'
        )
    )
)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 3. UPDATE EQUIPMENT CATALOG WITH MANUFACTURER IDs
-- ============================================================================

-- Canon Medical equipment
UPDATE equipment_catalog 
SET manufacturer_id = 'c8d3f4e5-6789-4abc-def0-123456789abc',
    manufacturer_name = 'Canon Medical Systems India'
WHERE manufacturer_name = 'Canon Medical';

-- DrÃ¤ger Medical equipment
UPDATE equipment_catalog 
SET manufacturer_id = 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
    manufacturer_name = 'DrÃ¤ger Medical India'
WHERE manufacturer_name ILIKE '%DrÃ¤ger%' OR manufacturer_name ILIKE '%Drager%';

-- Fresenius Medical Care equipment
UPDATE equipment_catalog 
SET manufacturer_id = 'e0f5b6c7-8901-4cde-f012-345678901cde',
    manufacturer_name = 'Fresenius Medical Care India'
WHERE manufacturer_name ILIKE '%Fresenius%';

-- Medtronic India equipment
UPDATE equipment_catalog 
SET manufacturer_id = 'f1a6b7c8-9012-4def-0123-456789012def',
    manufacturer_name = 'Medtronic India'
WHERE manufacturer_name = 'Medtronic India';

-- For Channel Partner/Sub-sub_SUB_DEALER equipment, keep them as is (no manufacturer_id)
-- These are resellers, not manufacturers:
-- - SouthCare Channel Partners (Channel Partner)
-- - Test Corp (test data)
-- - Test Medical Systems (test data)

-- ============================================================================
-- 4. VERIFY RESULTS
-- ============================================================================

-- Show updated counts
SELECT 
    'AFTER UPDATE' as status,
    COUNT(*) as total_equipment,
    COUNT(manufacturer_id) as with_manufacturer_id,
    COUNT(*) - COUNT(manufacturer_id) as missing_manufacturer_id
FROM equipment_catalog;

-- Show manufacturers with equipment counts
SELECT 
    o.name as manufacturer,
    o.org_type,
    COUNT(ec.id) as equipment_count
FROM organizations o
LEFT JOIN equipment_catalog ec ON ec.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name, o.org_type
ORDER BY equipment_count DESC, o.name;

-- Show equipment without manufacturer_id (should be Channel Partners/test data)
SELECT 
    id,
    product_name,
    manufacturer_name,
    'Reason: Not a manufacturer' as note
FROM equipment_catalog
WHERE manufacturer_id IS NULL
ORDER BY product_name;

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

SELECT 'âœ… Equipment-Manufacturer linking complete!' as result;
