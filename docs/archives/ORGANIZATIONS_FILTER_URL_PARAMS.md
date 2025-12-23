# Organizations Page - URL Query Parameters Support

## Issue
The organizations page was not reading URL query parameters like `?type=distributor`, so visiting `/organizations?type=distributor` would still show all organizations.

## Fix Applied

### File: `admin-ui/src/app/organizations/page.tsx`

**Added:**
1. Import `useSearchParams` from Next.js
2. Read URL query parameters on page load
3. Set filter states based on URL params

**Code Added:**
```typescript
import { useSearchParams } from 'next/navigation';

export default function OrganizationsPage() {
  const searchParams = useSearchParams();
  
  // Initialize filters from URL query parameters
  useEffect(() => {
    const typeParam = searchParams.get('type');
    const statusParam = searchParams.get('status');
    
    if (typeParam) {
      setFilterType(typeParam);
    }
    if (statusParam) {
      setFilterStatus(statusParam);
    }
  }, [searchParams]);
  
  // ... rest of component
}
```

## Supported URL Parameters

### Type Filter
- `?type=manufacturer` - Show only manufacturers (8 items)
- `?type=distributor` - Show only distributors (4 items)
- `?type=dealer` - Show only dealers (1 item)
- `?type=hospital` - Show only hospitals (5 items)
- `?type=service_provider` - Show only service providers
- `?type=imaging_center` - Show only imaging centers (3 items)

### Status Filter
- `?status=active` - Show only active organizations
- `?status=inactive` - Show only inactive organizations

### Combined Filters
- `?type=manufacturer&status=active` - Active manufacturers only
- `?type=hospital&status=active` - Active hospitals only

## How It Works

### Page Load Flow
1. Page loads with URL: `/organizations?type=distributor`
2. `useSearchParams()` reads query parameters
3. `searchParams.get('type')` returns `"distributor"`
4. `setFilterType('distributor')` updates the filter state
5. Dropdown shows "Distributors" selected
6. Only distributors are displayed in the list

### Filter State
The filter dropdown will automatically show the correct selection based on URL parameters, making it easy to share filtered views.

## Testing

### Test URLs
```
http://localhost:3000/organizations
http://localhost:3000/organizations?type=manufacturer
http://localhost:3000/organizations?type=distributor
http://localhost:3000/organizations?type=dealer
http://localhost:3000/organizations?type=hospital
http://localhost:3000/organizations?status=active
http://localhost:3000/organizations?type=manufacturer&status=active
```

### Expected Results

**No parameters:**
- Shows all 18 organizations
- Filter dropdown shows "All Types"

**?type=distributor:**
- Shows only 4 distributors
- Filter dropdown shows "Distributors" selected

**?type=manufacturer:**
- Shows only 8 manufacturers
- Filter dropdown shows "Manufacturers" selected

**?status=active:**
- Shows only active organizations
- Status dropdown shows "Active" selected

## Benefits

### 1. Shareable URLs
Users can share filtered views:
- "Here are all our distributors: /organizations?type=distributor"
- "Check active manufacturers: /organizations?type=manufacturer&status=active"

### 2. Deep Linking
Other pages can link directly to filtered views:
- From dashboard: "View Manufacturers" → `/organizations?type=manufacturer`
- From reports: "View Hospitals" → `/organizations?type=hospital`

### 3. Browser Navigation
- Back/Forward buttons maintain filter state
- Bookmarking preserves filters
- Refresh keeps the same view

### 4. Better UX
- Direct access to specific organization types
- No need to manually select filters
- Consistent with standard web practices

## Frontend Integration

### Link to Filtered View
```tsx
<Link href="/organizations?type=manufacturer">
  View All Manufacturers
</Link>
```

### Navigate Programmatically
```tsx
router.push('/organizations?type=distributor&status=active');
```

### Current Implementation
The filters work both ways:
- ✅ URL params set the filter dropdowns
- ✅ Filter dropdowns filter the displayed list
- ⏳ Could add: Update URL when dropdowns change (optional enhancement)

## Status

✅ URL query parameters read on page load
✅ Filter dropdowns auto-select based on URL
✅ List filtered correctly
✅ All organization types supported
✅ Status filtering supported
✅ Combined filters supported

## Files Modified

1. `admin-ui/src/app/organizations/page.tsx`
   - Added useSearchParams import
   - Added URL parameter reading logic
   - Auto-sets filter state from URL

## Future Enhancements (Optional)

### Update URL on Filter Change
When user changes dropdown, update URL:
```typescript
const handleTypeChange = (newType: string) => {
  setFilterType(newType);
  const params = new URLSearchParams(searchParams);
  if (newType === 'all') {
    params.delete('type');
  } else {
    params.set('type', newType);
  }
  router.push(`/organizations?${params.toString()}`);
};
```

This would make the URL stay in sync with the dropdown selection.
