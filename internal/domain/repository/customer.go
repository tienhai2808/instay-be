package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, schema string, customer *model.Customer) error

	FindAll(ctx context.Context, schema string, query model.CustomerQuery) ([]*model.Customer, error)

	FindAllPaginated(ctx context.Context, schema string, query model.CustomerQuery) ([]*model.Customer, int64, error)
}
