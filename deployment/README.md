# ServQR Production Deployment Guide

Complete production deployment setup for ServQR platform on a plain Linux VM.

## 📋 Prerequisites

- **Linux VM** (Ubuntu 20.04+ or RHEL 8+)
- **Root/Sudo Access**
- **Source Code** downloaded to `/opt/servqr`
- **Domain Name** (optional, can use IP)
- **Minimum Resources:**
  - 4 GB RAM
  - 2 CPU cores
  - 50 GB disk space
  - Open ports: 80, 443, 8081, 3000

## 🚀 Quick Start (One Command Deployment)

```bash
# 1. Download source code
cd /opt
sudo git clone <your-repo-url> servqr
cd /opt/servqr/deployment

# 2. Run deployment (installs everything)
sudo bash deploy-all.sh
```

**That's it!** The script will:
- Install all system dependencies
- Setup Docker and PostgreSQL
- Build and start backend + frontend
- Configure systemd services for auto-start
- Setup SSL (if domain provided)

## 📂 Deployment Scripts

### 1. `deploy-all.sh` - Master Deployment Script
One-command deployment that orchestrates everything.

```bash
sudo bash deploy-all.sh
```

**What it does:**
- Runs prerequisites installation
- Sets up Docker and database
- Builds and deploys application
- Configures services

### 2. `install-prerequisites.sh` - System Dependencies
Installs required system packages (Docker excluded).

```bash
sudo bash install-prerequisites.sh
```

**Installs:**
- Go 1.23+
- Node.js 20+
- Git
- Nginx
- Certbot (SSL)
- System utilities

### 3. `setup-docker.sh` - Docker & Database
Sets up Docker and PostgreSQL container.

```bash
sudo bash setup-docker.sh
```

**What it does:**
- Installs Docker + Docker Compose
- Creates PostgreSQL container
- Mounts data to `/opt/servqr/data/postgres`
- Creates database and user
- Runs migrations

### 4. `deploy-app.sh` - Application Deployment
Builds and starts the application.

```bash
sudo bash deploy-app.sh
```

**What it does:**
- Builds Go backend binary
- Builds Next.js frontend
- Configures environment variables
- Starts services via systemd

## 📁 Directory Structure

```
/opt/servqr/
├── deployment/               # Deployment scripts (this folder)
│   ├── README.md            # This file
│   ├── deploy-all.sh        # Master deployment script
│   ├── install-prerequisites.sh
│   ├── setup-docker.sh
│   ├── deploy-app.sh
│   ├── docker-compose.yml   # PostgreSQL container
│   ├── .env.production      # Environment template
│   ├── nginx.conf.template  # Nginx config
│   └── systemd/             # Service files
│       ├── servqr-backend.service
│       └── servqr-frontend.service
├── data/                     # Runtime data (created by scripts)
│   ├── postgres/            # PostgreSQL data (persistent)
│   ├── qrcodes/             # Generated QR codes
│   └── whatsapp/            # WhatsApp media
├── logs/                     # Application logs
│   ├── backend.log
│   └── frontend.log
├── backups/                  # Database backups
└── [source code...]          # Your application code
```

## 🔧 Configuration

### Environment Variables

Edit `/opt/servqr/deployment/.env.production`:

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=servqr
DATABASE_PASSWORD=<generate-secure-password>
DATABASE_NAME=servqr_production

# Application
PORT=8081
FRONTEND_PORT=3000
BASE_URL=https://yourdomain.com
ENVIRONMENT=production

# JWT
JWT_SECRET=<generate-secure-secret>

# AI (Optional)
OPENAI_API_KEY=your-key
ANTHROPIC_API_KEY=your-key

# Email (Optional)
SENDGRID_API_KEY=your-key
SENDGRID_FROM_EMAIL=noreply@yourdomain.com

# WhatsApp (Optional)
TWILIO_ACCOUNT_SID=your-sid
TWILIO_AUTH_TOKEN=your-token
```

**Generate Secure Secrets:**
```bash
# JWT Secret (32 characters)
openssl rand -base64 32

# Database Password (24 characters)
openssl rand -base64 24
```

### Domain Setup (Optional)

If you have a domain:

```bash
# Edit deploy-all.sh before running
DOMAIN="yourdomain.com"
EMAIL="admin@yourdomain.com"
```

The script will:
- Configure Nginx reverse proxy
- Obtain SSL certificate via Let's Encrypt
- Enable HTTPS

## 🐳 Docker Configuration

PostgreSQL runs in Docker with data persistence:

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:15-alpine
    volumes:
      - /opt/servqr/data/postgres:/var/lib/postgresql/data
    ports:
      - "5432:5432"
```

**Database Access:**
```bash
# Connect to PostgreSQL
docker exec -it servqr-postgres psql -U servqr -d servqr_production

# View logs
docker logs servqr-postgres

# Restart database
docker restart servqr-postgres
```

## 🔄 Service Management

Services run via systemd and auto-start on boot:

```bash
# Backend
sudo systemctl status servqr-backend
sudo systemctl start servqr-backend
sudo systemctl stop servqr-backend
sudo systemctl restart servqr-backend

# Frontend
sudo systemctl status servqr-frontend
sudo systemctl start servqr-frontend
sudo systemctl stop servqr-frontend
sudo systemctl restart servqr-frontend

# View logs
sudo journalctl -u servqr-backend -f
sudo journalctl -u servqr-frontend -f
```

## 🔍 Health Checks

```bash
# Backend health
curl http://localhost:8081/health

# Frontend health
curl http://localhost:3000

# Database health
docker exec servqr-postgres pg_isready -U servqr
```

## 💾 Database Management

### Migrations

```bash
cd /opt/servqr
./deployment/run-migrations.sh
```

### Backup

```bash
# Manual backup
./deployment/backup-database.sh

# Restore from backup
./deployment/restore-database.sh /opt/servqr/backups/backup-2026-02-06.sql
```

### Automated Backups

Configured via cron (runs daily at 2 AM):

```bash
# View cron jobs
sudo crontab -l

# Edit cron jobs
sudo crontab -e
```

## 🔒 Security Checklist

- [ ] Changed default database password
- [ ] Generated secure JWT secret
- [ ] Firewall configured (only 80, 443 open)
- [ ] SSL certificate installed
- [ ] Database not exposed externally
- [ ] Regular backups enabled
- [ ] Monitoring configured
- [ ] Log rotation enabled

## 🚨 Troubleshooting

### Backend Not Starting

```bash
# Check logs
sudo journalctl -u servqr-backend -n 50

# Check if port is in use
sudo netstat -tulpn | grep 8081

# Verify binary exists
ls -la /opt/servqr/platform

# Test manual start
cd /opt/servqr
./platform
```

### Frontend Not Starting

```bash
# Check logs
sudo journalctl -u servqr-frontend -n 50

# Verify build exists
ls -la /opt/servqr/admin-ui/.next

# Test manual start
cd /opt/servqr/admin-ui
npm start
```

### Database Connection Issues

```bash
# Check if container is running
docker ps | grep servqr-postgres

# Check database logs
docker logs servqr-postgres

# Test connection
psql -h localhost -U servqr -d servqr_production

# Restart database
docker restart servqr-postgres
```

### Port Already in Use

```bash
# Find process using port
sudo lsof -i :8081
sudo lsof -i :3000

# Kill process
sudo kill -9 <PID>
```

## 📊 Monitoring

### Application Logs

```bash
# View logs
tail -f /opt/servqr/logs/backend.log
tail -f /opt/servqr/logs/frontend.log

# View systemd logs
sudo journalctl -u servqr-backend -f
sudo journalctl -u servqr-frontend -f
```

### Resource Usage

```bash
# CPU and Memory
htop

# Disk usage
df -h
du -sh /opt/servqr/*

# Docker stats
docker stats servqr-postgres
```

### Database Monitoring

```bash
# Connect to database
docker exec -it servqr-postgres psql -U servqr -d servqr_production

# Check database size
SELECT pg_size_pretty(pg_database_size('servqr_production'));

# Check table sizes
SELECT schemaname, tablename, 
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables 
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

# Active connections
SELECT count(*) FROM pg_stat_activity;
```

## 🔄 Updates and Maintenance

### Update Application

```bash
cd /opt/servqr

# Pull latest code
git pull origin main

# Redeploy
sudo bash deployment/deploy-app.sh

# Restart services
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

### Database Maintenance

```bash
# Vacuum database (reclaim space)
docker exec servqr-postgres psql -U servqr -d servqr_production -c "VACUUM FULL;"

# Analyze tables (update statistics)
docker exec servqr-postgres psql -U servqr -d servqr_production -c "ANALYZE;"

# Reindex (improve performance)
docker exec servqr-postgres psql -U servqr -d servqr_production -c "REINDEX DATABASE servqr_production;"
```

## 🎯 Performance Tuning

### Backend

Edit `/opt/servqr/.env`:

```bash
# Increase connection pool
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Adjust timeouts
REQUEST_TIMEOUT=60
DB_TIMEOUT=30
```

### Frontend

Edit `/opt/servqr/admin-ui/.env.local`:

```bash
# Enable production optimizations
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1
```

### PostgreSQL

Edit `docker-compose.yml`:

```yaml
services:
  postgres:
    command:
      - "postgres"
      - "-c"
      - "max_connections=100"
      - "-c"
      - "shared_buffers=256MB"
      - "-c"
      - "effective_cache_size=1GB"
```

## 📞 Support

For issues or questions:

1. Check logs: `/opt/servqr/logs/`
2. Review documentation: `/opt/servqr/docs/`
3. Check systemd status: `sudo systemctl status servqr-*`
4. Contact support: Include logs and error messages

## ✅ Post-Deployment Checklist

- [ ] All services running: `sudo systemctl status servqr-*`
- [ ] Health checks passing: `curl http://localhost:8081/health`
- [ ] Database accessible: `docker exec -it servqr-postgres psql -U servqr`
- [ ] Frontend accessible: Visit `http://your-server-ip:3000`
- [ ] Backend API accessible: `curl http://your-server-ip:8081/api/v1/health`
- [ ] SSL configured (if domain): Visit `https://yourdomain.com`
- [ ] Backups working: Check `/opt/servqr/backups/`
- [ ] Monitoring configured: Check logs
- [ ] Firewall configured: `sudo ufw status`
- [ ] Services auto-start on boot: `sudo systemctl is-enabled servqr-*`

---

**Version:** 1.0.0  
**Last Updated:** 2026-02-06  
**Deployment Target:** Production Linux VM
