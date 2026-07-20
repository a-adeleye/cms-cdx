package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cms-builder/api/internal/ai"
)

const (
	maxAISuggestionRequestBytes   = 96 << 10
	maxAIArticleDraftRequestBytes = 12 << 10
	maxGeneratedImageBytes        = 12 << 20
)

type aiSuggestionRequest struct {
	Instruction     string `json:"instruction"`
	Title           string `json:"title"`
	Excerpt         string `json:"excerpt"`
	ContentMarkdown string `json:"contentMarkdown"`
}

type aiSuggestionResponse struct {
	Suggestion string `json:"suggestion"`
	Model      string `json:"model"`
}

type aiArticleDraftRequest struct {
	Topic string `json:"topic"`
}

type aiArticleDraftResponse struct {
	Title           string              `json:"title"`
	Slug            string              `json:"slug"`
	Category        string              `json:"category"`
	Featured        bool                `json:"featured"`
	Tags            []string            `json:"tags"`
	SEOTitle        string              `json:"seoTitle"`
	MetaDescription string              `json:"metaDescription"`
	CanonicalURL    string              `json:"canonicalUrl"`
	Excerpt         string              `json:"excerpt"`
	ContentMarkdown string              `json:"contentMarkdown"`
	CoverImage      *mediaAssetResponse `json:"coverImage,omitempty"`
	ImageError      string              `json:"imageError,omitempty"`
	Model           string              `json:"model"`
}

// handleAISuggestion generates a proposed editorial revision for a site's
// configured provider. It never accepts API keys from the client.
func (a *API) handleAISuggestion(w http.ResponseWriter, r *http.Request, siteID string, subroutes []string) {
	if len(subroutes) != 1 {
		http.NotFound(w, r)
		return
	}
	if subroutes[0] == "article-drafts" {
		a.handleAIArticleDraft(w, r, siteID)
		return
	}
	if subroutes[0] != "suggestions" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAISuggestionRequestBytes)
	var payload aiSuggestionRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid AI suggestion payload"})
		return
	}
	if err := validateAISuggestionRequest(payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	site, err := a.getSite(r.Context(), siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load AI configuration"})
		return
	}
	config, err := ai.ParseConfig(site.AIConfig)
	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "AI provider configuration is invalid. Update it in AI settings."})
		return
	}
	if a.Services.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI generation is unavailable. Try again shortly."})
		return
	}

	suggestion, err := a.Services.AI.GenerateSuggestion(r.Context(), config, ai.SuggestionInput{
		Instruction:     payload.Instruction,
		Title:           payload.Title,
		Excerpt:         payload.Excerpt,
		ContentMarkdown: payload.ContentMarkdown,
	})
	if err != nil {
		a.writeAISuggestionError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, aiSuggestionResponse{Suggestion: suggestion.Content, Model: suggestion.Model})
}

func (a *API) handleAIArticleDraft(w http.ResponseWriter, r *http.Request, siteID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAIArticleDraftRequestBytes)
	var payload aiArticleDraftRequest
	if err := decodeJSON(r, &payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid AI article draft payload"})
		return
	}
	if err := validateAIArticleDraftRequest(payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	site, err := a.getSite(r.Context(), siteID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load AI configuration"})
		return
	}
	config, err := ai.ParseConfig(site.AIConfig)
	if err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "AI provider configuration is invalid. Update it in AI settings."})
		return
	}
	if strings.TrimSpace(config.MasterPrompt) == "" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Add and save an AI master prompt in Site configuration before generating an article."})
		return
	}
	if a.Services.AI == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "AI generation is unavailable. Try again shortly."})
		return
	}

	categories, err := a.listCategories(r.Context(), siteID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load article categories"})
		return
	}
	if len(categories) == 0 {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Create at least one category before generating an article."})
		return
	}
	articles, err := a.listArticles(r.Context(), siteID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to load existing articles"})
		return
	}

	draft, err := a.Services.AI.GenerateArticleDraft(r.Context(), config, ai.GenerateArticleDraftInput{
		Topic:            strings.TrimSpace(payload.Topic),
		SiteName:         site.Name,
		SiteDescription:  site.Description,
		BlogURL:          strings.TrimRight(site.Domain, "/") + "/" + strings.Trim(site.BlogPath, "/"),
		ContentContext:   site.ContentContext,
		MasterPrompt:     config.MasterPrompt,
		Categories:       categoryNames(categories),
		ExistingArticles: articleReferences(articles, categories),
	})
	if err != nil {
		a.writeAISuggestionError(w, err)
		return
	}
	if !hasCategory(categories, draft.Category) {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "The AI response used an invalid article category. Try again."})
		return
	}

	response := aiArticleDraftResponse{
		Title:           draft.Title,
		Slug:            draft.Slug,
		Category:        draft.Category,
		Featured:        draft.Featured,
		Tags:            draft.Tags,
		SEOTitle:        draft.SEOTitle,
		MetaDescription: draft.MetaDescription,
		CanonicalURL:    draft.CanonicalURL,
		Excerpt:         draft.Excerpt,
		ContentMarkdown: draft.Content,
		ImageError:      draft.ImageError,
		Model:           draft.Model,
	}
	if draft.Image != nil {
		coverImage, imageError := a.storeGeneratedArticleImage(r.Context(), siteID, draft)
		if imageError != nil {
			response.ImageError = "The article was generated, but its featured image could not be saved."
		} else {
			response.CoverImage = &coverImage
		}
	}
	writeJSON(w, http.StatusOK, response)
}

func validateAIArticleDraftRequest(payload aiArticleDraftRequest) error {
	topic := strings.TrimSpace(payload.Topic)
	if len(topic) > 4000 {
		return fmt.Errorf("%w: article topic is too long", errValidation)
	}
	return nil
}

func categoryNames(categories []categoryResponse) []string {
	names := make([]string, 0, len(categories))
	for _, category := range categories {
		names = append(names, category.Name)
	}
	return names
}

func hasCategory(categories []categoryResponse, name string) bool {
	for _, category := range categories {
		if strings.EqualFold(strings.TrimSpace(category.Name), strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}

func articleReferences(articles []articleResponse, categories []categoryResponse) []ai.ArticleReference {
	categoryNamesByID := make(map[string]string, len(categories))
	for _, category := range categories {
		categoryNamesByID[category.ID] = category.Name
	}
	limit := min(len(articles), 24)
	references := make([]ai.ArticleReference, 0, limit)
	for _, article := range articles[:limit] {
		references = append(references, ai.ArticleReference{
			Title:    article.Title,
			Slug:     article.Slug,
			Excerpt:  article.Excerpt,
			Category: categoryNamesByID[article.CategoryID],
			Status:   article.Status,
		})
	}
	return references
}

func (a *API) storeGeneratedArticleImage(ctx context.Context, siteID string, draft ai.ArticleDraft) (mediaAssetResponse, error) {
	if draft.Image == nil || len(draft.Image.Contents) == 0 || len(draft.Image.Contents) > maxGeneratedImageBytes {
		return mediaAssetResponse{}, errors.New("generated image is invalid")
	}
	mimeType := strings.ToLower(strings.TrimSpace(draft.Image.MimeType))
	if mimeType != "image/png" && mimeType != "image/jpeg" && mimeType != "image/webp" {
		return mediaAssetResponse{}, errors.New("generated image type is invalid")
	}
	if detected := http.DetectContentType(draft.Image.Contents); detected != mimeType {
		return mediaAssetResponse{}, errors.New("generated image contents do not match its type")
	}
	return a.uploadMediaAsset(ctx, siteID, "ai-"+slugify(draft.Title)+imageFileExtension(mimeType), draft.Image.Contents, mimeType, draft.Image.AltText)
}

func imageFileExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}

func validateAISuggestionRequest(payload aiSuggestionRequest) error {
	payload.Instruction = strings.TrimSpace(payload.Instruction)
	if len(payload.Instruction) < 3 {
		return fmt.Errorf("%w: enter an instruction of at least 3 characters", errValidation)
	}
	if len(payload.Instruction) > 4000 || len(payload.Title) > 12000 || len(payload.Excerpt) > 12000 || len(payload.ContentMarkdown) > 64000 {
		return fmt.Errorf("%w: AI suggestion context is too large", errValidation)
	}
	return nil
}

func (a *API) writeAISuggestionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ai.ErrNotConfigured):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "AI provider is not configured. Choose a provider, model, and API key environment variable in AI settings."})
	case errors.Is(err, ai.ErrUnsupportedDraft):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Full AI article generation currently requires Google Gemini. Choose Google Gemini in AI settings."})
	case errors.Is(err, ai.ErrSecretUnavailable):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "The configured AI API key environment variable is unavailable on the server."})
	case errors.Is(err, ai.ErrProviderResponse), errors.Is(err, ai.ErrEmptyResponse):
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "The AI provider could not generate a suggestion. Check the provider settings and try again."})
	case errors.Is(err, context.DeadlineExceeded):
		writeJSON(w, http.StatusGatewayTimeout, map[string]string{"error": "The AI provider took too long to respond. Try again."})
	default:
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "Unable to reach the AI provider. Try again."})
	}
}
