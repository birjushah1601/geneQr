# Admin Dashboard Authentication Fix

## Problem

Admin dashboard showing 0 for equipment, tickets, and engineers despite data existing in database.

## Root Cause

Dashboard was making direct fetch() calls WITHOUT JWT authentication tokens.
Only sent X-Tenant-ID header, which caused 401 Unauthorized from AuthMiddleware.

## Solution

Changed admin-ui/src/app/dashboard/page.tsx to use proper API client utilities:

- equipmentApi.list() instead of fetch()
- ticketsApi.list() instead of fetch()  
- engineersApi.list() instead of fetch()

These utilities automatically include Authorization: Bearer <token> header from localStorage.

## Result

? Admin dashboard now shows correct counts
? JWT authentication working properly
? Multi-tenant filtering functional end-to-end

## Files Modified

- admin-ui/src/app/dashboard/page.tsx
