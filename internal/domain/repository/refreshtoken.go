package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, schema string, token *model.RefreshToken) error

	Update(ctx context.Context, schema, token string, updateData map[string]any) error

	Delete(ctx context.Context, schema, token string) error

	FindByToken(ctx context.Context, schema, token string) (*model.RefreshToken, error)

	DeleteAllByUserID(ctx context.Context, schema string, userID int64) error

	DeleteAllByUserIDs(ctx context.Context, schema string, userIDs []int64) error

	DeleteAllExpired(ctx context.Context, schema string) (int64, error)
}
