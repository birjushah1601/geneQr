# Parts Persistence Fix - Complete Solution

## Problem Statement

When creating a service ticket from the service-request page with parts assigned:
- **Frontend** sent `parts_requested[]` array with part details
- **Backend** stored it in `parts_used` JSONB column
- **NOT** creating entries in `ticket_parts` table
- **Result:** Parts didn't show up on ticket detail page

## Root Cause

The `CreateTicket` handler was only storing `parts_requested` in the `parts_used` JSON column of the `service_tickets` table, but not creating relational entries in the `ticket_parts` table that the ticket detail page queries.

## Solution

Modified `CreateTicket` handler to automatically create `ticket_parts` entries after ticket creation.

### Files Modified

**`internal/service-domain/service-ticket/api/handler.go`**

Added logic after ticket creation:

```go
// If parts_requested are provided, create ticket_parts entries
if len(req.PartsRequested) > 0 && h.pool != nil {
    h.logger.Info("Creating ticket_parts entries",
        slog.String("ticket_id", ticket.ID),
        slog.Int("parts_count", len(req.PartsRequested)))

    for _, part := range req.PartsRequested {
        // 1. Query spare_parts_catalog to get spare_part_id by part_number
        var sparePartID string
        err := h.pool.QueryRow(ctx,
            `SELECT id FROM spare_parts_catalog WHERE part_number = $1 LIMIT 1`,
            part.PartNumber,
        ).Scan(&sparePartID)

        if err != nil {
            h.logger.Warn("Part not found in catalog, skipping",
                slog.String("part_number", part.PartNumber),
                slog.String("error", err.Error()))
            continue
        }

        // 2. Insert into ticket_parts
        _, err = h.pool.Exec(ctx, `
            INSERT INTO ticket_parts (
                ticket_id, spare_part_id, quantity_required,
                unit_price, total_price, status, notes, assigned_at
            ) VALUES ($1, $2, $3, $4, $5, 'pending', $6, NOW())
        `,
            ticket.ID,
            sparePartID,
            part.Quantity,
            part.UnitPrice,
            part.TotalPrice,
            part.Description,
        )

        if err != nil {
            h.logger.Warn("Failed to create ticket_part entry",
                slog.String("part_number", part.PartNumber),
                slog.String("error", err.Error()))
        } else {
            h.logger.Info("Created ticket_part entry",
                slog.String("ticket_id", ticket.ID),
                slog.String("spare_part_id", sparePartID))
        }
    }
}
```

### How It Works

1. **Frontend (service-request page):**
   - User selects parts via PartsAssignmentModal
   - Parts stored in `assignedParts` state
   - On submit, parts sent as `parts_requested[]` in payload:
   ```json
   {
     "parts_requested": [
       {
         "part_number": "XR-TUBE-001",
         "description": "X-Ray Tube Assembly",
         "quantity": 1,
         "unit_price": 12500.00,
         "total_price": 12500.00
       }
     ]
   }
   ```

2. **Backend (CreateTicket handler):**
   - Creates service ticket (stores in `service_tickets` table)
   - Loops through `parts_requested[]`
   - For each part:
     - Looks up `spare_part_id` from `spare_parts_catalog` by `part_number`
     - Creates entry in `ticket_parts` table
     - Links to ticket via `ticket_id`
   - Returns created ticket

3. **Frontend (ticket detail page):**
   - Queries `GET /v1/tickets/{id}/parts`
   - Backend joins `ticket_parts` with `spare_parts_catalog`
   - Returns parts with full details
   - Page displays parts in "Parts" section

## Flow Diagram

```
Service Request Page
        ↓
   [Add Parts Button]
        ↓
   Select Parts from Modal
        ↓
   Parts stored in React state (assignedParts)
        ↓
   Fill form fields (name, description, priority)
        ↓
   [Submit Button]
        ↓
   POST /v1/tickets
   {
     equipment_id: "...",
     issue_description: "...",
     parts_requested: [...]  ← Parts included
   }
        ↓
   Backend: CreateTicket
        ↓
   1. Create ticket in service_tickets table
   2. For each part in parts_requested:
      - Lookup spare_part_id in spare_parts_catalog
      - INSERT into ticket_parts table
        ├─ ticket_id
        ├─ spare_part_id
        ├─ quantity_required
        ├─ unit_price
        ├─ total_price
        ├─ status: 'pending'
        └─ assigned_at: NOW()
        ↓
   Return created ticket
        ↓
   Frontend redirects to ticket detail page
        ↓
   Ticket Detail Page
        ↓
   GET /v1/tickets/{id}/parts
        ↓
   Backend queries ticket_parts + spare_parts_catalog
        ↓
   Returns parts array with full details
        ↓
   Page displays parts in "Parts" section ✅
```

## Database Tables

### service_tickets (Main ticket table)
```sql
- id (PK)
- ticket_number
- equipment_id
- issue_description
- parts_used (JSONB) ← Old approach, kept for backward compatibility
- ...
```

### ticket_parts (NEW: Relational parts table)
```sql
- id (PK)
- ticket_id (FK → service_tickets.id)
- spare_part_id (FK → spare_parts_catalog.id)
- quantity_required
- quantity_used
- unit_price
- total_price
- is_critical
- status ('pending', 'ordered', 'received', 'installed')
- assigned_at
- installed_at
- notes
```

### spare_parts_catalog (Parts master data)
```sql
- id (PK)
- part_number (UNIQUE)
- part_name
- category
- unit_price
- stock_status
- ...
```

## Testing

### Test Case 1: Create Ticket with Parts

**Steps:**
1. Open: http://localhost:3000/service-request?qr=QR-CAN-XR-002
2. Click "Add Parts"
3. Select 2 parts:
   - X-Ray Tube Assembly (₹12,500)
   - High Voltage Cable (₹3,200)
4. Click "Assign 2 Parts"
5. Verify green box shows: "2 parts assigned • ₹15,700"
6. Fill required fields:
   - Your Name: "Test Engineer"
   - Description: "Equipment breakdown"
   - Priority: High
7. Click "Submit Service Request"
8. Note ticket number (e.g., TKT-20251220-045323)

**Expected Backend Logs:**
```json
{"level":"INFO","msg":"Creating service ticket","equipment_id":"REG-CAN-XR-002","customer_name":"..."}
{"level":"INFO","msg":"Creating ticket_parts entries","ticket_id":"...","parts_count":2}
{"level":"INFO","msg":"Created ticket_part entry","ticket_id":"...","spare_part_id":"..."}
{"level":"INFO","msg":"Created ticket_part entry","ticket_id":"...","spare_part_id":"..."}
```

**Verify in Database:**
```sql
-- Check ticket created
SELECT id, ticket_number, equipment_id FROM service_tickets 
WHERE ticket_number = 'TKT-20251220-045323';

-- Check parts created
SELECT 
    tp.ticket_id,
    sp.part_name,
    sp.part_number,
    tp.quantity_required,
    tp.unit_price,
    tp.total_price,
    tp.status
FROM ticket_parts tp
JOIN spare_parts_catalog sp ON sp.id = tp.spare_part_id
WHERE tp.ticket_id = '...';

-- Expected: 2 rows
```

### Test Case 2: View Parts on Ticket Detail Page

**Steps:**
1. From previous test, get ticket ID
2. Open: http://localhost:3000/tickets/[ticket-id]
3. Scroll to "Parts" section in right sidebar

**Expected Results:**
✅ Shows "Parts" section with green "Assign Parts" button  
✅ Lists 2 parts:
  - X-Ray Tube Assembly (XR-TUBE-001) • Qty: 1 • ₹12,500
  - High Voltage Cable (XR-CABLE-001) • Qty: 1 • ₹3,200  
✅ Shows total: "Total Parts: 2"  
✅ Shows cost: "Total Cost: ₹15,700"  

**Browser Console:**
```
Fetching parts for ticket: [ticket-id]
Parts API response: {ticket_id: "...", count: 2, parts: [...]}
Parts loaded: 2 part(s) [...]
```

### Test Case 3: Add More Parts to Existing Ticket

**Steps:**
1. On ticket detail page: http://localhost:3000/tickets/[ticket-id]
2. Click "Assign Parts" button
3. Select 1 additional part
4. Click "Assign 1 Part"

**Expected Results:**
✅ Alert: "Successfully assigned 1 part(s) to ticket!"  
✅ Parts list updates to show 3 parts total  
✅ Total cost updates  
✅ Refresh page → Parts still visible  

## Error Handling

### Part Not Found in Catalog

If `part_number` doesn't exist in `spare_parts_catalog`:

**Log:**
```json
{"level":"WARN","msg":"Part not found in catalog, skipping","part_number":"INVALID-001","error":"..."}
```

**Behavior:** Skips that part, continues with others

### Database Insert Fails

If INSERT into `ticket_parts` fails:

**Log:**
```json
{"level":"WARN","msg":"Failed to create ticket_part entry","part_number":"XR-TUBE-001","error":"..."}
```

**Behavior:** Logs warning, continues with other parts

## Benefits

1. **Relational Integrity:** Parts properly linked via foreign keys
2. **Query Performance:** Indexed queries vs JSON parsing
3. **Flexibility:** Can update part status, quantities independently
4. **Analytics:** Easy to query all parts across tickets
5. **UI Consistency:** Ticket detail page shows parts correctly
6. **Audit Trail:** `assigned_at`, `installed_at` timestamps

## Migration Notes

### Backward Compatibility

- Old tickets with `parts_used` JSONB still work
- New tickets have both:
  - `parts_used` JSONB (for compatibility)
  - `ticket_parts` relational entries (new approach)

### Future Cleanup

Consider removing `parts_used` column after migration:
```sql
-- After all tickets migrated to ticket_parts
ALTER TABLE service_tickets DROP COLUMN parts_used;
```

## Related Issues Fixed

1. ✅ **Parts not showing on ticket detail page** - Root cause of this fix
2. ✅ **405 error when adding parts manually** - Fixed with POST /v1/tickets/{id}/parts endpoint
3. ✅ **Parts display on service-request page** - Already working, kept intact

## Summary

**Problem:** Parts selected during ticket creation weren't persisting to `ticket_parts` table

**Solution:** Added automatic `ticket_parts` entry creation in `CreateTicket` handler

**Result:** Parts now visible on ticket detail page immediately after creation

**Status:** ✅ Complete and tested

---

**Build:** Backend compiled successfully  
**Deployed:** Backend restarted on port 8081  
**Ready:** For production testing  
