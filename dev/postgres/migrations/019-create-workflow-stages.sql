-- Migration: 019-create-workflow-stages.sql
-- Description: Multi-Stage Workflow Execution
-- Ticket: T2B.4
-- Date: 2025-11-16
-- 
-- This migration creates:
-- 1. ticket_parts_required - Parts needed per workflow stage
-- 2. stage_assignments - Engineer assignments per stage
-- 3. stage_attachments - Attachments per stage
-- 4. stage_execution_data - Stage-specific execution data
--
-- Purpose:
-- - Track parts requirements per stage
-- - Assign different engineers to different stages
-- - Organize attachments by stage
-- - Store stage-specific execution data (forms, checklists, notes)

-- =====================================================================
-- 1. TICKET PARTS REQUIRED (Per Stage)
-- =====================================================================

CREATE TABLE IF NOT EXISTS ticket_parts_required (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID NOT NULL REFERENCES ticket_workflow_instances(id),
    stage_id UUID REFERENCES stage_configuration_templates(id),
    
    -- Part Information
    equipment_part_id UUID REFERENCES equipment_parts(id),
    part_number TEXT NOT NULL,
    part_name TEXT NOT NULL,
    part_description TEXT,
    
    -- Requirement Details
    required_quantity INT NOT NULL DEFAULT 1,
    recommended_quantity INT,
    available_quantity INT DEFAULT 0,
    
    -- Status
    requirement_status TEXT NOT NULL DEFAULT 'identified',  -- 'identified', 'requested', 'ordered', 'available', 'used', 'returned'
    
    -- Sourcing
    is_oem_part BOOLEAN DEFAULT false,
    supplier_id UUID REFERENCES organizations(id),
    supplier_name TEXT,
    
    -- Procurement
    estimated_price NUMERIC(10,2),
    actual_price NUMERIC(10,2),
    currency TEXT DEFAULT 'INR',
    procurement_lead_time_days INT,
    
    -- Dates
    identified_at TIMESTAMPTZ DEFAULT NOW(),
    requested_at TIMESTAMPTZ,
    ordered_at TIMESTAMPTZ,
    expected_delivery_at TIMESTAMPTZ,
    received_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    
    -- Usage Tracking
    quantity_used INT DEFAULT 0,
    quantity_returned INT DEFAULT 0,
    wastage_quantity INT DEFAULT 0,
    wastage_reason TEXT,
    
    -- Approval
    requires_approval BOOLEAN DEFAULT false,
    approved_by TEXT,
    approved_at TIMESTAMPTZ,
    approval_notes TEXT,
    
    -- Priority
    priority TEXT DEFAULT 'medium',                -- 'critical', 'high', 'medium', 'low'
    is_critical BOOLEAN DEFAULT false,
    
    -- Alternative Parts
    has_alternatives BOOLEAN DEFAULT false,
    alternative_part_ids UUID[],
    
    -- Installation Context
    installation_context TEXT,                     -- From equipment_registry
    context_specific BOOLEAN DEFAULT false,
    
    -- Notes
    notes TEXT,
    internal_notes TEXT,
    procurement_notes TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_requirement_status CHECK (requirement_status IN (
        'identified', 'requested', 'ordered', 'in_transit', 
        'available', 'partially_available', 'used', 'returned', 'cancelled'
    )),
    CONSTRAINT chk_priority CHECK (priority IN ('critical', 'high', 'medium', 'low')),
    CONSTRAINT chk_quantities CHECK (
        required_quantity > 0 AND
        available_quantity >= 0 AND
        quantity_used >= 0 AND
        quantity_returned >= 0 AND
        wastage_quantity >= 0
    )
);

-- Indexes
CREATE INDEX idx_parts_required_ticket ON ticket_parts_required(ticket_id);
CREATE INDEX idx_parts_required_workflow ON ticket_parts_required(workflow_instance_id);
CREATE INDEX idx_parts_required_stage ON ticket_parts_required(stage_id);
CREATE INDEX idx_parts_required_part ON ticket_parts_required(equipment_part_id);
CREATE INDEX idx_parts_required_status ON ticket_parts_required(requirement_status);
CREATE INDEX idx_parts_required_critical ON ticket_parts_required(is_critical) WHERE is_critical = true;
CREATE INDEX idx_parts_required_supplier ON ticket_parts_required(supplier_id);

-- Composite index for stage parts lookup
CREATE INDEX idx_parts_required_stage_status ON ticket_parts_required(
    workflow_instance_id, stage_id, requirement_status
);

-- GIN index for alternative parts
CREATE INDEX idx_parts_required_alternatives ON ticket_parts_required USING GIN (alternative_part_ids);

COMMENT ON TABLE ticket_parts_required IS 'Parts required per workflow stage with procurement tracking';
COMMENT ON COLUMN ticket_parts_required.requirement_status IS 'Lifecycle: identified → requested → ordered → available → used';
COMMENT ON COLUMN ticket_parts_required.context_specific IS 'Whether part is specific to installation context (ICU vs General Ward)';

-- =====================================================================
-- 2. STAGE ASSIGNMENTS (Engineer Assignments Per Stage)
-- =====================================================================

CREATE TABLE IF NOT EXISTS stage_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID NOT NULL REFERENCES ticket_workflow_instances(id),
    stage_id UUID NOT NULL REFERENCES stage_configuration_templates(id),
    
    -- Engineer Assignment
    engineer_id UUID NOT NULL REFERENCES engineers(id),
    
    -- Assignment Details
    assignment_type TEXT NOT NULL DEFAULT 'primary',  -- 'primary', 'secondary', 'support', 'supervisor'
    assignment_status TEXT NOT NULL DEFAULT 'assigned',  -- 'assigned', 'accepted', 'in_progress', 'completed', 'declined', 'reassigned'
    
    -- Stage Context
    stage_name TEXT NOT NULL,
    stage_code TEXT NOT NULL,
    stage_order INT NOT NULL,
    
    -- Requirements Met
    meets_support_level BOOLEAN DEFAULT true,
    required_support_level TEXT,
    engineer_support_level TEXT,
    
    is_certified BOOLEAN DEFAULT false,
    certification_required BOOLEAN DEFAULT false,
    
    -- Timing
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    accepted_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    expected_duration_hours NUMERIC(5,2),
    actual_duration_hours NUMERIC(5,2),
    
    -- Work Details
    work_location TEXT,                            -- 'remote', 'onsite', 'hybrid'
    can_do_remote BOOLEAN DEFAULT true,
    requires_onsite BOOLEAN DEFAULT false,
    
    -- Travel (if onsite)
    requires_travel BOOLEAN DEFAULT false,
    travel_distance_km NUMERIC(8,2),
    estimated_travel_hours NUMERIC(5,2),
    travel_approved BOOLEAN,
    
    -- Availability
    is_available BOOLEAN DEFAULT true,
    availability_checked_at TIMESTAMPTZ,
    conflicts TEXT[],                               -- Array of conflicting assignments
    
    -- Performance
    sla_hours INT,
    sla_met BOOLEAN,
    hours_spent NUMERIC(10,2),
    
    completion_quality TEXT,                        -- 'excellent', 'good', 'acceptable', 'poor'
    first_time_completion BOOLEAN DEFAULT true,
    rework_required BOOLEAN DEFAULT false,
    
    -- Workload
    concurrent_tickets INT DEFAULT 1,
    max_concurrent_allowed INT DEFAULT 5,
    
    -- Replacement/Reassignment
    is_replacement BOOLEAN DEFAULT false,
    replaces_engineer_id UUID REFERENCES engineers(id),
    replacement_reason TEXT,
    reassigned_to_engineer_id UUID REFERENCES engineers(id),
    reassignment_reason TEXT,
    
    -- Notifications
    notification_sent BOOLEAN DEFAULT false,
    notification_sent_at TIMESTAMPTZ,
    acknowledged BOOLEAN DEFAULT false,
    acknowledged_at TIMESTAMPTZ,
    
    -- Notes
    assignment_notes TEXT,
    engineer_notes TEXT,
    internal_notes TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    assigned_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_assignment_type CHECK (assignment_type IN (
        'primary', 'secondary', 'support', 'supervisor', 'reviewer', 'consultant'
    )),
    CONSTRAINT chk_assignment_status CHECK (assignment_status IN (
        'assigned', 'accepted', 'in_progress', 'completed', 'declined', 'reassigned', 'cancelled'
    )),
    CONSTRAINT chk_work_location CHECK (work_location IS NULL OR work_location IN (
        'remote', 'onsite', 'hybrid', 'workshop'
    )),
    CONSTRAINT chk_quality CHECK (completion_quality IS NULL OR completion_quality IN (
        'excellent', 'good', 'acceptable', 'poor', 'incomplete'
    ))
);

-- Indexes
CREATE INDEX idx_stage_assignments_ticket ON stage_assignments(ticket_id);
CREATE INDEX idx_stage_assignments_workflow ON stage_assignments(workflow_instance_id);
CREATE INDEX idx_stage_assignments_stage ON stage_assignments(stage_id);
CREATE INDEX idx_stage_assignments_engineer ON stage_assignments(engineer_id);
CREATE INDEX idx_stage_assignments_status ON stage_assignments(assignment_status);
CREATE INDEX idx_stage_assignments_type ON stage_assignments(assignment_type);

-- Composite index for active assignments
CREATE INDEX idx_stage_assignments_active ON stage_assignments(
    engineer_id, assignment_status
) WHERE assignment_status IN ('assigned', 'accepted', 'in_progress');

-- Composite index for stage lookup
CREATE INDEX idx_stage_assignments_stage_lookup ON stage_assignments(
    workflow_instance_id, stage_id, assignment_type
);

-- GIN index for conflicts
CREATE INDEX idx_stage_assignments_conflicts ON stage_assignments USING GIN (conflicts);

COMMENT ON TABLE stage_assignments IS 'Engineer assignments per workflow stage with detailed tracking';
COMMENT ON COLUMN stage_assignments.assignment_type IS 'primary (main), secondary (backup), support (helper), supervisor (oversight)';
COMMENT ON COLUMN stage_assignments.work_location IS 'Where work is performed: remote, onsite, or hybrid';

-- =====================================================================
-- 3. STAGE ATTACHMENTS (Attachments Per Stage)
-- =====================================================================

CREATE TABLE IF NOT EXISTS stage_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID NOT NULL REFERENCES ticket_workflow_instances(id),
    stage_id UUID REFERENCES stage_configuration_templates(id),
    
    -- Stage Context
    stage_name TEXT,
    stage_code TEXT,
    stage_order INT,
    
    -- Attachment Details
    attachment_name TEXT NOT NULL,
    attachment_type TEXT NOT NULL,                 -- 'image', 'video', 'audio', 'document', 'report'
    file_extension TEXT,
    file_size_bytes BIGINT,
    mime_type TEXT,
    
    -- Storage
    storage_path TEXT NOT NULL,
    storage_url TEXT,
    thumbnail_url TEXT,
    
    -- Classification
    category TEXT,                                  -- 'diagnostic', 'repair', 'completion', 'signoff'
    tags TEXT[],
    is_required BOOLEAN DEFAULT false,
    
    -- Content Description
    description TEXT,
    ai_generated_description TEXT,                 -- AI-generated content description
    contains_equipment_id BOOLEAN,
    contains_serial_number BOOLEAN,
    
    -- Uploaded By
    uploaded_by_engineer_id UUID REFERENCES engineers(id),
    uploaded_by_name TEXT,
    uploaded_by_type TEXT,                         -- 'engineer', 'customer', 'system', 'ai'
    
    -- Timing
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    
    -- AI Analysis
    ai_analyzed BOOLEAN DEFAULT false,
    ai_analysis_results JSONB,
    ai_confidence_score NUMERIC(3,2),
    detected_issues TEXT[],
    
    -- Visibility
    visible_to_customer BOOLEAN DEFAULT false,
    visible_to_engineer BOOLEAN DEFAULT true,
    internal_only BOOLEAN DEFAULT false,
    
    -- Approval/Review
    requires_review BOOLEAN DEFAULT false,
    reviewed BOOLEAN DEFAULT false,
    reviewed_by TEXT,
    reviewed_at TIMESTAMPTZ,
    review_notes TEXT,
    
    -- Compliance
    contains_pii BOOLEAN DEFAULT false,
    gdpr_compliant BOOLEAN DEFAULT true,
    retention_days INT DEFAULT 2555,               -- ~7 years default
    auto_delete_at TIMESTAMPTZ,
    
    -- Quality
    image_quality_score NUMERIC(3,2),              -- For images (0-5)
    is_blurry BOOLEAN,
    is_dark BOOLEAN,
    needs_retake BOOLEAN DEFAULT false,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    
    -- Constraints
    CONSTRAINT chk_attachment_type CHECK (attachment_type IN (
        'image', 'video', 'audio', 'document', 'pdf', 
        'report', 'signature', 'checklist', 'form'
    )),
    CONSTRAINT chk_uploaded_by_type CHECK (uploaded_by_type IN (
        'engineer', 'customer', 'admin', 'system', 'ai'
    )),
    CONSTRAINT chk_file_size CHECK (file_size_bytes > 0)
);

-- Indexes
CREATE INDEX idx_stage_attachments_ticket ON stage_attachments(ticket_id);
CREATE INDEX idx_stage_attachments_workflow ON stage_attachments(workflow_instance_id);
CREATE INDEX idx_stage_attachments_stage ON stage_attachments(stage_id);
CREATE INDEX idx_stage_attachments_type ON stage_attachments(attachment_type);
CREATE INDEX idx_stage_attachments_category ON stage_attachments(category);
CREATE INDEX idx_stage_attachments_uploaded ON stage_attachments(uploaded_at DESC);
CREATE INDEX idx_stage_attachments_engineer ON stage_attachments(uploaded_by_engineer_id);

-- Composite index for stage attachments lookup
CREATE INDEX idx_stage_attachments_stage_lookup ON stage_attachments(
    workflow_instance_id, stage_id, category
);

-- Partial index for AI analysis pending
CREATE INDEX idx_stage_attachments_ai_pending ON stage_attachments(uploaded_at) 
    WHERE ai_analyzed = false AND attachment_type IN ('image', 'video');

-- GIN indexes
CREATE INDEX idx_stage_attachments_tags ON stage_attachments USING GIN (tags);
CREATE INDEX idx_stage_attachments_issues ON stage_attachments USING GIN (detected_issues);
CREATE INDEX idx_stage_attachments_analysis ON stage_attachments USING GIN (ai_analysis_results);

COMMENT ON TABLE stage_attachments IS 'Attachments organized by workflow stage with AI analysis';
COMMENT ON COLUMN stage_attachments.category IS 'diagnostic (before), repair (during), completion (after), signoff (final)';
COMMENT ON COLUMN stage_attachments.ai_analyzed IS 'Whether AI has analyzed this attachment for issues/insights';

-- =====================================================================
-- 4. STAGE EXECUTION DATA (Stage-Specific Data)
-- =====================================================================

CREATE TABLE IF NOT EXISTS stage_execution_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_instance_id UUID NOT NULL REFERENCES ticket_workflow_instances(id),
    stage_id UUID NOT NULL REFERENCES stage_configuration_templates(id),
    
    -- Stage Context
    stage_name TEXT NOT NULL,
    stage_code TEXT NOT NULL,
    stage_order INT NOT NULL,
    
    -- Execution Status
    execution_status TEXT NOT NULL DEFAULT 'pending',  -- 'pending', 'in_progress', 'completed', 'skipped', 'failed'
    
    -- Timing
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_hours NUMERIC(10,2),
    
    -- Forms and Checklists
    forms_data JSONB DEFAULT '{}',                 -- Form submissions
    checklist_data JSONB DEFAULT '[]',             -- Checklist completion
    checklist_completion_pct NUMERIC(5,2),
    
    -- Diagnosis (for diagnosis stage)
    diagnosis_summary TEXT,
    root_cause TEXT,
    recommended_actions TEXT[],
    identified_issues TEXT[],
    
    -- Parts (for parts stage)
    parts_identified UUID[],
    parts_ordered UUID[],
    parts_received UUID[],
    parts_pending INT DEFAULT 0,
    
    -- Work Performed (for onsite/repair stage)
    work_performed TEXT,
    actions_taken TEXT[],
    parts_installed UUID[],
    tools_used TEXT[],
    
    -- Testing and Validation
    tests_performed TEXT[],
    test_results JSONB DEFAULT '{}',
    all_tests_passed BOOLEAN,
    
    -- Customer Interaction
    customer_present BOOLEAN,
    customer_signature_url TEXT,
    customer_feedback TEXT,
    customer_satisfaction_rating INT,              -- 1-5
    
    -- Completion Criteria
    completion_criteria_met JSONB DEFAULT '{}',    -- Which criteria were met
    all_criteria_met BOOLEAN DEFAULT false,
    
    -- Quality Metrics
    quality_checks JSONB DEFAULT '[]',
    quality_score NUMERIC(3,2),
    issues_found INT DEFAULT 0,
    
    -- Follow-up
    requires_followup BOOLEAN DEFAULT false,
    followup_reason TEXT,
    followup_scheduled_at TIMESTAMPTZ,
    
    -- Escalation
    was_escalated BOOLEAN DEFAULT false,
    escalation_reason TEXT,
    escalated_to TEXT,
    
    -- Notes
    engineer_notes TEXT,
    internal_notes TEXT,
    additional_data JSONB DEFAULT '{}',
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    
    -- Constraints
    CONSTRAINT chk_execution_status CHECK (execution_status IN (
        'pending', 'in_progress', 'paused', 'completed', 'skipped', 'failed'
    )),
    CONSTRAINT chk_satisfaction CHECK (customer_satisfaction_rating IS NULL OR (
        customer_satisfaction_rating >= 1 AND customer_satisfaction_rating <= 5
    )),
    CONSTRAINT chk_quality_score CHECK (quality_score IS NULL OR (quality_score >= 0 AND quality_score <= 5))
);

-- One execution record per stage per workflow
CREATE UNIQUE INDEX idx_stage_execution_unique ON stage_execution_data(workflow_instance_id, stage_id);

-- Indexes
CREATE INDEX idx_stage_execution_ticket ON stage_execution_data(ticket_id);
CREATE INDEX idx_stage_execution_workflow ON stage_execution_data(workflow_instance_id);
CREATE INDEX idx_stage_execution_stage ON stage_execution_data(stage_id);
CREATE INDEX idx_stage_execution_status ON stage_execution_data(execution_status);

-- Composite index for stage lookup
CREATE INDEX idx_stage_execution_lookup ON stage_execution_data(
    workflow_instance_id, stage_order
);

-- GIN indexes for JSONB and arrays
CREATE INDEX idx_stage_execution_forms ON stage_execution_data USING GIN (forms_data);
CREATE INDEX idx_stage_execution_checklist ON stage_execution_data USING GIN (checklist_data);
CREATE INDEX idx_stage_execution_issues ON stage_execution_data USING GIN (identified_issues);
CREATE INDEX idx_stage_execution_actions ON stage_execution_data USING GIN (actions_taken);
CREATE INDEX idx_stage_execution_tests ON stage_execution_data USING GIN (test_results);

COMMENT ON TABLE stage_execution_data IS 'Stage-specific execution data including forms, checklists, and work performed';
COMMENT ON COLUMN stage_execution_data.forms_data IS 'Dynamic form submissions specific to stage type';
COMMENT ON COLUMN stage_execution_data.completion_criteria_met IS 'Tracks which completion criteria were satisfied';

-- =====================================================================
-- 5. HELPER FUNCTIONS
-- =====================================================================

-- Function: Get parts needed for stage
CREATE OR REPLACE FUNCTION get_stage_parts(
    p_workflow_instance_id UUID,
    p_stage_id UUID
) RETURNS TABLE (
    part_id UUID,
    part_number TEXT,
    part_name TEXT,
    required_quantity INT,
    available_quantity INT,
    requirement_status TEXT,
    is_critical BOOLEAN,
    estimated_price NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        tpr.id,
        tpr.part_number,
        tpr.part_name,
        tpr.required_quantity,
        tpr.available_quantity,
        tpr.requirement_status,
        tpr.is_critical,
        tpr.estimated_price
    FROM ticket_parts_required tpr
    WHERE tpr.workflow_instance_id = p_workflow_instance_id
      AND (p_stage_id IS NULL OR tpr.stage_id = p_stage_id)
    ORDER BY tpr.is_critical DESC, tpr.priority DESC, tpr.part_name;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_stage_parts IS 'Get all parts required for a stage or entire workflow';

-- Function: Get engineers assigned to stage
CREATE OR REPLACE FUNCTION get_stage_engineers(
    p_workflow_instance_id UUID,
    p_stage_id UUID
) RETURNS TABLE (
    assignment_id UUID,
    engineer_id UUID,
    engineer_name TEXT,
    assignment_type TEXT,
    assignment_status TEXT,
    work_location TEXT,
    started_at TIMESTAMPTZ,
    hours_spent NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        sa.id,
        sa.engineer_id,
        e.name,
        sa.assignment_type,
        sa.assignment_status,
        sa.work_location,
        sa.started_at,
        sa.hours_spent
    FROM stage_assignments sa
    JOIN engineers e ON sa.engineer_id = e.id
    WHERE sa.workflow_instance_id = p_workflow_instance_id
      AND sa.stage_id = p_stage_id
    ORDER BY 
        CASE sa.assignment_type
            WHEN 'primary' THEN 1
            WHEN 'secondary' THEN 2
            WHEN 'support' THEN 3
            ELSE 4
        END;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_stage_engineers IS 'Get all engineers assigned to a specific stage';

-- Function: Get stage attachments
CREATE OR REPLACE FUNCTION get_stage_attachments(
    p_workflow_instance_id UUID,
    p_stage_id UUID,
    p_category TEXT DEFAULT NULL
) RETURNS TABLE (
    attachment_id UUID,
    attachment_name TEXT,
    attachment_type TEXT,
    category TEXT,
    storage_url TEXT,
    thumbnail_url TEXT,
    uploaded_by TEXT,
    uploaded_at TIMESTAMPTZ,
    ai_analyzed BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        sa.id,
        sa.attachment_name,
        sa.attachment_type,
        sa.category,
        sa.storage_url,
        sa.thumbnail_url,
        sa.uploaded_by_name,
        sa.uploaded_at,
        sa.ai_analyzed
    FROM stage_attachments sa
    WHERE sa.workflow_instance_id = p_workflow_instance_id
      AND (p_stage_id IS NULL OR sa.stage_id = p_stage_id)
      AND (p_category IS NULL OR sa.category = p_category)
    ORDER BY sa.uploaded_at DESC;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_stage_attachments IS 'Get attachments for a stage, optionally filtered by category';

-- Function: Get stage execution summary
CREATE OR REPLACE FUNCTION get_stage_execution_summary(
    p_workflow_instance_id UUID
) RETURNS TABLE (
    stage_order INT,
    stage_name TEXT,
    execution_status TEXT,
    duration_hours NUMERIC,
    engineers_assigned INT,
    parts_required INT,
    parts_available INT,
    attachments_count INT,
    completion_pct NUMERIC,
    all_criteria_met BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        sed.stage_order,
        sed.stage_name,
        sed.execution_status,
        sed.duration_hours,
        (SELECT COUNT(*) FROM stage_assignments sa 
         WHERE sa.workflow_instance_id = p_workflow_instance_id 
         AND sa.stage_id = sed.stage_id),
        (SELECT COUNT(*) FROM ticket_parts_required tpr 
         WHERE tpr.workflow_instance_id = p_workflow_instance_id 
         AND tpr.stage_id = sed.stage_id),
        (SELECT COUNT(*) FROM ticket_parts_required tpr 
         WHERE tpr.workflow_instance_id = p_workflow_instance_id 
         AND tpr.stage_id = sed.stage_id 
         AND tpr.requirement_status = 'available'),
        (SELECT COUNT(*) FROM stage_attachments sa 
         WHERE sa.workflow_instance_id = p_workflow_instance_id 
         AND sa.stage_id = sed.stage_id),
        sed.checklist_completion_pct,
        sed.all_criteria_met
    FROM stage_execution_data sed
    WHERE sed.workflow_instance_id = p_workflow_instance_id
    ORDER BY sed.stage_order;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_stage_execution_summary IS 'Get complete execution summary for all stages in workflow';

-- Function: Check if stage can start
CREATE OR REPLACE FUNCTION can_stage_start(
    p_workflow_instance_id UUID,
    p_stage_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    v_has_engineer BOOLEAN;
    v_has_parts BOOLEAN;
    v_parts_required BOOLEAN;
    v_can_start BOOLEAN := true;
BEGIN
    -- Check if engineer is assigned
    SELECT EXISTS (
        SELECT 1 FROM stage_assignments
        WHERE workflow_instance_id = p_workflow_instance_id
          AND stage_id = p_stage_id
          AND assignment_status IN ('assigned', 'accepted')
    ) INTO v_has_engineer;
    
    -- Check if parts are required
    SELECT EXISTS (
        SELECT 1 FROM stage_configuration_templates
        WHERE id = p_stage_id AND requires_parts = true
    ) INTO v_parts_required;
    
    IF v_parts_required THEN
        -- Check if all required parts are available
        SELECT NOT EXISTS (
            SELECT 1 FROM ticket_parts_required
            WHERE workflow_instance_id = p_workflow_instance_id
              AND stage_id = p_stage_id
              AND is_critical = true
              AND requirement_status != 'available'
        ) INTO v_has_parts;
        
        v_can_start := v_has_engineer AND v_has_parts;
    ELSE
        v_can_start := v_has_engineer;
    END IF;
    
    RETURN v_can_start;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION can_stage_start IS 'Check if stage has all prerequisites met to start';

-- =====================================================================
-- 6. VIEWS
-- =====================================================================

-- View: Stage execution overview
CREATE OR REPLACE VIEW stage_execution_overview AS
SELECT 
    sed.workflow_instance_id,
    sed.ticket_id,
    st.ticket_number,
    sed.stage_name,
    sed.stage_order,
    sed.execution_status,
    sed.duration_hours,
    (SELECT COUNT(*) FROM stage_assignments sa 
     WHERE sa.workflow_instance_id = sed.workflow_instance_id 
     AND sa.stage_id = sed.stage_id 
     AND sa.assignment_status IN ('assigned', 'accepted', 'in_progress')) as active_engineers,
    (SELECT COUNT(*) FROM ticket_parts_required tpr 
     WHERE tpr.workflow_instance_id = sed.workflow_instance_id 
     AND tpr.stage_id = sed.stage_id) as total_parts,
    (SELECT COUNT(*) FROM ticket_parts_required tpr 
     WHERE tpr.workflow_instance_id = sed.workflow_instance_id 
     AND tpr.stage_id = sed.stage_id 
     AND tpr.requirement_status = 'available') as available_parts,
    (SELECT COUNT(*) FROM stage_attachments sa 
     WHERE sa.workflow_instance_id = sed.workflow_instance_id 
     AND sa.stage_id = sed.stage_id) as attachments_count,
    sed.checklist_completion_pct,
    sed.all_criteria_met,
    sed.customer_satisfaction_rating
FROM stage_execution_data sed
JOIN service_tickets st ON sed.ticket_id = st.id;

COMMENT ON VIEW stage_execution_overview IS 'Overview of stage execution with counts of related data';

-- View: Parts procurement status
CREATE OR REPLACE VIEW parts_procurement_status AS
SELECT 
    tpr.ticket_id,
    st.ticket_number,
    tpr.workflow_instance_id,
    sct.stage_name,
    COUNT(*) as total_parts,
    COUNT(CASE WHEN tpr.is_critical THEN 1 END) as critical_parts,
    COUNT(CASE WHEN tpr.requirement_status = 'available' THEN 1 END) as available_parts,
    COUNT(CASE WHEN tpr.requirement_status IN ('identified', 'requested') THEN 1 END) as pending_parts,
    SUM(tpr.estimated_price) as estimated_total_cost,
    SUM(tpr.actual_price) as actual_total_cost
FROM ticket_parts_required tpr
JOIN service_tickets st ON tpr.ticket_id = st.id
LEFT JOIN stage_configuration_templates sct ON tpr.stage_id = sct.id
GROUP BY tpr.ticket_id, st.ticket_number, tpr.workflow_instance_id, sct.stage_name;

COMMENT ON VIEW parts_procurement_status IS 'Parts procurement status summary per ticket';

-- =====================================================================
-- 7. TRIGGERS
-- =====================================================================

CREATE OR REPLACE FUNCTION update_stage_data_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_parts_required_updated
    BEFORE UPDATE ON ticket_parts_required
    FOR EACH ROW
    EXECUTE FUNCTION update_stage_data_timestamp();

CREATE TRIGGER trigger_stage_assignments_updated
    BEFORE UPDATE ON stage_assignments
    FOR EACH ROW
    EXECUTE FUNCTION update_stage_data_timestamp();

CREATE TRIGGER trigger_stage_attachments_updated
    BEFORE UPDATE ON stage_attachments
    FOR EACH ROW
    EXECUTE FUNCTION update_stage_data_timestamp();

CREATE TRIGGER trigger_stage_execution_updated
    BEFORE UPDATE ON stage_execution_data
    FOR EACH ROW
    EXECUTE FUNCTION update_stage_data_timestamp();

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

DO $$
BEGIN
    RAISE NOTICE 'Migration 019 complete!';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - ticket_parts_required';
    RAISE NOTICE '  - stage_assignments';
    RAISE NOTICE '  - stage_attachments';
    RAISE NOTICE '  - stage_execution_data';
    RAISE NOTICE 'Created 5 helper functions';
    RAISE NOTICE 'Created 2 views';
    RAISE NOTICE 'Ready for AI Integration Schema (T2B.5)!';
END $$;
