package deploy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"cms-builder/api/internal/models"
)

var pagesProjectName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,58}[a-z0-9]$|^[a-z0-9]$`)
var pagesURL = regexp.MustCompile(`https://[^\s]+\.pages\.dev[^\s]*`)

type CloudflarePagesAdapter struct {
	APIToken  string
	AccountID string
	Command   string
}

func (a CloudflarePagesAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	if site.DeployProvider == "" || site.DeployProvider == "none" {
		return &DeployResult{Provider: "none", Message: "deployment not configured"}, nil
	}
	if site.DeployProvider != "cloudflare_pages" {
		return &DeployResult{Provider: "none", Message: "deployment is managed by the configured provider"}, nil
	}
	if a.APIToken == "" || a.AccountID == "" {
		return nil, fmt.Errorf("Cloudflare Pages requires CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID")
	}
	if info, err := os.Stat(outputPath); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("deployment output directory is unavailable")
	}

	projectName, productionBranch, err := pagesConfig(site.DeployConfig)
	if err != nil {
		return nil, err
	}

	branch := productionBranch
	if build.BuildType == "preview" {
		branch = "cms-preview"
	}
	deployCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	command := exec.CommandContext(deployCtx, a.Command, "pages", "deploy", filepath.Clean(outputPath), "--project-name", projectName, "--branch", branch, "--commit-message", "CMS build "+build.ID)
	command.Env = childEnvironment("CLOUDFLARE_API_TOKEN="+a.APIToken, "CLOUDFLARE_ACCOUNT_ID="+a.AccountID)
	output, err := command.CombinedOutput()
	message := redactToken(strings.TrimSpace(string(output)), a.APIToken)
	if err != nil {
		return nil, fmt.Errorf("Cloudflare Pages deployment failed: %s", message)
	}

	url := deploymentURL(message)
	if url == "" && build.BuildType != "preview" {
		url = "https://" + projectName + ".pages.dev"
	}
	if url == "" {
		return nil, fmt.Errorf("Cloudflare Pages deployment completed without a preview URL")
	}
	if err := verifyDeployment(deployCtx, url, outputPath); err != nil {
		return nil, err
	}
	return &DeployResult{Provider: "cloudflare_pages", URL: url, Message: message + "\nPost-deploy verification passed."}, nil
}

func pagesConfig(config map[string]any) (string, string, error) {
	projectName, _ := config["projectName"].(string)
	projectName = strings.TrimSpace(projectName)
	if !pagesProjectName.MatchString(projectName) {
		return "", "", fmt.Errorf("deployment config requires a valid projectName")
	}
	productionBranch, _ := config["productionBranch"].(string)
	productionBranch = strings.TrimSpace(productionBranch)
	if productionBranch == "" {
		productionBranch = "main"
	}
	if strings.ContainsAny(productionBranch, "\r\n") {
		return "", "", fmt.Errorf("deployment config has an invalid productionBranch")
	}
	return projectName, productionBranch, nil
}

func deploymentURL(output string) string {
	return pagesURL.FindString(output)
}

func redactToken(output, token string) string {
	return strings.ReplaceAll(output, token, "[REDACTED]")
}

func childEnvironment(values ...string) []string {
	environment := make([]string, 0, len(values)+3)
	for _, key := range []string{"PATH", "HOME", "TMPDIR"} {
		if value, ok := os.LookupEnv(key); ok {
			environment = append(environment, key+"="+value)
		}
	}
	return append(environment, values...)
}

func verifyDeployment(ctx context.Context, deploymentURL, outputPath string) error {
	parsedURL, err := url.Parse(deploymentURL)
	if err != nil || parsedURL.Scheme != "https" || !strings.HasSuffix(parsedURL.Hostname(), ".pages.dev") {
		return fmt.Errorf("Cloudflare Pages returned an invalid deployment URL")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	paths := []string{"/", "/articles/"}
	if articlePath, err := firstArticlePath(outputPath); err != nil {
		return err
	} else if articlePath != "" {
		paths = append(paths, articlePath)
	}
	for _, path := range paths {
		if err := verifyPath(ctx, client, deploymentURL+path); err != nil {
			return err
		}
	}
	return nil
}

func firstArticlePath(outputPath string) (string, error) {
	entries, err := os.ReadDir(filepath.Join(outputPath, "articles"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("inspect generated article pages: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if _, err := os.Stat(filepath.Join(outputPath, "articles", entry.Name(), "index.html")); err == nil {
				return "/articles/" + entry.Name() + "/", nil
			}
		}
	}
	return "", nil
}

func verifyPath(ctx context.Context, client *http.Client, requestURL string) error {
	var lastError error
	for attempt := 0; attempt < 3; attempt++ {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return fmt.Errorf("create post-deploy verification request: %w", err)
		}
		response, err := client.Do(request)
		if err == nil && response.StatusCode == http.StatusOK {
			response.Body.Close()
			return nil
		}
		if response != nil {
			response.Body.Close()
		}
		if err != nil {
			lastError = err
		} else {
			lastError = fmt.Errorf("received HTTP %d", response.StatusCode)
		}
		if attempt < 2 {
			select {
			case <-ctx.Done():
				return fmt.Errorf("verify deployed site: %w", ctx.Err())
			case <-time.After(time.Duration(attempt+1) * time.Second):
			}
		}
	}
	return fmt.Errorf("verify deployed site %s: %w", requestURL, lastError)
}
