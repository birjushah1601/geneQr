-- Sample Equipment Catalog for ServQR Platform
-- This file populates the equipment table with representative medical devices from Indian manufacturers

-- First ensure the schema exists
CREATE SCHEMA IF NOT EXISTS public;

-- Create the equipment table if it doesn't exist
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

-- Clear existing data if needed (commented out for safety)
-- TRUNCATE TABLE equipment;

-- Insert X-ray machines from BPL Medical and Allengers
INSERT INTO equipment (id, name, model, category_id, manufacturer_id, description, specifications, price_amount, price_currency, sku, images, is_active, tenant_id) VALUES
(
    '01HFPQZ6QBWK7NQZXT5G7JHA1',
    'BPL Digital X-Ray System',
    'DX-5100',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ14', -- X-ray Systems category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'Advanced digital X-ray system with high-resolution imaging and reduced radiation exposure. Suitable for general radiography in hospitals and diagnostic centers.',
    '{
        "detector_type": "Flat Panel Digital Detector",
        "detector_size": "43cm x 43cm (17\" x 17\")",
        "resolution": "3.5 lp/mm",
        "pixel_size": "140 Î¼m",
        "image_depth": "16-bit",
        "generator_power": "50 kW",
        "tube_voltage_range": "40-150 kV",
        "tube_current_range": "10-630 mA",
        "exposure_time": "1 ms to 5 sec",
        "weight_capacity": "200 kg",
        "dimensions": "2200 x 1800 x 2300 mm",
        "power_requirements": "380-400V, 50/60 Hz, three phase",
        "dicom_compatible": true,
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "AERB"]
    }',
    1850000.00,
    'INR',
    'BPL-XR-DX5100',
    ARRAY['https://ServQR.com/images/equipment/bpl-dx5100-1.jpg', 'https://ServQR.com/images/equipment/bpl-dx5100-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA2',
    'BPL Mobile X-Ray Unit',
    'MX-2020',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ14', -- X-ray Systems category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'Compact and portable X-ray unit designed for bedside radiography in ICUs, emergency rooms, and operating theaters. Features lightweight design and easy maneuverability.',
    '{
        "generator_type": "High Frequency",
        "power_output": "20 kW",
        "tube_voltage_range": "40-125 kV",
        "tube_current_range": "10-250 mA",
        "exposure_time": "0.001-5 seconds",
        "battery_capacity": "Up to 300 exposures on full charge",
        "charging_time": "3-4 hours",
        "collimator": "Manual with LED light",
        "arm_reach": "1200 mm",
        "vertical_travel": "500-2000 mm",
        "rotation": "Â±180Â°",
        "weight": "175 kg",
        "dimensions": "1200 x 700 x 1750 mm",
        "display": "10.4\" touchscreen",
        "warranty": "2 years",
        "certifications": ["CE", "ISO 13485", "AERB"]
    }',
    950000.00,
    'INR',
    'BPL-XR-MX2020',
    ARRAY['https://ServQR.com/images/equipment/bpl-mx2020-1.jpg', 'https://ServQR.com/images/equipment/bpl-mx2020-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA3',
    'Allengers Digital Radiography System',
    'MARS-40',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ14', -- X-ray Systems category
    '01HFPQ2Z5P8VF5PXZRT4K7MHG', -- Allengers Medical Systems
    'Fully integrated digital radiography system with advanced image processing capabilities. Provides exceptional image quality with minimal radiation dose.',
    '{
        "detector_type": "Amorphous Silicon Flat Panel Detector",
        "detector_size": "43cm x 43cm (17\" x 17\")",
        "resolution": "3.9 lp/mm",
        "pixel_size": "127 Î¼m",
        "image_depth": "16-bit",
        "generator_power": "65 kW",
        "tube_voltage_range": "40-150 kV",
        "tube_current_range": "10-800 mA",
        "exposure_time": "0.001-10 seconds",
        "weight_capacity": "250 kg",
        "table_movement": "Floating top with electromagnetic locks",
        "vertical_travel": "550-900 mm",
        "dimensions": "2300 x 1900 x 2400 mm",
        "power_requirements": "415V, 50/60 Hz, three phase",
        "image_processing": "Advanced MARS processing algorithm",
        "storage_capacity": "10,000 images (local)",
        "dicom_compatible": true,
        "warranty": "5 years",
        "certifications": ["CE", "ISO 13485", "FDA", "AERB"]
    }',
    2250000.00,
    'INR',
    'ALG-DR-MARS40',
    ARRAY['https://ServQR.com/images/equipment/allengers-mars40-1.jpg', 'https://ServQR.com/images/equipment/allengers-mars40-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA4',
    'Allengers C-Arm System',
    'FLEXIVIEW-3D',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ14', -- X-ray Systems category
    '01HFPQ2Z5P8VF5PXZRT4K7MHG', -- Allengers Medical Systems
    'Advanced mobile C-arm imaging system for real-time visualization during surgical, orthopedic, and vascular procedures. Features 3D reconstruction capabilities.',
    '{
        "detector_type": "CMOS Flat Panel",
        "detector_size": "30cm x 30cm (12\" x 12\")",
        "resolution": "4.0 lp/mm",
        "pixel_size": "125 Î¼m",
        "image_depth": "16-bit",
        "generator_power": "15 kW",
        "tube_voltage_range": "40-120 kV",
        "tube_current_range": "0.2-250 mA",
        "exposure_modes": ["Continuous", "Pulsed", "Single Shot"],
        "orbital_rotation": "135Â°",
        "horizontal_travel": "200 mm",
        "vertical_travel": "450 mm",
        "swivel_range": "Â±225Â°",
        "display": "Dual 19\" medical-grade LCD monitors",
        "cooling_system": "Liquid cooling with integrated heat exchanger",
        "3D_reconstruction": true,
        "storage": "1TB SSD",
        "dimensions": "1950 x 850 x 1750 mm",
        "weight": "320 kg",
        "battery_backup": "30 minutes",
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "AERB"]
    }',
    3500000.00,
    'INR',
    'ALG-CA-FLEX3D',
    ARRAY['https://ServQR.com/images/equipment/allengers-flexiview-1.jpg', 'https://ServQR.com/images/equipment/allengers-flexiview-2.jpg'],
    TRUE,
    'demo-hospital'
),

-- Insert Patient Monitors from Trivitron and BPL
(
    '01HFPQZ6QBWK7NQZXT5G7JHA5',
    'Trivitron Elite Patient Monitor',
    'Clarity-Pro 12',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ40', -- Multi-Parameter Monitors category
    '01HFPQ2Z5MXNVT9DAPZ3BWJAHR', -- Trivitron Healthcare
    'Advanced multi-parameter patient monitor with comprehensive monitoring capabilities for critical care, emergency, and general ward settings.',
    '{
        "display": "12.1\" TFT color touchscreen",
        "resolution": "1280 x 800 pixels",
        "parameters": ["ECG", "NIBP", "SpO2", "Respiration", "Temperature", "IBP", "EtCO2", "Anesthetic Agents", "BIS"],
        "ecg_leads": "3/5/12-lead selectable",
        "ecg_analysis": "ST segment analysis, arrhythmia detection",
        "nibp_measurement": "Oscillometric method with auto/manual/continuous modes",
        "spo2_technology": "Nellcor compatible",
        "temperature_channels": 2,
        "ibp_channels": 2,
        "etco2_method": "Sidestream/Mainstream selectable",
        "trends": "120 hours graphical and tabular",
        "alarms": "3-level visual and audible",
        "connectivity": ["Wi-Fi", "Ethernet", "HL7", "DICOM"],
        "battery": "Lithium-ion, up to 4 hours operation",
        "weight": "4.5 kg",
        "dimensions": "318 x 264 x 152 mm",
        "storage": "32GB internal, expandable via USB",
        "printer": "Built-in thermal printer",
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    285000.00,
    'INR',
    'TRV-PM-CP12',
    ARRAY['https://ServQR.com/images/equipment/trivitron-claritypro-1.jpg', 'https://ServQR.com/images/equipment/trivitron-claritypro-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA6',
    'Trivitron Transport Monitor',
    'Vitals-T7',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ40', -- Multi-Parameter Monitors category
    '01HFPQ2Z5MXNVT9DAPZ3BWJAHR', -- Trivitron Healthcare
    'Compact and lightweight patient monitor designed for intra-hospital transport and emergency care. Features rugged design and long battery life.',
    '{
        "display": "7\" TFT color touchscreen",
        "resolution": "800 x 480 pixels",
        "parameters": ["ECG", "NIBP", "SpO2", "Respiration", "Temperature"],
        "ecg_leads": "3/5-lead selectable",
        "ecg_analysis": "Basic arrhythmia detection",
        "nibp_measurement": "Oscillometric method with auto/manual modes",
        "spo2_technology": "Nellcor compatible",
        "temperature_channels": 1,
        "trends": "72 hours graphical and tabular",
        "alarms": "3-level visual and audible",
        "connectivity": ["Wi-Fi", "Bluetooth"],
        "battery": "Lithium-ion, up to 8 hours operation",
        "weight": "1.2 kg",
        "dimensions": "200 x 140 x 80 mm",
        "storage": "16GB internal",
        "drop_resistance": "Withstands drops from 1 meter height",
        "water_resistance": "IPX4 rated",
        "warranty": "2 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    125000.00,
    'INR',
    'TRV-PM-VT7',
    ARRAY['https://ServQR.com/images/equipment/trivitron-vitalst7-1.jpg', 'https://ServQR.com/images/equipment/trivitron-vitalst7-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA7',
    'BPL Intensive Care Monitor',
    'Penlon ICU-15',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ40', -- Multi-Parameter Monitors category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'High-end patient monitoring system designed for intensive care units with advanced parameter monitoring and analysis capabilities.',
    '{
        "display": "15\" medical-grade color touchscreen",
        "resolution": "1366 x 768 pixels",
        "parameters": ["ECG", "NIBP", "SpO2", "Respiration", "Temperature", "IBP", "EtCO2", "Cardiac Output", "BIS", "NMT"],
        "ecg_leads": "3/5/12-lead selectable",
        "ecg_analysis": "Advanced arrhythmia detection, ST/QT analysis, pacemaker detection",
        "nibp_measurement": "Oscillometric method with auto/manual/continuous modes",
        "spo2_technology": "Masimo SETÂ®",
        "temperature_channels": 2,
        "ibp_channels": 4,
        "etco2_method": "Mainstream and Sidestream",
        "cardiac_output": "Thermodilution method",
        "trends": "168 hours graphical and tabular",
        "alarms": "3-level visual and audible with smart alarm management",
        "connectivity": ["Wi-Fi", "Ethernet", "HL7", "DICOM"],
        "central_monitoring": true,
        "battery": "Lithium-ion, hot-swappable, up to 5 hours operation",
        "weight": "5.8 kg",
        "dimensions": "370 x 320 x 180 mm",
        "storage": "64GB SSD",
        "printer": "Built-in thermal printer",
        "warranty": "5 years",
        "certifications": ["CE", "ISO 13485", "BIS", "FDA"]
    }',
    450000.00,
    'INR',
    'BPL-PM-ICU15',
    ARRAY['https://ServQR.com/images/equipment/bpl-penlon-1.jpg', 'https://ServQR.com/images/equipment/bpl-penlon-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JHA8',
    'BPL Neonatal Monitor',
    'NeoGuard Plus',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ40', -- Multi-Parameter Monitors category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'Specialized patient monitor designed specifically for neonatal and pediatric care with gentle monitoring capabilities and specialized algorithms.',
    '{
        "display": "10.4\" TFT color touchscreen",
        "resolution": "1024 x 768 pixels",
        "parameters": ["ECG", "NIBP", "SpO2", "Respiration", "Temperature", "IBP", "EtCO2"],
        "ecg_leads": "3/5-lead selectable",
        "ecg_analysis": "Neonatal-specific arrhythmia detection",
        "nibp_measurement": "Oscillometric method with neonatal cuffs",
        "spo2_technology": "Masimo SETÂ® with neonatal sensors",
        "temperature_channels": 2,
        "ibp_channels": 2,
        "etco2_method": "Microstream technology (low flow)",
        "trends": "120 hours graphical and tabular",
        "alarms": "3-level visual and audible with gentle sound options",
        "connectivity": ["Wi-Fi", "Ethernet", "HL7"],
        "central_monitoring": true,
        "battery": "Lithium-ion, up to 6 hours operation",
        "weight": "3.2 kg",
        "dimensions": "280 x 230 x 150 mm",
        "storage": "32GB internal",
        "printer": "Built-in thermal printer",
        "special_features": ["Apnea detection", "Gentle alarm sounds", "Night mode"],
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    325000.00,
    'INR',
    'BPL-PM-NEOPLUS',
    ARRAY['https://ServQR.com/images/equipment/bpl-neoguard-1.jpg', 'https://ServQR.com/images/equipment/bpl-neoguard-2.jpg'],
    TRUE,
    'demo-hospital'
),

-- Insert Surgical Instruments from Hindustan Syringes and Poly Medicure
(
    '01HFPQZ6QBWK7NQZXT5G7JHA9',
    'Hindustan Syringes Auto-Disable Syringe',
    'Kojak AD',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ24', -- General Surgical Instruments category
    '01HFPQ2Z5RBWKPN8QZXT5G7JHF', -- Hindustan Syringes & Medical Devices
    'Auto-disable syringes that automatically lock after single use, preventing reuse and reducing the risk of cross-contamination and infections.',
    '{
        "sizes": ["0.5ml", "1ml", "2ml", "3ml", "5ml", "10ml"],
        "needle_sizes": ["23G", "24G", "25G", "26G"],
        "material": "Medical-grade polypropylene",
        "sterilization": "ETO sterilized",
        "shelf_life": "5 years",
        "locking_mechanism": "Automatic plunger lock after injection",
        "packaging": "Individual blister packs",
        "box_quantity": 100,
        "case_quantity": 2000,
        "latex_free": true,
        "phthalate_free": true,
        "color_coded": true,
        "certifications": ["ISO 13485", "CE", "WHO-PQS"]
    }',
    1200.00, -- Price per box of 100
    'INR',
    'HMD-SYR-KOJAK-2ML',
    ARRAY['https://ServQR.com/images/equipment/hmd-kojak-1.jpg', 'https://ServQR.com/images/equipment/hmd-kojak-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH10',
    'Hindustan Syringes Safety IV Cannula',
    'SafetyPro',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ24', -- General Surgical Instruments category
    '01HFPQ2Z5RBWKPN8QZXT5G7JHF', -- Hindustan Syringes & Medical Devices
    'Advanced safety IV cannula with needle protection mechanism to prevent needlestick injuries. Features smooth insertion and secure fixation.',
    '{
        "sizes": ["14G", "16G", "18G", "20G", "22G", "24G", "26G"],
        "color_coding": {
            "14G": "Orange",
            "16G": "Grey",
            "18G": "Green",
            "20G": "Pink",
            "22G": "Blue",
            "24G": "Yellow",
            "26G": "Purple"
        },
        "material": {
            "catheter": "FEP (Fluorinated Ethylene Propylene)",
            "needle": "Stainless steel with triple-facet bevel"
        },
        "safety_mechanism": "Automatic needle retraction after use",
        "flow_rate": {
            "14G": "270-330 ml/min",
            "16G": "180-210 ml/min",
            "18G": "90-110 ml/min",
            "20G": "60-80 ml/min",
            "22G": "36-45 ml/min",
            "24G": "22-28 ml/min",
            "26G": "13-17 ml/min"
        },
        "wings": "Flexible fixation wings",
        "injection_port": "Resealable",
        "sterilization": "ETO sterilized",
        "shelf_life": "5 years",
        "latex_free": true,
        "dehp_free": true,
        "packaging": "Individual blister packs",
        "box_quantity": 50,
        "case_quantity": 500,
        "certifications": ["ISO 13485", "CE", "FDA"]
    }',
    3500.00, -- Price per box of 50
    'INR',
    'HMD-CAN-SAFETY-20G',
    ARRAY['https://ServQR.com/images/equipment/hmd-safetypro-1.jpg', 'https://ServQR.com/images/equipment/hmd-safetypro-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH11',
    'Poly Medicure Infusion Set',
    'PolyFlo-IV',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ24', -- General Surgical Instruments category
    '01HFPQ2Z5RWTHGV9PZBS3F8KMN', -- Poly Medicure Limited
    'Precision infusion set for accurate and controlled delivery of fluids and medications. Features kink-resistant tubing and secure connections.',
    '{
        "tubing_length": "180 cm",
        "tubing_material": "Medical-grade PVC",
        "flow_regulator": "Precision roller clamp",
        "drip_chamber": "Transparent with filter",
        "drop_rate": "20 drops/ml",
        "connector": "Luer lock",
        "injection_port": "Y-type latex-free",
        "air_vent": "Hydrophobic filter",
        "sterilization": "ETO sterilized",
        "shelf_life": "5 years",
        "latex_free": true,
        "dehp_free": true,
        "packaging": "Individual peel-open packs",
        "box_quantity": 100,
        "case_quantity": 1000,
        "certifications": ["ISO 13485", "CE", "FDA"]
    }',
    2800.00, -- Price per box of 100
    'INR',
    'POLY-INF-PF-STD',
    ARRAY['https://ServQR.com/images/equipment/poly-polyflo-1.jpg', 'https://ServQR.com/images/equipment/poly-polyflo-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH12',
    'Poly Medicure Hemodialysis Blood Line Set',
    'PolyDial-Pro',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ24', -- General Surgical Instruments category
    '01HFPQ2Z5RWTHGV9PZBS3F8KMN', -- Poly Medicure Limited
    'Complete blood line set for hemodialysis procedures with universal compatibility with most dialysis machines. Features secure connections and optimal blood flow.',
    '{
        "compatibility": ["Fresenius", "Nipro", "B. Braun", "Nikkiso", "Gambro"],
        "tubing_material": "Medical-grade PVC",
        "tubing_length": {
            "arterial": "270 cm",
            "venous": "270 cm"
        },
        "connectors": "Universal luer lock",
        "blood_pump_segment": "Silicon rubber, 8mm internal diameter",
        "pressure_monitoring_lines": 4,
        "injection_ports": 6,
        "clamps": {
            "arterial": 2,
            "venous": 2
        },
        "drip_chambers": {
            "arterial": "With filter",
            "venous": "With filter and level adjustment"
        },
        "transducer_protectors": 4,
        "sterilization": "ETO sterilized",
        "shelf_life": "3 years",
        "latex_free": true,
        "dehp_free": true,
        "packaging": "Individual sterile packs",
        "box_quantity": 20,
        "case_quantity": 100,
        "certifications": ["ISO 13485", "CE", "FDA"]
    }',
    12000.00, -- Price per box of 20
    'INR',
    'POLY-HD-PRODIAL',
    ARRAY['https://ServQR.com/images/equipment/poly-dialysis-1.jpg', 'https://ServQR.com/images/equipment/poly-dialysis-2.jpg'],
    TRUE,
    'demo-hospital'
),

-- Insert Lab Analyzers from Transasia Bio-Medicals
(
    '01HFPQZ6QBWK7NQZXT5G7JH13',
    'Transasia Biochemistry Analyzer',
    'Erba XL-640',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ20', -- Biochemistry Analyzers category
    '01HFPQ2Z5N7YWGXS9JVKM6F8QT', -- Transasia Bio-Medicals
    'Fully automated clinical chemistry analyzer with high throughput and comprehensive test menu. Ideal for medium to large laboratories.',
    '{
        "throughput": "400 tests/hour (800 with ISE)",
        "test_positions": 64,
        "sample_positions": 80,
        "reagent_positions": 70,
        "sample_volume": "2-45 Î¼L",
        "reagent_volume": "20-350 Î¼L",
        "reaction_volume": "180-500 Î¼L",
        "wavelengths": ["340", "380", "405", "450", "480", "505", "546", "570", "600", "660", "700", "800"],
        "light_source": "Halogen lamp with 2000 hours life",
        "detection_system": "Photometric range: 0-3.5 OD",
        "temperature_control": "37Â°C Â± 0.1Â°C",
        "sample_types": ["Serum", "Plasma", "Urine", "CSF"],
        "test_methods": ["End-point", "Fixed-time", "Kinetic", "Two-point kinetic", "Non-linear multipoint calibration"],
        "onboard_refrigeration": true,
        "auto_dilution": true,
        "barcode_reader": true,
        "interface": "Bi-directional LIS interface",
        "display": "15\" color touchscreen",
        "data_storage": "100,000 patient results",
        "qc_management": "Westgard multi-rules, Levy-Jennings plots",
        "dimensions": "1200 x 750 x 1150 mm",
        "weight": "180 kg",
        "power_requirements": "220-240V, 50/60 Hz",
        "water_consumption": "< 5 L/hour",
        "noise_level": "< 65 dB",
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "FDA"]
    }',
    1850000.00,
    'INR',
    'TRA-BIO-XL640',
    ARRAY['https://ServQR.com/images/equipment/transasia-xl640-1.jpg', 'https://ServQR.com/images/equipment/transasia-xl640-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH14',
    'Transasia Hematology Analyzer',
    'H 560',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ19', -- Hematology Analyzers category
    '01HFPQ2Z5N7YWGXS9JVKM6F8QT', -- Transasia Bio-Medicals
    '5-part differential hematology analyzer with advanced technology for precise blood cell analysis. Suitable for mid to high-volume laboratories.',
    '{
        "parameters": 29,
        "throughput": "60 samples/hour",
        "sample_volume": {
            "whole_blood": "20 Î¼L",
            "pre-diluted": "20 Î¼L"
        },
        "measuring_principle": {
            "wbc_diff": "Laser light scatter, cytochemistry",
            "rbc_plt": "Impedance method",
            "hgb": "Cyanide-free colorimetric method"
        },
        "linearity_ranges": {
            "wbc": "0.3-100.0 x 10^9/L",
            "rbc": "0.3-8.00 x 10^12/L",
            "hgb": "10-250 g/L",
            "plt": "10-1000 x 10^9/L"
        },
        "precision": {
            "wbc": "â‰¤ 2.0%",
            "rbc": "â‰¤ 1.5%",
            "hgb": "â‰¤ 1.5%",
            "plt": "â‰¤ 4.0%"
        },
        "sample_modes": ["Whole blood", "Pre-diluted", "Capillary"],
        "data_storage": "100,000 results including graphics",
        "display": "10.4\" color touchscreen LCD",
        "connectivity": ["LIS", "HIS", "USB", "Ethernet"],
        "qc_management": "Levey-Jennings, XB, XR",
        "dimensions": "540 x 490 x 475 mm",
        "weight": "35 kg",
        "power_requirements": "100-240V, 50/60 Hz",
        "reagents": ["Diluent", "Lyse", "Cleaner"],
        "barcode_reader": true,
        "auto_sampler": "50 positions",
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "FDA"]
    }',
    1250000.00,
    'INR',
    'TRA-HEM-H560',
    ARRAY['https://ServQR.com/images/equipment/transasia-h560-1.jpg', 'https://ServQR.com/images/equipment/transasia-h560-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH15',
    'Transasia Coagulation Analyzer',
    'ECL 760',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ18', -- Clinical Laboratory Analyzers category
    '01HFPQ2Z5N7YWGXS9JVKM6F8QT', -- Transasia Bio-Medicals
    'Fully automated coagulation analyzer with high precision and comprehensive test menu. Features advanced optical detection system for accurate results.',
    '{
        "throughput": "120 PT tests/hour",
        "test_parameters": ["PT", "APTT", "Fibrinogen", "TT", "Factors", "Proteins C & S", "AT-III", "D-dimer", "Lupus anticoagulant"],
        "measuring_principle": "Photo-optical clot detection",
        "wavelength": "405 nm, 570 nm, 740 nm",
        "sample_positions": 40,
        "reagent_positions": 20,
        "cuvette_positions": 400,
        "sample_volume": "2-100 Î¼L",
        "reagent_volume": "5-300 Î¼L",
        "detection_methods": ["Clotting", "Chromogenic", "Immunological"],
        "onboard_refrigeration": true,
        "sample_identification": "Barcode reader",
        "calibration": "Auto-calibration with multi-point curves",
        "qc_management": "Levey-Jennings, Westgard multi-rules",
        "data_storage": "50,000 results",
        "connectivity": ["LIS", "HIS", "USB", "Ethernet"],
        "display": "12.1\" color touchscreen",
        "dimensions": "850 x 700 x 600 mm",
        "weight": "85 kg",
        "power_requirements": "220-240V, 50/60 Hz",
        "backup_power": "30 minutes UPS (optional)",
        "warranty": "2 years",
        "certifications": ["CE", "ISO 13485"]
    }',
    950000.00,
    'INR',
    'TRA-COA-ECL760',
    ARRAY['https://ServQR.com/images/equipment/transasia-ecl760-1.jpg', 'https://ServQR.com/images/equipment/transasia-ecl760-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH16',
    'Transasia Immunoassay Analyzer',
    'Erba ELAN 30s',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ18', -- Clinical Laboratory Analyzers category
    '01HFPQ2Z5N7YWGXS9JVKM6F8QT', -- Transasia Bio-Medicals
    'Compact chemiluminescence immunoassay analyzer for hormones, tumor markers, infectious diseases, and cardiac markers. Ideal for small to medium laboratories.',
    '{
        "technology": "Chemiluminescence",
        "throughput": "30 tests/hour",
        "test_menu": ["Thyroid panel", "Fertility hormones", "Tumor markers", "Cardiac markers", "Infectious diseases"],
        "sample_capacity": 15,
        "reagent_capacity": 6,
        "sample_volume": "5-100 Î¼L",
        "sample_types": ["Serum", "Plasma"],
        "reaction_time": "15-30 minutes",
        "detection_limit": "Down to 0.001 ng/mL",
        "precision": "CV < 5%",
        "calibration_stability": "28 days",
        "onboard_stability": "28 days",
        "data_storage": "10,000 results",
        "qc_management": "Levey-Jennings, Westgard multi-rules",
        "connectivity": ["LIS", "USB"],
        "display": "10.4\" color touchscreen",
        "dimensions": "600 x 550 x 600 mm",
        "weight": "65 kg",
        "power_requirements": "100-240V, 50/60 Hz",
        "noise_level": "< 60 dB",
        "warranty": "2 years",
        "certifications": ["CE", "ISO 13485"]
    }',
    1450000.00,
    'INR',
    'TRA-IMM-ELAN30',
    ARRAY['https://ServQR.com/images/equipment/transasia-elan30-1.jpg', 'https://ServQR.com/images/equipment/transasia-elan30-2.jpg'],
    TRUE,
    'demo-hospital'
),

-- Insert Hospital Beds from Janak Healthcare
(
    '01HFPQZ6QBWK7NQZXT5G7JH17',
    'Janak ICU Electric Bed',
    'Critical Care 5000',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ81', -- Hospital Beds category
    '01HFPQ2Z6EJKWOTZBS7M9HGF', -- Janak Healthcare
    'Advanced electric ICU bed with comprehensive patient positioning capabilities and integrated weighing system. Designed for critical care environments.',
    '{
        "positions": ["Trendelenburg", "Reverse Trendelenburg", "Cardiac chair", "CPR", "Vascular"],
        "controls": {
            "patient": "Handset with selective lockout",
            "nurse": "Control panel with full functions",
            "foot": "Emergency CPR lever"
        },
        "motors": "4 Linak actuators",
        "battery_backup": "24V, up to 8 hours",
        "sections": 4,
        "height_adjustment": "45-80 cm",
        "safe_working_load": "250 kg",
        "mattress_platform": "200 x 90 cm",
        "side_rails": "Split collapsible with safety lock",
        "castors": "150 mm diameter with central locking",
        "bumpers": "Corner protection on all sides",
        "iv_pole": "2 locations with height adjustment",
        "weighing_system": "Integrated with accuracy Â±500g",
        "x_ray_compatibility": "Cassette holder for chest X-ray",
        "dimensions": "220 x 100 x 45-80 cm",
        "frame_material": "Epoxy-coated steel",
        "mattress_surface": "Removable ABS panels",
        "head_foot_boards": "Detachable high-impact ABS",
        "accessories_included": ["IV pole", "Drainage bag hooks", "Lifting pole"],
        "warranty": "5 years on frame, 3 years on electrical components",
        "certifications": ["CE", "ISO 13485", "IEC 60601-2-52"]
    }',
    185000.00,
    'INR',
    'JNK-BED-CC5000',
    ARRAY['https://ServQR.com/images/equipment/janak-cc5000-1.jpg', 'https://ServQR.com/images/equipment/janak-cc5000-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH18',
    'Janak Semi-Electric Hospital Bed',
    'MediCare 3000',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ81', -- Hospital Beds category
    '01HFPQ2Z6EJKWOTZBS7M9HGF', -- Janak Healthcare
    'Semi-electric hospital bed with electric backrest and height adjustment, manual knee adjustment. Ideal for general wards and long-term care facilities.',
    '{
        "type": "Semi-electric",
        "electric_functions": ["Backrest", "Height adjustment"],
        "manual_functions": ["Knee break", "Trendelenburg"],
        "controls": "Handset with nurse lockout",
        "motors": "2 Linak actuators",
        "battery_backup": "24V, up to 4 hours",
        "sections": 3,
        "height_adjustment": "45-70 cm",
        "backrest_angle": "0-70Â°",
        "knee_break_angle": "0-35Â°",
        "trendelenburg": "Â±12Â°",
        "safe_working_load": "200 kg",
        "mattress_platform": "200 x 90 cm",
        "side_rails": "Collapsible full-length",
        "castors": "125 mm diameter with individual brakes",
        "iv_pole": "2 locations",
        "dimensions": "215 x 95 x 45-70 cm",
        "frame_material": "Epoxy-coated steel",
        "head_foot_boards": "Detachable ABS",
        "accessories_included": ["IV pole", "Drainage bag hooks"],
        "warranty": "5 years on frame, 2 years on electrical components",
        "certifications": ["CE", "ISO 13485"]
    }',
    85000.00,
    'INR',
    'JNK-BED-MC3000',
    ARRAY['https://ServQR.com/images/equipment/janak-mc3000-1.jpg', 'https://ServQR.com/images/equipment/janak-mc3000-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH19',
    'Janak Pediatric Hospital Bed',
    'KidCare 2000',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ81', -- Hospital Beds category
    '01HFPQ2Z6EJKWOTZBS7M9HGF', -- Janak Healthcare
    'Specialized pediatric hospital bed with transparent acrylic side panels for visibility and safety. Features colorful design and child-friendly features.',
    '{
        "type": "Electric pediatric bed",
        "electric_functions": ["Height adjustment", "Backrest"],
        "controls": "Nurse control panel with lockout",
        "motors": "2 Linak actuators",
        "battery_backup": "24V, up to 6 hours",
        "sections": 2,
        "height_adjustment": "65-95 cm",
        "backrest_angle": "0-60Â°",
        "safe_working_load": "100 kg",
        "mattress_platform": "160 x 75 cm",
        "side_rails": "Transparent acrylic panels with safety lock",
        "side_rail_height": "60 cm from mattress platform",
        "castors": "100 mm diameter with central locking",
        "trendelenburg": "Â±12Â°",
        "iv_pole": "2 locations",
        "dimensions": "175 x 85 x 65-95 cm",
        "frame_material": "Epoxy-coated steel with colorful panels",
        "head_foot_boards": "Detachable colorful ABS",
        "special_features": ["Under-bed lighting", "Night light", "Child-friendly graphics"],
        "accessories_included": ["IV pole", "Chart holder", "Toy storage"],
        "warranty": "5 years on frame, 2 years on electrical components",
        "certifications": ["CE", "ISO 13485", "IEC 60601-2-52"]
    }',
    120000.00,
    'INR',
    'JNK-BED-KC2000',
    ARRAY['https://ServQR.com/images/equipment/janak-kc2000-1.jpg', 'https://ServQR.com/images/equipment/janak-kc2000-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH20',
    'Janak Hydraulic Examination Table',
    'ExamPro 1500',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ84', -- Examination Tables category
    '01HFPQ2Z6EJKWOTZBS7M9HGF', -- Janak Healthcare
    'Versatile hydraulic examination table with multiple positioning options. Features sturdy construction and comfortable patient surface.',
    '{
        "type": "Hydraulic examination table",
        "adjustment_mechanism": "Hydraulic foot pump",
        "positions": ["Flat", "Chair", "Trendelenburg", "Reverse Trendelenburg"],
        "height_adjustment": "50-95 cm",
        "backrest_angle": "0-80Â°",
        "trendelenburg": "Â±15Â°",
        "safe_working_load": "180 kg",
        "table_top": "2-section",
        "padding": "65 mm high-density foam",
        "upholstery": "Seamless leatherette, antimicrobial",
        "upholstery_colors": ["Blue", "Green", "Grey", "Black", "Burgundy"],
        "paper_roll_holder": "Integrated at head end",
        "base": "Heavy-duty steel with protective cover",
        "mobility": "Retractable castors",
        "dimensions": "190 x 65 x 50-95 cm",
        "storage": "2 side drawers",
        "stirrups": "Adjustable and retractable",
        "accessories_included": ["Paper roll holder", "Drainage pan"],
        "optional_accessories": ["IV pole", "Side rails", "Arm boards"],
        "warranty": "5 years on structure, 2 years on hydraulic system",
        "certifications": ["CE", "ISO 13485"]
    }',
    65000.00,
    'INR',
    'JNK-EXT-EP1500',
    ARRAY['https://ServQR.com/images/equipment/janak-ep1500-1.jpg', 'https://ServQR.com/images/equipment/janak-ep1500-2.jpg'],
    TRUE,
    'demo-hospital'
),

-- Insert ECG Machines from Schiller and BPL
(
    '01HFPQZ6QBWK7NQZXT5G7JH21',
    'Schiller Cardiovit AT-102',
    'AT-102 Plus',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ65', -- Electrocardiographs (ECG/EKG) category
    '01HFPQ2Z5YQJK8TZBS7M9HGF', -- Schiller Healthcare India
    'Advanced 12-channel ECG machine with interpretation capabilities. Features high-resolution display and comprehensive analysis algorithms.',
    '{
        "channels": 12,
        "recording": "3/6/12 channel simultaneous",
        "display": "8.9\" high-resolution color LCD",
        "paper_speed": ["5 mm/s", "10 mm/s", "25 mm/s", "50 mm/s"],
        "sensitivity": ["2.5 mm/mV", "5 mm/mV", "10 mm/mV", "20 mm/mV"],
        "filters": {
            "baseline": "0.05-150 Hz",
            "myogram": "25/35 Hz",
            "ac": "50/60 Hz"
        },
        "sampling_rate": "32,000 samples/second/channel",
        "resolution": "5 Î¼V/LSB (24-bit ADC)",
        "cmrr": "> 115 dB",
        "pacemaker_detection": true,
        "interpretation": "SCHILLER ECG Analysis Program",
        "memory": "500 ECGs internal storage",
        "connectivity": ["USB", "LAN", "WLAN (optional)", "DICOM", "HL7", "GDT"],
        "paper": "Z-fold thermal paper, 210 mm width",
        "operation": {
            "mains": "100-240V, 50/60 Hz",
            "battery": "Lithium-ion, up to 6 hours operation"
        },
        "dimensions": "360 x 290 x 70 mm",
        "weight": "3.9 kg",
        "special_features": [
            "Resting ECG analysis",
            "Exercise ECG (optional)",
            "Spirometry (optional)",
            "Thrombolysis guidance"
        ],
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "FDA"]
    }',
    325000.00,
    'INR',
    'SCH-ECG-AT102P',
    ARRAY['https://ServQR.com/images/equipment/schiller-at102-1.jpg', 'https://ServQR.com/images/equipment/schiller-at102-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH22',
    'Schiller Cardiovit MS-2015',
    'MS-2015',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ65', -- Electrocardiographs (ECG/EKG) category
    '01HFPQ2Z5YQJK8TZBS7M9HGF', -- Schiller Healthcare India
    'Premium 16-channel ECG machine with touchscreen interface and advanced networking capabilities. Suitable for high-volume cardiac departments.',
    '{
        "channels": 16,
        "recording": "3/6/12/16 channel simultaneous",
        "display": "15\" high-resolution color touchscreen",
        "paper_speed": ["5 mm/s", "10 mm/s", "25 mm/s", "50 mm/s", "100 mm/s"],
        "sensitivity": ["2.5 mm/mV", "5 mm/mV", "10 mm/mV", "20 mm/mV"],
        "filters": {
            "baseline": "0.05-150 Hz",
            "myogram": "25/35/45 Hz",
            "ac": "50/60 Hz"
        },
        "sampling_rate": "32,000 samples/second/channel",
        "resolution": "1 Î¼V/LSB (32-bit ADC)",
        "cmrr": "> 140 dB",
        "pacemaker_detection": true,
        "interpretation": "SCHILLER ECG Analysis Program with C.A.R.E. algorithm",
        "memory": "10,000 ECGs internal storage",
        "connectivity": ["USB", "LAN", "WLAN", "DICOM", "HL7", "GDT", "MQTT"],
        "paper": "Z-fold thermal paper, 210 mm width",
        "operation": {
            "mains": "100-240V, 50/60 Hz",
            "battery": "Lithium-ion, up to 8 hours operation"
        },
        "dimensions": "400 x 330 x 95 mm",
        "weight": "5.8 kg",
        "special_features": [
            "Resting ECG analysis",
            "Exercise ECG",
            "Spirometry",
            "Thrombolysis guidance",
            "Vector ECG",
            "HRV analysis",
            "Signal-averaged ECG",
            "QT dispersion"
        ],
        "warranty": "5 years",
        "certifications": ["CE", "ISO 13485", "FDA"]
    }',
    575000.00,
    'INR',
    'SCH-ECG-MS2015',
    ARRAY['https://ServQR.com/images/equipment/schiller-ms2015-1.jpg', 'https://ServQR.com/images/equipment/schiller-ms2015-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH23',
    'BPL Cardiart 6208 View',
    '6208 View',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ65', -- Electrocardiographs (ECG/EKG) category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'Reliable 12-channel ECG machine with interpretation and large display. Features user-friendly interface and comprehensive connectivity options.',
    '{
        "channels": 12,
        "recording": "3/6/12 channel simultaneous",
        "display": "7\" color TFT LCD",
        "paper_speed": ["5 mm/s", "10 mm/s", "25 mm/s", "50 mm/s"],
        "sensitivity": ["2.5 mm/mV", "5 mm/mV", "10 mm/mV", "20 mm/mV", "Auto"],
        "filters": {
            "baseline": "0.05-150 Hz",
            "myogram": "25/35/45 Hz",
            "ac": "50/60 Hz"
        },
        "sampling_rate": "16,000 samples/second/channel",
        "resolution": "5 Î¼V/LSB (24-bit ADC)",
        "cmrr": "> 105 dB",
        "pacemaker_detection": true,
        "interpretation": "Glasgow algorithm",
        "memory": "200 ECGs internal storage",
        "connectivity": ["USB", "LAN", "Bluetooth (optional)"],
        "paper": "Z-fold thermal paper, 210 mm width",
        "operation": {
            "mains": "100-240V, 50/60 Hz",
            "battery": "Lithium-ion, up to 4 hours operation"
        },
        "dimensions": "320 x 240 x 80 mm",
        "weight": "3.2 kg",
        "special_features": [
            "Resting ECG analysis",
            "Arrhythmia detection",
            "Pediatric interpretation"
        ],
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    175000.00,
    'INR',
    'BPL-ECG-6208V',
    ARRAY['https://ServQR.com/images/equipment/bpl-6208-1.jpg', 'https://ServQR.com/images/equipment/bpl-6208-2.jpg'],
    TRUE,
    'demo-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH24',
    'BPL Cardiart 9108 Portable',
    '9108 Pocket',
    '01HFPQ3Z5MXNVT9DAPZ3BWJ65', -- Electrocardiographs (ECG/EKG) category
    '01HFPQ2Z5NWBCPXQ2RJVT3D7KF', -- BPL Medical Technologies
    'Compact and lightweight single-channel ECG machine for basic cardiac screening. Ideal for rural healthcare, home visits, and ambulatory care.',
    '{
        "channels": 1,
        "recording": "Single channel",
        "display": "4.3\" monochrome LCD",
        "paper_speed": ["25 mm/s", "50 mm/s"],
        "sensitivity": ["5 mm/mV", "10 mm/mV", "20 mm/mV"],
        "filters": {
            "baseline": "0.05-100 Hz",
            "myogram": "25/35 Hz",
            "ac": "50/60 Hz"
        },
        "sampling_rate": "8,000 samples/second/channel",
        "resolution": "10 Î¼V/LSB (16-bit ADC)",
        "cmrr": "> 95 dB",
        "pacemaker_detection": false,
        "memory": "50 ECGs internal storage",
        "connectivity": ["USB"],
        "paper": "Roll thermal paper, 50 mm width",
        "operation": {
            "mains": "100-240V, 50/60 Hz",
            "battery": "Lithium-ion, up to 8 hours operation"
        },
        "dimensions": "220 x 140 x 65 mm",
        "weight": "1.2 kg",
        "special_features": [
            "Auto/manual mode",
            "Rhythm recording",
            "Lead-off detection"
        ],
        "warranty": "2 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    45000.00,
    'INR',
    'BPL-ECG-9108P',
    ARRAY['https://ServQR.com/images/equipment/bpl-9108-1.jpg', 'https://ServQR.com/images/equipment/bpl-9108-2.jpg'],
    TRUE,
    'demo-hospital'
);

-- Add some equipment for city-hospital tenant
INSERT INTO equipment (id, name, model, category_id, manufacturer_id, description, specifications, price_amount, price_currency, sku, images, is_active, tenant_id) VALUES
(
    '01HFPQZ6QBWK7NQZXT5G7JH25',
    'Trivitron Elite Patient Monitor',
    'Clarity-Pro 12',
    '01HFPQ3Z5MXNVT9DAPZ3BW129', -- Vital Signs Monitors category (city-hospital)
    '01HFPQ2Z6QWTH8VPZBS3F8KMN', -- Trivitron Healthcare (city-hospital)
    'Advanced multi-parameter patient monitor with comprehensive monitoring capabilities for critical care, emergency, and general ward settings.',
    '{
        "display": "12.1\" TFT color touchscreen",
        "resolution": "1280 x 800 pixels",
        "parameters": ["ECG", "NIBP", "SpO2", "Respiration", "Temperature", "IBP", "EtCO2"],
        "ecg_leads": "3/5/12-lead selectable",
        "ecg_analysis": "ST segment analysis, arrhythmia detection",
        "nibp_measurement": "Oscillometric method with auto/manual/continuous modes",
        "spo2_technology": "Nellcor compatible",
        "temperature_channels": 2,
        "ibp_channels": 2,
        "etco2_method": "Sidestream/Mainstream selectable",
        "trends": "120 hours graphical and tabular",
        "alarms": "3-level visual and audible",
        "connectivity": ["Wi-Fi", "Ethernet", "HL7"],
        "battery": "Lithium-ion, up to 4 hours operation",
        "weight": "4.5 kg",
        "dimensions": "318 x 264 x 152 mm",
        "storage": "32GB internal, expandable via USB",
        "printer": "Built-in thermal printer",
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "BIS"]
    }',
    285000.00,
    'INR',
    'TRV-PM-CP12',
    ARRAY['https://ServQR.com/images/equipment/trivitron-claritypro-1.jpg', 'https://ServQR.com/images/equipment/trivitron-claritypro-2.jpg'],
    TRUE,
    'city-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH26',
    'BPL Digital X-Ray System',
    'DX-5100',
    '01HFPQ3Z5MXNVT9DAPZ3BW126', -- X-ray Systems category (city-hospital)
    '01HFPQ2Z6RJKW9TZBS7M9HGF', -- BPL Medical Technologies (city-hospital)
    'Advanced digital X-ray system with high-resolution imaging and reduced radiation exposure. Suitable for general radiography in hospitals and diagnostic centers.',
    '{
        "detector_type": "Flat Panel Digital Detector",
        "detector_size": "43cm x 43cm (17\" x 17\")",
        "resolution": "3.5 lp/mm",
        "pixel_size": "140 Î¼m",
        "image_depth": "16-bit",
        "generator_power": "50 kW",
        "tube_voltage_range": "40-150 kV",
        "tube_current_range": "10-630 mA",
        "exposure_time": "1 ms to 5 sec",
        "weight_capacity": "200 kg",
        "dimensions": "2200 x 1800 x 2300 mm",
        "power_requirements": "380-400V, 50/60 Hz, three phase",
        "dicom_compatible": true,
        "warranty": "3 years",
        "certifications": ["CE", "ISO 13485", "AERB"]
    }',
    1850000.00,
    'INR',
    'BPL-XR-DX5100',
    ARRAY['https://ServQR.com/images/equipment/bpl-dx5100-1.jpg', 'https://ServQR.com/images/equipment/bpl-dx5100-2.jpg'],
    TRUE,
    'city-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH27',
    'Hindustan Syringes Auto-Disable Syringe',
    'Kojak AD',
    '01HFPQ3Z5MXNVT9DAPZ3BW128', -- General Surgical Instruments category (city-hospital)
    '01HFPQ2Z6S4MXAGQNBT5F8HR', -- Hindustan Syringes & Medical Devices (city-hospital)
    'Auto-disable syringes that automatically lock after single use, preventing reuse and reducing the risk of cross-contamination and infections.',
    '{
        "sizes": ["0.5ml", "1ml", "2ml", "3ml", "5ml", "10ml"],
        "needle_sizes": ["23G", "24G", "25G", "26G"],
        "material": "Medical-grade polypropylene",
        "sterilization": "ETO sterilized",
        "shelf_life": "5 years",
        "locking_mechanism": "Automatic plunger lock after injection",
        "packaging": "Individual blister packs",
        "box_quantity": 100,
        "case_quantity": 2000,
        "latex_free": true,
        "phthalate_free": true,
        "color_coded": true,
        "certifications": ["ISO 13485", "CE", "WHO-PQS"]
    }',
    1200.00, -- Price per box of 100
    'INR',
    'HMD-SYR-KOJAK-2ML',
    ARRAY['https://ServQR.com/images/equipment/hmd-kojak-1.jpg', 'https://ServQR.com/images/equipment/hmd-kojak-2.jpg'],
    TRUE,
    'city-hospital'
),
(
    '01HFPQZ6QBWK7NQZXT5G7JH28',
    'Janak ICU Electric Bed',
    'Critical Care 5000',
    '01HFPQ3Z5MXNVT9DAPZ3BW131', -- Hospital Beds category (city-hospital)
    '01HFPQ2Z6T8VFCPXZRT4K7MHG', -- Janak Healthcare (city-hospital)
    'Advanced electric ICU bed with comprehensive patient positioning capabilities and integrated weighing system. Designed for critical care environments.',
    '{
        "positions": ["Trendelenburg", "Reverse Trendelenburg", "Cardiac chair", "CPR", "Vascular"],
        "controls": {
            "patient": "Handset with selective lockout",
            "nurse": "Control panel with full functions",
            "foot": "Emergency CPR lever"
        },
        "motors": "4 Linak actuators",
        "battery_backup": "24V, up to 8 hours",
        "sections": 4,
        "height_adjustment": "45-80 cm",
        "safe_working_load": "250 kg",
        "mattress_platform": "200 x 90 cm",
        "side_rails": "Split collapsible with safety lock",
        "castors": "150 mm diameter with central locking",
        "bumpers": "Corner protection on all sides",
        "iv_pole": "2 locations with height adjustment",
        "weighing_system": "Integrated with accuracy Â±500g",
        "x_ray_compatibility": "Cassette holder for chest X-ray",
        "dimensions": "220 x 100 x 45-80 cm",
        "frame_material": "Epoxy-coated steel",
        "mattress_surface": "Removable ABS panels",
        "head_foot_boards": "Detachable high-impact ABS",
        "accessories_included": ["IV pole", "Drainage bag hooks", "Lifting pole"],
        "warranty": "5 years on frame, 3 years on electrical components",
        "certifications": ["CE", "ISO 13485", "IEC 60601-2-52"]
    }',
    185000.00,
    'INR',
    'JNK-BED-CC5000',
    ARRAY['https://ServQR.com/images/equipment/janak-cc5000-1.jpg', 'https://ServQR.com/images/equipment/janak-cc5000-2.jpg'],
    TRUE,
    'city-hospital'
);
