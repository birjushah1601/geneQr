#!/bin/bash

################################################################################
# ServQR Platform - Docker & PostgreSQL Setup
#
# Installs Docker, creates PostgreSQL container with persistent storage
#
# Usage: sudo bash setup-docker.sh
################################################################################

set -e
set -u

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

INSTALL_DIR="/opt/servqr"
DATA_DIR="${INSTALL_DIR}/data/postgres"
DEPLOYMENT_DIR="${INSTALL_DIR}/deployment"

# Detect OS
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    else
        error "Cannot detect operating system"
    fi
    
    log "Detected OS: $OS $VER"
}

# Install Docker
install_docker() {
    log "Installing Docker..."
    
    # Check if Docker is already installed
    if command -v docker &> /dev/null; then
        log "Docker already installed: $(docker --version)"
        return 0
    fi
    
    if [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        # Install Docker on Ubuntu/Debian
        apt-get update
        apt-get install -y ca-certificates curl gnupg
        
        # Add Docker's official GPG key
        install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/${OS}/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        chmod a+r /etc/apt/keyrings/docker.gpg
        
        # Add Docker repository
        echo \
          "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/${OS} \
          $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
          tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        # Install Docker
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    elif [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        # Install Docker on CentOS/RHEL/Fedora
        yum install -y yum-utils
        yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
        yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        
    else
        error "Unsupported OS for Docker installation: $OS"
    fi
    
    # Start Docker service
    systemctl start docker
    systemctl enable docker
    
    # Verify installation
    docker --version || error "Docker installation failed"
    
    log "Docker installed successfully: $(docker --version)"
}

# Create Docker Compose file
create_docker_compose() {
    log "Creating Docker Compose configuration..."
    
    # Generate random database password if not exists
    if [[ ! -f "${INSTALL_DIR}/.db_password" ]]; then
        DB_PASSWORD=$(openssl rand -base64 24 | tr -d "=+/" | cut -c1-24)
        echo "$DB_PASSWORD" > "${INSTALL_DIR}/.db_password"
        chmod 600 "${INSTALL_DIR}/.db_password"
    else
        DB_PASSWORD=$(cat "${INSTALL_DIR}/.db_password")
    fi
    
    # Create docker-compose.yml
    cat > "${DEPLOYMENT_DIR}/docker-compose.yml" << EOF
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: servqr-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: servqr_production
      POSTGRES_USER: servqr
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --locale=en_US.UTF-8"
      PGDATA: /var/lib/postgresql/data/pgdata
    ports:
      - "5432:5432"
    volumes:
      - ${DATA_DIR}:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U servqr -d servqr_production"]
      interval: 10s
      timeout: 5s
      retries: 5
    command:
      - "postgres"
      - "-c"
      - "max_connections=100"
      - "-c"
      - "shared_buffers=256MB"
      - "-c"
      - "effective_cache_size=1GB"
      - "-c"
      - "maintenance_work_mem=64MB"
      - "-c"
      - "checkpoint_completion_target=0.9"
      - "-c"
      - "wal_buffers=16MB"
      - "-c"
      - "default_statistics_target=100"
      - "-c"
      - "random_page_cost=1.1"
      - "-c"
      - "effective_io_concurrency=200"
      - "-c"
      - "work_mem=2621kB"
      - "-c"
      - "min_wal_size=1GB"
      - "-c"
      - "max_wal_size=4GB"
      - "-c"
      - "log_statement=all"
      - "-c"
      - "log_duration=on"
    networks:
      - servqr-network

networks:
  servqr-network:
    driver: bridge
EOF
    
    log "Docker Compose configuration created"
}

# Start PostgreSQL container
start_postgres() {
    log "Starting PostgreSQL container..."
    
    cd "${DEPLOYMENT_DIR}"
    docker compose up -d postgres
    
    # Wait for PostgreSQL to be ready
    log "Waiting for PostgreSQL to be ready..."
    max_retries=30
    retry_count=0
    
    while ! docker exec servqr-postgres pg_isready -U servqr -d servqr_production &> /dev/null; do
        retry_count=$((retry_count + 1))
        if [[ $retry_count -ge $max_retries ]]; then
            error "PostgreSQL failed to start after ${max_retries} retries"
        fi
        log "Waiting for PostgreSQL... (${retry_count}/${max_retries})"
        sleep 2
    done
    
    log "PostgreSQL is ready"
}

# Initialize database
initialize_database() {
    log "Initializing database schema..."
    
    # Check if migrations directory exists
    MIGRATIONS_DIR="${INSTALL_DIR}/migrations"
    if [[ ! -d "$MIGRATIONS_DIR" ]]; then
        MIGRATIONS_DIR="${INSTALL_DIR}/database/migrations"
    fi
    
    if [[ ! -d "$MIGRATIONS_DIR" ]]; then
        warn "Migrations directory not found. Skipping database initialization."
        warn "You'll need to run migrations manually later."
        return 0
    fi
    
    # Run migrations
    log "Running database migrations..."
    
    # Copy migration files to container
    docker cp "$MIGRATIONS_DIR" servqr-postgres:/tmp/migrations
    
    # Execute migrations
    for migration in $(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort); do
        migration_file=$(basename "$migration")
        log "Applying migration: $migration_file"
        docker exec servqr-postgres psql -U servqr -d servqr_production -f "/tmp/migrations/$migration_file" || warn "Migration $migration_file failed"
    done
    
    log "Database initialization completed"
}

# Create database utilities
create_db_utilities() {
    log "Creating database utility scripts..."
    
    # Database backup script
    cat > "${DEPLOYMENT_DIR}/backup-database.sh" << 'EOFBACKUP'
#!/bin/bash
set -e

BACKUP_DIR="/opt/servqr/backups"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/servqr-backup-${TIMESTAMP}.sql"

mkdir -p "$BACKUP_DIR"

echo "Creating database backup..."
docker exec servqr-postgres pg_dump -U servqr servqr_production > "$BACKUP_FILE"

echo "Compressing backup..."
gzip "$BACKUP_FILE"

echo "Backup completed: ${BACKUP_FILE}.gz"
echo "Backup size: $(du -h ${BACKUP_FILE}.gz | cut -f1)"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "servqr-backup-*.sql.gz" -mtime +7 -delete
echo "Old backups cleaned up (kept last 7 days)"
EOFBACKUP
    
    chmod +x "${DEPLOYMENT_DIR}/backup-database.sh"
    
    # Database restore script
    cat > "${DEPLOYMENT_DIR}/restore-database.sh" << 'EOFRESTORE'
#!/bin/bash
set -e

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <backup-file.sql.gz>"
    echo "Example: $0 /opt/servqr/backups/servqr-backup-20260206-120000.sql.gz"
    exit 1
fi

BACKUP_FILE="$1"

if [[ ! -f "$BACKUP_FILE" ]]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "WARNING: This will replace the current database!"
echo "Backup file: $BACKUP_FILE"
read -p "Are you sure you want to continue? (yes/no): " confirm

if [[ "$confirm" != "yes" ]]; then
    echo "Restore cancelled"
    exit 0
fi

echo "Decompressing backup..."
gunzip -c "$BACKUP_FILE" > /tmp/restore.sql

echo "Restoring database..."
docker exec -i servqr-postgres psql -U servqr -d servqr_production < /tmp/restore.sql

echo "Cleaning up..."
rm -f /tmp/restore.sql

echo "Database restored successfully from: $BACKUP_FILE"
EOFRESTORE
    
    chmod +x "${DEPLOYMENT_DIR}/restore-database.sh"
    
    # Database connection script
    cat > "${DEPLOYMENT_DIR}/connect-database.sh" << 'EOFCONNECT'
#!/bin/bash
echo "Connecting to PostgreSQL..."
docker exec -it servqr-postgres psql -U servqr -d servqr_production
EOFCONNECT
    
    chmod +x "${DEPLOYMENT_DIR}/connect-database.sh"
    
    log "Database utility scripts created"
}

# Configure Docker auto-start
configure_docker_autostart() {
    log "Configuring Docker container auto-start..."
    
    # Create systemd service for Docker Compose
    cat > /etc/systemd/system/servqr-postgres.service << EOF
[Unit]
Description=ServQR PostgreSQL Database
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=${DEPLOYMENT_DIR}
ExecStart=/usr/bin/docker compose up -d postgres
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable servqr-postgres
    
    log "Docker auto-start configured"
}

# Main setup
main() {
    echo "=========================================================================="
    echo "  ServQR Platform - Docker & PostgreSQL Setup"
    echo "=========================================================================="
    echo ""
    
    # Check root
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use sudo)"
    fi
    
    # Check installation directory
    if [[ ! -d "$INSTALL_DIR" ]]; then
        error "Installation directory not found: $INSTALL_DIR"
    fi
    
    # Detect OS
    detect_os
    
    # Create data directory
    log "Creating data directory: $DATA_DIR"
    mkdir -p "$DATA_DIR"
    
    # Install Docker
    install_docker
    
    # Create Docker Compose configuration
    create_docker_compose
    
    # Start PostgreSQL
    start_postgres
    
    # Initialize database
    initialize_database
    
    # Create utility scripts
    create_db_utilities
    
    # Configure auto-start
    configure_docker_autostart
    
    # Display information
    DB_PASSWORD=$(cat "${INSTALL_DIR}/.db_password")
    
    log ""
    log "✓ Docker and PostgreSQL setup completed successfully!"
    log ""
    log "PostgreSQL Information:"
    log "  - Container: servqr-postgres"
    log "  - Database: servqr_production"
    log "  - User: servqr"
    log "  - Password: $DB_PASSWORD"
    log "  - Port: 5432"
    log "  - Data Directory: $DATA_DIR"
    log ""
    log "Useful Commands:"
    log "  - Connect to database:    ${DEPLOYMENT_DIR}/connect-database.sh"
    log "  - Backup database:        ${DEPLOYMENT_DIR}/backup-database.sh"
    log "  - Restore database:       ${DEPLOYMENT_DIR}/restore-database.sh <backup-file>"
    log "  - View logs:              docker logs servqr-postgres"
    log "  - Restart container:      docker restart servqr-postgres"
    log ""
}

main "$@"
