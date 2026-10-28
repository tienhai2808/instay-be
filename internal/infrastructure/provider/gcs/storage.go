package gcs

import (
	"context"
	"errors"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"golang.org/x/sync/errgroup"
)

type storageProviderImpl struct {
	cfg    config.StorageConfig
	client *storage.Client
}

func NewStorageProvider(
	cfg config.StorageConfig,
	client *storage.Client,
) port.StorageProvider {
	return &storageProviderImpl{
		cfg,
		client,
	}
}

func (s *storageProviderImpl) GenerateUploadURL(ctx context.Context, input port.PresignedUploadInput) (string, error) {
	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		ContentType: input.ContentType,
		Expires:     time.Now().Add(time.Duration(input.ExpiresIn) * time.Second),
		Scheme:      storage.SigningSchemeV4,
	}

	url, err := s.client.Bucket(s.cfg.Bucket).SignedURL(input.Key, opts)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *storageProviderImpl) GenerateViewURL(ctx context.Context, input port.PresignedViewInput) (string, error) {
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(input.ExpiresIn) * time.Second),
		Scheme:  storage.SigningSchemeV4,
	}

	url, err := s.client.Bucket(s.cfg.Bucket).SignedURL(input.Key, opts)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *storageProviderImpl) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	if err := s.client.Bucket(s.cfg.Bucket).Object(key).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}
		return err
	}

	return nil
}

func (s *storageProviderImpl) DeleteBatch(ctx context.Context, keys []string) error {
	filtered := make([]string, 0, len(keys))
	for _, k := range keys {
		if k != "" {
			filtered = append(filtered, k)
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	var mu sync.Mutex
	var failedKeys []string

	g, gCtx := errgroup.WithContext(ctx)
	bucket := s.client.Bucket(s.cfg.Bucket)

	for _, key := range filtered {
		sem <- struct{}{}
		g.Go(func() error {
			defer func() {
				<-sem
			}()

			if err := bucket.Object(key).Delete(gCtx); err != nil {
				if errors.Is(err, storage.ErrObjectNotExist) {
					return nil
				}
				mu.Lock()
				failedKeys = append(failedKeys, key)
				mu.Unlock()
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
