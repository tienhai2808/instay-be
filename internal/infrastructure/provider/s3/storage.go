package s3

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
)

type storageProviderImpl struct {
	bucket  string
	client  *s3.Client
	presign *s3.PresignClient
}

func NewStorageProvider(cfg config.StorageConfig, client *s3.Client) port.StorageProvider {
	return &storageProviderImpl{
		bucket:  cfg.Bucket,
		client:  client,
		presign: s3.NewPresignClient(client),
	}
}

func (s *storageProviderImpl) GenerateUploadURL(ctx context.Context, input port.PresignedUploadInput) (string, error) {
	request, err := s.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(input.Key),
		ContentType: aws.String(input.ContentType),
	}, s3.WithPresignExpires(time.Duration(input.ExpiresIn)*time.Second), keepUploadContentType)
	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func keepUploadContentType(options *s3.PresignOptions) {
	options.ClientOptions = append(options.ClientOptions, func(options *s3.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			// The browser supplies the body later. Preserve Content-Type in the signature
			// even though the SDK sees a bodyless request during presigning.
			_, err := stack.Build.Remove("RemoveContentTypeHeader")
			return err
		})
	})
}

func (s *storageProviderImpl) GenerateViewURL(ctx context.Context, input port.PresignedViewInput) (string, error) {
	request, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(input.Key),
	}, s3.WithPresignExpires(time.Duration(input.ExpiresIn)*time.Second))
	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func (s *storageProviderImpl) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *storageProviderImpl) DeleteBatch(ctx context.Context, keys []string) error {
	objects := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		if key != "" {
			objects = append(objects, types.ObjectIdentifier{Key: aws.String(key)})
		}
	}

	const maxDeleteObjects = 1000
	var failures []error
	for batch := range slices.Chunk(objects, maxDeleteObjects) {
		result, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &types.Delete{Objects: batch, Quiet: aws.Bool(true)},
		})
		if err != nil {
			return errors.Join(append(failures, err)...)
		}

		for _, failure := range result.Errors {
			failures = append(failures, fmt.Errorf("delete object %q: %s: %s",
				aws.ToString(failure.Key), aws.ToString(failure.Code), aws.ToString(failure.Message)))
		}
	}
	return errors.Join(failures...)
}
