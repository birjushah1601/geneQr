-- Migration: Store QR Codes in Database Instead of Filesystem
-- Date: 2025-10-05
-- Description: Changes QR code storage from filesystem to database (bytea column)

-- Add new column for storing QR code image data
ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_image BYTEA;

-- Add column for QR code format (png, svg, etc.)
ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png';

-- Add column for QR code generated timestamp
ALTER TABLE equipment 
ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP;

-- Keep qr_code_path for backward compatibility but make it nullable
ALTER TABLE equipment 
ALTER COLUMN qr_code_path DROP NOT NULL;

-- Add comment explaining the new approach
COMMENT ON COLUMN equipment.qr_code_image IS 'QR code image stored as binary data (PNG format)';
COMMENT ON COLUMN equipment.qr_code_format IS 'Format of stored QR code image (png, svg, jpg)';
COMMENT ON COLUMN equipment.qr_code_path IS 'DEPRECATED: Legacy filesystem path, use qr_code_image instead';

-- Create index for quick lookup of equipment with QR codes
CREATE INDEX IF NOT EXISTS idx_equipment_has_qr ON equipment(id) WHERE qr_code_image IS NOT NULL;
