#!/bin/bash

################################################################################
# Fix Database Migrations - Run this on the server
#
# This script:
# 1. Stops the deployment
# 2. Recreates the database with the correct init schema
# 3. Points migrations to the correct directory
#
# Usage: sudo bash fix-migrations.sh
################################################################################

set -e

echo "========================================================================"
echo "  ServQR - Fix Database Migrations"
echo "========================================================================"
echo ""

# Check root
if [[ $EUID -ne 0 ]]; then
    echo "ERROR: This script must be run as root"
    exit 1
fi

INSTALL_DIR="/opt/servqr"
DB_PASSWORD=$(cat "${INSTALL_DIR}/.db_password" 2>/dev/null || echo "")

if [[ -z "$DB_PASSWORD" ]]; then
    echo "ERROR: Database password not found"
    exit 1
fi

echo "Step 1: Stopping PostgreSQL container..."
docker stop servqr-postgres || true
docker rm servqr-postgres || true

echo ""
echo "Step 2: Removing old database data..."
rm -rf /opt/servqr/data/postgres/*

echo ""
echo "Step 3: Starting fresh PostgreSQL container..."
cd /opt/servqr/deployment
docker compose up -d postgres

echo ""
echo "Step 4: Waiting for PostgreSQL to be ready..."
sleep 5

max_retries=30
retry_count=0

while ! docker exec servqr-postgres pg_isready -U servqr -d servqr_production &> /dev/null; do
    retry_count=$((retry_count + 1))
    if [[ $retry_count -ge $max_retries ]]; then
        echo "ERROR: PostgreSQL failed to start"
        exit 1
    fi
    echo "Waiting... (${retry_count}/${max_retries})"
    sleep 2
done

echo "PostgreSQL is ready!"
echo ""

echo "Step 5: Running base schema initialization..."

# Check if init-database-schema.sql exists
if [[ -f "${INSTALL_DIR}/init-database-schema.sql" ]]; then
    echo "Using init-database-schema.sql..."
    docker exec -i servqr-postgres psql -U servqr -d servqr_production < "${INSTALL_DIR}/init-database-schema.sql"
elif [[ -f "${INSTALL_DIR}/database/migrations/001_full_organizations_schema.sql" ]]; then
    echo "Using 001_full_organizations_schema.sql..."
    docker exec -i servqr-postgres psql -U servqr -d servqr_production < "${INSTALL_DIR}/database/migrations/001_full_organizations_schema.sql"
else
    echo "WARNING: No base schema found, will try individual migrations"
fi

echo ""
echo "Step 6: Running migrations from database/migrations..."

# Copy correct migrations directory
if [[ -d "${INSTALL_DIR}/database/migrations" ]]; then
    docker cp "${INSTALL_DIR}/database/migrations" servqr-postgres:/tmp/db_migrations
    
    # Run migrations in order
    for migration in $(docker exec servqr-postgres ls -1 /tmp/db_migrations/*.sql | sort); do
        migration_file=$(basename "$migration")
        
        # Skip demo data and test files
        if [[ "$migration_file" == *"demo"* ]] || [[ "$migration_file" == *"test"* ]]; then
            echo "Skipping: $migration_file"
            continue
        fi
        
        echo "Applying: $migration_file"
        docker exec servqr-postgres psql -U servqr -d servqr_production -f "/tmp/db_migrations/$migration_file" 2>&1 | grep -v "NOTICE:" | grep -v "already exists" || true
    done
fi

echo ""
echo "Step 7: Verifying database setup..."
docker exec servqr-postgres psql -U servqr -d servqr_production -c "\dt" | head -20

echo ""
echo "========================================================================"
echo "  Database fixed and ready!"
echo "========================================================================"
echo ""
echo "Next steps:"
echo "  1. Continue with deployment: cd /opt/servqr/deployment && sudo bash deploy-app.sh"
echo "  2. Or run full deployment again: sudo bash deploy-all.sh"
echo ""
