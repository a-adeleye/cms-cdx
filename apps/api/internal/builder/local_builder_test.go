package builder

import (
	"strings"
	"testing"

	"cms-builder/api/internal/models"
)

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
}
