// Partner Association API Client

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081/api';

export interface Partner {
  id: string;
  partner_org_id: string;
  partner_name: string;
  org_type: 'channel_partner' | 'sub_dealer';
  equipment_id?: string;
  equipment_name?: string;
  association_type: 'general' | 'equipment-specific';
  engineers_count: number;
  created_at: string;
}

export interface Organization {
  id: string;
  name: string;
  org_type: string;
  location?: string;
  engineers_count: number;
  contact_email?: string;
}

export interface NetworkEngineer {
  id: string;
  name: string;
  phone: string;
  email: string;
  organization: {
    id: string;
    name: string;
    org_type: string;
  };
  category: string;
}

export interface NetworkEngineersResponse {
  equipment_id?: string;
  engineers: NetworkEngineer[];
  grouped: Record<string, Array<{
    id: string;
    name: string;
    phone: string;
    email: string;
  }>>;
  total_engineers: number;
  association_type: string;
}

class PartnersAPI {
  private getAuthHeaders() {
    const token = localStorage.getItem('access_token');
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : '',
    };
  }

  // List partners for a manufacturer
  async listPartners(
    manufacturerId: string, 
    filters?: { type?: string; association_type?: string }
  ): Promise<{ partners: Partner[]; total: number }> {
    const params = new URLSearchParams();
    if (filters?.type) params.append('type', filters.type);
    if (filters?.association_type) params.append('association_type', filters.association_type);

    const url = `${API_BASE_URL}/v1/organizations/${manufacturerId}/partners${params.toString() ? '?' + params.toString() : ''}`;
    
    const response = await fetch(url, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch partners');
    }

    return response.json();
  }

  // Get available partners (not yet associated)
  async getAvailablePartners(
    manufacturerId: string,
    search?: string
  ): Promise<{ organizations: Organization[]; total: number }> {
    const params = new URLSearchParams();
    if (search) params.append('search', search);

    const url = `${API_BASE_URL}/v1/organizations/${manufacturerId}/available-partners${params.toString() ? '?' + params.toString() : ''}`;
    
    const response = await fetch(url, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch available partners');
    }

    return response.json();
  }

  // Associate a partner with manufacturer
  async associatePartner(
    manufacturerId: string,
    data: {
      partner_org_id: string;
      equipment_id?: string;
      rel_type?: string;
    }
  ): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/v1/organizations/${manufacturerId}/partners`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify({
        ...data,
        rel_type: data.rel_type || 'services_for',
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to associate partner');
    }

    return response.json();
  }

  // Remove partner association
  async removePartner(
    manufacturerId: string,
    partnerId: string,
    equipmentId?: string
  ): Promise<void> {
    const params = new URLSearchParams();
    if (equipmentId) params.append('equipment_id', equipmentId);

    const url = `${API_BASE_URL}/v1/organizations/${manufacturerId}/partners/${partnerId}${params.toString() ? '?' + params.toString() : ''}`;
    
    const response = await fetch(url, {
      method: 'DELETE',
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to remove partner');
    }
  }

  // Get network engineers (with smart filtering)
  async getNetworkEngineers(
    manufacturerId: string,
    equipmentId?: string
  ): Promise<NetworkEngineersResponse> {
    const params = new URLSearchParams();
    if (equipmentId) params.append('equipment_id', equipmentId);

    const url = `${API_BASE_URL}/v1/engineers/network/${manufacturerId}${params.toString() ? '?' + params.toString() : ''}`;
    
    const response = await fetch(url, {
      headers: this.getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch network engineers');
    }

    return response.json();
  }
}

export const partnersApi = new PartnersAPI();
