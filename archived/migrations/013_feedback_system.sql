-- Migration: AI Feedback and Learning System
-- Description: Centralized feedback collection and continuous learning
-- Dependencies: ai_diagnoses, assignment_history, parts_recommendations tables
-- Phase: 2C - AI Services Layer

-- Create ai_feedback table (stores ALL feedback from human and machine sources)
CREATE TABLE IF NOT EXISTS ai_feedback (
    feedback_id BIGSERIAL PRIMARY KEY,
    
    -- Source and type
    source VARCHAR(20) NOT NULL CHECK (source IN ('human', 'machine')),
    type VARCHAR(50) NOT NULL,
    
    -- Context
    ticket_id BIGINT REFERENCES service_tickets(ticket_id) ON DELETE CASCADE,
    request_id VARCHAR(100),
    service_type VARCHAR(20) NOT NULL CHECK (service_type IN ('diagnosis', 'assignment', 'parts')),
    
    -- Human feedback fields
    user_id BIGINT REFERENCES users(user_id) ON DELETE SET NULL,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    sentiment VARCHAR(20) NOT NULL CHECK (sentiment IN ('positive', 'neutral', 'negative')),
    comments TEXT,
    
    -- Machine feedback (outcomes)
    outcomes JSONB,
    
    -- Corrections (what should have been)
    corrections JSONB,
    
    -- Metadata
    metadata JSONB,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Create feedback_improvements table (stores identified improvement opportunities)
CREATE TABLE IF NOT EXISTS feedback_improvements (
    improvement_id BIGSERIAL PRIMARY KEY,
    opportunity_id VARCHAR(100) UNIQUE NOT NULL,
    
    -- Details
    title TEXT NOT NULL,
    description TEXT,
    service_type VARCHAR(20) NOT NULL,
    
    -- Classification
    impact_level VARCHAR(20) CHECK (impact_level IN ('high', 'medium', 'low')),
    implementation_type VARCHAR(30) CHECK (implementation_type IN ('prompt_tuning', 'weight_adjustment', 'config_change', 'training_data')),
    
    -- Suggested changes
    suggested_changes JSONB NOT NULL,
    supporting_data BIGINT[], -- Array of feedback IDs
    
    -- Status tracking
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'applied', 'rejected')),
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create feedback_actions table (stores learning actions taken)
CREATE TABLE IF NOT EXISTS feedback_actions (
    action_id VARCHAR(100) PRIMARY KEY,
    opportunity_id VARCHAR(100) REFERENCES feedback_improvements(opportunity_id) ON DELETE CASCADE,
    
    -- Action details
    action_type VARCHAR(30) NOT NULL CHECK (action_type IN ('prompt_update', 'weight_adjustment', 'config_change')),
    service_type VARCHAR(20) NOT NULL,
    
    -- Changes made
    changes JSONB NOT NULL,
    
    -- Before/after metrics
    before_metrics JSONB,
    after_metrics JSONB,
    
    -- Status
    status VARCHAR(20) DEFAULT 'testing' CHECK (status IN ('testing', 'deployed', 'rolled_back')),
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    applied_by VARCHAR(100),
    
    -- Results
    result_notes TEXT,
    rolled_back_at TIMESTAMP WITH TIME ZONE,
    rollback_reason TEXT
);

-- Indexes for ai_feedback
CREATE INDEX idx_feedback_source ON ai_feedback(source);
CREATE INDEX idx_feedback_service_type ON ai_feedback(service_type);
CREATE INDEX idx_feedback_ticket ON ai_feedback(ticket_id);
CREATE INDEX idx_feedback_request ON ai_feedback(service_type, request_id);
CREATE INDEX idx_feedback_sentiment ON ai_feedback(sentiment);
CREATE INDEX idx_feedback_created ON ai_feedback(created_at DESC);
CREATE INDEX idx_feedback_user ON ai_feedback(user_id) WHERE user_id IS NOT NULL;

-- GIN indexes for JSONB
CREATE INDEX idx_feedback_outcomes ON ai_feedback USING GIN (outcomes);
CREATE INDEX idx_feedback_corrections ON ai_feedback USING GIN (corrections);
CREATE INDEX idx_feedback_metadata ON ai_feedback USING GIN (metadata);

-- Indexes for feedback_improvements
CREATE INDEX idx_improvements_service ON feedback_improvements(service_type);
CREATE INDEX idx_improvements_status ON feedback_improvements(status);
CREATE INDEX idx_improvements_impact ON feedback_improvements(impact_level);
CREATE INDEX idx_improvements_created ON feedback_improvements(created_at DESC);
CREATE INDEX idx_improvements_supporting_data ON feedback_improvements USING GIN (supporting_data);

-- Indexes for feedback_actions
CREATE INDEX idx_actions_opportunity ON feedback_actions(opportunity_id);
CREATE INDEX idx_actions_service ON feedback_actions(service_type);
CREATE INDEX idx_actions_status ON feedback_actions(status);
CREATE INDEX idx_actions_applied ON feedback_actions(applied_at DESC);

-- Trigger for updated_at
CREATE TRIGGER trg_update_feedback_improvements_timestamp
    BEFORE UPDATE ON feedback_improvements
    FOR EACH ROW
    EXECUTE FUNCTION update_parts_timestamp();

-- Comments
COMMENT ON TABLE ai_feedback IS 'Centralized feedback from both human and machine sources';
COMMENT ON COLUMN ai_feedback.source IS 'human = explicit user feedback, machine = implicit system outcomes';
COMMENT ON COLUMN ai_feedback.outcomes IS 'Machine-generated outcomes (what actually happened)';
COMMENT ON COLUMN ai_feedback.corrections IS 'Human corrections (what should have been)';

COMMENT ON TABLE feedback_improvements IS 'Identified improvement opportunities from feedback analysis';
COMMENT ON COLUMN feedback_improvements.implementation_type IS 'How to implement: prompt_tuning, weight_adjustment, config_change, training_data';
COMMENT ON COLUMN feedback_improvements.supporting_data IS 'Array of feedback IDs that support this improvement';

COMMENT ON TABLE feedback_actions IS 'Learning actions taken based on feedback';
COMMENT ON COLUMN feedback_actions.status IS 'testing = evaluating impact, deployed = in production, rolled_back = reverted';
COMMENT ON COLUMN feedback_actions.before_metrics IS 'Performance metrics before the change';
COMMENT ON COLUMN feedback_actions.after_metrics IS 'Performance metrics after the change';

-- View for feedback analytics
CREATE OR REPLACE VIEW v_feedback_analytics AS
SELECT 
    service_type,
    DATE_TRUNC('day', created_at) as feedback_date,
    source,
    COUNT(*) as total_feedback,
    COUNT(CASE WHEN sentiment = 'positive' THEN 1 END) as positive_count,
    COUNT(CASE WHEN sentiment = 'neutral' THEN 1 END) as neutral_count,
    COUNT(CASE WHEN sentiment = 'negative' THEN 1 END) as negative_count,
    AVG(rating) as avg_rating,
    AVG(CASE WHEN sentiment = 'positive' THEN 100.0 ELSE 0.0 END) as accuracy_rate
FROM ai_feedback
GROUP BY service_type, DATE_TRUNC('day', created_at), source
ORDER BY feedback_date DESC, service_type;

COMMENT ON VIEW v_feedback_analytics IS 'Daily feedback analytics by service type and source';

-- View for learning progress
CREATE OR REPLACE VIEW v_learning_progress AS
SELECT 
    service_type,
    COUNT(DISTINCT fi.improvement_id) as total_improvements,
    COUNT(DISTINCT CASE WHEN fi.status = 'pending' THEN fi.improvement_id END) as pending_improvements,
    COUNT(DISTINCT CASE WHEN fi.status = 'applied' THEN fi.improvement_id END) as applied_improvements,
    COUNT(DISTINCT fa.action_id) as total_actions,
    COUNT(DISTINCT CASE WHEN fa.status = 'deployed' THEN fa.action_id END) as deployed_actions,
    COUNT(DISTINCT CASE WHEN fa.status = 'rolled_back' THEN fa.action_id END) as rolled_back_actions,
    ROUND(
        CASE 
            WHEN COUNT(DISTINCT fa.action_id) > 0 
            THEN COUNT(DISTINCT CASE WHEN fa.status = 'deployed' THEN fa.action_id END) * 100.0 / COUNT(DISTINCT fa.action_id)
            ELSE 0 
        END, 
        2
    ) as success_rate_percent
FROM feedback_improvements fi
LEFT JOIN feedback_actions fa ON fi.opportunity_id = fa.opportunity_id
GROUP BY service_type;

COMMENT ON VIEW v_learning_progress IS 'Learning progress and success rates by service type';

-- View for top issues
CREATE OR REPLACE VIEW v_top_feedback_issues AS
SELECT 
    service_type,
    type as feedback_type,
    COUNT(*) as issue_count,
    COUNT(CASE WHEN sentiment = 'negative' THEN 1 END) as negative_count,
    AVG(rating) as avg_rating,
    ARRAY_AGG(feedback_id ORDER BY created_at DESC LIMIT 5) as recent_examples
FROM ai_feedback
WHERE sentiment = 'negative'
    AND created_at > NOW() - INTERVAL '30 days'
GROUP BY service_type, type
HAVING COUNT(*) >= 3
ORDER BY issue_count DESC;

COMMENT ON VIEW v_top_feedback_issues IS 'Top recurring issues from negative feedback';
