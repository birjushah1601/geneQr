# Phase 2: Channel Partners → Channel Partners
$ErrorActionPreference = "Continue"

Write-Host ""
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host "███   Phase 2: Channel Partners → Channel Partners          ███" -ForegroundColor Cyan
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host ""

$replacements = @(
    @{ Find = "ChannelPartnerDashboard"; Replace = "ChannelPartnerDashboard" },
    @{ Find = "channel_partner_"; Replace = "channel_partner_" },
    @{ Find = "_channel_partner"; Replace = "_channel_partner" },
    @{ Find = "Channel Partners"; Replace = "Channel Partners" },
    @{ Find = "Channel Partner"; Replace = "Channel Partner" },
    @{ Find = "Channel Partners"; Replace = "channel_partners" },
    @{ Find = "Channel Partner"; Replace = "channel_partner" },
    @{ Find = "Channel Partner"; Replace = "CHANNEL_PARTNER" }
)

$fileExtensions = @("*.go", "*.ts", "*.tsx", "*.js", "*.jsx", "*.sql", "*.md", "*.json", "*.ps1", "*.txt", "*.html", "*.csv")
$excludeDirs = @(".git", "node_modules", ".next", "dist", "build", "vendor", "bin", "storage")

function Should-Exclude {
    param($path)
    foreach ($dir in $excludeDirs) {
        if ($path -like "*\$dir\*") { return $true }
    }
    return $false
}

$allFiles = @()
foreach ($ext in $fileExtensions) {
    $allFiles += Get-ChildItem -Path . -Filter $ext -Recurse -File -ErrorAction SilentlyContinue | Where-Object { -not (Should-Exclude $_.FullName) }
}

$updated = 0
foreach ($file in $allFiles) {
    try {
        $content = Get-Content $file.FullName -Raw -ErrorAction Stop
        $originalContent = $content
        
        foreach ($r in $replacements) {
            $content = $content -replace [regex]::Escape($r.Find), $r.Replace
        }
        
        if ($content -ne $originalContent) {
            $content | Out-File $file.FullName -Encoding UTF8 -NoNewline
            $updated++
            if ($updated % 5 -eq 0) {
                Write-Host "  Updated $updated files..." -ForegroundColor Gray
            }
        }
    } catch {
        # Skip files that can't be read
    }
}

Write-Host ""
Write-Host "✅ Phase 2 Complete: Channel Partners → Channel Partners" -ForegroundColor Green
Write-Host "   Updated: $updated files" -ForegroundColor White
Write-Host ""
