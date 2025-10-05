'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Search, Plus, Upload, Download, Building2, CheckCircle } from 'lucide-react';

interface Manufacturer {
  id: string;
  name: string;
  contactPerson: string;
  email: string;
  phone: string;
  website: string;
  address: string;
  equipmentCount: number;
  engineersCount: number;
  status: string;
  createdAt: string;
}

export default function ManufacturersListPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Mock data - in production this would come from API
  const manufacturersData: Manufacturer[] = useMemo(() => {
    const manufacturers = [
      {
        id: 'MFR-001',
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
    ];

    // Check if current manufacturer exists in localStorage
    if (typeof window !== 'undefined') {
      const currentMfr = localStorage.getItem('current_manufacturer');
      if (currentMfr) {
        const mfrData = JSON.parse(currentMfr);
        // Check if already in list
        const exists = manufacturers.find(m => m.id === mfrData.id);
        if (!exists && mfrData.id) {
          manufacturers.unshift({
            id: mfrData.id,
            name: mfrData.name,
            contactPerson: mfrData.contact_person || '',
            email: mfrData.email || '',
            phone: mfrData.phone || '',
            website: mfrData.website || '',
            address: mfrData.address || '',
            equipmentCount: localStorage.getItem('equipment_imported') === 'true' ? 398 : 0,
            engineersCount: localStorage.getItem('engineers') ? JSON.parse(localStorage.getItem('engineers') || '[]').length : 0,
            status: 'Active',
            createdAt: mfrData.created_at || new Date().toISOString().split('T')[0],
          });
        }
      }
    }

    return manufacturers;
  }, []);

  // Filter manufacturers based on search and status
  const filteredManufacturers = useMemo(() => {
    return manufacturersData.filter(manufacturer => {
      const matchesSearch = searchQuery === '' || 
        manufacturer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        manufacturer.contactPerson.toLowerCase().includes(searchQuery.toLowerCase()) ||
        manufacturer.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
        manufacturer.phone.includes(searchQuery) ||
        manufacturer.address.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = filterStatus === 'all' || manufacturer.status.toLowerCase() === filterStatus.toLowerCase();

      return matchesSearch && matchesStatus;
    });
  }, [manufacturersData, searchQuery, filterStatus]);

  const statusCounts = useMemo(() => {
    return manufacturersData.reduce((acc, mfr) => {
      acc[mfr.status] = (acc[mfr.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
  }, [manufacturersData]);

  const totalEquipment = useMemo(() => {
    return manufacturersData.reduce((sum, mfr) => sum + mfr.equipmentCount, 0);
  }, [manufacturersData]);

  const totalEngineers = useMemo(() => {
    return manufacturersData.reduce((sum, mfr) => sum + mfr.engineersCount, 0);
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
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Manufacturers</CardDescription>
              <CardTitle className="text-3xl">{manufacturersData.length}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active</CardDescription>
              <CardTitle className="text-3xl text-green-600">{statusCounts['Active'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Equipment</CardDescription>
              <CardTitle className="text-3xl text-blue-600">{totalEquipment}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Engineers</CardDescription>
              <CardTitle className="text-3xl text-purple-600">{totalEngineers}</CardTitle>
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
                Showing {filteredManufacturers.length} of {manufacturersData.length} manufacturers
              </div>
            ) : null}
          </CardContent>
        </Card>

        {/* Manufacturers List */}
        {manufacturersData.length === 0 ? (
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
