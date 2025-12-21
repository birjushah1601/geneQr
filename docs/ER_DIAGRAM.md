# Medical Equipment Platform - ER Diagram

## Visual ER Diagram (Mermaid)

\\\mermaid
erDiagram
    %% Core Organization Entities
    organizations ||--o{ organization_facilities : has
    organizations ||--o{ org_relationships : "parent of"
    organizations ||--o{ org_relationships : "child of"
    organizations ||--o{ engineer_org_memberships : employs
    
    %% Engineers
    engineers ||--o{ engineer_org_memberships : "member of"
    engineers ||--o{ engineer_equipment_types : "certified for"
    engineers ||--o{ engineer_coverage : "covers region"
    engineers ||--o{ service_tickets : "assigned to"
    
    %% Equipment Catalog (Master Data)
    equipment_catalog ||--o{ equipment_part_assignments : "has parts"
    equipment_catalog ||--o{ equipment_registry : "installed as"
    spare_parts_catalog ||--o{ equipment_part_assignments : "used in"
    
    %% Equipment Registry (Installed Units)
    equipment_registry ||--o{ service_tickets : "support requests"
    equipment_registry ||--o{ equipment_maintenance_history : "maintenance"
    equipment_registry ||--o{ equipment_downtime : "downtime events"
    equipment_registry ||--o{ equipment_usage_logs : "usage"
    equipment_registry ||--o{ equipment_documents : "documents"
    equipment_registry }o--|| organizations : "owned by"
    
    %% Service Tickets
    service_tickets ||--o{ ticket_comments : "has comments"
    service_tickets ||--o{ ticket_attachments : "has attachments"
    service_tickets ||--o{ ticket_assignment_history : "assignment history"
    service_tickets ||--o{ ticket_status_history : "status changes"
    service_tickets ||--o{ ai_diagnosis_results : "AI diagnosis"
    service_tickets ||--o{ ai_vision_analysis : "vision analysis"
    
    %% Spare Parts & Marketplace
    spare_parts_catalog ||--o{ marketplace_listings : "listed in"
    spare_parts_catalog ||--o{ spare_parts_suppliers : "supplied by"
    spare_parts_catalog ||--o{ spare_parts_alternatives : "has alternatives"
    spare_parts_catalog ||--o{ spare_parts_bundle_items : "in bundles"
    spare_parts_bundles ||--o{ spare_parts_bundle_items : contains
    
    %% Key Tables
    organizations {
        uuid id PK
        string name
        string type
        int tier
        jsonb contact_info
    }
    
    equipment_catalog {
        uuid id PK
        string product_name
        string manufacturer_name
        string model_number
        string category
        jsonb specifications
    }
    
    equipment_registry {
        varchar id PK
        uuid equipment_catalog_id FK
        string qr_code UK
        string serial_number UK
        string customer_name
        string installation_location
        date installation_date
        date warranty_expiry
        string status
    }
    
    equipment_part_assignments {
        uuid id PK
        uuid equipment_catalog_id FK
        uuid spare_part_id FK
        int quantity_required
        boolean is_critical
    }
    
    spare_parts_catalog {
        uuid id PK
        string part_number
        string part_name
        decimal unit_price
        string currency
        string stock_status
        int lead_time_days
        boolean is_available
    }
    
    service_tickets {
        varchar id PK
        varchar equipment_id FK
        uuid assigned_engineer_id FK
        string ticket_number UK
        string status
        string priority
        text issue_description
        timestamp created_at
    }
    
    engineers {
        uuid id PK
        string name
        string email
        string phone
        int engineer_level
    }
    
    marketplace_listings {
        uuid id PK
        uuid part_id FK
        uuid seller_org_id FK
        decimal price
        string availability
    }
\\\

## Architecture Flow

\\\
MANUFACTURERS → equipment_catalog → equipment_part_assignments → spare_parts_catalog
                       ↓                                                ↓
                equipment_registry ←─────────────────────→ marketplace_listings
                       ↓
                service_tickets ←──────→ engineers
                       ↓
            AI diagnosis + parts assignment
\\\

## Key Relationships

1. **equipment_catalog** (Master) → **equipment_registry** (Instances)
   - One product model → Many installed units

2. **equipment_catalog** → **equipment_part_assignments** → **spare_parts_catalog**
   - Products have compatible parts defined

3. **equipment_registry** → **service_tickets** → **engineers**
   - Installed equipment generates support tickets assigned to engineers

4. **spare_parts_catalog** → **marketplace_listings**
   - Parts available for purchase in marketplace

5. **organizations** → **engineers** (via engineer_org_memberships)
   - Multi-org support for engineers

## Database Functions

- \get_parts_for_ticket(ticket_id)\ - Returns parts for ticket's equipment
- \get_parts_for_registry(registry_id)\ - Returns parts for equipment unit

