package tenant

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type TenantUseCase interface {
	GetBySchemaName(ctx context.Context, schema string) (*response.BasicTenantResponse, error)

	GetSettingsBySchemaName(ctx context.Context, schema string) (*response.TenantSettingsResponse, error)

	UpdateBySchemaName(ctx context.Context, schema string, req request.UpdateTenantRequest) error
	UpdateSettingsBySchemaName(ctx context.Context, schema string, req request.UpdateTenantSettingsRequest) error
}
