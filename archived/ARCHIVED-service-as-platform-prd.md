# Service-as-a-Platform (SaaP) – Product Requirements Document  
_File: ARCHIVED-service-as-platform-prd.md – Stand-alone service platform specification (superseded by unified docs, kept for historical reference)_  

---

## 1. EXECUTIVE SUMMARY  

Indian hospitals and laboratory chains suffer >48 h average downtime on critical medical devices because fewer than 20 % of installed equipment is covered by Annual Maintenance Contracts (AMCs).  
Manufacturers have limited field-service reach; hospitals rely on ad-hoc local technicians with no KPI transparency.  

The Service-as-a-Platform (SaaP) delivers a **nation-wide, SLA-driven service network** that manages AMCs, break-fix tickets, preventive maintenance and parts logistics on behalf of both hospitals (buyers) and manufacturers (sellers).

### 1.1 Business Objectives (Yr-1)  
1. Cover ≥ 500 high-value devices (MRI, CT, ICU ventilators) under AMC.  
2. Achieve ≥ 95 % SLA compliance (Platinum plan).  
3. Generate ₹ 3 Cr recurring revenue from AMC subscriptions.

### 1.2 Value Proposition  
• **Hospitals / Labs** – Guaranteed ≥98 % uptime, single portal for multi-OEM service, predictable AMC cost.  
• **Manufacturers / OEMs** – Extended service reach without CAPEX, consolidated field data, upsell parts revenue.  
• **Platform Operator** – Recurring subscription income + parts margin, strong customer lock-in.

---

## 2. BUSINESS MODEL & REVENUE STREAMS  

| Stream | Unit | Price | Target GM% | Notes |  
|--------|------|-------|------------|-------|  
| AMC Subscription | ₹ / device / month | ₹ 2–5 k | 45 % | Tiered by device class & SLA plan |  
| Break-Fix Ticket | ₹ / call + parts | ₹ 5 k avg | 35 % | For non-AMC fleet; convert to AMC |  
| Parts Logistics | Mark-up | 10 % | 25 % | Bulk purchasing from OEMs |  
| Installation & De-installation | One-time | ₹ 15 k | 40 % | New device installs & relocations |

---

## 3. USER PERSONAS  

| Persona | Role | Goals |  
|---------|------|-------|  
| Dr. Mehta | Hospital Biomed Head | Minimise downtime, track AMC renewals |  
| Anita Rao | Lab Operations Manager | Quick reagent-analyzer repair, single point of contact |  
| Vikram Singh | OEM Service Director | Increase AMC attach rate, reduce travel cost |  
| Ramesh Yadav | Field Engineer | Clear work orders, parts availability, mobile checklists |

---

## 4. FUNCTIONAL REQUIREMENTS  

### 4.1 Asset Registry  
• Unique serial/UDI record with make, model, install & warranty dates.  
• Contract linkage: warranty → AMC plan.  
• Attach service history & calibration certificates (PDF).

### 4.2 AMC Contract Management  
• SLA plans: Platinum, Gold, Silver (response / resolution targets).  
• Pricing engine supports tenure & device class.  
• Auto-renewal reminder 60 / 30 / 7 days.

### 4.3 Ticketing System  
• Omnichannel intake: web, mobile, phone (operator entry).  
• Ticket fields: priority, error code, device serial, location.  
• Status workflow: _Open → Assigned → In Progress → Awaiting Parts → Completed → Closed_.  
• SLA timer service; breach logs penalty ledger.

### 4.4 Dispatch & Workforce Management  
• Skill-matrix (modality, vendor certification).  
• Geo-routing to nearest qualified engineer.  
• Mobile app push notification; accept/decline within 5 min.  
• Route sheet with customer contact & checklist.

### 4.5 Parts Inventory & Logistics  
• Regional hubs (metros) + satellite depots.  
• Parts BOM per device model; alternates supported.  
• Reservation on ticket assignment; courier booking API.

### 4.6 Preventive Maintenance Scheduler  
• Calendar / usage-based triggers (run hours, test count).  
• Bulk PM plan generation each month; auto-create work orders.  
• Compliance report for NABH/NABL audits.

### 4.7 Mobile Technician App  
• Offline mode, QR scan to load asset.  
• Step-by-step checklist; capture readings, photos, customer signature.  
• Parts consumption & warranty claim flags.  
• Generate service report PDF on completion.

### 4.8 Reporting & Dashboards  
| Dashboard | Metrics | Audience |  
|-----------|---------|----------|  
| SLA Compliance | MTTR, response time, breach count | Hospital, OEM |  
| Fleet Utilisation | Uptime%, scheduled vs unscheduled calls | Hospital |  
| Parts Consumption | Top parts, stock-out risk, cost | Platform Ops |  
| Engineer Productivity | Calls/day, first-time-fix rate | Service Manager |

---

## 5. INTEGRATIONS  

| System | Direction | Objects | Protocol |  
|--------|-----------|---------|----------|  
| OEM ERP | Inbound | Parts price, error code docs | REST / CSV |  
| Hospital CMMS / HIS | Bi-dir | Asset master sync, ticket status | HL7 FHIR / REST |  
| Logistics 3PL | Inbound | Tracking events | Webhook JSON |  
| Accounting ERP | Outbound | Invoice, credit note | REST / CSV |

---

## 6. NON-FUNCTIONAL REQUIREMENTS  

| Area | Requirement |  
|------|-------------|  
| Availability | 99.9 % ticketing & dispatch APIs |  
| Performance | Assign engineer ≤30 s, ticket create ≤1 s |  
| Security | OAuth2 + MFA, AES-256 at rest, TLS 1.3 |  
| Audit | Immutable logs 10 years (21 CFR 11 style) |  
| Compliance | ISO 13485 service SOPs, DPDP data privacy |  

---

## 7. IMPLEMENTATION ROADMAP  

| Phase | Months | Deliverables | Exit Criteria |  
|-------|--------|-------------|---------------|  
| 0 Foundation | 0-2 | Asset DB, ticket MVP, 2 service hubs stocked | 10 pilot devices onboarded |  
| 1 Reactive Service | 2-6 | Dispatch engine, mobile app v1, SLA timers | MTTR ≤18 h, 100 AMC devices |  
| 2 Preventive Maint. | 6-9 | Scheduler, PM calendar, compliance reports | PM compliance ≥90 % |  
| 3 Nationwide Hubs | 9-14 | 8 hubs live, parts API with 3 OEMs | SLA compliance ≥95 % |  
| 4 Scale & Optimise | 14-24 | 8 000 devices, role-based dashboards | EBITDA break-even, uptime ≥98 % |

---

## 8. RISK & MITIGATION  

| Risk | Likelihood | Impact | Mitigation |  
|------|-----------|--------|------------|  
| OEM reluctance to share parts pricing | High | High | NDA + revenue share, manual CSV fallback |  
| Technician churn | Med | Med | Certification program, retention bonus |  
| Parts supply delays | Med | High | Safety stock, alternate suppliers |  
| Regulatory audit failure | Low | High | ISO 13485 alignment, internal audits |  

---

_End of stand-alone SaaP PRD (pre-AI version)_  
