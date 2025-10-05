'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Building2, 
  Package, 
  Users, 
  Ticket, 
  Upload, 
  Plus,
  CheckCircle,
  AlertCircle
} from 'lucide-react';

export default function ManufacturerDashboard() {
  const router = useRouter();
  const [manufacturer, setManufacturer] = useState<any>(null);
  const [stats, setStats] = useState({
    equipment: 0,
    engineers: 0,
    tickets: 0,
  });

  useEffect(() => {
    // Load data from localStorage
    const mfr = localStorage.getItem('current_manufacturer');
    const equipment = localStorage.getItem('equipment_imported');
    const engineers = localStorage.getItem('engineers');
    
    if (mfr) {
      setManufacturer(JSON.parse(mfr));
    }
    
    setStats({
      equipment: equipment === 'true' ? 398 : 0,
      engineers: engineers ? JSON.parse(engineers).length : 0,
      tickets: 0,
    });
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Manufacturer Portal</h1>
              <p className="text-gray-600">
                {manufacturer?.name || 'Welcome'}
              </p>
            </div>
            <div className="flex items-center gap-3">
              <div className="text-right">
                <p className="text-sm font-medium">{manufacturer?.contact_person || 'Admin'}</p>
                <p className="text-xs text-gray-500">{manufacturer?.email || ''}</p>
              </div>
              <div className="w-10 h-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
                {manufacturer?.name?.charAt(0) || 'A'}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        {/* Welcome Message */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold mb-2">Dashboard</h2>
          <p className="text-gray-600">
            Manage your equipment, engineers, and service tickets
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Equipment</p>
                  <p className="text-3xl font-bold mt-2">{stats.equipment}</p>
                  <p className="text-xs text-gray-400 mt-1">Registered devices</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
                  <Package className="w-6 h-6 text-blue-600" />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Engineers</p>
                  <p className="text-3xl font-bold mt-2">{stats.engineers}</p>
                  <p className="text-xs text-gray-400 mt-1">Service team</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
                  <Users className="w-6 h-6 text-green-600" />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Tickets</p>
                  <p className="text-3xl font-bold mt-2">{stats.tickets}</p>
                  <p className="text-xs text-gray-400 mt-1">Active requests</p>
                </div>
                <div className="w-12 h-12 rounded-full bg-orange-100 flex items-center justify-center">
                  <Ticket className="w-6 h-6 text-orange-600" />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {/* Equipment Import */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                    <Package className="w-5 h-5 text-blue-600" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Equipment Registry</CardTitle>
                    <CardDescription>Import and manage equipment</CardDescription>
                  </div>
                </div>
                {stats.equipment > 0 && (
                  <CheckCircle className="w-5 h-5 text-green-600" />
                )}
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {stats.equipment > 0 
                  ? `${stats.equipment} equipment items registered. Import more or manage existing equipment.`
                  : 'Upload a CSV file to import your equipment installations with QR codes.'}
              </p>
              <div className="flex gap-3">
                <Button 
                  onClick={() => router.push('/equipment/import')}
                  className="flex-1"
                >
                  <Upload className="w-4 h-4 mr-2" />
                  Import CSV
                </Button>
                {stats.equipment > 0 && (
                  <Button 
                    variant="outline"
                    onClick={() => router.push('/equipment')}
                  >
                    View All
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Engineers */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
                    <Users className="w-5 h-5 text-green-600" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Service Engineers</CardTitle>
                    <CardDescription>Manage your service team</CardDescription>
                  </div>
                </div>
                {stats.engineers > 0 && (
                  <CheckCircle className="w-5 h-5 text-green-600" />
                )}
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {stats.engineers > 0 
                  ? `${stats.engineers} engineers in your team. Add more or manage assignments.`
                  : 'Add field engineers who will handle service requests and maintenance.'}
              </p>
              <div className="flex gap-3">
                {stats.engineers > 0 ? (
                  <>
                    <Button 
                      onClick={() => router.push('/engineers')}
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
                  </>
                ) : (
                  <>
                    <Button 
                      onClick={() => router.push('/engineers/import')}
                      className="flex-1"
                    >
                      <Upload className="w-4 h-4 mr-2" />
                      Import CSV
                    </Button>
                    <Button 
                      variant="outline"
                      onClick={() => router.push('/engineers/add')}
                    >
                      <Plus className="w-4 h-4 mr-2" />
                      Add Manually
                    </Button>
                  </>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Getting Started Guide */}
        {(stats.equipment === 0 || stats.engineers === 0) && (
          <Card className="border-orange-200 bg-orange-50">
            <CardHeader>
              <div className="flex items-center gap-3">
                <AlertCircle className="w-6 h-6 text-orange-600" />
                <div>
                  <CardTitle>Getting Started</CardTitle>
                  <CardDescription>Complete your setup to start managing service requests</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {stats.equipment === 0 && (
                  <div className="flex items-center gap-3 p-3 bg-white rounded-lg">
                    <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
                      <Package className="w-4 h-4 text-blue-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium">Import Equipment</p>
                      <p className="text-sm text-gray-600">Upload your equipment registry</p>
                    </div>
                    <Button 
                      size="sm"
                      onClick={() => router.push('/equipment/import')}
                    >
                      Start
                    </Button>
                  </div>
                )}
                
                {stats.engineers === 0 && (
                  <div className="flex items-center gap-3 p-3 bg-white rounded-lg">
                    <div className="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
                      <Users className="w-4 h-4 text-green-600" />
                    </div>
                    <div className="flex-1">
                      <p className="font-medium">Add Engineers</p>
                      <p className="text-sm text-gray-600">Build your service team</p>
                    </div>
                    <Button 
                      size="sm"
                      onClick={() => router.push('/engineers/import')}
                    >
                      Start
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        )}

        {/* Service Tickets (Coming Soon) */}
        <Card className="border-gray-200">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                <Ticket className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <CardTitle className="text-lg">Service Tickets</CardTitle>
                <CardDescription>View and manage service requests</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-600 mb-4">
              Service tickets will appear here once customers start reporting issues via WhatsApp or web.
            </p>
            <Button variant="outline" disabled>
              View Tickets
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
