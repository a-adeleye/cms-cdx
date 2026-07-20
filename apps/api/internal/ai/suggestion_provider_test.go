package ai

import (
	"context"
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
