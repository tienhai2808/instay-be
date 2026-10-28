package devicetoken

import (
	"context"
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/application/port"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type deviceTokenUseCaseImpl struct {
	logPro          port.LogProvider
	deviceTokenRepo repository.DeviceTokenRepository
}

func NewDeviceTokenUseCase(
	logPro port.LogProvider,
	deviceTokenRepo repository.DeviceTokenRepository,
) DeviceTokenUseCase {
	return &deviceTokenUseCaseImpl{
		logPro,
		deviceTokenRepo,
	}
}

func (u *deviceTokenUseCaseImpl) Register(ctx context.Context, schema string, clientID int64, clientType model.DeviceTokenClientType, req request.DeviceTokenRequest) error {
	token := req.DeviceToken
	exists, err := u.deviceTokenRepo.ExistsByToken(ctx, schema, token)
	if err != nil {
		u.logPro.Error("Check device token by token failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if exists {
		if err = u.deviceTokenRepo.Update(ctx, schema, token, map[string]any{"client_id": clientID}); err != nil {
			if errors.Is(err, exceptions.ErrDeviceTokenNotFound) {
				return err
			}
			u.logPro.Error("Update device token by token failed", port.String("schema", schema), port.Error(err))
			return err
		}
	} else {
		deviceToken := &model.DeviceToken{
			ClientID:   clientID,
			ClientType: clientType,
			Token:      token,
		}
		if err = u.deviceTokenRepo.Create(ctx, schema, deviceToken); err != nil {
			u.logPro.Error("Create device token failed", port.String("schema", schema), port.Error(err))
			return err
		}
	}

	return nil
}

func (u *deviceTokenUseCaseImpl) Delete(ctx context.Context, schema string, req request.DeviceTokenRequest) error {
	if err := u.deviceTokenRepo.Delete(ctx, schema, req.DeviceToken); err != nil {
		if errors.Is(err, exceptions.ErrDeviceTokenNotFound) {
			return err
		}
		u.logPro.Error("Delete device token by token failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}
