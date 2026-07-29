package builder

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/models"
)

func TestAstroBuilderCommandHelper(t *testing.T) {
	if !isAstroBuilderHelperProcess() {
		return
	}
	dataPath := os.Getenv("CMS_BUILD_DATA_FILE")
	outputPath := os.Getenv("CMS_BUILD_OUTPUT_DIR")
	if dataPath == "" || outputPath == "" {
		t.Fatal("expected Astro build paths")
	}
	contents, err := os.ReadFile(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), `"coverImageUrl":"`) {
		t.Fatalf("expected CMS data to use the Astro cover image contract, got %s", contents)
	}
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputPath, "index.html"), []byte("<h1>Supromail</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(outputPath, "articles"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputPath, "articles", "index.html"), []byte("<h1>Articles</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputPath, "articles", "sitemap.xml"), []byte("<urlset></urlset>"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAstroBuilderGeneratesSupromailTemplatePreview(t *testing.T) {
	if isAstroBuilderHelperProcess() {
		return
	}
	builderDirectory := t.TempDir()
	if err := os.MkdirAll(filepath.Join(builderDirectory, "templates", SupromailTemplateKey), 0o755); err != nil {
		t.Fatal(err)
	}
	outputRoot := t.TempDir()
	builder := AstroBuilder{OutputRoot: outputRoot, BuilderDirectory: builderDirectory, NPMCommand: os.Args[0]}

	outputPath, err := builder.GenerateTemplatePreview(context.Background(), SupromailTemplateKey)
	if err != nil {
		t.Fatalf("GenerateTemplatePreview returned error: %v", err)
	}
	if outputPath != filepath.Join(outputRoot, "template-previews", SupromailTemplateKey) {
		t.Fatalf("unexpected preview output path: %q", outputPath)
	}
	if _, err := os.Stat(filepath.Join(outputPath, "index.html")); err != nil {
		t.Fatalf("expected preview output: %v", err)
	}
}

func TestAstroBuilderGeneratesSupromailOutputFromCMSData(t *testing.T) {
	if isAstroBuilderHelperProcess() {
		return
	}
	builderDirectory := t.TempDir()
	if err := os.MkdirAll(filepath.Join(builderDirectory, "templates", SupromailTemplateKey), 0o755); err != nil {
		t.Fatal(err)
	}
	outputRoot := t.TempDir()
	builder := AstroBuilder{
		OutputRoot:       outputRoot,
		BuilderDirectory: builderDirectory,
		NPMCommand:       os.Args[0],
	}

	outputPath, err := builder.GenerateSite(context.Background(), SiteContent{
		Site: models.Site{Name: "Supromail", Slug: "supromail-test", BlogPath: "/blog", TemplateKey: SupromailTemplateKey},
		Articles: []ArticleContent{{
			ID: "article-1", Title: "Article", Slug: "article", ContentMarkdown: "# Article", CoverImageURL: "https://cdn.example/cover.png",
		}},
	}, GenerateOptions{Preview: true})
	if err != nil {
		t.Fatalf("GenerateSite returned error: %v", err)
	}
	if outputPath != filepath.Join(outputRoot, "preview", "supromail-test") {
		t.Fatalf("unexpected output path: %q", outputPath)
	}
	contents, err := os.ReadFile(filepath.Join(outputPath, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "<h1>Supromail</h1>" {
		t.Fatalf("unexpected generated output: %q", contents)
	}
	if _, err := os.Stat(filepath.Join(outputPath, "blog", "index.html")); err != nil {
		t.Fatalf("expected CMS blog path to receive Astro article routes: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputPath, "blog", "sitemap.xml")); err != nil {
		t.Fatalf("expected blog sitemap to be relocated with Astro article routes: %v", err)
	}
}

func isAstroBuilderHelperProcess() bool {
	return len(os.Args) >= 3 && os.Args[len(os.Args)-2] == "run" && os.Args[len(os.Args)-1] == "build"
}

func TestNewAstroBuildDataUsesTemplateContractAndPreviewSelection(t *testing.T) {
	data, err := newAstroBuildData(SiteContent{
		Site: models.Site{Name: "Supromail", Domain: "https://supromail.example", BlogPath: "/blog"},
		Articles: []ArticleContent{
			{ID: "included", Title: "Included", Slug: "included", CoverImageURL: "https://cdn.example/included.png", IsFeatured: true},
			{ID: "excluded", Title: "Excluded", Slug: "excluded"},
		},
	}, "/blog", GenerateOptions{Preview: true, ArticleIDs: []string{"included"}})
	if err != nil {
		t.Fatalf("newAstroBuildData returned error: %v", err)
	}
	if len(data.Articles) != 1 || data.Articles[0].Title != "Included" {
		t.Fatalf("expected only selected article, got %#v", data.Articles)
	}
	if data.Articles[0].CoverImageURL != "https://cdn.example/included.png" {
		t.Fatalf("expected Astro cover image field, got %#v", data.Articles[0])
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"blogPath":"/blog"`) {
		t.Fatalf("expected camel-case Astro site contract, got %s", encoded)
	}
}

func TestAstroBuildEnvironmentExcludesProcessSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "not-for-node")
	environment := astroBuildEnvironment("data.json", "output", "https://example.com", "template")
	for _, value := range environment {
		if value == "JWT_SECRET=not-for-node" {
			t.Fatal("expected Node build environment to exclude API secrets")
		}
	}
}
