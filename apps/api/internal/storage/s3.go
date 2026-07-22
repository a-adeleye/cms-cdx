package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"cms-builder/api/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string
	enabled   bool
}

var unsafeObjectName = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func NewFromConfig(cfg config.Config) StorageProvider {
	return newS3Storage(cfg.S3Endpoint, cfg.S3Region, cfg.S3Bucket, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3PublicURL)
}

func newS3Storage(endpoint, region, bucket, accessKey, secretKey, publicURL string) StorageProvider {
	if strings.TrimSpace(endpoint) == "" ||
		strings.TrimSpace(bucket) == "" ||
		strings.TrimSpace(accessKey) == "" ||
		strings.TrimSpace(secretKey) == "" ||
		strings.TrimSpace(publicURL) == "" {
		return disabledStorage{reason: "storage is not configured"}
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return disabledStorage{reason: err.Error()}
	}

	client := s3.NewFromConfig(awsCfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.UsePathStyle = true
	})

	return &s3Storage{
		client:    client,
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
		enabled:   true,
	}
}

func (s *s3Storage) Upload(ctx context.Context, file UploadFile) (*StoredFile, error) {
	if !s.enabled {
		return nil, errors.New("storage is not configured")
	}
	if len(file.Contents) == 0 {
		return nil, errors.New("file contents are empty")
	}
	if err := s.ensureBucket(ctx); err != nil {
		return nil, err
	}

	key := storageObjectKey(file.SiteID, file.FileName)
	contentType := strings.TrimSpace(file.MimeType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file.Contents),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	return &StoredFile{
		Key:       key,
		PublicURL: s.GetPublicURL(key),
	}, nil
}

func (s *s3Storage) ensureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		_, createErr := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(s.bucket),
		})
		if createErr != nil {
			return fmt.Errorf("unable to prepare storage bucket: %w", createErr)
		}
	}

	policy, err := publicReadBucketPolicy(s.bucket)
	if err != nil {
		return fmt.Errorf("create public media policy: %w", err)
	}
	if _, err := s.client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(s.bucket),
		Policy: aws.String(policy),
	}); err != nil {
		return fmt.Errorf("enable public reads for media bucket: %w", err)
	}
	return nil
}

func publicReadBucketPolicy(bucket string) (string, error) {
	policy := struct {
		Version   string `json:"Version"`
		Statement []struct {
			Effect    string   `json:"Effect"`
			Principal string   `json:"Principal"`
			Action    []string `json:"Action"`
			Resource  []string `json:"Resource"`
		} `json:"Statement"`
	}{
		Version: "2012-10-17",
	}
	policy.Statement = append(policy.Statement, struct {
		Effect    string   `json:"Effect"`
		Principal string   `json:"Principal"`
		Action    []string `json:"Action"`
		Resource  []string `json:"Resource"`
	}{
		Effect:    "Allow",
		Principal: "*",
		Action:    []string{"s3:GetObject"},
		Resource:  []string{fmt.Sprintf("arn:aws:s3:::%s/*", bucket)},
	})

	encoded, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (s *s3Storage) Delete(ctx context.Context, key string) error {
	if !s.enabled {
		return errors.New("storage is not configured")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(strings.TrimSpace(key)),
	})
	return err
}

func (s *s3Storage) GetPublicURL(key string) string {
	if !s.enabled || s.publicURL == "" {
		return ""
	}

	return s.publicURL + "/" + strings.TrimLeft(key, "/")
}

type disabledStorage struct {
	reason string
}

func (d disabledStorage) Upload(ctx context.Context, file UploadFile) (*StoredFile, error) {
	return nil, errors.New(d.reason)
}

func (d disabledStorage) Delete(ctx context.Context, key string) error {
	return errors.New(d.reason)
}

func (d disabledStorage) GetPublicURL(key string) string {
	return ""
}

func storageObjectKey(siteID, fileName string) string {
	stem := strings.TrimSpace(fileName)
	if stem == "" {
		stem = "upload"
	}
	stem = path.Base(stem)
	stem = unsafeObjectName.ReplaceAllString(stem, "-")
	stem = strings.Trim(stem, "-_.")
	if stem == "" {
		stem = "upload"
	}

	return path.Join(strings.TrimSpace(siteID), "media", fmt.Sprintf("%d-%s", time.Now().UTC().UnixNano(), stem))
}
