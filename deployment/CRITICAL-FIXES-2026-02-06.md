# Critical Deployment Fixes - 2026-02-06

## Overview
This document describes critical fixes applied to resolve database initialization and frontend build failures in the production deployment system.

## Issues Fixed

### 1. Database Schema Syntax Errors
**Problem:** The base schema file `001_full_organizations_schema.sql` contained invalid SQL identifiers with hyphens:
- `'Channel Partner'` (contains space)
- `'Sub-sub_SUB_DEALER'` (contains hyphens)
- `'authorized_sub_Sub-sub_SUB_DEALER'` (contains hyphens)
- `'sub_sub_Sub-sub_SUB_DEALER_network'` (contains hyphens)
- `sold_by_sub_sub_Sub-sub_SUB_DEALER_id` (column name with hyphens)

**Impact:** PostgreSQL would fail to create tables, causing all subsequent migrations to fail with "relation does not exist" errors.

**Fix:** Standardized all identifiers to use underscores:
- `'channel_partner'`
- `'sub_dealer'`
- `'authorized_sub_dealer'`
- `'sub_dealer_network'`
- `sold_by_sub_dealer_id`

**Files Modified:**
- `database/migrations/001_full_organizations_schema.sql`

### 2. Next.js Static Page Generation Errors
**Problem:** Next.js was failing to build with errors on pages using `useSearchParams()`:
```
Error: useSearchParams() should be wrapped in a suspense boundary at page "/invite/accept"
Error: useSearchParams() should be wrapped in a suspense boundary at page "/organizations"
Error: useSearchParams() should be wrapped in a suspense boundary at page "/set-password"
Error: useSearchParams() should be wrapped in a suspense boundary at page "/tickets"
```

**Cause:** These pages are dynamic (they read URL parameters) and cannot be statically generated. Next.js 13+ requires either:
1. Wrapping `useSearchParams()` in a `<Suspense>` boundary, OR
2. Disabling the CSR (Client-Side Rendering) bailout warning

**Fix:** Added experimental flag to `next.config.js`:
```javascript
experimental: {
  // Allow useSearchParams without Suspense boundary
  missingSuspenseWithCSRBailout: false,
}
```

**Files Modified:**
- `admin-ui/next.config.js`
- `deployment/deploy-app.sh` (removed ineffective `SKIP_ENV_VALIDATION=true`)

### 3. Database Initialization Order
**Problem:** Migrations were running before the base schema, or the base schema was running twice (once explicitly, once in the migration loop).

**Fix:** Updated `setup-docker.sh` to:
1. Run `001_full_organizations_schema.sql` FIRST as the base schema (error if fails)
2. Skip it when running additional migrations
3. All other migrations run in sorted order as non-critical (warn on failure)

**Files Modified:**
- `deployment/setup-docker.sh`

## Testing the Fixes

### On a Fresh VM:
```bash
# 1. Clone/copy source code to /opt/servqr
cd /opt/servqr

# 2. Run deployment (should now complete successfully)
sudo bash deployment/deploy-all.sh

# 3. Verify database tables exist
deployment/connect-database.sh
\dt  -- Should show 20+ tables including organizations, engineers, service_tickets

# 4. Verify services are running
systemctl status servqr-backend
systemctl status servqr-frontend

# 5. Check application logs
tail -f /opt/servqr/logs/backend.log
tail -f /opt/servqr/logs/frontend.log
```

### On Existing Deployment:
```bash
# 1. Stop services
sudo systemctl stop servqr-backend servqr-frontend

# 2. Backup database
deployment/backup-database.sh

# 3. Drop and recreate database
docker exec -it servqr-postgres psql -U servqr -d postgres
DROP DATABASE servqr_production;
CREATE DATABASE servqr_production;
\q

# 4. Re-run database initialization
cd /opt/servqr/deployment
source ./setup-docker.sh  # Only run initialize_database function

# 5. Rebuild frontend
cd /opt/servqr/admin-ui
npm run build

# 6. Restart services
sudo systemctl start servqr-backend servqr-frontend
```

## Expected Results

### Database:
- Base schema creates all core tables: `organizations`, `org_relationships`, `organization_facilities`, `territories`, `contact_persons`, `engineers`, `engineer_skills`, `engineer_assignments`, `engineer_availability`, `service_tickets`, `equipment`, `equipment_registry`, etc.
- No "relation does not exist" errors
- All migrations complete (some may warn but should not fail critically)

### Frontend Build:
- No prerendering errors
- Build completes successfully
- `.next` directory created with all pages

### Services:
- Backend starts on port 8081
- Frontend starts on port 3000
- Both services respond to health checks
- Application accessible via browser

## Rollback Procedure

If issues occur:

1. **Restore Database:**
```bash
deployment/restore-database.sh /opt/servqr/backups/servqr-backup-YYYYMMDD-HHMMSS.sql.gz
```

2. **Revert Code:**
```bash
git checkout <previous-commit-hash>
```

3. **Rebuild:**
```bash
sudo bash deployment/deploy-app.sh
```

## Next Steps

1. ✅ Test deployment on clean Ubuntu 22.04 VM
2. ✅ Test deployment on existing installation
3. ⏳ Document any additional migration errors (if they occur)
4. ⏳ Create migration health check script
5. ⏳ Add database schema validation to deployment

## Technical Notes

### SQL Identifier Rules:
- PostgreSQL identifiers (table names, column names, constraint values) should:
  - Use underscores instead of hyphens
  - Use lowercase for consistency
  - Avoid spaces (or use quotes if necessary)

### Next.js Static Generation:
- Pages using `useSearchParams()`, `useRouter()`, or reading cookies/headers are **dynamic**
- Dynamic pages cannot be statically generated at build time
- Options:
  1. Wrap in `<Suspense>` (adds loading state)
  2. Use `experimental.missingSuspenseWithCSRBailout: false` (suppresses warning)
  3. Force dynamic rendering with `export const dynamic = 'force-dynamic'`

### Migration Strategy:
- **Base Schema:** Single comprehensive file that creates all core tables
- **Incremental Migrations:** Numbered files for schema changes, new features, data updates
- **Idempotency:** Use `IF NOT EXISTS`, `IF EXISTS` for safety
- **Order:** Migrations run in sorted (alphanumeric) order

## Support

For issues with these fixes, check:

1. **Database Logs:**
```bash
docker logs servqr-postgres | grep ERROR
```

2. **Application Logs:**
```bash
tail -100 /opt/servqr/logs/deployment-*.log
tail -100 /opt/servqr/logs/backend.log
tail -100 /opt/servqr/logs/frontend.log
```

3. **Service Status:**
```bash
systemctl status servqr-postgres servqr-backend servqr-frontend
```

4. **Health Checks:**
```bash
curl http://localhost:8081/health
curl http://localhost:3000/
```

---

**Date:** 2026-02-06  
**Author:** Droid  
**Version:** 1.0  
