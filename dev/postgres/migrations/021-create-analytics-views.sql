-- Migration: 021-create-analytics-views.sql
-- Description: Workflow Analytics & Monitoring Views
-- Ticket: T2B.6
-- Date: 2025-11-16
-- 
-- This migration creates comprehensive analytics and monitoring views:
-- 1. Ticket Performance Metrics
-- 2. Stage Performance Tracking
-- 3. Engineer Performance Analytics
-- 4. Parts Usage Analytics
-- 5. AI Performance Metrics
-- 6. SLA Compliance Monitoring
-- 7. Cost Analytics
-- 8. Workflow Efficiency Metrics
-- 9. Customer Satisfaction Tracking
-- 10. Real-time Dashboards

-- =====================================================================
-- 1. TICKET PERFORMANCE METRICS
-- =====================================================================

-- View: Ticket Overview with Complete Metrics
CREATE OR REPLACE VIEW ticket_performance_overview AS
SELECT 
    st.id as ticket_id,
    st.ticket_number,
    st.title,
    st.status,
    st.priority,
    st.severity,
    
    -- Equipment Info
    e.equipment_name,
    e.model_number,
    m.name as manufacturer_name,
    
    -- Customer Info
    c.name as customer_name,
    c.organization_type,
    
    -- Workflow Info
    wt.template_name as workflow_name,
    twi.current_stage_number,
    twi.total_stages,
    ROUND((twi.current_stage_number::NUMERIC / twi.total_stages) * 100, 2) as completion_percentage,
    
    -- Timing
    st.created_at,
    st.updated_at,
    EXTRACT(EPOCH FROM (COALESCE(st.resolved_at, NOW()) - st.created_at)) / 3600 as hours_open,
    EXTRACT(EPOCH FROM (st.resolved_at - st.created_at)) / 3600 as hours_to_resolve,
    st.sla_due_date,
    CASE 
        WHEN st.resolved_at IS NOT NULL AND st.resolved_at <= st.sla_due_date THEN 'Met'
        WHEN st.resolved_at IS NOT NULL AND st.resolved_at > st.sla_due_date THEN 'Breached'
        WHEN NOW() > st.sla_due_date THEN 'Breaching'
        ELSE 'On Track'
    END as sla_status,
    
    -- AI Assistance
    EXISTS(SELECT 1 FROM ai_diagnosis_results WHERE ticket_id = st.id) as has_ai_diagnosis,
    (SELECT confidence_score FROM ai_diagnosis_results WHERE ticket_id = st.id) as ai_confidence,
    (SELECT was_accurate FROM ai_diagnosis_results WHERE ticket_id = st.id) as ai_was_accurate,
    
    -- Assignment
    (SELECT COUNT(*) FROM stage_assignments WHERE ticket_id = st.id) as total_assignments,
    (SELECT COUNT(DISTINCT engineer_id) FROM stage_assignments WHERE ticket_id = st.id) as unique_engineers,
    
    -- Parts
    (SELECT COUNT(*) FROM ticket_parts_required WHERE ticket_id = st.id) as total_parts_required,
    (SELECT COUNT(*) FROM ticket_parts_required WHERE ticket_id = st.id AND procurement_status = 'used') as parts_used,
    (SELECT SUM(actual_price * quantity_used) FROM ticket_parts_required WHERE ticket_id = st.id) as parts_cost,
    
    -- Attachments
    (SELECT COUNT(*) FROM stage_attachments WHERE ticket_id = st.id) as total_attachments,
    
    -- Customer Satisfaction
    st.customer_satisfaction_rating,
    st.resolution_notes

FROM service_tickets st
LEFT JOIN equipment e ON st.equipment_id = e.id
LEFT JOIN manufacturers m ON e.manufacturer_id = m.id
LEFT JOIN customers c ON st.customer_id = c.id
LEFT JOIN ticket_workflow_instances twi ON st.id = twi.ticket_id
LEFT JOIN workflow_templates wt ON twi.template_id = wt.id;

CREATE INDEX idx_ticket_perf_status ON service_tickets(status);
CREATE INDEX idx_ticket_perf_priority ON service_tickets(priority);
CREATE INDEX idx_ticket_perf_created ON service_tickets(created_at DESC);

COMMENT ON VIEW ticket_performance_overview IS 'Complete ticket metrics with workflow, AI, parts, and SLA tracking';

-- =====================================================================
-- 2. STAGE PERFORMANCE TRACKING
-- =====================================================================

-- View: Stage Execution Performance
CREATE OR REPLACE VIEW stage_performance_metrics AS
SELECT 
    sct.id as stage_id,
    sct.stage_name,
    sct.stage_type,
    wt.template_name as workflow_name,
    
    -- Volume
    COUNT(DISTINCT wst.id) as total_executions,
    COUNT(DISTINCT CASE WHEN wst.status = 'completed' THEN wst.id END) as completed_count,
    COUNT(DISTINCT CASE WHEN wst.status = 'in_progress' THEN wst.id END) as in_progress_count,
    COUNT(DISTINCT CASE WHEN wst.status = 'blocked' THEN wst.id END) as blocked_count,
    
    -- Timing
    ROUND(AVG(EXTRACT(EPOCH FROM (wst.completed_at - wst.started_at)) / 3600), 2) as avg_hours_to_complete,
    ROUND(MIN(EXTRACT(EPOCH FROM (wst.completed_at - wst.started_at)) / 3600), 2) as min_hours,
    ROUND(MAX(EXTRACT(EPOCH FROM (wst.completed_at - wst.started_at)) / 3600), 2) as max_hours,
    
    -- SLA Performance
    COUNT(CASE WHEN wst.completed_at <= sct.target_duration_hours * INTERVAL '1 hour' + wst.started_at THEN 1 END) as sla_met_count,
    ROUND(
        (COUNT(CASE WHEN wst.completed_at <= sct.target_duration_hours * INTERVAL '1 hour' + wst.started_at THEN 1 END)::NUMERIC / 
         NULLIF(COUNT(CASE WHEN wst.status = 'completed' THEN 1 END), 0)) * 100, 
        2
    ) as sla_compliance_pct,
    
    -- Engineer Assignments
    COUNT(DISTINCT sa.engineer_id) as unique_engineers_assigned,
    ROUND(AVG((SELECT COUNT(*) FROM stage_assignments sa2 WHERE sa2.workflow_stage_id = wst.id)), 2) as avg_engineers_per_execution,
    
    -- Parts
    ROUND(AVG((SELECT COUNT(*) FROM ticket_parts_required tpr WHERE tpr.ticket_id = wst.ticket_id AND tpr.required_for_stage = sct.id)), 2) as avg_parts_per_execution,
    
    -- Quality
    ROUND(AVG(sed.quality_score), 2) as avg_quality_score,
    ROUND(AVG(sed.customer_satisfaction_rating), 2) as avg_customer_rating,
    
    -- Current Active
    COUNT(CASE WHEN wst.status = 'in_progress' AND wst.started_at < NOW() - INTERVAL '24 hours' THEN 1 END) as stuck_over_24h

FROM stage_configuration_templates sct
JOIN workflow_templates wt ON sct.workflow_template_id = wt.id
LEFT JOIN workflow_stage_transitions wst ON sct.id = wst.stage_id
LEFT JOIN stage_assignments sa ON wst.id = sa.workflow_stage_id
LEFT JOIN stage_execution_data sed ON wst.id = sed.workflow_stage_id
GROUP BY sct.id, sct.stage_name, sct.stage_type, wt.template_name;

COMMENT ON VIEW stage_performance_metrics IS 'Performance metrics per workflow stage with SLA compliance';

-- =====================================================================
-- 3. ENGINEER PERFORMANCE ANALYTICS
-- =====================================================================

-- View: Engineer Performance Dashboard
CREATE OR REPLACE VIEW engineer_performance_dashboard AS
SELECT 
    e.id as engineer_id,
    e.name as engineer_name,
    e.email,
    e.status as engineer_status,
    
    -- Expertise Summary
    (SELECT COUNT(*) FROM engineer_equipment_expertise WHERE engineer_id = e.id) as equipment_expertise_count,
    (SELECT COUNT(*) FROM engineer_equipment_expertise WHERE engineer_id = e.id AND support_level = 'L3') as l3_expertise_count,
    (SELECT COUNT(*) FROM engineer_certifications WHERE engineer_id = e.id AND is_active = true) as active_certifications,
    
    -- Workload
    (SELECT COUNT(*) FROM stage_assignments 
     WHERE engineer_id = e.id 
     AND assignment_status IN ('assigned', 'in_progress', 'traveling')) as current_active_tickets,
    
    (SELECT COUNT(*) FROM stage_assignments 
     WHERE engineer_id = e.id 
     AND assignment_status = 'completed'
     AND completed_at >= NOW() - INTERVAL '30 days') as tickets_last_30_days,
    
    -- Performance Metrics
    ROUND(AVG(sa.customer_rating), 2) as avg_customer_rating,
    ROUND(AVG(sa.performance_score), 2) as avg_performance_score,
    
    -- Timing
    ROUND(AVG(EXTRACT(EPOCH FROM (sa.completed_at - sa.assigned_at)) / 3600), 2) as avg_completion_hours,
    
    -- AI Recommendations
    (SELECT COUNT(*) FROM ai_engineer_recommendations 
     WHERE engineer_id = e.id 
     AND created_at >= NOW() - INTERVAL '30 days') as ai_recommended_count,
    
    (SELECT COUNT(*) FROM ai_engineer_recommendations 
     WHERE engineer_id = e.id 
     AND was_selected = true
     AND created_at >= NOW() - INTERVAL '30 days') as ai_selected_count,
    
    ROUND(
        (SELECT COUNT(*)::NUMERIC FROM ai_engineer_recommendations 
         WHERE engineer_id = e.id AND was_selected = true AND created_at >= NOW() - INTERVAL '30 days') /
        NULLIF((SELECT COUNT(*) FROM ai_engineer_recommendations 
                WHERE engineer_id = e.id AND created_at >= NOW() - INTERVAL '30 days'), 0) * 100,
        2
    ) as ai_selection_rate_pct,
    
    -- SLA Performance
    (SELECT COUNT(*) FROM stage_assignments sa2
     JOIN workflow_stage_transitions wst ON sa2.workflow_stage_id = wst.id
     JOIN stage_configuration_templates sct ON wst.stage_id = sct.id
     WHERE sa2.engineer_id = e.id 
     AND sa2.assignment_status = 'completed'
     AND sa2.completed_at <= sct.target_duration_hours * INTERVAL '1 hour' + sa2.assigned_at) as sla_met_count_30d,
    
    -- Availability
    e.is_available,
    e.current_location

FROM engineers e
LEFT JOIN stage_assignments sa ON e.id = sa.engineer_id AND sa.assignment_status = 'completed';

COMMENT ON VIEW engineer_performance_dashboard IS 'Comprehensive engineer performance with AI recommendations and SLA compliance';

-- =====================================================================
-- 4. PARTS USAGE ANALYTICS
-- =====================================================================

-- View: Parts Usage and Cost Analysis
CREATE OR REPLACE VIEW parts_usage_analytics AS
SELECT 
    ep.id as part_id,
    ep.part_number,
    ep.part_name,
    ep.part_category,
    ec.catalog_name as equipment_catalog,
    m.name as manufacturer_name,
    
    -- Usage Statistics
    COUNT(DISTINCT tpr.ticket_id) as tickets_used_in,
    SUM(tpr.quantity_used) as total_quantity_used,
    SUM(tpr.quantity_returned) as total_quantity_returned,
    SUM(tpr.quantity_wastage) as total_quantity_wastage,
    
    -- Financial
    SUM(tpr.estimated_price * tpr.quantity_identified) as total_estimated_cost,
    SUM(tpr.actual_price * tpr.quantity_used) as total_actual_cost,
    ROUND(AVG(tpr.actual_price), 2) as avg_actual_price,
    
    -- Procurement Performance
    ROUND(AVG(EXTRACT(EPOCH FROM (tpr.received_at - tpr.ordered_at)) / 86400), 2) as avg_lead_time_days,
    COUNT(CASE WHEN tpr.procurement_status = 'used' THEN 1 END) as successfully_used_count,
    
    -- AI Recommendations
    (SELECT COUNT(*) FROM ai_parts_recommendations 
     WHERE equipment_part_id = ep.id 
     AND created_at >= NOW() - INTERVAL '30 days') as ai_recommended_count_30d,
    
    (SELECT COUNT(*) FROM ai_parts_recommendations 
     WHERE equipment_part_id = ep.id 
     AND was_used = true
     AND created_at >= NOW() - INTERVAL '30 days') as ai_correct_count_30d,
    
    ROUND(
        (SELECT COUNT(*)::NUMERIC FROM ai_parts_recommendations 
         WHERE equipment_part_id = ep.id AND was_used = true AND created_at >= NOW() - INTERVAL '30 days') /
        NULLIF((SELECT COUNT(*) FROM ai_parts_recommendations 
                WHERE equipment_part_id = ep.id AND created_at >= NOW() - INTERVAL '30 days'), 0) * 100,
        2
    ) as ai_accuracy_pct,
    
    -- Context Usage
    (SELECT COUNT(DISTINCT installation_context) FROM ticket_parts_required WHERE part_id = ep.id) as used_in_contexts,
    
    -- Inventory Status
    ep.is_oem_part,
    ep.stock_status

FROM equipment_parts ep
JOIN equipment_catalog ec ON ep.equipment_catalog_id = ec.id
LEFT JOIN manufacturers m ON ec.manufacturer_id = m.id
LEFT JOIN ticket_parts_required tpr ON ep.id = tpr.part_id
GROUP BY ep.id, ep.part_number, ep.part_name, ep.part_category, ec.catalog_name, m.name, ep.is_oem_part, ep.stock_status;

COMMENT ON VIEW parts_usage_analytics IS 'Parts usage statistics with cost analysis and AI accuracy';

-- =====================================================================
-- 5. AI PERFORMANCE METRICS (Enhanced)
-- =====================================================================

-- View: AI Diagnosis Performance by Equipment
CREATE OR REPLACE VIEW ai_diagnosis_performance_by_equipment AS
SELECT 
    ec.catalog_name as equipment_type,
    m.name as manufacturer_name,
    adr.provider_name,
    adr.model_name,
    
    -- Volume
    COUNT(*) as total_diagnoses,
    COUNT(CASE WHEN adr.is_validated THEN 1 END) as validated_count,
    COUNT(CASE WHEN adr.was_accurate = true THEN 1 END) as accurate_count,
    
    -- Accuracy
    ROUND(
        (COUNT(CASE WHEN adr.was_accurate = true THEN 1 END)::NUMERIC / 
         NULLIF(COUNT(CASE WHEN adr.is_validated THEN 1 END), 0)) * 100, 
        2
    ) as accuracy_rate_pct,
    
    -- Confidence
    ROUND(AVG(adr.confidence_score), 2) as avg_confidence_score,
    ROUND(AVG(CASE WHEN adr.was_accurate = true THEN adr.confidence_score END), 2) as avg_confidence_when_correct,
    ROUND(AVG(CASE WHEN adr.was_accurate = false THEN adr.confidence_score END), 2) as avg_confidence_when_wrong,
    
    -- Performance
    ROUND(AVG(adr.processing_time_ms), 0) as avg_processing_ms,
    
    -- Cost
    SUM(adr.cost_usd) as total_cost_usd,
    ROUND(AVG(adr.cost_usd), 4) as avg_cost_per_diagnosis,
    
    -- Severity Distribution
    COUNT(CASE WHEN adr.severity_level = 'critical' THEN 1 END) as critical_count,
    COUNT(CASE WHEN adr.severity_level = 'high' THEN 1 END) as high_count,
    
    -- Support Level Recommendations
    COUNT(CASE WHEN adr.recommended_support_level = 'L3' THEN 1 END) as l3_recommended_count

FROM ai_diagnosis_results adr
JOIN service_tickets st ON adr.ticket_id = st.id
JOIN equipment e ON st.equipment_id = e.id
JOIN equipment_catalog ec ON e.catalog_id = ec.id
LEFT JOIN manufacturers m ON ec.manufacturer_id = m.id
GROUP BY ec.catalog_name, m.name, adr.provider_name, adr.model_name;

COMMENT ON VIEW ai_diagnosis_performance_by_equipment IS 'AI diagnosis accuracy and performance by equipment type';

-- =====================================================================
-- 6. SLA COMPLIANCE MONITORING
-- =====================================================================

-- View: SLA Compliance Dashboard
CREATE OR REPLACE VIEW sla_compliance_dashboard AS
SELECT 
    DATE_TRUNC('day', st.created_at) as date,
    st.priority,
    st.severity,
    c.organization_type as customer_type,
    ec.catalog_name as equipment_type,
    
    -- Volume
    COUNT(*) as total_tickets,
    COUNT(CASE WHEN st.status = 'resolved' THEN 1 END) as resolved_count,
    COUNT(CASE WHEN st.status IN ('open', 'in_progress', 'assigned') THEN 1 END) as open_count,
    
    -- SLA Status
    COUNT(CASE 
        WHEN st.resolved_at IS NOT NULL AND st.resolved_at <= st.sla_due_date THEN 1 
    END) as sla_met_count,
    
    COUNT(CASE 
        WHEN st.resolved_at IS NOT NULL AND st.resolved_at > st.sla_due_date THEN 1 
    END) as sla_breached_count,
    
    COUNT(CASE 
        WHEN st.status != 'resolved' AND NOW() > st.sla_due_date THEN 1 
    END) as sla_breaching_count,
    
    -- Compliance Rate
    ROUND(
        (COUNT(CASE WHEN st.resolved_at IS NOT NULL AND st.resolved_at <= st.sla_due_date THEN 1 END)::NUMERIC /
         NULLIF(COUNT(CASE WHEN st.resolved_at IS NOT NULL THEN 1 END), 0)) * 100,
        2
    ) as sla_compliance_pct,
    
    -- Timing
    ROUND(AVG(EXTRACT(EPOCH FROM (st.resolved_at - st.created_at)) / 3600), 2) as avg_resolution_hours,
    ROUND(AVG(EXTRACT(EPOCH FROM (st.first_response_at - st.created_at)) / 3600), 2) as avg_first_response_hours,
    
    -- At Risk
    COUNT(CASE 
        WHEN st.status != 'resolved' 
        AND st.sla_due_date BETWEEN NOW() AND NOW() + INTERVAL '4 hours' 
    THEN 1 END) as at_risk_4h_count

FROM service_tickets st
LEFT JOIN customers c ON st.customer_id = c.id
LEFT JOIN equipment e ON st.equipment_id = e.id
LEFT JOIN equipment_catalog ec ON e.catalog_id = ec.id
GROUP BY DATE_TRUNC('day', st.created_at), st.priority, st.severity, c.organization_type, ec.catalog_name;

COMMENT ON VIEW sla_compliance_dashboard IS 'Daily SLA compliance tracking with at-risk tickets';

-- =====================================================================
-- 7. COST ANALYTICS
-- =====================================================================

-- View: Comprehensive Cost Analysis
CREATE OR REPLACE VIEW ticket_cost_analytics AS
SELECT 
    st.id as ticket_id,
    st.ticket_number,
    c.name as customer_name,
    ec.catalog_name as equipment_type,
    
    -- Parts Cost
    COALESCE((SELECT SUM(actual_price * quantity_used) 
              FROM ticket_parts_required 
              WHERE ticket_id = st.id), 0) as parts_cost,
    
    -- AI Cost
    COALESCE((SELECT SUM(cost_usd)
              FROM ai_conversations
              WHERE ticket_id = st.id), 0) as ai_cost_usd,
    
    COALESCE((SELECT cost_usd
              FROM ai_diagnosis_results
              WHERE ticket_id = st.id), 0) as ai_diagnosis_cost,
    
    -- Engineer Cost (estimated)
    COALESCE((SELECT SUM(EXTRACT(EPOCH FROM (completed_at - assigned_at)) / 3600 * 50)
              FROM stage_assignments
              WHERE ticket_id = st.id 
              AND assignment_status = 'completed'), 0) as estimated_engineer_cost,
    
    -- Travel Cost (estimated)
    COALESCE((SELECT SUM(travel_distance_km * 0.5)
              FROM stage_assignments
              WHERE ticket_id = st.id 
              AND work_location = 'onsite'), 0) as estimated_travel_cost,
    
    -- Total Cost
    COALESCE((SELECT SUM(actual_price * quantity_used) FROM ticket_parts_required WHERE ticket_id = st.id), 0) +
    COALESCE((SELECT SUM(cost_usd) FROM ai_conversations WHERE ticket_id = st.id), 0) +
    COALESCE((SELECT SUM(EXTRACT(EPOCH FROM (completed_at - assigned_at)) / 3600 * 50) FROM stage_assignments WHERE ticket_id = st.id AND assignment_status = 'completed'), 0) +
    COALESCE((SELECT SUM(travel_distance_km * 0.5) FROM stage_assignments WHERE ticket_id = st.id AND work_location = 'onsite'), 0) as total_estimated_cost,
    
    -- Resolution Time
    EXTRACT(EPOCH FROM (st.resolved_at - st.created_at)) / 3600 as resolution_hours,
    
    -- Customer Satisfaction
    st.customer_satisfaction_rating

FROM service_tickets st
LEFT JOIN customers c ON st.customer_id = c.id
LEFT JOIN equipment e ON st.equipment_id = e.id
LEFT JOIN equipment_catalog ec ON e.catalog_id = ec.id
WHERE st.status = 'resolved';

COMMENT ON VIEW ticket_cost_analytics IS 'Complete cost breakdown per ticket including parts, AI, engineer time, and travel';

-- =====================================================================
-- 8. WORKFLOW EFFICIENCY METRICS
-- =====================================================================

-- View: Workflow Template Performance
CREATE OR REPLACE VIEW workflow_template_performance AS
SELECT 
    wt.id as template_id,
    wt.template_name,
    wt.template_type,
    ec.catalog_name as equipment_type,
    m.name as manufacturer_name,
    
    -- Usage
    COUNT(DISTINCT twi.id) as total_instances,
    COUNT(DISTINCT CASE WHEN twi.status = 'completed' THEN twi.id END) as completed_instances,
    COUNT(DISTINCT CASE WHEN twi.status = 'in_progress' THEN twi.id END) as in_progress_instances,
    
    -- Stages
    wt.total_stages,
    ROUND(AVG(twi.current_stage_number::NUMERIC / wt.total_stages * 100), 2) as avg_completion_pct,
    
    -- Timing
    ROUND(AVG(EXTRACT(EPOCH FROM (twi.completed_at - twi.started_at)) / 3600), 2) as avg_total_hours,
    ROUND(MIN(EXTRACT(EPOCH FROM (twi.completed_at - twi.started_at)) / 3600), 2) as best_time_hours,
    ROUND(MAX(EXTRACT(EPOCH FROM (twi.completed_at - twi.started_at)) / 3600), 2) as worst_time_hours,
    
    -- Quality
    ROUND(AVG(
        (SELECT AVG(customer_satisfaction_rating) 
         FROM stage_execution_data sed
         JOIN workflow_stage_transitions wst ON sed.workflow_stage_id = wst.id
         WHERE wst.workflow_instance_id = twi.id)
    ), 2) as avg_customer_satisfaction,
    
    -- Issues
    COUNT(CASE WHEN twi.has_issues = true THEN 1 END) as instances_with_issues,
    
    -- Active Count
    COUNT(CASE WHEN twi.status = 'in_progress' THEN 1 END) as currently_active

FROM workflow_templates wt
LEFT JOIN ticket_workflow_instances twi ON wt.id = twi.template_id
LEFT JOIN equipment_catalog ec ON wt.equipment_catalog_id = ec.id
LEFT JOIN manufacturers m ON wt.manufacturer_id = m.id
GROUP BY wt.id, wt.template_name, wt.template_type, ec.catalog_name, m.name, wt.total_stages;

COMMENT ON VIEW workflow_template_performance IS 'Performance metrics per workflow template with timing and quality';

-- =====================================================================
-- 9. CUSTOMER SATISFACTION TRACKING
-- =====================================================================

-- View: Customer Satisfaction Dashboard
CREATE OR REPLACE VIEW customer_satisfaction_dashboard AS
SELECT 
    c.id as customer_id,
    c.name as customer_name,
    c.organization_type,
    
    -- Ticket Volume
    COUNT(DISTINCT st.id) as total_tickets,
    COUNT(DISTINCT CASE WHEN st.status = 'resolved' THEN st.id END) as resolved_tickets,
    
    -- Satisfaction Ratings
    ROUND(AVG(st.customer_satisfaction_rating), 2) as avg_ticket_rating,
    ROUND(AVG(
        (SELECT AVG(customer_satisfaction_rating) 
         FROM stage_execution_data sed
         JOIN workflow_stage_transitions wst ON sed.workflow_stage_id = wst.id
         WHERE wst.ticket_id = st.id)
    ), 2) as avg_stage_rating,
    
    COUNT(CASE WHEN st.customer_satisfaction_rating >= 4 THEN 1 END) as positive_ratings_count,
    COUNT(CASE WHEN st.customer_satisfaction_rating <= 2 THEN 1 END) as negative_ratings_count,
    
    -- NPS Calculation
    ROUND(
        ((COUNT(CASE WHEN st.customer_satisfaction_rating >= 4 THEN 1 END)::NUMERIC - 
          COUNT(CASE WHEN st.customer_satisfaction_rating <= 2 THEN 1 END)) / 
         NULLIF(COUNT(st.customer_satisfaction_rating), 0)) * 100,
        2
    ) as nps_score,
    
    -- SLA Performance
    ROUND(
        (COUNT(CASE WHEN st.resolved_at IS NOT NULL AND st.resolved_at <= st.sla_due_date THEN 1 END)::NUMERIC /
         NULLIF(COUNT(CASE WHEN st.resolved_at IS NOT NULL THEN 1 END), 0)) * 100,
        2
    ) as sla_compliance_pct,
    
    -- Response Time
    ROUND(AVG(EXTRACT(EPOCH FROM (st.first_response_at - st.created_at)) / 3600), 2) as avg_first_response_hours,
    ROUND(AVG(EXTRACT(EPOCH FROM (st.resolved_at - st.created_at)) / 3600), 2) as avg_resolution_hours,
    
    -- Recent Activity
    MAX(st.created_at) as last_ticket_date,
    COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '30 days' THEN 1 END) as tickets_last_30_days

FROM customers c
LEFT JOIN service_tickets st ON c.id = st.customer_id
GROUP BY c.id, c.name, c.organization_type;

COMMENT ON VIEW customer_satisfaction_dashboard IS 'Customer satisfaction metrics with NPS scoring and SLA compliance';

-- =====================================================================
-- 10. REAL-TIME DASHBOARDS
-- =====================================================================

-- View: Real-Time Operations Dashboard
CREATE OR REPLACE VIEW realtime_operations_dashboard AS
SELECT 
    -- Current Active Tickets
    (SELECT COUNT(*) FROM service_tickets WHERE status IN ('open', 'assigned', 'in_progress')) as active_tickets_count,
    
    -- SLA At Risk
    (SELECT COUNT(*) FROM service_tickets 
     WHERE status != 'resolved' 
     AND sla_due_date BETWEEN NOW() AND NOW() + INTERVAL '4 hours') as sla_at_risk_4h,
    
    (SELECT COUNT(*) FROM service_tickets 
     WHERE status != 'resolved' 
     AND sla_due_date < NOW()) as sla_breaching_now,
    
    -- Engineers
    (SELECT COUNT(*) FROM engineers WHERE is_available = true) as available_engineers,
    (SELECT COUNT(*) FROM engineers WHERE status = 'active') as active_engineers,
    (SELECT COUNT(*) FROM stage_assignments 
     WHERE assignment_status IN ('assigned', 'in_progress', 'traveling')) as engineers_on_tickets,
    
    -- Today's Activity
    (SELECT COUNT(*) FROM service_tickets WHERE DATE(created_at) = CURRENT_DATE) as tickets_created_today,
    (SELECT COUNT(*) FROM service_tickets WHERE DATE(resolved_at) = CURRENT_DATE) as tickets_resolved_today,
    
    -- AI Activity (last hour)
    (SELECT COUNT(*) FROM ai_conversations WHERE created_at >= NOW() - INTERVAL '1 hour') as ai_requests_last_hour,
    (SELECT COUNT(*) FROM ai_diagnosis_results WHERE created_at >= NOW() - INTERVAL '1 hour') as ai_diagnoses_last_hour,
    
    -- Parts
    (SELECT COUNT(*) FROM ticket_parts_required 
     WHERE procurement_status IN ('requested', 'ordered')) as parts_in_procurement,
    
    -- Workflow Stages Stuck
    (SELECT COUNT(*) FROM workflow_stage_transitions 
     WHERE status = 'in_progress' 
     AND started_at < NOW() - INTERVAL '48 hours') as stages_stuck_over_48h,
    
    -- Average Response Time Today
    (SELECT ROUND(AVG(EXTRACT(EPOCH FROM (first_response_at - created_at)) / 3600), 2)
     FROM service_tickets 
     WHERE DATE(created_at) = CURRENT_DATE 
     AND first_response_at IS NOT NULL) as avg_response_hours_today,
    
    -- Critical Tickets
    (SELECT COUNT(*) FROM service_tickets 
     WHERE priority = 'critical' 
     AND status != 'resolved') as critical_tickets_open;

COMMENT ON VIEW realtime_operations_dashboard IS 'Real-time operational metrics for live dashboard';

-- View: Equipment Health Status
CREATE OR REPLACE VIEW equipment_health_status AS
SELECT 
    ec.catalog_name as equipment_type,
    m.name as manufacturer_name,
    e.equipment_name,
    e.serial_number,
    c.name as customer_name,
    
    -- Ticket History
    COUNT(st.id) as total_tickets,
    COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '30 days' THEN 1 END) as tickets_last_30_days,
    COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '7 days' THEN 1 END) as tickets_last_7_days,
    
    -- Issue Frequency
    CASE 
        WHEN COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '30 days' THEN 1 END) >= 5 THEN 'High Risk'
        WHEN COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '30 days' THEN 1 END) >= 3 THEN 'Medium Risk'
        WHEN COUNT(CASE WHEN st.created_at >= NOW() - INTERVAL '30 days' THEN 1 END) >= 1 THEN 'Low Risk'
        ELSE 'Healthy'
    END as health_status,
    
    -- Recent Issues
    (SELECT STRING_AGG(DISTINCT severity, ', ') 
     FROM service_tickets 
     WHERE equipment_id = e.id 
     AND created_at >= NOW() - INTERVAL '30 days') as recent_severities,
    
    -- Last Service
    MAX(st.resolved_at) as last_service_date,
    
    -- Parts Replaced (last 90 days)
    (SELECT COUNT(DISTINCT part_id) 
     FROM ticket_parts_required tpr
     JOIN service_tickets st2 ON tpr.ticket_id = st2.id
     WHERE st2.equipment_id = e.id
     AND tpr.procurement_status = 'used'
     AND st2.created_at >= NOW() - INTERVAL '90 days') as parts_replaced_90d,
    
    -- Cost (last 90 days)
    COALESCE((SELECT SUM(actual_price * quantity_used)
              FROM ticket_parts_required tpr
              JOIN service_tickets st2 ON tpr.ticket_id = st2.id
              WHERE st2.equipment_id = e.id
              AND st2.created_at >= NOW() - INTERVAL '90 days'), 0) as cost_last_90d

FROM equipment e
JOIN equipment_catalog ec ON e.catalog_id = ec.id
LEFT JOIN manufacturers m ON ec.manufacturer_id = m.id
LEFT JOIN customers c ON e.customer_id = c.id
LEFT JOIN service_tickets st ON e.id = st.equipment_id
GROUP BY e.id, ec.catalog_name, m.name, e.equipment_name, e.serial_number, c.name;

COMMENT ON VIEW equipment_health_status IS 'Equipment health scoring based on ticket frequency and severity';

-- =====================================================================
-- MATERIALIZED VIEWS FOR PERFORMANCE
-- =====================================================================

-- Materialized View: Daily Performance Metrics (refreshed hourly)
CREATE MATERIALIZED VIEW mv_daily_performance_metrics AS
SELECT 
    DATE_TRUNC('day', created_at) as date,
    COUNT(*) as tickets_created,
    COUNT(CASE WHEN status = 'resolved' THEN 1 END) as tickets_resolved,
    ROUND(AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)) / 3600), 2) as avg_resolution_hours,
    COUNT(CASE WHEN resolved_at <= sla_due_date THEN 1 END) as sla_met_count,
    ROUND(AVG(customer_satisfaction_rating), 2) as avg_satisfaction
FROM service_tickets
GROUP BY DATE_TRUNC('day', created_at);

CREATE UNIQUE INDEX idx_mv_daily_perf_date ON mv_daily_performance_metrics(date);

COMMENT ON MATERIALIZED VIEW mv_daily_performance_metrics IS 'Daily aggregated metrics - refresh hourly';

-- =====================================================================
-- REFRESH FUNCTION
-- =====================================================================

CREATE OR REPLACE FUNCTION refresh_analytics_materialized_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY mv_daily_performance_metrics;
    RAISE NOTICE 'Analytics materialized views refreshed at %', NOW();
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_analytics_materialized_views IS 'Refresh all analytics materialized views - run hourly via cron';

-- =====================================================================
-- MIGRATION COMPLETE
-- =====================================================================

DO $$
BEGIN
    RAISE NOTICE 'Migration 021 complete!';
    RAISE NOTICE 'Created 11 analytics views:';
    RAISE NOTICE '  - ticket_performance_overview';
    RAISE NOTICE '  - stage_performance_metrics';
    RAISE NOTICE '  - engineer_performance_dashboard';
    RAISE NOTICE '  - parts_usage_analytics';
    RAISE NOTICE '  - ai_diagnosis_performance_by_equipment';
    RAISE NOTICE '  - sla_compliance_dashboard';
    RAISE NOTICE '  - ticket_cost_analytics';
    RAISE NOTICE '  - workflow_template_performance';
    RAISE NOTICE '  - customer_satisfaction_dashboard';
    RAISE NOTICE '  - realtime_operations_dashboard';
    RAISE NOTICE '  - equipment_health_status';
    RAISE NOTICE 'Created 1 materialized view for performance';
    RAISE NOTICE 'Analytics and monitoring ready!';
END $$;
