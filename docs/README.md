# ServQR Medical Platform - Documentation

## ðŸ“š Documentation Structure

This documentation is organized by purpose and audience for easy navigation.

---

## ðŸ—‚ï¸ Directory Structure

```
docs/
â”œâ”€â”€ README.md (this file)
â”œâ”€â”€ 01-06 Master Documentation (6 files)
â”œâ”€â”€ Core Guides (9 files)
â”œâ”€â”€ guides/ (Implementation guides - 6 files)
â”œâ”€â”€ design/ (Design & planning docs - 5 files)
â””â”€â”€ archives/ (Historical logs - 103 files)
```

### Root Directory
- **Master Docs:** 01-GETTING-STARTED through 06-PERSONAS
- **Quick References:** QUICK-REFERENCE, EXECUTIVE-SUMMARY
- **Deployment:** DEPLOYMENT-GUIDE, PRODUCTION-CHECKLIST, EXTERNAL-SERVICES
- **Systems:** NOTIFICATIONS-SYSTEM, SECURITY-IMPLEMENTATION
- **Config:** LOGIN-PASSWORD-DEFAULT

### Subdirectories
- **guides/** - Implementation guides for features
- **design/** - Design documents and planning specs
- **archives/** - Historical progress logs (not needed for current work)

---

## ðŸ“– Quick Navigation

### For New Developers
**Start Here:** [`01-GETTING-STARTED.md`](./01-GETTING-STARTED.md)
- System overview
- Local development setup
- First time run guide
- Common commands

### For Architects & Tech Leads
**Read:** [`02-ARCHITECTURE.md`](./02-ARCHITECTURE.md)
- System architecture
- Technology stack
- Database schema
- Module structure
- Design decisions

### For Product Managers
**Read:** [`03-FEATURES.md`](./03-FEATURES.md)
- Feature catalog
- User stories
- Feature flags
- Roadmap

### For Frontend/Backend Developers
**Read:** [`04-API-REFERENCE.md`](./04-API-REFERENCE.md)
- API endpoints
- Request/Response formats
- Authentication
- Error codes

### For DevOps Engineers
**Read:** [`05-DEPLOYMENT.md`](./05-DEPLOYMENT.md)
- Deployment guide
- Environment setup
- CI/CD pipeline
- Monitoring

### For Stakeholders
**Read:** [`06-PERSONAS.md`](./06-PERSONAS.md)
- User personas
- Use cases
- Value proposition
- Success metrics

---

## ðŸŽ¯ Key Documents

| Document | Purpose | Audience |
|----------|---------|----------|
| [Getting Started](./01-GETTING-STARTED.md) | Quick setup and overview | All developers |
| [Architecture](./02-ARCHITECTURE.md) | System design and structure | Architects, Tech Leads |
| [Features](./03-FEATURES.md) | Feature documentation | PM, Developers |
| [API Reference](./04-API-REFERENCE.md) | API specifications | Frontend, Backend devs |
| [Deployment](./05-DEPLOYMENT.md) | Deployment procedures | DevOps, SRE |
| [Personas](./06-PERSONAS.md) | User perspectives | Stakeholders, PM |

---

## ðŸ—ï¸ System Overview

### What is ServQR?
Intelligent Medical Equipment Service Management Platform with:
- **Multi-tenant architecture** for manufacturers, hospitals, and service providers
- **AI-powered diagnostics** for equipment troubleshooting
- **WhatsApp integration** for ticket creation
- **QR code system** for equipment tracking
- **Parts marketplace** (coming soon)
- **Field service management** for engineers

### Technology Stack
- **Backend:** Go, PostgreSQL, Redis
- **Frontend:** Next.js 14, React, TypeScript, Tailwind CSS
- **Infrastructure:** Docker, Kubernetes (optional)
- **AI:** OpenAI GPT-4, Claude 3, Whisper (STT)

---

## ðŸ“Š Platform Metrics

- **Modules:** 8 core modules (Tickets, Equipment, Organizations, Engineers, etc.)
- **APIs:** 50+ REST endpoints
- **Database:** 40+ tables
- **User Types:** 8 organization types supported
- **Features:** Multi-tenant, AI diagnostics, WhatsApp, Email notifications

---

## ðŸš€ Quick Start

```bash
# 1. Clone repository
git clone <repo-url>
cd ServQR

# 2. Setup environment
cp .env.example .env
# Edit .env with your values

# 3. Start database
cd dev/compose
docker-compose up -d postgres

# 4. Run migrations
psql -h localhost -p 5430 -U postgres -d med_platform -f database/migrations/*.sql

# 5. Start backend
go run cmd/platform/main.go

# 6. Start frontend (new terminal)
cd admin-ui
npm install
npm run dev
```

Access application at: http://localhost:3000

---

## ðŸ“‚ Archives

Historical progress logs, session summaries, and old documentation have been moved to [`archives/`](./archives/) to keep the main docs clean. These are useful for understanding project evolution but not required for current development.

---

## ðŸ”„ Documentation Updates

**Last Updated:** December 23, 2025  
**Version:** 2.0  
**Status:** Active Development

### Recent Changes
- Reorganized structure (Dec 2025)
- Added consolidated documentation files
- Archived old progress logs
- Added persona-based documentation

---

## ðŸ¤ Contributing to Docs

When adding documentation:
1. Determine the appropriate main document (01-06)
2. Add content in the relevant section
3. Update this README if adding new major sections
4. Use clear headings and examples
5. Keep it concise and actionable

---

## ðŸ“ž Support

- **Technical Issues:** Check troubleshooting sections in respective documents
- **API Questions:** See [API Reference](./04-API-REFERENCE.md)
- **Deployment Issues:** See [Deployment Guide](./05-DEPLOYMENT.md)

---

**Happy Coding! ðŸš€**
