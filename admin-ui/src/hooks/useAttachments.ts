// React hooks for attachments data fetching
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { attachmentsApi, attachmentQueryKeys, AttachmentResponse } from '@/lib/api/attachments'
import { AttachmentInfo, AIAnalysisResult } from '@/lib/types/attachments'

export interface UseAttachmentsOptions {
  page?: number
  page_size?: number
  ticket_id?: string
  status?: string
  category?: string
  source?: string
  enabled?: boolean
  refetchInterval?: number
}

// Main hook for listing attachments
export function useAttachments(options: UseAttachmentsOptions = {}) {
  const {
    page = 1,
    page_size = 20,
    enabled = true,
    refetchInterval,
    ...filters
  } = options

  return useQuery({
    queryKey: attachmentQueryKeys.list({ page, page_size, ...filters }),
    queryFn: () => attachmentsApi.list({ page, page_size, ...filters }),
    enabled,
    refetchInterval,
    staleTime: 30 * 1000, // 30 seconds
    retry: 2,
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000)
  })
}

// Hook for single attachment
export function useAttachment(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: attachmentQueryKeys.detail(id),
    queryFn: () => attachmentsApi.getById(id),
    enabled: enabled && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 1
  })
}

// Hook for AI analysis
export function useAIAnalysis(attachmentId: string, enabled: boolean = true) {
  return useQuery({
    queryKey: attachmentQueryKeys.aiAnalysis(attachmentId),
    queryFn: () => attachmentsApi.getAIAnalysis(attachmentId),
    enabled: enabled && !!attachmentId,
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: 2
  })
}

// Hook for attachment statistics
export function useAttachmentStats(enabled: boolean = true) {
  return useQuery({
    queryKey: attachmentQueryKeys.stats(),
    queryFn: () => attachmentsApi.getStats(),
    enabled,
    staleTime: 60 * 1000, // 1 minute
    retry: 2
  })
}

// Hook with fallback to mock data for development
export function useAttachmentsWithFallback(options: UseAttachmentsOptions = {}) {
  const queryClient = useQueryClient()
  
  // Try real API first
  const realQuery = useAttachments({ ...options, enabled: true })
  
  // If real API fails, provide mock data
  if (realQuery.error && !realQuery.data) {
    // Create mock data that matches our type structure
    const mockData: AttachmentResponse = {
      items: [
        {
          id: '1',
          fileName: '20241118_153045_image_mri_scan.jpg',
          fileSize: 2456789,
          fileType: 'image/jpeg',
          uploadDate: '2024-11-18T15:30:45Z',
          ticketId: 'TK-2025-001',
          category: 'equipment_photo',
          status: 'completed',
          source: 'whatsapp'
        },
        {
          id: '2',
          fileName: 'ct_scanner_error_20241118_143022.jpg',
          fileSize: 3245678,
          fileType: 'image/jpeg',
          uploadDate: '2024-11-18T14:30:22Z',
          ticketId: 'TK-2025-002',
          category: 'document',
          status: 'processing',
          source: 'whatsapp'
        }
      ],
      total: 2,
      page: 1,
      page_size: 20,
      has_next: false,
      has_prev: false
    }

    // Set mock data in query cache
    queryClient.setQueryData(
      attachmentQueryKeys.list({ page: 1, page_size: 20 }),
      mockData
    )

    return {
      ...realQuery,
      data: mockData,
      isLoading: false,
      error: null,
      isError: false,
      isMockData: true
    }
  }

  return {
    ...realQuery,
    isMockData: false
  }
}

// Custom hook for refreshing attachment data
export function useRefreshAttachments() {
  const queryClient = useQueryClient()
  
  return {
    refreshAll: () => queryClient.invalidateQueries({ queryKey: attachmentQueryKeys.all }),
    refreshList: (params?: any) => queryClient.invalidateQueries({ queryKey: attachmentQueryKeys.list(params) }),
    refreshStats: () => queryClient.invalidateQueries({ queryKey: attachmentQueryKeys.stats() })
  }
}