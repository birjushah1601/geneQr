# Phase 2C: AI Services Layer - Implementation Plan

**Status:** 🚀 IN PROGRESS  
**Start Date:** 2025-11-16  
**Estimated Duration:** 2-3 weeks  
**Total Tickets:** 6  
**Dependencies:** Phase 2B Complete ✅  

---

## 📋 Executive Summary

Phase 2C implements the **Go services layer** that brings the Phase 2B database foundation to life with intelligent AI-powered capabilities. This phase creates:

- **Multi-provider AI abstraction** (OpenAI, Anthropic, extensible)
- **AI-powered diagnosis** from ticket descriptions and attachments
- **Intelligent engineer assignment** considering expertise, availability, location, workload
- **Context-aware parts recommendations** (ICU vs General Ward)
- **Feedback loop** for continuous AI learning
- **Complete integration tests** validating end-to-end workflows

---

## 🎯 Tickets Overview

### T2C.1: AI Service Foundation ✅ (3-4 days)
**Goal:** Create base AI client infrastructure with provider abstraction

**Deliverables:**
1. **AI Provider Interface** - Common interface for all AI providers
2. **OpenAI Client** - GPT-4o and GPT-4-Turbo implementation
3. **Anthropic Client** - Claude 3.5 Sonnet and Claude 3 Opus implementation
4. **Provider Manager** - Multi-provider orchestration with fallback
5. **Configuration** - API keys, rate limits, model selection
6. **Cost Tracking** - Token usage and cost calculation
7. **Retry Logic** - Exponential backoff with circuit breaker
8. **Health Checks** - Provider availability monitoring

**Key Features:**
- Provider abstraction (add new providers easily)
- Automatic fallback on provider failure
- Token usage tracking per request
- Cost calculation per provider/model
- Circuit breaker for failing providers
- Streaming support for real-time responses
- Function calling support

**Files to Create:**
- `internal/ai/provider.go` - Provider interface
- `internal/ai/openai/client.go` - OpenAI implementation
- `internal/ai/anthropic/client.go` - Anthropic implementation
- `internal/ai/manager.go` - Provider manager with fallback
- `internal/ai/config.go` - Configuration
- `internal/ai/models.go` - Common models and types
- `internal/ai/cost_tracker.go` - Cost tracking
- `pkg/ai/errors.go` - AI-specific errors

---

### T2C.2: Diagnosis Engine (3-4 days)
**Goal:** AI-powered ticket diagnosis from descriptions and attachments

**Deliverables:**
1. **Diagnosis Service** - Main service orchestrating AI diagnosis
2. **Prompt Engineering** - Optimized prompts for medical equipment diagnosis
3. **Vision Analysis** - Image/video analysis using vision models
4. **Context Enrichment** - Historical data and equipment specs
5. **Confidence Scoring** - 0-100% confidence calculation
6. **Issue Categorization** - Critical/High/Medium/Low severity
7. **Parts Identification** - Recommend parts based on diagnosis
8. **Support Level Detection** - Determine L1/L2/L3 requirement
9. **Database Persistence** - Save diagnosis results to DB

**Key Features:**
- Analyze ticket description for issues
- Vision AI for attachment analysis (images/videos)
- Enrich with equipment context (model, history)
- Enrich with historical similar tickets
- Generate diagnosis summary with confidence
- Identify root cause
- Recommend troubleshooting steps
- Recommend required parts
- Determine required support level
- Store complete results in `ai_diagnosis_results` table

**Files to Create:**
- `internal/ai-service/diagnosis/service.go` - Diagnosis orchestration
- `internal/ai-service/diagnosis/prompts.go` - Prompt templates
- `internal/ai-service/diagnosis/vision.go` - Vision analysis
- `internal/ai-service/diagnosis/context.go` - Context enrichment
- `internal/ai-service/diagnosis/repository.go` - DB persistence
- `internal/ai-service/diagnosis/models.go` - Domain models

**API Endpoints:**
- `POST /api/v1/ai/diagnose` - Diagnose ticket
- `GET /api/v1/ai/diagnosis/{ticket_id}` - Get diagnosis results
- `POST /api/v1/ai/diagnosis/{id}/validate` - Human validation

---

### T2C.3: Assignment Optimizer (3-4 days)
**Goal:** AI-based intelligent engineer assignment

**Deliverables:**
1. **Assignment Service** - Intelligent engineer matching
2. **Scoring Algorithm** - Multi-factor scoring (expertise, availability, location, performance, workload)
3. **Constraint Checking** - Availability, certifications, support level
4. **Ranking System** - Rank engineers 1st, 2nd, 3rd choice
5. **Reasoning Generation** - Explain why engineer recommended
6. **Database Persistence** - Save recommendations to DB

**Key Features:**
- Query qualified engineers from `engineer_equipment_expertise`
- Check availability and current workload
- Calculate distance from customer location
- Score based on:
  - Expertise level (L1/L2/L3)
  - Certifications (required vs actual)
  - Location proximity
  - Performance history
  - Current workload
- Generate reasoning (strengths/concerns)
- Rank top 3-5 engineers
- Store in `ai_engineer_recommendations` table

**Files to Create:**
- `internal/ai-service/assignment/service.go` - Assignment orchestration
- `internal/ai-service/assignment/scoring.go` - Scoring algorithms
- `internal/ai-service/assignment/constraints.go` - Constraint validation
- `internal/ai-service/assignment/repository.go` - DB persistence
- `internal/ai-service/assignment/models.go` - Domain models

**API Endpoints:**
- `POST /api/v1/ai/recommend-engineers` - Get engineer recommendations
- `GET /api/v1/ai/engineer-recommendations/{ticket_id}` - Get saved recommendations
- `POST /api/v1/ai/engineer-recommendations/{id}/select` - Mark engineer selected

---

### T2C.4: Parts Recommender (2-3 days)
**Goal:** Context-aware parts recommendations

**Deliverables:**
1. **Parts Service** - Intelligent parts matching
2. **Context Analysis** - ICU vs General Ward detection
3. **Compatibility Checking** - Alternative parts support
4. **Ranking System** - Rank parts by confidence
5. **Pricing Integration** - Include cost estimates
6. **Database Persistence** - Save recommendations to DB

**Key Features:**
- Extract parts from diagnosis results
- Detect installation context (ICU/Ward/OR/Lab)
- Query context-specific parts from `equipment_parts_context`
- Find alternatives using `equipment_compatibility`
- Check OEM requirements from `manufacturer_service_config`
- Rank parts by:
  - Confidence score
  - Context appropriateness
  - OEM requirement
  - Availability
  - Lead time
- Include pricing and lead time
- Store in `ai_parts_recommendations` table

**Files to Create:**
- `internal/ai-service/parts/service.go` - Parts recommendation
- `internal/ai-service/parts/context.go` - Context detection
- `internal/ai-service/parts/compatibility.go` - Compatibility checking
- `internal/ai-service/parts/repository.go` - DB persistence
- `internal/ai-service/parts/models.go` - Domain models

**API Endpoints:**
- `POST /api/v1/ai/recommend-parts` - Get parts recommendations
- `GET /api/v1/ai/parts-recommendations/{ticket_id}` - Get saved recommendations
- `POST /api/v1/ai/parts-recommendations/{id}/feedback` - Mark part used/not used

---

### T2C.5: Feedback Loop Manager (2-3 days)
**Goal:** Human validation and continuous learning

**Deliverables:**
1. **Feedback Service** - Collect and store feedback
2. **Validation API** - Human validation of AI results
3. **Accuracy Tracking** - Calculate accuracy rates
4. **Training Data Export** - Export for model fine-tuning
5. **Analytics Integration** - Feed into analytics views

**Key Features:**
- Validate diagnosis accuracy (was AI correct?)
- Validate engineer recommendations (was selection good?)
- Validate parts recommendations (were parts correct?)
- Track feedback from:
  - Engineers (field experience)
  - Managers (oversight)
  - Customers (satisfaction)
- Calculate accuracy rates per provider/model
- Store in `ai_feedback` table
- Mark training priority (high/medium/low)
- Export training data for fine-tuning

**Files to Create:**
- `internal/ai-service/feedback/service.go` - Feedback orchestration
- `internal/ai-service/feedback/validation.go` - Validation logic
- `internal/ai-service/feedback/analytics.go` - Accuracy calculation
- `internal/ai-service/feedback/repository.go` - DB persistence
- `internal/ai-service/feedback/models.go` - Domain models

**API Endpoints:**
- `POST /api/v1/ai/feedback/diagnosis` - Provide diagnosis feedback
- `POST /api/v1/ai/feedback/engineer-recommendation` - Provide assignment feedback
- `POST /api/v1/ai/feedback/parts-recommendation` - Provide parts feedback
- `GET /api/v1/ai/feedback/analytics` - Get accuracy metrics
- `GET /api/v1/ai/feedback/training-data` - Export training data

---

### T2C.6: Integration Tests (2 days)
**Goal:** End-to-end AI workflow testing

**Deliverables:**
1. **Unit Tests** - Test individual components
2. **Integration Tests** - Test complete workflows
3. **Mock Providers** - Test without calling real AI APIs
4. **Performance Tests** - Validate response times
5. **Cost Tests** - Validate cost tracking
6. **Error Handling Tests** - Test failure scenarios

**Test Scenarios:**
1. **Complete Diagnosis Flow:**
   - Create ticket → AI diagnoses → Parts recommended → Engineer assigned → Feedback collected

2. **Provider Fallback:**
   - Primary provider fails → Automatically falls back to secondary → Success

3. **Context Awareness:**
   - ICU ventilator → High-flow tubes recommended
   - General Ward ventilator → Standard tubes recommended

4. **Engineer Selection:**
   - L1 remote engineer for diagnosis
   - L2 onsite engineer for repair
   - L3 specialist for complex issues

5. **Accuracy Tracking:**
   - Correct diagnosis → Accuracy increases
   - Incorrect diagnosis → Accuracy tracked → Feedback stored

**Files to Create:**
- `internal/ai-service/diagnosis/service_test.go`
- `internal/ai-service/assignment/service_test.go`
- `internal/ai-service/parts/service_test.go`
- `internal/ai-service/feedback/service_test.go`
- `tests/integration/ai_workflow_test.go`
- `tests/mocks/ai_provider_mock.go`

---

## 🏗️ Architecture Overview

### Package Structure
```
internal/
├── ai/                          # AI provider abstraction
│   ├── provider.go              # Provider interface
│   ├── manager.go               # Provider manager with fallback
│   ├── config.go                # Configuration
│   ├── models.go                # Common models
│   ├── cost_tracker.go          # Cost tracking
│   ├── openai/                  # OpenAI implementation
│   │   ├── client.go
│   │   ├── models.go
│   │   └── vision.go
│   └── anthropic/               # Anthropic implementation
│       ├── client.go
│       ├── models.go
│       └── vision.go
├── ai-service/                  # AI services layer
│   ├── diagnosis/               # Diagnosis engine
│   │   ├── service.go
│   │   ├── prompts.go
│   │   ├── vision.go
│   │   ├── context.go
│   │   ├── repository.go
│   │   └── models.go
│   ├── assignment/              # Assignment optimizer
│   │   ├── service.go
│   │   ├── scoring.go
│   │   ├── constraints.go
│   │   ├── repository.go
│   │   └── models.go
│   ├── parts/                   # Parts recommender
│   │   ├── service.go
│   │   ├── context.go
│   │   ├── compatibility.go
│   │   ├── repository.go
│   │   └── models.go
│   └── feedback/                # Feedback loop
│       ├── service.go
│       ├── validation.go
│       ├── analytics.go
│       ├── repository.go
│       └── models.go
└── api/                         # HTTP handlers
    └── v1/
        └── ai/
            ├── diagnosis_handler.go
            ├── assignment_handler.go
            ├── parts_handler.go
            └── feedback_handler.go
```

### Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                         Ticket Created                          │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
                  ┌──────────────────────┐
                  │  Diagnosis Engine    │
                  │  (T2C.2)             │
                  └──────────┬───────────┘
                             │
                ┌────────────┴────────────┐
                │                         │
                ▼                         ▼
    ┌─────────────────────┐   ┌──────────────────────┐
    │  Assignment         │   │  Parts Recommender   │
    │  Optimizer (T2C.3)  │   │  (T2C.4)             │
    └──────────┬──────────┘   └──────────┬───────────┘
               │                          │
               │                          │
               ▼                          ▼
    ┌──────────────────────────────────────────────┐
    │         Engineer Assigned                    │
    │         Parts Ordered                        │
    └──────────────────┬───────────────────────────┘
                       │
                       ▼
            ┌──────────────────────┐
            │  Ticket Resolved     │
            └──────────┬───────────┘
                       │
                       ▼
            ┌──────────────────────┐
            │  Feedback Loop       │
            │  Manager (T2C.5)     │
            └──────────────────────┘
                       │
                       ▼
            ┌──────────────────────┐
            │  AI Accuracy         │
            │  Improves            │
            └──────────────────────┘
```

### AI Provider Flow

```
┌────────────────────────────────────────────────────────────┐
│                    AI Request                              │
└──────────────────────┬─────────────────────────────────────┘
                       │
                       ▼
            ┌──────────────────────┐
            │  Provider Manager    │
            │  (T2C.1)             │
            └──────────┬───────────┘
                       │
       ┌───────────────┼───────────────┐
       │               │               │
       ▼               ▼               ▼
┌──────────┐   ┌──────────┐   ┌──────────┐
│ OpenAI   │   │Anthropic │   │ Custom   │
│(Primary) │   │(Fallback)│   │ Provider │
└────┬─────┘   └────┬─────┘   └────┬─────┘
     │              │              │
     │   Success    │   Fallback   │
     ├──────────────┼──────────────┤
     │              │              │
     ▼              ▼              ▼
┌────────────────────────────────────────┐
│         Response + Metadata            │
│  - Tokens used                         │
│  - Cost                                │
│  - Latency                             │
│  - Provider/Model used                 │
└────────────────────────────────────────┘
```

---

## 🔑 Key Technologies

### AI Providers
- **OpenAI SDK:** `github.com/sashabaranov/go-openai`
- **Anthropic SDK:** Custom implementation (official Go SDK limited)

### Database
- **PostgreSQL Driver:** `github.com/lib/pq` or `pgx`
- **SQL Builder:** `github.com/Masterminds/squirrel` (optional)

### Configuration
- **Environment Variables:** `github.com/joho/godotenv`
- **Config Management:** `github.com/spf13/viper`

### HTTP Framework
- **Gin:** `github.com/gin-gonic/gin` (existing)

### Testing
- **Testify:** `github.com/stretchr/testify`
- **Mocking:** `github.com/stretchr/testify/mock`

### Logging
- **Zap:** `go.uber.org/zap` (structured logging)

---

## 📊 Success Criteria

### T2C.1: AI Service Foundation ✅
- [ ] OpenAI client successfully calls GPT-4o and GPT-4-Turbo
- [ ] Anthropic client successfully calls Claude 3.5 Sonnet and Claude 3 Opus
- [ ] Provider manager automatically falls back on failure
- [ ] Token usage tracked accurately
- [ ] Cost calculated correctly per provider/model
- [ ] Circuit breaker prevents cascading failures
- [ ] Configuration loads from environment variables
- [ ] Health checks detect provider availability

### T2C.2: Diagnosis Engine ✅
- [ ] Diagnosis generated from ticket description
- [ ] Vision AI analyzes uploaded images
- [ ] Confidence score calculated (0-100%)
- [ ] Root cause identified
- [ ] Parts recommended based on diagnosis
- [ ] Support level (L1/L2/L3) determined
- [ ] Results saved to `ai_diagnosis_results` table
- [ ] API endpoints functional and documented

### T2C.3: Assignment Optimizer ✅
- [ ] Qualified engineers identified based on expertise
- [ ] Multi-factor scoring works (expertise, availability, location, performance, workload)
- [ ] Top 3-5 engineers ranked
- [ ] Reasoning generated for each recommendation
- [ ] Results saved to `ai_engineer_recommendations` table
- [ ] API endpoints functional and documented

### T2C.4: Parts Recommender ✅
- [ ] Installation context detected (ICU/Ward/OR/Lab)
- [ ] Context-specific parts recommended
- [ ] Alternative parts identified via compatibility matrix
- [ ] OEM requirements respected
- [ ] Parts ranked by confidence and appropriateness
- [ ] Results saved to `ai_parts_recommendations` table
- [ ] API endpoints functional and documented

### T2C.5: Feedback Loop Manager ✅
- [ ] Human validation captured for all AI results
- [ ] Accuracy rates calculated per provider/model
- [ ] Feedback stored in `ai_feedback` table
- [ ] Training priority assigned
- [ ] Training data exportable for fine-tuning
- [ ] API endpoints functional and documented

### T2C.6: Integration Tests ✅
- [ ] All unit tests pass (>80% coverage)
- [ ] Integration tests validate end-to-end workflows
- [ ] Mock providers enable testing without API calls
- [ ] Performance tests validate <500ms response times
- [ ] Cost tracking validated
- [ ] Error handling tested (provider failures, timeouts, invalid inputs)

---

## 🎯 Timeline & Dependencies

| Ticket | Duration | Start After | End Result |
|--------|----------|-------------|------------|
| T2C.1  | 3-4 days | Phase 2B    | AI provider abstraction working |
| T2C.2  | 3-4 days | T2C.1       | Diagnosis engine functional |
| T2C.3  | 3-4 days | T2C.1       | Assignment optimizer functional |
| T2C.4  | 2-3 days | T2C.2       | Parts recommender functional |
| T2C.5  | 2-3 days | T2C.2, T2C.3, T2C.4 | Feedback loop functional |
| T2C.6  | 2 days   | All above   | Complete test coverage |

**Total Duration:** 15-20 days (3-4 weeks)

**Critical Path:** T2C.1 → T2C.2 → T2C.4 → T2C.5 → T2C.6

---

## 🚀 Getting Started

### Prerequisites
✅ Phase 2B database foundation complete  
✅ Go 1.21+ installed  
✅ PostgreSQL database accessible  
✅ OpenAI API key (environment variable: `OPENAI_API_KEY`)  
✅ Anthropic API key (environment variable: `ANTHROPIC_API_KEY`)  

### Environment Setup
```bash
# Install dependencies
go get github.com/sashabaranov/go-openai
go get github.com/lib/pq
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
go get go.uber.org/zap
go get github.com/stretchr/testify

# Set API keys
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export DATABASE_URL="postgres://user:pass@localhost:5432/geneqr"
```

### First Steps (T2C.1)
1. Create AI provider interface
2. Implement OpenAI client
3. Implement Anthropic client
4. Create provider manager
5. Add configuration
6. Add cost tracking
7. Write unit tests
8. Test with real API calls

---

## 📈 Business Value

### Faster Diagnosis
- **Before:** Manual diagnosis takes hours
- **After:** AI diagnosis in minutes
- **Impact:** 80-90% time reduction

### Better Engineer Matching
- **Before:** Assignment based on availability only
- **After:** AI considers expertise, location, performance, workload
- **Impact:** Optimal engineer for each ticket

### Context-Aware Parts
- **Before:** Generic parts recommendations
- **After:** ICU vs Ward context considered
- **Impact:** Reduced wrong orders, faster resolution

### Continuous Improvement
- **Before:** No feedback loop, static recommendations
- **After:** AI learns from corrections, improves over time
- **Impact:** Increasing accuracy over time

### Cost Optimization
- Complete cost tracking per AI provider/model
- Automatic fallback to cheaper providers when appropriate
- **Impact:** Optimized AI spend

---

## 🎉 Phase 2C Completion Criteria

**Phase 2C is COMPLETE when:**
✅ All 6 tickets implemented and tested  
✅ AI providers (OpenAI, Anthropic) working with fallback  
✅ Diagnosis engine generates accurate diagnoses  
✅ Assignment optimizer recommends optimal engineers  
✅ Parts recommender considers context (ICU/Ward)  
✅ Feedback loop captures human validation  
✅ Integration tests validate end-to-end workflows  
✅ API endpoints documented and functional  
✅ Code reviewed and merged  

**Ready for Phase 2D:** Application Services (Workflow orchestration, enhanced ticket management)

---

**Document Version:** 1.0  
**Created:** 2025-11-16  
**Author:** Droid (Factory AI Assistant)  
**Status:** 🚀 IN PROGRESS - Starting T2C.1
