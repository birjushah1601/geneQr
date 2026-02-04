# ServQR Database Rebranding Migration Script
# Renames distributor to channel_partner and dealer to sub_dealer

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  ServQR Database Rebranding Migration" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

$dbName = "servqr_dev"
$dbUser = "postgres"
$migrationFile = "database\migrations\026_rename_distributor_to_channel_partner.sql"

if (-not (Test-Path $migrationFile)) {
    Write-Host "Migration file not found!" -ForegroundColor Red
    exit 1
}

Write-Host "Database: $dbName" -ForegroundColor Gray
Write-Host "Migration: $migrationFile" -ForegroundColor Gray
Write-Host ""

Write-Host "Enter PostgreSQL password:" -ForegroundColor Yellow
$env:PGPASSWORD = Read-Host

Write-Host ""
Write-Host "Running migration..." -ForegroundColor Cyan

& "C:\Program Files\PostgreSQL\17\bin\psql.exe" -U $dbUser -d $dbName -f $migrationFile

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Migration completed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Changes applied:" -ForegroundColor Cyan
    Write-Host "  distributor -> channel_partner" -ForegroundColor Green
    Write-Host "  dealer -> sub_dealer" -ForegroundColor Green
    Write-Host ""
} else {
    Write-Host "Migration failed!" -ForegroundColor Red
}

$env:PGPASSWORD = $null
