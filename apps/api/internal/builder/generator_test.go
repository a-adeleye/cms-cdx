package builder

import (
	"context"
	"path/filepath"
	"testing"

	"cms-builder/api/internal/models"
)

func TestBuildOutputPathUsesPreviewDirectoryForPreviewBuilds(t *testing.T) {
	site := models.Site{Slug: "Anonime Blog"}

	if got := buildOutputPath("dist/builds", site, false); got != filepath.Join("dist", "builds", "sites", "anonime-blog") {
		t.Fatalf("expected published build path to use the site slug, got %q", got)
	}

	if got := buildOutputPath("dist/builds", site, true); got != filepath.Join("dist", "builds", "preview", "anonime-blog") {
		t.Fatalf("expected preview build path to use the preview directory, got %q", got)
	}
}

func TestThemeForSiteUsesAnonimePalette(t *testing.T) {
	theme := themeForSite(models.Site{TemplateKey: "anonime", Name: "Anonime"})

	if theme.Name != "anonime" {
		t.Fatalf("expected anonime theme, got %q", theme.Name)
	}
	if theme.LayoutClass != "anonime-layout" {
		t.Fatalf("expected anonime layout class, got %q", theme.LayoutClass)
	}
	if theme.Accent != "#10b26c" {
		t.Fatalf("expected anonime accent color, got %q", theme.Accent)
	}
}

func TestNoopBuilderUsesPublishedAndPreviewDestinations(t *testing.T) {
	builder := NoopBuilder{}
	content := SiteContent{Site: models.Site{Slug: "anonime-blog"}}

	publishedPath, err := builder.GenerateSite(context.Background(), content, GenerateOptions{})
	if err != nil {
		t.Fatalf("GenerateSite returned error for published build: %v", err)
	}
	if publishedPath != "dist/sites/anonime-blog" {
		t.Fatalf("expected published build path, got %q", publishedPath)
	}

	previewPath, err := builder.GenerateSite(context.Background(), content, GenerateOptions{Preview: true})
	if err != nil {
		t.Fatalf("GenerateSite returned error for preview build: %v", err)
	}
	if previewPath != "dist/preview/anonime-blog" {
		t.Fatalf("expected preview build path, got %q", previewPath)
	}
}
