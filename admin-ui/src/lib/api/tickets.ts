// ============================================================================
// Service Tickets API Service
// ============================================================================

import apiClient, { buildQueryString, handleApiError } from './client';
import type {
  ServiceTicket,
  CreateTicketRequest,
  AssignEngineerRequest,
  UpdateTicketStatusRequest,
  TicketListParams,
} from '@/types';

export const ticketsApi = {
  /**
   * List tickets with optional filters
   */
  async list(params?: TicketListParams) {
    try {
      const queryString = params ? buildQueryString(params) : '';
      const response = await apiClient.get<{ tickets: ServiceTicket[]; total: number; page: number; page_size: number; total_pages: number }>(
        `/tickets?${queryString}`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get ticket by ID
   */
  async getById(id: string) {
    try {
      const response = await apiClient.get<ServiceTicket>(`/tickets/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get ticket by ticket number
   */
  async getByTicketNumber(ticketNumber: string) {
    try {
      const response = await apiClient.get<ServiceTicket>(`/tickets/number/${ticketNumber}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Create new ticket
   */
  async create(data: CreateTicketRequest) {
    try {
      const response = await apiClient.post<ServiceTicket>('/tickets', data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Assign engineer to ticket
   */
  async assignEngineer(ticketId: string, data: AssignEngineerRequest) {
    try {
      const response = await apiClient.post<{ message: string }>(`/tickets/${ticketId}/assign`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Update ticket status
   */
  async updateStatus(ticketId: string, data: UpdateTicketStatusRequest) {
    try {
      const response = await apiClient.patch<{ message: string }>(`/tickets/${ticketId}/status`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Update ticket
   */
  async update(ticketId: string, data: Partial<ServiceTicket>) {
    try {
      const response = await apiClient.patch<{ message: string }>(`/tickets/${ticketId}`, data);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get tickets by engineer
   */
  async getByEngineer(engineerId: string, params?: TicketListParams) {
    try {
      const queryString = params ? buildQueryString(params) : '';
      const response = await apiClient.get<{ tickets: ServiceTicket[]; total: number }>(
        `/engineers/${engineerId}/tickets?${queryString}`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get tickets by equipment
   */
  async getByEquipment(equipmentId: string) {
    try {
      const response = await apiClient.get<{ tickets: ServiceTicket[] }>(
        `/equipment/${equipmentId}/tickets`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },
};
