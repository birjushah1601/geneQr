'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Package, 
  Users, 
  Wrench,
  ArrowRight,
  Loader2,
  Truck,
  CheckCircle2,
  Clock
} from 'lucide-react';

export default function ChannelPartnerDashboard() {
  const router = useRouter();

  // Fetch Channel Partner-specific data
  const { data: equipmentData, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'Channel Partner'],
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
    queryKey: ['tickets', 'Channel Partner'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/tickets?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      if (!response.ok) return { total: 0, assigned: 0, completed: 0 };
      const data = await response.json();
      const tickets = data.tickets || [];
      return {
        total: tickets.length,
        assigned: tickets.filter((t: any) => ['acknowledged', 'in_progress'].includes(t.status)).length,
        pending: tickets.filter((t: any) => t.status === 'pending').length,
        completed: tickets.filter((t: any) => ['resolved', 'closed'].includes(t.status)).length,
      };
    },
  });

  const { data: engineersData, isLoading: loadingEngineers } = useQuery({
    queryKey: ['engineers', 'Channel Partner'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/engineers?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      const data = await response.json();
      return { 
        total: data.engineers?.length || 0,
        engineers: data.engineers || []
      };
    },
  });

  const isLoading = loadingEquipment || loadingTickets || loadingEngineers;

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold mb-2">Channel Partner Dashboard</h1>
        <p className="text-gray-600">
          Manage service operations and engineer assignments
        </p>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Truck className="h-5 w-5 text-purple-600" />
                <p className="text-sm font-medium text-gray-500">Equipment Serviced</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-purple-600">{equipmentData?.total || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Under your service</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Wrench className="h-5 w-5 text-orange-600" />
                <p className="text-sm font-medium text-gray-500">Active Service Jobs</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-orange-600">{ticketsData?.assigned || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">In progress</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Users className="h-5 w-5 text-blue-600" />
                <p className="text-sm font-medium text-gray-500">Service Engineers</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-blue-600">{engineersData?.total || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">In your team</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <CheckCircle2 className="h-5 w-5 text-green-600" />
                <p className="text-sm font-medium text-gray-500">Completed</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-green-600">{ticketsData?.completed || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">This month</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center">
                <Wrench className="w-5 h-5 text-orange-600" />
              </div>
              <div>
                <CardTitle>Service Jobs</CardTitle>
                <CardDescription>Manage assigned service tickets</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {ticketsData?.pending || 0} new requests, {ticketsData?.assigned || 0} in progress
            </p>
            <Button 
              onClick={() => router.push('/tickets')}
              className="w-full"
              variant={ticketsData?.pending ? "default" : "secondary"}
            >
              Manage Jobs
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <Users className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <CardTitle>Engineers</CardTitle>
                <CardDescription>Manage your service team</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {engineersData?.total || 0} engineers available for assignments
            </p>
            <Button 
              onClick={() => router.push('/engineers')}
              className="w-full"
            >
              View Engineers
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                <Package className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <CardTitle>Equipment</CardTitle>
                <CardDescription>Equipment under your service</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {equipmentData?.total || 0} equipment pieces serviced
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
      </div>

      {/* Pending Assignments Alert */}
      {ticketsData && ticketsData.pending > 0 && (
        <Card className="border-orange-300 bg-orange-50">
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <Clock className="h-10 w-10 text-orange-600" />
              <div className="flex-1">
                <h3 className="text-lg font-semibold mb-1">Pending Assignments</h3>
                <p className="text-gray-600">
                  {ticketsData.pending} service request(s) need engineer assignment
                </p>
              </div>
              <Button onClick={() => router.push('/tickets?status=pending')}>
                Assign Now
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Performance Overview */}
      <Card>
        <CardHeader>
          <CardTitle>Service Performance</CardTitle>
          <CardDescription>Track your service team's performance metrics</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-500">
            <Truck className="h-12 w-12 mx-auto mb-4 text-gray-400" />
            <p>Service performance analytics coming soon</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
