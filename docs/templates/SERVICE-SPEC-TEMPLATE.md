# SERVICE SPECIFICATION TEMPLATE  
_Replace bracketed placeholders → ❰Example❱ with real values._

---

## 1. Service Overview  
| Field | Description |
|-------|-------------|
| Service Name | ❰catalog-svc❱ |
| Domain / Bounded Context | ❰Marketplace❱ |
| Purpose | _One sentence: “Provides CRUD & search for SKUs.”_ |
| Primary Responsibilities | • …<br>• … |
| Out-of-Scope | _Anything explicitly not handled._ |
| Success Metrics | e.g. p95 latency ≤300 ms, error rate < 0.1 % |

---

## 2. API Contract  

### 2.1 REST (OpenAPI)  
_Link to `openapi.yaml` or embed snippet._  
```
GET /v1/skus?search=... 
200 → [ { id, code, name, … } ]
```

### 2.2 gRPC / Protobuf  
```
service Catalog {
  rpc ListSKUs(ListSKUsRequest) returns (ListSKUsResponse);
}
```

### 2.3 WebSocket / Streaming  
Describe channel names and payloads.

---

## 3. Domain Model  

| Entity | Description | Key Business Invariants |
|--------|-------------|-------------------------|
| ❰SKU❱ | Stock keeping unit | `code` unique per tenant, `price>0` |
| ❰RFQ❱ | … | RFQ must expire in future |

Include class diagrams or Mermaid if helpful.

_Value Objects_: Money, Email, TenantID, etc.

---

## 4. Event Schema (Pub/Sub)  

| Event | Producer | Payload (Avro/JSON) | Versioning |
|-------|----------|---------------------|------------|
| `catalog.sku.created.v1` | catalog-svc | `{ skuId, tenantId, ... }` | SemVer |

Provide full schema file paths.

---

## 5. Persistence Schema & Migrations  

*Schema Definition*  
```
CREATE TABLE sku (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  price NUMERIC CHECK (price>0),
  created_at TIMESTAMPTZ DEFAULT now()
);
```

*Migration Tool*: golang-migrate / knex / Prisma.  
*RLS Policy*: `tenant_id = current_setting('app.tenant_id')`.

---

## 6. Error Handling & Response Codes  

| Code | HTTP | Description | Retryable |
|------|------|-------------|-----------|
| `SKU_NOT_FOUND` | 404 | SKU id invalid | ❌ |
| `TENANT_MISMATCH` | 403 | JWT tenant ≠ row tenant | ❌ |

Error format:
```
{
  "error": "SKU_NOT_FOUND",
  "message": "Sku xyz not found",
  "traceId": "abc-123"
}
```

---

## 7. Security Requirements  

* Auth: Keycloak OIDC (Authorization Code / Client-Creds).  
* Required scopes: `catalog.read`, `catalog.write`.  
* Tenant Isolation: Row-Level Security enforced.  
* RBAC Matrix:

| Role | Read | Write | Admin |
|------|------|-------|-------|
| buyer_admin | ✅ | ❌ | ❌ |
| seller_admin| ✅ | ✅ | ❌ |
| catalog_admin | ✅ | ✅ | ✅ |

* mTLS between services.

---

## 8. Testing Requirements  

| Layer | Coverage Target | Tooling |
|-------|-----------------|---------|
| Unit  | ≥85 % lines | `go test`, Vitest |
| Integration | Key use-cases | Testcontainers |
| Contract | 100 % provider pass | Pact |
| Performance | p95 ≤300 ms @5 k RPS | k6 |
| Security | 0 critical vuln | Semgrep, Trivy |

CI badge: tests must be green before merge.

---

## 9. Deployment Specification  

* Container: Distroless, health-check `/healthz`.  
* Helm Chart Path: `helm/<service>`  
* K8s Resources:

```yaml
resources:
  requests: { cpu: "200m", memory: "256Mi" }
  limits:   { cpu: "500m", memory: "512Mi" }
```

* Rollout: Argo Rollouts canary 10-50-100.  
* Env Vars: listed in §11.

---

## 10. Observability  

| Aspect | Implementation |
|--------|----------------|
| Logging | zerolog → Loki, include `tenant_id`, `trace_id` |
| Metrics | Prometheus: `http_requests_total`, `sku_created_total` |
| Tracing | OpenTelemetry auto-instrumentation → Tempo |
| Dashboards | `/grafana/dashboards/catalog.json` |

Alerts: p95 > 300 ms for 2 min → PagerDuty sev-2.

---

## 11. Configuration & Environment Variables  

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | HTTP listen |
| `KEYCLOAK_URL` | http://keycloak:8080 | OIDC issuer |
| `DB_DSN` | `postgres://...` | Postgres DSN |
| `KAFKA_BROKERS` | kafka:9092 | Comma separated |

Config file sample `config.yaml` embedded if needed.

---

## 12. Inter-Service Dependencies  

| Dependent On | Reason | Interface |
|--------------|--------|-----------|
| Keycloak | Auth | OIDC introspection |
| geo-location-svc | Resolve facility | REST `/geo/pincodes` |

Downstream events consumed: list topics.

---

## 13. Performance & SLA  

* p95 latency ≤300 ms  
* Throughput ≥5 k RPS sustained  
* Error rate <0.1 %  
* CPU ≤70 % at peak, memory ≤75 %.

Load-test plan: `perf/catalog-load.k6.js`.

---

## 14. Healthcare Compliance  

| Regulation | Artifact |
|------------|----------|
| IMDR 2017 | Audit field `mdi_code` stored |  
| DPDP 2023 | Data minimisation checklist |  
| HIPAA Safe‐Harbor | No PHI stored, only device data |

---

## 15. Migration & Data Seeding  

* Initial seed: default SKUs for demo tenant (see `migrations/202501011200_seed.sql`).  
* Safe roll-forward migrations only; rollback through backups.  
* Seed script runnable via `make seed-dev`.

---

### **Sign-Off**  
* Product Owner: __________  
* Tech Lead: __________  
* QA Lead: __________  
* Date Approved: __________
