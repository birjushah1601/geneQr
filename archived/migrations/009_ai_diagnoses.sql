-- Migration: AI Diagnoses Table
-- Description: Stores AI-powered ticket diagnoses with full diagnosis data and feedback
-- Dependencies: service_tickets table
-- Phase: 2C - AI Services Layer

-- Create ai_diagnoses table
CREATE TABLE IF NOT EXISTS ai_diagnoses (
    -- Primary key
    id BIGSERIAL PRIMARY KEY,
    
    -- Diagnosis identification
    diagnosis_id VARCHAR(100) UNIQUE NOT NULL,
    ticket_id BIGINT NOT NULL REFERENCES service_tickets(ticket_id) ON DELETE CASCADE,
    
    -- Diagnosis data (full JSON response)
    diagnosis_data JSONB NOT NULL,
    
    -- Quick access fields
    confidence_score DECIMAL(5,2) CHECK (confidence_score >= 0 AND confidence_score <= 100),
    problem_category VARCHAR(50),
    problem_type VARCHAR(100),
    severity VARCHAR(20),
    
    -- AI metadata
    provider VARCHAR(50) NOT NULL, -- openai, anthropic
    model VARCHAR(100) NOT NULL,
    tokens_used INTEGER,
    cost_usd DECIMAL(10,6),
    
    -- Feedback from engineers
    was_accurate BOOLEAN,
    accuracy_score INTEGER CHECK (accuracy_score >= 0 AND accuracy_score <= 100),
    feedback_notes TEXT,
    actual_resolution TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_ai_diagnoses_ticket ON ai_diagnoses(ticket_id);
CREATE INDEX idx_ai_diagnoses_created ON ai_diagnoses(created_at DESC);
CREATE INDEX idx_ai_diagnoses_provider ON ai_diagnoses(provider);
CREATE INDEX idx_ai_diagnoses_confidence ON ai_diagnoses(confidence_score DESC);
CREATE INDEX idx_ai_diagnoses_feedback ON ai_diagnoses(was_accurate) WHERE was_accurate IS NOT NULL;

-- GIN index for JSONB queries
CREATE INDEX idx_ai_diagnoses_data ON ai_diagnoses USING GIN (diagnosis_data);

-- Function to extract confidence from JSONB
CREATE OR REPLACE FUNCTION extract_diagnosis_metadata()
RETURNS TRIGGER AS $$
BEGIN
    -- Extract quick access fields from JSONB
    NEW.confidence_score := (NEW.diagnosis_data->'PrimaryDiagnosis'->>'Confidence')::DECIMAL;
    NEW.problem_category := NEW.diagnosis_data->'PrimaryDiagnosis'->>'ProblemCategory';
    NEW.problem_type := NEW.diagnosis_data->'PrimaryDiagnosis'->>'ProblemType';
    NEW.severity := NEW.diagnosis_data->'PrimaryDiagnosis'->>'Severity';
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-extract metadata
CREATE TRIGGER trg_extract_diagnosis_metadata
    BEFORE INSERT OR UPDATE ON ai_diagnoses
    FOR EACH ROW
    EXECUTE FUNCTION extract_diagnosis_metadata();

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_ai_diagnosis_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for updated_at
CREATE TRIGGER trg_update_ai_diagnosis_timestamp
    BEFORE UPDATE ON ai_diagnoses
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_diagnosis_timestamp();

-- Comments
COMMENT ON TABLE ai_diagnoses IS 'Stores AI-powered diagnoses for service tickets';
COMMENT ON COLUMN ai_diagnoses.diagnosis_id IS 'Unique UUID for the diagnosis';
COMMENT ON COLUMN ai_diagnoses.diagnosis_data IS 'Full diagnosis response as JSONB';
COMMENT ON COLUMN ai_diagnoses.confidence_score IS 'AI confidence level (0-100%)';
COMMENT ON COLUMN ai_diagnoses.provider IS 'AI provider used (openai, anthropic)';
COMMENT ON COLUMN ai_diagnoses.was_accurate IS 'Engineer feedback on accuracy';
COMMENT ON COLUMN ai_diagnoses.accuracy_score IS 'Engineer-provided accuracy rating';
COMMENT ON COLUMN ai_diagnoses.actual_resolution IS 'What actually fixed the issue';

-- View for diagnosis analytics
CREATE OR REPLACE VIEW v_diagnosis_analytics AS
SELECT 
    DATE_TRUNC('day', created_at) as diagnosis_date,
    provider,
    model,
    COUNT(*) as total_diagnoses,
    AVG(confidence_score) as avg_confidence,
    AVG(CASE WHEN was_accurate = true THEN 100.0 ELSE 0.0 END) as accuracy_rate,
    AVG(accuracy_score) as avg_accuracy_score,
    SUM(tokens_used) as total_tokens,
    SUM(cost_usd) as total_cost,
    COUNT(CASE WHEN was_accurate = true THEN 1 END) as accurate_count,
    COUNT(CASE WHEN was_accurate = false THEN 1 END) as inaccurate_count,
    COUNT(CASE WHEN was_accurate IS NULL THEN 1 END) as pending_feedback
FROM ai_diagnoses
GROUP BY DATE_TRUNC('day', created_at), provider, model
ORDER BY diagnosis_date DESC, provider;

COMMENT ON VIEW v_diagnosis_analytics IS 'Daily analytics on AI diagnosis performance';

-- View for diagnosis feedback summary
CREATE OR REPLACE VIEW v_diagnosis_feedback_summary AS
SELECT 
    problem_category,
    COUNT(*) as total_diagnoses,
    COUNT(CASE WHEN was_accurate = true THEN 1 END) as accurate_count,
    ROUND(AVG(CASE WHEN was_accurate = true THEN 100.0 ELSE 0.0 END), 2) as accuracy_percentage,
    ROUND(AVG(confidence_score), 2) as avg_confidence,
    ROUND(AVG(accuracy_score), 2) as avg_engineer_rating
FROM ai_diagnoses
WHERE was_accurate IS NOT NULL
GROUP BY problem_category
ORDER BY accuracy_percentage DESC;

COMMENT ON VIEW v_diagnosis_feedback_summary IS 'Accuracy metrics by problem category';

-- Sample query examples (commented)

-- Get recent diagnoses for a ticket:
-- SELECT diagnosis_id, confidence_score, problem_category, created_at
-- FROM ai_diagnoses
-- WHERE ticket_id = 123
-- ORDER BY created_at DESC;

-- Find low-confidence diagnoses needing review:
-- SELECT d.diagnosis_id, d.ticket_id, t.ticket_number, d.confidence_score
-- FROM ai_diagnoses d
-- JOIN service_tickets t ON d.ticket_id = t.ticket_id
-- WHERE d.confidence_score < 50 AND d.was_accurate IS NULL
-- ORDER BY d.created_at DESC;

-- Get AI cost summary:
-- SELECT 
--     provider,
--     model,
--     SUM(cost_usd) as total_cost,
--     SUM(tokens_used) as total_tokens,
--     COUNT(*) as request_count
-- FROM ai_diagnoses
-- WHERE created_at > NOW() - INTERVAL '30 days'
-- GROUP BY provider, model;

-- Accuracy by provider:
-- SELECT 
--     provider,
--     COUNT(*) as total,
--     AVG(CASE WHEN was_accurate = true THEN 100.0 ELSE 0.0 END) as accuracy_rate
-- FROM ai_diagnoses
-- WHERE was_accurate IS NOT NULL
-- GROUP BY provider;
