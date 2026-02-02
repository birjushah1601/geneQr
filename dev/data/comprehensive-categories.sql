-- =============================================================================
-- ServQR Platform - Comprehensive Medical Equipment Categories
-- =============================================================================
-- This file contains an expanded taxonomy of medical equipment categories
-- covering all aspects of healthcare including dental, laboratory, hospital
-- infrastructure, emergency care, and specialized medical equipment.
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS public;

-- Create categories table if it doesn't exist
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(26),
    description TEXT,
    tenant_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- =============================================================================
-- TOP LEVEL CATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Insert top-level categories for demo-hospital tenant
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGA', 'Diagnostic Equipment', NULL, 'Equipment used for diagnosing medical conditions including imaging systems, laboratory analyzers, and diagnostic tools.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB', 'Surgical Instruments & Devices', NULL, 'Tools and equipment used during surgical procedures including scalpels, forceps, retractors, and specialized surgical devices.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC', 'Patient Monitoring Systems', NULL, 'Equipment used to continuously monitor patient vital signs and physiological parameters in various clinical settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD', 'Rehabilitation Equipment', NULL, 'Devices and equipment used for physical therapy, rehabilitation, and recovery from injuries or medical conditions.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Dental Equipment', NULL, 'Specialized tools and devices used in dental procedures, examinations, and treatments.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Laboratory Equipment', NULL, 'Equipment and devices used in medical laboratories for testing, analysis, and research.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Hospital Infrastructure', NULL, 'Essential equipment and systems for hospital operations, patient care areas, and support services.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Emergency & Critical Care', NULL, 'Equipment specifically designed for emergency situations, critical care, and life support.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Laboratory Materials & Consumables', NULL, 'Consumable supplies, reagents, and materials used in laboratory procedures and testing.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized Medical Equipment', NULL, 'Equipment designed for specific medical specialties and specialized procedures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK', 'Medical Imaging Equipment', NULL, 'Advanced imaging technologies for diagnostic and therapeutic purposes.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGL', 'Medical Furniture', NULL, 'Specialized furniture designed for healthcare settings and patient care.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGM', 'Medical Consumables', NULL, 'Disposable and single-use medical supplies used in patient care and procedures.', 'demo-hospital');

-- =============================================================================
-- 1. DENTAL EQUIPMENT SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Dental Equipment Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHG1', 'Dental Chairs & Units', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Integrated dental treatment units including patient chair, delivery systems, and associated equipment.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG2', 'Dental X-Ray Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Radiographic imaging systems specifically designed for dental diagnostics including panoramic, cephalometric, and cone beam CT.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG3', 'Dental Sterilizers', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Equipment used for sterilization of dental instruments and materials including autoclaves and ultrasonic cleaners.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG4', 'Dental Handpieces', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'High-speed and low-speed rotary instruments used for various dental procedures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG5', 'Dental Compressors', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Air compression systems that power dental handpieces and other pneumatic dental equipment.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG6', 'Dental Lights', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Specialized lighting systems designed for optimal illumination during dental procedures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG7', 'Orthodontic Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Specialized tools and devices used in orthodontic treatments including braces, aligners, and associated instruments.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG8', 'Endodontic Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Specialized instruments for root canal procedures including files, reamers, and apex locators.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHG9', 'Periodontal Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Instruments and devices used for diagnosis and treatment of periodontal diseases.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGA1', 'Oral Surgery Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Specialized tools for dental extractions, implant placement, and other oral surgical procedures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGA2', 'Dental Imaging Software', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Specialized software for viewing, analyzing, and managing dental radiographs and 3D images.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGA3', 'Dental Laboratory Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Equipment used in dental labs for fabrication of crowns, bridges, dentures, and other dental prosthetics.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGA4', 'Dental CAD/CAM Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Computer-aided design and manufacturing systems for creating dental restorations.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGA5', 'Dental Lasers', '01HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Laser devices used for various dental procedures including soft tissue surgery and tooth whitening.', 'demo-hospital');

-- =============================================================================
-- 2. LABORATORY EQUIPMENT SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Laboratory Equipment Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGB1', 'Clinical Microscopes', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Optical instruments used for magnified observation of specimens in clinical and research settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB2', 'Centrifuges', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Devices that separate components of a fluid through centrifugal force for analysis and processing.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB3', 'Autoclaves & Sterilizers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Equipment used for sterilization of laboratory instruments, glassware, and materials.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB4', 'Incubators', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Devices that maintain optimal temperature, humidity, and other conditions for biological samples and cultures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB5', 'Spectrophotometers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Instruments that measure the intensity of light absorbed by a solution to determine concentration of substances.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB6', 'Analytical Balances', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Precision weighing instruments used for accurate measurement of mass in laboratory settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB7', 'Water Baths', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Laboratory equipment used to maintain water at a constant temperature for heating samples.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB8', 'Refrigerators & Freezers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Specialized cold storage equipment for preserving samples, reagents, and temperature-sensitive materials.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGB9', 'Biosafety Cabinets', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Enclosed, ventilated laboratory workspaces for safely working with materials contaminated with pathogens.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC1', 'PCR Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Equipment used for polymerase chain reaction to amplify DNA segments for analysis.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC2', 'ELISA Analyzers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Instruments that perform enzyme-linked immunosorbent assay tests for detecting and quantifying substances.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC3', 'Hematology Analyzers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Automated systems for counting and characterizing blood cells for diagnostic purposes.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC4', 'Biochemistry Analyzers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Instruments that measure various chemicals and biochemical markers in biological samples.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC5', 'Immunoassay Analyzers', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Systems that detect and measure specific proteins using antibody-antigen reactions.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC6', 'Microbiology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Specialized tools and instruments for culturing, identifying, and testing microorganisms.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC7', 'Molecular Diagnostic Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Advanced equipment for detecting and analyzing genetic material for diagnostic purposes.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGC8', 'Histology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Instruments used for preparing and examining tissue samples for microscopic study.', 'demo-hospital');

-- =============================================================================
-- 3. LABORATORY MATERIALS & CONSUMABLES SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Laboratory Materials & Consumables Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGD1', 'Test Tubes & Vials', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Glass or plastic containers used for collecting, holding, and processing laboratory samples.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD2', 'Pipettes & Tips', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Instruments used for transferring precise volumes of liquids in laboratory settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD3', 'Reagents & Chemicals', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Substances used in laboratory tests and procedures to detect, measure, or produce chemical reactions.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD4', 'Culture Media', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Nutrient preparations used for growing and cultivating microorganisms in laboratory settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD5', 'Disposable Gloves', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Single-use hand protection for laboratory and medical procedures.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD6', 'Laboratory Glassware', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Specialized glass containers and instruments used in laboratory procedures including beakers, flasks, and graduated cylinders.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD7', 'Sample Collection', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Devices and containers used for collecting biological samples including blood collection tubes and swabs.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD8', 'Safety Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Protective gear and equipment for laboratory safety including eye protection, lab coats, and spill kits.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGD9', 'Filtration Supplies', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Materials used for separating solids from liquids including filter papers, membranes, and filtration units.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGE1', 'Microbiology Consumables', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Specialized supplies for microbiology work including petri dishes, inoculation loops, and swabs.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGE2', 'Molecular Biology Reagents', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Specialized chemicals and kits used in DNA/RNA extraction, amplification, and analysis.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGE3', 'Chromatography Supplies', '01HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Materials used in separation and analysis techniques including columns, plates, and mobile phases.', 'demo-hospital');

-- =============================================================================
-- 4. HOSPITAL INFRASTRUCTURE SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Hospital Infrastructure Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGF1', 'Hospital Beds', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized beds designed for hospitalized patients with features for patient comfort and clinical care.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF2', 'Patient Trolleys', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Mobile platforms for transporting patients within healthcare facilities.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF3', 'Medical Furniture', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized furniture designed for healthcare settings including cabinets, carts, and storage systems.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF4', 'Operation Theatre Tables', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized tables designed for surgical procedures with positioning capabilities and accessories.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF5', 'Medical Lighting', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized lighting systems for examination rooms, operation theatres, and other clinical areas.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF6', 'Medical Gas Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Infrastructure for delivering medical gases including oxygen, nitrous oxide, and medical air to patient care areas.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF7', 'HVAC Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized heating, ventilation, and air conditioning systems for healthcare facilities with infection control features.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF8', 'Hospital Elevators', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized elevators designed for healthcare settings with features for patient transport and emergency use.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGF9', 'Nurse Call Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Communication systems that allow patients to alert nursing staff when assistance is needed.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGG1', 'Sterilization Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Equipment used for sterilizing medical instruments and supplies in hospital settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGG2', 'Waste Management Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Specialized systems for handling, treating, and disposing of medical waste.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGG3', 'Hospital Security Systems', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Access control, surveillance, and alarm systems designed for healthcare facility security.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGG4', 'Patient Room Fixtures', '01HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Built-in equipment and fixtures for patient rooms including headwalls and service columns.', 'demo-hospital');

-- =============================================================================
-- 5. EMERGENCY & CRITICAL CARE SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Emergency & Critical Care Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGH1', 'Defibrillators', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Devices that deliver electric shock to restore normal heart rhythm in cardiac emergencies.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH2', 'Ventilators', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Machines that provide mechanical ventilation by moving breathable air into and out of the lungs.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH3', 'Ambulance Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Specialized medical devices and supplies used in ambulances for pre-hospital emergency care.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH4', 'Emergency Trolleys', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Mobile carts containing essential emergency medications and equipment for rapid response to medical emergencies.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH5', 'Resuscitation Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Devices and supplies used in cardiopulmonary resuscitation including bag valve masks and airway management tools.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH6', 'Cardiac Monitors', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Devices that continuously record and display patient heart activity and other vital parameters.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH7', 'Infusion Pumps', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Devices that deliver fluids, medications, or nutrients into a patient's circulatory system in precisely controlled amounts.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH8', 'Suction Units', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Devices used to remove bodily fluids, secretions, or other substances through vacuum aspiration.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGH9', 'Trauma Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Specialized tools and devices for managing traumatic injuries including splints, immobilizers, and tourniquets.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGI1', 'ICU Beds', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Advanced hospital beds with specialized features for intensive care settings.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGI2', 'ECMO Machines', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Extracorporeal membrane oxygenation devices that provide cardiac and respiratory support to patients.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGI3', 'Dialysis Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Machines that filter blood to remove excess water and waste products when kidneys are damaged or dysfunctional.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGI4', 'Transport Monitors', '01HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Portable patient monitoring systems designed for use during patient transport.', 'demo-hospital');

-- =============================================================================
-- 6. SPECIALIZED MEDICAL EQUIPMENT SUBCATEGORIES - DEMO HOSPITAL TENANT
-- =============================================================================

-- Specialized Medical Equipment Subcategories
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ1', 'Cardiology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for diagnosis and treatment of heart conditions including ECG machines, stress test systems, and cardiac catheterization equipment.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ2', 'Neurology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Devices used for diagnosis and monitoring of neurological conditions including EEG, EMG, and nerve conduction study equipment.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ3', 'Orthopedic Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized tools and devices for orthopedic procedures including power tools, implants, and traction equipment.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ4', 'Gynecology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for gynecological examinations and procedures including colposcopes and hysteroscopes.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ5', 'Pediatric Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Medical equipment specifically designed or adapted for use with infants and children.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ6', 'Geriatric Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment designed for elderly patients including mobility aids and monitoring systems.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ7', 'Physical Therapy Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Equipment used in physical rehabilitation including exercise machines, therapeutic modalities, and mobility aids.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ8', 'Ophthalmology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for eye examinations and procedures including slit lamps, phoropters, and ophthalmic lasers.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGJ9', 'ENT Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized tools and devices for diagnosis and treatment of ear, nose, and throat conditions.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK1', 'Dermatology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for skin examination and treatment including dermatoscopes and light therapy systems.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK2', 'Urology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for urological examinations and procedures including cystoscopes and lithotripters.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK3', 'Pulmonology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Equipment for diagnosis and treatment of respiratory conditions including spirometers and bronchoscopes.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK4', 'Oncology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for cancer diagnosis and treatment including radiation therapy and chemotherapy delivery systems.', 'demo-hospital'),
('01HFPQ2Z5P8VF5PXZRT4K7MHGK5', 'Gastroenterology Equipment', '01HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized equipment for diagnosis and treatment of digestive system disorders including endoscopes and colonoscopes.', 'demo-hospital');

-- =============================================================================
-- TOP LEVEL CATEGORIES - CITY HOSPITAL TENANT
-- =============================================================================

-- Insert the same top-level categories for city-hospital tenant (multi-tenancy)
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('02HFPQ2Z5P8VF5PXZRT4K7MHGA', 'Diagnostic Equipment', NULL, 'Equipment used for diagnosing medical conditions including imaging systems, laboratory analyzers, and diagnostic tools.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGB', 'Surgical Instruments & Devices', NULL, 'Tools and equipment used during surgical procedures including scalpels, forceps, retractors, and specialized surgical devices.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGC', 'Patient Monitoring Systems', NULL, 'Equipment used to continuously monitor patient vital signs and physiological parameters in various clinical settings.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGD', 'Rehabilitation Equipment', NULL, 'Devices and equipment used for physical therapy, rehabilitation, and recovery from injuries or medical conditions.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGE', 'Dental Equipment', NULL, 'Specialized tools and devices used in dental procedures, examinations, and treatments.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGF', 'Laboratory Equipment', NULL, 'Equipment and devices used in medical laboratories for testing, analysis, and research.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGG', 'Hospital Infrastructure', NULL, 'Essential equipment and systems for hospital operations, patient care areas, and support services.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGH', 'Emergency & Critical Care', NULL, 'Equipment specifically designed for emergency situations, critical care, and life support.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGI', 'Laboratory Materials & Consumables', NULL, 'Consumable supplies, reagents, and materials used in laboratory procedures and testing.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGJ', 'Specialized Medical Equipment', NULL, 'Equipment designed for specific medical specialties and specialized procedures.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGK', 'Medical Imaging Equipment', NULL, 'Advanced imaging technologies for diagnostic and therapeutic purposes.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGL', 'Medical Furniture', NULL, 'Specialized furniture designed for healthcare settings and patient care.', 'city-hospital'),
('02HFPQ2Z5P8VF5PXZRT4K7MHGM', 'Medical Consumables', NULL, 'Disposable and single-use medical supplies used in patient care and procedures.', 'city-hospital');

-- =============================================================================
-- Create indexes for better performance
-- =============================================================================
CREATE INDEX IF NOT EXISTS idx_categories_tenant ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
