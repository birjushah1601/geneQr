#!/bin/bash

################################################################################
# Fix Build Errors - Complete Fix for Backend Build Issues
#
# This script:
# 1. Verifies all required files exist
# 2. Fixes Go module issues
# 3. Creates missing init files if needed
# 4. Attempts to build
#
# Usage: sudo bash fix-build-errors.sh
################################################################################

set -e

echo "========================================================================"
echo "  ServQR - Fix Backend Build Errors"
echo "========================================================================"
echo ""

cd /opt/servqr

echo "Step 1: Checking required files..."
echo ""

# Check if init files exist
if [[ ! -f "cmd/platform/init_auth.go" ]]; then
    echo "WARNING: init_auth.go not found, creating it..."
    
    cat > cmd/platform/init_auth.go << 'EOFAUTH'
package main

import (
	"log/slog"
	"github.com/aby-med/medical-platform/internal/core/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// initAuthModule initializes the authentication module if enabled
// Returns nil module if auth is disabled
func initAuthModule(router *chi.Mux, db *sqlx.DB, logger *slog.Logger) (*auth.Module, error) {
	// Auth module is always enabled in production
	logger.Info("Initializing auth module")
	
	authModule := auth.NewModule(db, logger)
	if authModule == nil {
		logger.Warn("Auth module initialization returned nil")
		return nil, nil
	}
	
	logger.Info("Auth module initialized successfully")
	return authModule, nil
}
EOFAUTH
    echo "✓ Created init_auth.go"
else
    echo "✓ init_auth.go exists"
fi

if [[ ! -f "cmd/platform/init_notifications.go" ]]; then
    echo "WARNING: init_notifications.go not found, creating it..."
    
    cat > cmd/platform/init_notifications.go << 'EOFNOTIF'
package main

import (
	"context"
	"log/slog"
	"os"
	
	"github.com/aby-med/medical-platform/internal/infrastructure/notification"
	"github.com/aby-med/medical-platform/internal/infrastructure/reports"
	"github.com/jmoiron/sqlx"
)

// initNotificationsAndReports initializes email notifications and daily reports systems
func initNotificationsAndReports(ctx context.Context, db *sqlx.DB, logger *slog.Logger) (*notification.Manager, *reports.ReportScheduler, error) {
	// Check if email notifications are enabled
	emailEnabled := os.Getenv("FEATURE_EMAIL_NOTIFICATIONS") == "true"
	reportsEnabled := os.Getenv("FEATURE_DAILY_REPORTS") == "true"
	
	var notifMgr *notification.Manager
	var reportScheduler *reports.ReportScheduler
	
	if emailEnabled {
		logger.Info("Initializing email notification manager")
		sendgridKey := os.Getenv("SENDGRID_API_KEY")
		fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
		fromName := os.Getenv("SENDGRID_FROM_NAME")
		
		if sendgridKey == "" {
			logger.Warn("SENDGRID_API_KEY not set, notifications disabled")
		} else {
			notifMgr = notification.NewManager(sendgridKey, fromEmail, fromName, logger)
			logger.Info("Email notification manager initialized")
		}
	}
	
	if reportsEnabled && db != nil {
		logger.Info("Initializing daily reports scheduler")
		reportScheduler = reports.NewReportScheduler(db, logger)
		
		// Start the scheduler
		if err := reportScheduler.Start(); err != nil {
			logger.Error("Failed to start report scheduler", slog.String("error", err.Error()))
		} else {
			logger.Info("Daily reports scheduler started")
		}
	}
	
	return notifMgr, reportScheduler, nil
}
EOFNOTIF
    echo "✓ Created init_notifications.go"
else
    echo "✓ init_notifications.go exists"
fi

echo ""
echo "Step 2: Verifying file contents..."
ls -lh cmd/platform/init*.go
echo ""

echo "Step 3: Cleaning Go build cache..."
go clean -cache -modcache || true
echo ""

echo "Step 4: Downloading Go modules..."
go mod download
go mod tidy
echo ""

echo "Step 5: Attempting to build backend..."
echo ""

# Set Go environment
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export CGO_ENABLED=0

# Try building
if go build -v -o platform ./cmd/platform/; then
    echo ""
    echo "========================================================================"
    echo "  ✓ Backend built successfully!"
    echo "========================================================================"
    echo ""
    echo "Binary: $(pwd)/platform"
    echo "Size: $(du -h platform | cut -f1)"
    echo ""
    echo "Next step: sudo bash deployment/deploy-app.sh"
else
    echo ""
    echo "========================================================================"
    echo "  ✗ Build failed"
    echo "========================================================================"
    echo ""
    echo "Diagnostic info:"
    echo "1. Go version: $(go version)"
    echo "2. Files in cmd/platform:"
    ls -la cmd/platform/*.go
    echo ""
    echo "3. Go module info:"
    go list -m all | head -10
    echo ""
    echo "Please check the error messages above."
    exit 1
fi
