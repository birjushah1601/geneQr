#!/bin/bash
# Restore Complete Database (Schema + Data) from SQL Dump
# This script drops the existing database and restores from a full pg_dump export

set -e

BACKUP_FILE="${1:-/tmp/local_schema_export.sql}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
LOG_FILE="/opt/servqr/logs/restore-$TIMESTAMP.log"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log() {
    echo -e "${CYAN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

# Header
echo ""
log "=== Database Restore (Schema + Data) ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    error "Please run as root (use sudo)"
    exit 1
fi

# Check if file exists
if [ ! -f "$BACKUP_FILE" ]; then
    error "Backup file not found: $BACKUP_FILE"
    echo ""
    echo "Usage: sudo bash $0 [path-to-dump-file.sql]"
    echo "Example: sudo bash $0 /tmp/local_schema_export.sql"
    exit 1
fi

log "Backup file: $BACKUP_FILE"
log "Size: $(du -h $BACKUP_FILE | cut -f1)"
echo ""

# Check if Docker container is running
if ! docker ps | grep -q servqr-postgres; then
    error "Docker container 'servqr-postgres' is not running"
    echo ""
    echo "Start it with:"
    echo "  cd /opt/servqr/deployment && bash setup-docker.sh"
    exit 1
fi

# Confirmation
warn "⚠️  This will:"
echo "   1. Create a backup of current database"
echo "   2. DROP the entire 'medical_equipment' database"
echo "   3. CREATE a new 'medical_equipment' database"
echo "   4. RESTORE all schema and data from: $BACKUP_FILE"
echo ""
warn "All existing data will be REPLACED!"
echo ""
read -p "Do you want to continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    log "Restore cancelled by user"
    exit 0
fi

echo ""
log "Step 1/5: Creating safety backup of current database..."

# Create backup directory
mkdir -p /opt/servqr/backups

# Backup current database
SAFETY_BACKUP="/opt/servqr/backups/pre-restore-backup-$TIMESTAMP.sql"
docker exec servqr-postgres pg_dump -U servqr -d medical_equipment > "$SAFETY_BACKUP" 2>/dev/null || {
    warn "Could not backup current database (might be empty or corrupted)"
}

if [ -f "$SAFETY_BACKUP" ] && [ -s "$SAFETY_BACKUP" ]; then
    success "Safety backup saved: $SAFETY_BACKUP"
else
    warn "No safety backup created (database might not exist yet)"
fi

echo ""
log "Step 2/5: Terminating active connections..."

# Terminate all connections to the database
docker exec servqr-postgres psql -U servqr -d postgres -c "
    SELECT pg_terminate_backend(pid) 
    FROM pg_stat_activity 
    WHERE datname = 'medical_equipment' AND pid <> pg_backend_pid();
" >> "$LOG_FILE" 2>&1 || warn "No active connections to terminate"

sleep 2

echo ""
log "Step 3/5: Dropping existing database..."

# Drop existing database
docker exec servqr-postgres psql -U servqr -d postgres -c "DROP DATABASE IF EXISTS medical_equipment;" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    success "Database dropped successfully"
else
    error "Failed to drop database"
    exit 1
fi

echo ""
log "Step 4/5: Creating fresh database..."

# Create new database
docker exec servqr-postgres psql -U servqr -d postgres -c "CREATE DATABASE medical_equipment OWNER servqr;" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    success "Database created successfully"
else
    error "Failed to create database"
    exit 1
fi

echo ""
log "Step 5/5: Restoring schema and data from backup..."
log "This may take a few minutes..."

# Restore from backup file
docker exec -i servqr-postgres psql -U servqr -d medical_equipment < "$BACKUP_FILE" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    echo ""
    success "Database restored successfully!"
    echo ""
    
    # Show statistics
    log "Verifying database..."
    echo ""
    
    echo "=== Tables Created ==="
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "\dt" 2>/dev/null | grep "public" || echo "No tables found"
    
    echo ""
    echo "=== Record Counts ==="
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "
        SELECT 'organizations' as table_name, COUNT(*) as records FROM organizations
        UNION ALL SELECT 'users', COUNT(*) FROM users
        UNION ALL SELECT 'equipment', COUNT(*) FROM equipment
        UNION ALL SELECT 'service_tickets', COUNT(*) FROM service_tickets
        UNION ALL SELECT 'engineers', COUNT(*) FROM engineers
        UNION ALL SELECT 'spare_parts', COUNT(*) FROM spare_parts
        ORDER BY table_name;
    " 2>/dev/null || warn "Could not get record counts (tables might not exist)"
    
    echo ""
    log "Restarting backend service..."
    systemctl restart servqr-backend
    
    sleep 3
    
    if systemctl is-active --quiet servqr-backend; then
        success "Backend service restarted successfully"
    else
        warn "Backend service may have issues. Check logs:"
        echo "  sudo journalctl -u servqr-backend -n 50 --no-pager"
    fi
    
    echo ""
    log "Restarting frontend service..."
    systemctl restart servqr-frontend
    
    echo ""
    success "=== Restore Complete ==="
    echo ""
    log "Log file: $LOG_FILE"
    log "Safety backup: $SAFETY_BACKUP"
    echo ""
    log "You can now access the application with your local data!"
    
else
    echo ""
    error "Database restore FAILED!"
    echo ""
    warn "Check the log file for details:"
    echo "  cat $LOG_FILE"
    echo ""
    
    if [ -f "$SAFETY_BACKUP" ] && [ -s "$SAFETY_BACKUP" ]; then
        warn "To restore from safety backup:"
        echo "  sudo bash $0 $SAFETY_BACKUP"
    fi
    
    exit 1
fi
