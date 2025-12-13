// ============================================================================
// AI Diagnosis API Service
// ============================================================================

import apiClient, { handleApiError } from './client';

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// AI DIAGNOSIS TYPES
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export interface DiagnosisRequest {
  ticket_id: number;
  equipment_id: string;
  symptoms: string[];
  description: string;
  images?: string[]; // Base64 encoded images
  user_context?: {
    user_id: number;
    username: string;
    role: string;
  };
  options?: {
    include_vision_analysis?: boolean;
    include_historical_context?: boolean;
    include_similar_tickets?: boolean;
    max_similar_tickets?: number;
    min_confidence_threshold?: number;
  };
}

export interface DiagnosisResult {
  problem_category: string;
  problem_type: string;
  description: string;
  confidence: number;
  severity: string;
  root_cause: string;
  symptoms: string[];
  possible_causes: string[];
  reasoning_explanation: string;
}

export interface AISuggestionMetadata {
  provider: string;
  model: string;
  confidence: number;
  confidence_factors: string[];
  alternatives_count: number;
  requires_feedback: boolean;
  suggestion_only: boolean;
}

export interface RecommendedAction {
  action: string;
  priority: string;
  description: string;
  estimated_time: string;
  requires_specialist: boolean;
  specialist_type?: string;
  required_tools: string[];
  required_parts: string[];
  safety_precautions: string[];
}

export interface RequiredPart {
  part_code: string;
  part_name: string;
  part_category: string;
  probability: number;
  quantity: number;
  is_oem_required: boolean;
  manufacturer: string;
  estimated_cost?: number;
}

export interface VisionFinding {
  attachment_id?: number;
  finding: string;
  confidence: number;
  category: string;
  location?: string;
}

export interface VisionAnalysisResult {
  attachments_analyzed: number;
  findings: VisionFinding[];
  overall_assessment: string;
  detected_components: string[];
  visible_damage: string[];
  confidence: number;
}

export interface DiagnosisResponse {
  diagnosis_id: string;
  ticket_id: number;
  primary_diagnosis: DiagnosisResult;
  alternate_diagnoses: DiagnosisResult[];
  
  // AI Confidence and Decision Tracking
  confidence: number;
  confidence_level: 'HIGH' | 'MEDIUM' | 'LOW';
  decision_status: 'pending' | 'accepted' | 'rejected';
  decided_by?: number;
  decided_at?: string;
  feedback_text?: string;
  
  // AI Metadata
  ai_metadata: AISuggestionMetadata;
  
  // Analysis results
  vision_analysis?: VisionAnalysisResult;
  recommended_actions: RecommendedAction[];
  required_parts: RequiredPart[];
  
  estimated_resolution_time?: string;
  created_at: string;
}

export interface DiagnosisDecisionFeedback {
  decision: 'accepted' | 'rejected';
  user_id: number;
  user_role: string;
  feedback_text?: string;
  corrections?: Record<string, any>;
}

export interface FeedbackResponse {
  success: boolean;
  message: string;
  feedback_id: string;
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// AI DIAGNOSIS SERVICE
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export const diagnosisApi = {
  /**
   * Run AI diagnosis analysis
   */
  async analyze(request: DiagnosisRequest): Promise<DiagnosisResponse> {
    try {
      const response = await apiClient.post<DiagnosisResponse>('/v1/diagnosis/analyze', request);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Submit user feedback on AI diagnosis decision
   */
  async submitFeedback(diagnosisId: string, feedback: DiagnosisDecisionFeedback): Promise<FeedbackResponse> {
    try {
      const response = await apiClient.post<FeedbackResponse>(
        `/v1/diagnosis/${diagnosisId}/feedback`, 
        feedback
      );
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get diagnosis by ID
   */
  async getById(diagnosisId: string): Promise<DiagnosisResponse> {
    try {
      const response = await apiClient.get<DiagnosisResponse>(`/v1/diagnosis/${diagnosisId}`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  },

  /**
   * Get diagnosis history for a ticket
   */
  async getHistoryByTicket(ticketId: number): Promise<DiagnosisResponse[]> {
    try {
      const response = await apiClient.get<DiagnosisResponse[]>(`/v1/diagnosis/ticket/${ticketId}/history`);
      return response.data;
    } catch (error) {
      throw new Error(handleApiError(error));
    }
  }
};

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// UTILITY FUNCTIONS
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

/**
 * Get confidence badge color based on confidence level
 */
export function getConfidenceBadgeColor(level: 'HIGH' | 'MEDIUM' | 'LOW') {
  switch (level) {
    case 'HIGH':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'MEDIUM':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'LOW':
      return 'bg-red-100 text-red-800 border-red-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
}

/**
 * Get confidence level description
 */
export function getConfidenceDescription(level: 'HIGH' | 'MEDIUM' | 'LOW') {
  switch (level) {
    case 'HIGH':
      return 'AI is very confident in this diagnosis';
    case 'MEDIUM':
      return 'AI has moderate confidence, human review recommended';
    case 'LOW':
      return 'AI has low confidence, human review required';
    default:
      return 'Unknown confidence level';
  }
}

/**
 * Extract symptoms from issue description
 */
export function extractSymptoms(description: string): string[] {
  // Simple extraction logic - could be enhanced with NLP
  const symptoms = description
    .toLowerCase()
    .split(/[.,;]/)
    .map(s => s.trim())
    .filter(s => s.length > 3)
    .slice(0, 5); // Limit to 5 symptoms
  
  return symptoms;
}

