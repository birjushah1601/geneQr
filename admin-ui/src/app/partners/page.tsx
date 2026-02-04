'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import DashboardLayout from '@/components/DashboardLayout';
import { partnersApi, Partner, Organization } from '@/lib/api/partners';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  Users,
  Trash2,
  Plus,
  Search,
  Loader2,
  Building2,
  AlertCircle,
} from 'lucide-react';

export default function PartnersPage() {
  const { organizationContext } = useAuth();
  const [activeTab, setActiveTab] = useState('general');
  const [partners, setPartners] = useState<Partner[]>([]);
  const [availablePartners, setAvailablePartners] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [error, setError] = useState('');

  const manufacturerId = organizationContext?.organization_id || '';

  useEffect(() => {
    if (manufacturerId) {
      loadPartners();
    }
  }, [manufacturerId, activeTab]);

  const loadPartners = async () => {
    try {
      setLoading(true);
      const filters = activeTab === 'equipment-specific' 
        ? { association_type: 'equipment-specific' }
        : { association_type: 'general' };
      
      const data = await partnersApi.listPartners(manufacturerId, filters);
      setPartners(data.partners || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const loadAvailablePartners = async (search: string) => {
    try {
      const data = await partnersApi.getAvailablePartners(manufacturerId, search);
      setAvailablePartners(data.organizations || []);
    } catch (err: any) {
      console.error(err);
    }
  };

  const handleAssociatePartner = async (partnerId: string) => {
    try {
      await partnersApi.associatePartner(manufacturerId, {
        partner_org_id: partnerId,
      });
      setShowAddModal(false);
      loadPartners();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleRemovePartner = async (partnerId: string, equipmentId?: string) => {
    if (!confirm('Are you sure you want to remove this partner association?')) {
      return;
    }

    try {
      await partnersApi.removePartner(manufacturerId, partnerId, equipmentId);
      loadPartners();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const getOrgTypeLabel = (orgType: string) => {
    return orgType === 'channel_partner' ? 'Channel Partner' : 'Sub-Dealer';
  };

  if (organizationContext?.organization_type !== 'manufacturer') {
    return (
      <DashboardLayout>
        <div className="flex items-center justify-center min-h-[400px]">
          <div className="text-center">
            <AlertCircle className="w-12 h-12 text-yellow-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Access Restricted</h2>
            <p className="text-gray-600">
              Partner management is only available for manufacturers.
            </p>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">Service Partners</h1>
            <p className="text-gray-600 mt-1">
              Manage channel partners and sub-dealers for service coverage
            </p>
          </div>
          <Button onClick={() => {
            setShowAddModal(true);
            loadAvailablePartners('');
          }}>
            <Plus className="w-4 h-4 mr-2" />
            Add Partner
          </Button>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded">
            {error}
          </div>
        )}

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList>
            <TabsTrigger value="general">General Partners</TabsTrigger>
            <TabsTrigger value="equipment-specific">Equipment-Specific</TabsTrigger>
          </TabsList>

          <TabsContent value="general" className="mt-6">
            {loading ? (
              <div className="flex justify-center py-12">
                <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
              </div>
            ) : partners.length === 0 ? (
              <Card>
                <CardContent className="py-12 text-center text-gray-500">
                  No general partners associated yet. Click "Add Partner" to get started.
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-4">
                {partners.map((partner) => (
                  <Card key={partner.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div className="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
                            <Users className="w-6 h-6 text-blue-600" />
                          </div>
                          <div>
                            <h3 className="text-lg font-semibold">{partner.partner_name}</h3>
                            <p className="text-sm text-gray-600">
                              {getOrgTypeLabel(partner.org_type)}
                            </p>
                            <p className="text-sm text-gray-500 mt-1">
                              {partner.engineers_count} engineers • Services ALL equipment
                            </p>
                          </div>
                        </div>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => handleRemovePartner(partner.partner_org_id)}
                        >
                          <Trash2 className="w-4 h-4 mr-2" />
                          Remove
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>

          <TabsContent value="equipment-specific" className="mt-6">
            {loading ? (
              <div className="flex justify-center py-12">
                <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
              </div>
            ) : partners.length === 0 ? (
              <Card>
                <CardContent className="py-12 text-center text-gray-500">
                  No equipment-specific associations yet.
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-4">
                {partners.map((partner) => (
                  <Card key={partner.id}>
                    <CardContent className="p-6">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold">{partner.equipment_name || 'Equipment'}</h3>
                          <p className="text-sm text-gray-600 mt-1">
                            → {partner.partner_name} ({getOrgTypeLabel(partner.org_type)})
                          </p>
                          <p className="text-xs text-gray-500 mt-1">
                            {partner.engineers_count} engineers
                          </p>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleRemovePartner(partner.partner_org_id, partner.equipment_id)}
                        >
                          Remove Override
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>
        </Tabs>

        {/* Add Partner Modal */}
        {showAddModal && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <Card className="w-full max-w-2xl max-h-[80vh] overflow-y-auto">
              <CardHeader>
                <CardTitle>Add Service Partner</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Input
                    placeholder="Search partners by name or location..."
                    value={searchTerm}
                    onChange={(e) => {
                      setSearchTerm(e.target.value);
                      loadAvailablePartners(e.target.value);
                    }}
                  />
                </div>

                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {availablePartners.map((org) => (
                    <div
                      key={org.id}
                      className="p-4 border rounded hover:bg-gray-50 cursor-pointer"
                      onClick={() => handleAssociatePartner(org.id)}
                    >
                      <h4 className="font-semibold">{org.name}</h4>
                      <p className="text-sm text-gray-600">
                        {getOrgTypeLabel(org.org_type)} • {org.engineers_count} engineers
                        {org.location && ` • ${org.location}`}
                      </p>
                    </div>
                  ))}
                  {availablePartners.length === 0 && (
                    <p className="text-center text-gray-500 py-8">
                      No available partners found. Try adjusting your search.
                    </p>
                  )}
                </div>

                <div className="flex justify-end gap-2">
                  <Button variant="outline" onClick={() => setShowAddModal(false)}>
                    Cancel
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
