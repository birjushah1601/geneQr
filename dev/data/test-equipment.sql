-- =============================================================================
-- ServQR Platform - Test Equipment Insertion
-- =============================================================================
-- This file inserts 5 test equipment items using existing category and manufacturer IDs
-- =============================================================================

-- First, let's get some valid category IDs from different sections
DO $$
DECLARE
    dental_chair_cat_id VARCHAR(26);
    xray_cat_id VARCHAR(26);
    microscope_cat_id VARCHAR(26);
    icu_bed_cat_id VARCHAR(26);
    ecg_cat_id VARCHAR(26);
    
    gnatus_mfr_id VARCHAR(26);
    olympus_mfr_id VARCHAR(26);
    skanray_mfr_id VARCHAR(26);
    paramount_mfr_id VARCHAR(26);
    bpl_mfr_id VARCHAR(26);
BEGIN
    -- Get category IDs (selecting first of each type)
    SELECT id INTO dental_chair_cat_id FROM categories WHERE name LIKE '%Dental Chair%' LIMIT 1;
    SELECT id INTO xray_cat_id FROM categories WHERE name LIKE '%X-Ray%' LIMIT 1;
    SELECT id INTO microscope_cat_id FROM categories WHERE name LIKE '%Microscope%' LIMIT 1;
    SELECT id INTO icu_bed_cat_id FROM categories WHERE name LIKE '%ICU Bed%' LIMIT 1;
    SELECT id INTO ecg_cat_id FROM categories WHERE name LIKE '%ECG%' LIMIT 1;
    
    -- Fallback to any categories if specific ones not found
    IF dental_chair_cat_id IS NULL THEN
        SELECT id INTO dental_chair_cat_id FROM categories LIMIT 1;
    END IF;
    
    IF xray_cat_id IS NULL THEN
        SELECT id INTO xray_cat_id FROM categories OFFSET 1 LIMIT 1;
    END IF;
    
    IF microscope_cat_id IS NULL THEN
        SELECT id INTO microscope_cat_id FROM categories OFFSET 2 LIMIT 1;
    END IF;
    
    IF icu_bed_cat_id IS NULL THEN
        SELECT id INTO icu_bed_cat_id FROM categories OFFSET 3 LIMIT 1;
    END IF;
    
    IF ecg_cat_id IS NULL THEN
        SELECT id INTO ecg_cat_id FROM categories OFFSET 4 LIMIT 1;
    END IF;
    
    -- Get manufacturer IDs
    SELECT id INTO gnatus_mfr_id FROM manufacturers WHERE name LIKE '%Gnatus%' LIMIT 1;
    SELECT id INTO olympus_mfr_id FROM manufacturers WHERE name LIKE '%Olympus%' LIMIT 1;
    SELECT id INTO skanray_mfr_id FROM manufacturers WHERE name LIKE '%Skanray%' LIMIT 1;
    SELECT id INTO paramount_mfr_id FROM manufacturers WHERE name LIKE '%Paramount%' LIMIT 1;
    SELECT id INTO bpl_mfr_id FROM manufacturers WHERE name LIKE '%BPL%' LIMIT 1;
    
    -- Fallback to any manufacturers if specific ones not found
    IF gnatus_mfr_id IS NULL THEN
        SELECT id INTO gnatus_mfr_id FROM manufacturers LIMIT 1;
    END IF;
    
    IF olympus_mfr_id IS NULL THEN
        SELECT id INTO olympus_mfr_id FROM manufacturers OFFSET 1 LIMIT 1;
    END IF;
    
    IF skanray_mfr_id IS NULL THEN
        SELECT id INTO skanray_mfr_id FROM manufacturers OFFSET 2 LIMIT 1;
    END IF;
    
    IF paramount_mfr_id IS NULL THEN
        SELECT id INTO paramount_mfr_id FROM manufacturers OFFSET 3 LIMIT 1;
    END IF;
    
    IF bpl_mfr_id IS NULL THEN
        SELECT id INTO bpl_mfr_id FROM manufacturers OFFSET 4 LIMIT 1;
    END IF;
    
    -- Insert test equipment items
    
    -- 1. Dental Chair (demo-hospital)
    INSERT INTO equipment (
        id, name, model, category_id, manufacturer_id, 
        description, specifications, price_amount, 
        price_currency, sku, images, is_active, tenant_id
    ) VALUES (
        '01TEST00000000000000000001', 
        'Test Dental Chair', 
        'DC-100', 
        dental_chair_cat_id, 
        gnatus_mfr_id,
        'Test dental chair for demo purposes',
        '{"chair_positions": "3 programmable", "backrest": "Standard", "headrest": "Adjustable", "warranty": "1 year"}',
        150000.00,
        'INR',
        'TST-DC-100',
        ARRAY['https://abymed.com/images/test/dental-chair-1.jpg'],
        TRUE,
        'demo-hospital'
    );
    
    -- 2. X-Ray System (demo-hospital)
    INSERT INTO equipment (
        id, name, model, category_id, manufacturer_id, 
        description, specifications, price_amount, 
        price_currency, sku, images, is_active, tenant_id
    ) VALUES (
        '01TEST00000000000000000002', 
        'Test X-Ray System', 
        'XR-200', 
        xray_cat_id, 
        skanray_mfr_id,
        'Test X-ray system for demo purposes',
        '{"voltage": "70 kV", "current": "7 mA", "focal_spot": "0.8 mm", "warranty": "2 years"}',
        250000.00,
        'INR',
        'TST-XR-200',
        ARRAY['https://abymed.com/images/test/xray-1.jpg'],
        TRUE,
        'demo-hospital'
    );
    
    -- 3. Laboratory Microscope (city-hospital)
    INSERT INTO equipment (
        id, name, model, category_id, manufacturer_id, 
        description, specifications, price_amount, 
        price_currency, sku, images, is_active, tenant_id
    ) VALUES (
        '01TEST00000000000000000003', 
        'Test Laboratory Microscope', 
        'LM-300', 
        microscope_cat_id, 
        olympus_mfr_id,
        'Test microscope for laboratory use',
        '{"magnification": "40x-1000x", "eyepiece": "10x", "objectives": ["4x", "10x", "40x", "100x"], "warranty": "1 year"}',
        180000.00,
        'INR',
        'TST-LM-300',
        ARRAY['https://abymed.com/images/test/microscope-1.jpg'],
        TRUE,
        'city-hospital'
    );
    
    -- 4. ICU Bed (city-hospital)
    INSERT INTO equipment (
        id, name, model, category_id, manufacturer_id, 
        description, specifications, price_amount, 
        price_currency, sku, images, is_active, tenant_id
    ) VALUES (
        '01TEST00000000000000000004', 
        'Test ICU Bed', 
        'IB-400', 
        icu_bed_cat_id, 
        paramount_mfr_id,
        'Test ICU bed with multiple positions',
        '{"sections": "4-section", "controls": "Electric", "positions": ["Trendelenburg", "Reverse Trendelenburg"], "warranty": "3 years"}',
        320000.00,
        'INR',
        'TST-IB-400',
        ARRAY['https://abymed.com/images/test/icu-bed-1.jpg'],
        TRUE,
        'city-hospital'
    );
    
    -- 5. ECG Machine (demo-hospital)
    INSERT INTO equipment (
        id, name, model, category_id, manufacturer_id, 
        description, specifications, price_amount, 
        price_currency, sku, images, is_active, tenant_id
    ) VALUES (
        '01TEST00000000000000000005', 
        'Test ECG Machine', 
        'EC-500', 
        ecg_cat_id, 
        bpl_mfr_id,
        'Test ECG machine with digital display',
        '{"channels": "12", "display": "7 inch LCD", "memory": "100 records", "battery": "4 hours", "warranty": "2 years"}',
        140000.00,
        'INR',
        'TST-EC-500',
        ARRAY['https://abymed.com/images/test/ecg-1.jpg'],
        TRUE,
        'demo-hospital'
    );
    
    RAISE NOTICE 'Successfully inserted 5 test equipment items';
    RAISE NOTICE 'Categories used: %, %, %, %, %', dental_chair_cat_id, xray_cat_id, microscope_cat_id, icu_bed_cat_id, ecg_cat_id;
    RAISE NOTICE 'Manufacturers used: %, %, %, %, %', gnatus_mfr_id, olympus_mfr_id, skanray_mfr_id, paramount_mfr_id, bpl_mfr_id;
END $$;

-- Verify the insertion
SELECT COUNT(*) AS test_equipment_count FROM equipment WHERE id LIKE '01TEST%';
