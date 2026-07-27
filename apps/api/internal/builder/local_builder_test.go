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
		filepath.Join("blog", "sitemap.xml"),
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

	blogLanding, err := os.ReadFile(filepath.Join(outputPath, "blog", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	firstPage, err := os.ReadFile(filepath.Join(outputPath, "blog", "articles", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	secondPage, err := os.ReadFile(filepath.Join(outputPath, "blog", "articles", "page", "2", "index.html"))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(blogLanding), `Thoughts that protect your right to`) || strings.Contains(string(blogLanding), `All <strong>Articles</strong>`) {
		t.Fatalf("expected /blog to render the Anonime blog landing page, got %s", blogLanding)
	}
	if !strings.Contains(string(firstPage), `href="/blog/articles/page/2/"`) {
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
	if !strings.Contains(string(secondPage), `href="/blog/articles/" aria-label="Previous page"`) {
		t.Fatalf("expected page 2 to link back to the first page, got %s", secondPage)
	}

	sitemap, err := os.ReadFile(filepath.Join(outputPath, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(sitemap), `/blog/articles/`) || !strings.Contains(string(sitemap), `/blog/articles/page/2/`) {
		t.Fatalf("expected sitemap to include paginated articles page, got %s", sitemap)
	}
	blogSitemap, err := os.ReadFile(filepath.Join(outputPath, "blog", "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(blogSitemap), `<loc>http://localhost:8081/</loc>`) || !strings.Contains(string(blogSitemap), `/blog/articles/page/2/`) {
		t.Fatalf("expected deployable blog sitemap to include only blog URLs, got %s", blogSitemap)
	}
}

func TestArticlePagesRenderSEOAndStructuredData(t *testing.T) {
	site := renderedSite{
		Site:          models.Site{Name: "Example", Domain: "https://example.com", BlogPath: "/blog"},
		BasePath:      "/blog",
		PublicBaseURL: "https://example.com",
		Title:         "Example",
		Description:   "Site description",
	}
	article := ArticleContent{
		Title:           "Article title",
		Slug:            "article-title",
		Excerpt:         "Article excerpt",
		SEOTitle:        "Search title",
		SEODescription:  "Search description",
		CoverImageURL:   "/blog/cover.webp",
		PublishedAt:     "2026-07-01T10:00:00Z",
		UpdatedAt:       "2026-07-02T10:00:00Z",
		AuthorName:      "Ada Example",
		Tags:            []TagContent{{Name: "Privacy"}, {Name: "Security"}},
		ContentMarkdown: "Body",
	}

	page := renderDefaultArticlePage(site, article)
	for _, expected := range []string{
		`<title>Search title | Example</title>`,
		`name="description" content="Search description"`,
		`rel="canonical" href="https://example.com/blog/article-title/"`,
		`property="og:type" content="article"`,
		`property="og:image" content="https://example.com/blog/cover.webp"`,
		`name="twitter:card" content="summary_large_image"`,
		`"@type":"BlogPosting"`,
		`"author":{"@type":"Person","name":"Ada Example"}`,
		`"keywords":"Privacy, Security"`,
	} {
		if !strings.Contains(page, expected) {
			t.Errorf("expected article SEO output to contain %q", expected)
		}
	}
}

func TestPreviewPagesAreNotIndexable(t *testing.T) {
	page := renderDefaultArticlePage(renderedSite{Preview: true, Title: "Example", Description: "Preview"}, ArticleContent{Title: "Draft", Slug: "draft"})
	if !strings.Contains(page, `name="robots" content="noindex,nofollow"`) {
		t.Fatalf("expected preview page to opt out of indexing, got %s", page)
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
	for _, expected := range []string{
		`class="anonime-hero-art"`,
		`anonime-blog-hero-white.webp`,
		`anonime-blog-hero-black.webp`,
	} {
		if !strings.Contains(home, expected) {
			t.Errorf("expected rendered anonime hero to contain %q", expected)
		}
	}
	if !strings.Contains(home, `href="/articles/articles/"`) {
		t.Fatalf("expected anonime browse link to remain relative, got %s", home)
	}
	if !strings.Contains(home, `data-art="shield"`) {
		t.Fatalf("expected anonime content cards to use anonime art variants, got %s", home)
	}
	if count := strings.Count(home, `id="latest-insights"`); count != 1 {
		t.Fatalf("expected one latest insights section, got %d", count)
	}
}

func TestAnonimeArticleEyebrowUsesArticleCategory(t *testing.T) {
	site := renderedSite{
		BasePath: "/blog",
		Theme:    themeForSite(models.Site{TemplateKey: "anonime", Name: "Anonime"}),
	}
	article := ArticleContent{
		Title:        "Protecting your identity",
		Slug:         "protecting-your-identity",
		CategoryName: "Identity",
	}

	page := renderAnonimeArticleLayout(site, article)
	if !strings.Contains(page, `<p class="anonime-eyebrow">Identity</p>`) {
		t.Fatalf("expected the article category in the eyebrow, got %s", page)
	}
	if strings.Contains(page, `<p class="anonime-eyebrow">Secure communication</p>`) {
		t.Fatalf("expected the hard-coded article eyebrow to be removed, got %s", page)
	}
}

func TestRenderAnonimeChromeUsesProductionLinksAndPersistentTheme(t *testing.T) {
	site := renderedSite{Site: models.Site{
		Name:    "Anonime",
		LogoURL: "https://cdn.example/anonime-logo.svg",
	}}

	header := renderAnonimeHeader(site)
	for _, expected := range []string{
		`class="brand-logo brand-logo-light" src="https://cdn.anonime.io/anonime-logo.svg" alt="Anonime"`,
		`class="brand-logo brand-logo-dark" src="https://cdn.anonime.io/anonime-logo-dark.svg" alt="" aria-hidden="true"`,
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

	styles := renderStyles(themeForSite(models.Site{TemplateKey: "anonime", Name: "Anonime"}))
	for _, expected := range []string{
		`html[data-theme="dark"] body.anonime-layout`,
		`.anonime-template .brand-logo-dark { display: none; }`,
		`html[data-theme="dark"] body.anonime-layout .brand-logo-light { display: none; }`,
		`html[data-theme="dark"] body.anonime-layout .brand-logo-dark { display: block; }`,
		`color-scheme: dark`,
		`.anonime-theme-checkbox:checked + .anonime-theme-toggle .anonime-theme-sun`,
		`@media (max-width: 920px)`,
		`@media (max-width: 620px)`,
		`font-family: Matter, Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`,
		`.anonime-hero-art .anonime-hero-art-dark { display: none; }`,
		`html[data-theme="dark"] body.anonime-layout .anonime-hero-art-light { display: none; }`,
		`html[data-theme="dark"] body.anonime-layout .anonime-hero-art-dark { display: block; }`,
		`@media (prefers-color-scheme: dark)`,
		`body.anonime-layout .anonime-hero-art-light { display: none; }`,
		`body.anonime-layout .anonime-hero-art-dark { display: block; }`,
	} {
		if !strings.Contains(styles, expected) {
			t.Errorf("anonimeStyles() did not contain %q", expected)
		}
	}
	page := renderAnonimeHomePage(renderedSite{Theme: themeForSite(models.Site{TemplateKey: "anonime", Name: "Anonime"})})
	for _, expected := range []string{`localStorage.getItem(key)`, `localStorage.setItem(key,dark?"dark":"light")`, `root.dataset.theme=dark?"dark":"light"`} {
		if !strings.Contains(page, expected) {
			t.Errorf("expected generated document to include persistent theme behavior %q", expected)
		}
	}
	if strings.Contains(styles, `.anonime-hero-art::before`) || strings.Contains(styles, `.anonime-hero-art::after`) {
		t.Error("expected hero art pseudo-elements to be removed")
	}

	footer := renderAnonimeFooter(site)
	for _, expected := range []string{
		`Private communication<br>without <strong>permanent trails.</strong>`,
		`class="anonime-footer-logo" src="https://cdn.anonime.io/anonime-logo-dark.svg" alt="Anonime"`,
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
