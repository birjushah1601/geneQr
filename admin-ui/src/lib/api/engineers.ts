// ============================================================================
// Engineers API Service
// ============================================================================

import apiClient, { buildQueryString, handleApiError } from './client';
import type {
  Engineer,
  CreateEngineerRequest,
  EngineerListParams,
  CSVImportResult,
} from '@/types';

export const engineersApi = {
  /**
   * List engineers with optional filters
   */
  async list(params?: EngineerListParams) {
    try {
      const queryString = params ? buildQueryString(params) : '';
      const response = await apiClient.get<{ engineers: Engineer[]; total: number; page: number; page_size: number; total_pages: number }>(
        `/v1/engineers?${queryString}`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get engineer by ID
   */
  async getById(id: string) {
    try {
      const response = await apiClient.get<Engineer>(`/v1/engineers/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Create new engineer
   */
  async create(data: any) {
    try {
      const response = await apiClient.post<Engineer>('/v1/engineers', data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Update engineer
   */
  async update(id: string, data: Partial<Engineer>) {
    try {
      const response = await apiClient.patch<{ message: string }>(`/v1/engineers/${id}`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Delete engineer
   */
  async delete(id: string) {
    try {
      const response = await apiClient.delete<{ message: string }>(`/v1/engineers/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Import engineers from CSV
   */
  async importCSV(file: File) {
    try {
      const formData = new FormData();
      formData.append('csv_file', file);
      
      const response = await apiClient.post<CSVImportResult>('/v1/engineers/import', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Update engineer availability
   */
  async updateAvailability(id: string, availability: 'available' | 'on_job' | 'off_duty') {
    try {
      const response = await apiClient.patch<{ message: string }>(`/v1/engineers/${id}/availability`, {
        availability,
      });
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },
};

export default engineersApi;
export type { Engineer };

