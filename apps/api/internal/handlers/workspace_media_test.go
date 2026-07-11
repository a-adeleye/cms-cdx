package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"cms-builder/api/internal/config"
	"cms-builder/api/internal/services"
	"cms-builder/api/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type fakeStorageProvider struct{}

func (fakeStorageProvider) Upload(ctx context.Context, file storage.UploadFile) (*storage.StoredFile, error) {
	return &storage.StoredFile{
		Key:       "site-example/media/cover.jpg",
		PublicURL: "https://cdn.example/cover.jpg",
	}, nil
}

func (fakeStorageProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (fakeStorageProvider) GetPublicURL(key string) string {
	return "https://cdn.example/" + key
}

func TestMediaUploadStoresUploadedImageMetadata(t *testing.T) {
	db := openTestDatabase(t)
	ctx := context.Background()
	api := &API{
		Services: services.Services{
			DB:      db,
			Storage: fakeStorageProvider{},
		},
		Config: config.Config{
			S3Endpoint:  "http://minio:9000",
			S3PublicURL: "http://localhost:9002/cms-builder",
		},
	}

	siteID := mustQueryText(t, db, ctx, `SELECT id::text FROM sites ORDER BY updated_at DESC LIMIT 1`)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filePart, err := writer.CreateFormFile("file", "cover.jpg")
	if err != nil {
		t.Fatalf("CreateFormFile returned error: %v", err)
	}
	if _, err := filePart.Write([]byte{0x89, 0x50, 0x4e, 0x47}); err != nil {
		t.Fatalf("write file content failed: %v", err)
	}
	if err := writer.WriteField("altText", "Cover image"); err != nil {
		t.Fatalf("WriteField returned error: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("writer close failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sites/"+siteID+"/media", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	api.handleMediaRoutes(rec, req, siteID)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var (
		fileURL         string
		storageProvider string
		storageKey      sql.NullString
	)
	if err := db.QueryRowContext(ctx, `
		SELECT file_url, storage_provider, storage_key
		FROM media_assets
		WHERE site_id = $1 AND file_name = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, siteID, "cover.jpg").Scan(&fileURL, &storageProvider, &storageKey); err != nil {
		t.Fatalf("query media asset failed: %v", err)
	}

	if fileURL != "https://cdn.example/cover.jpg" {
		t.Fatalf("expected uploaded file URL, got %q", fileURL)
	}
	if storageProvider != "minio" {
		t.Fatalf("expected storage provider minio, got %q", storageProvider)
	}
	if !storageKey.Valid || storageKey.String == "" {
		t.Fatal("expected storage key to be recorded")
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM media_assets WHERE site_id = $1 AND file_name = $2`, siteID, "cover.jpg")
		_ = db.Close()
	})
}
