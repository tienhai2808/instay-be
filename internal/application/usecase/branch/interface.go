package branch

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type BranchUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateBranchRequest) error

	ListBasic(ctx context.Context, schema string) ([]*response.BasicBranchResponse, error)

	ListFull(ctx context.Context, schema string, query request.BranchPaginationQuery) ([]*response.BranchResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, branchID, userID int64, req request.UpdateBranchRequest) error

	Delete(ctx context.Context, schema string, branchID int64) error
}
