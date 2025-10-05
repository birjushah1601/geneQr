'use client';

import { useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Building2, 
  Package, 
  Users, 
  Ticket, 
  Store,
  ArrowRight,
  QrCode,
  TestTube2 as TestTube
} from 'lucide-react';

export default function AdminDashboard() {
  const router = useRouter();

  // Platform-wide stats - hardcoded to avoid hydration issues
  const platformStats = {
    manufacturers: 5,
    suppliers: 7,
    totalEquipment: 505,
    totalEngineers: 90,
    activeTickets: 23,
  };

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
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-8">
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Manufacturers</p>
                <p className="text-3xl font-bold mt-2">{platformStats.manufacturers}</p>
                <p className="text-xs text-gray-400 mt-1">Active partners</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Suppliers</p>
                <p className="text-3xl font-bold mt-2">{platformStats.suppliers}</p>
                <p className="text-xs text-gray-400 mt-1">Registered vendors</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Equipment</p>
                <p className="text-3xl font-bold mt-2">{platformStats.totalEquipment}</p>
                <p className="text-xs text-gray-400 mt-1">Platform-wide</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Engineers</p>
                <p className="text-3xl font-bold mt-2">{platformStats.totalEngineers}</p>
                <p className="text-xs text-gray-400 mt-1">Service personnel</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Active Tickets</p>
                <p className="text-3xl font-bold mt-2 text-orange-600">{platformStats.activeTickets}</p>
                <p className="text-xs text-gray-400 mt-1">Open requests</p>
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

          {/* Suppliers Card */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                    <Store className="w-5 h-5 text-purple-600" />
                  </div>
                  <div>
                    <CardTitle className="text-lg">Suppliers</CardTitle>
                    <CardDescription>Manage supply chain partners</CardDescription>
                  </div>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {platformStats.suppliers} registered suppliers providing equipment, parts, and medical consumables to the platform with excellent service ratings.
              </p>
              <div className="space-y-2 mb-4">
                {/* Top Suppliers Preview */}
                <div className="text-sm">
                  <p className="font-medium text-gray-700 mb-2">Top Suppliers:</p>
                  <ul className="space-y-1 text-gray-600">
                    <li>• HealthCare Solutions - ⭐⭐⭐⭐⭐ 4.8</li>
                    <li>• Precision Med Parts - ⭐⭐⭐⭐⭐ 4.6</li>
                    <li>• MedTech Supplies - ⭐⭐⭐⭐ 4.5</li>
                  </ul>
                </div>
              </div>
              <Button 
                onClick={() => router.push('/suppliers')}
                className="w-full"
              >
                View All Suppliers
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
