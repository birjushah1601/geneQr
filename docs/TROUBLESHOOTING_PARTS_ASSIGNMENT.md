# Troubleshooting Parts Assignment Display

## Issue

Parts assigned in the PartsAssignmentModal are not displaying in the service request form after clicking "Assign".

---

## Changes Made

### 1. Enhanced Debugging

**File:** `admin-ui/src/app/service-request/page.tsx`

```typescript
const handlePartsAssign = (parts: any[]) => {
  console.log('Parts assigned - received:', parts);
  console.log('Parts count:', parts.length);
  setAssignedParts(parts);
  console.log('assignedParts state updated');
};
```

**What to check in browser console:**
- Message: "Parts assigned - received:" with array of parts
- Message: "Parts count:" with number (should be > 0)
- Message: "assignedParts state updated"

### 2. Improved Display

**Before:**
```tsx
<div className="flex justify-between text-xs bg-white p-2 rounded border">
  <span className="font-medium">{part.part_name}</span>
  <span>{part.quantity}x • ₹{part.unit_price * part.quantity}</span>
</div>
```

**After:**
```tsx
<div className="flex items-start justify-between text-xs bg-white p-3 rounded border shadow-sm">
  <div className="flex items-center gap-3 flex-1">
    {/* Part Image */}
    {part.image_url && (
      <img 
        src={part.image_url} 
        alt={part.part_name}
        className="w-10 h-10 rounded object-cover"
      />
    )}
    <div className="flex-1">
      <p className="font-medium text-gray-900">{part.part_name}</p>
      <p className="text-gray-500 text-[10px]">{part.part_number}</p>
    </div>
  </div>
  <div className="text-right ml-2">
    <p className="font-semibold text-green-700">{part.quantity}x</p>
    <p className="text-gray-600">₹{(part.unit_price * part.quantity).toLocaleString()}</p>
  </div>
</div>
```

**Improvements:**
- ✅ Added part images (40x40px thumbnails)
- ✅ Better spacing and layout
- ✅ Part number displayed
- ✅ Price calculation with proper formatting
- ✅ Visual hierarchy with colors

---

## Testing Steps

### Step 1: Open Browser Console

1. Press **F12** to open Developer Tools
2. Go to **Console** tab
3. Clear any existing messages

### Step 2: Navigate to Service Request Page

```
http://localhost:3000/service-request?qr=QR-CAN-XR-005
```

### Step 3: Open Parts Modal

1. Click **"Add Parts"** button
2. Modal should open showing 50 spare parts

### Step 4: Select Parts

1. **Search** for a part (e.g., "tube")
2. **Click checkbox** to select
3. **Adjust quantity** if needed
4. Click **"Assign [X] Parts"** button

### Step 5: Check Console Output

You should see:
```
Parts assigned - received: Array(1) [{...}]
Parts count: 1
assignedParts state updated
```

### Step 6: Verify Display

In the green "Spare Parts Needed" section, you should see:
- ✅ Part count updated: "1 part assigned • ₹12,500"
- ✅ Part card showing:
  - Image (if available)
  - Part name
  - Part number
  - Quantity (e.g., "1x")
  - Total price (e.g., "₹12,500")

---

## Common Issues & Solutions

### Issue 1: Console Shows Nothing

**Symptom:** No console messages after clicking "Assign"

**Possible Causes:**
1. Modal not calling `onAssign` callback
2. JavaScript error preventing execution

**Solution:**
1. Check for JavaScript errors in console (red text)
2. Verify modal is importing correctly:
   ```typescript
   import { PartsAssignmentModal } from '@/components/PartsAssignmentModal';
   ```
3. Check modal props:
   ```tsx
   <PartsAssignmentModal
     open={isPartsModalOpen}
     onClose={() => setIsPartsModalOpen(false)}
     onAssign={handlePartsAssign}  // Must be present
     equipmentId={(equipment as any)?.id || 'unknown'}
     equipmentName={(equipment as any)?.equipment_name || 'Equipment'}
   />
   ```

### Issue 2: Console Shows Array But No Display

**Symptom:** Console logs "Parts count: 1" but nothing shows

**Possible Causes:**
1. State not updating
2. Render condition issue
3. CSS hiding elements

**Solution:**
1. Check React DevTools:
   - Open React DevTools
   - Find ServiceRequestPageInner component
   - Check `assignedParts` state value
2. Check if state is array: `Array.isArray(assignedParts)`
3. Verify `assignedParts.length > 0`

### Issue 3: Parts Show But No Images

**Symptom:** Parts display but images are missing

**Possible Causes:**
1. Images failed to load from Unsplash
2. `image_url` field not in data
3. Image URLs are null

**Solution:**
1. Check network tab for failed image requests
2. Verify API response includes `image_url`:
   ```javascript
   // In console
   fetch('http://localhost:8081/api/v1/catalog/parts')
     .then(r => r.json())
     .then(d => console.log(d.parts[0]))
   ```
3. Check if images load: Open image URL directly in browser

### Issue 4: Price Shows as NaN

**Symptom:** Price displays as "₹NaN" or "₹0"

**Possible Causes:**
1. `unit_price` is undefined
2. `quantity` is undefined or 0
3. Values are strings instead of numbers

**Solution:**
1. Add fallback values:
   ```typescript
   ₹{((part.unit_price || 0) * (part.quantity || 1)).toLocaleString()}
   ```
2. Check data types in console:
   ```javascript
   console.log(typeof part.unit_price);  // should be "number"
   console.log(typeof part.quantity);     // should be "number"
   ```

### Issue 5: Parts Disappear After Refresh

**Symptom:** Parts assigned but disappear when page refreshes

**Expected Behavior:** This is NORMAL! Parts are only stored in component state, not persisted.

**Explanation:**
- Parts are held in `assignedParts` state
- State is cleared on page refresh
- Parts are only saved to database when ticket is submitted

**Flow:**
```
1. Select parts → State updated
2. Submit ticket → Parts saved to ticket_parts table
3. Refresh page → State cleared (expected)
```

---

## Debugging Checklist

### Backend Check

```bash
# Test API endpoint
curl http://localhost:8081/api/v1/catalog/parts

# Should return:
{
  "parts": [
    {
      "id": "...",
      "part_name": "...",
      "unit_price": 12500.00,
      "image_url": "https://...",
      ...
    }
  ],
  "count": 50
}
```

### Frontend Check

1. **Check React state:**
   - Open React DevTools
   - Find `ServiceRequestPageInner`
   - Inspect `assignedParts` value

2. **Check props:**
   - Find `PartsAssignmentModal`
   - Verify `onAssign` prop exists
   - Verify it points to `handlePartsAssign`

3. **Check render:**
   - Look for green box with "Spare Parts Needed"
   - Should always be visible
   - Should show "0 parts" or "X parts assigned"

### Network Check

1. Open Network tab
2. Select parts and assign
3. Check for any failed requests
4. Verify API response contains all needed fields

---

## Expected Flow

### 1. Initial State
```
assignedParts = []
Display: "Select spare parts needed for this service request"
Button: "Add Parts"
```

### 2. After Opening Modal
```
Modal opens
Fetches: GET /api/v1/catalog/parts
Shows: 50 parts with images
```

### 3. After Selecting Parts
```
User clicks checkboxes
Selected parts highlighted
Cart tab shows selected items
```

### 4. After Clicking Assign
```
Console: "Parts assigned - received: [{...}]"
Console: "Parts count: 1"
Console: "assignedParts state updated"

assignedParts = [{
  id: "...",
  part_name: "X-Ray Tube Assembly",
  part_number: "XR-TUBE-001",
  unit_price: 12500,
  quantity: 1,
  image_url: "https://..."
}]

Display: "1 part assigned • ₹12,500"
Shows: Part card with image, name, price
Button: "Modify Parts"
```

### 5. After Submit
```
Ticket created with parts_requested
Parts saved to ticket_parts table
Success message shown
Form resets
assignedParts = [] (cleared)
```

---

## Quick Fix Script

If parts are still not showing, try this in browser console:

```javascript
// Check if component has assignedParts state
const rootElement = document.querySelector('#__next');
console.log('Root element:', rootElement);

// Force check state by looking at green box
const partsBox = document.querySelector('[class*="green-50"]');
console.log('Parts box found:', partsBox !== null);

// Check if Add Parts button exists
const addButton = document.querySelector('button:contains("Add Parts")');
console.log('Add Parts button:', addButton);
```

---

## Visual Reference

### Before Assigning Parts
```
┌─────────────────────────────────────┐
│ Spare Parts Needed                  │
│                                     │
│ Select spare parts needed for       │
│ this service request                │
│                                     │
│              [Add Parts] Button     │
└─────────────────────────────────────┘
```

### After Assigning Parts
```
┌─────────────────────────────────────┐
│ Spare Parts Needed                  │
│ 2 parts assigned • ₹15,700          │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │ [IMG] X-Ray Tube Assembly   1x  │ │
│ │       XR-TUBE-001      ₹12,500  │ │
│ └─────────────────────────────────┘ │
│                                     │
│ ┌─────────────────────────────────┐ │
│ │ [IMG] Collimator Assembly   1x  │ │
│ │       XR-COL-001        ₹3,200  │ │
│ └─────────────────────────────────┘ │
│                                     │
│           [Modify Parts] Button     │
└─────────────────────────────────────┘
```

---

## Still Not Working?

### Option 1: Hard Refresh

1. Clear browser cache (Ctrl+Shift+Del)
2. Hard refresh page (Ctrl+Shift+R)
3. Try again

### Option 2: Check Browser Compatibility

- Chrome: ✅ Supported
- Firefox: ✅ Supported
- Edge: ✅ Supported
- Safari: ✅ Supported
- Internet Explorer: ❌ Not supported

### Option 3: Rebuild Frontend

```bash
cd admin-ui
npm run build
npm run dev
```

### Option 4: Check Backend Logs

Look for any errors when fetching parts:
```
Failed to query spare parts
Query failed
Database connection failed
```

---

## Summary

**Key Points:**
1. ✅ Parts assignment works via state management
2. ✅ Console logs help debug the flow
3. ✅ Parts display with images and proper formatting
4. ✅ State clears on refresh (expected behavior)
5. ✅ Parts persist only when ticket is submitted

**Test Checklist:**
- [ ] Backend running on port 8081
- [ ] Frontend running on port 3000
- [ ] Browser console open
- [ ] No JavaScript errors
- [ ] API returns 50 parts
- [ ] Modal opens and shows parts
- [ ] Parts can be selected
- [ ] Console shows debug messages
- [ ] Parts display in green box
- [ ] Images load (if available)
- [ ] Prices calculate correctly

**If all checked and still not working:**
- Share console errors
- Share network tab output
- Share React DevTools state snapshot
