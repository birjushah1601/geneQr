// ============================================================================
// AI Diagnosis Helpers
// ============================================================================

import { diagnosisApi } from '@/lib/api/diagnosis';
import { attachmentsApi } from '@/lib/api/attachments';
import { compressImage, extractVideoFrames, isImageFile, isVideoFile } from './imageUtils';

export interface DiagnosisOptions {
  equipment: any;
  description: string;
  priority?: string;
  ticketId?: number | string;
  files?: File[];
  onProgress?: (message: string) => void;
}

/**
 * Process files (images/videos) for AI analysis
 * Returns array of base64-encoded images
 */
export async function processFilesForAI(
  files: File[],
  onProgress?: (message: string) => void
): Promise<string[]> {
  const imageData: string[] = [];
  let imageCount = 0;
  let videoCount = 0;

  onProgress?.('Preparing files for analysis...');

  for (const file of files) {
    if (isImageFile(file)) {
      imageCount++;
      onProgress?.(Processing image ${imageCount}/${files.length}...);
      try {
        const base64 = await compressImage(file);
        imageData.push(base64);
      } catch (error) {
        console.error(Error processing image ${file.name}:, error);
      }
    } else if (isVideoFile(file)) {
      videoCount++;
      onProgress?.(Extracting frames from video ${videoCount}...);
      try {
        const frames = await extractVideoFrames(file, 3);
        imageData.push(...frames);
      } catch (error) {
        console.error(Error extracting video frames from ${file.name}:, error);
      }
    }
  }

  return imageData;
}

/**
 * Run AI diagnosis with optional file attachments
 */
export async function runAIDiagnosis(options: DiagnosisOptions): Promise<any> {
  const { equipment, description, priority, ticketId, files, onProgress } = options;

  try {
    onProgress?.('Starting AI diagnosis...');

    // Step 1: Process files if provided
    let imageData: string[] = [];
    
    if (files && files.length > 0) {
      onProgress?.(Processing ${files.length} file(s)...);
      imageData = await processFilesForAI(files, onProgress);
    } else if (ticketId) {
      // Fetch attachments from ticket
      onProgress?.('Fetching ticket attachments...');
      try {
        const attachments = await attachmentsApi.getByEntity(ticketId, 'ticket');
        
        if (attachments && attachments.length > 0) {
          onProgress?.(Found ${attachments.length} attachment(s)...);
          // Note: Backend needs to provide image URLs or base64 data
          // For now, we'll skip processing existing attachments
        }
      } catch (err) {
        console.log('No attachments found or error fetching:', err);
      }
    }

    // Step 2: Run diagnosis
    const analysisText = imageData.length > 0
      ? Analyzing ${imageData.length} image(s) with AI...
      : 'Analyzing with AI...';
    onProgress?.(analysisText);

    const diagnosisPayload = {
      equipment_id: equipment?.id,
      description: description,
      priority: priority || 'medium',
      symptoms: [description],
      equipment_context: {
        manufacturer: equipment?.manufacturer_name || equipment?.manufacturer,
        model: equipment?.model_name || equipment?.model,
        serial_number: equipment?.serial_number,
        installation_date: equipment?.installation_date,
      },
      // Include image data if available
      ...(imageData.length > 0 && { images: imageData }),
    };

    const diagnosis = await diagnosisApi.diagnose(diagnosisPayload);

    onProgress?.('Analysis complete!');

    return diagnosis;
  } catch (error) {
    console.error('AI Diagnosis failed:', error);
    onProgress?.('Analysis failed');
    throw error;
  }
}

/**
 * Add diagnosis as a comment to the ticket
 */
export async function addDiagnosisComment(ticketId: number | string, diagnosis: any): Promise<void> {
  try {
    const commentsApi = (await import('@/lib/api/tickets')).ticketsApi;
    
    // Format diagnosis results as comment
    let commentText = '🤖 **AI Diagnosis Results**\n\n';
    
    if (diagnosis.primary_diagnosis) {
      commentText += **Problem:** ${diagnosis.primary_diagnosis.problem_type}\n;
      commentText += **Description:** ${diagnosis.primary_diagnosis.description}\n;
      commentText += **Confidence:** ${Math.round(diagnosis.primary_diagnosis.confidence * 100)}%\n\n;
    }

    if (diagnosis.vision_analysis && diagnosis.vision_analysis.findings?.length > 0) {
      commentText += '**Visual Findings:**\n';
      diagnosis.vision_analysis.findings.forEach((finding: any, idx: number) => {
        commentText += ${idx + 1}. ${finding.finding} (${Math.round(finding.confidence * 100)}%)\n;
      });
      commentText += '\n';
    }

    if (diagnosis.recommended_actions && diagnosis.recommended_actions.length > 0) {
      commentText += '**Recommended Actions:**\n';
      diagnosis.recommended_actions.forEach((action: any, idx: number) => {
        commentText += ${idx + 1}. ${action.action} - ${action.description}\n;
      });
      commentText += '\n';
    }

    if (diagnosis.required_parts && diagnosis.required_parts.length > 0) {
      commentText += '**Required Parts:**\n';
      diagnosis.required_parts.forEach((part: any) => {
        commentText += - ${part.part_name} (${part.part_code}) - ${Math.round(part.probability * 100)}% probability\n;
      });
    }

    // Add comment via tickets API
    await commentsApi.addComment(ticketId, {
      comment_text: commentText,
      created_by: 'AI Diagnosis System',
    });

    console.log('✅ AI diagnosis added as comment');
  } catch (error) {
    console.error('Failed to add diagnosis comment:', error);
    throw error;
  }
}

/**
 * Suggest engineers based on diagnosis
 */
export function suggestEngineersFromDiagnosis(diagnosis: any): {
  specialization: string;
  reasoning: string;
} | null {
  if (!diagnosis.primary_diagnosis) return null;

  const category = diagnosis.primary_diagnosis.problem_category;
  const type = diagnosis.primary_diagnosis.problem_type;

  // Map problem categories to specializations
  const specializationMap: Record<string, string> = {
    'Electrical': 'electrical',
    'Mechanical': 'mechanical',
    'Software': 'software',
    'Hydraulic': 'hydraulic',
    'Pneumatic': 'mechanical',
    'Electronic': 'electrical',
  };

  const specialization = specializationMap[category] || 'general';

  return {
    specialization,
    reasoning: Based on ${category} issue: ${type},
  };
}

/**
 * Get parts suggestions from diagnosis
 */
export function getPartsSuggestionsFromDiagnosis(diagnosis: any): Array<{
  part_code: string;
  part_name: string;
  quantity: number;
  probability: number;
}> {
  if (!diagnosis.required_parts) return [];

  return diagnosis.required_parts.map((part: any) => ({
    part_code: part.part_code,
    part_name: part.part_name,
    quantity: part.quantity || 1,
    probability: part.probability || 0,
  }));
}
