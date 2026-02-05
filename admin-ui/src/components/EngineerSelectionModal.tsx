'use client';

import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { 
  User, 
  Building2, 
  Award, 
  MapPin, 
  Phone, 
  Mail,
  CheckCircle2,
  Loader2,
  AlertCircle,
  Star,
  Truck,
  Users as UsersIcon
} from 'lucide-react';
import apiClient from '@/lib/api/client';
import { partnersApi } from '@/lib/api/partners';
import { useAuth } from '@/contexts/AuthContext';

interface Engineer {
  engineer_id: string;
  engineer_name: string;
  organization_id: string;
  organization_name: string;
  engineer_level: string;
  match_score?: number;
  manufacturer_certified?: boolean;
  equipment_types?: string[];
  location?: string;
  phone?: string;
  email?: string;
  availability?: string;
  organization_type?: string; // NEW: for badge display
  category?: string; // NEW: 'Manufacturer', 'Channel Partner - [Name]', etc.
}

interface EngineerSelectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  ticketId: string;
  equipmentName?: string;
  onAssignmentSuccess?: () => void;
}

export default function EngineerSelectionModal({
  isOpen,
  onClose,
  ticketId,
  equipmentName,
  onAssignmentSuccess
}: EngineerSelectionModalProps) {
  const { organizationContext } = useAuth();
  const [engineers, setEngineers] = useState<Engineer[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isAssigning, setIsAssigning] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedEngineer, setSelectedEngineer] = useState<string | null>(null);
  const [activeCategory, setActiveCategory] = useState<string>('all');

  // Fetch suggested engineers when modal opens
  useEffect(() => {
    if (isOpen && ticketId) {
      fetchSuggestedEngineers();
    }
  }, [isOpen, ticketId]);

  const fetchSuggestedEngineers = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Get ticket details first to extract manufacturer and equipment info
      const ticketResponse = await apiClient.get(`/v1/tickets/${ticketId}`);
      const ticket = ticketResponse.data;
      
      // Fetch engineers with partners included
      const engineersResponse = await apiClient.get('/v1/engineers?limit=100&include_partners=true');
      const allEngineers = engineersResponse.data.engineers || [];
      
      // Get current user's organization
      const currentOrgId = organizationContext?.organization_id || ticket.manufacturer_id;
      
      // Transform to match Engineer interface
      const transformedEngineers = allEngineers.map((eng: any) => {
        const isOwnEngineer = eng.organization_id === currentOrgId;
        const category = isOwnEngineer ? 'Manufacturer' : `Partner - ${eng.organization_name}`;
        
        return {
          engineer_id: eng.id,
          engineer_name: eng.name,
          organization_id: eng.organization_id,
          organization_name: eng.organization_name,
          organization_type: isOwnEngineer ? 'manufacturer' : 'channel_partner',
          engineer_level: eng.engineer_level || 'L2',
          phone: eng.phone,
          email: eng.email,
          category: category,
          manufacturer_certified: !isOwnEngineer,
          match_score: null, // No AI scoring for now
          equipment_types: [],
          location: null,
        };
      });
      
      // Sort: Own engineers first, then partners, then by name
      transformedEngineers.sort((a, b) => {
        const aIsOwn = a.organization_id === currentOrgId ? 0 : 1;
        const bIsOwn = b.organization_id === currentOrgId ? 0 : 1;
        if (aIsOwn !== bIsOwn) return aIsOwn - bIsOwn;
        return a.engineer_name.localeCompare(b.engineer_name);
      });
      
      setEngineers(transformedEngineers);
    } catch (err: any) {
      console.error('Failed to fetch engineers:', err);
      setError(err.response?.data?.error || 'Failed to load engineers');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAssignEngineer = async (engineerId: string) => {
    setIsAssigning(true);
    setError(null);
    
    try {
      // Find the engineer to get their name
      const engineer = engineers.find(e => e.engineer_id === engineerId);
      if (!engineer) {
        throw new Error('Engineer not found');
      }
      
      await apiClient.post(`/v1/tickets/${ticketId}/assign`, {
        engineer_id: engineerId,
        engineer_name: engineer.engineer_name, // Include engineer name
        assignment_tier: 'tier_1' // Can be determined based on organization type
      });
      
      // Success!
      if (onAssignmentSuccess) {
        onAssignmentSuccess();
      }
      onClose();
    } catch (err: any) {
      console.error('Failed to assign engineer:', err);
      setError(err.response?.data?.error || 'Failed to assign engineer');
    } finally {
      setIsAssigning(false);
    }
  };

  const getLevelColor = (level: string) => {
    switch (level) {
      case 'L3':
        return 'bg-purple-100 text-purple-800 border-purple-300';
      case 'L2':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      case 'L1':
        return 'bg-green-100 text-green-800 border-green-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  const getLevelLabel = (level: string) => {
    switch (level) {
      case 'L3':
        return 'Senior Engineer';
      case 'L2':
        return 'Intermediate Engineer';
      case 'L1':
        return 'Junior Engineer';
      default:
        return level;
    }
  };

  const getMatchScoreColor = (score?: number) => {
    if (!score) return 'text-gray-500';
    if (score >= 90) return 'text-green-600';
    if (score >= 75) return 'text-blue-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-gray-600';
  };

  const getOrgTypeBadge = (category?: string, orgType?: string) => {
    if (!category) return null;
    
    if (category === 'Manufacturer' || orgType === 'manufacturer') {
      return (
        <Badge className="bg-blue-100 text-blue-800 border-blue-300">
          <Building2 className="h-3 w-3 mr-1" />
          Manufacturer
        </Badge>
      );
    }
    
    if (category.includes('Channel Partner') || orgType === 'channel_partner') {
      return (
        <Badge className="bg-orange-600 text-white border-orange-700">
          <Truck className="h-3 w-3 mr-1" />
          {category.replace('Channel Partner - ', 'Partner: ')}
        </Badge>
      );
    }
    
    if (category.includes('Sub-Dealer') || orgType === 'sub_dealer') {
      return (
        <Badge className="bg-purple-700 text-white border-purple-800">
          <UsersIcon className="h-3 w-3 mr-1" />
          {category.replace('Sub-Dealer - ', 'Dealer: ')}
        </Badge>
      );
    }
    
    if (category.includes('Partner')) {
      return (
        <Badge className="bg-orange-600 text-white border-orange-700">
          <Truck className="h-3 w-3 mr-1" />
          Partner
        </Badge>
      );
    }
    
    return null;
  };

  // Filter engineers by category
  const filteredEngineers = engineers.filter(eng => {
    if (activeCategory === 'all') return true;
    if (activeCategory === 'senior') return eng.engineer_level === 'L3';
    if (activeCategory === 'certified') return eng.manufacturer_certified;
    if (activeCategory === 'high-match') return (eng.match_score || 0) >= 80;
    if (activeCategory === 'partner-network') {
      return eng.category !== 'Manufacturer' || 
             eng.organization_type === 'channel_partner' || 
             eng.organization_type === 'sub_dealer';
    }
    return true;
  });

  const categories = [
    { id: 'all', label: 'Best Overall Match', count: engineers.length },
    { id: 'senior', label: 'Sr. Engineers Only', count: engineers.filter(e => e.engineer_level === 'L3').length },
    { id: 'certified', label: 'Manufacturer Certified', count: engineers.filter(e => e.manufacturer_certified).length },
    { id: 'high-match', label: 'Skills Match (80%+)', count: engineers.filter(e => (e.match_score || 0) >= 80).length },
    { 
      id: 'partner-network', 
      label: 'Partner Network', 
      count: engineers.filter(e => 
        e.category !== 'Manufacturer' || 
        e.organization_type === 'channel_partner' || 
        e.organization_type === 'sub_dealer'
      ).length 
    },
  ];

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <User className="h-5 w-5 text-blue-600" />
            Select Engineer for Assignment
          </DialogTitle>
          <DialogDescription>
            {equipmentName && (
              <span>Equipment: <strong>{equipmentName}</strong></span>
            )}
            {ticketId && (
              <span className="ml-2 text-xs text-gray-500">Ticket ID: {ticketId.slice(0, 8)}</span>
            )}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
            <div className="flex items-center gap-2 text-red-800">
              <AlertCircle className="h-5 w-5" />
              <span className="font-medium">{error}</span>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={fetchSuggestedEngineers}
              className="mt-2"
            >
              Retry
            </Button>
          </div>
        )}

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-center">
              <Loader2 className="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
              <p className="text-gray-600">Finding best engineers...</p>
            </div>
          </div>
        ) : engineers.length === 0 ? (
          <div className="text-center py-12">
            <User className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Engineers Available</h3>
            <p className="text-gray-600">
              No engineers found matching this equipment type and location.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Category Filters */}
            <div className="flex gap-2 flex-wrap">
              {categories.map(cat => (
                <Button
                  key={cat.id}
                  variant={activeCategory === cat.id ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveCategory(cat.id)}
                  className={activeCategory === cat.id ? 'bg-blue-600' : ''}
                >
                  {cat.label}
                  <Badge variant="secondary" className="ml-2">
                    {cat.count}
                  </Badge>
                </Button>
              ))}
            </div>

            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
              <p className="text-sm text-blue-800">
                <strong>{filteredEngineers.length}</strong> engineer{filteredEngineers.length !== 1 ? 's' : ''} available
                {' '}• Sorted by match score
              </p>
            </div>

            {filteredEngineers.map((engineer, index) => (
              <div
                key={engineer.engineer_id}
                className={`border rounded-lg p-4 transition-all ${
                  selectedEngineer === engineer.engineer_id
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-blue-300 hover:shadow-md'
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    {/* Header */}
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold text-lg">
                        {engineer.engineer_name.charAt(0)}
                      </div>
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2 flex-wrap">
                          {engineer.engineer_name}
                          {index === 0 && activeCategory === 'all' && (
                            <Badge className="bg-amber-100 text-amber-800 border-amber-300">
                              <Star className="h-3 w-3 mr-1" />
                              Recommended
                            </Badge>
                          )}
                          {getOrgTypeBadge(engineer.category, engineer.organization_type)}
                          {engineer.manufacturer_certified && (
                            <Badge className="bg-green-100 text-green-800 border-green-300">
                              <Award className="h-3 w-3 mr-1" />
                              Certified
                            </Badge>
                          )}
                        </h3>
                        <div className="flex items-center gap-4 mt-1 text-sm text-gray-600">
                          <div className="flex items-center gap-1">
                            <Building2 className="h-4 w-4" />
                            {engineer.organization_name}
                          </div>
                          {engineer.location && (
                            <div className="flex items-center gap-1">
                              <MapPin className="h-4 w-4" />
                              {engineer.location}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Details */}
                    <div className="grid grid-cols-2 gap-4 mb-3">
                      <div>
                        <span className="text-xs text-gray-500">Engineer Level</span>
                        <div className="mt-1">
                          <Badge className={getLevelColor(engineer.engineer_level)}>
                            {getLevelLabel(engineer.engineer_level)}
                          </Badge>
                        </div>
                      </div>
                      
                      {engineer.match_score !== undefined && (
                        <div>
                          <span className="text-xs text-gray-500">Match Score</span>
                          <div className={`mt-1 text-2xl font-bold ${getMatchScoreColor(engineer.match_score)}`}>
                            {engineer.match_score}%
                          </div>
                        </div>
                      )}

                      {engineer.equipment_types && engineer.equipment_types.length > 0 && (
                        <div className="col-span-2">
                          <span className="text-xs text-gray-500">Equipment Types</span>
                          <div className="mt-1 flex flex-wrap gap-2">
                            {engineer.equipment_types.map((type, idx) => (
                              <Badge key={idx} variant="outline" className="text-xs">
                                {type}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>

                    {/* Contact Info */}
                    {(engineer.phone || engineer.email) && (
                      <div className="flex items-center gap-4 text-sm text-gray-600">
                        {engineer.phone && (
                          <div className="flex items-center gap-1">
                            <Phone className="h-4 w-4" />
                            {engineer.phone}
                          </div>
                        )}
                        {engineer.email && (
                          <div className="flex items-center gap-1">
                            <Mail className="h-4 w-4" />
                            {engineer.email}
                          </div>
                        )}
                      </div>
                    )}
                  </div>

                  {/* Assign Button */}
                  <div className="ml-4">
                    <Button
                      onClick={() => {
                        setSelectedEngineer(engineer.engineer_id);
                        handleAssignEngineer(engineer.engineer_id);
                      }}
                      disabled={isAssigning}
                      className="bg-blue-600 hover:bg-blue-700"
                    >
                      {isAssigning && selectedEngineer === engineer.engineer_id ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Assigning...
                        </>
                      ) : (
                        <>
                          <CheckCircle2 className="h-4 w-4 mr-2" />
                          Assign
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="flex justify-end gap-2 mt-6 pt-4 border-t">
          <Button variant="outline" onClick={onClose} disabled={isAssigning}>
            Cancel
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
