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
