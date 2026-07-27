package storage

import (
	"encoding/json"
	"testing"
)

func TestPublicReadBucketPolicyAllowsOnlyObjectReads(t *testing.T) {
	policy, err := publicReadBucketPolicy("cms-builder")
	if err != nil {
		t.Fatalf("publicReadBucketPolicy returned error: %v", err)
	}

	var document struct {
		Version   string `json:"Version"`
		Statement []struct {
			Effect    string   `json:"Effect"`
			Principal string   `json:"Principal"`
			Action    []string `json:"Action"`
			Resource  []string `json:"Resource"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policy), &document); err != nil {
		t.Fatalf("public read policy must be valid JSON: %v", err)
	}
	if document.Version != "2012-10-17" || len(document.Statement) != 1 {
		t.Fatalf("unexpected policy structure: %s", policy)
	}

	statement := document.Statement[0]
	if statement.Effect != "Allow" || statement.Principal != "*" {
		t.Fatalf("expected public read access, got %#v", statement)
	}
	if len(statement.Action) != 1 || statement.Action[0] != "s3:GetObject" {
		t.Fatalf("expected only object reads to be public, got %#v", statement.Action)
	}
	if len(statement.Resource) != 1 || statement.Resource[0] != "arn:aws:s3:::cms-builder/*" {
		t.Fatalf("expected policy to cover the configured bucket only, got %#v", statement.Resource)
	}
}

func TestNormalizedS3EndpointRemovesAnR2BucketPath(t *testing.T) {
	got := normalizedS3Endpoint("https://account.r2.cloudflarestorage.com/anonime", "anonime")
	if got != "https://account.r2.cloudflarestorage.com" {
		t.Fatalf("unexpected R2 endpoint: %q", got)
	}
}

func TestRequiresBucketPolicySkipsCloudflareR2(t *testing.T) {
	if requiresBucketPolicy("https://account.r2.cloudflarestorage.com") {
		t.Fatal("expected Cloudflare R2 to skip unsupported bucket policy operations")
	}
	if !requiresBucketPolicy("http://minio:9000") {
		t.Fatal("expected MinIO to retain bucket policy operations")
	}
}

func TestUploadObjectKeyUsesAnExplicitKey(t *testing.T) {
	key, err := uploadObjectKey(UploadFile{ObjectKey: "blog/cover.png"})
	if err != nil {
		t.Fatalf("uploadObjectKey returned error: %v", err)
	}
	if key != "blog/cover.png" {
		t.Fatalf("unexpected explicit object key: %q", key)
	}
}
