package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func TestProviderGenerateSuggestionUsesConfiguredOpenAIProvider(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.String() != openAIResponsesURL {
			t.Fatalf("request URL = %q, want %q", request.URL.String(), openAIResponsesURL)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer test-secret" {
			t.Fatalf("Authorization = %q", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"output_text":"## Revised section\n\nA useful revision."}`)),
		}, nil
	})}
	provider := NewProvider(client, func(name string) (string, bool) {
		return "test-secret", name == "OPENAI_API_KEY"
	})

	suggestion, err := provider.GenerateSuggestion(context.Background(), Config{
		Provider: "openai", Model: "gpt-4.1-mini", APIKeySecretRef: "OPENAI_API_KEY",
	}, SuggestionInput{Instruction: "Improve the introduction", Title: "Privacy", ContentMarkdown: "Original copy."})
	if err != nil {
		t.Fatalf("GenerateSuggestion() error = %v", err)
	}
	if suggestion.Content != "## Revised section\n\nA useful revision." {
		t.Fatalf("suggestion.Content = %q", suggestion.Content)
	}
	if suggestion.Model != "gpt-4.1-mini" {
		t.Fatalf("suggestion.Model = %q", suggestion.Model)
	}
}

func TestProviderGenerateSuggestionRejectsMissingServerSecret(t *testing.T) {
	provider := NewProvider(nil, func(string) (string, bool) { return "", false })
	_, err := provider.GenerateSuggestion(context.Background(), Config{
		Provider: "openai", Model: "gpt-4.1-mini", APIKeySecretRef: "OPENAI_API_KEY",
	}, SuggestionInput{Instruction: "Improve this"})
	if !errors.Is(err, ErrSecretUnavailable) {
		t.Fatalf("GenerateSuggestion() error = %v, want ErrSecretUnavailable", err)
	}
}

func TestParseConfigRejectsUnsafeCompatibleBaseURL(t *testing.T) {
	_, err := ParseConfig(`{"provider":"openai_compatible","baseUrl":"http://example.test","model":"demo","apiKeySecretRef":"AI_KEY"}`)
	if err == nil {
		t.Fatal("ParseConfig() error = nil, want error for an insecure base URL")
	}
}

func TestProviderGeneratesStructuredGeminiArticleAndFeaturedImage(t *testing.T) {
	requestCount := 0
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requestCount++
		if got := request.Header.Get("x-goog-api-key"); got != "test-secret" {
			t.Fatalf("x-goog-api-key = %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode provider payload: %v", err)
		}
		switch requestCount {
		case 1:
			if request.URL.String() != googleGenerateURL+"gemini-2.5-flash:generateContent" {
				t.Fatalf("text request URL = %q", request.URL.String())
			}
			if config, ok := payload["generationConfig"].(map[string]any); !ok || config["responseMimeType"] != "application/json" {
				t.Fatalf("text request did not request JSON: %#v", payload)
			}
			articleJSON := `{"title":"Email aliases explained","slug":"email-aliases-explained","category":"Privacy Basics","featured":false,"tags":["email privacy","aliases"],"seo":{"seoTitle":"Email aliases explained","metaDescription":"Learn how email aliases protect your inbox.","canonical":"https://anonime.io/blog/email-aliases-explained"},"excerpt":"A useful privacy guide.","image":{"alt":"Abstract private email routing","caption":"A privacy-first inbox concept."},"content":"# Email aliases explained\n\nUseful guidance."}`
			response, err := json.Marshal(map[string]any{"candidates": []any{map[string]any{"content": map[string]any{"parts": []any{map[string]string{"text": articleJSON}}}}}})
			if err != nil {
				t.Fatalf("marshal provider response: %v", err)
			}
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(response)))}, nil
		case 2:
			if request.URL.String() != googleGenerateURL+"gemini-2.5-flash-image:generateContent" {
				t.Fatalf("image request URL = %q", request.URL.String())
			}
			encodedPayload, err := json.Marshal(payload)
			if err != nil || !strings.Contains(string(encodedPayload), "Example Site") {
				t.Fatalf("image request did not use the site's identity: %s", encodedPayload)
			}
			response := `{"candidates":[{"content":{"parts":[{"inlineData":{"mimeType":"image/png","data":"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAusB9Y9J1n0AAAAASUVORK5CYII="}}]}}]}`
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}, nil
		default:
			t.Fatalf("unexpected provider request %d", requestCount)
			return nil, nil
		}
	})}
	provider := NewProvider(client, func(name string) (string, bool) { return "test-secret", name == "GEMINI_API_KEY" })

	draft, err := provider.GenerateArticleDraft(context.Background(), Config{
		Provider: "google", Model: "gemini-2.5-flash", APIKeySecretRef: "GEMINI_API_KEY",
	}, GenerateArticleDraftInput{
		Topic: "How email aliases work", SiteName: "Example Site", MasterPrompt: "Write a publication-ready privacy article.", Categories: []string{"Privacy Basics"},
		ExistingArticles: []ArticleReference{{Title: "Protect your inbox", Slug: "protect-inbox", Excerpt: "Existing coverage."}},
	})
	if err != nil {
		t.Fatalf("GenerateArticleDraft() error = %v", err)
	}
	if requestCount != 2 || draft.Title != "Email aliases explained" || draft.Image == nil || draft.Image.MimeType != "image/png" {
		t.Fatalf("draft = %#v, requestCount = %d", draft, requestCount)
	}
	if draft.Image.AltText != "Abstract private email routing" || len(draft.Image.Contents) == 0 {
		t.Fatalf("image = %#v", draft.Image)
	}
}

func TestArticleDraftPromptSelectsATopicWhenNoBriefIsSupplied(t *testing.T) {
	prompt := articleDraftPrompt(GenerateArticleDraftInput{
		MasterPrompt: "Write useful articles.",
		Categories:   []string{"Guides"},
	})
	if !strings.Contains(prompt, "Choose a high-value, non-duplicative topic") {
		t.Fatalf("empty-topic draft prompt = %q", prompt)
	}
}
