package handlers

import (
	"strings"
	"testing"
)

func TestValidateDeploymentConfigRequiresRepositoryTokenReference(t *testing.T) {
	err := validateDeploymentConfig("git_repository", `{
		"repositoryUrl": "https://github.com/example/site.git",
		"branch": "main",
		"contentPath": "public/blog"
	}`)
	if err == nil || !strings.Contains(err.Error(), "tokenSecretRef") {
		t.Fatalf("expected repository deployment validation to require a token reference, got %v", err)
	}
}
