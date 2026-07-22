package storage

import (
	"errors"
	"testing"
)

func TestParseConfigAllowsEmptyConfig(t *testing.T) {
	cfg, err := ParseConfig("")
	if err != nil {
		t.Fatalf("ParseConfig returned error for empty config: %v", err)
	}
	if cfg.IsConfigured() {
		t.Fatalf("expected empty config to be unconfigured, got %#v", cfg)
	}

	cfg, err = ParseConfig("{}")
	if err != nil {
		t.Fatalf("ParseConfig returned error for {}: %v", err)
	}
	if cfg.IsConfigured() {
		t.Fatalf("expected {} config to be unconfigured, got %#v", cfg)
	}
}

func TestParseConfigRequiresBucketAndHTTPSPublicURL(t *testing.T) {
	if _, err := ParseConfig(`{"publicUrl":"https://cdn.example.com","accessKeySecretRef":"A","secretKeySecretRef":"B"}`); err == nil {
		t.Fatal("expected error when bucket is missing")
	}
	if _, err := ParseConfig(`{"bucket":"prod-media","publicUrl":"http://cdn.example.com","accessKeySecretRef":"A","secretKeySecretRef":"B"}`); err == nil {
		t.Fatal("expected error when publicUrl is not HTTPS")
	}
	if _, err := ParseConfig(`{"bucket":"prod-media","publicUrl":"https://cdn.example.com","accessKeySecretRef":"not a valid ref","secretKeySecretRef":"B"}`); err == nil {
		t.Fatal("expected error for invalid accessKeySecretRef")
	}
}

func TestParseConfigDefaultsRegionAndAcceptsValidConfig(t *testing.T) {
	cfg, err := ParseConfig(`{"bucket":"prod-media","publicUrl":"https://cdn.example.com/","accessKeySecretRef":"PROD_S3_ACCESS_KEY","secretKeySecretRef":"PROD_S3_SECRET_KEY"}`)
	if err != nil {
		t.Fatalf("ParseConfig returned error for valid config: %v", err)
	}
	if !cfg.IsConfigured() {
		t.Fatalf("expected config to be configured, got %#v", cfg)
	}
	if cfg.Region != "us-east-1" {
		t.Fatalf("expected default region us-east-1, got %q", cfg.Region)
	}
	if cfg.PublicURL != "https://cdn.example.com" {
		t.Fatalf("expected trailing slash trimmed from publicUrl, got %q", cfg.PublicURL)
	}
}

func TestParseConfigRejectsUnknownFields(t *testing.T) {
	if _, err := ParseConfig(`{"bucket":"prod-media","unexpected":"value"}`); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestNewFromSiteConfigRequiresConfiguredConfig(t *testing.T) {
	_, err := NewFromSiteConfig(Config{}, func(string) (string, bool) { return "", false })
	if !errors.Is(err, ErrNotConfigured) {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestNewFromSiteConfigRequiresSecretsToBeAvailable(t *testing.T) {
	cfg := Config{
		Bucket:             "prod-media",
		Region:             "us-east-1",
		PublicURL:          "https://cdn.example.com",
		AccessKeySecretRef: "PROD_S3_ACCESS_KEY",
		SecretKeySecretRef: "PROD_S3_SECRET_KEY",
	}
	_, err := NewFromSiteConfig(cfg, func(string) (string, bool) { return "", false })
	if !errors.Is(err, ErrSecretUnavailable) {
		t.Fatalf("expected ErrSecretUnavailable, got %v", err)
	}
}

func TestNewFromSiteConfigBuildsProviderWhenSecretsAvailable(t *testing.T) {
	cfg := Config{
		Bucket:             "prod-media",
		Region:             "us-east-1",
		PublicURL:          "https://cdn.example.com",
		AccessKeySecretRef: "PROD_S3_ACCESS_KEY",
		SecretKeySecretRef: "PROD_S3_SECRET_KEY",
	}
	secrets := map[string]string{
		"PROD_S3_ACCESS_KEY": "AKIAEXAMPLE",
		"PROD_S3_SECRET_KEY": "secret",
	}
	provider, err := NewFromSiteConfig(cfg, func(name string) (string, bool) {
		value, ok := secrets[name]
		return value, ok
	})
	if err != nil {
		t.Fatalf("NewFromSiteConfig returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected a non-nil storage provider")
	}
	if got := provider.GetPublicURL("site-1/media/cover.png"); got != "https://cdn.example.com/site-1/media/cover.png" {
		t.Fatalf("unexpected public URL: %q", got)
	}
}
