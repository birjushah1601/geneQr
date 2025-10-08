 # Multi-Org (Manufacturer/Distributor/Dealer) and Catalog Plan
 
 Status: Draft (brainstorming approved) — Orgs optional; no code committed yet.
 
 ## 1. Goals
 - Support manufacturers, distributors, dealers; dealers can be multi-brand.
 - Enable catalog publishing, channels, and pricing by partner.
 - Keep org references optional (nullable) so current flows keep working.
 - Prepare for service coverage (responsibility, assignment, SLAs) later.
 
 ## 2. Operating Mode (Optional Orgs)
 - Org fields are nullable. When absent, use global defaults/policies.
 - Feature flag: `ENABLE_ORG` (off by default). Enable per tenant gradually.
 
 ## 3. Domain Model (concise)
 - organizations: id, name, type (manufacturer|distributor|dealer|supplier|service_provider|hospital), status, metadata JSONB
 - org_relationships: id, parent_org_id, child_org_id, relation_type (distributor_of|dealer_of|supplier_of|partner_of), active_from/to
 - products: id, manufacturer_id, category, attrs JSONB, lifecycle_status
 - skus: id, product_id, sku_code, attrs JSONB, status
 - offerings: id, sku_id, offer_type (sale|subscription|rental), term, service_level, conditions JSONB
 - channels: id, name, type (direct|distributor|dealer|online)
 - channel_catalog: id, channel_id, org_id NULLABLE (global when null), item_ref (sku_id/offer_id), status, constraints JSONB
 - price_books: id, scope_type (global|org|channel|org_channel), scope_ids JSONB, currency, version, effective_from/to
 - price_rules: id, price_book_id, selector (sku_id|category|all), rule_type (fixed|margin|discount), payload JSONB
 - equipment (extend): sku_id NULLABLE, owner_org_id NULLABLE
 - service_policies: id, scope (global|org|manufacturer), rules JSONB (assignment order, SLAs, escalation)
 - tickets (extend): responsible_org_id NULLABLE, assigned_engineer_org_id NULLABLE, coverage_policy_id NULLABLE
 - engineers (later): engineers, engineer_org_memberships, engineer_skills, engineer_coverage
 
 Indexes (representative):
 - org_relationships: (parent_org_id), (child_org_id), (relation_type,parent_org_id)
 - channel_catalog: (channel_id, org_id NULLS FIRST, status), (item_ref)
 - price_books: (scope_type), price_rules: (price_book_id)
 - tickets: (responsible_org_id), equipment: (owner_org_id)
 
 ## 4. Inputs by Org Type (what we capture)
 - Common: name, type, status, contacts, addresses/regions, metadata JSONB
 - Manufacturer: catalog authority, warranty policy refs, partner tiers, price book templates, escalation policies
 - Distributor: territories, partner dealers, margin rules vs manufacturer, inventory capability, curation scope
 - Dealer: service coverage (geo/time), storefront channel, permitted offerings, preferred distributors, certifications
 - Supplier: parts lines, lead times, return policy, integrations
 - Service provider: engineer pool, skills/certs, capacity/availability, billing model
 - Hospital: sites/departments, approval matrix, contract coverage
 
 All optional; enforce via soft rules before adding hard constraints.
 
 ## 5. Catalog and Channels
 - Manufacturer owns master products/SKUs.
 - Distributors/dealers receive curated subsets via `channel_catalog`.
 - Pricing via `price_books` with rules (fixed/margin/discount) and effective windows.
 - Dealer can list multi-manufacturer SKUs (N:M via relationships).
 
 ## 6. Service Coverage & Tickets (later, optional-first)
 - Policies resolve `responsible_org_id` from manufacturer→distributor→dealer chain, territory, certifications, availability; fallback to global.
 - Track provenance (chosen rule, candidates, escalation path) for audits.
 
 ## 7. APIs (v1 – read-first, org filters optional)
 - Orgs: GET/POST /orgs; GET/POST /orgs/{id}/relationships
 - Channels: GET/POST /channels; GET/POST /channels/{id}/catalog
 - Catalog: GET/POST /products, /skus, /offerings
 - Pricing: GET/POST /price-books; GET/POST /price-rules
 - Policies: GET/PUT /policies/service-coverage; GET/PUT /policies/catalog-publishing
 - Stats (dashboards): GET /stats/tickets, /stats/equipment, /stats/catalog (support org_id/channel_id filters)
 
 ## 8. Dashboards by Org Type (fallback to Global)
 - Manufacturer: ticket SLA, product failure, partner performance
 - Distributor: tickets by dealer/region, backlog, catalog publishing
 - Dealer: local tickets, utilization, allowed listings
 - Hospital: asset uptime, turnaround, coverage/contracts
 - Service provider: jobs, utilization, certifications
 - Global (no org): system health, adoption
 
 Routing: if `org_type` claim present → org dashboard; else Global.
 
 ## 9. AuthZ
 - Keep tenant auth; org-scoped RBAC optional. If claims include org_id/org_type → enforce; else skip.
 
 ## 10. Migration & Backfill (safe)
 - Create new tables with nullable FKs; do not break existing flows.
 - Backfill manufacturers as organizations (optional); map relationships later.
 - Add nullable columns to equipment/tickets only; no destructive changes.
 
 ## 11. Phased Delivery & Acceptance
 - Phase 1: Orgs + relationships (nullable), Channels, Products/SKUs (read)
   - AC: list orgs/relationships; products/skus; channel exists; existing flows unchanged.
 - Phase 2: Channel Catalog + Offerings (read/write), basic Stats
   - AC: publish/unpublish items per channel/org; dashboards read with org filters.
 - Phase 3: Price Books + Rules (basic)
   - AC: resolve price by precedence (org_channel > channel > org > global); versioning/effective.
 - Phase 4: Service Policies + Ticket responsibility
   - AC: responsible_org_id set by policy; escalation recorded; ticket create 200.
 - Phase 5: Engineers + Memberships (optional)
   - AC: eligible-engineers endpoint filters by org/skills/coverage; manual assign works.
 - Phase 6: Analytics & Overrides
   - AC: partner KPIs; pricing overrides per manufacturer/category; audit logs.
 
 Addendum: Future-proofing integration
 - Introduce non-breaking contracts and agreements model alongside Phase 3 (read-only first).
 - Add SLA policy DSL and escalation matrix with Phase 4 (toggle-controlled rollout).
 - Catalog draft→published versioning aligns with Phase 2; enforce immutability post-publish.
 - Pricing enhancements (multi-currency/tax/discount stacking) follow Phase 3 as minor versions.
 - Event/webhook schema introduced early (Phase 2) to avoid later breaking changes.
 
 ## 12. Risks & Mitigations
 - No orgs provided → deterministic global policies; audit trail of decisions.
 - Multi-brand dealer conflicts → SKU dedupe + precedence rules (direct > preferred distributor > others).
 - Pricing complexity → start with fixed/margin rules; add stacking/tiers later.
 - Data quality (serial↔SKU) → allow NULL `sku_id` with best-effort matching.
 
 ## 13. Open Questions
 - Exclusivity: any brand/region exclusivity for distributors/dealers?
 - Pricing precedence: exact order for overlapping price books?
 - Dealer eligibility: mandatory certifications per manufacturer/category?
 - Territory model: ISO regions vs geofences? Source of truth?
 - Global policy defaults: acceptable assignment order and SLAs?
 
 ## 14. Client-specific notes & examples (for review)
 
 Context highlights
 - Multi-brand dealers: Client expects dealers to sell and service across multiple manufacturers; supported via N:M org_relationships and channel_catalog publishing per SKU/Offering.
 - Optional orgs now: Org fields remain nullable to avoid blocking current QR → service request flows already live.
 - QR-driven service: Keep current QR scan → equipment → service request path unchanged; responsible_org_id can remain NULL until policies are enabled.
 
 Example A — Catalog publishing (multi-brand dealer)
 - Orgs: M1 (Manufacturer: ACME), M2 (Manufacturer: BioMed), D1 (Distributor for ACME), DL1 (Dealer)
 - Relationships: D1 distributor_of M1; DL1 dealer_of D1; DL1 dealer_of M2 (direct)
 - Publishing:
   - ACME publishes Product P-A (sku S-A1) to D1; D1 publishes S-A1 to DL1
   - BioMed publishes Product P-B (sku S-B1) directly to DL1
 - Result: DL1 channel_catalog lists S-A1 (via D1) and S-B1 (direct from M2)
 
 Example B — Pricing precedence (simple)
 - Price books (highest to lowest precedence): org_channel (DL1@DealerChannel) > channel (DealerChannel) > org (DL1) > global
 - For S-A1: DL1 has an org_channel fixed price 950 → applied
 - For S-B1: no org_channel; channel price is base 1000 with -5% dealer discount → applied
 
 Example C — Service responsibility (optional-first)
 - Equipment E-123 (Manufacturer=M1) at Hospital H1; ticket created via QR
 - Policies disabled: responsible_org_id NULL; assignment falls back to internal/global rule (e.g., nearest available engineer)
 - Policies enabled later: order = dealer → distributor → manufacturer; DL1 eligible with cert; assign DL1; escalation to D1 then M1 if SLA breached
 
 Example D — Dashboards by org_type
 - Manufacturer (M1): SLA by product family, failures by model, partner performance (D1, DL1)
 - Distributor (D1): tickets by dealers, backlog, catalog coverage
 - Dealer (DL1): open jobs, utilization, published items
 - Global: system health and adoption (used when org is not in context)
 
 Rollout notes
 - Phase 1 can ship with organizations and relationships API (nullable usage) + catalog read; no change to existing ticketing.
 - Enable price books and publishing per partner only when client agrees precedence and discounting rules.
 - Service coverage policies remain opt-in per tenant to control change impact.
 
 ## 15. Future-proofing & Gaps (to document and implement incrementally)
 
 Agreements and pricing
 - Contracts/agreements: party A/B, scope (products/SKUs/categories), terms (warranty, AMC, subscription), effective windows, termination clauses.
 - Pricing precedence: explicit table for overlapping scopes (org_channel > channel > org > global); examples and tests.
 - Multi-currency and tax: currency per price book, FX strategy, tax rules per region; rounding policy.
 - Promotions/discounts: rule types (fixed, percent, tiered), stacking/compatibility policy, eligibility by org/channel/category.
 
 Policies and service
 - SLA policy DSL: response/resolve targets by severity, business hours/holidays, pause conditions, breach actions.
 - Escalation matrix: time-based and condition-based routing across dealer → distributor → manufacturer.
 - Preventive maintenance (PM): schedule templates per SKU/category; linkage to contracts.
 - RMA/returns: statuses, authorization workflow, parts logistics integration points.
 
 Catalog operations
 - Versioning: draft → published lifecycle; immutable published snapshots; deprecations and replacement mapping.
 - Auditability: who/when published, delta summaries, rollback guidelines (new version only).
 - Eligibility: certifications/territories to gate listings per dealer/distributor.
 
 Security and tenancy
 - RBAC matrix: roles (org_admin, ops_manager, engineer, viewer) × resources; org-optional behavior matrix.
 - Data visibility: examples for global vs org-scoped queries; soft-filters when org_id absent.
 - API versioning: v1 compatibility policy; idempotency keys for write endpoints.
 
 Integrations and events
 - Event schema: ticket.created, price.updated, catalog.published, agreement.updated; payload contracts and versioning.
 - Webhooks: retries/backoff, HMAC signatures, dead-letter queue strategy.
 - Plugin points: assignment strategies, pricing engines, import/export adapters.
 
 Migration and BCP
 - Nullable-first migrations; feature flags to gate behavior.
 - Rollback strategy: disable features, preserve data, re-enable after fix.
 - Data import/export playbooks (catalog, price books, orgs/relationships).
 
 Analytics and KPIs
 - KPI definitions: MTTR, FTF, SLA%, backlog burn, partner performance, coverage heatmaps.
 - Data marts or views for dashboards; retention and PII considerations.
 
 Acceptance for future-proofing docs
 - Each subsection includes: scope, non-goal, data model sketch, API impact, rollout/flags, and test checklist.

 ## 16. Gaps Closure Checklist (acceptance + owners)

 Product/Catalog
 - Catalog versioning documented with rollback rules and immutability of published versions
   - Acceptance: examples of draft→published, diffing, and rollback via new version
   - Owner: Product + Backend
 - Eligibility rules (certifications/territories) and enforcement points
   - Acceptance: policy examples; API filter behavior; failure messages
   - Owner: Product + Backend
 - Multi-brand/multi-channel conflict resolution and precedence
   - Acceptance: precedence table and test cases; SKU dedupe behavior
   - Owner: Product
 - Search/indexing and uniqueness constraints
   - Acceptance: defined keys (sku_code, product_id+attrs); indexing plan
   - Owner: Backend

 Pricing/Commercials
 - Precedence matrix finalized with tie-breakers and examples
   - Acceptance: matrix in docs + unit tests for overlapping scopes
   - Owner: Product + Backend
 - Multi-currency/tax/rounding and promotions stacking policy
   - Acceptance: price book currency, FX note, tax examples, stacking table
   - Owner: Product
 - Agreements/contracts schema (warranty/AMC/subscription)
   - Acceptance: ER sketch + lifecycle states + sample payloads
   - Owner: Product

 Service/SLA
 - SLA DSL and escalation matrix
   - Acceptance: severity table, business hours/holidays, breach actions
   - Owner: Product + Backend
 - PM schedules and RMA/returns workflows
   - Acceptance: state diagrams and minimal field sets
   - Owner: Product
 - Assignment plugin interface + audit trail
   - Acceptance: interface signature + audit event schema examples
   - Owner: Backend

 Security/Tenancy
 - RBAC matrix and org-optional behavior table
   - Acceptance: role×resource table with examples; no-org fallback rules
   - Owner: Product + Platform
 - Data visibility rules (global vs org-scoped)
   - Acceptance: query examples; soft-filters when org_id absent
   - Owner: Backend
 - PII/privacy/retention
   - Acceptance: retention windows; redaction policy; export/delete flows
   - Owner: Platform

 APIs/Contracts
 - Error model, pagination/filter standards, idempotency keys
   - Acceptance: style guide page; examples; conformance checklist
   - Owner: Platform
 - API versioning and deprecation process
   - Acceptance: version policy; sunset headers; changelog template
   - Owner: Platform
 - Events/webhooks schemas and retry/backoff/signing
   - Acceptance: schema registry; signature spec; DLQ strategy
   - Owner: Platform

 Operations/SRE
 - Observability SLOs, runbooks, alerting thresholds
   - Acceptance: SLO targets; dashboards; on-call runbooks
   - Owner: SRE
 - Backups/DR, migration rollback strategy, seed/demo data
   - Acceptance: RPO/RTO; rollback steps; demo dataset outline
   - Owner: SRE + Backend
 - Rate limits, caching, performance budgets
   - Acceptance: limits per endpoint; cache keys/TTL; p95 targets
   - Owner: Platform

 Data/Analytics
 - KPI definitions and canonical views
   - Acceptance: metrics dictionary; SQL/view samples
   - Owner: Data
 - Territory model source of truth (ISO vs geofence)
   - Acceptance: chosen model; validation rules; storage format
   - Owner: Product + Data
 - Import/export playbooks (catalog, price books, orgs)
   - Acceptance: CSV/JSON formats; idempotent import semantics
   - Owner: Platform