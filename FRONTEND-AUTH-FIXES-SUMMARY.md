# Frontend Authentication Fixes - Summary

## Overview

Multiple frontend pages were making direct fetch() API calls without JWT authentication,
causing 401 Unauthorized errors after implementing AuthMiddleware on the backend.

## Root Cause

Pages were using direct fetch() calls like:
\\\	ypescript
const response = await fetch(\\/v1/endpoint\, {
  headers: { 'X-Tenant-ID': 'default' }
});
\\\

This only sent X-Tenant-ID header but NO Authorization header with JWT token.

## Solution

Replaced all direct fetch() calls with proper API client utilities that automatically
include JWT authentication from localStorage:

\\\	ypescript
const response = await apiClient.get('/v1/endpoint');
// or
const response = await specificApi.list();
\\\

## Pages Fixed

### 1. Admin Dashboard (/dashboard)
**Problem:** Showing 0 for equipment, tickets, engineers
**Fixed:** 
- Equipment: Now uses equipmentApi.list()
- Tickets: Now uses ticketsApi.list() + fixed field name (tickets vs items)
- Engineers: Now uses engineersApi.list() + fixed field name (engineers vs items)

**Result:** Dashboard now shows:
- Equipment: 73
- Engineers: 34
- Service Tickets: 12

### 2. Manufacturers Page (/manufacturers)
**Problem:** Showing 0 manufacturers
**Fixed:** Now uses organizationsApi.list({ org_type: 'manufacturer' })

**Result:** Page now shows all 12 manufacturers

## API Response Field Name Inconsistencies

Different APIs return data in different field names:

| API | Field Name | Example |
|-----|------------|---------|
| Equipment | \items\ | \{ items: [], total: 73 }\ |
| Engineers | \engineers\ | \{ engineers: [], total: 34 }\ |
| Tickets | \	ickets\ | \{ tickets: [], total: 12 }\ |
| Organizations | \items\ | \{ items: [], total: 32 }\ |

**Solution:** Dashboard now checks for both field names:
\\\	ypescript
const list = response.engineers || response.items || [];
const list = response.tickets || response.items || [];
\\\

## How API Client Works

The apiClient (admin-ui/src/lib/api/client.ts) automatically:

1. ? Reads JWT token from localStorage
2. ? Adds Authorization: Bearer <token> to ALL requests
3. ? Adds X-Tenant-ID header
4. ? Handles 401 errors and token refresh
5. ? Redirects to /login if token refresh fails

## Commits Made

1. \ix: Admin dashboard now uses API client with JWT authentication\
2. \ix: Admin dashboard now correctly reads engineers and tickets API responses\
3. \ix: Manufacturers page now uses API client with JWT authentication\

## Testing Checklist

- [x] Admin dashboard shows correct counts
- [x] Manufacturers page shows all manufacturers
- [x] Equipment page works (already using API client)
- [ ] Engineers page (if exists)
- [ ] Tickets page (if exists)
- [ ] Other organization type pages

## Future Prevention

**Rule:** NEVER use direct fetch() for API calls in the frontend.

**Always use:**
- equipmentApi from '@/lib/api/equipment'
- organizationsApi from '@/lib/api/organizations'
- ticketsApi from '@/lib/api/tickets'
- engineersApi from '@/lib/api/engineers'
- Or apiClient.get/post/patch/delete for new endpoints

These utilities handle authentication automatically.
