'use client';

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { organizationsApi, Organization } from '@/lib/api/organizations';
import { Building2, MapPin, Users, Filter, Search, Loader2, ExternalLink } from 'lucide-react';
import Link from 'next/link';
import DashboardLayout from '@/components/DashboardLayout';

export default function OrganizationsPage() {
  const searchParams = useSearchParams();
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState<string>('all');
  const [filterStatus, setFilterStatus] = useState<string>('all');

  // Initialize filters from URL query parameters
  useEffect(() => {
    const typeParam = searchParams.get('type');
    const statusParam = searchParams.get('status');
    
    if (typeParam) {
      setFilterType(typeParam);
    }
    if (statusParam) {
      setFilterStatus(statusParam);
    }
  }, [searchParams]);

  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadOrganizations = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await organizationsApi.list({ limit: 100 });
      setOrganizations(data);
    } catch (err) {
      console.error('Failed to load organizations:', err);
      setError('Failed to load organizations. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const filteredOrganizations = organizations.filter((org) => {
    const matchesSearch = org.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = filterType === 'all' || org.org_type === filterType;
    const matchesStatus = filterStatus === 'all' || org.status === filterStatus;
    return matchesSearch && matchesType && matchesStatus;
  });

  const stats = {
    total: organizations.length,
    manufacturers: organizations.filter(o => o.org_type === 'manufacturer').length,
    channelPartners: organizations.filter(o => o.org_type === 'channel_partner').length,
    subDealers: organizations.filter(o => o.org_type === 'sub_dealer').length,
    hospitals: organizations.filter(o => o.org_type === 'hospital').length,
  };

  const getOrgTypeColor = (type: string) => {
    const colors = {
      manufacturer: 'bg-blue-100 text-blue-800',
      'channel_partner': 'bg-purple-100 text-purple-800',
      'sub_dealer': 'bg-green-100 text-green-800',
      hospital: 'bg-red-100 text-red-800',
      service_provider: 'bg-yellow-100 text-yellow-800',
      other: 'bg-gray-100 text-gray-800',
    };
    return colors[type as keyof typeof colors] || colors.other;
  };

  const getOrgTypeIcon = (type: string) => {
    switch (type) {
      case 'manufacturer': return 'Ã°Å¸ÂÂ­';
      case 'channel_partner': return 'Ã°Å¸â€œÂ¦';
      case 'sub_dealer': return 'Ã°Å¸ÂÂª';
      case 'hospital': return 'Ã°Å¸ÂÂ¥';
      case 'service_provider': return 'Ã°Å¸â€Â§';
      default: return 'Ã°Å¸ÂÂ¢';
    }
  };

  return (
    <DashboardLayout>
      <div>
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Organizations</h1>
        <p className="text-gray-600">
          Manage manufacturers, Channel Partners, Sub-Dealers, and hospitals
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-5 gap-4 mb-6">
        <div className="bg-white p-6 rounded-lg shadow border-l-4 border-blue-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">Total</p>
              <p className="text-2xl font-bold text-gray-900">{stats.total}</p>
            </div>
            <Building2 className="w-8 h-8 text-blue-500" />
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow border-l-4 border-purple-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">Manufacturers</p>
              <p className="text-2xl font-bold text-gray-900">{stats.manufacturers}</p>
            </div>
            <span className="text-2xl">Ã°Å¸ÂÂ­</span>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow border-l-4 border-green-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">Channel Partners</p>
              <p className="text-2xl font-bold text-gray-900">{stats.channelPartners}</p>
            </div>
            <span className="text-2xl">Ã°Å¸â€œÂ¦</span>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow border-l-4 border-amber-500">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-gray-500 text-sm">subDealers:grid-cols-3 gap-4">
          {/* Search */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search organizations..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {/* Type Filter */}
          <div className="relative">
            <Filter className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none"
            >
              <option value="all">All Types</option>
              <option value="manufacturer">Manufacturers</option>
              <option value="Channel Partner">Channel Partners</option>
              <option value="Sub-sub_SUB_DEALER">subDealers:ring-2 focus:ring-blue-500 focus:border-transparent appearance-none"
            >
              <option value="all">All Status</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>
        </div>
      </div>

      {/* Content */}
      {loading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
          <span className="ml-2 text-gray-600">Loading organizations...</span>
        </div>
      ) : error ? (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-800">
          {error}
        </div>
      ) : filteredOrganizations.length === 0 ? (
        <div className="bg-gray-50 border border-gray-200 rounded-lg p-8 text-center">
          <Building2 className="w-12 h-12 text-gray-400 mx-auto mb-2" />
          <p className="text-gray-600">No organizations found matching your filters</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredOrganizations.map((org) => (
            <Link
              key={org.id}
              href={`/organizations/${org.id}`}
              className="block bg-white rounded-lg shadow hover:shadow-lg transition-shadow border border-gray-200 overflow-hidden"
            >
              <div className="p-6">
                {/* Header */}
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-start space-x-3">
                    <span className="text-3xl">{getOrgTypeIcon(org.org_type)}</span>
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-900 text-lg mb-1 line-clamp-2">
                        {org.name}
                      </h3>
                      <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${getOrgTypeColor(org.org_type)}`}>
                        {org.org_type.replace('_', ' ')}
                      </span>
                    </div>
                  </div>
                  <ExternalLink className="w-5 h-5 text-gray-400" />
                </div>

                {/* Status */}
                <div className="flex items-center justify-between pt-4 border-t border-gray-100">
                  <div className="flex items-center space-x-2">
                    <div className={`w-2 h-2 rounded-full ${org.status === 'active' ? 'bg-green-500' : 'bg-gray-400'}`} />
                    <span className="text-sm text-gray-600 capitalize">{org.status}</span>
                  </div>
                  <span className="text-xs text-gray-400">ID: {org.id.slice(0, 8)}...</span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* Results Count */}
      {!loading && !error && (
        <div className="mt-6 text-center text-sm text-gray-600">
          Showing {filteredOrganizations.length} of {organizations.length} organizations
        </div>
      )}
      </div>
    </DashboardLayout>
  );
}
