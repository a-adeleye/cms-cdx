package config

import "os"

type Config struct {
	DatabaseURL string
	JWTSecret   string
	APIPort     string

	S3Endpoint   string
	S3Region     string
	S3Bucket     string
	S3AccessKey  string
	S3SecretKey  string
	S3PublicURL  string
	OpenAIAPIKey string
	GeminiAPIKey string
	AnthropicAPIKey string
}

func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-only-secret"),
		APIPort:     getEnv("API_PORT", "8080"),
		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3PublicURL: os.Getenv("S3_PUBLIC_URL"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
	}
}

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

