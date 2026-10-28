package repository

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DailyRoomTypeRateRepository interface {
	UpsertBulk(ctx context.Context, schema string, dailyRoomTypeRates []*model.DailyRoomTypeRate) error

	SumFinalPriceByRoomTypeIDAndDateRange(ctx context.Context, schema string, roomTypeID int64, fromDate, toDate time.Time) (float64, int64, error)
}
