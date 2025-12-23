# Manufacturers Page - Fixed to Show Only Manufacturers

## Issue
The /manufacturers page was showing all 18 organizations instead of filtering to show only manufacturers.

## Root Cause
Backend API was not filtering by 'type' query parameter. The ListOrgs method was ignoring the type filter.

## Fix Applied

### Backend Changes

**File 1:** internal/core/organizations/api/handler.go
- Added orgType and status query parameters extraction
- Pass parameters to repository

**File 2:** internal/core/organizations/infra/repository.go  
- Updated ListOrgs signature to accept orgType and status
- Built dynamic WHERE clause with filters
- Added fmt import for string formatting

### Frontend Already Correct
File: admin-ui/src/app/manufacturers/page.tsx
- Already calling: organizationsApi.list({ type: 'manufacturer' })
- Already filtering correctly

## Test Results

**Before Fix:**
GET /api/v1/organizations?type=manufacturer
Returns: 18 organizations (all types)

**After Fix:**  
GET /api/v1/organizations?type=manufacturer
Returns: 4 manufacturers only
- Philips Healthcare India
- Siemens Healthineers India
- Wipro GE Healthcare
- Global Manufacturer A

## Database Data

**Organizations by Type:**
- Manufacturers: 4
- Distributors: 4
- Dealers: 1
- Hospitals: 5
- Imaging Centers: 3
- Suppliers: 2
**Total: 18**

## API Filtering Now Supports

**By Type:**
- ?type=manufacturer
- ?type=distributor
- ?type=dealer
- ?type=hospital
- ?type=supplier
- ?type=imaging_center

**By Status:**
- ?status=active
- ?status=inactive

**Combined:**
- ?type=manufacturer&status=active

## Result

After restarting frontend, /manufacturers page will show:
✅ Only 4 manufacturers (not all 18 organizations)
✅ Philips Healthcare India
✅ Siemens Healthineers India  
✅ Wipro GE Healthcare
✅ Global Manufacturer A

## Files Modified
1. internal/core/organizations/api/handler.go
2. internal/core/organizations/infra/repository.go
3. admin-ui/src/app/manufacturers/page.tsx (already had correct API call)

## Status
✅ Backend rebuilt and restarted
✅ API filtering tested and working
✅ Frontend restart needed to see changes
