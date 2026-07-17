package deploy

import (
	"context"
	"strings"

	"cms-builder/api/internal/models"
)

type DeployResult struct {
	Provider string `json:"provider"`
	URL      string `json:"url,omitempty"`
	Message  string `json:"message,omitempty"`
	Revision string `json:"revision,omitempty"`
}

type DeployAdapter interface {
	Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error)
}

type NoopAdapter struct{}

func (NoopAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	return &DeployResult{
		Provider: fallbackProvider(site.DeployProvider),
		Message:  "deployment skipped in scaffold",
	}, nil
}

func fallbackProvider(value string) string {
	if value == "" {
		return "none"
	}
	return value
}

func configString(values map[string]any, key string) string {
	if values == nil {
		return ""
	}
	if value, ok := values[key].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}
