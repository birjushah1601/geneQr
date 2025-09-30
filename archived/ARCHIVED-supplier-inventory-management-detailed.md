# Supplier & Inventory Management – Detailed Specification  
_File: ARCHIVED-supplier-inventory-management-detailed.md (Superseded by FINAL tech guide – kept for historical reference)_  

---

## 1. SUPPLIER MANAGEMENT

### 1.1 Supplier On-boarding & Qualification  
1. Digital registration form (GSTIN, CDSCO licence, IEC for importers, ISO 13485 certificate).  
2. Document review workflow:  
   • Statuses – _Draft → Submitted → Under-Review → Approved/Rejected_.  
3. Risk-based site audit for Class C/D device vendors.  
4. Blacklist / watch-list check via Dun & Bradstreet API (credit rating ≥ 6).

### 1.2 Compliance Repository  
• Store licence copies (Form MD-15/MD-42, CE/FDA certificates).  
• Expiry alerts 90, 30, 7 days via email & dashboard.  
• Map each SKU → supplier-licence entry for traceability.

### 1.3 Performance Scorecards  
| KPI | Target | Formula | Review Cadence |  
|-----|--------|---------|----------------|  
| On-time Delivery | ≥ 95 % | Delivered POs / Total POs | Monthly |  
| Defect Rate | ≤ 0.5 % | QC failures / Received units | Monthly |  
| ASN Accuracy | ≥ 98 % | Mismatch qty / Received units | Monthly |  
| Response SLA | ≤ 24 h | Avg RFQ response | Quarterly |  

Auto-flag suppliers < threshold for Corrective Action (CAPA).

### 1.4 Contract & Pricing Management  
• Version-controlled contracts with clause tags (lead-time, warranty, penalties).  
• Tiered pricing tables (MOQ, volume breaks).  
• Renewal reminder 60 days before contract end.

### 1.5 Purchase Order (PO) Flow  
1. Re-order point triggers Draft PO.  
2. PO approval matrix (value-based).  
3. Electronic PO (cXML/EDI or PDF) sent; acknowledgment required within 12 h.  
4. Three-way match (PO ↔ GRN ↔ Invoice) before payment release.

### 1.6 Quality Control (QC)  
• Incoming QC checklist per device class: visual, functional, sterility.  
• Non-conformance module with 8D investigation workflow.  
• Hold & release process for suspect lots.

---

## 2. INVENTORY MANAGEMENT

### 2.1 Multi-Warehouse & Location Hierarchy  
Country → City → Warehouse → Zone → Bin.  
Stock ledger tracks Qty On-Hand, Reserved, In-Transit.

### 2.2 Batch / Lot / Serial Tracking  
• Mandatory serial capture for Class C/D devices.  
• GS1/UDI barcode scan at receipt, move, dispatch.  
• Recall module: query all customers impacted by batch.

### 2.3 Expiry & FEFO  
• FEFO enforced for reagents; FIFO fallback otherwise.  
• Near-expiry (<30 days) dashboard & auto-discount rules.  
• Quarantine & destruction workflow with biomedical-waste form.

### 2.4 Replenishment & Safety Stock  
• Min-Max and EOQ parameters per SKU & warehouse.  
• Draft PO generated when On-Hand + On-Order < Min.  
• Safety stock = (Avg daily demand × Lead-time) + Pandemic buffer 20 %.

### 2.5 Stock Transfers  
• Stock Transfer Order (STO) created; pick → pack → ship → receive.  
• In-transit virtual location.  
• Split STO for ambient vs cold-chain items.

### 2.6 Cycle Counting & Physical Inventory  
• ABC classification (A monthly, B quarterly, C bi-annual).  
• Blind count via mobile app; variance approval >±2 %.  
• Audit trail retained 10 years.

### 2.7 Rental Fleet Asset Tracking  
• States: _Available → Deployed → Service Due → Under Repair → Retired_.  
• Usage meter (hours or test-count) captured on return.  
• Preventive maintenance schedule based on OEM guidelines.

### 2.8 Cold-Chain Monitoring  
• Digital data-loggers / Bluetooth probes log temperature every 10 min.  
• Alert if excursion beyond 2-8 °C or 15-25 °C thresholds.  
• Excursion report attached to batch record.

---

## 3. SYSTEM INTEGRATIONS

| System | Direction | Data | Protocol |  
|--------|-----------|------|----------|  
| ERP (PO, GRN, Invoice) | Bi-directional | POHeader, POLine, GRN, Tax | REST/JSON, CSV |  
| Logistics 3PL | Inbound | Tracking, POD | REST webhook |  
| Barcode Scanners | Outbound | Scan events | HID / Bluetooth |  
| Quality LIMS | Inbound | QC result | CSV import |  

---

## 4. REPORTING & DASHBOARDS

• Supplier scorecard heat-map.  
• Inventory ageing & dead-stock report.  
• Expiry calendar & recall readiness.  
• Warehouse performance (PUT / PICK lines per hour).  
• Rental fleet utilisation & maintenance due chart.

---

## 5. NON-FUNCTIONAL REQUIREMENTS

| Category | Requirement |  
|----------|-------------|  
| Availability | 99.9 % for inventory APIs |  
| Performance | <200 ms query latency up to 1 M SKUs |  
| Auditability | Immutable logs retained 10 yrs |  
| Security | Row-level security by warehouse & tenant |  
| Compliance | IMDR traceability, GST e-invoice, DPDP encryption |  

---

## 6. ROADMAP SNAPSHOT (Original)  

| Phase | Months | Deliverables |  
|-------|--------|-------------|  
| 1 – Foundation | 0-2 | Supplier onboarding portal, core SKU master |  
| 2 – Inventory Core | 3-6 | Batch/serial tracking, FEFO engine, mobile WMS v1 |  
| 3 – Advanced Ops | 7-10 | Cycle counting, STO, supplier scorecards |  
| 4 – Rental Fleet | 11-14 | Asset lifecycle, maintenance scheduler |  
| 5 – Cold Chain | 15-18 | IoT probes integration, excursion dashboard |  

---

*End of archived document – superseded by FINAL-technical-implementation-guide.md which includes AI and advisory enhancements.*  
