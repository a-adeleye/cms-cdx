package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cms-builder/api/internal/middleware"
)

const (
	siteExportVersion  = 1
	maxSiteImportBytes = 10 << 20
	maxImportedRecords = 10000
)

// siteExportBundle is a portable, versioned representation of editable site
// content. Deployment, storage, build history, and secret references are kept
// out of the bundle so importing it cannot publish or reuse credentials.
type siteExportBundle struct {
	Version         int                  `json:"version"`
	ExportedAt      string               `json:"exportedAt"`
	Site            siteExportSite       `json:"site"`
	LandingSections []siteExportSection  `json:"landingSections"`
	Authors         []siteExportAuthor   `json:"authors"`
	Categories      []siteExportCategory `json:"categories"`
	MediaAssets     []siteExportMedia    `json:"mediaAssets"`
	Articles        []siteExportArticle  `json:"articles"`
}

type siteExportSite struct {
	Name           string         `json:"name"`
	Slug           string         `json:"slug"`
	Domain         string         `json:"domain"`
	BlogPath       string         `json:"blogPath"`
	Description    string         `json:"description"`
	ContentContext string         `json:"contentContext"`
	Status         string         `json:"status"`
	TemplateKey    string         `json:"templateKey"`
	ThemeConfig    map[string]any `json:"themeConfig"`
	AIConfig       map[string]any `json:"aiConfig"`
	LogoMediaID    string         `json:"logoMediaId"`
	FaviconMediaID string         `json:"faviconMediaId"`
}

type siteExportSection struct {
	SectionKey   string         `json:"sectionKey"`
	Title        string         `json:"title"`
	Subtitle     string         `json:"subtitle"`
	Content      map[string]any `json:"content"`
	DisplayOrder int            `json:"displayOrder"`
	IsEnabled    bool           `json:"isEnabled"`
}

type siteExportAuthor struct {
	SourceID string `json:"sourceId"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Bio      string `json:"bio"`
}

type siteExportCategory struct {
	SourceID    string `json:"sourceId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type siteExportMedia struct {
	SourceID  string `json:"sourceId"`
	FileName  string `json:"fileName"`
	FileURL   string `json:"fileUrl"`
	MimeType  string `json:"mimeType"`
	SizeBytes int64  `json:"sizeBytes"`
	AltText   string `json:"altText"`
}

type siteExportArticle struct {
	AuthorSourceID   string  `json:"authorSourceId"`
	CategorySourceID string  `json:"categorySourceId"`
	Title            string  `json:"title"`
	Slug             string  `json:"slug"`
	Excerpt          string  `json:"excerpt"`
	ContentMarkdown  string  `json:"contentMarkdown"`
	CoverImageURL    string  `json:"coverImageUrl"`
	Status           string  `json:"status"`
	IsFeatured       bool    `json:"isFeatured"`
	PublishedAt      *string `json:"publishedAt"`
	SEOTitle         string  `json:"seoTitle"`
	SEODescription   string  `json:"seoDescription"`
	CanonicalURL     string  `json:"canonicalUrl"`
	GeneratedByAI    bool    `json:"generatedByAi"`
	HumanReviewed    bool    `json:"humanReviewed"`
	AIPrompt         string  `json:"aiPrompt"`
	AIModel          string  `json:"aiModel"`
	Tags             string  `json:"tags"`
}

func (a *API) handleSiteExport(w http.ResponseWriter, r *http.Request, siteID string, remaining []string) {
	if len(remaining) != 0 || r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if !requireAdmin(w, r) {
		return
	}
	bundle, err := a.exportSite(r.Context(), siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "site not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to export site"})
		return
	}
	if err := a.recordSiteExport(r.Context(), middleware.UserIDFromContext(r.Context()), siteID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to record site export"})
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", bundle.Site.Slug+"-site-export.json"))
	writeJSON(w, http.StatusOK, bundle)
}

func (a *API) siteImports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if !requireAdmin(w, r) {
		return
	}
	if a.Services.DB == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}
	var bundle siteExportBundle
	if err := decodeSiteImport(w, r, &bundle); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid site export file"})
		return
	}
	site, err := a.importSite(r.Context(), bundle, middleware.UserIDFromContext(r.Context()))
	if err != nil {
		if errors.Is(err, errValidation) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to import site"})
		return
	}
	writeJSON(w, http.StatusCreated, site)
}

func decodeSiteImport(w http.ResponseWriter, r *http.Request, destination *siteExportBundle) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxSiteImportBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("site import contains multiple JSON values")
	}
	return nil
}

func (a *API) exportSite(ctx context.Context, siteID string) (siteExportBundle, error) {
	site, err := a.getSite(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	sections, err := a.listLandingSections(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	authors, err := a.listAuthors(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	categories, err := a.listCategories(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	mediaAssets, err := a.listMediaAssets(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	articles, err := a.listArticles(ctx, siteID)
	if err != nil {
		return siteExportBundle{}, err
	}
	themeConfig, err := parseExportConfig(site.ThemeConfig)
	if err != nil {
		return siteExportBundle{}, err
	}
	aiConfig, err := parseExportConfig(site.AIConfig)
	if err != nil {
		return siteExportBundle{}, err
	}

	bundle := siteExportBundle{
		Version:    siteExportVersion,
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Site: siteExportSite{
			Name: site.Name, Slug: site.Slug, Domain: site.Domain, BlogPath: site.BlogPath,
			Description: site.Description, ContentContext: site.ContentContext, Status: site.Status,
			TemplateKey: site.TemplateKey, ThemeConfig: themeConfig, AIConfig: stripSecretReferences(aiConfig),
			LogoMediaID: site.LogoMediaID, FaviconMediaID: site.FaviconMediaID,
		},
		LandingSections: make([]siteExportSection, 0, len(sections)),
		Authors:         make([]siteExportAuthor, 0, len(authors)),
		Categories:      make([]siteExportCategory, 0, len(categories)),
		MediaAssets:     make([]siteExportMedia, 0, len(mediaAssets)),
		Articles:        make([]siteExportArticle, 0, len(articles)),
	}
	for _, section := range sections {
		content, err := parseExportConfig(section.ContentJSON)
		if err != nil {
			return siteExportBundle{}, err
		}
		bundle.LandingSections = append(bundle.LandingSections, siteExportSection{SectionKey: section.SectionKey, Title: section.Title, Subtitle: section.Subtitle, Content: content, DisplayOrder: section.DisplayOrder, IsEnabled: section.IsEnabled})
	}
	for _, author := range authors {
		bundle.Authors = append(bundle.Authors, siteExportAuthor{SourceID: author.ID, Name: author.Name, Slug: author.Slug, Bio: author.Bio})
	}
	for _, category := range categories {
		bundle.Categories = append(bundle.Categories, siteExportCategory{SourceID: category.ID, Name: category.Name, Slug: category.Slug, Description: category.Description})
	}
	for _, asset := range mediaAssets {
		bundle.MediaAssets = append(bundle.MediaAssets, siteExportMedia{SourceID: asset.ID, FileName: asset.FileName, FileURL: asset.FileURL, MimeType: asset.MimeType, SizeBytes: asset.SizeBytes, AltText: asset.AltText})
	}
	for _, article := range articles {
		bundle.Articles = append(bundle.Articles, siteExportArticle{AuthorSourceID: article.AuthorID, CategorySourceID: article.CategoryID, Title: article.Title, Slug: article.Slug, Excerpt: article.Excerpt, ContentMarkdown: article.ContentMarkdown, CoverImageURL: article.CoverImageURL, Status: article.Status, IsFeatured: article.IsFeatured, PublishedAt: article.PublishedAt, SEOTitle: article.SEOTitle, SEODescription: article.SEODescription, CanonicalURL: article.CanonicalURL, GeneratedByAI: article.GeneratedByAI, HumanReviewed: article.HumanReviewed, AIPrompt: article.AIPrompt, AIModel: article.AIModel, Tags: article.Tags})
	}
	return bundle, nil
}

func (a *API) importSite(ctx context.Context, bundle siteExportBundle, actorID string) (siteResponse, error) {
	if err := validateSiteExportBundle(bundle); err != nil {
		return siteResponse{}, err
	}
	themeConfig, err := json.Marshal(bundle.Site.ThemeConfig)
	if err != nil {
		return siteResponse{}, fmt.Errorf("%w: invalid theme configuration", errValidation)
	}
	aiConfig, err := json.Marshal(stripSecretReferences(bundle.Site.AIConfig))
	if err != nil {
		return siteResponse{}, fmt.Errorf("%w: invalid AI configuration", errValidation)
	}
	payload := siteUpsertRequest{
		Name: bundle.Site.Name, Slug: bundle.Site.Slug, Domain: bundle.Site.Domain, BlogPath: bundle.Site.BlogPath,
		Description: bundle.Site.Description, ContentContext: bundle.Site.ContentContext, Status: bundle.Site.Status,
		TemplateKey: bundle.Site.TemplateKey, ThemeConfig: string(themeConfig), AIConfig: string(aiConfig),
		DeployProvider: "none", DeployConfig: "{}", PreviewDeployProvider: "none", PreviewDeployConfig: "{}", StorageConfig: "{}",
	}
	if err := validateSitePayload(payload); err != nil {
		return siteResponse{}, err
	}
	if err := a.templateExists(ctx, payload.TemplateKey); err != nil {
		return siteResponse{}, err
	}
	tx, err := a.Services.DB.BeginTx(ctx, nil)
	if err != nil {
		return siteResponse{}, err
	}
	defer func() { _ = tx.Rollback() }()

	name, slug, err := importedSiteIdentity(ctx, tx, bundle.Site.Name, bundle.Site.Slug)
	if err != nil {
		return siteResponse{}, err
	}
	var siteID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO sites (name, slug, domain, blog_path, description, content_context, status, template_key, theme_config, deploy_provider, deploy_config, preview_deploy_provider, preview_deploy_config, ai_config, storage_config)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, $8, $9::jsonb, 'none', '{}'::jsonb, 'none', '{}'::jsonb, $10::jsonb, '{}'::jsonb)
		RETURNING id::text
	`, name, slug, payload.Domain, payload.BlogPath, payload.Description, siteContentContext(payload.ContentContext), payload.Status, payload.TemplateKey, payload.ThemeConfig, payload.AIConfig).Scan(&siteID)
	if err != nil {
		return siteResponse{}, err
	}

	mediaIDs, err := importMediaAssets(ctx, tx, siteID, bundle.MediaAssets)
	if err != nil {
		return siteResponse{}, err
	}
	authorIDs, err := importAuthors(ctx, tx, siteID, bundle.Authors)
	if err != nil {
		return siteResponse{}, err
	}
	categoryIDs, err := importCategories(ctx, tx, siteID, bundle.Categories)
	if err != nil {
		return siteResponse{}, err
	}
	if err := importLandingSections(ctx, tx, siteID, bundle.LandingSections); err != nil {
		return siteResponse{}, err
	}
	if err := importArticles(ctx, tx, siteID, bundle.Articles, authorIDs, categoryIDs); err != nil {
		return siteResponse{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE sites SET logo_media_id = NULLIF($2, '')::uuid, favicon_media_id = NULLIF($3, '')::uuid WHERE id = $1`, siteID, mediaIDs[bundle.Site.LogoMediaID], mediaIDs[bundle.Site.FaviconMediaID]); err != nil {
		return siteResponse{}, err
	}
	if strings.TrimSpace(actorID) != "" {
		if _, err = tx.ExecContext(ctx, `INSERT INTO audit_logs (user_id, site_id, action, entity_type, entity_id, metadata) VALUES ($1::uuid, $2::uuid, 'import_site', 'site', $2::uuid, '{"version":1}'::jsonb)`, actorID, siteID); err != nil {
			return siteResponse{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return siteResponse{}, err
	}
	return a.getSite(ctx, siteID)
}

func (a *API) recordSiteExport(ctx context.Context, actorID, siteID string) error {
	if strings.TrimSpace(actorID) == "" {
		return errors.New("missing export actor")
	}
	_, err := a.Services.DB.ExecContext(ctx, `INSERT INTO audit_logs (user_id, site_id, action, entity_type, entity_id, metadata) VALUES ($1::uuid, $2::uuid, 'export_site', 'site', $2::uuid, '{"version":1}'::jsonb)`, actorID, siteID)
	return err
}

func validateSiteExportBundle(bundle siteExportBundle) error {
	if bundle.Version != siteExportVersion {
		return fmt.Errorf("%w: unsupported site export version", errValidation)
	}
	if strings.TrimSpace(bundle.Site.Name) == "" || strings.TrimSpace(bundle.Site.Slug) == "" {
		return fmt.Errorf("%w: exported site name and slug are required", errValidation)
	}
	if len(bundle.LandingSections) > maxImportedRecords || len(bundle.Authors) > maxImportedRecords || len(bundle.Categories) > maxImportedRecords || len(bundle.MediaAssets) > maxImportedRecords || len(bundle.Articles) > maxImportedRecords {
		return fmt.Errorf("%w: exported site has too many records", errValidation)
	}
	seenSections := make(map[string]struct{}, len(bundle.LandingSections))
	for _, section := range bundle.LandingSections {
		if strings.TrimSpace(section.SectionKey) == "" {
			return fmt.Errorf("%w: landing section key is required", errValidation)
		}
		if _, found := seenSections[section.SectionKey]; found {
			return fmt.Errorf("%w: exported landing section keys must be unique", errValidation)
		}
		seenSections[section.SectionKey] = struct{}{}
	}
	if err := validateExportIdentifiers(bundle.Authors, func(item siteExportAuthor) string { return item.SourceID }, func(item siteExportAuthor) string { return item.Name }, "author"); err != nil {
		return err
	}
	if err := validateExportIdentifiers(bundle.Categories, func(item siteExportCategory) string { return item.SourceID }, func(item siteExportCategory) string { return item.Name }, "category"); err != nil {
		return err
	}
	if err := validateExportIdentifiers(bundle.MediaAssets, func(item siteExportMedia) string { return item.SourceID }, func(item siteExportMedia) string { return item.FileURL }, "media asset"); err != nil {
		return err
	}
	seenArticleSlugs := make(map[string]struct{}, len(bundle.Articles))
	for _, article := range bundle.Articles {
		if err := validateArticlePayload(articleUpsertRequest{Title: article.Title, Slug: article.Slug, ContentMarkdown: article.ContentMarkdown}); err != nil {
			return err
		}
		if _, found := seenArticleSlugs[article.Slug]; found {
			return fmt.Errorf("%w: exported article slugs must be unique", errValidation)
		}
		seenArticleSlugs[article.Slug] = struct{}{}
		if article.PublishedAt != nil {
			if _, err := time.Parse(time.RFC3339, *article.PublishedAt); err != nil {
				return fmt.Errorf("%w: exported publishedAt must use RFC3339", errValidation)
			}
		}
	}
	return nil
}

func validateExportIdentifiers[T any](items []T, id func(T) string, required func(T) string, label string) error {
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		identifier := strings.TrimSpace(id(item))
		if identifier == "" || strings.TrimSpace(required(item)) == "" {
			return fmt.Errorf("%w: exported %s is incomplete", errValidation, label)
		}
		if _, found := seen[identifier]; found {
			return fmt.Errorf("%w: exported %s IDs must be unique", errValidation, label)
		}
		seen[identifier] = struct{}{}
	}
	return nil
}

func importedSiteIdentity(ctx context.Context, tx *sql.Tx, sourceName, sourceSlug string) (string, string, error) {
	for suffix := 0; suffix < 1000; suffix++ {
		name := sourceName + " (Imported)"
		slug := sourceSlug + "-imported"
		if suffix > 0 {
			name += " " + strconv.Itoa(suffix+1)
			slug += "-" + strconv.Itoa(suffix+1)
		}
		var exists bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM sites WHERE slug = $1)`, slug).Scan(&exists); err != nil {
			return "", "", err
		}
		if !exists {
			return name, slug, nil
		}
	}
	return "", "", fmt.Errorf("%w: unable to allocate an imported site slug", errValidation)
}

func importMediaAssets(ctx context.Context, tx *sql.Tx, siteID string, items []siteExportMedia) (map[string]string, error) {
	ids := make(map[string]string, len(items))
	for _, item := range items {
		var id string
		err := tx.QueryRowContext(ctx, `INSERT INTO media_assets (site_id, file_name, file_url, mime_type, size_bytes, storage_provider, storage_key, alt_text) VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, 0), 'imported', NULL, NULLIF($6, '')) RETURNING id::text`, siteID, item.FileName, item.FileURL, item.MimeType, item.SizeBytes, item.AltText).Scan(&id)
		if err != nil {
			return nil, err
		}
		ids[item.SourceID] = id
	}
	return ids, nil
}

func importAuthors(ctx context.Context, tx *sql.Tx, siteID string, items []siteExportAuthor) (map[string]string, error) {
	ids := make(map[string]string, len(items))
	for _, item := range items {
		var id string
		if err := tx.QueryRowContext(ctx, `INSERT INTO authors (site_id, name, slug, bio) VALUES ($1, $2, $3, NULLIF($4, '')) RETURNING id::text`, siteID, item.Name, item.Slug, item.Bio).Scan(&id); err != nil {
			return nil, err
		}
		ids[item.SourceID] = id
	}
	return ids, nil
}

func importCategories(ctx context.Context, tx *sql.Tx, siteID string, items []siteExportCategory) (map[string]string, error) {
	ids := make(map[string]string, len(items))
	for _, item := range items {
		var id string
		if err := tx.QueryRowContext(ctx, `INSERT INTO categories (site_id, name, slug, description) VALUES ($1, $2, $3, NULLIF($4, '')) RETURNING id::text`, siteID, item.Name, item.Slug, item.Description).Scan(&id); err != nil {
			return nil, err
		}
		ids[item.SourceID] = id
	}
	return ids, nil
}

func importLandingSections(ctx context.Context, tx *sql.Tx, siteID string, items []siteExportSection) error {
	for _, item := range items {
		content, err := json.Marshal(item.Content)
		if err != nil {
			return fmt.Errorf("%w: invalid landing section content", errValidation)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO landing_sections (site_id, section_key, title, subtitle, content_json, display_order, is_enabled) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5::jsonb, $6, $7)`, siteID, item.SectionKey, item.Title, item.Subtitle, string(content), item.DisplayOrder, item.IsEnabled); err != nil {
			return err
		}
	}
	return nil
}

func importArticles(ctx context.Context, tx *sql.Tx, siteID string, items []siteExportArticle, authorIDs, categoryIDs map[string]string) error {
	for _, item := range items {
		authorID, categoryID := authorIDs[item.AuthorSourceID], categoryIDs[item.CategorySourceID]
		if item.AuthorSourceID != "" && authorID == "" || item.CategorySourceID != "" && categoryID == "" {
			return fmt.Errorf("%w: article references an unknown author or category", errValidation)
		}
		var publishedAt any
		if item.PublishedAt != nil {
			value, err := time.Parse(time.RFC3339, *item.PublishedAt)
			if err != nil {
				return fmt.Errorf("%w: exported publishedAt must use RFC3339", errValidation)
			}
			publishedAt = value
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO articles (site_id, author_id, category_id, title, slug, excerpt, content_markdown, cover_image_url, status, is_featured, published_at, seo_title, seo_description, canonical_url, tags, generated_by_ai, human_reviewed, ai_prompt, ai_model)
			VALUES ($1, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, $4, $5, NULLIF($6, ''), $7, NULLIF($8, ''), $9, $10, $11, NULLIF($12, ''), NULLIF($13, ''), NULLIF($14, ''), $15, $16, $17, NULLIF($18, ''), NULLIF($19, ''))
		`, siteID, authorID, categoryID, item.Title, item.Slug, item.Excerpt, item.ContentMarkdown, item.CoverImageURL, fallbackString(item.Status, "draft"), item.IsFeatured, publishedAt, item.SEOTitle, item.SEODescription, item.CanonicalURL, normalizeTagsInput(item.Tags), item.GeneratedByAI, item.HumanReviewed, item.AIPrompt, item.AIModel); err != nil {
			return err
		}
	}
	return nil
}

func parseExportConfig(raw string) (map[string]any, error) {
	values := make(map[string]any)
	if err := json.Unmarshal([]byte(fallbackJSON(raw, `{}`)), &values); err != nil {
		return nil, err
	}
	return values, nil
}

func stripSecretReferences(values map[string]any) map[string]any {
	result := make(map[string]any, len(values))
	for key, value := range values {
		if strings.Contains(strings.ToLower(key), "secret") {
			continue
		}
		switch nested := value.(type) {
		case map[string]any:
			result[key] = stripSecretReferences(nested)
		case []any:
			items := make([]any, 0, len(nested))
			for _, entry := range nested {
				if nestedMap, ok := entry.(map[string]any); ok {
					items = append(items, stripSecretReferences(nestedMap))
					continue
				}
				items = append(items, entry)
			}
			result[key] = items
		default:
			result[key] = value
		}
	}
	return result
}
