-- Migration: Create comprehensive attachment system with AI analysis support
-- This enables WhatsApp media handling and AI-powered image/video analysis

-- ============================================================================
-- 1. TICKET ATTACHMENTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS ticket_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    
    -- File Information
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_type VARCHAR(100) NOT NULL, -- image/jpeg, video/mp4, application/pdf, etc.
    file_size_bytes BIGINT NOT NULL,
    storage_path TEXT NOT NULL, -- Path/URL to stored file
    
    -- Source Information  
    source VARCHAR(50) NOT NULL, -- whatsapp, web_upload, engineer_app, email
    source_message_id VARCHAR(255), -- WhatsApp message ID if from WhatsApp
    uploaded_by_user_id VARCHAR(32),
    uploaded_by_user_name VARCHAR(255),
    
    -- Categorization
    attachment_category VARCHAR(50) NOT NULL, -- equipment_photo, issue_photo, document, video, audio
    attachment_purpose VARCHAR(100), -- before_repair, after_repair, issue_evidence, parts_needed, etc.
    
    -- AI Analysis Status
    ai_analysis_status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed, skipped
    ai_analysis_requested_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ai_analysis_completed_at TIMESTAMP WITH TIME ZONE,
    ai_analysis_error TEXT,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT attachment_source_check CHECK (source IN ('whatsapp', 'web_upload', 'engineer_app', 'email', 'system')),
    CONSTRAINT attachment_category_check CHECK (attachment_category IN ('equipment_photo', 'issue_photo', 'repair_photo', 'document', 'video', 'audio', 'other')),
    CONSTRAINT attachment_ai_status_check CHECK (ai_analysis_status IN ('pending', 'processing', 'completed', 'failed', 'skipped', 'not_applicable'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_attachments_ticket ON ticket_attachments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_attachments_source ON ticket_attachments(source);
CREATE INDEX IF NOT EXISTS idx_attachments_category ON ticket_attachments(attachment_category);
CREATE INDEX IF NOT EXISTS idx_attachments_ai_status ON ticket_attachments(ai_analysis_status);
CREATE INDEX IF NOT EXISTS idx_attachments_ai_pending ON ticket_attachments(ai_analysis_requested_at) 
    WHERE ai_analysis_status = 'pending';
CREATE INDEX IF NOT EXISTS idx_attachments_created_at ON ticket_attachments(created_at DESC);

-- ============================================================================
-- 2. AI VISION ANALYSIS RESULTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS ai_vision_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES ticket_attachments(id) ON DELETE CASCADE,
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    
    -- AI Provider Information
    ai_provider VARCHAR(50) NOT NULL, -- openai, google_vision, aws_rekognition, custom
    ai_model VARCHAR(100) NOT NULL, -- gpt-4-vision, vision-api-v1, rekognition-v3, etc.
    
    -- Analysis Results
    overall_assessment TEXT, -- General description of what AI sees
    detected_objects JSONB DEFAULT '[]'::jsonb, -- Array of detected objects with confidence
    detected_issues JSONB DEFAULT '[]'::jsonb, -- Array of potential issues identified
    detected_components JSONB DEFAULT '[]'::jsonb, -- Medical equipment components identified
    visible_damage JSONB DEFAULT '[]'::jsonb, -- Damage or wear patterns detected
    text_extraction JSONB DEFAULT '[]'::jsonb, -- OCR results - error codes, labels, etc.
    
    -- Confidence and Quality
    analysis_confidence DECIMAL(5,4), -- 0.0000 to 1.0000
    image_quality_score DECIMAL(5,4), -- 0.0000 to 1.0000  
    analysis_quality VARCHAR(50), -- excellent, good, fair, poor
    
    -- Diagnostic Insights
    equipment_condition_assessment TEXT,
    suggested_focus_areas JSONB DEFAULT '[]'::jsonb, -- Areas that need attention
    repair_recommendations JSONB DEFAULT '[]'::jsonb, -- AI-suggested repair steps
    safety_concerns JSONB DEFAULT '[]'::jsonb, -- Safety issues identified
    
    -- Processing Information
    processing_duration_ms INTEGER, -- How long analysis took
    tokens_used INTEGER, -- For OpenAI pricing tracking
    cost_usd DECIMAL(10,6), -- Estimated cost of this analysis
    
    -- Raw AI Response (for debugging/improvement)
    raw_ai_response JSONB, -- Full response from AI provider
    
    -- Metadata
    analyzed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT analysis_provider_check CHECK (ai_provider IN ('openai', 'google_vision', 'aws_rekognition', 'azure_computer_vision', 'custom')),
    CONSTRAINT analysis_quality_check CHECK (analysis_quality IN ('excellent', 'good', 'fair', 'poor', 'failed'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_vision_analysis_attachment ON ai_vision_analysis(attachment_id);
CREATE INDEX IF NOT EXISTS idx_vision_analysis_ticket ON ai_vision_analysis(ticket_id);  
CREATE INDEX IF NOT EXISTS idx_vision_analysis_provider ON ai_vision_analysis(ai_provider);
CREATE INDEX IF NOT EXISTS idx_vision_analysis_confidence ON ai_vision_analysis(analysis_confidence DESC);
CREATE INDEX IF NOT EXISTS idx_vision_analysis_analyzed_at ON ai_vision_analysis(analyzed_at DESC);

-- GIN indexes for JSONB fields
CREATE INDEX IF NOT EXISTS idx_vision_detected_objects ON ai_vision_analysis USING GIN (detected_objects);
CREATE INDEX IF NOT EXISTS idx_vision_detected_issues ON ai_vision_analysis USING GIN (detected_issues);

-- ============================================================================
-- 3. AI DIAGNOSIS ENHANCEMENTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS ai_diagnosis_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    
    -- Diagnosis Information
    diagnosis_id VARCHAR(100) NOT NULL, -- Links to diagnosis system
    ai_provider VARCHAR(50) NOT NULL,
    ai_model VARCHAR(100) NOT NULL,
    
    -- Input Data Used
    symptoms_analyzed JSONB DEFAULT '[]'::jsonb,
    attachments_analyzed JSONB DEFAULT '[]'::jsonb, -- Array of attachment IDs used
    equipment_context JSONB, -- Equipment details used in analysis
    historical_data_used JSONB DEFAULT '[]'::jsonb, -- Similar cases referenced
    
    -- AI Diagnosis Results
    primary_diagnosis JSONB NOT NULL,
    alternate_diagnoses JSONB DEFAULT '[]'::jsonb,
    confidence_score DECIMAL(5,4) NOT NULL, -- 0.0000 to 1.0000
    confidence_level VARCHAR(20) NOT NULL, -- HIGH, MEDIUM, LOW
    confidence_factors JSONB DEFAULT '[]'::jsonb,
    
    -- Recommendations
    recommended_actions JSONB DEFAULT '[]'::jsonb,
    required_parts JSONB DEFAULT '[]'::jsonb,
    estimated_resolution_time VARCHAR(100),
    
    -- User Feedback and Learning
    user_decision VARCHAR(20), -- accepted, rejected, modified
    user_feedback TEXT,
    actual_diagnosis TEXT, -- What was actually wrong (for learning)
    user_corrections JSONB, -- Specific corrections provided by user
    
    -- Quality and Improvement
    diagnosis_accuracy_score DECIMAL(5,4), -- Post-resolution accuracy assessment
    feedback_provided_at TIMESTAMP WITH TIME ZONE,
    feedback_provided_by VARCHAR(255),
    
    -- Processing Information
    processing_duration_ms INTEGER,
    tokens_used INTEGER,
    cost_usd DECIMAL(10,6),
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT diagnosis_decision_check CHECK (user_decision IN ('accepted', 'rejected', 'modified', 'pending')),
    CONSTRAINT diagnosis_confidence_level_check CHECK (confidence_level IN ('HIGH', 'MEDIUM', 'LOW'))
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_diagnosis_ticket ON ai_diagnosis_results(ticket_id);
CREATE INDEX IF NOT EXISTS idx_diagnosis_confidence ON ai_diagnosis_results(confidence_score DESC);
CREATE INDEX IF NOT EXISTS idx_diagnosis_decision ON ai_diagnosis_results(user_decision);
CREATE INDEX IF NOT EXISTS idx_diagnosis_created_at ON ai_diagnosis_results(created_at DESC);

-- ============================================================================
-- 4. ATTACHMENT PROCESSING QUEUE TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS attachment_processing_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES ticket_attachments(id) ON DELETE CASCADE,
    ticket_id VARCHAR(32) NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    
    -- Processing Information
    processing_type VARCHAR(50) NOT NULL, -- ai_vision, ai_diagnosis, thumbnail, virus_scan, ocr
    priority INTEGER NOT NULL DEFAULT 5, -- 1=highest, 10=lowest
    status VARCHAR(50) NOT NULL DEFAULT 'queued', -- queued, processing, completed, failed, cancelled
    
    -- Processing Configuration
    processing_config JSONB DEFAULT '{}'::jsonb, -- Specific settings for this processing type
    max_retries INTEGER DEFAULT 3,
    retry_count INTEGER DEFAULT 0,
    
    -- Scheduling
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Results
    processing_result JSONB, -- Result data
    error_message TEXT,
    processing_duration_ms INTEGER,
    
    -- Worker Information
    worker_id VARCHAR(100), -- ID of worker/process handling this
    worker_node VARCHAR(255), -- Server/container processing this
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT queue_type_check CHECK (processing_type IN ('ai_vision', 'ai_diagnosis', 'thumbnail', 'virus_scan', 'ocr', 'backup')),
    CONSTRAINT queue_status_check CHECK (status IN ('queued', 'processing', 'completed', 'failed', 'cancelled'))
);

-- Indexes for efficient queue processing
CREATE INDEX IF NOT EXISTS idx_queue_status_priority ON attachment_processing_queue(status, priority, scheduled_at) 
    WHERE status = 'queued';
CREATE INDEX IF NOT EXISTS idx_queue_attachment ON attachment_processing_queue(attachment_id);
CREATE INDEX IF NOT EXISTS idx_queue_processing_type ON attachment_processing_queue(processing_type);
CREATE INDEX IF NOT EXISTS idx_queue_worker ON attachment_processing_queue(worker_id) WHERE status = 'processing';

-- ============================================================================
-- 5. AUTO-UPDATE TRIGGERS
-- ============================================================================

-- Update attachments updated_at
CREATE OR REPLACE FUNCTION update_attachments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER attachments_updated_at_trigger
    BEFORE UPDATE ON ticket_attachments
    FOR EACH ROW
    EXECUTE FUNCTION update_attachments_updated_at();

-- Update diagnosis results updated_at
CREATE TRIGGER diagnosis_updated_at_trigger
    BEFORE UPDATE ON ai_diagnosis_results
    FOR EACH ROW
    EXECUTE FUNCTION update_attachments_updated_at();

-- Update queue updated_at
CREATE TRIGGER queue_updated_at_trigger
    BEFORE UPDATE ON attachment_processing_queue
    FOR EACH ROW
    EXECUTE FUNCTION update_attachments_updated_at();

-- ============================================================================
-- 6. HELPER FUNCTIONS
-- ============================================================================

-- Function: Get pending AI analysis attachments
CREATE OR REPLACE FUNCTION get_pending_ai_analysis()
RETURNS TABLE (
    attachment_id UUID,
    ticket_id VARCHAR,
    filename VARCHAR,
    file_type VARCHAR,
    storage_path TEXT,
    attachment_category VARCHAR,
    hours_pending NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.id,
        a.ticket_id,
        a.filename,
        a.file_type,
        a.storage_path,
        a.attachment_category,
        EXTRACT(EPOCH FROM (NOW() - a.ai_analysis_requested_at)) / 3600 as hours_pending
    FROM ticket_attachments a
    WHERE a.ai_analysis_status = 'pending'
      AND a.file_type LIKE 'image/%' -- Only process images for now
    ORDER BY a.ai_analysis_requested_at ASC;
END;
$$ LANGUAGE plpgsql;

-- Function: Get ticket attachment summary
CREATE OR REPLACE FUNCTION get_ticket_attachment_summary(p_ticket_id VARCHAR)
RETURNS TABLE (
    total_attachments BIGINT,
    images_count BIGINT,
    videos_count BIGINT,
    documents_count BIGINT,
    ai_analyzed_count BIGINT,
    ai_pending_count BIGINT,
    total_file_size_mb NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_attachments,
        COUNT(*) FILTER (WHERE file_type LIKE 'image/%')::BIGINT as images_count,
        COUNT(*) FILTER (WHERE file_type LIKE 'video/%')::BIGINT as videos_count,
        COUNT(*) FILTER (WHERE file_type NOT LIKE 'image/%' AND file_type NOT LIKE 'video/%')::BIGINT as documents_count,
        COUNT(*) FILTER (WHERE ai_analysis_status = 'completed')::BIGINT as ai_analyzed_count,
        COUNT(*) FILTER (WHERE ai_analysis_status = 'pending')::BIGINT as ai_pending_count,
        COALESCE(ROUND(SUM(file_size_bytes) / 1024.0 / 1024.0, 2), 0)::NUMERIC as total_file_size_mb
    FROM ticket_attachments
    WHERE ticket_id = p_ticket_id;
END;
$$ LANGUAGE plpgsql;

-- Function: Queue attachment for AI processing
CREATE OR REPLACE FUNCTION queue_attachment_for_ai_processing(
    p_attachment_id UUID,
    p_processing_type VARCHAR DEFAULT 'ai_vision',
    p_priority INTEGER DEFAULT 5
) RETURNS UUID AS $$
DECLARE
    v_queue_id UUID;
    v_ticket_id VARCHAR;
BEGIN
    -- Get ticket ID
    SELECT ticket_id INTO v_ticket_id 
    FROM ticket_attachments 
    WHERE id = p_attachment_id;
    
    -- Insert into queue
    INSERT INTO attachment_processing_queue (
        attachment_id,
        ticket_id,
        processing_type,
        priority,
        status
    ) VALUES (
        p_attachment_id,
        v_ticket_id,
        p_processing_type,
        p_priority,
        'queued'
    ) RETURNING id INTO v_queue_id;
    
    -- Update attachment status
    UPDATE ticket_attachments 
    SET ai_analysis_status = 'pending',
        ai_analysis_requested_at = NOW()
    WHERE id = p_attachment_id;
    
    RETURN v_queue_id;
END;
$$ LANGUAGE plpgsql;

COMMENT ON TABLE ticket_attachments IS 'Stores all media files attached to service tickets with AI analysis tracking';
COMMENT ON TABLE ai_vision_analysis IS 'AI-powered analysis results for images and videos attached to tickets';
COMMENT ON TABLE ai_diagnosis_results IS 'AI diagnosis results enhanced with visual analysis from attachments';
COMMENT ON TABLE attachment_processing_queue IS 'Queue system for processing attachments with various AI and utility services';