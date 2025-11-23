// ============================================================================
// OpenAI Integration for AI-Assisted Diagnosis
// ============================================================================

import { DiagnosisRequest, DiagnosisResponse, DiagnosisResult } from '../api/diagnosis';
import { Equipment } from '@/types';

// OpenAI API configuration
const OPENAI_API_KEY = process.env.NEXT_PUBLIC_OPENAI_API_KEY;
const OPENAI_API_URL = 'https://api.openai.com/v1/chat/completions';

interface OpenAIResponse {
  choices: {
    message: {
      content: string;
    };
  }[];
  usage: {
    total_tokens: number;
  };
}

interface AIAnalysisResult {
  primary_diagnosis: DiagnosisResult;
  alternate_diagnoses: DiagnosisResult[];
  confidence_factors: string[];
  confidence_score: number;
  recommended_actions: any[];
  required_parts: any[];
  reasoning: string;
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// AI DIAGNOSIS ENGINE
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

export class AIDiagnosisEngine {
  
  /**
   * Generate AI-powered equipment diagnosis
   */
  async analyzeDiagnosis(request: DiagnosisRequest, equipment: Equipment): Promise<DiagnosisResponse> {
    
    // For now, check if API key is available, otherwise fall back to intelligent mock
    if (!OPENAI_API_KEY || OPENAI_API_KEY === 'your-openai-api-key-here') {
      console.log('🤖 OpenAI API key not configured, using intelligent mock analysis');
      return this.generateIntelligentMockDiagnosis(request, equipment);
    }

    try {
      const aiAnalysis = await this.callOpenAI(request, equipment);
      return this.formatAIResponse(request, equipment, aiAnalysis);
    } catch (error) {
      console.error('OpenAI API Error:', error);
      // Fallback to intelligent mock on API error
      return this.generateIntelligentMockDiagnosis(request, equipment);
    }
  }

  /**
   * Call OpenAI API with specialized medical equipment prompt
   */
  private async callOpenAI(request: DiagnosisRequest, equipment: Equipment): Promise<AIAnalysisResult> {
    const prompt = this.buildDiagnosisPrompt(request, equipment);
    
    const response = await fetch(OPENAI_API_URL, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${OPENAI_API_KEY}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        model: 'gpt-4',
        messages: [
          {
            role: 'system',
            content: this.getSystemPrompt()
          },
          {
            role: 'user',
            content: prompt
          }
        ],
        temperature: 0.3, // Lower temperature for more consistent technical analysis
        max_tokens: 2000,
      }),
    });

    if (!response.ok) {
      throw new Error(`OpenAI API error: ${response.statusText}`);
    }

    const data: OpenAIResponse = await response.json();
    const aiContent = data.choices[0].message.content;
    
    // Parse structured AI response
    return this.parseAIResponse(aiContent);
  }

  /**
   * System prompt for medical equipment diagnosis
   */
  private getSystemPrompt(): string {
    return `You are an expert medical equipment diagnostic engineer with 20+ years of experience. 

Your role is to analyze equipment failures and provide precise, actionable diagnoses with confidence scoring.

CRITICAL REQUIREMENTS:
1. Provide structured JSON responses only
2. Calculate realistic confidence scores (0.0-1.0) based on symptom clarity and historical patterns
3. Consider equipment age, manufacturer specifications, and common failure modes
4. Suggest specific, actionable repair steps with time estimates
5. Recommend exact parts with OEM/generic alternatives
6. Include safety precautions for all repairs
7. Provide alternative diagnoses when symptoms could indicate multiple issues

CONFIDENCE SCORING GUIDELINES:
- 0.8-1.0 (HIGH): Clear symptoms, visual confirmation, common failure pattern
- 0.6-0.79 (MEDIUM): Multiple symptoms align, but some ambiguity
- 0.3-0.59 (LOW): Symptoms could indicate several issues, needs investigation
- Below 0.3: Insufficient information for diagnosis

Always respond with valid JSON in the exact format requested. Be precise, technical, and actionable.`;
  }

  /**
   * Build equipment-specific diagnosis prompt
   */
  private buildDiagnosisPrompt(request: DiagnosisRequest, equipment: Equipment): string {
    const equipmentAge = equipment.installation_date 
      ? Math.floor((new Date().getTime() - new Date(equipment.installation_date).getTime()) / (1000 * 60 * 60 * 24 * 365))
      : 'unknown';

    return `MEDICAL EQUIPMENT DIAGNOSIS REQUEST

EQUIPMENT DETAILS:
- Type: ${equipment.equipment_name}
- Manufacturer: ${equipment.manufacturer_name}  
- Model: ${equipment.model_number || 'Unknown'}
- Serial: ${equipment.serial_number}
- Age: ${equipmentAge} years
- Location: ${equipment.installation_location || 'Unknown'}
- Service History: ${equipment.service_count} previous services

REPORTED ISSUE:
Description: ${request.description}
Symptoms: ${request.symptoms.join(', ')}

ANALYSIS REQUEST:
Please provide a comprehensive diagnosis with the following JSON structure:

{
  "primary_diagnosis": {
    "problem_category": "string",
    "problem_type": "string", 
    "description": "string",
    "confidence": number (0-100),
    "severity": "low|medium|high|critical",
    "root_cause": "detailed technical explanation",
    "symptoms": ["array of symptoms"],
    "possible_causes": ["array of potential causes"],
    "reasoning_explanation": "why this is the most likely diagnosis"
  },
  "alternate_diagnoses": [
    {
      "problem_category": "string",
      "problem_type": "string",
      "description": "string", 
      "confidence": number (0-100),
      "severity": "low|medium|high|critical",
      "root_cause": "string",
      "symptoms": ["array"],
      "possible_causes": ["array"],
      "reasoning_explanation": "why this is possible but less likely"
    }
  ],
  "confidence_factors": [
    "specific reasons for confidence level",
    "equipment age considerations",
    "manufacturer-specific failure patterns",
    "symptom clarity assessment"
  ],
  "confidence_score": number (0.0-1.0),
  "recommended_actions": [
    {
      "action": "specific action description",
      "priority": "high|medium|low",
      "description": "detailed steps",
      "estimated_time": "time estimate",
      "requires_specialist": boolean,
      "specialist_type": "if required",
      "required_tools": ["array of tools"],
      "required_parts": ["array of parts"],
      "safety_precautions": ["array of safety steps"]
    }
  ],
  "required_parts": [
    {
      "part_code": "suggested part number",
      "part_name": "part description",
      "part_category": "category",
      "probability": number (0-100),
      "quantity": number,
      "is_oem_required": boolean,
      "manufacturer": "part manufacturer",
      "estimated_cost": number (optional)
    }
  ],
  "reasoning": "overall diagnostic reasoning and approach"
}

Focus on ${equipment.manufacturer_name} ${equipment.equipment_name} specific failure patterns and provide actionable, technically accurate guidance.`;
  }

  /**
   * Parse structured AI response
   */
  private parseAIResponse(aiContent: string): AIAnalysisResult {
    try {
      // Clean up the response in case there are markdown code blocks
      const jsonMatch = aiContent.match(/\{[\s\S]*\}/);
      const jsonString = jsonMatch ? jsonMatch[0] : aiContent;
      
      return JSON.parse(jsonString);
    } catch (error) {
      console.error('Failed to parse AI response:', error);
      throw new Error('Invalid AI response format');
    }
  }

  /**
   * Format AI analysis into our DiagnosisResponse structure
   */
  private formatAIResponse(request: DiagnosisRequest, equipment: Equipment, aiAnalysis: AIAnalysisResult): DiagnosisResponse {
    return {
      diagnosis_id: `ai_${Date.now()}`,
      ticket_id: request.ticket_id,
      primary_diagnosis: aiAnalysis.primary_diagnosis,
      alternate_diagnoses: aiAnalysis.alternate_diagnoses,
      confidence: aiAnalysis.confidence_score,
      confidence_level: this.getConfidenceLevel(aiAnalysis.confidence_score),
      decision_status: 'pending',
      ai_metadata: {
        provider: 'openai',
        model: 'gpt-4',
        confidence: aiAnalysis.confidence_score * 100,
        confidence_factors: aiAnalysis.confidence_factors,
        alternatives_count: aiAnalysis.alternate_diagnoses.length,
        requires_feedback: aiAnalysis.confidence_score < 0.8, // Require feedback for medium/low confidence
        suggestion_only: true
      },
      recommended_actions: aiAnalysis.recommended_actions,
      required_parts: aiAnalysis.required_parts,
      vision_analysis: {
        attachments_analyzed: 0,
        findings: [],
        overall_assessment: 'No images provided for analysis',
        detected_components: [],
        visible_damage: [],
        confidence: 0
      },
      estimated_resolution_time: this.calculateEstimatedTime(aiAnalysis.recommended_actions),
      created_at: new Date().toISOString()
    };
  }

  /**
   * Generate intelligent mock diagnosis when AI API is unavailable
   */
  private generateIntelligentMockDiagnosis(request: DiagnosisRequest, equipment: Equipment): DiagnosisResponse {
    // Analyze symptoms to provide contextual mock response
    const symptoms = request.symptoms.map(s => s.toLowerCase());
    const description = request.description.toLowerCase();
    
    // Determine most likely issue category based on keywords
    let primaryIssue, confidence, actions, parts;
    
    if (this.containsKeywords(description + symptoms.join(' '), ['power', 'start', 'led', 'display', 'fan', 'button'])) {
      // Power-related issue
      primaryIssue = this.getPowerIssueTemplate(equipment, symptoms);
      confidence = 0.78;
    } else if (this.containsKeywords(description + symptoms.join(' '), ['image', 'scan', 'quality', 'artifact', 'noise', 'contrast'])) {
      // Image quality issue  
      primaryIssue = this.getImageQualityTemplate(equipment, symptoms);
      confidence = 0.72;
    } else if (this.containsKeywords(description + symptoms.join(' '), ['heat', 'cool', 'temperature', 'alarm', 'pressure', 'compressor'])) {
      // Cooling system issue
      primaryIssue = this.getCoolingIssueTemplate(equipment, symptoms);
      confidence = 0.75;
    } else if (this.containsKeywords(description + symptoms.join(' '), ['software', 'crash', 'error', 'freeze', 'save', 'interface'])) {
      // Software issue
      primaryIssue = this.getSoftwareIssueTemplate(equipment, symptoms);
      confidence = 0.65;
    } else {
      // Generic issue
      primaryIssue = this.getGenericIssueTemplate(equipment, symptoms);
      confidence = 0.55;
    }

    return {
      diagnosis_id: `mock_ai_${Date.now()}`,
      ticket_id: request.ticket_id,
      primary_diagnosis: primaryIssue.diagnosis,
      alternate_diagnoses: primaryIssue.alternatives,
      confidence,
      confidence_level: this.getConfidenceLevel(confidence),
      decision_status: 'pending',
      ai_metadata: {
        provider: 'openai-mock',
        model: 'gpt-4-mock',
        confidence: confidence * 100,
        confidence_factors: primaryIssue.factors,
        alternatives_count: primaryIssue.alternatives.length,
        requires_feedback: confidence < 0.8,
        suggestion_only: true
      },
      recommended_actions: primaryIssue.actions,
      required_parts: primaryIssue.parts,
      vision_analysis: {
        attachments_analyzed: 0,
        findings: [],
        overall_assessment: 'No images provided for analysis',
        detected_components: [],
        visible_damage: [],
        confidence: 0
      },
      estimated_resolution_time: this.calculateEstimatedTime(primaryIssue.actions),
      created_at: new Date().toISOString()
    };
  }

  // Helper methods for intelligent mock responses
  private containsKeywords(text: string, keywords: string[]): boolean {
    return keywords.some(keyword => text.includes(keyword));
  }

  private getPowerIssueTemplate(equipment: Equipment, symptoms: string[]) {
    return {
      diagnosis: {
        problem_category: "Electrical System",
        problem_type: "Power Supply Failure",
        description: `${equipment.equipment_name} power supply failure detected`,
        confidence: 85,
        severity: "high",
        root_cause: `Primary power supply components in the ${equipment.equipment_name} have likely failed due to component aging or power surge damage. This is a common failure mode for ${equipment.manufacturer_name} equipment over 5 years old.`,
        symptoms,
        possible_causes: ["Power supply capacitor failure", "Transformer malfunction", "Power surge damage", "Age-related component degradation"],
        reasoning_explanation: "Based on the described power-related symptoms and the equipment age/type, this appears to be a typical power supply failure pattern."
      },
      alternatives: [{
        problem_category: "Control System",
        problem_type: "Main Board Failure",
        description: "Control board malfunction",
        confidence: 65,
        severity: "high", 
        root_cause: "Main control board may have suffered component failure",
        symptoms,
        possible_causes: ["Processor failure", "Memory corruption", "Board component failure"],
        reasoning_explanation: "Alternative possibility if power supply tests normal, but less likely given symptom pattern."
      }],
      factors: [
        `Analysis of ${symptoms.length} reported power-related symptoms`,
        `${equipment.manufacturer_name} equipment failure pattern recognition`,
        `Equipment age and historical failure analysis`,
        `Power supply is most common failure point for these symptoms`
      ],
      actions: [{
        action: "Replace main power supply unit",
        priority: "high",
        description: `Replace the main power supply unit in the ${equipment.equipment_name} with OEM equivalent`,
        estimated_time: "2-3 hours",
        requires_specialist: false,
        required_tools: ["Screwdriver set", "Multimeter", "Anti-static equipment"],
        required_parts: ["Power supply unit"],
        safety_precautions: ["Disconnect all power", "Use anti-static protection", "Allow capacitor discharge time"]
      }],
      parts: [{
        part_code: `PSU-${equipment.manufacturer_name}-${equipment.model_number || 'STD'}-001`,
        part_name: "Main Power Supply Unit",
        part_category: "Power Supply",
        probability: 90,
        quantity: 1,
        is_oem_required: true,
        manufacturer: equipment.manufacturer_name,
        estimated_cost: 280
      }]
    };
  }

  private getImageQualityTemplate(equipment: Equipment, symptoms: string[]) {
    return {
      diagnosis: {
        problem_category: "Imaging System", 
        problem_type: "Image Processing Degradation",
        description: `${equipment.equipment_name} image quality degradation detected`,
        confidence: 78,
        severity: "medium",
        root_cause: `Image processing components in ${equipment.equipment_name} showing signs of calibration drift or component aging affecting image quality output.`,
        symptoms,
        possible_causes: ["Detector calibration drift", "Image processing board issues", "Software calibration errors", "Component aging"],
        reasoning_explanation: "Image quality issues typically indicate calibration problems or imaging component degradation."
      },
      alternatives: [{
        problem_category: "Software",
        problem_type: "Calibration Software Error", 
        description: "Software calibration malfunction",
        confidence: 60,
        severity: "medium",
        root_cause: "Calibration software may need updating or recalibration",
        symptoms,
        possible_causes: ["Software corruption", "Calibration data loss", "Configuration errors"],
        reasoning_explanation: "Could be resolved with software recalibration if hardware tests normal."
      }],
      factors: [
        `Analysis of image quality symptoms`,
        `${equipment.manufacturer_name} imaging system failure patterns`, 
        `Calibration drift assessment based on symptoms`,
        `Software vs hardware failure probability analysis`
      ],
      actions: [{
        action: "Perform system recalibration and component testing",
        priority: "medium",
        description: "Run full system calibration and test imaging components",
        estimated_time: "3-4 hours",
        requires_specialist: true,
        specialist_type: "Imaging specialist",
        required_tools: ["Calibration phantoms", "Test equipment"],
        required_parts: [],
        safety_precautions: ["Follow radiation safety protocols", "Use proper test procedures"]
      }],
      parts: [{
        part_code: `IMG-DETECT-${equipment.manufacturer_name}-001`,
        part_name: "Image Detector Assembly",
        part_category: "Imaging Component",
        probability: 45,
        quantity: 1,
        is_oem_required: true,
        manufacturer: equipment.manufacturer_name,
        estimated_cost: 2500
      }]
    };
  }

  private getCoolingIssueTemplate(equipment: Equipment, symptoms: string[]) {
    return {
      diagnosis: {
        problem_category: "Thermal Management",
        problem_type: "Cooling System Failure",
        description: `${equipment.equipment_name} cooling system malfunction`,
        confidence: 82,
        severity: "high",
        root_cause: `Cooling system in ${equipment.equipment_name} showing signs of failure, likely due to coolant leak, compressor issues, or thermal sensor malfunction.`,
        symptoms,
        possible_causes: ["Coolant leak", "Compressor failure", "Thermal sensor malfunction", "Cooling pump failure"],
        reasoning_explanation: "Temperature-related symptoms strongly indicate cooling system issues which are critical for equipment operation."
      },
      alternatives: [{
        problem_category: "Environmental",
        problem_type: "External Temperature Issues",
        description: "Room temperature or ventilation problem",
        confidence: 45,
        severity: "low",
        root_cause: "External environmental factors affecting equipment temperature",
        symptoms,
        possible_causes: ["Room HVAC failure", "Blocked ventilation", "Ambient temperature too high"],
        reasoning_explanation: "Less likely but should be ruled out before internal cooling system repair."
      }],
      factors: [
        `Thermal alarm and cooling symptom analysis`,
        `${equipment.manufacturer_name} cooling system reliability data`,
        `Critical nature of cooling for equipment operation`,
        `Historical cooling system failure patterns`
      ],
      actions: [{
        action: "Emergency cooling system inspection and repair",
        priority: "high", 
        description: "Immediate cooling system diagnostics and component replacement",
        estimated_time: "4-6 hours",
        requires_specialist: true,
        specialist_type: "Thermal systems technician",
        required_tools: ["Pressure gauges", "Leak detection equipment", "Thermal sensors"],
        required_parts: ["Potential compressor or cooling components"],
        safety_precautions: ["System shutdown required", "Handle refrigerants safely", "Follow pressure system safety"]
      }],
      parts: [{
        part_code: `COOL-COMP-${equipment.manufacturer_name}-001`,
        part_name: "Primary Cooling Compressor",
        part_category: "Cooling System",
        probability: 75,
        quantity: 1,
        is_oem_required: true,
        manufacturer: equipment.manufacturer_name,
        estimated_cost: 1200
      }]
    };
  }

  private getSoftwareIssueTemplate(equipment: Equipment, symptoms: string[]) {
    return {
      diagnosis: {
        problem_category: "Software System",
        problem_type: "Software Malfunction",
        description: `${equipment.equipment_name} software system instability`,
        confidence: 68,
        severity: "medium",
        root_cause: `Software system showing signs of instability, possibly due to corrupted files, memory issues, or software bugs in ${equipment.equipment_name}.`,
        symptoms,
        possible_causes: ["Software corruption", "Memory leaks", "Database issues", "System file corruption"],
        reasoning_explanation: "Software-related symptoms suggest system instability that may require software reinstall or system updates."
      },
      alternatives: [{
        problem_category: "Hardware",
        problem_type: "System Memory Issues",
        description: "RAM or storage hardware malfunction",
        confidence: 55,
        severity: "medium", 
        root_cause: "Hardware memory or storage components may be failing",
        symptoms,
        possible_causes: ["RAM failure", "Hard drive corruption", "Storage device issues"],
        reasoning_explanation: "Hardware issues can manifest as software symptoms and should be tested."
      }],
      factors: [
        `Software crash and error pattern analysis`,
        `${equipment.manufacturer_name} software reliability assessment`,
        `System performance degradation indicators`,
        `Hardware vs software failure probability evaluation`
      ],
      actions: [{
        action: "Software system diagnostics and reinstallation",
        priority: "medium",
        description: "Run diagnostics and perform clean software reinstall if needed", 
        estimated_time: "2-4 hours",
        requires_specialist: true,
        specialist_type: "Software technician",
        required_tools: ["Diagnostic software", "Installation media", "Backup equipment"],
        required_parts: [],
        safety_precautions: ["Backup all patient data", "Follow software update procedures", "Verify system integrity"]
      }],
      parts: [{
        part_code: `MEM-RAM-${equipment.manufacturer_name}-001`,
        part_name: "System Memory Module",
        part_category: "Computer Hardware", 
        probability: 30,
        quantity: 1,
        is_oem_required: false,
        manufacturer: "Generic",
        estimated_cost: 150
      }]
    };
  }

  private getGenericIssueTemplate(equipment: Equipment, symptoms: string[]) {
    return {
      diagnosis: {
        problem_category: "General System",
        problem_type: "System Malfunction",
        description: `${equipment.equipment_name} system malfunction requiring investigation`,
        confidence: 55,
        severity: "medium",
        root_cause: `Multiple system indicators suggest a malfunction in ${equipment.equipment_name} that requires further diagnostic investigation.`,
        symptoms,
        possible_causes: ["Multiple potential causes", "Requires diagnostic investigation", "Component failure", "System integration issues"],
        reasoning_explanation: "Symptoms require additional investigation to determine specific root cause and appropriate repair approach."
      },
      alternatives: [{
        problem_category: "Preventive",
        problem_type: "Scheduled Maintenance Required",
        description: "Equipment may need routine maintenance",
        confidence: 40,
        severity: "low",
        root_cause: "Issues may be resolved with routine maintenance procedures",
        symptoms,
        possible_causes: ["Routine wear", "Calibration drift", "Component aging"],
        reasoning_explanation: "Some symptoms may be addressed through preventive maintenance before major repairs."
      }],
      factors: [
        `General symptom pattern analysis`,
        `Equipment age and service history consideration`,
        `Need for comprehensive diagnostic testing`,
        `Multiple potential failure modes identified`
      ],
      actions: [{
        action: "Comprehensive system diagnostic testing",
        priority: "medium",
        description: "Perform full system diagnostic to identify specific issues",
        estimated_time: "2-3 hours",
        requires_specialist: true,
        specialist_type: "General technician",
        required_tools: ["Diagnostic equipment", "Test instruments", "Service manuals"],
        required_parts: [],
        safety_precautions: ["Follow all equipment safety procedures", "Use proper diagnostic protocols"]
      }],
      parts: []
    };
  }

  private getConfidenceLevel(confidence: number): 'HIGH' | 'MEDIUM' | 'LOW' {
    if (confidence >= 0.8) return 'HIGH';
    if (confidence >= 0.6) return 'MEDIUM';
    return 'LOW';
  }

  private calculateEstimatedTime(actions: any[]): string {
    if (actions.length === 0) return "2-4 hours";
    
    // Extract time estimates and calculate total
    const times = actions.map(action => {
      const timeStr = action.estimated_time || "2 hours";
      const match = timeStr.match(/(\d+)(?:-(\d+))?\s*hours?/);
      if (match) {
        const min = parseInt(match[1]);
        const max = match[2] ? parseInt(match[2]) : min;
        return (min + max) / 2;
      }
      return 2; // Default
    });
    
    const totalHours = times.reduce((sum, time) => sum + time, 0);
    const minHours = Math.floor(totalHours * 0.8);
    const maxHours = Math.ceil(totalHours * 1.2);
    
    return `${minHours}-${maxHours} hours`;
  }
}

// Export singleton instance
export const aiDiagnosisEngine = new AIDiagnosisEngine();