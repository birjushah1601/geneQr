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

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// TYPES FOR ADDITIONAL ENDPOINTS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export interface TicketComment {
  id: string;
  ticket_id: string;
  comment_type?: string;
  author_id?: string;
  author_name?: string;
  comment: string;
  attachments?: any[];
  created_by?: string;
  created_at: string;
}

export interface TicketStatusHistory {
  id: string;
  ticket_id: string;
  from_status: string;
  to_status: string;
  changed_by: string;
  changed_at: string;
  notes?: string;
}

export interface FollowupTask {
  id: string;
  ticket_id: string;
  title: string;
  description?: string;
  task_type: string;
  priority: string;
  assigned_to?: string;
  assigned_to_name?: string;
  assigned_at?: string;
  due_date: string;
  status: 'pending' | 'in_progress' | 'overdue' | 'completed' | string;
  completed_at?: string;
  completed_by?: string;
  completion_notes?: string;
  created_at: string;
  created_by?: string;
  updated_at?: string;
}

export interface AddCommentRequest {
  comment: string;
  comment_type?: string;
}

export const ticketsApi = {
  /**
   * List tickets with optional filters
   */
  async list(params?: TicketListParams) {
    try {
      const queryString = params ? buildQueryString(params) : '';
      const response = await apiClient.get<{ items: ServiceTicket[]; total: number; page: number; page_size: number }>(
        `/v1/tickets?${queryString}`
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
      const response = await apiClient.get<ServiceTicket>(`/v1/tickets/${id}`);
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
      const response = await apiClient.get<ServiceTicket>(`/v1/tickets/number/${ticketNumber}`);
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

      const response = await apiClient.post<ServiceTicket>('/api/v1/tickets', payload);
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
      const response = await apiClient.post<{ message: string }>(`/v1/tickets/${ticketId}/assign`, data);
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
      const response = await apiClient.patch<{ message: string }>(`/v1/tickets/${ticketId}/status`, data);
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
      const response = await apiClient.patch<{ message: string }>(`/v1/tickets/${ticketId}`, data);
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
        `/v1/engineers/${engineerId}/tickets?${queryString}`
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
        `/v1/equipment/${equipmentId}/tickets`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  // ------------------------------------------------------------------
  // Conversation & Comments
  // ------------------------------------------------------------------
  async getComments(ticketId: string): Promise<{ comments: TicketComment[] }> {
    try {
      const response = await apiClient.get<{ comments: TicketComment[] }>(
        `/v1/tickets/${ticketId}/comments`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  async addComment(ticketId: string, payload: AddCommentRequest): Promise<{ comment: TicketComment }> {
    try {
      const response = await apiClient.post<{ comment: TicketComment }>(
        `/v1/tickets/${ticketId}/comments`,
        payload
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  // ------------------------------------------------------------------
  // Follow-up Tasks
  // ------------------------------------------------------------------
  async getFollowupTasks(ticketId: string): Promise<{ tasks: FollowupTask[] }> {
    try {
      const response = await apiClient.get<{ tasks: FollowupTask[] }>(
        `/v1/tickets/${ticketId}/followups`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  async completeFollowupTask(ticketId: string, taskId: string, completionNotes?: string): Promise<{ task: FollowupTask }> {
    try {
      const response = await apiClient.post<{ task: FollowupTask }>(
        `/v1/tickets/${ticketId}/followups/${taskId}/complete`,
        completionNotes ? { completion_notes: completionNotes } : undefined
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  async deleteFollowupTask(ticketId: string, taskId: string): Promise<void> {
    try {
      await apiClient.delete(`/v1/tickets/${ticketId}/followups/${taskId}`);
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },
};


