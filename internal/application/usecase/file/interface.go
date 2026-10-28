package usecase

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type FileUseCase interface {
	CreateUploadURLs(ctx context.Context, schema string, req request.UploadPresignedURLsRequest) ([]*response.UploadPresignedURLResponse, error)

	CreateViewURLs(ctx context.Context, schema string, req request.ViewPresignedURLsRequest) ([]*response.ViewPresignedURLResponse, error)
}
