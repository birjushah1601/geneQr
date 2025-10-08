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

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
    h.respondJSON(w, status, map[string]string{"error": message})
}
