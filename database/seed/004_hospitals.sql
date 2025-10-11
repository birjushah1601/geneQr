-- ============================================================================
-- SEED DATA: HOSPITALS (10 Organizations)
-- ============================================================================
-- Major hospitals across India with in-house Biomedical Engineering (BME) teams

-- ============================================================================
-- MULTI-SPECIALTY HOSPITALS
-- ============================================================================

-- 1. Apollo Hospitals Delhi
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0100-000000000001',
  'Apollo Hospitals Delhi',
  'Apollo Delhi',
  'hospital',
  'active',
  true,
  1996,
  1200,
  '{"beds": 710, "specialties": ["Cardiology", "Oncology", "Neurology", "Orthopedics"], "bme_team_size": 12, "equipment_count": 450}',
  NOW()
);

-- 2. Fortis Hospital Bangalore
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0101-000000000002',
  'Fortis Hospital Bangalore',
  'Fortis Bangalore',
  'hospital',
  'active',
  true,
  2006,
  800,
  '{"beds": 400, "specialties": ["Cardiology", "Oncology", "Gastroenterology"], "bme_team_size": 8, "equipment_count": 320}',
  NOW()
);

-- 3. Manipal Hospitals Mumbai  
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0102-000000000003',
  'Manipal Hospitals Mumbai',
  'Manipal Mumbai',
  'hospital',
  'active',
  true,
  2012,
  650,
  '{"beds": 350, "specialties": ["Oncology", "Nephrology", "Neurology"], "bme_team_size": 7, "equipment_count": 280}',
  NOW()
);

-- 4. Max Super Speciality Hospital Delhi
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0103-000000000004',
  'Max Super Speciality Hospital Delhi',
  'Max Delhi',
  'hospital',
  'active',
  true,
  2001,
  950,
  '{"beds": 550, "specialties": ["Cardiology", "Oncology", "Neurosciences", "Orthopedics"], "bme_team_size": 10, "equipment_count": 400}',
  NOW()
);

-- 5. Narayana Health Bangalore
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0104-000000000005',
  'Narayana Health Bangalore',
  'Narayana Bangalore',
  'hospital',
  'active',
  true,
  2001,
  1100,
  '{"beds": 650, "specialties": ["Cardiology", "Cardiac Surgery", "Oncology"], "bme_team_size": 11, "equipment_count": 420}',
  NOW()
);

-- 6. KIMS Hospital Hyderabad
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0105-000000000006',
  'KIMS Hospital Hyderabad',
  'KIMS Hyderabad',
  'hospital',
  'active',
  true,
  2009,
  750,
  '{"beds": 450, "specialties": ["Oncology", "Nephrology", "Gastroenterology"], "bme_team_size": 9, "equipment_count": 350}',
  NOW()
);

-- 7. Medanta The Medicity Gurgaon
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0106-000000000007',
  'Medanta The Medicity Gurgaon',
  'Medanta Gurgaon',
  'hospital',
  'active',
  true,
  2009,
  1500,
  '{"beds": 1250, "specialties": ["Cardiology", "Oncology", "Neurosciences", "Transplants"], "bme_team_size": 15, "equipment_count": 650}',
  NOW()
);

-- 8. MGM Hospital Chennai
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0107-000000000008',
  'MGM Hospital Chennai',
  'MGM Chennai',
  'hospital',
  'active',
  true,
  2008,
  600,
  '{"beds": 400, "specialties": ["Cardiology", "Orthopedics", "Neurology"], "bme_team_size": 8, "equipment_count": 310}',
  NOW()
);

-- 9. Ruby Hall Clinic Pune
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0108-000000000009',
  'Ruby Hall Clinic Pune',
  'Ruby Hall Pune',
  'hospital',
  'active',
  true,
  1973,
  550,
  '{"beds": 350, "specialties": ["Cardiology", "Oncology", "Orthopedics"], "bme_team_size": 7, "equipment_count": 290}',
  NOW()
);

-- 10. AMRI Hospitals Kolkata
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  '00000000-0000-0000-0109-000000000010',
  'AMRI Hospitals Kolkata',
  'AMRI Kolkata',
  'hospital',
  'active',
  true,
  1996,
  700,
  '{"beds": 450, "specialties": ["Cardiology", "Oncology", "Nephrology"], "bme_team_size": 9, "equipment_count": 360}',
  NOW()
);

-- ============================================================================
-- HOSPITAL FACILITIES
-- ============================================================================

-- Apollo Delhi - Main Campus + BME Department
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0100-000000000101',
 '00000000-0000-0000-0100-000000000001',
 'Apollo Hospital Delhi Main Campus',
 'APO-DEL-MC',
 'hospital',
 '{"line1": "Sarita Vihar", "city": "New Delhi", "state": "Delhi", "pincode": "110076", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Diagnostics', 'Cardiology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Ventilators', 'Dialysis', 'Lab Equipment'],
 ARRAY['Delhi', 'NCR'],
 'active',
 '1996-01-15'),
('00000000-0000-0000-0100-000000000102',
 '00000000-0000-0000-0100-000000000001',
 'BME Department',
 'APO-DEL-BME',
 'service_center',
 '{"line1": "Basement 2, Main Building", "city": "New Delhi", "state": "Delhi", "pincode": "110076", "country": "India"}',
 '{"monday_saturday": "08:00-20:00", "emergency_247": true}',
 ARRAY['Preventive Maintenance', 'Repair', 'Calibration', 'Emergency Support'],
 ARRAY['All Equipment'],
 ARRAY['Delhi', 'NCR'],
 'active',
 '1996-01-15');

-- Fortis Bangalore
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0101-000000000103',
 '00000000-0000-0000-0101-000000000002',
 'Fortis Hospital Bangalore Main',
 'FOR-BLR-MC',
 'hospital',
 '{"line1": "154/9, Bannerghatta Road", "city": "Bangalore", "state": "Karnataka", "pincode": "560076", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Diagnostics'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Ventilators', 'Lab Equipment'],
 ARRAY['Karnataka'],
 'active',
 '2006-03-10');

-- Manipal Mumbai
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0102-000000000104',
 '00000000-0000-0000-0102-000000000003',
 'Manipal Hospital Mumbai',
 'MAN-MUM-MC',
 'hospital',
 '{"line1": "Mulund Goregaon Link Road", "city": "Mumbai", "state": "Maharashtra", "pincode": "400078", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Oncology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Lab Equipment'],
 ARRAY['Maharashtra'],
 'active',
 '2012-06-01');

-- Max Delhi
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0103-000000000105',
 '00000000-0000-0000-0103-000000000004',
 'Max Super Speciality Hospital Saket',
 'MAX-DEL-SKT',
 'hospital',
 '{"line1": "1-2, Press Enclave Road, Saket", "city": "New Delhi", "state": "Delhi", "pincode": "110017", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Cardiology', 'Oncology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Ventilators', 'Lab Equipment'],
 ARRAY['Delhi', 'NCR'],
 'active',
 '2001-09-15');

-- Narayana Bangalore
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0104-000000000106',
 '00000000-0000-0000-0104-000000000005',
 'Narayana Health City',
 'NAR-BLR-HC',
 'hospital',
 '{"line1": "258/A, Bommasandra Industrial Area", "city": "Bangalore", "state": "Karnataka", "pincode": "560099", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Cardiac Surgery', 'Cardiology', 'Oncology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Ventilators', 'Cardiac Cath Lab'],
 ARRAY['Karnataka'],
 'active',
 '2001-01-25');

-- KIMS Hyderabad
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0105-000000000107',
 '00000000-0000-0000-0105-000000000006',
 'KIMS Hospital Secunderabad',
 'KIMS-HYD-SEC',
 'hospital',
 '{"line1": "1-8-31/1, Minister Road", "city": "Hyderabad", "state": "Telangana", "pincode": "500003", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Oncology', 'Nephrology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Dialysis', 'Lab Equipment'],
 ARRAY['Telangana', 'Andhra Pradesh'],
 'active',
 '2009-05-20');

-- Medanta Gurgaon
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0106-000000000108',
 '00000000-0000-0000-0106-000000000007',
 'Medanta The Medicity',
 'MED-GRG-MC',
 'hospital',
 '{"line1": "CH Baktawar Singh Road, Sector 38", "city": "Gurgaon", "state": "Haryana", "pincode": "122001", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Cardiology', 'Transplants', 'Oncology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Ventilators', 'Cardiac Cath Lab', 'Lab Equipment'],
 ARRAY['Haryana', 'Delhi', 'NCR'],
 'active',
 '2009-11-16');

-- MGM Chennai
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0107-000000000109',
 '00000000-0000-0000-0107-000000000008',
 'MGM Healthcare',
 'MGM-CHE-HC',
 'hospital',
 '{"line1": "New No 72, Old No 54 Nelson Manickam Road", "city": "Chennai", "state": "Tamil Nadu", "pincode": "600029", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Cardiology', 'Orthopedics'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Lab Equipment'],
 ARRAY['Tamil Nadu'],
 'active',
 '2008-02-14');

-- Ruby Hall Pune
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0108-000000000110',
 '00000000-0000-0000-0108-000000000009',
 'Ruby Hall Clinic',
 'RUB-PUN-CLI',
 'hospital',
 '{"line1": "40, Sassoon Road", "city": "Pune", "state": "Maharashtra", "pincode": "411001", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Cardiology', 'Orthopedics'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Lab Equipment'],
 ARRAY['Maharashtra'],
 'active',
 '1973-01-01');

-- AMRI Kolkata
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('00000000-0000-0000-0109-000000000111',
 '00000000-0000-0000-0109-000000000010',
 'AMRI Hospital Salt Lake',
 'AMR-KOL-SLK',
 'hospital',
 '{"line1": "JC-16 & 17, Sector III, Salt Lake City", "city": "Kolkata", "state": "West Bengal", "pincode": "700098", "country": "India"}',
 '{"emergency_247": true}',
 ARRAY['Emergency', 'ICU', 'Surgery', 'Cardiology', 'Nephrology'],
 ARRAY['Diagnostic Imaging', 'Patient Monitors', 'Dialysis', 'Lab Equipment'],
 ARRAY['West Bengal'],
 'active',
 '1996-08-01');

-- ============================================================================
-- SUMMARY
-- ============================================================================
-- 10 hospitals added with:
--   - Major multi-specialty hospitals across India
--   - In-house BME teams (7-15 engineers per hospital)
--   - 290-650 beds per hospital
--   - 280-650 medical equipment items per hospital
--   - Total: ~86 BME engineers across all hospitals
--   - These act as Tier-5 fallback engineers in service routing
