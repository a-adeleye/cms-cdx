package builder

import (
	"testing"

	"cms-builder/api/internal/models"
)

func TestSiteURLUsesConfiguredDomain(t *testing.T) {
	url := siteURL(models.Site{Domain: "blog.example.com"})
	if url != "https://blog.example.com" {
		t.Fatalf("expected https://blog.example.com, got %q", url)
	}
}

func TestSiteURLUsesPagesDomainWhenSiteDomainIsEmpty(t *testing.T) {
	url := siteURL(models.Site{Slug: "example-site"})
	if url != "https://example-site.pages.dev" {
		t.Fatalf("expected Pages fallback URL, got %q", url)
	}
}

func TestBuildEnvironmentDoesNotInheritAPISecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "jwt-secret")
	t.Setenv("PATH", "test-path")

	environment := buildEnvironment("data.json", "output", "https://example.com")
	for _, value := range environment {
		if value == "JWT_SECRET=jwt-secret" {
			t.Fatal("expected unrelated API secret to be excluded from builder environment")
		}
	}
}
