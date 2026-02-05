-- ServQR: Align equipment-registry schema with application code
-- Purpose: Add columns expected by equipment repository and handlers
-- Safe: Uses IF NOT EXISTS; rerunnable

BEGIN;

ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS serial_number VARCHAR(200);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS equipment_id VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS equipment_name VARCHAR(500);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS manufacturer_name VARCHAR(200);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS model_number VARCHAR(200);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS category VARCHAR(200);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS customer_id VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS installation_location VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS installation_address JSONB DEFAULT '{}'::jsonb;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS installation_date DATE;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS contract_id VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS purchase_date DATE;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS purchase_price DECIMAL(15,2) DEFAULT 0;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS warranty_expiry DATE;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS amc_contract_id VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'operational';
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS last_service_date TIMESTAMP;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS next_service_date TIMESTAMP;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS service_count INTEGER DEFAULT 0;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS photos JSONB DEFAULT '[]'::jsonb;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS documents JSONB DEFAULT '[]'::jsonb;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_url VARCHAR(500);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS notes TEXT;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_image BYTEA;
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png';
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP;

COMMIT;

-- Verification
-- SELECT column_name FROM information_schema.columns WHERE table_name='equipment' ORDER BY ordinal_position;
