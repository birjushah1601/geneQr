# Medical Equipment & Materials Marketplace  
_Product Requirements Document — Version 0.1 (April 2025)_  

---

## 1. Executive Summary  

Hospitals, diagnostic laboratories and individual consumers in India purchase medical equipment and consumables through fragmented distributor networks, resulting in price opacity, long lead-times and manual paperwork.  
The proposed platform is an omni-channel **B2B + B2C ecommerce marketplace** that digitises the end-to-end buying process, giving buyers transparent pricing and real-time stock visibility while allowing manufacturers and authorised distributors to reach customers directly.  

Business Goals (Year 1):  
1. Process ₹40 crore Gross Merchandise Value (GMV).  
2. Onboard 150 healthcare institutions and 10 000 retail customers.  
3. Achieve ≥95 % on-time delivery within 48 h to tier-1 cities.  

---

## 2. Target Users & Personas  

| Persona | Segment | Pain Points | Goals |  
|---------|---------|-------------|-------|  
| Dr. Asha, Procurement Head | 250-bed private hospital | Phone & email RFQs, 7-day PO cycle, GST audit headaches | Faster sourcing, GST-ready invoices, single supplier portal |  
| Rajesh, Lab Owner | 3-branch pathology chain | Running out of reagents, no MOQ discounts | Subscription deliveries, bulk pricing |  
| Meera, Chronic Patient | Home-care consumer | Inconsistent strip availability, high shipping cost | Reliable deliveries, next-day service |  
| Pankaj, Distributor | Regional distributor | Limited digital reach, manual ledger | Wider market access, automated order capture |  

---

## 3. Product Scope & Features  

### 3.1 Marketplace Catalogue  
• Multi-level category tree (Imaging, ICU, Lab, Consumables, Home-care).  
• Rich product pages with IFU PDFs, images, technical specs.  
• IMDR class & HSN code stored for each SKU.  

### 3.2 Pricing & Promotions  
• List price + tiered B2B pricing tables.  
• Contract pricing linked to buyer account.  
• Discount coupons & limited-time promotions for B2C.  

### 3.3 Cart, Checkout & Payments  
• Mixed cart (multiple sellers) with automatic split orders.  
• Payment options: Net-30 (B2B on approval), UPI, credit/debit cards, COD (B2C).  
• GST split by CGST/SGST/IGST on invoice.  

### 3.4 Order Management System (OMS)  
• Real-time order status: Pending → Processing → Shipped → Delivered.  
• Partial shipment & back-order handling.  
• Returns & replacements workflow with RMA numbers.  

### 3.5 Inventory & Warehouse Integration  
• Stock ledger per SKU and warehouse.  
• Batch/lot & expiry date tracking for consumables.  
• FEFO picking logic; barcode/QR scan mobile app for warehouse staff.  

### 3.6 Logistics & Delivery  
• 3PL integrations (Delhivery, BlueDart) for label, tracking, Proof of Delivery.  
• Shipment tracking visible to buyers in portal and email notifications.  

### 3.7 Basic Reporting & Analytics  
• Sales dashboard (GMV, orders, AOV, top SKUs).  
• Inventory health (days on hand, near expiry).  
• Tax report export for finance teams.  

---

## 4. Regulatory & Compliance  

1. Indian Medical Device Rules 2017 (class listing, licence numbers on product pages).  
2. CDSCO import licence capture for imported devices.  
3. GST e-invoice and e-way bill integration (NIC API) for orders >₹50 000.  
4. Data privacy: DPDP Act 2023 compliance, AES-256 at rest, TLS 1.3 in transit.  

---

## 5. Technology Stack (Proposed)  

| Layer | Tech Choice | Notes |  
|-------|-------------|-------|  
| Front-end | React (Next.js) | SSR for SEO; Tailwind UI |  
| Mobile App | React Native | Consumer app for repeat orders |  
| Back-end | Node.js (NestJS) | REST APIs, modular services |  
| Database | PostgreSQL | ACID, JSONB for attributes |  
| Cache | Redis | Sessions, cart, rate-limit |  
| Search | Elasticsearch | Full-text, synonyms |  
| Queue | RabbitMQ | Email, SMS async tasks |  
| Hosting | AWS (ap-south-1) | RDS, S3, CloudFront CDN |  

---

## 6. Implementation Roadmap  

| Phase | Months | Scope | Key Deliverables |  
|-------|--------|-------|------------------|  
| 0 – Discovery & Design | 0-1 | Requirements, UX wireframes | Finalised PRD, UI mock-ups |  
| 1 – MVP Launch | 2-4 | Catalogue, cart/checkout, single warehouse, payment gateway | Web MVP live with 500 SKUs |  
| 2 – B2B Features | 5-7 | Tiered pricing, Net-30 credit, bulk RFQ module | First hospital group onboarded |  
| 3 – Multi-Warehouse | 8-10 | Stock transfer, batch/expiry, 3PL tracking | 95 % on-time delivery KPI |  
| 4 – Reporting & Scaling | 11-12 | BI dashboards, autoscaling infra | GMV ₹40 Cr target |  

---

## 7. Success Metrics  

| Category | KPI | Target |  
|----------|-----|--------|  
| Operational | Order-to-ship time (tier-1) | ≤48 h |  
| Financial | Gross margin | ≥22 % |  
| Customer | NPS (B2B buyers) | ≥50 |  
| Platform | p95 API latency | <300 ms |  

---

## 8. Risks & Mitigation  

| Risk | Probability | Impact | Mitigation |  
|------|-------------|--------|------------|  
| Distributor backlash | Med | Med | Present platform as additional channel; optional pricing anonymity |  
| Stock accuracy issues | Med | High | Barcode scanning, daily cycle counts |  
| Regulatory delays (CDSCO) | Low | Med | Pre-validate licences during seller onboarding |  
| Payment defaults (B2B) | Med | Med | Credit checks, trade insurance |  

---

## 9. Glossary  

* **IMDR** – Indian Medical Device Rules 2017  
* **HSN** – Harmonised System of Nomenclature code for GST  
* **FEFO** – First Expired, First Out (inventory)  
* **RMA** – Return Merchandise Authorisation  
* **GMV** – Gross Merchandise Value  

---

_End of Version 0.1 PRD (initial marketplace vision)_  
