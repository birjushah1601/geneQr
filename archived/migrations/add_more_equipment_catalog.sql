-- ============================================================================
-- Add More Equipment to Catalog for Each Manufacturer
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
    'SKYRA-3T', 'MRI', '3 Tesla', 45000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'SIE-CT-001', 'SOMATOM Definition AS',
    '11afdeec-5dee-44d4-aa5b-952703536f10', 'Siemens Healthineers',
    'AS-64', 'CT', '64-slice', 35000000.00, 'INR',
    180, 12, 'high', true
),
(
    gen_random_uuid(), 'SIE-XRAY-001', 'Multix Fusion',
    '11afdeec-5dee-44d4-aa5b-952703536f10', 'Siemens Healthineers',
    'MF-MAX', 'X-Ray', 'Digital Radiography', 8500000.00, 'INR',
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
    'MR450W-1.5T', 'MRI', '1.5 Tesla', 38000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'WGE-US-001', 'Voluson E10',
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad', 'Wipro GE Healthcare',
    'E10-BT20', 'Ultrasound', '4D Imaging', 12500000.00, 'INR',
    180, 10, 'medium', true
),
(
    gen_random_uuid(), 'WGE-XRAY-001', 'Brivo XR575',
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad', 'Wipro GE Healthcare',
    'XR575-PLUS', 'X-Ray', 'Digital', 7500000.00, 'INR',
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
    'ING-1.5T-S', 'MRI', '1.5 Tesla', 42000000.00, 'INR',
    180, 15, 'high', true
),
(
    gen_random_uuid(), 'PHI-US-001', 'EPIQ Elite',
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad', 'Philips Healthcare India',
    'EPIQ-7C', 'Ultrasound', 'Premium Cardiac', 18500000.00, 'INR',
    180, 12, 'medium', true
),
(
    gen_random_uuid(), 'PHI-PM-001', 'IntelliVue MX850',
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad', 'Philips Healthcare India',
    'MX850-ICU', 'Patient Monitor', 'ICU Monitor', 850000.00, 'INR',
    90, 8, 'low', true
);

-- Show final counts
SELECT 
    o.name as manufacturer,
    COUNT(DISTINCT ec.id) as equipment_catalog_count
FROM organizations o
LEFT JOIN equipment_catalog ec ON ec.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY o.name;
