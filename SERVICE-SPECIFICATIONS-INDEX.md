# MODULE-SPECIFICATIONS-INDEX.md  
_Master checklist for breaking the monolithic PRD into individual, actionable **module** specifications_

---

## Legend  

| Status | Meaning |
|--------|---------|
| ⬜ Not-Started | No spec document yet |
| ✍️ Drafting | Spec in progress |
| ✅ Ready | Spec approved & frozen |

Template reference: `docs/templates/SERVICE-SPEC-TEMPLATE.md` (still valid, but treat “Service” ⇒ “Module”).

---

## 1. Sprint-Based Breakdown & Domain Grouping  

### 🏁 Sprint 0 – Foundation (Weeks 0-2)  
| #  | Domain   | Module                        | Depends on (in-process) | DB Schema(s) | Comm Style | AI Thread | Status |
|---:|----------|------------------------------|-------------------------|--------------|------------|-----------|--------|
| 0.1| Identity | **keycloak-bootstrap**       | —                       | keycloak_rds | OIDC JWKS  | Thread-A  | ⬜ |
| 0.2| Shared   | api-gateway                  | 0.1                     | —            | function call | Thread-A | ⬜ |
| 0.3| Shared   | audit-trail                  | 0.1                     | audit_log    | Kafka `audit.*` | Thread-B | ⬜ |
| 0.4| Shared   | notification                 | 0.1                     | —            | SMTP / WA   | Thread-B | ⬜ |
| 0.5| Geography| geo-location                 | 0.2                     | geography    | func / REST | Thread-C | ⬜ |
| 0.6| Shared   | dev-bootstrap (compose)      | 0.1-0.5                 | —            | —          | Thread-Ops| ⬜ |

### 🟢 Sprint 1 – Core Marketplace (Weeks 2-4)  
| #  | Domain      | Module     | Depends on | DB Schema(s) | Comm Style | AI Thread | Status |
|---:|-------------|-----------|-----------|--------------|------------|-----------|--------|
| 1.1| Marketplace | catalog    | 0.2,0.3  | sku          | func/Kafka `catalog.*` | Thread-A | ⬜ |
| 1.2| Marketplace | rfq        | 1.1       | rfq          | func/Kafka `marketplace.rfq.*` | Thread-A | ⬜ |
| 1.3| Marketplace | quote      | 1.2       | quote        | func       | Thread-B | ⬜ |
| 1.4| Marketplace | contract   | 1.3       | contract     | func       | Thread-B | ⬜ |

### 🛠️ Sprint 2 – Service Ops Foundation (Weeks 4-6)  
| #  | Domain   | Module               | Depends on | DB Schema(s) | Comm Style | AI Thread | Status |
|---:|----------|----------------------|-----------|--------------|------------|-----------|--------|
| 2.1| Service  | asset-registry       | 0.2,0.5  | equipment    | func/Kafka `asset.*` | Thread-C | ⬜ |
| 2.2| Service  | device-registration  | 2.1      | device_reg   | func       | Thread-C | ⬜ |
| 2.3| Service  | qr-manager          | 2.2      | qr_log       | func       | Thread-D | ⬜ |
| 2.4| Service  | ticket              | 2.1      | ticket       | func/Kafka `ticket.*` | Thread-D | ⬜ |
| 2.5| Service  | whatsapp-gateway    | 0.2      | —            | webhook    | Thread-E | ⬜ |
| 2.6| Service  | workflow-engine     | 2.4      | workflow     | Kafka `workflow.*` | Thread-E | ⬜ |

### 🤖 Sprint 3 – AI Layer & Advanced Workflows (Weeks 6-10)  
| #  | Domain | Module            | Depends on | DB Schema(s) | Comm Style | AI Thread | Status |
|---:|--------|-------------------|-----------|--------------|------------|-----------|--------|
| 3.1| AI     | chat-ai           | 2.5      | —            | func/gRPC  | Thread-F | ⬜ |
| 3.2| AI     | negotiation-ai    | 1.3      | —            | func/REST  | Thread-F | ⬜ |
| 3.3| AI     | predictive-maint  | 2.1      | telemetry    | func/REST  | Thread-G | ⬜ |
| 3.4| AI     | dispatch-ai       | 2.6,0.5  | dispatch     | func/REST  | Thread-G | ⬜ |
| 3.5| Service| diagnostic-flow   | 2.6,3.1  | diagnostic   | func       | Thread-H | ⬜ |

### 🚀 Sprint 4 – Scale & Observability (Weeks 10-14)  
| #  | Domain | Module            | Depends on | DB Schema(s) | Comm Style | AI Thread | Status |
|---:|--------|-------------------|-----------|--------------|------------|-----------|--------|
| 4.1| Shared | reporting         | all      | reporting    | GraphQL    | Thread-I | ⬜ |
| 4.2| Shared | parts-inventory   | 2.4      | parts_stock  | func       | Thread-H | ⬜ |
| 4.3| AI     | demand-forecast   | 1.1      | forecast     | func/REST  | Thread-G | ⬜ |
| 4.4| Shared | ci-metrics-export | all      | —            | Prom push  | Thread-Ops | ⬜ |

---

## 2. Module Dependency / Build Order Graph  

1. **Identity → Gateway → Audit/Notification**  
2. Geo-Location → Asset Registry → Device-Reg → QR → Ticket → Workflow  
3. Catalog → RFQ → Quote → Contract  
4. WhatsApp-Gateway → Chat-AI & Ticket  
5. Workflow → Dispatch-AI & Diagnostic-Flow  
6. AI modules consume Kafka events, but **in-process calls preferred** when co-deployed.  

Compilation: single `go build ./cmd/platform` builds all modules; no per-module image builds.

---

## 3. AI Thread Assignment (Monorepo Context)  

| Thread | Focus Area | Initial Module Load |
|--------|------------|---------------------|
| Thread-A | Marketplace core | Keycloak bootstrap, Gateway, Catalog, RFQ |
| Thread-B | Marketplace pricing | Quote, Contract |
| Thread-C | Asset onboarding | Geo-Location, Asset, Device-Reg |
| Thread-D | Service operations | QR-Manager, Ticket |
| Thread-E | Communication & Orchestration | WhatsApp-GW, Workflow |
| Thread-F | Conversational AI | Chat-AI, Negotiation-AI |
| Thread-G | Predictive AI | Predict-Maint, Dispatch-AI, Demand-Forecast |
| Thread-H | Field ops & parts | Diagnostic-Flow, Parts-Inventory |
| Thread-I | Analytics & Reporting | Reporting |
| Thread-Ops | Tooling | Dev-Bootstrap, CI-Metrics |

_All threads commit to same repo; PRs must touch only their module folders plus shared libs._

---

## 4. Module Enablement & Runtime Configuration  

`ENABLED_MODULES` env controls runtime:  
* `"*"` → start all modules (dev/CI)  
* `"catalog,rfq,ticket"` → start subset in prod pod  
* Separate K8s Deployments reuse **same image**, differing only by `ENABLED_MODULES`.

Config file `config.yaml` supports:

```yaml
modules:
  enabled: ["catalog","rfq","ticket"]
gateway:
  listenPort: 8081
```

---

## 5. Inter-Module Communication  

| Interaction | Preferred Path | Notes |
|-------------|----------------|-------|
| Same-pod | Direct function call (no JSON) | via exported Service interface |
| Cross-pod (same image) | gRPC | Proto in `/api` folder |
| Async events | Kafka topic | Avro v⟂ semantic versioning |

---

## 6. Shared Database & Schemas  

Single Postgres cluster.  Schemas isolate bounded contexts:  

| Schema | Owned By Module(s) |
|--------|--------------------|
| `sku` | catalog |
| `rfq` | rfq |
| `quote` | quote |
| `contract` | contract |
| `equipment` | asset-registry |
| `device_reg` | device-registration |
| `ticket` | ticket |
| `workflow` | workflow-engine |
| `telemetry` | predictive-maint |
| `forecast` | demand-forecast |
| `reporting` | reporting |

Row-Level Security (`tenant_id`) enabled globally.

---

## 7. Deployment Order (Single Image Strategy)

1. Deploy **platform image** with `ENABLED_MODULES="*" (dev/staging)**  
2. Production pods created per domain slice:  
   * marketplace-pod → catalog, rfq, quote, contract  
   * service-ops-pod → asset, device-reg, qr, ticket, workflow  
   * ai-pod → chat-ai, negotiation-ai, predictive-maint, dispatch-ai, demand-forecast  
3. Reporting, parts-inventory can be added as dedicated pods if load requires.  
4. ArgoCD waves manage config-only changes (image digest identical).

---

## 8. Status Tracking  

Each **module spec PR** must link back to this table row.  Update icon:  
`⬜` → `✍️` when PR opened, `✍️` → `✅` after review & merge.

---

_End of index_
