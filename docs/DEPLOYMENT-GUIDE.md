# ServQR Onboarding System - Deployment Guide

Quick deployment guide for the complete onboarding system.

## ðŸš€ Quick Start

### Prerequisites
- PostgreSQL 15+ (running on port 5430 via Docker)
- Go 1.21+
- Node.js 18+
- Git

### 1. Database Setup

#### Start PostgreSQL (if using Docker)
```bash
cd dev/compose
docker-compose up -d postgres
```

#### Apply Migrations
```bash
# Navigate to project root
cd /path/to/ServQR

# Apply equipment FK migrations (Feb 2026)
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-01-maintenance.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-02-downtime.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-03-usage-logs.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-04-service-config.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-05-documents.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-06-attachments.sql

# Verify migrations
psql -h localhost -p 5430 -U postgres -d med_platform -c "\dt qr_*"
```

### 2. Backend Setup

#### Configure Environment
```bash
# Copy .env.example to .env
cp .env.example .env

# Edit .env and ensure these are set:
# ENABLE_ORG=true
# ENABLE_EQUIPMENT=true
# DB_HOST=localhost
# DB_PORT=5430
# DB_NAME=med_platform
```

#### Build & Run Backend
```bash
# Build
go build -o platform.exe ./cmd/platform

# Run
./platform.exe

# Or directly
go run cmd/platform/main.go
```

#### Verify Backend
```bash
# Check health
curl http://localhost:8081/health

# Expected output: {"status":"ok"}
```

### 3. Frontend Setup

#### Install Dependencies
```bash
cd admin-ui
npm install
```

#### Run Development Server
```bash
npm run dev

# Frontend will be available at:
# http://localhost:3000
```

#### Access Onboarding Wizard
```
http://localhost:3000/onboarding/wizard
```

## ðŸ§ª Testing

### Test Organizations Import

#### Dry Run (Validation)
```bash
curl -X POST http://localhost:8081/api/v1/organizations/import \
  -F "csv_file=@templates/csv/organizations-import-template.csv" \
  -F "dry_run=true" \
  -F "created_by=test-user"
```

#### Actual Import
```bash
curl -X POST http://localhost:8081/api/v1/organizations/import \
  -F "csv_file=@templates/csv/organizations-import-template.csv" \
  -F "dry_run=false" \
  -F "created_by=test-user"
```

### Test Equipment Import

#### Test Radiology Template
```bash
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-radiology-template.csv" \
  -F "dry_run=true" \
  -F "created_by=test-user"
```

#### Test All Industry Templates
```bash
# Cardiology
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-cardiology-template.csv" \
  -F "dry_run=false"

# Surgical
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-surgical-template.csv" \
  -F "dry_run=false"

# ICU
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-icu-template.csv" \
  -F "dry_run=false"

# Laboratory
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-laboratory-template.csv" \
  -F "dry_run=false"
```

### Verify Data in Database

#### Check Organizations
```sql
psql -h localhost -p 5430 -U postgres -d med_platform

SELECT name, org_type, status FROM organizations ORDER BY created_at DESC LIMIT 10;
```

#### Check Equipment Catalog
```sql
SELECT product_code, product_name, category, manufacturer_name 
FROM equipment_catalog 
ORDER BY created_at DESC LIMIT 10;
```

#### Check QR Codes
```sql
-- View unassigned QR codes
SELECT * FROM qr_codes_unassigned LIMIT 10;

-- View batch summary
SELECT * FROM qr_batches_summary;
```

## ðŸ“‹ Frontend Testing Checklist

### Wizard Flow
- [ ] Step 1: Company Profile
  - [ ] Form validation works
  - [ ] Required fields highlighted
  - [ ] Email format validation
  - [ ] Can proceed to next step

- [ ] Step 2: Organizations Upload
  - [ ] Template download works
  - [ ] Drag-and-drop file upload
  - [ ] Dry run validation
  - [ ] Error display (if any)
  - [ ] Success message
  - [ ] Can skip or continue

- [ ] Step 3: Equipment Catalog
  - [ ] Industry selector displays
  - [ ] All 5 industries visible
  - [ ] Template download for selected industry
  - [ ] CSV upload works
  - [ ] Validation before import
  - [ ] Can skip or continue

- [ ] Step 4: Completion
  - [ ] Success message displays
  - [ ] Statistics show correct counts
  - [ ] Next steps guide visible
  - [ ] "Go to Dashboard" button works

### Component Testing
- [ ] CSV Uploader
  - [ ] Drag-and-drop zone responds
  - [ ] File validation (size, format)
  - [ ] Progress indicators work
  - [ ] Error messages clear

- [ ] Industry Selector
  - [ ] All 5 industries display
  - [ ] Visual selection works
  - [ ] Icons render correctly
  - [ ] Selected state highlights

## ðŸ”§ Troubleshooting

### Backend Won't Start
```bash
# Check if port is already in use
netstat -ano | findstr :8081

# Check database connection
psql -h localhost -p 5430 -U postgres -d med_platform -c "SELECT 1;"

# Check environment variables
cat .env | grep ENABLE_
```

### Frontend Build Errors
```bash
# Clear cache and reinstall
cd admin-ui
rm -rf node_modules .next
npm install
npm run dev
```

### CSV Import Fails
- Verify file format is CSV
- Check file size < 10MB
- Ensure required columns exist
- Run dry run first to see errors
- Check backend logs for details

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Test connection
psql -h localhost -p 5430 -U postgres -d med_platform -c "\dt"
```

## ðŸ“Š Performance Benchmarks

### Expected Performance
- Organizations Import (6 items): < 1 second
- Equipment Import (8 items): < 1 second
- Equipment Import (40 items, all templates): < 3 seconds
- Wizard completion time: 5-10 minutes

### Database Queries
- QR code generation: < 10ms per code
- CSV parsing: ~1000 rows/second
- Validation: ~5000 rows/second

## ðŸ”’ Security Checklist

- [ ] Database passwords not in .env file committed to git
- [ ] API endpoints require authentication (production)
- [ ] File upload size limits enforced (10MB)
- [ ] CSV content sanitized
- [ ] SQL injection prevention (prepared statements)
- [ ] Feature flags enabled only for authorized users

## ðŸš€ Production Deployment

### Build for Production

#### Backend
```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o platform ./cmd/platform

# Or for Windows
go build -o platform.exe ./cmd/platform
```

#### Frontend
```bash
cd admin-ui
npm run build

# Output will be in .next/ directory
```

### Environment Variables (Production)

```bash
# .env.production
APP_ENV=production
PORT=8081
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=med_platform
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
ENABLE_ORG=true
ENABLE_EQUIPMENT=true
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

### Docker Deployment (Optional)

```dockerfile
# Dockerfile.backend
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o platform ./cmd/platform

FROM alpine:latest
COPY --from=builder /app/platform /platform
COPY --from=builder /app/templates /templates
CMD ["/platform"]
```

```dockerfile
# Dockerfile.frontend
FROM node:18-alpine AS builder
WORKDIR /app
COPY admin-ui/package*.json ./
RUN npm install
COPY admin-ui/ ./
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./
RUN npm install --production
CMD ["npm", "start"]
```

## ðŸ“ˆ Monitoring

### Health Checks
- Backend: `GET /health` (every 30s)
- Database: Connection pool status
- Frontend: Page load time

### Metrics to Track
- Import success rate
- Import duration
- Error rate by endpoint
- User completion rate (wizard)
- Template download count

## âœ… Pre-Launch Checklist

- [ ] All migrations applied successfully
- [ ] Backend compiles without errors
- [ ] Frontend builds without errors
- [ ] All 5 industry templates tested
- [ ] Organizations import tested
- [ ] Equipment import tested
- [ ] Database indexes created
- [ ] Feature flags configured
- [ ] Documentation complete
- [ ] Security review completed
- [ ] Performance benchmarks met
- [ ] Backup strategy in place

## ðŸŽ‰ Success Criteria

- [ ] Backend health check returns 200
- [ ] Frontend wizard loads in < 2 seconds
- [ ] CSV import completes in < 3 seconds
- [ ] Zero SQL errors in logs
- [ ] All tests pass
- [ ] Documentation accessible
- [ ] Templates downloadable

---

**Status**: Ready for Production Deployment

**Last Updated**: December 23, 2025

For issues or questions, refer to the [Onboarding System README](./ONBOARDING-SYSTEM-README.md).
