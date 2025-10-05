'use client';

import { useParams, useRouter } from 'next/navigation';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { 
  Store, 
  Package, 
  FileText, 
  TrendingUp,
  ArrowLeft,
  Star
} from 'lucide-react';

export default function SupplierDashboard() {
  const params = useParams();
  const router = useRouter();
  const supplierId = params.id as string;

  // Mock supplier data - in production this would come from API
  const supplierData: Record<string, any> = {
    'SUP-001': {
      id: 'SUP-001',
      name: 'MedTech Supplies India',
      contactPerson: 'Rajesh Kumar',
      email: 'rajesh@medtechsupplies.com',
      phone: '+91-9876543220',
      category: 'Medical Consumables',
      location: 'Mumbai, Maharashtra',
      rating: 4.5,
      totalOrders: 145,
      activeContracts: 3,
      completedOrders: 138,
      pendingOrders: 7,
      revenue: '₹12,50,000',
      status: 'Active',
      createdAt: '2024-01-10',
    },
    'SUP-002': {
      id: 'SUP-002',
      name: 'HealthCare Solutions Ltd',
      contactPerson: 'Anita Desai',
      email: 'anita@healthcaresolutions.com',
      phone: '+91-9876543221',
      category: 'Equipment Parts',
      location: 'Bangalore, Karnataka',
      rating: 4.8,
      totalOrders: 230,
      activeContracts: 5,
      completedOrders: 218,
      pendingOrders: 12,
      revenue: '₹22,80,000',
      status: 'Active',
      createdAt: '2023-11-15',
    },
    'SUP-003': {
      id: 'SUP-003',
      name: 'Bio Medical Instruments',
      contactPerson: 'Suresh Patel',
      email: 'suresh@biomedical.com',
      phone: '+91-9876543222',
      category: 'Diagnostic Equipment',
      location: 'Delhi, India',
      rating: 4.2,
      totalOrders: 89,
      activeContracts: 2,
      completedOrders: 82,
      pendingOrders: 7,
      revenue: '₹8,90,000',
      status: 'Active',
      createdAt: '2024-02-20',
    },
    'SUP-004': {
      id: 'SUP-004',
      name: 'Precision Med Parts',
      contactPerson: 'Meera Singh',
      email: 'meera@precisionmedparts.com',
      phone: '+91-9876543223',
      category: 'Spare Parts',
      location: 'Pune, Maharashtra',
      rating: 4.6,
      totalOrders: 178,
      activeContracts: 4,
      completedOrders: 165,
      pendingOrders: 13,
      revenue: '₹17,80,000',
      status: 'Active',
      createdAt: '2023-12-05',
    },
    'SUP-005': {
      id: 'SUP-005',
      name: 'Global Medical Supplies',
      contactPerson: 'Arjun Mehta',
      email: 'arjun@globalmedsupplies.com',
      phone: '+91-9876543224',
      category: 'General Supplies',
      location: 'Hyderabad, Telangana',
      rating: 4.4,
      totalOrders: 156,
      activeContracts: 3,
      completedOrders: 148,
      pendingOrders: 8,
      revenue: '₹15,60,000',
      status: 'Active',
      createdAt: '2024-03-12',
    },
    'SUP-006': {
      id: 'SUP-006',
      name: 'Advanced Healthcare Products',
      contactPerson: 'Kavita Reddy',
      email: 'kavita@advancedhealthcare.com',
      phone: '+91-9876543225',
      category: 'Medical Consumables',
      location: 'Chennai, Tamil Nadu',
      rating: 3.9,
      totalOrders: 67,
      activeContracts: 1,
      completedOrders: 60,
      pendingOrders: 7,
      revenue: '₹6,70,000',
      status: 'Pending',
      createdAt: '2024-08-01',
    },
    'SUP-007': {
      id: 'SUP-007',
      name: 'Quality Med Equipment',
      contactPerson: 'Vikram Sharma',
      email: 'vikram@qualitymedequip.com',
      phone: '+91-9876543226',
      category: 'Equipment Parts',
      location: 'Ahmedabad, Gujarat',
      rating: 3.2,
      totalOrders: 34,
      activeContracts: 0,
      completedOrders: 28,
      pendingOrders: 6,
      revenue: '₹3,40,000',
      status: 'Inactive',
      createdAt: '2023-09-15',
    },
  };

  const supplier = supplierData[supplierId];

  if (!supplier) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Card className="max-w-md">
          <CardContent className="pt-6">
            <div className="text-center">
              <Store className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h2 className="text-xl font-bold mb-2">Supplier Not Found</h2>
              <p className="text-gray-600 mb-4">
                The supplier with ID "{supplierId}" does not exist.
              </p>
              <Button onClick={() => router.push('/suppliers')}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Suppliers
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const getRatingStars = (rating: number) => {
    return '⭐'.repeat(Math.floor(rating));
  };

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
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-4">
          <Button
            variant="ghost"
            onClick={() => router.push('/suppliers')}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Suppliers
          </Button>
          
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full bg-purple-500 flex items-center justify-center text-white text-2xl font-bold">
                {supplier.name.substring(0, 2).toUpperCase()}
              </div>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{supplier.name}</h1>
                <p className="text-gray-600">{supplier.location}</p>
                <div className="flex items-center gap-2 mt-1">
                  <span className="text-sm">{getRatingStars(supplier.rating)} {supplier.rating}</span>
                  <span className={`px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(supplier.status)}`}>
                    {supplier.status}
                  </span>
                </div>
              </div>
            </div>
            
            <div className="text-right">
              <p className="text-sm font-medium">{supplier.contactPerson}</p>
              <p className="text-xs text-gray-500">{supplier.email}</p>
              <p className="text-xs text-gray-500">{supplier.phone}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        {/* Welcome Message */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold mb-2">Supplier Dashboard</h2>
          <p className="text-gray-600">
            Manage orders, contracts, and performance metrics for {supplier.name}
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Total Orders</p>
                <p className="text-3xl font-bold mt-2">{supplier.totalOrders}</p>
                <p className="text-xs text-gray-400 mt-1">Lifetime orders</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Active Contracts</p>
                <p className="text-3xl font-bold mt-2">{supplier.activeContracts}</p>
                <p className="text-xs text-gray-400 mt-1">Current agreements</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Pending Orders</p>
                <p className="text-3xl font-bold mt-2 text-orange-600">{supplier.pendingOrders}</p>
                <p className="text-xs text-gray-400 mt-1">Awaiting fulfillment</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6">
              <div className="flex flex-col">
                <p className="text-sm font-medium text-gray-500">Total Revenue</p>
                <p className="text-2xl font-bold mt-2 text-green-600">{supplier.revenue}</p>
                <p className="text-xs text-gray-400 mt-1">Lifetime value</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Management Sections */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {/* Orders Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
                  <Package className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Order Management</CardTitle>
                  <CardDescription>Track and manage orders</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3 mb-4">
                <div className="flex justify-between items-center p-3 bg-green-50 rounded-lg">
                  <span className="text-sm text-gray-700">Completed Orders</span>
                  <span className="text-sm font-semibold text-green-700">{supplier.completedOrders}</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-orange-50 rounded-lg">
                  <span className="text-sm text-gray-700">Pending Orders</span>
                  <span className="text-sm font-semibold text-orange-700">{supplier.pendingOrders}</span>
                </div>
              </div>
              <Button 
                onClick={() => router.push(`/orders?supplier=${supplierId}`)}
                className="w-full"
              >
                View All Orders
              </Button>
            </CardContent>
          </Card>

          {/* Contracts Management */}
          <Card>
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
                  <FileText className="w-5 h-5 text-purple-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Contract Management</CardTitle>
                  <CardDescription>View active agreements</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600 mb-4">
                {supplier.activeContracts > 0 
                  ? `${supplier.activeContracts} active contracts with ongoing supply agreements.`
                  : 'No active contracts. Create new agreements to start ordering.'}
              </p>
              <Button 
                onClick={() => router.push(`/contracts?supplier=${supplierId}`)}
                className="w-full"
                disabled={supplier.activeContracts === 0}
              >
                View All Contracts
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Performance Metrics */}
        <Card className="mb-6">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
                <TrendingUp className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <CardTitle className="text-lg">Performance Metrics</CardTitle>
                <CardDescription>Supplier performance indicators</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <Star className="w-4 h-4 text-yellow-500" />
                  <p className="text-sm font-medium text-gray-700">Rating</p>
                </div>
                <p className="text-2xl font-bold">{supplier.rating} / 5.0</p>
                <p className="text-xs text-gray-500 mt-1">{getRatingStars(supplier.rating)}</p>
              </div>
              
              <div className="p-4 bg-gray-50 rounded-lg">
                <p className="text-sm font-medium text-gray-700 mb-2">Order Fulfillment</p>
                <p className="text-2xl font-bold">{Math.round((supplier.completedOrders / supplier.totalOrders) * 100)}%</p>
                <p className="text-xs text-gray-500 mt-1">{supplier.completedOrders} of {supplier.totalOrders} orders</p>
              </div>
              
              <div className="p-4 bg-gray-50 rounded-lg">
                <p className="text-sm font-medium text-gray-700 mb-2">Category</p>
                <p className="text-lg font-semibold">{supplier.category}</p>
                <p className="text-xs text-gray-500 mt-1">Specialization</p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Supplier Information */}
        <Card>
          <CardHeader>
            <CardTitle>Supplier Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-gray-500">Supplier ID</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.id}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Category</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.category}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Contact Person</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.contactPerson}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Email</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.email}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Phone</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.phone}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Location</p>
                <p className="text-sm text-gray-900 mt-1">{supplier.location}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Member Since</p>
                <p className="text-sm text-gray-900 mt-1">
                  {new Date(supplier.createdAt).toLocaleDateString('en-US', { 
                    month: 'long', 
                    day: 'numeric', 
                    year: 'numeric' 
                  })}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Status</p>
                <span className={`inline-block px-2 py-1 text-xs font-semibold rounded-full mt-1 ${getStatusColor(supplier.status)}`}>
                  {supplier.status}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
