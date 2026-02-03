# Fix corrupted Unicode characters and rebrand GeneQR to ServQR

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  Unicode Fix & ServQR Rebranding Script" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

# Unicode replacements map
$unicodeReplacements = @{
    # Corrupted arrows
    "â†" = "←"
    "â†'" = "→"
    "â†�" = "↑"
    "â†"" = "↓"
    
    # Corrupted checkmarks and symbols
    "âœ"" = "✓"
    "âœ…" = "✅"
    "âœ–" = "✗"
    "âœ¨" = "✨"
    "â­�" = "⭐"
    "â˜…" = "★"
    "âš ï¸�" = "⚠️"
    "â„¹ï¸�" = "ℹ️"
    
    # Corrupted quotes and punctuation
    "â€œ" = '"'
    "â€�" = '"'
    "â€˜" = "'"
    "â€™" = "'"
    "â€"" = "—"
    "â€"" = "–"
    "â€¦" = "…"
    
    # Corrupted emojis (common ones)
    "ðŸ'‹" = "👋"
    "ðŸ"§" = "📧"
    "ðŸ"±" = "📱"
    "ðŸ"" = "📝"
    "ðŸ"Š" = "📊"
    "ðŸ"ˆ" = "📈"
    "ðŸš€" = "🚀"
    "ðŸ'¼" = "💼"
    "ðŸ'¡" = "💡"
    "âœ‰ï¸�" = "✉️"
    
    # Corrupted spaces
    "Â " = " "
    "Â" = ""
    
    # Unicode replacement character
    "ï¿½" = ""
}

# GeneQR rebranding
$rebrandReplacements = @{
    "GeneQR" = "ServQR"
    "genq-admin-ui" = "servqr-admin-ui"
    "genq" = "servqr"
    "GENQ" = "SERVQR"
}

$fixedFiles = 0
$totalReplacements = 0

# Get all source files
$files = Get-ChildItem -Path "src" -Recurse -Include "*.tsx","*.ts","*.jsx","*.js","*.json"

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Encoding UTF8 | Out-String
    $originalContent = $content
    $fileReplacements = 0
    
    # Fix Unicode
    foreach ($key in $unicodeReplacements.Keys) {
        if ($content -match [regex]::Escape($key)) {
            $count = ([regex]::Matches($content, [regex]::Escape($key))).Count
            $content = $content -replace [regex]::Escape($key), $unicodeReplacements[$key]
            $fileReplacements += $count
        }
    }
    
    # Rebrand GeneQR
    foreach ($key in $rebrandReplacements.Keys) {
        if ($content -match [regex]::Escape($key)) {
            $count = ([regex]::Matches($content, [regex]::Escape($key))).Count
            $content = $content -replace [regex]::Escape($key), $rebrandReplacements[$key]
            $fileReplacements += $count
        }
    }
    
    # If changes were made, save the file
    if ($content -ne $originalContent) {
        $utf8NoBom = New-Object System.Text.UTF8Encoding $false
        [System.IO.File]::WriteAllText($file.FullName, $content, $utf8NoBom)
        $fixedFiles++
        $totalReplacements += $fileReplacements
        Write-Host "✓ Fixed: $($file.Name) ($fileReplacements replacements)" -ForegroundColor Green
    }
}

# Also fix package.json
$packageJson = "package.json"
if (Test-Path $packageJson) {
    $content = Get-Content $packageJson -Encoding UTF8 | Out-String
    $originalContent = $content
    
    foreach ($key in $rebrandReplacements.Keys) {
        if ($content -match [regex]::Escape($key)) {
            $content = $content -replace [regex]::Escape($key), $rebrandReplacements[$key]
        }
    }
    
    if ($content -ne $originalContent) {
        $utf8NoBom = New-Object System.Text.UTF8Encoding $false
        [System.IO.File]::WriteAllText((Resolve-Path $packageJson).Path, $content, $utf8NoBom)
        Write-Host "✓ Fixed: package.json" -ForegroundColor Green
        $fixedFiles++
    }
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  Summary" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host "Files fixed: $fixedFiles" -ForegroundColor Cyan
Write-Host "Total replacements: $totalReplacements" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ Done!" -ForegroundColor Green
Write-Host ""
