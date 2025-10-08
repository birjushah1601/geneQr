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
      // Backend expects PascalCase keys (matching Go struct fields).
      // Map our snake_case CreateTicketRequest to the expected payload.
      const d: any = data as any;
      const payload: Record<string, any> = {
        EquipmentID: d.equipment_id ?? '',
        QRCode: d.qr_code ?? '',
        SerialNumber: d.serial_number ?? '',
        // Optional contextual fields if present
        EquipmentName: d.equipment_name ?? undefined,
        CustomerID: d.customer_id ?? undefined,
        CustomerName: d.customer_name ?? undefined,
        CustomerPhone: d.customer_phone,
        CustomerWhatsApp: d.customer_whatsapp ?? undefined,
        IssueCategory: d.issue_category,
        IssueDescription: d.issue_description,
        Priority: d.priority,
        Source: d.source,
        SourceMessageID: d.source_message_id ?? undefined,
        Photos: d.photos ?? undefined,
        Videos: d.videos ?? undefined,
        CreatedBy: d.created_by,
        InitialComment: d.notes ?? undefined,
      };
      // Remove undefined keys to keep payload clean
      Object.keys(payload).forEach((k) => payload[k] === undefined && delete payload[k]);

      const response = await apiClient.post<ServiceTicket>('/tickets', payload);
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
