// ============================================================================
// Equipment API Service
// ============================================================================

import apiClient, { buildQueryString, handleApiError } from './client';
import type {
  Equipment,
  RegisterEquipmentRequest,
  EquipmentListParams,
  CSVImportResult,
} from '@/types';

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// EQUIPMENT SERVICE
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export const equipmentApi = {
  /**
   * List equipment with optional filters
   */
  async list(params?: EquipmentListParams) {
    try {
      const queryString = params ? buildQueryString(params) : '';
      const response = await apiClient.get<{ equipment: Equipment[]; total: number; page: number; page_size: number; total_pages: number }>(
        `/equipment?${queryString}`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get equipment by ID
   */
  async getById(id: string) {
    try {
      const response = await apiClient.get<Equipment>(`/equipment/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get equipment by QR code
   */
  async getByQRCode(qrCode: string) {
    try {
      const response = await apiClient.get<Equipment>(`/equipment/qr/${qrCode}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get equipment by serial number
   */
  async getBySerial(serialNumber: string) {
    try {
      const response = await apiClient.get<Equipment>(`/equipment/serial/${serialNumber}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Register new equipment
   */
  async register(data: RegisterEquipmentRequest) {
    try {
      const response = await apiClient.post<Equipment>('/equipment', data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Update equipment
   */
  async update(id: string, data: Partial<Equipment>) {
    try {
      const response = await apiClient.patch<{ message: string }>(`/equipment/${id}`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Generate QR code for equipment
   */
  async generateQRCode(id: string) {
    try {
      const response = await apiClient.post<{ message: string; path: string }>(`/equipment/${id}/qr`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Bulk generate QR codes for all equipment without QR codes
   */
  async bulkGenerateQRCodes() {
    try {
      const response = await apiClient.post<{ 
        total_processed: number; 
        successful: number; 
        failed: number; 
        failed_ids?: string[];
        message: string 
      }>('/equipment/qr/bulk-generate');
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Download QR label PDF
   */
  async downloadQRLabel(id: string) {
    try {
      const response = await apiClient.get(`/equipment/${id}/qr/pdf`, {
        responseType: 'blob',
      });
      
      // Create download link
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `qr_label_${id}.pdf`);
      document.body.appendChild(link);
      link.click();
      link.remove();
      
      return true;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Import equipment from CSV
   */
  async importCSV(file: File, createdBy: string) {
    try {
      const formData = new FormData();
      formData.append('csv_file', file);
      formData.append('created_by', createdBy);
      
      const response = await apiClient.post<CSVImportResult>('/equipment/import', formData, {
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
   * Record service for equipment
   */
  async recordService(id: string, data: { service_date: string; notes: string }) {
    try {
      const response = await apiClient.post<{ message: string }>(`/equipment/${id}/service`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },
};
