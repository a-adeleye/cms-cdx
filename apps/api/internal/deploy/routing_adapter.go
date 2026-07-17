package deploy

import (
	"context"
	"path/filepath"
	"strings"

	"cms-builder/api/internal/models"
)

type routingAdapter struct {
	local      FilesystemAdapter
	firebase   FirebaseAdapter
	cloudflare CloudflarePagesAdapter
	repository RepositoryAdapter
}

func NewAdapter(deployRoot string) DeployAdapter {
	return NewAdapterWithCloudflare(deployRoot, CloudflarePagesAdapter{})
}

func NewAdapterWithCloudflare(deployRoot string, cloudflare CloudflarePagesAdapter) DeployAdapter {
	return NewAdapterWithCloudflareAndRepository(deployRoot, cloudflare, NewRepositoryAdapter())
}

func NewAdapterWithCloudflareAndRepository(deployRoot string, cloudflare CloudflarePagesAdapter, repository RepositoryAdapter) DeployAdapter {
	if strings.TrimSpace(deployRoot) == "" {
		deployRoot = filepath.Join("dist", "deployments")
	}

	return routingAdapter{
		local:      NewFilesystemAdapter(deployRoot),
		firebase:   NewFirebaseAdapter(),
		cloudflare: cloudflare,
		repository: repository,
	}
}

func (a routingAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	provider, _ := providerConfigForBuild(site, build)
	if strings.EqualFold(provider, "firebase") {
		return a.firebase.Deploy(ctx, site, build, outputPath)
	}
	if strings.EqualFold(provider, "cloudflare_pages") {
		return a.cloudflare.Deploy(ctx, site, build, outputPath)
	}
	if strings.EqualFold(provider, "git_repository") {
		return a.repository.Deploy(ctx, site, build, outputPath)
	}
	return a.local.Deploy(ctx, site, build, outputPath)
}
