package config

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	APIPort     string

	S3Endpoint          string
	S3Region            string
	S3Bucket            string
	S3AccessKey         string
	S3SecretKey         string
	S3PublicURL         string
	OpenAIAPIKey        string
	GeminiAPIKey        string
	AnthropicAPIKey     string
	CloudflareAPIToken  string
	CloudflareAccountID string
	BuilderDirectory    string
	BuildOutputRoot     string
	NPMCommand          string
	WranglerCommand     string
}

func Load() Config {
	return Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", "dev-only-secret"),
		APIPort:             getEnv("API_PORT", "8080"),
		S3Endpoint:          os.Getenv("S3_ENDPOINT"),
		S3Region:            getEnv("S3_REGION", "us-east-1"),
		S3Bucket:            os.Getenv("S3_BUCKET"),
		S3AccessKey:         os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:         os.Getenv("S3_SECRET_KEY"),
		S3PublicURL:         os.Getenv("S3_PUBLIC_URL"),
		OpenAIAPIKey:        os.Getenv("OPENAI_API_KEY"),
		GeminiAPIKey:        os.Getenv("GEMINI_API_KEY"),
		AnthropicAPIKey:     os.Getenv("ANTHROPIC_API_KEY"),
		CloudflareAPIToken:  os.Getenv("CLOUDFLARE_API_TOKEN"),
		CloudflareAccountID: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		BuilderDirectory:    getEnv("BUILDER_DIRECTORY", "../builder"),
		BuildOutputRoot:     getEnv("BUILD_OUTPUT_ROOT", "/tmp/cms-builder-builds"),
		NPMCommand:          getEnv("NPM_COMMAND", "npm"),
		WranglerCommand:     getEnv("WRANGLER_COMMAND", "wrangler"),
	}
}

func (c Config) Validate() error {
	if _, err := net.LookupPort("tcp", c.APIPort); err != nil {
		return fmt.Errorf("API_PORT must be a valid TCP port: %w", err)
	}
	if c.BuilderDirectory == "" || c.BuildOutputRoot == "" || c.NPMCommand == "" || c.WranglerCommand == "" {
		return fmt.Errorf("builder configuration cannot be empty")
	}
	if info, err := os.Stat(c.BuilderDirectory); err != nil || !info.IsDir() {
		return fmt.Errorf("BUILDER_DIRECTORY must reference a directory")
	}
	if _, err := exec.LookPath(c.NPMCommand); err != nil {
		return fmt.Errorf("NPM_COMMAND is unavailable: %w", err)
	}
	if err := os.MkdirAll(c.BuildOutputRoot, 0o755); err != nil {
		return fmt.Errorf("BUILD_OUTPUT_ROOT is unavailable: %w", err)
	}
	if _, err := filepath.Abs(c.BuildOutputRoot); err != nil {
		return fmt.Errorf("BUILD_OUTPUT_ROOT is invalid: %w", err)
	}
	if (c.CloudflareAPIToken == "") != (c.CloudflareAccountID == "") {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID must be configured together")
	}
	if c.CloudflareAPIToken != "" {
		if _, err := exec.LookPath(c.WranglerCommand); err != nil {
			return fmt.Errorf("WRANGLER_COMMAND is unavailable: %w", err)
		}
	}
	return nil
}

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
