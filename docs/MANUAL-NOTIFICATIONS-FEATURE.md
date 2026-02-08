# Manual Notifications Feature

## Overview

This feature enables admins to manually send email notifications to customers from the ticket details page. It's designed as a **Phase 1** approach while the system matures, with feature flags to control automatic notifications in the future.

## Feature Strategy

### Phase 1 (NOW - Manual Control)
- ✅ **First ticket creation email** - Automatic when ticket is created
- ✅ **Manual notification button** - Admin clicks to send updates
- ❌ **Daily digest** - OFF (will enable later)
- ❌ **Auto status change emails** - OFF (will enable later)

### Phase 2 (FUTURE - When Mature)
- ✅ All automatic notifications enabled
- ✅ Daily digest cron job running
- ✅ Auto emails on every ticket update
- ⚠️ Manual button can be hidden (optional)

## Feature Flags

Add these to your `.env` file:

```env
# Ticket Creation Email (First email - enable when ready)
FEATURE_TICKET_CREATED_EMAIL=false

# Manual Notifications (Admin can send anytime - safe to enable)
FEATURE_MANUAL_NOTIFICATIONS=true

# Auto notifications on ticket updates (enable when mature)
FEATURE_AUTO_TICKET_UPDATES=false

# Daily digest of open tickets (enable when mature)
FEATURE_TICKET_DAILY_DIGEST=false

# Tracking Configuration
TICKET_TRACKING_BASE_URL=https://servqr.com/track
FRONTEND_BASE_URL=https://servqr.com
TRACKING_TOKEN_EXPIRY_DAYS=1825  # 5 years
```

## How It Works

### 1. Manual Notification Button

**Location:** Ticket Details Page (`/tickets/[id]`)

**What it does:**
- Admin clicks "Send Notification" button
- Modal opens asking for:
  - Customer email address
  - Message/update to send
- Email is sent immediately
- Success message shown

**When to use:**
- Notify customer of status change
- Request more information
- Provide update on repair progress
- Confirm parts arrival
- Any custom message needed

### 2. Tracking URLs

**Generated on ticket creation:**
- Every new ticket gets a secure tracking token
- Token valid for 5 years (configurable)
- Customer can view ticket status anytime
- No login required

**Tracking URL format:**
```
https://servqr.com/track/abc123...
```

### 3. Ticket Creation Email (Future)

When `FEATURE_TICKET_CREATED_EMAIL=true`:
- Automatic email sent when ticket is created
- Includes tracking URL
- Customer can bookmark for later

## API Endpoints

### Manual Notification
```
POST /api/v1/tickets/:id/notify
```

**Request:**
```json
{
  "email": "customer@example.com",
  "comment": "Your equipment has been repaired and is ready for pickup."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Notification sent successfully",
  "email": "customer@example.com"
}
```

### Public Ticket Tracking
```
GET /api/v1/track/:token
```

**Response:**
```json
{
  "ticket_number": "TKT-20260207-001",
  "status": "in_progress",
  "priority": "high",
  "equipment_name": "X-Ray Machine",
  "issue_description": "Machine not powering on",
  "created_at": "2026-02-07T10:00:00Z",
  "updated_at": "2026-02-07T11:00:00Z",
  "assigned_engineer": "John Doe",
  "comments": []
}
```

## Frontend Components

### Ticket Details Page Button
```tsx
<button
  onClick={() => setShowNotificationModal(true)}
  className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg"
>
  <Mail className="h-4 w-4" />
  Send Notification
</button>
```

### SendNotificationModal Component
- Already exists at `src/components/SendNotificationModal.tsx`
- Handles email input and validation
- Shows success/error messages
- Calls `/api/v1/tickets/:id/notify`

## Database Tables

### ticket_tracking_tokens
```sql
CREATE TABLE ticket_tracking_tokens (
    id UUID PRIMARY KEY,
    ticket_id VARCHAR(50) REFERENCES service_tickets(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### notification_log
```sql
CREATE TABLE notification_log (
    id UUID PRIMARY KEY,
    ticket_id VARCHAR(50) REFERENCES service_tickets(id),
    notification_type VARCHAR(50),  -- 'manual', 'ticket_created', 'daily_digest'
    recipient_email VARCHAR(255),
    sent_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50),  -- 'sent', 'failed'
    error_message TEXT
);
```

## Rollout Strategy

### Week 1: Testing Phase
```env
FEATURE_TICKET_CREATED_EMAIL=false
FEATURE_MANUAL_NOTIFICATIONS=true
FEATURE_AUTO_TICKET_UPDATES=false
FEATURE_TICKET_DAILY_DIGEST=false
```

**Actions:**
- Enable manual notifications only
- Train admins on how to use
- Monitor logs for errors
- Collect feedback

### Week 2-3: Ticket Creation Emails
```env
FEATURE_TICKET_CREATED_EMAIL=true
FEATURE_MANUAL_NOTIFICATIONS=true
```

**Actions:**
- Enable automatic first email
- Monitor email delivery rates
- Check customer feedback
- Fix any issues

### Week 4-6: Mature Phase
```env
FEATURE_AUTO_TICKET_UPDATES=true
FEATURE_TICKET_DAILY_DIGEST=true
```

**Actions:**
- Enable all auto notifications
- Daily digest at 9 AM
- Reduce manual notifications
- Consider hiding manual button (optional)

## Testing

### Test Manual Notification

1. Go to any ticket: `http://localhost:3000/tickets/[id]`
2. Click "Send Notification" button
3. Enter customer email
4. Enter message
5. Click Send
6. Check backend logs for success

### Test Tracking URL

1. Create a new ticket
2. Copy the tracking URL from success page
3. Open in incognito/private window
4. Should see ticket details without login
5. Verify token doesn't expire (5 years)

### Test Feature Flags

1. Set `FEATURE_MANUAL_NOTIFICATIONS=false`
2. Restart backend
3. Try to send notification
4. Should get "Manual notifications are disabled" error

## Future Enhancements

### Daily Digest Email (Phase 2)
- Runs daily at 9 AM
- Sends summary of open tickets
- Only for tickets in: `new`, `assigned`, `in_progress`
- Excludes: `resolved`, `closed`, `on_hold`, `cancelled`

### Auto Status Change Emails (Phase 2)
- Email when ticket status changes
- Email when engineer assigned
- Email when parts ordered
- Email when resolved

### WhatsApp Integration (Phase 3)
- Send notifications via WhatsApp
- Customer can reply via WhatsApp
- Updates added as comments

## Troubleshooting

### "Notification service not available"
- Check if backend is running
- Verify notification service is initialized in `module.go`
- Check logs for initialization errors

### "Manual notifications are disabled"
- Set `FEATURE_MANUAL_NOTIFICATIONS=true` in `.env`
- Restart backend

### Tracking URL shows "Invalid or expired"
- Check if token exists in database
- Verify token hasn't expired
- Check `expires_at` column

### Email not received
- Check SendGrid API key
- Verify email address is valid
- Check notification_log table for errors
- Check SendGrid dashboard for delivery status

## Support

For issues or questions:
1. Check backend logs: `backend-output.log`
2. Check notification_log table in database
3. Verify feature flags are set correctly
4. Check SendGrid dashboard for email delivery

## References

- SendNotificationModal: `admin-ui/src/components/SendNotificationModal.tsx`
- Backend Handler: `internal/service-domain/service-ticket/api/handler.go`
- Routes: `internal/service-domain/service-ticket/module.go`
- Migration: `database/migrations/008_ticket_notifications.sql`
