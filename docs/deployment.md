# CICD-DEPLOYMENT-SPECIFICATIONS.md  
_Intelligent Medical Equipment Platform – CI/CD, GitOps & Deployment Blueprint (Monorepo Edition)_  

---

## 1. CI/CD PIPELINE ARCHITECTURE – **Single Repository, Single Image**

```
             ┌────────────┐  Pull-Request   ┌─────────────┐
Developer ──►│  GitHub PR ├────────────────►│  CI  (GHA)  │
             └────────────┘  status-checks  └──────┬──────┘
                                                  │
                   signed image + SBOM            │
                                                  ▼
                        OCI Registry  (ghcr.io/org/medical-platform:<sha>)
                                                  │
                                  GitOps commit (kustomize overlays)
                                                  ▼
                       ┌────────────────────────────────────┐
                       │       ArgoCD (CD Controller)       │
                       └────────────────────────────────────┘
                                  │ auto / manual sync
                 ┌──────────────┐ │            ┌──────────────┐
                 │ dev cluster  │◄┘────────────│ staging cls  │
                 └──────────────┘   promote PR └──────────────┘
                                                     ▲
                                                     │  manual gate
                                                ┌──────────────┐
                                                │  prod cls    │
                                                └──────────────┘
```

*Single source of truth*: same git repo –  
`/deploy/kustomize/<env>/` overlays referencing the **one image**.

---

## 2. BUILD & TEST STAGES (GITHUB ACTIONS – MONOREPO)

```yaml
name: ci
on:
  pull_request:
    paths-ignore: [ "docs/**" ]
  push:
    branches: [ main ]
jobs:
  build-test:
    runs-on: ubuntu-22.04
    permissions: { contents: read, id-token: write }
    steps:
    - uses: actions/checkout@v4

    # Go toolchain
    - uses: actions/setup-go@v5
      with: { go-version: '1.22' }

    # ----- Quality Gates --------------------------------------------------
    - run: make lint-all           # golangci-lint monorepo
    - run: make test-unit-all
    - run: make test-integration-all
    # (optional) run module-matrix fast:
    # strategy: matrix: { module: [catalog,rfq,ticket,…] }

    # ----- Build Single Platform Image -----------------------------------
    - run: make docker-build-platform TAG=${{ github.sha }}

    # SBOM + vulnerability scan
    - run: make docker-scan IMAGE_TAG=${{ github.sha }}

    # ----- Sign & Push ----------------------------------------------------
    - uses: sigstore/cosign-installer@v3
    - run: cosign sign ghcr.io/org/medical-platform:${{ github.sha }} --yes
    - run: docker push ghcr.io/org/medical-platform:${{ github.sha }}

    # ----- Update GitOps manifest ----------------------------------------
    - name: GitOps commit
      uses: EndBug/add-and-commit@v9
      with:
        add: 'deploy/kustomize/**/image-tag.yaml'
        message: 'ci: promote image ${{ github.sha }}'
```

---

## 3. SECURITY & COMPLIANCE SCANNING — **Single Codebase**

| Stage            | Tool                         | Gate |
|------------------|------------------------------|------|
| SAST             | Semgrep (`semgrep/scan-action@v1`) | Critical = 0 |
| Dependency       | Trivy / govulncheck          | Critical = 0 |
| Container Image  | Trivy image scan *(one image)* | CVSS ≥ 7 blocks |
| IaC (Terraform/K8s) | tfsec, kube-linter        | High = fail |
| Secrets          | Gitleaks                     | Any secret = fail |
| SBOM + Sigstore  | cosign + CycloneDX           | Required |
| Compliance Suite | `compliance-check@v1` (HIPAA/DPDP) | Pass |

---

## 4. CONTAINER IMAGE BUILD & SCAN – **Single Platform Image**

```
docker-build-platform
├─ Stage 1: Go compile → /out/platform       (≈60 MB)
└─ Stage 2: distroless/base:nonroot
      /usr/bin/platform   ENTRYPOINT
```

Image tags: `ghcr.io/org/medical-platform:<gitSha>` + `latest-dev`.

Single SBOM (`sbom-${sha}.json`) attached; one cosign signature.

---

## 5. GITOPS WITH ARGOCD – **One App, Many Overlays**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: medical-platform
spec:
  project: default
  source:
    repoURL: https://github.com/org/medical-platform
    targetRevision: main
    path: deploy/kustomize
    plugin:
      name: kustomize-v4
  destination:
    server: https://kubernetes.default.svc
    namespace: platform
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    retry: { limit: 5, backoff: { duration: 20s, factor: 2 } }
```

Folder structure:

```
deploy/kustomize/
  ├─ base/                    # common Deployment, Service, HPA
  ├─ overlays/
  │   ├─ dev/
  │   ├─ staging/
  │   └─ prod/
  └─ image-tag.yaml           # auto-patched by CI
```

---

## 6. DEPLOYMENT STRATEGY – **Single Image, Multiple Configs**

| Pod | ENABLED_MODULES env | Purpose |
|-----|---------------------|---------|
| marketplace-pod | `catalog,rfq,quote,contract` | Procurement |
| serviceops-pod  | `asset,device-reg,qr,ticket,workflow` | Field operations |
| ai-pod          | `chat-ai,negotiation-ai,predictive-maint,dispatch-ai,demand-forecast` | Intelligence |

All pods use **same container image** with different env config.

Rollout pattern: Argo Rollouts canary 10-50-100 across **each pod**.

---

## 7. ROLLBACK & PROMOTION

*Rollback*: change `image-tag.yaml` back to previous SHA → Argo sync — single action rolls back all pods (consistent version set).

*Promotion flow*:

1. PR to `main` builds image → auto-deploy to **dev** overlay.  
2. Merge `promote/staging` PR updates staging tag.  
3. CAB-approved PR to `promote/prod` updates prod tag.

---

## 8. PIPELINE COMPLIANCE – MONOREPO

Checks run once per PR on full codebase:

```
✓ Lint         ✓ Unit ≥80 %
✓ Integration  ✓ SAST clean
✓ SBOM         ✓ Image Scan
✓ Compliance   ✓ Signed image
```

Artifact bundle (zip) of test reports + compliance evidence uploaded to AWS Audit Bucket (immutable).

---

## 9. OBSERVABILITY & ALERTING (UNCHANGED)

Prometheus/Grafana/Tempo/Loki stack auto-deployed via Helm chart in `deploy/observability/`.

Default SLO dashboards target:
* p95 latency ≤300 ms for Gateway
* error-rate < 0.1 %

Alertmanager hooks PagerDuty for sev-1.

---

## 10. FUTURE ENHANCEMENTS

* Parallel CI cache optimisation by Go build cache restore  
* Progressive Delivery via Flagger on single image  
* OPA Gatekeeper enforcing module list vs namespace  

---

_This specification supersedes prior multi-service pipeline docs. Any deviation requires DevOps lead approval._  
