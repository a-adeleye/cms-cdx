package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"cms-builder/api/internal/builder"
	"cms-builder/api/internal/services"
)

type previewBuilderStub struct {
	outputPath string
	calls      int
}

func (b *previewBuilderStub) GenerateSite(context.Context, builder.SiteContent, builder.GenerateOptions) (string, error) {
	return "", nil
}

func (b *previewBuilderStub) GenerateTemplatePreview(_ context.Context, templateKey string) (string, error) {
	b.calls++
	if templateKey != builder.SupromailTemplateKey {
		return "", os.ErrNotExist
	}
	return b.outputPath, nil
}

func TestTemplatePreviewRendersAndCachesSupromailAstroOutput(t *testing.T) {
	outputPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(outputPath, "index.html"), []byte("<h1>Supromail preview</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}
	previewBuilder := &previewBuilderStub{outputPath: outputPath}
	api := &API{Services: services.Services{Builder: previewBuilder}}

	for attempt := 0; attempt < 2; attempt++ {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/template-previews/supromail", nil)
		response := httptest.NewRecorder()
		api.templatePreview(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
		if response.Body.String() != "<h1>Supromail preview</h1>" {
			t.Fatalf("unexpected preview body: %q", response.Body.String())
		}
		if response.Header().Get("Content-Security-Policy") == "" {
			t.Fatal("expected preview content security policy")
		}
	}
	if previewBuilder.calls != 1 {
		t.Fatalf("expected a single cached preview build, got %d", previewBuilder.calls)
	}
}

func TestTemplatePreviewRejectsNonGetRequests(t *testing.T) {
	api := &API{}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/template-previews/supromail", nil)
	response := httptest.NewRecorder()
	api.templatePreview(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}
}
