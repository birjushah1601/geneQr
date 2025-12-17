Set-Location 'C:\Users\birju\aby-med'

$env:DB_HOST = 'localhost'
$env:DB_PORT = '5430'
$env:DB_USER = 'postgres'
$env:DB_PASSWORD = 'postgres'
$env:DB_NAME = 'med_platform'
$env:ENABLED_MODULES = '*'
$env:ENABLE_ORG = 'true'
$env:ENABLE_ORG_SEED = 'true'
$env:PORT = '8082'

Write-Host '🚀 Starting Backend API...' -ForegroundColor Green
Write-Host "   Port: $env:PORT" -ForegroundColor Gray
Write-Host "   DB: $env:DB_USER@$env:DB_HOST:$env:DB_PORT/$env:DB_NAME" -ForegroundColor Gray
Write-Host ""

.\backend.exe
