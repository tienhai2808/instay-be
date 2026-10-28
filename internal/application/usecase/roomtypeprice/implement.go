package roomtypeprice

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	recalculateUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/recalculate"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type roomTypePriceUseCaseImpl struct {
	logPro            port.LogProvider
	idPro             port.IDProvider
	roomTypePriceRepo repository.RoomTypePriceRepository
	recalculateUC     recalculateUC.RecalculateUseCase
}

func NewRoomTypePriceUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	roomTypePriceRepo repository.RoomTypePriceRepository,
	recalculateUC recalculateUC.RecalculateUseCase,
) RoomTypePriceUseCase {
	return &roomTypePriceUseCaseImpl{
		logPro,
		idPro,
		roomTypePriceRepo,
		recalculateUC,
	}
}

func (u *roomTypePriceUseCaseImpl) Create(ctx context.Context, schema string, req request.CreateRoomTypePriceRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate room type price id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	roomTypePrice := &model.RoomTypePrice{
		ID:           id,
		RoomTypeID:   req.RoomTypeID,
		PricingTagID: req.PricingTagID,
		PricingType:  model.RoomTypePricePricingType(req.PricingType),
		PricingValue: req.PricingValue,
	}

	if err := u.roomTypePriceRepo.Create(ctx, schema, roomTypePrice); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypePriceAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrPricingTagNotFound) {
			return err
		}
		u.logPro.Error("Create room type price failed", port.String("schema", schema), port.Error(err))
		return err
	}

	u.recalculateDailyRatesAsync(schema, req.RoomTypeID)

	return nil
}

func (u *roomTypePriceUseCaseImpl) Update(ctx context.Context, schema string, roomTypePriceID int64, req request.UpdateRoomTypePriceRequest) error {
	existingRoomTypePrice, err := u.roomTypePriceRepo.FindByID(ctx, schema, roomTypePriceID)
	if err != nil {
		u.logPro.Error("Find room type price failed", port.String("schema", schema), port.Int64("room_type_price_id", roomTypePriceID), port.Error(err))
		return err
	}
	if existingRoomTypePrice == nil {
		return exceptions.ErrRoomTypePriceNotFound
	}

	updateData := map[string]any{
		"pricing_type":  req.PricingType,
		"pricing_value": req.PricingValue,
	}

	if err = u.roomTypePriceRepo.Update(ctx, schema, roomTypePriceID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypePriceNotFound) {
			return err
		}
		u.logPro.Error("Update room type price failed", port.String("schema", schema), port.Int64("room_type_price_id", roomTypePriceID), port.Error(err))
		return err
	}

	u.recalculateDailyRatesAsync(schema, existingRoomTypePrice.RoomTypeID)

	return nil
}

func (u *roomTypePriceUseCaseImpl) Delete(ctx context.Context, schema string, roomTypePriceID int64) error {
	existingRoomTypePrice, err := u.roomTypePriceRepo.FindByID(ctx, schema, roomTypePriceID)
	if err != nil {
		u.logPro.Error("Find room type price failed", port.String("schema", schema), port.Int64("room_type_price_id", roomTypePriceID), port.Error(err))
		return err
	}
	if existingRoomTypePrice == nil {
		return exceptions.ErrRoomTypePriceNotFound
	}

	if err = u.roomTypePriceRepo.Delete(ctx, schema, roomTypePriceID); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypePriceNotFound) {
			return err
		}
		u.logPro.Error("Delete room type price failed", port.String("schema", schema), port.Int64("room_type_price_id", roomTypePriceID), port.Error(err))
		return err
	}

	u.recalculateDailyRatesAsync(schema, existingRoomTypePrice.RoomTypeID)

	return nil
}

func (u *roomTypePriceUseCaseImpl) recalculateDailyRatesAsync(schema string, roomTypeID int64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := u.recalculateUC.RecalculateDailyRoomTypeRate(ctx, schema, roomTypeID); err != nil {
			u.logPro.Error(
				"recalculate daily room type rates failed",
				port.String("schema", schema),
				port.Int64("room_type_id", roomTypeID),
				port.Error(err),
			)
		}
	}()
}
