# ABY Medical Equipment Platform

> **Intelligent Medical Equipment Management & Service Platform**  
> QR-based equipment tracking | Multi-entity engineer management | Tier-based service routing

---

## 🎯 Overview

ABY Medical Platform is a comprehensive B2B medical equipment lifecycle management system designed for the Indian healthcare ecosystem. It connects **manufacturers**, **distributors**, **dealers**, and **hospitals** with intelligent service routing and equipment tracking.

### Key Features

✅ **QR Code Equipment Registry** - Generate, store, and scan QR codes for equipment tracking  
✅ **Service Request Management** - Customer-initiated service requests via QR scan  
✅ **Multi-Entity Organizations** - Manufacturers, Distributors, Dealers, Hospitals with complex relationships  
✅ **Engineer Management** - 50+ service engineers with skills, certifications, and availability  
✅ **Tier-Based Routing** - Intelligent engineer assignment with fallback to client in-house BME teams  
✅ **Real-Time Dashboards** - Role-specific views for all stakeholders  

---

## 🏗️ Architecture

### Technology Stack

**Backend:**
- Go 1.21+
- PostgreSQL 12+ with PostGIS
- RESTful APIs

**Frontend:**
- Next.js 13+ (App Router)
- React 18
- TypeScript
- Tailwind CSS
- Shadcn/ui components

### Database Schema

**Organizations Architecture:**
- 10 core tables for multi-entity management
- 4 engineer management tables
- Complex relationship modeling (manufacturer → distributor → dealer)
- Geographic territory management

**Current Data:**
- 10 real manufacturers (Siemens, GE, Philips, Medtronic, Abbott, etc.)
- 20 distributors across India
- 15 dealers in major cities
- 38 manufacturer-distributor relationships
- 50+ facilities (manufacturing plants, service centers, warehouses)

[📖 Full Architecture Documentation](docs/architecture/organizations-architecture.md)

---

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL 12+
- Docker (optional, for containerized database)

### 1. Database Setup

```bash
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Run migrations
psql -U postgres -d medplatform -f database/migrations/002_organizations_simple.sql

# Load seed data
psql -U postgres -d medplatform -f database/seed/001_manufacturers.sql
psql -U postgres -d medplatform -f database/seed/002_distributors.sql
psql -U postgres -d medplatform -f database/seed/003_dealers.sql
```

### 2. Backend Setup

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5433
export DB_NAME=medplatform
export DB_USER=postgres
export DB_PASSWORD=postgres

# Build and run
cd cmd/platform
go run main.go
```

Backend will start on **http://localhost:8081**

### 3. Frontend Setup

```bash
cd admin-ui
npm install
npm run dev
```

Frontend will start on **http://localhost:3000**

---

## 📊 Current Status

### ✅ Phase 1: Database Foundation (COMPLETE)

- [x] Organizations architecture (10 tables)
- [x] Engineer management (4 tables)
- [x] Equipment & service tickets enhancement
- [x] Seed data: 10 manufacturers, 20 distributors, 15 dealers
- [x] 38 B2B relationships with business terms

[📖 Phase 1 Details](docs/database/phase1-complete.md)

### 🚧 Phase 2: Backend APIs (IN PROGRESS)

- [x] Equipment Registry API (working)
- [x] QR Generation & Storage (working)
- [ ] Organizations Module API
- [ ] Engineer Management API
- [ ] Service Request Routing API

### 🚧 Phase 3: Frontend Development (IN PROGRESS)

- [x] Equipment Registry UI
- [x] QR Code Generation & Display
- [x] Service Request Page
- [ ] Organizations Management UI
- [ ] Engineer Management UI

### ⏳ Phase 4: Dashboards (PENDING)

- [ ] Manufacturer Dashboard
- [ ] Distributor Dashboard
- [ ] Dealer Dashboard
- [ ] Hospital Dashboard
- [ ] Service Provider Dashboard
- [ ] Platform Admin Dashboard

[📖 Implementation Roadmap](docs/architecture/implementation-roadmap.md)

---

## 🎓 Key Concepts

### Multi-Entity Engineer Management

Engineers can belong to different organizations:
- **Manufacturer Engineers**: OEM-certified, Tier-1 routing
- **Dealer Engineers**: Multi-brand trained, Tier-2 routing
- **Distributor Engineers**: Regional coverage, Tier-3 routing
- **Service Provider Engineers**: Independent, Tier-4 routing
- **Hospital BME Engineers**: In-house, Tier-5 fallback

[📖 Engineer Management Design](docs/architecture/engineer-management.md)

### Tier-Based Service Routing

```
Service Request
    ↓
1. OEM Engineer (if covered)
    ↓
2. Authorized Dealer Engineer
    ↓
3. Distributor Service Team
    ↓
4. Third-Party Service Provider
    ↓
5. Hospital In-House BME (Fallback)
```

### QR Code System

Each equipment has a unique QR code that encodes:
```json
{
  "url": "http://localhost:3000/service-request?qr=QR-eq-001",
  "id": "EQ-123456",
  "serial": "SN-2024-001",
  "qr": "QR-eq-001"
}
```

Scanning triggers service request with auto-filled equipment details.

---

## 📁 Project Structure

```
aby-med/
├── cmd/
│   └── platform/           # Backend main entry point
├── internal/
│   └── core/              # Business logic modules
│       ├── equipment-registry/
│       ├── organizations/
│       └── service-ticket/
├── admin-ui/              # Next.js frontend
│   ├── app/              # Next.js 13 app router
│   ├── components/       # React components
│   └── lib/             # Utilities & API clients
├── database/
│   ├── migrations/      # SQL schema migrations
│   └── seed/           # Seed data SQL files
├── docs/
│   ├── architecture/   # Architecture docs
│   └── database/      # Database docs
└── dev/
    └── compose/       # Docker compose files
```

---

## 🧪 Testing

### Backend API Testing

```bash
# Test equipment API
curl http://localhost:8081/api/v1/equipment

# Test QR generation
curl -X POST http://localhost:8081/api/v1/equipment/EQ-001/qr

# Test QR retrieval
curl http://localhost:8081/api/v1/equipment/QR-eq-001/qr-image
```

### Database Queries

```sql
-- Check organizations
SELECT org_type, COUNT(*) FROM organizations GROUP BY org_type;

-- Check relationships
SELECT COUNT(*) FROM org_relationships;

-- Check equipment with QR codes
SELECT id, equipment_name, qr_code_id, 
       CASE WHEN qr_code_image IS NOT NULL THEN 'Yes' ELSE 'No' END as has_qr
FROM equipment;
```

---

## 📖 Documentation

- [Organizations Architecture](docs/architecture/organizations-architecture.md) - Complete multi-entity design
- [Engineer Management](docs/architecture/engineer-management.md) - Tier-based routing system
- [Implementation Roadmap](docs/architecture/implementation-roadmap.md) - 4-week execution plan
- [Phase 1 Complete](docs/database/phase1-complete.md) - Database foundation summary

---

## 🤝 Contributing

This is a private project. For access or questions, contact the development team.

---

## 📝 License

Proprietary - ABY Medical Platform  
© 2024 All Rights Reserved

---

## 🔗 Quick Links

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081/api/v1
- **Database**: PostgreSQL on port 5433

---

**Built with ❤️ for the Indian Healthcare Ecosystem**
