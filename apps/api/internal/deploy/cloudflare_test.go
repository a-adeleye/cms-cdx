package deploy

import "testing"

import "cms-builder/api/internal/models"

func TestPagesConfigAcceptsProjectAndDefaultBranch(t *testing.T) {
	projectName, branch, err := pagesConfig(map[string]any{"projectName": "example-site"})
	if err != nil {
		t.Fatalf("pagesConfig returned error: %v", err)
	}
	if projectName != "example-site" {
		t.Fatalf("expected project name example-site, got %q", projectName)
	}
	if branch != "main" {
		t.Fatalf("expected default branch main, got %q", branch)
	}
}

func TestPagesConfigRejectsInvalidProjectName(t *testing.T) {
	_, _, err := pagesConfig(map[string]any{"projectName": "invalid project"})
	if err == nil {
		t.Fatal("expected invalid project name to be rejected")
	}
}

func TestRedactTokenRemovesCredentialFromDeploymentOutput(t *testing.T) {
	message := redactToken("upload failed for token secret-token", "secret-token")
	if message != "upload failed for token [REDACTED]" {
		t.Fatalf("expected token to be redacted, got %q", message)
	}
}

func TestChildEnvironmentDoesNotInheritUnrelatedSecrets(t *testing.T) {
	t.Setenv("JWT_SECRET", "jwt-secret")
	t.Setenv("PATH", "test-path")

	environment := childEnvironment("CLOUDFLARE_API_TOKEN=cloudflare-token")
	for _, value := range environment {
		if value == "JWT_SECRET=jwt-secret" {
			t.Fatal("expected unrelated API secret to be excluded from child environment")
		}
	}
}

func TestNonCloudflareProviderDoesNotReportDeployment(t *testing.T) {
	result, err := (CloudflarePagesAdapter{}).Deploy(t.Context(), models.Site{DeployProvider: "firebase"}, models.Build{}, "")
	if err != nil {
		t.Fatalf("Deploy returned error: %v", err)
	}
	if result.Provider != "none" {
		t.Fatalf("expected a skipped deployment, got provider %q", result.Provider)
	}
}
