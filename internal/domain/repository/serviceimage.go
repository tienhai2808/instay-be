package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ServiceImageRepository interface {
	CreateBatch(ctx context.Context, schema string, serviceImages []*model.ServiceImage) error

	DeleteByServiceID(ctx context.Context, schema string, serviceID int64) error
}
