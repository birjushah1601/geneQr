'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Search, Plus, Upload, Download, UserCheck, UserX, Users } from 'lucide-react';

interface Engineer {
  id: string;
  name: string;
  phone: string;
  email: string;
  location: string;
  specializations: string;
  status: string;
  assignedTickets: number;
  completedTickets: number;
  rating: number;
}

export default function EngineersListPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Get manufacturer data
  const manufacturer = useMemo(() => {
    if (typeof window !== 'undefined') {
      const data = localStorage.getItem('current_manufacturer');
      return data ? JSON.parse(data) : null;
    }
    return null;
  }, []);

  // Mock data - in production this would come from API
  const engineersData: Engineer[] = useMemo(() => {
    if (typeof window === 'undefined') return [];
    
    const stored = localStorage.getItem('engineers');
    if (!stored) return [];

    try {
      const engineers = JSON.parse(stored);
      
      // Enhance with additional fields
      return engineers.map((eng: any, i: number) => ({
        id: eng.id || `ENG-${String(i + 1).padStart(6, '0')}`,
        name: eng.name,
        phone: eng.phone,
        email: eng.email,
        location: eng.location || 'Not specified',
        specializations: eng.specializations || 'General',
        status: ['Available', 'Available', 'Busy', 'Available'][i % 4],
        assignedTickets: Math.floor(Math.random() * 10),
        completedTickets: Math.floor(Math.random() * 50) + 10,
        rating: +(4 + Math.random()).toFixed(1),
      }));
    } catch (e) {
      return [];
    }
  }, []);

  // Filter engineers based on search and status
  const filteredEngineers = useMemo(() => {
    return engineersData.filter(engineer => {
      const matchesSearch = searchQuery === '' || 
        engineer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.phone.includes(searchQuery) ||
        engineer.email.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.location.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.specializations.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = filterStatus === 'all' || engineer.status.toLowerCase() === filterStatus.toLowerCase();

      return matchesSearch && matchesStatus;
    });
  }, [engineersData, searchQuery, filterStatus]);

  const statusCounts = useMemo(() => {
    return engineersData.reduce((acc, eng) => {
      acc[eng.status] = (acc[eng.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
  }, [engineersData]);

  const totalTickets = useMemo(() => {
    return engineersData.reduce((sum, eng) => sum + eng.assignedTickets, 0);
  }, [engineersData]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'available':
        return 'bg-green-100 text-green-800';
      case 'busy':
        return 'bg-yellow-100 text-yellow-800';
      case 'offline':
        return 'bg-red-100 text-red-800';
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
              <h1 className="text-3xl font-bold text-gray-900">Service Engineers</h1>
              <p className="text-gray-600 mt-1">
                {manufacturer ? `${manufacturer.name} • ` : ''}Manage field service engineers
              </p>
            </div>
            
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => router.push('/engineers/import')}>
                <Upload className="mr-2 h-4 w-4" />
                Import CSV
              </Button>
              <Button onClick={() => router.push('/engineers/add')}>
                <Plus className="mr-2 h-4 w-4" />
                Add Engineer
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Engineers</CardDescription>
              <CardTitle className="text-3xl">{engineersData.length}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Available</CardDescription>
              <CardTitle className="text-3xl text-green-600">{statusCounts['Available'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Busy</CardDescription>
              <CardTitle className="text-3xl text-yellow-600">{statusCounts['Busy'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active Tickets</CardDescription>
              <CardTitle className="text-3xl text-blue-600">{totalTickets}</CardTitle>
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
                  placeholder="Search engineers by name, phone, email, location, or specialization..."
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
                  <option value="available">Available</option>
                  <option value="busy">Busy</option>
                  <option value="offline">Offline</option>
                </select>
                <Button variant="outline">
                  <Download className="mr-2 h-4 w-4" />
                  Export
                </Button>
              </div>
            </div>
            
            {searchQuery || filterStatus !== 'all' ? (
              <div className="mt-4 text-sm text-gray-600">
                Showing {filteredEngineers.length} of {engineersData.length} engineers
              </div>
            ) : null}
          </CardContent>
        </Card>

        {/* Engineers List */}
        {engineersData.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="text-gray-400 mb-4">
                <Users className="h-12 w-12 mx-auto" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No Engineers Found</h3>
              <p className="text-gray-600 mb-4">
                Get started by importing engineers from a CSV file or adding them manually.
              </p>
              <div className="flex justify-center gap-2">
                <Button onClick={() => router.push('/engineers/import')}>
                  <Upload className="mr-2 h-4 w-4" />
                  Import CSV
                </Button>
                <Button variant="outline" onClick={() => router.push('/engineers/add')}>
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
                        Engineer
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Contact
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Location
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Specializations
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Performance
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredEngineers.map((engineer) => (
                      <tr key={engineer.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <div className="flex-shrink-0 h-10 w-10">
                              <div className="h-10 w-10 rounded-full bg-blue-500 flex items-center justify-center text-white font-semibold">
                                {engineer.name.split(' ').map(n => n[0]).join('').substring(0, 2)}
                              </div>
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-gray-900">
                                {engineer.name}
                              </div>
                              <div className="text-xs text-gray-500">
                                ID: {engineer.id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{engineer.phone}</div>
                          <div className="text-xs text-gray-500">{engineer.email}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {engineer.location}
                        </td>
                        <td className="px-6 py-4">
                          <div className="text-sm text-gray-900 max-w-xs truncate">
                            {engineer.specializations}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(engineer.status)}`}>
                            {engineer.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">
                            {getRatingStars(engineer.rating)} {engineer.rating}
                          </div>
                          <div className="text-xs text-gray-500">
                            {engineer.completedTickets} completed • {engineer.assignedTickets} active
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => alert(`View details for ${engineer.name}`)}
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
