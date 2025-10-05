# Quick Reference - ABY-MED API Testing

## 🚀 Quick Start

### Test Equipment Registry (✅ Working)
```bash
curl -H "X-Tenant-ID: city-hospital" http://localhost:8081/api/v1/equipment
```

### Test Health
```bash
curl http://localhost:8081/health
```

---

## 📡 All API Endpoints

### Base URL
```
http://localhost:8081/api/v1
```

### Common Headers
```
X-Tenant-ID: city-hospital
X-User-ID: test-user
Content-Type: application/json
```

---

## 🎯 Service Endpoints Quick Reference

| Service | Method | Endpoint | Status |
|---------|--------|----------|--------|
| **Health** | GET | `/health` | ✅ |
| **Metrics** | GET | `/metrics` | ✅ |
| **Equipment** | GET | `/equipment` | ✅ |
| **Equipment** | POST | `/equipment` | ✅ |
| **RFQ** | GET | `/rfq` | ⚠️ |
| **RFQ** | POST | `/rfq` | ⚠️ |
| **Catalog** | GET | `/catalog` | ⚠️ |
| **Catalog** | POST | `/catalog` | ⚠️ |
| **Suppliers** | GET | `/suppliers` | ⚠️ |
| **Suppliers** | POST | `/suppliers` | ⚠️ |
| **Quotes** | GET | `/quotes` | 🔄 |
| **Quotes** | POST | `/quotes` | 🔄 |
| **Tickets** | GET | `/tickets` | 🔄 |
| **Tickets** | POST | `/tickets` | 🔄 |

**Legend:** ✅ Working | ⚠️ Schema issue | 🔄 Not tested yet

---

## 💻 PowerShell Testing Commands

### Test All Services
```powershell
# Health
Invoke-RestMethod http://localhost:8081/health

# Equipment (Working)
Invoke-RestMethod -Uri http://localhost:8081/api/v1/equipment -Headers @{"X-Tenant-ID"="city-hospital"}

# RFQ
Invoke-RestMethod -Uri http://localhost:8081/api/v1/rfq -Headers @{"X-Tenant-ID"="city-hospital"}

# Catalog
Invoke-RestMethod -Uri http://localhost:8081/api/v1/catalog -Headers @{"X-Tenant-ID"="city-hospital"}

# Suppliers
Invoke-RestMethod -Uri http://localhost:8081/api/v1/suppliers -Headers @{"X-Tenant-ID"="city-hospital"}
```

### Create Equipment
```powershell
$body = @{
    name = "New MRI Scanner"
    serial_number = "MRI-2025-001"
    model = "MAGNETOM Skyra"
    manufacturer = "Siemens"
    status = "active"
} | ConvertTo-Json

Invoke-RestMethod -Method POST `
    -Uri http://localhost:8081/api/v1/equipment `
    -Headers @{"X-Tenant-ID"="city-hospital";"Content-Type"="application/json"} `
    -Body $body
```

---

## 🌐 Web Dashboards

| Dashboard | URL | Credentials | Purpose |
|-----------|-----|-------------|---------|
| Grafana | http://localhost:3000 | admin/admin | Monitoring |
| Prometheus | http://localhost:9090 | - | Metrics |
| MailHog | http://localhost:8025 | - | Email testing |

---

## 📦 Import Postman Collection

1. Open Postman
2. Click **Import**
3. Select file: `ABY-MED-Postman-Collection.json`
4. Collection variables already set:
   - `baseUrl`: http://localhost:8081
   - `tenantId`: city-hospital
   - `userId`: test-user

---

## 🔍 Debugging Commands

### Check Platform Status
```powershell
# Platform process
Get-Process *platform* | Select-Object Name, Id, CPU, WorkingSet

# Check if port 8081 is listening
Test-NetConnection -ComputerName localhost -Port 8081
```

### Check Docker Services
```powershell
cd dev/compose
docker compose -p med-platform ps
docker compose -p med-platform logs platform --tail=20
docker compose -p med-platform logs postgres --tail=20
```

### Database Commands
```powershell
# List all tables
docker exec -it med-platform-postgres psql -U postgres -d medplatform -c "\dt"

# Check sample data
docker exec -it med-platform-postgres psql -U postgres -d medplatform -c "SELECT COUNT(*) FROM catalog_items;"

# Connect to database
docker exec -it med-platform-postgres psql -U postgres -d medplatform
```

---

## 📝 Sample Test Scenarios

### Scenario 1: List Equipment (✅ Working)
```bash
curl -H "X-Tenant-ID: city-hospital" \
     http://localhost:8081/api/v1/equipment
```

**Expected:** JSON with equipment list

### Scenario 2: Create RFQ (⚠️ Schema issue)
```bash
curl -X POST http://localhost:8081/api/v1/rfq \
  -H "X-Tenant-ID: city-hospital" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Need MRI Scanner",
    "description": "Urgent requirement",
    "deadline": "2025-12-31T00:00:00Z"
  }'
```

**Current Status:** Schema mismatch error

### Scenario 3: List Catalog Items (⚠️ Schema issue)
```bash
curl -H "X-Tenant-ID: city-hospital" \
     http://localhost:8081/api/v1/catalog
```

**Current Status:** Schema mismatch error

---

## 🛠️ Troubleshooting Guide

### Issue: "relation does not exist"
**Solution:** Database schema mismatch. Check:
1. What tables the code expects
2. What tables actually exist
3. Run migrations if available

### Issue: Connection refused on port 8081
**Solution:** 
```powershell
# Check if platform is running
Get-Process *platform*

# Restart if needed
.\bin\platform.exe
```

### Issue: Empty response
**Solution:**
1. Check platform logs in terminal
2. Verify tenant ID header is set
3. Check database has data

---

## 📁 Files Created for You

| File | Purpose |
|------|---------|
| `API-TESTING-GUIDE.md` | Complete API documentation |
| `ABY-MED-Postman-Collection.json` | Postman collection |
| `QUICK-START-TESTING.md` | Step-by-step guide |
| `init-database-schema.sql` | Database schema (executed) |
| `TESTING-STATUS-SUMMARY.md` | Current status report |
| `QUICK-REFERENCE-TESTING.md` | This file |

---

## 🎯 Testing Priorities

1. **First**: Test Equipment Registry (already working)
2. **Second**: Fix schema issues for RFQ, Catalog, Supplier
3. **Third**: Test Quote, Contract, Comparison services
4. **Fourth**: Test Service Tickets
5. **Fifth**: Test complete workflows end-to-end

---

## ✨ Quick Wins

Things you can test RIGHT NOW:

```bash
# 1. Health check
curl http://localhost:8081/health

# 2. List equipment (will return data!)
curl -H "X-Tenant-ID: city-hospital" http://localhost:8081/api/v1/equipment

# 3. View metrics
curl http://localhost:8081/metrics | grep "go_"

# 4. Open Grafana
start http://localhost:3000

# 5. Open Prometheus
start http://localhost:9090
```

---

## 📞 Need Help?

1. Check `TESTING-STATUS-SUMMARY.md` for detailed status
2. Check `API-TESTING-GUIDE.md` for complete API docs
3. Check platform logs in terminal
4. Check Docker logs: `docker compose -p med-platform logs`

---

**Last Updated:** October 1, 2025  
**Platform Version:** 0.1.0  
**Status:** 60% Services Tested ✅
