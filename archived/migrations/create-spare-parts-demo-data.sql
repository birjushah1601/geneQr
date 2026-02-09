-- Create comprehensive spare parts catalog for medical equipment

-- X-Ray Machine Parts
INSERT INTO spare_parts_catalog (part_number, part_name, manufacturer_part_number, category, subcategory, part_type, description, unit_price, currency, is_available, stock_status, lead_time_days, requires_engineer, engineer_level_required, installation_time_minutes) VALUES
('XR-TUBE-001', 'X-Ray Tube Assembly', 'VARIAN-G1592', 'X-Ray', 'Core Component', 'Critical', 'High-voltage X-ray tube assembly, 150kV capacity', 12500.00, 'USD', true, 'in_stock', 7, true, 'L3', 180),
('XR-DET-001', 'Flat Panel Detector', 'CANON-CXDI-DET', 'X-Ray', 'Imaging', 'Critical', 'Digital flat panel detector 17x17 inch', 28000.00, 'USD', true, 'in_stock', 14, true, 'L3', 120),
('XR-COL-001', 'Collimator Assembly', 'CANON-COL-450', 'X-Ray', 'Beam Control', 'Important', 'Automatic collimator with light field', 3500.00, 'USD', true, 'in_stock', 5, true, 'L2', 90),
('XR-FILT-001', 'X-Ray Filter Set', 'GEN-FILT-AL', 'X-Ray', 'Beam Control', 'Consumable', 'Aluminum filtration set 0.5-3mm', 450.00, 'USD', true, 'in_stock', 2, false, 'L1', 15),
('XR-GRID-001', 'Anti-Scatter Grid', 'GRID-103L', 'X-Ray', 'Imaging', 'Important', 'Focused anti-scatter grid 12:1 ratio', 1200.00, 'USD', true, 'in_stock', 3, false, 'L2', 30),

-- CT Scanner Parts
('CT-TUBE-001', 'CT X-Ray Tube', 'SIEMENS-STRATON', 'CT', 'Core Component', 'Critical', 'Straton ceramic tube for CT scanner', 85000.00, 'USD', true, 'in_stock', 21, true, 'L3', 240),
('CT-DET-001', 'CT Detector Module', 'SIE-UFC-DET', 'CT', 'Imaging', 'Critical', 'Ultra-fast ceramic detector module', 42000.00, 'USD', true, 'low_stock', 30, true, 'L3', 180),
('CT-SLIP-001', 'Slip Ring Assembly', 'SLIP-RING-CT', 'CT', 'Power Transfer', 'Critical', 'High-speed slip ring for continuous rotation', 18000.00, 'USD', true, 'in_stock', 14, true, 'L3', 300),
('CT-COL-001', 'CT Collimator', 'SIE-COL-128', 'CT', 'Beam Control', 'Important', '128-slice collimator assembly', 8500.00, 'USD', true, 'in_stock', 10, true, 'L3', 120),

-- MRI Scanner Parts  
('MRI-COIL-HEAD', 'MRI Head Coil', 'SIE-HC-32CH', 'MRI', 'RF Coil', 'Important', '32-channel head coil for 3T MRI', 15000.00, 'USD', true, 'in_stock', 14, false, 'L2', 30),
('MRI-COIL-BODY', 'MRI Body Coil', 'SIE-BC-18CH', 'MRI', 'RF Coil', 'Important', '18-channel body coil for 3T MRI', 22000.00, 'USD', true, 'in_stock', 14, false, 'L2', 45),
('MRI-GRAD-001', 'Gradient Coil', 'SIE-GRAD-120', 'MRI', 'Gradient System', 'Critical', 'High-performance gradient coil 80 mT/m', 125000.00, 'USD', true, 'low_stock', 60, true, 'L3', 480),
('MRI-CRYO-001', 'Cryogen System', 'CRYO-HE-500L', 'MRI', 'Cooling', 'Critical', 'Liquid helium cryogen system 500L', 8500.00, 'USD', true, 'in_stock', 3, true, 'L3', 240),
('MRI-RF-AMP', 'RF Power Amplifier', 'RF-AMP-35KW', 'MRI', 'RF System', 'Critical', '35kW RF power amplifier for 3T', 45000.00, 'USD', true, 'in_stock', 30, true, 'L3', 180),

-- Ultrasound Parts
('US-PROBE-C60', 'Convex Ultrasound Probe', 'GE-C1-6D', 'Ultrasound', 'Transducer', 'Critical', 'Convex probe 1-6 MHz for abdominal imaging', 8500.00, 'USD', true, 'in_stock', 7, false, 'L1', 10),
('US-PROBE-L38', 'Linear Ultrasound Probe', 'GE-L3-12', 'Ultrasound', 'Transducer', 'Critical', 'Linear probe 3-12 MHz for vascular', 9200.00, 'USD', true, 'in_stock', 7, false, 'L1', 10),
('US-GEL-001', 'Ultrasound Gel 5L', 'AQUASONIC-5L', 'Ultrasound', 'Consumable', 'Consumable', 'Medical ultrasound transmission gel', 45.00, 'USD', true, 'in_stock', 1, false, 'L1', 2),
('US-BATT-001', 'Ultrasound Battery', 'BAT-LI-14V', 'Ultrasound', 'Power', 'Important', 'Lithium-ion battery pack 14.8V 6800mAh', 450.00, 'USD', true, 'in_stock', 3, false, 'L1', 15),

-- Ventilator Parts
('VENT-VALVE-001', 'Expiratory Valve', 'DRAGER-EXP-V', 'Ventilator', 'Breathing Circuit', 'Critical', 'Expiratory valve assembly with sensor', 850.00, 'USD', true, 'in_stock', 2, false, 'L2', 20),
('VENT-VALVE-002', 'Inspiratory Valve', 'DRAGER-INS-V', 'Ventilator', 'Breathing Circuit', 'Critical', 'Inspiratory valve assembly with flow sensor', 920.00, 'USD', true, 'in_stock', 2, false, 'L2', 20),
('VENT-SENS-O2', 'Oxygen Sensor', 'OXY-SENS-PO2', 'Ventilator', 'Monitoring', 'Critical', 'Galvanic oxygen sensor 0-100%', 280.00, 'USD', true, 'in_stock', 1, false, 'L1', 10),
('VENT-SENS-CO2', 'CO2 Sensor', 'CAPNO-SENS-IR', 'Ventilator', 'Monitoring', 'Critical', 'Infrared CO2 sensor module', 1200.00, 'USD', true, 'in_stock', 3, false, 'L2', 30),
('VENT-FILT-001', 'HEPA Filter', 'HEPA-H14-MED', 'Ventilator', 'Filtration', 'Consumable', 'H14 HEPA filter for breathing circuit', 85.00, 'USD', true, 'in_stock', 1, false, 'L1', 5),
('VENT-TUBE-001', 'Breathing Circuit Tubing', 'TUBE-22MM-150', 'Ventilator', 'Breathing Circuit', 'Consumable', '22mm disposable breathing circuit', 35.00, 'USD', true, 'in_stock', 1, false, 'L1', 5),
('VENT-BATT-001', 'Ventilator Battery', 'BAT-VENT-12V', 'Ventilator', 'Power', 'Important', 'Rechargeable battery 12V 7.2Ah', 380.00, 'USD', true, 'in_stock', 2, false, 'L1', 15),

-- Patient Monitor Parts
('PM-ECG-CABLE', 'ECG Cable 5-Lead', 'PM-ECG-5L-AHA', 'Patient Monitor', 'ECG', 'Important', '5-lead ECG cable AHA standard', 145.00, 'USD', true, 'in_stock', 2, false, 'L1', 5),
('PM-SPO2-SENSOR', 'SpO2 Sensor Adult', 'SPO2-ADULT-CLIP', 'Patient Monitor', 'SpO2', 'Important', 'Reusable adult SpO2 clip sensor', 95.00, 'USD', true, 'in_stock', 1, false, 'L1', 2),
('PM-NIBP-CUFF', 'NIBP Cuff Adult', 'CUFF-NIBP-ADULT', 'Patient Monitor', 'NIBP', 'Consumable', 'Disposable adult NIBP cuff 25-35cm', 25.00, 'USD', true, 'in_stock', 1, false, 'L1', 2),
('PM-TEMP-PROBE', 'Temperature Probe', 'TEMP-PROBE-ORAL', 'Patient Monitor', 'Temperature', 'Consumable', 'Disposable oral temperature probe', 8.00, 'USD', true, 'in_stock', 1, false, 'L1', 1),
('PM-IBP-CABLE', 'IBP Cable', 'IBP-CABLE-DUAL', 'Patient Monitor', 'IBP', 'Important', 'Invasive blood pressure cable dual channel', 180.00, 'USD', true, 'in_stock', 3, false, 'L1', 5),
('PM-BATT-001', 'Monitor Battery', 'BAT-PM-14V-5AH', 'Patient Monitor', 'Power', 'Important', 'Li-ion battery 14.4V 5.0Ah', 420.00, 'USD', true, 'in_stock', 2, false, 'L1', 15),
('PM-DISPLAY-001', 'LCD Display Module', 'LCD-15-TOUCH', 'Patient Monitor', 'Display', 'Critical', '15" touchscreen LCD module', 2800.00, 'USD', true, 'in_stock', 14, true, 'L2', 90),

-- Dialysis Machine Parts
('DIAL-PUMP-001', 'Blood Pump Head', 'FRES-BP-HEAD', 'Dialysis', 'Blood Circuit', 'Critical', 'Peristaltic blood pump head', 320.00, 'USD', true, 'in_stock', 3, false, 'L2', 20),
('DIAL-FILT-001', 'Dialyzer Filter High-Flux', 'FX80-HEMOFILTER', 'Dialysis', 'Filtration', 'Consumable', 'High-flux polysulfone dialyzer 1.8m²', 65.00, 'USD', true, 'in_stock', 1, false, 'L1', 5),
('DIAL-LINE-001', 'Bloodline Set', 'BLOODLINE-STERILE', 'Dialysis', 'Blood Circuit', 'Consumable', 'Sterile disposable bloodline set', 28.00, 'USD', true, 'in_stock', 1, false, 'L1', 10),
('DIAL-CONC-BIC', 'Bicarbonate Concentrate', 'BICARB-CONC-5L', 'Dialysis', 'Dialysate', 'Consumable', 'Bicarbonate concentrate 5L container', 45.00, 'USD', true, 'in_stock', 1, false, 'L1', 5),
('DIAL-CONC-ACID', 'Acid Concentrate', 'ACID-CONC-5L', 'Dialysis', 'Dialysate', 'Consumable', 'Acid concentrate 5L container', 38.00, 'USD', true, 'in_stock', 1, false, 'L1', 5),
('DIAL-PRES-001', 'Pressure Transducer', 'PRESS-TRANS-300', 'Dialysis', 'Monitoring', 'Important', 'Blood pressure transducer 0-300mmHg', 180.00, 'USD', true, 'in_stock', 2, false, 'L1', 15),
('DIAL-VALVE-001', 'Solenoid Valve', 'VALVE-SOL-24V', 'Dialysis', 'Fluid Control', 'Important', '24V solenoid valve for dialysate', 145.00, 'USD', true, 'in_stock', 3, false, 'L2', 30),

-- Anesthesia Machine Parts
('ANES-VAPOR-ISO', 'Isoflurane Vaporizer', 'VAPOR-ISO-TEC7', 'Anesthesia', 'Vaporizer', 'Critical', 'Temperature-compensated isoflurane vaporizer', 8500.00, 'USD', true, 'in_stock', 14, true, 'L3', 60),
('ANES-VAPOR-SEV', 'Sevoflurane Vaporizer', 'VAPOR-SEV-TEC7', 'Anesthesia', 'Vaporizer', 'Critical', 'Temperature-compensated sevoflurane vaporizer', 8800.00, 'USD', true, 'in_stock', 14, true, 'L3', 60),
('ANES-CO2-ABS', 'CO2 Absorbent Canister', 'SODASORB-1.5KG', 'Anesthesia', 'Breathing Circuit', 'Consumable', 'Soda lime CO2 absorbent 1.5kg', 55.00, 'USD', true, 'in_stock', 1, false, 'L1', 10),
('ANES-O2-SENS', 'Oxygen Analyzer Sensor', 'O2-CELL-PSR11', 'Anesthesia', 'Monitoring', 'Important', 'Paramagnetic oxygen sensor', 480.00, 'USD', true, 'in_stock', 5, false, 'L2', 20),
('ANES-BELLOW', 'Ventilator Bellows', 'BELLOW-ANES-2L', 'Anesthesia', 'Ventilation', 'Important', 'Ascending bellows assembly 2L', 650.00, 'USD', true, 'in_stock', 7, false, 'L2', 30);

-- Display results
SELECT COUNT(*) as "Total Spare Parts" FROM spare_parts_catalog;

SELECT category, COUNT(*) as count 
FROM spare_parts_catalog 
GROUP BY category 
ORDER BY count DESC;
