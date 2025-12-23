# Master Documentation Index  
*(Last updated: 13 December 2025)*  

This index is the **single authoritative map** of project documents.  
– "OFFICIAL" files are current and must be used by every team.  
– "ARCHIVED / OBSOLETE" files are retained only for historical reference and **must not** be used for engineering or business decisions.  
– "IMPLEMENTATION" files track recent development work and production features.

---

## 1 OFFICIAL DOCUMENTS (USE THESE)

| File | Purpose |
|------|---------|
| [FINAL-intelligent-medical-platform-prd.md](./FINAL-intelligent-medical-platform-prd.md) | Master Product Requirements Document – business model, AI–advisory vision, user journeys, success metrics |
| [FINAL-technical-implementation-guide.md](./FINAL-technical-implementation-guide.md) | Engineering handbook – micro-services, APIs, DB schemas, ML pipeline, security & DevOps specs |
| [FINAL-project-roadmap-and-execution.md](./FINAL-project-roadmap-and-execution.md) | 24-month execution plan – phases, resources, budget, GTM, risk management |
| [platform-engineering.md](./platform-engineering.md) | Platform and monorepo engineering overview (Makefile, modules, CI/CD, observability) |
| [deployment.md](./deployment.md) | Deployment and environment runbooks (local, staging, prod) |
| [postman-verification.md](./postman-verification.md) | API verification with Postman; includes collection and steps |

> Always start with these files. Any new specs, epics or designs **must align** with them.

---

## 2 RECENT IMPLEMENTATION SESSIONS (PRODUCTION-READY)

### December 2025 Sessions

| File | Status | Description |
|------|--------|-------------|
| [SESSION-DEC-12-2025-SUMMARY.md](./SESSION-DEC-12-2025-SUMMARY.md) | ✅ Complete | **Multi-Model Engineer Assignment System** - 5 assignment algorithms, side-by-side UI, API fixes, dashboard cleanup, comment system fixes |
| [features/MULTI-MODEL-ENGINEER-ASSIGNMENT.md](./features/MULTI-MODEL-ENGINEER-ASSIGNMENT.md) | ✅ Production | Complete feature documentation for intelligent engineer assignment with 5 models |
| [TICKETS-PARTS-INTEGRATION-COMPLETE.md](./TICKETS-PARTS-INTEGRATION-COMPLETE.md) | ✅ Complete | Parts management integrated with service tickets |
| [PARTS-MANAGEMENT-COMPLETE.md](./PARTS-MANAGEMENT-COMPLETE.md) | ✅ Complete | Spare parts catalog and management system |

### November 2025 Sessions

| File | Status | Description |
|------|--------|-------------|
| [GIT-PUSH-SUMMARY.md](./GIT-PUSH-SUMMARY.md) | ✅ Complete | Git workflow and push summary |
| [AI_INTEGRATION_STATUS.md](./AI_INTEGRATION_STATUS.md) | ✅ Complete | AI diagnosis system integration status |
| [INTEGRATION_PLAN.md](./INTEGRATION_PLAN.md) | 📋 Planning | Integration plan for various systems |

---

## 3 FEATURE DOCUMENTATION

### Field Service Management

| File | Status | Description |
|------|--------|-------------|
| [field-service-management-implementation.md](./field-service-management-implementation.md) | ✅ Complete | Core service ticket system implementation |
| [features/MULTI-MODEL-ENGINEER-ASSIGNMENT.md](./features/MULTI-MODEL-ENGINEER-ASSIGNMENT.md) | ✅ Production | **NEW**: 5 intelligent assignment algorithms with scoring |
| [ENGINEER-ASSIGNMENT-COMPLETE-WITH-POSTMAN.md](./ENGINEER-ASSIGNMENT-COMPLETE-WITH-POSTMAN.md) | ✅ Complete | Engineer assignment APIs with Postman tests |
| [ENGINEER-ASSIGNMENT-BACKEND-COMPLETE.md](./ENGINEER-ASSIGNMENT-BACKEND-COMPLETE.md) | ✅ Complete | Backend implementation for engineer assignment |
| [ENGINEER-ASSIGNMENT-TESTED-WORKING.md](./ENGINEER-ASSIGNMENT-TESTED-WORKING.md) | ✅ Verified | Testing verification for assignment system |
| [SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md](./SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md) | ⚠️ Superseded | Replaced by multi-model system |

### AI Services

| File | Status | Description |
|------|--------|-------------|
| [features/AI-SERVICE-MANAGEMENT-PLAN.md](./features/AI-SERVICE-MANAGEMENT-PLAN.md) | ✅ Complete | AI service planning and architecture |
| [features/T2C1-AI-FOUNDATION-COMPLETE.md](./features/T2C1-AI-FOUNDATION-COMPLETE.md) | ✅ Complete | AI foundation implementation |
| [features/PHASE2C-AI-SERVICES-PLAN.md](./features/PHASE2C-AI-SERVICES-PLAN.md) | 📋 Planning | Phase 2C AI services roadmap |
| [AI_ASSISTED_IMPLEMENTATION.md](./AI_ASSISTED_IMPLEMENTATION.md) | ✅ Reference | AI-assisted development practices |

### Parts Management

| File | Status | Description |
|------|--------|-------------|
| [PARTS-MANAGEMENT-COMPLETE.md](./PARTS-MANAGEMENT-COMPLETE.md) | ✅ Complete | Spare parts catalog system |
| [TICKETS-PARTS-INTEGRATION-COMPLETE.md](./TICKETS-PARTS-INTEGRATION-COMPLETE.md) | ✅ Complete | Integration of parts with tickets |

### QR Code System

| File | Status | Description |
|------|--------|-------------|
| [features/qr-code-feature.md](./features/qr-code-feature.md) | ✅ Complete | QR code generation and scanning |
| [features/qr-database-storage.md](./features/qr-database-storage.md) | ✅ Complete | QR code database schema |

---

## 4 API DOCUMENTATION

Complete API reference documents for all services:

| File | Endpoints | Status |
|------|-----------|--------|
| [api/ATTACHMENT-API.md](./api/ATTACHMENT-API.md) | 4 endpoints | ✅ Production |
| [api/ASSIGNMENT-API.md](./api/ASSIGNMENT-API.md) | AI-powered assignment | ⚠️ Superseded by multi-model |

### API Endpoints Summary (Current)

**Service Tickets:**
- `GET /api/v1/tickets` - List tickets
- `GET /api/v1/tickets/{id}` - Get ticket details
- `POST /api/v1/tickets` - Create ticket
- `GET /api/v1/tickets/{id}/assignment-suggestions` - **NEW**: Get all 5 assignment models
- `POST /api/v1/tickets/{id}/assign-engineer` - **NEW**: Assign engineer to ticket
- `POST /api/v1/tickets/{id}/comments` - Add comment (fixed with type validation)

**Engineers:**
- `GET /api/v1/engineers` - List engineers (path normalized)
- `GET /api/v1/engineers/{id}` - Get engineer details

**Attachments:**
- `POST /api/v1/attachments` - Upload attachment
- `GET /api/v1/attachments` - List attachments

**Equipment:**
- `GET /api/v1/equipment` - List equipment
- `GET /api/v1/equipment/{id}` - Get equipment details

---

## 5 DATABASE DOCUMENTATION

Per-table docs with fields, relationships, and improvement suggestions.

**Start here:** [database/README.md](./database/README.md)

### Recent Schema Changes (Dec 2025)

**Service Tickets Table:**
```sql
-- Extended engineer ID field
ALTER TABLE service_tickets 
ALTER COLUMN assigned_engineer_id TYPE VARCHAR(255);

-- Fixed NULL handling
UPDATE service_tickets SET
  severity = COALESCE(severity, ''),
  assigned_engineer_id = COALESCE(assigned_engineer_id, '');
```

**Ticket Comments Table:**
```sql
-- Enforces comment types: 'customer', 'engineer', 'internal', 'system'
CONSTRAINT comment_type_check CHECK (comment_type IN (...))
```

**Engineer Assignment Tracking:**
```sql
CREATE TABLE ticket_engineer_assignments (
    id VARCHAR(32) PRIMARY KEY,
    ticket_id VARCHAR(32) NOT NULL,
    engineer_id VARCHAR(32) NOT NULL,
    assignment_tier VARCHAR(10),
    assignment_tier_name VARCHAR(100),
    assigned_by VARCHAR(255),
    assigned_at TIMESTAMP DEFAULT NOW()
);
```

---

## 6 TESTING DOCUMENTATION

| File | Coverage | Status |
|------|----------|--------|
| [TESTING.md](./TESTING.md) | General testing strategy | ✅ Complete |
| [API-TEST-RESULTS.md](./API-TEST-RESULTS.md) | API test results | ✅ Current |
| [MIGRATION-COMPLETE-TESTING-NEXT.md](./MIGRATION-COMPLETE-TESTING-NEXT.md) | Migration testing | ✅ Complete |

---

## 7 PHASE COMPLETION DOCUMENTS

| File | Phase | Status |
|------|-------|--------|
| [PHASE_2C_COMPLETE.md](./PHASE_2C_COMPLETE.md) | Phase 2C | ✅ Complete |
| [PHASE2-ENGINEER-ASSIGNMENT-APIS-COMPLETE.md](./PHASE2-ENGINEER-ASSIGNMENT-APIS-COMPLETE.md) | Phase 2 Assignment | ✅ Complete |

---

## 8 ARCHIVED / OBSOLETE DOCUMENTS (DO NOT USE)

<details>
<summary>Click to expand archived documents</summary>

| File | Reason Archived |
|------|-----------------| 
| medical-equipment-platform-prd.md | Pre-AI draft; outdated business model |
| supplier-inventory-management-detailed.md | Superseded by FINAL technical guide |
| comprehensive-medical-platform-specification.md | Broad speculative scope; replaced by unified PRD |
| medical-equipment-business-playbook.md | Strategy notes now folded into FINAL roadmap |
| pain-point-analysis-and-solutions.md | Discovery artefact; insights integrated |
| digital-marketplace-solution-architecture.md | Early architecture; obsolete |
| service-as-platform-prd.md | Separate PRD; merged into unified PRD |
| digital-procurement-marketplace-prd.md | Separate PRD; merged into unified PRD |
| integrated-project-implementation-plan.md | Timeline superseded |
| medical-platform-implementation-plan.md | Duplicate of earlier plan |
| detailed-task-breakdown-spreadsheet.md | Granular tasks no longer aligned |
| project-execution-roadmap.md | Outdated execution view |
| intelligent-advisory-platform-enhancement.md | Ideas incorporated in FINAL PRD |
| ai-enhanced-ecosystem-architecture.md | Replaced by FINAL technical guide |
| ai-use-cases-detailed-specification.md | Details now in FINAL technical guide |
| unified-platform-prd-engineering-ready.md | Older consolidation; superseded |
| SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md | Replaced by multi-model system |

</details>

---

## 9 QUICK REFERENCE

### Current Production Features (December 2025)

✅ **Service Ticket Management** - Complete CRUD, status workflow  
✅ **Multi-Model Engineer Assignment** - 5 intelligent algorithms  
✅ **AI Diagnosis** - Computer vision + text analysis  
✅ **Attachments System** - Upload, AI analysis, storage  
✅ **Parts Management** - Catalog, assignment to tickets  
✅ **Comments System** - Type-validated, multi-source  
✅ **QR Code System** - Generation, scanning, tracking  

### Recent UI Improvements

✅ **Ticket Detail Page** - Reorganized 2-column layout  
✅ **Engineer Cards** - Simplified, clean design  
✅ **Assignment Interface** - Side-by-side master-detail  
✅ **Dashboard** - Removed all mock data  
✅ **API Paths** - Normalized to `/api/v1/` prefix  

### Recent Backend Improvements

✅ **Type Safety** - Fixed SQL type mismatches  
✅ **Null Handling** - Proper defaults and coalescing  
✅ **Validation** - Comment type constraints  
✅ **Scoring Algorithms** - 100-point multi-factor system  
✅ **Equipment Context** - Manufacturer, category extraction  
✅ **Workload Calculation** - Real-time ticket counting  

---

## 10 HOW TO USE THIS INDEX

1. **Read the OFFICIAL documents first.**  
2. **Check IMPLEMENTATION sections for recent features.**  
3. When creating new tickets, reference the **section/heading** in the relevant OFFICIAL file.  
4. If you find conflicting information, default to OFFICIAL docs and raise an issue with Product & Architecture leads.  
5. Archived files may be browsed for historical context but **must not** drive decisions.  
6. **Session summaries** document implementation decisions and can be referenced for context.

---

## 11 CHANGE CONTROL

Any modification to OFFICIAL documents requires:  
1. Product & Architecture sign-off.  
2. Version bump in filename (e.g., `FINAL-technical-implementation-guide_v1.1.md`).  
3. Update of this index within the same pull request.

**For Implementation Documents:**
- Session summaries are dated and versioned automatically
- Feature docs should be updated when features change
- API docs must be updated when endpoints change

---

## 12 CURRENT DEVELOPMENT STATUS

**Platform Version:** 2.0.0  
**Last Major Update:** December 13, 2025  
**Current Phase:** Production-ready field service management  

**Running Services:**
- Backend: `localhost:8081` (Go)
- Frontend: `localhost:3002` (Next.js)
- Database: `localhost:5432` (PostgreSQL)

**Latest Commit:** `06660392` - Multi-model engineer assignment + UI improvements

---

*End of Master Index v2.0*
