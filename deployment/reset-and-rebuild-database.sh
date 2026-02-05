#!/bin/bash

################################################################################
# ServQR Platform - Database Reset and Rebuild Script
#
# This script:
# 1. Stops backend and frontend services
# 2. Backs up existing database (if it exists)
# 3. Drops and recreates the database
# 4. Runs the base schema (001_full_organizations_schema.sql)
# 5. Runs additional migrations
# 6. Restarts services
#
# Usage: sudo bash deployment/reset-and-rebuild-database.sh
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

INSTALL_DIR="/opt/servqr"
DEPLOYMENT_DIR="${INSTALL_DIR}/deployment"
BACKUP_DIR="${INSTALL_DIR}/backups"
DB_NAME="servqr_production"
DB_USER="servqr"

# Check root
if [[ $EUID -ne 0 ]]; then
   error "This script must be run as root (use sudo)"
fi

# Check installation directory
if [[ ! -d "$INSTALL_DIR" ]]; then
    error "Installation directory not found: $INSTALL_DIR"
fi

echo "=========================================================================="
echo "  ServQR Platform - Database Reset and Rebuild"
echo "=========================================================================="
echo ""
warn "This will DROP and RECREATE the database!"
warn "All existing data will be LOST (a backup will be created)."
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [[ "$confirm" != "yes" ]]; then
    log "Operation cancelled by user"
    exit 0
fi

# Step 1: Stop services
log "Stopping services..."
systemctl stop servqr-backend servqr-frontend || warn "Services may not be running"
log "Services stopped"

# Step 2: Backup existing database
log "Creating database backup..."
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="${BACKUP_DIR}/pre-reset-backup-$(date +%Y%m%d-%H%M%S).sql"

if docker exec servqr-postgres psql -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    log "Database exists, creating backup..."
    docker exec servqr-postgres pg_dump -U $DB_USER $DB_NAME > "$BACKUP_FILE" 2>/dev/null || warn "Backup failed (database may be corrupted)"
    if [[ -f "$BACKUP_FILE" ]]; then
        gzip "$BACKUP_FILE"
        log "Backup created: ${BACKUP_FILE}.gz"
    fi
else
    log "No existing database found, skipping backup"
fi

# Step 3: Drop and recreate database
log "Dropping existing database..."
docker exec servqr-postgres psql -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true

log "Creating fresh database..."
docker exec servqr-postgres psql -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" || error "Failed to create database"

log "Database recreated successfully"

# Step 4: Run base schema
log "Running base schema (001_full_organizations_schema.sql)..."
BASE_SCHEMA="${INSTALL_DIR}/database/migrations/001_full_organizations_schema.sql"

if [[ ! -f "$BASE_SCHEMA" ]]; then
    error "Base schema not found: $BASE_SCHEMA"
fi

info "Applying base schema..."
if docker exec -i servqr-postgres psql -U $DB_USER -d $DB_NAME < "$BASE_SCHEMA"; then
    log "✓ Base schema applied successfully"
else
    error "Base schema failed to apply. Check the SQL syntax."
fi

# Verify tables were created
TABLE_COUNT=$(docker exec servqr-postgres psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
log "Tables created: $TABLE_COUNT"

if [[ $TABLE_COUNT -lt 10 ]]; then
    warn "Only $TABLE_COUNT tables created. Expected more. There may be issues."
fi

# Step 5: Run additional migrations
log "Running additional migrations..."
MIGRATIONS_DIR="${INSTALL_DIR}/database/migrations"

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
    warn "Migrations directory not found: $MIGRATIONS_DIR"
else
    # Copy migrations to container
    docker cp "$MIGRATIONS_DIR" servqr-postgres:/tmp/migrations
    
    # Run each migration (skip base schema)
    MIGRATION_COUNT=0
    SUCCESS_COUNT=0
    FAIL_COUNT=0
    
    for migration_file in $(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort); do
        migration_name=$(basename "$migration_file")
        
        # Skip base schema (already applied)
        if [[ "$migration_name" == "001_full_organizations_schema.sql" ]]; then
            log "Skipping $migration_name (already applied)"
            continue
        fi
        
        MIGRATION_COUNT=$((MIGRATION_COUNT + 1))
        info "Applying migration: $migration_name"
        
        if docker exec servqr-postgres psql -U $DB_USER -d $DB_NAME -f "/tmp/migrations/$migration_name" &>/dev/null; then
            SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
            echo "  ✓ Success"
        else
            FAIL_COUNT=$((FAIL_COUNT + 1))
            echo "  ✗ Failed (non-critical)"
        fi
    done
    
    log "Migration summary: $MIGRATION_COUNT total, $SUCCESS_COUNT succeeded, $FAIL_COUNT failed"
fi

# Step 6: Verify critical tables exist
log "Verifying critical tables..."
CRITICAL_TABLES=("organizations" "users" "service_tickets" "equipment" "engineers")
MISSING_TABLES=()

for table in "${CRITICAL_TABLES[@]}"; do
    if docker exec servqr-postgres psql -U $DB_USER -d $DB_NAME -t -c "SELECT 1 FROM information_schema.tables WHERE table_name = '$table';" | grep -q 1; then
        echo "  ✓ $table"
    else
        echo "  ✗ $table (MISSING)"
        MISSING_TABLES+=("$table")
    fi
done

if [[ ${#MISSING_TABLES[@]} -gt 0 ]]; then
    error "Critical tables missing: ${MISSING_TABLES[*]}"
fi

# Step 7: Display table list
log "Database schema overview:"
docker exec servqr-postgres psql -U $DB_USER -d $DB_NAME -c "\dt" | head -40

# Step 8: Restart services
log "Starting services..."
systemctl start servqr-backend
systemctl start servqr-frontend

# Wait for backend
log "Waiting for backend to start..."
for i in {1..30}; do
    if systemctl is-active --quiet servqr-backend; then
        log "✓ Backend service is active"
        break
    fi
    if [[ $i -eq 30 ]]; then
        warn "Backend service failed to start. Check logs: journalctl -u servqr-backend -n 50"
    fi
    sleep 2
done

# Wait for frontend
log "Waiting for frontend to start..."
sleep 5
if systemctl is-active --quiet servqr-frontend; then
    log "✓ Frontend service is active"
else
    warn "Frontend service failed to start. Check logs: journalctl -u servqr-frontend -n 50"
fi

echo ""
echo "=========================================================================="
echo "  Database Reset Complete"
echo "=========================================================================="
echo ""
log "Summary:"
log "  - Database dropped and recreated"
log "  - Base schema applied: $(basename $BASE_SCHEMA)"
log "  - Tables created: $TABLE_COUNT"
log "  - Migrations applied: $SUCCESS_COUNT succeeded, $FAIL_COUNT failed"
log "  - Backup saved: ${BACKUP_FILE}.gz (if existed)"
echo ""
log "Next Steps:"
log "  1. Check backend logs: journalctl -u servqr-backend -n 100"
log "  2. Check backend status: systemctl status servqr-backend"
log "  3. Test API: curl http://localhost:8081/health"
log "  4. Access frontend: http://$(hostname -I | awk '{print $1}'):3000"
echo ""
log "If backend still fails, check:"
log "  - Environment variables: cat /opt/servqr/.env"
log "  - Database connection: docker exec servqr-postgres psql -U servqr -d servqr_production -c 'SELECT version();'"
echo ""
