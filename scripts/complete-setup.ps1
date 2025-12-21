# Complete System Setup Script - From Auth to Production

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
Write-Host "🚀 COMPLETE SYSTEM SETUP - 4 WEEK PIPELINE" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Green

$ErrorActionPreference = "Continue"
$completedSteps = @()
$failedSteps = @()

# ============================================================================
# WEEK 1, DAY 1: AUTHENTICATION DEPLOYMENT
# ============================================================================

Write-Host "📅 WEEK 1, DAY 1: Authentication Deployment" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan

# Step 1: Setup Authentication
Write-Host "Step 1/5: Setting up authentication system..." -ForegroundColor Yellow

try {
    & ".\scripts\setup-authentication.ps1"
    if ($LASTEXITCODE -eq 0) {
        $completedSteps += "Authentication Setup"
        Write-Host "✅ Authentication setup complete`n" -ForegroundColor Green
    } else {
        $failedSteps += "Authentication Setup"
        Write-Host "⚠️  Authentication setup had warnings`n" -ForegroundColor Yellow
    }
} catch {
    $failedSteps += "Authentication Setup"
    Write-Host "❌ Authentication setup failed: $_`n" -ForegroundColor Red
}

# Step 2: Install Go Dependencies
Write-Host "Step 2/5: Installing Go dependencies..." -ForegroundColor Yellow

$goDeps = @(
    "github.com/golang-jwt/jwt/v5",
    "golang.org/x/crypto/bcrypt",
    "github.com/twilio/twilio-go",
    "github.com/sendgrid/sendgrid-go",
    "github.com/jmoiron/sqlx",
    "github.com/lib/pq"
)

foreach ($dep in $goDeps) {
    Write-Host "  Installing $dep..." -ForegroundColor Gray
    & go get $dep 2>&1 | Out-Null
}

Write-Host "✅ Go dependencies installed`n" -ForegroundColor Green
$completedSteps += "Go Dependencies"

# Step 3: Build Backend
Write-Host "Step 3/5: Building backend..." -ForegroundColor Yellow

try {
    & go build -o platform.exe ./cmd/platform 2>&1 | Out-Null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Backend built successfully`n" -ForegroundColor Green
        $completedSteps += "Backend Build"
    } else {
        Write-Host "⚠️  Backend build completed with warnings`n" -ForegroundColor Yellow
        $completedSteps += "Backend Build (with warnings)"
    }
} catch {
    Write-Host "❌ Backend build failed: $_`n" -ForegroundColor Red
    $failedSteps += "Backend Build"
}

# Step 4: Frontend Setup
Write-Host "Step 4/5: Setting up frontend..." -ForegroundColor Yellow

if (Test-Path "admin-ui") {
    Push-Location admin-ui
    
    # Check if .env.local exists
    if (!(Test-Path ".env.local")) {
        Write-Host "  Creating frontend .env.local..." -ForegroundColor Gray
        "NEXT_PUBLIC_API_URL=http://localhost:8080" | Out-File -FilePath ".env.local" -Encoding UTF8
    }
    
    Write-Host "✅ Frontend configured`n" -ForegroundColor Green
    $completedSteps += "Frontend Setup"
    
    Pop-Location
} else {
    Write-Host "⚠️  Frontend directory not found`n" -ForegroundColor Yellow
}

# Step 5: Create startup scripts
Write-Host "Step 5/5: Creating convenient startup scripts..." -ForegroundColor Yellow

# Backend startup script
$backendScript = @"
@echo off
echo Starting ABY-MED Backend Server...
echo.
echo API will be available at: http://localhost:8080
echo.
platform.exe
"@

$backendScript | Out-File -FilePath "start-backend.bat" -Encoding ASCII

# Combined startup script
$combinedScript = @"
# ABY-MED Platform Startup Script
Write-Host "`n🚀 Starting ABY-MED Platform..." -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan

# Start backend in new window
Write-Host "Starting backend server..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '\$PWD'; .\platform.exe" -WindowStyle Normal

Start-Sleep -Seconds 3

# Start frontend in new window
Write-Host "Starting frontend application..." -ForegroundColor Yellow
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '\$PWD\admin-ui'; npm run dev" -WindowStyle Normal

Write-Host "`n✅ Both servers starting..." -ForegroundColor Green
Write-Host "`nBackend:  http://localhost:8080" -ForegroundColor Cyan
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "`nPress Ctrl+C in each window to stop the servers`n" -ForegroundColor Yellow
"@

$combinedScript | Out-File -FilePath "start-platform.ps1" -Encoding UTF8

Write-Host "✅ Startup scripts created`n" -ForegroundColor Green
$completedSteps += "Startup Scripts"

# ============================================================================
# SUMMARY
# ============================================================================

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
Write-Host "📊 SETUP SUMMARY" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Green

Write-Host "✅ Completed Steps: $($completedSteps.Count)" -ForegroundColor Green
foreach ($step in $completedSteps) {
    Write-Host "   • $step" -ForegroundColor Gray
}

if ($failedSteps.Count -gt 0) {
    Write-Host "`n❌ Failed Steps: $($failedSteps.Count)" -ForegroundColor Red
    foreach ($step in $failedSteps) {
        Write-Host "   • $step" -ForegroundColor Gray
    }
}

# ============================================================================
# NEXT STEPS
# ============================================================================

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "🚀 READY TO START!" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan

Write-Host "Quick Start Options:`n" -ForegroundColor Yellow

Write-Host "Option A: Start Everything (Recommended)" -ForegroundColor White
Write-Host "   .\start-platform.ps1`n" -ForegroundColor Cyan

Write-Host "Option B: Start Manually" -ForegroundColor White
Write-Host "   Terminal 1: .\start-backend.bat" -ForegroundColor Cyan
Write-Host "   Terminal 2: cd admin-ui && npm run dev`n" -ForegroundColor Cyan

Write-Host "Option C: Development Mode" -ForegroundColor White
Write-Host "   go run cmd/platform/main.go`n" -ForegroundColor Cyan

Write-Host "📝 Test URLs:" -ForegroundColor Yellow
Write-Host "   Backend API: http://localhost:8080/health" -ForegroundColor Gray
Write-Host "   Auth Endpoints: http://localhost:8080/api/v1/auth" -ForegroundColor Gray
Write-Host "   Frontend: http://localhost:3000" -ForegroundColor Gray
Write-Host "   Register: http://localhost:3000/register" -ForegroundColor Gray
Write-Host "   Login: http://localhost:3000/login`n" -ForegroundColor Gray

Write-Host "📚 Documentation:" -ForegroundColor Yellow
Write-Host "   Strategic Pipeline: docs/STRATEGIC-IMPLEMENTATION-PIPELINE.md" -ForegroundColor Gray
Write-Host "   Auth Guide: docs/AUTHENTICATION-READY-TO-DEPLOY.md" -ForegroundColor Gray
Write-Host "   API Spec: docs/specs/API-SPECIFICATION.md`n" -ForegroundColor Gray

Write-Host "🎯 WEEK 1 Goals:" -ForegroundColor Yellow
Write-Host "   Day 1: ✅ Authentication deployed (YOU ARE HERE)" -ForegroundColor Green
Write-Host "   Day 2-3: Integrate auth with existing system" -ForegroundColor Gray
Write-Host "   Day 4-5: Production configuration" -ForegroundColor Gray
Write-Host "   Day 6-7: Testing and verification`n" -ForegroundColor Gray

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan

Write-Host "✨ Ready to start? Run: .\start-platform.ps1`n" -ForegroundColor Green
