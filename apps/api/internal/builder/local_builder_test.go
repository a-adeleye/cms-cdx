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
