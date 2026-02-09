-- ============================================================================
-- Add Demo Equipment for Each Manufacturer
-- ============================================================================
-- This migration adds sample equipment_registry entries for each manufacturer
-- to enable full demo workflow: Equipment → QR Codes → Service Tickets

-- ============================================================================
-- 1. HELPER FUNCTION TO GENERATE QR CODES
-- ============================================================================

CREATE OR REPLACE FUNCTION generate_qr_code(equipment_id TEXT) 
RETURNS TEXT AS $$
BEGIN
    RETURN 'QR-' || UPPER(SUBSTRING(MD5(equipment_id) FROM 1 FOR 12));
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- 2. SIEMENS HEALTHINEERS INDIA - Add 6 more units (total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-SIE-MRI-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-SIE-MRI-' || LPAD(seq::TEXT, 3, '0')),
    'SIE-VIDA-' || LPAD((1000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440001'::uuid, -- MAGNETOM Vida 3T
    '11afdeec-5dee-44d4-aa5b-952703536f10'::uuid, -- Siemens
    'MAGNETOM Vida 3T MRI Scanner',
    'Siemens Healthineers India',
    'VIDA-3T',
    'MRI',
    customer_id,
    customer_name,
    'Radiology Department - Floor ' || floor_num,
    CURRENT_DATE - (seq * 30)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-SIE-MRI-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE seq
            WHEN 1 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
            WHEN 2 THEN '9305562f-2129-4ff5-a1e3-913bdff39b52'::VARCHAR -- Apollo
            WHEN 3 THEN '8002dec5-340b-4c8e-a28f-9dd80d151944'::VARCHAR -- Fortis
        END as customer_id,
        CASE seq
            WHEN 1 THEN 'AIIMS New Delhi'
            WHEN 2 THEN 'Apollo Hospitals Chennai'
            WHEN 3 THEN 'Fortis Hospital Mumbai'
        END as customer_name,
        CASE seq
            WHEN 1 THEN '3'
            WHEN 2 THEN '2'
            WHEN 3 THEN '4'
        END as floor_num
    FROM generate_series(1, 3) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- Siemens CT Scanners
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-SIE-CT-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-SIE-CT-' || LPAD(seq::TEXT, 3, '0')),
    'SIE-SOMATOM-' || LPAD((2000 + seq)::TEXT, 6, '0'),
    'e51ab917-09f3-4112-8b2d-3f80e1d69366'::uuid, -- SOMATOM Definition AS
    '11afdeec-5dee-44d4-aa5b-952703536f10'::uuid, -- Siemens
    'SOMATOM Definition AS 64-slice CT',
    'Siemens Healthineers India',
    'AS-64',
    'CT',
    customer_id,
    customer_name,
    'CT Suite ' || seq,
    CURRENT_DATE - (seq * 45)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-SIE-CT-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE seq
            WHEN 1 THEN '265c476d-cc8d-4767-934a-c9d895434f2b'::VARCHAR -- Manipal
            WHEN 2 THEN 'f21f6bbc-fe27-475f-876f-d7cf7cdea278'::VARCHAR -- Yashoda
        END as customer_id,
        CASE seq
            WHEN 1 THEN 'Manipal Hospitals Bengaluru'
            WHEN 2 THEN 'Yashoda Hospitals Hyderabad'
        END as customer_name
    FROM generate_series(1, 2) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 3. WIPRO GE HEALTHCARE - Add 4 more units (total 10)
-- ============================================================================

-- GE Ultrasound Systems
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-WGE-US-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-WGE-US-' || LPAD(seq::TEXT, 3, '0')),
    'WGE-LOGIQ-' || LPAD((3000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440007'::uuid, -- LOGIQ E10
    'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad'::uuid, -- Wipro GE
    'LOGIQ E10 Ultrasound System',
    'Wipro GE Healthcare',
    'E10-BT20',
    'Ultrasound',
    customer_id,
    customer_name,
    'OB/GYN Department - Room ' || (100 + seq),
    CURRENT_DATE - (seq * 60)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-WGE-US-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE seq
            WHEN 1 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
            WHEN 2 THEN '9305562f-2129-4ff5-a1e3-913bdff39b52'::VARCHAR -- Apollo
            WHEN 3 THEN 'f66d202d-9d0d-4116-8b42-98c7cbf7d38c'::VARCHAR -- SRL Diagnostics
        END as customer_id,
        CASE seq
            WHEN 1 THEN 'AIIMS New Delhi'
            WHEN 2 THEN 'Apollo Hospitals Chennai'
            WHEN 3 THEN 'SRL Diagnostics Imaging - Mumbai'
        END as customer_name
    FROM generate_series(1, 3) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 4. PHILIPS HEALTHCARE INDIA - Add 7 more units (total 10)
-- ============================================================================

-- Philips Patient Monitors
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-PHI-PM-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-PHI-PM-' || LPAD(seq::TEXT, 3, '0')),
    'PHI-MX850-' || LPAD((4000 + seq)::TEXT, 6, '0'),
    'c6b31de7-e369-46c7-9061-c28414287bca'::uuid, -- IntelliVue MX850
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad'::uuid, -- Philips
    'IntelliVue MX850 Patient Monitor',
    'Philips Healthcare India',
    'MX850-ICU',
    'Patient Monitor',
    customer_id,
    customer_name,
    'ICU - Bed ' || seq,
    CURRENT_DATE - (seq * 90)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-PHI-PM-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE 
            WHEN seq <= 2 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
            WHEN seq <= 4 THEN '8002dec5-340b-4c8e-a28f-9dd80d151944'::VARCHAR -- Fortis
            ELSE '265c476d-cc8d-4767-934a-c9d895434f2b'::VARCHAR -- Manipal
        END as customer_id,
        CASE 
            WHEN seq <= 2 THEN 'AIIMS New Delhi'
            WHEN seq <= 4 THEN 'Fortis Hospital Mumbai'
            ELSE 'Manipal Hospitals Bengaluru'
        END as customer_name
    FROM generate_series(1, 5) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- Philips Ultrasound
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
(
    'REG-PHI-US-001',
    generate_qr_code('REG-PHI-US-001'),
    'PHI-EPIQ-005001',
    '25b91ef2-8d5e-4d9b-8935-a2adac8b2282'::uuid, -- EPIQ Elite
    'f1c1ebfb-57fd-4307-93db-2f72e9d004ad'::uuid, -- Philips
    'EPIQ Elite Ultrasound System',
    'Philips Healthcare India',
    'EPIQ-7C',
    'Ultrasound',
    '9305562f-2129-4ff5-a1e3-913bdff39b52', -- Apollo
    'Apollo Hospitals Chennai',
    'Cardiology Department',
    CURRENT_DATE - 120,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-PHI-US-001',
    'system'
)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 5. MEDTRONIC INDIA - Add 8 more units (total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-MDT-PM-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-MDT-PM-' || LPAD(seq::TEXT, 3, '0')),
    'MDT-VISION-' || LPAD((5000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440008'::uuid, -- Patient Monitor Visionary
    'f1a6b7c8-9012-4def-0123-456789012def'::uuid, -- Medtronic
    'Patient Monitor Visionary',
    'Medtronic India',
    'VISION-ICU-PRO',
    'Patient Monitor',
    customer_id,
    customer_name,
    'Ward ' || CASE WHEN seq <= 3 THEN 'A' WHEN seq <= 6 THEN 'B' ELSE 'C' END || ' - Bed ' || ((seq - 1) % 3 + 1),
    CURRENT_DATE - (seq * 75)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-MDT-PM-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE 
            WHEN seq <= 3 THEN 'f21f6bbc-fe27-475f-876f-d7cf7cdea278'::VARCHAR -- Yashoda
            WHEN seq <= 6 THEN '265c476d-cc8d-4767-934a-c9d895434f2b'::VARCHAR -- Manipal
            ELSE 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
        END as customer_id,
        CASE 
            WHEN seq <= 3 THEN 'Yashoda Hospitals Hyderabad'
            WHEN seq <= 6 THEN 'Manipal Hospitals Bengaluru'
            ELSE 'AIIMS New Delhi'
        END as customer_name
    FROM generate_series(1, 8) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 6. DRÄGER MEDICAL INDIA - Add 8 more units (total 10)
-- ============================================================================

-- Ventilators
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-DRG-VNT-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-DRG-VNT-' || LPAD(seq::TEXT, 3, '0')),
    'DRG-SAVINA-' || LPAD((6000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440009'::uuid, -- Savina 300
    'd9e4a5b6-7890-4bcd-ef01-234567890bcd'::uuid, -- Dräger
    'Savina 300 Ventilator',
    'Dräger Medical India',
    'SAVINA-300',
    'Ventilator',
    customer_id,
    customer_name,
    'ICU - Ventilator Bay ' || seq,
    CURRENT_DATE - (seq * 50)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-DRG-VNT-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE 
            WHEN seq <= 2 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
            WHEN seq <= 4 THEN '8002dec5-340b-4c8e-a28f-9dd80d151944'::VARCHAR -- Fortis
            ELSE 'f21f6bbc-fe27-475f-876f-d7cf7cdea278'::VARCHAR -- Yashoda
        END as customer_id,
        CASE 
            WHEN seq <= 2 THEN 'AIIMS New Delhi'
            WHEN seq <= 4 THEN 'Fortis Hospital Mumbai'
            ELSE 'Yashoda Hospitals Hyderabad'
        END as customer_name
    FROM generate_series(1, 5) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- Anesthesia Workstations
INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-DRG-ANS-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-DRG-ANS-' || LPAD(seq::TEXT, 3, '0')),
    'DRG-PRIMUS-' || LPAD((7000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440010'::uuid, -- Primus
    'd9e4a5b6-7890-4bcd-ef01-234567890bcd'::uuid, -- Dräger
    'Primus Anesthesia Workstation',
    'Dräger Medical India',
    'PRIMUS-IE',
    'Anesthesia',
    customer_id,
    customer_name,
    'Operating Theatre ' || seq,
    CURRENT_DATE - (seq * 100)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-DRG-ANS-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE seq
            WHEN 1 THEN '9305562f-2129-4ff5-a1e3-913bdff39b52'::VARCHAR -- Apollo
            WHEN 2 THEN '265c476d-cc8d-4767-934a-c9d895434f2b'::VARCHAR -- Manipal
            WHEN 3 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
        END as customer_id,
        CASE seq
            WHEN 1 THEN 'Apollo Hospitals Chennai'
            WHEN 2 THEN 'Manipal Hospitals Bengaluru'
            WHEN 3 THEN 'AIIMS New Delhi'
        END as customer_name
    FROM generate_series(1, 3) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 7. FRESENIUS MEDICAL CARE - Add 8 more units (total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-FMC-DLY-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-FMC-DLY-' || LPAD(seq::TEXT, 3, '0')),
    'FMC-5008-' || LPAD((8000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440011'::uuid, -- Fresenius 5008
    'e0f5b6c7-8901-4cde-f012-345678901cde'::uuid, -- Fresenius
    'Fresenius 5008 Dialysis Machine',
    'Fresenius Medical Care India',
    '5008S-CORDIAX',
    'Dialysis',
    customer_id,
    customer_name,
    'Dialysis Center - Station ' || seq,
    CURRENT_DATE - (seq * 80)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-FMC-DLY-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE 
            WHEN seq <= 3 THEN 'a078de20-ea2f-4f7b-a6eb-6f00e0eb66eb'::VARCHAR -- AIIMS
            WHEN seq <= 6 THEN 'f21f6bbc-fe27-475f-876f-d7cf7cdea278'::VARCHAR -- Yashoda
            ELSE '9305562f-2129-4ff5-a1e3-913bdff39b52'::VARCHAR -- Apollo
        END as customer_id,
        CASE 
            WHEN seq <= 3 THEN 'AIIMS New Delhi'
            WHEN seq <= 6 THEN 'Yashoda Hospitals Hyderabad'
            ELSE 'Apollo Hospitals Chennai'
        END as customer_name
    FROM generate_series(1, 8) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 8. CANON MEDICAL SYSTEMS - Add 9 more units (total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
)
SELECT 
    'REG-CAN-XR-' || LPAD(seq::TEXT, 3, '0'),
    generate_qr_code('REG-CAN-XR-' || LPAD(seq::TEXT, 3, '0')),
    'CAN-CXDI-' || LPAD((9000 + seq)::TEXT, 6, '0'),
    '550e8400-e29b-41d4-a716-446655440005'::uuid, -- CXDI-410C
    'c8d3f4e5-6789-4abc-def0-123456789abc'::uuid, -- Canon
    'Digital X-Ray System CXDI-410C',
    'Canon Medical Systems India',
    'CXDI-410C',
    'X-Ray',
    customer_id,
    customer_name,
    'X-Ray Room ' || seq,
    CURRENT_DATE - (seq * 40)::INTEGER,
    'operational',
    'https://api.qrserver.com/v1/create-qr-code/?data=REG-CAN-XR-' || LPAD(seq::TEXT, 3, '0'),
    'system'
FROM (
    SELECT 
        seq,
        CASE 
            WHEN seq <= 2 THEN '795f79cd-ee55-41b4-8259-ede84424c41a'::VARCHAR -- Aarthi Scans
            WHEN seq <= 4 THEN 'f66d202d-9d0d-4116-8b42-98c7cbf7d38c'::VARCHAR -- SRL
            WHEN seq <= 6 THEN '953706ff-f020-4f5c-93d7-3141c76aa750'::VARCHAR -- Vijaya
            ELSE '8002dec5-340b-4c8e-a28f-9dd80d151944'::VARCHAR -- Fortis
        END as customer_id,
        CASE 
            WHEN seq <= 2 THEN 'Aarthi Scans & Labs - Chennai'
            WHEN seq <= 4 THEN 'SRL Diagnostics Imaging - Mumbai'
            WHEN seq <= 6 THEN 'Vijaya Diagnostic Centre - Hyderabad'
            ELSE 'Fortis Hospital Mumbai'
        END as customer_name
    FROM generate_series(1, 9) seq
) AS customers
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- 9. VERIFY RESULTS
-- ============================================================================

SELECT 
    o.name as manufacturer,
    COUNT(er.id) as total_equipment,
    COUNT(DISTINCT er.customer_id) as unique_customers,
    MIN(er.installation_date) as oldest_install,
    MAX(er.installation_date) as newest_install
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY total_equipment DESC;

-- Show equipment distribution by category
SELECT 
    o.name as manufacturer,
    er.category,
    COUNT(*) as count
FROM organizations o
JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.name, er.category
ORDER BY o.name, count DESC;

-- Clean up helper function
DROP FUNCTION IF EXISTS generate_qr_code(TEXT);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

SELECT '✅ Demo equipment added for all manufacturers!' as result,
       COUNT(*) as total_equipment
FROM equipment_registry
WHERE manufacturer_id IS NOT NULL;
