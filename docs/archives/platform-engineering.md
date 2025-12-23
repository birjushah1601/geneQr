# MASTER-ARCHITECTURE-GUIDE.md  
_Intelligent Medical Equipment Platform_ — **Authoritative Engineering Blueprint**

---

## 1. Architectural Philosophy  

**Guiding Principles**  
1. Domain-Driven Design (DDD) – domain language shapes code.  
2. Hexagonal / Clean Architecture – dependencies point inward.  
3. Modular Monolith – single codebase, modular run-time.  
4. SOLID + KISS – small, composable, test-first units.  
5. Event-Driven – publish immutable domain events.  
6. Security-First – zero-trust, tenant isolation, audited flows.  
7. Observability-Everywhere – logs, metrics, traces as first-class citizens.

---

## 2. Bounded Contexts & Domains  

| Domain | Core Responsibility | Ubiquitous Language | Key Modules |
|--------|--------------------|---------------------|-------------|
| Marketplace | Equipment procurement | RFQ, Quote, Contract, SKU | catalog, rfq, quote, contract |
| Service Ops | Device uptime & field ops | Ticket, Asset, SLA, Workflow | ticket, asset, workflow |
| AI / ML | Intelligence & automation | Prediction, Recommendation, Insight | triage-ai, negotiation-ai, predictive-ai |
| Identity | AuthN / AuthZ | User, Role, Tenant, Realm | external Keycloak |
| Geography | Coverage & routing | Facility, Pincode, Hub | geo-location |

Each bounded context owns its data and publishes events for other contexts.

---

## 3. Service Architecture Pattern — **Monorepo Modules**

```
medical-platform/
├── cmd/
│   ├── platform/            # Single-binary entrypoint (all or selective)
│   └── gateway/             # Stand-alone API gateway (optional)
├── internal/
│   ├── shared/              # Auth, tenant, observability, common VO
│   ├── marketplace/         # ← bounded-context root
│   │   ├── catalog/         #   module
│   │   ├── rfq/
│   │   └── quote/
│   ├── serviceops/
│   │   ├── ticket/
│   │   ├── asset/
│   │   └── workflow/
│   ├── ai/
│   │   ├── triage/
│   │   ├── negotiation/
│   │   └── predictive/
│   └── geography/
├── pkg/                     # Public interfaces (SDK / CLI)
├── scripts/                 # Dev-ops helpers
└── deployments/             # k8s, Docker, Helm templates
```

Module skeleton:

```
internal/<context>/<module>/
├── domain/          # Entities, VOs, domain services
├── app/             # Command/query handlers, DTOs
├── infra/           # DB, MQ, external adapters
├── http/            # REST / gRPC handlers
└── module.go        # Wire dependencies & expose Service interface
```

---

## 4. Modular Monolith Pattern & Benefits  

**Definition**: Single deployable binary/container hosting multiple isolated service modules that communicate _in-process_.  

**Benefits**  
• Shared code, atomic commits → rapid refactorings  
• One CI/CD pipeline & image → infra simplicity  
• No network hops → lower latency, lower cost  
• Still respects bounded contexts via package boundaries + go modules  
• Progressive extraction: any module can later run standalone by toggling ENV  

**Pitfalls Guardrails**  
• Enforce boundaries via internal visibility & lint rules  
• Only cross via published interfaces & domain events  
• No direct DB joins across contexts – use events & projections  

---

## 5. Dependency Management (Single Repo)  

1. **Go Workspaces**: `go work sync` keeps per-module `go.mod` tidy but shares cache.  
2. Central **toolchain version** (`go 1.22`) and shared `Makefile`.  
3. One `golangci-lint` config, one `go.sum` security scan.  
4. Atomic PR ensures schema + code + tests travel together.  

---

## 6. Inter-Service Communication  

| Scenario | Primary Path | Fallback / External |
|----------|--------------|---------------------|
| Module ↔ Module (same binary) | Direct function call via exported Service interface (0-copy) | n/a |
| Gateway ↔ Module | In-process call (mux router delegates) | HTTP/gRPC if module extracted |
| Cross-pod / extracted | gRPC or REST | Kafka domain events |

```
┌──────────┐       Fn Call       ┌──────────┐
│ catalog  │────────────────────►│   rfq    │
└──────────┘                    └──────────┘
(if extracted → HTTP or gRPC)
```

---

## 7. Deployment Strategies  

| Mode | Container Count | ENV example | Use-case |
|------|-----------------|-------------|----------|
| **Monolith** | 1 | `ENABLED_MODULES=*` | Dev, small prod |
| **Selective** | 1 × N (same image) | `ENABLED_MODULES=catalog,rfq` etc. | Domain scaling |
| **Distributed** | Many images (same build) | Extract module binary | Hot paths, high load |

Single Dockerfile builds once; runtime ENV decides what boots.

---

## 8. Configuration-Driven Module Enablement  

```yaml
# config.yaml
enabledModules:
  - catalog
  - rfq
  - ticket
```

```go
func loadEnabled(cfg Config) []service.Module {
    registry := service.AllModules()
    return registry.Filter(cfg.EnabledModules)
}
```

`ENABLED_MODULES="catalog,rfq"` or `"*"` for all.

---

## 9. Code Samples – Registry & Selective Startup  

```go
// internal/service/registry.go
type Module interface {
    Name() string
    MountRoutes(mux *chi.Mux)
    Start(context.Context) error
}

func AllModules() []Module {
    return []Module{
        catalog.New(),
        rfq.New(),
        quote.New(),
        ticket.New(),
        // ...
    }
}

func Boot(ctx context.Context, enabled []string) error {
    mux := chi.NewRouter()
    eg, ctx := errgroup.WithContext(ctx)

    for _, m := range AllModules() {
        if !slice.Contains(enabled, m.Name()) && enabled[0] != "*" {
            continue
        }
        m.MountRoutes(mux)
        eg.Go(func() error { return m.Start(ctx) })
        log.Info("module started", slog.String("module", m.Name()))
    }
    return eg.Wait()
}
```

`cmd/platform/main.go`

```go
func main() {
 cfg := loadConfig()
 ctx := signalctx.New()
 if err := service.Boot(ctx, cfg.EnabledModules); err != nil {
     log.Fatal(err)
 }
}
```

---

## 10. Testing Strategy (Monorepo)  

```
            ┌────────────┐ 5 %  E2E  (Playwright, k6)
            └────────────┘
      ┌─────────────────────┐ 15 % Integration (Testcontainers per module)
      └─────────────────────┘
┌──────────────────────────────┐ 80 % Unit (domain + app pkg)
└──────────────────────────────┘
```

* All tests live under module folder (`*_test.go`).  
* `go test ./...` from repo root runs everything.  
* Tag `integration` uses Docker; skipped in fast CI stage.  
* Coverage aggregated by `go tool cover` and must be ≥ 80 %.  

---

## 11. CI/CD – Single Pipeline  

1. **Lint → Unit → Build → Integration → Image → Scan → Helm Package**  
2. Matrix builds run modules in parallel _within_ pipeline, not separate pipelines.  
3. Provenance: **one SBOM**, **one cosign signature**, **one ArgoCD Application** with ENV overlay selecting run mode.  

GitHub Actions excerpt:

```yaml
strategy:
  matrix: { stage: [lint,test,build] }

steps:
- uses: actions/setup-go@v5
- run: make lint-all
- run: make test-unit-all
- run: make docker-build-platform
```

---

## 12. Data, Security & Observability  

Unchanged from prior version; still DB-per-bounded-context, RLS, Keycloak JWT, Prom/OTel. Only difference: fewer network hops, single sidecar.

---

## 13. How to Extend / Extract a Module  

1. Set `ENABLED_MODULES` to desired list in new Deployment manifest.  
2. (Optional) Compile `cmd/<module>` for slim image; reuse same code.  
3. Update Gateway routing if HTTP boundary introduced.  
4. Keep domain events identical → zero consumer impact.

---

### **Quick Reference**

| Task | Command |
|------|---------|
| Run all modules locally | `make dev-up` (monolith) |
| Run only catalog+rfq | `ENABLED_MODULES=catalog,rfq make dev-up` |
| Generate new module skeleton | `make init-service SERVICE=parts` |
| Run entire test suite | `go test ./...` |
| Build prod image | `make docker-build-platform` |

---

### **TL;DR**

We now operate a **modular monolith**:

• One repo • One image • Config-driven modules  
• DDD boundaries are enforced by Go packages  
• We can scale horizontally by deploying the _same image_ multiple times with different module sets or later extract true microservices — no lock-in.

_All teams must follow this guide. Divergence requires architecture-guild approval._  
