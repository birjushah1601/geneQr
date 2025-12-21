'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Search, Plus, Upload, Download, Building2, CheckCircle, Loader2, AlertCircle } from 'lucide-react';
import { organizationsApi } from '@/lib/api/organizations';

interface Manufacturer {
  id: string;
  name: string;
  contactPerson?: string;
  email?: string;
  phone?: string;
  website?: string;
  address?: string;
  equipmentCount?: number;
  engineersCount?: number;
  status: string;
  createdAt?: string;
  org_type: string;
}

export default function ManufacturersListPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Fetch manufacturers from API with equipment counts
  const { data: organizationsData, isLoading, error } = useQuery({
    queryKey: ['organizations', 'manufacturer', 'with-counts'],
    queryFn: async () => {
      // Use the same base URL as apiClient (includes /api)
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const response = await fetch(`${apiBaseUrl}/v1/organizations?type=manufacturer&include_counts=true&limit=1000`, {
        headers: { 'X-Tenant-ID': 'default' }
      });
      if (!response.ok) throw new Error('Failed to fetch manufacturers');
      const data = await response.json();
      return data.items || [];
    },
  });

  // Transform API data to match expected format
  const manufacturersData: Manufacturer[] = useMemo(() => {
    if (!organizationsData || !Array.isArray(organizationsData)) return [];
    
    return organizationsData.map((org: any) => ({
      id: org.id,
      name: org.name,
      org_type: org.org_type,
      status: org.status === 'active' ? 'Active' : 'Inactive',
      contactPerson: org.metadata?.contact_person || 'N/A',
      email: org.metadata?.email || 'N/A',
      phone: org.metadata?.phone || 'N/A',
      website: org.metadata?.website || 'N/A',
      address: org.metadata?.address?.city || org.metadata?.city || 'N/A',
      equipmentCount: org.equipment_count || 0,
      engineersCount: org.engineers_count || 0,
      activeTickets: org.active_tickets || 0,
      createdAt: org.created_at || new Date().toISOString(),
    }));
  }, [organizationsData]);

  // Keep the mock data as fallback only if API fails
  const fallbackManufacturers: Manufacturer[] = useMemo(() => [
      {
        id: 'MFR-001-OLD',
        name: 'Siemens Healthineers',
        contactPerson: 'John Smith',
        email: 'john.smith@siemens.com',
        phone: '+91-9876543210',
        website: 'https://www.siemens-healthineers.com',
        address: 'Mumbai, Maharashtra, India',
        equipmentCount: 150,
        engineersCount: 25,
        status: 'Active',
        createdAt: '2024-01-15',
      },
      {
        id: 'MFR-002',
        name: 'GE Healthcare',
        contactPerson: 'Sarah Johnson',
        email: 'sarah.j@gehealthcare.com',
        phone: '+91-9876543211',
        website: 'https://www.gehealthcare.com',
        address: 'Bangalore, Karnataka, India',
        equipmentCount: 120,
        engineersCount: 20,
        status: 'Active',
        createdAt: '2024-02-10',
      },
      {
        id: 'MFR-003',
        name: 'Philips Healthcare',
        contactPerson: 'Michael Chen',
        email: 'michael.chen@philips.com',
        phone: '+91-9876543212',
        website: 'https://www.philips.com/healthcare',
        address: 'Delhi, India',
        equipmentCount: 95,
        engineersCount: 18,
        status: 'Active',
        createdAt: '2024-03-05',
      },
      {
        id: 'MFR-004',
        name: 'Medtronic India',
        contactPerson: 'Priya Sharma',
        email: 'priya.sharma@medtronic.com',
        phone: '+91-9876543213',
        website: 'https://www.medtronic.com',
        address: 'Hyderabad, Telangana, India',
        equipmentCount: 80,
        engineersCount: 15,
        status: 'Active',
        createdAt: '2024-04-20',
      },
      {
        id: 'MFR-005',
        name: 'Carestream Health',
        contactPerson: 'David Lee',
        email: 'david.lee@carestream.com',
        phone: '+91-9876543214',
        website: 'https://www.carestream.com',
        address: 'Chennai, Tamil Nadu, India',
        equipmentCount: 60,
        engineersCount: 12,
        status: 'Inactive',
        createdAt: '2023-12-01',
      },
    ], []);

  // Use real data if available, fallback to mock data if API fails
  const displayManufacturers = useMemo(() => {
    return manufacturersData.length > 0 ? manufacturersData : (error ? fallbackManufacturers : []);
  }, [manufacturersData, fallbackManufacturers, error]);

  // Filter manufacturers based on search and status
  const filteredManufacturers = useMemo(() => {
    return displayManufacturers.filter(manufacturer => {
      const matchesSearch = searchQuery === '' || 
        manufacturer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (manufacturer.contactPerson && manufacturer.contactPerson.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (manufacturer.email && manufacturer.email.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (manufacturer.phone && manufacturer.phone.includes(searchQuery)) ||
        (manufacturer.address && manufacturer.address.toLowerCase().includes(searchQuery.toLowerCase()));

      const matchesStatus = filterStatus === 'all' || manufacturer.status.toLowerCase() === filterStatus.toLowerCase();

      return matchesSearch && matchesStatus;
    });
  }, [displayManufacturers, searchQuery, filterStatus]);

  const statusCounts = useMemo(() => {
    return displayManufacturers.reduce((acc, mfr) => {
      acc[mfr.status] = (acc[mfr.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
  }, [displayManufacturers]);

  const totalEquipment = useMemo(() => {
    return manufacturersData.reduce((sum, mfr) => sum + (mfr.equipmentCount || 0), 0);
  }, [manufacturersData]);

  const totalEngineers = useMemo(() => {
    return manufacturersData.reduce((sum, mfr) => sum + (mfr.engineersCount || 0), 0);
  }, [manufacturersData]);
  
  const totalActiveTickets = useMemo(() => {
    return manufacturersData.reduce((sum, mfr) => sum + (mfr.activeTickets || 0), 0);
  }, [manufacturersData]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'inactive':
        return 'bg-red-100 text-red-800';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-6">
          <Button
            variant="ghost"
            onClick={() => router.push('/dashboard')}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Dashboard
          </Button>
          
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Manufacturers</h1>
              <p className="text-gray-600 mt-1">
                Manage manufacturer partners and their equipment
              </p>
            </div>
            
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => router.push('/manufacturers/import')}>
                <Upload className="mr-2 h-4 w-4" />
                Import CSV
              </Button>
              <Button onClick={() => router.push('/manufacturers/add')}>
                <Plus className="mr-2 h-4 w-4" />
                Add Manufacturer
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Manufacturers</CardDescription>
              <CardTitle className="text-3xl">
                {isLoading ? <Loader2 className="h-8 w-8 animate-spin" /> : manufacturersData.length}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active</CardDescription>
              <CardTitle className="text-3xl text-green-600">
                {isLoading ? <Loader2 className="h-8 w-8 animate-spin" /> : (statusCounts['Active'] || 0)}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Equipment</CardDescription>
              <CardTitle className="text-3xl text-blue-600">
                {isLoading ? <Loader2 className="h-8 w-8 animate-spin" /> : totalEquipment}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Engineers</CardDescription>
              <CardTitle className="text-3xl text-purple-600">
                {isLoading ? <Loader2 className="h-8 w-8 animate-spin" /> : totalEngineers}
              </CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active Tickets</CardDescription>
              <CardTitle className="text-3xl text-orange-600">
                {isLoading ? <Loader2 className="h-8 w-8 animate-spin" /> : totalActiveTickets}
              </CardTitle>
            </CardHeader>
          </Card>
        </div>

        {/* Search and Filters */}
        <Card className="mb-6">
          <CardContent className="pt-6">
            <div className="flex flex-col md:flex-row gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search manufacturers by name, contact person, email, phone, or address..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              <div className="flex gap-2">
                <select
                  value={filterStatus}
                  onChange={(e) => setFilterStatus(e.target.value)}
                  className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">All Status</option>
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                  <option value="pending">Pending</option>
                </select>
                <Button variant="outline">
                  <Download className="mr-2 h-4 w-4" />
                  Export
                </Button>
              </div>
            </div>
            
            {searchQuery || filterStatus !== 'all' ? (
              <div className="mt-4 text-sm text-gray-600">
                Showing {filteredManufacturers.length} of {displayManufacturers.length} manufacturers
              </div>
            ) : null}
            
            {error && (
              <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg flex items-center gap-2">
                <AlertCircle className="h-5 w-5 text-yellow-600" />
                <p className="text-sm text-yellow-800">Using fallback data. API error: {error.message}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Manufacturers List */}
        {isLoading ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Loader2 className="h-12 w-12 animate-spin mx-auto text-gray-400 mb-4" />
              <p className="text-gray-600">Loading manufacturers...</p>
            </CardContent>
          </Card>
        ) : displayManufacturers.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="text-gray-400 mb-4">
                <Building2 className="h-12 w-12 mx-auto" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No Manufacturers Found</h3>
              <p className="text-gray-600 mb-4">
                Get started by importing manufacturers from a CSV file or adding them manually.
              </p>
              <div className="flex justify-center gap-2">
                <Button onClick={() => router.push('/manufacturers/import')}>
                  <Upload className="mr-2 h-4 w-4" />
                  Import CSV
                </Button>
                <Button variant="outline" onClick={() => router.push('/manufacturers/add')}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Manually
                </Button>
              </div>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <CardContent className="p-0">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Manufacturer
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact Person
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact Info
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Location
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Resources
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredManufacturers.map((manufacturer) => (
                      <tr key={manufacturer.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10">
                              <div className="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                                {manufacturer.name.substring(0, 2).toUpperCase()}
                              </div>
                            </div>
                            <div className="ml-4">
                              <button
                                onClick={() => router.push(`/manufacturers/${manufacturer.id}/dashboard`)}
                                className="text-sm font-medium text-blue-600 hover:text-blue-800 hover:underline text-left"
                              >
                                {manufacturer.name}
                              </button>
                              <div className="text-xs text-gray-500">
                                ID: {manufacturer.id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {manufacturer.contactPerson}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{manufacturer.phone}</div>
                          <div className="text-xs text-gray-500">{manufacturer.email}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {manufacturer.address}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">
                            {manufacturer.equipmentCount} Equipment
                          </div>
                          <div className="text-xs text-gray-500">
                            {manufacturer.engineersCount} Engineers
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(manufacturer.status)}`}>
                            {manufacturer.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => alert(`View details for ${manufacturer.name}`)}
                          >
                            View
                          </Button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
