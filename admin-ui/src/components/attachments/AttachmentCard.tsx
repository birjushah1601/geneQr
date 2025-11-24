'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { 
  FileImage, 
  FileVideo, 
  FileText, 
  Download, 
  Eye, 
  Clock, 
  User, 
  MessageSquare,
  Bot,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Loader2
} from 'lucide-react'

// Import shared types
import { AttachmentInfo, AIAnalysisResult } from '@/lib/types/attachments'

// Local component-specific types that extend shared types
interface Attachment extends AttachmentInfo {
  ticketId: string
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
  analysisQuality: 'excellent' | 'good' | 'fair' | 'poor'
  processingDurationMs: number
  tokensUsed: number
  costUsd: number
  status: 'completed' | 'failed'
  errorMessage?: string
  analyzedAt: string
  
  // Parsed analysis results
  detectedObjects: Array<{
    name: string
    confidence: number
    boundingBox?: {
      x: number
      y: number
      width: number
      height: number
    }
  }>
  
  detectedIssues: Array<{
    issueType: string
    severity: 'low' | 'medium' | 'high' | 'critical'
    confidence: number
    evidence: string
    location?: string
  }>
  
  safetyConcerns: Array<{
    concernType: string
    severity: 'low' | 'medium' | 'high' | 'critical'
    action: string
    urgency: string
  }>
  
  repairRecommendations: Array<{
    action: string
    priority: 'low' | 'medium' | 'high' | 'urgent'
    estimatedTime: string
    difficulty: 'easy' | 'medium' | 'hard' | 'expert'
  }>
}

interface AttachmentCardProps {
  attachment: Attachment
  aiAnalysis?: AIAnalysisResult
  onViewAttachment?: (attachment: Attachment) => void
  onDownloadAttachment?: (attachment: Attachment) => void
  showAIAnalysis?: boolean
}

export function AttachmentCard({ 
  attachment, 
  aiAnalysis, 
  onViewAttachment, 
  onDownloadAttachment,
  showAIAnalysis = true 
}: AttachmentCardProps) {
  const [expanded, setExpanded] = useState(false)

  const getFileIcon = (fileType: string) => {
    if (fileType.startsWith('image/')) return <FileImage className="h-5 w-5 text-blue-500" />
    if (fileType.startsWith('video/')) return <FileVideo className="h-5 w-5 text-purple-500" />
    return <FileText className="h-5 w-5 text-gray-500" />
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
      case 'processed':
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case 'processing':
        return <Loader2 className="h-4 w-4 text-blue-500 animate-spin" />
      case 'failed':
        return <XCircle className="h-4 w-4 text-red-500" />
      default:
        return <Clock className="h-4 w-4 text-yellow-500" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
      case 'processed':
        return 'bg-green-100 text-green-800 border-green-200'
      case 'processing':
        return 'bg-blue-100 text-blue-800 border-blue-200'
      case 'failed':
        return 'bg-red-100 text-red-800 border-red-200'
      default:
        return 'bg-yellow-100 text-yellow-800 border-yellow-200'
    }
  }

  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'issue_photo':
        return 'bg-red-100 text-red-800 border-red-200'
      case 'equipment_photo':
        return 'bg-blue-100 text-blue-800 border-blue-200'
      case 'repair_photo':
        return 'bg-green-100 text-green-800 border-green-200'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200'
    }
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-900 border-red-300'
      case 'high':
        return 'bg-orange-100 text-orange-900 border-orange-300'
      case 'medium':
        return 'bg-yellow-100 text-yellow-900 border-yellow-300'
      case 'low':
        return 'bg-green-100 text-green-900 border-green-300'
      default:
        return 'bg-gray-100 text-gray-900 border-gray-300'
    }
  }

  const formatFileSize = (bytes: number): string => {
    const sizes = ['B', 'KB', 'MB', 'GB']
    if (bytes === 0) return '0 B'
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDateTime = (dateString: string): string => {
    return new Date(dateString).toLocaleString()
  }

  const getSourceIcon = (source: string) => {
    switch (source) {
      case 'whatsapp':
        return <MessageSquare className="h-4 w-4 text-green-600" />
      default:
        return <User className="h-4 w-4 text-gray-600" />
    }
  }

  return (
    <Card className="w-full hover:shadow-lg transition-shadow duration-200">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-start space-x-3">
            {getFileIcon(attachment.fileType)}
            <div className="flex-1 min-w-0">
              <h3 className="text-sm font-semibold text-gray-900 truncate">
                {attachment.originalFilename}
              </h3>
              <div className="flex items-center space-x-2 mt-1">
                <Badge variant="outline" className={getCategoryColor(attachment.category)}>
                  {attachment.category.replace('_', ' ')}
                </Badge>
                <Badge variant="outline" className={getStatusColor(attachment.processingStatus)}>
                  {getStatusIcon(attachment.processingStatus)}
                  <span className="ml-1">{attachment.processingStatus}</span>
                </Badge>
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            {onViewAttachment && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onViewAttachment(attachment)}
              >
                <Eye className="h-4 w-4" />
              </Button>
            )}
            {onDownloadAttachment && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onDownloadAttachment(attachment)}
              >
                <Download className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        {/* File Details */}
        <div className="grid grid-cols-2 gap-4 text-sm text-gray-600 mb-4">
          <div className="flex items-center space-x-2">
            {getSourceIcon(attachment.source)}
            <span>Source: {attachment.source}</span>
          </div>
          <div>Size: {formatFileSize(attachment.fileSizeBytes)}</div>
          <div>Type: {attachment.fileType}</div>
          <div>Uploaded: {formatDateTime(attachment.uploadedAt)}</div>
        </div>

        {/* AI Analysis Section */}
        {showAIAnalysis && aiAnalysis && (
          <>
            <Separator className="my-4" />
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Bot className="h-4 w-4 text-blue-600" />
                  <span className="font-medium text-gray-900">AI Analysis</span>
                  <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">
                    {Math.round(aiAnalysis.confidence * 100)}% confidence
                  </Badge>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setExpanded(!expanded)}
                >
                  {expanded ? 'Less' : 'More'} Details
                </Button>
              </div>

              {/* Quick Summary */}
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="text-gray-600">Quality:</span>
                  <Badge variant="outline" className="ml-2">
                    {aiAnalysis.analysisQuality}
                  </Badge>
                </div>
                <div>
                  <span className="text-gray-600">Processing time:</span>
                  <span className="ml-2">{aiAnalysis.processingDurationMs}ms</span>
                </div>
              </div>

              {/* Detected Issues (Always show if present) */}
              {aiAnalysis.detectedIssues && aiAnalysis.detectedIssues.length > 0 && (
                <div>
                  <h4 className="font-medium text-gray-900 mb-2 flex items-center">
                    <AlertTriangle className="h-4 w-4 text-orange-500 mr-2" />
                    Detected Issues ({aiAnalysis.detectedIssues.length})
                  </h4>
                  <div className="space-y-2">
                    {aiAnalysis.detectedIssues.slice(0, expanded ? undefined : 2).map((issue, index) => (
                      <div key={index} className="bg-gray-50 p-3 rounded-lg">
                        <div className="flex items-center justify-between mb-1">
                          <span className="font-medium text-gray-900">{issue.issueType}</span>
                          <Badge variant="outline" className={getSeverityColor(issue.severity)}>
                            {issue.severity}
                          </Badge>
                        </div>
                        <p className="text-sm text-gray-600">{issue.evidence}</p>
                        {issue.location && (
                          <p className="text-xs text-gray-500 mt-1">Location: {issue.location}</p>
                        )}
                      </div>
                    ))}
                    {!expanded && aiAnalysis.detectedIssues.length > 2 && (
                      <p className="text-sm text-gray-500 text-center">
                        +{aiAnalysis.detectedIssues.length - 2} more issues
                      </p>
                    )}
                  </div>
                </div>
              )}

              {/* Safety Concerns (Always show if present) */}
              {aiAnalysis.safetyConcerns && aiAnalysis.safetyConcerns.length > 0 && (
                <div>
                  <h4 className="font-medium text-gray-900 mb-2 flex items-center">
                    <AlertTriangle className="h-4 w-4 text-red-500 mr-2" />
                    Safety Concerns ({aiAnalysis.safetyConcerns.length})
                  </h4>
                  <div className="space-y-2">
                    {aiAnalysis.safetyConcerns.slice(0, expanded ? undefined : 1).map((concern, index) => (
                      <div key={index} className="bg-red-50 p-3 rounded-lg border border-red-200">
                        <div className="flex items-center justify-between mb-1">
                          <span className="font-medium text-gray-900">{concern.concernType}</span>
                          <Badge variant="outline" className={getSeverityColor(concern.severity)}>
                            {concern.severity}
                          </Badge>
                        </div>
                        <p className="text-sm text-gray-600">{concern.action}</p>
                        <p className="text-xs text-gray-500 mt-1">Urgency: {concern.urgency}</p>
                      </div>
                    ))}
                    {!expanded && aiAnalysis.safetyConcerns.length > 1 && (
                      <p className="text-sm text-gray-500 text-center">
                        +{aiAnalysis.safetyConcerns.length - 1} more concerns
                      </p>
                    )}
                  </div>
                </div>
              )}

              {/* Expanded Details */}
              {expanded && (
                <div className="space-y-4 border-t pt-4">
                  {/* Detected Objects */}
                  {aiAnalysis.detectedObjects && aiAnalysis.detectedObjects.length > 0 && (
                    <div>
                      <h4 className="font-medium text-gray-900 mb-2">
                        Detected Objects ({aiAnalysis.detectedObjects.length})
                      </h4>
                      <div className="flex flex-wrap gap-2">
                        {aiAnalysis.detectedObjects.map((obj, index) => (
                          <Badge key={index} variant="outline" className="bg-blue-50 text-blue-700">
                            {obj.name} ({Math.round(obj.confidence * 100)}%)
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Repair Recommendations */}
                  {aiAnalysis.repairRecommendations && aiAnalysis.repairRecommendations.length > 0 && (
                    <div>
                      <h4 className="font-medium text-gray-900 mb-2">
                        Repair Recommendations ({aiAnalysis.repairRecommendations.length})
                      </h4>
                      <div className="space-y-2">
                        {aiAnalysis.repairRecommendations.map((rec, index) => (
                          <div key={index} className="bg-green-50 p-3 rounded-lg">
                            <div className="flex items-center justify-between mb-1">
                              <span className="font-medium text-gray-900">{rec.action}</span>
                              <div className="flex space-x-2">
                                <Badge variant="outline" className={getSeverityColor(rec.priority)}>
                                  {rec.priority}
                                </Badge>
                                <Badge variant="outline" className="bg-gray-100 text-gray-700">
                                  {rec.difficulty}
                                </Badge>
                              </div>
                            </div>
                            <p className="text-xs text-gray-500">Est. time: {rec.estimatedTime}</p>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Analysis Metadata */}
                  <div className="bg-gray-50 p-3 rounded-lg">
                    <h4 className="font-medium text-gray-900 mb-2">Analysis Details</h4>
                    <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
                      <div>Provider: {aiAnalysis.aiProvider}</div>
                      <div>Model: {aiAnalysis.aiModel}</div>
                      <div>Tokens: {aiAnalysis.tokensUsed}</div>
                      <div>Cost: ${aiAnalysis.costUsd.toFixed(4)}</div>
                      <div className="col-span-2">
                        Analyzed: {formatDateTime(aiAnalysis.analyzedAt)}
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </>
        )}

        {/* Show message when AI analysis is pending */}
        {showAIAnalysis && !aiAnalysis && attachment.processingStatus === 'pending' && (
          <>
            <Separator className="my-4" />
            <div className="flex items-center space-x-2 text-gray-600">
              <Clock className="h-4 w-4" />
              <span className="text-sm">AI analysis queued for processing...</span>
            </div>
          </>
        )}

        {/* Show message when AI analysis is processing */}
        {showAIAnalysis && !aiAnalysis && attachment.processingStatus === 'processing' && (
          <>
            <Separator className="my-4" />
            <div className="flex items-center space-x-2 text-blue-600">
              <Loader2 className="h-4 w-4 animate-spin" />
              <span className="text-sm">AI analysis in progress...</span>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}