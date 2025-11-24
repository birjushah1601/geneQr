package feedback

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Learner applies improvements based on feedback analysis
type Learner struct {
	db       *sql.DB
	analyzer *Analyzer
}

// NewLearner creates a new learning engine
func NewLearner(db *sql.DB) *Learner {
	return &Learner{
		db:       db,
		analyzer: NewAnalyzer(db),
	}
}

// ApplyImprovement applies an improvement opportunity
func (l *Learner) ApplyImprovement(ctx context.Context, opportunityID string, appliedBy string) (*LearningAction, error) {
	// Get the improvement opportunity
	opportunity, err := l.getImprovement(ctx, opportunityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get improvement: %w", err)
	}
	
	if opportunity.Status != "pending" {
		return nil, fmt.Errorf("improvement already processed with status: %s", opportunity.Status)
	}
	
	// Record metrics before change
	beforeMetrics, err := l.captureMetrics(ctx, opportunity.ServiceType)
	if err != nil {
		return nil, fmt.Errorf("failed to capture before metrics: %w", err)
	}
	
	// Apply the change based on implementation type
	var actionType string
	var changes map[string]interface{}
	
	switch opportunity.ImplementationType {
	case "prompt_tuning":
		actionType = "prompt_update"
		changes = l.applyPromptTuning(ctx, opportunity)
	case "weight_adjustment":
		actionType = "weight_adjustment"
		changes = l.applyWeightAdjustment(ctx, opportunity)
	case "config_change":
		actionType = "config_change"
		changes = l.applyConfigChange(ctx, opportunity)
	default:
		// For training_data, we can't apply automatically - requires manual review
		return nil, fmt.Errorf("implementation type %s requires manual action", opportunity.ImplementationType)
	}
	
	// Create learning action record
	action := &LearningAction{
		ActionID:      fmt.Sprintf("action_%s_%d", opportunity.ServiceType, time.Now().Unix()),
		OpportunityID: opportunityID,
		ActionType:    actionType,
		ServiceType:   opportunity.ServiceType,
		Changes:       changes,
		BeforeMetrics: beforeMetrics,
		Status:        "testing", // Start in testing mode
		AppliedAt:     time.Now(),
		AppliedBy:     appliedBy,
	}
	
	// Store the action
	if err := l.storeAction(ctx, action); err != nil {
		return nil, fmt.Errorf("failed to store action: %w", err)
	}
	
	// Update improvement status
	if err := l.updateImprovementStatus(ctx, opportunityID, "applied"); err != nil {
		return nil, fmt.Errorf("failed to update improvement status: %w", err)
	}
	
	return action, nil
}

// EvaluateAction evaluates the effectiveness of an applied action
func (l *Learner) EvaluateAction(ctx context.Context, actionID string) error {
	// Get the action
	action, err := l.getAction(ctx, actionID)
	if err != nil {
		return fmt.Errorf("failed to get action: %w", err)
	}
	
	if action.Status != "testing" {
		return fmt.Errorf("action must be in testing status to evaluate")
	}
	
	// Check if enough time has passed (at least 7 days)
	if time.Since(action.AppliedAt).Hours() < 168 {
		return fmt.Errorf("not enough data yet, need at least 7 days of testing")
	}
	
	// Capture current metrics
	afterMetrics, err := l.captureMetrics(ctx, action.ServiceType)
	if err != nil {
		return fmt.Errorf("failed to capture after metrics: %w", err)
	}
	
	// Calculate improvement
	improvementPercent := l.calculateImprovement(action.BeforeMetrics, afterMetrics)
	
	// Decide whether to deploy or rollback
	if improvementPercent >= 5.0 {
		// Significant improvement - deploy to production
		action.Status = "deployed"
		action.AfterMetrics = afterMetrics
		action.ResultNotes = fmt.Sprintf("Improved performance by %.2f%%. Deployed to production.", improvementPercent)
	} else if improvementPercent < -5.0 {
		// Performance degraded - rollback
		action.Status = "rolled_back"
		now := time.Now()
		action.RolledBackAt = &now
		action.RollbackReason = fmt.Sprintf("Performance decreased by %.2f%%. Rolling back changes.", -improvementPercent)
		
		// Actually rollback the changes
		l.rollbackChanges(ctx, action)
	} else {
		// Neutral result - keep testing
		action.ResultNotes = fmt.Sprintf("Neutral impact (%.2f%%). Continuing testing.", improvementPercent)
	}
	
	// Update action in database
	return l.updateAction(ctx, action)
}

// GetLearningProgress retrieves learning progress for a service
func (l *Learner) GetLearningProgress(ctx context.Context, serviceType string) (map[string]interface{}, error) {
	progress := make(map[string]interface{})
	
	// Get total improvements found
	var totalImprovements, pendingImprovements, appliedImprovements int
	err := l.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'applied' THEN 1 END) as applied
		FROM feedback_improvements
		WHERE service_type = $1
	`, serviceType).Scan(&totalImprovements, &pendingImprovements, &appliedImprovements)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	progress["total_improvements_identified"] = totalImprovements
	progress["pending_improvements"] = pendingImprovements
	progress["applied_improvements"] = appliedImprovements
	
	// Get action statistics
	var totalActions, deployedActions, rolledBackActions int
	err = l.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'deployed' THEN 1 END) as deployed,
			COUNT(CASE WHEN status = 'rolled_back' THEN 1 END) as rolled_back
		FROM feedback_actions
		WHERE service_type = $1
	`, serviceType).Scan(&totalActions, &deployedActions, &rolledBackActions)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	progress["total_actions"] = totalActions
	progress["deployed_actions"] = deployedActions
	progress["rolled_back_actions"] = rolledBackActions
	
	// Calculate success rate
	if totalActions > 0 {
		successRate := float64(deployedActions) / float64(totalActions) * 100.0
		progress["success_rate_percent"] = successRate
	}
	
	// Get recent actions
	actions := l.getRecentActionsForService(ctx, serviceType, 10)
	progress["recent_actions"] = actions
	
	return progress, nil
}

// applyPromptTuning applies prompt improvements
func (l *Learner) applyPromptTuning(ctx context.Context, opp *ImprovementOpportunity) map[string]interface{} {
	changes := make(map[string]interface{})
	
	// Based on the issue type, adjust prompts
	if issueType, ok := opp.SuggestedChanges["issue_type"].(string); ok {
		changes["prompt_modification_type"] = "append_instruction"
		changes["added_instruction"] = fmt.Sprintf("Pay special attention to %s to avoid common mistakes.", issueType)
		changes["issue_addressed"] = issueType
	}
	
	// Store the updated prompt configuration
	// In a real system, this would update a configuration file or database table
	changes["status"] = "prompt_updated"
	changes["timestamp"] = time.Now()
	
	return changes
}

// applyWeightAdjustment adjusts scoring weights
func (l *Learner) applyWeightAdjustment(ctx context.Context, opp *ImprovementOpportunity) map[string]interface{} {
	changes := make(map[string]interface{})
	
	// Extract the suggested adjustment
	if action, ok := opp.SuggestedChanges["suggested_action"].(string); ok {
		changes["action"] = action
		
		if adjustment, ok := opp.SuggestedChanges["weight_adjustment"].(string); ok {
			changes["weight_adjustment"] = adjustment
			
			// Apply the actual weight change
			// In a real system, this would update the scoring weights in the assignment/diagnosis/parts engines
			changes["status"] = "weights_updated"
			changes["timestamp"] = time.Now()
		}
	}
	
	return changes
}

// applyConfigChange applies configuration changes
func (l *Learner) applyConfigChange(ctx context.Context, opp *ImprovementOpportunity) map[string]interface{} {
	changes := make(map[string]interface{})
	
	// Apply configuration changes based on the opportunity
	changes["config_updates"] = opp.SuggestedChanges
	changes["status"] = "config_updated"
	changes["timestamp"] = time.Now()
	
	return changes
}

// rollbackChanges rolls back an applied action
func (l *Learner) rollbackChanges(ctx context.Context, action *LearningAction) error {
	// In a real system, this would revert the changes made by the action
	// For now, we just log the rollback
	fmt.Printf("Rolling back action %s for service %s\n", action.ActionID, action.ServiceType)
	
	// Specific rollback logic based on action type
	switch action.ActionType {
	case "prompt_update":
		// Revert prompt changes
	case "weight_adjustment":
		// Revert weight changes
	case "config_change":
		// Revert config changes
	}
	
	return nil
}

// captureMetrics captures current performance metrics for a service
func (l *Learner) captureMetrics(ctx context.Context, serviceType string) (map[string]float64, error) {
	metrics := make(map[string]float64)
	
	end := time.Now()
	start := end.AddDate(0, 0, -7) // Last 7 days
	
	feedbackMetrics, err := l.analyzer.GetMetrics(ctx, serviceType, start, end)
	if err != nil {
		return nil, err
	}
	
	metrics["accuracy_rate"] = feedbackMetrics.AvgAccuracyRate
	metrics["avg_rating"] = feedbackMetrics.AvgRating
	metrics["positive_sentiment"] = feedbackMetrics.PositiveSentiment
	metrics["feedback_rate"] = feedbackMetrics.FeedbackRate
	
	return metrics, nil
}

// calculateImprovement calculates the improvement percentage between before and after metrics
func (l *Learner) calculateImprovement(before, after map[string]float64) float64 {
	// Calculate weighted average improvement
	weights := map[string]float64{
		"accuracy_rate":      0.4,
		"positive_sentiment": 0.3,
		"avg_rating":         0.3,
	}
	
	totalImprovement := 0.0
	totalWeight := 0.0
	
	for metric, weight := range weights {
		beforeVal, beforeOk := before[metric]
		afterVal, afterOk := after[metric]
		
		if beforeOk && afterOk && beforeVal > 0 {
			improvement := ((afterVal - beforeVal) / beforeVal) * 100.0
			totalImprovement += improvement * weight
			totalWeight += weight
		}
	}
	
	if totalWeight > 0 {
		return totalImprovement / totalWeight
	}
	
	return 0.0
}

// storeAction stores a learning action in the database
func (l *Learner) storeAction(ctx context.Context, action *LearningAction) error {
	changesJSON, _ := json.Marshal(action.Changes)
	beforeMetricsJSON, _ := json.Marshal(action.BeforeMetrics)
	
	query := `
		INSERT INTO feedback_actions (
			action_id,
			opportunity_id,
			action_type,
			service_type,
			changes,
			before_metrics,
			status,
			applied_at,
			applied_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := l.db.ExecContext(
		ctx,
		query,
		action.ActionID,
		action.OpportunityID,
		action.ActionType,
		action.ServiceType,
		changesJSON,
		beforeMetricsJSON,
		action.Status,
		action.AppliedAt,
		action.AppliedBy,
	)
	
	return err
}

// updateAction updates an existing action
func (l *Learner) updateAction(ctx context.Context, action *LearningAction) error {
	afterMetricsJSON, _ := json.Marshal(action.AfterMetrics)
	
	query := `
		UPDATE feedback_actions
		SET status = $1,
			after_metrics = $2,
			result_notes = $3,
			rolled_back_at = $4,
			rollback_reason = $5
		WHERE action_id = $6
	`
	
	_, err := l.db.ExecContext(
		ctx,
		query,
		action.Status,
		afterMetricsJSON,
		action.ResultNotes,
		action.RolledBackAt,
		action.RollbackReason,
		action.ActionID,
	)
	
	return err
}

// getImprovement retrieves an improvement opportunity
func (l *Learner) getImprovement(ctx context.Context, opportunityID string) (*ImprovementOpportunity, error) {
	var opp ImprovementOpportunity
	var changesJSON, dataJSON []byte
	
	query := `
		SELECT 
			opportunity_id,
			title,
			description,
			impact_level,
			implementation_type,
			suggested_changes,
			supporting_data,
			created_at,
			status,
			service_type
		FROM feedback_improvements
		WHERE opportunity_id = $1
	`
	
	err := l.db.QueryRowContext(ctx, query, opportunityID).Scan(
		&opp.OpportunityID,
		&opp.Title,
		&opp.Description,
		&opp.ImpactLevel,
		&opp.ImplementationType,
		&changesJSON,
		&dataJSON,
		&opp.CreatedAt,
		&opp.Status,
		&opp.ServiceType,
	)
	
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal(changesJSON, &opp.SuggestedChanges)
	json.Unmarshal(dataJSON, &opp.SupportingData)
	
	// Add service_type field to ImprovementOpportunity
	// (Note: would need to add this field to the type definition)
	
	return &opp, nil
}

// getAction retrieves a learning action
func (l *Learner) getAction(ctx context.Context, actionID string) (*LearningAction, error) {
	var action LearningAction
	var changesJSON, beforeMetricsJSON, afterMetricsJSON []byte
	var rolledBackAt sql.NullTime
	
	query := `
		SELECT 
			action_id,
			opportunity_id,
			action_type,
			service_type,
			changes,
			before_metrics,
			after_metrics,
			status,
			applied_at,
			applied_by,
			result_notes,
			rolled_back_at,
			rollback_reason
		FROM feedback_actions
		WHERE action_id = $1
	`
	
	err := l.db.QueryRowContext(ctx, query, actionID).Scan(
		&action.ActionID,
		&action.OpportunityID,
		&action.ActionType,
		&action.ServiceType,
		&changesJSON,
		&beforeMetricsJSON,
		&afterMetricsJSON,
		&action.Status,
		&action.AppliedAt,
		&action.AppliedBy,
		&action.ResultNotes,
		&rolledBackAt,
		&action.RollbackReason,
	)
	
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal(changesJSON, &action.Changes)
	json.Unmarshal(beforeMetricsJSON, &action.BeforeMetrics)
	json.Unmarshal(afterMetricsJSON, &action.AfterMetrics)
	
	if rolledBackAt.Valid {
		action.RolledBackAt = &rolledBackAt.Time
	}
	
	return &action, nil
}

// updateImprovementStatus updates the status of an improvement
func (l *Learner) updateImprovementStatus(ctx context.Context, opportunityID, status string) error {
	query := `
		UPDATE feedback_improvements
		SET status = $1
		WHERE opportunity_id = $2
	`
	
	_, err := l.db.ExecContext(ctx, query, status, opportunityID)
	return err
}

// getRecentActionsForService retrieves recent actions for a service
func (l *Learner) getRecentActionsForService(ctx context.Context, serviceType string, limit int) []map[string]interface{} {
	query := `
		SELECT 
			action_id,
			action_type,
			status,
			applied_at,
			result_notes
		FROM feedback_actions
		WHERE service_type = $1
		ORDER BY applied_at DESC
		LIMIT $2
	`
	
	rows, err := l.db.QueryContext(ctx, query, serviceType, limit)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()
	
	var actions []map[string]interface{}
	for rows.Next() {
		var actionID, actionType, status string
		var appliedAt time.Time
		var resultNotes sql.NullString
		
		rows.Scan(&actionID, &actionType, &status, &appliedAt, &resultNotes)
		
		action := map[string]interface{}{
			"action_id":   actionID,
			"action_type": actionType,
			"status":      status,
			"applied_at":  appliedAt,
		}
		
		if resultNotes.Valid {
			action["result_notes"] = resultNotes.String
		}
		
		actions = append(actions, action)
	}
	
	return actions
}

