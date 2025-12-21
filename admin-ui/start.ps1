# Add Node.js to PATH
$env:Path = "C:\Program Files\nodejs;$env:Path"

Set-Location "C:\Users\birju\aby-med\admin-ui"

# Restore tsconfig
if (Test-Path "tsconfig.json.bak") {
    Rename-Item "tsconfig.json.bak" "tsconfig.json" -Force
}

# Install TypeScript if missing
if (!(Test-Path "node_modules\typescript")) {
    Write-Host "Installing TypeScript..." -ForegroundColor Yellow
    npm install typescript @types/node @types/react @types/react-dom --save-dev --force
}

Write-Host ""
Write-Host "Starting Frontend on http://localhost:3000" -ForegroundColor Green
Write-Host ""

# Start Next.js
npm run dev
