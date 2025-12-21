# Authentication & Multi-Tenancy API Specification

**Version:** 1.0  
**Date:** December 20, 2025  
**Base URL:** `https://api.aby-med.com/api/v1`  
**Format:** REST API with JSON  

---

## Table of Contents

1. [Authentication Endpoints](#1-authentication-endpoints)
2. [Organization Endpoints](#2-organization-endpoints)
3. [User Management Endpoints](#3-user-management-endpoints)
4. [Ticket Endpoints](#4-ticket-endpoints)
5. [WhatsApp Endpoints](#5-whatsapp-endpoints)
6. [Notification Endpoints](#6-notification-endpoints)
7. [Common Responses](#7-common-responses)
8. [Error Codes](#8-error-codes)

---

## 1. Authentication Endpoints

### 1.1 Send OTP

**Endpoint:** `POST /auth/send-otp`  
**Description:** Send OTP to user's email or phone  
**Authentication:** None (public endpoint with rate limiting)

**Request:**
```json
{
  "identifier": "user@hospital.com",  // or phone: "+1234567890"
  "type": "email",                     // or "sms"
  "purpose": "login"                   // or "verify", "reset"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP sent to user@hospital.com",
  "expires_in": 300,
  "retry_after": 60,
  "request_id": "req-uuid"
}
```

**Error Responses:**
- `429 Too Many Requests` - Rate limit exceeded
- `400 Bad Request` - Invalid identifier format
- `500 Internal Server Error` - OTP delivery failed

**Rate Limiting:**
- 3 OTPs per hour per identifier
- 60-second cooldown between requests

---

### 1.2 Verify OTP

**Endpoint:** `POST /auth/verify-otp`  
**Description:** Verify OTP and login user  
**Authentication:** None

**Request:**
```json
{
  "identifier": "user@hospital.com",
  "otp": "123456",
  "device_info": {
    "user_agent": "Chrome 120",
    "platform": "Windows",
    "ip": "192.168.1.1"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": {
    "id": "user-uuid",
    "email": "user@hospital.com",
    "phone": "+1234567890",
    "name": "John Doe",
    "status": "active",
    "email_verified": true,
    "phone_verified": true,
    "created_at": "2025-01-01T00:00:00Z"
  },
  "organizations": [
    {
      "id": "org-uuid",
      "name": "City Hospital",
      "type": "hospital",
      "role": "engineer",
      "permissions": ["view_tickets", "update_tickets"]
    }
  ]
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid OTP
- `429 Too Many Requests` - Too many attempts (max 3)
- `410 Gone` - OTP expired

---

### 1.3 Login with Password

**Endpoint:** `POST /auth/login-password`  
**Description:** Login with email and password (fallback method)  
**Authentication:** None

**Request:**
```json
{
  "email": "user@hospital.com",
  "password": "SecurePass123!",
  "device_info": {
    "user_agent": "Chrome 120",
    "platform": "Windows",
    "ip": "192.168.1.1"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": { ... },
  "organizations": [ ... ]
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials
- `423 Locked` - Account locked due to failed attempts
- `403 Forbidden` - Account suspended

---

### 1.4 Register User

**Endpoint:** `POST /auth/register`  
**Description:** Register new user account  
**Authentication:** None

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@hospital.com",
  "phone": "+1234567890",
  "organization_id": "org-uuid",
  "role": "engineer",
  "password": "SecurePass123!",  // Optional
  "metadata": {
    "department": "Biomedical Engineering",
    "employee_id": "EMP-123"
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "user_id": "user-uuid",
  "verification_required": true,
  "message": "Verification OTP sent to john@hospital.com",
  "next_step": "verify-otp"
}
```

**Error Responses:**
- `409 Conflict` - Email or phone already exists
- `400 Bad Request` - Validation errors
- `404 Not Found` - Organization not found

---

### 1.5 Refresh Token

**Endpoint:** `POST /auth/refresh`  
**Description:** Refresh access token using refresh token  
**Authentication:** Refresh token required

**Request:**
```json
{
  "refresh_token": "eyJhbG..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "access_token": "eyJhbG...",     // New access token
  "refresh_token": "eyJhbG...",    // New refresh token (rotation)
  "expires_in": 900,
  "token_type": "Bearer"
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or expired refresh token
- `403 Forbidden` - Token revoked

---

### 1.6 Logout

**Endpoint:** `POST /auth/logout`  
**Description:** Logout current session  
**Authentication:** Bearer token required

**Request:**
```json
{
  "revoke_all": false  // If true, revoke all sessions
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

### 1.7 Get Current User

**Endpoint:** `GET /auth/me`  
**Description:** Get current authenticated user profile  
**Authentication:** Bearer token required

**Response (200 OK):**
```json
{
  "success": true,
  "user": {
    "id": "user-uuid",
    "email": "user@hospital.com",
    "phone": "+1234567890",
    "name": "John Doe",
    "status": "active",
    "email_verified": true,
    "phone_verified": true,
    "preferred_auth_method": "otp",
    "last_login": "2025-12-20T10:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-12-20T09:00:00Z"
  },
  "organizations": [ ... ],
  "current_organization": {
    "id": "org-uuid",
    "name": "City Hospital",
    "type": "hospital",
    "role": "engineer",
    "permissions": ["view_tickets", "update_tickets", "create_tickets"]
  }
}
```

---

### 1.8 List Active Sessions

**Endpoint:** `GET /auth/sessions`  
**Description:** List all active sessions for current user  
**Authentication:** Bearer token required

**Response (200 OK):**
```json
{
  "success": true,
  "sessions": [
    {
      "id": "session-uuid",
      "device_info": {
        "user_agent": "Chrome 120",
        "platform": "Windows",
        "ip": "192.168.1.1"
      },
      "created_at": "2025-12-20T10:00:00Z",
      "last_activity": "2025-12-20T10:30:00Z",
      "expires_at": "2025-12-27T10:00:00Z",
      "is_current": true
    }
  ]
}
```

---

### 1.9 Revoke Session

**Endpoint:** `DELETE /auth/sessions/{session_id}`  
**Description:** Revoke a specific session  
**Authentication:** Bearer token required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Session revoked successfully"
}
```

---

### 1.10 Forgot Password

**Endpoint:** `POST /auth/forgot-password`  
**Description:** Request password reset OTP  
**Authentication:** None

**Request:**
```json
{
  "email": "user@hospital.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset OTP sent to user@hospital.com",
  "expires_in": 300,
  "request_id": "req-uuid"
}
```

---

### 1.11 Verify Reset OTP

**Endpoint:** `POST /auth/verify-reset-otp`  
**Description:** Verify password reset OTP  
**Authentication:** None

**Request:**
```json
{
  "email": "user@hospital.com",
  "otp": "123456"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "reset_token": "temp-reset-token-uuid",
  "expires_in": 300,
  "message": "OTP verified. Proceed to reset password."
}
```

---

### 1.12 Reset Password

**Endpoint:** `POST /auth/reset-password`  
**Description:** Set new password after OTP verification  
**Authentication:** Reset token required

**Request:**
```json
{
  "reset_token": "temp-reset-token-uuid",
  "new_password": "NewSecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successfully. All sessions have been revoked.",
  "next_step": "login"
}
```

---

## 2. Organization Endpoints

### 2.1 List Organizations

**Endpoint:** `GET /organizations`  
**Description:** List all organizations (with pagination)  
**Authentication:** Bearer token required  
**Permissions:** `view_organizations` or own organization

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20, max: 100)
- `type` (filter: manufacturer, hospital, laboratory, distributor, dealer)
- `status` (filter: active, pending, suspended)
- `search` (search in name)

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": "org-uuid",
      "name": "City Hospital",
      "type": "hospital",
      "status": "active",
      "onboarding_status": "completed",
      "address": {
        "street": "123 Main St",
        "city": "Springfield",
        "state": "IL",
        "zip": "62701",
        "country": "USA"
      },
      "contact": {
        "email": "contact@cityhospital.com",
        "phone": "+1234567890"
      },
      "metadata": {
        "license_number": "HOSP-123",
        "bed_count": 500,
        "specialties": ["Cardiology", "Neurology"]
      },
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-12-20T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

---

### 2.2 Create Organization

**Endpoint:** `POST /organizations`  
**Description:** Create new organization  
**Authentication:** Bearer token required  
**Permissions:** `create_organization` (admin only)

**Request:**
```json
{
  "name": "Metro Health System",
  "type": "hospital",
  "address": {
    "street": "456 Oak Ave",
    "city": "Chicago",
    "state": "IL",
    "zip": "60601",
    "country": "USA"
  },
  "contact": {
    "email": "admin@metrohealth.com",
    "phone": "+1987654321"
  },
  "metadata": {
    "license_number": "HOSP-456",
    "bed_count": 800
  }
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "organization": {
    "id": "org-uuid",
    "name": "Metro Health System",
    "type": "hospital",
    "status": "pending",
    "onboarding_status": "created",
    ...
  },
  "message": "Organization created. Awaiting verification."
}
```

---

### 2.3 Get Organization Details

**Endpoint:** `GET /organizations/{id}`  
**Description:** Get detailed organization information  
**Authentication:** Bearer token required  
**Permissions:** Member of organization or admin

**Response (200 OK):**
```json
{
  "success": true,
  "organization": {
    "id": "org-uuid",
    "name": "City Hospital",
    "type": "hospital",
    "status": "active",
    "stats": {
      "total_users": 45,
      "total_equipment": 120,
      "active_tickets": 15,
      "engineers": 12
    },
    ...
  }
}
```

---

### 2.4 Update Organization

**Endpoint:** `PUT /organizations/{id}`  
**Description:** Update organization details  
**Authentication:** Bearer token required  
**Permissions:** `manage_organization` role in that org

**Request:**
```json
{
  "name": "City Hospital - Main Campus",
  "contact": {
    "email": "newcontact@cityhospital.com"
  },
  "metadata": {
    "bed_count": 550
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "organization": { ... },
  "message": "Organization updated successfully"
}
```

---

### 2.5 Get Organization Dashboard

**Endpoint:** `GET /organizations/{id}/dashboard`  
**Description:** Get organization-specific dashboard data  
**Authentication:** Bearer token required  
**Permissions:** Member of organization

**Response (200 OK):**
```json
{
  "success": true,
  "dashboard": {
    "cards": [
      {
        "id": "total_equipment",
        "title": "Total Equipment",
        "type": "stats",
        "value": 120,
        "change": "+5",
        "change_type": "increase"
      },
      {
        "id": "active_tickets",
        "title": "Active Service Tickets",
        "type": "stats",
        "value": 15,
        "urgent": 3
      },
      {
        "id": "equipment_uptime",
        "title": "Equipment Uptime",
        "type": "chart",
        "data": {
          "current": 98.5,
          "target": 99.0,
          "trend": [98.2, 98.5, 98.7, 98.5]
        }
      }
    ],
    "recent_activity": [ ... ]
  }
}
```

---

### 2.6 List Organization Users

**Endpoint:** `GET /organizations/{id}/users`  
**Description:** List all users in organization  
**Authentication:** Bearer token required  
**Permissions:** `view_users` in organization

**Query Parameters:**
- `page`, `limit`
- `role` (filter by role)
- `status` (filter by status)

**Response (200 OK):**
```json
{
  "success": true,
  "users": [
    {
      "id": "user-uuid",
      "name": "John Doe",
      "email": "john@cityhospital.com",
      "role": "engineer",
      "status": "active",
      "permissions": ["view_tickets", "update_tickets"],
      "joined_at": "2025-01-15T00:00:00Z",
      "last_login": "2025-12-20T10:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

### 2.7 Add User to Organization

**Endpoint:** `POST /organizations/{id}/users`  
**Description:** Add existing user to organization or invite new user  
**Authentication:** Bearer token required  
**Permissions:** `manage_users` in organization

**Request:**
```json
{
  "user_id": "existing-user-uuid",  // Or
  "email": "newuser@cityhospital.com",  // For new user
  "role": "engineer",
  "permissions": ["view_tickets", "update_tickets"]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User added to organization",
  "user_organization": {
    "user_id": "user-uuid",
    "organization_id": "org-uuid",
    "role": "engineer",
    "joined_at": "2025-12-20T10:00:00Z"
  }
}
```

---

### 2.8 Remove User from Organization

**Endpoint:** `DELETE /organizations/{id}/users/{user_id}`  
**Description:** Remove user from organization  
**Authentication:** Bearer token required  
**Permissions:** `manage_users` in organization

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User removed from organization"
}
```

---

## 3. User Management Endpoints

### 3.1 Update User Profile

**Endpoint:** `PUT /users/me`  
**Description:** Update current user's profile  
**Authentication:** Bearer token required

**Request:**
```json
{
  "name": "John Smith",
  "phone": "+1234567890",
  "preferred_auth_method": "otp",
  "metadata": {
    "department": "Maintenance"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "user": { ... },
  "message": "Profile updated successfully"
}
```

---

### 3.2 Change Password

**Endpoint:** `POST /users/me/change-password`  
**Description:** Change password (requires current password)  
**Authentication:** Bearer token required

**Request:**
```json
{
  "current_password": "OldPass123!",
  "new_password": "NewPass123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully. All other sessions revoked."
}
```

---

### 3.3 Update Notification Preferences

**Endpoint:** `PUT /users/me/notification-preferences`  
**Description:** Update notification settings  
**Authentication:** Bearer token required

**Request:**
```json
{
  "email_notifications": true,
  "sms_notifications": false,
  "whatsapp_notifications": true,
  "events": {
    "ticket_created": true,
    "ticket_assigned": true,
    "ticket_updated": true,
    "ticket_completed": true
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "preferences": { ... },
  "message": "Preferences updated"
}
```

---

## 4. Ticket Endpoints (Enhanced)

### 4.1 Create Ticket (Authenticated)

**Endpoint:** `POST /tickets`  
**Description:** Create service ticket (authenticated user)  
**Authentication:** Bearer token required

**Request:**
```json
{
  "equipment_id": "eq-uuid",
  "issue_description": "X-ray machine not powering on",
  "priority": "high",
  "attachments": ["attachment-uuid-1", "attachment-uuid-2"]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "ticket": {
    "id": "ticket-uuid",
    "ticket_number": "TKT-12345",
    "equipment_id": "eq-uuid",
    "created_by": "user-uuid",
    "organization_id": "org-uuid",
    "issue_description": "X-ray machine not powering on",
    "priority": "high",
    "status": "pending",
    "created_at": "2025-12-20T10:00:00Z"
  },
  "message": "Ticket created successfully",
  "tracking_url": "https://app.aby-med.com/track/TKT-12345"
}
```

---

### 4.2 Create Anonymous Ticket

**Endpoint:** `POST /tickets/anonymous`  
**Description:** Create ticket without authentication (QR/WhatsApp)  
**Authentication:** None (reCAPTCHA required)

**Request:**
```json
{
  "equipment_id": "eq-uuid",
  "issue_description": "Equipment malfunction",
  "contact_info": {
    "email": "user@hospital.com",
    "phone": "+1234567890",
    "whatsapp": "+1234567890",
    "name": "John Maintenance Staff",
    "preferred_channel": "whatsapp"
  },
  "recaptcha_token": "recaptcha-token",
  "source": "qr"  // or "whatsapp"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "ticket": {
    "id": "ticket-uuid",
    "ticket_number": "TKT-12346",
    ...
  },
  "message": "Ticket created. You'll receive updates via WhatsApp.",
  "tracking_url": "https://app.aby-med.com/track/TKT-12346"
}
```

---

### 4.3 Track Ticket (Public)

**Endpoint:** `GET /tickets/{ticket_number}/track`  
**Description:** Track ticket status without authentication  
**Authentication:** None

**Response (200 OK):**
```json
{
  "success": true,
  "ticket": {
    "ticket_number": "TKT-12345",
    "status": "in_progress",
    "created_at": "2025-12-20T10:00:00Z",
    "equipment": {
      "name": "X-Ray Machine - Model XR-500",
      "location": "Radiology Department"
    },
    "assigned_engineer": {
      "name": "Engineer John",
      "phone": "+1234567890",
      "eta": "2025-12-20T14:00:00Z"
    },
    "timeline": [
      {
        "status": "created",
        "timestamp": "2025-12-20T10:00:00Z",
        "message": "Ticket created"
      },
      {
        "status": "assigned",
        "timestamp": "2025-12-20T10:30:00Z",
        "message": "Engineer John assigned"
      }
    ]
  }
}
```

---

### 4.4 Update Ticket Status

**Endpoint:** `PUT /tickets/{id}/status`  
**Description:** Update ticket status  
**Authentication:** Bearer token required  
**Permissions:** `update_tickets`

**Request:**
```json
{
  "status": "in_progress",
  "notes": "Arrived on-site. Diagnosing issue.",
  "attachments": ["photo-uuid-1"]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "ticket": { ... },
  "message": "Ticket status updated"
}
```

---

### 4.5 Assign Parts to Ticket

**Endpoint:** `POST /tickets/{id}/parts`  
**Description:** Assign spare parts to ticket  
**Authentication:** Bearer token required  
**Permissions:** `update_tickets`

**Request:**
```json
{
  "parts": [
    {
      "spare_part_id": "part-uuid",
      "quantity": 2,
      "unit_price": 150.00
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "parts_assigned": 1,
  "total_cost": 300.00,
  "message": "Parts assigned to ticket"
}
```

---

### 4.6 List Tickets

**Endpoint:** `GET /tickets`  
**Description:** List tickets with filters  
**Authentication:** Bearer token required

**Query Parameters:**
- `status` (filter)
- `priority` (filter)
- `assigned_to` (engineer_id)
- `organization_id` (filter)
- `page`, `limit`
- `sort` (created_at, priority, status)

**Response (200 OK):**
```json
{
  "success": true,
  "tickets": [ ... ],
  "pagination": { ... }
}
```

---

## 5. WhatsApp Endpoints

### 5.1 WhatsApp Webhook

**Endpoint:** `POST /whatsapp/webhook`  
**Description:** Receive incoming WhatsApp messages (Twilio webhook)  
**Authentication:** Twilio signature validation

**Request (from Twilio):**
```json
{
  "From": "whatsapp:+1234567890",
  "To": "whatsapp:+1800ABYMED",
  "Body": "My machine is broken",
  "MessageSid": "SM..."
}
```

**Response (200 OK):**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<Response>
  <Message>Thank you! I'll help you create a ticket. Please provide the Equipment ID from the QR code.</Message>
</Response>
```

---

### 5.2 Send WhatsApp Message

**Endpoint:** `POST /whatsapp/send`  
**Description:** Send WhatsApp message (internal use)  
**Authentication:** Bearer token required  
**Permissions:** `send_notifications`

**Request:**
```json
{
  "to": "+1234567890",
  "message": "Your ticket TKT-12345 has been assigned to Engineer John.",
  "template": "ticket_assigned",
  "variables": {
    "ticket_number": "TKT-12345",
    "engineer_name": "John"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message_sid": "SM...",
  "status": "queued"
}
```

---

## 6. Notification Endpoints

### 6.1 List Notifications

**Endpoint:** `GET /notifications`  
**Description:** Get user notifications  
**Authentication:** Bearer token required

**Query Parameters:**
- `unread` (boolean)
- `page`, `limit`

**Response (200 OK):**
```json
{
  "success": true,
  "notifications": [
    {
      "id": "notif-uuid",
      "type": "ticket_assigned",
      "title": "New Ticket Assigned",
      "message": "Ticket TKT-12345 has been assigned to you",
      "data": {
        "ticket_id": "ticket-uuid",
        "ticket_number": "TKT-12345"
      },
      "read": false,
      "created_at": "2025-12-20T10:00:00Z"
    }
  ],
  "unread_count": 5,
  "pagination": { ... }
}
```

---

### 6.2 Mark Notification as Read

**Endpoint:** `PUT /notifications/{id}/read`  
**Description:** Mark notification as read  
**Authentication:** Bearer token required

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Notification marked as read"
}
```

---

### 6.3 Send Notification (Internal)

**Endpoint:** `POST /notifications/send`  
**Description:** Send notification to user (internal use)  
**Authentication:** Bearer token required  
**Permissions:** `send_notifications`

**Request:**
```json
{
  "user_id": "user-uuid",
  "type": "ticket_assigned",
  "title": "New Ticket Assigned",
  "message": "Ticket TKT-12345 assigned to you",
  "channels": ["email", "sms", "whatsapp", "in_app"],
  "data": {
    "ticket_id": "ticket-uuid"
  }
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "sent_via": ["email", "whatsapp", "in_app"],
  "failed_via": ["sms"],
  "message": "Notification sent"
}
```

---

## 7. Common Responses

### Success Response Structure
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful",
  "metadata": {
    "request_id": "req-uuid",
    "timestamp": "2025-12-20T10:00:00Z"
  }
}
```

### Error Response Structure
```json
{
  "success": false,
  "error": {
    "code": "INVALID_INPUT",
    "message": "Invalid email format",
    "field": "email",
    "details": { ... }
  },
  "request_id": "req-uuid",
  "timestamp": "2025-12-20T10:00:00Z"
}
```

---

## 8. Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_INPUT` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required |
| `INVALID_TOKEN` | 401 | Token invalid or expired |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily down |

---

## Authentication

All authenticated endpoints require:
```
Authorization: Bearer {access_token}
```

Token must be included in every request header.

---

**END OF API SPECIFICATION**

Total Endpoints: 40+  
Version: 1.0  
Last Updated: December 20, 2025
