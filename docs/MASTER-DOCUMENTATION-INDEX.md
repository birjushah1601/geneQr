# Master Documentation Index  
*(Last updated: 19 Nov 2024)*  

This index is the **single authoritative map** of project documents.  
– “OFFICIAL” files are current and must be used by every team.  
– “ARCHIVED / OBSOLETE” files are retained only for historical reference and **must not** be used for engineering or business decisions.

---

## 1  OFFICIAL DOCUMENTS (USE THESE)

| File | Purpose |
|------|---------|
| [FINAL-intelligent-medical-platform-prd.md](./FINAL-intelligent-medical-platform-prd.md) | Master Product Requirements Document – business model, AI–advisory vision, user journeys, success metrics |
| [FINAL-technical-implementation-guide.md](./FINAL-technical-implementation-guide.md) | Engineering handbook – micro-services, APIs, DB schemas, ML pipeline, security & DevOps specs |
| [FINAL-project-roadmap-and-execution.md](./FINAL-project-roadmap-and-execution.md) | 24-month execution plan – phases, resources, budget, GTM, risk management |
| [platform-engineering.md](./platform-engineering.md) | Platform and monorepo engineering overview (Makefile, modules, CI/CD, observability) |
| [deployment.md](./deployment.md) | Deployment and environment runbooks (local, staging, prod) |
| [qa.md](./qa.md) | QA strategy, test flows, and acceptance criteria |
| [postman-verification.md](./postman-verification.md) | API verification with Postman; includes collection and steps |

> Always start with these three files. Any new specs, epics or designs **must align** with them.

---

## 2  ARCHIVED / OBSOLETE DOCUMENTS (DO NOT USE)

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

---

## 3  SERVICE DOCUMENTATION

Each service has a concise overview and API surface:

- [services/catalog.md](./services/catalog.md)
- [services/rfq.md](./services/rfq.md)
- [services/quote.md](./services/quote.md)
- [services/comparison.md](./services/comparison.md)
- [services/contract.md](./services/contract.md)
- [services/supplier.md](./services/supplier.md)
- [services/equipment-registry.md](./services/equipment-registry.md)
- [services/service-ticket.md](./services/service-ticket.md)

## 4  API DOCUMENTATION

Complete API reference documents for all services:

- [api/ATTACHMENT-API.md](./api/ATTACHMENT-API.md) - 📎 **Complete Attachment System API** (4 endpoints, frontend integration, database schema, production-ready)
- [api/ASSIGNMENT-API.md](./api/ASSIGNMENT-API.md) - 🎯 AI-Powered Smart Assignment API

## 5  DATABASE DOCUMENTATION

Per-table docs with fields, relationships, and improvement suggestions live under `docs/database/`.
Start here: [database/README.md](./database/README.md)

---

## 6  HOW TO USE THIS INDEX

1. **Read the OFFICIAL documents first.**  
2. When creating new tickets, user stories or designs, reference the **section/heading** in the relevant OFFICIAL file.  
3. If you find conflicting information, default to OFFICIAL docs and raise an issue with Product & Architecture leads.  
4. Archived files may be browsed for historical context but **must not** drive decisions.  

---

## 7  CHANGE CONTROL

Any modification to OFFICIAL documents requires:  
1. Product & Architecture sign-off.  
2. Version bump in filename (e.g., `FINAL-technical-implementation-guide_v1.1.md`).  
3. Update of this index within the same pull request.

---

*End of Master Index*
