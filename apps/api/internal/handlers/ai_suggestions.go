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

const maxAISuggestionRequestBytes = 96 << 10

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

// handleAISuggestion generates a proposed editorial revision for a site's
// configured provider. It never accepts API keys from the client.
func (a *API) handleAISuggestion(w http.ResponseWriter, r *http.Request, siteID string, subroutes []string) {
	if len(subroutes) != 1 || subroutes[0] != "suggestions" {
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
