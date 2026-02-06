#!/bin/bash
# Import Local Data to Production Database
# Run this script on the production server after uploading the data export

set -e

LOG_FILE="/opt/servqr/logs/data-import-$(date +%Y%m%d-%H%M%S).log"
DATA_FILE="${1:-/tmp/local_data_export.sql}"

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
log "=== Production Data Import ==="
echo ""

# Check if data file exists
if [ ! -f "$DATA_FILE" ]; then
    error "Data file not found: $DATA_FILE"
    echo ""
    echo "Usage: sudo bash $0 [path-to-data-file.sql]"
    echo "Example: sudo bash $0 /tmp/local_data_export.sql"
    exit 1
fi

log "Data file: $DATA_FILE"
log "Size: $(du -h $DATA_FILE | cut -f1)"
echo ""

# Confirmation
warn "This will import data into the production database."
warn "Existing data may be affected if there are conflicts."
echo ""
read -p "Do you want to continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    log "Import cancelled by user"
    exit 0
fi

echo ""
log "Creating backup of current database..."

# Backup current database
BACKUP_FILE="/opt/servqr/backups/pre-import-backup-$(date +%Y%m%d-%H%M%S).sql"
mkdir -p /opt/servqr/backups

docker exec servqr-postgres pg_dump -U servqr -d medical_equipment > "$BACKUP_FILE"
if [ $? -eq 0 ]; then
    success "Backup saved: $BACKUP_FILE"
else
    error "Backup failed! Aborting import."
    exit 1
fi

echo ""
log "Importing data..."

# Import data
docker exec -i servqr-postgres psql -U servqr -d medical_equipment < "$DATA_FILE" >> "$LOG_FILE" 2>&1

if [ $? -eq 0 ]; then
    echo ""
    success "Data import completed successfully!"
    echo ""
    
    # Show statistics
    log "Database Statistics:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "
        SELECT 
            schemaname,
            tablename,
            n_tup_ins as inserted_rows,
            n_tup_upd as updated_rows,
            n_tup_del as deleted_rows
        FROM pg_stat_user_tables
        WHERE schemaname = 'public'
        ORDER BY n_tup_ins DESC
        LIMIT 20;
    "
    
    echo ""
    log "Verifying critical tables..."
    
    # Count records in critical tables
    echo ""
    echo "Organizations:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "SELECT COUNT(*) FROM organizations;"
    
    echo "Users:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "SELECT COUNT(*) FROM users;"
    
    echo "Equipment:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "SELECT COUNT(*) FROM equipment;"
    
    echo "Service Tickets:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "SELECT COUNT(*) FROM service_tickets;"
    
    echo "Engineers:"
    docker exec servqr-postgres psql -U servqr -d medical_equipment -c "SELECT COUNT(*) FROM engineers;"
    
    echo ""
    success "Import verification complete!"
    echo ""
    log "Restarting backend service..."
    systemctl restart servqr-backend
    
    sleep 3
    
    if systemctl is-active --quiet servqr-backend; then
        success "Backend service restarted successfully!"
    else
        warn "Backend service may have issues. Check logs:"
        echo "  sudo journalctl -u servqr-backend -n 50"
    fi
    
    echo ""
    success "Data import completed!"
    echo ""
    log "Log file: $LOG_FILE"
    log "Backup file: $BACKUP_FILE"
    
else
    echo ""
    error "Data import failed!"
    echo ""
    warn "To restore from backup:"
    echo "  docker exec -i servqr-postgres psql -U servqr -d medical_equipment < $BACKUP_FILE"
    echo ""
    log "Check log file for details: $LOG_FILE"
    exit 1
fi
