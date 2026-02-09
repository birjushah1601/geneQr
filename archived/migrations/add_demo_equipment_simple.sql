-- ============================================================================
-- Add Demo Equipment for Manufacturers (Simple Version)
-- ============================================================================
-- Adds sample equipment for each manufacturer using existing schema constraints

-- ============================================================================
-- SIEMENS HEALTHINEERS INDIA - Add 6 more units (currently 4, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
-- MRI Scanners
('REG-SIE-MRI-002', 'QR-SIE-MRI-002', 'SIE-VIDA-001002', '550e8400-e29b-41d4-a716-446655440001', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'MAGNETOM Vida 3T MRI Scanner', 'Siemens Healthineers India', 'VIDA-3T', 'MRI',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Radiology Dept - Floor 3', '2024-01-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-MRI-002', 'system'),

('REG-SIE-MRI-003', 'QR-SIE-MRI-003', 'SIE-VIDA-001003', '550e8400-e29b-41d4-a716-446655440001', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'MAGNETOM Vida 3T MRI Scanner', 'Siemens Healthineers India', 'VIDA-3T', 'MRI',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'MRI Suite 2', '2024-02-20', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-MRI-003', 'system'),

('REG-SIE-MRI-004', 'QR-SIE-MRI-004', 'SIE-SKYRA-001001', 'b2c0fde3-40bd-45ae-8422-6cfb7d9165e5', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'MAGNETOM Skyra 3T', 'Siemens Healthineers India', 'SKYRA-3T', 'MRI',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'MRI Center', '2024-03-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-MRI-004', 'system'),

-- CT Scanners
('REG-SIE-CT-001', 'QR-SIE-CT-001', 'SIE-SOMATOM-002001', 'e51ab917-09f3-4112-8b2d-3f80e1d69366', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'SOMATOM Definition AS', 'Siemens Healthineers India', 'AS-64', 'CT',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'CT Suite 1', '2023-11-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-CT-001', 'system'),

('REG-SIE-CT-002', 'QR-SIE-CT-002', 'SIE-SOMATOM-002002', 'e51ab917-09f3-4112-8b2d-3f80e1d69366', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'SOMATOM Definition AS', 'Siemens Healthineers India', 'AS-64', 'CT',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'CT Suite 2', '2023-12-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-CT-002', 'system'),

-- X-Ray
('REG-SIE-XR-001', 'QR-SIE-XR-001', 'SIE-MULTIX-003001', 'd2e959d0-4215-4397-b115-a11ce5027ef9', '11afdeec-5dee-44d4-aa5b-952703536f10',
 'Multix Fusion', 'Siemens Healthineers India', 'MF-MAX', 'X-Ray',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'X-Ray Room 1', '2024-04-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-SIE-XR-001', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- WIPRO GE HEALTHCARE - Add 4 more units (currently 6, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
-- Ultrasound Systems
('REG-WGE-US-001', 'QR-WGE-US-001', 'WGE-LOGIQ-003001', '550e8400-e29b-41d4-a716-446655440007', 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad',
 'LOGIQ E10 Ultrasound', 'Wipro GE Healthcare', 'E10-BT20', 'Ultrasound',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'OB/GYN Room 101', '2024-01-20', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-WGE-US-001', 'system'),

('REG-WGE-US-002', 'QR-WGE-US-002', 'WGE-LOGIQ-003002', '550e8400-e29b-41d4-a716-446655440007', 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad',
 'LOGIQ E10 Ultrasound', 'Wipro GE Healthcare', 'E10-BT20', 'Ultrasound',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'Ultrasound Suite', '2024-02-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-WGE-US-002', 'system'),

('REG-WGE-US-003', 'QR-WGE-US-003', 'WGE-VOLUSON-003003', '7d17a501-830b-4698-bf87-34327b770a37', 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad',
 'Voluson E10', 'Wipro GE Healthcare', 'E10-BT20', 'Ultrasound',
 'CUST-SRL-001', 'SRL Diagnostics Imaging', 'Imaging Center', '2024-03-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-WGE-US-003', 'system'),

-- MRI
('REG-WGE-MRI-001', 'QR-WGE-MRI-001', 'WGE-OPTIMA-004001', '7a571a32-d0d8-4b07-85b5-14ae66fdb0ac', 'aa0cbe3a-7e35-4cc9-88f8-2dcfdc0909ad',
 'Optima MR450w', 'Wipro GE Healthcare', 'MR450W-1.5T', 'MRI',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'MRI Suite 3', '2024-01-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-WGE-MRI-001', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- PHILIPS HEALTHCARE INDIA - Add 7 more units (currently 3, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
-- Patient Monitors
('REG-PHI-PM-001', 'QR-PHI-PM-001', 'PHI-MX850-004001', 'c6b31de7-e369-46c7-9061-c28414287bca', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'IntelliVue MX850', 'Philips Healthcare India', 'MX850-ICU', 'Patient Monitor',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'ICU Bed 1', '2024-01-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-PM-001', 'system'),

('REG-PHI-PM-002', 'QR-PHI-PM-002', 'PHI-MX850-004002', 'c6b31de7-e369-46c7-9061-c28414287bca', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'IntelliVue MX850', 'Philips Healthcare India', 'MX850-ICU', 'Patient Monitor',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'ICU Bed 2', '2024-01-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-PM-002', 'system'),

('REG-PHI-PM-003', 'QR-PHI-PM-003', 'PHI-MX850-004003', 'c6b31de7-e369-46c7-9061-c28414287bca', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'IntelliVue MX850', 'Philips Healthcare India', 'MX850-ICU', 'Patient Monitor',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'ICU Bed 1', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-PM-003', 'system'),

('REG-PHI-PM-004', 'QR-PHI-PM-004', 'PHI-MX850-004004', 'c6b31de7-e369-46c7-9061-c28414287bca', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'IntelliVue MX850', 'Philips Healthcare India', 'MX850-ICU', 'Patient Monitor',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'ICU Bed 2', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-PM-004', 'system'),

('REG-PHI-PM-005', 'QR-PHI-PM-005', 'PHI-MX850-004005', 'c6b31de7-e369-46c7-9061-c28414287bca', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'IntelliVue MX850', 'Philips Healthcare India', 'MX850-ICU', 'Patient Monitor',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'ICU Bed 1', '2024-03-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-PM-005', 'system'),

-- MRI
('REG-PHI-MRI-001', 'QR-PHI-MRI-001', 'PHI-INGENIA-005001', '91ad0a66-e900-4918-93a0-b9684068e7e3', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'Ingenia 1.5T', 'Philips Healthcare India', 'ING-1.5T-S', 'MRI',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'MRI Room', '2023-12-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-MRI-001', 'system'),

-- Ultrasound
('REG-PHI-US-001', 'QR-PHI-US-001', 'PHI-EPIQ-005002', '25b91ef2-8d5e-4d9b-8935-a2adac8b2282', 'f1c1ebfb-57fd-4307-93db-2f72e9d004ad',
 'EPIQ Elite', 'Philips Healthcare India', 'EPIQ-7C', 'Ultrasound',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'Cardiology Dept', '2024-02-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-PHI-US-001', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- MEDTRONIC INDIA - Add 8 more units (currently 2, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
('REG-MDT-PM-001', 'QR-MDT-PM-001', 'MDT-VISION-005001', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Ward A Bed 1', '2024-01-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-001', 'system'),

('REG-MDT-PM-002', 'QR-MDT-PM-002', 'MDT-VISION-005002', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Ward A Bed 2', '2024-01-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-002', 'system'),

('REG-MDT-PM-003', 'QR-MDT-PM-003', 'MDT-VISION-005003', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Ward A Bed 3', '2024-01-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-003', 'system'),

('REG-MDT-PM-004', 'QR-MDT-PM-004', 'MDT-VISION-005004', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'Ward B Bed 1', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-004', 'system'),

('REG-MDT-PM-005', 'QR-MDT-PM-005', 'MDT-VISION-005005', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'Ward B Bed 2', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-005', 'system'),

('REG-MDT-PM-006', 'QR-MDT-PM-006', 'MDT-VISION-005006', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'Ward B Bed 3', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-006', 'system'),

('REG-MDT-PM-007', 'QR-MDT-PM-007', 'MDT-VISION-005007', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Ward C Bed 1', '2024-03-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-007', 'system'),

('REG-MDT-PM-008', 'QR-MDT-PM-008', 'MDT-VISION-005008', '550e8400-e29b-41d4-a716-446655440008', 'f1a6b7c8-9012-4def-0123-456789012def',
 'Patient Monitor Visionary', 'Medtronic India', 'VISION-ICU', 'Patient Monitor',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Ward C Bed 2', '2024-03-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-MDT-PM-008', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- CONTINUE IN NEXT PART...
-- ============================================================================

SELECT '✅ Part 1 complete - Added equipment for Siemens, Wipro GE, Philips, Medtronic' as status;
