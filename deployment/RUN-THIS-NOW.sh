#!/bin/bash
################################################################################
# CRITICAL FIX - Run this on your server NOW
################################################################################

echo "=========================================="
echo "  ServQR - Applying Critical Database Fix"
echo "=========================================="
echo ""

# Step 1: Get latest fixes
echo "Step 1/3: Pulling latest fixes from GitHub..."
cd /opt/servqr
git pull origin main

if [ $? -ne 0 ]; then
    echo "ERROR: Failed to pull latest code"
    exit 1
fi

echo "✓ Latest code pulled"
echo ""

# Step 2: Reset database
echo "Step 2/3: Resetting database with corrected schema..."
echo "(This will backup existing data first)"
echo ""

sudo bash deployment/reset-and-rebuild-database.sh << EOF
yes
EOF

if [ $? -ne 0 ]; then
    echo "ERROR: Database reset failed"
    exit 1
fi

echo ""
echo "✓ Database reset complete"
echo ""

# Step 3: Verify
echo "Step 3/3: Verifying deployment..."
echo ""

# Check services
echo "Backend status:"
sudo systemctl status servqr-backend --no-pager | head -5

echo ""
echo "Frontend status:"
sudo systemctl status servqr-frontend --no-pager | head -5

echo ""
echo "Table count:"
docker exec servqr-postgres psql -U servqr -d servqr_production -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"

echo ""
echo "=========================================="
echo "  Fix Applied!"
echo "=========================================="
echo ""
echo "Test URLs:"
echo "  Backend:  http://$(hostname -I | awk '{print $1}'):8081/health"
echo "  Frontend: http://$(hostname -I | awk '{print $1}'):3000"
echo ""
echo "Login:"
echo "  Email: admin@geneqr.com"
echo "  Password: admin123"
echo ""
