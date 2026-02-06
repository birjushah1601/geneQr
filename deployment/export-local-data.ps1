# Export Local Database Data to Production
# This script exports data from local PostgreSQL and prepares it for production import

param(
    [string]$LocalHost = "localhost",
    [string]$LocalPort = "5432",
    [string]$LocalUser = "postgres",
    [string]$LocalDB = "medical_equipment",
    [string]$OutputFile = "local_data_export.sql"
)

Write-Host "=== Local Database Export ===" -ForegroundColor Cyan
Write-Host ""

# Check if pg_dump is available
$pgDump = Get-Command pg_dump -ErrorAction SilentlyContinue
if (-not $pgDump) {
    Write-Host "ERROR: pg_dump not found. Please install PostgreSQL client tools." -ForegroundColor Red
    Write-Host "Download from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    exit 1
}

Write-Host "Exporting data from local database..." -ForegroundColor Yellow
Write-Host "  Host: $LocalHost"
Write-Host "  Database: $LocalDB"
Write-Host "  Output: $OutputFile"
Write-Host ""

# Export data only (no schema, using INSERT statements for better compatibility)
$env:PGPASSWORD = Read-Host "Enter PostgreSQL password for user '$LocalUser'" -AsSecureString | ConvertFrom-SecureString -AsPlainText

try {
    # Export all data
    pg_dump -h $LocalHost -p $LocalPort -U $LocalUser -d $LocalDB `
        --data-only `
        --inserts `
        --disable-triggers `
        --no-owner `
        --no-privileges `
        -f $OutputFile

    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✓ Data exported successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "File: $OutputFile" -ForegroundColor Cyan
        Write-Host "Size: $([math]::Round((Get-Item $OutputFile).Length / 1MB, 2)) MB" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Next steps:" -ForegroundColor Yellow
        Write-Host "1. Upload this file to your production server"
        Write-Host "2. Run the import script on the server"
        Write-Host ""
        Write-Host "Upload command:" -ForegroundColor Cyan
        Write-Host "  scp $OutputFile root@158.69.118.34:/tmp/" -ForegroundColor White
    } else {
        Write-Host "ERROR: Export failed!" -ForegroundColor Red
        exit 1
    }
} catch {
    $errorMsg = $_.Exception.Message
    Write-Host "ERROR: $errorMsg" -ForegroundColor Red
    exit 1
} finally {
    $env:PGPASSWORD = ""
}
