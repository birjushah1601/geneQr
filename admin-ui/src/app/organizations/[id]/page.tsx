'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { organizationsApi, Organization, Facility, OrgRelationship } from '@/lib/api/organizations';
import { Building2, MapPin, Users, ArrowLeft, Loader2, Building, Phone, Mail, Globe, ChevronRight } from 'lucide-react';
import Link from 'next/link';

export default function OrganizationDetailPage() {
  const params = useParams();
  const router = useRouter();
  const orgId = params.id as string;

  const [organization, setOrganization] = useState<Organization | null>(null);
  const [facilities, setFacilities] = useState<Facility[]>([]);
  const [relationships, setRelationships] = useState<OrgRelationship[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'overview' | 'facilities' | 'relationships'>('overview');

  useEffect(() => {
    if (orgId) {
      loadOrganizationDetails();
    }
  }, [orgId]);

  const loadOrganizationDetails = async () => {
    try {
      setLoading(true);
      setError(null);

      const [orgData, facilitiesData, relationshipsData] = await Promise.all([
        organizationsApi.get(orgId),
        organizationsApi.listFacilities(orgId),
        organizationsApi.listRelationships(orgId),
      ]);

      setOrganization(orgData);
      setFacilities(facilitiesData);
      setRelationships(relationshipsData);
    } catch (err) {
      console.error('Failed to load organization:', err);
      setError('Failed to load organization details. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const getOrgTypeColor = (type: string) => {
    const colors = {
      manufacturer: 'bg-blue-100 text-blue-800 border-blue-200',
      distributor: 'bg-purple-100 text-purple-800 border-purple-200',
      dealer: 'bg-green-100 text-green-800 border-green-200',
      hospital: 'bg-red-100 text-red-800 border-red-200',
      service_provider: 'bg-yellow-100 text-yellow-800 border-yellow-200',
      other: 'bg-gray-100 text-gray-800 border-gray-200',
    };
    return colors[type as keyof typeof colors] || colors.other;
  };

  const getOrgTypeIcon = (type: string) => {
    switch (type) {
      case 'manufacturer': return '🏭';
      case 'distributor': return '📦';
      case 'dealer': return '🏪';
      case 'hospital': return '🏥';
      case 'service_provider': return '🔧';
      default: return '🏢';
    }
  };

  if (loading) {
    return (
      <div className="p-8">
        <div className="flex items-center justify-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
          <span className="ml-2 text-gray-600">Loading organization...</span>
        </div>
      </div>
    );
  }

  if (error || !organization) {
    return (
      <div className="p-8">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-800">
          {error || 'Organization not found'}
        </div>
        <button
          onClick={() => router.push('/organizations')}
          className="mt-4 text-blue-600 hover:text-blue-800"
        >
          ← Back to Organizations
        </button>
      </div>
    );
  }

  return (
    <div className="p-8">
      {/* Back Button */}
      <Link
        href="/organizations"
        className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="w-4 h-4 mr-2" />
        Back to Organizations
      </Link>

      {/* Header */}
      <div className="bg-white rounded-lg shadow-lg p-8 mb-6">
        <div className="flex items-start justify-between">
          <div className="flex items-start space-x-4">
            <div className="text-5xl">{getOrgTypeIcon(organization.org_type)}</div>
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">{organization.name}</h1>
              <div className="flex items-center space-x-4">
                <span className={`inline-block px-3 py-1 rounded-full text-sm font-medium border ${getOrgTypeColor(organization.org_type)}`}>
                  {organization.org_type.replace('_', ' ')}
                </span>
                <div className="flex items-center space-x-2">
                  <div className={`w-2 h-2 rounded-full ${organization.status === 'active' ? 'bg-green-500' : 'bg-gray-400'}`} />
                  <span className="text-sm text-gray-600 capitalize">{organization.status}</span>
                </div>
              </div>
            </div>
          </div>

          <div className="text-right">
            <p className="text-xs text-gray-400 mb-1">Organization ID</p>
            <p className="text-sm font-mono text-gray-600">{organization.id}</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="bg-white rounded-lg shadow mb-6">
        <div className="border-b border-gray-200">
          <nav className="flex space-x-8 px-6" aria-label="Tabs">
            <button
              onClick={() => setActiveTab('overview')}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'overview'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Overview
            </button>
            <button
              onClick={() => setActiveTab('facilities')}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'facilities'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Facilities ({facilities.length})
            </button>
            <button
              onClick={() => setActiveTab('relationships')}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'relationships'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Relationships ({relationships.length})
            </button>
          </nav>
        </div>

        {/* Tab Content */}
        <div className="p-6">
          {/* Overview Tab */}
          {activeTab === 'overview' && (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-blue-50 p-6 rounded-lg">
                  <div className="flex items-center space-x-3 mb-2">
                    <Building className="w-5 h-5 text-blue-600" />
                    <h3 className="font-semibold text-gray-900">Facilities</h3>
                  </div>
                  <p className="text-3xl font-bold text-blue-600">{facilities.length}</p>
                  <p className="text-sm text-gray-600 mt-1">Total locations</p>
                </div>

                <div className="bg-purple-50 p-6 rounded-lg">
                  <div className="flex items-center space-x-3 mb-2">
                    <Users className="w-5 h-5 text-purple-600" />
                    <h3 className="font-semibold text-gray-900">Relationships</h3>
                  </div>
                  <p className="text-3xl font-bold text-purple-600">{relationships.length}</p>
                  <p className="text-sm text-gray-600 mt-1">Business connections</p>
                </div>

                <div className="bg-green-50 p-6 rounded-lg">
                  <div className="flex items-center space-x-3 mb-2">
                    <Building2 className="w-5 h-5 text-green-600" />
                    <h3 className="font-semibold text-gray-900">Type</h3>
                  </div>
                  <p className="text-lg font-bold text-green-600 capitalize">{organization.org_type.replace('_', ' ')}</p>
                  <p className="text-sm text-gray-600 mt-1">Organization category</p>
                </div>
              </div>

              {/* Metadata */}
              {organization.metadata && Object.keys(organization.metadata).length > 0 && (
                <div className="bg-gray-50 p-6 rounded-lg">
                  <h3 className="font-semibold text-gray-900 mb-4">Additional Information</h3>
                  <pre className="text-sm text-gray-600 overflow-auto">
                    {JSON.stringify(organization.metadata, null, 2)}
                  </pre>
                </div>
              )}
            </div>
          )}

          {/* Facilities Tab */}
          {activeTab === 'facilities' && (
            <div>
              {facilities.length === 0 ? (
                <div className="text-center py-12">
                  <MapPin className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                  <p className="text-gray-600">No facilities registered yet</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {facilities.map((facility) => (
                    <div key={facility.id} className="border border-gray-200 rounded-lg p-6 hover:shadow-md transition-shadow">
                      <div className="flex items-start justify-between mb-4">
                        <div>
                          <h3 className="font-semibold text-gray-900 text-lg mb-1">
                            {facility.facility_name}
                          </h3>
                          <p className="text-sm text-gray-600">{facility.facility_code}</p>
                        </div>
                        <span className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                          {facility.facility_type}
                        </span>
                      </div>

                      {facility.address && (
                        <div className="flex items-start space-x-2 text-sm text-gray-600">
                          <MapPin className="w-4 h-4 mt-0.5 flex-shrink-0" />
                          <div>
                            {typeof facility.address === 'string' ? (
                              <p>{facility.address}</p>
                            ) : (
                              <div className="space-y-1">
                                {facility.address.street && <p>{facility.address.street}</p>}
                                {facility.address.city && <p>{facility.address.city}, {facility.address.state} {facility.address.pincode}</p>}
                                {facility.address.country && <p>{facility.address.country}</p>}
                              </div>
                            )}
                          </div>
                        </div>
                      )}

                      <div className="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <div className={`w-2 h-2 rounded-full ${facility.status === 'active' ? 'bg-green-500' : 'bg-gray-400'}`} />
                          <span className="text-xs text-gray-600 capitalize">{facility.status}</span>
                        </div>
                        <span className="text-xs text-gray-400">ID: {facility.id.slice(0, 8)}...</span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Relationships Tab */}
          {activeTab === 'relationships' && (
            <div>
              {relationships.length === 0 ? (
                <div className="text-center py-12">
                  <Users className="w-12 h-12 text-gray-400 mx-auto mb-2" />
                  <p className="text-gray-600">No business relationships established yet</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {relationships.map((rel) => (
                    <div key={rel.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
                      <div className="flex items-center space-x-4">
                        <ChevronRight className="w-5 h-5 text-gray-400" />
                        <div>
                          <p className="font-medium text-gray-900">{rel.rel_type.replace('_', ' ')}</p>
                          <p className="text-sm text-gray-600">
                            Parent: {rel.parent_org_id.slice(0, 8)}... → Child: {rel.child_org_id.slice(0, 8)}...
                          </p>
                        </div>
                      </div>
                      <span className="text-xs text-gray-400">ID: {rel.id.slice(0, 8)}...</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
