-- Migration: Create equipment registry schema for field service management
-- This creates tables for tracking installed equipment and QR codes

-- Equipment registry table
CREATE TABLE IF NOT EXISTS equipment_registry (
    id VARCHAR(32) PRIMARY KEY,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    serial_number VARCHAR(255) UNIQUE NOT NULL,
    
    -- Equipment details
    equipment_id VARCHAR(32),  -- Link to catalog (optional)
    equipment_name VARCHAR(500) NOT NULL,
    manufacturer_name VARCHAR(255) NOT NULL,
    model_number VARCHAR(255),
    category VARCHAR(255),
    
    -- Installation details
    customer_id VARCHAR(32),
    customer_name VARCHAR(500) NOT NULL,
    installation_location TEXT,
    installation_address JSONB,
    installation_date DATE,
    
    -- Contract details
    contract_id VARCHAR(32),  -- Link to procurement contract
    purchase_date DATE,
    purchase_price DECIMAL(15,2),
    warranty_expiry DATE,
    amc_contract_id VARCHAR(32),
    
    -- Status and service
    status VARCHAR(50) NOT NULL DEFAULT 'operational',
    last_service_date DATE,
    next_service_date DATE,
    service_count INT NOT NULL DEFAULT 0,
    
    -- Technical details (stored as JSONB for flexibility)
    specifications JSONB DEFAULT '{}'::jsonb,
    photos JSONB DEFAULT '[]'::jsonb,  -- Array of photo URLs
    documents JSONB DEFAULT '[]'::jsonb,  -- Array of document URLs
    
    -- QR Code URL
    qr_code_url TEXT NOT NULL,
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL,
    
    CONSTRAINT equipment_status_check CHECK (status IN ('operational', 'down', 'under_maintenance', 'decommissioned'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_equipment_qr_code ON equipment_registry(qr_code);
CREATE INDEX IF NOT EXISTS idx_equipment_serial_number ON equipment_registry(serial_number);
CREATE INDEX IF NOT EXISTS idx_equipment_customer ON equipment_registry(customer_id);
CREATE INDEX IF NOT EXISTS idx_equipment_manufacturer ON equipment_registry(manufacturer_name);
CREATE INDEX IF NOT EXISTS idx_equipment_status ON equipment_registry(status);
CREATE INDEX IF NOT EXISTS idx_equipment_category ON equipment_registry(category);
CREATE INDEX IF NOT EXISTS idx_equipment_warranty_expiry ON equipment_registry(warranty_expiry);
CREATE INDEX IF NOT EXISTS idx_equipment_next_service ON equipment_registry(next_service_date);
CREATE INDEX IF NOT EXISTS idx_equipment_created_at ON equipment_registry(created_at DESC);

-- GIN indexes for JSONB queries
CREATE INDEX IF NOT EXISTS idx_equipment_specifications ON equipment_registry USING GIN (specifications);
CREATE INDEX IF NOT EXISTS idx_equipment_address ON equipment_registry USING GIN (installation_address);

-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_equipment_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER equipment_updated_at_trigger
    BEFORE UPDATE ON equipment_registry
    FOR EACH ROW
    EXECUTE FUNCTION update_equipment_updated_at();

-- Helper function to generate QR codes (just the unique identifier part)
CREATE OR REPLACE FUNCTION generate_qr_code()
RETURNS VARCHAR AS $$
DECLARE
    v_qr_code VARCHAR;
    v_exists BOOLEAN;
BEGIN
    LOOP
        -- Generate format: QR-YYYYMMDD-XXXXXX (random 6 digit)
        v_qr_code := 'QR-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || LPAD(FLOOR(RANDOM() * 1000000)::TEXT, 6, '0');
        
        -- Check if exists
        SELECT EXISTS(SELECT 1 FROM equipment_registry WHERE qr_code = v_qr_code) INTO v_exists;
        
        -- If doesn't exist, use it
        IF NOT v_exists THEN
            EXIT;
        END IF;
    END LOOP;
    
    RETURN v_qr_code;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get equipment statistics
CREATE OR REPLACE FUNCTION get_equipment_statistics()
RETURNS TABLE (
    total_equipment BIGINT,
    operational_count BIGINT,
    down_count BIGINT,
    under_maintenance_count BIGINT,
    under_warranty_count BIGINT,
    with_amc_count BIGINT,
    avg_service_count NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_equipment,
        COUNT(*) FILTER (WHERE status = 'operational')::BIGINT as operational_count,
        COUNT(*) FILTER (WHERE status = 'down')::BIGINT as down_count,
        COUNT(*) FILTER (WHERE status = 'under_maintenance')::BIGINT as under_maintenance_count,
        COUNT(*) FILTER (WHERE warranty_expiry IS NOT NULL AND warranty_expiry > NOW())::BIGINT as under_warranty_count,
        COUNT(*) FILTER (WHERE amc_contract_id IS NOT NULL AND amc_contract_id != '')::BIGINT as with_amc_count,
        COALESCE(AVG(service_count), 0)::NUMERIC as avg_service_count
    FROM equipment_registry;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get equipment needing service
CREATE OR REPLACE FUNCTION get_equipment_needing_service(p_days INT DEFAULT 7)
RETURNS TABLE (
    equipment_id VARCHAR,
    qr_code VARCHAR,
    serial_number VARCHAR,
    equipment_name VARCHAR,
    customer_name VARCHAR,
    next_service_date DATE,
    days_until_service INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        e.id,
        e.qr_code,
        e.serial_number,
        e.equipment_name,
        e.customer_name,
        e.next_service_date,
        EXTRACT(DAY FROM (e.next_service_date - NOW()::DATE))::INT as days_until_service
    FROM equipment_registry e
    WHERE e.next_service_date IS NOT NULL
      AND e.next_service_date BETWEEN NOW()::DATE AND (NOW()::DATE + INTERVAL '1 day' * p_days)
      AND e.status != 'decommissioned'
    ORDER BY e.next_service_date ASC;
END;
$$ LANGUAGE plpgsql;

-- Helper function to get equipment with expired warranty
CREATE OR REPLACE FUNCTION get_expired_warranty_equipment()
RETURNS TABLE (
    equipment_id VARCHAR,
    qr_code VARCHAR,
    serial_number VARCHAR,
    equipment_name VARCHAR,
    customer_name VARCHAR,
    warranty_expiry DATE,
    days_since_expiry INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        e.id,
        e.qr_code,
        e.serial_number,
        e.equipment_name,
        e.customer_name,
        e.warranty_expiry,
        EXTRACT(DAY FROM (NOW()::DATE - e.warranty_expiry))::INT as days_since_expiry
    FROM equipment_registry e
    WHERE e.warranty_expiry IS NOT NULL
      AND e.warranty_expiry < NOW()::DATE
      AND e.amc_contract_id IS NULL  -- No AMC coverage
      AND e.status != 'decommissioned'
    ORDER BY e.warranty_expiry DESC;
END;
$$ LANGUAGE plpgsql;

-- Sample CSV template comment (for documentation)
COMMENT ON TABLE equipment_registry IS 'Equipment registry for field service management. 
CSV Import Columns: serial_number, equipment_name, manufacturer_name, model_number, category, 
customer_name, customer_id, installation_location, installation_date (YYYY-MM-DD), 
purchase_date (YYYY-MM-DD), purchase_price, warranty_months, notes';
