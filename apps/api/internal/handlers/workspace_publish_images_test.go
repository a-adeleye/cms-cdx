package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
	pngBytes := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(pngBytes)
	}))
	defer server.Close()

	fake := &capturingStorageProvider{}
	newURL, err := migrateImageToProductionStorage(context.Background(), server.Client(), fake, "site-1", server.URL+"/media/cover.png")
	if err != nil {
		t.Fatalf("migrateImageToProductionStorage returned error: %v", err)
	}
	if newURL != "https://cdn.example.com/site-1/media/cover.png" {
		t.Fatalf("unexpected migrated URL: %q", newURL)
	}
	if len(fake.uploaded) != 1 {
		t.Fatalf("expected exactly one upload, got %d", len(fake.uploaded))
	}
	uploaded := fake.uploaded[0]
	if uploaded.FileName != "cover.png" {
		t.Fatalf("expected filename cover.png, got %q", uploaded.FileName)
	}
	if uploaded.SiteID != "site-1" {
		t.Fatalf("expected siteID site-1, got %q", uploaded.SiteID)
	}
	if string(uploaded.Contents) != string(pngBytes) {
		t.Fatalf("expected uploaded contents to match fetched bytes")
	}
}

func TestMigrateImageToProductionStorageReturnsErrorOnFetchFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	fake := &capturingStorageProvider{}
	if _, err := migrateImageToProductionStorage(context.Background(), server.Client(), fake, "site-1", server.URL+"/missing.png"); err == nil {
		t.Fatal("expected an error when the source image cannot be fetched")
	}
	if len(fake.uploaded) != 0 {
		t.Fatalf("expected no uploads on fetch failure, got %d", len(fake.uploaded))
	}
}
