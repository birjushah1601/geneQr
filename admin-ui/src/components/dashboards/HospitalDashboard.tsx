'use client';

import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Package, 
  Wrench, 
  Clock,
  ArrowRight,
  Loader2,
  Hospital,
  CheckCircle2,
  AlertTriangle
} from 'lucide-react';

export default function HospitalDashboard() {
  const router = useRouter();

  // Fetch hospital-specific data
  const { data: equipmentData, isLoading: loadingEquipment } = useQuery({
    queryKey: ['equipment', 'hospital'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/equipment?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      const data = await response.json();
      const equipment = data.equipment || [];
      return {
        total: equipment.length,
        operational: equipment.filter((e: any) => e.status === 'operational').length,
        maintenance: equipment.filter((e: any) => e.status === 'under_maintenance').length,
        underWarranty: equipment.filter((e: any) => e.warranty_expiry && new Date(e.warranty_expiry) > new Date()).length,
      };
    },
  });

  const { data: ticketsData, isLoading: loadingTickets } = useQuery({
    queryKey: ['tickets', 'hospital'],
    queryFn: async () => {
      const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';
      const token = localStorage.getItem('access_token');
      const response = await fetch(`${apiBaseUrl}/v1/tickets?limit=1000`, {
        headers: { 
          'Authorization': `Bearer ${token}`
        }
      });
      if (!response.ok) return { total: 0, open: 0, pending: 0 };
      const data = await response.json();
      const tickets = data.tickets || [];
      return {
        total: tickets.length,
        open: tickets.filter((t: any) => ['pending', 'acknowledged', 'in_progress'].includes(t.status)).length,
        pending: tickets.filter((t: any) => t.status === 'pending').length,
        inProgress: tickets.filter((t: any) => t.status === 'in_progress').length,
      };
    },
  });

  const isLoading = loadingEquipment || loadingTickets;

  return (
    <div className="space-y-6">
      {/* Welcome Section */}
      <div>
        <h1 className="text-3xl font-bold mb-2">Hospital Dashboard</h1>
        <p className="text-gray-600">
          Manage your medical equipment and service requests
        </p>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Package className="h-5 w-5 text-blue-600" />
                <p className="text-sm font-medium text-gray-500">Total Equipment</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-blue-600">{equipmentData?.total || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">In your facility</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <CheckCircle2 className="h-5 w-5 text-green-600" />
                <p className="text-sm font-medium text-gray-500">Operational</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-green-600">{equipmentData?.operational || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Ready for use</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="h-5 w-5 text-orange-600" />
                <p className="text-sm font-medium text-gray-500">Service Requests</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-orange-600">{ticketsData?.open || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Active tickets</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex flex-col">
              <div className="flex items-center gap-2 mb-2">
                <Clock className="h-5 w-5 text-purple-600" />
                <p className="text-sm font-medium text-gray-500">Under Warranty</p>
              </div>
              {isLoading ? (
                <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
              ) : (
                <p className="text-3xl font-bold text-purple-600">{equipmentData?.underWarranty || 0}</p>
              )}
              <p className="text-xs text-gray-400 mt-1">Equipment covered</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="hover:shadow-md transition-shadow">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                <Package className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <CardTitle>My Equipment</CardTitle>
                <CardDescription>View and manage your medical equipment</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {equipmentData?.operational || 0} operational, {equipmentData?.maintenance || 0} under maintenance
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
                <CardTitle>Service Requests</CardTitle>
                <CardDescription>Track your maintenance and repair requests</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              {ticketsData?.pending || 0} pending, {ticketsData?.inProgress || 0} in progress
            </p>
            <Button 
              onClick={() => router.push('/tickets')}
              className="w-full"
              variant={ticketsData?.pending ? "default" : "secondary"}
            >
              View Requests
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Create New Request */}
      <Card className="border-2 border-dashed border-blue-300 bg-blue-50">
        <CardContent className="pt-6">
          <div className="text-center">
            <Hospital className="h-12 w-12 mx-auto mb-4 text-blue-600" />
            <h3 className="text-lg font-semibold mb-2">Need Equipment Service?</h3>
            <p className="text-gray-600 mb-4">
              Create a new service request for any equipment issue
            </p>
            <Button 
              onClick={() => router.push('/service-request')}
              size="lg"
            >
              Create Service Request
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Equipment Status Overview */}
      <Card>
        <CardHeader>
          <CardTitle>Equipment Status</CardTitle>
          <CardDescription>Real-time status of all your medical equipment</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-500">
            <Package className="h-12 w-12 mx-auto mb-4 text-gray-400" />
            <p>Equipment status visualization coming soon</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
