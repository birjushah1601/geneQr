// API client for attachments and AI analysis
import { AttachmentInfo, AIAnalysisResult } from '@/lib/types/attachments'
import { apiClient } from '@/lib/api/client'

export interface AttachmentResponse {
  items: AttachmentInfo[]
  total: number
  page: number
  page_size: number
  has_next: boolean
  has_prev: boolean
}

export interface APIError {
  error: string
  message: string
  status: number
}

class AttachmentsApiError extends Error {
  constructor(public status: number, public apiError: APIError) {
    super(apiError.message)
    this.name = 'AttachmentsApiError'
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData: APIError = await response.json().catch(() => ({
      error: 'Unknown error',
      message: `HTTP ${response.status}`,
      status: response.status
    }))
    throw new AttachmentsApiError(response.status, errorData)
  }
  return response.json()
}

export const attachmentsApi = {
  // List attachments with optional filtering
  async list(params?: {
    page?: number
    page_size?: number
    ticket_id?: string
    status?: string
    category?: string
    source?: string
    unassigned?: boolean
  }): Promise<AttachmentResponse> {
    const searchParams = new URLSearchParams()
    if (params?.page) searchParams.set('page', params.page.toString())
    if (params?.page_size) searchParams.set('page_size', params.page_size.toString())
    if (params?.ticket_id) searchParams.set('ticket_id', params.ticket_id)
    if (params?.status) searchParams.set('status', params.status)
    if (params?.category) searchParams.set('category', params.category)
    if (params?.source) searchParams.set('source', params.source)
    if (params?.unassigned) searchParams.set('unassigned', 'true')

    const { data } = await apiClient.get(`/v1/attachments`, { params: Object.fromEntries(searchParams) })
    return data
  },

  // Get single attachment by ID
  async getById(id: string): Promise<AttachmentInfo> {
    const { data } = await apiClient.get(`/v1/attachments/${id}`)
    return data
  },

  // Get AI analysis for an attachment
  async getAIAnalysis(attachmentId: string): Promise<AIAnalysisResult> {
    const { data } = await apiClient.get(`/v1/attachments/${attachmentId}/ai-analysis`)
    return data
  },

  // Upload new attachment
  async upload(data: {
    file: File
    ticketId: string
    category?: string
    source?: string
  }): Promise<AttachmentInfo> {
    const formData = new FormData()
    formData.append('file', data.file)
    formData.append('ticket_id', data.ticketId)
    if (data.category) formData.append('category', data.category)
    if (data.source) formData.append('source', data.source)

    // Don't set Content-Type manually - let axios set it with boundary
    const { data: resp } = await apiClient.post(`/v1/attachments`, formData)
    return resp
  },

  // Delete attachment
  async delete(id: string): Promise<void> {
    await apiClient.delete(`/v1/attachments/${id}`)
  },

  // Link an existing attachment to a ticket (for pre-creation uploads)
  async link(id: string, ticketId: string): Promise<void> {
    await apiClient.post(`/v1/attachments/${id}/link`, { ticket_id: ticketId })
  },

  // Get attachment statistics
  async getStats(): Promise<{
    total: number
    by_status: Record<string, number>
    by_category: Record<string, number>
    processing_queue_size: number
    avg_confidence: number
  }> {
    const { data } = await apiClient.get('/v1/attachments/stats')
    return data
  },

  // Health check specifically for attachments service
  async healthCheck(): Promise<{ status: string; database: boolean; ai: boolean }> {
    const { data } = await apiClient.get('/v1/attachments/health')
    return data
  }
}

// Hook for React Query integration
export const attachmentQueryKeys = {
  all: ['attachments'] as const,
  lists: () => [...attachmentQueryKeys.all, 'list'] as const,
  list: (params: any) => [...attachmentQueryKeys.lists(), params] as const,
  details: () => [...attachmentQueryKeys.all, 'detail'] as const,
  detail: (id: string) => [...attachmentQueryKeys.details(), id] as const,
  aiAnalysis: (id: string) => [...attachmentQueryKeys.detail(id), 'ai-analysis'] as const,
  stats: () => [...attachmentQueryKeys.all, 'stats'] as const
}

