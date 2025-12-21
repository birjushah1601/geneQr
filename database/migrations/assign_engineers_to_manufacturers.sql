-- ============================================================================
-- Assign Engineers to Manufacturers
-- ============================================================================
-- This migration assigns existing engineers to manufacturers based on their
-- skills and equipment specializations
-- ============================================================================

-- Clear existing manufacturer memberships to start fresh
DELETE FROM engineer_org_memberships 
WHERE org_id IN (SELECT id FROM organizations WHERE org_type = 'manufacturer');

-- ============================================================================
-- Siemens Healthineers India - MRI and CT specialists
-- ============================================================================
-- Assign engineers with MRI and CT skills to Siemens
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Field Service Engineer'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Siemens Healthineers India'
  AND o.org_type = 'manufacturer'
  AND (
    'MRI' = ANY(e.skills) OR 
    'CT Scanner' = ANY(e.skills)
  );

-- ============================================================================
-- Philips Healthcare India - Patient Monitor and Ultrasound specialists
-- ============================================================================
-- Assign engineers with Ultrasound and general equipment skills
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Service Technician'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Philips Healthcare India'
  AND o.org_type = 'manufacturer'
  AND (
    'Ultrasound' = ANY(e.skills) OR
    'MRI' = ANY(e.skills)
  )
  AND NOT EXISTS (
    -- Don't duplicate if already assigned
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id AND m2.org_id = o.id
  );

-- ============================================================================
-- Wipro GE Healthcare - CT and Ultrasound specialists
-- ============================================================================
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Field Service Engineer'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Wipro GE Healthcare'
  AND o.org_type = 'manufacturer'
  AND (
    'CT Scanner' = ANY(e.skills) OR
    'Ultrasound' = ANY(e.skills)
  )
  AND NOT EXISTS (
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id AND m2.org_id = o.id
  );

-- ============================================================================
-- Canon Medical Systems India - X-Ray specialists
-- ============================================================================
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'X-Ray Technician'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Canon Medical Systems India'
  AND o.org_type = 'manufacturer'
  AND 'X-Ray' = ANY(e.skills);

-- ============================================================================
-- Dräger Medical India - Ventilator and Critical Care specialists
-- ============================================================================
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Critical Care Specialist'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Dräger Medical India'
  AND o.org_type = 'manufacturer'
  AND (
    'Ventilator' = ANY(e.skills) OR
    'ECG' = ANY(e.skills)
  );

-- ============================================================================
-- Medtronic India - General equipment specialists (CT, MRI, X-Ray)
-- ============================================================================
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Biomedical Engineer'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Medtronic India'
  AND o.org_type = 'manufacturer'
  AND (
    'CT Scanner' = ANY(e.skills) OR
    'MRI' = ANY(e.skills) OR
    'X-Ray' = ANY(e.skills)
  )
  AND NOT EXISTS (
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id AND m2.org_id = o.id
  )
LIMIT 2; -- Limit to 2 engineers for Medtronic

-- ============================================================================
-- Fresenius Medical Care India - Add some general engineers
-- ============================================================================
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Dialysis Equipment Specialist'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Fresenius Medical Care India'
  AND o.org_type = 'manufacturer'
  AND NOT EXISTS (
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id
  )
LIMIT 1; -- Assign 1 unassigned engineer

-- ============================================================================
-- Verification and Stats
-- ============================================================================

-- Show engineer assignments by manufacturer
SELECT 
    o.name as manufacturer,
    COUNT(DISTINCT m.engineer_id) as engineer_count,
    string_agg(DISTINCT m.role, ', ') as roles
FROM organizations o
LEFT JOIN engineer_org_memberships m ON m.org_id = o.id
WHERE o.org_type = 'manufacturer'
GROUP BY o.name
ORDER BY engineer_count DESC, o.name;

-- Show which engineers are assigned to which manufacturers
SELECT 
    e.name as engineer_name,
    e.skills,
    o.name as manufacturer,
    m.role
FROM engineer_org_memberships m
JOIN engineers e ON e.id = m.engineer_id
JOIN organizations o ON o.id = m.org_id
WHERE o.org_type = 'manufacturer'
ORDER BY o.name, e.name;
