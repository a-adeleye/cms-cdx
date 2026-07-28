package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"cms-builder/api/internal/services"
)

func TestSiteExportImportPreservesEditableContentWithoutSecretsOrDeployments(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	source, err := api.createSite(ctx, siteUpsertRequest{
		Name: "Transfer source", Slug: "transfer-source", Domain: "https://source.example", BlogPath: "/blog",
		ContentContext: "standalone_blog", Status: "active", TemplateKey: "default-blog", ThemeConfig: `{"accent":"#123456"}`,
		DeployProvider: "none", DeployConfig: `{}`, PreviewDeployProvider: "none", PreviewDeployConfig: `{}`,
		AIConfig: `{"masterPrompt":"Keep the original voice.","apiKeySecretRef":"SOURCE_AI_KEY"}`, StorageConfig: `{}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM sites WHERE id = $1 OR slug LIKE 'transfer-source-imported%'`, source.ID)
	})

	author, err := api.createAuthor(ctx, source.ID, authorUpsertRequest{Name: "Transfer author", Bio: "Original author bio"})
	if err != nil {
		t.Fatal(err)
	}
	category, err := api.createCategory(ctx, source.ID, categoryUpsertRequest{Name: "Transfer category", Description: "Original category description"})
	if err != nil {
		t.Fatal(err)
	}
	asset, err := api.createMediaAsset(ctx, source.ID, mediaUpsertRequest{FileName: "transfer.png", FileURL: "https://cdn.example/transfer.png", MimeType: "image/png", SizeBytes: 42, StorageProvider: "s3", StorageKey: "source/transfer.png", AltText: "Transfer logo"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE sites SET logo_media_id = $2::uuid, favicon_media_id = $2::uuid WHERE id = $1`, source.ID, asset.ID); err != nil {
		t.Fatal(err)
	}
	article, err := api.upsertArticleWithSite(ctx, source.ID, articleUpsertRequest{Title: "Transfer article", Slug: "transfer-article", Excerpt: "Original excerpt", ContentMarkdown: "# Transfer article\n\nOriginal content.", CoverImageURL: asset.FileURL, SEOTitle: "Transfer SEO", SEODescription: "Transfer description", CanonicalURL: "https://source.example/blog/transfer-article", AuthorID: author.ID, CategoryID: category.ID, Tags: "one, two", IsFeatured: true, Status: "published"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE articles SET generated_by_ai = TRUE, human_reviewed = TRUE, ai_prompt = 'source prompt', ai_model = 'source model' WHERE id = $1`, article.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO builds (site_id, status, build_type) VALUES ($1, 'success', 'published')`, source.ID); err != nil {
		t.Fatal(err)
	}

	bundle, err := api.exportSite(ctx, source.ID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(mustJSON(t, bundle)), "SOURCE_AI_KEY") {
		t.Fatal("site export contains an AI secret reference")
	}

	imported, err := api.importSite(ctx, bundle, "")
	if err != nil {
		t.Fatal(err)
	}
	if imported.ID == source.ID || imported.Slug != "transfer-source-imported" || imported.DeployProvider != "none" || imported.PreviewDeployProvider != "none" {
		t.Fatalf("unexpected imported site: %#v", imported)
	}
	if strings.Contains(imported.AIConfig, "SOURCE_AI_KEY") || !strings.Contains(imported.AIConfig, "Keep the original voice") {
		t.Fatalf("expected imported AI configuration without secret reference, got %s", imported.AIConfig)
	}
	articles, err := api.listArticles(ctx, imported.ID)
	if err != nil || len(articles) != 1 {
		t.Fatalf("expected one imported article, got %#v, %v", articles, err)
	}
	if articles[0].ID == article.ID || articles[0].AuthorID == author.ID || articles[0].CategoryID == category.ID || articles[0].ContentMarkdown != article.ContentMarkdown || articles[0].Tags != "one, two" {
		t.Fatalf("article relationships or content were not remapped: %#v", articles[0])
	}
	media, err := api.listMediaAssets(ctx, imported.ID)
	if err != nil || len(media) != 1 || media[0].StorageProvider != "imported" || media[0].StorageKey != "" || imported.LogoMediaID != media[0].ID {
		t.Fatalf("expected imported media reference and branding map, got %#v, %#v, %v", imported, media, err)
	}
	builds, err := api.listBuilds(ctx, imported.ID)
	if err != nil || len(builds) != 0 {
		t.Fatalf("expected build history to be excluded, got %#v, %v", builds, err)
	}
}

func TestSiteImportRollsBackWhenArticleReferencesMissingTaxonomy(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}
	before := mustQueryText(t, db, ctx, `SELECT COUNT(*)::text FROM sites`)
	bundle := siteExportBundle{
		Version:  siteExportVersion,
		Site:     siteExportSite{Name: "Broken import", Slug: "broken-import", BlogPath: "/articles", ContentContext: "standalone_blog", Status: "active", TemplateKey: "default-blog", ThemeConfig: map[string]any{}, AIConfig: map[string]any{}},
		Articles: []siteExportArticle{{Title: "Broken article", Slug: "broken-article", ContentMarkdown: "# Broken", Status: "draft", AuthorSourceID: "missing-author"}},
	}
	if _, err := api.importSite(ctx, bundle, ""); err == nil {
		t.Fatal("expected import with unknown author reference to fail")
	}
	after := mustQueryText(t, db, ctx, `SELECT COUNT(*)::text FROM sites`)
	if after != before {
		t.Fatalf("failed import left a site behind: before=%s after=%s", before, after)
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	contents, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}
