package deploy

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"cms-builder/api/internal/models"
)

type RepositoryAdapter struct {
	GitCommand             string
	AllowedHosts           []string
	SecretResolver         func(string) (string, error)
	AllowLocalRepositories bool
}

type repositoryConfig struct {
	RepositoryURL  string
	Branch         string
	ContentPath    string
	TokenSecretRef string
	PublicURL      string
}

func NewRepositoryAdapter() RepositoryAdapter {
	return RepositoryAdapter{
		GitCommand:     "git",
		AllowedHosts:   []string{"github.com", "gitlab.com"},
		SecretResolver: resolveRepositorySecret,
	}
}

func (a RepositoryAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	provider, values := providerConfigForBuild(site, build)
	if !strings.EqualFold(provider, "git_repository") {
		return nil, fmt.Errorf("repository adapter received unsupported provider %q", provider)
	}
	config, err := a.parseConfig(values)
	if err != nil {
		return nil, err
	}
	blogPath, err := models.CanonicalBlogPath(site.BlogPath)
	if err != nil {
		return nil, err
	}
	source := filepath.Join(outputPath, filepath.FromSlash(strings.TrimPrefix(blogPath, "/")))
	if info, statErr := os.Stat(source); statErr != nil || !info.IsDir() {
		return nil, errors.New("generated blog directory is unavailable")
	}

	deployCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	workRoot, err := os.MkdirTemp("", "cms-repository-deploy-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workRoot)
	repositoryRoot := filepath.Join(workRoot, "repository")
	env, err := a.gitEnvironment(workRoot, config.TokenSecretRef)
	if err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, "", "clone", "--depth", "1", "--single-branch", "--branch", config.Branch, config.RepositoryURL, repositoryRoot); err != nil {
		return nil, err
	}
	target, err := confinedRepositoryPath(repositoryRoot, config.ContentPath)
	if err != nil {
		return nil, err
	}
	if err := os.RemoveAll(target); err != nil {
		return nil, err
	}
	if err := copyTree(source, target); err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "config", "user.name", "CMS Publisher"); err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "config", "user.email", "cms-publisher@localhost"); err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "add", "--", filepath.ToSlash(config.ContentPath)); err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "diff", "--cached", "--quiet"); err == nil {
		revision, revErr := a.runGit(deployCtx, env, repositoryRoot, "rev-parse", "HEAD")
		if revErr != nil {
			return nil, revErr
		}
		return &DeployResult{Provider: "git_repository", URL: repositoryDeploymentURL(site, config), Revision: strings.TrimSpace(revision), Message: "repository already contains the generated blog"}, nil
	}
	message := fmt.Sprintf("Publish %s CMS build %s", build.BuildType, build.ID)
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "commit", "-m", message); err != nil {
		return nil, err
	}
	revision, err := a.runGit(deployCtx, env, repositoryRoot, "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	if _, err := a.runGit(deployCtx, env, repositoryRoot, "push", "origin", "HEAD:refs/heads/"+config.Branch); err != nil {
		return nil, err
	}
	return &DeployResult{Provider: "git_repository", URL: repositoryDeploymentURL(site, config), Revision: strings.TrimSpace(revision), Message: "generated blog committed and pushed to " + config.Branch}, nil
}

func (a RepositoryAdapter) parseConfig(values map[string]any) (repositoryConfig, error) {
	config := repositoryConfig{
		RepositoryURL: configString(values, "repositoryUrl"), Branch: configString(values, "branch"),
		ContentPath: configString(values, "contentPath"), TokenSecretRef: configString(values, "tokenSecretRef"),
		PublicURL: configString(values, "publicUrl"),
	}
	if config.Branch == "" {
		config.Branch = "main"
	}
	if config.ContentPath == "" {
		config.ContentPath = "public/blog"
	}
	parsed, err := url.Parse(config.RepositoryURL)
	localAllowed := a.AllowLocalRepositories && parsed != nil && parsed.Scheme == "file"
	if err != nil || (!localAllowed && (parsed.Scheme != "https" || parsed.User != nil || parsed.Hostname() == "")) {
		return repositoryConfig{}, errors.New("repositoryUrl must be an HTTPS URL without embedded credentials")
	}
	allowed := localAllowed
	for _, host := range a.AllowedHosts {
		if strings.EqualFold(strings.TrimSpace(host), parsed.Hostname()) {
			allowed = true
			break
		}
	}
	if !allowed {
		return repositoryConfig{}, fmt.Errorf("repository host %s is not allowed", parsed.Hostname())
	}
	if strings.ContainsAny(config.Branch, "\r\n ~^:?*[\\") || strings.HasPrefix(config.Branch, "-") {
		return repositoryConfig{}, errors.New("repository branch is invalid")
	}
	if _, err := safeRelativePath(config.ContentPath); err != nil {
		return repositoryConfig{}, err
	}
	return config, nil
}

func safeRelativePath(value string) (string, error) {
	cleaned := filepath.Clean(filepath.FromSlash(strings.TrimSpace(value)))
	if cleaned == "." || filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || strings.Contains(cleaned, ".git") {
		return "", errors.New("contentPath must be a safe relative directory outside .git")
	}
	return cleaned, nil
}

func confinedRepositoryPath(root, relative string) (string, error) {
	cleaned, err := safeRelativePath(relative)
	if err != nil {
		return "", err
	}
	target := filepath.Join(root, cleaned)
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("contentPath escapes the repository")
	}
	return target, nil
}

func (a RepositoryAdapter) gitEnvironment(workRoot, secretRef string) ([]string, error) {
	env := append([]string{}, os.Environ()...)
	env = append(env, "GIT_TERMINAL_PROMPT=0")
	if strings.TrimSpace(secretRef) == "" {
		return env, nil
	}
	resolver := a.SecretResolver
	if resolver == nil {
		resolver = resolveRepositorySecret
	}
	token, err := resolver(secretRef)
	if err != nil {
		return nil, err
	}
	name, contents := "askpass.sh", "#!/bin/sh\ncase \"$1\" in *Username*) echo x-access-token;; *) echo \"$CMS_GIT_TOKEN\";; esac\n"
	if runtime.GOOS == "windows" {
		name, contents = "askpass.bat", "@echo off\r\necho %CMS_GIT_TOKEN%\r\n"
	}
	path := filepath.Join(workRoot, name)
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		return nil, err
	}
	return append(env, "GIT_ASKPASS="+path, "CMS_GIT_TOKEN="+token), nil
}

func (a RepositoryAdapter) runGit(ctx context.Context, env []string, directory string, arguments ...string) (string, error) {
	command := a.GitCommand
	if strings.TrimSpace(command) == "" {
		command = "git"
	}
	cmd := exec.CommandContext(ctx, command, arguments...)
	cmd.Dir, cmd.Env = directory, env
	output, err := cmd.CombinedOutput()
	message := strings.TrimSpace(string(output))
	if err != nil {
		return "", fmt.Errorf("git %s failed: %s", arguments[0], message)
	}
	return message, nil
}

func resolveRepositorySecret(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", errors.New("repository token secret reference is required")
	}
	value := strings.TrimSpace(os.Getenv(ref))
	if value == "" {
		return "", fmt.Errorf("repository token secret %s is not set", ref)
	}
	return value, nil
}

func repositoryDeploymentURL(site models.Site, config repositoryConfig) string {
	if config.PublicURL != "" {
		return strings.TrimRight(config.PublicURL, "/")
	}
	return strings.TrimRight(site.Domain, "/") + site.BlogPath + "/"
}
