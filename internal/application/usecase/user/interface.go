package usecase

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type UserUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateUserRequest) (int64, error)

	GetByID(ctx context.Context, schema string, userID int64) (*response.UserDetailsResponse, error)

	List(ctx context.Context, schema string, query request.UserPaginationQuery) ([]*response.SimpleUserResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, userID, currentUserID int64, req request.UpdateUserRequest) error

	UpdatePassword(ctx context.Context, schema string, userID, currentUserID int64, newPassword string) error

	Delete(ctx context.Context, schema string, userID, currentUserID int64) error

	DeleteBatch(ctx context.Context, schema string, currentUserID int64, userIDs []int64) (int64, error)
}
