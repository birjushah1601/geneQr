# Phase 3: Sub-Sub-sub_sub_SUB_DEALERs → Sub-Sub-Sub-sub_sub_SUB_DEALERs
$ErrorActionPreference = "Continue"

Write-Host ""
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host "███   Phase 3: Sub-Sub-sub_sub_SUB_DEALERs → Sub-Sub-Sub-sub_sub_SUB_DEALERs                    ███" -ForegroundColor Cyan
Write-Host "████████████████████████████████████████████████████████████" -ForegroundColor Cyan
Write-Host ""

$replacements = @(
    @{ Find = "sub_sub_Sub-sub_SUB_DEALER_"; Replace = "sub_sub_sub_Sub-sub_SUB_DEALER_" },
    @{ Find = "_sub_Sub-sub_SUB_DEALER"; Replace = "_sub_sub_Sub-sub_SUB_DEALER" },
    @{ Find = "Sub-Sub-sub_sub_SUB_DEALERs"; Replace = "Sub-Sub-Sub-sub_sub_SUB_DEALERs" },
    @{ Find = "Sub-sub_SUB_DEALER"; Replace = "Sub-Sub-sub_SUB_DEALER" },
    @{ Find = "Sub-Sub-sub_sub_SUB_DEALERs"; Replace = "sub_sub_Sub-Sub-sub_sub_SUB_DEALERs" },
    @{ Find = "Sub-sub_SUB_DEALER"; Replace = "sub_sub_Sub-sub_SUB_DEALER" },
    @{ Find = "Sub-sub_SUB_DEALER"; Replace = "SUB_sub_Sub-sub_SUB_DEALER" }
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
Write-Host "✅ Phase 3 Complete: Sub-Sub-sub_sub_SUB_DEALERs → Sub-Sub-Sub-sub_sub_SUB_DEALERs" -ForegroundColor Green
Write-Host "   Updated: $updated files" -ForegroundColor White
Write-Host ""
