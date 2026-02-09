-- ============================================================================
-- Add Demo Equipment Part 2 - Remaining Manufacturers
-- ============================================================================

-- ============================================================================
-- DRÄGER MEDICAL INDIA - Add 8 more units (currently 2, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
-- Ventilators
('REG-DRG-VNT-001', 'QR-DRG-VNT-001', 'DRG-SAVINA-006001', '550e8400-e29b-41d4-a716-446655440009', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Savina 300 Ventilator', 'Dräger Medical India', 'SAVINA-300', 'Ventilator',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'ICU Ventilator Bay 1', '2024-01-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-001', 'system'),

('REG-DRG-VNT-002', 'QR-DRG-VNT-002', 'DRG-SAVINA-006002', '550e8400-e29b-41d4-a716-446655440009', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Savina 300 Ventilator', 'Dräger Medical India', 'SAVINA-300', 'Ventilator',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'ICU Ventilator Bay 2', '2024-01-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-002', 'system'),

('REG-DRG-VNT-003', 'QR-DRG-VNT-003', 'DRG-SAVINA-006003', '550e8400-e29b-41d4-a716-446655440009', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Savina 300 Ventilator', 'Dräger Medical India', 'SAVINA-300', 'Ventilator',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'ICU Ventilator Bay 1', '2024-02-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-003', 'system'),

('REG-DRG-VNT-004', 'QR-DRG-VNT-004', 'DRG-SAVINA-006004', '550e8400-e29b-41d4-a716-446655440009', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Savina 300 Ventilator', 'Dräger Medical India', 'SAVINA-300', 'Ventilator',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'ICU Ventilator Bay 2', '2024-02-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-004', 'system'),

('REG-DRG-VNT-005', 'QR-DRG-VNT-005', 'DRG-SAVINA-006005', '550e8400-e29b-41d4-a716-446655440009', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Savina 300 Ventilator', 'Dräger Medical India', 'SAVINA-300', 'Ventilator',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'ICU Ventilator Bay 1', '2024-03-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-VNT-005', 'system'),

-- Anesthesia Workstations
('REG-DRG-ANS-001', 'QR-DRG-ANS-001', 'DRG-PRIMUS-007001', '550e8400-e29b-41d4-a716-446655440010', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Primus Anesthesia Workstation', 'Dräger Medical India', 'PRIMUS-IE', 'Anesthesia',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'OT 1', '2023-11-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-ANS-001', 'system'),

('REG-DRG-ANS-002', 'QR-DRG-ANS-002', 'DRG-PRIMUS-007002', '550e8400-e29b-41d4-a716-446655440010', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Primus Anesthesia Workstation', 'Dräger Medical India', 'PRIMUS-IE', 'Anesthesia',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'OT 2', '2023-12-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-ANS-002', 'system'),

('REG-DRG-ANS-003', 'QR-DRG-ANS-003', 'DRG-PRIMUS-007003', '550e8400-e29b-41d4-a716-446655440010', 'd9e4a5b6-7890-4bcd-ef01-234567890bcd',
 'Primus Anesthesia Workstation', 'Dräger Medical India', 'PRIMUS-IE', 'Anesthesia',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'OT 3', '2024-01-20', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-DRG-ANS-003', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- FRESENIUS MEDICAL CARE - Add 8 more units (currently 2, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
('REG-FMC-DLY-001', 'QR-FMC-DLY-001', 'FMC-5008-008001', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Dialysis Center Station 1', '2023-10-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-001', 'system'),

('REG-FMC-DLY-002', 'QR-FMC-DLY-002', 'FMC-5008-008002', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Dialysis Center Station 2', '2023-10-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-002', 'system'),

('REG-FMC-DLY-003', 'QR-FMC-DLY-003', 'FMC-5008-008003', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-AIIMS-001', 'AIIMS New Delhi', 'Dialysis Center Station 3', '2023-10-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-003', 'system'),

('REG-FMC-DLY-004', 'QR-FMC-DLY-004', 'FMC-5008-008004', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Dialysis Center Station 1', '2023-11-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-004', 'system'),

('REG-FMC-DLY-005', 'QR-FMC-DLY-005', 'FMC-5008-008005', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Dialysis Center Station 2', '2023-11-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-005', 'system'),

('REG-FMC-DLY-006', 'QR-FMC-DLY-006', 'FMC-5008-008006', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-YASHODA-001', 'Yashoda Hospitals Hyderabad', 'Dialysis Center Station 3', '2023-11-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-006', 'system'),

('REG-FMC-DLY-007', 'QR-FMC-DLY-007', 'FMC-5008-008007', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'Dialysis Center Station 1', '2024-01-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-007', 'system'),

('REG-FMC-DLY-008', 'QR-FMC-DLY-008', 'FMC-5008-008008', '550e8400-e29b-41d4-a716-446655440011', 'e0f5b6c7-8901-4cde-f012-345678901cde',
 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care India', '5008S', 'Dialysis',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'Dialysis Center Station 2', '2024-01-10', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-FMC-DLY-008', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- CANON MEDICAL SYSTEMS - Add 9 more units (currently 1, total 10)
-- ============================================================================

INSERT INTO equipment_registry (
    id, qr_code, serial_number, equipment_catalog_id, manufacturer_id,
    equipment_name, manufacturer_name, model_number, category,
    customer_id, customer_name, installation_location,
    installation_date, status, qr_code_url, created_by
) VALUES
('REG-CAN-XR-001', 'QR-CAN-XR-001', 'CAN-CXDI-009001', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-AARTHI-001', 'Aarthi Scans Chennai', 'X-Ray Room 1', '2023-09-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-001', 'system'),

('REG-CAN-XR-002', 'QR-CAN-XR-002', 'CAN-CXDI-009002', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-AARTHI-001', 'Aarthi Scans Chennai', 'X-Ray Room 2', '2023-09-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-002', 'system'),

('REG-CAN-XR-003', 'QR-CAN-XR-003', 'CAN-CXDI-009003', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-SRL-001', 'SRL Diagnostics Imaging', 'X-Ray Room 1', '2023-10-20', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-003', 'system'),

('REG-CAN-XR-004', 'QR-CAN-XR-004', 'CAN-CXDI-009004', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-SRL-001', 'SRL Diagnostics Imaging', 'X-Ray Room 2', '2023-10-20', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-004', 'system'),

('REG-CAN-XR-005', 'QR-CAN-XR-005', 'CAN-CXDI-009005', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-VIJAYA-001', 'Vijaya Diagnostic Centre', 'X-Ray Room 1', '2023-11-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-005', 'system'),

('REG-CAN-XR-006', 'QR-CAN-XR-006', 'CAN-CXDI-009006', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-VIJAYA-001', 'Vijaya Diagnostic Centre', 'X-Ray Room 2', '2023-11-05', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-006', 'system'),

('REG-CAN-XR-007', 'QR-CAN-XR-007', 'CAN-CXDI-009007', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-FORTIS-001', 'Fortis Hospital Mumbai', 'Emergency X-Ray', '2023-12-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-007', 'system'),

('REG-CAN-XR-008', 'QR-CAN-XR-008', 'CAN-CXDI-009008', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-MANIPAL-001', 'Manipal Hospitals Bengaluru', 'X-Ray Room', '2024-01-15', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-008', 'system'),

('REG-CAN-XR-009', 'QR-CAN-XR-009', 'CAN-CXDI-009009', '550e8400-e29b-41d4-a716-446655440005', 'c8d3f4e5-6789-4abc-def0-123456789abc',
 'Digital X-Ray System CXDI-410C', 'Canon Medical Systems India', 'CXDI-410C', 'X-Ray',
 'CUST-APOLLO-001', 'Apollo Hospitals Chennai', 'Portable X-Ray Unit', '2024-02-01', 'operational',
 'https://api.qrserver.com/v1/create-qr-code/?data=QR-CAN-XR-009', 'system')
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- FINAL VERIFICATION
-- ============================================================================

SELECT 
    o.name as manufacturer,
    COUNT(er.id) as total_equipment,
    COUNT(DISTINCT er.customer_id) as unique_customers
FROM organizations o
LEFT JOIN equipment_registry er ON er.manufacturer_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.id, o.name
ORDER BY total_equipment DESC;

SELECT '✅ Demo equipment complete for all manufacturers!' as status,
       COUNT(*) as total_equipment
FROM equipment_registry
WHERE manufacturer_id IS NOT NULL;
