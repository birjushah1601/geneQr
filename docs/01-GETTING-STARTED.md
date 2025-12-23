# Getting Started with ABY-MED Platform

## 🎯 Overview

ABY-MED is an intelligent medical equipment service management platform that connects manufacturers, hospitals, service providers, and field engineers through a unified system with AI-powered diagnostics and automation.

### Key Capabilities
- **Multi-tenant Architecture:** Separate data per organization
- **AI-Powered Diagnostics:** Automated equipment troubleshooting
- **WhatsApp Integration:** Create tickets via messaging
- **QR Code System:** Track and manage equipment
- **Field Service Management:** Engineer assignment and tracking
- **Parts Management:** Catalog and marketplace (coming soon)

---

## 🏗️ System Architecture (High-Level)

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │Dashboard │  │ Tickets  │  │Equipment │   ...       │
│  └──────────┘  └──────────┘  └──────────┘             │
└────────────────────┬────────────────────────────────────┘
                     │ REST API
┌────────────────────┴────────────────────────────────────┐
│              Backend (Go) - Port 8081                   │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Modular Architecture (8+ modules)               │  │
│  │  • Tickets  • Equipment  • Organizations        │  │
│  │  • Engineers • Parts • WhatsApp • AI • Auth     │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│        Database (PostgreSQL) - Port 5430                │
│  • Multi-tenant data isolation                          │
│  • 40+ tables with relationships                        │
│  • Audit logging • Event tracking                       │
└─────────────────────────────────────────────────────────┘
```

---

## 💻 Prerequisites

### Required Software
- **Go** 1.21+ ([Download](https://golang.org/dl/))
- **Node.js** 18+ and npm ([Download](https://nodejs.org/))
- **PostgreSQL** 15+ ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/downloads))

### Optional (Recommended)
- **Docker** & Docker Compose (for database)
- **VS Code** with Go and TypeScript extensions
- **Postman** or similar API testing tool

---

## 🚀 Quick Setup (10 Minutes)

### Step 1: Clone Repository
```bash
git clone <repository-url>
cd aby-med
```

### Step 2: Setup Environment
```bash
# Copy environment template
cp .env.example .env

# Edit .env and configure:
# - Database credentials
# - AI API keys (OpenAI, Anthropic)
# - Feature flags
```

**Minimum Required in .env:**
```bash
DB_HOST=localhost
DB_PORT=5430
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=med_platform

AI_OPENAI_API_KEY=sk-...your-key
OPENAI_API_KEY=sk-...your-key  # For Whisper STT

PORT=8081
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8081
```

### Step 3: Start Database

**Option A: Using Docker (Recommended)**
```bash
cd dev/compose
docker-compose up -d postgres
```

**Option B: Local PostgreSQL**
```bash
# Create database
createdb -U postgres med_platform

# Or with psql
psql -U postgres
CREATE DATABASE med_platform;
\q
```

### Step 4: Run Database Migrations
```bash
# Apply all migrations
psql -h localhost -p 5430 -U postgres -d med_platform -f database/migrations/001_init_schema.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f database/migrations/002_add_multi_tenant.sql
# ... apply all migration files in order

# Or use migration script if available
./scripts/migrate.sh
```

### Step 5: Start Backend
```bash
# From project root
go mod download
go run cmd/platform/main.go

# Backend will start on http://localhost:8081
```

**Verify Backend:**
```bash
curl http://localhost:8081/health
# Should return: {"status":"ok"}
```

### Step 6: Start Frontend
```bash
# Open new terminal
cd admin-ui
npm install
npm run dev

# Frontend will start on http://localhost:3000
```

### Step 7: Access Application
Open browser: **http://localhost:3000**

**Default Login (Test):**
- Email: `admin@example.com`
- Password: `admin123` (change in production!)

---

## 📁 Project Structure

```
aby-med/
├── cmd/
│   └── platform/
│       └── main.go                # Backend entry point
├── internal/
│   ├── service-domain/            # Core modules
│   │   ├── service-ticket/        # Ticket management
│   │   ├── equipment-registry/    # Equipment tracking
│   │   ├── organizations/         # Org management
│   │   ├── engineers/             # Engineer management
│   │   ├── whatsapp/              # WhatsApp integration
│   │   └── ...
│   ├── infrastructure/            # Cross-cutting concerns
│   │   ├── email/                 # Email notifications
│   │   ├── reports/               # Daily reports
│   │   └── audit/                 # Audit logging
│   └── shared/                    # Shared utilities
│       ├── middleware/            # Rate limiting, security
│       └── ...
├── admin-ui/
│   ├── src/
│   │   ├── app/                   # Next.js app router
│   │   │   ├── tickets/
│   │   │   ├── equipment/
│   │   │   ├── organizations/
│   │   │   └── ...
│   │   ├── components/            # React components
│   │   ├── lib/                   # API clients
│   │   └── types/                 # TypeScript types
│   └── package.json
├── database/
│   └── migrations/                # SQL migration files
├── docs/                          # Documentation (you are here)
├── .env.example                   # Environment template
└── go.mod                         # Go dependencies
```

---

## 🧪 Verify Setup

### Test Backend APIs
```bash
# Health check
curl http://localhost:8081/health

# List tickets (requires auth)
curl http://localhost:8081/api/v1/tickets

# Create test organization
curl -X POST http://localhost:8081/api/v1/organizations \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Hospital","org_type":"hospital"}'
```

### Test Frontend
1. Navigate to http://localhost:3000
2. Login with default credentials
3. Check dashboard loads
4. Navigate to Tickets, Equipment, Organizations

### Test Database Connection
```bash
psql -h localhost -p 5430 -U postgres -d med_platform

# Inside psql
\dt              # List tables
\d organizations # Describe organizations table
SELECT COUNT(*) FROM tickets;
\q
```

---

## 🔧 Common Commands

### Backend
```bash
# Run backend
go run cmd/platform/main.go

# Build backend
go build -o platform.exe ./cmd/platform

# Run tests
go test ./...

# Format code
go fmt ./...
```

### Frontend
```bash
# Development
npm run dev

# Build
npm run build

# Lint
npm run lint

# Type check
npx tsc --noEmit
```

### Database
```bash
# Connect to database
psql -h localhost -p 5430 -U postgres -d med_platform

# Backup database
pg_dump -h localhost -p 5430 -U postgres med_platform > backup.sql

# Restore database
psql -h localhost -p 5430 -U postgres med_platform < backup.sql

# Reset database (CAUTION!)
psql -h localhost -p 5430 -U postgres -c "DROP DATABASE med_platform;"
psql -h localhost -p 5430 -U postgres -c "CREATE DATABASE med_platform;"
```

---

## 🎓 Next Steps

After successful setup:

1. **Explore Features**
   - Read [03-FEATURES.md](./03-FEATURES.md) for feature overview
   - Check feature flags in `.env`

2. **Understand Architecture**
   - Read [02-ARCHITECTURE.md](./02-ARCHITECTURE.md)
   - Review database schema

3. **API Development**
   - Read [04-API-REFERENCE.md](./04-API-REFERENCE.md)
   - Import Postman collection from `postman/`

4. **Create Sample Data**
   ```bash
   # Run seed script
   go run scripts/seed_data.go
   
   # Or manually create via UI
   # Organizations → Add → Equipment → Add → Tickets → Create
   ```

5. **Enable Advanced Features**
   ```bash
   # Edit .env
   ENABLE_WHATSAPP=true
   ENABLE_AI_DIAGNOSIS=true
   FEATURE_EMAIL_NOTIFICATIONS=true
   ```

---

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check if port 8081 is in use
netstat -ano | findstr :8081  # Windows
lsof -i :8081                 # Mac/Linux

# Check database connection
psql -h localhost -p 5430 -U postgres -d med_platform -c "SELECT 1;"

# Check Go modules
go mod tidy
go mod download
```

### Frontend won't start
```bash
# Clear cache and reinstall
rm -rf node_modules .next
npm install
npm run dev

# Check Node version
node --version  # Should be 18+
```

### Database connection errors
```bash
# Check if PostgreSQL is running
docker ps | grep postgres  # If using Docker

# Verify connection details in .env
# Ensure DB_HOST, DB_PORT, DB_USER, DB_PASSWORD are correct
```

### Module not found errors
```bash
# Backend
go mod tidy
go mod download

# Frontend
npm install
```

---

## 📚 Additional Resources

- **Architecture:** [02-ARCHITECTURE.md](./02-ARCHITECTURE.md)
- **API Reference:** [04-API-REFERENCE.md](./04-API-REFERENCE.md)
- **Deployment:** [05-DEPLOYMENT.md](./05-DEPLOYMENT.md)
- **Features:** [03-FEATURES.md](./03-FEATURES.md)

---

## 🤝 Development Workflow

```bash
# 1. Pull latest changes
git pull origin main

# 2. Create feature branch
git checkout -b feature/your-feature

# 3. Make changes and test
# ... develop ...

# 4. Commit changes
git add .
git commit -m "feat: your feature description"

# 5. Push and create PR
git push origin feature/your-feature
```

---

**Ready to build! 🚀**

For questions or issues, refer to other documentation files or check the troubleshooting section.
