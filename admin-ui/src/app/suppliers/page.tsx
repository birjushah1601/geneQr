'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Search, Plus, Upload, Download, Store, TrendingUp } from 'lucide-react';

interface Supplier {
  id: string;
  name: string;
  contactPerson: string;
  email: string;
  phone: string;
  category: string;
  location: string;
  rating: number;
  totalOrders: number;
  activeContracts: number;
  status: string;
  createdAt: string;
}

export default function SuppliersListPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Mock data - in production this would come from API
  const suppliersData: Supplier[] = useMemo(() => {
    return [
      {
        id: 'SUP-001',
        name: 'MedTech Supplies India',
        contactPerson: 'Rajesh Kumar',
        email: 'rajesh@medtechsupplies.com',
        phone: '+91-9876543220',
        category: 'Medical Consumables',
        location: 'Mumbai, Maharashtra',
        rating: 4.5,
        totalOrders: 145,
        activeContracts: 3,
        status: 'Active',
        createdAt: '2024-01-10',
      },
      {
        id: 'SUP-002',
        name: 'HealthCare Solutions Ltd',
        contactPerson: 'Anita Desai',
        email: 'anita@healthcaresolutions.com',
        phone: '+91-9876543221',
        category: 'Equipment Parts',
        location: 'Bangalore, Karnataka',
        rating: 4.8,
        totalOrders: 230,
        activeContracts: 5,
        status: 'Active',
        createdAt: '2023-11-15',
      },
      {
        id: 'SUP-003',
        name: 'Bio Medical Instruments',
        contactPerson: 'Suresh Patel',
        email: 'suresh@biomedical.com',
        phone: '+91-9876543222',
        category: 'Diagnostic Equipment',
        location: 'Delhi, India',
        rating: 4.2,
        totalOrders: 89,
        activeContracts: 2,
        status: 'Active',
        createdAt: '2024-02-20',
      },
      {
        id: 'SUP-004',
        name: 'Precision Med Parts',
        contactPerson: 'Meera Singh',
        email: 'meera@precisionmedparts.com',
        phone: '+91-9876543223',
        category: 'Spare Parts',
        location: 'Pune, Maharashtra',
        rating: 4.6,
        totalOrders: 178,
        activeContracts: 4,
        status: 'Active',
        createdAt: '2023-12-05',
      },
      {
        id: 'SUP-005',
        name: 'Global Medical Supplies',
        contactPerson: 'Arjun Mehta',
        email: 'arjun@globalmedsupplies.com',
        phone: '+91-9876543224',
        category: 'General Supplies',
        location: 'Hyderabad, Telangana',
        rating: 4.4,
        totalOrders: 156,
        activeContracts: 3,
        status: 'Active',
        createdAt: '2024-03-12',
      },
      {
        id: 'SUP-006',
        name: 'Advanced Healthcare Products',
        contactPerson: 'Kavita Reddy',
        email: 'kavita@advancedhealthcare.com',
        phone: '+91-9876543225',
        category: 'Medical Consumables',
        location: 'Chennai, Tamil Nadu',
        rating: 3.9,
        totalOrders: 67,
        activeContracts: 1,
        status: 'Pending',
        createdAt: '2024-08-01',
      },
      {
        id: 'SUP-007',
        name: 'Quality Med Equipment',
        contactPerson: 'Vikram Sharma',
        email: 'vikram@qualitymedequip.com',
        phone: '+91-9876543226',
        category: 'Equipment Parts',
        location: 'Ahmedabad, Gujarat',
        rating: 3.2,
        totalOrders: 34,
        activeContracts: 0,
        status: 'Inactive',
        createdAt: '2023-09-15',
      },
    ];
  }, []);

  // Filter suppliers based on search and status
  const filteredSuppliers = useMemo(() => {
    return suppliersData.filter(supplier => {
      const matchesSearch = searchQuery === '' || 
        supplier.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        supplier.contactPerson.toLowerCase().includes(searchQuery.toLowerCase()) ||
        supplier.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
        supplier.phone.includes(searchQuery) ||
        supplier.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
        supplier.location.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = filterStatus === 'all' || supplier.status.toLowerCase() === filterStatus.toLowerCase();

      return matchesSearch && matchesStatus;
    });
  }, [suppliersData, searchQuery, filterStatus]);

  const statusCounts = useMemo(() => {
    return suppliersData.reduce((acc, sup) => {
      acc[sup.status] = (acc[sup.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
  }, [suppliersData]);

  const totalOrders = useMemo(() => {
    return suppliersData.reduce((sum, sup) => sum + sup.totalOrders, 0);
  }, [suppliersData]);

  const activeContracts = useMemo(() => {
    return suppliersData.reduce((sum, sup) => sum + sup.activeContracts, 0);
  }, [suppliersData]);

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

  const getRatingStars = (rating: number) => {
    return '⭐'.repeat(Math.floor(rating));
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
              <h1 className="text-3xl font-bold text-gray-900">Suppliers</h1>
              <p className="text-gray-600 mt-1">
                Manage supplier relationships and procurement
              </p>
            </div>
            
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => router.push('/suppliers/import')}>
                <Upload className="mr-2 h-4 w-4" />
                Import CSV
              </Button>
              <Button onClick={() => router.push('/suppliers/add')}>
                <Plus className="mr-2 h-4 w-4" />
                Add Supplier
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Suppliers</CardDescription>
              <CardTitle className="text-3xl">{suppliersData.length}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active Suppliers</CardDescription>
              <CardTitle className="text-3xl text-green-600">{statusCounts['Active'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Orders</CardDescription>
              <CardTitle className="text-3xl text-blue-600">{totalOrders}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active Contracts</CardDescription>
              <CardTitle className="text-3xl text-purple-600">{activeContracts}</CardTitle>
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
                  placeholder="Search suppliers by name, contact person, email, phone, category, or location..."
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
                Showing {filteredSuppliers.length} of {suppliersData.length} suppliers
              </div>
            ) : null}
          </CardContent>
        </Card>

        {/* Suppliers List */}
        {suppliersData.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="text-gray-400 mb-4">
                <Store className="h-12 w-12 mx-auto" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No Suppliers Found</h3>
              <p className="text-gray-600 mb-4">
                Get started by importing suppliers from a CSV file or adding them manually.
              </p>
              <div className="flex justify-center gap-2">
                <Button onClick={() => router.push('/suppliers/import')}>
                  <Upload className="mr-2 h-4 w-4" />
                  Import CSV
                </Button>
                <Button variant="outline" onClick={() => router.push('/suppliers/add')}>
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
                        Supplier
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact Person
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact Info
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Category
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Performance
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
                    {filteredSuppliers.map((supplier) => (
                      <tr key={supplier.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10">
                              <div className="h-10 w-10 rounded-full bg-purple-500 flex items-center justify-center text-white font-semibold">
                                {supplier.name.substring(0, 2).toUpperCase()}
                              </div>
                            </div>
                            <div className="ml-4">
                              <button
                                onClick={() => router.push(`/suppliers/${supplier.id}/dashboard`)}
                                className="text-sm font-medium text-purple-600 hover:text-purple-800 hover:underline text-left"
                              >
                                {supplier.name}
                              </button>
                              <div className="text-xs text-gray-500">
                                {supplier.location}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {supplier.contactPerson}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{supplier.phone}</div>
                          <div className="text-xs text-gray-500">{supplier.email}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className="px-2 py-1 inline-flex text-xs leading-5 font-medium rounded-full bg-blue-100 text-blue-800">
                            {supplier.category}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">
                            {getRatingStars(supplier.rating)} {supplier.rating}
                          </div>
                          <div className="text-xs text-gray-500">
                            {supplier.totalOrders} orders • {supplier.activeContracts} contracts
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(supplier.status)}`}>
                            {supplier.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => alert(`View details for ${supplier.name}`)}
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
