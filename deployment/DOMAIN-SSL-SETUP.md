# Domain and SSL Certificate Setup for servqr.com

This guide explains how to configure your deployment to use `servqr.com` domain with automatic SSL certificates from Let's Encrypt.

## Prerequisites

Before starting, ensure:

1. ✅ **Domain DNS is configured:**
   - `servqr.com` A record → Points to your server's IP address
   - `www.servqr.com` A record → Points to your server's IP address (optional)
   - DNS propagation completed (check with `nslookup servqr.com`)

2. ✅ **Ports are open:**
   - Port 80 (HTTP) - Required for Let's Encrypt verification
   - Port 443 (HTTPS) - For secure traffic
   - Check firewall: `sudo ufw status` or cloud provider security groups

3. ✅ **Valid email address:**
   - For Let's Encrypt certificate notifications
   - Example: `admin@servqr.com` or your email

## Configuration Steps

### Option 1: Configure Before Deployment (Recommended)

**Edit the deployment script:**

```bash
# Edit deploy-all.sh
sudo nano /opt/servqr/deployment/deploy-all.sh
```

**Update these lines (around line 37-38):**

```bash
# Before:
DOMAIN=""  # Leave empty to skip SSL setup
EMAIL=""   # Leave empty to skip SSL setup

# After:
DOMAIN="servqr.com"                    # Your domain
EMAIL="admin@servqr.com"               # Your email
```

**Then run deployment:**

```bash
cd /opt/servqr/deployment
sudo bash deploy-all.sh
```

The script will:
1. ✅ Install Nginx
2. ✅ Create configuration for servqr.com
3. ✅ Install certbot
4. ✅ Obtain SSL certificate from Let's Encrypt
5. ✅ Configure auto-renewal
6. ✅ Set up HTTPS redirect

---

### Option 2: Configure After Deployment (Manual)

If you already deployed with IP address, you can add the domain later:

**Step 1: Update deployment script**

```bash
sudo nano /opt/servqr/deployment/deploy-all.sh
```

Update DOMAIN and EMAIL variables as shown above.

**Step 2: Run Nginx configuration script**

```bash
cd /opt/servqr/deployment
sudo bash configure-nginx.sh servqr.com admin@servqr.com
```

**Or manually configure Nginx:**

```bash
# Create Nginx configuration
sudo nano /etc/nginx/sites-available/servqr
```

**Add this configuration:**

```nginx
# HTTP - Redirect to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name servqr.com www.servqr.com;
    
    # Let's Encrypt verification
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }
    
    # Redirect all other traffic to HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS - Main application
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name servqr.com www.servqr.com;
    
    # SSL certificates (certbot will create these)
    ssl_certificate /etc/letsencrypt/live/servqr.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/servqr.com/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-CHACHA20-POLY1305;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Frontend (Next.js)
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_read_timeout 86400;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # Health check
    location /health {
        proxy_pass http://localhost:8081/health;
        access_log off;
    }
    
    # File upload size limit
    client_max_body_size 50M;
    client_body_buffer_size 128k;
    
    # Logging
    access_log /var/log/nginx/servqr-access.log;
    error_log /var/log/nginx/servqr-error.log;
}
```

**Step 3: Enable site and test**

```bash
# Enable site
sudo ln -sf /etc/nginx/sites-available/servqr /etc/nginx/sites-enabled/

# Remove default site if exists
sudo rm -f /etc/nginx/sites-enabled/default

# Test configuration
sudo nginx -t
```

**Step 4: Install certbot and obtain certificate**

```bash
# Install certbot
sudo apt-get update
sudo apt-get install -y certbot python3-certbot-nginx

# Obtain certificate (interactive)
sudo certbot --nginx -d servqr.com -d www.servqr.com --email admin@servqr.com --agree-tos --no-eff-email

# Or non-interactive
sudo certbot --nginx -d servqr.com -d www.servqr.com --email admin@servqr.com --agree-tos --no-eff-email --non-interactive
```

**Step 5: Reload Nginx**

```bash
sudo systemctl reload nginx
```

---

## Verification

### 1. Check DNS Resolution

```bash
# Check if domain resolves to your server
nslookup servqr.com
dig servqr.com +short

# Should show your server's IP address
```

### 2. Check Nginx Configuration

```bash
# Test Nginx config
sudo nginx -t

# Check Nginx is running
sudo systemctl status nginx
```

### 3. Check SSL Certificate

```bash
# List certificates
sudo certbot certificates

# Should show:
# Certificate Name: servqr.com
# Domains: servqr.com www.servqr.com
# Expiry Date: (90 days from now)
# Certificate Path: /etc/letsencrypt/live/servqr.com/fullchain.pem
```

### 4. Test HTTPS Access

```bash
# Test from server
curl -I https://servqr.com

# Should return: HTTP/2 200

# Test SSL
curl -v https://servqr.com 2>&1 | grep "SSL certificate verify"
# Should return: SSL certificate verify ok
```

### 5. Test in Browser

1. Visit: `http://servqr.com`
   - Should redirect to `https://servqr.com`
   
2. Visit: `https://servqr.com`
   - Should show green padlock 🔒
   - Certificate issued by Let's Encrypt
   - Valid for 90 days

3. Check certificate details:
   - Click padlock → Certificate
   - Issued to: servqr.com
   - Issued by: Let's Encrypt

---

## Certificate Auto-Renewal

Certbot automatically sets up a cron job for certificate renewal.

### Check Auto-Renewal

```bash
# Check renewal timer
sudo systemctl status certbot.timer

# Test renewal (dry run)
sudo certbot renew --dry-run

# Manual renewal (if needed)
sudo certbot renew
```

### Renewal Cron Job

Certbot creates: `/etc/cron.d/certbot`

```bash
# View cron job
cat /etc/cron.d/certbot

# Typical content:
# 0 */12 * * * root test -x /usr/bin/certbot && perl -e 'sleep int(rand(43200))' && certbot -q renew
```

---

## Update Application URLs

After SSL is configured, update your application to use HTTPS URLs:

### 1. Update .env File

```bash
sudo nano /opt/servqr/.env
```

**Update these values:**

```bash
# API URLs
NEXT_PUBLIC_API_BASE_URL=https://servqr.com/api
API_BASE_URL=https://servqr.com/api

# Application URLs
APP_URL=https://servqr.com
NEXT_PUBLIC_APP_URL=https://servqr.com

# Tracking URL (already updated in code)
TRACKING_URL=https://servqr.com

# Email Links (already updated in code)
# These are now hardcoded in the backend
```

### 2. Restart Services

```bash
# Restart backend
sudo systemctl restart servqr-backend

# Restart frontend
sudo systemctl restart servqr-frontend

# Verify both are running
sudo systemctl status servqr-backend
sudo systemctl status servqr-frontend
```

---

## Troubleshooting

### Issue: "Connection refused" on port 80

```bash
# Check if port 80 is open
sudo netstat -tulpn | grep :80

# Check firewall
sudo ufw status
sudo ufw allow 80
sudo ufw allow 443
```

### Issue: "DNS resolution failed"

```bash
# Check DNS
nslookup servqr.com

# Wait for DNS propagation (can take up to 48 hours)
# Check propagation: https://www.whatsmydns.net/#A/servqr.com
```

### Issue: "Certificate validation failed"

```bash
# Check Nginx error logs
sudo tail -f /var/log/nginx/error.log

# Check Let's Encrypt logs
sudo tail -f /var/log/letsencrypt/letsencrypt.log

# Common fixes:
# 1. Ensure port 80 is accessible (Let's Encrypt verification)
# 2. Ensure domain points to correct IP
# 3. Check Nginx is running: sudo systemctl status nginx
```

### Issue: "Too many certificates already issued"

Let's Encrypt has rate limits:
- 50 certificates per registered domain per week
- 5 certificates per FQDN per week

**Solution:** Wait a week or use staging environment:

```bash
# Test with staging (doesn't count against rate limit)
sudo certbot --nginx -d servqr.com --staging
```

### Issue: Nginx test fails

```bash
# Check Nginx config syntax
sudo nginx -t

# View detailed errors
sudo nginx -T

# Common issues:
# - Missing semicolons
# - Incorrect paths
# - Port conflicts
```

---

## Security Best Practices

### 1. Enable Firewall

```bash
# Install UFW
sudo apt-get install ufw

# Allow SSH (important!)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

### 2. Setup Fail2Ban (Optional)

```bash
# Install
sudo apt-get install fail2ban

# Create Nginx jail
sudo nano /etc/fail2ban/jail.local
```

Add:
```ini
[nginx-http-auth]
enabled = true

[nginx-noscript]
enabled = true

[nginx-badbots]
enabled = true
```

### 3. Regular Updates

```bash
# Update system
sudo apt-get update && sudo apt-get upgrade -y

# Update certbot
sudo apt-get install --only-upgrade certbot python3-certbot-nginx
```

---

## Summary Checklist

Before going live with `servqr.com`:

- [ ] DNS A record configured and propagated
- [ ] Ports 80 and 443 open in firewall
- [ ] DOMAIN and EMAIL set in deploy-all.sh
- [ ] Nginx installed and configured
- [ ] SSL certificate obtained from Let's Encrypt
- [ ] HTTPS redirect working
- [ ] Auto-renewal configured
- [ ] Application .env updated with HTTPS URLs
- [ ] Services restarted
- [ ] Browser shows green padlock
- [ ] API endpoints accessible via HTTPS

---

## Support

If you encounter issues:

1. Check logs:
   - Nginx: `/var/log/nginx/error.log`
   - Certbot: `/var/log/letsencrypt/letsencrypt.log`
   - Application: `/opt/servqr/logs/`

2. Test connectivity:
   ```bash
   curl -I https://servqr.com
   curl -I https://servqr.com/api/health
   ```

3. Verify DNS:
   ```bash
   nslookup servqr.com
   ```

---

**Document Version:** 1.0
**Last Updated:** February 9, 2026
**Deployment Scripts:** Compatible with deploy-all.sh
