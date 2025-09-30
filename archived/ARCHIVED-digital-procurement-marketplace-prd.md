# Digital Procurement Marketplace – Product Requirements Document  
_File: ARCHIVED-digital-procurement-marketplace-prd.md (stand-alone version, May 2025 – superseded by unified docs)_  

---

## 1. EXECUTIVE SUMMARY  

Indian hospitals and laboratory groups purchase medical equipment through two-to-three layers of distributors, adding 9-15 % mark-ups and extending the quotation–PO cycle to 7-10 days.  
The **Digital Procurement Marketplace** is a B2B platform that connects buyers directly with manufacturers/OEMs, digitises RFQ-to-PO workflows, and integrates with hospital SAP systems.  
This document captures the marketplace requirements **before** subsequent AI and service-platform expansions.

### 1.1 Business Objectives (Year 1)  
1. Cut order finalisation cycle from 7-10 days to **≤ 24 h**.  
2. Achieve **₹40 Cr GMV** with 2–4 % commission.  
3. Deliver **≥8 % landed-cost savings** to buyers versus traditional distribution.

---

## 2. STAKEHOLDERS & PERSONAS  

| Persona | Organisation | Goals | Pain Points |
|---------|--------------|-------|-------------|
| Procurement Head (Dr Asha) | 250-bed hospital | Faster PO processing, price transparency | Email RFQ chaos, manual SAP entry |
| Operations Manager (Rajesh) | Pathology chain | Bulk reagent pricing, repeat ordering | MOQ hurdles, delayed quotes |
| OEM Sales Manager (Kiran) | Manufacturer | Direct demand visibility, quick deal closure | Distributor dependency, credit risk |
| Finance Officer (Neha) | Hospital | GST-ready invoices, audit trail | Manual data re-entry, errors |

---

## 3. VALUE PROPOSITION  

### 3.1 For Buyers  
• Transparent OEM pricing and real-time stock.  
• One-click PO creation with GST split and SAP upload.  
• Reduced working-capital lock via faster quote turnaround.

### 3.2 For Sellers  
• Direct access to hospital demand without distributor fees.  
• Digital quote builder and contract management.  
• Lower DSO through milestone-linked payment terms.

---

## 4. PRODUCT FEATURES  

### 4.1 Catalogue & Search  
• Multi-level medical category tree, filter by modality, IMDR class, price.  
• Rich SKU page with technical specs, IFU PDF, downloadable compliance docs.  
• Real-time inventory indicator (Green > 10, Amber ≤ 10, Red = 0).

### 4.2 RFQ & Negotiation Workflow  
State machine:

```
Draft → Published → Quotes Received → Counter Offer → Accepted/Rejected → Closed
```

• Buyers can publish RFQ to specific OEMs or open marketplace.  
• Sellers submit structured quotes (price, lead-time, warranty).  
• Buyers can issue **one counter**; sellers accept or decline.  
• Time-box: default 24 h expiration; auto-expire closes RFQ.

### 4.3 Purchase Order & Contract  
• “Accept” action converts winning quote to Purchase Order (PO).  
• PO PDF auto-generated, digitally signed, and emailed + stored.  
• SAP/MM adapter pushes `ORDERS05` IDoc into hospital SAP system.  
• Contract repository stores T&Cs, warranty clauses, delivery schedules.

### 4.4 Pricing & Commission Model  
• Platform commission 2–4 % charged to seller; configurable per SKU family.  
• Volume rebate engine tracks cumulative buyer spend → auto-issues credit notes monthly.  
• Early-payment discount configuration (e.g., 2 % 10/Net 30).

### 4.5 Logistics & Tracking  
• Seller books shipment via integrated 3PL widget; tracking ID stored.  
• Delivery status events (`In-Transit`, `Out for Delivery`, `Delivered`) update buyer dashboard.  
• Proof-of-Delivery (POD) PDF stored against PO.

### 4.6 Compliance Console (Basic)  
• Mandatory licence fields: CDSCO import licence #, GSTIN.  
• Expiry alerts inside seller portal (30 d, 7 d).  
• Attachments validated for PDF/A format.

### 4.7 Reporting  
| Report | Metrics | Audience |
|--------|---------|----------|
| Buyer Savings | OEM list vs negotiated price | Hospital CFO |  
| Quote Cycle | Avg time per RFQ stage | Ops |  
| Seller Performance | Quote response time, fulfilment SLA | OEM mgr |

---

## 5. NON-FUNCTIONAL REQUIREMENTS  

| Area | Requirement |
|------|-------------|
| Availability | 99.9 % for RFQ & PO APIs |
| Performance | p95 < 300 ms API latency at 5 k concurrent users |
| Security | OAuth 2.1, AES-256 at rest, TLS 1.3 |
| Compliance | DPDP 2023 data privacy, IMDR licence display |
| Audit | Immutable log retention 10 years |

---

## 6. TECHNICAL ARCHITECTURE (NON-AI)  

| Layer | Component | Tech Choice |
|-------|-----------|-------------|
| Front-end | Web portal | React + Next.js |
| Back-end | RFQ-svc, Quote-svc, Contract-svc | Node.js (Nest) |
| DB | Catalogue & transaction data | PostgreSQL |
| Search | Full-text keyword | Elasticsearch |
| Integration | SAP connector | Apache Camel + IDoc |
| Messaging | Order events | RabbitMQ |
| File Storage | Documents | AWS S3 |

Event examples:

```
RFQCreated → QuoteSubmitted → QuoteAccepted → POGenerated → ShipmentDispatched → Delivered
```

---

## 7. IMPLEMENTATION ROADMAP  

| Phase | Months | Deliverables | Exit Criteria |
|-------|--------|-------------|---------------|
| 0 Discovery | 0-1 | Final PRD, UX wireframes | Stakeholder sign-off |
| 1 MVP | 2-4 | Catalogue, RFQ, quote submit, manual PO PDF | 5 live RFQs completed |
| 2 SAP & 3PL | 5-7 | SAP IDoc push, 3PL tracking integration | First auto-synced PO; shipment tracked |
| 3 Contract & Rebate | 8-10 | Contract repo, rebate engine | Volume rebate issued to pilot buyers |
| 4 Scaling | 11-12 | Multi-tenant onboarding portal | 50 hospitals, 20 OEMs onboarded |

---

## 8. KPIs & SUCCESS METRICS  

| KPI | Target |
|-----|--------|
| RFQ → PO cycle time | ≤ 24 h median |
| Average price saving | ≥ 8 % vs baseline |
| Commission revenue | ₹2 Cr first year |
| API error rate | < 0.1 % |
| Net Promoter Score | ≥ 50 |

---

## 9. RISKS & MITIGATION  

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Distributor push-back | Medium | Position platform as optional channel; NDA pricing |
| OEM data reluctance | High | Manual SKU upload support, NDA, early-success case study |
| Buyer SAP integration delays | Medium | Provide sandbox IDoc generator, professional services |
| Regulatory changes (GST slabs) | Low | Configurable tax engine |

---

## 10. CHANGE LOG  

| Version | Date | Author | Notes |
|---------|------|--------|-------|
| 0.1 | 05 May 2025 | Product Lead | Initial draft |
| 0.2 | 18 May 2025 | Tech Architect | Added SAP connector scope |
| 0.3 | 10 Jun 2025 | Ops | Incorporated logistics tracking |
| **Archived** | 21 Sept 2025 | PMO | Superseded by unified AI-enabled PRD |

---

_End of archived standalone Digital Procurement Marketplace PRD_  
