package builder

import (
	"context"

	"cms-builder/api/internal/models"
)

type GenerateOptions struct {
	SiteID  string
	Preview bool
}

type Builder interface {
	GenerateSite(ctx context.Context, site models.Site, options GenerateOptions) (string, error)
}

type NoopBuilder struct{}

func (NoopBuilder) GenerateSite(ctx context.Context, site models.Site, options GenerateOptions) (string, error) {
	return "dist/sites/" + site.Slug, nil
}

