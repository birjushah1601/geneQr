-- Medical Equipment Categories Dataset
-- For ABY-Med Platform Catalog Module
-- This file populates the categories table with hierarchical medical equipment categories

-- First ensure the schema exists
CREATE SCHEMA IF NOT EXISTS public;

-- Create the categories table if it doesn't exist
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

-- Clear existing data if needed (commented out for safety)
-- TRUNCATE TABLE categories CASCADE;

-- Insert main parent categories for demo-hospital tenant
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
-- Main Categories (Parent = NULL)
('01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Diagnostic Equipment', NULL, 'Equipment used for diagnosing medical conditions including imaging systems, laboratory analyzers, and diagnostic tools.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Surgical Instruments & Devices', NULL, 'Tools and equipment used during surgical procedures including scalpels, forceps, retractors, and specialized surgical devices.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Patient Monitoring Systems', NULL, 'Equipment used to continuously monitor patient vital signs and physiological parameters in various clinical settings.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Rehabilitation Equipment', NULL, 'Devices and equipment used for physical therapy, rehabilitation, and recovery from injuries or medical conditions.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Dental Equipment', NULL, 'Specialized tools and devices used in dental procedures, examinations, and treatments.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Ophthalmic Equipment', NULL, 'Specialized instruments and devices used for eye examinations, diagnoses, and ophthalmic surgeries.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Cardiology Equipment', NULL, 'Equipment used for diagnosing and treating heart conditions including ECG machines, cardiac monitors, and catheterization lab equipment.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Orthopedic Devices', NULL, 'Implants, instruments, and equipment used in orthopedic procedures and treatments for bone and joint conditions.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Hospital Furniture & Infrastructure', NULL, 'Furniture and infrastructure components designed specifically for healthcare facilities including beds, stretchers, and cabinets.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Biomedical Equipment', NULL, 'Equipment used in biomedical research, analysis, and applications in healthcare settings.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Respiratory Equipment', NULL, 'Devices used for treating and monitoring respiratory conditions including ventilators, oxygen delivery systems, and nebulizers.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Anesthesia Equipment', NULL, 'Equipment used for administering anesthesia during surgical procedures and monitoring patients under anesthesia.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Neonatal Equipment', NULL, 'Specialized equipment designed for the care of newborn infants, particularly premature babies and those requiring intensive care.', 'demo-hospital');

-- Insert sub-categories for Diagnostic Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ14', 'X-ray Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Equipment that uses X-ray radiation to create images of structures inside the body.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ15', 'Ultrasound Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Imaging devices that use high-frequency sound waves to create images of structures inside the body.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ16', 'MRI Machines', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Magnetic Resonance Imaging equipment that uses magnetic fields and radio waves to create detailed images of organs and tissues.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ17', 'CT Scanners', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Computed Tomography scanners that use X-rays to create cross-sectional images of the body.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ18', 'Clinical Laboratory Analyzers', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Equipment used to analyze blood, urine, and other body fluids for diagnostic purposes.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ19', 'Hematology Analyzers', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Specialized equipment for analyzing blood cells and components.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ20', 'Biochemistry Analyzers', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Equipment for measuring biochemical markers in blood and other body fluids.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ21', 'Endoscopy Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Optical instruments used to examine the interior of hollow organs and cavities of the body.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ22', 'Mammography Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Specialized imaging equipment used for breast cancer screening and diagnosis.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ23', 'PET-CT Scanners', '01HFPQ3Z5MXNVT9DAPZ3BWJA1', 'Combined Positron Emission Tomography and Computed Tomography scanners for functional and anatomical imaging.', 'demo-hospital');

-- Insert sub-categories for Surgical Instruments & Devices
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ24', 'General Surgical Instruments', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Basic instruments used in various surgical procedures including scalpels, forceps, scissors, and clamps.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ25', 'Laparoscopic Instruments', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Specialized instruments used in minimally invasive laparoscopic surgeries.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ26', 'Electrosurgical Units', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Devices that use high-frequency electrical currents to cut tissue and control bleeding during surgery.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ27', 'Surgical Staplers', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Devices used to close wounds or connect tissues during surgical procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ28', 'Surgical Sutures', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Medical devices used to hold body tissues together after an injury or surgery.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ29', 'Surgical Microscopes', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Optical microscopes specifically designed for use in surgical settings.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ30', 'Surgical Navigation Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Computer-assisted technologies that help surgeons navigate during complex procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ31', 'Robotic Surgical Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Advanced robotic platforms that assist surgeons in performing minimally invasive procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ32', 'Surgical Lasers', '01HFPQ3Z5MXNVT9DAPZ3BWJA2', 'Devices that use focused light beams for precise cutting and cauterizing during surgery.', 'demo-hospital');

-- Insert sub-categories for Patient Monitoring Systems
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ33', 'Vital Signs Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Devices that measure and display patient vital signs including heart rate, blood pressure, temperature, and oxygen saturation.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ34', 'ECG Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Equipment that records the electrical activity of the heart over a period of time.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ35', 'Pulse Oximeters', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Devices that measure oxygen saturation levels in the blood.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ36', 'Capnography Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Equipment that measures and displays the concentration of carbon dioxide in respiratory gases.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ37', 'Fetal Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Devices used to monitor fetal heart rate and uterine contractions during pregnancy and labor.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ38', 'Telemetry Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Wireless monitoring systems that allow continuous monitoring of patients while they move around.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ39', 'Intracranial Pressure Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Devices that measure pressure within the skull and brain.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ40', 'Multi-Parameter Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA3', 'Comprehensive monitoring systems that track multiple vital parameters simultaneously.', 'demo-hospital');

-- Insert sub-categories for Rehabilitation Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ41', 'Physical Therapy Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Equipment used in physical therapy for rehabilitation after injury or surgery.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ42', 'Mobility Aids', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Devices that assist individuals with mobility impairments including wheelchairs, walkers, and canes.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ43', 'Orthopedic Braces & Supports', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Devices that provide support, alignment, or immobilization to limbs or joints.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ44', 'Exercise & Fitness Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Specialized equipment used for therapeutic exercise and rehabilitation.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ45', 'Electrotherapy Devices', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Equipment that uses electrical stimulation for pain management and muscle rehabilitation.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ46', 'Hydrotherapy Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Equipment that uses water for physical therapy and rehabilitation.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ47', 'Traction Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Devices used to apply a pulling force to parts of the body for therapeutic purposes.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ48', 'Prosthetics & Orthotics', '01HFPQ3Z5MXNVT9DAPZ3BWJA4', 'Artificial limbs and supportive devices designed to improve function and mobility.', 'demo-hospital');

-- Insert sub-categories for Dental Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ49', 'Dental Chairs & Delivery Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Specialized chairs and integrated systems for dental procedures and patient comfort.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ50', 'Dental Imaging Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'X-ray and imaging equipment specifically designed for dental diagnostics.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ51', 'Dental Handpieces', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'High-speed rotary instruments used for drilling and polishing in dental procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ52', 'Dental Lasers', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Laser devices specifically designed for dental treatments and procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ53', 'Dental CAD/CAM Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Computer-aided design and manufacturing systems for creating dental restorations.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ54', 'Dental Sterilization Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Equipment used for sterilizing dental instruments and materials.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ55', 'Dental Laboratory Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Equipment used in dental laboratories for creating prosthetics and restorations.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ56', 'Endodontic Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA5', 'Specialized tools and devices used for root canal procedures.', 'demo-hospital');

-- Insert sub-categories for Ophthalmic Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ57', 'Ophthalmic Diagnostic Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Equipment used for diagnosing eye conditions and diseases.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ58', 'Slit Lamps', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Microscopes with a light source used for eye examinations.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ59', 'Refraction Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Devices used to measure refractive errors and determine eyeglass prescriptions.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ60', 'Ophthalmic Surgical Microscopes', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Specialized microscopes used during eye surgeries.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ61', 'Phacoemulsification Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Equipment used for cataract surgery that uses ultrasonic energy to break up and remove cataracts.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ62', 'Ophthalmic Lasers', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Laser systems specifically designed for eye surgeries and treatments.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ63', 'Tonometers', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Devices used to measure intraocular pressure for glaucoma screening and monitoring.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ64', 'Visual Field Analyzers', '01HFPQ3Z5MXNVT9DAPZ3BWJA6', 'Equipment used to test and measure a patient''s entire scope of vision.', 'demo-hospital');

-- Insert sub-categories for Cardiology Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ65', 'Electrocardiographs (ECG/EKG)', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Devices that record the electrical activity of the heart.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ66', 'Cardiac Ultrasound (Echocardiography)', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Ultrasound equipment specifically designed for imaging the heart.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ67', 'Cardiac Catheterization Lab Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Specialized equipment used in cardiac catheterization procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ68', 'Holter Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Portable devices that record heart activity continuously for 24-48 hours.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ69', 'Stress Test Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Equipment used to measure heart function during physical activity.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ70', 'Defibrillators', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Devices that deliver electrical shocks to restore normal heart rhythm.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ71', 'Cardiac Output Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Devices that measure the volume of blood pumped by the heart.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ72', 'Pacemakers & Implantable Devices', '01HFPQ3Z5MXNVT9DAPZ3BWJA7', 'Electronic devices implanted to regulate heart rhythm.', 'demo-hospital');

-- Insert sub-categories for Orthopedic Devices
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ73', 'Orthopedic Implants', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Medical devices surgically implanted to replace or support damaged bone or joint structures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ74', 'Joint Replacement Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Prosthetic devices used to replace damaged joints, particularly knees and hips.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ75', 'Spine Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Implants and devices used in spine surgeries and treatments.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ76', 'Trauma Fixation Devices', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Plates, screws, and rods used to stabilize and heal bone fractures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ77', 'Orthopedic Power Tools', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Drills, saws, and other power tools specifically designed for orthopedic surgeries.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ78', 'Casting & Splinting Materials', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Materials used to immobilize and protect injured bones during healing.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ79', 'Arthroscopy Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Minimally invasive surgical equipment used to visualize, diagnose, and treat joint problems.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ80', 'Bone Growth Stimulators', '01HFPQ3Z5MXNVT9DAPZ3BWJA8', 'Devices that use electrical or ultrasound stimulation to promote bone healing.', 'demo-hospital');

-- Insert sub-categories for Hospital Furniture & Infrastructure
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ81', 'Hospital Beds', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Specialized beds designed for hospitalized patients with features for comfort, safety, and care accessibility.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ82', 'Stretchers & Trolleys', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Mobile platforms used to transport patients within healthcare facilities.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ83', 'Medical Carts & Workstations', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Mobile carts and workstations used by healthcare professionals for various clinical tasks.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ84', 'Examination Tables', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Tables specifically designed for patient examinations in clinical settings.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ85', 'Operating Tables', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Specialized tables used during surgical procedures with positioning capabilities.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ86', 'Medical Cabinets & Storage', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Storage solutions designed for medical supplies, medications, and equipment.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ87', 'Patient Room Furniture', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Specialized furniture designed for patient rooms in healthcare facilities.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ88', 'Modular Healthcare Facilities', '01HFPQ3Z5MXNVT9DAPZ3BWJA9', 'Prefabricated and modular solutions for healthcare facility construction and expansion.', 'demo-hospital');

-- Insert sub-categories for Biomedical Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ89', 'Laboratory Centrifuges', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Equipment that separates fluids, gases, or liquids based on density using centrifugal force.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ90', 'Incubators & Culture Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Devices used to grow and maintain microbiological cultures or cell cultures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ91', 'PCR & Molecular Diagnostic Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Equipment used for polymerase chain reaction and other molecular diagnostic techniques.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ92', 'Spectrophotometers', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Instruments that measure the intensity of light as a function of its wavelength.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ93', 'Microscopes', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Optical instruments used to view objects that are too small to be seen by the naked eye.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ94', 'Flow Cytometers', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Instruments used to analyze the physical and chemical characteristics of particles in a fluid.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ95', 'Biosafety Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Equipment designed to protect laboratory personnel and the environment from exposure to biohazards.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ96', 'Cryogenic Storage Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJ10', 'Equipment used for storing biological samples at extremely low temperatures.', 'demo-hospital');

-- Insert sub-categories for Respiratory Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BWJ97', 'Ventilators', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Machines designed to mechanically move breathable air into and out of the lungs.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ98', 'CPAP & BiPAP Machines', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices that deliver air pressure to keep airways open during sleep for patients with sleep apnea.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BWJ99', 'Oxygen Concentrators', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Medical devices that concentrate oxygen from ambient air by removing nitrogen.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW100', 'Nebulizers', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices that convert liquid medication into a fine mist for inhalation into the lungs.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW101', 'Spirometers', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices that measure lung function by measuring the volume and flow of air during breathing.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW102', 'Oxygen Therapy Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices used to deliver supplemental oxygen to patients who cannot get enough oxygen naturally.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW103', 'Suction Machines', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices used to remove substances like mucus or saliva from a person''s airway.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW104', 'Respiratory Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJ11', 'Devices that monitor respiratory rate, effort, and other breathing parameters.', 'demo-hospital');

-- Insert sub-categories for Anesthesia Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BW105', 'Anesthesia Machines', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Equipment used to deliver anesthetic agents during surgical procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW106', 'Anesthesia Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Devices that monitor patient vital signs and parameters during anesthesia.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW107', 'Anesthesia Ventilators', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Specialized ventilators designed for use during anesthesia.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW108', 'Anesthesia Gas Scavenging Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Systems designed to collect and remove waste anesthetic gases from the operating room.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW109', 'Vaporizers', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Devices that convert liquid anesthetic agents into vapor for inhalation.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW110', 'Regional Anesthesia Equipment', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Equipment used for administering regional anesthesia such as spinal or epidural blocks.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW111', 'Anesthesia Carts', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Mobile carts containing medications and equipment needed for anesthesia procedures.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW112', 'Airway Management Devices', '01HFPQ3Z5MXNVT9DAPZ3BWJ12', 'Equipment used to maintain or secure patient airways during anesthesia.', 'demo-hospital');

-- Insert sub-categories for Neonatal Equipment
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BW113', 'Infant Incubators', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Enclosed apparatus that provides a controlled environment for newborns, especially premature infants.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW114', 'Infant Radiant Warmers', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Open beds with overhead heating elements to maintain infant body temperature.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW115', 'Neonatal Ventilators', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Specialized ventilators designed for the unique respiratory needs of newborns.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW116', 'Phototherapy Units', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Devices that use light therapy to treat neonatal jaundice.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW117', 'Neonatal Monitors', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Specialized monitoring equipment designed for newborns and premature infants.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW118', 'Infant Scales', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Specialized scales designed for weighing newborns and infants.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW119', 'Infant Transport Systems', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Mobile incubators and equipment for safely transporting critically ill newborns.', 'demo-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW120', 'Bilirubin Meters', '01HFPQ3Z5MXNVT9DAPZ3BWJ13', 'Devices that measure bilirubin levels in newborns to detect and monitor jaundice.', 'demo-hospital');

-- Add main categories for city-hospital tenant (different tenant)
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BW121', 'Diagnostic Equipment', NULL, 'Equipment used for diagnosing medical conditions including imaging systems, laboratory analyzers, and diagnostic tools.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW122', 'Surgical Instruments & Devices', NULL, 'Tools and equipment used during surgical procedures including scalpels, forceps, retractors, and specialized surgical devices.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW123', 'Patient Monitoring Systems', NULL, 'Equipment used to continuously monitor patient vital signs and physiological parameters in various clinical settings.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW124', 'Rehabilitation Equipment', NULL, 'Devices and equipment used for physical therapy, rehabilitation, and recovery from injuries or medical conditions.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW125', 'Hospital Furniture & Infrastructure', NULL, 'Furniture and infrastructure components designed specifically for healthcare facilities including beds, stretchers, and cabinets.', 'city-hospital');

-- Add some sub-categories for city-hospital tenant
INSERT INTO categories (id, name, parent_id, description, tenant_id) VALUES
('01HFPQ3Z5MXNVT9DAPZ3BW126', 'X-ray Systems', '01HFPQ3Z5MXNVT9DAPZ3BW121', 'Equipment that uses X-ray radiation to create images of structures inside the body.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW127', 'Ultrasound Systems', '01HFPQ3Z5MXNVT9DAPZ3BW121', 'Imaging devices that use high-frequency sound waves to create images of structures inside the body.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW128', 'General Surgical Instruments', '01HFPQ3Z5MXNVT9DAPZ3BW122', 'Basic instruments used in various surgical procedures including scalpels, forceps, scissors, and clamps.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW129', 'Vital Signs Monitors', '01HFPQ3Z5MXNVT9DAPZ3BW123', 'Devices that measure and display patient vital signs including heart rate, blood pressure, temperature, and oxygen saturation.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW130', 'Physical Therapy Equipment', '01HFPQ3Z5MXNVT9DAPZ3BW124', 'Equipment used in physical therapy for rehabilitation after injury or surgery.', 'city-hospital'),
('01HFPQ3Z5MXNVT9DAPZ3BW131', 'Hospital Beds', '01HFPQ3Z5MXNVT9DAPZ3BW125', 'Specialized beds designed for hospitalized patients with features for comfort, safety, and care accessibility.', 'city-hospital');
