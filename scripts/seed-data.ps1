param(
    [string]$ApiBase = "http://localhost:8082/api/v1",
    [string]$CreatedBy = "seed-script"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function PostJson($url, $body) {
    $json = $body | ConvertTo-Json -Depth 6
    Invoke-RestMethod -Method Post -Uri $url -ContentType 'application/json' -Body $json
}

Write-Host "Seeding organizations (via ENABLE_ORG_SEED=true on backend startup)" -ForegroundColor Cyan

# 1) Seed Equipment via API
Write-Host "Seeding equipment..." -ForegroundColor Cyan

$equipmentPayloads = @(
    @{ serial_number = "SN-ECG-0001"; equipment_name = "ECG Machine Pro"; manufacturer_name = "Global Manufacturer A"; model_number = "ECG-1000"; category = "ECG"; customer_name = "City Hospital"; installation_location = "Cardiology, Floor 2"; purchase_price = 250000; warranty_months = 12; created_by = $CreatedBy },
    @{ serial_number = "SN-VTL-0001"; equipment_name = "Ventilator Max"; manufacturer_name = "Global Manufacturer A"; model_number = "VTL-500"; category = "Ventilator"; customer_name = "City Hospital"; installation_location = "ICU"; purchase_price = 850000; warranty_months = 24; created_by = $CreatedBy },
    @{ serial_number = "SN-INF-0001"; equipment_name = "Infusion Pump Lite"; manufacturer_name = "Supplier S"; model_number = "INF-200"; category = "Infusion Pump"; customer_name = "Metro Clinic"; installation_location = "Ward 3"; purchase_price = 65000; warranty_months = 12; created_by = $CreatedBy },
    @{ serial_number = "SN-XRY-0001"; equipment_name = "X-Ray System Alpha"; manufacturer_name = "Regional Distributor X"; model_number = "XR-900"; category = "X-Ray"; customer_name = "Metro Clinic"; installation_location = "Radiology"; purchase_price = 420000; warranty_months = 18; created_by = $CreatedBy },
    @{ serial_number = "SN-CT-0001"; equipment_name = "CT Scanner Nova"; manufacturer_name = "Global Manufacturer A"; model_number = "CT-16S"; category = "CT"; customer_name = "General Hospital"; installation_location = "Imaging"; purchase_price = 12500000; warranty_months = 12; created_by = $CreatedBy }
)

$equipmentIds = @()
foreach ($eq in $equipmentPayloads) {
    try {
        $res = PostJson "$ApiBase/equipment" $eq
        $equipmentIds += $res.id
        Write-Host "  + Equipment created: $($res.id) [$($eq.serial_number)]" -ForegroundColor Green
    } catch {
        Write-Host "  ! Failed to create equipment [$($eq.serial_number)]: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 2) Generate QR codes in bulk
try {
    $qrRes = Invoke-RestMethod -Method Post -Uri "$ApiBase/equipment/qr/bulk-generate"
    Write-Host "Bulk QR generation: $($qrRes.Message)" -ForegroundColor Green
} catch {
    Write-Host "Failed bulk QR generation: $($_.Exception.Message)" -ForegroundColor Yellow
}

# 3) Seed a couple of service tickets referencing first equipment
if ($equipmentIds.Count -gt 0) {
    Write-Host "Seeding service tickets..." -ForegroundColor Cyan
    $firstEq = $equipmentIds[0]
    $secondEq = if ($equipmentIds.Count -gt 1) { $equipmentIds[1] } else { $equipmentIds[0] }

    $ticket1 = @{
        EquipmentID      = $firstEq
        SerialNumber     = "SN-ECG-0001"
        EquipmentName    = "ECG Machine Pro"
        CustomerName     = "City Hospital"
        IssueCategory    = "breakdown"
        IssueDescription = "ECG display flickers intermittently"
        Priority         = "high"
        Source           = "web"
        CreatedBy        = $CreatedBy
        InitialComment   = "Reported by cardiology department"
    }

    $ticket2 = @{
        EquipmentID      = $secondEq
        SerialNumber     = "SN-VTL-0001"
        EquipmentName    = "Ventilator Max"
        CustomerName     = "City Hospital"
        IssueCategory    = "maintenance"
        IssueDescription = "Scheduled preventive maintenance"
        Priority         = "medium"
        Source           = "scheduled"
        CreatedBy        = $CreatedBy
        InitialComment   = "Auto-created for PM cycle"
    }

    foreach ($t in @($ticket1, $ticket2)) {
        try {
            $tRes = PostJson "$ApiBase/tickets" $t
            Write-Host "  + Ticket created: $($tRes.ticket_number) (ID: $($tRes.id))" -ForegroundColor Green
        } catch {
            Write-Host "  ! Failed to create ticket: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
}

Write-Host "Seeding complete." -ForegroundColor Green
