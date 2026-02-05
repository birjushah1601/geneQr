'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Package, 
  Wrench, 
  TrendingUp,
  ArrowRight,
  Loader2,
  Factory,
  AlertCircle,
  Sparkles,
  Upload
} from 'lucide-react';

export default function ManufacturerDashboard() {
  const router = useRouter();

  // Fetch manufacturer-specific data
  const { data: equipmentData, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'manufacturer'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/equipment?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      const data = await response.json();
      return { equipment: data.equipment || [], total: data.equipment?.length || 0 };
    },
  });

  const { data: ticketsData, isLoading: loadingTickets } = useQuery({
    queryKey: ['tickets', 'manufacturer'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/tickets?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      if (!response.ok) return { total: 0, pending: 0, inProgress: 0 };
      const data = await response.json();
      const tickets = data.tickets || [];
      return {
        total: tickets.length,
        pending: tickets.filter((t: any) => t.status === 'pending').length,
        inProgress: tickets.filter((t: any) => t.status === 'in_progress').length,
        resolved: tickets.filter((t: any) => t.status === 'resolved').length,
      };
    },
  });

  const isLoading = loadingEquipment || loadingTickets;

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold mb-2">Manufacturer Dashboard</h1>
        <p className="text-gray-600">
          Monitor your equipment installations and service tickets
        </p>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Factory className="h-5 w-5 text-blue-600" />
                <p className="text-sm font-medium text-gray-500">Equipment Manufactured</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-blue-600">{equipmentData?.total || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Total installations</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Wrench className="h-5 w-5 text-orange-600" />
                <p className="text-sm font-medium text-gray-500">Active Service Tickets</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-orange-600">{ticketsData?.total || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Requires attention</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <AlertCircle className="h-5 w-5 text-red-600" />
                <p className="text-sm font-medium text-gray-500">Pending Tickets</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-red-600">{ticketsData?.pending || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Awaiting response</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <TrendingUp className="h-5 w-5 text-green-600" />
                <p className="text-sm font-medium text-gray-500">Resolution Rate</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-green-600">
                  {ticketsData?.total ? Math.round((ticketsData.resolved / ticketsData.total) * 100) : 0}%
                </p>
              )}
              <p className="text-xs text-gray-400 mt-1">This month</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* AI Onboarding CTA */}
      {(!equipmentData?.total || equipmentData.total === 0) && (
        <Card className="border-blue-200 bg-gradient-to-br from-blue-50 to-indigo-50">
          <CardContent className="pt-6">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-xl bg-blue-600 flex items-center justify-center flex-shrink-0">
                <Sparkles className="w-6 h-6 text-white" />
              </div>
              <div className="flex-1">
                <h3 className="text-xl font-bold text-gray-900 mb-2">
                  Ã°Å¸Å¡â‚¬ Get Started with AI Onboarding
                </h3>
                <p className="text-gray-700 mb-4">
                  Import your equipment, parts catalog, hospitals, and engineers in minutes using our AI-powered assistant. 
                  Simply upload CSV files and let AI handle the rest!
                </p>
                <div className="flex gap-3">
                  <Button 
                    onClick={() => router.push('/onboarding/ai-wizard')}
                    className="bg-blue-600 hover:bg-blue-700"
                  >
                    <Sparkles className="mr-2 h-4 w-4" />
                    Start AI Onboarding
                  </Button>
                  <Button 
                    onClick={() => router.push('/onboarding/wizard')}
                    variant="outline"
                  >
                    <Upload className="mr-2 h-4 w-4" />
                    Manual Import
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <Package className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <CardTitle>Equipment Registry</CardTitle>
                <CardDescription>View all your manufactured equipment</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              Track installations, warranty status, and service history across all your equipment
            </p>
            <Button 
              onClick={() => router.push('/equipment')}
              className="w-full"
            >
              View Equipment
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center">
                <Wrench className="w-5 h-5 text-orange-600" />
              </div>
              <div>
                <CardTitle>Service Tickets</CardTitle>
                <CardDescription>Manage service requests for your equipment</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {ticketsData?.pending || 0} tickets awaiting your response, {ticketsData?.inProgress || 0} in progress
            </p>
            <Button 
              onClick={() => router.push('/tickets')}
              className="w-full"
              variant={ticketsData?.pending ? "default" : "secondary"}
            >
              Manage Tickets
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle>Equipment Performance</CardTitle>
          <CardDescription>Overview of your equipment across all installations</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-500">
            <TrendingUp className="h-12 w-12 mx-auto mb-4 text-gray-400" />
            <p>Equipment performance analytics coming soon</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
