package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/models"
)

func TestGenerateSiteWritesArticlesToConfiguredBlogPath(t *testing.T) {
	outputRoot := t.TempDir()
	outputPath, err := NewLocalBuilder(outputRoot).GenerateSite(context.Background(), SiteContent{
		Site: models.Site{
			Slug:     "example-site",
			BlogPath: "/blog",
		},
		Articles: []ArticleContent{{
			Slug: "hello-world",
		}},
	}, GenerateOptions{})
	if err != nil {
		t.Fatalf("GenerateSite returned error: %v", err)
	}

	for _, relativePath := range []string{
		filepath.Join("blog", "index.html"),
		filepath.Join("blog", "hello-world", "index.html"),
	} {
		if _, err := os.Stat(filepath.Join(outputPath, relativePath)); err != nil {
			t.Fatalf("expected generated page at %s: %v", relativePath, err)
		}
	}
	if _, err := os.Stat(filepath.Join(outputPath, "articles", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("expected legacy articles output not to exist, got %v", err)
	}
}

func TestGenerateAnonimeArticlesPagination(t *testing.T) {
	articles := make([]ArticleContent, 7)
	for index := range articles {
		articles[index] = ArticleContent{
			Title:       fmt.Sprintf("Page article %d", index+1),
			Slug:        fmt.Sprintf("page-article-%d", index+1),
			PublishedAt: fmt.Sprintf("2026-06-%02d", index+1),
		}
	}

	outputPath, err := NewLocalBuilder(t.TempDir()).GenerateSite(context.Background(), SiteContent{
		Site: models.Site{
			Name: "Anonime", Slug: "anonime", BlogPath: "/blog", TemplateKey: "anonime", ThemeConfig: map[string]any{},
		},
		Articles: articles,
	}, GenerateOptions{})
	if err != nil {
		t.Fatalf("GenerateSite returned error: %v", err)
	}

	firstPage, err := os.ReadFile(filepath.Join(outputPath, "blog", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	secondPage, err := os.ReadFile(filepath.Join(outputPath, "blog", "page", "2", "index.html"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(firstPage), `href="/blog/page/2/"`) {
		t.Fatalf("expected the first page to link to page 2, got %s", firstPage)
	}
	if count := strings.Count(string(firstPage), `<article class="anonime-card anonime-article-card">`); count != anonimeArticlesPerPage {
		t.Fatalf("expected %d article cards on the first page, got %d", anonimeArticlesPerPage, count)
	}
	if !strings.Contains(string(secondPage), `Page article 7`) {
		t.Fatalf("expected second page to render its final article, got %s", secondPage)
	}
	if count := strings.Count(string(secondPage), `<article class="anonime-card anonime-article-card">`); count != 1 {
		t.Fatalf("expected one article card on the second page, got %d", count)
	}
	if !strings.Contains(string(secondPage), `href="/blog/" aria-label="Previous page"`) {
		t.Fatalf("expected page 2 to link back to the first page, got %s", secondPage)
	}

	sitemap, err := os.ReadFile(filepath.Join(outputPath, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(sitemap), `/blog/page/2/`) {
		t.Fatalf("expected sitemap to include paginated articles page, got %s", sitemap)
	}
}

func TestGeneratedPagesUseConfiguredBrandAssets(t *testing.T) {
	outputPath, err := NewLocalBuilder(t.TempDir()).GenerateSite(context.Background(), SiteContent{Site: models.Site{
		Name: "Example", Slug: "example", BlogPath: "/blog", TemplateKey: "default-blog", ThemeConfig: map[string]any{},
		Description: "A real description", LogoURL: "https://cdn.example/logo.png", FaviconURL: "https://cdn.example/favicon.png",
	}}, GenerateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(outputPath, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(contents)
	for _, expected := range []string{"A real description", `src="https://cdn.example/logo.png"`, `rel="icon" href="https://cdn.example/favicon.png"`} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected generated HTML to contain %q", expected)
		}
	}
}

func TestRenderHomeAndArticleLinksAreRelative(t *testing.T) {
	site := renderedSite{
		Site: models.Site{
			Slug:         "anonime-blog",
			Domain:       "https://anonime-dev.web.app",
			BlogPath:     "/articles",
			TemplateKey:  "default-blog",
			ThemeConfig:  map[string]any{},
			DeployConfig: map[string]any{},
		},
		BasePath:      "/articles",
		PublicBaseURL: "https://anonime-dev.web.app",
		Title:         "Anonime",
	}
	article := ArticleContent{
		Title:  "Hello world",
		Slug:   "hello-world",
		Status: "published",
	}

	home := renderHomePage(site)
	if !strings.Contains(home, `href="/articles/"`) {
		t.Fatalf("expected browse articles link to be relative, got %s", home)
	}
	if strings.Contains(home, "https://anonime-dev.web.app/articles/") {
		t.Fatalf("expected home page browse link not to hardcode the public base url")
	}

	card := renderArticleCard(site, article, false)
	if !strings.Contains(card, `href="/articles/hello-world/"`) {
		t.Fatalf("expected article card link to be relative, got %s", card)
	}
	if strings.Contains(card, "https://anonime-dev.web.app/articles/hello-world/") {
		t.Fatalf("expected article card not to hardcode the public base url")
	}

	if got := articleURL(site, article); got != "https://anonime-dev.web.app/articles/hello-world/" {
		t.Fatalf("expected canonical article URL to remain absolute, got %q", got)
	}
}

func TestRenderMarkdownSupportsRichTextEditorBlocks(t *testing.T) {
	output := renderMarkdown("## A heading\n\n**Bold** and *italic* with a [link](https://example.com).\n\n> A quoted thought\n\n- First item\n- Second item")

	for _, expected := range []string{
		"<h2>A heading</h2>",
		"<strong>Bold</strong>",
		"<em>italic</em>",
		`<a href="https://example.com">link</a>`,
		"<blockquote>",
		"<ul>",
		"<li>First item</li>",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("renderMarkdown() output did not contain %q:\n%s", expected, output)
		}
	}
}

func TestRenderAnonimeTemplateUsesAnonimeLayout(t *testing.T) {
	site := renderedSite{
		Site: models.Site{
			Slug:         "anonime-blog",
			Domain:       "https://anonime-dev.web.app",
			BlogPath:     "/articles",
			TemplateKey:  "anonime",
			ThemeConfig:  map[string]any{},
			DeployConfig: map[string]any{},
		},
		BasePath:      "/articles",
		PublicBaseURL: "https://anonime-dev.web.app",
		Title:         "Anonime",
		Theme:         themeForSite(models.Site{TemplateKey: "anonime", Name: "Anonime"}),
		Articles: []ArticleContent{
			{
				Title:           "Who Needs Encrypted Communication?",
				Slug:            "who-needs-encrypted-communication",
				Excerpt:         "Encrypted communication matters for everyone.",
				ContentMarkdown: "# Heading\n\nBody copy for testing.",
				Status:          "published",
				CategoryName:    "Security",
				AuthorName:      "Rene Carter",
				PublishedAt:     "Feb 7, 2026",
			},
		},
		RecentArticles: []ArticleContent{
			{
				Title:           "Latest insight",
				Slug:            "latest-insight",
				Excerpt:         "Latest insight summary",
				ContentMarkdown: "Latest insight body",
				Status:          "published",
				CategoryName:    "Guides",
				AuthorName:      "Rene Carter",
				PublishedAt:     "Feb 2, 2026",
			},
		},
		FeaturedArticles: []ArticleContent{
			{
				Title:           "Featured insight",
				Slug:            "featured-insight",
				Excerpt:         "Featured summary",
				ContentMarkdown: "Featured body",
				Status:          "published",
				CategoryName:    "Privacy Basics",
				AuthorName:      "Rene Carter",
				PublishedAt:     "Feb 1, 2026",
			},
		},
	}

	home := renderHomePage(site)
	if !strings.Contains(home, `class="anonime-template"`) {
		t.Fatalf("expected anonime template wrapper, got %s", home)
	}
	if !strings.Contains(home, `Thoughts that protect your right to`) {
		t.Fatalf("expected anonime hero copy, got %s", home)
	}
	if !strings.Contains(home, `href="/articles/"`) {
		t.Fatalf("expected anonime browse link to remain relative, got %s", home)
	}
	if !strings.Contains(home, `data-art="shield"`) {
		t.Fatalf("expected anonime content cards to use anonime art variants, got %s", home)
	}
	if count := strings.Count(home, `id="latest-insights"`); count != 1 {
		t.Fatalf("expected one latest insights section, got %d", count)
	}
}

func TestRenderAnonimeChromeUsesProductionLinksAndCSSOnlyTheme(t *testing.T) {
	site := renderedSite{Site: models.Site{
		Name:    "Anonime",
		LogoURL: "https://cdn.example/anonime-logo.svg",
	}}

	header := renderAnonimeHeader(site)
	for _, expected := range []string{
		`class="brand-logo" src="https://cdn.anonime.io/anonime-logo.svg" alt="Anonime"`,
		`href="https://anonime.io/#top">Home</a>`,
		`href="https://anonime.io/#how-it-works">How it Works</a>`,
		`href="https://anonime.io/pricing">Plans &amp; Pricing</a>`,
		`href="https://anonime.io/blog" aria-current="page">Blog</a>`,
		`type="checkbox" id="anonime-theme-toggle"`,
		`for="anonime-theme-toggle"`,
		`Toggle dark mode`,
		`href="https://app.anonime.io">Log in</a>`,
		`href="https://app.anonime.io">Get Started</a>`,
		`aria-label="Open navigation"`,
		`aria-label="Mobile navigation"`,
	} {
		if !strings.Contains(header, expected) {
			t.Errorf("renderAnonimeHeader() did not contain %q", expected)
		}
	}

	styles := anonimeStyles()
	for _, expected := range []string{
		`body.anonime-layout:has(#anonime-theme-toggle:checked)`,
		`color-scheme: dark`,
		`.anonime-theme-checkbox:checked + .anonime-theme-toggle .anonime-theme-sun`,
		`@media (max-width: 920px)`,
		`@media (max-width: 620px)`,
		`font-family: Matter, Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`,
	} {
		if !strings.Contains(styles, expected) {
			t.Errorf("anonimeStyles() did not contain %q", expected)
		}
	}

	footer := renderAnonimeFooter(site)
	for _, expected := range []string{
		`Private communication<br>without <strong>permanent trails.</strong>`,
		`<h3>Product</h3>`,
		`href="https://anonime.io/#security">Whispers</a>`,
		`href="https://anonime.io/#private-drops">Private Drops</a>`,
		`href="https://anonime.io/#how-it-works">Nyms</a>`,
		`href="https://anonime.io/pricing">Pricing</a>`,
		`<h3>Company</h3>`,
		`href="https://anonime.io/#about">About</a>`,
		`href="https://anonime.io/blog">Blog</a>`,
		`href="https://anonime.io/contact">Contact</a>`,
		`<h3>Legal</h3>`,
		`href="https://anonime.io/privacy">Privacy Policy</a>`,
		`href="https://anonime.io/terms">Terms of Service</a>`,
		`href="https://anonime.io/security">Security</a>`,
		`href="https://anonime.io/acceptable-usage-policy">Acceptable Use Policy</a>`,
		`href="https://anonime.io/data-processing">Data Processing</a>`,
		`Built with <strong>privacy by design.</strong>`,
		`End-to-end encrypted`,
		`No personal data required`,
		`No logs. No tracking.`,
		`Hosted with privacy-first infra`,
		`&copy; 2026 Anonime, Inc. All rights reserved.`,
		`class="footer-social-link" href="https://x.com/anonimehq" aria-label="X" target="_blank" rel="noopener noreferrer"`,
		`Your privacy is your right. <strong>We protect it every day.</strong>`,
	} {
		if !strings.Contains(footer, expected) {
			t.Errorf("renderAnonimeFooter() did not contain %q", expected)
		}
	}
}

func TestRenderAnonimeMetaUsesDateOnlyForTimestampValues(t *testing.T) {
	site := renderedSite{BasePath: "/articles"}
	article := ArticleContent{
		Title:       "Timestamp article",
		Slug:        "timestamp-article",
		AuthorName:  "Rene Carter",
		PublishedAt: "2026-06-04T08:35:39Z",
	}

	featuredMeta := renderAnonimeMeta(site, article, true)
	articleFooter := renderAnonimeMeta(site, article, false)
	sidebarMeta := renderAnonimeCompactArticle(site, article)
	for _, output := range []string{featuredMeta, articleFooter} {
		if !strings.Contains(output, `>2026-06-04<`) {
			t.Fatalf("expected date-only timestamp label, got %s", output)
		}
		if strings.Contains(output, article.PublishedAt) {
			t.Fatalf("expected timestamp to be omitted, got %s", output)
		}
	}
	if !strings.Contains(sidebarMeta, `>2026-06-04 <span`) || strings.Contains(sidebarMeta, article.PublishedAt) {
		t.Fatalf("expected sidebar date-only timestamp label, got %s", sidebarMeta)
	}
}

func TestAnonimeArticleCardsKeepArtLeftAndCopyRight(t *testing.T) {
	styles := anonimeStyles()
	for _, expected := range []string{
		`grid-template-areas: "media copy"`,
		`.anonime-article-card > a:first-child { grid-area: media; min-height: 100%; }`,
		`.anonime-article-card .anonime-article-art { height: 100%; min-height: 100%; }`,
		`.anonime-sidebar-item { display: grid; grid-template-columns: 84px minmax(0, 1fr);`,
		`-webkit-line-clamp: 2; line-clamp: 2;`,
	} {
		if !strings.Contains(styles, expected) {
			t.Fatalf("expected Anonime layout rule %q", expected)
		}
	}
}
