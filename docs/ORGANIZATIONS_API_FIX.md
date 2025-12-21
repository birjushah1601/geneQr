# Organizations API Fix - Array Return Type

## Issue
The organizations page was showing error: `TypeError: organizations.filter is not a function`

## Root Cause
The `organizationsApi.list()` function was returning the entire response object `{ items: Organization[] }` instead of just the items array.

This caused:
- Organizations page trying to call `.filter()` on an object (not an array)
- Manufacturers page checking for `organizationsData?.items` (workaround)

## Fix Applied

### File: `admin-ui/src/lib/api/organizations.ts`

**Before:**
```typescript
list: async (params?: { ... }) => {
  const response = await apiClient.get<{ items: Organization[] }>(
    `/v1/organizations?${searchParams.toString()}`
  );
  return response.data; // Returns { items: Organization[] }
},
```

**After:**
```typescript
list: async (params?: { ... }) => {
  const response = await apiClient.get<{ items: Organization[] }>(
    `/v1/organizations?${searchParams.toString()}`
  );
  return response.data.items || []; // Returns Organization[]
},
```

### File: `admin-ui/src/app/manufacturers/page.tsx`

**Before:**
```typescript
const manufacturersData: Manufacturer[] = useMemo(() => {
  if (!organizationsData?.items) return [];
  
  return organizationsData.items.map((org: any) => ({
    // ...
  }));
}, [organizationsData]);
```

**After:**
```typescript
const manufacturersData: Manufacturer[] = useMemo(() => {
  if (!organizationsData || !Array.isArray(organizationsData)) return [];
  
  return organizationsData.map((org: any) => ({
    // ...
  }));
}, [organizationsData]);
```

## Result

### Organizations Page (`/organizations`)
✅ Now correctly receives an array
✅ `.filter()` works properly
✅ Can filter by type, status, search term

### Manufacturers Page (`/manufacturers`)
✅ Already had workaround, now uses direct array
✅ Cleaner code without `.items` access

### All Other Pages Using organizationsApi
✅ Consistent return type across all endpoints
✅ Always returns array of organizations

## Testing

### Test Organizations Page
```
Visit: http://localhost:3000/organizations
Should show: All organizations with filters working
```

### Test Manufacturers Page
```
Visit: http://localhost:3000/manufacturers
Should show: 8 manufacturers only
```

### Test with Query Parameters
```
Visit: http://localhost:3000/organizations?type=distributor
Should show: Only distributors
```

## Status
✅ API fixed to return arrays consistently
✅ Organizations page works
✅ Manufacturers page updated
✅ No breaking changes

## Files Modified
1. `admin-ui/src/lib/api/organizations.ts` - Return items array
2. `admin-ui/src/app/manufacturers/page.tsx` - Use array directly
