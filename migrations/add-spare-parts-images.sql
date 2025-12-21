-- Add sample images and videos for spare parts
-- Using placeholder images from various sources

-- X-Ray Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=800',
        'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800'
    ]
WHERE part_number = 'XR-TUBE-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=800',
        'https://images.unsplash.com/photo-1559757175-5700dde675bc?w=800'
    ]
WHERE part_number = 'XR-DET-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'XR-COL-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'XR-FILT-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=800'
    ]
WHERE part_number = 'XR-GRID-001';

-- CT Scanner Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1516549655169-df83a0774514?w=800',
        'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800'
    ]
WHERE part_number = 'CT-TUBE-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=800'
    ]
WHERE part_number = 'CT-DET-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'CT-SLIP-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'CT-COL-001';

-- MRI Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=800',
        'https://images.unsplash.com/photo-1582719471137-c3967ffb1c42?w=800'
    ]
WHERE part_number = 'MRI-COIL-HEAD';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1582719471137-c3967ffb1c42?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1582719471137-c3967ffb1c42?w=800'
    ]
WHERE part_number = 'MRI-COIL-BODY';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=800',
        'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=800'
    ]
WHERE part_number = 'MRI-GRAD-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'MRI-CRYO-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'MRI-RF-AMP';

-- Ultrasound Parts Images  
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=800'
    ]
WHERE part_number = 'US-PROBE-C60';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1631815588090-d4bfec5b1ccb?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1631815588090-d4bfec5b1ccb?w=800'
    ]
WHERE part_number = 'US-PROBE-L38';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number = 'US-GEL-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=800'
    ]
WHERE part_number = 'US-BATT-001';

-- Ventilator Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=800'
    ]
WHERE part_number IN ('VENT-VALVE-001', 'VENT-VALVE-002');

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number IN ('VENT-SENS-O2', 'VENT-SENS-CO2');

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number = 'VENT-FILT-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'VENT-TUBE-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=800'
    ]
WHERE part_number = 'VENT-BATT-001';

-- Patient Monitor Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1631815589968-fdb09a223b1e?w=800'
    ]
WHERE part_number = 'PM-ECG-CABLE';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'PM-SPO2-SENSOR';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number = 'PM-NIBP-CUFF';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'PM-TEMP-PROBE';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1631815588090-d4bfec5b1ccb?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1631815588090-d4bfec5b1ccb?w=800'
    ]
WHERE part_number = 'PM-IBP-CABLE';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1609069700247-0a6a8c2f6535?w=800'
    ]
WHERE part_number = 'PM-BATT-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1530497610245-94d3c16cda28?w=800',
        'https://images.unsplash.com/photo-1559757175-5700dde675bc?w=800'
    ]
WHERE part_number = 'PM-DISPLAY-001';

-- Dialysis Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=800'
    ]
WHERE part_number = 'DIAL-PUMP-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number = 'DIAL-FILT-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581093588401-fbb62a02f120?w=800'
    ]
WHERE part_number = 'DIAL-LINE-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number IN ('DIAL-CONC-BIC', 'DIAL-CONC-ACID');

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'DIAL-PRES-001';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=800'
    ]
WHERE part_number = 'DIAL-VALVE-001';

-- Anesthesia Parts Images
UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800',
        'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=800'
    ]
WHERE part_number IN ('ANES-VAPOR-ISO', 'ANES-VAPOR-SEV');

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1585435557343-3b092031a831?w=800'
    ]
WHERE part_number = 'ANES-CO2-ABS';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=800'
    ]
WHERE part_number = 'ANES-O2-SENS';

UPDATE spare_parts_catalog 
SET image_url = 'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=400',
    photos = ARRAY[
        'https://images.unsplash.com/photo-1584982751601-97dcc096659c?w=800'
    ]
WHERE part_number = 'ANES-BELLOW';

-- Show updated counts
SELECT 
    COUNT(*) as total_parts,
    COUNT(image_url) as with_images,
    COUNT(photos) as with_photo_arrays
FROM spare_parts_catalog;

-- Show sample with images
SELECT part_number, part_name, category, image_url 
FROM spare_parts_catalog 
WHERE image_url IS NOT NULL 
LIMIT 10;
