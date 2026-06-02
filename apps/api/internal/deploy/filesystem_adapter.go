package deploy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"cms-builder/api/internal/models"
)

type FilesystemAdapter struct {
	DeployRoot string
}

func NewFilesystemAdapter(deployRoot string) FilesystemAdapter {
	if strings.TrimSpace(deployRoot) == "" {
		deployRoot = filepath.Join("dist", "deployments")
	}
	return FilesystemAdapter{DeployRoot: deployRoot}
}

func (a FilesystemAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(outputPath) == "" {
		return nil, errors.New("output path is required")
	}
	if strings.TrimSpace(site.Slug) == "" {
		return nil, errors.New("site slug is required")
	}

	provider := fallbackProvider(site.DeployProvider)
	targetDir := filepath.Join(a.DeployRoot, provider, safePathSegment(site.Slug), safePathSegment(build.BuildType))
	if err := os.RemoveAll(targetDir); err != nil {
		return nil, err
	}
	if err := copyTree(outputPath, targetDir); err != nil {
		return nil, err
	}

	return &DeployResult{
		Provider: provider,
		URL:      deploymentURL(site, build, provider),
		Message:  fmt.Sprintf("deployed to %s", targetDir),
	}, nil
}

func deploymentURL(site models.Site, build models.Build, provider string) string {
	base := ""
	if publicURL, ok := site.DeployConfig["publicUrl"].(string); ok && strings.TrimSpace(publicURL) != "" {
		base = strings.TrimRight(strings.TrimSpace(publicURL), "/")
	}
	if base == "" {
		base = "http://localhost:8081/deployments"
	}
	return fmt.Sprintf("%s/%s/%s/%s/", base, safePathSegment(provider), safePathSegment(site.Slug), safePathSegment(build.BuildType))
}

func copyTree(source, target string) error {
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}

	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}

		destination := filepath.Join(target, relative)
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}

		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}

		input, err := os.Open(path)
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := os.Create(destination)
		if err != nil {
			return err
		}

		if _, err := io.Copy(output, input); err != nil {
			output.Close()
			return err
		}
		return output.Close()
	})
}

func safePathSegment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "default"
	}
	value = strings.ToLower(value)
	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_':
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		default:
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(builder.String(), "-")
}
