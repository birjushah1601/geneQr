# ServQR Deployment - Critical Fixes

## 🚨 Issues Found

Your deployment is failing at two points:

### Issue 1: Database Migrations Failing
**Error:** `ERROR: relation "service_tickets" does not exist`

**Root Cause:** Migrations are trying to modify tables that don't exist yet because:
- The base schema wasn't run first
- Script is looking in `/opt/servqr/migrations` but schemas are in `/opt/servqr/database/migrations`

### Issue 2: Backend Build Failing
**Error:** `undefined: initAuthModule` and `undefined: initNotificationsAndReports`

**Root Cause:** Missing files in the build directory:
- `cmd/platform/init_auth.go`
- `cmd/platform/init_notifications.go`

---

## 🔧 **SOLUTION: Run These Commands on Your Server**

### Quick Fix (All-in-One)

```bash
cd /opt/servqr

# Pull latest code (if you pushed fixes)
git pull origin main

# Run the fix script
cd deployment
sudo bash fix-migrations.sh

# Then continue with deployment
sudo bash deploy-app.sh
```

---

## 📝 **Detailed Fix Steps**

### Step 1: Fix Database

```bash
# Stop and remove old database
sudo docker stop servqr-postgres
sudo docker rm servqr-postgres
sudo rm -rf /opt/servqr/data/postgres/*

# Start fresh database
cd /opt/servqr/deployment
sudo docker compose up -d postgres

# Wait for it to be ready
sleep 10

# Run base schema
sudo docker exec -i servqr-postgres psql -U servqr -d servqr_production < /opt/servqr/init-database-schema.sql

# If that file doesn't exist, try:
sudo docker exec -i servqr-postgres psql -U servqr -d servqr_production < /opt/servqr/database/migrations/001_full_organizations_schema.sql

# Verify tables created
sudo docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt"
```

### Step 2: Fix Backend Build

Check if init files exist:

```bash
cd /opt/servqr/cmd/platform
ls -la init*.go
```

**If files are missing:**

```bash
# Pull latest code
cd /opt/servqr
git pull origin main

# Or check if they're in a different location
find /opt/servqr -name "init_auth.go"
find /opt/servqr -name "init_notifications.go"
```

**If they exist but build still fails:**

```bash
# Try building manually to see the full error
cd /opt/servqr
go build -v ./cmd/platform/main.go
```

### Step 3: Continue Deployment

```bash
cd /opt/servqr/deployment
sudo bash deploy-app.sh
```

---

## 🎯 **Alternative: Complete Fresh Start**

If the above doesn't work, do a complete fresh installation:

```bash
# 1. Stop everything
sudo systemctl stop servqr-backend servqr-frontend 2>/dev/null || true
sudo docker stop servqr-postgres 2>/dev/null || true
sudo docker rm servqr-postgres 2>/dev/null || true

# 2. Clean data
sudo rm -rf /opt/servqr/data/postgres/*
sudo rm -rf /opt/servqr/.db_password
sudo rm -rf /opt/servqr/.jwt_secret
sudo rm -rf /opt/servqr/.env
sudo rm -rf /opt/servqr/admin-ui/.env.local

# 3. Pull latest code
cd /opt/servqr
git pull origin main

# 4. Run deployment from scratch
cd deployment
sudo bash deploy-all.sh
```

---

## 🔍 **Verification Commands**

### Check Database

```bash
# Connect to database
sudo docker exec -it servqr-postgres psql -U servqr -d servqr_production

# List all tables
\dt

# Should see tables like:
# - organizations
# - users
# - service_tickets
# - equipment_registry
# - engineers
# - spare_parts_catalog

# Exit
\q
```

### Check Backend Build

```bash
# Check if init files exist
ls -la /opt/servqr/cmd/platform/init*.go

# Try manual build
cd /opt/servqr
go build -o test-platform ./cmd/platform/main.go

# If successful, binary created
ls -lh test-platform
```

---

## 📊 **Expected Migration Order**

Migrations should run in this order:

1. **Base Schema** (creates all tables)
   - `001_full_organizations_schema.sql` OR
   - `init-database-schema.sql`

2. **Additional Migrations** (modify existing tables)
   - All other `.sql` files in `database/migrations/`

3. **Demo Data** (optional)
   - `*demo*.sql` files

---

## 🆘 **If Still Failing**

### Check These Files Exist

```bash
# On the server
ls -la /opt/servqr/init-database-schema.sql
ls -la /opt/servqr/database/migrations/001_full_organizations_schema.sql
ls -la /opt/servqr/cmd/platform/init_auth.go
ls -la /opt/servqr/cmd/platform/init_notifications.go
ls -la /opt/servqr/cmd/platform/main.go
```

### Get Full Error Output

```bash
# Run deployment with full output
cd /opt/servqr/deployment
sudo bash deploy-all.sh 2>&1 | tee deployment-debug.log

# Send the deployment-debug.log file for analysis
```

### Check Go Module Issues

```bash
cd /opt/servqr
go mod tidy
go mod download
```

---

## 💡 **Quick Diagnosis**

Run this to check your setup:

```bash
#!/bin/bash
echo "=== ServQR Deployment Diagnosis ==="
echo ""
echo "1. Database Status:"
docker ps | grep servqr-postgres || echo "   NOT RUNNING"
echo ""
echo "2. Init Schema Files:"
ls -lh /opt/servqr/init-database-schema.sql 2>/dev/null || echo "   NOT FOUND"
ls -lh /opt/servqr/database/migrations/001_full_organizations_schema.sql 2>/dev/null || echo "   NOT FOUND"
echo ""
echo "3. Backend Init Files:"
ls -lh /opt/servqr/cmd/platform/init_auth.go 2>/dev/null || echo "   NOT FOUND"
ls -lh /opt/servqr/cmd/platform/init_notifications.go 2>/dev/null || echo "   NOT FOUND"
echo ""
echo "4. Go Version:"
go version
echo ""
echo "5. Node Version:"
node -v
echo ""
echo "6. Docker Version:"
docker --version
echo ""
```

Save this as `diagnose.sh`, run it, and share the output.

---

## 📞 **Need More Help?**

Provide this information:

1. **Full deployment log:**
   ```bash
   cat /opt/servqr/logs/deployment-*.log | tail -100
   ```

2. **Database status:**
   ```bash
   docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt"
   ```

3. **File structure:**
   ```bash
   ls -R /opt/servqr/cmd/platform/
   ls -R /opt/servqr/database/migrations/ | head -20
   ```

4. **Build errors:**
   ```bash
   cd /opt/servqr
   go build -v ./cmd/platform/main.go 2>&1
   ```

---

**Quick Command Summary:**

```bash
# Fix database
cd /opt/servqr/deployment
sudo bash fix-migrations.sh

# Continue deployment
sudo bash deploy-app.sh

# Or start fresh
sudo bash deploy-all.sh
```
