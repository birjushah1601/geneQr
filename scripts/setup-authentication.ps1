# Complete Authentication Setup Script

Write-Host "`nâ”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”" -ForegroundColor Green
Write-Host "ðŸš€ AUTHENTICATION SYSTEM SETUP" -ForegroundColor Green
Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Green

$ErrorActionPreference = "Continue"
$setupErrors = @()

# Step 1: Generate JWT Keys
Write-Host "Step 1/3: Generating JWT Keys..." -ForegroundColor Cyan
Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Cyan

try {
    & ".\scripts\generate-jwt-keys.ps1"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "âœ… JWT keys generated successfully`n" -ForegroundColor Green
    } else {
        $setupErrors += "Failed to generate JWT keys"
        Write-Host "âŒ JWT key generation failed`n" -ForegroundColor Red
    }
} catch {
    $setupErrors += "Error running key generation script: $_"
    Write-Host "âŒ Error: $_`n" -ForegroundColor Red
}

# Step 2: Apply Database Migrations
Write-Host "`nStep 2/3: Applying Database Migrations..." -ForegroundColor Cyan
Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Cyan

try {
    Write-Host "Running migration script..." -ForegroundColor Yellow
    $migrationOutput = & go run scripts/apply-auth-migrations.go 2>&1
    Write-Host $migrationOutput
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nâœ… Database migrations applied successfully`n" -ForegroundColor Green
    } else {
        $setupErrors += "Database migration failed"
        Write-Host "`nâŒ Database migration failed`n" -ForegroundColor Red
    }
} catch {
    $setupErrors += "Error applying migrations: $_"
    Write-Host "âŒ Error: $_`n" -ForegroundColor Red
}

# Step 3: Update .env file
Write-Host "`nStep 3/3: Checking Environment Configuration..." -ForegroundColor Cyan
Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Cyan

$envFile = ".env"
$envLocalFile = ".env.local"
$envPath = $envLocalFile

if (!(Test-Path $envFile) -and !(Test-Path $envLocalFile)) {
    Write-Host "Creating .env.local file..." -ForegroundColor Yellow
    
    $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5430/med_platform?sslmode=disable

# JWT Configuration
JWT_PRIVATE_KEY_PATH=./keys/jwt-private.pem
JWT_PUBLIC_KEY_PATH=./keys/jwt-public.pem
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h
JWT_ISSUER=servqr-platform

# OTP Configuration
OTP_LENGTH=6
OTP_EXPIRY_MINUTES=5
OTP_MAX_ATTEMPTS=3
OTP_RATE_LIMIT_PER_HOUR=3
OTP_COOLDOWN_SECONDS=60

# Password Configuration
PASSWORD_BCRYPT_COST=12
PASSWORD_MIN_LENGTH=8

# Auth Configuration
MAX_FAILED_ATTEMPTS=5
LOCKOUT_DURATION=30m
ALLOW_REGISTRATION=true

# Twilio (SMS/WhatsApp) - Optional for development
# TWILIO_ACCOUNT_SID=your_account_sid
# TWILIO_AUTH_TOKEN=your_auth_token
# TWILIO_PHONE_NUMBER=+1234567890
# TWILIO_WHATSAPP_NUMBER=+1234567890

# SendGrid (Email) - Optional for development
# SENDGRID_API_KEY=your_api_key
# SENDGRID_FROM_EMAIL=noreply@ServQR.com
# SENDGRID_FROM_NAME=ServQR Platform

# Server Configuration
SERVER_PORT=8080
"@
    
    $envContent | Out-File -FilePath $envLocalFile -Encoding UTF8
    Write-Host "âœ… Created .env.local with default configuration`n" -ForegroundColor Green
} else {
    Write-Host "âœ… Environment file exists ($envPath)`n" -ForegroundColor Green
    Write-Host "âš ï¸  Make sure it includes JWT configuration:
   JWT_PRIVATE_KEY_PATH=./keys/jwt-private.pem
   JWT_PUBLIC_KEY_PATH=./keys/jwt-public.pem`n" -ForegroundColor Yellow
}

# Summary
Write-Host "`nâ”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”" -ForegroundColor Green
if ($setupErrors.Count -eq 0) {
    Write-Host "âœ… SETUP COMPLETE!" -ForegroundColor Green
    Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Green
    
    Write-Host "ðŸŽ‰ Authentication system is ready!`n" -ForegroundColor Cyan
    
    Write-Host "ðŸ“ Next Steps:" -ForegroundColor Cyan
    Write-Host "   1. Configure external services (optional):" -ForegroundColor White
    Write-Host "      â€¢ Twilio (SMS/WhatsApp): Edit .env.local" -ForegroundColor Gray
    Write-Host "      â€¢ SendGrid (Email): Edit .env.local`n" -ForegroundColor Gray
    
    Write-Host "   2. Start the backend server:" -ForegroundColor White
    Write-Host "      go run cmd/platform/main.go`n" -ForegroundColor Gray
    
    Write-Host "   3. Start the frontend (in another terminal):" -ForegroundColor White
    Write-Host "      cd admin-ui" -ForegroundColor Gray
    Write-Host "      npm run dev`n" -ForegroundColor Gray
    
    Write-Host "   4. Test authentication:" -ForegroundColor White
    Write-Host "      â€¢ Open: http://localhost:3000/register" -ForegroundColor Gray
    Write-Host "      â€¢ Open: http://localhost:3000/login`n" -ForegroundColor Gray
    
    Write-Host "ðŸ“š Documentation:" -ForegroundColor Cyan
    Write-Host "   â€¢ Complete Guide: docs/PHASE1-COMPLETE.md" -ForegroundColor White
    Write-Host "   â€¢ API Reference: docs/specs/API-SPECIFICATION.md`n" -ForegroundColor White
    
    Write-Host "ðŸ’¡ Development Mode:" -ForegroundColor Yellow
    Write-Host "   â€¢ OTP codes will be logged to console (mock services)" -ForegroundColor Yellow
    Write-Host "   â€¢ Configure Twilio/SendGrid for real email/SMS`n" -ForegroundColor Yellow
    
} else {
    Write-Host "âš ï¸  SETUP COMPLETED WITH WARNINGS" -ForegroundColor Yellow
    Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Yellow
    
    Write-Host "âŒ Errors encountered:" -ForegroundColor Red
    foreach ($error in $setupErrors) {
        Write-Host "   â€¢ $error" -ForegroundColor Red
    }
    Write-Host "`nðŸ“ Please review and fix the errors above.`n" -ForegroundColor Yellow
}

Write-Host "â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”`n" -ForegroundColor Green
