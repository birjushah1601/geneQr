'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Building2, 
  Package, 
  Users, 
  Ticket, 
  Store,
  ArrowRight,
  Loader2,
  Factory,
  Truck,
  ShoppingBag,
  Hospital,
  Brain,
  Sparkles
} from 'lucide-react';
import { organizationsApi } from '@/lib/api/organizations';
import { equipmentApi } from '@/lib/api/equipment';
import { ticketsApi } from '@/lib/api/tickets';
import { engineersApi } from '@/lib/api/engineers';

export default function AdminDashboard() {
  const router = useRouter();

  // Fetch organizations stats from new unified API
  const { data: organizationsData, isLoading: loadingOrganizations } = useQuery({
    queryKey: ['organizations', 'all'],
    queryFn: () => organizationsApi.list(),
  });

  const { data: equipmentData, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'count'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const response = await fetch(`${apiBaseUrl}/v1/equipment?limit=1000`, {
        headers: { 'X-Tenant-ID': 'default' }
      });
      const data = await response.json();
      return { equipment: data.equipment || [], total: data.equipment?.length || 0 };
    },
  });

  const { data: ticketsData, isLoading: loadingTickets } = useQuery({
    queryKey: ['tickets', 'count', 'active'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const response = await fetch(`${apiBaseUrl}/v1/tickets?limit=1000`, {
        headers: { 'X-Tenant-ID': 'default' }
      });
      if (!response.ok) return { total: 0 };
      const data = await response.json();
      // Filter active tickets (not closed)
      const activeTickets = (data.tickets || []).filter((t: any) => t.status !== 'closed');
      return { total: activeTickets.length, tickets: activeTickets };
    },
    retry: false,
    throwOnError: false,
  });

  const { data: engineersData, isLoading: loadingEngineers } = useQuery({
    queryKey: ['engineers', 'count'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const response = await fetch(`${apiBaseUrl}/v1/engineers?limit=1000`, {
        headers: { 'X-Tenant-ID': 'default' }
      });
      const data = await response.json();
      return { total: data.engineers?.length || 0, engineers: data.engineers || [] };
    },
  });

  // Calculate organization breakdown
  const orgsArray = Array.isArray(organizationsData) ? organizationsData : [];
  const orgsByType = {
    manufacturer: orgsArray.filter((o: any) => o.org_type === 'manufacturer').length,
    distributor: orgsArray.filter((o: any) => o.org_type === 'distributor').length,
    dealer: orgsArray.filter((o: any) => o.org_type === 'dealer').length,
    hospital: orgsArray.filter((o: any) => o.org_type === 'hospital').length,
  };

  // Calculate platform stats from API responses
  const platformStats = {
    totalOrganizations: orgsArray.length,
    manufacturers: orgsByType.manufacturer,
    distributors: orgsByType.distributor,
    dealers: orgsByType.dealer,
    hospitals: orgsByType.hospital,
    totalEquipment: equipmentData?.total || 0,
    totalEngineers: engineersData?.total || 0,
    activeTickets: ticketsData?.total || ticketsData?.tickets?.length || 0,
  };

  const isLoading = loadingOrganizations || loadingEquipment || loadingTickets || loadingEngineers;

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">GenQ Admin Portal</h1>
              <p className="text-gray-600">
                Platform Administration
              </p>
            </div>
            <div className="flex items-center gap-3">
              <div className="text-right">
                <p className="text-sm font-medium">Admin</p>
                <p className="text-xs text-gray-500">admin@genq.com</p>
              </div>
              <div className="w-10 h-10 rounded-full bg-purple-600 text-white flex items-center justify-center font-bold">
                A
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        {/* Welcome Message */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold mb-2">Admin Dashboard</h2>
          <p className="text-gray-600">
            Manage manufacturers, suppliers, and monitor platform activity
          </p>
        </div>

        {/* Platform-Wide Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-6 gap-4 mb-8">
          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/organizations')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <Building2 className="h-5 w-5 text-blue-600" />
                  <p className="text-sm font-medium text-gray-500">Organizations</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2 text-blue-600">{platformStats.totalOrganizations}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">All partners</p>
              </div>
            </CardContent>
          </Card>

          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/organizations?type=manufacturer')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <Factory className="h-5 w-5 text-indigo-600" />
                  <p className="text-sm font-medium text-gray-500">Manufacturers</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2 text-indigo-600">{platformStats.manufacturers}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">OEMs</p>
              </div>
            </CardContent>
          </Card>

          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/organizations?type=distributor')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <Truck className="h-5 w-5 text-purple-600" />
                  <p className="text-sm font-medium text-gray-500">Distributors</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2 text-purple-600">{platformStats.distributors}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">Partners</p>
              </div>
            </CardContent>
          </Card>

          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/organizations?type=dealer')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <ShoppingBag className="h-5 w-5 text-green-600" />
                  <p className="text-sm font-medium text-gray-500">Dealers</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2 text-green-600">{platformStats.dealers}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">Retailers</p>
              </div>
            </CardContent>
          </Card>

          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/organizations?type=hospital')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <Hospital className="h-5 w-5 text-red-600" />
                  <p className="text-sm font-medium text-gray-500">Hospitals</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2 text-red-600">{platformStats.hospitals}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">Clients</p>
              </div>
            </CardContent>
          </Card>

          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/equipment')}>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <div className="flex items-center gap-2 mb-2">
                  <Package className="h-5 w-5 text-blue-600" />
                  <p className="text-sm font-medium text-gray-500">Equipment</p>
                </div>
                {isLoading ? (
                  <div className="flex items-center gap-2 mt-2">
                    <Loader2 className="h-6 w-6 animate-spin text-gray-400" />
                  </div>
                ) : (
                  <p className="text-3xl font-bold mt-2">{platformStats.totalEquipment}</p>
                )}
                <p className="text-xs text-gray-400 mt-1">Platform-wide</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Manufacturers & Suppliers Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {/* Manufacturers Card */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                    <Building2 className="w-5 h-5 text-blue-600" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Manufacturers</CardTitle>
                    <CardDescription>Manage equipment manufacturers</CardDescription>
                  </div>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {platformStats.manufacturers} active manufacturers with {platformStats.totalEquipment} total equipment pieces and {platformStats.totalEngineers} service engineers across the platform.
              </p>
              <Button 
                onClick={() => router.push('/manufacturers')}
                className="w-full"
              >
                View All Manufacturers
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </CardContent>
          </Card>

          {/* Distributors & Dealers Card */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                    <Truck className="w-5 h-5 text-purple-600" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Distribution Network</CardTitle>
                    <CardDescription>Distributors and dealers network</CardDescription>
                  </div>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {platformStats.distributors + platformStats.dealers} partners ({platformStats.distributors} distributors, {platformStats.dealers} dealers) providing equipment distribution and service across India.
              </p>
              <Button 
                onClick={() => router.push('/organizations?type=distributor')}
                className="w-full"
              >
                View Distribution Network
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* AI Systems Section */}
        <div className="mb-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* AI Diagnosis Demo */}
          <Card className="bg-gradient-to-r from-purple-50 to-blue-50 border-purple-200">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-r from-purple-600 to-blue-600 flex items-center justify-center">
                  <Brain className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1">
                  <CardTitle className="text-lg text-purple-900 flex items-center gap-2">
                    🤖 AI-Assisted Diagnosis
                    <Sparkles className="w-4 h-4 text-purple-600" />
                  </CardTitle>
                  <CardDescription className="text-purple-700">
                    Intelligent equipment diagnostics with confidence scoring
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-700 mb-4">
                AI-powered diagnosis system that analyzes equipment issues with confidence scoring and repair recommendations.
              </p>
              <Button 
                onClick={() => router.push('/ai-diagnosis-demo')}
                className="w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700"
              >
                <Brain className="mr-2 h-4 w-4" />
                Try AI Diagnosis Demo
                <Sparkles className="ml-2 h-3 w-3" />
              </Button>
            </CardContent>
          </Card>

          {/* Attachments & AI Analysis */}
          <Card className="bg-gradient-to-r from-green-50 to-teal-50 border-green-200">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-r from-green-600 to-teal-600 flex items-center justify-center">
                  <Hospital className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1">
                  <CardTitle className="text-lg text-green-900 flex items-center gap-2">
                    📸 Attachments & AI Analysis
                    <Brain className="w-4 h-4 text-green-600" />
                  </CardTitle>
                  <CardDescription className="text-green-700">
                    Automated visual analysis of equipment photos and documents
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-700 mb-4">
                Complete attachment management with automated AI analysis, safety concern detection, and repair recommendations.
              </p>
              <Button 
                onClick={() => router.push('/attachments')}
                className="w-full bg-gradient-to-r from-green-600 to-teal-600 hover:from-green-700 hover:to-teal-700"
              >
                <Hospital className="mr-2 h-4 w-4" />
                View Attachments & AI Analysis
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Quick Links Section */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Equipment Quick Link */}
          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/equipment')}>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center">
                  <Package className="w-4 h-4 text-blue-600" />
                </div>
                <CardTitle className="text-base">Equipment</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-3">
                Manage all equipment across the platform
              </p>
              <div className="flex items-center justify-between">
                <span className="text-2xl font-bold text-blue-600">
                  {isLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : platformStats.totalEquipment}
                </span>
                <Button variant="ghost" size="sm">
                  View All <ArrowRight className="ml-1 h-3 w-3" />
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Engineers Quick Link */}
          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/engineers')}>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center">
                  <Users className="w-4 h-4 text-green-600" />
                </div>
                <CardTitle className="text-base">Engineers</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-3">
                Manage service engineers and assignments
              </p>
              <div className="flex items-center justify-between">
                <span className="text-2xl font-bold text-green-600">
                  {isLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : platformStats.totalEngineers}
                </span>
                <Button variant="ghost" size="sm">
                  View All <ArrowRight className="ml-1 h-3 w-3" />
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Tickets Quick Link */}
          <Card className="hover:shadow-md transition-shadow cursor-pointer" onClick={() => router.push('/tickets')}>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-orange-100 flex items-center justify-center">
                  <Ticket className="w-4 h-4 text-orange-600" />
                </div>
                <CardTitle className="text-base">Service Tickets</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-3">
                Track and manage service requests
              </p>
              <div className="flex items-center justify-between">
                <span className="text-2xl font-bold text-orange-600">
                  {isLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : platformStats.activeTickets}
                </span>
                <Button variant="ghost" size="sm">
                  View All <ArrowRight className="ml-1 h-3 w-3" />
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
