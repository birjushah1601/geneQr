 # Implementation Tracker: Multi‑Org, Catalog, Pricing (Backend‑first, Feature‑flagged)

 Owner: birju shah • Mode: backend‑first • Feature flags: default OFF

 Conventions
 - Flags: ENABLE_ORG, ENABLE_CHANNELS, ENABLE_CATALOG_PUBLISHING, ENABLE_PRICE_BOOKS, ENABLE_SERVICE_POLICIES, ENABLE_ENGINEERS, ENABLE_DUAL_WRITE, ENABLE_READ_FROM_ORG_GRAPH, ENABLE_PRICE_RESOLUTION, ENABLE_RESP_ORG_ASSIGNMENT
 - Scope now: backend only (no frontend changes yet)
 - Safety: additive migrations (nullable), dual‑write optional, orgs optional
 - Security posture: least‑privilege DB roles, audit logging, PII classification, migration rollback plans

 M03 — Phase 1: Orgs + Relationships + Catalog Core (read‑only)
 - Objectives
   - Add: organizations, org_relationships, channels, products, skus (nullable FKs)
   - APIs (read): GET /orgs, /orgs/{id}/relationships, /channels, /products, /skus
   - Flags: ENABLE_ORG, ENABLE_CHANNELS (read path only)
 - Tasks
   - Schema: create tables; indexes; RLS feasibility check (postpone if risky)
   - Backfill prep: crosswalk plan for existing manufacturers/suppliers
   - Security: DB users/roles; secrets sourcing; audit tables baseline
   - Tests/seed: minimal seed for demo; read‑API tests
   - Docs: ERD sketch; API spec; runbook entries
 - Deliverable: tag milestone-03

 Progress (2025‑10‑08)
 - DONE: ENABLE_ORG module skeleton with EnsureOrgSchema (organizations, org_relationships)
 - DONE: Read APIs: GET /orgs, /orgs/{id}/relationships
 - DONE: ENABLE_CHANNELS read endpoints + schema for channels/products/skus
- DONE: Seed (ENABLE_ORG_SEED) for demo orgs/channels/products/skus
- NEXT: Unit tests; backfill plan draft; update runbook

 M04 — Phase 2: Offerings + Channel Catalog (publish flow)
 - Objectives
   - Add: offerings; channel_catalog (list/unlist); draft→published versioning
   - APIs: GET/POST offerings; GET/POST channels/{id}/catalog
   - Flags: ENABLE_CATALOG_PUBLISHING (write path behind flag)
 - Tasks
   - Schema: offering, channel_catalog, publish audit
   - Services: publish/unpublish with validation; immutable published snapshots
   - Observability: metrics for listings, failures; audit trail entries
   - Tests: publish workflow, versioning guarantees
   - Docs: catalog lifecycle; rollback (new version only)
 - Deliverable: tag milestone-04

 Progress (2025‑10‑08)
 - DONE: Schema for offerings + channel_catalog (versioned)
 - DONE: APIs under ENABLE_CATALOG_PUBLISHING
   - GET /offerings, POST /offerings
   - POST /channels/{id}/catalog/publish, POST /channels/{id}/catalog/unlist
 - NEXT: Publish audit trail; tests; metrics
 - Added: basic audit logs on publish/unlist endpoints

 M05 — Phase 3: Price Books + Rules (resolver)
 - Objectives
   - Add: price_books, price_rules; resolver precedence org_channel > channel > org > global
   - APIs: GET/POST price-books, price-rules; GET /prices/resolve?sku_id&org_id&channel_id
   - Flags: ENABLE_PRICE_BOOKS, ENABLE_PRICE_RESOLUTION
 - Tasks
   - Schema: price entities; effective windows; currency field (multi‑currency ready)
   - Resolver: overlap/tie‑breaker rules with unit tests
   - Security: guard write endpoints; audit decisions
   - Docs: precedence matrix with examples; rounding/tax placeholders
 - Deliverable: tag milestone-05

 Progress (2025‑10‑08)
 - DONE: Schema for price_books + price_rules
 - DONE: APIs under ENABLE_PRICE_BOOKS and ENABLE_PRICE_RESOLUTION
   - POST /price-books, POST /price-rules, GET /prices/resolve
 - NEXT: Overlap windows/unit tests; rounding/tax placeholders
 - Added: audit logs for price book/rule creation and resolution

 M06 — Phase 4: Service Policies + Ticket Responsibility (optional‑first)
 - Objectives
   - Policy engine to compute responsible_org_id; store provenance
   - APIs: GET/PUT policies; integrate optional resolution in ticket create
   - Flags: ENABLE_SERVICE_POLICIES, ENABLE_RESP_ORG_ASSIGNMENT (tenant‑scoped)
 - Tasks
   - Schema: service_policies; nullable columns on tickets
   - Engine: manufacturer→distributor→dealer fallbacks; territory/cert placeholders
   - Observability: assignment latency; breach counters (placeholders)
   - Docs: policy DSL outline; escalation stubs
 - Deliverable: tag milestone-06

Progress (2025‑10‑08)
- DONE: Schema added: service_policies; tickets.responsible_org_id + policy_provenance (nullable)
- DONE: Minimal provenance write on ticket create behind ENABLE_RESP_ORG_ASSIGNMENT
- NEXT: Basic policy rules and actual responsible_org assignment
 - Added: Basic resolver to assign default_org_id from enabled service_policies

 M07 — Phase 5: Engineers + Eligibility (optional)
 - Objectives
   - Add: engineers, engineer_org_memberships, skills, coverage (all optional)
   - API: GET /tickets/{id}/eligible-engineers; manual assign
   - Flags: ENABLE_ENGINEERS
 - Tasks
   - Schema: engineer entities; indexes (geo, skills)
   - Logic: eligibility by org/skills/coverage (no auto‑assign yet)
   - Docs: assignment plugin interface; audit events
 - Deliverable: tag milestone-07

Progress (2025‑10‑08)
- DONE: Schema for engineers, memberships, coverage
- DONE: Read API under ENABLE_ENGINEERS: GET /engineers
- DONE: Eligibility endpoint: GET /engineers/eligible?skills=...&region=...
- NEXT: Ranking/scoring; pagination tuning; tests

 M08 — Phase 6: Agreements + SLA DSL + Events
 - Objectives
   - Add: contracts/agreements; SLA DSL; event/webhook schemas
 - Tasks
   - Schema: agreements; SLA policy; event registry
   - APIs: contracts, SLA policies; webhook delivery with retries
   - Security: HMAC signatures; DLQ strategy
   - Docs: versioning policy; idempotency keys
 - Deliverable: tag milestone-08

 Progress (2025‑10‑08)
 - DONE: Agreements table skeleton (backend-only)
 - DONE: SLA policies schema (sla_policies) with org scoping, effective windows, rules JSON
 - DONE: Events/webhooks outbox schema (service_events, webhook_subscriptions, webhook_deliveries)
 - DONE: Event outbox repository and ticket lifecycle event emission
 - NEXT: Dispatcher worker behind ENABLE_EVENT_DISPATCHER + signing/HMAC + retries/backoff

 Manufacturer & Supplier coexistence (no breakage)
 - Keep legacy tables/APIs intact
 - Profiles approach: organizations + manufacturer_profiles/supplier_profiles
 - Crosswalk tables: manufacturer_organization_map, supplier_organization_map
 - Dual‑write (flagged): writes mirror into orgs+profiles; reads unchanged until ENABLE_READ_FROM_ORG_GRAPH

 Data architecture guardrails (avoid column bloat)
 - Normalize core entities; move variable attributes to profile tables or validated JSONB
 - PII classification; masking/redaction; encryption in transit; backups/DR plan
 - Row‑level security feasibility; least‑privilege roles; migration rollback playbooks

 Software architecture guardrails
 - Domain modules: orgs, catalog, pricing, policy, engineers, agreements
 - Clear service boundaries; DI for resolvers; audit logging for decisions
 - Evented integrations; API versioning; pagination and idempotency standards

 Frontend architecture (deferred)
 - No UI scope now; later: dashboard routing by org_type, catalog manager, price books UI
 - SSR/edge‑safe APIs; loading states; access control by claims

 Tracking and updates
 - Update this file on each milestone completion: summary, links to migrations, APIs, tests, and tags
 - Maintain acceptance criteria per phase; record feature flag states per environment
