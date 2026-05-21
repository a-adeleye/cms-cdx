package storage

import "context"

type UploadFile struct {
	FileName string
	Contents []byte
	MimeType string
	SiteID   string
}

type StoredFile struct {
	Key       string
	PublicURL string
}

type StorageProvider interface {
	Upload(ctx context.Context, file UploadFile) (*StoredFile, error)
	Delete(ctx context.Context, key string) error
	GetPublicURL(key string) string
}

type NoopStorage struct{}

func (n NoopStorage) Upload(ctx context.Context, file UploadFile) (*StoredFile, error) {
	return &StoredFile{Key: file.FileName, PublicURL: ""}, nil
}

func (n NoopStorage) Delete(ctx context.Context, key string) error { return nil }

func (n NoopStorage) GetPublicURL(key string) string { return "" }

