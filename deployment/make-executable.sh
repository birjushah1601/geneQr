#!/bin/bash
# Make all deployment scripts executable

cd "$(dirname "$0")"

echo "Making deployment scripts executable..."

chmod +x deploy-all.sh
chmod +x install-prerequisites.sh
chmod +x setup-docker.sh
chmod +x deploy-app.sh
chmod +x backup-database.sh 2>/dev/null || true
chmod +x restore-database.sh 2>/dev/null || true
chmod +x connect-database.sh 2>/dev/null || true

echo "✓ All deployment scripts are now executable"
echo ""
echo "You can now run:"
echo "  sudo bash deploy-all.sh"
