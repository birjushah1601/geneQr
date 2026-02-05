Set-Location 'C:\Users\birju\ServQR'

# Database Configuration
$env:DB_HOST = 'localhost'
$env:DB_PORT = '5430'
$env:DB_USER = 'postgres'
$env:DB_PASSWORD = 'postgres'
$env:DB_NAME = 'med_platform'

# Storage Configuration
$env:STORAGE_BASE_PATH = 'C:\Users\birju\ServQR\storage'

# Module Configuration
$env:ENABLED_MODULES = '*'
$env:ENABLE_ORG = 'true'
$env:ENABLE_ORG_SEED = 'true'
$env:PORT = '8081'

# AI Configuration (from .env file)
$env:AI_PROVIDER = 'openai'
$env:OPENAI_API_KEY = 'your-openai-api-key-here'
$env:OPENAI_ORG_ID = ''
$env:OPENAI_MODEL = 'gpt-4-vision-preview'
$env:OPENAI_MAX_TOKENS = '4000'
$env:OPENAI_TEMPERATURE = '0.3'

Write-Host '🚀 Starting Backend API...' -ForegroundColor Green
Write-Host "   Port: $env:PORT" -ForegroundColor Gray
Write-Host "   DB: $env:DB_USER@$env:DB_HOST:$env:DB_PORT/$env:DB_NAME" -ForegroundColor Gray
Write-Host "   AI Provider: $env:AI_PROVIDER" -ForegroundColor Gray
Write-Host "   OpenAI Key: Configured ✓" -ForegroundColor Gray
Write-Host ''

.\backend.exe

