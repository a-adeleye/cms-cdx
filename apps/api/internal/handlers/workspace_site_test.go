package handlers

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/services"
)

func TestUpdateSitePersistsPreviewDeploySettingsAndPreviewBuildUsesThem(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	original, err := api.getSite(ctx, siteID)
	if err != nil {
		t.Fatalf("getSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = api.updateSite(ctx, siteID, siteUpsertRequest{
			Name:                  original.Name,
			Slug:                  original.Slug,
			Domain:                original.Domain,
			BlogPath:              original.BlogPath,
			Status:                original.Status,
			TemplateKey:           original.TemplateKey,
			ThemeConfig:           original.ThemeConfig,
			DeployProvider:        original.DeployProvider,
			DeployConfig:          original.DeployConfig,
			PreviewDeployProvider: original.PreviewDeployProvider,
			PreviewDeployConfig:   original.PreviewDeployConfig,
			AIConfig:              original.AIConfig,
			StorageConfig:         original.StorageConfig,
		})
		_ = db.Close()
	})

	updated, err := api.updateSite(ctx, siteID, siteUpsertRequest{
		Name:                  original.Name,
		Slug:                  original.Slug,
		Domain:                original.Domain,
		BlogPath:              original.BlogPath,
		Status:                original.Status,
		TemplateKey:           original.TemplateKey,
		ThemeConfig:           original.ThemeConfig,
		DeployProvider:        "netlify",
		DeployConfig:          `{"production":"true"}`,
		PreviewDeployProvider: "cloudflare",
		PreviewDeployConfig:   `{"branch":"preview"}`,
		AIConfig:              original.AIConfig,
		StorageConfig:         original.StorageConfig,
	})
	if err != nil {
		t.Fatalf("updateSite returned error: %v", err)
	}
	if updated.PreviewDeployProvider != "cloudflare" {
		t.Fatalf("expected preview deploy provider cloudflare, got %q", updated.PreviewDeployProvider)
	}
	if updated.PreviewDeployConfig != `{"branch":"preview"}` {
		t.Fatalf("expected preview deploy config to persist, got %q", updated.PreviewDeployConfig)
	}

	build, err := api.createBuild(ctx, siteID, buildCreateRequest{BuildType: "preview"})
	if err != nil {
		t.Fatalf("createBuild returned error: %v", err)
	}
	if build.BuildType != "preview" {
		t.Fatalf("expected preview build type, got %q", build.BuildType)
	}
	if !strings.HasSuffix(build.OutputPath, filepath.Join("preview", "site")) {
		t.Fatalf("expected preview output path to end with preview/site, got %q", build.OutputPath)
	}
	if build.DeployProvider != "cloudflare" {
		t.Fatalf("expected preview deploy provider cloudflare, got %q", build.DeployProvider)
	}
	if build.DeployStatus != "deployed" {
		t.Fatalf("expected deployed status, got %q", build.DeployStatus)
	}
	if build.DeployURL == "" || !strings.Contains(build.DeployURL, "/deployments/") {
		t.Fatalf("expected deployed URL to point to local deployments, got %q", build.DeployURL)
	}
	if _, err := os.Stat(filepath.Join(build.OutputPath, "index.html")); err != nil {
		t.Fatalf("expected generated index.html: %v", err)
	}
	if _, err := os.Stat(filepath.Join(deployRoot, "cloudflare", "example", "preview", "index.html")); err != nil {
		t.Fatalf("expected deployed index.html: %v", err)
	}
}

func TestGetSiteIncludesDeploymentWarningsForMissingFirebaseSecret(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	original, err := api.getSite(ctx, siteID)
	if err != nil {
		t.Fatalf("getSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = api.updateSite(ctx, siteID, siteUpsertRequest{
			Name:                  original.Name,
			Slug:                  original.Slug,
			Domain:                original.Domain,
			BlogPath:              original.BlogPath,
			Status:                original.Status,
			TemplateKey:           original.TemplateKey,
			ThemeConfig:           original.ThemeConfig,
			DeployProvider:        original.DeployProvider,
			DeployConfig:          original.DeployConfig,
			PreviewDeployProvider: original.PreviewDeployProvider,
			PreviewDeployConfig:   original.PreviewDeployConfig,
			AIConfig:              original.AIConfig,
			StorageConfig:         original.StorageConfig,
		})
		_ = db.Close()
	})

	updated, err := api.updateSite(ctx, siteID, siteUpsertRequest{
		Name:                  original.Name,
		Slug:                  original.Slug,
		Domain:                original.Domain,
		BlogPath:              original.BlogPath,
		Status:                original.Status,
		TemplateKey:           original.TemplateKey,
		ThemeConfig:           original.ThemeConfig,
		DeployProvider:        "firebase",
		DeployConfig:          `{"provider":"firebase","projectId":"demo-project","siteId":"demo-site","serviceAccountSecretRef":"MISSING_FIREBASE_SERVICE_ACCOUNT_JSON_FOR_TEST"}`,
		PreviewDeployProvider: "firebase",
		PreviewDeployConfig:   `{"provider":"firebase","projectId":"demo-project","siteId":"demo-site-preview","serviceAccountSecretRef":"MISSING_FIREBASE_SERVICE_ACCOUNT_JSON_FOR_TEST"}`,
		AIConfig:              original.AIConfig,
		StorageConfig:         original.StorageConfig,
	})
	if err != nil {
		t.Fatalf("updateSite returned error: %v", err)
	}
	if len(updated.DeploymentWarnings) != 2 {
		t.Fatalf("expected production and preview deployment warnings, got %d: %#v", len(updated.DeploymentWarnings), updated.DeploymentWarnings)
	}
	if !strings.Contains(updated.DeploymentWarnings[0], "MISSING_FIREBASE_SERVICE_ACCOUNT_JSON_FOR_TEST") && !strings.Contains(updated.DeploymentWarnings[1], "MISSING_FIREBASE_SERVICE_ACCOUNT_JSON_FOR_TEST") {
		t.Fatalf("expected deployment warnings to mention the missing firebase secret ref, got %#v", updated.DeploymentWarnings)
	}
}
