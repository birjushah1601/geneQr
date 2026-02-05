#!/bin/bash

################################################################################
# ServQR Platform - Master Deployment Script
# 
# This script orchestrates the complete deployment of ServQR platform on a
# plain Linux VM. It installs all dependencies, sets up Docker, PostgreSQL,
# builds the application, and configures services.
#
# Usage: sudo bash deploy-all.sh
#
# Prerequisites:
# - Source code at /opt/servqr
# - Root/sudo access
# - Internet connectivity
################################################################################

set -e  # Exit on error
set -u  # Exit on undefined variable

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/opt/servqr"
DEPLOYMENT_DIR="${INSTALL_DIR}/deployment"

# Create logs directory first (before using LOG_FILE)
mkdir -p "${INSTALL_DIR}/logs"

LOG_FILE="${INSTALL_DIR}/logs/deployment-$(date +%Y%m%d-%H%M%S).log"

# Optional: Set your domain and email for SSL
DOMAIN=""  # Leave empty to skip SSL setup
EMAIL=""   # Leave empty to skip SSL setup

################################################################################
# Helper Functions
################################################################################

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
    exit 1
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

print_header() {
    echo ""
    echo "=========================================================================="
    echo "  $1"
    echo "=========================================================================="
    echo ""
}

check_root() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use sudo)"
    fi
}

check_directory() {
    if [[ ! -d "$INSTALL_DIR" ]]; then
        error "Source code not found at $INSTALL_DIR. Please clone the repository first."
    fi
    
    if [[ ! -d "$DEPLOYMENT_DIR" ]]; then
        error "Deployment directory not found at $DEPLOYMENT_DIR"
    fi
}

################################################################################
# Main Deployment Steps
################################################################################

main() {
    print_header "ServQR Platform - Production Deployment"
    
    log "Starting deployment at $(date)"
    log "Installation directory: $INSTALL_DIR"
    
    # Pre-flight checks
    check_root
    check_directory
    
    # Create required directories
    log "Creating required directories..."
    mkdir -p "${INSTALL_DIR}/logs"
    mkdir -p "${INSTALL_DIR}/data/postgres"
    mkdir -p "${INSTALL_DIR}/data/qrcodes"
    mkdir -p "${INSTALL_DIR}/data/whatsapp"
    mkdir -p "${INSTALL_DIR}/backups"
    mkdir -p "${INSTALL_DIR}/storage"
    
    # Step 1: Install system prerequisites
    print_header "Step 1/4: Installing System Prerequisites"
    log "Installing Go, Node.js, Nginx, and other dependencies..."
    bash "${DEPLOYMENT_DIR}/install-prerequisites.sh" || error "Prerequisites installation failed"
    
    # Step 2: Setup Docker and PostgreSQL
    print_header "Step 2/4: Setting up Docker and PostgreSQL"
    log "Installing Docker and creating PostgreSQL container..."
    bash "${DEPLOYMENT_DIR}/setup-docker.sh" || error "Docker setup failed"
    
    # Step 3: Deploy application
    print_header "Step 3/4: Building and Deploying Application"
    log "Building backend and frontend..."
    bash "${DEPLOYMENT_DIR}/deploy-app.sh" || error "Application deployment failed"
    
    # Step 4: Configure services
    print_header "Step 4/4: Configuring System Services"
    log "Setting up systemd services and auto-start..."
    
    # Copy systemd service files
    cp "${DEPLOYMENT_DIR}/systemd/servqr-backend.service" /etc/systemd/system/
    cp "${DEPLOYMENT_DIR}/systemd/servqr-frontend.service" /etc/systemd/system/
    
    # Reload systemd and enable services
    systemctl daemon-reload
    systemctl enable servqr-backend
    systemctl enable servqr-frontend
    systemctl start servqr-backend
    systemctl start servqr-frontend
    
    log "Services configured and started"
    
    # Setup Nginx reverse proxy
    log "Configuring Nginx reverse proxy..."
    SERVER_IP=$(hostname -I | awk '{print $1}')
    
    if [[ -n "$DOMAIN" ]]; then
        log "Using domain: $DOMAIN"
        bash "${DEPLOYMENT_DIR}/configure-nginx.sh" "$DOMAIN" || warn "Nginx setup failed"
    else
        log "Using IP address: $SERVER_IP"
        bash "${DEPLOYMENT_DIR}/configure-nginx.sh" "$SERVER_IP" || warn "Nginx setup failed"
    fi
    
    # Setup automated backups
    log "Configuring automated database backups..."
    setup_backups
    
    # Setup log rotation
    log "Configuring log rotation..."
    setup_log_rotation
    
    # Final health checks
    print_header "Running Health Checks"
    run_health_checks
    
    # Print summary
    print_summary
    
    log "Deployment completed successfully at $(date)"
}

setup_nginx() {
    # Install Nginx if not already installed
    if ! command -v nginx &> /dev/null; then
        log "Nginx not found. Installing..."
        apt-get install -y nginx || yum install -y nginx
    fi
    
    # Create Nginx configuration
    cat > /etc/nginx/sites-available/servqr << EOF
server {
    listen 80;
    server_name ${DOMAIN};
    
    # Redirect to HTTPS
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name ${DOMAIN};
    
    # SSL certificates (will be generated by certbot)
    ssl_certificate /etc/letsencrypt/live/${DOMAIN}/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/${DOMAIN}/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # Frontend proxy
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # Backend API proxy
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8081/health;
    }
    
    # File upload size limit
    client_max_body_size 50M;
}
EOF
    
    # Enable site
    ln -sf /etc/nginx/sites-available/servqr /etc/nginx/sites-enabled/
    
    # Test Nginx configuration
    nginx -t || error "Nginx configuration test failed"
    
    # Install certbot and obtain SSL certificate
    if [[ -n "$EMAIL" ]]; then
        log "Installing certbot and obtaining SSL certificate..."
        if command -v apt-get &> /dev/null; then
            apt-get install -y certbot python3-certbot-nginx
        else
            yum install -y certbot python3-certbot-nginx
        fi
        
        certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "$EMAIL" || warn "SSL certificate setup failed. You can run 'sudo certbot --nginx' manually."
    fi
    
    # Reload Nginx
    systemctl reload nginx
    log "Nginx configured successfully"
}

setup_backups() {
    # Create backup script
    cat > "${INSTALL_DIR}/deployment/backup-database.sh" << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/servqr/backups"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/servqr-backup-${TIMESTAMP}.sql"

# Create backup
docker exec servqr-postgres pg_dump -U servqr servqr_production > "$BACKUP_FILE"

# Compress backup
gzip "$BACKUP_FILE"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: ${BACKUP_FILE}.gz"
EOF
    
    chmod +x "${INSTALL_DIR}/deployment/backup-database.sh"
    
    # Add cron job for daily backups at 2 AM
    (crontab -l 2>/dev/null; echo "0 2 * * * ${INSTALL_DIR}/deployment/backup-database.sh >> ${INSTALL_DIR}/logs/backup.log 2>&1") | crontab -
    
    log "Automated backups configured (daily at 2 AM)"
}

setup_log_rotation() {
    cat > /etc/logrotate.d/servqr << 'EOF'
/opt/servqr/logs/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0644 root root
    sharedscripts
    postrotate
        systemctl reload servqr-backend > /dev/null 2>&1 || true
        systemctl reload servqr-frontend > /dev/null 2>&1 || true
    endscript
}
EOF
    
    log "Log rotation configured"
}

run_health_checks() {
    log "Running health checks..."
    
    sleep 5  # Give services time to start
    
    # Check backend
    if curl -f http://localhost:8081/health &> /dev/null; then
        log "✓ Backend health check passed"
    else
        warn "Backend health check failed. Check logs: journalctl -u servqr-backend -n 50"
    fi
    
    # Check frontend
    if curl -f http://localhost:3000 &> /dev/null; then
        log "✓ Frontend health check passed"
    else
        warn "Frontend health check failed. Check logs: journalctl -u servqr-frontend -n 50"
    fi
    
    # Check database
    if docker exec servqr-postgres pg_isready -U servqr &> /dev/null; then
        log "✓ Database health check passed"
    else
        warn "Database health check failed. Check logs: docker logs servqr-postgres"
    fi
    
    # Check services status
    systemctl is-active --quiet servqr-backend && log "✓ Backend service is active" || warn "Backend service is not active"
    systemctl is-active --quiet servqr-frontend && log "✓ Frontend service is active" || warn "Frontend service is not active"
}

print_summary() {
    print_header "Deployment Summary"
    
    echo "✓ ServQR Platform deployed successfully!"
    echo ""
    echo "Service Status:"
    systemctl status servqr-backend --no-pager -l | head -5
    echo ""
    systemctl status servqr-frontend --no-pager -l | head -5
    echo ""
    
    echo "Access Information:"
    if [[ -n "$DOMAIN" ]]; then
        echo "  Frontend: https://${DOMAIN}"
        echo "  Backend API: https://${DOMAIN}/api"
    else
        SERVER_IP=$(hostname -I | awk '{print $1}')
        echo "  Frontend: http://${SERVER_IP}:3000"
        echo "  Backend API: http://${SERVER_IP}:8081/api"
    fi
    echo ""
    
    echo "Default Login Credentials:"
    echo "  Email: admin@servqr.com"
    echo "  Password: Check /opt/servqr/LOGIN-CREDENTIALS.txt"
    echo ""
    
    echo "Management Commands:"
    echo "  View backend logs:   sudo journalctl -u servqr-backend -f"
    echo "  View frontend logs:  sudo journalctl -u servqr-frontend -f"
    echo "  Restart backend:     sudo systemctl restart servqr-backend"
    echo "  Restart frontend:    sudo systemctl restart servqr-frontend"
    echo "  Database backup:     sudo ${INSTALL_DIR}/deployment/backup-database.sh"
    echo "  Database access:     docker exec -it servqr-postgres psql -U servqr -d servqr_production"
    echo ""
    
    echo "Important Files:"
    echo "  Deployment logs:     ${LOG_FILE}"
    echo "  Environment config:  ${INSTALL_DIR}/.env"
    echo "  Database backups:    ${INSTALL_DIR}/backups/"
    echo "  Application logs:    ${INSTALL_DIR}/logs/"
    echo ""
    
    echo "Next Steps:"
    echo "  1. Update environment variables: vim ${INSTALL_DIR}/.env"
    echo "  2. Configure external services (SendGrid, Twilio, AI APIs)"
    echo "  3. Test the application by visiting the URLs above"
    echo "  4. Review security checklist in deployment/README.md"
    echo "  5. Setup monitoring and alerts"
    echo ""
    
    warn "IMPORTANT: Change default passwords and update JWT secret before going to production!"
    echo ""
}

# Run main deployment
main "$@"
