package api

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"
    "strings"

    "github.com/aby-med/medical-platform/internal/core/organizations/infra"
    "github.com/go-chi/chi/v5"
)

type Handler struct {
    repo   *infra.Repository
    logger *slog.Logger
}

func NewHandler(repo *infra.Repository, logger *slog.Logger) *Handler {
    return &Handler{repo: repo, logger: logger.With(slog.String("component", "org_handler"))}
}

func (h *Handler) ListOrgs(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListOrgs(ctx, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list orgs: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) ListRelationships(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")
    if id == "" {
        h.respondError(w, http.StatusBadRequest, "id required")
        return
    }
    rels, err := h.repo.ListRelationships(ctx, id)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list relationships: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": rels})
}

func (h *Handler) ListChannels(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListChannels(ctx, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list channels: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListProducts(ctx, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list products: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) ListSkus(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListSkus(ctx, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list skus: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

// Phase 2: Offerings + Channel Catalog
func (h *Handler) ListOfferings(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListOfferings(ctx, limit, offset)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to list offerings: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) CreateOffering(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    var req struct{
        SkuID      string  `json:"sku_id"`
        OwnerOrgID *string `json:"owner_org_id"`
        Data       json.RawMessage `json:"data"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid body: "+err.Error())
        return
    }
    o, err := h.repo.CreateOffering(ctx, req.SkuID, req.OwnerOrgID, req.Data)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to create offering: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusCreated, o)
}

func (h *Handler) PublishToChannel(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    channelID := chi.URLParam(r, "id")
    var req struct{ OfferingID string `json:"offering_id"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid body: "+err.Error())
        return
    }
    h.logger.Info("catalog.publish", slog.String("channel_id", channelID), slog.String("offering_id", req.OfferingID))
    if err := h.repo.PublishToChannel(ctx, channelID, req.OfferingID); err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to publish: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]string{"status":"published"})
}

func (h *Handler) UnlistFromChannel(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    channelID := chi.URLParam(r, "id")
    var req struct{ OfferingID string `json:"offering_id"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid body: "+err.Error())
        return
    }
    h.logger.Info("catalog.unlist", slog.String("channel_id", channelID), slog.String("offering_id", req.OfferingID))
    if err := h.repo.UnlistFromChannel(ctx, channelID, req.OfferingID); err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to unlist: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]string{"status":"unlisted"})
}

// Phase 3: Price books + rules + resolve
func (h *Handler) CreatePriceBook(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    var req struct{
        Name string `json:"name"`
        OrgID *string `json:"org_id"`
        ChannelID *string `json:"channel_id"`
        Currency string `json:"currency"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid body: "+err.Error())
        return
    }
    h.logger.Info("pricing.create_price_book", slog.String("name", req.Name))
    if req.Currency == "" { req.Currency = "INR" }
    b, err := h.repo.CreatePriceBook(ctx, req.Name, req.OrgID, req.ChannelID, req.Currency)
    if err != nil { h.respondError(w, http.StatusInternalServerError, "failed to create price book: "+err.Error()); return }
    h.respondJSON(w, http.StatusCreated, b)
}

func (h *Handler) AddPriceRule(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    var req struct{ BookID, SkuID string; Price float64 }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { h.respondError(w, http.StatusBadRequest, "invalid body: "+err.Error()); return }
    h.logger.Info("pricing.add_rule", slog.String("book_id", req.BookID), slog.String("sku_id", req.SkuID))
    if err := h.repo.AddPriceRule(ctx, req.BookID, req.SkuID, req.Price); err != nil { h.respondError(w, http.StatusInternalServerError, "failed to add price rule: "+err.Error()); return }
    h.respondJSON(w, http.StatusOK, map[string]string{"status":"ok"})
}

func (h *Handler) ResolvePrice(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    skuID := r.URL.Query().Get("sku_id")
    orgID := r.URL.Query().Get("org_id")
    channelID := r.URL.Query().Get("channel_id")
    var oPtr, cPtr *string
    if orgID != "" { oPtr = &orgID }
    if channelID != "" { cPtr = &channelID }
    h.logger.Info("pricing.resolve", slog.String("sku_id", skuID), slog.String("org_id", orgID), slog.String("channel_id", channelID))
    res, err := h.repo.ResolvePrice(ctx, skuID, oPtr, cPtr)
    if err != nil { h.respondError(w, http.StatusNotFound, "price not found"); return }
    h.respondJSON(w, http.StatusOK, res)
}

// Phase 5: Engineers
func (h *Handler) ListEngineers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    items, err := h.repo.ListEngineers(ctx, limit, offset)
    if err != nil { h.respondError(w, http.StatusInternalServerError, "failed to list engineers: "+err.Error()); return }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) ListEligibleEngineers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    region := r.URL.Query().Get("region")
    skillsCSV := r.URL.Query().Get("skills")
    var skills []string
    if skillsCSV != "" {
        // split by comma and trim spaces
        raw := strings.Split(skillsCSV, ",")
        for _, s := range raw { if t := strings.TrimSpace(s); t != "" { skills = append(skills, t) } }
    }
    items, err := h.repo.EligibleEngineers(ctx, skills, region, limit)
    if err != nil { h.respondError(w, http.StatusInternalServerError, "failed to compute eligibility: "+err.Error()); return }
    h.respondJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
    h.respondJSON(w, status, map[string]string{"error": message})
}
