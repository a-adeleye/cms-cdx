package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	openAIResponsesURL   = "https://api.openai.com/v1/responses"
	anthropicMessagesURL = "https://api.anthropic.com/v1/messages"
	googleGenerateURL    = "https://generativelanguage.googleapis.com/v1beta/models/"
)

var (
	ErrNotConfigured     = errors.New("AI provider is not configured")
	ErrSecretUnavailable = errors.New("AI provider secret is unavailable")
	ErrProviderResponse  = errors.New("AI provider returned an unsuccessful response")
	ErrEmptyResponse     = errors.New("AI provider returned an empty suggestion")
	secretRefPattern     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,127}$`)
)

// Config is the non-secret portion of a site's AI configuration. API keys are
// resolved only from the named server-side environment variable.
type Config struct {
	Provider        string `json:"provider"`
	Model           string `json:"model"`
	APIKeySecretRef string `json:"apiKeySecretRef"`
	BaseURL         string `json:"baseUrl"`
}

// ParseConfig validates the stored site configuration before it can influence
// an outbound request. An empty configuration is valid but cannot generate.
func ParseConfig(raw string) (Config, error) {
	if strings.TrimSpace(raw) == "" {
		raw = "{}"
	}

	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("AI config must be JSON with supported fields: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return Config{}, errors.New("AI config must contain one JSON object")
	}

	config.Provider = strings.TrimSpace(config.Provider)
	config.Model = strings.TrimSpace(config.Model)
	config.APIKeySecretRef = strings.TrimSpace(config.APIKeySecretRef)
	config.BaseURL = strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")

	switch config.Provider {
	case "", "none", "openai", "anthropic", "google":
	case "openai_compatible":
		parsed, err := url.Parse(config.BaseURL)
		if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Hostname() == "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			return Config{}, errors.New("compatible AI provider requires an HTTPS baseUrl without credentials, query, or fragment")
		}
	default:
		return Config{}, fmt.Errorf("unsupported AI provider %q", config.Provider)
	}

	if config.APIKeySecretRef != "" && !secretRefPattern.MatchString(config.APIKeySecretRef) {
		return Config{}, errors.New("apiKeySecretRef must be an environment variable name")
	}
	if len(config.Model) > 256 || len(config.BaseURL) > 2048 {
		return Config{}, errors.New("AI configuration value is too long")
	}
	return config, nil
}

type SuggestionInput struct {
	Instruction     string
	Title           string
	Excerpt         string
	ContentMarkdown string
}

type Suggestion struct {
	Content string
	Model   string
}

type SuggestionGenerator interface {
	GenerateSuggestion(ctx context.Context, config Config, input SuggestionInput) (Suggestion, error)
}

type Provider struct {
	client    *http.Client
	lookupEnv func(string) (string, bool)
}

func NewProvider(client *http.Client, lookupEnv func(string) (string, bool)) *Provider {
	if client == nil {
		client = &http.Client{
			Timeout: 45 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("AI provider redirects are not allowed")
			},
		}
	}
	if lookupEnv == nil {
		lookupEnv = os.LookupEnv
	}
	return &Provider{client: client, lookupEnv: lookupEnv}
}

func (p *Provider) GenerateSuggestion(ctx context.Context, config Config, input SuggestionInput) (Suggestion, error) {
	if config.Provider == "" || config.Provider == "none" || config.Model == "" || config.APIKeySecretRef == "" {
		return Suggestion{}, ErrNotConfigured
	}
	apiKey, ok := p.lookupEnv(config.APIKeySecretRef)
	if !ok || strings.TrimSpace(apiKey) == "" {
		return Suggestion{}, ErrSecretUnavailable
	}

	requestContext, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	prompt := editorialPrompt(input)

	var content string
	var err error
	switch config.Provider {
	case "openai":
		content, err = p.generateOpenAI(requestContext, openAIResponsesURL, config.Model, apiKey, prompt, false)
	case "openai_compatible":
		content, err = p.generateOpenAI(requestContext, config.BaseURL+"/chat/completions", config.Model, apiKey, prompt, true)
	case "anthropic":
		content, err = p.generateAnthropic(requestContext, config.Model, apiKey, prompt)
	case "google":
		content, err = p.generateGoogle(requestContext, config.Model, apiKey, prompt)
	default:
		return Suggestion{}, ErrNotConfigured
	}
	if err != nil {
		return Suggestion{}, err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return Suggestion{}, ErrEmptyResponse
	}
	return Suggestion{Content: content, Model: config.Model}, nil
}

func editorialPrompt(input SuggestionInput) string {
	return "You are an editorial writing assistant. Return only the proposed Markdown text, with no preamble, explanations, or code fences. Preserve factual claims unless the instruction explicitly asks to change them.\n\n" +
		"Instruction:\n" + strings.TrimSpace(input.Instruction) + "\n\n" +
		"Article title:\n" + strings.TrimSpace(input.Title) + "\n\n" +
		"Article excerpt:\n" + strings.TrimSpace(input.Excerpt) + "\n\n" +
		"Current article Markdown:\n" + strings.TrimSpace(input.ContentMarkdown)
}

func (p *Provider) generateOpenAI(ctx context.Context, endpoint, model, apiKey, prompt string, compatible bool) (string, error) {
	var payload any
	if compatible {
		payload = map[string]any{
			"model":    model,
			"messages": []map[string]string{{"role": "system", "content": "You are an editorial writing assistant."}, {"role": "user", "content": prompt}},
		}
	} else {
		payload = map[string]any{"model": model, "input": prompt}
	}

	body, err := p.postJSON(ctx, endpoint, map[string]string{"Authorization": "Bearer " + apiKey}, payload)
	if err != nil {
		return "", err
	}
	if compatible {
		var response struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", fmt.Errorf("invalid AI provider response: %w", err)
		}
		if len(response.Choices) == 0 {
			return "", ErrEmptyResponse
		}
		return response.Choices[0].Message.Content, nil
	}

	var response struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("invalid AI provider response: %w", err)
	}
	if strings.TrimSpace(response.OutputText) != "" {
		return response.OutputText, nil
	}
	for _, output := range response.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				return content.Text, nil
			}
		}
	}
	return "", ErrEmptyResponse
}

func (p *Provider) generateAnthropic(ctx context.Context, model, apiKey, prompt string) (string, error) {
	body, err := p.postJSON(ctx, anthropicMessagesURL, map[string]string{
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
	}, map[string]any{
		"model": model, "max_tokens": 1200,
		"system":   "You are an editorial writing assistant. Return only the proposed Markdown text.",
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	})
	if err != nil {
		return "", err
	}
	var response struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("invalid AI provider response: %w", err)
	}
	for _, content := range response.Content {
		if content.Type == "text" && strings.TrimSpace(content.Text) != "" {
			return content.Text, nil
		}
	}
	return "", ErrEmptyResponse
}

func (p *Provider) generateGoogle(ctx context.Context, model, apiKey, prompt string) (string, error) {
	endpoint := googleGenerateURL + url.PathEscape(model) + ":generateContent"
	body, err := p.postJSON(ctx, endpoint, map[string]string{"x-goog-api-key": apiKey}, map[string]any{
		"contents":          []map[string]any{{"role": "user", "parts": []map[string]string{{"text": prompt}}}},
		"systemInstruction": map[string]any{"parts": []map[string]string{{"text": "You are an editorial writing assistant. Return only the proposed Markdown text."}}},
	})
	if err != nil {
		return "", err
	}
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("invalid AI provider response: %w", err)
	}
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", ErrEmptyResponse
	}
	return response.Candidates[0].Content.Parts[0].Text, nil
}

func (p *Provider) postJSON(ctx context.Context, endpoint string, headers map[string]string, payload any) ([]byte, error) {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode AI provider request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("create AI provider request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call AI provider: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("read AI provider response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%w (status %d)", ErrProviderResponse, response.StatusCode)
	}
	return body, nil
}
