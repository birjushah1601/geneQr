# Field Service Management (FSM) Module - Implementation Plan

## Overview
This document outlines the implementation of the Field Service Management module for medical equipment after-sales service. The system enables QR code-based equipment registration, WhatsApp-based service requests, and complete ticket management.

## Business Model
- **Equipment Registration**: Track all installed medical equipment with unique QR codes
- **AMC Contracts**: Manage Annual Maintenance Contracts with SLA tracking
- **Service Requests**: Receive service requests via WhatsApp (photos, videos, QR codes)
- **Ticket Management**: Full lifecycle from creation to resolution
- **Engineer Dispatch**: Assign and track field engineers

---

## Module Architecture

```
FSM Module
├── Equipment Registry Service (Foundation)
├── QR Code Generation Service
├── Service Ticket Management Service
├── WhatsApp Integration Service
├── Engineer Management Service
└── AMC Contract Service
```

---

## Implementation Status

### ✅ COMPLETED (Phase 1)

#### 1. Equipment Registry - Domain Model
**File**: `internal/service-domain/equipment-registry/domain/equipment.go`

**Features**:
- Equipment entity with complete lifecycle
- Status management: operational, down, under_maintenance, decommissioned
- Service history tracking
- Warranty and AMC linkage
- CSV import structures

**Business Methods**:
- `MarkAsDown()` - Equipment failure
- `MarkUnderMaintenance()` - Service in progress
- `RecordService()` - Complete service visit
- `ScheduleNextService()` - Plan maintenance
- `IsUnderWarranty()` - Check warranty status
- `HasAMC()` - Check AMC coverage

#### 2. Repository Interface
**File**: `internal/service-domain/equipment-registry/domain/repository.go`

**Methods**:
- CRUD operations
- Query by QR code
- Query by serial number
- List with filtering (status, manufacturer, customer)
- Bulk create for CSV import

#### 3. Database Schema
**File**: `dev/postgres/migrations/008-create-equipment-registry.sql`

**Tables**:
- `equipment_registry` - Main equipment table
  - 30+ columns for complete equipment tracking
  - JSONB for specifications, photos, documents
  - Unique constraints on QR code and serial number

**Indexes**:
- 11 standard indexes for performance
- 2 GIN indexes for JSONB queries

**Database Functions**:
- `generate_qr_code()` - Auto-generate unique QR codes
- `get_equipment_statistics()` - Equipment analytics
- `get_equipment_needing_service()` - Service reminders
- `get_expired_warranty_equipment()` - Warranty expiry alerts

---

## 📋 REMAINING IMPLEMENTATION (This Week)

### Phase 2: Core Services (2-3 days)

#### 1. QR Code Generation Service
**File**: `internal/service-domain/equipment-registry/qrcode/generator.go`

**Features**:
- Generate QR code with unique URL
- URL format: `https://service.example.com/eq/{qr_code}`
- PDF label generation for printing
- Batch generation for CSV imports

**Libraries Needed**:
```go
github.com/skip2/go-qrcode  // QR generation
github.com/jung-kurt/gofpdf  // PDF generation
```

#### 2. PostgreSQL Repository
**File**: `internal/service-domain/equipment-registry/infra/repository.go`

**Implementation**:
- Full CRUD with JSONB handling
- Bulk insert for CSV import
- Query optimization
- Transaction support

#### 3. CSV Import Service
**File**: `internal/service-domain/equipment-registry/app/csv_import.go`

**Features**:
- Parse CSV file
- Validate rows
- Bulk insert with error handling
- Return import results (success/failure counts)

**CSV Format**:
```
serial_number,equipment_name,manufacturer_name,model_number,category,customer_name,customer_id,installation_location,installation_date,purchase_date,purchase_price,warranty_months,notes
SN001,X-Ray Machine,GE Healthcare,DX-9000,Radiology,City Hospital,CUST001,Radiology Dept,2024-01-15,2024-01-10,300000,24,New installation
```

#### 4. Application Service
**File**: `internal/service-domain/equipment-registry/app/service.go`

**Use Cases**:
- `RegisterEquipment` - Single equipment registration
- `BulkImportFromCSV` - CSV bulk import
- `GenerateQRCode` - Generate and print QR
- `GetEquipmentByQR` - Lookup by QR scan
- `GetEquipmentBySerial` - Lookup by serial
- `ListEquipment` - With filtering
- `UpdateEquipment` - Update details
- `RecordService` - Log service completion
- `GenerateQRLabels` - PDF labels for printing

#### 5. REST API
**File**: `internal/service-domain/equipment-registry/api/handler.go`

**Endpoints**:
```
POST   /equipment                    - Register equipment
POST   /equipment/import            - CSV import
GET    /equipment                    - List equipment
GET    /equipment/{id}              - Get by ID
GET    /equipment/qr/{qr_code}      - Get by QR code
GET    /equipment/serial/{serial}   - Get by serial
PATCH  /equipment/{id}              - Update
DELETE /equipment/{id}              - Delete

POST   /equipment/{id}/qr           - Generate QR code
GET    /equipment/{id}/qr/pdf       - Download QR PDF label
POST   /equipment/bulk/qr           - Batch generate QR codes

POST   /equipment/{id}/service      - Record service
GET    /equipment/statistics        - Get statistics
GET    /equipment/service-due       - Equipment needing service
```

---

### Phase 3: WhatsApp Integration (2-3 days)

#### 1. WhatsApp Webhook Service
**File**: `internal/service-domain/whatsapp/webhook/handler.go`

**Features**:
- Receive webhook from WhatsApp Business API
- Parse message types (text, image, video, document)
- Extract QR codes from images
- Extract serial numbers from text (OCR)
- Download media attachments
- Send responses back to customer

**Integration Options**:
1. **Twilio WhatsApp API** (Recommended for MVP)
   - Easy setup
   - Good documentation
   - Pay-per-message
   
2. **Meta WhatsApp Cloud API** (Free tier available)
   - 1000 free conversations/month
   - Official Meta API

**Webhook Flow**:
```
Customer sends message
    ↓
Webhook receives POST request
    ↓
Parse message (text/image/video)
    ↓
If image: Decode QR code or OCR serial
    ↓
Lookup equipment in registry
    ↓
Create service ticket
    ↓
Send acknowledgment to customer
```

#### 2. QR Code Recognition
**File**: `internal/service-domain/whatsapp/qrcode/decoder.go`

**Features**:
- Download image from WhatsApp
- Decode QR code
- Extract equipment ID/URL
- Fallback to OCR for serial numbers

**Libraries**:
```go
github.com/liyue201/goqr       // QR decoding
github.com/otiai10/gosseract   // OCR
```

#### 3. Media Processing
**File**: `internal/service-domain/whatsapp/media/processor.go`

**Features**:
- Download photos/videos from WhatsApp
- Compress images
- Store in S3 or local storage
- Generate thumbnails
- Link to service tickets

---

### Phase 4: Service Ticket Management (2-3 days)

#### 1. Ticket Domain Model
**File**: `internal/service-domain/service-ticket/domain/ticket.go`

**Features**:
- Ticket lifecycle: new → assigned → in_progress → resolved → closed
- Priority levels: critical, high, medium, low
- SLA tracking
- Engineer assignment
- Media attachments
- Customer communication history

#### 2. Database Schema
**File**: `dev/postgres/migrations/009-create-service-tickets.sql`

**Tables**:
- `service_tickets` - Main ticket table
- `ticket_comments` - Comments and updates
- `ticket_status_history` - Audit trail

#### 3. Application Service
**File**: `internal/service-domain/service-ticket/app/service.go`

**Use Cases**:
- `CreateTicket` - From WhatsApp or manual
- `AssignEngineer` - Manual assignment
- `UpdateStatus` - Lifecycle management
- `AddComment` - Customer/engineer notes
- `AttachMedia` - Photos/videos
- `ResolveTicket` - Mark as resolved
- `CloseTicket` - Final closure
- `ReopenTicket` - If issue persists

---

### Phase 5: Engineer Management (1-2 days)

#### 1. Engineer Domain Model
**File**: `internal/service-domain/engineer/domain/engineer.go`

**Features**:
- Engineer profiles
- Skills and specializations
- Availability status
- Geographic coverage
- Performance metrics

#### 2. Manual Assignment (MVP)
- Simple dropdown to select engineer
- No auto-routing yet
- Basic availability check

---

## Technical Stack

### Required Go Packages
```bash
# QR Code generation and recognition
go get github.com/skip2/go-qrcode
go get github.com/liyue201/goqr

# PDF generation
go get github.com/jung-kurt/gofpdf

# CSV processing
# Use standard library encoding/csv

# Image processing
go get github.com/disintegration/imaging

# OCR (optional for MVP)
go get github.com/otiai10/gosseract

# WhatsApp (choose one)
go get github.com/twilio/twilio-go  # If using Twilio
```

### External Services
1. **WhatsApp Business API**
   - Twilio account (recommended)
   - OR Meta WhatsApp Cloud API

2. **Storage** (for media)
   - Local filesystem (MVP)
   - OR MinIO (S3-compatible)
   - OR AWS S3 (production)

---

## Testing Plan

### Week 1 Testing
1. Register 10 equipment via API
2. Generate QR codes and PDF labels
3. Import 100 equipment from CSV
4. Query by QR code and serial number

### Week 2 Testing
1. Send QR code image to WhatsApp
2. Verify equipment recognition
3. Check auto-ticket creation
4. Test customer acknowledgment

### Week 3 Testing
1. Create tickets manually
2. Assign engineers
3. Update ticket status
4. Close tickets with notes

---

## MVP Scope (2 weeks)

### Must Have ✅
1. Equipment Registration (manual + CSV)
2. QR Code Generation
3. WhatsApp Integration (receive messages)
4. QR Code Recognition from images
5. Auto-Ticket Creation
6. Basic Ticket Management
7. Engineer Manual Assignment

### Nice to Have (Future)
1. AMC Contract Management
2. Preventive Maintenance Scheduling
3. Parts Inventory
4. Mobile App for Engineers
5. AI-powered issue analysis
6. Auto-engineer routing
7. Customer portal

---

## Database Migrations Checklist

- [x] 008-create-equipment-registry.sql
- [ ] 009-create-service-tickets.sql
- [ ] 010-create-engineers.sql
- [ ] 011-create-amc-contracts.sql (future)

---

## API Endpoints Summary

### Equipment Registry (18 endpoints)
- 8 core CRUD endpoints
- 5 query endpoints
- 3 QR code endpoints
- 2 statistics endpoints

### Service Tickets (15 endpoints)
- 5 core CRUD endpoints
- 4 lifecycle endpoints
- 3 comment endpoints
- 3 query endpoints

### Engineers (8 endpoints)
- 5 core CRUD endpoints
- 3 query endpoints

**Total New Endpoints**: ~41

---

## Next Steps

1. **Today**: Complete Equipment Registry Service + QR Generation
2. **Tomorrow**: Build WhatsApp Integration
3. **Day 3**: Service Ticket Management
4. **Day 4-5**: Engineer Management + Testing
5. **Week 2**: Polish, bug fixes, documentation

---

## Questions to Resolve

1. **WhatsApp Number**: Do you have WhatsApp Business API access?
2. **QR URL Format**: What should be the base URL? (e.g., `https://service.yourcompany.com/eq/{qr_code}`)
3. **Storage**: Local filesystem OK for MVP, or prefer S3/MinIO?
4. **Notification Preferences**: Email + WhatsApp, or WhatsApp only?
5. **Engineer Count**: How many engineers for MVP testing?

---

## Success Criteria

✅ 100 equipment registered with QR codes  
✅ QR PDF labels generated and printable  
✅ WhatsApp message → Ticket created  
✅ Customer receives acknowledgment  
✅ Engineer can view and accept tickets  
✅ End-to-end service workflow completed  

---

*This is a living document. Update as implementation progresses.*
