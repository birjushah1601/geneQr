# Admin Dashboard - Next Steps & Implementation Plan

**Date:** October 10, 2025  
**Status:** Analysis Complete - Ready for Implementation

---

## 📋 Executive Summary

The GenQ Admin Dashboard has a solid foundation with 10+ pages built, covering onboarding, equipment, engineers, manufacturers, and suppliers. However, critical features are missing:

- **Service Tickets Management UI** (backend exists, no UI)
- **Orders & Procurement UI** (RFQ, Quotes, Contracts - backend exists, no UI)
- **Real Backend Integration** (currently using localStorage/mock data)
- **Detail/Edit Pages** (equipment, engineers, manufacturers, suppliers)

---

## ✅ Current State

### **Built & Working:**
1. ✅ Platform Admin Dashboard (`/dashboard`)
2. ✅ Manufacturer Portal (`/manufacturer/dashboard`, `/manufacturers/[id]/dashboard`)
3. ✅ Supplier Portal (`/suppliers/[id]/dashboard`)
4. ✅ Manufacturers List (`/manufacturers`)
5. ✅ Suppliers List (`/suppliers`)
6. ✅ Equipment List (`/equipment`)
7. ✅ Engineers List (`/engineers`)
8. ✅ Onboarding Flow (3 steps)
9. ✅ CSV Import (equipment, engineers)
10. ✅ Manual Add (engineers)

### **Backend Services Available:**
- ✅ equipment-registry (Go service with API)
- ✅ service-ticket (Go service with API)
- ✅ rfq (Go service with API)
- ✅ quote (Go service with API)
- ✅ comparison (Go service with API)
- ✅ contract (Go service with API)
- ✅ supplier (Go service with API)
- ✅ organizations (Go service with API)

### **API Client Setup:**
- ✅ `admin-ui/src/lib/api/client.ts` - Axios client configured
- ✅ `admin-ui/src/lib/api/equipment.ts` - Equipment API functions
- ✅ `admin-ui/src/lib/api/engineers.ts` - Engineers API functions
- ⚠️ `admin-ui/src/lib/api/tickets.ts` - Exists but incomplete
- ❌ Missing: orders.ts, rfq.ts, quotes.ts, contracts.ts, suppliers.ts

---

## 🎯 Priority 1: Service Tickets Management

**Why First:** Core functionality for the medical equipment platform. Backend fully implemented with dispatcher, SLA monitoring, and webhook support.

### **Pages to Build:**

#### **1.1 Tickets List Page** (`/tickets`)
**File:** `admin-ui/src/app/tickets/page.tsx`

**Features:**
- Stats cards: Total, Open, In Progress, Resolved, Overdue
- Search: by ticket ID, equipment name, serial number, engineer name
- Filters: 
  - Status: open, assigned, in_progress, resolved, closed
  - Priority: critical, high, medium, low
  - SLA Status: within_sla, near_breach, breached
  - Date range
- Table columns:
  - Ticket ID
  - Equipment (name + serial)
  - Priority badge
  - Status badge
  - SLA indicator
  - Assigned Engineer
  - Created date
  - Actions: View, Assign, Update Status
- Pagination
- Mobile-responsive card layout

**API Integration:**
```typescript
// admin-ui/src/lib/api/tickets.ts
export const ticketsApi = {
  list: (filters?: {
    status?: string;
    priority?: string;
    sla_status?: string;
    manufacturer_id?: string;
    search?: string;
    page?: number;
    limit?: number;
  }) => apiClient.get('/v1/tickets', { params: filters }),
  
  getById: (ticketId: string) => apiClient.get(`/v1/tickets/${ticketId}`),
  
  create: (data: {
    equipment_id: string;
    priority: string;
    description: string;
    reported_by?: string;
  }) => apiClient.post('/v1/tickets', data),
  
  assignEngineer: (ticketId: string, engineerId: string) => 
    apiClient.post(`/v1/tickets/${ticketId}/assign`, { engineer_id: engineerId }),
  
  updateStatus: (ticketId: string, status: string, notes?: string) =>
    apiClient.patch(`/v1/tickets/${ticketId}/status`, { status, notes }),
  
  addNote: (ticketId: string, note: string) =>
    apiClient.post(`/v1/tickets/${ticketId}/notes`, { note }),
};
```

---

#### **1.2 Ticket Detail Page** (`/tickets/[id]`)
**File:** `admin-ui/src/app/tickets/[id]/page.tsx`

**Features:**
- **Header Section:**
  - Ticket ID + Status badge
  - Priority badge
  - SLA countdown timer (color-coded: green/yellow/red)
  - Equipment info with QR link
  - Created date/time
  
- **Main Info Card:**
  - Equipment details (name, serial, model, location)
  - Manufacturer info
  - Problem description
  - Reported by (customer/WhatsApp/manual)
  
- **Engineer Assignment Card:**
  - Currently assigned engineer (if any)
  - "Assign Engineer" button → Modal with eligible engineers list
  - Engineer contact info
  - Assignment history
  
- **Status Updates Card:**
  - Current status
  - "Update Status" button → Modal with status dropdown + notes
  - Status history timeline
  
- **Notes & Activity Log:**
  - All notes (engineer, customer, admin)
  - System events (assigned, status changed, SLA breach)
  - Add note functionality
  
- **Actions:**
  - Close ticket
  - Escalate ticket
  - Print/Export
  - Back to list

**State Management:**
```typescript
// Use React Query for real-time updates
const { data: ticket, isLoading } = useQuery({
  queryKey: ['ticket', ticketId],
  queryFn: () => ticketsApi.getById(ticketId),
  refetchInterval: 30000, // Refresh every 30 seconds
});
```

---

#### **1.3 Ticket Creation Flow**

**Option A: From Equipment Page**
- Button on equipment detail page
- Pre-fill equipment info
- Select priority
- Enter description
- Auto-assign or manual assign

**Option B: Manual Create**
- `/tickets/new` page
- Search and select equipment
- Fill details
- Submit

---

### **1.4 SLA Visual Indicators**

**Color Coding:**
- 🟢 **Green:** > 50% time remaining
- 🟡 **Yellow:** 10-50% time remaining
- 🔴 **Red:** < 10% time remaining or breached
- ⚫ **Gray:** No SLA (closed/resolved)

**Display Format:**
```
SLA: 🟢 4h 32m remaining
SLA: 🟡 45m remaining (Near Breach)
SLA: 🔴 BREACHED (2h 15m overdue)
```

---

## 🎯 Priority 2: Orders & Procurement Management

**Why Second:** Enables the marketplace functionality. Backend services (RFQ, Quote, Comparison, Contract) are built.

### **2.1 Orders List Page** (`/orders`)
**File:** `admin-ui/src/app/orders/page.tsx`

**Features:**
- Stats: Total Orders, Pending, Fulfilled, Cancelled
- Search: by order ID, equipment, supplier
- Filters: Status, Date range, Supplier
- Table: Order ID, Equipment, Supplier, Quantity, Amount, Status, Date
- Export orders report

---

### **2.2 RFQ Management** (`/rfq`)
**File:** `admin-ui/src/app/rfq/page.tsx`

**Features:**
- Create new RFQ
- RFQ list with status
- Send to multiple suppliers
- Set deadline
- View responses

**API Integration:**
```typescript
// admin-ui/src/lib/api/rfq.ts
export const rfqApi = {
  list: (filters?: { status?: string; page?: number }) => 
    apiClient.get('/v1/rfqs', { params: filters }),
  
  create: (data: {
    equipment_catalog_id: string;
    quantity: number;
    specifications: any;
    deadline: string;
    supplier_ids: string[];
  }) => apiClient.post('/v1/rfqs', data),
  
  getById: (rfqId: string) => apiClient.get(`/v1/rfqs/${rfqId}`),
  
  close: (rfqId: string) => apiClient.post(`/v1/rfqs/${rfqId}/close`),
};
```

---

### **2.3 Quote Comparison** (`/rfq/[id]/compare`)
**File:** `admin-ui/src/app/rfq/[id]/compare/page.tsx`

**Features:**
- Side-by-side comparison table
- Price comparison
- Delivery time comparison
- Rating/reviews of suppliers
- Select winner button
- Create contract from selected quote

---

### **2.4 Contract Management** (`/contracts`)
**File:** `admin-ui/src/app/contracts/page.tsx`

**Features:**
- Active contracts list
- Expired contracts
- Contract details (terms, payment, delivery)
- Upload contract documents
- Track deliveries

---

## 🎯 Priority 3: Real Backend Integration

**Current Issue:** All data is in localStorage or hardcoded mock data.

### **3.1 Replace Mock Data**

**Files to Update:**
1. `admin-ui/src/app/manufacturers/page.tsx` - Fetch from `/v1/organizations/manufacturers`
2. `admin-ui/src/app/suppliers/page.tsx` - Fetch from `/v1/suppliers`
3. `admin-ui/src/app/equipment/page.tsx` - Fetch from `/v1/equipment`
4. `admin-ui/src/app/engineers/page.tsx` - Fetch from `/v1/engineers`

**Example with React Query:**
```typescript
// Before (mock data)
const equipment = generateMockEquipment();

// After (real API)
import { useQuery } from '@tanstack/react-query';
import { equipmentApi } from '@/lib/api/equipment';

const { data: equipment, isLoading, error } = useQuery({
  queryKey: ['equipment', filters],
  queryFn: () => equipmentApi.list(filters),
});
```

---

### **3.2 API Client Enhancements**

**Add Missing API Files:**

```typescript
// admin-ui/src/lib/api/manufacturers.ts
export const manufacturersApi = {
  list: (filters?: { search?: string; page?: number }) => 
    apiClient.get('/v1/organizations/manufacturers', { params: filters }),
  
  getById: (id: string) => apiClient.get(`/v1/organizations/manufacturers/${id}`),
  
  create: (data: any) => apiClient.post('/v1/organizations/manufacturers', data),
  
  update: (id: string, data: any) => 
    apiClient.put(`/v1/organizations/manufacturers/${id}`, data),
};

// admin-ui/src/lib/api/suppliers.ts
export const suppliersApi = {
  list: (filters?: { category?: string; rating_min?: number }) => 
    apiClient.get('/v1/suppliers', { params: filters }),
  
  getById: (id: string) => apiClient.get(`/v1/suppliers/${id}`),
  
  create: (data: any) => apiClient.post('/v1/suppliers', data),
};

// admin-ui/src/lib/api/orders.ts
export const ordersApi = {
  list: (filters?: { supplier_id?: string; status?: string }) => 
    apiClient.get('/v1/orders', { params: filters }),
  
  getById: (id: string) => apiClient.get(`/v1/orders/${id}`),
};
```

---

## 🎯 Priority 4: Detail & Edit Pages

### **4.1 Equipment Detail Page** (`/equipment/[id]`)
**File:** `admin-ui/src/app/equipment/[id]/page.tsx`

**Features:**
- Equipment photo/image
- QR code display (download button)
- Full specifications
- Installation details
- Service history (all tickets)
- Maintenance schedule
- Related documents
- Edit button
- Generate service request button

---

### **4.2 Engineer Profile Page** (`/engineers/[id]`)
**File:** `admin-ui/src/app/engineers/[id]/page.tsx`

**Features:**
- Profile photo/avatar
- Contact info
- Specializations badges
- Coverage area map
- Performance metrics:
  - Total tickets handled
  - Average resolution time
  - Rating
  - SLA compliance %
- Current assignments
- Ticket history
- Availability status
- Edit button

---

### **4.3 Edit Functionality**

**Pages to Add:**
- `/manufacturers/[id]/edit` - Edit manufacturer details
- `/suppliers/[id]/edit` - Edit supplier details
- `/equipment/[id]/edit` - Edit equipment details
- `/engineers/[id]/edit` - Edit engineer profile

**Pattern:**
```typescript
const { data: manufacturer } = useQuery(['manufacturer', id], () => 
  manufacturersApi.getById(id)
);

const updateMutation = useMutation({
  mutationFn: (data) => manufacturersApi.update(id, data),
  onSuccess: () => {
    queryClient.invalidateQueries(['manufacturer', id]);
    router.push(`/manufacturers/${id}/dashboard`);
  },
});
```

---

## 🎯 Priority 5: Enhanced Features

### **5.1 Dashboard Analytics**

**Add to `/dashboard`:**
- Equipment health trend chart (7-day)
- Ticket volume chart (daily/weekly)
- Engineer utilization chart
- SLA compliance metrics
- Top manufacturers by ticket volume
- Top suppliers by order value

**Libraries to Use:**
- `recharts` or `chart.js` for charts
- `date-fns` for date formatting

---

### **5.2 WhatsApp Monitoring**

**Enhance `/test-qr`:**
- Real-time webhook log viewer
- Show incoming WhatsApp messages
- Display ticket creation from WhatsApp
- Test QR scanning flow
- Monitor webhook delivery status

---

### **5.3 Authentication & Authorization**

**Keycloak Integration:**
- Login page
- Role-based access control:
  - Super Admin (full access)
  - Manufacturer Admin (own data only)
  - Engineer (assigned tickets only)
  - Supplier (own orders only)
- Protected routes with middleware
- JWT token management

---

## 📁 File Structure Plan

```
admin-ui/src/app/
├── dashboard/
│   └── page.tsx ✅ (exists - needs charts)
├── tickets/
│   ├── page.tsx ❌ (NEW - list)
│   ├── new/
│   │   └── page.tsx ❌ (NEW - create)
│   └── [id]/
│       ├── page.tsx ❌ (NEW - detail)
│       └── edit/
│           └── page.tsx ❌ (NEW - edit)
├── orders/
│   ├── page.tsx ❌ (NEW - list)
│   └── [id]/
│       └── page.tsx ❌ (NEW - detail)
├── rfq/
│   ├── page.tsx ❌ (NEW - list)
│   ├── new/
│   │   └── page.tsx ❌ (NEW - create)
│   └── [id]/
│       ├── page.tsx ❌ (NEW - detail)
│       └── compare/
│           └── page.tsx ❌ (NEW - comparison)
├── contracts/
│   ├── page.tsx ❌ (NEW - list)
│   └── [id]/
│       └── page.tsx ❌ (NEW - detail)
├── equipment/
│   ├── page.tsx ✅ (exists)
│   ├── import/
│   │   └── page.tsx ✅ (exists)
│   └── [id]/
│       ├── page.tsx ❌ (NEW - detail)
│       └── edit/
│           └── page.tsx ❌ (NEW - edit)
├── engineers/
│   ├── page.tsx ✅ (exists)
│   ├── add/
│   │   └── page.tsx ✅ (exists)
│   ├── import/
│   │   └── page.tsx ✅ (exists)
│   └── [id]/
│       ├── page.tsx ❌ (NEW - profile)
│       └── edit/
│           └── page.tsx ❌ (NEW - edit)
├── manufacturers/
│   ├── page.tsx ✅ (exists)
│   └── [id]/
│       ├── dashboard/
│       │   └── page.tsx ✅ (exists)
│       └── edit/
│           └── page.tsx ❌ (NEW - edit)
└── suppliers/
    ├── page.tsx ✅ (exists)
    └── [id]/
        ├── dashboard/
        │   └── page.tsx ✅ (exists)
        └── edit/
            └── page.tsx ❌ (NEW - edit)

admin-ui/src/lib/api/
├── client.ts ✅ (exists)
├── equipment.ts ✅ (exists)
├── engineers.ts ✅ (exists)
├── tickets.ts ⚠️ (exists - incomplete)
├── manufacturers.ts ❌ (NEW)
├── suppliers.ts ❌ (NEW)
├── orders.ts ❌ (NEW)
├── rfq.ts ❌ (NEW)
├── quotes.ts ❌ (NEW)
└── contracts.ts ❌ (NEW)
```

---

## 🎨 UI Components Needed

### **New Components to Create:**

```
admin-ui/src/components/ui/
├── badge.tsx ❌ (NEW - for status/priority badges)
├── select.tsx ❌ (NEW - dropdown select)
├── dialog.tsx ❌ (NEW - modals)
├── table.tsx ❌ (NEW - reusable table component)
├── tabs.tsx ❌ (NEW - for detail pages)
├── progress.tsx ❌ (NEW - SLA progress bars)
├── avatar.tsx ❌ (NEW - engineer avatars)
└── chart.tsx ❌ (NEW - dashboard charts wrapper)
```

Install from shadcn/ui:
```bash
cd admin-ui
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add select
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add table
npx shadcn-ui@latest add tabs
npx shadcn-ui@latest add progress
npx shadcn-ui@latest add avatar
```

---

## 📦 Dependencies to Add

```json
{
  "dependencies": {
    "@tanstack/react-query": "^5.0.0",
    "recharts": "^2.10.0",
    "date-fns": "^3.0.0",
    "zod": "^3.22.0",
    "react-hook-form": "^7.49.0"
  },
  "devDependencies": {
    "@tanstack/react-query-devtools": "^5.0.0"
  }
}
```

---

## 🚀 Implementation Phases

### **Phase 1: Service Tickets (Week 1-2)**
1. Build API client functions (tickets.ts)
2. Create tickets list page
3. Create ticket detail page
4. Add assign engineer modal
5. Add update status modal
6. Test with backend APIs

### **Phase 2: Backend Integration (Week 2-3)**
1. Replace localStorage with React Query
2. Connect manufacturers API
3. Connect suppliers API
4. Connect equipment API
5. Connect engineers API
6. Test data flow end-to-end

### **Phase 3: Orders & Procurement (Week 3-4)**
1. Build RFQ list and create pages
2. Build quote comparison page
3. Build contracts list page
4. Build orders list page
5. Test procurement workflow

### **Phase 4: Detail Pages (Week 4-5)**
1. Equipment detail page
2. Engineer profile page
3. Add edit functionality
4. Test CRUD operations

### **Phase 5: Enhancements (Week 5-6)**
1. Dashboard charts
2. WhatsApp monitoring
3. Authentication
4. Mobile responsiveness refinement
5. Performance optimization

---

## ✅ Success Criteria

### **Tickets Management:**
- [ ] Can view all tickets with filters
- [ ] Can create new ticket
- [ ] Can assign engineer to ticket
- [ ] Can update ticket status
- [ ] SLA indicators working
- [ ] Real-time updates via React Query

### **Procurement:**
- [ ] Can create RFQ
- [ ] Can view quotes
- [ ] Can compare quotes side-by-side
- [ ] Can create contract from quote
- [ ] Can track orders

### **Backend Integration:**
- [ ] All pages fetch from real APIs
- [ ] No localStorage/mock data
- [ ] Error handling implemented
- [ ] Loading states working
- [ ] Real-time sync working

### **User Experience:**
- [ ] Mobile responsive
- [ ] Fast load times (< 2s)
- [ ] Smooth interactions
- [ ] Clear error messages
- [ ] Intuitive navigation

---

## 🧪 Testing Strategy

### **Unit Tests:**
- API client functions
- Form validation
- Data transformation

### **Integration Tests:**
- API calls with mock server
- Navigation flows
- State management

### **E2E Tests (Playwright):**
- Complete ticket creation flow
- RFQ to contract flow
- Equipment import flow
- Engineer assignment flow

---

## 📚 Documentation Needed

1. **API Integration Guide** - How to connect new endpoints
2. **Component Library** - Reusable component docs
3. **State Management** - React Query patterns
4. **Deployment Guide** - Build and deploy process
5. **User Manual** - End-user documentation

---

## 🎯 Immediate Next Action

**Start with Priority 1: Service Tickets Management**

**First Task:** Build the Tickets List Page
1. Create `admin-ui/src/app/tickets/page.tsx`
2. Complete `admin-ui/src/lib/api/tickets.ts`
3. Install required dependencies (React Query, shadcn components)
4. Test with backend API at `http://localhost:8080/v1/tickets`

**Command to run:**
```bash
cd admin-ui
npm install @tanstack/react-query recharts date-fns
npx shadcn-ui@latest add badge select dialog
```

---

**Ready to implement?** Let me know which priority you'd like to tackle first, and I'll help build it! 🚀
