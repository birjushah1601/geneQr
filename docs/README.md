# ServQR Documentation

Welcome to the ServQR platform documentation! This guide will help you navigate all available documentation.

---

## Quick Start

**New to ServQR?** Start here:
1. **[Getting Started](01-GETTING-STARTED.md)** - Installation and setup (15 minutes)
2. **[Quick Reference](QUICK-REFERENCE.md)** - Common commands and workflows

**For Developers:**
- **[Architecture](02-ARCHITECTURE.md)** - System design and structure
- **[API Reference](04-API-REFERENCE.md)** - REST API documentation

**For Product/Business:**
- **[Features](03-FEATURES.md)** - Complete feature catalog
- **[Executive Summary](EXECUTIVE-SUMMARY.md)** - Business overview

---

## Documentation Structure

### Core Documentation (Start Here)


docs/
├── 01-GETTING-STARTED.md       # Installation, setup, first run
├── 02-ARCHITECTURE.md          # System design, tech stack
├── 03-FEATURES.md              # Feature catalog and status
├── 04-API-REFERENCE.md         # API endpoints documentation
├── 05-DEPLOYMENT.md            # Production deployment
├── 05-TESTING.md               # Testing guide
├── QUICK-REFERENCE.md          # Command cheat sheet
├── EXECUTIVE-SUMMARY.md        # Business overview
└── DOCUMENTATION-INDEX.md      # Complete docs index


### How-To Guides


docs/guides/
├── engineer-management.md      # Manage engineers (all org types)
├── csv-imports.md              # Bulk import organizations/equipment
├── qr-code-setup.md            # QR code generation and usage
└── (more guides...)


### Feature Specifications


docs/specs/
├── PARTNER-ASSOCIATION-SPECIFICATION.md       # Partner network (NEW!)
├── QR-CODE-MIGRATION-PLAN.md                  # QR system architecture
└── DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md  # Multi-entity design


### Configuration & Setup


docs/
├── EXTERNAL-SERVICES-SETUP.md  # Configure SendGrid, OpenAI, etc.
├── DEPLOYMENT-GUIDE.md         # Detailed deployment procedures
├── PRODUCTION-DEPLOYMENT-CHECKLIST.md  # Pre-launch checklist
├── SECURITY-IMPLEMENTATION-COMPLETE.md  # Security features
└── NOTIFICATIONS-SYSTEM.md     # Email notification system


### Historical & Archived


docs/archived/
├── implementation-status/      # Completed feature status docs
├── README-OLD.md               # Previous README
└── (other archived docs...)


---

## Documentation by Audience

### For New Developers
**Goal:** Get up and running quickly

1. [01-GETTING-STARTED.md](01-GETTING-STARTED.md)
2. [QUICK-REFERENCE.md](QUICK-REFERENCE.md)
3. [02-ARCHITECTURE.md](02-ARCHITECTURE.md)
4. [guides/engineer-management.md](guides/engineer-management.md)

### For Frontend Developers
**Goal:** Build UI components and integrate APIs

1. [02-ARCHITECTURE.md](02-ARCHITECTURE.md) - Frontend stack
2. [04-API-REFERENCE.md](04-API-REFERENCE.md) - API endpoints
3. [03-FEATURES.md](03-FEATURES.md) - Feature requirements

### For Backend Developers
**Goal:** Build APIs and business logic

1. [02-ARCHITECTURE.md](02-ARCHITECTURE.md) - Backend structure
2. [04-API-REFERENCE.md](04-API-REFERENCE.md) - API specs
3. [specs/](specs/) - Feature specifications

### For DevOps Engineers
**Goal:** Deploy and maintain the platform

1. [05-DEPLOYMENT.md](05-DEPLOYMENT.md)
2. [DEPLOYMENT-GUIDE.md](DEPLOYMENT-GUIDE.md)
3. [PRODUCTION-DEPLOYMENT-CHECKLIST.md](PRODUCTION-DEPLOYMENT-CHECKLIST.md)
4. [EXTERNAL-SERVICES-SETUP.md](EXTERNAL-SERVICES-SETUP.md)

### For Product Managers
**Goal:** Understand features and capabilities

1. [EXECUTIVE-SUMMARY.md](EXECUTIVE-SUMMARY.md)
2. [03-FEATURES.md](03-FEATURES.md)
3. [06-PERSONAS.md](06-PERSONAS.md)
4. [specs/](specs/) - Feature specifications

### For Architects & Tech Leads
**Goal:** Understand system design decisions

1. [02-ARCHITECTURE.md](02-ARCHITECTURE.md)
2. [specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md](specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md)
3. [specs/QR-CODE-MIGRATION-PLAN.md](specs/QR-CODE-MIGRATION-PLAN.md)
4. [specs/PARTNER-ASSOCIATION-SPECIFICATION.md](specs/PARTNER-ASSOCIATION-SPECIFICATION.md)

---

## Key Features Overview

### Multi-Tenancy
Complete data isolation per organization with 8 organization types supported.

### Service Tickets
Full lifecycle management with WhatsApp integration, AI diagnostics, and engineer assignment.

### Equipment Registry
QR code-based tracking with bulk import support.

### Engineer Management
All organization types (manufacturers, channel partners, sub-dealers) can manage engineers.

### Partner Network (NEW!)
Associate channel partners and sub-dealers with manufacturers for expanded service coverage.

**Spec:** [specs/PARTNER-ASSOCIATION-SPECIFICATION.md](specs/PARTNER-ASSOCIATION-SPECIFICATION.md)

---

## Technology Stack

- **Backend:** Go 1.21+, PostgreSQL 15+, Gin Framework
- **Frontend:** Next.js 14+, TypeScript, React 18, Tailwind CSS
- **AI/ML:** OpenAI GPT-4, Anthropic Claude, Whisper STT
- **Integrations:** SendGrid (Email), Twilio (WhatsApp)

---

## Quick Commands


# Start development
go run cmd/platform/main.go           # Backend
cd admin-ui && npm run dev             # Frontend

# Run tests
go test ./...                          # Backend tests
cd admin-ui && npm test                # Frontend tests

# Build production
go build -o platform cmd/platform/main.go
cd admin-ui && npm run build

# Database
./scripts/run-migrations.sh           # Run migrations
psql -h localhost -p 5430 -U postgres -d med_platform


---

## Getting Help

- **Technical Issues:** Check the relevant guide in [guides/](guides/)
- **API Questions:** See [04-API-REFERENCE.md](04-API-REFERENCE.md)
- **Deployment Issues:** See [05-DEPLOYMENT.md](05-DEPLOYMENT.md)

---

## Contributing to Docs

When updating documentation:
1. Keep it concise and actionable
2. Use clear headings and code examples
3. Update this README if adding major sections
4. Archive outdated docs to [archived/](archived/)

---

## Recent Updates

**February 2026:**
- ✅ Documentation cleanup and reorganization
- ✅ New README with clear structure
- ✅ Partner Association specification added
- ✅ Guides moved to docs/guides/
- ✅ Specs moved to docs/specs/
- ✅ Archived old status documents

**December 2025:**
- Engineer management for all org types
- AI diagnostics with multi-model support
- WhatsApp integration complete

---

**Last Updated:** February 4, 2026  
**Platform Version:** 2.0  
**Status:** Active Development

---

Happy Coding! 🚀
