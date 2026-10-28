package supportcategory

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type SupportCategoryUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateSupportCategoryRequest) error

	Update(ctx context.Context, schema string, supportCategoryID, userID int64, req request.UpdateSupportCategoryRequest) error

	ListFullForAdmin(ctx context.Context, schema string, query request.SupportCategoryPaginationQuery) ([]*response.SupportCategoryResponse, *response.MetaResponse, error)

	ListBasicForAdmin(ctx context.Context, schema string) ([]*response.BasicSupportCategoryResponse, error)

	Delete(ctx context.Context, schema string, supportCategoryID int64) error

	ListForGuest(ctx context.Context, schema string) ([]*response.BasicSupportCategoryResponse, error)
}
