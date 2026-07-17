package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/models"
)

func TestNewRoutesCloudflarePagesDeploymentsWithRuntimeConfiguration(t *testing.T) {
	outputPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(outputPath, "index.html"), []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write deployment output: %v", err)
	}

	services := New(nil, config.Config{
		CloudflareAPIToken:  "test-token",
		CloudflareAccountID: "test-account",
		WranglerCommand:     "command-that-does-not-exist",
	})
	_, err := services.Deploy.Deploy(t.Context(), models.Site{
		Slug:           "example-site",
		DeployProvider: "cloudflare_pages",
		DeployConfig: map[string]any{
			"projectName": "example-site",
		},
	}, models.Build{BuildType: "published"}, outputPath)
	if err == nil {
		t.Fatal("expected the configured Cloudflare Pages adapter to invoke Wrangler")
	}
	if !strings.Contains(err.Error(), "Cloudflare Pages deployment failed") {
		t.Fatalf("expected a Cloudflare Pages deployment error, got %v", err)
	}
}
