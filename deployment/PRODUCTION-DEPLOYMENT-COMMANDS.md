# Production Deployment Commands - servqr.com

Complete step-by-step commands to deploy ServQR platform on a production Linux VM.

## Pre-Deployment Checklist

Before running any commands:

- [ ] Fresh Ubuntu 20.04+ or RHEL 8+ VM
- [ ] Root/sudo access
- [ ] DNS configured: `servqr.com` → Your server IP
- [ ] Ports open: 22 (SSH), 80 (HTTP), 443 (HTTPS)
- [ ] Valid email for SSL: `admin@servqr.com`

---

## Step 1: Initial Server Setup

**Connect to your server:**

```bash
# SSH into your server
ssh root@YOUR_SERVER_IP
# or
ssh ubuntu@YOUR_SERVER_IP
```

**Update system:**

```bash
sudo apt-get update
sudo apt-get upgrade -y
```

**Set hostname (optional but recommended):**

```bash
sudo hostnamectl set-hostname servqr
```

---

## Step 2: Clone Repository

**Create installation directory:**

```bash
sudo mkdir -p /opt/servqr
sudo chown -R $USER:$USER /opt/servqr
cd /opt/servqr
```

**Clone the repository:**

```bash
# If you have git credentials
git clone https://github.com/birjushah1601/geneQr.git .

# Or download and extract
# wget https://github.com/birjushah1601/geneQr/archive/refs/heads/main.zip
# unzip main.zip
# mv geneQr-main/* .
```

**Verify files:**

```bash
ls -la
# Should see: deployment/, admin-ui/, cmd/, internal/, etc.
```

---

## Step 3: Configure Domain and Email

**Edit deployment script:**

```bash
sudo nano /opt/servqr/deployment/deploy-all.sh
```

**Find and update these lines (around line 37-38):**

```bash
# Change from:
DOMAIN=""
EMAIL=""

# To:
DOMAIN="servqr.com"
EMAIL="admin@servqr.com"
```

**Save and exit:** `Ctrl+X`, then `Y`, then `Enter`

---

## Step 4: Run Full Deployment

**Execute deployment script:**

```bash
cd /opt/servqr/deployment
sudo bash deploy-all.sh
```

**This will automatically:**
1. ✅ Install prerequisites (Go, Node.js, Nginx)
2. ✅ Setup Docker and PostgreSQL
3. ✅ Initialize database with working schema
4. ✅ Build backend (Go)
5. ✅ Build frontend (Next.js)
6. ✅ Create systemd services
7. ✅ Configure Nginx with SSL
8. ✅ Obtain Let's Encrypt certificate
9. ✅ Setup auto-renewal
10. ✅ Configure backups

**Expected duration:** 10-15 minutes

**Watch for:**
- Green checkmarks ✓
- "Deployment completed successfully"
- No red errors

---

## Step 5: Verify Deployment

**Check services status:**

```bash
# Backend status
sudo systemctl status servqr-backend

# Frontend status
sudo systemctl status servqr-frontend

# Should both show: "active (running)"
```

**Check database:**

```bash
# Connect to database
docker exec -it servqr-postgres psql -U servqr -d servqr_production

# List tables
\dt

# Should show 65 tables

# Exit
\q
```

**Check Nginx:**

```bash
# Nginx status
sudo systemctl status nginx

# Test SSL certificate
sudo certbot certificates

# Should show:
# Certificate Name: servqr.com
# Domains: servqr.com
# Expiry Date: (90 days from now)
```

**Check logs:**

```bash
# Backend logs
sudo journalctl -u servqr-backend -n 50

# Frontend logs
sudo journalctl -u servqr-frontend -n 50

# Nginx logs
sudo tail -f /var/log/nginx/servqr-access.log
```

---

## Step 6: Test Application

**Test from server:**

```bash
# Test backend health
curl http://localhost:8081/health

# Test frontend
curl http://localhost:3000

# Test HTTPS
curl -I https://servqr.com

# Should return: HTTP/2 200
```

**Test in browser:**

1. **Visit:** `http://servqr.com`
   - Should redirect to `https://servqr.com`

2. **Visit:** `https://servqr.com`
   - Should show ServQR login page
   - Green padlock 🔒 in browser
   - No certificate warnings

3. **Check certificate:**
   - Click padlock → Certificate
   - Issued by: Let's Encrypt
   - Valid for: 90 days

4. **Test login:**
   - Check `/opt/servqr/LOGIN-CREDENTIALS.txt` for credentials
   - Try logging in

---

## Step 7: Configure Firewall (Important!)

**Setup UFW firewall:**

```bash
# Install UFW
sudo apt-get install -y ufw

# Allow SSH (IMPORTANT - don't lock yourself out!)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

**Expected output:**
```
Status: active

To                         Action      From
--                         ------      ----
22/tcp                     ALLOW       Anywhere
80/tcp                     ALLOW       Anywhere
443/tcp                    ALLOW       Anywhere
```

---

## Step 8: Update Application Configuration

**Edit environment variables:**

```bash
sudo nano /opt/servqr/.env
```

**Update these values:**

```bash
# API URLs (if needed)
NEXT_PUBLIC_API_BASE_URL=https://servqr.com/api
APP_URL=https://servqr.com

# Email Configuration (if using SendGrid)
SENDGRID_API_KEY=your_sendgrid_api_key
SENDGRID_FROM_EMAIL=noreply@servqr.com
SENDGRID_FROM_NAME=ServQR

# WhatsApp Configuration (if using)
WHATSAPP_API_TOKEN=your_whatsapp_token
WHATSAPP_PHONE_NUMBER=your_phone_number

# AI Configuration (if using)
OPENAI_API_KEY=your_openai_key
```

**Save and exit**

**Restart services to apply changes:**

```bash
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend

# Verify both restarted successfully
sudo systemctl status servqr-backend
sudo systemctl status servqr-frontend
```

---

## Step 9: Setup Monitoring (Optional but Recommended)

**Setup log rotation (already configured):**

```bash
# Check log rotation config
cat /etc/logrotate.d/servqr
```

**Setup automated backups (already configured):**

```bash
# Check backup cron
sudo crontab -l

# Should show: 0 2 * * * /opt/servqr/deployment/backup-database.sh

# Test backup manually
sudo /opt/servqr/deployment/backup-database.sh

# Check backup created
ls -lh /opt/servqr/backups/
```

**Setup monitoring (optional):**

```bash
# Install monitoring tools
sudo apt-get install -y htop iotop nethogs

# Monitor resources
htop

# Monitor network
sudo nethogs

# Monitor disk
df -h
```

---

## Management Commands

### Service Management

```bash
# Restart services
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
sudo systemctl restart nginx

# Stop services
sudo systemctl stop servqr-backend
sudo systemctl stop servqr-frontend

# Start services
sudo systemctl start servqr-backend
sudo systemctl start servqr-frontend

# Check service status
sudo systemctl status servqr-backend
sudo systemctl status servqr-frontend

# View service logs (live)
sudo journalctl -u servqr-backend -f
sudo journalctl -u servqr-frontend -f
```

### Database Management

```bash
# Connect to database
docker exec -it servqr-postgres psql -U servqr -d servqr_production

# Backup database
sudo /opt/servqr/deployment/backup-database.sh

# List backups
ls -lh /opt/servqr/backups/

# Restore database (if needed)
sudo bash /opt/servqr/deployment/restore-database.sh

# Check database size
docker exec servqr-postgres psql -U servqr -d servqr_production -c "
  SELECT pg_size_pretty(pg_database_size('servqr_production'));
"
```

### SSL Certificate Management

```bash
# Check certificate
sudo certbot certificates

# Test renewal
sudo certbot renew --dry-run

# Force renewal (if needed)
sudo certbot renew --force-renewal

# Revoke certificate (if needed)
sudo certbot revoke --cert-path /etc/letsencrypt/live/servqr.com/fullchain.pem
```

### Log Management

```bash
# View application logs
sudo tail -f /opt/servqr/logs/backend.log
sudo tail -f /opt/servqr/logs/frontend.log

# View Nginx logs
sudo tail -f /var/log/nginx/servqr-access.log
sudo tail -f /var/log/nginx/servqr-error.log

# View system logs
sudo journalctl -xe

# Clear old logs (if disk space low)
sudo journalctl --vacuum-time=7d
```

### Update Application

```bash
# Pull latest code
cd /opt/servqr
git pull origin main

# Rebuild backend
cd /opt/servqr
go build -o platform.exe ./cmd/platform

# Rebuild frontend
cd /opt/servqr/admin-ui
npm install
npm run build

# Restart services
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

---

## Troubleshooting

### Issue: Backend not starting

```bash
# Check logs
sudo journalctl -u servqr-backend -n 100

# Check if port 8081 is in use
sudo netstat -tulpn | grep 8081

# Check database connection
docker exec servqr-postgres psql -U servqr -d servqr_production -c "SELECT 1;"

# Check environment file
cat /opt/servqr/.env | grep DB_
```

### Issue: Frontend not starting

```bash
# Check logs
sudo journalctl -u servqr-frontend -n 100

# Check if port 3000 is in use
sudo netstat -tulpn | grep 3000

# Check Node.js version
node --version

# Rebuild frontend
cd /opt/servqr/admin-ui
npm install
npm run build
sudo systemctl restart servqr-frontend
```

### Issue: SSL certificate not obtained

```bash
# Check DNS
nslookup servqr.com

# Check port 80 is accessible
sudo netstat -tulpn | grep :80

# Check Nginx error log
sudo tail -f /var/log/nginx/error.log

# Check Let's Encrypt log
sudo tail -f /var/log/letsencrypt/letsencrypt.log

# Try manual certificate
sudo certbot --nginx -d servqr.com --email admin@servqr.com
```

### Issue: Database connection failed

```bash
# Check PostgreSQL container
docker ps | grep postgres

# Check PostgreSQL logs
docker logs servqr-postgres

# Restart PostgreSQL
docker restart servqr-postgres

# Check database exists
docker exec servqr-postgres psql -U servqr -l
```

### Issue: Out of disk space

```bash
# Check disk usage
df -h

# Find large files
sudo du -h /opt/servqr | sort -rh | head -20

# Clean old logs
sudo journalctl --vacuum-time=7d
sudo find /opt/servqr/logs -name "*.log" -mtime +7 -delete

# Clean old backups (keep last 7 days)
sudo find /opt/servqr/backups -name "*.sql.gz" -mtime +7 -delete

# Clean Docker
docker system prune -a
```

---

## Security Checklist

After deployment, ensure:

- [ ] Firewall enabled (UFW)
- [ ] Only necessary ports open (22, 80, 443)
- [ ] SSH key-based authentication (disable password)
- [ ] Regular backups configured
- [ ] SSL certificate obtained and auto-renewal working
- [ ] Default passwords changed
- [ ] JWT secret updated in .env
- [ ] Fail2Ban installed (optional)
- [ ] Regular system updates scheduled

**Change default passwords:**

```bash
# Generate new secure password
openssl rand -base64 32

# Update .env with new passwords
sudo nano /opt/servqr/.env

# Restart services
sudo systemctl restart servqr-backend
sudo systemctl restart servqr-frontend
```

---

## Production URLs

After successful deployment:

- **Application:** https://servqr.com
- **API:** https://servqr.com/api
- **Health Check:** https://servqr.com/api/health
- **Admin Panel:** https://servqr.com/login

---

## Quick Reference Card

**Most Common Commands:**

```bash
# Restart everything
sudo systemctl restart servqr-backend servqr-frontend nginx

# Check everything is running
sudo systemctl status servqr-backend servqr-frontend nginx

# View logs
sudo journalctl -u servqr-backend -f

# Backup database
sudo /opt/servqr/deployment/backup-database.sh

# Connect to database
docker exec -it servqr-postgres psql -U servqr -d servqr_production

# Update SSL certificate
sudo certbot renew

# Check disk space
df -h

# Check memory
free -h

# Check CPU
top
```

---

## Support

**Log Locations:**
- Application: `/opt/servqr/logs/`
- Nginx: `/var/log/nginx/`
- System: `journalctl -u servqr-backend`
- Deployment: `/opt/servqr/logs/deployment-*.log`

**Configuration Files:**
- Environment: `/opt/servqr/.env`
- Nginx: `/etc/nginx/sites-available/servqr`
- Services: `/etc/systemd/system/servqr-*.service`

**Important Directories:**
- Installation: `/opt/servqr/`
- Backups: `/opt/servqr/backups/`
- QR Codes: `/opt/servqr/data/qrcodes/`
- Database Data: `/opt/servqr/data/postgres/`

---

**Deployment Script:** `/opt/servqr/deployment/deploy-all.sh`
**Domain Guide:** `/opt/servqr/deployment/DOMAIN-SSL-SETUP.md`
**Last Updated:** February 9, 2026
