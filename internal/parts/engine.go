package parts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aby-med/medical-platform/internal/ai"
	"github.com/google/uuid"
)

const partsEngineVersion = "1.0.0"

// Engine handles intelligent parts recommendations
type Engine struct {
	aiManager *ai.Manager
	db        *sql.DB
}

// NewEngine creates a new parts recommendation engine
func NewEngine(aiManager *ai.Manager, db *sql.DB) *Engine {
	return &Engine{
		aiManager: aiManager,
		db:        db,
	}
}

// RecommendParts returns intelligent parts recommendations
func (e *Engine) RecommendParts(ctx context.Context, req *RecommendationRequest) (*RecommendationResponse, error) {
	startTime := time.Now()

	// Generate request ID
	requestID := uuid.New().String()

	response := &RecommendationResponse{
		RequestID: requestID,
		TicketID:  req.TicketID,
		Metadata: RecommendationMetadata{
			Version:  partsEngineVersion,
			Currency: "USD",
		},
		CreatedAt: time.Now(),
	}

	// Step 1: Get replacement parts (if diagnosis/problem indicates)
	if req.Options.IncludeReplacementParts {
		parts, err := e.getReplacementParts(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get replacement parts: %w", err)
		}
		response.ReplacementParts = parts
		response.Metadata.TotalReplacementParts = len(parts)
	}

	// Step 2: Get preventive maintenance parts
	if req.Options.IncludePreventiveParts {
		parts, err := e.getPreventiveParts(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get preventive parts: %w", err)
		}
		response.PreventiveParts = parts
		response.Metadata.TotalPreventiveParts = len(parts)
	}

	// Step 3: Get accessories for upselling
	if req.Options.IncludeAccessories {
		accessories, err := e.getAccessories(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to get accessories: %w", err)
		}
		response.Accessories = accessories
		response.Metadata.TotalAccessories = len(accessories)
	}

	// Step 4: Use AI to refine recommendations (if enabled)
	if req.Options.UseAI && len(response.ReplacementParts) > 0 {
		refined, cost, provider, model, err := e.aiRefineRecommendations(ctx, req, response)
		if err == nil && refined != nil {
			response.ReplacementParts = refined
			response.Metadata.UsedAI = true
			response.Metadata.CostUSD = cost
			response.Metadata.AIProvider = provider
			response.Metadata.AIModel = model
		}
	}

	// Step 5: Assign ranks
	e.assignRanks(response)

	// Step 6: Calculate costs
	e.calculateCosts(response)

	// Step 7: Save to database
	response.Metadata.ProcessingTime = time.Since(startTime)
	if err := e.saveRecommendation(ctx, response); err != nil {
		// Log but don't fail
		fmt.Printf("Warning: failed to save recommendation: %v\n", err)
	}

	return response, nil
}

// getReplacementParts gets parts likely needed for repair
func (e *Engine) getReplacementParts(ctx context.Context, req *RecommendationRequest) ([]PartRecommendation, error) {
	var parts []PartRecommendation

	// Query 1: Get parts matching diagnosis/problem
	diagnosisParts, err := e.getPartsFromDiagnosis(ctx, req)
	if err == nil {
		parts = append(parts, diagnosisParts...)
	}

	// Query 2: Get parts from historical patterns
	historicalParts, err := e.getPartsFromHistory(ctx, req)
	if err == nil {
		parts = append(parts, historicalParts...)
	}

	// Query 3: Get critical parts
	criticalParts, err := e.getCriticalParts(ctx, req)
	if err == nil {
		parts = append(parts, criticalParts...)
	}

	// Deduplicate and merge
	parts = e.deduplicateParts(parts)

	// Enrich with inventory and pricing
	if req.Options.CheckInventory || req.Options.IncludePricing {
		for i := range parts {
			e.enrichPartWithInventoryAndPricing(ctx, &parts[i], req.Options)
		}
	}

	// Filter by confidence threshold
	filtered := []PartRecommendation{}
	for _, part := range parts {
		if part.Confidence >= req.Options.MinConfidence {
			filtered = append(filtered, part)
		}
	}

	// Sort by confidence
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Confidence > filtered[j].Confidence
	})

	// Limit results
	if len(filtered) > req.Options.MaxRecommendations {
		filtered = filtered[:req.Options.MaxRecommendations]
	}

	return filtered, nil
}

// getPartsFromDiagnosis gets parts based on diagnosis
func (e *Engine) getPartsFromDiagnosis(ctx context.Context, req *RecommendationRequest) ([]PartRecommendation, error) {
	var parts []PartRecommendation

	// Build search terms from problem description and identified issues
	searchTerms := e.extractSearchTerms(req)

	if len(searchTerms) == 0 {
		return parts, nil
	}

	// Query parts that match the problem
	query := `
		SELECT 
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.subcategory,
			pc.is_oem_part,
			ep.is_critical_part,
			ep.installation_notes,
			ep.compatibility_notes
		FROM parts_catalog pc
		JOIN equipment_parts ep ON pc.part_id = ep.part_id
		WHERE ep.equipment_type_id = $1
			AND pc.is_active = true
			AND ep.is_active = true
			AND (
				ep.variant_id = $2 OR ep.variant_id IS NULL
			)
			AND (
				pc.part_name ILIKE ANY($3)
				OR pc.description ILIKE ANY($3)
				OR pc.category ILIKE ANY($3)
			)
		ORDER BY ep.is_critical_part DESC, pc.part_name
	`

	variantID := sql.NullInt64{}
	if req.VariantID != nil {
		variantID.Valid = true
		variantID.Int64 = *req.VariantID
	}

	// Convert search terms to ILIKE patterns
	searchPatterns := make([]string, len(searchTerms))
	for i, term := range searchTerms {
		searchPatterns[i] = "%" + term + "%"
	}

	rows, err := e.db.QueryContext(ctx, query, req.EquipmentTypeID, variantID, searchPatterns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var part PartRecommendation
		var subcategory sql.NullString
		var installNotes, compatNotes sql.NullString

		err := rows.Scan(
			&part.PartID,
			&part.PartNumber,
			&part.PartName,
			&part.Description,
			&part.Category,
			&subcategory,
			&part.IsOEMPart,
			&part.IsCriticalPart,
			&installNotes,
			&compatNotes,
		)
		if err != nil {
			continue
		}

		if subcategory.Valid {
			part.Subcategory = &subcategory.String
		}
		if installNotes.Valid {
			part.InstallationNotes = &installNotes.String
		}
		if compatNotes.Valid {
			part.CompatibilityNotes = &compatNotes.String
		}

		// Calculate confidence based on diagnosis
		part.Confidence = e.calculateDiagnosisConfidence(req, part)
		part.ReasonCode = ReasonDiagnosisMatch
		part.ReasonText = fmt.Sprintf("Matches problem: %s", req.ProblemType)
		part.Evidence = []string{fmt.Sprintf("Related to %s issues", req.ProblemType)}
		
		if req.DiagnosisConfidence != nil {
			part.Evidence = append(part.Evidence, 
				fmt.Sprintf("AI diagnosis confidence: %.0f%%", *req.DiagnosisConfidence))
		}

		part.RecommendedQuantity = 1
		part.QuantityReasoning = "Standard replacement quantity"
		part.Currency = "USD"

		parts = append(parts, part)
	}

	return parts, nil
}

// getPartsFromHistory gets parts based on historical ticket patterns
func (e *Engine) getPartsFromHistory(ctx context.Context, req *RecommendationRequest) ([]PartRecommendation, error) {
	var parts []PartRecommendation

	// Query: Find parts frequently used for similar problems
	query := `
		SELECT 
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.is_oem_part,
			ep.is_critical_part,
			COUNT(*) as usage_count,
			AVG(CASE WHEN st.status = 'Resolved' THEN 1.0 ELSE 0.0 END) as success_rate
		FROM parts_catalog pc
		JOIN equipment_parts ep ON pc.part_id = ep.part_id
		JOIN ticket_parts tp ON pc.part_id = tp.part_id
		JOIN service_tickets st ON tp.ticket_id = st.ticket_id
		WHERE st.equipment_type_id = $1
			AND st.problem_type ILIKE $2
			AND st.created_at > NOW() - INTERVAL '2 years'
			AND pc.is_active = true
		GROUP BY pc.part_id, pc.part_number, pc.part_name, pc.description, 
			pc.category, pc.is_oem_part, ep.is_critical_part
		HAVING COUNT(*) >= 2
		ORDER BY usage_count DESC, success_rate DESC
		LIMIT 5
	`

	problemPattern := "%" + req.ProblemType + "%"
	rows, err := e.db.QueryContext(ctx, query, req.EquipmentTypeID, problemPattern)
	if err != nil {
		return parts, nil // Not critical, return empty
	}
	defer rows.Close()

	for rows.Next() {
		var part PartRecommendation
		var usageCount int
		var successRate float64

		err := rows.Scan(
			&part.PartID,
			&part.PartNumber,
			&part.PartName,
			&part.Description,
			&part.Category,
			&part.IsOEMPart,
			&part.IsCriticalPart,
			&usageCount,
			&successRate,
		)
		if err != nil {
			continue
		}

		// Calculate confidence based on historical success
		part.Confidence = successRate * 80.0 // Scale to 0-80
		if usageCount > 5 {
			part.Confidence += 10.0 // Bonus for frequency
		}
		part.Confidence = min(part.Confidence, 95.0)

		part.ReasonCode = ReasonHistoricalPattern
		part.ReasonText = fmt.Sprintf("Used in %d similar cases with %.0f%% success", 
			usageCount, successRate*100)
		part.Evidence = []string{
			fmt.Sprintf("Used %d times for similar problems", usageCount),
			fmt.Sprintf("Success rate: %.0f%%", successRate*100),
		}

		part.RecommendedQuantity = 1
		part.QuantityReasoning = "Based on historical usage"
		part.Currency = "USD"

		parts = append(parts, part)
	}

	return parts, nil
}

// getCriticalParts gets critical parts for this equipment
func (e *Engine) getCriticalParts(ctx context.Context, req *RecommendationRequest) ([]PartRecommendation, error) {
	var parts []PartRecommendation

	// Only include critical parts if severity is high
	if req.Severity != "High" && req.Severity != "Critical" {
		return parts, nil
	}

	query := `
		SELECT 
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.is_oem_part,
			ep.is_critical_part
		FROM parts_catalog pc
		JOIN equipment_parts ep ON pc.part_id = ep.part_id
		WHERE ep.equipment_type_id = $1
			AND ep.is_critical_part = true
			AND pc.is_active = true
			AND ep.is_active = true
			AND (ep.variant_id = $2 OR ep.variant_id IS NULL)
		ORDER BY pc.part_name
		LIMIT 3
	`

	variantID := sql.NullInt64{}
	if req.VariantID != nil {
		variantID.Valid = true
		variantID.Int64 = *req.VariantID
	}

	rows, err := e.db.QueryContext(ctx, query, req.EquipmentTypeID, variantID)
	if err != nil {
		return parts, nil
	}
	defer rows.Close()

	for rows.Next() {
		var part PartRecommendation

		err := rows.Scan(
			&part.PartID,
			&part.PartNumber,
			&part.PartName,
			&part.Description,
			&part.Category,
			&part.IsOEMPart,
			&part.IsCriticalPart,
		)
		if err != nil {
			continue
		}

		part.Confidence = 60.0 // Moderate confidence
		part.ReasonCode = ReasonCriticalPart
		part.ReasonText = "Critical part for equipment operation"
		part.Evidence = []string{"Identified as critical component"}
		part.RecommendedQuantity = 1
		part.QuantityReasoning = "Keep as spare for critical equipment"
		part.Currency = "USD"

		parts = append(parts, part)
	}

	return parts, nil
}

// getPreventiveParts gets parts due for scheduled replacement
func (e *Engine) getPreventiveParts(ctx context.Context, req *RecommendationRequest) ([]PartRecommendation, error) {
	var parts []PartRecommendation

	if req.LastMaintenanceDate == nil {
		return parts, nil
	}

	query := `
		SELECT 
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.is_oem_part,
			ep.recommended_replacement_interval,
			ep.replacement_interval_hours
		FROM parts_catalog pc
		JOIN equipment_parts ep ON pc.part_id = ep.part_id
		WHERE ep.equipment_type_id = $1
			AND pc.is_active = true
			AND ep.is_active = true
			AND (ep.variant_id = $2 OR ep.variant_id IS NULL)
			AND (
				ep.recommended_replacement_interval IS NOT NULL
				OR ep.replacement_interval_hours IS NOT NULL
			)
	`

	variantID := sql.NullInt64{}
	if req.VariantID != nil {
		variantID.Valid = true
		variantID.Int64 = *req.VariantID
	}

	rows, err := e.db.QueryContext(ctx, query, req.EquipmentTypeID, variantID)
	if err != nil {
		return parts, nil
	}
	defer rows.Close()

	for rows.Next() {
		var part PartRecommendation
		var interval sql.NullString
		var hours sql.NullInt64

		err := rows.Scan(
			&part.PartID,
			&part.PartNumber,
			&part.PartName,
			&part.Description,
			&part.Category,
			&part.IsOEMPart,
			&interval,
			&hours,
		)
		if err != nil {
			continue
		}

		isDue := false
		reason := ""

		// Check time-based interval
		if interval.Valid {
			// Parse interval (simplified - in production use proper interval parsing)
			if strings.Contains(interval.String, "6 months") {
				if time.Since(*req.LastMaintenanceDate) > 6*30*24*time.Hour {
					isDue = true
					reason = "Due for 6-month replacement"
				}
			} else if strings.Contains(interval.String, "1 year") {
				if time.Since(*req.LastMaintenanceDate) > 365*24*time.Hour {
					isDue = true
					reason = "Due for annual replacement"
				}
			}
		}

		// Check hours-based interval
		if hours.Valid && req.OperatingHours != nil {
			if *req.OperatingHours >= int(hours.Int64) {
				isDue = true
				reason = fmt.Sprintf("Operating hours exceeded (%d hrs)", hours.Int64)
			}
		}

		if isDue {
			part.Confidence = 85.0
			part.ReasonCode = ReasonPreventiveMaintenance
			part.ReasonText = reason
			part.Evidence = []string{"Scheduled preventive maintenance"}
			part.RecommendedQuantity = 1
			part.QuantityReasoning = "Preventive replacement"
			part.Currency = "USD"

			parts = append(parts, part)
		}
	}

	return parts, nil
}

// getAccessories gets accessories for upselling based on variant
func (e *Engine) getAccessories(ctx context.Context, req *RecommendationRequest) ([]AccessoryRecommendation, error) {
	var accessories []AccessoryRecommendation

	query := `
		SELECT 
			ea.accessory_id,
			pc.part_id,
			pc.part_number,
			pc.part_name,
			pc.description,
			pc.category,
			pc.unit_price,
			ea.upsell_priority,
			ea.is_recommended,
			ea.is_required_for_variant,
			ea.bundle_discount_percent,
			ea.marketing_description,
			ea.benefits,
			pc.image_url,
			pi.quantity_available
		FROM equipment_accessories ea
		JOIN parts_catalog pc ON ea.part_id = pc.part_id
		LEFT JOIN parts_inventory pi ON pc.part_id = pi.part_id
		WHERE ea.equipment_type_id = $1
			AND ea.is_active = true
			AND pc.is_active = true
			AND (ea.variant_id = $2 OR ea.variant_id IS NULL)
		ORDER BY ea.upsell_priority DESC, ea.is_required_for_variant DESC
	`

	variantID := sql.NullInt64{}
	if req.VariantID != nil {
		variantID.Valid = true
		variantID.Int64 = *req.VariantID
	}

	rows, err := e.db.QueryContext(ctx, query, req.EquipmentTypeID, variantID)
	if err != nil {
		return accessories, nil
	}
	defer rows.Close()

	for rows.Next() {
		var acc AccessoryRecommendation
		var accessoryID int64
		var bundleDiscount sql.NullFloat64
		var marketingDesc sql.NullString
		var imageURL sql.NullString
		var qtyAvailable sql.NullInt64

		// Use pq.Array for PostgreSQL array type
		var benefits []string

		err := rows.Scan(
			&accessoryID,
			&acc.PartID,
			&acc.PartNumber,
			&acc.PartName,
			&acc.Description,
			&acc.Category,
			&acc.UnitPrice,
			&acc.UpsellPriority,
			&acc.IsRecommended,
			&acc.IsRequiredForVariant,
			&bundleDiscount,
			&marketingDesc,
			&benefits,
			&imageURL,
			&qtyAvailable,
		)
		if err != nil {
			continue
		}

		acc.Benefits = benefits
		acc.Currency = "USD"

		if bundleDiscount.Valid {
			acc.BundleDiscountPercent = &bundleDiscount.Float64
			discounted := acc.UnitPrice * (1.0 - bundleDiscount.Float64/100.0)
			acc.DiscountedPrice = &discounted
		}

		if marketingDesc.Valid {
			acc.MarketingDescription = &marketingDesc.String
		}

		if imageURL.Valid {
			acc.ImageURL = &imageURL.String
		}

		// Stock status
		if qtyAvailable.Valid {
			acc.QuantityAvailable = int(qtyAvailable.Int64)
			if qtyAvailable.Int64 == 0 {
				acc.StockStatus = StockOutOfStock
			} else if qtyAvailable.Int64 <= 5 {
				acc.StockStatus = StockLowStock
			} else {
				acc.StockStatus = StockInStock
			}
		} else {
			acc.StockStatus = StockUnknown
		}

		// Build reason text
		if acc.IsRequiredForVariant {
			acc.ReasonText = fmt.Sprintf("Required for %s installation", *req.VariantName)
		} else if acc.IsRecommended {
			acc.ReasonText = "Recommended accessory for this equipment"
		} else {
			acc.ReasonText = "Optional accessory"
		}

		accessories = append(accessories, acc)
	}

	// Limit to reasonable number
	if len(accessories) > 10 {
		accessories = accessories[:10]
	}

	return accessories, nil
}

// Helper functions

func (e *Engine) extractSearchTerms(req *RecommendationRequest) []string {
	terms := []string{}

	// From problem type
	if req.ProblemType != "" {
		terms = append(terms, strings.ToLower(req.ProblemType))
	}

	// From identified issues
	for _, issue := range req.IdentifiedIssues {
		terms = append(terms, strings.ToLower(issue))
	}

	// Extract keywords from description
	keywords := []string{"filter", "valve", "sensor", "pump", "motor", "circuit", "battery", 
		"display", "cable", "switch", "fan", "belt", "tube", "hose"}
	
	desc := strings.ToLower(req.ProblemDescription)
	for _, keyword := range keywords {
		if strings.Contains(desc, keyword) {
			terms = append(terms, keyword)
		}
	}

	return terms
}

func (e *Engine) calculateDiagnosisConfidence(req *RecommendationRequest, part PartRecommendation) float64 {
	confidence := 50.0 // Base confidence

	// Boost if critical part
	if part.IsCriticalPart {
		confidence += 15.0
	}

	// Boost from diagnosis confidence
	if req.DiagnosisConfidence != nil {
		confidence += (*req.DiagnosisConfidence * 0.3) // Up to 30 points
	}

	// Boost for severity match
	if req.Severity == "High" || req.Severity == "Critical" {
		confidence += 10.0
	}

	return min(confidence, 95.0)
}

func (e *Engine) deduplicateParts(parts []PartRecommendation) []PartRecommendation {
	seen := make(map[int64]*PartRecommendation)

	for _, part := range parts {
		if existing, found := seen[part.PartID]; found {
			// Merge: keep higher confidence and combine evidence
			if part.Confidence > existing.Confidence {
				existing.Confidence = part.Confidence
				existing.ReasonCode = part.ReasonCode
				existing.ReasonText = part.ReasonText
			}
			existing.Evidence = append(existing.Evidence, part.Evidence...)
		} else {
			partCopy := part
			seen[part.PartID] = &partCopy
		}
	}

	// Convert back to slice
	result := []PartRecommendation{}
	for _, part := range seen {
		result = append(result, *part)
	}

	return result
}

func (e *Engine) enrichPartWithInventoryAndPricing(ctx context.Context, part *PartRecommendation, opts RecommendationOptions) {
	// Get inventory
	if opts.CheckInventory {
		var qtyAvailable sql.NullInt64
		err := e.db.QueryRowContext(ctx, `
			SELECT quantity_available 
			FROM parts_inventory 
			WHERE part_id = $1 AND is_active = true
			LIMIT 1
		`, part.PartID).Scan(&qtyAvailable)

		if err == nil && qtyAvailable.Valid {
			part.QuantityAvailable = int(qtyAvailable.Int64)
			if qtyAvailable.Int64 == 0 {
				part.StockStatus = StockOutOfStock
			} else if qtyAvailable.Int64 <= 5 {
				part.StockStatus = StockLowStock
			} else {
				part.StockStatus = StockInStock
			}
		} else {
			part.StockStatus = StockUnknown
		}
	}

	// Get pricing from preferred supplier
	if opts.IncludePricing {
		var price sql.NullFloat64
		var leadTime sql.NullInt64
		var supplierID sql.NullInt64
		var supplierName sql.NullString
		var isOEM sql.NullBool
		var isPreferred sql.NullBool
		var inStock sql.NullBool

		err := e.db.QueryRowContext(ctx, `
			SELECT 
				sp.unit_price,
				sp.lead_time_days,
				ps.supplier_id,
				ps.supplier_name,
				ps.is_oem_supplier,
				sp.is_preferred,
				sp.is_in_stock
			FROM supplier_parts sp
			JOIN parts_suppliers ps ON sp.supplier_id = ps.supplier_id
			WHERE sp.part_id = $1 
				AND sp.is_active = true 
				AND ps.is_active = true
			ORDER BY sp.is_preferred DESC, sp.unit_price ASC
			LIMIT 1
		`, part.PartID).Scan(&price, &leadTime, &supplierID, &supplierName, &isOEM, &isPreferred, &inStock)

		if err == nil && price.Valid {
			part.UnitPrice = &price.Float64
			totalPrice := price.Float64 * float64(part.RecommendedQuantity)
			part.TotalPrice = &totalPrice

			if leadTime.Valid {
				days := int(leadTime.Int64)
				part.LeadTimeDays = &days
			}

			if supplierID.Valid {
				part.SupplierInfo = &SupplierInfo{
					SupplierID:   supplierID.Int64,
					SupplierName: supplierName.String,
					IsOEMSupplier: isOEM.Bool,
					IsPreferred:  isPreferred.Bool,
					LeadTimeDays: int(leadTime.Int64),
					InStock:      inStock.Bool,
				}
			}
		}
	}
}

func (e *Engine) assignRanks(response *RecommendationResponse) {
	for i := range response.ReplacementParts {
		response.ReplacementParts[i].Rank = i + 1
	}
	for i := range response.PreventiveParts {
		response.PreventiveParts[i].Rank = i + 1
	}
	for i := range response.Accessories {
		response.Accessories[i].Rank = i + 1
	}
}

func (e *Engine) calculateCosts(response *RecommendationResponse) {
	var partsCost, accessoriesCost float64

	for _, part := range response.ReplacementParts {
		if part.TotalPrice != nil {
			partsCost += *part.TotalPrice
		}
	}

	for _, part := range response.PreventiveParts {
		if part.TotalPrice != nil {
			partsCost += *part.TotalPrice
		}
	}

	for _, acc := range response.Accessories {
		if acc.DiscountedPrice != nil {
			accessoriesCost += *acc.DiscountedPrice
		} else {
			accessoriesCost += acc.UnitPrice
		}
	}

	if partsCost > 0 {
		response.Metadata.EstimatedPartsCost = &partsCost
	}
	if accessoriesCost > 0 {
		response.Metadata.EstimatedAccessoriesCost = &accessoriesCost
	}

	totalCost := partsCost + accessoriesCost
	if totalCost > 0 {
		response.Metadata.EstimatedTotalCost = &totalCost
	}
}

// AI refinement (continued in next file due to length)
func (e *Engine) aiRefineRecommendations(ctx context.Context, req *RecommendationRequest, response *RecommendationResponse) ([]PartRecommendation, float64, string, string, error) {
	prompt := e.buildRefinementPrompt(req, response.ReplacementParts)

	result, err := e.aiManager.Chat(ctx, &ai.ChatRequest{
		Messages: []ai.Message{
			{
				Role:    "system",
				Content: "You are an expert in medical equipment parts and repair. Review parts recommendations and suggest improvements.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: func() *float32 { v := float32(0.3); return &v }(),
		MaxTokens:   func() *int { v := 800; return &v }(),
	})

	if err != nil {
		return nil, 0, "", "", err
	}

	// Parse and apply refinements
	refined := e.parseAIRefinements(result.Content, response.ReplacementParts)

	return refined, result.Cost, result.Provider, result.Model, nil
}

func (e *Engine) buildRefinementPrompt(req *RecommendationRequest, parts []PartRecommendation) string {
	prompt := fmt.Sprintf(`Review parts recommendations for equipment repair:

**Equipment:** %s
**Problem:** %s
**Severity:** %s

**Recommended Parts:**
`, req.EquipmentType, req.ProblemType, req.Severity)

	for i, part := range parts {
		if i >= 5 {
			break
		}
		prompt += fmt.Sprintf(`
%d. %s (%s)
   - Category: %s
   - Confidence: %.0f%%
   - Reason: %s
`, i+1, part.PartName, part.PartNumber, part.Category, part.Confidence, part.ReasonText)
	}

	prompt += `
**Task:** Review and respond in JSON:
{
  "refinements": [
    {"part_number": "...", "confidence_adjustment": +5 or -10, "reason": "..."}
  ],
  "concerns": ["..."],
  "additional_suggestions": ["..."]
}
`
	return prompt
}

func (e *Engine) parseAIRefinements(aiResponse string, parts []PartRecommendation) []PartRecommendation {
	var response struct {
		Refinements []struct {
			PartNumber           string  `json:"part_number"`
			ConfidenceAdjustment float64 `json:"confidence_adjustment"`
			Reason               string  `json:"reason"`
		} `json:"refinements"`
	}

	if err := json.Unmarshal([]byte(aiResponse), &response); err != nil {
		return parts // Return original on parse error
	}

	// Apply refinements
	refined := make([]PartRecommendation, len(parts))
	copy(refined, parts)

	for i := range refined {
		for _, ref := range response.Refinements {
			if refined[i].PartNumber == ref.PartNumber {
				refined[i].Confidence += ref.ConfidenceAdjustment
				refined[i].Confidence = max(0, min(100, refined[i].Confidence))
				refined[i].Evidence = append(refined[i].Evidence, 
					fmt.Sprintf("AI: %s", ref.Reason))
			}
		}
	}

	// Re-sort by confidence
	sort.Slice(refined, func(i, j int) bool {
		return refined[i].Confidence > refined[j].Confidence
	})

	return refined
}

func (e *Engine) saveRecommendation(ctx context.Context, response *RecommendationResponse) error {
	partsJSON, _ := json.Marshal(response.ReplacementParts)
	accessoriesJSON, _ := json.Marshal(response.Accessories)
	preventiveJSON, _ := json.Marshal(response.PreventiveParts)
	metadataJSON, _ := json.Marshal(response.Metadata)

	query := `
		INSERT INTO parts_recommendations
		(request_id, ticket_id, replacement_parts, accessories, preventive_parts, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := e.db.ExecContext(ctx, query,
		response.RequestID,
		response.TicketID,
		partsJSON,
		accessoriesJSON,
		preventiveJSON,
		metadataJSON,
		response.CreatedAt,
	)

	return err
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}



