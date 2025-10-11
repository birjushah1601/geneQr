# React Query Implementation Examples

This document provides copy-paste ready examples for updating frontend pages to use React Query with the new API clients.

---

## Setup React Query Provider

### 1. Install Dependencies

```bash
cd admin-ui
npm install @tanstack/react-query @tanstack/react-query-devtools
```

### 2. Create Query Provider

**File:** `admin-ui/src/providers/QueryProvider.tsx`

```typescript
'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { useState } from 'react';

export default function QueryProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000, // 1 minute
            retry: 1,
            refetchOnWindowFocus: false,
          },
        },
      })
  );

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
```

### 3. Wrap App in Provider

**File:** `admin-ui/src/app/layout.tsx`

```typescript
import QueryProvider from '@/providers/QueryProvider';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>
          {children}
        </QueryProvider>
      </body>
    </html>
  );
}
```

---

## Example 1: Dashboard Page with Stats

**File:** `admin-ui/src/app/dashboard/page.tsx`

```typescript
'use client';

import { useQuery } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';
import { suppliersApi } from '@/lib/api/suppliers';
import { equipmentApi } from '@/lib/api/equipment';
import { ticketsApi } from '@/lib/api/tickets';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Building2, Package, Wrench, AlertCircle } from 'lucide-react';

export default function DashboardPage() {
  // Fetch all stats in parallel
  const { data: manufacturers, isLoading: loadingManufacturers } = useQuery({
    queryKey: ['manufacturers', 'count'],
    queryFn: () => manufacturersApi.list({ limit: 1 }),
  });

  const { data: suppliers, isLoading: loadingSuppliers } = useQuery({
    queryKey: ['suppliers', 'count'],
    queryFn: () => suppliersApi.list({ page: 1, page_size: 1 }),
  });

  const { data: equipment, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'count'],
    queryFn: () => equipmentApi.list({ page: 1, page_size: 1 }),
  });

  const { data: tickets, isLoading: loadingTickets } = useQuery({
    queryKey: ['tickets', 'count'],
    queryFn: () => ticketsApi.list({ page: 1, page_size: 1, status: 'open' }),
  });

  const isLoading =
    loadingManufacturers || loadingSuppliers || loadingEquipment || loadingTickets;

  const stats = [
    {
      title: 'Total Manufacturers',
      value: manufacturers?.total || 0,
      icon: Building2,
      loading: loadingManufacturers,
    },
    {
      title: 'Total Suppliers',
      value: suppliers?.total || 0,
      icon: Package,
      loading: loadingSuppliers,
    },
    {
      title: 'Total Equipment',
      value: equipment?.total || 0,
      icon: Wrench,
      loading: loadingEquipment,
    },
    {
      title: 'Active Tickets',
      value: tickets?.total || 0,
      icon: AlertCircle,
      loading: loadingTickets,
    },
  ];

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Platform Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              {stat.loading ? (
                <div className="h-8 w-20 bg-gray-200 animate-pulse rounded" />
              ) : (
                <div className="text-2xl font-bold">{stat.value}</div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
```

---

## Example 2: Manufacturers List Page with Filtering

**File:** `admin-ui/src/app/manufacturers/page.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card } from '@/components/ui/card';
import Link from 'next/link';
import { Building2, Search } from 'lucide-react';

export default function ManufacturersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState('');

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['manufacturers', page, search, status],
    queryFn: () =>
      manufacturersApi.list({
        limit: 20,
        offset: (page - 1) * 20,
        search: search || undefined,
        status: status || undefined,
      }),
    keepPreviousData: true, // Keep old data while loading new page
  });

  if (isError) {
    return (
      <div className="p-8">
        <Card className="p-6 bg-red-50 border-red-200">
          <h2 className="text-red-800 font-semibold mb-2">Error Loading Manufacturers</h2>
          <p className="text-red-600">{error.message}</p>
          <Button 
            onClick={() => window.location.reload()} 
            className="mt-4"
            variant="outline"
          >
            Retry
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Manufacturers</h1>
        <Link href="/manufacturers/new">
          <Button>
            <Building2 className="mr-2 h-4 w-4" />
            Add Manufacturer
          </Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="mb-6 flex gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
          <Input
            placeholder="Search manufacturers..."
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1); // Reset to first page on search
            }}
            className="pl-10"
          />
        </div>
        <select
          value={status}
          onChange={(e) => {
            setStatus(e.target.value);
            setPage(1);
          }}
          className="border rounded-md px-3 py-2"
        >
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="inactive">Inactive</option>
          <option value="pending">Pending</option>
        </select>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <Card key={i} className="p-6">
              <div className="animate-pulse space-y-3">
                <div className="h-4 bg-gray-200 rounded w-1/4" />
                <div className="h-3 bg-gray-200 rounded w-1/2" />
                <div className="h-3 bg-gray-200 rounded w-1/3" />
              </div>
            </Card>
          ))}
        </div>
      ) : (
        <>
          {/* Manufacturers List */}
          <div className="space-y-4">
            {data?.items?.length === 0 ? (
              <Card className="p-12 text-center">
                <Building2 className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  No manufacturers found
                </h3>
                <p className="text-gray-500">
                  {search || status
                    ? 'Try adjusting your filters'
                    : 'Get started by adding a manufacturer'}
                </p>
              </Card>
            ) : (
              data?.items?.map((manufacturer) => (
                <Link key={manufacturer.id} href={`/manufacturers/${manufacturer.id}/dashboard`}>
                  <Card className="p-6 hover:shadow-lg transition-shadow cursor-pointer">
                    <div className="flex items-start justify-between">
                      <div>
                        <h3 className="text-lg font-semibold">{manufacturer.name}</h3>
                        <p className="text-sm text-gray-600 mt-1">
                          {manufacturer.email || 'No email'}
                        </p>
                        <p className="text-sm text-gray-600">
                          {manufacturer.phone || 'No phone'}
                        </p>
                      </div>
                      <div>
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-medium ${
                            manufacturer.status === 'active'
                              ? 'bg-green-100 text-green-800'
                              : manufacturer.status === 'pending'
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {manufacturer.status}
                        </span>
                      </div>
                    </div>
                  </Card>
                </Link>
              ))
            )}
          </div>

          {/* Pagination */}
          {data && data.total > 20 && (
            <div className="mt-6 flex justify-center gap-2">
              <Button
                variant="outline"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
              >
                Previous
              </Button>
              <span className="px-4 py-2">
                Page {page} of {Math.ceil(data.total / 20)}
              </span>
              <Button
                variant="outline"
                onClick={() => setPage((p) => p + 1)}
                disabled={page >= Math.ceil(data.total / 20)}
              >
                Next
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
```

---

## Example 3: Manufacturer Dashboard with Multiple Queries

**File:** `admin-ui/src/app/manufacturers/[id]/dashboard/page.tsx`

```typescript
'use client';

import { useQuery } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Package, Users, AlertCircle } from 'lucide-react';

export default function ManufacturerDashboard({
  params,
}: {
  params: { id: string };
}) {
  const manufacturerId = params.id;

  // Fetch manufacturer details
  const {
    data: manufacturer,
    isLoading: loadingManufacturer,
    isError: errorManufacturer,
  } = useQuery({
    queryKey: ['manufacturer', manufacturerId],
    queryFn: () => manufacturersApi.getById(manufacturerId),
  });

  // Fetch manufacturer stats
  const { data: stats, isLoading: loadingStats } = useQuery({
    queryKey: ['manufacturer', manufacturerId, 'stats'],
    queryFn: () => manufacturersApi.getStats(manufacturerId),
    enabled: !!manufacturerId, // Only run if we have an ID
  });

  // Fetch equipment
  const { data: equipment, isLoading: loadingEquipment } = useQuery({
    queryKey: ['manufacturer', manufacturerId, 'equipment'],
    queryFn: () =>
      manufacturersApi.getEquipment(manufacturerId, {
        page: 1,
        page_size: 5,
      }),
    enabled: !!manufacturerId,
  });

  // Fetch tickets
  const { data: tickets, isLoading: loadingTickets } = useQuery({
    queryKey: ['manufacturer', manufacturerId, 'tickets'],
    queryFn: () =>
      manufacturersApi.getTickets(manufacturerId, {
        page: 1,
        page_size: 5,
      }),
    enabled: !!manufacturerId,
  });

  if (errorManufacturer) {
    return (
      <div className="p-8">
        <Card className="p-6 bg-red-50 border-red-200">
          <h2 className="text-red-800 font-semibold">Manufacturer not found</h2>
        </Card>
      </div>
    );
  }

  if (loadingManufacturer) {
    return (
      <div className="p-8">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4" />
          <div className="grid grid-cols-3 gap-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">{manufacturer?.name}</h1>
        <p className="text-gray-600 mt-2">{manufacturer?.email}</p>
        <p className="text-gray-600">{manufacturer?.phone}</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Total Equipment</CardTitle>
            <Package className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            {loadingStats ? (
              <div className="h-8 w-16 bg-gray-200 animate-pulse rounded" />
            ) : (
              <div className="text-2xl font-bold">{stats?.equipmentCount || 0}</div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Engineers</CardTitle>
            <Users className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            {loadingStats ? (
              <div className="h-8 w-16 bg-gray-200 animate-pulse rounded" />
            ) : (
              <div className="text-2xl font-bold">{stats?.engineersCount || 0}</div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">Active Tickets</CardTitle>
            <AlertCircle className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            {loadingStats ? (
              <div className="h-8 w-16 bg-gray-200 animate-pulse rounded" />
            ) : (
              <div className="text-2xl font-bold">{stats?.activeTickets || 0}</div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent Equipment */}
      <Card className="mb-6">
        <CardHeader>
          <CardTitle>Recent Equipment</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingEquipment ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="h-12 bg-gray-200 animate-pulse rounded" />
              ))}
            </div>
          ) : equipment?.items?.length === 0 ? (
            <p className="text-gray-500 text-center py-4">No equipment found</p>
          ) : (
            <div className="space-y-2">
              {equipment?.items?.map((item) => (
                <div
                  key={item.id}
                  className="flex justify-between items-center p-3 border rounded-lg"
                >
                  <div>
                    <p className="font-medium">{item.name}</p>
                    <p className="text-sm text-gray-600">{item.serial_number}</p>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs rounded ${
                      item.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {item.status}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Recent Tickets */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Tickets</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingTickets ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="h-12 bg-gray-200 animate-pulse rounded" />
              ))}
            </div>
          ) : tickets?.items?.length === 0 ? (
            <p className="text-gray-500 text-center py-4">No tickets found</p>
          ) : (
            <div className="space-y-2">
              {tickets?.items?.map((ticket) => (
                <div
                  key={ticket.id}
                  className="flex justify-between items-center p-3 border rounded-lg"
                >
                  <div>
                    <p className="font-medium">{ticket.ticket_number}</p>
                    <p className="text-sm text-gray-600">{ticket.description}</p>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs rounded ${
                      ticket.status === 'open'
                        ? 'bg-yellow-100 text-yellow-800'
                        : ticket.status === 'closed'
                        ? 'bg-gray-100 text-gray-800'
                        : 'bg-blue-100 text-blue-800'
                    }`}
                  >
                    {ticket.status}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## Example 4: Creating/Updating with Mutations

**File:** `admin-ui/src/app/manufacturers/new/page.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { manufacturersApi } from '@/lib/api/manufacturers';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card } from '@/components/ui/card';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner'; // Or your toast library

export default function NewManufacturerPage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    website: '',
    address: '',
    contact_person: '',
  });

  // Create mutation
  const createMutation = useMutation({
    mutationFn: (data: typeof formData) => manufacturersApi.create(data),
    onSuccess: (newManufacturer) => {
      // Invalidate manufacturers list to refetch
      queryClient.invalidateQueries(['manufacturers']);
      
      // Show success message
      toast.success('Manufacturer created successfully!');
      
      // Redirect to manufacturer dashboard
      router.push(`/manufacturers/${newManufacturer.id}/dashboard`);
    },
    onError: (error: any) => {
      // Show error message
      toast.error(error.response?.data?.message || 'Failed to create manufacturer');
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createMutation.mutate(formData);
  };

  return (
    <div className="p-8 max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Add New Manufacturer</h1>

      <Card className="p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">
              Name <span className="text-red-500">*</span>
            </label>
            <Input
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
              placeholder="Enter manufacturer name"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Email</label>
            <Input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              placeholder="email@example.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Phone</label>
            <Input
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              placeholder="+1 (555) 000-0000"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Website</label>
            <Input
              value={formData.website}
              onChange={(e) => setFormData({ ...formData, website: e.target.value })}
              placeholder="https://example.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Address</label>
            <Input
              value={formData.address}
              onChange={(e) => setFormData({ ...formData, address: e.target.value })}
              placeholder="123 Main St, City, State"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Contact Person</label>
            <Input
              value={formData.contact_person}
              onChange={(e) =>
                setFormData({ ...formData, contact_person: e.target.value })
              }
              placeholder="John Doe"
            />
          </div>

          <div className="flex gap-4 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => router.back()}
              disabled={createMutation.isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createMutation.isLoading}>
              {createMutation.isLoading ? 'Creating...' : 'Create Manufacturer'}
            </Button>
          </div>

          {createMutation.isError && (
            <div className="p-4 bg-red-50 border border-red-200 rounded text-red-800">
              {createMutation.error?.response?.data?.message || 'An error occurred'}
            </div>
          )}
        </form>
      </Card>
    </div>
  );
}
```

---

## Example 5: Equipment Page with CSV Import

**File:** `admin-ui/src/app/equipment/page.tsx` (excerpt)

```typescript
'use client';

import { useState, useRef } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { equipmentApi } from '@/lib/api/equipment';
import { Button } from '@/components/ui/button';
import { Upload } from 'lucide-react';
import { toast } from 'sonner';

export default function EquipmentPage() {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);

  // ... other state and queries

  // CSV Import mutation
  const importMutation = useMutation({
    mutationFn: (file: File) => equipmentApi.importCSV(file),
    onSuccess: (result) => {
      toast.success(
        `Import completed: ${result.success} succeeded, ${result.failed} failed`
      );
      queryClient.invalidateQueries(['equipment']);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Import failed');
    },
  });

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      importMutation.mutate(file);
    }
  };

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Equipment</h1>
        <div className="flex gap-2">
          <input
            ref={fileInputRef}
            type="file"
            accept=".csv"
            onChange={handleFileUpload}
            className="hidden"
          />
          <Button
            variant="outline"
            onClick={() => fileInputRef.current?.click()}
            disabled={importMutation.isLoading}
          >
            <Upload className="mr-2 h-4 w-4" />
            {importMutation.isLoading ? 'Importing...' : 'Import CSV'}
          </Button>
        </div>
      </div>

      {/* Rest of the component */}
    </div>
  );
}
```

---

## Common Patterns

### 1. Loading Skeleton

```typescript
{isLoading ? (
  <div className="space-y-4">
    {[...Array(5)].map((_, i) => (
      <div key={i} className="animate-pulse bg-gray-200 h-20 rounded" />
    ))}
  </div>
) : (
  // Your content
)}
```

### 2. Error Display with Retry

```typescript
{isError && (
  <Card className="p-6 bg-red-50 border-red-200">
    <h2 className="text-red-800 font-semibold mb-2">Error</h2>
    <p className="text-red-600 mb-4">{error.message}</p>
    <Button onClick={() => refetch()} variant="outline">
      Retry
    </Button>
  </Card>
)}
```

### 3. Empty State

```typescript
{data?.items?.length === 0 && (
  <Card className="p-12 text-center">
    <Icon className="mx-auto h-12 w-12 text-gray-400 mb-4" />
    <h3 className="text-lg font-medium mb-2">No data found</h3>
    <p className="text-gray-500">Get started by adding items</p>
  </Card>
)}
```

### 4. Optimistic Updates

```typescript
const updateMutation = useMutation({
  mutationFn: (data) => api.update(id, data),
  onMutate: async (newData) => {
    // Cancel outgoing refetches
    await queryClient.cancelQueries(['item', id]);
    
    // Snapshot previous value
    const previous = queryClient.getQueryData(['item', id]);
    
    // Optimistically update
    queryClient.setQueryData(['item', id], newData);
    
    // Return context with previous value
    return { previous };
  },
  onError: (err, newData, context) => {
    // Rollback on error
    queryClient.setQueryData(['item', id], context.previous);
  },
  onSettled: () => {
    // Refetch after error or success
    queryClient.invalidateQueries(['item', id]);
  },
});
```

---

## Next Steps

1. Copy the QueryProvider setup to your project
2. Pick a page (start with Dashboard)
3. Copy the relevant example code
4. Adjust to your specific needs
5. Test with backend running
6. Move to next page

**Happy coding!** 🚀
