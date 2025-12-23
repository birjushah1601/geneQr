# ABY-MED Platform Architecture

## 🏗️ System Overview

ABY-MED is a **modular, multi-tenant medical equipment service management platform** built with modern microservices patterns, clean architecture, and domain-driven design principles.

### Core Principles
- **Multi-tenancy:** Complete data isolation per organization
- **Modularity:** Pluggable modules with clear boundaries
- **Scalability:** Horizontal scaling support
- **Security:** Role-based access, audit logging, rate limiting
- **AI-First:** Integrated AI for diagnostics and automation

---

## 🎯 Technology Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** Chi router (lightweight, fast)
- **Database:** PostgreSQL 15+ (primary), Redis (caching)
- **ORM/Driver:** pgx/v5 (native PostgreSQL driver)
- **Authentication:** JWT tokens, session-based
- **API:** RESTful JSON APIs

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript 5+
- **UI Library:** React 18
- **Styling:** Tailwind CSS 3+, shadcn/ui components
- **State:** React Query (server state), Context API
- **Forms:** React Hook Form, Zod validation

### Infrastructure
- **Containerization:** Docker, Docker Compose
- **Database:** PostgreSQL container
- **Caching:** Redis (optional)
- **Storage:** Local filesystem / S3-compatible
- **Email:** SendGrid
- **AI:** OpenAI (GPT-4, Whisper), Anthropic (Claude 3)

---

## 📐 System Architecture

### High-Level Architecture
```
┌──────────────────────────────────────────────────────────────────┐
│                         Frontend Layer                           │
│  Next.js 14 App Router | React 18 | TypeScript | Tailwind       │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                │
│  │ Dashboard  │  │  Tickets   │  │ Equipment  │  [+ 10 pages]  │
│  └────────────┘  └────────────┘  └────────────┘                │
└────────────────────────┬─────────────────────────────────────────┘
                         │ HTTP/REST (Port 3000)
                         ├─────────────────────────────────────────┐
                         ↓                                         │
┌──────────────────────────────────────────────────────────────────┤
│                        API Gateway Layer                         │
│  Chi Router | Middleware (Auth, CORS, Rate Limiting)            │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  /api/v1/tickets    /api/v1/equipment                    │   │
│  │  /api/v1/organizations   /api/v1/engineers               │   │
│  │  /api/v1/diagnosis  /api/v1/whatsapp                     │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────────────────────┬─────────────────────────────────────────┘
                         │
┌────────────────────────┴─────────────────────────────────────────┐
│                    Application Layer (Go)                        │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                   Module Architecture                     │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │   │
│  │  │ Tickets  │ │Equipment │ │  Orgs    │ │Engineers │   │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │   │
│  │  │ WhatsApp │ │   AI     │ │  Parts   │ │Marketplace│  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
│  Infrastructure: Email, Reports, Audit, Notifications           │
└────────────────────────┬─────────────────────────────────────────┘
                         │
┌────────────────────────┴─────────────────────────────────────────┐
│                      Data Layer                                  │
│  PostgreSQL 15+ (Multi-tenant)  │  Redis (Caching)               │
│  40+ Tables | JSONB | FTS       │  Sessions | Rate Limiting      │
└──────────────────────────────────────────────────────────────────┘
```

---

## 🧩 Module Architecture

### Clean Architecture Pattern
Each module follows clean architecture with clear separation:

```
module/
├── domain/          # Business entities, interfaces, rules
│   ├── entities.go  # Domain models (Equipment, Ticket, etc.)
│   ├── repository.go# Repository interfaces
│   └── services.go  # Domain services
├── app/             # Application use cases
│   ├── service.go   # Business logic orchestration
│   └── dto.go       # Data transfer objects
├── api/             # HTTP handlers
│   └── handler.go   # REST API endpoints
├── infra/           # Infrastructure implementations
│   └── postgres_repository.go # Database implementation
└── module.go        # Module initialization & routing
```

### Module List (8 Core Modules)

1. **Service Ticket** - Ticket lifecycle management
2. **Equipment Registry** - Equipment tracking, QR codes
3. **Organizations** - Multi-tenant organization management
4. **Engineers** - Field engineer management, assignment
5. **WhatsApp** - Message handling, ticket creation
6. **AI Diagnosis** - Intelligent diagnostics
7. **Parts** - Spare parts catalog
8. **Marketplace** - E-commerce (coming soon)

---

## 🗄️ Database Architecture

### Multi-Tenant Design
```
┌─────────────────────────────────────────────────────────────┐
│                    Organizations Table                      │
│  (Tenant Root)                                              │
│  - id (UUID)                                                │
│  - name, type, status                                       │
└────────────┬────────────────────────────────────────────────┘
             │
       ┌─────┴──────┬──────────┬──────────┐
       │            │          │          │
   ┌───▼────┐  ┌───▼────┐ ┌──▼────┐ ┌──▼────┐
   │Tickets │  │Equipment│ │Engineers│ │ Parts│
   │        │  │         │ │        │ │       │
   │org_id  │  │org_id   │ │org_id  │ │org_id │
   └────────┘  └─────────┘ └────────┘ └───────┘
```

### Key Tables (40+)

**Core Entities:**
- `organizations` - Tenant root
- `users` - System users
- `equipment_registry` - Equipment tracking
- `service_tickets` - Ticket management
- `engineers` - Field engineers
- `spare_parts_catalog` - Parts inventory

**Relationships:**
- `ticket_parts` - Parts per ticket
- `ticket_status_history` - Audit trail
- `ticket_comments` - Communication
- `engineer_equipment_types` - Capabilities
- `equipment_service_config` - Service SLAs

**System:**
- `audit_logs` - All user actions
- `qr_codes`, `qr_batches` - QR management
- `attachments` - File uploads
- `notifications` - Email/SMS queue

### Multi-Tenancy Implementation
```sql
-- Every tenant-scoped table has org_id
CREATE TABLE service_tickets (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    ticket_number VARCHAR(50) UNIQUE,
    ...
);

-- Row-level security (planned)
ALTER TABLE service_tickets ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON service_tickets
    USING (organization_id = current_setting('app.current_org_id')::UUID);
```

---

## 🔐 Security Architecture

### Authentication & Authorization
```
User Login
    ↓
JWT Token Generation (24h expiry)
    ↓
Token stored in httpOnly cookie
    ↓
Every API request → Middleware validates token
    ↓
Extract user_id, org_id, role
    ↓
Set in request context
    ↓
Business logic uses context for authorization
```

### Security Layers

1. **Transport Security**
   - HTTPS in production
   - CORS policy enforcement
   - Secure cookie flags

2. **Input Validation**
   - Request size limits (10MB)
   - Input sanitization middleware
   - SQL injection prevention (prepared statements)
   - XSS protection

3. **Rate Limiting**
   - IP-based: 20 tickets/hour
   - QR-based: 5 tickets/hour per QR code
   - API-level: 100 requests/minute per user

4. **Audit Logging**
   - Every CREATE/UPDATE/DELETE logged
   - User, IP, timestamp, changes tracked
   - Immutable audit trail

5. **Data Isolation**
   - org_id filtering on all queries
   - No cross-tenant data access
   - API responses filtered by tenant

---

## 🔄 Data Flow Examples

### Ticket Creation Flow
```
User (Frontend)
    │
    │ POST /api/v1/tickets
    │ { equipment_id, description, ... }
    ↓
API Handler (Middleware Chain)
    │ 1. CORS Check
    │ 2. Auth Check (JWT)
    │ 3. Rate Limiting
    │ 4. Input Sanitization
    ↓
Ticket Service (Business Logic)
    │ 1. Validate equipment exists
    │ 2. Generate ticket number
    │ 3. Set default priority
    │ 4. Calculate SLA deadlines
    ↓
Repository (Database)
    │ 1. BEGIN TRANSACTION
    │ 2. INSERT ticket
    │ 3. INSERT status_history
    │ 4. INSERT audit_log
    │ 5. COMMIT
    ↓
Async Tasks (Background)
    │ 1. Send email notification
    │ 2. Queue for engineer assignment
    │ 3. Update analytics
    ↓
Response
    │ 201 Created
    │ { ticket_id, ticket_number, ... }
    ↓
Frontend
    │ Show success message
    │ Navigate to ticket details
```

### WhatsApp Message → Ticket
```
WhatsApp (User sends message with QR code)
    ↓
Webhook Handler
    │ Parse message
    │ Extract QR code
    │ Extract issue description
    ↓
Equipment Service
    │ Lookup equipment by QR
    ↓
Ticket Service
    │ Create ticket
    │ source = "whatsapp"
    ↓
If audio message:
    │ Download audio file
    │ Call Whisper API
    │ Transcribe to text
    │ Attach audio + transcript
    ↓
Confirmation
    │ Send WhatsApp message back
    │ "Ticket #TKT-20251223-001 created"
```

---

## 🚀 Deployment Architecture

### Development Environment
```
Developer Machine
├── Backend: go run cmd/platform/main.go (port 8081)
├── Frontend: npm run dev (port 3000)
└── Database: Docker PostgreSQL (port 5430)
```

### Production Environment (Proposed)
```
Load Balancer (NGINX/Caddy)
    │
    ├─→ Frontend Server(s) (Next.js)
    │   └─ Static assets (CDN)
    │
    ├─→ Backend Server(s) (Go binary)
    │   ├─ Auto-scaling (CPU/Memory based)
    │   └─ Health checks (/health)
    │
    ├─→ Database (PostgreSQL)
    │   ├─ Primary (Read/Write)
    │   └─ Replica(s) (Read-only)
    │
    └─→ Cache (Redis)
        └─ Sessions, rate limiting
```

---

## 📊 Performance Considerations

### Database Optimization
- **Indexes:** All foreign keys, query columns
- **JSONB:** For flexible metadata storage
- **Full-text search:** PostgreSQL FTS on descriptions
- **Connection pooling:** Max 10 connections
- **Query timeout:** 10 seconds

### Caching Strategy
- **Redis:** Session storage, rate limit counters
- **In-memory:** Configuration, feature flags
- **CDN:** Static assets, images

### Scalability
- **Horizontal:** Multiple backend instances
- **Vertical:** Database replica for reads
- **Async:** Background jobs for email, reports
- **Queue:** Future: RabbitMQ/Kafka for events

---

## 🔌 Integration Points

### External Services
- **OpenAI:** GPT-4 (diagnosis), Whisper (audio→text)
- **Anthropic:** Claude 3 (diagnosis)
- **SendGrid:** Email notifications
- **Twilio:** WhatsApp Business API
- **Storage:** Local FS / S3 (future)

### Webhooks
- **WhatsApp:** Incoming message webhook
- **Payment:** (Future) Payment gateway webhooks

### APIs (Outbound)
- **AI Models:** Diagnosis, transcription
- **Email Service:** Transactional emails
- **SMS:** (Future) Status updates

---

## 📁 Project Structure

```
aby-med/
├── cmd/
│   └── platform/
│       └── main.go              # Application entry point
├── internal/
│   ├── service-domain/          # Business modules
│   │   ├── service-ticket/
│   │   ├── equipment-registry/
│   │   ├── organizations/
│   │   ├── engineers/
│   │   ├── whatsapp/
│   │   ├── attachment/
│   │   └── marketplace/ (future)
│   ├── infrastructure/          # Cross-cutting
│   │   ├── email/
│   │   ├── reports/
│   │   ├── notification/
│   │   └── audit/
│   └── shared/                  # Utilities
│       ├── middleware/
│       ├── database/
│       └── config/
├── admin-ui/                    # Frontend
│   ├── src/
│   │   ├── app/                 # Next.js pages
│   │   ├── components/          # React components
│   │   ├── lib/                 # API clients
│   │   └── types/               # TypeScript types
│   └── public/                  # Static assets
├── database/
│   └── migrations/              # SQL migrations
├── docs/                        # Documentation
├── .env.example                 # Environment template
├── go.mod                       # Go dependencies
└── docker-compose.yml           # Local dev setup
```

---

## 🎯 Design Decisions

### Why Go for Backend?
- **Performance:** Fast, compiled binary
- **Concurrency:** Goroutines for async tasks
- **Type Safety:** Strong typing prevents bugs
- **Deployment:** Single binary, no runtime needed
- **Libraries:** Excellent PostgreSQL, HTTP support

### Why Next.js for Frontend?
- **SSR/SSG:** Better SEO, faster initial load
- **TypeScript:** Type safety end-to-end
- **Developer Experience:** Hot reload, fast refresh
- **Production Ready:** Optimized builds, image optimization
- **Community:** Large ecosystem, frequent updates

### Why PostgreSQL?
- **JSONB:** Flexible schema when needed
- **Full-text Search:** Built-in, no external service
- **ACID:** Strong consistency guarantees
- **Extensions:** PostGIS (future location features)
- **Performance:** Handles millions of rows easily

### Why Modular Architecture?
- **Maintainability:** Clear boundaries, easy to understand
- **Testability:** Mock interfaces, unit test modules
- **Scalability:** Extract modules to microservices later
- **Team Work:** Multiple devs work on different modules
- **Feature Flags:** Enable/disable modules independently

---

## 🔮 Future Architecture Considerations

### Microservices Migration
```
Monolith (Current)
    ↓
Modular Monolith (Current state)
    ↓
Microservices (Future)
    ├─ Ticket Service
    ├─ Equipment Service
    ├─ AI Service
    └─ API Gateway
```

### Event-Driven Architecture
- **Message Queue:** RabbitMQ/Kafka
- **Events:** TicketCreated, EngineerAssigned, etc.
- **Consumers:** Email service, analytics, webhooks

### Advanced Features
- **GraphQL:** Flexible queries for mobile apps
- **WebSockets:** Real-time ticket updates
- **gRPC:** Inter-service communication
- **Service Mesh:** Istio for microservices

---

## 📚 Related Documentation

- **Getting Started:** [01-GETTING-STARTED.md](./01-GETTING-STARTED.md)
- **Features:** [03-FEATURES.md](./03-FEATURES.md)
- **API Reference:** [04-API-REFERENCE.md](./04-API-REFERENCE.md)
- **Deployment:** [05-DEPLOYMENT.md](./05-DEPLOYMENT.md)
- **Multi-Tenant:** [MULTI-TENANT-IMPLEMENTATION-PLAN.md](./MULTI-TENANT-IMPLEMENTATION-PLAN.md)
- **Security:** [SECURITY-IMPLEMENTATION-COMPLETE.md](./SECURITY-IMPLEMENTATION-COMPLETE.md)

---

**Last Updated:** December 23, 2025  
**Version:** 2.0  
**Status:** Production Architecture
