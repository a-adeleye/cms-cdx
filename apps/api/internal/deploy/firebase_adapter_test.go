package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cms-builder/api/internal/models"
)

func TestFirebaseAdapterDeploysPublishedBuildToLiveChannel(t *testing.T) {
	outputPath := t.TempDir()
	mustWriteFile(t, filepath.Join(outputPath, "index.html"), "<html>live</html>")
	mustWriteFile(t, filepath.Join(outputPath, "articles", "index.html"), "<html>articles</html>")

	uploaded := make(map[string]bool)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/versions":
			writeJSON(t, w, map[string]any{"name": "sites/site-example/versions/version-1", "status": "CREATED"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/versions/version-1:populateFiles":
			var request firebasePopulateFilesRequest
			decodeBody(t, r, &request)
			if len(request.Files) != 2 {
				t.Fatalf("expected 2 files in populate request, got %d", len(request.Files))
			}
			writeJSON(t, w, firebasePopulateFilesResponse{
				UploadRequiredHashes: values(request.Files),
				UploadURL:            serverURL(r) + "/upload/sites/site-example/versions/version-1/files",
			})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/upload/sites/site-example/versions/version-1/files/"):
			uploaded[strings.TrimPrefix(r.URL.Path, "/upload/sites/site-example/versions/version-1/files/")] = true
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1beta1/sites/site-example/versions/version-1":
			if got := r.URL.Query().Get("update_mask"); got != "status" {
				t.Fatalf("expected update_mask=status, got %q", got)
			}
			writeJSON(t, w, map[string]any{"name": "sites/site-example/versions/version-1", "status": "FINALIZED"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/releases":
			if got := r.URL.Query().Get("versionName"); got != "sites/site-example/versions/version-1" {
				t.Fatalf("unexpected versionName: %q", got)
			}
			writeJSON(t, w, map[string]any{"name": "sites/site-example/releases/release-1"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	adapter := FirebaseAdapter{
		APIBaseURL:    server.URL + "/v1beta1",
		UploadBaseURL: server.URL + "/upload",
		HTTPClient:    server.Client(),
		SecretResolver: func(ref string) (string, error) {
			return "ignored", nil
		},
		TokenFetcher: func(ctx context.Context, serviceAccountJSON string) (string, error) {
			return "test-token", nil
		},
		PreviewChannelTTL: time.Hour,
	}

	site := models.Site{
		Slug:           "site-example",
		DeployProvider: "firebase",
		DeployConfig: map[string]any{
			"siteId":                  "site-example",
			"serviceAccountSecretRef": "FIREBASE_SERVICE_ACCOUNT_JSON",
		},
	}

	build := models.Build{BuildType: "published"}
	result, err := adapter.Deploy(context.Background(), site, build, outputPath)
	if err != nil {
		t.Fatalf("Deploy returned error: %v", err)
	}
	if result.Provider != "firebase" {
		t.Fatalf("expected firebase provider, got %q", result.Provider)
	}
	if result.URL != "https://site-example.web.app/" {
		t.Fatalf("expected live url, got %q", result.URL)
	}
	if len(uploaded) != 2 {
		t.Fatalf("expected 2 uploaded files, got %d", len(uploaded))
	}
}

func TestFirebaseAdapterDeploysPreviewBuildToPreviewChannel(t *testing.T) {
	outputPath := t.TempDir()
	mustWriteFile(t, filepath.Join(outputPath, "index.html"), "<html>preview</html>")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/versions":
			writeJSON(t, w, map[string]any{"name": "sites/site-example/versions/version-2", "status": "CREATED"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/versions/version-2:populateFiles":
			var request firebasePopulateFilesRequest
			decodeBody(t, r, &request)
			writeJSON(t, w, firebasePopulateFilesResponse{
				UploadRequiredHashes: values(request.Files),
				UploadURL:            serverURL(r) + "/upload/sites/site-example/versions/version-2/files",
			})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/upload/sites/site-example/versions/version-2/files/"):
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1beta1/sites/site-example/versions/version-2":
			writeJSON(t, w, map[string]any{"name": "sites/site-example/versions/version-2", "status": "FINALIZED"})
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/channels":
			writeJSON(t, w, firebaseChannelResponse{
				Name: "sites/site-example/channels/preview-site-example",
				URL:  "https://site-example--preview-site-example-random.web.app/",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v1beta1/sites/site-example/channels/preview-site-example/releases":
			writeJSON(t, w, map[string]any{"name": "sites/site-example/channels/preview-site-example/releases/release-1"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	adapter := FirebaseAdapter{
		APIBaseURL:    server.URL + "/v1beta1",
		UploadBaseURL: server.URL + "/upload",
		HTTPClient:    server.Client(),
		SecretResolver: func(ref string) (string, error) {
			return "ignored", nil
		},
		TokenFetcher: func(ctx context.Context, serviceAccountJSON string) (string, error) {
			return "test-token", nil
		},
		PreviewChannelTTL: time.Hour,
	}

	site := models.Site{
		Slug:                  "site-example",
		DeployProvider:        "firebase",
		PreviewDeployProvider: "firebase",
		PreviewDeployConfig: map[string]any{
			"siteId":                  "site-example",
			"serviceAccountSecretRef": "FIREBASE_SERVICE_ACCOUNT_JSON",
		},
	}

	build := models.Build{BuildType: "preview"}
	result, err := adapter.Deploy(context.Background(), site, build, outputPath)
	if err != nil {
		t.Fatalf("Deploy returned error: %v", err)
	}
	if result.Provider != "firebase" {
		t.Fatalf("expected firebase provider, got %q", result.Provider)
	}
	if result.URL != "https://site-example--preview-site-example-random.web.app/" {
		t.Fatalf("expected preview channel url, got %q", result.URL)
	}
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

func decodeBody[T any](t *testing.T, r *http.Request, target *T) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		t.Fatalf("decode body failed: %v", err)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("encode json failed: %v", err)
	}
}

func serverURL(r *http.Request) string {
	return "http://" + r.Host
}

func values(items map[string]string) []string {
	result := make([]string, 0, len(items))
	for _, value := range items {
		result = append(result, value)
	}
	return result
}
