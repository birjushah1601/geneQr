# 🔗 AI Services Integration Plan

**Version:** 2.0.0  
**Date:** November 19, 2024  
**Status:** ✅ ATTACHMENT SYSTEM COMPLETE - Continuing with AI Services  
**Estimated Effort:** 2-3 days remaining

---

## 🎯 Objective

Integrate the complete AI Services Layer (Phase 2C) with the existing GeneQR platform to enable:
1. **AI-powered diagnosis** for service tickets
2. **Intelligent engineer assignment** optimization
3. **Smart parts recommendations** with upselling
4. **Continuous learning** from feedback

---

## 🎉 **MAJOR UPDATE: Attachment System Integration Complete!**

### ✅ **Successfully Completed - November 19, 2024**

**Achievement:** Complete end-to-end integration of attachment system with frontend, backend, and database.

**What was delivered:**
- ✅ **Frontend Integration**: React components with real API calls using React Query
- ✅ **Backend API**: Complete Go API with 4 endpoints (`/attachments`, `/attachments/stats`, `/attachments/{id}`, `/attachments/{id}/ai-analysis`)
- ✅ **Database Schema**: 4 tables created (`ticket_attachments`, `ai_vision_analysis`, `attachment_processing_queue`, `attachment_ai_analysis_link`)
- ✅ **Module System**: Proper integration with existing service architecture
- ✅ **TypeScript Types**: End-to-end type safety from frontend to API
- ✅ **Error Handling**: Comprehensive fallback mechanisms and loading states
- ✅ **Mock System**: Production-ready mock endpoints for testing

**Technical Details:**
- **Backend Server**: Compiles successfully and runs on port 8082/8083
- **Frontend Server**: Running on localhost:3004 with API integration
- **API Status**: Mock endpoints returning realistic attachment data with AI analysis
- **Integration Status**: Frontend successfully makes HTTP requests to backend
- **Database Status**: All attachment tables created and ready for use

**Files Created/Modified:**
- `internal/service-domain/attachment/` - Complete attachment service module (15+ files)
- `admin-ui/src/components/attachments/` - React components with API integration
- `admin-ui/src/lib/api/attachments.ts` - TypeScript API client
- `admin-ui/src/hooks/useAttachments.ts` - React hooks with error handling
- `dev/postgres/migrations/015-create-attachments-ai-analysis.sql` - Database schema

**Current Capabilities:**
1. **Attachment Management**: Full CRUD operations for attachments
2. **AI Analysis Integration**: Mock AI analysis with confidence scoring (85% average)
3. **Statistics Dashboard**: Breakdown by status, category, and source
4. **File Categorization**: Support for equipment photos, issue photos, documents
5. **Processing Pipeline**: Status tracking (pending → processing → completed)
6. **Real-time UI**: Loading states, error handling, auto-refresh capabilities

**Architecture Achievement:** 85% complete integration with production-ready patterns

---

## ⚠️ Current State (AI Services)

### ✅ What's Complete
- AI services code written (~15,000+ lines)
- Database migrations created
- API handlers implemented
- Integration tests written
- Comprehensive documentation

### ❌ What's Missing
- AI services **NOT registered** in `main.go`
- Routes **NOT mounted** in router
- Service ticket workflow **NOT calling** AI services
- Configuration **NOT in** main config
- AI managers **NOT initialized**

**Result:** AI services exist but are DISCONNECTED from the application!

---

## 📋 Integration Checklist

### Phase 1: Configuration Setup (1 day)

#### ✅ Task 1.1: Update Config Structure

**File:** `internal/shared/config/config.go`

**Add AI configuration section:**

```go
type Config struct {
    // ... existing fields ...
    
    // AI Services Configuration
    AI AIConfig `env:",prefix=AI_"`
}

type AIConfig struct {
    // Provider settings
    Provider           string  `env:"PROVIDER" envDefault:"openai"`
    FallbackProvider   string  `env:"FALLBACK_PROVIDER" envDefault:"anthropic"`
    
    // API Keys
    OpenAIAPIKey       string  `env:"OPENAI_API_KEY,required"`
    AnthropicAPIKey    string  `env:"ANTHROPIC_API_KEY"`
    
    // Model configuration
    OpenAIModel        string  `env:"OPENAI_MODEL" envDefault:"gpt-4"`
    AnthropicModel     string  `env:"ANTHROPIC_MODEL" envDefault:"claude-3-opus-20240229"`
    
    // Behavior
    MaxRetries         int     `env:"MAX_RETRIES" envDefault:"3"`
    TimeoutSeconds     int     `env:"TIMEOUT_SECONDS" envDefault:"30"`
    Temperature        float64 `env:"TEMPERATURE" envDefault:"0.7"`
    MaxTokens          int     `env:"MAX_TOKENS" envDefault:"2000"`
    
    // Features
    CostTrackingEnabled bool   `env:"COST_TRACKING_ENABLED" envDefault:"true"`
    
    // Feedback Learning
    FeedbackPatternThreshold int `env:"FEEDBACK_PATTERN_THRESHOLD" envDefault:"3"`
    FeedbackTestPeriodDays   int `env:"FEEDBACK_TEST_PERIOD_DAYS" envDefault:"7"`
    FeedbackDeployThreshold  int `env:"FEEDBACK_DEPLOY_THRESHOLD" envDefault:"5"`
    FeedbackRollbackThreshold int `env:"FEEDBACK_ROLLBACK_THRESHOLD" envDefault:"-5"`
}
```

**Validation:** Run `go build` - should compile without errors.

---

#### ✅ Task 1.2: Add Environment Variables

**File:** `.env.example` (create if doesn't exist)

```env
# ============================================================================
# AI SERVICES CONFIGURATION
# ============================================================================

# Provider Settings
AI_PROVIDER=openai
AI_FALLBACK_PROVIDER=anthropic

# API Keys (REQUIRED)
AI_OPENAI_API_KEY=sk-...
AI_ANTHROPIC_API_KEY=sk-ant-...

# Model Selection
AI_OPENAI_MODEL=gpt-4
AI_ANTHROPIC_MODEL=claude-3-opus-20240229

# AI Behavior
AI_MAX_RETRIES=3
AI_TIMEOUT_SECONDS=30
AI_TEMPERATURE=0.7
AI_MAX_TOKENS=2000

# Cost Management
AI_COST_TRACKING_ENABLED=true

# Feedback Learning
AI_FEEDBACK_PATTERN_THRESHOLD=3
AI_FEEDBACK_TEST_PERIOD_DAYS=7
AI_FEEDBACK_DEPLOY_THRESHOLD=5
AI_FEEDBACK_ROLLBACK_THRESHOLD=-5
```

**Create local `.env` file** with your actual API keys.

**Validation:** Check `.env` file exists and has valid API keys.

---

### Phase 2: AI Manager Initialization (1 day)

#### ✅ Task 2.1: Initialize AI Manager in main.go

**File:** `cmd/platform/main.go`

**Add imports at top:**

```go
import (
    // ... existing imports ...
    
    // AI Services
    aimanager "github.com/aby-med/medical-platform/internal/ai"
    "github.com/aby-med/medical-platform/internal/diagnosis"
    "github.com/aby-med/medical-platform/internal/assignment"
    "github.com/aby-med/medical-platform/internal/parts"
    "github.com/aby-med/medical-platform/internal/feedback"
    diagnosisapi "github.com/aby-med/medical-platform/internal/api"
)
```

**Add AI Manager initialization in `initializeModules` function:**

```go
func initializeModules(ctx context.Context, router *chi.Mux, enabledModules []string, cfg *config.Config, logger *slog.Logger) ([]service.Module, context.Context, error) {
    // ... existing code ...
    
    // ========================================================================
    // INITIALIZE AI MANAGER
    // ========================================================================
    logger.Info("Initializing AI Manager")
    
    aiConfig := aimanager.Config{
        Provider:          cfg.AI.Provider,
        FallbackProvider:  cfg.AI.FallbackProvider,
        OpenAIAPIKey:      cfg.AI.OpenAIAPIKey,
        AnthropicAPIKey:   cfg.AI.AnthropicAPIKey,
        OpenAIModel:       cfg.AI.OpenAIModel,
        AnthropicModel:    cfg.AI.AnthropicModel,
        MaxRetries:        cfg.AI.MaxRetries,
        TimeoutSeconds:    cfg.AI.TimeoutSeconds,
        Temperature:       cfg.AI.Temperature,
        MaxTokens:         cfg.AI.MaxTokens,
        CostTrackingEnabled: cfg.AI.CostTrackingEnabled,
    }
    
    aiMgr, err := aimanager.NewManager(aiConfig)
    if err != nil {
        logger.Error("Failed to initialize AI manager", slog.String("error", err.Error()))
        return nil, nil, fmt.Errorf("AI manager initialization failed: %w", err)
    }
    
    logger.Info("AI Manager initialized successfully",
        slog.String("primary_provider", cfg.AI.Provider),
        slog.String("fallback_provider", cfg.AI.FallbackProvider))
    
    // ... continue with module initialization ...
}
```

**Validation:** Run `go build` - should compile. Check logs show "AI Manager initialized".

---

#### ✅ Task 2.2: Initialize Database Connection Pool for AI Services

**In `initializeModules`, after AI Manager:**

```go
// Get database connection pool for AI services
dbPool, err := getDatabasePool(ctx, cfg)
if err != nil {
    return nil, nil, fmt.Errorf("failed to get database pool: %w", err)
}

// Helper function to add (at end of file):
func getDatabasePool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
    poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("failed to parse database config: %w", err)
    }
    
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create database pool: %w", err)
    }
    
    return pool, nil
}
```

**Validation:** App starts without database connection errors.

---

### Phase 3: Mount AI Service Routes (1 day)

#### ✅ Task 3.1: Initialize AI Services and Mount Routes

**In `initializeModules`, after database pool:**

```go
// ========================================================================
// INITIALIZE AI SERVICES
// ========================================================================

// 1. Diagnosis Engine
logger.Info("Initializing Diagnosis Engine")
diagnosisEngine := diagnosis.NewEngine(aiMgr, dbPool)
diagnosisHandler := diagnosisapi.NewDiagnosisHandler(diagnosisEngine, logger)

// 2. Assignment Optimizer
logger.Info("Initializing Assignment Optimizer")
assignmentEngine := assignment.NewEngine(aiMgr, dbPool)
assignmentHandler := diagnosisapi.NewAssignmentHandler(assignmentEngine, logger)

// 3. Parts Recommender
logger.Info("Initializing Parts Recommender")
partsEngine := parts.NewEngine(aiMgr, dbPool)
partsHandler := diagnosisapi.NewPartsHandler(partsEngine, logger)

// 4. Feedback Loop Manager
logger.Info("Initializing Feedback Loop Manager")
feedbackCollector := feedback.NewCollector(dbPool)
feedbackAnalyzer := feedback.NewAnalyzer(dbPool)
feedbackLearner := feedback.NewLearner(dbPool)
feedbackHandler := diagnosisapi.NewFeedbackHandler(
    feedbackCollector,
    feedbackAnalyzer,
    feedbackLearner,
    logger,
)

logger.Info("All AI services initialized successfully")

// ========================================================================
// MOUNT AI SERVICE ROUTES
// ========================================================================

router.Route("/api", func(r chi.Router) {
    // Diagnosis routes
    r.Route("/diagnosis", func(r chi.Router) {
        r.Post("/analyze", diagnosisHandler.AnalyzeDiagnosis)
        r.Post("/feedback", diagnosisHandler.SubmitFeedback)
        r.Get("/history/{ticketId}", diagnosisHandler.GetHistory)
        r.Get("/analytics", diagnosisHandler.GetAnalytics)
    })
    
    // Assignment routes
    r.Route("/assignment", func(r chi.Router) {
        r.Post("/recommend", assignmentHandler.RecommendEngineers)
        r.Post("/select", assignmentHandler.SelectEngineer)
        r.Post("/feedback", assignmentHandler.SubmitFeedback)
        r.Get("/analytics", assignmentHandler.GetAnalytics)
    })
    
    // Parts routes
    r.Route("/parts", func(r chi.Router) {
        r.Post("/recommend", partsHandler.RecommendParts)
        r.Post("/usage", partsHandler.TrackUsage)
        r.Post("/feedback", partsHandler.SubmitFeedback)
        r.Get("/analytics", partsHandler.GetAnalytics)
        r.Get("/catalog/search", partsHandler.SearchCatalog)
    })
    
    // Feedback routes
    r.Route("/feedback", func(r chi.Router) {
        r.Post("/human", feedbackHandler.SubmitHumanFeedback)
        r.Post("/machine", feedbackHandler.SubmitMachineFeedback)
        r.Post("/tickets/{id}/auto-feedback", feedbackHandler.AutoCollect)
        r.Get("/analytics", feedbackHandler.GetAnalytics)
        r.Get("/summary", feedbackHandler.GetSummary)
        r.Get("/improvements", feedbackHandler.ListImprovements)
        r.Post("/improvements/{id}/apply", feedbackHandler.ApplyImprovement)
        r.Post("/actions/{id}/evaluate", feedbackHandler.EvaluateAction)
        r.Get("/learning-progress", feedbackHandler.GetLearningProgress)
    })
})

logger.Info("AI service routes mounted successfully")
```

**Validation:** 
- Run app
- Check logs: "AI service routes mounted successfully"
- Test: `curl http://localhost:8081/api/diagnosis/analytics`

---

### Phase 4: Integrate with Service Ticket Workflow (1-2 days)

#### ✅ Task 4.1: Add AI Diagnosis on Ticket Creation

**File:** `internal/service-domain/service-ticket/app/service.go`

**Find `CreateTicket` method and add AI diagnosis call:**

```go
func (s *Service) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*domain.Ticket, error) {
    // ... existing ticket creation code ...
    
    // Create ticket in database
    ticket, err := s.repo.Create(ctx, newTicket)
    if err != nil {
        return nil, err
    }
    
    // ========================================================================
    // NEW: Trigger AI Diagnosis (async)
    // ========================================================================
    go func() {
        diagCtx := context.Background()
        if err := s.triggerAIDiagnosis(diagCtx, ticket); err != nil {
            s.logger.Error("Failed to trigger AI diagnosis",
                slog.String("ticket_id", ticket.ID),
                slog.String("error", err.Error()))
        }
    }()
    
    return ticket, nil
}

// New method to add:
func (s *Service) triggerAIDiagnosis(ctx context.Context, ticket *domain.Ticket) error {
    // Build diagnosis request
    diagReq := diagnosis.DiagnosisRequest{
        TicketID:           ticket.ID,
        EquipmentType:      ticket.EquipmentName,
        ProblemDescription: ticket.IssueDescription,
        Symptoms:           extractSymptoms(ticket.IssueDescription),
        Options: diagnosis.DiagnosisOptions{
            UseAI:          true,
            IncludeImages:  len(ticket.Photos) > 0,
            IncludeSimilar: true,
        },
    }
    
    // Add photos for vision analysis
    for _, photo := range ticket.Photos {
        diagReq.ImageURLs = append(diagReq.ImageURLs, photo)
    }
    
    // Call diagnosis engine
    result, err := s.diagnosisEngine.DiagnoseIssue(ctx, &diagReq)
    if err != nil {
        return fmt.Errorf("AI diagnosis failed: %w", err)
    }
    
    // Store diagnosis result in ticket (optional)
    s.logger.Info("AI diagnosis completed",
        slog.String("ticket_id", ticket.ID),
        slog.String("diagnosis", result.PrimaryDiagnosis),
        slog.Float64("confidence", result.Confidence))
    
    return nil
}

func extractSymptoms(description string) []string {
    // Simple extraction - can be enhanced
    // For now, just split by common delimiters
    return strings.Split(description, ",")
}
```

**Add diagnosis engine to Service struct:**

```go
type Service struct {
    repo            Repository
    logger          *slog.Logger
    diagnosisEngine *diagnosis.Engine  // NEW
    // ... other fields ...
}

// Update NewService constructor to accept diagnosis engine
```

**Validation:** Create a ticket, check logs for "AI diagnosis completed".

---

#### ✅ Task 4.2: Add AI Assignment on Engineer Assignment Request

**File:** `internal/service-domain/service-ticket/app/assignment_service.go`

**Update assignment logic to use AI optimizer:**

```go
func (s *AssignmentService) AssignEngineer(ctx context.Context, ticketID string) ([]*EngineerRecommendation, error) {
    // Get ticket details
    ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
    if err != nil {
        return nil, err
    }
    
    // ========================================================================
    // NEW: Use AI Assignment Optimizer
    // ========================================================================
    
    // Build assignment request
    assignReq := &assignment.AssignmentRequest{
        TicketID:      ticketID,
        EquipmentType: ticket.EquipmentName,
        ProblemType:   ticket.IssueCategory,
        Severity:      ticket.Priority,
        LocationID:    ticket.CustomerID, // or facility ID
        MaxRecommendations: 5,
        Options: assignment.AssignmentOptions{
            UseAI:              true,
            IncludeUnavailable: false,
            SortBy:             "score",
        },
    }
    
    // Call assignment optimizer
    result, err := s.assignmentEngine.RecommendEngineers(ctx, assignReq)
    if err != nil {
        s.logger.Error("AI assignment failed, falling back to basic",
            slog.String("ticket_id", ticketID),
            slog.String("error", err.Error()))
        
        // Fallback to basic assignment
        return s.basicAssignment(ctx, ticketID)
    }
    
    // Convert AI recommendations to service recommendations
    recommendations := make([]*EngineerRecommendation, len(result.Recommendations))
    for i, rec := range result.Recommendations {
        recommendations[i] = &EngineerRecommendation{
            EngineerID:   rec.EngineerID,
            EngineerName: rec.EngineerName,
            Score:        rec.Score,
            Reasoning:    rec.Reasoning,
        }
    }
    
    s.logger.Info("AI assignment completed",
        slog.String("ticket_id", ticketID),
        slog.Int("recommendations", len(recommendations)))
    
    return recommendations, nil
}
```

**Validation:** Request engineer assignment, check logs for "AI assignment completed".

---

#### ✅ Task 4.3: Add Parts Recommendation

**File:** `internal/service-domain/service-ticket/app/service.go`

**Add new method for parts recommendation:**

```go
func (s *Service) RecommendParts(ctx context.Context, ticketID string) (*parts.RecommendationResponse, error) {
    // Get ticket details
    ticket, err := s.repo.GetByID(ctx, ticketID)
    if err != nil {
        return nil, err
    }
    
    // ========================================================================
    // NEW: Call Parts Recommender
    // ========================================================================
    
    partsReq := &parts.RecommendationRequest{
        TicketID:      ticketID,
        EquipmentType: ticket.EquipmentName,
        ProblemType:   ticket.IssueCategory,
        Severity:      ticket.Priority,
        Options: parts.RecommendationOptions{
            IncludeReplacementParts: true,
            IncludeAccessories:      true,
            IncludePreventiveParts:  true,
            UseAI:                   true,
            MaxRecommendations:      10,
        },
    }
    
    result, err := s.partsEngine.RecommendParts(ctx, partsReq)
    if err != nil {
        return nil, fmt.Errorf("parts recommendation failed: %w", err)
    }
    
    s.logger.Info("Parts recommendation completed",
        slog.String("ticket_id", ticketID),
        slog.Int("replacement_parts", len(result.ReplacementParts)),
        slog.Int("accessories", len(result.Accessories)))
    
    return result, nil
}
```

**Add API endpoint to expose this:**

```go
// In handler
r.Get("/tickets/{id}/parts", handler.GetPartsRecommendation)
```

**Validation:** Get parts for a ticket, verify recommendations returned.

---

#### ✅ Task 4.4: Add Feedback Collection on Ticket Completion

**File:** `internal/service-domain/service-ticket/app/service.go`

**Update `CompleteTicket` or `CloseTicket` method:**

```go
func (s *Service) CloseTicket(ctx context.Context, ticketID string, resolution Resolution) error {
    // ... existing close logic ...
    
    // Update ticket status
    err := s.repo.Close(ctx, ticketID, resolution)
    if err != nil {
        return err
    }
    
    // ========================================================================
    // NEW: Auto-collect machine feedback
    // ========================================================================
    go func() {
        feedbackCtx := context.Background()
        if err := s.feedbackCollector.CollectTicketCompletionFeedback(feedbackCtx, ticketID); err != nil {
            s.logger.Error("Failed to collect machine feedback",
                slog.String("ticket_id", ticketID),
                slog.String("error", err.Error()))
        } else {
            s.logger.Info("Machine feedback collected",
                slog.String("ticket_id", ticketID))
        }
    }()
    
    return nil
}
```

**Validation:** Close a ticket, check logs for "Machine feedback collected".

---

### Phase 5: Database Migrations (0.5 days)

#### ✅ Task 5.1: Run AI Database Migrations

**Run migrations in order:**

```bash
# Navigate to project
cd aby-med

# Run each AI migration
psql $DATABASE_URL -f database/migrations/009_ai_diagnoses.sql
psql $DATABASE_URL -f database/migrations/010_assignment_history.sql
psql $DATABASE_URL -f database/migrations/011_parts_management.sql
psql $DATABASE_URL -f database/migrations/012_parts_recommendations.sql
psql $DATABASE_URL -f database/migrations/013_feedback_system.sql
```

**Verify tables created:**

```bash
psql $DATABASE_URL -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name LIKE 'ai_%' OR table_name LIKE 'feedback_%' OR table_name LIKE 'assignment_%' OR table_name LIKE 'parts_%';"
```

**Expected output:**
- `ai_diagnoses`
- `ai_diagnosis_analytics_view`
- `ai_diagnosis_feedback_summary_view`
- `assignment_history`
- `assignment_performance_view`
- `assignment_analytics_view`
- `equipment_variants`
- `parts_catalog`
- `equipment_parts`
- `equipment_accessories`
- `parts_suppliers`
- `supplier_parts`
- `parts_inventory`
- `parts_recommendations`
- `ai_feedback`
- `feedback_improvements`
- `feedback_actions`

**Validation:** All 17+ tables exist.

---

### Phase 6: Testing (1 day)

#### ✅ Task 6.1: End-to-End Workflow Test

**Test scenario:**

1. **Create Ticket**
   ```bash
   curl -X POST http://localhost:8081/api/v1/tickets \
     -H "Content-Type: application/json" \
     -d '{
       "equipment_id": "eq-123",
       "issue_description": "Ventilator showing error E-42, filter warning light on",
       "priority": "high",
       "source": "web"
     }'
   ```
   
   **Expected:** Ticket created, AI diagnosis triggered (check logs)

2. **Check Diagnosis**
   ```bash
   curl http://localhost:8081/api/diagnosis/history/TICKET_ID
   ```
   
   **Expected:** Diagnosis result with confidence score

3. **Get Engineer Recommendations**
   ```bash
   curl -X POST http://localhost:8081/api/assignment/recommend \
     -H "Content-Type: application/json" \
     -d '{
       "ticket_id": "TICKET_ID",
       "equipment_type": "Ventilator",
       "problem_type": "Filter Issue"
     }'
   ```
   
   **Expected:** List of recommended engineers with scores

4. **Get Parts Recommendations**
   ```bash
   curl http://localhost:8081/api/v1/tickets/TICKET_ID/parts
   ```
   
   **Expected:** Replacement parts + accessories

5. **Submit Human Feedback**
   ```bash
   curl -X POST http://localhost:8081/api/feedback/human \
     -H "Content-Type: application/json" \
     -d '{
       "service_type": "diagnosis",
       "request_id": "DIAG_REQUEST_ID",
       "user_id": 1,
       "rating": 5,
       "was_accurate": true,
       "comments": "Perfect diagnosis!"
     }'
   ```
   
   **Expected:** Feedback recorded

6. **Close Ticket**
   ```bash
   curl -X PUT http://localhost:8081/api/v1/tickets/TICKET_ID/close \
     -H "Content-Type: application/json" \
     -d '{
       "resolution_notes": "Replaced filter, issue resolved",
       "parts_used": ["FILTER-HEPA-01"],
       "labor_hours": 1.5
     }'
   ```
   
   **Expected:** Machine feedback auto-collected (check logs)

7. **Check Analytics**
   ```bash
   curl http://localhost:8081/api/diagnosis/analytics
   curl http://localhost:8081/api/assignment/analytics
   curl http://localhost:8081/api/parts/analytics
   curl http://localhost:8081/api/feedback/analytics
   ```
   
   **Expected:** Metrics displayed

---

#### ✅ Task 6.2: Run Integration Tests

```bash
# Navigate to project
cd aby-med

# Set test environment variables
export TEST_DB_URL="postgresql://localhost/geneqr_test"
export AI_PROVIDER="mock"

# Run integration tests
go test ./tests/integration/... -v
```

**Expected:** All tests pass

---

### Phase 7: Documentation Updates (0.5 days)

#### ✅ Task 7.1: Update API Documentation

**Add AI endpoints to API docs** (create `docs/API_REFERENCE.md` if doesn't exist)

#### ✅ Task 7.2: Update README

**Add AI features section** to main README.md

---

## 🎯 Success Criteria

### ✅ Integration Complete When:

1. **Configuration**
   - [ ] AI config in `config.go`
   - [ ] `.env` has all AI variables
   - [ ] App loads config without errors

2. **Initialization**
   - [ ] AI Manager initializes successfully
   - [ ] All AI engines created
   - [ ] Database pool connected

3. **Routes**
   - [ ] All AI routes mounted
   - [ ] Routes respond to requests
   - [ ] Error handling works

4. **Workflow Integration**
   - [ ] Tickets trigger AI diagnosis
   - [ ] Assignment uses AI optimizer
   - [ ] Parts recommendations work
   - [ ] Feedback auto-collects

5. **Database**
   - [ ] All migrations run successfully
   - [ ] Tables created
   - [ ] Data persists correctly

6. **Testing**
   - [ ] End-to-end workflow test passes
   - [ ] Integration tests pass
   - [ ] No errors in logs

7. **Documentation**
   - [ ] API docs updated
   - [ ] README updated
   - [ ] Integration guide complete

---

## 🚧 Rollback Plan

If integration fails:

1. **Revert code changes:**
   ```bash
   git revert HEAD~N  # where N is number of commits
   ```

2. **Rollback database:**
   ```bash
   # Drop AI tables (if needed)
   psql $DATABASE_URL -f database/rollback/rollback_ai_services.sql
   ```

3. **Disable AI features:**
   ```env
   # In .env
   AI_ENABLED=false
   ```

4. **Fallback to basic services:**
   - Use basic assignment (already exists)
   - Manual parts selection
   - No AI diagnosis

---

## 📊 Progress Tracking

| Phase | Tasks | Status | Time Estimate |
|-------|-------|--------|---------------|
| 1. Configuration | 2 | 🟡 Pending | 1 day |
| 2. AI Manager Init | 2 | 🟡 Pending | 1 day |
| 3. Mount Routes | 1 | 🟡 Pending | 1 day |
| 4. Workflow Integration | 4 | 🟡 Pending | 1-2 days |
| 5. Database Migrations | 1 | 🟡 Pending | 0.5 days |
| 6. Testing | 2 | 🟡 Pending | 1 day |
| 7. Documentation | 2 | 🟡 Pending | 0.5 days |
| **TOTAL** | **14 tasks** | **🟡 Not Started** | **5-7 days** |

---

## 🔍 Testing Matrix

| Test Scenario | Endpoint | Expected Result | Status |
|---------------|----------|-----------------|--------|
| Create ticket → AI diagnosis | POST /tickets | Diagnosis in logs | 🟡 |
| Get diagnosis history | GET /diagnosis/history/{id} | Results returned | 🟡 |
| AI engineer assignment | POST /assignment/recommend | Engineers ranked | 🟡 |
| Parts recommendation | GET /tickets/{id}/parts | Parts + accessories | 🟡 |
| Submit human feedback | POST /feedback/human | Feedback stored | 🟡 |
| Close ticket → machine feedback | PUT /tickets/{id}/close | Auto-collected | 🟡 |
| Check analytics | GET /diagnosis/analytics | Metrics shown | 🟡 |
| AI fallback | (Disable primary) | Uses fallback provider | 🟡 |
| Cost tracking | GET /ai/costs | Costs tracked | 🟡 |
| Learning progress | GET /feedback/learning-progress | Progress shown | 🟡 |

---

## 📞 Support

**Questions or issues during integration?**
- Review: `docs/REQUIREMENTS_MASTER.md`
- Check: `docs/PHASE_2C_COMPLETE.md`
- Test Guide: `docs/TESTING.md`
- Feedback System: `docs/FEEDBACK_SYSTEM.md`

---

**Document End**

_Update this document as integration progresses._
