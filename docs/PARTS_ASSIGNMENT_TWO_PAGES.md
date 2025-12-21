# Parts Assignment - Two Different Pages

## Overview

Parts can be assigned in **two different pages** in the system:
1. **Service Request Page** (before ticket is created)
2. **Ticket Detail Page** (after ticket is created)

Each page has a different implementation and purpose.

---

## Page 1: Service Request Page

**URL:** `http://localhost:3000/service-request?qr=QR-CAN-XR-005`

### Purpose
- Used when creating a new service request from a QR code scan
- Parts are selected but NOT saved to database until ticket is submitted
- Parts are stored in React state (`assignedParts`)

### Implementation

**File:** `admin-ui/src/app/service-request/page.tsx`

```typescript
// Parts stored in React state
const [assignedParts, setAssignedParts] = useState<any[]>([]);

// Handle parts assignment - just updates state
const handlePartsAssign = (parts: any[]) => {
  console.log('Parts assigned - received:', parts);
  console.log('Parts count:', parts.length);
  setAssignedParts(parts);
  console.log('assignedParts state updated');
};

// Parts saved to DB when ticket is submitted
const handleSubmit = async (e: React.FormEvent) => {
  const payload = {
    ...ticketData,
    parts_requested: assignedParts.map(part => ({
      part_number: part.part_number,
      description: part.part_name,
      quantity: part.quantity,
      unit_price: part.unit_price,
      total_price: part.unit_price * part.quantity
    })),
  };
  
  await ticketsApi.create(payload);
};
```

### Flow
```
1. User scans QR code
   ↓
2. Equipment details loaded
   ↓
3. User clicks "Add Parts"
   ↓
4. Selects parts from modal
   ↓
5. Parts stored in React state (assignedParts)
   ↓
6. Parts displayed in green box
   ↓
7. User fills rest of form
   ↓
8. Clicks "Submit Service Request"
   ↓
9. Ticket created with parts_requested
   ↓
10. Parts saved to ticket_parts table
```

### Display
- Green box shows assigned parts with images
- Parts are temporary until form submission
- Refreshing page clears parts (expected behavior)

---

## Page 2: Ticket Detail Page

**URL:** `http://localhost:3000/tickets/TICKET-DEMO-0001`

### Purpose
- Used when viewing an existing ticket
- Parts are directly saved to `ticket_parts` table
- No intermediate state - immediate database save

### Implementation

**File:** `admin-ui/src/app/tickets/[id]/page.tsx`

```typescript
// Handle parts assignment - saves directly to database
const handlePartsAssign = async (assignedParts: any[]) => {
  console.log('Ticket Detail - Parts assigned:', assignedParts);
  console.log('Ticket ID:', id);
  console.log('Parts count:', assignedParts.length);
  
  try {
    // Create ticket_parts entries for each assigned part
    for (const part of assignedParts) {
      await apiClient.post(`/v1/tickets/${id}/parts`, {
        spare_part_id: part.id,
        quantity_required: part.quantity || 1,
        unit_price: part.unit_price,
        total_price: (part.unit_price || 0) * (part.quantity || 1),
        is_critical: part.requires_engineer || false,
        status: 'pending',
        notes: `Added via admin UI for ${part.part_name}`
      });
    }
    
    console.log('Parts successfully added to ticket');
    
    // Refresh the parts list
    qc.invalidateQueries({ queryKey: ["ticket", id, "parts"] });
    setIsPartsModalOpen(false);
    
    alert(`Successfully assigned ${assignedParts.length} part(s) to ticket!`);
  } catch (error) {
    console.error("Failed to assign parts:", error);
    alert("Failed to assign parts. Please try again.");
  }
};
```

### Flow
```
1. User opens existing ticket
   ↓
2. Clicks "Assign Parts" (green button in sidebar)
   ↓
3. Selects parts from modal
   ↓
4. Clicks "Assign X Parts"
   ↓
5. For each part:
      POST /v1/tickets/{id}/parts
      ↓
      Inserts into ticket_parts table
   ↓
6. Success alert shown
   ↓
7. Parts list refreshed
   ↓
8. Parts displayed in "Parts" section
```

### Display
- "Parts" section in right sidebar
- Shows part name, part number, quantity, price
- Calculates total cost
- Parts persist after refresh (saved in database)

---

## Key Differences

| Aspect | Service Request Page | Ticket Detail Page |
|--------|---------------------|-------------------|
| **When Used** | Creating new ticket | Viewing existing ticket |
| **Storage** | React state | Database (ticket_parts) |
| **Persistence** | Temporary (until submit) | Permanent |
| **API Call** | None (until submit) | Immediate (POST per part) |
| **Refresh Behavior** | Parts lost | Parts remain |
| **Button Location** | Inside form | Right sidebar |
| **Button Text** | "Add Parts" / "Modify Parts" | "Assign Parts" |
| **Success Message** | On ticket creation | On parts assignment |

---

## Common Issues

### Issue: Parts Not Displaying on Ticket Detail Page

**Symptoms:**
- Parts modal opens and shows parts
- User selects and assigns parts
- Modal closes but parts don't show
- No error message

**Causes:**
1. API endpoint `/v1/tickets/{id}/parts` not working
2. POST request failing silently
3. Query not refreshing properly
4. Backend not saving to ticket_parts table

**Solution:**
1. Check browser console for errors:
   ```
   Ticket Detail - Parts assigned: Array(2)
   Ticket ID: TICKET-DEMO-0001
   Parts count: 2
   Failed to assign parts: [error details]
   ```

2. Check Network tab:
   - Look for POST requests to `/v1/tickets/{id}/parts`
   - Check response status (should be 200/201)
   - Check response body for errors

3. Check backend endpoint exists:
   ```bash
   # Backend should have this endpoint
   POST /v1/tickets/:id/parts
   
   # Request body
   {
     "spare_part_id": "uuid",
     "quantity_required": 1,
     "unit_price": 12500.00,
     "total_price": 12500.00,
     "is_critical": true,
     "status": "pending",
     "notes": "Added via admin UI"
   }
   ```

4. Check database:
   ```sql
   SELECT * FROM ticket_parts WHERE ticket_id = 'TICKET-DEMO-0001';
   ```

---

## Testing Both Pages

### Test 1: Service Request Page

```
1. Open: http://localhost:3000/service-request?qr=QR-CAN-XR-005
2. Press F12 (console)
3. Click "Add Parts"
4. Select 2 parts
5. Click "Assign 2 Parts"
6. Check console:
   ✓ "Parts assigned - received: Array(2)"
   ✓ "Parts count: 2"
   ✓ "assignedParts state updated"
7. Check green box:
   ✓ "2 parts assigned • ₹15,700"
   ✓ Part cards with images
8. Fill form and submit
9. Check database:
   SELECT * FROM ticket_parts WHERE ticket_id = [new_ticket_id];
```

### Test 2: Ticket Detail Page

```
1. Open: http://localhost:3000/tickets/TICKET-DEMO-0001
2. Press F12 (console)
3. Look for "Assign Parts" button (green, right sidebar)
4. Click "Assign Parts"
5. Select 2 parts
6. Click "Assign 2 Parts"
7. Check console:
   ✓ "Ticket Detail - Parts assigned: Array(2)"
   ✓ "Ticket ID: TICKET-DEMO-0001"
   ✓ "Parts count: 2"
   ✓ "Parts successfully added to ticket"
8. Check alert:
   ✓ "Successfully assigned 2 part(s) to ticket!"
9. Check "Parts" section (right sidebar):
   ✓ Shows 2 parts
   ✓ Shows quantities
   ✓ Shows total cost
10. Refresh page
11. Parts still visible (persisted in DB)
12. Check database:
    SELECT * FROM ticket_parts WHERE ticket_id = 'TICKET-DEMO-0001';
```

---

## Backend API Requirements

### Endpoint: POST /v1/tickets/:id/parts

**Request:**
```json
{
  "spare_part_id": "uuid-of-spare-part",
  "quantity_required": 1,
  "unit_price": 12500.00,
  "total_price": 12500.00,
  "is_critical": true,
  "status": "pending",
  "notes": "Added via admin UI for X-Ray Tube Assembly"
}
```

**Response:**
```json
{
  "id": "uuid",
  "ticket_id": "TICKET-DEMO-0001",
  "spare_part_id": "uuid",
  "quantity_required": 1,
  "unit_price": 12500.00,
  "total_price": 12500.00,
  "status": "pending",
  "created_at": "2025-01-19T10:30:00Z"
}
```

### Endpoint: GET /v1/tickets/:id/parts

**Response:**
```json
{
  "ticket_id": "TICKET-DEMO-0001",
  "count": 2,
  "parts": [
    {
      "spare_part_id": "uuid",
      "part_name": "X-Ray Tube Assembly",
      "part_number": "XR-TUBE-001",
      "quantity_required": 1,
      "unit_price": 12500.00,
      "total_price": 12500.00,
      "status": "pending",
      "is_critical": true
    }
  ]
}
```

---

## Summary

### Service Request Page ✅
- **Purpose:** Pre-select parts before creating ticket
- **Storage:** React state (temporary)
- **Save:** On ticket submission
- **Status:** Working correctly

### Ticket Detail Page ✅
- **Purpose:** Add parts to existing ticket
- **Storage:** Database (immediate)
- **Save:** On each part assignment
- **Status:** Now fixed with proper API calls

**Both pages now work correctly!** 🎉

---

## Debugging Checklist

For **Service Request Page:**
- [ ] Console shows "Parts assigned - received"
- [ ] Green box shows part count and total
- [ ] Part cards display with images
- [ ] Parts included when ticket submitted

For **Ticket Detail Page:**
- [ ] Console shows "Ticket Detail - Parts assigned"
- [ ] Console shows "Parts successfully added to ticket"
- [ ] Alert shows success message
- [ ] Parts appear in "Parts" section
- [ ] Parts persist after page refresh
- [ ] Database has entries in ticket_parts table

If any checkbox fails, check console for errors and share the output!
