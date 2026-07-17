package handlers

import (
	"strings"
	"testing"
)

func TestValidateSitePayloadRejectsCloudflarePagesPreviewWithoutProjectName(t *testing.T) {
	err := validateSitePayload(siteUpsertRequest{
		Name:                  "Example Site",
		Slug:                  "example-site",
		BlogPath:              "/articles",
		DeployProvider:        "none",
		PreviewDeployProvider: "cloudflare_pages",
		PreviewDeployConfig:   `{}`,
	})
	if err == nil {
		t.Fatal("expected Cloudflare Pages preview config without a project name to be rejected")
	}
	if !strings.Contains(err.Error(), "Cloudflare deployment config requires a valid projectName") {
		t.Fatalf("expected missing project name error, got %v", err)
	}
}

func TestValidateSitePayloadAcceptsCloudflarePagesPreviewWithCustomBlogPath(t *testing.T) {
	err := validateSitePayload(siteUpsertRequest{
		Name:                  "Example Site",
		Slug:                  "example-site",
		BlogPath:              "/blog",
		DeployProvider:        "none",
		PreviewDeployProvider: "cloudflare_pages",
		PreviewDeployConfig:   `{"projectName":"example-site","productionBranch":"main"}`,
	})
	if err != nil {
		t.Fatalf("expected valid Cloudflare Pages preview config, got %v", err)
	}
}

func TestValidateSitePayloadRejectsCloudflarePagesPreviewSetupWithoutProject(t *testing.T) {
	err := validateSitePayload(siteUpsertRequest{
		Name:                  "Example Site",
		Slug:                  "example-site",
		BlogPath:              "/articles",
		DeployProvider:        "none",
		PreviewDeployProvider: "cloudflare_pages",
		PreviewDeployConfig:   `{"projectName":"","productionBranch":"main"}`,
	})
	if err == nil {
		t.Fatal("expected an incomplete Cloudflare Pages setup to be rejected")
	}
}

func TestCloudflarePagesDeploymentWarningsReportIncompleteSetup(t *testing.T) {
	t.Setenv("CLOUDFLARE_API_TOKEN", "")
	t.Setenv("CLOUDFLARE_ACCOUNT_ID", "")

	warnings := cloudflarePagesDeploymentWarnings("preview", "cloudflare_pages", `{"projectName":""}`)
	if len(warnings) != 2 {
		t.Fatalf("expected two Cloudflare Pages setup warnings, got %#v", warnings)
	}
	if !strings.Contains(warnings[0], "projectName") {
		t.Fatalf("expected a project name warning, got %q", warnings[0])
	}
	if !strings.Contains(warnings[1], "CLOUDFLARE_API_TOKEN") {
		t.Fatalf("expected a credentials warning, got %q", warnings[1])
	}
}

func TestValidateSitePayloadRejectsUnsupportedContentContext(t *testing.T) {
	err := validateSitePayload(siteUpsertRequest{
		Name:           "Example Site",
		Slug:           "example-site",
		BlogPath:       "/blog",
		ContentContext: "marketing_site",
	})
	if err == nil {
		t.Fatal("expected unsupported content context to be rejected")
	}
	if !strings.Contains(err.Error(), "contentContext") {
		t.Fatalf("expected a content context validation error, got %v", err)
	}
}

func TestSiteContentContextDefaultsToStandaloneBlog(t *testing.T) {
	if context := siteContentContext(""); context != "standalone_blog" {
		t.Fatalf("expected standalone_blog default, got %q", context)
	}
}
