'use client';

import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Building2, 
  Package, 
  Users, 
  Ticket,
  ArrowLeft,
  Upload,
  Plus,
  Loader2
} from 'lucide-react';
import { organizationsApi } from '@/lib/api/organizations';

export default function ManufacturerDashboard() {
  const params = useParams();
  const router = useRouter();
  const manufacturerId = params.id as string;

  // Fetch manufacturer data from API
  const { data: manufacturer, isLoading, error } = useQuery({
    queryKey: ['manufacturer', manufacturerId],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      
      // Fetch organization details with all counts in one API call
      let org;
      let equipmentCount = 0;
      let engineersCount = 0;
      let activeTickets = 0;
      
      try {
        const orgResponse = await fetch(
          `${apiBaseUrl}/v1/organizations/${manufacturerId}?include_counts=true`,
          { headers: { 'X-Tenant-ID': 'default' } }
        );
        if (orgResponse.ok) {
          org = await orgResponse.json();
          equipmentCount = org.equipment_count || 0;
          engineersCount = org.engineers_count || 0;
          activeTickets = org.active_tickets || 0;
        } else {
          // Fallback to basic org API without counts
          org = await organizationsApi.get(manufacturerId);
        }
      } catch (e) {
        console.error('Failed to fetch organization with counts:', e);
        org = await organizationsApi.get(manufacturerId);
      }
      
      return {
        id: org.id,
        name: org.name,
        contactPerson: org.metadata?.contact_person || 'N/A',
        email: org.metadata?.email || 'N/A',
        phone: org.metadata?.phone || 'N/A',
        website: org.metadata?.website || 'N/A',
        address: org.metadata?.address?.city || 'N/A',
        equipmentCount,
        engineersCount,
        activeTickets,
        createdAt: new Date().toISOString().split('T')[0],
        metadata: org.metadata,
      };
    },
    enabled: !!manufacturerId,
  });

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-12 h-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading manufacturer dashboard...</p>
        </div>
      </div>
    );
  }

  // Error state
  if (error || !manufacturer) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Card className="max-w-md">
          <CardContent className="pt-6">
            <div className="text-center">
              <Building2 className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h2 className="text-xl font-bold mb-2">Manufacturer Not Found</h2>
              <p className="text-gray-600 mb-4">
                The manufacturer with ID "{manufacturerId}" does not exist.
              </p>
              <Button onClick={() => router.push('/manufacturers')}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Manufacturers
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <Button
            variant="ghost"
            onClick={() => router.push('/manufacturers')}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Manufacturers
          </Button>
          
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full bg-blue-500 flex items-center justify-center text-white text-2xl font-bold">
                {manufacturer.name.substring(0, 2).toUpperCase()}
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{manufacturer.name}</h1>
                <p className="text-gray-600">{manufacturer.address}</p>
              </div>
            </div>
            
            <div className="text-right">
              <p className="text-sm font-medium">{manufacturer.contactPerson}</p>
              <p className="text-xs text-gray-500">{manufacturer.email}</p>
              <p className="text-xs text-gray-500">{manufacturer.phone}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        {/* Welcome Message */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold mb-2">Manufacturer Dashboard</h2>
          <p className="text-gray-600">
            Manage equipment, engineers, and service operations for {manufacturer.name}
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Equipment</p>
                <p className="text-3xl font-bold mt-2">{manufacturer.equipmentCount}</p>
                <p className="text-xs text-gray-400 mt-1">Registered devices</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Engineers</p>
                <p className="text-3xl font-bold mt-2">{manufacturer.engineersCount}</p>
                <p className="text-xs text-gray-400 mt-1">Service team</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Active Tickets</p>
                <p className="text-3xl font-bold mt-2 text-orange-600">{manufacturer.activeTickets}</p>
                <p className="text-xs text-gray-400 mt-1">Open requests</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Member Since</p>
                <p className="text-lg font-bold mt-2">{new Date(manufacturer.createdAt).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}</p>
                <p className="text-xs text-gray-400 mt-1">Partner status</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Management Sections */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {/* Equipment Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                  <Package className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Equipment Registry</CardTitle>
                  <CardDescription>Manage equipment installations</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {manufacturer.equipmentCount} equipment items registered. View, import, or manage equipment for this manufacturer.
              </p>
              <div className="flex gap-3">
                <Button 
                  onClick={() => router.push(`/equipment?manufacturer=${manufacturerId}`)}
                  className="flex-1"
                >
                  View All Equipment
                </Button>
                <Button 
                  variant="outline"
                  onClick={() => router.push('/equipment/import')}
                >
                  <Upload className="w-4 h-4 mr-2" />
                  Import
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Engineers Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
                  <Users className="w-5 h-5 text-green-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Service Engineers</CardTitle>
                  <CardDescription>Manage service team</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {manufacturer.engineersCount} engineers in the service team. View, add, or manage engineer assignments.
              </p>
              <div className="flex gap-3">
                <Button 
                  onClick={() => router.push(`/engineers?manufacturer=${manufacturerId}`)}
                  className="flex-1"
                >
                  View All Engineers
                </Button>
                <Button 
                  variant="outline"
                  onClick={() => router.push('/engineers/add')}
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Add
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Service Tickets */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center">
                <Ticket className="w-5 h-5 text-orange-600" />
              </div>
              <div>
                <CardTitle className="text-lg">Service Tickets</CardTitle>
                <CardDescription>Active service requests and maintenance</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {manufacturer.activeTickets > 0 
                ? `${manufacturer.activeTickets} active service tickets requiring attention.`
                : 'No active service tickets at the moment. All equipment is running smoothly!'}
            </p>
            <Button 
              onClick={() => router.push(`/tickets?manufacturer=${manufacturerId}`)}
              disabled={manufacturer.activeTickets === 0}
            >
              View All Tickets
            </Button>
          </CardContent>
        </Card>

        {/* Company Information */}
        <Card className="mt-6">
          <CardHeader>
            <CardTitle>Company Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-gray-500">Manufacturer ID</p>
                <p className="text-sm text-gray-900 mt-1">{manufacturer.id}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Website</p>
                <p className="text-sm text-blue-600 mt-1">
                  <a href={`https://${manufacturer.website}`} target="_blank" rel="noopener noreferrer">
                    {manufacturer.website}
                  </a>
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Contact Person</p>
                <p className="text-sm text-gray-900 mt-1">{manufacturer.contactPerson}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Email</p>
                <p className="text-sm text-gray-900 mt-1">{manufacturer.email}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Phone</p>
                <p className="text-sm text-gray-900 mt-1">{manufacturer.phone}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Location</p>
                <p className="text-sm text-gray-900 mt-1">{manufacturer.address}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
