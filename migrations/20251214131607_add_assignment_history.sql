-- Migration: Add ticket_assignment_history table
-- Purpose: Track all engineer assignments for audit and history

CREATE TABLE IF NOT EXISTS ticket_assignment_history (
    id VARCHAR(26) PRIMARY KEY,
    ticket_id VARCHAR(26) NOT NULL,
    engineer_id VARCHAR(255),
    engineer_name VARCHAR(255) NOT NULL,
    assigned_by VARCHAR(255) NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason VARCHAR(500),
    previous_engineer_id VARCHAR(255),
    previous_engineer_name VARCHAR(255),
    
    -- Foreign key
    CONSTRAINT fk_assignment_ticket 
        FOREIGN KEY (ticket_id) 
        REFERENCES service_tickets(id) 
        ON DELETE CASCADE
);

-- Create indexes separately
CREATE INDEX IF NOT EXISTS idx_assignment_ticket ON ticket_assignment_history(ticket_id);
CREATE INDEX IF NOT EXISTS idx_assignment_engineer ON ticket_assignment_history(engineer_id);
CREATE INDEX IF NOT EXISTS idx_assignment_date ON ticket_assignment_history(assigned_at);

-- Add comments
COMMENT ON TABLE ticket_assignment_history IS 'Tracks all engineer assignments for tickets including reassignments';
COMMENT ON COLUMN ticket_assignment_history.ticket_id IS 'Reference to the service ticket';
COMMENT ON COLUMN ticket_assignment_history.engineer_id IS 'ID of the engineer being assigned';
COMMENT ON COLUMN ticket_assignment_history.engineer_name IS 'Name of the engineer being assigned';
COMMENT ON COLUMN ticket_assignment_history.assigned_by IS 'User who made the assignment';
COMMENT ON COLUMN ticket_assignment_history.previous_engineer_id IS 'Previous engineer if this was a reassignment';
COMMENT ON COLUMN ticket_assignment_history.reason IS 'Optional reason for assignment/reassignment';
