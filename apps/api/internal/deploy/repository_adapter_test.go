package deploy

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"cms-builder/api/internal/models"
)

func TestRepositoryAdapterPreservesLandingSiteAndReplacesOnlyBlogDirectory(t *testing.T) {
	root := t.TempDir()
	bare := filepath.Join(root, "remote.git")
	seed := filepath.Join(root, "seed")
	runTestGit(t, root, "init", "--bare", bare)
	runTestGit(t, root, "init", "-b", "main", seed)
	runTestGit(t, seed, "config", "user.name", "Test")
	runTestGit(t, seed, "config", "user.email", "test@example.com")
	writeTestFile(t, filepath.Join(seed, "index.html"), "landing page")
	writeTestFile(t, filepath.Join(seed, "public", "blog", "old.html"), "old blog")
	runTestGit(t, seed, "add", ".")
	runTestGit(t, seed, "commit", "-m", "seed")
	runTestGit(t, seed, "remote", "add", "origin", bare)
	runTestGit(t, seed, "push", "origin", "main")

	output := filepath.Join(root, "output")
	writeTestFile(t, filepath.Join(output, "blog", "index.html"), "new blog")
	writeTestFile(t, filepath.Join(output, "blog", "sitemap.xml"), "new sitemap")
	adapter := NewRepositoryAdapter()
	adapter.AllowLocalRepositories = true
	site := models.Site{Domain: "https://example.com", BlogPath: "/blog", DeployProvider: "git_repository", DeployConfig: map[string]any{
		"repositoryUrl": "file:///" + filepath.ToSlash(bare), "branch": "main", "contentPath": "public/blog",
	}}
	result, err := adapter.Deploy(context.Background(), site, models.Build{ID: "build-1", BuildType: "published"}, output)
	if err != nil {
		t.Fatalf("Deploy returned error: %v", err)
	}
	if result.Revision == "" {
		t.Fatal("expected pushed revision")
	}

	checkout := filepath.Join(root, "checkout")
	runTestGit(t, root, "clone", "--branch", "main", bare, checkout)
	assertTestFile(t, filepath.Join(checkout, "index.html"), "landing page")
	assertTestFile(t, filepath.Join(checkout, "public", "blog", "index.html"), "new blog")
	assertTestFile(t, filepath.Join(checkout, "public", "blog", "sitemap.xml"), "new sitemap")
	if _, err := os.Stat(filepath.Join(checkout, "public", "blog", "old.html")); !os.IsNotExist(err) {
		t.Fatalf("expected old blog file to be removed, got %v", err)
	}
}

func TestRepositoryConfigRejectsTraversalAndUnapprovedHosts(t *testing.T) {
	adapter := NewRepositoryAdapter()
	if _, err := adapter.parseConfig(map[string]any{"repositoryUrl": "https://example.com/site.git", "branch": "main", "contentPath": "public/blog"}); err == nil {
		t.Fatal("expected host allowlist error")
	}
	if _, err := adapter.parseConfig(map[string]any{"repositoryUrl": "https://github.com/example/site.git", "branch": "main", "contentPath": "../landing"}); err == nil {
		t.Fatal("expected traversal error")
	}
}

func TestRepositoryConfigRequiresTokenReferenceForHTTPS(t *testing.T) {
	adapter := NewRepositoryAdapter()
	_, err := adapter.parseConfig(map[string]any{
		"repositoryUrl": "https://github.com/example/site.git",
		"branch":        "main",
		"contentPath":   "public/blog",
	})
	if err == nil || !strings.Contains(err.Error(), "tokenSecretRef") {
		t.Fatalf("expected HTTPS repository configuration to require a token reference, got %v", err)
	}
}

func runTestGit(t *testing.T, directory string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = directory
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertTestFile(t *testing.T, path, expected string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != expected {
		t.Fatalf("expected %q, got %q", expected, contents)
	}
}
