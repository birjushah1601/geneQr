-- ============================================================================
-- Add More Engineers and Complete Manufacturer Assignments
-- ============================================================================

-- Add specialized engineers for manufacturers that need them
INSERT INTO engineers (name, phone, email, skills, home_region, engineer_level) VALUES
-- Dialysis specialists for Fresenius
('Sanjay Mehta', '+91-9876543215', 'sanjay.mehta@engineer.com', ARRAY['Dialysis', 'Renal Equipment'], 'Chennai', 2),
('Neha Kulkarni', '+91-9876543216', 'neha.kulkarni@engineer.com', ARRAY['Dialysis', 'Blood Processing'], 'Hyderabad', 3),

-- Additional patient monitor specialists for Medtronic
('Arjun Malhotra', '+91-9876543217', 'arjun.malhotra@engineer.com', ARRAY['Patient Monitor', 'ECG', 'Vital Signs'], 'Delhi', 2),
('Shreya Patel', '+91-9876543218', 'shreya.patel@engineer.com', ARRAY['Patient Monitor', 'Infusion Pump'], 'Mumbai', 2),

-- General equipment specialists
('Karthik Raghavan', '+91-9876543219', 'karthik.r@engineer.com', ARRAY['MRI', 'CT Scanner', 'X-Ray'], 'Bengaluru', 3),
('Divya Krishnan', '+91-9876543220', 'divya.k@engineer.com', ARRAY['Ultrasound', 'ECG'], 'Pune', 1);

-- ============================================================================
-- Assign new engineers to manufacturers
-- ============================================================================

-- Fresenius Medical Care - Dialysis specialists
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Dialysis Equipment Specialist'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Fresenius Medical Care India'
  AND o.org_type = 'manufacturer'
  AND 'Dialysis' = ANY(e.skills);

-- Medtronic India - Add patient monitor specialists
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Patient Monitoring Specialist'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Medtronic India'
  AND o.org_type = 'manufacturer'
  AND ('Patient Monitor' = ANY(e.skills) OR 'Infusion Pump' = ANY(e.skills))
  AND NOT EXISTS (
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id AND m2.org_id = o.id
  );

-- Add general specialists to Global Manufacturer A
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Service Engineer'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Global Manufacturer A'
  AND o.org_type = 'manufacturer'
  AND e.name IN ('Karthik Raghavan', 'Divya Krishnan');

-- Add more engineers to smaller teams
INSERT INTO engineer_org_memberships (engineer_id, org_id, role)
SELECT e.id, o.id, 'Senior Field Engineer'
FROM engineers e
CROSS JOIN organizations o
WHERE o.name = 'Dräger Medical India'
  AND o.org_type = 'manufacturer'
  AND e.name IN ('Karthik Raghavan')
  AND NOT EXISTS (
    SELECT 1 FROM engineer_org_memberships m2
    WHERE m2.engineer_id = e.id AND m2.org_id = o.id
  );

-- ============================================================================
-- Final Statistics
-- ============================================================================

-- Engineer count by manufacturer
SELECT 
    o.name as manufacturer,
    COUNT(DISTINCT m.engineer_id) as engineer_count,
    string_agg(DISTINCT e.name, ', ' ORDER BY e.name) as engineers
FROM organizations o
LEFT JOIN engineer_org_memberships m ON m.org_id = o.id
LEFT JOIN engineers e ON e.id = m.engineer_id
WHERE o.org_type = 'manufacturer'
GROUP BY o.name
ORDER BY engineer_count DESC, o.name;

-- Total summary
SELECT 
    COUNT(DISTINCT e.id) as total_engineers,
    COUNT(DISTINCT m.org_id) as manufacturers_with_engineers,
    COUNT(*) as total_assignments
FROM engineers e
LEFT JOIN engineer_org_memberships m ON m.engineer_id = e.id
JOIN organizations o ON o.id = m.org_id
WHERE o.org_type = 'manufacturer';
