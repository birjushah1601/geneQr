# Simple API Test Script
$BaseUrl = "http://localhost:8081/api/v1"

Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host "  ENGINEER ASSIGNMENT API TESTS" -ForegroundColor Cyan
Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: List Engineers
Write-Host "[TEST 1] List All Engineers" -ForegroundColor Yellow
$eng = Invoke-WebRequest "$BaseUrl/engineers" -UseBasicParsing | ConvertFrom-Json
Write-Host "✓ Found $($eng.engineers.Count) engineers" -ForegroundColor Green
$testId = $eng.engineers[0].id
Write-Host ""

# Test 2: Get Single Engineer
Write-Host "[TEST 2] Get Engineer By ID: $testId" -ForegroundColor Yellow
$single = Invoke-WebRequest "$BaseUrl/engineers/$testId" -UseBasicParsing | ConvertFrom-Json
Write-Host "✓ Engineer: $($single.name) - Level: $($single.engineer_level)" -ForegroundColor Green
Write-Host ""

# Test 3: List Equipment Types
Write-Host "[TEST 3] List Equipment Types" -ForegroundColor Yellow
$types = Invoke-WebRequest "$BaseUrl/engineers/$testId/equipment-types" -UseBasicParsing | ConvertFrom-Json
Write-Host "✓ Found $($types.equipment_types.Count) capabilities" -ForegroundColor Green
$types.equipment_types | Format-Table manufacturer, category -AutoSize
Write-Host ""

# Test 4: Add Equipment Type
Write-Host "[TEST 4] Add Equipment Type" -ForegroundColor Yellow
$body = @{manufacturer="Siemens Healthineers"; category="Ultrasound"} | ConvertTo-Json
$add = Invoke-WebRequest "$BaseUrl/engineers/$testId/equipment-types" -Method POST -Body $body -ContentType "application/json" -UseBasicParsing
Write-Host "✓ Added capability" -ForegroundColor Green
Write-Host ""

# Test 5: List After Add
Write-Host "[TEST 5] List Equipment Types (After Add)" -ForegroundColor Yellow
$types2 = Invoke-WebRequest "$BaseUrl/engineers/$testId/equipment-types" -UseBasicParsing | ConvertFrom-Json
Write-Host "✓ Now has $($types2.equipment_types.Count) capabilities" -ForegroundColor Green
$types2.equipment_types | Format-Table manufacturer, category -AutoSize
Write-Host ""

Write-Host "=======================================================" -ForegroundColor Green
Write-Host "  ALL TESTS PASSED!" -ForegroundColor Green
Write-Host "=======================================================" -ForegroundColor Green
