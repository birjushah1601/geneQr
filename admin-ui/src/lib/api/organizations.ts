import { apiClient } from './client';

export interface Organization {
  id: string;
  name: string;
  org_type: 'manufacturer' | 'distributor' | 'dealer' | 'hospital' | 'service_provider' | 'other';
  status: 'active' | 'inactive';
  metadata: any;
}

export interface Facility {
  id: string;
  org_id: string;
  facility_name: string;
  facility_code: string;
  facility_type: string;
  address: any;
  status: string;
}

export interface OrgRelationship {
  id: string;
  parent_org_id: string;
  child_org_id: string;
  rel_type: string;
}

export const organizationsApi = {
  list: async (params?: { 
    limit?: number; 
    offset?: number; 
    type?: string;
    status?: string;
  }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    if (params?.type) searchParams.set('type', params.type);
    if (params?.status) searchParams.set('status', params.status);
    
    const response = await apiClient.get<{ items: Organization[] }>(
      `/v1/organizations?${searchParams.toString()}`
    );
    return response.data.items || [];
  },

  get: async (id: string) => {
    const response = await apiClient.get<Organization>(`/v1/organizations/${id}`);
    return response.data;
  },

  listFacilities: async (orgId: string) => {
    const response = await apiClient.get<{ items: Facility[] }>(
      `/v1/organizations/${orgId}/facilities`
    );
    return response.data.items;
  },

  listRelationships: async (orgId: string) => {
    const response = await apiClient.get<{ items: OrgRelationship[] }>(
      `/v1/organizations/${orgId}/relationships`
    );
    return response.data.items;
  },

  // Stats endpoint (will calculate client-side for now)
  getStats: async () => {
    const orgs = await organizationsApi.list({ limit: 1000 });
    return {
      total: orgs.length,
      manufacturers: orgs.filter(o => o.org_type === 'manufacturer').length,
      distributors: orgs.filter(o => o.org_type === 'distributor').length,
      dealers: orgs.filter(o => o.org_type === 'dealer').length,
      hospitals: orgs.filter(o => o.org_type === 'hospital').length,
      active: orgs.filter(o => o.status === 'active').length,
    };
  },
};

