# Documentation Cleanup Script
# Safely removes outdated documentation files

Write-Host "`n🧹 Starting Documentation Cleanup...`n" -ForegroundColor Cyan

# Create backup directory
$backupDir = "archived/old-docs-backup-$(Get-Date -Format 'yyyy-MM-dd-HHmmss')"
New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
Write-Host "✅ Created backup directory: $backupDir" -ForegroundColor Green

# Files to delete (will be moved to archive first)
$filesToArchive = @(
    "API-FIX-SUMMARY.md",
    "BACKEND-API-STATUS.md",
    "BACKEND-DEBUG-STATUS.md",
    "CODE-AUDIT-AND-IMPROVEMENTS.md",
    "COMPREHENSIVE-ARCHITECTURE-ANALYSIS.md",
    "DATABASE-SAMPLE-DATA.md",
    "DEMO-READY-STATUS.md",
    "FINAL-STATUS.md",
    "FRONTEND-DEBUG-INSTRUCTIONS.md",
    "IMPLEMENTATION-CHECKLIST.md",
    "IMPLEMENTATION-COMPLETE.md",
    "IMPLEMENTATION-GUIDE.md",
    "MANUFACTURERS-CLARIFICATION.md",
    "MOCK-DATA-AUDIT-AND-REDESIGN-PLAN.md",
    "PROGRESS-UPDATE.md",
    "QA-TESTING-SPECIFICATIONS.md",
    "QR-CODE-CONTENT-EXPLAINED.md",
    "QR-DATABASE-STORAGE-COMPLETE.md",
    "QR-GENERATION-FIX.md",
    "QR-SERVICE-REQUEST-FIX.md",
    "QR-SYSTEM-STATUS-FINAL.md",
    "QR-URL-FIX-COMPLETE.md",
    "QUICK-START.md",
    "REACT-QUERY-EXAMPLES.md",
    "SERVICE-SPECIFICATIONS-INDEX.md",
    "SERVICES-RUNNING.md",
    "SYSTEM-READY-FOR-TESTING.md",
    "DOCS-CLEANUP-PLAN.md"
)

# Log and build artifact files to delete permanently
$filesToDelete = @(
    "backend.log",
    "backend-error.log",
    "platform.log",
    "platform-stdout.log",
    "platform-stderr.log",
    "platform-err.log",
    "platform_runtime.log",
    "platform_runtime.err",
    "ui_dev.err",
    "ui_dev.out",
    "medical-platform.exe",
    "platform.exe"
)

# Temporary SQL files to archive
$sqlFilesToArchive = @(
    "add-remaining-tables.sql",
    "apply-qr-migration.sql",
    "fix-contract-comparison-schema.sql",
    "fix-database-schema.sql",
    "init-database-schema.sql"
)

# Test files to move to tests folder
$testFiles = @(
    "test-csv-import.ps1",
    "test-equipment-registration.ps1",
    "test-qr-eq-001.png",
    "manufacturer-installations-sample.csv",
    "equipment-ids.json"
)

Write-Host "`n📦 Archiving outdated documentation files..." -ForegroundColor Yellow
$archivedCount = 0
foreach ($file in $filesToArchive) {
    if (Test-Path $file) {
        Move-Item $file $backupDir -Force
        Write-Host "  ✓ Archived: $file" -ForegroundColor Gray
        $archivedCount++
    }
}

Write-Host "`n📦 Archiving old SQL files..." -ForegroundColor Yellow
foreach ($file in $sqlFilesToArchive) {
    if (Test-Path $file) {
        Move-Item $file $backupDir -Force
        Write-Host "  ✓ Archived: $file" -ForegroundColor Gray
        $archivedCount++
    }
}

Write-Host "`n🗑️  Deleting logs and build artifacts..." -ForegroundColor Yellow
$deletedCount = 0
foreach ($file in $filesToDelete) {
    if (Test-Path $file) {
        Remove-Item $file -Force
        Write-Host "  ✓ Deleted: $file" -ForegroundColor Gray
        $deletedCount++
    }
}

Write-Host "`n📁 Moving test files to tests/ folder..." -ForegroundColor Yellow
if (!(Test-Path "tests")) {
    New-Item -ItemType Directory -Path "tests" -Force | Out-Null
    New-Item -ItemType Directory -Path "tests\fixtures" -Force | Out-Null
}

$movedCount = 0
foreach ($file in $testFiles) {
    if (Test-Path $file) {
        if ($file -like "*.png" -or $file -like "*.csv" -or $file -like "*.json") {
            Move-Item $file "tests\fixtures\" -Force
        } else {
            Move-Item $file "tests\" -Force
        }
        Write-Host "  ✓ Moved: $file" -ForegroundColor Gray
        $movedCount++
    }
}

Write-Host "`n📂 Creating organized docs structure..." -ForegroundColor Yellow
$docsDirs = @(
    "docs\architecture",
    "docs\database",
    "docs\api"
)

foreach ($dir in $docsDirs) {
    if (!(Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force | Out-Null
        Write-Host "  ✓ Created: $dir" -ForegroundColor Gray
    }
}

# Move architecture docs
if (Test-Path "DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md") {
    Copy-Item "DETAILED-ORGANIZATIONS-ARCHITECTURE-DESIGN.md" "docs\architecture\organizations-architecture.md" -Force
    Write-Host "  ✓ Copied to docs\architecture\organizations-architecture.md" -ForegroundColor Gray
}

if (Test-Path "ENGINEER-MANAGEMENT-DESIGN.md") {
    Copy-Item "ENGINEER-MANAGEMENT-DESIGN.md" "docs\architecture\engineer-management.md" -Force
    Write-Host "  ✓ Copied to docs\architecture\engineer-management.md" -ForegroundColor Gray
}

if (Test-Path "IMPLEMENTATION-ROADMAP.md") {
    Copy-Item "IMPLEMENTATION-ROADMAP.md" "docs\architecture\implementation-roadmap.md" -Force
    Write-Host "  ✓ Copied to docs\architecture\implementation-roadmap.md" -ForegroundColor Gray
}

if (Test-Path "PHASE1-DATABASE-COMPLETE.md") {
    Copy-Item "PHASE1-DATABASE-COMPLETE.md" "docs\database\phase1-complete.md" -Force
    Write-Host "  ✓ Copied to docs\database\phase1-complete.md" -ForegroundColor Gray
}

Write-Host "`n📊 Cleanup Summary:" -ForegroundColor Cyan
Write-Host "  📦 Archived: $archivedCount files" -ForegroundColor Green
Write-Host "  🗑️  Deleted: $deletedCount files" -ForegroundColor Green
Write-Host "  📁 Moved: $movedCount test files" -ForegroundColor Green
Write-Host "  📂 Organized docs structure created" -ForegroundColor Green

Write-Host ""
Write-Host "✨ Cleanup Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "  Backup location: $backupDir" -ForegroundColor Yellow
Write-Host "  Next: Create README.md and ARCHITECTURE.md" -ForegroundColor Yellow
Write-Host ""
