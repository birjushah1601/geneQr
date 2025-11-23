package diagnosis

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ContextEnricher enriches diagnosis requests with historical data
type ContextEnricher struct {
	db *sql.DB
}

// NewContextEnricher creates a new context enricher
func NewContextEnricher(db *sql.DB) *ContextEnricher {
	return &ContextEnricher{
		db: db,
	}
}

// EnrichedContext represents enriched context for diagnosis
type EnrichedContext struct {
	// EquipmentHistory past issues with this equipment
	EquipmentHistory []HistoricalIssue

	// SimilarTickets similar past tickets
	SimilarTickets []SimilarTicketInfo

	// ManufacturerInfo manufacturer-specific information
	ManufacturerInfo *ManufacturerInfo

	// KnownIssues known issues for this equipment type
	KnownIssues []KnownIssue

	// LocationContext location-specific context
	LocationContext *LocationContextInfo
}

// HistoricalIssue represents a past issue with the equipment
type HistoricalIssue struct {
	TicketID       int64
	IssueDate      time.Time
	ProblemType    string
	Description    string
	Resolution     string
	ResolutionTime time.Duration
	PartsReplaced  []string
	AssignedTo     string
	WasSuccessful  bool
}

// ManufacturerInfo represents manufacturer-specific information
type ManufacturerInfo struct {
	ManufacturerCode string
	ManufacturerName string
	ContactInfo      string
	CommonIssues     []string
	MaintenanceTips  []string
	WarrantyInfo     *string
}

// KnownIssue represents a known issue for equipment type
type KnownIssue struct {
	IssueID         int64
	EquipmentType   string
	IssueTitle      string
	Description     string
	Symptoms        []string
	Solution        string
	Frequency       int // How often this occurs
	Severity        string
	RequiresParts   []string
	EstimatedTime   time.Duration
}

// LocationContextInfo represents location-specific context
type LocationContextInfo struct {
	LocationType         string
	TypicalIssues        []string
	EnvironmentalFactors []string
	UsageIntensity       string // High, Medium, Low
}

// Enrich enriches a diagnosis request with context
func (ce *ContextEnricher) Enrich(ctx context.Context, req *DiagnosisRequest) (*EnrichedContext, error) {
	enriched := &EnrichedContext{}

	// Get equipment history if equipment ID is known
	if req.EquipmentID != nil && req.Options.IncludeHistoricalContext {
		history, err := ce.getEquipmentHistory(ctx, *req.EquipmentID)
		if err != nil {
			// Log but don't fail
			fmt.Printf("Warning: failed to get equipment history: %v\n", err)
		} else {
			enriched.EquipmentHistory = history
		}
	}

	// Find similar tickets
	if req.Options.IncludeSimilarTickets {
		similar, err := ce.findSimilarTickets(ctx, req)
		if err != nil {
			fmt.Printf("Warning: failed to find similar tickets: %v\n", err)
		} else {
			enriched.SimilarTickets = similar
		}
	}

	// Get manufacturer info
	if req.Manufacturer != nil {
		mfgInfo, err := ce.getManufacturerInfo(ctx, *req.Manufacturer)
		if err != nil {
			fmt.Printf("Warning: failed to get manufacturer info: %v\n", err)
		} else {
			enriched.ManufacturerInfo = mfgInfo
		}
	}

	// Get known issues for equipment type
	knownIssues, err := ce.getKnownIssues(ctx, req.EquipmentType)
	if err != nil {
		fmt.Printf("Warning: failed to get known issues: %v\n", err)
	} else {
		enriched.KnownIssues = knownIssues
	}

	// Get location context
	if req.LocationType != nil {
		locContext := ce.getLocationContext(*req.LocationType)
		enriched.LocationContext = locContext
	}

	return enriched, nil
}

// getEquipmentHistory retrieves historical issues for equipment
func (ce *ContextEnricher) getEquipmentHistory(ctx context.Context, equipmentID int64) ([]HistoricalIssue, error) {
	query := `
		SELECT 
			t.ticket_id,
			t.created_at,
			t.problem_type,
			t.description,
			t.resolution_notes,
			EXTRACT(EPOCH FROM (t.closed_at - t.created_at)) as resolution_seconds,
			COALESCE(array_agg(DISTINCT p.part_name) FILTER (WHERE p.part_name IS NOT NULL), ARRAY[]::text[]) as parts,
			COALESCE(e.full_name, 'Unknown') as engineer_name,
			CASE WHEN t.status = 'closed' AND t.resolution_notes IS NOT NULL THEN true ELSE false END as successful
		FROM service_tickets t
		LEFT JOIN ticket_parts tp ON t.ticket_id = tp.ticket_id
		LEFT JOIN parts p ON tp.part_id = p.part_id
		LEFT JOIN engineers e ON t.assigned_engineer_id = e.engineer_id
		WHERE t.equipment_id = $1
		  AND t.created_at > NOW() - INTERVAL '1 year'
		  AND t.status IN ('closed', 'resolved')
		GROUP BY t.ticket_id, t.created_at, t.problem_type, t.description, 
		         t.resolution_notes, t.closed_at, e.full_name, t.status
		ORDER BY t.created_at DESC
		LIMIT 10
	`

	rows, err := ce.db.QueryContext(ctx, query, equipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query equipment history: %w", err)
	}
	defer rows.Close()

	var history []HistoricalIssue
	for rows.Next() {
		var issue HistoricalIssue
		var resolutionSeconds *float64
		var parts []string

		err := rows.Scan(
			&issue.TicketID,
			&issue.IssueDate,
			&issue.ProblemType,
			&issue.Description,
			&issue.Resolution,
			&resolutionSeconds,
			&parts,
			&issue.AssignedTo,
			&issue.WasSuccessful,
		)
		if err != nil {
			continue
		}

		if resolutionSeconds != nil {
			issue.ResolutionTime = time.Duration(*resolutionSeconds) * time.Second
		}
		issue.PartsReplaced = parts

		history = append(history, issue)
	}

	return history, nil
}

// findSimilarTickets finds similar past tickets using text similarity
func (ce *ContextEnricher) findSimilarTickets(ctx context.Context, req *DiagnosisRequest) ([]SimilarTicketInfo, error) {
	// Use PostgreSQL's full-text search or similarity functions
	query := `
		SELECT 
			t.ticket_id,
			t.description,
			t.resolution_notes,
			EXTRACT(EPOCH FROM (t.closed_at - t.created_at)) as resolution_seconds,
			similarity(t.description, $1) as sim_score
		FROM service_tickets t
		WHERE t.equipment_type_id = (SELECT equipment_type_id FROM equipment_types WHERE name = $2)
		  AND t.status IN ('closed', 'resolved')
		  AND t.resolution_notes IS NOT NULL
		  AND t.created_at > NOW() - INTERVAL '2 years'
		  AND similarity(t.description, $1) > 0.3
		ORDER BY sim_score DESC
		LIMIT $3
	`

	maxTickets := req.Options.MaxSimilarTickets
	if maxTickets == 0 {
		maxTickets = 5
	}

	rows, err := ce.db.QueryContext(ctx, query, req.Description, req.EquipmentType, maxTickets)
	if err != nil {
		// If similarity function not available, fall back to simple search
		return ce.findSimilarTicketsSimple(ctx, req)
	}
	defer rows.Close()

	var similar []SimilarTicketInfo
	for rows.Next() {
		var ticket SimilarTicketInfo
		var resolutionSeconds *float64
		var simScore float64

		err := rows.Scan(
			&ticket.TicketID,
			&ticket.Description,
			&ticket.Resolution,
			&resolutionSeconds,
			&simScore,
		)
		if err != nil {
			continue
		}

		if resolutionSeconds != nil {
			ticket.TimeTaken = time.Duration(*resolutionSeconds) * time.Second
		}
		ticket.Similarity = simScore

		similar = append(similar, ticket)
	}

	return similar, nil
}

// findSimilarTicketsSimple simple fallback for finding similar tickets
func (ce *ContextEnricher) findSimilarTicketsSimple(ctx context.Context, req *DiagnosisRequest) ([]SimilarTicketInfo, error) {
	query := `
		SELECT 
			t.ticket_id,
			t.description,
			t.resolution_notes,
			EXTRACT(EPOCH FROM (t.closed_at - t.created_at)) as resolution_seconds
		FROM service_tickets t
		WHERE t.equipment_type_id = (SELECT equipment_type_id FROM equipment_types WHERE name = $1)
		  AND t.status IN ('closed', 'resolved')
		  AND t.resolution_notes IS NOT NULL
		  AND t.created_at > NOW() - INTERVAL '1 year'
		ORDER BY t.created_at DESC
		LIMIT $2
	`

	maxTickets := req.Options.MaxSimilarTickets
	if maxTickets == 0 {
		maxTickets = 5
	}

	rows, err := ce.db.QueryContext(ctx, query, req.EquipmentType, maxTickets)
	if err != nil {
		return nil, fmt.Errorf("failed to query similar tickets: %w", err)
	}
	defer rows.Close()

	var similar []SimilarTicketInfo
	for rows.Next() {
		var ticket SimilarTicketInfo
		var resolutionSeconds *float64

		err := rows.Scan(
			&ticket.TicketID,
			&ticket.Description,
			&ticket.Resolution,
			&resolutionSeconds,
		)
		if err != nil {
			continue
		}

		if resolutionSeconds != nil {
			ticket.TimeTaken = time.Duration(*resolutionSeconds) * time.Second
		}
		ticket.Similarity = 0.5 // Default similarity

		similar = append(similar, ticket)
	}

	return similar, nil
}

// getManufacturerInfo retrieves manufacturer-specific information
func (ce *ContextEnricher) getManufacturerInfo(ctx context.Context, manufacturerCode string) (*ManufacturerInfo, error) {
	query := `
		SELECT 
			manufacturer_code,
			manufacturer_name,
			contact_info,
			common_issues,
			maintenance_tips
		FROM manufacturers
		WHERE manufacturer_code = $1
	`

	var info ManufacturerInfo
	var commonIssues, maintenanceTips []string

	err := ce.db.QueryRowContext(ctx, query, manufacturerCode).Scan(
		&info.ManufacturerCode,
		&info.ManufacturerName,
		&info.ContactInfo,
		&commonIssues,
		&maintenanceTips,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query manufacturer info: %w", err)
	}

	info.CommonIssues = commonIssues
	info.MaintenanceTips = maintenanceTips

	return &info, nil
}

// getKnownIssues retrieves known issues for equipment type
func (ce *ContextEnricher) getKnownIssues(ctx context.Context, equipmentType string) ([]KnownIssue, error) {
	// This would query a known_issues table if it exists
	// For now, return empty array
	return []KnownIssue{}, nil
}

// getLocationContext provides location-specific context
func (ce *ContextEnricher) getLocationContext(locationType string) *LocationContextInfo {
	// Hardcoded location context for now
	contexts := map[string]LocationContextInfo{
		"ICU": {
			LocationType: "ICU",
			TypicalIssues: []string{
				"High-frequency use wear",
				"Urgent response required",
				"Critical patient impact",
			},
			EnvironmentalFactors: []string{
				"24/7 operation",
				"High humidity from medical procedures",
				"Frequent cleaning with harsh chemicals",
			},
			UsageIntensity: "High",
		},
		"Ward": {
			LocationType: "Ward",
			TypicalIssues: []string{
				"Moderate use wear",
				"Standard maintenance needs",
			},
			EnvironmentalFactors: []string{
				"Regular business hours",
				"Standard hospital environment",
			},
			UsageIntensity: "Medium",
		},
		"Emergency": {
			LocationType: "Emergency",
			TypicalIssues: []string{
				"Rough handling",
				"Urgent repairs needed",
				"Equipment moved frequently",
			},
			EnvironmentalFactors: []string{
				"24/7 high-stress environment",
				"Frequent urgent use",
			},
			UsageIntensity: "High",
		},
	}

	if ctx, exists := contexts[locationType]; exists {
		return &ctx
	}

	return &LocationContextInfo{
		LocationType:   locationType,
		UsageIntensity: "Medium",
	}
}

// FormatContextForAI formats enriched context for AI prompt
func (ec *EnrichedContext) FormatContextForAI() string {
	var context string

	// Equipment history
	if len(ec.EquipmentHistory) > 0 {
		context += "\n## Equipment History:\n"
		for i, issue := range ec.EquipmentHistory {
			if i >= 5 {
				break // Limit to 5 most recent
			}
			context += fmt.Sprintf("- %s: %s (Resolved: %s)\n",
				issue.IssueDate.Format("2006-01-02"),
				issue.ProblemType,
				issue.Resolution)
			if len(issue.PartsReplaced) > 0 {
				context += fmt.Sprintf("  Parts: %v\n", issue.PartsReplaced)
			}
		}
	}

	// Similar tickets
	if len(ec.SimilarTickets) > 0 {
		context += "\n## Similar Past Tickets:\n"
		for i, ticket := range ec.SimilarTickets {
			if i >= 3 {
				break
			}
			context += fmt.Sprintf("- Ticket #%d (%.0f%% similar): %s\n",
				ticket.TicketID,
				ticket.Similarity*100,
				ticket.Description)
			if ticket.Resolution != "" {
				context += fmt.Sprintf("  Resolution: %s\n", ticket.Resolution)
			}
		}
	}

	// Manufacturer info
	if ec.ManufacturerInfo != nil {
		context += "\n## Manufacturer Information:\n"
		context += fmt.Sprintf("Manufacturer: %s\n", ec.ManufacturerInfo.ManufacturerName)
		if len(ec.ManufacturerInfo.CommonIssues) > 0 {
			context += "Common Issues:\n"
			for _, issue := range ec.ManufacturerInfo.CommonIssues {
				context += fmt.Sprintf("- %s\n", issue)
			}
		}
	}

	// Known issues
	if len(ec.KnownIssues) > 0 {
		context += "\n## Known Issues for Equipment Type:\n"
		for _, issue := range ec.KnownIssues {
			context += fmt.Sprintf("- %s: %s\n", issue.IssueTitle, issue.Description)
			if issue.Solution != "" {
				context += fmt.Sprintf("  Solution: %s\n", issue.Solution)
			}
		}
	}

	// Location context
	if ec.LocationContext != nil {
		context += "\n## Location Context:\n"
		context += fmt.Sprintf("Location Type: %s (Usage Intensity: %s)\n",
			ec.LocationContext.LocationType,
			ec.LocationContext.UsageIntensity)
		if len(ec.LocationContext.TypicalIssues) > 0 {
			context += "Typical Issues:\n"
			for _, issue := range ec.LocationContext.TypicalIssues {
				context += fmt.Sprintf("- %s\n", issue)
			}
		}
	}

	return context
}

