# Backend Troubleshooting Guide

## Quick Diagnosis

Run these commands on your server to diagnose the backend issue:

```bash
# 1. Check backend service status
sudo systemctl status servqr-backend

# 2. View backend logs (last 100 lines)
sudo journalctl -u servqr-backend -n 100 --no-pager

# 3. Check if database tables exist
docker exec servqr-postgres psql -U servqr -d servqr_production -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"

# 4. List all tables
docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt"

# 5. Check environment variables
cat /opt/servqr/.env | head -20
```

## Common Issues and Solutions

### Issue 1: Database Tables Not Created

**Symptoms:**
- Backend fails with errors like "relation does not exist"
- Health check fails
- Migration errors during deployment

**Solution:**
```bash
cd /opt/servqr
sudo bash deployment/reset-and-rebuild-database.sh
```

This will:
1. Backup existing database
2. Drop and recreate the database
3. Apply base schema cleanly
4. Run all migrations
5. Restart services

### Issue 2: Environment Variables Missing

**Symptoms:**
- Backend logs show "DATABASE_URL not found"
- Connection errors

**Solution:**
```bash
# Check if .env file exists
ls -la /opt/servqr/.env

# Verify database password
cat /opt/servqr/.db_password

# Verify JWT secret
cat /opt/servqr/.jwt_secret

# If missing, regenerate:
cd /opt/servqr
sudo bash deployment/deploy-app.sh
```

### Issue 3: Port Already in Use

**Symptoms:**
- Backend logs show "address already in use"
- Backend exits immediately

**Solution:**
```bash
# Check what's using port 8081
sudo lsof -i :8081

# If another process is using it, kill it:
sudo kill -9 <PID>

# Or change the port in .env:
sudo vim /opt/servqr/.env
# Change: PORT=8082

# Restart backend
sudo systemctl restart servqr-backend
```

### Issue 4: Database Connection Failed

**Symptoms:**
- Backend logs show "connection refused"
- "could not connect to database"

**Solution:**
```bash
# 1. Check if PostgreSQL container is running
docker ps | grep servqr-postgres

# 2. If not running, start it:
cd /opt/servqr/deployment
docker compose up -d postgres

# 3. Check database is accessible
docker exec servqr-postgres psql -U servqr -d servqr_production -c "SELECT version();"

# 4. Verify DATABASE_URL in .env matches database password
cat /opt/servqr/.db_password
cat /opt/servqr/.env | grep DATABASE_URL
```

### Issue 5: Permission Errors

**Symptoms:**
- Backend logs show "permission denied"
- Cannot read/write files

**Solution:**
```bash
# Fix ownership
sudo chown -R root:root /opt/servqr
sudo chmod +x /opt/servqr/platform

# Fix log directory permissions
sudo mkdir -p /opt/servqr/logs
sudo chmod 755 /opt/servqr/logs

# Restart backend
sudo systemctl restart servqr-backend
```

## Step-by-Step Debugging

### 1. Start with a clean slate

```bash
cd /opt/servqr
git pull origin main
sudo bash deployment/reset-and-rebuild-database.sh
```

### 2. Check database schema

```bash
# Count tables
docker exec servqr-postgres psql -U servqr -d servqr_production -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"

# Should show 30+ tables. If less than 20, schema didn't apply correctly.

# Check critical tables exist
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('organizations', 'users', 'service_tickets', 'equipment', 'engineers')
ORDER BY table_name;"
```

### 3. Test backend manually

```bash
# Stop systemd service
sudo systemctl stop servqr-backend

# Run backend manually to see full error output
cd /opt/servqr
./platform

# If it starts successfully, check health:
curl http://localhost:8081/health

# Press Ctrl+C to stop

# Start service again
sudo systemctl start servqr-backend
```

### 4. Check logs in real-time

```bash
# Terminal 1: Backend logs
sudo journalctl -u servqr-backend -f

# Terminal 2: Database logs
docker logs servqr-postgres -f

# Terminal 3: Restart backend
sudo systemctl restart servqr-backend
```

## Quick Fixes

### Reset everything and start fresh

```bash
# Complete reset
cd /opt/servqr
sudo systemctl stop servqr-backend servqr-frontend
sudo docker compose -f deployment/docker-compose.yml down -v
sudo rm -rf data/postgres/*
git pull origin main
sudo bash deployment/deploy-all.sh
```

### Rebuild just the backend

```bash
cd /opt/servqr
go build -o platform ./cmd/platform/
sudo systemctl restart servqr-backend
```

### Check backend configuration

```bash
# View full environment
sudo cat /opt/servqr/.env

# Critical variables to verify:
# - DATABASE_URL (should match .db_password)
# - JWT_SECRET (should exist)
# - PORT=8081
# - GIN_MODE=release
```

## Health Check Commands

```bash
# Backend health (should return 200 OK)
curl -v http://localhost:8081/health

# Database health
docker exec servqr-postgres pg_isready -U servqr -d servqr_production

# Frontend health (should return HTML)
curl -v http://localhost:3000

# All services status
sudo systemctl status servqr-postgres servqr-backend servqr-frontend
```

## Log Locations

- **Backend logs:** `journalctl -u servqr-backend`
- **Frontend logs:** `journalctl -u servqr-frontend`
- **Database logs:** `docker logs servqr-postgres`
- **Deployment logs:** `/opt/servqr/logs/deployment-*.log`
- **Systemd service file:** `/etc/systemd/system/servqr-backend.service`

## Getting Help

When asking for help, provide:

1. **Backend logs:**
```bash
sudo journalctl -u servqr-backend -n 200 --no-pager > backend-logs.txt
```

2. **Database table count:**
```bash
docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt" > db-tables.txt
```

3. **Service status:**
```bash
sudo systemctl status servqr-backend servqr-frontend servqr-postgres > services-status.txt
```

4. **Environment (without secrets):**
```bash
cat /opt/servqr/.env | grep -v PASSWORD | grep -v SECRET | grep -v KEY > env-config.txt
```

Send these 4 files for diagnosis.

---

## Most Likely Fix (Based on Your Deployment)

The deployment logs show the base schema encountered errors. The quickest fix is:

```bash
cd /opt/servqr
git pull origin main  # Get the latest fixes
sudo bash deployment/reset-and-rebuild-database.sh  # Reset and rebuild database
```

This will properly create all tables and restart services.
