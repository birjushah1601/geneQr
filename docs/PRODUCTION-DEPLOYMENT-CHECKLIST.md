# Production Deployment Checklist

**Date:** December 21, 2025  
**Purpose:** Complete checklist for production deployment  
**Estimated Time:** 4-8 hours for full production setup  

---

## 📋 **PRE-DEPLOYMENT CHECKLIST**

### **1. Code & Build** ✅
- [ ] All features complete and tested
- [ ] No debug code or console.logs in production
- [ ] Error handling comprehensive
- [ ] Logging configured properly
- [ ] Backend builds successfully
- [ ] Frontend builds successfully
- [ ] All dependencies up to date
- [ ] Security vulnerabilities resolved

### **2. Database** ✅
- [ ] All migrations applied
- [ ] Database schema verified
- [ ] Indexes created for performance
- [ ] Backup strategy in place
- [ ] Connection pool configured
- [ ] SSL/TLS enabled for connections
- [ ] Read replicas configured (if needed)
- [ ] Automated backups scheduled

### **3. Authentication & Security** ✅
- [ ] JWT keys generated (RSA 2048-bit minimum)
- [ ] Keys stored securely (not in code)
- [ ] Token expiry times configured
- [ ] Rate limiting enabled
- [ ] Security headers configured
- [ ] CORS properly configured
- [ ] SQL injection protection verified
- [ ] XSS protection enabled
- [ ] CSRF protection (if using cookies)
- [ ] Password hashing verified (bcrypt cost 12+)

### **4. External Services**
- [ ] Twilio account created
- [ ] Twilio credentials secured
- [ ] SendGrid account created
- [ ] SendGrid API key secured
- [ ] Email sender verified
- [ ] Phone numbers purchased
- [ ] WhatsApp configured
- [ ] SMS templates approved
- [ ] Email templates approved
- [ ] Service health checks configured

### **5. Environment Configuration**
- [ ] Production .env file created
- [ ] All required variables set
- [ ] Secrets stored in secrets manager
- [ ] Environment variables documented
- [ ] No hardcoded credentials
- [ ] Database URLs correct
- [ ] API URLs correct
- [ ] Feature flags configured

### **6. Monitoring & Logging**
- [ ] Application logging configured
- [ ] Log aggregation set up (ELK, Splunk, etc.)
- [ ] Error tracking configured (Sentry, Rollbar)
- [ ] Performance monitoring (APM)
- [ ] Uptime monitoring
- [ ] Alert rules configured
- [ ] On-call rotation established
- [ ] Dashboard created

### **7. Performance**
- [ ] Load testing completed
- [ ] Database query optimization
- [ ] API response times < 200ms
- [ ] Frontend bundle optimized
- [ ] CDN configured for static assets
- [ ] Caching strategy implemented
- [ ] Connection pooling optimized
- [ ] Rate limits appropriate

### **8. Disaster Recovery**
- [ ] Backup strategy documented
- [ ] Recovery procedures tested
- [ ] Database backups automated
- [ ] Backup retention policy
- [ ] Disaster recovery plan
- [ ] Incident response plan
- [ ] Communication plan
- [ ] Team training completed

---

## 🚀 **DEPLOYMENT STEPS**

### **Phase 1: Infrastructure Setup** (2-3 hours)

#### **1.1 Database Setup**
```bash
# Create production database
createdb -h prod-db.amazonaws.com -U admin med_platform_prod

# Apply all migrations
psql -h prod-db.amazonaws.com -U admin -d med_platform_prod \
  -f database/migrations/*.sql

# Verify tables created
psql -h prod-db.amazonaws.com -U admin -d med_platform_prod \
  -c "\dt"

# Seed initial data (roles, etc.)
psql -h prod-db.amazonaws.com -U admin -d med_platform_prod \
  -f database/seed/production.sql
```

#### **1.2 Generate Production Keys**
```bash
# Generate JWT keys
cd keys
openssl genrsa -out jwt-private.pem 2048
openssl rsa -in jwt-private.pem -pubout -out jwt-public.pem

# Secure permissions
chmod 600 jwt-private.pem
chmod 644 jwt-public.pem

# Store in secrets manager
aws secretsmanager create-secret \
  --name prod/jwt/private-key \
  --secret-string file://jwt-private.pem

aws secretsmanager create-secret \
  --name prod/jwt/public-key \
  --secret-string file://jwt-public.pem
```

#### **1.3 Configure External Services**

**Twilio:**
```bash
# Set production environment variables
export TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export TWILIO_AUTH_TOKEN=your_production_token
export TWILIO_PHONE_NUMBER=+1234567890
export TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890
```

**SendGrid:**
```bash
export SENDGRID_API_KEY=SG.production_key_here
export SENDGRID_FROM_EMAIL=noreply@yourcompany.com
export SENDGRID_FROM_NAME="Your Company Name"
```

### **Phase 2: Application Deployment** (1-2 hours)

#### **2.1 Build Production Binary**
```bash
# Backend
cd backend
CGO_ENABLED=0 GOOS=linux go build -a \
  -ldflags '-extldflags "-static"' \
  -o platform \
  cmd/platform/main.go cmd/platform/init_auth.go

# Frontend
cd admin-ui
npm run build
```

#### **2.2 Deploy to Server**
```bash
# Copy binary to server
scp platform user@prod-server:/opt/med-platform/

# Copy frontend build
scp -r admin-ui/out user@prod-server:/var/www/med-platform/

# Set permissions
ssh user@prod-server "chmod +x /opt/med-platform/platform"
```

#### **2.3 Configure Systemd Service**
```bash
# Create service file
cat > /etc/systemd/system/med-platform.service <<EOF
[Unit]
Description=Medical Equipment Platform
After=network.target

[Service]
Type=simple
User=med-platform
WorkingDirectory=/opt/med-platform
EnvironmentFile=/opt/med-platform/.env
ExecStart=/opt/med-platform/platform
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
systemctl enable med-platform
systemctl start med-platform
systemctl status med-platform
```

#### **2.4 Configure Nginx**
```nginx
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com;
    
    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # Frontend
    root /var/www/med-platform;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### **Phase 3: Verification** (1 hour)

#### **3.1 Health Checks**
```bash
# Check backend health
curl https://yourdomain.com/health
# Expected: {"status":"ok"}

# Check metrics
curl https://yourdomain.com/metrics

# Check database connection
curl https://yourdomain.com/api/v1/equipment
```

#### **3.2 Authentication Testing**
```bash
# Test registration
curl -X POST https://yourdomain.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "full_name": "Test User",
    "password": "SecurePass123!"
  }'

# Check email/SMS received
# Verify OTP
# Test login
```

#### **3.3 Load Testing**
```bash
# Install k6 or Apache Bench
apt-get install apache2-utils

# Run load test
ab -n 1000 -c 10 https://yourdomain.com/api/v1/equipment

# Monitor response times, error rates
```

### **Phase 4: Monitoring Setup** (1-2 hours)

#### **4.1 Configure Prometheus**
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'med-platform'
    static_configs:
      - targets: ['localhost:8081']
```

#### **4.2 Configure Grafana**
- Import dashboard for Go applications
- Create custom dashboards for business metrics
- Set up alert rules

#### **4.3 Configure Alerts**
```yaml
# Alert rules
groups:
  - name: med-platform
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "High error rate detected"
          
      - alert: HighLatency
        expr: histogram_quantile(0.99, http_request_duration_seconds) > 1
        for: 5m
        annotations:
          summary: "High latency detected"
```

---

## 🔒 **SECURITY HARDENING**

### **1. Server Security**
```bash
# Update system
apt-get update && apt-get upgrade -y

# Configure firewall
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable

# Disable root login
sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
systemctl restart sshd

# Install fail2ban
apt-get install fail2ban
systemctl enable fail2ban
```

### **2. Application Security**
```bash
# Set secure permissions
chown -R med-platform:med-platform /opt/med-platform
chmod 700 /opt/med-platform
chmod 600 /opt/med-platform/.env

# Restrict file access
chmod 400 /opt/med-platform/keys/jwt-private.pem
```

### **3. Database Security**
```sql
-- Create read-only user for reporting
CREATE USER med_platform_readonly WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE med_platform_prod TO med_platform_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO med_platform_readonly;

-- Revoke unnecessary permissions
REVOKE ALL ON SCHEMA public FROM PUBLIC;
```

---

## 📊 **POST-DEPLOYMENT MONITORING**

### **Metrics to Monitor (First 24 Hours):**
- [ ] CPU usage < 70%
- [ ] Memory usage < 80%
- [ ] Disk usage < 80%
- [ ] API response times < 200ms (p99)
- [ ] Error rate < 0.1%
- [ ] Database connections stable
- [ ] No memory leaks
- [ ] No goroutine leaks

### **Business Metrics:**
- [ ] User registrations
- [ ] Login success rate
- [ ] OTP delivery rate
- [ ] Token refresh rate
- [ ] API usage patterns

### **Alert Thresholds:**
- **Critical:** Error rate > 5%, Downtime, Database unreachable
- **Warning:** Error rate > 1%, High latency, High CPU
- **Info:** New deployments, Configuration changes

---

## 🔄 **ROLLBACK PLAN**

If issues occur, follow this rollback procedure:

### **1. Immediate Actions**
```bash
# Stop new service
systemctl stop med-platform

# Start old service
systemctl start med-platform-old

# Verify old service working
curl https://yourdomain.com/health

# Update DNS if needed (if using blue-green)
```

### **2. Database Rollback**
```bash
# Revert migrations if needed
psql -h prod-db.amazonaws.com -U admin -d med_platform_prod \
  -f database/migrations/rollback.sql
```

### **3. Communication**
```
# Notify team
# Update status page
# Send customer communication if needed
```

---

## ✅ **POST-DEPLOYMENT CHECKLIST**

### **Immediate (0-4 hours):**
- [ ] All health checks passing
- [ ] Monitoring active
- [ ] Alerts configured
- [ ] Team notified
- [ ] Documentation updated
- [ ] Smoke tests passing

### **Day 1:**
- [ ] Monitor error rates
- [ ] Monitor performance
- [ ] Check log aggregation
- [ ] Verify backups running
- [ ] Test critical user flows
- [ ] Review metrics dashboard

### **Week 1:**
- [ ] Performance optimization
- [ ] Cost analysis
- [ ] User feedback review
- [ ] Security scan
- [ ] Load test results reviewed
- [ ] Documentation updated

---

## 📚 **DOCUMENTATION REQUIRED**

### **Operations Documentation:**
- [ ] Deployment procedure
- [ ] Rollback procedure
- [ ] Incident response plan
- [ ] On-call runbook
- [ ] Troubleshooting guide
- [ ] Architecture diagram

### **Development Documentation:**
- [ ] API documentation
- [ ] Database schema
- [ ] Environment variables
- [ ] Feature flags
- [ ] Testing procedures
- [ ] Development setup

---

## 🎯 **SUCCESS CRITERIA**

Deployment is considered successful when:
- ✅ All health checks passing
- ✅ Zero critical errors in first 24 hours
- ✅ API response times < 200ms (p99)
- ✅ User registration working
- ✅ Authentication working
- ✅ External services operational
- ✅ Monitoring active with alerts
- ✅ Backups running successfully

---

**Document:** Production Deployment Checklist  
**Last Updated:** December 21, 2025  
**Status:** Ready for production deployment  
**Estimated Deployment Time:** 4-8 hours for complete setup
