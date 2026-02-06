#!/bin/bash
# Apply quick schema fixes to make production work right now

echo "=== Applying Quick Schema Fixes ==="
echo ""

cd /opt/servqr

# Apply fixes
docker exec -i servqr-postgres psql -U servqr -d medical_equipment < deployment/fix-production-schema-now.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Schema fixes applied!"
    echo ""
    echo "Restarting backend..."
    systemctl restart servqr-backend
    
    sleep 3
    
    echo ""
    echo "Testing endpoints..."
    echo "Engineers:"
    curl -s http://localhost:8081/api/v1/engineers?page_size=1 | jq '.total' || echo "  (auth required)"
    
    echo "Tickets:"
    curl -s http://localhost:8081/api/v1/tickets?page_size=1 | jq '.total' || echo "  (auth required)"
    
    echo "Equipment:"
    curl -s http://localhost:8081/api/v1/equipment?page_size=1 | jq '.total' || echo "  (auth required)"
    
    echo ""
    echo "✅ Production should now work!"
else
    echo "❌ Failed to apply fixes"
    exit 1
fi
