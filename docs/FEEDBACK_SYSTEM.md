# AI Feedback and Continuous Learning System

## Overview

The GeneQR AI Feedback System collects feedback from **both human users and system outcomes** to continuously improve the accuracy and effectiveness of all AI services (Diagnosis, Assignment, Parts Recommendation).

This document explains:
1. How humans can provide feedback
2. How the system automatically learns from outcomes
3. How AI continuously improves over time
4. The complete feedback loop architecture

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                     FEEDBACK COLLECTION LAYER                         │
├────────────────────────────────┬─────────────────────────────────────┤
│     HUMAN FEEDBACK (Explicit)  │  MACHINE FEEDBACK (Implicit)        │
│                                │                                      │
│  • Engineers rating diagnosis  │  • Actual outcomes vs predictions   │
│  • Dispatchers evaluating      │  • Which parts were actually used   │
│    assignments                 │  • Assignment acceptance rates      │
│  • Technicians confirming      │  • Ticket resolution times          │
│    parts accuracy              │  • Customer satisfaction scores     │
│  • Manual corrections          │  • First-time fix rates             │
│  • Written comments            │  • Cost variances                   │
└────────────────────────────────┴─────────────────────────────────────┘
                                 ↓
┌──────────────────────────────────────────────────────────────────────┐
│                       FEEDBACK ANALYZER                               │
│                                                                        │
│  • Aggregates feedback from all sources                              │
│  • Identifies patterns and recurring issues                          │
│  • Detects correlations (e.g., "diagnosis always wrong for X")       │
│  • Calculates accuracy rates and trends                              │
│  • Generates improvement opportunities                                │
└────────────────────────────────────────────────────────────────────────┘
                                 ↓
┌──────────────────────────────────────────────────────────────────────┐
│                      LEARNING ENGINE                                  │
│                                                                        │
│  • Evaluates improvement opportunities                                │
│  • Applies changes (prompt tuning, weight adjustments)               │
│  • Tests changes in staging                                           │
│  • Measures impact (before/after metrics)                            │
│  • Auto-deploys successful improvements                               │
│  • Auto-rolls back unsuccessful changes                               │
└────────────────────────────────────────────────────────────────────────┘
                                 ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    IMPROVED AI SERVICES                               │
│                                                                        │
│  Diagnosis → Assignment → Parts → Better Outcomes → More Feedback    │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Human Feedback: How Users Provide Feedback

### 1.1 Diagnosis Feedback

**When:** After reviewing an AI diagnosis

**How:** 

```http
POST /api/feedback/human
Content-Type: application/json

{
  "service_type": "diagnosis",
  "request_id": "diag_req_123456",
  "ticket_id": 789,
  "user_id": 42,
  "user_role": "field_engineer",
  "rating": 4,
  "was_accurate": true,
  "comments": "Good diagnosis but missed the secondary issue with the humidity sensor",
  "corrections": {
    "additional_problems": ["Humidity sensor calibration error"]
  }
}
```

**UI Integration Points:**
- ✅/❌ buttons after diagnosis is shown
- 1-5 star rating widget
- Optional comment box for details
- "Suggest correction" button for experts

**Example UI:**
```
╔════════════════════════════════════════════════╗
║ AI Diagnosis: Filter Clogged (92% confidence) ║
║                                                 ║
║ Was this diagnosis accurate?                   ║
║  [ Yes, Accurate ✓ ]  [ No, Incorrect ✗ ]      ║
║                                                 ║
║ Rate this diagnosis: ⭐⭐⭐⭐☆                    ║
║                                                 ║
║ Additional comments (optional):                ║
║  ┌──────────────────────────────────────┐      ║
║  │                                      │      ║
║  └──────────────────────────────────────┘      ║
║                                                 ║
║          [ Submit Feedback ]                    ║
╚════════════════════════════════════════════════╝
```

### 1.2 Assignment Feedback

**When:** After engineer selection or ticket completion

**How:**

```http
POST /api/feedback/human
Content-Type: application/json

{
  "service_type": "assignment",
  "request_id": "assign_req_789012",
  "ticket_id": 789,
  "user_id": 15,
  "user_role": "dispatcher",
  "rating": 5,
  "was_accurate": true,
  "comments": "Perfect match! Engineer had exact expertise needed",
  "corrections": {}
}
```

**UI Integration Points:**
- Quick feedback after assignment acceptance
- Post-completion survey for dispatcher
- Engineer self-assessment of assignment fit

### 1.3 Parts Feedback

**When:** After parts are ordered/used

**How:**

```http
POST /api/feedback/human
Content-Type: application/json

{
  "service_type": "parts",
  "request_id": "parts_req_345678",
  "ticket_id": 789,
  "user_id": 42,
  "user_role": "field_engineer",
  "rating": 3,
  "was_accurate": false,
  "comments": "Recommended part was correct but missed a required gasket",
  "corrections": {
    "missing_parts": ["Gasket-234-A"]
  }
}
```

**UI Integration Points:**
- Checkbox list: "Which recommended parts did you use?"
- "Add missing part" button
- Quick rating after job completion

---

## 2. Machine Feedback: Automatic Learning from Outcomes

### 2.1 Automatic Collection on Ticket Closure

When a ticket is closed/resolved, the system **automatically** collects feedback by comparing AI predictions with actual outcomes.

**Triggered by:**
```http
POST /api/tickets/{ticketId}/auto-feedback
```

**What it collects:**

#### Diagnosis Outcomes
- Was the AI-predicted problem correct?
- Did the engineer confirm the diagnosis?
- Were there additional issues not caught by AI?

#### Assignment Outcomes
- Was the recommended engineer selected?
- If not, which engineer was assigned?
- Did the assignment result in successful resolution?
- How long did it take?

#### Parts Outcomes
- Which recommended parts were actually used?
- Were additional parts needed?
- Were accessories sold (upsell success)?
- Cost accuracy (estimated vs actual)

### 2.2 Performance Metrics Collected

```json
{
  "service_type": "diagnosis",
  "request_id": "diag_req_123456",
  "ticket_id": 789,
  "outcomes": {
    "diagnosis_matched_ai": true,
    "actual_problem": "Filter Clogged",
    "resolution_time_minutes": 45,
    "first_time_fix": true,
    "customer_satisfaction": 5,
    "actual_cost": 120.50,
    "estimated_cost": 115.00,
    "cost_variance_percent": 4.8
  }
}
```

---

## 3. How AI Learns and Improves

### 3.1 Feedback Analysis Process

The **Analyzer** runs daily (or triggered manually) to:

1. **Aggregate feedback** from the last 30 days
2. **Identify patterns:**
   - "Diagnosis fails for equipment model X"
   - "Parts recommendations missing gaskets in 15 cases"
   - "Assignment algorithm underweighting location"

3. **Generate improvement opportunities:**

```json
{
  "opportunity_id": "diagnosis_equipment_model_x_1234567",
  "title": "Improve diagnosis for Equipment Model X",
  "description": "AI diagnosis has only 65% accuracy for Equipment Model X (vs 90% overall). Need to improve model-specific diagnosis logic.",
  "impact_level": "high",
  "implementation_type": "prompt_tuning",
  "suggested_changes": {
    "add_instruction": "For Equipment Model X, pay special attention to humidity sensor and air filter issues which are more common.",
    "update_prompts": true
  },
  "supporting_data": [feedback_id_1, feedback_id_2, ...],
  "status": "pending"
}
```

### 3.2 Learning Actions: What the System Can Change

#### A. Prompt Tuning (Automatic)
**What:** Modifies AI prompts to address specific issues  
**When:** Pattern of similar errors detected (e.g., always missing a specific component)  
**Example:**
```
Before: "Analyze the equipment issue and provide diagnosis"

After: "Analyze the equipment issue and provide diagnosis. 
       Pay special attention to humidity sensors in Equipment Model X, 
       which fail frequently and are often missed."
```

#### B. Weight Adjustment (Automatic)
**What:** Adjusts scoring weights in recommendation algorithms  
**When:** One factor is consistently over/under-valued  
**Example:**
```
Assignment Scoring Weights:
Before: Expertise=35%, Location=15%, Performance=25%, Workload=15%, Availability=10%
After:  Expertise=35%, Location=20%, Performance=25%, Workload=10%, Availability=10%

Reason: Analysis showed location proximity had bigger impact on 
        resolution time than workload consideration.
```

#### C. Config Change (Automatic)
**What:** Updates thresholds, confidence levels, filtering rules  
**When:** Patterns suggest different optimal values  
**Example:**
```
Parts Recommendation Confidence Threshold:
Before: Only show parts with >80% confidence
After:  Only show parts with >75% confidence

Reason: 75-80% confidence parts were actually used 85% of the time, 
        so threshold was too strict.
```

#### D. Training Data (Manual Review Required)
**What:** Flags cases for model retraining  
**When:** Systemic issues that require model updates  
**Example:**
```
Issue: AI consistently misdiagnoses "sensor error" as "filter clogged"
Action: Generate training dataset from corrected cases
Status: Requires ML engineer review
```

### 3.3 Testing and Deployment Cycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. IMPROVEMENT IDENTIFIED                                    │
│    Pattern detected: "Parts recommendations missing gaskets" │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. CAPTURE BEFORE METRICS                                    │
│    - Accuracy Rate: 85%                                      │
│    - Avg Rating: 3.8/5                                       │
│    - Positive Sentiment: 78%                                 │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. APPLY CHANGE (Status: TESTING)                           │
│    Action: Update parts logic to include related gaskets    │
│    Applied At: 2025-11-17 10:00:00                          │
│    Applied By: system                                        │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. MONITOR FOR 7 DAYS                                        │
│    Collect feedback with the new logic...                   │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. EVALUATE ACTION                                           │
│    - After Accuracy Rate: 92% (+7%)                          │
│    - After Avg Rating: 4.3/5 (+0.5)                          │
│    - After Positive Sentiment: 88% (+10%)                    │
│    - Overall Improvement: +8.2%                              │
└────────────────────┬────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────────┐
│ 6. DECISION                                                  │
│    ✅ Improvement >= +5% → Deploy to Production              │
│    ⚠️  Improvement -5% to +5% → Continue Testing             │
│    ❌ Improvement < -5% → Rollback Changes                   │
└─────────────────────────────────────────────────────────────┘
```

**Auto-Deployment Criteria:**
- ✅ Accuracy improved by ≥5%
- ✅ No decrease in any metric >3%
- ✅ At least 7 days of testing data
- ✅ Minimum 20 feedback samples collected

**Auto-Rollback Triggers:**
- ❌ Accuracy decreased by ≥5%
- ❌ Customer satisfaction dropped significantly
- ❌ Error rate increased
- ❌ Negative feedback spike

---

## 4. API Reference

### 4.1 Submit Human Feedback

```http
POST /api/feedback/human
Content-Type: application/json

{
  "service_type": "diagnosis|assignment|parts",
  "request_id": "string",
  "ticket_id": 123,
  "user_id": 456,
  "user_role": "engineer|dispatcher|technician",
  "rating": 1-5,
  "was_accurate": boolean,
  "comments": "string",
  "corrections": {
    "field": "corrected_value"
  }
}
```

**Response:**
```json
{
  "success": true,
  "feedback_id": 789,
  "message": "Thank you for your feedback! It helps us improve our AI systems."
}
```

### 4.2 Auto-Collect Feedback (Internal)

```http
POST /api/tickets/{ticketId}/auto-feedback
```

Automatically collects machine feedback when ticket is closed.

### 4.3 Get Feedback for Request

```http
GET /api/feedback/{serviceType}/{requestId}
```

Returns all feedback (human + machine) for a specific AI request.

### 4.4 Get Analytics

```http
GET /api/feedback/analytics?service_type=diagnosis&days=30
```

**Response:**
```json
{
  "service_type": "diagnosis",
  "period": "30_days",
  "total_feedback": 150,
  "human_feedback_count": 45,
  "machine_feedback_count": 105,
  "positive_feedback": 120,
  "neutral_feedback": 20,
  "negative_feedback": 10,
  "avg_rating": 4.2,
  "accuracy_rate": 87.5,
  "common_issues": [
    {
      "issue_type": "incorrect_confidence",
      "frequency": 8,
      "severity": "medium"
    }
  ],
  "improvements": [
    {
      "opportunity_id": "...",
      "title": "Improve confidence calibration",
      "impact_level": "medium",
      "status": "pending"
    }
  ]
}
```

### 4.5 Get Learning Progress

```http
GET /api/feedback/learning-progress?service_type=diagnosis
```

**Response:**
```json
{
  "total_improvements_identified": 15,
  "pending_improvements": 5,
  "applied_improvements": 10,
  "total_actions": 10,
  "deployed_actions": 7,
  "rolled_back_actions": 1,
  "success_rate_percent": 70.0,
  "recent_actions": [
    {
      "action_id": "action_123",
      "action_type": "prompt_update",
      "status": "deployed",
      "applied_at": "2025-11-10T10:00:00Z",
      "result_notes": "Improved performance by 8.2%. Deployed to production."
    }
  ]
}
```

### 4.6 Apply Improvement (Admin)

```http
POST /api/feedback/improvements/{opportunityId}/apply
Content-Type: application/json

{
  "applied_by": "admin_user_123"
}
```

Manually trigger an improvement to be applied and tested.

### 4.7 Evaluate Action (System/Admin)

```http
POST /api/feedback/actions/{actionId}/evaluate
```

Evaluates a testing action and decides whether to deploy, continue testing, or rollback.

---

## 5. Integration Guide

### 5.1 Frontend Integration

**Step 1: Add Feedback Widgets to AI Response UI**

```javascript
// After showing AI diagnosis
<DiagnosisResult diagnosis={aiDiagnosis} />

<FeedbackWidget
  serviceType="diagnosis"
  requestId={aiDiagnosis.request_id}
  ticketId={ticketId}
  userId={currentUser.id}
  userRole={currentUser.role}
  onSubmit={handleFeedbackSubmit}
/>
```

**Step 2: Implement Feedback Component**

```jsx
function FeedbackWidget({ serviceType, requestId, ticketId, userId, userRole }) {
  const [rating, setRating] = useState(null);
  const [wasAccurate, setWasAccurate] = useState(null);
  const [comments, setComments] = useState('');

  const handleSubmit = async () => {
    const response = await fetch('/api/feedback/human', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        service_type: serviceType,
        request_id: requestId,
        ticket_id: ticketId,
        user_id: userId,
        user_role: userRole,
        rating,
        was_accurate: wasAccurate,
        comments
      })
    });

    if (response.ok) {
      showToast('Thank you for your feedback!');
    }
  };

  return (
    <div className="feedback-widget">
      <h4>Was this {serviceType} accurate?</h4>
      
      <div className="accuracy-buttons">
        <button onClick={() => setWasAccurate(true)}>
          ✓ Yes, Accurate
        </button>
        <button onClick={() => setWasAccurate(false)}>
          ✗ No, Incorrect
        </button>
      </div>

      <StarRating value={rating} onChange={setRating} />

      <textarea
        placeholder="Additional comments (optional)"
        value={comments}
        onChange={(e) => setComments(e.target.value)}
      />

      <button onClick={handleSubmit}>Submit Feedback</button>
    </div>
  );
}
```

### 5.2 Backend Integration

**Auto-collect feedback on ticket closure:**

```go
// When ticket is marked as resolved
func (s *TicketService) ResolveTicket(ctx context.Context, ticketID int64) error {
    // ... update ticket status ...

    // Auto-collect machine feedback
    err = s.feedbackCollector.CollectTicketCompletionFeedback(ctx, ticketID)
    if err != nil {
        // Log but don't fail ticket resolution
        log.Printf("Failed to collect feedback for ticket %d: %v", ticketID, err)
    }

    return nil
}
```

---

## 6. Monitoring and Dashboards

### 6.1 Feedback Dashboard Metrics

**Key Metrics to Display:**

1. **Feedback Volume**
   - Total feedback received (human + machine)
   - Feedback rate (% of AI requests with feedback)
   - Trend over time

2. **Accuracy Metrics**
   - Overall accuracy rate
   - By service type (diagnosis/assignment/parts)
   - Trend: improving/stable/declining

3. **Sentiment Analysis**
   - Positive/Neutral/Negative breakdown
   - Average rating
   - Sentiment trend

4. **Learning Progress**
   - Improvements identified
   - Actions applied
   - Success rate
   - Active testing actions

5. **Top Issues**
   - Most common problems
   - Severity distribution
   - Time to resolution

### 6.2 Admin Dashboard Example

```
╔═══════════════════════════════════════════════════════════╗
║           AI FEEDBACK & LEARNING DASHBOARD                 ║
╠═══════════════════════════════════════════════════════════╣
║                                                            ║
║  Last 30 Days Overview                                     ║
║  ┌──────────────┬──────────────┬──────────────┐          ║
║  │  Diagnosis   │  Assignment  │    Parts     │          ║
║  ├──────────────┼──────────────┼──────────────┤          ║
║  │ Accuracy     │ Accuracy     │ Accuracy     │          ║
║  │   87.5% ▲    │   92.0% ▲    │   85.0% ►    │          ║
║  │              │              │              │          ║
║  │ Avg Rating   │ Avg Rating   │ Avg Rating   │          ║
║  │   4.2 / 5    │   4.5 / 5    │   4.0 / 5    │          ║
║  │              │              │              │          ║
║  │ Feedback     │ Feedback     │ Feedback     │          ║
║  │  150 total   │  98 total    │  120 total   │          ║
║  └──────────────┴──────────────┴──────────────┘          ║
║                                                            ║
║  Active Learning Actions (Testing)                         ║
║  • Improve diagnosis for Equipment Model X (Day 3/7)       ║
║  • Adjust location weight in assignment (Day 5/7)          ║
║                                                            ║
║  Recent Deployments (Last 7 Days)                          ║
║  ✅ Parts logic: Include related gaskets (+8.2%)           ║
║  ✅ Diagnosis prompts: Focus on humidity sensors (+6.5%)   ║
║                                                            ║
║  Top Issues Requiring Attention                            ║
║  ⚠️  Diagnosis struggling with Equipment Model Y (12 cases)║
║  ⚠️  Parts missing tool accessories (8 cases)              ║
║                                                            ║
╚═══════════════════════════════════════════════════════════╝
```

---

## 7. Best Practices

### 7.1 For Users Providing Feedback

✅ **DO:**
- Provide feedback as soon as possible (while details are fresh)
- Be specific in comments ("missed humidity sensor" vs "wrong")
- Rate honestly - both positive and negative feedback are valuable
- Use corrections field to suggest what should have been recommended

❌ **DON'T:**
- Skip feedback when AI is wrong (negative feedback is most valuable!)
- Give low ratings without explanation
- Provide generic comments like "bad" or "good"

### 7.2 For System Administrators

✅ **DO:**
- Review improvement opportunities weekly
- Monitor learning progress dashboard
- Investigate rolled-back actions to understand failures
- Manually review high-impact improvements before deployment

❌ **DON'T:**
- Ignore pending improvements for too long
- Override system rollbacks without investigation
- Disable automatic feedback collection
- Skip evaluation of testing actions

### 7.3 For Developers

✅ **DO:**
- Call auto-collect feedback on every ticket closure
- Add feedback widgets to all AI response UIs
- Log feedback submission failures for debugging
- Test feedback flows in staging environment

❌ **DON'T:**
- Block ticket operations on feedback failures
- Expose internal improvement IDs to end users
- Modify learning thresholds without data analysis

---

## 8. Troubleshooting

### Q: Feedback submission failing?
**A:** Check:
1. User has valid user_id and authentication
2. request_id matches an existing AI request
3. service_type is one of: diagnosis, assignment, parts
4. Required fields are not empty

### Q: Auto-collection not working?
**A:** Verify:
1. Ticket has AI requests associated with it
2. Ticket is in "resolved" or "closed" status
3. Database triggers are functioning
4. Check logs for specific errors

### Q: No improvements being generated?
**A:** Check:
1. Sufficient feedback volume (need at least 20 samples)
2. Patterns are detectable (need 3+ similar issues)
3. Analyzer is running (scheduled daily)
4. Review analyzer logs for errors

### Q: Action stuck in "testing" status?
**A:** Requirements:
1. Must have at least 7 days of data
2. Need minimum 20 feedback samples after change
3. Manually trigger evaluation via API if needed

---

## 9. Privacy and Data Retention

- Feedback is anonymized after 90 days (user_id removed)
- Personal comments are sanitized to remove PII
- Machine feedback (outcomes) retained indefinitely for learning
- Users can request their feedback to be deleted (GDPR compliance)

---

## 10. Summary

The GeneQR Feedback System creates a **complete learning loop**:

1. **Collect** - Feedback from humans + system outcomes
2. **Analyze** - Identify patterns and improvement opportunities
3. **Learn** - Apply changes and test their effectiveness
4. **Improve** - Deploy successful changes, rollback failures
5. **Repeat** - Continuous cycle of improvement

**Result:** AI systems that get smarter over time, with measurable improvements in accuracy, user satisfaction, and business outcomes!

---

**For more information:**
- API Reference: `/api/docs`
- Support: support@geneqr.com
- System Status: https://status.geneqr.com
