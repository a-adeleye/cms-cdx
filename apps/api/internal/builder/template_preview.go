package builder

import (
	"time"

	"cms-builder/api/internal/models"
)

// RenderTemplatePreview uses the production renderer with deterministic sample
// content so the template catalogue never drifts into hand-drawn mockups.
func RenderTemplatePreview(templateKey string) string {
	siteModel := models.Site{
		Name: "Northstar Journal", Slug: "northstar-journal", Domain: "https://example.com",
		BlogPath: "/blog", Description: "Ideas, field notes, and practical guides for thoughtful teams.",
		TemplateKey: templateKey, ThemeConfig: map[string]any{"accent": "#2563eb"},
	}
	article := ArticleContent{
		ID: "preview-article", Title: "A calmer way to build on the web", Slug: "calmer-web",
		Excerpt: "A practical guide to making focused products that remain fast, useful, and easy to maintain.",
		Status:  "published", IsFeatured: true, AuthorName: "Northstar Editors", CategoryName: "Guides",
		PublishedAt: "2026-07-17T00:00:00Z", UpdatedAt: "2026-07-17T00:00:00Z",
	}
	rendered := renderedSite{
		Site: siteModel, Theme: themeForSite(siteModel), BuildType: "published", BasePath: "/blog",
		PublicBaseURL: siteModel.Domain, Title: siteModel.Name, Description: siteModel.Description,
		Articles: []ArticleContent{article}, FeaturedArticles: []ArticleContent{article}, RecentArticles: []ArticleContent{article},
		Now: time.Date(2026, time.July, 17, 0, 0, 0, 0, time.UTC),
	}
	return renderHomePage(rendered)
}
