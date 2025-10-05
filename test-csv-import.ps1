# Test CSV Import for Equipment Registry
$csvPath = "manufacturer-installations-sample.csv"
$uri = "http://localhost:8081/api/v1/equipment/import"

# Read CSV file content
$csvContent = Get-Content -Path $csvPath -Raw

# Create boundary
$boundary = [System.Guid]::NewGuid().ToString()

# Build multipart form data
$LF = "`r`n"
$bodyLines = @(
    "--$boundary",
    "Content-Disposition: form-data; name=`"csv_file`"; filename=`"$csvPath`"",
    "Content-Type: text/csv",
    "",
    $csvContent,
    "--$boundary",
    "Content-Disposition: form-data; name=`"created_by`"",
    "",
    "manufacturer-onboard",
    "--$boundary--"
) -join $LF

# Make request
try {
    $response = Invoke-RestMethod -Uri $uri -Method Post -ContentType "multipart/form-data; boundary=$boundary" -Body $bodyLines -Headers @{"X-Tenant-ID"="city-hospital"}
    
    Write-Host "`n=== CSV IMPORT SUCCESS ===" -ForegroundColor Green
    Write-Host "Total Rows: $($response.total_rows)" -ForegroundColor Yellow
    Write-Host "Success Count: $($response.success_count)" -ForegroundColor Green
    Write-Host "Failure Count: $($response.failure_count)" -ForegroundColor $(if ($response.failure_count -gt 0) {"Red"} else {"Gray"})
    
    if ($response.errors -and $response.errors.Count -gt 0) {
        Write-Host "`nErrors:" -ForegroundColor Red
        $response.errors | ForEach-Object { Write-Host "  - $_" -ForegroundColor Red }
    }
    
    Write-Host "`nImported Equipment IDs:" -ForegroundColor Cyan
    $response.imported_ids | Select-Object -First 5 | ForEach-Object { Write-Host "  - $_" -ForegroundColor Gray }
    if ($response.imported_ids.Count -gt 5) {
        Write-Host "  ... and $($response.imported_ids.Count - 5) more" -ForegroundColor Gray
    }
    
    # Save full response
    $response | ConvertTo-Json -Depth 3 | Out-File "import-result.json"
    Write-Host "`nFull results saved to import-result.json" -ForegroundColor Green
    
    # Return first ID for testing
    return $response.imported_ids[0]
    
} catch {
    Write-Host "`n=== CSV IMPORT FAILED ===" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails) {
        Write-Host "Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
    }
}
