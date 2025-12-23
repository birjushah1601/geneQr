# Test Organizations Bulk Import API
Write-Host 'Testing Organizations Bulk Import API' -ForegroundColor Cyan
Write-Host 'Endpoint: POST /api/v1/organizations/import' -ForegroundColor Yellow
Write-Host ''
Write-Host 'Usage: Invoke-RestMethod -Uri http://localhost:8081/api/v1/organizations/import -Method Post -Form @{ csv_file = Get-Item templates/csv/organizations-import-template.csv }'
