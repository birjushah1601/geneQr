# ðŸ“‹ Assignment API Documentation

**Version:** 1.0  
**Base URL:** `/api/v1`

---

## ðŸŽ¯ Overview

The Assignment API manages engineer assignments for service tickets, including assignment creation, escalation, and tracking.

---

## ðŸ“š Endpoints

### **Ticket Assignment**

#### **POST /tickets/{ticketId}/assign**
Assign a ticket to an engineer.

**Request Body:**
```json
{
  "engineer_id": "ENG-001",
  "equipment_id": "EQ-12345",
  "tier": 1,
  "tier_name": "OEM Engineer",
  "reason": "Initial assignment",
  "type": "manual",
  "assigned_by": "USER-123"
}
```

**Response:** `201 Created`
```json
{
  "id": "ASSIGN-abc123",
  "ticket_id": "TKT-20231116-001",
  "engineer_id": "ENG-001",
  "equipment_id": "EQ-12345",
  "assignment_sequence": 1,
  "assignment_tier": 1,
  "assignment_tier_name": "OEM Engineer",
  "assignment_reason": "Initial assignment",
  "assignment_type": "manual",
  "status": "assigned",
  "assigned_by": "USER-123",
  "assigned_at": "2025-11-16T10:30:00Z",
  "created_at": "2025-11-16T10:30:00Z",
  "updated_at": "2025-11-16T10:30:00Z"
}
```

---

#### **POST /tickets/{ticketId}/escalate**
Escalate a ticket to the next tier.

**Request Body:**
```json
{
  "reason": "Issue requires specialized expertise",
  "next_engineer_id": "ENG-002",
  "next_tier_name": "Sub-sub_SUB_DEALER Engineer",
  "escalated_by": "ENG-001"
}
```

**Response:** `201 Created`
```json
{
  "id": "ASSIGN-xyz789",
  "ticket_id": "TKT-20231116-001",
  "engineer_id": "ENG-002",
  "assignment_sequence": 2,
  "assignment_tier": 2,
  "assignment_tier_name": "Sub-sub_SUB_DEALER Engineer",
  "assignment_reason": "Escalation from tier 1",
  "status": "assigned",
  "escalation_reason": "Issue requires specialized expertise",
  "assigned_at": "2025-11-16T11:00:00Z"
}
```

---

#### **GET /tickets/{ticketId}/current-assignment**
Get the current active assignment for a ticket.

**Response:** `200 OK`
```json
{
  "id": "ASSIGN-abc123",
  "ticket_id": "TKT-20231116-001",
  "engineer_id": "ENG-001",
  "status": "in_progress",
  "assigned_at": "2025-11-16T10:30:00Z",
  "started_at": "2025-11-16T10:45:00Z"
}
```

---

#### **GET /tickets/{ticketId}/assignments**
Get complete assignment history for a ticket.

**Response:** `200 OK`
```json
[
  {
    "id": "ASSIGN-abc123",
    "assignment_sequence": 1,
    "engineer_id": "ENG-001",
    "status": "escalated",
    "assigned_at": "2025-11-16T10:30:00Z",
    "completed_at": "2025-11-16T11:00:00Z",
    "escalation_reason": "Issue requires specialized expertise"
  },
  {
    "id": "ASSIGN-xyz789",
    "assignment_sequence": 2,
    "engineer_id": "ENG-002",
    "status": "in_progress",
    "assigned_at": "2025-11-16T11:00:00Z",
    "started_at": "2025-11-16T11:15:00Z"
  }
]
```

---

### **Assignment Actions**

#### **POST /assignments/{assignmentId}/accept**
Engineer accepts an assignment.

**Request Body:**
```json
{
  "engineer_id": "ENG-001"
}
```

**Response:** `200 OK`
```json
{
  "message": "Assignment accepted successfully"
}
```

---

#### **POST /assignments/{assignmentId}/reject**
Engineer rejects an assignment.

**Request Body:**
```json
{
  "engineer_id": "ENG-001",
  "reason": "Not available at this time"
}
```

**Response:** `200 OK`
```json
{
  "message": "Assignment rejected successfully"
}
```

---

#### **POST /assignments/{assignmentId}/start**
Engineer starts working on an assignment.

**Request Body:**
```json
{
  "engineer_id": "ENG-001"
}
```

**Response:** `200 OK`
```json
{
  "message": "Assignment started successfully"
}
```

---

#### **POST /assignments/{assignmentId}/complete**
Engineer completes an assignment.

**Request Body:**
```json
{
  "engineer_id": "ENG-001",
  "completion_status": "success",
  "diagnosis": "Faulty sensor detected",
  "actions_taken": "Replaced sensor and recalibrated",
  "parts_used": [
    {
      "part_number": "SEN-001",
      "part_name": "Temperature Sensor",
      "quantity": 1,
      "cost": 250.00
    }
  ],
  "time_spent_hours": 2.5
}
```

**Response:** `200 OK`
```json
{
  "message": "Assignment completed successfully"
}
```

---

#### **POST /assignments/{assignmentId}/feedback**
Add customer feedback for a completed assignment.

**Request Body:**
```json
{
  "rating": 5,
  "feedback": "Excellent service, very professional"
}
```

**Response:** `200 OK`
```json
{
  "message": "Feedback added successfully"
}
```

---

### **Engineer Workload**

#### **GET /engineers/{engineerId}/assignments**
Get assignments for an engineer.

**Query Parameters:**
- `limit` (optional): Maximum number of assignments to return (default: 50)

**Response:** `200 OK`
```json
[
  {
    "id": "ASSIGN-abc123",
    "ticket_id": "TKT-20231116-001",
    "status": "completed",
    "assigned_at": "2025-11-16T10:30:00Z",
    "completed_at": "2025-11-16T13:00:00Z",
    "time_spent_hours": 2.5,
    "customer_rating": 5
  }
]
```

---

#### **GET /engineers/{engineerId}/assignments/active**
Get active assignments for an engineer.

**Response:** `200 OK`
```json
[
  {
    "id": "ASSIGN-xyz789",
    "ticket_id": "TKT-20231116-002",
    "status": "in_progress",
    "assigned_at": "2025-11-16T14:00:00Z",
    "started_at": "2025-11-16T14:15:00Z"
  }
]
```

---

#### **GET /engineers/{engineerId}/workload**
Get workload statistics for an engineer.

**Response:** `200 OK`
```json
{
  "engineer_id": "ENG-001",
  "active_count": 3,
  "completed_count": 45,
  "avg_hours": 2.8
}
```

---

## ðŸ”‘ Data Models

### **Assignment Status**
- `assigned` - Ticket assigned to engineer
- `accepted` - Engineer accepted the assignment
- `rejected` - Engineer rejected the assignment
- `in_progress` - Engineer is working on it
- `completed` - Work completed
- `failed` - Assignment failed
- `escalated` - Escalated to next tier

### **Completion Status**
- `success` - Successfully resolved
- `failed` - Could not resolve
- `escalated` - Escalated to higher tier
- `parts_required` - Waiting for parts
- `customer_unavailable` - Customer not available

### **Assignment Type**
- `auto` - Automatically assigned by system
- `manual` - Manually assigned by dispatcher
- `escalation` - Created through escalation

---

## âš ï¸ Error Responses

### **400 Bad Request**
```json
{
  "error": "Invalid request body: field 'engineer_id' is required"
}
```

### **404 Not Found**
```json
{
  "error": "No active assignment found for ticket"
}
```

### **500 Internal Server Error**
```json
{
  "error": "Failed to assign ticket: database connection error"
}
```

---

## ðŸ”„ Typical Workflow

### **Scenario: Ticket Assignment & Escalation**

1. **Assign ticket to L1 engineer:**
   ```
   POST /tickets/TKT-001/assign
   ```

2. **Engineer accepts:**
   ```
   POST /assignments/ASSIGN-001/accept
   ```

3. **Engineer starts work:**
   ```
   POST /assignments/ASSIGN-001/start
   ```

4. **Issue requires escalation:**
   ```
   POST /tickets/TKT-001/escalate
   ```

5. **L2 engineer accepts:**
   ```
   POST /assignments/ASSIGN-002/accept
   ```

6. **L2 engineer completes:**
   ```
   POST /assignments/ASSIGN-002/complete
   ```

7. **Customer provides feedback:**
   ```
   POST /assignments/ASSIGN-002/feedback
   ```

8. **View complete history:**
   ```
   GET /tickets/TKT-001/assignments
   ```

---

## ðŸ“Š Performance Considerations

- Assignment history queries are optimized with indexes on `ticket_id` and `assignment_sequence`
- Current assignment lookup uses a partial index for active assignments only
- Engineer workload queries cache results for 5 minutes
- Recommend pagination for large assignment lists (use `limit` parameter)

---

## ðŸ”’ Authentication

All endpoints require authentication. Include the bearer token in the Authorization header:

```
Authorization: Bearer <your-token>
```

---

## ðŸ“ Changelog

### Version 1.0 (2025-11-16)
- Initial release
- Complete assignment lifecycle management
- Escalation support
- Workload tracking
- Customer feedback

---

**Questions?** Contact the engineering team or refer to the [Master Fix Plan](../database/MASTER-FIX-PLAN.md).
