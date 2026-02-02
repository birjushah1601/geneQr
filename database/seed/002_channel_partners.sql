-- ============================================================================
-- SEED DATA: Channel Partners (20 Organizations)
-- ============================================================================
-- Medical equipment Channel Partners across India with manufacturer relationships

-- ============================================================================
-- NORTH INDIA Channel Partners (Delhi, UP, Punjab, Haryana, Rajasthan)
-- ============================================================================

-- 1. MedEquip Channel Partners (Multi-Brand, North India)
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, website, metadata, created_at)
VALUES (
  'dddd0001-0001-0001-0001-000000000001',
  'MedEquip Channel Partners Private Limited',
  'MedEquip Channel Partners',
  'Channel Partner',
  'active',
  true,
  2005,
  120,
  'https://www.medequip-india.com',
  '{"regions": ["Delhi", "UP", "Punjab", "Haryana"], "brands": ["Siemens", "GE", "Philips"]}',
  NOW()
);

-- Facilities
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, operational_hours, services_offered, equipment_types, coverage_states, status, operational_since)
VALUES 
('dddd0001-0001-0001-0001-000000000101',
 'dddd0001-0001-0001-0001-000000000001',
 'MedEquip Delhi Warehouse',
 'MED-DEL-WH',
 'warehouse',
 '{"line1": "Plot 45, Sector 8, IMT Manesar", "city": "Gurgaon", "state": "Haryana", "pincode": "122050", "country": "India"}',
 '{"monday_saturday": "09:00-18:00"}',
 ARRAY['Warehousing', 'Distribution', 'Basic Service'],
 ARRAY['Diagnostic Imaging', 'Patient Monitoring'],
 ARRAY['Delhi', 'Haryana', 'Punjab', 'Chandigarh'],
 'active',
 '2005-06-01'),
('dddd0001-0001-0001-0001-000000000102',
 'dddd0001-0001-0001-0001-000000000001',
 'MedEquip Delhi Service Center',
 'MED-DEL-SVC',
 'service_center',
 '{"line1": "B-14, Okhla Industrial Area Phase 2", "city": "New Delhi", "state": "Delhi", "pincode": "110020", "country": "India"}',
 '{"monday_saturday": "09:00-19:00"}',
 ARRAY['Installation', 'Maintenance', 'AMC Support'],
 ARRAY['All Equipment'],
 ARRAY['Delhi', 'NCR'],
 'active',
 '2006-01-01');

-- Contacts
INSERT INTO contact_persons (org_id, name, designation, email, primary_phone, is_primary, can_approve_orders)
VALUES
('dddd0001-0001-0001-0001-000000000001', 'Arun Mehta', 'Managing Director', 
 'arun.mehta@medequip-india.com', '+91-11-4567-9000', true, true),
('dddd0001-0001-0001-0001-000000000001', 'Ritu Sharma', 'Service Head', 
 'ritu.sharma@medequip-india.com', '+91-11-4567-9100', false, false);

-- Relationships with Manufacturers
INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, relationship_status,
  exclusive, commission_percentage, credit_limit, annual_target, priority_level, metadata)
VALUES
-- Siemens â†’ MedEquip
('11111111-1111-1111-1111-111111111111', 'dddd0001-0001-0001-0001-000000000001',
 'authorized_channel_partner', 'active', false, 12.5, 50000000, 200000000, 1,
 '{"product_lines": ["CT Scanners", "X-Ray Systems", "Lab Diagnostics"]}'),
-- GE â†’ MedEquip
('22222222-2222-2222-2222-222222222222', 'dddd0001-0001-0001-0001-000000000001',
 'authorized_channel_partner', 'active', false, 10.0, 40000000, 150000000, 2,
 '{"product_lines": ["Ultrasound", "Patient Monitors"]}'),
-- Philips â†’ MedEquip
('33333333-3333-3333-3333-333333333333', 'dddd0001-0001-0001-0001-000000000001',
 'regional_channel_partner', 'active', false, 11.0, 30000000, 100000000, 2,
 '{"product_lines": ["Patient Monitoring", "Respiratory Care"]}');

-- 2. HealthTech Solutions (Delhi/NCR specialist)
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, metadata, created_at)
VALUES (
  'dddd0002-0002-0002-0002-000000000002',
  'HealthTech Solutions India',
  'HealthTech Solutions',
  'Channel Partner',
  'active',
  true,
  2010,
  75,
  '{"regions": ["Delhi-NCR"], "brands": ["Medtronic", "Abbott", "BD"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, coverage_states, status, operational_since)
VALUES 
('dddd0002-0002-0002-0002-000000000201',
 'dddd0002-0002-0002-0002-000000000002',
 'HealthTech Noida Distribution Center',
 'HTH-NOI-DC',
 'distribution_center',
 '{"line1": "Plot 23, Sector 63", "city": "Noida", "state": "Uttar Pradesh", "pincode": "201301", "country": "India"}',
 ARRAY['Distribution', 'Logistics', 'Technical Support'],
 ARRAY['Delhi', 'Uttar Pradesh'],
 'active',
 '2010-03-01');

INSERT INTO contact_persons (org_id, name, designation, email, primary_phone, is_primary)
VALUES
('dddd0002-0002-0002-0002-000000000002', 'Vikram Singh', 'CEO', 
 'vikram@healthtech-solutions.in', '+91-120-456-7890', true);

-- Relationships
INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, relationship_status,
  commission_percentage, credit_limit, annual_target)
VALUES
('44444444-4444-4444-4444-444444444444', 'dddd0002-0002-0002-0002-000000000002',
 'authorized_channel_partner', 'active', 15.0, 25000000, 80000000),
('55555555-5555-5555-5555-555555555555', 'dddd0002-0002-0002-0002-000000000002',
 'authorized_channel_partner', 'active', 13.0, 20000000, 60000000);

-- ============================================================================
-- SOUTH INDIA Channel Partners (Bangalore, Chennai, Hyderabad, Kerala)
-- ============================================================================

-- 3. Southern Medical Supplies (Bangalore-based, Multi-state)
INSERT INTO organizations (id, name, display_name, org_type, status, verified, 
  year_established, employee_count, website, metadata, created_at)
VALUES (
  'dddd0003-0003-0003-0003-000000000003',
  'Southern Medical Supplies',
  'SMS',
  'Channel Partner',
  'active',
  true,
  2002,
  150,
  'https://www.southernmed.co.in',
  '{"regions": ["Karnataka", "Tamil Nadu", "Kerala", "Andhra Pradesh"], "brands": ["Siemens", "Philips", "Stryker"]}',
  NOW()
);

INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, 
  address, services_offered, coverage_states, status, operational_since)
VALUES 
('dddd0003-0003-0003-0003-000000000301',
 'dddd0003-0003-0003-0003-000000000003',
 'SMS Bangalore Hub',
 'SMS-BLR-HUB',
 'distribution_center',
 '{"line1": "Survey No. 45, Jigani Industrial Area", "city": "Bangalore", "state": "Karnataka", "pincode": "562106", "country": "India"}',
 ARRAY['Distribution', 'Service', 'Training'],
 ARRAY['Karnataka', 'Tamil Nadu', 'Kerala'],
 'active',
 '2002-08-01'),
('dddd0003-0003-0003-0003-000000000302',
 'dddd0003-0003-0003-0003-000000000003',
 'SMS Chennai Branch',
 'SMS-CHE-BR',
 'service_center',
 '{"line1": "No. 45, Ambattur Industrial Estate", "city": "Chennai", "state": "Tamil Nadu", "pincode": "600058", "country": "India"}',
 ARRAY['Service', 'Parts Supply'],
 ARRAY['Tamil Nadu'],
 'active',
 '2005-04-01');

INSERT INTO contact_persons (org_id, name, designation, email, primary_phone, is_primary)
VALUES
('dddd0003-0003-0003-0003-000000000003', 'Ramesh Kumar', 'Director', 
 'ramesh@southernmed.co.in', '+91-80-2734-5600', true);

INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, relationship_status,
  exclusive, commission_percentage, credit_limit, annual_target)
VALUES
('11111111-1111-1111-1111-111111111111', 'dddd0003-0003-0003-0003-000000000003',
 'regional_channel_partner', 'active', false, 13.0, 60000000, 250000000),
('33333333-3333-3333-3333-333333333333', 'dddd0003-0003-0003-0003-000000000003',
 'authorized_channel_partner', 'active', false, 12.0, 45000000, 180000000);

-- 4-20: Additional Channel Partners (Simplified)

-- West India
INSERT INTO organizations (id, name, display_name, org_type, status, verified, year_established, employee_count, metadata, created_at)
VALUES 
('dddd0004-0004-0004-0004-000000000004', 'Western Medical Supplies', 'WMS', 'Channel Partner', 'active', true, 2008, 90, '{"regions": ["Maharashtra", "Gujarat"]}', NOW()),
('dddd0005-0005-0005-0005-000000000005', 'Mumbai Healthcare Channel Partners', 'MHD', 'Channel Partner', 'active', true, 2006, 110, '{"regions": ["Mumbai", "Pune"]}', NOW()),

-- East India
('dddd0006-0006-0006-0006-000000000006', 'Eastern Medical Equipment Co.', 'EMEC', 'Channel Partner', 'active', true, 2007, 65, '{"regions": ["West Bengal", "Odisha", "Bihar"]}', NOW()),
('dddd0007-0007-0007-0007-000000000007', 'Kolkata MedTech Channel Partners', 'KMD', 'Channel Partner', 'active', true, 2011, 55, '{"regions": ["Kolkata"]}', NOW()),

-- Central India
('dddd0008-0008-0008-0008-000000000008', 'Central India MedEquip', 'CIME', 'Channel Partner', 'active', true, 2009, 70, '{"regions": ["Madhya Pradesh", "Chhattisgarh"]}', NOW()),

-- Regional specialists
('dddd0009-0009-0009-0009-000000000009', 'Hyderabad Medical Systems', 'HMS', 'Channel Partner', 'active', true, 2004, 85, '{"regions": ["Telangana", "Andhra Pradesh"]}', NOW()),
('dddd0010-0010-0010-0010-000000000010', 'Pune Diagnostics Distribution', 'PDD', 'Channel Partner', 'active', true, 2012, 60, '{"regions": ["Pune"]}', NOW()),
('dddd0011-0011-0011-0011-000000000011', 'Ahmedabad Healthcare Solutions', 'AHS', 'Channel Partner', 'active', true, 2010, 75, '{"regions": ["Gujarat"]}', NOW()),
('dddd0012-0012-0012-0012-000000000012', 'Kerala Medical Channel Partners', 'KMD', 'Channel Partner', 'active', true, 2008, 50, '{"regions": ["Kerala"]}', NOW()),
('dddd0013-0013-0013-0013-000000000013', 'Jaipur Med Solutions', 'JMS', 'Channel Partner', 'active', true, 2013, 45, '{"regions": ["Rajasthan"]}', NOW()),
('dddd0014-0014-0014-0014-000000000014', 'Chandigarh Medical Equipment', 'CME', 'Channel Partner', 'active', true, 2007, 55, '{"regions": ["Punjab", "Chandigarh"]}', NOW()),
('dddd0015-0015-0015-0015-000000000015', 'Lucknow Health Systems', 'LHS', 'Channel Partner', 'active', true, 2009, 48, '{"regions": ["Uttar Pradesh"]}', NOW()),
('dddd0016-0016-0016-0016-000000000016', 'Coimbatore Medical Supplies', 'CMS', 'Channel Partner', 'active', true, 2011, 40, '{"regions": ["Tamil Nadu"]}', NOW()),
('dddd0017-0017-0017-0017-000000000017', 'Visakhapatnam MedTech', 'VMT', 'Channel Partner', 'active', true, 2010, 42, '{"regions": ["Andhra Pradesh"]}', NOW()),
('dddd0018-0018-0018-0018-000000000018', 'Indore Healthcare Channel Partners', 'IHD', 'Channel Partner', 'active', true, 2012, 38, '{"regions": ["Madhya Pradesh"]}', NOW()),
('dddd0019-0019-0019-0019-000000000019', 'Surat Medical Equipment', 'SME', 'Channel Partner', 'active', true, 2014, 35, '{"regions": ["Gujarat"]}', NOW()),
('dddd0020-0020-0020-0020-000000000020', 'Nagpur Diagnostics Distribution', 'NDD', 'Channel Partner', 'active', true, 2013, 32, '{"regions": ["Maharashtra"]}', NOW());

-- Facilities for Channel Partners 4-20 (1 facility each)
INSERT INTO organization_facilities (id, org_id, facility_name, facility_code, facility_type, address, services_offered, status, operational_since)
VALUES 
('dddd0004-0004-0004-0004-000000000401', 'dddd0004-0004-0004-0004-000000000004', 'WMS Mumbai Warehouse', 'WMS-MUM-WH', 'warehouse', '{"city": "Mumbai", "state": "Maharashtra"}', ARRAY['Distribution'], 'active', '2008-01-01'),
('dddd0005-0005-0005-0005-000000000501', 'dddd0005-0005-0005-0005-000000000005', 'MHD Pune Center', 'MHD-PUN-CTR', 'distribution_center', '{"city": "Pune", "state": "Maharashtra"}', ARRAY['Distribution'], 'active', '2006-01-01'),
('dddd0006-0006-0006-0006-000000000601', 'dddd0006-0006-0006-0006-000000000006', 'EMEC Kolkata Hub', 'EMEC-KOL-HUB', 'distribution_center', '{"city": "Kolkata", "state": "West Bengal"}', ARRAY['Distribution'], 'active', '2007-01-01'),
('dddd0007-0007-0007-0007-000000000701', 'dddd0007-0007-0007-0007-000000000007', 'KMD Kolkata Branch', 'KMD-KOL-BR', 'service_center', '{"city": "Kolkata", "state": "West Bengal"}', ARRAY['Service'], 'active', '2011-01-01'),
('dddd0008-0008-0008-0008-000000000801', 'dddd0008-0008-0008-0008-000000000008', 'CIME Bhopal Center', 'CIME-BHO-CTR', 'distribution_center', '{"city": "Bhopal", "state": "Madhya Pradesh"}', ARRAY['Distribution'], 'active', '2009-01-01'),
('dddd0009-0009-0009-0009-000000000901', 'dddd0009-0009-0009-0009-000000000009', 'HMS Hyderabad Hub', 'HMS-HYD-HUB', 'distribution_center', '{"city": "Hyderabad", "state": "Telangana"}', ARRAY['Distribution', 'Service'], 'active', '2004-01-01'),
('dddd0010-0010-0010-0010-000000001001', 'dddd0010-0010-0010-0010-000000000010', 'PDD Pune Center', 'PDD-PUN-CTR', 'warehouse', '{"city": "Pune", "state": "Maharashtra"}', ARRAY['Distribution'], 'active', '2012-01-01'),
('dddd0011-0011-0011-0011-000000001101', 'dddd0011-0011-0011-0011-000000000011', 'AHS Ahmedabad Hub', 'AHS-AMD-HUB', 'distribution_center', '{"city": "Ahmedabad", "state": "Gujarat"}', ARRAY['Distribution'], 'active', '2010-01-01'),
('dddd0012-0012-0012-0012-000000001201', 'dddd0012-0012-0012-0012-000000000012', 'KMD Kochi Center', 'KMD-KOC-CTR', 'distribution_center', '{"city": "Kochi", "state": "Kerala"}', ARRAY['Distribution'], 'active', '2008-01-01'),
('dddd0013-0013-0013-0013-000000001301', 'dddd0013-0013-0013-0013-000000000013', 'JMS Jaipur Hub', 'JMS-JAI-HUB', 'warehouse', '{"city": "Jaipur", "state": "Rajasthan"}', ARRAY['Distribution'], 'active', '2013-01-01'),
('dddd0014-0014-0014-0014-000000001401', 'dddd0014-0014-0014-0014-000000000014', 'CME Chandigarh Center', 'CME-CHD-CTR', 'distribution_center', '{"city": "Chandigarh", "state": "Chandigarh"}', ARRAY['Distribution'], 'active', '2007-01-01'),
('dddd0015-0015-0015-0015-000000001501', 'dddd0015-0015-0015-0015-000000000015', 'LHS Lucknow Hub', 'LHS-LKO-HUB', 'warehouse', '{"city": "Lucknow", "state": "Uttar Pradesh"}', ARRAY['Distribution'], 'active', '2009-01-01'),
('dddd0016-0016-0016-0016-000000001601', 'dddd0016-0016-0016-0016-000000000016', 'CMS Coimbatore Center', 'CMS-CBE-CTR', 'distribution_center', '{"city": "Coimbatore", "state": "Tamil Nadu"}', ARRAY['Distribution'], 'active', '2011-01-01'),
('dddd0017-0017-0017-0017-000000001701', 'dddd0017-0017-0017-0017-000000000017', 'VMT Vizag Hub', 'VMT-VIZ-HUB', 'warehouse', '{"city": "Visakhapatnam", "state": "Andhra Pradesh"}', ARRAY['Distribution'], 'active', '2010-01-01'),
('dddd0018-0018-0018-0018-000000001801', 'dddd0018-0018-0018-0018-000000000018', 'IHD Indore Center', 'IHD-IND-CTR', 'distribution_center', '{"city": "Indore", "state": "Madhya Pradesh"}', ARRAY['Distribution'], 'active', '2012-01-01'),
('dddd0019-0019-0019-0019-000000001901', 'dddd0019-0019-0019-0019-000000000019', 'SME Surat Hub', 'SME-SUR-HUB', 'warehouse', '{"city": "Surat", "state": "Gujarat"}', ARRAY['Distribution'], 'active', '2014-01-01'),
('dddd0020-0020-0020-0020-000000002001', 'dddd0020-0020-0020-0020-000000000020', 'NDD Nagpur Center', 'NDD-NAG-CTR', 'distribution_center', '{"city": "Nagpur", "state": "Maharashtra"}', ARRAY['Distribution'], 'active', '2013-01-01');

-- Additional Channel Partner-manufacturer relationships
INSERT INTO org_relationships (parent_org_id, child_org_id, rel_type, relationship_status, commission_percentage, credit_limit, annual_target)
VALUES
-- GE relationships
('22222222-2222-2222-2222-222222222222', 'dddd0004-0004-0004-0004-000000000004', 'authorized_channel_partner', 'active', 11.0, 35000000, 120000000),
('22222222-2222-2222-2222-222222222222', 'dddd0006-0006-0006-0006-000000000006', 'regional_channel_partner', 'active', 10.5, 28000000, 90000000),
('22222222-2222-2222-2222-222222222222', 'dddd0009-0009-0009-0009-000000000009', 'authorized_channel_partner', 'active', 11.5, 32000000, 110000000),

-- Philips relationships
('33333333-3333-3333-3333-333333333333', 'dddd0005-0005-0005-0005-000000000005', 'authorized_channel_partner', 'active', 12.0, 30000000, 100000000),
('33333333-3333-3333-3333-333333333333', 'dddd0008-0008-0008-0008-000000000008', 'regional_channel_partner', 'active', 11.0, 25000000, 80000000),

-- Medtronic relationships
('44444444-4444-4444-4444-444444444444', 'dddd0009-0009-0009-0009-000000000009', 'authorized_channel_partner', 'active', 14.0, 22000000, 70000000),
('44444444-4444-4444-4444-444444444444', 'dddd0011-0011-0011-0011-000000000011', 'regional_channel_partner', 'active', 13.5, 18000000, 55000000),

-- Smaller manufacturers with regional Channel Partners
('66666666-6666-6666-6666-666666666666', 'dddd0010-0010-0010-0010-000000000010', 'authorized_channel_partner', 'active', 15.0, 15000000, 45000000),
('77777777-7777-7777-7777-777777777777', 'dddd0012-0012-0012-0012-000000000012', 'authorized_channel_partner', 'active', 16.0, 12000000, 35000000),
('88888888-8888-8888-8888-888888888888', 'dddd0014-0014-0014-0014-000000000014', 'regional_channel_partner', 'active', 14.5, 20000000, 60000000),
('99999999-9999-9999-9999-999999999999', 'dddd0016-0016-0016-0016-000000000016', 'authorized_channel_partner', 'active', 17.0, 18000000, 50000000),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'dddd0018-0018-0018-0018-000000000018', 'regional_channel_partner', 'active', 15.5, 10000000, 30000000);

-- ============================================================================
-- SEED DATA COMPLETE: 20 Channel Partners with 21 Facilities and 35+ Relationships
-- ============================================================================
