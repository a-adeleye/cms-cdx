package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"cms-builder/api/internal/models"
)

type GenerateOptions struct {
	SiteID    string
	BuildID   string
	Preview   bool
	BuildData []byte
}

type Builder interface {
	GenerateSite(ctx context.Context, site models.Site, options GenerateOptions) (string, error)
}

type NoopBuilder struct{}

func (NoopBuilder) GenerateSite(ctx context.Context, site models.Site, options GenerateOptions) (string, error) {
	return "dist/sites/" + site.Slug, nil
}

type AstroBuilder struct {
	Directory  string
	OutputRoot string
	NPMCommand string
}

func (b AstroBuilder) GenerateSite(ctx context.Context, site models.Site, options GenerateOptions) (string, error) {
	if len(options.BuildData) == 0 {
		return "", fmt.Errorf("build data is required")
	}

	outputPath := filepath.Join(b.OutputRoot, site.ID, options.BuildID)
	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return "", fmt.Errorf("create build output directory: %w", err)
	}

	dataPath := filepath.Join(b.OutputRoot, ".build-data", options.BuildID+".json")
	if err := os.MkdirAll(filepath.Dir(dataPath), 0o700); err != nil {
		return "", fmt.Errorf("create build data directory: %w", err)
	}
	if err := os.WriteFile(dataPath, options.BuildData, 0o600); err != nil {
		return "", fmt.Errorf("write build data: %w", err)
	}
	defer os.Remove(dataPath)

	buildCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	command := exec.CommandContext(buildCtx, b.NPMCommand, "--prefix", b.Directory, "run", "build")
	command.Env = buildEnvironment(dataPath, outputPath, siteURL(site))
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("generate static site: %s", strings.TrimSpace(string(output)))
	}

	if _, err := os.Stat(filepath.Join(outputPath, "index.html")); err != nil {
		return "", fmt.Errorf("generated site is missing index.html: %w", err)
	}
	return outputPath, nil
}

func siteURL(site models.Site) string {
	if strings.HasPrefix(site.Domain, "https://") || strings.HasPrefix(site.Domain, "http://") {
		return site.Domain
	}
	if site.Domain != "" {
		return "https://" + site.Domain
	}
	return "https://" + site.Slug + ".pages.dev"
}

func buildEnvironment(dataPath, outputPath, url string) []string {
	environment := make([]string, 0, 6)
	for _, key := range []string{"PATH", "HOME", "TMPDIR"} {
		if value, ok := os.LookupEnv(key); ok {
			environment = append(environment, key+"="+value)
		}
	}
	return append(environment, "CMS_BUILD_DATA_FILE="+dataPath, "CMS_BUILD_OUTPUT_DIR="+outputPath, "CMS_SITE_URL="+url)
}
