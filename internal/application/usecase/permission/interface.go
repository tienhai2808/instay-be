package permission

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type PermissionUseCase interface {
	List(ctx context.Context, schema string) ([]*response.PermissionResponse, error)
}
