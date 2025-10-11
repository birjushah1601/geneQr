# 🚀 Quick Start Guide - GenQ Platform

**Last Updated:** October 10, 2025

---

## ✅ What's Ready

- ✅ **4 API client files** created/updated (manufacturers, suppliers, equipment, tickets)
- ✅ **React Query** already set up and working
- ✅ **40+ backend endpoints** documented
- ✅ **Multi-tenant support** configured
- ✅ **Error handling** and interceptors in place

---

## 🎯 To Start Using Real APIs (3 Steps)

### Step 1: Start Backend

```bash
# Terminal 1
cd cmd/platform
go run main.go
```

**Expected output:**
```
INFO  Starting medical platform version=dev environment=development
INFO  HTTP server listening port=8080
```

### Step 2: Verify APIs Work

```bash
# Terminal 2
curl -H "X-Tenant-ID: default" http://localhost:8080/v1/equipment
```

**Expected:** JSON response (even if empty array)

### Step 3: Start Frontend

```bash
# Terminal 3
cd admin-ui
npm run dev
```

**Open:** http://localhost:3000

---

## 📂 Files You Need to Know

### API Clients (Ready to Use):
```
admin-ui/src/lib/api/
├── client.ts         ← Base client with interceptors
├── manufacturers.ts  ← NEW! Manufacturers CRUD + stats
├── suppliers.ts      ← NEW! Suppliers CRUD + actions
├── equipment.ts      ← Updated with /v1 paths
└── tickets.ts        ← Updated with /v1 paths
```

### Documentation:
```
├── CODE-AUDIT-AND-IMPROVEMENTS.md  ← Full backend API docs
├── IMPLEMENTATION-CHECKLIST.md     ← Phase-by-phase plan
├── REACT-QUERY-EXAMPLES.md         ← Copy-paste examples
├── IMPLEMENTATION-COMPLETE.md      ← What was done
└── QUICK-START.md (this file)      ← Quick reference
```

---

## 💻 Example: Update a Page

Let's update the **Dashboard** page to use real APIs:

### Before (Mock Data):
```typescript
// dashboard/page.tsx
const stats = {
  manufacturers: 15,
  suppliers: 42,
  equipment: 1247,
  tickets: 38
};
```

### After (Real APIs):
```typescript
'use client';

import { useQuery } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';
import { suppliersApi } from '@/lib/api/suppliers';
import { equipmentApi } from '@/lib/api/equipment';
import { ticketsApi } from '@/lib/api/tickets';

export default function Dashboard() {
  // Fetch all stats in parallel
  const { data: manufacturers, isLoading: loadingM } = useQuery({
    queryKey: ['manufacturers', 'count'],
    queryFn: () => manufacturersApi.list({ limit: 1 }),
  });

  const { data: suppliers, isLoading: loadingS } = useQuery({
    queryKey: ['suppliers', 'count'],
    queryFn: () => suppliersApi.list({ page: 1, page_size: 1 }),
  });

  const { data: equipment, isLoading: loadingE } = useQuery({
    queryKey: ['equipment', 'count'],
    queryFn: () => equipmentApi.list({ page: 1, page_size: 1 }),
  });

  const { data: tickets, isLoading: loadingT } = useQuery({
    queryKey: ['tickets', 'count'],
    queryFn: () => ticketsApi.list({ page: 1, page_size: 1 }),
  });

  const isLoading = loadingM || loadingS || loadingE || loadingT;

  return (
    <div>
      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <div className="grid grid-cols-4 gap-4">
          <StatCard title="Manufacturers" value={manufacturers?.total || 0} />
          <StatCard title="Suppliers" value={suppliers?.total || 0} />
          <StatCard title="Equipment" value={equipment?.total || 0} />
          <StatCard title="Tickets" value={tickets?.total || 0} />
        </div>
      )}
    </div>
  );
}
```

**That's it!** Your dashboard now uses real data. 🎉

---

## 🔥 Common API Calls

### List Manufacturers
```typescript
import { manufacturersApi } from '@/lib/api/manufacturers';

const { data } = useQuery({
  queryKey: ['manufacturers'],
  queryFn: () => manufacturersApi.list({ limit: 20 }),
});
// data.items = Manufacturer[]
// data.total = number
```

### List Suppliers with Filters
```typescript
import { suppliersApi } from '@/lib/api/suppliers';

const { data } = useQuery({
  queryKey: ['suppliers', filters],
  queryFn: () => suppliersApi.list({
    status: ['active'],
    verification_status: ['verified'],
    search: 'medical',
    page: 1,
    page_size: 20,
  }),
});
```

### List Equipment
```typescript
import { equipmentApi } from '@/lib/api/equipment';

const { data } = useQuery({
  queryKey: ['equipment', manufacturerId],
  queryFn: () => equipmentApi.list({
    customer_id: manufacturerId,
    status: 'active',
    page: 1,
    page_size: 20,
  }),
});
```

### List Tickets
```typescript
import { ticketsApi } from '@/lib/api/tickets';

const { data } = useQuery({
  queryKey: ['tickets', filters],
  queryFn: () => ticketsApi.list({
    status: 'open',
    priority: 'high',
    page: 1,
    page_size: 20,
  }),
});
```

### Create/Update with Mutations
```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';

const queryClient = useQueryClient();

const createMutation = useMutation({
  mutationFn: (data) => manufacturersApi.create(data),
  onSuccess: () => {
    queryClient.invalidateQueries(['manufacturers']);
    toast.success('Created successfully!');
  },
});

// Use it:
createMutation.mutate({
  name: 'New Manufacturer',
  email: 'contact@example.com',
  status: 'active',
});
```

---

## 🐛 Troubleshooting

### Problem: "Network Error" in console

**Check:**
1. Backend is running on port 8080
2. Frontend is running on port 3000
3. CORS is enabled in backend

**Fix CORS** (add to `cmd/platform/main.go`):
```go
router.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000"},
    AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"*"},
    AllowCredentials: true,
}))
```

### Problem: "404 Not Found"

**Check endpoint paths:**
- Should be `/v1/equipment` not `/equipment`
- Should be `/v1/tickets` not `/tickets`
- Check `baseURL` in `client.ts` is `http://localhost:8080`

### Problem: No data showing

**Check:**
1. Backend database has data (run seed script if needed)
2. Network tab shows successful 200 responses
3. Response has `items` array (not `equipment` or `tickets`)

---

## 📚 Next Steps

1. ✅ **Backend running?** → Test with curl
2. ✅ **Frontend running?** → Check http://localhost:3000
3. ✅ **APIs working?** → Check Network tab
4. 📖 **Read:** `REACT-QUERY-EXAMPLES.md` for more examples
5. 🎨 **Update:** Start with one page (dashboard recommended)
6. 🔄 **Repeat:** Move to next page

---

## 🎉 You're All Set!

The foundation is ready. Now you can:
- Replace mock data with real APIs
- Add filtering and search
- Implement create/edit forms
- Add real-time updates

**Need help?** Check:
- `CODE-AUDIT-AND-IMPROVEMENTS.md` for API docs
- `REACT-QUERY-EXAMPLES.md` for code samples
- `IMPLEMENTATION-CHECKLIST.md` for the full plan

**Happy coding!** 🚀
