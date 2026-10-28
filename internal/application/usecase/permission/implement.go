package permission

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/port"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type permissionUseCaseImpl struct {
	logPro         port.LogProvider
	permissionRepo repository.PermissionRepository
}

func NewPermissionUseCase(
	logPro port.LogProvider,
	permissionRepo repository.PermissionRepository,
) PermissionUseCase {
	return &permissionUseCaseImpl{
		logPro,
		permissionRepo,
	}
}

func (u *permissionUseCaseImpl) List(ctx context.Context, schema string) ([]*response.PermissionResponse, error) {
	permissions, err := u.permissionRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all permission failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	permsRes := mapper.ToPermissionsResponse(permissions)

	return permsRes, nil
}
