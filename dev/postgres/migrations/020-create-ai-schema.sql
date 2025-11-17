-- Migration: 020-create-ai-schema.sql
-- Description: AI Integration Schema for Multi-Provider Orchestration
-- Ticket: T2B.5
-- Date: 2025-11-16
-- 
-- This migration creates:
-- 1. ai_providers - Provider configuration (OpenAI, Anthropic, etc.)
-- 2. ai_models - Available models per provider
-- 3. ai_conversations - Chat/conversation history with AI
-- 4. ai_diagnosis_results - AI-generated diagnosis from tickets
-- 5. ai_engineer_recommendations - AI-recommended engineers
-- 6. ai_parts_recommendations - AI-recommended parts
-- 7. ai_attachment_analysis - AI analysis of attachments
-- 8. ai_feedback - Feedback loop for AI learning
--
-- Purpose:
-- - Support multiple AI providers (OpenAI, Anthropic, custom)
-- - Track all AI interactions and results
-- - Enable feedback loop for continuous improvement
-- - Cost tracking per provider/model
-- - Confidence scoring and validation

-- =====================================================================
-- 1. AI PROVIDERS (OpenAI, Anthropic, etc.)
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Provider Identity
    provider_name TEXT NOT NULL UNIQUE,            -- 'openai', 'anthropic', 'custom'
    provider_code TEXT NOT NULL UNIQUE,            -- 'openai-gpt', 'anthropic-claude'
    display_name TEXT NOT NULL,
    
    -- Configuration
    api_base_url TEXT NOT NULL,
    api_version TEXT,
    requires_api_key BOOLEAN DEFAULT true,
    api_key_env_var TEXT,                          -- Environment variable name for API key
    
    -- Capabilities
    supports_chat BOOLEAN DEFAULT true,
    supports_vision BOOLEAN DEFAULT false,
    supports_function_calling BOOLEAN DEFAULT false,
    supports_streaming BOOLEAN DEFAULT false,
    supports_embeddings BOOLEAN DEFAULT false,
    
    -- Limits
    max_tokens_per_request INT,
    max_context_window INT,
    rate_limit_requests_per_minute INT,
    rate_limit_tokens_per_minute INT,
    
    -- Pricing (per 1M tokens)
    default_input_cost_per_1m NUMERIC(10,4),
    default_output_cost_per_1m NUMERIC(10,4),
    currency TEXT DEFAULT 'USD',
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_primary BOOLEAN DEFAULT false,              -- Primary provider to use
    priority INT DEFAULT 5,                        -- Fallback priority
    
    -- Retry and Timeout
    default_timeout_seconds INT DEFAULT 30,
    max_retries INT DEFAULT 3,
    retry_backoff_multiplier NUMERIC(3,2) DEFAULT 2.0,
    
    -- Configuration
    default_temperature NUMERIC(3,2) DEFAULT 0.7,
    default_top_p NUMERIC(3,2) DEFAULT 1.0,
    provider_config JSONB DEFAULT '{}',            -- Provider-specific config
    
    -- Health Check
    last_health_check TIMESTAMPTZ,
    is_healthy BOOLEAN DEFAULT true,
    health_check_failure_count INT DEFAULT 0,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    notes TEXT,
    
    -- Constraints
    CONSTRAINT chk_priority CHECK (priority >= 1 AND priority <= 10),
    CONSTRAINT chk_temperature CHECK (default_temperature >= 0 AND default_temperature <= 2),
    CONSTRAINT chk_top_p CHECK (default_top_p >= 0 AND default_top_p <= 1)
);

-- Indexes
CREATE INDEX idx_ai_providers_active ON ai_providers(is_active) WHERE is_active = true;
CREATE INDEX idx_ai_providers_primary ON ai_providers(is_primary) WHERE is_primary = true;
CREATE INDEX idx_ai_providers_priority ON ai_providers(priority ASC) WHERE is_active = true;

-- GIN index for provider config
CREATE INDEX idx_ai_providers_config ON ai_providers USING GIN (provider_config);

COMMENT ON TABLE ai_providers IS 'AI provider configuration (OpenAI, Anthropic, custom)';
COMMENT ON COLUMN ai_providers.is_primary IS 'Primary provider to use by default';
COMMENT ON COLUMN ai_providers.priority IS 'Fallback priority (lower number = higher priority)';

-- =====================================================================
-- 2. AI MODELS (Available models per provider)
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Provider Relationship
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    
    -- Model Identity
    model_name TEXT NOT NULL,                      -- 'gpt-4', 'claude-3-opus', etc.
    model_code TEXT NOT NULL,                      -- Unique code for API calls
    display_name TEXT NOT NULL,
    model_version TEXT,
    
    -- Capabilities
    supports_vision BOOLEAN DEFAULT false,
    supports_function_calling BOOLEAN DEFAULT false,
    max_context_tokens INT NOT NULL,
    max_output_tokens INT,
    
    -- Pricing (per 1M tokens)
    input_cost_per_1m NUMERIC(10,4) NOT NULL,
    output_cost_per_1m NUMERIC(10,4) NOT NULL,
    currency TEXT DEFAULT 'USD',
    
    -- Performance
    avg_latency_ms INT,
    avg_tokens_per_second INT,
    
    -- Use Cases
    recommended_for TEXT[],                        -- ['diagnosis', 'chat', 'vision', 'analysis']
    
    -- Configuration
    default_temperature NUMERIC(3,2) DEFAULT 0.7,
    default_max_tokens INT DEFAULT 1000,
    model_config JSONB DEFAULT '{}',
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,              -- Default model for this provider
    is_deprecated BOOLEAN DEFAULT false,
    deprecation_date DATE,
    replacement_model_id UUID REFERENCES ai_models(id),
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    notes TEXT,
    
    -- Constraints
    CONSTRAINT chk_tokens CHECK (max_context_tokens > 0 AND max_output_tokens > 0)
);

-- Unique constraint: model_code per provider
CREATE UNIQUE INDEX idx_ai_models_unique ON ai_models(provider_id, model_code);

-- Indexes
CREATE INDEX idx_ai_models_provider ON ai_models(provider_id);
CREATE INDEX idx_ai_models_active ON ai_models(is_active) WHERE is_active = true;
CREATE INDEX idx_ai_models_default ON ai_models(provider_id, is_default) WHERE is_default = true;

-- GIN index for recommended_for and config
CREATE INDEX idx_ai_models_recommended ON ai_models USING GIN (recommended_for);
CREATE INDEX idx_ai_models_config ON ai_models USING GIN (model_config);

COMMENT ON TABLE ai_models IS 'Available AI models per provider with pricing and capabilities';
COMMENT ON COLUMN ai_models.is_default IS 'Default model to use for this provider';
COMMENT ON COLUMN ai_models.recommended_for IS 'Use cases this model is recommended for';

-- =====================================================================
-- 3. AI CONVERSATIONS (Chat history with AI)
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Context
    ticket_id UUID REFERENCES service_tickets(id),
    workflow_instance_id UUID REFERENCES ticket_workflow_instances(id),
    stage_id UUID REFERENCES stage_configuration_templates(id),
    
    -- Conversation Identity
    conversation_id TEXT NOT NULL,                 -- Session ID
    conversation_type TEXT NOT NULL,               -- 'diagnosis', 'parts_recommendation', 'engineer_assignment', 'general'
    
    -- AI Configuration
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    provider_name TEXT NOT NULL,
    model_name TEXT NOT NULL,
    
    -- Message
    message_role TEXT NOT NULL,                    -- 'system', 'user', 'assistant', 'function'
    message_content TEXT NOT NULL,
    message_order INT NOT NULL,                    -- Order in conversation
    
    -- Function Calling
    function_name TEXT,
    function_arguments JSONB,
    function_result JSONB,
    
    -- Attachments/Context
    includes_attachments BOOLEAN DEFAULT false,
    attachment_ids UUID[],
    context_data JSONB DEFAULT '{}',               -- Additional context provided
    
    -- Token Usage
    prompt_tokens INT,
    completion_tokens INT,
    total_tokens INT,
    
    -- Cost
    cost_usd NUMERIC(10,6),
    
    -- Performance
    latency_ms INT,
    
    -- Status
    is_successful BOOLEAN DEFAULT true,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    
    -- User Context
    initiated_by_user_id TEXT,
    initiated_by_engineer_id UUID REFERENCES engineers(id),
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    
    -- Constraints
    CONSTRAINT chk_message_role CHECK (message_role IN ('system', 'user', 'assistant', 'function', 'tool')),
    CONSTRAINT chk_conversation_type CHECK (conversation_type IN (
        'diagnosis', 'parts_recommendation', 'engineer_assignment', 
        'chat', 'analysis', 'support', 'general'
    ))
);

-- Indexes
CREATE INDEX idx_ai_conversations_ticket ON ai_conversations(ticket_id);
CREATE INDEX idx_ai_conversations_workflow ON ai_conversations(workflow_instance_id);
CREATE INDEX idx_ai_conversations_session ON ai_conversations(conversation_id, message_order);
CREATE INDEX idx_ai_conversations_type ON ai_conversations(conversation_type);
CREATE INDEX idx_ai_conversations_provider ON ai_conversations(provider_id);
CREATE INDEX idx_ai_conversations_model ON ai_conversations(model_id);
CREATE INDEX idx_ai_conversations_time ON ai_conversations(created_at DESC);

-- Composite index for conversation lookup
CREATE INDEX idx_ai_conversations_lookup ON ai_conversations(
    conversation_id, message_order
);

-- GIN indexes
CREATE INDEX idx_ai_conversations_attachments ON ai_conversations USING GIN (attachment_ids);
CREATE INDEX idx_ai_conversations_context ON ai_conversations USING GIN (context_data);

COMMENT ON TABLE ai_conversations IS 'Complete conversation history with AI providers';
COMMENT ON COLUMN ai_conversations.conversation_id IS 'Session ID to group related messages';
COMMENT ON COLUMN ai_conversations.message_order IS 'Order of messages in conversation';

-- =====================================================================
-- 4. AI DIAGNOSIS RESULTS
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_diagnosis_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID REFERENCES ticket_workflow_instances(id),
    conversation_id TEXT,                          -- Link to ai_conversations
    
    -- AI Configuration
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    provider_name TEXT NOT NULL,
    model_name TEXT NOT NULL,
    
    -- Input Data
    input_description TEXT NOT NULL,               -- Ticket description
    input_attachments UUID[],                      -- Attachment IDs analyzed
    equipment_context JSONB,                       -- Equipment details
    historical_context JSONB,                      -- Similar past tickets
    
    -- Diagnosis Output
    diagnosis_summary TEXT NOT NULL,
    root_cause TEXT,
    confidence_score NUMERIC(5,2),                 -- 0-100%
    
    -- Identified Issues
    identified_issues JSONB NOT NULL DEFAULT '[]', -- Array of issue objects
    issue_categories TEXT[],
    severity_level TEXT,                           -- 'critical', 'high', 'medium', 'low'
    
    -- Recommendations
    recommended_actions TEXT[],
    recommended_parts UUID[],                      -- References to equipment_parts
    estimated_resolution_time_hours NUMERIC(5,2),
    
    -- Support Level Required
    recommended_support_level TEXT,                -- 'L1', 'L2', 'L3'
    requires_specialist BOOLEAN DEFAULT false,
    specialist_skills TEXT[],
    
    -- Troubleshooting Steps
    troubleshooting_steps JSONB DEFAULT '[]',
    
    -- Token Usage
    prompt_tokens INT,
    completion_tokens INT,
    total_tokens INT,
    cost_usd NUMERIC(10,6),
    
    -- Performance
    processing_time_ms INT,
    
    -- Validation
    is_validated BOOLEAN DEFAULT false,
    validated_by TEXT,
    validated_at TIMESTAMPTZ,
    validation_score INT,                          -- 1-5 rating by human
    validation_notes TEXT,
    
    -- Actual Outcome
    was_accurate BOOLEAN,
    actual_root_cause TEXT,
    accuracy_feedback TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_confidence CHECK (confidence_score >= 0 AND confidence_score <= 100),
    CONSTRAINT chk_severity CHECK (severity_level IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT chk_support_level CHECK (recommended_support_level IN ('L1', 'L2', 'L3')),
    CONSTRAINT chk_validation_score CHECK (validation_score IS NULL OR (validation_score >= 1 AND validation_score <= 5))
);

-- One diagnosis per ticket (can be updated)
CREATE UNIQUE INDEX idx_ai_diagnosis_ticket ON ai_diagnosis_results(ticket_id);

-- Indexes
CREATE INDEX idx_ai_diagnosis_workflow ON ai_diagnosis_results(workflow_instance_id);
CREATE INDEX idx_ai_diagnosis_provider ON ai_diagnosis_results(provider_id);
CREATE INDEX idx_ai_diagnosis_model ON ai_diagnosis_results(model_id);
CREATE INDEX idx_ai_diagnosis_confidence ON ai_diagnosis_results(confidence_score DESC);
CREATE INDEX idx_ai_diagnosis_severity ON ai_diagnosis_results(severity_level);
CREATE INDEX idx_ai_diagnosis_validated ON ai_diagnosis_results(is_validated);
CREATE INDEX idx_ai_diagnosis_accurate ON ai_diagnosis_results(was_accurate) WHERE was_accurate IS NOT NULL;

-- GIN indexes
CREATE INDEX idx_ai_diagnosis_issues ON ai_diagnosis_results USING GIN (identified_issues);
CREATE INDEX idx_ai_diagnosis_categories ON ai_diagnosis_results USING GIN (issue_categories);
CREATE INDEX idx_ai_diagnosis_actions ON ai_diagnosis_results USING GIN (recommended_actions);
CREATE INDEX idx_ai_diagnosis_parts ON ai_diagnosis_results USING GIN (recommended_parts);

COMMENT ON TABLE ai_diagnosis_results IS 'AI-generated diagnosis from ticket descriptions and attachments';
COMMENT ON COLUMN ai_diagnosis_results.confidence_score IS 'AI confidence in diagnosis (0-100%)';
COMMENT ON COLUMN ai_diagnosis_results.was_accurate IS 'Feedback: Was AI diagnosis correct?';

-- =====================================================================
-- 5. AI ENGINEER RECOMMENDATIONS
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_engineer_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID REFERENCES ticket_workflow_instances(id),
    stage_id UUID REFERENCES stage_configuration_templates(id),
    conversation_id TEXT,
    
    -- AI Configuration
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    
    -- Recommendation
    engineer_id UUID NOT NULL REFERENCES engineers(id),
    recommendation_rank INT NOT NULL,              -- 1=best, 2=second best, etc.
    
    -- Scoring
    overall_score NUMERIC(5,2) NOT NULL,           -- 0-100
    confidence_score NUMERIC(5,2),
    
    -- Score Breakdown
    expertise_score NUMERIC(5,2),
    availability_score NUMERIC(5,2),
    location_proximity_score NUMERIC(5,2),
    performance_history_score NUMERIC(5,2),
    workload_score NUMERIC(5,2),
    
    -- Reasoning
    recommendation_reason TEXT NOT NULL,
    strengths TEXT[],
    concerns TEXT[],
    
    -- Context
    required_support_level TEXT,
    engineer_support_level TEXT,
    is_certified BOOLEAN,
    certification_required BOOLEAN,
    distance_km NUMERIC(8,2),
    estimated_travel_time_hours NUMERIC(5,2),
    current_workload_tickets INT,
    
    -- Token Usage
    total_tokens INT,
    cost_usd NUMERIC(10,6),
    
    -- Outcome
    was_selected BOOLEAN DEFAULT false,
    selected_at TIMESTAMPTZ,
    performance_rating INT,                        -- 1-5 after completion
    was_good_recommendation BOOLEAN,
    outcome_feedback TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_rank CHECK (recommendation_rank > 0),
    CONSTRAINT chk_scores CHECK (
        overall_score >= 0 AND overall_score <= 100 AND
        (expertise_score IS NULL OR (expertise_score >= 0 AND expertise_score <= 100)) AND
        (availability_score IS NULL OR (availability_score >= 0 AND availability_score <= 100))
    )
);

-- Indexes
CREATE INDEX idx_ai_engineer_recs_ticket ON ai_engineer_recommendations(ticket_id);
CREATE INDEX idx_ai_engineer_recs_workflow ON ai_engineer_recommendations(workflow_instance_id);
CREATE INDEX idx_ai_engineer_recs_stage ON ai_engineer_recommendations(stage_id);
CREATE INDEX idx_ai_engineer_recs_engineer ON ai_engineer_recommendations(engineer_id);
CREATE INDEX idx_ai_engineer_recs_rank ON ai_engineer_recommendations(ticket_id, recommendation_rank);
CREATE INDEX idx_ai_engineer_recs_selected ON ai_engineer_recommendations(was_selected) WHERE was_selected = true;

-- Composite index for recommendation lookup
CREATE INDEX idx_ai_engineer_recs_lookup ON ai_engineer_recommendations(
    ticket_id, recommendation_rank
);

-- GIN indexes
CREATE INDEX idx_ai_engineer_recs_strengths ON ai_engineer_recommendations USING GIN (strengths);
CREATE INDEX idx_ai_engineer_recs_concerns ON ai_engineer_recommendations USING GIN (concerns);

COMMENT ON TABLE ai_engineer_recommendations IS 'AI-recommended engineers with scoring and reasoning';
COMMENT ON COLUMN ai_engineer_recommendations.recommendation_rank IS '1=best match, 2=second best, etc.';
COMMENT ON COLUMN ai_engineer_recommendations.was_selected IS 'Whether this engineer was actually assigned';

-- =====================================================================
-- 6. AI PARTS RECOMMENDATIONS
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_parts_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID REFERENCES ticket_workflow_instances(id),
    diagnosis_result_id UUID REFERENCES ai_diagnosis_results(id),
    conversation_id TEXT,
    
    -- AI Configuration
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    
    -- Part Recommendation
    equipment_part_id UUID REFERENCES equipment_parts(id),
    part_number TEXT NOT NULL,
    part_name TEXT NOT NULL,
    recommendation_rank INT NOT NULL,
    
    -- Scoring
    confidence_score NUMERIC(5,2),                 -- 0-100%
    relevance_score NUMERIC(5,2),
    
    -- Quantity
    recommended_quantity INT NOT NULL DEFAULT 1,
    is_critical BOOLEAN DEFAULT false,
    priority TEXT DEFAULT 'medium',
    
    -- Context
    installation_context TEXT,                     -- ICU, General Ward, etc.
    is_context_specific BOOLEAN DEFAULT false,
    
    -- Reasoning
    recommendation_reason TEXT NOT NULL,
    why_this_part TEXT,
    issue_addressed TEXT,
    
    -- Alternatives
    has_alternatives BOOLEAN DEFAULT false,
    alternative_part_ids UUID[],
    alternative_reasons TEXT,
    
    -- Pricing
    estimated_price NUMERIC(10,2),
    currency TEXT DEFAULT 'INR',
    
    -- Sourcing
    is_oem_part BOOLEAN,
    estimated_lead_time_days INT,
    
    -- Token Usage
    total_tokens INT,
    cost_usd NUMERIC(10,6),
    
    -- Outcome
    was_ordered BOOLEAN DEFAULT false,
    was_used BOOLEAN DEFAULT false,
    was_correct_part BOOLEAN,
    outcome_feedback TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_parts_priority CHECK (priority IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT chk_parts_quantity CHECK (recommended_quantity > 0),
    CONSTRAINT chk_parts_confidence CHECK (confidence_score IS NULL OR (confidence_score >= 0 AND confidence_score <= 100))
);

-- Indexes
CREATE INDEX idx_ai_parts_recs_ticket ON ai_parts_recommendations(ticket_id);
CREATE INDEX idx_ai_parts_recs_workflow ON ai_parts_recommendations(workflow_instance_id);
CREATE INDEX idx_ai_parts_recs_diagnosis ON ai_parts_recommendations(diagnosis_result_id);
CREATE INDEX idx_ai_parts_recs_part ON ai_parts_recommendations(equipment_part_id);
CREATE INDEX idx_ai_parts_recs_rank ON ai_parts_recommendations(ticket_id, recommendation_rank);
CREATE INDEX idx_ai_parts_recs_context ON ai_parts_recommendations(installation_context);
CREATE INDEX idx_ai_parts_recs_critical ON ai_parts_recommendations(is_critical) WHERE is_critical = true;

-- GIN index for alternatives
CREATE INDEX idx_ai_parts_recs_alternatives ON ai_parts_recommendations USING GIN (alternative_part_ids);

COMMENT ON TABLE ai_parts_recommendations IS 'AI-recommended parts with context-aware suggestions';
COMMENT ON COLUMN ai_parts_recommendations.is_context_specific IS 'Whether recommendation considers installation context (ICU vs Ward)';

-- =====================================================================
-- 7. AI ATTACHMENT ANALYSIS
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_attachment_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    attachment_id UUID NOT NULL REFERENCES stage_attachments(id),
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    conversation_id TEXT,
    
    -- AI Configuration
    provider_id UUID NOT NULL REFERENCES ai_providers(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    
    -- Analysis Type
    analysis_type TEXT NOT NULL,                   -- 'image_recognition', 'issue_detection', 'text_extraction', 'quality_check'
    
    -- Vision Analysis Results
    detected_objects JSONB DEFAULT '[]',           -- Objects detected in image
    detected_issues TEXT[],
    issue_severity TEXT,
    confidence_score NUMERIC(5,2),
    
    -- Content Description
    ai_description TEXT,
    technical_details TEXT,
    
    -- Issue Detection
    issues_found INT DEFAULT 0,
    issue_categories TEXT[],
    recommended_actions TEXT[],
    
    -- Quality Assessment
    image_quality_score NUMERIC(3,2),              -- 0-5
    is_usable BOOLEAN DEFAULT true,
    quality_issues TEXT[],                         -- 'blurry', 'dark', 'partial_view'
    
    -- Text Extraction (if applicable)
    extracted_text TEXT,
    extracted_data JSONB,
    
    -- Equipment Identification
    equipment_identified BOOLEAN,
    equipment_type TEXT,
    equipment_model TEXT,
    serial_number TEXT,
    
    -- Anomalies
    anomalies_detected BOOLEAN DEFAULT false,
    anomaly_descriptions TEXT[],
    
    -- Token Usage
    prompt_tokens INT,
    completion_tokens INT,
    total_tokens INT,
    cost_usd NUMERIC(10,6),
    
    -- Performance
    processing_time_ms INT,
    
    -- Validation
    is_validated BOOLEAN DEFAULT false,
    validated_by TEXT,
    validation_score INT,
    was_accurate BOOLEAN,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    
    -- Constraints
    CONSTRAINT chk_analysis_type CHECK (analysis_type IN (
        'image_recognition', 'issue_detection', 'text_extraction', 
        'quality_check', 'equipment_identification', 'anomaly_detection'
    )),
    CONSTRAINT chk_attachment_confidence CHECK (confidence_score IS NULL OR (confidence_score >= 0 AND confidence_score <= 100))
);

-- One analysis per attachment per type
CREATE UNIQUE INDEX idx_ai_attachment_analysis_unique ON ai_attachment_analysis(attachment_id, analysis_type);

-- Indexes
CREATE INDEX idx_ai_attachment_analysis_attachment ON ai_attachment_analysis(attachment_id);
CREATE INDEX idx_ai_attachment_analysis_ticket ON ai_attachment_analysis(ticket_id);
CREATE INDEX idx_ai_attachment_analysis_type ON ai_attachment_analysis(analysis_type);
CREATE INDEX idx_ai_attachment_analysis_provider ON ai_attachment_analysis(provider_id);
CREATE INDEX idx_ai_attachment_analysis_confidence ON ai_attachment_analysis(confidence_score DESC);

-- GIN indexes
CREATE INDEX idx_ai_attachment_analysis_objects ON ai_attachment_analysis USING GIN (detected_objects);
CREATE INDEX idx_ai_attachment_analysis_issues ON ai_attachment_analysis USING GIN (detected_issues);
CREATE INDEX idx_ai_attachment_analysis_categories ON ai_attachment_analysis USING GIN (issue_categories);

COMMENT ON TABLE ai_attachment_analysis IS 'AI analysis of attachments (images, videos, documents)';
COMMENT ON COLUMN ai_attachment_analysis.detected_objects IS 'JSONB array of detected objects with bounding boxes';

-- =====================================================================
-- 8. AI FEEDBACK (Learning Loop)
-- =====================================================================

CREATE TABLE IF NOT EXISTS ai_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    diagnosis_result_id UUID REFERENCES ai_diagnosis_results(id),
    engineer_recommendation_id UUID REFERENCES ai_engineer_recommendations(id),
    parts_recommendation_id UUID REFERENCES ai_parts_recommendations(id),
    attachment_analysis_id UUID REFERENCES ai_attachment_analysis(id),
    
    -- Feedback Type
    feedback_type TEXT NOT NULL,                   -- 'diagnosis', 'engineer', 'parts', 'analysis'
    feedback_category TEXT NOT NULL,               -- 'accuracy', 'usefulness', 'completeness', 'relevance'
    
    -- Rating
    rating INT NOT NULL,                           -- 1-5 stars
    was_helpful BOOLEAN NOT NULL,
    was_accurate BOOLEAN,
    
    -- Detailed Feedback
    feedback_text TEXT,
    what_was_good TEXT,
    what_was_wrong TEXT,
    what_was_missing TEXT,
    suggestions TEXT,
    
    -- Corrections
    correct_diagnosis TEXT,
    correct_root_cause TEXT,
    correct_parts UUID[],
    correct_actions TEXT[],
    
    -- Context
    actual_outcome TEXT,
    resolution_time_hours NUMERIC(10,2),
    actual_cost NUMERIC(10,2),
    
    -- Feedback Provider
    feedback_by_user_id TEXT,
    feedback_by_engineer_id UUID REFERENCES engineers(id),
    feedback_by_type TEXT,                         -- 'engineer', 'manager', 'customer', 'system'
    
    -- Training Usage
    can_use_for_training BOOLEAN DEFAULT true,
    training_priority TEXT DEFAULT 'medium',       -- 'high', 'medium', 'low'
    training_category TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    
    -- Constraints
    CONSTRAINT chk_feedback_type CHECK (feedback_type IN (
        'diagnosis', 'engineer_recommendation', 'parts_recommendation', 
        'attachment_analysis', 'general'
    )),
    CONSTRAINT chk_feedback_category CHECK (feedback_category IN (
        'accuracy', 'usefulness', 'completeness', 'relevance', 'performance', 'cost'
    )),
    CONSTRAINT chk_rating CHECK (rating >= 1 AND rating <= 5),
    CONSTRAINT chk_feedback_by_type CHECK (feedback_by_type IN (
        'engineer', 'manager', 'customer', 'admin', 'system'
    )),
    CONSTRAINT chk_training_priority CHECK (training_priority IN ('high', 'medium', 'low'))
);

-- Indexes
CREATE INDEX idx_ai_feedback_ticket ON ai_feedback(ticket_id);
CREATE INDEX idx_ai_feedback_diagnosis ON ai_feedback(diagnosis_result_id);
CREATE INDEX idx_ai_feedback_engineer_rec ON ai_feedback(engineer_recommendation_id);
CREATE INDEX idx_ai_feedback_parts_rec ON ai_feedback(parts_recommendation_id);
CREATE INDEX idx_ai_feedback_attachment ON ai_feedback(attachment_analysis_id);
CREATE INDEX idx_ai_feedback_type ON ai_feedback(feedback_type);
CREATE INDEX idx_ai_feedback_rating ON ai_feedback(rating);
CREATE INDEX idx_ai_feedback_helpful ON ai_feedback(was_helpful);
CREATE INDEX idx_ai_feedback_training ON ai_feedback(can_use_for_training) WHERE can_use_for_training = true;

-- GIN indexes
CREATE INDEX idx_ai_feedback_correct_parts ON ai_feedback USING GIN (correct_parts);
CREATE INDEX idx_ai_feedback_correct_actions ON ai_feedback USING GIN (correct_actions);

COMMENT ON TABLE ai_feedback IS 'Feedback on AI recommendations for continuous learning';
COMMENT ON COLUMN ai_feedback.can_use_for_training IS 'Whether this feedback can be used to improve AI models';

-- =====================================================================
-- 9. HELPER FUNCTIONS
-- =====================================================================

-- Function: Get AI diagnosis for ticket
CREATE OR REPLACE FUNCTION get_ai_diagnosis(
    p_ticket_id UUID
) RETURNS TABLE (
    diagnosis_id UUID,
    diagnosis_summary TEXT,
    root_cause TEXT,
    confidence_score NUMERIC,
    severity_level TEXT,
    recommended_support_level TEXT,
    recommended_parts UUID[],
    was_accurate BOOLEAN,
    provider_name TEXT,
    model_name TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        adr.id,
        adr.diagnosis_summary,
        adr.root_cause,
        adr.confidence_score,
        adr.severity_level,
        adr.recommended_support_level,
        adr.recommended_parts,
        adr.was_accurate,
        adr.provider_name,
        adr.model_name
    FROM ai_diagnosis_results adr
    WHERE adr.ticket_id = p_ticket_id;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_ai_diagnosis IS 'Get AI diagnosis results for a ticket';

-- Function: Get AI engineer recommendations
CREATE OR REPLACE FUNCTION get_ai_engineer_recommendations(
    p_ticket_id UUID,
    p_stage_id UUID DEFAULT NULL
) RETURNS TABLE (
    recommendation_id UUID,
    engineer_id UUID,
    engineer_name TEXT,
    recommendation_rank INT,
    overall_score NUMERIC,
    recommendation_reason TEXT,
    distance_km NUMERIC,
    was_selected BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        aer.id,
        aer.engineer_id,
        e.name,
        aer.recommendation_rank,
        aer.overall_score,
        aer.recommendation_reason,
        aer.distance_km,
        aer.was_selected
    FROM ai_engineer_recommendations aer
    JOIN engineers e ON aer.engineer_id = e.id
    WHERE aer.ticket_id = p_ticket_id
      AND (p_stage_id IS NULL OR aer.stage_id = p_stage_id)
    ORDER BY aer.recommendation_rank;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_ai_engineer_recommendations IS 'Get AI-recommended engineers ranked by score';

-- Function: Get AI parts recommendations
CREATE OR REPLACE FUNCTION get_ai_parts_recommendations(
    p_ticket_id UUID
) RETURNS TABLE (
    recommendation_id UUID,
    part_id UUID,
    part_number TEXT,
    part_name TEXT,
    recommendation_rank INT,
    confidence_score NUMERIC,
    recommended_quantity INT,
    is_critical BOOLEAN,
    recommendation_reason TEXT,
    estimated_price NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        apr.id,
        apr.equipment_part_id,
        apr.part_number,
        apr.part_name,
        apr.recommendation_rank,
        apr.confidence_score,
        apr.recommended_quantity,
        apr.is_critical,
        apr.recommendation_reason,
        apr.estimated_price
    FROM ai_parts_recommendations apr
    WHERE apr.ticket_id = p_ticket_id
    ORDER BY apr.is_critical DESC, apr.recommendation_rank;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_ai_parts_recommendations IS 'Get AI-recommended parts ranked by priority';

-- Function: Get conversation history
CREATE OR REPLACE FUNCTION get_conversation_history(
    p_conversation_id TEXT,
    p_limit INT DEFAULT 50
) RETURNS TABLE (
    message_id UUID,
    message_role TEXT,
    message_content TEXT,
    message_order INT,
    provider_name TEXT,
    model_name TEXT,
    created_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ac.id,
        ac.message_role,
        ac.message_content,
        ac.message_order,
        ac.provider_name,
        ac.model_name,
        ac.created_at
    FROM ai_conversations ac
    WHERE ac.conversation_id = p_conversation_id
    ORDER BY ac.message_order ASC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_conversation_history IS 'Get conversation history in order';

-- Function: Calculate AI accuracy rate
CREATE OR REPLACE FUNCTION calculate_ai_accuracy(
    p_provider_id UUID DEFAULT NULL,
    p_days_back INT DEFAULT 30
) RETURNS TABLE (
    provider_name TEXT,
    total_diagnoses INT,
    validated_count INT,
    accurate_count INT,
    accuracy_rate NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        adr.provider_name,
        COUNT(*)::INT as total,
        COUNT(CASE WHEN adr.is_validated THEN 1 END)::INT as validated,
        COUNT(CASE WHEN adr.was_accurate = true THEN 1 END)::INT as accurate,
        ROUND(
            (COUNT(CASE WHEN adr.was_accurate = true THEN 1 END)::NUMERIC / 
             NULLIF(COUNT(CASE WHEN adr.is_validated THEN 1 END), 0)) * 100, 
            2
        ) as accuracy
    FROM ai_diagnosis_results adr
    WHERE (p_provider_id IS NULL OR adr.provider_id = p_provider_id)
      AND adr.created_at >= NOW() - (p_days_back || ' days')::INTERVAL
    GROUP BY adr.provider_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION calculate_ai_accuracy IS 'Calculate AI diagnosis accuracy rate by provider';

-- =====================================================================
-- 10. VIEWS
-- =====================================================================

-- View: AI Performance Summary
CREATE OR REPLACE VIEW ai_performance_summary AS
SELECT 
    ap.provider_name,
    am.model_name,
    COUNT(DISTINCT adr.id) as total_diagnoses,
    COUNT(CASE WHEN adr.is_validated THEN 1 END) as validated_diagnoses,
    COUNT(CASE WHEN adr.was_accurate = true THEN 1 END) as accurate_diagnoses,
    ROUND(AVG(adr.confidence_score), 2) as avg_confidence,
    ROUND(AVG(adr.processing_time_ms), 0) as avg_processing_ms,
    SUM(adr.cost_usd) as total_cost_usd,
    ROUND(
        (COUNT(CASE WHEN adr.was_accurate = true THEN 1 END)::NUMERIC / 
         NULLIF(COUNT(CASE WHEN adr.is_validated THEN 1 END), 0)) * 100, 
        2
    ) as accuracy_rate_pct
FROM ai_providers ap
JOIN ai_models am ON ap.id = am.provider_id
LEFT JOIN ai_diagnosis_results adr ON am.id = adr.model_id
WHERE ap.is_active = true
GROUP BY ap.provider_name, am.model_name;

COMMENT ON VIEW ai_performance_summary IS 'AI provider and model performance metrics';

-- View: AI Cost Summary
CREATE OR REPLACE VIEW ai_cost_summary AS
SELECT 
    DATE_TRUNC('day', created_at) as date,
    provider_name,
    COUNT(*) as total_requests,
    SUM(total_tokens) as total_tokens,
    SUM(cost_usd) as total_cost_usd,
    ROUND(AVG(cost_usd), 6) as avg_cost_per_request
FROM ai_conversations
GROUP BY DATE_TRUNC('day', created_at), provider_name
ORDER BY date DESC, provider_name;

COMMENT ON VIEW ai_cost_summary IS 'Daily AI usage and cost summary by provider';

-- =====================================================================
-- 11. SEED DEFAULT PROVIDERS
-- =====================================================================

-- Insert OpenAI provider
INSERT INTO ai_providers (
    provider_name, provider_code, display_name,
    api_base_url, api_version, api_key_env_var,
    supports_chat, supports_vision, supports_function_calling, supports_streaming,
    max_context_window, rate_limit_requests_per_minute,
    default_input_cost_per_1m, default_output_cost_per_1m,
    is_active, is_primary, priority
) VALUES (
    'openai', 'openai-gpt', 'OpenAI ChatGPT',
    'https://api.openai.com/v1', 'v1', 'OPENAI_API_KEY',
    true, true, true, true,
    128000, 10000,
    5.00, 15.00,
    true, true, 1
);

-- Insert Anthropic provider
INSERT INTO ai_providers (
    provider_name, provider_code, display_name,
    api_base_url, api_version, api_key_env_var,
    supports_chat, supports_vision, supports_function_calling, supports_streaming,
    max_context_window, rate_limit_requests_per_minute,
    default_input_cost_per_1m, default_output_cost_per_1m,
    is_active, is_primary, priority
) VALUES (
    'anthropic', 'anthropic-claude', 'Anthropic Claude',
    'https://api.anthropic.com/v1', 'v1', 'ANTHROPIC_API_KEY',
    true, true, true, true,
    200000, 5000,
    3.00, 15.00,
    true, false, 2
);

-- Insert OpenAI models
INSERT INTO ai_models (
    provider_id, model_name, model_code, display_name,
    supports_vision, supports_function_calling,
    max_context_tokens, max_output_tokens,
    input_cost_per_1m, output_cost_per_1m,
    recommended_for, is_active, is_default
)
SELECT 
    id, 'GPT-4o', 'gpt-4o', 'GPT-4 Optimized',
    true, true,
    128000, 16384,
    5.00, 15.00,
    ARRAY['diagnosis', 'chat', 'vision', 'analysis'], true, true
FROM ai_providers WHERE provider_code = 'openai-gpt';

INSERT INTO ai_models (
    provider_id, model_name, model_code, display_name,
    supports_vision, supports_function_calling,
    max_context_tokens, max_output_tokens,
    input_cost_per_1m, output_cost_per_1m,
    recommended_for, is_active, is_default
)
SELECT 
    id, 'GPT-4-Turbo', 'gpt-4-turbo', 'GPT-4 Turbo',
    true, true,
    128000, 4096,
    10.00, 30.00,
    ARRAY['diagnosis', 'complex_analysis'], true, false
FROM ai_providers WHERE provider_code = 'openai-gpt';

-- Insert Anthropic models
INSERT INTO ai_models (
    provider_id, model_name, model_code, display_name,
    supports_vision, supports_function_calling,
    max_context_tokens, max_output_tokens,
    input_cost_per_1m, output_cost_per_1m,
    recommended_for, is_active, is_default
)
SELECT 
    id, 'Claude 3.5 Sonnet', 'claude-3-5-sonnet-20241022', 'Claude 3.5 Sonnet',
    true, true,
    200000, 8192,
    3.00, 15.00,
    ARRAY['diagnosis', 'chat', 'vision', 'analysis'], true, true
FROM ai_providers WHERE provider_code = 'anthropic-claude';

INSERT INTO ai_models (
    provider_id, model_name, model_code, display_name,
    supports_vision, supports_function_calling,
    max_context_tokens, max_output_tokens,
    input_cost_per_1m, output_cost_per_1m,
    recommended_for, is_active, is_default
)
SELECT 
    id, 'Claude 3 Opus', 'claude-3-opus-20240229', 'Claude 3 Opus',
    true, true,
    200000, 4096,
    15.00, 75.00,
    ARRAY['complex_analysis', 'critical_diagnosis'], true, false
FROM ai_providers WHERE provider_code = 'anthropic-claude';

-- =====================================================================
-- 12. TRIGGERS
-- =====================================================================

CREATE OR REPLACE FUNCTION update_ai_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_ai_providers_updated
    BEFORE UPDATE ON ai_providers
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_timestamp();

CREATE TRIGGER trigger_ai_models_updated
    BEFORE UPDATE ON ai_models
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_timestamp();

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

DO $$
BEGIN
    RAISE NOTICE 'Migration 020 complete!';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - ai_providers (seeded: OpenAI, Anthropic)';
    RAISE NOTICE '  - ai_models (seeded: GPT-4o, Claude 3.5 Sonnet, etc.)';
    RAISE NOTICE '  - ai_conversations';
    RAISE NOTICE '  - ai_diagnosis_results';
    RAISE NOTICE '  - ai_engineer_recommendations';
    RAISE NOTICE '  - ai_parts_recommendations';
    RAISE NOTICE '  - ai_attachment_analysis';
    RAISE NOTICE '  - ai_feedback';
    RAISE NOTICE 'Created 5 helper functions';
    RAISE NOTICE '  Created 2 views';
    RAISE NOTICE 'AI orchestration ready for OpenAI and Anthropic!';
    RAISE NOTICE 'Ready for Phase 2C: AI Services Layer!';
END $$;
