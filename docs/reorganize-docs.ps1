# Documentation Reorganization Script
# Moves old progress logs and summaries to archives/

Write-Host "`n╔═══════════════════════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                 📁 DOCUMENTATION REORGANIZATION SCRIPT 📁                     ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

$docsRoot = "C:\Users\birju\aby-med\docs"
$archivesPath = Join-Path $docsRoot "archives"

# Create archives directory if it doesn't exist
if (!(Test-Path $archivesPath)) {
    New-Item -ItemType Directory -Path $archivesPath -Force | Out-Null
    Write-Host "✅ Created archives directory" -ForegroundColor Green
}

# Files to move to archives (progress logs, summaries, old docs)
$filesToArchive = @(
    "ACTIVE_TICKETS_API_COMPLETE.md",
    "ALL_CARDS_UPDATED_SUMMARY.md",
    "API-TEST-RESULTS.md",
    "AUDIT-LOGGING-COMPLETE.md",
    "BACKEND_METADATA_FIX.md",
    "COMPLETE-SYSTEM-READY.md",
    "CURRENT-STATUS-DEC-2025.md",
    "DASHBOARD_ORGANIZATIONS_COUNT_FIX.md",
    "DEMO_EQUIPMENT_COMPLETE.md",
    "DEMO_QUICK_REFERENCE.md",
    "ENGINEER-ASSIGNMENT-BACKEND-COMPLETE.md",
    "ENGINEER-ASSIGNMENT-COMPLETE-WITH-POSTMAN.md",
    "ENGINEER-ASSIGNMENT-TESTED-WORKING.md",
    "ENGINEERS_ASSIGNED_TO_MANUFACTURERS.md",
    "ENGINEERS_COUNT_API_COMPLETE.md",
    "ENGINEERS_LEVEL_CARDS_FIX.md",
    "ENGINEERS_PAGE_FIX.md",
    "ENGINEER_REASSIGNMENT_FIX.md",
    "EQUIPMENT_COUNTS_API_FIX.md",
    "EQUIPMENT_MANUFACTURER_LINKING_COMPLETE.md",
    "EQUIPMENT_MANUFACTURER_MAPPING.md",
    "EQUIPMENT_PAGE_REFERENCE_ERROR_FIX.md",
    "EQUIPMENT_QR_CODE_FIXED.md",
    "EQUIPMENT_QR_TROUBLESHOOTING.md",
    "EQUIPMENT_REGISTRY_SCHEMA_FIX.md",
    "GIT-PUSH-SUMMARY.md",
    "MANUFACTURERS_API_URL_FIX.md",
    "MANUFACTURERS_COMPLETE_SETUP.md",
    "MANUFACTURERS_FILTER_FIX.md",
    "MANUFACTURERS_PAGE_FIX.md",
    "MANUFACTURER_DASHBOARD_FIX.md",
    "MANUFACTURER_DASHBOARD_REAL_DATA.md",
    "MANUFACTURER_DATA_COMPLETE.md",
    "MIGRATION-COMPLETE-TESTING-NEXT.md",
    "MULTI-TENANT-AUTH-SETUP.md",
    "MULTI-TENANT-IMPLEMENTATION-COMPLETE.md",
    "NAVIGATION-ALL-PAGES-COMPLETE.md",
    "NAVIGATION-FIX-SUMMARY.md",
    "NAVIGATION-IMPROVEMENTS-COMPLETE.md",
    "NOTIFICATIONS-COMPLETE-SUMMARY.md",
    "OPTION2-ENGINEER-UI-COMPLETE.md",
    "ORGANIZATIONS_API_FIX.md",
    "ORGANIZATIONS_FILTER_URL_PARAMS.md",
    "PARTS-MANAGEMENT-COMPLETE.md",
    "PARTS_ASSIGNMENT_FIX_COMPLETE.md",
    "PARTS_ASSIGNMENT_TWO_PAGES.md",
    "PARTS_PERSISTENCE_FIX.md",
    "PHASE-2-COMPLETE-SUMMARY.md",
    "PHASE1-COMPLETE.md",
    "PHASE1-IMPLEMENTATION-STARTED.md",
    "PHASE1-PROGRESS-SUMMARY.md",
    "PHASE2-ENGINEER-ASSIGNMENT-APIS-COMPLETE.md",
    "PHASE_2C_COMPLETE.md",
    "QR-RATE-LIMITING-COMPLETE.md",
    "QR_CODES_SERVICE_REQUEST_COMPLETE.md",
    "QR_CODE_DISPLAY_FIX.md",
    "QR_IMAGES_EXTERNAL_URL_FIX.md",
    "SESSION-DEC-12-2025-SUMMARY.md",
    "SESSION-NOTIFICATIONS-AND-REPORTS-SUMMARY.md",
    "SESSION_COMPLETE_SUMMARY.md",
    "SPARE_PARTS_CATALOG_COMPLETE.md",
    "SPARE_PARTS_FLOW_COMPLETE.md",
    "SPARE_PARTS_WITH_IMAGES.md",
    "TESTING-VERIFICATION.md",
    "TICKETS-PARTS-INTEGRATION-COMPLETE.md",
    "TROUBLESHOOTING_PARTS_ASSIGNMENT.md",
    "UI-WHATSAPP-INTEGRATION-COMPLETE.md",
    "WEEK1-DAY2-INTEGRATION-COMPLETE.md",
    "WEEK1-DAY3-FRONTEND-INTEGRATION.md",
    "WEEK1-DAY4-5-PRODUCTION-READY.md",
    "WEEK1-IMPLEMENTATION-GUIDE.md",
    "WEEK2-DASHBOARD-STATUS.md",
    "WEEK3-ENGINEER-ASSIGNMENT-COMPLETE.md",
    "MASTER-DOCUMENTATION-INDEX.md",
    "MASTER-DOCUMENTATION-INDEX-V2.md",
    "admin-ui-quick-start.md",
    "dev-setup.md"
)

$movedCount = 0
foreach ($file in $filesToArchive) {
    $sourcePath = Join-Path $docsRoot $file
    if (Test-Path $sourcePath) {
        $destPath = Join-Path $archivesPath $file
        Move-Item -Path $sourcePath -Destination $destPath -Force
        $movedCount++
    }
}

Write-Host "✅ Moved $movedCount files to archives/" -ForegroundColor Green
Write-Host ""

# Keep these important files in root
$keepFiles = @(
    "README.md",
    "01-GETTING-STARTED.md",
    "02-ARCHITECTURE.md",
    "03-FEATURES.md",
    "04-API-REFERENCE.md",
    "05-DEPLOYMENT.md",
    "06-PERSONAS.md",
    "MARKETPLACE-BRAINSTORMING.md",
    "TICKET-ENHANCEMENTS-IMPLEMENTATION.md",
    "ONBOARDING-SYSTEM-BRAINSTORM.md",
    "ONBOARDING-SYSTEM-README.md",
    "ONBOARDING-IMPLEMENTATION-ROADMAP.md",
    "MANUFACTURER-ONBOARDING-UX-DESIGN.md",
    "QR-CODE-TABLE-DESIGN-ANALYSIS.md",
    "WEEK-1-PROGRESS.md",
    "DEPLOYMENT-GUIDE.md",
    "EXECUTIVE-SUMMARY.md",
    "QUICK-REFERENCE.md",
    "AUTHENTICATION-MULTITENANCY-PRD.md",
    "MULTI-TENANT-IMPLEMENTATION-PLAN.md",
    "QR-CODE-PUBLIC-ACCESS-ANALYSIS.md",
    "SECURITY-IMPLEMENTATION-COMPLETE.md",
    "PRODUCTION-DEPLOYMENT-CHECKLIST.md",
    "EXTERNAL-SERVICES-SETUP.md",
    "LOGIN-PASSWORD-DEFAULT.md",
    "COMMIT-INSTRUCTIONS.md",
    "EMAIL-NOTIFICATIONS-SYSTEM.md",
    "DAILY-REPORTS-SYSTEM.md",
    "FEATURE-FLAGS-NOTIFICATIONS.md",
    "SECURITY-ROADMAP.md",
    "TESTING-GUIDE-MULTI-TENANT.md",
    "TICKET-CREATION-SECURITY-ASSESSMENT.md",
    "AI_ASSISTED_IMPLEMENTATION.md",
    "AI_INTEGRATION_STATUS.md",
    "CLIENT_CAPABILITIES.md",
    "EQUIPMENT_AND_PARTS_SYSTEM.md",
    "ER_DIAGRAM.md",
    "FEEDBACK_SYSTEM.md",
    "INTEGRATION_PLAN.md",
    "MANUFACTURER_ONBOARDING.md",
    "REQUIREMENTS_MASTER.md",
    "SIMPLIFIED-ENGINEER-ASSIGNMENT-IMPLEMENTATION.md",
    "STRATEGIC-IMPLEMENTATION-PIPELINE.md",
    "TESTING.md",
    "OPTION3-WHATSAPP-IMPLEMENTATION-GUIDE.md",
    "AUTHENTICATION-READY-TO-DEPLOY.md"
)

Write-Host "✅ Keeping important documentation in root" -ForegroundColor Green
Write-Host ""

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "📊 REORGANIZATION SUMMARY:" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Archived: $movedCount progress logs and summaries" -ForegroundColor White
Write-Host "  Kept: Main documentation and specifications" -ForegroundColor White
Write-Host ""
Write-Host "New Structure:" -ForegroundColor Green
Write-Host "  ├── README.md (navigation hub)" -ForegroundColor White
Write-Host "  ├── 01-GETTING-STARTED.md" -ForegroundColor White
Write-Host "  ├── 02-ARCHITECTURE.md (coming next)" -ForegroundColor White
Write-Host "  ├── 03-FEATURES.md (coming next)" -ForegroundColor White
Write-Host "  ├── 04-API-REFERENCE.md (coming next)" -ForegroundColor White
Write-Host "  ├── 05-DEPLOYMENT.md (coming next)" -ForegroundColor White
Write-Host "  ├── 06-PERSONAS.md (coming next)" -ForegroundColor White
Write-Host "  ├── Feature-specific docs (marketplace, onboarding, etc.)" -ForegroundColor White
Write-Host "  └── archives/ (old progress logs)" -ForegroundColor Gray
Write-Host ""
Write-Host "✅ Documentation cleanup complete!" -ForegroundColor Green
Write-Host ""
