package handlers

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"cms-builder/api/internal/services"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUpsertArticleCreatesDraftWithAuthorAndCategory(t *testing.T) {
	db := openTestDatabase(t)

	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Regression test article",
		Slug:            "regression-test-article",
		Excerpt:         "A short regression test article",
		ContentMarkdown: "This article verifies draft saves work.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}
	if article.ID == "" {
		t.Fatal("expected created article ID")
	}
	if article.Status != "draft" {
		t.Fatalf("expected status draft, got %q", article.Status)
	}
	if article.AuthorID != authorID {
		t.Fatalf("expected author %q, got %q", authorID, article.AuthorID)
	}
	if article.CategoryID != categoryID {
		t.Fatalf("expected category %q, got %q", categoryID, article.CategoryID)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_ = db.Close()
	})
}

func TestUpsertArticleNormalizesFreeTextTags(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Regression test article with tags",
		Slug:            "regression-test-article-with-tags",
		Excerpt:         "A short regression test article",
		ContentMarkdown: "This article verifies free-text tags are stored as-is.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		Tags:            "Privacy,  privacy, VPN , security",
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_ = db.Close()
	})

	if got, want := article.Tags, "Privacy, VPN, security"; got != want {
		t.Fatalf("expected normalized tags %q, got %q", want, got)
	}
}

func TestDeleteArticleRemovesArticle(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Delete test article",
		Slug:            "delete-test-article",
		Excerpt:         "A short article to delete",
		ContentMarkdown: "This article exists for deletion testing.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		Tags:            "privacy, security",
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}

	if err := api.deleteArticle(ctx, article.ID); err != nil {
		t.Fatalf("deleteArticle returned error: %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM articles WHERE id = $1`, article.ID).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected article to be deleted, found %d rows", count)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})
}

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://cms:cms@localhost:5433/cms_builder?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Skipf("skipping database-backed test: %v", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Skipf("skipping database-backed test: %v", err)
	}

	return db
}

func mustQueryText(t *testing.T, db *sql.DB, ctx context.Context, query string, args ...any) string {
	t.Helper()

	var value string
	if err := db.QueryRowContext(ctx, query, args...).Scan(&value); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	return value
}
