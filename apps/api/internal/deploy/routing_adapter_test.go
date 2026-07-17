package deploy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/models"
)

func TestRoutingAdapterRejectsUnconfiguredCloudflarePagesDeployment(t *testing.T) {
	outputPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(outputPath, "index.html"), []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write deployment output: %v", err)
	}

	_, err := NewAdapter(t.TempDir()).Deploy(t.Context(), models.Site{
		Slug:           "example",
		DeployProvider: "cloudflare_pages",
		DeployConfig: map[string]any{
			"projectName": "example-site",
		},
	}, models.Build{BuildType: "published"}, outputPath)
	if err == nil {
		t.Fatal("expected an unconfigured Cloudflare Pages deployment to be rejected")
	}
	if !strings.Contains(err.Error(), "CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID") {
		t.Fatalf("expected missing Cloudflare credentials error, got %v", err)
	}
}

func TestCloudflarePagesAdapterUsesPreviewDeploymentSettings(t *testing.T) {
	_, err := (CloudflarePagesAdapter{}).Deploy(t.Context(), models.Site{
		DeployProvider:        "firebase",
		PreviewDeployProvider: "cloudflare_pages",
		PreviewDeployConfig: map[string]any{
			"projectName": "example-site",
		},
	}, models.Build{BuildType: "preview"}, "")
	if err == nil {
		t.Fatal("expected missing Cloudflare credentials error for a preview deployment")
	}
	if !strings.Contains(err.Error(), "CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID") {
		t.Fatalf("expected missing Cloudflare credentials error, got %v", err)
	}
}
