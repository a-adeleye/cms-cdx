package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"cms-builder/api/internal/ai"
	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/media"
	"cms-builder/api/internal/middleware"
	"cms-builder/api/internal/models"
	"cms-builder/api/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

type workspaceResponse struct {
	User              userResponse         `json:"user"`
	SelectedSiteID    string               `json:"selectedSiteId"`
	SelectedArticleID string               `json:"selectedArticleId"`
	Sites             []siteResponse       `json:"sites"`
	Templates         []templateResponse   `json:"templates"`
	LandingSections   []landingSectionResp `json:"landingSections"`
	Articles          []articleResponse    `json:"articles"`
	Authors           []authorResponse     `json:"authors"`
	Categories        []categoryResponse   `json:"categories"`
	MediaAssets       []mediaAssetResponse `json:"mediaAssets"`
	Builds            []buildResponse      `json:"builds"`
}

type siteResponse struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Slug                  string   `json:"slug"`
	Domain                string   `json:"domain"`
	BlogPath              string   `json:"blogPath"`
	Description           string   `json:"description"`
	ContentContext        string   `json:"contentContext"`
	LogoMediaID           string   `json:"logoMediaId"`
	LogoURL               string   `json:"logoUrl"`
	FaviconMediaID        string   `json:"faviconMediaId"`
	FaviconURL            string   `json:"faviconUrl"`
	Status                string   `json:"status"`
	TemplateKey           string   `json:"templateKey"`
	ThemeConfig           string   `json:"themeConfig"`
	DeployProvider        string   `json:"deployProvider"`
	DeployConfig          string   `json:"deployConfig"`
	PreviewDeployProvider string   `json:"previewDeployProvider"`
	PreviewDeployConfig   string   `json:"previewDeployConfig"`
	AIConfig              string   `json:"aiConfig"`
	StorageConfig         string   `json:"storageConfig"`
	DeploymentWarnings    []string `json:"deploymentWarnings,omitempty"`
	UpdatedAt             string   `json:"updatedAt"`
}

type templateResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	UpdatedAt  string `json:"updatedAt"`
	PreviewURL string `json:"previewUrl"`
}

type landingSectionResp struct {
	ID           string `json:"id"`
	SiteID       string `json:"siteId"`
	SectionKey   string `json:"sectionKey"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	ContentJSON  string `json:"contentJson"`
	DisplayOrder int    `json:"displayOrder"`
	IsEnabled    bool   `json:"isEnabled"`
}

type articleResponse struct {
	ID              string  `json:"id"`
	SiteID          string  `json:"siteId"`
	AuthorID        string  `json:"authorId"`
	CategoryID      string  `json:"categoryId"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Excerpt         string  `json:"excerpt"`
	ContentMarkdown string  `json:"contentMarkdown"`
	CoverImageURL   string  `json:"coverImageUrl"`
	Status          string  `json:"status"`
	IsFeatured      bool    `json:"isFeatured"`
	PublishedAt     *string `json:"publishedAt"`
	SEOTitle        string  `json:"seoTitle"`
	SEODescription  string  `json:"seoDescription"`
	CanonicalURL    string  `json:"canonicalUrl"`
	GeneratedByAI   bool    `json:"generatedByAi"`
	HumanReviewed   bool    `json:"humanReviewed"`
	AIPrompt        string  `json:"aiPrompt"`
	AIModel         string  `json:"aiModel"`
	Tags            string  `json:"tags"`
	UpdatedAt       string  `json:"updatedAt"`
}

type authorResponse struct {
	ID     string `json:"id"`
	SiteID string `json:"siteId"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Bio    string `json:"bio"`
}

type categoryResponse struct {
	ID          string `json:"id"`
	SiteID      string `json:"siteId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type mediaAssetResponse struct {
	ID              string `json:"id"`
	SiteID          string `json:"siteId"`
	FileName        string `json:"fileName"`
	FileURL         string `json:"fileUrl"`
	MimeType        string `json:"mimeType"`
	SizeBytes       int64  `json:"sizeBytes"`
	StorageProvider string `json:"storageProvider"`
	StorageKey      string `json:"storageKey"`
	AltText         string `json:"altText"`
	CreatedAt       string `json:"createdAt"`
}

type buildResponse struct {
	ID             string  `json:"id"`
	SiteID         string  `json:"siteId"`
	Status         string  `json:"status"`
	BuildType      string  `json:"buildType"`
	Logs           string  `json:"logs"`
	OutputPath     string  `json:"outputPath"`
	DeployProvider string  `json:"deployProvider"`
	DeployStatus   string  `json:"deployStatus"`
	DeployURL      string  `json:"deployUrl"`
	DeployRevision string  `json:"deployRevision"`
	StartedAt      *string `json:"startedAt"`
	FinishedAt     *string `json:"finishedAt"`
}

type siteUpsertRequest struct {
	Name                  string `json:"name"`
	Slug                  string `json:"slug"`
	Domain                string `json:"domain"`
	BlogPath              string `json:"blogPath"`
	Description           string `json:"description"`
	ContentContext        string `json:"contentContext"`
	LogoMediaID           string `json:"logoMediaId"`
	FaviconMediaID        string `json:"faviconMediaId"`
	Status                string `json:"status"`
	TemplateKey           string `json:"templateKey"`
	ThemeConfig           string `json:"themeConfig"`
	DeployProvider        string `json:"deployProvider"`
	DeployConfig          string `json:"deployConfig"`
	PreviewDeployProvider string `json:"previewDeployProvider"`
	PreviewDeployConfig   string `json:"previewDeployConfig"`
	AIConfig              string `json:"aiConfig"`
	StorageConfig         string `json:"storageConfig"`
}

type articleUpsertRequest struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Excerpt         string `json:"excerpt"`
	ContentMarkdown string `json:"contentMarkdown"`
	CoverImageURL   string `json:"coverImageUrl"`
	SEOTitle        string `json:"seoTitle"`
	SEODescription  string `json:"seoDescription"`
	CanonicalURL    string `json:"canonicalUrl"`
	AuthorID        string `json:"authorId"`
	CategoryID      string `json:"categoryId"`
	Tags            string `json:"tags"`
	IsFeatured      bool   `json:"isFeatured"`
	Status          string `json:"status"`
}

type categoryUpsertRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type authorUpsertRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type landingReorderRequest struct {
	SectionIDs []string `json:"sectionIds"`
}

type landingUpdateRequest struct {
	IsEnabled    *bool `json:"isEnabled"`
	DisplayOrder *int  `json:"displayOrder"`
}

type mediaUpsertRequest struct {
	FileName        string `json:"fileName"`
	FileURL         string `json:"fileUrl"`
	MimeType        string `json:"mimeType"`
	SizeBytes       int64  `json:"sizeBytes"`
	StorageProvider string `json:"storageProvider"`
	StorageKey      string `json:"storageKey"`
	AltText         string `json:"altText"`
}

type mediaUpdateRequest struct {
	AltText string `json:"altText"`
}

const (
	maxMediaUploadBytes  = 12 << 20
	maxMediaAltTextRunes = 500
)

type buildCreateRequest struct {
	BuildType  string   `json:"buildType"`
	ArticleIDs []string `json:"articleIds"`
}

func (a *API) workspace(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := a.findUserByID(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load workspace"})
		return
	}

	siteID := strings.TrimSpace(r.URL.Query().Get("siteId"))
	if siteID == "" {
		siteID, err = a.firstSiteID(r.Context())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load workspace"})
			return
		}
	}

	response, err := a.loadWorkspace(r.Context(), claims.UserID, siteID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load workspace"})
		return
	}
	response.User = toUserResponse(*user)
	writeJSON(w, http.StatusOK, response)
}

func (a *API) templates(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := a.listTemplates(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load templates"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	case http.MethodPost:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "templates must be implemented in the production renderer"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (a *API) sites(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		sites, err := a.listSites(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load sites"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": sites})
	case http.MethodPost:
		if !requireAdmin(w, r) {
			return
		}
		var payload siteUpsertRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		site, err := a.createSite(r.Context(), payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to create site"})
			return
		}
		writeJSON(w, http.StatusCreated, site)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) siteSubroutes(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, "/api/v1/sites/")
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	siteID := parts[0]
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			site, err := a.getSite(r.Context(), siteID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load site"})
				return
			}
			writeJSON(w, http.StatusOK, site)
		case http.MethodPatch:
			if !requireAdmin(w, r) {
				return
			}
			var payload siteUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			site, err := a.updateSite(r.Context(), siteID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update site"})
				return
			}
			writeJSON(w, http.StatusOK, site)
		case http.MethodDelete:
			if !requireAdmin(w, r) {
				return
			}
			if err := a.deleteSite(r.Context(), siteID); err != nil {
				if errors.Is(err, errConflict) {
					writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete site"})
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch parts[1] {
	case "export":
		a.handleSiteExport(w, r, siteID, parts[2:])
	case "workspace":
		response, err := a.loadWorkspace(r.Context(), "", siteID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load workspace"})
			return
		}
		writeJSON(w, http.StatusOK, response)
	case "landing-sections":
		a.handleLandingSectionRoutes(w, r, siteID, parts[2:])
	case "articles":
		a.handleSiteArticleRoutes(w, r, siteID, parts[2:])
	case "authors":
		a.handleSiteAuthorRoutes(w, r, siteID, parts[2:])
	case "categories":
		a.handleSiteCategoryRoutes(w, r, siteID, parts[2:])
	case "builds":
		a.handleBuildRoutes(w, r, siteID)
	case "media":
		a.handleMediaRoutes(w, r, siteID, parts[2:])
	case "ai":
		a.handleAISuggestion(w, r, siteID, parts[2:])
	default:
		http.NotFound(w, r)
	}
}

func (a *API) articleSubroutes(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, "/api/v1/articles/")
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	articleID := parts[0]
	switch r.Method {
	case http.MethodGet:
		article, err := a.getArticle(r.Context(), articleID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load article"})
			return
		}
		writeJSON(w, http.StatusOK, article)
	case http.MethodDelete:
		if err := a.deleteArticle(r.Context(), articleID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete article"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPatch:
		var payload articleUpsertRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		payload.ID = articleID
		article, err := a.upsertArticle(r.Context(), payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update article"})
			return
		}
		writeJSON(w, http.StatusOK, article)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) buildSubroutes(w http.ResponseWriter, r *http.Request) {
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	rest := strings.TrimPrefix(r.URL.Path, "/api/v1/builds/")
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	buildID := parts[0]
	switch r.Method {
	case http.MethodGet:
		build, err := a.getBuild(r.Context(), buildID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load build"})
			return
		}
		writeJSON(w, http.StatusOK, build)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleLandingSectionRoutes(w http.ResponseWriter, r *http.Request, siteID string, parts []string) {
	switch {
	case len(parts) == 0 || parts[0] == "":
		switch r.Method {
		case http.MethodGet:
			sections, err := a.listLandingSections(r.Context(), siteID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load landing sections"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": sections})
		case http.MethodPut:
			var payload landingReorderRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			sections, err := a.reorderLandingSections(r.Context(), siteID, payload.SectionIDs)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update landing sections"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": sections})
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	case parts[0] == "reorder":
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	default:
		sectionID := parts[0]
		if r.Method != http.MethodPatch {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload landingUpdateRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		section, err := a.updateLandingSection(r.Context(), siteID, sectionID, payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update landing section"})
			return
		}
		writeJSON(w, http.StatusOK, section)
	}
}

func (a *API) handleSiteArticleRoutes(w http.ResponseWriter, r *http.Request, siteID string, parts []string) {
	switch {
	case len(parts) == 0 || parts[0] == "":
		switch r.Method {
		case http.MethodGet:
			articles, err := a.listArticles(r.Context(), siteID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load articles"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": articles})
		case http.MethodPost:
			var payload articleUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			article, err := a.upsertArticleWithSite(r.Context(), siteID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to save article"})
				return
			}
			writeJSON(w, http.StatusCreated, article)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleSiteAuthorRoutes(w http.ResponseWriter, r *http.Request, siteID string, parts []string) {
	switch {
	case len(parts) == 0 || parts[0] == "":
		switch r.Method {
		case http.MethodGet:
			authors, err := a.listAuthors(r.Context(), siteID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load authors"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": authors})
		case http.MethodPost:
			var payload authorUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			author, err := a.createAuthor(r.Context(), siteID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to create author"})
				return
			}
			writeJSON(w, http.StatusCreated, author)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		authorID := parts[0]
		switch r.Method {
		case http.MethodPatch:
			var payload authorUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			author, err := a.updateAuthor(r.Context(), siteID, authorID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update author"})
				return
			}
			writeJSON(w, http.StatusOK, author)
		case http.MethodDelete:
			if err := a.deleteAuthor(r.Context(), siteID, authorID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete author"})
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (a *API) handleSiteCategoryRoutes(w http.ResponseWriter, r *http.Request, siteID string, parts []string) {
	switch {
	case len(parts) == 0 || parts[0] == "":
		switch r.Method {
		case http.MethodGet:
			categories, err := a.listCategories(r.Context(), siteID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load categories"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": categories})
		case http.MethodPost:
			var payload categoryUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			category, err := a.createCategory(r.Context(), siteID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to create category"})
				return
			}
			writeJSON(w, http.StatusCreated, category)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		categoryID := parts[0]
		switch r.Method {
		case http.MethodPatch:
			var payload categoryUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			category, err := a.updateCategory(r.Context(), siteID, categoryID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update category"})
				return
			}
			writeJSON(w, http.StatusOK, category)
		case http.MethodDelete:
			if err := a.deleteCategory(r.Context(), siteID, categoryID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete category"})
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (a *API) handleBuildRoutes(w http.ResponseWriter, r *http.Request, siteID string) {
	switch r.Method {
	case http.MethodGet:
		builds, err := a.listBuilds(r.Context(), siteID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load builds"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": builds})
	case http.MethodPost:
		if !requireAdmin(w, r) {
			return
		}
		var payload buildCreateRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		build, err := a.createBuild(r.Context(), siteID, payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create build: %v", err)})
			return
		}
		writeJSON(w, http.StatusCreated, build)
	case http.MethodDelete:
		if !requireAdmin(w, r) {
			return
		}
		if err := a.clearBuildHistory(r.Context(), siteID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to clear deployment history"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleMediaRoutes(w http.ResponseWriter, r *http.Request, siteID string, routeParts ...[]string) {
	parts := []string{}
	if len(routeParts) > 0 {
		parts = routeParts[0]
	}
	if len(parts) > 1 {
		http.NotFound(w, r)
		return
	}
	if len(parts) == 1 && parts[0] != "" {
		a.handleMediaAssetRoute(w, r, siteID, parts[0])
		return
	}

	switch r.Method {
	case http.MethodGet:
		media, err := a.listMediaAssets(r.Context(), siteID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load media assets"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": media})
	case http.MethodPost:
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			fileName, contents, mimeType, altText, status, err := readMediaUpload(w, r)
			if err != nil {
				writeJSON(w, status, map[string]string{"error": err.Error()})
				return
			}
			media, err := a.uploadMediaAsset(r.Context(), siteID, fileName, contents, mimeType, altText)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to upload media asset"})
				return
			}
			writeJSON(w, http.StatusCreated, media)
			return
		}

		var payload mediaUpsertRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		media, err := a.createMediaAsset(r.Context(), siteID, payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to create media asset"})
			return
		}
		writeJSON(w, http.StatusCreated, media)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleMediaAssetRoute(w http.ResponseWriter, r *http.Request, siteID, assetID string) {
	switch r.Method {
	case http.MethodPatch:
		var payload mediaUpdateRequest
		if err := decodeJSON(r, &payload); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		media, err := a.updateMediaAsset(r.Context(), siteID, assetID, payload)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update media asset"})
			return
		}
		writeJSON(w, http.StatusOK, media)
	case http.MethodPut:
		fileName, contents, mimeType, altText, status, err := readMediaUpload(w, r)
		if err != nil {
			writeJSON(w, status, map[string]string{"error": err.Error()})
			return
		}
		media, err := a.replaceMediaAsset(r.Context(), siteID, assetID, fileName, contents, mimeType, altText)
		if err != nil {
			if errors.Is(err, errValidation) {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to replace media asset"})
			return
		}
		writeJSON(w, http.StatusOK, media)
	case http.MethodDelete:
		if err := a.deleteMediaAsset(r.Context(), siteID, assetID); err != nil {
			if errors.Is(err, errConflict) {
				writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
				return
			}
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete media asset"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func readMediaUpload(w http.ResponseWriter, r *http.Request) (string, []byte, string, string, int, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxMediaUploadBytes)
	if err := r.ParseMultipartForm(maxMediaUploadBytes); err != nil {
		return "", nil, "", "", http.StatusBadRequest, errors.New("invalid multipart payload")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", nil, "", "", http.StatusBadRequest, errors.New("file is required")
	}
	defer file.Close()

	contents, err := io.ReadAll(io.LimitReader(file, maxMediaUploadBytes+1))
	if err != nil {
		return "", nil, "", "", http.StatusInternalServerError, errors.New("unable to read uploaded file")
	}
	if len(contents) > maxMediaUploadBytes {
		return "", nil, "", "", http.StatusRequestEntityTooLarge, errors.New("file is too large")
	}
	if len(contents) == 0 {
		return "", nil, "", "", http.StatusBadRequest, errors.New("file is empty")
	}

	mimeType := http.DetectContentType(contents)
	allowedImageTypes := map[string]bool{
		"image/gif": true, "image/jpeg": true, "image/png": true, "image/webp": true, "image/x-icon": true,
	}
	if !allowedImageTypes[mimeType] {
		return "", nil, "", "", http.StatusBadRequest, errors.New("only image uploads are supported")
	}
	extensionTypes := map[string]string{
		".gif": "image/gif", ".jpeg": "image/jpeg", ".jpg": "image/jpeg", ".png": "image/png", ".webp": "image/webp", ".ico": "image/x-icon",
	}
	if extensionTypes[strings.ToLower(filepath.Ext(header.Filename))] != mimeType {
		return "", nil, "", "", http.StatusBadRequest, errors.New("file extension does not match image content")
	}

	return header.Filename, contents, mimeType, strings.TrimSpace(r.FormValue("altText")), http.StatusOK, nil
}

var (
	errValidation = errors.New("validation error")
	errConflict   = errors.New("conflict")
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "23505"
}

func (a *API) findUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, email, password_hash, COALESCE(full_name, ''), role, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *API) findUserByID(ctx context.Context, id string) (*models.User, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, email, password_hash, COALESCE(full_name, ''), role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *API) firstSiteID(ctx context.Context) (string, error) {
	row := a.Services.DB.QueryRowContext(ctx, `SELECT id::text FROM sites ORDER BY name ASC LIMIT 1`)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (a *API) loadWorkspace(ctx context.Context, userID, selectedSiteID string) (workspaceResponse, error) {
	sites, err := a.listSites(ctx)
	if err != nil {
		return workspaceResponse{}, err
	}
	templates, err := a.listTemplates(ctx)
	if err != nil {
		return workspaceResponse{}, err
	}

	if selectedSiteID == "" && len(sites) > 0 {
		selectedSiteID = sites[0].ID
	}

	articles, err := a.listArticles(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}
	landingSections, err := a.listLandingSections(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}
	authors, err := a.listAuthors(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}
	categories, err := a.listCategories(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}
	mediaAssets, err := a.listMediaAssets(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}
	builds, err := a.listBuilds(ctx, selectedSiteID)
	if err != nil {
		return workspaceResponse{}, err
	}

	selectedArticleID := ""
	if len(articles) > 0 {
		selectedArticleID = articles[0].ID
	}

	return workspaceResponse{
		SelectedSiteID:    selectedSiteID,
		SelectedArticleID: selectedArticleID,
		Sites:             sites,
		Templates:         templates,
		LandingSections:   landingSections,
		Articles:          articles,
		Authors:           authors,
		Categories:        categories,
		MediaAssets:       mediaAssets,
		Builds:            builds,
	}, nil
}

func (a *API) listSites(ctx context.Context) ([]siteResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			id::text,
			name,
			slug,
			COALESCE(domain, ''),
			blog_path,
			COALESCE(description, ''),
			content_context,
			COALESCE(logo_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.logo_media_id), ''),
			COALESCE(favicon_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.favicon_media_id), ''),
			status,
			template_key,
			COALESCE(theme_config::text, '{}'),
			COALESCE(deploy_provider, ''),
			COALESCE(deploy_config::text, '{}'),
			COALESCE(preview_deploy_provider, ''),
			COALESCE(preview_deploy_config::text, '{}'),
			COALESCE(ai_config::text, '{}'),
			COALESCE(storage_config::text, '{}'),
			updated_at
		FROM sites
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sites := make([]siteResponse, 0)
	for rows.Next() {
		var site siteResponse
		var updatedAt time.Time
		if err := rows.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Description, &site.ContentContext, &site.LogoMediaID, &site.LogoURL, &site.FaviconMediaID, &site.FaviconURL, &site.Status, &site.TemplateKey, &site.ThemeConfig, &site.DeployProvider, &site.DeployConfig, &site.PreviewDeployProvider, &site.PreviewDeployConfig, &site.AIConfig, &site.StorageConfig, &updatedAt); err != nil {
			return nil, err
		}
		site.UpdatedAt = updatedAt.Format(time.RFC3339)
		site.DeploymentWarnings = deploymentWarningsForSite(site)
		sites = append(sites, site)
	}

	return sites, rows.Err()
}

func (a *API) getSite(ctx context.Context, siteID string) (siteResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT
			id::text,
			name,
			slug,
			COALESCE(domain, ''),
			blog_path,
			COALESCE(description, ''),
			content_context,
			COALESCE(logo_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.logo_media_id), ''),
			COALESCE(favicon_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.favicon_media_id), ''),
			status,
			template_key,
			COALESCE(theme_config::text, '{}'),
			COALESCE(deploy_provider, ''),
			COALESCE(deploy_config::text, '{}'),
			COALESCE(preview_deploy_provider, ''),
			COALESCE(preview_deploy_config::text, '{}'),
			COALESCE(ai_config::text, '{}'),
			COALESCE(storage_config::text, '{}'),
			updated_at
		FROM sites
		WHERE id = $1
	`, siteID)

	var site siteResponse
	var updatedAt time.Time
	if err := row.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Description, &site.ContentContext, &site.LogoMediaID, &site.LogoURL, &site.FaviconMediaID, &site.FaviconURL, &site.Status, &site.TemplateKey, &site.ThemeConfig, &site.DeployProvider, &site.DeployConfig, &site.PreviewDeployProvider, &site.PreviewDeployConfig, &site.AIConfig, &site.StorageConfig, &updatedAt); err != nil {
		return siteResponse{}, err
	}
	site.UpdatedAt = updatedAt.Format(time.RFC3339)
	site.DeploymentWarnings = deploymentWarningsForSite(site)
	return site, nil
}

func (a *API) listTemplates(ctx context.Context) ([]templateResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT id::text, name, slug, updated_at
		FROM templates
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]templateResponse, 0)
	for rows.Next() {
		var item templateResponse
		var updatedAt time.Time
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug, &updatedAt); err != nil {
			return nil, err
		}
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		item.PreviewURL = "/api/v1/template-previews/" + item.Slug
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) createSite(ctx context.Context, payload siteUpsertRequest) (siteResponse, error) {
	if err := validateSitePayload(payload); err != nil {
		return siteResponse{}, err
	}
	if err := a.templateExists(ctx, fallbackString(payload.TemplateKey, "default-blog")); err != nil {
		return siteResponse{}, err
	}

	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return siteResponse{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var siteID string
	var updatedAt time.Time
	err = tx.QueryRowContext(ctx, `
		INSERT INTO sites (
			name, slug, domain, blog_path, description, content_context, status, template_key, theme_config, deploy_provider, deploy_config, preview_deploy_provider, preview_deploy_config, ai_config, storage_config
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, $8, $9::jsonb, NULLIF($10, ''), $11::jsonb, NULLIF($12, ''), $13::jsonb, $14::jsonb, $15::jsonb)
		RETURNING id::text, updated_at
	`, payload.Name, payload.Slug, payload.Domain, fallbackString(payload.BlogPath, "/articles"), strings.TrimSpace(payload.Description), siteContentContext(payload.ContentContext), fallbackString(payload.Status, "active"), fallbackString(payload.TemplateKey, "default-blog"), fallbackJSON(payload.ThemeConfig, `{"tone":"professional"}`), payload.DeployProvider, fallbackJSON(payload.DeployConfig, `{}`), payload.PreviewDeployProvider, fallbackJSON(payload.PreviewDeployConfig, `{}`), fallbackJSON(payload.AIConfig, `{}`), fallbackJSON(payload.StorageConfig, `{}`)).Scan(&siteID, &updatedAt)
	if err != nil {
		return siteResponse{}, err
	}

	if err := seedSiteDefaults(ctx, tx, siteID); err != nil {
		return siteResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return siteResponse{}, err
	}

	return a.getSite(ctx, siteID)
}

func (a *API) updateSite(ctx context.Context, siteID string, payload siteUpsertRequest) (siteResponse, error) {
	if err := validateSitePayload(payload); err != nil {
		return siteResponse{}, err
	}
	if err := a.templateExists(ctx, fallbackString(payload.TemplateKey, "default-blog")); err != nil {
		return siteResponse{}, err
	}
	if err := a.validateBrandingAssets(ctx, siteID, payload.LogoMediaID, payload.FaviconMediaID); err != nil {
		return siteResponse{}, err
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE sites
		SET
			name = $2,
			slug = $3,
			domain = NULLIF($4, ''),
			blog_path = $5,
			description = $6,
			content_context = COALESCE(NULLIF($7, ''), content_context),
			logo_media_id = NULLIF($8, '')::uuid,
			favicon_media_id = NULLIF($9, '')::uuid,
			status = $10,
			template_key = $11,
			theme_config = $12::jsonb,
			deploy_provider = NULLIF($13, ''),
			deploy_config = $14::jsonb,
			preview_deploy_provider = NULLIF($15, ''),
			preview_deploy_config = $16::jsonb,
			ai_config = $17::jsonb,
			storage_config = $18::jsonb,
			updated_at = NOW()
		WHERE id = $1
	`, siteID, payload.Name, payload.Slug, payload.Domain, fallbackString(payload.BlogPath, "/articles"), strings.TrimSpace(payload.Description), strings.TrimSpace(payload.ContentContext), payload.LogoMediaID, payload.FaviconMediaID, fallbackString(payload.Status, "active"), fallbackString(payload.TemplateKey, "default-blog"), fallbackJSON(payload.ThemeConfig, `{}`), payload.DeployProvider, fallbackJSON(payload.DeployConfig, `{}`), payload.PreviewDeployProvider, fallbackJSON(payload.PreviewDeployConfig, `{}`), fallbackJSON(payload.AIConfig, `{}`), fallbackJSON(payload.StorageConfig, `{}`))
	if err != nil {
		return siteResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return siteResponse{}, sql.ErrNoRows
	}
	return a.getSite(ctx, siteID)
}

func (a *API) deleteSite(ctx context.Context, siteID string) error {
	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var siteCount int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sites`).Scan(&siteCount); err != nil {
		return err
	}
	if siteCount <= 1 {
		return fmt.Errorf("%w: create another site before deleting the last site", errConflict)
	}
	// Keep the administrative history while removing the foreign-key reference
	// that would otherwise prevent a site and its content from being deleted.
	if _, err = tx.ExecContext(ctx, `UPDATE audit_logs SET site_id = NULL WHERE site_id = $1`, siteID); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM sites WHERE id = $1`, siteID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func validateSitePayload(payload siteUpsertRequest) error {
	if strings.TrimSpace(payload.Name) == "" || strings.TrimSpace(payload.Slug) == "" {
		return fmt.Errorf("%w: name and slug are required", errValidation)
	}
	if _, err := models.CanonicalBlogPath(payload.BlogPath); err != nil {
		return fmt.Errorf("%w: %v", errValidation, err)
	}
	if context := strings.TrimSpace(payload.ContentContext); context != "" && !models.IsValidSiteContentContext(context) {
		return fmt.Errorf("%w: contentContext must be application_blog or standalone_blog", errValidation)
	}
	if err := validateAIConfig(payload.AIConfig); err != nil {
		return err
	}
	if err := validateStorageConfig(payload.StorageConfig); err != nil {
		return err
	}
	if err := validateDeploymentConfig(payload.DeployProvider, payload.DeployConfig); err != nil {
		return err
	}
	return validateDeploymentConfig(payload.PreviewDeployProvider, payload.PreviewDeployConfig)
}

func validateAIConfig(rawConfig string) error {
	if _, err := ai.ParseConfig(fallbackJSON(rawConfig, `{}`)); err != nil {
		return fmt.Errorf("%w: %v", errValidation, err)
	}
	return nil
}

func validateStorageConfig(rawConfig string) error {
	if _, err := storage.ParseConfig(fallbackJSON(rawConfig, `{}`)); err != nil {
		return fmt.Errorf("%w: %v", errValidation, err)
	}
	return nil
}

func siteContentContext(value string) string {
	if strings.TrimSpace(value) == "" {
		return string(models.SiteContentContextStandaloneBlog)
	}
	return strings.TrimSpace(value)
}

func validateDeploymentConfig(provider, rawConfig string) error {
	provider = fallbackString(strings.TrimSpace(provider), "none")
	var values map[string]any
	if err := json.Unmarshal([]byte(fallbackJSON(rawConfig, `{}`)), &values); err != nil {
		return fmt.Errorf("%w: deployment config must be JSON", errValidation)
	}
	switch provider {
	case "none":
		return nil
	case "cloudflare_pages":
		return validateCloudflarePagesConfig(provider, rawConfig)
	case "firebase":
		if configValue(values, "siteId") == "" && configValue(values, "projectId") == "" {
			return fmt.Errorf("%w: Firebase deployment requires siteId or projectId", errValidation)
		}
		if configValue(values, "serviceAccountSecretRef") == "" {
			return fmt.Errorf("%w: Firebase deployment requires serviceAccountSecretRef", errValidation)
		}
		return nil
	case "git_repository":
		repositoryURL := configValue(values, "repositoryUrl")
		parsed, err := url.Parse(repositoryURL)
		if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Hostname() == "" {
			return fmt.Errorf("%w: repositoryUrl must be an HTTPS URL without credentials", errValidation)
		}
		if configValue(values, "branch") == "" || configValue(values, "contentPath") == "" {
			return fmt.Errorf("%w: repository deployment requires branch and contentPath", errValidation)
		}
		if configValue(values, "tokenSecretRef") == "" {
			return fmt.Errorf("%w: repository deployment requires tokenSecretRef", errValidation)
		}
		if _, err := deploySafeRelativePath(configValue(values, "contentPath")); err != nil {
			return fmt.Errorf("%w: %v", errValidation, err)
		}
		return nil
	default:
		return fmt.Errorf("%w: unsupported deployment provider %q", errValidation, provider)
	}
}

func deploySafeRelativePath(value string) (string, error) {
	cleaned := filepath.Clean(filepath.FromSlash(strings.TrimSpace(value)))
	if cleaned == "." || filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || strings.Contains(cleaned, ".git") {
		return "", errors.New("contentPath must be a safe relative directory")
	}
	return cleaned, nil
}

func isCloudflarePagesProvider(provider string) bool {
	return strings.EqualFold(strings.TrimSpace(provider), "cloudflare_pages")
}

func validateCloudflarePagesConfig(provider, rawConfig string) error {
	if !isCloudflarePagesProvider(provider) {
		return nil
	}

	var deployConfig map[string]any
	if err := json.Unmarshal([]byte(fallbackJSON(rawConfig, `{}`)), &deployConfig); err != nil {
		return fmt.Errorf("%w: deployment config must be JSON", errValidation)
	}
	for key := range deployConfig {
		if key != "projectName" && key != "productionBranch" {
			return fmt.Errorf("%w: Cloudflare deployment config only supports projectName and productionBranch", errValidation)
		}
	}
	projectName, isString := deployConfig["projectName"].(string)
	if !isString {
		return fmt.Errorf("%w: Cloudflare deployment config requires a valid projectName", errValidation)
	}
	if branch, ok := deployConfig["productionBranch"]; ok {
		value, isString := branch.(string)
		if !isString || strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("%w: Cloudflare deployment config has an invalid productionBranch", errValidation)
		}
	}
	if !cloudflareProjectName.MatchString(strings.TrimSpace(projectName)) {
		return fmt.Errorf("%w: Cloudflare deployment config requires a valid projectName", errValidation)
	}
	return nil
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil || claims.Role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "administrator access is required"})
		return false
	}
	return true
}

var cloudflareProjectName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,58}[a-z0-9]$|^[a-z0-9]$`)

func (a *API) templateExists(ctx context.Context, slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("%w: template key is required", errValidation)
	}

	var existing string
	err := a.Services.DB.QueryRowContext(ctx, `SELECT slug FROM templates WHERE slug = $1`, slug).Scan(&existing)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: unknown template %q", errValidation, slug)
		}
		return err
	}
	return nil
}

func (a *API) listLandingSections(ctx context.Context, siteID string) ([]landingSectionResp, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			id::text,
			site_id::text,
			section_key,
			COALESCE(title, ''),
			COALESCE(subtitle, ''),
			COALESCE(content_json::text, '{}'),
			display_order,
			is_enabled
		FROM landing_sections
		WHERE site_id = $1
		ORDER BY display_order ASC, created_at ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections := make([]landingSectionResp, 0)
	for rows.Next() {
		var section landingSectionResp
		if err := rows.Scan(&section.ID, &section.SiteID, &section.SectionKey, &section.Title, &section.Subtitle, &section.ContentJSON, &section.DisplayOrder, &section.IsEnabled); err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, rows.Err()
}

func (a *API) updateLandingSection(ctx context.Context, siteID, sectionID string, payload landingUpdateRequest) (landingSectionResp, error) {
	setters := []string{}
	args := []any{siteID, sectionID}
	if payload.IsEnabled != nil {
		setters = append(setters, fmt.Sprintf("is_enabled = $%d", len(args)+1))
		args = append(args, *payload.IsEnabled)
	}
	if payload.DisplayOrder != nil {
		setters = append(setters, fmt.Sprintf("display_order = $%d", len(args)+1))
		args = append(args, *payload.DisplayOrder)
	}
	if len(setters) == 0 {
		return landingSectionResp{}, fmt.Errorf("%w: no landing section fields provided", errValidation)
	}
	args = append(args, time.Now())
	query := fmt.Sprintf(`
		UPDATE landing_sections
		SET %s, updated_at = $%d
		WHERE site_id = $1 AND id = $2
		RETURNING id::text, site_id::text, section_key, COALESCE(title, ''), COALESCE(subtitle, ''), COALESCE(content_json::text, '{}'), display_order, is_enabled
	`, strings.Join(setters, ", "), len(args))

	row := a.Services.DB.QueryRowContext(ctx, query, args...)
	var section landingSectionResp
	if err := row.Scan(&section.ID, &section.SiteID, &section.SectionKey, &section.Title, &section.Subtitle, &section.ContentJSON, &section.DisplayOrder, &section.IsEnabled); err != nil {
		return landingSectionResp{}, err
	}
	return section, nil
}

func (a *API) reorderLandingSections(ctx context.Context, siteID string, sectionIDs []string) ([]landingSectionResp, error) {
	sections, err := a.listLandingSections(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if len(sectionIDs) != len(sections) {
		return nil, fmt.Errorf("%w: section ordering is incomplete", errValidation)
	}

	sectionMap := make(map[string]landingSectionResp, len(sections))
	for _, section := range sections {
		sectionMap[section.ID] = section
	}

	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for index, sectionID := range sectionIDs {
		if _, ok := sectionMap[sectionID]; !ok {
			return nil, fmt.Errorf("%w: unknown landing section %s", errValidation, sectionID)
		}
		if _, err = tx.ExecContext(ctx, `
			UPDATE landing_sections
			SET display_order = $3, updated_at = NOW()
			WHERE site_id = $1 AND id = $2
		`, siteID, sectionID, index); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return a.listLandingSections(ctx, siteID)
}

func (a *API) listArticles(ctx context.Context, siteID string) ([]articleResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			a.id::text,
			a.site_id::text,
			COALESCE(a.author_id::text, ''),
			COALESCE(a.category_id::text, ''),
			a.title,
			a.slug,
			COALESCE(a.excerpt, ''),
			a.content_markdown,
			COALESCE(a.cover_image_url, ''),
			a.status,
			a.is_featured,
			a.published_at,
			COALESCE(a.seo_title, ''),
			COALESCE(a.seo_description, ''),
			COALESCE(a.canonical_url, ''),
			a.generated_by_ai,
			a.human_reviewed,
			COALESCE(a.ai_prompt, ''),
			COALESCE(a.ai_model, ''),
			a.tags,
			a.updated_at
		FROM articles a
		WHERE a.site_id = $1
		ORDER BY a.updated_at DESC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]articleResponse, 0)
	for rows.Next() {
		var article articleResponse
		var publishedAt sql.NullTime
		var updatedAt time.Time
		if err := rows.Scan(&article.ID, &article.SiteID, &article.AuthorID, &article.CategoryID, &article.Title, &article.Slug, &article.Excerpt, &article.ContentMarkdown, &article.CoverImageURL, &article.Status, &article.IsFeatured, &publishedAt, &article.SEOTitle, &article.SEODescription, &article.CanonicalURL, &article.GeneratedByAI, &article.HumanReviewed, &article.AIPrompt, &article.AIModel, &article.Tags, &updatedAt); err != nil {
			return nil, err
		}
		if publishedAt.Valid {
			value := publishedAt.Time.UTC().Format(time.RFC3339)
			article.PublishedAt = &value
		}
		article.UpdatedAt = updatedAt.Format(time.RFC3339)
		articles = append(articles, article)
	}

	return articles, rows.Err()
}

func (a *API) getArticle(ctx context.Context, articleID string) (articleResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			a.id::text,
			a.site_id::text,
			COALESCE(a.author_id::text, ''),
			COALESCE(a.category_id::text, ''),
			a.title,
			a.slug,
			COALESCE(a.excerpt, ''),
			a.content_markdown,
			COALESCE(a.cover_image_url, ''),
			a.status,
			a.is_featured,
			a.published_at,
			COALESCE(a.seo_title, ''),
			COALESCE(a.seo_description, ''),
			COALESCE(a.canonical_url, ''),
			a.generated_by_ai,
			a.human_reviewed,
			COALESCE(a.ai_prompt, ''),
			COALESCE(a.ai_model, ''),
			a.tags,
			a.updated_at
		FROM articles a
		WHERE a.id = $1
		LIMIT 1
	`, articleID)
	if err != nil {
		return articleResponse{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return articleResponse{}, sql.ErrNoRows
	}

	var article articleResponse
	var publishedAt sql.NullTime
	var updatedAt time.Time
	if err := rows.Scan(&article.ID, &article.SiteID, &article.AuthorID, &article.CategoryID, &article.Title, &article.Slug, &article.Excerpt, &article.ContentMarkdown, &article.CoverImageURL, &article.Status, &article.IsFeatured, &publishedAt, &article.SEOTitle, &article.SEODescription, &article.CanonicalURL, &article.GeneratedByAI, &article.HumanReviewed, &article.AIPrompt, &article.AIModel, &article.Tags, &updatedAt); err != nil {
		return articleResponse{}, err
	}
	if publishedAt.Valid {
		value := publishedAt.Time.UTC().Format(time.RFC3339)
		article.PublishedAt = &value
	}
	article.UpdatedAt = updatedAt.Format(time.RFC3339)
	return article, nil
}

func (a *API) upsertArticleWithSite(ctx context.Context, siteID string, payload articleUpsertRequest) (articleResponse, error) {
	payload.ID = strings.TrimSpace(payload.ID)
	return a.upsertArticle(ctx, payload, siteID)
}

func (a *API) upsertArticle(ctx context.Context, payload articleUpsertRequest, siteIDOverride ...string) (articleResponse, error) {
	siteID := ""
	if len(siteIDOverride) > 0 {
		siteID = strings.TrimSpace(siteIDOverride[0])
	}
	if siteID == "" {
		if payload.ID == "" {
			return articleResponse{}, fmt.Errorf("%w: site id is required", errValidation)
		}
		row := a.Services.DB.QueryRowContext(ctx, `SELECT site_id::text FROM articles WHERE id = $1`, payload.ID)
		if err := row.Scan(&siteID); err != nil {
			return articleResponse{}, err
		}
	}
	if err := validateArticlePayload(payload); err != nil {
		return articleResponse{}, err
	}

	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return articleResponse{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	tags := normalizeTagsInput(payload.Tags)

	articleID := strings.TrimSpace(payload.ID)
	if articleID == "" {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO articles (
				site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, tags, generated_by_ai, human_reviewed, ai_prompt, ai_model, updated_at
			)
			VALUES (
				$1, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, $4, $5, $6, $7, NULLIF($8, ''), $9, $10, CASE WHEN $9 = 'published' THEN NOW() ELSE NULL END, $11, $12, NULLIF($13, ''), $14, FALSE, CASE WHEN $9 IN ('review', 'published') THEN TRUE ELSE FALSE END, '', '', NOW()
			)
			RETURNING id::text
		`, siteID, payload.AuthorID, payload.CategoryID, payload.Title, payload.Slug, payload.Excerpt, payload.ContentMarkdown, payload.CoverImageURL, fallbackString(payload.Status, "draft"), payload.IsFeatured, payload.SEOTitle, payload.SEODescription, payload.CanonicalURL, tags).Scan(&articleID)
		if err != nil {
			return articleResponse{}, err
		}
	} else {
		result, err := tx.ExecContext(ctx, `
			UPDATE articles
			SET
				author_id = NULLIF($2, '')::uuid,
				category_id = NULLIF($3, '')::uuid,
				title = $4,
				slug = $5,
				excerpt = $6,
				content_markdown = $7,
				cover_image_url = NULLIF($8, ''),
				status = $9,
				is_featured = $10,
				published_at = CASE
					WHEN $9 = 'published' AND published_at IS NULL THEN NOW()
					WHEN $9 = 'published' THEN published_at
					ELSE published_at
				END,
				seo_title = $11,
				seo_description = $12,
				canonical_url = NULLIF($13, ''),
				tags = $14,
				human_reviewed = CASE WHEN $9 IN ('review', 'published') THEN TRUE ELSE human_reviewed END,
				updated_at = NOW()
			WHERE id = $1 AND site_id = $15
		`, payload.ID, payload.AuthorID, payload.CategoryID, payload.Title, payload.Slug, payload.Excerpt, payload.ContentMarkdown, payload.CoverImageURL, fallbackString(payload.Status, "draft"), payload.IsFeatured, payload.SEOTitle, payload.SEODescription, payload.CanonicalURL, tags, siteID)
		if err != nil {
			return articleResponse{}, err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return articleResponse{}, sql.ErrNoRows
		}
		articleID = payload.ID
	}

	if err = tx.Commit(); err != nil {
		return articleResponse{}, err
	}

	return a.getArticle(ctx, articleID)
}

func (a *API) deleteArticle(ctx context.Context, articleID string) error {
	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingID string
	if err = tx.QueryRowContext(ctx, `
		SELECT id::text
		FROM articles
		WHERE id = $1
	`, articleID).Scan(&existingID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, articleID); err != nil {
		return err
	}

	return tx.Commit()
}

func validateArticlePayload(payload articleUpsertRequest) error {
	if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Slug) == "" || strings.TrimSpace(payload.ContentMarkdown) == "" {
		return fmt.Errorf("%w: title, slug and content are required", errValidation)
	}
	return nil
}

// normalizeTagsInput trims, dedupes (case-insensitively), and rejoins a
// comma-separated tag string, capping both the count and length of entries
// so free-text tags (including AI-suggested ones) can't grow unbounded.
func normalizeTagsInput(raw string) string {
	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag == "" || len(tag) > 80 {
			continue
		}
		key := strings.ToLower(tag)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		tags = append(tags, tag)
		if len(tags) == 24 {
			break
		}
	}
	return strings.Join(tags, ", ")
}

func (a *API) listAuthors(ctx context.Context, siteID string) ([]authorResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT id::text, site_id::text, name, slug, COALESCE(bio, '')
		FROM authors
		WHERE site_id = $1
		ORDER BY name ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]authorResponse, 0)
	for rows.Next() {
		var item authorResponse
		if err := rows.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Bio); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) listCategories(ctx context.Context, siteID string) ([]categoryResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT id::text, site_id::text, name, slug, COALESCE(description, '')
		FROM categories
		WHERE site_id = $1
		ORDER BY name ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]categoryResponse, 0)
	for rows.Next() {
		var item categoryResponse
		if err := rows.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) getCategory(ctx context.Context, siteID, categoryID string) (categoryResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, site_id::text, name, slug, COALESCE(description, '')
		FROM categories
		WHERE site_id = $1 AND id = $2
	`, siteID, categoryID)

	var item categoryResponse
	if err := row.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Description); err != nil {
		return categoryResponse{}, err
	}
	return item, nil
}

func (a *API) createCategory(ctx context.Context, siteID string, payload categoryUpsertRequest) (categoryResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return categoryResponse{}, fmt.Errorf("%w: category name is required", errValidation)
	}

	slug, err := a.uniqueCategorySlug(ctx, siteID, "", name)
	if err != nil {
		return categoryResponse{}, err
	}

	var item categoryResponse
	err = a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO categories (site_id, name, slug, description)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id::text, site_id::text, name, slug, COALESCE(description, '')
	`, siteID, name, slug, strings.TrimSpace(payload.Description)).Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Description)
	if err != nil {
		if isUniqueViolation(err) {
			return categoryResponse{}, fmt.Errorf("%w: category slug already exists", errValidation)
		}
		return categoryResponse{}, err
	}
	return item, nil
}

func (a *API) updateCategory(ctx context.Context, siteID, categoryID string, payload categoryUpsertRequest) (categoryResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return categoryResponse{}, fmt.Errorf("%w: category name is required", errValidation)
	}

	slug, err := a.uniqueCategorySlug(ctx, siteID, categoryID, name)
	if err != nil {
		return categoryResponse{}, err
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE categories
		SET name = $3, slug = $4, description = NULLIF($5, ''), updated_at = NOW()
		WHERE id = $1 AND site_id = $2
	`, categoryID, siteID, name, slug, strings.TrimSpace(payload.Description))
	if err != nil {
		if isUniqueViolation(err) {
			return categoryResponse{}, fmt.Errorf("%w: category slug already exists", errValidation)
		}
		return categoryResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return categoryResponse{}, sql.ErrNoRows
	}
	return a.getCategory(ctx, siteID, categoryID)
}

func (a *API) deleteCategory(ctx context.Context, siteID, categoryID string) error {
	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingID string
	if err = tx.QueryRowContext(ctx, `
		SELECT id::text
		FROM categories
		WHERE id = $1 AND site_id = $2
	`, categoryID, siteID).Scan(&existingID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE articles
		SET category_id = NULL, updated_at = NOW()
		WHERE site_id = $1 AND category_id = $2
	`, siteID, categoryID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		DELETE FROM categories
		WHERE id = $1 AND site_id = $2
	`, categoryID, siteID); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *API) uniqueCategorySlug(ctx context.Context, siteID, categoryID, name string) (string, error) {
	return a.uniqueSlug(ctx, "categories", siteID, categoryID, name)
}

func (a *API) uniqueAuthorSlug(ctx context.Context, siteID, authorID, name string) (string, error) {
	return a.uniqueSlug(ctx, "authors", siteID, authorID, name)
}

func (a *API) uniqueSlug(ctx context.Context, table, siteID, excludeID, name string) (string, error) {
	base := slugify(name)
	if base == "" {
		return "", fmt.Errorf("%w: name must contain letters or numbers", errValidation)
	}

	query := fmt.Sprintf(`SELECT slug FROM %s WHERE site_id = $1`, table)
	args := []any{siteID}
	if strings.TrimSpace(excludeID) != "" {
		query += " AND id <> $2"
		args = append(args, excludeID)
	}

	rows, err := a.Services.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	used := make(map[string]struct{})
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return "", err
		}
		used[slug] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	candidate := base
	if _, exists := used[candidate]; !exists {
		return candidate, nil
	}

	for suffix := 2; ; suffix++ {
		candidate = fmt.Sprintf("%s-%d", base, suffix)
		if _, exists := used[candidate]; !exists {
			return candidate, nil
		}
	}
}

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = slugPattern.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func (a *API) listMediaAssets(ctx context.Context, siteID string) ([]mediaAssetResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			id::text,
			site_id::text,
			file_name,
			file_url,
			COALESCE(mime_type, ''),
			COALESCE(size_bytes, 0),
			storage_provider,
			COALESCE(storage_key, ''),
			COALESCE(alt_text, ''),
			created_at
		FROM media_assets
		WHERE site_id = $1
		ORDER BY created_at DESC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]mediaAssetResponse, 0)
	for rows.Next() {
		var item mediaAssetResponse
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.SiteID, &item.FileName, &item.FileURL, &item.MimeType, &item.SizeBytes, &item.StorageProvider, &item.StorageKey, &item.AltText, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) getAuthor(ctx context.Context, siteID, authorID string) (authorResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, site_id::text, name, slug, COALESCE(bio, '')
		FROM authors
		WHERE site_id = $1 AND id = $2
	`, siteID, authorID)

	var item authorResponse
	if err := row.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Bio); err != nil {
		return authorResponse{}, err
	}
	return item, nil
}

func (a *API) createAuthor(ctx context.Context, siteID string, payload authorUpsertRequest) (authorResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return authorResponse{}, fmt.Errorf("%w: author name is required", errValidation)
	}

	slug, err := a.uniqueAuthorSlug(ctx, siteID, "", name)
	if err != nil {
		return authorResponse{}, err
	}

	var item authorResponse
	err = a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO authors (site_id, name, slug, bio)
		VALUES ($1, $2, $3, NULLIF($4, ''))
		RETURNING id::text, site_id::text, name, slug, COALESCE(bio, '')
	`, siteID, name, slug, strings.TrimSpace(payload.Bio)).Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug, &item.Bio)
	if err != nil {
		if isUniqueViolation(err) {
			return authorResponse{}, fmt.Errorf("%w: author slug already exists", errValidation)
		}
		return authorResponse{}, err
	}
	return item, nil
}

func (a *API) updateAuthor(ctx context.Context, siteID, authorID string, payload authorUpsertRequest) (authorResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return authorResponse{}, fmt.Errorf("%w: author name is required", errValidation)
	}

	slug, err := a.uniqueAuthorSlug(ctx, siteID, authorID, name)
	if err != nil {
		return authorResponse{}, err
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE authors
		SET name = $3, slug = $4, bio = NULLIF($5, ''), updated_at = NOW()
		WHERE id = $1 AND site_id = $2
	`, authorID, siteID, name, slug, strings.TrimSpace(payload.Bio))
	if err != nil {
		if isUniqueViolation(err) {
			return authorResponse{}, fmt.Errorf("%w: author slug already exists", errValidation)
		}
		return authorResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return authorResponse{}, sql.ErrNoRows
	}
	return a.getAuthor(ctx, siteID, authorID)
}

func (a *API) deleteAuthor(ctx context.Context, siteID, authorID string) error {
	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingID string
	if err = tx.QueryRowContext(ctx, `
		SELECT id::text
		FROM authors
		WHERE id = $1 AND site_id = $2
	`, authorID, siteID).Scan(&existingID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE articles
		SET author_id = NULL, updated_at = NOW()
		WHERE site_id = $1 AND author_id = $2
	`, siteID, authorID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		DELETE FROM authors
		WHERE id = $1 AND site_id = $2
	`, authorID, siteID); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *API) createMediaAsset(ctx context.Context, siteID string, payload mediaUpsertRequest) (mediaAssetResponse, error) {
	if strings.TrimSpace(payload.FileName) == "" || strings.TrimSpace(payload.FileURL) == "" {
		return mediaAssetResponse{}, fmt.Errorf("%w: file name and URL are required", errValidation)
	}

	return a.persistMediaAsset(ctx, siteID, payload)
}

func (a *API) uploadMediaAsset(ctx context.Context, siteID, fileName string, contents []byte, mimeType, altText string) (mediaAssetResponse, error) {
	if strings.TrimSpace(fileName) == "" {
		return mediaAssetResponse{}, fmt.Errorf("%w: file name is required", errValidation)
	}
	if len(contents) == 0 {
		return mediaAssetResponse{}, fmt.Errorf("%w: file contents are required", errValidation)
	}
	altText, err := validateMediaAltText(altText)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	optimized, err := media.OptimizeForUpload(fileName, mimeType, contents)
	if err != nil {
		return mediaAssetResponse{}, fmt.Errorf("%w: unable to optimize image", errValidation)
	}

	stored, err := a.Services.Storage.Upload(ctx, storage.UploadFile{
		FileName: optimized.FileName,
		Contents: optimized.Contents,
		MimeType: optimized.MimeType,
		SiteID:   siteID,
	})
	if err != nil {
		return mediaAssetResponse{}, err
	}
	if stored == nil || strings.TrimSpace(stored.PublicURL) == "" {
		return mediaAssetResponse{}, fmt.Errorf("%w: storage did not return a public URL", errValidation)
	}

	return a.persistMediaAsset(ctx, siteID, mediaUpsertRequest{
		FileName:        optimized.FileName,
		FileURL:         stored.PublicURL,
		MimeType:        optimized.MimeType,
		SizeBytes:       int64(len(optimized.Contents)),
		StorageProvider: a.storageProviderName(),
		StorageKey:      stored.Key,
		AltText:         altText,
	})
}

func (a *API) persistMediaAsset(ctx context.Context, siteID string, payload mediaUpsertRequest) (mediaAssetResponse, error) {
	var item mediaAssetResponse
	var createdAt time.Time
	err := a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO media_assets (
			site_id, file_name, file_url, mime_type, size_bytes, storage_provider, storage_key, alt_text
		)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, 0), $6, NULLIF($7, ''), NULLIF($8, ''))
		RETURNING id::text, site_id::text, file_name, file_url, COALESCE(mime_type, ''), COALESCE(size_bytes, 0), storage_provider, COALESCE(storage_key, ''), COALESCE(alt_text, ''), created_at
	`, siteID, payload.FileName, payload.FileURL, payload.MimeType, payload.SizeBytes, fallbackString(payload.StorageProvider, "s3"), payload.StorageKey, payload.AltText).Scan(&item.ID, &item.SiteID, &item.FileName, &item.FileURL, &item.MimeType, &item.SizeBytes, &item.StorageProvider, &item.StorageKey, &item.AltText, &createdAt)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return item, nil
}

func (a *API) updateMediaAsset(ctx context.Context, siteID, assetID string, payload mediaUpdateRequest) (mediaAssetResponse, error) {
	altText, err := validateMediaAltText(payload.AltText)
	if err != nil {
		return mediaAssetResponse{}, err
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE media_assets
		SET alt_text = NULLIF($3, '')
		WHERE id = $1 AND site_id = $2
	`, assetID, siteID, altText)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return mediaAssetResponse{}, sql.ErrNoRows
	}
	return a.getMediaAsset(ctx, siteID, assetID)
}

func (a *API) replaceMediaAsset(ctx context.Context, siteID, assetID, fileName string, contents []byte, mimeType, altText string) (mediaAssetResponse, error) {
	altText, err := validateMediaAltText(altText)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	current, err := a.getMediaAsset(ctx, siteID, assetID)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	optimized, err := media.OptimizeForUpload(fileName, mimeType, contents)
	if err != nil {
		return mediaAssetResponse{}, fmt.Errorf("%w: unable to optimize image", errValidation)
	}

	stored, err := a.Services.Storage.Upload(ctx, storage.UploadFile{
		FileName:  optimized.FileName,
		Contents:  optimized.Contents,
		MimeType:  optimized.MimeType,
		SiteID:    siteID,
		ObjectKey: current.StorageKey,
	})
	if err != nil {
		return mediaAssetResponse{}, err
	}
	if stored == nil || strings.TrimSpace(stored.PublicURL) == "" {
		return mediaAssetResponse{}, fmt.Errorf("%w: storage did not return a public URL", errValidation)
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE media_assets
		SET file_name = $3, file_url = $4, mime_type = $5, size_bytes = $6, storage_provider = $7, storage_key = $8, alt_text = NULLIF($9, '')
		WHERE id = $1 AND site_id = $2
	`, assetID, siteID, optimized.FileName, stored.PublicURL, optimized.MimeType, len(optimized.Contents), a.storageProviderName(), stored.Key, altText)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return mediaAssetResponse{}, sql.ErrNoRows
	}
	return a.getMediaAsset(ctx, siteID, assetID)
}

func (a *API) deleteMediaAsset(ctx context.Context, siteID, assetID string) error {
	asset, err := a.getMediaAsset(ctx, siteID, assetID)
	if err != nil {
		return err
	}

	var inUse bool
	if err := a.Services.DB.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM articles WHERE site_id = $1 AND cover_image_url = $2
			UNION ALL
			SELECT 1 FROM sites WHERE id = $1 AND (logo_media_id = $3::uuid OR favicon_media_id = $3::uuid)
		)
	`, siteID, asset.FileURL, assetID).Scan(&inUse); err != nil {
		return err
	}
	if inUse {
		return fmt.Errorf("%w: remove this asset from articles or site branding before deleting it", errConflict)
	}
	if asset.StorageKey != "" {
		if err := a.Services.Storage.Delete(ctx, asset.StorageKey); err != nil {
			return err
		}
	}

	result, err := a.Services.DB.ExecContext(ctx, `DELETE FROM media_assets WHERE id = $1 AND site_id = $2`, assetID, siteID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (a *API) getMediaAsset(ctx context.Context, siteID, assetID string) (mediaAssetResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, site_id::text, file_name, file_url, COALESCE(mime_type, ''), COALESCE(size_bytes, 0), storage_provider, COALESCE(storage_key, ''), COALESCE(alt_text, ''), created_at
		FROM media_assets
		WHERE id = $1 AND site_id = $2
	`, assetID, siteID)
	var item mediaAssetResponse
	var createdAt time.Time
	if err := row.Scan(&item.ID, &item.SiteID, &item.FileName, &item.FileURL, &item.MimeType, &item.SizeBytes, &item.StorageProvider, &item.StorageKey, &item.AltText, &createdAt); err != nil {
		return mediaAssetResponse{}, err
	}
	item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return item, nil
}

func validateMediaAltText(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len([]rune(value)) > maxMediaAltTextRunes {
		return "", fmt.Errorf("%w: alt text must not exceed %d characters", errValidation, maxMediaAltTextRunes)
	}
	return value, nil
}

func (a *API) storageProviderName() string {
	endpoint := strings.ToLower(strings.TrimSpace(a.Config.S3Endpoint))
	publicURL := strings.ToLower(strings.TrimSpace(a.Config.S3PublicURL))
	switch {
	case strings.Contains(endpoint, "minio") || strings.Contains(publicURL, "localhost:9002"):
		return "minio"
	case strings.Contains(endpoint, "r2") || strings.Contains(publicURL, "r2"):
		return "r2"
	case endpoint != "":
		return "s3"
	default:
		return "s3"
	}
}

func (a *API) listBuilds(ctx context.Context, siteID string) ([]buildResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT
			id::text,
			site_id::text,
			status,
			build_type,
			COALESCE(logs, ''),
			COALESCE(output_path, ''),
			COALESCE(deploy_provider, ''),
			COALESCE(deploy_status, ''),
			COALESCE(deploy_url, ''),
			COALESCE(deploy_revision, ''),
			started_at,
			finished_at
		FROM builds
		WHERE site_id = $1
		ORDER BY created_at DESC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]buildResponse, 0)
	for rows.Next() {
		var item buildResponse
		var startedAt sql.NullTime
		var finishedAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.SiteID, &item.Status, &item.BuildType, &item.Logs, &item.OutputPath, &item.DeployProvider, &item.DeployStatus, &item.DeployURL, &item.DeployRevision, &startedAt, &finishedAt); err != nil {
			return nil, err
		}
		if startedAt.Valid {
			value := startedAt.Time.UTC().Format(time.RFC3339)
			item.StartedAt = &value
		}
		if finishedAt.Valid {
			value := finishedAt.Time.UTC().Format(time.RFC3339)
			item.FinishedAt = &value
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (a *API) clearBuildHistory(ctx context.Context, siteID string) error {
	var exists bool
	if err := a.Services.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sites WHERE id = $1)`, siteID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	_, err := a.Services.DB.ExecContext(ctx, `DELETE FROM builds WHERE site_id = $1`, siteID)
	return err
}

func (a *API) getBuild(ctx context.Context, buildID string) (buildResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT
			id::text,
			site_id::text,
			status,
			build_type,
			COALESCE(logs, ''),
			COALESCE(output_path, ''),
			COALESCE(deploy_provider, ''),
			COALESCE(deploy_status, ''),
			COALESCE(deploy_url, ''),
			COALESCE(deploy_revision, ''),
			started_at,
			finished_at
		FROM builds
		WHERE id = $1
	`, buildID)

	var item buildResponse
	var startedAt sql.NullTime
	var finishedAt sql.NullTime
	if err := row.Scan(&item.ID, &item.SiteID, &item.Status, &item.BuildType, &item.Logs, &item.OutputPath, &item.DeployProvider, &item.DeployStatus, &item.DeployURL, &item.DeployRevision, &startedAt, &finishedAt); err != nil {
		return buildResponse{}, err
	}
	if startedAt.Valid {
		value := startedAt.Time.UTC().Format(time.RFC3339)
		item.StartedAt = &value
	}
	if finishedAt.Valid {
		value := finishedAt.Time.UTC().Format(time.RFC3339)
		item.FinishedAt = &value
	}
	return item, nil
}

func (a *API) createBuild(ctx context.Context, siteID string, payload buildCreateRequest) (buildResponse, error) {
	buildType := fallbackString(payload.BuildType, "published")
	if buildType != "preview" && buildType != "published" {
		return buildResponse{}, fmt.Errorf("%w: invalid build type", errValidation)
	}

	site, err := a.getSiteModel(ctx, siteID)
	if err != nil {
		return buildResponse{}, err
	}
	buildSite := site
	provider := fallbackString(site.DeployProvider, "none")
	if buildType == "preview" {
		provider = fallbackString(site.PreviewDeployProvider, provider)
	}

	var buildID string
	if err := a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO builds (site_id, status, build_type, logs, deploy_provider, deploy_status, started_at)
		VALUES ($1, 'running', $2, 'Build started.', $3, 'pending', NOW())
		RETURNING id::text
	`, siteID, buildType, provider).Scan(&buildID); err != nil {
		return buildResponse{}, err
	}
	fail := func(cause error) (buildResponse, error) {
		message := strings.TrimSpace(cause.Error())
		failureCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = a.Services.DB.ExecContext(failureCtx, `
			UPDATE builds SET status = 'failed', deploy_status = 'failed', logs = $2, finished_at = NOW() WHERE id = $1
		`, buildID, message)
		return buildResponse{}, cause
	}

	if buildType == "published" && len(payload.ArticleIDs) > 0 {
		if err := a.publishArticles(ctx, siteID, payload.ArticleIDs); err != nil {
			return fail(err)
		}
	}

	content, err := a.siteBuildContent(ctx, siteID, buildType == "published")
	if err != nil {
		return fail(err)
	}

	if buildType == "published" {
		if err := a.rewritePublishedCoverImages(ctx, &content, buildSite); err != nil {
			return fail(err)
		}
	}

	outputPath, err := a.Services.Builder.GenerateSite(ctx, content, builder.GenerateOptions{
		SiteID:     siteID,
		Preview:    buildType == "preview",
		ArticleIDs: payload.ArticleIDs,
	})
	if err != nil {
		return fail(err)
	}

	build := models.Build{
		ID:         buildID,
		SiteID:     siteID,
		Status:     "success",
		BuildType:  buildType,
		Logs:       buildLogs(buildType, buildArticleCount(buildType, len(payload.ArticleIDs), len(content.Articles))),
		OutputPath: outputPath,
	}
	deployResult, err := a.Services.Deploy.Deploy(ctx, buildSite, build, outputPath)
	if err != nil {
		return fail(err)
	}

	deployProvider := fallbackString(buildSite.DeployProvider, "none")
	deployStatus := "deployed"
	deployURL := ""
	deployRevision := ""
	if deployResult != nil {
		deployProvider = fallbackString(deployResult.Provider, deployProvider)
		deployStatus = "deployed"
		deployURL = strings.TrimSpace(deployResult.URL)
		deployRevision = strings.TrimSpace(deployResult.Revision)
		if message := strings.TrimSpace(deployResult.Message); message != "" {
			build.Logs = strings.TrimSpace(build.Logs + " " + message)
		}
	}

	err = func() error {
		_, updateErr := a.Services.DB.ExecContext(ctx, `
			UPDATE builds SET status = 'success', logs = $2, output_path = $3, deploy_provider = $4,
				deploy_status = $5, deploy_url = $6, deploy_revision = $7, finished_at = NOW()
			WHERE id = $1
		`, buildID, build.Logs, outputPath, deployProvider, deployStatus, deployURL, deployRevision)
		return updateErr
	}()
	if err != nil {
		return fail(err)
	}
	return a.getBuild(ctx, buildID)
}

// rewritePublishedCoverImages re-uploads any article cover image still
// pointing at local/dev storage to the site's configured production S3
// bucket and rewrites the URL in place, so published output never links
// back to a dev-only host. It is a no-op when the site has no production
// storage configured.
func (a *API) rewritePublishedCoverImages(ctx context.Context, content *builder.SiteContent, site models.Site) error {
	rawConfig, err := json.Marshal(site.StorageConfig)
	if err != nil {
		return err
	}
	prodConfig, err := storage.ParseConfig(string(rawConfig))
	if err != nil || !prodConfig.IsConfigured() {
		return nil
	}

	prodStorage, err := storage.NewFromSiteConfig(prodConfig, os.LookupEnv)
	if err != nil {
		return fmt.Errorf("production storage is misconfigured: %w", err)
	}
	return a.rewritePublishedCoverImagesWithStorage(ctx, content, site, prodStorage)
}

func (a *API) rewritePublishedCoverImagesWithStorage(ctx context.Context, content *builder.SiteContent, site models.Site, prodStorage storage.StorageProvider) error {
	blogPath, err := models.CanonicalBlogPath(site.BlogPath)
	if err != nil {
		return fmt.Errorf("production storage requires a valid blog path: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	for i := range content.Articles {
		article := &content.Articles[i]
		coverURL := strings.TrimSpace(article.CoverImageURL)
		if coverURL == "" || !isDevStorageImageURL(coverURL, a.Config.S3PublicURL) {
			continue
		}
		sourceURL, err := devStorageObjectURL(coverURL, a.Config.S3PublicURL, a.Config.S3Endpoint, a.Config.S3Bucket)
		if err != nil {
			return fmt.Errorf("unable to resolve dev cover image source for %q: %w", article.Title, err)
		}
		objectKey, err := productionCoverObjectKey(blogPath, coverURL)
		if err != nil {
			return fmt.Errorf("unable to resolve production cover image key for %q: %w", article.Title, err)
		}

		newURL, err := migrateImageToProductionStorage(ctx, client, prodStorage, site.ID, sourceURL, objectKey)
		if err != nil {
			return fmt.Errorf("unable to migrate cover image for %q: %w", article.Title, err)
		}
		if err := a.persistProductionCoverImageURL(ctx, site.ID, article.ID, newURL); err != nil {
			return fmt.Errorf("unable to save production cover image URL for %q: %w", article.Title, err)
		}
		article.CoverImageURL = newURL
	}
	return nil
}

func (a *API) persistProductionCoverImageURL(ctx context.Context, siteID, articleID, coverImageURL string) error {
	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE articles
		SET cover_image_url = $1, updated_at = NOW()
		WHERE id = $2 AND site_id = $3
	`, coverImageURL, articleID, siteID)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// productionCoverObjectKey keeps production cover images directly beneath the
// generated blog path, for example blog/cover.png.
func productionCoverObjectKey(blogPath, imageURL string) (string, error) {
	canonicalBlogPath, err := models.CanonicalBlogPath(blogPath)
	if err != nil {
		return "", err
	}
	parsed, err := url.Parse(strings.TrimSpace(imageURL))
	if err != nil {
		return "", errors.New("cover image URL is invalid")
	}
	fileName := path.Base(parsed.Path)
	if fileName == "." || fileName == "/" || fileName == "" {
		return "", errors.New("cover image file name is missing")
	}
	return path.Join(strings.TrimPrefix(canonicalBlogPath, "/"), fileName), nil
}

// devStorageObjectURL maps a browser-facing development storage URL to the
// S3 endpoint reachable by the API server. This prevents a container from
// attempting to reach its own localhost port while publishing cover images.
func devStorageObjectURL(imageURL, devPublicURL, devEndpoint, bucket string) (string, error) {
	image, err := url.Parse(strings.TrimSpace(imageURL))
	if err != nil || image.Scheme == "" || image.Host == "" {
		return "", errors.New("cover image URL is invalid")
	}
	public, err := url.Parse(strings.TrimRight(strings.TrimSpace(devPublicURL), "/"))
	if err != nil || public.Scheme == "" || public.Host == "" {
		return "", errors.New("development storage public URL is invalid")
	}
	endpoint, err := url.Parse(strings.TrimRight(strings.TrimSpace(devEndpoint), "/"))
	if err != nil || (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.Host == "" {
		return "", errors.New("development storage endpoint is invalid")
	}
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return "", errors.New("development storage bucket is not configured")
	}

	publicPath := strings.TrimRight(public.Path, "/")
	if image.Scheme != public.Scheme || !strings.EqualFold(image.Host, public.Host) || !strings.HasPrefix(image.Path, publicPath+"/") {
		return "", errors.New("cover image does not use the configured development storage public URL")
	}
	objectKey := strings.TrimPrefix(image.Path, publicPath+"/")
	if objectKey == "" {
		return "", errors.New("cover image storage key is missing")
	}

	endpoint.Path = "/" + strings.TrimLeft(path.Join(endpoint.Path, bucket, objectKey), "/")
	endpoint.RawPath = ""
	endpoint.RawQuery = ""
	endpoint.Fragment = ""
	return endpoint.String(), nil
}

// isDevStorageImageURL reports whether an image URL points at this server's
// local/dev storage rather than an already-production host.
func isDevStorageImageURL(rawURL, devPublicURL string) bool {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return false
	}
	devPublicURL = strings.TrimRight(strings.TrimSpace(devPublicURL), "/")
	if devPublicURL != "" && strings.HasPrefix(trimmed, devPublicURL+"/") {
		return true
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == "localhost" || host == "127.0.0.1"
}

func migrateImageToProductionStorage(ctx context.Context, client *http.Client, prodStorage storage.StorageProvider, siteID, imageURL, objectKey string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch image: unexpected status %d", resp.StatusCode)
	}

	contents, err := io.ReadAll(io.LimitReader(resp.Body, 25<<20+1))
	if err != nil {
		return "", err
	}
	if len(contents) > 25<<20 {
		return "", errors.New("image exceeds the maximum size for migration")
	}

	mimeType := http.DetectContentType(contents)
	fileName := "cover"
	if parsed, parseErr := url.Parse(imageURL); parseErr == nil {
		if base := filepath.Base(parsed.Path); base != "." && base != "/" {
			fileName = base
		}
	}
	optimized, err := media.OptimizeForUpload(fileName, mimeType, contents)
	if err != nil {
		return "", err
	}
	objectKey = path.Join(path.Dir(objectKey), optimized.FileName)

	stored, err := prodStorage.Upload(ctx, storage.UploadFile{
		FileName:  optimized.FileName,
		Contents:  optimized.Contents,
		MimeType:  optimized.MimeType,
		SiteID:    siteID,
		ObjectKey: objectKey,
	})
	if err != nil {
		return "", err
	}
	return stored.PublicURL, nil
}

func buildLogs(buildType string, articleCount int) string {
	if buildType == "preview" {
		return fmt.Sprintf("Preview build completed successfully for %d articles.", articleCount)
	}
	return fmt.Sprintf("Published build completed successfully for %d published articles.", articleCount)
}

func buildArticleCount(buildType string, selectedCount, contentCount int) int {
	if buildType == "preview" && selectedCount > 0 {
		return selectedCount
	}
	return contentCount
}

func (a *API) publishArticles(ctx context.Context, siteID string, articleIDs []string) error {
	if len(articleIDs) == 0 {
		return nil
	}

	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, articleID := range articleIDs {
		articleID = strings.TrimSpace(articleID)
		if articleID == "" {
			continue
		}

		result, err := tx.ExecContext(ctx, `
			UPDATE articles
			SET
				status = 'published',
				published_at = COALESCE(published_at, NOW()),
				human_reviewed = TRUE,
				updated_at = NOW()
			WHERE id = $1 AND site_id = $2
		`, articleID, siteID)
		if err != nil {
			return err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return sql.ErrNoRows
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *API) siteBuildContent(ctx context.Context, siteID string, publishedOnly bool) (builder.SiteContent, error) {
	site, err := a.getSiteModel(ctx, siteID)
	if err != nil {
		return builder.SiteContent{}, err
	}

	articles, err := a.listArticles(ctx, siteID)
	if err != nil {
		return builder.SiteContent{}, err
	}
	if publishedOnly {
		filtered := make([]articleResponse, 0, len(articles))
		for _, article := range articles {
			if article.Status == "published" {
				filtered = append(filtered, article)
			}
		}
		articles = filtered
	}

	authors, err := a.listAuthors(ctx, siteID)
	if err != nil {
		return builder.SiteContent{}, err
	}
	categories, err := a.listCategories(ctx, siteID)
	if err != nil {
		return builder.SiteContent{}, err
	}
	sections, err := a.listLandingSections(ctx, siteID)
	if err != nil {
		return builder.SiteContent{}, err
	}

	authorNames := make(map[string]authorResponse, len(authors))
	for _, author := range authors {
		authorNames[author.ID] = author
	}
	categoryNames := make(map[string]categoryResponse, len(categories))
	for _, category := range categories {
		categoryNames[category.ID] = category
	}

	content := builder.SiteContent{
		Site:            site,
		LandingSections: make([]builder.LandingSectionContent, 0, len(sections)),
		Articles:        make([]builder.ArticleContent, 0, len(articles)),
	}

	for _, section := range sections {
		content.LandingSections = append(content.LandingSections, builder.LandingSectionContent{
			SectionKey:   section.SectionKey,
			Title:        section.Title,
			Subtitle:     section.Subtitle,
			ContentJSON:  parseJSONMap(section.ContentJSON),
			DisplayOrder: section.DisplayOrder,
			IsEnabled:    section.IsEnabled,
		})
	}

	for _, article := range articles {
		articleContent := builder.ArticleContent{
			ID:              article.ID,
			Title:           article.Title,
			Slug:            article.Slug,
			Excerpt:         article.Excerpt,
			ContentMarkdown: article.ContentMarkdown,
			CoverImageURL:   article.CoverImageURL,
			Status:          article.Status,
			IsFeatured:      article.IsFeatured,
			PublishedAt:     fallbackStringPtr(article.PublishedAt),
			UpdatedAt:       article.UpdatedAt,
			SEOTitle:        article.SEOTitle,
			SEODescription:  article.SEODescription,
			CanonicalURL:    article.CanonicalURL,
		}

		if author, ok := authorNames[article.AuthorID]; ok {
			articleContent.AuthorName = author.Name
			articleContent.AuthorSlug = author.Slug
		}
		if category, ok := categoryNames[article.CategoryID]; ok {
			articleContent.CategoryName = category.Name
			articleContent.CategorySlug = category.Slug
		}
		for _, tagName := range strings.Split(article.Tags, ",") {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}
			articleContent.Tags = append(articleContent.Tags, builder.TagContent{Name: tagName, Slug: slugify(tagName)})
		}
		content.Articles = append(content.Articles, articleContent)
	}

	return content, nil
}

func fallbackStringPtr(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func (a *API) getSiteModel(ctx context.Context, siteID string) (models.Site, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT
			id::text,
			name,
			slug,
			COALESCE(domain, ''),
			blog_path,
			COALESCE(description, ''),
			content_context,
			COALESCE(logo_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.logo_media_id), ''),
			COALESCE(favicon_media_id::text, ''),
			COALESCE((SELECT file_url FROM media_assets WHERE id = sites.favicon_media_id), ''),
			status,
			template_key,
			COALESCE(theme_config::text, '{}'),
			COALESCE(deploy_provider, ''),
			COALESCE(deploy_config::text, '{}'),
			COALESCE(preview_deploy_provider, ''),
			COALESCE(preview_deploy_config::text, '{}'),
			COALESCE(ai_config::text, '{}'),
			COALESCE(storage_config::text, '{}'),
			created_at,
			updated_at
		FROM sites
		WHERE id = $1
	`, siteID)

	var site models.Site
	var themeConfig string
	var deployConfig string
	var previewDeployConfig string
	var aiConfig string
	var storageConfig string
	if err := row.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Description, &site.ContentContext, &site.LogoMediaID, &site.LogoURL, &site.FaviconMediaID, &site.FaviconURL, &site.Status, &site.TemplateKey, &themeConfig, &site.DeployProvider, &deployConfig, &site.PreviewDeployProvider, &previewDeployConfig, &aiConfig, &storageConfig, &site.CreatedAt, &site.UpdatedAt); err != nil {
		return models.Site{}, err
	}
	site.ThemeConfig = parseJSONMap(themeConfig)
	site.DeployConfig = parseJSONMap(deployConfig)
	site.PreviewDeployConfig = parseJSONMap(previewDeployConfig)
	site.AIConfig = parseJSONMap(aiConfig)
	site.StorageConfig = parseJSONMap(storageConfig)
	return site, nil
}

func (a *API) validateBrandingAssets(ctx context.Context, siteID string, mediaIDs ...string) error {
	for _, mediaID := range mediaIDs {
		mediaID = strings.TrimSpace(mediaID)
		if mediaID == "" {
			continue
		}
		var exists bool
		if err := a.Services.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM media_assets WHERE id = $1 AND site_id = $2)`, mediaID, siteID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("%w: branding assets must belong to the selected site", errValidation)
		}
	}
	return nil
}

func parseJSONMap(raw string) map[string]any {
	if strings.TrimSpace(raw) == "" {
		return map[string]any{}
	}

	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return map[string]any{}
	}
	if value == nil {
		return map[string]any{}
	}
	return value
}

func deploymentWarningsForSite(site siteResponse) []string {
	warnings := make([]string, 0, 4)
	warnings = append(warnings, firebaseDeploymentWarnings("production", site.DeployProvider, site.DeployConfig)...)
	warnings = append(warnings, firebaseDeploymentWarnings("preview", site.PreviewDeployProvider, site.PreviewDeployConfig)...)
	warnings = append(warnings, cloudflarePagesDeploymentWarnings("production", site.DeployProvider, site.DeployConfig)...)
	warnings = append(warnings, cloudflarePagesDeploymentWarnings("preview", site.PreviewDeployProvider, site.PreviewDeployConfig)...)
	warnings = append(warnings, repositoryDeploymentWarnings("production", site.DeployProvider, site.DeployConfig)...)
	warnings = append(warnings, repositoryDeploymentWarnings("preview", site.PreviewDeployProvider, site.PreviewDeployConfig)...)
	return warnings
}

func repositoryDeploymentWarnings(channel, provider, rawConfig string) []string {
	if !strings.EqualFold(strings.TrimSpace(provider), "git_repository") {
		return nil
	}
	values := parseJSONMap(rawConfig)
	secretRef := strings.TrimSpace(configValue(values, "tokenSecretRef"))
	if secretRef == "" {
		return []string{fmt.Sprintf("Repository %s deployment requires a token environment variable.", channel)}
	}
	if strings.TrimSpace(os.Getenv(secretRef)) == "" {
		return []string{fmt.Sprintf("Repository %s deploy secret %s is not set on the API server.", channel, secretRef)}
	}
	return nil
}

func cloudflarePagesDeploymentWarnings(channel, provider, rawConfig string) []string {
	if !isCloudflarePagesProvider(provider) {
		return nil
	}

	values := parseJSONMap(rawConfig)
	warnings := make([]string, 0, 2)
	if !cloudflareProjectName.MatchString(strings.TrimSpace(configValue(values, "projectName"))) {
		warnings = append(warnings, fmt.Sprintf("Cloudflare Pages %s deployment config is missing a valid projectName.", channel))
	}
	if strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN")) == "" || strings.TrimSpace(os.Getenv("CLOUDFLARE_ACCOUNT_ID")) == "" {
		warnings = append(warnings, fmt.Sprintf("Cloudflare Pages %s deployment requires CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID on the API server.", channel))
	}
	return warnings
}

func firebaseDeploymentWarnings(channel, provider, rawConfig string) []string {
	values := parseJSONMap(rawConfig)
	resolvedProvider := strings.TrimSpace(provider)
	if resolvedProvider == "" {
		resolvedProvider = strings.TrimSpace(configValue(values, "provider"))
	}
	if resolvedProvider != "firebase" {
		return nil
	}

	secretRef := strings.TrimSpace(configValue(values, "serviceAccountSecretRef"))
	if secretRef == "" {
		secretRef = strings.TrimSpace(configValue(values, "tokenSecretRef"))
	}
	if secretRef == "" {
		return []string{fmt.Sprintf("Firebase %s deployment config is missing serviceAccountSecretRef.", channel)}
	}
	if strings.TrimSpace(os.Getenv(secretRef)) == "" {
		return []string{fmt.Sprintf("Firebase %s deploy secret %s is not set on the API server.", channel, secretRef)}
	}
	return nil
}

func configValue(values map[string]any, key string) string {
	raw, ok := values[key]
	if !ok || raw == nil {
		return ""
	}
	text, ok := raw.(string)
	if !ok {
		return ""
	}
	return text
}

func seedSiteDefaults(ctx context.Context, tx *sql.Tx, siteID string) error {
	var authorID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO authors (site_id, name, slug, bio)
		VALUES ($1, 'Site Editor', 'site-editor', 'Default author seeded for the new site.')
		RETURNING id::text
	`, siteID).Scan(&authorID); err != nil {
		return err
	}

	var categoryID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO categories (site_id, name, slug, description)
		VALUES ($1, 'General', 'general', 'Default category for the new site.')
		RETURNING id::text
	`, siteID).Scan(&categoryID); err != nil {
		return err
	}

	sections := []struct {
		key       string
		title     string
		subtitle  string
		order     int
		isEnabled bool
		content   string
	}{
		{"hero", "Welcome to your new site", "Update the landing page sections to match your brand.", 0, true, `{"headline":"Build your homepage"}`},
		{"cta", "Publish your first article", "Seeded starting point for the builder", 1, true, `{"buttonText":"Open articles"}`},
	}

	for _, section := range sections {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO landing_sections (site_id, section_key, title, subtitle, content_json, display_order, is_enabled)
			VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7)
			ON CONFLICT (site_id, section_key) DO NOTHING
		`, siteID, section.key, section.title, section.subtitle, section.content, section.order, section.isEnabled); err != nil {
			return err
		}
	}

	_ = authorID
	_ = categoryID
	return nil
}

func fallbackString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func fallbackJSON(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
