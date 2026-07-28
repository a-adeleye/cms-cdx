package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/config"
	"cms-builder/api/internal/models"
	"cms-builder/api/internal/services"
	"cms-builder/api/internal/storage"
)

func TestIsDevStorageImageURL(t *testing.T) {
	devPublicURL := "http://localhost:9002/cms-builder"

	cases := []struct {
		name string
		url  string
		want bool
	}{
		{"matches dev public url prefix", devPublicURL + "/site-1/media/cover.png", true},
		{"matches bare localhost host", "http://localhost:4000/media/cover.png", true},
		{"matches loopback ip", "http://127.0.0.1:4000/media/cover.png", true},
		{"does not match production host", "https://cdn.example.com/site-1/media/cover.png", false},
		{"empty url", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isDevStorageImageURL(tc.url, devPublicURL); got != tc.want {
				t.Fatalf("isDevStorageImageURL(%q) = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}

func TestDevStorageObjectURLUsesEndpointReachableByAPI(t *testing.T) {
	imageURL, err := devStorageObjectURL(
		"http://localhost:9002/cms-builder/site-1/media/cover.png",
		"http://localhost:9002/cms-builder",
		"http://minio:9000",
		"cms-builder",
	)
	if err != nil {
		t.Fatalf("devStorageObjectURL returned error: %v", err)
	}
	if imageURL != "http://minio:9000/cms-builder/site-1/media/cover.png" {
		t.Fatalf("unexpected API-reachable image URL: %q", imageURL)
	}
}

func TestDevStorageObjectURLRejectsAnUnconfiguredLocalhostURL(t *testing.T) {
	_, err := devStorageObjectURL(
		"http://localhost:4000/media/cover.png",
		"http://localhost:9002/cms-builder",
		"http://minio:9000",
		"cms-builder",
	)
	if err == nil {
		t.Fatal("expected an error for an image outside configured development storage")
	}
}

func TestProductionCoverObjectKeyUsesTheBlogFolderWithoutMediaSubfolder(t *testing.T) {
	key, err := productionCoverObjectKey("/blog", "http://localhost:9002/cms-builder/site-1/media/cover.png")
	if err != nil {
		t.Fatalf("productionCoverObjectKey returned error: %v", err)
	}
	if key != "blog/cover.png" {
		t.Fatalf("unexpected production object key: %q", key)
	}
}

type capturingStorageProvider struct {
	uploaded []storage.UploadFile
}

func (f *capturingStorageProvider) Upload(ctx context.Context, file storage.UploadFile) (*storage.StoredFile, error) {
	f.uploaded = append(f.uploaded, file)
	return &storage.StoredFile{Key: "site-1/media/" + file.FileName, PublicURL: "https://cdn.example.com/site-1/media/" + file.FileName}, nil
}

func (f *capturingStorageProvider) Delete(ctx context.Context, key string) error { return nil }

func (f *capturingStorageProvider) GetPublicURL(key string) string {
	return "https://cdn.example.com/" + key
}

func TestMigrateImageToProductionStorageUploadsFetchedBytes(t *testing.T) {
	pngBytes := testPNG(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(pngBytes)
	}))
	defer server.Close()

	fake := &capturingStorageProvider{}
	newURL, err := migrateImageToProductionStorage(context.Background(), server.Client(), fake, "site-1", server.URL+"/media/cover.png", "blog/cover.png")
	if err != nil {
		t.Fatalf("migrateImageToProductionStorage returned error: %v", err)
	}
	if newURL != "https://cdn.example.com/site-1/media/cover.webp" {
		t.Fatalf("unexpected migrated URL: %q", newURL)
	}
	if len(fake.uploaded) != 1 {
		t.Fatalf("expected exactly one upload, got %d", len(fake.uploaded))
	}
	uploaded := fake.uploaded[0]
	if uploaded.FileName != "cover.webp" {
		t.Fatalf("expected filename cover.webp, got %q", uploaded.FileName)
	}
	if uploaded.SiteID != "site-1" {
		t.Fatalf("expected siteID site-1, got %q", uploaded.SiteID)
	}
	if uploaded.ObjectKey != "blog/cover.webp" {
		t.Fatalf("expected object key blog/cover.webp, got %q", uploaded.ObjectKey)
	}
	if detected := http.DetectContentType(uploaded.Contents); detected != "image/webp" {
		t.Fatalf("expected an uploaded WebP image, got %q", detected)
	}
}

func TestMigrateImageToProductionStorageReturnsErrorOnFetchFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	fake := &capturingStorageProvider{}
	if _, err := migrateImageToProductionStorage(context.Background(), server.Client(), fake, "site-1", server.URL+"/missing.png", "blog/missing.png"); err == nil {
		t.Fatal("expected an error when the source image cannot be fetched")
	}
	if len(fake.uploaded) != 0 {
		t.Fatalf("expected no uploads on fetch failure, got %d", len(fake.uploaded))
	}
}

func TestRewritePublishedCoverImagesPersistsR2URLAndSkipsLaterUpload(t *testing.T) {
	db := openTestDatabase(t)
	t.Cleanup(func() { _ = db.Close() })
	ctx := context.Background()
	siteID := mustQueryText(t, db, ctx, `
		INSERT INTO sites (name, slug, blog_path)
		VALUES ('R2 persistence test', 'r2-persistence-' || replace(gen_random_uuid()::text, '-', ''), '/blog')
		RETURNING id::text
	`)
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM sites WHERE id = $1`, siteID)
	})
	articleID := mustQueryText(t, db, ctx, `
		INSERT INTO articles (site_id, title, slug, content_markdown, cover_image_url, status)
		VALUES ($1, 'Existing article', 'existing-article', 'Body', 'http://development.invalid/cover.png', 'published')
		RETURNING id::text
	`, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(testPNG(t))
	}))
	defer server.Close()

	devPublicURL := server.URL + "/cms-builder"
	api := &API{
		Services: services.Services{DB: db},
		Config: config.Config{
			S3Endpoint:  server.URL,
			S3Bucket:    "cms-builder",
			S3PublicURL: devPublicURL,
		},
	}
	content := builder.SiteContent{Articles: []builder.ArticleContent{{
		ID: articleID, Title: "Existing article", CoverImageURL: devPublicURL + "/" + siteID + "/media/cover.png",
	}}}
	provider := &capturingStorageProvider{}
	site := models.Site{ID: siteID, BlogPath: "/blog"}

	if err := api.rewritePublishedCoverImagesWithStorage(ctx, &content, site, provider); err != nil {
		t.Fatalf("first migration returned error: %v", err)
	}
	if len(provider.uploaded) != 1 {
		t.Fatalf("expected one initial upload, got %d", len(provider.uploaded))
	}
	if content.Articles[0].CoverImageURL != "https://cdn.example.com/site-1/media/cover.webp" {
		t.Fatalf("unexpected rewritten cover URL: %q", content.Articles[0].CoverImageURL)
	}
	var persistedURL string
	if err := db.QueryRowContext(ctx, `SELECT cover_image_url FROM articles WHERE id = $1`, articleID).Scan(&persistedURL); err != nil {
		t.Fatal(err)
	}
	if persistedURL != content.Articles[0].CoverImageURL {
		t.Fatalf("expected persisted cover URL %q, got %q", content.Articles[0].CoverImageURL, persistedURL)
	}

	if err := api.rewritePublishedCoverImagesWithStorage(ctx, &content, site, provider); err != nil {
		t.Fatalf("later production build returned error: %v", err)
	}
	if len(provider.uploaded) != 1 {
		t.Fatalf("expected later build to reuse the persisted R2 URL, got %d uploads", len(provider.uploaded))
	}
}
