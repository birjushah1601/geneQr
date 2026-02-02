package assignment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/google/uuid"
)

const assignmentEngineVersion = "1.0.0"

// Engine handles intelligent engineer assignment
type Engine struct {
	aiManager *ai.Manager
	scorer    *Scorer
	db        *sql.DB
}

// NewEngine creates a new assignment engine
func NewEngine(aiManager *ai.Manager, db *sql.DB) *Engine {
	return &Engine{
		aiManager: aiManager,
		scorer:    NewScorer(db),
		db:        db,
	}
}

// RecommendEngineers returns ranked engineer recommendations for a ticket
func (e *Engine) RecommendEngineers(ctx context.Context, req *AssignmentRequest) (*AssignmentResponse, error) {
	startTime := time.Now()

	// Set defaults
	if req.MaxRecommendations == 0 {
		req.MaxRecommendations = 5
	}
	if req.MinScoreThreshold == 0 {
		req.MinScoreThreshold = 40.0
	}

	// Generate request ID
	requestID := uuid.New().String()

	// Step 1: Load available engineers
	engineers, err := e.loadEngineers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to load engineers: %w", err)
	}

	if len(engineers) == 0 {
		return &AssignmentResponse{
			RequestID:       requestID,
			TicketID:        req.TicketID,
			Recommendations: []EngineerRecommendation{},
			Metadata: AssignmentMetadata{
				EngineersEvaluated: 0,
				EngineersFiltered:  0,
				ProcessingTime:     time.Since(startTime),
				Version:            assignmentEngineVersion,
			},
			CreatedAt: time.Now(),
		}, nil
	}

	// Step 2: Score each engineer
	var recommendations []EngineerRecommendation
	for _, engineer := range engineers {
		rec, err := e.scoreAndBuildRecommendation(ctx, &engineer, req)
		if err != nil {
			// Log error but continue with other engineers
			continue
		}

		// Filter by minimum threshold
		if rec.OverallScore >= req.MinScoreThreshold {
			recommendations = append(recommendations, *rec)
		}
	}

	// Step 3: Sort by score (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].OverallScore > recommendations[j].OverallScore
	})

	// Step 4: Optionally use AI for final ranking adjustment
	var aiCost float64
	var aiProvider, aiModel string
	if req.Options.UseAI && len(recommendations) > 1 {
		// Use AI to review and potentially adjust rankings
		adjusted, cost, provider, model, err := e.aiRankingAdjustment(ctx, req, recommendations)
		if err == nil && adjusted != nil {
			recommendations = adjusted
			aiCost = cost
			aiProvider = provider
			aiModel = model
		}
	}

	// Step 5: Assign ranks and limit to max recommendations
	if len(recommendations) > req.MaxRecommendations {
		recommendations = recommendations[:req.MaxRecommendations]
	}

	for i := range recommendations {
		recommendations[i].Rank = i + 1
	}

	// Step 6: Save assignment history
	response := &AssignmentResponse{
		RequestID:       requestID,
		TicketID:        req.TicketID,
		Recommendations: recommendations,
		Metadata: AssignmentMetadata{
			EngineersEvaluated: len(engineers),
			EngineersFiltered:  len(recommendations),
			UsedAI:             req.Options.UseAI && aiCost > 0,
			AIProvider:         aiProvider,
			AIModel:            aiModel,
			CostUSD:            aiCost,
			ProcessingTime:     time.Since(startTime),
			Version:            assignmentEngineVersion,
		},
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := e.saveAssignmentHistory(ctx, response); err != nil {
		// Log but don't fail the request
		fmt.Printf("Warning: failed to save assignment history: %v\n", err)
	}

	return response, nil
}

// scoreAndBuildRecommendation scores an engineer and builds recommendation
func (e *Engine) scoreAndBuildRecommendation(ctx context.Context, engineer *EngineerProfile, req *AssignmentRequest) (*EngineerRecommendation, error) {
	// Calculate scores
	breakdown, reasons, warnings, err := e.scorer.ScoreEngineer(ctx, engineer, req)
	if err != nil {
		return nil, err
	}

	// Calculate overall score
	overallScore := e.scorer.CalculateOverallScore(breakdown)

	// Find matching and missing skills
	matchingSkills, missingSkills := e.analyzeSkills(engineer, req)

	// Build availability info
	availInfo := e.buildAvailabilityInfo(engineer)

	rec := &EngineerRecommendation{
		EngineerID:            engineer.EngineerID,
		EngineerName:          engineer.FullName,
		Email:                 engineer.Email,
		Phone:                 engineer.Phone,
		OverallScore:          overallScore,
		Scores:                *breakdown,
		MatchReasons:          reasons,
		Warnings:              warnings,
		CurrentWorkload:       engineer.OpenTicketsCount,
		AvailabilityInfo:      &availInfo,
		CurrentLocation:       engineer.CurrentLocation,
		MatchingSkills:        matchingSkills,
		MissingSkills:         missingSkills,
		AverageResolutionTime: &engineer.AverageResolutionTime,
		SuccessRate:           &engineer.SuccessRate,
		RecentTicketsCount:    len(engineer.RecentTickets),
	}

	return rec, nil
}

// analyzeSkills finds matching and missing skills
func (e *Engine) analyzeSkills(engineer *EngineerProfile, req *AssignmentRequest) ([]string, []string) {
	var matching []string
	var missing []string

	for _, reqSkill := range req.RequiredSkills {
		found := false
		for _, engSkill := range engineer.Skills {
			if engSkill.SkillName == reqSkill {
				matching = append(matching, reqSkill)
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, reqSkill)
		}
	}

	return matching, missing
}

// buildAvailabilityInfo creates availability description
func (e *Engine) buildAvailabilityInfo(engineer *EngineerProfile) string {
	switch engineer.AvailabilityStatus {
	case AvailabilityAvailable:
		if engineer.OpenTicketsCount == 0 {
			return "Available now - no current assignments"
		}
		return fmt.Sprintf("Available now - %d ticket(s) in progress", engineer.OpenTicketsCount)
	case AvailabilityBusy:
		return "Busy with current ticket - can be assigned as backup"
	case AvailabilityOnLeave:
		return "On leave - not available"
	case AvailabilityOffline:
		return "Offline - availability unknown"
	default:
		return "Status unknown"
	}
}

// aiRankingAdjustment uses AI to review and adjust rankings
func (e *Engine) aiRankingAdjustment(ctx context.Context, req *AssignmentRequest, recommendations []EngineerRecommendation) ([]EngineerRecommendation, float64, string, string, error) {
	// Build context for AI
	prompt := e.buildRankingPrompt(req, recommendations)

	// Call AI
	result, err := e.aiManager.Chat(ctx, &ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: "You are an expert in medical equipment service management. Review engineer assignments and provide ranking adjustments based on context.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: func() *float32 { v := float32(0.3); return &v }(),
		MaxTokens:   func() *int { v := 1000; return &v }(),
	})

	if err != nil {
		return nil, 0, "", "", err
	}

	// Parse AI response for ranking adjustments
	adjusted := e.parseAIRankingAdjustments(result.Content, recommendations)

	return adjusted, result.Cost, result.Provider, result.Model, nil
}

// buildRankingPrompt creates AI prompt for ranking review
func (e *Engine) buildRankingPrompt(req *AssignmentRequest, recommendations []EngineerRecommendation) string {
	prompt := fmt.Sprintf(`Review these engineer recommendations for a service ticket:

**Ticket Details:**
- Equipment: %s
- Priority: %s
- Problem: %s
- Location: %s

**Current Top Recommendations (by algorithm):**
`, req.EquipmentType, req.Priority, req.ProblemType, req.LocationName)

	for i, rec := range recommendations {
		if i >= 5 {
			break // Only show top 5
		}
		prompt += fmt.Sprintf(`
%d. %s (Score: %.1f)
   - Expertise: %.1f, Location: %.1f, Performance: %.1f
   - Workload: %d tickets
   - Availability: %s
   - Key strength: %s
`,
			i+1,
			rec.EngineerName,
			rec.OverallScore,
			rec.Scores.ExpertiseScore,
			rec.Scores.LocationScore,
			rec.Scores.PerformanceScore,
			rec.CurrentWorkload,
			*rec.AvailabilityInfo,
			getTopReason(rec.MatchReasons),
		)
	}

	prompt += `
**Task:** Review these rankings considering:
1. Does the priority level (%s) suggest different weightings?
2. Are there any red flags (overloaded, on leave, poor recent performance)?
3. Would you swap any rankings based on the context?

Respond in JSON format:
{
  "recommended_changes": [
    {"from_rank": 1, "to_rank": 2, "reason": "..."}
  ],
  "concerns": ["..."],
  "endorsement": "Agree with ranking" or "Suggest adjustment"
}
`
	return fmt.Sprintf(prompt, req.Priority)
}

// parseAIRankingAdjustments parses AI response and adjusts rankings
func (e *Engine) parseAIRankingAdjustments(aiResponse string, recommendations []EngineerRecommendation) []EngineerRecommendation {
	// Parse JSON response
	var response struct {
		RecommendedChanges []struct {
			FromRank int    `json:"from_rank"`
			ToRank   int    `json:"to_rank"`
			Reason   string `json:"reason"`
		} `json:"recommended_changes"`
		Concerns    []string `json:"concerns"`
		Endorsement string   `json:"endorsement"`
	}

	if err := json.Unmarshal([]byte(aiResponse), &response); err != nil {
		// If parsing fails, return original rankings
		return recommendations
	}

	// Apply changes
	adjusted := make([]EngineerRecommendation, len(recommendations))
	copy(adjusted, recommendations)

	for _, change := range response.RecommendedChanges {
		if change.FromRank > 0 && change.ToRank > 0 &&
			change.FromRank <= len(adjusted) && change.ToRank <= len(adjusted) {
			// Swap positions
			fromIdx := change.FromRank - 1
			toIdx := change.ToRank - 1
			adjusted[fromIdx], adjusted[toIdx] = adjusted[toIdx], adjusted[fromIdx]

			// Add AI reason to warnings
			adjusted[toIdx].Warnings = append(adjusted[toIdx].Warnings,
				fmt.Sprintf("AI adjustment: %s", change.Reason))
		}
	}

	return adjusted
}

// getTopReason gets the highest impact reason
func getTopReason(reasons []MatchReason) string {
	for _, r := range reasons {
		if r.Impact == ImpactHigh {
			return r.Reason
		}
	}
	if len(reasons) > 0 {
		return reasons[0].Reason
	}
	return "General match"
}

// loadEngineers loads available engineers from database
func (e *Engine) loadEngineers(ctx context.Context, req *AssignmentRequest) ([]EngineerProfile, error) {
	query := `
		SELECT 
			u.user_id,
			u.full_name,
			u.email,
			u.phone,
			u.is_active,
			u.hire_date,
			u.department,
			u.specialization,
			-- Add aggregated metrics
			COALESCE(COUNT(DISTINCT st.ticket_id), 0) as open_tickets
		FROM users u
		LEFT JOIN service_tickets st ON st.assigned_to = u.user_id 
			AND st.status IN ('Open', 'In Progress', 'On Hold')
		WHERE u.role = 'Engineer' 
			AND u.is_active = true
		GROUP BY u.user_id, u.full_name, u.email, u.phone, 
			u.is_active, u.hire_date, u.department, u.specialization
		ORDER BY u.full_name
	`

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var engineers []EngineerProfile
	for rows.Next() {
		var eng EngineerProfile
		var hireDate time.Time
		var openTickets int

		err := rows.Scan(
			&eng.EngineerID,
			&eng.FullName,
			&eng.Email,
			&eng.Phone,
			&eng.IsActive,
			&hireDate,
			&eng.Department,
			&eng.Specialization,
			&openTickets,
		)
		if err != nil {
			continue
		}

		eng.HireDate = hireDate
		eng.OpenTicketsCount = openTickets
		eng.CurrentWorkload = openTickets

		// Load additional details
		e.enrichEngineerProfile(ctx, &eng)

		engineers = append(engineers, eng)
	}

	return engineers, nil
}

// enrichEngineerProfile adds skills, expertise, and performance metrics
func (e *Engine) enrichEngineerProfile(ctx context.Context, engineer *EngineerProfile) {
	// Load equipment expertise from ticket history
	expertiseQuery := `
		SELECT 
			et.equipment_name,
			COUNT(*) as tickets_handled,
			AVG(CASE WHEN st.status = 'Resolved' THEN 1.0 ELSE 0.0 END) as success_rate,
			AVG(EXTRACT(EPOCH FROM (st.resolved_at - st.created_at))/3600.0) as avg_hours
		FROM service_tickets st
		JOIN equipment_types et ON st.equipment_type_id = et.equipment_type_id
		WHERE st.assigned_to = $1
			AND st.created_at > NOW() - INTERVAL '1 year'
		GROUP BY et.equipment_name
		ORDER BY tickets_handled DESC
	`

	rows, err := e.db.QueryContext(ctx, expertiseQuery, engineer.EngineerID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var exp EquipmentExpertise
			var avgHours float64
			rows.Scan(&exp.EquipmentType, &exp.TicketsHandled, &exp.SuccessRate, &avgHours)
			exp.AverageResolutionTime = time.Duration(avgHours * float64(time.Hour))
			engineer.EquipmentExpertise = append(engineer.EquipmentExpertise, exp)
		}
	}

	// Load overall performance metrics
	perfQuery := `
		SELECT 
			COUNT(*) as total_tickets,
			AVG(CASE WHEN status = 'Resolved' THEN 1.0 ELSE 0.0 END) as success_rate,
			AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))/3600.0) as avg_hours
		FROM service_tickets
		WHERE assigned_to = $1
			AND created_at > NOW() - INTERVAL '1 year'
	`

	var totalTickets int
	var successRate float64
	var avgHours *float64

	err = e.db.QueryRowContext(ctx, perfQuery, engineer.EngineerID).Scan(
		&totalTickets,
		&successRate,
		&avgHours,
	)

	if err == nil {
		engineer.TotalTicketsResolved = totalTickets
		engineer.SuccessRate = successRate
		if avgHours != nil {
			engineer.AverageResolutionTime = time.Duration(*avgHours * float64(time.Hour))
		}
	}

	// Set default availability (in production, track this in database)
	engineer.AvailabilityStatus = AvailabilityAvailable
	if engineer.OpenTicketsCount > 5 {
		engineer.AvailabilityStatus = AvailabilityBusy
	}
}

// saveAssignmentHistory saves recommendations to database
func (e *Engine) saveAssignmentHistory(ctx context.Context, response *AssignmentResponse) error {
	recsJSON, err := json.Marshal(response.Recommendations)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(response.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO assignment_history 
		(request_id, ticket_id, recommendations, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = e.db.ExecContext(
		ctx,
		query,
		response.RequestID,
		response.TicketID,
		recsJSON,
		metadataJSON,
		response.CreatedAt,
	)

	return err
}

// GetAssignmentHistory retrieves assignment history for a ticket
func (e *Engine) GetAssignmentHistory(ctx context.Context, ticketID int64) ([]*AssignmentResponse, error) {
	query := `
		SELECT 
			request_id,
			ticket_id,
			recommendations,
			metadata,
			selected_engineer_id,
			was_successful,
			created_at
		FROM assignment_history
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := e.db.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*AssignmentResponse
	for rows.Next() {
		var resp AssignmentResponse
		var recsJSON, metadataJSON []byte
		var selectedEngineerID *int64
		var wasSuccessful *bool

		err := rows.Scan(
			&resp.RequestID,
			&resp.TicketID,
			&recsJSON,
			&metadataJSON,
			&selectedEngineerID,
			&wasSuccessful,
			&resp.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(recsJSON, &resp.Recommendations)
		json.Unmarshal(metadataJSON, &resp.Metadata)

		history = append(history, &resp)
	}

	return history, nil
}



