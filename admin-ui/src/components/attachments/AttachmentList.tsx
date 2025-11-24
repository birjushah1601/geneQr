'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { AttachmentInfo, AIAnalysisResult } from '@/lib/types/attachments'
import { useAttachmentsWithFallback, useAttachmentStats, useRefreshAttachments } from '@/hooks/useAttachments'
import { 
  Search, 
  Filter, 
  RefreshCw,
  SortAsc,
  SortDesc,
  FileImage,
  Bot,
  AlertTriangle,
  TrendingUp,
  Clock,
  CheckCircle2,
  Loader2,
  WifiOff
} from 'lucide-react'
import { AttachmentCard } from './AttachmentCard'

// Use shared types with local extensions
interface Attachment extends AttachmentInfo {
  filename: string
  originalFilename: string
  fileSizeBytes: number
  storagePath: string
  uploadedById?: string
  sourceMessageId?: string
  processingStatus: 'pending' | 'processing' | 'processed' | 'completed' | 'failed'
  uploadedAt: string
  createdAt: string
  updatedAt: string
}

interface AttachmentListProps {
  ticketId?: string
  maxItems?: number
  showFilters?: boolean
  showStats?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
  onAttachmentSelect?: (attachment: Attachment) => void
}

export function AttachmentList({
  ticketId,
  maxItems,
  showFilters = true,
  showStats = true,
  autoRefresh = true,
  refreshInterval = 30000, // 30 seconds
  onAttachmentSelect
}: AttachmentListProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [categoryFilter, setCategoryFilter] = useState<string>('all')
  const [sourceFilter, setSourceFilter] = useState<string>('all')
  const [sortBy, setSortBy] = useState<'uploadedAt' | 'fileName' | 'fileSize'>('uploadedAt')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc')

  // Use real API with fallback to mock data
  const { 
    data: attachmentResponse, 
    isLoading, 
    error, 
    isMockData,
    refetch 
  } = useAttachmentsWithFallback({
    ticket_id: ticketId,
    status: statusFilter !== 'all' ? statusFilter : undefined,
    category: categoryFilter !== 'all' ? categoryFilter : undefined,
    source: sourceFilter !== 'all' ? sourceFilter : undefined,
    page_size: maxItems || 50,
    refetchInterval: autoRefresh ? refreshInterval : undefined
  })

  // Get attachment stats
  const { data: statsData } = useAttachmentStats(showStats)
  const { refreshAll } = useRefreshAttachments()

  // Get attachments from API response or empty array
  const attachments = attachmentResponse?.items || []
  
  // Create AI analyses map (in real implementation, this would come from separate API calls)
  const aiAnalyses = new Map<string, AIAnalysisResult>()

  // Convert API attachments to local format for compatibility
  const convertedAttachments: Attachment[] = attachments.map(item => ({
    ...item,
    filename: item.fileName,
    originalFilename: item.fileName,
    fileSizeBytes: item.fileSize,
    storagePath: `/storage/${item.source}/${item.id}`,
    processingStatus: item.status as any,
    uploadedAt: item.uploadDate,
    createdAt: item.uploadDate,
    updatedAt: item.uploadDate
  }))

  // Handle loading and error states
  const generateMockData = () => {
    const mockAttachments: Attachment[] = [
      {
        id: '1',
        ticketId: ticketId || 'TK-2025-001',
        filename: '20241118_153045_image_mri_scan.jpg',
        originalFilename: 'IMG-20241118-WA0001.jpg',
        fileType: 'image/jpeg',
        fileSizeBytes: 2456789,
        storagePath: '/storage/whatsapp/2024/11/18/mri_scan.jpg',
        source: 'whatsapp',
        sourceMessageId: 'wa_msg_001',
        category: 'issue_photo',
        processingStatus: 'processed',
        uploadedAt: '2024-11-18T15:30:45Z',
        createdAt: '2024-11-18T15:30:45Z',
        updatedAt: '2024-11-18T15:32:15Z'
      },
      {
        id: '2',
        ticketId: ticketId || 'TK-2025-002',
        filename: '20241118_143022_image_ct_scanner.jpg',
        originalFilename: 'IMG-20241118-WA0002.jpg',
        fileType: 'image/jpeg',
        fileSizeBytes: 3124567,
        storagePath: '/storage/whatsapp/2024/11/18/ct_scanner.jpg',
        source: 'whatsapp',
        sourceMessageId: 'wa_msg_002',
        category: 'equipment_photo',
        processingStatus: 'completed',
        uploadedAt: '2024-11-18T14:30:22Z',
        createdAt: '2024-11-18T14:30:22Z',
        updatedAt: '2024-11-18T14:31:45Z'
      },
      {
        id: '3',
        ticketId: ticketId || 'TK-2025-003',
        filename: '20241118_162103_image_xray_machine.jpg',
        originalFilename: 'IMG-20241118-WA0003.jpg',
        fileType: 'image/jpeg',
        fileSizeBytes: 1834256,
        storagePath: '/storage/whatsapp/2024/11/18/xray_machine.jpg',
        source: 'whatsapp',
        sourceMessageId: 'wa_msg_003',
        category: 'issue_photo',
        processingStatus: 'processing',
        uploadedAt: '2024-11-18T16:21:03Z',
        createdAt: '2024-11-18T16:21:03Z',
        updatedAt: '2024-11-18T16:21:03Z'
      },
      {
        id: '4',
        ticketId: ticketId || 'TK-2025-004',
        filename: '20241118_120815_image_ultrasound.jpg',
        originalFilename: 'IMG-20241118-WA0004.jpg',
        fileType: 'image/jpeg',
        fileSizeBytes: 2789345,
        storagePath: '/storage/whatsapp/2024/11/18/ultrasound.jpg',
        source: 'whatsapp',
        sourceMessageId: 'wa_msg_004',
        category: 'repair_photo',
        processingStatus: 'pending',
        uploadedAt: '2024-11-18T12:08:15Z',
        createdAt: '2024-11-18T12:08:15Z',
        updatedAt: '2024-11-18T12:08:15Z'
      }
    ]

    const mockAiAnalyses: AIAnalysisResult[] = [
      {
        id: 'ai_1',
        attachmentId: '1',
        ticketId: 'TK-2025-001',
        aiProvider: 'openai',
        aiModel: 'gpt-4-vision-preview',
        confidence: 0.87,
        imageQualityScore: 0.92,
        analysisQuality: 'good',
        processingDurationMs: 2340,
        tokensUsed: 1245,
        costUsd: 0.0089,
        status: 'completed',
        analyzedAt: '2024-11-18T15:32:15Z',
        detectedObjects: [
          { name: 'MRI Scanner Gantry', confidence: 0.95 },
          { name: 'Control Console', confidence: 0.88 },
          { name: 'Patient Table', confidence: 0.92 },
          { name: 'Emergency Stop Button', confidence: 0.78 }
        ],
        detectedIssues: [
          {
            issueType: 'Cooling System Alert',
            severity: 'medium',
            confidence: 0.72,
            evidence: 'Yellow indicator light visible on control panel',
          },
          {
            issueType: 'Display Screen Flickering',
            severity: 'low',
            confidence: 0.65,
            evidence: 'Intermittent display issues observed on main monitor'
          }
        ],
        safetyConcerns: [
          {
            concernType: 'Magnetic Safety',
            severity: 'medium',
            action: 'Verify safety signage placement',
            urgency: 'routine'
          }
        ],
        repairRecommendations: [
          {
            action: 'Check cooling system filters',
            priority: 'medium',
            estimatedTime: '30 minutes',
            difficulty: 'easy'
          },
          {
            action: 'Inspect display connections',
            priority: 'low',
            estimatedTime: '15 minutes',
            difficulty: 'easy'
          }
        ]
      },
      {
        id: 'ai_2',
        attachmentId: '2',
        ticketId: 'TK-2025-002',
        aiProvider: 'openai',
        aiModel: 'gpt-4-vision-preview',
        confidence: 0.91,
        imageQualityScore: 0.88,
        analysisQuality: 'good',
        processingDurationMs: 1890,
        tokensUsed: 987,
        costUsd: 0.0067,
        status: 'completed',
        analyzedAt: '2024-11-18T14:31:45Z',
        detectedObjects: [
          { name: 'CT Scanner', confidence: 0.96 },
          { name: 'Patient Bed', confidence: 0.89 },
          { name: 'Control Station', confidence: 0.85 }
        ],
        detectedIssues: [],
        safetyConcerns: [],
        repairRecommendations: [
          {
            action: 'Routine calibration check',
            priority: 'low',
            estimatedTime: '1 hour',
            difficulty: 'medium'
          }
        ]
      }
    ]

    return { mockAttachments, mockAiAnalyses }
  }

  const loadData = async () => {
    setLoading(true)
    try {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      const { mockAttachments, mockAiAnalyses } = generateMockData()
      
      // Filter by ticketId if provided
      const filteredAttachments = ticketId 
        ? mockAttachments.filter(att => att.ticketId === ticketId)
        : mockAttachments
      
      // Apply maxItems limit
      const limitedAttachments = maxItems 
        ? filteredAttachments.slice(0, maxItems)
        : filteredAttachments
      
      setAttachments(limitedAttachments)
      
      // Create AI analyses map
      const aiMap = new Map<string, AIAnalysisResult>()
      mockAiAnalyses.forEach(analysis => {
        aiMap.set(analysis.attachmentId, analysis)
      })
      setAiAnalyses(aiMap)
      
    } catch (error) {
      console.error('Error loading attachments:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [ticketId, maxItems])

  // Auto-refresh functionality
  useEffect(() => {
    if (!autoRefresh) return

    const interval = setInterval(loadData, refreshInterval)
    return () => clearInterval(interval)
  }, [autoRefresh, refreshInterval])

  // Filter and sort attachments using converted data
  const filteredAndSortedAttachments = React.useMemo(() => {
    let filtered = convertedAttachments.filter(attachment => {
      const matchesSearch = attachment.originalFilename.toLowerCase().includes(searchTerm.toLowerCase()) ||
                          attachment.ticketId.toLowerCase().includes(searchTerm.toLowerCase())
      
      const matchesStatus = statusFilter === 'all' || attachment.processingStatus === statusFilter
      const matchesCategory = categoryFilter === 'all' || attachment.category === categoryFilter
      const matchesSource = sourceFilter === 'all' || attachment.source === sourceFilter
      
      return matchesSearch && matchesStatus && matchesCategory && matchesSource
    })

    // Sort
    filtered.sort((a, b) => {
      let aValue, bValue
      
      switch (sortBy) {
        case 'fileName':
          aValue = a.originalFilename.toLowerCase()
          bValue = b.originalFilename.toLowerCase()
          break
        case 'fileSize':
          aValue = a.fileSizeBytes
          bValue = b.fileSizeBytes
          break
        default: // uploadedAt
          aValue = new Date(a.uploadedAt).getTime()
          bValue = new Date(b.uploadedAt).getTime()
      }
      
      if (sortOrder === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
      }
    })

    return filtered
  }, [convertedAttachments, searchTerm, statusFilter, categoryFilter, sourceFilter, sortBy, sortOrder])

  // Calculate statistics from API data or use stats from API
  const stats = React.useMemo(() => {
    if (statsData) {
      return {
        total: statsData.total,
        processed: statsData.by_status?.completed || 0,
        processing: statsData.by_status?.processing || 0,
        pending: statsData.by_status?.pending || 0,
        failed: statsData.by_status?.failed || 0,
        withIssues: 0, // Would come from AI analysis API
        withSafetyConcerns: 0, // Would come from AI analysis API  
        avgConfidence: statsData.avg_confidence || 0
      }
    }
    
    // Fallback to client-side calculation
    const total = convertedAttachments.length
    const processed = convertedAttachments.filter(a => a.processingStatus === 'processed' || a.processingStatus === 'completed').length
    const processing = convertedAttachments.filter(a => a.processingStatus === 'processing').length
    const pending = convertedAttachments.filter(a => a.processingStatus === 'pending').length
    const failed = convertedAttachments.filter(a => a.processingStatus === 'failed').length
    
    return {
      total,
      processed,
      processing,
      pending,
      failed,
      withIssues: 0,
      withSafetyConcerns: 0,
      avgConfidence: 0.87 // Mock average
    }
  }, [convertedAttachments, statsData])

  const handleViewAttachment = (attachment: Attachment) => {
    if (onAttachmentSelect) {
      onAttachmentSelect(attachment)
    }
    // Default behavior - could open in modal or new tab
    console.log('View attachment:', attachment.storagePath)
  }

  const handleDownloadAttachment = (attachment: Attachment) => {
    // Simulate download
    console.log('Download attachment:', attachment.filename)
    // In real implementation: window.open(attachment.downloadUrl)
  }

  const handleRefresh = () => {
    refetch()
    refreshAll()
  }

  // Loading state
  if (isLoading) {
    return (
      <Card className="w-full">
        <CardContent className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 animate-spin mr-2" />
          <span>Loading attachments...</span>
        </CardContent>
      </Card>
    )
  }

  // Error state (with fallback data)
  const showErrorBanner = error && !isMockData

  return (
    <div className="space-y-6">
      {/* Statistics Dashboard */}
      {showStats && (
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <FileImage className="h-4 w-4 text-blue-500" />
                <div>
                  <p className="text-sm font-medium">Total</p>
                  <p className="text-2xl font-bold">{stats.total}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <CheckCircle2 className="h-4 w-4 text-green-500" />
                <div>
                  <p className="text-sm font-medium">Processed</p>
                  <p className="text-2xl font-bold text-green-600">{stats.processed}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <RefreshCw className="h-4 w-4 text-blue-500" />
                <div>
                  <p className="text-sm font-medium">Processing</p>
                  <p className="text-2xl font-bold text-blue-600">{stats.processing}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Clock className="h-4 w-4 text-yellow-500" />
                <div>
                  <p className="text-sm font-medium">Pending</p>
                  <p className="text-2xl font-bold text-yellow-600">{stats.pending}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <AlertTriangle className="h-4 w-4 text-orange-500" />
                <div>
                  <p className="text-sm font-medium">Issues</p>
                  <p className="text-2xl font-bold text-orange-600">{stats.withIssues}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <AlertTriangle className="h-4 w-4 text-red-500" />
                <div>
                  <p className="text-sm font-medium">Safety</p>
                  <p className="text-2xl font-bold text-red-600">{stats.withSafetyConcerns}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-4">
              <div className="flex items-center space-x-2">
                <Bot className="h-4 w-4 text-purple-500" />
                <div>
                  <p className="text-sm font-medium">AI Avg</p>
                  <p className="text-2xl font-bold text-purple-600">{Math.round(stats.avgConfidence * 100)}%</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Filters and Search */}
      {showFilters && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-lg">Attachments & AI Analysis</CardTitle>
              <div className="flex items-center space-x-2">
                {isMockData && (
                  <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
                    <WifiOff className="h-3 w-3 mr-1" />
                    Using Mock Data
                  </Badge>
                )}
                <Button variant="outline" size="sm" onClick={handleRefresh}>
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Refresh
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Search */}
            <div className="flex items-center space-x-2">
              <Search className="h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search by filename or ticket ID..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="flex-1"
              />
            </div>

            {/* Filters */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label className="text-sm font-medium mb-1 block">Status</label>
                <select
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-200 rounded-md text-sm"
                >
                  <option value="all">All Status</option>
                  <option value="pending">Pending</option>
                  <option value="processing">Processing</option>
                  <option value="processed">Processed</option>
                  <option value="completed">Completed</option>
                  <option value="failed">Failed</option>
                </select>
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Category</label>
                <select
                  value={categoryFilter}
                  onChange={(e) => setCategoryFilter(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-200 rounded-md text-sm"
                >
                  <option value="all">All Categories</option>
                  <option value="issue_photo">Issue Photo</option>
                  <option value="equipment_photo">Equipment Photo</option>
                  <option value="repair_photo">Repair Photo</option>
                  <option value="document">Document</option>
                  <option value="video">Video</option>
                  <option value="audio">Audio</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Source</label>
                <select
                  value={sourceFilter}
                  onChange={(e) => setSourceFilter(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-200 rounded-md text-sm"
                >
                  <option value="all">All Sources</option>
                  <option value="whatsapp">WhatsApp</option>
                  <option value="web_upload">Web Upload</option>
                  <option value="email">Email</option>
                  <option value="api">API</option>
                </select>
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Sort</label>
                <div className="flex space-x-2">
                  <select
                    value={sortBy}
                    onChange={(e) => setSortBy(e.target.value as any)}
                    className="flex-1 px-3 py-2 border border-gray-200 rounded-md text-sm"
                  >
                    <option value="uploadedAt">Upload Date</option>
                    <option value="filename">Filename</option>
                    <option value="fileSize">File Size</option>
                  </select>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
                  >
                    {sortOrder === 'asc' ? <SortAsc className="h-4 w-4" /> : <SortDesc className="h-4 w-4" />}
                  </Button>
                </div>
              </div>
            </div>

            {/* Results count */}
            <div className="flex items-center justify-between text-sm text-gray-600">
              <span>
                Showing {filteredAndSortedAttachments.length} of {stats.total} attachments
              </span>
              {searchTerm && (
                <Button variant="ghost" size="sm" onClick={() => setSearchTerm('')}>
                  Clear search
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Attachment List */}
      <div className="space-y-4">
        {filteredAndSortedAttachments.length === 0 ? (
          <Card>
            <CardContent className="py-8 text-center">
              <FileImage className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No attachments found</h3>
              <p className="text-gray-600">
                {searchTerm ? 'Try adjusting your search terms or filters.' : 'No attachments have been uploaded yet.'}
              </p>
            </CardContent>
          </Card>
        ) : (
          filteredAndSortedAttachments.map((attachment) => (
            <AttachmentCard
              key={attachment.id}
              attachment={attachment}
              aiAnalysis={aiAnalyses.get(attachment.id)}
              onViewAttachment={handleViewAttachment}
              onDownloadAttachment={handleDownloadAttachment}
              showAIAnalysis={true}
            />
          ))
        )}
      </div>
    </div>
  )
}
