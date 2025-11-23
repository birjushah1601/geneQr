package feedback

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Analyzer analyzes feedback to identify patterns and improvement opportunities
type Analyzer struct {
	db        *sql.DB
	collector *Collector
}

// NewAnalyzer creates a new feedback analyzer
func NewAnalyzer(db *sql.DB) *Analyzer {
	return &Analyzer{
		db:        db,
		collector: NewCollector(db),
	}
}

// AnalyzeFeedback analyzes recent feedback for a service type
func (a *Analyzer) AnalyzeFeedback(ctx context.Context, serviceType string, days int) (*FeedbackAnalysis, error) {
	feedback, err := a.collector.GetRecentFeedback(ctx, serviceType, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}
	
	analysis := &FeedbackAnalysis{
		ServiceType:          serviceType,
		Period:              fmt.Sprintf("%d_days", days),
		TotalFeedback:       len(feedback),
		GeneratedAt:         time.Now(),
	}
	
	// Count by source
	for _, f := range feedback {
		if f.Source == SourceHuman {
			analysis.HumanFeedbackCount++
		} else {
			analysis.MachineFeedbackCount++
		}
	}
	
	// Count by sentiment
	totalRating := 0.0
	ratingCount := 0
	accurateCount := 0
	totalCount := 0
	
	for _, f := range feedback {
		switch f.Sentiment {
		case SentimentPositive:
			analysis.PositiveFeedback++
			accurateCount++
		case SentimentNeutral:
			analysis.NeutralFeedback++
		case SentimentNegative:
			analysis.NegativeFeedback++
		}
		
		if f.Rating != nil {
			totalRating += float64(*f.Rating)
			ratingCount++
		}
		
		totalCount++
	}
	
	if ratingCount > 0 {
		analysis.AvgRating = totalRating / float64(ratingCount)
	}
	
	if totalCount > 0 {
		analysis.AccuracyRate = float64(accurateCount) / float64(totalCount) * 100.0
	}
	
	// Identify common issues
	analysis.CommonIssues = a.identifyCommonIssues(feedback)
	
	// Generate improvement opportunities
	analysis.Improvements = a.generateImprovementOpportunities(ctx, serviceType, feedback, analysis.CommonIssues)
	
	return analysis, nil
}

// GetMetrics calculates feedback metrics for a service type
func (a *Analyzer) GetMetrics(ctx context.Context, serviceType string, start, end time.Time) (*FeedbackMetrics, error) {
	metrics := &FeedbackMetrics{
		ServiceType: serviceType,
		DateRange: DateRange{
			Start: start,
			End:   end,
		},
	}
	
	// Get total AI requests in period
	err := a.db.QueryRowContext(ctx, a.getTotalRequestsQuery(serviceType), start, end).Scan(&metrics.TotalRequests)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	// Get feedback counts and metrics
	query := `
		SELECT 
			COUNT(*) as feedback_count,
			AVG(CASE WHEN sentiment = 'positive' THEN 100.0 ELSE 0.0 END) as accuracy_rate,
			AVG(rating) as avg_rating,
			COUNT(CASE WHEN sentiment = 'positive' THEN 1 END) * 100.0 / COUNT(*) as positive_sentiment
		FROM ai_feedback
		WHERE service_type = $1
			AND created_at BETWEEN $2 AND $3
	`
	
	err = a.db.QueryRowContext(ctx, query, serviceType, start, end).Scan(
		&metrics.FeedbackReceived,
		&metrics.AvgAccuracyRate,
		&metrics.AvgRating,
		&metrics.PositiveSentiment,
	)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	if metrics.TotalRequests > 0 {
		metrics.FeedbackRate = float64(metrics.FeedbackReceived) / float64(metrics.TotalRequests) * 100.0
	}
	
	// Get improvement counts
	err = a.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as improvements_found,
			COUNT(CASE WHEN status = 'applied' THEN 1 END) as actions_applied
		FROM feedback_improvements
		WHERE service_type = $1
			AND created_at BETWEEN $2 AND $3
	`, serviceType, start, end).Scan(&metrics.ImprovementsFound, &metrics.ActionsApplied)
	
	if err != nil && err != sql.ErrNoRows {
		// Table might not exist yet, ignore
	}
	
	// Calculate trends
	metrics.AccuracyTrend = a.calculateTrend(ctx, serviceType, "accuracy", start, end)
	metrics.SentimentTrend = a.calculateTrend(ctx, serviceType, "sentiment", start, end)
	
	return metrics, nil
}

// GetSummary provides a comprehensive feedback summary
func (a *Analyzer) GetSummary(ctx context.Context, days int) (*FeedbackSummary, error) {
	summary := &FeedbackSummary{
		ByServiceType:  make(map[string]FeedbackMetrics),
		GeneratedAt:    time.Now(),
	}
	
	end := time.Now()
	start := end.AddDate(0, 0, -days)
	
	// Get metrics for each service type
	serviceTypes := []string{"diagnosis", "assignment", "parts"}
	
	var totalRequests, totalFeedback int
	var sumAccuracy, sumRating, sumPositiveSentiment float64
	var sumImprovements, sumActions int
	
	for _, serviceType := range serviceTypes {
		metrics, err := a.GetMetrics(ctx, serviceType, start, end)
		if err != nil {
			continue
		}
		
		summary.ByServiceType[serviceType] = *metrics
		
		// Aggregate for overall
		totalRequests += metrics.TotalRequests
		totalFeedback += metrics.FeedbackReceived
		sumAccuracy += metrics.AvgAccuracyRate
		sumRating += metrics.AvgRating
		sumPositiveSentiment += metrics.PositiveSentiment
		sumImprovements += metrics.ImprovementsFound
		sumActions += metrics.ActionsApplied
	}
	
	// Calculate overall metrics
	count := len(serviceTypes)
	summary.OverallMetrics = FeedbackMetrics{
		ServiceType:       "all",
		DateRange:         DateRange{Start: start, End: end},
		TotalRequests:     totalRequests,
		FeedbackReceived:  totalFeedback,
		AvgAccuracyRate:   sumAccuracy / float64(count),
		AvgRating:         sumRating / float64(count),
		PositiveSentiment: sumPositiveSentiment / float64(count),
		ImprovementsFound: sumImprovements,
		ActionsApplied:    sumActions,
	}
	
	if totalRequests > 0 {
		summary.OverallMetrics.FeedbackRate = float64(totalFeedback) / float64(totalRequests) * 100.0
	}
	
	// Get recent improvements
	summary.RecentImprovements = a.getRecentImprovements(ctx, 10)
	
	// Get active actions
	summary.ActiveActions = a.getActiveActions(ctx)
	
	// Get top issues
	summary.TopIssues = a.getTopIssues(ctx, days)
	
	return summary, nil
}

// identifyCommonIssues identifies recurring problems from feedback
func (a *Analyzer) identifyCommonIssues(feedback []*FeedbackEntry) []FeedbackIssue {
	// Group negative feedback by keywords in comments
	issueKeywords := make(map[string][]int64)
	
	for _, f := range feedback {
		if f.Sentiment != SentimentNegative {
			continue
		}
		
		// Extract keywords from comments
		keywords := extractKeywords(f.Comments)
		for _, keyword := range keywords {
			issueKeywords[keyword] = append(issueKeywords[keyword], f.FeedbackID)
		}
		
		// Check corrections for patterns
		if len(f.Corrections) > 0 {
			for key := range f.Corrections {
				issueKeywords["incorrect_"+key] = append(issueKeywords["incorrect_"+key], f.FeedbackID)
			}
		}
	}
	
	// Convert to FeedbackIssue list
	var issues []FeedbackIssue
	for keyword, feedbackIDs := range issueKeywords {
		if len(feedbackIDs) < 2 {
			// Skip issues that only occurred once
			continue
		}
		
		severity := "low"
		if len(feedbackIDs) >= 5 {
			severity = "high"
		} else if len(feedbackIDs) >= 3 {
			severity = "medium"
		}
		
		issues = append(issues, FeedbackIssue{
			IssueType:   keyword,
			Description: fmt.Sprintf("Multiple occurrences of: %s", keyword),
			Frequency:   len(feedbackIDs),
			Severity:    severity,
			Examples:    feedbackIDs[:min(3, len(feedbackIDs))], // Include up to 3 examples
		})
	}
	
	// Sort by frequency
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Frequency > issues[j].Frequency
	})
	
	// Return top 10 issues
	if len(issues) > 10 {
		return issues[:10]
	}
	
	return issues
}

// generateImprovementOpportunities generates actionable improvements from feedback analysis
func (a *Analyzer) generateImprovementOpportunities(ctx context.Context, serviceType string, feedback []*FeedbackEntry, issues []FeedbackIssue) []ImprovementOpportunity {
	var opportunities []ImprovementOpportunity
	
	// Analyze human corrections for patterns
	corrections := a.analyzeCorrections(feedback)
	
	for correctionType, examples := range corrections {
		if len(examples) < 2 {
			continue
		}
		
		impactLevel := "low"
		if len(examples) >= 5 {
			impactLevel = "high"
		} else if len(examples) >= 3 {
			impactLevel = "medium"
		}
		
		opportunity := ImprovementOpportunity{
			OpportunityID:   fmt.Sprintf("%s_%s_%d", serviceType, correctionType, time.Now().Unix()),
			Title:           fmt.Sprintf("Improve %s recommendations", correctionType),
			Description:     fmt.Sprintf("Users frequently correct %s, suggesting the AI model needs adjustment", correctionType),
			ImpactLevel:     impactLevel,
			ImplementationType: a.determineImplementationType(correctionType),
			SuggestedChanges: a.generateSuggestedChanges(serviceType, correctionType, examples),
			SupportingData:  examples,
			CreatedAt:       time.Now(),
			Status:          "pending",
		}
		
		opportunities = append(opportunities, opportunity)
	}
	
	// Generate opportunities from common issues
	for _, issue := range issues {
		if issue.Severity == "low" {
			continue
		}
		
		opportunity := ImprovementOpportunity{
			OpportunityID:   fmt.Sprintf("%s_issue_%s_%d", serviceType, issue.IssueType, time.Now().Unix()),
			Title:           fmt.Sprintf("Address: %s", issue.Description),
			Description:     fmt.Sprintf("Common issue identified in %d cases", issue.Frequency),
			ImpactLevel:     issue.Severity,
			ImplementationType: "prompt_tuning",
			SuggestedChanges: map[string]interface{}{
				"issue_type": issue.IssueType,
				"frequency":  issue.Frequency,
			},
			SupportingData: issue.Examples,
			CreatedAt:      time.Now(),
			Status:         "pending",
		}
		
		opportunities = append(opportunities, opportunity)
	}
	
	return opportunities
}

// analyzeCorrections analyzes human corrections to find patterns
func (a *Analyzer) analyzeCorrections(feedback []*FeedbackEntry) map[string][]int64 {
	corrections := make(map[string][]int64)
	
	for _, f := range feedback {
		if f.Source != SourceHuman || len(f.Corrections) == 0 {
			continue
		}
		
		for key := range f.Corrections {
			corrections[key] = append(corrections[key], f.FeedbackID)
		}
	}
	
	return corrections
}

// determineImplementationType determines how an improvement should be implemented
func (a *Analyzer) determineImplementationType(correctionType string) string {
	switch {
	case strings.Contains(correctionType, "priority") || strings.Contains(correctionType, "score"):
		return "weight_adjustment"
	case strings.Contains(correctionType, "confidence"):
		return "weight_adjustment"
	case strings.Contains(correctionType, "diagnosis") || strings.Contains(correctionType, "recommendation"):
		return "prompt_tuning"
	default:
		return "training_data"
	}
}

// generateSuggestedChanges generates specific changes based on feedback patterns
func (a *Analyzer) generateSuggestedChanges(serviceType, correctionType string, examples []int64) map[string]interface{} {
	changes := map[string]interface{}{
		"correction_type": correctionType,
		"example_count":   len(examples),
	}
	
	// Add service-specific suggested changes
	switch serviceType {
	case "assignment":
		if strings.Contains(correctionType, "expertise") {
			changes["suggested_action"] = "increase_expertise_weight"
			changes["weight_adjustment"] = "+5%"
		} else if strings.Contains(correctionType, "location") {
			changes["suggested_action"] = "increase_location_weight"
			changes["weight_adjustment"] = "+5%"
		}
	case "diagnosis":
		if strings.Contains(correctionType, "confidence") {
			changes["suggested_action"] = "adjust_confidence_thresholds"
			changes["threshold_adjustment"] = "-10%"
		}
	case "parts":
		if strings.Contains(correctionType, "priority") {
			changes["suggested_action"] = "adjust_recommendation_logic"
		}
	}
	
	return changes
}

// getTotalRequestsQuery returns the query to count total AI requests
func (a *Analyzer) getTotalRequestsQuery(serviceType string) string {
	switch serviceType {
	case "diagnosis":
		return "SELECT COUNT(*) FROM ai_diagnoses WHERE created_at BETWEEN $1 AND $2"
	case "assignment":
		return "SELECT COUNT(*) FROM assignment_history WHERE created_at BETWEEN $1 AND $2"
	case "parts":
		return "SELECT COUNT(*) FROM parts_recommendations WHERE created_at BETWEEN $1 AND $2"
	default:
		return "SELECT 0"
	}
}

// calculateTrend calculates trend direction for a metric
func (a *Analyzer) calculateTrend(ctx context.Context, serviceType, metricType string, start, end time.Time) string {
	// Split period in half
	mid := start.Add(end.Sub(start) / 2)
	
	var firstHalf, secondHalf float64
	
	query := a.getTrendQuery(metricType)
	
	a.db.QueryRowContext(ctx, query, serviceType, start, mid).Scan(&firstHalf)
	a.db.QueryRowContext(ctx, query, serviceType, mid, end).Scan(&secondHalf)
	
	if secondHalf > firstHalf*1.05 {
		return "improving"
	} else if secondHalf < firstHalf*0.95 {
		return "declining"
	}
	
	return "stable"
}

// getTrendQuery returns the query for trend calculation
func (a *Analyzer) getTrendQuery(metricType string) string {
	switch metricType {
	case "accuracy":
		return `
			SELECT AVG(CASE WHEN sentiment = 'positive' THEN 100.0 ELSE 0.0 END)
			FROM ai_feedback
			WHERE service_type = $1 AND created_at BETWEEN $2 AND $3
		`
	case "sentiment":
		return `
			SELECT COUNT(CASE WHEN sentiment = 'positive' THEN 1 END) * 100.0 / COUNT(*)
			FROM ai_feedback
			WHERE service_type = $1 AND created_at BETWEEN $2 AND $3
		`
	default:
		return "SELECT 0"
	}
}

// getRecentImprovements retrieves recent improvement opportunities
func (a *Analyzer) getRecentImprovements(ctx context.Context, limit int) []ImprovementOpportunity {
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
			status
		FROM feedback_improvements
		ORDER BY created_at DESC
		LIMIT $1
	`
	
	rows, err := a.db.QueryContext(ctx, query, limit)
	if err != nil {
		return []ImprovementOpportunity{}
	}
	defer rows.Close()
	
	var improvements []ImprovementOpportunity
	for rows.Next() {
		var imp ImprovementOpportunity
		var changesJSON, dataJSON []byte
		
		rows.Scan(
			&imp.OpportunityID,
			&imp.Title,
			&imp.Description,
			&imp.ImpactLevel,
			&imp.ImplementationType,
			&changesJSON,
			&dataJSON,
			&imp.CreatedAt,
			&imp.Status,
		)
		
		// Parse JSON fields (simplified for now)
		improvements = append(improvements, imp)
	}
	
	return improvements
}

// getActiveActions retrieves active learning actions
func (a *Analyzer) getActiveActions(ctx context.Context) []LearningAction {
	query := `
		SELECT 
			action_id,
			opportunity_id,
			action_type,
			service_type,
			changes,
			status,
			applied_at,
			applied_by
		FROM feedback_actions
		WHERE status IN ('testing', 'deployed')
		ORDER BY applied_at DESC
	`
	
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return []LearningAction{}
	}
	defer rows.Close()
	
	var actions []LearningAction
	for rows.Next() {
		var action LearningAction
		var changesJSON []byte
		
		rows.Scan(
			&action.ActionID,
			&action.OpportunityID,
			&action.ActionType,
			&action.ServiceType,
			&changesJSON,
			&action.Status,
			&action.AppliedAt,
			&action.AppliedBy,
		)
		
		actions = append(actions, action)
	}
	
	return actions
}

// getTopIssues retrieves top issues from recent feedback
func (a *Analyzer) getTopIssues(ctx context.Context, days int) []FeedbackIssue {
	// Get all recent negative feedback
	feedback, err := a.collector.GetRecentFeedback(ctx, "", days)
	if err != nil {
		return []FeedbackIssue{}
	}
	
	// Filter to negative only
	var negativeFeedback []*FeedbackEntry
	for _, f := range feedback {
		if f.Sentiment == SentimentNegative {
			negativeFeedback = append(negativeFeedback, f)
		}
	}
	
	// Identify common issues
	issues := a.identifyCommonIssues(negativeFeedback)
	
	// Return top 5
	if len(issues) > 5 {
		return issues[:5]
	}
	
	return issues
}

// extractKeywords extracts meaningful keywords from text
func extractKeywords(text string) []string {
	// Simple keyword extraction (could be enhanced with NLP)
	words := strings.Fields(strings.ToLower(text))
	
	// Filter out common words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true,
		"but": true, "in": true, "on": true, "at": true, "to": true,
		"for": true, "of": true, "is": true, "was": true, "it": true,
	}
	
	var keywords []string
	for _, word := range words {
		if len(word) > 3 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

