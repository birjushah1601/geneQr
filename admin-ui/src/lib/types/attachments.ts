// Shared type definitions for attachments and AI analysis

export interface AttachmentInfo {
  id: string
  fileName: string
  fileSize: number
  fileType: string
  uploadDate: string
  ticketId: string
  category: 'equipment_photo' | 'document' | 'error_screen' | 'manual' | 'other'
  status: 'pending' | 'processing' | 'completed' | 'failed'
  source: 'whatsapp' | 'web_upload' | 'email' | 'mobile_app'
}

export interface DetectedObject {
  name: string
  confidence: number
  location?: string
}

export interface DetectedIssue {
  issueType: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  confidence: number
  evidence: string
  location?: string
}

export interface SafetyConcern {
  concern: string
  concernType: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  evidence: string
  recommendation: string
  action: string
  urgency: string
}

export interface RepairRecommendation {
  task: string
  action: string
  priority: 'low' | 'medium' | 'high' | 'critical' | 'urgent'
  estimatedTime: string
  difficulty: 'easy' | 'medium' | 'hard' | 'expert'
  requiredParts?: string[]
  cost?: string
}

export interface AIAnalysisResult {
  id: string
  analysisId: string
  attachmentId: string
  ticketId: string
  aiProvider: string
  aiModel: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  confidence: number
  imageQualityScore: number
  analysisQuality: 'excellent' | 'good' | 'fair' | 'poor'
  processingDurationMs: number
  tokensUsed: number
  costUsd: number
  errorMessage?: string
  analyzedAt: string
  analysisDate: string
  processingTime?: number
  detectedObjects: DetectedObject[]
  detectedIssues: DetectedIssue[]
  safetyConcerns: SafetyConcern[]
  repairRecommendations: RepairRecommendation[]
  rawResponse?: string
}