package diagnosis

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/google/uuid"
)

// Engine is the main diagnosis engine that orchestrates AI-powered diagnosis
type Engine struct {
	aiManager       *ai.Manager
	contextEnricher *ContextEnricher
	visionAnalyzer  *VisionAnalyzer
	db              *sql.DB
	version         string
}

// NewEngine creates a new diagnosis engine
func NewEngine(aiManager *ai.Manager, db *sql.DB) *Engine {
	return &Engine{
		aiManager:       aiManager,
		contextEnricher: NewContextEnricher(db),
		visionAnalyzer:  NewVisionAnalyzer(aiManager),
		db:              db,
		version:         "1.0.0",
	}
}

// Diagnose performs AI-powered diagnosis on a ticket
func (e *Engine) Diagnose(ctx context.Context, req *DiagnosisRequest) (*DiagnosisResponse, error) {
	startTime := time.Now()

	// Generate diagnosis ID
	diagnosisID := uuid.New().String()

	// Initialize response
	response := &DiagnosisResponse{
		DiagnosisID:        diagnosisID,
		TicketID:           req.TicketID,
		AlternateDiagnoses: []DiagnosisResult{},
		RecommendedActions: []RecommendedAction{},
		RequiredParts:      []RequiredPart{},
		CreatedAt:          time.Now(),
	}

	// Step 1: Enrich context with historical data
	var enrichedContext *EnrichedContext
	if req.Options.IncludeHistoricalContext || req.Options.IncludeSimilarTickets {
		ctx, err := e.contextEnricher.Enrich(ctx, req)
		if err != nil {
			fmt.Printf("Warning: context enrichment failed: %v\n", err)
		} else {
			enrichedContext = ctx
			response.ContextUsed = e.buildContextUsed(enrichedContext)
		}
	}

	// Step 2: Analyze images if present
	if len(req.Attachments) > 0 && req.Options.IncludeVisionAnalysis {
		visionResult, err := e.visionAnalyzer.AnalyzeAttachments(ctx, req)
		if err != nil {
			fmt.Printf("Warning: vision analysis failed: %v\n", err)
		} else if visionResult != nil {
			response.VisionAnalysis = visionResult
		}
	}

	// Step 3: Perform AI diagnosis
	diagnosisResult, metadata, err := e.performAIDiagnosis(ctx, req, enrichedContext, response.VisionAnalysis)
	if err != nil {
		return nil, fmt.Errorf("AI diagnosis failed: %w", err)
	}

	// Parse diagnosis result
	e.parseDiagnosisResult(diagnosisResult, response)

	// Set metadata
	response.Metadata = metadata
	response.Metadata.Latency = time.Since(startTime)
	response.Metadata.Version = e.version
	response.Metadata.VisionAnalysisPerformed = response.VisionAnalysis != nil
	response.Metadata.ContextEnrichmentPerformed = enrichedContext != nil

	// Save diagnosis to database
	if err := e.saveDiagnosis(ctx, response); err != nil {
		fmt.Printf("Warning: failed to save diagnosis: %v\n", err)
	}

	return response, nil
}

// performAIDiagnosis performs the main AI diagnosis
func (e *Engine) performAIDiagnosis(ctx context.Context, req *DiagnosisRequest, enrichedCtx *EnrichedContext, visionResult *VisionAnalysisResult) (string, DiagnosisMetadata, error) {
	// Build comprehensive prompt
	prompt := e.buildDiagnosisPrompt(req, enrichedCtx, visionResult)

	// Build AI request
	aiReq := &ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role:    ai.RoleSystem,
				Content: e.getSystemPrompt(),
			},
			{
				Role:    ai.RoleUser,
				Content: prompt,
			},
		},
	}

	// Use options from request
	if req.Options.Model != "" {
		aiReq.Model = req.Options.Model
	}
	if req.Options.Temperature != nil {
		aiReq.Temperature = req.Options.Temperature
	}
	if req.Options.MaxTokens != nil {
		aiReq.MaxTokens = req.Options.MaxTokens
	}

	// Call AI
	resp, err := e.aiManager.Chat(ctx, aiReq)
	if err != nil {
		return "", DiagnosisMetadata{}, fmt.Errorf("AI chat failed: %w", err)
	}

	// Build metadata
	metadata := DiagnosisMetadata{
		Provider:   resp.Provider,
		Model:      resp.Model,
		TokensUsed: resp.Usage.TotalTokens,
		CostUSD:    resp.Cost,
		Latency:    resp.Latency,
	}

	return resp.Content, metadata, nil
}

// getSystemPrompt returns the system prompt for diagnosis
func (e *Engine) getSystemPrompt() string {
	return `You are an expert medical equipment diagnostic system. Your role is to analyze equipment issues and provide accurate, actionable diagnoses.

**Your Capabilities:**
- Analyze equipment problems across hardware, software, and configuration issues
- Consider historical context and similar past issues
- Interpret visual evidence from images
- Recommend specific repair actions with safety considerations
- Identify required parts with confidence levels
- Estimate resolution time based on issue complexity

**Output Requirements:**
Provide your diagnosis in the following structured format:

## PRIMARY DIAGNOSIS
**Category:** [Hardware/Software/Configuration/User Error/Network/Power/Environmental]
**Problem Type:** [Specific classification]
**Confidence:** [0-100%]
**Severity:** [Low/Medium/High/Critical]

**Description:**
[Detailed description of the problem]

**Root Cause:**
[Identified root cause]

**Symptoms:**
- [Observed symptom 1]
- [Observed symptom 2]

**Reasoning:**
[Explain why you believe this is the issue, referencing evidence]

## ALTERNATE DIAGNOSES
[If applicable, list 1-2 alternative diagnoses with lower confidence]

## RECOMMENDED ACTIONS
1. **[Action Type]** [Description]
   - Estimated Time: [duration]
   - Required Tools: [list]
   - Safety: [precautions]

2. [Continue...]

## REQUIRED PARTS
- **[Part Name]** (Probability: [0-100]%)
  - Part Code: [code]
  - OEM Required: [Yes/No]
  - Quantity: [number]

## ESTIMATED RESOLUTION TIME
[Time estimate with reasoning]

**Guidelines:**
- Be specific and technical
- Reference visual evidence when available
- Consider equipment age and maintenance history
- Prioritize safety
- Be realistic about confidence levels
- Focus on actionable steps`
}

// buildDiagnosisPrompt builds the comprehensive diagnosis prompt
func (e *Engine) buildDiagnosisPrompt(req *DiagnosisRequest, enrichedCtx *EnrichedContext, visionResult *VisionAnalysisResult) string {
	var prompt strings.Builder

	prompt.WriteString("# Equipment Issue Diagnosis Request\n\n")

	// Equipment information
	prompt.WriteString("## Equipment Information\n")
	prompt.WriteString(fmt.Sprintf("- **Type:** %s\n", req.EquipmentType))

	if req.Manufacturer != nil {
		prompt.WriteString(fmt.Sprintf("- **Manufacturer:** %s\n", *req.Manufacturer))
	}
	if req.ModelNumber != nil {
		prompt.WriteString(fmt.Sprintf("- **Model:** %s\n", *req.ModelNumber))
	}

	prompt.WriteString(fmt.Sprintf("- **Location:** %s", req.Location))
	if req.LocationType != nil {
		prompt.WriteString(fmt.Sprintf(" (%s)", *req.LocationType))
	}
	prompt.WriteString("\n")

	prompt.WriteString(fmt.Sprintf("- **Priority:** %s\n", req.Priority))

	// Issue description
	prompt.WriteString("\n## Reported Issue\n")
	prompt.WriteString(fmt.Sprintf("%s\n", req.Description))

	// Reporter information
	prompt.WriteString(fmt.Sprintf("\n**Reported by:** %s (%s)\n", req.ReportedBy.Username, req.ReportedBy.Role))

	// Vision analysis results
	if visionResult != nil && len(visionResult.Findings) > 0 {
		prompt.WriteString("\n## Visual Analysis Results\n")
		prompt.WriteString(fmt.Sprintf("**Overall Assessment:** %s (Confidence: %.0f%%)\n\n", 
			visionResult.OverallAssessment, visionResult.Confidence))

		if len(visionResult.DetectedComponents) > 0 {
			prompt.WriteString("**Visible Components:**\n")
			for _, component := range visionResult.DetectedComponents {
				prompt.WriteString(fmt.Sprintf("- %s\n", component))
			}
			prompt.WriteString("\n")
		}

		if len(visionResult.VisibleDamage) > 0 {
			prompt.WriteString("**Visible Damage:**\n")
			for _, damage := range visionResult.VisibleDamage {
				prompt.WriteString(fmt.Sprintf("- %s: %s (Severity: %s)\n", 
					damage.Type, damage.Description, damage.Severity))
			}
			prompt.WriteString("\n")
		}

		prompt.WriteString("**Detailed Findings:**\n")
		for i, finding := range visionResult.Findings {
			if i < 3 { // Limit to first 3 findings
				prompt.WriteString(fmt.Sprintf("- [%s] %s\n", finding.Category, finding.Finding))
			}
		}
		prompt.WriteString("\n")
	}

	// Enriched context
	if enrichedCtx != nil {
		contextStr := enrichedCtx.FormatContextForAI()
		if contextStr != "" {
			prompt.WriteString(contextStr)
		}
	}

	// Additional context
	if len(req.AdditionalContext) > 0 {
		prompt.WriteString("\n## Additional Context\n")
		for key, value := range req.AdditionalContext {
			prompt.WriteString(fmt.Sprintf("- **%s:** %v\n", key, value))
		}
	}

	prompt.WriteString("\n---\n\n")
	prompt.WriteString("Based on all the information above, provide your comprehensive diagnosis following the required format.\n")

	return prompt.String()
}

// parseDiagnosisResult parses AI response into structured diagnosis
func (e *Engine) parseDiagnosisResult(result string, response *DiagnosisResponse) {
	// This is a simple parser - in production, you might use more sophisticated NLP
	// or have the AI return structured JSON

	lines := strings.Split(result, "\n")

	var currentSection string
	var primaryDiagnosis DiagnosisResult
	var alternateDiagnoses []DiagnosisResult
	var actions []RecommendedAction
	var parts []RequiredPart

	actionOrder := 1

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect sections
		if strings.Contains(line, "## PRIMARY DIAGNOSIS") {
			currentSection = "primary"
			continue
		} else if strings.Contains(line, "## ALTERNATE DIAGNOSES") {
			currentSection = "alternate"
			continue
		} else if strings.Contains(line, "## RECOMMENDED ACTIONS") {
			currentSection = "actions"
			continue
		} else if strings.Contains(line, "## REQUIRED PARTS") {
			currentSection = "parts"
			continue
		} else if strings.Contains(line, "## ESTIMATED RESOLUTION TIME") {
			currentSection = "time"
			continue
		}

		// Parse based on section
		switch currentSection {
		case "primary":
			e.parsePrimaryDiagnosisLine(line, &primaryDiagnosis)

		case "actions":
			if strings.HasPrefix(line, fmt.Sprintf("%d.", actionOrder)) ||
				strings.HasPrefix(line, fmt.Sprintf("%d)", actionOrder)) {
				// New action
				action := e.parseAction(line, actionOrder)
				if action != nil {
					actions = append(actions, *action)
					actionOrder++
				}
			}

		case "parts":
			if strings.HasPrefix(line, "- **") || strings.HasPrefix(line, "* **") {
				part := e.parsePart(line)
				if part != nil {
					parts = append(parts, *part)
				}
			}
		}
	}

	// Set defaults if not parsed
	if primaryDiagnosis.ProblemCategory == "" {
		primaryDiagnosis.ProblemCategory = CategoryUnknown
		primaryDiagnosis.Description = result
		primaryDiagnosis.Confidence = 50.0
		primaryDiagnosis.Severity = SeverityMedium
	}

	response.PrimaryDiagnosis = primaryDiagnosis
	response.AlternateDiagnoses = alternateDiagnoses
	response.RecommendedActions = actions
	response.RequiredParts = parts
}

// parsePrimaryDiagnosisLine parses a line from primary diagnosis
func (e *Engine) parsePrimaryDiagnosisLine(line string, diagnosis *DiagnosisResult) {
	if strings.HasPrefix(line, "**Category:**") {
		diagnosis.ProblemCategory = strings.TrimSpace(strings.TrimPrefix(line, "**Category:**"))
	} else if strings.HasPrefix(line, "**Problem Type:**") {
		diagnosis.ProblemType = strings.TrimSpace(strings.TrimPrefix(line, "**Problem Type:**"))
	} else if strings.HasPrefix(line, "**Confidence:**") {
		confStr := strings.TrimSpace(strings.TrimPrefix(line, "**Confidence:**"))
		confStr = strings.TrimSuffix(confStr, "%")
		var conf float64
		fmt.Sscanf(confStr, "%f", &conf)
		diagnosis.Confidence = conf
	} else if strings.HasPrefix(line, "**Severity:**") {
		diagnosis.Severity = strings.TrimSpace(strings.TrimPrefix(line, "**Severity:**"))
	} else if strings.HasPrefix(line, "**Description:**") {
		diagnosis.Description = strings.TrimSpace(strings.TrimPrefix(line, "**Description:**"))
	} else if strings.HasPrefix(line, "**Root Cause:**") {
		diagnosis.RootCause = strings.TrimSpace(strings.TrimPrefix(line, "**Root Cause:**"))
	} else if strings.HasPrefix(line, "**Reasoning:**") {
		diagnosis.ReasoningExplanation = strings.TrimSpace(strings.TrimPrefix(line, "**Reasoning:**"))
	} else if strings.HasPrefix(line, "- ") {
		// Symptom or possible cause
		symptom := strings.TrimPrefix(line, "- ")
		symptom = strings.TrimPrefix(symptom, "* ")
		diagnosis.Symptoms = append(diagnosis.Symptoms, strings.TrimSpace(symptom))
	}
}

// parseAction parses an action from text
func (e *Engine) parseAction(line string, order int) *RecommendedAction {
	// Simple parsing - extract action text
	parts := strings.SplitN(line, "**", 3)
	if len(parts) < 3 {
		return nil
	}

	actionType := strings.TrimSpace(parts[1])
	description := strings.TrimSpace(parts[2])

	return &RecommendedAction{
		Order:      order,
		Action:     description,
		ActionType: actionType,
	}
}

// parsePart parses a part from text
func (e *Engine) parsePart(line string) *RequiredPart {
	// Extract part name and probability
	// Format: - **[Part Name]** (Probability: [0-100]%)

	if !strings.Contains(line, "**") {
		return nil
	}

	parts := strings.Split(line, "**")
	if len(parts) < 2 {
		return nil
	}

	partName := strings.TrimSpace(parts[1])

	// Extract probability if present
	probability := 70.0 // Default
	if strings.Contains(line, "Probability:") {
		probStr := strings.Split(line, "Probability:")[1]
		probStr = strings.TrimSpace(probStr)
		probStr = strings.TrimSuffix(probStr, "%")
		probStr = strings.TrimSuffix(probStr, ")")
		var prob float64
		fmt.Sscanf(probStr, "%f", &prob)
		if prob > 0 && prob <= 100 {
			probability = prob
		}
	}

	return &RequiredPart{
		PartName:    partName,
		Probability: probability,
		Quantity:    1,
	}
}

// buildContextUsed builds the context used metadata
func (e *Engine) buildContextUsed(enriched *EnrichedContext) DiagnosisContext {
	ctx := DiagnosisContext{}

	if len(enriched.EquipmentHistory) > 0 {
		ctx.EquipmentHistoryUsed = true
		ctx.EquipmentHistoryCount = len(enriched.EquipmentHistory)
	}

	if len(enriched.SimilarTickets) > 0 {
		ctx.SimilarTicketsUsed = true
		ctx.SimilarTicketsCount = len(enriched.SimilarTickets)
		ctx.SimilarTickets = enriched.SimilarTickets
	}

	ctx.ManufacturerGuidelinesUsed = enriched.ManufacturerInfo != nil
	ctx.KnownIssuesUsed = len(enriched.KnownIssues) > 0

	return ctx
}

// saveDiagnosis saves the diagnosis to database
func (e *Engine) saveDiagnosis(ctx context.Context, response *DiagnosisResponse) error {
	// Serialize diagnosis data to JSON
	diagnosisJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal diagnosis: %w", err)
	}

	query := `
		INSERT INTO ai_diagnoses (
			diagnosis_id,
			ticket_id,
			diagnosis_data,
			confidence_score,
			provider,
			model,
			tokens_used,
			cost_usd,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = e.db.ExecContext(ctx, query,
		response.DiagnosisID,
		response.TicketID,
		diagnosisJSON,
		response.PrimaryDiagnosis.Confidence,
		response.Metadata.Provider,
		response.Metadata.Model,
		response.Metadata.TokensUsed,
		response.Metadata.CostUSD,
		response.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save diagnosis: %w", err)
	}

	return nil
}

// GetDiagnosis retrieves a diagnosis by ID
func (e *Engine) GetDiagnosis(ctx context.Context, diagnosisID string) (*DiagnosisResponse, error) {
	query := `
		SELECT diagnosis_data
		FROM ai_diagnoses
		WHERE diagnosis_id = $1
	`

	var diagnosisJSON []byte
	err := e.db.QueryRowContext(ctx, query, diagnosisID).Scan(&diagnosisJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("diagnosis not found")
		}
		return nil, fmt.Errorf("failed to retrieve diagnosis: %w", err)
	}

	var response DiagnosisResponse
	if err := json.Unmarshal(diagnosisJSON, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal diagnosis: %w", err)
	}

	return &response, nil
}

// GetTicketDiagnoses retrieves all diagnoses for a ticket
func (e *Engine) GetTicketDiagnoses(ctx context.Context, ticketID int64) ([]*DiagnosisResponse, error) {
	query := `
		SELECT diagnosis_data
		FROM ai_diagnoses
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := e.db.QueryContext(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query diagnoses: %w", err)
	}
	defer rows.Close()

	var diagnoses []*DiagnosisResponse
	for rows.Next() {
		var diagnosisJSON []byte
		if err := rows.Scan(&diagnosisJSON); err != nil {
			continue
		}

		var response DiagnosisResponse
		if err := json.Unmarshal(diagnosisJSON, &response); err != nil {
			continue
		}

		diagnoses = append(diagnoses, &response)
	}

	return diagnoses, nil
}

// ProvideFeedback records feedback on a diagnosis
func (e *Engine) ProvideFeedback(ctx context.Context, diagnosisID string, wasAccurate bool, accuracyScore int, notes string, actualResolution string) error {
	query := `
		UPDATE ai_diagnoses
		SET 
			was_accurate = $2,
			accuracy_score = $3,
			feedback_notes = $4,
			actual_resolution = $5,
			updated_at = NOW()
		WHERE diagnosis_id = $1
	`

	_, err := e.db.ExecContext(ctx, query, diagnosisID, wasAccurate, accuracyScore, notes, actualResolution)
	if err != nil {
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	return nil
}

