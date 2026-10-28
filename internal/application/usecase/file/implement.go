package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type fileUseCaseImpl struct {
	logPro     port.LogProvider
	storagePro port.StorageProvider
}

func NewFileUseCase(
	logPro port.LogProvider,
	storagePro port.StorageProvider,
) FileUseCase {
	return &fileUseCaseImpl{
		logPro,
		storagePro,
	}
}

func (u *fileUseCaseImpl) CreateUploadURLs(ctx context.Context, schema string, req request.UploadPresignedURLsRequest) ([]*response.UploadPresignedURLResponse, error) {
	result := make([]*response.UploadPresignedURLResponse, 0, len(req.Files))

	for _, file := range req.Files {
		name := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)

		key := fmt.Sprintf("%s/%s-%s%s", schema, utils.GenerateRandom(12), utils.GenerateSlug(name), ext)
		url, err := u.storagePro.GenerateUploadURL(ctx, port.PresignedUploadInput{
			Key:         key,
			ContentType: file.ContentType,
			ExpiresIn:   900,
		})
		if err != nil {
			u.logPro.Error("Generate upload presigned URL failed",
				port.String("schema", schema),
				port.String("content_type", file.ContentType),
				port.Error(err))
			return nil, err
		}

		result = append(result, &response.UploadPresignedURLResponse{
			Key: key,
			Url: url,
		})
	}

	return result, nil
}

func (u *fileUseCaseImpl) CreateViewURLs(ctx context.Context, schema string, req request.ViewPresignedURLsRequest) ([]*response.ViewPresignedURLResponse, error) {
	result := make([]*response.ViewPresignedURLResponse, 0, len(req.Keys))

	for _, key := range req.Keys {
		url, err := u.storagePro.GenerateViewURL(ctx, port.PresignedViewInput{
			Key:       key,
			ExpiresIn: 900,
		})
		if err != nil {
			u.logPro.Error("Generate view presigned URL failed", port.String("schema", schema), port.Error(err))
			return nil, err
		}

		result = append(result, &response.ViewPresignedURLResponse{
			Url: url,
		})
	}

	return result, nil
}
