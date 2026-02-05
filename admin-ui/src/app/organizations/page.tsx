'use client';

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import DashboardLayout from '@/components/DashboardLayout';

import { organizationsApi, Organization } from '@/lib/api/organizations';
import {
  Building2,
  Filter,
  Search,
  Loader2,
  ExternalLink,
} from 'lucide-react';

export default function OrganizationsPage() {
  const searchParams = useSearchParams();

  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterType, setFilterType] = useState('all');
  const [filterStatus, setFilterStatus] = useState('all');

  /* ----------------------------------
     Init filters from URL
  ---------------------------------- */
  useEffect(() => {
    const type = searchParams.get('type');
    const status = searchParams.get('status');
    if (type) setFilterType(type);
    if (status) setFilterStatus(status);
  }, [searchParams]);

  /* ----------------------------------
     Load data
  ---------------------------------- */
  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadOrganizations = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await organizationsApi.list({ page_size: 100 });
      setOrganizations(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
      setError('Failed to load organizations.');
    } finally {
      setLoading(false);
    }
  };

  /* ----------------------------------
     Derived data
  ---------------------------------- */
  const filteredOrganizations = organizations.filter((org) => {
    const matchesSearch = org.name
      .toLowerCase()
      .includes(searchTerm.toLowerCase());
    const matchesType = filterType === 'all' || org.org_type === filterType;
    const matchesStatus =
      filterStatus === 'all' || org.status === filterStatus;
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
    const colors: Record<string, string> = {
      manufacturer: 'bg-blue-100 text-blue-800',
      channel_partner: 'bg-purple-100 text-purple-800',
      sub_dealer: 'bg-green-100 text-green-800',
      hospital: 'bg-red-100 text-red-800',
    };
    return colors[type] || 'bg-gray-100 text-gray-800';
  };

  const getOrgTypeIcon = (type: string) => {
    switch (type) {
      case 'manufacturer': return '🏭';
      case 'channel_partner': return '🚚';
      case 'sub_dealer': return '🛍️';
      case 'hospital': return '🏥';
      default: return '🏢';
    }
  };

  return (
    <DashboardLayout>
      <div>
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Organizations
          </h1>
          <p className="text-gray-600">
            Manage manufacturers, channel partners, sub-dealers, and hospitals
          </p>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
          <StatCard label="Total" value={stats.total} icon={<Building2 />} />
          <StatCard label="Manufacturers" value={stats.manufacturers} emoji="🏭" />
          <StatCard label="Channel Partners" value={stats.channelPartners} emoji="🚚" />
          <StatCard label="Sub-dealers" value={stats.subDealers} emoji="🛍️" />
          <StatCard label="Hospitals" value={stats.hospitals} emoji="🏥" />
        </div>

        {/* Filters */}
        <div className="bg-white p-6 rounded-lg shadow mb-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search organizations..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Types</option>
              <option value="manufacturer">Manufacturers</option>
              <option value="channel_partner">Channel Partners</option>
              <option value="sub_dealer">Sub-dealers</option>
              <option value="hospital">Hospitals</option>
            </select>

            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Status</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>
        </div>

        {/* Content */}
        {loading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
            <span className="ml-2 text-gray-600">Loading organizations…</span>
          </div>
        ) : error ? (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-800">
            {error}
          </div>
        ) : filteredOrganizations.length === 0 ? (
          <div className="bg-gray-50 border border-gray-200 rounded-lg p-8 text-center">
            <Building2 className="w-12 h-12 text-gray-400 mx-auto mb-2" />
            <p className="text-gray-600">No organizations found</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredOrganizations.map((org) => (
              <Link
                key={org.id}
                href={`/organizations/${org.id}`}
                className="block bg-white rounded-lg shadow hover:shadow-lg transition-shadow border"
              >
                <div className="p-6">
                  <div className="flex justify-between mb-4">
                    <div className="flex gap-3">
                      <span className="text-3xl">{getOrgTypeIcon(org.org_type)}</span>
                      <div>
                        <h3 className="font-semibold text-lg">{org.name}</h3>
                        <span
                          className={`inline-block px-2 py-1 rounded-full text-xs ${getOrgTypeColor(org.org_type)}`}
                        >
                          {org.org_type.replace('_', ' ')}
                        </span>
                      </div>
                    </div>
                    <ExternalLink className="w-5 h-5 text-gray-400" />
                  </div>

                  <div className="flex justify-between pt-4 border-t text-sm text-gray-600">
                    <span className="capitalize">{org.status}</span>
                    <span>ID: {org.id.slice(0, 8)}…</span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}

        {!loading && !error && (
          <div className="mt-6 text-center text-sm text-gray-600">
            Showing {filteredOrganizations.length} of {organizations.length}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

/* ----------------------------------
   Helper component
---------------------------------- */
function StatCard({
  label,
  value,
  emoji,
  icon,
}: {
  label: string;
  value: number;
  emoji?: string;
  icon?: React.ReactNode;
}) {
  return (
    <div className="bg-white p-6 rounded-lg shadow border-l-4 border-blue-500">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-500 text-sm">{label}</p>
          <p className="text-2xl font-bold">{value}</p>
        </div>
        {emoji ? <span className="text-2xl">{emoji}</span> : icon}
      </div>
    </div>
  );
}
