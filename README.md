# ServQR - Medical Equipment Service Management Platform

> **Intelligent Service Management for Medical Equipment**  
> Multi-tenant | AI-Powered | QR-Based Tracking | Engineer Management

---

## Overview

**ServQR** is a comprehensive B2B medical equipment lifecycle management platform for the healthcare ecosystem. It connects manufacturers, channel partners, sub-dealers, hospitals, and service engineers through intelligent service routing and equipment tracking.

### Platform Highlights

- **Multi-Tenant Architecture** - Complete data isolation per organization  
- **QR Code Equipment Registry** - Track equipment with scannable QR codes  
- **Service Ticket Management** - Complete lifecycle from creation to resolution  
- **Engineer Management** - All organization types can manage engineers  
- **AI-Powered Diagnostics** - GPT-4/Claude integration for troubleshooting  
- **WhatsApp Integration** - Create tickets via text or voice messages  
- **Parts Catalog** - Track parts, pricing, and inventory  
- **Partner Network** - Associate channel partners and sub-dealers (NEW!)  

---

## Technology Stack

**Backend:** Go 1.21+ | PostgreSQL 15+ | Gin Framework  
**Frontend:** Next.js 14+ | TypeScript | React 18 | Tailwind CSS  
**Integration:** OpenAI GPT-4 | Anthropic Claude | SendGrid | Twilio

---

## Documentation

### Core Documentation
- **[Getting Started](docs/01-GETTING-STARTED.md)** - Setup and installation
- **[Architecture](docs/02-ARCHITECTURE.md)** - System design
- **[Features](docs/03-FEATURES.md)** - Complete feature catalog
- **[API Reference](docs/04-API-REFERENCE.md)** - API documentation
- **[Deployment](docs/05-DEPLOYMENT.md)** - Production deployment
- **[Testing](docs/05-TESTING.md)** - Testing guide

### How-To Guides
- **[Engineer Management](docs/guides/engineer-management.md)** - Manage engineers (all org types)
- **[CSV Imports](docs/guides/csv-imports.md)** - Bulk import data
- **[QR Code Setup](docs/guides/qr-code-setup.md)** - QR code generation
- **[External Services](docs/EXTERNAL-SERVICES-SETUP.md)** - Configure integrations

### Feature Specifications
- **[Partner Association](docs/specs/PARTNER-ASSOCIATION-SPECIFICATION.md)** - Partner network (hybrid approach)
- **[QR Code Migration](docs/specs/QR-CODE-MIGRATION-PLAN.md)** - QR system architecture
- **[Organizations Architecture](docs/specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md)** - Multi-entity design

---

## Quick Start

### Prerequisites
- Go 1.21+, Node.js 18+, PostgreSQL 15+, Git

### Installation

bash
# 1. Clone and setup
git clone <repo-url>
cd aby-med
cp .env.example .env

# 2. Start database
docker-compose up -d postgres

# 3. Run migrations
./scripts/run-migrations.sh

# 4. Start backend
go run cmd/platform/main.go

# 5. Start frontend (new terminal)
cd admin-ui
npm install
npm run dev


**Access:** http://localhost:3000 (Frontend) | http://localhost:8081 (API)

---

## Project Structure


ServQR/
├── cmd/platform/          # Backend entry point
├── internal/              # Business logic
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business services
│   └── models/           # Data models
├── admin-ui/             # Next.js frontend
│   ├── src/app/         # App router pages
│   └── src/components/  # React components
├── database/
│   ├── migrations/      # SQL migrations
│   └── seed/           # Seed data
├── docs/                # Documentation
│   ├── guides/         # How-to guides
│   ├── specs/          # Feature specs
│   └── archived/       # Historical docs
└── storage/            # File uploads


---

## Key Features

### Completed Features

**Multi-Tenancy**
- 8 organization types (manufacturer, channel_partner, sub_dealer, hospital, etc.)
- Complete data isolation
- Role-based access control

**Service Tickets**
- Full lifecycle management
- WhatsApp integration (text/voice)
- AI-powered diagnostics (GPT-4, Claude)
- Engineer assignment
- File attachments

**Equipment Registry**
- QR code generation and tracking
- Bulk CSV import
- Maintenance history

**Engineer Management**
- ALL organization types can manage engineers (manufacturers, channel partners, sub-dealers)
- Assignment tracking
- Skills and certifications

**Parts Management**
- Parts catalog
- Ticket-parts association
- Pricing and inventory

**Notifications**
- Email notifications (SendGrid)
- Daily summary reports

### In Development

**Partner Network Management** (Specification Complete)
- Associate channel partners/sub-dealers with manufacturers
- Hybrid approach: organization-level (default) + equipment-level (override)
- Categorized engineer display
- Smart filtering logic

**Spec:** [docs/specs/PARTNER-ASSOCIATION-SPECIFICATION.md](docs/specs/PARTNER-ASSOCIATION-SPECIFICATION.md)

---

## Testing


# Backend API
curl http://localhost:8081/health
curl http://localhost:8081/api/v1/equipment

# Frontend
cd admin-ui
npm run test
npm run build

# Database
psql -h localhost -p 5430 -U postgres -d med_platform


---

## Deployment

See **[Deployment Guide](docs/05-DEPLOYMENT.md)** for production deployment instructions.

---

## License

Proprietary - ServQR Platform  
© 2024-2026 All Rights Reserved

---

**Built for Healthcare** 🏥
