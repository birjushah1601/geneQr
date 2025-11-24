-- Migration: 022-data-migration-seeding.sql
-- Description: Data Migration & Seeding for New Tables
-- Ticket: T2B.7
-- Date: 2025-11-16
--
-- This migration:
-- 1. Seeds equipment catalog with common medical equipment
-- 2. Seeds common parts for each equipment type
-- 3. Migrates existing engineer expertise data
-- 4. Creates default workflow templates
-- 5. Seeds manufacturer service configurations
-- 6. Creates sample data for testing

-- =====================================================================
-- 1. SEED EQUIPMENT CATALOG
-- =====================================================================

-- Common Medical Equipment Categories
INSERT INTO equipment_catalog (
    catalog_name, equipment_category, equipment_type, description,
    typical_use_cases, maintenance_requirements, is_active
) VALUES
-- Imaging Equipment
('CT Scanner', 'imaging', 'diagnostic', 'Computed Tomography Scanner for detailed internal imaging',
 ARRAY['Emergency diagnosis', 'Cancer detection', 'Trauma assessment'], 
 ARRAY['Daily calibration', 'Monthly preventive maintenance', 'Annual certification'],
 true),

('MRI Machine', 'imaging', 'diagnostic', 'Magnetic Resonance Imaging for soft tissue visualization',
 ARRAY['Brain imaging', 'Joint analysis', 'Tumor detection'],
 ARRAY['Daily helium check', 'Weekly gradient calibration', 'Quarterly system check'],
 true),

('X-Ray Machine', 'imaging', 'diagnostic', 'Radiography equipment for bone and basic imaging',
 ARRAY['Bone fractures', 'Chest X-rays', 'Dental imaging'],
 ARRAY['Weekly calibration', 'Monthly tube inspection', 'Annual radiation safety check'],
 true),

('Ultrasound', 'imaging', 'diagnostic', 'Ultrasound imaging for real-time visualization',
 ARRAY['Obstetrics', 'Cardiology', 'Abdominal imaging'],
 ARRAY['Daily probe cleaning', 'Monthly transducer check', 'Annual calibration'],
 true),

-- Life Support Equipment
('Ventilator', 'life_support', 'critical_care', 'Mechanical ventilation for respiratory support',
 ARRAY['ICU patients', 'Surgery', 'Emergency respiratory support'],
 ARRAY['Daily leak test', 'Weekly calibration', 'Monthly filter replacement'],
 true),

('Patient Monitor', 'monitoring', 'critical_care', 'Vital signs monitoring system',
 ARRAY['ICU monitoring', 'OR monitoring', 'Post-op care'],
 ARRAY['Daily zero calibration', 'Weekly probe check', 'Monthly system test'],
 true),

('Defibrillator', 'life_support', 'emergency', 'Cardiac defibrillation device',
 ARRAY['Cardiac arrest', 'Emergency response', 'Surgery standby'],
 ARRAY['Daily battery check', 'Weekly self-test', 'Monthly electrode check'],
 true),

('Infusion Pump', 'drug_delivery', 'critical_care', 'Controlled medication and fluid delivery',
 ARRAY['IV medication', 'Chemotherapy', 'Pain management'],
 ARRAY['Daily accuracy check', 'Weekly occlusion test', 'Monthly calibration'],
 true),

-- Laboratory Equipment
('Blood Gas Analyzer', 'laboratory', 'diagnostic', 'Blood gas and electrolyte analysis',
 ARRAY['ICU diagnostics', 'Emergency lab', 'Surgery monitoring'],
 ARRAY['Daily calibration', 'Weekly QC', 'Monthly maintenance'],
 true),

('Centrifuge', 'laboratory', 'sample_processing', 'Sample separation equipment',
 ARRAY['Blood separation', 'Urine analysis', 'Research'],
 ARRAY['Weekly speed check', 'Monthly balance verification', 'Annual motor service'],
 true),

-- Surgical Equipment
('Surgical Light', 'surgical', 'or_equipment', 'Operating room illumination system',
 ARRAY['Surgery lighting', 'Examination lighting', 'Procedure lighting'],
 ARRAY['Weekly bulb check', 'Monthly cleaning', 'Annual electrical safety test'],
 true),

('Anesthesia Machine', 'anesthesia', 'critical_care', 'Anesthesia delivery and monitoring',
 ARRAY['Surgery anesthesia', 'Sedation', 'Pain management'],
 ARRAY['Daily leak test', 'Weekly vaporizer check', 'Monthly calibration'],
 true);

-- =====================================================================
-- 2. SEED COMMON PARTS FOR EQUIPMENT
-- =====================================================================

-- CT Scanner Parts
INSERT INTO equipment_parts (
    equipment_catalog_id, part_number, part_name, part_category, description,
    is_oem_part, typical_lifespan_months, stock_status
)
SELECT 
    id, 
    part_number,
    part_name,
    part_category,
    description,
    is_oem,
    lifespan,
    'in_stock'
FROM equipment_catalog ec
CROSS JOIN LATERAL (
    VALUES
        ('CT-TUBE-001', 'X-Ray Tube', 'primary_component', 'High voltage X-ray generation tube', true, 36),
        ('CT-DET-001', 'Detector Array', 'primary_component', 'Multi-slice detector system', true, 60),
        ('CT-FIL-001', 'Air Filter', 'consumable', 'HEPA air filter for cooling system', false, 6),
        ('CT-COOL-001', 'Cooling System Pump', 'maintenance_part', 'Tube cooling circulation pump', true, 24)
) AS parts(part_number, part_name, part_category, description, is_oem, lifespan)
WHERE ec.catalog_name = 'CT Scanner';

-- Ventilator Parts
INSERT INTO equipment_parts (
    equipment_catalog_id, part_number, part_name, part_category, description,
    is_oem_part, typical_lifespan_months, stock_status
)
SELECT 
    id,
    part_number,
    part_name,
    part_category,
    description,
    is_oem,
    lifespan,
    'in_stock'
FROM equipment_catalog ec
CROSS JOIN LATERAL (
    VALUES
        ('VENT-TUBE-001', 'High-Flow Breathing Tube', 'consumable', 'ICU-grade breathing circuit', true, 3),
        ('VENT-TUBE-002', 'Standard Breathing Tube', 'consumable', 'General ward breathing circuit', false, 3),
        ('VENT-VALVE-001', 'Exhalation Valve', 'maintenance_part', 'Pressure control exhalation valve', true, 12),
        ('VENT-FIL-001', 'Bacterial Filter', 'consumable', 'Disposable bacterial/viral filter', true, 1),
        ('VENT-SENS-001', 'Flow Sensor', 'sensor', 'Airflow measurement sensor', true, 18),
        ('VENT-BATT-001', 'Backup Battery', 'maintenance_part', 'Emergency power battery pack', true, 24)
) AS parts(part_number, part_name, part_category, description, is_oem, lifespan)
WHERE ec.catalog_name = 'Ventilator';

-- Patient Monitor Parts
INSERT INTO equipment_parts (
    equipment_catalog_id, part_number, part_name, part_category, description,
    is_oem_part, typical_lifespan_months, stock_status
)
SELECT 
    id,
    part_number,
    part_name,
    part_category,
    description,
    is_oem,
    lifespan,
    'in_stock'
FROM equipment_catalog ec
CROSS JOIN LATERAL (
    VALUES
        ('MON-ECG-001', 'ECG Lead Set', 'consumable', '5-lead ECG cable set', true, 12),
        ('MON-SPO2-001', 'SpO2 Sensor Adult', 'consumable', 'Pulse oximetry finger sensor', true, 6),
        ('MON-SPO2-002', 'SpO2 Sensor Pediatric', 'consumable', 'Pediatric pulse oximetry sensor', true, 6),
        ('MON-NIBP-001', 'NIBP Cuff Adult', 'consumable', 'Non-invasive blood pressure cuff', true, 12),
        ('MON-BATT-001', 'Battery Pack', 'maintenance_part', 'Rechargeable battery module', true, 24),
        ('MON-SCREEN-001', 'Display Screen', 'primary_component', 'LCD touchscreen display', true, 60)
) AS parts(part_number, part_name, part_category, description, is_oem, lifespan)
WHERE ec.catalog_name = 'Patient Monitor';

-- Defibrillator Parts
INSERT INTO equipment_parts (
    equipment_catalog_id, part_number, part_name, part_category, description,
    is_oem_part, typical_lifespan_months, stock_status
)
SELECT 
    id,
    part_number,
    part_name,
    part_category,
    description,
    is_oem,
    lifespan,
    'in_stock'
FROM equipment_catalog ec
CROSS JOIN LATERAL (
    VALUES
        ('DEFIB-PAD-001', 'Defibrillation Pads Adult', 'consumable', 'Disposable defib electrode pads', true, 24),
        ('DEFIB-PAD-002', 'Defibrillation Pads Pediatric', 'consumable', 'Pediatric defib pads', true, 24),
        ('DEFIB-BATT-001', 'Battery Pack', 'maintenance_part', 'High-capacity battery module', true, 36),
        ('DEFIB-CAP-001', 'Capacitor', 'primary_component', 'High voltage capacitor', true, 60)
) AS parts(part_number, part_name, part_category, description, is_oem, lifespan)
WHERE ec.catalog_name = 'Defibrillator';

-- =====================================================================
-- 3. SEED EQUIPMENT PARTS CONTEXT (ICU vs General Ward)
-- =====================================================================

-- Ventilator breathing tubes - ICU context
INSERT INTO equipment_parts_context (
    part_id, installation_context, context_priority, context_specific_notes
)
SELECT 
    ep.id,
    'ICU',
    'required',
    'High-flow tubes mandatory for ICU ventilators due to critical patient requirements'
FROM equipment_parts ep
JOIN equipment_catalog ec ON ep.equipment_catalog_id = ec.id
WHERE ec.catalog_name = 'Ventilator'
  AND ep.part_name = 'High-Flow Breathing Tube';

-- Standard tubes for general ward
INSERT INTO equipment_parts_context (
    part_id, installation_context, context_priority, context_specific_notes
)
SELECT 
    ep.id,
    'General Ward',
    'preferred',
    'Standard tubes sufficient for general ward ventilators'
FROM equipment_parts ep
JOIN equipment_catalog ec ON ep.equipment_catalog_id = ec.id
WHERE ec.catalog_name = 'Ventilator'
  AND ep.part_name = 'Standard Breathing Tube';

-- =====================================================================
-- 4. MIGRATE ENGINEER EXPERTISE (From existing engineers table)
-- =====================================================================

-- This assumes engineers exist with some equipment-related data
-- Create L2 expertise for all engineers on basic equipment
INSERT INTO engineer_equipment_expertise (
    engineer_id, equipment_catalog_id, support_level,
    can_diagnose, can_repair, can_install, can_train,
    expertise_notes
)
SELECT 
    e.id,
    ec.id,
    'L2',
    true,
    true,
    false,
    false,
    'Default L2 expertise assignment - requires validation'
FROM engineers e
CROSS JOIN equipment_catalog ec
WHERE ec.catalog_name IN ('Patient Monitor', 'Infusion Pump')
  AND e.status = 'active'
ON CONFLICT (engineer_id, equipment_catalog_id) DO NOTHING;

-- =====================================================================
-- 5. CREATE DEFAULT WORKFLOW TEMPLATES
-- =====================================================================

-- Default Workflow for General Equipment Service
INSERT INTO workflow_templates (
    template_name, template_code, template_type, description,
    total_stages, is_default, is_active
) VALUES (
    'Standard Service Workflow',
    'STANDARD-SERVICE-V1',
    'default',
    'Default multi-stage service workflow for general equipment',
    5,
    true,
    true
);

-- Get the template ID
DO $$
DECLARE
    v_template_id UUID;
    v_stage_1_id UUID;
    v_stage_2_id UUID;
    v_stage_3_id UUID;
    v_stage_4_id UUID;
    v_stage_5_id UUID;
BEGIN
    SELECT id INTO v_template_id 
    FROM workflow_templates 
    WHERE template_code = 'STANDARD-SERVICE-V1';

    -- Stage 1: Remote Diagnosis
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_number, stage_name, stage_type,
        stage_description, is_mandatory, target_duration_hours,
        requires_customer_approval, requires_parts
    ) VALUES (
        v_template_id, 1, 'Remote Diagnosis', 'diagnosis',
        'Initial remote assessment and diagnosis',
        true, 4, false, false
    ) RETURNING id INTO v_stage_1_id;

    -- Stage 2: Parts Identification
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_number, stage_name, stage_type,
        stage_description, is_mandatory, target_duration_hours,
        requires_customer_approval, requires_parts
    ) VALUES (
        v_template_id, 2, 'Parts Identification', 'assessment',
        'Identify and request required parts',
        true, 2, false, false
    ) RETURNING id INTO v_stage_2_id;

    -- Stage 3: Parts Procurement
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_number, stage_name, stage_type,
        stage_description, is_mandatory, target_duration_hours,
        requires_customer_approval, requires_parts
    ) VALUES (
        v_template_id, 3, 'Parts Procurement', 'waiting',
        'Order and receive required parts',
        false, 48, false, true
    ) RETURNING id INTO v_stage_3_id;

    -- Stage 4: Onsite Repair
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_number, stage_name, stage_type,
        stage_description, is_mandatory, target_duration_hours,
        requires_customer_approval, requires_parts
    ) VALUES (
        v_template_id, 4, 'Onsite Repair', 'repair',
        'Onsite equipment repair and testing',
        true, 4, false, true
    ) RETURNING id INTO v_stage_4_id;

    -- Stage 5: Completion & Documentation
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_number, stage_name, stage_type,
        stage_description, is_mandatory, target_duration_hours,
        requires_customer_approval, requires_parts
    ) VALUES (
        v_template_id, 5, 'Completion & Sign-off', 'completion',
        'Final testing, documentation, and customer sign-off',
        true, 2, true, false
    ) RETURNING id INTO v_stage_5_id;

    -- Create stage transitions
    INSERT INTO stage_transition_rules (
        from_stage_id, to_stage_id, transition_condition,
        auto_transition, is_active
    ) VALUES
        (v_stage_1_id, v_stage_2_id, 'diagnosis_complete', false, true),
        (v_stage_2_id, v_stage_3_id, 'parts_identified', true, true),
        (v_stage_3_id, v_stage_4_id, 'parts_available', false, true),
        (v_stage_4_id, v_stage_5_id, 'repair_complete', false, true);

    RAISE NOTICE 'Default workflow template created with 5 stages';
END $$;

-- =====================================================================
-- 6. SEED MANUFACTURER SERVICE CONFIGURATIONS
-- =====================================================================

-- Create default manufacturer service configs
-- This assumes manufacturers table has data
INSERT INTO manufacturer_service_config (
    manufacturer_id, equipment_catalog_id,
    service_model, who_provides_service, who_provides_parts,
    service_config, is_active
)
SELECT 
    m.id,
    ec.id,
    'manufacturer_service',
    'manufacturer',
    'manufacturer',
    jsonb_build_object(
        'response_time_hours', 24,
        'parts_lead_time_days', 7,
        'certification_required', true,
        'uses_oem_parts_only', true
    ),
    true
FROM manufacturers m
CROSS JOIN equipment_catalog ec
WHERE m.name IN ('Philips', 'GE Healthcare', 'Siemens')
  AND ec.catalog_name IN ('CT Scanner', 'MRI Machine')
ON CONFLICT DO NOTHING;

-- GeneQR-provided service for general equipment
INSERT INTO manufacturer_service_config (
    manufacturer_id, equipment_catalog_id,
    service_model, who_provides_service, who_provides_parts,
    service_config, is_active
)
SELECT 
    m.id,
    ec.id,
    'geneqr_service',
    'geneqr',
    'geneqr',
    jsonb_build_object(
        'response_time_hours', 4,
        'parts_lead_time_days', 2,
        'certification_required', false,
        'can_use_compatible_parts', true
    ),
    true
FROM manufacturers m
CROSS JOIN equipment_catalog ec
WHERE ec.catalog_name IN ('Patient Monitor', 'Infusion Pump', 'Defibrillator')
ON CONFLICT DO NOTHING;

-- =====================================================================
-- 7. SAMPLE TEST DATA (for development/testing only)
-- =====================================================================

-- Insert sample customers (if not exists)
INSERT INTO customers (
    name, customer_code, organization_type, primary_contact_name,
    primary_contact_phone, primary_contact_email, status
) VALUES
    ('City General Hospital', 'CGH-001', 'government_hospital', 'Dr. Rajesh Kumar',
     '+91-9876543210', 'rajesh.kumar@cityhospital.in', 'active'),
    ('Apollo Hospitals', 'APOLLO-001', 'private_hospital', 'Ms. Priya Sharma',
     '+91-9876543211', 'priya.sharma@apollo.in', 'active'),
    ('AIIMS Delhi', 'AIIMS-001', 'government_hospital', 'Dr. Amit Verma',
     '+91-9876543212', 'amit.verma@aiims.in', 'active')
ON CONFLICT (customer_code) DO NOTHING;

-- Insert sample equipment
INSERT INTO equipment (
    customer_id, equipment_name, serial_number, model_number,
    catalog_id, manufacturer_id, installation_location, status
)
SELECT 
    c.id,
    ec.catalog_name || ' - ' || c.name,
    'SN-' || substr(md5(random()::text), 1, 10),
    'MODEL-' || substr(md5(random()::text), 1, 8),
    ec.id,
    m.id,
    CASE 
        WHEN ec.catalog_name IN ('Ventilator', 'Patient Monitor') THEN 'ICU Ward 3'
        WHEN ec.catalog_name IN ('CT Scanner', 'MRI Machine') THEN 'Radiology Department'
        ELSE 'General Equipment Room'
    END,
    'operational'
FROM customers c
CROSS JOIN equipment_catalog ec
CROSS JOIN manufacturers m
WHERE c.customer_code IN ('CGH-001', 'APOLLO-001', 'AIIMS-001')
  AND ec.catalog_name IN ('Ventilator', 'Patient Monitor', 'CT Scanner')
  AND m.name IN ('Philips', 'GE Healthcare')
LIMIT 20
ON CONFLICT DO NOTHING;

-- =====================================================================
-- 8. COMPATIBILITY MATRIX
-- =====================================================================

-- Create compatibility between ventilator parts
INSERT INTO equipment_compatibility (
    part_id, compatible_with_part_id, compatibility_type,
    compatibility_notes, is_interchangeable
)
SELECT 
    ep1.id,
    ep2.id,
    'alternative',
    'Can be used as alternative in non-critical situations',
    true
FROM equipment_parts ep1
JOIN equipment_parts ep2 ON ep1.equipment_catalog_id = ep2.equipment_catalog_id
JOIN equipment_catalog ec ON ep1.equipment_catalog_id = ec.id
WHERE ec.catalog_name = 'Ventilator'
  AND ep1.part_name = 'High-Flow Breathing Tube'
  AND ep2.part_name = 'Standard Breathing Tube'
ON CONFLICT DO NOTHING;

-- =====================================================================
-- MIGRATION STATISTICS
-- =====================================================================

DO $$
DECLARE
    v_catalog_count INT;
    v_parts_count INT;
    v_customers_count INT;
    v_equipment_count INT;
    v_workflow_count INT;
BEGIN
    SELECT COUNT(*) INTO v_catalog_count FROM equipment_catalog;
    SELECT COUNT(*) INTO v_parts_count FROM equipment_parts;
    SELECT COUNT(*) INTO v_customers_count FROM customers;
    SELECT COUNT(*) INTO v_equipment_count FROM equipment;
    SELECT COUNT(*) INTO v_workflow_count FROM workflow_templates;

    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Migration 022 complete!';
    RAISE NOTICE '==============================================';
    RAISE NOTICE 'Data Seeding Summary:';
    RAISE NOTICE '  - Equipment Catalog: % entries', v_catalog_count;
    RAISE NOTICE '  - Equipment Parts: % entries', v_parts_count;
    RAISE NOTICE '  - Sample Customers: % entries', v_customers_count;
    RAISE NOTICE '  - Sample Equipment: % entries', v_equipment_count;
    RAISE NOTICE '  - Workflow Templates: % templates', v_workflow_count;
    RAISE NOTICE '';
    RAISE NOTICE 'Seeded Equipment Types:';
    RAISE NOTICE '  - CT Scanner, MRI, X-Ray, Ultrasound';
    RAISE NOTICE '  - Ventilator, Patient Monitor, Defibrillator';
    RAISE NOTICE '  - Infusion Pump, Blood Gas Analyzer';
    RAISE NOTICE '  - Centrifuge, Surgical Light, Anesthesia Machine';
    RAISE NOTICE '';
    RAISE NOTICE 'Ready for production use!';
    RAISE NOTICE '==============================================';
END $$;
