package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	// ErrNotConfigured indicates a site has no production storage configured.
	ErrNotConfigured = errors.New("production storage is not configured")
	// ErrSecretUnavailable indicates a configured secret env var is unset on the server.
	ErrSecretUnavailable = errors.New("production storage secret is unavailable")

	secretRefPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,127}$`)
)

// Config is a site's production S3 bucket configuration. API credentials are
// resolved only from named server-side environment variables, never stored
// or returned to the browser.
type Config struct {
	Bucket             string `json:"bucket"`
	Region             string `json:"region"`
	Endpoint           string `json:"endpoint"`
	PublicURL          string `json:"publicUrl"`
	AccessKeySecretRef string `json:"accessKeySecretRef"`
	SecretKeySecretRef string `json:"secretKeySecretRef"`
}

// IsConfigured reports whether enough fields are present to use this config
// for production uploads.
func (c Config) IsConfigured() bool {
	return c.Bucket != "" && c.PublicURL != "" && c.AccessKeySecretRef != "" && c.SecretKeySecretRef != ""
}

// ParseConfig validates a site's stored production-storage configuration. An
// empty configuration is valid but leaves IsConfigured() false.
func ParseConfig(raw string) (Config, error) {
	if strings.TrimSpace(raw) == "" {
		raw = "{}"
	}

	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("production storage config must be JSON with supported fields: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return Config{}, errors.New("production storage config must contain one JSON object")
	}

	cfg.Bucket = strings.TrimSpace(cfg.Bucket)
	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.Endpoint = strings.TrimRight(strings.TrimSpace(cfg.Endpoint), "/")
	cfg.PublicURL = strings.TrimRight(strings.TrimSpace(cfg.PublicURL), "/")
	cfg.AccessKeySecretRef = strings.TrimSpace(cfg.AccessKeySecretRef)
	cfg.SecretKeySecretRef = strings.TrimSpace(cfg.SecretKeySecretRef)

	if cfg.Bucket == "" && cfg.PublicURL == "" && cfg.Endpoint == "" && cfg.AccessKeySecretRef == "" && cfg.SecretKeySecretRef == "" {
		return cfg, nil
	}

	if cfg.Bucket == "" {
		return Config{}, errors.New("production storage requires a bucket")
	}
	if parsed, err := url.Parse(cfg.PublicURL); err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" {
		return Config{}, errors.New("production storage requires an HTTPS public URL")
	}
	if cfg.Endpoint != "" {
		if parsed, err := url.Parse(cfg.Endpoint); err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" {
			return Config{}, errors.New("production storage endpoint must be an HTTPS URL")
		}
	}
	if !secretRefPattern.MatchString(cfg.AccessKeySecretRef) || !secretRefPattern.MatchString(cfg.SecretKeySecretRef) {
		return Config{}, errors.New("production storage requires accessKeySecretRef and secretKeySecretRef as environment variable names")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	return cfg, nil
}

// NewFromSiteConfig builds a StorageProvider for a site's configured
// production bucket, resolving credentials from the server environment.
func NewFromSiteConfig(cfg Config, lookupEnv func(string) (string, bool)) (StorageProvider, error) {
	if !cfg.IsConfigured() {
		return nil, ErrNotConfigured
	}
	if lookupEnv == nil {
		lookupEnv = os.LookupEnv
	}

	accessKey, ok := lookupEnv(cfg.AccessKeySecretRef)
	if !ok || strings.TrimSpace(accessKey) == "" {
		return nil, ErrSecretUnavailable
	}
	secretKey, ok := lookupEnv(cfg.SecretKeySecretRef)
	if !ok || strings.TrimSpace(secretKey) == "" {
		return nil, ErrSecretUnavailable
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://s3.%s.amazonaws.com", cfg.Region)
	}

	provider := newS3Storage(endpoint, cfg.Region, cfg.Bucket, accessKey, secretKey, cfg.PublicURL)
	if _, disabled := provider.(disabledStorage); disabled {
		return nil, errors.New("unable to configure production storage")
	}
	return provider, nil
}
