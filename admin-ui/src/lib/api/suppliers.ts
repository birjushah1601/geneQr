// ============================================================================
// Suppliers API Client
// ============================================================================

import apiClient from './client';

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TYPES
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export interface Supplier {
  id: string;
  name: string;
  contact_person?: string;
  email?: string;
  phone?: string;
  category?: string;
  location?: string;
  status: 'active' | 'inactive' | 'pending' | 'suspended';
  verification_status?: 'pending' | 'verified' | 'rejected';
  rating?: number;
  created_at: string;
  updated_at: string;
}

export interface ListSuppliersResponse {
  items: Supplier[];
  total: number;
  page: number;
  page_size: number;
}

export interface SupplierCertification {
  name: string;
  issuer: string;
  issue_date: string;
  expiry_date?: string;
  certificate_url?: string;
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// API CLIENT
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export const suppliersApi = {
  /**
   * List suppliers with optional filters
   */
  list: async (params?: {
    status?: string[];
    verification_status?: string[];
    category_id?: string;
    search?: string;
    sort_by?: string;
    sort_direction?: string;
    page?: number;
    page_size?: number;
  }): Promise<ListSuppliersResponse> => {
    const response = await apiClient.get('/v1/suppliers', { params });
    return response.data;
  },

  /**
   * Get single supplier by ID
   */
  getById: async (id: string): Promise<Supplier> => {
    const response = await apiClient.get(`/v1/suppliers/${id}`);
    return response.data;
  },

  /**
   * Create a new supplier
   */
  create: async (
    data: Omit<Supplier, 'id' | 'created_at' | 'updated_at'>
  ): Promise<Supplier> => {
    const response = await apiClient.post('/v1/suppliers', data);
    return response.data;
  },

  /**
   * Update an existing supplier
   */
  update: async (id: string, data: Partial<Supplier>): Promise<Supplier> => {
    const response = await apiClient.put(`/v1/suppliers/${id}`, data);
    return response.data;
  },

  /**
   * Delete a supplier
   */
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/v1/suppliers/${id}`);
  },

  /**
   * Verify a supplier
   */
  verify: async (id: string): Promise<Supplier> => {
    const response = await apiClient.post(`/v1/suppliers/${id}/verify`);
    return response.data;
  },

  /**
   * Reject a supplier
   */
  reject: async (id: string): Promise<Supplier> => {
    const response = await apiClient.post(`/v1/suppliers/${id}/reject`);
    return response.data;
  },

  /**
   * Suspend a supplier
   */
  suspend: async (id: string): Promise<Supplier> => {
    const response = await apiClient.post(`/v1/suppliers/${id}/suspend`);
    return response.data;
  },

  /**
   * Activate a supplier
   */
  activate: async (id: string): Promise<Supplier> => {
    const response = await apiClient.post(`/v1/suppliers/${id}/activate`);
    return response.data;
  },

  /**
   * Add certification to supplier
   */
  addCertification: async (
    id: string,
    certification: SupplierCertification
  ): Promise<Supplier> => {
    const response = await apiClient.post(
      `/v1/suppliers/${id}/certifications`,
      certification
    );
    return response.data;
  },

  /**
   * Get suppliers by category
   */
  getByCategory: async (categoryId: string): Promise<ListSuppliersResponse> => {
    const response = await apiClient.get(`/v1/suppliers/category/${categoryId}`);
    return response.data;
  },
};
