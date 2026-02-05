# CSV Import Templates

This folder contains sample CSV templates for bulk importing data into the ServQR Platform.

## ðŸ“‹ Available Templates

### 1. Equipment Import
**File:** `equipment-catalog-template.csv`

**Fields (13 columns):**
```csv
serial_number, equipment_name, manufacturer_name, model_number, category,
customer_name, customer_id, installation_location, installation_date,
purchase_date, purchase_price, warranty_months, notes
```

**Required Fields:**
- `serial_number` - Unique equipment identifier
- `equipment_name` - Full equipment name
- `manufacturer_name` - Manufacturer/brand name
- `customer_name` - Hospital/customer name

**Example:**
```csv
SN-MRI-2024-001,MRI Machine - MAGNETOM Vida,Siemens Healthineers,MAGNETOM Vida,Radiology,Apollo Hospital Delhi,CUST-APO-DEL-001,Radiology Department - Room 201,2024-01-15,2023-12-10,45000000,36,3 Tesla MRI Scanner
```

---

### 2. Engineers Import
**File:** `engineers-import-template.csv`

**Fields (7 columns):**
```csv
name, phone, email, location, engineer_level, equipment_types, experience_years
```

**Required Fields:**
- `name` - Engineer full name
- `phone` - Contact number (format: +919876543210)
- `email` - Email address

**Engineer Levels:**
- `1` - Junior Engineer
- `2` - Senior Engineer
- `3` - Expert Engineer

**Equipment Types:**
- Pipe-separated list: `MRI|CT Scanner|X-Ray`
- Multiple types supported

**Example:**
```csv
Rajesh Kumar,+919876543210,rajesh.kumar@company.com,Delhi NCR,3,MRI|CT Scanner|X-Ray,12
```

**Note:** Backend implementation for CSV import is pending. Currently, engineers must be created through the UI.

---

### 3. Parts Catalog Import
**File:** `parts-catalog-template.csv`

**Fields (16 columns):**
```csv
part_number, part_name, description, category, subcategory, part_type,
is_oem_part, manufacturer_name, unit_price, currency, minimum_stock,
lead_time_days, weight_kg, dimensions, warranty_months, specifications
```

**Required Fields:**
- `part_number` - Unique part identifier
- `part_name` - Part name

**Part Types:**
- `Component` - Replaceable component
- `Consumable` - Consumable item
- `Accessory` - Accessory/attachment

**Example:**
```csv
VENT-FILT-001,HEPA Filter H13,High-efficiency particulate air filter,Filter,Air Filter,Consumable,true,Medtronic,450.00,INR,20,7,0.15,10x10x2 cm,6,"{""efficiency"": ""99.97%""}"
```

---

### 4. Team Members Import
**File:** `team-members-template.csv`

**Fields (4 columns):**
```csv
name, email, role, phone
```

**Roles:**
- `admin` - Full access
- `manager` - Management access
- `viewer` - Read-only access

**Example:**
```csv
Rajesh Kumar,ceo@company.com,admin,+919876543210
```

**Note:** Team members can also be added via invitation system (recommended).

---

### 5. Organizations Import
**File:** `organizations-import-template.csv`

For importing multiple hospitals/customers at once.

---

## ðŸ”§ Specialized Equipment Templates

Category-specific templates with pre-filled equipment types:

- `equipment-catalog-cardiology-template.csv` - Cardiology equipment
- `equipment-catalog-radiology-template.csv` - Radiology equipment
- `equipment-catalog-icu-template.csv` - ICU equipment
- `equipment-catalog-surgical-template.csv` - Surgical equipment
- `equipment-catalog-laboratory-template.csv` - Laboratory equipment

---

## ðŸ“ Usage Instructions

### Step 1: Download Template
Download the appropriate template from this folder.

### Step 2: Fill Data
- Open in Excel, Google Sheets, or any CSV editor
- Fill in your data following the format
- Keep column order exactly as shown
- Don't remove or rename column headers

### Step 3: Validate Data
- **Required fields** must not be empty
- **Dates** use format: `YYYY-MM-DD` (e.g., 2024-01-15)
- **Phone numbers** include country code (e.g., +919876543210)
- **Email addresses** must be valid format
- **Numeric fields** use numbers only (no currency symbols)

### Step 4: Import
1. Log into ServQR Platform
2. Navigate to appropriate section (Equipment, Engineers, etc.)
3. Click "Import CSV" or "Bulk Upload"
4. Select your filled CSV file
5. Review preview and import

---

## âš ï¸ Common Issues & Solutions

### Issue: Import Fails
**Solution:** Check that:
- Column headers match exactly (case-sensitive)
- No extra columns added
- Required fields are not empty
- Date format is YYYY-MM-DD
- No special characters in IDs

### Issue: Data Not Appearing
**Solution:**
- Refresh the page
- Check if import completed successfully
- Review error messages in import results

### Issue: Duplicate Entries
**Solution:**
- Use unique values for ID fields (serial_number, part_number, etc.)
- Check existing data before importing

---

## ðŸ”— Related Documentation

- **Full Template Review:** See `../CSV-TEMPLATE-REVIEW.md` for complete analysis
- **Backend API:** Equipment import endpoint: `POST /api/v1/equipment/import`
- **Engineer Import:** Pending backend implementation
- **Parts Import:** Backend verification needed

---

## ðŸ“ž Support

For issues with CSV imports:
1. Check error messages carefully
2. Validate data format matches template
3. Review backend logs for detailed errors
4. Contact support with:
   - CSV file sample
   - Error message
   - Number of rows attempted

---

## ðŸ“… Template Version

**Last Updated:** January 26, 2026  
**Version:** 2.0  
**Breaking Changes:** Equipment template format corrected (v2.0)

---

## ðŸŽ¯ Quick Start

**New to CSV imports?**

1. Start with **equipment-catalog-template.csv**
2. Add 2-3 sample rows
3. Test import in development environment
4. Review results
5. Scale up to full data set

**Best Practice:** Always test with small sample before importing thousands of rows!
