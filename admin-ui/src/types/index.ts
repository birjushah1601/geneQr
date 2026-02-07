// ============================================================================
// ServQR Admin UI - TypeScript Type Definitions
// ============================================================================

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// COMMON TYPES
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export interface PaginationParams {
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_dir?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiError {
  error: string;
  details?: string;
  code?: string;
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// MANUFACTURER
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export interface Manufacturer {
  id: string;
  name: string;
  contact_person: string;
  email: string;
  phone: string;
  address?: string;
  website?: string;
  logo_url?: string;
  equipment_count?: number;
  engineer_count?: number;
  active_tickets?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateManufacturerRequest {
  name: string;
  contact_person: string;
  email: string;
  phone: string;
  address?: string;
  website?: string;
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// EQUIPMENT
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export type EquipmentStatus = 'operational' | 'down' | 'under_maintenance' | 'decommissioned';

export interface Equipment {
  id: string;
  qr_code: string;
  serial_number: string;
  equipment_id?: string;
  equipment_name: string;
  manufacturer_name: string;
  model_number?: string;
  category?: string;
  
  // Customer info
  customer_id?: string;
  customer_name: string;
  installation_location?: string;
  installation_address?: Record<string, any>;
  installation_date?: string;
  
  // Purchase info
  contract_id?: string;
  purchase_date?: string;
  purchase_price?: number;
  warranty_expiry?: string;
  amc_contract_id?: string;
  
  // Status
  status: EquipmentStatus;
  last_service_date?: string;
  next_service_date?: string;
  service_count: number;
  
  // Technical
  specifications?: Record<string, any>;
  photos?: string[];
  documents?: string[];
  qr_code_url: string;
  
  // Metadata
  notes?: string;
  created_at: string;
  updated_at: string;
  created_by: string;
}

export interface RegisterEquipmentRequest {
  serial_number: string;
  equipment_name: string;
  manufacturer_name: string;
  model_number?: string;
  category?: string;
  customer_name: string;
  customer_id?: string;
  installation_location?: string;
  installation_address?: Record<string, any>;
  installation_date?: string;
  purchase_date?: string;
  purchase_price?: number;
  warranty_months?: number;
  amc_contract_id?: string;
  specifications?: Record<string, any>;
  notes?: string;
  created_by: string;
}

export interface EquipmentListParams extends PaginationParams {
  customer_id?: string;
  manufacturer?: string;
  category?: string;
  status?: EquipmentStatus;
  has_amc?: boolean;
  under_warranty?: boolean;
}

export interface CSVImportResult {
  total_rows: number;
  success_count: number;
  failure_count: number;
  errors: string[];
  imported_ids: string[];
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// FIELD ENGINEER
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export type EngineerStatus = 'active' | 'inactive' | 'on_leave';

export interface Engineer {
  id: string;
  name: string;
  phone: string;
  whatsapp?: string;
  email: string;
  
  // Location
  location: string;
  latitude?: number;
  longitude?: number;
  
  // Skills
  specializations: string[];
  certifications?: string[];
  experience_years?: number;
  
  // Assignment
  manufacturer_id?: string;
  manufacturer_name?: string;
  status: EngineerStatus;
  availability: 'available' | 'on_job' | 'off_duty';
  
  // Performance
  rating?: number;
  total_tickets?: number;
  completed_tickets?: number;
  avg_resolution_time?: number; // hours
  
  // Metadata
  photo_url?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEngineerRequest {
  name: string;
  phone: string;
  whatsapp?: string;
  email: string;
  location: string;
  latitude?: number;
  longitude?: number;
  specializations: string[];
  certifications?: string[];
  experience_years?: number;
  manufacturer_id: string; // REQUIRED: Must belong to a manufacturer
  employee_id?: string;
  notes?: string;
}

export interface EngineerListParams extends PaginationParams {
  manufacturer_id?: string;
  location?: string;
  status?: EngineerStatus;
  availability?: string;
  specialization?: string;
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// SERVICE TICKET
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export type TicketStatus = 'new' | 'assigned' | 'in_progress' | 'on_hold' | 'resolved' | 'closed' | 'cancelled';
export type TicketPriority = 'critical' | 'high' | 'medium' | 'low';
export type TicketSource = 'whatsapp' | 'web' | 'phone' | 'email' | 'scheduled';

export interface ServiceTicket {
  id: string;
  ticket_number: string;
  
  // Equipment & Customer
  equipment_id: string;
  qr_code?: string;
  serial_number: string;
  equipment_name: string;
  customer_id: string;
  customer_name: string;
  customer_phone: string;
  customer_email?: string;
  customer_whatsapp?: string;
  
  // Issue
  issue_category: string;
  issue_description: string;
  priority: TicketPriority;
  severity?: string;
  
  // Source
  source: TicketSource;
  source_message_id?: string;
  
  // Assignment
  assigned_engineer_id?: string;
  assigned_engineer_name?: string;
  assigned_at?: string;
  
  // Status & Timeline
  status: TicketStatus;
  created_at: string;
  acknowledged_at?: string;
  started_at?: string;
  resolved_at?: string;
  closed_at?: string;
  
  // SLA
  sla_response_due?: string;
  sla_resolution_due?: string;
  sla_breached: boolean;
  
  // Resolution
  resolution_notes?: string;
  parts_used?: Part[];
  labor_hours?: number;
  cost?: number;
  
  // Media
  photos?: string[];
  videos?: string[];
  documents?: string[];
  
  // Metadata
  notes?: string;
  created_by?: string;
  updated_at: string;
}

export interface Part {
  name: string;
  part_number?: string;
  quantity: number;
  cost?: number;
}

export interface CreateTicketRequest {
  equipment_id?: string;
  qr_code?: string;
  serial_number?: string;
  customer_phone: string;
  customer_whatsapp?: string;
  issue_category: string;
  issue_description: string;
  priority: TicketPriority;
  source: TicketSource;
  source_message_id?: string;
  notes?: string;
  created_by: string;
}

export interface AssignEngineerRequest {
  engineer_id: string;
  priority?: TicketPriority;
  sla_hours?: number;
  notes?: string;
}

export interface UpdateTicketStatusRequest {
  status: TicketStatus;
  notes?: string;
  resolution_notes?: string;
  parts_used?: Part[];
  labor_hours?: number;
  cost?: number;
}

export interface TicketListParams extends PaginationParams {
  status?: TicketStatus;
  priority?: TicketPriority;
  source?: TicketSource;
  engineer_id?: string;
  customer_id?: string;
  equipment_id?: string;
  date_from?: string;
  date_to?: string;
  sla_breached?: boolean;
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// WHATSAPP
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export interface WhatsAppMessage {
  id: string;
  from: string;
  to: string;
  text: string;
  timestamp: string;
  type: 'text' | 'image' | 'document' | 'location';
  media_url?: string;
}

export interface WhatsAppWebhookPayload {
  event: 'message' | 'status';
  message?: WhatsAppMessage;
  status?: {
    id: string;
    status: 'sent' | 'delivered' | 'read' | 'failed';
    timestamp: string;
  };
}

// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”
// DASHBOARD STATS
// â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”â”

export interface DashboardStats {
  equipment: {
    total: number;
    operational: number;
    down: number;
    under_maintenance: number;
    with_qr_codes: number;
  };
  tickets: {
    total: number;
    new: number;
    assigned: number;
    in_progress: number;
    resolved_today: number;
    sla_breached: number;
    avg_resolution_time: number; // hours
  };
  engineers: {
    total: number;
    available: number;
    on_job: number;
    off_duty: number;
    avg_rating: number;
  };
  manufacturers: {
    total: number;
    active: number;
  };
}
