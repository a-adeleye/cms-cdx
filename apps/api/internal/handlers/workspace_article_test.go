package handlers

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
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
		_, _ = db.ExecContext(ctx, `DELETE FROM article_tags WHERE article_id = $1`, article.ID)
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_ = db.Close()
	})
}

func TestUpsertArticleReturnsValidationErrorForInvalidTagIds(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	_, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Regression test article with tags",
		Slug:            "regression-test-article-with-tags",
		Excerpt:         "A short regression test article",
		ContentMarkdown: "This article verifies invalid tag ids are rejected.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		TagIDs:          []string{"not-a-uuid", "still-not-a-uuid"},
		IsFeatured:      false,
		Status:          "draft",
	})
	if err == nil {
		t.Fatal("expected validation error for invalid tag ids")
	}
	if !errors.Is(err, errValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if got := err.Error(); !strings.Contains(got, "invalid tag id") {
		t.Fatalf("expected invalid tag id message, got %q", got)
	}

	_ = db.Close()
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
