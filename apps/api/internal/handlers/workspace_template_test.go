package handlers

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/database"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/services"
)

func TestWorkspaceListsOnlyRenderedTemplatesWithPreviewURLs(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	if err := database.RunMigrations(ctx, db, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatalf("apply template migrations: %v", err)
	}
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}
	workspace, err := api.loadWorkspace(ctx, "", "")
	if err != nil {
		t.Fatalf("loadWorkspace returned error: %v", err)
	}

	foundAnonime := false
	foundSupromail := false
	for _, template := range workspace.Templates {
		if template.Slug == "anonime" {
			if template.PreviewURL != "/api/v1/template-previews/anonime" {
				t.Fatalf("expected renderer preview URL, got %q", template.PreviewURL)
			}
			foundAnonime = true
		}
		if template.Slug == builder.SupromailTemplateKey {
			if template.PreviewURL != "/api/v1/template-previews/supromail" {
				t.Fatalf("expected Supromail preview URL, got %q", template.PreviewURL)
			}
			foundSupromail = true
		}
	}
	if !foundAnonime {
		t.Fatalf("expected workspace templates to include anonime, got %#v", workspace.Templates)
	}
	if !foundSupromail {
		t.Fatalf("expected workspace templates to include Supromail, got %#v", workspace.Templates)
	}
}

func TestCreateSiteRejectsUnknownTemplateSlug(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}

	_, err := api.createSite(ctx, siteUpsertRequest{
		Name:        "Template Validation Site",
		Slug:        "template-validation-site",
		Domain:      "https://example.test",
		BlogPath:    "/articles",
		Status:      "active",
		TemplateKey: "missing-template",
	})
	if err == nil {
		t.Fatal("expected createSite to reject unknown template slug")
	}
	if !strings.Contains(err.Error(), "unknown template") {
		t.Fatalf("expected unknown template error, got %v", err)
	}
}
