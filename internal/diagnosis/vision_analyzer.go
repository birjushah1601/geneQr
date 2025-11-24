package diagnosis

import (
	"context"
	"fmt"
	"strings"

	"github.com/aby-med/medical-platform/internal/ai"
)

// VisionAnalyzer analyzes images using AI vision capabilities
type VisionAnalyzer struct {
	aiManager *ai.Manager
}

// NewVisionAnalyzer creates a new vision analyzer
func NewVisionAnalyzer(aiManager *ai.Manager) *VisionAnalyzer {
	return &VisionAnalyzer{
		aiManager: aiManager,
	}
}

// AnalyzeAttachments analyzes ticket attachments for visual issues
func (va *VisionAnalyzer) AnalyzeAttachments(ctx context.Context, req *DiagnosisRequest) (*VisionAnalysisResult, error) {
	if len(req.Attachments) == 0 {
		return nil, nil
	}

	// Filter image attachments
	imageAttachments := filterImageAttachments(req.Attachments)
	if len(imageAttachments) == 0 {
		return nil, nil
	}

	result := &VisionAnalysisResult{
		AttachmentsAnalyzed: len(imageAttachments),
		Findings:            []VisualFinding{},
		DetectedComponents:  []string{},
		VisibleDamage:       []DamageDescription{},
	}

	// Analyze each image
	for _, attachment := range imageAttachments {
		finding, err := va.analyzeImage(ctx, req, attachment)
		if err != nil {
			fmt.Printf("Warning: failed to analyze attachment %d: %v\n", attachment.ID, err)
			continue
		}

		if finding != nil {
			result.Findings = append(result.Findings, *finding)
		}
	}

	// If we got findings, generate overall assessment
	if len(result.Findings) > 0 {
		overall, confidence := va.generateOverallAssessment(result.Findings)
		result.OverallAssessment = overall
		result.Confidence = confidence

		// Extract components and damage from findings
		result.DetectedComponents = va.extractComponents(result.Findings)
		result.VisibleDamage = va.extractDamage(result.Findings)
	}

	return result, nil
}

// analyzeImage analyzes a single image
func (va *VisionAnalyzer) analyzeImage(ctx context.Context, req *DiagnosisRequest, attachment Attachment) (*VisualFinding, error) {
	// Build vision prompt
	prompt := va.buildVisionPrompt(req)

	// Prepare image input
	

	// Call AI vision API
	visionReq := &ai.VisionRequest{
		ImageData: []string{attachment.Base64Data},
		Prompt:      prompt,
		
	}

	// Use model from options if specified
	if req.Options.Model != "" {
		visionReq.Model = req.Options.Model
	}

	if req.Options.Temperature != nil {
		visionReq.Temperature = req.Options.Temperature
	}

	if req.Options.MaxTokens != nil {
		visionReq.MaxTokens = req.Options.MaxTokens
	}

	resp, err := va.aiManager.Analyze(ctx, visionReq)
	if err != nil {
		return nil, fmt.Errorf("AI vision analysis failed: %w", err)
	}

	// Parse vision response into structured finding
	finding := va.parseVisionResponse(attachment.ID, resp.Analysis)
	finding.Confidence = 0.85 // Default confidence

	return finding, nil
}

// buildVisionPrompt builds the prompt for vision analysis
func (va *VisionAnalyzer) buildVisionPrompt(req *DiagnosisRequest) string {
	prompt := fmt.Sprintf(`You are analyzing an image of medical equipment to help diagnose an issue.

**Equipment Information:**
- Type: %s`, req.EquipmentType)

	if req.Manufacturer != nil {
		prompt += fmt.Sprintf("\n- Manufacturer: %s", *req.Manufacturer)
	}

	if req.ModelNumber != nil {
		prompt += fmt.Sprintf("\n- Model: %s", *req.ModelNumber)
	}

	prompt += fmt.Sprintf(`
- Location: %s`, req.Location)

	if req.LocationType != nil {
		prompt += fmt.Sprintf(" (%s)", *req.LocationType)
	}

	prompt += fmt.Sprintf(`

**Reported Issue:**
%s

**Analysis Instructions:**
1. **Identify Equipment Components**: What parts/components are visible in the image?
2. **Detect Visual Issues**: Look for:
   - Physical damage (cracks, dents, burns, corrosion)
   - Error messages or warning lights on displays
   - Unusual wear or discoloration
   - Loose or disconnected parts
   - Fluid leaks
   - Overheating signs (discoloration, melted parts)
3. **Assess Condition**: Rate the overall condition
4. **Link to Reported Issue**: How does what you see relate to the reported problem?
5. **Diagnostic Clues**: What visual clues help diagnose the issue?

**Response Format:**
Provide a detailed analysis covering:
- Components visible
- Issues observed (be specific about location and severity)
- Condition assessment
- Diagnostic insights
- Recommended actions based on visual findings`, req.Description)

	return prompt
}

// parseVisionResponse parses AI vision response into structured finding
func (va *VisionAnalyzer) parseVisionResponse(attachmentID int64, analysis string) *VisualFinding {
	finding := &VisualFinding{
		AttachmentID: attachmentID,
		Finding:      analysis,
	}

	// Categorize based on keywords
	lowerAnalysis := strings.ToLower(analysis)

	if strings.Contains(lowerAnalysis, "damage") || strings.Contains(lowerAnalysis, "broken") ||
		strings.Contains(lowerAnalysis, "crack") || strings.Contains(lowerAnalysis, "dent") {
		finding.Category = "Damage"
	} else if strings.Contains(lowerAnalysis, "error") || strings.Contains(lowerAnalysis, "warning") ||
		strings.Contains(lowerAnalysis, "alert") {
		finding.Category = "Error Display"
	} else if strings.Contains(lowerAnalysis, "normal") || strings.Contains(lowerAnalysis, "good condition") {
		finding.Category = "Normal"
	} else if strings.Contains(lowerAnalysis, "wear") || strings.Contains(lowerAnalysis, "aged") {
		finding.Category = "Wear"
	} else {
		finding.Category = "General"
	}

	return finding
}

// generateOverallAssessment creates an overall assessment from findings
func (va *VisionAnalyzer) generateOverallAssessment(findings []VisualFinding) (string, float64) {
	if len(findings) == 0 {
		return "", 0
	}

	// Simple assessment based on categories
	damageCount := 0
	errorCount := 0
	wearCount := 0
	normalCount := 0

	var assessmentParts []string

	for _, finding := range findings {
		switch finding.Category {
		case "Damage":
			damageCount++
		case "Error Display":
			errorCount++
		case "Wear":
			wearCount++
		case "Normal":
			normalCount++
		}
	}

	// Build assessment
	if damageCount > 0 {
		assessmentParts = append(assessmentParts, fmt.Sprintf("%d image(s) show physical damage", damageCount))
	}
	if errorCount > 0 {
		assessmentParts = append(assessmentParts, fmt.Sprintf("%d image(s) show error displays or warnings", errorCount))
	}
	if wearCount > 0 {
		assessmentParts = append(assessmentParts, fmt.Sprintf("%d image(s) show wear or aging", wearCount))
	}
	if normalCount > 0 && damageCount == 0 && errorCount == 0 {
		assessmentParts = append(assessmentParts, fmt.Sprintf("%d image(s) appear normal", normalCount))
	}

	assessment := strings.Join(assessmentParts, "; ")

	// Calculate confidence (simple heuristic)
	confidence := 70.0 // Base confidence
	if damageCount > 0 || errorCount > 0 {
		confidence += 15.0 // Higher confidence if issues found
	}
	if len(findings) > 1 {
		confidence += 10.0 // Multiple images increase confidence
	}
	if confidence > 95.0 {
		confidence = 95.0
	}

	return assessment, confidence
}

// extractComponents extracts detected components from findings
func (va *VisionAnalyzer) extractComponents(findings []VisualFinding) []string {
	componentsMap := make(map[string]bool)

	// Common component keywords
	componentKeywords := []string{
		"display", "screen", "panel", "power supply", "circuit board",
		"fan", "cooling system", "cable", "connector", "button",
		"LED", "indicator light", "sensor", "motor", "control unit",
	}

	for _, finding := range findings {
		lowerFinding := strings.ToLower(finding.Finding)
		for _, keyword := range componentKeywords {
			if strings.Contains(lowerFinding, keyword) {
				componentsMap[keyword] = true
			}
		}
	}

	// Convert map to slice
	components := make([]string, 0, len(componentsMap))
	for component := range componentsMap {
		components = append(components, component)
	}

	return components
}

// extractDamage extracts damage descriptions from findings
func (va *VisionAnalyzer) extractDamage(findings []VisualFinding) []DamageDescription {
	var damages []DamageDescription

	for _, finding := range findings {
		lowerFinding := strings.ToLower(finding.Finding)

		// Check for various damage types
		if strings.Contains(lowerFinding, "crack") {
			damages = append(damages, DamageDescription{
				Type:        "Physical",
				Description: "Cracks detected",
				Severity:    va.inferSeverity(lowerFinding),
			})
		}

		if strings.Contains(lowerFinding, "burn") || strings.Contains(lowerFinding, "scorch") {
			damages = append(damages, DamageDescription{
				Type:        "Physical",
				Description: "Burn marks detected",
				Severity:    "Severe",
			})
		}

		if strings.Contains(lowerFinding, "corrosi") || strings.Contains(lowerFinding, "rust") {
			damages = append(damages, DamageDescription{
				Type:        "Physical",
				Description: "Corrosion detected",
				Severity:    va.inferSeverity(lowerFinding),
			})
		}

		if strings.Contains(lowerFinding, "error") && strings.Contains(lowerFinding, "display") {
			damages = append(damages, DamageDescription{
				Type:        "Display Error",
				Description: "Error message visible on display",
				Severity:    va.inferSeverity(lowerFinding),
			})
		}

		if strings.Contains(lowerFinding, "leak") {
			damages = append(damages, DamageDescription{
				Type:        "Physical",
				Description: "Fluid leak detected",
				Severity:    "Moderate",
			})
		}
	}

	return damages
}

// inferSeverity infers damage severity from text
func (va *VisionAnalyzer) inferSeverity(text string) string {
	if strings.Contains(text, "severe") || strings.Contains(text, "critical") ||
		strings.Contains(text, "major") || strings.Contains(text, "extensive") {
		return "Severe"
	}

	if strings.Contains(text, "moderate") || strings.Contains(text, "significant") {
		return "Moderate"
	}

	return "Minor"
}

// filterImageAttachments filters attachments to only images
func filterImageAttachments(attachments []Attachment) []Attachment {
	var images []Attachment

	for _, attachment := range attachments {
		if isImageContentType(attachment.ContentType) {
			images = append(images, attachment)
		}
	}

	return images
}

// isImageContentType checks if content type is an image
func isImageContentType(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
	}

	lowerType := strings.ToLower(contentType)
	for _, imgType := range imageTypes {
		if strings.Contains(lowerType, imgType) {
			return true
		}
	}

	return false
}


