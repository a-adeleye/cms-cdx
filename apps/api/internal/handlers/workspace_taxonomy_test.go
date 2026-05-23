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

func TestTagCRUDCascadesArticleTagLinksWhenDeleted(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{Services: services.Services{DB: db}}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)
	authorID := mustQueryText(t, db, ctx, `SELECT id::text FROM authors WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)
	categoryID := mustQueryText(t, db, ctx, `SELECT id::text FROM categories WHERE site_id = $1 ORDER BY name ASC LIMIT 1`, siteID)

	tag, err := api.createTag(ctx, siteID, tagUpsertRequest{
		Name: "Launch Plan",
	})
	if err != nil {
		t.Fatalf("createTag returned error: %v", err)
	}
	if tag.ID == "" {
		t.Fatal("expected created tag ID")
	}
	if tag.Slug == "" {
		t.Fatal("expected created tag slug")
	}

	updatedTag, err := api.updateTag(ctx, siteID, tag.ID, tagUpsertRequest{
		Name: "Launch Strategy",
	})
	if err != nil {
		t.Fatalf("updateTag returned error: %v", err)
	}
	if updatedTag.Name != "Launch Strategy" {
		t.Fatalf("expected updated tag name, got %q", updatedTag.Name)
	}
	if updatedTag.Slug == "" {
		t.Fatal("expected updated tag slug")
	}

	article, err := api.upsertArticleWithSite(ctx, siteID, articleUpsertRequest{
		Title:           "Tag regression article",
		Slug:            "tag-regression-article",
		Excerpt:         "An article used to verify tag cleanup.",
		ContentMarkdown: "This article references the created tag.",
		AuthorID:        authorID,
		CategoryID:      categoryID,
		TagIDs:          []string{updatedTag.ID},
		IsFeatured:      false,
		Status:          "draft",
	})
	if err != nil {
		t.Fatalf("upsertArticleWithSite returned error: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, article.ID)
		_, _ = db.ExecContext(ctx, `DELETE FROM tags WHERE id = $1`, updatedTag.ID)
		_ = db.Close()
	})

	if err := api.deleteTag(ctx, siteID, updatedTag.ID); err != nil {
		t.Fatalf("deleteTag returned error: %v", err)
	}

	var linkCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM article_tags WHERE article_id = $1`, article.ID).Scan(&linkCount); err != nil {
		t.Fatalf("query article tag links failed: %v", err)
	}
	if linkCount != 0 {
		t.Fatalf("expected deleted tag to cascade from article_tags, got %d links", linkCount)
	}
}
