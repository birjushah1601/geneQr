-- Indian Medical Device Manufacturers Dataset
-- For ServQR Platform Catalog Module
-- This file populates the manufacturers table with Indian medical device companies

-- First ensure the schema exists
CREATE SCHEMA IF NOT EXISTS public;

-- Create the manufacturers table if it doesn't exist
CREATE TABLE IF NOT EXISTS manufacturers (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    headquarters VARCHAR(255) NOT NULL,
    website VARCHAR(255),
    specialization VARCHAR(255) NOT NULL,
    established INT,
    description TEXT,
    country VARCHAR(50) DEFAULT 'India',
    tenant_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Clear existing data if needed (commented out for safety)
-- TRUNCATE TABLE manufacturers;

-- Demo tenant for development purposes
INSERT INTO manufacturers (id, name, headquarters, website, specialization, established, description, tenant_id) VALUES
-- Diagnostic Equipment (imaging, lab equipment)
('01HFPQ2Z5MXNVT9DAPZ3BWJAHR', 'Trivitron Healthcare', 'Chennai, Tamil Nadu', 'https://www.trivitron.com', 'Diagnostic Equipment', 1997, 'Leading medical technology company offering comprehensive solutions in imaging, lab diagnostics, and critical care equipment across India and globally.', 'demo-hospital'),
('01HFPQ2Z5N7YWGXS9JVKM6F8QT', 'Transasia Bio-Medicals', 'Mumbai, Maharashtra', 'https://www.transasia.co.in', 'Diagnostic Equipment', 1979, 'India''s largest in-vitro diagnostic company with a comprehensive product portfolio in biochemistry, hematology, coagulation, immunology, and microbiology.', 'demo-hospital'),
('01HFPQ2Z5NWBCPXQ2RJVT3D7KF', 'BPL Medical Technologies', 'Bengaluru, Karnataka', 'https://www.bplmedicaltechnologies.com', 'Diagnostic Equipment', 1967, 'Pioneer in medical equipment manufacturing with focus on patient monitoring systems, ECG machines, and defibrillators.', 'demo-hospital'),
('01HFPQ2Z5P4MXVGZ3QNBT5F8HR', 'Agappe Diagnostics', 'Kochi, Kerala', 'https://www.agappe.com', 'Diagnostic Equipment', 1994, 'Specializes in manufacturing biochemistry reagents, analyzers, and ELISA kits with presence in over 65 countries.', 'demo-hospital'),
('01HFPQ2Z5PQJKWT4XZBS7M9HGF', 'J. Mitra & Co.', 'New Delhi, Delhi', 'https://www.jmitra.co.in', 'Diagnostic Equipment', 1969, 'Leading biotechnology company specializing in rapid test kits, ELISA kits, and blood grouping reagents.', 'demo-hospital'),
('01HFPQ2Z5Q8VFNP6XZRT4K7MHG', 'Meril Diagnostics', 'Vapi, Gujarat', 'https://www.merillife.com', 'Diagnostic Equipment', 2006, 'Manufactures high-quality diagnostic reagents, instruments, and rapid test kits for infectious diseases and pregnancy testing.', 'demo-hospital'),
('01HFPQ2Z5QXTHGV7PZBS3F9KMN', 'Skanray Technologies', 'Mysuru, Karnataka', 'https://www.skanray.com', 'Diagnostic Equipment', 2007, 'Designs and manufactures advanced medical imaging and life-supporting equipment including X-ray systems and critical care devices.', 'demo-hospital'),

-- Surgical Instruments & Devices
('01HFPQ2Z5RBWKPN8QZXT5G7JHF', 'Hindustan Syringes & Medical Devices', 'Faridabad, Haryana', 'https://www.hmdhealthcare.com', 'Surgical Instruments & Devices', 1957, 'World''s largest manufacturer of auto-disable syringes and a leading producer of disposable medical devices.', 'demo-hospital'),
('01HFPQ2Z5RWTHGV9PZBS3F8KMN', 'Poly Medicure Limited', 'Faridabad, Haryana', 'https://www.polymedicure.com', 'Surgical Instruments & Devices', 1995, 'Manufactures and exports medical devices in over 100 countries, specializing in infusion therapy, surgery, anesthesia, and urology products.', 'demo-hospital'),
('01HFPQ2Z5SJKWT0XZBS7M9HGF', 'Sutures India', 'Bengaluru, Karnataka', 'https://www.suturesindia.com', 'Surgical Instruments & Devices', 1992, 'Leading manufacturer of surgical sutures, mesh, skin staplers, and other wound closure products.', 'demo-hospital'),
('01HFPQ2Z5T4MXVG1QNBT5F8HR', 'Healthium Medtech', 'Bengaluru, Karnataka', 'https://www.healthiummedtech.com', 'Surgical Instruments & Devices', 1992, 'Develops, manufactures and markets surgical devices focusing on wound closure, minimally invasive surgeries, and infection prevention.', 'demo-hospital'),
('01HFPQ2Z5TQJKW2XZBS7M9HGF', 'Iscon Surgicals', 'Jodhpur, Rajasthan', 'https://www.isconsurgical.com', 'Surgical Instruments & Devices', 1975, 'Manufactures high-quality surgical and dental instruments exported to over 40 countries worldwide.', 'demo-hospital'),
('01HFPQ2Z5V8VFN3XZRT4K7MHG', 'Kalelker Surgical Industries', 'Mumbai, Maharashtra', 'https://www.kalelker.com', 'Surgical Instruments & Devices', 1930, 'One of India''s oldest surgical instrument manufacturers specializing in general surgery, gynecology, and ENT instruments.', 'demo-hospital'),

-- Patient Monitoring Systems
('01HFPQ2Z5VXTH4VPZBS3F9KMN', 'BPL Medical Technologies', 'Bengaluru, Karnataka', 'https://www.bplmedicaltechnologies.com', 'Patient Monitoring Systems', 1967, 'Leading manufacturer of patient monitoring systems, ECG machines, and defibrillators with nationwide service network.', 'demo-hospital'),
('01HFPQ2Z5WBWK5NQZXT5G7JHF', 'Larsen & Toubro Medical Equipment', 'Mumbai, Maharashtra', 'https://www.lntmedical.com', 'Patient Monitoring Systems', 1995, 'Division of L&T focusing on advanced patient monitoring systems and critical care equipment.', 'demo-hospital'),
('01HFPQ2Z5XJKW6XZBS7M9HGF', 'Opto Circuits India', 'Bengaluru, Karnataka', 'https://www.optoindia.com', 'Patient Monitoring Systems', 1992, 'Global medical technology group focusing on interventional devices, monitoring, diagnostics, and emergency cardiac care.', 'demo-hospital'),
('01HFPQ2Z5Y4MX7GQNBT5F8HR', 'Mindray Medical India', 'Gurugram, Haryana', 'https://www.mindray.com/in', 'Patient Monitoring Systems', 2008, 'Indian subsidiary of global medical device company specializing in patient monitoring, diagnostic, and ultrasound systems.', 'demo-hospital'),
('01HFPQ2Z5YQJK8TZBS7M9HGF', 'Schiller Healthcare India', 'Mumbai, Maharashtra', 'https://www.schillerindia.com', 'Patient Monitoring Systems', 1989, 'Manufactures cardiovascular diagnostic systems, patient monitors, and external defibrillators.', 'demo-hospital'),

-- Rehabilitation Equipment
('01HFPQ2Z5Z8VF9PXZRT4K7MHG', 'Vissco Rehabilitation Aids', 'Mumbai, Maharashtra', 'https://www.vissco.com', 'Rehabilitation Equipment', 1963, 'Leading manufacturer of orthopedic and rehabilitation products including supports, braces, and mobility aids.', 'demo-hospital'),
('01HFPQ2Z60XTH0VPZBS3F9KMN', 'Tynor Orthotics', 'Mohali, Punjab', 'https://www.tynorindia.com', 'Rehabilitation Equipment', 1991, 'Manufactures orthopedic and fracture aids, body supports, and rehabilitation products.', 'demo-hospital'),
('01HFPQ2Z61BWK1NQZXT5G7JHF', 'Physiomed Devices', 'New Delhi, Delhi', 'https://www.physiomedindia.com', 'Rehabilitation Equipment', 1982, 'Manufactures physiotherapy, rehabilitation, and electrotherapy equipment for hospitals and clinics.', 'demo-hospital'),
('01HFPQ2Z61WTH2VPZBS3F8KMN', 'Asco Medicare', 'Kolkata, West Bengal', 'https://www.ascomedicare.com', 'Rehabilitation Equipment', 1984, 'Specializes in manufacturing mobility aids, hospital furniture, and rehabilitation equipment.', 'demo-hospital'),
('01HFPQ2Z62JKW3TZBS7M9HGF', 'Mobility India', 'Bengaluru, Karnataka', 'https://www.mobility-india.org', 'Rehabilitation Equipment', 1994, 'Focuses on prosthetics, orthotics, and rehabilitation aids with emphasis on accessibility and affordability.', 'demo-hospital'),

-- Dental Equipment
('01HFPQ2Z634MX4GQNBT5F8HR', 'Confident Dental Equipment', 'Bengaluru, Karnataka', 'https://www.confident.in', 'Dental Equipment', 1959, 'Manufactures comprehensive range of dental chairs, equipment, and accessories exported to over 100 countries.', 'demo-hospital'),
('01HFPQ2Z63QJK5TZBS7M9HGF', 'Dentech Dental Care', 'New Delhi, Delhi', 'https://www.dentechdentalcare.com', 'Dental Equipment', 1996, 'Manufactures dental chairs, delivery systems, and specialized dental equipment for clinics and hospitals.', 'demo-hospital'),
('01HFPQ2Z648VF6PXZRT4K7MHG', 'Gnatus India', 'Mumbai, Maharashtra', 'https://www.gnatusindia.com', 'Dental Equipment', 1998, 'Indian arm of global dental equipment manufacturer specializing in dental chairs and imaging systems.', 'demo-hospital'),
('01HFPQ2Z64XTH7VPZBS3F9KMN', 'Dentsply Sirona India', 'Gurugram, Haryana', 'https://www.dentsplysirona.com/en-in', 'Dental Equipment', 1995, 'Indian subsidiary offering comprehensive dental solutions including CAD/CAM systems, imaging equipment, and dental materials.', 'demo-hospital'),
('01HFPQ2Z65BWK8NQZXT5G7JHF', 'Septodont Healthcare India', 'Mumbai, Maharashtra', 'https://www.septodont.in', 'Dental Equipment', 1989, 'Specializes in dental anesthetics, endodontics, and dental surgical products.', 'demo-hospital'),

-- Ophthalmic Equipment
('01HFPQ2Z65WTH9VPZBS3F8KMN', 'Appasamy Associates', 'Chennai, Tamil Nadu', 'https://www.appasamy.com', 'Ophthalmic Equipment', 1978, 'Leading manufacturer of ophthalmic equipment including surgical microscopes, slit lamps, and phaco systems.', 'demo-hospital'),
('01HFPQ2Z66JKWATZBS7M9HGF', 'Suraj Eye Institute', 'Nagpur, Maharashtra', 'https://www.surajeyeinstitute.org', 'Ophthalmic Equipment', 1995, 'Develops innovative, affordable ophthalmic devices and instruments for cataract and retinal surgeries.', 'demo-hospital'),
('01HFPQ2Z674MXBGQNBT5F8HR', 'Aurolab', 'Madurai, Tamil Nadu', 'https://www.aurolab.com', 'Ophthalmic Equipment', 1992, 'Manufactures affordable intraocular lenses, suture needles, pharmaceutical products, and equipment for eye care.', 'demo-hospital'),
('01HFPQ2Z67QJKCTZBS7M9HGF', 'Care Group', 'Vadodara, Gujarat', 'https://www.caregroup.in', 'Ophthalmic Equipment', 1987, 'Manufactures ophthalmic surgical instruments, diagnostic equipment, and vision care products.', 'demo-hospital'),
('01HFPQ2Z688VFDPXZRT4K7MHG', 'Ophthalmic Instruments & Equipments', 'Ambala, Haryana', 'https://www.oieindia.com', 'Ophthalmic Equipment', 1962, 'Manufactures comprehensive range of ophthalmic instruments and diagnostic equipment.', 'demo-hospital'),

-- Cardiology Equipment
('01HFPQ2Z68XTHEVPZBS3F9KMN', 'BPL Medical Technologies', 'Bengaluru, Karnataka', 'https://www.bplmedicaltechnologies.com', 'Cardiology Equipment', 1967, 'Manufactures ECG machines, cardiac monitors, and defibrillators with nationwide service network.', 'demo-hospital'),
('01HFPQ2Z69BWKFNQZXT5G7JHF', 'Schiller Healthcare India', 'Mumbai, Maharashtra', 'https://www.schillerindia.com', 'Cardiology Equipment', 1989, 'Specializes in cardiovascular diagnostic systems, stress test systems, and Holter monitors.', 'demo-hospital'),
('01HFPQ2Z69WTHGVPZBS3F8KMN', 'Opto Circuits India', 'Bengaluru, Karnataka', 'https://www.optoindia.com', 'Cardiology Equipment', 1992, 'Manufactures cardiac diagnostic equipment, patient monitoring systems, and interventional products.', 'demo-hospital'),
('01HFPQ2Z6AJKWHTZBS7M9HGF', 'Meril Life Sciences', 'Vapi, Gujarat', 'https://www.merillife.com', 'Cardiology Equipment', 2006, 'Develops and manufactures cardiovascular devices including stents, balloons, and cardiac diagnostic equipment.', 'demo-hospital'),
('01HFPQ2Z6B4MXIGQNBT5F8HR', 'Sahajanand Medical Technologies', 'Surat, Gujarat', 'https://www.smtpl.com', 'Cardiology Equipment', 1998, 'Specializes in cardiovascular devices with focus on coronary stents and balloon catheters.', 'demo-hospital'),

-- Orthopedic Devices
('01HFPQ2Z6BQJKJTZBS7M9HGF', 'Sharma Orthopedic', 'Ahmedabad, Gujarat', 'https://www.sharmaortho.com', 'Orthopedic Devices', 1965, 'Manufactures orthopedic implants, trauma products, and surgical instruments for bone and joint surgeries.', 'demo-hospital'),
('01HFPQ2Z6C8VFKPXZRT4K7MHG', 'Auxein Medical', 'Sonipat, Haryana', 'https://www.auxein.com', 'Orthopedic Devices', 2005, 'Manufactures orthopedic implants and instruments including trauma, spine, and joint replacement systems.', 'demo-hospital'),
('01HFPQ2Z6CXTHLVPZBS3F9KMN', 'Narang Medical Limited', 'New Delhi, Delhi', 'https://www.narang.com', 'Orthopedic Devices', 1950, 'Manufactures orthopedic implants, instruments, and hospital furniture with exports to over 90 countries.', 'demo-hospital'),
('01HFPQ2Z6DBWKMNQZXT5G7JHF', 'Meril Life Sciences', 'Vapi, Gujarat', 'https://www.merillife.com', 'Orthopedic Devices', 2006, 'Develops orthopedic implants and instruments for trauma, spine, and joint reconstruction.', 'demo-hospital'),
('01HFPQ2Z6DWTHNVPZBS3F8KMN', 'Yogeshwar Implants', 'Thane, Maharashtra', 'https://www.yogeshwarimplants.com', 'Orthopedic Devices', 1992, 'Manufactures orthopedic implants including bone plates, screws, and intramedullary nails.', 'demo-hospital'),

-- Hospital Furniture & Infrastructure
('01HFPQ2Z6EJKWOTZBS7M9HGF', 'Janak Healthcare', 'Mumbai, Maharashtra', 'https://www.janakhealthcare.com', 'Hospital Furniture & Infrastructure', 1951, 'Manufactures hospital beds, OT tables, examination tables, and other medical furniture.', 'demo-hospital'),
('01HFPQ2Z6F4MXPGQNBT5F8HR', 'Midmark India', 'Mumbai, Maharashtra', 'https://www.midmark.in', 'Hospital Furniture & Infrastructure', 2007, 'Manufactures examination tables, procedure chairs, and medical cabinetry for healthcare facilities.', 'demo-hospital'),
('01HFPQ2Z6FQJKQTZBS7M9HGF', 'Narang Medical Limited', 'New Delhi, Delhi', 'https://www.narang.com', 'Hospital Furniture & Infrastructure', 1950, 'Comprehensive manufacturer of hospital furniture, equipment, and surgical instruments.', 'demo-hospital'),
('01HFPQ2Z6G8VFRPXZRT4K7MHG', 'Godrej Interio Healthcare', 'Mumbai, Maharashtra', 'https://www.godrejinterio.com/healthcare', 'Hospital Furniture & Infrastructure', 1971, 'Division of Godrej specializing in hospital furniture, modular OT solutions, and healthcare infrastructure.', 'demo-hospital'),
('01HFPQ2Z6GXTHSVPZBS3F9KMN', 'Paramount Bed India', 'Bengaluru, Karnataka', 'https://www.paramount.in', 'Hospital Furniture & Infrastructure', 2000, 'Manufactures advanced hospital beds, stretchers, and medical furniture with focus on patient comfort and safety.', 'demo-hospital'),

-- Biomedical Equipment
('01HFPQ2Z6HBWKTNQZXT5G7JHF', 'Skanray Technologies', 'Mysuru, Karnataka', 'https://www.skanray.com', 'Biomedical Equipment', 2007, 'Designs and manufactures critical care equipment including ventilators, anesthesia systems, and patient monitors.', 'demo-hospital'),
('01HFPQ2Z6HWTHUVPZBS3F8KMN', 'Advanced Micronic Devices', 'Bengaluru, Karnataka', 'https://www.micronicindia.com', 'Biomedical Equipment', 1989, 'Manufactures laboratory equipment, analytical instruments, and biomedical devices.', 'demo-hospital'),
('01HFPQ2Z6IJKWVTZBS7M9HGF', 'Nasan Medical Electronics', 'Chennai, Tamil Nadu', 'https://www.nasanmedical.com', 'Biomedical Equipment', 1984, 'Manufactures electrosurgical units, diathermy machines, and other biomedical equipment.', 'demo-hospital'),
('01HFPQ2Z6J4MXWGQNBT5F8HR', 'Maestros Electronics & Telecommunications', 'Navi Mumbai, Maharashtra', 'https://www.maestros.net', 'Biomedical Equipment', 1973, 'Designs and manufactures medical electronics, telemedicine solutions, and biomedical equipment.', 'demo-hospital'),
('01HFPQ2Z6JQJKXTZBS7M9HGF', 'Opto Eurocor Healthcare', 'Bengaluru, Karnataka', 'https://www.optoeurocor.com', 'Biomedical Equipment', 2005, 'Manufactures cardiovascular devices, biomedical sensors, and monitoring equipment.', 'demo-hospital'),

-- Additional Mixed Specialization Companies
('01HFPQ2Z6K8VFYPXZRT4K7MHG', 'Wipro GE Healthcare', 'Bengaluru, Karnataka', 'https://www.wipro-ge.com', 'Diagnostic Equipment', 1990, 'Joint venture offering comprehensive medical technology solutions including imaging, diagnostics, and monitoring systems.', 'demo-hospital'),
('01HFPQ2Z6KXTHZVPZBS3F9KMN', 'Philips India Healthcare', 'Gurugram, Haryana', 'https://www.philips.co.in/healthcare', 'Patient Monitoring Systems', 1996, 'Indian arm offering advanced healthcare technologies in imaging, monitoring, and informatics.', 'demo-hospital'),
('01HFPQ2Z6LBWK0NQZXT5G7JHF', 'Phoenix Medical Systems', 'Chennai, Tamil Nadu', 'https://www.phoenixmedicalsystems.com', 'Neonatal Equipment', 1989, 'Specializes in maternal and neonatal care equipment including incubators, phototherapy units, and radiant warmers.', 'demo-hospital'),
('01HFPQ2Z6LWTH1VPZBS3F8KMN', 'Nidek Medical India', 'Bengaluru, Karnataka', 'https://www.nidekmedical.in', 'Respiratory Equipment', 1995, 'Manufactures oxygen concentrators, ventilators, and respiratory care equipment.', 'demo-hospital'),
('01HFPQ2Z6MJKW2TZBS7M9HGF', 'Oxymed India', 'Chennai, Tamil Nadu', 'https://www.oxymedindia.com', 'Respiratory Equipment', 1992, 'Manufactures oxygen delivery systems, concentrators, and respiratory care products.', 'demo-hospital'),
('01HFPQ2Z6N4MX3GQNBT5F8HR', 'Eastern Medikit', 'Gurugram, Haryana', 'https://www.easternmedikit.com', 'Surgical Instruments & Devices', 1991, 'Manufactures IV cannulas, blood collection tubes, and other medical disposables.', 'demo-hospital'),
('01HFPQ2Z6NQJK4TZBS7M9HGF', 'Advin Health Care', 'Bengaluru, Karnataka', 'https://www.advinhealth.com', 'Surgical Instruments & Devices', 1998, 'Manufactures surgical disposables, wound care products, and infection control solutions.', 'demo-hospital'),
('01HFPQ2Z6P8VF5PXZRT4K7MHG', 'Allengers Medical Systems', 'Mohali, Punjab', 'https://www.allengers.com', 'Diagnostic Equipment', 1987, 'Manufactures X-ray machines, C-arms, CT scanners, and other imaging equipment.', 'demo-hospital'),
('01HFPQ2Z6PXTH6VPZBS3F9KMN', 'IPA Medical', 'Mumbai, Maharashtra', 'https://www.ipamedical.com', 'Anesthesia Equipment', 1978, 'Manufactures anesthesia machines, ventilators, and critical care equipment.', 'demo-hospital'),
('01HFPQ2Z6QBWK7NQZXT5G7JHF', 'Clarity Medical Systems', 'Mohali, Punjab', 'https://www.claritymedical.in', 'Ophthalmic Equipment', 2003, 'Specializes in ophthalmic diagnostic and surgical equipment with focus on retinal imaging.', 'demo-hospital');

-- Add more tenant-specific manufacturers
INSERT INTO manufacturers (id, name, headquarters, website, specialization, established, description, tenant_id) VALUES
('01HFPQ2Z6QWTH8VPZBS3F8KMN', 'Trivitron Healthcare', 'Chennai, Tamil Nadu', 'https://www.trivitron.com', 'Diagnostic Equipment', 1997, 'Leading medical technology company offering comprehensive solutions in imaging, lab diagnostics, and critical care equipment across India and globally.', 'city-hospital'),
('01HFPQ2Z6RJKW9TZBS7M9HGF', 'BPL Medical Technologies', 'Bengaluru, Karnataka', 'https://www.bplmedicaltechnologies.com', 'Patient Monitoring Systems', 1967, 'Pioneer in medical equipment manufacturing with focus on patient monitoring systems, ECG machines, and defibrillators.', 'city-hospital'),
('01HFPQ2Z6S4MXAGQNBT5F8HR', 'Hindustan Syringes & Medical Devices', 'Faridabad, Haryana', 'https://www.hmdhealthcare.com', 'Surgical Instruments & Devices', 1957, 'World''s largest manufacturer of auto-disable syringes and a leading producer of disposable medical devices.', 'city-hospital'),
('01HFPQ2Z6SQJKBTZBS7M9HGF', 'Appasamy Associates', 'Chennai, Tamil Nadu', 'https://www.appasamy.com', 'Ophthalmic Equipment', 1978, 'Leading manufacturer of ophthalmic equipment including surgical microscopes, slit lamps, and phaco systems.', 'city-hospital'),
('01HFPQ2Z6T8VFCPXZRT4K7MHG', 'Janak Healthcare', 'Mumbai, Maharashtra', 'https://www.janakhealthcare.com', 'Hospital Furniture & Infrastructure', 1951, 'Manufactures hospital beds, OT tables, examination tables, and other medical furniture.', 'city-hospital');
