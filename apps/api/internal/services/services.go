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
		DB:     db,
		Config: cfg,
		AI:     ai.NoopProvider{},
		Builder: builder.AstroBuilder{
			Directory:  cfg.BuilderDirectory,
			OutputRoot: cfg.BuildOutputRoot,
			NPMCommand: cfg.NPMCommand,
		},
		Deploy: deploy.CloudflarePagesAdapter{
			APIToken:  cfg.CloudflareAPIToken,
			AccountID: cfg.CloudflareAccountID,
			Command:   cfg.WranglerCommand,
		},
		Storage: storage.NoopStorage{},
		Sites:   repositories.NoopSites{},
	}
}

func (s Services) ExampleSite() models.Site {
	return models.Site{
		Name:        "Example Site",
		Slug:        "example",
		Domain:      "http://localhost:4321",
		BlogPath:    "/articles",
		Status:      "active",
		TemplateKey: "default-blog",
	}
}
