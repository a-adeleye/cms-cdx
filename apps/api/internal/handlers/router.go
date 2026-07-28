package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/config"
	"cms-builder/api/internal/middleware"
	"cms-builder/api/internal/services"
)

type API struct {
	Services         services.Services
	Config           config.Config
	previewMutex     sync.Mutex
	templatePreviews map[string]string
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
	mux.Handle("/api/v1/templates", protected(http.HandlerFunc(api.templates)))
	mux.HandleFunc("/api/v1/template-previews/", api.templatePreview)
	mux.Handle("/api/v1/site-imports", protected(http.HandlerFunc(api.siteImports)))
	mux.Handle("/api/v1/sites", protected(http.HandlerFunc(api.sites)))
	mux.Handle("/api/v1/sites/", protected(http.HandlerFunc(api.siteSubroutes)))
	mux.Handle("/api/v1/articles/", protected(http.HandlerFunc(api.articleSubroutes)))
	mux.Handle("/api/v1/builds/", protected(http.HandlerFunc(api.buildSubroutes)))
	mux.Handle("/deployments/", http.StripPrefix("/deployments/", http.FileServer(http.Dir(api.deploymentsRoot()))))
	return mux
}

func (a *API) templatePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/template-previews/"), "/")
	var page string
	switch slug {
	case "default-blog", "premium-saas", "anonime":
		page = builder.RenderTemplatePreview(slug)
	case builder.SupromailTemplateKey:
		var err error
		page, err = a.supromailTemplatePreview(r.Context())
		if err != nil {
			http.Error(w, "template preview is temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
	default:
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; img-src https: data:; font-src https: data:; frame-ancestors 'self'")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}

func (a *API) supromailTemplatePreview(ctx context.Context) (string, error) {
	a.previewMutex.Lock()
	defer a.previewMutex.Unlock()
	if page, found := a.templatePreviews[builder.SupromailTemplateKey]; found {
		return page, nil
	}
	previewer, ok := a.Services.Builder.(builder.TemplatePreviewer)
	if !ok {
		return "", os.ErrNotExist
	}
	outputPath, err := previewer.GenerateTemplatePreview(ctx, builder.SupromailTemplateKey)
	if err != nil {
		return "", err
	}
	page, err := os.ReadFile(filepath.Join(outputPath, "index.html"))
	if err != nil {
		return "", err
	}
	if a.templatePreviews == nil {
		a.templatePreviews = make(map[string]string)
	}
	a.templatePreviews[builder.SupromailTemplateKey] = string(page)
	return string(page), nil
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
