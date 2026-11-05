package initialization

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	s3 *s3.Client
}

func InitStorage(cfg config.StorageConfig) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		}
	})
	return &Storage{s3: client}, nil
}

func (s *Storage) Client() *s3.Client {
	return s.s3
}
