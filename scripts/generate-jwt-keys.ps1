# Generate JWT RSA Keys for Authentication

Write-Host "`n🔑 Generating JWT RSA Keys..." -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan

# Create keys directory if it doesn't exist
$keysDir = "keys"
if (!(Test-Path $keysDir)) {
    New-Item -ItemType Directory -Path $keysDir | Out-Null
    Write-Host "✅ Created 'keys' directory`n" -ForegroundColor Green
}

# Check if OpenSSL is available
$opensslPath = (Get-Command openssl -ErrorAction SilentlyContinue).Source

if ($opensslPath) {
    Write-Host "✅ OpenSSL found: $opensslPath`n" -ForegroundColor Green
    
    # Generate private key
    Write-Host "📝 Generating RSA private key (2048 bits)..." -ForegroundColor Yellow
    & openssl genrsa -out "$keysDir/jwt-private.pem" 2048 2>&1 | Out-Null
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Private key generated: $keysDir/jwt-private.pem`n" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to generate private key" -ForegroundColor Red
        exit 1
    }
    
    # Generate public key
    Write-Host "📝 Extracting RSA public key..." -ForegroundColor Yellow
    & openssl rsa -in "$keysDir/jwt-private.pem" -pubout -out "$keysDir/jwt-public.pem" 2>&1 | Out-Null
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Public key generated: $keysDir/jwt-public.pem`n" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to generate public key" -ForegroundColor Red
        exit 1
    }
    
    # Set permissions (read-only for private key)
    $privateKeyPath = Join-Path (Get-Location) "$keysDir\jwt-private.pem"
    $publicKeyPath = Join-Path (Get-Location) "$keysDir\jwt-public.pem"
    
    Write-Host "🔒 Setting file permissions..." -ForegroundColor Yellow
    icacls $privateKeyPath /inheritance:r /grant:r "$env:USERNAME:(R)" | Out-Null
    Write-Host "✅ Private key secured (read-only)`n" -ForegroundColor Green
    
    # Display key information
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "✅ JWT Keys Generated Successfully!" -ForegroundColor Green
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`n" -ForegroundColor Cyan
    
    Write-Host "📁 Key Files:" -ForegroundColor Cyan
    Write-Host "   Private Key: $privateKeyPath" -ForegroundColor White
    Write-Host "   Public Key:  $publicKeyPath`n" -ForegroundColor White
    
    Write-Host "⚠️  IMPORTANT:" -ForegroundColor Yellow
    Write-Host "   • Keep the private key secure and never commit it to git" -ForegroundColor Yellow
    Write-Host "   • The private key is used to sign JWT tokens" -ForegroundColor Yellow
    Write-Host "   • The public key is used to verify JWT tokens`n" -ForegroundColor Yellow
    
    Write-Host "📝 Next Steps:" -ForegroundColor Cyan
    Write-Host "   1. Add to .env file:" -ForegroundColor White
    Write-Host "      JWT_PRIVATE_KEY_PATH=./keys/jwt-private.pem" -ForegroundColor Gray
    Write-Host "      JWT_PUBLIC_KEY_PATH=./keys/jwt-public.pem`n" -ForegroundColor Gray
    Write-Host "   2. Apply database migrations:" -ForegroundColor White
    Write-Host "      go run scripts/apply-auth-migrations.go`n" -ForegroundColor Gray
    Write-Host "   3. Start the backend server`n" -ForegroundColor White
    
} else {
    Write-Host "❌ OpenSSL not found!" -ForegroundColor Red
    Write-Host "`nOpenSSL is required to generate RSA keys.`n" -ForegroundColor Yellow
    
    Write-Host "Installation Options:" -ForegroundColor Cyan
    Write-Host "  1. Install via Chocolatey:" -ForegroundColor White
    Write-Host "     choco install openssl`n" -ForegroundColor Gray
    
    Write-Host "  2. Download from:" -ForegroundColor White
    Write-Host "     https://slproweb.com/products/Win32OpenSSL.html`n" -ForegroundColor Gray
    
    Write-Host "  3. Use Git Bash (comes with Git for Windows):" -ForegroundColor White
    Write-Host "     git-bash.exe" -ForegroundColor Gray
    Write-Host "     openssl genrsa -out keys/jwt-private.pem 2048" -ForegroundColor Gray
    Write-Host "     openssl rsa -in keys/jwt-private.pem -pubout -out keys/jwt-public.pem`n" -ForegroundColor Gray
    
    exit 1
}
