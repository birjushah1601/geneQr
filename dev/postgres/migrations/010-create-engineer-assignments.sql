-- ============================================================================
-- Migration 010: Create Engineer Assignments Table
-- Purpose: Track complete assignment history for service tickets
-- Author: Database Refactor Team
-- Date: 2025-11-16
-- ============================================================================

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- ENGINEER ASSIGNMENTS TABLE
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE TABLE IF NOT EXISTS engineer_assignments (
    id VARCHAR(32) PRIMARY KEY,
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    engineer_id VARCHAR(32) NOT NULL REFERENCES engineers(id),
    equipment_id VARCHAR(32) NOT NULL,
    
    -- Sequence tracking
    assignment_sequence INT NOT NULL DEFAULT 1,
    assignment_tier INT NOT NULL DEFAULT 1,
    assignment_tier_name VARCHAR(100),
    assignment_reason TEXT,
    
    -- Workflow
    assignment_type VARCHAR(50) NOT NULL DEFAULT 'manual',
    status VARCHAR(50) NOT NULL DEFAULT 'assigned',
    assigned_by VARCHAR(32),
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    accepted_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    
    -- Execution
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    completion_status VARCHAR(50),
    escalation_reason TEXT,
    time_spent_hours DECIMAL(5,2) DEFAULT 0,
    
    -- Details
    diagnosis TEXT,
    actions_taken TEXT,
    parts_used JSONB DEFAULT '[]'::jsonb,
    
    -- Customer feedback
    customer_rating INT,
    customer_feedback TEXT,
    
    -- Metadata
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT ea_assignment_type_check CHECK (assignment_type IN ('auto', 'manual', 'escalation')),
    CONSTRAINT ea_status_check CHECK (status IN ('assigned', 'accepted', 'rejected', 'in_progress', 'completed', 'failed', 'escalated')),
    CONSTRAINT ea_completion_status_check CHECK (completion_status IN ('success', 'failed', 'escalated', 'parts_required', 'customer_unavailable')),
    CONSTRAINT ea_customer_rating_check CHECK (customer_rating >= 1 AND customer_rating <= 5)
);

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- INDEXES
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE INDEX IF NOT EXISTS idx_ea_ticket ON engineer_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ea_engineer ON engineer_assignments(engineer_id);
CREATE INDEX IF NOT EXISTS idx_ea_equipment ON engineer_assignments(equipment_id);
CREATE INDEX IF NOT EXISTS idx_ea_status ON engineer_assignments(status);
CREATE INDEX IF NOT EXISTS idx_ea_assigned_at ON engineer_assignments(assigned_at DESC);
CREATE INDEX IF NOT EXISTS idx_ea_sequence ON engineer_assignments(ticket_id, assignment_sequence);

-- Index for fast current assignment lookup
CREATE INDEX IF NOT EXISTS idx_ea_current_assignment 
ON engineer_assignments(ticket_id, assigned_at DESC) 
WHERE status NOT IN ('completed', 'rejected', 'failed', 'escalated');

-- GIN index for parts_used JSONB
CREATE INDEX IF NOT EXISTS idx_ea_parts_used ON engineer_assignments USING GIN (parts_used);

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- TRIGGERS
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE OR REPLACE FUNCTION update_engineer_assignment_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER engineer_assignment_updated_at_trigger
    BEFORE UPDATE ON engineer_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_engineer_assignment_updated_at();

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- VIEW FOR CURRENT ASSIGNMENTS
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

CREATE OR REPLACE VIEW current_ticket_assignments AS
SELECT DISTINCT ON (ticket_id)
    ticket_id,
    engineer_id,
    assignment_sequence,
    assignment_tier,
    assignment_tier_name,
    status,
    assigned_at,
    accepted_at,
    started_at,
    completed_at
FROM engineer_assignments
WHERE status NOT IN ('completed', 'rejected', 'failed', 'escalated')
ORDER BY ticket_id, assigned_at DESC;

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- BACKFILL EXISTING ASSIGNMENTS
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

-- Backfill existing assignment data from service_tickets table
INSERT INTO engineer_assignments (
    id,
    ticket_id,
    engineer_id,
    equipment_id,
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
    'ASSIGN-' || st.id AS id,
    st.id,
    st.assigned_engineer_id,
    st.equipment_id,
    1,  -- First assignment
    1,  -- Default tier
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
      WHERE ea.ticket_id = st.id 
        AND ea.engineer_id = st.assigned_engineer_id
  );

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- DEPRECATION COMMENTS
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

-- Mark old columns as deprecated (don't drop yet for safety during transition)
COMMENT ON COLUMN service_tickets.assigned_engineer_id IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON COLUMN service_tickets.assigned_engineer_name IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON COLUMN service_tickets.assigned_at IS 
'DEPRECATED: Use engineer_assignments table. Will be removed in next release.';

COMMENT ON TABLE engineer_assignments IS 
'Complete assignment history for service tickets. 
Tracks all engineers assigned to a ticket including escalations.
Current assignment is the latest row where status NOT IN (completed, rejected, failed, escalated).';

-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
-- VERIFICATION
-- ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SELECT 
    'engineer_assignments' as table_name,
    COUNT(*) as total_rows,
    COUNT(DISTINCT ticket_id) as unique_tickets,
    COUNT(DISTINCT engineer_id) as unique_engineers
FROM engineer_assignments;

SELECT 
    'Indexes' as component,
    COUNT(*) as count
FROM pg_indexes
WHERE tablename = 'engineer_assignments';

-- Show assignment statistics
SELECT 
    status,
    COUNT(*) as count,
    ROUND(AVG(time_spent_hours), 2) as avg_hours
FROM engineer_assignments
GROUP BY status
ORDER BY count DESC;
