package devicetoken

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DeviceTokenUseCase interface {
	Register(ctx context.Context, schema string, clientID int64, clientType model.DeviceTokenClientType, req request.DeviceTokenRequest) error

	Delete(ctx context.Context, schema string, req request.DeviceTokenRequest) error
}
