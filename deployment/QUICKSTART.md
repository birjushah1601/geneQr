# ServQR Production Deployment - Quick Start Guide

## 🚀 One-Command Deployment

Complete deployment from a plain Linux VM to a running ServQR platform in **under 10 minutes**.

### Prerequisites

✅ Linux VM (Ubuntu 20.04+ or RHEL 8+)  
✅ Root/sudo access  
✅ 4 GB RAM, 2 CPU cores, 50 GB disk  
✅ Internet connectivity

### Step 1: Download Source Code

```bash
# Clone repository to /opt/servqr
cd /opt
sudo git clone <your-repo-url> servqr

# Or if you already have the code elsewhere, move it
sudo mv /path/to/your/code /opt/servqr
```

### Step 2: Run Deployment

```bash
cd /opt/servqr/deployment
sudo bash deploy-all.sh
```

**That's it!** Wait 5-10 minutes for the deployment to complete.

### Step 3: Access Application

```bash
# Get your server IP
hostname -I
```

Visit: `http://YOUR_SERVER_IP:3000`

---

## 📋 What Gets Deployed

The `deploy-all.sh` script automatically:

1. **Installs System Dependencies**
   - Go 1.23+
   - Node.js 20+
   - Nginx
   - Git, curl, build tools
   - Certbot (for SSL)

2. **Sets Up Docker & Database**
   - Docker + Docker Compose
   - PostgreSQL 15 container
   - Persistent data storage at `/opt/servqr/data/postgres`
   - Database migrations applied

3. **Builds & Deploys Application**
   - Compiles Go backend binary
   - Builds Next.js frontend
   - Configures environment variables
   - Generates secure secrets (JWT, DB password)

4. **Configures Services**
   - Systemd services (auto-start on boot)
   - Automated database backups (daily at 2 AM)
   - Log rotation
   - Firewall rules

---

## 🔐 Default Credentials

**Database:**
- Host: localhost:5432
- Database: servqr_production
- User: servqr
- Password: Check `/opt/servqr/.db_password`

**Application:**
- Check `/opt/servqr/LOGIN-CREDENTIALS.txt` for default admin credentials

---

## 🎯 Service Management

### View Status
```bash
sudo systemctl status servqr-backend
sudo systemctl status servqr-frontend
```

### View Logs
```bash
# Real-time logs
sudo journalctl -u servqr-backend -f
sudo journalctl -u servqr-frontend -f

# Last 50 lines
sudo journalctl -u servqr-backend -n 50
```

### Restart Services
```bash
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

### Stop/Start Services
```bash
sudo systemctl stop servqr-backend
sudo systemctl start servqr-backend
```

---

## 🗄️ Database Management

### Connect to Database
```bash
/opt/servqr/deployment/connect-database.sh
```

### Backup Database
```bash
/opt/servqr/deployment/backup-database.sh
```

### Restore Database
```bash
/opt/servqr/deployment/restore-database.sh /opt/servqr/backups/backup-file.sql.gz
```

### View Database Logs
```bash
docker logs servqr-postgres
```

---

## ⚙️ Configuration

### Backend Configuration
```bash
sudo vim /opt/servqr/.env
```

**Critical Settings:**
- `DATABASE_PASSWORD` - Auto-generated, stored in `.db_password`
- `JWT_SECRET` - Auto-generated, stored in `.jwt_secret`
- `OPENAI_API_KEY` - Add your OpenAI key for AI features
- `SENDGRID_API_KEY` - Add your SendGrid key for emails
- `TWILIO_*` - Add Twilio credentials for WhatsApp

### Frontend Configuration
```bash
sudo vim /opt/servqr/admin-ui/.env.local
```

**After Changing Config:**
```bash
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

---

## 🔍 Health Checks

```bash
# Backend API
curl http://localhost:8081/health

# Frontend
curl http://localhost:3000

# Database
docker exec servqr-postgres pg_isready -U servqr
```

**Expected Response:**
- Backend: `{"status":"ok"}`
- Frontend: HTTP 200 (HTML response)
- Database: `servqr-postgres:5432 - accepting connections`

---

## 🌐 Domain & SSL Setup (Optional)

If you have a domain name:

### Before Deployment
Edit `deploy-all.sh` and set:
```bash
DOMAIN="yourdomain.com"
EMAIL="admin@yourdomain.com"
```

Then run deployment.

### After Deployment
The script will automatically:
- Configure Nginx reverse proxy
- Obtain SSL certificate from Let's Encrypt
- Enable HTTPS

**Access:** `https://yourdomain.com`

### Manual SSL Setup
```bash
sudo certbot --nginx -d yourdomain.com
```

---

## 🚨 Troubleshooting

### Backend Not Starting

```bash
# Check logs
sudo journalctl -u servqr-backend -n 100

# Check if port is in use
sudo netstat -tulpn | grep 8081

# Test manual start
cd /opt/servqr
./platform
```

**Common Issues:**
- Database not ready → Wait 30 seconds and retry
- Port in use → `sudo lsof -i :8081` and kill process
- Missing .env → Run `deploy-app.sh` again

### Frontend Not Starting

```bash
# Check logs
sudo journalctl -u servqr-frontend -n 100

# Verify build
ls -la /opt/servqr/admin-ui/.next

# Test manual start
cd /opt/servqr/admin-ui
npm start
```

**Common Issues:**
- Build failed → Check Node.js version (`node -v` should be 20+)
- Port in use → `sudo lsof -i :3000` and kill process

### Database Connection Failed

```bash
# Check container status
docker ps | grep servqr-postgres

# Check database logs
docker logs servqr-postgres

# Restart container
docker restart servqr-postgres

# Test connection
psql -h localhost -U servqr -d servqr_production
# Password: cat /opt/servqr/.db_password
```

### Services Not Auto-Starting on Boot

```bash
# Enable services
sudo systemctl enable servqr-postgres
sudo systemctl enable servqr-backend
sudo systemctl enable servqr-frontend

# Verify
sudo systemctl is-enabled servqr-backend
```

---

## 📊 Monitoring

### Resource Usage
```bash
# CPU and Memory
htop

# Disk space
df -h

# Docker stats
docker stats servqr-postgres
```

### Application Logs
```bash
# Backend logs
tail -f /opt/servqr/logs/backend.log
tail -f /opt/servqr/logs/backend-error.log

# Frontend logs
tail -f /opt/servqr/logs/frontend.log
```

### Database Stats
```bash
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
SELECT 
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 10;
"
```

---

## 🔄 Updates

### Update Application Code

```bash
cd /opt/servqr

# Pull latest code
git pull origin main

# Rebuild and redeploy
sudo bash deployment/deploy-app.sh

# Restart services
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

### Update System Packages

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# CentOS/RHEL
sudo yum update -y

# Restart if kernel updated
sudo reboot
```

---

## 🔒 Security Checklist

After deployment, ensure you:

- [ ] Changed database password: Edit `/opt/servqr/.env`
- [ ] Changed JWT secret: Edit `/opt/servqr/.env`
- [ ] Updated admin login credentials
- [ ] Configured firewall: Only ports 80, 443 open externally
- [ ] Installed SSL certificate (if domain)
- [ ] Set up regular backups: Check cron with `sudo crontab -l`
- [ ] Configured monitoring/alerts
- [ ] Reviewed logs for errors

---

## 📞 Support

### Useful Commands

```bash
# View all services
sudo systemctl list-units servqr-*

# View all logs
sudo journalctl -u servqr-* -f

# Restart everything
sudo systemctl restart servqr-postgres servqr-backend servqr-frontend

# Stop everything
sudo systemctl stop servqr-backend servqr-frontend
docker stop servqr-postgres
```

### Log Locations

- Deployment: `/opt/servqr/logs/deployment-*.log`
- Backend: `/opt/servqr/logs/backend*.log`
- Frontend: `/opt/servqr/logs/frontend*.log`
- Database: `docker logs servqr-postgres`
- Systemd: `journalctl -u servqr-*`

### Getting Help

1. Check logs first
2. Review troubleshooting section above
3. Check documentation: `/opt/servqr/docs/`
4. Check GitHub issues: [repository-url]/issues

---

## ✅ Deployment Verification

Run this checklist to verify successful deployment:

```bash
# 1. Services running
sudo systemctl is-active servqr-backend   # Should return: active
sudo systemctl is-active servqr-frontend  # Should return: active
docker ps | grep servqr-postgres          # Should show running container

# 2. Health checks
curl http://localhost:8081/health         # Should return: {"status":"ok"}
curl http://localhost:3000                # Should return: HTML page

# 3. Database
docker exec servqr-postgres pg_isready -U servqr  # Should return: accepting connections

# 4. Logs (should have no critical errors)
sudo journalctl -u servqr-backend -n 20 | grep -i error
sudo journalctl -u servqr-frontend -n 20 | grep -i error

# 5. Firewall
sudo ufw status  # Ubuntu
sudo firewall-cmd --list-all  # RHEL/CentOS

# 6. Backups configured
sudo crontab -l | grep backup

# 7. Files exist
ls -la /opt/servqr/platform              # Backend binary
ls -la /opt/servqr/admin-ui/.next        # Frontend build
ls -la /opt/servqr/.env                  # Backend config
```

**If all checks pass, your deployment is successful!** 🎉

---

**Questions?** Check the full documentation: `/opt/servqr/deployment/README.md`
