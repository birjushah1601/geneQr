# BOOTSTRAP-DEVELOPMENT-SETUP.md  
_A 30-minute path from zero to a fully-running local stack (monorepo edition)_

---

## 1. Prerequisites & System Requirements  

| Component      | Minimum           | Recommended        |
| -------------- | ----------------- | ------------------ |
| OS             | macOS 12 / Ubuntu 22 / WSL2 | macOS 13 / Ubuntu 24 |
| CPU            | 4 vCPU            | 8 vCPU             |
| RAM            | 8 GB              | 16 GB              |
| Disk           | 20 GB free (SSD)  | 40 GB              |
| Tools          | `git` ≥2.37, `make`, `curl`, `docker` ≥24, Docker Compose plugin, **Go 1.22** |

> TIP Docker Desktop on macOS/Windows ships with everything except `go` and `make`.

---

## 2. Directory Layout (Monorepo)  

```
medical-platform/
├─ cmd/                  # Binary entrypoints
│   ├─ platform/         # Single binary starting N modules
│   └─ gateway/          # (optional) standalone gateway binary
├─ internal/             # All bounded-context modules
│   ├─ shared/           # Auth, tenant, observability
│   ├─ marketplace/ …    # catalog, rfq, quote, contract
│   ├─ serviceops/ …     # asset, ticket, workflow …
│   ├─ ai/ …             # chat-ai, predictive-ai …
│   └─ geography/ …      # geo-location
├─ dev/                  # Local-dev artefacts
│   ├─ compose/          # docker-compose*.yml
│   ├─ keycloak/         # realm exports
│   ├─ postgres/         # init SQL
│   ├─ kafka/            # topic scripts
│   ├─ redis/            # redis.conf
│   └─ scripts/          # helper bash files
├─ docs/                 # Architecture, specs, ADRs
├─ Makefile              # 💚 single entry for every task
└─ go.work               # Go workspace file
```

---

## 3. Quick Start – One Command  

```bash
git clone https://github.com/your-org/medical-platform.git
cd medical-platform
make dev-up            # builds & starts entire stack
```

`make dev-up` performs:

1. Validate prerequisites  
2. Build the **single platform image** (`docker-build-platform`)  
3. Launch docker-compose stack  
4. Post-up hooks (Keycloak import, DB schema, Kafka topics, etc.)  

Stop / reset:

```bash
make dev-down          # graceful stop
make dev-reset         # nuke volumes & images (CAUTION)
```

---

## 4. Infrastructure Stack (Docker Compose)  

| Service  | Image                                       | Ports | Volume |
|----------|---------------------------------------------|-------|--------|
| platform | **medical-platform:local** (built locally)  | 8081  | — |
| keycloak | quay.io/keycloak/keycloak:22                | 8080  | kc-db |
| postgres | citusdata/citus:12.1                        | 5432  | pg-data |
| kafka    | confluentinc/cp-kafka:7.5                   | 9092  | kafka-data |
| zookeeper| confluentinc/cp-zookeeper:7.5               | 2181  | — |
| redis    | redis:7-alpine                              | 6379  | redis-data |
| prometheus| prom/prometheus:v2.48                      | 9090  | prom-data |
| grafana  | grafana/grafana:10                          | 3000  | graf-data |
| otel-col | otel/opentelemetry-collector-contrib:0.91   | 4317/4318 | — |
| mailhog  | mailhog/mailhog:latest                      | 8025  | — |

`platform` is the **only business-logic container**.  
Runtime modules are selected via the `ENABLED_MODULES` env (default `*`).

---

## 5. Development Workflow (Single Repo)  

```
make dev-up        # spin infra + platform
make docker-build-platform
go test ./...      # run all unit tests
ENABLED_MODULES="catalog,rfq" make dev-up   # run subset
```

* Hot-reload: use `air`, `reflex` or VS Code Go-tools for local binary rebuild, then `docker compose restart platform`.  
* One CI workflow (`ci.yaml`) runs lint → unit → integration → image → scan.

---

## 6. Keycloak Configuration (Multi-Tenant)  

Unchanged, but **platform** container calls Keycloak directly (no sidecar).  
Use `make kc-add-tenant TENANT=<name>` to duplicate template realm.

---

## 7. PostgreSQL Multi-Tenant Schema  

Same RLS policies.  All modules share the same database server but separate schemas; platform container connects once and each module gets its own repo instance.

---

## 8. Kafka Topic Bootstrap  

`dev/kafka/create-topics.sh` still valid.  Platform container uses **in-process function calls** for same-pod communication; Kafka remains for cross-pod & audit.

---

## 9. Redis Feature Store  

No changes.

---

## 10. Build & Image  

| Command | Result |
|---------|--------|
| `make docker-build-platform` | Multi-stage Go build → `medical-platform:latest` |
| `make docker-build-platform TAG=v0.3.1` | Tagged image |

Single binary located at `/usr/bin/platform` inside image. Size ≈ 60 MB (distroless).

---

## 11. Module-Specific Development  

| Scenario | How-To |
|----------|--------|
| Start only marketplace modules | `ENABLED_MODULES="catalog,rfq,quote,contract" make dev-up` |
| Add new module skeleton | `make init-service SERVICE=parts` |
| Run integration tests for module | `make test-integration SERVICE=rfq` |
| Debug module locally | `make dev-up` then `dlv attach $(docker compose ps -q platform)` |

---

## 12. Testing Infrastructure (Monorepo Pattern)  

| Layer | Tool / Target | Command |
|-------|---------------|---------|
| Unit (all Go) | `go test ./...` | `make test-unit-all` |
| Integration (module) | Testcontainers | `make test-integration SERVICE=ticket` |
| Contract | Pact | `make test-contract` |
| E2E | Playwright (API Gateway) | `make test-e2e` |

Coverage aggregated; gate ≥ 80 %.

---

## 13. Observability  

URLs:

| Tool      | URL                           |
|-----------|------------------------------|
| API Gateway | http://localhost:8081        |
| Keycloak  | http://localhost:8080         |
| Grafana   | http://localhost:3000 (admin/admin) |
| Prometheus| http://localhost:9090         |
| MailHog   | http://localhost:8025         |

Traces displayed in Grafana Tempo data-source (already provisioned).

---

## 14. Deployment References (Single Image)  

Kubernetes example:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: marketplace-pod
spec:
  template:
    spec:
      containers:
      - name: platform
        image: ghcr.io/org/medical-platform:sha-abc123
        env:
        - name: ENABLED_MODULES
          value: "catalog,rfq,quote,contract"
```

Create additional deployments (service-ops, ai) using **same image**, different `ENABLED_MODULES`.

---

## 15. Troubleshooting Cheatsheet (Monorepo)  

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| `platform exited` with `unknown module` | Typo in `ENABLED_MODULES` | Check module list via `make list-services` |
| Hot-reload doesn’t reflect | Binary cached | `docker compose restart platform` |
| `go test ./...` slow | Stale build cache | `go clean -testcache` |
| DB errors across schemas | Missing `set_tenant_context` | Ensure gateway injects tenant header |
| Route 404 | Gateway path not mounted | Confirm module’s `MountRoutes` registration |

---

## 16. Next Steps for Developers  

1. `make dev-up`  
2. Obtain token:

```bash
TOKEN=$(make kc-get-token TENANT=demo-hospital CLIENT_ID=api-gateway \
       CLIENT_SECRET=api-gateway-secret)
```

3. Call API:

```bash
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8081/api/v1/rfqs
```

4. Develop inside module folder, run `go test ./...` and watch CI pass ✨

Happy hacking — and remember: **one repo, one image, many modules**!  
