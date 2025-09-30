-- =============================================================================
-- ABY-Med Platform - Comprehensive Medical Equipment Catalog (Simplified)
-- =============================================================================
-- This file contains a comprehensive catalog of 60 medical equipment items
-- covering dental, laboratory, hospital infrastructure, emergency care,
-- and specialized medical equipment with realistic Indian manufacturers.
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS public;

-- Create equipment table if it doesn't exist
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(100) NOT NULL,
    category_id VARCHAR(26) NOT NULL,
    manufacturer_id VARCHAR(26) NOT NULL,
    description TEXT,
    specifications JSONB NOT NULL,
    price_amount DECIMAL(12, 2) NOT NULL,
    price_currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    sku VARCHAR(50),
    images TEXT[],
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tenant_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(id)
);

-- =============================================================================
-- 1. DENTAL EQUIPMENT (15 items) - DEMO HOSPITAL TENANT
-- =============================================================================

INSERT INTO equipment (id, name, model, category_id, manufacturer_id, description, specifications, price_amount, price_currency, sku, images, is_active, tenant_id)
VALUES
-- Dental Chairs
('01HFPQ3Z5P8VF5PXZRT4K7MHG1', 'Gnatus Dental Chair', 'G3 New', '01HFPQ2Z5P8VF5PXZRT4K7MHG1', '01HFPQ2Z5P8VF5PXZRT4K7MHG1', 'Premium dental chair with integrated delivery system', 
'{"chair_positions": "5 programmable", "backrest": "Ultra-thin", "headrest": "Articulated", "delivery_system": "Over-patient", "light": "LED 35,000 lux", "warranty": "3 years"}', 
385000.00, 'INR', 'GNT-DC-G3N-001', 
ARRAY['https://abymed.com/images/equipment/dental/gnatus-g3-new-1.jpg'], 
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHG2', 'Confident Dental Chair', 'Chesa Neo', '01HFPQ2Z5P8VF5PXZRT4K7MHG1', '01HFPQ2Z5P8VF5PXZRT4K7MHG2', 'Ergonomic dental chair with seamless upholstery',
'{"chair_positions": "4 programmable", "backrest": "Contoured", "headrest": "Double-articulated", "delivery_system": "Side delivery", "light": "LED 30,000 lux", "warranty": "2 years"}',
295000.00, 'INR', 'CNF-DC-CN-001',
ARRAY['https://abymed.com/images/equipment/dental/confident-chesa-neo-1.jpg'],
TRUE, 'demo-hospital'),

-- Dental X-Ray Systems
('01HFPQ3Z5P8VF5PXZRT4K7MHG6', 'Planmeca ProMax 3D', 'Classic', '01HFPQ2Z5P8VF5PXZRT4K7MHG2', '01HFPQ2Z5P8VF5PXZRT4K7MHG6', 'Advanced panoramic and 3D dental imaging system',
'{"imaging_modes": ["Panoramic", "Cephalometric", "3D CBCT"], "field_of_view": "8x8 cm", "resolution": "75-600 μm", "voltage": "60-90 kV", "warranty": "2 years"}',
2850000.00, 'INR', 'PLM-DX-PC-001',
ARRAY['https://abymed.com/images/equipment/dental/planmeca-promax-3d-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHG8', 'Allengers Intraoral X-Ray', 'Mars-70', '01HFPQ2Z5P8VF5PXZRT4K7MHG2', '01HFPQ2Z5P8VF5PXZRT4K7MHG8', 'Wall-mounted intraoral X-ray unit',
'{"imaging_mode": "Intraoral", "tube_voltage": "70 kV", "tube_current": "7 mA", "exposure_time": "0.02-3.2 sec", "arm_reach": "165 cm", "warranty": "3 years"}',
125000.00, 'INR', 'ALG-DX-M70-001',
ARRAY['https://abymed.com/images/equipment/dental/allengers-mars70-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGA', 'Skanray Intraoral X-Ray', 'Numx Plus', '01HFPQ2Z5P8VF5PXZRT4K7MHG2', '01HFPQ2Z5P8VF5PXZRT4K7MHGA', 'Indian-made intraoral X-ray unit',
'{"imaging_mode": "Intraoral", "tube_voltage": "65 kV", "tube_current": "6 mA", "exposure_time": "0.02-2.0 sec", "arm_reach": "150 cm", "warranty": "2 years"}',
85000.00, 'INR', 'SKN-DX-NP-001',
ARRAY['https://abymed.com/images/equipment/dental/skanray-numx-1.jpg'],
TRUE, 'demo-hospital'),

-- Dental Sterilizers
('01HFPQ3Z5P8VF5PXZRT4K7MHGB', 'Tuttnauer Autoclave', '3870EA', '01HFPQ2Z5P8VF5PXZRT4K7MHG3', '01HFPQ2Z5P8VF5PXZRT4K7MHGB', 'Automatic autoclave with closed door drying',
'{"chamber_size": "85 liters", "temperature_range": "105-138°C", "cycle_options": ["Unwrapped", "Wrapped", "Liquid"], "drying_system": "Closed door", "warranty": "2 years"}',
385000.00, 'INR', 'TTN-DS-3870-001',
ARRAY['https://abymed.com/images/equipment/dental/tuttnauer-3870ea-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGC', 'Confident Autoclave', 'Class B 22L', '01HFPQ2Z5P8VF5PXZRT4K7MHG3', '01HFPQ2Z5P8VF5PXZRT4K7MHG2', 'Class B autoclave with vacuum system',
'{"chamber_size": "22 liters", "temperature_range": "121-134°C", "cycle_options": ["Unwrapped", "Wrapped", "Prion"], "vacuum_system": "Pre and post vacuum", "warranty": "1 year"}',
195000.00, 'INR', 'CNF-DS-B22-001',
ARRAY['https://abymed.com/images/equipment/dental/confident-autoclave-1.jpg'],
TRUE, 'demo-hospital'),

-- Dental Handpieces
('01HFPQ3Z5P8VF5PXZRT4K7MHGF', 'NSK High-Speed Handpiece', 'Pana-Max Plus', '01HFPQ2Z5P8VF5PXZRT4K7MHG4', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Premium high-speed dental handpiece',
'{"speed": "Up to 400,000 rpm", "spray_type": "Quattro spray", "light_source": "LED", "bearing_type": "Ceramic", "sterilization": "Autoclavable 135°C", "warranty": "1 year"}',
28500.00, 'INR', 'NSK-DH-PMP-001',
ARRAY['https://abymed.com/images/equipment/dental/nsk-panamax-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGG', 'Being Dental Handpiece', 'Rose-H2', '01HFPQ2Z5P8VF5PXZRT4K7MHG4', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Affordable high-speed handpiece',
'{"speed": "Up to 350,000 rpm", "spray_type": "Triple spray", "light_source": "None", "bearing_type": "Steel", "sterilization": "Autoclavable 135°C", "warranty": "6 months"}',
12500.00, 'INR', 'BNG-DH-RH2-001',
ARRAY['https://abymed.com/images/equipment/dental/being-rose-1.jpg'],
TRUE, 'demo-hospital'),

-- Dental Lights & Compressors
('01HFPQ3Z5P8VF5PXZRT4K7MHGI', 'Faro Dental Light', 'Maia LED', '01HFPQ2Z5P8VF5PXZRT4K7MHG6', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Advanced LED dental operatory light',
'{"light_source": "LED", "intensity": "3,000-35,000 lux", "color_temp": "5,000K", "control": "No-touch sensor", "power": "9W", "warranty": "2 years"}',
125000.00, 'INR', 'FAR-DL-MLED-001',
ARRAY['https://abymed.com/images/equipment/dental/faro-maia-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGJ', 'Anand Dental Compressor', 'Oilless 2HP', '01HFPQ2Z5P8VF5PXZRT4K7MHG5', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Oil-free dental air compressor',
'{"type": "Oil-free", "motor": "2 HP", "air_displacement": "340 L/min", "tank": "60 liters", "noise": "65 dB", "warranty": "2 years"}',
85000.00, 'INR', 'AND-DC-O2HP-001',
ARRAY['https://abymed.com/images/equipment/dental/anand-compressor-1.jpg'],
TRUE, 'demo-hospital'),

-- Dental CAD/CAM & Lasers
('01HFPQ3Z5P8VF5PXZRT4K7MHGL', 'Planmeca CAD/CAM System', 'PlanMill 40 S', '01HFPQ2Z5P8VF5PXZRT4K7MHGA4', '01HFPQ2Z5P8VF5PXZRT4K7MHG6', 'Complete chairside CAD/CAM system',
'{"components": ["Scanner", "Design Software", "Milling Unit"], "milling_axes": "4-axis", "materials": ["Ceramics", "Composites", "Zirconia"], "warranty": "2 years"}',
3500000.00, 'INR', 'PLM-CAD-PM40S-001',
ARRAY['https://abymed.com/images/equipment/dental/planmeca-planmill-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGN', 'Biolase Dental Laser', 'Waterlase iPlus', '01HFPQ2Z5P8VF5PXZRT4K7MHGA5', '01HFPQ2Z5P8VF5PXZRT4K7MHGN', 'All-tissue dental laser',
'{"laser_type": "Er,Cr:YSGG", "wavelength": "2780 nm", "power": "0.1-10.0 W", "applications": ["Hard Tissue", "Soft Tissue", "Bone"], "warranty": "2 years"}',
2500000.00, 'INR', 'BIO-DL-WIP-001',
ARRAY['https://abymed.com/images/equipment/dental/biolase-waterlase-1.jpg'],
TRUE, 'demo-hospital'),

-- Additional Dental Equipment
('01HFPQ3Z5P8VF5PXZRT4K7MHGQ', 'Woodpecker Ultrasonic Scaler', 'UDS-K LED', '01HFPQ2Z5P8VF5PXZRT4K7MHG9', '01HFPQ2Z5P8VF5PXZRT4K7MHGQ', 'Piezoelectric ultrasonic scaler',
'{"technology": "Piezoelectric", "frequency": "28 kHz", "power_settings": "10 levels", "handpiece": "LED light", "tips": "5 scaling tips", "warranty": "1 year"}',
18500.00, 'INR', 'WDP-DS-UDSK-001',
ARRAY['https://abymed.com/images/equipment/dental/woodpecker-udsk-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGR', 'Coltene Endo Motor', 'CanalPro Jeni', '01HFPQ2Z5P8VF5PXZRT4K7MHG8', '01HFPQ2Z5P8VF5PXZRT4K7MHGR', 'Intelligent endodontic motor',
'{"technology": "Digital motor with AI", "speed": "200-1000 rpm", "torque": "0.1-5.0 Ncm", "apex_locator": "Integrated", "battery": "Rechargeable", "warranty": "2 years"}',
125000.00, 'INR', 'CLT-DE-CPJ-001',
ARRAY['https://abymed.com/images/equipment/dental/coltene-canalpro-1.jpg'],
TRUE, 'demo-hospital'),

-- =============================================================================
-- 2. LABORATORY EQUIPMENT (15 items) - DEMO HOSPITAL TENANT
-- =============================================================================

-- Microscopes
('01HFPQ3Z5P8VF5PXZRT4K7MHGU', 'Olympus Compound Microscope', 'CX43', '01HFPQ2Z5P8VF5PXZRT4K7MHGB1', '01HFPQ2Z5P8VF5PXZRT4K7MHGU', 'Professional biological microscope',
'{"optical_system": "Infinity-corrected", "eyepiece": "10x", "objectives": ["4x", "10x", "40x", "100x"], "illumination": "LED", "warranty": "3 years"}',
185000.00, 'INR', 'OLY-LM-CX43-001',
ARRAY['https://abymed.com/images/equipment/lab/olympus-cx43-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGV', 'Labomed Digital Microscope', 'Lx500', '01HFPQ2Z5P8VF5PXZRT4K7MHGB1', '01HFPQ2Z5P8VF5PXZRT4K7MHGV', 'Digital microscope with camera',
'{"optical_system": "Infinity-corrected", "eyepiece": "10x", "objectives": ["4x", "10x", "40x", "100x"], "camera": "3MP CMOS", "warranty": "2 years"}',
225000.00, 'INR', 'LBM-LM-LX500-001',
ARRAY['https://abymed.com/images/equipment/lab/labomed-lx500-1.jpg'],
TRUE, 'demo-hospital'),

-- Centrifuges
('01HFPQ3Z5P8VF5PXZRT4K7MHGY', 'Remi Microcentrifuge', 'RM-12C Plus', '01HFPQ2Z5P8VF5PXZRT4K7MHGB2', '01HFPQ2Z5P8VF5PXZRT4K7MHGY', 'High-speed microcentrifuge',
'{"max_speed": "15,000 rpm", "max_rcf": "21,380 x g", "capacity": "24 x 1.5/2.0 ml", "timer": "0-99 min", "programs": "9", "warranty": "2 years"}',
125000.00, 'INR', 'RMI-LC-RM12C-001',
ARRAY['https://abymed.com/images/equipment/lab/remi-micro-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHGZ', 'Thermo Scientific Clinical Centrifuge', 'Medifuge', '01HFPQ2Z5P8VF5PXZRT4K7MHGB2', '01HFPQ2Z5P8VF5PXZRT4K7MHGZ', 'Clinical centrifuge for blood samples',
'{"max_speed": "4,500 rpm", "max_rcf": "2,490 x g", "capacity": "24 x 15 ml", "timer": "0-99 min", "safety": ["Imbalance detection", "Auto-lock"], "warranty": "2 years"}',
185000.00, 'INR', 'THS-LC-MDF-001',
ARRAY['https://abymed.com/images/equipment/lab/thermo-medifuge-1.jpg'],
TRUE, 'demo-hospital'),

-- Autoclaves & Sterilizers
('01HFPQ3Z5P8VF5PXZRT4K7MHH3', 'Tuttnauer Laboratory Autoclave', '3870ELV', '01HFPQ2Z5P8VF5PXZRT4K7MHGB3', '01HFPQ2Z5P8VF5PXZRT4K7MHGB', 'Vertical laboratory autoclave',
'{"chamber_size": "85 liters", "temperature": "105-138°C", "cycles": ["Unwrapped", "Wrapped", "Liquid", "Agar"], "drying": "Vacuum-assisted", "warranty": "2 years"}',
450000.00, 'INR', 'TTN-LS-3870ELV-001',
ARRAY['https://abymed.com/images/equipment/lab/tuttnauer-3870elv-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHH4', 'Equitron Autoclave', 'EQ-75L', '01HFPQ2Z5P8VF5PXZRT4K7MHGB3', '01HFPQ2Z5P8VF5PXZRT4K7MHH4', 'Indian-made vertical autoclave',
'{"chamber_size": "75 liters", "temperature": "121-134°C", "cycles": ["Standard", "Liquid"], "control": "Digital PID", "warranty": "1 year"}',
125000.00, 'INR', 'EQT-LS-75L-001',
ARRAY['https://abymed.com/images/equipment/lab/equitron-75l-1.jpg'],
TRUE, 'demo-hospital'),

-- Incubators & Spectrophotometers
('01HFPQ3Z5P8VF5PXZRT4K7MHH6', 'Thermo Scientific CO2 Incubator', 'Heracell 150i', '01HFPQ2Z5P8VF5PXZRT4K7MHGB4', '01HFPQ2Z5P8VF5PXZRT4K7MHGZ', 'CO2 incubator for cell culture',
'{"chamber_size": "150 liters", "temperature": "Ambient +3°C to 55°C", "CO2_range": "0-20%", "humidity": "Up to 95% RH", "sterilization": "180°C dry heat", "warranty": "2 years"}',
850000.00, 'INR', 'THS-LI-H150I-001',
ARRAY['https://abymed.com/images/equipment/lab/thermo-heracell-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHH7', 'Remi Bacteriological Incubator', 'RI-175', '01HFPQ2Z5P8VF5PXZRT4K7MHGB4', '01HFPQ2Z5P8VF5PXZRT4K7MHGY', 'General purpose incubator',
'{"chamber_size": "175 liters", "temperature": "Ambient +5°C to 70°C", "accuracy": "±0.5°C", "shelves": "3 adjustable", "circulation": "Natural convection", "warranty": "1 year"}',
65000.00, 'INR', 'RMI-LI-RI175-001',
ARRAY['https://abymed.com/images/equipment/lab/remi-incubator-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHH8', 'Shimadzu UV-Vis Spectrophotometer', 'UV-1800', '01HFPQ2Z5P8VF5PXZRT4K7MHGB5', '01HFPQ2Z5P8VF5PXZRT4K7MHH8', 'Double-beam UV-Visible spectrophotometer',
'{"wavelength": "190-1100 nm", "bandwidth": "1 nm", "accuracy": "±0.1 nm", "photometric_range": "-4 to +4 Abs", "light_source": "Deuterium and Tungsten", "warranty": "1 year"}',
650000.00, 'INR', 'SHM-LS-UV1800-001',
ARRAY['https://abymed.com/images/equipment/lab/shimadzu-uv1800-1.jpg'],
TRUE, 'demo-hospital'),

-- Analytical Balances
('01HFPQ3Z5P8VF5PXZRT4K7MHHA', 'Mettler Toledo Analytical Balance', 'XS204', '01HFPQ2Z5P8VF5PXZRT4K7MHGB6', '01HFPQ2Z5P8VF5PXZRT4K7MHHA', 'High-precision analytical balance',
'{"capacity": "220 g", "readability": "0.1 mg", "repeatability": "0.1 mg", "calibration": "Internal automatic", "display": "Touchscreen", "warranty": "2 years"}',
450000.00, 'INR', 'MTT-LB-XS204-001',
ARRAY['https://abymed.com/images/equipment/lab/mettler-xs204-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHHC', 'Wensar Precision Balance', 'PGB-300', '01HFPQ2Z5P8VF5PXZRT4K7MHGB6', '01HFPQ2Z5P8VF5PXZRT4K7MHHC', 'Indian-made precision balance',
'{"capacity": "300 g", "readability": "1 mg", "repeatability": "2 mg", "calibration": "External", "display": "LCD with backlight", "warranty": "1 year"}',
85000.00, 'INR', 'WNS-LB-PGB300-001',
ARRAY['https://abymed.com/images/equipment/lab/wensar-pgb300-1.jpg'],
TRUE, 'demo-hospital'),

-- Water Baths, Shakers, PCR
('01HFPQ3Z5P8VF5PXZRT4K7MHHF', 'Thermo Scientific Water Bath', 'Precision 28', '01HFPQ2Z5P8VF5PXZRT4K7MHGB7', '01HFPQ2Z5P8VF5PXZRT4K7MHGZ', 'Digital laboratory water bath',
'{"capacity": "28 liters", "temperature": "Ambient +5°C to 100°C", "uniformity": "±0.1°C at 37°C", "display": "LCD", "safety": ["Over-temp cutoff", "Low-water alarm"], "warranty": "2 years"}', 
135000.00, 'INR', 'THS-LW-P28-001',
ARRAY['https://abymed.com/images/equipment/lab/thermo-waterbath-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHHG', 'Remi Orbital Shaker', 'OS-10', '01HFPQ2Z5P8VF5PXZRT4K7MHGB8', '01HFPQ2Z5P8VF5PXZRT4K7MHGY', 'Bench-top orbital shaker',
'{"speed": "40-300 rpm", "orbit": "20 mm", "capacity": "7.5 kg", "timer": "0-99 hrs", "platform": "330 x 330 mm", "warranty": "1 year"}',
65000.00, 'INR', 'RMI-LS-OS10-001',
ARRAY['https://abymed.com/images/equipment/lab/remi-shaker-1.jpg'],
TRUE, 'demo-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHHH', 'Bio-Rad Gradient PCR', 'T100', '01HFPQ2Z5P8VF5PXZRT4K7MHGB9', '01HFPQ2Z5P8VF5PXZRT4K7MHHH', '96-well gradient thermal cycler',
'{"capacity": "96 x 0.2 ml", "temperature": "4-100 °C", "gradient": "30 °C", "interface": "Touchscreen", "storage": "500 methods", "warranty": "2 years"}',
385000.00, 'INR', 'BRD-LP-T100-001',
ARRAY['https://abymed.com/images/equipment/lab/biorad-t100-1.jpg'],
TRUE, 'demo-hospital'),

-- =============================================================================
-- 3. HOSPITAL INFRASTRUCTURE (10 items) - CITY HOSPITAL TENANT
-- =============================================================================

-- Hospital Beds
('01HFPQ3Z5P8VF5PXZRT4K7MHHJ', 'Paramount ICU Bed', 'Eleganza 5', '01HFPQ2Z5P8VF5PXZRT4K7MHGC1', '01HFPQ2Z5P8VF5PXZRT4K7MHHJ', 'Fully-electric ICU bed with weighing system',
'{"sections": "4-section platform", "controls": "Bedside & nurse panel", "positions": ["Trendelenburg", "Reverse Trendelenburg"], "load": "250 kg", "battery": "Yes", "warranty": "5 years"}',
325000.00, 'INR', 'PRM-HI-EG5-001',
ARRAY['https://abymed.com/images/equipment/infra/paramount-eleganza5-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHI1', 'Midmark Hospital Bed', 'RG-7000', '01HFPQ2Z5P8VF5PXZRT4K7MHGC1', '01HFPQ2Z5P8VF5PXZRT4K7MHI1', 'Semi-electric hospital bed',
'{"sections": "3-section platform", "controls": "Hand pendant", "positions": ["Hi-Lo", "Fowler"], "load": "200 kg", "side_rails": "Collapsible", "warranty": "3 years"}',
175000.00, 'INR', 'MDM-HI-RG7000-001',
ARRAY['https://abymed.com/images/equipment/infra/midmark-rg7000-1.jpg'],
TRUE, 'city-hospital'),

-- OT Tables
('01HFPQ3Z5P8VF5PXZRT4K7MHHK', 'Hospitech OT Table', 'C-MAX Pro', '01HFPQ2Z5P8VF5PXZRT4K7MHGC2', '01HFPQ2Z5P8VF5PXZRT4K7MHHK', 'Electro-hydraulic operating table',
'{"tabletop": "Radiolucent, modular", "height": "600-1000 mm", "tilt": "±20°", "trendelenburg": "±30°", "load": "350 kg", "warranty": "3 years"}',
495000.00, 'INR', 'HST-HI-CMP-001',
ARRAY['https://abymed.com/images/equipment/infra/hospitech-cmax-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHI2', 'Narang OT Table', 'Surgimax', '01HFPQ2Z5P8VF5PXZRT4K7MHGC2', '01HFPQ2Z5P8VF5PXZRT4K7MHI2', 'Hydraulic operating table',
'{"tabletop": "5-section", "height": "750-950 mm", "tilt": "±15°", "trendelenburg": "±25°", "load": "200 kg", "warranty": "2 years"}',
285000.00, 'INR', 'NRG-HI-SGM-001',
ARRAY['https://abymed.com/images/equipment/infra/narang-surgimax-1.jpg'],
TRUE, 'city-hospital'),

-- OT Lights
('01HFPQ3Z5P8VF5PXZRT4K7MHI3', 'Mindray OT Light', 'HyLED X8', '01HFPQ2Z5P8VF5PXZRT4K7MHGC3', '01HFPQ2Z5P8VF5PXZRT4K7MHI3', 'Surgical LED light with camera',
'{"illumination": "160,000 lux", "color_temp": "3500-5000K", "diameter": "70 cm", "camera": "HD integrated", "control": "Wall panel & wireless", "warranty": "2 years"}',
750000.00, 'INR', 'MND-HI-HLX8-001',
ARRAY['https://abymed.com/images/equipment/infra/mindray-hyled-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHI4', 'BPL OT Light', 'SurgLED 500', '01HFPQ2Z5P8VF5PXZRT4K7MHGC3', '01HFPQ2Z5P8VF5PXZRT4K7MHHM', 'Indian-made surgical LED light',
'{"illumination": "120,000 lux", "color_temp": "4300K", "diameter": "60 cm", "control": "Wall panel", "sterilizable_handle": "Yes", "warranty": "2 years"}',
350000.00, 'INR', 'BPL-HI-SL500-001',
ARRAY['https://abymed.com/images/equipment/infra/bpl-surgled-1.jpg'],
TRUE, 'city-hospital'),

-- Hospital Furniture
('01HFPQ3Z5P8VF5PXZRT4K7MHI5', 'Godrej Hospital Furniture', 'Premium Ward Set', '01HFPQ2Z5P8VF5PXZRT4K7MHGC4', '01HFPQ2Z5P8VF5PXZRT4K7MHI5', 'Complete ward furniture set',
'{"includes": ["Bed", "Bedside Table", "Over-bed Table", "Visitor Chair"], "material": "Powder-coated steel", "finish": "Anti-bacterial", "warranty": "5 years"}',
125000.00, 'INR', 'GDJ-HI-PWS-001',
ARRAY['https://abymed.com/images/equipment/infra/godrej-ward-1.jpg'],
TRUE, 'city-hospital'),

-- HVAC & Air Purification
('01HFPQ3Z5P8VF5PXZRT4K7MHI6', 'Blue Star HVAC System', 'MediClean Series', '01HFPQ2Z5P8VF5PXZRT4K7MHGC5', '01HFPQ2Z5P8VF5PXZRT4K7MHI6', 'Hospital-grade HVAC system',
'{"capacity": "20 TR", "filtration": "HEPA H14", "air_changes": "15-20 per hour", "controls": "Digital touchscreen", "monitoring": "Remote IoT", "warranty": "3 years"}',
1850000.00, 'INR', 'BLS-HI-MCS-001',
ARRAY['https://abymed.com/images/equipment/infra/bluestar-hvac-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHI7', 'Eureka Forbes Air Purifier', 'Aeroguard SCPR', '01HFPQ2Z5P8VF5PXZRT4K7MHGC5', '01HFPQ2Z5P8VF5PXZRT4K7MHI7', 'Hospital air purification system',
'{"coverage": "1500 sq ft", "filtration": "6-stage with HEPA", "CADR": "500 m³/hr", "noise": "32-58 dB", "indicators": "AQI display", "warranty": "1 year"}',
85000.00, 'INR', 'EFK-HI-AGSCPR-001',
ARRAY['https://abymed.com/images/equipment/infra/eureka-aeroguard-1.jpg'],
TRUE, 'city-hospital'),

-- Medical Gas Pipeline
('01HFPQ3Z5P8VF5PXZRT4K7MHI8', 'Bharat Medical Gas System', 'Central Pipeline', '01HFPQ2Z5P8VF5PXZRT4K7MHGC6', '01HFPQ2Z5P8VF5PXZRT4K7MHI8', 'Medical gas pipeline system',
'{"gases": ["Oxygen", "Vacuum", "Medical Air", "Nitrous Oxide"], "capacity": "50 bed hospital", "outlets": "Imported quick-connect", "alarms": "Digital monitoring", "warranty": "5 years"}',
1250000.00, 'INR', 'BRT-HI-CMGP-001',
ARRAY['https://abymed.com/images/equipment/infra/bharat-gaspipe-1.jpg'],
TRUE, 'city-hospital'),

-- =============================================================================
-- 4. EMERGENCY & CRITICAL CARE (10 items) - CITY HOSPITAL TENANT
-- =============================================================================

-- Defibrillators
('01HFPQ3Z5P8VF5PXZRT4K7MHHL', 'Philips Defibrillator', 'HeartStart XL+', '01HFPQ2Z5P8VF5PXZRT4K7MHGD1', '01HFPQ2Z5P8VF5PXZRT4K7MHHL', 'Biphasic manual/AED defibrillator',
'{"energy": "0-200 J biphasic", "modes": ["Manual", "AED", "Sync cardioversion"], "monitoring": ["ECG", "SpO2", "NIBP"], "battery": "4 hrs", "warranty": "2 years"}',
285000.00, 'INR', 'PHP-EC-HSXL-001',
ARRAY['https://abymed.com/images/equipment/emergency/philips-xl-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHI9', 'BPL Defibrillator', 'Relife 700', '01HFPQ2Z5P8VF5PXZRT4K7MHGD1', '01HFPQ2Z5P8VF5PXZRT4K7MHHM', 'Indian-made defibrillator',
'{"energy": "0-360 J biphasic", "modes": ["Manual", "AED"], "monitoring": ["ECG"], "battery": "3 hrs", "display": "7 inch color", "warranty": "2 years"}',
175000.00, 'INR', 'BPL-EC-RL700-001',
ARRAY['https://abymed.com/images/equipment/emergency/bpl-relife-1.jpg'],
TRUE, 'city-hospital'),

-- Ventilators
('01HFPQ3Z5P8VF5PXZRT4K7MHIA', 'Drager Ventilator', 'Evita V500', '01HFPQ2Z5P8VF5PXZRT4K7MHGD2', '01HFPQ2Z5P8VF5PXZRT4K7MHIA', 'Advanced ICU ventilator',
'{"modes": ["VC-CMV", "PC-CMV", "CPAP", "PSV", "APRV"], "display": "15 inch touchscreen", "monitoring": ["Flow", "Pressure", "Volume", "EtCO2"], "battery": "30 min", "warranty": "2 years"}',
1250000.00, 'INR', 'DRG-EC-EV500-001',
ARRAY['https://abymed.com/images/equipment/emergency/drager-evita-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIB', 'AgVa Ventilator', 'Pro', '01HFPQ2Z5P8VF5PXZRT4K7MHGD2', '01HFPQ2Z5P8VF5PXZRT4K7MHIB', 'Indian-made compact ventilator',
'{"modes": ["VC-CMV", "PC-CMV", "CPAP", "PSV"], "display": "10 inch touchscreen", "monitoring": ["Flow", "Pressure", "Volume"], "battery": "4 hrs", "warranty": "1 year"}',
350000.00, 'INR', 'AGV-EC-PRO-001',
ARRAY['https://abymed.com/images/equipment/emergency/agva-pro-1.jpg'],
TRUE, 'city-hospital'),

-- Patient Monitors
('01HFPQ3Z5P8VF5PXZRT4K7MHIC', 'GE Patient Monitor', 'CARESCAPE B650', '01HFPQ2Z5P8VF5PXZRT4K7MHGD3', '01HFPQ2Z5P8VF5PXZRT4K7MHIC', 'Advanced patient monitor',
'{"parameters": ["ECG", "SpO2", "NIBP", "IBP", "Temp", "EtCO2"], "display": "15 inch color", "trends": "72 hours", "alarms": "3-level priority", "warranty": "2 years"}',
450000.00, 'INR', 'GEC-EC-B650-001',
ARRAY['https://abymed.com/images/equipment/emergency/ge-carescape-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHID', 'Skanray Patient Monitor', 'Star 65', '01HFPQ2Z5P8VF5PXZRT4K7MHGD3', '01HFPQ2Z5P8VF5PXZRT4K7MHGA', 'Indian-made patient monitor',
'{"parameters": ["ECG", "SpO2", "NIBP", "Temp"], "display": "12.1 inch color", "trends": "48 hours", "battery": "4 hrs", "warranty": "2 years"}',
175000.00, 'INR', 'SKN-EC-S65-001',
ARRAY['https://abymed.com/images/equipment/emergency/skanray-star65-1.jpg'],
TRUE, 'city-hospital'),

-- Infusion Pumps
('01HFPQ3Z5P8VF5PXZRT4K7MHIE', 'B Braun Infusion Pump', 'Infusomat Space', '01HFPQ2Z5P8VF5PXZRT4K7MHGD4', '01HFPQ2Z5P8VF5PXZRT4K7MHIE', 'Volumetric infusion pump',
'{"flow_rate": "0.1-1200 ml/h", "accuracy": "±2%", "display": "Color", "drug_library": "Yes", "alarms": ["Air", "Occlusion", "KVO"], "warranty": "2 years"}',
125000.00, 'INR', 'BBR-EC-IS-001',
ARRAY['https://abymed.com/images/equipment/emergency/bbraun-infusomat-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIF', 'Smiths Medical Syringe Pump', 'Graseby 2000', '01HFPQ2Z5P8VF5PXZRT4K7MHGD4', '01HFPQ2Z5P8VF5PXZRT4K7MHIF', 'Precision syringe pump',
'{"flow_rate": "0.1-999 ml/h", "accuracy": "±2%", "syringe_size": "5-60 ml", "alarms": ["Occlusion", "Near-empty", "End"], "battery": "10 hrs", "warranty": "2 years"}',
85000.00, 'INR', 'SMT-EC-G2000-001',
ARRAY['https://abymed.com/images/equipment/emergency/smiths-graseby-1.jpg'],
TRUE, 'city-hospital'),

-- Crash Carts & Emergency Equipment
('01HFPQ3Z5P8VF5PXZRT4K7MHIG', 'Pedigo Crash Cart', 'RC-2110', '01HFPQ2Z5P8VF5PXZRT4K7MHGD5', '01HFPQ2Z5P8VF5PXZRT4K7MHIG', 'Emergency crash cart',
'{"drawers": "6 color-coded", "accessories": ["Cardiac board", "Oxygen tank holder", "IV pole"], "material": "Powder-coated steel", "locking": "Breakaway seal", "warranty": "5 years"}',
95000.00, 'INR', 'PDG-EC-RC2110-001',
ARRAY['https://abymed.com/images/equipment/emergency/pedigo-cart-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIH', 'Allied Medical Transport Ventilator', 'MiniVent', '01HFPQ2Z5P8VF5PXZRT4K7MHGD6', '01HFPQ2Z5P8VF5PXZRT4K7MHIH', 'Emergency transport ventilator',
'{"modes": ["CMV", "SIMV", "CPAP"], "controls": "Simple dial", "weight": "3.5 kg", "battery": "6 hrs", "oxygen": "100% or Air-mix", "warranty": "2 years"}',
225000.00, 'INR', 'ALM-EC-MV-001',
ARRAY['https://abymed.com/images/equipment/emergency/allied-minivent-1.jpg'],
TRUE, 'city-hospital'),

-- =============================================================================
-- 5. SPECIALIZED EQUIPMENT (10 items) - CITY HOSPITAL TENANT
-- =============================================================================

-- Diagnostic Equipment
('01HFPQ3Z5P8VF5PXZRT4K7MHHM', 'BPL ECG Machine', 'Cardiart 9108D', '01HFPQ2Z5P8VF5PXZRT4K7MHGE1', '01HFPQ2Z5P8VF5PXZRT4K7MHHM', '12-channel digital ECG',
'{"channels": "12", "sampling": "1000 Hz", "display": "7 inch color", "memory": "200 records", "interfaces": ["USB", "LAN"], "battery": "3 hrs", "warranty": "2 years"}',
145000.00, 'INR', 'BPL-SE-9108D-001',
ARRAY['https://abymed.com/images/equipment/special/bpl-9108d-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHII', 'Schiller Stress Test System', 'CARDIOVIT CS-200', '01HFPQ2Z5P8VF5PXZRT4K7MHGE1', '01HFPQ2Z5P8VF5PXZRT4K7MHII', 'Cardiac stress test system',
'{"channels": "16", "sampling": "16,000 Hz", "display": "24 inch touchscreen", "protocols": ["Bruce", "Modified Bruce", "Naughton"], "interfaces": ["USB", "LAN", "DICOM"], "warranty": "2 years"}',
850000.00, 'INR', 'SCH-SE-CS200-001',
ARRAY['https://abymed.com/images/equipment/special/schiller-cs200-1.jpg'],
TRUE, 'city-hospital'),

-- Dialysis Systems
('01HFPQ3Z5P8VF5PXZRT4K7MHIJ', 'Fresenius Dialysis Machine', '4008S Classic', '01HFPQ2Z5P8VF5PXZRT4K7MHGE2', '01HFPQ2Z5P8VF5PXZRT4K7MHIJ', 'Hemodialysis system',
'{"treatment_modes": ["HD", "HF", "HDF"], "display": "10.4 inch TFT", "blood_flow": "30-600 ml/min", "dialysate_flow": "300-800 ml/min", "disinfection": "Heat and chemical", "warranty": "1 year"}',
950000.00, 'INR', 'FRS-SE-4008S-001',
ARRAY['https://abymed.com/images/equipment/special/fresenius-4008s-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIK', 'Nipro Dialysis Machine', 'Surdial X', '01HFPQ2Z5P8VF5PXZRT4K7MHGE2', '01HFPQ2Z5P8VF5PXZRT4K7MHIK', 'Compact hemodialysis system',
'{"treatment_modes": ["HD", "HF"], "display": "10.1 inch LCD", "blood_flow": "30-500 ml/min", "dialysate_flow": "300-700 ml/min", "disinfection": "Heat", "warranty": "1 year"}',
750000.00, 'INR', 'NPR-SE-SDX-001',
ARRAY['https://abymed.com/images/equipment/special/nipro-surdial-1.jpg'],
TRUE, 'city-hospital'),

-- Endoscopy Systems
('01HFPQ3Z5P8VF5PXZRT4K7MHIL', 'Olympus Endoscopy System', 'EVIS X1', '01HFPQ2Z5P8VF5PXZRT4K7MHGE3', '01HFPQ2Z5P8VF5PXZRT4K7MHGU', 'Advanced endoscopy platform',
'{"components": ["Processor", "Light source", "Monitor", "Endoscopes"], "imaging": ["NBI", "TXI", "RDI"], "display": "31 inch 4K", "recording": "Full HD", "warranty": "1 year"}',
3500000.00, 'INR', 'OLY-SE-EVISX1-001',
ARRAY['https://abymed.com/images/equipment/special/olympus-evisx1-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIM', 'Pentax Medical Endoscopy System', 'EPK-i7010', '01HFPQ2Z5P8VF5PXZRT4K7MHGE3', '01HFPQ2Z5P8VF5PXZRT4K7MHIM', 'High-definition endoscopy system',
'{"components": ["Processor", "Light source", "Monitor"], "imaging": ["i-scan", "HD+"], "display": "27 inch Full HD", "recording": "HD", "warranty": "1 year"}',
2500000.00, 'INR', 'PTX-SE-EPK7010-001',
ARRAY['https://abymed.com/images/equipment/special/pentax-epk-1.jpg'],
TRUE, 'city-hospital'),

-- Physiotherapy Equipment
('01HFPQ3Z5P8VF5PXZRT4K7MHIN', 'BTL Physiotherapy System', 'BTL-4000 Premium', '01HFPQ2Z5P8VF5PXZRT4K7MHGE4', '01HFPQ2Z5P8VF5PXZRT4K7MHIN', 'Multi-modality therapy system',
'{"therapies": ["Electrotherapy", "Ultrasound", "Laser", "Magnetotherapy"], "channels": "2 independent", "display": "8.4 inch color touchscreen", "protocols": "50+ preset", "warranty": "2 years"}',
375000.00, 'INR', 'BTL-SE-4000P-001',
ARRAY['https://abymed.com/images/equipment/special/btl-4000-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIO', 'Chattanooga Physiotherapy', 'Intelect Advanced', '01HFPQ2Z5P8VF5PXZRT4K7MHGE4', '01HFPQ2Z5P8VF5PXZRT4K7MHIO', 'Combo therapy unit',
'{"therapies": ["Electrotherapy", "Ultrasound"], "channels": "4 stim", "display": "10 inch color touchscreen", "protocols": "200+ clinical", "warranty": "2 years"}',
225000.00, 'INR', 'CHT-SE-IA-001',
ARRAY['https://abymed.com/images/equipment/special/chattanooga-intelect-1.jpg'],
TRUE, 'city-hospital'),

-- Ophthalmic Equipment
('01HFPQ3Z5P8VF5PXZRT4K7MHIP', 'Carl Zeiss Slit Lamp', 'SL 220', '01HFPQ2Z5P8VF5PXZRT4K7MHGE5', '01HFPQ2Z5P8VF5PXZRT4K7MHGW', 'Advanced slit lamp biomicroscope',
'{"magnification": "6x to 40x", "slit_width": "0-14 mm", "slit_length": "1-14 mm", "filters": ["Cobalt blue", "Red-free", "Grey"], "illumination": "LED", "warranty": "1 year"}',
650000.00, 'INR', 'CZS-SE-SL220-001',
ARRAY['https://abymed.com/images/equipment/special/zeiss-sl220-1.jpg'],
TRUE, 'city-hospital'),

('01HFPQ3Z5P8VF5PXZRT4K7MHIQ', 'Appasamy Auto Refractometer', 'ARK-900', '01HFPQ2Z5P8VF5PXZRT4K7MHGE5', '01HFPQ2Z5P8VF5PXZRT4K7MHIQ', 'Indian-made auto refractometer',
'{"measurement_range": "Sphere: -25D to +22D, Cylinder: 0D to ±10D", "minimum_pupil": "2.0 mm", "display": "7 inch LCD touchscreen", "printer": "Built-in thermal", "warranty": "1 year"}',
450000.00, 'INR', 'APS-SE-ARK900-001',
ARRAY['https://abymed.com/images/equipment/special/appasamy-ark900-1.jpg'],
TRUE, 'city-hospital');
