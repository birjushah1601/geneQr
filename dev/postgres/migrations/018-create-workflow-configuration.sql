-- Migration: 018-create-workflow-configuration.sql
-- Description: Configurable Workflow Foundation
-- Ticket: T2B.3
-- Date: 2025-11-16
-- 
-- This migration creates:
-- 1. workflow_templates - Template definitions with hierarchy
-- 2. stage_configuration_templates - Stage templates for workflows
-- 3. ticket_workflow_instances - Active workflow instances per ticket
-- 4. workflow_stage_transitions - Track stage transitions and history
--
-- Purpose:
-- - Enable configurable workflows per equipment/manufacturer
-- - Support Equipment-specific → Manufacturer-specific → Default hierarchy
-- - Track workflow execution and stage transitions
-- - Enable multi-stage service workflows (Remote → Parts → Onsite)

-- =====================================================================
-- 1. WORKFLOW TEMPLATES
-- =====================================================================

CREATE TABLE IF NOT EXISTS workflow_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Template Identity
    template_name TEXT NOT NULL,
    template_code TEXT NOT NULL,                -- Unique code for programmatic reference
    version INT DEFAULT 1,
    
    -- Scope (Hierarchy)
    scope_type TEXT NOT NULL,                   -- 'equipment', 'manufacturer', 'default'
    equipment_catalog_id UUID REFERENCES equipment_catalog(id),
    manufacturer_id UUID REFERENCES organizations(id),
    
    -- Priority for hierarchy resolution
    priority INT NOT NULL,                      -- Equipment=10, Manufacturer=5, Default=1
    
    -- Workflow Configuration
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    
    -- Stage Configuration
    stages_config JSONB NOT NULL DEFAULT '[]',  -- Array of stage configurations
    
    -- Triggers and Conditions
    trigger_conditions JSONB DEFAULT '{}',      -- Conditions to auto-apply this workflow
    
    -- SLA Configuration
    total_sla_hours INT,                        -- Total workflow SLA
    sla_priority TEXT,                          -- 'critical', 'high', 'medium', 'low'
    
    -- Features and Capabilities
    requires_remote_diagnosis BOOLEAN DEFAULT false,
    requires_parts_procurement BOOLEAN DEFAULT false,
    requires_onsite_visit BOOLEAN DEFAULT true,
    requires_followup BOOLEAN DEFAULT false,
    allows_parallel_stages BOOLEAN DEFAULT false,
    
    -- Escalation
    escalation_workflow_id UUID REFERENCES workflow_templates(id),
    escalation_after_hours INT,
    
    -- Approval Requirements
    requires_manager_approval BOOLEAN DEFAULT false,
    requires_client_approval BOOLEAN DEFAULT false,
    
    -- Notifications
    notification_config JSONB DEFAULT '{}',     -- Email/SMS notification settings
    
    -- Validity
    effective_from DATE DEFAULT CURRENT_DATE,
    effective_to DATE,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    updated_by TEXT,
    notes TEXT,
    
    -- Constraints
    CONSTRAINT chk_scope_type CHECK (scope_type IN ('equipment', 'manufacturer', 'default')),
    CONSTRAINT chk_priority CHECK (priority >= 1 AND priority <= 10),
    CONSTRAINT chk_equipment_scope CHECK (
        (scope_type = 'equipment' AND equipment_catalog_id IS NOT NULL) OR
        (scope_type = 'manufacturer' AND manufacturer_id IS NOT NULL AND equipment_catalog_id IS NULL) OR
        (scope_type = 'default' AND equipment_catalog_id IS NULL AND manufacturer_id IS NULL)
    ),
    CONSTRAINT chk_sla_priority CHECK (sla_priority IS NULL OR sla_priority IN (
        'critical', 'high', 'medium', 'low'
    )),
    CONSTRAINT chk_effective_dates CHECK (effective_to IS NULL OR effective_to >= effective_from)
);

-- Unique constraint: One active template per scope
CREATE UNIQUE INDEX idx_workflow_template_unique ON workflow_templates(
    scope_type, 
    COALESCE(equipment_catalog_id::TEXT, 'NULL'),
    COALESCE(manufacturer_id::TEXT, 'NULL'),
    version
) WHERE is_active = true AND effective_to IS NULL;

-- Unique template code
CREATE UNIQUE INDEX idx_workflow_template_code ON workflow_templates(template_code, version);

-- Indexes for hierarchy lookup
CREATE INDEX idx_workflow_templates_equipment ON workflow_templates(equipment_catalog_id) WHERE equipment_catalog_id IS NOT NULL;
CREATE INDEX idx_workflow_templates_manufacturer ON workflow_templates(manufacturer_id) WHERE manufacturer_id IS NOT NULL;
CREATE INDEX idx_workflow_templates_priority ON workflow_templates(priority DESC);
CREATE INDEX idx_workflow_templates_active ON workflow_templates(is_active) WHERE is_active = true;

-- Composite index for hierarchy resolution (most common query)
CREATE INDEX idx_workflow_templates_hierarchy ON workflow_templates(
    equipment_catalog_id, manufacturer_id, priority DESC
) WHERE is_active = true AND effective_to IS NULL;

-- GIN index for JSONB
CREATE INDEX idx_workflow_templates_stages ON workflow_templates USING GIN (stages_config);
CREATE INDEX idx_workflow_templates_triggers ON workflow_templates USING GIN (trigger_conditions);

COMMENT ON TABLE workflow_templates IS 'Configurable workflow templates with Equipment → Manufacturer → Default hierarchy';
COMMENT ON COLUMN workflow_templates.scope_type IS 'equipment (priority 10), manufacturer (priority 5), or default (priority 1)';
COMMENT ON COLUMN workflow_templates.stages_config IS 'JSONB array defining stages in order';
COMMENT ON COLUMN workflow_templates.trigger_conditions IS 'Auto-apply conditions based on ticket properties';

-- =====================================================================
-- 2. STAGE CONFIGURATION TEMPLATES
-- =====================================================================

CREATE TABLE IF NOT EXISTS stage_configuration_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationship
    workflow_template_id UUID NOT NULL REFERENCES workflow_templates(id) ON DELETE CASCADE,
    
    -- Stage Identity
    stage_name TEXT NOT NULL,
    stage_code TEXT NOT NULL,                   -- Unique within workflow
    stage_order INT NOT NULL,                   -- Execution order (1, 2, 3...)
    
    -- Stage Type
    stage_type TEXT NOT NULL,                   -- 'diagnosis', 'parts_procurement', 'onsite', 'followup', 'closure'
    
    -- Configuration
    description TEXT,
    is_optional BOOLEAN DEFAULT false,
    is_parallel BOOLEAN DEFAULT false,          -- Can run parallel with other stages
    
    -- Requirements
    requires_engineer_assignment BOOLEAN DEFAULT true,
    required_support_level TEXT,                -- L1, L2, L3 - minimum level required
    requires_certification BOOLEAN DEFAULT false,
    requires_parts BOOLEAN DEFAULT false,
    requires_client_approval BOOLEAN DEFAULT false,
    
    -- SLA
    sla_hours INT,                              -- Stage-specific SLA
    sla_warning_hours INT,                      -- Warning threshold
    
    -- Completion Criteria
    completion_criteria JSONB DEFAULT '{}',     -- Rules to mark stage complete
    auto_complete BOOLEAN DEFAULT false,
    
    -- Transitions
    next_stage_id UUID REFERENCES stage_configuration_templates(id),
    next_stage_conditions JSONB DEFAULT '{}',   -- Conditions for next stage
    allowed_transitions UUID[],                 -- Array of allowed next stage IDs
    
    -- Actions
    on_start_actions JSONB DEFAULT '[]',        -- Actions when stage starts
    on_complete_actions JSONB DEFAULT '[]',     -- Actions when stage completes
    
    -- Forms and Checklists
    required_forms TEXT[],                      -- Forms that must be filled
    checklist JSONB DEFAULT '[]',               -- Checklist items
    
    -- Attachments
    required_attachments TEXT[],                -- Required attachment types
    optional_attachments TEXT[],
    
    -- Notifications
    notify_on_start BOOLEAN DEFAULT false,
    notify_on_complete BOOLEAN DEFAULT false,
    notification_recipients TEXT[],             -- Roles/users to notify
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_stage_type CHECK (stage_type IN (
        'diagnosis', 'parts_procurement', 'onsite', 'installation', 
        'calibration', 'testing', 'training', 'followup', 'closure'
    )),
    CONSTRAINT chk_stage_order CHECK (stage_order > 0),
    CONSTRAINT chk_support_level CHECK (required_support_level IS NULL OR required_support_level IN ('L1', 'L2', 'L3')),
    CONSTRAINT chk_sla_hours CHECK (sla_hours IS NULL OR sla_hours > 0)
);

-- Unique constraint: stage_code within workflow
CREATE UNIQUE INDEX idx_stage_config_code ON stage_configuration_templates(workflow_template_id, stage_code);

-- Unique constraint: stage_order within workflow
CREATE UNIQUE INDEX idx_stage_config_order ON stage_configuration_templates(workflow_template_id, stage_order);

-- Indexes
CREATE INDEX idx_stage_config_workflow ON stage_configuration_templates(workflow_template_id, stage_order);
CREATE INDEX idx_stage_config_type ON stage_configuration_templates(stage_type);
CREATE INDEX idx_stage_config_next ON stage_configuration_templates(next_stage_id);

-- GIN indexes for JSONB and arrays
CREATE INDEX idx_stage_config_completion ON stage_configuration_templates USING GIN (completion_criteria);
CREATE INDEX idx_stage_config_transitions ON stage_configuration_templates USING GIN (allowed_transitions);

COMMENT ON TABLE stage_configuration_templates IS 'Stage definitions for workflow templates';
COMMENT ON COLUMN stage_configuration_templates.stage_type IS 'Type of stage: diagnosis, parts_procurement, onsite, etc.';
COMMENT ON COLUMN stage_configuration_templates.is_parallel IS 'Whether stage can run parallel with others';
COMMENT ON COLUMN stage_configuration_templates.completion_criteria IS 'Rules to determine stage completion';

-- =====================================================================
-- 3. TICKET WORKFLOW INSTANCES
-- =====================================================================

CREATE TABLE IF NOT EXISTS ticket_workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    workflow_template_id UUID NOT NULL REFERENCES workflow_templates(id),
    
    -- Current State
    current_stage_id UUID REFERENCES stage_configuration_templates(id),
    current_stage_order INT,
    workflow_status TEXT NOT NULL DEFAULT 'pending',  -- 'pending', 'in_progress', 'paused', 'completed', 'cancelled'
    
    -- Timing
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    paused_at TIMESTAMPTZ,
    expected_completion_at TIMESTAMPTZ,
    
    -- SLA Tracking
    total_sla_hours INT,
    hours_elapsed NUMERIC(10,2),
    hours_remaining NUMERIC(10,2),
    is_sla_breached BOOLEAN DEFAULT false,
    sla_breach_reason TEXT,
    
    -- Progress
    total_stages INT NOT NULL,
    completed_stages INT DEFAULT 0,
    progress_percentage NUMERIC(5,2) DEFAULT 0,
    
    -- Configuration Snapshot (at time of creation)
    workflow_config JSONB NOT NULL,            -- Snapshot of workflow template
    stages_config JSONB NOT NULL,              -- Snapshot of stage configurations
    
    -- Modifications
    is_modified BOOLEAN DEFAULT false,         -- Workflow modified after start
    modification_reason TEXT,
    modified_at TIMESTAMPTZ,
    modified_by TEXT,
    
    -- Escalation
    is_escalated BOOLEAN DEFAULT false,
    escalation_workflow_id UUID REFERENCES workflow_templates(id),
    escalated_at TIMESTAMPTZ,
    escalation_reason TEXT,
    
    -- Approvals
    requires_approval BOOLEAN DEFAULT false,
    approved_by TEXT,
    approved_at TIMESTAMPTZ,
    approval_notes TEXT,
    
    -- Metrics
    actual_vs_planned_variance NUMERIC(10,2),  -- Hours difference
    quality_score NUMERIC(3,2),                 -- 0-5 rating
    customer_feedback TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT,
    notes TEXT,
    
    -- Constraints
    CONSTRAINT chk_workflow_status CHECK (workflow_status IN (
        'pending', 'in_progress', 'paused', 'completed', 'cancelled', 'failed'
    )),
    CONSTRAINT chk_completed_stages CHECK (completed_stages >= 0 AND completed_stages <= total_stages),
    CONSTRAINT chk_progress CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    CONSTRAINT chk_quality_score CHECK (quality_score IS NULL OR (quality_score >= 0 AND quality_score <= 5))
);

-- One workflow per ticket
CREATE UNIQUE INDEX idx_workflow_instance_ticket ON ticket_workflow_instances(ticket_id);

-- Indexes
CREATE INDEX idx_workflow_instance_template ON ticket_workflow_instances(workflow_template_id);
CREATE INDEX idx_workflow_instance_status ON ticket_workflow_instances(workflow_status);
CREATE INDEX idx_workflow_instance_stage ON ticket_workflow_instances(current_stage_id);
CREATE INDEX idx_workflow_instance_sla ON ticket_workflow_instances(is_sla_breached) WHERE is_sla_breached = true;
CREATE INDEX idx_workflow_instance_active ON ticket_workflow_instances(workflow_status) 
    WHERE workflow_status IN ('pending', 'in_progress', 'paused');

-- GIN indexes for JSONB
CREATE INDEX idx_workflow_instance_config ON ticket_workflow_instances USING GIN (workflow_config);
CREATE INDEX idx_workflow_instance_stages ON ticket_workflow_instances USING GIN (stages_config);

COMMENT ON TABLE ticket_workflow_instances IS 'Active workflow instances running for service tickets';
COMMENT ON COLUMN ticket_workflow_instances.workflow_config IS 'Immutable snapshot of workflow template at creation time';
COMMENT ON COLUMN ticket_workflow_instances.stages_config IS 'Immutable snapshot of stage configurations';
COMMENT ON COLUMN ticket_workflow_instances.is_modified IS 'Whether workflow was modified after start (deviates from template)';

-- =====================================================================
-- 4. WORKFLOW STAGE TRANSITIONS
-- =====================================================================

CREATE TABLE IF NOT EXISTS workflow_stage_transitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Relationships
    workflow_instance_id UUID NOT NULL REFERENCES ticket_workflow_instances(id),
    ticket_id UUID NOT NULL REFERENCES service_tickets(id),
    
    -- Transition Details
    from_stage_id UUID REFERENCES stage_configuration_templates(id),
    to_stage_id UUID NOT NULL REFERENCES stage_configuration_templates(id),
    from_stage_name TEXT,
    to_stage_name TEXT,
    
    -- Timing
    transitioned_at TIMESTAMPTZ DEFAULT NOW(),
    from_stage_started_at TIMESTAMPTZ,
    from_stage_duration_hours NUMERIC(10,2),
    
    -- Stage Completion
    stage_completed BOOLEAN DEFAULT true,
    completion_percentage NUMERIC(5,2),
    completion_status TEXT,                     -- 'completed', 'skipped', 'failed', 'auto_completed'
    
    -- Assignment
    assigned_engineer_id UUID REFERENCES engineers(id),
    engineer_name TEXT,
    
    -- Actions Taken
    actions_taken TEXT[],
    parts_used UUID[],                          -- References to equipment_parts
    attachments_added INT DEFAULT 0,
    forms_submitted INT DEFAULT 0,
    
    -- Performance
    sla_met BOOLEAN,
    sla_hours INT,
    actual_hours NUMERIC(10,2),
    sla_variance_hours NUMERIC(10,2),
    
    -- Quality
    first_time_completion BOOLEAN DEFAULT true,
    rework_required BOOLEAN DEFAULT false,
    rework_reason TEXT,
    
    -- Transition Reason
    transition_reason TEXT NOT NULL,            -- Why transition occurred
    transition_type TEXT NOT NULL,              -- 'auto', 'manual', 'approval', 'escalation'
    triggered_by TEXT,                          -- User/system that triggered
    
    -- Conditions Met
    conditions_evaluated JSONB DEFAULT '{}',    -- Conditions that were checked
    all_conditions_met BOOLEAN DEFAULT true,
    
    -- Notes and Feedback
    engineer_notes TEXT,
    internal_notes TEXT,
    customer_feedback TEXT,
    
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_completion_status CHECK (completion_status IN (
        'completed', 'skipped', 'failed', 'auto_completed', 'partial'
    )),
    CONSTRAINT chk_transition_type CHECK (transition_type IN (
        'auto', 'manual', 'approval', 'escalation', 'skip', 'fail'
    )),
    CONSTRAINT chk_completion_pct CHECK (completion_percentage IS NULL OR (completion_percentage >= 0 AND completion_percentage <= 100))
);

-- Indexes
CREATE INDEX idx_stage_transition_workflow ON workflow_stage_transitions(workflow_instance_id, transitioned_at DESC);
CREATE INDEX idx_stage_transition_ticket ON workflow_stage_transitions(ticket_id, transitioned_at DESC);
CREATE INDEX idx_stage_transition_from ON workflow_stage_transitions(from_stage_id);
CREATE INDEX idx_stage_transition_to ON workflow_stage_transitions(to_stage_id);
CREATE INDEX idx_stage_transition_engineer ON workflow_stage_transitions(assigned_engineer_id);
CREATE INDEX idx_stage_transition_time ON workflow_stage_transitions(transitioned_at DESC);
CREATE INDEX idx_stage_transition_sla ON workflow_stage_transitions(sla_met) WHERE sla_met = false;

-- GIN indexes
CREATE INDEX idx_stage_transition_actions ON workflow_stage_transitions USING GIN (actions_taken);
CREATE INDEX idx_stage_transition_parts ON workflow_stage_transitions USING GIN (parts_used);
CREATE INDEX idx_stage_transition_conditions ON workflow_stage_transitions USING GIN (conditions_evaluated);

COMMENT ON TABLE workflow_stage_transitions IS 'History of all stage transitions in workflows';
COMMENT ON COLUMN workflow_stage_transitions.transition_type IS 'auto (rule-based), manual (user), approval, escalation';
COMMENT ON COLUMN workflow_stage_transitions.first_time_completion IS 'Whether stage was completed on first attempt';

-- =====================================================================
-- 5. HELPER FUNCTIONS
-- =====================================================================

-- Function: Get workflow template with hierarchy resolution
CREATE OR REPLACE FUNCTION get_workflow_template(
    p_equipment_catalog_id UUID,
    p_manufacturer_id UUID DEFAULT NULL
) RETURNS TABLE (
    template_id UUID,
    template_name TEXT,
    template_code TEXT,
    scope_type TEXT,
    priority INT,
    stages_config JSONB,
    total_sla_hours INT,
    requires_remote_diagnosis BOOLEAN,
    requires_parts_procurement BOOLEAN,
    requires_onsite_visit BOOLEAN
) AS $$
BEGIN
    -- Try equipment-specific first (priority 10)
    RETURN QUERY
    SELECT 
        wt.id,
        wt.template_name,
        wt.template_code,
        wt.scope_type,
        wt.priority,
        wt.stages_config,
        wt.total_sla_hours,
        wt.requires_remote_diagnosis,
        wt.requires_parts_procurement,
        wt.requires_onsite_visit
    FROM workflow_templates wt
    WHERE wt.equipment_catalog_id = p_equipment_catalog_id
      AND wt.is_active = true
      AND wt.effective_to IS NULL
    ORDER BY wt.priority DESC, wt.version DESC
    LIMIT 1;
    
    -- If not found, try manufacturer-level (priority 5)
    IF NOT FOUND AND p_manufacturer_id IS NOT NULL THEN
        RETURN QUERY
        SELECT 
            wt.id,
            wt.template_name,
            wt.template_code,
            wt.scope_type,
            wt.priority,
            wt.stages_config,
            wt.total_sla_hours,
            wt.requires_remote_diagnosis,
            wt.requires_parts_procurement,
            wt.requires_onsite_visit
        FROM workflow_templates wt
        WHERE wt.manufacturer_id = p_manufacturer_id
          AND wt.equipment_catalog_id IS NULL
          AND wt.is_active = true
          AND wt.effective_to IS NULL
        ORDER BY wt.priority DESC, wt.version DESC
        LIMIT 1;
    END IF;
    
    -- If still not found, use default (priority 1)
    IF NOT FOUND THEN
        RETURN QUERY
        SELECT 
            wt.id,
            wt.template_name,
            wt.template_code,
            wt.scope_type,
            wt.priority,
            wt.stages_config,
            wt.total_sla_hours,
            wt.requires_remote_diagnosis,
            wt.requires_parts_procurement,
            wt.requires_onsite_visit
        FROM workflow_templates wt
        WHERE wt.scope_type = 'default'
          AND wt.equipment_catalog_id IS NULL
          AND wt.manufacturer_id IS NULL
          AND wt.is_active = true
          AND wt.effective_to IS NULL
        ORDER BY wt.priority DESC, wt.version DESC
        LIMIT 1;
    END IF;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_workflow_template IS 'Get workflow template with hierarchy: Equipment → Manufacturer → Default';

-- Function: Get stages for workflow template
CREATE OR REPLACE FUNCTION get_workflow_stages(
    p_workflow_template_id UUID
) RETURNS TABLE (
    stage_id UUID,
    stage_name TEXT,
    stage_code TEXT,
    stage_order INT,
    stage_type TEXT,
    description TEXT,
    is_optional BOOLEAN,
    required_support_level TEXT,
    sla_hours INT,
    requires_engineer_assignment BOOLEAN,
    completion_criteria JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        sct.id,
        sct.stage_name,
        sct.stage_code,
        sct.stage_order,
        sct.stage_type,
        sct.description,
        sct.is_optional,
        sct.required_support_level,
        sct.sla_hours,
        sct.requires_engineer_assignment,
        sct.completion_criteria
    FROM stage_configuration_templates sct
    WHERE sct.workflow_template_id = p_workflow_template_id
    ORDER BY sct.stage_order ASC;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_workflow_stages IS 'Get all stages for a workflow template in order';

-- Function: Get current workflow progress
CREATE OR REPLACE FUNCTION get_workflow_progress(
    p_ticket_id UUID
) RETURNS TABLE (
    workflow_id UUID,
    workflow_status TEXT,
    current_stage TEXT,
    progress_percentage NUMERIC,
    completed_stages INT,
    total_stages INT,
    hours_elapsed NUMERIC,
    hours_remaining NUMERIC,
    is_sla_breached BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        twi.id,
        twi.workflow_status,
        sct.stage_name,
        twi.progress_percentage,
        twi.completed_stages,
        twi.total_stages,
        twi.hours_elapsed,
        twi.hours_remaining,
        twi.is_sla_breached
    FROM ticket_workflow_instances twi
    LEFT JOIN stage_configuration_templates sct ON twi.current_stage_id = sct.id
    WHERE twi.ticket_id = p_ticket_id;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_workflow_progress IS 'Get current progress of workflow for a ticket';

-- Function: Get workflow history
CREATE OR REPLACE FUNCTION get_workflow_history(
    p_ticket_id UUID
) RETURNS TABLE (
    transition_id UUID,
    from_stage TEXT,
    to_stage TEXT,
    transitioned_at TIMESTAMPTZ,
    duration_hours NUMERIC,
    completion_status TEXT,
    engineer_name TEXT,
    sla_met BOOLEAN,
    transition_reason TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        wst.id,
        wst.from_stage_name,
        wst.to_stage_name,
        wst.transitioned_at,
        wst.from_stage_duration_hours,
        wst.completion_status,
        wst.engineer_name,
        wst.sla_met,
        wst.transition_reason
    FROM workflow_stage_transitions wst
    WHERE wst.ticket_id = p_ticket_id
    ORDER BY wst.transitioned_at ASC;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION get_workflow_history IS 'Get complete stage transition history for a ticket';

-- Function: Check if stage transition is allowed
CREATE OR REPLACE FUNCTION is_transition_allowed(
    p_workflow_instance_id UUID,
    p_to_stage_id UUID
) RETURNS BOOLEAN AS $$
DECLARE
    v_current_stage_id UUID;
    v_allowed BOOLEAN;
BEGIN
    -- Get current stage
    SELECT current_stage_id INTO v_current_stage_id
    FROM ticket_workflow_instances
    WHERE id = p_workflow_instance_id;
    
    IF v_current_stage_id IS NULL THEN
        -- First stage, always allowed
        RETURN true;
    END IF;
    
    -- Check if transition is in allowed_transitions
    SELECT EXISTS (
        SELECT 1
        FROM stage_configuration_templates
        WHERE id = v_current_stage_id
          AND (next_stage_id = p_to_stage_id 
               OR p_to_stage_id = ANY(allowed_transitions))
    ) INTO v_allowed;
    
    RETURN v_allowed;
END;
$$ LANGUAGE plpgsql STABLE;

COMMENT ON FUNCTION is_transition_allowed IS 'Check if stage transition is allowed per workflow configuration';

-- =====================================================================
-- 6. VIEWS
-- =====================================================================

-- View: Active workflows summary
CREATE OR REPLACE VIEW active_workflows_summary AS
SELECT 
    twi.id as workflow_instance_id,
    twi.ticket_id,
    st.ticket_number,
    st.title as ticket_title,
    wt.template_name as workflow_name,
    sct.stage_name as current_stage,
    twi.workflow_status,
    twi.progress_percentage,
    twi.completed_stages,
    twi.total_stages,
    twi.hours_elapsed,
    twi.hours_remaining,
    twi.is_sla_breached,
    twi.started_at,
    twi.expected_completion_at,
    e.name as assigned_engineer
FROM ticket_workflow_instances twi
JOIN service_tickets st ON twi.ticket_id = st.id
JOIN workflow_templates wt ON twi.workflow_template_id = wt.id
LEFT JOIN stage_configuration_templates sct ON twi.current_stage_id = sct.id
LEFT JOIN engineer_assignments ea ON (st.id = ea.ticket_id AND ea.status = 'active')
LEFT JOIN engineers e ON ea.engineer_id = e.id
WHERE twi.workflow_status IN ('pending', 'in_progress', 'paused');

COMMENT ON VIEW active_workflows_summary IS 'Summary of all active workflows with current status';

-- View: Workflow performance metrics
CREATE OR REPLACE VIEW workflow_performance_metrics AS
SELECT 
    wt.id as workflow_template_id,
    wt.template_name,
    wt.scope_type,
    COUNT(twi.id) as total_executions,
    COUNT(CASE WHEN twi.workflow_status = 'completed' THEN 1 END) as completed_count,
    COUNT(CASE WHEN twi.is_sla_breached = true THEN 1 END) as sla_breach_count,
    ROUND(AVG(twi.progress_percentage), 2) as avg_progress,
    ROUND(AVG(twi.hours_elapsed), 2) as avg_hours_elapsed,
    ROUND(AVG(twi.quality_score), 2) as avg_quality_score,
    ROUND(AVG(twi.actual_vs_planned_variance), 2) as avg_variance_hours
FROM workflow_templates wt
LEFT JOIN ticket_workflow_instances twi ON wt.id = twi.workflow_template_id
WHERE wt.is_active = true
GROUP BY wt.id, wt.template_name, wt.scope_type;

COMMENT ON VIEW workflow_performance_metrics IS 'Performance metrics per workflow template';

-- View: Stage performance metrics
CREATE OR REPLACE VIEW stage_performance_metrics AS
SELECT 
    sct.id as stage_id,
    sct.stage_name,
    sct.stage_type,
    wt.template_name as workflow_name,
    COUNT(wst.id) as total_executions,
    COUNT(CASE WHEN wst.stage_completed = true THEN 1 END) as completed_count,
    COUNT(CASE WHEN wst.sla_met = false THEN 1 END) as sla_miss_count,
    ROUND(AVG(wst.from_stage_duration_hours), 2) as avg_duration_hours,
    ROUND(AVG(wst.sla_variance_hours), 2) as avg_sla_variance,
    COUNT(CASE WHEN wst.first_time_completion = false THEN 1 END) as rework_count
FROM stage_configuration_templates sct
JOIN workflow_templates wt ON sct.workflow_template_id = wt.id
LEFT JOIN workflow_stage_transitions wst ON sct.id = wst.to_stage_id
GROUP BY sct.id, sct.stage_name, sct.stage_type, wt.template_name;

COMMENT ON VIEW stage_performance_metrics IS 'Performance metrics per workflow stage';

-- =====================================================================
-- 7. TRIGGERS
-- =====================================================================

CREATE OR REPLACE FUNCTION update_workflow_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_workflow_template_updated
    BEFORE UPDATE ON workflow_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_timestamp();

CREATE TRIGGER trigger_workflow_instance_updated
    BEFORE UPDATE ON ticket_workflow_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_timestamp();

CREATE TRIGGER trigger_stage_config_updated
    BEFORE UPDATE ON stage_configuration_templates
    FOR EACH ROW
    EXECUTE FUNCTION update_workflow_timestamp();

-- =====================================================================
-- 8. DEFAULT WORKFLOW TEMPLATE
-- =====================================================================

-- Insert default workflow template
DO $$
DECLARE
    v_default_workflow_id UUID;
    v_stage1_id UUID;
    v_stage2_id UUID;
    v_stage3_id UUID;
BEGIN
    -- Create default workflow
    INSERT INTO workflow_templates (
        template_name, template_code, version, scope_type, priority,
        description, total_sla_hours, requires_remote_diagnosis,
        requires_onsite_visit, stages_config
    ) VALUES (
        'Standard Service Workflow',
        'STD-SERVICE-01',
        1,
        'default',
        1,
        'Default 3-stage workflow: Remote Diagnosis → Parts Procurement → Onsite Visit',
        72,  -- 72 hours (3 days)
        true,
        true,
        '[]'::JSONB
    ) RETURNING id INTO v_default_workflow_id;
    
    -- Stage 1: Remote Diagnosis
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_name, stage_code, stage_order, stage_type,
        description, is_optional, required_support_level, sla_hours,
        requires_engineer_assignment, requires_parts, completion_criteria
    ) VALUES (
        v_default_workflow_id,
        'Remote Diagnosis',
        'REMOTE_DIAG',
        1,
        'diagnosis',
        'Initial remote diagnosis to identify issue and required parts',
        false,
        'L1',
        24,
        true,
        false,
        '{"requires_diagnosis_notes": true, "requires_parts_list": true}'::JSONB
    ) RETURNING id INTO v_stage1_id;
    
    -- Stage 2: Parts Procurement
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_name, stage_code, stage_order, stage_type,
        description, is_optional, required_support_level, sla_hours,
        requires_engineer_assignment, requires_parts, completion_criteria,
        next_stage_id
    ) VALUES (
        v_default_workflow_id,
        'Parts Procurement',
        'PARTS_PROC',
        2,
        'parts_procurement',
        'Procure required parts and accessories',
        true,  -- Optional if no parts needed
        'L1',
        24,
        false,  -- Internal team handles
        true,
        '{"parts_available": true}'::JSONB,
        NULL  -- Set later
    ) RETURNING id INTO v_stage2_id;
    
    -- Stage 3: Onsite Visit
    INSERT INTO stage_configuration_templates (
        workflow_template_id, stage_name, stage_code, stage_order, stage_type,
        description, is_optional, required_support_level, sla_hours,
        requires_engineer_assignment, requires_parts, completion_criteria
    ) VALUES (
        v_default_workflow_id,
        'Onsite Visit & Repair',
        'ONSITE_VISIT',
        3,
        'onsite',
        'Onsite visit to perform repair/installation',
        false,
        'L2',
        24,
        true,
        true,
        '{"issue_resolved": true, "equipment_functional": true, "customer_signoff": true}'::JSONB
    ) RETURNING id INTO v_stage3_id;
    
    -- Update next_stage references
    UPDATE stage_configuration_templates SET next_stage_id = v_stage2_id WHERE id = v_stage1_id;
    UPDATE stage_configuration_templates SET next_stage_id = v_stage3_id WHERE id = v_stage2_id;
    
    RAISE NOTICE 'Default workflow template created successfully!';
    RAISE NOTICE 'Workflow ID: %', v_default_workflow_id;
    RAISE NOTICE 'Stages: Remote Diagnosis → Parts Procurement → Onsite Visit';
END $$;

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

-- Verification
DO $$
BEGIN
    RAISE NOTICE 'Migration 018 complete!';
    RAISE NOTICE 'Created tables:';
    RAISE NOTICE '  - workflow_templates';
    RAISE NOTICE '  - stage_configuration_templates';
    RAISE NOTICE '  - ticket_workflow_instances';
    RAISE NOTICE '  - workflow_stage_transitions';
    RAISE NOTICE 'Created 5 helper functions';
    RAISE NOTICE '  Created 3 views';
    RAISE NOTICE 'Created default workflow template with 3 stages';
    RAISE NOTICE 'Ready for multi-stage workflow execution (T2B.4)!';
END $$;
