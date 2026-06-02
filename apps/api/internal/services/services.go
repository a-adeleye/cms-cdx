package services

import (
	"database/sql"

	"cms-builder/api/internal/ai"
	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/config"
	"cms-builder/api/internal/deploy"
	"cms-builder/api/internal/models"
	"cms-builder/api/internal/repositories"
	"cms-builder/api/internal/storage"
)

type Services struct {
	DB       *sql.DB
	Config   config.Config
	AI       ai.AIProvider
	Builder  builder.Builder
	Deploy   deploy.DeployAdapter
	Storage  storage.StorageProvider
	Sites    repositories.SiteRepository
	Articles repositories.ArticleRepository
}

func New(db *sql.DB, cfg config.Config) Services {
	return Services{
		DB:      db,
		Config:  cfg,
		AI:      ai.NoopProvider{},
		Builder: builder.NewLocalBuilder(""),
		Deploy:  deploy.NewFilesystemAdapter(""),
		Storage: storage.NewFromConfig(cfg),
		Sites:   repositories.NoopSites{},
	}
}

func (s Services) ExampleSite() models.Site {
	return models.Site{
		Name:                  "Example Site",
		Slug:                  "example",
		Domain:                "http://localhost:4321",
		BlogPath:              "/articles",
		Status:                "active",
		TemplateKey:           "default-blog",
		DeployProvider:        "none",
		PreviewDeployProvider: "none",
	}
}
