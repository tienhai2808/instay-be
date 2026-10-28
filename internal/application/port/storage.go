package port

import "context"

type PresignedUploadInput struct {
	Key         string
	ContentType string
	ExpiresIn   int64
}

type PresignedViewInput struct {
	Key       string
	ExpiresIn int64
}

type StorageProvider interface {
	GenerateUploadURL(ctx context.Context, input PresignedUploadInput) (string, error)

	GenerateViewURL(ctx context.Context, input PresignedViewInput) (string, error)

	Delete(ctx context.Context, key string) error

	DeleteBatch(ctx context.Context, keys []string) error
}
