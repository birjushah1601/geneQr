# Quick Start Testing Guide

## Current Status ✅

Your ABY-MED platform is **running successfully**! Here's what's confirmed:

1. ✅ Health endpoint responding: `{"status":"ok"}`
2. ✅ All 8 modules loaded and initialized
3. ✅ HTTP server listening on port 8081
4. ✅ Infrastructure services (PostgreSQL, Kafka, Redis) healthy

## ⚠️ Database Schema Issue

The database tables haven't been created yet. You're seeing errors like:
- `relation "rfqs" does not exist`
- `Failed to retrieve equipment list`

This is normal for a fresh installation!

## Solution: Initialize Database Schema

### Option 1: Use SQL Initialization Scripts (Recommended)

Check if initialization scripts exist:

```powershell
# Check for SQL migration files
Get-ChildItem -Path internal -Filter "*.sql" -Recurse
Get-ChildItem -Path dev/compose/postgres -Filter "*.sql" -Recurse
```

### Option 2: Manual Database Setup

Connect to PostgreSQL and create tables:

```powershell
# Connect to PostgreSQL
docker exec -it med-platform-postgres psql -U postgres -d medplatform

# Or from Windows:
psql -h localhost -p 5433 -U postgres -d medplatform
# Password: postgres
```

### Option 3: Simple Testing Without Full Schema

Let me create a minimal schema creation script for you:

```sql
-- Create basic tables for testing
CREATE SCHEMA IF NOT EXISTS public;

-- RFQ Tables
CREATE TABLE IF NOT EXISTS rfqs (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    deadline TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_rfqs_tenant ON rfqs(tenant_id);

-- Catalog Tables
CREATE TABLE IF NOT EXISTS catalog_items (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(200),
    manufacturer VARCHAR(200),
    description TEXT,
    base_price DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'INR',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_catalog_tenant ON catalog_items(tenant_id);

-- Supplier Tables
CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    company_name VARCHAR(500) NOT NULL,
    email VARCHAR(200),
    phone VARCHAR(50),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_suppliers_tenant ON suppliers(tenant_id);

-- Equipment Registry
CREATE TABLE IF NOT EXISTS equipment (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(500) NOT NULL,
    serial_number VARCHAR(200) UNIQUE,
    model VARCHAR(200),
    manufacturer VARCHAR(200),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_equipment_tenant ON equipment(tenant_id);

-- Service Tickets
CREATE TABLE IF NOT EXISTS service_tickets (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    equipment_id VARCHAR(255),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    priority VARCHAR(50) DEFAULT 'medium',
    status VARCHAR(50) DEFAULT 'open',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (equipment_id) REFERENCES equipment(id)
);

CREATE INDEX idx_tickets_tenant ON service_tickets(tenant_id);
CREATE INDEX idx_tickets_equipment ON service_tickets(equipment_id);
```

## Testing Without Database (Monitoring Only)

While the database is being set up, you can test these endpoints:

### 1. Health Check ✅ Working
```bash
curl http://localhost:8081/health
```

**Response:**
```json
{"status":"ok"}
```

### 2. Prometheus Metrics ✅ Working
```bash
curl http://localhost:8081/metrics
```

### 3. Monitoring Dashboards

- **Grafana:** http://localhost:3000 (admin/admin)
- **Prometheus:** http://localhost:9090
- **MailHog:** http://localhost:8025

## Complete Testing Script

I'll create a PowerShell script to test all services once the database is initialized:

```powershell
# test-all-services.ps1
Write-Host "`n=== ABY-MED Platform Testing Suite ===" -ForegroundColor Green

$baseUrl = "http://localhost:8081"
$tenant = "city-hospital"

# Test 1: Health Check
Write-Host "`n1. Testing Health Endpoint..." -ForegroundColor Cyan
$health = Invoke-RestMethod -Uri "$baseUrl/health"
Write-Host "   Status: $($health.status)" -ForegroundColor Green

# Test 2: List RFQs
Write-Host "`n2. Testing RFQ Service..." -ForegroundColor Cyan
try {
    $rfqs = Invoke-RestMethod -Uri "$baseUrl/api/v1/rfq" -Headers @{"X-Tenant-ID"=$tenant}
    Write-Host "   RFQs found: $($rfqs.Count)" -ForegroundColor Green
} catch {
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 3: List Catalog
Write-Host "`n3. Testing Catalog Service..." -ForegroundColor Cyan
try {
    $catalog = Invoke-RestMethod -Uri "$baseUrl/api/v1/catalog" -Headers @{"X-Tenant-ID"=$tenant}
    Write-Host "   Catalog items: $($catalog.Count)" -ForegroundColor Green
} catch {
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 4: List Suppliers
Write-Host "`n4. Testing Supplier Service..." -ForegroundColor Cyan
try {
    $suppliers = Invoke-RestMethod -Uri "$baseUrl/api/v1/suppliers" -Headers @{"X-Tenant-ID"=$tenant}
    Write-Host "   Suppliers found: $($suppliers.Count)" -ForegroundColor Green
} catch {
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 5: List Equipment
Write-Host "`n5. Testing Equipment Registry..." -ForegroundColor Cyan
try {
    $equipment = Invoke-RestMethod -Uri "$baseUrl/api/v1/equipment" -Headers @{"X-Tenant-ID"=$tenant}
    Write-Host "   Equipment registered: $($equipment.Count)" -ForegroundColor Green
} catch {
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 6: List Service Tickets
Write-Host "`n6. Testing Service Ticket..." -ForegroundColor Cyan
try {
    $tickets = Invoke-RestMethod -Uri "$baseUrl/api/v1/tickets" -Headers @{"X-Tenant-ID"=$tenant}
    Write-Host "   Tickets found: $($tickets.Count)" -ForegroundColor Green
} catch {
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "`n=== Testing Complete ===" -ForegroundColor Green
```

## Next Steps

1. **Initialize Database Schema**
   - Run the SQL script above in PostgreSQL
   - Or check for existing migration scripts in the project

2. **Import Postman Collection**
   - File: `ABY-MED-Postman-Collection.json`
   - Import into Postman
   - Update environment variables if needed

3. **Run Tests**
   - Execute the PowerShell test script
   - Or use Postman collection
   - Or use curl commands from API-TESTING-GUIDE.md

4. **Create Sample Data**
   - Start with Catalog items
   - Register Suppliers
   - Create RFQs
   - Test complete workflows

## Available Resources

1. **API-TESTING-GUIDE.md** - Complete API documentation with curl examples
2. **ABY-MED-Postman-Collection.json** - Ready-to-import Postman collection
3. **Platform logs** - Check terminal where platform is running
4. **Service logs** - Docker compose logs for infrastructure

## Getting Help

If you encounter issues:
1. Check platform logs in the terminal
2. Verify database connection: `docker exec -it med-platform-postgres psql -U postgres -d medplatform -c "\\dt"`
3. Check service health: `docker compose -p med-platform ps`
4. Review error messages in the API responses

## Alternative Testing Tools

### Using HTTPie (Prettier than curl)
```bash
# Install: pip install httpie
http GET http://localhost:8081/health
http GET http://localhost:8081/api/v1/rfq X-Tenant-ID:city-hospital
```

### Using Insomnia
- Similar to Postman
- Import our Postman collection (Insomnia can read it)

### Using VS Code REST Client
- Install "REST Client" extension
- Create `.http` files with requests

Happy Testing! 🚀
