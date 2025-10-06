# FINAL-technical-implementation-guide.md  
_Engineering-ready implementation handbook – v1.0 (Sept 2025)_

---

## 1. SYSTEM ARCHITECTURE

### 1.1 High-Level View
```
Clients ──► API-Gateway (REST / gRPC / WS)
                 │
                 ├─► Marketplace Domain (catalog-svc, rfq-svc …)
                 ├─► Service Domain (ticket-svc, sla-svc …)
                 ├─► AI Service Layer (negotiation-ai-svc …)
                 ├─► Shared Services (audit-svc …)
                 └─► Keycloak (External Identity Provider)
                          │
                   Kafka Event Bus
                          │
                 Data Platform (S3 lake → Redshift, Redis FS)
```

### 1.2 Micro-services Breakdown
| Service | Purpose | Language | DB | Key Events |
|---------|---------|----------|----|------------|
| catalog-svc | SKU CRUD, search | Go | Postgres | CatalogUpdated |
| rfq-svc | RFQ lifecycle | Node | Postgres | RFQCreated |
| quote-svc | Quote mgmt | Node | Postgres | QuoteSubmitted |
| contract-svc | PO/contract | Java | Postgres | ContractSigned |
| order-svc | Fulfilment | Java | Postgres | OrderShipped |
| asset-registry-svc | Device master | Node | Postgres | AssetCreated |
| ticket-svc | Service tickets | Node | Postgres | TicketOpened |
| sla-svc | Timers, penalties | Python | Redis | SLABreached |
| dispatch-svc | Tech/parts routing | Go | Postgres | DispatchAssigned |
| parts-inventory-svc | Parts stock | Java | Postgres | PartReserved |
| negotiation-ai-svc | NLP coach | Python | — | InsightCreated |
| demand-forecast-svc | SKU forecast | Python | — | ForecastReady |
| predictive-maint-svc | Failure pred. | Python | — | PreemptiveTicket |
| audit-trail-svc | Immutable logs | Go | Postgres-WORM | — |
| notification-svc | Email/SMS/WS push | Node | — | — |
| **whatsapp-gateway-svc** | WhatsApp Business API integration | Node | — | WhatsAppMessageReceived |
| **qr-service-svc** | QR code gen/validation | Go | Postgres | QRScanned |
| **workflow-orchestrator-svc** | Device-type workflows | Node | Postgres | WorkflowStageCompleted |
| **chat-ai-svc** | Conversational AI triage | Python | — | TriageCompleted |
| **device-context-svc** | Rich device data aggregation | Node | — | DeviceContextLoaded |
| **diagnostic-workflow-svc** | Diagnostic processes | Node | Postgres | DiagnosticCompleted |
| **parts-workflow-svc** | Parts procurement | Java | Postgres | PartsOrdered |
| **device-registration-svc** | Device onboarding workflows | Node | Postgres | DeviceOnboarded |
| **geo-location-svc** | Indian pincode and hospital/lab mapping | Node | Postgres | LocationResolved |

_Note: Identity management is handled by Keycloak (external identity provider) instead of a custom identity-svc_

### 1.3 Data Platform
- **Raw zone (S3 Parquet)**: orders, telemetry, chat, IDocs  
- **Processed (dbt)**: star schemas in Redshift  
- **Online Feature Store**: Redis cluster (tenant-sharded)  
- **Offline Feature Store**: Hive on S3

### 1.4 Integration Specifications
| System | Protocol | Objects | Frequency |
|--------|----------|---------|-----------|
| SAP MM/PM | IDoc ORDERS05, DESADV, INVOIC02 | PO, ASN, Invoice | near-real-time |
| OEM ERP | REST + OAuth2 | Catalog, stock, error codes | ≤5 min |
| IoT Gateways | MQTT TLS | Telemetry JSON | 5–60 s |
| 3PL | REST webhook | Tracking, POD | event-driven |
| WhatsApp Business API | HTTPS Webhook | Chat messages, templates | sub-second |
| LIS/LIMS | HL7/FHIR | Test results, equipment status | near-real-time |
| Keycloak | OIDC/OAuth 2.0 | Authentication, authorization | real-time |

---

## 2. API SPECIFICATIONS

### 2.1 REST (sample set)

**POST /rfqs**  
Request  
`{ "buyerId":"uuid","items":[{"sku":"MRI123","qty":1}],"expiresAt":"2025-11-10T18:00:00Z" }`  
Response 201  
`{ "rfqId":"uuid" }`

**POST /whatsapp/webhook**  
Consumes WhatsApp webhook payloads.  
Responds 200 with `{ "status":"ack" }`

**POST /qr/scan**  
`{ "qrPayload":"<encrypted>" }` → Returns device context  
`{ "deviceId":"uuid","facility":"uuid","skuCode":"MRI123" }`

**POST /devices/register**  
`{ "serialNumber":"SN12345","skuCode":"MRI123","facilityId":"uuid","installationDate":"2025-09-21" }` → Registers new device  
Response 201  
`{ "deviceId":"uuid","qrCode":"base64_image_data" }`

**POST /devices/bulk-import**  
`{ "manufacturerId":"uuid","batchId":"uuid","devices":[{"serialNumber":"SN12345","skuCode":"MRI123"}] }` → Bulk import from manufacturer  
Response 201  
`{ "batchId":"uuid","totalDevices":100,"processedDevices":100,"qrBatchUrl":"https://..." }`

**POST /qr/generate**  
`{ "deviceId":"uuid" }` → Generate QR code for existing device  
Response 200  
`{ "qrCode":"base64_image_data","qrId":"uuid" }`

**POST /qr/validate**  
`{ "qrId":"uuid","scanLocation":{"lat":28.6139,"lng":77.2090} }` → Validate QR deployment  
Response 200  
`{ "isValid":true,"deviceId":"uuid","validatedAt":"2025-09-21T14:30:00Z" }`

**GET /geo/pincodes/{pincode}**  
Returns location details for a specific pincode  
Response 200  
`{ "pincode":"110001","city":"New Delhi","state":"Delhi","district":"Central Delhi","latitude":28.6139,"longitude":77.2090 }`

**GET /geo/healthcare-facilities/nearby**  
Request  
`?latitude=28.6139&longitude=77.2090&radius=10&limit=5&facilityType=all`  
Response 200  
`{ "facilities":[{"id":"uuid","name":"AIIMS Delhi","address":"Ansari Nagar East","pincode":"110029","facilityType":"hospital","bedCount":2500,"specialties":["cardiology","neurology"]}] }`

**GET /geo/laboratories/by-tests**  
Request  
`?testCategories=["pathology","radiology"]&pincode=110001&radius=5`  
Response 200  
`{ "laboratories":[{"id":"uuid","name":"Dr Lal PathLabs","address":"CP Branch","pincode":"110001","labType":"full_service","testCategories":["pathology","biochemistry","microbiology"],"accreditations":["NABL","CAP"]}] }`

**POST /geo/service-areas/coverage**  
Request  
`{ "pincode":"110001","deviceType":"MRI" }`  
Response 200  
`{ "covered":true,"serviceHubId":"uuid","estimatedResponseTime":2.5,"techniciansAvailable":3 }`

**PATCH /tickets/{id}/status**  
`{ "status":"in_progress" }` → 204

_Error codes_: 400 (validation), 403 (auth), 404, 409 (state), 500.

### 2.2 GraphQL
```
type Query {
  catalog(search:String!, first:Int=20, after:ID): CatalogConnection!
  ticket(id:ID!): Ticket!
}

type Mutation {
  createRFQ(input: CreateRFQInput!): RFQ!
  suggestCounterOffer(rfqId:ID!): NegotiationTip   # AI hook
}
```

### 2.3 WebSocket Channels
- `/ws/negotiation/{rfqId}` → messages + `coachTip` events
- `/ws/tickets/{tenantId}`   → real-time ticket status
Frames: JSON, max 4 KB.

### 2.4 Authentication & Authorization
- **OAuth 2.0/OIDC** via Keycloak for all authentication flows
- **Authorization Code Flow** for web/mobile user authentication
- **Client Credentials Flow** for service-to-service communication
- **JWT RS256 tokens** with tenant context and role claims
- **Token validation** at API Gateway with caching for performance
- **RBAC** managed through Keycloak roles (`buyer_admin`, `oem_user`, `tech`)
- **FIDO2/WebAuthn** integration through Keycloak for strong MFA
- **Healthcare compliance**: HIPAA-aligned audit logging of all auth events

### 2.5 Device Registration API Flows

**Flow 1: Existing Device Registration**
```
Client ──► POST /devices/register
           └─► device-registration-svc validates & creates device profile
                └─► POST /qr/generate (internal)
                     └─► qr-service-svc generates secure QR code
                          └─► Event: DeviceOnboarded
                               └─► asset-registry-svc updates inventory
```

**Flow 2: Manufacturer Bulk Import**
```
OEM ERP ──► POST /devices/bulk-import
             └─► device-registration-svc creates batch
                └─► Async processing of each device
                     └─► POST /qr/generate (batch mode)
                          └─► Event: BatchProcessed
                               └─► notification-svc alerts completion
                                    └─► QR deployment kit generated
```

**Flow 3: QR Deployment & Validation**
```
Mobile App ──► POST /qr/validate
                └─► qr-service-svc verifies signature & location
                     └─► Event: QRDeployed
                          └─► device-registration-svc updates status
                               └─► asset-registry-svc activates device
```

### 2.6 Keycloak Integration Architecture

**Multi-Tenant Realm Configuration:**
```
Keycloak
├── Master Realm (platform administration)
├── Tenant Realms (one per healthcare organization)
│   ├── Hospital Chain Realm
│   │   ├── Users (staff, admins)
│   │   ├── Roles (buyer_admin, technician, viewer)
│   │   ├── Groups (by department, location)
│   │   └── Client applications (web, mobile, API)
│   ├── Laboratory Chain Realm
│   │   ├── Users (lab staff, managers)
│   │   ├── Roles (lab_admin, technician)
│   │   └── Client applications
│   └── Manufacturer/OEM Realm
│       ├── Users (sales, service)
│       ├── Roles (seller_admin, service_manager)
│       └── Client applications
└── Shared Services Realm (platform microservices)
    └── Service accounts (for service-to-service auth)
```

**OAuth 2.0/OIDC Flows:**

**1. User Authentication Flow (Authorization Code):**
```
User ──► Web/Mobile App ──► Keycloak Authorization Endpoint
                               │
                               ▼
                            Login Form
                               │
                               ▼
                         Authentication
                               │
                               ▼
                        Authorization Code
                               │
                               ▼
Web/Mobile App ──► Keycloak Token Endpoint ──► Access/Refresh Tokens
        │
        ▼
API Gateway ──► Token Introspection/Validation
        │
        ▼
Microservices (with tenant context & roles)
```

**2. Service-to-Service Authentication (Client Credentials):**
```
Microservice A ──► Keycloak Token Endpoint ──► Service Account Token
                                                      │
                                                      ▼
                                           Microservice B (validates token)
```

**Token Validation Strategy:**
- API Gateway validates all incoming tokens
- JWT signature verification using Keycloak's public keys (cached, auto-rotated)
- Claims extraction for tenant_id, roles, and permissions
- Token introspection for high-security operations
- Token propagation for service-to-service calls

**Tenant Isolation:**
- Each tenant (healthcare organization) has a dedicated Keycloak realm
- JWT contains tenant_id claim used for row-level security in databases
- API Gateway enriches requests with tenant context headers
- Cross-tenant access controlled via explicit permissions

**Healthcare Compliance Features:**
- Comprehensive authentication audit logs (HIPAA 164.312(b))
- Automatic session timeout (HIPAA 164.312(a)(2)(iii))
- Role-based access control (HIPAA 164.308(a)(4)(i))
- Emergency access procedure support ("break glass")
- Strong MFA enforcement for sensitive operations
- Password policies aligned with HIPAA requirements

---

## 3. DATABASE DESIGN

### 3.1 Core Schemas (PostgreSQL)
```
tenants(id PK, name, tier, created_at)
users(id PK, tenant_id FK, email, keycloak_id, role, created_at, ...)
sku(id PK, tenant_id, code, name, imdr_class, list_price, attrs JSONB)
rfq(id PK, tenant_id, status, expires_at, created_at)
rfq_item(id PK, rfq_id FK, sku_id FK, qty)
quote(id PK, rfq_id FK, seller_tenant_id, amount, validity)
contract(id PK, buyer_tenant_id, seller_tenant_id, status, signed_at)
equipment(id PK, tenant_id, udi, sku_id, install_date, amc_status, current_sla)
ticket(id PK, tenant_id, equipment_id, priority, status, error_code, opened_at)
device_registration_batch(id PK, manufacturer_id FK, status, total_count, processed_count, created_at)
qr_deployment_log(id PK, qr_id FK, device_id FK, location_lat, location_lng, deployed_by_user_id FK, validated_at)
manufacturer_device_feed(id PK, manufacturer_id FK, serial_number, sku_code, import_batch_id FK, status)
indian_pincodes(pincode PK, city, state, district, latitude, longitude)
healthcare_facilities(id PK, name, address, pincode FK, facility_type, bed_count, specialties JSONB, contact_info JSONB)
laboratories(id PK, name, address, pincode FK, lab_type, test_categories JSONB, accreditations JSONB, contact_info JSONB)
service_areas(id PK, hub_location, covered_pincodes JSONB, max_response_time_hours)
healthcare_facility_equipment(id PK, facility_id FK, device_count_by_category JSONB, lab_equipment_types JSONB, last_updated)
```
`facility_type` enum _(exhaustive for Phase-1 rollout)_:  
`hospital`, `diagnostic_lab`, `pathology_lab`, `radiology_center`, `imaging_center`, `cardiac_center`, `oncology_center`, `primary_health_center`, `community_health_center`, `dispensary`, `eye_care_center`, `dental_clinic`, `dialysis_center`, `physiotherapy_center`, `sleep_study_center`, `emergency_center`, `trauma_center`, `corporate_health_center`, `telemedicine_center`

### 3.2 AI Feature Store Keys
`feature_store(tenant_id, entity_id, ts, feature_group, feature_name, value)`  
Redis key pattern: `fs:{tenant}:{entity}:{feature}`

### 3.3 Multi-Tenancy
- Row-level security (RLS) on `tenant_id`.  
- Citus shard per tenant group (>50 GB).  
- Sequence ID range partitioned to avoid hot spots.

### 3.4 Performance & Indexing
- `btree` on `(tenant_id, id)` for OLTP tables.  
- `GIN` on `attrs` JSONB.  
- TimescaleDB hypertables for telemetry (partition by hour).  
- Regular `VACUUM ANALYZE`; autovacuum tuned for 1 GB tables.

### 3.5 Sharding Strategy
- Hash-partitioning on `tenant_id` for multi-tenant tables.
- Range-partitioning on `created_at` for time-series data.
- Colocation groups for tenant-specific joins.

### 3.6 Indian Geography Data Model
```
                                ┌───────────────────┐
                                │  indian_pincodes  │
                                │ ──────────────── │
                                │ pincode (PK)     │
                                │ city             │
                                │ state            │
                                │ district         │
                                │ latitude         │
                                │ longitude        │
                                └─────────┬─────────┘
                                          │
                                          │
         ┌──────────────┬─────────┴───────────┬──────────────┐
         │              │                     │              │
┌────────▼────────┐     │                     │     ┌────────▼────────┐
│  healthcare_    │     │                     │     │  service_areas  │
│   facilities    │     │                     │     │ ────────────── │
│ ────────────── │     │                     │     │ id (PK)        │
│ id (PK)        │     │                     │     │ hub_location    │
│ name           │     │                     │     │ covered_pincodes│
│ address        │◄────┘                     └────►│ max_response_   │
│ pincode (FK)   │                                 │   time_hours    │
│ facility_type  │                                 └────────┬────────┘
│ bed_count      │                                          │
│ specialties    │                                          │
│ contact_info   │                                          │
└────────┬───────┘                                          │
         │                                                  │
┌────────▼────────┐    ┌────────────────┐                   │
│ healthcare_     │    │  laboratories  │                   │
│ facility_       │    │ ────────────── │                   │
│ equipment       │    │ id (PK)        │                   │
│ ────────────── │    │ name           │                   │
│ id (PK)        │    │ address        │◄──────────────────┘
│ facility_id(FK) │    │ pincode (FK)   │                   │
│ device_count_   │    │ lab_type       │                   │
│  by_category    │    │ test_categories│                   │
│ lab_equipment_  │    │ accreditations │                   │
│  types          │    │ contact_info   │                   │
│ last_updated    │    └────────────────┘                   │
└────────────────┘                                          │
                                                            │
                                               ┌────────────▼───────┐
                                               │    technician      │
                                               │     coverage       │
                                               │ ────────────────── │
                                               │ tech_id (FK)       │
                                               │ service_area_id    │
                                               │ specialties        │
                                               │ availability       │
                                               └────────────────────┘
```

**Key Relationships:**
- Each healthcare facility and laboratory is associated with a specific pincode
- Service areas define coverage regions based on pincodes
- Healthcare facility equipment tracks medical devices across hospitals and labs
- Technician coverage maps specialists to service areas

**Major Indian Cities Covered:**
- Delhi NCR (110001-110096, 122001-122018, 201301-201310)
- Mumbai (400001-400104)
- Bangalore (560001-560114)
- Chennai (600001-600119)
- Kolkata (700001-700156)
- Hyderabad (500001-500096)
- Ahmedabad (380001-380061)
- Pune (411001-411047)
 Plus **100+** additional tier-2/3 cities – overall dataset **≥ 35 000 facilities** (hospitals, labs, diagnostic/imaging centres, PHCs, clinics, emergency hubs, corporate health sites)

**Major Healthcare Facility Networks:**
- AIIMS Network (Delhi, Bhopal, Bhubaneswar, etc.)
- Apollo Hospitals (Pan-India)
- Fortis Healthcare
- Max Healthcare
- Manipal Hospitals
- Narayana Health
- Medanta
- Government Medical Colleges & Hospitals
- **Diagnostic & Imaging Chains** – Scans and Care, Mahajan Imaging, Aarthi Scans, Vijaya Diagnostics, Lucid Diagnostics
- **Primary Healthcare Facilities** – PHCs, CHCs, Urban PHCs, Mohalla Clinics (Delhi), BMC Dispensaries (Mumbai)
- **Specialty Clinics & Centres** – Sankara Eye, Dr. Agarwal's Eye, Clove Dental, NephroPlus Dialysis, Portea Physio
- **Emergency & Trauma Networks** – StanPlus, Ziqitza, Apollo 24×7 ERs, state ambulance fleets (108/102)
- **Corporate & Tele-health Networks** – HealthCube, 1mg Health Centres, Practo Clinics, Jio Health Hub

**Major Laboratory Chains:**
- Dr. Lal PathLabs (1,200+ labs nationwide)
- SRL Diagnostics (400+ labs)
- Metropolis Healthcare (2,500+ collection centers)
- Thyrocare (3,000+ collection centers)
- Quest Diagnostics India
- Suburban Diagnostics
- Neuberg Diagnostics
- Vijaya Diagnostic Centre
- Aarthi Scans and Labs
- Core Diagnostics

**Laboratory Equipment Categories:**
- Analyzers (biochemistry, hematology, immunoassay)
- PCR & molecular diagnostic systems
- Histopathology equipment
- Imaging systems (X-ray, CT, MRI, ultrasound)
- Microbiology automation
- Blood banking equipment
- Point-of-care testing devices
- Sample processing systems
- LIMS (Laboratory Information Management Systems)
- Cold storage & sample management
- **Specialty equipment** – Dialysis machines, phaco lasers, dental CBCT, physiotherapy lasers, polysomnography systems
- **Emergency / POCT** – Portable ventilators, transport monitors, defibrillators, handheld ultrasound, cardiac marker readers
- **Primary-care equipment** – Vaccine refrigerators, multiparameter analyzers, digital otoscopes, digital BP/ECG kiosks
- **Corporate-health / Tele-medicine** – Digital stethoscopes, dermatoscopes, tele-ICU carts, remote vitals gateways

---

## 4. AI TECHNICAL SPECIFICATIONS

### 4.1 Model Overview
| Model | Type | Input | Output | Target |
|-------|------|-------|--------|--------|
| price_band | XGBoost reg. | SKU, qty, historic deals | ₹ band | MAE ≤6 % |
| negotiation_tip | LLM prompt | Chat context | tip string | latency ≤500 ms |
| demand_forecast | Prophet + XGB | time-series | weekly forecast | MAPE ≤12 % |
| failure_pred | LSTM | telemetry seq | P(fail)<30d | recall ≥0.8 |
| dispatch_rl | DQN | state graph | tech/parts assignment | MTTR ↓35 % |
| triage_intent | BERT-CLS | chat text, media | intent label | F1 ≥0.88 |

**dispatch_rl (Enhanced for Indian Geography)**
- **Input Features**: 
  - Technician location (lat/long)
  - Device location (hospital/lab pincode)
  - Traffic conditions (time-of-day dependent)
  - Technician expertise match score
  - Parts availability at nearest hub
  - Facility tier and SLA level
  - Historical travel times between pincodes
  - Monsoon season impact factor
- **Constraints**:
  - City-specific traffic patterns (e.g., Delhi peak hours 9-11 AM, 5-8 PM)
  - Travel restrictions in certain zones (e.g., Old Delhi narrow lanes)
  - Hospital/lab access hours and protocols
  - Pincode-based routing optimization with real-time traffic API integration
  - Service hub proximity weightings
- **Optimization Targets**:
  - Minimize travel time considering Indian traffic patterns
  - Maximize first-time fix rate
  - Optimize technician utilization across service areas
  - Balance urgent vs. routine service calls
  - Respect SLA tiers while minimizing total system cost

### 4.2 ML Pipeline
1. **Ingest** via Kafka Connect → S3 Raw.  
2. **dbt build** processed views.  
3. **Feature engineering** (PySpark) → offline FS.  
4. **Training** in SageMaker pipelines (nightly or drift).  
5. **Registry** MLflow, model card, license.  
6. **Deployment** via KFServing → K8s GPU/CPU inference pods.

### 4.3 Real-Time Inference
- REST `/predict` p95 ≤120 ms.  
- Horizontal pod autoscaling on CPU/GPU util.  
- Canary 5 % traffic for new model versions.

### 4.4 Monitoring & Governance
- Prometheus + custom exporter: latency, throughput, accuracy.  
- Evidently-AI drift monitor PSI/KS.  
- Alert thresholds: latency >200 ms 95p, PSI >0.2, accuracy drop >5 pp.  
- Rollback via MLflow stage change.

---

## 5. DEVELOPMENT EPICS

| Epic | Points | Dep | Definition of Done |
|------|--------|-----|--------------------|
| EP-CAT-01 Catalog Core | 34 | — | CRUD + search API, 80 % tests |
| EP-RFQ-02 RFQ Flow | 40 | CAT-01 | RFQ->Quote->Contract states, WS events |
| EP-SAP-03 SAP Connector | 26 | RFQ-02 | ORDERS05 send, DESADV recv, unit tests |
| EP-AI-04 Negotiation Coach MVP | 21 | RFQ-02 | /analyzeQuote, latency ≤500 ms |
| EP-TCK-05 Ticket Engine | 32 | CAT-01 | CRUD, RLS, SLA hook |
| EP-DSP-06 Dispatch AI | 24 | TCK-05, PART-07 | REST /dispatch, RL model v0 |
| PART-07 Parts Inventory | 28 | CAT-01 | BOM tables, reserve API |
| EP-IOT-08 Telemetry Ingest | 18 | TCK-05 | MQTT broker, device auth |
| EP-PM-09 Predictive Maint | 30 | IOT-08 | LSTM model, preemptive tickets |
| EP-KC-10 Keycloak Setup | 22 | — | Multi-tenant realms, OIDC flows, FIDO2 MFA |
| EP-DEVID-11 Multi-Modal Device ID | 24 | TCK-05 | OCR, barcode recognition, secure payload, unified API |
| EP-WA-12 WhatsApp Integration | 28 | QR-11 | whatsapp-gateway-svc, webhook, template flow |
| EP-WF-13 Workflow Engine | 32 | WA-12 | workflow-orchestrator-svc, stage DSL |
| EP-DIAG-14 Diagnostic Workflow | 30 | WF-13 | diagnostic-workflow-svc, triage AI, mobile checklist |
| EP-REG-15 Device Registration System | 26 | DEVID-11 | Bulk import, QR generation, deployment tracking |
| EP-GEO-16 Indian Geography System | 22 | TCK-05 | Pincode database, hospital mapping, service area coverage |
| EP-LAB-17 Laboratory Integration | 28 | GEO-16 | Lab database, test catalog, equipment mapping |

_Points assume 1 PD = 1 sp._

**Definition of Ready**  
• Story has AC, API contract, mock UI, size ≤8 sp.  

**Critical Path**: EP-CAT-01 → RFQ-02 → SAP-03 / AI-04.

---

## 6. DEPLOYMENT & DEVOPS

### 6.1 Infrastructure
| Layer | Tech | Notes |
|-------|------|-------|
| Runtime | AWS EKS (K8s 1.29) | gp3 + GPUSpot for ML |
| DB | AWS RDS Postgres 15 (Citus) | Multi-AZ, 15 k IOPS |
| Cache | Elasticache Redis cluster (6.2) | TLS, auth |
| Message | MSK Kafka 3.x | 3×broker, tiered storage |
| Object | S3 with WORM lock | Audit logs |
| Identity | Keycloak 22.x | Multi-AZ, RDS Postgres backend |

### 6.2 CI/CD
- **GitHub Actions**: build → test → container scan → push to ECR.  
- **Argo CD** GitOps deploy to EKS, blue/green with Argo Rollouts.  
- **Argo Workflows** for ML batch & data jobs.

### 6.3 Observability
- **Prometheus + Grafana** dashboards per service SLO.  
- **Loki** for logs, **Tempo** for traces (OpenTelemetry).  
- SLO alert: error budget burn >2 % triggers PagerDuty.

### 6.4 Security & Compliance
- AWS KMS CMKs, Secrets Manager.  
- IAM least privilege; service-mesh mTLS (Istio).  
- Quarterly VAPT; dependency scanning (Snyk).  
- SOC 2 & ISO 13485 documentation stored in Confluence.

### 6.5 Disaster Recovery
- Multi-AZ default; cross-region read replica for Postgres.  
- RPO ≤ 5 min (WAL streaming), RTO ≤ 15 min.  
- Monthly failover drills.

---

## 7. ACTION ITEMS FOR SPRINT 0

1. Stand up EKS cluster & Argo CD bootstrap (Dev env).  
2. Create Citus Postgres with RLS template.  
3. Deploy Keycloak with master realm and initial tenant realm templates.
4. Configure OIDC clients and test OAuth flows with API Gateway.
5. Scaffold catalog-svc & RFQ-svc repos with CI template.  
6. Set up Kafka topics & schema registry.
7. Import Indian pincode database and setup geo-location-svc.
8. Create initial hospital and laboratory mapping for Delhi NCR, Mumbai, and Bangalore.
9. Define equipment categories for both hospital and laboratory devices.

**Target:** Dev environment operational by end of Sprint 0 (2 weeks).

---

_End of guide_
