package initialization

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

type Storage struct {
	gcs *storage.Client
}

func InitStorage(cfg config.StorageConfig) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Storage{client}, nil
}

func (s *Storage) Client() *storage.Client {
	return s.gcs
}

func (s *Storage) Close() {
	_ = s.gcs.Close()
}
