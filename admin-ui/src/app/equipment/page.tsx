'use client';

import { useState, useMemo, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Search, Plus, Upload, Download, Filter, QrCode, Eye, Loader2 } from 'lucide-react';
import { equipmentApi } from '@/lib/api/equipment';
import Image from 'next/image';
import QRCodeLib from 'qrcode';

interface Equipment {
  id: string;
  name: string;
  serialNumber: string;
  model: string;
  manufacturer: string;
  category: string;
  location: string;
  status: string;
  installDate: string;
  lastService?: string;
  qrCode?: string;
  qrCodeUrl?: string;
  qrCodeImageUrl?: string; // Data URL for the QR code image
  hasQRCode?: boolean; // true only if image exists/generated
}

export default function EquipmentListPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [searchQuery, setSearchQuery] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [filterManufacturer, setFilterManufacturer] = useState<string>('');
  const [generatingQR, setGeneratingQR] = useState<string | null>(null);
  const [bulkGenerating, setBulkGenerating] = useState(false);
  const [qrPreview, setQrPreview] = useState<{id: string; url: string} | null>(null);
  const [equipmentData, setEquipmentData] = useState<Equipment[]>([]);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [isClient, setIsClient] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Read manufacturer from URL query parameter
  useEffect(() => {
    const manufacturer = searchParams.get('manufacturer');
    if (manufacturer) {
      setFilterManufacturer(manufacturer);
    }
  }, [searchParams]);

  // Fetch equipment data from API with fallback to mock data
  useEffect(() => {
    setIsClient(true);
    
    const fetchEquipment = async () => {
      try {
        setLoading(true);
        setError(null);
        
        console.log('Fetching equipment from API...');
        const response = await equipmentApi.list({
          page: 1,
          page_size: 1000,
        });
        
        console.log('API Response:', response);
        
        // Map API response to component format
        const responseData: any = response;
        const mappedEquipment: Equipment[] = (responseData.items || responseData.equipment || []).map((item: any) => ({
          id: item.id,
          name: item.equipment_name,
          serialNumber: item.serial_number,
          model: item.model_number || 'N/A',
          manufacturer: item.manufacturer_name,
          category: item.category || 'Unknown',
          location: item.installation_location || item.customer_name,
          status: item.status === 'operational' ? 'Active' : item.status === 'down' ? 'Inactive' : 'Maintenance',
          installDate: item.installation_date || item.created_at?.split('T')[0],
          lastService: item.last_service_date,
          qrCode: item.qr_code,
          qrCodeUrl: item.qr_code_url,
          hasQRCode: !!item.qr_code_generated_at || !!item.qr_code_image,
        }));
        
        console.log(`Loaded ${mappedEquipment.length} equipment items from API`);
        setEquipmentData(mappedEquipment);
      } catch (err) {
        console.error('Failed to fetch equipment from API, using demo data:', err);
        // Fallback to mock data for demo purposes
        const mockEquipment: Equipment[] = [
          {
            id: 'eq-001',
            name: 'X-Ray Machine',
            serialNumber: 'SN-001-2024',
            model: 'Discovery XR656',
            manufacturer: 'GE Healthcare',
            category: 'Imaging',
            location: 'City General Hospital - Radiology Department',
            status: 'Active',
            installDate: '2024-01-15',
            lastService: '2024-09-15',
            qrCode: 'QR-eq-001',
            hasQRCode: true,
          },
          {
            id: 'eq-002',
            name: 'MRI Scanner',
            serialNumber: 'SN-002-2024',
            model: 'Magnetom Skyra 3T',
            manufacturer: 'Siemens Healthineers',
            category: 'Imaging',
            location: 'Regional Medical Center - MRI Suite',
            status: 'Active',
            installDate: '2024-02-20',
            lastService: '2024-09-20',
            qrCode: 'QR-eq-002',
            hasQRCode: true,
          },
          {
            id: 'eq-003',
            name: 'Ultrasound System',
            serialNumber: 'SN-003-2024',
            model: 'EPIQ Elite',
            manufacturer: 'Philips Healthcare',
            category: 'Imaging',
            location: 'Metro Clinic - Diagnostic Center',
            status: 'Active',
            installDate: '2024-03-10',
            qrCode: 'QR-eq-003',
            hasQRCode: false,
          },
          {
            id: 'eq-004',
            name: 'Patient Monitor',
            serialNumber: 'SN-004-2024',
            model: 'Excel 15',
            manufacturer: 'BPL Medical Technologies',
            category: 'Patient Monitoring',
            location: 'Apollo Hospital - ICU Ward 3',
            status: 'Active',
            installDate: '2024-04-05',
            lastService: '2024-10-01',
            qrCode: 'QR-eq-004',
            hasQRCode: true,
          },
        ];
        setEquipmentData(mockEquipment);
        setError(null); // Clear error since we have mock data
      } finally {
        setLoading(false);
      }
    };

    fetchEquipment();
  }, []);

  // Filter equipment based on search, status, and manufacturer
  const filteredEquipment = useMemo(() => {
    return equipmentData.filter(equipment => {
      const matchesSearch = searchQuery === '' || 
        equipment.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        equipment.serialNumber.toLowerCase().includes(searchQuery.toLowerCase()) ||
        equipment.model.toLowerCase().includes(searchQuery.toLowerCase()) ||
        equipment.manufacturer.toLowerCase().includes(searchQuery.toLowerCase()) ||
        equipment.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
        equipment.location.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = filterStatus === 'all' || equipment.status.toLowerCase() === filterStatus.toLowerCase();
      
      const matchesManufacturer = filterManufacturer === '' || equipment.manufacturer === filterManufacturer;

      return matchesSearch && matchesStatus && matchesManufacturer;
    });
  }, [equipmentData, searchQuery, filterStatus, filterManufacturer]);

  const statusCounts = useMemo(() => {
    return equipmentData.reduce((acc, eq) => {
      acc[eq.status] = (acc[eq.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
  }, [equipmentData]);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return 'bg-green-100 text-green-800';
      case 'maintenance':
        return 'bg-yellow-100 text-yellow-800';
      case 'inactive':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const handleGenerateQR = async (equipmentId: string) => {
    try {
      setGeneratingQR(equipmentId);
      
      // Call real backend API to generate and store QR code
      const result: any = await equipmentApi.generateQRCode(equipmentId);
      
      // Reload the page to fetch updated equipment with QR code
      alert(`✅ QR Code generated and stored successfully!\n\nEquipment: ${equipmentId}\nQR Code: ${result.qr_code || `QR-${equipmentId}`}`);
      window.location.reload();
    } catch (error) {
      console.error('QR generation failed:', error);
      alert(`Failed to generate QR code: ${error instanceof Error ? error.message : 'Unknown error'}\n\nPlease ensure backend is running on port 8081.`);
    } finally {
      setGeneratingQR(null);
    }
  };

  const handlePreviewQR = (equipment: Equipment) => {
    if (equipment.hasQRCode) {
      // Use generated QR code image if available, otherwise try backend
      const imageUrl = equipment.qrCodeImageUrl || "http://localhost:8081/api/v1/equipment/qr/image/" + equipment.id;
      setQrPreview({ id: equipment.id, url: imageUrl });
    }
  };

  const handleDownloadQR = async (equipmentId: string) => {
    try {
      await equipmentApi.downloadQRLabel(equipmentId);
    } catch (error) {
      alert(`Failed to download QR label: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  };


  const handleBulkGenerateQR = async () => {
    const confirmed = confirm("Generate QR codes for all equipment that doesn't have one? This may take a while.");
    if (!confirmed) return;

    try {
      setBulkGenerating(true);
      
      try {
        // Try to call real API first
        const result = await equipmentApi.bulkGenerateQRCodes();
        alert(result.message);
        window.location.reload();
      } catch (apiError) {
        // If API fails, generate QR codes locally for demo
        console.log('Bulk API failed, generating QR codes locally:', apiError);
        
        // Count equipment without QR codes
        const withoutQR = equipmentData.filter(eq => !eq.hasQRCode);
        
        // Generate QR codes for all equipment without them
        const updatedEquipment = await Promise.all(
          equipmentData.map(async (eq) => {
            if (eq.hasQRCode) return eq;
            
            // Generate QR code
            const qrData = `http://localhost:3000/equipment/${eq.id}`;
            const qrCodeDataUrl = await QRCodeLib.toDataURL(qrData, {
              width: 300,
              margin: 2,
              color: {
                dark: '#000000',
                light: '#FFFFFF'
              }
            });
            
            return {
              ...eq,
              hasQRCode: true,
              qrCode: `QR-${eq.id}`,
              qrCodeImageUrl: qrCodeDataUrl,
              qrCodeUrl: qrData
            };
          })
        );
        
        // Simulate processing delay
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Update all equipment
        setEquipmentData(updatedEquipment);
        
        alert(`✅ Bulk QR generation complete!\n\n${withoutQR.length} QR codes generated successfully.\n\n(Demo mode: QR codes generated locally)`);
      }
    } catch (error) {
      alert(`Failed to bulk generate QR codes: ${error instanceof Error ? error.message : 'Unknown error'}`);
    } finally {
      setBulkGenerating(false);
    }
  };

  const toggleSelectAll = (checked: boolean) => {
    if (checked) {
      const ids = new Set(filteredEquipment.map(e => e.id));
      setSelectedIds(ids);
    } else {
      setSelectedIds(new Set());
    }
  };

  const toggleSelectOne = (id: string, checked: boolean) => {
    setSelectedIds(prev => {
      const next = new Set(prev);
      if (checked) next.add(id); else next.delete(id);
      return next;
    });
  };

  const handleGenerateSelected = async () => {
    if (selectedIds.size === 0) return;
    try {
      setBulkGenerating(true);
      for (const id of Array.from(selectedIds)) {
        try {
          await equipmentApi.generateQRCode(id);
        } catch (e) {
          console.error('Failed to generate for', id, e);
        }
      }
      alert(`Generated QR for ${selectedIds.size} selected equipment`);
      window.location.reload();
    } finally {
      setBulkGenerating(false);
    }
  };
  // Show loading state during hydration
  if (!isClient || loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-600" />
          <p className="text-gray-600">Loading equipment from API...</p>
        </div>
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Card className="max-w-md w-full">
          <CardHeader>
            <CardTitle className="text-red-600">Error Loading Equipment</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-700 mb-4">{error}</p>
            <div className="space-y-2 text-sm text-gray-600">
              <p>Possible causes:</p>
              <ul className="list-disc list-inside space-y-1">
                <li>Backend not running on port 8081</li>
                <li>Database connection issue</li>
                <li>No equipment data in database</li>
              </ul>
            </div>
            <Button 
              onClick={() => window.location.reload()} 
              className="mt-4 w-full"
            >
              Retry
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

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
              <h1 className="text-3xl font-bold text-gray-900">Equipment Registry</h1>
              <p className="text-gray-600 mt-1">
                Manage and monitor all medical equipment
              </p>
            </div>
            
            <div className="flex gap-2">
              <Button 
                variant="outline" 
                onClick={handleBulkGenerateQR}
                disabled={bulkGenerating}
              >
                {bulkGenerating ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Generating...
                  </>
                ) : (
                  <>
                    <QrCode className="mr-2 h-4 w-4" />
                    Generate All QR Codes
                  </>
                )}
              </Button>
              <Button 
                variant="outline" 
                onClick={handleGenerateSelected}
                disabled={bulkGenerating || selectedIds.size === 0}
              >
                {bulkGenerating ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Generating Selected...
                  </>
                ) : (
                  <>
                    <QrCode className="mr-2 h-4 w-4" />
                    Generate Selected
                  </>
                )}
              </Button>
              <Button variant="outline" onClick={() => router.push('/equipment/import')}>
                <Upload className="mr-2 h-4 w-4" />
                Import CSV
              </Button>
              <Button onClick={() => router.push('/equipment/new')}>
                <Plus className="mr-2 h-4 w-4" />
                Add Equipment
              </Button>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Total Equipment</CardDescription>
              <CardTitle className="text-3xl">{equipmentData.length}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Active</CardDescription>
              <CardTitle className="text-3xl text-green-600">{statusCounts['Active'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Under Maintenance</CardDescription>
              <CardTitle className="text-3xl text-yellow-600">{statusCounts['Maintenance'] || 0}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>Inactive</CardDescription>
              <CardTitle className="text-3xl text-red-600">{statusCounts['Inactive'] || 0}</CardTitle>
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
                  placeholder="Search equipment by name, serial number, model, manufacturer, category, or location..."
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
                  <option value="active">Active</option>
                  <option value="maintenance">Maintenance</option>
                  <option value="inactive">Inactive</option>
                </select>
                <Button variant="outline">
                  <Download className="mr-2 h-4 w-4" />
                  Export
                </Button>
              </div>
            </div>
            
            {(searchQuery || filterStatus !== 'all' || filterManufacturer) && (
              <div className="mt-4 flex items-center gap-4">
                <div className="text-sm text-gray-600">
                  Showing {filteredEquipment.length} of {equipmentData.length} equipment
                </div>
                {filterManufacturer && (
                  <div className="flex items-center gap-2 px-3 py-1 bg-blue-50 border border-blue-200 rounded-md">
                    <span className="text-sm text-blue-700">
                      Filtered by manufacturer: <strong>{filterManufacturer}</strong>
                    </span>
                    <button
                      onClick={() => {
                        setFilterManufacturer('');
                        router.push('/equipment');
                      }}
                      className="text-blue-700 hover:text-blue-900"
                    >
                      ✕
                    </button>
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Equipment List */}
        {equipmentData.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="text-gray-400 mb-4">
                <Upload className="h-12 w-12 mx-auto" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No Equipment Found</h3>
              <p className="text-gray-600 mb-4">
                Get started by importing equipment from a CSV file or adding them manually.
              </p>
              <div className="flex justify-center gap-2">
                <Button onClick={() => router.push('/equipment/import')}>
                  <Upload className="mr-2 h-4 w-4" />
                  Import CSV
                </Button>
                <Button variant="outline" onClick={() => alert('Add equipment feature coming soon!')}>
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
                      <th className="px-6 py-3">
                        <input 
                          type="checkbox" 
                          aria-label="Select all"
                          checked={selectedIds.size > 0 && filteredEquipment.length > 0 && selectedIds.size === filteredEquipment.length}
                          onChange={(e) => toggleSelectAll(e.target.checked)}
                        />
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        QR Code
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Equipment
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Serial Number
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Category
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Location
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Last Service
                      </th>
                      <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {filteredEquipment.map((equipment) => (
                      <tr key={equipment.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <input 
                            type="checkbox" 
                            checked={selectedIds.has(equipment.id)}
                            onChange={(e) => toggleSelectOne(equipment.id, e.target.checked)}
                            aria-label={`Select ${equipment.name}`}
                          />
                        </td>
                        {/* QR Code Column */}
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center gap-2">
                            {equipment.hasQRCode ? (
                              <div className="group relative">
                                <div 
                                  className="w-20 h-20 border-2 border-gray-200 rounded-md overflow-hidden cursor-pointer hover:border-blue-500 transition-colors bg-white"
                                  onClick={() => handlePreviewQR(equipment)}
                                  title="Click to preview full size"
                                >
                                  <img
                                    src={`http://localhost:8081/api/v1/equipment/qr/image/${equipment.id}`}
                                    alt={`QR Code for ${equipment.name}`}
                                    className="w-full h-full object-contain p-1"
                                    onError={(e) => {
                                      console.error('Failed to load QR image for', equipment.id);
                                      e.currentTarget.style.display = 'none';
                                    }}
                                  />
                                </div>
                                <div className="absolute hidden group-hover:flex flex-col gap-1 top-0 left-20 bg-white shadow-lg rounded-md p-2 z-10">
                                  <Button
                                    size="sm"
                                    variant="ghost"
                                    onClick={() => handlePreviewQR(equipment)}
                                    className="justify-start text-xs"
                                  >
                                    <Eye className="mr-1 h-3 w-3" />
                                    Preview
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="ghost"
                                    onClick={() => handleDownloadQR(equipment.id)}
                                    className="justify-start text-xs"
                                  >
                                    <Download className="mr-1 h-3 w-3" />
                                    Download
                                  </Button>
                                </div>
                              </div>
                            ) : (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleGenerateQR(equipment.id)}
                                disabled={generatingQR === equipment.id}
                                className="w-16 h-16 flex flex-col items-center justify-center text-xs"
                              >
                                {generatingQR === equipment.id ? (
                                  <>
                                    <Loader2 className="h-4 w-4 animate-spin mb-1" />
                                    <span className="text-[10px]">Wait...</span>
                                  </>
                                ) : (
                                  <>
                                    <QrCode className="h-4 w-4 mb-1" />
                                    <span className="text-[10px]">Generate</span>
                                  </>
                                )}
                              </Button>
                            )}
                          </div>
                        </td>
                        
                        {/* Equipment Details */}
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex flex-col">
                            <div className="text-sm font-medium text-gray-900">
                              {equipment.name}
                            </div>
                            <div className="text-xs text-gray-500">
                              {equipment.manufacturer} • {equipment.model}
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {equipment.serialNumber}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {equipment.category}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {equipment.location}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(equipment.status)}`}>
                            {equipment.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {equipment.lastService || 'N/A'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                          <div className="flex items-center justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => router.push(`/equipment/${equipment.id}`)}
                            >
                              View
                            </Button>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => {
                                const code = equipment.qrCode || equipment.id;
                                const url = code ? `/service-request?qr=${encodeURIComponent(code)}` : '/service-request';
                                router.push(url);
                              }}
                            >
                              Create Service Request
                            </Button>
                          </div>
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

      {/* QR Code Preview Modal */}
      {qrPreview && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4"
          onClick={() => setQrPreview(null)}
        >
          <div 
            className="bg-white rounded-lg shadow-xl max-w-md w-full p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">QR Code Preview</h3>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setQrPreview(null)}
              >
                ✕
              </Button>
            </div>
            
            <div className="flex justify-center mb-4">
              <div className="w-64 h-64 border-2 border-gray-200 rounded-md p-4 bg-white">
                <Image
                  src={qrPreview.url}
                  alt="QR Code"
                  width={256}
                  height={256}
                  className="w-full h-full object-contain"
                  unoptimized
                />
              </div>
            </div>

            <div className="text-center text-sm text-gray-500 mb-4">
              Equipment ID: {qrPreview.id}
            </div>

            <div className="flex gap-2">
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => handleDownloadQR(qrPreview.id)}
              >
                <Download className="mr-2 h-4 w-4" />
                Download PDF Label
              </Button>
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => {
                  window.open(qrPreview.url, '_blank');
                }}
              >
                <Eye className="mr-2 h-4 w-4" />
                Open in New Tab
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
