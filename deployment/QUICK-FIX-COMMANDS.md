# Quick Fix Commands - Run These on Your Server

## The Problem
Your deployment completed but the backend is failing to start because the database schema has errors from previous incomplete migrations.

## The Solution (3 Simple Steps)

### Step 1: Pull Latest Fixes
```bash
cd /opt/servqr
git pull origin main
```

Expected output: Should show the latest commits with SQL syntax fixes.

### Step 2: Reset and Rebuild Database
```bash
sudo bash deployment/reset-and-rebuild-database.sh
```

This will:
- Backup your current database
- Drop and recreate it cleanly
- Apply the corrected base schema
- Run all migrations properly
- Restart services

When prompted "Are you sure you want to continue? (yes/no):", type: **yes**

### Step 3: Verify Services
```bash
# Check backend status
sudo systemctl status servqr-backend

# Check backend logs
sudo journalctl -u servqr-backend -n 50

# Test backend health
curl http://localhost:8081/health

# Test frontend
curl http://localhost:3000
```

Expected results:
- Backend should show "active (running)"
- Health check should return JSON response
- Frontend should return HTML

---

## If Backend Still Fails

### Check Backend Logs
```bash
sudo journalctl -u servqr-backend -n 200 --no-pager
```

Look for specific error messages like:
- "relation ... does not exist" → Database tables missing
- "connection refused" → Database not accessible
- "address already in use" → Port conflict

### Manual Backend Test
```bash
# Stop service
sudo systemctl stop servqr-backend

# Run manually to see errors
cd /opt/servqr
./platform

# If it works, press Ctrl+C and restart service
sudo systemctl start servqr-backend
```

### Check Database Tables
```bash
# Count tables (should be 30+)
docker exec servqr-postgres psql -U servqr -d servqr_production -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"

# List all tables
docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt"
```

---

## Alternative: Complete Clean Reinstall

If the reset script doesn't work, do a complete reinstall:

```bash
cd /opt/servqr

# Stop everything
sudo systemctl stop servqr-backend servqr-frontend

# Remove database data
sudo docker compose -f deployment/docker-compose.yml down -v
sudo rm -rf data/postgres/*

# Get latest code
git pull origin main

# Redeploy everything
sudo bash deployment/deploy-all.sh
```

---

## Quick Access After Fix

Once backend is running:

- **Frontend:** http://YOUR_SERVER_IP:3000
- **Backend API:** http://YOUR_SERVER_IP:8081/api/v1
- **Health Check:** http://YOUR_SERVER_IP:8081/health

Default login:
- Email: `admin@geneqr.com`
- Password: `admin123` (or check `/opt/servqr/LOGIN-CREDENTIALS.txt`)

---

## What Was Fixed

1. **SQL Syntax Errors:** Removed invalid hyphens from identifiers in:
   - `001_full_organizations_schema.sql`
   - `002_organizations_simple.sql`

2. **Next.js Build:** Added proper configuration for dynamic pages with `useSearchParams()`

3. **Migration Order:** Base schema now runs first, then additional migrations

All these fixes are in the latest commit (a1b0c571).
