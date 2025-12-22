# Email Notification System - Implementation Guide

**Date:** December 22, 2025  
**Status:** ✅ **Infrastructure Ready** | ⚠️ **Integration Pending**

---

## 🎯 Overview

Email notification system for alerting ticket creators, admins, and engineers about ticket lifecycle events.

---

## ✅ What's Available

### **1. Email Infrastructure (✅ Complete)**

**File:** `internal/infrastructure/email/sendgrid.go`

**Features:**
- SendGrid API integration
- OTP email sending (currently used for authentication)
- HTML and plain text email support
- Error handling and status checking

**Usage:**
```go
sender := email.NewSendGridSender(apiKey, fromEmail, fromName)
err := sender.SendOTP(ctx, toEmail, otpCode)
```

---

### **2. Notification Service (✅ Just Created)**

**File:** `internal/infrastructure/email/notification.go`

**Features:**
- **Ticket Created Notifications**
  - Email to customer (confirmation)
  - Email to admin (new ticket alert)

- **Engineer Assigned Notifications**
  - Email to customer (engineer details)
  - Email to engineer (ticket assignment)

- **Status Changed Notifications**
  - Email to customer (status update)
  - Email to admin (important status changes)

**Email Templates:**
- Professional HTML emails with styling
- Plain text fallback
- Mobile-responsive design
- Company branding

---

## 📧 Notification Types

### **1. Ticket Created**

**Triggers:** When a customer creates a service ticket

**Recipients:**
- ✅ Customer (confirmation)
- ✅ Admin (new ticket alert)

**Customer Email Includes:**
- Ticket number
- Equipment details
- Priority level
- Description
- Contact information
- Next steps

**Admin Email Includes:**
- Ticket number
- Customer details (name, phone, email)
- Equipment details
- Priority level
- Description
- Action required alert

---

### **2. Engineer Assigned**

**Triggers:** When an engineer is assigned to a ticket

**Recipients:**
- ✅ Customer (engineer details)
- ✅ Engineer (new assignment)

**Customer Email Includes:**
- Ticket number
- Engineer name
- Engineer phone
- Engineer email
- "Engineer will contact you" message

**Engineer Email Includes:**
- Ticket number
- Customer details
- Equipment details
- Priority level
- Issue description
- Action required (contact customer)

---

### **3. Status Changed**

**Triggers:** When ticket status is updated

**Recipients:**
- ✅ Customer (always)
- ✅ Admin (for resolved/closed/cancelled)

**Customer Email Includes:**
- Ticket number
- Old status → New status
- Equipment details
- Updated by (who made the change)

**Admin Email Includes:**
- Ticket number
- Customer name
- Status change
- Updated by

---

## 🔧 Integration Points (⚠️ To Be Implemented)

### **Where to Add Notifications:**

1. **Ticket Creation**
   - **File:** `internal/service-domain/service-ticket/api/handler.go`
   - **Function:** `CreateServiceTicket`
   - **Action:** Call `notificationService.SendTicketCreatedNotification()`

2. **Engineer Assignment**
   - **File:** `internal/service-domain/service-ticket/app/assignment_service.go`
   - **Function:** `AssignEngineer`
   - **Action:** Call `notificationService.SendTicketAssignedNotification()`

3. **Status Change**
   - **File:** `internal/service-domain/service-ticket/api/handler.go`
   - **Function:** `UpdateTicketStatus`
   - **Action:** Call `notificationService.SendTicketStatusChangedNotification()`

---

## 📝 Integration Example

### **Step 1: Initialize Notification Service**

Add to `internal/service-domain/service-ticket/module.go`:

```go
import (
    "github.com/aby-med/internal/infrastructure/email"
    "os"
)

// In NewTicketModule function
func NewTicketModule(db *sql.DB, logger *slog.Logger) *TicketModule {
    // ... existing code ...

    // Initialize notification service
    notificationService := email.NewNotificationService(
        os.Getenv("SENDGRID_API_KEY"),
        os.Getenv("SENDGRID_FROM_EMAIL"),
        os.Getenv("SENDGRID_FROM_NAME"),
    )

    // Pass to handler
    handler := NewTicketHandler(service, notificationService, logger)

    return &TicketModule{
        Handler: handler,
    }
}
```

### **Step 2: Add to Ticket Creation**

Update `CreateServiceTicket` handler:

```go
func (h *TicketHandler) CreateServiceTicket(w http.ResponseWriter, r *http.Request) {
    // ... existing ticket creation code ...

    // After ticket is created successfully:
    go func() {
        ctx := context.Background()
        err := h.notificationService.SendTicketCreatedNotification(ctx, email.TicketCreatedData{
            TicketNumber:  ticket.TicketNumber,
            CustomerName:  ticket.CustomerName,
            CustomerEmail: ticket.CustomerEmail,
            CustomerPhone: ticket.CustomerPhone,
            EquipmentName: ticket.EquipmentName,
            Description:   ticket.Description,
            Priority:      ticket.Priority,
            AdminEmail:    "admin@aby-med.com", // Or fetch from config
        })
        if err != nil {
            h.logger.Error("Failed to send ticket created notification", "error", err)
        }
    }()

    // ... rest of handler ...
}
```

### **Step 3: Add to Engineer Assignment**

Update `AssignEngineer` in assignment service:

```go
func (s *AssignmentService) AssignEngineer(ctx context.Context, ticketID, engineerID string) error {
    // ... existing assignment code ...

    // After assignment is successful:
    go func() {
        // Fetch ticket and engineer details
        ticket, _ := s.ticketRepo.GetByID(ctx, ticketID)
        engineer, _ := s.engineerRepo.GetByID(ctx, engineerID)

        err := s.notificationService.SendTicketAssignedNotification(ctx, email.TicketAssignedData{
            TicketNumber:  ticket.TicketNumber,
            CustomerName:  ticket.CustomerName,
            CustomerEmail: ticket.CustomerEmail,
            EngineerName:  engineer.Name,
            EngineerEmail: engineer.Email,
            EngineerPhone: engineer.Phone,
            EquipmentName: ticket.EquipmentName,
            Description:   ticket.Description,
            Priority:      ticket.Priority,
        })
        if err != nil {
            s.logger.Error("Failed to send assignment notification", "error", err)
        }
    }()

    return nil
}
```

### **Step 4: Add to Status Change**

Update `UpdateTicketStatus` handler:

```go
func (h *TicketHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
    // ... existing status update code ...

    // After status is updated:
    go func() {
        ctx := context.Background()
        err := h.notificationService.SendTicketStatusChangedNotification(ctx, email.TicketStatusChangedData{
            TicketNumber:  ticket.TicketNumber,
            CustomerName:  ticket.CustomerName,
            CustomerEmail: ticket.CustomerEmail,
            OldStatus:     oldStatus,
            NewStatus:     newStatus,
            EquipmentName: ticket.EquipmentName,
            UpdatedBy:     updatedByUserName,
            AdminEmail:    "admin@aby-med.com",
        })
        if err != nil {
            h.logger.Error("Failed to send status change notification", "error", err)
        }
    }()

    // ... rest of handler ...
}
```

---

## 🔐 Environment Configuration

### **Required Environment Variables:**

```bash
# SendGrid API Configuration
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
SENDGRID_FROM_EMAIL=noreply@aby-med.com
SENDGRID_FROM_NAME=ABY-MED Service Platform
```

### **How to Get SendGrid API Key:**

1. Sign up at https://sendgrid.com
2. Go to Settings → API Keys
3. Create New API Key
4. Select "Full Access" or "Mail Send" permissions
5. Copy the API key (shown only once)
6. Add to `.env` file

### **For Development (Optional):**

If you don't have SendGrid configured, the system will use a mock email service that logs to console:

```
📧 MOCK EMAIL to=customer@example.com subject=Service Ticket Created
```

---

## 🎨 Email Templates

### **Visual Design:**

**Colors:**
- Header backgrounds: Blue (#2563eb), Green (#10b981), Red (#dc2626), Orange (#f59e0b)
- Content background: Light gray (#f9fafb)
- Text: Dark gray (#333)
- Borders: Blue/matching theme colors

**Layout:**
- Responsive (max-width: 600px)
- Clear header with icon
- Clean content area
- Info boxes with left border
- Call-to-action sections
- Footer with company info

**Typography:**
- Font: Arial, sans-serif
- Line height: 1.6
- Readable font sizes
- Bold for emphasis

---

## 📊 Notification Flow Diagram

```
Customer Creates Ticket
         ↓
    [Ticket Saved]
         ↓
    ┌─────────────┬─────────────┐
    ↓             ↓             ↓
[Customer]    [Admin]      [Audit Log]
Confirmation  Alert        Record
Email         Email        Event

---

Admin Assigns Engineer
         ↓
  [Assignment Saved]
         ↓
    ┌─────────────┬─────────────┐
    ↓             ↓             ↓
[Customer]    [Engineer]   [Audit Log]
Engineer      Assignment   Record
Details       Email        Event

---

Status Changes
         ↓
  [Status Updated]
         ↓
    ┌─────────────┬──────────────┐
    ↓             ↓              ↓
[Customer]    [Admin]*      [Audit Log]
Status        Status         Record
Update        Update*        Event
Email         Email

* Admin email only for resolved/closed/cancelled
```

---

## 🧪 Testing

### **Manual Testing:**

1. **Test Ticket Creation Email:**
   ```bash
   # Create a ticket via API or UI
   # Check SendGrid dashboard for sent emails
   # Check customer inbox
   # Check admin inbox
   ```

2. **Test Engineer Assignment Email:**
   ```bash
   # Assign an engineer to a ticket
   # Check customer inbox (engineer details)
   # Check engineer inbox (assignment)
   ```

3. **Test Status Change Email:**
   ```bash
   # Update ticket status
   # Check customer inbox (status update)
   # Check admin inbox (if status is resolved/closed/cancelled)
   ```

### **Testing Without SendGrid:**

If `SENDGRID_API_KEY` is not set, emails will be logged to console:

```
📧 MOCK EMAIL to=customer@example.com
Subject: Service Ticket Created - TKT-001
Ticket Number: TKT-001
Equipment: MRI Scanner
Priority: high
```

---

## 🚀 Future Enhancements (SMS/WhatsApp)

### **SMS Notifications (Planned)**

**File:** `internal/infrastructure/sms/twilio.go` (to be created)

**Features:**
- SMS for critical updates
- SMS for engineer assignment
- SMS for status changes

**Integration:**
```go
// SMS notification service
smsService := sms.NewTwilioService(accountSID, authToken, fromNumber)
err := smsService.SendTicketCreatedSMS(phone, ticketNumber)
```

### **WhatsApp Notifications (Planned)**

**File:** `internal/service-domain/whatsapp/service.go` (already exists)

**Features:**
- WhatsApp messages for updates
- Rich formatting with images
- Interactive buttons

**Integration:**
```go
// WhatsApp notification
whatsappService := whatsapp.NewService(twilioClient)
err := whatsappService.SendTicketUpdate(phone, ticketData)
```

---

## 📋 Implementation Checklist

### **Phase 1: Email Notifications (Current)**

- [x] Create email notification service
- [x] Design email templates (HTML + plain text)
- [x] Add ticket creation notifications
- [x] Add engineer assignment notifications
- [x] Add status change notifications
- [ ] Integrate into ticket creation handler
- [ ] Integrate into assignment service
- [ ] Integrate into status update handler
- [ ] Test with real SendGrid API
- [ ] Test all notification types
- [ ] Document integration points

### **Phase 2: SMS Notifications (Future)**

- [ ] Create SMS notification service
- [ ] Integrate Twilio SMS API
- [ ] Add SMS templates
- [ ] Integrate into ticket lifecycle
- [ ] Test SMS delivery

### **Phase 3: WhatsApp Notifications (Future)**

- [ ] Use existing WhatsApp service
- [ ] Create WhatsApp message templates
- [ ] Integrate into ticket lifecycle
- [ ] Test WhatsApp delivery

---

## 🎯 Summary

### **What's Ready:**
- ✅ Email infrastructure (SendGrid)
- ✅ Notification service with templates
- ✅ Professional HTML emails
- ✅ Three notification types (created, assigned, status changed)
- ✅ Mock email service for development

### **What's Needed:**
- ⚠️ Integration into ticket handlers
- ⚠️ Integration into assignment service
- ⚠️ SendGrid API key configuration
- ⚠️ Testing with real emails
- ⚠️ Admin email configuration

### **Estimated Integration Time:**
- **30-45 minutes** to integrate all three notification points
- **15 minutes** to test with SendGrid
- **Total: ~1 hour**

---

## 📞 Support

**For SendGrid Issues:**
- Documentation: https://docs.sendgrid.com
- API Reference: https://www.twilio.com/docs/sendgrid/api-reference

**For Integration Help:**
- See integration examples above
- Check existing OTP email implementation in `sendgrid.go`
- Follow the pattern: initialize → call in handler → log errors

---

**Last Updated:** December 22, 2025  
**Status:** Infrastructure Ready - Integration Pending  
**Next Step:** Integrate notification calls into ticket handlers
