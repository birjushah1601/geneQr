# Fix Production Database - UTF-8 Import Error

## Issue

During production deployment, the database init script failed with:
```
ERROR: invalid byte sequence for encoding "UTF8": 0xff
```

This caused partial database creation - some columns like `customer_email` were not created, breaking the application.

## Root Cause

The original init script was created on Windows and had UTF-8 encoding issues that PostgreSQL on Linux rejected.

## Solution

A clean UTF-8 compliant init script has been created and pushed to main.

---

## Fix Steps for Production

### Step 1: Pull Latest Code

```bash
cd /opt/servqr
git pull origin main
```

### Step 2: Backup Existing Database

```bash
# Create backup
docker exec servqr-postgres pg_dump -U servqr servqr_production > /opt/servqr/backups/before-utf8-fix-$(date +%Y%m%d-%H%M%S).sql

# Verify backup created
ls -lh /opt/servqr/backups/
```

### Step 3: Stop Application Services

```bash
sudo systemctl stop servqr-backend servqr-frontend
```

### Step 4: Recreate Database

```bash
# Connect to PostgreSQL
docker exec servqr-postgres psql -U servqr -d postgres << 'EOF'

-- Terminate all connections
SELECT pg_terminate_backend(pid) 
FROM pg_stat_activity 
WHERE datname = 'servqr_production' 
  AND pid <> pg_backend_pid();

-- Drop and recreate database
DROP DATABASE IF EXISTS servqr_production;
CREATE DATABASE servqr_production OWNER servqr;

EOF
```

### Step 5: Import Clean Schema

```bash
# Import the new clean init script
docker exec -i servqr-postgres psql -U servqr -d servqr_production < /opt/servqr/database/migrations/001_init_schema.sql

# Check for errors
echo "Exit code: $?"
# Should show: Exit code: 0
```

### Step 6: Verify Database

```bash
# Check table count
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
SELECT COUNT(*) as table_count 
FROM information_schema.tables 
WHERE table_schema = 'public';
"
# Should show: 73 tables

# Verify customer_email column exists
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'service_tickets' 
  AND column_name = 'customer_email';
"
# Should show:
#  column_name    | data_type
# ----------------+------------------------
#  customer_email | character varying

# Check for any errors in logs
docker logs servqr-postgres --tail 50 | grep -i error
```

### Step 7: Restart Services

```bash
# Start services
sudo systemctl start servqr-backend servqr-frontend

# Wait 5 seconds
sleep 5

# Check status
sudo systemctl status servqr-backend servqr-frontend

# Both should show: active (running)
```

### Step 8: Test Application

```bash
# Test backend health
curl http://localhost:8081/health
# Should return: {"status":"healthy"}

# Test frontend
curl http://localhost:3001
# Should return HTML

# Test via Nginx
curl http://servqr.com
# Should return login page
```

**Test in browser:**
1. Visit: `http://servqr.com`
2. Try to create a service request
3. Should work without "customer_email" error

---

## Verification Checklist

After completing the fix:

- [ ] Database has 73 tables (not 65)
- [ ] `customer_email` column exists in `service_tickets` table
- [ ] No UTF-8 errors in PostgreSQL logs
- [ ] Backend service running
- [ ] Frontend service running
- [ ] Can access application at `http://servqr.com`
- [ ] Can create service request without errors
- [ ] Login works
- [ ] Dashboard loads

---

## Rollback (If Needed)

If something goes wrong, restore from backup:

```bash
# Stop services
sudo systemctl stop servqr-backend servqr-frontend

# Drop current database
docker exec servqr-postgres psql -U servqr -d postgres -c "DROP DATABASE servqr_production;"

# Recreate database
docker exec servqr-postgres psql -U servqr -d postgres -c "CREATE DATABASE servqr_production OWNER servqr;"

# Restore from backup
docker exec -i servqr-postgres psql -U servqr -d servqr_production < /opt/servqr/backups/before-utf8-fix-YYYYMMDD-HHMMSS.sql

# Restart services
sudo systemctl start servqr-backend servqr-frontend
```

---

## What Changed

**Old init script:**
- File: `database/migrations/001_init_schema_OLD_WITH_UTF8_ISSUES.sql`
- Size: ~1.4 MB
- Issue: Invalid UTF-8 byte sequences
- Result: Partial import, 65 tables, missing columns

**New init script:**
- File: `database/migrations/001_init_schema.sql`
- Size: ~1.4 MB
- Encoding: Clean UTF-8
- Result: Complete import, 73 tables, all columns present

**Key improvements:**
1. ✅ Explicit UTF-8 encoding specified during dump
2. ✅ No invalid byte sequences
3. ✅ Works on both Windows and Linux
4. ✅ Complete schema (73 tables vs 65)
5. ✅ All columns including `customer_email`

---

## Quick Fix (Alternative)

If you don't want to recreate the entire database, just add the missing column:

```bash
# Add customer_email column
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
ALTER TABLE service_tickets 
ADD COLUMN IF NOT EXISTS customer_email VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_service_tickets_customer_email 
ON service_tickets(customer_email);
"

# Restart backend
sudo systemctl restart servqr-backend
```

**Note:** This quick fix only adds the missing column. There may be other missing columns from the partial import. The full database recreation is recommended.

---

## Troubleshooting

### Issue: "Database does not exist"

```bash
# Check if database exists
docker exec servqr-postgres psql -U servqr -l | grep servqr_production

# If not, create it
docker exec servqr-postgres psql -U servqr -d postgres -c "CREATE DATABASE servqr_production OWNER servqr;"
```

### Issue: "Permission denied"

```bash
# Grant permissions
docker exec servqr-postgres psql -U postgres -d servqr_production -c "
GRANT ALL PRIVILEGES ON DATABASE servqr_production TO servqr;
GRANT ALL ON ALL TABLES IN SCHEMA public TO servqr;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO servqr;
"
```

### Issue: "Still getting UTF-8 errors"

```bash
# Check database encoding
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
SHOW server_encoding;
SHOW client_encoding;
"
# Both should show: UTF8

# If not, recreate with explicit encoding
docker exec servqr-postgres psql -U servqr -d postgres -c "
CREATE DATABASE servqr_production 
OWNER servqr 
ENCODING 'UTF8' 
LC_COLLATE 'en_US.UTF-8' 
LC_CTYPE 'en_US.UTF-8';
"
```

### Issue: "Import takes too long"

The import should take 10-30 seconds. If it takes longer:

```bash
# Check PostgreSQL logs
docker logs servqr-postgres -f

# Check if database is locked
docker exec servqr-postgres psql -U servqr -d postgres -c "
SELECT pid, usename, state, query 
FROM pg_stat_activity 
WHERE datname = 'servqr_production';
"
```

---

## Prevention

To prevent this issue in future:

1. **Always use explicit UTF-8 encoding when creating dumps:**
   ```bash
   pg_dump --encoding=UTF8 -U user -d database > dump.sql
   ```

2. **Test imports on target platform before production:**
   ```bash
   # Create test database
   createdb test_import
   
   # Test import
   psql -U user -d test_import < dump.sql
   
   # Check for errors
   echo $?
   ```

3. **Use Docker for consistent environments:**
   - Same PostgreSQL version locally and in production
   - Same encoding settings
   - Test on Linux containers even when developing on Windows

4. **Validate dumps before committing:**
   ```bash
   # Check for invalid UTF-8
   grep -axv '.*' dump.sql
   
   # Should show no output if all valid
   ```

---

## Summary

**Problem:** UTF-8 encoding error broke production database import
**Cause:** Windows-created dump had invalid byte sequences
**Fix:** Created clean UTF-8 dump, recreate production database
**Time:** 5-10 minutes downtime
**Risk:** Low (backup created before changes)

**Status:** ✅ Fixed and pushed to main (commit: 7b1cb96d)

---

**Last Updated:** February 9, 2026
**Tested On:** PostgreSQL 15 (Docker)
**Platforms:** Windows (development), Linux (production)
