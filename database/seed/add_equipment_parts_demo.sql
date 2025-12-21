-- ============================================================================
-- Equipment Parts Assignments for Demo
-- Purpose: Link spare parts to all equipment for comprehensive demo
-- ============================================================================

-- Clear existing assignments (if any)
DELETE FROM equipment_part_assignments;

-- Get equipment IDs by category for reference
-- CT Scanners: 2 units
-- ECG Machines: 2 units  
-- Infusion Pumps: 2 units
-- Ventilators: 2 units
-- X-Ray Systems: 2 units

-- ============================================================================
-- CT SCANNER PARTS (2 units)
-- ============================================================================

-- CT Scanner parts assignments
INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'CT-TUBE-XRAY' THEN 1
        WHEN p.part_number = 'CT-DETECTOR-MODULE' THEN 1
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 20  -- Consumable quantity
    END as quantity_required,
    CASE 
        WHEN p.part_number IN ('CT-TUBE-XRAY', 'CT-DETECTOR-MODULE') THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'CT-TUBE-XRAY' THEN 730  -- 2 years
        WHEN p.part_number = 'CT-DETECTOR-MODULE' THEN 1825  -- 5 years
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 90  -- 3 months
    END as replacement_frequency_days,
    CASE 
        WHEN p.part_number IN ('CT-TUBE-XRAY', 'CT-DETECTOR-MODULE') THEN 'high'
        ELSE 'low'
    END as installation_complexity,
    CASE 
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'CT'
  AND p.part_number IN ('CT-TUBE-XRAY', 'CT-DETECTOR-MODULE', 'MRI-CRYOGEN-HELIUM');

-- ============================================================================
-- X-RAY SYSTEM PARTS (2 units)
-- ============================================================================

INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'XR-TUBE-STANDARD' THEN 1
        WHEN p.part_number = 'XR-PANEL-FLAT' THEN 1
    END as quantity_required,
    true as is_critical,
    false as is_consumable,
    CASE 
        WHEN p.part_number = 'XR-TUBE-STANDARD' THEN 1095  -- 3 years
        WHEN p.part_number = 'XR-PANEL-FLAT' THEN 1825  -- 5 years
    END as replacement_frequency_days,
    'high' as installation_complexity,
    'scheduled' as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'X-Ray'
  AND p.part_number IN ('XR-TUBE-STANDARD', 'XR-PANEL-FLAT');

-- ============================================================================
-- ULTRASOUND/ECG MACHINE PARTS (2 units)
-- ============================================================================

INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'PM-ECG-CABLE' THEN 2
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 5
        WHEN p.part_number = 'US-PROBE-LINEAR' THEN 1
        WHEN p.part_number = 'US-PROBE-CONVEX' THEN 1
    END as quantity_required,
    CASE 
        WHEN p.part_number IN ('PM-ECG-CABLE', 'US-PROBE-LINEAR', 'US-PROBE-CONVEX') THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'PM-ECG-CABLE' THEN 365  -- 1 year
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 180  -- 6 months
        WHEN p.part_number IN ('US-PROBE-LINEAR', 'US-PROBE-CONVEX') THEN 1095  -- 3 years
    END as replacement_frequency_days,
    CASE 
        WHEN p.part_number IN ('US-PROBE-LINEAR', 'US-PROBE-CONVEX') THEN 'medium'
        ELSE 'low'
    END as installation_complexity,
    CASE 
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'ECG'
  AND p.part_number IN ('PM-ECG-CABLE', 'PM-SPO2-SENSOR', 'US-PROBE-LINEAR', 'US-PROBE-CONVEX');

-- ============================================================================
-- VENTILATOR PARTS (2 units)
-- ============================================================================

INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'VEN-FILTER-HEPA' THEN 4
        WHEN p.part_number = 'VEN-SENSOR-FLOW' THEN 1
    END as quantity_required,
    CASE 
        WHEN p.part_number = 'VEN-SENSOR-FLOW' THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'VEN-FILTER-HEPA' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'VEN-FILTER-HEPA' THEN 30  -- Monthly
        WHEN p.part_number = 'VEN-SENSOR-FLOW' THEN 730  -- 2 years
    END as replacement_frequency_days,
    CASE 
        WHEN p.part_number = 'VEN-SENSOR-FLOW' THEN 'medium'
        ELSE 'low'
    END as installation_complexity,
    CASE 
        WHEN p.part_number = 'VEN-FILTER-HEPA' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'Ventilator'
  AND p.part_number IN ('VEN-FILTER-HEPA', 'VEN-SENSOR-FLOW');

-- ============================================================================
-- INFUSION PUMP PARTS (2 units)
-- ============================================================================

INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'INF-TUBING-IV' THEN 100  -- Consumable stock
        WHEN p.part_number = 'INF-BATTERY-PACK' THEN 2
    END as quantity_required,
    CASE 
        WHEN p.part_number = 'INF-BATTERY-PACK' THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'INF-TUBING-IV' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'INF-TUBING-IV' THEN 7  -- Weekly
        WHEN p.part_number = 'INF-BATTERY-PACK' THEN 365  -- Yearly
    END as replacement_frequency_days,
    'low' as installation_complexity,
    CASE 
        WHEN p.part_number = 'INF-TUBING-IV' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'Infusion Pump'
  AND p.part_number IN ('INF-TUBING-IV', 'INF-BATTERY-PACK');

-- ============================================================================
-- ADDITIONAL COMMON PARTS (Apply to multiple equipment types)
-- ============================================================================

-- Add MRI Head Coil to CT and X-Ray equipment (as optional accessory)
INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_optional,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    1 as quantity_required,
    false as is_critical,
    true as is_optional,
    'low' as installation_complexity,
    'on_demand' as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category IN ('CT', 'X-Ray')
  AND p.part_number = 'MRI-COIL-HEAD-8CH';

-- Add Dialyzer components to Infusion Pumps (optional upgrades)
INSERT INTO equipment_part_assignments (
    equipment_catalog_id,
    spare_part_id,
    quantity_required,
    is_critical,
    is_consumable,
    is_optional,
    replacement_frequency_days,
    installation_complexity,
    replacement_type
)
SELECT 
    e.id::uuid as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'DIA-FILTER-DIALYZER' THEN 10
        WHEN p.part_number = 'DIA-TUBING-SET' THEN 50
    END as quantity_required,
    false as is_critical,
    true as is_consumable,
    true as is_optional,
    30 as replacement_frequency_days,
    'low' as installation_complexity,
    'routine' as replacement_type
FROM equipment e
CROSS JOIN spare_parts_catalog p
WHERE e.category = 'Infusion Pump'
  AND p.part_number IN ('DIA-FILTER-DIALYZER', 'DIA-TUBING-SET');

-- ============================================================================
-- VERIFICATION QUERY
-- ============================================================================

-- Show summary of parts per equipment
SELECT 
    e.equipment_name,
    e.category,
    COUNT(DISTINCT epa.spare_part_id) as total_parts_assigned,
    COUNT(DISTINCT CASE WHEN epa.is_critical THEN epa.spare_part_id END) as critical_parts,
    COUNT(DISTINCT CASE WHEN epa.is_consumable THEN epa.spare_part_id END) as consumable_parts,
    SUM(epa.quantity_required) as total_quantity
FROM equipment e
LEFT JOIN equipment_part_assignments epa ON e.id::uuid = epa.equipment_catalog_id
GROUP BY e.id, e.equipment_name, e.category
ORDER BY e.category, e.equipment_name;

-- Show total assignments
SELECT 
    COUNT(*) as total_assignments,
    COUNT(DISTINCT equipment_catalog_id) as equipment_with_parts,
    COUNT(DISTINCT spare_part_id) as unique_parts_used
FROM equipment_part_assignments;
