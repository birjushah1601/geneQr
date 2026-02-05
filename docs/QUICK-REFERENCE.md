# ServQR Onboarding System - Quick Reference Card

## ðŸš€ Quick Start (60 seconds)

```bash
# 1. Start Database
docker-compose up -d postgres

# 2. Apply Equipment FK Migrations (Feb 2026)
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-01-maintenance.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-02-downtime.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-03-usage-logs.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-04-service-config.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-05-documents.sql
psql -h localhost -p 5430 -U postgres -d med_platform -f migrations/fix-equipment-fk-06-attachments.sql

# 3. Start Backend (port 8081)
./platform.exe

# 4. Start Frontend (port 3000)
cd admin-ui && npm run dev
```

## ðŸ“ Key URLs

- **Frontend Wizard:** http://localhost:3000/onboarding/wizard
- **Backend Health:** http://localhost:8081/health
- **Organizations API:** http://localhost:8081/api/v1/organizations/import
- **Equipment API:** http://localhost:8081/api/v1/equipment/catalog/import

## ðŸ”§ Environment Variables

```bash
# .env
ENABLE_ORG=true
ENABLE_EQUIPMENT=true
DB_HOST=localhost
DB_PORT=5430
DB_NAME=med_platform
```

## ðŸ“¤ API Testing

### Organizations Import
```bash
curl -X POST http://localhost:8081/api/v1/organizations/import \
  -F "csv_file=@templates/csv/organizations-import-template.csv" \
  -F "dry_run=true"
```

### Equipment Import
```bash
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/equipment-catalog-radiology-template.csv" \
  -F "dry_run=false"
```

## ðŸ“ File Locations

### Backend
- **APIs:** `internal/core/{organizations,equipment}/api/`
- **Modules:** `internal/core/{organizations,equipment}/module.go`
- **Migrations:** `database/migrations/028-030_*.sql`

### Frontend
- **Wizard:** `admin-ui/src/app/onboarding/wizard/page.tsx`
- **Components:** `admin-ui/src/components/onboarding/`
- **Templates:** `admin-ui/public/templates/*.csv`

### Documentation
- **README:** `docs/ONBOARDING-SYSTEM-README.md`
- **Deployment:** `docs/DEPLOYMENT-GUIDE.md`
- **Executive:** `docs/EXECUTIVE-SUMMARY.md`

## ðŸ­ Industry Templates

| Industry | File | Items |
|----------|------|-------|
| Radiology | `equipment-catalog-radiology-template.csv` | 8 |
| Cardiology | `equipment-catalog-cardiology-template.csv` | 8 |
| Surgical | `equipment-catalog-surgical-template.csv` | 8 |
| ICU | `equipment-catalog-icu-template.csv` | 8 |
| Laboratory | `equipment-catalog-laboratory-template.csv` | 8 |

## ðŸ—„ï¸ Database Queries

```sql
-- Check organizations
SELECT name, org_type FROM organizations ORDER BY created_at DESC LIMIT 10;

-- Check equipment
SELECT product_code, product_name, category FROM equipment_catalog ORDER BY created_at DESC LIMIT 10;

-- Check QR codes
SELECT * FROM qr_codes_unassigned LIMIT 10;

-- Check batch summary
SELECT * FROM qr_batches_summary;
```

## ðŸ› Troubleshooting

| Issue | Solution |
|-------|----------|
| Backend won't start | Check `.env` file, verify DB connection |
| Port 8081 in use | Kill process: `netstat -ano \| findstr :8081` |
| Migration fails | Check PostgreSQL running, verify connection |
| CSV upload fails | Verify file < 10MB, format is `.csv` |
| Import errors | Run dry run first, check error messages |

## ðŸ“Š Performance Benchmarks

| Operation | Expected Time |
|-----------|---------------|
| Org import (6 items) | < 1 second |
| Equipment import (8 items) | < 1 second |
| Equipment import (40 items) | < 3 seconds |
| Wizard completion | 5-10 minutes |

## âœ… Pre-Launch Checklist

- [ ] Migrations applied
- [ ] Backend compiles
- [ ] Frontend builds
- [ ] Feature flags enabled
- [ ] Templates accessible
- [ ] Health check returns 200
- [ ] Documentation reviewed

## ðŸŽ¯ Common Tasks

### Add New Industry Template
1. Create CSV in `templates/csv/equipment-catalog-{industry}-template.csv`
2. Copy to `admin-ui/public/templates/`
3. Update `EquipmentUploadStep.tsx` INDUSTRY_TEMPLATES array

### Test New Template
```bash
curl -X POST http://localhost:8081/api/v1/equipment/catalog/import \
  -F "csv_file=@templates/csv/your-template.csv" \
  -F "dry_run=true"
```

### Check Logs
```bash
# Backend logs (if redirected)
tail -f backend.log

# Database logs
docker logs med_platform_pg

# Frontend logs
# Check console in browser DevTools
```

## ðŸ“ˆ Metrics to Monitor

- Import success rate
- Import duration
- Error rate by endpoint
- User completion rate
- Template download count

## ðŸ” Security Notes

- File upload limit: 10MB
- CSV validation: automatic
- SQL injection: prevented (prepared statements)
- Transaction safety: enabled
- Feature flags: production-ready

## ðŸŽ¨ Wizard Steps

1. **Company Profile** - Organization details + validation
2. **Organizations Upload** - CSV import (optional, can skip)
3. **Equipment Catalog** - Industry templates (optional, can skip)
4. **Completion** - Success + statistics + next steps

## ðŸ’¡ Tips

- Always run **dry run first** to validate
- Use **industry templates** for fastest onboarding
- Check **error messages** for specific row issues
- **Skip steps** if not ready (flexibility)
- Download **templates** for correct format

## ðŸ“ž Quick Help

- **Full Documentation:** `docs/ONBOARDING-SYSTEM-README.md`
- **Deployment Guide:** `docs/DEPLOYMENT-GUIDE.md`
- **API Details:** `docs/ONBOARDING-SYSTEM-README.md#api-endpoints`
- **Troubleshooting:** `docs/DEPLOYMENT-GUIDE.md#troubleshooting`

---

**Version:** 1.0  
**Last Updated:** December 23, 2025  
**Status:** âœ… Production Ready
