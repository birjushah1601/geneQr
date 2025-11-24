# 🏭 Manufacturer Onboarding Guide

**Version:** 1.0.0  
**Date:** November 17, 2024  
**Status:** Ready for Implementation  
**Owner:** GeneQR Platform Team

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Pre-Onboarding Checklist](#pre-onboarding-checklist)
3. [Required Information](#required-information)
4. [Onboarding Flow](#onboarding-flow)
5. [Data Collection Forms](#data-collection-forms)
6. [Validation & Approval Process](#validation--approval-process)
7. [Post-Onboarding Activities](#post-onboarding-activities)
8. [Technical Integration](#technical-integration)
9. [Support & Training](#support--training)

---

## 🎯 Overview

### What is Manufacturer Onboarding?

Manufacturer onboarding is the process of registering and activating equipment manufacturers in the GeneQR platform so that:

✅ **Service Providers** can access manufacturer documentation, parts catalogs, and support  
✅ **Hospitals** can link equipment to manufacturers for warranty and service tracking  
✅ **Parts Suppliers** can identify authorized parts and compatibility  
✅ **AI Services** can access manufacturer knowledge bases for better diagnosis  

### Benefits for Manufacturers

- 📊 **Real-time visibility** into equipment performance across installations
- 🔧 **Better service quality** through platform integration
- 📈 **Data-driven insights** on equipment issues and patterns
- 💰 **Revenue opportunities** through parts sales and service contracts
- 🤝 **Stronger relationships** with service providers and hospitals

### Types of Manufacturers

The platform supports different manufacturer types:

1. **OEM (Original Equipment Manufacturer)** - Primary equipment manufacturer
2. **Component Manufacturer** - Makes specific components for equipment
3. **Contract Manufacturer** - Manufactures under license
4. **Private Label Manufacturer** - Branded equipment from other manufacturers

---

## ✅ Pre-Onboarding Checklist

**Before starting the onboarding process, verify the manufacturer has:**

### 📄 Legal & Business Requirements

- [ ] **Valid business registration** - Certificate/registration number
- [ ] **Tax identification** - GST/Tax ID number
- [ ] **Business license** - License to manufacture medical equipment
- [ ] **ISO certifications** - ISO 13485 (Medical Devices Quality Management)
- [ ] **FDA/CE certifications** - If applicable for their equipment
- [ ] **Company PAN card** (India)
- [ ] **Bank account details** - For potential revenue sharing

### 🏢 Organizational Requirements

- [ ] **Primary contact person** - Authorized representative
- [ ] **Technical support contact** - For escalations
- [ ] **Parts/warranty contact** - For parts catalog and warranty queries
- [ ] **Legal/compliance contact** - For contracts and agreements
- [ ] **Corporate address** - Registered business address
- [ ] **Website** - Company website (optional but recommended)

### 📚 Documentation Requirements

- [ ] **Equipment catalog** - List of all equipment models manufactured
- [ ] **Parts catalog** - Parts list with SKUs and specifications
- [ ] **Service manuals** - Technical documentation for each equipment
- [ ] **Warranty policies** - Standard warranty terms and conditions
- [ ] **Service guidelines** - Recommended maintenance schedules
- [ ] **Training materials** - For technicians (optional)

### 🔧 Technical Requirements

- [ ] **Equipment specifications** - Technical specs for each model
- [ ] **Variant information** - Different variants/configurations per model
- [ ] **Parts compatibility matrix** - Which parts work with which models
- [ ] **Error code dictionary** - Common error codes and meanings
- [ ] **Service bulletins** - Known issues and fixes (optional)
- [ ] **API access** (optional) - If manufacturer has APIs for real-time data

### 💰 Commercial Requirements

- [ ] **Parts pricing** - Standard parts pricing (if selling through platform)
- [ ] **Service rates** - Recommended service labor rates
- [ ] **Warranty coverage** - What's covered and for how long
- [ ] **Payment terms** - If participating in platform commerce
- [ ] **Revenue sharing agreement** (optional) - For platform-facilitated sales

---

## 📝 Required Information

### Section 1: Basic Organization Information

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Organization Name** | Text | ✅ Yes | Legal business name |
| **Brand Name(s)** | Text | ⚠️ Conditional | If different from org name |
| **Organization Type** | Dropdown | ✅ Yes | "manufacturer" |
| **Manufacturer Subtype** | Dropdown | ✅ Yes | OEM / Component / Contract / Private Label |
| **Registration Number** | Text | ✅ Yes | Business registration ID |
| **GST/Tax ID** | Text | ✅ Yes | Tax identification |
| **PAN Number** | Text | ✅ Yes | India specific |
| **Incorporation Date** | Date | ⚠️ Recommended | Company establishment date |
| **Website** | URL | ⚠️ Recommended | Company website |
| **Logo** | Image | ⚠️ Recommended | Company logo (for UI) |

### Section 2: Address Information

**Corporate/Registered Address:**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Address Line 1** | Text | ✅ Yes | |
| **Address Line 2** | Text | ⬜ Optional | |
| **City** | Text | ✅ Yes | |
| **State/Province** | Text | ✅ Yes | |
| **Postal Code** | Text | ✅ Yes | |
| **Country** | Dropdown | ✅ Yes | |
| **Landmark** | Text | ⬜ Optional | |

**Manufacturing Facilities:**  
(Can add multiple facilities - see Facilities section)

### Section 3: Contact Information

**Primary Contact:**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Full Name** | Text | ✅ Yes | Authorized representative |
| **Designation** | Text | ✅ Yes | Job title |
| **Email** | Email | ✅ Yes | Primary email (must be verified) |
| **Phone** | Phone | ✅ Yes | Mobile with country code |
| **Alternate Phone** | Phone | ⬜ Optional | |

**Additional Contacts:**
- Technical Support Contact
- Parts/Warranty Contact
- Legal/Compliance Contact
- Accounts/Finance Contact

(Each with same fields as primary)

### Section 4: Certifications & Compliance

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **ISO 13485 Certified** | Yes/No | ✅ Yes | Medical devices quality mgmt |
| **ISO 13485 Certificate Number** | Text | ⚠️ Conditional | If yes above |
| **ISO 13485 Valid Until** | Date | ⚠️ Conditional | Expiry date |
| **FDA Registered** | Yes/No | ⬜ Optional | US market |
| **FDA Registration Number** | Text | ⚠️ Conditional | If yes |
| **CE Marked** | Yes/No | ⬜ Optional | EU market |
| **CE Certificate Number** | Text | ⚠️ Conditional | If yes |
| **Other Certifications** | Text Array | ⬜ Optional | List other certs |

**Document Uploads:**
- [ ] ISO 13485 Certificate (PDF)
- [ ] FDA Registration (PDF) - if applicable
- [ ] CE Certificate (PDF) - if applicable
- [ ] Business License (PDF)
- [ ] GST Registration Certificate (PDF)

### Section 5: Banking Information

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Bank Name** | Text | ⚠️ Conditional | If participating in commerce |
| **Account Holder Name** | Text | ⚠️ Conditional | Must match business name |
| **Account Number** | Text | ⚠️ Conditional | |
| **IFSC Code** | Text | ⚠️ Conditional | India specific |
| **Bank Branch** | Text | ⚠️ Conditional | |
| **Account Type** | Dropdown | ⚠️ Conditional | Current / Savings |
| **Cancelled Cheque** | Image/PDF | ⚠️ Conditional | For verification |

### Section 6: Equipment Catalog

For each equipment model manufactured:

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Equipment Category** | Dropdown | ✅ Yes | Ventilator, MRI, CT Scanner, etc. |
| **Model Name** | Text | ✅ Yes | |
| **Model Number** | Text | ✅ Yes | Unique identifier |
| **Model Code/SKU** | Text | ✅ Yes | Internal SKU |
| **Description** | Text | ✅ Yes | |
| **Equipment Image** | Image | ⚠️ Recommended | Product photo |
| **Year Introduced** | Year | ⚠️ Recommended | Market introduction year |
| **Is Active** | Yes/No | ✅ Yes | Currently manufactured? |
| **Warranty Period** | Number | ✅ Yes | In months |
| **Expected Lifespan** | Number | ⚠️ Recommended | In years |
| **Service Manual** | PDF | ✅ Yes | Technical documentation |

**Specifications (JSONB):**
```json
{
  "technical": {
    "dimensions": "120cm x 80cm x 90cm",
    "weight": "150 kg",
    "power_requirements": "220V, 50Hz, 15A",
    "operating_temperature": "15-30°C",
    "humidity_range": "30-75%"
  },
  "features": [
    "Feature 1",
    "Feature 2"
  ],
  "compliance": {
    "standards": ["IEC 60601", "ISO 13485"]
  }
}
```

### Section 7: Equipment Variants

Each model can have multiple variants (e.g., ICU vs General Ward):

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Variant Name** | Text | ✅ Yes | E.g., "ICU Grade", "Standard" |
| **Variant Code** | Text | ✅ Yes | Unique code |
| **Base Model** | Dropdown | ✅ Yes | Link to equipment model |
| **Description** | Text | ✅ Yes | How it differs from base |
| **Price Difference** | Number | ⬜ Optional | vs base model |
| **Specification Overrides** | JSONB | ⬜ Optional | Variant-specific specs |

**Example:**
- **Base Model:** Ventilator V-100
- **Variants:**
  - V-100-ICU (ICU Grade with advanced monitoring)
  - V-100-STD (Standard for general wards)
  - V-100-PED (Pediatric variant)

### Section 8: Parts Catalog

For each part:

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Part Name** | Text | ✅ Yes | |
| **Part Number/SKU** | Text | ✅ Yes | Manufacturer part number |
| **Part Category** | Dropdown | ✅ Yes | Filter, Sensor, Battery, etc. |
| **Description** | Text | ✅ Yes | |
| **Compatible Equipment Models** | Multi-select | ✅ Yes | Which models use this part |
| **Compatible Variants** | Multi-select | ⚠️ Conditional | Variant-specific parts |
| **Unit of Measure** | Text | ✅ Yes | Piece, Pair, Set, etc. |
| **Recommended Stock Level** | Number | ⬜ Optional | Min inventory |
| **Lead Time** | Number | ⬜ Optional | Days to procure |
| **Expected Lifespan** | Number | ⬜ Optional | In hours/cycles |
| **List Price** | Number | ⬜ Optional | MSRP |
| **Part Image** | Image | ⚠️ Recommended | Photo of part |
| **Installation Instructions** | PDF | ⬜ Optional | How to install |

### Section 9: Parts Compatibility Matrix

**Related/Accessory Parts:**
- Filter → Filter Housing Seal (always replace together)
- Battery → Battery Connector Cable (optional accessory)

**Variant-Specific Parts:**
- V-100-ICU requires HEPA-V100-ICU filter (higher grade)
- V-100-STD uses HEPA-V100-STD filter (standard grade)

### Section 10: Service & Warranty Information

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Standard Warranty Period** | Number | ✅ Yes | In months |
| **Extended Warranty Available** | Yes/No | ⬜ Optional | |
| **Extended Warranty Period** | Number | ⚠️ Conditional | Additional months |
| **Warranty Terms & Conditions** | Text/PDF | ✅ Yes | What's covered |
| **Recommended Service Interval** | Number | ⬜ Optional | In months |
| **Service Manual URL** | URL | ⬜ Optional | Link to online docs |
| **Technical Support Email** | Email | ✅ Yes | |
| **Technical Support Phone** | Phone | ✅ Yes | |
| **Support Hours** | Text | ⚠️ Recommended | E.g., "24/7" or "9AM-6PM IST" |

### Section 11: Error Codes & Troubleshooting

For each equipment model, provide common error codes:

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| **Error Code** | Text | ✅ Yes | E.g., "E-42" |
| **Error Message** | Text | ✅ Yes | What displays |
| **Description** | Text | ✅ Yes | What the error means |
| **Severity** | Dropdown | ✅ Yes | Critical / High / Medium / Low |
| **Probable Causes** | Text Array | ✅ Yes | List of possible causes |
| **Troubleshooting Steps** | Text | ✅ Yes | Step-by-step guide |
| **Parts Typically Required** | Multi-select | ⬜ Optional | Link to parts |
| **Estimated Repair Time** | Number | ⬜ Optional | In minutes |

**Example:**
```
Error Code: E-42
Message: "Filter Warning"
Description: HEPA filter has reached end of life
Severity: Medium
Probable Causes:
  - Filter clogged with particles
  - Filter housing seal damaged
  - Air pressure sensor malfunction
Troubleshooting:
  1. Check filter for visible damage
  2. Inspect housing seal for cracks
  3. Test air pressure sensor
  4. Replace filter if needed
Parts Required: HEPA-V100, SEAL-V100-KIT (optional)
Est. Repair Time: 30-45 minutes
```

---

## 🔄 Onboarding Flow

### **Phase 1: Initial Contact & Assessment** (Day 1-2)

```
┌─────────────────────────────────────┐
│   Manufacturer Inquiry              │
│   (Email/Form/Call)                 │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Initial Assessment                │
│   - Business validation             │
│   - Product fit check               │
│   - Market relevance                │
└──────────────┬──────────────────────┘
               │
         ┌─────┴─────┐
         │           │
    ✅ Qualified  ❌ Not Qualified
         │           │
         │           └──► Send Rejection Email
         │
         ▼
┌─────────────────────────────────────┐
│   Send Welcome Email                │
│   - Onboarding guide attached       │
│   - Information checklist           │
│   - Portal access link              │
└─────────────────────────────────────┘
```

**Actions:**
1. Business team receives inquiry
2. Validate basic business credentials
3. Check if equipment is relevant to platform
4. If qualified, send welcome email with onboarding portal link
5. If not qualified, send polite rejection with reasons

**Deliverables:**
- Welcome email sent
- Onboarding portal access credentials
- Dedicated account manager assigned

---

### **Phase 2: Information Collection** (Day 3-7)

```
┌─────────────────────────────────────┐
│   Manufacturer Portal Login         │
│   (Credentials from welcome email)  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Onboarding Dashboard              │
│   Shows 11 sections to complete     │
│   Progress tracker: 0/11            │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Section-by-Section Data Entry     │
│   ☐ 1. Basic Organization Info      │
│   ☐ 2. Address Information          │
│   ☐ 3. Contact Information          │
│   ☐ 4. Certifications & Compliance  │
│   ☐ 5. Banking Information          │
│   ☐ 6. Equipment Catalog            │
│   ☐ 7. Equipment Variants           │
│   ☐ 8. Parts Catalog                │
│   ☐ 9. Parts Compatibility          │
│   ☐ 10. Service & Warranty Info     │
│   ☐ 11. Error Codes                 │
└──────────────┬──────────────────────┘
               │
         ┌─────┴─────┐
         │           │
    Manufacturer    System
    Fills Forms     Auto-saves
         │           │
         │           ▼
         │   ┌───────────────┐
         │   │ Validation    │
         │   │ Rules Run     │
         │   └───────┬───────┘
         │           │
         │      ┌────┴────┐
         │      │         │
         │   ✅ Valid  ❌ Invalid
         │      │         │
         │      │         └──► Show Errors
         │      │              User Corrects
         │      │
         └──────┴──────────────────┐
                                   │
                                   ▼
                ┌──────────────────────────────┐
                │   Progress: 11/11 Complete   │
                │   [Submit for Review]         │
                └──────────────────────────────┘
```

**Actions:**
1. Manufacturer logs into onboarding portal
2. Completes 11 sections (can save and resume)
3. System validates each section in real-time
4. Manufacturer uploads required documents
5. Account manager monitors progress and offers help
6. Once all sections complete, manufacturer submits for review

**Validation Rules:**
- Required fields must be filled
- Email addresses must be verified (OTP)
- Phone numbers validated
- GST/Tax ID format checked
- Certifications must have valid expiry dates (not expired)
- At least 1 equipment model must be added
- Each equipment must have at least 1 part
- Parts must link to equipment models
- File uploads must be valid PDFs/images

**Deliverables:**
- All 11 sections completed
- Documents uploaded
- Data validated by system
- Submission for review

---

### **Phase 3: Verification & Approval** (Day 8-10)

```
┌─────────────────────────────────────┐
│   Submission Received               │
│   Status: "Pending Verification"    │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Automated Verification            │
│   - Duplicate check                 │
│   - GST validation (API)            │
│   - Email/phone verification        │
│   - Document quality check          │
└──────────────┬──────────────────────┘
               │
         ┌─────┴─────┐
         │           │
    ✅ Pass      ❌ Fail
         │           │
         │           └──► Notify Manufacturer
         │                Request Corrections
         │
         ▼
┌─────────────────────────────────────┐
│   Manual Review                     │
│   Business team validates:          │
│   - Business authenticity           │
│   - Certification validity          │
│   - Equipment relevance             │
│   - Parts catalog quality           │
└──────────────┬──────────────────────┘
               │
         ┌─────┴─────┐
         │           │
    ✅ Approved  ❌ Rejected
         │           │
         │           └──► Send Rejection
         │                With Reasons
         │
         ▼
┌─────────────────────────────────────┐
│   Create Organization Record        │
│   - org_type: "manufacturer"        │
│   - status: "active"                │
│   - All data imported               │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Create Related Records            │
│   - Facilities                      │
│   - Contacts                        │
│   - Equipment models                │
│   - Parts catalog                   │
│   - Certifications                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Send Approval Email               │
│   Status: "Active"                  │
│   Portal access: ENABLED            │
└─────────────────────────────────────┘
```

**Actions:**
1. System runs automated checks
2. Business team manually reviews
3. Legal team reviews contracts/agreements (if needed)
4. If approved, organization created in database
5. Approval email sent to manufacturer
6. Status changed to "Active"

**Verification Checklist:**
- [ ] No duplicate organization exists
- [ ] GST/Tax ID is valid and active
- [ ] Business registration verified
- [ ] ISO 13485 certificate valid and not expired
- [ ] Contact email verified (OTP sent)
- [ ] Phone number verified
- [ ] Documents are clear and readable
- [ ] Equipment models have proper documentation
- [ ] Parts catalog is complete
- [ ] Banking information validated (if provided)

**Deliverables:**
- Organization approved and activated
- Approval email sent
- Manufacturer can now access full platform

---

### **Phase 4: Platform Setup & Configuration** (Day 11-12)

```
┌─────────────────────────────────────┐
│   Organization Activated            │
│   Status: "Active"                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Technical Setup                   │
│   - User accounts created           │
│   - Role assignments                │
│   - Permissions configured          │
│   - Dashboard access enabled        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Data Integration                  │
│   - Equipment catalog indexed       │
│   - Parts catalog searchable        │
│   - AI knowledge base updated       │
│   - Search indexed                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Training & Onboarding             │
│   - Platform tour                   │
│   - Video tutorials                 │
│   - Documentation                   │
│   - Live training session (optional)│
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Go Live                           │
│   Manufacturer can now:             │
│   - View equipment installations    │
│   - Track service tickets           │
│   - Manage parts catalog            │
│   - Update documentation            │
│   - View analytics                  │
└─────────────────────────────────────┘
```

**Actions:**
1. Create user accounts for all contacts
2. Assign roles (Admin, Technical Support, Parts Manager)
3. Configure permissions
4. Index all data for search
5. Update AI knowledge base with manufacturer data
6. Schedule training session
7. Provide documentation
8. Go live!

**User Roles Created:**
- **Manufacturer Admin** - Full access
- **Technical Support** - View tickets, provide guidance
- **Parts Manager** - Manage parts catalog, pricing
- **Quality Manager** - View analytics, feedback
- **Finance** - View transactions (if applicable)

**Deliverables:**
- User accounts created
- Training completed
- Manufacturer is live on platform!

---

### **Phase 5: Post-Onboarding Support** (Ongoing)

```
┌─────────────────────────────────────┐
│   30-Day Check-in                   │
│   - Usage review                    │
│   - Issues identification           │
│   - Training gaps                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Quarterly Business Review         │
│   - Performance metrics             │
│   - Service quality                 │
│   - Parts sales (if applicable)     │
│   - AI accuracy for their equipment │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│   Continuous Improvement            │
│   - Catalog updates                 │
│   - New equipment models            │
│   - Error code additions            │
│   - Documentation improvements      │
└─────────────────────────────────────┘
```

**Activities:**
- Weekly: Monitor platform usage
- Monthly: Review analytics and feedback
- Quarterly: Business review meeting
- As needed: Catalog updates, new model additions

---

## 📋 Data Collection Forms

### Form 1: Basic Information Form

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   MANUFACTURER ONBOARDING - BASIC INFORMATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ORGANIZATION DETAILS
────────────────────────────────────────────

Legal Name: [________________________]  *Required
Brand Name: [________________________]  Optional
             (If different from legal name)

Organization Type: [▼ Manufacturer]  *Required

Manufacturer Subtype: [▼ Select One]  *Required
  ☐ OEM (Original Equipment Manufacturer)
  ☐ Component Manufacturer
  ☐ Contract Manufacturer
  ☐ Private Label Manufacturer

Registration Number: [________________________]  *Required
GST/Tax ID: [________________________]  *Required
PAN Number: [________________________]  *Required

Incorporation Date: [DD/MM/YYYY]  Recommended

Website: [________________________]  Recommended
Company Logo: [Choose File]  Recommended
              (PNG/JPG, Max 2MB)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[Save Draft]  [Continue to Next Section →]
```

### Form 2: Contact Information Form

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   MANUFACTURER ONBOARDING - CONTACT INFORMATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PRIMARY CONTACT (Authorized Representative)
────────────────────────────────────────────

Full Name: [________________________]  *Required
Designation: [________________________]  *Required

Email: [________________________]  *Required
       [Send Verification OTP]
       OTP: [______]  [Verify]

Mobile Phone: [+91] [__________]  *Required
              [Send SMS OTP]
              OTP: [______]  [Verify]

Alternate Phone: [+91] [__________]  Optional

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

TECHNICAL SUPPORT CONTACT
────────────────────────────────────────────

Full Name: [________________________]  *Required
Designation: [________________________]  *Required
Email: [________________________]  *Required
Phone: [+91] [__________]  *Required

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PARTS/WARRANTY CONTACT
────────────────────────────────────────────

Full Name: [________________________]  *Required
Designation: [________________________]  *Required
Email: [________________________]  *Required
Phone: [+91] [__________]  *Required

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

LEGAL/COMPLIANCE CONTACT
────────────────────────────────────────────

[☐ Same as Primary Contact]

Full Name: [________________________]  Optional
Designation: [________________________]
Email: [________________________]
Phone: [+91] [__________]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[← Back]  [Save Draft]  [Continue →]
```

### Form 3: Equipment Catalog Form

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   MANUFACTURER ONBOARDING - EQUIPMENT CATALOG
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

EQUIPMENT MODELS
────────────────────────────────────────────

[+ Add New Equipment Model]

┌──────────────────────────────────────────┐
│ Equipment #1                       [✕ Remove] │
├──────────────────────────────────────────┤
│                                          │
│ Category: [▼ Select Category]  *Required│
│   ☐ Ventilator                          │
│   ☐ MRI Scanner                         │
│   ☐ CT Scanner                          │
│   ☐ X-Ray Machine                       │
│   ☐ Ultrasound                          │
│   ☐ Patient Monitor                     │
│   ☐ Dialysis Machine                    │
│   ☐ Anesthesia Machine                  │
│   ☐ Infusion Pump                       │
│   ☐ Other                               │
│                                          │
│ Model Name: [_____________________]  *Req│
│                                          │
│ Model Number: [___________________]  *Req│
│                                          │
│ Model Code/SKU: [_________________]  *Req│
│                                          │
│ Description:                             │
│ [______________________________________ ] │
│ [______________________________________ ] │
│ [______________________________________ ]  *Required
│                                          │
│ Year Introduced: [YYYY]  Recommended     │
│                                          │
│ Currently Manufactured: ○ Yes  ○ No  *Req│
│                                          │
│ Warranty Period: [__] months  *Required  │
│                                          │
│ Expected Lifespan: [__] years  Recommended│
│                                          │
│ Product Image: [Choose File]  Recommended│
│                (JPG/PNG, Max 5MB)        │
│                                          │
│ Service Manual: [Choose File]  *Required │
│                 (PDF, Max 20MB)          │
│                                          │
│ Technical Specifications (JSON):         │
│ [______________________________________ ] │
│ [______________________________________ ] │
│                                          │
│ [+ Add Variant for this Model]          │
│                                          │
└──────────────────────────────────────────┘

[+ Add Another Equipment Model]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[← Back]  [Save Draft]  [Continue →]
```

### Form 4: Parts Catalog Form

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
   MANUFACTURER ONBOARDING - PARTS CATALOG
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PARTS LIST
────────────────────────────────────────────

💡 TIP: You can bulk upload parts using our Excel template
    [Download Template] [Upload Excel File]

Or add parts manually:

[+ Add New Part]

┌──────────────────────────────────────────┐
│ Part #1                          [✕ Remove] │
├──────────────────────────────────────────┤
│                                          │
│ Part Name: [_______________________]  *Req│
│                                          │
│ Part Number/SKU: [________________]  *Req│
│                                          │
│ Part Category: [▼ Select Category]  *Req │
│   ☐ Filter                              │
│   ☐ Sensor                              │
│   ☐ Battery                             │
│   ☐ Circuit Board                       │
│   ☐ Cable/Wire                          │
│   ☐ Seal/Gasket                         │
│   ☐ Valve                               │
│   ☐ Pump                                │
│   ☐ Display/Screen                      │
│   ☐ Other                               │
│                                          │
│ Description:                             │
│ [______________________________________ ] │
│ [______________________________________ ]  *Required
│                                          │
│ Compatible Equipment:  *Required         │
│   ☑ Ventilator V-100                    │
│   ☐ Ventilator V-200                    │
│   ☐ Patient Monitor PM-500              │
│                                          │
│ Compatible Variants:  Optional           │
│   ☑ V-100-ICU                           │
│   ☑ V-100-STD                           │
│   ☐ V-100-PED                           │
│                                          │
│ Unit of Measure: [▼ Piece]  *Required   │
│                                          │
│ Recommended Stock Level: [__]  Optional  │
│                                          │
│ Lead Time (Days): [__]  Optional         │
│                                          │
│ Expected Lifespan: [__] hours  Optional  │
│                                          │
│ List Price (₹): [________]  Optional     │
│                                          │
│ Part Image: [Choose File]  Recommended   │
│             (JPG/PNG, Max 2MB)           │
│                                          │
│ Installation Guide: [Choose File]  Opt   │
│                     (PDF, Max 10MB)      │
│                                          │
│ Related/Accessory Parts:                 │
│   [☐ Select related parts to recommend]  │
│                                          │
└──────────────────────────────────────────┘

[+ Add Another Part]

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[← Back]  [Save Draft]  [Continue →]
```

---

## ✅ Validation & Approval Process

### Automated Validations

1. **Business Registration**
   - GST/Tax ID format validation
   - GST API verification (if available)
   - PAN card format check

2. **Contact Verification**
   - Email OTP verification
   - Phone OTP verification
   - Email deliverability check

3. **Certification Validation**
   - Expiry date not in past
   - Certificate number format
   - Document quality check (readable PDF)

4. **Data Completeness**
   - All required fields filled
   - At least 1 equipment model
   - At least 1 part per equipment
   - Parts linked to equipment models

5. **Document Quality**
   - PDFs are readable (not corrupted)
   - Images are clear (min resolution)
   - File sizes within limits

### Manual Review Checklist

**Business Validation:**
- [ ] Company exists and is legitimate (Google search, LinkedIn)
- [ ] Business registration number verified
- [ ] ISO 13485 certificate authentic (check issuing body website)
- [ ] FDA/CE certificates authentic (if provided)
- [ ] No red flags in background check

**Technical Validation:**
- [ ] Equipment models are medically relevant
- [ ] Service manuals are comprehensive
- [ ] Parts catalog is complete
- [ ] Error codes provided are helpful
- [ ] Specifications are detailed enough

**Commercial Validation:**
- [ ] Banking information correct (name matches)
- [ ] Pricing is reasonable (if provided)
- [ ] Warranty terms are standard
- [ ] No conflicts with existing manufacturers

**Legal Validation:**
- [ ] Terms & conditions accepted
- [ ] Data sharing agreement signed
- [ ] HIPAA/GDPR compliance acknowledged (if applicable)
- [ ] No legal disputes or litigation

### Approval Criteria

**Must meet ALL of:**
1. ✅ Valid business registration
2. ✅ ISO 13485 certified
3. ✅ All required information provided
4. ✅ At least 3 equipment models OR 50+ parts
5. ✅ Service manuals uploaded
6. ✅ Technical support contact verified
7. ✅ No duplicate organization
8. ✅ Manual review passed

### Rejection Reasons

Common reasons for rejection:
- ❌ Business not legitimate
- ❌ ISO 13485 certificate expired or fake
- ❌ Incomplete information
- ❌ Equipment not medically relevant
- ❌ Poor quality documentation
- ❌ Failed background check
- ❌ Duplicate organization

---

## 🚀 Post-Onboarding Activities

### Immediate (Day 1 After Activation)

1. **Welcome Email** - "You're Live on GeneQR!"
2. **Platform Tour** - Interactive walkthrough
3. **Dashboard Access** - Full feature access enabled
4. **Documentation** - User guides and FAQs sent

### Week 1

1. **Training Session** - 1-hour live training
2. **Q&A Session** - Address any questions
3. **Data Verification** - Double-check all entered data
4. **Usage Monitoring** - Track initial activity

### Week 2-4

1. **Usage Review** - Analyze how manufacturer is using platform
2. **Issue Resolution** - Fix any problems encountered
3. **Feature Training** - Additional features explained
4. **Feedback Collection** - What's working, what's not

### Monthly

1. **Performance Report** - Monthly usage and metrics
2. **Catalog Updates** - Any new models/parts to add
3. **Support Ticket Review** - Common issues related to their equipment
4. **AI Accuracy Review** - How well AI diagnoses their equipment

### Quarterly

1. **Business Review** - Executive presentation
2. **Analytics Deep Dive** - Equipment performance data
3. **Strategic Planning** - Future roadmap discussion
4. **Contract Renewal** - If applicable

### As Needed

1. **New Equipment Onboarding** - Add new models as they're released
2. **Parts Catalog Updates** - Update pricing, add new parts
3. **Documentation Updates** - Update manuals, error codes
4. **Training Updates** - New technician training materials

---

## 🔧 Technical Integration

### Database Records Created

When a manufacturer is approved, the following records are created:

**1. Organization Record**
```sql
INSERT INTO organizations (
    id,
    name,
    org_type,
    status,
    metadata
) VALUES (
    'uuid',
    'Acme Medical Equipment Pvt Ltd',
    'manufacturer',
    'active',
    '{
        "subtype": "OEM",
        "registration_number": "U12345MH2010PTC123456",
        "gst_number": "27AAAAA0000A1Z5",
        "pan_number": "AAAAA0000A",
        "website": "https://acmemedical.com",
        "incorporation_date": "2010-01-15",
        "certifications": {
            "iso_13485": {
                "certified": true,
                "certificate_number": "ISO-12345",
                "valid_until": "2025-12-31"
            },
            "fda": {
                "registered": false
            },
            "ce": {
                "marked": false
            }
        }
    }'::jsonb
);
```

**2. Facility Records**
```sql
INSERT INTO organization_facilities (
    id,
    org_id,
    facility_name,
    facility_code,
    facility_type,
    address,
    status
) VALUES (
    'uuid',
    'org-uuid',
    'Acme Manufacturing Plant - Mumbai',
    'AMP-MUM-01',
    'manufacturing',
    '{
        "line1": "Plot 123, MIDC Industrial Area",
        "line2": "Andheri East",
        "city": "Mumbai",
        "state": "Maharashtra",
        "postal_code": "400093",
        "country": "India"
    }'::jsonb,
    'active'
);
```

**3. Contact Records**
```sql
INSERT INTO organization_contacts (
    id,
    org_id,
    contact_name,
    contact_role,
    email,
    phone,
    is_primary
) VALUES 
('uuid1', 'org-uuid', 'Rajesh Kumar', 'CEO', 'rajesh@acmemedical.com', '+919876543210', true),
('uuid2', 'org-uuid', 'Priya Sharma', 'Technical Head', 'priya@acmemedical.com', '+919876543211', false);
```

**4. Equipment Models**
```sql
INSERT INTO equipment_models (
    id,
    manufacturer_id,
    category,
    model_name,
    model_number,
    model_code,
    description,
    year_introduced,
    is_active,
    warranty_months,
    expected_lifespan_years,
    specifications
) VALUES (
    'uuid',
    'org-uuid',
    'Ventilator',
    'AcmeVent Pro',
    'AV-PRO-2023',
    'AVPRO2023',
    'Advanced ICU ventilator with AI-powered monitoring',
    2023,
    true,
    24,
    10,
    '{
        "technical": {
            "dimensions": "120cm x 80cm x 90cm",
            "weight": "150 kg",
            "power": "220V, 50Hz, 15A"
        },
        "features": ["AI Monitoring", "Touch Screen", "Remote Access"]
    }'::jsonb
);
```

**5. Equipment Variants**
```sql
INSERT INTO equipment_variants (
    id,
    model_id,
    variant_name,
    variant_code,
    description,
    specification_overrides
) VALUES
('uuid1', 'model-uuid', 'ICU Grade', 'AV-PRO-ICU', 'High-end variant for ICU', '{"features": ["Advanced Monitoring"]}'::jsonb),
('uuid2', 'model-uuid', 'Standard', 'AV-PRO-STD', 'Standard variant for general wards', '{}'::jsonb);
```

**6. Parts Catalog**
```sql
INSERT INTO parts_catalog (
    id,
    manufacturer_id,
    part_name,
    part_number,
    part_category,
    description,
    unit_of_measure,
    list_price
) VALUES (
    'uuid',
    'org-uuid',
    'HEPA Filter V-100',
    'HEPA-V100',
    'Filter',
    'High-efficiency particulate air filter for AcmeVent Pro',
    'Piece',
    8500.00
);
```

**7. Equipment-Parts Mapping**
```sql
INSERT INTO equipment_parts (
    id,
    model_id,
    part_id,
    is_required,
    expected_lifespan_hours
) VALUES (
    'uuid',
    'model-uuid',
    'part-uuid',
    true,
    8760  -- 1 year
);
```

**8. Error Codes**
```sql
INSERT INTO equipment_error_codes (
    id,
    model_id,
    error_code,
    error_message,
    description,
    severity,
    probable_causes,
    troubleshooting_steps
) VALUES (
    'uuid',
    'model-uuid',
    'E-42',
    'Filter Warning',
    'HEPA filter has reached end of life',
    'medium',
    '["Filter clogged", "Housing seal damaged", "Sensor malfunction"]'::jsonb,
    'Check filter for visible damage, inspect housing seal, test sensor, replace if needed'
);
```

### API Endpoints Used

**During Onboarding:**
- `POST /api/v1/manufacturers/onboarding/create` - Create onboarding record
- `PUT /api/v1/manufacturers/onboarding/:id/section/:sectionName` - Update section
- `POST /api/v1/manufacturers/onboarding/:id/documents/upload` - Upload documents
- `GET /api/v1/manufacturers/onboarding/:id/progress` - Check progress
- `POST /api/v1/manufacturers/onboarding/:id/submit` - Submit for review

**After Activation:**
- `GET /api/v1/organizations/:id` - Get organization details
- `PUT /api/v1/organizations/:id` - Update organization
- `GET /api/v1/organizations/:id/equipment` - List equipment models
- `POST /api/v1/equipment-models` - Add new equipment model
- `GET /api/v1/parts` - List parts catalog
- `PUT /api/v1/parts/:id` - Update part

---

## 📚 Support & Training

### Training Materials Provided

1. **Platform Overview Video** (15 minutes)
   - Platform tour
   - Key features
   - Navigation guide

2. **Equipment Management Video** (10 minutes)
   - How to add new models
   - Managing variants
   - Updating specifications

3. **Parts Catalog Video** (10 minutes)
   - Adding parts
   - Updating pricing
   - Managing inventory

4. **Analytics Dashboard Video** (8 minutes)
   - Reading metrics
   - Exporting reports
   - Insights interpretation

5. **User Guides** (PDF)
   - Complete platform documentation
   - Step-by-step screenshots
   - FAQs

### Support Channels

1. **Email Support** - support@geneqr.com
2. **Phone Support** - +91-XXXX-XXXXX (9 AM - 6 PM IST)
3. **Help Center** - help.geneqr.com
4. **Live Chat** - Available in portal
5. **Dedicated Account Manager** - For enterprise manufacturers

### Training Schedule

**Week 1:**
- Day 1: Platform overview and navigation
- Day 3: Equipment and parts management
- Day 5: Analytics and reporting

**Week 2:**
- Day 1: Q&A session
- Day 3: Advanced features
- Day 5: Best practices

**Ongoing:**
- Monthly webinars on new features
- Quarterly best practices sessions

---

## 📊 Success Metrics

### For Manufacturer

- **Time to Onboard:** <10 days (goal)
- **Data Completeness:** 100% of required fields
- **First Week Activity:** >50% of features used
- **30-Day Adoption:** Regular platform usage
- **Catalog Completeness:** >80% of equipment models added

### For GeneQR

- **Approval Rate:** >70% of submissions approved
- **Time to Review:** <48 hours
- **Manufacturer Satisfaction:** >4.5/5 rating
- **Platform Usage:** >80% weekly active
- **Catalog Quality:** >90% accuracy in AI diagnosis

---

## 🎯 Summary

### Onboarding Timeline

| Phase | Duration | Key Activities |
|-------|----------|----------------|
| **1. Initial Contact** | Day 1-2 | Inquiry, assessment, welcome email |
| **2. Information Collection** | Day 3-7 | Forms, documents, data entry |
| **3. Verification** | Day 8-10 | Validation, approval, setup |
| **4. Platform Setup** | Day 11-12 | Configuration, training |
| **5. Go Live** | Day 13+ | Active on platform |
| **TOTAL** | **10-15 days** | From inquiry to live |

### Key Takeaways

✅ **Clear Process:** Step-by-step flow with defined outcomes  
✅ **Comprehensive Data:** All information needed for platform operation  
✅ **Validation:** Multiple checks to ensure quality  
✅ **Support:** Training and ongoing assistance  
✅ **Integration:** Technical setup for seamless operation  

### Next Steps

1. **Review this document** with your team
2. **Design onboarding portal** (UI/UX)
3. **Implement forms and workflows** (development)
4. **Create training materials** (videos, docs)
5. **Test with pilot manufacturer** (validation)
6. **Launch onboarding process** (go live)

---

**Document End**

_For questions or clarifications, contact: GeneQR Platform Team_
