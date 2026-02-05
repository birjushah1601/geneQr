# Feature Flags - Notification System

**Date:** December 22, 2025  
**Status:** âœ… **Complete**

---

## ðŸŽ¯ Overview

Feature flag system to enable/disable notifications independently for each event type and channel (Email, SMS, WhatsApp).

---

## ðŸš€ Feature Flags

### **Master Switches**

Control entire notification channels:

| Flag | Default | Description |
|------|---------|-------------|
| `FEATURE_EMAIL_NOTIFICATIONS` | `false` | Master switch for all email notifications |
| `FEATURE_SMS_NOTIFICATIONS` | `false` | Master switch for all SMS notifications (future) |
| `FEATURE_WHATSAPP_NOTIFICATIONS` | `false` | Master switch for all WhatsApp notifications (future) |

### **Email Notification Events**

Individual toggles for each event type:

| Flag | Default | Description |
|------|---------|-------------|
| `FEATURE_EMAIL_TICKET_CREATED` | `false` | Email when ticket is created |
| `FEATURE_EMAIL_ENGINEER_ASSIGNED` | `false` | Email when engineer is assigned |
| `FEATURE_EMAIL_STATUS_CHANGED` | `false` | Email when ticket status changes |

### **SMS Notification Events (Future)**

| Flag | Default | Description |
|------|---------|-------------|
| `FEATURE_SMS_TICKET_CREATED` | `false` | SMS when ticket is created |
| `FEATURE_SMS_ENGINEER_ASSIGNED` | `false` | SMS when engineer is assigned |
| `FEATURE_SMS_STATUS_CHANGED` | `false` | SMS when ticket status changes |

### **WhatsApp Notification Events (Future)**

| Flag | Default | Description |
|------|---------|-------------|
| `FEATURE_WHATSAPP_TICKET_CREATED` | `false` | WhatsApp when ticket is created |
| `FEATURE_WHATSAPP_ENGINEER_ASSIGNED` | `false` | WhatsApp when engineer is assigned |
| `FEATURE_WHATSAPP_STATUS_CHANGED` | `false` | WhatsApp when ticket status changes |

---

## ðŸ“ Configuration

### **Quick Start - Enable All Email Notifications**

Add to `.env`:

```bash
# Enable Email Notifications
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true

# SendGrid Configuration (required for email)
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
SENDGRID_FROM_EMAIL=noreply@ServQR.com
SENDGRID_FROM_NAME=ServQR Service Platform

# Admin Email (for alerts)
ADMIN_EMAIL=admin@ServQR.com
```

### **Granular Control Examples**

**Example 1: Only Ticket Created Emails**
```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=false
FEATURE_EMAIL_STATUS_CHANGED=false
```

**Example 2: Only Engineer Assignment Emails**
```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=false
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=false
```

**Example 3: Status Changes Only**
```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=false
FEATURE_EMAIL_ENGINEER_ASSIGNED=false
FEATURE_EMAIL_STATUS_CHANGED=true
```

**Example 4: All Disabled (Development Mode)**
```bash
FEATURE_EMAIL_NOTIFICATIONS=false
# All individual flags ignored when master switch is off
```

---

## ðŸ”§ Boolean Value Formats

The system accepts multiple formats for boolean values:

**TRUE values:**
- `true`, `True`, `TRUE`
- `yes`, `Yes`, `YES`
- `y`, `Y`
- `1`
- `on`, `On`, `ON`
- `enabled`, `Enabled`, `ENABLED`

**FALSE values:**
- `false`, `False`, `FALSE`
- `no`, `No`, `NO`
- `n`, `N`
- `0`
- `off`, `Off`, `OFF`
- `disabled`, `Disabled`, `DISABLED`

---

## ðŸ’» Usage in Code

### **Initialization**

```go
import (
    "github.com/ServQR/internal/infrastructure/config"
    "github.com/ServQR/internal/infrastructure/email"
    "github.com/ServQR/internal/infrastructure/notification"
)

// Load feature flags
featureFlags := config.LoadFeatureFlags()

// Initialize email service
emailService := email.NewNotificationService(
    os.Getenv("SENDGRID_API_KEY"),
    os.Getenv("SENDGRID_FROM_EMAIL"),
    os.Getenv("SENDGRID_FROM_NAME"),
)

// Create notification manager
notificationManager := notification.NewManager(
    emailService,
    featureFlags,
    logger,
    os.Getenv("ADMIN_EMAIL"),
)
```

### **Sending Notifications**

The manager automatically checks feature flags:

```go
// Ticket Created
err := notificationManager.SendTicketCreatedNotifications(ctx, notification.TicketCreatedData{
    TicketNumber:  "TKT-001",
    CustomerName:  "John Doe",
    CustomerEmail: "john@example.com",
    // ... other fields
})

// Engineer Assigned
err := notificationManager.SendEngineerAssignedNotifications(ctx, notification.EngineerAssignedData{
    TicketNumber:  "TKT-001",
    EngineerName:  "Jane Engineer",
    EngineerEmail: "jane@engineer.com",
    // ... other fields
})

// Status Changed
err := notificationManager.SendStatusChangedNotifications(ctx, notification.StatusChangedData{
    TicketNumber: "TKT-001",
    OldStatus:    "new",
    NewStatus:    "assigned",
    // ... other fields
})
```

### **Checking Feature Status**

```go
// Get all feature flag states
status := notificationManager.GetFeatureStatus()

// Check if specific notification is enabled
if featureFlags.ShouldSendEmailNotification("ticket_created") {
    // Will send email
}
```

---

## ðŸ“Š Notification Flow with Feature Flags

```
Event Occurs (e.g., Ticket Created)
         â†“
[Notification Manager]
         â†“
Check Master Switch (FEATURE_EMAIL_NOTIFICATIONS)
         â†“
    Enabled? â”€â”€â”€â”€NOâ”€â”€â”€â†’ Skip, Log Debug
         â†“
        YES
         â†“
Check Event Flag (FEATURE_EMAIL_TICKET_CREATED)
         â†“
    Enabled? â”€â”€â”€â”€NOâ”€â”€â”€â†’ Skip, Log Debug
         â†“
        YES
         â†“
Send Notification (Email Service)
         â†“
    â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”´â”€â”€â”€â”€â”€â”€â”€â”€â”
    â†“                 â†“
[Success]        [Failure]
Log Info         Log Error
```

---

## ðŸ§ª Testing Feature Flags

### **Test 1: All Disabled (Default)**

```bash
# No flags set (or all false)
```

**Expected:**
- No emails sent
- Logs show: "notifications disabled by feature flag"

### **Test 2: Master Switch Only**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
```

**Expected:**
- No emails sent (individual flags still false)
- Logs show: "notifications disabled by feature flag"

### **Test 3: Individual Flag Without Master**

```bash
FEATURE_EMAIL_TICKET_CREATED=true
```

**Expected:**
- No emails sent (master switch is false)
- Logs show: "notifications disabled by feature flag"

### **Test 4: Full Enable**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
```

**Expected:**
- âœ… Emails sent on ticket creation
- Logs show: "Sending ticket created email notifications"

---

## ðŸŽ›ï¸ Runtime Configuration

### **Restart Required?**

**YES** - Feature flags are loaded on application startup from environment variables.

To change flags:
1. Update `.env` file
2. Restart backend server
3. Verify with logs

### **Dynamic Configuration (Future)**

For runtime changes without restart:
- Use configuration service (e.g., Consul, etcd)
- Implement config reload endpoint
- Use feature flag service (e.g., LaunchDarkly)

---

## ðŸ“‹ Environment File Examples

### **Development (.env.development)**

```bash
# Disable all notifications in dev
FEATURE_EMAIL_NOTIFICATIONS=false

# Use mock email service (logs to console)
# No SendGrid key needed
```

### **Staging (.env.staging)**

```bash
# Enable for testing
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true

# Use test SendGrid account
SENDGRID_API_KEY=SG.test_key_here
SENDGRID_FROM_EMAIL=staging@ServQR.com
SENDGRID_FROM_NAME=ServQR Staging

# Test admin email
ADMIN_EMAIL=admin-staging@ServQR.com
```

### **Production (.env.production)**

```bash
# Enable all email notifications
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true

# Production SendGrid
SENDGRID_API_KEY=SG.production_key_here
SENDGRID_FROM_EMAIL=noreply@ServQR.com
SENDGRID_FROM_NAME=ServQR Service Platform

# Production admin
ADMIN_EMAIL=admin@ServQR.com
```

---

## ðŸ” Monitoring & Logging

### **Log Levels**

**INFO:** Notification sent successfully
```
INFO: Sending ticket created email notifications ticket=TKT-001 customer=john@example.com
INFO: Ticket created email notifications sent successfully ticket=TKT-001
```

**DEBUG:** Notification skipped by feature flag
```
DEBUG: Ticket created email notifications disabled by feature flag ticket=TKT-001
```

**ERROR:** Notification failed
```
ERROR: Failed to send ticket created email ticket=TKT-001 error="sendgrid error: status 401"
```

### **Metrics to Track**

- Notifications attempted
- Notifications sent successfully
- Notifications failed
- Notifications skipped by feature flag
- Average send time
- Failure rate by type

---

## ðŸš€ Rollout Strategy

### **Phase 1: Silent Mode (Week 1)**

```bash
FEATURE_EMAIL_NOTIFICATIONS=false
```

- All notifications disabled
- System logs what would be sent
- Monitor logs for errors
- Verify data collection

### **Phase 2: Admin Only (Week 2)**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true  # Only admin emails
# Customer emails disabled
```

- Only admin receives notifications
- Verify email content and formatting
- Test all scenarios
- Collect feedback

### **Phase 3: One Event Type (Week 3)**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true  # Enable for customers too
FEATURE_EMAIL_ENGINEER_ASSIGNED=false
FEATURE_EMAIL_STATUS_CHANGED=false
```

- Enable ticket creation emails to customers
- Monitor email delivery rate
- Check for bounces/complaints
- Gather customer feedback

### **Phase 4: Gradual Rollout (Week 4)**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true  # Add engineer assignments
FEATURE_EMAIL_STATUS_CHANGED=false
```

### **Phase 5: Full Enable (Week 5+)**

```bash
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true  # All notifications
```

---

## ðŸ“Š Feature Flag Dashboard (Future)

### **Potential UI Features**

- Visual toggles for each flag
- Real-time status display
- Notification volume metrics
- Success/failure rates
- Toggle history/audit log

### **API Endpoint (Future)**

```
GET /api/v1/admin/feature-flags
POST /api/v1/admin/feature-flags/{flag}/toggle
```

---

## âœ… Checklist

### **Implementation**

- [x] Create feature flag configuration system
- [x] Create notification manager with flag support
- [x] Add master switches for each channel
- [x] Add individual event type flags
- [x] Document all flags and usage
- [ ] Integrate into ticket creation handler
- [ ] Integrate into engineer assignment service
- [ ] Integrate into status update handler
- [ ] Test all flag combinations
- [ ] Monitor in production

### **Documentation**

- [x] Feature flag reference
- [x] Configuration examples
- [x] Testing guide
- [x] Rollout strategy
- [ ] Operational runbook

---

## ðŸŽ¯ Summary

**Feature Flags Implemented:**
- âœ… 3 master switches (Email, SMS, WhatsApp)
- âœ… 9 individual event flags (3 per channel)
- âœ… Flexible boolean parsing
- âœ… Runtime checking
- âœ… Comprehensive logging

**Benefits:**
- ðŸŽ›ï¸ Granular control over notifications
- ðŸ”„ Easy rollback if issues occur
- ðŸ§ª Safe testing in production
- ðŸ“Š Better monitoring and metrics
- ðŸš€ Gradual rollout capability

**Next Steps:**
1. Add to `.env` file
2. Integrate notification calls
3. Test with flags disabled
4. Test with flags enabled
5. Gradual rollout to production

---

**Last Updated:** December 22, 2025  
**Status:** âœ… Complete - Ready for Integration  
**Files Created:** 
- `internal/infrastructure/config/feature_flags.go`
- `internal/infrastructure/notification/manager.go`
