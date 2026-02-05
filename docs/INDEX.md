# Documentation Index

**Last Updated:** February 5, 2026

---

## 📚 Quick Start

1. [Getting Started](01-GETTING-STARTED.md) - Setup and installation
2. [Quick Reference](QUICK-REFERENCE.md) - Common commands and tasks
3. [Architecture Overview](02-ARCHITECTURE.md) - System architecture

---

## 🎯 Core Documentation

### Architecture & Design
- [02-ARCHITECTURE.md](02-ARCHITECTURE.md) - Complete system architecture
- [EQUIPMENT-ARCHITECTURE-FINAL.md](EQUIPMENT-ARCHITECTURE-FINAL.md) ⭐ **NEW** - Equipment tables explained
- [EQUIPMENT-RELATIONSHIPS-DIAGRAM.md](EQUIPMENT-RELATIONSHIPS-DIAGRAM.md) ⭐ **NEW** - Critical FK relationships
- [EXECUTIVE-SUMMARY.md](EXECUTIVE-SUMMARY.md) - Executive overview

### Features
- [03-FEATURES.md](03-FEATURES.md) - All features overview
- [PARTNER-ENGINEERS-FEATURE.md](PARTNER-ENGINEERS-FEATURE.md) ⭐ **NEW** - Partner engineers implementation
- [SERVICE-REQUEST-ENHANCEMENTS.md](SERVICE-REQUEST-ENHANCEMENTS.md) ⭐ **NEW** - Service request contact fields
- [NOTIFICATIONS-SYSTEM.md](NOTIFICATIONS-SYSTEM.md) - Notification system

### API Reference
- [04-API-REFERENCE.md](04-API-REFERENCE.md) - API overview
- [api/ASSIGNMENT-API.md](api/ASSIGNMENT-API.md) - Engineer assignment API (includes include_partners)
- [api/ATTACHMENT-API.md](api/ATTACHMENT-API.md) - Attachment handling API

### Deployment & Testing
- [05-DEPLOYMENT.md](05-DEPLOYMENT.md) - Deployment guide
- [05-TESTING.md](05-TESTING.md) - Testing procedures
- [DEPLOYMENT-GUIDE.md](DEPLOYMENT-GUIDE.md) - Detailed deployment
- [PRODUCTION-DEPLOYMENT-CHECKLIST.md](PRODUCTION-DEPLOYMENT-CHECKLIST.md) - Production checklist

### Personas & Users
- [06-PERSONAS.md](06-PERSONAS.md) - User personas and workflows

---

## 🔧 Implementation Guides

### Equipment & QR Codes
- [guides/qr-code-setup.md](guides/qr-code-setup.md) - QR code setup
- [specs/QR-CODE-MIGRATION-PLAN.md](specs/QR-CODE-MIGRATION-PLAN.md) - QR migration (✅ COMPLETED)

### Engineer Management
- [guides/engineer-management.md](guides/engineer-management.md) - Engineer management
- [guides/SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md](guides/SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md) - Assignment logic

### Multi-Tenancy & Organizations
- [guides/MULTI-TENANT-IMPLEMENTATION-PLAN.md](guides/MULTI-TENANT-IMPLEMENTATION-PLAN.md) - Multi-tenant architecture
- [specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md](specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md) - Organizations design
- [specs/PARTNER-ASSOCIATION-SPECIFICATION.md](specs/PARTNER-ASSOCIATION-SPECIFICATION.md) - Partner relationships

### Data Import
- [guides/csv-imports.md](guides/csv-imports.md) - CSV import guide
- [template_csv/INDEX.md](template_csv/INDEX.md) - CSV templates index
- [template_csv/README.md](template_csv/README.md) - Template documentation

### Tickets & Workflows
- [guides/TICKET-ENHANCEMENTS-IMPLEMENTATION.md](guides/TICKET-ENHANCEMENTS-IMPLEMENTATION.md) - Ticket features
- [guides/OPTION3-WHATSAPP-IMPLEMENTATION-GUIDE.md](guides/OPTION3-WHATSAPP-IMPLEMENTATION-GUIDE.md) - WhatsApp integration

### Onboarding
- [guides/ONBOARDING-SYSTEM-README.md](guides/ONBOARDING-SYSTEM-README.md) - Onboarding system

---

## 🔐 Security

- [SECURITY-IMPLEMENTATION-COMPLETE.md](SECURITY-IMPLEMENTATION-COMPLETE.md) - Security implementation
- [specs/SECURITY-CHECKLIST.md](specs/SECURITY-CHECKLIST.md) - Security checklist
- [EXTERNAL-SERVICES-SETUP.md](EXTERNAL-SERVICES-SETUP.md) - External services (SendGrid, Twilio, etc.)

---

## 📋 Specifications

- [specs/API-SPECIFICATION.md](specs/API-SPECIFICATION.md) - Complete API spec
- [specs/SPECIFICATION-SUMMARY.md](specs/SPECIFICATION-SUMMARY.md) - Spec summary
- [specs/PARTNER-ASSOCIATION-SPECIFICATION.md](specs/PARTNER-ASSOCIATION-SPECIFICATION.md) - Partner relationships

---

## 🎨 Design Documents

- [design/AUTHENTICATION-MULTITENANCY-PRD.md](design/AUTHENTICATION-MULTITENANCY-PRD.md) - Auth & multitenancy PRD
- [design/MANUFACTURER-ONBOARDING-UX-DESIGN.md](design/MANUFACTURER-ONBOARDING-UX-DESIGN.md) - Onboarding UX
- [design/MARKETPLACE-BRAINSTORMING.md](design/MARKETPLACE-BRAINSTORMING.md) - Marketplace ideas
- [design/ONBOARDING-SYSTEM-BRAINSTORM.md](design/ONBOARDING-SYSTEM-BRAINSTORM.md) - Onboarding brainstorm
- [design/QR-CODE-TABLE-DESIGN-ANALYSIS.md](design/QR-CODE-TABLE-DESIGN-ANALYSIS.md) - QR design decisions

---

## 📊 Status & Tracking

- [CHANGELOG-2026-02-05.md](../CHANGELOG-2026-02-05.md) ⭐ **NEW** - Latest session changes
- [DOCUMENTATION-AUDIT-2026-02-05.md](DOCUMENTATION-AUDIT-2026-02-05.md) ⭐ **NEW** - Documentation audit
- [LOGIN-PASSWORD-DEFAULT.md](LOGIN-PASSWORD-DEFAULT.md) - Default credentials

---

## 📁 Archived Documentation

See [archived/](archived/) folder for:
- Historical implementation status
- Completed migration documentation
- Old README versions
- Past project status reports

---

## 🆕 Recent Additions (February 5, 2026)

### Equipment Architecture
- **EQUIPMENT-ARCHITECTURE-FINAL.md** - Comprehensive explanation of equipment vs equipment_registry tables
- **EQUIPMENT-RELATIONSHIPS-DIAGRAM.md** - FK relationships and data flow

### Features
- **PARTNER-ENGINEERS-FEATURE.md** - Complete partner engineers implementation guide
- **SERVICE-REQUEST-ENHANCEMENTS.md** - Contact fields for service requests

### Project Management
- **CHANGELOG-2026-02-05.md** - Session summary (19 commits)
- **DOCUMENTATION-AUDIT-2026-02-05.md** - Comprehensive doc audit

---

## 🔍 Finding Documentation

### By Topic

**Equipment Management:**
- EQUIPMENT-ARCHITECTURE-FINAL.md
- EQUIPMENT-RELATIONSHIPS-DIAGRAM.md
- guides/csv-imports.md

**Engineer Assignment:**
- PARTNER-ENGINEERS-FEATURE.md
- guides/engineer-management.md
- guides/SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md
- api/ASSIGNMENT-API.md

**QR Codes:**
- guides/qr-code-setup.md
- specs/QR-CODE-MIGRATION-PLAN.md (completed)

**Service Requests:**
- SERVICE-REQUEST-ENHANCEMENTS.md
- guides/TICKET-ENHANCEMENTS-IMPLEMENTATION.md

**Multi-Tenancy:**
- guides/MULTI-TENANT-IMPLEMENTATION-PLAN.md
- specs/DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md
- design/AUTHENTICATION-MULTITENANCY-PRD.md

**Notifications:**
- NOTIFICATIONS-SYSTEM.md
- guides/FEATURE-FLAGS-NOTIFICATIONS.md
- guides/OPTION3-WHATSAPP-IMPLEMENTATION-GUIDE.md

---

## 📖 Documentation Standards

### Status Indicators
- ⭐ **NEW** - Created in latest session
- ✅ **COMPLETED** - Feature/migration completed
- ⚠️ **NEEDS UPDATE** - Requires review/update
- 📁 **ARCHIVED** - Historical reference only

### Document Types
- **Guides** - Step-by-step implementation
- **Specs** - Technical specifications
- **Design** - Design decisions and brainstorming
- **API** - API endpoint documentation
- **Status** - Project tracking and changelogs

---

## 🤝 Contributing

When adding new documentation:
1. Add entry to this index
2. Use consistent formatting
3. Add status indicator if needed
4. Update "Last Updated" date
5. Link related documents

---

**Questions?** See [README.md](README.md) for overview or [01-GETTING-STARTED.md](01-GETTING-STARTED.md) for setup.
