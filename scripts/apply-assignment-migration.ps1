# Apply Engineer Assignment Migration Script
# This script applies the engineer assignment migration and seed data

param(
    [string]$Host = "localhost",
    [string]$Port = "5430",
    [string]$User = "postgres",
    [string]$Password = "postgres",
    [string]$Database = "med_platform"
)

Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "Engineer Assignment Migration Tool" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host ""

# Check if PostgreSQL client is available
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue

if (-not $psqlPath) {
    Write-Host "ERROR: psql command not found in PATH" -ForegroundColor Red
    Write-Host "Checking for Docker alternative..." -ForegroundColor Yellow
    
    # Try Docker approach
    $dockerRunning = docker ps 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Docker is not running either" -ForegroundColor Red
        Write-Host ""
        Write-Host "Please start Docker Desktop and ensure PostgreSQL container is running:" -ForegroundColor Yellow
        Write-Host "  docker-compose up -d" -ForegroundColor White
        Write-Host ""
        Write-Host "OR install PostgreSQL client tools and add to PATH" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "Using Docker to apply migrations..." -ForegroundColor Green
    
    # Find PostgreSQL container
    $containerName = docker ps --filter "ancestor=postgres" --format "{{.Names}}" 2>&1
    if (-not $containerName) {
        Write-Host "ERROR: No PostgreSQL container found" -ForegroundColor Red
        Write-Host "Start the database with: docker-compose up -d" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "Found PostgreSQL container: $containerName" -ForegroundColor Green
    
    # Apply migration via Docker
    Write-Host ""
    Write-Host "Step 1: Applying migration 003..." -ForegroundColor Cyan
    docker exec -i $containerName psql -U $User -d $Database < database/migrations/003_simplified_engineer_assignment.sql
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Migration failed" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Migration 003 applied successfully" -ForegroundColor Green
    
    # Apply seed data
    Write-Host ""
    Write-Host "Step 2: Applying seed data 005..." -ForegroundColor Cyan
    docker exec -i $containerName psql -U $User -d $Database < database/seed/005_engineer_assignment_data.sql
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "WARNING: Seed data partially applied (may be due to missing foreign keys)" -ForegroundColor Yellow
    } else {
        Write-Host "✓ Seed data applied successfully" -ForegroundColor Green
    }
    
} else {
    Write-Host "Using psql to apply migrations..." -ForegroundColor Green
    
    # Set password environment variable
    $env:PGPASSWORD = $Password
    
    # Apply migration
    Write-Host ""
    Write-Host "Step 1: Applying migration 003..." -ForegroundColor Cyan
    psql -h $Host -p $Port -U $User -d $Database -f database/migrations/003_simplified_engineer_assignment.sql
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Migration failed" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Migration 003 applied successfully" -ForegroundColor Green
    
    # Apply seed data
    Write-Host ""
    Write-Host "Step 2: Applying seed data 005..." -ForegroundColor Cyan
    psql -h $Host -p $Port -U $User -d $Database -f database/seed/005_engineer_assignment_data.sql
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "WARNING: Seed data partially applied" -ForegroundColor Yellow
    } else {
        Write-Host "✓ Seed data applied successfully" -ForegroundColor Green
    }
    
    # Clear password
    $env:PGPASSWORD = $null
}

# Verify migration
Write-Host ""
Write-Host "Step 3: Verifying migration..." -ForegroundColor Cyan

$verifyQuery = @"
SELECT 
    'engineer_level column' as check_item,
    CASE WHEN EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='engineers' AND column_name='engineer_level'
    ) THEN '✓ EXISTS' ELSE '✗ MISSING' END as status
UNION ALL
SELECT 
    'engineer_equipment_types table',
    CASE WHEN EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_name='engineer_equipment_types'
    ) THEN '✓ EXISTS' ELSE '✗ MISSING' END
UNION ALL
SELECT 
    'equipment_service_config table',
    CASE WHEN EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_name='equipment_service_config'
    ) THEN '✓ EXISTS' ELSE '✗ MISSING' END
UNION ALL
SELECT 
    'Engineers count',
    COALESCE(COUNT(*)::TEXT || ' engineers', '0 engineers')
FROM engineers WHERE engineer_level IS NOT NULL;
"@

if ($psqlPath) {
    $env:PGPASSWORD = $Password
    Write-Output $verifyQuery | psql -h $Host -p $Port -U $User -d $Database -t -A -F ' | '
    $env:PGPASSWORD = $null
} else {
    Write-Output $verifyQuery | docker exec -i $containerName psql -U $User -d $Database -t -A -F ' | '
}

Write-Host ""
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "Migration Complete!" -ForegroundColor Green
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Start backend server: .\start-backend.ps1 or go run cmd/platform/main.go" -ForegroundColor White
Write-Host "  2. Test APIs: curl http://localhost:8081/api/v1/engineers" -ForegroundColor White
Write-Host ""
