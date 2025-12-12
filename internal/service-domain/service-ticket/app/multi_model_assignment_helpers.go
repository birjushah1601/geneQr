package app

import (
	"context"
	"fmt"
	"sort"
	"strings"

	equipmentDomain "github.com/aby-med/medical-platform/internal/service-domain/equipment-registry/domain"
	"github.com/aby-med/medical-platform/internal/service-domain/service-ticket/domain"
)

// getEngineerWorkloads calculates workload for all engineers
func (s *MultiModelAssignmentService) getEngineerWorkloads(ctx context.Context, engineers []*domain.Engineer) (map[string]*WorkloadInfo, error) {
	workloadMap := make(map[string]*WorkloadInfo)
	
	// Query workload for each engineer
	query := `
		SELECT 
			assigned_engineer_name,
			COUNT(*) FILTER (WHERE status IN ('new', 'assigned', 'in_progress', 'on_hold')) as active_count,
			COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count
		FROM service_tickets
		WHERE assigned_engineer_name != '' AND assigned_engineer_name IS NOT NULL
		GROUP BY assigned_engineer_name
	`
	
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	workloadByName := make(map[string]*WorkloadInfo)
	for rows.Next() {
		var name string
		var active, inProgress int
		if err := rows.Scan(&name, &active, &inProgress); err != nil {
			continue
		}
		workloadByName[name] = &WorkloadInfo{
			ActiveTickets:     active,
			InProgressTickets: inProgress,
		}
	}
	
	// Map to engineer IDs
	for _, eng := range engineers {
		if wl, found := workloadByName[eng.Name]; found {
			workloadMap[eng.ID] = wl
		} else {
			workloadMap[eng.ID] = &WorkloadInfo{
				ActiveTickets:     0,
				InProgressTickets: 0,
			}
		}
	}
	
	return workloadMap, nil
}

// getEngineerCertifications fetches certifications for engineers
func (s *MultiModelAssignmentService) getEngineerCertifications(ctx context.Context, engineers []*domain.Engineer, manufacturer, category string) (map[string][]*CertificationInfo, error) {
	certMap := make(map[string][]*CertificationInfo)
	
	if len(engineers) == 0 {
		return certMap, nil
	}
	
	// Build engineer IDs for query
	engineerIDs := make([]string, 0, len(engineers))
	for _, eng := range engineers {
		engineerIDs = append(engineerIDs, eng.ID)
	}
	
	query := `
		SELECT 
			engineer_id,
			manufacturer_name,
			equipment_category,
			is_certified,
			COALESCE(certification_number, '') as cert_number,
			COALESCE(certification_expiry::TEXT, '') as expiry
		FROM engineer_equipment_types
		WHERE engineer_id = ANY($1)
	`
	
	rows, err := s.pool.Query(ctx, query, engineerIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var engineerID, mfr, cat, certNum, expiry string
		var isCert bool
		if err := rows.Scan(&engineerID, &mfr, &cat, &isCert, &certNum, &expiry); err != nil {
			continue
		}
		
		cert := &CertificationInfo{
			Manufacturer:        mfr,
			Category:            cat,
			IsCertified:         isCert,
			CertificationNumber: certNum,
			Expiry:              expiry,
		}
		
		certMap[engineerID] = append(certMap[engineerID], cert)
	}
	
	return certMap, nil
}

// getTierInformation builds tier information from engineer organizations
func (s *MultiModelAssignmentService) getTierInformation(ctx context.Context, engineers []*domain.Engineer) []*TierInfo {
	tierMap := make(map[int]map[string]bool) // tier -> orgID -> exists
	tierMap[1] = make(map[string]bool)
	tierMap[2] = make(map[string]bool)
	tierMap[3] = make(map[string]bool)
	
	for _, eng := range engineers {
		if eng.OrganizationID != "" {
			// For now, assume tier 1 (can be enhanced with org metadata)
			tierMap[1][eng.OrganizationID] = true
		}
	}
	
	result := []*TierInfo{
		{
			Tier:            1,
			Name:            "Manufacturer",
			OrganizationIDs: keysToSlice(tierMap[1]),
			AvailableCount:  len(tierMap[1]),
		},
		{
			Tier:            2,
			Name:            "Distributor",
			OrganizationIDs: keysToSlice(tierMap[2]),
			AvailableCount:  len(tierMap[2]),
		},
		{
			Tier:            3,
			Name:            "Dealer",
			OrganizationIDs: keysToSlice(tierMap[3]),
			AvailableCount:  len(tierMap[3]),
		},
	}
	
	return result
}

func keysToSlice(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MODEL 1: Best Match - combines all factors
func (s *MultiModelAssignmentService) getBestMatchModel(
	engineers []*domain.Engineer,
	equipment *equipmentDomain.Equipment,
	ticket *domain.ServiceTicket,
	workloadMap map[string]*WorkloadInfo,
	certMap map[string][]*CertificationInfo,
	minLevel int,
) *AssignmentModel {
	suggestions := make([]*EngineerSuggestion, 0)
	
	for _, eng := range engineers {
		// Get level as int
		level := getLevelInt(eng.EngineerLevel)
		
		// Skip if below minimum level
		if level < minLevel {
			continue
		}
		
		// Build suggestion
		sug := &EngineerSuggestion{
			ID:               eng.ID,
			Name:             eng.Name,
			Email:            eng.Email,
			Phone:            eng.Phone,
			EngineerLevel:    level,
			HomeRegion:       "", // TODO: fetch from engineers table
			OrganizationID:   eng.OrganizationID,
			OrganizationName: eng.OrganizationName,
			MatchReasons:     make([]string, 0),
			Workload:         workloadMap[eng.ID],
			Certifications:   certMap[eng.ID],
		}
		
		// Calculate match score
		score := 0
		
		// Level match (20 points)
		if level >= minLevel {
			score += 20
			sug.MatchReasons = append(sug.MatchReasons, fmt.Sprintf("✅ Level %d engineer (meets requirement)", level))
		}
		
		// Certification match (40 points)
		hasCert := false
		for _, cert := range certMap[eng.ID] {
			if cert.Manufacturer == equipment.ManufacturerName && cert.Category == equipment.Category && cert.IsCertified {
				score += 40
				sug.MatchReasons = append(sug.MatchReasons, fmt.Sprintf("✅ Certified for %s %s equipment", equipment.ManufacturerName, equipment.Category))
				hasCert = true
				break
			}
		}
		
		if !hasCert {
			// Check for category match without certification (10 points)
			for _, cert := range certMap[eng.ID] {
				if cert.Category == equipment.Category {
					score += 10
					sug.MatchReasons = append(sug.MatchReasons, fmt.Sprintf("⚠️ Familiar with %s equipment (not certified)", equipment.Category))
					break
				}
			}
		}
		
		// Workload score (20 points for low workload)
		workload := workloadMap[eng.ID]
		if workload != nil {
			if workload.ActiveTickets == 0 {
				score += 20
				sug.MatchReasons = append(sug.MatchReasons, "✅ No active tickets - fully available")
			} else if workload.ActiveTickets <= 2 {
				score += 10
				sug.MatchReasons = append(sug.MatchReasons, fmt.Sprintf("✅ Low workload (%d active tickets)", workload.ActiveTickets))
			} else {
				sug.MatchReasons = append(sug.MatchReasons, fmt.Sprintf("⚠️ Moderate workload (%d active tickets)", workload.ActiveTickets))
			}
		}
		
		// Organization tier (10 points for manufacturer)
		if strings.Contains(strings.ToLower(eng.OrganizationName), "siemens") || 
		   strings.Contains(strings.ToLower(eng.OrganizationName), "philips") ||
		   strings.Contains(strings.ToLower(eng.OrganizationName), "ge") {
			score += 10
			sug.MatchReasons = append(sug.MatchReasons, "✅ Manufacturer organization")
		}
		
		// Senior engineer bonus (10 points)
		if level == 3 {
			score += 10
			sug.MatchReasons = append(sug.MatchReasons, "✅ Senior engineer (Level 3)")
		}
		
		sug.MatchScore = score
		suggestions = append(suggestions, sug)
	}
	
	// Sort by score descending
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].MatchScore > suggestions[j].MatchScore
	})
	
	// Return top 10
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}
	
	return &AssignmentModel{
		ModelName:   "Best Overall Match",
		Description: "Combines certification, experience level, workload, and organization tier",
		Engineers:   suggestions,
		Count:       len(suggestions),
	}
}

// MODEL 2: Manufacturer Certified
func (s *MultiModelAssignmentService) getManufacturerCertifiedModel(
	engineers []*domain.Engineer,
	equipment *equipmentDomain.Equipment,
	certMap map[string][]*CertificationInfo,
	minLevel int,
) *AssignmentModel {
	suggestions := make([]*EngineerSuggestion, 0)
	
	for _, eng := range engineers {
		level := getLevelInt(eng.EngineerLevel)
		if level < minLevel {
			continue
		}
		
		// Check for valid certification
		hasCert := false
		for _, cert := range certMap[eng.ID] {
			if cert.Manufacturer == equipment.ManufacturerName && 
			   cert.Category == equipment.Category && 
			   cert.IsCertified {
				hasCert = true
				break
			}
		}
		
		if !hasCert {
			continue
		}
		
		sug := &EngineerSuggestion{
			ID:               eng.ID,
			Name:             eng.Name,
			Email:            eng.Email,
			Phone:            eng.Phone,
			EngineerLevel:    level,
			OrganizationID:   eng.OrganizationID,
			OrganizationName: eng.OrganizationName,
			Certifications:   certMap[eng.ID],
			MatchScore:       100,
			MatchReasons: []string{
				fmt.Sprintf("✅ Certified by %s for %s equipment", equipment.ManufacturerName, equipment.Category),
				fmt.Sprintf("✅ Level %d engineer", level),
			},
		}
		
		suggestions = append(suggestions, sug)
	}
	
	// Sort by level descending
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].EngineerLevel > suggestions[j].EngineerLevel
	})
	
	return &AssignmentModel{
		ModelName:   "Manufacturer Certified",
		Description: "Engineers officially certified by the equipment manufacturer",
		Engineers:   suggestions,
		Count:       len(suggestions),
	}
}

// MODEL 3: Skills Match
func (s *MultiModelAssignmentService) getSkillsMatchModel(
	engineers []*domain.Engineer,
	equipment *equipmentDomain.Equipment,
	minLevel int,
) *AssignmentModel {
	suggestions := make([]*EngineerSuggestion, 0)
	
	// This would require fetching skills from engineers table
	// For now, return empty as skills are stored as array in DB
	
	return &AssignmentModel{
		ModelName:   "Skills Match",
		Description: "Engineers with matching equipment category skills",
		Engineers:   suggestions,
		Count:       0,
	}
}

// MODEL 4: Low Workload
func (s *MultiModelAssignmentService) getLowWorkloadModel(
	engineers []*domain.Engineer,
	workloadMap map[string]*WorkloadInfo,
	minLevel int,
) *AssignmentModel {
	suggestions := make([]*EngineerSuggestion, 0)
	
	for _, eng := range engineers {
		level := getLevelInt(eng.EngineerLevel)
		if level < minLevel {
			continue
		}
		
		workload := workloadMap[eng.ID]
		if workload == nil {
			workload = &WorkloadInfo{ActiveTickets: 0, InProgressTickets: 0}
		}
		
		sug := &EngineerSuggestion{
			ID:               eng.ID,
			Name:             eng.Name,
			Email:            eng.Email,
			Phone:            eng.Phone,
			EngineerLevel:    level,
			OrganizationID:   eng.OrganizationID,
			OrganizationName: eng.OrganizationName,
			Workload:         workload,
			MatchScore:       100 - (workload.ActiveTickets * 10),
			MatchReasons: []string{
				fmt.Sprintf("✅ %d active tickets", workload.ActiveTickets),
				fmt.Sprintf("✅ Level %d engineer", level),
			},
		}
		
		suggestions = append(suggestions, sug)
	}
	
	// Sort by workload ascending
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Workload.ActiveTickets < suggestions[j].Workload.ActiveTickets
	})
	
	// Return top 10
	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}
	
	return &AssignmentModel{
		ModelName:   "Low Workload",
		Description: "Engineers with fewest active tickets - most available",
		Engineers:   suggestions,
		Count:       len(suggestions),
	}
}

// MODEL 5: High Seniority
func (s *MultiModelAssignmentService) getHighSeniorityModel(
	engineers []*domain.Engineer,
	equipment *equipmentDomain.Equipment,
) *AssignmentModel {
	suggestions := make([]*EngineerSuggestion, 0)
	
	for _, eng := range engineers {
		level := getLevelInt(eng.EngineerLevel)
		
		// Only Level 3 engineers
		if level != 3 {
			continue
		}
		
		sug := &EngineerSuggestion{
			ID:               eng.ID,
			Name:             eng.Name,
			Email:            eng.Email,
			Phone:            eng.Phone,
			EngineerLevel:    level,
			OrganizationID:   eng.OrganizationID,
			OrganizationName: eng.OrganizationName,
			MatchScore:       100,
			MatchReasons: []string{
				"✅ Senior engineer (Level 3)",
				"✅ Highest experience level available",
			},
		}
		
		suggestions = append(suggestions, sug)
	}
	
	return &AssignmentModel{
		ModelName:   "Senior Engineers Only",
		Description: "Level 3 (Senior) engineers with highest experience",
		Engineers:   suggestions,
		Count:       len(suggestions),
	}
}

// Helper: convert EngineerLevel to int
func getLevelInt(level domain.EngineerLevel) int {
	switch level {
	case domain.EngineerLevelL1:
		return 1
	case domain.EngineerLevelL2:
		return 2
	case domain.EngineerLevelL3:
		return 3
	default:
		return 1
	}
}
