package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/services"
)

type API struct {
	Services services.Services
	Config   config.Config
}

func NewRouter(svc services.Services, cfg config.Config) http.Handler {
	api := &API{Services: svc, Config: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthz", api.healthz)
	mux.HandleFunc("/api/v1/auth/", api.authSubroutes)
	mux.HandleFunc("/api/v1/sites", api.sites)
	mux.HandleFunc("/api/v1/sites/", api.siteSubroutes)
	mux.HandleFunc("/api/v1/pages/", api.notImplemented)
	mux.HandleFunc("/api/v1/articles/", api.notImplemented)
	mux.HandleFunc("/api/v1/builds/", api.notImplemented)
	return mux
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (a *API) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"user": nil, "auth": "placeholder"})
}

func (a *API) authSubroutes(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/auth/login":
		a.login(w, r)
	case "/api/v1/auth/logout":
		a.logout(w, r)
	case "/api/v1/auth/me":
		a.me(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (a *API) sites(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": []any{a.Services.ExampleSite()}})
}

func (a *API) siteSubroutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeJSON(w, http.StatusOK, map[string]any{"message": "site subroutes placeholder", "path": r.URL.Path})
		return
	}

	siteID := parts[3]
	if len(parts) == 4 {
		writeJSON(w, http.StatusOK, map[string]any{"message": "site root placeholder", "site_id": siteID})
		return
	}

	switch parts[4] {
	case "pages":
		writeJSON(w, http.StatusOK, map[string]any{"message": "pages placeholder", "site_id": siteID})
	case "articles":
		writeJSON(w, http.StatusOK, map[string]any{"message": "articles placeholder", "site_id": siteID})
	case "sections":
		writeJSON(w, http.StatusOK, map[string]any{"message": "landing sections placeholder", "site_id": siteID})
	case "build":
		writeJSON(w, http.StatusOK, map[string]any{"message": "build requested", "site_id": siteID})
	case "ai":
		if len(parts) > 5 {
			writeJSON(w, http.StatusOK, map[string]any{"message": "ai placeholder", "site_id": siteID, "action": parts[5]})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"message": "ai placeholder", "site_id": siteID})
	default:
		writeJSON(w, http.StatusOK, map[string]any{"message": "site subroutes placeholder", "site_id": siteID, "path": r.URL.Path})
	}
}

func (a *API) notImplemented(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"message": "placeholder endpoint", "path": r.URL.Path})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
