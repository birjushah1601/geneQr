# QA-TESTING-SPECIFICATIONS.md  
_Intelligent Medical Equipment Platform – Quality Assurance Manifest_

---

## 1. Global Testing Strategy  

| Layer | Goal | Tooling | Frequency |
|-------|------|---------|-----------|
| Unit (≈80 %) | Validate domain & app logic in isolation | Go `testing`, `vitest`, `jest` | On every commit |
| Integration (≈15 %) | Verify service-to-DB / service-to-MQ contracts | Testcontainers, `docker-compose`, Postgres, Kafka | PR & nightly |
| Contract (<3 %) | Guarantee API compatibility across teams | **Pact** (consumer/provider) | PR & nightly |
| End-to-End (≈2 %) | Validate critical user journeys | Cypress, Playwright | Nightly & pre-release |
| Performance | Ensure SLA adherence (p95 latency, throughput) | k6, Locust | Weekly & pre-release |
| Security | Detect vulns, authZ bypass, OWASP | ZAP, Semgrep, Trivy, kube-bench | PR & quarterly |
| Chaos / Resilience | Validate graceful degradation | Litmus, Gremlin | Monthly |
| Compliance | HIPAA/DPDP rule enforcement | Custom suites + Evidence capture | Quarterly audit |

---

## 2. Unit Testing Standards  

1. Follow **Arrange-Act-Assert** pattern.  
2. Cover all _business branches_, **≥80 % statement coverage**; CI gate fails below.  
3. No DB/network in unit tests – use fakes/mocks.  
4. Deterministic; must run in <100 ms each.  
5. Naming: `Test<Subject>_<Scenario>_<Expectation>`.

_Example – Go (RFQ domain)_  
```go
func TestRFQ_AddItem_WhenExpired_ReturnsErr(t *testing.T) {
  rfq := domain.NewRFQ("rfq-1", tenant, time.Now().Add(-1*time.Hour))
  err := rfq.AddItem("MRI-001", 1)
  assert.Equal(t, domain.ErrRFQExpired, err)
}
```

---

## 3. Integration Testing Approach  

* Goal: verify real adapters (Postgres, Kafka, Keycloak) behave as expected.  
* **Testcontainers** spins ephemeral Docker services per test suite.  
* Use tenant-scoped JWT from Keycloak test realm.  
* Clear data after each test (SQL `TRUNCATE ... CASCADE`).  

```yaml
# docker-compose.test.yml excerpt
postgres:
  image: postgres:15-alpine
kafka:
  image: confluentinc/cp-kafka:7.5
keycloak:
  image: quay.io/keycloak/keycloak:22
```

---

## 4. Contract Testing  

* **Pactflow** workflow:  
  1. Consumers publish contracts (e.g., UI → catalog-svc).  
  2. Providers verify & tag with git SHA.  
  3. CI blocks merge if contract break detected.  

Versioning rules: MAJOR change requires new endpoint or new event version (`*.v2`).

---

## 5. End-to-End Testing  

* **Playwright** for Web, **Cypress** for Mobile-Web.  
* Runs against docker-compose “mini-env” + seeded data.  
* Focus journeys:  
  - “Buyer creates RFQ, receives quote, converts PO”  
  - “Technician scans QR, closes ticket within SLA”  
  - “Lab manager gets predictive alert & schedules maintenance”  
* Screenshots + trace uploaded to S3.

---

## 6. Performance Testing Requirements  

| Metric | Target | Tool | Threshold Gate |
|--------|--------|------|----------------|
| p95 API Latency | ≤300 ms | k6 | ❌ fail if >10 % over |
| Throughput | 5 k RPS per svc | k6 | ❌ fail if error rate >0.1 % |
| WhatsApp Triage | <3 s end-to-end | Locust + Twilio sandbox | ❌ |

Scripts live in `perf/` and executed via `make perf`.

---

## 7. Security Testing Protocols  

1. **Static Analysis** – Semgrep (PR gate).  
2. **Dependency Scanning** – Trivy / `govulncheck`.  
3. **Dynamic Scans** – OWASP ZAP nightly against staging.  
4. **Container Hardening** – Trivy image scan in CI.  
5. **Kubernetes** – kube-bench + Kyverno policies.  
6. **Pen-Test** – External vendor bi-annually.

---

## 8. Healthcare Compliance Testing  

| Standard | Test Focus | Evidence |
|----------|------------|----------|
| HIPAA 164.312(b) | Auth audit log present in Keycloak & forwarded to SIEM | Log dump + SIEM alert test |
| HIPAA 164.312(a)(2)(iii) | Auto logout after 15 min idle | Cypress script verifying session expiry |
| DPDP 2023 | Consent capture & data minimisation | API contract + privacy impact test |
| ISO 13485 | Service process traceability | Ticket lifecycle E2E test |

---

## 9. Multi-Tenant Isolation Testing  

* **RLS Verification**:  
```sql
SET app.tenant_id = 'tenant-A';
SELECT COUNT(*) FROM rfqs; -- expect 3
SET app.tenant_id = 'tenant-B';
SELECT COUNT(*) FROM rfqs; -- expect 0
```  
CI runs this query set after every schema migration.

* **JWT Swap Attack** – ensure tenant claim mismatch returns 403.  
* **Keycloak Realm Isolation** – user of hospital realm cannot introspect lab realm tokens.

---

## 10. Load & Chaos Engineering  

* **k6 soak test** – 24 h run at 60 % peak traffic.  
* **Gremlin attacks**:  
  - CPU hog on ticket-svc container (expect auto-scale).  
  - Kafka broker network latency (expect retry logic).  
  - Postgres failover simulation (RTO ≤15 min).  

Chaos run monthly in non-prod env.

---

## 11. Test Data Management & Privacy  

1. Synthetic data via **GoFaker** + **Mockaroo** schema = no PHI.  
2. Production data never copied; use **DB-Anon** tool for masked subsets.  
3. Data seeding scripts live under `dev/seed/`.  
4. GDPR “right to be forgotten” test ensures anonymisation job clears traces.  

---

## 12. CI/CD Pipeline Integration  

```
.github/workflows/ci.yaml
  ├─ lint-go
  ├─ test-unit-go
  ├─ test-unit-node
  ├─ test-integration (docker-compose)
  ├─ pact-verify
  ├─ coverage-check (>80 %)
  ├─ trivy-scan
  └─ build-push-image
```
Failed stage blocks merge; results posted to PR.

---

## 13. Quality Gates  

| Gate | Threshold |
|------|-----------|
| Unit Coverage | ≥80 % |
| Integration Pass | 100 % |
| Lint / Static | 0 error, 0 critical warn |
| Performance | ≤110 % latency budget |
| Security | 0 critical/high vuln |
| Contract | 0 breaking changes |
| SonarQube Maintainability | A |

---

## 14. Tools & Framework Matrix  

| Layer | Language | Framework |
|-------|----------|-----------|
| Unit (Go) | Go testing + Testify | Coverage via `go tool cover` |
| Unit (Node) | Vitest | NYC coverage |
| Integration | Testcontainers (Go + Node) | |
| Contract | Pact Go / Pact JS | Pactflow broker |
| E2E | Playwright, Cypress | |
| Performance | k6 (scripted), Locust (Python) | |
| Security | Semgrep, Trivy, ZAP | |
| Chaos | Gremlin, Litmus | |

---

## 15. Test Automation Strategy  

1. **`make test`** umbrella command runs fast unit tests.  
2. **Git Hooks** – pre-push executes `make lint && make test-unit`.  
3. **Nightly Jenkins** triggers heavy suites: integration, contract, e2e, perf.  
4. **Tag-based Selection** – Smoke @smoke, Regression @regress, PCI @pci.  
5. **Allure** reporting auto-published to S3.  
6. **Slack Bot** posts red pipeline summaries with owning squad.  

---

### **Appendix A – Sample Compliance Test (HIPAA Audit Log)**  

```go
func TestAuditLog_HIPAA_Compliance(t *testing.T) {
  // Simulate user login via Keycloak test realm
  token := testkc.LoginAs("nurse@demo-hospital.com", "Pass123!")
  // Call protected endpoint
  res := httpclient.Get("/v1/asset/123", token)
  assert.Equal(t, 200, res.StatusCode)

  // Validate audit record
  log := auditRepo.FindLastByUser("nurse@demo-hospital.com")
  assert.Equal(t, "READ_ASSET", log.Action)
  assert.WithinDuration(t, time.Now(), log.Timestamp, time.Minute)
  assert.NotEmpty(t, log.TraceID)
}
```

---

### **Appendix B – Tenant Isolation Smoke Test (Bash)**  

```bash
#!/usr/bin/env bash
TOKEN_A=$(get_token tenantA)
TOKEN_B=$(get_token tenantB)

curl -s -H "Authorization: Bearer $TOKEN_A" http://localhost:8081/rfqs > /tmp/a.json
curl -s -H "Authorization: Bearer $TOKEN_B" http://localhost:8081/rfqs > /tmp/b.json

jq -e 'length==0' /tmp/b.json   # expect empty for Tenant B
diff /tmp/a.json /tmp/b.json && exit 1  # should differ
echo "Tenant isolation passed"
```

---

_This document is the governing reference for QA across all services. Any deviation requires approval from the QA chapter lead._  
