package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cms-builder/api/internal/ai"
	"cms-builder/api/internal/services"
)

type fakeSuggestionGenerator struct {
	config      ai.Config
	input       ai.SuggestionInput
	draftInput  ai.GenerateArticleDraftInput
	draftResult ai.ArticleDraft
}

func (generator *fakeSuggestionGenerator) GenerateSuggestion(_ context.Context, config ai.Config, input ai.SuggestionInput) (ai.Suggestion, error) {
	generator.config = config
	generator.input = input
	return ai.Suggestion{Content: "## Suggested copy\n\nA clearer introduction.", Model: config.Model}, nil
}

func (generator *fakeSuggestionGenerator) GenerateArticleDraft(_ context.Context, config ai.Config, input ai.GenerateArticleDraftInput) (ai.ArticleDraft, error) {
	generator.config = config
	generator.draftInput = input
	if generator.draftResult.Title != "" {
		return generator.draftResult, nil
	}
	return ai.ArticleDraft{
		Title:           "A practical privacy guide",
		Slug:            "practical-privacy-guide",
		Category:        input.Categories[0],
		Tags:            []string{"privacy", "security"},
		SEOTitle:        "A practical privacy guide",
		MetaDescription: "Practical steps for protecting your privacy online.",
		Excerpt:         "A concise guide to stronger online privacy.",
		Content:         "# A practical privacy guide\n\nUseful privacy guidance.",
		Model:           config.Model,
	}, nil
}

func TestSiteAISuggestionUsesStoredProviderConfig(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	originalConfig := mustQueryText(t, db, ctx, `SELECT COALESCE(ai_config::text, '{}') FROM sites WHERE id = $1`, siteID)
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `UPDATE sites SET ai_config = $2::jsonb WHERE id = $1`, siteID, originalConfig)
	})
	if _, err := db.ExecContext(ctx, `UPDATE sites SET ai_config = $2::jsonb WHERE id = $1`, siteID, `{"provider":"openai","model":"gpt-4.1-mini","apiKeySecretRef":"OPENAI_API_KEY"}`); err != nil {
		t.Fatalf("configure test site: %v", err)
	}

	generator := &fakeSuggestionGenerator{}
	api := &API{Services: services.Services{DB: db, AI: generator}}
	body, err := json.Marshal(aiSuggestionRequest{
		Instruction:     "Make the introduction clearer",
		Title:           "Privacy basics",
		Excerpt:         "A short primer.",
		ContentMarkdown: "Original copy.",
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+siteID+"/ai/suggestions", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	api.siteSubroutes(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if generator.config.Provider != "openai" || generator.config.Model != "gpt-4.1-mini" {
		t.Fatalf("generator config = %#v", generator.config)
	}
	if generator.input.Instruction != "Make the introduction clearer" {
		t.Fatalf("generator instruction = %q", generator.input.Instruction)
	}
	var response aiSuggestionResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Suggestion == "" || response.Model != "gpt-4.1-mini" {
		t.Fatalf("response = %#v", response)
	}
}

func TestAISuggestionRejectsEmptyInstructionBeforeCallingProvider(t *testing.T) {
	if err := validateAISuggestionRequest(aiSuggestionRequest{Instruction: "  "}); err == nil {
		t.Fatal("validateAISuggestionRequest() error = nil, want validation error")
	}
}

func TestAIArticleDraftUsesMasterPromptAndExistingArticleReferences(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	originalConfig := mustQueryText(t, db, ctx, `SELECT COALESCE(ai_config::text, '{}') FROM sites WHERE id = $1`, siteID)
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `UPDATE sites SET ai_config = $2::jsonb WHERE id = $1`, siteID, originalConfig)
	})
	if _, err := db.ExecContext(ctx, `UPDATE sites SET ai_config = $2::jsonb WHERE id = $1`, siteID, `{"provider":"google","model":"gemini-2.5-flash","apiKeySecretRef":"GEMINI_API_KEY","masterPrompt":"Write accurate, useful Anonime articles."}`); err != nil {
		t.Fatalf("configure test site: %v", err)
	}

	generator := &fakeSuggestionGenerator{}
	api := &API{Services: services.Services{DB: db, AI: generator}}
	body, err := json.Marshal(aiArticleDraftRequest{Topic: "How email aliases protect privacy"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+siteID+"/ai/article-drafts", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	api.siteSubroutes(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if generator.draftInput.MasterPrompt != "Write accurate, useful Anonime articles." {
		t.Fatalf("master prompt = %q", generator.draftInput.MasterPrompt)
	}
	if generator.draftInput.Topic != "How email aliases protect privacy" {
		t.Fatalf("topic = %q", generator.draftInput.Topic)
	}
	if len(generator.draftInput.ExistingArticles) == 0 {
		t.Fatal("expected existing articles to be included in draft context")
	}
	var response aiArticleDraftResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Title != "A practical privacy guide" || response.ContentMarkdown == "" || response.Model != "gemini-2.5-flash" {
		t.Fatalf("response = %#v", response)
	}
}

func TestAIArticleDraftAllowsAnEmptyTopic(t *testing.T) {
	if err := validateAIArticleDraftRequest(aiArticleDraftRequest{Topic: "  "}); err != nil {
		t.Fatalf("validateAIArticleDraftRequest() error = %v, want nil", err)
	}
}
