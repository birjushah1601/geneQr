# Test Equipment Registration Manually (bypassing CSV for now)
Write-Host "`n=== TESTING EQUIPMENT REGISTRATION ===" -ForegroundColor Cyan

$uri = "http://localhost:8081/api/v1/equipment"
$headers = @{"X-Tenant-ID"="city-hospital"; "Content-Type"="application/json"}

# Test Equipment 1: MRI Scanner
$equipment1 = @{
    serial_number = "MED-MRI-001"
    equipment_name = "MRI Scanner 1.5T"
    manufacturer_name = "Siemens Healthineers"
    model_number = "MAGNETOM Sola"
    category = "Diagnostic Imaging"
    customer_name = "Apollo Hospitals Mumbai"
    customer_id = "CUST-001"
    installation_location = "Radiology Department - Building A"
    installation_date = "2023-06-15T00:00:00Z"
    purchase_date = "2023-05-20T00:00:00Z"
    purchase_price = 125000000
    warranty_months = 24
    notes = "Primary MRI unit for radiology department"
    created_by = "manufacturer-onboard"
} | ConvertTo-Json

Write-Host "`nTest 1: Registering MRI Scanner..." -ForegroundColor Yellow
try {
    $result1 = Invoke-RestMethod -Uri $uri -Method Post -Headers $headers -Body $equipment1
    Write-Host "✅ SUCCESS - Equipment ID: $($result1.id)" -ForegroundColor Green
    Write-Host "   QR Code: $($result1.qr_code)" -ForegroundColor Gray
    Write-Host "   Serial: $($result1.serial_number)" -ForegroundColor Gray
    $equipmentId1 = $result1.id
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Equipment 2: CT Scanner
$equipment2 = @{
    serial_number = "MED-CT-002"
    equipment_name = "CT Scanner 64-Slice"
    manufacturer_name = "GE Healthcare"
    model_number = "Revolution EVO"
    category = "Diagnostic Imaging"
    customer_name = "Fortis Hospital Delhi"
    customer_id = "CUST-002"
    installation_location = "Emergency Wing - 2nd Floor"
    installation_date = "2023-07-10T00:00:00Z"
    purchase_date = "2023-06-25T00:00:00Z"
    purchase_price = 89000000
    warranty_months = 24
    notes = "Emergency CT scanner with cardiac imaging"
    created_by = "manufacturer-onboard"
} | ConvertTo-Json

Write-Host "`nTest 2: Registering CT Scanner..." -ForegroundColor Yellow
try {
    $result2 = Invoke-RestMethod -Uri $uri -Method Post -Headers $headers -Body $equipment2
    Write-Host "✅ SUCCESS - Equipment ID: $($result2.id)" -ForegroundColor Green
    Write-Host "   QR Code: $($result2.qr_code)" -ForegroundColor Gray
    Write-Host "   Serial: $($result2.serial_number)" -ForegroundColor Gray
    $equipmentId2 = $result2.id
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

# Test Equipment 3: ICU Ventilator
$equipment3 = @{
    serial_number = "MED-VENT-003"
    equipment_name = "ICU Ventilator"
    manufacturer_name = "Medtronic"
    model_number = "PB980"
    category = "Critical Care"
    customer_name = "AIIMS New Delhi"
    customer_id = "CUST-003"
    installation_location = "ICU Ward 3 - Bed 12"
    installation_date = "2023-08-05T00:00:00Z"
    purchase_date = "2023-07-15T00:00:00Z"
    purchase_price = 1800000
    warranty_months = 36
    notes = "Advanced ventilator with NAVA mode"
    created_by = "manufacturer-onboard"
} | ConvertTo-Json

Write-Host "`nTest 3: Registering ICU Ventilator..." -ForegroundColor Yellow
try {
    $result3 = Invoke-RestMethod -Uri $uri -Method Post -Headers $headers -Body $equipment3
    Write-Host "✅ SUCCESS - Equipment ID: $($result3.id)" -ForegroundColor Green
    Write-Host "   QR Code: $($result3.qr_code)" -ForegroundColor Gray
    Write-Host "   Serial: $($result3.serial_number)" -ForegroundColor Gray
    $equipmentId3 = $result3.id
} catch {
    Write-Host "❌ FAILED: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== REGISTRATION SUMMARY ===" -ForegroundColor Cyan
Write-Host "Registered 3 test equipment items" -ForegroundColor Green

# Save equipment IDs for QR code generation
@{
    equipment_ids = @($equipmentId1, $equipmentId2, $equipmentId3)
} | ConvertTo-Json | Out-File "equipment-ids.json"

Write-Host "Equipment IDs saved to equipment-ids.json" -ForegroundColor Gray

# Return first ID
return $equipmentId1
