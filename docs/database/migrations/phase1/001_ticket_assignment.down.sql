-- ============================================================================
-- Rollback 001: Service Ticket Assignment Refactor
-- Author: Database Refactor Team
-- Date: 2025-11-16
-- ============================================================================

BEGIN;

-- Step 1: Copy current assignments back to service_tickets
UPDATE service_tickets st
SET 
    assigned_engineer_id = cta.engineer_id::text,
    assignment_tier = cta.assignment_tier,
    assigned_at = cta.assigned_at
FROM current_ticket_assignments cta
WHERE st.id = cta.ticket_id::text;

-- Step 2: Drop view
DROP VIEW IF EXISTS current_ticket_assignments;

-- Step 3: Remove new columns from engineer_assignments
ALTER TABLE engineer_assignments
DROP COLUMN IF EXISTS assignment_sequence,
DROP COLUMN IF EXISTS assignment_reason,
DROP COLUMN IF EXISTS rejection_reason,
DROP COLUMN IF EXISTS escalation_reason,
DROP COLUMN IF EXISTS completion_status,
DROP COLUMN IF EXISTS time_spent_hours,
DROP COLUMN IF EXISTS notes;

-- Step 4: Drop constraints
ALTER TABLE engineer_assignments
DROP CONSTRAINT IF EXISTS chk_completion_status;

-- Step 5: Drop new indexes
DROP INDEX IF EXISTS idx_ea_current_assignment;

-- Step 6: Remove comments
COMMENT ON COLUMN service_tickets.assigned_engineer_id IS NULL;
COMMENT ON COLUMN service_tickets.assigned_engineer_name IS NULL;
COMMENT ON COLUMN service_tickets.assignment_tier IS NULL;
COMMENT ON COLUMN service_tickets.assignment_tier_name IS NULL;
COMMENT ON TABLE engineer_assignments IS NULL;

COMMIT;

-- Verify rollback
SELECT 
    'service_tickets' as table_name,
    COUNT(*) as total_tickets,
    COUNT(assigned_engineer_id) as tickets_with_assignment
FROM service_tickets;
