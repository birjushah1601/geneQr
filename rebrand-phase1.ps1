# ServQR Rebranding Script
$ErrorActionPreference = "Continue"

Write-Host ""
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host "███         ServQR REBRANDING SCRIPT                    ███" -ForegroundColor Cyan
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host ""

$replacements = @(
    # ServQR variations
    @{ Find = "ServQR Platform"; Replace = "ServQR Platform" },
    @{ Find = "ServQR Platform"; Replace = "ServQR Platform" },
    @{ Find = "ServQR Platform"; Replace = "ServQR Platform" },
    @{ Find = "servqr-platform"; Replace = "servqr-platform" },
    @{ Find = "ServQR"; Replace = "ServQR" },
    @{ Find = "ServQR"; Replace = "ServQR" },
    @{ Find = "ServQR"; Replace = "ServQR" },
    @{ Find = "ServQR"; Replace = "ServQR" },
    @{ Find = "ServQR"; Replace = "servqr" },
    @{ Find = "ServQR"; Replace = "servqr" }
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
            if ($updated % 10 -eq 0) {
                Write-Host "  Updated $updated files..." -ForegroundColor Gray
            }
        }
    } catch {
        # Skip files that can't be read
    }
}

Write-Host ""
Write-Host "✅ Phase 1 Complete: ServQR → ServQR" -ForegroundColor Green
Write-Host "   Updated: $updated files" -ForegroundColor White
Write-Host ""
