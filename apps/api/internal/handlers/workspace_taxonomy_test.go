package handlers

import (
	"context"
	"database/sql"
	"testing"

	"cms-builder/api/internal/services"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestCategoryCRUDClearsArticleCategoryWhenDeleted(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	category, err := api.createCategory(ctx, siteID, categoryUpsertRequest{
		Name:        "Editorial Flow",
		Description: "Used by regression tests.",
	})
	if err != nil {
		t.Fatalf("createCategory returned error: %v", err)
	}
	if category.ID == "" {
		t.Fatal("expected created category ID")
	}
	if category.Slug == "" {
		t.Fatal("expected created category slug")
	}

	updatedCategory, err := api.updateCategory(ctx, siteID, category.ID, categoryUpsertRequest{
		Name:        "Editorial Workflow",
		Description: "Updated taxonomy description.",
	})
	if err != nil {
		t.Fatalf("updateCategory returned error: %v", err)
	}
	if updatedCategory.Name != "Editorial Workflow" {
		t.Fatalf("expected updated category name, got %q", updatedCategory.Name)
	}
	if updatedCategory.Slug == "" {
		t.Fatal("expected updated category slug")
	}

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Category regression article",
		Slug:            "category-regression-article",
		Excerpt:         "An article used to verify category cleanup.",
		ContentMarkdown: "This article references the created category.",
		AuthorID:        authorID,
		CategoryID:      updatedCategory.ID,
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_, _ = db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, updatedCategory.ID)
		_ = db.Close()
	})

	if err := api.deleteCategory(ctx, siteID, updatedCategory.ID); err != nil {
		t.Fatalf("deleteCategory returned error: %v", err)
	}

	var categoryID sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT category_id::text FROM articles WHERE id = $1`, article.ID).Scan(&categoryID); err != nil {
		t.Fatalf("query article category failed: %v", err)
	}
	if categoryID.Valid {
		t.Fatalf("expected deleted category to be cleared from article, got %q", categoryID.String)
	}
}

func TestArticleFreeTextTagsRoundTripThroughUpdate(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Tag regression article",
		Slug:            "tag-regression-article",
		Excerpt:         "An article used to verify free-text tags.",
		ContentMarkdown: "This article carries AI-suggested free-text tags.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		Tags:            "launch plan",
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

	if article.Tags != "launch plan" {
		t.Fatalf("expected tags %q, got %q", "launch plan", article.Tags)
	}

	updated, err := api.upsertArticle(ctx, articleUpsertRequest{
		ID:              article.ID,
		Title:           article.Title,
		Slug:            article.Slug,
		Excerpt:         article.Excerpt,
		ContentMarkdown: article.ContentMarkdown,
		AuthorID:        article.AuthorID,
		CategoryID:      article.CategoryID,
		Tags:            "launch strategy, growth",
		IsFeatured:      article.IsFeatured,
		Status:          article.Status,
	})
	if err != nil {
		t.Fatalf("upsertArticle returned error: %v", err)
	}
	if updated.Tags != "launch strategy, growth" {
		t.Fatalf("expected updated tags %q, got %q", "launch strategy, growth", updated.Tags)
	}
}

func TestAuthorCRUDClearsArticleAuthorWhenDeleted(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	author, err := api.createAuthor(ctx, siteID, authorUpsertRequest{
		Name: "Test Author",
		Bio:  "Used by regression tests.",
	})
	if err != nil {
		t.Fatalf("createAuthor returned error: %v", err)
	}
	if author.ID == "" {
		t.Fatal("expected created author ID")
	}
	if author.Slug == "" {
		t.Fatal("expected created author slug")
	}

	updatedAuthor, err := api.updateAuthor(ctx, siteID, author.ID, authorUpsertRequest{
		Name: "Updated Author",
		Bio:  "Updated author bio.",
	})
	if err != nil {
		t.Fatalf("updateAuthor returned error: %v", err)
	}
	if updatedAuthor.Name != "Updated Author" {
		t.Fatalf("expected updated author name, got %q", updatedAuthor.Name)
	}
	if updatedAuthor.Slug == "" {
		t.Fatal("expected updated author slug")
	}

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Author regression article",
		Slug:            "author-regression-article",
		Excerpt:         "An article used to verify author cleanup.",
		ContentMarkdown: "This article references the created author.",
		AuthorID:        updatedAuthor.ID,
		CategoryID:      categoryID,
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_, _ = db.ExecContext(ctx, `DELETE FROM authors WHERE id = $1`, updatedAuthor.ID)
		_ = db.Close()
	})

	if err := api.deleteAuthor(ctx, siteID, updatedAuthor.ID); err != nil {
		t.Fatalf("deleteAuthor returned error: %v", err)
	}

	var authorID sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT author_id::text FROM articles WHERE id = $1`, article.ID).Scan(&authorID); err != nil {
		t.Fatalf("query article author failed: %v", err)
	}
	if authorID.Valid {
		t.Fatalf("expected deleted author to be cleared from article, got %q", authorID.String)
	}
}
