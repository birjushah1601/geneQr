#!/bin/bash

################################################################################
# Quick Fix and Deploy Script
# 
# This script fixes the logs directory issue and runs deployment
#
# Usage: sudo bash fix-and-deploy.sh
################################################################################

set -e

echo "=========================================================================="
echo "  ServQR Platform - Fix and Deploy"
echo "=========================================================================="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo "ERROR: This script must be run as root (use sudo)"
    exit 1
fi

# Create required directories first
echo "Creating required directories..."
mkdir -p /opt/servqr/logs
mkdir -p /opt/servqr/data/postgres
mkdir -p /opt/servqr/data/qrcodes
mkdir -p /opt/servqr/data/whatsapp
mkdir -p /opt/servqr/backups
mkdir -p /opt/servqr/storage

echo "Directories created successfully"
echo ""

# Make scripts executable
echo "Making deployment scripts executable..."
cd /opt/servqr/deployment
chmod +x deploy-all.sh
chmod +x install-prerequisites.sh
chmod +x setup-docker.sh
chmod +x deploy-app.sh
chmod +x make-executable.sh

echo "Scripts are now executable"
echo ""

# Run deployment
echo "Starting deployment..."
echo ""
bash /opt/servqr/deployment/deploy-all.sh

exit $?
