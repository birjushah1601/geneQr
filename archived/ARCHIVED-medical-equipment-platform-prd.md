# Medical Equipment & Materials Marketplace  
_Product Requirements Document â€” Version 0.1 (April 2025)_  

---

## 1. Executive Summary  

Hospitals, diagnostic laboratories and individual consumers in India purchase medical equipment and consumables through fragmented Channel Partner networks, resulting in price opacity, long lead-times and manual paperwork.  
The proposed platform is an omni-channel **B2B + B2C ecommerce marketplace** that digitises the end-to-end buying process, giving buyers transparent pricing and real-time stock visibility while allowing manufacturers and authorised Channel Partners to reach customers directly.  

Business Goals (Year 1):  
1. Process â‚¹40 crore Gross Merchandise Value (GMV).  
2. Onboard 150 healthcare institutions and 10 000 retail customers.  
3. Achieve â‰¥95 % on-time delivery within 48 h to tier-1 cities.  

---

## 2. Target Users & Personas  

| Persona | Segment | Pain Points | Goals |  
|---------|---------|-------------|-------|  
| Dr. Asha, Procurement Head | 250-bed private hospital | Phone & email RFQs, 7-day PO cycle, GST audit headaches | Faster sourcing, GST-ready invoices, single supplier portal |  
| Rajesh, Lab Owner | 3-branch pathology chain | Running out of reagents, no MOQ discounts | Subscription deliveries, bulk pricing |  
| Meera, Chronic Patient | Home-care consumer | Inconsistent strip availability, high shipping cost | Reliable deliveries, next-day service |  
| Pankaj, Channel Partner | Regional Channel Partner | Limited digital reach, manual ledger | Wider market access, automated order capture |  

---

## 3. Product Scope & Features  

### 3.1 Marketplace Catalogue  
â€¢ Multi-level category tree (Imaging, ICU, Lab, Consumables, Home-care).  
â€¢ Rich product pages with IFU PDFs, images, technical specs.  
â€¢ IMDR class & HSN code stored for each SKU.  

### 3.2 Pricing & Promotions  
â€¢ List price + tiered B2B pricing tables.  
â€¢ Contract pricing linked to buyer account.  
â€¢ Discount coupons & limited-time promotions for B2C.  

### 3.3 Cart, Checkout & Payments  
â€¢ Mixed cart (multiple sellers) with automatic split orders.  
â€¢ Payment options: Net-30 (B2B on approval), UPI, credit/debit cards, COD (B2C).  
â€¢ GST split by CGST/SGST/IGST on invoice.  

### 3.4 Order Management System (OMS)  
â€¢ Real-time order status: Pending â†’ Processing â†’ Shipped â†’ Delivered.  
â€¢ Partial shipment & back-order handling.  
â€¢ Returns & replacements workflow with RMA numbers.  

### 3.5 Inventory & Warehouse Integration  
â€¢ Stock ledger per SKU and warehouse.  
â€¢ Batch/lot & expiry date tracking for consumables.  
â€¢ FEFO picking logic; barcode/QR scan mobile app for warehouse staff.  

### 3.6 Logistics & Delivery  
â€¢ 3PL integrations (Delhivery, BlueDart) for label, tracking, Proof of Delivery.  
â€¢ Shipment tracking visible to buyers in portal and email notifications.  

### 3.7 Basic Reporting & Analytics  
â€¢ Sales dashboard (GMV, orders, AOV, top SKUs).  
â€¢ Inventory health (days on hand, near expiry).  
â€¢ Tax report export for finance teams.  

---

## 4. Regulatory & Compliance  

1. Indian Medical Device Rules 2017 (class listing, licence numbers on product pages).  
2. CDSCO import licence capture for imported devices.  
3. GST e-invoice and e-way bill integration (NIC API) for orders >â‚¹50 000.  
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
| 0 â€“ Discovery & Design | 0-1 | Requirements, UX wireframes | Finalised PRD, UI mock-ups |  
| 1 â€“ MVP Launch | 2-4 | Catalogue, cart/checkout, single warehouse, payment gateway | Web MVP live with 500 SKUs |  
| 2 â€“ B2B Features | 5-7 | Tiered pricing, Net-30 credit, bulk RFQ module | First hospital group onboarded |  
| 3 â€“ Multi-Warehouse | 8-10 | Stock transfer, batch/expiry, 3PL tracking | 95 % on-time delivery KPI |  
| 4 â€“ Reporting & Scaling | 11-12 | BI dashboards, autoscaling infra | GMV â‚¹40 Cr target |  

---

## 7. Success Metrics  

| Category | KPI | Target |  
|----------|-----|--------|  
| Operational | Order-to-ship time (tier-1) | â‰¤48 h |  
| Financial | Gross margin | â‰¥22 % |  
| Customer | NPS (B2B buyers) | â‰¥50 |  
| Platform | p95 API latency | <300 ms |  

---

## 8. Risks & Mitigation  

| Risk | Probability | Impact | Mitigation |  
|------|-------------|--------|------------|  
| Channel Partner backlash | Med | Med | Present platform as additional channel; optional pricing anonymity |  
| Stock accuracy issues | Med | High | Barcode scanning, daily cycle counts |  
| Regulatory delays (CDSCO) | Low | Med | Pre-validate licences during seller onboarding |  
| Payment defaults (B2B) | Med | Med | Credit checks, trade insurance |  

---

## 9. Glossary  

* **IMDR** â€“ Indian Medical Device Rules 2017  
* **HSN** â€“ Harmonised System of Nomenclature code for GST  
* **FEFO** â€“ First Expired, First Out (inventory)  
* **RMA** â€“ Return Merchandise Authorisation  
* **GMV** â€“ Gross Merchandise Value  

---

_End of Version 0.1 PRD (initial marketplace vision)_  
