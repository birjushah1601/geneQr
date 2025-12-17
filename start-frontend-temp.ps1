Set-Location 'C:\Users\birju\aby-med\admin-ui'

Write-Host '🚀 Starting Frontend (Next.js)...' -ForegroundColor Green
Write-Host '   URL: http://localhost:3000' -ForegroundColor Cyan
Write-Host '   API: http://localhost:8082/api' -ForegroundColor Cyan
Write-Host ''

npm run dev
