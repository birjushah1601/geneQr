-- ============================================================================
-- Data Correction: Migrate to Proper Equipment Architecture
-- ============================================================================
-- This script:
-- 1. Creates equipment_registry entries for all catalog items
-- 2. Updates existing tickets to use equipment_registry IDs
-- 3. Verifies all tickets show parts correctly
-- ============================================================================

BEGIN;

-- Step 1: Create equipment_registry entries for all catalog items (with parts)
-- We'll create 2-3 installed units per catalog item at different locations

INSERT INTO equipment_registry (
  id, 
  equipment_catalog_id, 
  qr_code,
  serial_number, 
  equipment_name,
  manufacturer_name,
  customer_name,
  installation_location,
  installation_date,
  warranty_expiry,
  status,
  qr_code_url,
  created_by
)
SELECT 
  v.id,
  v.equipment_catalog_id_text::uuid,
  v.qr_code,
  v.serial_number,
  v.equipment_name,
  v.manufacturer_name,
  v.customer_name,
  v.installation_location,
  v.installation_date::date,
  v.warranty_expiry::date,
  v.status,
  v.qr_code_url,
  v.created_by
FROM (VALUES
  -- MRI Scanners
  ('REG-MRI-VIDA-001', '550e8400-e29b-41d4-a716-446655440001', 'QR-MRI-VIDA-001', 'SN-MRI-VIDA-001', 'MAGNETOM Vida 3T MRI Scanner', 'Siemens Healthineers', 'Apollo Hospital', 'Apollo Hospital - Imaging Wing', '2023-01-15', '2026-01-15', 'operational', 'https://app.example.com/qr/REG-MRI-VIDA-001', 'admin'),
  ('REG-MRI-VIDA-002', '550e8400-e29b-41d4-a716-446655440001', 'QR-MRI-VIDA-002', 'SN-MRI-VIDA-002', 'MAGNETOM Vida 3T MRI Scanner', 'Siemens Healthineers', 'Fortis Hospital', 'Fortis Hospital - Radiology Dept', '2023-03-20', '2026-03-20', 'operational', 'https://app.example.com/qr/REG-MRI-VIDA-002', 'admin'),
  ('REG-MRI-SIGNA-001', '550e8400-e29b-41d4-a716-446655440002', 'QR-MRI-SIGNA-001', 'SN-MRI-SIGNA-001', 'SIGNA Explorer 1.5T', 'GE Healthcare', 'Max Healthcare', 'Max Healthcare - MRI Center', '2022-11-10', '2025-11-10', 'operational', 'https://app.example.com/qr/REG-MRI-SIGNA-001', 'admin'),
  
  -- CT Scanners
  ('REG-CT-NOVA-001', '550e8400-e29b-41d4-a716-446655440004', 'QR-CT-NOVA-001', 'SN-CT-NOVA-001', 'CT Scanner Nova', 'Wipro GE Healthcare', 'AIIMS Delhi', 'AIIMS - CT Imaging Suite', '2023-05-12', '2026-05-12', 'operational', 'https://app.example.com/qr/REG-CT-NOVA-001', 'admin'),
  ('REG-CT-NOVA-002', '550e8400-e29b-41d4-a716-446655440004', 'QR-CT-NOVA-002', 'SN-CT-NOVA-002', 'CT Scanner Nova', 'Wipro GE Healthcare', 'Medanta Hospital', 'Medanta Hospital - Emergency Wing', '2023-07-18', '2026-07-18', 'operational', 'https://app.example.com/qr/REG-CT-NOVA-002', 'admin'),
  ('REG-CT-PHIL-001', '550e8400-e29b-41d4-a716-446655440003', 'QR-CT-PHIL-001', 'SN-CT-PHIL-001', 'Ingenuity CT 128-slice', 'Philips Healthcare', 'Columbia Asia', 'Columbia Asia - Diagnostics Center', '2022-09-25', '2025-09-25', 'operational', 'https://app.example.com/qr/REG-CT-PHIL-001', 'admin'),
  
  -- X-Ray Systems
  ('REG-XR-ALPHA-001', '550e8400-e29b-41d4-a716-446655440006', 'QR-XR-ALPHA-001', 'SN-XR-ALPHA-001', 'X-Ray System Alpha', 'SouthCare Distributors', 'Local Clinic', 'Local Clinic - X-Ray Room', '2023-02-14', '2026-02-14', 'operational', 'https://app.example.com/qr/REG-XR-ALPHA-001', 'admin'),
  ('REG-XR-ALPHA-002', '550e8400-e29b-41d4-a716-446655440006', 'QR-XR-ALPHA-002', 'SN-XR-ALPHA-002', 'X-Ray System Alpha', 'SouthCare Distributors', 'Community Hospital', 'Community Hospital - Radiology', '2023-04-22', '2026-04-22', 'operational', 'https://app.example.com/qr/REG-XR-ALPHA-002', 'admin'),
  ('REG-XR-CANON-001', '550e8400-e29b-41d4-a716-446655440005', 'QR-XR-CANON-001', 'SN-XR-CANON-001', 'Digital X-Ray System CXDI-410C', 'Canon Medical', 'Regional Medical Center', 'Regional Medical Center', '2022-12-05', '2025-12-05', 'operational', 'https://app.example.com/qr/REG-XR-CANON-001', 'admin'),
  
  -- Ultrasound
  ('REG-US-LOGIQ-001', '550e8400-e29b-41d4-a716-446655440007', 'QR-US-LOGIQ-001', 'SN-US-LOGIQ-001', 'LOGIQ E10 Ultrasound', 'GE Healthcare', 'Maternity Hospital', 'Maternity Hospital - OB/GYN', '2023-06-30', '2026-06-30', 'operational', 'https://app.example.com/qr/REG-US-LOGIQ-001', 'admin'),
  ('REG-US-LOGIQ-002', '550e8400-e29b-41d4-a716-446655440007', 'QR-US-LOGIQ-002', 'SN-US-LOGIQ-002', 'LOGIQ E10 Ultrasound', 'GE Healthcare', 'Primary Health Center', 'Primary Health Center', '2023-08-15', '2026-08-15', 'operational', 'https://app.example.com/qr/REG-US-LOGIQ-002', 'admin'),
  
  -- Patient Monitors
  ('REG-PM-VIS-001', '550e8400-e29b-41d4-a716-446655440008', 'QR-PM-VIS-001', 'SN-PM-VIS-001', 'Patient Monitor Visionary', 'Medtronic India', 'City Hospital', 'ICU Ward A - Bed 5', '2023-01-20', '2026-01-20', 'operational', 'https://app.example.com/qr/REG-PM-VIS-001', 'admin'),
  ('REG-PM-VIS-002', '550e8400-e29b-41d4-a716-446655440008', 'QR-PM-VIS-002', 'SN-PM-VIS-002', 'Patient Monitor Visionary', 'Medtronic India', 'Emergency Hospital', 'Emergency Department - Station 3', '2023-03-10', '2026-03-10', 'operational', 'https://app.example.com/qr/REG-PM-VIS-002', 'admin'),
  
  -- Ventilators
  ('REG-VENT-SAV-001', '550e8400-e29b-41d4-a716-446655440009', 'QR-VENT-SAV-001', 'SN-VENT-SAV-001', 'Savina 300 Ventilator', 'Dräger Medical', 'Metro Hospital', 'ICU - Critical Care Unit', '2023-02-28', '2026-02-28', 'operational', 'https://app.example.com/qr/REG-VENT-SAV-001', 'admin'),
  ('REG-VENT-SAV-002', '550e8400-e29b-41d4-a716-446655440009', 'QR-VENT-SAV-002', 'SN-VENT-SAV-002', 'Savina 300 Ventilator', 'Dräger Medical', 'Care Hospital', 'Respiratory Care Unit', '2023-04-05', '2026-04-05', 'operational', 'https://app.example.com/qr/REG-VENT-SAV-002', 'admin'),
  
  -- Dialysis
  ('REG-DIAL-FRES-001', '550e8400-e29b-41d4-a716-446655440011', 'QR-DIAL-FRES-001', 'SN-DIAL-FRES-001', 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care', 'Nephrology Center', 'Nephrology Department - Station 1', '2023-05-18', '2026-05-18', 'operational', 'https://app.example.com/qr/REG-DIAL-FRES-001', 'admin'),
  ('REG-DIAL-FRES-002', '550e8400-e29b-41d4-a716-446655440011', 'QR-DIAL-FRES-002', 'SN-DIAL-FRES-002', 'Fresenius 5008 Dialysis Machine', 'Fresenius Medical Care', 'Dialysis Center', 'Dialysis Center - Bay 4', '2023-07-22', '2026-07-22', 'operational', 'https://app.example.com/qr/REG-DIAL-FRES-002', 'admin'),
  
  -- Infusion Pumps
  ('REG-INF-LITE-001', '550e8400-e29b-41d4-a716-446655440012', 'QR-INF-LITE-001', 'SN-INF-LITE-001', 'Infusion Pump Lite', 'Philips Healthcare India', 'General Hospital', 'General Ward - Nursing Station', '2023-03-15', '2026-03-15', 'operational', 'https://app.example.com/qr/REG-INF-LITE-001', 'admin'),
  ('REG-INF-LITE-002', '550e8400-e29b-41d4-a716-446655440012', 'QR-INF-LITE-002', 'SN-INF-LITE-002', 'Infusion Pump Lite', 'Philips Healthcare India', 'Childrens Hospital', 'Pediatric Ward - Room 203', '2023-05-25', '2026-05-25', 'operational', 'https://app.example.com/qr/REG-INF-LITE-002', 'admin')
) AS v(id, equipment_catalog_id_text, qr_code, serial_number, equipment_name, manufacturer_name, customer_name, installation_location, installation_date, warranty_expiry, status, qr_code_url, created_by)
WHERE NOT EXISTS (SELECT 1 FROM equipment_registry WHERE equipment_registry.id = v.id)
RETURNING id, equipment_name, serial_number;

-- Step 2: Update existing tickets to use new equipment_registry entries
-- First, let's see which tickets exist and update them to valid registry IDs

-- Update tickets to use appropriate equipment types
UPDATE service_tickets 
SET 
  equipment_id = 'REG-VENT-SAV-001',
  equipment_name = 'Savina 300 Ventilator'
WHERE equipment_name LIKE '%Ventilator%' 
  AND equipment_id NOT IN (SELECT id FROM equipment_registry)
RETURNING id, ticket_number, equipment_name;

UPDATE service_tickets 
SET 
  equipment_id = 'REG-PM-VIS-001',
  equipment_name = 'Patient Monitor Visionary'
WHERE equipment_name LIKE '%ECG%' OR equipment_name LIKE '%Monitor%'
  AND equipment_id NOT IN (SELECT id FROM equipment_registry)
RETURNING id, ticket_number, equipment_name;

UPDATE service_tickets 
SET 
  equipment_id = 'REG-MRI-VIDA-001',
  equipment_name = 'MAGNETOM Vida 3T MRI Scanner'
WHERE equipment_name LIKE '%MRI%' 
  AND equipment_id NOT IN (SELECT id FROM equipment_registry)
RETURNING id, ticket_number, equipment_name;

UPDATE service_tickets 
SET 
  equipment_id = 'REG-CT-NOVA-001',
  equipment_name = 'CT Scanner Nova'
WHERE equipment_name LIKE '%CT%' 
  AND equipment_id NOT IN (SELECT id FROM equipment_registry)
RETURNING id, ticket_number, equipment_name;

-- Step 3: Verify all tickets now have parts
SELECT 
  '=== VERIFICATION: Tickets with Parts ===' as verification_step;

SELECT 
  t.ticket_number,
  t.equipment_name,
  er.serial_number,
  ec.product_name as catalog_product,
  COUNT(epa.id) as parts_available
FROM service_tickets t
LEFT JOIN equipment_registry er ON t.equipment_id = er.id
LEFT JOIN equipment_catalog ec ON er.equipment_catalog_id = ec.id
LEFT JOIN equipment_part_assignments epa ON epa.equipment_catalog_id = ec.id
GROUP BY t.ticket_number, t.equipment_name, er.serial_number, ec.product_name
ORDER BY parts_available DESC, t.ticket_number;

-- Step 4: Show sample parts for verification
SELECT 
  '=== SAMPLE: Parts for First Ticket ===' as sample_check;

SELECT * FROM get_parts_for_ticket(
  (SELECT id FROM service_tickets LIMIT 1)
) LIMIT 5;

-- Step 5: Summary statistics
SELECT 
  '=== SUMMARY STATISTICS ===' as summary;

SELECT 
  'Equipment Registry Entries' as metric,
  COUNT(*) as count
FROM equipment_registry
UNION ALL
SELECT 
  'Catalog Items with Parts',
  COUNT(DISTINCT equipment_catalog_id)
FROM equipment_part_assignments
UNION ALL
SELECT 
  'Total Part Assignments',
  COUNT(*)
FROM equipment_part_assignments
UNION ALL
SELECT 
  'Service Tickets',
  COUNT(*)
FROM service_tickets
UNION ALL
SELECT 
  'Tickets Using Registry',
  COUNT(*)
FROM service_tickets t
JOIN equipment_registry er ON t.equipment_id = er.id;

COMMIT;

-- Final message
SELECT 'Data correction complete! All tickets should now show parts.' as status;
