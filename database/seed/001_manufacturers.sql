-- ============================================================================
-- SEED DATA: MANUFACTURERS (10 Organizations)
-- ============================================================================
-- Real-world medical equipment manufacturers operating in India

-- ============================================================================
-- 1. SIEMENS HEALTHINEERS INDIA
-- ============================================================================

-- Organization
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  legal_entity_name, registration_number, tax_id, year_established, 
  employee_count, website, logo_url, metadata, created_at)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Siemens Healthineers India',
  'Siemens Healthineers',
  'manufacturer',
  'active',
  true,
  'Siemens Healthcare Private Limited',
  'U33112MH2005FTC156346',
  'AABCS1234F',
  2005,
  2500,
  'https://www.siemens-healthineers.com/en-in',
  'https://www.siemens-healthineers.com/logo.png',
  '{"specializations": ["Diagnostic Imaging", "Laboratory Diagnostics", "Point of Care Testing"], "iso_certified": true, "ce_certified": true}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, status, operational_since)
VALUES 
-- Mumbai Manufacturing
('11111111-1111-1111-1111-111111111101',
 '11111111-1111-1111-1111-111111111111',
 'Siemens Mumbai Manufacturing Plant',
 'SIE-MUM-MFG',
 'manufacturing_plant',
 '{"line1": "Plot No. 2, MIDC Industrial Area", "city": "Mumbai", "state": "Maharashtra", "pincode": "400093", "country": "India"}',
 '{"monday_friday": "08:00-18:00", "saturday": "08:00-14:00"}',
 ARRAY['Manufacturing', 'Assembly', 'Quality Testing'],
 ARRAY['CT Scanners', 'MRI Machines', 'X-Ray Systems'],
 'active',
 '2005-06-01'),

-- Delhi Service Center
('11111111-1111-1111-1111-111111111102',
 '11111111-1111-1111-1111-111111111111',
 'Siemens Delhi Service Center',
 'SIE-DEL-SVC',
 'service_center',
 '{"line1": "A-26, Mohan Cooperative Industrial Estate", "city": "New Delhi", "state": "Delhi", "pincode": "110044", "country": "India"}',
 '{"monday_friday": "09:00-18:00", "saturday": "09:00-13:00"}',
 ARRAY['Installation', 'Maintenance', 'Calibration', 'Repair'],
 ARRAY['All Diagnostic Equipment'],
 'active',
 '2007-01-15'),

-- Bangalore R&D Center
('11111111-1111-1111-1111-111111111103',
 '11111111-1111-1111-1111-111111111111',
 'Siemens Bangalore R&D Center',
 'SIE-BLR-RND',
 'rnd_center',
 '{"line1": "RMZ Infinity, Old Madras Road", "city": "Bangalore", "state": "Karnataka", "pincode": "560016", "country": "India"}',
 '{"monday_friday": "09:00-18:00"}',
 ARRAY['Research', 'Product Development', 'Training'],
 ARRAY['Software Development', 'AI/ML Applications'],
 'active',
 '2010-03-01'),

-- Chennai Training Center
('11111111-1111-1111-1111-111111111104',
 '11111111-1111-1111-1111-111111111111',
 'Siemens Chennai Training Center',
 'SIE-CHE-TRN',
 'training_center',
 '{"line1": "Mount Poonamallee Road, Manapakkam", "city": "Chennai", "state": "Tamil Nadu", "pincode": "600089", "country": "India"}',
 '{"monday_friday": "09:00-17:00"}',
 ARRAY['Engineer Training', 'Customer Training', 'Certification Programs'],
 ARRAY['All Equipment Categories'],
 'active',
 '2012-06-01');

-- Contact Persons
INSERT INTO contact_persons (org_id, name, designation, department, email, primary_phone, 
  whatsapp_number, is_primary, can_approve_orders, can_raise_tickets, preferred_contact_method)
VALUES
('11111111-1111-1111-1111-111111111111', 'Rajesh Kumar', 'Managing Director', 'Executive', 
 'rajesh.kumar@siemens-healthineers.com', '+91-22-6715-8000', '+91-98200-12345', true, true, true, 'email'),
('11111111-1111-1111-1111-111111111111', 'Priya Sharma', 'Head of Service', 'Service', 
 'priya.sharma@siemens-healthineers.com', '+91-22-6715-8100', '+91-98200-12346', false, true, true, 'phone'),
('11111111-1111-1111-1111-111111111111', 'Amit Patel', 'Regional Service Manager - North', 'Service', 
 'amit.patel@siemens-healthineers.com', '+91-11-4567-8900', '+91-98100-12347', false, false, true, 'whatsapp');

-- Certifications
INSERT INTO organization_certifications (org_id, facility_id, certification_type, 
  certification_number, issued_by, issue_date, expiry_date, status)
VALUES
('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111101', 
 'ISO 13485:2016', 'ISO-MH-2023-1234', 'TUV SUD', '2023-01-15', '2026-01-15', 'active'),
('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111101', 
 'CE Mark', 'CE-EU-2023-5678', 'European Commission', '2023-03-01', '2028-03-01', 'active'),
('11111111-1111-1111-1111-111111111111', NULL, 
 'FDA 510(k)', 'FDA-510K-2023-9012', 'US FDA', '2023-06-01', '2026-06-01', 'active');

-- ============================================================================
-- 2. GE HEALTHCARE INDIA
-- ============================================================================

INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  legal_entity_name, registration_number, tax_id, year_established, 
  employee_count, website, logo_url, metadata, created_at)
VALUES (
  '22222222-2222-2222-2222-222222222222',
  'GE Healthcare India',
  'GE Healthcare',
  'manufacturer',
  'active',
  true,
  'GE Healthcare India Private Limited',
  'U33112KA2004PTC034567',
  'AABCG5678H',
  2004,
  3000,
  'https://www.gehealthcare.co.in',
  'https://www.gehealthcare.com/logo.png',
  '{"specializations": ["Medical Imaging", "Ultrasound", "Patient Monitoring", "Life Support"], "global_presence": true}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, status, operational_since)
VALUES 
-- Bangalore Manufacturing
('22222222-2222-2222-2222-222222222201',
 '22222222-2222-2222-2222-222222222222',
 'GE Bangalore Multi-Modal Manufacturing',
 'GEH-BLR-MFG',
 'manufacturing_plant',
 '{"line1": "John F Welch Technology Centre", "line2": "EPIP Zone, Whitefield", "city": "Bangalore", "state": "Karnataka", "pincode": "560066", "country": "India"}',
 '{"monday_friday": "08:00-20:00", "saturday": "08:00-14:00"}',
 ARRAY['Manufacturing', 'Assembly', 'Export Operations'],
 ARRAY['Ultrasound Systems', 'Patient Monitors', 'ECG Machines'],
 'active',
 '2004-08-01'),

-- Mumbai Service Hub
('22222222-2222-2222-2222-222222222202',
 '22222222-2222-2222-2222-222222222222',
 'GE Mumbai Service Hub',
 'GEH-MUM-SVC',
 'service_center',
 '{"line1": "Peninsula Corporate Park", "line2": "Lower Parel", "city": "Mumbai", "state": "Maharashtra", "pincode": "400013", "country": "India"}',
 '{"monday_friday": "09:00-18:00", "247_emergency": true}',
 ARRAY['24x7 Service', 'Remote Monitoring', 'Preventive Maintenance'],
 ARRAY['All GE Equipment'],
 'active',
 '2005-01-01'),

-- Pune Manufacturing
('22222222-2222-2222-2222-222222222203',
 '22222222-2222-2222-2222-222222222222',
 'GE Pune Excellence Center',
 'GEH-PUN-MFG',
 'manufacturing_plant',
 '{"line1": "Hinjewadi Phase 2", "city": "Pune", "state": "Maharashtra", "pincode": "411057", "country": "India"}',
 '{"monday_friday": "08:00-18:00"}',
 ARRAY['Manufacturing', 'Quality Assurance', 'Training'],
 ARRAY['Anesthesia Systems', 'Ventilators'],
 'active',
 '2015-04-01');

-- Contact Persons
INSERT INTO contact_persons (org_id, name, designation, department, email, primary_phone, 
  whatsapp_number, is_primary, can_approve_orders, can_raise_tickets)
VALUES
('22222222-2222-2222-2222-222222222222', 'Vineeta Sharma', 'President & CEO', 'Executive', 
 'vineeta.sharma@ge.com', '+91-80-6721-5000', '+91-98450-23456', true, true, true),
('22222222-2222-2222-2222-222222222222', 'Karthik Reddy', 'Director - Service Operations', 'Service', 
 'karthik.reddy@ge.com', '+91-80-6721-5100', '+91-98450-23457', false, true, true);

-- Certifications
INSERT INTO organization_certifications (org_id, facility_id, certification_type, 
  certification_number, issued_by, issue_date, expiry_date, status)
VALUES
('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222201', 
 'ISO 13485:2016', 'ISO-KA-2023-2345', 'BSI Group', '2023-02-01', '2026-02-01', 'active'),
('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222201', 
 'CE Mark', 'CE-EU-2023-6789', 'Notified Body 0123', '2023-04-01', '2028-04-01', 'active');

-- ============================================================================
-- 3. PHILIPS HEALTHCARE INDIA
-- ============================================================================

INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  legal_entity_name, registration_number, tax_id, year_established, 
  employee_count, website, logo_url, metadata, created_at)
VALUES (
  '33333333-3333-3333-3333-333333333333',
  'Philips Healthcare India',
  'Philips Healthcare',
  'manufacturer',
  'active',
  true,
  'Philips India Limited',
  'L32109MH1930PLC001153',
  'AABCP1234M',
  1930,
  4500,
  'https://www.philips.co.in/healthcare',
  'https://www.philips.com/logo.png',
  '{"specializations": ["Diagnostic Imaging", "Image-Guided Therapy", "Patient Monitoring", "Sleep & Respiratory Care"], "heritage": "90+ years in India"}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, status, operational_since)
VALUES 
-- Pune Manufacturing
('33333333-3333-3333-3333-333333333301',
 '33333333-3333-3333-3333-333333333333',
 'Philips Innovation Campus Pune',
 'PHI-PUN-MFG',
 'manufacturing_plant',
 '{"line1": "Manyata Tech Park", "city": "Pune", "state": "Maharashtra", "pincode": "411045", "country": "India"}',
 '{"monday_friday": "08:00-18:00"}',
 ARRAY['Manufacturing', 'R&D', 'Innovation Hub'],
 ARRAY['MRI', 'CT', 'X-Ray', 'Ultrasound'],
 'active',
 '1983-01-01'),

-- Gurgaon Office & Service
('33333333-3333-3333-3333-333333333302',
 '33333333-3333-3333-3333-333333333333',
 'Philips Gurgaon Regional Office',
 'PHI-GGN-OFF',
 'service_center',
 '{"line1": "Building 5, Tower A, DLF Cyber City", "city": "Gurgaon", "state": "Haryana", "pincode": "122002", "country": "India"}',
 '{"monday_friday": "09:00-18:00", "customer_support_247": true}',
 ARRAY['Regional Management', 'Service Coordination', 'Customer Support'],
 ARRAY['All Philips Equipment'],
 'active',
 '2000-06-01');

-- Contact Persons
INSERT INTO contact_persons (org_id, name, designation, department, email, primary_phone, 
  whatsapp_number, is_primary, can_approve_orders, can_raise_tickets)
VALUES
('33333333-3333-3333-3333-333333333333', 'Daniel Mazon', 'Managing Director', 'Executive', 
 'daniel.mazon@philips.com', '+91-20-6725-9000', '+91-98230-34567', true, true, true),
('33333333-3333-3333-3333-333333333333', 'Sunita Iyer', 'Head - Customer Services', 'Service', 
 'sunita.iyer@philips.com', '+91-20-6725-9100', '+91-98230-34568', false, true, true);

-- Certifications
INSERT INTO organization_certifications (org_id, facility_id, certification_type, 
  certification_number, issued_by, issue_date, expiry_date, status)
VALUES
('33333333-3333-3333-3333-333333333333', '33333333-3333-3333-3333-333333333301', 
 'ISO 13485:2016', 'ISO-MH-2023-3456', 'DNV GL', '2023-01-20', '2026-01-20', 'active');

-- ============================================================================
-- 4. MEDTRONIC INDIA
-- ============================================================================

INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  legal_entity_name, registration_number, tax_id, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  '44444444-4444-4444-4444-444444444444',
  'Medtronic India',
  'Medtronic',
  'manufacturer',
  'active',
  true,
  'Medtronic India Private Limited',
  'U33112MH1998PTC115678',
  'AABCM9012P',
  1998,
  1200,
  'https://www.medtronic.com/in-en',
  '{"specializations": ["Cardiac Devices", "Spine & Neurosurgery", "Diabetes Management", "Surgical Innovations"]}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, equipment_types, status, operational_since)
VALUES 
('44444444-4444-4444-4444-444444444401',
 '44444444-4444-4444-4444-444444444444',
 'Medtronic Hyderabad Manufacturing',
 'MED-HYD-MFG',
 'manufacturing_plant',
 '{"line1": "Plot No. 18, IDA Bollaram", "city": "Hyderabad", "state": "Telangana", "pincode": "502325", "country": "India"}',
 ARRAY['Manufacturing', 'Quality Control', 'Regulatory Affairs'],
 ARRAY['Cardiac Rhythm Devices', 'Surgical Equipment'],
 'active',
 '2008-09-01'),
('44444444-4444-4444-4444-444444444402',
 '44444444-4444-4444-4444-444444444444',
 'Medtronic Mumbai Regional Office',
 'MED-MUM-OFF',
 'sales_office',
 '{"line1": "601-604, Hallmark Business Plaza", "line2": "Sant Dyaneshwar Marg, Bandra East", "city": "Mumbai", "state": "Maharashtra", "pincode": "400051", "country": "India"}',
 ARRAY['Sales', 'Clinical Support', 'Training'],
 ARRAY['All Medtronic Products'],
 'active',
 '1998-01-01');

-- Contact Persons
INSERT INTO contact_persons (org_id, name, designation, department, email, primary_phone, is_primary, can_approve_orders)
VALUES
('44444444-4444-4444-4444-444444444444', 'Madan Krishnan', 'Vice President & Managing Director', 'Executive', 
 'madan.krishnan@medtronic.com', '+91-22-6112-7000', true, true);

-- Certifications
INSERT INTO organization_certifications (org_id, facility_id, certification_type, 
  certification_number, issued_by, issue_date, expiry_date, status)
VALUES
('44444444-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444401', 
 'ISO 13485:2016', 'ISO-TG-2023-4567', 'TUV Rheinland', '2023-03-15', '2026-03-15', 'active');

-- ============================================================================
-- 5. ABBOTT LABORATORIES INDIA
-- ============================================================================

INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  legal_entity_name, year_established, employee_count, website, metadata, created_at)
VALUES (
  '55555555-5555-5555-5555-555555555555',
  'Abbott Laboratories India',
  'Abbott India',
  'manufacturer',
  'active',
  true,
  'Abbott India Limited',
  1944,
  3500,
  'https://www.abbott.co.in',
  '{"specializations": ["Diagnostics", "Medical Devices", "Nutrition", "Established Pharmaceuticals"]}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, equipment_types, status, operational_since)
VALUES 
('55555555-5555-5555-5555-555555555501',
 '55555555-5555-5555-5555-555555555555',
 'Abbott Goa Manufacturing Plant',
 'ABT-GOA-MFG',
 'manufacturing_plant',
 '{"line1": "Verna Industrial Estate", "city": "Verna", "state": "Goa", "pincode": "403722", "country": "India"}',
 ARRAY['Diagnostics Manufacturing', 'Point-of-Care Testing'],
 ARRAY['Blood Glucose Monitors', 'Cardiac Markers', 'Rapid Diagnostic Tests'],
 'active',
 '2000-01-01');

-- Contact Persons
INSERT INTO contact_persons (org_id, name, designation, department, email, primary_phone, is_primary)
VALUES
('55555555-5555-5555-5555-555555555555', 'Anil Bhalla', 'Managing Director', 'Executive', 
 'anil.bhalla@abbott.com', '+91-22-5046-1000', true);

-- ============================================================================
-- 6-10. ADDITIONAL MANUFACTURERS (Simplified for seed data)
-- ============================================================================

-- 6. B. Braun Medical India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  '66666666-6666-6666-6666-666666666666',
  'B. Braun Medical India',
  'B. Braun',
  'manufacturer',
  'active',
  true,
  2009,
  800,
  'https://www.bbraun.co.in',
  '{"specializations": ["Infusion Therapy", "Clinical Nutrition", "Surgical Instruments"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, status, operational_since)
VALUES 
('66666666-6666-6666-6666-666666666601',
 '66666666-6666-6666-6666-666666666666',
 'B. Braun Bangalore Manufacturing',
 'BBR-BLR-MFG',
 'manufacturing_plant',
 '{"line1": "Plot No. 1A, Bommasandra Industrial Area", "city": "Bangalore", "state": "Karnataka", "pincode": "560099", "country": "India"}',
 ARRAY['Manufacturing', 'Quality Assurance'],
 'active',
 '2010-06-01');

-- 7. Baxter India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  '77777777-7777-7777-7777-777777777777',
  'Baxter India',
  'Baxter',
  'manufacturer',
  'active',
  true,
  1997,
  1500,
  'https://www.baxter.co.in',
  '{"specializations": ["Renal Care", "Medication Delivery", "Pharmaceuticals", "Clinical Nutrition"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, status, operational_since)
VALUES 
('77777777-7777-7777-7777-777777777701',
 '77777777-7777-7777-7777-777777777777',
 'Baxter Ahmedabad Manufacturing',
 'BAX-AMD-MFG',
 'manufacturing_plant',
 '{"line1": "Village Moraiya, Sanand", "city": "Ahmedabad", "state": "Gujarat", "pincode": "382110", "country": "India"}',
 ARRAY['Renal Products Manufacturing', 'Quality Control'],
 'active',
 '1997-01-01');

-- 8. Becton Dickinson (BD) India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  '88888888-8888-8888-8888-888888888888',
  'Becton Dickinson India',
  'BD India',
  'manufacturer',
  'active',
  true,
  1996,
  2200,
  'https://www.bd.com/en-in',
  '{"specializations": ["Medical Devices", "Laboratory Equipment", "Diagnostic Systems"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, status, operational_since)
VALUES 
('88888888-8888-8888-8888-888888888801',
 '88888888-8888-8888-8888-888888888888',
 'BD Bawal Manufacturing',
 'BDI-BAW-MFG',
 'manufacturing_plant',
 '{"line1": "Plot No. SP-383-386, RIICO Industrial Area", "city": "Bawal", "state": "Haryana", "pincode": "123501", "country": "India"}',
 ARRAY['Syringes & Needles Manufacturing', 'Laboratory Equipment'],
 'active',
 '2009-01-01');

-- 9. Stryker India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  '99999999-9999-9999-9999-999999999999',
  'Stryker India',
  'Stryker',
  'manufacturer',
  'active',
  true,
  2006,
  900,
  'https://www.stryker.com/in/en',
  '{"specializations": ["Orthopedic Implants", "Surgical Equipment", "Neurotechnology", "Spine"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, status, operational_since)
VALUES 
('99999999-9999-9999-9999-999999999901',
 '99999999-9999-9999-9999-999999999999',
 'Stryker Gurgaon Office',
 'STR-GGN-OFF',
 'sales_office',
 '{"line1": "Level 3, Tower B, Global Business Park", "city": "Gurgaon", "state": "Haryana", "pincode": "122002", "country": "India"}',
 ARRAY['Sales', 'Clinical Training', 'Technical Support'],
 'active',
 '2006-03-01');

-- 10. Nihon Kohden India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, 
  employee_count, website, metadata, created_at)
VALUES (
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'Nihon Kohden India',
  'Nihon Kohden',
  'manufacturer',
  'active',
  true,
  2012,
  450,
  'https://www.nihonkohden.com/in',
  '{"specializations": ["Patient Monitors", "Defibrillators", "ECG Systems", "Neurology Equipment"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, status, operational_since)
VALUES 
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
 'Nihon Kohden Delhi Regional Office',
 'NKI-DEL-OFF',
 'service_center',
 '{"line1": "F-26, Okhla Industrial Area Phase 1", "city": "New Delhi", "state": "Delhi", "pincode": "110020", "country": "India"}',
 ARRAY['Sales', 'Service', 'Training'],
 'active',
 '2012-06-01');

-- ============================================================================
-- SEED DATA COMPLETE: 10 Manufacturers with 18 Facilities
-- ============================================================================
