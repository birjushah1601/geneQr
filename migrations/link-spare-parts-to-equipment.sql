-- Link spare parts to equipment catalog entries
-- This creates associations between spare parts and equipment types

-- First, update spare parts to have proper part_type categories for filtering
UPDATE spare_parts_catalog 
SET part_type = 'component' 
WHERE part_type = 'Critical';

UPDATE spare_parts_catalog 
SET part_type = 'component' 
WHERE part_type = 'Important';

UPDATE spare_parts_catalog 
SET part_type = 'consumable' 
WHERE part_type = 'Consumable';

-- Now link X-Ray parts to X-Ray equipment catalog entries
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Core Component' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L3' THEN 'expert'
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        WHEN sp.engineer_level_required = 'L1' THEN 'basic'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'X-Ray' 
  AND sp.category = 'X-Ray'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link CT parts to CT equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Core Component' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L3' THEN 'expert'
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        WHEN sp.engineer_level_required = 'L1' THEN 'basic'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'CT' 
  AND sp.category = 'CT'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link MRI parts to MRI equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory IN ('Gradient System', 'RF System', 'Cooling') THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L3' THEN 'expert'
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        WHEN sp.engineer_level_required = 'L1' THEN 'basic'
        ELSE 'intermediate'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'MRI' 
  AND sp.category = 'MRI'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link Ultrasound parts to Ultrasound equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Transducer' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    'basic' as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'Ultrasound' 
  AND sp.category = 'Ultrasound'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link Ventilator parts to Ventilator equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory IN ('Breathing Circuit', 'Monitoring') THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'Ventilator' 
  AND sp.category = 'Ventilator'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link Patient Monitor parts to Patient Monitor equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Display' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'Patient Monitor' 
  AND sp.category = 'Patient Monitor'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link Dialysis parts to Dialysis equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Blood Circuit' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'Dialysis' 
  AND sp.category = 'Dialysis'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Link Anesthesia parts to Anesthesia equipment
INSERT INTO equipment_spare_parts (equipment_catalog_id, spare_part_id, is_critical, is_consumable, part_category, installation_complexity)
SELECT 
    ec.id as equipment_catalog_id,
    sp.id as spare_part_id,
    CASE WHEN sp.subcategory = 'Vaporizer' THEN true ELSE false END as is_critical,
    CASE WHEN sp.part_type = 'consumable' THEN true ELSE false END as is_consumable,
    sp.subcategory as part_category,
    CASE 
        WHEN sp.engineer_level_required = 'L3' THEN 'expert'
        WHEN sp.engineer_level_required = 'L2' THEN 'intermediate'
        ELSE 'basic'
    END as installation_complexity
FROM equipment_catalog ec
CROSS JOIN spare_parts_catalog sp
WHERE ec.category = 'Anesthesia' 
  AND sp.category = 'Anesthesia'
ON CONFLICT (equipment_catalog_id, spare_part_id) DO NOTHING;

-- Display summary
SELECT 
    ec.category as equipment_type,
    COUNT(DISTINCT esp.spare_part_id) as parts_linked,
    COUNT(DISTINCT esp.equipment_catalog_id) as equipment_models
FROM equipment_spare_parts esp
JOIN equipment_catalog ec ON esp.equipment_catalog_id = ec.id
GROUP BY ec.category
ORDER BY ec.category;

SELECT COUNT(*) as "Total Equipment-Part Associations" FROM equipment_spare_parts;
