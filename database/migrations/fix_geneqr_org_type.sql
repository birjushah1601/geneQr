-- Update GENEQR organization type from 'manufacturer' to 'system'
-- This ensures admin@geneqr.com is recognized as a system administrator
-- rather than a manufacturer user

UPDATE organizations
SET org_type = 'system',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'GENEQR';

-- Verify the change
SELECT 
    name,
    org_type,
    status,
    updated_at
FROM organizations
WHERE name = 'GENEQR';
