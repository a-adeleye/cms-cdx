package repositories

import (
	"context"

	"cms-builder/api/internal/models"
)

type SiteRepository interface {
	List(ctx context.Context) ([]models.Site, error)
	Get(ctx context.Context, id string) (*models.Site, error)
}

type ArticleRepository interface {
	ListBySite(ctx context.Context, siteID string, publishedOnly bool) ([]models.Article, error)
	Get(ctx context.Context, id string) (*models.Article, error)
}

type NoopSites struct{}

func (NoopSites) List(ctx context.Context) ([]models.Site, error) { return nil, nil }
func (NoopSites) Get(ctx context.Context, id string) (*models.Site, error) { return nil, nil }

