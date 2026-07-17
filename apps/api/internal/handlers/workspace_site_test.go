package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/models"
	"cms-builder/api/internal/services"
)

type failingDeployAdapter struct{}

func (failingDeployAdapter) Deploy(context.Context, models.Site, models.Build, string) (*deploy.DeployResult, error) {
	return nil, errors.New("simulated provider failure")
}

func TestCreateBuildRecordsDeploymentFailure(t *testing.T) {
	db := openTestDatabase(t)
	defer db.Close()
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(t.TempDir()), Deploy: failingDeployAdapter{}}}

	if _, err := api.createBuild(ctx, siteID, buildCreateRequest{BuildType: "preview"}); err == nil {
		t.Fatal("expected simulated deployment failure")
	}
	var status, deployStatus, logs string
	if err := db.QueryRowContext(ctx, `SELECT status, deploy_status, logs FROM builds WHERE site_id = $1 ORDER BY created_at DESC LIMIT 1`, siteID).Scan(&status, &deployStatus, &logs); err != nil {
		t.Fatal(err)
	}
	if status != "failed" || deployStatus != "failed" || !strings.Contains(logs, "simulated provider failure") {
		t.Fatalf("expected persisted failure, got status=%q deploy=%q logs=%q", status, deployStatus, logs)
	}
}

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
		DeployProvider:        "none",
		DeployConfig:          `{}`,
		PreviewDeployProvider: "cloudflare_pages",
		PreviewDeployConfig:   `{"projectName":"preview-project","productionBranch":"preview"}`,
		AIConfig:              original.AIConfig,
		StorageConfig:         original.StorageConfig,
	})
	if err != nil {
		t.Fatalf("updateSite returned error: %v", err)
	}
	if updated.PreviewDeployProvider != "cloudflare_pages" {
		t.Fatalf("expected preview deploy provider cloudflare_pages, got %q", updated.PreviewDeployProvider)
	}
	var previewConfig map[string]string
	if err := json.Unmarshal([]byte(updated.PreviewDeployConfig), &previewConfig); err != nil {
		t.Fatalf("expected valid preview deploy config JSON, got %q: %v", updated.PreviewDeployConfig, err)
	}
	if previewConfig["productionBranch"] != "preview" {
		t.Fatalf("expected preview deploy branch to persist, got %q", updated.PreviewDeployConfig)
	}

	build, err := api.createBuild(ctx, siteID, buildCreateRequest{BuildType: "preview"})
	if err != nil {
		t.Fatalf("createBuild returned error: %v", err)
	}
	if build.BuildType != "preview" {
		t.Fatalf("expected preview build type, got %q", build.BuildType)
	}
	if !strings.HasSuffix(build.OutputPath, filepath.Join("preview", original.Slug)) {
		t.Fatalf("expected preview output path to be isolated by site, got %q", build.OutputPath)
	}
	if build.DeployProvider != "cloudflare_pages" {
		t.Fatalf("expected preview deploy provider cloudflare_pages, got %q", build.DeployProvider)
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
	if _, err := os.Stat(filepath.Join(deployRoot, "cloudflare_pages", original.Slug, "preview", "index.html")); err != nil {
		t.Fatalf("expected deployed index.html (logs: %s): %v", build.Logs, err)
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

func TestUpdateSitePersistsContentContext(t *testing.T) {
	db := openTestDatabase(t)
	defer db.Close()
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}
	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	original, err := api.getSite(ctx, siteID)
	if err != nil {
		t.Fatalf("getSite returned error: %v", err)
	}
	t.Cleanup(func() {
		_, _ = api.updateSite(context.Background(), siteID, siteUpsertRequest{
			Name: original.Name, Slug: original.Slug, Domain: original.Domain, BlogPath: original.BlogPath,
			Description: original.Description, ContentContext: original.ContentContext, Status: original.Status,
			TemplateKey: original.TemplateKey, ThemeConfig: original.ThemeConfig, DeployProvider: original.DeployProvider,
			DeployConfig: original.DeployConfig, PreviewDeployProvider: original.PreviewDeployProvider,
			PreviewDeployConfig: original.PreviewDeployConfig, AIConfig: original.AIConfig, StorageConfig: original.StorageConfig,
		})
	})

	updated, err := api.updateSite(ctx, siteID, siteUpsertRequest{
		Name: original.Name, Slug: original.Slug, Domain: original.Domain, BlogPath: original.BlogPath,
		Description: original.Description, ContentContext: "application_blog", Status: original.Status,
		TemplateKey: original.TemplateKey, ThemeConfig: original.ThemeConfig, DeployProvider: original.DeployProvider,
		DeployConfig: original.DeployConfig, PreviewDeployProvider: original.PreviewDeployProvider,
		PreviewDeployConfig: original.PreviewDeployConfig, AIConfig: original.AIConfig, StorageConfig: original.StorageConfig,
	})
	if err != nil {
		t.Fatalf("updateSite returned error: %v", err)
	}
	if updated.ContentContext != "application_blog" {
		t.Fatalf("expected application_blog context, got %q", updated.ContentContext)
	}
}
