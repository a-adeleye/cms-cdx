package deploy

import (
	"context"
	"path/filepath"
	"strings"

	"cms-builder/api/internal/models"
)

type routingAdapter struct {
	local    FilesystemAdapter
	firebase FirebaseAdapter
}

func NewAdapter(deployRoot string) DeployAdapter {
	if strings.TrimSpace(deployRoot) == "" {
		deployRoot = filepath.Join("dist", "deployments")
	}

	return routingAdapter{
		local:    NewFilesystemAdapter(deployRoot),
		firebase: NewFirebaseAdapter(),
	}
}

func (a routingAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	provider, _ := providerConfigForBuild(site, build)
	if strings.EqualFold(provider, "firebase") {
		return a.firebase.Deploy(ctx, site, build, outputPath)
	}
	return a.local.Deploy(ctx, site, build, outputPath)
}
