-- ============================================================================
-- Migration 001: Service Ticket Assignment Refactor
-- Purpose: Remove inline assignment from service_tickets, enhance engineer_assignments
-- Author: Database Refactor Team
-- Date: 2025-11-16
-- ============================================================================

BEGIN;

-- Step 1: Enhance engineer_assignments table with new fields
ALTER TABLE engineer_assignments
ADD COLUMN IF NOT EXISTS assignment_sequence INT DEFAULT 1,
ADD COLUMN IF NOT EXISTS assignment_reason TEXT,
ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
ADD COLUMN IF NOT EXISTS escalation_reason TEXT,
ADD COLUMN IF NOT EXISTS completion_status TEXT,
ADD COLUMN IF NOT EXISTS time_spent_hours NUMERIC(5,2),
ADD COLUMN IF NOT EXISTS notes TEXT;

-- Add check constraint for completion_status
ALTER TABLE engineer_assignments
ADD CONSTRAINT chk_completion_status 
CHECK (completion_status IN ('success', 'failed', 'escalated', 'parts_required', 'customer_unavailable'));

-- Create index for fast current assignment lookup
CREATE INDEX IF NOT EXISTS idx_ea_current_assignment 
ON engineer_assignments(ticket_id, assigned_at DESC) 
WHERE status NOT IN ('completed', 'rejected', 'failed');

-- Step 2: Create view for current assignments (backward compatibility)
CREATE OR REPLACE VIEW current_ticket_assignments AS
SELECT DISTINCT ON (ticket_id)
    ticket_id,
    engineer_id,
    assignment_sequence,
    assignment_tier,
    status,
    assigned_at,
    completed_at
FROM engineer_assignments
WHERE status NOT IN ('completed', 'rejected', 'failed')
ORDER BY ticket_id, assigned_at DESC;

-- Step 3: Backfill existing assignment data
INSERT INTO engineer_assignments (
    id,
    ticket_id,
    engineer_id,
    assignment_sequence,
    assignment_tier,
    assignment_type,
    status,
    assigned_at,
    accepted_at,
    started_at,
    completed_at,
    created_at,
    updated_at
)
SELECT 
    gen_random_uuid(),
    st.id::uuid,
    st.assigned_engineer_id::uuid,
    1,  -- First assignment
    st.assignment_tier,
    'manual',
    CASE 
        WHEN st.status IN ('resolved', 'closed') THEN 'completed'
        WHEN st.status IN ('assigned', 'in_progress') THEN 'in_progress'
        ELSE 'assigned'
    END,
    st.assigned_at,
    st.assigned_at,  -- Assume immediate acceptance for old data
    st.started_at,
    st.resolved_at,
    st.created_at,
    st.updated_at
FROM service_tickets st
WHERE st.assigned_engineer_id IS NOT NULL
  AND st.assigned_engineer_id != ''
  -- Don't duplicate if already exists
  AND NOT EXISTS (
      SELECT 1 FROM engineer_assignments ea 
      WHERE ea.ticket_id::text = st.id 
        AND ea.engineer_id::text = st.assigned_engineer_id
  );

-- Step 4: Add comment for documentation
COMMENT ON TABLE engineer_assignments IS 
'Complete assignment history for service tickets. 
Tracks all engineers assigned to a ticket including escalations.
Current assignment is the latest row where status NOT IN (completed, rejected, failed).';

-- Step 5: Mark old columns as deprecated (don't drop yet for safety)
COMMENT ON COLUMN service_tickets.assigned_engineer_id IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON COLUMN service_tickets.assigned_engineer_name IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON COLUMN service_tickets.assignment_tier IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON COLUMN service_tickets.assignment_tier_name IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMIT;

-- Verify results
SELECT 
    'engineer_assignments' as table_name,
    COUNT(*) as total_rows,
    COUNT(DISTINCT ticket_id) as unique_tickets
FROM engineer_assignments;
