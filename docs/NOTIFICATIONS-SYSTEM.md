# Notifications System - Complete Guide

Comprehensive guide for all notification features: Email notifications and Daily reports.

---

## ðŸ“§ Email Notifications

### Overview
Real-time email notifications for ticket lifecycle events using SendGrid.

### Notification Types

**1. Ticket Created**
- **Trigger:** New service ticket created
- **Recipients:** Customer (confirmation) + Admin (alert)
- **Includes:** Ticket #, equipment, priority, description, next steps

**2. Engineer Assigned**
- **Trigger:** Engineer assigned to ticket
- **Recipients:** Customer (engineer info) + Engineer (assignment)
- **Includes:** Ticket #, engineer contact, expected timeline

**3. Status Changed**
- **Trigger:** Ticket status update
- **Recipients:** Customer (update) + Admin (critical changes)
- **Includes:** Old/new status, updated by, next actions

### Configuration

```bash
# .env
SENDGRID_API_KEY=SG.xxx
SENDGRID_FROM_EMAIL=noreply@ServQR.com
SENDGRID_FROM_NAME="ServQR Platform"

# Feature flags
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true
```

### Implementation

**Files:**
- `internal/infrastructure/email/sendgrid.go` - SendGrid client
- `internal/infrastructure/email/notification.go` - Email templates & sending

**Usage Example:**
```go
// In ticket creation handler
notificationService := email.NewNotificationService(sendgridClient)
err := notificationService.SendTicketCreatedNotifications(ctx, ticket, equipment, customer)

// In engineer assignment handler
err := notificationService.SendEngineerAssignedNotifications(ctx, ticket, engineer, customer)

// In status update handler
err := notificationService.SendStatusChangedNotifications(ctx, ticket, oldStatus, customer)
```

### Email Templates
- Professional HTML with styling
- Plain text fallback
- Mobile-responsive
- Company branding included

---

## ðŸ“Š Daily Reports

### Overview
Automated daily email reports sent twice daily with platform metrics and summaries.

### Report Schedule

**Morning Report:** 8:00 AM
- Fresh start data for the day
- Pending items requiring attention

**Evening Report:** 6:00 PM
- End of day summary
- Today's accomplishments

### Report Contents

Reports include 8 key data categories:

1. **Tickets Summary**
   - Total active tickets
   - By status (new, assigned, in_progress, resolved)
   - By priority (critical, high, medium, low)
   - Overdue tickets

2. **Engineers Performance**
   - Total engineers
   - Available vs busy
   - By skill level
   - Active assignments

3. **Equipment Status**
   - Total equipment registered
   - Equipment with active tickets
   - Equipment by manufacturer
   - Recently registered

4. **Organizations**
   - Total organizations
   - By type (manufacturer, hospital, clinic)
   - Active vs inactive
   - Recently added

5. **Today's Activity**
   - New tickets created today
   - Tickets resolved today
   - Engineers assigned today
   - Equipment registered today

6. **SLA Tracking**
   - Tickets approaching deadline
   - SLA breaches
   - Average resolution time

7. **Parts Management**
   - Parts requested today
   - Parts catalog size
   - Low stock alerts (future)

8. **System Health**
   - API uptime
   - Database connections
   - Error rates

### Configuration

```bash
# .env
FEATURE_DAILY_REPORTS=true
DAILY_REPORT_RECIPIENT_EMAIL=admin@ServQR.com
REPORT_SCHEDULE_MORNING="0 8 * * *"   # 8 AM daily
REPORT_SCHEDULE_EVENING="0 18 * * *"  # 6 PM daily
SENDGRID_API_KEY=SG.xxx
```

### Implementation

**Files:**
- `internal/infrastructure/reports/daily_report.go` - Report generation
- `internal/infrastructure/reports/scheduler.go` - Cron scheduling

**How it Works:**
1. Cron job triggers at scheduled times
2. Report service queries database for metrics
3. HTML email generated with data
4. Email sent via SendGrid
5. Logs sent/failed status

**Manual Trigger:**
```go
reportService := reports.NewDailyReportService(db, sendgridClient)
err := reportService.SendDailyReport(ctx, "admin@example.com")
```

### Report Format
- Professional HTML email
- Tables for metrics
- Color-coded status indicators
- Charts (future enhancement)
- Mobile-responsive layout

---

## ðŸŽšï¸ Feature Flags

Control notifications at granular level:

```bash
# Master switches
FEATURE_EMAIL_NOTIFICATIONS=true
FEATURE_DAILY_REPORTS=true

# Email notification types
FEATURE_EMAIL_TICKET_CREATED=true
FEATURE_EMAIL_ENGINEER_ASSIGNED=true
FEATURE_EMAIL_STATUS_CHANGED=true

# Report schedules
REPORT_SCHEDULE_MORNING="0 8 * * *"
REPORT_SCHEDULE_EVENING="0 18 * * *"
```

### Disabling Notifications

**Disable all emails:**
```bash
FEATURE_EMAIL_NOTIFICATIONS=false
```

**Disable specific notification:**
```bash
FEATURE_EMAIL_TICKET_CREATED=false
```

**Disable daily reports:**
```bash
FEATURE_DAILY_REPORTS=false
```

---

## ðŸ§ª Testing

### Test Email Notifications

```bash
# Create test ticket
curl -X POST http://localhost:8081/api/v1/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": "uuid",
    "issue_description": "Test ticket",
    "priority": "medium"
  }'

# Check SendGrid activity
# Visit: https://app.sendgrid.com/email_activity
```

### Test Daily Reports

```go
// In code or test file
reportService := reports.NewDailyReportService(db, sendgridClient)
err := reportService.SendDailyReport(context.Background(), "your-email@example.com")
```

---

## ðŸ“ˆ Monitoring

### Email Delivery Tracking

**SendGrid Dashboard:**
- Delivery rate
- Open rate
- Bounce rate
- Click-through rate

**Application Logs:**
```bash
# Success logs
INFO: Email notification sent successfully, type=ticket_created, recipient=customer@example.com

# Failure logs
ERROR: Failed to send email notification, type=engineer_assigned, error=invalid API key
```

### Report Delivery Status

**Daily Report Logs:**
```bash
INFO: Daily report sent successfully, time=08:00, recipient=admin@ServQR.com, metrics_count=8
ERROR: Failed to send daily report, error=database timeout
```

---

## ðŸ”§ Troubleshooting

### Email Not Sending

1. **Check SendGrid API Key:**
   ```bash
   echo $SENDGRID_API_KEY
   ```

2. **Verify feature flag:**
   ```bash
   echo $FEATURE_EMAIL_NOTIFICATIONS
   ```

3. **Check SendGrid sender verification:**
   - Visit SendGrid â†’ Settings â†’ Sender Authentication
   - Verify domain or single sender email

4. **Review logs:**
   ```bash
   grep "email notification" backend.log
   ```

### Daily Reports Not Arriving

1. **Check cron schedule:**
   ```bash
   echo $REPORT_SCHEDULE_MORNING
   echo $REPORT_SCHEDULE_EVENING
   ```

2. **Verify feature flag:**
   ```bash
   echo $FEATURE_DAILY_REPORTS
   ```

3. **Check recipient email:**
   ```bash
   echo $DAILY_REPORT_RECIPIENT_EMAIL
   ```

4. **Test manual trigger:**
   ```go
   reportService.SendDailyReport(ctx, "your-email@example.com")
   ```

---

## ðŸš€ Future Enhancements

### Email Notifications
- [ ] SMS notifications (Twilio)
- [ ] Push notifications (mobile app)
- [ ] Slack/Teams integration
- [ ] Webhook callbacks
- [ ] Email templates customization (per org)
- [ ] Notification preferences (per user)
- [ ] Digest mode (batch notifications)

### Daily Reports
- [ ] Weekly summary reports
- [ ] Monthly executive reports
- [ ] Custom report builder
- [ ] Export to PDF
- [ ] Interactive charts
- [ ] Predictive analytics
- [ ] Scheduled reports per organization
- [ ] Custom recipient lists

---

## ðŸ“š Related Documentation

- **Feature Flags:** [FEATURE-FLAGS-NOTIFICATIONS.md](./FEATURE-FLAGS-NOTIFICATIONS.md)
- **Email Setup:** [EXTERNAL-SERVICES-SETUP.md](./EXTERNAL-SERVICES-SETUP.md)
- **Architecture:** [02-ARCHITECTURE.md](./02-ARCHITECTURE.md)
- **Features:** [03-FEATURES.md](./03-FEATURES.md)

---

**Last Updated:** December 23, 2025  
**Status:** Production Ready
