package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToRefreshTokenEntity(rt *model.RefreshToken) *entity.RefreshTokenEntity {
	if rt == nil {
		return nil
	}

	return &entity.RefreshTokenEntity{
		UserID:    rt.UserID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
	}
}

func ToDeviceTokenEntity(dt *model.DeviceToken) *entity.DeviceTokenEntity {
	if dt == nil {
		return nil
	}

	return &entity.DeviceTokenEntity{
		ClientID:   dt.ClientID,
		ClientType: string(dt.ClientType),
		Token:      dt.Token,
		LastUsedAt: dt.LastUsedAt,
	}
}

func ToRefreshToken(rt *entity.RefreshTokenEntity) *model.RefreshToken {
	if rt == nil {
		return nil
	}

	return &model.RefreshToken{
		Token:     rt.Token,
		UserID:    rt.UserID,
		ExpiresAt: rt.ExpiresAt,
		User:      ToUser(rt.User),
	}
}

func ToRefreshTokens(rts []*entity.RefreshTokenEntity) []*model.RefreshToken {
	if len(rts) == 0 {
		return make([]*model.RefreshToken, 0)
	}

	rtsRes := make([]*model.RefreshToken, 0, len(rts))
	for _, rt := range rts {
		rtsRes = append(rtsRes, ToRefreshToken(rt))
	}

	return rtsRes
}

func ToDeviceToken(dt *entity.DeviceTokenEntity) *model.DeviceToken {
	if dt == nil {
		return nil
	}

	return &model.DeviceToken{
		ClientID:   dt.ClientID,
		Token:      dt.Token,
		ClientType: model.DeviceTokenClientType(dt.ClientType),
		LastUsedAt: dt.LastUsedAt,
	}
}

func ToDeviceTokens(dts []*entity.DeviceTokenEntity) []*model.DeviceToken {
	if len(dts) == 0 {
		return make([]*model.DeviceToken, 0)
	}

	dtsRes := make([]*model.DeviceToken, 0, len(dts))
	for _, dt := range dts {
		dtsRes = append(dtsRes, ToDeviceToken(dt))
	}

	return dtsRes
}
