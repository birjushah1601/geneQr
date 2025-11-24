# ==============================================================================
# COMPREHENSIVE ENGINEER ASSIGNMENT API TEST SUITE
# ==============================================================================

$ErrorActionPreference = "Continue"
$BaseUrl = "http://localhost:8081/api/v1"
$script:TestResults = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Url,
        [object]$Body = $null
    )
    
    Write-Host "Testing: $Name" -ForegroundColor Yellow
    
    try {
        $params = @{
            Uri = "$BaseUrl$Url"
            Method = $Method
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json -Depth 10)
            $params.ContentType = "application/json"
        }
        
        $response = Invoke-WebRequest @params
        Write-Host "✓ SUCCESS - Status: $($response.StatusCode)" -ForegroundColor Green
        
        $script:TestResults += [PSCustomObject]@{
            Test = $Name
            Status = "PASSED"
            StatusCode = $response.StatusCode
        }
        
        return ($response.Content | ConvertFrom-Json)
    }
    catch {
        Write-Host "✗ FAILED - $($_.Exception.Message)" -ForegroundColor Red
        
        $script:TestResults += [PSCustomObject]@{
            Test = $Name
            Status = "FAILED"
            Error = $_.Exception.Message
        }
        
        return $null
    }
}

Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host "   ENGINEER ASSIGNMENT API TEST SUITE" -ForegroundColor Cyan
Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host ""

# ==============================================================================
# SECTION 1: ENGINEER MANAGEMENT
# ==============================================================================
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "SECTION 1: Engineer Management APIs" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Test 1.1: List All Engineers
$engineers = Test-Endpoint -Name "List All Engineers" -Method GET -Url "/engineers"
if ($engineers) {
    Write-Host "  Found: $($engineers.engineers.Count) engineers" -ForegroundColor Gray
    $testEngineer = $engineers.engineers[0]
    $testEngineerId = $testEngineer.id
    Write-Host "  Test Engineer: $($testEngineer.name) (ID: $testEngineerId)" -ForegroundColor Gray
}
Write-Host ""

# Test 1.2: Get Single Engineer
if ($testEngineerId) {
    $engineer = Test-Endpoint -Name "Get Engineer By ID" -Method GET -Url "/engineers/$testEngineerId"
    if ($engineer) {
        Write-Host "  Engineer: $($engineer.name)" -ForegroundColor Gray
        Write-Host "  Level: $($engineer.engineer_level)" -ForegroundColor Gray
        Write-Host "  Email: $($engineer.email)" -ForegroundColor Gray
    }
}
Write-Host ""

# Test 1.3: List Engineer Equipment Types (Before)
if ($testEngineerId) {
    $typesBefore = Test-Endpoint -Name "List Engineer Equipment Types" -Method GET -Url "/engineers/$testEngineerId/equipment-types"
    if ($typesBefore) {
        Write-Host "  Existing Capabilities: $($typesBefore.equipment_types.Count)" -ForegroundColor Gray
        if ($typesBefore.equipment_types.Count -gt 0) {
            $typesBefore.equipment_types | Select-Object -First 3 manufacturer, category | Format-Table -AutoSize
        }
    }
}
Write-Host ""

# Test 1.4: Add Equipment Type Capability
if ($testEngineerId) {
    $newCapability = @{
        manufacturer = "GE Healthcare"
        category = "Ultrasound"
    }
    $addResult = Test-Endpoint -Name "Add Equipment Type Capability" -Method POST -Url "/engineers/$testEngineerId/equipment-types" -Body $newCapability
    if ($addResult) {
        Write-Host "  Added: GE Healthcare Ultrasound" -ForegroundColor Gray
    }
}
Write-Host ""

# Test 1.5: List Engineer Equipment Types (After)
if ($testEngineerId) {
    $typesAfter = Test-Endpoint -Name "List Equipment Types (After Add)" -Method GET -Url "/engineers/$testEngineerId/equipment-types"
    if ($typesAfter) {
        Write-Host "  Total Capabilities: $($typesAfter.equipment_types.Count)" -ForegroundColor Gray
        $typesAfter.equipment_types | Format-Table manufacturer, category -AutoSize
    }
}
Write-Host ""

# ==============================================================================
# SECTION 2: ASSIGNMENT SUGGESTION & MANUAL ASSIGNMENT
# ==============================================================================
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "SECTION 2: Assignment Suggestion Algorithm" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""

# Test 2.1: Get/Create a Test Service Ticket
Write-Host "Preparing Test Ticket..." -ForegroundColor Yellow
try {
    $tickets = Invoke-WebRequest "$BaseUrl/tickets?limit=1" -UseBasicParsing | ConvertFrom-Json
    $testTicketId = $null

    if ($tickets.tickets -and $tickets.tickets.Count -gt 0) {
        $testTicketId = $tickets.tickets[0].id
        Write-Host "✓ Using existing ticket: $testTicketId" -ForegroundColor Green
    } else {
        Write-Host "✗ No tickets found - Create a test ticket first" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Failed to get tickets: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2.2: Get Engineer Suggestions for Ticket
if ($testTicketId) {
    $suggestions = Test-Endpoint -Name "Get Assignment Suggestions" -Method GET -Url "/tickets/$testTicketId/suggested-engineers"
    if ($suggestions -and $suggestions.suggested_engineers) {
        Write-Host "  Found: $($suggestions.suggested_engineers.Count) engineer suggestions" -ForegroundColor Gray
        
        if ($suggestions.suggested_engineers.Count -gt 0) {
            Write-Host ""
            Write-Host "  Top Suggestions:" -ForegroundColor Cyan
            $suggestions.suggested_engineers | Select-Object -First 3 priority, engineer_name, engineer_level, assignment_tier_name, match_reason | Format-Table -AutoSize
            
            $topSuggestion = $suggestions.suggested_engineers[0]
            
            # Test 2.3: Manually Assign Engineer
            Write-Host ""
            $assignmentRequest = @{
                engineer_id = $topSuggestion.engineer_id
                engineer_name = $topSuggestion.engineer_name
                organization_id = $topSuggestion.organization_id
                assignment_tier = $topSuggestion.assignment_tier
                assignment_tier_name = $topSuggestion.assignment_tier_name
            }
            
            $assignResult = Test-Endpoint -Name "Manual Engineer Assignment" -Method POST -Url "/tickets/$testTicketId/assign-engineer" -Body $assignmentRequest
            if ($assignResult) {
                Write-Host "  Assigned: $($topSuggestion.engineer_name) to ticket" -ForegroundColor Gray
            }
        }
    }
}
Write-Host ""

# ==============================================================================
# TEST SUMMARY
# ==============================================================================
Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host "   TEST SUMMARY" -ForegroundColor Cyan
Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host ""

$passed = ($TestResults | Where-Object { $_.Status -eq "PASSED" }).Count
$failed = ($TestResults | Where-Object { $_.Status -eq "FAILED" }).Count
$total = $TestResults.Count

Write-Host "Total Tests: $total" -ForegroundColor White
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -gt 0) { "Red" } else { "Green" })
Write-Host ""

if ($failed -gt 0) {
    Write-Host "Failed Tests:" -ForegroundColor Red
    $TestResults | Where-Object { $_.Status -eq "FAILED" } | Format-Table Test, Error -AutoSize
}

Write-Host "=======================================================" -ForegroundColor Cyan
Write-Host ""
