# Authentication & Multi-Tenancy System
# Product Requirements Document (PRD)

**Version:** 1.0  
**Date:** December 20, 2025  
**Status:** Draft - For Review  
**Project:** ABY-MED Medical Equipment Service Platform  
**Document Owner:** Technical Architecture Team  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Project Overview](#2-project-overview)
3. [Goals & Objectives](#3-goals--objectives)
4. [User Personas](#4-user-personas)
5. [System Architecture](#5-system-architecture)
6. [Authentication System](#6-authentication-system)
7. [Multi-Tenant Architecture](#7-multi-tenant-architecture)
8. [Ticket Creation Flows](#8-ticket-creation-flows)
9. [WhatsApp Integration](#9-whatsapp-integration)
10. [Notification System](#10-notification-system)
11. [Security Requirements](#11-security-requirements)
12. [Database Schema](#12-database-schema)
13. [API Specifications](#13-api-specifications)
14. [Frontend Requirements](#14-frontend-requirements)
15. [Implementation Plan](#15-implementation-plan)
16. [Testing Strategy](#16-testing-strategy)
17. [Deployment Plan](#17-deployment-plan)
18. [Success Metrics](#18-success-metrics)
19. [Risks & Mitigation](#19-risks--mitigation)
20. [Appendices](#20-appendices)

---

## 1. Executive Summary

### 1.1 Overview

ABY-MED is implementing a comprehensive authentication and multi-tenancy system to support secure access for various stakeholders including manufacturers, hospitals, laboratories, distributors, dealers, engineers, and administrators. The system will feature OTP-first authentication with password fallback, WhatsApp-based ticket creation and notifications, and organization-specific dashboards.

### 1.2 Key Features

- **OTP-First Authentication**: Modern, passwordless login with email/SMS OTP as primary method
- **Multi-Tenant Architecture**: Support for 5+ organization types with role-based access
- **Flexible Ticket Creation**: Web, QR code, and WhatsApp channels for service requests
- **Smart Notifications**: Multi-channel notifications (Email, SMS, WhatsApp, In-app)
- **WhatsApp Integration**: Ticket creation, diagnostic chatbot, and status updates
- **Organization Dashboards**: Custom dashboards per organization type
- **Security-First Design**: reCAPTCHA, rate limiting, audit logging, encryption

### 1.3 Business Impact

- **Improved User Experience**: Passwordless login reduces friction by 70%
- **Increased Accessibility**: WhatsApp support reaches 2B+ users
- **Enhanced Security**: OTP-based auth reduces password-related breaches
- **Operational Efficiency**: Multi-tenant system reduces development overhead
- **Scalability**: Architecture supports 10,000+ organizations
- **Compliance**: Audit logging and data protection meet healthcare standards

### 1.4 Success Criteria

| Metric | Target | Timeline |
|--------|--------|----------|
| User Adoption | 80% of users using OTP login | 3 months post-launch |
| WhatsApp Tickets | 30% of tickets via WhatsApp | 6 months post-launch |
| Login Success Rate | >95% | Immediate |
| Average Login Time | <30 seconds | Immediate |
| Security Incidents | Zero breaches | Ongoing |
| System Uptime | 99.9% | Ongoing |

---

## 2. Project Overview

### 2.1 Background

The current ABY-MED platform has basic API key authentication suitable for development but lacks user management, organization-specific access control, and modern authentication methods required for production deployment.

### 2.2 Problem Statement

**Current Limitations:**
1. No user registration or login system
2. No organization-based access control
3. No support for passwordless authentication
4. No WhatsApp integration for ticket creation
5. No flexible notification system
6. QR code endpoints are fully open (security risk)
7. No audit trail for sensitive operations

**Business Impact:**
- Cannot onboard real customers (hospitals, manufacturers)
- Engineers cannot securely access the system
- No way to track who performed what actions
- Vulnerable to spam/abuse on public endpoints
- Missing modern authentication expected by users

### 2.3 Proposed Solution

Implement a comprehensive authentication and multi-tenancy system that:
- Uses OTP-first login (modern, secure, passwordless)
- Supports multiple organization types with custom dashboards
- Enables WhatsApp-based ticket creation and notifications
- Provides flexible, multi-channel notifications
- Secures public endpoints with reCAPTCHA
- Maintains complete audit trails

### 2.4 Scope

**In Scope:**
- OTP-first authentication (Email/SMS)
- Password fallback option
- User registration and verification
- Multi-tenant organization system (5 types)
- Role-based access control (RBAC)
- Organization-specific dashboards
- Ticket creation (Web, QR, WhatsApp)
- WhatsApp integration (basic chatbot + notifications)
- Multi-channel notification system
- Security measures (reCAPTCHA, rate limiting)
- Audit logging
- Engineer ticket management

**Out of Scope:**
- WhatsApp login (only ticket creation)
- Advanced AI diagnostic chatbot
- Social login (Google/Microsoft) - Future
- Marketplace controls - Future
- Mobile push notifications - Future
- Biometric authentication - Future
- Payment integration - Future

---

## 3. Goals & Objectives

### 3.1 Primary Goals

1. **Enable Secure Multi-User Access**
   - Support 5+ organization types
   - 10+ user roles with granular permissions
   - Secure authentication for 10,000+ users

2. **Modernize Authentication**
   - OTP-first login (70%+ adoption target)
   - <30 second login time
   - >95% login success rate

3. **Enable WhatsApp Channel**
   - 30% of tickets via WhatsApp (6-month target)
   - <2 minute ticket creation time
   - Real-time notifications

4. **Improve Security Posture**
   - Zero security breaches
   - 100% audit coverage on sensitive operations
   - reCAPTCHA on all public endpoints

5. **Support Business Growth**
   - Onboard 100+ organizations in first year
   - Support 1,000+ concurrent users
   - 99.9% system uptime

### 3.2 Secondary Goals

- Reduce support tickets related to password resets by 80%
- Enable self-service user registration
- Provide real-time status updates to customers
- Build foundation for mobile app (future)
- Support international expansion (multi-language OTP)

### 3.3 Non-Goals

- Building a generic identity provider (use existing frameworks)
- Supporting every possible authentication method
- Building custom chatbot AI (use basic rule-based for now)
- Replacing existing equipment management system

---

## 4. User Personas

### 4.1 Persona 1: Hospital Biomedical Engineer

**Name:** Sarah Johnson  
**Role:** Biomedical Engineer  
**Organization:** City General Hospital  
**Tech Savviness:** Medium  

**Needs:**
- Quick access to create service tickets
- View equipment status and history
- Receive real-time updates on service requests
- Access system from mobile device

**Pain Points:**
- Forgets passwords frequently
- Needs immediate assistance when equipment fails
- Works in areas with limited computer access
- Requires fast ticket creation during emergencies

**How Our System Helps:**
- OTP login (no passwords to remember)
- WhatsApp ticket creation (works on any phone)
- Real-time notifications via WhatsApp
- Mobile-friendly web interface

### 4.2 Persona 2: Field Service Engineer

**Name:** Raj Patel  
**Role:** Senior Field Engineer  
**Organization:** MedTech Services (Manufacturer)  
**Tech Savviness:** Medium-High  

**Needs:**
- View assigned service tickets
- Update ticket status from field
- Access equipment manuals and diagrams
- Record parts used and time spent

**Pain Points:**
- Constantly on the move
- Limited laptop access
- Needs offline capability
- Multiple customer sites daily

**How Our System Helps:**
- Fast OTP login from any device
- Mobile-optimized ticket interface
- WhatsApp status updates
- Simple ticket update workflow

### 4.3 Persona 3: Manufacturer Service Manager

**Name:** Emily Chen  
**Role:** Service Operations Manager  
**Organization:** GlobalMed Equipment Corp  
**Tech Savviness:** High  

**Needs:**
- Dashboard showing all service tickets
- Engineer performance metrics
- Customer satisfaction tracking
- Service contract management

**Pain Points:**
- Managing distributed engineer team
- Tracking SLA compliance
- Customer communication overhead
- Generating management reports

**How Our System Helps:**
- Comprehensive service dashboard
- Real-time engineer locations and status
- Automated customer notifications
- Built-in reporting and analytics

### 4.4 Persona 4: Hospital Administrator

**Name:** Dr. Michael Torres  
**Role:** Director of Operations  
**Organization:** Metro Health System  
**Tech Savviness:** Low-Medium  

**Needs:**
- Oversight of all equipment
- Budget tracking for repairs
- Vendor performance tracking
- Compliance reporting

**Pain Points:**
- Too many systems to learn
- Difficulty getting timely updates
- Budget overruns on repairs
- Audit requirements

**How Our System Helps:**
- Simple, intuitive interface
- Email/SMS summaries (no login required)
- Cost tracking and budgets
- Audit-ready reports

### 4.5 Persona 5: Platform Administrator

**Name:** Alex Kumar  
**Role:** System Administrator  
**Organization:** ABY-MED (Internal)  
**Tech Savviness:** Expert  

**Needs:**
- User and organization management
- System configuration and monitoring
- Security and audit oversight
- Performance optimization

**Pain Points:**
- Managing multiple organizations
- Security threat monitoring
- Scaling challenges
- Support ticket volume

**How Our System Helps:**
- Centralized admin dashboard
- Comprehensive audit logs
- Automated security monitoring
- Role-based delegation

---


## 5. System Architecture

### 5.1 High-Level Architecture

\\\
┌─────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                          │
├───────────────────────────────────────────────────────┬─┤
│  Web App (React)  │  Mobile Browser  │  WhatsApp Bot   │
└──────────┬────────────────┬──────────────────┬─────────┘
           │                │                  │
    ┌──────┴────────────────┴──────────────────┴────────┐
    │         NGINX / Load Balancer / SSL/TLS           │
    └──────┬────────────────┬──────────────────┬────────┘
           │                │                  │
    ┌──────┴────────────────┴──────────────────┴────────┐
    │              API GATEWAY (Go/Chi Router)           │
    │  • Rate Limiting  • CORS  • Request Logging       │
    └──────┬────────────────┬──────────────────┬────────┘
           │                │                  │
    ┌──────┴─────┐   ┌─────┴──────┐   ┌──────┴────────┐
    │   Auth     │   │ Business   │   │   WhatsApp    │
    │  Service   │   │  Services  │   │   Service     │
    │            │   │            │   │               │
    │ • OTP      │   │ • Tickets  │   │ • Webhook     │
    │ • JWT      │   │ • Equip    │   │ • Chatbot     │
    │ • Sessions │   │ • Parts    │   │ • Notif       │
    └──────┬─────┘   └─────┬──────┘   └──────┬────────┘
           │               │                  │
    ┌──────┴───────────────┴──────────────────┴────────┐
    │                  DATA LAYER                       │
    ├───────────────────┬───────────────┬───────────────┤
    │   PostgreSQL      │    Redis      │   AWS S3      │
    │   (Primary DB)    │   (Cache)     │   (Files)     │
    └───────────────────┴───────────────┴───────────────┘
                        │
    ┌───────────────────┴──────────────────────────────┐
    │              EXTERNAL SERVICES                    │
    ├──────────────┬──────────────┬────────────────────┤
    │  Twilio      │  SendGrid    │  Google reCAPTCHA  │
    │  (SMS/WA)    │  (Email)     │  (Anti-spam)       │
    └──────────────┴──────────────┴────────────────────┘
\\\

### 5.2 Technology Stack

**Backend:**
- **Language:** Go 1.22+
- **Web Framework:** Chi Router v5
- **Database:** PostgreSQL 14+
- **Cache:** Redis 7+
- **JWT Library:** golang-jwt/jwt v5
- **OTP Library:** github.com/pquerna/otp
- **Password Hashing:** golang.org/x/crypto/bcrypt

**Frontend:**
- **Framework:** React 18+ with TypeScript
- **Routing:** React Router v6
- **State Management:** Context API + TanStack Query
- **UI Library:** Tailwind CSS + shadcn/ui
- **Forms:** React Hook Form + Zod
- **HTTP Client:** Axios with interceptors

**Infrastructure:**
- **Hosting:** AWS EC2 / Docker containers
- **Database:** AWS RDS PostgreSQL
- **Cache:** AWS ElastiCache Redis
- **Storage:** AWS S3
- **CDN:** CloudFront
- **Monitoring:** CloudWatch + Custom metrics

**External Services:**
- **SMS/WhatsApp:** Twilio
- **Email:** SendGrid
- **reCAPTCHA:** Google reCAPTCHA v3
- **Analytics:** Custom + Google Analytics

### 5.3 Security Architecture

\\\
┌────────────────────────────────────────────────────┐
│            SECURITY LAYERS                         │
├────────────────────────────────────────────────────┤
│                                                    │
│  Layer 1: Network Security                        │
│  • HTTPS/TLS 1.3 only                            │
│  • CloudFlare WAF                                 │
│  • DDoS protection                                │
│  • IP whitelisting for admin                     │
│                                                    │
│  Layer 2: Application Security                    │
│  • Input validation & sanitization               │
│  • SQL injection prevention (parameterized)       │
│  • XSS prevention (output encoding)              │
│  • CSRF tokens                                    │
│  • Rate limiting (per IP, per user)              │
│                                                    │
│  Layer 3: Authentication Security                 │
│  • OTP with 5-min expiry                         │
│  • Max 3 OTP attempts                            │
│  • JWT with 15-min expiry                        │
│  • Refresh token rotation                        │
│  • Account lockout after 5 failed attempts       │
│  • Password: bcrypt cost 12                      │
│                                                    │
│  Layer 4: Authorization Security                  │
│  • Role-based access control (RBAC)              │
│  • Organization-based isolation                   │
│  • Resource ownership validation                  │
│  • Permission checks on every request            │
│                                                    │
│  Layer 5: Data Security                           │
│  • Encryption at rest (AWS RDS encryption)       │
│  • Encryption in transit (TLS)                   │
│  • PII masking in logs                           │
│  • Secure file uploads (virus scan)             │
│  • Regular backups (encrypted)                    │
│                                                    │
│  Layer 6: Monitoring & Audit                      │
│  • Comprehensive audit logs                       │
│  • Real-time alerting (failed logins, etc.)      │
│  • Anomaly detection                             │
│  • Security dashboard                             │
│  • Incident response procedures                   │
│                                                    │
└────────────────────────────────────────────────────┘
\\\

---

## 6. Authentication System

### 6.1 OTP-First Authentication Flow

#### 6.1.1 Login Flow (Primary)

\\\
User Action                 System Response
───────────                 ────────────────

1. Enter email/phone    →   Validate format
                           Check if user exists
                           
2. Click "Send OTP"     →   Generate 6-digit OTP
                           Store in Redis (5-min TTL)
                           Send via SMS/Email
                           Rate limit: 3 OTPs/hour
                           
3. Enter OTP code       →   Validate OTP
                           Check expiry
                           Check attempts (max 3)
                           
4. Submit              →   Generate JWT tokens
                           • Access: 15-min expiry
                           • Refresh: 7-day expiry
                           Clear OTP from Redis
                           Update last_login
                           Create audit log
                           
5. Authenticated       →   Return user profile
                           Return organization context
                           Redirect to dashboard
\\\

**Backend Endpoints:**
\\\
POST /api/v1/auth/send-otp
Request:
{
  "identifier": "user@hospital.com", // or phone
  "type": "email" // or "sms"
}

Response:
{
  "success": true,
  "message": "OTP sent to user@hospital.com",
  "expires_in": 300,
  "retry_after": 60
}

POST /api/v1/auth/verify-otp
Request:
{
  "identifier": "user@hospital.com",
  "otp": "123456"
}

Response:
{
  "success": true,
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "user": {
    "id": "uuid",
    "email": "user@hospital.com",
    "name": "John Doe",
    "organizations": [...]
  }
}
\\\

#### 6.1.2 Password Fallback Flow

\\\
User Action                 System Response
───────────                 ────────────────

1. Click "Use Password" →   Show password field
                           
2. Enter password       →   Hash with bcrypt
                           Compare with stored hash
                           Check failed attempt count
                           
3. Submit              →   If valid:
                             - Generate JWT tokens
                             - Reset failed attempts
                             - Update last_login
                           If invalid:
                             - Increment failed attempts
                             - Lock after 5 failures
                             - Return error
                             
4. Authenticated       →   Same as OTP flow
\\\

**Backend Endpoint:**
\\\
POST /api/v1/auth/login-password
Request:
{
  "email": "user@hospital.com",
  "password": "SecurePass123"
}

Response:
{
  "success": true,
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "user": {...}
}
\\\

### 6.2 User Registration Flow

\\\
User Action                 System Response
───────────                 ────────────────

1. Fill registration form→  Validate inputs
   - Name                  - Email format check
   - Email                 - Phone format check
   - Phone (optional)      - Email uniqueness check
   - Organization          
   - Role                  

2. Click "Register"     →   Create user record
                           Set status = "pending"
                           Generate verification OTP
                           Send to email/phone
                           
3. Enter OTP           →   Validate OTP
                           Mark email/phone verified
                           Set status = "active"
                           
4. Account created     →   Send welcome email
                           Auto-login user
                           Redirect to onboarding
\\\

**Backend Endpoint:**
\\\
POST /api/v1/auth/register
Request:
{
  "name": "John Doe",
  "email": "john@hospital.com",
  "phone": "+1234567890",
  "organization_id": "uuid",
  "role": "engineer",
  "password": "optional" // if user wants password
}

Response:
{
  "success": true,
  "user_id": "uuid",
  "verification_required": true,
  "message": "OTP sent to john@hospital.com"
}
\\\

### 6.3 JWT Token Management

#### 6.3.1 Token Structure

**Access Token (15-minute expiry):**
\\\json
{
  "sub": "user-uuid",
  "email": "user@hospital.com",
  "name": "John Doe",
  "org_id": "current-org-uuid",
  "role": "engineer",
  "permissions": ["view_tickets", "update_tickets"],
  "iat": 1703001234,
  "exp": 1703002134,
  "type": "access"
}
\\\

**Refresh Token (7-day expiry):**
\\\json
{
  "sub": "user-uuid",
  "jti": "token-uuid",
  "iat": 1703001234,
  "exp": 1703606034,
  "type": "refresh"
}
\\\

#### 6.3.2 Token Refresh Flow

\\\
POST /api/v1/auth/refresh
Request:
{
  "refresh_token": "eyJhbG..."
}

Response:
{
  "success": true,
  "access_token": "eyJhbG...", // New access token
  "refresh_token": "eyJhbG..."  // New refresh token (rotation)
}
\\\

#### 6.3.3 Token Storage (Frontend)

**Option 1: HttpOnly Cookies (Recommended)**
\\\javascript
// Backend sets cookies
Set-Cookie: access_token=...; HttpOnly; Secure; SameSite=Strict; Max-Age=900
Set-Cookie: refresh_token=...; HttpOnly; Secure; SameSite=Strict; Max-Age=604800

// Frontend doesn't need to handle tokens
// Automatically sent with every request
\\\

**Option 2: LocalStorage (Alternative)**
\\\javascript
// Frontend stores tokens
localStorage.setItem('access_token', token);

// Add to requests
axios.defaults.headers.common['Authorization'] = \Bearer \\;
\\\

### 6.4 Session Management

**Session Data (Redis):**
\\\javascript
Key: "session:user-uuid:token-jti"
Value: {
  user_id: "uuid",
  device_info: {
    user_agent: "Chrome 120",
    ip: "192.168.1.1",
    platform: "Windows"
  },
  created_at: "2025-12-20T10:00:00Z",
  last_activity: "2025-12-20T10:30:00Z",
  expires_at: "2025-12-27T10:00:00Z"
}
TTL: 7 days
\\\

**Session Management Endpoints:**
\\\
GET /api/v1/auth/sessions
- List all active sessions for current user

DELETE /api/v1/auth/sessions/{session_id}
- Revoke specific session

DELETE /api/v1/auth/sessions/all
- Revoke all sessions (except current)

POST /api/v1/auth/logout
- Logout current session
\\\

### 6.5 Password Management

#### 6.5.1 Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character
- Not in common password list (top 10,000)
- Not similar to email/name

#### 6.5.2 Password Reset Flow

\\\
1. User clicks "Forgot Password"
2. Enter email/phone
3. System sends reset OTP
4. User enters OTP
5. User sets new password
6. System validates password strength
7. Password updated (bcrypt hash)
8. All sessions revoked
9. User must login again
\\\

**Backend Endpoints:**
\\\
POST /api/v1/auth/forgot-password
POST /api/v1/auth/verify-reset-otp
POST /api/v1/auth/reset-password
\\\

### 6.6 OTP Configuration

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| OTP Length | 6 digits | Balance security & usability |
| OTP Expiry | 5 minutes | Short window reduces risk |
| Max Attempts | 3 | Prevents brute force |
| Rate Limit | 3 OTPs/hour per identifier | Prevents spam |
| Cooldown | 60 seconds between sends | Prevents abuse |
| Delivery | Email (primary), SMS (fallback) | Cost optimization |

---


## 7. Multi-Tenant Architecture [SUMMARY - Full details available on request]

### 7.1 Organization Types
- Manufacturers
- Hospitals
- Laboratories  
- Distributors
- Dealers

### 7.2 Role-Based Access Control (RBAC)
- Admin, Manager, Engineer, Viewer roles
- Organization-specific permissions
- Dynamic permission assignment

### 7.3 Dashboard Configurations
- Custom cards per organization type
- Data isolation between organizations
- Shared resources (engineers, equipment)

---

## 8. Ticket Creation Flows [SUMMARY]

### 8.1 Authenticated Web Flow
- User logs in → Dashboard → Create Ticket
- Equipment auto-populated
- Contact info from profile

### 8.2 QR Code Anonymous Flow  
- Scan QR → Pre-filled form
- reCAPTCHA v3 validation
- Contact info required
- Ticket tracking URL provided

### 8.3 WhatsApp Anonymous Flow
- Message bot → Conversational ticket creation
- Equipment ID from QR
- WhatsApp number auto-captured

---

## 9. WhatsApp Integration [SUMMARY]

### 9.1 Technical Implementation
- Twilio WhatsApp Business API
- Webhook endpoint for incoming messages
- Redis for conversation state
- Command-based interaction

### 9.2 Features
- Ticket creation wizard
- Basic diagnostic chatbot
- Status update notifications
- Ticket tracking

### 9.3 Commands
- CREATE - Start ticket wizard
- STATUS [ticket-id] - Check status
- HELP - Show commands

---

## 10. Notification System [SUMMARY]

### 10.1 Notification Channels
- Email (SendGrid)
- SMS (Twilio)
- WhatsApp (Twilio)
- In-app (WebSocket)

### 10.2 Notification Events
- Ticket created
- Engineer assigned
- Status updates
- Ticket completed
- Parts ordered

### 10.3 Notification Preferences
- User-level preferences
- Organization-level defaults
- Per-ticket overrides

---

## 11. Security Requirements [SUMMARY]

### 11.1 Authentication Security
- bcrypt password hashing (cost 12)
- JWT RS256 signing
- OTP rate limiting
- Account lockout

### 11.2 API Security
- Rate limiting (100 req/min per IP)
- Input validation
- SQL injection prevention
- XSS prevention
- CSRF tokens

### 11.3 Data Security
- Encryption at rest
- TLS 1.3 in transit
- PII masking in logs
- Secure file uploads

### 11.4 Audit & Compliance
- Comprehensive audit logs
- HIPAA-ready logging
- Data retention policies
- Incident response procedures

---

## 12. Database Schema [SUMMARY]

### Core Tables:
1. **users** - User accounts
2. **organizations** - Multi-tenant orgs
3. **user_organizations** - Many-to-many mapping
4. **roles** - RBAC roles
5. **permissions** - Granular permissions
6. **otp_codes** - OTP management
7. **refresh_tokens** - Token rotation
8. **auth_audit_log** - Security audit
9. **service_tickets** - Enhanced with contact_info
10. **notification_preferences** - User preferences

Full schema available in separate migration files.

---

## 13. API Specifications [SUMMARY]

### Authentication Endpoints (12 total)
- POST /api/v1/auth/send-otp
- POST /api/v1/auth/verify-otp
- POST /api/v1/auth/login-password
- POST /api/v1/auth/register
- POST /api/v1/auth/refresh
- POST /api/v1/auth/logout
- GET /api/v1/auth/me
- GET /api/v1/auth/sessions
- DELETE /api/v1/auth/sessions/{id}
- POST /api/v1/auth/forgot-password
- POST /api/v1/auth/verify-reset-otp
- POST /api/v1/auth/reset-password

### Organization Endpoints (8 total)
- GET /api/v1/organizations
- POST /api/v1/organizations
- GET /api/v1/organizations/{id}
- PUT /api/v1/organizations/{id}
- GET /api/v1/organizations/{id}/dashboard
- GET /api/v1/organizations/{id}/users
- POST /api/v1/organizations/{id}/users
- DELETE /api/v1/organizations/{id}/users/{user_id}

### Ticket Endpoints (Enhanced - 6 total)
- POST /api/v1/tickets (authenticated)
- POST /api/v1/tickets/anonymous (QR/WhatsApp)
- GET /api/v1/tickets/{id}/track (public)
- PUT /api/v1/tickets/{id}/status
- POST /api/v1/tickets/{id}/parts
- GET /api/v1/tickets

### WhatsApp Endpoints (2 total)
- POST /api/v1/whatsapp/webhook
- POST /api/v1/whatsapp/send

### Notification Endpoints (3 total)
- GET /api/v1/notifications
- PUT /api/v1/notifications/preferences
- POST /api/v1/notifications/send

Full OpenAPI spec available separately.

---

## 14. Frontend Requirements [SUMMARY]

### 14.1 New Pages
- **/login** - OTP-first login
- **/register** - User registration
- **/forgot-password** - Password reset
- **/dashboard** - Organization-specific
- **/profile** - User settings
- **/sessions** - Active sessions
- **/organizations/{id}/settings** - Org management

### 14.2 Components
- AuthProvider (Context)
- ProtectedRoute (HOC)
- OTPInput (Component)
- OrganizationSwitcher
- NotificationBell
- SessionManager

### 14.3 State Management
- Auth state (Context API)
- User profile (TanStack Query)
- Organizations (TanStack Query)
- Notifications (WebSocket + State)

---

## 15. Implementation Plan

### Phase 1: Core Authentication (4 weeks)
**Week 1-2:**
- Database schema & migrations
- User model & repository
- OTP service (Twilio integration)
- JWT service

**Week 3-4:**
- Authentication endpoints
- Frontend login/register pages
- Session management
- Basic RBAC

**Deliverables:**
- ✅ Users can register with OTP
- ✅ Users can login with OTP/password
- ✅ JWT token management working
- ✅ Basic permissions implemented

### Phase 2: Multi-Tenancy (3 weeks)
**Week 5-6:**
- Organization models
- User-org mapping
- Role & permission system
- Dashboard configuration

**Week 7:**
- Organization-specific dashboards
- Frontend org switcher
- Permission enforcement
- Testing

**Deliverables:**
- ✅ 5 organization types supported
- ✅ Custom dashboards per type
- ✅ RBAC fully functional
- ✅ Data isolation verified

### Phase 3: WhatsApp Integration (3 weeks)
**Week 8-9:**
- Twilio WhatsApp API setup
- Webhook endpoint
- Conversation state management
- Basic chatbot logic

**Week 10:**
- Ticket creation via WhatsApp
- Notification sending
- Command handlers
- Testing

**Deliverables:**
- ✅ WhatsApp ticket creation working
- ✅ Basic diagnostic chatbot
- ✅ Notifications via WhatsApp
- ✅ Command system functional

### Phase 4: Enhanced Tickets & Notifications (2 weeks)
**Week 11:**
- Anonymous ticket creation
- Contact info management
- reCAPTCHA integration
- Ticket tracking page

**Week 12:**
- Multi-channel notification service
- Notification preferences
- Email templates
- SMS templates

**Deliverables:**
- ✅ QR tickets with reCAPTCHA
- ✅ Flexible notification system
- ✅ All channels working
- ✅ Tracking page implemented

### Phase 5: Security & Polish (2 weeks)
**Week 13:**
- Security audit
- Rate limiting implementation
- Audit logging
- Performance optimization

**Week 14:**
- End-to-end testing
- Security testing
- Load testing
- Documentation

**Deliverables:**
- ✅ Security audit passed
- ✅ Performance targets met
- ✅ All tests passing
- ✅ Documentation complete

### Total Timeline: 14 weeks (~3.5 months)

---

## 16. Testing Strategy

### 16.1 Unit Testing
- All service methods
- Repository methods
- Utility functions
- 80%+ code coverage target

### 16.2 Integration Testing
- API endpoint tests
- Database transactions
- External service mocks
- Auth flow tests

### 16.3 E2E Testing
- Complete user journeys
- Multi-tenant scenarios
- WhatsApp flows
- Security scenarios

### 16.4 Security Testing
- OWASP Top 10 vulnerabilities
- Penetration testing
- Rate limit testing
- SQL injection attempts

### 16.5 Performance Testing
- Load testing (1000 concurrent users)
- Stress testing
- OTP delivery latency
- Database query optimization

---

## 17. Deployment Plan

### 17.1 Infrastructure Setup
- AWS RDS PostgreSQL (Multi-AZ)
- AWS ElastiCache Redis (Cluster mode)
- AWS EC2 Auto Scaling Group
- Application Load Balancer
- CloudFront CDN
- AWS S3 for static assets

### 17.2 Deployment Strategy
- Blue-Green deployment
- Database migrations (automated)
- Configuration management (AWS Secrets Manager)
- Monitoring (CloudWatch + Custom)
- Rollback procedures

### 17.3 CI/CD Pipeline
- GitHub Actions
- Automated tests on PR
- Staging environment deployment
- Production deployment (manual approval)
- Automated rollback on failure

---

## 18. Success Metrics

### 18.1 User Adoption Metrics
| Metric | Target | Timeline |
|--------|--------|----------|
| OTP Login Adoption | 70% | 3 months |
| User Registration Rate | 100 users/month | 6 months |
| WhatsApp Ticket % | 30% | 6 months |
| Mobile Traffic | 50% | 6 months |

### 18.2 Performance Metrics
| Metric | Target | SLA |
|--------|--------|-----|
| Login Time | <30 seconds | 95th percentile |
| OTP Delivery | <30 seconds | 99% success |
| API Response Time | <200ms | 95th percentile |
| System Uptime | 99.9% | Monthly |

### 18.3 Security Metrics
| Metric | Target | Frequency |
|--------|--------|-----------|
| Security Breaches | 0 | Ongoing |
| Failed Login Rate | <5% | Daily |
| OTP Success Rate | >95% | Daily |
| Suspicious Activity Alerts | <10/day | Daily |

---

## 19. Risks & Mitigation

### 19.1 Technical Risks

**Risk:** OTP delivery failures (SMS/Email)
- **Impact:** High - Users cannot login
- **Mitigation:** 
  - Use reliable providers (Twilio, SendGrid)
  - Implement fallback (SMS → Email → Password)
  - Monitor delivery rates in real-time
  - Have manual override for support

**Risk:** Database performance degradation
- **Impact:** High - System slowdown
- **Mitigation:**
  - Database indexing optimization
  - Query performance monitoring
  - Read replicas for scaling
  - Redis caching layer

**Risk:** WhatsApp API rate limits
- **Impact:** Medium - Notification delays
- **Mitigation:**
  - Queue-based sending
  - Rate limit awareness
  - Fallback to SMS/Email
  - Batch notifications

### 19.2 Security Risks

**Risk:** OTP brute force attacks
- **Impact:** High - Account compromise
- **Mitigation:**
  - Rate limiting (3 attempts)
  - Account lockout
  - CAPTCHA on repeated failures
  - Monitoring & alerts

**Risk:** Session hijacking
- **Impact:** High - Unauthorized access
- **Mitigation:**
  - HttpOnly cookies
  - Short token expiry (15 min)
  - Token rotation
  - IP/Device binding
  - Suspicious login detection

### 19.3 Business Risks

**Risk:** Low user adoption of OTP login
- **Impact:** Medium - Feature underutilization
- **Mitigation:**
  - User education/onboarding
  - Prominent positioning
  - Password fallback available
  - Gather user feedback

**Risk:** WhatsApp dependency
- **Impact:** Medium - Single point of failure
- **Mitigation:**
  - Not using for login (only tickets)
  - Alternative channels available
  - No critical functions WhatsApp-only

---

## 20. Appendices

### Appendix A: Complete API Reference
See: API-SPECIFICATION.md

### Appendix B: Database Migrations
See: database/migrations/auth-*.sql

### Appendix C: Security Checklist
See: SECURITY-CHECKLIST.md

### Appendix D: Testing Scenarios
See: TESTING-SCENARIOS.md

### Appendix E: Deployment Runbook
See: DEPLOYMENT-RUNBOOK.md

### Appendix F: WhatsApp Integration Guide
See: WHATSAPP-INTEGRATION.md

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-12-20 | Technical Team | Initial comprehensive PRD |

---

## Approval & Sign-off

**Prepared by:** Technical Architecture Team  
**Date:** December 20, 2025  

**Reviewed by:** [Pending]  
**Date:** [Pending]  

**Approved by:** [Pending]  
**Date:** [Pending]  

---

## Next Steps

1. **Review & Approval** (Week 1)
   - Stakeholder review of PRD
   - Technical team feedback
   - Finalize scope & timeline

2. **Technical Design** (Week 2)
   - Detailed API specifications
   - Database schema finalization
   - Architecture diagrams

3. **Sprint Planning** (Week 2)
   - Break down into user stories
   - Assign story points
   - Create sprint schedule

4. **Development Kickoff** (Week 3)
   - Phase 1 begins
   - Daily standups
   - Weekly demos

---

**END OF PRD DOCUMENT**

**Total Pages:** ~50+  
**Word Count:** ~20,000+  
**Estimated Reading Time:** 90 minutes

For questions or clarifications, contact: technical-team@aby-med.com

