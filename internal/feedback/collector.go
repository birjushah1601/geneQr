package feedback

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Collector collects feedback from human and machine sources
type Collector struct {
	db *sql.DB
}

// NewCollector creates a new feedback collector
func NewCollector(db *sql.DB) *Collector {
	return &Collector{db: db}
}

// CollectHumanFeedback processes explicit human feedback
func (c *Collector) CollectHumanFeedback(ctx context.Context, req *HumanFeedbackRequest) (*FeedbackEntry, error) {
	// Determine feedback type based on service
	feedbackType := c.determineFeedbackType(req.ServiceType, true)
	
	// Determine sentiment based on rating and accuracy
	sentiment := c.determineSentiment(req.Rating, req.WasAccurate)
	
	// Create feedback entry
	entry := &FeedbackEntry{
		Source:      SourceHuman,
		Type:        feedbackType,
		RequestID:   &req.RequestID,
		TicketID:    req.TicketID,
		ServiceType: req.ServiceType,
		UserID:      &req.UserID,
		Rating:      req.Rating,
		Sentiment:   sentiment,
		Comments:    req.Comments,
		Corrections: req.Corrections,
		Metadata: map[string]interface{}{
			"user_role":   req.UserRole,
			"was_accurate": req.WasAccurate,
		},
		CreatedAt: time.Now(),
	}
	
	// Store in database
	if err := c.storeFeedback(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to store human feedback: %w", err)
	}
	
	return entry, nil
}

// CollectMachineFeedback processes implicit machine-generated feedback
func (c *Collector) CollectMachineFeedback(ctx context.Context, req *MachineFeedbackRequest) (*FeedbackEntry, error) {
	// Determine feedback type based on service
	feedbackType := c.determineFeedbackType(req.ServiceType, false)
	
	// Analyze outcomes to determine sentiment
	sentiment := c.analyzeOutcomesSentiment(&req.Outcomes)
	
	// Convert outcomes to map
	outcomesMap, err := structToMap(req.Outcomes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert outcomes: %w", err)
	}
	
	// Create feedback entry
	entry := &FeedbackEntry{
		Source:      SourceMachine,
		Type:        feedbackType,
		RequestID:   &req.RequestID,
		TicketID:    &req.TicketID,
		ServiceType: req.ServiceType,
		Sentiment:   sentiment,
		Outcomes:    outcomesMap,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
	}
	
	// Store in database
	if err := c.storeFeedback(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to store machine feedback: %w", err)
	}
	
	return entry, nil
}

// CollectTicketCompletionFeedback automatically collects feedback when ticket closes
func (c *Collector) CollectTicketCompletionFeedback(ctx context.Context, ticketID int64) error {
	// Fetch all AI requests for this ticket
	requests, err := c.getTicketAIRequests(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket AI requests: %w", err)
	}
	
	// Fetch ticket outcomes
	outcomes, err := c.getTicketOutcomes(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket outcomes: %w", err)
	}
	
	// Create machine feedback for each AI service used
	for serviceType, requestID := range requests {
		feedback := &MachineFeedbackRequest{
			ServiceType: serviceType,
			RequestID:   requestID,
			TicketID:    ticketID,
			Outcomes:    outcomes,
			Metadata: map[string]interface{}{
				"auto_collected": true,
				"collection_source": "ticket_completion",
			},
		}
		
		_, err := c.CollectMachineFeedback(ctx, feedback)
		if err != nil {
			// Log error but continue processing other services
			fmt.Printf("Error collecting feedback for %s: %v\n", serviceType, err)
		}
	}
	
	return nil
}

// GetFeedbackByRequest retrieves all feedback for a specific AI request
func (c *Collector) GetFeedbackByRequest(ctx context.Context, serviceType, requestID string) ([]*FeedbackEntry, error) {
	query := `
		SELECT 
			feedback_id,
			source,
			type,
			ticket_id,
			request_id,
			service_type,
			user_id,
			rating,
			sentiment,
			comments,
			outcomes,
			corrections,
			metadata,
			created_at,
			processed_at
		FROM ai_feedback
		WHERE service_type = $1 AND request_id = $2
		ORDER BY created_at DESC
	`
	
	rows, err := c.db.QueryContext(ctx, query, serviceType, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*FeedbackEntry
	for rows.Next() {
		entry, err := c.scanFeedbackEntry(rows)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// GetFeedbackByTicket retrieves all feedback for a ticket
func (c *Collector) GetFeedbackByTicket(ctx context.Context, ticketID int64) ([]*FeedbackEntry, error) {
	query := `
		SELECT 
			feedback_id,
			source,
			type,
			ticket_id,
			request_id,
			service_type,
			user_id,
			rating,
			sentiment,
			comments,
			outcomes,
			corrections,
			metadata,
			created_at,
			processed_at
		FROM ai_feedback
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := c.db.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*FeedbackEntry
	for rows.Next() {
		entry, err := c.scanFeedbackEntry(rows)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// GetRecentFeedback retrieves recent feedback for analysis
func (c *Collector) GetRecentFeedback(ctx context.Context, serviceType string, days int) ([]*FeedbackEntry, error) {
	query := `
		SELECT 
			feedback_id,
			source,
			type,
			ticket_id,
			request_id,
			service_type,
			user_id,
			rating,
			sentiment,
			comments,
			outcomes,
			corrections,
			metadata,
			created_at,
			processed_at
		FROM ai_feedback
		WHERE service_type = $1 
			AND created_at > NOW() - INTERVAL '$2 days'
		ORDER BY created_at DESC
	`
	
	rows, err := c.db.QueryContext(ctx, query, serviceType, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var entries []*FeedbackEntry
	for rows.Next() {
		entry, err := c.scanFeedbackEntry(rows)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	
	return entries, nil
}

// storeFeedback stores a feedback entry in the database
func (c *Collector) storeFeedback(ctx context.Context, entry *FeedbackEntry) error {
	outcomesJSON, _ := json.Marshal(entry.Outcomes)
	correctionsJSON, _ := json.Marshal(entry.Corrections)
	metadataJSON, _ := json.Marshal(entry.Metadata)
	
	query := `
		INSERT INTO ai_feedback (
			source,
			type,
			ticket_id,
			request_id,
			service_type,
			user_id,
			rating,
			sentiment,
			comments,
			outcomes,
			corrections,
			metadata,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING feedback_id
	`
	
	err := c.db.QueryRowContext(
		ctx,
		query,
		entry.Source,
		entry.Type,
		entry.TicketID,
		entry.RequestID,
		entry.ServiceType,
		entry.UserID,
		entry.Rating,
		entry.Sentiment,
		entry.Comments,
		outcomesJSON,
		correctionsJSON,
		metadataJSON,
		entry.CreatedAt,
	).Scan(&entry.FeedbackID)
	
	return err
}

// scanFeedbackEntry scans a database row into a FeedbackEntry
func (c *Collector) scanFeedbackEntry(row interface {
	Scan(dest ...interface{}) error
}) (*FeedbackEntry, error) {
	var entry FeedbackEntry
	var outcomesJSON, correctionsJSON, metadataJSON []byte
	var processedAt sql.NullTime
	
	err := row.Scan(
		&entry.FeedbackID,
		&entry.Source,
		&entry.Type,
		&entry.TicketID,
		&entry.RequestID,
		&entry.ServiceType,
		&entry.UserID,
		&entry.Rating,
		&entry.Sentiment,
		&entry.Comments,
		&outcomesJSON,
		&correctionsJSON,
		&metadataJSON,
		&entry.CreatedAt,
		&processedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if processedAt.Valid {
		entry.ProcessedAt = &processedAt.Time
	}
	
	json.Unmarshal(outcomesJSON, &entry.Outcomes)
	json.Unmarshal(correctionsJSON, &entry.Corrections)
	json.Unmarshal(metadataJSON, &entry.Metadata)
	
	return &entry, nil
}

// determineFeedbackType determines the feedback type based on service and source
func (c *Collector) determineFeedbackType(serviceType string, isHuman bool) FeedbackType {
	if isHuman {
		switch serviceType {
		case "diagnosis":
			return FeedbackDiagnosisAccuracy
		case "assignment":
			return FeedbackAssignmentAcceptance
		case "parts":
			return FeedbackPartsAccuracy
		default:
			return FeedbackGeneral
		}
	} else {
		switch serviceType {
		case "diagnosis":
			return FeedbackDiagnosisCorrection
		case "assignment":
			return FeedbackAssignmentPerformance
		case "parts":
			return FeedbackPartsUsed
		default:
			return FeedbackGeneral
		}
	}
}

// determineSentiment determines sentiment based on rating and accuracy
func (c *Collector) determineSentiment(rating *int, wasAccurate bool) FeedbackSentiment {
	if rating != nil {
		if *rating >= 4 {
			return SentimentPositive
		} else if *rating == 3 {
			return SentimentNeutral
		} else {
			return SentimentNegative
		}
	}
	
	if wasAccurate {
		return SentimentPositive
	}
	
	return SentimentNegative
}

// analyzeOutcomesSentiment analyzes machine outcomes to determine sentiment
func (c *Collector) analyzeOutcomesSentiment(outcomes *MachineFeedbackOutcomes) FeedbackSentiment {
	positiveSignals := 0
	negativeSignals := 0
	
	// Diagnosis outcomes
	if outcomes.DiagnosisMatchedAI != nil {
		if *outcomes.DiagnosisMatchedAI {
			positiveSignals++
		} else {
			negativeSignals++
		}
	}
	
	// Assignment outcomes
	if outcomes.AssignmentAccepted != nil {
		if *outcomes.AssignmentAccepted {
			positiveSignals++
		} else {
			negativeSignals++
		}
	}
	
	if outcomes.WasTopRecommendation != nil && *outcomes.WasTopRecommendation {
		positiveSignals++
	}
	
	// Parts outcomes
	if len(outcomes.RecommendedPartsUsed) > 0 {
		positiveSignals++
	}
	if len(outcomes.UnrecommendedPartsUsed) > 0 {
		negativeSignals++
	}
	
	// Resolution outcomes
	if outcomes.FirstTimeFixRate != nil {
		if *outcomes.FirstTimeFixRate {
			positiveSignals++
		} else {
			negativeSignals++
		}
	}
	
	if outcomes.CustomerSatisfaction != nil {
		if *outcomes.CustomerSatisfaction >= 4 {
			positiveSignals++
		} else if *outcomes.CustomerSatisfaction <= 2 {
			negativeSignals++
		}
	}
	
	// Determine overall sentiment
	if positiveSignals > negativeSignals {
		return SentimentPositive
	} else if positiveSignals < negativeSignals {
		return SentimentNegative
	}
	
	return SentimentNeutral
}

// getTicketAIRequests retrieves all AI request IDs for a ticket
func (c *Collector) getTicketAIRequests(ctx context.Context, ticketID int64) (map[string]string, error) {
	requests := make(map[string]string)
	
	// Get diagnosis request
	var diagnosisRequestID sql.NullString
	c.db.QueryRowContext(ctx, `
		SELECT request_id FROM ai_diagnoses WHERE ticket_id = $1 ORDER BY created_at DESC LIMIT 1
	`, ticketID).Scan(&diagnosisRequestID)
	if diagnosisRequestID.Valid {
		requests["diagnosis"] = diagnosisRequestID.String
	}
	
	// Get assignment request
	var assignmentRequestID sql.NullString
	c.db.QueryRowContext(ctx, `
		SELECT request_id FROM assignment_history WHERE ticket_id = $1 ORDER BY created_at DESC LIMIT 1
	`, ticketID).Scan(&assignmentRequestID)
	if assignmentRequestID.Valid {
		requests["assignment"] = assignmentRequestID.String
	}
	
	// Get parts request
	var partsRequestID sql.NullString
	c.db.QueryRowContext(ctx, `
		SELECT request_id FROM parts_recommendations WHERE ticket_id = $1 ORDER BY created_at DESC LIMIT 1
	`, ticketID).Scan(&partsRequestID)
	if partsRequestID.Valid {
		requests["parts"] = partsRequestID.String
	}
	
	return requests, nil
}

// getTicketOutcomes retrieves actual outcomes from a completed ticket
func (c *Collector) getTicketOutcomes(ctx context.Context, ticketID int64) (MachineFeedbackOutcomes, error) {
	var outcomes MachineFeedbackOutcomes
	
	// Get ticket details
	var (
		actualProblem      sql.NullString
		resolutionTime     sql.NullInt64
		customerSat        sql.NullInt64
		actualCost         sql.NullFloat64
		assignedEngineerID sql.NullInt64
	)
	
	err := c.db.QueryRowContext(ctx, `
		SELECT 
			problem_description,
			EXTRACT(EPOCH FROM (resolved_at - created_at))/60 as resolution_minutes,
			customer_satisfaction_rating,
			total_cost,
			assigned_engineer_id
		FROM service_tickets
		WHERE ticket_id = $1
	`, ticketID).Scan(&actualProblem, &resolutionTime, &customerSat, &actualCost, &assignedEngineerID)
	
	if err != nil && err != sql.ErrNoRows {
		return outcomes, err
	}
	
	if actualProblem.Valid {
		outcomes.ActualProblem = &actualProblem.String
	}
	if resolutionTime.Valid {
		resolutionTimeInt := int(resolutionTime.Int64)
		outcomes.ResolutionTime = &resolutionTimeInt
	}
	if customerSat.Valid {
		customerSatInt := int(customerSat.Int64)
		outcomes.CustomerSatisfaction = &customerSatInt
	}
	if actualCost.Valid {
		outcomes.ActualCost = &actualCost.Float64
	}
	if assignedEngineerID.Valid {
		outcomes.AssignedEngineerID = &assignedEngineerID.Int64
	}
	
	// Get parts used
	rows, err := c.db.QueryContext(ctx, `
		SELECT part_id, was_recommended, cost
		FROM ticket_parts
		WHERE ticket_id = $1
	`, ticketID)
	
	if err == nil {
		defer rows.Close()
		
		var allParts []int64
		var recommendedParts []int64
		var unrecommendedParts []int64
		
		for rows.Next() {
			var partID int64
			var wasRecommended bool
			var cost sql.NullFloat64
			
			rows.Scan(&partID, &wasRecommended, &cost)
			allParts = append(allParts, partID)
			
			if wasRecommended {
				recommendedParts = append(recommendedParts, partID)
			} else {
				unrecommendedParts = append(unrecommendedParts, partID)
			}
		}
		
		outcomes.PartsUsedIDs = allParts
		outcomes.RecommendedPartsUsed = recommendedParts
		outcomes.UnrecommendedPartsUsed = unrecommendedParts
	}
	
	return outcomes, nil
}

// structToMap converts a struct to map[string]interface{}
func structToMap(s interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// MarkProcessed marks feedback as processed
func (c *Collector) MarkProcessed(ctx context.Context, feedbackID int64) error {
	query := `
		UPDATE ai_feedback
		SET processed_at = NOW()
		WHERE feedback_id = $1
	`
	
	_, err := c.db.ExecContext(ctx, query, feedbackID)
	return err
}


