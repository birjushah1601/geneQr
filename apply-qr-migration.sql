-- Quick Migration Script: Store QR Codes in Database
-- Run this to update your equipment table

BEGIN;

-- Add new columns for storing QR code in database
ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_image BYTEA;

ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png';

ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP;

-- Add comments
COMMENT ON COLUMN equipment.qr_code_image IS 'QR code image stored as binary data (PNG format)';
COMMENT ON COLUMN equipment.qr_code_format IS 'Format of stored QR code image (png, svg, jpg)';
COMMENT ON COLUMN equipment.qr_code_path IS 'DEPRECATED: Legacy filesystem path, use qr_code_image instead';

-- Create index for equipment with QR codes
CREATE INDEX IF NOT EXISTS idx_equipment_has_qr ON equipment(id) WHERE qr_code_image IS NOT NULL;

COMMIT;

-- Verify the changes
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'equipment' 
    AND column_name IN ('qr_code_image', 'qr_code_format', 'qr_code_generated_at', 'qr_code_path')
ORDER BY ordinal_position;
