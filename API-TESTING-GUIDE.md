# API Testing Guide - ABY-MED Medical Platform

## Overview
This guide provides step-by-step instructions to test all services running on the ABY-MED platform.

**Base URL:** `http://localhost:8081`

**Common Headers:**
- `X-Tenant-ID: city-hospital` (defaults to this if not provided)
- `X-User-ID: test-user` (optional)
- `Content-Type: application/json`

---

## Table of Contents
1. [Health & Metrics](#1-health--metrics)
2. [Catalog Service](#2-catalog-service)
3. [RFQ Service](#3-rfq-service)
4. [Supplier Service](#4-supplier-service)
5. [Quote Service](#5-quote-service)
6. [Comparison Service](#6-comparison-service)
7. [Contract Service](#7-contract-service)
8. [Equipment Registry](#8-equipment-registry)
9. [Service Ticket](#9-service-ticket)

---

## 1. Health & Metrics

### Check Platform Health
```bash
curl http://localhost:8081/health
```

**Expected Response:**
```json
{"status":"ok"}
```

### Get Prometheus Metrics
```bash
curl http://localhost:8081/metrics
```

---

## 2. Catalog Service

### List All Catalog Items
```bash
curl -X GET http://localhost:8081/api/v1/catalog \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json"
```

### Get Catalog Item by ID
```bash
curl -X GET http://localhost:8081/api/v1/catalog/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Create Catalog Item
```bash
curl -X POST http://localhost:8081/api/v1/catalog \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MRI Scanner - Siemens Magnetom",
    "sku": "MRI-001-SIEMENS",
    "category": "Diagnostic Imaging",
    "manufacturer": "Siemens Healthineers",
    "description": "1.5T MRI Scanner with advanced imaging capabilities",
    "specifications": {
      "field_strength": "1.5 Tesla",
      "bore_diameter": "70 cm",
      "max_patient_weight": "250 kg"
    },
    "base_price": 1500000.00,
    "currency": "INR",
    "compliance_certifications": ["CDSCO", "ISO 13485", "CE Mark"]
  }'
```

### Update Catalog Item
```bash
curl -X PUT http://localhost:8081/api/v1/catalog/{id} \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MRI Scanner - Siemens Magnetom Skyra",
    "base_price": 1750000.00
  }'
```

### Delete Catalog Item
```bash
curl -X DELETE http://localhost:8081/api/v1/catalog/{id} \
  -H "X-Tenant-ID: city-hospital"
```

---

## 3. RFQ Service

### List All RFQs
```bash
curl -X GET http://localhost:8081/api/v1/rfq \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user"
```

### Get RFQ by ID
```bash
curl -X GET http://localhost:8081/api/v1/rfq/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Create New RFQ
```bash
curl -X POST http://localhost:8081/api/v1/rfq \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Request for MRI Scanner",
    "description": "Need 1.5T MRI Scanner for radiology department",
    "deadline": "2025-11-01T00:00:00Z",
    "delivery_location": {
      "facility_name": "City Hospital Radiology",
      "address": "123 Medical Lane, Mumbai, Maharashtra",
      "pincode": "400001",
      "contact_person": "Dr. Sharma",
      "phone": "+91-9876543210"
    },
    "items": [
      {
        "catalog_id": "catalog-item-id",
        "quantity": 1,
        "specifications": "1.5T with advanced cardiac imaging"
      }
    ]
  }'
```

### Add Item to RFQ
```bash
curl -X POST http://localhost:8081/api/v1/rfq/{id}/items \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "catalog_id": "catalog-item-id",
    "quantity": 2,
    "specifications": "Additional requirements"
  }'
```

### Publish RFQ (Send to Suppliers)
```bash
curl -X POST http://localhost:8081/api/v1/rfq/{id}/publish \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user"
```

### Close RFQ
```bash
curl -X POST http://localhost:8081/api/v1/rfq/{id}/close \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user"
```

### Cancel RFQ
```bash
curl -X POST http://localhost:8081/api/v1/rfq/{id}/cancel \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user"
```

---

## 4. Supplier Service

### List All Suppliers
```bash
curl -X GET http://localhost:8081/api/v1/suppliers \
  -H "X-Tenant-ID: city-hospital"
```

### Get Supplier by ID
```bash
curl -X GET http://localhost:8081/api/v1/suppliers/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Register New Supplier
```bash
curl -X POST http://localhost:8081/api/v1/suppliers \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "MedTech Supplies Pvt Ltd",
    "gstin": "27AABCU9603R1ZM",
    "pan": "AABCU9603R",
    "contact_person": "Rajesh Kumar",
    "email": "rajesh@medtechsupplies.com",
    "phone": "+91-9876543210",
    "address": {
      "street": "456 Business Park",
      "city": "Mumbai",
      "state": "Maharashtra",
      "pincode": "400051",
      "country": "India"
    },
    "categories": ["Diagnostic Imaging", "Laboratory Equipment"],
    "certifications": ["ISO 9001", "CDSCO Registered"]
  }'
```

### Update Supplier
```bash
curl -X PUT http://localhost:8081/api/v1/suppliers/{id} \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "contact@medtechsupplies.com",
    "phone": "+91-9876543211"
  }'
```

---

## 5. Quote Service

### List All Quotes
```bash
curl -X GET http://localhost:8081/api/v1/quotes \
  -H "X-Tenant-ID: city-hospital"
```

### Get Quote by ID
```bash
curl -X GET http://localhost:8081/api/v1/quotes/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Create Quote (Supplier Response to RFQ)
```bash
curl -X POST http://localhost:8081/api/v1/quotes \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: supplier-user" \
  -H "Content-Type: application/json" \
  -d '{
    "rfq_id": "rfq-id-here",
    "supplier_id": "supplier-id-here",
    "validity_days": 30,
    "payment_terms": "30 days from delivery",
    "delivery_timeline": "60-90 days",
    "warranty": "24 months comprehensive",
    "items": [
      {
        "rfq_item_id": "rfq-item-id",
        "unit_price": 1650000.00,
        "discount_percent": 5.0,
        "tax_percent": 18.0,
        "delivery_charges": 50000.00,
        "installation_charges": 100000.00
      }
    ],
    "notes": "Installation and training included"
  }'
```

### Update Quote
```bash
curl -X PUT http://localhost:8081/api/v1/quotes/{id} \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "validity_days": 45,
    "notes": "Extended validity period"
  }'
```

---

## 6. Comparison Service

### Create Quote Comparison
```bash
curl -X POST http://localhost:8081/api/v1/comparisons \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "rfq_id": "rfq-id-here",
    "quote_ids": ["quote-1-id", "quote-2-id", "quote-3-id"],
    "comparison_criteria": [
      "total_cost",
      "delivery_timeline",
      "warranty",
      "supplier_rating"
    ]
  }'
```

### Get Comparison by ID
```bash
curl -X GET http://localhost:8081/api/v1/comparisons/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### List Comparisons for RFQ
```bash
curl -X GET "http://localhost:8081/api/v1/comparisons?rfq_id={rfq_id}" \
  -H "X-Tenant-ID: city-hospital"
```

---

## 7. Contract Service

### Create Contract from Quote
```bash
curl -X POST http://localhost:8081/api/v1/contracts \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "quote_id": "quote-id-here",
    "buyer_signatory": "Dr. Ramesh Verma, Chief Administrator",
    "supplier_signatory": "Rajesh Kumar, Director",
    "special_terms": [
      "Installation within 90 days",
      "Training for 5 staff members",
      "Annual maintenance contract included for first year"
    ],
    "payment_schedule": [
      {
        "milestone": "PO Confirmation",
        "percentage": 30,
        "due_date": "2025-10-15T00:00:00Z"
      },
      {
        "milestone": "Pre-Delivery Inspection",
        "percentage": 40,
        "due_date": "2025-12-01T00:00:00Z"
      },
      {
        "milestone": "Installation Complete",
        "percentage": 30,
        "due_date": "2026-01-15T00:00:00Z"
      }
    ]
  }'
```

### Get Contract by ID
```bash
curl -X GET http://localhost:8081/api/v1/contracts/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### List All Contracts
```bash
curl -X GET http://localhost:8081/api/v1/contracts \
  -H "X-Tenant-ID: city-hospital"
```

### Sign Contract
```bash
curl -X POST http://localhost:8081/api/v1/contracts/{id}/sign \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "signature": "digital-signature-hash",
    "signed_at": "2025-10-01T12:00:00Z"
  }'
```

---

## 8. Equipment Registry

### Register Equipment Asset
```bash
curl -X POST http://localhost:8081/api/v1/equipment \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Siemens MRI Magnetom Skyra",
    "serial_number": "MRI-SKYRA-2025-001",
    "model": "Magnetom Skyra 1.5T",
    "manufacturer": "Siemens Healthineers",
    "purchase_date": "2025-01-15",
    "installation_date": "2025-03-01",
    "location": {
      "building": "Main Hospital Block A",
      "floor": "Ground Floor",
      "room": "Radiology Suite 1"
    },
    "warranty_expiry": "2027-03-01",
    "maintenance_schedule": "quarterly"
  }'
```

### List All Equipment
```bash
curl -X GET http://localhost:8081/api/v1/equipment \
  -H "X-Tenant-ID: city-hospital"
```

### Get Equipment by ID
```bash
curl -X GET http://localhost:8081/api/v1/equipment/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Generate QR Code for Equipment
```bash
curl -X POST http://localhost:8081/api/v1/equipment/{id}/qr-code \
  -H "X-Tenant-ID: city-hospital"
```

### Update Equipment Status
```bash
curl -X PUT http://localhost:8081/api/v1/equipment/{id}/status \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "operational",
    "notes": "Regular maintenance completed"
  }'
```

---

## 9. Service Ticket

### Create Service Ticket
```bash
curl -X POST http://localhost:8081/api/v1/tickets \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": "equipment-id-here",
    "title": "MRI Scanner - Strange Noise During Scan",
    "description": "Scanner making unusual clicking noise during cardiac imaging sequences",
    "priority": "high",
    "reported_by": "Dr. Anjali Mehta",
    "contact_phone": "+91-9876543210",
    "contact_email": "anjali.mehta@cityhospital.com"
  }'
```

### List All Tickets
```bash
curl -X GET http://localhost:8081/api/v1/tickets \
  -H "X-Tenant-ID: city-hospital"
```

### Get Ticket by ID
```bash
curl -X GET http://localhost:8081/api/v1/tickets/{id} \
  -H "X-Tenant-ID: city-hospital"
```

### Update Ticket Status
```bash
curl -X PUT http://localhost:8081/api/v1/tickets/{id}/status \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "assigned_to": "technician-id",
    "notes": "Technician dispatched, ETA 2 hours"
  }'
```

### Add Ticket Comment
```bash
curl -X POST http://localhost:8081/api/v1/tickets/{id}/comments \
  -H "X-Tenant-ID: city-hospital" \
  -H "X-User-ID: test-user" \
  -H "Content-Type: application/json" \
  -d '{
    "comment": "Technician identified issue with cooling system fan",
    "is_internal": false
  }'
```

### Close Ticket
```bash
curl -X POST http://localhost:8081/api/v1/tickets/{id}/close \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "resolution": "Replaced cooling fan, tested successfully",
    "parts_used": ["Cooling Fan Assembly - Part #CF-001"],
    "labor_hours": 2.5
  }'
```

### WhatsApp Webhook (For receiving messages)
```bash
# This endpoint is called by WhatsApp Business API
curl -X POST http://localhost:8081/api/v1/whatsapp/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "entry": [{
      "changes": [{
        "value": {
          "messages": [{
            "from": "919876543210",
            "text": {
              "body": "My MRI scanner is not working"
            }
          }]
        }
      }]
    }]
  }'
```

---

## Testing Workflow Examples

### Complete Procurement Workflow

1. **Create Catalog Items** (if not exists)
2. **Create RFQ** with items
3. **Publish RFQ** to suppliers
4. **Suppliers Create Quotes** (multiple suppliers)
5. **Create Comparison** of all quotes
6. **Select Best Quote** and **Create Contract**
7. **Sign Contract** by both parties
8. **Register Equipment** after delivery
9. **Generate QR Code** for asset tracking

### Service Management Workflow

1. **Register Equipment Asset**
2. **Generate QR Code** for equipment
3. **Create Service Ticket** when issue occurs
4. **Update Ticket Status** as technician works
5. **Add Comments** for progress tracking
6. **Close Ticket** when resolved

---

## Monitoring & Observability

### Grafana Dashboards
- **URL:** http://localhost:3000
- **Credentials:** admin / admin
- View real-time metrics, request rates, latencies

### Prometheus Metrics
- **URL:** http://localhost:9090
- Query custom metrics and alerts

### MailHog (Email Testing)
- **URL:** http://localhost:8025
- View all outgoing emails during testing

---

## Troubleshooting

### Common Issues

1. **401 Unauthorized**: Add `X-Tenant-ID` header
2. **404 Not Found**: Check endpoint path matches `/api/v1/...`
3. **500 Internal Server Error**: Check platform logs
4. **Connection Refused**: Ensure platform is running on port 8081

### Check Logs
```bash
# Platform logs are showing in the terminal where you started the platform
# Or check Docker logs for infrastructure services
cd dev/compose
docker compose -p med-platform logs -f postgres
docker compose -p med-platform logs -f kafka
```

---

## Next Steps

1. Import the Postman collection (see `ABY-MED-Postman-Collection.json`)
2. Set up environment variables in Postman
3. Run the collection to test all endpoints
4. Create custom test scenarios based on your requirements
