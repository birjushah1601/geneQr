# Service Request Page Enhancements

## Overview
Documentation for the enhanced service request page (QR code workflow) with optional contact fields for receiving ticket updates.

---

## Feature: Optional Contact Fields

### Purpose
Allow users scanning QR codes to provide contact information for receiving updates about their service requests via email or SMS/WhatsApp.

---

## Implementation

### Frontend Changes

**File:** `admin-ui/src/app/service-request/page.tsx`

#### 1. State Management

```typescript
const [formData, setFormData] = useState({
  description: '',
  requestedBy: '',
  contactName: '',
  contactPhone: '',
  contactEmail: '',  // NEW
});
```

#### 2. Form Fields

**Email Field:**
```tsx
<div>
  <label htmlFor="contactEmail" className="block text-sm font-medium text-gray-700 mb-2">
    Email (Optional)
  </label>
  <input
    type="email"
    id="contactEmail"
    value={formData.contactEmail}
    onChange={(e) => setFormData({ ...formData, contactEmail: e.target.value })}
    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
    placeholder="your@email.com"
  />
  <p className="text-xs text-gray-500 mt-1">Get updates via email</p>
</div>
```

**Phone Number Field:**
```tsx
<div>
  <label htmlFor="contactPhone" className="block text-sm font-medium text-gray-700 mb-2">
    Phone Number (Optional)
  </label>
  <input
    type="tel"
    id="contactPhone"
    value={formData.contactPhone}
    onChange={(e) => setFormData({ ...formData, contactPhone: e.target.value })}
    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
    placeholder="+91-98765-43210"
  />
  <p className="text-xs text-gray-500 mt-1">Get updates via SMS/WhatsApp</p>
</div>
```

#### 3. Layout Structure

```
┌──────────────────────────────────────────────────────┐
│ Equipment Details (Auto-populated from QR)           │
├──────────────────────────────────────────────────────┤
│ Your Name * (Required)                               │
│ [Input Field]                                        │
├──────────────────────────────────────────────────────┤
│ Issue Description * (Required)                       │
│ [Textarea]                                           │
├──────────────────────────┬───────────────────────────┤
│ Email (Optional)         │ Phone Number (Optional)   │ ← NEW
│ [your@email.com]         │ [+91-98765-43210]         │
│ Get updates via email    │ Get updates via SMS...    │
├──────────────────────────┴───────────────────────────┤
│ Attach Photos/Audio/Documents (Optional)             │
│ [File Upload Area]                                   │
├──────────────────────────────────────────────────────┤
│ [AI Analysis Section]                                │
├──────────────────────────────────────────────────────┤
│ [Submit Service Request Button]                      │
└──────────────────────────────────────────────────────┘
```

**Two-Column Grid:**
```tsx
<div className="grid grid-cols-2 gap-4">
  <div>{/* Email field */}</div>
  <div>{/* Phone field */}</div>
</div>
```

#### 4. Form Submission

```typescript
const payload = {
  equipment_id: equipmentId,
  qr_code: equipment.qr_code,
  serial_number: equipment.serial_number,
  equipment_name: equipment.equipment_name,
  customer_id: equipment.customer_id,
  customer_name: equipment.customer_name,
  customer_phone: formData.contactPhone || '9999999999',  // Use provided or default
  customer_email: formData.contactEmail || undefined,     // NEW: Only if provided
  issue_category: 'breakdown',
  issue_description: formData.description,
  priority: 'medium',
  source: 'web',
  created_by: formData.requestedBy || 'web-user',
  // ... other fields
};

const created = await ticketsApi.create(payload);
```

#### 5. Form Reset

```typescript
setSuccess(true);
setFormData({ 
  description: '', 
  requestedBy: '', 
  contactName: '', 
  contactPhone: '', 
  contactEmail: ''  // Reset email field
});
```

---

## User Experience

### Field Characteristics

| Field | Type | Required | Default | Validation |
|-------|------|----------|---------|------------|
| Email | email | No | empty | HTML5 email validation |
| Phone | tel | No | empty | None (flexible format) |

### Helper Text

**Email Field:**
- Text: "Get updates via email"
- Color: Gray-500
- Size: Extra small (text-xs)

**Phone Number Field:**
- Text: "Get updates via SMS/WhatsApp"
- Color: Gray-500
- Size: Extra small (text-xs)

### Visual Design

**Input Styling:**
```css
- Full width (w-full)
- Padding: 12px horizontal, 8px vertical
- Border: 1px solid gray-300
- Border radius: 6px (rounded-md)
- Focus state: 2px blue-500 ring
- Font size: 14px (text-sm)
```

**Label Styling:**
```css
- Font weight: Medium (font-medium)
- Color: Gray-700
- Margin bottom: 8px
- Font size: 14px (text-sm)
```

---

## Backend Integration

### Expected Payload Structure

```json
{
  "equipment_id": "EQ-DEMO-0002",
  "qr_code": "QR-123456",
  "serial_number": "SN-789",
  "equipment_name": "Demo Equipment",
  "customer_id": "customer-id",
  "customer_name": "Hospital Name",
  "customer_phone": "+91-98765-43210",
  "customer_email": "user@example.com",
  "issue_category": "breakdown",
  "issue_description": "Equipment not working...",
  "priority": "medium",
  "source": "web",
  "created_by": "John Doe"
}
```

### Database Schema

**Table:** `service_tickets`

```sql
-- Existing fields
customer_phone VARCHAR(20),

-- May need to add (if not exists)
customer_email VARCHAR(255),

-- Constraints
-- customer_email is optional (can be NULL)
```

### Notification System (Future)

When contact info is provided:

**Email Notifications:**
- Ticket created confirmation
- Status updates (acknowledged, in progress, resolved)
- Engineer assigned notification
- Parts ordered notification
- Completion notification

**SMS/WhatsApp Notifications:**
- Ticket number and confirmation
- Engineer on the way
- Ticket status changes
- Quick links for tracking

---

## Usage Scenarios

### Scenario 1: User Provides Both Email and Phone

**Input:**
- Email: `john.doe@hospital.com`
- Phone: `+91-98765-43210`

**Result:**
- User receives email updates
- User receives SMS/WhatsApp updates
- Maximum communication coverage

### Scenario 2: User Provides Only Email

**Input:**
- Email: `john.doe@hospital.com`
- Phone: (empty)

**Result:**
- User receives email updates
- Phone defaults to `9999999999` (system default)
- Email-only communication

### Scenario 3: User Provides Only Phone

**Input:**
- Email: (empty)
- Phone: `+91-98765-43210`

**Result:**
- User receives SMS/WhatsApp updates
- Email not sent (undefined)
- Phone-only communication

### Scenario 4: User Provides Neither

**Input:**
- Email: (empty)
- Phone: (empty)

**Result:**
- Phone defaults to `9999999999`
- Email not sent
- No direct user communication
- Standard workflow continues

---

## Validation

### Frontend Validation

**Email Field:**
```html
<input type="email" ... />
```
- HTML5 validates email format
- Not required - can be empty
- Invalid format prevents submission

**Phone Field:**
```html
<input type="tel" ... />
```
- No automatic validation
- Accepts any format
- Flexible for international numbers

### Backend Validation

**Recommended Checks:**
1. Email format validation (if provided)
2. Phone number length check
3. Sanitize inputs to prevent XSS
4. Check for spam/abuse patterns

---

## Accessibility

### Labels
- All inputs have associated `<label>` elements
- Labels use `htmlFor` to link to input `id`
- Screen readers can identify fields

### Placeholder Text
- Email: `your@email.com`
- Phone: `+91-98765-43210`
- Shows expected format

### Helper Text
- Explains purpose of each field
- Low contrast (gray-500) to not distract
- Positioned below input

### Tab Order
1. Your Name
2. Issue Description
3. Email
4. Phone Number
5. File Upload
6. Submit Button

---

## Mobile Responsiveness

### Grid Behavior

**Desktop (≥768px):**
```css
grid-cols-2  /* Two columns side by side */
gap-4        /* 16px gap between columns */
```

**Mobile (<768px):**
```css
grid-cols-1  /* Stack vertically */
gap-4        /* 16px gap between fields */
```

### Input Sizing

- Full width on all screen sizes
- Touch-friendly height (py-2 = 8px padding)
- Minimum 44px tap target

---

## Testing Checklist

### Functional Tests
- [ ] Email field accepts valid email addresses
- [ ] Email field shows error for invalid format
- [ ] Phone field accepts various formats
- [ ] Phone field accepts international numbers
- [ ] Fields can be left empty
- [ ] Form submits with both fields filled
- [ ] Form submits with only email filled
- [ ] Form submits with only phone filled
- [ ] Form submits with neither filled
- [ ] Fields reset after successful submission

### Visual Tests
- [ ] Two-column layout on desktop
- [ ] Single-column layout on mobile
- [ ] Helper text displays correctly
- [ ] Focus state shows blue ring
- [ ] Placeholder text visible
- [ ] Labels aligned properly

### Integration Tests
- [ ] Email sent to backend in payload
- [ ] Phone sent to backend in payload
- [ ] Default phone used when empty
- [ ] Email undefined when empty
- [ ] Ticket created successfully
- [ ] Confirmation screen shows

---

## Future Enhancements

### Potential Improvements

1. **Phone Number Formatting**
   - Auto-format based on country code
   - Validate format for specific regions
   - International phone library

2. **Email Verification**
   - Send verification code
   - Confirm email before updates
   - Bounce detection

3. **Communication Preferences**
   - Checkbox: "Send email updates"
   - Checkbox: "Send SMS updates"
   - Preferred communication method

4. **Saved Contacts**
   - Remember user (browser localStorage)
   - Auto-fill on return visit
   - "Use same contact info"

5. **Rich Notifications**
   - Email templates with branding
   - SMS shortlinks for tracking
   - WhatsApp Business API integration

6. **Notification Settings**
   - Frequency preferences
   - Notification types selection
   - Quiet hours

---

## Related Features

### QR Code Workflow
- User scans QR code
- Equipment details auto-populated
- Contact fields for updates
- Creates ticket automatically

### Ticket Tracking
- Users can track via ticket ID
- Email/SMS contains tracking link
- Real-time status updates

### Communication System
- Email service (SendGrid, etc.)
- SMS gateway integration
- WhatsApp Business API

---

## Commit History

**Main Commit:**
- `8b218187 - feat(ui): Add optional contact fields to service request page`

**Changes:**
- Added email field (optional)
- Added phone number field (optional)
- Two-column layout implementation
- Helper text for both fields
- Updated form submission logic
- Updated form reset logic

---

## Troubleshooting

### Issue: Email Not Sent

**Check:**
1. Field has valid email format
2. Backend receives `customer_email` field
3. Email service configured
4. Email template exists

### Issue: Phone Defaults to 9999999999

**Expected Behavior:**
- This is correct when user doesn't provide phone
- Backend requires phone field
- Default prevents errors

**To Change:**
- Modify default value in payload
- Or make backend accept null/undefined

### Issue: Form Not Submitting

**Check:**
1. Required fields filled (Name, Description)
2. Email format valid (if provided)
3. No JavaScript errors
4. Network connectivity

---

## Security Considerations

### Data Privacy
- Email/phone are personal data
- Follow GDPR/privacy regulations
- Secure transmission (HTTPS)
- Encrypted storage

### Spam Prevention
- Rate limiting on form submission
- CAPTCHA for high-volume sites
- Block known spam patterns
- Monitor abuse

### Data Retention
- Define retention period
- Delete after ticket closure + X days
- User can request deletion
- Comply with privacy laws

---

## Performance

### Impact
- Minimal (2 additional form fields)
- No external API calls on page load
- Standard form validation

### Optimization
- Fields lazy-load validation
- No blocking operations
- Fast form submission

---

## Documentation Updates

**Files Updated:**
- `SERVICE-REQUEST-ENHANCEMENTS.md` (this file)
- `CHANGELOG-2026-02-05.md`
- Code comments in `page.tsx`

**Related Docs:**
- `QR-CODE-FUNCTIONALITY.md`
- `QUICK-ACCESS-GUIDE.md`

---

## Support

For questions or issues with this feature:
1. Check this documentation first
2. Review commit `8b218187`
3. Test with different input combinations
4. Verify backend receives fields correctly
