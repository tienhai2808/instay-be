package recalculate

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/domain/service"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type recalculateUseCaseImpl struct {
	logPro                port.LogProvider
	pricingService        service.PricingService
	roomTypeRepo          repository.RoomTypeRepository
	roomTypePriceRepo     repository.RoomTypePriceRepository
	pricingCalendarRepo   repository.PricingCalendarRepository
	dailyRoomTypeRateRepo repository.DailyRoomTypeRateRepository
}

func NewRecalculateUseCase(
	logPro port.LogProvider,
	roomTypeRepo repository.RoomTypeRepository,
	roomTypePriceRepo repository.RoomTypePriceRepository,
	pricingCalendarRepo repository.PricingCalendarRepository,
	dailyRoomTypeRateRepo repository.DailyRoomTypeRateRepository,
) RecalculateUseCase {
	return &recalculateUseCaseImpl{
		logPro:                logPro,
		pricingService:        service.PricingService{},
		roomTypeRepo:          roomTypeRepo,
		roomTypePriceRepo:     roomTypePriceRepo,
		pricingCalendarRepo:   pricingCalendarRepo,
		dailyRoomTypeRateRepo: dailyRoomTypeRateRepo,
	}
}

func (u *recalculateUseCaseImpl) RecalculateDailyRoomTypeRate(ctx context.Context, schema string, roomTypeID int64) error {
	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type for recalculation failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}

	if roomType == nil {
		return nil
	}

	dateRangeStart, dateRangeEnd := currentAndNextYearRange(time.Now().UTC())

	if err = u.pricingCalendarRepo.EnsureDateRange(ctx, schema, dateRangeStart, dateRangeEnd); err != nil {
		u.logPro.Error("Ensure pricing calendar date range failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}

	calendars, err := u.pricingCalendarRepo.FindAllByDateRange(ctx, schema, dateRangeStart, dateRangeEnd)
	if err != nil {
		u.logPro.Error("Find pricing calendars for recalculation failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}

	if len(calendars) == 0 {
		return nil
	}

	roomTypePrices, err := u.roomTypePriceRepo.FindAllByRoomTypeID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type prices for recalculation failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}

	roomTypePriceByTagID := make(map[int64]*model.RoomTypePrice, len(roomTypePrices))
	for _, roomTypePrice := range roomTypePrices {
		if roomTypePrice == nil {
			continue
		}
		roomTypePriceByTagID[roomTypePrice.PricingTagID] = roomTypePrice
	}

	rates := make([]*model.DailyRoomTypeRate, 0, len(calendars))
	for _, calendar := range calendars {
		if calendar == nil {
			continue
		}

		finalPrice := roomType.Price
		if calendar.PricingTagID != nil {
			if roomTypePrice, ok := roomTypePriceByTagID[*calendar.PricingTagID]; ok {
				finalPrice = u.pricingService.CalculateFinalPrice(roomType.Price, roomTypePrice.PricingType, roomTypePrice.PricingValue)
			}
		}

		rates = append(rates, &model.DailyRoomTypeRate{
			RoomTypeID:          roomTypeID,
			PricingCalendarDate: utils.NormalizeDateOnly(calendar.Date),
			FinalPrice:          finalPrice,
		})
	}

	if err := u.dailyRoomTypeRateRepo.UpsertBulk(ctx, schema, rates); err != nil {
		u.logPro.Error("Upsert daily room type rates failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}

	return nil
}

func currentAndNextYearRange(now time.Time) (time.Time, time.Time) {
	startDate := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year()+2, time.January, 1, 0, 0, 0, 0, time.UTC)
	return startDate, endDate
}
