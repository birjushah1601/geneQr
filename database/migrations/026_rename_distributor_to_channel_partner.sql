-- ============================================================================
-- Migration: Rename Distributor to Channel Partner and Dealer to Sub-Dealer
-- ============================================================================
-- Date: February 4, 2026
-- Description: Updates org_type values to align with ServQR rebranding
-- ============================================================================

BEGIN;

-- Step 1: Update existing data
UPDATE organizations 
SET org_type = 'channel_partner' 
WHERE org_type IN ('distributor', 'Channel Partner', 'Distributor');

UPDATE organizations 
SET org_type = 'sub_dealer' 
WHERE org_type IN ('dealer', 'Sub-Dealer', 'Dealer', 'sub_SUB_DEALER');

-- Step 2: Drop old constraint
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS chk_org_type;

-- Step 3: Add new constraint with proper values (including all existing types)
ALTER TABLE organizations ADD CONSTRAINT chk_org_type CHECK (org_type IN (
  'manufacturer',
  'channel_partner',
  'sub_dealer',
  'supplier',
  'hospital',
  'laboratory',
  'diagnostic_center',
  'imaging_center',
  'clinic',
  'service_provider',
  'logistics_partner',
  'insurance_provider',
  'government_body',
  'system_admin',
  'system',
  'other'
));

-- Step 4: Update any references in other tables (if role column exists)
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name='users' AND column_name='role') THEN
        UPDATE users 
        SET role = 'channel_partner_admin' 
        WHERE role IN ('distributor_admin', 'Distributor_admin');
        
        UPDATE users 
        SET role = 'sub_dealer_admin' 
        WHERE role IN ('dealer_admin', 'Dealer_admin');
    END IF;
END $$;

-- Step 5: Create indexes for new org_type values if not exists
CREATE INDEX IF NOT EXISTS idx_org_type_channel_partner ON organizations(org_type) WHERE org_type = 'channel_partner';
CREATE INDEX IF NOT EXISTS idx_org_type_sub_dealer ON organizations(org_type) WHERE org_type = 'sub_dealer';

COMMIT;

-- Verify the changes
SELECT org_type, COUNT(*) as count 
FROM organizations 
GROUP BY org_type 
ORDER BY org_type;
