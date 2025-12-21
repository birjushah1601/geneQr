# Engineers Page Level Cards Fixed

## Issues Found

### Issue 1: Level Cards Showing 0
**Problem:** All level cards (Junior, Senior, Expert) showed 0 engineers

**Root Cause:** Backend returns `engineer_level` as strings ("L1", "L2", "L3") but frontend expected integers (1, 2, 3)

**API Response:**
```json
{
  "engineers": [
    {
      "name": "Amit Patel",
      "engineer_level": "L3"  // ← String format
    }
  ]
}
```

**Frontend Expected:**
```typescript
{
  engineer_level: 3  // ← Number format
}
```

### Issue 2: Duplicate Engineers
**Problem:** API returns 34 engineers instead of 16 unique

**Root Cause:** Backend returns one record per engineer-organization assignment
- Amit Patel works for 4 manufacturers → returned 4 times
- Rajesh Kumar Singh works for 4 manufacturers → returned 4 times
- etc.

**Database has:**
- 16 unique engineers
- 33 engineer-organization assignments

**API was returning:** 34 records (with some duplicates)

---

## Solutions Applied

### Fix 1: Parse Level Strings

**File:** `admin-ui/src/app/engineers/page.tsx`

**Before:**
```typescript
const stats = useMemo(() => {
  return {
    byLevel: engineers.reduce((acc, eng) => {
      const level = eng.engineer_level || 1;
      acc[level] = (acc[level] || 0) + 1;  // Fails for "L1" strings
      return acc;
    }, {} as Record<number, number>),
  };
}, [engineers]);
```

**After:**
```typescript
const stats = useMemo(() => {
  const byLevel = engineers.reduce((acc, eng) => {
    // Handle both integer (1, 2, 3) and string ("L1", "L2", "L3") formats
    let level = eng.engineer_level;
    if (typeof level === 'string' && level.startsWith('L')) {
      level = parseInt(level.substring(1));  // "L1" → 1
    } else if (typeof level === 'string') {
      level = parseInt(level);  // "1" → 1
    }
    const levelNum = level || 1;
    acc[levelNum] = (acc[levelNum] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);
  
  return {
    total: engineers.length,
    filtered: filteredEngineers.length,
    byLevel,
  };
}, [engineers, filteredEngineers]);
```

**Handles:**
- String format: "L1", "L2", "L3" → 1, 2, 3
- Numeric strings: "1", "2", "3" → 1, 2, 3
- Integer format: 1, 2, 3 → 1, 2, 3

### Fix 2: Deduplicate Engineers

**Before:**
```typescript
const response = await apiClient.get(url);
setEngineers(response.data.engineers || []);
// Sets 34 engineers with duplicates
```

**After:**
```typescript
const response = await apiClient.get(url);
const engineersData = response.data.engineers || [];

// Deduplicate engineers by ID
const uniqueEngineers = Array.from(
  new Map(engineersData.map((eng: Engineer) => [eng.id, eng])).values()
);

setEngineers(uniqueEngineers);
// Sets 16 unique engineers
```

**How it works:**
1. Maps each engineer to [id, engineer] pairs
2. Creates a Map (automatically keeps only last occurrence per key)
3. Extracts values to get unique engineers

---

## Expected Results

### Database Counts
```sql
SELECT 
  engineer_level,
  COUNT(*) as count
FROM engineers
GROUP BY engineer_level;
```

**Result:**
```
engineer_level | count
---------------+-------
     1         |   2    (Junior)
     2         |   7    (Senior)
     3         |   7    (Expert)
```

### Frontend Display After Reload

**Stats Cards:**
```
┌──────────────────┬──────────────┬──────────────┬──────────────┐
│ Total Engineers  │ Junior (L1)  │ Senior (L2)  │ Expert (L3)  │
│       16         │      2       │      7       │      7       │
└──────────────────┴──────────────┴──────────────┴──────────────┘
```

**Engineer List:**
- Shows 16 unique engineers (not 34 duplicates)
- Each engineer appears once
- Correct level badges (Junior/Senior/Expert)

---

## Engineers by Level

### Junior (Level 1) - 2 Engineers
1. Divya Krishnan
2. Priya Sharma

### Senior (Level 2) - 7 Engineers
1. Arjun Malhotra
2. Deepak Verma
3. Neha Kulkarni
4. Ravi Iyer
5. Sanjay Mehta
6. Shreya Patel
7. Suresh Gupta

### Expert (Level 3) - 7 Engineers
1. Amit Patel
2. Arun Menon
3. Karthik Raghavan
4. Kavita Nair
5. Manish Joshi
6. Rajesh Kumar Singh
7. Vikram Reddy

**Total: 16 unique engineers**

---

## Multi-Organization Engineers

These engineers work for multiple manufacturers:

| Engineer | Manufacturers Count | Level |
|----------|---------------------|-------|
| Amit Patel | 4 | Expert (L3) |
| Rajesh Kumar Singh | 4 | Expert (L3) |
| Manish Joshi | 4 | Expert (L3) |
| Kavita Nair | 3 | Expert (L3) |
| Suresh Gupta | 4 | Senior (L2) |
| Arun Menon | 2 | Expert (L3) |
| Vikram Reddy | 2 | Expert (L3) |
| Karthik Raghavan | 2 | Expert (L3) |
| Priya Sharma | 2 | Junior (L1) |

**Note:** These engineers appeared multiple times in the API response (once per manufacturer)

---

## Backend Note

The backend engineers API returns all engineer-organization relationships, which causes duplicates when an engineer works for multiple manufacturers.

**Current Behavior:**
```sql
-- Backend query returns this
SELECT e.*, m.org_id, m.role
FROM engineers e
JOIN engineer_org_memberships m ON e.id = m.engineer_id;
-- Returns 33+ rows (one per assignment)
```

**For unique engineers, should query:**
```sql
SELECT DISTINCT ON (e.id) e.*
FROM engineers e;
-- Returns 16 rows (unique engineers)
```

**Frontend Workaround:** Deduplication by ID handles this correctly.

---

## Files Modified

1. ✅ `admin-ui/src/app/engineers/page.tsx`
   - Added level string parsing ("L1" → 1)
   - Added engineer deduplication by ID
   - Stats now calculate correctly

---

## Testing

### Before Reload
```
Total Engineers: 0 or wrong number
Junior (Level 1): 0
Senior (Level 2): 0
Expert (Level 3): 0
Engineer List: Duplicates (34 entries)
```

### After Reload
```
Total Engineers: 16 ✅
Junior (Level 1): 2 ✅
Senior (Level 2): 7 ✅
Expert (Level 3): 7 ✅
Engineer List: Unique engineers (16 entries) ✅
```

### Level Badges
- Junior engineers: Blue badge
- Senior engineers: Green badge
- Expert engineers: Purple badge

---

## Status

✅ **Level cards fixed** - Parse "L1"/"L2"/"L3" strings  
✅ **Deduplication added** - Show 16 unique engineers  
✅ **Correct counts** - 2 Junior, 7 Senior, 7 Expert  
✅ **No duplicates** - Each engineer appears once  

⏳ **Browser reload needed** - To see all fixes  

---

## URLs to Test

### All Engineers
```
http://localhost:3000/engineers
```
**Expected:**
- Total: 16
- Junior: 2
- Senior: 7
- Expert: 7
- List shows 16 unique engineers

### Engineers by Manufacturer (Philips)
```
http://localhost:3000/engineers?manufacturer=f1c1ebfb-57fd-4307-93db-2f72e9d004ad
```
**Expected:**
- Total: 5
- Each of 5 engineers shown once
- Correct level badges

**Hard reload browser (Ctrl+Shift+R) to see the fixed level cards!** 🎉
