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

export default function AdminDashboard() {
  const router = useRouter();

  // Fetch organizations stats from new unified API
  const { data: organizationsData, isLoading: loadingOrganizations } = useQuery({
    queryKey: ['organizations', 'all'],
    queryFn: () => organizationsApi.list(),
  });

  const { data: equipmentData, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'count'],
    queryFn: () => equipmentApi.list({ page: 1, page_size: 1 }),
  });

  const { data: ticketsData, isLoading: loadingTickets } = useQuery({
    queryKey: ['tickets', 'count', 'active'],
    queryFn: () => ticketsApi.list({ page: 1, page_size: 1 }),
  });

  // Calculate organization breakdown
  const orgsData: any = organizationsData;
  const orgsByType = {
    manufacturer: orgsData?.items?.filter((o: any) => o.org_type === 'manufacturer').length || 0,
    distributor: orgsData?.items?.filter((o: any) => o.org_type === 'distributor').length || 0,
    dealer: orgsData?.items?.filter((o: any) => o.org_type === 'dealer').length || 0,
    hospital: orgsData?.items?.filter((o: any) => o.org_type === 'hospital').length || 0,
  };

  // Calculate platform stats from API responses
  const platformStats = {
    totalOrganizations: orgsData?.total || 0,
    manufacturers: orgsByType.manufacturer,
    distributors: orgsByType.distributor,
    dealers: orgsByType.dealer,
    hospitals: orgsByType.hospital,
    totalEquipment: equipmentData?.total || 0,
    totalEngineers: 90, // TODO: Add engineers API endpoint
    activeTickets: ticketsData?.total || 0,
  };

  const isLoading = loadingOrganizations || loadingEquipment || loadingTickets;

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
              <div className="space-y-2 mb-4">
                {/* Top Manufacturers Preview */}
                <div className="text-sm">
                  <p className="font-medium text-gray-700 mb-2">Top Manufacturers:</p>
                  <ul className="space-y-1 text-gray-600">
                    <li>• Siemens Healthineers - 150 equipment</li>
                    <li>• GE Healthcare - 120 equipment</li>
                    <li>• Philips Healthcare - 95 equipment</li>
                  </ul>
                </div>
              </div>
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
              <div className="space-y-2 mb-4">
                {/* Top Partners Preview */}
                <div className="text-sm">
                  <p className="font-medium text-gray-700 mb-2">Top Partners:</p>
                  <ul className="space-y-1 text-gray-600">
                    <li>• MedEquip Distributors - 15 locations</li>
                    <li>• Healthcare Solutions - 12 locations</li>
                    <li>• MediCare Dealers - 8 locations</li>
                  </ul>
                </div>
              </div>
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
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-purple-800">87%</div>
                  <div className="text-sm text-purple-600">Avg Confidence</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-800">4</div>
                  <div className="text-sm text-blue-600">Analysis Factors</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-indigo-800">2-4hr</div>
                  <div className="text-sm text-indigo-600">Est Resolution</div>
                </div>
              </div>
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
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div className="text-center">
                  <div className="text-2xl font-bold text-green-800">4</div>
                  <div className="text-sm text-green-600">Processing</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-teal-800">2</div>
                  <div className="text-sm text-teal-600">With Issues</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-emerald-800">3</div>
                  <div className="text-sm text-emerald-600">Queue Workers</div>
                </div>
              </div>
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

        {/* Platform Activity Section */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Equipment Overview */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center">
                  <Package className="w-4 h-4 text-blue-600" />
                </div>
                <CardTitle className="text-base">Equipment Overview</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Active</span>
                  <span className="text-sm font-semibold text-green-600">405 (80%)</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Maintenance</span>
                  <span className="text-sm font-semibold text-yellow-600">79 (16%)</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Inactive</span>
                  <span className="text-sm font-semibold text-red-600">21 (4%)</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Engineers Overview */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center">
                  <Users className="w-4 h-4 text-green-600" />
                </div>
                <CardTitle className="text-base">Engineers Overview</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Available</span>
                  <span className="text-sm font-semibold text-green-600">68 (76%)</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Busy</span>
                  <span className="text-sm font-semibold text-yellow-600">20 (22%)</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Offline</span>
                  <span className="text-sm font-semibold text-red-600">2 (2%)</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Tickets Overview */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-orange-100 flex items-center justify-center">
                  <Ticket className="w-4 h-4 text-orange-600" />
                </div>
                <CardTitle className="text-base">Tickets Overview</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Open</span>
                  <span className="text-sm font-semibold text-orange-600">23</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">In Progress</span>
                  <span className="text-sm font-semibold text-blue-600">15</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Resolved Today</span>
                  <span className="text-sm font-semibold text-green-600">12</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
