package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/middleware"
	"cms-builder/api/internal/services"
)

type API struct {
	Services services.Services
	Config   config.Config
}

func NewRouter(svc services.Services, cfg config.Config) http.Handler {
	api := &API{Services: svc, Config: cfg}
	protected := middleware.Auth([]byte(cfg.JWTSecret))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthz", api.healthz)
	mux.HandleFunc("/api/v1/auth/login", api.login)
	mux.Handle("/api/v1/auth/logout", protected(http.HandlerFunc(api.logout)))
	mux.Handle("/api/v1/auth/me", protected(http.HandlerFunc(api.me)))
	mux.Handle("/api/v1/workspace", protected(http.HandlerFunc(api.workspace)))
	mux.Handle("/api/v1/sites", protected(http.HandlerFunc(api.sites)))
	mux.Handle("/api/v1/sites/", protected(http.HandlerFunc(api.siteSubroutes)))
	mux.Handle("/api/v1/articles/", protected(http.HandlerFunc(api.articleSubroutes)))
	mux.Handle("/api/v1/builds/", protected(http.HandlerFunc(api.buildSubroutes)))
	mux.Handle("/deployments/", http.StripPrefix("/deployments/", http.FileServer(http.Dir(api.deploymentsRoot()))))
	return mux
}

func (a *API) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) deploymentsRoot() string {
	return filepath.Join("dist", "deployments")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
