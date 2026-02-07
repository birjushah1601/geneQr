-- Migration: Add customer_email to service_tickets
-- Description: Add customer_email column for email notifications
-- Date: 2026-02-07

BEGIN;

-- Add customer_email column to service_tickets
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS customer_email VARCHAR(255);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_service_tickets_customer_email 
ON service_tickets(customer_email);

-- Add comment
COMMENT ON COLUMN service_tickets.customer_email IS 'Customer email address for notifications';

COMMIT;

-- Verification
SELECT 
    'customer_email column added' as status,
    COUNT(*) as total_tickets,
    COUNT(customer_email) as tickets_with_email,
    COUNT(*) - COUNT(customer_email) as tickets_without_email
FROM service_tickets;
