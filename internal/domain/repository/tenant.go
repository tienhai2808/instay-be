package repository

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type TenantRepository interface {
	FindByDomain(ctx context.Context, domain string) (*model.Tenant, error)

	FindAll(ctx context.Context) ([]*model.Tenant, error)

	FindActive(ctx context.Context) ([]*model.Tenant, error)

	FindBySchemaName(ctx context.Context, schema string) (*model.Tenant, error)

	UpdateBySchemaName(ctx context.Context, schema string, updateData map[string]any) error
}
