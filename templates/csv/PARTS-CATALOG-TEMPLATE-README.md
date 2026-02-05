# Parts Catalog CSV Template - Field Guide

## Template File
ðŸ“ **File:** `parts-catalog-template.csv`

## Purpose
This template helps you prepare your parts catalog data for bulk upload into the ServQR Platform. Use this format to import spare parts for your medical equipment.

---

## ðŸ“‹ Field Descriptions

### Required Fields (*)

| Field | Type | Description | Example | Required |
|-------|------|-------------|---------|----------|
| **part_number** * | Text | Unique identifier for the part | VENT-FILT-001 | âœ… Yes |
| **part_name** * | Text | Display name of the part | HEPA Filter H13 | âœ… Yes |
| **description** | Text | Detailed description | High-efficiency particulate air filter... | No |
| **category** * | Text | Main category | Filter, Sensor, Valve, Battery, Cable, Circuit | âœ… Yes |
| **subcategory** | Text | Subcategory | Air Filter, Gas Sensor, Pressure Valve | No |
| **part_type** * | Text | Type of part | Component, Consumable, Accessory, Tool | âœ… Yes |
| **is_oem_part** | Boolean | Is this an OEM part? | true, false | No (default: true) |
| **manufacturer_name** | Text | Manufacturer name | Medtronic, GE Healthcare, Philips | No |
| **unit_price** | Decimal | Price per unit | 450.00 | No |
| **currency** | Text | Currency code (ISO 4217) | INR, USD, EUR | No (default: INR) |
| **minimum_stock** | Integer | Minimum inventory level | 20 | No |
| **lead_time_days** | Integer | Days to receive part | 7 | No |
| **weight_kg** | Decimal | Weight in kilograms | 0.15 | No |
| **dimensions** | Text | Physical dimensions | 10x10x2 cm | No |
| **warranty_months** | Integer | Warranty period in months | 6, 12, 24 | No |
| **specifications** | JSON | Technical specifications (JSON format) | See below | No |

---

## ðŸŽ¯ Part Categories

**Common Categories:**
- Filter (Air Filter, Oil Filter, Water Filter)
- Sensor (Gas Sensor, Flow Sensor, Temperature Sensor, Pressure Sensor)
- Valve (Pressure Valve, Safety Valve, Control Valve)
- Battery (Power Supply, Backup Battery)
- Cable (Power Cable, Signal Cable, Gas Connection)
- Circuit (Breathing Circuit, Electronic Circuit)
- Display (LCD, LED, Touchscreen)
- Tubing (Patient Tubing, Gas Line)
- Mask (Patient Interface, Ventilation Mask)
- Adapter (Connector, Y-Connector, Coupling)

---

## ðŸ“¦ Part Types

| Type | Description | Examples |
|------|-------------|----------|
| **Component** | Permanent parts that are replaced when broken | Sensors, Valves, Batteries, Circuit Boards |
| **Consumable** | Parts that need regular replacement | Filters, Tubing, Masks, Circuits |
| **Accessory** | Optional or enhancement parts | Humidifiers, Stands, Covers, Adapters |
| **Tool** | Tools needed for maintenance | Calibration tools, Wrenches, Test kits |

---

## ðŸ”§ Technical Specifications (JSON Format)

The `specifications` field accepts JSON format for technical details:

```json
{
  "efficiency": "99.97%",
  "particle_size": "0.3 microns",
  "flow_rate": "60 L/min",
  "pressure_range": "0-100 psi",
  "temperature_range": "31-41Â°C",
  "material": "silicone",
  "sterilizable": true,
  "autoclavable": true
}
```

**Common Specification Keys:**
- `efficiency` - Performance percentage
- `range` - Operating range
- `accuracy` - Measurement accuracy
- `response_time` - Sensor response time
- `capacity` - Battery/fluid capacity
- `voltage` - Electrical voltage
- `material` - Construction material
- `size` - Physical size
- `weight` - Item weight
- `sterilizable` - Can be sterilized (true/false)
- `autoclavable` - Can be autoclaved (true/false)

---

## ðŸ“ Example Rows

### Example 1: Filter (Consumable)
```csv
VENT-FILT-001,HEPA Filter H13,High-efficiency particulate air filter for ventilator breathing circuit,Filter,Air Filter,Consumable,true,Medtronic,450.00,INR,20,7,0.15,"10x10x2 cm",6,"{""efficiency"": ""99.97%"", ""particle_size"": ""0.3 microns"", ""flow_rate"": ""60 L/min""}"
```

### Example 2: Sensor (Component)
```csv
VENT-SENS-002,Oxygen Sensor,Galvanic oxygen sensor for O2 monitoring,Sensor,Gas Sensor,Component,true,GE Healthcare,2500.00,INR,10,14,0.05,"5x3x2 cm",12,"{""range"": ""0-100%"", ""accuracy"": ""Â±2%"", ""response_time"": ""<15 sec""}"
```

### Example 3: Battery (Component)
```csv
VENT-BATT-005,Backup Battery Pack,Lithium-ion battery for emergency backup,Battery,Power Supply,Component,true,Siemens,4500.00,INR,5,21,1.20,"20x10x5 cm",24,"{""capacity"": ""2200mAh"", ""voltage"": ""14.4V"", ""runtime"": ""4 hours""}"
```

---

## âœ… Validation Rules

1. **part_number**: Must be unique across all parts
2. **part_name**: Required, max 200 characters
3. **category**: Required, should match common categories
4. **part_type**: Must be one of: Component, Consumable, Accessory, Tool
5. **is_oem_part**: true or false (case-insensitive)
6. **unit_price**: Positive decimal number (2 decimal places)
7. **currency**: 3-letter ISO code (INR, USD, EUR, etc.)
8. **minimum_stock**: Positive integer
9. **lead_time_days**: Positive integer
10. **weight_kg**: Positive decimal number
11. **warranty_months**: Positive integer
12. **specifications**: Valid JSON format (use double quotes)

---

## ðŸš€ How to Use

### Step 1: Download Template
Download `parts-catalog-template.csv` from the platform

### Step 2: Fill in Your Data
- Open in Excel or Google Sheets
- Fill in your parts information
- Use the examples as reference
- Ensure part_number is unique

### Step 3: Validate Data
- Check all required fields are filled
- Verify JSON format in specifications column
- Confirm prices and dimensions are correct
- Review part categorization

### Step 4: Save as CSV
- File â†’ Save As â†’ CSV (Comma delimited)
- Keep the header row intact
- Use UTF-8 encoding

### Step 5: Upload to Platform
- Go to AI Onboarding Wizard or Parts Management
- Click "Upload Parts Catalog"
- Select your CSV file
- Review the preview
- Confirm import

---

## ðŸ’¡ Tips & Best Practices

### Naming Conventions
- **Part Numbers**: Use consistent prefix (e.g., VENT-FILT-001, VENT-SENS-002)
- **Part Names**: Clear, descriptive names with model if applicable
- **Categories**: Stick to standard categories for better organization

### Pricing
- Include your standard selling price
- Use consistent currency (INR recommended for India)
- Update prices regularly

### Inventory Management
- Set realistic minimum_stock levels
- Account for lead_time_days for reordering
- Consider criticality of parts

### Technical Details
- Add comprehensive specifications
- Include model numbers in description
- Reference OEM part numbers when available

### Quality Data
- Double-check all measurements
- Verify compatibility information
- Keep warranty information current
- Include datasheet URLs if available

---

## âŒ Common Mistakes to Avoid

1. **Duplicate Part Numbers**: Each part_number must be unique
2. **Invalid JSON**: Use double quotes in specifications, not single quotes
3. **Missing Required Fields**: part_number, part_name, category, part_type are mandatory
4. **Wrong Part Type**: Must be exactly: Component, Consumable, Accessory, or Tool
5. **Currency Mismatch**: Use consistent currency throughout
6. **Negative Values**: Prices, weights, and quantities must be positive
7. **Excel Formatting**: Save as CSV, not XLSX
8. **Special Characters**: Avoid special characters in part_number

---

## ðŸ“ž Support

Need help with the template?
- Check the platform documentation
- Contact support team
- Refer to existing parts in the system for examples

---

## ðŸ”„ Template Version
**Version:** 1.0  
**Last Updated:** January 2026  
**Compatible With:** ServQR Platform v2.0+

---

**Ready to import your parts catalog?** ðŸš€  
Follow the steps above and your parts data will be ready in minutes!
