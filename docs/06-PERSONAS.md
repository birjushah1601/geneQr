# ABY-MED User Personas & Use Cases

Documentation organized by user perspective and real-world scenarios.

---

## 👥 User Personas

### 1. Hospital Administrator (Customer)
**Profile:** Manages hospital equipment, tracks service requests  
**Goals:** Quick ticket creation, track repair status, manage equipment inventory

**Key Features:**
- ✅ Equipment registry with QR codes
- ✅ Create tickets via web or WhatsApp
- ✅ Track ticket status in real-time
- ✅ View equipment maintenance history
- ✅ Receive email notifications

**User Journey:**
```
1. Equipment malfunction detected
2. Scan QR code on equipment
3. Create ticket (web/WhatsApp)
4. Receive confirmation + ticket number
5. Track engineer assignment
6. Receive status updates (email)
7. Review and close ticket
```

---

### 2. Field Engineer (Service Provider)
**Profile:** Repairs medical equipment on-site  
**Goals:** View assigned tickets, update status, request parts, close tickets

**Key Features:**
- ✅ View assigned tickets dashboard
- ✅ Update ticket status
- ✅ Add comments/photos
- ✅ Request spare parts
- ✅ Mark tickets resolved

**User Journey:**
```
1. Login to dashboard
2. View assigned tickets (filtered list)
3. Click ticket → See equipment details
4. Update status: "In Progress"
5. Add diagnostic comments
6. Request parts if needed
7. Complete repair → Mark "Resolved"
8. Upload completion photos
```

---

### 3. Manufacturer Admin (Platform Owner)
**Profile:** Manages service operations, assigns engineers, monitors performance  
**Goals:** Efficient ticket assignment, team management, analytics

**Key Features:**
- ✅ Dashboard with all tickets
- ✅ Manual engineer assignment
- ✅ AI-powered assignment suggestions
- ✅ Update ticket priority (admin-only)
- ✅ Daily reports (email)
- ✅ Organization management
- ✅ Equipment catalog management
- ✅ Onboarding system (bulk import)

**User Journey:**
```
1. Login → Dashboard overview
2. See new tickets requiring assignment
3. Click ticket → View details
4. Use AI suggestions for engineer
5. Assign engineer → Notification sent
6. Update priority if critical
7. Monitor progress
8. Review daily reports (morning/evening)
```

---

### 4. Service Manager (Manufacturer)
**Profile:** Oversees service operations, analytics, reporting  
**Goals:** Monitor SLAs, team performance, customer satisfaction

**Key Features:**
- ✅ Analytics dashboard (coming soon)
- ✅ SLA tracking and alerts
- ✅ Daily email reports
- ✅ Ticket history and trends
- ✅ Engineer performance metrics

**User Journey:**
```
1. Review morning report email
2. Check SLA breaches
3. Identify bottlenecks
4. Reassign overloaded engineers
5. Review evening report
6. Plan next day assignments
```

---

### 5. Equipment Manufacturer (Platform Seller - Future)
**Profile:** Sells spare parts via marketplace  
**Goals:** Manage products, fulfill orders, track sales

**Key Features:**
- 🚧 Product management dashboard
- 🚧 Inventory tracking
- 🚧 Order fulfillment
- 🚧 Sales analytics

**User Journey (Planned):**
```
1. Login → Seller dashboard
2. Add/update products
3. Receive order notification
4. Process order
5. Update tracking info
6. View sales reports
```

---

## 🎯 Use Cases by Scenario

### Scenario 1: Emergency Equipment Failure (Hospital)
**Persona:** Hospital Administrator  
**Urgency:** Critical

**Flow:**
1. MRI machine stops working during patient scan
2. Staff scans QR code on equipment
3. Creates ticket via web form: "MRI not starting, error code E-503"
4. Priority auto-set to "medium", admin escalates to "critical"
5. System notifies manufacturer immediately
6. AI suggests senior engineer with MRI experience
7. Engineer assigned within 5 minutes
8. Engineer arrives, diagnoses, requests part
9. Part ordered from marketplace (future)
10. Equipment repaired, ticket closed

**Features Used:**
- QR code system
- Web ticket creation
- Admin priority update
- AI engineer assignment
- Email notifications
- Parts management

---

### Scenario 2: Routine Maintenance via WhatsApp (Clinic)
**Persona:** Clinic Owner (Non-technical)  
**Urgency:** Low

**Flow:**
1. Clinic owner notices ventilator making noise
2. Takes photo of QR code on ventilator
3. Sends WhatsApp message: "QR-20251223-005 making strange noise"
4. Records voice note describing the sound
5. System auto-creates ticket
6. Whisper AI transcribes audio
7. Ticket created with text + audio + transcript
8. Confirmation sent via WhatsApp: "Ticket #TKT-20251223-042 created"
9. Engineer assigned next day
10. Engineer visits, performs maintenance

**Features Used:**
- WhatsApp integration
- Audio message handling
- Whisper STT transcription
- Auto-ticket creation
- QR code recognition

---

### Scenario 3: Bulk Manufacturer Onboarding (Manufacturer)
**Persona:** New Manufacturer Admin  
**Urgency:** One-time setup

**Flow:**
1. Manufacturer signs up
2. Navigates to onboarding wizard
3. Step 1: Enters company details
4. Step 2: Bulk imports 50 hospital organizations (CSV)
5. Step 3: Selects "Radiology" industry template
6. Downloads pre-configured template with 8 equipment types
7. Uploads CSV with 200 equipment items
8. System imports in 2 minutes
9. Generates QR codes for all equipment
10. Completion page shows: 50 orgs, 200 equipment imported
11. Ready to start operations

**Features Used:**
- Onboarding wizard
- Bulk CSV import
- Industry templates
- QR batch generation
- 97% time reduction (5h → 5-10 min)

---

### Scenario 4: AI-Assisted Diagnosis (Engineer)
**Persona:** Junior Engineer  
**Urgency:** Medium

**Flow:**
1. Engineer views ticket: "CT Scanner error E-304"
2. Clicks "Get AI Diagnosis" button
3. Selects model: GPT-4
4. AI analyzes:
   - Equipment model: Siemens Somatom
   - Error code: E-304
   - Maintenance history
5. AI suggests: "Power supply module failure, replace PSU-204"
6. Engineer accepts diagnosis
7. Requests part PSU-204
8. Part shipped, arrives next day
9. Engineer replaces module
10. Marks ticket resolved
11. Feedback: "Diagnosis accurate" → Improves AI model

**Features Used:**
- AI diagnosis (GPT-4)
- Multi-model support
- Parts suggestion
- Parts request workflow
- Feedback loop

---

## 🔄 Cross-Persona Workflows

### Workflow 1: Ticket Lifecycle (All Personas)
```
Hospital Staff (Create)
    ↓
Manufacturer Admin (Assign)
    ↓
Field Engineer (Work)
    ↓
Parts Supplier (Provide parts - future)
    ↓
Field Engineer (Complete)
    ↓
Hospital Staff (Close)
```

### Workflow 2: Equipment Lifecycle
```
Manufacturer (Produce equipment)
    ↓
Distributor (Ship to hospital)
    ↓
Hospital (Install, register in system)
    ↓
Service Provider (Maintain)
    ↓
Engineer (Repair when needed)
    ↓
Parts Supplier (Provide spare parts)
```

---

## 📚 Persona-Specific Documentation

### For Hospital Admins
- [01-GETTING-STARTED.md](./01-GETTING-STARTED.md) - Section: "Access Application"
- [03-FEATURES.md](./03-FEATURES.md) - Section: "Service Ticket Features"
- [QUICK-REFERENCE.md](./QUICK-REFERENCE.md)

### For Engineers
- [03-FEATURES.md](./03-FEATURES.md) - Section: "Engineer Management"
- [AI_INTEGRATION_STATUS.md](./AI_INTEGRATION_STATUS.md)

### For Manufacturer Admins
- [MANUFACTURER-ONBOARDING-UX-DESIGN.md](./MANUFACTURER-ONBOARDING-UX-DESIGN.md)
- [ONBOARDING-SYSTEM-README.md](./ONBOARDING-SYSTEM-README.md)
- [MULTI-TENANT-IMPLEMENTATION-PLAN.md](./MULTI-TENANT-IMPLEMENTATION-PLAN.md)

### For Service Managers
- [DAILY-REPORTS-SYSTEM.md](./DAILY-REPORTS-SYSTEM.md)
- [EMAIL-NOTIFICATIONS-SYSTEM.md](./EMAIL-NOTIFICATIONS-SYSTEM.md)

---

## 🎭 Persona Comparison Matrix

| Feature | Hospital Admin | Engineer | Mfr Admin | Service Mgr |
|---------|----------------|----------|-----------|-------------|
| Create Ticket | ✅ Primary | ❌ | ✅ | ❌ |
| Assign Engineer | ❌ | ❌ | ✅ Primary | ✅ |
| Update Status | ❌ | ✅ Primary | ✅ | ✅ |
| Update Priority | ❌ | ❌ | ✅ Admin only | ✅ Admin |
| View Analytics | ⚠️ Limited | ⚠️ Own | ✅ All | ✅ Primary |
| Manage Equipment | ✅ Own | ❌ | ✅ All | ⚠️ View |
| Request Parts | ⚠️ Indirect | ✅ Primary | ✅ | ❌ |
| Onboard Orgs | ❌ | ❌ | ✅ Admin only | ❌ |

---

**Last Updated:** December 23, 2025  
**Status:** Production Personas
