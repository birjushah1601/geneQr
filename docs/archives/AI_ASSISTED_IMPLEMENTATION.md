# AI-Assisted Implementation Guide

## 🎯 Goal: Suggestion-Only AI with Feedback Loop

Transform the AI services from automatic execution to **AI-Assisted suggestions** with human decision-making and continuous learning through feedback.

---

## 📋 Implementation Phases

### Phase 1: Enhance Response Structures ✅ (Ready to implement)
### Phase 2: Add Feedback Capture API ✅ (Ready to implement)
### Phase 3: Implement Learning Loop ✅ (Ready to implement)
### Phase 4: Add Analytics Dashboard (Future)

---

## 🔧 Phase 1: Enhance Response Structures

### Current State
All AI services return direct results without confidence scoring or suggestion metadata.

### Target State
All AI responses include:
- ✅ AI suggestion with confidence score
- ✅ Decision status (pending/accepted/rejected)
- ✅ Feedback collection mechanism
- ✅ Alternative suggestions

---

### 1.1 Enhanced Diagnosis Response

**File:** `internal/diagnosis/types.go`

```go
// Add to DiagnosisResponse struct
type DiagnosisResponse struct {
	DiagnosisID        string            `json:"diagnosis_id"`
	TicketID           int64             `json:"ticket_id"`
	
	// AI Suggestion with confidence
	PrimaryDiagnosis   DiagnosisResult   `json:"primary_diagnosis"`
	Confidence         float64           `json:"confidence"`          // 0.0 - 1.0
	ConfidenceLevel    string            `json:"confidence_level"`    // HIGH/MEDIUM/LOW
	
	// Alternatives
	AlternateDiagnoses []DiagnosisResult `json:"alternate_diagnoses"`
	
	// Actions and parts
	RecommendedActions []RecommendedAction `json:"recommended_actions"`
	RequiredParts      []RequiredPart      `json:"required_parts"`
	
	// AI Metadata
	AIMetadata         AISuggestionMetadata `json:"ai_metadata"`
	
	// Decision tracking
	DecisionStatus     string            `json:"decision_status"`     // pending/accepted/rejected
	DecidedBy          *int64            `json:"decided_by,omitempty"` // user_id
	DecidedAt          *time.Time        `json:"decided_at,omitempty"`
	FeedbackText       string            `json:"feedback_text,omitempty"`
	
	// Existing fields...
	VisionAnalysis     *VisionAnalysisResult `json:"vision_analysis,omitempty"`
	ContextUsed        map[string]interface{} `json:"context_used,omitempty"`
	Metadata           DiagnosisMetadata      `json:"metadata"`
	CreatedAt          time.Time              `json:"created_at"`
}

// New struct for AI suggestion metadata
type AISuggestionMetadata struct {
	Provider           string        `json:"provider"`            // openai/anthropic
	Model              string        `json:"model"`
	Confidence         float64       `json:"confidence"`
	ConfidenceFactors  []string      `json:"confidence_factors"`  // What influenced confidence
	AlternativesCount  int           `json:"alternatives_count"`
	RequiresFeedback   bool          `json:"requires_feedback"`   // If confidence < threshold
	SuggestionOnly     bool          `json:"suggestion_only"`     // Always true for now
}
```

**Confidence Calculation:**
```go
// Add to internal/diagnosis/engine.go

func (e *Engine) calculateConfidence(result *DiagnosisResponse) float64 {
	score := 0.0
	factors := 0
	
	// Factor 1: Vision analysis confidence (if available)
	if result.VisionAnalysis != nil && len(result.VisionAnalysis.Findings) > 0 {
		// Use default confidence of 0.85 that was set
		score += 0.85
		factors++
	}
	
	// Factor 2: Historical match (if similar tickets found)
	if result.ContextUsed != nil {
		if similarCount, ok := result.ContextUsed["similar_tickets_count"].(int); ok && similarCount > 0 {
			score += 0.80
			factors++
		}
	}
	
	// Factor 3: Symptom clarity (more symptoms = higher confidence)
	if len(result.PrimaryDiagnosis.Symptoms) >= 3 {
		score += 0.75
		factors++
	}
	
	// Factor 4: AI model confidence (based on response characteristics)
	// Check if diagnosis has specific details (not generic)
	if len(result.PrimaryDiagnosis.RootCause) > 50 {
		score += 0.70
		factors++
	}
	
	if factors == 0 {
		return 0.5 // Default medium confidence
	}
	
	return score / float64(factors)
}

func (e *Engine) getConfidenceLevel(confidence float64) string {
	if confidence >= 0.80 {
		return "HIGH"
	} else if confidence >= 0.60 {
		return "MEDIUM"
	}
	return "LOW"
}

func (e *Engine) getConfidenceFactors(result *DiagnosisResponse) []string {
	factors := []string{}
	
	if result.VisionAnalysis != nil {
		factors = append(factors, "Visual analysis of equipment images")
	}
	
	if result.ContextUsed != nil {
		if count, ok := result.ContextUsed["similar_tickets_count"].(int); ok && count > 0 {
			factors = append(factors, fmt.Sprintf("Matched with %d similar historical cases", count))
		}
	}
	
	if len(result.PrimaryDiagnosis.Symptoms) >= 3 {
		factors = append(factors, "Multiple symptoms analyzed")
	}
	
	if result.Metadata.Provider != "" {
		factors = append(factors, fmt.Sprintf("AI model: %s", result.Metadata.Model))
	}
	
	return factors
}
```

---

### 1.2 Enhanced Assignment Response

**File:** `internal/assignment/types.go`

```go
// Add to AssignmentResponse struct
type AssignmentResponse struct {
	RequestID          string                     `json:"request_id"`
	TicketID           int64                      `json:"ticket_id"`
	
	// AI Recommendations with confidence
	Recommendations    []EngineerRecommendation   `json:"recommendations"`
	TopRecommendation  *EngineerRecommendation    `json:"top_recommendation"`
	Confidence         float64                    `json:"confidence"`
	ConfidenceLevel    string                     `json:"confidence_level"`
	
	// AI Metadata
	AIMetadata         AISuggestionMetadata       `json:"ai_metadata"`
	
	// Decision tracking
	DecisionStatus     string                     `json:"decision_status"`     // pending/accepted/rejected
	SelectedEngineer   *int64                     `json:"selected_engineer,omitempty"`
	DecidedBy          *int64                     `json:"decided_by,omitempty"`
	DecidedAt          *time.Time                 `json:"decided_at,omitempty"`
	FeedbackText       string                     `json:"feedback_text,omitempty"`
	
	// Metadata
	Metadata           AssignmentMetadata         `json:"metadata"`
	CreatedAt          time.Time                  `json:"created_at"`
}

// Enhanced EngineerRecommendation
type EngineerRecommendation struct {
	EngineerID         int64                      `json:"engineer_id"`
	EngineerName       string                     `json:"engineer_name"`
	MatchScore         float64                    `json:"match_score"`
	Confidence         float64                    `json:"confidence"`          // NEW
	Rank               int                        `json:"rank"`
	
	// Scoring breakdown
	Scores             ScoringBreakdown           `json:"scores"`
	
	// Availability
	IsAvailable        bool                       `json:"is_available"`
	AvailableFrom      *time.Time                 `json:"available_from,omitempty"`
	CurrentLoad        int                        `json:"current_load"`
	
	// Experience
	RelevantExperience []EquipmentExpertise       `json:"relevant_experience"`
	SuccessRate        float64                    `json:"success_rate"`
	AverageResolutionTime float64                 `json:"avg_resolution_time_hours"`
	
	// Location
	DistanceKM         float64                    `json:"distance_km"`
	EstimatedTravelMin int                        `json:"estimated_travel_min"`
	
	// Reasoning (NEW)
	Reasoning          []string                   `json:"reasoning"`           // Why recommended
	Concerns           []string                   `json:"concerns,omitempty"`  // Why not perfect
}
```

---

### 1.3 Enhanced Parts Response

**File:** `internal/parts/types.go`

```go
// Add to PartsResponse struct
type PartsResponse struct {
	RequestID          string                     `json:"request_id"`
	TicketID           int64                      `json:"ticket_id"`
	DiagnosisID        string                     `json:"diagnosis_id,omitempty"`
	
	// AI Recommendations
	RecommendedParts   []PartRecommendation       `json:"recommended_parts"`
	TotalEstimatedCost float64                    `json:"total_estimated_cost"`
	Confidence         float64                    `json:"confidence"`
	ConfidenceLevel    string                     `json:"confidence_level"`
	
	// AI Metadata
	AIMetadata         AISuggestionMetadata       `json:"ai_metadata"`
	
	// Decision tracking
	DecisionStatus     string                     `json:"decision_status"`     // pending/accepted/rejected
	ApprovedParts      []int64                    `json:"approved_parts,omitempty"` // part IDs
	DecidedBy          *int64                     `json:"decided_by,omitempty"`
	DecidedAt          *time.Time                 `json:"decided_at,omitempty"`
	FeedbackText       string                     `json:"feedback_text,omitempty"`
	
	// Metadata
	Metadata           PartsMetadata              `json:"metadata"`
	CreatedAt          time.Time                  `json:"created_at"`
}

// Enhanced PartRecommendation
type PartRecommendation struct {
	PartID             int64                      `json:"part_id"`
	PartNumber         string                     `json:"part_number"`
	PartName           string                     `json:"part_name"`
	Quantity           int                        `json:"quantity"`
	UnitPrice          float64                    `json:"unit_price"`
	TotalPrice         float64                    `json:"total_price"`
	
	// AI Confidence
	Confidence         float64                    `json:"confidence"`          // NEW
	Necessity          string                     `json:"necessity"`           // CRITICAL/RECOMMENDED/OPTIONAL
	
	// Reasoning (NEW)
	Reason             string                     `json:"reason"`              // Why needed
	AlternativeParts   []AlternativePart          `json:"alternatives,omitempty"`
	
	// Availability
	InStock            bool                       `json:"in_stock"`
	StockQuantity      int                        `json:"stock_quantity"`
	LeadTimeDays       int                        `json:"lead_time_days"`
	
	// Compatibility
	CompatibleModels   []string                   `json:"compatible_models"`
	IsOEM              bool                       `json:"is_oem"`
}
```

---

## 🔧 Phase 2: Feedback Capture API

### 2.1 Unified Feedback Endpoint

**File:** `internal/api/feedback_handler.go` (already exists, enhance it)

```go
// Add new endpoint for AI decision feedback
func (h *FeedbackHandler) RegisterRoutes(r *mux.Router) {
	// ... existing routes ...
	
	// NEW: Unified AI decision feedback
	r.HandleFunc("/api/ai-decisions/{serviceType}/{requestId}/feedback", h.SubmitAIDecisionFeedback).Methods("POST")
}

// SubmitAIDecisionFeedback handles feedback when user accepts/rejects AI suggestion
func (h *FeedbackHandler) SubmitAIDecisionFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceType := vars["serviceType"]   // diagnosis, assignment, parts
	requestID := vars["requestId"]
	
	var req AIDecisionFeedback
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}
	
	// Validate service type
	if serviceType != "diagnosis" && serviceType != "assignment" && serviceType != "parts" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid service type",
			Message: "Must be one of: diagnosis, assignment, parts",
		})
		return
	}
	
	// Process feedback
	feedbackReq := &feedback.HumanFeedbackRequest{
		ServiceType: serviceType,
		RequestID:   requestID,
		TicketID:    req.TicketID,
		UserID:      req.UserID,
		UserRole:    req.UserRole,
		WasAccurate: req.Decision == "accepted",
		Comments:    req.FeedbackText,
		Corrections: req.Corrections,
	}
	
	// Add rating based on decision
	if req.Decision == "accepted" {
		rating := 5
		feedbackReq.Rating = &rating
	} else {
		rating := 2
		feedbackReq.Rating = &rating
	}
	
	entry, err := h.collector.CollectHumanFeedback(r.Context(), feedbackReq)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to submit feedback",
			Message: err.Error(),
		})
		return
	}
	
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success":     true,
		"feedback_id": entry.FeedbackID,
		"message":     "Thank you for your feedback! This helps our AI learn and improve.",
	})
}
```

### 2.2 Feedback Request Types

**File:** `internal/api/types.go` (create if doesn't exist)

```go
package api

// AIDecisionFeedback represents feedback on an AI suggestion
type AIDecisionFeedback struct {
	// Decision
	Decision     string `json:"decision"`      // "accepted" or "rejected"
	
	// Context
	TicketID     *int64 `json:"ticket_id,omitempty"`
	UserID       int64  `json:"user_id"`
	UserRole     string `json:"user_role"`     // engineer, dispatcher, admin
	
	// Feedback (required if rejected)
	FeedbackText string `json:"feedback_text,omitempty"`
	
	// Corrections (what should have been suggested)
	Corrections  map[string]interface{} `json:"corrections,omitempty"`
}

// Example corrections for diagnosis:
// {
//   "correct_diagnosis": "Sensor malfunction",
//   "correct_root_cause": "Dust accumulation on sensor",
//   "correct_parts": ["SENSOR-X100", "CLEANING-KIT"]
// }

// Example corrections for assignment:
// {
//   "correct_engineer_id": 456,
//   "reason": "This engineer has certification for this specific equipment model"
// }

// Example corrections for parts:
// {
//   "correct_parts": [
//     {"part_number": "ABC-123", "quantity": 2, "reason": "Original part is outdated"}
//   ]
// }
```

---

## 🔧 Phase 3: Update Engine Logic

### 3.1 Update Diagnosis Engine

**File:** `internal/diagnosis/engine.go`

Add confidence calculation after diagnosis:

```go
func (e *Engine) Diagnose(ctx context.Context, req *DiagnosisRequest) (*DiagnosisResponse, error) {
	// ... existing diagnosis logic ...
	
	// NEW: Calculate confidence
	confidence := e.calculateConfidence(response)
	response.Confidence = confidence
	response.ConfidenceLevel = e.getConfidenceLevel(confidence)
	
	// NEW: Set AI metadata
	response.AIMetadata = AISuggestionMetadata{
		Provider:          response.Metadata.Provider,
		Model:             response.Metadata.Model,
		Confidence:        confidence,
		ConfidenceFactors: e.getConfidenceFactors(response),
		AlternativesCount: len(response.AlternateDiagnoses),
		RequiresFeedback:  confidence < 0.80, // Require feedback if not high confidence
		SuggestionOnly:    true,              // Always true for now
	}
	
	// NEW: Initialize as pending decision
	response.DecisionStatus = "pending"
	
	return response, nil
}
```

### 3.2 Update Assignment Engine

**File:** `internal/assignment/engine.go`

```go
func (e *Engine) RecommendEngineers(ctx context.Context, req *AssignmentRequest) (*AssignmentResponse, error) {
	// ... existing recommendation logic ...
	
	// NEW: Calculate overall confidence based on top recommendation
	var confidence float64
	if len(response.Recommendations) > 0 {
		topRec := response.Recommendations[0]
		confidence = topRec.MatchScore // Use match score as confidence
		
		// Adjust based on availability and distance
		if !topRec.IsAvailable {
			confidence *= 0.8
		}
		if topRec.DistanceKM > 50 {
			confidence *= 0.9
		}
		
		response.TopRecommendation = &topRec
	} else {
		confidence = 0.0
	}
	
	response.Confidence = confidence
	response.ConfidenceLevel = e.getConfidenceLevel(confidence)
	
	// NEW: Set AI metadata
	response.AIMetadata = AISuggestionMetadata{
		Provider:          response.Metadata.Provider,
		Model:             response.Metadata.Model,
		Confidence:        confidence,
		ConfidenceFactors: e.getConfidenceFactors(response),
		AlternativesCount: len(response.Recommendations) - 1,
		RequiresFeedback:  confidence < 0.85,
		SuggestionOnly:    true,
	}
	
	// NEW: Initialize as pending
	response.DecisionStatus = "pending"
	
	return response, nil
}
```

---

## 📱 Frontend Integration Example

### Example API Response (Diagnosis)

```json
{
  "diagnosis_id": "diag_abc123",
  "ticket_id": 12345,
  "confidence": 0.87,
  "confidence_level": "HIGH",
  "primary_diagnosis": {
    "issue": "Power supply failure",
    "root_cause": "Capacitor degradation",
    "severity": "High"
  },
  "alternate_diagnoses": [
    {
      "issue": "Motherboard fault",
      "probability": 0.15
    }
  ],
  "ai_metadata": {
    "provider": "openai",
    "model": "gpt-4o",
    "confidence": 0.87,
    "confidence_factors": [
      "Visual analysis of equipment images",
      "Matched with 5 similar historical cases",
      "Multiple symptoms analyzed"
    ],
    "requires_feedback": false,
    "suggestion_only": true
  },
  "decision_status": "pending",
  "created_at": "2025-11-18T10:30:00Z"
}
```

### Example Frontend UI Flow

```typescript
// 1. Display AI Suggestion
<Card>
  <CardHeader>
    <h3>AI Diagnosis Suggestion</h3>
    <Badge color={getConfidenceColor(diagnosis.confidence_level)}>
      {diagnosis.confidence_level} Confidence ({(diagnosis.confidence * 100).toFixed(0)}%)
    </Badge>
  </CardHeader>
  
  <CardBody>
    <h4>{diagnosis.primary_diagnosis.issue}</h4>
    <p>{diagnosis.primary_diagnosis.root_cause}</p>
    
    <div className="confidence-factors">
      <h5>Based on:</h5>
      <ul>
        {diagnosis.ai_metadata.confidence_factors.map(factor => (
          <li key={factor}>{factor}</li>
        ))}
      </ul>
    </div>
    
    {diagnosis.alternate_diagnoses.length > 0 && (
      <div className="alternatives">
        <h5>Alternative diagnoses:</h5>
        {/* Show alternatives */}
      </div>
    )}
  </CardBody>
  
  <CardFooter>
    <Button onClick={() => acceptSuggestion(diagnosis.diagnosis_id)}>
      Accept Suggestion
    </Button>
    <Button variant="secondary" onClick={() => showFeedbackForm()}>
      Reject & Provide Feedback
    </Button>
  </CardFooter>
</Card>

// 2. Feedback Form (shown when rejected)
<FeedbackModal>
  <h3>Why are you rejecting this suggestion?</h3>
  <Textarea 
    placeholder="Please explain what the AI got wrong and what the correct diagnosis should be..."
    value={feedbackText}
    onChange={(e) => setFeedbackText(e.target.value)}
  />
  
  <h4>What is the correct diagnosis?</h4>
  <Input placeholder="Correct issue..." />
  <Input placeholder="Correct root cause..." />
  
  <Button onClick={submitFeedback}>Submit Feedback</Button>
</FeedbackModal>
```

---

## 🎯 Implementation Checklist

### ✅ Quick Wins (1-2 days)
- [ ] Add confidence scoring to all response types
- [ ] Add `decision_status` field to all responses
- [ ] Create unified feedback endpoint
- [ ] Update database schemas to store decision status

### ⚙️ Core Features (3-5 days)
- [ ] Implement confidence calculation for diagnosis
- [ ] Implement confidence calculation for assignment
- [ ] Implement confidence calculation for parts
- [ ] Add confidence factors explanation
- [ ] Wire up feedback collection to learning engine

### 📊 Learning Loop (1 week)
- [ ] Enhance feedback analyzer to process accept/reject patterns
- [ ] Update learning engine to identify improvement opportunities
- [ ] Create periodic learning reports
- [ ] Add A/B testing framework for model improvements

### 📈 Analytics (Ongoing)
- [ ] Track acceptance rate by confidence level
- [ ] Track feedback quality and patterns
- [ ] Monitor AI accuracy improvement over time
- [ ] Create dashboard showing AI vs Human decisions

---

## 🔄 Learning Loop Flow

```
┌─────────────────┐
│  AI Suggestion  │
│  (with conf.)   │
└────────┬────────┘
         │
    ┌────▼─────┐
    │  Human   │
    │ Decision │
    └────┬─────┘
         │
    ┌────▼──────────┐
    │   Accept?     │
    └───┬───────┬───┘
        │       │
   YES  │       │ NO
        │       │
        ▼       ▼
    ┌───────┐ ┌──────────┐
    │Track  │ │ Collect  │
    │Success│ │ Feedback │
    └───┬───┘ └────┬─────┘
        │          │
        └────┬─────┘
             ▼
    ┌────────────────┐
    │ Feedback DB    │
    │ (both sources) │
    └────────┬───────┘
             │
             ▼
    ┌────────────────┐
    │   Analyzer     │
    │  (patterns)    │
    └────────┬───────┘
             │
             ▼
    ┌────────────────┐
    │ Learning Engine│
    │ (improvements) │
    └────────┬───────┘
             │
             ▼
    ┌────────────────┐
    │  Model Update  │
    │ (prompt tuning)│
    └────────────────┘
```

---

## 📊 Success Metrics

### Week 1-2
- [ ] All AI responses include confidence scores
- [ ] Feedback form deployed and functional
- [ ] 100% of decisions tracked in database

### Month 1
- [ ] 80%+ feedback submission rate on rejections
- [ ] Confidence calibration (high conf. = high acceptance)
- [ ] Initial learning patterns identified

### Month 3
- [ ] 85%+ acceptance rate for HIGH confidence suggestions
- [ ] <20% override rate on accepted suggestions
- [ ] Measurable improvement in AI accuracy

---

## 🚀 Next Steps

1. **Review this document** - Make sure this aligns with your vision
2. **Prioritize features** - Which phase do you want to start with?
3. **Database updates** - Update schemas to support new fields
4. **Implement Phase 1** - Enhanced response structures
5. **Test & Iterate** - Start collecting real feedback

**Would you like me to start implementing any specific phase now?**
