# Onboarding System - Complete Brainstorming & Implementation Plan

**Date:** December 23, 2025  
**Status:** ðŸ§  **Brainstorming Phase**

---

## ðŸŽ¯ Overview

Comprehensive onboarding system for all organization types in the ServQR multi-tenant platform with CSV import capabilities, QR code generation, and data management.

---

## ðŸ“Š Current Schema Analysis

### **Core Tables Identified**

#### **1. Organizations (`organizations`)**
```sql
- id (UUID)
- name (TEXT)
- org_type (manufacturer|supplier|Channel Partner|Sub-sub_SUB_DEALER|hospital|service_provider|other)
- status (active|inactive|suspended)
- external_ref (TEXT) - For external system integration
- metadata (JSONB) - Flexible data storage
- created_at, updated_at
```

#### **2. Organization Relationships (`org_relationships`)**
```sql
- parent_org_id â†’ child_org_id
- rel_type (manufacturer_of|channel_partner_of|sub_sub_Sub-sub_SUB_DEALER_of|supplier_of|partner_of)
- metadata (JSONB)
```

#### **3. Organization Facilities (`organization_facilities`)**
```sql
- org_id (UUID)
- facility_name, facility_code
- facility_type (branch|department|lab|imaging_center|other)
- address (JSONB)
- status
```

#### **4. Equipment Catalog (`equipment_catalog`)**
```sql
- manufacturer_id (UUID â†’ organizations)
- equipment_type, model_number, model_name
- category (Diagnostic|Life Support|Surgical|Laboratory|Monitoring|Therapeutic|Imaging)
- specifications (JSONB)
- service_manual_url, user_manual_url, brochure_url
- image_urls (TEXT[])
- regulatory_approvals (JSONB)
- is_active, discontinued_date
```

#### **5. Equipment Parts (`equipment_parts`)**
```sql
- equipment_catalog_id (UUID)
- part_number, part_name
- part_category (consumable|replaceable|optional|tool)
- part_type (accessory|component|attachment|tool|supply)
- specifications (JSONB)
- is_oem, is_universal, is_critical
- standard_price, currency, lead_time_days
```

#### **6. Equipment Registry (`equipment_registry`)**
```sql
- id, qr_code, serial_number
- equipment_id (â†’ equipment_catalog)
- equipment_name, manufacturer_name, model_number
- customer_id, customer_name
- installation_location, installation_address (JSONB)
- installation_date, purchase_date, warranty_expiry
- status (operational|down|under_maintenance|decommissioned)
- qr_code_url
- specifications (JSONB), photos (JSONB), documents (JSONB)
```

#### **7. Contact Persons (`contact_persons`)**
```sql
- org_id (UUID â†’ organizations)
- contact_type (technical|billing|sales|support|management)
- name, email, phone, mobile
- designation, department
- is_primary
- metadata (JSONB)
```

#### **8. Engineers (`engineers`)**
```sql
- id (UUID)
- name, phone, email
- skills (TEXT[])
- home_region
- metadata (JSONB)
- org_id (service provider organization)
```

#### **9. Users (`users`)**
```sql
- id (UUID)
- email, phone
- password_hash
- full_name
- role (super_admin|org_admin|org_user|engineer|customer)
- status (active|inactive|suspended|pending)
- email_verified, phone_verified
```

#### **10. User Organizations (`user_organizations`)**
```sql
- user_id (UUID â†’ users)
- organization_id (UUID â†’ organizations)
- role (admin|manager|engineer|viewer)
- permissions (TEXT[])
- is_primary
- status
```

---

## ðŸ¢ Organization Type Onboarding Requirements

### **1. MANUFACTURER Onboarding**

#### **Data Required:**

**A. Organization Profile**
- Company name
- Legal entity details (GSTIN, PAN, CIN)
- Registration number
- Headquarters address
- Website, email, phone
- Company logo
- Year established
- Number of employees
- Annual revenue range
- Certifications (ISO, CE, FDA)

**B. Login Credentials**
- Admin user email
- Admin user full name
- Admin user phone
- Password (auto-generated or user-set)
- Role assignment (manufacturer_admin)
- Multi-factor authentication setup (optional)

**C. Equipment Types & Catalog**
- Equipment types manufactured
- Model numbers & names
- Category assignments
- Technical specifications
- Service manuals (PDFs)
- User manuals (PDFs)
- Product brochures
- Product images
- Regulatory approvals
- Compliance standards
- Typical lifespan
- Maintenance intervals
- Warranty terms

**D. Parts & Accessories Catalog**
- Part numbers
- Part names & descriptions
- Part categories (consumable/replaceable/optional)
- Compatible models
- Specifications
- Pricing (standard price list)
- Lead times
- Storage conditions
- Is OEM/Universal/Critical flags
- Replacement frequencies
- Stock availability

**E. Service Contacts**
- Technical support contacts
  - Name, email, phone
  - Areas of expertise
  - Available hours
- Billing/Finance contacts
  - Name, email, phone
  - Payment terms
  - Invoicing details
- Sales contacts (optional)
- Management contacts

**F. Facilities/Branches**
- Branch/facility locations
- Facility types (factory|warehouse|service_center|office)
- Full addresses with geolocations
- Contact details per facility
- Service coverage areas

**G. Equipment Installations (Optional)**
- List of equipment already installed in field
- Installation locations (hospital/clinic details)
- Serial numbers
- Installation dates
- Warranty status
- Contract details

**H. QR Code Management**
- Generate QR codes for new equipment
- QR codes WITHOUT location assignment (for stock/inventory)
- Bulk QR generation for batch production
- QR code format customization
- QR code branding (logo, colors)

---

### **2. HOSPITAL Onboarding**

#### **Data Required:**

**A. Organization Profile**
- Hospital name
- Hospital type (government|private|trust|corporate)
- Registration number
- Accreditations (NABH, JCI, ISO)
- Bed capacity
- Specializations
- Address, contact details
- Website, social media

**B. Login Credentials**
- Hospital admin email
- Biomedical engineer email(s)
- IT admin email
- Department heads (optional)

**C. Departments/Facilities**
- Department names (ICU, OT, Radiology, Lab, etc.)
- Department heads & contacts
- Department locations within hospital
- Equipment needs per department

**D. Installed Equipment Inventory**
- All equipment currently installed
- Equipment locations (department-wise)
- Purchase dates
- Warranty status
- AMC contracts
- Service history

**E. Preferred Manufacturers/Vendors**
- List of approved manufacturers
- Preferred service providers
- Contract terms
- SLA agreements

**F. Service Contacts**
- Biomedical engineering team
- IT department
- Procurement department
- Facilities management

---

### **3. SERVICE PROVIDER Onboarding**

#### **Data Required:**

**A. Organization Profile**
- Company name
- Service areas covered (geographic)
- Types of equipment serviced
- Certifications
- Years in business

**B. Login Credentials**
- Admin users
- Field engineers
- Back-office staff

**C. Engineers Database**
- Engineer names, contacts
- Skill sets
- Certifications
- Experience levels
- Equipment expertise
- Geographic coverage
- Availability schedules

**D. Service Coverage**
- Geographic territories
- Equipment types serviced
- Manufacturer partnerships
- Response time commitments
- 24/7 availability

**E. Contract Terms**
- Service rates
- Response times
- SLA commitments
- Payment terms
- Warranty policies

---

### **4. Channel Partner/Sub-sub_SUB_DEALER Onboarding**

#### **Data Required:**

**A. Organization Profile**
- Company details
- Distribution license
- Territory coverage
- Product lines carried

**B. Login Credentials**
- Sales team
- Technical support
- Warehouse staff

**C. Product Catalog**
- Manufacturer partnerships
- Product lines distributed
- Inventory locations
- Pricing agreements

**D. Service Capabilities**
- Installation services
- After-sales support
- Spare parts inventory
- Technical training

---

## ðŸ“¥ CSV Import System - Current State & Extensions

### **Current CSV Import Functionality**

#### **1. Equipment Registry Import** (`test-csv-import.ps1`)
**Endpoint:** `POST /api/v1/equipment/import`

**Current Fields Supported:**
```csv
serial_number,equipment_name,manufacturer_name,model_number,category,
customer_name,installation_location,installation_date,warranty_expiry,status
```

**Issues Identified:**
- âŒ No manufacturer_id field (needs to be added)
- âŒ No equipment_catalog_id field
- âŒ No organization_id for customer
- âŒ No QR code pre-generation
- âŒ No bulk QR code generation before assignment
- âŒ No photo/document upload via CSV
- âŒ No specifications in structured format

---

### **Required CSV Import Extensions**

#### **1. Equipment Catalog Bulk Import**
**New Endpoint:** `POST /api/v1/equipment-catalog/import`

**CSV Format:**
```csv
manufacturer_id,equipment_type,model_number,model_name,category,sub_category,
specifications_json,dimensions_json,weight_kg,power_requirements_json,
description,features_json,service_manual_url,user_manual_url,brochure_url,
image_urls_json,typical_lifespan_years,maintenance_interval_months,
requires_certification,regulatory_approvals_json,compliance_standards_json,
is_active
```

**Features:**
- Support JSONB fields via JSON strings
- Validate manufacturer_id exists
- Auto-generate UUID
- Handle image URLs (comma-separated in JSON array)
- Validate category enums
- Check for duplicate model numbers per manufacturer
- Support update mode (upsert)

#### **2. Equipment Parts Bulk Import**
**New Endpoint:** `POST /api/v1/equipment-parts/import`

**CSV Format:**
```csv
equipment_catalog_id,part_number,part_name,part_category,part_type,
description,specifications_json,dimensions_json,weight_kg,material,
compatible_models_json,replaces_part_number,is_oem,is_universal,is_critical,
lifespan_hours,replacement_frequency_months,unit_of_measure,min_order_quantity,
standard_price,currency,lead_time_days,storage_conditions,is_active
```

**Features:**
- Link parts to equipment catalog
- Support array fields (compatible_models)
- Validate equipment_catalog_id
- Price validation
- Stock availability tracking
- Support bulk updates

#### **3. Equipment Registry Import (Enhanced)**
**Endpoint:** `POST /api/v1/equipment/import` (Enhanced)

**New CSV Format:**
```csv
serial_number,qr_code,equipment_catalog_id,manufacturer_id,customer_org_id,
installation_location,installation_address_json,installation_date,
purchase_date,purchase_price,warranty_expiry,amc_contract_id,
status,specifications_json,photos_json,documents_json,notes
```

**Features:**
- Auto-generate QR codes if not provided
- Link to equipment_catalog_id
- Link to manufacturer organization
- Link to customer organization
- Support unassigned equipment (no customer_org_id = inventory)
- Generate QR code URL automatically
- Support photo/document URLs
- Validate status enum
- Set installation address as JSONB

#### **4. QR Code Bulk Generation**
**New Endpoint:** `POST /api/v1/qr-codes/bulk-generate`

**CSV Format:**
```csv
equipment_catalog_id,serial_number,quantity,batch_id,manufacturer_id
```

**Features:**
- Generate QR codes WITHOUT equipment assignment
- For new equipment batches/stock
- Support quantity-based generation
- Batch tracking
- Export generated QR codes to CSV
- Generate printable QR code sheets (PDF)
- Include manufacturer logo on QR codes
- Customizable QR code format

**Response:**
```json
{
  "total_generated": 100,
  "qr_codes": [
    {
      "qr_code": "QR-20251223-000001",
      "serial_number": "SN123456",
      "equipment_catalog_id": "uuid-here",
      "qr_code_url": "https://...",
      "qr_image_url": "https://..."
    }
  ],
  "pdf_url": "https://.../qr-batch-20251223.pdf"
}
```

#### **5. Organizations Bulk Import**
**New Endpoint:** `POST /api/v1/organizations/import`

**CSV Format:**
```csv
name,org_type,status,external_ref,gstin,pan,registration_number,
headquarters_address_json,website,primary_email,primary_phone,
year_established,employee_count,certifications_json,metadata_json
```

**Features:**
- Create organizations in bulk
- Support all org types
- Validate unique constraints
- Auto-generate UUID
- Support metadata JSONB
- Link relationships via separate CSV

#### **6. Organization Contacts Bulk Import**
**New Endpoint:** `POST /api/v1/organizations/{org_id}/contacts/import`

**CSV Format:**
```csv
org_id,contact_type,name,email,phone,mobile,designation,department,
is_primary,metadata_json
```

**Features:**
- Import multiple contacts per organization
- Validate contact types
- Ensure at least one primary contact
- Validate email/phone formats

#### **7. Engineers Bulk Import**
**New Endpoint:** `POST /api/v1/engineers/import`

**CSV Format:**
```csv
org_id,name,email,phone,mobile,skills_json,home_region,
certifications_json,experience_years,metadata_json
```

**Features:**
- Link to service provider organization
- Support skills array
- Create user accounts automatically
- Set engineer role and permissions
- Geographic coverage assignment

#### **8. Users Bulk Import**
**New Endpoint:** `POST /api/v1/users/bulk-import`

**CSV Format:**
```csv
email,phone,full_name,role,organization_id,org_role,permissions_json,
send_welcome_email,auto_generate_password
```

**Features:**
- Create users for multiple organizations
- Auto-generate passwords with email notification
- Set organization roles and permissions
- Support multiple organization memberships via separate import
- Email verification triggers

---

## ðŸ”„ Onboarding Workflows

### **Workflow 1: Manufacturer Onboarding (Complete)**

```
Step 1: Organization Creation
â”œâ”€â”€ Create organization record (type=manufacturer)
â”œâ”€â”€ Add headquarters facility
â”œâ”€â”€ Upload company documents
â””â”€â”€ Set status=pending_verification

Step 2: Admin User Creation
â”œâ”€â”€ Create admin user account
â”œâ”€â”€ Link to organization (org_admin role)
â”œâ”€â”€ Auto-generate password
â”œâ”€â”€ Send welcome email with setup link
â””â”€â”€ Enable MFA setup (optional)

Step 3: Equipment Catalog Import
â”œâ”€â”€ Prepare equipment catalog CSV
â”œâ”€â”€ Import equipment types via CSV
â”œâ”€â”€ Upload product images
â”œâ”€â”€ Upload manuals & brochures
â””â”€â”€ Verify and activate catalog items

Step 4: Parts Catalog Import
â”œâ”€â”€ Prepare parts catalog CSV
â”œâ”€â”€ Import parts linked to equipment
â”œâ”€â”€ Set pricing and lead times
â””â”€â”€ Mark critical parts

Step 5: Service Contacts Setup
â”œâ”€â”€ Import technical contacts CSV
â”œâ”€â”€ Import billing contacts CSV
â”œâ”€â”€ Set primary contacts
â””â”€â”€ Verify contact details

Step 6: Facilities Setup
â”œâ”€â”€ Import facilities/branches CSV
â”œâ”€â”€ Set geographic coverage
â””â”€â”€ Assign service territories

Step 7: QR Code Generation
â”œâ”€â”€ Bulk generate QR codes for equipment types
â”œâ”€â”€ Download QR code sheets (PDF)
â”œâ”€â”€ Print and apply to equipment
â””â”€â”€ Track QR code usage

Step 8: Equipment Registry (Optional)
â”œâ”€â”€ Import installed equipment CSV
â”œâ”€â”€ Link to equipment catalog
â”œâ”€â”€ Assign to customer organizations
â”œâ”€â”€ Set warranty and service dates
â””â”€â”€ Upload installation photos

Step 9: Verification & Activation
â”œâ”€â”€ Admin reviews all data
â”œâ”€â”€ Verifies certifications
â”œâ”€â”€ Approves manufacturer account
â””â”€â”€ Status â†’ active
```

### **Workflow 2: Hospital Onboarding**

```
Step 1: Organization Creation
â”œâ”€â”€ Create hospital organization
â”œâ”€â”€ Add hospital details
â”œâ”€â”€ Set accreditations
â””â”€â”€ Status=pending_verification

Step 2: Departments Setup
â”œâ”€â”€ Create facility records for departments
â”œâ”€â”€ Set department types
â””â”€â”€ Assign department heads

Step 3: Admin Users Setup
â”œâ”€â”€ Create biomedical engineer accounts
â”œâ”€â”€ Create department admin accounts
â”œâ”€â”€ Send welcome emails
â””â”€â”€ Set permissions

Step 4: Equipment Inventory Import
â”œâ”€â”€ Prepare equipment inventory CSV
â”œâ”€â”€ Import installed equipment
â”œâ”€â”€ Link to departments
â”œâ”€â”€ Set warranty/AMC status
â””â”€â”€ Upload equipment photos

Step 5: Preferred Vendors Setup
â”œâ”€â”€ Select preferred manufacturers
â”œâ”€â”€ Set service providers
â”œâ”€â”€ Define SLA agreements
â””â”€â”€ Configure notification preferences

Step 6: Verification & Activation
â”œâ”€â”€ Verify hospital credentials
â”œâ”€â”€ Review equipment inventory
â””â”€â”€ Activate account
```

### **Workflow 3: Service Provider Onboarding**

```
Step 1: Organization Creation
â”œâ”€â”€ Create service provider org
â”œâ”€â”€ Set service areas
â””â”€â”€ Upload certifications

Step 2: Engineers Bulk Import
â”œâ”€â”€ Prepare engineers CSV
â”œâ”€â”€ Import engineer profiles
â”œâ”€â”€ Create user accounts
â”œâ”€â”€ Set skills and expertise
â””â”€â”€ Assign territories

Step 3: Service Coverage Setup
â”œâ”€â”€ Define geographic territories
â”œâ”€â”€ Set equipment types serviced
â”œâ”€â”€ Configure response times
â””â”€â”€ Set availability schedules

Step 4: Contract Terms Setup
â”œâ”€â”€ Set service rates
â”œâ”€â”€ Define SLA commitments
â””â”€â”€ Configure payment terms

Step 5: Manufacturer Partnerships
â”œâ”€â”€ Link to manufacturers
â”œâ”€â”€ Set authorized service status
â””â”€â”€ Configure warranty servicing

Step 6: Activation
â””â”€â”€ Activate service provider account
```

---

## ðŸ› ï¸ Technical Implementation Plan

### **Phase 1: Schema Extensions** (Week 1)

**Tasks:**
1. Add missing fields to equipment_registry
   - manufacturer_id (UUID â†’ organizations)
   - organization_id for customer (UUID â†’ organizations)
   - equipment_catalog_id (UUID â†’ equipment_catalog)

2. Create organization_metadata table
   - Extended fields (GSTIN, PAN, CIN, etc.)
   - Certifications
   - Compliance documents

3. Add QR code batch tracking table
   ```sql
   CREATE TABLE qr_code_batches (
     id UUID PRIMARY KEY,
     manufacturer_id UUID REFERENCES organizations(id),
     batch_number TEXT,
     equipment_catalog_id UUID REFERENCES equipment_catalog(id),
     quantity_generated INT,
     generated_at TIMESTAMPTZ,
     pdf_url TEXT,
     status TEXT
   );
   ```

4. Add unassigned equipment table
   ```sql
   CREATE TABLE equipment_inventory (
     id UUID PRIMARY KEY,
     qr_code TEXT UNIQUE,
     serial_number TEXT UNIQUE,
     equipment_catalog_id UUID,
     manufacturer_id UUID,
     batch_id UUID,
     status TEXT (unassigned|reserved|assigned|sold),
     created_at TIMESTAMPTZ
   );
   ```

### **Phase 2: CSV Import Endpoints** (Week 2-3)

**Implement 8 new CSV import endpoints:**
1. Equipment Catalog Import
2. Equipment Parts Import
3. Enhanced Equipment Registry Import
4. QR Code Bulk Generation
5. Organizations Bulk Import
6. Organization Contacts Import
7. Engineers Bulk Import
8. Users Bulk Import

**Each endpoint needs:**
- CSV parsing with validation
- Transaction support (rollback on errors)
- Progress tracking
- Error reporting per row
- Success/failure counts
- Imported IDs return
- Duplicate detection and handling
- Update vs insert mode

### **Phase 3: QR Code System Enhancement** (Week 3-4)

**Features to implement:**
1. Bulk QR generation without equipment assignment
2. QR code batch management
3. PDF generation with multiple QR codes per page
4. QR code branding (manufacturer logo)
5. QR code format customization
6. Serial number linking
7. Equipment catalog linking
8. QR code usage tracking

### **Phase 4: Onboarding UI** (Week 4-5)

**Create onboarding wizard:**
1. Multi-step form for each organization type
2. CSV upload interface with preview
3. Drag-and-drop file uploads
4. Real-time validation
5. Progress tracking
6. Error display and correction
7. Bulk operation status
8. Success summary

### **Phase 5: Documentation & Testing** (Week 5-6)

**Deliverables:**
1. Onboarding guides per organization type
2. CSV template files with examples
3. API documentation for all endpoints
4. Video tutorials
5. Testing scripts
6. Sample data sets
7. Troubleshooting guides

---

## ðŸ“‹ CSV Templates Required

### **1. Equipment Catalog Template**
`equipment-catalog-template.csv`

### **2. Equipment Parts Template**
`equipment-parts-template.csv`

### **3. Equipment Registry Template** (Enhanced)
`equipment-registry-template.csv`

### **4. QR Code Generation Template**
`qr-bulk-generation-template.csv`

### **5. Organizations Template**
`organizations-template.csv`

### **6. Organization Contacts Template**
`organization-contacts-template.csv`

### **7. Engineers Template**
`engineers-template.csv`

### **8. Users Template**
`users-bulk-import-template.csv`

### **9. Facilities Template**
`organization-facilities-template.csv`

### **10. Organization Relationships Template**
`organization-relationships-template.csv`

---

## ðŸŽ¯ Priority Matrix

### **High Priority (Must Have)**
1. âœ… Equipment Catalog CSV Import
2. âœ… Equipment Parts CSV Import
3. âœ… Enhanced Equipment Registry Import
4. âœ… QR Code Bulk Generation
5. âœ… Organizations Bulk Import
6. âœ… Users Bulk Import

### **Medium Priority (Should Have)**
7. âš ï¸ Organization Contacts Import
8. âš ï¸ Engineers Bulk Import
9. âš ï¸ Facilities Import
10. âš ï¸ QR Code PDF Generation

### **Low Priority (Nice to Have)**
11. ðŸ”µ Organization Relationships Import
12. ðŸ”µ QR Code Branding/Customization
13. ðŸ”µ Advanced validation rules
14. ðŸ”µ Auto-detection of CSV format

---

## ðŸ“Š Success Metrics

### **Onboarding Efficiency**
- Time to complete manufacturer onboarding: < 30 minutes
- Time to complete hospital onboarding: < 20 minutes
- CSV import success rate: > 95%
- QR code generation speed: > 1000 codes/minute

### **Data Quality**
- Duplicate prevention rate: 100%
- Data validation error rate: < 5%
- Relationship integrity: 100%

### **User Experience**
- Wizard completion rate: > 90%
- Error recovery success: > 85%
- User satisfaction score: > 4.0/5.0

---

## ðŸš€ Next Steps

### **Immediate Actions:**
1. Review and approve brainstorming document
2. Prioritize features for MVP
3. Create detailed technical specifications
4. Design database schema changes
5. Create CSV template files
6. Develop API endpoints
7. Build onboarding UI
8. Create documentation
9. Test with real data
10. Deploy to production

---

## ðŸ“ Questions to Resolve

1. **QR Code Format:** What format do we want? (EAN-13, QR Code, Data Matrix?)
2. **Batch Size Limits:** Max equipment items per CSV import?
3. **File Size Limits:** Max CSV file size?
4. **Image Storage:** Where to store product images? (S3, local, CDN?)
5. **PDF Generation:** Use which library? (wkhtmltopdf, pdfkit, go-pdf?)
6. **Validation Rules:** How strict should data validation be?
7. **Duplicate Handling:** Update, skip, or error on duplicates?
8. **User Notifications:** Email/SMS on import completion?
9. **Audit Trail:** Log all bulk operations?
10. **Rollback:** Support transaction rollback on partial failures?

---

**Status:** Ready for technical specification and implementation planning  
**Estimated Total Effort:** 6 weeks (with 2 developers)  
**Next Document:** Technical Specification & API Design
