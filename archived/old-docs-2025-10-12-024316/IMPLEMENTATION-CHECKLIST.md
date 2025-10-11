# Implementation Checklist - GenQ Platform API Integration

**Reference:** See CODE-AUDIT-AND-IMPROVEMENTS.md for detailed code and documentation

---

## ✅ Phase 1: Verify Backend (Days 1-2)

### Backend Verification
- [ ] Start backend server: `cd cmd/platform && go run main.go`
- [ ] Test health endpoint: `curl http://localhost:8080/health`
- [ ] Test equipment API: `curl -H "X-Tenant-ID: default" http://localhost:8080/v1/equipment`
- [ ] Test tickets API: `curl -H "X-Tenant-ID: default" http://localhost:8080/v1/tickets`
- [ ] Test suppliers API: `curl -H "X-Tenant-ID: default" http://localhost:8080/v1/suppliers`
- [ ] Test organizations API: `curl -H "X-Tenant-ID: default" http://localhost:8080/v1/organizations`
- [ ] Verify all endpoints return proper responses (even if empty data)

### Environment Setup
- [ ] Create `.env` file in admin-ui with: `NEXT_PUBLIC_API_URL=http://localhost:8080`
- [ ] Verify CORS is enabled in backend for `http://localhost:3000`
- [ ] Check database is populated with test data

---

## ✅ Phase 2: Create API Client Files (Days 3-5)

### Update Base Client
- [ ] Update `admin-ui/src/lib/api/client.ts` with code from audit document
  - Add axios interceptors
  - Add X-Tenant-ID header
  - Add error handling
  - Add auth token support

### Create New API Clients
- [ ] Create `admin-ui/src/lib/api/manufacturers.ts` (Copy from audit doc)
- [ ] Create `admin-ui/src/lib/api/suppliers.ts` (Copy from audit doc)
- [ ] Update `admin-ui/src/lib/api/equipment.ts` (Copy updated version from audit doc)
- [ ] Update `admin-ui/src/lib/api/tickets.ts` (Copy complete version from audit doc)

### Optional API Clients (If needed)
- [ ] Create `admin-ui/src/lib/api/organizations.ts`
- [ ] Create `admin-ui/src/lib/api/rfq.ts`
- [ ] Create `admin-ui/src/lib/api/quotes.ts`
- [ ] Create `admin-ui/src/lib/api/comparisons.ts`
- [ ] Create `admin-ui/src/lib/api/contracts.ts`

---

## ✅ Phase 3: Install Dependencies (Day 5)

### React Query Setup
```bash
cd admin-ui
npm install @tanstack/react-query @tanstack/react-query-devtools
npm install axios  # If not already installed
```

### Create Query Client Provider
- [ ] Create `admin-ui/src/providers/QueryProvider.tsx`
- [ ] Wrap app in QueryProvider in `admin-ui/src/app/layout.tsx`

---

## ✅ Phase 4: Update Dashboard Page (Week 2, Days 1-2)

**File:** `admin-ui/src/app/dashboard/page.tsx`

### Changes Needed:
- [ ] Remove hardcoded stats (lines 20-26)
- [ ] Add React Query hooks to fetch:
  - Total manufacturers count
  - Total suppliers count
  - Total equipment count
  - Total active tickets count
- [ ] Add loading state
- [ ] Add error handling
- [ ] Update stats display with real data

---

## ✅ Phase 5: Update Manufacturers Pages (Week 2, Days 2-3)

### Manufacturers List Page
**File:** `admin-ui/src/app/manufacturers/page.tsx`

- [ ] Remove mock manufacturer data
- [ ] Import `manufacturersApi` from API client
- [ ] Use React Query `useQuery` to fetch manufacturers list
- [ ] Add loading skeleton
- [ ] Add error display
- [ ] Add pagination controls
- [ ] Add search functionality
- [ ] Add filter by status

### Manufacturer Dashboard Page
**File:** `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`

- [ ] Remove mock manufacturers object (lines 22-92)
- [ ] Use `useQuery` to fetch manufacturer by ID
- [ ] Use `useQuery` to fetch manufacturer stats
- [ ] Use `useQuery` to fetch manufacturer equipment
- [ ] Use `useQuery` to fetch manufacturer engineers
- [ ] Use `useQuery` to fetch manufacturer tickets
- [ ] Add loading states for each section
- [ ] Add error handling

---

## ✅ Phase 6: Update Suppliers Pages (Week 2, Days 3-4)

### Suppliers List Page
**File:** `admin-ui/src/app/suppliers/page.tsx`

- [ ] Remove mock supplier data
- [ ] Import `suppliersApi` from API client
- [ ] Use React Query `useQuery` to fetch suppliers list
- [ ] Add loading skeleton
- [ ] Add error display
- [ ] Add pagination
- [ ] Add search
- [ ] Add filter by status and verification status

### Supplier Dashboard Page
**File:** `admin-ui/src/app/suppliers/[id]/dashboard/page.tsx`

- [ ] Remove mock suppliers object (lines 20-148)
- [ ] Use `useQuery` to fetch supplier by ID
- [ ] Add loading state
- [ ] Add error handling
- [ ] Implement verify/reject/suspend/activate actions using mutations

---

## ✅ Phase 7: Update Equipment Page (Week 2, Day 4)

**File:** `admin-ui/src/app/equipment/page.tsx`

- [ ] Remove `generateMockEquipment()` function
- [ ] Import `equipmentApi` from API client
- [ ] Use `useQuery` to fetch equipment list
- [ ] Add loading skeleton
- [ ] Add error display
- [ ] Add pagination
- [ ] Add filters (manufacturer, category, status, AMC, warranty)
- [ ] Add search by serial number
- [ ] Implement QR code generation using mutations
- [ ] Implement CSV import

---

## ✅ Phase 8: Update Engineers Page (Week 2, Day 5)

**File:** `admin-ui/src/app/engineers/page.tsx`

- [ ] Remove mock engineers data
- [ ] Create `engineers.ts` API client if needed
- [ ] Use `useQuery` to fetch engineers from organizations API
- [ ] Add loading state
- [ ] Add error handling
- [ ] Add filters (region, skills)

---

## ✅ Phase 9: Update Onboarding Flow (Week 3, Days 1-2)

### Files to Update:
- `admin-ui/src/app/onboarding/manufacturer/basic-info/page.tsx`
- `admin-ui/src/app/onboarding/manufacturer/location/page.tsx`
- `admin-ui/src/app/onboarding/manufacturer/contact/page.tsx`
- `admin-ui/src/app/onboarding/manufacturer/review/page.tsx`

### Changes:
- [ ] Remove localStorage usage
- [ ] Use React Context or state management for form data
- [ ] On final submit, call `manufacturersApi.create()`
- [ ] Show success/error messages
- [ ] Redirect to dashboard on success

---

## ✅ Phase 10: Add Error Handling & Loading States (Week 3, Days 3-4)

### Global Error Boundary
- [ ] Create `admin-ui/src/components/ErrorBoundary.tsx`
- [ ] Wrap app in error boundary

### Loading Components
- [ ] Create skeleton loaders for:
  - Table rows
  - Cards
  - Stats
  - Forms

### Error Components
- [ ] Create error display components
- [ ] Add retry buttons
- [ ] Add user-friendly error messages

---

## ✅ Phase 11: Testing (Week 3, Day 5)

### Manual Testing
- [ ] Test all pages load correctly
- [ ] Test filtering works on all list pages
- [ ] Test pagination works
- [ ] Test search functionality
- [ ] Test create/update operations
- [ ] Test error scenarios (backend down, 404s, etc.)
- [ ] Test loading states appear correctly

### API Testing
- [ ] Verify all API calls use correct headers (X-Tenant-ID)
- [ ] Verify auth tokens are included when available
- [ ] Test multi-tenancy works (switch tenant IDs)

---

## ✅ Phase 12: Documentation (Week 4)

### Update Documentation
- [ ] Update README.md with:
  - API endpoints documentation
  - Environment variables needed
  - Setup instructions
  - How to run backend and frontend
- [ ] Create API.md with complete API reference
- [ ] Update architecture diagram
- [ ] Document multi-tenancy setup

### Code Documentation
- [ ] Add JSDoc comments to API client functions
- [ ] Add comments to complex React components
- [ ] Document custom hooks

---

## 🚀 Quick Start Commands

### Backend
```bash
cd cmd/platform
go run main.go
# Or build: go build -o platform && ./platform
```

### Frontend
```bash
cd admin-ui
npm install
npm run dev
```

### Verify APIs
```bash
# Health check
curl http://localhost:8080/health

# List equipment
curl -H "X-Tenant-ID: default" http://localhost:8080/v1/equipment

# List tickets
curl -H "X-Tenant-ID: default" http://localhost:8080/v1/tickets
```

---

## 📊 Progress Tracking

- [ ] Phase 1: Backend Verification (2 days)
- [ ] Phase 2: API Client Files (3 days)
- [ ] Phase 3: Dependencies (0.5 days)
- [ ] Phase 4: Dashboard Page (1.5 days)
- [ ] Phase 5: Manufacturers Pages (2 days)
- [ ] Phase 6: Suppliers Pages (2 days)
- [ ] Phase 7: Equipment Page (1 day)
- [ ] Phase 8: Engineers Page (1 day)
- [ ] Phase 9: Onboarding Flow (2 days)
- [ ] Phase 10: Error Handling (2 days)
- [ ] Phase 11: Testing (1 day)
- [ ] Phase 12: Documentation (variable)

**Total Estimated Time:** 18 days (~4 weeks)

---

## 📞 Support

- Backend API Documentation: See CODE-AUDIT-AND-IMPROVEMENTS.md Part 1
- API Client Code: See CODE-AUDIT-AND-IMPROVEMENTS.md Phase 2
- Frontend Updates: TBD (will be added to audit document)

**Last Updated:** October 10, 2025
