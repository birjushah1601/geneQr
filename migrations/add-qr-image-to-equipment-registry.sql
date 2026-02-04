-- Migration: Add QR code image storage to equipment_registry table
-- Date: 2026-02-05
-- Author: Droid AI
-- Purpose: Enable QR code image storage for PDF generation and image serving
-- Issue: equipment_registry table missing qr_code_image column causing 404 errors

BEGIN;

-- Add qr_code_image column to store binary PNG data
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS qr_code_image BYTEA;

-- Add qr_code_format column to track image format (default: png)
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS qr_code_format VARCHAR(10) DEFAULT 'png';

-- Add qr_code_generated_at timestamp to track when QR was generated
ALTER TABLE equipment_registry 
ADD COLUMN IF NOT EXISTS qr_code_generated_at TIMESTAMP;

-- Add index for quick lookup of equipment with QR codes
CREATE INDEX IF NOT EXISTS idx_equipment_registry_has_qr 
ON equipment_registry(id) 
WHERE qr_code_image IS NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN equipment_registry.qr_code_image IS 
    'Binary PNG image data for QR code, served via /api/v1/equipment/qr/image/{id}';
    
COMMENT ON COLUMN equipment_registry.qr_code_format IS 
    'Image format (png, svg, etc.)';
    
COMMENT ON COLUMN equipment_registry.qr_code_generated_at IS 
    'Timestamp when QR code image was last generated';

COMMIT;

-- Verification query (run separately)
-- SELECT column_name, data_type, character_maximum_length, column_default
-- FROM information_schema.columns
-- WHERE table_name = 'equipment_registry'
-- AND column_name IN ('qr_code_image', 'qr_code_format', 'qr_code_generated_at')
-- ORDER BY column_name;
