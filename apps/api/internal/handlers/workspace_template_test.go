package handlers

import (
	"context"
	"strings"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/services"
)

func TestWorkspaceListsOnlyRenderedTemplatesWithPreviewURLs(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}
	workspace, err := api.loadWorkspace(ctx, "", "")
	if err != nil {
		t.Fatalf("loadWorkspace returned error: %v", err)
	}

	found := false
	for _, template := range workspace.Templates {
		if template.Slug == "anonime" {
			if template.PreviewURL != "/api/v1/template-previews/anonime" {
				t.Fatalf("expected renderer preview URL, got %q", template.PreviewURL)
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected workspace templates to include anonime, got %#v", workspace.Templates)
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
