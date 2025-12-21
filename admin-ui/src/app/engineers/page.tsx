'use client';

import { useState, useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, Search, Plus, Upload, Download, UserCheck, UserX, Users, Loader2, Filter } from 'lucide-react';
import apiClient from '@/lib/api/client';

interface Engineer {
  id: string;
  name: string;
  phone: string;
  email: string;
  skills: string[];
  engineer_level: number | string; // Can be number (1,2,3) or string ("L1","L2","L3")
  home_region: string;
  metadata?: any;
  created_at: string;
  updated_at: string;
}

export default function EngineersListPage() {
  const router = useRouter();
  const [engineers, setEngineers] = useState<Engineer[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterSkill, setFilterSkill] = useState<string>('all');
  const [filterRegion, setFilterRegion] = useState<string>('all');
  const [filterLevel, setFilterLevel] = useState<string>('all');
  const [manufacturerName, setManufacturerName] = useState<string>('');
  
  // Get manufacturer filter from URL params
  const searchParams = typeof window !== 'undefined' ? new URLSearchParams(window.location.search) : null;
  const manufacturerFilter = searchParams?.get('manufacturer') || '';
  
  // Fetch manufacturer name if filtering
  useEffect(() => {
    if (manufacturerFilter) {
      const fetchManufacturerName = async () => {
        try {
          const response = await apiClient.get(`/v1/organizations/${manufacturerFilter}`);
          setManufacturerName(response.data.name || '');
        } catch (err) {
          console.error('Failed to fetch manufacturer name:', err);
        }
      };
      fetchManufacturerName();
    }
  }, [manufacturerFilter]);

  // Fetch engineers from API
  useEffect(() => {
    loadEngineers();
  }, [manufacturerFilter]);

  const loadEngineers = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Build URL with manufacturer filter if present
      let url = '/v1/engineers?limit=100';
      if (manufacturerFilter) {
        url += `&organization_id=${manufacturerFilter}`;
      }
      
      const response = await apiClient.get(url);
      // Backend returns 'engineers' array, not 'items'
      const engineersData = response.data.engineers || response.data.items || [];
      
      // Deduplicate engineers by ID (backend may return duplicates for multi-org assignments)
      const uniqueEngineers = Array.from(
        new Map(engineersData.map((eng: Engineer) => [eng.id, eng])).values()
      ) as Engineer[];
      
      setEngineers(uniqueEngineers);
    } catch (err) {
      console.error('Failed to load engineers:', err);
      setError('Failed to load engineers. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  // Get unique skills for filter
  const allSkills = useMemo(() => {
    const skills = new Set<string>();
    engineers.forEach(eng => {
      if (eng.skills) {
        eng.skills.forEach(skill => skills.add(skill));
      }
    });
    return Array.from(skills).sort();
  }, [engineers]);

  // Get unique regions for filter
  const allRegions = useMemo(() => {
    const regions = new Set<string>();
    engineers.forEach(eng => {
      if (eng.home_region) regions.add(eng.home_region);
    });
    return Array.from(regions).sort();
  }, [engineers]);

  // Filter engineers based on search and filters
  const filteredEngineers = useMemo(() => {
    return engineers.filter(engineer => {
      // Search filter
      const matchesSearch = searchQuery === '' || 
        engineer.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.phone?.includes(searchQuery) ||
        engineer.email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.home_region?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        engineer.skills?.some(s => s.toLowerCase().includes(searchQuery.toLowerCase()));

      // Skill filter
      const matchesSkill = filterSkill === 'all' || 
        engineer.skills?.includes(filterSkill);

      // Region filter
      const matchesRegion = filterRegion === 'all' || 
        engineer.home_region === filterRegion;

      // Level filter
      const matchesLevel = filterLevel === 'all' || 
        String(engineer.engineer_level) === filterLevel;

      return matchesSearch && matchesSkill && matchesRegion && matchesLevel;
    });
  }, [engineers, searchQuery, filterSkill, filterRegion, filterLevel]);

  // Statistics
  const stats = useMemo(() => {
    const byLevel = engineers.reduce((acc, eng) => {
      // Handle both integer (1, 2, 3) and string ("L1", "L2", "L3") formats
      let level: number;
      if (typeof eng.engineer_level === 'string') {
        if (eng.engineer_level.startsWith('L')) {
          level = parseInt(eng.engineer_level.substring(1), 10);
        } else {
          level = parseInt(eng.engineer_level, 10);
        }
      } else {
        level = eng.engineer_level;
      }
      const levelNum = level || 1;
      acc[levelNum] = (acc[levelNum] || 0) + 1;
      return acc;
    }, {} as Record<number, number>);
    
    return {
      total: engineers.length,
      filtered: filteredEngineers.length,
      byLevel,
    };
  }, [engineers, filteredEngineers]);

  const getLevelBadge = (level: number | string) => {
    // Convert string levels to numbers
    let numLevel: number;
    if (typeof level === 'string') {
      if (level.startsWith('L')) {
        numLevel = parseInt(level.substring(1), 10);
      } else {
        numLevel = parseInt(level, 10);
      }
    } else {
      numLevel = level;
    }
    
    const colors = {
      1: 'bg-blue-100 text-blue-800',
      2: 'bg-green-100 text-green-800',
      3: 'bg-purple-100 text-purple-800',
    };
    const labels = {
      1: 'Junior',
      2: 'Senior',
      3: 'Expert',
    };
    return {
      color: colors[numLevel as keyof typeof colors] || 'bg-gray-100 text-gray-800',
      label: labels[numLevel as keyof typeof labels] || `Level ${numLevel}`,
    };
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
          <p className="text-gray-600">Loading engineers...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-red-600">Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-700 mb-4">{error}</p>
            <Button onClick={loadEngineers} className="w-full">
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
        {/* Back button if filtering by manufacturer */}
        {manufacturerFilter && (
          <Button
            variant="ghost"
            onClick={() => router.push(`/manufacturers/${manufacturerFilter}/dashboard`)}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to {manufacturerName || 'Manufacturer'} Dashboard
          </Button>
        )}
        
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-2">
              <Users className="h-8 w-8 text-blue-600" />
              Engineers Management
              {manufacturerName && <span className="text-blue-600">- {manufacturerName}</span>}
            </h1>
            <p className="text-gray-600 mt-1">
              {manufacturerFilter 
                ? `Showing ${filteredEngineers.length} engineers for ${manufacturerName || 'selected manufacturer'}`
                : 'Manage field engineers and service technicians'
              }
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => router.push('/engineers/import')}
            >
              <Upload className="h-4 w-4 mr-2" />
              Import CSV
            </Button>
            <Button
              onClick={() => router.push('/engineers/add')}
            >
              <Plus className="h-4 w-4 mr-2" />
              Add Engineer
            </Button>
          </div>
        </div>

        {/* Statistics Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <Card>
            <CardContent className="pt-6">
              <div className="text-center">
                <p className="text-sm text-gray-600">Total Engineers</p>
                <p className="text-3xl font-bold text-gray-900">{stats.total}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-center">
                <p className="text-sm text-gray-600">Junior (Level 1)</p>
                <p className="text-3xl font-bold text-blue-600">{stats.byLevel[1] || 0}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-center">
                <p className="text-sm text-gray-600">Senior (Level 2)</p>
                <p className="text-3xl font-bold text-green-600">{stats.byLevel[2] || 0}</p>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="pt-6">
              <div className="text-center">
                <p className="text-sm text-gray-600">Expert (Level 3)</p>
                <p className="text-3xl font-bold text-purple-600">{stats.byLevel[3] || 0}</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Search and Filters */}
        <Card className="mb-6">
          <CardContent className="pt-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              {/* Search */}
              <div className="md:col-span-2">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                  <Input
                    type="text"
                    placeholder="Search by name, email, phone, region, or skills..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </div>

              {/* Skill Filter */}
              <div>
                <select
                  value={filterSkill}
                  onChange={(e) => setFilterSkill(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">All Skills</option>
                  {allSkills.map(skill => (
                    <option key={skill} value={skill}>{skill}</option>
                  ))}
                </select>
              </div>

              {/* Region Filter */}
              <div className="flex gap-2">
                <select
                  value={filterRegion}
                  onChange={(e) => setFilterRegion(e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">All Regions</option>
                  {allRegions.map(region => (
                    <option key={region} value={region}>{region}</option>
                  ))}
                </select>

                {/* Level Filter */}
                <select
                  value={filterLevel}
                  onChange={(e) => setFilterLevel(e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">All Levels</option>
                  <option value="1">Junior</option>
                  <option value="2">Senior</option>
                  <option value="3">Expert</option>
                </select>
              </div>
            </div>

            {/* Active Filters */}
            {(searchQuery || filterSkill !== 'all' || filterRegion !== 'all' || filterLevel !== 'all') && (
              <div className="mt-4 flex items-center gap-2">
                <span className="text-sm text-gray-600">Active filters:</span>
                <div className="flex gap-2 flex-wrap">
                  {searchQuery && (
                    <Badge variant="secondary" className="text-xs">
                      Search: {searchQuery}
                      <button onClick={() => setSearchQuery('')} className="ml-1 hover:text-red-600">×</button>
                    </Badge>
                  )}
                  {filterSkill !== 'all' && (
                    <Badge variant="secondary" className="text-xs">
                      Skill: {filterSkill}
                      <button onClick={() => setFilterSkill('all')} className="ml-1 hover:text-red-600">×</button>
                    </Badge>
                  )}
                  {filterRegion !== 'all' && (
                    <Badge variant="secondary" className="text-xs">
                      Region: {filterRegion}
                      <button onClick={() => setFilterRegion('all')} className="ml-1 hover:text-red-600">×</button>
                    </Badge>
                  )}
                  {filterLevel !== 'all' && (
                    <Badge variant="secondary" className="text-xs">
                      Level: {filterLevel}
                      <button onClick={() => setFilterLevel('all')} className="ml-1 hover:text-red-600">×</button>
                    </Badge>
                  )}
                  <button
                    onClick={() => {
                      setSearchQuery('');
                      setFilterSkill('all');
                      setFilterRegion('all');
                      setFilterLevel('all');
                    }}
                    className="text-xs text-blue-600 hover:text-blue-800 underline"
                  >
                    Clear all
                  </button>
                </div>
              </div>
            )}

            <div className="mt-4 text-sm text-gray-600">
              Showing {filteredEngineers.length} of {stats.total} engineers
            </div>
          </CardContent>
        </Card>

        {/* Engineers Table */}
        <Card>
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Engineer
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Contact
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Skills
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Level
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Region
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {filteredEngineers.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                        <Users className="h-12 w-12 mx-auto mb-3 text-gray-300" />
                        <p>No engineers found</p>
                        <p className="text-sm mt-1">Try adjusting your search or filters</p>
                      </td>
                    </tr>
                  ) : (
                    filteredEngineers.map((engineer) => (
                      <tr
                        key={engineer.id}
                        className="hover:bg-gray-50 cursor-pointer transition-colors"
                        onClick={() => router.push(`/engineers/${engineer.id}`)}
                      >
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div>
                            <div className="text-sm font-medium text-gray-900">
                              {engineer.name}
                            </div>
                            <div className="text-sm text-gray-500">
                              ID: {engineer.id.slice(0, 8)}...
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{engineer.email || 'N/A'}</div>
                          <div className="text-sm text-gray-500">{engineer.phone || 'N/A'}</div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex flex-wrap gap-1">
                            {engineer.skills && engineer.skills.length > 0 ? (
                              engineer.skills.slice(0, 3).map((skill, idx) => (
                                <Badge key={idx} variant="outline" className="text-xs">
                                  {skill}
                                </Badge>
                              ))
                            ) : (
                              <span className="text-xs text-gray-400">No skills listed</span>
                            )}
                            {engineer.skills && engineer.skills.length > 3 && (
                              <Badge variant="outline" className="text-xs">
                                +{engineer.skills.length - 3}
                              </Badge>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <Badge className={getLevelBadge(engineer.engineer_level).color}>
                            {getLevelBadge(engineer.engineer_level).label}
                          </Badge>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {engineer.home_region || 'N/A'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={(e) => {
                              e.stopPropagation();
                              router.push(`/engineers/${engineer.id}`);
                            }}
                          >
                            View Details
                          </Button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
