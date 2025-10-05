-- ============================================================================
-- FIELD ENGINEERS DATABASE SCHEMA
-- ============================================================================
-- Engineers table for field technicians who service equipment
-- ============================================================================

-- Drop existing tables
DROP TABLE IF EXISTS engineer_certifications CASCADE;
DROP TABLE IF EXISTS engineers CASCADE;

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- ENGINEERS TABLE
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE TABLE engineers (
    id VARCHAR(32) PRIMARY KEY,
    
    -- Personal Info
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    whatsapp VARCHAR(20),
    email VARCHAR(255) NOT NULL,
    
    -- Location
    location VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'India',
    pincode VARCHAR(10),
    
    -- Skills & Qualifications
    specializations TEXT[] NOT NULL DEFAULT '{}',
    certifications JSONB DEFAULT '[]',
    experience_years INTEGER DEFAULT 0,
    qualification VARCHAR(255),
    
    -- Assignment
    manufacturer_id VARCHAR(255),
    manufacturer_name VARCHAR(255),
    employee_id VARCHAR(50),
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    availability VARCHAR(50) NOT NULL DEFAULT 'available',
    
    -- Performance Metrics
    rating DECIMAL(3, 2),
    total_tickets INTEGER DEFAULT 0,
    completed_tickets INTEGER DEFAULT 0,
    in_progress_tickets INTEGER DEFAULT 0,
    avg_resolution_time DECIMAL(10, 2), -- hours
    customer_satisfaction_score DECIMAL(3, 2),
    
    -- Documents
    photo_url TEXT,
    documents JSONB DEFAULT '[]',
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    
    -- Constraints
    CONSTRAINT engineers_status_check CHECK (status IN ('active', 'inactive', 'on_leave', 'terminated')),
    CONSTRAINT engineers_availability_check CHECK (availability IN ('available', 'on_job', 'off_duty', 'on_leave')),
    CONSTRAINT engineers_rating_check CHECK (rating >= 0 AND rating <= 5),
    CONSTRAINT engineers_satisfaction_check CHECK (customer_satisfaction_score >= 0 AND customer_satisfaction_score <= 5)
);

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- INDEXES
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE INDEX idx_engineers_phone ON engineers(phone);
CREATE INDEX idx_engineers_email ON engineers(email);
CREATE INDEX idx_engineers_location ON engineers(location);
CREATE INDEX idx_engineers_manufacturer ON engineers(manufacturer_id);
CREATE INDEX idx_engineers_status ON engineers(status);
CREATE INDEX idx_engineers_availability ON engineers(availability);
CREATE INDEX idx_engineers_specializations ON engineers USING GIN(specializations);
CREATE INDEX idx_engineers_created_at ON engineers(created_at DESC);
CREATE INDEX idx_engineers_rating ON engineers(rating DESC);

-- Geo-spatial index for location-based queries
CREATE INDEX idx_engineers_coordinates ON engineers USING gist(
    ll_to_earth(latitude, longitude)
) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- TRIGGER FOR UPDATED_AT
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE OR REPLACE FUNCTION update_engineer_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER engineer_updated_at_trigger
    BEFORE UPDATE ON engineers
    FOR EACH ROW
    EXECUTE FUNCTION update_engineer_updated_at();

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- SAMPLE DATA
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INSERT INTO engineers (
    id, name, phone, whatsapp, email, location, city, state,
    specializations, experience_years, status, availability,
    rating, total_tickets, completed_tickets, created_by
) VALUES
(
    'ENG-001',
    'Raj Kumar Sharma',
    '+91-9876543210',
    '+91-9876543210',
    'raj.sharma@example.com',
    'Delhi NCR',
    'New Delhi',
    'Delhi',
    ARRAY['MRI Scanner', 'CT Scanner', 'X-Ray'],
    8,
    'active',
    'available',
    4.7,
    145,
    142,
    'system'
),
(
    'ENG-002',
    'Priya Shah',
    '+91-9876543211',
    '+91-9876543211',
    'priya.shah@example.com',
    'Mumbai',
    'Mumbai',
    'Maharashtra',
    ARRAY['Ultrasound', 'ECG', 'Patient Monitoring'],
    5,
    'active',
    'on_job',
    4.9,
    98,
    96,
    'system'
),
(
    'ENG-003',
    'Amit Patel',
    '+91-9876543212',
    '+91-9876543212',
    'amit.patel@example.com',
    'Bangalore',
    'Bangalore',
    'Karnataka',
    ARRAY['ICU Ventilator', 'Anesthesia Machine', 'Critical Care'],
    10,
    'active',
    'available',
    4.8,
    203,
    198,
    'system'
),
(
    'ENG-004',
    'Sneha Reddy',
    '+91-9876543213',
    '+91-9876543213',
    'sneha.reddy@example.com',
    'Hyderabad',
    'Hyderabad',
    'Telangana',
    ARRAY['Laboratory Equipment', 'Diagnostic Tools'],
    6,
    'active',
    'available',
    4.6,
    87,
    84,
    'system'
),
(
    'ENG-005',
    'Vikram Singh',
    '+91-9876543214',
    '+91-9876543214',
    'vikram.singh@example.com',
    'Pune',
    'Pune',
    'Maharashtra',
    ARRAY['MRI Scanner', 'CT Scanner', 'PET Scanner'],
    12,
    'active',
    'off_duty',
    4.9,
    267,
    263,
    'system'
);

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- VERIFICATION
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SELECT 
    'Engineers Table Created' as status,
    COUNT(*) as sample_data_count
FROM engineers;

SELECT 
    'Indexes' as component,
    COUNT(*) as count
FROM pg_indexes
WHERE tablename = 'engineers';

-- Show engineer summary by manufacturer
SELECT 
    manufacturer_name,
    status,
    availability,
    COUNT(*) as count,
    ROUND(AVG(rating), 2) as avg_rating
FROM engineers
GROUP BY manufacturer_name, status, availability
ORDER BY manufacturer_name, status, availability;

-- Verify multi-tenant isolation
SELECT 
    'Multi-tenant Check' as test,
    COUNT(DISTINCT manufacturer_id) as unique_manufacturers,
    COUNT(*) as total_engineers
FROM engineers;
