Set-Location 'C:\Users\birju\aby-med'

# Use med_platform_pg (docker container) on localhost:5430 only
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5430'
$env:DB_USER = 'postgres'
$env:DB_PASSWORD = 'postgres'
$env:DB_NAME = 'med_platform'

# Enable all modules
$env:ENABLED_MODULES = '*'

# Enable Organizations module (behind feature flag)
$env:ENABLE_ORG = 'true'

# API port (change if 8081 is busy)
$env:PORT = '8082'

Write-Host '?? Starting Backend API...' -ForegroundColor Green
Write-Host "?? Port: $env:PORT" -ForegroundColor Gray
Write-Host "???  DB: $env:DB_USER@$env:DB_HOST:$env:DB_PORT/$env:DB_NAME" -ForegroundColor Gray
Write-Host "?? Modules: $env:ENABLED_MODULES" -ForegroundColor Gray
Write-Host '?? Running: cmd/platform/main.go' -ForegroundColor Gray
Write-Host ''

go run cmd/platform/main.go
