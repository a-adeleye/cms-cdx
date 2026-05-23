package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cms-builder/api/internal/auth"
)

func TestAuthRejectsMissingToken(t *testing.T) {
	t.Parallel()

	handler := Auth([]byte("secret"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthAllowsValidToken(t *testing.T) {
	t.Parallel()

	token, err := auth.CreateToken([]byte("secret"), "user-1", "admin", time.Hour)
	if err != nil {
		t.Fatalf("CreateToken() error = %v", err)
	}

	var nextCalled bool
	handler := Auth([]byte("secret"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		if got := UserIDFromContext(r.Context()); got != "user-1" {
			t.Fatalf("UserIDFromContext() = %q, want %q", got, "user-1")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("next handler was not called")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
