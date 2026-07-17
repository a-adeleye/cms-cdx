package builder

import (
	"context"

	"cms-builder/api/internal/models"
)

type GenerateOptions struct {
	SiteID     string
	Preview    bool
	ArticleIDs []string
}

type TagContent struct {
	Name string
	Slug string
}

type ArticleContent struct {
	ID              string
	Title           string
	Slug            string
	Excerpt         string
	ContentMarkdown string
	CoverImageURL   string
	Status          string
	IsFeatured      bool
	PublishedAt     string
	UpdatedAt       string
	SEOTitle        string
	SEODescription  string
	CanonicalURL    string
	AuthorName      string
	AuthorSlug      string
	CategoryName    string
	CategorySlug    string
	Tags            []TagContent
}

type LandingSectionContent struct {
	SectionKey   string
	Title        string
	Subtitle     string
	ContentJSON  map[string]any
	DisplayOrder int
	IsEnabled    bool
}

type SiteContent struct {
	Site            models.Site
	LandingSections []LandingSectionContent
	Articles        []ArticleContent
}

type Builder interface {
	GenerateSite(ctx context.Context, content SiteContent, options GenerateOptions) (string, error)
}

type NoopBuilder struct{}

func (NoopBuilder) GenerateSite(ctx context.Context, content SiteContent, options GenerateOptions) (string, error) {
	if options.Preview {
		return "dist/preview/" + content.Site.Slug, nil
	}
	return "dist/sites/" + content.Site.Slug, nil
}
