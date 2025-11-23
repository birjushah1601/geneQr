-- ============================================================================
-- SEED DATA: ENGINEER ASSIGNMENT SYSTEM (MINIMAL VERSION)
-- ============================================================================
-- Create sample engineers matching the actual schema
-- ============================================================================

-- Insert Engineers (Simple)
INSERT INTO engineers (id, name, phone, email, skills, home_region, engineer_level) VALUES
  (gen_random_uuid(), 'Rajesh Kumar Singh', '+91-98765-43210', 'rajesh.singh@siemens.com', ARRAY['MRI', 'CT Scanner'], 'Mumbai', 3),
  (gen_random_uuid(), 'Priya Sharma', '+91-98765-43211', 'priya.sharma@siemens.com', ARRAY['X-Ray'], 'Delhi', 2),
  (gen_random_uuid(), 'Arun Menon', '+91-98765-43230', 'arun.menon@philips.com', ARRAY['MRI'], 'Chennai', 3),
  (gen_random_uuid(), 'Vikram Reddy', '+91-98765-43220', 'vikram.reddy@ge.com', ARRAY['CT Scanner'], 'Hyderabad', 3),
  (gen_random_uuid(), 'Suresh Gupta', '+91-98765-43250', 'suresh.gupta@dealer.com', ARRAY['X-Ray', 'Ultrasound'], 'Delhi', 2)
ON CONFLICT (id) DO NOTHING;

-- Link Engineers to Organizations
DO $$
DECLARE
  siemens_org_id UUID;
  philips_org_id UUID;
  ge_org_id UUID;
  dealer_org_id UUID;
  rajesh_id UUID;
  priya_id UUID;
  arun_id UUID;
  vikram_id UUID;
  suresh_id UUID;
BEGIN
  -- Get organization IDs
  SELECT id INTO siemens_org_id FROM organizations WHERE name LIKE '%Siemens%' LIMIT 1;
  SELECT id INTO philips_org_id FROM organizations WHERE name LIKE '%Philips%' LIMIT 1;
  SELECT id INTO ge_org_id FROM organizations WHERE name LIKE '%GE%' OR name LIKE '%Wipro%' LIMIT 1;
  SELECT id INTO dealer_org_id FROM organizations WHERE org_type = 'dealer' LIMIT 1;

  -- Get engineer IDs
  SELECT id INTO rajesh_id FROM engineers WHERE email = 'rajesh.singh@siemens.com';
  SELECT id INTO priya_id FROM engineers WHERE email = 'priya.sharma@siemens.com';
  SELECT id INTO arun_id FROM engineers WHERE email = 'arun.menon@philips.com';
  SELECT id INTO vikram_id FROM engineers WHERE email = 'vikram.reddy@ge.com';
  SELECT id INTO suresh_id FROM engineers WHERE email = 'suresh.gupta@dealer.com';

  -- Create org memberships
  IF siemens_org_id IS NOT NULL THEN
    IF rajesh_id IS NOT NULL THEN
      INSERT INTO engineer_org_memberships (engineer_id, org_id, role) 
      VALUES (rajesh_id, siemens_org_id, 'field_engineer') 
      ON CONFLICT DO NOTHING;
      
      INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, is_certified)
      VALUES 
        (rajesh_id, 'Siemens Healthineers', 'MRI', true),
        (rajesh_id, 'Siemens Healthineers', 'CT Scanner', true)
      ON CONFLICT DO NOTHING;
    END IF;

    IF priya_id IS NOT NULL THEN
      INSERT INTO engineer_org_memberships (engineer_id, org_id, role) 
      VALUES (priya_id, siemens_org_id, 'field_engineer') 
      ON CONFLICT DO NOTHING;
      
      INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, is_certified)
      VALUES (priya_id, 'Siemens Healthineers', 'X-Ray', true)
      ON CONFLICT DO NOTHING;
    END IF;
  END IF;

  IF philips_org_id IS NOT NULL AND arun_id IS NOT NULL THEN
    INSERT INTO engineer_org_memberships (engineer_id, org_id, role) 
    VALUES (arun_id, philips_org_id, 'field_engineer') 
    ON CONFLICT DO NOTHING;
    
    INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, is_certified)
    VALUES (arun_id, 'Philips Healthcare', 'MRI', true)
    ON CONFLICT DO NOTHING;
  END IF;

  IF ge_org_id IS NOT NULL AND vikram_id IS NOT NULL THEN
    INSERT INTO engineer_org_memberships (engineer_id, org_id, role) 
    VALUES (vikram_id, ge_org_id, 'field_engineer') 
    ON CONFLICT DO NOTHING;
    
    INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, is_certified)
    VALUES (vikram_id, 'GE Healthcare', 'CT Scanner', true)
    ON CONFLICT DO NOTHING;
  END IF;

  IF dealer_org_id IS NOT NULL AND suresh_id IS NOT NULL THEN
    INSERT INTO engineer_org_memberships (engineer_id, org_id, role) 
    VALUES (suresh_id, dealer_org_id, 'field_engineer') 
    ON CONFLICT DO NOTHING;
    
    INSERT INTO engineer_equipment_types (engineer_id, manufacturer_name, equipment_category, is_certified)
    VALUES 
      (suresh_id, 'Siemens Healthineers', 'X-Ray', true),
      (suresh_id, 'GE Healthcare', 'X-Ray', true)
    ON CONFLICT DO NOTHING;
  END IF;

  RAISE NOTICE 'Engineers and memberships created';
END $$;

-- Verify
SELECT 
  'Total engineers:' as info, 
  COUNT(*)::TEXT as count 
FROM engineers WHERE engineer_level IS NOT NULL
UNION ALL
SELECT 
  'Engineer equipment types:', 
  COUNT(*)::TEXT 
FROM engineer_equipment_types
UNION ALL
SELECT 
  'Engineer org memberships:', 
  COUNT(*)::TEXT 
FROM engineer_org_memberships
UNION ALL
SELECT 
  'Equipment service configs:', 
  COUNT(*)::TEXT 
FROM equipment_service_config;
