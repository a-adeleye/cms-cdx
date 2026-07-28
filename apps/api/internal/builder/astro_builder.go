package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cms-builder/api/internal/models"
)

const (
	SupromailTemplateKey = "supromail"
	astroBuildTimeout    = 2 * time.Minute
)

// TemplateBuilder routes trusted Astro templates to the Node-based builder and
// leaves the existing Go renderer responsible for all other registered templates.
type TemplateBuilder struct {
	Fallback LocalBuilder
	Astro    AstroBuilder
}

type AstroBuilder struct {
	OutputRoot       string
	BuilderDirectory string
	NPMCommand       string
	BuildTimeout     time.Duration
}

type TemplatePreviewer interface {
	GenerateTemplatePreview(ctx context.Context, templateKey string) (string, error)
}

type astroBuildData struct {
	Site     astroSite      `json:"site"`
	Articles []astroArticle `json:"articles"`
	Sections []astroSection `json:"sections"`
}

type astroSite struct {
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	BlogPath string `json:"blogPath"`
}

type astroArticle struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Excerpt         string `json:"excerpt"`
	ContentMarkdown string `json:"contentMarkdown"`
	SEOTitle        string `json:"seoTitle"`
	SEODescription  string `json:"seoDescription"`
	CanonicalURL    string `json:"canonicalUrl"`
	CategoryName    string `json:"categoryName"`
	AuthorName      string `json:"authorName"`
	PublishedAt     string `json:"publishedAt"`
	ReadingTime     int    `json:"readingTime"`
	CoverImageURL   string `json:"coverImageUrl"`
	IsFeatured      bool   `json:"isFeatured"`
}

type astroSection struct {
	SectionKey  string `json:"sectionKey"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	ContentJSON string `json:"contentJson"`
	IsEnabled   bool   `json:"isEnabled"`
}

func NewTemplateBuilder(outputRoot, builderDirectory, npmCommand string) TemplateBuilder {
	fallback := NewLocalBuilder(outputRoot)
	return TemplateBuilder{
		Fallback: fallback,
		Astro: AstroBuilder{
			OutputRoot:       fallback.OutputRoot,
			BuilderDirectory: builderDirectory,
			NPMCommand:       npmCommand,
			BuildTimeout:     astroBuildTimeout,
		},
	}
}

func (b TemplateBuilder) GenerateSite(ctx context.Context, content SiteContent, options GenerateOptions) (string, error) {
	if strings.EqualFold(strings.TrimSpace(content.Site.TemplateKey), SupromailTemplateKey) {
		return b.Astro.GenerateSite(ctx, content, options)
	}
	return b.Fallback.GenerateSite(ctx, content, options)
}

func (b TemplateBuilder) GenerateTemplatePreview(ctx context.Context, templateKey string) (string, error) {
	return b.Astro.GenerateTemplatePreview(ctx, templateKey)
}

func (b AstroBuilder) GenerateSite(ctx context.Context, content SiteContent, options GenerateOptions) (string, error) {
	if strings.TrimSpace(content.Site.Slug) == "" {
		return "", errors.New("site slug is required")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	blogPath, err := models.CanonicalBlogPath(content.Site.BlogPath)
	if err != nil {
		return "", fmt.Errorf("invalid blog path: %w", err)
	}
	outputPath := buildOutputPath(b.OutputRoot, content.Site, options.Preview)
	return b.generateSiteToOutput(ctx, content, options, blogPath, outputPath)
}

func (b AstroBuilder) GenerateTemplatePreview(ctx context.Context, templateKey string) (string, error) {
	if templateKey != SupromailTemplateKey {
		return "", errors.New("Astro template preview is unavailable")
	}
	content := SiteContent{
		Site: models.Site{
			Name:        "Northstar Journal",
			Slug:        "template-preview-supromail",
			Domain:      "https://example.com",
			BlogPath:    "/blog",
			TemplateKey: SupromailTemplateKey,
		},
		Articles: []ArticleContent{{
			ID: "preview-article", Title: "A calmer way to build on the web", Slug: "calmer-web",
			Excerpt:         "A practical guide to making focused products that remain fast, useful, and easy to maintain.",
			ContentMarkdown: "## A calmer way to build\n\nPractical systems make space for careful work.",
			Status:          "published", IsFeatured: true, AuthorName: "Northstar Editors", CategoryName: "Guides",
			PublishedAt: "2026-07-17T00:00:00Z", CoverImageURL: "https://images.example.com/preview.jpg",
		}},
	}
	return b.generateSiteToOutput(ctx, content, GenerateOptions{Preview: true}, "/blog", filepath.Join(b.OutputRoot, "template-previews", templateKey))
}

func (b AstroBuilder) generateSiteToOutput(ctx context.Context, content SiteContent, options GenerateOptions, blogPath, outputPath string) (string, error) {
	builderDirectory, err := absoluteDirectory(b.BuilderDirectory)
	if err != nil {
		return "", err
	}
	templateRoot, err := supromailTemplateRoot(builderDirectory)
	if err != nil {
		return "", err
	}
	if err := ensureOutputPath(b.OutputRoot, outputPath); err != nil {
		return "", err
	}
	buildInput, err := newAstroBuildData(content, blogPath, options)
	if err != nil {
		return "", err
	}
	buildData, err := json.Marshal(buildInput)
	if err != nil {
		return "", fmt.Errorf("encode Astro build data: %w", err)
	}
	dataFile, err := os.CreateTemp("", "cms-builder-astro-data-*.json")
	if err != nil {
		return "", fmt.Errorf("create Astro build data file: %w", err)
	}
	dataPath := dataFile.Name()
	defer os.Remove(dataPath)
	if _, err := dataFile.Write(buildData); err != nil {
		dataFile.Close()
		return "", fmt.Errorf("write Astro build data: %w", err)
	}
	if err := dataFile.Close(); err != nil {
		return "", fmt.Errorf("close Astro build data: %w", err)
	}

	timeout := b.BuildTimeout
	if timeout <= 0 {
		timeout = astroBuildTimeout
	}
	buildCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	command := exec.CommandContext(buildCtx, firstNonEmpty(b.NPMCommand, "npm"), "run", "build")
	command.Dir = builderDirectory
	command.Env = astroBuildEnvironment(dataPath, outputPath, content.Site.Domain, templateRoot)
	output, err := command.CombinedOutput()
	if err != nil {
		if errors.Is(buildCtx.Err(), context.DeadlineExceeded) {
			return "", errors.New("Astro template build timed out")
		}
		return "", fmt.Errorf("Astro template build failed: %s", boundedBuildOutput(output))
	}
	if err := relocateAstroBlogOutput(outputPath, blogPath); err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(outputPath, "index.html")); err != nil {
		return "", errors.New("Astro template build did not generate index.html")
	}
	return outputPath, nil
}

// relocateAstroBlogOutput adapts Astro's fixed /articles source route to the
// per-site CMS blog path after the static build has completed.
func relocateAstroBlogOutput(outputPath, blogPath string) error {
	if blogPath == "/articles" {
		return nil
	}
	source := filepath.Join(outputPath, "articles")
	destination := filepath.Join(outputPath, filepath.FromSlash(strings.TrimPrefix(blogPath, "/")))
	if _, err := os.Stat(source); err != nil {
		return errors.New("Astro template build did not generate article routes")
	}
	if _, err := os.Stat(destination); err == nil {
		return errors.New("Astro template build generated a conflicting blog route")
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Astro blog route: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create blog route directory: %w", err)
	}
	if err := os.Rename(source, destination); err != nil {
		return fmt.Errorf("move Astro blog route: %w", err)
	}
	return nil
}

func newAstroBuildData(content SiteContent, blogPath string, options GenerateOptions) (astroBuildData, error) {
	articles := selectedAstroArticles(content.Articles, options)
	sort.SliceStable(articles, func(left, right int) bool {
		if articles[left].IsFeatured != articles[right].IsFeatured {
			return articles[left].IsFeatured
		}
		return articles[left].PublishedAt > articles[right].PublishedAt
	})

	sections := make([]astroSection, 0, len(content.LandingSections))
	for _, section := range content.LandingSections {
		contentJSON, err := json.Marshal(section.ContentJSON)
		if err != nil {
			return astroBuildData{}, fmt.Errorf("encode landing section %q: %w", section.SectionKey, err)
		}
		sections = append(sections, astroSection{
			SectionKey:  section.SectionKey,
			Title:       section.Title,
			Subtitle:    section.Subtitle,
			ContentJSON: string(contentJSON),
			IsEnabled:   section.IsEnabled,
		})
	}

	return astroBuildData{
		Site:     astroSite{Name: content.Site.Name, Domain: content.Site.Domain, BlogPath: blogPath},
		Articles: articles,
		Sections: sections,
	}, nil
}

func selectedAstroArticles(articles []ArticleContent, options GenerateOptions) []astroArticle {
	selected := make(map[string]struct{}, len(options.ArticleIDs))
	for _, articleID := range options.ArticleIDs {
		if articleID = strings.TrimSpace(articleID); articleID != "" {
			selected[articleID] = struct{}{}
		}
	}

	items := make([]astroArticle, 0, len(articles))
	for _, article := range articles {
		if options.Preview && len(selected) > 0 {
			if _, found := selected[article.ID]; !found {
				continue
			}
		}
		items = append(items, astroArticle{
			Title:           article.Title,
			Slug:            article.Slug,
			Excerpt:         article.Excerpt,
			ContentMarkdown: article.ContentMarkdown,
			SEOTitle:        article.SEOTitle,
			SEODescription:  article.SEODescription,
			CanonicalURL:    article.CanonicalURL,
			CategoryName:    article.CategoryName,
			AuthorName:      article.AuthorName,
			PublishedAt:     article.PublishedAt,
			CoverImageURL:   article.CoverImageURL,
			IsFeatured:      article.IsFeatured,
		})
	}
	return items
}

func absoluteDirectory(path string) (string, error) {
	directory, err := filepath.Abs(strings.TrimSpace(path))
	if err != nil {
		return "", fmt.Errorf("resolve Astro builder directory: %w", err)
	}
	info, err := os.Stat(directory)
	if err != nil || !info.IsDir() {
		return "", errors.New("Astro builder directory is unavailable")
	}
	return directory, nil
}

func supromailTemplateRoot(builderDirectory string) (string, error) {
	for _, candidate := range []string{
		filepath.Join(builderDirectory, "templates", SupromailTemplateKey),
		filepath.Join(builderDirectory, "..", "..", "packages", "templates", SupromailTemplateKey),
	} {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return filepath.Abs(candidate)
		}
	}
	return "", errors.New("Supromail Astro template is unavailable")
}

func ensureOutputPath(outputRoot, outputPath string) error {
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return fmt.Errorf("resolve build output root: %w", err)
	}
	path, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("resolve build output path: %w", err)
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return errors.New("Astro build output escapes configured output root")
	}
	return nil
}

func astroBuildEnvironment(dataPath, outputPath, siteURL, templateRoot string) []string {
	environment := make([]string, 0, 12)
	for _, key := range []string{"PATH", "HOME", "TMPDIR", "TEMP", "TMP", "SystemRoot", "SYSTEMROOT", "ComSpec", "COMSPEC", "PATHEXT"} {
		if value, ok := os.LookupEnv(key); ok {
			environment = append(environment, key+"="+value)
		}
	}
	return append(environment,
		"CMS_BUILD_DATA_FILE="+dataPath,
		"CMS_BUILD_OUTPUT_DIR="+outputPath,
		"CMS_SITE_URL="+firstNonEmpty(siteURL, "http://localhost:8081"),
		"CMS_TEMPLATE_KEY="+SupromailTemplateKey,
		"CMS_TEMPLATE_ROOT="+templateRoot,
	)
}

func firstNonEmpty(value, fallback string) string {
	if value = strings.TrimSpace(value); value != "" {
		return value
	}
	return fallback
}

func boundedBuildOutput(output []byte) string {
	message := strings.TrimSpace(strings.ReplaceAll(string(output), "\x00", ""))
	if message == "" {
		return "command returned no output"
	}
	const maxLength = 4096
	if len(message) > maxLength {
		return message[:maxLength] + "…"
	}
	return message
}
