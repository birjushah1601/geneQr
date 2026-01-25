// AI Onboarding Helper Functions
import apiClient from './api/client';

export interface OnboardingStep {
  id: string;
  completed: boolean;
  data?: any;
}

export interface OnboardingSession {
  sessionId: string;
  manufacturerId?: string;
  currentStep: number;
  steps: OnboardingStep[];
  createdAt: string;
  updatedAt: string;
}

// Template for AI conversation flow
export const ONBOARDING_TEMPLATES = {
  company_profile: {
    questions: [
      "What's your company's registered name?",
      "What's your GSTIN number? (This helps verify your business)",
      "What's your primary email address for communication?",
      "What's your contact phone number?",
      "Where is your head office located? (City, State)",
    ],
    validation: {
      gstin: /^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$/,
      email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
      phone: /^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$/,
    },
  },
  equipment_types: {
    questions: [
      "What type of medical equipment do you manufacture?",
      "What are your main product models?",
      "What's the typical warranty period for your equipment?",
    ],
    suggestions: [
      "Imaging Equipment (MRI, CT, X-Ray)",
      "Cardiology Equipment (ECG, Monitors)",
      "Laboratory Equipment (Analyzers)",
      "Surgical Equipment (OR Tables, Lights)",
      "ICU Equipment (Ventilators, Monitors)",
    ],
  },
  parts_catalog: {
    questions: [
      "Do you have a parts catalog ready?",
      "Upload your parts list or I can help you create one",
    ],
    csvColumns: {
      required: ['part_number', 'part_name', 'category', 'unit_price'],
      optional: ['description', 'stock_quantity', 'supplier', 'lead_time_days'],
    },
  },
  hospitals: {
    questions: [
      "Do you have existing customers/hospitals?",
      "Upload your customer list or add them one by one",
    ],
    csvColumns: {
      required: ['name', 'city', 'state'],
      optional: ['phone', 'email', 'contact_person', 'gstin'],
    },
  },
  installations: {
    questions: [
      "Do you have equipment already installed at hospitals?",
      "Upload your equipment registry with serial numbers and locations",
    ],
    csvColumns: {
      required: ['hospital_name', 'equipment_name', 'serial_number', 'installation_date'],
      optional: ['model', 'warranty_end_date', 'location', 'department'],
    },
  },
  engineers: {
    questions: [
      "How many service engineers do you have?",
      "Upload your engineers list with their expertise and coverage areas",
    ],
    csvColumns: {
      required: ['name', 'phone', 'email', 'engineer_level'],
      optional: ['equipment_types', 'experience_years', 'location', 'coverage_cities'],
    },
  },
};

// Parse CSV file
export async function parseCSVFile(file: File): Promise<any[]> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    
    reader.onload = (e) => {
      try {
        const text = e.target?.result as string;
        const lines = text.split('\n').filter(line => line.trim());
        
        if (lines.length < 2) {
          reject(new Error('CSV file is empty or has no data rows'));
          return;
        }

        // Parse header
        const headers = lines[0].split(',').map(h => h.trim().toLowerCase());
        
        // Parse rows
        const data = lines.slice(1).map((line, index) => {
          const values = line.split(',').map(v => v.trim());
          const row: any = { _rowNumber: index + 2 }; // +2 because header is row 1
          
          headers.forEach((header, i) => {
            row[header] = values[i] || '';
          });
          
          return row;
        });

        resolve(data);
      } catch (error) {
        reject(error);
      }
    };

    reader.onerror = () => reject(new Error('Failed to read file'));
    reader.readAsText(file);
  });
}

// Validate CSV data against template
export function validateCSVData(data: any[], template: typeof ONBOARDING_TEMPLATES[keyof typeof ONBOARDING_TEMPLATES]): {
  valid: boolean;
  errors: string[];
  warnings: string[];
  rowCount: number;
} {
  const errors: string[] = [];
  const warnings: string[] = [];

  if (!template.csvColumns) {
    return { valid: true, errors, warnings, rowCount: data.length };
  }

  const headers = Object.keys(data[0] || {}).filter(k => k !== '_rowNumber');
  const missingRequired = template.csvColumns.required.filter(col => !headers.includes(col));

  if (missingRequired.length > 0) {
    errors.push(`Missing required columns: ${missingRequired.join(', ')}`);
  }

  // Validate each row
  data.forEach((row, index) => {
    template.csvColumns!.required.forEach(col => {
      if (!row[col] || row[col].trim() === '') {
        errors.push(`Row ${row._rowNumber || index + 2}: Missing required field '${col}'`);
      }
    });
  });

  return {
    valid: errors.length === 0,
    errors,
    warnings,
    rowCount: data.length,
  };
}

// Call GPT-4/Claude API for intelligent responses
export async function getAIResponse(userMessage: string, context: {
  currentStep: string;
  conversationHistory: string[];
  onboardingData: any;
}): Promise<string> {
  try {
    // In production, call OpenAI API here
    // For now, return template-based responses
    
    const step = context.currentStep;
    const template = ONBOARDING_TEMPLATES[step as keyof typeof ONBOARDING_TEMPLATES];
    
    if (!template) {
      return "I'm not sure how to help with that. Can you provide more details?";
    }

    // Simple keyword matching (replace with GPT-4 in production)
    if (userMessage.toLowerCase().includes('skip') || userMessage.toLowerCase().includes('later')) {
      return `No problem! We can ${step.replace('_', ' ')} later. Let's move to the next step.`;
    }

    if (userMessage.toLowerCase().includes('help') || userMessage.toLowerCase().includes('what')) {
      return template.questions[0] + '\n\n' + (
        'suggestions' in template
          ? `Common options:\n${(template.suggestions as string[]).map((s, i) => `${i + 1}. ${s}`).join('\n')}`
          : ''
      );
    }

    return template.questions[0];
  } catch (error) {
    console.error('AI API error:', error);
    throw error;
  }
}

// Extract structured data from natural language (using GPT-4)
export async function extractStructuredData(text: string, dataType: 'company' | 'equipment' | 'engineer'): Promise<any> {
  // In production, use GPT-4 to extract structured data
  // For now, return simple parsing
  
  switch (dataType) {
    case 'company':
      return {
        name: text,
        organization_type: 'manufacturer',
        status: 'active',
      };
    
    case 'equipment':
      return {
        name: text,
        category: 'medical_equipment',
      };
    
    case 'engineer':
      return {
        name: text,
        status: 'active',
      };
    
    default:
      return { raw: text };
  }
}

// Smart field mapping for uploaded CSVs
export function mapCSVFields(headers: string[], targetFields: string[]): Record<string, string> {
  const mapping: Record<string, string> = {};

  // Simple fuzzy matching (in production, use GPT-4 for better mapping)
  headers.forEach(header => {
    const normalized = header.toLowerCase().replace(/[^a-z0-9]/g, '');
    
    targetFields.forEach(target => {
      const targetNormalized = target.toLowerCase().replace(/[^a-z0-9]/g, '');
      
      if (normalized.includes(targetNormalized) || targetNormalized.includes(normalized)) {
        mapping[header] = target;
      }
    });
  });

  return mapping;
}

// Import data via existing APIs
export async function importData(dataType: string, data: any[], sessionId: string) {
  try {
    let endpoint = '';
    let payload: any = {};

    switch (dataType) {
      case 'hospitals':
        endpoint = '/v1/organizations/import';
        // Format for organizations import
        payload = {
          organizations: data.map(row => ({
            name: row.name || row.hospital_name,
            organization_type: 'hospital',
            city: row.city,
            state: row.state,
            phone: row.phone,
            email: row.email,
            contact_person: row.contact_person,
          })),
        };
        break;

      case 'equipment':
        endpoint = '/v1/equipment/import';
        // Use existing equipment import
        break;

      case 'engineers':
        endpoint = '/v1/engineers/import';
        // Use existing engineer import
        break;

      case 'parts':
        endpoint = '/v1/parts/import';
        // TODO: Create parts import endpoint
        break;

      default:
        throw new Error(`Unknown data type: ${dataType}`);
    }

    const response = await apiClient.post(endpoint, payload);
    return response.data;
  } catch (error) {
    console.error(`Import error for ${dataType}:`, error);
    throw error;
  }
}

// Save onboarding session progress
export function saveOnboardingSession(session: OnboardingSession): void {
  localStorage.setItem('onboarding_session', JSON.stringify(session));
}

// Load onboarding session
export function loadOnboardingSession(): OnboardingSession | null {
  const saved = localStorage.getItem('onboarding_session');
  return saved ? JSON.parse(saved) : null;
}

// Clear onboarding session
export function clearOnboardingSession(): void {
  localStorage.removeItem('onboarding_session');
}
