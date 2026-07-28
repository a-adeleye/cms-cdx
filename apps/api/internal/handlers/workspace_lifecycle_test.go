package handlers

import (
	"context"
	"errors"
	"strings"
	"testing"

	"cms-builder/api/internal/services"
)

func TestDeleteSiteCascadesItsRecordsButPreservesAnotherSite(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `
		INSERT INTO sites (name, slug, blog_path)
		VALUES ('Delete test', 'delete-test-' || replace(gen_random_uuid()::text, '-', ''), '/articles')
		RETURNING id::text
	`)
	if _, err := db.ExecContext(ctx, `INSERT INTO builds (site_id, status, build_type) VALUES ($1, 'success', 'preview')`, siteID); err != nil {
		t.Fatal(err)
	}

	api := &API{Services: services.Services{DB: db}}
	if err := api.deleteSite(ctx, siteID); err != nil {
		t.Fatalf("deleteSite returned error: %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sites WHERE id = $1`, siteID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected deleted site to be absent, found %d records", count)
	}
}

func TestClearBuildHistoryDeletesOnlyTheSelectedSiteBuilds(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `
		INSERT INTO sites (name, slug, blog_path)
		VALUES ('History test', 'history-test-' || replace(gen_random_uuid()::text, '-', ''), '/articles')
		RETURNING id::text
	`)
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, `DELETE FROM sites WHERE id = $1`, siteID) })
	if _, err := db.ExecContext(ctx, `INSERT INTO builds (site_id, status, build_type) VALUES ($1, 'success', 'preview'), ($1, 'success', 'published')`, siteID); err != nil {
		t.Fatal(err)
	}

	api := &API{Services: services.Services{DB: db}}
	if err := api.clearBuildHistory(ctx, siteID); err != nil {
		t.Fatalf("clearBuildHistory returned error: %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM builds WHERE site_id = $1`, siteID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no builds after clearing history, found %d", count)
	}
}

func TestMediaAssetManagementUpdatesReplacesAndProtectsReferencedAssets(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY created_at LIMIT 1`)
	storage := &fakeStorageProvider{}
	api := &API{Services: services.Services{DB: db, Storage: storage}}

	asset, err := api.createMediaAsset(ctx, siteID, mediaUpsertRequest{
		FileName: "old.png", FileURL: "https://cdn.example/old.png", MimeType: "image/png", SizeBytes: 12,
		StorageProvider: "test", StorageKey: "site-example/media/old.png", AltText: "Old alt text",
	})
	if err != nil {
		t.Fatalf("createMediaAsset returned error: %v", err)
	}
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, `DELETE FROM media_assets WHERE id = $1`, asset.ID) })

	updated, err := api.updateMediaAsset(ctx, siteID, asset.ID, mediaUpdateRequest{AltText: "Updated alt text"})
	if err != nil {
		t.Fatalf("updateMediaAsset returned error: %v", err)
	}
	if updated.AltText != "Updated alt text" {
		t.Fatalf("expected updated alt text, got %q", updated.AltText)
	}
	if _, err := api.updateMediaAsset(ctx, siteID, asset.ID, mediaUpdateRequest{AltText: strings.Repeat("a", maxMediaAltTextRunes+1)}); !errors.Is(err, errValidation) {
		t.Fatalf("expected oversized alt text to be rejected, got %v", err)
	}

	replaced, err := api.replaceMediaAsset(ctx, siteID, asset.ID, "replacement.png", testPNG(t), "image/png", "Replacement alt text")
	if err != nil {
		t.Fatalf("replaceMediaAsset returned error: %v", err)
	}
	if replaced.AltText != "Replacement alt text" || len(storage.uploaded) != 1 || storage.uploaded[0].ObjectKey != asset.StorageKey {
		t.Fatalf("expected replacement to preserve the storage key and alt text, got %#v and %#v", replaced, storage.uploaded)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO articles (site_id, title, slug, content_markdown, cover_image_url) VALUES ($1, 'Uses media', 'uses-media-' || replace(gen_random_uuid()::text, '-', ''), 'Content', $2)`, siteID, replaced.FileURL); err != nil {
		t.Fatal(err)
	}
	if err := api.deleteMediaAsset(ctx, siteID, asset.ID); !errors.Is(err, errConflict) {
		t.Fatalf("expected referenced asset deletion to be rejected, got %v", err)
	}

	if _, err := db.ExecContext(ctx, `DELETE FROM articles WHERE site_id = $1 AND cover_image_url = $2`, siteID, replaced.FileURL); err != nil {
		t.Fatal(err)
	}
	if err := api.deleteMediaAsset(ctx, siteID, asset.ID); err != nil {
		t.Fatalf("deleteMediaAsset returned error: %v", err)
	}
	if len(storage.deleted) != 1 || storage.deleted[0] != replaced.StorageKey {
		t.Fatalf("expected storage object to be deleted, got %#v", storage.deleted)
	}
}
