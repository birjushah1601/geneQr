# Testing Guide - GeneQR AI Services

## Overview

This document describes the testing strategy for the GeneQR AI-enhanced service management system, including unit tests, integration tests, and end-to-end workflow validation.

---

## Test Structure

```
tests/
├── integration/          # Integration tests (cross-service)
│   ├── ai_workflow_test.go
│   └── ...
├── unit/                 # Unit tests (single component)
│   ├── diagnosis_test.go
│   ├── assignment_test.go
│   ├── parts_test.go
│   └── feedback_test.go
└── fixtures/             # Test data and mocks
    ├── test_tickets.json
    └── mock_ai_responses.json
```

---

## Integration Tests

### Purpose
Validate that all AI services work together correctly in real-world scenarios.

### Test Coverage

#### 1. Complete AI Workflow Test
**File:** `tests/integration/ai_workflow_test.go`  
**Function:** `TestCompleteAIWorkflow`

Tests the full workflow:
```
Ticket Created
    ↓
AI Diagnosis (Step 1)
    ↓
Engineer Assignment (Step 2)
    ↓
Parts Recommendation (Step 3)
    ↓
Human Feedback (Step 4)
    ↓
Machine Feedback Auto-Collection (Step 5)
    ↓
Feedback Analysis (Step 6)
    ↓
Learning Progress (Step 7)
```

**Expected Results:**
- ✅ Diagnosis completes with confidence score >0
- ✅ Assignment returns ranked engineer recommendations
- ✅ Parts includes replacement parts + accessories
- ✅ Human feedback is collected and stored
- ✅ Machine feedback is auto-collected on ticket closure
- ✅ Analysis detects patterns and generates improvements
- ✅ Learning progress is tracked

#### 2. AI Workflow with Corrections Test
**Function:** `TestAIWorkflowWithCorrections`

Tests the learning loop when AI is wrong:
- AI provides incorrect diagnosis
- Human corrects the error
- System detects pattern from corrections
- Improvement opportunity is generated

**Expected Results:**
- ✅ Negative feedback creates negative sentiment
- ✅ Corrections are stored properly
- ✅ Analyzer identifies common issues
- ✅ Improvement opportunities are generated

#### 3. Parallel AI Requests Test
**Function:** `TestParallelAIRequests`

Tests concurrent AI requests:
- Creates 5 tickets simultaneously
- Runs diagnosis for all in parallel
- Validates all complete successfully

**Expected Results:**
- ✅ All parallel requests succeed
- ✅ No race conditions or deadlocks
- ✅ Results are correct for each ticket

#### 4. Feedback Loop Improvement Cycle Test
**Function:** `TestFeedbackLoopImprovementCycle`

Tests the complete learning cycle:
- Collects multiple feedback entries with same issue
- Analyzer detects pattern
- Improvement opportunity is generated
- Change is applied and tested

**Expected Results:**
- ✅ Pattern detection works (3+ similar issues)
- ✅ Improvements are generated with correct type
- ✅ Learning actions can be applied

---

## Running Tests

### Prerequisites

```bash
# Install Go testing tools
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/require

# Setup test database
createdb geneqr_test
psql geneqr_test < database/schema.sql
psql geneqr_test < database/migrations/*.sql
```

### Run All Tests

```bash
# Run all tests with verbose output
go test ./... -v

# Run only integration tests
go test ./tests/integration/... -v

# Run specific test
go test ./tests/integration/... -v -run TestCompleteAIWorkflow

# Run with coverage
go test ./... -v -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Run Tests with Database

```bash
# Set test database connection
export TEST_DB_URL="postgresql://user:pass@localhost/geneqr_test"

# Run tests
go test ./tests/integration/... -v
```

### Run Tests with Mock AI

```bash
# Use mock AI provider for faster tests
export AI_PROVIDER="mock"
export MOCK_AI_RESPONSES="tests/fixtures/mock_ai_responses.json"

# Run tests
go test ./tests/integration/... -v
```

---

## Test Configuration

### Environment Variables

```bash
# Database
TEST_DB_URL="postgresql://localhost/geneqr_test"

# AI Provider
AI_PROVIDER="openai"           # or "mock" for testing
OPENAI_API_KEY="test_key"      # for integration testing with real AI
ANTHROPIC_API_KEY="test_key"

# Test Options
RUN_SLOW_TESTS="false"         # Skip slow AI API calls
PARALLEL_TESTS="4"             # Number of parallel test workers
```

### Test Database Setup

```sql
-- Create test database
CREATE DATABASE geneqr_test;

-- Run all migrations
\i database/schema.sql
\i database/migrations/001_initial_schema.sql
-- ... all migrations ...

-- Seed test data
INSERT INTO users (user_id, username, email, role) 
VALUES (1, 'test_engineer', 'test@example.com', 'field_engineer');

INSERT INTO users (user_id, username, email, role) 
VALUES (2, 'test_dispatcher', 'dispatch@example.com', 'dispatcher');

-- Add test equipment, engineers, parts, etc.
```

---

## Unit Tests

### Diagnosis Engine Tests

```go
// Test diagnosis accuracy
func TestDiagnosisAccuracy(t *testing.T)

// Test diagnosis with images
func TestDiagnosisWithVisionAnalysis(t *testing.T)

// Test similar tickets enrichment
func TestSimilarTicketsEnrichment(t *testing.T)

// Test error handling
func TestDiagnosisErrorHandling(t *testing.T)
```

### Assignment Engine Tests

```go
// Test multi-factor scoring
func TestMultiFactorScoring(t *testing.T)

// Test location weight impact
func TestLocationWeightImpact(t *testing.T)

// Test expertise matching
func TestExpertiseMatching(t *testing.T)

// Test workload balancing
func TestWorkloadBalancing(t *testing.T)
```

### Parts Engine Tests

```go
// Test diagnosis-based recommendations
func TestDiagnosisBasedRecommendations(t *testing.T)

// Test variant-specific accessories
func TestVariantSpecificAccessories(t *testing.T)

// Test preventive maintenance parts
func TestPreventiveMaintenanceParts(t *testing.T)

// Test upselling logic
func TestUpsellingLogic(t *testing.T)
```

### Feedback System Tests

```go
// Test human feedback collection
func TestHumanFeedbackCollection(t *testing.T)

// Test machine feedback collection
func TestMachineFeedbackCollection(t *testing.T)

// Test sentiment analysis
func TestSentimentAnalysis(t *testing.T)

// Test pattern detection
func TestPatternDetection(t *testing.T)

// Test improvement generation
func TestImprovementGeneration(t *testing.T)

// Test learning action application
func TestLearningActionApplication(t *testing.T)
```

---

## Continuous Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Run Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_DB: geneqr_test
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: |
        go mod download
        go install github.com/stretchr/testify
    
    - name: Setup test database
      run: |
        psql -h localhost -U postgres -d geneqr_test -f database/schema.sql
        for f in database/migrations/*.sql; do
          psql -h localhost -U postgres -d geneqr_test -f $f
        done
      env:
        PGPASSWORD: postgres
    
    - name: Run unit tests
      run: go test ./internal/... -v -cover
      env:
        AI_PROVIDER: mock
    
    - name: Run integration tests
      run: go test ./tests/integration/... -v
      env:
        TEST_DB_URL: postgresql://postgres:postgres@localhost/geneqr_test
        AI_PROVIDER: mock
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
```

---

## Performance Testing

### Load Test: Concurrent Diagnoses

```go
func BenchmarkConcurrentDiagnoses(b *testing.B) {
    db, _ := setupTestDB(b)
    aiManager := setupTestAIManager(b)
    diagnosisEngine := diagnosis.NewEngine(aiManager, db)
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            req := &diagnosis.DiagnosisRequest{
                TicketID: rand.Int63(),
                EquipmentType: "Ventilator",
                ProblemDescription: "Test issue",
                Options: diagnosis.DiagnosisOptions{UseAI: true},
            }
            
            _, err := diagnosisEngine.DiagnoseIssue(context.Background(), req)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

**Run:**
```bash
go test -bench=BenchmarkConcurrentDiagnoses -benchtime=30s
```

### Load Test: Feedback Processing

```go
func BenchmarkFeedbackProcessing(b *testing.B) {
    db, _ := setupTestDB(b)
    collector := feedback.NewCollector(db)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := &feedback.HumanFeedbackRequest{
            ServiceType: "diagnosis",
            RequestID: fmt.Sprintf("req_%d", i),
            UserID: 1,
            WasAccurate: true,
            Rating: intPtr(5),
        }
        
        _, err := collector.CollectHumanFeedback(context.Background(), req)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## Test Data Fixtures

### Mock AI Responses

```json
// tests/fixtures/mock_ai_responses.json
{
  "diagnosis": {
    "ventilator_filter": {
      "primary_diagnosis": "Filter Clogged",
      "confidence": 92.5,
      "reasoning": "Error code E-42 indicates filter issue...",
      "recommended_actions": [
        "Replace HEPA filter",
        "Check air pressure readings",
        "Inspect filter housing"
      ]
    }
  },
  "assignment": {
    "expert_recommendation": {
      "reasoning": "Engineer has 5+ years ventilator experience...",
      "adjustments": {
        "expertise_boost": 10,
        "location_boost": 5
      }
    }
  },
  "parts": {
    "refinement": {
      "additional_parts": [],
      "removed_parts": [],
      "reasoning": "All recommended parts are appropriate..."
    }
  }
}
```

### Test Tickets

```json
// tests/fixtures/test_tickets.json
{
  "tickets": [
    {
      "ticket_id": 1001,
      "equipment_type": "Ventilator",
      "problem_description": "Filter warning light on, error E-42",
      "priority": "High",
      "expected_diagnosis": "Filter Clogged",
      "expected_parts": ["HEPA Filter Assembly"]
    },
    {
      "ticket_id": 1002,
      "equipment_type": "X-Ray Machine",
      "problem_description": "Blurry images, poor quality",
      "priority": "Medium",
      "expected_diagnosis": "Detector Calibration",
      "expected_parts": ["Calibration Kit"]
    }
  ]
}
```

---

## Test Coverage Goals

| Component | Target Coverage | Current |
|-----------|----------------|---------|
| AI Manager | 90% | TBD |
| Diagnosis Engine | 85% | TBD |
| Assignment Engine | 85% | TBD |
| Parts Engine | 85% | TBD |
| Feedback System | 90% | TBD |
| API Handlers | 80% | TBD |
| **Overall** | **85%** | **TBD** |

---

## Troubleshooting Tests

### Common Issues

#### 1. Database Connection Errors

```bash
# Check if test database exists
psql -l | grep geneqr_test

# Recreate test database
dropdb geneqr_test
createdb geneqr_test
psql geneqr_test < database/schema.sql
```

#### 2. AI API Rate Limits

```bash
# Use mock AI provider
export AI_PROVIDER="mock"

# Or add delays between tests
export TEST_AI_DELAY="1000"  # milliseconds
```

#### 3. Flaky Tests

```bash
# Run specific test multiple times
go test -run TestCompleteAIWorkflow -count=10

# Increase timeouts
export TEST_TIMEOUT="60s"
```

#### 4. Missing Test Data

```bash
# Seed test database
psql geneqr_test < tests/fixtures/seed_test_data.sql
```

---

## Best Practices

### 1. Test Isolation
- Each test should be independent
- Use database transactions and rollback
- Clean up test data after each test

### 2. Mocking External Services
- Mock AI providers for unit tests
- Use test API keys for integration tests
- Never use production API keys in tests

### 3. Test Naming
- Use descriptive test names
- Follow pattern: `Test<Component><Scenario>`
- Example: `TestDiagnosisEngineWithImages`

### 4. Assertions
- Use meaningful assertion messages
- Check both success and error cases
- Validate all important fields in responses

### 5. Test Documentation
- Add comments explaining complex test scenarios
- Document expected behaviors
- Link to related tickets/issues

---

## Future Improvements

- [ ] Add E2E tests with real UI interactions
- [ ] Implement chaos testing for AI fallback
- [ ] Add performance regression tests
- [ ] Create test data generation tools
- [ ] Implement visual regression tests for UI
- [ ] Add mutation testing for code quality
- [ ] Create automated test reports dashboard

---

## Running Tests in Development

```bash
# Quick test during development (unit tests only)
go test ./internal/diagnosis/... -v -short

# Full test before commit
go test ./... -v -cover

# Integration test specific feature
go test ./tests/integration/... -v -run TestCompleteAIWorkflow

# Watch mode (with external tool)
# Install: go install github.com/cespare/reflex@latest
reflex -r '\.go$' -s -- go test ./internal/diagnosis/... -v
```

---

**For questions or issues:**
- GitHub Issues: https://github.com/birjushah1601/geneQr/issues
- Wiki: https://github.com/birjushah1601/geneQr/wiki/Testing
- Slack: #geneqr-testing
