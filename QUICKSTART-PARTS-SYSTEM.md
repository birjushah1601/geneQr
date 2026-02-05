# ðŸš€ Parts Management System - Quick Start

## âš¡ 30-Second Setup

```powershell
# 1. Start Database (if not running)
cd dev/compose
docker-compose up -d postgres

# 2. Start Backend
cd C:\Users\birju\ServQR
$env:DB_HOST="localhost"; $env:DB_PORT="5430"; $env:DB_USER="postgres"; $env:DB_PASSWORD="postgres"; $env:DB_NAME="med_platform"; $env:ENABLE_ORG="true"
.\backend.exe

# 3. Start Frontend (new terminal)
cd admin-ui
npm run dev
```

## ðŸŽ¯ Test It Immediately

**Visit:** http://localhost:3000/parts-demo

**What You'll See:**
1. Sample MRI equipment
2. "Open Parts Browser" button
3. Click it â†’ Browse 16 real spare parts
4. Select parts, adjust quantities, see cost
5. Click "Assign" to complete

## ðŸ“Š What's Working

âœ… **16 Spare Parts** - Battery packs, filters, sensors, etc.  
âœ… **Real-time API** - Live data from backend  
âœ… **Smart Filtering** - Search, category, engineer requirements  
âœ… **Cost Calculator** - Automatic totaling  
âœ… **Engineer Detection** - Auto-identifies skill level needed  

## ðŸŽ¨ Key Features to Try

### 1. Search Parts
Type "battery" or "filter" in search box

### 2. Filter by Category
Select from: component, consumable, accessory, sensor, filter, battery

### 3. Engineer Filter
- Click "Needs Engineer" - Shows only parts requiring technician
- Click "Self-Service" - Shows user-serviceable parts

### 4. Multi-Select
Click multiple part cards or use checkboxes

### 5. Shopping Cart
- Switch to "Cart" tab
- Use +/- buttons to adjust quantities
- See real-time cost updates
- View engineer requirements summary

## ðŸ“± Screenshots

```
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚  Parts Assignment Demo                  â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”   â”‚
â”‚  â”‚ [Open Parts Browser]            â”‚   â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜   â”‚
â”‚                                         â”‚
â”‚  Sample Equipment: MRI Scanner          â”‚
â”‚  Model: GE Discovery MR750              â”‚
â”‚  Serial: MRI-3T-001                     â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜

After clicking "Open Parts Browser":

â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚  ðŸ” Assign Spare Parts                  â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”               â”‚
â”‚  â”‚ Browse   â”‚  Cart    â”‚               â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”´â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜               â”‚
â”‚  [Search...]  [Category â–¼]  [Filters]  â”‚
â”‚                                         â”‚
â”‚  â˜‘ Battery Pack - â‚¹350                 â”‚
â”‚  â–¡ Blood Tubing Set - â‚¹25              â”‚
â”‚  â–¡ Convex Probe - â‚¹9,500               â”‚
â”‚  â–¡ Detector Module - â‚¹25,000           â”‚
â”‚                                         â”‚
â”‚  [Cancel]          [Assign 1 Part]     â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

## ðŸ”¥ Pro Tips

1. **Filter Combo** - Use search + category + engineer filter together
2. **Cart Preview** - Switch to Cart tab to see summary before assigning
3. **Quantity Shortcuts** - Click +/- multiple times or click the part card again to remove
4. **Cost Tracking** - Total cost updates instantly as you adjust quantities

## ðŸ“ž Quick Troubleshooting

**Backend not starting?**
```powershell
# Check if database is running
docker ps | findstr med_platform_pg

# Verify connection
docker exec med_platform_pg psql -U postgres -d med_platform -c "SELECT COUNT(*) FROM spare_parts_catalog;"
```

**No parts showing?**
```powershell
# Test API directly
curl -H "X-Tenant-ID: default" http://localhost:8081/api/v1/catalog/parts
```

**Frontend errors?**
```powershell
# Check if backend is responding
curl http://localhost:8081/health
```

## ðŸŽ¯ What's Next?

1. **Test the UI** - Spend 5 minutes playing with it
2. **Check the data** - All 16 parts are real
3. **Review the code** - See `admin-ui/src/components/PartsAssignmentModal.tsx`
4. **Read full docs** - See `docs/PARTS-MANAGEMENT-COMPLETE.md`

## ðŸ“¦ System Status

| Component | Status | URL |
|-----------|--------|-----|
| Database | âœ… Running | localhost:5430 |
| Backend | âœ… Running | http://localhost:8081 |
| Frontend | âœ… Running | http://localhost:3000 |
| Parts API | âœ… Working | /api/v1/catalog/parts |
| Demo Page | âœ… Ready | /parts-demo |

## ðŸŽŠ You're All Set!

The system is **100% functional** for the core workflow:
- Browse parts âœ…
- Search & filter âœ…  
- Multi-select âœ…
- Calculate costs âœ…
- Detect engineer needs âœ…

**Go try it now:** http://localhost:3000/parts-demo

---

**Built with â¤ï¸ by Factory AI**
