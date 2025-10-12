// ============================================================================
// Manufacturers API Client
// ============================================================================

import apiClient from './client';

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TYPES
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export interface Manufacturer {
  id: string;
  name: string;
  contact_person?: string;
  email?: string;
  phone?: string;
  website?: string;
  address?: string;
  status: 'active' | 'inactive' | 'pending';
  created_at: string;
  updated_at: string;
}

export interface ListManufacturersResponse {
  items: Manufacturer[];
  total: number;
  page: number;
  page_size: number;
}

export interface ManufacturerStats {
  equipmentCount: number;
  engineersCount: number;
  activeTickets: number;
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// API CLIENT
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export const manufacturersApi = {
  /**
   * List all manufacturers with optional filters
   */
  list: async (params?: {
    search?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<ListManufacturersResponse> => {
    const response = await apiClient.get('/v1/organizations', { params });
    return response.data;
  },

  /**
   * Get single manufacturer by ID
   */
  getById: async (id: string): Promise<Manufacturer> => {
    const response = await apiClient.get(`/v1/organizations/${id}`);
    return response.data;
  },

  /**
   * Create a new manufacturer
   */
  create: async (
    data: Omit<Manufacturer, 'id' | 'created_at' | 'updated_at'>
  ): Promise<Manufacturer> => {
    const response = await apiClient.post('/v1/organizations', data);
    return response.data;
  },

  /**
   * Update an existing manufacturer
   */
  update: async (id: string, data: Partial<Manufacturer>): Promise<Manufacturer> => {
    const response = await apiClient.put(`/v1/organizations/${id}`, data);
    return response.data;
  },

  /**
   * Get manufacturer's equipment
   */
  getEquipment: async (
    id: string,
    params?: {
      page?: number;
      page_size?: number;
      status?: string;
    }
  ) => {
    const response = await apiClient.get('/v1/equipment', {
      params: {
        ...params,
        customer_id: id,
      },
    });
    return response.data;
  },

  /**
   * Get manufacturer's engineers
   */
  getEngineers: async (
    id: string,
    params?: {
      limit?: number;
      offset?: number;
    }
  ) => {
    const response = await apiClient.get('/v1/organizations/engineers', {
      params: {
        ...params,
        org_id: id,
      },
    });
    return response.data;
  },

  /**
   * Get manufacturer's service tickets
   */
  getTickets: async (
    id: string,
    params?: {
      page?: number;
      page_size?: number;
      status?: string;
    }
  ) => {
    const response = await apiClient.get('/v1/tickets', {
      params: {
        ...params,
        customer_id: id,
      },
    });
    return response.data;
  },

  /**
   * Get manufacturer statistics (aggregated data)
   */
  getStats: async (id: string): Promise<ManufacturerStats> => {
    // Fetch all data in parallel
    const [equipment, engineers, tickets] = await Promise.all([
      manufacturersApi.getEquipment(id, { page: 1, page_size: 1 }),
      manufacturersApi.getEngineers(id, { limit: 1, offset: 0 }),
      manufacturersApi.getTickets(id, { page: 1, page_size: 1 }),
    ]);

    return {
      equipmentCount: equipment.total || 0,
      engineersCount: engineers.total || 0,
      activeTickets: tickets.total || 0,
    };
  },
};
