package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"cms-builder/api/internal/models"
	"github.com/yuin/goldmark"
)

type LocalBuilder struct {
	OutputRoot string
}

type themeVariant struct {
	Name                   string
	Background             string
	Surface                string
	SurfaceStrong          string
	Text                   string
	Muted                  string
	Accent                 string
	AccentContrast         string
	Border                 string
	LayoutClass            string
	CardClass              string
	HeroClass              string
	ArticleAsideClass      string
	ArticleBodyClass       string
	ArticleGridColumns     string
	SecondaryAccent        string
	FeatureBadgeBackground string
}

type renderedSite struct {
	Site             models.Site
	Theme            themeVariant
	Preview          bool
	BuildType        string
	BasePath         string
	PublicBaseURL    string
	Title            string
	Description      string
	LandingSections  []LandingSectionContent
	Articles         []ArticleContent
	FeaturedArticles []ArticleContent
	RecentArticles   []ArticleContent
	Now              time.Time
}

type documentMetadata struct {
	Title       string
	Description string
	Canonical   string
	Image       string
	Type        string
	PublishedAt string
	ModifiedAt  string
	Author      string
	Tags        []TagContent
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type urlSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
	GUID        string `xml:"guid"`
}

func NewLocalBuilder(outputRoot string) LocalBuilder {
	if strings.TrimSpace(outputRoot) == "" {
		outputRoot = filepath.Join("dist", "builds")
	}
	return LocalBuilder{OutputRoot: outputRoot}
}

func (b LocalBuilder) GenerateSite(ctx context.Context, content SiteContent, options GenerateOptions) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	if strings.TrimSpace(content.Site.Slug) == "" {
		return "", errors.New("site slug is required")
	}
	blogPath, err := models.CanonicalBlogPath(content.Site.BlogPath)
	if err != nil {
		return "", fmt.Errorf("invalid blog path: %w", err)
	}

	articles := content.Articles
	if options.Preview && len(options.ArticleIDs) > 0 {
		selected := make(map[string]struct{}, len(options.ArticleIDs))
		for _, id := range options.ArticleIDs {
			if trimmed := strings.TrimSpace(id); trimmed != "" {
				selected[trimmed] = struct{}{}
			}
		}

		filtered := make([]ArticleContent, 0, len(selected))
		for _, article := range articles {
			if _, ok := selected[article.ID]; ok {
				filtered = append(filtered, article)
			}
		}
		articles = filtered
	}

	sort.SliceStable(articles, func(left, right int) bool {
		if articles[left].IsFeatured != articles[right].IsFeatured {
			return articles[left].IsFeatured
		}
		if articles[left].PublishedAt != articles[right].PublishedAt {
			return articles[left].PublishedAt > articles[right].PublishedAt
		}
		return articles[left].UpdatedAt > articles[right].UpdatedAt
	})

	rendered := renderedSite{
		Site:             content.Site,
		Theme:            themeForSite(content.Site),
		Preview:          options.Preview,
		BuildType:        buildTypeLabel(options.Preview),
		BasePath:         blogPath,
		PublicBaseURL:    publicBaseURL(content.Site),
		Title:            content.Site.Name,
		Description:      siteDescription(content.Site),
		LandingSections:  content.LandingSections,
		Articles:         articles,
		FeaturedArticles: featuredArticles(articles),
		RecentArticles:   recentArticles(articles),
		Now:              time.Now().UTC(),
	}

	outputPath := buildOutputPath(b.OutputRoot, content.Site, options.Preview)
	if err := os.RemoveAll(outputPath); err != nil {
		return "", err
	}
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return "", err
	}

	if err := writeFile(filepath.Join(outputPath, "index.html"), renderHomePage(rendered)); err != nil {
		return "", err
	}
	blogOutputPath := filepath.Join(outputPath, filepath.FromSlash(strings.TrimPrefix(rendered.BasePath, "/")))
	articlesOutputPath := blogOutputPath
	if rendered.Theme.Name == "anonime" {
		if err := writeFile(filepath.Join(blogOutputPath, "index.html"), renderAnonimeBlogLandingPage(rendered)); err != nil {
			return "", err
		}
		articlesOutputPath = filepath.Join(blogOutputPath, "articles")
	}
	if err := writeFile(filepath.Join(articlesOutputPath, "index.html"), renderArticlesPage(rendered)); err != nil {
		return "", err
	}
	if rendered.Theme.Name == "anonime" {
		for page := 2; page <= anonimeArticlePageCount(len(rendered.Articles)); page++ {
			pageDir := filepath.Join(articlesOutputPath, "page", strconv.Itoa(page))
			if err := writeFile(filepath.Join(pageDir, "index.html"), renderAnonimeArticlesPageNumber(rendered, page)); err != nil {
				return "", err
			}
		}
	}
	for _, article := range articles {
		articleDir := filepath.Join(blogOutputPath, safePathSegment(article.Slug))
		if err := writeFile(filepath.Join(articleDir, "index.html"), renderArticlePage(rendered, article)); err != nil {
			return "", err
		}
	}
	if err := writeFile(filepath.Join(blogOutputPath, "sitemap.xml"), renderBlogSitemap(rendered)); err != nil {
		return "", err
	}
	if err := writeFile(filepath.Join(outputPath, "sitemap.xml"), renderSitemap(rendered)); err != nil {
		return "", err
	}
	if err := writeFile(filepath.Join(outputPath, "robots.txt"), renderRobots(rendered)); err != nil {
		return "", err
	}
	if err := writeFile(filepath.Join(outputPath, "rss.xml"), renderRSS(rendered)); err != nil {
		return "", err
	}

	return outputPath, nil
}

func buildOutputPath(root string, site models.Site, preview bool) string {
	if preview {
		return filepath.Join(root, "preview", safePathSegment(site.Slug))
	}
	return filepath.Join(root, "sites", safePathSegment(site.Slug))
}

func buildTypeLabel(preview bool) string {
	if preview {
		return "preview"
	}
	return "published"
}

func siteDescription(site models.Site) string {
	if value := strings.TrimSpace(site.Description); value != "" {
		return value
	}
	if brand, ok := site.ThemeConfig["brand"].(string); ok && strings.TrimSpace(brand) != "" {
		return brand
	}
	return site.Name
}

func publicBaseURL(site models.Site) string {
	if value := strings.TrimSpace(site.Domain); value != "" {
		return strings.TrimRight(value, "/")
	}
	return "http://localhost:8081"
}

func featuredArticles(articles []ArticleContent) []ArticleContent {
	items := make([]ArticleContent, 0)
	for _, article := range articles {
		if article.IsFeatured {
			items = append(items, article)
		}
	}
	return items
}

func recentArticles(articles []ArticleContent) []ArticleContent {
	limit := 4
	if len(articles) < limit {
		limit = len(articles)
	}
	items := make([]ArticleContent, 0, limit)
	items = append(items, articles[:limit]...)
	return items
}

func themeForSite(site models.Site) themeVariant {
	accent := themeString(site.ThemeConfig, "accent", "#1d4ed8")
	brand := themeString(site.ThemeConfig, "brand", site.Name)
	_ = brand
	switch strings.ToLower(strings.TrimSpace(site.TemplateKey)) {
	case "anonime":
		return themeVariant{
			Name:                   "anonime",
			Background:             "#f7fbfa",
			Surface:                "rgba(255, 255, 255, 0.92)",
			SurfaceStrong:          "#ffffff",
			Text:                   "#0f1728",
			Muted:                  "#5b6474",
			Accent:                 "#10b26c",
			AccentContrast:         "#ffffff",
			Border:                 "rgba(18, 28, 43, 0.1)",
			LayoutClass:            "anonime-layout",
			CardClass:              "anonime-card",
			HeroClass:              "anonime-hero",
			ArticleAsideClass:      "anonime-aside",
			ArticleBodyClass:       "anonime-body",
			ArticleGridColumns:     "repeat(2, minmax(0, 1fr))",
			SecondaryAccent:        "#059669",
			FeatureBadgeBackground: "rgba(16, 178, 108, 0.1)",
		}
	case "premium-saas":
		return themeVariant{
			Name:                   "premium-saas",
			Background:             "#07111f",
			Surface:                "#0f172a",
			SurfaceStrong:          "#111c34",
			Text:                   "#e5eefb",
			Muted:                  "#9fb3d1",
			Accent:                 accent,
			AccentContrast:         "#ffffff",
			Border:                 "#20304d",
			LayoutClass:            "premium-layout",
			CardClass:              "premium-card",
			HeroClass:              "premium-hero",
			ArticleAsideClass:      "premium-aside",
			ArticleBodyClass:       "premium-body",
			ArticleGridColumns:     "repeat(2, minmax(0, 1fr))",
			SecondaryAccent:        "#22c55e",
			FeatureBadgeBackground: "rgba(34, 197, 94, 0.16)",
		}
	default:
		return themeVariant{
			Name:                   "default-blog",
			Background:             "#f7f8fb",
			Surface:                "#ffffff",
			SurfaceStrong:          "#eef2ff",
			Text:                   "#111827",
			Muted:                  "#5b6472",
			Accent:                 accent,
			AccentContrast:         "#ffffff",
			Border:                 "#d7ddea",
			LayoutClass:            "default-layout",
			CardClass:              "default-card",
			HeroClass:              "default-hero",
			ArticleAsideClass:      "default-aside",
			ArticleBodyClass:       "default-body",
			ArticleGridColumns:     "repeat(2, minmax(0, 1fr))",
			SecondaryAccent:        "#6b7280",
			FeatureBadgeBackground: "rgba(29, 78, 216, 0.08)",
		}
	}
}

func themeString(values map[string]any, key, fallback string) string {
	if values == nil {
		return fallback
	}
	if value, ok := values[key].(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func writeFile(path, contents string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func renderHomePage(site renderedSite) string {
	switch site.Theme.Name {
	case "anonime":
		return renderAnonimeHomePage(site)
	default:
		return renderDefaultHomePage(site)
	}
}

func renderDefaultHomePage(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="hero">`)
	if logo := strings.TrimSpace(site.Site.LogoURL); logo != "" {
		body.WriteString(`<img class="site-logo" src="`)
		body.WriteString(html.EscapeString(logo))
		body.WriteString(`" alt="`)
		body.WriteString(html.EscapeString(site.Site.Name))
		body.WriteString(` logo">`)
	}
	body.WriteString(`<p class="eyebrow">`)
	body.WriteString(html.EscapeString(strings.ToUpper(site.BuildType)))
	body.WriteString(` BUILD</p>`)
	body.WriteString(`<h1>`)
	body.WriteString(html.EscapeString(site.Title))
	body.WriteString(`</h1>`)
	body.WriteString(`<p class="lead">`)
	body.WriteString(html.EscapeString(site.Description))
	body.WriteString(`</p>`)
	if site.Preview {
		body.WriteString(`<p class="notice">Preview build generated from the CMS.</p>`)
	}
	body.WriteString(`<div class="hero-actions">`)
	body.WriteString(`<a class="button primary" href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">Browse articles</a>`)
	body.WriteString(`</div></section>`)

	body.WriteString(`<section class="grid">`)
	if len(site.LandingSections) > 0 {
		body.WriteString(`<article class="panel">`)
		body.WriteString(`<h2>Landing sections</h2>`)
		body.WriteString(`<div class="section-list">`)
		for _, section := range site.LandingSections {
			if !section.IsEnabled {
				continue
			}
			body.WriteString(`<div class="section-card">`)
			body.WriteString(`<strong>`)
			body.WriteString(html.EscapeString(section.Title))
			body.WriteString(`</strong>`)
			if strings.TrimSpace(section.Subtitle) != "" {
				body.WriteString(`<p>`)
				body.WriteString(html.EscapeString(section.Subtitle))
				body.WriteString(`</p>`)
			}
			if len(section.ContentJSON) > 0 {
				body.WriteString(`<pre>`)
				body.WriteString(html.EscapeString(renderJSON(section.ContentJSON)))
				body.WriteString(`</pre>`)
			}
			body.WriteString(`</div>`)
		}
		body.WriteString(`</div></article>`)
	}

	body.WriteString(`<article class="panel">`)
	body.WriteString(`<h2>Featured articles</h2>`)
	if len(site.FeaturedArticles) == 0 {
		body.WriteString(`<p class="muted">No featured articles yet.</p>`)
	} else {
		body.WriteString(renderArticleCards(site, site.FeaturedArticles, true))
	}
	body.WriteString(`</article>`)

	body.WriteString(`<article class="panel">`)
	body.WriteString(`<h2>Recent articles</h2>`)
	if len(site.RecentArticles) == 0 {
		body.WriteString(`<p class="muted">No published articles yet.</p>`)
	} else {
		body.WriteString(renderArticleCards(site, site.RecentArticles, false))
	}
	body.WriteString(`</article>`)
	body.WriteString(`</section>`)

	return renderDocument(site, siteDocumentMetadata(site, "Home", site.PublicBaseURL+"/"), body.String())
}

func renderArticlesPage(site renderedSite) string {
	switch site.Theme.Name {
	case "anonime":
		return renderAnonimeArticlesPage(site)
	default:
		return renderDefaultArticlesPage(site)
	}
}

func renderDefaultArticlesPage(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="hero compact">`)
	body.WriteString(`<p class="eyebrow">ARTICLES</p>`)
	body.WriteString(`<h1>`)
	body.WriteString(html.EscapeString(site.Title))
	body.WriteString(` articles</h1>`)
	body.WriteString(`<p class="lead">Published posts for the selected site.</p></section>`)

	body.WriteString(`<section class="panel">`)
	body.WriteString(`<div class="article-list">`)
	if len(site.Articles) == 0 {
		body.WriteString(`<p class="muted">No articles available.</p>`)
	} else {
		for _, article := range site.Articles {
			body.WriteString(renderArticleCard(site, article, false))
		}
	}
	body.WriteString(`</div></section>`)

	return renderDocument(site, siteDocumentMetadata(site, "Articles", articlesURL(site)), body.String())
}

func renderArticlePage(site renderedSite, article ArticleContent) string {
	switch site.Theme.Name {
	case "anonime":
		return renderAnonimeArticlePage(site, article)
	default:
		return renderDefaultArticlePage(site, article)
	}
}

func renderDefaultArticlePage(site renderedSite, article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<article class="article-layout">`)
	body.WriteString(`<header class="article-hero">`)
	body.WriteString(`<p class="eyebrow">`)
	body.WriteString(html.EscapeString(strings.ToUpper(article.Status)))
	body.WriteString(`</p>`)
	body.WriteString(`<h1>`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</h1>`)
	if strings.TrimSpace(article.Excerpt) != "" {
		body.WriteString(`<p class="lead">`)
		body.WriteString(html.EscapeString(article.Excerpt))
		body.WriteString(`</p>`)
	}
	body.WriteString(`<div class="article-meta">`)
	if strings.TrimSpace(article.AuthorName) != "" {
		body.WriteString(`<span>By ` + html.EscapeString(article.AuthorName) + `</span>`)
	}
	if strings.TrimSpace(article.PublishedAt) != "" {
		body.WriteString(`<span>Published ` + html.EscapeString(article.PublishedAt) + `</span>`)
	}
	if strings.TrimSpace(article.CategoryName) != "" {
		body.WriteString(`<span>Category ` + html.EscapeString(article.CategoryName) + `</span>`)
	}
	body.WriteString(`</div>`)
	if strings.TrimSpace(article.CoverImageURL) != "" {
		body.WriteString(`<figure class="cover"><img src="`)
		body.WriteString(html.EscapeString(article.CoverImageURL))
		body.WriteString(`" alt="`)
		body.WriteString(html.EscapeString(article.Title))
		body.WriteString(` cover image"></figure>`)
	}
	body.WriteString(`</header>`)

	body.WriteString(`<div class="article-body">`)
	body.WriteString(renderMarkdown(article.ContentMarkdown))
	body.WriteString(`</div>`)

	if len(article.Tags) > 0 {
		body.WriteString(`<footer class="tag-list">`)
		for _, tag := range article.Tags {
			body.WriteString(`<span class="tag">`)
			body.WriteString(html.EscapeString(tag.Name))
			body.WriteString(`</span>`)
		}
		body.WriteString(`</footer>`)
	}
	body.WriteString(`</article>`)

	return renderDocument(site, articleDocumentMetadata(site, article), body.String())
}

func renderArticleCards(site renderedSite, articles []ArticleContent, featured bool) string {
	var body strings.Builder
	for _, article := range articles {
		body.WriteString(renderArticleCard(site, article, featured))
	}
	return body.String()
}

func renderArticleCard(site renderedSite, article ArticleContent, featured bool) string {
	var body strings.Builder
	body.WriteString(`<a class="article-card `)
	body.WriteString(site.Theme.CardClass)
	if featured {
		body.WriteString(` featured`)
	}
	body.WriteString(`" href="`)
	body.WriteString(html.EscapeString(articlePath(site, article)))
	body.WriteString(`">`)
	if strings.TrimSpace(article.CoverImageURL) != "" {
		body.WriteString(`<img class="article-image" src="`)
		body.WriteString(html.EscapeString(article.CoverImageURL))
		body.WriteString(`" alt="`)
		body.WriteString(html.EscapeString(article.Title))
		body.WriteString(` cover image">`)
	}
	body.WriteString(`<div class="article-card-content">`)
	body.WriteString(`<p class="article-status">`)
	body.WriteString(html.EscapeString(strings.ToUpper(article.Status)))
	body.WriteString(`</p>`)
	body.WriteString(`<h3>`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</h3>`)
	if strings.TrimSpace(article.Excerpt) != "" {
		body.WriteString(`<p>`)
		body.WriteString(html.EscapeString(article.Excerpt))
		body.WriteString(`</p>`)
	}
	if strings.TrimSpace(article.AuthorName) != "" || strings.TrimSpace(article.CategoryName) != "" {
		body.WriteString(`<div class="article-chip-row">`)
		if strings.TrimSpace(article.AuthorName) != "" {
			body.WriteString(`<span>`)
			body.WriteString(html.EscapeString(article.AuthorName))
			body.WriteString(`</span>`)
		}
		if strings.TrimSpace(article.CategoryName) != "" {
			body.WriteString(`<span>`)
			body.WriteString(html.EscapeString(article.CategoryName))
			body.WriteString(`</span>`)
		}
		body.WriteString(`</div>`)
	}
	body.WriteString(`</div></a>`)
	return body.String()
}

func renderDocument(site renderedSite, metadata documentMetadata, body string) string {
	var builder strings.Builder
	builder.WriteString(`<!doctype html><html lang="en"><head><meta charset="utf-8">`)
	builder.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	if favicon := strings.TrimSpace(site.Site.FaviconURL); favicon != "" {
		builder.WriteString(`<link rel="icon" href="`)
		builder.WriteString(html.EscapeString(favicon))
		builder.WriteString(`">`)
	}
	if site.Theme.Name == "anonime" {
		builder.WriteString(renderAnonimeThemeScript())
	}
	builder.WriteString(`<title>`)
	builder.WriteString(html.EscapeString(pageTitle(site, metadata.Title)))
	builder.WriteString(`</title>`)
	if strings.TrimSpace(metadata.Canonical) != "" {
		builder.WriteString(`<link rel="canonical" href="`)
		builder.WriteString(html.EscapeString(metadata.Canonical))
		builder.WriteString(`">`)
	}
	builder.WriteString(`<meta name="description" content="`)
	builder.WriteString(html.EscapeString(metadata.Description))
	builder.WriteString(`">`)
	if site.Preview {
		builder.WriteString(`<meta name="robots" content="noindex,nofollow">`)
	} else {
		builder.WriteString(`<meta name="robots" content="index,follow,max-image-preview:large,max-snippet:-1,max-video-preview:-1">`)
	}
	builder.WriteString(`<meta property="og:type" content="`)
	builder.WriteString(html.EscapeString(metadata.Type))
	builder.WriteString(`"><meta property="og:site_name" content="`)
	builder.WriteString(html.EscapeString(site.Title))
	builder.WriteString(`"><meta property="og:title" content="`)
	builder.WriteString(html.EscapeString(pageTitle(site, metadata.Title)))
	builder.WriteString(`"><meta property="og:description" content="`)
	builder.WriteString(html.EscapeString(metadata.Description))
	builder.WriteString(`">`)
	if strings.TrimSpace(metadata.Canonical) != "" {
		builder.WriteString(`<meta property="og:url" content="`)
		builder.WriteString(html.EscapeString(metadata.Canonical))
		builder.WriteString(`">`)
	}
	if strings.TrimSpace(metadata.Image) != "" {
		builder.WriteString(`<meta property="og:image" content="`)
		builder.WriteString(html.EscapeString(metadata.Image))
		builder.WriteString(`">`)
	}
	if metadata.Type == "article" {
		if strings.TrimSpace(metadata.PublishedAt) != "" {
			builder.WriteString(`<meta property="article:published_time" content="`)
			builder.WriteString(html.EscapeString(metadata.PublishedAt))
			builder.WriteString(`">`)
		}
		if strings.TrimSpace(metadata.ModifiedAt) != "" {
			builder.WriteString(`<meta property="article:modified_time" content="`)
			builder.WriteString(html.EscapeString(metadata.ModifiedAt))
			builder.WriteString(`">`)
		}
	}
	builder.WriteString(`<meta name="twitter:card" content="`)
	if strings.TrimSpace(metadata.Image) != "" {
		builder.WriteString(`summary_large_image`)
	} else {
		builder.WriteString(`summary`)
	}
	builder.WriteString(`"><meta name="twitter:title" content="`)
	builder.WriteString(html.EscapeString(pageTitle(site, metadata.Title)))
	builder.WriteString(`"><meta name="twitter:description" content="`)
	builder.WriteString(html.EscapeString(metadata.Description))
	builder.WriteString(`">`)
	if strings.TrimSpace(metadata.Image) != "" {
		builder.WriteString(`<meta name="twitter:image" content="`)
		builder.WriteString(html.EscapeString(metadata.Image))
		builder.WriteString(`">`)
	}
	builder.WriteString(renderStructuredData(site, metadata))
	builder.WriteString(`<style>`)
	builder.WriteString(renderStyles(site.Theme))
	builder.WriteString(`</style></head><body class="`)
	builder.WriteString(site.Theme.LayoutClass)
	builder.WriteString(`">`)
	builder.WriteString(`<main class="page-shell">`)
	builder.WriteString(body)
	builder.WriteString(`</main></body></html>`)
	return builder.String()
}

func siteDocumentMetadata(site renderedSite, title, canonical string) documentMetadata {
	return documentMetadata{
		Title:       title,
		Description: site.Description,
		Canonical:   canonical,
		Type:        "website",
	}
}

func articleDocumentMetadata(site renderedSite, article ArticleContent) documentMetadata {
	title := strings.TrimSpace(article.SEOTitle)
	if title == "" {
		title = article.Title
	}
	description := strings.TrimSpace(article.SEODescription)
	if description == "" {
		description = strings.TrimSpace(article.Excerpt)
	}
	if description == "" {
		description = site.Description
	}
	return documentMetadata{
		Title:       title,
		Description: description,
		Canonical:   canonicalURL(site, article),
		Image:       absolutePublicURL(site, article.CoverImageURL),
		Type:        "article",
		PublishedAt: article.PublishedAt,
		ModifiedAt:  article.UpdatedAt,
		Author:      article.AuthorName,
		Tags:        article.Tags,
	}
}

func absolutePublicURL(site renderedSite, value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "/") {
		return site.PublicBaseURL + value
	}
	return value
}

func renderStructuredData(site renderedSite, metadata documentMetadata) string {
	data := map[string]any{
		"@context":    "https://schema.org",
		"@type":       "WebSite",
		"name":        site.Title,
		"description": metadata.Description,
		"url":         metadata.Canonical,
	}
	if metadata.Type == "article" {
		data["@type"] = "BlogPosting"
		data["headline"] = metadata.Title
		data["mainEntityOfPage"] = map[string]string{"@type": "WebPage", "@id": metadata.Canonical}
		if metadata.Image != "" {
			data["image"] = []string{metadata.Image}
		}
		if metadata.PublishedAt != "" {
			data["datePublished"] = metadata.PublishedAt
		}
		if metadata.ModifiedAt != "" {
			data["dateModified"] = metadata.ModifiedAt
		}
		if metadata.Author != "" {
			data["author"] = map[string]string{"@type": "Person", "name": metadata.Author}
		}
		if len(metadata.Tags) > 0 {
			keywords := make([]string, 0, len(metadata.Tags))
			for _, tag := range metadata.Tags {
				if name := strings.TrimSpace(tag.Name); name != "" {
					keywords = append(keywords, name)
				}
			}
			if len(keywords) > 0 {
				data["keywords"] = strings.Join(keywords, ", ")
			}
		}
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return `<script type="application/ld+json">` + string(encoded) + `</script>`
}

func pageTitle(site renderedSite, title string) string {
	if strings.TrimSpace(title) == "" {
		return site.Title
	}
	return title + " | " + site.Title
}

func renderStyles(theme themeVariant) string {
	styles := strings.NewReplacer(
		"{{background}}", theme.Background,
		"{{surface}}", theme.Surface,
		"{{surfaceStrong}}", theme.SurfaceStrong,
		"{{text}}", theme.Text,
		"{{muted}}", theme.Muted,
		"{{accent}}", theme.Accent,
		"{{accentContrast}}", theme.AccentContrast,
		"{{border}}", theme.Border,
		"{{secondaryAccent}}", theme.SecondaryAccent,
		"{{featureBadgeBackground}}", theme.FeatureBadgeBackground,
	).Replace(baseStyles())
	if theme.Name == "anonime" {
		styles += strings.ReplaceAll(
			anonimeStyles(),
			`body.anonime-layout:has(#anonime-theme-toggle:checked)`,
			`html[data-theme="dark"] body.anonime-layout`,
		)
	}
	return styles
}

func baseStyles() string {
	return `
		:root {
			color-scheme: light;
		}
		body {
			margin: 0;
			font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			background: linear-gradient(180deg, {{background}} 0%, {{surfaceStrong}} 100%);
			color: {{text}};
		}
		a {
			color: inherit;
			text-decoration: none;
		}
		.page-shell {
			width: min(1120px, calc(100vw - 32px));
			margin: 0 auto;
			padding: 32px 0 64px;
		}
		.hero,
		.panel,
		.article-layout {
			background: {{surface}};
			border: 1px solid {{border}};
			border-radius: 24px;
			box-shadow: 0 18px 48px rgba(15, 23, 42, 0.08);
		}
		.site-logo { display: block; max-width: 180px; max-height: 64px; object-fit: contain; margin-bottom: 20px; }
		.hero {
			padding: 40px;
		}
		.hero.compact {
			padding: 32px 40px;
		}
		.eyebrow {
			margin: 0 0 10px;
			font-size: 0.78rem;
			font-weight: 800;
			letter-spacing: 0.24em;
			text-transform: uppercase;
			color: {{accent}};
		}
		h1, h2, h3, p {
			margin-top: 0;
		}
		h1 {
			font-size: clamp(2.4rem, 5vw, 4.4rem);
			line-height: 0.95;
			margin-bottom: 16px;
		}
		h2 {
			font-size: 1.25rem;
			margin-bottom: 20px;
		}
		h3 {
			font-size: 1.1rem;
			margin-bottom: 10px;
		}
		.lead {
			font-size: 1.05rem;
			line-height: 1.7;
			color: {{muted}};
			max-width: 58ch;
		}
		.notice {
			display: inline-flex;
			padding: 10px 14px;
			border-radius: 999px;
			background: {{featureBadgeBackground}};
			color: {{accent}};
			font-weight: 700;
		}
		.hero-actions {
			margin-top: 24px;
			display: flex;
			gap: 12px;
			flex-wrap: wrap;
		}
		.button {
			display: inline-flex;
			align-items: center;
			justify-content: center;
			gap: 8px;
			padding: 14px 20px;
			border-radius: 16px;
			font-weight: 700;
		}
		.button.primary {
			background: {{accent}};
			color: {{accentContrast}};
		}
		.grid {
			margin-top: 24px;
			display: grid;
			grid-template-columns: 1fr;
			gap: 20px;
		}
		.panel {
			padding: 28px;
		}
		.section-list,
		.article-list,
		.article-chip-row,
		.article-meta {
			display: flex;
			flex-wrap: wrap;
			gap: 12px;
		}
		.section-card,
		.article-card {
			background: {{surfaceStrong}};
			border: 1px solid {{border}};
			border-radius: 20px;
			padding: 18px;
		}
		.article-card {
			display: grid;
			grid-template-columns: 160px minmax(0, 1fr);
			gap: 18px;
			align-items: stretch;
		}
		.article-card.featured {
			border-color: {{accent}};
			box-shadow: inset 0 0 0 1px {{accent}};
		}
		.article-image {
			width: 100%;
			height: 100%;
			min-height: 150px;
			object-fit: cover;
			border-radius: 16px;
			background: {{surface}};
		}
		.article-card-content {
			min-width: 0;
		}
		.article-status {
			display: inline-flex;
			padding: 6px 10px;
			margin-bottom: 10px;
			font-size: 0.72rem;
			font-weight: 800;
			letter-spacing: 0.12em;
			text-transform: uppercase;
			border-radius: 999px;
			background: {{featureBadgeBackground}};
			color: {{accent}};
		}
		.article-chip-row span,
		.article-meta span,
		.tag {
			display: inline-flex;
			align-items: center;
			padding: 8px 12px;
			border-radius: 999px;
			background: {{surface}};
			border: 1px solid {{border}};
			color: {{muted}};
			font-size: 0.88rem;
		}
		.article-layout {
			padding: 36px;
		}
		.article-hero {
			margin-bottom: 32px;
		}
		.article-body {
			font-size: 1.03rem;
			line-height: 1.8;
			color: {{text}};
		}
		.article-body img {
			max-width: 100%;
		}
		.cover img {
			width: 100%;
			height: auto;
			border-radius: 20px;
			border: 1px solid {{border}};
		}
		.muted {
			color: {{muted}};
		}
		.tag-list {
			margin-top: 28px;
			display: flex;
			flex-wrap: wrap;
			gap: 10px;
		}
		@media (max-width: 720px) {
			.hero,
			.panel,
			.article-layout {
				padding: 22px;
				border-radius: 20px;
			}
			.article-card {
				grid-template-columns: 1fr;
			}
		}
	`
}

func renderMarkdown(markdownContent string) string {
	if strings.TrimSpace(markdownContent) == "" {
		return `<p class="muted">No article body yet.</p>`
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdownContent), &buf); err != nil {
		return "<pre>" + html.EscapeString(markdownContent) + "</pre>"
	}
	return buf.String()
}

func articleURL(site renderedSite, article ArticleContent) string {
	return site.PublicBaseURL + articlePath(site, article)
}

func articlesURL(site renderedSite) string {
	return site.PublicBaseURL + articlesPath(site)
}

func articlesPath(site renderedSite) string {
	if site.Theme.Name == "anonime" {
		return strings.TrimRight(site.BasePath, "/") + "/articles/"
	}
	return site.BasePath + "/"
}

func anonimeBlogLandingPath(site renderedSite) string {
	return site.BasePath + "/"
}

func articlePath(site renderedSite, article ArticleContent) string {
	return site.BasePath + "/" + safePathSegment(article.Slug) + "/"
}

func canonicalURL(site renderedSite, article ArticleContent) string {
	if strings.TrimSpace(article.CanonicalURL) != "" {
		return article.CanonicalURL
	}
	return articleURL(site, article)
}

func renderSitemap(site renderedSite) string {
	urls := sitemapURLs(site, true)
	return renderSitemapURLs(urls)
}

func renderBlogSitemap(site renderedSite) string {
	return renderSitemapURLs(sitemapURLs(site, false))
}

func sitemapURLs(site renderedSite, includeHome bool) []sitemapURL {
	urls := make([]sitemapURL, 0, len(site.Articles)+4)
	if includeHome {
		urls = append(urls, sitemapURL{Loc: site.PublicBaseURL + "/"})
	}
	urls = append(urls, sitemapURL{Loc: site.PublicBaseURL + site.BasePath + "/"})
	if site.Theme.Name == "anonime" {
		urls = append(urls, sitemapURL{Loc: articlesURL(site)})
		for page := 2; page <= anonimeArticlePageCount(len(site.Articles)); page++ {
			urls = append(urls, sitemapURL{Loc: anonimeArticlesPageURL(site, page)})
		}
	}
	for _, article := range site.Articles {
		urls = append(urls, sitemapURL{Loc: articleURL(site, article), LastMod: article.UpdatedAt})
	}
	return urls
}

func renderSitemapURLs(urls []sitemapURL) string {
	doc := urlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		return xml.Header + "<urlset></urlset>"
	}
	return buf.String()
}

func renderRobots(site renderedSite) string {
	var builder strings.Builder
	builder.WriteString("User-agent: *\n")
	builder.WriteString("Allow: /\n")
	builder.WriteString("Sitemap: ")
	builder.WriteString(site.PublicBaseURL)
	builder.WriteString("/sitemap.xml\n")
	return builder.String()
}

func renderRSS(site renderedSite) string {
	items := make([]rssItem, 0, len(site.Articles))
	for _, article := range site.Articles {
		items = append(items, rssItem{
			Title:       article.Title,
			Link:        articleURL(site, article),
			Description: article.Excerpt,
			PubDate:     article.PublishedAt,
			GUID:        articleURL(site, article),
		})
	}

	doc := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         site.Title,
			Link:          site.PublicBaseURL + site.BasePath + "/",
			Description:   site.Description,
			Language:      "en",
			LastBuildDate: site.Now.Format(time.RFC1123Z),
			Items:         items,
		},
	}
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		return xml.Header + "<rss version=\"2.0\"></rss>"
	}
	return buf.String()
}

func safePathSegment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "index"
	}
	value = strings.ToLower(value)
	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_':
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		default:
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(builder.String(), "-")
}

func renderJSON(value map[string]any) string {
	if len(value) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var builder strings.Builder
	builder.WriteString("{")
	for index, key := range keys {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(key)
		builder.WriteString(": ")
		switch entry := value[key].(type) {
		case string:
			builder.WriteString(strconv.Quote(entry))
		default:
			builder.WriteString(fmt.Sprint(entry))
		}
	}
	builder.WriteString("}")
	return builder.String()
}
