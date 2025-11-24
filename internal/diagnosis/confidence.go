package diagnosis

import (
	"fmt"
	"strings"
)

// calculateConfidence computes the AI confidence score based on multiple factors
func (e *Engine) calculateConfidence(response *DiagnosisResponse) (*ConfidenceCalculation, error) {
	calc := &ConfidenceCalculation{}
	var totalScore float64
	var factorCount int

	// Factor 1: Vision analysis confidence (if available)
	if response.VisionAnalysis != nil && len(response.VisionAnalysis.Findings) > 0 {
		visionScore := 0.85 // Default high confidence for vision analysis
		calc.VisionConfidence = &visionScore
		totalScore += visionScore
		factorCount++
	}

	// Factor 2: Historical match (if similar tickets found)
	if response.ContextUsed.SimilarTickets != nil && len(response.ContextUsed.SimilarTickets) > 0 {
		// Higher confidence with more similar cases
		similarCount := len(response.ContextUsed.SimilarTickets)
		historicalScore := 0.70 + (float64(similarCount) * 0.05) // 0.70-0.95 based on similar cases
		if historicalScore > 0.95 {
			historicalScore = 0.95
		}
		calc.HistoricalMatch = &historicalScore
		totalScore += historicalScore
		factorCount++
	}

	// Factor 3: Symptom clarity (more symptoms = higher confidence)
	symptomCount := len(response.PrimaryDiagnosis.Symptoms)
	if symptomCount >= 2 {
		var symptomScore float64
		if symptomCount >= 5 {
			symptomScore = 0.90
		} else if symptomCount >= 3 {
			symptomScore = 0.80
		} else {
			symptomScore = 0.70
		}
		calc.SymptomClarity = &symptomScore
		totalScore += symptomScore
		factorCount++
	}

	// Factor 4: AI model confidence (based on response characteristics)
	// Check if diagnosis has specific details (not generic)
	if len(response.PrimaryDiagnosis.RootCause) > 50 && 
		!strings.Contains(strings.ToLower(response.PrimaryDiagnosis.Description), "unknown") &&
		!strings.Contains(strings.ToLower(response.PrimaryDiagnosis.Description), "unclear") {
		modelScore := 0.75
		calc.ModelConfidence = &modelScore
		totalScore += modelScore
		factorCount++
	}

	// Calculate final confidence
	if factorCount == 0 {
		calc.FinalConfidence = 0.5 // Default medium confidence
	} else {
		calc.FinalConfidence = totalScore / float64(factorCount)
	}
	
	calc.ConfidenceFactorCount = factorCount

	return calc, nil
}

// getConfidenceLevel converts numeric confidence to level string
func getConfidenceLevel(confidence float64) string {
	if confidence >= 0.80 {
		return "HIGH"
	} else if confidence >= 0.60 {
		return "MEDIUM"
	}
	return "LOW"
}

// getConfidenceFactors generates human-readable confidence factors
func (e *Engine) getConfidenceFactors(response *DiagnosisResponse, calc *ConfidenceCalculation) []string {
	factors := []string{}

	if calc.VisionConfidence != nil {
		factors = append(factors, "Visual analysis of equipment images")
	}

	if calc.HistoricalMatch != nil {
		similarCount := len(response.ContextUsed.SimilarTickets)
		factors = append(factors, fmt.Sprintf("Matched with %d similar historical cases", similarCount))
	}

	if calc.SymptomClarity != nil {
		symptomCount := len(response.PrimaryDiagnosis.Symptoms)
		factors = append(factors, fmt.Sprintf("Analysis of %d reported symptoms", symptomCount))
	}

	if calc.ModelConfidence != nil {
		factors = append(factors, fmt.Sprintf("AI model analysis (%s)", response.Metadata.Model))
	}

	if len(factors) == 0 {
		factors = append(factors, "Basic symptom analysis")
	}

	return factors
}

// EnhanceResponseWithConfidence adds confidence scoring to a DiagnosisResponse (exported for testing)
func (e *Engine) EnhanceResponseWithConfidence(response *DiagnosisResponse) (*EnhancedDiagnosisResponse, error) {
	// Calculate confidence
	calc, err := e.calculateConfidence(response)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate confidence: %w", err)
	}

	// Create enhanced response
	enhanced := &EnhancedDiagnosisResponse{
		DiagnosisID:        response.DiagnosisID,
		TicketID:           response.TicketID,
		PrimaryDiagnosis:   response.PrimaryDiagnosis,
		AlternateDiagnoses: response.AlternateDiagnoses,
		
		// Set confidence data
		Confidence:      calc.FinalConfidence,
		ConfidenceLevel: getConfidenceLevel(calc.FinalConfidence),
		DecisionStatus:  "pending",

		// Set AI metadata
		AIMetadata: AISuggestionMetadata{
			Provider:          response.Metadata.Provider,
			Model:             response.Metadata.Model,
			Confidence:        calc.FinalConfidence,
			ConfidenceFactors: e.getConfidenceFactors(response, calc),
			AlternativesCount: len(response.AlternateDiagnoses),
			RequiresFeedback:  calc.FinalConfidence < 0.80,
			SuggestionOnly:    true,
		},

		// Copy original fields
		VisionAnalysis:          response.VisionAnalysis,
		ContextUsed:            response.ContextUsed,
		RecommendedActions:     response.RecommendedActions,
		RequiredParts:          response.RequiredParts,
		EstimatedResolutionTime: response.EstimatedResolutionTime,
		Metadata:               response.Metadata,
		CreatedAt:              response.CreatedAt,
	}

	return enhanced, nil
}