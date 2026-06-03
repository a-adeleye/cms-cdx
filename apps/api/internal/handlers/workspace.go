package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"cms-builder/api/internal/builder"
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
	LandingSections   []landingSectionResp `json:"landingSections"`
	Articles          []articleResponse    `json:"articles"`
	Authors           []authorResponse     `json:"authors"`
	Categories        []categoryResponse   `json:"categories"`
	Tags              []tagResponse        `json:"tags"`
	MediaAssets       []mediaAssetResponse `json:"mediaAssets"`
	Builds            []buildResponse      `json:"builds"`
}

type siteResponse struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Slug                  string `json:"slug"`
	Domain                string `json:"domain"`
	BlogPath              string `json:"blogPath"`
	Status                string `json:"status"`
	TemplateKey           string `json:"templateKey"`
	ThemeConfig           string `json:"themeConfig"`
	DeployProvider        string `json:"deployProvider"`
	DeployConfig          string `json:"deployConfig"`
	PreviewDeployProvider string `json:"previewDeployProvider"`
	PreviewDeployConfig   string `json:"previewDeployConfig"`
	AIConfig              string `json:"aiConfig"`
	StorageConfig         string `json:"storageConfig"`
	UpdatedAt             string `json:"updatedAt"`
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
	ID              string   `json:"id"`
	SiteID          string   `json:"siteId"`
	AuthorID        string   `json:"authorId"`
	CategoryID      string   `json:"categoryId"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	ContentMarkdown string   `json:"contentMarkdown"`
	CoverImageURL   string   `json:"coverImageUrl"`
	Status          string   `json:"status"`
	IsFeatured      bool     `json:"isFeatured"`
	PublishedAt     *string  `json:"publishedAt"`
	SEOTitle        string   `json:"seoTitle"`
	SEODescription  string   `json:"seoDescription"`
	CanonicalURL    string   `json:"canonicalUrl"`
	GeneratedByAI   bool     `json:"generatedByAi"`
	HumanReviewed   bool     `json:"humanReviewed"`
	AIPrompt        string   `json:"aiPrompt"`
	AIModel         string   `json:"aiModel"`
	TagIDs          []string `json:"tagIds"`
	UpdatedAt       string   `json:"updatedAt"`
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

type tagResponse struct {
	ID     string `json:"id"`
	SiteID string `json:"siteId"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
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
	StartedAt      *string `json:"startedAt"`
	FinishedAt     *string `json:"finishedAt"`
}

type siteUpsertRequest struct {
	Name                  string `json:"name"`
	Slug                  string `json:"slug"`
	Domain                string `json:"domain"`
	BlogPath              string `json:"blogPath"`
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
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	ContentMarkdown string   `json:"contentMarkdown"`
	CoverImageURL   string   `json:"coverImageUrl"`
	SEOTitle        string   `json:"seoTitle"`
	SEODescription  string   `json:"seoDescription"`
	CanonicalURL    string   `json:"canonicalUrl"`
	AuthorID        string   `json:"authorId"`
	CategoryID      string   `json:"categoryId"`
	TagIDs          []string `json:"tagIds"`
	IsFeatured      bool     `json:"isFeatured"`
	Status          string   `json:"status"`
}

type categoryUpsertRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type authorUpsertRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type tagUpsertRequest struct {
	Name string `json:"name"`
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
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch parts[1] {
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
	case "tags":
		a.handleSiteTagRoutes(w, r, siteID, parts[2:])
	case "builds":
		a.handleBuildRoutes(w, r, siteID)
	case "media":
		a.handleMediaRoutes(w, r, siteID)
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

func (a *API) handleSiteTagRoutes(w http.ResponseWriter, r *http.Request, siteID string, parts []string) {
	switch {
	case len(parts) == 0 || parts[0] == "":
		switch r.Method {
		case http.MethodGet:
			tags, err := a.listTags(r.Context(), siteID)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load tags"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": tags})
		case http.MethodPost:
			var payload tagUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			tag, err := a.createTag(r.Context(), siteID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to create tag"})
				return
			}
			writeJSON(w, http.StatusCreated, tag)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		tagID := parts[0]
		switch r.Method {
		case http.MethodPatch:
			var payload tagUpsertRequest
			if err := decodeJSON(r, &payload); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
				return
			}
			tag, err := a.updateTag(r.Context(), siteID, tagID, payload)
			if err != nil {
				if errors.Is(err, errValidation) {
					writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to update tag"})
				return
			}
			writeJSON(w, http.StatusOK, tag)
		case http.MethodDelete:
			if err := a.deleteTag(r.Context(), siteID, tagID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.NotFound(w, r)
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to delete tag"})
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
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleMediaRoutes(w http.ResponseWriter, r *http.Request, siteID string) {
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
			r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
			if err := r.ParseMultipartForm(12 << 20); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart payload"})
				return
			}

			file, header, err := r.FormFile("file")
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file is required"})
				return
			}
			defer file.Close()

			contents, err := io.ReadAll(io.LimitReader(file, 12<<20+1))
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to read uploaded file"})
				return
			}
			if len(contents) > 12<<20 {
				writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "file is too large"})
				return
			}
			if len(contents) == 0 {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file is empty"})
				return
			}

			mimeType := strings.TrimSpace(header.Header.Get("Content-Type"))
			if mimeType == "" {
				mimeType = http.DetectContentType(contents)
			}
			if !strings.HasPrefix(mimeType, "image/") {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only image uploads are supported"})
				return
			}

			altText := strings.TrimSpace(r.FormValue("altText"))
			media, err := a.uploadMediaAsset(r.Context(), siteID, header.Filename, contents, mimeType, altText)
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

var errValidation = errors.New("validation error")

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
	tags, err := a.listTags(ctx, selectedSiteID)
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
		LandingSections:   landingSections,
		Articles:          articles,
		Authors:           authors,
		Categories:        categories,
		Tags:              tags,
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
		if err := rows.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Status, &site.TemplateKey, &site.ThemeConfig, &site.DeployProvider, &site.DeployConfig, &site.PreviewDeployProvider, &site.PreviewDeployConfig, &site.AIConfig, &site.StorageConfig, &updatedAt); err != nil {
			return nil, err
		}
		site.UpdatedAt = updatedAt.Format(time.RFC3339)
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
	if err := row.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Status, &site.TemplateKey, &site.ThemeConfig, &site.DeployProvider, &site.DeployConfig, &site.PreviewDeployProvider, &site.PreviewDeployConfig, &site.AIConfig, &site.StorageConfig, &updatedAt); err != nil {
		return siteResponse{}, err
	}
	site.UpdatedAt = updatedAt.Format(time.RFC3339)
	return site, nil
}

func (a *API) createSite(ctx context.Context, payload siteUpsertRequest) (siteResponse, error) {
	if err := validateSitePayload(payload); err != nil {
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
			name, slug, domain, blog_path, status, template_key, theme_config, deploy_provider, deploy_config, preview_deploy_provider, preview_deploy_config, ai_config, storage_config
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7::jsonb, NULLIF($8, ''), $9::jsonb, NULLIF($10, ''), $11::jsonb, $12::jsonb, $13::jsonb)
		RETURNING id::text, updated_at
	`, payload.Name, payload.Slug, payload.Domain, fallbackString(payload.BlogPath, "/articles"), fallbackString(payload.Status, "active"), fallbackString(payload.TemplateKey, "default-blog"), fallbackJSON(payload.ThemeConfig, `{"tone":"professional"}`), payload.DeployProvider, fallbackJSON(payload.DeployConfig, `{}`), payload.PreviewDeployProvider, fallbackJSON(payload.PreviewDeployConfig, `{}`), fallbackJSON(payload.AIConfig, `{}`), fallbackJSON(payload.StorageConfig, `{}`)).Scan(&siteID, &updatedAt)
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

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE sites
		SET
			name = $2,
			slug = $3,
			domain = NULLIF($4, ''),
			blog_path = $5,
			status = $6,
			template_key = $7,
			theme_config = $8::jsonb,
			deploy_provider = NULLIF($9, ''),
			deploy_config = $10::jsonb,
			preview_deploy_provider = NULLIF($11, ''),
			preview_deploy_config = $12::jsonb,
			ai_config = $13::jsonb,
			storage_config = $14::jsonb,
			updated_at = NOW()
		WHERE id = $1
	`, siteID, payload.Name, payload.Slug, payload.Domain, fallbackString(payload.BlogPath, "/articles"), fallbackString(payload.Status, "active"), fallbackString(payload.TemplateKey, "default-blog"), fallbackJSON(payload.ThemeConfig, `{}`), payload.DeployProvider, fallbackJSON(payload.DeployConfig, `{}`), payload.PreviewDeployProvider, fallbackJSON(payload.PreviewDeployConfig, `{}`), fallbackJSON(payload.AIConfig, `{}`), fallbackJSON(payload.StorageConfig, `{}`))
	if err != nil {
		return siteResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return siteResponse{}, sql.ErrNoRows
	}
	return a.getSite(ctx, siteID)
}

func validateSitePayload(payload siteUpsertRequest) error {
	if strings.TrimSpace(payload.Name) == "" || strings.TrimSpace(payload.Slug) == "" {
		return fmt.Errorf("%w: name and slug are required", errValidation)
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
			COALESCE(string_agg(at.tag_id::text, ',' ORDER BY at.tag_id), ''),
			a.updated_at
		FROM articles a
		LEFT JOIN article_tags at ON at.article_id = a.id
		WHERE a.site_id = $1
		GROUP BY a.id
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
		var tagIDs string
		var updatedAt time.Time
		if err := rows.Scan(&article.ID, &article.SiteID, &article.AuthorID, &article.CategoryID, &article.Title, &article.Slug, &article.Excerpt, &article.ContentMarkdown, &article.CoverImageURL, &article.Status, &article.IsFeatured, &publishedAt, &article.SEOTitle, &article.SEODescription, &article.CanonicalURL, &article.GeneratedByAI, &article.HumanReviewed, &article.AIPrompt, &article.AIModel, &tagIDs, &updatedAt); err != nil {
			return nil, err
		}
		if publishedAt.Valid {
			value := publishedAt.Time.UTC().Format(time.RFC3339)
			article.PublishedAt = &value
		}
		article.TagIDs = splitIDs(tagIDs)
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
			COALESCE(string_agg(at.tag_id::text, ',' ORDER BY at.tag_id), ''),
			a.updated_at
		FROM articles a
		LEFT JOIN article_tags at ON at.article_id = a.id
		WHERE a.id = $1
		GROUP BY a.id
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
	var tagIDs string
	var updatedAt time.Time
	if err := rows.Scan(&article.ID, &article.SiteID, &article.AuthorID, &article.CategoryID, &article.Title, &article.Slug, &article.Excerpt, &article.ContentMarkdown, &article.CoverImageURL, &article.Status, &article.IsFeatured, &publishedAt, &article.SEOTitle, &article.SEODescription, &article.CanonicalURL, &article.GeneratedByAI, &article.HumanReviewed, &article.AIPrompt, &article.AIModel, &tagIDs, &updatedAt); err != nil {
		return articleResponse{}, err
	}
	if publishedAt.Valid {
		value := publishedAt.Time.UTC().Format(time.RFC3339)
		article.PublishedAt = &value
	}
	article.TagIDs = splitIDs(tagIDs)
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

	articleID := strings.TrimSpace(payload.ID)
	if articleID == "" {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO articles (
				site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, generated_by_ai, human_reviewed, ai_prompt, ai_model, updated_at
			)
			VALUES (
				$1, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, $4, $5, $6, $7, NULLIF($8, ''), $9, $10, CASE WHEN $9 = 'published' THEN NOW() ELSE NULL END, $11, $12, NULLIF($13, ''), FALSE, CASE WHEN $9 IN ('review', 'published') THEN TRUE ELSE FALSE END, '', '', NOW()
			)
			RETURNING id::text
		`, siteID, payload.AuthorID, payload.CategoryID, payload.Title, payload.Slug, payload.Excerpt, payload.ContentMarkdown, payload.CoverImageURL, fallbackString(payload.Status, "draft"), payload.IsFeatured, payload.SEOTitle, payload.SEODescription, payload.CanonicalURL).Scan(&articleID)
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
				human_reviewed = CASE WHEN $9 IN ('review', 'published') THEN TRUE ELSE human_reviewed END,
				updated_at = NOW()
			WHERE id = $1 AND site_id = $14
		`, payload.ID, payload.AuthorID, payload.CategoryID, payload.Title, payload.Slug, payload.Excerpt, payload.ContentMarkdown, payload.CoverImageURL, fallbackString(payload.Status, "draft"), payload.IsFeatured, payload.SEOTitle, payload.SEODescription, payload.CanonicalURL, siteID)
		if err != nil {
			return articleResponse{}, err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return articleResponse{}, sql.ErrNoRows
		}
		articleID = payload.ID
	}

	if err := replaceArticleTags(ctx, tx, articleID, payload.TagIDs); err != nil {
		return articleResponse{}, err
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

	if _, err = tx.ExecContext(ctx, `DELETE FROM article_tags WHERE article_id = $1`, articleID); err != nil {
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

func replaceArticleTags(ctx context.Context, tx *sql.Tx, articleID string, tagIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM article_tags WHERE article_id = $1`, articleID); err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		tagID = strings.TrimSpace(tagID)
		if tagID == "" {
			continue
		}
		if !isUUID(tagID) {
			return fmt.Errorf("%w: invalid tag id %q", errValidation, tagID)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, articleID, tagID); err != nil {
			return err
		}
	}
	return nil
}

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func isUUID(value string) bool {
	return uuidPattern.MatchString(strings.TrimSpace(value))
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

func (a *API) listTags(ctx context.Context, siteID string) ([]tagResponse, error) {
	rows, err := a.Services.DB.QueryContext(ctx, `
		SELECT id::text, site_id::text, name, slug
		FROM tags
		WHERE site_id = $1
		ORDER BY name ASC
	`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]tagResponse, 0)
	for rows.Next() {
		var item tagResponse
		if err := rows.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug); err != nil {
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

func (a *API) getTag(ctx context.Context, siteID, tagID string) (tagResponse, error) {
	row := a.Services.DB.QueryRowContext(ctx, `
		SELECT id::text, site_id::text, name, slug
		FROM tags
		WHERE site_id = $1 AND id = $2
	`, siteID, tagID)

	var item tagResponse
	if err := row.Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug); err != nil {
		return tagResponse{}, err
	}
	return item, nil
}

func (a *API) createTag(ctx context.Context, siteID string, payload tagUpsertRequest) (tagResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return tagResponse{}, fmt.Errorf("%w: tag name is required", errValidation)
	}

	slug, err := a.uniqueTagSlug(ctx, siteID, "", name)
	if err != nil {
		return tagResponse{}, err
	}

	var item tagResponse
	err = a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO tags (site_id, name, slug)
		VALUES ($1, $2, $3)
		RETURNING id::text, site_id::text, name, slug
	`, siteID, name, slug).Scan(&item.ID, &item.SiteID, &item.Name, &item.Slug)
	if err != nil {
		if isUniqueViolation(err) {
			return tagResponse{}, fmt.Errorf("%w: tag slug already exists", errValidation)
		}
		return tagResponse{}, err
	}
	return item, nil
}

func (a *API) updateTag(ctx context.Context, siteID, tagID string, payload tagUpsertRequest) (tagResponse, error) {
	name := strings.TrimSpace(payload.Name)
	if name == "" {
		return tagResponse{}, fmt.Errorf("%w: tag name is required", errValidation)
	}

	slug, err := a.uniqueTagSlug(ctx, siteID, tagID, name)
	if err != nil {
		return tagResponse{}, err
	}

	result, err := a.Services.DB.ExecContext(ctx, `
		UPDATE tags
		SET name = $3, slug = $4, updated_at = NOW()
		WHERE id = $1 AND site_id = $2
	`, tagID, siteID, name, slug)
	if err != nil {
		if isUniqueViolation(err) {
			return tagResponse{}, fmt.Errorf("%w: tag slug already exists", errValidation)
		}
		return tagResponse{}, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return tagResponse{}, sql.ErrNoRows
	}
	return a.getTag(ctx, siteID, tagID)
}

func (a *API) deleteTag(ctx context.Context, siteID, tagID string) error {
	result, err := a.Services.DB.ExecContext(ctx, `
		DELETE FROM tags
		WHERE id = $1 AND site_id = $2
	`, tagID, siteID)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (a *API) uniqueCategorySlug(ctx context.Context, siteID, categoryID, name string) (string, error) {
	return a.uniqueSlug(ctx, "categories", siteID, categoryID, name)
}

func (a *API) uniqueAuthorSlug(ctx context.Context, siteID, authorID, name string) (string, error) {
	return a.uniqueSlug(ctx, "authors", siteID, authorID, name)
}

func (a *API) uniqueTagSlug(ctx context.Context, siteID, tagID, name string) (string, error) {
	return a.uniqueSlug(ctx, "tags", siteID, tagID, name)
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
			COALESCE(alt_text, '')
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
		if err := rows.Scan(&item.ID, &item.SiteID, &item.FileName, &item.FileURL, &item.MimeType, &item.SizeBytes, &item.StorageProvider, &item.StorageKey, &item.AltText); err != nil {
			return nil, err
		}
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

	stored, err := a.Services.Storage.Upload(ctx, storage.UploadFile{
		FileName: fileName,
		Contents: contents,
		MimeType: mimeType,
		SiteID:   siteID,
	})
	if err != nil {
		return mediaAssetResponse{}, err
	}
	if stored == nil || strings.TrimSpace(stored.PublicURL) == "" {
		return mediaAssetResponse{}, fmt.Errorf("%w: storage did not return a public URL", errValidation)
	}

	return a.persistMediaAsset(ctx, siteID, mediaUpsertRequest{
		FileName:        fileName,
		FileURL:         stored.PublicURL,
		MimeType:        mimeType,
		SizeBytes:       int64(len(contents)),
		StorageProvider: a.storageProviderName(),
		StorageKey:      stored.Key,
		AltText:         altText,
	})
}

func (a *API) persistMediaAsset(ctx context.Context, siteID string, payload mediaUpsertRequest) (mediaAssetResponse, error) {
	var item mediaAssetResponse
	err := a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO media_assets (
			site_id, file_name, file_url, mime_type, size_bytes, storage_provider, storage_key, alt_text
		)
		VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, 0), $6, NULLIF($7, ''), NULLIF($8, ''))
		RETURNING id::text, site_id::text, file_name, file_url, COALESCE(mime_type, ''), COALESCE(size_bytes, 0), storage_provider, COALESCE(storage_key, ''), COALESCE(alt_text, '')
	`, siteID, payload.FileName, payload.FileURL, payload.MimeType, payload.SizeBytes, fallbackString(payload.StorageProvider, "s3"), payload.StorageKey, payload.AltText).Scan(&item.ID, &item.SiteID, &item.FileName, &item.FileURL, &item.MimeType, &item.SizeBytes, &item.StorageProvider, &item.StorageKey, &item.AltText)
	if err != nil {
		return mediaAssetResponse{}, err
	}
	return item, nil
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
		if err := rows.Scan(&item.ID, &item.SiteID, &item.Status, &item.BuildType, &item.Logs, &item.OutputPath, &item.DeployProvider, &item.DeployStatus, &item.DeployURL, &startedAt, &finishedAt); err != nil {
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
			started_at,
			finished_at
		FROM builds
		WHERE id = $1
	`, buildID)

	var item buildResponse
	var startedAt sql.NullTime
	var finishedAt sql.NullTime
	if err := row.Scan(&item.ID, &item.SiteID, &item.Status, &item.BuildType, &item.Logs, &item.OutputPath, &item.DeployProvider, &item.DeployStatus, &item.DeployURL, &startedAt, &finishedAt); err != nil {
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

	if buildType == "published" && len(payload.ArticleIDs) > 0 {
		if err := a.publishArticles(ctx, siteID, payload.ArticleIDs); err != nil {
			return buildResponse{}, err
		}
	}

	content, err := a.siteBuildContent(ctx, siteID, buildType == "published")
	if err != nil {
		return buildResponse{}, err
	}

	buildSite := site
	if buildType == "preview" {
		buildSite.DeployProvider = fallbackString(site.PreviewDeployProvider, site.DeployProvider)
		buildSite.DeployConfig = site.PreviewDeployConfig
	}

	outputPath, err := a.Services.Builder.GenerateSite(ctx, content, builder.GenerateOptions{
		SiteID:     siteID,
		Preview:    buildType == "preview",
		ArticleIDs: payload.ArticleIDs,
	})
	if err != nil {
		return buildResponse{}, err
	}

	build := models.Build{
		SiteID:     siteID,
		Status:     "success",
		BuildType:  buildType,
		Logs:       buildLogs(buildType, buildArticleCount(buildType, len(payload.ArticleIDs), len(content.Articles))),
		OutputPath: outputPath,
	}
	deployResult, err := a.Services.Deploy.Deploy(ctx, buildSite, build, outputPath)
	if err != nil {
		return buildResponse{}, err
	}

	deployProvider := fallbackString(buildSite.DeployProvider, "none")
	deployStatus := "deployed"
	deployURL := ""
	if deployResult != nil {
		deployProvider = fallbackString(deployResult.Provider, deployProvider)
		deployStatus = "deployed"
		deployURL = strings.TrimSpace(deployResult.URL)
		if message := strings.TrimSpace(deployResult.Message); message != "" {
			build.Logs = strings.TrimSpace(build.Logs + " " + message)
		}
	}

	var item buildResponse
	now := time.Now().UTC().Format(time.RFC3339)
	var startedAt sql.NullTime
	var finishedAt sql.NullTime
	err = a.Services.DB.QueryRowContext(ctx, `
		INSERT INTO builds (
			site_id, status, build_type, logs, output_path, deploy_provider, deploy_status, deploy_url, started_at, finished_at
		)
		VALUES ($1, 'success', $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id::text, site_id::text, status, build_type, logs, output_path, COALESCE(deploy_provider, ''), COALESCE(deploy_status, ''), COALESCE(deploy_url, ''), started_at, finished_at
	`, siteID, buildType, build.Logs, outputPath, deployProvider, deployStatus, deployURL).Scan(&item.ID, &item.SiteID, &item.Status, &item.BuildType, &item.Logs, &item.OutputPath, &item.DeployProvider, &item.DeployStatus, &item.DeployURL, &startedAt, &finishedAt)
	if err != nil {
		return buildResponse{}, err
	}
	item.DeployProvider = deployProvider
	item.DeployStatus = deployStatus
	item.DeployURL = deployURL
	item.Logs = build.Logs
	if startedAt.Valid {
		value := startedAt.Time.UTC().Format(time.RFC3339)
		item.StartedAt = &value
	} else {
		item.StartedAt = &now
	}
	if finishedAt.Valid {
		value := finishedAt.Time.UTC().Format(time.RFC3339)
		item.FinishedAt = &value
	} else {
		item.FinishedAt = &now
	}
	return item, nil
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
	tags, err := a.listTags(ctx, siteID)
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
	tagNames := make(map[string]tagResponse, len(tags))
	for _, tag := range tags {
		tagNames[tag.ID] = tag
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
		for _, tagID := range article.TagIDs {
			if tag, ok := tagNames[tagID]; ok {
				articleContent.Tags = append(articleContent.Tags, builder.TagContent{Name: tag.Name, Slug: tag.Slug})
			}
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
	if err := row.Scan(&site.ID, &site.Name, &site.Slug, &site.Domain, &site.BlogPath, &site.Status, &site.TemplateKey, &themeConfig, &site.DeployProvider, &deployConfig, &site.PreviewDeployProvider, &previewDeployConfig, &aiConfig, &storageConfig, &site.CreatedAt, &site.UpdatedAt); err != nil {
		return models.Site{}, err
	}
	site.ThemeConfig = parseJSONMap(themeConfig)
	site.DeployConfig = parseJSONMap(deployConfig)
	site.PreviewDeployConfig = parseJSONMap(previewDeployConfig)
	site.AIConfig = parseJSONMap(aiConfig)
	site.StorageConfig = parseJSONMap(storageConfig)
	return site, nil
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

	var tagID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO tags (site_id, name, slug)
		VALUES ($1, 'starter', 'starter')
		RETURNING id::text
	`, siteID).Scan(&tagID); err != nil {
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
	_ = tagID
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

func splitIDs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			ids = append(ids, part)
		}
	}
	return ids
}
