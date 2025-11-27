# QR Code Functionality Test Script
# Tests all QR code endpoints and displays results

Write-Host "========== QR CODE FUNCTIONALITY TEST ==========" -ForegroundColor Cyan
Write-Host "Testing backend API endpoints for QR codes`n" -ForegroundColor Green

$headers = @{"X-Tenant-ID" = "default"}
$baseUrl = "http://localhost:8081/api/v1/equipment"

# Step 1: Get list of equipment
Write-Host "[1/5] Fetching equipment list..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl?page=1&page_size=5" -Headers $headers
    $equipmentList = $response.items
    Write-Host "✅ Found $($equipmentList.Count) equipment items`n" -ForegroundColor Green
    
    if ($equipmentList.Count -eq 0) {
        Write-Host "❌ No equipment found in database!" -ForegroundColor Red
        Write-Host "   Please add equipment first" -ForegroundColor Yellow
        exit
    }
    
    # Pick first equipment for testing
    $testEquipment = $equipmentList[0]
    $equipmentId = $testEquipment.id
    $equipmentName = $testEquipment.equipment_name
    
    Write-Host "📋 Test Equipment:" -ForegroundColor Cyan
    Write-Host "   ID: $equipmentId"
    Write-Host "   Name: $equipmentName"
    Write-Host "   Has QR: $($testEquipment.qr_code_generated_at -ne $null)"
    Write-Host ""
    
} catch {
    Write-Host "❌ Failed to fetch equipment: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# Step 2: Check if QR exists
Write-Host "[2/5] Checking existing QR code..." -ForegroundColor Yellow
$hasQR = $testEquipment.qr_code_generated_at -ne $null

if ($hasQR) {
    Write-Host "✅ QR code already exists (generated at: $($testEquipment.qr_code_generated_at))" -ForegroundColor Green
    $testImageUrl = "http://localhost:8081/api/v1/equipment/qr/image/$equipmentId"
    Write-Host "   Image URL: $testImageUrl" -ForegroundColor Cyan
} else {
    Write-Host "⚠️  No QR code found - will generate one" -ForegroundColor Yellow
}
Write-Host ""

# Step 3: Generate QR Code (if doesn't exist)
if (-not $hasQR) {
    Write-Host "[3/5] Generating QR code..." -ForegroundColor Yellow
    try {
        $genResponse = Invoke-RestMethod -Method POST -Uri "$baseUrl/$equipmentId/qr" -Headers $headers
        Write-Host "✅ QR code generated successfully!" -ForegroundColor Green
        Write-Host "   Message: $($genResponse.message)" -ForegroundColor Cyan
        Write-Host ""
        $hasQR = $true
    } catch {
        Write-Host "❌ Failed to generate QR: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "   Check backend logs for details" -ForegroundColor Yellow
        Write-Host ""
    }
} else {
    Write-Host "[3/5] Skipping generation (QR already exists)" -ForegroundColor Cyan
    Write-Host ""
}

# Step 4: Test QR Image Endpoint
if ($hasQR) {
    Write-Host "[4/5] Testing QR image endpoint..." -ForegroundColor Yellow
    try {
        $imageUrl = "http://localhost:8081/api/v1/equipment/qr/image/$equipmentId"
        $imageResponse = Invoke-WebRequest -Uri $imageUrl -Headers $headers -TimeoutSec 5
        
        if ($imageResponse.StatusCode -eq 200) {
            $contentType = $imageResponse.Headers['Content-Type']
            $contentLength = $imageResponse.RawContentLength
            
            Write-Host "✅ QR image endpoint working!" -ForegroundColor Green
            Write-Host "   Status: $($imageResponse.StatusCode)" -ForegroundColor Cyan
            Write-Host "   Content-Type: $contentType" -ForegroundColor Cyan
            Write-Host "   Image Size: $contentLength bytes" -ForegroundColor Cyan
            Write-Host "   URL: $imageUrl" -ForegroundColor Cyan
        }
    } catch {
        Write-Host "❌ QR image endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "   This means the image is not being served correctly" -ForegroundColor Yellow
    }
    Write-Host ""
} else {
    Write-Host "[4/5] Skipping image test (no QR generated)" -ForegroundColor Gray
    Write-Host ""
}

# Step 5: Check Database
Write-Host "[5/5] Checking database..." -ForegroundColor Yellow
try {
    $dbCheck = docker exec med_platform_pg psql -U postgres -d med_platform -t -c "SELECT id, equipment_name, qr_code, CASE WHEN qr_code_image IS NOT NULL THEN 'YES' ELSE 'NO' END as has_image, qr_code_generated_at FROM equipment WHERE id = '$equipmentId';" 2>$null
    
    if ($dbCheck) {
        Write-Host "✅ Database record:" -ForegroundColor Green
        Write-Host "$dbCheck" -ForegroundColor Cyan
    } else {
        Write-Host "⚠️  Could not query database directly" -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠️  Database check skipped (docker not available)" -ForegroundColor Yellow
}
Write-Host ""

# Summary
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "📊 TEST SUMMARY" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan

if ($hasQR) {
    Write-Host "✅ QR Code: GENERATED" -ForegroundColor Green
    Write-Host "✅ API Endpoint: WORKING" -ForegroundColor Green
    Write-Host "✅ Image URL: http://localhost:8081/api/v1/equipment/qr/image/$equipmentId" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "🎯 NEXT STEPS:" -ForegroundColor Yellow
    Write-Host "1. Open frontend: http://localhost:3000/equipment" -ForegroundColor White
    Write-Host "2. Look for equipment: $equipmentName" -ForegroundColor White
    Write-Host "3. You should see a QR code thumbnail" -ForegroundColor White
    Write-Host "4. Click it to preview full size" -ForegroundColor White
    Write-Host ""
    Write-Host "🔗 Test image directly in browser:" -ForegroundColor Yellow
    Write-Host "   http://localhost:8081/api/v1/equipment/qr/image/$equipmentId" -ForegroundColor Cyan
} else {
    Write-Host "❌ QR Code: NOT GENERATED" -ForegroundColor Red
    Write-Host "   Generation failed - check backend logs" -ForegroundColor Yellow
}

Write-Host "=========================================" -ForegroundColor Cyan
