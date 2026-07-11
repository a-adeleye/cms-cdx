package handlers

import (
	"context"
	"strings"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/services"
)

func TestCreateTemplateRegistersTemplateInWorkspace(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	buildRoot := t.TempDir()
	deployRoot := t.TempDir()
	api := &API{Services: services.Services{DB: db, Builder: builder.NewLocalBuilder(buildRoot), Deploy: deploy.NewFilesystemAdapter(deployRoot)}}

	created, err := api.createTemplate(ctx, templateUpsertRequest{Name: "Magazine", Slug: "magazine"})
	if err != nil {
		t.Fatalf("createTemplate returned error: %v", err)
	}
	if created.Slug != "magazine" {
		t.Fatalf("expected slug magazine, got %q", created.Slug)
	}

	workspace, err := api.loadWorkspace(ctx, "", "")
	if err != nil {
		t.Fatalf("loadWorkspace returned error: %v", err)
	}

	found := false
	for _, template := range workspace.Templates {
		if template.Slug == "magazine" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected workspace templates to include magazine, got %#v", workspace.Templates)
	}
	found = false
	for _, template := range workspace.Templates {
		if template.Slug == "anonime" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected workspace templates to include anonime, got %#v", workspace.Templates)
	}

	_, _ = db.ExecContext(ctx, `DELETE FROM templates WHERE slug = $1`, "magazine")
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
