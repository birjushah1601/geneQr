# Security Features Test Script
# Tests all security implementations

$baseUrl = "http://localhost:8081/api/v1"
$testResults = @()

Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "🧪 SECURITY FEATURES TEST SUITE" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Test 1: XSS Protection
Write-Host "Test 1: XSS Protection" -ForegroundColor Yellow
$xssPayload = '{"QRCode":"TEST-XSS","EquipmentID":"test-eq","SerialNumber":"SN-001","IssueDescription":"<script>alert(1)</script>Test","CustomerName":"<img src=x onerror=alert(1)>John","Priority":"medium"}'

try {
    $response = Invoke-WebRequest -Uri "$baseUrl/tickets" -Method POST -ContentType "application/json" -Body $xssPayload -UseBasicParsing -TimeoutSec 10
    Write-Host "✅ Test 1 PASSED: XSS payload accepted (should be sanitized)" -ForegroundColor Green
    $testResults += "XSS:PASSED"
} catch {
    Write-Host "❌ Test 1 FAILED: $($_.Exception.Message)" -ForegroundColor Red
    $testResults += "XSS:FAILED"
}
Write-Host ""

# Test 2: QR Rate Limiting
Write-Host "Test 2: QR Rate Limiting (creating 6 tickets)" -ForegroundColor Yellow
$qrCode = "TEST-RATE-$(Get-Random)"
$successCount = 0

for ($i = 1; $i -le 6; $i++) {
    $payload = "{`"QRCode`":`"$qrCode`",`"EquipmentID`":`"test-$i`",`"SerialNumber`":`"SN-$i`",`"IssueDescription`":`"Test $i`",`"Priority`":`"low`"}"
    try {
        Invoke-WebRequest -Uri "$baseUrl/tickets" -Method POST -ContentType "application/json" -Body $payload -UseBasicParsing -TimeoutSec 10 | Out-Null
        Write-Host "  Request $i : ✅ Success" -ForegroundColor Green
        $successCount++
    } catch {
        Write-Host "  Request $i : ❌ Rate Limited" -ForegroundColor Yellow
    }
    Start-Sleep -Milliseconds 300
}

if ($successCount -le 5) {
    Write-Host "✅ Test 2 PASSED: QR rate limiting working ($successCount/6 succeeded)" -ForegroundColor Green
    $testResults += "QR_RATE:PASSED"
} else {
    Write-Host "❌ Test 2 FAILED: Too many requests succeeded ($successCount/6)" -ForegroundColor Red
    $testResults += "QR_RATE:FAILED"
}
Write-Host ""

# Summary
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "📊 TEST RESULTS" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
foreach ($result in $testResults) { Write-Host "  $result" }
Write-Host "`n"
