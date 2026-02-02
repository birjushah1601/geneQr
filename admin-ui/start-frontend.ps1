$ErrorActionPreference = "Stop"

# Ensure Node.js is FIRST in PATH (critical for subprocesses)
$nodePath = "C:\Program Files\nodejs"
$env:Path = "$nodePath;$env:Path"

# Navigate to admin-ui directory
Set-Location "C:\Users\birju\ServQR\admin-ui"

# Restore tsconfig if it was renamed
if (Test-Path "tsconfig.json.bak") {
    Rename-Item "tsconfig.json.bak" "tsconfig.json" -Force
}

# Install typescript if missing
if (!(Test-Path "node_modules\typescript")) {
    Write-Host "ðŸ“¦ Installing TypeScript..." -ForegroundColor Yellow
    npm install typescript @types/node @types/react @types/react-dom --save-dev --force
}

Write-Host ""
Write-Host "ðŸš€ Starting Frontend (Next.js)..." -ForegroundColor Green
Write-Host "   URL: http://localhost:3000" -ForegroundColor Cyan
Write-Host "   Press Ctrl+C to stop" -ForegroundColor Gray
Write-Host ""

# Start with npm (which will use node from PATH)
npm run dev
