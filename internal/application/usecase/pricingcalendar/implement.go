package pricingcalendar

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	recalculateUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/recalculate"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type pricingCalendarUseCaseImpl struct {
	logPro              port.LogProvider
	txPro               port.TxProvider
	pricingCalendarRepo repository.PricingCalendarRepository
	roomTypeRepo        repository.RoomTypeRepository
	recalculateUC       recalculateUC.RecalculateUseCase
}

func NewPricingCalendarUseCase(
	logPro port.LogProvider,
	txPro port.TxProvider,
	pricingCalendarRepo repository.PricingCalendarRepository,
	roomTypeRepo repository.RoomTypeRepository,
	recalculateUC recalculateUC.RecalculateUseCase,
) PricingCalendarUseCase {
	return &pricingCalendarUseCaseImpl{
		logPro,
		txPro,
		pricingCalendarRepo,
		roomTypeRepo,
		recalculateUC,
	}
}

func (u *pricingCalendarUseCaseImpl) List(ctx context.Context, schema string) ([]*response.PricingCalendarResponse, error) {
	pricingCalendars, err := u.pricingCalendarRepo.FindAllWithPricingTag(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all pricing calendars failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToPricingCalendarsResponse(pricingCalendars), nil
}

func (u *pricingCalendarUseCaseImpl) UpdateBulk(ctx context.Context, schema string, req request.UpdatePricingCalendarsRequest) error {
	pricingTagIDs := make([]int64, 0, len(req.PricingCalendars))
	pricingTagIDSeen := make(map[int64]struct{}, len(req.PricingCalendars))

	dateToTag := make(map[time.Time]int64)
	calendars := make([]*model.PricingCalendar, 0)

	for _, item := range req.PricingCalendars {
		if _, ok := pricingTagIDSeen[item.PricingTagID]; !ok {
			pricingTagIDSeen[item.PricingTagID] = struct{}{}
			pricingTagIDs = append(pricingTagIDs, item.PricingTagID)
		}

		for _, d := range item.Dates {
			dateOnly := utils.NormalizeDateOnly(d)
			if existingTagID, ok := dateToTag[dateOnly]; ok {
				if existingTagID != item.PricingTagID {
					return exceptions.ErrPricingCalendarDateConflict
				}
				continue
			}
			dateToTag[dateOnly] = item.PricingTagID

			tagID := item.PricingTagID
			calendars = append(calendars, &model.PricingCalendar{
				Date:         dateOnly,
				PricingTagID: &tagID,
			})
		}
	}

	if err := u.txPro.WithinTx(ctx, func(txCtx context.Context) error {
		if err := u.pricingCalendarRepo.UpdateAllByPricingTagIDs(txCtx, schema, pricingTagIDs, map[string]any{"pricing_tag_id": nil}); err != nil {
			u.logPro.Error("Clear pricing calendar tag ids failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if err := u.pricingCalendarRepo.UpsertBulk(txCtx, schema, calendars); err != nil {
			if errors.Is(err, exceptions.ErrPricingTagNotFound) {
				return err
			}
			u.logPro.Error("Update bulk pricing calendar failed", port.String("schema", schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	u.recalculateAllRoomTypesAsync(schema)

	return nil
}

func (u *pricingCalendarUseCaseImpl) recalculateAllRoomTypesAsync(schema string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		roomTypes, err := u.roomTypeRepo.FindAll(ctx, schema, model.RoomTypeQuery{})
		if err != nil {
			u.logPro.Error("Find room types for recalculation failed", port.String("schema", schema), port.Error(err))
			return
		}

		for _, roomType := range roomTypes {
			if roomType == nil {
				continue
			}

			if err = u.recalculateUC.RecalculateDailyRoomTypeRate(ctx, schema, roomType.ID); err != nil {
				u.logPro.Error(
					"recalculate daily room type rates failed",
					port.String("schema", schema),
					port.Int64("room_type_id", roomType.ID),
					port.Error(err),
				)
			}
		}
	}()
}
