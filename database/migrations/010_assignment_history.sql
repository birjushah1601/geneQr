-- Migration: Assignment History Table
-- Description: Stores AI-powered engineer assignment recommendations
-- Dependencies: service_tickets, users tables
-- Phase: 2C - AI Services Layer

-- Create assignment_history table
CREATE TABLE IF NOT EXISTS assignment_history (
    -- Primary key
    id BIGSERIAL PRIMARY KEY,
    
    -- Assignment identification
    request_id VARCHAR(100) UNIQUE NOT NULL,
    ticket_id BIGINT NOT NULL REFERENCES service_tickets(ticket_id) ON DELETE CASCADE,
    
    -- Recommendations data (full JSON)
    recommendations JSONB NOT NULL,
    metadata JSONB NOT NULL,
    
    -- Selection tracking
    selected_engineer_id BIGINT REFERENCES users(user_id),
    selection_reason TEXT,
    selection_time TIMESTAMP WITH TIME ZONE,
    
    -- Feedback on recommendation quality
    was_successful BOOLEAN,
    actual_resolution_time INTERVAL,
    feedback_notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_assignment_history_ticket ON assignment_history(ticket_id);
CREATE INDEX idx_assignment_history_created ON assignment_history(created_at DESC);
CREATE INDEX idx_assignment_history_selected ON assignment_history(selected_engineer_id) WHERE selected_engineer_id IS NOT NULL;

-- GIN indexes for JSONB queries
CREATE INDEX idx_assignment_history_recommendations ON assignment_history USING GIN (recommendations);
CREATE INDEX idx_assignment_history_metadata ON assignment_history USING GIN (metadata);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_assignment_history_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for updated_at
CREATE TRIGGER trg_update_assignment_history_timestamp
    BEFORE UPDATE ON assignment_history
    FOR EACH ROW
    EXECUTE FUNCTION update_assignment_history_timestamp();

-- Comments
COMMENT ON TABLE assignment_history IS 'Stores AI-powered engineer assignment recommendations';
COMMENT ON COLUMN assignment_history.request_id IS 'Unique UUID for the assignment request';
COMMENT ON COLUMN assignment_history.recommendations IS 'Full recommendation list as JSONB array';
COMMENT ON COLUMN assignment_history.metadata IS 'Processing metadata (AI usage, timing, etc.)';
COMMENT ON COLUMN assignment_history.selected_engineer_id IS 'Engineer actually assigned';
COMMENT ON COLUMN assignment_history.was_successful IS 'Whether the assignment was successful';

-- View for assignment analytics
CREATE OR REPLACE VIEW v_assignment_analytics AS
SELECT 
    DATE_TRUNC('day', created_at) as assignment_date,
    COUNT(*) as total_assignments,
    COUNT(selected_engineer_id) as assignments_made,
    COUNT(CASE WHEN was_successful = true THEN 1 END) as successful_assignments,
    AVG(CASE WHEN was_successful = true THEN 100.0 ELSE 0.0 END) as success_rate,
    AVG(EXTRACT(EPOCH FROM actual_resolution_time)/3600.0) as avg_resolution_hours,
    SUM((metadata->>'cost_usd')::DECIMAL) as total_ai_cost,
    COUNT(CASE WHEN (metadata->>'used_ai')::BOOLEAN = true THEN 1 END) as ai_assisted_count
FROM assignment_history
GROUP BY DATE_TRUNC('day', created_at)
ORDER BY assignment_date DESC;

COMMENT ON VIEW v_assignment_analytics IS 'Daily analytics on assignment recommendations';

-- View for engineer assignment success rates
CREATE OR REPLACE VIEW v_engineer_assignment_success AS
SELECT 
    u.user_id,
    u.full_name,
    COUNT(*) as total_assignments,
    COUNT(CASE WHEN ah.was_successful = true THEN 1 END) as successful_count,
    ROUND(AVG(CASE WHEN ah.was_successful = true THEN 100.0 ELSE 0.0 END), 2) as success_percentage,
    ROUND(AVG(EXTRACT(EPOCH FROM ah.actual_resolution_time)/3600.0), 2) as avg_resolution_hours,
    COUNT(CASE WHEN ah.created_at > NOW() - INTERVAL '30 days' THEN 1 END) as assignments_last_30_days
FROM users u
JOIN assignment_history ah ON ah.selected_engineer_id = u.user_id
WHERE ah.was_successful IS NOT NULL
GROUP BY u.user_id, u.full_name
ORDER BY success_percentage DESC;

COMMENT ON VIEW v_engineer_assignment_success IS 'Success rates by engineer for assignments';

-- Sample query examples (commented)

-- Get recent assignments for a ticket:
-- SELECT request_id, selected_engineer_id, was_successful, created_at
-- FROM assignment_history
-- WHERE ticket_id = 123
-- ORDER BY created_at DESC;

-- Find recommendations where top choice was selected:
-- SELECT 
--     ah.request_id,
--     ah.ticket_id,
--     ah.selected_engineer_id,
--     (ah.recommendations->0->>'engineer_id')::BIGINT as recommended_engineer_id,
--     CASE 
--         WHEN ah.selected_engineer_id = (ah.recommendations->0->>'engineer_id')::BIGINT 
--         THEN 'Top choice selected'
--         ELSE 'Different engineer selected'
--     END as selection_type
-- FROM assignment_history ah
-- WHERE ah.selected_engineer_id IS NOT NULL;

-- AI cost summary:
-- SELECT 
--     DATE_TRUNC('month', created_at) as month,
--     SUM((metadata->>'cost_usd')::DECIMAL) as total_cost,
--     COUNT(*) as requests
-- FROM assignment_history
-- WHERE (metadata->>'used_ai')::BOOLEAN = true
-- GROUP BY DATE_TRUNC('month', created_at)
-- ORDER BY month DESC;

-- Success rate by recommendation rank:
-- WITH ranked_selections AS (
--     SELECT 
--         ah.request_id,
--         ah.selected_engineer_id,
--         ah.was_successful,
--         jsonb_array_elements(ah.recommendations) as rec
--     FROM assignment_history ah
--     WHERE ah.selected_engineer_id IS NOT NULL
-- )
-- SELECT 
--     (rec->>'rank')::INT as rank,
--     COUNT(*) as selected_count,
--     AVG(CASE WHEN was_successful = true THEN 100.0 ELSE 0.0 END) as success_rate
-- FROM ranked_selections
-- WHERE (rec->>'engineer_id')::BIGINT = selected_engineer_id
-- GROUP BY (rec->>'rank')::INT
-- ORDER BY rank;
