#!/bin/bash

################################################################################
# ServQR Platform - Application Deployment Script
#
# Builds backend and frontend, configures environment, starts services
#
# Usage: sudo bash deploy-app.sh
################################################################################

set -e
set -u

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
info() { echo -e "${BLUE}[INFO]${NC} $1"; }

INSTALL_DIR="/opt/servqr"
DEPLOYMENT_DIR="${INSTALL_DIR}/deployment"
BACKEND_DIR="${INSTALL_DIR}"
FRONTEND_DIR="${INSTALL_DIR}/admin-ui"

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if running as root
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use sudo)"
    fi
    
    # Check if installation directory exists
    if [[ ! -d "$INSTALL_DIR" ]]; then
        error "Installation directory not found: $INSTALL_DIR"
    fi
    
    # Check Go
    if ! command -v go &> /dev/null; then
        error "Go not found. Run install-prerequisites.sh first"
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        error "Node.js not found. Run install-prerequisites.sh first"
    fi
    
    # Check if PostgreSQL is running
    if ! docker ps | grep servqr-postgres &> /dev/null; then
        error "PostgreSQL container not running. Run setup-docker.sh first"
    fi
    
    log "All prerequisites met"
}

# Configure environment variables
configure_environment() {
    log "Configuring environment variables..."
    
    # Load database password
    if [[ -f "${INSTALL_DIR}/.db_password" ]]; then
        DB_PASSWORD=$(cat "${INSTALL_DIR}/.db_password")
    else
        error "Database password file not found. Run setup-docker.sh first"
    fi
    
    # Generate JWT secret if not exists
    if [[ ! -f "${INSTALL_DIR}/.jwt_secret" ]]; then
        JWT_SECRET=$(openssl rand -base64 32)
        echo "$JWT_SECRET" > "${INSTALL_DIR}/.jwt_secret"
        chmod 600 "${INSTALL_DIR}/.jwt_secret"
    else
        JWT_SECRET=$(cat "${INSTALL_DIR}/.jwt_secret")
    fi
    
    # Get server IP
    SERVER_IP=$(hostname -I | awk '{print $1}')
    
    # Create .env file for backend
    cat > "${INSTALL_DIR}/.env" << EOF
# ========================================================================
# ServQR Platform - Production Environment Configuration
# Generated: $(date)
# ========================================================================

# Environment
ENVIRONMENT=production
VERSION=1.0.0

# Server Configuration
PORT=8081
BASE_URL=http://${SERVER_IP}:8081
FRONTEND_URL=http://${SERVER_IP}:3000

# Database Configuration (DATABASE_ prefix)
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=servqr
DATABASE_PASSWORD=${DB_PASSWORD}
DATABASE_NAME=servqr_production
DATABASE_SSL_MODE=disable

# Database Configuration (DB_ prefix - required by Go code)
DB_HOST=localhost
DB_PORT=5432
DB_USER=servqr
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=servqr_production
DB_SSLMODE=disable

# Database Connection String
DATABASE_URL=postgresql://servqr:${DB_PASSWORD}@localhost:5432/servqr_production?sslmode=disable

# JWT Authentication
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRATION_HOURS=24

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://${SERVER_IP}:3000

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# File Storage
QR_OUTPUT_DIR=${INSTALL_DIR}/data/qrcodes
WHATSAPP_MEDIA_DIR=${INSTALL_DIR}/data/whatsapp
STORAGE_PATH=${INSTALL_DIR}/storage

# Feature Flags
ENABLE_ORG=true
ENABLE_EQUIPMENT=true
ENABLE_WHATSAPP=false
ENABLE_AI_DIAGNOSIS=true

# AI Configuration (Optional - Add your API keys)
AI_PROVIDER=openai
AI_FALLBACK_PROVIDER=anthropic
AI_OPENAI_API_KEY=
AI_OPENAI_MODEL=gpt-4
AI_ANTHROPIC_API_KEY=
AI_ANTHROPIC_MODEL=claude-3-opus-20240229
AI_TIMEOUT_SECONDS=30
AI_MAX_RETRIES=3

# Email Configuration (Optional - Add your SendGrid key)
SENDGRID_API_KEY=
SENDGRID_FROM_EMAIL=noreply@servqr.com
SENDGRID_FROM_NAME=ServQR Platform
FEATURE_EMAIL_NOTIFICATIONS=false

# WhatsApp Configuration (Optional - Add your Twilio credentials)
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_WHATSAPP_NUMBER=
WHATSAPP_VERIFY_TOKEN=

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# Observability
OBSERVABILITY_LOG_LEVEL=info
OBSERVABILITY_TRACING_ENABLED=false
OBSERVABILITY_METRICS_ENABLED=true

# Security
SECURITY_HEADERS_ENABLED=true
SECURITY_RATE_LIMITING_ENABLED=true
EOF
    
    # Create .env.local for frontend
    cat > "${FRONTEND_DIR}/.env.local" << EOF
# ServQR Frontend Configuration
NEXT_PUBLIC_API_URL=http://${SERVER_IP}:8081/api/v1
NEXT_PUBLIC_BASE_URL=http://${SERVER_IP}:3000
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1
EOF
    
    log "Environment configuration created"
    info "Database password saved to: ${INSTALL_DIR}/.db_password"
    info "JWT secret saved to: ${INSTALL_DIR}/.jwt_secret"
    info "Backend config: ${INSTALL_DIR}/.env"
    info "Frontend config: ${FRONTEND_DIR}/.env.local"
}

# Generate JWT keys
generate_jwt_keys() {
    log "Generating JWT RSA keys..."
    
    # Create keys directory
    mkdir -p "${INSTALL_DIR}/keys"
    
    # Generate private key if not exists
    if [[ ! -f "${INSTALL_DIR}/keys/jwt-private.pem" ]]; then
        openssl genrsa -out "${INSTALL_DIR}/keys/jwt-private.pem" 4096
        chmod 600 "${INSTALL_DIR}/keys/jwt-private.pem"
        log "JWT private key generated"
    else
        log "JWT private key already exists"
    fi
    
    # Generate public key if not exists
    if [[ ! -f "${INSTALL_DIR}/keys/jwt-public.pem" ]]; then
        openssl rsa -in "${INSTALL_DIR}/keys/jwt-private.pem" -pubout -out "${INSTALL_DIR}/keys/jwt-public.pem"
        chmod 644 "${INSTALL_DIR}/keys/jwt-public.pem"
        log "JWT public key generated"
    else
        log "JWT public key already exists"
    fi
    
    log "JWT keys ready"
}

# Fix database schema issues
fix_database_schema() {
    log "Fixing database schema issues..."
    
    # Fix ticket_comments table to use UUID instead of VARCHAR
    docker exec servqr-postgres psql -U servqr -d servqr_production << 'SQLFIX' 2>/dev/null || true
-- Drop and recreate ticket_comments with correct UUID type
DROP TABLE IF EXISTS ticket_comments CASCADE;

CREATE TABLE IF NOT EXISTS ticket_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES service_tickets(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_ticket_comments_ticket ON ticket_comments(ticket_id);
SQLFIX
    
    log "Database schema fixes applied"
}

# Build backend
build_backend() {
    log "Building backend..."
    
    cd "$BACKEND_DIR"
    
    # Set Go environment
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export CGO_ENABLED=0
    
    # Clean previous build
    rm -f platform platform.exe
    
    # Build binary
    log "Compiling Go binary..."
    go build -o platform \
        -ldflags="-s -w -X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y%m%d.%H%M%S)" \
        ./cmd/platform/ || error "Backend build failed"
    
    # Make executable
    chmod +x platform
    
    # Verify build
    if [[ ! -f "platform" ]]; then
        error "Backend binary not found after build"
    fi
    
    log "Backend built successfully"
    log "Binary size: $(du -h platform | cut -f1)"
}

# Build frontend
build_frontend() {
    log "Building frontend..."
    
    cd "$FRONTEND_DIR"
    
    # Install dependencies
    log "Installing npm dependencies..."
    # Use npm install instead of npm ci (ci requires package-lock.json)
    npm install || error "npm install failed"
    
    # Build Next.js application
    log "Building Next.js application..."
    npm run build || error "Frontend build failed"
    
    # Verify build
    if [[ ! -d ".next" ]]; then
        error "Frontend build directory (.next) not found"
    fi
    
    log "Frontend built successfully"
    log "Build directory: ${FRONTEND_DIR}/.next"
}

# Create systemd services
create_systemd_services() {
    log "Creating systemd service files..."
    
    # Backend service
    cat > /etc/systemd/system/servqr-backend.service << EOF
[Unit]
Description=ServQR Backend API
After=network.target servqr-postgres.service
Requires=servqr-postgres.service

[Service]
Type=simple
User=root
WorkingDirectory=${INSTALL_DIR}
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
EnvironmentFile=${INSTALL_DIR}/.env
ExecStart=${INSTALL_DIR}/platform
Restart=always
RestartSec=10
StandardOutput=append:${INSTALL_DIR}/logs/backend.log
StandardError=append:${INSTALL_DIR}/logs/backend-error.log

# Resource Limits
LimitNOFILE=65536
LimitNPROC=32768

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
    
    # Frontend service
    cat > /etc/systemd/system/servqr-frontend.service << EOF
[Unit]
Description=ServQR Frontend (Next.js)
After=network.target servqr-backend.service

[Service]
Type=simple
User=root
WorkingDirectory=${FRONTEND_DIR}
Environment="NODE_ENV=production"
Environment="PORT=3000"
ExecStart=/usr/bin/npm start
Restart=always
RestartSec=10
StandardOutput=append:${INSTALL_DIR}/logs/frontend.log
StandardError=append:${INSTALL_DIR}/logs/frontend-error.log

# Resource Limits
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
    
    log "Systemd service files created"
}

# Start services
start_services() {
    log "Starting services..."
    
    # Reload systemd
    systemctl daemon-reload
    
    # Enable services to start on boot
    systemctl enable servqr-backend
    systemctl enable servqr-frontend
    
    # Stop services if running
    systemctl stop servqr-backend 2>/dev/null || true
    systemctl stop servqr-frontend 2>/dev/null || true
    
    # Start backend
    log "Starting backend service..."
    systemctl start servqr-backend
    
    # Wait for backend to be ready
    log "Waiting for backend to be ready..."
    max_retries=30
    retry_count=0
    
    while ! curl -f http://localhost:8081/health &> /dev/null; do
        retry_count=$((retry_count + 1))
        if [[ $retry_count -ge $max_retries ]]; then
            warn "Backend health check timeout. Check logs: journalctl -u servqr-backend -n 50"
            break
        fi
        sleep 2
    done
    
    if curl -f http://localhost:8081/health &> /dev/null; then
        log "✓ Backend is healthy"
    fi
    
    # Start frontend
    log "Starting frontend service..."
    systemctl start servqr-frontend
    
    # Wait for frontend to be ready
    log "Waiting for frontend to be ready..."
    sleep 10
    
    if curl -f http://localhost:3000 &> /dev/null; then
        log "✓ Frontend is healthy"
    else
        warn "Frontend health check failed. Check logs: journalctl -u servqr-frontend -n 50"
    fi
    
    log "Services started"
}

# Display service status
display_status() {
    log "Service Status:"
    
    echo ""
    systemctl status servqr-backend --no-pager -l | head -10
    echo ""
    systemctl status servqr-frontend --no-pager -l | head -10
    echo ""
}

# Main deployment
main() {
    echo "=========================================================================="
    echo "  ServQR Platform - Application Deployment"
    echo "=========================================================================="
    echo ""
    
    check_prerequisites
    configure_environment
    generate_jwt_keys
    fix_database_schema
    build_backend
    build_frontend
    create_systemd_services
    start_services
    display_status
    
    # Get server IP
    SERVER_IP=$(hostname -I | awk '{print $1}')
    
    log ""
    log "✓ Application deployment completed successfully!"
    log ""
    log "Access Information:"
    log "  Frontend: http://${SERVER_IP}:3000"
    log "  Backend API: http://${SERVER_IP}:8081/api/v1"
    log "  Health Check: http://${SERVER_IP}:8081/health"
    log ""
    log "Default Login Credentials:"
    log "  Check: ${INSTALL_DIR}/LOGIN-CREDENTIALS.txt"
    log ""
    log "Service Management:"
    log "  Backend:"
    log "    - Status:  sudo systemctl status servqr-backend"
    log "    - Logs:    sudo journalctl -u servqr-backend -f"
    log "    - Restart: sudo systemctl restart servqr-backend"
    log ""
    log "  Frontend:"
    log "    - Status:  sudo systemctl status servqr-frontend"
    log "    - Logs:    sudo journalctl -u servqr-frontend -f"
    log "    - Restart: sudo systemctl restart servqr-frontend"
    log ""
    log "Configuration Files:"
    log "  Backend:  ${INSTALL_DIR}/.env"
    log "  Frontend: ${FRONTEND_DIR}/.env.local"
    log "  Secrets:  ${INSTALL_DIR}/.db_password, ${INSTALL_DIR}/.jwt_secret"
    log ""
    warn "IMPORTANT: Update API keys (AI, SendGrid, Twilio) in ${INSTALL_DIR}/.env"
    log ""
}

main "$@"
