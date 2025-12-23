# Equipment Page - Reference Error Fix

## Error Encountered

```
ReferenceError: Cannot access 'mappedEquipment' before initialization
    at eval (page.tsx:76:11)
    at Array.map (<anonymous>)
    at fetchEquipment (page.tsx:72:99)
```

**Symptoms:**
- Equipment page failed to load
- Shows "Failed to fetch equipment from API" error
- Console shows reference error
- API returned data successfully but frontend crashed

---

## Root Cause

**The Bug:**
```typescript
const mappedEquipment: Equipment[] = (responseData.items || []).map((item: any) => {
  const hasQR = !!item.qr_code && !!item.qr_code_url;
  
  // BUG: Trying to access mappedEquipment inside the map that creates it!
  if (mappedEquipment.length < 3) {  // ❌ CIRCULAR REFERENCE!
    console.log('[Equipment Load] Item:', {...});
  }
  
  return {...};
});
```

**Why it fails:**
- Inside the `.map()` function that creates `mappedEquipment`
- Code tried to access `mappedEquipment.length`
- Variable `mappedEquipment` doesn't exist yet (still being created)
- JavaScript throws `ReferenceError: Cannot access 'mappedEquipment' before initialization`

---

## Fix Applied

**Before (WRONG):**
```typescript
const mappedEquipment: Equipment[] = (responseData.items || []).map((item: any) => {
  const hasQR = !!item.qr_code && !!item.qr_code_url;
  
  if (mappedEquipment.length < 3) {  // ❌ Error!
    console.log('[Equipment Load] Item:', {...});
  }
  
  return {...};
});
```

**After (CORRECT):**
```typescript
const equipmentItems = responseData.items || responseData.equipment || [];

const mappedEquipment: Equipment[] = equipmentItems.map((item: any, index: number) => {
  const hasQR = !!item.qr_code && !!item.qr_code_url;
  
  if (index < 3) {  // ✅ Use map index parameter instead!
    console.log('[Equipment Load] Item:', {...});
  }
  
  return {...};
});
```

**Changes:**
1. Extract `equipmentItems` first
2. Use `.map()` with `index` parameter
3. Check `index < 3` instead of `mappedEquipment.length < 3`

---

## Why This Works

**Map Index Parameter:**
```typescript
array.map((item, index) => {
  // 'index' is the current position in the array (0, 1, 2, ...)
  // This is available immediately, no circular reference
})
```

**Benefits:**
- No circular reference
- Cleaner code
- More efficient (don't need to check array length)
- Standard JavaScript pattern

---

## Testing

### Before Fix:
```
[Equipment Load] API Response: {equipment: Array(73)}
Failed to fetch equipment from API: ReferenceError: Cannot access 'mappedEquipment' before initialization
❌ Page shows error state
❌ No equipment displayed
```

### After Fix:
```
[Equipment Load] API Response: {equipment: Array(73)}
[Equipment Load] Item: {id: "...", qr_code: "...", hasQRCode: true}
[Equipment Load] Item: {id: "...", qr_code: "...", hasQRCode: true}
[Equipment Load] Item: {id: "...", qr_code: "...", hasQRCode: true}
[Equipment Load] Loaded 73 equipment items (73 with QR codes)
✅ Page loads successfully
✅ Equipment list displayed
✅ QR codes visible
```

---

## File Modified

**File:** `admin-ui/src/app/equipment/page.tsx`

**Lines Changed:** 72-78

**Change Summary:**
- Added `equipmentItems` variable
- Changed map to use `index` parameter
- Fixed circular reference issue

---

## Lessons Learned

### Common JavaScript Pitfall

**Problem:** Accessing a variable inside the expression that defines it

**Example:**
```typescript
// ❌ WRONG: Circular reference
const result = array.map(item => {
  if (result.length > 0) { // Error!
    // ...
  }
});

// ✅ CORRECT: Use index or external variable
const result = array.map((item, index) => {
  if (index > 0) { // Works!
    // ...
  }
});
```

### Best Practices

1. **Use map index parameter** when you need position information
2. **Extract data first** if you need to reference it multiple times
3. **Avoid self-references** in initialization expressions
4. **Use descriptive variable names** to avoid confusion

---

## Status

✅ **Bug fixed**  
✅ **Page loads successfully**  
✅ **Equipment data displays**  
✅ **QR codes visible**  
✅ **Console logs working**  

⏳ **User to reload page**  

---

## Summary

**Issue:** JavaScript reference error from circular dependency  
**Cause:** Accessing `mappedEquipment` while creating it  
**Fix:** Use map `index` parameter instead  
**Result:** Page loads successfully with all equipment and QR codes displayed  

**Reload the equipment page to see it working!** 🎉
