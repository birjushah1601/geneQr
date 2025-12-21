-- ============================================================================
-- Equipment Catalog Parts Assignments for Demo
-- Purpose: Link spare parts to equipment catalog items for comprehensive demo
-- ============================================================================

-- Clear existing assignments (if any)
DELETE FROM equipment_part_assignments;

-- Get equipment catalog IDs by category for reference
-- MRI: 2 items
-- CT: 2 items
-- X-Ray: 1 item
-- Ultrasound: 1 item
-- Patient Monitor: 1 item
-- Ventilator: 1 item
-- Anesthesia: 2 items
-- Dialysis: 1 item
-- Infusion Pump: 1 item
-- Laboratory: 1 item

-- ============================================================================
-- MRI SCANNER PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'MRI-COIL-HEAD-8CH' THEN 2
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 50
    END as quantity_required,
    CASE 
        WHEN p.part_number = 'MRI-COIL-HEAD-8CH' THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'MRI-COIL-HEAD-8CH' THEN 1095  -- 3 years
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 90  -- 3 months
    END as replacement_frequency_days,
    CASE 
        WHEN p.part_number = 'MRI-COIL-HEAD-8CH' THEN 'medium'
        ELSE 'low'
    END as installation_complexity,
    CASE 
        WHEN p.part_number = 'MRI-CRYOGEN-HELIUM' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'MRI'
  AND p.part_number IN ('MRI-COIL-HEAD-8CH', 'MRI-CRYOGEN-HELIUM');

-- ============================================================================
-- CT SCANNER PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'CT-TUBE-XRAY' THEN 1
        WHEN p.part_number = 'CT-DETECTOR-MODULE' THEN 1
    END as quantity_required,
    true as is_critical,
    false as is_consumable,
    CASE 
        WHEN p.part_number = 'CT-TUBE-XRAY' THEN 730  -- 2 years
        WHEN p.part_number = 'CT-DETECTOR-MODULE' THEN 1825  -- 5 years
    END as replacement_frequency_days,
    'high' as installation_complexity,
    'scheduled' as replacement_type
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'CT'
  AND p.part_number IN ('CT-TUBE-XRAY', 'CT-DETECTOR-MODULE');

-- ============================================================================
-- X-RAY SYSTEM PARTS
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
    ec.id as equipment_catalog_id,
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
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'X-Ray'
  AND p.part_number IN ('XR-TUBE-STANDARD', 'XR-PANEL-FLAT');

-- ============================================================================
-- ULTRASOUND MACHINE PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'US-PROBE-LINEAR' THEN 1
        WHEN p.part_number = 'US-PROBE-CONVEX' THEN 1
    END as quantity_required,
    true as is_critical,
    false as is_consumable,
    1095 as replacement_frequency_days,  -- 3 years
    'medium' as installation_complexity,
    'scheduled' as replacement_type
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'Ultrasound'
  AND p.part_number IN ('US-PROBE-LINEAR', 'US-PROBE-CONVEX');

-- ============================================================================
-- PATIENT MONITOR PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'PM-ECG-CABLE' THEN 3
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 10
    END as quantity_required,
    CASE 
        WHEN p.part_number = 'PM-ECG-CABLE' THEN true
        ELSE false
    END as is_critical,
    CASE 
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN true
        ELSE false
    END as is_consumable,
    CASE 
        WHEN p.part_number = 'PM-ECG-CABLE' THEN 365  -- 1 year
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 180  -- 6 months
    END as replacement_frequency_days,
    'low' as installation_complexity,
    CASE 
        WHEN p.part_number = 'PM-SPO2-SENSOR' THEN 'routine'
        ELSE 'scheduled'
    END as replacement_type
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'Patient Monitor'
  AND p.part_number IN ('PM-ECG-CABLE', 'PM-SPO2-SENSOR');

-- ============================================================================
-- VENTILATOR PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'VEN-FILTER-HEPA' THEN 6
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
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'Ventilator'
  AND p.part_number IN ('VEN-FILTER-HEPA', 'VEN-SENSOR-FLOW');

-- ============================================================================
-- DIALYSIS MACHINE PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'DIA-FILTER-DIALYZER' THEN 20
        WHEN p.part_number = 'DIA-TUBING-SET' THEN 50
    END as quantity_required,
    false as is_critical,
    true as is_consumable,
    30 as replacement_frequency_days,  -- Monthly
    'low' as installation_complexity,
    'routine' as replacement_type
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'Dialysis'
  AND p.part_number IN ('DIA-FILTER-DIALYZER', 'DIA-TUBING-SET');

-- ============================================================================
-- INFUSION PUMP PARTS
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
    ec.id as equipment_catalog_id,
    p.id as spare_part_id,
    CASE 
        WHEN p.part_number = 'INF-TUBING-IV' THEN 100
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
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog p
WHERE ec.category = 'Infusion Pump'
  AND p.part_number IN ('INF-TUBING-IV', 'INF-BATTERY-PACK');

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================

-- Show summary of parts per equipment catalog
SELECT 
    ec.product_name,
    ec.category,
    COUNT(DISTINCT epa.spare_part_id) as total_parts_assigned,
    COUNT(DISTINCT CASE WHEN epa.is_critical THEN epa.spare_part_id END) as critical_parts,
    COUNT(DISTINCT CASE WHEN epa.is_consumable THEN epa.spare_part_id END) as consumable_parts,
    SUM(epa.quantity_required) as total_quantity
FROM equipment_catalog ec
LEFT JOIN equipment_part_assignments epa ON ec.id = epa.equipment_catalog_id
GROUP BY ec.id, ec.product_name, ec.category
ORDER BY ec.category, ec.product_name;

-- Show total assignments
SELECT 
    COUNT(*) as total_assignments,
    COUNT(DISTINCT equipment_catalog_id) as equipment_with_parts,
    COUNT(DISTINCT spare_part_id) as unique_parts_used
FROM equipment_part_assignments;

-- Show parts breakdown by category
SELECT 
    ec.category,
    COUNT(DISTINCT epa.spare_part_id) as unique_parts,
    SUM(epa.quantity_required) as total_quantity,
    COUNT(CASE WHEN epa.is_critical THEN 1 END) as critical_assignments
FROM equipment_catalog ec
JOIN equipment_part_assignments epa ON ec.id = epa.equipment_catalog_id
GROUP BY ec.category
ORDER BY ec.category;
