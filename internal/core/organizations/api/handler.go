package api

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"

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
    if err := h.repo.UnlistFromChannel(ctx, channelID, req.OfferingID); err != nil {
        h.respondError(w, http.StatusInternalServerError, "failed to unlist: "+err.Error())
        return
    }
    h.respondJSON(w, http.StatusOK, map[string]string{"status":"unlisted"})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
    h.respondJSON(w, status, map[string]string{"error": message})
}
