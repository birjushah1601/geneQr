# ðŸ§ª LIVE TESTING GUIDE - Real Data from Database

## âœ… CONFIRMED: REAL DATA READY

**Database Status:**
- âœ… 16 spare parts in `spare_parts_catalog` table
- âœ… 3 bundles in `spare_parts_bundles` table  
- âœ… 2 suppliers in `spare_parts_suppliers` table
- âœ… Backend API serving real data from PostgreSQL

**Sample Real Parts:**
```
1. Battery Pack Rechargeable    - â‚¹350   - component
2. Blood Tubing Set             - â‚¹25    - consumable
3. Convex Array Probe           - â‚¹9,500 - accessory
4. Detector Module 16-slice     - â‚¹25,000 - component
5. Flat Panel Detector          - â‚¹45,000 - component
6. Head Coil 8-Channel          - â‚¹12,500 - accessory
7. X-Ray Tube Assembly          - â‚¹65,000 - component
... +9 more parts
```

---

## ðŸš€ STEP 1: Start Frontend

**Open a PowerShell terminal and run:**
```powershell
cd C:\Users\birju\ServQR\admin-ui
npm run dev
```

**Wait for:**
```
âœ“ Ready in 3.2s
â—‹ Local: http://localhost:3000
```

---

## ðŸŽ¯ STEP 2: Test Parts Demo Page (Standalone)

### 2A. Open Demo Page
```
http://localhost:3000/parts-demo
```

### 2B. What You'll See
- Sample MRI equipment details
- "Open Parts Browser" button
- Empty parts list

### 2C. Click "Open Parts Browser"
âœ… **Modal opens with REAL DATA from database:**
- 16 parts loaded from API
- Real prices from `spare_parts_catalog` table
- Real categories (component, consumable, accessory, etc.)

### 2D. Test Features

**Search:**
- Type "battery" â†’ Should show "Battery Pack Rechargeable" (â‚¹350)
- Type "filter" â†’ Should show filter-related parts
- Type "probe" â†’ Should show "Convex Array Probe" (â‚¹9,500)

**Category Filter:**
- Select "component" â†’ Shows 6 parts (Battery, X-Ray Tube, Detector, etc.)
- Select "consumable" â†’ Shows consumables (Blood Tubing, etc.)
- Select "accessory" â†’ Shows accessories (Probes, Coils, etc.)

**Engineer Filter:**
- Toggle "Needs Engineer" â†’ Shows parts requiring technician
- Toggle "Self-Service" â†’ Shows user-serviceable parts

**Add to Cart:**
- Click on "Battery Pack Rechargeable" card â†’ Selected âœ…
- Click on "Blood Tubing Set" â†’ Selected âœ…
- Switch to "Cart" tab

**Adjust Quantities:**
- Battery Pack: Click + to increase (1 â†’ 2 â†’ 3)
- Blood Tubing: Click - to decrease
- See total cost update in real-time

**Cost Calculation:**
```
Battery Pack: 2x @ â‚¹350 = â‚¹700
Blood Tubing: 5x @ â‚¹25 = â‚¹125
Total: â‚¹825
```

**Engineer Detection:**
- If you select parts with `requires_engineer=true`
- See "Engineer Required: L2" or "L3" in summary

**Click "Assign":**
- Parts added to equipment
- Summary shows on main page
- Cost displayed: "2 parts assigned â€¢ â‚¹825"

---

## ðŸŽ¯ STEP 3: Test Service Request Integration (Full Workflow)

### 3A. Open Service Request Page
```
http://localhost:3000/service-request?qr=HOSP001-CT001
```

### 3B. Fill Service Request Form

**Equipment Details** (Auto-loaded from QR):
- Should show equipment info

**Your Name:**
```
John Doe
```

**Priority:**
```
Select: High
```

**Issue Description:**
```
CT Scanner showing error code E203. Not starting up. 
Suspected power supply issue.
```

### 3C. Add Parts to Request

**Look for green section:**
```
ðŸ“¦ Spare Parts Needed
Select spare parts needed for this service request
[Add Parts]
```

**Click "Add Parts" button:**
- âœ… Modal opens with 16 REAL parts from database
- Same parts as demo page

**Select Parts:**
1. Search "battery" â†’ Select "Battery Pack Rechargeable"
2. Search "detector" â†’ Select "Detector Module 16-slice"
3. Switch to Cart tab
4. Battery: Set quantity to 2
5. Detector: Keep quantity at 1

**See Real-Time Calculation:**
```
Battery Pack: 2x @ â‚¹350 = â‚¹700
Detector Module: 1x @ â‚¹25,000 = â‚¹25,000
Total Cost: â‚¹25,700
```

**Click "Assign":**
- Modal closes
- Parts added to service request

### 3D. Verify Parts Summary

**In the green "Spare Parts Needed" section:**
```
ðŸ“¦ Spare Parts Needed
2 parts assigned â€¢ â‚¹25,700

â€¢ Battery Pack Rechargeable - 2x â€¢ â‚¹700
â€¢ Detector Module 16-slice - 1x â€¢ â‚¹25,000

[Modify Parts]
```

### 3E. Submit Service Request

**Click "Submit Service Request" button:**
- âœ… Request submitted successfully
- Success message appears
- Shows equipment details
- Shows "Create Another Request" button

---

## ðŸ“Š STEP 4: Verify Data Flow

### 4A. Check API Response
**Open PowerShell and run:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8081/api/v1/catalog/parts" -Headers @{"X-Tenant-ID"="default"} | ConvertTo-Json -Depth 3
```

**Should return 16 parts with:**
- `id` - UUID from database
- `part_name` - Real part names
- `part_number` - Part numbers
- `category` - component/consumable/accessory
- `unit_price` - Real prices from DB
- `requires_engineer` - Boolean flag
- `engineer_level_required` - L1/L2/L3

### 4B. Check Database Directly
```powershell
docker exec med_platform_pg psql -U postgres -d med_platform -c "SELECT part_name, unit_price, category FROM spare_parts_catalog LIMIT 5;"
```

**Should show exact same data as API**

---

## ðŸŽ¨ STEP 5: Advanced Features Test

### Filter Combinations
1. Search "pack" + Category "component" â†’ Battery Pack
2. Search "tube" + Category "consumable" â†’ Blood Tubing
3. Engineer "Needs Engineer" + Category "component" â†’ High-skill parts

### Multi-Select Test
1. Select 5 different parts
2. Adjust quantities (1-10 each)
3. See total cost climb to â‚¹50,000+
4. Remove 2 parts
5. See cost recalculate

### Edge Cases
1. Search with no results: "xyz123"
2. Add 0 parts and try to assign
3. Select same part twice (shouldn't duplicate)
4. Adjust quantity to 99

---

## âœ… EXPECTED RESULTS

### Parts Demo Page
âœ… 16 parts loaded from database  
âœ… Search works across part names  
âœ… Filters work (category, engineer)  
âœ… Multi-select with checkboxes  
âœ… Cart shows selected parts  
âœ… Quantity adjusters work (+/-)  
âœ… Cost calculates in real-time  
âœ… Engineer level detected automatically  
âœ… Assign button adds parts  

### Service Request Page
âœ… Form loads equipment details  
âœ… "Add Parts" button visible (green section)  
âœ… Clicking opens parts modal  
âœ… Modal loads same 16 real parts  
âœ… Parts can be selected and assigned  
âœ… Summary displays in request form  
âœ… Cost shown: "X parts â€¢ â‚¹Y"  
âœ… Parts list preview (first 3)  
âœ… "Modify Parts" works to reopen  
âœ… Submit includes parts data  

---

## ðŸ› TROUBLESHOOTING

### Modal Shows "No parts available"
**Problem:** API not returning data  
**Fix:**
```powershell
# Check backend is running
curl http://localhost:8081/api/v1/catalog/parts
```

### Parts list is empty
**Problem:** Database has no parts  
**Fix:**
```powershell
# Check database
docker exec med_platform_pg psql -U postgres -d med_platform -c "SELECT COUNT(*) FROM spare_parts_catalog;"
# Should show: 16
```

### "Failed to fetch parts" error
**Problem:** Backend not running or wrong port  
**Fix:**
```powershell
# Verify backend running on 8081
netstat -ano | findstr ":8081"

# Restart if needed
.\backend.exe
```

### Modal doesn't open
**Problem:** React component error  
**Fix:**
```powershell
# Check browser console (F12)
# Look for errors in Components tab
```

---

## ðŸ“ˆ SUCCESS METRICS

After testing, you should have:

âœ… **Seen 16 real parts** from spare_parts_catalog table  
âœ… **Filtered by category** (component, consumable, accessory)  
âœ… **Searched parts** by name  
âœ… **Added to cart** with quantities  
âœ… **Calculated costs** in real-time (â‚¹)  
âœ… **Detected engineer levels** (L1/L2/L3)  
âœ… **Assigned parts** to service request  
âœ… **Viewed summary** with cost and parts list  
âœ… **Modified assignments** after initial selection  
âœ… **Submitted request** with all data  

---

## ðŸŽ¯ TEST SCENARIOS

### Scenario 1: Simple Repair
**Parts:** Battery Pack (2x), Blood Tubing (5x)  
**Cost:** â‚¹825  
**Engineer:** Not required  
**Use Case:** Quick user-serviceable fix

### Scenario 2: Major Component Replacement
**Parts:** X-Ray Tube Assembly (1x), Detector Module (1x)  
**Cost:** â‚¹90,000  
**Engineer:** L3 required  
**Use Case:** Complex repair needing expert

### Scenario 3: Routine Maintenance
**Parts:** Filter Element (2x), Head Coil (1x), Battery (1x)  
**Cost:** â‚¹13,800  
**Engineer:** L1 or L2  
**Use Case:** Scheduled preventive maintenance

---

## ðŸŽŠ FINAL CHECKLIST

Before considering testing complete:

- [ ] Frontend running on :3000
- [ ] Backend running on :8081
- [ ] Database accessible with 16 parts
- [ ] Parts demo page loads
- [ ] Modal opens and shows parts
- [ ] Search/filter works
- [ ] Cart functionality works
- [ ] Cost calculation accurate
- [ ] Parts assigned to service request
- [ ] Summary displayed correctly
- [ ] Service request submits successfully

---

**Ready to test? Start with Step 1!** ðŸš€

**Questions or issues? Check TROUBLESHOOTING section above.**
