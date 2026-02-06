# Cleanup Dead Files - Review what will be removed
# This script shows files that are ignored by .gitignore and can be safely deleted

Write-Host "=== Files that will be IGNORED by Git (safe to delete) ===" -ForegroundColor Cyan
Write-Host ""

# Get all files that match .gitignore patterns
$deadFiles = @(
    Get-ChildItem -Path ".\" -Filter "backend-*.log*" -File
    Get-ChildItem -Path ".\" -Filter "*.sql" -File | Where-Object { 
        $_.Name -notlike "*migrations*" -and $_.Name -ne "init-database-schema.sql"
    }
    Get-ChildItem -Path ".\" -Filter "*.exe" -File
    Get-ChildItem -Path ".\" -Filter "test-*.ps1" -File
    Get-ChildItem -Path ".\" -Filter "start-*.ps1" -File
    Get-ChildItem -Path ".\" -Filter "*-*.ps1" -File | Where-Object {
        $_.Name -notlike "test-*.ps1" -and $_.Name -notlike "start-*.ps1"
    }
)

$totalSize = 0
Write-Host "Dead files found:" -ForegroundColor Yellow
foreach ($file in $deadFiles) {
    $size = [math]::Round($file.Length / 1KB, 2)
    $totalSize += $size
    Write-Host "  $($file.Name) - $size KB"
}

Write-Host ""
Write-Host "Total: $($deadFiles.Count) files, $([math]::Round($totalSize/1024, 2)) MB" -ForegroundColor Cyan
Write-Host ""
Write-Host "Run 'git clean -fdx' to remove all untracked/ignored files" -ForegroundColor Yellow
Write-Host "(Or delete manually)"
