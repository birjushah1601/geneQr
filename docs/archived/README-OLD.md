# ServQRical Equipment Platform

> **Intelligent Medical Equipment Management & Service Platform**  
> QR-based equipment tracking | Multi-entity engineer management | Tier-based service routing

---

## ðŸŽ¯ Overview

ServQRical Platform is a comprehensive B2B medical equipment lifecycle management system designed for the Indian healthcare ecosystem. It connects **manufacturers**, **Channel Partners**, **Sub-Sub-sub_sub_SUB_DEALERs**, and **hospitals** with intelligent service routing and equipment tracking.

### Key Features

âœ… **QR Code Equipment Registry** - Generate, store, and scan QR codes for equipment tracking  
âœ… **Service Request Management** - Customer-initiated service requests via QR scan  
âœ… **Multi-Entity Organizations** - Manufacturers, Channel Partners, Sub-Sub-sub_sub_SUB_DEALERs, Hospitals with complex relationships  
âœ… **Engineer Management** - 50+ service engineers with skills, certifications, and availability  
âœ… **Tier-Based Routing** - Intelligent engineer assignment with fallback to client in-house BME teams  
âœ… **Real-Time Dashboards** - Role-specific views for all stakeholders  

---

## ðŸ—ï¸ Architecture

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
- Complex relationship modeling (manufacturer â†’ Channel Partner â†’ Sub-sub_SUB_DEALER)
- Geographic territory management

**Current Data:**
- 10 real manufacturers (Siemens, GE, Philips, Medtronic, Abbott, etc.)
- 20 Channel Partners across India
- 15 Sub-Sub-sub_sub_SUB_DEALERs in major cities
- 38 manufacturer-Channel Partner relationships
- 50+ facilities (manufacturing plants, service centers, warehouses)

[ðŸ“– Full Architecture Documentation](docs/architecture/organizations-architecture.md)

---

## ðŸš€ Quick Start

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
psql -U postgres -d medplatform -f database/seed/002_channel_partners.sql
psql -U postgres -d medplatform -f database/seed/003_sub_Sub-Sub-sub_sub_SUB_DEALERs.sql
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

## ðŸ“Š Current Status

### âœ… Phase 1: Database Foundation (COMPLETE)

- [x] Organizations architecture (10 tables)
- [x] Engineer management (4 tables)
- [x] Equipment & service tickets enhancement
- [x] Seed data: 10 manufacturers, 20 Channel Partners, 15 Sub-Sub-sub_sub_SUB_DEALERs
- [x] 38 B2B relationships with business terms

[ðŸ“– Phase 1 Details](docs/database/phase1-complete.md)

### ðŸš§ Phase 2: Backend APIs (IN PROGRESS)

- [x] Equipment Registry API (working)
- [x] QR Generation & Storage (working)
- [ ] Organizations Module API
- [ ] Engineer Management API
- [ ] Service Request Routing API

### ðŸš§ Phase 3: Frontend Development (IN PROGRESS)

- [x] Equipment Registry UI
- [x] QR Code Generation & Display
- [x] Service Request Page
- [ ] Organizations Management UI
- [ ] Engineer Management UI

### â³ Phase 4: Dashboards (PENDING)

- [ ] Manufacturer Dashboard
- [ ] Channel Partner Dashboard
- [ ] Sub-sub_SUB_DEALER Dashboard
- [ ] Hospital Dashboard
- [ ] Service Provider Dashboard
- [ ] Platform Admin Dashboard

[ðŸ“– Implementation Roadmap](docs/architecture/implementation-roadmap.md)

---

## ðŸŽ“ Key Concepts

### Multi-Entity Engineer Management

Engineers can belong to different organizations:
- **Manufacturer Engineers**: OEM-certified, Tier-1 routing
- **Sub-sub_SUB_DEALER Engineers**: Multi-brand trained, Tier-2 routing
- **Channel Partner Engineers**: Regional coverage, Tier-3 routing
- **Service Provider Engineers**: Independent, Tier-4 routing
- **Hospital BME Engineers**: In-house, Tier-5 fallback

[ðŸ“– Engineer Management Design](docs/architecture/engineer-management.md)

### Tier-Based Service Routing

```
Service Request
    â†“
1. OEM Engineer (if covered)
    â†“
2. Authorized Sub-sub_SUB_DEALER Engineer
    â†“
3. Channel Partner Service Team
    â†“
4. Third-Party Service Provider
    â†“
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

## ðŸ“ Project Structure

```
ServQR/
â”œâ”€â”€ cmd/
â”‚   â””â”€â”€ platform/           # Backend main entry point
â”œâ”€â”€ internal/
â”‚   â””â”€â”€ core/              # Business logic modules
â”‚       â”œâ”€â”€ equipment-registry/
â”‚       â”œâ”€â”€ organizations/
â”‚       â””â”€â”€ service-ticket/
â”œâ”€â”€ admin-ui/              # Next.js frontend
â”‚   â”œâ”€â”€ app/              # Next.js 13 app router
â”‚   â”œâ”€â”€ components/       # React components
â”‚   â””â”€â”€ lib/             # Utilities & API clients
â”œâ”€â”€ database/
â”‚   â”œâ”€â”€ migrations/      # SQL schema migrations
â”‚   â””â”€â”€ seed/           # Seed data SQL files
â”œâ”€â”€ docs/
â”‚   â”œâ”€â”€ architecture/   # Architecture docs
â”‚   â””â”€â”€ database/      # Database docs
â””â”€â”€ dev/
    â””â”€â”€ compose/       # Docker compose files
```

---

## ðŸ§ª Testing

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

## ðŸ“– Documentation

- [Organizations Architecture](docs/architecture/organizations-architecture.md) - Complete multi-entity design
- [Engineer Management](docs/architecture/engineer-management.md) - Tier-based routing system
- [Implementation Roadmap](docs/architecture/implementation-roadmap.md) - 4-week execution plan
- [Phase 1 Complete](docs/database/phase1-complete.md) - Database foundation summary

---

## ðŸ¤ Contributing

This is a private project. For access or questions, contact the development team.

---

## ðŸ“ License

Proprietary - ServQRical Platform  
Â© 2024 All Rights Reserved

---

## ðŸ”— Quick Links

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081/api/v1
- **Database**: PostgreSQL on port 5433

---

**Built with â¤ï¸ for the Indian Healthcare Ecosystem**
