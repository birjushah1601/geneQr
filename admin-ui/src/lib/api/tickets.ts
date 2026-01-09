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

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// TYPES FOR ADDITIONAL ENDPOINTS
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

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
  comment_type: 'customer' | 'engineer' | 'internal' | 'system';
  author_name?: string;
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// MULTI-MODEL ASSIGNMENT TYPES
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export interface AssignmentSuggestionsResponse {
  ticket_id: string;
  equipment: {
    id: string;
    name: string;
    manufacturer: string;
    category: string;
    model_number: string;
    location?: {
      region: string;
      address: string;
      lat?: number;
      lng?: number;
    };
  };
  ticket: {
    priority: string;
    min_level_required: number;
    requires_certification: boolean;
  };
  suggestions_by_model: {
    [modelKey: string]: AssignmentModel;
  };
  assignment_tiers: {
    tier: number;
    name: string;
    organization_ids: string[];
    available_count: number;
  }[];
}

export interface AssignmentModel {
  model_name: string;
  description: string;
  engineers: EngineerSuggestion[];
  count: number;
}

export interface EngineerSuggestion {
  id: string;
  name: string;
  email: string;
  phone: string;
  engineer_level: number;
  skills?: string[];
  home_region?: string;
  organization_id?: string;
  organization_name?: string;
  organization_tier?: number;
  match_score: number;
  match_reasons: string[];
  workload?: {
    active_tickets: number;
    in_progress_tickets: number;
    avg_resolution_hours?: number;
  };
  certifications?: {
    manufacturer: string;
    category: string;
    is_certified: boolean;
    certification_number?: string;
    expiry?: string;
  }[];
  distance_km?: number;
  estimated_travel_time_mins?: number;
}

export interface AssignEngineerPayload {
  ticket_id: string;
  engineer_id: string;
  assignment_tier: string;
  assignment_tier_name: string;
  assigned_by: string;
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
      // Backend now accepts snake_case with JSON tags
      // Send data as-is, just clean up undefined fields
      const payload: Record<string, any> = {
        equipment_id: data.equipment_id,
        qr_code: data.qr_code,
        serial_number: data.serial_number,
        equipment_name: data.equipment_name,
        customer_id: data.customer_id,
        customer_name: data.customer_name,
        customer_phone: data.customer_phone,
        customer_whatsapp: data.customer_whatsapp,
        issue_category: data.issue_category,
        issue_description: data.issue_description,
        priority: data.priority,
        source: data.source,
        source_message_id: data.source_message_id,
        photos: data.photos,
        videos: data.videos,
        created_by: data.created_by,
        initial_comment: (data as any).notes || data.initial_comment,
        parts_requested: (data as any).parts_requested,
      };
      
      // Remove undefined keys to keep payload clean
      Object.keys(payload).forEach((k) => payload[k] === undefined && delete payload[k]);

      const response = await apiClient.post<ServiceTicket>('/v1/tickets', payload);
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

  // ------------------------------------------------------------------
  // Multi-Model Engineer Assignment
  // ------------------------------------------------------------------
  async getAssignmentSuggestions(ticketId: string): Promise<AssignmentSuggestionsResponse> {
    try {
      const response = await apiClient.get<AssignmentSuggestionsResponse>(
        `/v1/tickets/${ticketId}/assignment-suggestions`
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  async assignEngineerToTicket(ticketId: string, payload: AssignEngineerPayload): Promise<void> {
    try {
      await apiClient.post(
        `/v1/tickets/${ticketId}/assign-engineer`,
        payload
      );
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

  // Delete comment from a ticket
  async deleteComment(ticketId: string, commentId: string): Promise<void> {
    await apiClient.delete(`/v1/tickets/${ticketId}/comments/${commentId}`);
  },
};


