#!/bin/bash

################################################################################
# ServQR Platform - Nginx Reverse Proxy Configuration
#
# Configures Nginx as a reverse proxy for frontend and backend
# - Frontend: / -> localhost:3000
# - Backend API: /api -> localhost:8081
#
# Usage: sudo bash configure-nginx.sh [domain-or-ip]
################################################################################

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   error "This script must be run as root (use sudo)"
fi

# Get server address (domain or IP)
if [[ $# -eq 1 ]]; then
    SERVER_ADDRESS="$1"
else
    # Use public IP if no argument provided
    SERVER_ADDRESS=$(hostname -I | awk '{print $1}')
fi

log "Configuring Nginx for: $SERVER_ADDRESS"

# Check if Nginx is installed
if ! command -v nginx &> /dev/null; then
    log "Installing Nginx..."
    apt-get update
    apt-get install -y nginx
fi

# Create Nginx configuration
log "Creating Nginx configuration..."

cat > /etc/nginx/sites-available/servqr << EOF
# ServQR Platform - Nginx Configuration
# Generated: $(date)

# Rate limiting zone
limit_req_zone \$binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone \$binary_remote_addr zone=general_limit:10m rate=30r/s;

# Upstream servers
upstream servqr_backend {
    server localhost:8081 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

upstream servqr_frontend {
    server localhost:3000 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

server {
    listen 3000;
    listen [::]:3000;
    
    server_name $SERVER_ADDRESS;
    
    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Logging
    access_log /var/log/nginx/servqr_access.log;
    error_log /var/log/nginx/servqr_error.log warn;
    
    # Max upload size
    client_max_body_size 50M;
    
    # Backend API (proxy to backend without modifying path)
    location /api/ {
        # Rate limiting
        limit_req zone=api_limit burst=20 nodelay;
        
        # Proxy settings (pass full path including /api/)
        proxy_pass http://servqr_backend;
        proxy_http_version 1.1;
        
        # Headers
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
    }
    
    # Health check endpoint
    location /health {
        proxy_pass http://servqr_backend/health;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header Connection "";
        access_log off;
    }
    
    # WebSocket support (if needed)
    location /ws {
        proxy_pass http://servqr_backend/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_read_timeout 86400;
    }
    
    # Static files (served by Next.js)
    location /_next/static/ {
        proxy_pass http://servqr_frontend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        
        # Cache static assets
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # Next.js internal routes
    location /_next/ {
        proxy_pass http://servqr_frontend;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
    }
    
    # Frontend application
    location / {
        # Rate limiting
        limit_req zone=general_limit burst=50 nodelay;
        
        # Proxy to Next.js
        proxy_pass http://servqr_frontend;
        proxy_http_version 1.1;
        
        # Headers
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering on;
    }
}
EOF

# Enable the site
log "Enabling ServQR site..."
ln -sf /etc/nginx/sites-available/servqr /etc/nginx/sites-enabled/servqr

# Remove default site if it exists
if [[ -L /etc/nginx/sites-enabled/default ]]; then
    rm /etc/nginx/sites-enabled/default
fi

# Test Nginx configuration
log "Testing Nginx configuration..."
if nginx -t; then
    log "✓ Nginx configuration is valid"
else
    error "Nginx configuration test failed"
fi

# Reload Nginx
log "Reloading Nginx..."
systemctl reload nginx

# Enable Nginx to start on boot
systemctl enable nginx

# Update backend .env
log "Updating backend configuration..."
sed -i "s|BASE_URL=.*|BASE_URL=http://${SERVER_ADDRESS}|g" /opt/servqr/.env
sed -i "s|FRONTEND_URL=.*|FRONTEND_URL=http://${SERVER_ADDRESS}|g" /opt/servqr/.env

# Update frontend .env.local
log "Updating frontend configuration..."
sed -i "s|NEXT_PUBLIC_API_URL=.*|NEXT_PUBLIC_API_URL=http://${SERVER_ADDRESS}/api/v1|g" /opt/servqr/admin-ui/.env.local
sed -i "s|NEXT_PUBLIC_BASE_URL=.*|NEXT_PUBLIC_BASE_URL=http://${SERVER_ADDRESS}|g" /opt/servqr/admin-ui/.env.local

# Restart services to pick up new config
log "Restarting services..."
systemctl restart servqr-backend
systemctl restart servqr-frontend

# Display status
log ""
log "✓ Nginx configuration completed successfully!"
log ""
log "Access Information:"
log "  Frontend:    http://${SERVER_ADDRESS}/"
log "  Backend API: http://${SERVER_ADDRESS}/api/"
log "  Health:      http://${SERVER_ADDRESS}/health"
log ""
log "Services:"
log "  Nginx:    systemctl status nginx"
log "  Backend:  systemctl status servqr-backend"
log "  Frontend: systemctl status servqr-frontend"
log ""
log "Logs:"
log "  Nginx access:  tail -f /var/log/nginx/servqr_access.log"
log "  Nginx error:   tail -f /var/log/nginx/servqr_error.log"
log ""

# Test the setup
log "Testing endpoints..."
sleep 3

if curl -f -s http://localhost/health > /dev/null 2>&1; then
    log "✓ Health check: OK"
else
    warn "Health check failed - backend may still be starting"
fi

if curl -f -s http://localhost/ > /dev/null 2>&1; then
    log "✓ Frontend: OK"
else
    warn "Frontend check failed - service may still be starting"
fi

log ""
log "Configuration complete!"
log "Access your application at: http://${SERVER_ADDRESS}"
