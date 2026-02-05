# ServQR Production Deployment Checklist

Complete checklist for deploying ServQR platform to production.

## 📋 Pre-Deployment

### Server Requirements
- [ ] Linux VM provisioned (Ubuntu 20.04+ or RHEL 8+)
- [ ] Minimum 4 GB RAM, 2 CPU cores
- [ ] 50 GB disk space available
- [ ] Root/sudo access configured
- [ ] Internet connectivity verified
- [ ] Firewall configured (ports 80, 443 open)

### Domain & DNS (Optional but Recommended)
- [ ] Domain name registered
- [ ] DNS A record pointing to server IP
- [ ] DNS propagation verified (`nslookup yourdomain.com`)
- [ ] Email address for SSL certificate

### Credentials & API Keys
- [ ] OpenAI API key (for AI diagnosis)
- [ ] Anthropic API key (optional fallback)
- [ ] SendGrid API key (for email notifications)
- [ ] Twilio credentials (for WhatsApp integration)
- [ ] Admin email addresses collected

### Source Code
- [ ] Repository access configured
- [ ] Latest code pulled/downloaded
- [ ] Code placed at `/opt/servqr`
- [ ] Migrations directory verified

---

## 🚀 Deployment Steps

### 1. Initial Setup
```bash
cd /opt/servqr/deployment
```

- [ ] Reviewed `README.md`
- [ ] Reviewed `QUICKSTART.md`
- [ ] Updated `deploy-all.sh` with domain (if applicable)

### 2. Run Deployment Script
```bash
sudo bash deploy-all.sh
```

**Monitor for:**
- [ ] Prerequisites installation completed
- [ ] Docker installed successfully
- [ ] PostgreSQL container started
- [ ] Database migrations applied
- [ ] Backend built successfully
- [ ] Frontend built successfully
- [ ] Services started

**Expected Duration:** 5-10 minutes

### 3. Verify Deployment
- [ ] Backend health check: `curl http://localhost:8081/health`
- [ ] Frontend accessible: `curl http://localhost:3000`
- [ ] Database connection: `docker exec servqr-postgres pg_isready -U servqr`
- [ ] Services running: `sudo systemctl status servqr-*`
- [ ] No errors in logs: `sudo journalctl -u servqr-* -n 50`

---

## ⚙️ Post-Deployment Configuration

### 1. Update Environment Variables
```bash
sudo vim /opt/servqr/.env
```

**Required Updates:**
- [ ] `AI_OPENAI_API_KEY` - Add your OpenAI key
- [ ] `SENDGRID_API_KEY` - Add your SendGrid key
- [ ] `SENDGRID_FROM_EMAIL` - Update sender email
- [ ] Review and update `BASE_URL` (if domain configured)

**Optional Updates:**
- [ ] `AI_ANTHROPIC_API_KEY` - For AI fallback
- [ ] `TWILIO_*` - For WhatsApp integration
- [ ] Feature flags as needed

**After Updates:**
```bash
sudo systemctl restart servqr-backend
```

### 2. Update Frontend Configuration (if needed)
```bash
sudo vim /opt/servqr/admin-ui/.env.local
```

- [ ] Verify `NEXT_PUBLIC_API_URL` matches your domain
- [ ] Restart frontend: `sudo systemctl restart servqr-frontend`

### 3. Create Admin User
- [ ] Access frontend: `http://your-server-ip:3000` or `https://yourdomain.com`
- [ ] Use default credentials from `/opt/servqr/LOGIN-CREDENTIALS.txt`
- [ ] Change default password immediately
- [ ] Create additional admin users as needed

### 4. Database Verification
```bash
/opt/servqr/deployment/connect-database.sh
```

**Verify:**
- [ ] Tables created: `\dt`
- [ ] No errors in schema
- [ ] Sample data (if any) loaded
- [ ] Exit: `\q`

---

## 🔒 Security Hardening

### 1. Change Default Credentials
- [ ] Database password changed
- [ ] JWT secret regenerated
- [ ] Admin password changed
- [ ] Root password changed (if default)

### 2. Firewall Configuration
```bash
# Ubuntu/Debian
sudo ufw status
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# CentOS/RHEL
sudo firewall-cmd --list-all
```

- [ ] Firewall enabled
- [ ] Only necessary ports open (22, 80, 443)
- [ ] Ports 8081, 3000 not exposed externally (proxied via Nginx)
- [ ] Database port 5432 not exposed externally

### 3. SSL Certificate (if domain)
```bash
sudo certbot certificates
```

- [ ] SSL certificate installed
- [ ] HTTPS working: `https://yourdomain.com`
- [ ] HTTP redirects to HTTPS
- [ ] Certificate auto-renewal configured

### 4. System Security
- [ ] System packages updated: `sudo apt update && sudo apt upgrade -y`
- [ ] SSH key-based authentication configured
- [ ] Password authentication disabled (optional)
- [ ] Fail2ban installed (optional): `sudo apt install fail2ban`
- [ ] Automatic security updates enabled

---

## 💾 Backup & Monitoring

### 1. Automated Backups
```bash
sudo crontab -l | grep backup
```

- [ ] Cron job for daily backups verified (2 AM)
- [ ] Backup script tested: `sudo /opt/servqr/deployment/backup-database.sh`
- [ ] Backup location verified: `/opt/servqr/backups/`
- [ ] Backup retention policy configured (7 days)

### 2. Manual Backup Test
```bash
sudo /opt/servqr/deployment/backup-database.sh
ls -lh /opt/servqr/backups/
```

- [ ] Backup created successfully
- [ ] Backup file size reasonable (not 0 bytes)
- [ ] Backup compressed (.gz)

### 3. Restore Test (Optional but Recommended)
```bash
# On test environment only!
sudo /opt/servqr/deployment/restore-database.sh /opt/servqr/backups/latest-backup.sql.gz
```

- [ ] Restore procedure tested (on dev/staging)
- [ ] Restore successful
- [ ] Application works after restore

### 4. Monitoring Setup
- [ ] Log rotation configured: `cat /etc/logrotate.d/servqr`
- [ ] Disk space monitoring setup
- [ ] Service health monitoring setup
- [ ] Alert notifications configured (optional)

---

## 🔍 Testing & Validation

### 1. Application Testing
- [ ] Frontend loads without errors
- [ ] Login works with admin credentials
- [ ] Dashboard displays data
- [ ] Navigation works (all pages load)
- [ ] No JavaScript console errors

### 2. API Testing
```bash
# Health check
curl http://localhost:8081/health

# Organizations API
curl http://localhost:8081/api/v1/organizations

# Equipment API
curl http://localhost:8081/api/v1/equipment
```

- [ ] Health endpoint returns `{"status":"ok"}`
- [ ] API endpoints respond (may require authentication)
- [ ] No 500 errors

### 3. Database Testing
```bash
docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt"
```

- [ ] All tables present
- [ ] No missing migrations
- [ ] Indexes created

### 4. Integration Testing
- [ ] Create test organization
- [ ] Add test equipment
- [ ] Create test service ticket
- [ ] Assign test engineer
- [ ] Email notification sent (if configured)

### 5. Performance Testing
- [ ] Page load times acceptable (<3 seconds)
- [ ] API response times good (<500ms)
- [ ] Database queries optimized
- [ ] No memory leaks: `htop`

---

## 📊 Production Readiness

### 1. Documentation
- [ ] Deployment documentation reviewed
- [ ] API documentation accessible: `/opt/servqr/docs/`
- [ ] Troubleshooting guide reviewed
- [ ] Runbooks created for common issues

### 2. Access & Credentials
- [ ] Admin credentials documented (securely)
- [ ] Database credentials stored in password manager
- [ ] API keys stored securely
- [ ] Team access configured

### 3. Monitoring Dashboard (Optional)
- [ ] Server metrics monitored (CPU, RAM, Disk)
- [ ] Application logs monitored
- [ ] Database performance monitored
- [ ] Alert thresholds configured

### 4. Disaster Recovery
- [ ] Backup strategy documented
- [ ] Recovery procedures tested
- [ ] RTO/RPO defined
- [ ] Contact list for emergencies

---

## ✅ Final Verification

### Run Full Health Check
```bash
cd /opt/servqr/deployment

# 1. Services
sudo systemctl is-active servqr-backend
sudo systemctl is-active servqr-frontend
docker ps | grep servqr-postgres

# 2. Health endpoints
curl http://localhost:8081/health
curl http://localhost:3000

# 3. Database
docker exec servqr-postgres pg_isready -U servqr

# 4. Logs (no critical errors)
sudo journalctl -u servqr-backend -n 50 | grep -i error
sudo journalctl -u servqr-frontend -n 50 | grep -i error

# 5. Disk space
df -h /opt/servqr

# 6. Backups
ls -lh /opt/servqr/backups/
```

**All Checks Passed:**
- [ ] Backend service: ✓ Active
- [ ] Frontend service: ✓ Active
- [ ] Database: ✓ Running
- [ ] Health checks: ✓ Passing
- [ ] Logs: ✓ No errors
- [ ] Disk space: ✓ Adequate
- [ ] Backups: ✓ Configured

---

## 🎉 Go Live

### Pre-Launch
- [ ] All checklist items completed
- [ ] Team briefed on deployment
- [ ] Support contacts ready
- [ ] Rollback plan prepared

### Launch
- [ ] DNS updated (if switching domains)
- [ ] Application accessible to users
- [ ] Monitoring active
- [ ] First users onboarded

### Post-Launch
- [ ] Monitor for first 24-48 hours
- [ ] Check for errors/issues
- [ ] Collect user feedback
- [ ] Plan for improvements

---

## 📞 Support Contacts

**Emergency Contacts:**
- System Admin: _________________
- Database Admin: _________________
- DevOps: _________________
- On-call: _________________

**External Services:**
- Hosting Provider: _________________
- Domain Registrar: _________________
- SSL Provider: Let's Encrypt (auto-renew)

---

## 📝 Notes

**Deployment Date:** _________________  
**Deployed By:** _________________  
**Server IP:** _________________  
**Domain:** _________________  
**Version:** 1.0.0

**Issues Encountered:**
- 
- 

**Resolutions:**
- 
- 

---

**Status:** ☐ Not Started | ☐ In Progress | ☐ Completed | ☐ Verified

**Sign-off:**
- Technical Lead: _________________ Date: _________
- Operations: _________________ Date: _________
- Management: _________________ Date: _________
